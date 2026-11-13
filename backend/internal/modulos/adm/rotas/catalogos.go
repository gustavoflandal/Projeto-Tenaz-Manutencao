package rotas

import (
	"context"
	"net/http"
	"time"

	"tenaz/internal/modulos/adm"
	"tenaz/internal/plataforma/api"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/erros"
)

type entradaListarEventos struct {
	Entidade   string `query:"entidade" maxLength:"60" example:"usuario"`
	EntidadeID string `query:"entidade_id" format:"uuid"`
	UsuarioID  string `query:"usuario_id" format:"uuid"`
	Acao       string `query:"acao" maxLength:"40" example:"alteracao"`
	De         string `query:"ocorrido_em[gte]" format:"date-time" doc:"A partir deste instante, inclusive"`
	Ate        string `query:"ocorrido_em[lt]" format:"date-time" doc:"Antes deste instante"`
	Cursor     string `query:"cursor" maxLength:"512"`
	Limite     int32  `query:"limite" minimum:"1" maximum:"200" default:"50"`
}

func instanteDoFiltro(campo, s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, erros.EntradaInvalida.ComCampo(campo, erros.CampoFormatoInvalido, "Use data e hora ISO 8601 em UTC, com Z.")
	}
	return &t, nil
}

func (r *rotas) registrarCatalogos(s *api.Servidor) {
	api.Registrar(s, rotaADM(api.Rota{
		ID: "listarPerfis", Metodo: http.MethodGet, Caminho: "/v1/perfis", Permissoes: []string{autorizacao.AdmLerUsuarios},
		Resumo: "Listar perfis", Descricao: "Os 12 perfis padrão e os personalizados, com as permissões. Lista curta, inteira numa página.",
	}), func(ctx context.Context, _ *semEntrada) (*saidaLista[Perfil], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		l, err := svc.ListarPerfis(ctx, a)
		if err != nil {
			return nil, err
		}
		return &saidaLista[Perfil]{Body: api.ListaCompleta(l, paraPerfil)}, nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "listarPermissoes", Metodo: http.MethodGet, Caminho: "/v1/permissoes", Permissoes: []string{autorizacao.AdmGerenciarPerfis},
		Resumo: "Catálogo de permissões", Descricao: "O catálogo inteiro de RF-ADM-005, com fase e marcação de sensível. Lista curta, inteira numa página.",
	}), func(ctx context.Context, _ *semEntrada) (*saidaLista[Permissao], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		l, err := svc.ListarPermissoes(ctx, a)
		if err != nil {
			return nil, err
		}
		return &saidaLista[Permissao]{Body: api.ListaCompleta(l, paraPermissao)}, nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "listarEventosAuditoria", Metodo: http.MethodGet, Caminho: "/v1/eventos-auditoria", Permissoes: []string{autorizacao.AdmLerAuditoria},
		Resumo:    "Trilha de auditoria",
		Descricao: "Filtros por entidade, registro, usuário, ação e período; paginada por cursor, do mais recente (RF-ADM-011).",
	}), func(ctx context.Context, in *entradaListarEventos) (*saidaLista[ResumoEvento], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		f := adm.FiltroAuditoria{Entidade: textoOpcional(in.Entidade), Acao: textoOpcional(in.Acao), Cursor: in.Cursor, Limite: in.Limite}
		if f.EntidadeID, err = uuidDoFiltro("entidade_id", in.EntidadeID); err != nil {
			return nil, err
		}
		if f.UsuarioID, err = uuidDoFiltro("usuario_id", in.UsuarioID); err != nil {
			return nil, err
		}
		if f.De, err = instanteDoFiltro("ocorrido_em[gte]", in.De); err != nil {
			return nil, err
		}
		if f.Ate, err = instanteDoFiltro("ocorrido_em[lt]", in.Ate); err != nil {
			return nil, err
		}
		p, err := svc.ListarEventosAuditoria(ctx, a, f)
		if err != nil {
			return nil, err
		}
		return &saidaLista[ResumoEvento]{Body: api.DePagina(p, paraResumoEvento)}, nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "obterEventoAuditoria", Metodo: http.MethodGet, Caminho: "/v1/eventos-auditoria/{id}", Permissoes: []string{autorizacao.AdmLerAuditoria},
		Resumo: "Evento de auditoria", Descricao: "Com os valores antes e depois. O evento é imutável.",
	}), func(ctx context.Context, in *entradaID) (*saida[Evento], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		id, err := idDoCaminho(in.ID)
		if err != nil {
			return nil, err
		}
		e, err := svc.ObterEventoAuditoria(ctx, a, id)
		if err != nil {
			return nil, err
		}
		return &saida[Evento]{Body: paraEvento(e)}, nil
	})
}
