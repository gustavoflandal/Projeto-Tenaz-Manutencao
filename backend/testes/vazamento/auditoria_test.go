package vazamento

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
)

// Trilha de auditoria (RF-ADM-011; seções 9.2 a 9.5 do documento 2):
// imutável, gravada na mesma transação do fato e isolada por empresa.

// UPDATE e DELETE recusados por todos: pela falta de privilégio (aplicação),
// pelo RLS sem policy (dono, que só tem o RLS forçado) e pelo trigger
// (superusuário, que ignora o RLS). Cada camada sozinha basta; o teste prova
// as três.
func TestAuditoriaNaoSeAlteraNemSeApaga(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)

	var eventos int64
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM app.evento_auditoria WHERE empresa_id = $1`, f.A.ID).Scan(&eventos); err != nil {
		t.Fatalf("contar eventos de A: %v", err)
	}
	if eventos == 0 {
		t.Fatal("as fixtures de A não geraram nenhum evento de auditoria; o trigger não está ligado")
	}

	for _, papel := range []string{papelApp, papelJob, papelLeitura} {
		t.Run(papel, func(t *testing.T) {
			for _, sql := range []string{
				`UPDATE app.evento_auditoria SET acao = acao WHERE empresa_id = $1`,
				`DELETE FROM app.evento_auditoria WHERE empresa_id = $1`,
			} {
				err := comoPapel(ctx, papel, f.internoA(), func(tx pgx.Tx) error {
					_, err := tx.Exec(ctx, sql, f.A.ID)
					return err
				})
				exigirErro(t, papel+": "+sql, err, erroPrivilegio)
			}
		})
	}

	t.Run("dono", func(t *testing.T) {
		// O dono tem o privilégio, mas o RLS forçado e a falta de policy para
		// ele o deixam sem linha nenhuma.
		for _, sql := range []string{
			`UPDATE app.evento_auditoria SET acao = acao WHERE empresa_id = $1`,
			`DELETE FROM app.evento_auditoria WHERE empresa_id = $1`,
		} {
			var n int64
			err := comoPapel(ctx, papelDono, f.internoA(), func(tx pgx.Tx) error {
				var err error
				n, err = alterou(ctx, tx, sql, f.A.ID)
				return err
			})
			naoMudouNada(t, "dono: "+sql, n, err)
		}
	})

	t.Run("superusuário", func(t *testing.T) {
		// Fora do RLS, sobra o trigger: ele recusa até quem ignora as policies.
		for _, sql := range []string{
			`UPDATE app.evento_auditoria SET acao = acao WHERE empresa_id = $1`,
			`DELETE FROM app.evento_auditoria WHERE empresa_id = $1`,
		} {
			err := comoPapel(ctx, "", contexto{}, func(tx pgx.Tx) error {
				_, err := tx.Exec(ctx, sql, f.A.ID)
				return err
			})
			exigirErro(t, "superusuário: "+sql, err, erroPrivilegio)
		}
	})
}

// A partição não é atalho: o acesso direto, sem passar pelo pai, é recusado
// para a aplicação.
func TestParticaoDeAuditoriaNaoEAtalho(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)

	particoes := listarTexto(t, ctx, pool, `
		SELECT c.oid::regclass::text FROM pg_inherits i JOIN pg_class c ON c.oid = i.inhrelid
		WHERE i.inhparent = 'app.evento_auditoria'::regclass ORDER BY 1`)
	if len(particoes) == 0 {
		t.Fatal("app.evento_auditoria não tem partição; a migração 00003 não criou as mensais")
	}
	if !contem(particoes, "app.evento_auditoria_padrao") {
		t.Error("falta a partição padrão app.evento_auditoria_padrao (seção 9.5)")
	}
	for _, p := range particoes {
		for _, papel := range []string{papelApp, papelJob, papelLeitura} {
			err := comoPapel(ctx, papel, f.internoA(), func(tx pgx.Tx) error {
				_, err := contar(ctx, tx, "SELECT count(*) FROM "+p)
				return err
			})
			exigirErro(t, papel+" lendo "+p+" direto", err, erroPrivilegio)
		}
	}
}

// O evento nasce na transação do fato, com autor e valores antes e depois, e
// some com ela. E o evento de A não aparece para B.
func TestAuditoriaNaMesmaTransacaoEIsoladaPorEmpresa(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)

	c := f.internoA()
	c.Usuario = f.A.UsuarioInterno

	var antes, depois int64
	contarEventosDaUnidade := func(destino *int64) {
		t.Helper()
		if err := pool.QueryRow(ctx, `SELECT count(*) FROM app.evento_auditoria WHERE empresa_id = $1 AND entidade_id = $2`,
			f.A.ID, f.A.Unidade).Scan(destino); err != nil {
			t.Fatalf("contar eventos: %v", err)
		}
	}
	contarEventosDaUnidade(&antes)

	err := comoPapel(ctx, papelApp, c, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `UPDATE app.unidade SET nome = 'Nome auditado' WHERE id = $1`, f.A.Unidade); err != nil {
			return err
		}
		var acao, usuario, nomeAntes, nomeDepois string
		var campos []string
		err := tx.QueryRow(ctx, `
			SELECT acao, coalesce(usuario_id::text, ''), coalesce(antes->>'nome', ''), coalesce(depois->>'nome', ''),
			       coalesce(campos_alterados, '{}')
			FROM app.evento_auditoria
			WHERE entidade = 'unidade' AND entidade_id = $1 AND transacao = pg_current_xact_id()`, f.A.Unidade).
			Scan(&acao, &usuario, &nomeAntes, &nomeDepois, &campos)
		if err != nil {
			t.Errorf("o evento da alteração não está na mesma transação: %v", err)
			return nil
		}
		if acao != "alteracao" || usuario != f.A.UsuarioInterno || nomeAntes != "Unidade" ||
			nomeDepois != "Nome auditado" || !mesmoConjunto(campos, []string{"nome"}) {
			t.Errorf("evento incompleto: acao=%s usuario=%s antes=%q depois=%q campos=%v",
				acao, usuario, nomeAntes, nomeDepois, campos)
		}

		// Na mesma transação, trocando o contexto para B: o evento de A some.
		if _, err := tx.Exec(ctx, `SELECT set_config('app.empresa_id', $1, true), set_config('app.usuario_id', '', true)`,
			f.B.ID); err != nil {
			return err
		}
		n, err := contar(ctx, tx, `SELECT count(*) FROM app.evento_auditoria WHERE entidade_id = $1`, f.A.Unidade)
		if err != nil {
			return err
		}
		if n != 0 {
			t.Errorf("no contexto de B, %d evento(s) da unidade de A ficaram visíveis", n)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("alterar a unidade: %v", err)
	}

	contarEventosDaUnidade(&depois)
	if depois != antes {
		t.Errorf("a transação foi desfeita e a auditoria ficou com %d evento(s) a mais; "+
			"evento sem fato é trilha falsa", depois-antes)
	}
}

// A empresa B não grava evento em nome de A, nem por INSERT direto.
func TestAuditoriaNaoRecebeEventoDeOutraEmpresa(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)

	err := comoPapel(ctx, papelApp, contexto{Empresa: f.B.ID, Escopo: "interno"}, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `INSERT INTO app.evento_auditoria (empresa_id, ator_tipo, acao, entidade)
		                        VALUES ($1, 'sistema', 'alteracao', 'unidade')`, f.A.ID)
		return err
	})
	exigirRLS(t, "B grava evento com empresa A", err)
}
