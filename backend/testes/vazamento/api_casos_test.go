package vazamento

import (
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Casos da parte 3. É AQUI que uma rota nova pede o que a suíte não deduz do
// contrato, e só isso; o resto (quem pode chamar, se aceita o app, quais
// parâmetros são identificadores, quais campos do corpo são uuid) vem do
// openapi.yaml. A varredura quebra o build, com a instrução, quando:
//
//   - um parâmetro de caminho vem depois de um segmento que não está em
//     recursosDeCaminho (a suíte não sabe de que tabela tirar o id de B);
//   - uma rota autenticada tem corpo e não está em casosDeCorpo (a suíte não
//     sabe montar um corpo válido, e um corpo inválido daria 400/422 antes de
//     o caso de uso procurar o registro, o que não prova nada);
//   - um campo uuid de corpo ou de consulta não está em referenciasPorCampo;
//   - um caso aponta para rota que não existe mais no contrato.
//
// Revisor obrigatório de toda mudança neste arquivo: tenaz-isolamento.

// recursoAPI diz de onde tirar o identificador de um parâmetro de caminho.
type recursoAPI struct {
	// Tabela do registro (com esquema). Vazia quando o parâmetro não é
	// identificador de registro.
	Tabela string
	// ColunaPai liga o sub-recurso ao recurso do parâmetro anterior
	// (/v1/unidades/{id}/turnos/{turno_id}: turno_unidade.unidade_id).
	ColunaPai string
	// Valor, para parâmetro que não é identificador (código de catálogo): vai
	// igual nas duas variantes, e a variante muda só os identificadores.
	Valor string
	// Justificativa, obrigatória quando Valor não é vazio.
	Justificativa string
}

// recursosDeCaminho: segmento literal antes do parâmetro -> recurso.
var recursosDeCaminho = map[string]recursoAPI{
	"unidades":          {Tabela: "app.unidade"},
	"turnos":            {Tabela: "app.turno_unidade", ColunaPai: "unidade_id"},
	"feriados":          {Tabela: "app.feriado_unidade", ColunaPai: "unidade_id"},
	"pessoas":           {Tabela: "app.pessoa"},
	"usuarios":          {Tabela: "app.usuario"},
	"eventos-auditoria": {Tabela: "app.evento_auditoria"},
	"capturas-auditoria": {Valor: telaDeCaptura, Justificativa: "código do catálogo fechado de telas, não identificador; " +
		"B configurou a tela no cenário, e uma escrita de A que alcançasse B mudaria o retrato de B"},
	"permissoes": {Valor: "os.encerrar", Justificativa: "código do catálogo de permissões, não identificador; " +
		"o técnico de A e o de B têm a permissão avulsa os.encerrar, então a retirada em B seria visível no retrato do banco"},
}

// referenciaAPI classifica um campo uuid de corpo ou de consulta.
type referenciaAPI struct {
	// Recurso: chave de recursosDeCaminho, ou "empresa" (app.empresa).
	Recurso string
	// Solto: justificativa de por que o campo não é testado como referência
	// a registro de outra empresa.
	Solto string
}

// referenciasPorCampo: nome do campo (o último trecho da notação de ponto).
var referenciasPorCampo = map[string]referenciaAPI{
	"empresa_escolhida_id": {Recurso: "empresa"},
	"pessoa_id":            {Recurso: "pessoas"},
	"unidade_id":           {Recurso: "unidades"},
	"unidade_padrao_id":    {Recurso: "unidades"},
	"usuario_id":           {Recurso: "usuarios"},
	// Na trilha, entidade_id aponta para qualquer entidade: a pessoa de B tem
	// o evento da própria criação na trilha de B.
	"entidade_id": {Recurso: "pessoas"},
	"fornecedor_id": {Solto: "as tabelas fornecedor (E06) e locatario (fase 2) não existem: o id é solto, e o " +
		"mesmo fornecedor_id em duas empresas é o caso que a suíte simula de propósito (seção 9.3 do documento 5)"},
	"locatario_id": {Solto: "as tabelas fornecedor (E06) e locatario (fase 2) não existem: o id é solto"},
}

// casoCorpo monta um corpo válido para a rota.
type casoCorpo struct {
	// corpo devolve um corpo válido, sem campo de referência (a variante
	// acrescenta o campo). É chamado a cada requisição: e-mail novo, etc.
	corpo func(e *empresaAPI) map[string]any
	// porCampo, para campo de referência aninhado: o corpo inteiro com o
	// valor no lugar certo.
	porCampo map[string]func(e *empresaAPI, ref string) map[string]any
	// escritaPropria: a varredura executa a rota com a credencial de A e o
	// corpo válido (exige sucesso), para que o retrato de B e de C prove que a
	// escrita da própria empresa não as alcança. Na rota declarada para o
	// perfil externo, também com a credencial dele (web e app), sem
	// identificador no caminho. Só para rota sem efeito sobre a sessão e
	// repetível.
	escritaPropria bool
}

// telaDeCaptura e capturaDeB: a configuração de captura que só B tem.
const telaDeCaptura = "adm.pessoas"

var capturaDeB = map[string]any{"visualizacao": true, "impressao": false, "extracao": false}

// consultasPadrao: parâmetros de consulta sem os quais a rota não responde
// 200 (período obrigatório). Entram em toda requisição da varredura à rota.
var consultasPadrao = map[string]func() url.Values{
	"exportarEventosAuditoria": func() url.Values {
		agora := time.Now().UTC()
		return url.Values{"ocorrido_em[gte]": {agora.Add(-48 * time.Hour).Format(time.RFC3339)},
			"ocorrido_em[lt]": {agora.Add(time.Hour).Format(time.RFC3339)}}
	},
	"verificarIntegridadeAuditoria": func() url.Values {
		hoje := time.Now().UTC()
		return url.Values{"de": {hoje.AddDate(0, 0, -7).Format("2006-01-02")}, "ate": {hoje.Format("2006-01-02")}}
	},
}

// verificacoesDeLeitura: o que a varredura genérica (UUID e texto alheios)
// não vê numa leitura. e é a empresa de quem lê; o retorno são problemas.
var verificacoesDeLeitura = map[string]func(s *sondador, e *empresaAPI, r *respostaAPI) []string{
	// A captura é booleana: um vazamento da configuração de B não traz UUID
	// nem texto. B configurou a tela com uma combinação fora do padrão; A, que
	// não configurou, não pode vê-la, e B, que configurou, precisa vê-la.
	"listarCapturasAuditoria": func(s *sondador, e *empresaAPI, r *respostaAPI) []string {
		d, _ := r.dados()
		var tela map[string]any
		for _, x := range d {
			if m, ok := x.(map[string]any); ok && m["tela"] == telaDeCaptura {
				tela = m
			}
		}
		if tela == nil {
			return []string{"GET /v1/capturas-auditoria: a tela " + telaDeCaptura + " sumiu do catálogo; troque telaDeCaptura"}
		}
		igualB := true
		for k, v := range capturaDeB {
			igualB = igualB && tela[k] == v
		}
		switch {
		case e == s.cen.B && (!igualB || tela["configurada"] != true):
			return []string{fmt.Sprintf("controle de GET /v1/capturas-auditoria: B não vê a própria configuração (%v)", tela)}
		case e != s.cen.B && igualB:
			return []string{fmt.Sprintf("GET /v1/capturas-auditoria por %s mostra a configuração de captura de B: %v", e.rotulo, tela)}
		}
		return nil
	},
}

var turnoDeTeste = []any{map[string]any{"dia_semana": 1, "hora_inicio": "08:00", "hora_fim": "17:00"}}

func corpoFixo(m map[string]any) func(*empresaAPI) map[string]any {
	return func(*empresaAPI) map[string]any {
		c := make(map[string]any, len(m))
		for k, v := range m {
			c[k] = v
		}
		return c
	}
}

// casosDeCorpo: operationId -> caso. Toda rota autenticada com corpo precisa
// de um.
var casosDeCorpo = map[string]casoCorpo{
	"autorizarApp": {corpo: func(*empresaAPI) map[string]any {
		_, desafio := desafioPKCE()
		return map[string]any{"response_type": "code", "client_id": clienteOAuthApp, "redirect_uri": retornoOAuthApp,
			"code_challenge": desafio, "code_challenge_method": "S256", "state": "vazamento"}
	}},
	"confirmarMFA":  {corpo: corpoFixo(map[string]any{"codigo": "123456"})},
	"desativarMFA":  {corpo: corpoFixo(map[string]any{"codigo": "123456"})},
	"trocarEmpresa": {corpo: corpoFixo(map[string]any{})},
	"trocarSenha":   {corpo: corpoFixo(map[string]any{"senha_atual": "senha errada da suíte", "nova_senha": "Outra senha longa da suíte 42"})},
	// As três escritas de rota única (PATCH /v1/empresa e /v1/usuarios/eu, POST
	// /v1/unidades) rodam com sucesso: sem isso, o retrato de B não provaria
	// nada sobre elas. O If-Match da rota única vem do GET da mesma rota.
	"alterarEmpresa": {corpo: corpoFixo(map[string]any{"nome_fantasia": "Fantasia da varredura"}), escritaPropria: true},
	"criarPessoa":    {corpo: corpoFixo(map[string]any{"nome": "Pessoa da varredura"}), escritaPropria: true},
	"alterarCapturaAuditoria": {corpo: corpoFixo(map[string]any{"visualizacao": true, "impressao": true, "extracao": true}),
		escritaPropria: true},
	"alterarPessoa":  {corpo: corpoFixo(map[string]any{"cargo": "Mecânico"})},
	"incluirFeriado": {corpo: corpoFixo(map[string]any{"data": "2027-09-07", "descricao": "Independência"})},
	"incluirTurno":   {corpo: corpoFixo(map[string]any{"dia_semana": 6, "hora_inicio": "08:00", "hora_fim": "12:00"})},
	"alterarTurno":   {corpo: corpoFixo(map[string]any{"hora_fim": "18:00"})},
	"criarUnidade": {corpo: func(e *empresaAPI) map[string]any {
		return map[string]any{"codigo": "V" + aleatorioHex(3), "nome": "Unidade da varredura", "turnos": turnoDeTeste,
			"endereco": e.endereco("Cidade da varredura")}
	}, escritaPropria: true},
	"alterarUnidade": {corpo: corpoFixo(map[string]any{"nome": "Invadida"})},
	// O idioma é o OUTRO, não o da empresa: a trava do perfil externo
	// (app.tg_usuario_colunas_externo, 00009) só age quando o valor muda, e a
	// primeira escrita de cada usuário muda o idioma de fato. As comparações
	// da varredura são entre respostas ao mesmo usuário, no mesmo idioma.
	"alterarUsuarioEu": {corpo: func(e *empresaAPI) map[string]any {
		idioma := "es-419"
		if e.pa.idioma == idioma {
			idioma = "pt-BR"
		}
		return map[string]any{"nome": "Eu mesmo", "idioma": idioma}
	}, escritaPropria: true},
	"concederPermissaoUsuario": {corpo: corpoFixo(map[string]any{"permissao": "est.decidir_destino"})},
	"alterarUsuario":           {corpo: corpoFixo(map[string]any{"telefone": "11999990000"})},
	"criarUsuario": {
		corpo: func(*empresaAPI) map[string]any {
			return map[string]any{"nome": "Usuário da varredura", "email": emailAPI("varredura"), "perfis": []any{"tecnico"}}
		},
		porCampo: map[string]func(*empresaAPI, string) map[string]any{
			"nova_pessoa.unidade_id": func(_ *empresaAPI, ref string) map[string]any {
				return map[string]any{"nome": "Usuário da varredura", "email": emailAPI("varredura"), "perfis": []any{"tecnico"},
					"nova_pessoa": map[string]any{"nome": "Pessoa nova da varredura", "unidade_id": ref}}
			},
		},
	},
}

// refCampo é um campo de referência testado.
type refCampo struct {
	Campo   string // notação de ponto
	Recurso string
}

// planoDaRota é o que a varredura faz com a rota.
type planoDaRota struct {
	op           operacao
	recursos     []recursoAPI // um por parâmetro de caminho
	caso         *casoCorpo
	refsCorpo    []refCampo
	refsConsulta []refCampo
	busca        bool // aceita busca textual
	cursor       bool // aceita cursor e limite
}

// get: a rota é de leitura. Se é lista ({dados, paginacao}), a resposta diz.
func (p planoDaRota) get() bool { return p.op.Metodo == "GET" }

func ultimoTrecho(campo string) string {
	campo = strings.TrimSuffix(campo, "[]")
	if i := strings.LastIndex(campo, "."); i >= 0 {
		campo = campo[i+1:]
	}
	return strings.TrimSuffix(campo, "[]")
}

// planejar monta o plano de cada rota autenticada e devolve, como problema,
// tudo o que a suíte não sabe testar.
func planejar(ops []operacao) ([]planoDaRota, []string) {
	var (
		planos []planoDaRota
		p      []string
		usados = map[string]bool{}
	)
	for _, op := range ops {
		if op.Publica {
			continue
		}
		pl := planoDaRota{op: op}
		for _, pc := range op.Caminhos {
			r, ok := recursosDeCaminho[pc.Segmento]
			switch {
			case !ok:
				p = append(p, fmt.Sprintf("%s: o parâmetro {%s} vem depois de /%s/, que não está em recursosDeCaminho "+
					"(api_casos_test.go); acrescente a tabela do recurso e povoe A e B no cenário", op.nome(), pc.Nome, pc.Segmento))
			case r.Valor != "" && r.Justificativa == "":
				p = append(p, fmt.Sprintf("%s: recurso %q com valor fixo e sem justificativa", op.nome(), pc.Segmento))
			}
			pl.recursos = append(pl.recursos, r)
		}
		if op.Corpo != nil {
			c, ok := casosDeCorpo[op.ID]
			if !ok {
				p = append(p, fmt.Sprintf("%s (%s): rota autenticada com corpo e sem caso em casosDeCorpo (api_casos_test.go)",
					op.nome(), op.ID))
			} else {
				usados[op.ID] = true
				pl.caso = &c
			}
			for _, campo := range op.CamposUUID {
				ref, ok := referenciasPorCampo[ultimoTrecho(campo)]
				switch {
				case !ok:
					p = append(p, fmt.Sprintf("%s: campo uuid %q do corpo sem classificação em referenciasPorCampo", op.nome(), campo))
				case ref.Solto != "":
				case strings.Contains(campo, "[]"):
					p = append(p, fmt.Sprintf("%s: campo uuid em lista %q: acrescente o caso na varredura", op.nome(), campo))
				case strings.Contains(campo, ".") && (pl.caso == nil || pl.caso.porCampo[campo] == nil):
					p = append(p, fmt.Sprintf("%s: campo de referência aninhado %q sem porCampo no caso", op.nome(), campo))
				default:
					pl.refsCorpo = append(pl.refsCorpo, refCampo{Campo: campo, Recurso: ref.Recurso})
				}
			}
		}
		nomes := make([]string, 0, len(op.Consulta))
		for n := range op.Consulta {
			nomes = append(nomes, n)
		}
		sort.Strings(nomes)
		for _, n := range nomes {
			s := op.Consulta[n]
			switch {
			case n == "busca":
				pl.busca = true
			case n == "cursor":
				_, temLimite := op.Consulta["limite"]
				pl.cursor = temLimite
			case s != nil && s.Format == "uuid":
				ref, ok := referenciasPorCampo[n]
				switch {
				case !ok:
					p = append(p, fmt.Sprintf("%s: parâmetro de consulta uuid %q sem classificação em referenciasPorCampo", op.nome(), n))
				case ref.Solto == "":
					pl.refsConsulta = append(pl.refsConsulta, refCampo{Campo: n, Recurso: ref.Recurso})
				}
			}
		}
		planos = append(planos, pl)
	}
	ids := map[string]bool{}
	for _, op := range ops {
		ids[op.ID] = true
	}
	for _, m := range []map[string]bool{chaves(consultasPadrao), chaves(verificacoesDeLeitura)} {
		for id := range m {
			if !ids[id] {
				p = append(p, fmt.Sprintf("caso para %q: nenhuma rota do contrato tem esse operationId; retire o caso", id))
			}
		}
	}
	for id := range casosDeCorpo {
		if !usados[id] {
			p = append(p, fmt.Sprintf("casosDeCorpo[%q]: nenhuma rota autenticada com corpo tem esse operationId; retire o caso", id))
		}
	}
	for campo, r := range referenciasPorCampo {
		if r.Solto == "" && r.Recurso != "empresa" {
			if _, ok := recursosDeCaminho[r.Recurso]; !ok {
				p = append(p, fmt.Sprintf("referenciasPorCampo[%q]: recurso %q fora de recursosDeCaminho", campo, r.Recurso))
			}
		}
	}
	sort.Strings(p)
	return planos, p
}

func chaves[V any](m map[string]V) map[string]bool {
	r := make(map[string]bool, len(m))
	for k := range m {
		r[k] = true
	}
	return r
}

// corpoComReferencia devolve o corpo da rota com ref no campo.
func (p planoDaRota) corpoComReferencia(e *empresaAPI, campo, ref string) map[string]any {
	if f := p.caso.porCampo[campo]; f != nil {
		return f(e, ref)
	}
	c := p.caso.corpo(e)
	c[campo] = ref
	return c
}

// corpo devolve o corpo base da rota, ou nil.
func (p planoDaRota) corpo(e *empresaAPI) any {
	if p.caso == nil {
		return nil
	}
	return p.caso.corpo(e)
}

func idInexistente() string { return uuid.NewString() }
