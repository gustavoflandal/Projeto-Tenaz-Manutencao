package db_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/db/dbteste"
	"tenaz/internal/plataforma/erros"
)

// NaEmpresa contra PostgreSQL real (seção 2.4 do documento 2): o erro da
// função desfaz tudo, inclusive o evento da trilha; o contexto de empresa é
// local à transação; transação aninhada é recusada sem abrir a segunda.
func TestNaEmpresaIntegracao(t *testing.T) {
	ctx := context.Background()
	amb, err := dbteste.Preparar(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(amb.Encerrar)

	contar := func(t *testing.T, sql string, args ...any) int {
		t.Helper()
		var n int
		if err := amb.Migrador.QueryRow(ctx, sql, args...).Scan(&n); err != nil {
			t.Fatal(err)
		}
		return n
	}
	// criarEmpresa grava a empresa pelo pool da aplicação, como CriarConta.
	criarEmpresa := func(ctx context.Context, id uuid.UUID) error {
		tx, err := db.Tx(ctx)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `INSERT INTO app.empresa (id, razao_social) VALUES ($1, 'Teste NaEmpresa')`, id)
		return err
	}

	t.Run("erro na função desfaz a escrita e o evento", func(t *testing.T) {
		id := uuid.Must(uuid.NewV7())
		falha := errors.New("falha do caso de uso")
		err := amb.App.NaEmpresa(ctx, id, db.Interno(nil), func(ctx context.Context) error {
			if err := criarEmpresa(ctx, id); err != nil {
				return err
			}
			return falha
		})
		if !errors.Is(err, falha) {
			t.Fatalf("o erro da função volta intacto: %v", err)
		}
		if n := contar(t, `SELECT count(*) FROM app.empresa WHERE id = $1`, id); n != 0 {
			t.Errorf("empresa gravada apesar do erro: %d", n)
		}
		if n := contar(t, `SELECT count(*) FROM app.evento_auditoria WHERE empresa_id = $1`, id); n != 0 {
			t.Errorf("evento da trilha sobreviveu ao rollback: %d", n)
		}
	})

	t.Run("pânico na função desfaz a escrita", func(t *testing.T) {
		id := uuid.Must(uuid.NewV7())
		func() {
			defer func() { _ = recover() }()
			_ = amb.App.NaEmpresa(ctx, id, db.Interno(nil), func(ctx context.Context) error {
				if err := criarEmpresa(ctx, id); err != nil {
					t.Error(err)
				}
				panic("erro de programação")
			})
		}()
		if n := contar(t, `SELECT count(*) FROM app.empresa WHERE id = $1`, id); n != 0 {
			t.Errorf("empresa gravada apesar do pânico: %d", n)
		}
	})

	t.Run("sucesso confirma a escrita e o evento", func(t *testing.T) {
		id := uuid.Must(uuid.NewV7())
		if err := amb.App.NaEmpresa(ctx, id, db.Interno(nil), func(ctx context.Context) error {
			return criarEmpresa(ctx, id)
		}); err != nil {
			t.Fatal(err)
		}
		if n := contar(t, `SELECT count(*) FROM app.empresa WHERE id = $1`, id); n != 1 {
			t.Errorf("empresa não gravada: %d", n)
		}
		if n := contar(t, `SELECT count(*) FROM app.evento_auditoria WHERE empresa_id = $1 AND acao = 'inclusao'`, id); n != 1 {
			t.Errorf("evento de inclusão: %d", n)
		}
	})

	t.Run("aninhamento recusado sem abrir a segunda transação", func(t *testing.T) {
		a, b := uuid.Must(uuid.NewV7()), uuid.Must(uuid.NewV7())
		chamou := false
		var interno error
		err := amb.App.NaEmpresa(ctx, a, db.Interno(nil), func(ctx context.Context) error {
			for _, alvo := range []uuid.UUID{a, b} {
				interno = amb.App.NaEmpresa(ctx, alvo, db.Interno(nil), func(context.Context) error {
					chamou = true
					return nil
				})
				if !errors.Is(interno, db.ErrTransacaoAninhada) || !errors.Is(interno, erros.SemContextoEmpresa) {
					t.Errorf("aninhada na empresa %s: %v", alvo, interno)
				}
			}
			// A transação de fora continua na empresa dela.
			tx, err := db.Tx(ctx)
			if err != nil {
				return err
			}
			var atual string
			if err := tx.QueryRow(ctx, `SELECT current_setting('app.empresa_id', true)`).Scan(&atual); err != nil {
				return err
			}
			if atual != a.String() {
				t.Errorf("a transação de fora mudou de empresa: %s", atual)
			}
			return interno
		})
		if chamou {
			t.Error("a função aninhada não pode rodar")
		}
		if !errors.Is(err, db.ErrTransacaoAninhada) {
			t.Errorf("o erro do aninhamento sobe: %v", err)
		}
	})

	t.Run("contexto de empresa é local à transação", func(t *testing.T) {
		id := uuid.Must(uuid.NewV7())
		if err := amb.App.NaEmpresa(ctx, id, db.Interno(nil), func(context.Context) error { return nil }); err != nil {
			t.Fatal(err)
		}
		// Todas as conexões do pool: nenhuma ficou com a empresa da anterior.
		for i := 0; i < 8; i++ {
			var atual *string
			if err := amb.App.Pool().QueryRow(ctx, `SELECT nullif(current_setting('app.empresa_id', true), '')`).Scan(&atual); err != nil {
				t.Fatal(err)
			}
			if atual != nil {
				t.Fatalf("app.empresa_id vazou para fora da transação: %s", *atual)
			}
		}
	})

	t.Run("empresa vazia e escopo inválido recusados antes do banco", func(t *testing.T) {
		chamou := false
		fn := func(context.Context) error { chamou = true; return nil }
		if err := amb.App.NaEmpresa(ctx, uuid.Nil, db.Interno(nil), fn); !errors.Is(err, erros.SemContextoEmpresa) {
			t.Errorf("empresa vazia: %v", err)
		}
		if err := amb.App.NaEmpresa(ctx, uuid.New(), db.Escopo{Tipo: db.EscopoPrestador}, fn); !errors.Is(err, erros.SemContextoEmpresa) {
			t.Errorf("escopo prestador sem fornecedor: %v", err)
		}
		if chamou {
			t.Error("a função não pode rodar")
		}
	})
}
