package vazamento

import (
	"context"
	"strings"
	"testing"
)

// Consulta que encontra toda tabela de dados de cliente sem a defesa completa.
//
// É a verificação mais importante da suíte: ela não testa um caso, testa a
// COBERTURA. Um teste de caso prova que a tabela que alguém lembrou de testar
// está protegida; esta consulta prova que não existe tabela esquecida.
const sqlTabelasSemDefesa = `
WITH tabelas AS (
  SELECT c.oid, c.relname, c.relrowsecurity, c.relforcerowsecurity
  FROM pg_class c
  JOIN pg_namespace n ON n.oid = c.relnamespace
  WHERE n.nspname = 'app'
    AND c.relkind = 'r'
    AND EXISTS (
      SELECT 1 FROM pg_attribute a
      WHERE a.attrelid = c.oid AND a.attname = 'empresa_id' AND a.attnum > 0 AND NOT a.attisdropped
    )
)
SELECT t.relname,
       t.relrowsecurity      AS rls_habilitado,
       t.relforcerowsecurity AS rls_forcado,
       EXISTS (SELECT 1 FROM pg_policy p
               WHERE p.polrelid = t.oid AND p.polname = 'isolamento_empresa'
                 AND p.polpermissive AND p.polqual IS NOT NULL AND p.polwithcheck IS NOT NULL
              ) AS politica_empresa_completa,
       EXISTS (SELECT 1 FROM pg_policy p
               WHERE p.polrelid = t.oid AND p.polname = 'escopo_externo'
                 AND NOT p.polpermissive
              ) AS politica_escopo_restritiva
FROM tabelas t
WHERE NOT (t.relrowsecurity AND t.relforcerowsecurity)
   OR NOT EXISTS (SELECT 1 FROM pg_policy p
                  WHERE p.polrelid = t.oid AND p.polname = 'isolamento_empresa'
                    AND p.polpermissive AND p.polqual IS NOT NULL AND p.polwithcheck IS NOT NULL)
   OR NOT EXISTS (SELECT 1 FROM pg_policy p
                  WHERE p.polrelid = t.oid AND p.polname = 'escopo_externo'
                    AND NOT p.polpermissive)
ORDER BY t.relname
`

type tabelaSemDefesa struct {
	Nome               string
	RLSHabilitado      bool
	RLSForcado         bool
	PoliticaEmpresa    bool
	PoliticaRestritiva bool
}

func tabelasSemDefesa(t *testing.T, ctx context.Context) []tabelaSemDefesa {
	t.Helper()

	linhas, err := pool.Query(ctx, sqlTabelasSemDefesa)
	if err != nil {
		t.Fatalf("consultar cobertura de RLS: %v", err)
	}
	defer linhas.Close()

	var achados []tabelaSemDefesa
	for linhas.Next() {
		var a tabelaSemDefesa
		if err := linhas.Scan(&a.Nome, &a.RLSHabilitado, &a.RLSForcado, &a.PoliticaEmpresa, &a.PoliticaRestritiva); err != nil {
			t.Fatalf("ler linha: %v", err)
		}
		achados = append(achados, a)
	}
	if err := linhas.Err(); err != nil {
		t.Fatalf("iterar: %v", err)
	}
	return achados
}

// RNF-013: toda tabela com empresa_id tem RLS habilitado E forçado, política de
// empresa com USING e WITH CHECK, e política de escopo externo restritiva.
func TestCoberturaDeRLSEmTodaTabelaDeCliente(t *testing.T) {
	ctx := context.Background()

	for _, a := range tabelasSemDefesa(t, ctx) {
		t.Errorf("app.%s está sem defesa completa: rls_habilitado=%v rls_forcado=%v politica_empresa=%v politica_escopo_restritiva=%v",
			a.Nome, a.RLSHabilitado, a.RLSForcado, a.PoliticaEmpresa, a.PoliticaRestritiva)
	}
}

// Prova de que o detector funciona.
//
// Sem este teste, a verificação acima passaria por vacuidade enquanto não
// houvesse tabela de negócio — e passaria também se a consulta estivesse
// errada. Aqui a suíte cria de propósito uma tabela com empresa_id e sem
// nenhuma política e exige que o detector a acuse.
func TestDetectorAcusaTabelaSemPolitica(t *testing.T) {
	ctx := context.Background()

	const nome = "tabela_proposital_sem_politica"

	t.Cleanup(func() {
		if _, err := pool.Exec(ctx, "DROP TABLE IF EXISTS app."+nome); err != nil {
			t.Errorf("limpar tabela de prova: %v", err)
		}
	})

	if _, err := pool.Exec(ctx,
		"CREATE TABLE app."+nome+" (id uuid PRIMARY KEY, empresa_id uuid NOT NULL)"); err != nil {
		t.Fatalf("criar tabela de prova: %v", err)
	}

	achados := tabelasSemDefesa(t, ctx)

	var encontrou bool
	for _, a := range achados {
		if a.Nome == nome {
			encontrou = true
			if a.RLSHabilitado || a.RLSForcado || a.PoliticaEmpresa || a.PoliticaRestritiva {
				t.Errorf("tabela de prova deveria estar sem nenhuma defesa, veio %+v", a)
			}
		}
	}

	if !encontrou {
		t.Fatal("o detector de cobertura NÃO acusou uma tabela criada sem política; " +
			"a verificação de RNF-013 está cega e precisa ser corrigida antes de qualquer tabela de negócio entrar")
	}
}

// Catálogos globais ficam numa lista explícita, revisada em revisão de código.
// Tabela fora da lista e sem empresa_id quebra o build (seção 9.3).
func TestTabelasSemEmpresaIDEstaoNaListaDeGlobais(t *testing.T) {
	ctx := context.Background()

	permitidas, err := lerListaDeGlobais()
	if err != nil {
		t.Fatalf("ler lista de tabelas globais: %v", err)
	}

	linhas, err := pool.Query(ctx, `
		SELECT n.nspname || '.' || c.relname
		FROM pg_class c
		JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE n.nspname IN ('app', 'plataforma')
		  AND c.relkind = 'r'
		  AND NOT EXISTS (
		    SELECT 1 FROM pg_attribute a
		    WHERE a.attrelid = c.oid AND a.attname = 'empresa_id'
		      AND a.attnum > 0 AND NOT a.attisdropped
		  )
		ORDER BY 1`)
	if err != nil {
		t.Fatalf("consultar tabelas sem empresa_id: %v", err)
	}
	defer linhas.Close()

	for linhas.Next() {
		var nome string
		if err := linhas.Scan(&nome); err != nil {
			t.Fatalf("ler linha: %v", err)
		}
		if !permitidas[nome] {
			t.Errorf("%s não tem empresa_id e não está em testes/vazamento/tabelas_globais.txt; "+
				"se é catálogo global, acrescente-a à lista com justificativa; se não é, dê-lhe empresa_id", nome)
		}
	}
	if err := linhas.Err(); err != nil {
		t.Fatalf("iterar: %v", err)
	}
}

// Todo objeto de app e plataforma pertence a tenaz_dono.
//
// Não é organização: FORCE ROW LEVEL SECURITY não vale para o dono da tabela,
// então uma tabela de dono diferente — em especial um superusuário — contorna
// o RLS silenciosamente. E ALTER DEFAULT PRIVILEGES FOR ROLE tenaz_dono só
// alcança objetos criados por tenaz_dono.
func TestTodoObjetoPertenceAoDono(t *testing.T) {
	ctx := context.Background()

	linhas, err := pool.Query(ctx, `
		SELECT 'tabela ' || n.nspname || '.' || c.relname, pg_get_userbyid(c.relowner)
		FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE n.nspname IN ('app', 'plataforma')
		  AND c.relkind IN ('r', 'v', 'm', 'p')
		  AND pg_get_userbyid(c.relowner) <> 'tenaz_dono'
		UNION ALL
		SELECT 'função ' || n.nspname || '.' || p.proname, pg_get_userbyid(p.proowner)
		FROM pg_proc p JOIN pg_namespace n ON n.oid = p.pronamespace
		WHERE n.nspname IN ('app', 'plataforma')
		  AND pg_get_userbyid(p.proowner) <> 'tenaz_dono'
		UNION ALL
		SELECT 'esquema ' || n.nspname, pg_get_userbyid(n.nspowner)
		FROM pg_namespace n
		WHERE n.nspname IN ('app', 'plataforma')
		  AND pg_get_userbyid(n.nspowner) <> 'tenaz_dono'
		ORDER BY 1`)
	if err != nil {
		t.Fatalf("consultar donos: %v", err)
	}
	defer linhas.Close()

	for linhas.Next() {
		var objeto, dono string
		if err := linhas.Scan(&objeto, &dono); err != nil {
			t.Fatalf("ler linha: %v", err)
		}
		t.Errorf("%s pertence a %q, deveria pertencer a tenaz_dono; "+
			"a migração que o criou esqueceu o SET ROLE tenaz_dono", objeto, dono)
	}
	if err := linhas.Err(); err != nil {
		t.Fatalf("iterar: %v", err)
	}
}

// Nenhum papel da aplicação contorna o RLS. Só o papel que aplica migração
// pode, e ele não atende requisição.
func TestNenhumPapelDaAplicacaoContornaRLS(t *testing.T) {
	ctx := context.Background()

	linhas, err := pool.Query(ctx, `
		SELECT rolname, rolbypassrls, rolsuper
		FROM pg_roles
		WHERE rolname IN ('tenaz_app', 'tenaz_job', 'tenaz_leitura', 'tenaz_sinc', 'tenaz_sinc_conexao')
		ORDER BY rolname`)
	if err != nil {
		t.Fatalf("consultar papéis: %v", err)
	}
	defer linhas.Close()

	vistos := map[string]bool{}
	for linhas.Next() {
		var nome string
		var bypass, super bool
		if err := linhas.Scan(&nome, &bypass, &super); err != nil {
			t.Fatalf("ler linha: %v", err)
		}
		vistos[nome] = true
		if bypass {
			t.Errorf("papel %s tem BYPASSRLS; nenhum papel que atende requisição pode ter", nome)
		}
		if super {
			t.Errorf("papel %s é superusuário; superusuário ignora RLS", nome)
		}
	}
	if err := linhas.Err(); err != nil {
		t.Fatalf("iterar: %v", err)
	}

	for _, esperado := range []string{"tenaz_app", "tenaz_job", "tenaz_leitura", "tenaz_sinc", "tenaz_sinc_conexao"} {
		if !vistos[esperado] {
			t.Errorf("papel %s não existe; a migração de bootstrap não foi aplicada por completo", esperado)
		}
	}
}

// O papel do sincronizador não pode ganhar privilégio por tabela inteira nem
// ficar dono de tabela: ele recebe GRANT por coluna, só nos campos da lista
// fechada da conexão (risco R5). Em E00 ele ainda não tem privilégio nenhum, e
// é exatamente isso que se verifica.
func TestSincronizadorNaoEDonoNemTemPrivilegioAmplo(t *testing.T) {
	ctx := context.Background()

	var donoDeAlgo int
	if err := pool.QueryRow(ctx, `
		SELECT count(*) FROM pg_class c
		JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE n.nspname IN ('app', 'plataforma')
		  AND pg_get_userbyid(c.relowner) IN ('tenaz_sinc', 'tenaz_sinc_conexao')`).Scan(&donoDeAlgo); err != nil {
		t.Fatalf("consultar propriedade do sincronizador: %v", err)
	}
	if donoDeAlgo != 0 {
		t.Errorf("tenaz_sinc é dono de %d objeto(s); o dono ignora RLS forçado", donoDeAlgo)
	}

	// Privilégio concedido à tabela toda, em vez de por coluna.
	linhas, err := pool.Query(ctx, `
		SELECT table_schema || '.' || table_name, privilege_type
		FROM information_schema.table_privileges
		WHERE grantee = 'tenaz_sinc' AND table_schema IN ('app', 'plataforma')
		ORDER BY 1, 2`)
	if err != nil {
		t.Fatalf("consultar privilégios do sincronizador: %v", err)
	}
	defer linhas.Close()

	for linhas.Next() {
		var tabela, privilegio string
		if err := linhas.Scan(&tabela, &privilegio); err != nil {
			t.Fatalf("ler linha: %v", err)
		}
		t.Errorf("tenaz_sinc tem %s na tabela inteira %s; o privilégio da sincronização é por coluna, "+
			"restrito aos campos da lista fechada da conexão (RNF-041)", privilegio, tabela)
	}
	if err := linhas.Err(); err != nil {
		t.Fatalf("iterar: %v", err)
	}
}

// Toda migração que cria objeto assume o papel do dono. A 00001 é a exceção
// declarada: é ela que cria os papéis.
func TestMigracoesAssumemOPapelDoDono(t *testing.T) {
	arquivos, err := listarMigracoes()
	if err != nil {
		t.Fatalf("listar migrações: %v", err)
	}
	if len(arquivos) == 0 {
		t.Fatal("nenhuma migração encontrada")
	}

	for _, a := range arquivos {
		if strings.HasPrefix(a.Nome, "00001_") {
			continue // bootstrap: cria os próprios papéis
		}
		if !strings.Contains(a.Conteudo, "SET ROLE tenaz_dono") {
			t.Errorf("%s não contém `SET ROLE tenaz_dono`; os objetos que ela criar ficarão com o dono errado "+
				"e escaparão do RLS forçado", a.Nome)
		}
	}
}
