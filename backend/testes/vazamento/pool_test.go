package vazamento

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"tenaz/internal/plataforma/db"
)

// O contexto não vaza pelo pool, provado com o helper REAL (db.NaEmpresa) e
// um pool de UMA conexão: toda transação, inclusive a que falha no meio, usa a
// mesma conexão física, e a seguinte não pode herdar nada dela (seção 2.4 do
// documento 2; armadilha "SET de sessão com pool de conexões", seção 2.6).
func TestNaEmpresaNaoVazaContextoPeloPool(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)

	banco, err := db.AbrirComOpcoes(ctx, db.Opcoes{URL: urlBanco, MaxConn: 1, Papel: papelApp})
	if err != nil {
		t.Fatalf("abrir pool de uma conexão: %v", err)
	}
	defer banco.Fechar()

	empresaA := uuid.MustParse(f.A.ID)
	empresaB := uuid.MustParse(f.B.ID)
	usuarioF1 := uuid.MustParse(f.UsuarioF1)
	fornecedorF1 := uuid.MustParse(f.F1)
	tabelas := tabelasDeEmpresa(t, ctx)

	// semResto confere, na MESMA conexão, fora de transação de NaEmpresa, que
	// não sobrou contexto nem linha visível.
	semResto := func(depois string) {
		t.Helper()
		var empresaNula, escopoNulo, usuarioNulo, fornecedorNulo bool
		if err := banco.Pool().QueryRow(ctx, `SELECT app.empresa_atual() IS NULL, app.escopo_atual() IS NULL,
			app.usuario_atual() IS NULL, app.fornecedor_atual() IS NULL`).
			Scan(&empresaNula, &escopoNulo, &usuarioNulo, &fornecedorNulo); err != nil {
			t.Fatalf("depois de %s: ler o contexto: %v", depois, err)
		}
		if !empresaNula || !escopoNulo || !usuarioNulo || !fornecedorNulo {
			t.Errorf("depois de %s, a conexão ficou com contexto: empresa=%v escopo=%v usuario=%v fornecedor=%v",
				depois, !empresaNula, !escopoNulo, !usuarioNulo, !fornecedorNulo)
		}
		for _, tb := range tabelas {
			var n int64
			if err := banco.Pool().QueryRow(ctx, "SELECT count(*) FROM "+tb.Nome).Scan(&n); err != nil {
				t.Fatalf("depois de %s: contar %s: %v", depois, tb.Nome, err)
			}
			if n != 0 {
				t.Errorf("depois de %s, a conexão sem contexto viu %d linha(s) de %s", depois, n, tb.Nome)
			}
		}
	}

	contarUnidades := func(c context.Context) (int64, error) {
		tx, err := db.Tx(c)
		if err != nil {
			return 0, err
		}
		var n int64
		err = tx.QueryRow(c, `SELECT count(*) FROM app.unidade`).Scan(&n)
		return n, err
	}

	// 1. Transação que dá certo.
	var vistas int64
	if err := banco.NaEmpresa(ctx, empresaA, db.Interno(nil), func(c context.Context) error {
		var err error
		vistas, err = contarUnidades(c)
		return err
	}); err != nil {
		t.Fatalf("NaEmpresa(A): %v", err)
	}
	if vistas == 0 {
		t.Fatal("NaEmpresa(A) não viu a unidade de A: o teste não prova nada se o contexto nunca valeu")
	}
	semResto("uma transação confirmada")

	// 2. Transação que falha no banco, no meio (transação abortada).
	errSQL := banco.NaEmpresa(ctx, empresaA, db.Interno(nil), func(c context.Context) error {
		tx, err := db.Tx(c)
		if err != nil {
			return err
		}
		_, err = tx.Exec(c, `SELECT 1/0`)
		return err
	})
	if errSQL == nil {
		t.Fatal("a divisão por zero deveria ter falhado")
	}
	semResto("uma transação abortada por erro de SQL")

	// 3. Transação que o caso de uso desfaz com erro próprio.
	errProprio := errors.New("regra de negócio recusou")
	if err := banco.NaEmpresa(ctx, empresaA, db.Interno(nil), func(c context.Context) error {
		return errProprio
	}); !errors.Is(err, errProprio) {
		t.Fatalf("o erro do caso de uso não voltou: %v", err)
	}
	semResto("uma transação desfeita pelo caso de uso")

	// 4. Escopo externo, que fixa mais variáveis.
	if err := banco.NaEmpresa(ctx, empresaA, db.Escopo{Tipo: db.EscopoPrestador, FornecedorID: &fornecedorF1,
		UsuarioID: &usuarioF1}, func(c context.Context) error {
		_, err := contarUnidades(c)
		return err
	}); err != nil {
		t.Fatalf("NaEmpresa(A, prestador F1): %v", err)
	}
	semResto("uma transação de prestador")

	// 5. A empresa seguinte, na mesma conexão, não vê nada da anterior.
	if err := banco.NaEmpresa(ctx, empresaB, db.Interno(nil), func(c context.Context) error {
		tx, err := db.Tx(c)
		if err != nil {
			return err
		}
		var n int64
		if err := tx.QueryRow(c, `SELECT count(*) FROM app.unidade WHERE id = $1`, f.A.Unidade).Scan(&n); err != nil {
			return err
		}
		if n != 0 {
			t.Errorf("NaEmpresa(B), logo depois de A na mesma conexão, viu a unidade de A")
		}
		return nil
	}); err != nil {
		t.Fatalf("NaEmpresa(B): %v", err)
	}
}
