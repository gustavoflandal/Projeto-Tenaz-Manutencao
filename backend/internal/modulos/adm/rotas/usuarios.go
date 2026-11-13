package rotas

import (
	"context"
	"net/http"
	"net/url"

	"tenaz/internal/modulos/adm"
	"tenaz/internal/plataforma/api"
	"tenaz/internal/plataforma/autorizacao"
)

type entradaListarUsuarios struct {
	Busca    string `query:"busca" maxLength:"100" doc:"Trecho do nome ou do e-mail"`
	Situacao string `query:"situacao" enum:"convidado,ativo,inativo"`
	Perfil   string `query:"perfil" maxLength:"40" doc:"Código do perfil de origem"`
	Cursor   string `query:"cursor" maxLength:"512"`
	Limite   int32  `query:"limite" minimum:"1" maximum:"200" default:"50"`
}

type entradaRetirarPermissao struct {
	ID        string `path:"id"`
	Permissao string `path:"permissao" maxLength:"60" example:"os.encerrar"`
}

func (r *rotas) registrarUsuarios(s *api.Servidor) {
	ler := []string{autorizacao.AdmLerUsuarios}
	gerenciar := []string{autorizacao.AdmGerenciarUsuarios}

	api.Registrar(s, rotaADM(api.Rota{
		ID: "listarUsuarios", Metodo: http.MethodGet, Caminho: "/v1/usuarios", Permissoes: ler,
		Resumo: "Listar usuários", Descricao: "Paginada por cursor, ordem -criado_em.",
	}), func(ctx context.Context, in *entradaListarUsuarios) (*saidaLista[ResumoUsuario], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		p, err := svc.ListarUsuarios(ctx, a, adm.FiltroUsuarios{
			Busca: textoOpcional(in.Busca), Situacao: textoOpcional(in.Situacao), Perfil: textoOpcional(in.Perfil),
			Cursor: in.Cursor, Limite: in.Limite,
		})
		if err != nil {
			return nil, err
		}
		return &saidaLista[ResumoUsuario]{Body: api.DePagina(p, paraResumoUsuario)}, nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "criarUsuario", Metodo: http.MethodPost, Caminho: "/v1/usuarios", Status: http.StatusCreated, Permissoes: gerenciar,
		Resumo: "Criar usuário e enviar convite",
		Descricao: "O usuário nasce convidado e o convite (7 dias) vai por e-mail (RF-ADM-004). Regras anti-escalada: só " +
			"se atribui perfil cujas permissões se tem (403 permissao_negada). prestador exige fornecedor_id; locatario, " +
			"locatario_id. 409 conflito em e-mail já usado na empresa.",
	}), func(ctx context.Context, in *entradaCorpo[EntradaUsuario]) (*saida[UsuarioCriado], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		b := in.Body
		e := adm.EntradaUsuario{
			Nome: b.Nome, Email: b.Email, Telefone: b.Telefone, Perfis: b.Perfis, FornecedorID: b.FornecedorID,
			LocatarioID: b.LocatarioID, PessoaID: b.PessoaID, CustoHora: b.CustoHora, UnidadePadraoID: b.UnidadePadraoID,
		}
		if b.NovaPessoa != nil {
			p := dePessoa(*b.NovaPessoa)
			e.NovaPessoa = &p
		}
		u, err := svc.CriarUsuario(ctx, a, e)
		if err != nil {
			return nil, err
		}
		corpo := UsuarioCriado{Usuario: paraUsuario(u.Usuario), ConviteEnviado: u.ConviteEnviado}
		return criado(corpo, u.Usuario.Versao, "/v1/usuarios/"+u.Usuario.ID.String()), nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "obterUsuarioEu", Metodo: http.MethodGet, Caminho: "/v1/usuarios/eu", Escopos: api.EscoposDoProprioUsuario,
		Resumo: "Próprio usuário", Descricao: "Usuário da credencial, com perfis e permissões avulsas. Vale para os perfis externos. Traz ETag.",
	}), func(ctx context.Context, _ *semEntrada) (*saida[Usuario], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		u, err := svc.ObterEu(ctx, a)
		if err != nil {
			return nil, err
		}
		return comETag(paraUsuario(u), u.Versao), nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "alterarUsuarioEu", Metodo: http.MethodPatch, Caminho: "/v1/usuarios/eu", Escopos: api.EscoposDoProprioUsuario,
		Resumo: "Alterar os próprios dados", Descricao: "Nome, telefone e idioma (pt-BR ou es-419): só o próprio usuário altera o idioma. JSON Merge Patch com If-Match.",
	}), func(ctx context.Context, in *entradaPatch[PatchEu]) (*saida[Usuario], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		esperada, err := api.VersaoEsperada(in.IfMatch)
		if err != nil {
			return nil, err
		}
		atual, err := svc.ObterEu(ctx, a)
		if err != nil {
			return nil, err
		}
		if err := api.ConferirVersao(esperada, atual.Versao, func() any { return paraUsuario(atual) }); err != nil {
			return nil, err
		}
		var n camposEu
		if _, err := api.MesclarPatch(camposEu{Nome: atual.Nome, Telefone: atual.Telefone, Idioma: atual.Idioma}, in.RawBody, &n); err != nil {
			return nil, err
		}
		u, err := svc.EditarEu(ctx, a, adm.EntradaEu{Nome: n.Nome, Telefone: n.Telefone, Idioma: &n.Idioma})
		if err != nil {
			return nil, err
		}
		return comETag(paraUsuario(u), u.Versao), nil
	})

	obter := func(ctx context.Context, svc adm.Servico, a autorizacao.Ator, idTexto string) (adm.Usuario, error) {
		id, err := idDoCaminho(idTexto)
		if err != nil {
			return adm.Usuario{}, err
		}
		return svc.ObterUsuario(ctx, a, id)
	}
	conferido := func(ctx context.Context, svc adm.Servico, a autorizacao.Ator, idTexto, ifMatch string) (adm.Usuario, error) {
		esperada, err := api.VersaoEsperada(ifMatch)
		if err != nil {
			return adm.Usuario{}, err
		}
		atual, err := obter(ctx, svc, a, idTexto)
		if err != nil {
			return adm.Usuario{}, err
		}
		return atual, api.ConferirVersao(esperada, atual.Versao, func() any { return paraUsuario(atual) })
	}

	api.Registrar(s, rotaADM(api.Rota{
		ID: "obterUsuario", Metodo: http.MethodGet, Caminho: "/v1/usuarios/{id}", Permissoes: ler,
		Resumo: "Detalhar usuário", Descricao: "Com perfis e permissões avulsas. Traz ETag.",
	}), func(ctx context.Context, in *entradaID) (*saida[Usuario], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		u, err := obter(ctx, svc, a, in.ID)
		if err != nil {
			return nil, err
		}
		return comETag(paraUsuario(u), u.Versao), nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "alterarUsuario", Metodo: http.MethodPatch, Caminho: "/v1/usuarios/{id}", Permissoes: gerenciar,
		Resumo: "Alterar perfis e dados do usuário",
		Descricao: "JSON Merge Patch com If-Match. perfis é a lista completa e substitui a atual. Ninguém altera os " +
			"próprios perfis; só administrador altera administrador; 409 ultimo_administrador ao tirar o perfil do " +
			"último administrador ativo. E-mail e vínculo externo não mudam (422 campo_somente_leitura no e-mail).",
	}), func(ctx context.Context, in *entradaPatchID[PatchUsuario]) (*saida[Usuario], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		atual, err := conferido(ctx, svc, a, in.ID, in.IfMatch)
		if err != nil {
			return nil, err
		}
		esperada := atual.Versao
		base := camposUsuario{
			Nome: atual.Nome, Telefone: atual.Telefone, CustoHora: precoTexto(atual.CustoHora),
			UnidadePadraoID: atual.UnidadePadraoID, PessoaID: atual.PessoaID,
		}
		var n camposUsuario
		if _, err := api.MesclarPatch(base, in.RawBody, &n); err != nil {
			return nil, err
		}
		u := atual
		// Os perfis vão primeiro: a regra anti-escalada recusa antes de
		// qualquer dado mudar. A troca de perfis não muda a versão do usuário.
		if in.Body.Perfis != nil {
			if u, err = svc.DefinirPerfis(ctx, a, atual.ID, in.Body.Perfis); err != nil {
				return nil, err
			}
		}
		if api.PatchTemCampoAlem(in.RawBody, "perfis") {
			u, err = svc.EditarUsuario(ctx, a, atual.ID, adm.EntradaEditarUsuario{
				Nome: n.Nome, Telefone: n.Telefone, CustoHora: n.CustoHora, UnidadePadraoID: n.UnidadePadraoID,
				PessoaID: n.PessoaID, VersaoEsperada: &esperada,
			})
			if err != nil {
				return nil, desatualizado(err, func() (any, int64, error) {
					x, e := svc.ObterUsuario(ctx, a, atual.ID)
					return paraUsuario(x), x.Versao, e
				})
			}
		}
		return comETag(paraUsuario(u), u.Versao), nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "desativarUsuario", Metodo: http.MethodPost, Caminho: "/v1/usuarios/{id}/desativar", Permissoes: gerenciar,
		Resumo: "Desativar usuário",
		Descricao: "situacao → inativo, com If-Match; as sessões e os refresh tokens do usuário caem na mesma transação. " +
			"409 ultimo_administrador; 409 transicao_invalida se já está inativo.",
	}), func(ctx context.Context, in *entradaIDIfMatch) (*saida[Usuario], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		atual, err := conferido(ctx, svc, a, in.ID, in.IfMatch)
		if err != nil {
			return nil, err
		}
		u, err := svc.InativarUsuario(ctx, a, atual.ID)
		if err != nil {
			return nil, err
		}
		return comETag(paraUsuario(u), u.Versao), nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "reativarUsuario", Metodo: http.MethodPost, Caminho: "/v1/usuarios/{id}/reativar", Permissoes: gerenciar,
		Resumo: "Reativar usuário",
		Descricao: "inativo → ativo, com If-Match; quem nunca aceitou o convite volta a convidado e recebe convite novo. " +
			"409 transicao_invalida se não está inativo.",
	}), func(ctx context.Context, in *entradaIDIfMatch) (*saida[UsuarioCriado], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		atual, err := conferido(ctx, svc, a, in.ID, in.IfMatch)
		if err != nil {
			return nil, err
		}
		u, err := svc.ReativarUsuario(ctx, a, atual.ID)
		if err != nil {
			return nil, err
		}
		return comETag(UsuarioCriado{Usuario: paraUsuario(u.Usuario), ConviteEnviado: u.ConviteEnviado}, u.Usuario.Versao), nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "reenviarConvite", Metodo: http.MethodPost, Caminho: "/v1/usuarios/{id}/reenviar-convite", Status: http.StatusNoContent, Permissoes: gerenciar,
		Resumo: "Reenviar o convite",
		Descricao: "Revoga o link anterior e envia outro (7 dias). 409 transicao_invalida se o usuário não está convidado; " +
			"503 servico_indisponivel se o e-mail não saiu.",
	}), func(ctx context.Context, in *entradaID) (*saidaVazia, error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		id, err := idDoCaminho(in.ID)
		if err != nil {
			return nil, err
		}
		return &saidaVazia{}, svc.ReenviarConvite(ctx, a, id)
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "encerrarSessoesUsuario", Metodo: http.MethodPost, Caminho: "/v1/usuarios/{id}/encerrar-sessoes", Permissoes: gerenciar,
		Resumo:    "Encerrar as sessões do usuário",
		Descricao: "Revoga as sessões web e do app e os refresh tokens do usuário nesta empresa; o logout vai para a trilha.",
	}), func(ctx context.Context, in *entradaID) (*saida[SessoesEncerradas], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		id, err := idDoCaminho(in.ID)
		if err != nil {
			return nil, err
		}
		n, err := svc.EncerrarSessoes(ctx, a, id)
		if err != nil {
			return nil, err
		}
		return &saida[SessoesEncerradas]{Body: SessoesEncerradas{SessoesEncerradas: n}}, nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "concederPermissaoUsuario", Metodo: http.MethodPost, Caminho: "/v1/usuarios/{id}/permissoes", Status: http.StatusCreated, Permissoes: gerenciar,
		Resumo: "Conceder permissão avulsa",
		Descricao: "Permissão do catálogo de RF-ADM-005 concedida a este usuário (usuario_permissao), como " +
			"os.encerrar. Ninguém concede a si mesmo; só concede quem tem a permissão, salvo as duas que o " +
			"administrador concede sem ter (os.encerrar e est.decidir_destino); usuário externo nunca recebe permissão vedada. A concessão vai para a trilha.",
	}), func(ctx context.Context, in *entradaIDCorpo[EntradaPermissao]) (*saida[Usuario], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		id, err := idDoCaminho(in.ID)
		if err != nil {
			return nil, err
		}
		u, err := svc.ConcederPermissao(ctx, a, id, in.Body.Permissao)
		if err != nil {
			return nil, err
		}
		return criado(paraUsuario(u), u.Versao, "/v1/usuarios/"+id.String()+"/permissoes/"+url.PathEscape(in.Body.Permissao)), nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "retirarPermissaoUsuario", Metodo: http.MethodDelete, Caminho: "/v1/usuarios/{id}/permissoes/{permissao}", Status: http.StatusNoContent, Permissoes: gerenciar,
		Resumo: "Retirar permissão avulsa", Descricao: "A retirada vai para a trilha. Permissão que o usuário não tem: 404.",
	}), func(ctx context.Context, in *entradaRetirarPermissao) (*saidaVazia, error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		id, err := idDoCaminho(in.ID)
		if err != nil {
			return nil, err
		}
		_, err = svc.RetirarPermissao(ctx, a, id, in.Permissao)
		return &saidaVazia{}, err
	})
}
