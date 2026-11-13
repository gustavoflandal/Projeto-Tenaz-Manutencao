package integracao_test

import (
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"
)

// Criar conta → aceitar convite → login → MFA → trocar de empresa → logout.
func TestFluxoContaLoginMFATrocaLogout(t *testing.T) {
	em := emailAleatorio("fluxo")
	novaConta(t, "essencial", em)

	// O mesmo e-mail é responsável por uma segunda empresa: o convite não
	// pede senha, porque o login já existe.
	c := anonimo(t)
	c.post("/v1/empresas", map[string]any{
		"cnpj": cnpjAleatorio(), "razao_social": "Segunda Empresa Ltda", "plano": "profissional",
		"responsavel": map[string]any{"nome": "Ana", "email": em},
		"unidade":     map[string]any{"codigo": "FIL", "nome": "Filial", "turnos": turnoUtil},
	}).ok(http.StatusCreated)
	tok := caixa.token(t, em)
	conv := c.post("/v1/auth/convites/consultar", map[string]any{"token": tok}).ok(http.StatusOK)
	if conv["pede_senha"] != false || conv["email"] != em {
		t.Fatalf("convite de quem já tem login: %v", conv)
	}
	if a := c.post("/v1/auth/convites/aceitar", map[string]any{"token": tok}).ok(http.StatusOK); a["identidade_nova"] != false {
		t.Fatalf("aceite: %v", a)
	}
	// O link é de uso único.
	c.post("/v1/auth/convites/aceitar", map[string]any{"token": tok}).problema(http.StatusGone, "link_email_invalido")

	web := entrar(t, em, "")
	sessao := web.get("/v1/auth/sessao").ok(http.StatusOK)
	empresas := sessao["empresas"].([]any)
	if len(empresas) != 2 || sessao["escopo_acesso"] != "interno" || sessao["canal"] != "web" || texto(sessao, "csrf_token") != web.csrf {
		t.Fatalf("sessão: %v", sessao)
	}
	if _, tem := sessao["empresa_id"]; tem {
		t.Fatal("empresa_id não sai na sessão (seção 2.2)")
	}

	// Ativa o MFA.
	cfg := web.post("/v1/auth/mfa/iniciar", nil).ok(http.StatusOK)
	segredo := texto(cfg, "segredo_base32")
	if !strings.HasPrefix(texto(cfg, "uri_otpauth"), "otpauth://totp/") {
		t.Fatalf("configuração de MFA: %v", cfg)
	}
	web.post("/v1/auth/mfa/confirmar", map[string]any{"codigo": "000000"}).problema(http.StatusUnauthorized, "mfa_codigo_invalido")
	web.post("/v1/auth/mfa/confirmar", map[string]any{"codigo": codigoTOTP(t, segredo, time.Now())}).ok(http.StatusNoContent)
	web.post("/v1/auth/logout", nil).ok(http.StatusNoContent)
	web.get("/v1/auth/sessao").problema(http.StatusUnauthorized, "credencial_invalida")

	// Login com MFA: 401 mfa_obrigatorio com desafio_id.
	p := anonimo(t).post("/v1/auth/login", map[string]any{"email": em, "senha": senhaBoa}).problema(http.StatusUnauthorized, "mfa_obrigatorio")
	desafio := texto(p, "desafio_id")
	if len(desafio) != 43 || p["desafio_expira_em"] == nil {
		t.Fatalf("desafio: %v", p)
	}
	e := anonimo(t).post("/v1/auth/mfa/verificar", map[string]any{"desafio_id": desafio, "codigo": "000000"}).problema(http.StatusUnauthorized, "mfa_codigo_invalido")
	if e["tentativas_restantes"] != float64(4) {
		t.Fatalf("tentativas_restantes: %v", e)
	}
	// O código do passo usado na confirmação não vale de novo (RFC 6238,
	// seção 5.2): o login usa o do passo seguinte, ainda dentro da janela.
	r := anonimo(t).post("/v1/auth/mfa/verificar", map[string]any{"desafio_id": desafio, "codigo": codigoTOTP(t, segredo, time.Now().Add(30*time.Second))})
	s := r.ok(http.StatusOK)
	web = &cliente{t: t, origem: origemWeb, cookie: cookieDaResposta(t, r), csrf: texto(s, "csrf_token")}
	// Desafio gasto não vale de novo.
	anonimo(t).post("/v1/auth/mfa/verificar", map[string]any{"desafio_id": desafio, "codigo": codigoTOTP(t, segredo, time.Now())}).
		problema(http.StatusUnauthorized, "credencial_invalida")

	// Troca de empresa: credencial nova, a antiga cai.
	var outra string
	for _, x := range s["empresas"].([]any) {
		if m := x.(map[string]any); m["atual"] == false {
			outra = texto(m, "id")
		}
	}
	antigo := web.cookie
	web.post("/v1/auth/trocar-empresa", map[string]any{"empresa_id": outra}).problema(http.StatusBadRequest, "campo_nao_permitido")
	web.post("/v1/auth/trocar-empresa", map[string]any{"empresa_escolhida_id": "0193a1b2-7c4d-7e8f-9a0b-1c2d3e4f5a6b"}).
		problema(http.StatusNotFound, "recurso_nao_encontrado")
	tr := web.post("/v1/auth/trocar-empresa", map[string]any{"empresa_escolhida_id": outra})
	ts := tr.ok(http.StatusOK)
	novo := &cliente{t: t, origem: origemWeb, cookie: cookieDaResposta(t, tr), csrf: texto(ts, "csrf_token")}
	if ts["tokens"] != nil || empresaAtual(t, novo) != outra {
		t.Fatalf("troca de empresa: %v", ts)
	}
	(&cliente{t: t, cookie: antigo}).get("/v1/auth/sessao").problema(http.StatusUnauthorized, "credencial_invalida")
	if e := novo.get("/v1/empresa").ok(http.StatusOK); texto(e, "id") != outra || e["razao_social"] != "Segunda Empresa Ltda" {
		t.Fatalf("empresa após a troca: %v", e)
	}

	// Logout: cookie apagado e sessão revogada no servidor.
	lo := novo.post("/v1/auth/logout", nil)
	lo.ok(http.StatusNoContent)
	if sc := lo.cab.Get("Set-Cookie"); !strings.Contains(sc, "__Host-tenaz_sessao=;") || !strings.Contains(sc, "Max-Age=0") {
		t.Fatalf("logout sem apagar o cookie: %q", sc)
	}
	novo.get("/v1/auth/sessao").problema(http.StatusUnauthorized, "credencial_invalida")
}

func TestRecuperarERedefinirSenha(t *testing.T) {
	em := novaConta(t, "essencial", "")
	c := anonimo(t)
	// Mesma resposta para e-mail inexistente.
	c.post("/v1/auth/recuperar-senha", map[string]any{"email": emailAleatorio("ninguem")}).ok(http.StatusAccepted)
	c.post("/v1/auth/recuperar-senha", map[string]any{"email": em}).ok(http.StatusAccepted)
	tok := caixa.token(t, em)
	c.post("/v1/auth/redefinir-senha", map[string]any{"token": tok, "nova_senha": "curta"}).problema(http.StatusUnprocessableEntity, "senha_fraca")
	nova := "Outra senha longa de teste 42"
	c.post("/v1/auth/redefinir-senha", map[string]any{"token": tok, "nova_senha": nova}).ok(http.StatusNoContent)
	c.post("/v1/auth/redefinir-senha", map[string]any{"token": tok, "nova_senha": nova}).problema(http.StatusGone, "link_email_invalido")
	c.post("/v1/auth/login", map[string]any{"email": em, "senha": senhaBoa}).problema(http.StatusUnauthorized, "credencial_invalida")
	r := c.post("/v1/auth/login", map[string]any{"email": em, "senha": nova})
	s := r.ok(http.StatusOK)
	web := &cliente{t: t, origem: origemWeb, cookie: cookieDaResposta(t, r), csrf: texto(s, "csrf_token")}
	web.post("/v1/auth/trocar-senha", map[string]any{"senha_atual": "errada errada errada", "nova_senha": senhaBoa}).
		problema(http.StatusUnauthorized, "credencial_invalida")
	web.post("/v1/auth/trocar-senha", map[string]any{"senha_atual": nova, "nova_senha": senhaBoa}).ok(http.StatusNoContent)
	web.get("/v1/auth/sessao").ok(http.StatusOK) // a sessão atual continua
}

// OAuth 2.1 com PKCE do app, com renovação e detecção de reuso.
func TestOAuthPKCEDoApp(t *testing.T) {
	em := novaConta(t, "essencial", "")
	web := entrar(t, em, "")
	empresa := empresaAtual(t, web)
	verificador, desafio := pkce()
	q := url.Values{
		"response_type": {"code"}, "client_id": {"tenaz-app"}, "redirect_uri": {"br.com.tenaz.app:/oauth"},
		"code_challenge": {desafio}, "code_challenge_method": {"S256"}, "state": {"xyz"},
	}

	// GET /oauth/authorize redireciona para a tela da web, com os parâmetros.
	r := anonimo(t).get("/oauth/authorize?" + q.Encode())
	if r.status != http.StatusFound || !strings.HasPrefix(r.cab.Get("Location"), origemWeb+"/autorizar-app?") ||
		!strings.Contains(r.cab.Get("Location"), "code_challenge="+desafio) {
		t.Fatalf("authorize: %d %q", r.status, r.cab.Get("Location"))
	}
	// URI de retorno não registrada nunca redireciona.
	ruim := url.Values{}
	for k, v := range q {
		ruim[k] = v
	}
	ruim.Set("redirect_uri", "https://atacante.exemplo/oauth")
	anonimo(t).get("/oauth/authorize?"+ruim.Encode()).problema(http.StatusUnprocessableEntity, "validacao_falhou")

	corpo := map[string]any{
		"response_type": "code", "client_id": "tenaz-app", "redirect_uri": "br.com.tenaz.app:/oauth",
		"code_challenge": desafio, "code_challenge_method": "S256", "state": "xyz", "empresa_escolhida_id": empresa,
	}
	// Só a sessão web autoriza; sem CSRF, 403.
	(&cliente{t: t, cookie: web.cookie, origem: origemWeb}).post("/oauth/authorize", corpo).problema(http.StatusForbidden, "csrf_invalido")
	anonimo(t).post("/oauth/authorize", corpo).problema(http.StatusUnauthorized, "nao_autenticado")
	a := web.post("/oauth/authorize", corpo).ok(http.StatusOK)
	retorno, err := url.Parse(texto(a, "uri_redirecionamento"))
	if err != nil || retorno.Query().Get("state") != "xyz" || retorno.Query().Get("code") == "" {
		t.Fatalf("redirecionamento: %v", a)
	}
	codigo := retorno.Query().Get("code")

	app := anonimo(t)
	app.origem = ""
	troca := url.Values{"grant_type": {"authorization_code"}, "code": {codigo}, "redirect_uri": {"br.com.tenaz.app:/oauth"},
		"client_id": {"tenaz-app"}, "code_verifier": {"verificador-errado-" + verificador[:20]}}
	app.post("/oauth/token", troca).problema(http.StatusUnauthorized, "credencial_invalida")

	// Código gasto pela tentativa errada: precisa de outro.
	a = web.post("/oauth/authorize", corpo).ok(http.StatusOK)
	retorno, _ = url.Parse(texto(a, "uri_redirecionamento"))
	troca.Set("code", retorno.Query().Get("code"))
	troca.Set("code_verifier", verificador)
	tr := app.post("/oauth/token", troca)
	tk := tr.ok(http.StatusOK)
	if tk["token_type"] != "Bearer" || tk["expires_in"] != float64(900) || !strings.HasPrefix(texto(tk, "refresh_token"), "tnz_rt_") ||
		tr.cab.Get("Cache-Control") != "no-store" {
		t.Fatalf("tokens: %v", tk)
	}
	app.bearer = texto(tk, "access_token")
	if eu := app.get("/v1/usuarios/eu").ok(http.StatusOK); eu["email"] != em {
		t.Fatalf("eu pelo app: %v", eu)
	}
	if s := app.get("/v1/auth/sessao").ok(http.StatusOK); s["canal"] != "app" || s["csrf_token"] != nil {
		t.Fatalf("sessão do app: %v", s)
	}
	// O app não precisa de CSRF.
	app.post("/v1/auth/mfa/iniciar", nil).ok(http.StatusOK)

	// Renovação rotativa.
	r1 := texto(tk, "refresh_token")
	ren := url.Values{"grant_type": {"refresh_token"}, "refresh_token": {r1}, "client_id": {"tenaz-app"}}
	tk2 := app.post("/oauth/token", ren).ok(http.StatusOK)
	r2 := texto(tk2, "refresh_token")
	if r2 == r1 {
		t.Fatal("o refresh token não foi trocado")
	}
	// Reuso do r1: a família inteira cai, inclusive o r2 e o access token.
	p := app.post("/oauth/token", ren).problema(http.StatusUnauthorized, "credencial_invalida")
	if p["error"] != "invalid_grant" {
		t.Fatalf("reuso sem error=invalid_grant: %v", p)
	}
	ren.Set("refresh_token", r2)
	app.post("/oauth/token", ren).problema(http.StatusUnauthorized, "credencial_invalida")
	app.bearer = texto(tk2, "access_token")
	app.get("/v1/usuarios/eu").problema(http.StatusUnauthorized, "credencial_invalida")

	// Revogação (sair no app) e token desconhecido: 200.
	a = web.post("/oauth/authorize", corpo).ok(http.StatusOK)
	retorno, _ = url.Parse(texto(a, "uri_redirecionamento"))
	troca.Set("code", retorno.Query().Get("code"))
	tk3 := app.post("/oauth/token", troca).ok(http.StatusOK)
	app.post("/oauth/revoke", url.Values{"token": {texto(tk3, "refresh_token")}}).ok(http.StatusOK)
	app.post("/oauth/revoke", url.Values{"token": {"tnz_rt_desconhecido"}}).ok(http.StatusOK)
	app.bearer = texto(tk3, "access_token")
	app.get("/v1/usuarios/eu").problema(http.StatusUnauthorized, "credencial_invalida")

	// O locatario não usa o app: a autorização volta com access_denied.
	_, _, emLoc := convidar(t, web, []string{"locatario"}, map[string]any{"locatario_id": "0193a1b2-7c4d-7e8f-9a0b-1c2d3e4f5a6b"})
	loc := entrar(t, emLoc, empresa)
	a = loc.post("/oauth/authorize", corpo).ok(http.StatusOK)
	retorno, _ = url.Parse(texto(a, "uri_redirecionamento"))
	if retorno.Query().Get("error") != "access_denied" || retorno.Query().Get("code") != "" {
		t.Fatalf("locatario no app: %v", a)
	}
}
