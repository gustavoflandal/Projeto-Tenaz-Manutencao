package autenticacao

import (
	"go/scanner"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// As sete tabelas de credencial (ADR-011; seção 6.9 do documento 5).
var tabelasCredencial = []string{
	"identidade", "acesso_empresa", "sessao", "token_renovacao", "codigo_autorizacao", "desafio_mfa", "token_email",
}

const dirRepo = "internal/repo"

var (
	// citação qualificada: plataforma.sessao
	reCredencialQualificada = regexp.MustCompile(`\bplataforma\s*\.\s*(` + strings.Join(tabelasCredencial, "|") + `)\b`)
	// citação sem esquema, depois de uma palavra que introduz tabela (FROM
	// sessao) ou de uma vírgula da lista do FROM (FROM app.usuario u, sessao s)
	reCredencialSolta = regexp.MustCompile(`(?:\b(?:from|join|into|update|table)\s+|,\s*)(?:only\s+)?(?:lateral\s+)?(` +
		strings.Join(tabelasCredencial, "|") + `)\b`)
)

func citaCredencial(sql string) bool {
	return reCredencialQualificada.MatchString(sql) || reCredencialSolta.MatchString(sql)
}

// semComentarios tira os comentários de linha (--) e de bloco (/* */) do SQL,
// respeitando literais entre aspas simples. Um filtro comentado não pode
// passar por filtro de verdade.
func semComentarios(sql string) string {
	var b strings.Builder
	for i := 0; i < len(sql); i++ {
		switch {
		case sql[i] == '\'':
			j := i + 1
			for j < len(sql) {
				if sql[j] == '\'' {
					if j+1 < len(sql) && sql[j+1] == '\'' {
						j += 2 // aspas dobradas dentro do literal
						continue
					}
					break
				}
				j++
			}
			b.WriteString(sql[i:min(j+1, len(sql))])
			i = j
		case strings.HasPrefix(sql[i:], "--"):
			for i < len(sql) && sql[i] != '\n' {
				i++
			}
			b.WriteByte('\n')
		case strings.HasPrefix(sql[i:], "/*"):
			fim := strings.Index(sql[i+2:], "*/")
			if fim < 0 {
				return b.String()
			}
			i += fim + 3
			b.WriteByte(' ')
		default:
			b.WriteByte(sql[i])
		}
	}
	return b.String()
}

// nivelZero troca o conteúdo entre parênteses por espaços, para que a
// análise veja só o nível de cima da expressão. Devolve também se os
// parênteses estão equilibrados.
func nivelZero(s string) (string, bool) {
	b := []byte(s)
	prof := 0
	for i, c := range b {
		switch c {
		case '(':
			prof++
		case ')':
			prof--
			if prof < 0 {
				return "", false
			}
		default:
			if prof > 0 {
				b[i] = ' '
			}
		}
	}
	return string(b), prof == 0
}

var (
	reFimDoWhere  = regexp.MustCompile(`\b(order\s+by|group\s+by|limit|returning|for\s+update|offset)\b|;`)
	reAnd         = regexp.MustCompile(`\band\b`)
	reOr          = regexp.MustCompile(`\bor\b`)
	reNot         = regexp.MustCompile(`\bnot\b`)
	reIsNot       = regexp.MustCompile(`\bis\s+not\b`)
	reSelect      = regexp.MustCompile(`\bselect\b`)
	reWhere       = regexp.MustCompile(`\bwhere\b`)
	reIgualArg    = regexp.MustCompile(`^(?:(` + reIdent + `)\.)?(` + reIdent + `) = sqlc\.n?arg\((` + reIdent + `)\)(?:::[a-z0-9_]+(?:\[\])?)?$`)
	reFamilia     = regexp.MustCompile(`^(?:(` + reIdent + `)\.)?familia_id = any \(sqlc\.arg\(familias\)::uuid\[\]\)$`)
	chavesExatas  = []string{"token_hash", "codigo_hash", "id", "identidade_id", "email"}
	reEspacos     = regexp.MustCompile(`\s+`)
	reEspacoParen = regexp.MustCompile(`\(\s+|\s+\)`)
	reEspacoPonto = regexp.MustCompile(`\s*\.\s*`)

	reFrom     = regexp.MustCompile(`\bfrom\b`)
	reSet      = regexp.MustCompile(`\bset\b`)
	reUsing    = regexp.MustCompile(`\busing\b`)
	reLateral  = regexp.MustCompile(`\blateral\b`)
	reJoin     = regexp.MustCompile(`\b(?:(?:inner|left|right|full|cross)\s+)?(?:outer\s+)?join\b`)
	reItemFrom = regexp.MustCompile(`^(` + reIdent + `(?:\.` + reIdent + `)?)(?:\s+(?:as\s+)?(` + reIdent + `))?(?:\s+on\s+(.+))?$`)
	reLigacao  = regexp.MustCompile(`^(` + reIdent + `)\.(` + reIdent + `) = (` + reIdent + `)\.(` + reIdent + `)$`)
)

const reIdent = `[a-z_][a-z0-9_]*`

// normalizar põe o SQL em minúsculas, com um espaço entre as palavras, sem
// espaço junto dos parênteses e do ponto, e sem as aspas duplas dos
// identificadores: "plataforma"."sessao" vira plataforma.sessao.
func normalizar(s string) string {
	s = strings.ReplaceAll(strings.ToLower(s), `"`, "")
	s = strings.TrimSpace(reEspacos.ReplaceAllString(s, " "))
	s = reEspacoPonto.ReplaceAllString(s, ".")
	return reEspacoParen.ReplaceAllStringFunc(s, strings.TrimSpace)
}

// ehCredencial diz se o nome de tabela (com ou sem o esquema plataforma) é
// uma das tabelas de credencial.
func ehCredencial(tabela string) bool {
	nome, ok := strings.CutPrefix(tabela, "plataforma.")
	if !ok && strings.Contains(tabela, ".") {
		return false
	}
	return slices.Contains(tabelasCredencial, nome)
}

// itemFrom é uma tabela da lista do FROM (ou o alvo do UPDATE e do DELETE).
type itemFrom struct {
	tabela, alias, on string
	credencial        bool
}

// itensDoFrom lê a lista de tabelas: só tabelas com nome, separadas por JOIN.
// Vírgula (junção implícita), LATERAL, subconsulta, função e USING são
// recusados em consulta de credencial: com eles, a ligação entre o filtro de
// chave e a tabela de credencial deixa de ser verificável.
func itensDoFrom(trecho, trechoTopo string) ([]itemFrom, string) {
	switch {
	case strings.Contains(trechoTopo, ","):
		return nil, "lista do FROM separada por vírgula (junção implícita)"
	case reLateral.MatchString(trecho):
		return nil, "LATERAL no FROM"
	case reSelect.MatchString(trecho):
		return nil, "subconsulta no FROM"
	case reUsing.MatchString(trechoTopo):
		return nil, "JOIN ... USING"
	}
	var pedacos []string
	ini := 0
	for _, m := range reJoin.FindAllStringIndex(trechoTopo, -1) {
		pedacos = append(pedacos, trecho[ini:m[0]])
		ini = m[1]
	}
	pedacos = append(pedacos, trecho[ini:])
	itens := make([]itemFrom, 0, len(pedacos))
	for i, p := range pedacos {
		m := reItemFrom.FindStringSubmatch(strings.TrimSpace(p))
		if m == nil || (i == 0 && m[3] != "") {
			return nil, "item do FROM que não é uma tabela com nome: " + strings.TrimSpace(p)
		}
		alias := m[2]
		if alias == "" {
			alias = m[1][strings.LastIndex(m[1], ".")+1:]
		}
		for _, outro := range itens {
			if outro.alias == alias {
				return nil, "alias repetido no FROM: " + alias
			}
		}
		itens = append(itens, itemFrom{tabela: m[1], alias: alias, on: m[3], credencial: ehCredencial(m[1])})
	}
	return itens, ""
}

// conjuncao corta a expressão nos AND do nível de cima e recusa OR, NOT e
// BETWEEN nesse nível.
func conjuncao(expr string) ([]string, string) {
	topo, ok := nivelZero(expr)
	if !ok {
		return nil, "parênteses desequilibrados"
	}
	switch {
	case reOr.MatchString(topo):
		return nil, "OR no nível de cima"
	case reNot.MatchString(reIsNot.ReplaceAllString(topo, "is")):
		return nil, "NOT no nível de cima"
	case regexp.MustCompile(`\bbetween\b`).MatchString(topo):
		return nil, "BETWEEN (a conjunção não é analisável)"
	}
	var termos []string
	ini := 0
	for _, m := range reAnd.FindAllStringIndex(topo, -1) {
		termos = append(termos, strings.TrimSpace(expr[ini:m[0]]))
		ini = m[1]
	}
	termos = append(termos, strings.TrimSpace(expr[ini:]))
	for _, t := range termos {
		if t == "" {
			return nil, "termo vazio na conjunção"
		}
	}
	return termos, ""
}

// problemaDaConsulta devolve por que a consulta não filtra uma tabela de
// credencial por chave exata, ou "" se ela filtra. A regra (seção 6.9 do
// documento 5):
//
//   - o WHERE é uma conjunção (só AND no nível de cima, sem OR, NOT nem
//     subconsulta);
//   - toda tabela de credencial do FROM tem um termo de chave exata —
//     `<alias>.<col> = sqlc.arg(<col>)` com hash, id, identidade_id, e-mail,
//     o par (empresa_id, usuario_id) ou a família por lista de ids — cujo
//     alias é o dela (sem alias, só com uma tabela no FROM); ou é ligada por
//     JOIN ... ON `<alias>.<chave> = <outro>.<chave>` a uma tabela de
//     credencial já filtrada assim;
//   - todo termo `<col> = sqlc.arg(x)` usa x igual ao nome da coluna, com
//     fronteira de palavra.
//
// Os demais termos só estreitam o resultado, porque a conjunção não alarga.
// INSERT ... SELECT, CTE, UPDATE ... FROM, DELETE ... USING, vírgula, LATERAL
// e subconsulta no FROM, e tabela de credencial fora do FROM são recusados.
func problemaDaConsulta(corpo string) string {
	sql := normalizar(semComentarios(corpo))
	if !citaCredencial(sql) {
		return ""
	}
	topo, ok := nivelZero(sql)
	if !ok {
		return "parênteses desequilibrados"
	}
	primeira := strings.Fields(sql)[0]
	switch primeira {
	case "insert":
		if reSelect.MatchString(sql) {
			return "INSERT ... SELECT em tabela de credencial"
		}
		if strings.Contains(sql, "do update") {
			return "INSERT ... ON CONFLICT DO UPDATE em tabela de credencial"
		}
		return ""
	case "select", "update", "delete":
	default:
		return "comando " + primeira + " (CTE ou outro) em tabela de credencial"
	}
	loc := reWhere.FindStringIndex(topo)
	if loc == nil {
		return "sem WHERE"
	}
	// UPDATE ... FROM: um FROM entre o UPDATE e o WHERE (o "IS DISTINCT FROM"
	// do próprio WHERE não conta).
	if primeira == "update" && regexp.MustCompile(`\bfrom\b`).MatchString(topo[:loc[0]]) {
		return "UPDATE ... FROM em tabela de credencial"
	}
	if len(reWhere.FindAllStringIndex(topo, -1)) > 1 {
		return "mais de um WHERE no nível de cima"
	}
	fim := len(sql)
	if f := reFimDoWhere.FindStringIndex(topo[loc[1]:]); f != nil {
		fim = loc[1] + f[0]
	}
	onde := strings.TrimSpace(sql[loc[1]:fim])
	if reSelect.MatchString(onde) {
		return "subconsulta no WHERE"
	}
	termos, p := conjuncao(onde)
	if p != "" {
		return p + " do WHERE"
	}

	// A lista de tabelas: depois do FROM (SELECT, DELETE) ou entre UPDATE e
	// SET, sempre antes do WHERE.
	var ini, fimFrom int
	switch primeira {
	case "update":
		s := reSet.FindStringIndex(topo[:loc[0]])
		if s == nil {
			return "UPDATE sem SET"
		}
		ini, fimFrom = len("update"), s[0]
	default:
		f := reFrom.FindStringIndex(topo[:loc[0]])
		if f == nil {
			return "sem FROM no nível de cima"
		}
		ini, fimFrom = f[1], loc[0]
	}
	itens, p := itensDoFrom(sql[ini:fimFrom], topo[ini:fimFrom])
	if p != "" {
		return p
	}
	credenciais := 0
	porAlias := map[string]itemFrom{}
	for _, it := range itens {
		porAlias[it.alias] = it
		if it.credencial {
			credenciais++
		}
	}
	if credenciais == 0 {
		return "tabela de credencial fora do FROM (lista de colunas, RETURNING ou outra parte)"
	}

	// Termos de chave do WHERE, ligados ao alias da tabela que filtram.
	filtrada := map[string]bool{}
	empresa, usuario := map[string]bool{}, map[string]bool{}
	aliasDo := func(a string) (string, string) {
		if a == "" {
			if len(itens) != 1 {
				return "", "termo de chave sem alias numa consulta com mais de uma tabela"
			}
			return itens[0].alias, ""
		}
		if _, ok := porAlias[a]; !ok {
			return "", "termo de chave com alias desconhecido: " + a
		}
		return a, ""
	}
	for _, t := range termos {
		if m := reFamilia.FindStringSubmatch(t); m != nil {
			al, p := aliasDo(m[1])
			if p != "" {
				return p
			}
			filtrada[al] = true
			continue
		}
		m := reIgualArg.FindStringSubmatch(t)
		if m == nil {
			continue // termo que só estreita: IS NULL, > now(), < 5, booleano...
		}
		if m[2] != m[3] {
			return "termo " + t + " compara a coluna com argumento de outro nome"
		}
		if strings.Contains(t, "sqlc.narg(") {
			continue // argumento anulável não é chave
		}
		al, p := aliasDo(m[1])
		if p != "" {
			return p
		}
		switch m[2] {
		case "empresa_id":
			empresa[al] = true
		case "usuario_id":
			usuario[al] = true
		default:
			if slices.Contains(chavesExatas, m[2]) {
				filtrada[al] = true
			}
		}
	}
	for al := range empresa {
		if usuario[al] {
			filtrada[al] = true
		}
	}
	// Tabela de credencial ligada por JOIN ... ON a outra já filtrada, pela
	// mesma coluna de chave: AuthUltimaEmpresa junta sessao (filtrada por
	// identidade_id) e acesso_empresa (a.identidade_id = s.identidade_id).
	for mudou := true; mudou; {
		mudou = false
		for _, it := range itens {
			if !it.credencial || filtrada[it.alias] || it.on == "" {
				continue
			}
			ligacoes, p := conjuncao(it.on)
			if p != "" {
				return p + " do ON de " + it.alias
			}
			for _, l := range ligacoes {
				m := reLigacao.FindStringSubmatch(l)
				if m == nil || m[2] != m[4] || !slices.Contains(chavesExatas, m[2]) {
					continue
				}
				if (m[1] == it.alias && m[3] != it.alias && filtrada[m[3]] && porAlias[m[3]].credencial) ||
					(m[3] == it.alias && m[1] != it.alias && filtrada[m[1]] && porAlias[m[1]].credencial) {
					filtrada[it.alias], mudou = true, true
					break
				}
			}
		}
	}
	for _, it := range itens {
		if it.credencial && !filtrada[it.alias] {
			return "a tabela " + it.tabela + " (" + it.alias + ") não é filtrada por chave exata"
		}
	}
	return ""
}

// O próprio analisador recusa o que deve recusar e aceita o que deve aceitar.
func TestAnalisadorDeConsultas(t *testing.T) {
	recusadas := map[string]string{
		"sem where":                   `SELECT id FROM plataforma.identidade`,
		"where comentado":             "SELECT id FROM plataforma.identidade -- WHERE id = sqlc.arg(id)\nWHERE true",
		"chave em bloco comentado":    `SELECT id FROM plataforma.sessao WHERE /* id = sqlc.arg(id) AND */ revogada_em IS NULL`,
		"or no topo":                  `SELECT id FROM plataforma.sessao WHERE id = sqlc.arg(id) OR true`,
		"or em maiúsculas":            `UPDATE plataforma.sessao SET revogada_em = now() WHERE token_hash = sqlc.arg(token_hash) OR 1 = 1`,
		"not no topo":                 `SELECT id FROM plataforma.sessao WHERE NOT id = sqlc.arg(id)`,
		"not negando a chave":         `SELECT id FROM plataforma.sessao WHERE id = sqlc.arg(id) AND NOT (revogada_em IS NULL AND id = sqlc.arg(id))`,
		"chave só dentro de or":       `SELECT id FROM plataforma.sessao WHERE (id = sqlc.arg(id) OR true)`,
		"prefixo de outra coluna":     `SELECT id FROM plataforma.sessao WHERE usuario_xid = sqlc.arg(id)`,
		"argumento de outro nome":     `SELECT id FROM plataforma.identidade WHERE id = sqlc.arg(identidade_id)`,
		"coluna sufixo da chave":      `SELECT id FROM plataforma.sessao WHERE outro_token_hash = sqlc.arg(outro_token_hash)`,
		"só empresa":                  `SELECT id FROM plataforma.sessao WHERE empresa_id = sqlc.arg(empresa_id)`,
		"chave anulável":              `SELECT id FROM plataforma.sessao WHERE id = sqlc.narg(id)`,
		"subconsulta":                 `SELECT id FROM plataforma.sessao WHERE id = sqlc.arg(id) AND usuario_id IN (SELECT usuario_id FROM plataforma.acesso_empresa)`,
		"insert select":               `INSERT INTO app.copia (id) SELECT id FROM plataforma.identidade`,
		"insert select na credencial": `INSERT INTO plataforma.sessao (id) SELECT id FROM app.usuario`,
		"cte":                         `WITH x AS (SELECT id FROM plataforma.identidade) SELECT * FROM x WHERE id = sqlc.arg(id)`,
		"update from":                 `UPDATE plataforma.sessao SET revogada_em = now() FROM app.usuario u WHERE id = sqlc.arg(id)`,
		"between":                     `SELECT id FROM plataforma.sessao WHERE id = sqlc.arg(id) AND criado_em BETWEEN now() AND now()`,
		"tabela sem esquema":          `SELECT id FROM sessao WHERE true`,
		// Achados do tenaz-qa-automacao na revisão dos testes de E01.
		"identificador entre aspas":                   `SELECT senha_hash FROM "plataforma"."sessao"`,
		"aspas com chave em outra tabela":             `SELECT s.* FROM "plataforma" . "sessao" s JOIN app.usuario u ON true WHERE u.id = sqlc.arg(id)`,
		"vírgula no FROM":                             `SELECT s.* FROM app.usuario u, sessao s WHERE u.id = sqlc.arg(id)`,
		"vírgula com esquema":                         `SELECT s.* FROM app.usuario u, plataforma.sessao s WHERE u.id = sqlc.arg(id)`,
		"join sem ligação":                            `SELECT s.* FROM app.usuario u JOIN sessao s ON true WHERE u.id = sqlc.arg(id)`,
		"chave no alias da outra tabela":              `SELECT s.* FROM plataforma.sessao s JOIN app.usuario u ON true WHERE u.id = sqlc.arg(id)`,
		"lateral":                                     `SELECT s.* FROM app.usuario u, LATERAL (SELECT * FROM plataforma.sessao) s WHERE u.id = sqlc.arg(id)`,
		"lateral com join":                            `SELECT s.* FROM app.usuario u JOIN LATERAL (SELECT * FROM plataforma.sessao) s ON true WHERE u.id = sqlc.arg(id)`,
		"subconsulta no from":                         `SELECT x.* FROM (SELECT * FROM plataforma.sessao) x WHERE x.id = sqlc.arg(id)`,
		"chave sem alias com duas tabelas":            `SELECT s.id FROM plataforma.sessao s JOIN app.usuario u ON true WHERE id = sqlc.arg(id)`,
		"alias desconhecido":                          `SELECT s.id FROM plataforma.sessao s WHERE x.id = sqlc.arg(id)`,
		"segunda credencial sem ligação":              `SELECT a.* FROM plataforma.sessao s JOIN plataforma.acesso_empresa a ON a.habilitado WHERE s.id = sqlc.arg(id)`,
		"ligação por coluna que não é chave":          `SELECT a.* FROM plataforma.sessao s JOIN plataforma.acesso_empresa a ON a.empresa_id = s.empresa_id WHERE s.id = sqlc.arg(id)`,
		"ligação com or":                              `SELECT a.* FROM plataforma.sessao s JOIN plataforma.acesso_empresa a ON a.identidade_id = s.identidade_id OR true WHERE s.id = sqlc.arg(id)`,
		"join using":                                  `SELECT a.* FROM plataforma.sessao s JOIN plataforma.acesso_empresa a USING (identidade_id) WHERE s.id = sqlc.arg(id)`,
		"credencial na lista de colunas":              `SELECT (SELECT senha_hash FROM plataforma.identidade LIMIT 1) FROM app.usuario WHERE id = sqlc.arg(id)`,
		"delete using":                                `DELETE FROM plataforma.sessao USING app.usuario u WHERE u.id = sqlc.arg(id)`,
		"par empresa e usuário em aliases diferentes": `SELECT s.id FROM plataforma.sessao s JOIN plataforma.sessao t ON t.id = s.id WHERE s.empresa_id = sqlc.arg(empresa_id) AND t.usuario_id = sqlc.arg(usuario_id)`,
	}
	for nome, sql := range recusadas {
		if problemaDaConsulta(sql) == "" {
			t.Errorf("%s: o analisador aceitou %q", nome, sql)
		}
	}
	aceitas := []string{
		`SELECT id FROM plataforma.identidade WHERE email = sqlc.arg(email)`,
		`UPDATE plataforma.identidade SET x = 1 WHERE id = sqlc.arg(id) AND (mfa_ultimo_passo IS NULL OR mfa_ultimo_passo < sqlc.arg(passo))`,
		`UPDATE plataforma.identidade SET x = 1 WHERE id = sqlc.arg(id) AND mfa_segredo_cifrado IS NOT NULL`,
		`SELECT s.empresa_id FROM plataforma.sessao s WHERE s.identidade_id = sqlc.arg(identidade_id) ORDER BY 1 LIMIT 1`,
		`UPDATE plataforma.sessao SET r = now() WHERE empresa_id = sqlc.arg(empresa_id) AND usuario_id = sqlc.arg(usuario_id) RETURNING id`,
		`UPDATE plataforma.token_renovacao SET r = now() WHERE familia_id = ANY (sqlc.arg(familias)::uuid[]) AND revogado_em IS NULL`,
		`INSERT INTO plataforma.identidade (id) VALUES (sqlc.arg(id)) ON CONFLICT (email) DO NOTHING RETURNING id`,
		`SELECT 1 FROM app.usuario WHERE true`, // não toca credencial
		`SELECT s.empresa_id FROM plataforma.sessao s JOIN plataforma.acesso_empresa a ON a.identidade_id = s.identidade_id AND a.habilitado
		 WHERE s.identidade_id = sqlc.arg(identidade_id) ORDER BY s.criado_em DESC LIMIT 1`,
		`SELECT id FROM "plataforma"."identidade" WHERE "id" = sqlc.arg(id)`,
		`DELETE FROM plataforma.sessao WHERE id = sqlc.arg(id)`,
		`UPDATE plataforma.sessao SET r = now() WHERE identidade_id = sqlc.arg(identidade_id) AND id IS DISTINCT FROM sqlc.narg(exceto)::uuid`,
	}
	for _, sql := range aceitas {
		if p := problemaDaConsulta(sql); p != "" {
			t.Errorf("o analisador recusou %q: %s", sql, p)
		}
	}
	if got := semComentarios("SELECT '--não é comentário' -- é\nFROM x /* bloco */ WHERE y"); strings.Contains(got, "é\n") ||
		strings.Contains(got, "bloco") || !strings.Contains(got, "'--não é comentário'") {
		t.Errorf("semComentarios: %q", got)
	}
}

// Nenhum .sql do backend fora do repositório da autenticação cita as tabelas
// de credencial: o pool da aplicação não as alcança, e uma consulta que as
// cite em outro módulo falharia em produção por falta de privilégio.
func TestCredenciaisSoNoRepositorioDaAutenticacao(t *testing.T) {
	raiz := filepath.Join("..", "..") // backend/internal
	vistos := 0
	err := filepath.WalkDir(raiz, func(caminho string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(caminho, ".sql") {
			return err
		}
		if strings.Contains(filepath.ToSlash(caminho), "plataforma/autenticacao/"+dirRepo+"/") {
			return nil
		}
		b, err := os.ReadFile(caminho)
		if err != nil {
			return err
		}
		vistos++
		if citaCredencial(normalizar(semComentarios(string(b)))) {
			t.Errorf("%s cita uma tabela de credencial, que só o repositório da autenticação lê", caminho)
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if vistos == 0 {
		t.Fatal("nenhum .sql varrido: o caminho da raiz mudou?")
	}
}

// Nenhum literal de texto em código Go de produção — fora do repositório da
// autenticação e do código gerado por ele — cita as tabelas de credencial. Os
// testes ficam de fora: preparam cenário com o pool do migrador.
func TestNenhumSQLLiteralDeCredencialNoGo(t *testing.T) {
	raiz := filepath.Join("..", "..", "..") // backend
	vistos := 0
	err := filepath.WalkDir(raiz, func(caminho string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		barra := filepath.ToSlash(caminho)
		if d.IsDir() {
			if strings.HasSuffix(barra, "plataforma/autenticacao/"+dirRepo) || d.Name() == "node_modules" || d.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}
		if !strings.HasSuffix(caminho, ".go") || strings.HasSuffix(caminho, "_test.go") {
			return nil
		}
		fonte, err := os.ReadFile(caminho)
		if err != nil {
			return err
		}
		vistos++
		fset := token.NewFileSet()
		var s scanner.Scanner
		s.Init(fset.AddFile(caminho, -1, len(fonte)), fonte, nil, 0)
		for {
			pos, tok, lit := s.Scan()
			if tok == token.EOF {
				break
			}
			if tok != token.STRING {
				continue
			}
			texto, err := strconv.Unquote(lit)
			if err != nil {
				texto = lit
			}
			if citaCredencial(normalizar(semComentarios(texto))) {
				t.Errorf("%s: literal cita uma tabela de credencial: %.80q", fset.Position(pos), texto)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if vistos < 50 {
		t.Fatalf("só %d arquivos .go varridos: o caminho da raiz mudou?", vistos)
	}
}

// Toda consulta do repositório que lê ou altera uma tabela de credencial
// filtra por chave exata (problemaDaConsulta). Nem o pacote de autenticação
// enumera credenciais (seção 6.9 do documento 5).
func TestConsultasDoRepositorioFiltramPorChaveExata(t *testing.T) {
	b, err := os.ReadFile(filepath.Join(dirRepo, "autenticacao.sql"))
	if err != nil {
		t.Fatal(err)
	}
	consultas := strings.Split(string(b), "-- name: ")[1:]
	if len(consultas) < 20 {
		t.Fatalf("esperadas as consultas do repositório; achadas %d", len(consultas))
	}
	tocam := 0
	for _, c := range consultas {
		nome := strings.Fields(c)[0]
		corpo := c[strings.Index(c, "\n"):]
		if citaCredencial(normalizar(semComentarios(corpo))) {
			tocam++
		}
		if p := problemaDaConsulta(corpo); p != "" {
			t.Errorf("%s: %s", nome, p)
		}
	}
	if tocam < 20 {
		t.Fatalf("só %d consultas tocam credencial: o analisador deixou de reconhecê-las?", tocam)
	}
}
