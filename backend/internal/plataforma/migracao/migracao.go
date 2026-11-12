// Package migracao aplica as migrações do banco.
//
// As migrações vão embutidas no binário: o artefato que roda é o mesmo que
// carrega o esquema que ele espera, o que elimina a chance de aplicar a
// migração de uma versão e subir o código de outra.
//
// Padrão expandir e contrair (seção 14 de docs/2_ARQUITETURA_BANCO_DADOS.md):
// toda migração é compatível com a versão anterior do código, porque as duas
// convivem durante o deploy contínuo.
package migracao

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pressly/goose/v3"

	_ "github.com/jackc/pgx/v5/stdlib" // driver database/sql para o goose

	esquema "tenaz/db"
)

func abrir(url string) (*sql.DB, error) {
	bd, err := sql.Open("pgx", url)
	if err != nil {
		return nil, fmt.Errorf("abrir conexão: %w", err)
	}
	return bd, nil
}

func preparar() error {
	goose.SetBaseFS(esquema.Migracoes)
	goose.SetTableName("plataforma.migracao")
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("definir dialeto: %w", err)
	}
	return nil
}

// Aplicar sobe o banco até a última migração.
func Aplicar(ctx context.Context, url string) error {
	if err := preparar(); err != nil {
		return err
	}
	bd, err := abrir(url)
	if err != nil {
		return err
	}
	defer bd.Close()

	// O esquema plataforma precisa existir antes de o goose criar a própria
	// tabela de controle dentro dele.
	if _, err := bd.ExecContext(ctx, `CREATE SCHEMA IF NOT EXISTS plataforma`); err != nil {
		return fmt.Errorf("criar esquema plataforma: %w", err)
	}

	if err := goose.UpContext(ctx, bd, esquema.DiretorioMigracoes); err != nil {
		return fmt.Errorf("aplicar migrações: %w", err)
	}
	return nil
}

// Situacao imprime o estado das migrações.
func Situacao(ctx context.Context, url string) error {
	if err := preparar(); err != nil {
		return err
	}
	bd, err := abrir(url)
	if err != nil {
		return err
	}
	defer bd.Close()

	return goose.StatusContext(ctx, bd, esquema.DiretorioMigracoes)
}

// Versao devolve a versão aplicada.
func Versao(ctx context.Context, url string) (int64, error) {
	if err := preparar(); err != nil {
		return 0, err
	}
	bd, err := abrir(url)
	if err != nil {
		return 0, err
	}
	defer bd.Close()

	return goose.GetDBVersionContext(ctx, bd)
}
