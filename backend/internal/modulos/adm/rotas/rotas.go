// Package rotas é a camada HTTP do módulo ADM (seção 6.1 do documento 3):
// conta da empresa, unidades com turnos e feriados, pessoas, usuários,
// perfis, catálogo de permissões, trilha de auditoria e as rotas de convite.
//
// Fica fora de adm/internal porque os DTOs do caso de uso estão no pacote
// raiz adm, e um adm/internal/http que o pacote raiz registrasse formaria
// ciclo de importação. Importa só a API pública do módulo e a plataforma.
package rotas

import (
	"context"
	"net/http"

	"github.com/google/uuid"

	"tenaz/internal/modulos/adm"
	"tenaz/internal/plataforma/api"
	"tenaz/internal/plataforma/autenticacao"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/erros"
	"tenaz/internal/plataforma/textos"
)

// Convites é o que as rotas de convite usam do serviço de autenticação.
type Convites interface {
	LocalizarConvite(ctx context.Context, token string) (autenticacao.ConviteValido, error)
}

type rotas struct {
	svc      adm.Servico
	convites Convites
}

func (r *rotas) servico() (adm.Servico, error) {
	if r.svc == nil {
		return nil, erros.ServicoIndisponivel.ComChave("detalhe.adm_nao_configurado", nil)
	}
	return r.svc, nil
}

// contexto devolve o serviço e o ator da requisição.
func (r *rotas) contexto(ctx context.Context) (adm.Servico, autorizacao.Ator, error) {
	s, err := r.servico()
	if err != nil {
		return nil, autorizacao.Ator{}, err
	}
	a, err := api.AtorDe(ctx)
	return s, a, err
}

// idDoCaminho lê o UUID do caminho. Identificador malformado não pode existir:
// 404, como o inexistente.
func idDoCaminho(s string) (uuid.UUID, error) {
	id, err := uuid.Parse(s)
	if err != nil {
		return uuid.Nil, erros.NaoEncontrado
	}
	return id, nil
}

// ------------------------------------------------------------ entradas

type semEntrada struct{}

type entradaID struct {
	ID string `path:"id" doc:"Identificador (UUID)"`
}

type entradaIDIfMatch struct {
	ID      string `path:"id"`
	IfMatch string `header:"If-Match" doc:"ETag lido do recurso; obrigatório"`
}

type entradaCorpo[T any] struct {
	Body T
}

type entradaIDCorpo[T any] struct {
	ID   string `path:"id"`
	Body T
}

type entradaPatch[T any] struct {
	IfMatch string `header:"If-Match" doc:"ETag lido do recurso; obrigatório"`
	Body    T
	RawBody []byte
}

type entradaPatchID[T any] struct {
	ID      string `path:"id"`
	IfMatch string `header:"If-Match" doc:"ETag lido do recurso; obrigatório"`
	Body    T
	RawBody []byte
}

type entradaSub struct {
	ID    string `path:"id"`
	SubID string `path:"sub_id"`
}

type entradaSubIfMatch struct {
	ID      string `path:"id"`
	SubID   string `path:"sub_id"`
	IfMatch string `header:"If-Match" doc:"ETag lido do recurso; obrigatório"`
}

type entradaPatchSub[T any] struct {
	ID      string `path:"id"`
	SubID   string `path:"sub_id"`
	IfMatch string `header:"If-Match" doc:"ETag lido do recurso; obrigatório"`
	Body    T
	RawBody []byte
}

// ------------------------------------------------------------- saídas

type saida[T any] struct {
	ETag     *string `header:"ETag"`
	Location *string `header:"Location"`
	Body     T
}

// saidaIdioma é a resposta que leva texto traduzido pelo catálogo do servidor:
// Content-Language e Vary: Accept-Language (ADR-013).
type saidaIdioma[T any] struct {
	ContentLanguage string `header:"Content-Language"`
	Vary            string `header:"Vary"`
	Body            T
}

func comIdioma[T any](ctx context.Context, corpo T) *saidaIdioma[T] {
	return &saidaIdioma[T]{ContentLanguage: string(textos.DoContexto(ctx)), Vary: textos.CabecalhoVary, Body: corpo}
}

type saidaLista[T any] struct {
	Body api.Lista[T]
}

type saidaVazia struct{}

func comETag[T any](corpo T, versao int64) *saida[T] {
	return &saida[T]{ETag: api.ETag(versao), Body: corpo}
}

func criado[T any](corpo T, versao int64, local string) *saida[T] {
	return &saida[T]{ETag: api.ETag(versao), Location: &local, Body: corpo}
}

// Registrar registra as rotas do ADM no servidor. svc é nil só na geração
// do contrato.
func Registrar(s *api.Servidor, svc adm.Servico) {
	r := &rotas{svc: svc}
	if a := s.Autenticacao(); a != nil {
		r.convites = a
	}
	r.registrarConta(s)
	r.registrarUnidades(s)
	r.registrarPessoas(s)
	r.registrarUsuarios(s)
	r.registrarCatalogos(s)
}

func rotaADM(ro api.Rota) api.Rota {
	ro.Tag, ro.Modulo = api.TagADM, "ADM"
	return ro
}

// ------------------------------------------------------------- conta

type entradaConviteConsultar struct {
	Body struct {
		Token string `json:"token" minLength:"43" maxLength:"43" doc:"Token do link do convite, lido do fragmento #token="`
	}
}

// Convite é o convite localizado pelo link.
type Convite struct {
	Email     string `json:"email" doc:"E-mail convidado"`
	PedeSenha bool   `json:"pede_senha" doc:"Verdadeiro: o e-mail ainda não tem login e a tela pede a senha nova; falso: entra com a senha que já tem"`
}

type entradaConviteAceitar struct {
	Body struct {
		Token string `json:"token" minLength:"43" maxLength:"43"`
		Senha string `json:"senha,omitempty" maxLength:"1024" doc:"Obrigatória quando pede_senha; ignorada quando o e-mail já tem login"`
		// Idioma em que a pessoa aceitou (o da tela, do navegador ou do aparelho).
		Idioma *string `json:"idioma,omitempty" enum:"pt-BR,es-419" doc:"Idioma em que a pessoa aceitou; passa a ser o dela. Ausente: o do convite (padrão da empresa)"`
	}
}

// ConviteAceito é a resposta do aceite.
type ConviteAceito struct {
	IdentidadeNova bool `json:"identidade_nova" doc:"Falso: o e-mail já entrava em outra empresa e usa a senha que tem"`
}

func (r *rotas) registrarConta(s *api.Servidor) {
	api.Registrar(s, rotaADM(api.Rota{
		ID: "criarEmpresa", Metodo: http.MethodPost, Caminho: "/v1/empresas", Status: http.StatusCreated, Publica: true,
		Resumo: "Criar a conta da empresa (cadastro pelo site)",
		Descricao: "Cria a empresa com a unidade inicial, os 12 perfis padrão com a matriz de RF-ADM-005 e o responsável " +
			"como administrador convidado, e envia o convite, que confirma o e-mail antes do primeiro acesso " +
			"(RF-ADM-001; RF-SRV-001 com srv_habilitado). O pais define a moeda e o idioma padrão. 409 conflito em identificação já cadastrada.",
	}), func(ctx context.Context, in *entradaCorpo[NovaEmpresa]) (*saida[ContaCriada], error) {
		svc, err := r.servico()
		if err != nil {
			return nil, err
		}
		b := in.Body
		// O responsável fica com o idioma da tela de cadastro (Accept-Language).
		idiomaResp, _ := textos.DeAcceptLanguage(textos.PreferenciasDe(ctx).AcceptLanguage)
		c, err := svc.CriarConta(ctx, adm.EntradaCriarConta{
			Pais: b.Pais, IdentificacaoTipo: b.IdentificacaoTipo, Identificacao: b.Identificacao,
			IdiomaPadrao: b.IdiomaPadrao, IdiomaResponsavel: string(idiomaResp), RazaoSocial: b.RazaoSocial, NomeFantasia: b.NomeFantasia, FusoHorario: b.FusoHorario, Plano: b.Plano,
			SrvHabilitado: b.SrvHabilitado,
			Responsavel:   adm.EntradaResponsavel{Nome: b.Responsavel.Nome, Email: b.Responsavel.Email, Telefone: b.Responsavel.Telefone},
			Unidade:       deUnidade(b.Unidade),
		})
		if err != nil {
			return nil, err
		}
		return &saida[ContaCriada]{Body: ContaCriada{ConviteEnviado: c.ConviteEnviado}}, nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "consultarConvite", Metodo: http.MethodPost, Caminho: "/v1/auth/convites/consultar", Publica: true, Tag: api.TagAutenticacao,
		Resumo: "Consultar o convite do link",
		Descricao: "Confere o link sem gastá-lo, para a tela de aceite saber se pede senha. O token vai no corpo, " +
			"nunca na URL. 410 link_email_invalido.",
	}), func(ctx context.Context, in *entradaConviteConsultar) (*saida[Convite], error) {
		if r.convites == nil {
			return nil, erros.ServicoIndisponivel
		}
		c, err := r.convites.LocalizarConvite(ctx, in.Body.Token)
		if err != nil {
			return nil, err
		}
		return &saida[Convite]{Body: Convite{Email: c.Email, PedeSenha: !c.IdentidadeExiste}}, nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "aceitarConvite", Metodo: http.MethodPost, Caminho: "/v1/auth/convites/aceitar", Publica: true, Tag: api.TagAutenticacao,
		Resumo: "Aceitar o convite",
		Descricao: "Liga o login ao usuário convidado e o ativa; o link é gasto. Depois, a pessoa entra por " +
			"POST /v1/auth/login. 410 link_email_invalido; 422 senha_fraca; 409 conflito se o e-mail já tem usuário na empresa.",
	}), func(ctx context.Context, in *entradaConviteAceitar) (*saida[ConviteAceito], error) {
		svc, err := r.servico()
		if err != nil {
			return nil, err
		}
		c, err := svc.AceitarConvite(ctx, in.Body.Token, in.Body.Senha, in.Body.Idioma)
		if err != nil {
			return nil, err
		}
		return &saida[ConviteAceito]{Body: ConviteAceito{IdentidadeNova: c.IdentidadeNova}}, nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "obterEmpresa", Metodo: http.MethodGet, Caminho: "/v1/empresa", Permissoes: []string{autorizacao.AdmLer},
		Resumo: "Dados da empresa", Descricao: "A empresa da credencial. Traz ETag.",
	}), func(ctx context.Context, _ *semEntrada) (*saida[Empresa], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		e, err := svc.ObterEmpresa(ctx, a)
		if err != nil {
			return nil, err
		}
		return comETag(paraEmpresa(e), e.Versao), nil
	})

	api.Registrar(s, rotaADM(api.Rota{
		ID: "alterarEmpresa", Metodo: http.MethodPatch, Caminho: "/v1/empresa", Permissoes: []string{autorizacao.AdmGerenciarEmpresa},
		Resumo: "Alterar dados cadastrais",
		Descricao: "JSON Merge Patch com If-Match. Plano e módulo SRV mudam pela cobrança (RF-ADM-009), não aqui; país e " +
			"moeda não mudam pela API: 422 campo_somente_leitura. A identificação é a do tipo do país da empresa e " +
			"idioma_padrao é o que o usuário convidado recebe. mfa_obrigatorio liga o segundo fator obrigatório: quem " +
			"não o tem passa a entrar com sessão restrita à ativação (RN-ADM-027).",
	}), func(ctx context.Context, in *entradaPatch[PatchEmpresa]) (*saida[Empresa], error) {
		svc, a, err := r.contexto(ctx)
		if err != nil {
			return nil, err
		}
		esperada, err := api.VersaoEsperada(in.IfMatch)
		if err != nil {
			return nil, err
		}
		atual, err := svc.ObterEmpresa(ctx, a)
		if err != nil {
			return nil, err
		}
		if err := api.ConferirVersao(esperada, atual.Versao, func() any { return paraEmpresa(atual) }); err != nil {
			return nil, err
		}
		base := camposEmpresa{RazaoSocial: atual.RazaoSocial, NomeFantasia: atual.NomeFantasia, FusoHorario: atual.FusoHorario,
			MFAObrigatorio: atual.MFAObrigatorio}
		base.IdiomaPadrao = atual.IdiomaPadrao
		if atual.IdentificacaoTipo != nil {
			base.IdentificacaoTipo = *atual.IdentificacaoTipo
		}
		if atual.Identificacao != nil {
			base.Identificacao = *atual.Identificacao
		}
		var novo camposEmpresa
		if _, err := api.MesclarPatch(base, in.RawBody, &novo); err != nil {
			return nil, err
		}
		e, err := svc.EditarEmpresa(ctx, a, adm.EntradaEmpresa{
			RazaoSocial: novo.RazaoSocial, NomeFantasia: novo.NomeFantasia, IdentificacaoTipo: novo.IdentificacaoTipo,
			Identificacao: novo.Identificacao, IdiomaPadrao: novo.IdiomaPadrao, FusoHorario: novo.FusoHorario,
			MFAObrigatorio: novo.MFAObrigatorio, VersaoEsperada: &esperada,
		})
		if err != nil {
			return nil, desatualizado(err, func() (any, int64, error) {
				x, e := svc.ObterEmpresa(ctx, a)
				return paraEmpresa(x), x.Versao, e
			})
		}
		return comETag(paraEmpresa(e), e.Versao), nil
	})
}

// desatualizado acrescenta o recurso atual ao 412 que o caso de uso devolveu
// (a versão mudou entre a leitura e a gravação).
func desatualizado(err error, ler func() (any, int64, error)) error {
	if !api.EhDesatualizado(err) {
		return err
	}
	atual, versao, errLer := ler()
	if errLer != nil {
		return err
	}
	return api.Desatualizado(versao, atual)
}
