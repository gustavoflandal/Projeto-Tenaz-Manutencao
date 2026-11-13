package vazamento

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
)

// Papel próprio das credenciais de autenticação (ADR-011, decisão 9;
// documento 5, seções 6.9 e 9.3). Cada regra tem um verificador que devolve a
// lista de violações e uma prova de detector: dentro de uma transação desfeita,
// a suíte estraga o banco de propósito (concede privilégio, derruba a policy,
// troca a função) e exige que o verificador acuse.

// As sete tabelas de credencial.
var tabelasCredencial = []struct{ nome, coluna string }{
	{"plataforma.identidade", "id"},
	{"plataforma.acesso_empresa", "identidade_id"},
	{"plataforma.sessao", "id"},
	{"plataforma.token_renovacao", "id"},
	{"plataforma.codigo_autorizacao", "id"},
	{"plataforma.desafio_mfa", "id"},
	{"plataforma.token_email", "id"},
}

// emTransacaoDesfeita abre uma transação do superusuário, sempre desfeita.
// Os verificadores rodam dentro dela, e as provas de detector a estragam.
func emTransacaoDesfeita(t *testing.T, ctx context.Context, fn func(tx pgx.Tx)) {
	t.Helper()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("abrir transação: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	fn(tx)
}

// trocarPapel volta ao superusuário, limpa o contexto e assume papel e
// contexto novos, na transação corrente.
func trocarPapel(ctx context.Context, tx pgx.Tx, papel string, c contexto) error {
	if _, err := tx.Exec(ctx, "RESET ROLE"); err != nil {
		return erroPreparacao{"RESET ROLE: " + err.Error()}
	}
	for _, v := range []string{"app.empresa_id", "app.usuario_id", "app.escopo", "app.fornecedor_id",
		"app.locatario_id", "app.link_hash", "app.conexao_empresa_id"} {
		if _, err := tx.Exec(ctx, `SELECT set_config($1, '', true)`, v); err != nil {
			return erroPreparacao{"limpar " + v + ": " + err.Error()}
		}
	}
	return prepararTransacao(ctx, tx, papel, c)
}

// tentarComo executa um comando como papel, num savepoint desfeito em
// seguida: o papel, o contexto e o efeito do comando somem, e a transação de
// fora segue utilizável mesmo depois de um erro.
func tentarComo(ctx context.Context, tx pgx.Tx, papel string, c contexto, sql string, args ...any) error {
	sub, err := tx.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = sub.Rollback(ctx) }()
	if err := trocarPapel(ctx, sub, papel, c); err != nil {
		// %v, e não %w: o 42501 de "permission denied to set role" não pode
		// contar como recusa do ataque (erroPreparacao não tem SQLSTATE).
		return erroPreparacao{fmt.Sprintf("preparar %s: %v", papel, err)}
	}
	_, err = sub.Exec(ctx, sql, args...)
	return err
}

// ------------------------------------------- ataques do ISO-E01-02

type tentativa struct {
	rotulo, sql string
	args        []any
}

// ataquesDoAchado são os ataques do ISO-E01-02, mais os quatro comandos em
// cada uma das sete tabelas.
func ataquesDoAchado(f *fixturesE01) []tentativa {
	a := []tentativa{
		{"lê a identidade de B", `SELECT email, senha_hash, mfa_segredo_cifrado FROM plataforma.identidade WHERE id = $1`,
			[]any{f.B.Identidade}},
		{"lê os acessos de B", `SELECT * FROM plataforma.acesso_empresa WHERE empresa_id = $1`, []any{f.B.ID}},
		{"lê as sessões de B", `SELECT ip, user_agent FROM plataforma.sessao WHERE empresa_id = $1`, []any{f.B.ID}},
		{"cria sessão para o usuário de B",
			`INSERT INTO plataforma.sessao (canal, token_hash, csrf_hash, identidade_id, empresa_id, usuario_id, expira_em)
			 VALUES ('web', $1, $2, $3, $4, $5, now() + interval '1 hour')`,
			[]any{hashAleatorio(), hashAleatorio(), f.B.Identidade, f.B.ID, f.B.UsuarioInterno}},
		{"cria código OAuth para o usuário de B",
			`INSERT INTO plataforma.codigo_autorizacao (codigo_hash, identidade_id, empresa_id, usuario_id,
			   cliente_oauth, uri_retorno, desafio_pkce, expira_em)
			 VALUES ($1, $2, $3, $4, 'tenaz-app', 'tenaz://retorno', 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQ',
			         now() + interval '60 seconds')`, []any{hashAleatorio(), f.B.Identidade, f.B.ID, f.B.UsuarioInterno}},
		{"troca a senha da identidade de B", `UPDATE plataforma.identidade SET senha_hash = 'trocada' WHERE id = $1`,
			[]any{f.B.Identidade}},
	}
	for _, tb := range tabelasCredencial {
		col := pgx.Identifier{tb.coluna}.Sanitize()
		a = append(a,
			tentativa{"SELECT em " + tb.nome, "SELECT count(*) FROM " + tb.nome, nil},
			tentativa{"INSERT em " + tb.nome, "INSERT INTO " + tb.nome + " DEFAULT VALUES", nil},
			tentativa{"UPDATE em " + tb.nome, "UPDATE " + tb.nome + " SET " + col + " = " + col + " WHERE false", nil},
			tentativa{"DELETE em " + tb.nome, "DELETE FROM " + tb.nome + " WHERE false", nil},
		)
	}
	return a
}

// ataquesAceitos devolve as tentativas que NÃO receberam "permissão negada".
// Zero linhas não basta: a decisão 1 do ADR-011 é que o papel não alcança a
// tabela, e não que a tabela lhe pareça vazia.
func ataquesAceitos(t *testing.T, ctx context.Context, tx pgx.Tx, papel string, f *fixturesE01) []string {
	t.Helper()
	var aceitos []string
	for _, a := range ataquesDoAchado(f) {
		err := tentarComo(ctx, tx, papel, f.internoA(), a.sql, a.args...)
		if sqlstate(err) != erroPrivilegio {
			aceitos = append(aceitos, fmt.Sprintf("%s (resultado: %v)", a.rotulo, err))
		}
	}
	return aceitos
}

// Todo papel, menos o da autenticação, recebe "permissão negada" nas
// credenciais — inclusive o dono do esquema e a manutenção, que é membro dele
// (ISO-E01-08: as tabelas são de tenaz_credencial).
var papeisSemCredencial = []string{papelApp, papelJob, papelLeitura, papelSinc, "tenaz_sinc_conexao",
	papelDono, papelManutencao}

func TestAtaquesDoISOE0102RecebemPermissaoNegada(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		for _, papel := range papeisSemCredencial {
			for _, a := range ataquesAceitos(t, ctx, tx, papel, f) {
				t.Errorf("%s: %s; esperado %s (ISO-E01-02, ADR-011)", papel, a, erroPrivilegio)
			}
		}
	})
}

func TestDetectorAcusaPrivilegioDeCredencialNoApp(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		estragar(t, ctx, tx, "GRANT SELECT ON plataforma.identidade TO tenaz_app")
		aceitos := strings.Join(ataquesAceitos(t, ctx, tx, papelApp, f), " | ")
		for _, esperado := range []string{"lê a identidade de B (", "SELECT em plataforma.identidade ("} {
			if !strings.Contains(aceitos, esperado) {
				t.Errorf("com SELECT em plataforma.identidade concedido ao tenaz_app, o verificador não acusou %q; "+
					"acusou: %s", esperado, aceitos)
			}
		}
		// E só isso: o INSERT continua recusado, e o verificador não acusa o que não mudou.
		if strings.Contains(aceitos, "INSERT em plataforma.identidade") {
			t.Errorf("o verificador acusou INSERT, que não foi concedido: %s", aceitos)
		}
	})
}

// ------------------------------------------------ atributos do papel

// problemasDoPapel confere o par papel/conexão (ADR-011, decisão 1): o papel é
// NOLOGIN, NOBYPASSRLS, sem atributo de administração, sem objeto, com USAGE
// exatamente em app, plataforma e public; a conexão é LOGIN e só herda o
// papel; nenhum dos dois tem relação de membro, em nenhum sentido, com os
// demais papéis do sistema; e só a conexão é membro do papel.
func problemasDoPapel(t *testing.T, ctx context.Context, tx pgx.Tx, papel, conexao string) []string {
	t.Helper()
	var problemas []string
	conferir := func(problema, sql string, args ...any) {
		var ok bool
		if err := tx.QueryRow(ctx, sql, args...).Scan(&ok); err != nil {
			t.Fatalf("conferir %q: %v", problema, err)
		}
		if !ok {
			problemas = append(problemas, problema)
		}
	}

	conferir("o papel não existe, faz login, contorna RLS ou tem atributo de administração", `
		SELECT coalesce((SELECT NOT rolcanlogin AND NOT rolbypassrls AND NOT rolsuper AND NOT rolcreaterole
		                        AND NOT rolcreatedb AND NOT rolreplication
		                 FROM pg_roles WHERE rolname = $1), false)`, papel)
	conferir("a conexão não existe, não faz login, contorna RLS ou tem atributo de administração", `
		SELECT coalesce((SELECT rolcanlogin AND NOT rolbypassrls AND NOT rolsuper AND NOT rolcreaterole
		                        AND NOT rolcreatedb AND NOT rolreplication
		                 FROM pg_roles WHERE rolname = $1), false)`, conexao)
	conferir("a conexão é membro de outro papel além do das credenciais, ou não herda dele", `
		SELECT coalesce(array_agg(r.rolname::text ORDER BY r.rolname) FILTER (WHERE m.inherit_option OR m.set_option), '{}')
		       = ARRAY[$1::text]
		       AND bool_and(m.inherit_option) FILTER (WHERE r.rolname = $1)
		FROM pg_auth_members m
		JOIN pg_roles r ON r.oid = m.roleid
		WHERE m.member = $2::regrole`, papel, conexao)
	conferir("outro papel além da conexão herda ou assume o papel das credenciais", `
		SELECT NOT EXISTS (SELECT 1 FROM pg_auth_members m
		                   WHERE m.roleid = $1::regrole AND m.member <> $2::regrole
		                     AND (m.inherit_option OR m.set_option))`, papel, conexao)
	conferir("há relação de membro com tenaz_dono, tenaz_app, tenaz_job ou outro papel do sistema", `
		SELECT NOT EXISTS (
		  SELECT 1
		  FROM unnest(ARRAY[$1, $2]) p,
		       unnest(ARRAY['tenaz_dono', 'tenaz_app', 'tenaz_job', 'tenaz_leitura', 'tenaz_sinc',
		                    'tenaz_sinc_conexao', 'tenaz_manutencao']) o
		  WHERE pg_has_role(p, o, 'MEMBER') OR pg_has_role(o, p, 'MEMBER'))`, papel, conexao)
	conferir("o papel ou a conexão é dono de algum objeto", `
		SELECT NOT EXISTS (SELECT 1 FROM pg_class WHERE relowner IN ($1::regrole, $2::regrole))
		   AND NOT EXISTS (SELECT 1 FROM pg_proc WHERE proowner IN ($1::regrole, $2::regrole))
		   AND NOT EXISTS (SELECT 1 FROM pg_namespace WHERE nspowner IN ($1::regrole, $2::regrole))
		   AND NOT EXISTS (SELECT 1 FROM pg_type WHERE typowner IN ($1::regrole, $2::regrole))`, papel, conexao)
	conferir("o papel não tem USAGE exatamente em app, plataforma e public", `
		SELECT coalesce(array_agg(nspname::text ORDER BY nspname), '{}') = ARRAY['app', 'plataforma', 'public']
		FROM pg_namespace
		WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg\_%'
		  AND has_schema_privilege($1, oid, 'USAGE')`, papel)
	return problemas
}

func TestPapelDeAutenticacaoTemOsAtributosDoADR011(t *testing.T) {
	ctx := context.Background()
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		for _, p := range problemasDoPapel(t, ctx, tx, papelAutenticacao, papelAutenticacaoConexao) {
			t.Errorf("%s/%s: %s (ADR-011, decisão 1)", papelAutenticacao, papelAutenticacaoConexao, p)
		}
	})
}

// A prova cria, numa transação desfeita (CREATE ROLE é transacional), um par
// com tudo errado: o papel faz login e contorna RLS, e a conexão é também
// membro do tenaz_app e não herda o papel.
func TestDetectorAcusaPapelDeAutenticacaoMalformado(t *testing.T) {
	ctx := context.Background()
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		for _, c := range []string{
			"CREATE ROLE tenaz_prova_aut LOGIN BYPASSRLS",
			"CREATE ROLE tenaz_prova_aut_conexao NOLOGIN",
			"GRANT tenaz_prova_aut TO tenaz_prova_aut_conexao WITH INHERIT FALSE",
			"GRANT tenaz_app TO tenaz_prova_aut_conexao",
		} {
			if _, err := tx.Exec(ctx, c); err != nil {
				t.Fatalf("preparar prova (%s): %v", c, err)
			}
		}
		p := problemasDoPapel(t, ctx, tx, "tenaz_prova_aut", "tenaz_prova_aut_conexao")
		for _, esperado := range []string{
			"o papel não existe, faz login, contorna RLS ou tem atributo de administração",
			"a conexão não existe, não faz login, contorna RLS ou tem atributo de administração",
			"a conexão é membro de outro papel além do das credenciais, ou não herda dele",
			"há relação de membro com tenaz_dono, tenaz_app, tenaz_job ou outro papel do sistema",
			"o papel não tem USAGE exatamente em app, plataforma e public",
		} {
			if !contem(p, esperado) {
				t.Errorf("o par malformado não foi acusado por %q; acusado: %v", esperado, p)
			}
		}
	})
}

// --------------------------------------------------- molde da auditoria

// eventoIdentidade monta o INSERT de evento que o serviço de autenticação faz.
func eventoIdentidade(empresa, usuario, identidade, ator, acao, entidade string) tentativa {
	return tentativa{
		fmt.Sprintf("empresa=%.8s usuario=%.8s identidade=%.8s ator=%s acao=%s entidade=%s",
			empresa, usuario, identidade, ator, acao, entidade),
		`INSERT INTO app.evento_auditoria (empresa_id, ator_tipo, usuario_id, acao, entidade, entidade_id, origem)
		 VALUES ($1, $2, $3, $4, $5, $6, 'web')`,
		[]any{empresa, ator, usuario, acao, entidade, identidade},
	}
}

// problemasDoMolde devolve o que o papel das credenciais conseguiu e não
// deveria (ou o contrário) na trilha de auditoria (ADR-011, decisão 5).
func problemasDoMolde(t *testing.T, ctx context.Context, tx pgx.Tx, papel string, f *fixturesE01) []string {
	t.Helper()
	ctxA := contexto{Empresa: f.A.ID}
	ctxB := contexto{Empresa: f.B.ID}
	a, b := f.A, f.B
	var problemas []string

	// Permitidos: login, falha, logout e alteração da identidade, com o usuário
	// que acesso_empresa liga àquela identidade naquela empresa — em A e em B.
	for _, ok := range []struct {
		c contexto
		e tentativa
	}{
		{ctxA, eventoIdentidade(a.ID, a.UsuarioInterno, a.Identidade, "usuario", "login", "identidade")},
		{ctxA, eventoIdentidade(a.ID, a.UsuarioInterno, a.Identidade, "usuario", "login_falho", "identidade")},
		{ctxA, eventoIdentidade(a.ID, a.UsuarioInterno, a.Identidade, "usuario", "logout", "identidade")},
		{ctxA, eventoIdentidade(a.ID, a.UsuarioInterno, a.Identidade, "usuario", "alteracao", "identidade")},
		{ctxB, eventoIdentidade(b.ID, b.UsuarioInterno, b.Identidade, "usuario", "login", "identidade")},
	} {
		if err := tentarComo(ctx, tx, papel, ok.c, ok.e.sql, ok.e.args...); err != nil {
			problemas = append(problemas, "evento legítimo recusado: "+ok.e.rotulo+": "+err.Error())
		}
	}

	// Recusados com "permissão negada".
	recusados := []struct {
		c contexto
		e tentativa
	}{
		{ctxA, eventoIdentidade(a.ID, a.UsuarioInterno, a.Identidade, "usuario", "login", "usuario")},
		{ctxA, eventoIdentidade(a.ID, a.UsuarioInterno, a.Identidade, "usuario", "exclusao", "identidade")},
		{ctxA, eventoIdentidade(a.ID, a.UsuarioInterno, a.Identidade, "usuario", "exportacao", "identidade")},
		{ctxA, eventoIdentidade(a.ID, a.UsuarioInterno, a.Identidade, "sistema", "login", "identidade")},
		// usuário de B, identidade de A, empresa A
		{ctxA, eventoIdentidade(a.ID, b.UsuarioInterno, a.Identidade, "usuario", "login", "identidade")},
		// usuário de A que não é o daquela identidade
		{ctxA, eventoIdentidade(a.ID, f.UsuarioF1, a.Identidade, "usuario", "login", "identidade")},
		{ctxA, eventoIdentidade(a.ID, a.UsuarioInterno, f.IdentidadeF1, "usuario", "login", "identidade")},
		// identidade de B em A
		{ctxA, eventoIdentidade(a.ID, b.UsuarioInterno, b.Identidade, "usuario", "login", "identidade")},
		// evento de B gravado na transação de A
		{ctxA, eventoIdentidade(b.ID, b.UsuarioInterno, b.Identidade, "usuario", "login", "identidade")},
		// sem empresa no contexto
		{contexto{}, eventoIdentidade(a.ID, a.UsuarioInterno, a.Identidade, "usuario", "login", "identidade")},
		{ctxA, tentativa{"coluna antes, fora do GRANT",
			`INSERT INTO app.evento_auditoria (empresa_id, ator_tipo, usuario_id, acao, entidade, entidade_id, antes)
			 VALUES ($1, 'usuario', $2, 'alteracao', 'identidade', $3, '{}')`, []any{a.ID, a.UsuarioInterno, a.Identidade}}},
		{ctxA, tentativa{"coluna conexao_empresa_id, fora do GRANT",
			`INSERT INTO app.evento_auditoria (empresa_id, ator_tipo, usuario_id, acao, entidade, entidade_id, conexao_empresa_id)
			 VALUES ($1, 'usuario', $2, 'login', 'identidade', $3, $1)`, []any{a.ID, a.UsuarioInterno, a.Identidade}}},
		{ctxA, tentativa{"lê a trilha", `SELECT count(*) FROM app.evento_auditoria`, nil}},
		{ctxA, tentativa{"altera a trilha", `UPDATE app.evento_auditoria SET acao = acao WHERE false`, nil}},
		{ctxA, tentativa{"apaga a trilha", `DELETE FROM app.evento_auditoria WHERE false`, nil}},
		{ctxA, tentativa{"lê usuários", `SELECT count(*) FROM app.usuario`, nil}},
		{ctxA, tentativa{"lê empresas", `SELECT count(*) FROM app.empresa`, nil}},
	}
	for _, r := range recusados {
		if err := tentarComo(ctx, tx, papel, r.c, r.e.sql, r.e.args...); sqlstate(err) != erroPrivilegio {
			problemas = append(problemas, fmt.Sprintf("não recusado com %s: %s (resultado: %v)", erroPrivilegio, r.e.rotulo, err))
		}
	}
	return problemas
}

func TestAuditoriaDoPapelDeAutenticacaoSegueOMolde(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		// O usuário de conexão herda o papel: as policies TO tenaz_autenticacao
		// precisam valer para ele, que é quem de fato conecta.
		for _, papel := range []string{papelAutenticacao, papelAutenticacaoConexao} {
			for _, p := range problemasDoMolde(t, ctx, tx, papel, f) {
				t.Errorf("%s: %s (ADR-011, decisão 5)", papel, p)
			}
		}
	})
}

func TestDetectorAcusaMoldeDeAuditoriaSemARestritiva(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		estragar(t, ctx, tx, "DROP POLICY autenticacao_identidade ON app.evento_auditoria")
		p := strings.Join(problemasDoMolde(t, ctx, tx, papelAutenticacao, f), " | ")
		for _, esperado := range []string{"entidade=usuario", "acao=exclusao", "ator=sistema"} {
			if !strings.Contains(p, "não recusado com 42501: ") || !strings.Contains(p, esperado) {
				t.Errorf("sem a restritiva autenticacao_identidade, o verificador não acusou o evento com %s; acusou: %s",
					esperado, p)
			}
		}
		if strings.Contains(p, "evento legítimo recusado") {
			t.Errorf("sem a restritiva, um evento legítimo foi recusado: %s", p)
		}
	})
}

// --------------------------------------- revogação pela inativação

// problemasDaRevogacao monta, na transação, uma identidade com acesso a A e a
// B, sessões web e do app nas duas, inativa (ou exclui) o usuário em A pelo
// caminho do ADM e confere: tudo de A revogado, nada de B tocado.
func problemasDaRevogacao(t *testing.T, ctx context.Context, tx pgx.Tx, f *fixturesE01, inativar string) []string {
	t.Helper()
	var identidade, usuarioA, usuarioB string
	sessoes := map[string]string{} // rótulo -> id
	tokens := map[string]string{}
	passo := func(rotulo, papel string, c contexto, sql string, saida *string, args ...any) {
		t.Helper()
		if err := trocarPapel(ctx, tx, papel, c); err != nil {
			t.Fatalf("%s: %v", rotulo, err)
		}
		var err error
		if saida != nil {
			err = tx.QueryRow(ctx, sql, args...).Scan(saida)
		} else {
			_, err = tx.Exec(ctx, sql, args...)
		}
		if err != nil {
			t.Fatalf("%s: %v", rotulo, err)
		}
	}

	passo("identidade", papelAutenticacao, contexto{}, `INSERT INTO plataforma.identidade (email) VALUES ($1) RETURNING id`,
		&identidade, "duas-empresas-"+novoID()+"@teste.tenaz")
	passo("usuário em A", papelApp, f.internoA(),
		`INSERT INTO app.usuario (empresa_id, nome, identidade_id) VALUES ($1, 'Duas empresas', $2) RETURNING id`,
		&usuarioA, f.A.ID, identidade)
	passo("usuário em B", papelApp, contexto{Empresa: f.B.ID, Escopo: "interno"},
		`INSERT INTO app.usuario (empresa_id, nome, identidade_id) VALUES ($1, 'Duas empresas', $2) RETURNING id`,
		&usuarioB, f.B.ID, identidade)
	for _, e := range []struct{ rotulo, empresa, usuario string }{{"A", f.A.ID, usuarioA}, {"B", f.B.ID, usuarioB}} {
		var web, app, token string
		passo("sessão web "+e.rotulo, papelAutenticacao, contexto{},
			`INSERT INTO plataforma.sessao (canal, token_hash, csrf_hash, identidade_id, empresa_id, usuario_id, expira_em)
			 VALUES ('web', $1, $2, $3, $4, $5, now() + interval '1 hour') RETURNING id`,
			&web, hashAleatorio(), hashAleatorio(), identidade, e.empresa, e.usuario)
		passo("sessão app "+e.rotulo, papelAutenticacao, contexto{},
			`INSERT INTO plataforma.sessao (canal, identidade_id, empresa_id, usuario_id, expira_em)
			 VALUES ('app', $1, $2, $3, now() + interval '30 days') RETURNING id`,
			&app, identidade, e.empresa, e.usuario)
		passo("refresh token "+e.rotulo, papelAutenticacao, contexto{},
			`INSERT INTO plataforma.token_renovacao (familia_id, token_hash, expira_em)
			 VALUES ($1, $2, now() + interval '30 days') RETURNING id`, &token, app, hashAleatorio())
		sessoes["web "+e.rotulo], sessoes["app "+e.rotulo], tokens[e.rotulo] = web, app, token
	}

	// A inativação é do ADM, pelo tenaz_app, que não tem privilégio nenhum nas
	// credenciais: quem revoga é o trigger.
	passo("inativar em A", papelApp, f.internoA(), `UPDATE app.usuario SET `+inativar+` WHERE id = $1`, nil, usuarioA)

	if _, err := tx.Exec(ctx, "RESET ROLE"); err != nil {
		t.Fatalf("voltar ao superusuário: %v", err)
	}
	var problemas []string
	revogada := func(sql, id string) bool {
		var r bool
		if err := tx.QueryRow(ctx, sql, id).Scan(&r); err != nil {
			t.Fatalf("conferir revogação: %v", err)
		}
		return r
	}
	for rotulo, id := range sessoes {
		r := revogada(`SELECT revogada_em IS NOT NULL FROM plataforma.sessao WHERE id = $1`, id)
		if deA := rotulo[len(rotulo)-1] == 'A'; r != deA {
			problemas = append(problemas, fmt.Sprintf("sessão %s: revogada=%v, esperado %v", rotulo, r, deA))
		}
	}
	for rotulo, id := range tokens {
		r := revogada(`SELECT revogado_em IS NOT NULL FROM plataforma.token_renovacao WHERE id = $1`, id)
		if deA := rotulo == "A"; r != deA {
			problemas = append(problemas, fmt.Sprintf("refresh token de %s: revogado=%v, esperado %v", rotulo, r, deA))
		}
	}
	for _, e := range []struct {
		empresa  string
		esperado bool
	}{{f.A.ID, false}, {f.B.ID, true}} {
		var habilitado bool
		if err := tx.QueryRow(ctx, `SELECT habilitado FROM plataforma.acesso_empresa WHERE identidade_id = $1 AND empresa_id = $2`,
			identidade, e.empresa).Scan(&habilitado); err != nil {
			t.Fatalf("ler acesso_empresa: %v", err)
		}
		if habilitado != e.esperado {
			problemas = append(problemas, fmt.Sprintf("acesso à empresa %.8s: habilitado=%v, esperado %v",
				e.empresa, habilitado, e.esperado))
		}
	}
	return problemas
}

func TestInativacaoRevogaCredenciaisSoNaEmpresa(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	for _, inativar := range []string{"situacao = 'inativo'", "excluido_em = now()"} {
		t.Run(inativar, func(t *testing.T) {
			emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
				for _, p := range problemasDaRevogacao(t, ctx, tx, f, inativar) {
					t.Errorf("%s (ADR-011, decisão 6)", p)
				}
			})
		})
	}
}

// A prova troca o trigger por uma versão que revoga pela IDENTIDADE, em todas
// as empresas: o erro que derrubaria a pessoa em B por uma decisão de A.
func TestDetectorAcusaRevogacaoQueAtravessaEmpresas(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		estragar(t, ctx, tx, `
			CREATE OR REPLACE FUNCTION plataforma.tg_acesso_empresa() RETURNS trigger
			  LANGUAGE plpgsql SECURITY DEFINER SET search_path = pg_catalog, pg_temp
			AS $$
			BEGIN
			  INSERT INTO plataforma.acesso_empresa (identidade_id, empresa_id, usuario_id, habilitado)
			  VALUES (NEW.identidade_id, NEW.empresa_id, NEW.id, NEW.situacao = 'ativo' AND NEW.excluido_em IS NULL)
			  ON CONFLICT (identidade_id, empresa_id) DO UPDATE SET habilitado = EXCLUDED.habilitado;
			  IF NOT (NEW.situacao = 'ativo' AND NEW.excluido_em IS NULL) THEN
			    UPDATE plataforma.sessao SET revogada_em = pg_catalog.now() WHERE identidade_id = NEW.identidade_id;
			  END IF;
			  RETURN NULL;
			END
			$$`)
		p := strings.Join(problemasDaRevogacao(t, ctx, tx, f, "situacao = 'inativo'"), " | ")
		if !strings.Contains(p, "sessão web B: revogada=true") || !strings.Contains(p, "sessão app B: revogada=true") {
			t.Errorf("um trigger que revoga as sessões de B ao inativar em A não foi acusado; acusou: %s", p)
		}
	})
}

// --------------------------------------------------------------- purga

// problemasDaPurga grava credenciais vencidas há mais de 30 dias e outras
// vigentes ou vencidas há pouco, roda a purga como tenaz_job e confere que
// só as primeiras sumiram.
func problemasDaPurga(t *testing.T, ctx context.Context, tx pgx.Tx, f *fixturesE01) []string {
	t.Helper()
	ident, emp, usu := f.A.Identidade, f.A.ID, f.A.UsuarioInterno
	type linha struct{ tabela, id string }
	var velhas, vivas []linha

	inserir := func(destino *[]linha, tabela, sql string, args ...any) {
		t.Helper()
		var id string
		if err := tx.QueryRow(ctx, sql, args...).Scan(&id); err != nil {
			t.Fatalf("preparar %s: %v", tabela, err)
		}
		*destino = append(*destino, linha{tabela, id})
	}
	if _, err := tx.Exec(ctx, "RESET ROLE"); err != nil {
		t.Fatalf("voltar ao superusuário: %v", err)
	}
	var sessaoVelha string
	if err := tx.QueryRow(ctx, `INSERT INTO plataforma.sessao (canal, identidade_id, empresa_id, usuario_id, criado_em, ultimo_uso_em, expira_em)
		VALUES ('app', $1, $2, $3, now() - interval '70 days', now() - interval '70 days', now() - interval '31 days') RETURNING id`,
		ident, emp, usu).Scan(&sessaoVelha); err != nil {
		t.Fatalf("preparar sessão velha: %v", err)
	}
	velhas = append(velhas, linha{"plataforma.sessao", sessaoVelha})
	inserir(&velhas, "plataforma.token_renovacao", `INSERT INTO plataforma.token_renovacao (familia_id, token_hash, criado_em, expira_em)
		VALUES ($1, $2, now() - interval '70 days', now() - interval '41 days') RETURNING id`, sessaoVelha, hashAleatorio())
	inserir(&velhas, "plataforma.codigo_autorizacao", `INSERT INTO plataforma.codigo_autorizacao (codigo_hash, identidade_id, empresa_id,
		usuario_id, cliente_oauth, uri_retorno, desafio_pkce, criado_em, expira_em)
		VALUES ($1, $2, $3, $4, 'tenaz-app', 'tenaz://retorno', 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQ',
		        now() - interval '32 days', now() - interval '32 days' + interval '60 seconds') RETURNING id`,
		hashAleatorio(), ident, emp, usu)
	inserir(&velhas, "plataforma.desafio_mfa", `INSERT INTO plataforma.desafio_mfa (token_hash, identidade_id, canal, criado_em, expira_em)
		VALUES ($1, $2, 'web', now() - interval '32 days', now() - interval '32 days' + interval '5 minutes') RETURNING id`,
		hashAleatorio(), ident)
	inserir(&velhas, "plataforma.token_email", `INSERT INTO plataforma.token_email (token_hash, finalidade, email, identidade_id,
		criado_em, expira_em)
		VALUES ($1, 'redefinicao_senha', 'x@teste.tenaz', $2, now() - interval '32 days', now() - interval '32 days' + interval '1 hour')
		RETURNING id`, hashAleatorio(), ident)
	// Vencida há 10 dias: ainda serve à investigação de incidente.
	inserir(&vivas, "plataforma.sessao", `INSERT INTO plataforma.sessao (canal, identidade_id, empresa_id, usuario_id, criado_em, expira_em)
		VALUES ('app', $1, $2, $3, now() - interval '20 days', now() - interval '10 days') RETURNING id`, ident, emp, usu)
	vivas = append(vivas,
		linha{"plataforma.sessao", f.A.SessaoWeb}, linha{"plataforma.sessao", f.A.SessaoApp},
		linha{"plataforma.token_renovacao", f.A.TokenRenovacao}, linha{"plataforma.codigo_autorizacao", f.A.CodigoAutoriz},
		linha{"plataforma.desafio_mfa", f.A.DesafioMFA}, linha{"plataforma.token_email", f.A.TokenEmail},
		linha{"plataforma.sessao", f.B.SessaoWeb})

	if err := trocarPapel(ctx, tx, papelJob, contexto{}); err != nil {
		t.Fatalf("assumir tenaz_job: %v", err)
	}
	var apagadas int64
	if err := tx.QueryRow(ctx, `SELECT plataforma.purgar_credenciais_vencidas()`).Scan(&apagadas); err != nil {
		t.Fatalf("tenaz_job não conseguiu purgar: %v", err)
	}
	if _, err := tx.Exec(ctx, "RESET ROLE"); err != nil {
		t.Fatalf("voltar ao superusuário: %v", err)
	}

	var problemas []string
	existe := func(l linha) bool {
		var n int
		if err := tx.QueryRow(ctx, "SELECT count(*) FROM "+l.tabela+" WHERE id = $1", l.id).Scan(&n); err != nil {
			t.Fatalf("conferir %s: %v", l.tabela, err)
		}
		return n == 1
	}
	for _, l := range velhas {
		if existe(l) {
			problemas = append(problemas, "credencial vencida há mais de 30 dias sobreviveu: "+l.tabela)
		}
	}
	for _, l := range vivas {
		if !existe(l) {
			problemas = append(problemas, "credencial vigente ou vencida há pouco foi apagada: "+l.tabela+" "+l.id)
		}
	}
	if apagadas < int64(len(velhas)) {
		problemas = append(problemas, fmt.Sprintf("a purga devolveu %d, e apagou ao menos %d", apagadas, len(velhas)))
	}
	return problemas
}

func TestPurgaSoPeloJobESoDoVencido(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		for _, p := range problemasDaPurga(t, ctx, tx, f) {
			t.Errorf("%s (ADR-011, decisão 2)", p)
		}
		for _, papel := range []string{papelApp, papelLeitura, papelSinc, papelAutenticacao, papelAutenticacaoConexao} {
			err := tentarComo(ctx, tx, papel, contexto{}, `SELECT plataforma.purgar_credenciais_vencidas()`)
			exigirErro(t, papel+" executa a purga", err, erroPrivilegio)
		}
	})
}

// A prova troca a função por uma que apaga tudo o que vence no próximo
// século: credencial viva apagada precisa ser acusada.
func TestDetectorAcusaPurgaQueApagaCredencialViva(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		estragar(t, ctx, tx, `
			CREATE OR REPLACE FUNCTION plataforma.purgar_credenciais_vencidas() RETURNS bigint
			  LANGUAGE plpgsql VOLATILE SECURITY DEFINER SET search_path = pg_catalog, pg_temp
			AS $$
			BEGIN
			  DELETE FROM plataforma.codigo_autorizacao WHERE expira_em < pg_catalog.now() + interval '100 years';
			  DELETE FROM plataforma.desafio_mfa WHERE expira_em < pg_catalog.now() + interval '100 years';
			  DELETE FROM plataforma.token_email WHERE expira_em < pg_catalog.now() + interval '100 years';
			  DELETE FROM plataforma.token_renovacao WHERE expira_em < pg_catalog.now() + interval '100 years';
			  DELETE FROM plataforma.sessao WHERE expira_em < pg_catalog.now() + interval '100 years';
			  RETURN 0;
			END
			$$`)
		p := strings.Join(problemasDaPurga(t, ctx, tx, f), " | ")
		for _, esperado := range []string{"foi apagada: plataforma.sessao " + f.A.SessaoWeb,
			"foi apagada: plataforma.token_email " + f.A.TokenEmail, "a purga devolveu 0"} {
			if !strings.Contains(p, esperado) {
				t.Errorf("a purga que apaga credencial vigente não foi acusada por %q; acusou: %s", esperado, p)
			}
		}
	})
}
