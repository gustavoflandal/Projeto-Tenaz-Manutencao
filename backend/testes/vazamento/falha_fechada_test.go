package vazamento

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Sem contexto de empresa, nada é visível. É a garantia de "falha fechada":
// um caminho de código que esqueça de abrir a transação por NaEmpresa não
// devolve dado nenhum, em vez de devolver tudo.
func TestSemContextoNaoHaEmpresa(t *testing.T) {
	ctx := context.Background()

	// app.escopo_interno() devolve NULL, e não falso, quando não há escopo:
	// NULL = 'interno' é NULL, pela lógica de três valores do SQL. Isso
	// continua falhando fechado, porque uma policy só libera a linha quando o
	// predicado é verdadeiro, e NULL não é verdadeiro. A asserção correta é
	// portanto "não é verdadeiro", não "é falso".
	var empresaNula, escopoNulo, internoNaoVerdadeiro bool
	err := pool.QueryRow(ctx, `
		SELECT app.empresa_atual()  IS NULL,
		       app.escopo_atual()   IS NULL,
		       app.escopo_interno() IS NOT TRUE`).Scan(&empresaNula, &escopoNulo, &internoNaoVerdadeiro)
	if err != nil {
		t.Fatalf("consultar funções de contexto: %v", err)
	}

	if !empresaNula {
		t.Error("app.empresa_atual() devolveu valor sem contexto definido; a policy não falharia fechada")
	}
	if !escopoNulo {
		t.Error("app.escopo_atual() devolveu valor sem contexto definido")
	}
	if !internoNaoVerdadeiro {
		t.Error("app.escopo_interno() foi verdadeiro sem escopo definido; " +
			"a policy restritiva liberaria a linha para quem não tem escopo")
	}
}

// Confirma, no próprio banco, que um predicado nulo não libera linha. É a
// premissa em que a falha fechada se apoia; se uma versão do PostgreSQL
// mudasse isso, todas as policies restritivas passariam a liberar tudo.
func TestPredicadoNuloNaoLiberaLinha(t *testing.T) {
	ctx := context.Background()

	const nome = "tabela_prova_predicado_nulo"

	t.Cleanup(func() {
		if _, err := pool.Exec(ctx, "DROP TABLE IF EXISTS app."+nome); err != nil {
			t.Errorf("limpar tabela de prova: %v", err)
		}
	})

	for _, comando := range []string{
		"CREATE TABLE app." + nome + " (id int PRIMARY KEY, empresa_id uuid NOT NULL)",
		"INSERT INTO app." + nome + " VALUES (1, gen_random_uuid())",
		"ALTER TABLE app." + nome + " ENABLE ROW LEVEL SECURITY",
		"ALTER TABLE app." + nome + " FORCE ROW LEVEL SECURITY",
		"GRANT SELECT ON app." + nome + " TO tenaz_app",
		// Predicado que resulta em NULL sem contexto definido.
		"CREATE POLICY isolamento_empresa ON app." + nome +
			" FOR ALL TO tenaz_app USING (empresa_id = app.empresa_atual()) WITH CHECK (empresa_id = app.empresa_atual())",
	} {
		if _, err := pool.Exec(ctx, comando); err != nil {
			t.Fatalf("preparar prova (%s): %v", comando, err)
		}
	}

	var visiveis int
	err := pgx.BeginTxFunc(ctx, pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, "SET LOCAL ROLE tenaz_app"); err != nil {
			return err
		}
		return tx.QueryRow(ctx, "SELECT count(*) FROM app."+nome).Scan(&visiveis)
	})
	if err != nil {
		t.Fatalf("consultar como tenaz_app: %v", err)
	}

	if visiveis != 0 {
		t.Errorf("tenaz_app viu %d linha(s) sem contexto de empresa; a falha fechada não está valendo", visiveis)
	}
}

// Depois do fim da transação, a variável local volta a string vazia, e não a
// nulo. Sem o nullif nas funções, o cast estouraria com erro em vez de
// devolver zero linhas — e o próprio erro já contaria ao chamador que havia
// algo ali. Este teste tranca esse comportamento.
func TestVariavelVaziaNaoEstouraOCast(t *testing.T) {
	ctx := context.Background()

	conexao, err := pool.Acquire(ctx)
	if err != nil {
		t.Fatalf("obter conexão: %v", err)
	}
	defer conexao.Release()

	// Define e encerra a transação, deixando a variável como string vazia.
	tx, err := conexao.Begin(ctx)
	if err != nil {
		t.Fatalf("abrir transação: %v", err)
	}
	if _, err := tx.Exec(ctx, `SELECT set_config('app.empresa_id', $1, true)`, uuid.New().String()); err != nil {
		t.Fatalf("definir contexto: %v", err)
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatalf("encerrar transação: %v", err)
	}

	var nula bool
	if err := conexao.QueryRow(ctx, `SELECT app.empresa_atual() IS NULL`).Scan(&nula); err != nil {
		t.Fatalf("app.empresa_atual() estourou depois do commit, em vez de devolver nulo: %v", err)
	}
	if !nula {
		t.Error("app.empresa_atual() não voltou a nulo depois do fim da transação")
	}
}

// app.conferir_escopo() recusa a transação quando o contexto está incoerente.
// Cada caso abaixo é um jeito real de errar.
func TestConferirEscopoRecusaContextoInvalido(t *testing.T) {
	ctx := context.Background()
	empresa := uuid.New().String()
	usuario := uuid.New().String()

	casos := []struct {
		nome     string
		contexto map[string]string
		recusa   bool
		motivo   string
	}{
		{
			nome:     "sem nada definido",
			contexto: map[string]string{},
			recusa:   true,
			motivo:   "escopo ausente",
		},
		{
			nome:     "empresa sem escopo",
			contexto: map[string]string{"app.empresa_id": empresa},
			recusa:   true,
			motivo:   "escopo ausente",
		},
		{
			nome:     "escopo inventado",
			contexto: map[string]string{"app.empresa_id": empresa, "app.escopo": "administrador"},
			recusa:   true,
			motivo:   "escopo não está entre os quatro valores canônicos",
		},
		{
			nome:     "escopo sem empresa",
			contexto: map[string]string{"app.escopo": "interno"},
			recusa:   true,
			motivo:   "empresa ausente",
		},
		{
			nome:     "prestador sem usuário",
			contexto: map[string]string{"app.empresa_id": empresa, "app.escopo": "prestador"},
			recusa:   true,
			motivo:   "escopo externo exige usuário",
		},
		{
			nome:     "locatario sem usuário",
			contexto: map[string]string{"app.empresa_id": empresa, "app.escopo": "locatario"},
			recusa:   true,
			motivo:   "escopo externo exige usuário",
		},
		{
			nome:     "link sem segredo",
			contexto: map[string]string{"app.empresa_id": empresa, "app.escopo": "link"},
			recusa:   true,
			motivo:   "link exige o hash do segredo",
		},
		{
			nome:     "interno com usuário",
			contexto: map[string]string{"app.empresa_id": empresa, "app.escopo": "interno", "app.usuario_id": usuario},
			recusa:   false,
		},
		{
			nome:     "interno sem usuário (job)",
			contexto: map[string]string{"app.empresa_id": empresa, "app.escopo": "interno"},
			recusa:   false,
		},
		{
			nome: "link com segredo",
			contexto: map[string]string{
				"app.empresa_id": empresa,
				"app.escopo":     "link",
				"app.link_hash":  strings.Repeat("ab", 32),
			},
			recusa: false,
		},
	}

	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			err := pgx.BeginTxFunc(ctx, pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
				for chave, valor := range c.contexto {
					if _, err := tx.Exec(ctx, `SELECT set_config($1, $2, true)`, chave, valor); err != nil {
						return err
					}
				}
				_, err := tx.Exec(ctx, `SELECT app.conferir_escopo()`)
				return err
			})

			if c.recusa && err == nil {
				t.Errorf("app.conferir_escopo() aceitou o contexto, mas deveria recusar: %s", c.motivo)
			}
			if !c.recusa && err != nil {
				t.Errorf("app.conferir_escopo() recusou um contexto válido: %v", err)
			}
		})
	}
}

// O trigger das colunas padrão torna empresa_id imutável. Uma linha nunca muda
// de empresa, nem por erro de código: é a última barreira depois do RLS.
func TestEmpresaIDEImutavel(t *testing.T) {
	ctx := context.Background()

	const nome = "tabela_prova_imutabilidade"

	t.Cleanup(func() {
		if _, err := pool.Exec(ctx, "DROP TABLE IF EXISTS app."+nome); err != nil {
			t.Errorf("limpar tabela de prova: %v", err)
		}
	})

	_, err := pool.Exec(ctx, `
		CREATE TABLE app.`+nome+` (
		  id             uuid PRIMARY KEY DEFAULT app.uuid_v7(),
		  empresa_id     uuid NOT NULL,
		  versao         bigint      NOT NULL DEFAULT 1,
		  alterado_xid   xid8        NOT NULL DEFAULT pg_current_xact_id(),
		  criado_em      timestamptz NOT NULL DEFAULT now(),
		  criado_por     uuid,
		  atualizado_em  timestamptz NOT NULL DEFAULT now(),
		  atualizado_por uuid,
		  excluido_em    timestamptz
		)`)
	if err != nil {
		t.Fatalf("criar tabela de prova: %v", err)
	}

	if _, err := pool.Exec(ctx, `
		CREATE TRIGGER trg_padrao BEFORE INSERT OR UPDATE ON app.`+nome+`
		FOR EACH ROW EXECUTE FUNCTION app.tg_padrao()`); err != nil {
		t.Fatalf("criar trigger: %v", err)
	}

	empresaA := uuid.New()
	empresaB := uuid.New()
	var id uuid.UUID
	if err := pool.QueryRow(ctx,
		"INSERT INTO app."+nome+" (empresa_id) VALUES ($1) RETURNING id", empresaA).Scan(&id); err != nil {
		t.Fatalf("inserir: %v", err)
	}

	// A versão nasce em 1.
	var versao int64
	if err := pool.QueryRow(ctx, "SELECT versao FROM app."+nome+" WHERE id = $1", id).Scan(&versao); err != nil {
		t.Fatalf("ler versão: %v", err)
	}
	if versao != 1 {
		t.Errorf("versão inicial = %d, esperado 1", versao)
	}

	// Trocar a empresa da linha é recusado.
	if _, err := pool.Exec(ctx,
		"UPDATE app."+nome+" SET empresa_id = $1 WHERE id = $2", empresaB, id); err == nil {
		t.Error("o trigger aceitou trocar empresa_id; uma linha jamais muda de empresa")
	}

	// Atualização comum incrementa a versão.
	if _, err := pool.Exec(ctx,
		"UPDATE app."+nome+" SET atualizado_por = $1 WHERE id = $2", uuid.New(), id); err != nil {
		t.Fatalf("atualizar: %v", err)
	}
	if err := pool.QueryRow(ctx, "SELECT versao FROM app."+nome+" WHERE id = $1", id).Scan(&versao); err != nil {
		t.Fatalf("ler versão: %v", err)
	}
	if versao != 2 {
		t.Errorf("versão depois do update = %d, esperado 2", versao)
	}
}

// UUID v7 é ordenável no tempo, com granularidade de milissegundo: os 48 bits
// iniciais são o instante, e o resto é aleatório. É o que dá boa localidade no
// índice B-tree e permite o app gerar identificador offline sem mapeamento
// posterior.
//
// A garantia é que o PREFIXO de tempo nunca retrocede — não que dois
// identificadores gerados no mesmo milissegundo saiam em ordem. Exigir ordem
// estrita aqui seria afirmar mais do que a função entrega, e o teste passaria
// ou falharia conforme a velocidade da máquina.
func TestUUIDv7TemPrefixoDeTempoNaoDecrescente(t *testing.T) {
	ctx := context.Background()

	const amostras = 50
	var prefixoAnterior string

	for i := 0; i < amostras; i++ {
		var atual string
		if err := pool.QueryRow(ctx, "SELECT app.uuid_v7()::text").Scan(&atual); err != nil {
			t.Fatalf("gerar uuid: %v", err)
		}

		partes := strings.Split(atual, "-")
		if len(partes) != 5 {
			t.Fatalf("uuid %q malformado", atual)
		}
		// Versão 7: o primeiro dígito hexadecimal do terceiro grupo é 7.
		if !strings.HasPrefix(partes[2], "7") {
			t.Fatalf("uuid %q não é versão 7 (terceiro grupo: %q)", atual, partes[2])
		}

		// Os 48 bits de tempo são os dois primeiros grupos: 8 + 4 dígitos hex.
		prefixo := partes[0] + partes[1]
		if prefixoAnterior != "" && prefixo < prefixoAnterior {
			t.Errorf("prefixo de tempo retrocedeu: %q veio depois de %q", prefixo, prefixoAnterior)
		}
		prefixoAnterior = prefixo
	}

	// E o prefixo tem de ANDAR: se todas as amostras tivessem o mesmo
	// instante, o teste acima passaria sem provar nada sobre tempo.
	//
	// A pausa é feita aqui, entre consultas, e não dentro de uma só. Num
	// `generate_series(...) , pg_sleep(...)` o pg_sleep é avaliado uma única
	// vez para o cruzamento inteiro, então as amostras saem todas no mesmo
	// milissegundo e o resultado depende da velocidade da máquina — teste
	// instável, que é pior que teste ausente.
	prefixos := map[string]bool{}
	for i := 0; i < 3; i += 1 {
		if i > 0 {
			time.Sleep(5 * time.Millisecond)
		}
		var u string
		if err := pool.QueryRow(ctx, "SELECT app.uuid_v7()::text").Scan(&u); err != nil {
			t.Fatalf("gerar uuid: %v", err)
		}
		partes := strings.Split(u, "-")
		prefixos[partes[0]+partes[1]] = true
	}

	if len(prefixos) < 2 {
		t.Errorf("o prefixo de tempo do uuid_v7 não variou em amostras separadas por 5 ms "+
			"(%d prefixo distinto); a ordenação temporal não está funcionando", len(prefixos))
	}
}
