package api

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"

	"tenaz/internal/plataforma/autenticacao"
	"tenaz/internal/plataforma/erros"
)

// Rotas de /v1/auth/* e /oauth/* (seções 4.2, 4.3 e 6.1 do documento 3).
// Todas usam o serviço da plataforma; as de convite e de criação de conta
// ficam com o módulo ADM, que é dono do usuário convidado.

// ---------------------------------------------------------------- tipos

// Sessao descreve a sessão aberta: quem é, em que empresa, com que perfis e
// permissões, e as empresas para o seletor.
type Sessao struct {
	Canal        string          `json:"canal" enum:"web,app" doc:"Canal da credencial"`
	CSRFToken    *string         `json:"csrf_token" doc:"Token CSRF da sessão web, para o cabeçalho X-CSRF-Token em POST, PATCH, PUT e DELETE; nulo no app"`
	ExpiraEm     *time.Time      `json:"expira_em" doc:"Teto absoluto da sessão web (7 dias); a inatividade de 12 horas vence antes. Nulo fora do login e da troca de empresa"`
	EscopoAcesso string          `json:"escopo_acesso" enum:"interno,prestador,locatario,link" doc:"Escopo de acesso da credencial (seção 4.9)"`
	Usuario      UsuarioSessao   `json:"usuario"`
	Empresa      EmpresaSessao   `json:"empresa"`
	Perfis       []string        `json:"perfis" doc:"Códigos de origem dos perfis ativos"`
	Permissoes   []string        `json:"permissoes" doc:"Permissões efetivas, em ordem alfabética (perfis e permissões avulsas)"`
	Empresas     []EmpresaAcesso `json:"empresas" doc:"Empresas em que a pessoa entra, para o seletor; a atual vem primeiro"`
	MFAPendente  bool            `json:"mfa_pendente" doc:"A empresa exige segundo fator e a pessoa ainda não o ativou: a sessão só serve para ler a sessão, sair, trocar de empresa e ativar o MFA; as demais rotas respondem 401 nao_autenticado com mfa_pendente (RN-ADM-027)"`
}

// UsuarioSessao é o usuário da sessão, na empresa atual.
type UsuarioSessao struct {
	ID       uuid.UUID  `json:"id" format:"uuid"`
	Nome     string     `json:"nome"`
	Email    string     `json:"email" format:"email"`
	PessoaID *uuid.UUID `json:"pessoa_id" format:"uuid"`
	MFAAtivo bool       `json:"mfa_ativo" doc:"Segundo fator por TOTP ativo"`
}

// EmpresaSessao é a empresa atual. O identificador não aparece: a empresa
// vem da credencial (seção 4.7) e sai só em GET /v1/empresa.
type EmpresaSessao struct {
	Nome          string  `json:"nome"`
	Plano         *string `json:"plano" enum:"essencial,profissional,enterprise" doc:"Nulo na conta só do plano Serviços"`
	SrvHabilitado bool    `json:"srv_habilitado" doc:"Módulo SRV habilitado (conta de subcontratada)"`
}

// EmpresaAcesso é uma empresa do seletor.
type EmpresaAcesso struct {
	ID    uuid.UUID `json:"id" format:"uuid" doc:"Identificador para empresa_escolhida_id"`
	Nome  string    `json:"nome"`
	Atual bool      `json:"atual" doc:"Empresa da credencial atual"`
}

// TokensApp é a resposta de POST /oauth/token (seção 4.3).
type TokensApp struct {
	AccessToken  string `json:"access_token" doc:"JWT EdDSA, 15 minutos"`
	TokenType    string `json:"token_type" enum:"Bearer"`
	ExpiresIn    int64  `json:"expires_in" doc:"Segundos até o access token vencer" example:"900"`
	RefreshToken string `json:"refresh_token" doc:"Refresh token rotativo (tnz_rt_...), 30 dias de inatividade; cada uso devolve outro"`
	Scope        string `json:"scope" enum:"app_campo"`
}

func paraTokens(t autenticacao.TokensApp) TokensApp {
	return TokensApp{AccessToken: t.AccessToken, TokenType: "Bearer", ExpiresIn: t.ExpiraEmSegs, RefreshToken: t.RefreshToken, Scope: "app_campo"}
}

func montarSessao(info autenticacao.InfoSessao) Sessao {
	a := info.Ator
	s := Sessao{
		Canal: a.Canal, EscopoAcesso: a.Escopo.Tipo,
		Usuario: UsuarioSessao{ID: a.UsuarioID, Nome: info.Nome, Email: info.Email, PessoaID: a.PessoaID, MFAAtivo: info.MFAAtivo},
		Perfis:  append([]string{}, a.Perfis...), Permissoes: a.Permissoes.Lista(),
		Empresas: []EmpresaAcesso{}, MFAPendente: info.MFAPendente,
	}
	if s.Canal == "" {
		s.Canal = autenticacao.CanalWeb
	}
	if info.CSRF != "" {
		c := info.CSRF
		s.CSRFToken = &c
	}
	for _, e := range info.Empresas {
		ea := EmpresaAcesso{ID: e.EmpresaID, Nome: e.Nome, Atual: e.EmpresaID == a.EmpresaID}
		if ea.Atual {
			s.Empresa = EmpresaSessao{Nome: e.Nome, Plano: a.Plano, SrvHabilitado: a.SrvHabilitado}
			s.Empresas = append([]EmpresaAcesso{ea}, s.Empresas...)
		} else {
			s.Empresas = append(s.Empresas, ea)
		}
	}
	return s
}

// ------------------------------------------------------- entradas e saídas

type entradaLogin struct {
	Body struct {
		Email              string     `json:"email" maxLength:"254" doc:"E-mail do login"`
		Senha              string     `json:"senha" maxLength:"1024"`
		EmpresaEscolhidaID *uuid.UUID `json:"empresa_escolhida_id,omitempty" format:"uuid" doc:"Empresa da sessão, entre as da pessoa; ausente: a da última sessão"`
	}
}

type saidaSessao struct {
	SetCookie *string `header:"Set-Cookie" doc:"Cookie da sessão web: HttpOnly, SameSite=Lax, Path=/; Secure e prefixo __Host- fora do desenvolvimento"`
	Cache     string  `header:"Cache-Control"`
	Body      Sessao
}

type entradaMFA struct {
	Body struct {
		DesafioID          string     `json:"desafio_id" minLength:"43" maxLength:"43" doc:"desafio_id da resposta 401 mfa_obrigatorio"`
		Codigo             string     `json:"codigo" pattern:"^[0-9]{6}$" doc:"Código de 6 dígitos do aplicativo autenticador"`
		EmpresaEscolhidaID *uuid.UUID `json:"empresa_escolhida_id,omitempty" format:"uuid"`
	}
}

type entradaSemCorpo struct{}

type saidaLogout struct {
	SetCookie *string `header:"Set-Cookie" doc:"Cookie da sessão web: HttpOnly, SameSite=Lax, Path=/; Secure e prefixo __Host- fora do desenvolvimento"`
}

type entradaSessao struct{}

type saidaObterSessao struct {
	Cache string `header:"Cache-Control"`
	Body  Sessao
}

type entradaTrocarEmpresa struct {
	Body struct {
		EmpresaEscolhidaID uuid.UUID `json:"empresa_escolhida_id" format:"uuid" doc:"Empresa de destino, entre as de GET /v1/auth/sessao"`
	}
}

// TrocaEmpresa é a sessão nova; no app, com os tokens novos.
type TrocaEmpresa struct {
	Sessao
	Tokens *TokensApp `json:"tokens" doc:"Tokens novos do app; nulo na web, que recebe o cookie novo"`
}

type saidaTrocarEmpresa struct {
	SetCookie *string `header:"Set-Cookie" doc:"Cookie da sessão web: HttpOnly, SameSite=Lax, Path=/; Secure e prefixo __Host- fora do desenvolvimento"`
	Cache     string  `header:"Cache-Control"`
	Body      TrocaEmpresa
}

type entradaRecuperarSenha struct {
	Body struct {
		Email string `json:"email" maxLength:"254"`
	}
}

type entradaRedefinirSenha struct {
	Body struct {
		Token     string `json:"token" minLength:"43" maxLength:"43" doc:"Token do link do e-mail, lido do fragmento #token="`
		NovaSenha string `json:"nova_senha" maxLength:"1024"`
	}
}

type entradaTrocarSenha struct {
	Body struct {
		SenhaAtual string `json:"senha_atual" maxLength:"1024"`
		NovaSenha  string `json:"nova_senha" maxLength:"1024"`
	}
}

// ConfiguracaoMFA é o segredo TOTP pendente, mostrado uma vez.
type ConfiguracaoMFA struct {
	SegredoBase32 string `json:"segredo_base32" doc:"Segredo para digitar no aplicativo autenticador"`
	URIOTPAuth    string `json:"uri_otpauth" doc:"URI otpauth:// para o QR code"`
}

type saidaIniciarMFA struct {
	Cache string `header:"Cache-Control"`
	Body  ConfiguracaoMFA
}

type entradaCodigoMFA struct {
	Body struct {
		Codigo string `json:"codigo" pattern:"^[0-9]{6}$"`
	}
}

type saidaVazia struct{}

// --------------------------------------------------------------- registro

const semCache = "no-store"

func registrarAutenticacao(s *Servidor) {
	publicoAuth := func(r Rota) Rota {
		r.Tag, r.Modulo, r.Publica, r.SemIdempotencia = TagAutenticacao, "ADM", true, true
		r.Erros = []int{http.StatusUnauthorized, http.StatusForbidden, http.StatusNotFound, http.StatusGone}
		return r
	}
	proprio := func(r Rota) Rota {
		r.Tag, r.Modulo, r.Escopos, r.SemIdempotencia = TagAutenticacao, "ADM", EscoposDoProprioUsuario, true
		return r
	}

	Registrar(s, publicoAuth(Rota{
		ID: "entrar", Metodo: http.MethodPost, Caminho: "/v1/auth/login",
		Resumo: "Entrar com e-mail e senha",
		Descricao: "Cria a sessão web e devolve o cookie de sessão e o csrf_token. Com MFA ativo, responde " +
			"401 mfa_obrigatorio com desafio_id, e o login termina em POST /v1/auth/mfa/verificar. Toda falha de " +
			"credencial responde 401 credencial_invalida, sem dizer se o e-mail existe (RF-ADM-006).",
	}), func(ctx context.Context, in *entradaLogin) (*saidaSessao, error) {
		if s.aut == nil {
			return nil, erros.ServicoIndisponivel
		}
		r, err := s.aut.Login(ctx, autenticacao.EntradaLogin{Email: in.Body.Email, Senha: in.Body.Senha, EmpresaID: in.Body.EmpresaEscolhidaID})
		if errors.Is(err, erros.CredencialInvalida) {
			// A mesma mensagem para e-mail inexistente, senha errada e conta
			// bloqueada: a resposta não revela qual foi (documento 9, 5.3).
			return nil, erros.CredencialInvalida.Com("E-mail ou senha incorretos, ou acesso bloqueado por algum tempo. Confira os dados ou peça um link de redefinição de senha.")
		}
		if err != nil {
			return nil, err
		}
		return s.concluirLogin(ctx, r)
	})

	Registrar(s, publicoAuth(Rota{
		ID: "verificarMFA", Metodo: http.MethodPost, Caminho: "/v1/auth/mfa/verificar",
		Resumo: "Concluir o login com o segundo fator",
		Descricao: "Confere o código TOTP do desafio. Código errado responde 401 mfa_codigo_invalido com " +
			"tentativas_restantes; esgotadas as 5 tentativas, ou vencido o desafio (5 minutos), o login recomeça pela senha.",
	}), func(ctx context.Context, in *entradaMFA) (*saidaSessao, error) {
		if s.aut == nil {
			return nil, erros.ServicoIndisponivel
		}
		r, err := s.aut.VerificarMFA(ctx, autenticacao.EntradaVerificarMFA{
			DesafioID: in.Body.DesafioID, Codigo: in.Body.Codigo, EmpresaID: in.Body.EmpresaEscolhidaID,
		})
		var e *erros.Erro
		if errors.As(err, &e) && errors.Is(e, erros.MFACodigoInvalido) {
			return nil, e.Com("Código incorreto. Confira o código atual do aplicativo autenticador.")
		}
		if errors.Is(err, erros.CredencialInvalida) {
			return nil, erros.CredencialInvalida.Com("O desafio venceu, já foi usado ou esgotou as tentativas; entre de novo com e-mail e senha.")
		}
		if err != nil {
			return nil, err
		}
		return s.concluirLogin(ctx, r)
	})

	Registrar(s, proprio(Rota{
		ID: "sair", Metodo: http.MethodPost, Caminho: "/v1/auth/logout", Status: http.StatusNoContent, AceitaMFAPendente: true,
		Resumo:    "Sair",
		Descricao: "Encerra a sessão no servidor (web ou app; no app, com a família de refresh tokens) e apaga o cookie.",
	}), func(ctx context.Context, _ *entradaSemCorpo) (*saidaLogout, error) {
		a, err := AtorDe(ctx)
		if err != nil {
			return nil, err
		}
		if err := s.aut.Logout(ctx, a); err != nil {
			return nil, err
		}
		out := &saidaLogout{}
		if a.Canal != autenticacao.CanalApp {
			out.SetCookie = s.cookieApagado()
		}
		return out, nil
	})

	Registrar(s, proprio(Rota{
		ID: "obterSessao", Metodo: http.MethodGet, Caminho: "/v1/auth/sessao", AceitaMFAPendente: true,
		Resumo:    "Sessão atual",
		Descricao: "Usuário, empresa, escopo, perfis, permissões, empresas do seletor e, na web, o csrf_token.",
	}), func(ctx context.Context, _ *entradaSessao) (*saidaObterSessao, error) {
		a, err := AtorDe(ctx)
		if err != nil {
			return nil, err
		}
		var info autenticacao.InfoSessao
		if a.Canal == autenticacao.CanalApp {
			info, err = s.aut.ObterSessao(ctx, a)
		} else {
			info, err = s.aut.ObterSessaoWeb(ctx, cookieDoContexto(ctx))
		}
		if err != nil {
			return nil, err
		}
		return &saidaObterSessao{Cache: semCache, Body: montarSessao(info)}, nil
	})

	Registrar(s, proprio(Rota{
		ID: "trocarEmpresa", Metodo: http.MethodPost, Caminho: "/v1/auth/trocar-empresa", AceitaMFAPendente: true,
		Resumo: "Trocar de empresa",
		Descricao: "Emite credencial NOVA na empresa escolhida e revoga a atual: uma credencial nunca vale para duas " +
			"empresas (seção 4.7). Na web, cookie e csrf_token novos; no app, tokens novos. Empresa fora do acesso: 404.",
	}), func(ctx context.Context, in *entradaTrocarEmpresa) (*saidaTrocarEmpresa, error) {
		a, err := AtorDe(ctx)
		if err != nil {
			return nil, err
		}
		r, err := s.aut.TrocarEmpresa(ctx, a, in.Body.EmpresaEscolhidaID)
		if err != nil {
			return nil, err
		}
		out := &saidaTrocarEmpresa{Cache: semCache}
		if r.App != nil {
			novo, err := s.aut.AutenticarApp(ctx, r.App.AccessToken)
			if err != nil {
				return nil, err
			}
			info, err := s.aut.ObterSessao(ctx, novo)
			if err != nil {
				return nil, err
			}
			t := paraTokens(*r.App)
			out.Body = TrocaEmpresa{Sessao: montarSessao(info), Tokens: &t}
			return out, nil
		}
		info, err := s.aut.ObterSessaoWeb(ctx, r.Web.Cookie)
		if err != nil {
			return nil, err
		}
		out.SetCookie = s.cookieSessao(*r.Web)
		out.Body = TrocaEmpresa{Sessao: montarSessao(info)}
		exp := Instante(r.Web.ExpiraEm)
		out.Body.ExpiraEm = &exp
		return out, nil
	})

	Registrar(s, publicoAuth(Rota{
		ID: "recuperarSenha", Metodo: http.MethodPost, Caminho: "/v1/auth/recuperar-senha", Status: http.StatusAccepted,
		Resumo: "Pedir link de redefinição de senha",
		Descricao: "Envia o link (válido por 1 hora) se o e-mail tiver login. Responde 202 igual, exista ou não o " +
			"e-mail, para não revelar quais e-mails existem.",
	}), func(ctx context.Context, in *entradaRecuperarSenha) (*saidaVazia, error) {
		if s.aut == nil {
			return nil, erros.ServicoIndisponivel
		}
		return &saidaVazia{}, s.aut.RecuperarSenha(ctx, in.Body.Email)
	})

	Registrar(s, publicoAuth(Rota{
		ID: "redefinirSenha", Metodo: http.MethodPost, Caminho: "/v1/auth/redefinir-senha", Status: http.StatusNoContent,
		Resumo: "Redefinir a senha com o link do e-mail",
		Descricao: "O token vem do fragmento do link (#token=...) e vai no corpo, nunca na URL. Redefinir encerra " +
			"todas as sessões da pessoa, em todas as empresas. 410 link_email_invalido; 422 senha_fraca.",
	}), func(ctx context.Context, in *entradaRedefinirSenha) (*saidaVazia, error) {
		if s.aut == nil {
			return nil, erros.ServicoIndisponivel
		}
		return &saidaVazia{}, s.aut.RedefinirSenha(ctx, in.Body.Token, in.Body.NovaSenha)
	})

	Registrar(s, proprio(Rota{
		ID: "trocarSenha", Metodo: http.MethodPost, Caminho: "/v1/auth/trocar-senha", Status: http.StatusNoContent,
		Resumo: "Trocar a própria senha",
		Descricao: "Confere a senha atual (errada conta como falha de login) e encerra as outras sessões da pessoa; " +
			"a atual continua. 422 senha_fraca.",
	}), func(ctx context.Context, in *entradaTrocarSenha) (*saidaVazia, error) {
		a, err := AtorDe(ctx)
		if err != nil {
			return nil, err
		}
		return &saidaVazia{}, s.aut.TrocarSenha(ctx, a, in.Body.SenhaAtual, in.Body.NovaSenha)
	})

	Registrar(s, proprio(Rota{
		ID: "iniciarMFA", Metodo: http.MethodPost, Caminho: "/v1/auth/mfa/iniciar", AceitaMFAPendente: true,
		Resumo: "Iniciar a ativação do segundo fator",
		Descricao: "Gera o segredo TOTP pendente e devolve o segredo e a URI do QR code, uma única vez. Só vale depois " +
			"de POST /v1/auth/mfa/confirmar. 409 transicao_invalida se o segundo fator já está ativo.",
	}), func(ctx context.Context, _ *entradaSemCorpo) (*saidaIniciarMFA, error) {
		a, err := AtorDe(ctx)
		if err != nil {
			return nil, err
		}
		c, err := s.aut.IniciarMFA(ctx, a)
		if err != nil {
			return nil, err
		}
		return &saidaIniciarMFA{Cache: semCache, Body: ConfiguracaoMFA{SegredoBase32: c.SegredoBase32, URIOTPAuth: c.URIOTPAuth}}, nil
	})

	Registrar(s, proprio(Rota{
		ID: "confirmarMFA", Metodo: http.MethodPost, Caminho: "/v1/auth/mfa/confirmar", Status: http.StatusNoContent,
		AceitaMFAPendente: true,
		Resumo:            "Confirmar a ativação do segundo fator",
		Descricao:         "Ativa o segundo fator com o primeiro código certo. 401 mfa_codigo_invalido.",
	}), func(ctx context.Context, in *entradaCodigoMFA) (*saidaVazia, error) {
		a, err := AtorDe(ctx)
		if err != nil {
			return nil, err
		}
		return &saidaVazia{}, s.aut.ConfirmarMFA(ctx, a, in.Body.Codigo)
	})

	Registrar(s, proprio(Rota{
		ID: "desativarMFA", Metodo: http.MethodPost, Caminho: "/v1/auth/mfa/desativar", Status: http.StatusNoContent,
		Resumo: "Desativar o segundo fator",
		Descricao: "Exige um código válido. Código errado soma uma falha ao bloqueio progressivo, como no login, e código " +
			"de um passo TOTP já usado vale como errado: 401 mfa_codigo_invalido. 409 transicao_invalida com " +
			"estado_atual \"inativo\" se o segundo fator não está ativo, e com estado_atual \"exigido_pela_empresa\" " +
			"enquanto a identidade tiver acesso a alguma empresa que exige MFA (RN-ADM-027).",
	}), func(ctx context.Context, in *entradaCodigoMFA) (*saidaVazia, error) {
		a, err := AtorDe(ctx)
		if err != nil {
			return nil, err
		}
		return &saidaVazia{}, s.aut.DesativarMFA(ctx, a, in.Body.Codigo)
	})

	registrarOAuth(s)
}

// concluirLogin transforma o resultado do login em resposta: desafio de MFA
// (401) ou sessão com cookie.
func (s *Servidor) concluirLogin(ctx context.Context, r autenticacao.ResultadoLogin) (*saidaSessao, error) {
	if r.DesafioMFA != nil {
		return nil, erros.MFAObrigatorio.
			Com("Informe o código do aplicativo autenticador em POST /v1/auth/mfa/verificar.").
			ComDado("desafio_id", r.DesafioMFA.DesafioID).
			ComDado("desafio_expira_em", r.DesafioMFA.ExpiraEm)
	}
	if r.Sessao == nil {
		return nil, erros.CredencialInvalida
	}
	info, err := s.aut.ObterSessaoWeb(ctx, r.Sessao.Cookie)
	if err != nil {
		return nil, err
	}
	sessao := montarSessao(info)
	exp := Instante(r.Sessao.ExpiraEm)
	sessao.ExpiraEm = &exp
	return &saidaSessao{SetCookie: s.cookieSessao(*r.Sessao), Cache: semCache, Body: sessao}, nil
}

// cookieSessao monta o cookie da sessão web (seção 4.2): HttpOnly,
// SameSite=Lax, Path=/, sem Domain; Secure e prefixo __Host- fora do
// desenvolvimento.
func (s *Servidor) cookieSessao(w autenticacao.SessaoWeb) *string {
	c := (&http.Cookie{
		Name: s.cfg.NomeCookie(), Value: w.Cookie, Path: "/", HttpOnly: true, Secure: s.cfg.CookieSeguro,
		SameSite: http.SameSiteLaxMode, Expires: w.ExpiraEm.UTC(), MaxAge: int(time.Until(w.ExpiraEm) / time.Second),
	}).String()
	return &c
}

func (s *Servidor) cookieApagado() *string {
	c := (&http.Cookie{
		Name: s.cfg.NomeCookie(), Value: "", Path: "/", HttpOnly: true, Secure: s.cfg.CookieSeguro,
		SameSite: http.SameSiteLaxMode, MaxAge: -1,
	}).String()
	return &c
}

// --------------------------------------------------------------- OAuth

type entradaAutorizarGET struct {
	ResponseType        string `query:"response_type" required:"true" enum:"code"`
	ClientID            string `query:"client_id" required:"true" maxLength:"64"`
	RedirectURI         string `query:"redirect_uri" required:"true" maxLength:"512"`
	CodeChallenge       string `query:"code_challenge" required:"true" maxLength:"128"`
	CodeChallengeMethod string `query:"code_challenge_method" required:"true" enum:"S256"`
	State               string `query:"state" maxLength:"512"`
}

type saidaRedirecionar struct {
	Status   int
	Location string `header:"Location"`
	Cache    string `header:"Cache-Control"`
}

// Autorizacao é a resposta de POST /oauth/authorize.
type Autorizacao struct {
	URIRedirecionamento string `json:"uri_redirecionamento" doc:"redirect_uri com code e state (ou error=access_denied), para a web navegar até ela e o sistema entregar ao app"`
}

type entradaAutorizarPOST struct {
	Body struct {
		ResponseType        string    `json:"response_type" enum:"code"`
		ClientID            string    `json:"client_id" maxLength:"64"`
		RedirectURI         string    `json:"redirect_uri" maxLength:"512"`
		CodeChallenge       string    `json:"code_challenge" maxLength:"128"`
		CodeChallengeMethod string    `json:"code_challenge_method" enum:"S256"`
		State               string    `json:"state,omitempty" maxLength:"512"`
		EmpresaEscolhidaID  uuid.UUID `json:"empresa_escolhida_id" format:"uuid" doc:"Empresa da sessão do app, entre as de GET /v1/auth/sessao"`
	}
}

type saidaAutorizacao struct {
	Cache string `header:"Cache-Control"`
	Body  Autorizacao
}

type entradaToken struct {
	Dispositivo string `header:"X-Tenaz-Dispositivo-Id" format:"uuid" doc:"Dispositivo do app, gravado na sessão"`
	Body        struct {
		GrantType    string `json:"grant_type" enum:"authorization_code,refresh_token"`
		Code         string `json:"code,omitempty" maxLength:"64"`
		RedirectURI  string `json:"redirect_uri,omitempty" maxLength:"512"`
		ClientID     string `json:"client_id" maxLength:"64"`
		CodeVerifier string `json:"code_verifier,omitempty" maxLength:"128"`
		RefreshToken string `json:"refresh_token,omitempty" maxLength:"128"`
	} `contentType:"application/x-www-form-urlencoded"`
}

type saidaToken struct {
	Cache  string `header:"Cache-Control"`
	Pragma string `header:"Pragma"`
	Body   TokensApp
}

type entradaRevogar struct {
	Body struct {
		Token         string `json:"token" maxLength:"128"`
		TokenTypeHint string `json:"token_type_hint,omitempty" enum:"refresh_token"`
		ClientID      string `json:"client_id,omitempty" maxLength:"64"`
	} `contentType:"application/x-www-form-urlencoded"`
}

func registrarOAuth(s *Servidor) {
	Registrar(s, Rota{
		ID: "iniciarAutorizacaoApp", Metodo: http.MethodGet, Caminho: "/oauth/authorize", Status: http.StatusFound,
		Tag: TagOAuth, Modulo: "ADM", Publica: true,
		Resumo: "Início do OAuth 2.1 com PKCE do app",
		Descricao: "Confere client_id e redirect_uri e redireciona (302) para a tela de autorização da web, " +
			"/autorizar-app, com os mesmos parâmetros. Lá a pessoa entra com a sessão web normal (senha e MFA), " +
			"escolhe a empresa e a web chama POST /oauth/authorize. Cliente ou URI de retorno inválidos respondem " +
			"422 e nunca redirecionam.",
	}, func(_ context.Context, in *entradaAutorizarGET) (*saidaRedirecionar, error) {
		if s.aut == nil {
			return nil, erros.ServicoIndisponivel
		}
		if err := s.aut.ValidarCliente(in.ClientID, in.RedirectURI); err != nil {
			return nil, err
		}
		q := url.Values{}
		q.Set("response_type", in.ResponseType)
		q.Set("client_id", in.ClientID)
		q.Set("redirect_uri", in.RedirectURI)
		q.Set("code_challenge", in.CodeChallenge)
		q.Set("code_challenge_method", in.CodeChallengeMethod)
		if in.State != "" {
			q.Set("state", in.State)
		}
		return &saidaRedirecionar{Status: http.StatusFound, Cache: semCache, Location: s.cfg.URLWeb + "/autorizar-app?" + q.Encode()}, nil
	})

	Registrar(s, Rota{
		ID: "autorizarApp", Metodo: http.MethodPost, Caminho: "/oauth/authorize",
		Tag: TagOAuth, Modulo: "ADM", ApenasSessaoWeb: true, Escopos: EscoposDoProprioUsuario, SemIdempotencia: true,
		Resumo: "Autorizar o app na empresa escolhida",
		Descricao: "Chamada pela tela /autorizar-app da web, com a sessão web (cookie, X-CSRF-Token e Origin). Emite o " +
			"código de uso único (60 s) e devolve a redirect_uri com code e state. Perfil sem mob.usar_app (o locatario) " +
			"recebe a redirect_uri com error=access_denied. Empresa fora do acesso: 404.",
	}, func(ctx context.Context, in *entradaAutorizarPOST) (*saidaAutorizacao, error) {
		a, err := AtorDe(ctx)
		if err != nil {
			return nil, err
		}
		b := in.Body
		if err := s.aut.ValidarCliente(b.ClientID, b.RedirectURI); err != nil {
			return nil, err
		}
		codigo, err := s.aut.EmitirCodigoAutorizacao(ctx, a, autenticacao.EntradaAutorizacao{
			ClienteID: b.ClientID, URIRetorno: b.RedirectURI, DesafioPKCE: b.CodeChallenge, MetodoPKCE: b.CodeChallengeMethod,
			EmpresaID: b.EmpresaEscolhidaID,
		})
		q := url.Values{}
		switch {
		case errors.Is(err, erros.SemPermissao):
			q.Set("error", "access_denied")
			q.Set("error_description", "o perfil nesta empresa não usa o app")
		case err != nil:
			return nil, err
		default:
			q.Set("code", codigo)
		}
		if b.State != "" {
			q.Set("state", b.State)
		}
		return &saidaAutorizacao{Cache: semCache, Body: Autorizacao{URIRedirecionamento: b.RedirectURI + "?" + q.Encode()}}, nil
	})

	Registrar(s, Rota{
		ID: "emitirTokensApp", Metodo: http.MethodPost, Caminho: "/oauth/token",
		Tag: TagOAuth, Modulo: "ADM", Publica: true, SemIdempotencia: true,
		Resumo: "Trocar o código ou renovar o refresh token",
		Descricao: "grant_type=authorization_code (code, redirect_uri, client_id, code_verifier) ou " +
			"grant_type=refresh_token (refresh_token, client_id). Cada renovação devolve refresh token novo; um " +
			"refresh token já usado e reapresentado revoga a família e a sessão do aparelho. Falha: 401 " +
			"credencial_invalida com error=invalid_grant (RFC 6749). O app serializa as renovações.",
	}, func(ctx context.Context, in *entradaToken) (*saidaToken, error) {
		if s.aut == nil {
			return nil, erros.ServicoIndisponivel
		}
		b := in.Body
		var (
			t   autenticacao.TokensApp
			err error
		)
		switch b.GrantType {
		case "authorization_code":
			var disp *uuid.UUID
			if d, errD := uuid.Parse(in.Dispositivo); errD == nil {
				disp = &d
			}
			t, err = s.aut.TrocarCodigo(ctx, autenticacao.EntradaTrocaCodigo{
				Codigo: b.Code, ClienteID: b.ClientID, URIRetorno: b.RedirectURI, VerificadorPKCE: b.CodeVerifier, DispositivoID: disp,
			})
		default:
			t, err = s.aut.Renovar(ctx, b.RefreshToken, b.ClientID)
		}
		if err != nil {
			if errors.Is(err, erros.CredencialInvalida) {
				return nil, erros.CredencialInvalida.Com("Código ou refresh token inválido, vencido ou já usado; entre de novo.").
					ComDado("error", "invalid_grant")
			}
			return nil, err
		}
		return &saidaToken{Cache: semCache, Pragma: "no-cache", Body: paraTokens(t)}, nil
	})

	Registrar(s, Rota{
		ID: "revogarTokenApp", Metodo: http.MethodPost, Caminho: "/oauth/revoke", Status: http.StatusOK,
		Tag: TagOAuth, Modulo: "ADM", Publica: true, SemIdempotencia: true,
		Resumo: "Revogar o refresh token (sair no app)",
		Descricao: "RFC 7009: revoga a família de refresh tokens e a sessão do aparelho. O refresh token no corpo é a " +
			"credencial. Token desconhecido também responde 200, para não servir de oráculo.",
	}, func(ctx context.Context, in *entradaRevogar) (*saidaVazia, error) {
		if s.aut == nil {
			return nil, erros.ServicoIndisponivel
		}
		return &saidaVazia{}, s.aut.Revogar(ctx, in.Body.Token)
	})
}

// cookieDoContexto devolve o cookie de sessão da requisição atual. O
// middleware guarda o valor no contexto ao autenticar pela web.
func cookieDoContexto(ctx context.Context) string {
	c, _ := ctx.Value(chaveCookie{}).(string)
	return c
}

type chaveCookie struct{}
