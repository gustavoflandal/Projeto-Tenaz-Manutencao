package api

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/danielgtaylor/huma/v2"

	"tenaz/internal/plataforma/erros"
)

// Idempotência de POST (seção 2.9 do documento 3).
//
// A chave do cabeçalho Idempotency-Key vale por empresa e por credencial
// (sessão), por 24 horas. Repetição com o mesmo corpo devolve a resposta
// original com Idempotency-Replayed: true; com outro corpo, 422
// idempotencia_conflito; enquanto a primeira ainda roda, 409
// idempotencia_em_processamento com Retry-After: 1.
//
// O depósito desta etapa fica na memória do processo, o que só é correto com
// uma instância da API. O depósito compartilhado (tabela no banco ou Valkey)
// é pendência registrada no documento 3, seção 26, e precisa existir antes
// de a API rodar com mais de uma instância.

const (
	cabecalhoIdempotencia = "Idempotency-Key"
	retencaoIdempotencia  = 24 * time.Hour
	tamanhoMaxChave       = 64
	limiteCorpoIdem       = 1 << 20
)

type respostaGravada struct {
	status    int
	cabecalho http.Header
	corpo     []byte
}

type registroIdem struct {
	hashCorpo string
	expira    time.Time
	resposta  *respostaGravada // nil: em processamento
}

// Idempotencia é o depósito das chaves.
type Idempotencia struct {
	mu     sync.Mutex
	chaves map[string]*registroIdem
	agora  func() time.Time
}

// NovaIdempotencia cria o depósito em memória.
func NovaIdempotencia() *Idempotencia {
	return &Idempotencia{chaves: map[string]*registroIdem{}, agora: time.Now}
}

// reservar devolve a resposta gravada (repetição), ou reserva a chave.
func (d *Idempotencia) reservar(escopo, hash string) (*respostaGravada, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	agora := d.agora()
	for k, r := range d.chaves {
		if agora.After(r.expira) {
			delete(d.chaves, k)
		}
	}
	if r, ok := d.chaves[escopo]; ok {
		switch {
		case r.hashCorpo != hash:
			return nil, erros.IdempotenciaConflito.ComChave("detalhe.idempotencia_conflito", nil)
		case r.resposta == nil:
			return nil, erros.IdempotenciaEmProcessamento.ComChave("detalhe.idempotencia_processando", nil)
		default:
			return r.resposta, nil
		}
	}
	d.chaves[escopo] = &registroIdem{hashCorpo: hash, expira: agora.Add(retencaoIdempotencia)}
	return nil, nil
}

func (d *Idempotencia) concluir(escopo string, r *respostaGravada) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if reg, ok := d.chaves[escopo]; ok {
		reg.resposta = r
	}
}

func (d *Idempotencia) liberar(escopo string) {
	d.mu.Lock()
	defer d.mu.Unlock()
	delete(d.chaves, escopo)
}

func chaveIdemValida(c string) bool {
	if c == "" || len(c) > tamanhoMaxChave {
		return false
	}
	for i := 0; i < len(c); i++ {
		if b := c[i]; b < 0x21 || b > 0x7e {
			return false
		}
	}
	return true
}

// mwIdempotencia aplica a idempotência aos POST autenticados. Roda depois de
// mwCredencial, porque a chave vale por empresa e por sessão.
func (s *Servidor) mwIdempotencia(ctx huma.Context, next func(huma.Context)) {
	r := rotaDe(ctx.Operation())
	chave := ctx.Header(cabecalhoIdempotencia)
	if r == nil || ctx.Method() != http.MethodPost || chave == "" || r.Publica || r.SemIdempotencia {
		next(ctx)
		return
	}
	if !chaveIdemValida(chave) {
		s.falhar(ctx, erros.RequisicaoInvalida.ComChave("detalhe.idempotencia_formato", nil).
			ComCampoChave("Idempotency-Key", erros.CampoFormatoInvalido, "campo.idempotencia_formato", nil))
		return
	}
	ator, err := AtorDe(ctx.Context())
	if err != nil {
		next(ctx)
		return
	}
	corpo, err := io.ReadAll(io.LimitReader(ctx.BodyReader(), limiteCorpoIdem+1))
	if err != nil {
		s.falhar(ctx, erros.RequisicaoInvalida.ComChave("detalhe.corpo_ilegivel", nil))
		return
	}
	if len(corpo) > limiteCorpoIdem {
		s.falhar(ctx, erros.PayloadMuitoGrande.ComChave("detalhe.corpo_grande", nil))
		return
	}
	u := ctx.URL()
	soma := sha256.Sum256(append([]byte(ctx.Method()+" "+u.Path+"\n"), corpo...))
	hash := hex.EncodeToString(soma[:])
	escopo := ator.EmpresaID.String() + "|" + ator.SessaoID.String() + "|" + chave

	gravada, err := s.idem.reservar(escopo, hash)
	if err != nil {
		if errors.Is(err, erros.IdempotenciaEmProcessamento) {
			ctx.SetHeader("Retry-After", "1")
		}
		s.falhar(ctx, err)
		return
	}
	if gravada != nil {
		for k, vs := range gravada.cabecalho {
			for _, v := range vs {
				ctx.AppendHeader(k, v)
			}
		}
		ctx.SetHeader("Idempotency-Replayed", "true")
		ctx.SetStatus(gravada.status)
		_, _ = ctx.BodyWriter().Write(gravada.corpo)
		return
	}

	g := &gravador{contextoHuma: ctx, corpo: corpo, cabecalho: http.Header{}}
	concluiu := false
	defer func() {
		if !concluiu {
			s.idem.liberar(escopo)
		}
	}()
	next(g)
	if g.status() >= http.StatusInternalServerError {
		return
	}
	s.idem.concluir(escopo, &respostaGravada{status: g.status(), cabecalho: g.cabecalho, corpo: g.saida.Bytes()})
	concluiu = true
}

// contextoHuma existe para o campo embutido não se chamar Context, nome que
// esconderia o método Context() da interface.
type contextoHuma = huma.Context

// gravador repassa a resposta e guarda uma cópia para a repetição.
type gravador struct {
	contextoHuma
	corpo     []byte
	cabecalho http.Header
	saida     bytes.Buffer
	codigo    int
}

// cabeçalhos que a repetição devolve.
var cabecalhosRepetidos = map[string]bool{"Content-Type": true, "Location": true, "Etag": true}

func (g *gravador) BodyReader() io.Reader { return bytes.NewReader(g.corpo) }

func (g *gravador) SetStatus(c int) {
	g.codigo = c
	g.contextoHuma.SetStatus(c)
}

func (g *gravador) status() int {
	if g.codigo == 0 {
		return http.StatusOK
	}
	return g.codigo
}

func (g *gravador) SetHeader(nome, valor string) {
	if cabecalhosRepetidos[http.CanonicalHeaderKey(nome)] {
		g.cabecalho.Set(nome, valor)
	}
	g.contextoHuma.SetHeader(nome, valor)
}

func (g *gravador) AppendHeader(nome, valor string) {
	if cabecalhosRepetidos[http.CanonicalHeaderKey(nome)] {
		g.cabecalho.Add(nome, valor)
	}
	g.contextoHuma.AppendHeader(nome, valor)
}

func (g *gravador) BodyWriter() io.Writer {
	return io.MultiWriter(g.contextoHuma.BodyWriter(), &g.saida)
}
