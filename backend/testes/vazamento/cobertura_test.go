package vazamento

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
)

// Expressões exatas das policies, como o PostgreSQL as devolve em
// pg_get_expr. Comparar o texto, e não só "a policy existe", é o que pega uma
// isolamento_empresa escrita com USING (true) ou com a coluna errada.
const (
	exprEmpresa       = "(empresa_id = app.empresa_atual())"
	exprEscopoInterno = "app.escopo_interno()"
)

// Consulta que encontra toda tabela de dados de cliente sem a defesa completa.
//
// É a verificação mais importante da suíte: ela não testa um caso, testa a
// COBERTURA. Um teste de caso prova que a tabela que alguém lembrou de testar
// está protegida; esta consulta prova que não existe tabela esquecida.
//
// Tabela particionada: olha o PAI (relkind 'p') e pula as partições
// (relispartition). O RLS de uma consulta pelo pai vale pelas policies do
// pai; a partição é cobrada à parte, em TestParticoesHerdamADefesaDoPai
// (consultas 1 e 4 da seção 2.8 do documento 2).
//
// $1: tabelas com exceção versionada em excecoes_cobertura.txt, que têm a
// forma conferida por TestExcecoesDeCoberturaTemAFormaDeclarada.
const sqlTabelasSemDefesa = `
WITH tabelas AS (
  SELECT c.oid, n.nspname || '.' || c.relname AS nome, c.relrowsecurity, c.relforcerowsecurity
  FROM pg_class c
  JOIN pg_namespace n ON n.oid = c.relnamespace
  WHERE n.nspname = 'app'
    AND c.relkind IN ('r', 'p')
    AND NOT c.relispartition
    AND EXISTS (
      SELECT 1 FROM pg_attribute a
      WHERE a.attrelid = c.oid AND a.attname = 'empresa_id' AND a.attnum > 0 AND NOT a.attisdropped
    )
    AND n.nspname || '.' || c.relname <> ALL ($1::text[])
), avaliadas AS (
  SELECT t.*,
         EXISTS (SELECT 1 FROM pg_policy p
                 WHERE p.polrelid = t.oid AND p.polname = 'isolamento_empresa'
                   AND p.polpermissive AND p.polcmd = '*'
                   AND pg_get_expr(p.polqual, p.polrelid) = '` + exprEmpresa + `'
                   AND pg_get_expr(p.polwithcheck, p.polrelid) = '` + exprEmpresa + `'
                   AND (p.polroles @> ARRAY['tenaz_app'::regrole::oid] OR p.polroles = '{0}')
                ) AS politica_empresa,
         EXISTS (SELECT 1 FROM pg_policy p
                 WHERE p.polrelid = t.oid AND p.polname = 'escopo_externo'
                   AND NOT p.polpermissive AND p.polcmd = '*'
                   AND (p.polroles @> ARRAY['tenaz_app'::regrole::oid] OR p.polroles = '{0}')
                ) AS politica_escopo
  FROM tabelas t
)
SELECT nome, relrowsecurity, relforcerowsecurity, politica_empresa, politica_escopo
FROM avaliadas
WHERE NOT (relrowsecurity AND relforcerowsecurity AND politica_empresa AND politica_escopo)
ORDER BY nome
`

type tabelaSemDefesa struct {
	Nome               string
	RLSHabilitado      bool
	RLSForcado         bool
	PoliticaEmpresa    bool
	PoliticaRestritiva bool
}

func tabelasSemDefesa(t *testing.T, ctx context.Context, q consultor) []tabelaSemDefesa {
	t.Helper()

	excecoes, err := lerLista(arquivoExcecoesCobertura)
	if err != nil {
		t.Fatalf("ler %s: %v", arquivoExcecoesCobertura, err)
	}

	linhas, err := q.Query(ctx, sqlTabelasSemDefesa, excecoes.nomes())
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
// empresa com USING e WITH CHECK exatos, e política de escopo externo restritiva.
func TestCoberturaDeRLSEmTodaTabelaDeCliente(t *testing.T) {
	ctx := context.Background()

	for _, a := range tabelasSemDefesa(t, ctx, pool) {
		t.Errorf("%s está sem defesa completa: rls_habilitado=%v rls_forcado=%v isolamento_empresa_exata=%v escopo_externo_restritiva=%v",
			a.Nome, a.RLSHabilitado, a.RLSForcado, a.PoliticaEmpresa, a.PoliticaRestritiva)
	}
}

// Prova de que o detector funciona.
//
// Sem este teste, a verificação acima passaria por vacuidade enquanto não
// houvesse tabela de negócio — e passaria também se a consulta estivesse
// errada. Aqui a suíte cria de propósito tabelas malformadas e exige que o
// detector as acuse: uma comum sem nenhuma política, uma com a política de
// empresa permissiva demais, e uma particionada sem RLS (o caso que a
// consulta da E00, com relkind = 'r', deixava passar).
func TestDetectorAcusaTabelaSemPolitica(t *testing.T) {
	ctx := context.Background()

	const (
		semPolitica  = "app.tabela_proposital_sem_politica"
		politicaFrou = "app.tabela_proposital_politica_frouxa"
		particionada = "app.tabela_proposital_particionada"
		particao     = "app.tabela_proposital_particionada_p1"
	)

	tx := criarTabelasDeProva(t, ctx, []string{
		"CREATE TABLE " + semPolitica + " (id uuid PRIMARY KEY, empresa_id uuid NOT NULL)",

		"CREATE TABLE " + politicaFrou + " (id uuid PRIMARY KEY, empresa_id uuid NOT NULL)",
		"ALTER TABLE " + politicaFrou + " ENABLE ROW LEVEL SECURITY",
		"ALTER TABLE " + politicaFrou + " FORCE ROW LEVEL SECURITY",
		"CREATE POLICY isolamento_empresa ON " + politicaFrou + " FOR ALL TO tenaz_app USING (true) WITH CHECK (true)",
		"CREATE POLICY escopo_externo ON " + politicaFrou + " AS RESTRICTIVE FOR ALL TO tenaz_app " +
			"USING (app.escopo_interno()) WITH CHECK (app.escopo_interno())",

		"CREATE TABLE " + particionada + " (empresa_id uuid NOT NULL, dia date NOT NULL) PARTITION BY RANGE (dia)",
		"CREATE TABLE " + particao + " PARTITION OF " + particionada + " FOR VALUES FROM ('2000-01-01') TO ('2100-01-01')",
	})

	acusadas := map[string]tabelaSemDefesa{}
	for _, a := range tabelasSemDefesa(t, ctx, tx) {
		acusadas[a.Nome] = a
	}

	for _, nome := range []string{semPolitica, politicaFrou, particionada} {
		if _, ok := acusadas[nome]; !ok {
			t.Errorf("o detector de cobertura NÃO acusou %s, criada sem defesa; a verificação de RNF-013 está cega "+
				"e precisa ser corrigida antes de qualquer tabela de negócio entrar", nome)
		}
	}
	if a := acusadas[semPolitica]; a.RLSHabilitado || a.RLSForcado || a.PoliticaEmpresa || a.PoliticaRestritiva {
		t.Errorf("tabela sem política deveria estar sem nenhuma defesa, veio %+v", a)
	}
	if a := acusadas[politicaFrou]; a.PoliticaEmpresa {
		t.Errorf("isolamento_empresa com USING (true) foi aceita como política de empresa: %+v", a)
	}
	// A partição não é avaliada como tabela própria: herda a cobrança do pai.
	if _, ok := acusadas[particao]; ok {
		t.Errorf("a partição %s foi avaliada como tabela solta; a cobrança é do pai", particao)
	}
	// E é acusada pela verificação de partições, por estar sem RLS.
	if !contem(particoesSemDefesa(t, ctx, tx), particao) {
		t.Errorf("a partição %s, sem RLS, não foi acusada pela verificação de partições", particao)
	}
}

// criarTabelasDeProva abre uma transação do superusuário, executa nela os
// comandos que estragam o banco e a devolve; a transação é desfeita no fim do
// teste. Nada da prova chega a existir para outra conexão: o banco
// compartilhado (TENAZ_TESTE_BANCO_URL) nunca fica com objeto de prova, nem
// se o teste for interrompido.
func criarTabelasDeProva(t *testing.T, ctx context.Context, comandos []string) pgx.Tx {
	t.Helper()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("abrir transação da prova: %v", err)
	}
	t.Cleanup(func() { _ = tx.Rollback(context.Background()) })
	estragar(t, ctx, tx, comandos...)
	return tx
}

// estragar executa os comandos da prova num savepoint, com lock_timeout curto,
// e repete se a trava não vier ou se o banco escolher esta transação como
// vítima de deadlock (55P03, 40P01). Com um banco compartilhado entre pacotes
// (tarefas.sh testar-rapido), outra conexão pode estar escrevendo na tabela
// que a prova altera: a prova espera a vez, nunca segura a outra por muito
// tempo, e não falha por concorrência alheia. Erro de outro tipo é falha.
func estragar(t *testing.T, ctx context.Context, tx pgx.Tx, comandos ...string) {
	t.Helper()
	const tentativas = 8
	for i := 1; ; i++ {
		err := func() error {
			sub, err := tx.Begin(ctx)
			if err != nil {
				return err
			}
			defer func() { _ = sub.Rollback(ctx) }()
			if _, err := sub.Exec(ctx, "SET LOCAL lock_timeout = '1s'"); err != nil {
				return err
			}
			for _, c := range comandos {
				if _, err := sub.Exec(ctx, c); err != nil {
					return fmt.Errorf("%s: %w", c, err)
				}
			}
			// O SET LOCAL de um savepoint liberado sobrevive na transação: o
			// limite volta ao padrão para o resto da prova.
			if _, err := sub.Exec(ctx, "SET LOCAL lock_timeout = 0"); err != nil {
				return err
			}
			return sub.Commit(ctx)
		}()
		if err == nil {
			return
		}
		if codigo := sqlstate(err); (codigo == "55P03" || codigo == "40P01") && i < tentativas {
			time.Sleep(time.Duration(i*150) * time.Millisecond)
			continue
		}
		t.Fatalf("preparar prova: %v", err)
	}
}

// sqlPrivilegioComColuna é o teste de privilégio que vê também o concedido
// por coluna (SELECT, INSERT, UPDATE, REFERENCES): has_table_privilege só vê o
// da tabela inteira, e um GRANT UPDATE (coluna) passaria calado.
func sqlPrivilegioComColuna(papel, objeto, privilegio string) string {
	return "CASE WHEN " + privilegio + " IN ('SELECT', 'INSERT', 'UPDATE', 'REFERENCES') THEN has_any_column_privilege(" +
		papel + ", " + objeto + ", " + privilegio + ") ELSE has_table_privilege(" + papel + ", " + objeto + ", " + privilegio + ") END"
}

func contem(lista []string, item string) bool {
	for _, x := range lista {
		if x == item {
			return true
		}
	}
	return false
}

// Partição: RLS habilitado e forçado, sem privilégio para nenhum papel da
// aplicação, e com os triggers do pai (o trg_somente_insercao, na auditoria).
// O acesso pelo pai usa as policies do pai; o acesso direto à partição não
// pode ser caminho para contorná-las (consulta 4 da seção 2.8; seção 9.5).
const sqlParticoesSemDefesa = `
SELECT n.nspname || '.' || c.relname
FROM pg_class c
JOIN pg_namespace n ON n.oid = c.relnamespace
WHERE n.nspname = 'app' AND c.relispartition AND c.relkind IN ('r', 'p')
  AND (NOT c.relrowsecurity OR NOT c.relforcerowsecurity
       OR EXISTS (SELECT 1 FROM pg_policy p WHERE p.polrelid = c.oid AND p.polpermissive)
       OR EXISTS (SELECT 1
                  FROM unnest(` + sqlPapeisDaAplicacao + `) r,
                       unnest(ARRAY['SELECT', 'INSERT', 'UPDATE', 'DELETE', 'TRUNCATE', 'REFERENCES', 'TRIGGER']) p
                  WHERE CASE WHEN p IN ('SELECT', 'INSERT', 'UPDATE', 'REFERENCES') THEN has_any_column_privilege(r, c.oid, p) ELSE has_table_privilege(r, c.oid, p) END)
       OR EXISTS (SELECT 1 FROM pg_trigger tp
                  WHERE tp.tgrelid = (SELECT i.inhparent FROM pg_inherits i WHERE i.inhrelid = c.oid)
                    AND NOT tp.tgisinternal
                    AND NOT EXISTS (SELECT 1 FROM pg_trigger tc
                                    WHERE tc.tgrelid = c.oid AND tc.tgparentid = tp.oid)))
ORDER BY 1`

func particoesSemDefesa(t *testing.T, ctx context.Context, q consultor) []string {
	t.Helper()
	return listarTexto(t, ctx, q, sqlParticoesSemDefesa)
}

func TestParticoesHerdamADefesaDoPai(t *testing.T) {
	ctx := context.Background()
	for _, p := range particoesSemDefesa(t, ctx, pool) {
		t.Errorf("a partição %s está sem RLS forçado, tem policy permissiva, privilégio de papel da aplicação "+
			"ou falta-lhe um trigger do pai; crie partições só por plataforma.criar_particoes_mensais()", p)
	}
}

// listarTexto devolve a primeira coluna de cada linha.
func listarTexto(t *testing.T, ctx context.Context, q consultor, sql string, args ...any) []string {
	t.Helper()
	linhas, err := q.Query(ctx, sql, args...)
	if err != nil {
		t.Fatalf("consultar catálogo: %v", err)
	}
	defer linhas.Close()
	var saida []string
	for linhas.Next() {
		var s string
		if err := linhas.Scan(&s); err != nil {
			t.Fatalf("ler linha: %v", err)
		}
		saida = append(saida, s)
	}
	if err := linhas.Err(); err != nil {
		t.Fatalf("iterar: %v", err)
	}
	return saida
}

// Policies permissivas se somam com OR. Uma segunda policy permissiva com
// predicado mais largo que o da empresa abre a tabela inteira, mesmo com a
// isolamento_empresa intacta. Toda policy permissiva de tabela com empresa_id
// tem de usar exatamente o predicado da empresa.
const sqlPermissivasLargas = `
SELECT c.oid::regclass::text || '.' || p.polname || ': ' ||
       coalesce(pg_get_expr(p.polqual, p.polrelid), '-') || ' / ' || coalesce(pg_get_expr(p.polwithcheck, p.polrelid), '-')
FROM pg_policy p
JOIN pg_class c ON c.oid = p.polrelid
JOIN pg_namespace n ON n.oid = c.relnamespace
WHERE n.nspname = 'app' AND p.polpermissive
  AND EXISTS (SELECT 1 FROM pg_attribute a
              WHERE a.attrelid = c.oid AND a.attname = 'empresa_id' AND a.attnum > 0 AND NOT a.attisdropped)
  AND (coalesce(pg_get_expr(p.polqual, p.polrelid), '` + exprEmpresa + `') <> '` + exprEmpresa + `'
       OR coalesce(pg_get_expr(p.polwithcheck, p.polrelid), '` + exprEmpresa + `') <> '` + exprEmpresa + `')
ORDER BY 1`

func TestPoliticasPermissivasSoComparamAEmpresa(t *testing.T) {
	ctx := context.Background()
	for _, p := range listarTexto(t, ctx, pool, sqlPermissivasLargas) {
		t.Errorf("policy permissiva com predicado diferente de %s: %s; permissivas se somam com OR "+
			"e esta alarga o corte por empresa", exprEmpresa, p)
	}
}

func TestDetectorAcusaPermissivaLarga(t *testing.T) {
	ctx := context.Background()
	const nome = "app.tabela_proposital_permissiva_larga"
	tx := criarTabelasDeProva(t, ctx, []string{
		"CREATE TABLE " + nome + " (id uuid PRIMARY KEY, empresa_id uuid NOT NULL)",
		"CREATE POLICY isolamento_empresa ON " + nome + " FOR ALL TO tenaz_app " +
			"USING (empresa_id = app.empresa_atual()) WITH CHECK (empresa_id = app.empresa_atual())",
		"CREATE POLICY leitura_geral ON " + nome + " FOR SELECT TO tenaz_app USING (true)",
	})

	var acusou bool
	for _, p := range listarTexto(t, ctx, tx, sqlPermissivasLargas) {
		if strings.HasPrefix(p, nome+".leitura_geral") {
			acusou = true
		}
	}
	if !acusou {
		t.Fatal("uma policy permissiva USING (true) ao lado da isolamento_empresa não foi acusada")
	}
}

// ------------------------------------------------------------ exceções

// Tabela de app com empresa_id que, por desenho, não tem a isolamento_empresa
// FOR ALL nem a escopo_externo: está em excecoes_cobertura.txt com o nome da
// forma que a substitui, e a forma é conferida aqui. A exceção troca uma
// regra por outra mais estrita; nunca dispensa a defesa.
func TestExcecoesDeCoberturaTemAFormaDeclarada(t *testing.T) {
	ctx := context.Background()

	excecoes, err := lerLista(arquivoExcecoesCobertura)
	if err != nil {
		t.Fatalf("ler %s: %v", arquivoExcecoesCobertura, err)
	}
	for _, e := range excecoes {
		if len(e.Campos) != 1 {
			t.Errorf("%s: a linha precisa de exatamente uma forma, veio %v", e.Nome, e.Campos)
			continue
		}
		switch e.Campos[0] {
		case "somente_insercao":
			for _, p := range problemasSomenteInsercao(t, ctx, pool, e.Nome) {
				t.Errorf("%s (somente_insercao): %s", e.Nome, p)
			}
		default:
			t.Errorf("%s: forma %q desconhecida; formas aceitas: somente_insercao", e.Nome, e.Campos[0])
		}
	}
}

// problemasSomenteInsercao confere a forma de evento_auditoria (seções 9.2 e
// 9.4 do documento 2): RLS forçado; leitura e inserção só da empresa; nenhuma
// permissiva de UPDATE ou DELETE; leitura só interna; nenhum papel da
// aplicação com UPDATE, DELETE ou TRUNCATE; e o trigger de somente inserção.
func problemasSomenteInsercao(t *testing.T, ctx context.Context, q consultor, tabela string) []string {
	t.Helper()
	var problemas []string
	falhou := func(sql, problema string) {
		var ok bool
		if err := q.QueryRow(ctx, sql, tabela).Scan(&ok); err != nil {
			t.Fatalf("conferir %s em %s: %v", problema, tabela, err)
		}
		if !ok {
			problemas = append(problemas, problema)
		}
	}

	falhou(`SELECT relrowsecurity AND relforcerowsecurity FROM pg_class WHERE oid = $1::regclass`,
		"RLS não está habilitado e forçado")
	falhou(`SELECT NOT EXISTS (SELECT 1 FROM pg_policy WHERE polrelid = $1::regclass AND polpermissive
	                            AND polcmd IN ('*', 'w', 'd'))`,
		"há policy permissiva que alcança UPDATE ou DELETE")
	falhou(`SELECT EXISTS (SELECT 1 FROM pg_policy WHERE polrelid = $1::regclass AND polpermissive AND polcmd = 'r'
	                        AND pg_get_expr(polqual, polrelid) = '`+exprEmpresa+`')
	           AND EXISTS (SELECT 1 FROM pg_policy WHERE polrelid = $1::regclass AND polpermissive AND polcmd = 'a'
	                        AND pg_get_expr(polwithcheck, polrelid) = '`+exprEmpresa+`')`,
		"falta a permissiva de leitura ou de inserção com o predicado exato da empresa")
	falhou(`SELECT EXISTS (SELECT 1 FROM pg_policy WHERE polrelid = $1::regclass AND NOT polpermissive
	                        AND polcmd IN ('*', 'r') AND pg_get_expr(polqual, polrelid) = '`+exprEscopoInterno+`'
	                        AND polroles @> ARRAY['tenaz_app'::regrole::oid])`,
		"a leitura não é restrita ao escopo interno")
	falhou(`SELECT NOT EXISTS (SELECT 1
	          FROM unnest(`+sqlPapeisDaAplicacao+`) r,
	               unnest(ARRAY['UPDATE', 'DELETE', 'TRUNCATE']) p
	          WHERE `+sqlPrivilegioComColuna("r", "$1::regclass", "p")+`)`,
		"algum papel da aplicação tem UPDATE, DELETE ou TRUNCATE")
	falhou(`SELECT EXISTS (SELECT 1 FROM pg_trigger WHERE tgrelid = $1::regclass
	                        AND tgfoid = 'app.tg_somente_insercao'::regproc
	                        AND tgtype & 8 <> 0 AND tgtype & 16 <> 0)`,
		"falta o trigger app.tg_somente_insercao em UPDATE e DELETE")
	return problemas
}

// A prova do detector da forma: uma tabela "somente inserção" que aceita
// UPDATE por uma permissiva FOR ALL precisa ser acusada.
func TestDetectorAcusaSomenteInsercaoMalformada(t *testing.T) {
	ctx := context.Background()
	const nome = "app.tabela_proposital_somente_insercao"
	tx := criarTabelasDeProva(t, ctx, []string{
		"CREATE TABLE " + nome + " (id uuid PRIMARY KEY, empresa_id uuid NOT NULL)",
		"ALTER TABLE " + nome + " ENABLE ROW LEVEL SECURITY",
		"ALTER TABLE " + nome + " FORCE ROW LEVEL SECURITY",
		"CREATE POLICY isolamento_empresa ON " + nome + " FOR ALL TO tenaz_app " +
			"USING (empresa_id = app.empresa_atual()) WITH CHECK (empresa_id = app.empresa_atual())",
	})

	if len(problemasSomenteInsercao(t, ctx, tx, nome)) == 0 {
		t.Fatal("uma tabela com permissiva FOR ALL e sem trigger passou como somente_insercao")
	}
}

// ------------------------------------------------------ tabelas globais

// Toda tabela fora do RLS por empresa está numa lista explícita, revisada em
// revisão de código: as de app sem empresa_id (consulta 2 da seção 2.8) e
// TODAS as de qualquer outro esquema que não seja do sistema (consulta 8,
// estendida): em plataforma não há RLS, e um esquema novo com uma cópia de
// pessoa escaparia de todas as outras verificações (mutação m24).
const sqlTabelasForaDoRLS = `
SELECT n.nspname || '.' || c.relname
FROM pg_class c
JOIN pg_namespace n ON n.oid = c.relnamespace
WHERE c.relkind IN ('r', 'p', 'f') AND NOT c.relispartition
  AND n.nspname NOT IN ('pg_catalog', 'information_schema') AND n.nspname NOT LIKE 'pg\_%'
  AND (n.nspname <> 'app'
       OR NOT EXISTS (SELECT 1 FROM pg_attribute a
                      WHERE a.attrelid = c.oid AND a.attname = 'empresa_id' AND a.attnum > 0 AND NOT a.attisdropped))
ORDER BY 1`

// diferencasDasGlobais devolve as tabelas fora do RLS que não estão na lista e
// as linhas da lista que não existem no banco.
func diferencasDasGlobais(t *testing.T, ctx context.Context, q consultor) (foraDaLista, inexistentes []string) {
	t.Helper()
	permitidas, err := lerListaDeGlobais()
	if err != nil {
		t.Fatalf("ler lista de tabelas globais: %v", err)
	}
	existentes := map[string]bool{}
	for _, nome := range listarTexto(t, ctx, q, sqlTabelasForaDoRLS) {
		existentes[nome] = true
		if !permitidas[nome] {
			foraDaLista = append(foraDaLista, nome)
		}
	}
	for nome := range permitidas {
		if !existentes[nome] {
			inexistentes = append(inexistentes, nome)
		}
	}
	return foraDaLista, inexistentes
}

func TestTabelasSemEmpresaIDEstaoNaListaDeGlobais(t *testing.T) {
	ctx := context.Background()
	fora, inexistentes := diferencasDasGlobais(t, ctx, pool)
	for _, nome := range fora {
		t.Errorf("%s está fora do RLS por empresa e não está em testes/vazamento/tabelas_globais.txt; "+
			"se é global, acrescente-a à lista com justificativa e aprovação de isolamento; se não é, "+
			"dê-lhe empresa_id e mova-a para app", nome)
	}
	for _, nome := range inexistentes {
		t.Errorf("%s está em tabelas_globais.txt mas não existe no banco (ou ganhou empresa_id em app); "+
			"retire a linha", nome)
	}
}

func TestDetectorAcusaTabelaGlobalNaoDeclarada(t *testing.T) {
	ctx := context.Background()
	tx := criarTabelasDeProva(t, ctx, []string{
		"CREATE TABLE plataforma.tabela_proposital_global (id int PRIMARY KEY)",
		"CREATE SCHEMA rel_prova",
		"CREATE TABLE rel_prova.pessoa (empresa_id uuid, cpf text)",
		"CREATE TABLE app.tabela_proposital_sem_empresa (id int PRIMARY KEY)",
		"DROP TABLE plataforma.desafio_mfa",
	})
	fora, inexistentes := diferencasDasGlobais(t, ctx, tx)
	for _, esperado := range []string{"plataforma.tabela_proposital_global", "rel_prova.pessoa",
		"app.tabela_proposital_sem_empresa"} {
		if !contem(fora, esperado) {
			t.Errorf("a tabela %s, fora do RLS e fora da lista, não foi acusada (acusadas: %v)", esperado, fora)
		}
	}
	if !contem(inexistentes, "plataforma.desafio_mfa") {
		t.Errorf("a linha da lista sem tabela no banco não foi acusada (acusadas: %v)", inexistentes)
	}
}

// ------------------------------------------------------ dono

// Todo objeto de app e plataforma pertence a tenaz_dono, exceto as sete
// tabelas de credencial e as duas funções SECURITY DEFINER que as tocam, que
// pertencem a tenaz_credencial (ISO-E01-08): sendo do dono, elas seriam lidas
// pelo tenaz_manutencao, membro do dono, já que plataforma não tem RLS. A
// exceção é exata — nem a mais, nem a menos.
//
// Não é organização: FORCE ROW LEVEL SECURITY não vale para o dono da tabela,
// então uma tabela de dono diferente — em especial um superusuário — contorna
// o RLS silenciosamente. E ALTER DEFAULT PRIVILEGES FOR ROLE tenaz_dono só
// alcança objetos criados por tenaz_dono.
const sqlDonoErrado = `
WITH credencial(nome) AS (VALUES
  ('plataforma.identidade'), ('plataforma.acesso_empresa'), ('plataforma.sessao'), ('plataforma.token_renovacao'),
  ('plataforma.codigo_autorizacao'), ('plataforma.desafio_mfa'), ('plataforma.token_email'),
  ('plataforma.tg_acesso_empresa()'), ('plataforma.purgar_credenciais_vencidas()')
), objetos AS (
  -- base: a tabela de que o objeto depende (o próprio nome, para tabela e
  -- função); índice e tipo de linha seguem o dono da tabela base
  SELECT 'tabela' AS tipo, n.nspname || '.' || c.relname AS nome, pg_get_userbyid(c.relowner) AS dono,
         n.nspname || '.' || c.relname AS base
  FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
  WHERE n.nspname IN ('app', 'plataforma') AND c.relkind IN ('r', 'v', 'm', 'p', 'S', 'f')
  UNION ALL
  SELECT 'índice', n.nspname || '.' || ci.relname, pg_get_userbyid(ci.relowner), n.nspname || '.' || ct.relname
  FROM pg_index i
  JOIN pg_class ci ON ci.oid = i.indexrelid
  JOIN pg_class ct ON ct.oid = i.indrelid
  JOIN pg_namespace n ON n.oid = ct.relnamespace
  WHERE n.nspname IN ('app', 'plataforma')
  UNION ALL
  SELECT 'tipo de linha', format_type(ty.oid, NULL), pg_get_userbyid(ty.typowner), n.nspname || '.' || c.relname
  FROM pg_type ty
  JOIN pg_class c ON c.oid = ty.typrelid
  JOIN pg_namespace n ON n.oid = c.relnamespace
  WHERE n.nspname IN ('app', 'plataforma') AND c.relkind IN ('r', 'p', 'v', 'm', 'f')
  UNION ALL
  SELECT 'função', p.oid::regprocedure::text, pg_get_userbyid(p.proowner), p.oid::regprocedure::text
  FROM pg_proc p JOIN pg_namespace n ON n.oid = p.pronamespace
  WHERE n.nspname IN ('app', 'plataforma')
)
SELECT o.tipo || ' ' || o.nome || ' de ' || o.dono
FROM objetos o
WHERE o.dono <> CASE WHEN o.base IN (SELECT nome FROM credencial) THEN 'tenaz_credencial' ELSE 'tenaz_dono' END
UNION ALL
SELECT 'objeto de credencial ausente: ' || c.nome FROM credencial c
WHERE c.nome NOT IN (SELECT nome FROM objetos WHERE tipo IN ('tabela', 'função'))
UNION ALL
SELECT 'esquema ' || n.nspname || ' de ' || pg_get_userbyid(n.nspowner)
FROM pg_namespace n
WHERE n.nspname IN ('app', 'plataforma')
  AND pg_get_userbyid(n.nspowner) <> 'tenaz_dono'
ORDER BY 1`

func TestTodoObjetoPertenceAoDono(t *testing.T) {
	ctx := context.Background()
	for _, o := range listarTexto(t, ctx, pool, sqlDonoErrado) {
		t.Errorf("%s; deveria pertencer a tenaz_dono: a migração que o criou esqueceu o SET ROLE tenaz_dono", o)
	}
}

func TestDetectorAcusaObjetoDeOutroDono(t *testing.T) {
	ctx := context.Background()
	// Criados pelo superusuário da suíte: o dono é ele, não tenaz_dono.
	tx := criarTabelasDeProva(t, ctx, []string{
		"CREATE TABLE app.tabela_proposital_dono (empresa_id uuid)",
		"CREATE FUNCTION plataforma.funcao_proposital_dono() RETURNS int LANGUAGE sql AS 'SELECT 1'",
		// A credencial de volta ao dono: o tenaz_manutencao, membro dele, a leria.
		"ALTER TABLE plataforma.sessao OWNER TO tenaz_dono",
		// E uma tabela comum dada ao dono das credenciais.
		"ALTER TABLE app.unidade OWNER TO tenaz_credencial",
	})
	o := strings.Join(listarTexto(t, ctx, tx, sqlDonoErrado), " | ")
	for _, esperado := range []string{"tabela app.tabela_proposital_dono", "função plataforma.funcao_proposital_dono()",
		"tabela plataforma.sessao de tenaz_dono", "tabela app.unidade de tenaz_credencial",
		// o ALTER ... OWNER leva junto índices e tipo de linha, que também são acusados
		"índice plataforma.pk_sessao de tenaz_dono", "tipo de linha plataforma.sessao de tenaz_dono"} {
		if !strings.Contains(o, esperado) {
			t.Errorf("não acusou %q; acusados: %s", esperado, o)
		}
	}
}

// ------------------------------------------------------ papéis

// Nenhum papel de execução contorna o RLS nem é superusuário — inclusive o
// tenaz_manutencao, que é membro do dono por desenho (mutação m7).
const sqlPapeisQueContornamRLS = `
SELECT r.nome || CASE WHEN p.oid IS NULL THEN ' não existe'
                      WHEN p.rolsuper THEN ' é superusuário'
                      ELSE ' tem BYPASSRLS' END
FROM unnest(` + sqlPapeisDeExecucaoTexto + `) r(nome)
LEFT JOIN pg_roles p ON p.rolname = r.nome
WHERE p.oid IS NULL OR p.rolsuper OR p.rolbypassrls
ORDER BY 1`

func TestNenhumPapelDaAplicacaoContornaRLS(t *testing.T) {
	ctx := context.Background()
	for _, p := range listarTexto(t, ctx, pool, sqlPapeisQueContornamRLS) {
		t.Errorf("papel de execução: %s; nenhum papel que atende requisição ou job contorna o RLS", p)
	}
}

func TestDetectorAcusaPapelQueContornaRLS(t *testing.T) {
	ctx := context.Background()
	tx := criarTabelasDeProva(t, ctx, []string{
		"ALTER ROLE tenaz_manutencao BYPASSRLS",
		"ALTER ROLE tenaz_leitura SUPERUSER",
	})
	p := listarTexto(t, ctx, tx, sqlPapeisQueContornamRLS)
	if !contem(p, "tenaz_manutencao tem BYPASSRLS") || !contem(p, "tenaz_leitura é superusuário") {
		t.Errorf("BYPASSRLS e superusuário em papel de execução não foram acusados: %v", p)
	}
}

// Relações de membro de e para os papéis de execução: só as declaradas em
// relacoes_de_membro.txt, com as mesmas opções (INHERIT, SET). Uma relação a
// mais herda privilégios silenciosamente: GRANT tenaz_dono TO tenaz_app torna o
// tenaz_app dono de tudo (m19); GRANT tenaz_leitura TO tenaz_sinc dá ao
// sincronizador leitura de toda tabela de app (m18).
//
// Fica de fora só a concessão pura de ADMIN (sem INHERIT nem SET), que o
// PostgreSQL 16+ dá a quem cria o papel quando ele não é superusuário: ela
// não transfere privilégio.
const sqlRelacoesDeMembro = `
SELECT m.rolname || ' ' || r.rolname || ' ' ||
       concat_ws(',', CASE WHEN am.inherit_option THEN 'INHERIT' END, CASE WHEN am.set_option THEN 'SET' END)
FROM pg_auth_members am
JOIN pg_roles r ON r.oid = am.roleid
JOIN pg_roles m ON m.oid = am.member
WHERE (m.rolname = ANY (` + sqlPapeisDeExecucaoTexto + `) OR r.rolname = ANY (` + sqlPapeisDeExecucaoTexto + `))
  AND (am.inherit_option OR am.set_option)
ORDER BY 1`

func diferencasDasRelacoes(t *testing.T, ctx context.Context, q consultor) (aMais, aMenos []string) {
	t.Helper()
	l, err := lerListaDe(arquivoRelacoesDeMembro, false)
	if err != nil {
		t.Fatalf("ler %s: %v", arquivoRelacoesDeMembro, err)
	}
	declarado := map[string]bool{}
	for _, e := range l {
		declarado[e.Nome+" "+strings.Join(e.Campos, " ")] = true
	}
	existentes := map[string]bool{}
	for _, rel := range listarTexto(t, ctx, q, sqlRelacoesDeMembro) {
		existentes[rel] = true
		if !declarado[rel] {
			aMais = append(aMais, rel)
		}
	}
	for rel := range declarado {
		if !existentes[rel] {
			aMenos = append(aMenos, rel)
		}
	}
	return aMais, aMenos
}

func TestRelacoesDeMembroSaoSoAsDeclaradas(t *testing.T) {
	ctx := context.Background()
	aMais, aMenos := diferencasDasRelacoes(t, ctx, pool)
	for _, r := range aMais {
		t.Errorf("relação de membro não declarada (membro papel opções): %s; um membro herda os privilégios "+
			"do papel, inclusive o de dono", r)
	}
	for _, r := range aMenos {
		t.Errorf("relação declarada em %s e ausente do banco: %s", arquivoRelacoesDeMembro, r)
	}
}

func TestDetectorAcusaRelacaoDeMembroNaoDeclarada(t *testing.T) {
	ctx := context.Background()
	tx := criarTabelasDeProva(t, ctx, []string{
		"GRANT tenaz_leitura TO tenaz_sinc",
		"GRANT tenaz_dono TO tenaz_app",
		"REVOKE tenaz_autenticacao FROM tenaz_autenticacao_conexao",
	})
	aMais, aMenos := diferencasDasRelacoes(t, ctx, tx)
	// tenaz_app é NOINHERIT: o GRANT sem opção explícita sai só com SET, e
	// ainda assim é acusado (SET ROLE tenaz_dono basta para ser o dono).
	for _, esperado := range []string{"tenaz_sinc tenaz_leitura INHERIT,SET", "tenaz_app tenaz_dono SET"} {
		if !contem(aMais, esperado) {
			t.Errorf("a relação %q não foi acusada; acusadas: %v", esperado, aMais)
		}
	}
	if !contem(aMenos, "tenaz_autenticacao_conexao tenaz_autenticacao INHERIT,SET") {
		t.Errorf("a relação declarada e retirada não foi acusada; acusadas: %v", aMenos)
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
