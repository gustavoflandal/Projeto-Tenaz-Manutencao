package servidor_test

import (
	"regexp"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"

	"tenaz/internal/servidor"
)

var reSegmento = regexp.MustCompile(`^([a-z0-9]+(-[a-z0-9]+)*|\{[a-z_]+\})$`)

// O contrato gerado segue as convenções do documento 3: rota em kebab-case,
// escopo de acesso declarado em toda operação (a suíte de isolamento lê
// x-escopos-acesso), erro sempre em problem+json e nenhum empresa_id no
// corpo de requisição ou resposta.
func TestConvencoesDoContrato(t *testing.T) {
	oapi := servidor.Contrato().API().OpenAPI()
	if len(oapi.Paths) == 0 {
		t.Fatal("contrato sem rotas")
	}
	for caminho, item := range oapi.Paths {
		if !strings.HasPrefix(caminho, "/v1/") && !strings.HasPrefix(caminho, "/oauth/") {
			t.Errorf("%s: fora de /v1 e /oauth", caminho)
		}
		for _, seg := range strings.Split(strings.TrimPrefix(caminho, "/"), "/") {
			if !reSegmento.MatchString(seg) {
				t.Errorf("%s: segmento %q fora do kebab-case", caminho, seg)
			}
		}
		for metodo, op := range map[string]*huma.Operation{"GET": item.Get, "POST": item.Post, "PATCH": item.Patch, "PUT": item.Put, "DELETE": item.Delete} {
			if op == nil {
				continue
			}
			nome := metodo + " " + caminho
			if op.OperationID == "" || op.Extensions["x-modulo"] == "" || op.Extensions["x-fase"] == nil {
				t.Errorf("%s: sem operationId, x-modulo ou x-fase", nome)
			}
			escopos, _ := op.Extensions["x-escopos-acesso"].([]string)
			if len(escopos) == 0 {
				t.Errorf("%s: sem x-escopos-acesso", nome)
			}
			publica := len(escopos) == 1 && escopos[0] == "publico"
			if publica != (op.Security != nil && len(op.Security) == 0) {
				t.Errorf("%s: x-escopos-acesso %v não confere com security %v", nome, escopos, op.Security)
			}
			for status, r := range op.Responses {
				if strings.HasPrefix(status, "2") || strings.HasPrefix(status, "3") {
					continue
				}
				if _, ok := r.Content["application/problem+json"]; !ok {
					t.Errorf("%s: resposta %s sem application/problem+json", nome, status)
				}
			}
		}
	}
	for nome, s := range oapi.Components.Schemas.Map() {
		if _, tem := s.Properties["empresa_id"]; tem {
			t.Errorf("esquema %s tem empresa_id: a empresa vem só da credencial (seção 4.7)", nome)
		}
	}
}
