// Package auditoria grava os eventos da trilha que não nascem de uma escrita
// em tabela, e por isso não passam pelo trg_auditoria: visualização,
// impressão, extração e a alteração no nível do fato (seção 9.1 de
// docs/2_ARQUITETURA_BANCO_DADOS.md; RF-ADM-011). Login, logout e a falha de
// login são gravados pelo papel da autenticação (ADR-011).
//
// Também é dono da captura por tela (captura.go), do catálogo fechado de
// telas (catalogo.go), do selo diário e da verificação de integridade
// (selo.go, verificacao.go).
//
// O evento é gravado na transação do contexto, aberta por NaEmpresa: se o fato
// for desfeito, o evento também é. Nunca leva senha, hash, token nem segredo,
// nem o conteúdo que foi mostrado, impresso ou extraído.
package auditoria

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"

	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/db/consultas"
)

// Tipos de operação aceitos por ck_evento_auditoria_acao (lista fechada,
// sumário 3.6).
const (
	AcaoLogin        = "login"
	AcaoLogout       = "logout"
	AcaoInclusao     = "inclusao"
	AcaoAlteracao    = "alteracao"
	AcaoExclusao     = "exclusao"
	AcaoDelecao      = "delecao"
	AcaoVisualizacao = "visualizacao"
	AcaoImpressao    = "impressao"
	AcaoExtracao     = "extracao"
)

// Evento é um evento gravado pelo serviço.
type Evento struct {
	EmpresaID       uuid.UUID
	AtorTipo        string // db.AtorUsuario, db.AtorSistema...
	UsuarioID       *uuid.UUID
	Acao            string
	Entidade        string // nome da tabela ou do recurso
	EntidadeID      *uuid.UUID
	EntidadeChave   string
	Antes           any // serializado em JSON; nil grava nulo
	Depois          any
	CamposAlterados []string
	// Sensivel: o registro mostrado, impresso ou extraído tem dado sensível
	// (CPF completo etc.). Só vale em visualizacao, impressao e extracao.
	Sensivel bool
	// Detalhe: filtros, formato e quantidade de uma impressão ou extração.
	// Nunca o conteúdo.
	Detalhe any
}

// Registrar grava o evento na transação do contexto.
func Registrar(ctx context.Context, e Evento) error {
	q, err := db.Q(ctx)
	if err != nil {
		return err
	}
	antes, err := jsonOuNulo(e.Antes)
	if err != nil {
		return err
	}
	depois, err := jsonOuNulo(e.Depois)
	if err != nil {
		return err
	}
	detalhe, err := jsonOuNulo(e.Detalhe)
	if err != nil {
		return err
	}
	var chave *string
	if e.EntidadeChave != "" {
		chave = &e.EntidadeChave
	}
	if e.AtorTipo == "" {
		e.AtorTipo = db.AtorSistema
	}
	err = q.InserirEventoAuditoria(ctx, consultas.InserirEventoAuditoriaParams{
		EmpresaID:       e.EmpresaID,
		AtorTipo:        e.AtorTipo,
		UsuarioID:       e.UsuarioID,
		Acao:            e.Acao,
		Entidade:        e.Entidade,
		EntidadeID:      e.EntidadeID,
		EntidadeChave:   chave,
		Antes:           antes,
		Depois:          depois,
		CamposAlterados: e.CamposAlterados,
		Sensivel:        e.Sensivel,
		Detalhe:         detalhe,
	})
	if err != nil {
		return fmt.Errorf("gravar evento de auditoria %s/%s: %w", e.Entidade, e.Acao, err)
	}
	return nil
}

func jsonOuNulo(v any) ([]byte, error) {
	if v == nil {
		return nil, nil
	}
	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("serializar evento de auditoria: %w", err)
	}
	return b, nil
}
