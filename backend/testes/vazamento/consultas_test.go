package vazamento

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"testing"
)

// Toda consulta sqlc sobre tabela de cliente repete o filtro da empresa
// (seção 2.4 do documento 2; seção 9.3 do documento 5, "Consultas sqlc"). O
// RLS é a primeira camada; o filtro na consulta é a segunda, e é também o que
// deixa o planejador usar o índice que começa por empresa_id.
//
// A suíte só LÊ os .sql de backend/internal; não os altera.
//
// Regra, por consulta (bloco "-- name:"), depois de retirar os comentários
// (-- e /* */) e os literais de texto:
//   - nome entre aspas ("app"."pessoa") e tabela de cliente sem esquema
//     (FROM pessoa) são acusados: escapam da leitura da regra;
//   - INSERT em tabela de app com empresa_id: empresa_id na lista de colunas e
//     o parâmetro sqlc.arg(empresa_id);
//   - SELECT, UPDATE e DELETE: todo WHERE do nível de cima (fora de
//     parênteses: subconsulta não conta) tem, como conjunção AND, empresa_id
//     (com ou sem alias) igual ao parâmetro empresa_id, e nenhum OR no nível de
//     cima (OR true anularia o filtro);
//   - app.empresa: o mesmo, com id igual ao parâmetro id (ou empresa_id).
//
// Limite conhecido: a regra é por consulta, não por tabela da junção. Uma
// junção de duas tabelas de cliente com um só filtro passa; o RLS segura a
// outra. Fica registrado como pendência da suíte.

const raizConsultas = "../../internal"

var (
	reNomeConsulta = regexp.MustCompile(`(?m)^--\s*name:\s*(\S+)`)
	reTabelaApp    = regexp.MustCompile(`\bapp\.([a-z_][a-z0-9_]*)\b\s*(\()?`)
	reSemEsquema   = regexp.MustCompile(`\b(from|join|update|into)\s+([a-z_][a-z0-9_]*)\s*(\.)?`)
	reInsert       = regexp.MustCompile(`^\s*insert\s+into\s+app\.([a-z_][a-z0-9_]*)\s*\(([^)]*)\)`)
	reWhere        = regexp.MustCompile(`\bwhere\b`)
	reFimDoWhere   = regexp.MustCompile(`\b(order\s+by|group\s+by|having|limit|offset|returning|for\s+(update|share|no\s+key|key)|union|intersect|except|window|on\s+conflict)\b|;`)
	reOuNoTopo     = regexp.MustCompile(`\bor\b`)
	reFiltroEmp    = regexp.MustCompile(`^\(*\s*([a-z_][a-z0-9_]*\.)?empresa_id\s*=\s*(sqlc\.n?arg\(empresa_id\)|@empresa_id)\s*\)*$`)
	reFiltroID     = regexp.MustCompile(`^\(*\s*([a-z_][a-z0-9_]*\.)?id\s*=\s*(sqlc\.n?arg\((id|empresa_id)\)|@id|@empresa_id)\s*\)*$`)
)

// limparSQL tira comentários e o conteúdo dos literais de texto (que viram
// ”), e devolve em minúsculas. Um filtro dentro de comentário ou de texto não
// é filtro.
func limparSQL(sql string) string {
	var b strings.Builder
	for i := 0; i < len(sql); i++ {
		switch {
		case strings.HasPrefix(sql[i:], "--"):
			for i < len(sql) && sql[i] != '\n' {
				i++
			}
			b.WriteByte('\n')
		case strings.HasPrefix(sql[i:], "/*"):
			profundidade := 0
			for ; i < len(sql); i++ {
				if strings.HasPrefix(sql[i:], "/*") {
					profundidade++
					i++
				} else if strings.HasPrefix(sql[i:], "*/") {
					profundidade--
					i++
					if profundidade == 0 {
						break
					}
				}
			}
			b.WriteByte(' ')
		case sql[i] == '\'':
			i++
			for i < len(sql) && !(sql[i] == '\'' && (i+1 >= len(sql) || sql[i+1] != '\'')) {
				if sql[i] == '\'' {
					i++ // aspas dobradas dentro do literal
				}
				i++
			}
			b.WriteString("''")
		default:
			b.WriteByte(sql[i])
		}
	}
	return strings.ToLower(b.String())
}

// profundidades devolve, para cada byte, quantos parênteses estão abertos.
func profundidades(sql string) []int {
	p := make([]int, len(sql))
	d := 0
	for i := 0; i < len(sql); i++ {
		if sql[i] == '(' {
			d++
		}
		p[i] = d
		if sql[i] == ')' {
			d--
		}
	}
	return p
}

// noTopo devolve as ocorrências de re fora de parênteses.
func noTopo(re *regexp.Regexp, sql string, prof []int, desde int) [][]int {
	var saida [][]int
	for _, m := range re.FindAllStringIndex(sql[desde:], -1) {
		if prof[desde+m[0]] == 0 {
			saida = append(saida, []int{desde + m[0], desde + m[1]})
		}
	}
	return saida
}

// clausulasWhere devolve o texto de cada WHERE do nível de cima.
func clausulasWhere(sql string) []string {
	prof := profundidades(sql)
	var clausulas []string
	for _, w := range noTopo(reWhere, sql, prof, 0) {
		fim := len(sql)
		if f := noTopo(reFimDoWhere, sql, prof, w[1]); len(f) > 0 {
			fim = f[0][0]
		}
		clausulas = append(clausulas, sql[w[1]:fim])
	}
	return clausulas
}

// whereFiltra diz se a cláusula tem, como conjunção AND do nível de cima, um
// termo que casa com filtro, e nenhum OR no nível de cima.
func whereFiltra(clausula string, filtro *regexp.Regexp) (bool, string) {
	prof := profundidades(clausula)
	if len(noTopo(reOuNoTopo, clausula, prof, 0)) > 0 {
		return false, "OR no nível de cima do WHERE"
	}
	inicio := 0
	termos := []string{}
	for _, e := range noTopo(regexp.MustCompile(`\band\b`), clausula, prof, 0) {
		termos = append(termos, clausula[inicio:e[0]])
		inicio = e[1]
	}
	termos = append(termos, clausula[inicio:])
	for _, termo := range termos {
		if filtro.MatchString(strings.TrimSpace(termo)) {
			return true, ""
		}
	}
	return false, "sem o filtro no nível de cima do WHERE"
}

// consultasSemFiltro devolve "arquivo:consulta (motivo)" de toda consulta que
// cita tabela de cliente sem o filtro. clientes: tabela -> tem empresa_id
// (falso só em app.empresa).
func consultasSemFiltro(arquivo, texto string, clientes map[string]bool) []string {
	var achados []string
	acusar := func(nome, motivo string) { achados = append(achados, arquivo+":"+nome+" ("+motivo+")") }

	posicoes := reNomeConsulta.FindAllStringSubmatchIndex(texto, -1)
	for i, p := range posicoes {
		fim := len(texto)
		if i+1 < len(posicoes) {
			fim = posicoes[i+1][0]
		}
		nome := texto[p[2]:p[3]]
		sql := limparSQL(texto[p[1]:fim])

		if strings.Contains(sql, `"`) {
			acusar(nome, "nome entre aspas")
			continue
		}
		for _, m := range reSemEsquema.FindAllStringSubmatch(sql, -1) {
			if _, cliente := clientes[m[2]]; cliente && m[3] == "" {
				acusar(nome, "tabela de cliente sem esquema: "+m[2])
			}
		}

		citaCliente, citaSoEmpresa := false, true
		for _, m := range reTabelaApp.FindAllStringSubmatch(sql, -1) {
			if m[2] != "" { // função
				continue
			}
			if temEmpresa, cliente := clientes[m[1]]; cliente {
				citaCliente = true
				citaSoEmpresa = citaSoEmpresa && !temEmpresa
			}
		}
		if !citaCliente {
			continue
		}

		if ins := reInsert.FindStringSubmatch(sql); ins != nil {
			if temEmpresa, cliente := clientes[ins[1]]; cliente {
				coluna, parametro := "empresa_id", "sqlc.arg(empresa_id)"
				if !temEmpresa {
					coluna, parametro = "id", "sqlc.arg(id)"
				}
				if !regexp.MustCompile(`\b`+coluna+`\b`).MatchString(ins[2]) || !strings.Contains(sql, parametro) {
					acusar(nome, "INSERT sem "+coluna+" = "+parametro)
				}
				continue
			}
		}

		filtro := reFiltroEmp
		if citaSoEmpresa {
			filtro = reFiltroID
		}
		clausulas := clausulasWhere(sql)
		if len(clausulas) == 0 {
			acusar(nome, "sem WHERE no nível de cima")
			continue
		}
		for _, c := range clausulas {
			if ok, motivo := whereFiltra(c, filtro); !ok {
				acusar(nome, motivo)
				break
			}
		}
	}
	return achados
}

func tabelasDeCliente(t *testing.T, ctx context.Context) map[string]bool {
	t.Helper()
	clientes := map[string]bool{}
	for _, tb := range tabelasDeEmpresa(t, ctx) {
		clientes[strings.TrimPrefix(tb.Nome, "app.")] = tb.Coluna == "empresa_id"
	}
	return clientes
}

func TestConsultasSqlcFiltramPorEmpresa(t *testing.T) {
	ctx := context.Background()
	clientes := tabelasDeCliente(t, ctx)

	var arquivos []string
	err := filepath.WalkDir(raizConsultas, func(caminho string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() && strings.HasSuffix(caminho, ".sql") {
			arquivos = append(arquivos, caminho)
		}
		return nil
	})
	if err != nil {
		t.Fatalf("percorrer %s: %v", raizConsultas, err)
	}
	sort.Strings(arquivos)
	if len(arquivos) == 0 {
		t.Fatalf("nenhum .sql em %s: a regra passaria por vacuidade", raizConsultas)
	}
	consultas := 0
	for _, a := range arquivos {
		b, err := os.ReadFile(a)
		if err != nil {
			t.Fatalf("ler %s: %v", a, err)
		}
		consultas += len(reNomeConsulta.FindAllString(string(b), -1))
		for _, achado := range consultasSemFiltro(filepath.ToSlash(a), string(b), clientes) {
			t.Errorf("consulta sobre tabela de cliente sem filtro de empresa: %s; repita "+
				"empresa_id = sqlc.arg(empresa_id) como AND no nível de cima (seção 2.4 do documento 2)", achado)
		}
	}
	if consultas == 0 {
		t.Fatal("nenhuma consulta sqlc encontrada: a regra passaria por vacuidade")
	}
}

// As sete consultas montadas pelo tenaz-qa-automacao na 3ª revisão (P-QA-2),
// mais as corretas que a regra precisa deixar passar.
func TestDetectorAcusaConsultaSemFiltroDeEmpresa(t *testing.T) {
	clientes := map[string]bool{"pessoa": true, "usuario": true, "unidade": true, "empresa": false}
	texto := `
-- name: OrTrue :many
SELECT cpf FROM app.pessoa WHERE empresa_id = sqlc.arg(empresa_id) OR true;

-- name: ComentarioNaLinha :many
SELECT cpf FROM app.pessoa /* todas */ ORDER BY nome; -- empresa_id = sqlc.arg(empresa_id)

-- name: Bloco :many
SELECT cpf FROM app.pessoa /* WHERE empresa_id = sqlc.arg(empresa_id) */;

-- name: SoNaSubconsulta :many
SELECT p.cpf FROM app.pessoa p WHERE p.id IN (SELECT pessoa_id FROM app.usuario WHERE empresa_id = sqlc.arg(empresa_id));

-- name: SemEsquema :many
SELECT cpf FROM pessoa;

-- name: ParametroErrado :many
SELECT cpf FROM app.pessoa WHERE empresa_id = sqlc.arg(fornecedor_id);

-- name: Aspas :many
SELECT cpf FROM "app"."pessoa";

-- name: FiltroEmTexto :many
SELECT cpf FROM app.pessoa WHERE nome = 'empresa_id = sqlc.arg(empresa_id)';

-- name: ListarEmpresas :many
SELECT id, razao_social FROM app.empresa;

-- name: Correta :many
SELECT cpf FROM app.pessoa WHERE empresa_id = sqlc.arg(empresa_id);

-- name: CorretaComAliasEJuncao :many
SELECT u.nome FROM app.usuario u JOIN app.pessoa p ON p.empresa_id = u.empresa_id AND p.id = u.pessoa_id
WHERE (u.empresa_id = sqlc.arg(empresa_id)) AND (u.situacao = 'ativo' OR u.situacao = 'convidado')
ORDER BY u.nome;

-- name: CorretaInsert :exec
INSERT INTO app.unidade (empresa_id, codigo, nome) VALUES (sqlc.arg(empresa_id), sqlc.arg(codigo), sqlc.arg(nome));

-- name: CorretaEmpresa :one
SELECT razao_social FROM app.empresa WHERE id = sqlc.arg(id) AND excluido_em IS NULL;

-- name: Contexto :exec
SELECT app.conferir_escopo();
`
	got := consultasSemFiltro("prova.sql", texto, clientes)
	acusadas := map[string]bool{}
	for _, a := range got {
		acusadas[strings.SplitN(strings.TrimPrefix(a, "prova.sql:"), " ", 2)[0]] = true
	}
	for _, nome := range []string{"OrTrue", "ComentarioNaLinha", "Bloco", "SoNaSubconsulta", "SemEsquema",
		"ParametroErrado", "Aspas", "FiltroEmTexto", "ListarEmpresas"} {
		if !acusadas[nome] {
			t.Errorf("a consulta %s, sem filtro de empresa de verdade, não foi acusada; acusadas: %v", nome, got)
		}
	}
	for _, nome := range []string{"Correta", "CorretaComAliasEJuncao", "CorretaInsert", "CorretaEmpresa", "Contexto"} {
		if acusadas[nome] {
			t.Errorf("a consulta correta %s foi acusada: %v", nome, got)
		}
	}
}
