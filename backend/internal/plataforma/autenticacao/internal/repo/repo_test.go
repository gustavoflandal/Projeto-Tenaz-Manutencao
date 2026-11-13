package repo

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"tenaz/internal/plataforma/autenticacao/internal/repo/gerado"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/db/dbteste"
)

// O repositório da autenticação recusa rodar dentro de NaEmpresa (seção 6.9
// do documento 5; ISO-E01-02): nem Tx nem NaTrilha abrem transação, e a
// função não é chamada. Fora de NaEmpresa, as duas funcionam no pool do papel
// tenaz_autenticacao.
func TestRepositorioRecusaRodarDentroDeNaEmpresa(t *testing.T) {
	ctx := context.Background()
	amb, err := dbteste.Preparar(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(amb.Encerrar)
	r := Novo(amb.Auth)
	empresa := uuid.Must(uuid.NewV7())

	chamou := false
	fn := func(*gerado.Queries) error { chamou = true; return nil }
	err = amb.App.NaEmpresa(ctx, empresa, db.Interno(nil), func(ctx context.Context) error {
		if err := r.Tx(ctx, fn); !errors.Is(err, ErrDentroDeTransacaoDaAplicacao) {
			t.Errorf("Tx dentro de NaEmpresa: %v", err)
		}
		if err := r.NaTrilha(ctx, empresa, fn); !errors.Is(err, ErrDentroDeTransacaoDaAplicacao) {
			t.Errorf("NaTrilha dentro de NaEmpresa: %v", err)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if chamou {
		t.Error("a função do repositório não pode rodar dentro de NaEmpresa")
	}

	// Fora de NaEmpresa: Tx roda como tenaz_autenticacao; NaTrilha exige empresa.
	if err := r.Tx(ctx, func(q *gerado.Queries) error {
		_, err := q.AuthBuscarIdentidade(ctx, uuid.New())
		return err
	}); err == nil || errors.Is(err, ErrDentroDeTransacaoDaAplicacao) {
		t.Errorf("fora de NaEmpresa, Tx chega ao banco (sem linhas): %v", err)
	}
	if err := r.NaTrilha(ctx, uuid.Nil, fn); err == nil || chamou {
		t.Errorf("NaTrilha sem empresa: %v", err)
	}
	var papel string
	if err := amb.Auth.QueryRow(ctx, `SELECT current_user`).Scan(&papel); err != nil || papel != "tenaz_autenticacao" {
		t.Errorf("pool da autenticação: %q %v", papel, err)
	}
}
