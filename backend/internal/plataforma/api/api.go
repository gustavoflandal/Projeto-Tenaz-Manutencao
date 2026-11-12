// Package api monta o roteador HTTP e o contrato OpenAPI.
//
// O contrato é gerado a partir dos tipos Go anotados (princípio 4 da
// arquitetura: um contrato, vários clientes). Web, app e integradores usam os
// mesmos endpoints; não existe API interna paralela.
package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"

	"tenaz/internal/plataforma/erros"
	"tenaz/internal/plataforma/versao"
)

// VersaoContrato é a versão do CONTRATO, não a do build.
//
// Fica fixa em "1" porque é a versão do caminho `/v1`: muda só quando houver
// `/v2`, com convivência mínima de 12 meses (seção 8.3 de
// docs/5_ARQUITETURA_SOFTWARE.md). Usar aqui a versão do build faria o
// contrato gerado mudar a cada compilação, e a verificação de divergência do
// CI — que compara o gerado com o versionado — acusaria falso positivo em
// todo pull request. A versão do build é exposta por GET /v1/versao.
const VersaoContrato = "1"

func init() {
	// Todo erro da API sai no formato RFC 9457 com o campo `codigo` estável
	// (seção 3 de docs/3_ESPECIFICACAO_APIS.md). Sem esta troca, o roteador
	// usaria o modelo de erro padrão da biblioteca, que não tem `codigo`.
	huma.NewError = func(status int, msg string, errs ...error) huma.StatusError {
		p := erros.ParaProblema(mapearStatus(status), "")
		if msg != "" {
			p.Detail = msg
		}
		for _, e := range errs {
			if e == nil {
				continue
			}
			var detalhe huma.ErrorDetail
			if ok := comoDetalhe(e, &detalhe); ok {
				p.Campos = append(p.Campos, erros.Campo{Local: detalhe.Location, Mensagem: detalhe.Message})
				continue
			}
			p.Campos = append(p.Campos, erros.Campo{Mensagem: e.Error()})
		}
		return &p
	}
}

func comoDetalhe(err error, destino *huma.ErrorDetail) bool {
	if d, ok := err.(*huma.ErrorDetail); ok {
		*destino = *d
		return true
	}
	return false
}

// mapearStatus traduz o status que a biblioteca escolheu para o erro de
// domínio correspondente, para que o `codigo` seja o da Tenaz.
func mapearStatus(status int) error {
	switch status {
	case http.StatusNotFound:
		return erros.NaoEncontrado
	case http.StatusUnauthorized:
		return erros.NaoAutenticado
	case http.StatusForbidden:
		return erros.SemPermissao
	case http.StatusUnprocessableEntity, http.StatusBadRequest:
		return erros.EntradaInvalida
	case http.StatusPreconditionFailed, http.StatusConflict:
		return erros.Conflito
	case http.StatusTooManyRequests:
		return erros.LimiteExcedido
	default:
		return erros.Interno
	}
}

// Novo monta o mux e registra as rotas da plataforma.
func Novo() (*http.ServeMux, huma.API) {
	mux := http.NewServeMux()

	cfg := huma.DefaultConfig("Tenaz", VersaoContrato)
	cfg.Info.Description = "API da Tenaz — Gestão de Manutenção e Frota. " +
		"Contrato único usado pela web, pelo aplicativo de campo e pelos integradores."
	cfg.DocsPath = "/v1/documentacao"
	cfg.OpenAPIPath = "/v1/openapi"

	// Sem o campo "$schema" nas respostas. GET /v1/versao devolve exatamente
	// versao, commit e data_build (seção 9.5 de docs/5_ARQUITETURA_SOFTWARE.md);
	// um campo a mais, ainda que inofensivo, já é contrato diferente do
	// especificado. Vale para todas as respostas, por coerência.
	//
	// O campo vem de um transformador que o Huma acrescenta dentro de New(),
	// por um CreateHook. Zerar SchemasPath só muda a URL do link, e zerar
	// Transformers antes de New() não resolve, porque o hook roda depois e
	// acrescenta o transformador de novo. É o hook que precisa sair.
	cfg.CreateHooks = nil
	cfg.SchemasPath = ""

	api := humago.New(mux, cfg)

	registrarVersao(api)
	registrarRotaDesconhecida(mux)

	return mux, api
}

// registrarRotaDesconhecida faz caminho inexistente responder no formato
// RFC 9457, como qualquer outro erro da API.
//
// Sem isto, o roteador da biblioteca padrão responde texto puro, e um cliente
// que trata erro lendo o campo `codigo` quebraria justamente no caso mais
// comum: uma rota digitada errada por um integrador.
func registrarRotaDesconhecida(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		p := erros.ParaProblema(erros.NaoEncontrado, r.URL.Path)
		corpo, err := json.Marshal(p)
		if err != nil {
			http.Error(w, "erro interno", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/problem+json")
		w.WriteHeader(p.Status)
		_, _ = w.Write(corpo)
	})
}

type respostaVersao struct {
	Body versao.Dados
}

// registrarVersao publica GET /v1/versao.
//
// Exceção aceita à regra de rota no plural (sumário executivo, seção 3.6).
// É público e sem autenticação, e devolve só versão, commit curto e data do
// build — nada de dependência, sistema operacional, região ou infraestrutura.
func registrarVersao(api huma.API) {
	huma.Register(api, huma.Operation{
		OperationID: "obterVersao",
		Method:      http.MethodGet,
		Path:        "/v1/versao",
		Summary:     "Versão publicada",
		Description: "Devolve a versão, o commit curto e a data do build do release em produção (RNF-039).",
		Tags:        []string{"Plataforma"},
		Security:    []map[string][]string{},
	}, func(_ context.Context, _ *struct{}) (*respostaVersao, error) {
		return &respostaVersao{Body: versao.Publico()}, nil
	})
}

// EscreverOpenAPI devolve o contrato em YAML, para o build versioná-lo em
// backend/api/openapi.yaml. O CI falha se o gerado divergir do versionado.
func EscreverOpenAPI(api huma.API) ([]byte, error) {
	return api.OpenAPI().YAML()
}
