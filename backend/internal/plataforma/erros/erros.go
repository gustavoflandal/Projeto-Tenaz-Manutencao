// Package erros define os erros de domínio da Tenaz e a tradução para o
// formato RFC 9457 (application/problem+json).
//
// Regra da seção 8.2 de docs/5_ARQUITETURA_SOFTWARE.md: erros de domínio são
// tipados e têm código estável. O código nunca muda, porque integradores
// passam a depender dele; a mensagem é em português e diz o que fazer.
package erros

import (
	"errors"
	"fmt"
	"net/http"
)

// BaseTipo é o prefixo da URL que documenta cada código de erro.
const BaseTipo = "https://docs.tenaz.com.br/erros/"

// Erro é um erro de domínio com código estável.
type Erro struct {
	Codigo  string // estável, em MAIUSCULA_COM_SUBLINHADO
	Titulo  string // curto, em português
	Status  int    // status HTTP
	Detalhe string // o que aconteceu e o que fazer
	causa   error
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

// Tipo é a URL de documentação do código.
func (e *Erro) Tipo() string {
	return BaseTipo + slug(e.Codigo)
}

func slug(codigo string) string {
	b := make([]byte, 0, len(codigo))
	for i := 0; i < len(codigo); i++ {
		c := codigo[i]
		switch {
		case c == '_':
			b = append(b, '-')
		case c >= 'A' && c <= 'Z':
			b = append(b, c+('a'-'A'))
		default:
			b = append(b, c)
		}
	}
	return string(b)
}

func novo(codigo, titulo string, status int) *Erro {
	return &Erro{Codigo: codigo, Titulo: titulo, Status: status}
}

// Erros da plataforma. Os erros de cada módulo ficam no pacote do módulo.
var (
	// NaoEncontrado é a resposta para recurso inexistente E para recurso de
	// outra empresa ou fora do escopo. Nunca 403 nesse caso: 403 confirmaria
	// a existência do recurso (docs/3_ESPECIFICACAO_APIS.md, seção 3).
	NaoEncontrado = novo("NAO_ENCONTRADO", "Recurso não encontrado", http.StatusNotFound)

	NaoAutenticado  = novo("NAO_AUTENTICADO", "Autenticação necessária", http.StatusUnauthorized)
	SemPermissao    = novo("SEM_PERMISSAO", "Sem permissão para esta ação", http.StatusForbidden)
	EntradaInvalida = novo("ENTRADA_INVALIDA", "Dados inválidos", http.StatusUnprocessableEntity)
	Conflito        = novo("CONFLITO", "Conflito de versão", http.StatusPreconditionFailed)
	LimiteExcedido  = novo("LIMITE_EXCEDIDO", "Limite de requisições excedido", http.StatusTooManyRequests)
	Interno         = novo("ERRO_INTERNO", "Erro interno", http.StatusInternalServerError)

	// SemContextoEmpresa é falha de programação: alguém tocou dado de cliente
	// fora de NaEmpresa. Falha fechada, com 500, e nunca devolve dado.
	SemContextoEmpresa = novo("SEM_CONTEXTO_EMPRESA", "Erro interno", http.StatusInternalServerError)
)

// Problema é o corpo de resposta no formato RFC 9457.
//
// O campo `codigo` é a extensão da Tenaz e é o que integradores usam para
// tratar o erro em código: ele é estável e nunca muda, enquanto `title` e
// `detail` são texto para pessoas e podem ser reescritos.
type Problema struct {
	Type     string  `json:"type" doc:"URI que documenta o código do erro" example:"https://docs.tenaz.com.br/erros/os-transicao-invalida"`
	Title    string  `json:"title" doc:"Resumo curto do tipo de erro" example:"Transição de estado inválida"`
	Status   int     `json:"status" doc:"Código de status HTTP" example:"409"`
	Codigo   string  `json:"codigo" doc:"Código estável do erro, para tratamento em código" example:"OS_TRANSICAO_INVALIDA"`
	Detail   string  `json:"detail,omitempty" doc:"O que aconteceu e o que fazer" example:"A OS OS-2026-004812 está cancelada e não pode ser concluída."`
	Instance string  `json:"instance,omitempty" doc:"Caminho da requisição que falhou" example:"/v1/ordens-servico/0193f2a1-7c3e-7b10-9a2e-5d1c3f0a8b77/concluir"`
	Campos   []Campo `json:"campos,omitempty" doc:"Erros por campo, quando a entrada é inválida"`
}

// Campo é um erro ligado a um campo específico da entrada.
type Campo struct {
	Local    string `json:"local" doc:"Onde o erro ocorreu" example:"body.prioridade"`
	Mensagem string `json:"mensagem" doc:"O que está errado com o campo"`
}

// Error faz Problema satisfazer a interface de erro, para que a camada HTTP
// possa devolvê-lo diretamente.
func (p *Problema) Error() string { return p.Codigo + ": " + p.Detail }

// GetStatus devolve o status HTTP. É o que o roteador usa para escolher o
// código da resposta.
func (p *Problema) GetStatus() int { return p.Status }

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

// Add acrescenta erros de campo ao problema.
func (p *Problema) Add(err error) {
	p.Campos = append(p.Campos, Campo{Mensagem: err.Error()})
}

// ParaProblema traduz um erro para o corpo RFC 9457.
// Erro desconhecido vira erro interno, sem expor a mensagem original.
func ParaProblema(err error, instancia string) Problema {
	var e *Erro
	if !errors.As(err, &e) {
		e = Interno
	}
	return Problema{
		Type:     e.Tipo(),
		Title:    e.Titulo,
		Status:   e.Status,
		Codigo:   e.Codigo,
		Detail:   e.Detalhe,
		Instance: instancia,
	}
}
