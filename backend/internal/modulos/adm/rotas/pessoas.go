package rotas

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"tenaz/internal/modulos/adm"
	"tenaz/internal/plataforma/api"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/erros"
)

type entradaListarPessoas struct {
	Busca        string `query:"busca" maxLength:"100" doc:"Trecho do nome"`
	Situacao     string `query:"situacao" enum:"ativa,inativa"`
	FornecedorID string `query:"fornecedor_id" format:"uuid"`
	UnidadeID    string `query:"unidade_id" format:"uuid"`
	ComUsuario   string `query:"com_usuario" enum:"true,false" doc:"true: só pessoas com usuário; false: só sem usuário"`
	Cursor       string `query:"cursor" maxLength:"512" doc:"proximo_cursor da página anterior"`
	Limite       int32  `query:"limite" minimum:"1" maximum:"200" default:"50"`
}

// uuidDoFiltro lê um filtro UUID opcional, já validado pelo contrato.
func uuidDoFiltro(campo, s string) (*uuid.UUID, error) {
	if s == "" {
		return nil, nil
	}
	id, err := uuid.Parse(s)
	if err != nil {
		return nil, erros.EntradaInvalida.ComCampo(campo, erros.CampoFormatoInvalido, "Use um identificador UUID.")
	}
	return &id, nil
}

func (r *rotas) registrarPessoas(s *api.Servidor) {
	ler := []string{autorizacao.AdmLerPessoas}
	gerenciar := []string{autorizacao.AdmGerenciarPessoas}

	api.Registrar(s, rotaADM(api.Rota{
		ID: "listarPessoas", Metodo: http.MethodGet, Caminho: "/v1/pessoas", Permissoes: ler,
		Resumo:    "Listar pessoas",
		Descricao: "Paginada por cursor, ordem -criado_em. O CPF sai mascarado; o inteiro só no detalhe.",
	}), func(ctx context.Context, in *entradaListarPessoas) (*saidaLista[ResumoPessoa], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		f := adm.FiltroPessoas{Busca: textoOpcional(in.Busca), Situacao: textoOpcional(in.Situacao), Cursor: in.Cursor, Limite: in.Limite}
		if f.FornecedorID, err = uuidDoFiltro("fornecedor_id", in.FornecedorID); err != nil {
			return nil, err
		}
		if f.UnidadeID, err = uuidDoFiltro("unidade_id", in.UnidadeID); err != nil {
			return nil, err
		}
		if in.ComUsuario != "" {
			v := in.ComUsuario == "true"
			f.ComUsuario = &v
		}
		p, err := svc.ListarPessoas(ctx, a, f)
		if err != nil {
			return nil, err
		}
		return &saidaLista[ResumoPessoa]{Body: api.DePagina(p, paraResumoPessoa)}, nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "criarPessoa", Metodo: http.MethodPost, Caminho: "/v1/pessoas", Status: http.StatusCreated, Permissoes: gerenciar,
		Resumo: "Cadastrar pessoa", Descricao: "Com ou sem usuário (RF-ADM-015). 409 conflito em CPF repetido na empresa.",
	}), func(ctx context.Context, in *entradaCorpo[EntradaPessoa]) (*saida[Pessoa], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		p, err := svc.CriarPessoa(ctx, a, dePessoa(in.Body))
		if err != nil {
			return nil, err
		}
		return criado(paraPessoa(p), p.Versao, "/v1/pessoas/"+p.ID.String()), nil
	})

	obter := func(ctx context.Context, svc adm.Servico, a autorizacao.Ator, idTexto string) (adm.Pessoa, error) {
		id, err := idDoCaminho(idTexto)
		if err != nil {
			return adm.Pessoa{}, err
		}
		return svc.ObterPessoa(ctx, a, id)
	}

	api.Registrar(s, rotaADM(api.Rota{
		ID: "obterPessoa", Metodo: http.MethodGet, Caminho: "/v1/pessoas/{id}", Permissoes: ler, AlcanceNoServico: true,
		Resumo: "Detalhar pessoa",
		Descricao: "Com o CPF inteiro; a leitura é registrada na trilha. O titular lê a própria pessoa sem " +
			"adm.ler_pessoas. Traz ETag.",
	}), func(ctx context.Context, in *entradaID) (*saida[Pessoa], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		p, err := obter(ctx, svc, a, in.ID)
		if err != nil {
			return nil, err
		}
		return comETag(paraPessoa(p), p.Versao), nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "alterarPessoa", Metodo: http.MethodPatch, Caminho: "/v1/pessoas/{id}", Permissoes: gerenciar,
		Resumo: "Alterar pessoa", Descricao: "JSON Merge Patch com If-Match.",
	}), func(ctx context.Context, in *entradaPatchID[PatchPessoa]) (*saida[Pessoa], error) {
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
		if err := api.ConferirVersao(esperada, atual.Versao, func() any { return paraPessoa(atual) }); err != nil {
			return nil, err
		}
		base := camposPessoa{
			Nome: atual.Nome, CPF: atual.CPF, DataNascimento: atual.DataNascimento, Matricula: atual.Matricula, Cargo: atual.Cargo,
			Telefone: atual.Telefone, Email: atual.Email, Vinculo: atual.Vinculo, FornecedorID: atual.FornecedorID,
			UnidadeID: atual.UnidadeID, RegistroProfissional: atual.RegistroProfissional,
		}
		var n camposPessoa
		if _, err := api.MesclarPatch(base, in.RawBody, &n); err != nil {
			return nil, err
		}
		p, err := svc.EditarPessoa(ctx, a, atual.ID, adm.EntradaPessoa{
			Nome: n.Nome, CPF: n.CPF, DataNascimento: n.DataNascimento, Matricula: n.Matricula, Cargo: n.Cargo, Telefone: n.Telefone,
			Email: n.Email, Vinculo: n.Vinculo, FornecedorID: n.FornecedorID, UnidadeID: n.UnidadeID,
			RegistroProfissional: n.RegistroProfissional, VersaoEsperada: &esperada,
		})
		if err != nil {
			return nil, desatualizado(err, func() (any, int64, error) {
				x, e := svc.ObterPessoa(ctx, a, atual.ID)
				return paraPessoa(x), x.Versao, e
			})
		}
		return comETag(paraPessoa(p), p.Versao), nil
	})

	conferida := func(ctx context.Context, svc adm.Servico, a autorizacao.Ator, idTexto, ifMatch string) (adm.Pessoa, error) {
		esperada, err := api.VersaoEsperada(ifMatch)
		if err != nil {
			return adm.Pessoa{}, err
		}
		atual, err := obter(ctx, svc, a, idTexto)
		if err != nil {
			return adm.Pessoa{}, err
		}
		return atual, api.ConferirVersao(esperada, atual.Versao, func() any { return paraPessoa(atual) })
	}

	api.Registrar(s, rotaADM(api.Rota{
		ID: "inativarPessoa", Metodo: http.MethodDelete, Caminho: "/v1/pessoas/{id}", Status: http.StatusNoContent, Permissoes: gerenciar,
		Resumo: "Inativar pessoa",
		Descricao: "A pessoa nunca é apagada: é inativada (seção 2.12). Com If-Match. 409 transicao_invalida se a " +
			"pessoa já está inativa ou se está ligada a usuário ativo ou convidado: o usuário é inativado antes " +
			"(RN-ADM-021).",
	}), func(ctx context.Context, in *entradaIDIfMatch) (*saidaVazia, error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		p, err := conferida(ctx, svc, a, in.ID, in.IfMatch)
		if err != nil {
			return nil, err
		}
		_, err = svc.InativarPessoa(ctx, a, p.ID)
		return &saidaVazia{}, err
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "reativarPessoa", Metodo: http.MethodPost, Caminho: "/v1/pessoas/{id}/reativar", Permissoes: gerenciar,
		Resumo: "Reativar pessoa", Descricao: "inativa → ativa, com If-Match. 409 transicao_invalida se já está ativa.",
	}), func(ctx context.Context, in *entradaIDIfMatch) (*saida[Pessoa], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		atual, err := conferida(ctx, svc, a, in.ID, in.IfMatch)
		if err != nil {
			return nil, err
		}
		p, err := svc.ReativarPessoa(ctx, a, atual.ID)
		if err != nil {
			return nil, err
		}
		return comETag(paraPessoa(p), p.Versao), nil
	})
}
