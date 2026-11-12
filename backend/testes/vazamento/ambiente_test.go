// Package vazamento é a suíte de isolamento da Tenaz.
//
// Roda em TODO build. Falha bloqueia o merge e o deploy, e é tratada como
// incidente de segurança de severidade alta mesmo quando acontece só no CI
// (seção 9.3 de docs/5_ARQUITETURA_SOFTWARE.md). Não se reexecuta um vermelho
// daqui esperando que passe.
//
// Cobre os três riscos críticos do sumário executivo:
//
//	R2  vazamento entre empresas          — RNF-013
//	R4  prestador ou locatário vendo dado alheio — RNF-037, RNF-038
//	R5  sincronização copiando além do permitido — RNF-041, RNF-042
//
// Nesta etapa (E00) ainda não existe tabela de negócio, então o que a suíte
// verifica é a COBERTURA: que o mecanismo de detecção existe e quebra o build
// no instante em que uma tabela nascer sem defesa. A prova de que o detector
// funciona está em TestDetectorAcusaTabelaSemPolitica, que cria de propósito
// uma tabela malformada e exige que a verificação a acuse.
package vazamento

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"tenaz/internal/plataforma/migracao"
)

// imagemPostgres acompanha a versão usada em desenvolvimento e em produção.
// Testar RLS contra uma versão diferente da de produção não prova nada.
const imagemPostgres = "postgres:17-alpine"

var pool *pgxpool.Pool

func TestMain(m *testing.M) {
	ctx := context.Background()

	url, encerrar, err := subirBanco(ctx)
	if err != nil {
		fmt.Fprintf(os.Stderr, "não foi possível preparar o banco de teste: %v\n", err)
		os.Exit(1)
	}

	if err := migracao.Aplicar(ctx, url); err != nil {
		fmt.Fprintf(os.Stderr, "aplicar migrações: %v\n", err)
		encerrar()
		os.Exit(1)
	}

	pool, err = pgxpool.New(ctx, url)
	if err != nil {
		fmt.Fprintf(os.Stderr, "abrir pool: %v\n", err)
		encerrar()
		os.Exit(1)
	}

	codigo := m.Run()

	pool.Close()
	encerrar()
	os.Exit(codigo)
}

// subirBanco usa um PostgreSQL descartável em contêiner. Quando
// TENAZ_TESTE_BANCO_URL está definida, usa esse banco — útil no ciclo local
// com o compose de desenvolvimento, em que subir contêiner a cada execução
// custa mais que o teste.
func subirBanco(ctx context.Context) (string, func(), error) {
	if url := os.Getenv("TENAZ_TESTE_BANCO_URL"); url != "" {
		return url, func() {}, nil
	}

	cont, err := postgres.Run(ctx, imagemPostgres,
		postgres.WithDatabase("tenaz"),
		postgres.WithUsername("tenaz_migrador"),
		postgres.WithPassword("tenaz_teste"),
		postgres.BasicWaitStrategies(),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(90*time.Second),
		),
	)
	if err != nil {
		return "", func() {}, fmt.Errorf("subir contêiner: %w", err)
	}

	url, err := cont.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = testcontainers.TerminateContainer(cont)
		return "", func() {}, fmt.Errorf("obter URL: %w", err)
	}

	return url, func() { _ = testcontainers.TerminateContainer(cont) }, nil
}
