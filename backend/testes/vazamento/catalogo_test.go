package vazamento

import (
	"context"
	"sort"
	"strings"
	"testing"
)

// Verificações de catálogo que não dependem de dados: cada uma é uma regra
// que vale para TODA tabela, inclusive as que ainda não existem, e cada uma
// tem a sua prova de detector.

// ------------------------------------------------------ FK composta

// Toda FK que aponta para tabela de empresa leva empresa_id junto, na mesma
// posição (seção 2.5 do documento 2). A verificação de integridade referencial
// ignora o RLS: com FK simples, um id de outra empresa seria aceito. Vale
// também para as FKs de plataforma para app (sessão, código, convite, acesso).
const sqlFKSimples = `
WITH fks AS (
  SELECT con.conname, con.conrelid, con.confrelid,
         (SELECT array_agg(a.attname ORDER BY k.ord)
            FROM unnest(con.conkey) WITH ORDINALITY k(attnum, ord)
            JOIN pg_attribute a ON a.attrelid = con.conrelid AND a.attnum = k.attnum) AS colunas,
         (SELECT array_agg(a.attname ORDER BY k.ord)
            FROM unnest(con.confkey) WITH ORDINALITY k(attnum, ord)
            JOIN pg_attribute a ON a.attrelid = con.confrelid AND a.attnum = k.attnum) AS referenciadas
  FROM pg_constraint con
  JOIN pg_class c  ON c.oid = con.conrelid
  JOIN pg_namespace n  ON n.oid = c.relnamespace
  JOIN pg_class rc ON rc.oid = con.confrelid
  JOIN pg_namespace rn ON rn.oid = rc.relnamespace
  WHERE con.contype = 'f' AND con.conparentid = 0
    AND n.nspname IN ('app', 'plataforma') AND rn.nspname = 'app'
)
SELECT conrelid::regclass::text || '.' || conname || ' ' || array_to_string(colunas, ',') || ' -> ' ||
       confrelid::regclass::text || ' ' || array_to_string(referenciadas, ',')
FROM fks
WHERE NOT (
  -- para app.empresa: a própria coluna empresa_id aponta para o id
  (confrelid = 'app.empresa'::regclass AND colunas = ARRAY['empresa_id'::name] AND referenciadas = ARRAY['id'::name])
  -- para tabela de empresa: empresa_id casado com empresa_id
  OR EXISTS (SELECT 1 FROM generate_subscripts(colunas, 1) i
             WHERE colunas[i] = 'empresa_id' AND referenciadas[i] = 'empresa_id'))
ORDER BY 1`

func TestFKParaTabelaDeEmpresaEComposta(t *testing.T) {
	ctx := context.Background()
	for _, fk := range listarTexto(t, ctx, pool, sqlFKSimples) {
		t.Errorf("FK sem empresa_id casado: %s; a FK ignora o RLS e aceitaria linha de outra empresa (seção 2.5)", fk)
	}
}

func TestDetectorAcusaFKSimples(t *testing.T) {
	ctx := context.Background()
	const pai, filho = "app.tabela_proposital_pai", "app.tabela_proposital_filho_fk_simples"
	tx := criarTabelasDeProva(t, ctx, []string{
		"CREATE TABLE " + pai + " (empresa_id uuid NOT NULL, id uuid NOT NULL, PRIMARY KEY (empresa_id, id), UNIQUE (id))",
		"CREATE TABLE " + filho + " (empresa_id uuid NOT NULL, id uuid PRIMARY KEY, pai_id uuid REFERENCES " + pai + " (id))",
	})

	var acusou bool
	for _, fk := range listarTexto(t, ctx, tx, sqlFKSimples) {
		if strings.HasPrefix(fk, filho+".") {
			acusou = true
		}
	}
	if !acusou {
		t.Fatal("uma FK simples (pai_id -> id) entre tabelas de empresa não foi acusada")
	}
}

// ------------------------------------------------------ unicidade

// Unicidade de negócio é sempre composta com empresa_id na frente (consulta 7
// da seção 2.8): uma unicidade global responde "já existe" sobre dado de outra
// empresa. app.empresa é a exceção declarada (o próprio tenant, com CNPJ único).
const sqlUnicidadeGlobal = `
SELECT i.indexrelid::regclass::text
FROM pg_index i
JOIN pg_class c ON c.oid = i.indrelid
JOIN pg_namespace n ON n.oid = c.relnamespace
LEFT JOIN pg_attribute a ON a.attrelid = c.oid AND a.attnum = i.indkey[0]
WHERE n.nspname = 'app' AND i.indisunique AND c.oid <> 'app.empresa'::regclass
  AND coalesce(a.attname, '') <> 'empresa_id'
ORDER BY 1`

func TestUnicidadeComecaPorEmpresa(t *testing.T) {
	ctx := context.Background()
	for _, ix := range listarTexto(t, ctx, pool, sqlUnicidadeGlobal) {
		t.Errorf("índice único %s não começa por empresa_id; a violação revelaria dado de outra empresa", ix)
	}
}

// ------------------------------------------------------ escopo externo

var ordemComandos = []string{"SELECT", "INSERT", "UPDATE", "DELETE"}

// sqlPoliticasDoApp lista, por tabela de app, as policies que valem para o
// tenaz_app e os privilégios dele. É a matéria-prima de acessoExterno.
const sqlPoliticasDoApp = `
SELECT c.oid::regclass::text AS tabela,
       has_table_privilege('tenaz_app', c.oid, 'SELECT') AS p_select,
       has_table_privilege('tenaz_app', c.oid, 'INSERT') AS p_insert,
       has_table_privilege('tenaz_app', c.oid, 'UPDATE') AS p_update,
       has_table_privilege('tenaz_app', c.oid, 'DELETE') AS p_delete,
       coalesce(p.polname, '') AS politica,
       coalesce(p.polpermissive, true) AS permissiva,
       coalesce(p.polcmd::text, '') AS comando,
       coalesce(pg_get_expr(p.polqual, p.polrelid), '') AS usando,
       coalesce(pg_get_expr(p.polwithcheck, p.polrelid), '') AS verificando
FROM pg_class c
JOIN pg_namespace n ON n.oid = c.relnamespace
LEFT JOIN pg_policy p ON p.polrelid = c.oid
                     AND (p.polroles @> ARRAY['tenaz_app'::regrole::oid] OR p.polroles = '{0}')
WHERE n.nspname = 'app' AND c.relkind IN ('r', 'p') AND NOT c.relispartition
ORDER BY 1, 6`

type politica struct {
	permissiva          bool
	comando             string // '*', 'r', 'a', 'w', 'd'
	usando, verificando string
}

func (p politica) vale(cmd string) bool {
	letra := map[string]string{"SELECT": "r", "INSERT": "a", "UPDATE": "w", "DELETE": "d"}[cmd]
	return p.comando == "*" || p.comando == letra
}

// travaInterno diz se a restritiva prende o comando ao escopo interno. Numa
// policy FOR ALL sem WITH CHECK, o USING vale também para a linha nova.
func (p politica) travaInterno(cmd string) bool {
	verif := p.verificando
	if verif == "" && p.comando == "*" {
		verif = p.usando
	}
	switch cmd {
	case "SELECT", "DELETE":
		return p.usando == exprEscopoInterno
	case "INSERT":
		return verif == exprEscopoInterno
	default: // UPDATE: USING e WITH CHECK precisam passar
		return p.usando == exprEscopoInterno || verif == exprEscopoInterno
	}
}

// acessoExterno calcula, para cada tabela de app, os comandos que o tenaz_app
// pode executar sem que uma restritiva o prenda ao escopo interno.
func acessoExterno(t *testing.T, ctx context.Context, q consultor) map[string]string {
	t.Helper()
	linhas, err := q.Query(ctx, sqlPoliticasDoApp)
	if err != nil {
		t.Fatalf("consultar policies: %v", err)
	}
	defer linhas.Close()

	privilegio := map[string]map[string]bool{}
	politicas := map[string][]politica{}
	for linhas.Next() {
		var tabela, nome, cmd, usando, verif string
		var ps, pi, pu, pd, perm bool
		if err := linhas.Scan(&tabela, &ps, &pi, &pu, &pd, &nome, &perm, &cmd, &usando, &verif); err != nil {
			t.Fatalf("ler linha: %v", err)
		}
		privilegio[tabela] = map[string]bool{"SELECT": ps, "INSERT": pi, "UPDATE": pu, "DELETE": pd}
		if nome != "" {
			politicas[tabela] = append(politicas[tabela], politica{perm, cmd, usando, verif})
		}
	}
	if err := linhas.Err(); err != nil {
		t.Fatalf("iterar: %v", err)
	}

	saida := map[string]string{}
	for tabela, priv := range privilegio {
		var abertos []string
		for _, cmd := range ordemComandos {
			if !priv[cmd] {
				continue
			}
			var temPermissiva, travado bool
			for _, p := range politicas[tabela] {
				if !p.vale(cmd) {
					continue
				}
				if p.permissiva {
					temPermissiva = true
				} else if p.travaInterno(cmd) {
					travado = true
				}
			}
			if temPermissiva && !travado {
				abertos = append(abertos, cmd)
			}
		}
		if len(abertos) > 0 {
			saida[tabela] = strings.Join(abertos, ",")
		}
	}
	return saida
}

func TestListaDeAcessoExternoBateComAsPoliticas(t *testing.T) {
	ctx := context.Background()

	l, err := lerLista(arquivoAcessoExterno)
	if err != nil {
		t.Fatalf("ler %s: %v", arquivoAcessoExterno, err)
	}
	declarado := map[string]string{}
	for _, e := range l {
		if len(e.Campos) != 1 {
			t.Errorf("%s: esperado um campo de comandos, veio %v", e.Nome, e.Campos)
			continue
		}
		declarado[e.Nome] = e.Campos[0]
	}

	calculado := acessoExterno(t, ctx, pool)
	for tabela, cmds := range calculado {
		if declarado[tabela] != cmds {
			t.Errorf("%s: as policies abrem %s a escopo externo, e a lista declara %q; "+
				"ou a policy está larga demais (falta escopo_externo_exclusao?) ou a lista precisa de revisão de isolamento",
				tabela, cmds, declarado[tabela])
		}
	}
	for tabela, cmds := range declarado {
		if _, ok := calculado[tabela]; !ok {
			t.Errorf("%s está em %s com %s, mas o banco não a abre a escopo externo; retire a linha",
				tabela, arquivoAcessoExterno, cmds)
		}
	}
}

// Prova do detector, com o erro que a revisão de E01 confirmou: uma
// restritiva FOR ALL que libera a leitura a escopo externo libera também o
// DELETE, porque o WITH CHECK não alcança o DELETE.
func TestDetectorAcusaExclusaoPorEscopoExterno(t *testing.T) {
	ctx := context.Background()
	const nome = "app.tabela_proposital_leitura_externa"
	tx := criarTabelasDeProva(t, ctx, []string{
		"CREATE TABLE " + nome + " (empresa_id uuid NOT NULL, id uuid PRIMARY KEY)",
		"GRANT SELECT, INSERT, UPDATE, DELETE ON " + nome + " TO tenaz_app",
		"CREATE POLICY isolamento_empresa ON " + nome + " FOR ALL TO tenaz_app " +
			"USING (empresa_id = app.empresa_atual()) WITH CHECK (empresa_id = app.empresa_atual())",
		"CREATE POLICY escopo_externo ON " + nome + " AS RESTRICTIVE FOR ALL TO tenaz_app " +
			"USING (app.escopo_atual() IN ('interno', 'prestador', 'locatario', 'link')) WITH CHECK (app.escopo_interno())",
	})

	if got := acessoExterno(t, ctx, tx)[nome]; got != "SELECT,DELETE" {
		t.Fatalf("a tabela de prova deveria abrir SELECT e DELETE a escopo externo, o detector calculou %q", got)
	}
}

// ------------------------------------------------------ SECURITY DEFINER

const sqlSecurityDefiner = `
SELECT p.oid::regprocedure::text,
       coalesce(array_to_string(p.proconfig, '|'), ''),
       EXISTS (SELECT 1 FROM aclexplode(coalesce(p.proacl, acldefault('f', p.proowner))) x
               WHERE x.grantee = 0 AND x.privilege_type = 'EXECUTE'),
       array_to_string(ARRAY(
         SELECT r FROM unnest(` + sqlPapeisDaAplicacao + `) r
         WHERE has_function_privilege(r, p.oid, 'EXECUTE') ORDER BY r), ',')
FROM pg_proc p
JOIN pg_namespace n ON n.oid = p.pronamespace
WHERE n.nspname IN ('app', 'plataforma') AND p.prosecdef
ORDER BY 1`

type funcaoDefiner struct {
	Assinatura, Config, Executores string
	Publica                        bool
}

func funcoesDefiner(t *testing.T, ctx context.Context, q consultor) []funcaoDefiner {
	t.Helper()
	linhas, err := q.Query(ctx, sqlSecurityDefiner)
	if err != nil {
		t.Fatalf("consultar funções SECURITY DEFINER: %v", err)
	}
	defer linhas.Close()
	var saida []funcaoDefiner
	for linhas.Next() {
		var f funcaoDefiner
		if err := linhas.Scan(&f.Assinatura, &f.Config, &f.Publica, &f.Executores); err != nil {
			t.Fatalf("ler linha: %v", err)
		}
		saida = append(saida, f)
	}
	if err := linhas.Err(); err != nil {
		t.Fatalf("iterar: %v", err)
	}
	return saida
}

// problemasDefiner confere uma função contra a lista aprovada.
func problemasDefiner(f funcaoDefiner, aprovadas map[string]string) []string {
	var p []string
	executores, ok := aprovadas[f.Assinatura]
	if !ok {
		p = append(p, "fora da lista aprovada em "+arquivoSecurityDefiner)
	}
	if !contem(strings.Split(f.Config, "|"), "search_path=pg_catalog, pg_temp") {
		p = append(p, "search_path não é exatamente \"pg_catalog, pg_temp\" (config: "+f.Config+")")
	}
	if f.Publica {
		p = append(p, "EXECUTE concedido a PUBLIC")
	}
	if executores == "-" {
		executores = ""
	}
	if ok && f.Executores != executores {
		p = append(p, "executam: "+f.Executores+"; aprovado: "+executores)
	}
	return p
}

func TestFuncoesSecurityDefinerAprovadas(t *testing.T) {
	ctx := context.Background()
	l, err := lerLista(arquivoSecurityDefiner)
	if err != nil {
		t.Fatalf("ler %s: %v", arquivoSecurityDefiner, err)
	}
	aprovadas := map[string]string{}
	for _, e := range l {
		if len(e.Campos) != 1 {
			t.Errorf("%s: esperado o campo de executores, veio %v", e.Nome, e.Campos)
			continue
		}
		aprovadas[e.Nome] = e.Campos[0]
	}
	existentes := map[string]bool{}
	for _, f := range funcoesDefiner(t, ctx, pool) {
		existentes[f.Assinatura] = true
		for _, p := range problemasDefiner(f, aprovadas) {
			t.Errorf("SECURITY DEFINER %s: %s", f.Assinatura, p)
		}
	}
	for nome := range aprovadas {
		if !existentes[nome] {
			t.Errorf("%s está em %s mas não existe (ou deixou de ser SECURITY DEFINER); retire a linha",
				nome, arquivoSecurityDefiner)
		}
	}
}

func TestDetectorAcusaSecurityDefinerSemSearchPath(t *testing.T) {
	ctx := context.Background()
	const nome = "plataforma.funcao_proposital_definer()"
	tx := criarTabelasDeProva(t, ctx, []string{
		"CREATE FUNCTION " + nome + " RETURNS int LANGUAGE sql SECURITY DEFINER AS 'SELECT 1'",
		// E uma aprovada que passa a ser executada por quem não devia.
		"GRANT EXECUTE ON FUNCTION plataforma.purgar_credenciais_vencidas() TO tenaz_app",
	})
	aprovadas := map[string]string{"plataforma.purgar_credenciais_vencidas()": "tenaz_job"}
	var achouProva, achouPurga bool
	for _, f := range funcoesDefiner(t, ctx, tx) {
		p := strings.Join(problemasDefiner(f, aprovadas), " | ")
		switch f.Assinatura {
		case nome:
			achouProva = true
			for _, esperado := range []string{"fora da lista aprovada", "search_path", "EXECUTE concedido a PUBLIC"} {
				if !strings.Contains(p, esperado) {
					t.Errorf("a função de prova não foi acusada por %q; problemas: %s", esperado, p)
				}
			}
		case "plataforma.purgar_credenciais_vencidas()":
			achouPurga = true
			if !strings.Contains(p, "executam: tenaz_app,tenaz_job") {
				t.Errorf("EXECUTE da purga concedido ao tenaz_app não foi acusado; problemas: %q", p)
			}
		}
	}
	if !achouProva || !achouPurga {
		t.Fatal("as funções da prova não apareceram na consulta do catálogo")
	}
}

// ------------------------------------------------------ privilégios

// Mapa fechado de privilégios onde o RLS não é a defesa. Entram:
//   - tabela, visão, sequência e função FORA de app, em qualquer esquema que
//     não seja do sistema, concedidas a qualquer papel tenaz_* ou a PUBLIC
//     (em app, o privilégio padrão do tenaz_app, do tenaz_job e do
//     tenaz_leitura é governado pelo RLS e pelas outras regras da suíte);
//   - qualquer objeto, em qualquer esquema, concedido ao sincronizador, ao
//     papel da autenticação ou ao dono das credenciais (o privilégio exato é a
//     defesa deles: R5, ADR-011, ISO-E01-08);
//   - todo privilégio de esquema (USAGE, CREATE) dado a papel de execução ou a
//     PUBLIC, como "esquema:<nome>".
//
// Um esquema novo com uma cópia de pessoa e SELECT para o tenaz_app aparece
// aqui duas vezes: o USAGE no esquema e o SELECT na tabela.
const sqlPrivilegios = `
WITH papeis AS (
  SELECT oid, rolname::text AS rolname,
         rolname IN ('tenaz_sinc', 'tenaz_sinc_conexao', 'tenaz_autenticacao', 'tenaz_autenticacao_conexao',
                     'tenaz_credencial') AS fechado
  FROM pg_roles WHERE rolname LIKE 'tenaz\_%'
  UNION ALL SELECT 0, 'PUBLIC', false
), esquemas AS (
  SELECT oid, nspname, nspacl, nspowner FROM pg_namespace
  WHERE nspname NOT IN ('pg_catalog', 'information_schema') AND nspname NOT LIKE 'pg\_%'
), objetos AS (
  SELECT e.nspname AS esquema, c.oid::regclass::text AS nome, c.relacl AS acl, c.relowner AS dono, c.oid AS rel,
         false AS funcao
  FROM pg_class c JOIN esquemas e ON e.oid = c.relnamespace
  WHERE c.relkind IN ('r', 'p', 'v', 'm', 'f', 'S') AND NOT c.relispartition
  UNION ALL
  SELECT e.nspname, p.oid::regprocedure::text, p.proacl, p.proowner, NULL::oid, true
  FROM pg_proc p JOIN esquemas e ON e.oid = p.pronamespace
), privs AS (
  SELECT o.esquema, o.nome, o.funcao, x.grantee, x.privilege_type, NULL::text AS colunas
  FROM objetos o, aclexplode(o.acl) x WHERE x.grantee <> o.dono
  UNION ALL
  SELECT o.esquema, o.nome, o.funcao, x.grantee, x.privilege_type, string_agg(att.attname, ',' ORDER BY att.attname)
  FROM objetos o
  JOIN pg_attribute att ON att.attrelid = o.rel AND att.attnum > 0 AND NOT att.attisdropped,
       aclexplode(att.attacl) x
  WHERE x.grantee <> o.dono
  GROUP BY o.esquema, o.nome, o.funcao, x.grantee, x.privilege_type
)
SELECT pr.nome || ' ' || pa.rolname || ' ' || pr.privilege_type || coalesce('(' || pr.colunas || ')', '')
FROM privs pr JOIN papeis pa ON pa.oid = pr.grantee
-- PUBLIC tem EXECUTE em toda função por padrão do PostgreSQL; nas SECURITY
-- DEFINER, a lista aprovada exige a revogação (TestFuncoesSecurityDefinerAprovadas)
WHERE NOT (pa.oid = 0 AND pr.funcao)
  AND (pr.esquema <> 'app' OR pa.fechado OR pa.oid = 0)
UNION ALL
SELECT 'esquema:' || e.nspname || ' ' || pa.rolname || ' ' || x.privilege_type
FROM esquemas e
CROSS JOIN LATERAL aclexplode(e.nspacl) x
JOIN papeis pa ON pa.oid = x.grantee
WHERE x.grantee <> e.nspowner
ORDER BY 1`

// diferencasDoMapa compara o catálogo com privilegios.txt nos dois sentidos.
func diferencasDoMapa(t *testing.T, ctx context.Context, q consultor) (aMais, aMenos []string) {
	t.Helper()
	l, err := lerLista(arquivoPrivilegios)
	if err != nil {
		t.Fatalf("ler %s: %v", arquivoPrivilegios, err)
	}
	declarado := map[string]bool{}
	for _, e := range l {
		declarado[e.Nome+" "+strings.Join(e.Campos, " ")] = true
	}
	existentes := map[string]bool{}
	for _, p := range listarTexto(t, ctx, q, sqlPrivilegios) {
		existentes[p] = true
		if !declarado[p] {
			aMais = append(aMais, p)
		}
	}
	for p := range declarado {
		if !existentes[p] {
			aMenos = append(aMenos, p)
		}
	}
	sort.Strings(aMenos)
	return aMais, aMenos
}

func TestPrivilegiosForaDoRLSBatemComOMapa(t *testing.T) {
	ctx := context.Background()
	aMais, aMenos := diferencasDoMapa(t, ctx, pool)
	for _, p := range aMais {
		t.Errorf("privilégio fora do mapa: %s; fora de app o privilégio é a defesa, e no sincronizador e na "+
			"autenticação é a lista fechada. Revisão de isolamento antes de acrescentar a %s", p, arquivoPrivilegios)
	}
	for _, p := range aMenos {
		t.Errorf("privilégio declarado em %s e ausente do banco: %s", arquivoPrivilegios, p)
	}
}

// A mutação m24 (esquema novo com cópia de pessoa) e uma concessão a mais
// numa credencial precisam ser acusadas; e uma linha do mapa que some do banco
// também.
func TestDetectorAcusaPrivilegioForaDoMapa(t *testing.T) {
	ctx := context.Background()
	tx := criarTabelasDeProva(t, ctx, []string{
		"CREATE SCHEMA rel_prova",
		"CREATE TABLE rel_prova.pessoa AS SELECT * FROM app.pessoa",
		"GRANT USAGE ON SCHEMA rel_prova TO tenaz_app",
		"GRANT SELECT ON rel_prova.pessoa TO tenaz_app",
		"GRANT SELECT ON plataforma.sessao TO tenaz_leitura",
		"REVOKE UPDATE (revogada_em) ON plataforma.sessao FROM tenaz_autenticacao",
	})
	aMais, aMenos := diferencasDoMapa(t, ctx, tx)
	for _, esperado := range []string{"esquema:rel_prova tenaz_app USAGE", "rel_prova.pessoa tenaz_app SELECT",
		"plataforma.sessao tenaz_leitura SELECT"} {
		if !contem(aMais, esperado) {
			t.Errorf("o mapa não acusou %q; a mais: %v", esperado, aMais)
		}
	}
	if len(aMenos) == 0 {
		t.Error("o mapa não acusou o privilégio de coluna retirado do tenaz_autenticacao")
	}
}

// Consulta 5 da seção 2.8: nenhuma visão materializada, e toda visão com
// security_invoker, em qualquer esquema que não seja do sistema. Uma visão
// comum roda com os privilégios do DONO: sobre plataforma.identidade, ela
// entregaria senha_hash a quem tivesse SELECT nela (mutação m6).
const sqlVisoesPerigosas = `
SELECT c.oid::regclass::text || CASE c.relkind WHEN 'm' THEN ' (materializada)' ELSE ' (sem security_invoker)' END
FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
WHERE n.nspname NOT IN ('pg_catalog', 'information_schema') AND n.nspname NOT LIKE 'pg\_%'
  AND (c.relkind = 'm'
       OR (c.relkind = 'v' AND NOT coalesce(c.reloptions @> ARRAY['security_invoker=true'], false)))
ORDER BY 1`

func TestVisoesComSecurityInvokerESemMaterializada(t *testing.T) {
	ctx := context.Background()
	for _, v := range listarTexto(t, ctx, pool, sqlVisoesPerigosas) {
		t.Errorf("visão proibida: %s; use WITH (security_invoker = true), e nunca visão materializada", v)
	}
}

func TestDetectorAcusaVisaoDoDono(t *testing.T) {
	ctx := context.Background()
	tx := criarTabelasDeProva(t, ctx, []string{
		"CREATE VIEW app.v_identidade_prova AS SELECT id, email, senha_hash FROM plataforma.identidade",
		"CREATE MATERIALIZED VIEW plataforma.mv_prova AS SELECT 1 AS x",
		"CREATE VIEW app.v_invocador_prova WITH (security_invoker = true) AS SELECT 1 AS x",
	})
	v := listarTexto(t, ctx, tx, sqlVisoesPerigosas)
	if !contem(v, "app.v_identidade_prova (sem security_invoker)") || !contem(v, "plataforma.mv_prova (materializada)") {
		t.Errorf("as visões de prova não foram acusadas: %v", v)
	}
	if contem(v, "app.v_invocador_prova (sem security_invoker)") {
		t.Error("a visão com security_invoker foi acusada sem motivo")
	}
}

// TRUNCATE ignora o RLS; TRIGGER deixaria o papel pendurar código numa tabela
// alheia; REFERENCES deixaria sondar existência por FK. Nenhum papel da
// aplicação recebe nenhum dos três, em tabela nenhuma de esquema nenhum.
const sqlPrivilegiosProibidos = `
SELECT c.oid::regclass::text || ' ' || r || ' ' || pr
FROM pg_class c
JOIN pg_namespace n ON n.oid = c.relnamespace,
     unnest(` + sqlPapeisDaAplicacao + `) r,
     unnest(ARRAY['TRUNCATE', 'TRIGGER', 'REFERENCES']) pr
WHERE n.nspname NOT IN ('pg_catalog', 'information_schema') AND n.nspname NOT LIKE 'pg\_%'
  AND c.relkind IN ('r', 'p')
  AND has_table_privilege(r, c.oid, pr)
ORDER BY 1`

func TestNenhumPapelDaAplicacaoTemTruncateTriggerOuReferences(t *testing.T) {
	ctx := context.Background()
	for _, p := range listarTexto(t, ctx, pool, sqlPrivilegiosProibidos) {
		t.Errorf("privilégio proibido: %s", p)
	}
}

func TestDetectorAcusaTruncateTriggerReferences(t *testing.T) {
	ctx := context.Background()
	tx := criarTabelasDeProva(t, ctx, []string{
		"GRANT TRUNCATE ON app.unidade TO tenaz_app",
		"GRANT TRIGGER ON app.pessoa TO tenaz_job",
		"GRANT REFERENCES ON plataforma.identidade TO tenaz_autenticacao",
	})
	p := listarTexto(t, ctx, tx, sqlPrivilegiosProibidos)
	for _, esperado := range []string{"app.unidade tenaz_app TRUNCATE", "app.pessoa tenaz_job TRIGGER",
		"plataforma.identidade tenaz_autenticacao REFERENCES"} {
		if !contem(p, esperado) {
			t.Errorf("não acusou %q; acusou %v", esperado, p)
		}
	}
}

// Nenhum papel da aplicação cria objeto em esquema nenhum: um objeto criado
// por ele seria dele, e o dono ignora o RLS da própria tabela.
const sqlCreateEmEsquema = `
SELECT r || ' em ' || n.nspname
FROM unnest(` + sqlPapeisDaAplicacao + `) r, pg_namespace n
WHERE n.nspname NOT IN ('pg_catalog', 'information_schema') AND n.nspname NOT LIKE 'pg\_%'
  AND has_schema_privilege(r, n.oid, 'CREATE')
ORDER BY 1`

func TestNenhumPapelDaAplicacaoCriaObjeto(t *testing.T) {
	ctx := context.Background()
	for _, p := range listarTexto(t, ctx, pool, sqlCreateEmEsquema) {
		t.Errorf("CREATE concedido: %s", p)
	}
}

func TestDetectorAcusaCreateEmEsquema(t *testing.T) {
	ctx := context.Background()
	tx := criarTabelasDeProva(t, ctx, []string{
		"GRANT CREATE ON SCHEMA app TO tenaz_job",
		"CREATE SCHEMA rel_prova",
		"GRANT CREATE ON SCHEMA rel_prova TO tenaz_app",
	})
	p := listarTexto(t, ctx, tx, sqlCreateEmEsquema)
	if !contem(p, "tenaz_job em app") || !contem(p, "tenaz_app em rel_prova") {
		t.Errorf("CREATE concedido não foi acusado: %v", p)
	}
}

func TestDetectorAcusaUnicidadeGlobal(t *testing.T) {
	ctx := context.Background()
	tx := criarTabelasDeProva(t, ctx, []string{
		"CREATE TABLE app.tabela_proposital_unica (empresa_id uuid NOT NULL, placa text, PRIMARY KEY (empresa_id, placa))",
		"CREATE UNIQUE INDEX uq_prova_placa_global ON app.tabela_proposital_unica (placa)",
	})
	if !contem(listarTexto(t, ctx, tx, sqlUnicidadeGlobal), "app.uq_prova_placa_global") {
		t.Error("um índice único sem empresa_id na frente não foi acusado")
	}
}
