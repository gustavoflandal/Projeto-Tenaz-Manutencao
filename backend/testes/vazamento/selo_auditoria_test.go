package vazamento

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
)

// Trilha configurável e selo diário (00007; seções 9.2.1 e 9.4 do documento
// 2): o hash do evento não se forja, o selo só o tenaz_job grava e ninguém
// altera, o selo não se forja, e a verificação de integridade é por empresa.
// A captura por tela (app.captura_auditoria) é tabela de empresa comum, só
// interna: entra pelos testes gerados do catálogo (empresas_test.go e
// escopo_test.go), com linha nas fixtures.

const hashZerado = `decode(repeat('00', 32), 'hex')`

// colunasDoServidor são as colunas da trilha que o servidor preenche e que
// nenhum papel que grava evento em nome da empresa pode mandar (00008): id,
// ocorrido_em e transacao (valores padrão) e hash (trigger). Com ocorrido_em à
// escolha, um evento era gravado com data retroativa num dia ainda não selado,
// e o selo o selaria sem perceber (achado do tenaz-isolamento sobre a 00007).
var colunasDoServidor = map[string]string{
	"id":          "app.uuid_v7()",
	"ocorrido_em": "now() - interval '2 days'",
	"transacao":   "pg_current_xact_id()",
	"hash":        hashZerado,
}

// problemasDoEventoForjado: o tenaz_app e o tenaz_job, no contexto de A, não
// mandam nenhuma coluna do servidor (permissão negada pelo privilégio, e não
// pela policy); sem elas, o evento grava, com ocorrido_em = now() (o início da
// transação), a transação corrente e o hash do conteúdo.
func problemasDoEventoForjado(t *testing.T, ctx context.Context, tx pgx.Tx, f *fixturesE01) []string {
	t.Helper()
	var p []string
	for _, papel := range []string{papelApp, papelJob} {
		for col, valor := range colunasDoServidor {
			err := tentarComo(ctx, tx, papel, f.internoA(), `INSERT INTO app.evento_auditoria (empresa_id, ator_tipo, acao, entidade, `+
				col+`) VALUES ($1, 'sistema', 'inclusao', 'unidade', `+valor+`)`, f.A.ID)
			if err == nil || recusaDeRLS(err) || sqlstate(err) != erroPrivilegio {
				p = append(p, fmt.Sprintf("%s mandou a coluna %s (resultado: %v)", papel, col, err))
			}
		}
		sub, err := tx.Begin(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if err := trocarPapel(ctx, sub, papel, f.internoA()); err != nil {
			t.Fatal(err)
		}
		var id string
		if err := sub.QueryRow(ctx, `INSERT INTO app.evento_auditoria (empresa_id, ator_tipo, acao, entidade)
		                             VALUES ($1, 'sistema', 'inclusao', 'unidade') RETURNING id`, f.A.ID).Scan(&id); err != nil {
			p = append(p, fmt.Sprintf("%s não grava o evento sem as colunas do servidor: %v", papel, err))
			_ = sub.Rollback(ctx)
			continue
		}
		if err := trocarPapel(ctx, sub, "", contexto{}); err != nil {
			t.Fatal(err)
		}
		var agora, transacao, conteudo bool
		if err := sub.QueryRow(ctx, `SELECT ocorrido_em = now(), transacao = pg_current_xact_id(), hash = app.hash_evento_auditoria(e)
		                             FROM app.evento_auditoria e WHERE id = $1`, id).Scan(&agora, &transacao, &conteudo); err != nil {
			t.Fatalf("ler o evento: %v", err)
		}
		_ = sub.Rollback(ctx)
		if !agora || !transacao || !conteudo {
			p = append(p, fmt.Sprintf("%s: o servidor não preencheu o evento (ocorrido_em=now(): %v, transação corrente: %v, hash do conteúdo: %v)",
				papel, agora, transacao, conteudo))
		}
	}
	return p
}

func TestEventoRetroativoERecusado(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		for _, p := range problemasDoEventoForjado(t, ctx, tx, f) {
			t.Errorf("%s (00008; seção 9.2 do documento 2)", p)
		}
	})
}

// A prova devolve o INSERT de tabela inteira que a 00008 retirou.
func TestDetectorAcusaEventoRetroativo(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		estragar(t, ctx, tx, "GRANT INSERT ON app.evento_auditoria TO tenaz_app")
		p := strings.Join(problemasDoEventoForjado(t, ctx, tx, f), " | ")
		for _, col := range []string{"ocorrido_em", "transacao", "hash", "id"} {
			if !strings.Contains(p, "tenaz_app mandou a coluna "+col+" ") {
				t.Errorf("com o INSERT de tabela inteira, o detector não acusou a coluna %s: %q", col, p)
			}
		}
		if strings.Contains(p, "tenaz_job mandou") {
			t.Errorf("a prova estragou só o tenaz_app, e o tenaz_job foi acusado: %q", p)
		}
	})
}

// problemasDoHash: ninguém além do migrador consegue mandar a coluna hash
// (acima); o migrador manda o hash zerado e o trigger o troca pelo do
// conteúdo. É o que garante que nem a conversão de uma migração grava hash
// forjado.
func problemasDoHash(t *testing.T, ctx context.Context, tx pgx.Tx, f *fixturesE01) []string {
	t.Helper()
	if err := trocarPapel(ctx, tx, "", contexto{}); err != nil {
		t.Fatal(err)
	}
	var forjado, confere bool
	if err := tx.QueryRow(ctx, `WITH e AS (
	                              INSERT INTO app.evento_auditoria (empresa_id, ator_tipo, acao, entidade, hash)
	                              VALUES ($1, 'sistema', 'inclusao', 'unidade', `+hashZerado+`) RETURNING *)
	                            SELECT hash = `+hashZerado+`, hash = app.hash_evento_auditoria(e) FROM e`, f.A.ID).Scan(&forjado, &confere); err != nil {
		t.Fatalf("gravar evento com hash: %v", err)
	}
	var p []string
	if forjado {
		p = append(p, "o hash mandado pelo chamador foi gravado")
	}
	if !confere {
		p = append(p, "o hash gravado não é o do conteúdo do evento")
	}
	return p
}

func TestHashDoEventoNaoSeForja(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		for _, p := range problemasDoHash(t, ctx, tx, f) {
			t.Errorf("%s (seção 9.2.1)", p)
		}
	})
}

func TestDetectorAcusaHashForjavel(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		estragar(t, ctx, tx, "DROP TRIGGER trg_evento_preparar ON app.evento_auditoria")
		if p := strings.Join(problemasDoHash(t, ctx, tx, f), " | "); !strings.Contains(p, "foi gravado") {
			t.Errorf("sem o trigger do hash, o hash forjado não foi acusado: %q", p)
		}
	})
}

// Privilégio concedido por coluna (SELECT, INSERT, UPDATE, REFERENCES) não
// aparece em has_table_privilege: as verificações de privilégio da suíte usam
// has_any_column_privilege, e esta prova concede por coluna e exige a acusação.
func TestDetectorAcusaPrivilegioPorColuna(t *testing.T) {
	ctx := context.Background()
	var particao string
	if err := pool.QueryRow(ctx, `SELECT i.inhrelid::regclass::text FROM pg_inherits i
	                              WHERE i.inhparent = 'app.evento_auditoria'::regclass ORDER BY 1 LIMIT 1`).Scan(&particao); err != nil {
		t.Fatalf("achar uma partição da trilha: %v", err)
	}
	tx := criarTabelasDeProva(t, ctx, []string{
		"GRANT REFERENCES (id) ON app.pessoa TO tenaz_app",
		"GRANT UPDATE (acao) ON app.evento_auditoria TO tenaz_app",
		"GRANT SELECT (empresa_id) ON " + particao + " TO tenaz_app",
	})
	if !contem(listarTexto(t, ctx, tx, sqlPrivilegiosProibidos), "app.pessoa tenaz_app REFERENCES") {
		t.Error("REFERENCES por coluna em app.pessoa não foi acusado")
	}
	if !contem(problemasSomenteInsercao(t, ctx, tx, "app.evento_auditoria"), "algum papel da aplicação tem UPDATE, DELETE ou TRUNCATE") {
		t.Error("UPDATE por coluna em app.evento_auditoria não foi acusado pela forma somente_insercao")
	}
	if !contem(particoesSemDefesa(t, ctx, tx), particao) {
		t.Errorf("SELECT por coluna na partição %s não foi acusado", particao)
	}
	if _, abertos := acessoExterno(t, ctx, tx)["app.evento_auditoria"]; !abertos {
		t.Error("o cálculo do acesso externo não vê o INSERT por coluna da trilha")
	}
}

// O selo: só o tenaz_job insere, ninguém altera nem apaga (nem o
// superusuário, nem por TRUNCATE), e a função do job só o job executa.
func TestSeloSoOJobGravaENinguemAltera(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	semInsercao := []string{papelApp, papelLeitura, papelSinc, "tenaz_sinc_conexao", papelAutenticacao, papelAutenticacaoConexao}
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		for _, papel := range semInsercao {
			err := tentarComo(ctx, tx, papel, f.internoA(), `INSERT INTO app.selo_auditoria (empresa_id, dia) VALUES ($1, current_date - 40)`, f.A.ID)
			exigirErro(t, papel+": INSERT em app.selo_auditoria", err, erroPrivilegio)
		}
		for _, papel := range papeisDeExecucao {
			for _, priv := range []string{"UPDATE", "DELETE", "TRUNCATE"} {
				var tem bool
				if err := tx.QueryRow(ctx, `SELECT `+sqlPrivilegioComColuna("$1::text", "'app.selo_auditoria'::regclass", "$2::text"), papel, priv).Scan(&tem); err != nil {
					t.Fatal(err)
				}
				if tem && papel != papelManutencao {
					t.Errorf("%s tem %s em app.selo_auditoria", papel, priv)
				}
			}
		}
		// Nem quem passa por cima do RLS altera ou apaga um selo.
		for _, sql := range []string{
			`UPDATE app.selo_auditoria SET quantidade = quantidade WHERE empresa_id = $1`,
			`DELETE FROM app.selo_auditoria WHERE empresa_id = $1`,
		} {
			if err := tentarComo(ctx, tx, "", contexto{}, sql, f.A.ID); err == nil {
				t.Errorf("o superusuário executou %q no selo: a trava de somente inserção não existe", sql)
			}
		}
		if err := tentarComo(ctx, tx, "", contexto{}, `TRUNCATE app.selo_auditoria`); err == nil {
			t.Error("o superusuário truncou app.selo_auditoria: falta a trava de TRUNCATE")
		}
		// O dono passa pelo RLS forçado, sem policy de UPDATE nem de DELETE.
		for _, sql := range []string{
			`UPDATE app.selo_auditoria SET quantidade = quantidade WHERE empresa_id = $1`,
			`DELETE FROM app.selo_auditoria WHERE empresa_id = $1`,
		} {
			var n int64
			err := func() error {
				sub, err := tx.Begin(ctx)
				if err != nil {
					return err
				}
				defer func() { _ = sub.Rollback(ctx) }()
				if err := trocarPapel(ctx, sub, papelDono, f.internoA()); err != nil {
					return err
				}
				n, err = alterou(ctx, sub, sql, f.A.ID)
				return err
			}()
			naoMudouNada(t, "tenaz_dono: "+sql, n, err)
		}
		// EXECUTE da função do job.
		for _, papel := range append([]string{papelJob}, semInsercao...) {
			var pode bool
			if err := tx.QueryRow(ctx, `SELECT has_function_privilege($1, 'app.selar_dias_pendentes()', 'EXECUTE')`, papel).Scan(&pode); err != nil {
				t.Fatal(err)
			}
			if pode != (papel == papelJob) {
				t.Errorf("EXECUTE em app.selar_dias_pendentes() para %s = %v; só o tenaz_job executa", papel, pode)
			}
		}
		var publico bool
		if err := tx.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM pg_proc p, aclexplode(coalesce(p.proacl, acldefault('f', p.proowner))) a
		                            WHERE p.oid = 'app.selar_dias_pendentes()'::regprocedure AND a.grantee = 0)`).Scan(&publico); err != nil {
			t.Fatal(err)
		}
		if publico {
			t.Error("PUBLIC executa app.selar_dias_pendentes()")
		}
	})
}

// problemasDoSeloForjado cria, na transação, uma empresa com um evento de três
// dias atrás e manda o tenaz_job gravar o primeiro selo com valores forjados:
// o banco precisa recalcular tudo.
func problemasDoSeloForjado(t *testing.T, ctx context.Context, tx pgx.Tx) []string {
	t.Helper()
	x := novoID()
	for _, sql := range []string{
		`INSERT INTO app.empresa (id, razao_social) VALUES ($1, 'Empresa do selo forjado')`,
		`INSERT INTO app.evento_auditoria (empresa_id, ator_tipo, acao, entidade, ocorrido_em)
		 VALUES ($1, 'sistema', 'inclusao', 'unidade', now() - interval '3 days')`,
	} {
		if _, err := tx.Exec(ctx, sql, x); err != nil {
			t.Fatalf("preparar a empresa do selo: %v", err)
		}
	}
	if err := trocarPapel(ctx, tx, papelJob, contexto{Empresa: x, Escopo: "interno"}); err != nil {
		t.Fatal(err)
	}
	var (
		quantidade      int64
		zeradoDia       bool
		zeradoCadeia    bool
		chave, esperada string
	)
	err := tx.QueryRow(ctx, `
		INSERT INTO app.selo_auditoria (empresa_id, dia, quantidade, hash_dia, hash_encadeado, chave_objeto)
		VALUES ($1, ((now() - interval '3 days') AT TIME ZONE 'UTC')::date, 999, `+hashZerado+`, `+hashZerado+`, 'forjado')
		RETURNING quantidade, hash_dia = `+hashZerado+`, hash_encadeado = `+hashZerado+`, chave_objeto,
		          'selos/empresas/' || empresa_id::text || '/' || to_char(dia, 'YYYY-MM-DD') || '.selo'`, x).
		Scan(&quantidade, &zeradoDia, &zeradoCadeia, &chave, &esperada)
	if err != nil {
		t.Fatalf("o tenaz_job não gravou o primeiro selo: %v", err)
	}
	var p []string
	if quantidade != 1 {
		p = append(p, fmt.Sprintf("quantidade %d gravada; o dia tem 1 evento", quantidade))
	}
	if zeradoDia || zeradoCadeia {
		p = append(p, "hash do selo mandado pelo job foi gravado")
	}
	if chave != esperada {
		p = append(p, fmt.Sprintf("chave_objeto %q gravada; esperada %q", chave, esperada))
	}
	return p
}

func TestSeloNaoSeForja(t *testing.T) {
	ctx := context.Background()
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		for _, p := range problemasDoSeloForjado(t, ctx, tx) {
			t.Errorf("%s (seção 9.4)", p)
		}
	})
}

func TestDetectorAcusaSeloForjavel(t *testing.T) {
	ctx := context.Background()
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		estragar(t, ctx, tx, "DROP TRIGGER trg_selo_calcular ON app.selo_auditoria")
		p := strings.Join(problemasDoSeloForjado(t, ctx, tx), " | ")
		for _, esperado := range []string{"quantidade 999", "hash do selo", "chave_objeto"} {
			if !strings.Contains(p, esperado) {
				t.Errorf("sem o trigger do selo, o detector não acusou %q: %q", esperado, p)
			}
		}
	})
}

// A verificação de integridade é da empresa da transação: um evento incluído
// num dia já selado de A aparece na verificação de A, e a de B, no mesmo dia,
// continua íntegra; e a verificação de B não fala de A.
func TestVerificacaoDeIntegridadeEPorEmpresa(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	verificar := func(tx pgx.Tx, e empresaFixture, dia string) (estado, motivo string) {
		t.Helper()
		sub, err := tx.Begin(ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = sub.Rollback(ctx) }()
		if err := trocarPapel(ctx, sub, papelApp, contexto{Empresa: e.ID, Escopo: "interno"}); err != nil {
			t.Fatal(err)
		}
		var m *string
		if err := sub.QueryRow(ctx, `SELECT estado, motivo FROM app.verificar_integridade_auditoria($1::date, $1::date)`, dia).
			Scan(&estado, &m); err != nil {
			t.Fatalf("verificar a empresa %.8s: %v", e.ID, err)
		}
		if m != nil {
			motivo = *m
		}
		return estado, motivo
	}
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		if f.A.SeloDia == "" || f.A.SeloDia != f.B.SeloDia {
			t.Fatalf("as fixtures precisam do mesmo dia selado em A e em B: %q, %q", f.A.SeloDia, f.B.SeloDia)
		}
		dia := f.A.SeloDia
		for _, e := range []empresaFixture{f.A, f.B} {
			if estado, motivo := verificar(tx, e, dia); estado != "integro" {
				t.Fatalf("antes da adulteração, o dia %s da empresa %.8s está %s (%s)", dia, e.ID, estado, motivo)
			}
		}
		// Um evento incluído em A no dia selado (só o migrador consegue: a
		// trilha não tem UPDATE, e incluir é o ataque que sobra).
		if err := trocarPapel(ctx, tx, "", contexto{}); err != nil {
			t.Fatal(err)
		}
		if _, err := tx.Exec(ctx, `INSERT INTO app.evento_auditoria (empresa_id, ator_tipo, acao, entidade, ocorrido_em)
		                           VALUES ($1, 'sistema', 'inclusao', 'unidade', $2::date + interval '12 hours')`, f.A.ID, dia); err != nil {
			t.Fatal(err)
		}
		if estado, motivo := verificar(tx, f.A, dia); estado != "violado" || motivo != "eventos_divergentes" {
			t.Errorf("a inclusão no dia selado de A não foi acusada na verificação de A: %s (%s)", estado, motivo)
		}
		if estado, motivo := verificar(tx, f.B, dia); estado != "integro" {
			t.Errorf("a adulteração de A apareceu na verificação de B: %s (%s); a verificação não é por empresa", estado, motivo)
		}
	})
}
