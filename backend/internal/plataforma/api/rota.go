package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/erros"
	"tenaz/internal/plataforma/obs"
)

// Etiquetas do contrato.
const (
	TagPlataforma   = "Plataforma"
	TagAutenticacao = "Autenticação"
	TagOAuth        = "OAuth do app"
	TagADM          = "ADM — Administração"
)

const (
	esquemaSessao = "sessao"
	esquemaBearer = "bearer"
)

// Escopos de acesso (sumário 3.6; seção 4.9 do documento 3).
var (
	// EscoposInternos: a rota existe só para usuário interno. É o padrão.
	EscoposInternos = []string{db.EscopoInterno}
	// EscoposDoProprioUsuario: rotas de sessão e de infraestrutura do próprio
	// usuário, que os perfis externos também usam (seção 4.8).
	EscoposDoProprioUsuario = []string{db.EscopoInterno, db.EscopoPrestador, db.EscopoLocatario}
)

// Rota descreve uma operação do contrato e as regras que a camada HTTP aplica
// antes do handler.
type Rota struct {
	ID        string // operationId
	Metodo    string
	Caminho   string
	Resumo    string
	Descricao string
	Tag       string
	Modulo    string // código do módulo (x-modulo)
	Fase      string // MVP, F2, F3 (x-fase); vazio: MVP
	Plano     string // plano mínimo (x-plano-minimo); vazio: E

	// Publica: sem credencial.
	Publica bool
	// ApenasSessaoWeb: só a sessão web por cookie (POST /oauth/authorize).
	ApenasSessaoWeb bool
	// Escopos aceitos; vazio: só interno. Escopo fora da lista: 404.
	Escopos []string
	// Permissoes exigidas pela rota (x-permissoes), conferidas antes do
	// handler, salvo com AlcanceNoServico.
	Permissoes []string
	// AlcanceNoServico: a permissão depende do registro (o titular lê a
	// própria pessoa) e é conferida pelo caso de uso.
	AlcanceNoServico bool
	// SemIdempotencia: POST que não guarda resposta para repetição, porque
	// a resposta traz credencial nova (login, troca de empresa, tokens).
	// A chave é aceita e ignorada.
	SemIdempotencia bool

	Status int // status de sucesso; vazio: 200 (201 em criação, 204 sem corpo)

	// Erros: status de erro documentados além dos padrão da rota.
	Erros []int

	// AceitaMFAPendente: a rota atende a sessão restrita à ativação do MFA
	// (empresa com mfa_obrigatorio e identidade sem MFA, RN-ADM-027): ler a
	// sessão, sair, trocar de empresa e ativar o MFA. As demais respondem 401
	// nao_autenticado com mfa_pendente.
	AceitaMFAPendente bool
}

type chaveRota struct{}

// Registrar registra a operação no contrato e no roteador.
func Registrar[I, O any](s *Servidor, r Rota, h func(context.Context, *I) (*O, error)) {
	escopos := r.Escopos
	if len(escopos) == 0 && !r.Publica {
		escopos = EscoposInternos
	}
	r.Escopos = escopos
	fase, plano := r.Fase, r.Plano
	if fase == "" {
		fase = "MVP"
	}
	if plano == "" {
		plano = "E"
	}
	ext := map[string]any{"x-modulo": r.Modulo, "x-fase": fase, "x-plano-minimo": plano}
	if r.Publica {
		ext["x-escopos-acesso"] = []string{"publico"}
	} else {
		ext["x-escopos-acesso"] = escopos
	}
	if len(r.Permissoes) > 0 {
		ext["x-permissoes"] = r.Permissoes
	}

	op := huma.Operation{
		OperationID:   r.ID,
		Method:        r.Metodo,
		Path:          r.Caminho,
		Summary:       r.Resumo,
		Description:   r.Descricao,
		Tags:          []string{r.Tag},
		DefaultStatus: r.Status,
		Extensions:    ext,
		Metadata:      map[string]any{"rota": &r},
	}
	switch {
	case r.Publica:
		op.Security = []map[string][]string{}
	case r.ApenasSessaoWeb:
		op.Security = []map[string][]string{{esquemaSessao: {}}}
	default:
		op.Security = []map[string][]string{{esquemaSessao: {}}, {esquemaBearer: {}}}
	}
	op.Errors = errosDaRota(r)
	huma.Register(s.huma, op, h)
}

// errosDaRota lista os status de erro documentados da operação.
func errosDaRota(r Rota) []int {
	l := []int{http.StatusBadRequest, http.StatusUnprocessableEntity}
	for _, e := range r.Erros {
		if !slices.Contains(l, e) {
			l = append(l, e)
		}
	}
	if !r.Publica {
		for _, e := range []int{http.StatusUnauthorized, http.StatusNotFound, http.StatusForbidden} {
			if !slices.Contains(l, e) {
				l = append(l, e)
			}
		}
	}
	switch r.Metodo {
	case http.MethodPatch, http.MethodDelete:
		l = append(l, http.StatusConflict, http.StatusPreconditionFailed, http.StatusPreconditionRequired)
	case http.MethodPost:
		if !r.Publica || r.Caminho == "/v1/empresas" {
			l = append(l, http.StatusConflict)
		}
		if strings.HasSuffix(r.Caminho, "/desativar") || strings.HasSuffix(r.Caminho, "/reativar") {
			// Ação de estado: If-Match obrigatório (seção 2.11).
			l = append(l, http.StatusPreconditionFailed, http.StatusPreconditionRequired)
		}
	}
	slices.Sort(l)
	return slices.Compact(l)
}

func rotaDe(op *huma.Operation) *Rota {
	if op == nil || op.Metadata == nil {
		return nil
	}
	r, _ := op.Metadata["rota"].(*Rota)
	return r
}

type chaveAtor struct{}

// AtorDe devolve o ator resolvido pela credencial nesta requisição. Só existe
// em rota não pública; sem ele, a rota não deveria ter chegado ao handler.
func AtorDe(ctx context.Context) (autorizacao.Ator, error) {
	a, ok := ctx.Value(chaveAtor{}).(autorizacao.Ator)
	if !ok {
		return autorizacao.Ator{}, erros.NaoAutenticado
	}
	return a, nil
}

// --------------------------------------------------------------- erros

func init() {
	// Todo erro da API sai no formato RFC 9457 com o campo `codigo` estável
	// (seção 3 de docs/3_ESPECIFICACAO_APIS.md). Sem esta troca, o roteador
	// usaria o modelo de erro padrão da biblioteca, que não tem `codigo`.
	huma.NewError = func(status int, msg string, errs ...error) huma.StatusError {
		return montarProblema(nil, status, msg, errs...)
	}
	huma.NewErrorWithContext = func(ctx huma.Context, status int, msg string, errs ...error) huma.StatusError {
		return montarProblema(ctx, status, msg, errs...)
	}
}

// montarProblema traduz o erro (do caso de uso ou da validação da
// biblioteca) para o corpo RFC 9457.
func montarProblema(ctx huma.Context, status int, _ string, errs ...error) *erros.Problema {
	var p erros.Problema
	if dominio := erroDeDominio(errs); dominio != nil {
		p = erros.ParaProblema(dominio, "")
		if dominio.Status >= http.StatusInternalServerError && ctx != nil {
			obs.Logger(ctx.Context()).Error("erro na requisição", "codigo", dominio.Codigo, "erro", errors.Unwrap(dominio))
		}
	} else if pr := problemaPronto(errs); pr != nil {
		p = *pr
	} else {
		p = problemaDaBiblioteca(ctx, status, errs)
	}
	if ctx != nil {
		u := ctx.URL()
		if p.Instance == "" {
			p.Instance = u.Path
		}
		p.IDRastreio = ctx.Header(cabecalhoRequestID)
	}
	return &p
}

func erroDeDominio(errs []error) *erros.Erro {
	for _, e := range errs {
		var d *erros.Erro
		if e != nil && errors.As(e, &d) {
			return d
		}
	}
	for _, e := range errs {
		if e != nil && !isDetalhe(e) {
			// Erro técnico sem tradução: erro_interno, com a causa no log.
			var pr *erros.Problema
			if !errors.As(e, &pr) {
				return erros.Interno.Por(e)
			}
		}
	}
	return nil
}

func problemaPronto(errs []error) *erros.Problema {
	for _, e := range errs {
		var pr *erros.Problema
		if e != nil && errors.As(e, &pr) {
			return pr
		}
	}
	return nil
}

func isDetalhe(e error) bool {
	var d *huma.ErrorDetail
	return errors.As(e, &d)
}

// somenteLeitura são campos que nenhuma rota de escrita aceita: tentar mudá-los
// é 422 campo_somente_leitura, e não 400 (seção 2.1 do documento 3).
var somenteLeitura = map[string]bool{
	"id": true, "versao": true, "situacao": true, "status": true, "numero": true,
	"criado_em": true, "atualizado_em": true, "criado_por_id": true, "plano": true,
	"srv_habilitado": true, "escopo_acesso": true, "email": true, "usuario_id": true,
	"permissoes_avulsas": true, "teste_ate": true, "retencao_auditoria_anos": true,
}

var (
	reObrigatorio  = regexp.MustCompile(`^expected required property (\S+) to be present$`)
	reTamanhoMin   = regexp.MustCompile(`^expected length >= (\d+)$`)
	reTamanhoMax   = regexp.MustCompile(`^expected length <= (\d+)$`)
	reItensMin     = regexp.MustCompile(`^expected array length >= (\d+)$`)
	reItensMax     = regexp.MustCompile(`^expected array length <= (\d+)$`)
	reNumeroMin    = regexp.MustCompile(`^expected number >=? (\S+)$`)
	reNumeroMax    = regexp.MustCompile(`^expected number <=? (\S+)$`)
	reUmDe         = regexp.MustCompile(`^expected value to be one of "(.*)"$`)
	reTipo         = regexp.MustCompile(`^expected (string|integer|number|boolean|array|object|null)$`)
	reFormatoData  = regexp.MustCompile(`RFC 3339 date-time`)
	reFormatoDia   = regexp.MustCompile(`RFC 3339 date$`)
	reFormatoUUID  = regexp.MustCompile(`RFC 4122 uuid`)
	reFormatoEmail = regexp.MustCompile(`RFC 5322 email`)
)

// problemaDaBiblioteca monta o problema de um erro gerado pela própria
// biblioteca: validação do contrato, corpo malformado, tipo de conteúdo,
// limite do corpo, rota sem método.
func problemaDaBiblioteca(_ huma.Context, status int, errs []error) erros.Problema {
	base := erros.PorStatus(status)
	var (
		campos                      []erros.ErroCampo
		naoPermitido, desconhecido  bool
		filtro, tipoErrado, leitura bool
		malformado                  bool
	)
	for _, e := range errs {
		var d *huma.ErrorDetail
		if e == nil || !errors.As(e, &d) {
			continue
		}
		local, msg := d.Location, d.Message
		campo := campoDoLocal(local)
		c := erros.ErroCampo{Campo: campo, Codigo: erros.CampoValorInvalido, Mensagem: "Valor inválido."}
		switch {
		case msg == "unknown query parameter":
			filtro = true
			c.Mensagem = "Parâmetro de consulta não aceito nesta rota."
		case msg == "unexpected property":
			ultimo := campo
			if i := strings.LastIndexAny(campo, ".]"); i >= 0 {
				ultimo = campo[i+1:]
			}
			switch {
			case ultimo == "empresa_id":
				naoPermitido = true
				c.Mensagem = "A empresa vem da credencial; não envie empresa_id."
			case somenteLeitura[ultimo]:
				leitura = true
				c.Mensagem = "Campo somente leitura."
			default:
				desconhecido = true
				c.Mensagem = "Campo desconhecido nesta rota."
			}
		case reObrigatorio.MatchString(msg):
			m := reObrigatorio.FindStringSubmatch(msg)
			c.Campo = juntarCampo(campo, m[1])
			c.Codigo, c.Mensagem = erros.CampoObrigatorio, "Campo obrigatório."
		case reTamanhoMin.MatchString(msg):
			c.Codigo, c.Mensagem = erros.CampoTamanhoMinimo, fmt.Sprintf("Use pelo menos %s caracteres.", reTamanhoMin.FindStringSubmatch(msg)[1])
		case reTamanhoMax.MatchString(msg):
			c.Codigo, c.Mensagem = erros.CampoTamanhoMaximo, fmt.Sprintf("Use no máximo %s caracteres.", reTamanhoMax.FindStringSubmatch(msg)[1])
		case reItensMin.MatchString(msg):
			c.Codigo, c.Mensagem = erros.CampoTamanhoMinimo, fmt.Sprintf("Informe pelo menos %s itens.", reItensMin.FindStringSubmatch(msg)[1])
		case reItensMax.MatchString(msg):
			c.Codigo, c.Mensagem = erros.CampoTamanhoMaximo, fmt.Sprintf("Informe no máximo %s itens.", reItensMax.FindStringSubmatch(msg)[1])
		case reNumeroMin.MatchString(msg):
			c.Codigo, c.Mensagem = erros.CampoForaDaFaixa, fmt.Sprintf("O valor mínimo é %s.", reNumeroMin.FindStringSubmatch(msg)[1])
		case reNumeroMax.MatchString(msg):
			c.Codigo, c.Mensagem = erros.CampoForaDaFaixa, fmt.Sprintf("O valor máximo é %s.", reNumeroMax.FindStringSubmatch(msg)[1])
		case reUmDe.MatchString(msg):
			c.Mensagem = "Valor não aceito."
			if vs := reUmDe.FindStringSubmatch(msg)[1]; vs != "" {
				c.Mensagem = "Valor não aceito. Use um destes: " + vs + "."
			}
		case reTipo.MatchString(msg):
			tipoErrado = true
			c.Mensagem = "Tipo de valor errado: esperado " + nomeTipo(reTipo.FindStringSubmatch(msg)[1]) + "."
		case reFormatoData.MatchString(msg):
			c.Codigo, c.Mensagem = erros.CampoFormatoInvalido, "Use data e hora ISO 8601 em UTC, com Z: 2026-11-12T13:45:00Z."
		case reFormatoDia.MatchString(msg):
			c.Codigo, c.Mensagem = erros.CampoFormatoInvalido, "Use a data no formato AAAA-MM-DD."
		case reFormatoUUID.MatchString(msg):
			c.Codigo, c.Mensagem = erros.CampoFormatoInvalido, "Use um identificador UUID."
		case reFormatoEmail.MatchString(msg):
			c.Codigo, c.Mensagem = erros.CampoFormatoInvalido, "Informe um e-mail válido."
		case strings.HasPrefix(msg, "expected string to match pattern"):
			c.Codigo, c.Mensagem = erros.CampoFormatoInvalido, "Formato inválido."
		case local == "body" && status == http.StatusBadRequest:
			// Corpo que não é JSON (ou formulário) válido. A mensagem da
			// biblioteca e o valor recebido não voltam: o corpo pode ter senha.
			malformado = true
			c.Campo, c.Mensagem = "", "O corpo da requisição não é um JSON válido."
		}
		campos = append(campos, c)
	}

	detalhe := ""
	switch {
	case naoPermitido:
		base, detalhe = erros.CampoNaoPermitido, "A empresa vem sempre da credencial; remova empresa_id da requisição."
	case filtro:
		base, detalhe = erros.FiltroNaoSuportado, "A rota não aceita um dos parâmetros de consulta; confira o contrato."
	case malformado:
		base, detalhe = erros.RequisicaoInvalida, "O corpo da requisição não é um JSON válido."
		campos = nil
	case desconhecido || tipoErrado:
		base, detalhe = erros.RequisicaoInvalida, "A requisição tem campo desconhecido ou valor de tipo errado; confira o contrato."
	case leitura:
		base, detalhe = erros.CampoSomenteLeitura, "A requisição tenta alterar campo somente leitura."
	case status == http.StatusUnprocessableEntity:
		base = erros.EntradaInvalida
		if len(campos) == 1 {
			detalhe = "1 campo com problema."
		} else {
			detalhe = fmt.Sprintf("%d campos com problema.", len(campos))
		}
	case status == http.StatusRequestEntityTooLarge:
		detalhe = "O corpo passa do limite aceito pela rota."
	case status == http.StatusUnsupportedMediaType:
		detalhe = "Envie o corpo como application/json."
	case status == http.StatusBadRequest:
		detalhe = "A requisição não pôde ser lida; confira o contrato."
	case status == http.StatusNotFound || status == http.StatusMethodNotAllowed:
		detalhe = "A rota não existe."
	}
	e := base.Com("%s", detalhe)
	for _, c := range campos {
		e = e.ComCampo(c.Campo, c.Codigo, c.Mensagem)
	}
	return erros.ParaProblema(e, "")
}

func nomeTipo(t string) string {
	switch t {
	case "string":
		return "texto"
	case "integer":
		return "número inteiro"
	case "number":
		return "número"
	case "boolean":
		return "verdadeiro ou falso"
	case "array":
		return "lista"
	case "object":
		return "objeto"
	}
	return t
}

// campoDoLocal tira o prefixo do local (body., query., path., header.).
func campoDoLocal(local string) string {
	for _, p := range []string{"body.", "query.", "path.", "header."} {
		if strings.HasPrefix(local, p) {
			return local[len(p):]
		}
	}
	if local == "body" {
		return ""
	}
	return local
}

func juntarCampo(base, nome string) string {
	if base == "" {
		return nome
	}
	return base + "." + nome
}
