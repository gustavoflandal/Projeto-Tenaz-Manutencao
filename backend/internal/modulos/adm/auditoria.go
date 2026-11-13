package adm

import (
	"context"
	"time"

	"github.com/google/uuid"

	"tenaz/internal/modulos/adm/internal/repo"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/db/consultas"
	"tenaz/internal/plataforma/erros"
	"tenaz/internal/plataforma/paginacao"
)

// Trilha de auditoria: RF-ADM-011. Só leitura; nenhum perfil altera evento.

var acoesAuditoria = []string{"inclusao", "alteracao", "exclusao", "login", "login_falho", "logout", "exportacao",
	"visualizacao_sensivel", "acesso_suporte", "revisao", "sincronizacao"}

func (s *servico) ListarEventosAuditoria(ctx context.Context, a autorizacao.Ator, f FiltroAuditoria) (paginacao.Pagina[ResumoEvento], error) {
	var pag paginacao.Pagina[ResumoEvento]
	if err := apenasInterno(a, autorizacao.AdmLerAuditoria); err != nil {
		return pag, err
	}
	if f.Acao != nil {
		ok := false
		for _, ac := range acoesAuditoria {
			ok = ok || ac == *f.Acao
		}
		if !ok {
			return pag, erros.EntradaInvalida.ComCampo("acao", erros.CampoValorInvalido, "Ação desconhecida.")
		}
	}
	if f.De != nil && f.Ate != nil && !f.De.Before(*f.Ate) {
		return pag, erros.EntradaInvalida.ComCampo("ate", erros.CampoValorInvalido, "O fim do período tem de ser depois do início.")
	}
	instante := func(t *time.Time) string {
		if t == nil {
			return ""
		}
		return t.UTC().Format(time.RFC3339Nano)
	}
	consulta := paginacao.Assinatura("eventos-auditoria",
		textoOuVazio(f.Entidade, func(s string) string { return s }),
		textoOuVazio(f.EntidadeID, uuid.UUID.String), textoOuVazio(f.UsuarioID, uuid.UUID.String),
		textoOuVazio(f.Acao, func(s string) string { return s }), instante(f.De), instante(f.Ate))
	agora := s.agora()
	pos, err := paginacao.Decodificar(f.Cursor, consulta, agora)
	if err != nil {
		return pag, err
	}
	limite := paginacao.Limite(f.Limite)
	cEm, cID := cursorParams(pos)
	err = s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		linhas, err := q.AdmListarEventosAuditoria(ctx, consultas.AdmListarEventosAuditoriaParams{
			EmpresaID: a.EmpresaID, Entidade: f.Entidade, EntidadeID: f.EntidadeID, UsuarioID: f.UsuarioID,
			Acao: f.Acao, De: f.De, Ate: f.Ate, CursorOcorridoEm: cEm, CursorID: cID, Limite: limite + 1,
		})
		if err != nil {
			return err
		}
		res := make([]ResumoEvento, 0, len(linhas))
		for _, r := range linhas {
			res = append(res, ResumoEvento{
				ID: r.ID, OcorridoEm: r.OcorridoEm, AtorTipo: r.AtorTipo, UsuarioID: r.UsuarioID, TokenAPIID: r.TokenApiID,
				Acao: r.Acao, Entidade: r.Entidade, EntidadeID: r.EntidadeID, EntidadeChave: r.EntidadeChave,
				CamposAlterados: r.CamposAlterados, Origem: r.Origem, IP: vazioParaNil(r.Ip), DispositivoID: r.DispositivoID,
				RequestID: r.RequestID, ConexaoEmpresaID: r.ConexaoEmpresaID, OcorridoNoDispositivoEm: r.OcorridoNoDispositivoEm,
			})
		}
		pag = paginacao.Montar(res, limite, consulta, agora, func(e ResumoEvento) paginacao.Posicao {
			return paginacao.Posicao{Instante: e.OcorridoEm, ID: e.ID}
		})
		return nil
	})
	return pag, err
}

func (s *servico) ObterEventoAuditoria(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (Evento, error) {
	if err := apenasInterno(a, autorizacao.AdmLerAuditoria); err != nil {
		return Evento{}, err
	}
	var ev Evento
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		r, err := q.AdmBuscarEventoAuditoria(ctx, consultas.AdmBuscarEventoAuditoriaParams{EmpresaID: a.EmpresaID, ID: id})
		if err != nil {
			return naoEncontrado(err)
		}
		ev = Evento{
			ResumoEvento: ResumoEvento{
				ID: r.ID, OcorridoEm: r.OcorridoEm, AtorTipo: r.AtorTipo, UsuarioID: r.UsuarioID, TokenAPIID: r.TokenApiID,
				Acao: r.Acao, Entidade: r.Entidade, EntidadeID: r.EntidadeID, EntidadeChave: r.EntidadeChave,
				CamposAlterados: r.CamposAlterados, Origem: r.Origem, IP: vazioParaNil(r.Ip), DispositivoID: r.DispositivoID,
				RequestID: r.RequestID, ConexaoEmpresaID: r.ConexaoEmpresaID, OcorridoNoDispositivoEm: r.OcorridoNoDispositivoEm,
			},
			Antes: r.Antes, Depois: r.Depois, UserAgent: r.UserAgent,
		}
		return nil
	})
	return ev, err
}
