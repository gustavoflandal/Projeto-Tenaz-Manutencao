package auditoria

import (
	"context"

	"github.com/google/uuid"

	"tenaz/internal/plataforma/db"
)

// Operacao descreve uma visualização, impressão ou extração a registrar. Não
// leva o conteúdo mostrado: só o que foi acessado (entidade, registro, filtros).
type Operacao struct {
	EmpresaID     uuid.UUID
	UsuarioID     *uuid.UUID
	Entidade      string
	EntidadeID    *uuid.UUID
	EntidadeChave string
	// Sensivel: mostra dado sensível (CPF completo, dados de sinistro, CPF de
	// locatário pessoa física). Visualização sensível é sempre gravada.
	Sensivel bool
	// Detalhe: filtros, formato e quantidade. Nunca o conteúdo.
	Detalhe any
	// Sempre ignora a captura da tela: o pedido e o download da exportação
	// completa (RF-ADM-013).
	Sempre bool
}

// Visualizou registra a abertura de um registro, na transação do contexto.
// Grava quando o dado é sensível, quando o ator é a equipe de suporte da Tenaz
// ou quando a captura de visualização da tela de origem está ligada. A tela
// vem da requisição (cabeçalho X-Tenaz-Tela, validado contra o catálogo); sem
// tela, vale o padrão (visualização desligada). Devolve se gravou.
func (s *Servico) Visualizou(ctx context.Context, o Operacao) (bool, error) {
	return s.registrar(ctx, TipoVisualizacao, o)
}

// Imprimiu registra a impressão (PDF, tela) quando a captura da tela liga.
func (s *Servico) Imprimiu(ctx context.Context, o Operacao) (bool, error) {
	return s.registrar(ctx, TipoImpressao, o)
}

// Extraiu registra a extração (lista, relatório, exportação, inclusive a da
// própria trilha) quando a captura da tela liga.
func (s *Servico) Extraiu(ctx context.Context, o Operacao) (bool, error) {
	return s.registrar(ctx, TipoExtracao, o)
}

func (s *Servico) registrar(ctx context.Context, tipo Tipo, o Operacao) (bool, error) {
	req := db.RequisicaoDe(ctx)
	grava := o.Sempre
	if tipo == TipoVisualizacao && (o.Sensivel || req.AtorTipo == db.AtorSuporteTenaz) {
		grava = true
	}
	if !grava {
		c, err := s.ChavesDaTela(ctx, o.EmpresaID, req.TelaDoModulo())
		if err != nil {
			return false, err
		}
		grava = c.liga(tipo)
	}
	if !grava {
		return false, nil
	}
	ator := req.AtorTipo
	if ator == "" {
		ator = db.AtorSistema
		if o.UsuarioID != nil {
			ator = db.AtorUsuario
		}
	}
	err := Registrar(ctx, Evento{
		EmpresaID: o.EmpresaID, AtorTipo: ator, UsuarioID: o.UsuarioID, Acao: string(tipo),
		Entidade: o.Entidade, EntidadeID: o.EntidadeID, EntidadeChave: o.EntidadeChave,
		Sensivel: o.Sensivel, Detalhe: o.Detalhe,
	})
	return err == nil, err
}
