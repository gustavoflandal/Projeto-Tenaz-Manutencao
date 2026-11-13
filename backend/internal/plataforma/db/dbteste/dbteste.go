// Package dbteste prepara um PostgreSQL real para os testes de integração,
// no mesmo padrão da suíte de vazamento: contêiner descartável, ou o banco de
// TENAZ_TESTE_BANCO_URL quando definida (ciclo local mais rápido; o banco
// precisa estar com as migrações atuais).
//
// Os pools de teste assumem os papéis de produção por SET ROLE sobre a URL do
// migrador: o da aplicação como tenaz_app (RLS valendo) e o da autenticação
// como tenaz_autenticacao (ADR-011). Só código de teste importa este pacote.
package dbteste

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/migracao"
)

// Imagem igual à de desenvolvimento e produção.
const imagemPostgres = "postgres:17-alpine"

// Ambiente são os três pools de um teste de integração.
type Ambiente struct {
	URL      string        // migrador
	App      *db.DB        // tenaz_app, só NaEmpresa
	Auth     *pgxpool.Pool // tenaz_autenticacao
	Migrador *pgxpool.Pool // superusuário: só para preparar e inspecionar o cenário
	encerrar func()
}

// Encerrar fecha os pools e o contêiner.
func (a *Ambiente) Encerrar() {
	if a.App != nil {
		a.App.Fechar()
	}
	if a.Auth != nil {
		a.Auth.Close()
	}
	if a.Migrador != nil {
		a.Migrador.Close()
	}
	if a.encerrar != nil {
		a.encerrar()
	}
}

// Preparar sobe o banco, aplica as migrações e abre os pools.
func Preparar(ctx context.Context) (*Ambiente, error) {
	url, encerrar, err := subir(ctx)
	if err != nil {
		return nil, err
	}
	amb := &Ambiente{URL: url, encerrar: encerrar}
	if err := migracao.Aplicar(ctx, url); err != nil {
		amb.Encerrar()
		return nil, fmt.Errorf("aplicar migrações: %w", err)
	}
	if amb.Migrador, err = pgxpool.New(ctx, url); err != nil {
		amb.Encerrar()
		return nil, err
	}
	if amb.App, err = db.AbrirComOpcoes(ctx, db.Opcoes{URL: url, MaxConn: 8, Papel: "tenaz_app"}); err != nil {
		amb.Encerrar()
		return nil, err
	}
	if amb.Auth, err = db.AbrirPool(ctx, db.Opcoes{URL: url, MaxConn: 4, Papel: "tenaz_autenticacao"}); err != nil {
		amb.Encerrar()
		return nil, fmt.Errorf("pool da autenticação (o papel tenaz_autenticacao existe?): %w", err)
	}
	return amb, nil
}

func subir(ctx context.Context) (string, func(), error) {
	if url := os.Getenv("TENAZ_TESTE_BANCO_URL"); url != "" {
		return url, func() {}, nil
	}
	cont, err := postgres.Run(ctx, imagemPostgres,
		postgres.WithDatabase("tenaz"),
		postgres.WithUsername("tenaz_migrador"),
		postgres.WithPassword("tenaz_teste"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(90*time.Second),
		),
	)
	if err != nil {
		return "", nil, fmt.Errorf("subir PostgreSQL de teste: %w", err)
	}
	encerrar := func() { _ = testcontainers.TerminateContainer(cont) }
	url, err := cont.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		encerrar()
		return "", nil, err
	}
	return url, encerrar, nil
}
