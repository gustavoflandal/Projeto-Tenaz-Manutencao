package api

import (
	"errors"
	"net"
	"net/http"
	"slices"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"

	"tenaz/internal/plataforma/auditoria"
	"tenaz/internal/plataforma/autenticacao"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/erros"
	"tenaz/internal/plataforma/obs"
	"tenaz/internal/plataforma/textos"
	"tenaz/internal/plataforma/versao"
)

const (
	cabecalhoRequestID    = "X-Request-Id"
	cabecalhoCSRF         = "X-CSRF-Token"
	cabecalhoDispositivo  = "X-Tenaz-Dispositivo-Id"
	cabecalhoTela         = "X-Tenaz-Tela"
	cabecalhoVersaoTenaz  = "Tenaz-Versao"
	tamanhoMaxRequestID   = 128
	prefixoAutorizacaoApp = "Bearer "
)

// envolver aplica, a toda requisição, o que vem antes do roteador:
// identificador de rastreio, versão do produto na resposta, recusa de
// empresa_id na consulta, contexto da auditoria e recuperação de pânico.
func envolver(proximo http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(cabecalhoRequestID)
		if !requestIDValido(id) {
			id = uuid.Must(uuid.NewV7()).String()
		}
		r.Header.Set(cabecalhoRequestID, id)
		w.Header().Set(cabecalhoRequestID, id)
		// Tenaz-Versao em toda resposta, inclusive de erro (seção 18.8).
		w.Header().Set(cabecalhoVersaoTenaz, versao.Publico().Versao)

		defer func() {
			if v := recover(); v != nil {
				if v == http.ErrAbortHandler {
					panic(v)
				}
				obs.Logger(r.Context()).Error("pânico na requisição", "valor", v, "request_id", id)
				escreverProblema(w, r, erros.Interno)
			}
		}()

		// Não existe parâmetro de empresa (seção 4.7): a empresa vem só da
		// credencial. No corpo, a validação do contrato faz a mesma recusa.
		if r.URL.Query().Has("empresa_id") {
			escreverProblema(w, r, erros.CampoNaoPermitido.ComChave("detalhe.empresa_id_requisicao", nil).
				ComCampoChave("empresa_id", erros.CampoValorInvalido, "validacao.empresa_id", nil))
			return
		}

		req := db.Requisicao{
			Origem:    db.OrigemWeb,
			IP:        ipDe(r),
			UserAgent: r.UserAgent(),
			RequestID: id,
		}
		if strings.HasPrefix(r.Header.Get("Authorization"), prefixoAutorizacaoApp) {
			req.Origem = db.OrigemApp
			if d, err := uuid.Parse(r.Header.Get(cabecalhoDispositivo)); err == nil {
				req.DispositivoID = &d
			}
		}
		// Tela de origem (RF-ADM-011): a web diz em qual tela o usuário está.
		// Vale só se estiver no catálogo fechado; qualquer outro valor é
		// ignorado, e o evento grava só o módulo. A tela é rótulo da trilha,
		// nunca decide acesso.
		if t := strings.TrimSpace(r.Header.Get(cabecalhoTela)); auditoria.TelaValida(t) {
			req.Tela = t
		}
		ctx := db.ComRequisicao(r.Context(), req)
		// Idioma da resposta (ADR-013): o cabeçalho vale desde já; o idioma do
		// usuário e o da empresa entram quando a credencial é resolvida.
		ctx = textos.ComPreferencias(ctx, textos.Preferencias{AcceptLanguage: r.Header.Get("Accept-Language")})
		proximo.ServeHTTP(w, r.WithContext(ctx))
	})
}

func requestIDValido(id string) bool {
	if id == "" || len(id) > tamanhoMaxRequestID {
		return false
	}
	for i := 0; i < len(id); i++ {
		if c := id[i]; c < 0x21 || c > 0x7e {
			return false
		}
	}
	return true
}

// ipDe usa o endereço da conexão. X-Forwarded-For não é lido: atrás do
// balanceador, a leitura do cabeçalho exige a lista de proxies confiáveis
// (pendência de E01 no documento 3, seção 26).
func ipDe(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}

func metodoSeguro(m string) bool {
	return m == http.MethodGet || m == http.MethodHead || m == http.MethodOptions
}

// origemAceita confere Origin (ou, na falta dele, Referer) contra as
// origens da web. obrigatoria: sem os dois cabeçalhos, recusa.
func (s *Servidor) origemAceita(ctx huma.Context, obrigatoria bool) bool {
	o := ctx.Header("Origin")
	if o == "" || o == "null" {
		o = OrigemDe(ctx.Header("Referer"))
	}
	if o == "" {
		return !obrigatoria
	}
	return slices.Contains(s.cfg.OrigensWeb, o)
}

// mwCredencial resolve a credencial, o escopo e a permissão da rota, nesta
// ordem (seção 4.8 do documento 3), antes de ler o corpo.
//
// A sessão é resolvida pelo serviço de autenticação em dois passos, nunca na
// mesma transação (seção 6.3 do documento 5): a credencial no pool da
// autenticação e o ator no pool da aplicação, com a empresa da sessão. O
// caso de uso abre a terceira transação, com o escopo do cadastro.
func (s *Servidor) mwCredencial(ctx huma.Context, next func(huma.Context)) {
	r := rotaDe(ctx.Operation())
	if r == nil {
		next(ctx)
		return
	}
	bearer := strings.HasPrefix(ctx.Header("Authorization"), prefixoAutorizacaoApp)

	// Requisição de navegador que altera dados vinda de outra origem é
	// recusada mesmo em rota pública (login CSRF, criação de conta).
	if !bearer && !metodoSeguro(ctx.Method()) && !s.origemAceita(ctx, false) {
		s.falhar(ctx, erros.CSRFInvalido.ComChave("detalhe.origem_invalida", nil))
		return
	}
	if r.Publica {
		next(ctx)
		return
	}
	if s.aut == nil {
		s.falhar(ctx, erros.ServicoIndisponivel.ComChave("detalhe.autenticacao_nao_configurada", nil))
		return
	}

	var (
		ator   autorizacao.Ator
		err    error
		cookie string
	)
	switch {
	case bearer && r.ApenasSessaoWeb:
		err = erros.NaoAutenticado.ComChave("detalhe.so_sessao_web", nil)
	case bearer:
		token := strings.TrimPrefix(ctx.Header("Authorization"), prefixoAutorizacaoApp)
		ator, err = s.aut.AutenticarApp(ctx.Context(), token)
	default:
		c, errCookie := huma.ReadCookie(ctx, s.cfg.NomeCookie())
		if errCookie != nil || c.Value == "" {
			err = erros.NaoAutenticado.ComChave("detalhe.entre_para_continuar", nil)
			break
		}
		cookie = c.Value
		alteraDados := !metodoSeguro(ctx.Method())
		if alteraDados && !s.origemAceita(ctx, true) {
			err = erros.CSRFInvalido.ComChave("detalhe.csrf_origem", nil)
			break
		}
		ator, err = s.aut.AutenticarWeb(ctx.Context(), cookie, ctx.Header(cabecalhoCSRF), alteraDados)
		if errors.Is(err, erros.CSRFInvalido) {
			err = erros.CSRFInvalido.ComChave("detalhe.csrf_cabecalho", nil)
		}
	}
	if err != nil {
		if errors.Is(err, erros.CredencialInvalida) {
			err = erros.CredencialInvalida.ComChave("detalhe.sessao_invalida", nil)
		}
		s.falhar(ctx, err)
		return
	}

	// Sessão restrita à ativação do MFA (RN-ADM-027).
	if ator.SomenteAtivarMFA && !r.AceitaMFAPendente {
		s.falhar(ctx, autenticacao.ErrMFAPendente())
		return
	}

	// Escopo: rota não declarada para o escopo responde 404 (seção 4.8).
	if !slices.Contains(r.Escopos, ator.Escopo.Tipo) {
		s.falhar(ctx, erros.NaoEncontrado)
		return
	}
	if !r.AlcanceNoServico {
		if err := ator.Exigir(r.Permissoes...); err != nil {
			s.falhar(ctx, err)
			return
		}
	}
	ctx = huma.WithContext(ctx, textos.ComUsuario(ctx.Context(), ator.Idioma, ator.IdiomaEmpresa))
	ctx = huma.WithValue(ctx, chaveAtor{}, ator)
	if cookie != "" {
		ctx = huma.WithValue(ctx, chaveCookie{}, cookie)
	}
	next(ctx)
}

// falhar escreve o erro no formato RFC 9457.
func (s *Servidor) falhar(ctx huma.Context, err error) {
	var e *erros.Erro
	status := http.StatusInternalServerError
	if errors.As(err, &e) {
		status = e.Status
	}
	_ = huma.WriteErr(s.huma, ctx, status, "", err)
}
