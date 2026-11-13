package vazamento

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
)

// Isolamento entre empresas (R2, RNF-013), orientado a tabela.
//
// A lista de tabelas vem do catálogo, não deste arquivo: toda tabela de app
// que existir entra sozinha, e o teste FALHA se alguma não tiver linha de A e
// de B nas fixtures — é o que obriga a manter o gerador junto com as
// migrações (parte 2 da seção 2.8 do documento 2).

type tabelaDeEmpresa struct {
	Nome   string // esquema.tabela
	Coluna string // coluna da empresa: empresa_id, ou id em app.empresa
	TemID  bool
}

func tabelasDeEmpresa(t *testing.T, ctx context.Context) []tabelaDeEmpresa {
	t.Helper()
	linhas, err := pool.Query(ctx, `
		SELECT c.oid::regclass::text,
		       EXISTS (SELECT 1 FROM pg_attribute a WHERE a.attrelid = c.oid AND a.attname = 'empresa_id'
		                 AND a.attnum > 0 AND NOT a.attisdropped),
		       EXISTS (SELECT 1 FROM pg_attribute a WHERE a.attrelid = c.oid AND a.attname = 'id'
		                 AND a.attnum > 0 AND NOT a.attisdropped)
		FROM pg_class c
		JOIN pg_namespace n ON n.oid = c.relnamespace
		WHERE n.nspname = 'app' AND c.relkind IN ('r', 'p') AND NOT c.relispartition
		ORDER BY 1`)
	if err != nil {
		t.Fatalf("listar tabelas de app: %v", err)
	}
	defer linhas.Close()

	var saida []tabelaDeEmpresa
	for linhas.Next() {
		var nome string
		var temEmpresa, temID bool
		if err := linhas.Scan(&nome, &temEmpresa, &temID); err != nil {
			t.Fatalf("ler linha: %v", err)
		}
		switch {
		case temEmpresa:
			saida = append(saida, tabelaDeEmpresa{nome, "empresa_id", temID})
		case nome == "app.empresa":
			saida = append(saida, tabelaDeEmpresa{nome, "id", true})
		default:
			t.Errorf("%s está em app sem empresa_id e não é app.empresa: a suíte não sabe isolá-la", nome)
		}
	}
	if err := linhas.Err(); err != nil {
		t.Fatalf("iterar: %v", err)
	}
	return saida
}

// amostra é o que o superusuário (fora do RLS) sabe de uma tabela: quantas
// linhas cada empresa tem, um id de B e uma linha de B inteira, em JSON.
type amostra struct {
	LinhasA, LinhasB int64
	IDB              string
	LinhaB           []byte
	LinhaA           []byte
}

func amostrar(t *testing.T, ctx context.Context, tb tabelaDeEmpresa, f *fixturesE01) amostra {
	t.Helper()
	var a amostra
	col := pgx.Identifier{tb.Coluna}.Sanitize()
	q := func(sql string, destino any, empresa string) {
		t.Helper()
		if err := pool.QueryRow(ctx, sql, empresa).Scan(destino); err != nil {
			t.Fatalf("amostrar %s: %v", tb.Nome, err)
		}
	}
	q("SELECT count(*) FROM "+tb.Nome+" WHERE "+col+" = $1", &a.LinhasA, f.A.ID)
	q("SELECT count(*) FROM "+tb.Nome+" WHERE "+col+" = $1", &a.LinhasB, f.B.ID)
	if a.LinhasA == 0 || a.LinhasB == 0 {
		return a
	}
	if tb.TemID {
		q("SELECT id::text FROM "+tb.Nome+" WHERE "+col+" = $1 LIMIT 1", &a.IDB, f.B.ID)
	}
	q("SELECT to_jsonb(x) FROM "+tb.Nome+" x WHERE "+col+" = $1 LIMIT 1", &a.LinhaB, f.B.ID)
	q("SELECT to_jsonb(x) FROM "+tb.Nome+" x WHERE "+col+" = $1 LIMIT 1", &a.LinhaA, f.A.ID)
	return a
}

func TestFixturesPovoamTodaTabela(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	for _, tb := range tabelasDeEmpresa(t, ctx) {
		a := amostrar(t, ctx, tb, f)
		if a.LinhasA == 0 || a.LinhasB == 0 {
			t.Errorf("%s sem linha nas fixtures (A=%d, B=%d); acrescente-a em fixtures_test.go — "+
				"tabela vazia passa em todo teste de vazamento sem provar nada", tb.Nome, a.LinhasA, a.LinhasB)
		}
	}
}

// Para cada tabela e cada papel da aplicação: sem contexto não há linha; no
// contexto de A, só as de A, nem por id conhecido de B; e nenhuma escrita
// alcança B.
func TestEmpresaANaoAlcancaB(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)

	for _, tb := range tabelasDeEmpresa(t, ctx) {
		a := amostrar(t, ctx, tb, f)
		if a.LinhasA == 0 || a.LinhasB == 0 {
			continue // TestFixturesPovoamTodaTabela acusa
		}
		col := pgx.Identifier{tb.Coluna}.Sanitize()

		for _, papel := range []string{papelApp, papelJob, papelLeitura} {
			t.Run(tb.Nome+"/"+papel, func(t *testing.T) {
				leituraSemContexto(t, ctx, tb, papel)
				leituraNoContextoDeA(t, ctx, tb, papel, col, a, f)
				escritaCruzada(t, ctx, tb, papel, col, a, f)
			})
		}
	}
}

func leituraSemContexto(t *testing.T, ctx context.Context, tb tabelaDeEmpresa, papel string) {
	t.Helper()
	// A listagem de empresas do job é sem contexto POR DESENHO (seção 2.7):
	// é dela que o job tira as empresas que vai percorrer.
	if papel == papelJob && tb.Nome == "app.empresa" {
		return
	}
	err := comoPapel(ctx, papel, contexto{}, func(tx pgx.Tx) error {
		n, err := contar(ctx, tx, "SELECT count(*) FROM "+tb.Nome)
		if err != nil {
			return err
		}
		if n != 0 {
			t.Errorf("sem contexto, %s viu %d linha(s) de %s; a policy não falha fechada", papel, n, tb.Nome)
		}
		return nil
	})
	if err != nil {
		t.Errorf("ler %s sem contexto: %v", tb.Nome, err)
	}
}

func leituraNoContextoDeA(t *testing.T, ctx context.Context, tb tabelaDeEmpresa, papel, col string, a amostra, f *fixturesE01) {
	t.Helper()
	var deB, peloID, total int64
	err := comoPapel(ctx, papel, f.internoA(), func(tx pgx.Tx) error {
		var err error
		if deB, err = contar(ctx, tx, "SELECT count(*) FROM "+tb.Nome+" WHERE "+col+" = $1", f.B.ID); err != nil {
			return err
		}
		if a.IDB != "" {
			// Sem filtro de empresa: é o RLS, e não o WHERE, que precisa cortar.
			if peloID, err = contar(ctx, tx, "SELECT count(*) FROM "+tb.Nome+" WHERE id = $1", a.IDB); err != nil {
				return err
			}
		}
		total, err = contar(ctx, tx, "SELECT count(*) FROM "+tb.Nome)
		return err
	})
	if err != nil {
		t.Errorf("ler %s no contexto de A: %v", tb.Nome, err)
		return
	}

	if papel == papelJob && tb.Nome == "app.empresa" {
		t.Run("job_lista_empresas_no_contexto", func(t *testing.T) {
			confirmarAchado(t, "ISO-E01-04", deB != 0 || peloID != 0 || total != a.LinhasA,
				fmt.Sprintf("tenaz_job no contexto de A viu %d empresa(s), entre elas B (%d)", total, deB))
		})
		return
	}
	if deB != 0 {
		t.Errorf("no contexto de A, %s leu %d linha(s) de B em %s", papel, deB, tb.Nome)
	}
	if peloID != 0 {
		t.Errorf("no contexto de A, %s achou em %s a linha de B pelo id", papel, tb.Nome)
	}
	if total != a.LinhasA {
		t.Errorf("no contexto de A, %s viu %d linha(s) em %s; A tem %d — "+
			"a diferença é linha de outra empresa ou policy que corta demais", papel, total, tb.Nome, a.LinhasA)
	}
}

func escritaCruzada(t *testing.T, ctx context.Context, tb tabelaDeEmpresa, papel, col string, a amostra, f *fixturesE01) {
	t.Helper()

	// Cada comando numa transação própria: um erro aborta a transação, e o
	// comando seguinte precisa ser avaliado por si.
	tentar := func(sql string, args ...any) (int64, error) {
		var linhas int64
		err := comoPapel(ctx, papel, f.internoA(), func(tx pgx.Tx) error {
			var err error
			linhas, err = alterou(ctx, tx, sql, args...)
			return err
		})
		return linhas, err
	}

	n, err := tentar("UPDATE "+tb.Nome+" SET "+col+" = "+col+" WHERE "+col+" = $1", f.B.ID)
	naoMudouNada(t, papel+": UPDATE de B em "+tb.Nome, n, err)

	n, err = tentar("DELETE FROM "+tb.Nome+" WHERE "+col+" = $1", f.B.ID)
	naoMudouNada(t, papel+": DELETE de B em "+tb.Nome, n, err)

	if a.IDB != "" {
		n, err = tentar("UPDATE "+tb.Nome+" SET "+col+" = "+col+" WHERE id = $1", a.IDB)
		naoMudouNada(t, papel+": UPDATE da linha de B pelo id em "+tb.Nome, n, err)

		n, err = tentar("DELETE FROM "+tb.Nome+" WHERE id = $1", a.IDB)
		naoMudouNada(t, papel+": DELETE da linha de B pelo id em "+tb.Nome, n, err)
	}

	// Cópia exata de uma linha de B, inserida no contexto de A: o WITH CHECK
	// recusa antes de qualquer unicidade (que, se viesse primeiro, confirmaria
	// a existência da linha).
	copia := "INSERT INTO " + tb.Nome + " SELECT * FROM jsonb_populate_record(NULL::" + tb.Nome + ", $1::jsonb)"
	_, err = tentar(copia, string(a.LinhaB))
	// Quem tem INSERT na tabela precisa ser barrado pela POLICY (WITH CHECK da
	// empresa); quem não tem, pelo privilégio. Um 42501 de outra origem não
	// prova o WITH CHECK.
	rotulo := papel + ": INSERT de linha com a empresa B em " + tb.Nome
	if !temInsert(t, ctx, papel, tb.Nome) {
		exigirErro(t, rotulo, err, erroPrivilegio)
		return
	}
	if cols := colunasSemInsert(t, ctx, papel, tb.Nome); len(cols) > 0 {
		// INSERT por coluna (evento_auditoria desde a 00008: sem id,
		// ocorrido_em, transacao e hash). A cópia inteira tem de ser barrada
		// pelo PRIVILÉGIO, e a policy é provada com a cópia só das colunas que
		// o papel pode mandar.
		if err == nil || recusaDeRLS(err) || sqlstate(err) != erroPrivilegio {
			t.Errorf("%s (todas as colunas, inclusive %v, fora do GRANT): esperado privilégio negado, veio %v", rotulo, cols, err)
		}
		lista := colunasComInsert(t, ctx, papel, tb.Nome)
		copia = "INSERT INTO " + tb.Nome + " (" + lista + ") SELECT " + lista +
			" FROM jsonb_populate_record(NULL::" + tb.Nome + ", $1::jsonb)"
		_, err = tentar(copia, string(a.LinhaB))
		rotulo += " (só as colunas do GRANT)"
	}
	if !recusaDeRLS(err) && err != nil && temGatilhoAntesDaInsercao(t, ctx, tb.Nome) {
		// Um trigger BEFORE INSERT (o do selo, que recalcula a linha) recusou
		// antes do WITH CHECK, que o PostgreSQL avalia sobre a linha que sai dos
		// triggers. A recusa veio, mas não prova a policy: a suíte repete a
		// inserção com os triggers de usuário desligados, só nesta transação.
		err = inserirSemGatilhos(t, ctx, papel, f.internoA(), tb.Nome, copia, string(a.LinhaB))
		rotulo += " (com os triggers de usuário desligados)"
	}
	exigirRLS(t, rotulo, err)
}

// colunasSemInsert lista as colunas em que o papel não tem INSERT (vazio
// quando o INSERT é da tabela inteira ou de todas as colunas).
func colunasSemInsert(t *testing.T, ctx context.Context, papel, tabela string) []string {
	t.Helper()
	return listarTexto(t, ctx, pool, `SELECT attname FROM pg_attribute WHERE attrelid = $2::regclass AND attnum > 0
	                                  AND NOT attisdropped AND NOT has_column_privilege($1::name, $2::regclass, attname, 'INSERT')
	                                  ORDER BY attnum`, papel, tabela)
}

// colunasComInsert devolve, prontas para SQL, as colunas em que o papel tem
// INSERT.
func colunasComInsert(t *testing.T, ctx context.Context, papel, tabela string) string {
	t.Helper()
	l := listarTexto(t, ctx, pool, `SELECT quote_ident(attname) FROM pg_attribute WHERE attrelid = $2::regclass AND attnum > 0
	                                  AND NOT attisdropped AND has_column_privilege($1::name, $2::regclass, attname, 'INSERT')
	                                  ORDER BY attnum`, papel, tabela)
	return strings.Join(l, ", ")
}

// temGatilhoAntesDaInsercao: a tabela tem trigger de usuário BEFORE INSERT FOR
// EACH ROW.
func temGatilhoAntesDaInsercao(t *testing.T, ctx context.Context, tabela string) bool {
	t.Helper()
	var tem bool
	if err := pool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM pg_trigger WHERE tgrelid = $1::regclass AND NOT tgisinternal
	                                  AND tgtype & 1 <> 0 AND tgtype & 2 <> 0 AND tgtype & 4 <> 0)`, tabela).Scan(&tem); err != nil {
		t.Fatalf("conferir triggers de %s: %v", tabela, err)
	}
	return tem
}

// inserirSemGatilhos executa a inserção como o papel, no contexto, com os
// triggers de usuário da tabela desligados numa transação desfeita. Serve só
// para observar a policy que um trigger anterior esconde; as chaves
// estrangeiras (triggers internos) continuam valendo.
func inserirSemGatilhos(t *testing.T, ctx context.Context, papel string, c contexto, tabela, sql string, args ...any) error {
	t.Helper()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("abrir transação: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	estragar(t, ctx, tx, "ALTER TABLE "+tabela+" DISABLE TRIGGER USER")
	if err := prepararTransacao(ctx, tx, papel, c); err != nil {
		return err
	}
	_, err = tx.Exec(ctx, sql, args...)
	return err
}

// O contexto não sobrevive ao fim da transação na mesma conexão física: é a
// conexão que volta ao pool e atende a requisição seguinte.
func TestContextoNaoVazaPeloPool(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)

	conexao, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("obter conexão: %v", err)
	}
	defer conexao.Release()

	tx, err := conexao.Begin(ctx)
	if err != nil {
		t.Fatalf("abrir transação: %v", err)
	}
	for _, v := range f.internoA().variaveis() {
		if _, err := tx.Exec(ctx, `SELECT set_config($1, $2, true)`, v[0], v[1]); err != nil {
			t.Fatalf("definir contexto: %v", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatalf("encerrar transação: %v", err)
	}

	tx, err = conexao.Begin(ctx)
	if err != nil {
		t.Fatalf("abrir segunda transação: %v", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()
	if _, err := tx.Exec(ctx, "SET LOCAL ROLE tenaz_app"); err != nil {
		t.Fatalf("assumir tenaz_app: %v", err)
	}
	for _, tb := range tabelasDeEmpresa(t, ctx) {
		n, err := contar(ctx, tx, "SELECT count(*) FROM "+tb.Nome)
		if err != nil {
			t.Fatalf("contar %s: %v", tb.Nome, err)
		}
		if n != 0 {
			t.Errorf("a conexão reaproveitada viu %d linha(s) de %s: o contexto da transação anterior vazou", n, tb.Nome)
		}
	}
}

// FK composta na prática: no contexto de A, uma linha de A apontando para
// linha de B é recusada pela FK, mesmo com todos os ids corretos.
func TestFKCompostaRecusaReferenciaABNaPratica(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)

	casos := []struct {
		rotulo, sql string
		args        []any
	}{
		{"turno em unidade de B", `INSERT INTO app.turno_unidade (empresa_id, unidade_id, dia_semana, hora_inicio, hora_fim)
		                            VALUES ($1, $2, 2, '08:00', '12:00')`, []any{f.A.ID, f.B.Unidade}},
		{"feriado em unidade de B", `INSERT INTO app.feriado_unidade (empresa_id, unidade_id, data_feriado, descricao)
		                              VALUES ($1, $2, '2026-01-01', 'Ano novo')`, []any{f.A.ID, f.B.Unidade}},
		{"perfil de B para usuário de A", `INSERT INTO app.usuario_perfil (empresa_id, usuario_id, perfil_id) VALUES ($1, $2, $3)`,
			[]any{f.A.ID, f.A.UsuarioInterno, f.B.PerfilGestor}},
		{"perfil de A para usuário de B", `INSERT INTO app.usuario_perfil (empresa_id, usuario_id, perfil_id) VALUES ($1, $2, $3)`,
			[]any{f.A.ID, f.B.UsuarioInterno, f.A.PerfilGestor}},
		{"permissão avulsa a usuário de B", `INSERT INTO app.usuario_permissao (empresa_id, usuario_id, permissao) VALUES ($1, $2, 'os.executar')`,
			[]any{f.A.ID, f.B.UsuarioInterno}},
		{"usuário de A ligado a pessoa de B", `INSERT INTO app.usuario (empresa_id, nome, pessoa_id, moeda) VALUES ($1, 'x', $2, $3)`,
			[]any{f.A.ID, f.B.PessoaInterna, f.A.Moeda}},
		{"pessoa de A na unidade de B", `INSERT INTO app.pessoa (empresa_id, nome, unidade_id) VALUES ($1, 'x', $2)`,
			[]any{f.A.ID, f.B.Unidade}},
	}
	for _, c := range casos {
		t.Run(c.rotulo, func(t *testing.T) {
			err := comoPapel(ctx, papelApp, f.internoA(), func(tx pgx.Tx) error {
				_, err := tx.Exec(ctx, c.sql, c.args...)
				return err
			})
			exigirErro(t, c.rotulo, err, erroFK)
		})
	}
}

// temInsert diz se o papel tem INSERT na tabela (em alguma coluna).
func temInsert(t *testing.T, ctx context.Context, papel, tabela string) bool {
	t.Helper()
	var tem bool
	if err := pool.QueryRow(ctx, `SELECT has_table_privilege($1, $2::regclass, 'INSERT')
	                                  OR has_any_column_privilege($1, $2::regclass, 'INSERT')`, papel, tabela).Scan(&tem); err != nil {
		t.Fatalf("conferir INSERT de %s em %s: %v", papel, tabela, err)
	}
	return tem
}
