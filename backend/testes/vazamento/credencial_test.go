package vazamento

import (
	"context"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
)

// ISO-E01-08: as tabelas de credencial eram do tenaz_dono, e o
// tenaz_manutencao, membro do dono com INHERIT, as lia (plataforma não tem
// RLS). A correção as deu a tenaz_credencial, papel NOLOGIN sem membro de
// runtime. O ataque em si (dono e manutenção lendo e escrevendo credenciais)
// está em TestAtaquesDoISOE0102RecebemPermissaoNegada, com papelDono e
// papelManutencao em papeisSemCredencial; aqui fica o papel.

// problemasDoDonoDasCredenciais: NOLOGIN, NOBYPASSRLS, não superusuário, não
// membro de papel nenhum, e nenhum membro (com INHERIT ou SET) além de quem
// aplica migração — superusuário ou tenaz_migrador.
func problemasDoDonoDasCredenciais(t *testing.T, ctx context.Context, q consultor) []string {
	t.Helper()
	return listarTexto(t, ctx, q, `
		SELECT 'tenaz_credencial ' || p
		FROM (
		  SELECT 'não existe' AS p WHERE NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'tenaz_credencial')
		  UNION ALL
		  SELECT 'faz login, contorna RLS ou é superusuário' FROM pg_roles
		  WHERE rolname = 'tenaz_credencial' AND (rolcanlogin OR rolbypassrls OR rolsuper OR rolcreaterole)
		  UNION ALL
		  SELECT 'tem como membro ' || m.rolname
		  FROM pg_auth_members am JOIN pg_roles m ON m.oid = am.member
		  WHERE am.roleid = 'tenaz_credencial'::regrole AND (am.inherit_option OR am.set_option)
		    AND NOT m.rolsuper AND m.rolname <> 'tenaz_migrador'
		  UNION ALL
		  SELECT 'é membro de ' || r.rolname
		  FROM pg_auth_members am JOIN pg_roles r ON r.oid = am.roleid
		  WHERE am.member = 'tenaz_credencial'::regrole AND (am.inherit_option OR am.set_option)
		  UNION ALL
		  -- lista explícita da decisão 9 do ADR-011, com herança transitiva
		  -- (pg_has_role segue a cadeia: credencial dada ao dono chega à manutenção)
		  SELECT 'alcançável por ' || r
		  FROM unnest(ARRAY['tenaz_app', 'tenaz_job', 'tenaz_leitura', 'tenaz_sinc', 'tenaz_sinc_conexao',
		                    'tenaz_autenticacao', 'tenaz_autenticacao_conexao', 'tenaz_manutencao', 'tenaz_dono']) r
		  WHERE pg_has_role(r, 'tenaz_credencial', 'USAGE') OR pg_has_role(r, 'tenaz_credencial', 'SET')
		) x
		ORDER BY 1`)
}

func TestDonoDasCredenciaisSemMembroDeRuntime(t *testing.T) {
	ctx := context.Background()
	for _, p := range problemasDoDonoDasCredenciais(t, ctx, pool) {
		t.Errorf("%s (ISO-E01-08)", p)
	}
}

func TestDetectorAcusaDonoDasCredenciaisAlcancavel(t *testing.T) {
	ctx := context.Background()
	tx := criarTabelasDeProva(t, ctx, []string{
		"GRANT tenaz_credencial TO tenaz_dono",
		"GRANT tenaz_leitura TO tenaz_credencial",
		"ALTER ROLE tenaz_credencial LOGIN",
	})
	p := problemasDoDonoDasCredenciais(t, ctx, tx)
	for _, esperado := range []string{"tenaz_credencial tem como membro tenaz_dono",
		// transitivo: a manutenção, membro do dono, passa a alcançar as credenciais
		"tenaz_credencial alcançável por tenaz_manutencao",
		"tenaz_credencial é membro de tenaz_leitura", "tenaz_credencial faz login, contorna RLS ou é superusuário"} {
		if !contem(p, esperado) {
			t.Errorf("não acusou %q; acusou %v", esperado, p)
		}
	}
}

// O próprio defeito do ISO-E01-08, pela porta que ele abria: o tenaz_manutencao
// (e o dono) lendo a identidade de B. Fica aqui, além do ataque genérico, para
// que o achado tenha um teste com o seu nome.
func TestManutencaoEDonoNaoLeemCredenciais(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		for _, papel := range []string{papelManutencao, papelDono} {
			for _, tb := range tabelasCredencial {
				err := tentarComo(ctx, tx, papel, contexto{}, "SELECT count(*) FROM "+tb.nome)
				exigirErro(t, papel+" lê "+tb.nome, err, erroPrivilegio)
			}
			err := tentarComo(ctx, tx, papel, contexto{},
				`UPDATE plataforma.identidade SET senha_hash = 'x' WHERE id = $1`, f.B.Identidade)
			exigirErro(t, papel+" troca a senha da identidade de B", err, erroPrivilegio)
		}
	})
}

// O trigger que mantém acesso_empresa e revoga as credenciais na inativação
// existe e está habilitado: desligado, a inativação deixa de revogar a sessão
// (ADR-011, decisão 6) sem nenhum erro aparente.
const sqlTriggerDeAcesso = `
SELECT coalesce((SELECT CASE WHEN tg.tgenabled = 'D' THEN 'desabilitado'
                             WHEN tg.tgfoid <> 'plataforma.tg_acesso_empresa()'::regprocedure THEN 'com outra função'
                             ELSE 'ok' END
                 FROM pg_trigger tg
                 WHERE tg.tgrelid = 'app.usuario'::regclass AND tg.tgname = 'trg_acesso_empresa'), 'ausente')`

func TestTriggerDeAcessoExisteEHabilitado(t *testing.T) {
	ctx := context.Background()
	var estado string
	if err := pool.QueryRow(ctx, sqlTriggerDeAcesso).Scan(&estado); err != nil {
		t.Fatalf("ler o trigger: %v", err)
	}
	if estado != "ok" {
		t.Errorf("trg_acesso_empresa em app.usuario está %s", estado)
	}
}

func TestDetectorAcusaTriggerDeAcessoDesligado(t *testing.T) {
	ctx := context.Background()
	tx := criarTabelasDeProva(t, ctx, []string{"ALTER TABLE app.usuario DISABLE TRIGGER trg_acesso_empresa"})
	var estado string
	if err := tx.QueryRow(ctx, sqlTriggerDeAcesso).Scan(&estado); err != nil {
		t.Fatalf("ler o trigger: %v", err)
	}
	if estado != "desabilitado" {
		t.Errorf("o trigger desabilitado não foi acusado (estado: %s)", estado)
	}
}

// Regra geral (decisão 9 do ADR-011): tabela de plataforma sem RLS de dono
// tenaz_dono só pode ser catálogo de leitura ou controle; a de credencial é de
// tenaz_credencial. Senão, o tenaz_manutencao, membro do dono, a lê.
func problemasDeDonoEmPlataforma(t *testing.T, ctx context.Context, q consultor) []string {
	t.Helper()
	classes, err := lerClassesDeGlobais()
	if err != nil {
		t.Fatalf("ler %s: %v", arquivoListaGlobais, err)
	}
	linhas, err := q.Query(ctx, `
		SELECT n.nspname || '.' || c.relname, pg_get_userbyid(c.relowner), c.relrowsecurity AND c.relforcerowsecurity
		FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE n.nspname = 'plataforma' AND c.relkind IN ('r', 'p') AND NOT c.relispartition
		ORDER BY 1`)
	if err != nil {
		t.Fatalf("listar tabelas de plataforma: %v", err)
	}
	defer linhas.Close()
	var p []string
	for linhas.Next() {
		var nome, dono string
		var rls bool
		if err := linhas.Scan(&nome, &dono, &rls); err != nil {
			t.Fatalf("ler linha: %v", err)
		}
		classe := classes[nome]
		switch {
		case classe == "credencial" && dono != papelCredencial:
			p = append(p, nome+": credencial de dono "+dono+", esperado tenaz_credencial")
		case !rls && dono == papelDono && classe != "catalogo" && classe != "controle":
			p = append(p, nome+": sem RLS, de dono tenaz_dono, e não é catálogo de leitura (classe "+classe+")")
		}
	}
	if err := linhas.Err(); err != nil {
		t.Fatalf("iterar: %v", err)
	}
	return p
}

func TestTabelaDePlataformaDoDonoSoSeForCatalogo(t *testing.T) {
	ctx := context.Background()
	for _, p := range problemasDeDonoEmPlataforma(t, ctx, pool) {
		t.Errorf("%s (ADR-011, decisão 9; ISO-E01-08)", p)
	}
}

func TestDetectorAcusaTabelaDePlataformaDoDono(t *testing.T) {
	ctx := context.Background()
	tx := criarTabelasDeProva(t, ctx, []string{
		"SET LOCAL ROLE tenaz_dono",
		"CREATE TABLE plataforma.tabela_proposital_do_dono (id int)",
		"RESET ROLE",
		"ALTER TABLE plataforma.token_email OWNER TO tenaz_dono",
	})
	p := strings.Join(problemasDeDonoEmPlataforma(t, ctx, tx), " | ")
	for _, esperado := range []string{"plataforma.tabela_proposital_do_dono: sem RLS, de dono tenaz_dono",
		"plataforma.token_email: credencial de dono tenaz_dono"} {
		if !strings.Contains(p, esperado) {
			t.Errorf("não acusou %q; acusou: %s", esperado, p)
		}
	}
}
