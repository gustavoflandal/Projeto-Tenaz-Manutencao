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
// Organização (a partir de E01):
//
//	cobertura_test.go  RLS forçado e policies exatas em toda tabela; partições;
//	                   exceções versionadas; tabelas globais
//	catalogo_test.go   FK composta, unicidade, lista de acesso externo,
//	                   SECURITY DEFINER, mapa de privilégios fora do RLS
//	empresas_test.go   A contra B em toda tabela de app, gerada do catálogo
//	escopo_test.go     prestador, locatário e link dentro de A
//	auditoria_test.go  imutabilidade, mesma transação, isolamento da trilha
//	plataforma_test.go credenciais de autenticação, que não têm RLS
//
// Toda regra de catálogo tem a sua prova de detector (TestDetector...): a
// suíte cria de propósito um objeto malformado e exige ser acusada. Sem isso,
// uma verificação errada passaria por vacuidade. Falhas de isolamento já
// demonstradas e devolvidas ao dono da migração ficam em achadosAbertos
// (banco_test.go): o teste roda, aparece como SKIP com o ID do achado e, enquanto
// a lista não estiver vazia, o portão de isolamento da etapa está vermelho.
package vazamento

import (
	"context"
	"flag"
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

var (
	pool     *pgxpool.Pool
	urlBanco string // para os testes que abrem pool próprio (db.NaEmpresa real)
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	url, encerrar, err := subirBanco(ctx)
	urlBanco = url
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

	// O portão dos achados abertos vale para a suíte inteira, não para um
	// teste: chave sem prova, ou achado aberto no CI, encerra com erro.
	execucaoCompleta := true
	if f := flag.Lookup("test.run"); f != nil && f.Value.String() != "" {
		execucaoCompleta = false
	}
	for _, p := range problemasDosAchados(execucaoCompleta) {
		fmt.Fprintln(os.Stderr, "FALHA DO PORTÃO DE ISOLAMENTO:", p)
		codigo = 1
	}

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
