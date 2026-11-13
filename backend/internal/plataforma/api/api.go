// Package api monta o roteador HTTP e o contrato OpenAPI.
//
// O contrato é gerado a partir dos tipos Go anotados (princípio 4 da
// arquitetura: um contrato, vários clientes). Web, app e integradores usam os
// mesmos endpoints; não existe API interna paralela.
//
// Este pacote é a camada HTTP da plataforma: credencial (sessão web por
// cookie com CSRF, access token do app), escopo de acesso por rota,
// permissão, idempotência, ETag e If-Match, paginação e o formato de erro
// RFC 9457. As rotas de cada módulo ficam no pacote rotas do módulo e se
// registram com Registrar; a plataforma não importa módulo nenhum.
package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/danielgtaylor/huma/v2"
	"github.com/danielgtaylor/huma/v2/adapters/humago"

	"tenaz/internal/plataforma/autenticacao"
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

// Nomes do cookie de sessão web (seção 4.2 do documento 3).
const (
	// CookieSessao é o nome em homologação e produção: o prefixo __Host-
	// obriga Secure, Path=/ e ausência de Domain, e o navegador recusa o
	// cookie se qualquer um faltar.
	CookieSessao = "__Host-tenaz_sessao"
	// CookieSessaoDesenvolvimento é o nome em desenvolvimento sobre http,
	// em que Secure não funciona fora de localhost e o prefixo __Host- é
	// recusado. Só é aceito com TENAZ_AMBIENTE=desenvolvimento.
	CookieSessaoDesenvolvimento = "tenaz_sessao"
)

// Config é a configuração da camada HTTP.
type Config struct {
	// CookieSeguro: verdadeiro em homologação e produção (nome
	// __Host-tenaz_sessao, atributo Secure). Falso só em desenvolvimento.
	CookieSeguro bool
	// OrigensWeb são as origens aceitas em Origin nas requisições do
	// navegador que alteram dados (conferência de origem, seção 4.2).
	OrigensWeb []string
	// URLWeb é a base da web, para onde GET /oauth/authorize redireciona.
	URLWeb string
}

// NomeCookie devolve o nome do cookie de sessão conforme a configuração.
func (c Config) NomeCookie() string {
	if c.CookieSeguro {
		return CookieSessao
	}
	return CookieSessaoDesenvolvimento
}

// ConfigContrato é a configuração usada só para gerar o contrato.
func ConfigContrato() Config {
	return Config{CookieSeguro: true, URLWeb: "https://app.tenaz.com.br"}
}

// OrigemDe devolve o esquema e o host de uma URL ("https://app.tenaz.com.br").
func OrigemDe(u string) string {
	p, err := url.Parse(strings.TrimSpace(u))
	if err != nil || p.Scheme == "" || p.Host == "" {
		return ""
	}
	return p.Scheme + "://" + p.Host
}

// Servidor é a API montada: o mux, o contrato e as dependências da camada
// HTTP. Os módulos registram as suas rotas nele com Registrar.
type Servidor struct {
	cfg  Config
	mux  *http.ServeMux
	huma huma.API
	aut  *autenticacao.Servico
	idem *Idempotencia
}

// Novo monta o servidor e registra as rotas da plataforma (versão,
// autenticação e OAuth do app). aut pode ser nil só na geração do contrato:
// as rotas existem, e responderiam 503 servico_indisponivel.
func Novo(cfg Config, aut *autenticacao.Servico) *Servidor {
	mux := http.NewServeMux()

	hc := huma.DefaultConfig("Tenaz", VersaoContrato)
	hc.Info.Description = "API da Tenaz — Gestão de Manutenção e Frota. " +
		"Contrato único usado pela web, pelo aplicativo de campo e pelos integradores."
	hc.DocsPath = "/v1/documentacao"
	hc.OpenAPIPath = "/v1/openapi"

	// Sem o campo "$schema" nas respostas. GET /v1/versao devolve exatamente
	// versao, commit e data_build (seção 9.5 de docs/5_ARQUITETURA_SOFTWARE.md);
	// um campo a mais, ainda que inofensivo, já é contrato diferente do
	// especificado. Vale para todas as respostas, por coerência.
	//
	// O campo vem de um transformador que o Huma acrescenta dentro de New(),
	// por um CreateHook. Zerar SchemasPath só muda a URL do link, e zerar
	// Transformers antes de New() não resolve, porque o hook roda depois e
	// acrescenta o transformador de novo. É o hook que precisa sair.
	hc.CreateHooks = nil
	hc.SchemasPath = ""

	// PATCH aceita application/merge-patch+json com a mesma semântica de
	// application/json (seção 2.2 do documento 3); /oauth/token e
	// /oauth/revoke recebem formulário, como manda o OAuth.
	// Cópia: DefaultConfig devolve o mapa global da biblioteca.
	formatos := make(map[string]huma.Format, len(hc.Formats)+2)
	for k, v := range hc.Formats {
		formatos[k] = v
	}
	hc.Formats = formatos
	// Filtro não declarado responde 400 filtro_nao_suportado (seção 2.7).
	hc.RejectUnknownQueryParameters = true
	hc.Formats["application/merge-patch+json"] = huma.DefaultJSONFormat
	hc.Formats["application/x-www-form-urlencoded"] = formatoFormulario

	hc.Components.SecuritySchemes = map[string]*huma.SecurityScheme{
		esquemaSessao: {
			Type: "apiKey", In: "cookie", Name: CookieSessao,
			Description: "Sessão web (seção 4.2 do documento 3). Métodos que alteram dados exigem também " +
				"o cabeçalho X-CSRF-Token e Origin de uma origem da web.",
		},
		esquemaBearer: {
			Type: "http", Scheme: "bearer", BearerFormat: "JWT",
			Description: "Access token do app de campo, emitido por POST /oauth/token (seção 4.3 do documento 3).",
		},
	}

	s := &Servidor{cfg: cfg, mux: mux, huma: humago.New(mux, hc), aut: aut, idem: NovaIdempotencia()}
	s.huma.UseMiddleware(s.mwCredencial, s.mwIdempotencia)

	registrarVersao(s)
	registrarAutenticacao(s)
	registrarRotaDesconhecida(mux)
	return s
}

// Handler devolve o handler HTTP completo, com os envoltórios de toda
// requisição (identificador de rastreio, versão, contexto da auditoria).
func (s *Servidor) Handler() http.Handler {
	return envolver(s.mux)
}

// API devolve a API Huma, para gerar o contrato.
func (s *Servidor) API() huma.API { return s.huma }

// Config devolve a configuração da camada HTTP.
func (s *Servidor) Config() Config { return s.cfg }

// Autenticacao devolve o serviço de autenticação (nil na geração do
// contrato).
func (s *Servidor) Autenticacao() *autenticacao.Servico { return s.aut }

// registrarRotaDesconhecida faz caminho inexistente responder no formato
// RFC 9457, como qualquer outro erro da API.
//
// Sem isto, o roteador da biblioteca padrão responde texto puro, e um cliente
// que trata erro lendo o campo `codigo` quebraria justamente no caso mais
// comum: uma rota digitada errada por um integrador.
func registrarRotaDesconhecida(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		escreverProblema(w, r, erros.NaoEncontrado)
	})
}

// escreverProblema responde um erro fora do roteador do Huma.
func escreverProblema(w http.ResponseWriter, r *http.Request, err error) {
	p := erros.ParaProblema(err, r.URL.Path)
	p.IDRastreio = r.Header.Get(cabecalhoRequestID)
	corpo, errJSON := json.Marshal(p)
	if errJSON != nil {
		http.Error(w, "erro interno", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/problem+json")
	w.WriteHeader(int(p.Status))
	_, _ = w.Write(corpo)
}

// Versao é a resposta de GET /v1/versao.
type Versao struct {
	Versao    string `json:"versao" doc:"Versão do release em execução (SemVer)" example:"1.4.2"`
	Commit    string `json:"commit" doc:"Commit curto que gerou o build" example:"a1b2c3d"`
	DataBuild string `json:"data_build" doc:"Instante do build em UTC (RFC 3339)" example:"2026-11-12T14:05:00Z"`
}

type respostaVersao struct {
	Body Versao
}

// registrarVersao publica GET /v1/versao.
//
// Exceção aceita à regra de rota no plural (sumário executivo, seção 3.6).
// É público e sem autenticação, e devolve só versão, commit curto e data do
// build — nada de dependência, sistema operacional, região ou infraestrutura.
func registrarVersao(s *Servidor) {
	Registrar(s, Rota{
		ID: "obterVersao", Metodo: http.MethodGet, Caminho: "/v1/versao",
		Resumo:    "Versão publicada",
		Descricao: "Devolve a versão, o commit curto e a data do build do release em produção (RNF-039).",
		Tag:       TagPlataforma, Modulo: "ADM", Publica: true,
	}, func(_ context.Context, _ *struct{}) (*respostaVersao, error) {
		d := versao.Publico()
		return &respostaVersao{Body: Versao{Versao: d.Versao, Commit: d.Commit, DataBuild: d.DataBuild}}, nil
	})
}

// EscreverOpenAPI devolve o contrato em YAML, para o build versioná-lo em
// backend/api/openapi.yaml. O CI falha se o gerado divergir do versionado.
func EscreverOpenAPI(a huma.API) ([]byte, error) {
	return a.OpenAPI().YAML()
}

// formatoFormulario lê application/x-www-form-urlencoded como um objeto JSON
// de textos, para que a validação do contrato valha igual. Não serve para
// resposta: a API responde JSON.
var formatoFormulario = huma.Format{
	Marshal: func(_ io.Writer, _ any) error {
		return fmt.Errorf("a API não responde em formulário")
	},
	Unmarshal: func(dados []byte, v any) error {
		valores, err := url.ParseQuery(string(dados))
		if err != nil {
			return err
		}
		m := make(map[string]any, len(valores))
		for k, vs := range valores {
			if len(vs) > 0 {
				m[k] = vs[0]
			}
		}
		b, err := json.Marshal(m)
		if err != nil {
			return err
		}
		return json.Unmarshal(b, v)
	},
}
