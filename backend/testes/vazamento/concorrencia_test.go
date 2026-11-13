package vazamento

import (
	"context"
	"testing"
)

// ISO-E01-06: corrida em READ COMMITTED entre a atribuição de perfil ou de
// permissão avulsa (tg_usuario_perfil_externo, tg_usuario_permissao_externo) e
// a troca de vínculo (tg_usuario_vinculo). Sem trava, cada transação vê a
// outra no estado antigo e as duas passam: o interno vira prestador levando o
// perfil interno, ou a permissão vedada.
//
// O teste é determinístico: a conexão 1 faz a atribuição e NÃO confirma; a
// conexão 2 tenta trocar o vínculo com lock_timeout curto. Com o FOR SHARE dos
// triggers de atribuição, a linha do usuário está travada e a troca recebe
// 55P03 (lock_not_available). Sem ele, a FK só trava a chave (FOR KEY SHARE),
// a troca de fornecedor_id não é de chave, e passa.

const erroTrava = "55P03"

// resultadoDaCorrida devolve o erro da troca de vínculo na conexão 2 enquanto
// a atribuição da conexão 1 está aberta. comandoDeProva roda na conexão 1, antes da
// atribuição, como superusuário (prova do detector).
func resultadoDaCorrida(t *testing.T, ctx context.Context, f *fixturesE01, atribuicao string, comandoDeProva string) error {
	t.Helper()
	usuario := f.A.UsuarioInativo // interno, sem perfil nem permissão avulsa

	c1, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("conexão 1: %v", err)
	}
	defer c1.Release()
	c2, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("conexão 2: %v", err)
	}
	defer c2.Release()

	tx1, err := c1.Begin(ctx)
	if err != nil {
		t.Fatalf("transação 1: %v", err)
	}
	defer func() { _ = tx1.Rollback(context.Background()) }()
	if comandoDeProva != "" {
		estragar(t, ctx, tx1, comandoDeProva)
	}
	if err := prepararTransacao(ctx, tx1, papelApp, f.internoA()); err != nil {
		t.Fatal(err)
	}
	if _, err := tx1.Exec(ctx, atribuicao, f.A.ID, usuario); err != nil {
		t.Fatalf("a atribuição legítima falhou: %v", err)
	}

	tx2, err := c2.Begin(ctx)
	if err != nil {
		t.Fatalf("transação 2: %v", err)
	}
	defer func() { _ = tx2.Rollback(context.Background()) }()
	if err := prepararTransacao(ctx, tx2, papelApp, f.internoA()); err != nil {
		t.Fatal(err)
	}
	if _, err := tx2.Exec(ctx, "SET LOCAL lock_timeout = '300ms'"); err != nil {
		t.Fatalf("lock_timeout: %v", err)
	}
	_, err = tx2.Exec(ctx, `UPDATE app.usuario SET fornecedor_id = $2 WHERE empresa_id = $1 AND id = $3`,
		f.A.ID, f.F1, usuario)
	return err
}

var atribuicoes = []struct{ nome, sql string }{
	{"perfil interno", `INSERT INTO app.usuario_perfil (empresa_id, usuario_id, perfil_id)
	                    SELECT $1, $2, id FROM app.perfil WHERE empresa_id = $1 AND codigo = 'gestor'`},
	{"permissão avulsa vedada", `INSERT INTO app.usuario_permissao (empresa_id, usuario_id, permissao)
	                             VALUES ($1, $2, 'os.encerrar')`},
}

func TestAtribuicaoTravaATrocaDeVinculo(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	for _, a := range atribuicoes {
		t.Run(a.nome, func(t *testing.T) {
			err := resultadoDaCorrida(t, ctx, f, a.sql, "")
			if sqlstate(err) != erroTrava {
				t.Errorf("com a atribuição de %s aberta, a troca de vínculo deveria esperar a trava (%s); veio %v "+
					"(ISO-E01-06: as duas transações passariam)", a.nome, erroTrava, err)
			}
		})
	}
}

// A prova troca o trigger de atribuição de perfil por uma versão sem o FOR
// SHARE, só dentro da transação 1: a troca de vínculo da conexão 2 passa.
func TestDetectorAcusaAtribuicaoSemTrava(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	semTrava := `
		CREATE OR REPLACE FUNCTION app.tg_usuario_perfil_externo() RETURNS trigger LANGUAGE plpgsql AS $$
		DECLARE v_codigo text; v_fornecedor uuid; v_locatario uuid;
		BEGIN
		  SELECT coalesce(p.perfil_origem, p.codigo) INTO v_codigo FROM app.perfil p
		  WHERE p.empresa_id = NEW.empresa_id AND p.id = NEW.perfil_id;
		  SELECT u.fornecedor_id, u.locatario_id INTO v_fornecedor, v_locatario FROM app.usuario u
		  WHERE u.empresa_id = NEW.empresa_id AND u.id = NEW.usuario_id;
		  RETURN NEW;
		END $$`
	err := resultadoDaCorrida(t, ctx, f, atribuicoes[0].sql, semTrava)
	if err != nil {
		t.Fatalf("sem o FOR SHARE, a troca de vínculo deveria passar e o teste acusar a corrida; veio %v", err)
	}
}
