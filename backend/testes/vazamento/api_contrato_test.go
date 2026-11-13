package vazamento

import (
	"bytes"
	"fmt"
	"os"
	"slices"
	"sort"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	"tenaz/internal/plataforma/api"
	"tenaz/internal/servidor"
)

// Parte 3 da suíte (seção 2.8 do documento 2): a API, rota a rota, guiada
// pelo contrato. A varredura não escolhe rota à mão: lê backend/api/openapi.yaml
// e percorre TODA operação autenticada. Rota nova entra sozinha; rota cujo
// teste precisa de um dado que a suíte não sabe montar (segmento de caminho
// sem recurso, corpo sem caso, campo uuid sem classificação) quebra o build
// com a instrução do que acrescentar em api_casos_test.go.

// arquivoContrato é o contrato versionado, relativo ao diretório do pacote.
const arquivoContrato = "../../api/openapi.yaml"

// parametroCaminho é um parâmetro do caminho e o segmento literal que vem
// antes dele, que diz de que recurso é o identificador.
type parametroCaminho struct {
	Nome     string
	Segmento string
}

// operacao é uma operação do contrato, com o que a varredura precisa saber.
type operacao struct {
	ID      string
	Metodo  string
	Caminho string
	Escopos []string
	Publica bool
	Sessao  bool // aceita a sessão web (cookie, CSRF e Origin)
	Bearer  bool // aceita o access token do app

	Caminhos []parametroCaminho
	Consulta map[string]*huma.Schema
	IfMatch  bool

	// Corpo é o esquema do corpo, já resolvido; nil quando a rota não tem.
	Corpo *huma.Schema
	// CamposUUID são os campos de formato uuid do corpo, em notação de ponto
	// ("nova_pessoa.unidade_id"), inclusive em objeto aninhado.
	CamposUUID []string
}

func (o operacao) nome() string { return o.Metodo + " " + o.Caminho }

func (o operacao) aceita(escopo string) bool { return slices.Contains(o.Escopos, escopo) }

// temID diz se algum parâmetro do caminho é identificador de registro.
func (o operacao) temID() bool {
	for _, p := range o.Caminhos {
		if r, ok := recursosDeCaminho[p.Segmento]; ok && r.Valor == "" {
			return true
		}
	}
	return false
}

// lerContratoVersionado lê o contrato versionado e confere que ele é
// exatamente o que o servidor gera. A varredura percorre a estrutura gerada
// (tipada), e a igualdade garante que ela é a do arquivo: um contrato
// desatualizado esconderia rota nova da suíte.
func lerContratoVersionado() (*huma.OpenAPI, error) { return lerContrato(arquivoContrato) }

func lerContrato(arquivo string) (*huma.OpenAPI, error) {
	versionado, err := os.ReadFile(arquivo)
	if err != nil {
		return nil, fmt.Errorf("ler o contrato versionado: %w", err)
	}
	gerador := servidor.Contrato()
	gerado, err := api.EscreverOpenAPI(gerador.API())
	if err != nil {
		return nil, fmt.Errorf("gerar o contrato: %w", err)
	}
	semCR := func(b []byte) []byte { return bytes.ReplaceAll(b, []byte("\r\n"), []byte("\n")) }
	if !bytes.Equal(semCR(versionado), semCR(gerado)) {
		return nil, fmt.Errorf("backend/api/openapi.yaml difere do contrato gerado pelo servidor: " +
			"a varredura por rota não pode confiar nele. Rode ferramentas/tarefas.sh contrato")
	}
	return gerador.API().OpenAPI(), nil
}

// operacoesDoContrato devolve todas as operações, em ordem de caminho e método.
func operacoesDoContrato(oapi *huma.OpenAPI) ([]operacao, error) {
	var ops []operacao
	for caminho, item := range oapi.Paths {
		for metodo, op := range map[string]*huma.Operation{
			"GET": item.Get, "POST": item.Post, "PATCH": item.Patch, "PUT": item.Put, "DELETE": item.Delete,
			"HEAD": item.Head, "OPTIONS": item.Options, "TRACE": item.Trace,
		} {
			if op == nil {
				continue
			}
			o, err := lerOperacao(oapi, metodo, caminho, op)
			if err != nil {
				return nil, err
			}
			ops = append(ops, o)
		}
	}
	sort.Slice(ops, func(i, j int) bool {
		if ops[i].Caminho != ops[j].Caminho {
			return ops[i].Caminho < ops[j].Caminho
		}
		return ops[i].Metodo < ops[j].Metodo
	})
	return ops, nil
}

func lerOperacao(oapi *huma.OpenAPI, metodo, caminho string, op *huma.Operation) (operacao, error) {
	o := operacao{ID: op.OperationID, Metodo: metodo, Caminho: caminho, Consulta: map[string]*huma.Schema{}}
	escopos, ok := op.Extensions["x-escopos-acesso"].([]string)
	if !ok || len(escopos) == 0 {
		return o, fmt.Errorf("%s: sem x-escopos-acesso; a suíte não sabe quem pode chamar a rota", o.nome())
	}
	o.Escopos = escopos
	o.Publica = len(escopos) == 1 && escopos[0] == "publico"
	for _, req := range op.Security {
		if _, tem := req["sessao"]; tem {
			o.Sessao = true
		}
		if _, tem := req["bearer"]; tem {
			o.Bearer = true
		}
	}
	if !o.Publica && !o.Sessao && !o.Bearer {
		return o, fmt.Errorf("%s: rota autenticada sem esquema de segurança conhecido (%v)", o.nome(), op.Security)
	}

	segmentos := strings.Split(strings.TrimPrefix(caminho, "/"), "/")
	for i, s := range segmentos {
		if strings.HasPrefix(s, "{") {
			anterior := ""
			if i > 0 {
				anterior = segmentos[i-1]
			}
			o.Caminhos = append(o.Caminhos, parametroCaminho{Nome: strings.Trim(s, "{}"), Segmento: anterior})
		}
	}
	for _, p := range op.Parameters {
		switch {
		case p.In == "query":
			o.Consulta[p.Name] = resolverEsquema(oapi, p.Schema)
		case p.In == "header" && strings.EqualFold(p.Name, "If-Match"):
			o.IfMatch = true
		}
	}
	if op.RequestBody != nil {
		for _, tipo := range []string{"application/json", "application/merge-patch+json", "application/x-www-form-urlencoded"} {
			if mt, ok := op.RequestBody.Content[tipo]; ok && mt.Schema != nil {
				o.Corpo = resolverEsquema(oapi, mt.Schema)
				break
			}
		}
		if o.Corpo == nil {
			return o, fmt.Errorf("%s: corpo sem esquema JSON nem formulário", o.nome())
		}
		o.CamposUUID = camposUUID(oapi, o.Corpo, "", 0)
	}
	return o, nil
}

// resolverEsquema segue a referência para o componente.
func resolverEsquema(oapi *huma.OpenAPI, s *huma.Schema) *huma.Schema {
	for s != nil && s.Ref != "" {
		nome := strings.TrimPrefix(s.Ref, "#/components/schemas/")
		s = oapi.Components.Schemas.Map()[nome]
	}
	return s
}

// camposUUID lista os campos de formato uuid, descendo em objetos e listas.
func camposUUID(oapi *huma.OpenAPI, s *huma.Schema, prefixo string, profundidade int) []string {
	s = resolverEsquema(oapi, s)
	if s == nil || profundidade > 4 {
		return nil
	}
	var campos []string
	nomes := make([]string, 0, len(s.Properties))
	for n := range s.Properties {
		nomes = append(nomes, n)
	}
	sort.Strings(nomes)
	for _, n := range nomes {
		p := resolverEsquema(oapi, s.Properties[n])
		if p == nil {
			continue
		}
		caminho := prefixo + n
		if p.Format == "uuid" {
			campos = append(campos, caminho)
		}
		if p.Items != nil {
			it := resolverEsquema(oapi, p.Items)
			if it != nil && it.Format == "uuid" {
				campos = append(campos, caminho+"[]")
			}
			campos = append(campos, camposUUID(oapi, it, caminho+"[].", profundidade+1)...)
		}
		campos = append(campos, camposUUID(oapi, p, caminho+".", profundidade+1)...)
	}
	return campos
}

// mesmasRotas confere que o servidor montado atende exatamente as rotas do
// contrato: rota registrada fora do contrato escaparia da varredura.
func mesmasRotas(contrato, montado *huma.OpenAPI) []string {
	chaves := func(o *huma.OpenAPI) map[string]bool {
		m := map[string]bool{}
		for caminho, item := range o.Paths {
			for metodo, op := range map[string]*huma.Operation{
				"GET": item.Get, "POST": item.Post, "PATCH": item.Patch, "PUT": item.Put, "DELETE": item.Delete,
				"HEAD": item.Head, "OPTIONS": item.Options, "TRACE": item.Trace,
			} {
				if op != nil {
					m[metodo+" "+caminho] = true
				}
			}
		}
		return m
	}
	c, s := chaves(contrato), chaves(montado)
	var p []string
	for k := range s {
		if !c[k] {
			p = append(p, "rota atendida pelo servidor e ausente do contrato: "+k)
		}
	}
	for k := range c {
		if !s[k] {
			p = append(p, "rota do contrato que o servidor montado não atende: "+k)
		}
	}
	sort.Strings(p)
	return p
}
