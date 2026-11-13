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
	Acao       string `query:"acao" maxLength:"40" example:"alteracao" doc:"Tipo de operação: login, logout, inclusao, alteracao, exclusao, delecao, visualizacao, impressao ou extracao"`
	Modulo     string `query:"modulo" maxLength:"3" example:"adm" doc:"Código do módulo, em minúsculas"`
	Tela       string `query:"tela" maxLength:"80" example:"adm.pessoas" doc:"<modulo>.<tela> do catálogo fechado"`
	IP         string `query:"ip" maxLength:"45" example:"203.0.113.7" doc:"Endereço exato; exige ocorrido_em[gte]"`
	De         string `query:"ocorrido_em[gte]" format:"date-time" doc:"A partir deste instante, inclusive"`
	Ate        string `query:"ocorrido_em[lt]" format:"date-time" doc:"Antes deste instante"`
	Cursor     string `query:"cursor" maxLength:"512"`
	Limite     int32  `query:"limite" minimum:"1" maximum:"200" default:"50"`
}

type entradaExportarEventos struct {
	Entidade   string `query:"entidade" maxLength:"60"`
	EntidadeID string `query:"entidade_id" format:"uuid"`
	UsuarioID  string `query:"usuario_id" format:"uuid"`
	Acao       string `query:"acao" maxLength:"40"`
	Modulo     string `query:"modulo" maxLength:"3"`
	Tela       string `query:"tela" maxLength:"80"`
	IP         string `query:"ip" maxLength:"45"`
	De         string `query:"ocorrido_em[gte]" format:"date-time" doc:"Obrigatório"`
	Ate        string `query:"ocorrido_em[lt]" format:"date-time" doc:"Obrigatório"`
}

type saidaCSV struct {
	ContentType        string `header:"Content-Type"`
	ContentDisposition string `header:"Content-Disposition"`
	Body               []byte
}

type entradaVerificar struct {
	De  string `query:"de" format:"date" doc:"Primeiro dia UTC, inclusive"`
	Ate string `query:"ate" format:"date" doc:"Último dia UTC, inclusive; no máximo 400 dias depois de de"`
}

type entradaAlterarCaptura struct {
	Tela string `path:"tela" example:"adm.pessoas" doc:"<modulo>.<tela> do catálogo fechado"`
	Body EntradaCapturaAuditoria
}

func filtroDaConsulta(in entradaListarEventos) (adm.FiltroAuditoria, error) {
	f := adm.FiltroAuditoria{Entidade: textoOpcional(in.Entidade), Acao: textoOpcional(in.Acao), Modulo: textoOpcional(in.Modulo),
		Tela: textoOpcional(in.Tela), IP: textoOpcional(in.IP), Cursor: in.Cursor, Limite: in.Limite}
	var err error
	if f.EntidadeID, err = uuidDoFiltro("entidade_id", in.EntidadeID); err != nil {
		return f, err
	}
	if f.UsuarioID, err = uuidDoFiltro("usuario_id", in.UsuarioID); err != nil {
		return f, err
	}
	if f.De, err = instanteDoFiltro("ocorrido_em[gte]", in.De); err != nil {
		return f, err
	}
	f.Ate, err = instanteDoFiltro("ocorrido_em[lt]", in.Ate)
	return f, err
}

func diaDoFiltro(campo, s string) (time.Time, error) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Time{}, erros.EntradaInvalida.ComCampoChave(campo, erros.CampoFormatoInvalido, "campo.data_formato", nil)
	}
	return t, nil
}

func instanteDoFiltro(campo, s string) (*time.Time, error) {
	if s == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, erros.EntradaInvalida.ComCampoChave(campo, erros.CampoFormatoInvalido, "campo.data_hora_z", nil)
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
		Resumo: "Trilha de auditoria",
		Descricao: "Filtros por usuário, tela ou módulo, tipo de operação (acao), período, IP, entidade e registro; todos juntos " +
			"valem como E. Paginada por cursor, do mais recente (RF-ADM-011). O filtro por IP exige ocorrido_em[gte].",
	}), func(ctx context.Context, in *entradaListarEventos) (*saidaLista[ResumoEvento], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		f, err := filtroDaConsulta(*in)
		if err != nil {
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
		Resumo: "Evento de auditoria",
		Descricao: "Com antes, depois, detalhe, hash e sensivel. O evento é imutável. Abrir o detalhe é uma visualizacao, gravada " +
			"quando a captura da tela de origem está ligada.",
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

	api.Registrar(s, rotaADM(api.Rota{
		ID: "exportarEventosAuditoria", Metodo: http.MethodGet, Caminho: "/v1/exportacoes-eventos-auditoria",
		Permissoes: []string{autorizacao.AdmLerAuditoria}, TipoSaida: "text/csv",
		Resumo: "Exportar a trilha em CSV",
		Descricao: "As mesmas linhas e os mesmos filtros de GET /v1/eventos-auditoria, do mais recente ao mais antigo, em CSV UTF-8. " +
			"Exige o período inteiro (ocorrido_em[gte] e ocorrido_em[lt]) e no máximo 50.000 eventos. A exportação é uma extracao: " +
			"grava um evento quando a captura de extração da tela de origem está ligada, inclusive na tela de auditoria (RF-ADM-011).",
	}), func(ctx context.Context, in *entradaExportarEventos) (*saidaCSV, error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		f, err := filtroDaConsulta(entradaListarEventos{Entidade: in.Entidade, EntidadeID: in.EntidadeID, UsuarioID: in.UsuarioID,
			Acao: in.Acao, Modulo: in.Modulo, Tela: in.Tela, IP: in.IP, De: in.De, Ate: in.Ate})
		if err != nil {
			return nil, err
		}
		x, err := svc.ExportarEventosAuditoria(ctx, a, f)
		if err != nil {
			return nil, err
		}
		return &saidaCSV{ContentType: "text/csv; charset=utf-8",
			ContentDisposition: `attachment; filename="trilha-de-auditoria.csv"`, Body: x.CSV}, nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "verificarIntegridadeAuditoria", Metodo: http.MethodGet, Caminho: "/v1/verificacoes-auditoria",
		Permissoes: []string{autorizacao.AdmLerAuditoria}, Erros: []int{http.StatusServiceUnavailable},
		Resumo: "Verificar a integridade da trilha",
		Descricao: "Para cada dia UTC do período (de..ate, inclusive, até 400 dias), recalcula no banco o hash dos eventos, o selo e a " +
			"cadeia, e compara o selo com a cópia no armazenamento imutável. Estados: integro, violado (motivo e primeiro evento " +
			"alterado, quando há como apontar) e nao_selado (dia que ainda não tem selo). 503 servico_indisponivel se o " +
			"armazenamento do selo não responde: sem a cópia não se afirma integridade (RF-ADM-011).",
	}), func(ctx context.Context, in *entradaVerificar) (*saidaLista[DiaVerificado], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		de, err := diaDoFiltro("de", in.De)
		if err != nil {
			return nil, err
		}
		ate, err := diaDoFiltro("ate", in.Ate)
		if err != nil {
			return nil, err
		}
		dias, err := svc.VerificarIntegridadeAuditoria(ctx, a, de, ate)
		if err != nil {
			return nil, err
		}
		return &saidaLista[DiaVerificado]{Body: api.ListaCompleta(dias, paraDiaVerificado)}, nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "listarCapturasAuditoria", Metodo: http.MethodGet, Caminho: "/v1/capturas-auditoria",
		Permissoes: []string{autorizacao.AdmGerenciarParametros},
		Resumo:     "Captura da trilha por tela",
		Descricao: "Todas as telas do catálogo, por módulo, com o nome da tela no idioma da resposta (Accept-Language, depois o idioma do usuário e o da empresa; Content-Language informa qual), com as chaves de visualização, impressão e extração da empresa (o padrão " +
			"onde ela não mudou nada). Login, logout, inclusão, alteração, exclusão, deleção e a visualização de dado sensível " +
			"nunca se desligam e não têm chave. Lista curta, inteira numa página.",
	}), func(ctx context.Context, _ *semEntrada) (*saidaIdioma[api.Lista[CapturaAuditoria]], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		l, err := svc.ListarCapturasAuditoria(ctx, a)
		if err != nil {
			return nil, err
		}
		return comIdioma(ctx, api.ListaCompleta(l, paraCapturaAuditoria)), nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "alterarCapturaAuditoria", Metodo: http.MethodPut, Caminho: "/v1/capturas-auditoria/{tela}",
		Permissoes: []string{autorizacao.AdmGerenciarParametros},
		Resumo:     "Ligar ou desligar a captura de uma tela",
		Descricao: "Grava as três chaves da tela (todas obrigatórias). A troca é uma alteração sempre gravada na trilha e só vale para " +
			"eventos futuros. 422 validacao_falhou se a tela não está no catálogo.",
	}), func(ctx context.Context, in *entradaAlterarCaptura) (*saidaIdioma[CapturaAuditoria], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		c, err := svc.AlterarCapturaAuditoria(ctx, a, in.Tela, adm.EntradaCaptura{
			Visualizacao: in.Body.Visualizacao, Impressao: in.Body.Impressao, Extracao: in.Body.Extracao,
		})
		if err != nil {
			return nil, err
		}
		return comIdioma(ctx, paraCapturaAuditoria(c)), nil
	})
}
