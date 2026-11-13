package rotas

import (
	"context"
	"net/http"

	"tenaz/internal/modulos/adm"
	"tenaz/internal/plataforma/api"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/erros"
)

type entradaTurnoID struct {
	ID      string `path:"id" doc:"Unidade"`
	TurnoID string `path:"turno_id"`
	IfMatch string `header:"If-Match" doc:"ETag do turno; obrigatório"`
}

type entradaPatchTurno struct {
	ID      string `path:"id" doc:"Unidade"`
	TurnoID string `path:"turno_id"`
	IfMatch string `header:"If-Match" doc:"ETag do turno; obrigatório"`
	Body    PatchTurno
	RawBody []byte
}

type entradaFeriadoID struct {
	ID        string `path:"id" doc:"Unidade"`
	FeriadoID string `path:"feriado_id"`
	IfMatch   string `header:"If-Match" doc:"ETag do feriado; obrigatório"`
}

type entradaListarFeriados struct {
	ID  string `path:"id" doc:"Unidade"`
	De  string `query:"data[gte]" format:"date" doc:"A partir desta data, inclusive"`
	Ate string `query:"data[lte]" format:"date" doc:"Até esta data, inclusive"`
}

func textoOpcional(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (r *rotas) registrarUnidades(s *api.Servidor) {
	ler := []string{autorizacao.AdmLer}
	gerenciar := []string{autorizacao.AdmGerenciarEmpresa}

	api.Registrar(s, rotaADM(api.Rota{
		ID: "listarUnidades", Metodo: http.MethodGet, Caminho: "/v1/unidades", Permissoes: ler,
		Resumo: "Listar unidades", Descricao: "Lista curta, inteira numa página (tem_mais falso), com as unidades ativas e inativas.",
	}), func(ctx context.Context, _ *semEntrada) (*saidaLista[ResumoUnidade], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		l, err := svc.ListarUnidades(ctx, a)
		if err != nil {
			return nil, err
		}
		return &saidaLista[ResumoUnidade]{Body: api.ListaCompleta(l, paraResumoUnidade)}, nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "criarUnidade", Metodo: http.MethodPost, Caminho: "/v1/unidades", Status: http.StatusCreated, Permissoes: gerenciar,
		Resumo: "Criar unidade",
		Descricao: "Com pelo menos um turno (RF-ADM-002). O Essencial tem 1 unidade: a segunda responde 403 " +
			"plano_nao_inclui_recurso com plano_necessario. 409 conflito em código repetido.",
	}), func(ctx context.Context, in *entradaCorpo[EntradaUnidade]) (*saida[Unidade], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		u, err := svc.CriarUnidade(ctx, a, deUnidade(in.Body))
		if err != nil {
			return nil, err
		}
		return criado(paraUnidade(u), u.Versao, "/v1/unidades/"+u.ID.String()), nil
	})

	obter := func(ctx context.Context, svc adm.Servico, a autorizacao.Ator, idTexto string) (adm.Unidade, error) {
		id, err := idDoCaminho(idTexto)
		if err != nil {
			return adm.Unidade{}, err
		}
		return svc.ObterUnidade(ctx, a, id)
	}

	api.Registrar(s, rotaADM(api.Rota{
		ID: "obterUnidade", Metodo: http.MethodGet, Caminho: "/v1/unidades/{id}", Permissoes: ler,
		Resumo: "Detalhar unidade", Descricao: "Com os turnos. Traz ETag.",
	}), func(ctx context.Context, in *entradaID) (*saida[Unidade], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		u, err := obter(ctx, svc, a, in.ID)
		if err != nil {
			return nil, err
		}
		return comETag(paraUnidade(u), u.Versao), nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "alterarUnidade", Metodo: http.MethodPatch, Caminho: "/v1/unidades/{id}", Permissoes: gerenciar,
		Resumo: "Alterar unidade", Descricao: "Nome, código, endereço, fuso, identificação (do tipo do país da empresa) e coordenadas. JSON Merge Patch com If-Match.",
	}), func(ctx context.Context, in *entradaPatchID[PatchUnidade]) (*saida[Unidade], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		esperada, err := api.VersaoEsperada(in.IfMatch)
		if err != nil {
			return nil, err
		}
		atual, err := obter(ctx, svc, a, in.ID)
		if err != nil {
			return nil, err
		}
		if err := api.ConferirVersao(esperada, atual.Versao, func() any { return paraUnidade(atual) }); err != nil {
			return nil, err
		}
		base := camposUnidade{
			Codigo: atual.Codigo, Nome: atual.Nome, IdentificacaoTipo: atual.IdentificacaoTipo, Identificacao: atual.Identificacao, Endereco: paraEndereco(atual.Endereco),
			FusoHorario: atual.FusoHorario, Latitude: atual.Latitude, Longitude: atual.Longitude,
		}
		var novo camposUnidade
		if _, err := api.MesclarPatch(base, in.RawBody, &novo); err != nil {
			return nil, err
		}
		u, err := svc.EditarUnidade(ctx, a, atual.ID, adm.EntradaUnidade{
			Codigo: novo.Codigo, Nome: novo.Nome, IdentificacaoTipo: novo.IdentificacaoTipo, Identificacao: novo.Identificacao, Endereco: deEndereco(novo.Endereco), FusoHorario: novo.FusoHorario,
			Latitude: novo.Latitude, Longitude: novo.Longitude, VersaoEsperada: &esperada,
		})
		if err != nil {
			return nil, desatualizado(err, func() (any, int64, error) {
				x, e := svc.ObterUnidade(ctx, a, atual.ID)
				return paraUnidade(x), x.Versao, e
			})
		}
		return comETag(paraUnidade(u), u.Versao), nil
	})

	conferida := func(ctx context.Context, svc adm.Servico, a autorizacao.Ator, idTexto, ifMatch string) (adm.Unidade, error) {
		esperada, err := api.VersaoEsperada(ifMatch)
		if err != nil {
			return adm.Unidade{}, err
		}
		atual, err := obter(ctx, svc, a, idTexto)
		if err != nil {
			return adm.Unidade{}, err
		}
		return atual, api.ConferirVersao(esperada, atual.Versao, func() any { return paraUnidade(atual) })
	}

	api.Registrar(s, rotaADM(api.Rota{
		ID: "inativarUnidade", Metodo: http.MethodDelete, Caminho: "/v1/unidades/{id}", Status: http.StatusNoContent, Permissoes: gerenciar,
		Resumo: "Inativar unidade",
		Descricao: "ativa → inativa, com If-Match; a unidade nunca é apagada (seção 2.12). A última unidade ativa não é " +
			"inativada (422 validacao_falhou). 409 transicao_invalida se já está inativa.",
	}), func(ctx context.Context, in *entradaIDIfMatch) (*saidaVazia, error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		u, err := conferida(ctx, svc, a, in.ID, in.IfMatch)
		if err != nil {
			return nil, err
		}
		_, err = svc.InativarUnidade(ctx, a, u.ID)
		return &saidaVazia{}, err
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "reativarUnidade", Metodo: http.MethodPost, Caminho: "/v1/unidades/{id}/reativar", Permissoes: gerenciar,
		Resumo:    "Reativar unidade",
		Descricao: "inativa → ativa, com If-Match; respeita o limite de unidades do plano. 409 transicao_invalida se já está ativa.",
	}), func(ctx context.Context, in *entradaIDIfMatch) (*saida[Unidade], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		atual, err := conferida(ctx, svc, a, in.ID, in.IfMatch)
		if err != nil {
			return nil, err
		}
		u, err := svc.ReativarUnidade(ctx, a, atual.ID)
		if err != nil {
			return nil, err
		}
		return comETag(paraUnidade(u), u.Versao), nil
	})

	// Turnos (turno_unidade).
	api.Registrar(s, rotaADM(api.Rota{
		ID: "incluirTurno", Metodo: http.MethodPost, Caminho: "/v1/unidades/{id}/turnos", Status: http.StatusCreated, Permissoes: gerenciar,
		Resumo: "Incluir turno", Descricao: "Turno por dia da semana no calendário da unidade (RF-ADM-002).",
	}), func(ctx context.Context, in *entradaIDCorpo[EntradaTurno]) (*saida[Turno], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		uid, err := idDoCaminho(in.ID)
		if err != nil {
			return nil, err
		}
		t, err := svc.IncluirTurno(ctx, a, uid, deTurno(in.Body))
		if err != nil {
			return nil, err
		}
		return criado(paraTurno(t), t.Versao, "/v1/unidades/"+uid.String()+"/turnos/"+t.ID.String()), nil
	})

	turnoAtual := func(ctx context.Context, svc adm.Servico, a autorizacao.Ator, idU, idT, ifMatch string) (adm.Turno, error) {
		esperada, err := api.VersaoEsperada(ifMatch)
		if err != nil {
			return adm.Turno{}, err
		}
		u, err := obter(ctx, svc, a, idU)
		if err != nil {
			return adm.Turno{}, err
		}
		tid, err := idDoCaminho(idT)
		if err != nil {
			return adm.Turno{}, err
		}
		for _, t := range u.Turnos {
			if t.ID == tid {
				return t, api.ConferirVersao(esperada, t.Versao, func() any { return paraTurno(t) })
			}
		}
		return adm.Turno{}, erros.NaoEncontrado
	}

	api.Registrar(s, rotaADM(api.Rota{
		ID: "alterarTurno", Metodo: http.MethodPatch, Caminho: "/v1/unidades/{id}/turnos/{turno_id}", Permissoes: gerenciar,
		Resumo: "Alterar turno", Descricao: "JSON Merge Patch com If-Match (ETag do turno).",
	}), func(ctx context.Context, in *entradaPatchTurno) (*saida[Turno], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		t, err := turnoAtual(ctx, svc, a, in.ID, in.TurnoID, in.IfMatch)
		if err != nil {
			return nil, err
		}
		var novo camposTurno
		if _, err := api.MesclarPatch(camposTurno{DiaSemana: t.DiaSemana, HoraInicio: t.HoraInicio, HoraFim: t.HoraFim, Nome: t.Nome}, in.RawBody, &novo); err != nil {
			return nil, err
		}
		alterado, err := svc.AlterarTurno(ctx, a, t.UnidadeID, t.ID, adm.EntradaTurno{DiaSemana: novo.DiaSemana, HoraInicio: novo.HoraInicio, HoraFim: novo.HoraFim, Nome: novo.Nome})
		if err != nil {
			return nil, err
		}
		return comETag(paraTurno(alterado), alterado.Versao), nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "removerTurno", Metodo: http.MethodDelete, Caminho: "/v1/unidades/{id}/turnos/{turno_id}", Status: http.StatusNoContent, Permissoes: gerenciar,
		Resumo: "Remover turno", Descricao: "Com If-Match. O último turno da unidade não é removido (422 validacao_falhou).",
	}), func(ctx context.Context, in *entradaTurnoID) (*saidaVazia, error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		t, err := turnoAtual(ctx, svc, a, in.ID, in.TurnoID, in.IfMatch)
		if err != nil {
			return nil, err
		}
		return &saidaVazia{}, svc.RemoverTurno(ctx, a, t.UnidadeID, t.ID)
	})

	// Feriados (feriado_unidade).
	api.Registrar(s, rotaADM(api.Rota{
		ID: "listarFeriados", Metodo: http.MethodGet, Caminho: "/v1/unidades/{id}/feriados", Permissoes: ler,
		Resumo: "Listar feriados da unidade", Descricao: "Lista curta, inteira numa página, com filtro de período opcional.",
	}), func(ctx context.Context, in *entradaListarFeriados) (*saidaLista[Feriado], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		uid, err := idDoCaminho(in.ID)
		if err != nil {
			return nil, err
		}
		l, err := svc.ListarFeriados(ctx, a, uid, textoOpcional(in.De), textoOpcional(in.Ate))
		if err != nil {
			return nil, err
		}
		return &saidaLista[Feriado]{Body: api.ListaCompleta(l, paraFeriado)}, nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "incluirFeriado", Metodo: http.MethodPost, Caminho: "/v1/unidades/{id}/feriados", Status: http.StatusCreated, Permissoes: gerenciar,
		Resumo: "Incluir feriado", Descricao: "409 conflito em data repetida na unidade.",
	}), func(ctx context.Context, in *entradaIDCorpo[EntradaFeriado]) (*saida[Feriado], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		uid, err := idDoCaminho(in.ID)
		if err != nil {
			return nil, err
		}
		f, err := svc.IncluirFeriado(ctx, a, uid, adm.EntradaFeriado{Data: in.Body.Data, Descricao: in.Body.Descricao})
		if err != nil {
			return nil, err
		}
		return criado(paraFeriado(f), f.Versao, "/v1/unidades/"+uid.String()+"/feriados/"+f.ID.String()), nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "removerFeriado", Metodo: http.MethodDelete, Caminho: "/v1/unidades/{id}/feriados/{feriado_id}", Status: http.StatusNoContent, Permissoes: gerenciar,
		Resumo: "Remover feriado", Descricao: "Com If-Match (ETag do feriado, que vem na lista).",
	}), func(ctx context.Context, in *entradaFeriadoID) (*saidaVazia, error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		esperada, err := api.VersaoEsperada(in.IfMatch)
		if err != nil {
			return nil, err
		}
		uid, err := idDoCaminho(in.ID)
		if err != nil {
			return nil, err
		}
		fid, err := idDoCaminho(in.FeriadoID)
		if err != nil {
			return nil, err
		}
		l, err := svc.ListarFeriados(ctx, a, uid, nil, nil)
		if err != nil {
			return nil, err
		}
		for _, f := range l {
			if f.ID == fid {
				if err := api.ConferirVersao(esperada, f.Versao, func() any { return paraFeriado(f) }); err != nil {
					return nil, err
				}
				return &saidaVazia{}, svc.RemoverFeriado(ctx, a, uid, fid)
			}
		}
		return nil, erros.NaoEncontrado
	})
}
