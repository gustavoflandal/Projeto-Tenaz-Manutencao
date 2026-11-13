package vazamento

import (
	"context"
	"errors"
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// contexto é o que o helper NaEmpresa fixa em cada transação (seção 2.4 do
// documento 2). Campo vazio não é definido, exatamente como um esquecimento
// no código real — que é o caso que a suíte precisa cobrir.
type contexto struct {
	Empresa    string
	Usuario    string
	Escopo     string
	Fornecedor string
	Locatario  string
	LinkHash   string
	Conexao    string // só o sincronizador: app.conexao_empresa_id
}

func (c contexto) variaveis() [][2]string {
	var v [][2]string
	for _, par := range [][2]string{
		{"app.empresa_id", c.Empresa},
		{"app.usuario_id", c.Usuario},
		{"app.escopo", c.Escopo},
		{"app.fornecedor_id", c.Fornecedor},
		{"app.locatario_id", c.Locatario},
		{"app.link_hash", c.LinkHash},
		{"app.conexao_empresa_id", c.Conexao},
	} {
		if par[1] != "" {
			v = append(v, par)
		}
	}
	return v
}

// Papéis de banco que atendem a aplicação. O pool da suíte conecta como o
// migrador (superusuário no contêiner) e assume o papel com SET LOCAL ROLE,
// que vale só dentro da transação.
const (
	papelApp     = "tenaz_app"
	papelJob     = "tenaz_job"
	papelLeitura = "tenaz_leitura"
	papelDono    = "tenaz_dono"
	papelSinc    = "tenaz_sinc"

	// Credenciais de autenticação (ADR-011): o papel e o usuário de conexão.
	papelAutenticacao        = "tenaz_autenticacao"
	papelAutenticacaoConexao = "tenaz_autenticacao_conexao"

	// Dono das tabelas de credencial (ISO-E01-08): NOLOGIN, sem membro de
	// runtime. Nem o dono nem a manutenção alcançam as credenciais.
	papelCredencial = "tenaz_credencial"
	papelManutencao = "tenaz_manutencao"
)

// sqlPapeisDaAplicacao é a lista, em SQL, de todo papel que atende a
// aplicação e não é membro do dono. Toda regra de catálogo "nenhum papel da
// aplicação..." usa esta lista, para que um papel novo entre em todas.
const sqlPapeisDaAplicacao = `ARRAY['tenaz_app', 'tenaz_job', 'tenaz_leitura', 'tenaz_sinc', 'tenaz_sinc_conexao',
                                   'tenaz_autenticacao', 'tenaz_autenticacao_conexao']`

// papeisDeExecucao são todos os papéis que se conectam ou que uma conexão
// herda, inclusive tenaz_manutencao, que é membro do dono por desenho. As
// regras de atributo (BYPASSRLS, superusuário) e de relação de membro valem
// para todos eles.
var papeisDeExecucao = []string{"tenaz_app", "tenaz_job", "tenaz_leitura", "tenaz_sinc", "tenaz_sinc_conexao",
	"tenaz_autenticacao", "tenaz_autenticacao_conexao", "tenaz_manutencao"}

// sqlPapeisDeExecucaoTexto é papeisDeExecucao em SQL.
const sqlPapeisDeExecucaoTexto = `ARRAY['tenaz_app', 'tenaz_job', 'tenaz_leitura', 'tenaz_sinc', 'tenaz_sinc_conexao',
                                       'tenaz_autenticacao', 'tenaz_autenticacao_conexao', 'tenaz_manutencao']`

// consultor é o que os verificadores de catálogo precisam: o pool, para o
// teste real, ou uma transação, para a prova do detector, que estraga o banco
// e desfaz tudo no fim.
type consultor interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

// comoPapel abre uma transação, assume o papel, fixa o contexto e executa fn.
// A transação é SEMPRE desfeita: nenhum teste deixa rastro nas fixtures, e um
// teste que apague ou altere dado não contamina o seguinte.
func comoPapel(ctx context.Context, papel string, c contexto, fn func(tx pgx.Tx) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("abrir transação: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if err := prepararTransacao(ctx, tx, papel, c); err != nil {
		return err
	}
	return fn(tx)
}

// confirmarComoPapel é comoPapel com COMMIT. Só as fixtures usam.
func confirmarComoPapel(ctx context.Context, papel string, c contexto, fn func(tx pgx.Tx) error) error {
	tx, err := pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("abrir transação: %w", err)
	}
	defer func() { _ = tx.Rollback(ctx) }()

	if err := prepararTransacao(ctx, tx, papel, c); err != nil {
		return err
	}
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

// erroPreparacao é a falha ao montar a transação de teste (assumir o papel,
// fixar o contexto). Ela NÃO carrega o SQLSTATE de propósito: "permission
// denied to set role" também é 42501, e contaria como "o banco recusou o
// ataque" quando o ataque nem chegou a rodar.
type erroPreparacao struct{ msg string }

func (e erroPreparacao) Error() string { return "preparação da transação de teste: " + e.msg }

func prepararTransacao(ctx context.Context, tx pgx.Tx, papel string, c contexto) error {
	if papel != "" {
		if _, err := tx.Exec(ctx, "SET LOCAL ROLE "+pgx.Identifier{papel}.Sanitize()); err != nil {
			return erroPreparacao{fmt.Sprintf("assumir %s: %v", papel, err)}
		}
	}
	for _, v := range c.variaveis() {
		if _, err := tx.Exec(ctx, `SELECT set_config($1, $2, true)`, v[0], v[1]); err != nil {
			return erroPreparacao{fmt.Sprintf("definir %s: %v", v[0], err)}
		}
	}
	return nil
}

// sqlstate devolve o código do erro do PostgreSQL, ou "" se o erro não veio do banco.
func sqlstate(err error) string {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		return pgErr.Code
	}
	return ""
}

// Códigos usados nas asserções.
const (
	erroPrivilegio = "42501" // insufficient_privilege: policy, GRANT ou trigger de segurança
	erroFK         = "23503" // foreign_key_violation
	erroCheck      = "23514" // check_violation, inclusive dos triggers de perfil
)

// contar executa um count(*) e devolve o número, ou o erro.
func contar(ctx context.Context, tx pgx.Tx, sql string, args ...any) (int64, error) {
	var n int64
	err := tx.QueryRow(ctx, sql, args...).Scan(&n)
	return n, err
}

// alterou executa um comando de escrita e devolve quantas linhas mudaram.
func alterou(ctx context.Context, tx pgx.Tx, sql string, args ...any) (int64, error) {
	tag, err := tx.Exec(ctx, sql, args...)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

// naoMudouNada é o critério de "escrita cruzada recusada" para UPDATE e
// DELETE: zero linhas afetadas, ou erro de privilégio. Qualquer outro erro é
// falha do teste, e não aprovação silenciosa: um erro de sintaxe passaria por
// "recusado" e esconderia o vazamento.
func naoMudouNada(t *testing.T, rotulo string, linhas int64, err error) {
	t.Helper()
	switch {
	case err == nil && linhas == 0:
	case err == nil:
		t.Errorf("%s: alterou %d linha(s); deveria não alcançar nenhuma", rotulo, linhas)
	case sqlstate(err) == erroPrivilegio:
	default:
		t.Errorf("%s: erro inesperado (o teste não prova a recusa): %v", rotulo, err)
	}
}

// exigirRLS confere que a recusa veio de uma POLICY, e não de outro 42501 (falta
// de GRANT, trigger): quando o alvo do teste é a policy, recusa por outro
// motivo esconderia a policy ausente.
func exigirRLS(t *testing.T, rotulo string, err error) {
	t.Helper()
	if !recusaDeRLS(err) {
		t.Errorf("%s: esperado recusa de policy (42501, row-level security); veio %v", rotulo, err)
	}
}

func recusaDeRLS(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) && pgErr.Code == erroPrivilegio &&
		strings.Contains(pgErr.Message, "row-level security")
}

// exigirErro confere que o comando falhou com o código esperado.
func exigirErro(t *testing.T, rotulo string, err error, codigo string) {
	t.Helper()
	if err == nil {
		t.Errorf("%s: o banco aceitou; deveria recusar com %s", rotulo, codigo)
		return
	}
	if got := sqlstate(err); got != codigo {
		t.Errorf("%s: recusou com %q (%v); esperado %s", rotulo, got, err, codigo)
	}
}

// ------------------------------------------------------------------ achados

// achadosAbertos são falhas de isolamento JÁ demonstradas por um teste desta
// suíte e devolvidas ao dono da migração (tenaz-banco-dados), com o ID do
// achado da revisão de isolamento da etapa.
//
// Isto NÃO é lista de exceções. Enquanto houver item aqui, o portão de
// isolamento da etapa está VERMELHO:
//   - no CI (variável CI definida), o teste do achado FALHA, e o TestMain
//     encerra com erro se o mapa não estiver vazio;
//   - fora do CI, o teste aparece como SKIP com o ID, para o desenvolvedor
//     poder trabalhar no resto;
//   - chave do mapa que nenhum teste exercitou numa execução completa é erro:
//     achado sem prova não fica aqui;
//   - quando o defeito some, o teste falha pedindo que o item saia.
var achadosAbertos = map[string]string{}

var (
	achadosMu        sync.Mutex
	achadosExercidos = map[string]bool{}
)

const mensagemPortaoCI = "portão de isolamento vermelho: há achados abertos e a suíte roda no CI"

// achadosSemPortaCI é verdadeiro fora do CI (variável CI vazia).
func achadosSemPortaCI() bool { return os.Getenv("CI") == "" }

// confirmarAchado aplica a regra de achadosAbertos. violou diz se o defeito
// apareceu nesta execução.
func confirmarAchado(t *testing.T, id string, violou bool, detalhe string) {
	t.Helper()
	achadosMu.Lock()
	achadosExercidos[id] = true
	achadosMu.Unlock()

	descricao, aberto := achadosAbertos[id]
	switch {
	case violou && aberto && achadosSemPortaCI():
		t.Skipf("ACHADO ABERTO %s (portão de isolamento vermelho): %s — %s", id, descricao, detalhe)
	case violou && aberto:
		t.Errorf("ACHADO ABERTO %s, %s: %s — %s", id, mensagemPortaoCI, descricao, detalhe)
	case violou && !aberto:
		t.Errorf("%s: %s", id, detalhe)
	case !violou && aberto:
		t.Errorf("o achado %s não se reproduz mais (%s); confirme a correção na revisão de isolamento "+
			"e retire-o de achadosAbertos", id, descricao)
	}
}

// problemasDosAchados é chamado pelo TestMain depois de todos os testes.
// execucaoCompleta é falsa quando a suíte rodou com -run: aí uma chave não
// exercitada pode ser só um teste filtrado.
func problemasDosAchados(execucaoCompleta bool) []string {
	achadosMu.Lock()
	defer achadosMu.Unlock()
	var p []string
	for id := range achadosAbertos {
		if execucaoCompleta && !achadosExercidos[id] {
			p = append(p, "achado "+id+" está em achadosAbertos e nenhum teste o exercitou; achado sem prova não fica na lista")
		}
		if !achadosSemPortaCI() {
			p = append(p, "achado "+id+" aberto: "+mensagemPortaoCI)
		}
	}
	sort.Strings(p)
	return p
}
