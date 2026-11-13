package vazamento

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"sync"
)

// O que a parte 3 lê do banco, sempre pelo pool do migrador (superusuário,
// que enxerga todas as empresas): de quem é cada UUID, e o retrato de B antes
// e depois da varredura. Tudo é calculado do catálogo, tabela a tabela, para
// que tabela nova (a 00007 da auditoria configurável, por exemplo) entre sem
// lista à mão.

var reUUID = regexp.MustCompile(`[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`)

// uuidsNoTexto devolve os UUIDs do texto, em minúsculas e sem repetição.
func uuidsNoTexto(s string) []string {
	vistos := map[string]bool{}
	var l []string
	for _, u := range reUUID.FindAllString(s, -1) {
		u = strings.ToLower(u)
		if !vistos[u] {
			vistos[u] = true
			l = append(l, u)
		}
	}
	return l
}

// tabelasComColuna lista as tabelas (o pai, nas particionadas) que têm a
// coluna, em todo esquema que não seja do sistema.
func tabelasComColuna(ctx context.Context, c consultor, coluna string) ([]string, error) {
	rows, err := c.Query(ctx, `
		SELECT format('%I.%I', n.nspname, cl.relname)
		FROM pg_class cl
		JOIN pg_namespace n ON n.oid = cl.relnamespace
		JOIN pg_attribute a ON a.attrelid = cl.oid AND a.attname = $1 AND NOT a.attisdropped
		WHERE cl.relkind IN ('r', 'p') AND NOT cl.relispartition
		  AND n.nspname NOT IN ('pg_catalog', 'information_schema') AND n.nspname NOT LIKE 'pg\_%'
		ORDER BY 1`, coluna)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var l []string
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return nil, err
		}
		l = append(l, t)
	}
	return l, rows.Err()
}

// ------------------------------------------------------------ classificador

// classificador sabe, para cada UUID que aparece no texto de alguma linha do
// banco, em linhas de quais empresas ele aparece. "Aparece no texto" pega id,
// chave estrangeira, request_id da trilha e valor de antes/depois: qualquer
// UUID que só aparece em linha de outra empresa é dado de outra empresa.
// Fica de fora o detalhe do evento da trilha: ele guarda os filtros que o
// próprio usuário mandou (a exportação de A filtrada pela pessoa de B grava o
// id de B na trilha de A), e eco de pedido não dá posse do UUID.
type classificador struct {
	mu    sync.Mutex
	donos map[string]map[string]bool // uuid -> empresas
}

func novoClassificador(ctx context.Context) (*classificador, error) {
	c := &classificador{}
	return c, c.carregar(ctx)
}

func (c *classificador) carregar(ctx context.Context) error {
	tabelas, err := tabelasComColuna(ctx, pool, "empresa_id")
	if err != nil {
		return err
	}
	donos := map[string]map[string]bool{}
	ler := func(sql string) error {
		rows, err := pool.Query(ctx, sql)
		if err != nil {
			return fmt.Errorf("%s: %w", sql, err)
		}
		defer rows.Close()
		for rows.Next() {
			var empresa, linha string
			if err := rows.Scan(&empresa, &linha); err != nil {
				return err
			}
			for _, u := range uuidsNoTexto(linha) {
				if donos[u] == nil {
					donos[u] = map[string]bool{}
				}
				donos[u][empresa] = true
			}
		}
		return rows.Err()
	}
	for _, t := range tabelas {
		if err := ler(`SELECT empresa_id::text, (to_jsonb(t) - 'detalhe')::text FROM ` + t + ` t WHERE empresa_id IS NOT NULL`); err != nil {
			return err
		}
	}
	if err := ler(`SELECT id::text, t::text FROM app.empresa t`); err != nil {
		return err
	}
	c.mu.Lock()
	c.donos = donos
	c.mu.Unlock()
	return nil
}

// atorAPI é quem faz a requisição, para decidir o que lhe é alheio.
type atorAPI struct {
	rotulo  string
	empresa string
	externo bool
	// proprios: para o externo, os UUIDs do próprio alcance (o resto do
	// banco, inclusive a própria empresa, é alheio).
	proprios map[string]bool
	// textos que não podem aparecer em resposta nenhuma deste ator.
	textos []string
}

// estado diz se o UUID é alheio ao ator e se ele é conhecido no banco.
func (c *classificador) estado(a atorAPI, u string) (alheio, conhecido bool) {
	c.mu.Lock()
	donos := c.donos[u]
	c.mu.Unlock()
	if a.externo && a.proprios[u] {
		return false, true
	}
	if len(donos) == 0 {
		return false, false
	}
	if a.externo {
		return true, true
	}
	return !donos[a.empresa], true
}

// alheios filtra os UUIDs alheios ao ator. UUID desconhecido ou alheio
// provoca uma releitura do banco (linha nova da própria empresa, criada na
// varredura) antes da decisão.
func (c *classificador) alheios(ctx context.Context, a atorAPI, uuids []string) ([]string, error) {
	recarregado := false
	var fora []string
	for _, u := range uuids {
		alheio, conhecido := c.estado(a, u)
		if (alheio || !conhecido) && !recarregado {
			if err := c.carregar(ctx); err != nil {
				return nil, err
			}
			recarregado = true
			alheio, _ = c.estado(a, u)
		}
		if alheio {
			fora = append(fora, u)
		}
	}
	return fora, nil
}

// propriosDoExterno: a empresa, o usuário, a pessoa, a identidade, o vínculo
// (fornecedor ou locatário), a unidade padrão, os perfis atribuídos e as
// linhas com usuario_id do próprio usuário (sessões, atribuições, eventos de
// que ele é o ator).
func propriosDoExterno(ctx context.Context, x *externoAPI) (map[string]bool, error) {
	p := map[string]bool{strings.ToLower(x.empresa.ID): true, strings.ToLower(x.usuario): true}
	var linha string
	if err := pool.QueryRow(ctx, `SELECT to_jsonb(u)::text FROM app.usuario u WHERE id = $1`, x.usuario).Scan(&linha); err != nil {
		return nil, err
	}
	var u map[string]any
	if err := json.Unmarshal([]byte(linha), &u); err != nil {
		return nil, err
	}
	for _, k := range []string{"pessoa_id", "identidade_id", "fornecedor_id", "locatario_id", "unidade_padrao_id"} {
		if s, ok := u[k].(string); ok {
			p[strings.ToLower(s)] = true
		}
	}
	comUsuario, err := tabelasComColuna(ctx, pool, "usuario_id")
	if err != nil {
		return nil, err
	}
	comID, err := tabelasComColuna(ctx, pool, "id")
	if err != nil {
		return nil, err
	}
	temID := map[string]bool{}
	for _, t := range comID {
		temID[t] = true
	}
	consultas := []string{`SELECT perfil_id::text FROM app.usuario_perfil WHERE usuario_id = $1`}
	for _, t := range comUsuario {
		if temID[t] {
			consultas = append(consultas, `SELECT id::text FROM `+t+` WHERE usuario_id = $1`)
		}
	}
	for _, sql := range consultas {
		rows, err := pool.Query(ctx, sql, x.usuario)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", sql, err)
		}
		for rows.Next() {
			var id string
			if err := rows.Scan(&id); err != nil {
				rows.Close()
				return nil, err
			}
			p[strings.ToLower(id)] = true
		}
		rows.Close()
		if err := rows.Err(); err != nil {
			return nil, err
		}
	}
	return p, nil
}

// --------------------------------------------------------------- retrato

// retratoDaEmpresa resume todas as linhas da empresa, tabela a tabela:
// contagem e md5 do texto das linhas em ordem. Entram toda tabela com
// empresa_id, app.empresa, as identidades dos usuários da empresa, as tabelas
// por identidade sem empresa_id (desafio de MFA) e os refresh tokens das
// sessões da empresa. Uma escrita que alcance a empresa muda o retrato.
func retratoDaEmpresa(ctx context.Context, c consultor, empresa string) (map[string]string, error) {
	const resumo = `SELECT count(*)::text || ':' || coalesce(md5(string_agg(t::text, E'\n' ORDER BY t::text)), '') FROM `
	comEmpresa, err := tabelasComColuna(ctx, c, "empresa_id")
	if err != nil {
		return nil, err
	}
	comIdentidade, err := tabelasComColuna(ctx, c, "identidade_id")
	if err != nil {
		return nil, err
	}
	consultas := map[string]string{
		"app.empresa":                resumo + `app.empresa t WHERE id = $1`,
		"plataforma.identidade":      resumo + `plataforma.identidade t WHERE id IN (SELECT identidade_id FROM app.usuario WHERE empresa_id = $1)`,
		"plataforma.token_renovacao": resumo + `plataforma.token_renovacao t WHERE familia_id IN (SELECT id FROM plataforma.sessao WHERE empresa_id = $1)`,
	}
	temEmpresa := map[string]bool{}
	for _, t := range comEmpresa {
		temEmpresa[t] = true
		consultas[t] = resumo + t + ` t WHERE empresa_id = $1`
	}
	for _, t := range comIdentidade {
		if !temEmpresa[t] {
			consultas[t] = resumo + t + ` t WHERE identidade_id IN (SELECT identidade_id FROM app.usuario WHERE empresa_id = $1)`
		}
	}
	r := map[string]string{}
	for t, sql := range consultas {
		var v string
		if err := c.QueryRow(ctx, sql, empresa).Scan(&v); err != nil {
			return nil, fmt.Errorf("retrato de %s: %w", t, err)
		}
		r[t] = v
	}
	return r, nil
}

// diferencasDoRetrato lista as tabelas cujas linhas da empresa mudaram.
func diferencasDoRetrato(antes, depois map[string]string) []string {
	var p []string
	for t, a := range antes {
		if d, ok := depois[t]; !ok {
			p = append(p, t+": sumiu do retrato")
		} else if d != a {
			p = append(p, fmt.Sprintf("%s: as linhas da empresa mudaram (antes %s, depois %s)", t, a, d))
		}
	}
	for t := range depois {
		if _, ok := antes[t]; !ok {
			p = append(p, t+": tabela nova no meio da varredura")
		}
	}
	sort.Strings(p)
	return p
}

// etagDoRegistro lê a versão do registro para o If-Match.
func etagDoRegistro(ctx context.Context, tabela, id string) (string, error) {
	var v int64
	if err := pool.QueryRow(ctx, `SELECT versao FROM `+tabela+` WHERE id = $1`, id).Scan(&v); err != nil {
		return "", fmt.Errorf("versão de %s %s: %w", tabela, id, err)
	}
	return fmt.Sprintf(`"%d"`, v), nil
}
