package adm_test

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1" //nolint:gosec // TOTP da RFC 6238
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"math/big"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"

	"tenaz/internal/modulos/adm"
	"tenaz/internal/plataforma/armazenamento"
	"tenaz/internal/plataforma/auditoria"
	"tenaz/internal/plataforma/autenticacao"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/db/dbteste"
	"tenaz/internal/plataforma/email"
)

var (
	amb  *dbteste.Ambiente
	auth *autenticacao.Servico
	svc  adm.Servico
	aud  *auditoria.Servico

	verificador     *auditoria.Verificador
	bucket          *armazenamento.S3 // nil: sem armazenamento (motivoSemBucket diz por quê)
	motivoSemBucket string
	caixa           = &caixaDeEntrada{}
)

// caixaDeEntrada é o Enviador dos testes: guarda as mensagens em memória.
type caixaDeEntrada struct {
	mu   sync.Mutex
	msgs []email.Mensagem
}

func (c *caixaDeEntrada) Enviar(_ context.Context, m email.Mensagem) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.msgs = append(c.msgs, m)
	return nil
}

var reToken = regexp.MustCompile(`#token=([A-Za-z0-9_-]{43})`)

// token devolve o segredo do último e-mail para o destino.
func (c *caixaDeEntrada) token(t *testing.T, para string) string {
	t.Helper()
	auth.Aguardar()
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := len(c.msgs) - 1; i >= 0; i-- {
		if c.msgs[i].Para == para {
			if m := reToken.FindStringSubmatch(c.msgs[i].Texto); m != nil {
				return m[1]
			}
		}
	}
	t.Fatalf("nenhum e-mail com link para %s", para)
	return ""
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	var err error
	amb, err = dbteste.Preparar(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	chave := make([]byte, 32)
	_, _ = rand.Read(chave)
	cifrador, _ := autenticacao.NovoCifrador(chave)
	semente := make([]byte, 32)
	_, _ = rand.Read(semente)
	assinador, _ := autenticacao.NovoAssinador("teste", semente, nil)
	// Argon2id com o custo mínimo aceito: o teste confere o fluxo, não o custo.
	hasher, err := autenticacao.NovoHasher(autenticacao.ParametrosArgon2{MemoriaKiB: 8 * 1024, Iteracoes: 1, Paralelo: 1, TamSal: 16, TamChave: 32}, 4)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	auth, err = autenticacao.Novo(autenticacao.ConfigPadrao("http://tenaz.teste"), autenticacao.Dependencias{
		Banco: amb.App, PoolAutenticacao: amb.Auth, Hasher: hasher, Politica: autenticacao.NovaPolitica(),
		Cifrador: cifrador, Assinador: assinador, Email: caixa,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	aud, err = auditoria.Novo(auditoria.Dependencias{Banco: amb.App})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	// O armazenamento do selo é opcional no TestMain: os testes que o exigem
	// se pulam se ele não responde.
	var arm armazenamento.Armazenamento
	if s3, _, errS3 := dbteste.ArmazenamentoS3(ctx); errS3 == nil {
		arm = s3
		bucket = s3
	} else {
		motivoSemBucket = errS3.Error()
	}
	verificador, err = auditoria.NovoVerificador(amb.App, arm)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	svc, err = adm.Novo(adm.Dependencias{Banco: amb.App, Autenticacao: auth, Auditoria: aud, Verificador: verificador})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	codigo := m.Run()
	auth.Aguardar()
	amb.Encerrar()
	os.Exit(codigo)
}

// ctxWeb é o contexto de uma requisição web, como o middleware o montaria.
func ctxWeb() context.Context {
	return db.ComRequisicao(context.Background(), db.Requisicao{
		Origem: db.OrigemWeb, IP: "203.0.113.7", UserAgent: "teste-integracao", RequestID: uuid.NewString(),
	})
}

func cnpjAleatorio() string {
	d := make([]int, 14)
	for i := 0; i < 12; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		d[i] = int(n.Int64())
	}
	dv := func(pesos []int) int {
		s := 0
		for i, p := range pesos {
			s += d[i] * p
		}
		if r := s % 11; r >= 2 {
			return 11 - r
		}
		return 0
	}
	d[12] = dv([]int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2})
	d[13] = dv([]int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2})
	var b strings.Builder
	for _, x := range d {
		b.WriteByte(byte('0' + x))
	}
	return b.String()
}

func cpfAleatorio() string {
	d := make([]int, 11)
	for i := 0; i < 9; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		d[i] = int(n.Int64())
	}
	for _, n := range []int{9, 10} {
		s := 0
		for i := 0; i < n; i++ {
			s += d[i] * (n + 1 - i)
		}
		r := s * 10 % 11
		if r == 10 {
			r = 0
		}
		d[n] = r
	}
	var b strings.Builder
	for _, x := range d {
		b.WriteByte(byte('0' + x))
	}
	return b.String()
}

func emailAleatorio(prefixo string) string {
	return prefixo + "." + strings.ReplaceAll(uuid.NewString()[:13], "-", "") + "@teste.tenaz.dev"
}

const senhaBoa = "Manutenção preditiva 7x"

// turnosComerciais é o modelo "comercial" que a web oferece: segunda a sexta,
// 08:00 às 17:00. O backend não impõe calendário (RN-ADM-019).
func turnosComerciais() []adm.EntradaTurno {
	var t []adm.EntradaTurno
	for d := int32(1); d <= 5; d++ {
		t = append(t, adm.EntradaTurno{DiaSemana: d, HoraInicio: "08:00", HoraFim: "17:00"})
	}
	return t
}

// conta é uma empresa criada pelo fluxo real, com o administrador já logado.
type conta struct {
	criada adm.ContaCriada
	email  string
	adm    autorizacao.Ator
	cookie string
}

func novaConta(t *testing.T, plano string) conta {
	t.Helper()
	ctx := ctxWeb()
	em := emailAleatorio("adm")
	var p *string
	if plano != "" {
		p = &plano
	}
	c, err := svc.CriarConta(ctx, adm.EntradaCriarConta{
		Pais: "BR", IdentificacaoTipo: "cnpj", Identificacao: cnpjAleatorio(), RazaoSocial: "Transportes Teste Ltda", Plano: p,
		Responsavel: adm.EntradaResponsavel{Nome: "Ana Administradora", Email: em},
		Unidade:     adm.EntradaUnidade{Codigo: "MTZ", Nome: "Matriz", Turnos: turnosComerciais()},
	})
	if err != nil {
		t.Fatalf("criar conta: %v", err)
	}
	if !c.ConviteEnviado {
		t.Fatal("o convite do responsável deveria ter saído")
	}
	if _, err := svc.AceitarConvite(ctx, caixa.token(t, em), senhaBoa, nil); err != nil {
		t.Fatalf("aceitar convite: %v", err)
	}
	a, cookie := entrar(t, em, senhaBoa, &c.EmpresaID)
	return conta{criada: c, email: em, adm: a, cookie: cookie}
}

// entrar faz login sem MFA e devolve o Ator da sessão web.
func entrar(t *testing.T, em, senha string, empresa *uuid.UUID) (autorizacao.Ator, string) {
	t.Helper()
	r, err := auth.Login(ctxWeb(), autenticacao.EntradaLogin{Email: em, Senha: senha, EmpresaID: empresa})
	if err != nil {
		t.Fatalf("login de %s: %v", em, err)
	}
	if r.Sessao == nil {
		t.Fatalf("login de %s sem sessão (MFA?)", em)
	}
	a, err := auth.AutenticarWeb(ctxWeb(), r.Sessao.Cookie, r.Sessao.CSRF, true)
	if err != nil {
		t.Fatalf("autenticar sessão: %v", err)
	}
	return a, r.Sessao.Cookie
}

// convidarEEntrar cria o usuário na conta, aceita o convite e faz login.
func convidarEEntrar(t *testing.T, c conta, e adm.EntradaUsuario) (adm.Usuario, autorizacao.Ator) {
	t.Helper()
	ctx := ctxWeb()
	if e.Email == "" {
		e.Email = emailAleatorio(strings.Join(e.Perfis, "-"))
	}
	if e.Nome == "" {
		e.Nome = "Usuário " + strings.Join(e.Perfis, " ")
	}
	u, err := svc.CriarUsuario(ctx, c.adm, e)
	if err != nil {
		t.Fatalf("criar usuário: %v", err)
	}
	if _, err := svc.AceitarConvite(ctx, caixa.token(t, e.Email), senhaBoa, nil); err != nil {
		t.Fatalf("aceitar convite: %v", err)
	}
	a, _ := entrar(t, e.Email, senhaBoa, &c.criada.EmpresaID)
	return u.Usuario, a
}

// codigoTOTP calcula o código atual (RFC 6238) a partir do segredo base32.
func codigoTOTP(t *testing.T, segredoBase32 string, quando time.Time) string {
	t.Helper()
	s, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(segredoBase32)
	if err != nil {
		t.Fatal(err)
	}
	var msg [8]byte
	binary.BigEndian.PutUint64(msg[:], uint64(quando.Unix()/30))
	h := hmac.New(sha1.New, s)
	h.Write(msg[:])
	soma := h.Sum(nil)
	off := soma[len(soma)-1] & 0x0f
	return fmt.Sprintf("%06d", (binary.BigEndian.Uint32(soma[off:off+4])&0x7fffffff)%1_000_000)
}
