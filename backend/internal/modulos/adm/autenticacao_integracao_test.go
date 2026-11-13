package adm_test

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"tenaz/internal/modulos/adm"
	"tenaz/internal/plataforma/autenticacao"
	"tenaz/internal/plataforma/erros"
)

// Testes de integração da autenticação (RF-ADM-004, RF-ADM-006, RNF-017;
// documento 9, seção 5.3), pedidos pelo tenaz-qa-automacao na etapa E01.
// Os prazos que o banco confere (expira_em > now()) são vencidos movendo
// criado_em e expira_em para o passado pelo migrador, dentro dos CHECKs.

func hashDoSegredo(texto string) []byte {
	h := sha256.Sum256([]byte(texto))
	return h[:]
}

// relogioDoBanco é o clock_timestamp() do banco: instantes gravados pelo
// banco são comparados com o relógio dele, e não com o da máquina do teste
// (o relógio da VM do banco pode se afastar do da máquina).
func relogioDoBanco(t *testing.T) time.Time {
	t.Helper()
	var agora time.Time
	if err := amb.Migrador.QueryRow(context.Background(), `SELECT clock_timestamp()`).Scan(&agora); err != nil {
		t.Fatal(err)
	}
	return agora
}

func executar(t *testing.T, sql string, args ...any) {
	t.Helper()
	tag, err := amb.Migrador.Exec(context.Background(), sql, args...)
	if err != nil {
		t.Fatalf("%s: %v", sql, err)
	}
	if tag.RowsAffected() != 1 {
		t.Fatalf("%s: %d linhas, esperada 1", sql, tag.RowsAffected())
	}
}

func comoErro(t *testing.T, err error) *erros.Erro {
	t.Helper()
	var e *erros.Erro
	if !errors.As(err, &e) {
		t.Fatalf("esperado *erros.Erro, veio %v", err)
	}
	return e
}

// CSRF: com o cookie válido, POST/PATCH/DELETE (exigirCSRF) sem o token, com
// o token de outra sessão válida, com lixo ou com o próprio cookie é
// csrf_invalido; GET (sem exigência) passa sem o token.
func TestCSRFComCookieValido(t *testing.T) {
	c := novaConta(t, "")
	ctx := ctxWeb()
	r1, err := auth.Login(ctx, autenticacao.EntradaLogin{Email: c.email, Senha: senhaBoa})
	if err != nil {
		t.Fatal(err)
	}
	r2, err := auth.Login(ctx, autenticacao.EntradaLogin{Email: c.email, Senha: senhaBoa})
	if err != nil {
		t.Fatal(err)
	}
	s1, s2 := r1.Sessao, r2.Sessao
	if s1.CSRF == s2.CSRF {
		t.Fatal("cada sessão tem o seu CSRF")
	}
	lixo := base64.RawURLEncoding.EncodeToString(make([]byte, 32))
	for nome, csrf := range map[string]string{
		"ausente": "", "de outra sessão válida": s2.CSRF, "lixo bem formado": lixo, "o próprio cookie": s1.Cookie,
		"com espaço": " " + s1.CSRF, "truncado": s1.CSRF[:len(s1.CSRF)-1],
	} {
		if _, err := auth.AutenticarWeb(ctx, s1.Cookie, csrf, true); !errors.Is(err, erros.CSRFInvalido) {
			t.Errorf("CSRF %s: esperado csrf_invalido, veio %v", nome, err)
		}
	}
	if _, err := auth.AutenticarWeb(ctx, s1.Cookie, s1.CSRF, true); err != nil {
		t.Fatalf("CSRF certo: %v", err)
	}
	if _, err := auth.AutenticarWeb(ctx, s1.Cookie, "", false); err != nil {
		t.Fatalf("leitura sem CSRF: %v", err)
	}
	// Cookie inválido com CSRF certo: a credencial cai antes do CSRF.
	_, err = auth.AutenticarWeb(ctx, lixo, s1.CSRF, true)
	exigirErro(t, err, erros.CredencialInvalida)
}

// Convite (7 dias, uso único, o reenvio revoga o anterior): expirado,
// reapresentado e revogado respondem link_email_invalido, e o usuário
// continua convidado.
func TestConviteExpiradoReapresentadoERevogado(t *testing.T) {
	c := novaConta(t, "")
	ctx := ctxWeb()
	novo := func(t *testing.T) (adm.Usuario, string) {
		t.Helper()
		em := emailAleatorio("convite")
		u, err := svc.CriarUsuario(ctx, c.adm, adm.EntradaUsuario{Nome: "Convidado", Email: em, Perfis: []string{"tecnico"}})
		if err != nil {
			t.Fatal(err)
		}
		return u.Usuario, caixa.token(t, em)
	}
	situacao := func(t *testing.T, id adm.Usuario) string {
		t.Helper()
		u, err := svc.ObterUsuario(ctx, c.adm, id.ID)
		if err != nil {
			t.Fatal(err)
		}
		return u.Situacao
	}

	t.Run("expirado", func(t *testing.T) {
		u, tok := novo(t)
		executar(t, `UPDATE plataforma.token_email SET criado_em = now() - interval '8 days', expira_em = now() - interval '1 day'
WHERE token_hash = $1`, hashDoSegredo(tok))
		_, err := auth.LocalizarConvite(ctx, tok)
		exigirErro(t, err, erros.LinkEmailInvalido)
		_, err = svc.AceitarConvite(ctx, tok, senhaBoa, nil)
		exigirErro(t, err, erros.LinkEmailInvalido)
		if s := situacao(t, u); s != "convidado" {
			t.Errorf("convite vencido não ativa: %s", s)
		}
	})

	t.Run("no último segundo ainda vale", func(t *testing.T) {
		u, tok := novo(t)
		executar(t, `UPDATE plataforma.token_email SET criado_em = now() - interval '7 days' + interval '1 minute',
  expira_em = now() + interval '1 minute' WHERE token_hash = $1`, hashDoSegredo(tok))
		if _, err := svc.AceitarConvite(ctx, tok, senhaBoa, nil); err != nil {
			t.Fatalf("convite a um minuto do fim: %v", err)
		}
		if s := situacao(t, u); s != "ativo" {
			t.Errorf("situação: %s", s)
		}
	})

	t.Run("reapresentado", func(t *testing.T) {
		u, tok := novo(t)
		if _, err := svc.AceitarConvite(ctx, tok, senhaBoa, nil); err != nil {
			t.Fatal(err)
		}
		_, err := svc.AceitarConvite(ctx, tok, senhaBoa, nil)
		exigirErro(t, err, erros.LinkEmailInvalido)
		_, err = svc.AceitarConvite(ctx, tok, "", nil)
		exigirErro(t, err, erros.LinkEmailInvalido)
		if s := situacao(t, u); s != "ativo" {
			t.Errorf("situação: %s", s)
		}
	})

	t.Run("revogado no reenvio", func(t *testing.T) {
		u, antigo := novo(t)
		if err := svc.ReenviarConvite(ctx, c.adm, u.ID); err != nil {
			t.Fatal(err)
		}
		atual := caixa.token(t, *u.Email)
		if atual == antigo {
			t.Fatal("o reenvio gera link novo")
		}
		_, err := svc.AceitarConvite(ctx, antigo, senhaBoa, nil)
		exigirErro(t, err, erros.LinkEmailInvalido)
		if s := situacao(t, u); s != "convidado" {
			t.Fatalf("link revogado não ativa: %s", s)
		}
		if _, err := svc.AceitarConvite(ctx, atual, senhaBoa, nil); err != nil {
			t.Fatalf("o link do reenvio vale: %v", err)
		}
	})

	t.Run("malformado", func(t *testing.T) {
		for _, tok := range []string{"", "curto", strings.Repeat("*", 43)} {
			_, err := svc.AceitarConvite(ctx, tok, senhaBoa, nil)
			exigirErro(t, err, erros.LinkEmailInvalido)
		}
	})
}

// Redefinição de senha (1 hora, uso único, novo pedido revoga o anterior):
// link vencido e link substituído respondem link_email_invalido, e a senha
// antiga continua valendo.
func TestLinkDeRedefinicaoExpiradoOuSubstituido(t *testing.T) {
	c := novaConta(t, "")
	ctx := ctxWeb()
	nova := "Outra senha boa 2026!"

	if err := auth.RecuperarSenha(ctx, c.email); err != nil {
		t.Fatal(err)
	}
	vencido := caixa.token(t, c.email)
	executar(t, `UPDATE plataforma.token_email SET criado_em = now() - interval '2 hours', expira_em = now() - interval '1 hour'
WHERE token_hash = $1`, hashDoSegredo(vencido))
	exigirErro(t, auth.RedefinirSenha(ctx, vencido, nova), erros.LinkEmailInvalido)

	if err := auth.RecuperarSenha(ctx, c.email); err != nil {
		t.Fatal(err)
	}
	primeiro := caixa.token(t, c.email)
	if err := auth.RecuperarSenha(ctx, c.email); err != nil {
		t.Fatal(err)
	}
	segundo := caixa.token(t, c.email)
	if primeiro == segundo {
		t.Fatal("cada pedido gera link novo")
	}
	exigirErro(t, auth.RedefinirSenha(ctx, primeiro, nova), erros.LinkEmailInvalido)
	entrar(t, c.email, senhaBoa, nil) // nada mudou
	if err := auth.RedefinirSenha(ctx, segundo, nova); err != nil {
		t.Fatalf("o link mais novo vale: %v", err)
	}
	entrar(t, c.email, nova, nil)
}

// falhasDe lê falhas_login e bloqueada_ate da identidade, e o instante do banco.
func falhasDe(t *testing.T, email string) (int, *time.Time, time.Time) {
	t.Helper()
	var (
		n     int
		ate   *time.Time
		agora time.Time
	)
	if err := amb.Migrador.QueryRow(context.Background(),
		`SELECT falhas_login, bloqueada_ate, now() FROM plataforma.identidade WHERE email = $1`, email).Scan(&n, &ate, &agora); err != nil {
		t.Fatal(err)
	}
	return n, ate, agora
}

// Bloqueio progressivo no fluxo real (RF-ADM-006; documento 9, seção 5.3):
// 4 falhas ainda entram e o login certo zera a contagem; a 5ª bloqueia por
// 15 minutos; tentativa durante o bloqueio não soma, nem com a senha certa;
// a 10ª bloqueia por 30 minutos.
func TestFronteiraDoBloqueioProgressivo(t *testing.T) {
	c := novaConta(t, "")
	ctx := ctxWeb()
	errar := func(t *testing.T, n int) {
		t.Helper()
		for i := 0; i < n; i++ {
			_, err := auth.Login(ctx, autenticacao.EntradaLogin{Email: c.email, Senha: "errada-errada-1"})
			exigirErro(t, err, erros.CredencialInvalida)
		}
	}
	conferir := func(t *testing.T, falhas int, bloqueio time.Duration) {
		t.Helper()
		n, ate, agora := falhasDe(t, c.email)
		if n != falhas {
			t.Fatalf("falhas_login = %d, esperado %d", n, falhas)
		}
		switch {
		case bloqueio == 0 && ate != nil && ate.After(agora):
			t.Fatalf("bloqueada até %v sem motivo", ate)
		case bloqueio > 0 && (ate == nil || ate.Sub(agora) < bloqueio-time.Minute || ate.Sub(agora) > bloqueio):
			t.Fatalf("bloqueio esperado de %v; bloqueada_ate %v, agora %v", bloqueio, ate, agora)
		}
	}

	errar(t, 4)
	conferir(t, 4, 0)
	entrar(t, c.email, senhaBoa, nil) // 4 falhas ainda entram
	conferir(t, 0, 0)                 // e o login certo zera
	errar(t, 4)
	entrar(t, c.email, senhaBoa, nil) // sem zerar, seriam 8: bloqueada na 5ª
	conferir(t, 0, 0)

	errar(t, 5)
	conferir(t, 5, 15*time.Minute)
	_, ateAntes, _ := falhasDe(t, c.email)
	errar(t, 3) // durante o bloqueio
	_, err := auth.Login(ctx, autenticacao.EntradaLogin{Email: c.email, Senha: senhaBoa})
	exigirErro(t, err, erros.CredencialInvalida)
	conferir(t, 5, 15*time.Minute)
	if _, ateDepois, _ := falhasDe(t, c.email); !ateDepois.Equal(*ateAntes) {
		t.Errorf("tentativa durante o bloqueio estendeu o prazo: %v -> %v", ateAntes, ateDepois)
	}

	// Fim do bloqueio (simulado): as falhas seguem contando; a 10ª dobra.
	executar(t, `UPDATE plataforma.identidade SET bloqueada_ate = now() - interval '1 second' WHERE email = $1`, c.email)
	errar(t, 4)
	conferir(t, 9, 0)
	errar(t, 1)
	conferir(t, 10, 30*time.Minute)
}

// Mesma resposta para senha errada, conta bloqueada (com a senha certa ou
// errada) e e-mail inexistente: o mesmo erro, com o mesmo texto, os mesmos
// campos e os mesmos dados (documento 9, seção 5.3).
func TestMesmaRespostaParaBloqueioESenhaErrada(t *testing.T) {
	c := novaConta(t, "")
	ctx := ctxWeb()
	login := func(email, senha string) *erros.Erro {
		_, err := auth.Login(ctx, autenticacao.EntradaLogin{Email: email, Senha: senha})
		return comoErro(t, err)
	}
	senhaErrada := login(c.email, "errada-errada-1")
	inexistente := login(emailAleatorio("ninguem"), senhaBoa)
	for i := 0; i < 4; i++ {
		login(c.email, "errada-errada-1")
	}
	bloqueadaCerta := login(c.email, senhaBoa)
	bloqueadaErrada := login(c.email, "errada-errada-1")
	senhaVazia := login(c.email, "")
	for nome, e := range map[string]*erros.Erro{
		"e-mail inexistente": inexistente, "bloqueada com a senha certa": bloqueadaCerta,
		"bloqueada com a senha errada": bloqueadaErrada, "senha vazia": senhaVazia,
	} {
		if e.Codigo != senhaErrada.Codigo || e.Status != senhaErrada.Status || e.Titulo != senhaErrada.Titulo ||
			e.Detalhe != senhaErrada.Detalhe || e.Error() != senhaErrada.Error() ||
			!reflect.DeepEqual(e.Campos, senhaErrada.Campos) || !reflect.DeepEqual(e.Dados, senhaErrada.Dados) {
			t.Errorf("%s responde diferente da senha errada:\n %+v\n %+v", nome, e, senhaErrada)
		}
	}
	if senhaErrada.Codigo != erros.CredencialInvalida.Codigo {
		t.Errorf("código: %s", senhaErrada.Codigo)
	}
}

// Sessão web (RNF-017; documento 3, seção 4.2): 12 horas de inatividade e
// teto de 7 dias desde a criação, conferidos a cada requisição.
func TestSessaoWebExpirada(t *testing.T) {
	c := novaConta(t, "")
	ctx := ctxWeb()
	nova := func(t *testing.T) *autenticacao.SessaoWeb {
		t.Helper()
		r, err := auth.Login(ctx, autenticacao.EntradaLogin{Email: c.email, Senha: senhaBoa})
		if err != nil {
			t.Fatal(err)
		}
		return r.Sessao
	}

	t.Run("inatividade", func(t *testing.T) {
		s := nova(t)
		executar(t, `UPDATE plataforma.sessao SET ultimo_uso_em = now() - interval '11 hours 50 minutes' WHERE token_hash = $1`,
			hashDoSegredo(s.Cookie))
		if _, err := auth.AutenticarWeb(ctx, s.Cookie, s.CSRF, true); err != nil {
			t.Fatalf("antes das 12 horas a sessão vale: %v", err)
		}
		// O uso renova a inatividade.
		var ultimo time.Time
		if err := amb.Migrador.QueryRow(context.Background(),
			`SELECT ultimo_uso_em FROM plataforma.sessao WHERE token_hash = $1`, hashDoSegredo(s.Cookie)).Scan(&ultimo); err != nil {
			t.Fatal(err)
		}
		if relogioDoBanco(t).Sub(ultimo) > time.Minute {
			t.Errorf("o uso não renovou ultimo_uso_em: %v", ultimo)
		}
		executar(t, `UPDATE plataforma.sessao SET ultimo_uso_em = now() - interval '12 hours 1 minute' WHERE token_hash = $1`,
			hashDoSegredo(s.Cookie))
		_, err := auth.AutenticarWeb(ctx, s.Cookie, s.CSRF, true)
		exigirErro(t, err, erros.CredencialInvalida)
		_, err = auth.ObterSessaoWeb(ctx, s.Cookie)
		exigirErro(t, err, erros.CredencialInvalida)
	})

	t.Run("máximo de 7 dias", func(t *testing.T) {
		s := nova(t)
		executar(t, `UPDATE plataforma.sessao SET criado_em = now() - interval '7 days 1 minute',
  expira_em = now() - interval '1 minute', ultimo_uso_em = now() WHERE token_hash = $1`, hashDoSegredo(s.Cookie))
		_, err := auth.AutenticarWeb(ctx, s.Cookie, s.CSRF, true)
		exigirErro(t, err, erros.CredencialInvalida)
	})

	t.Run("o teto é conferido pelo banco", func(t *testing.T) {
		s := nova(t)
		_, err := amb.Migrador.Exec(context.Background(),
			`UPDATE plataforma.sessao SET expira_em = criado_em + interval '7 days 1 second' WHERE token_hash = $1`, hashDoSegredo(s.Cookie))
		if err == nil || !strings.Contains(err.Error(), "ck_sessao_validade") {
			t.Fatalf("sessão web acima de 7 dias: esperado ck_sessao_validade, veio %v", err)
		}
	})
}

// Refresh token vencido (30 dias de inatividade, RNF-017) não renova; o
// vencimento não é tratado como reuso.
func TestRefreshTokenVencido(t *testing.T) {
	c := novaConta(t, "")
	ctx := ctxWeb()
	verificador := strings.Repeat("v3rif1cad0r-", 5)
	h := sha256.Sum256([]byte(verificador))
	entrada := autenticacao.EntradaAutorizacao{ClienteID: "tenaz-app", URIRetorno: "br.com.tenaz.app:/oauth",
		DesafioPKCE: base64.RawURLEncoding.EncodeToString(h[:]), MetodoPKCE: "S256", EmpresaID: c.criada.EmpresaID}
	codigo, err := auth.EmitirCodigoAutorizacao(ctx, c.adm, entrada)
	if err != nil {
		t.Fatal(err)
	}
	tok, err := auth.TrocarCodigo(ctx, autenticacao.EntradaTrocaCodigo{Codigo: codigo, ClienteID: "tenaz-app",
		URIRetorno: entrada.URIRetorno, VerificadorPKCE: verificador})
	if err != nil {
		t.Fatal(err)
	}
	texto := strings.TrimPrefix(tok.RefreshToken, autenticacao.PrefixoRefreshToken)
	executar(t, `UPDATE plataforma.token_renovacao SET criado_em = now() - interval '31 days', expira_em = now() - interval '1 day'
WHERE token_hash = $1`, hashDoSegredo(texto))
	for i := 0; i < 2; i++ {
		_, err = auth.Renovar(ctx, tok.RefreshToken, "tenaz-app")
		exigirErro(t, err, erros.CredencialInvalida)
	}
	var revogado *time.Time
	if err := amb.Migrador.QueryRow(context.Background(),
		`SELECT revogado_em FROM plataforma.token_renovacao WHERE token_hash = $1`, hashDoSegredo(texto)).Scan(&revogado); err != nil {
		t.Fatal(err)
	}
	if revogado != nil {
		t.Error("token vencido não é reuso: a família não é revogada por isso")
	}
}

// Código TOTP reapresentado (RFC 6238, seção 5.2; mfa_ultimo_passo): o código
// aceito num login não vale no seguinte, nem o de um passo anterior, mesmo
// dentro da janela de ±1 passo; vale também para desativar o MFA.
func TestCodigoTOTPReapresentado(t *testing.T) {
	c := novaConta(t, "")
	ctx := ctxWeb()
	cfg, err := auth.IniciarMFA(ctx, c.adm)
	if err != nil {
		t.Fatal(err)
	}
	agora := time.Now()
	if err := auth.ConfirmarMFA(ctx, c.adm, codigoTOTP(t, cfg.SegredoBase32, agora)); err != nil {
		t.Fatal(err)
	}
	desafio := func(t *testing.T) string {
		t.Helper()
		r, err := auth.Login(ctx, autenticacao.EntradaLogin{Email: c.email, Senha: senhaBoa})
		if err != nil || r.DesafioMFA == nil {
			t.Fatalf("login com MFA: %+v %v", r, err)
		}
		return r.DesafioMFA.DesafioID
	}
	seguinte := codigoTOTP(t, cfg.SegredoBase32, agora.Add(30*time.Second))
	r, err := auth.VerificarMFA(ctx, autenticacao.EntradaVerificarMFA{DesafioID: desafio(t), Codigo: seguinte})
	if err != nil || r.Sessao == nil {
		t.Fatalf("passo seguinte, dentro da janela: %v", err)
	}

	d := desafio(t)
	for nome, cod := range map[string]string{
		"o mesmo código":      seguinte,
		"o passo da ativação": codigoTOTP(t, cfg.SegredoBase32, agora),
		"um passo antes":      codigoTOTP(t, cfg.SegredoBase32, agora.Add(-30*time.Second)),
	} {
		_, err := auth.VerificarMFA(ctx, autenticacao.EntradaVerificarMFA{DesafioID: d, Codigo: cod})
		if !errors.Is(err, erros.MFACodigoInvalido) {
			t.Errorf("%s, reapresentado: esperado mfa_codigo_invalido, veio %v", nome, err)
		}
	}
	// Desativar também consome o passo: o código já usado não desliga o MFA.
	exigirErro(t, auth.DesativarMFA(ctx, c.adm, seguinte), erros.MFACodigoInvalido)
}
