// Package erros define os erros de domínio da Tenaz e a tradução para o
// formato RFC 9457 (application/problem+json).
//
// O dono dos códigos é a seção 3 de docs/3_ESPECIFICACAO_APIS.md: o código é
// estável, em snake_case e português, e nunca muda de significado em /v1,
// porque integradores passam a depender dele; a mensagem é em português e diz
// o que fazer. Código novo entra primeiro na tabela da seção 3.2 do documento.
package erros

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

// BaseTipo é o prefixo da URI que documenta cada código de erro (seção 3.1
// do documento 3). O domínio depende da pendência P5 do sumário.
const BaseTipo = "https://api.tenaz.com.br/erros/"

// Erro é um erro de domínio com código estável.
type Erro struct {
	Codigo  string // estável, em snake_case (seção 3.2 do documento 3)
	Titulo  string // curto, em português, fixo por código
	Status  int    // status HTTP
	Detalhe string // o que aconteceu e o que fazer

	// Campos lista os erros por campo da entrada (membro `erros` da seção
	// 3.1 do documento 3: `campo`, `codigo` do vocabulário curto e
	// `mensagem`).
	Campos []ErroCampo

	// Dados são os membros adicionais da RFC 9457 que o código exige, como
	// `plano_necessario` em plano_nao_inclui_recurso ou `estado_atual` em
	// transicao_invalida (seção 3.1 do documento 3). Saem no corpo do
	// problema, no nível de cima.
	Dados map[string]any

	causa error
}

// ErroCampo é um erro ligado a um campo da entrada.
type ErroCampo struct {
	Campo    string // notação de ponto: "turnos[0].hora_inicio"
	Codigo   string // obrigatorio, valor_invalido, tamanho_minimo, formato_invalido, duplicado, nao_encontrado, inativo...
	Mensagem string
}

// Vocabulário curto e estável dos erros por campo (seção 3.2 do documento 3).
const (
	CampoObrigatorio     = "obrigatorio"
	CampoValorInvalido   = "valor_invalido"
	CampoTamanhoMaximo   = "tamanho_maximo"
	CampoTamanhoMinimo   = "tamanho_minimo"
	CampoFormatoInvalido = "formato_invalido"
	CampoForaDaFaixa     = "fora_da_faixa"
	CampoNaoEncontrado   = "nao_encontrado"
	CampoDuplicado       = "duplicado"
	CampoInativo         = "inativo"
)

// ComCampo devolve uma cópia do erro com mais um erro de campo.
func (e *Erro) ComCampo(campo, codigo, mensagem string) *Erro {
	copia := *e
	copia.Campos = append(append([]ErroCampo(nil), e.Campos...), ErroCampo{Campo: campo, Codigo: codigo, Mensagem: mensagem})
	return &copia
}

// ComDado devolve uma cópia do erro com um membro adicional.
func (e *Erro) ComDado(chave string, valor any) *Erro {
	copia := *e
	copia.Dados = make(map[string]any, len(e.Dados)+1)
	for k, v := range e.Dados {
		copia.Dados[k] = v
	}
	copia.Dados[chave] = valor
	return &copia
}

// Is faz errors.Is comparar pelo código, e não pelo ponteiro: Com, Por e
// ComCampo devolvem cópias, e errors.Is(err, erros.NaoEncontrado) tem de
// continuar verdadeiro para elas.
func (e *Erro) Is(alvo error) bool {
	a, ok := alvo.(*Erro)
	return ok && a.Codigo == e.Codigo
}

func (e *Erro) Error() string {
	if e.Detalhe != "" {
		return fmt.Sprintf("%s: %s", e.Codigo, e.Detalhe)
	}
	return e.Codigo
}

func (e *Erro) Unwrap() error { return e.causa }

// Com devolve uma cópia do erro com o detalhe preenchido.
func (e *Erro) Com(detalhe string, args ...any) *Erro {
	copia := *e
	copia.Detalhe = fmt.Sprintf(detalhe, args...)
	return &copia
}

// Por devolve uma cópia do erro embrulhando a causa técnica. A causa não vai
// para a resposta HTTP: fica no log.
func (e *Erro) Por(causa error) *Erro {
	copia := *e
	copia.causa = causa
	return &copia
}

// Tipo é a URI de documentação do código.
func (e *Erro) Tipo() string {
	return BaseTipo + e.Codigo
}

// Novo cria um erro de domínio. Os módulos declaram os seus com ele, como
// variáveis de pacote que nunca são alteradas (Com, Por e ComCampo devolvem
// cópias). O código precisa estar na seção 3.2 do documento 3.
func Novo(codigo, titulo string, status int) *Erro {
	return &Erro{Codigo: codigo, Titulo: titulo, Status: status}
}

// Erros transversais da seção 3.2 do documento 3. Os erros de cada módulo
// ficam no pacote do módulo.
var (
	// RequisicaoInvalida: JSON malformado, tipo errado, parâmetro ou campo
	// desconhecido.
	RequisicaoInvalida = Novo("requisicao_invalida", "Requisição inválida", http.StatusBadRequest)

	// CampoNaoPermitido: campo proibido no corpo ou na consulta, como
	// empresa_id (seção 4.7 do documento 3: a empresa vem só da credencial).
	CampoNaoPermitido = Novo("campo_nao_permitido", "Campo não permitido", http.StatusBadRequest)

	// FiltroNaoSuportado: filtro, ordenação ou parâmetro de consulta não
	// declarado no contrato da rota.
	FiltroNaoSuportado = Novo("filtro_nao_suportado", "Filtro não suportado", http.StatusBadRequest)

	CursorInvalido = Novo("cursor_invalido", "Cursor inválido", http.StatusBadRequest)

	NaoAutenticado = Novo("nao_autenticado", "Autenticação necessária", http.StatusUnauthorized)

	// CredencialInvalida: sessão, token, código ou senha inválidos, expirados
	// ou revogados. É também a resposta para e-mail inexistente e conta
	// bloqueada, para não revelar quais e-mails existem.
	CredencialInvalida = Novo("credencial_invalida", "Credencial inválida", http.StatusUnauthorized)

	// MFAObrigatorio: a senha está certa e falta o segundo fator. A camada
	// HTTP devolve `desafio_id` e `desafio_expira_em` junto.
	MFAObrigatorio = Novo("mfa_obrigatorio", "Segundo fator necessário", http.StatusUnauthorized)

	// MFACodigoInvalido: código TOTP errado ao ativar, desativar ou concluir o
	// login. No login, cada erro consome uma das 5 tentativas do desafio e a
	// resposta traz `tentativas_restantes`.
	MFACodigoInvalido = Novo("mfa_codigo_invalido", "Código do segundo fator inválido", http.StatusUnauthorized)

	CSRFInvalido = Novo("csrf_invalido", "Token CSRF ausente ou inválido", http.StatusForbidden)

	// SemPermissao: o perfil não permite a operação.
	SemPermissao = Novo("permissao_negada", "Sem permissão para esta ação", http.StatusForbidden)

	// PlanoNaoIncluiRecurso traz `plano_necessario` em Dados.
	PlanoNaoIncluiRecurso = Novo("plano_nao_inclui_recurso", "Recurso fora do plano contratado", http.StatusForbidden)

	// NaoEncontrado é a resposta para recurso inexistente E para recurso de
	// outra empresa ou fora do alcance. Nunca 403 nesse caso: 403 confirmaria
	// a existência do recurso (seções 3.2 e 4.7 do documento 3).
	NaoEncontrado = Novo("recurso_nao_encontrado", "Recurso não encontrado", http.StatusNotFound)

	// Conflito: conflito genérico de unicidade, com o campo em Campos.
	Conflito = Novo("conflito", "Registro duplicado", http.StatusConflict)

	// Duplicado é o mesmo código de Conflito, com o nome que os repositórios
	// usam ao traduzir violação de unicidade.
	Duplicado = Conflito

	// TransicaoInvalida traz `estado_atual` em Dados.
	TransicaoInvalida = Novo("transicao_invalida", "Ação não permitida no estado atual", http.StatusConflict)

	IdempotenciaEmProcessamento = Novo("idempotencia_em_processamento", "Requisição ainda em processamento", http.StatusConflict)

	// LinkEmailInvalido: link de convite, redefinição de senha ou verificação
	// de e-mail inexistente, vencido, já usado ou substituído por outro envio.
	LinkEmailInvalido = Novo("link_email_invalido", "Link inválido ou vencido", http.StatusGone)

	// VersaoDesatualizada: If-Match não confere com a versão gravada. A camada
	// HTTP devolve o recurso atual em `recurso_atual`.
	VersaoDesatualizada = Novo("versao_desatualizada", "O registro mudou desde a leitura", http.StatusPreconditionFailed)

	PayloadMuitoGrande    = Novo("payload_muito_grande", "Corpo acima do limite", http.StatusRequestEntityTooLarge)
	TipoMidiaNaoSuportado = Novo("tipo_midia_nao_suportado", "Tipo de conteúdo não aceito", http.StatusUnsupportedMediaType)

	// EntradaInvalida: erros de campo, com o detalhe em Campos.
	EntradaInvalida = Novo("validacao_falhou", "Os dados enviados não passaram na validação", http.StatusUnprocessableEntity)

	CampoSomenteLeitura  = Novo("campo_somente_leitura", "Campo somente leitura", http.StatusUnprocessableEntity)
	IdempotenciaConflito = Novo("idempotencia_conflito", "Chave de idempotência já usada com outro corpo", http.StatusUnprocessableEntity)

	// SenhaFraca: a senha não atende à política (mínimo de 10 caracteres,
	// fora da lista de senhas comuns e vazadas, sem repetir o e-mail).
	SenhaFraca = Novo("senha_fraca", "Senha fraca", http.StatusUnprocessableEntity)

	// PreCondicaoObrigatoria: PATCH, DELETE ou ação de estado sem If-Match.
	PreCondicaoObrigatoria = Novo("pre_condicao_obrigatoria", "Cabeçalho If-Match obrigatório", http.StatusPreconditionRequired)

	LimiteExcedido = Novo("limite_requisicoes_excedido", "Limite de requisições excedido", http.StatusTooManyRequests)

	Interno = Novo("erro_interno", "Erro interno", http.StatusInternalServerError)

	// SemContextoEmpresa é falha de programação: alguém tocou dado de cliente
	// fora de NaEmpresa. Falha fechada, com 500, e nunca devolve dado. O
	// código só existe no log: a resposta sai como erro_interno, porque 5xx
	// não expõe detalhe interno (seção 3.2 do documento 3).
	SemContextoEmpresa = Novo("sem_contexto_empresa", "Erro interno", http.StatusInternalServerError)

	// ServicoIndisponivel: dependência externa fora (SMTP, por exemplo).
	ServicoIndisponivel = Novo("servico_indisponivel", "Serviço indisponível", http.StatusServiceUnavailable)
)

// Problema é o corpo de resposta no formato RFC 9457 (seção 3.1 do
// documento 3).
//
// O campo `codigo` é a extensão da Tenaz e é o que integradores usam para
// tratar o erro em código: ele é estável e nunca muda, enquanto `title` e
// `detail` são texto para pessoas e podem ser reescritos.
//
// Os membros adicionais conhecidos estão declarados, para que o contrato os
// documente; um membro adicional ainda não declarado sai em Extras, no mesmo
// nível, e o esquema aceita propriedades adicionais.
type Problema struct {
	_ struct{} `json:"-" additionalProperties:"true"`

	Type  string `json:"type" doc:"URI que documenta o código do erro" example:"https://api.tenaz.com.br/erros/validacao_falhou"`
	Title string `json:"title" doc:"Resumo legível, fixo por código" example:"Os dados enviados não passaram na validação"`
	// int64, e não int: o tamanho de int muda com a arquitetura (32 bits num Go
	// de 32 bits), e o contrato gerado sairia int32 numa máquina e int64 noutra.
	// Todo número exposto no contrato usa tipo de tamanho fixo.
	Status     int64       `json:"status" doc:"Código de status HTTP" example:"422"`
	Detail     string      `json:"detail,omitempty" doc:"Explicação desta ocorrência e o que fazer"`
	Instance   string      `json:"instance,omitempty" doc:"Caminho da requisição" example:"/v1/pessoas"`
	Codigo     string      `json:"codigo" doc:"Código estável do erro, em snake_case; é o que o cliente trata" example:"validacao_falhou"`
	IDRastreio string      `json:"id_rastreio,omitempty" doc:"Igual ao X-Request-Id; o suporte usa para localizar o log"`
	Erros      []CampoErro `json:"erros,omitempty" doc:"Erros por campo, com o campo em notação de ponto"`

	// Membros adicionais (seção 3.1 do documento 3).
	DesafioID            *string    `json:"desafio_id,omitempty" doc:"mfa_obrigatorio: desafio do segundo fator, segredo de 256 bits e uso único; vai no corpo de POST /v1/auth/mfa/verificar"`
	DesafioExpiraEm      *time.Time `json:"desafio_expira_em,omitempty" doc:"mfa_obrigatorio: fim da validade do desafio (5 minutos)"`
	TentativasRestantes  *int32     `json:"tentativas_restantes,omitempty" doc:"mfa_codigo_invalido: tentativas que restam no desafio; em zero, o login recomeça pela senha"`
	PlanoNecessario      *string    `json:"plano_necessario,omitempty" doc:"plano_nao_inclui_recurso: plano mínimo do recurso"`
	EstadoAtual          *string    `json:"estado_atual,omitempty" doc:"transicao_invalida: estado atual do registro"`
	TransicoesPermitidas []string   `json:"transicoes_permitidas,omitempty" doc:"transicao_invalida: ações aceitas no estado atual"`
	VersaoAtual          *int64     `json:"versao_atual,omitempty" doc:"versao_desatualizada: versão gravada do registro"`
	RecursoAtual         any        `json:"recurso_atual,omitempty" doc:"versao_desatualizada: o recurso como está gravado, para o cliente mostrar a diferença"`
	CamposFaltantes      []string   `json:"campos_faltantes,omitempty" doc:"Campos exigidos que faltaram"`

	// Extras são membros adicionais sem declaração própria.
	Extras map[string]any `json:"-"`
}

// CampoErro é um erro ligado a um campo da entrada (membro `erros`).
type CampoErro struct {
	Campo          string   `json:"campo" doc:"Campo em notação de ponto" example:"responsavel.email"`
	Codigo         string   `json:"codigo" doc:"Vocabulário curto: obrigatorio, valor_invalido, tamanho_maximo, tamanho_minimo, formato_invalido, fora_da_faixa, nao_encontrado, duplicado, inativo" example:"obrigatorio"`
	Mensagem       string   `json:"mensagem" doc:"O que está errado e o que fazer" example:"Informe o e-mail."`
	ValoresAceitos []string `json:"valores_aceitos,omitempty" doc:"Valores aceitos, quando o campo é uma enumeração"`
}

// MarshalJSON acrescenta os Extras no nível de cima do objeto.
func (p Problema) MarshalJSON() ([]byte, error) {
	type semMetodo Problema
	base, err := json.Marshal(semMetodo(p))
	if err != nil || len(p.Extras) == 0 {
		return base, err
	}
	var m map[string]any
	if err := json.Unmarshal(base, &m); err != nil {
		return nil, err
	}
	for k, v := range p.Extras {
		if _, existe := m[k]; !existe {
			m[k] = v
		}
	}
	return json.Marshal(m)
}

// Error faz Problema satisfazer a interface de erro, para que a camada HTTP
// possa devolvê-lo diretamente.
func (p *Problema) Error() string { return p.Codigo + ": " + p.Detail }

// GetStatus devolve o status HTTP. É o que o roteador usa para escolher o
// código da resposta.
func (p *Problema) GetStatus() int { return int(p.Status) }

// ContentType troca o tipo de conteúdo da resposta de erro para o exigido pela
// RFC 9457. Sem isto, o erro sairia como application/json, e um cliente que
// distingue problema de resposta normal pelo tipo de conteúdo não o
// reconheceria.
func (p *Problema) ContentType(tipo string) string {
	if tipo == "application/json" {
		return "application/problem+json"
	}
	return tipo
}

// Add acrescenta um erro de campo sem campo conhecido.
func (p *Problema) Add(err error) {
	p.Erros = append(p.Erros, CampoErro{Campo: "", Codigo: CampoValorInvalido, Mensagem: err.Error()})
}

// ParaProblema traduz um erro para o corpo RFC 9457.
//
// Erro desconhecido vira erro_interno, sem expor a mensagem original. Erro
// de domínio com status 500 também sai como erro_interno: 5xx nunca expõe
// detalhe interno (seção 3.2 do documento 3). O código original fica no log.
func ParaProblema(err error, instancia string) Problema {
	var e *Erro
	if !errors.As(err, &e) {
		e = Interno
	}
	if e.Status == http.StatusInternalServerError {
		e = Interno
	}
	p := Problema{
		Type:     e.Tipo(),
		Title:    e.Titulo,
		Status:   int64(e.Status),
		Codigo:   e.Codigo,
		Detail:   e.Detalhe,
		Instance: instancia,
	}
	for _, c := range e.Campos {
		p.Erros = append(p.Erros, CampoErro{Campo: c.Campo, Codigo: c.Codigo, Mensagem: c.Mensagem})
	}
	for k, v := range e.Dados {
		p.definirMembro(k, v)
	}
	return p
}

// definirMembro põe o membro adicional no campo declarado, quando houver, ou
// em Extras.
func (p *Problema) definirMembro(chave string, valor any) {
	switch chave {
	case "desafio_id":
		if s, ok := valor.(string); ok {
			p.DesafioID = &s
			return
		}
	case "desafio_expira_em":
		if t, ok := valor.(time.Time); ok {
			t = t.UTC().Truncate(time.Second)
			p.DesafioExpiraEm = &t
			return
		}
	case "tentativas_restantes":
		if n, ok := inteiro(valor); ok {
			v := int32(n)
			p.TentativasRestantes = &v
			return
		}
	case "plano_necessario":
		if s, ok := valor.(string); ok {
			p.PlanoNecessario = &s
			return
		}
	case "estado_atual":
		if s, ok := valor.(string); ok {
			p.EstadoAtual = &s
			return
		}
	case "transicoes_permitidas":
		if l, ok := valor.([]string); ok {
			p.TransicoesPermitidas = l
			return
		}
	case "versao_atual":
		if n, ok := inteiro(valor); ok {
			p.VersaoAtual = &n
			return
		}
	case "recurso_atual":
		p.RecursoAtual = valor
		return
	case "campos_faltantes":
		if l, ok := valor.([]string); ok {
			p.CamposFaltantes = l
			return
		}
	}
	if p.Extras == nil {
		p.Extras = map[string]any{}
	}
	p.Extras[chave] = valor
}

func inteiro(v any) (int64, bool) {
	switch n := v.(type) {
	case int:
		return int64(n), true
	case int32:
		return int64(n), true
	case int64:
		return n, true
	}
	return 0, false
}

// PorStatus devolve o erro transversal que corresponde a um status HTTP, para
// respostas que a biblioteca de roteamento gera sozinha (rota inexistente,
// corpo malformado, tipo de conteúdo não aceito).
func PorStatus(status int) *Erro {
	switch status {
	case http.StatusBadRequest:
		return RequisicaoInvalida
	case http.StatusUnauthorized:
		return NaoAutenticado
	case http.StatusForbidden:
		return SemPermissao
	case http.StatusNotFound, http.StatusMethodNotAllowed:
		return NaoEncontrado
	case http.StatusConflict:
		return Conflito
	case http.StatusGone:
		return LinkEmailInvalido
	case http.StatusPreconditionFailed:
		return VersaoDesatualizada
	case http.StatusRequestEntityTooLarge:
		return PayloadMuitoGrande
	case http.StatusUnsupportedMediaType:
		return TipoMidiaNaoSuportado
	case http.StatusUnprocessableEntity:
		return EntradaInvalida
	case http.StatusPreconditionRequired:
		return PreCondicaoObrigatoria
	case http.StatusTooManyRequests:
		return LimiteExcedido
	case http.StatusServiceUnavailable:
		return ServicoIndisponivel
	}
	return Interno
}
