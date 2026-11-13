package adm

import (
	"context"
	"errors"
	"slices"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"tenaz/internal/modulos/adm/internal/dominio"
	"tenaz/internal/modulos/adm/internal/repo"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/db/consultas"
	"tenaz/internal/plataforma/erros"
	"tenaz/internal/plataforma/paginacao"
)

// Usuários, perfis e permissões avulsas: RF-ADM-004 e RF-ADM-005.

func escopoDoUsuario(fornecedor, locatario *uuid.UUID) string {
	switch {
	case fornecedor != nil:
		return db.EscopoPrestador
	case locatario != nil:
		return db.EscopoLocatario
	}
	return db.EscopoInterno
}

// lerUsuario lê o usuário com os perfis e as permissões avulsas. Em escopo
// externo, o RLS esconde usuario_perfil e usuario_permissao; ObterEu completa
// os perfis pelo Ator.
func lerUsuario(ctx context.Context, empresaID, id uuid.UUID) (Usuario, error) {
	q, err := repo.Q(ctx)
	if err != nil {
		return Usuario{}, err
	}
	r, err := q.AdmBuscarUsuario(ctx, consultas.AdmBuscarUsuarioParams{EmpresaID: empresaID, ID: id})
	if err != nil {
		return Usuario{}, naoEncontrado(err)
	}
	u := Usuario{
		ID: r.ID, IdentidadeLigada: r.IdentidadeID != nil, PessoaID: r.PessoaID, FornecedorID: r.FornecedorID,
		LocatarioID: r.LocatarioID, EscopoAcesso: escopoDoUsuario(r.FornecedorID, r.LocatarioID), Nome: r.Nome,
		Email: r.Email, Telefone: r.Telefone, CustoHoraCentavos: centavos(r.CustoHora),
		UnidadePadraoID: r.UnidadePadraoID, Situacao: r.Situacao, Versao: r.Versao, CriadoEm: r.CriadoEm,
		AtualizadoEm: r.AtualizadoEm, Perfis: []PerfilAtribuido{}, PermissoesAvulsas: []string{},
	}
	perfis, err := q.AdmPerfisDoUsuario(ctx, consultas.AdmPerfisDoUsuarioParams{EmpresaID: empresaID, UsuarioID: id})
	if err != nil {
		return Usuario{}, err
	}
	for _, p := range perfis {
		u.Perfis = append(u.Perfis, PerfilAtribuido{PerfilID: p.PerfilID, Codigo: p.Codigo, Nome: p.Nome,
			CodigoOrigem: p.CodigoOrigem, UnidadeID: p.UnidadeID})
	}
	avulsas, err := q.AdmPermissoesAvulsas(ctx, consultas.AdmPermissoesAvulsasParams{EmpresaID: empresaID, UsuarioID: id})
	if err != nil {
		return Usuario{}, err
	}
	for _, p := range avulsas {
		u.PermissoesAvulsas = append(u.PermissoesAvulsas, p.Permissao)
	}
	return u, nil
}

// perfisPorCodigo busca os perfis da empresa pelos códigos e devolve as
// permissões de cada um, pelo código.
func perfisPorCodigo(ctx context.Context, q *consultas.Queries, empresaID uuid.UUID, codigos []string) (map[string]consultas.AdmBuscarPerfisPorCodigosRow, error) {
	linhas, err := q.AdmBuscarPerfisPorCodigos(ctx, consultas.AdmBuscarPerfisPorCodigosParams{EmpresaID: empresaID, Codigos: codigos})
	if err != nil {
		return nil, err
	}
	m := make(map[string]consultas.AdmBuscarPerfisPorCodigosRow, len(linhas))
	for _, l := range linhas {
		m[l.Codigo] = l
	}
	for i, c := range codigos {
		if _, ok := m[c]; !ok {
			return nil, erros.EntradaInvalida.ComCampo("perfis", erros.CampoNaoEncontrado, "Perfil não encontrado: "+codigos[i]+".")
		}
	}
	return m, nil
}

// tipoDosPerfis devolve o código de origem de cada perfil (o tipo).
func tipoDosPerfis(m map[string]consultas.AdmBuscarPerfisPorCodigosRow, codigos []string) []string {
	tipos := make([]string, 0, len(codigos))
	for _, c := range codigos {
		p := m[c]
		t := p.Codigo
		if p.PerfilOrigem != nil {
			t = *p.PerfilOrigem
		}
		tipos = append(tipos, t)
	}
	return tipos
}

// conferirPessoa aplica a regra de RF-ADM-015: pessoa com fornecedor só se
// liga a usuário prestador daquele fornecedor. A unicidade (uma pessoa, um
// usuário) fica com o índice uq_usuario_pessoa.
func conferirPessoa(ctx context.Context, q *consultas.Queries, empresaID, pessoaID uuid.UUID, fornecedorUsuario *uuid.UUID) error {
	p, err := q.AdmTravarPessoa(ctx, consultas.AdmTravarPessoaParams{EmpresaID: empresaID, ID: pessoaID})
	if errors.Is(err, pgx.ErrNoRows) {
		return erros.EntradaInvalida.ComCampo("pessoa_id", erros.CampoNaoEncontrado, "Pessoa não encontrada.")
	}
	if err != nil {
		return err
	}
	if p.FornecedorID != nil && !mesmoID(p.FornecedorID, fornecedorUsuario) {
		return erros.EntradaInvalida.ComCampo("pessoa_id", erros.CampoValorInvalido,
			"Pessoa de fornecedor só se liga a usuário prestador do mesmo fornecedor.")
	}
	return nil
}

func validarCusto(c *int64) error {
	if c != nil && (*c < 0 || *c > 9_999_999_999) {
		return erros.EntradaInvalida.ComCampo("custo_hora_centavos", erros.CampoForaDaFaixa, "Custo por hora de 0 a 99.999.999,99.")
	}
	return nil
}

func (s *servico) CriarUsuario(ctx context.Context, a autorizacao.Ator, e EntradaUsuario) (UsuarioCriado, error) {
	if err := apenasInterno(a, autorizacao.AdmGerenciarUsuarios); err != nil {
		return UsuarioCriado{}, err
	}
	nome, ee := dominio.Texto("nome", e.Nome, 200)
	if ee != nil {
		return UsuarioCriado{}, ee
	}
	email, ee := dominio.Email("email", e.Email)
	if ee != nil {
		return UsuarioCriado{}, ee
	}
	tel, ee := dominio.TextoOpcional("telefone", e.Telefone, 30)
	if ee != nil {
		return UsuarioCriado{}, ee
	}
	if err := validarCusto(e.CustoHoraCentavos); err != nil {
		return UsuarioCriado{}, err
	}
	if e.PessoaID != nil && e.NovaPessoa != nil {
		return UsuarioCriado{}, erros.EntradaInvalida.ComCampo("nova_pessoa", erros.CampoValorInvalido,
			"Informe pessoa_id ou nova_pessoa, não os dois.")
	}
	var novaPessoa *pessoaValida
	if e.NovaPessoa != nil {
		pv, err := validarPessoa("nova_pessoa.", *e.NovaPessoa)
		if err != nil {
			return UsuarioCriado{}, err
		}
		novaPessoa = &pv
	}

	var (
		criado      Usuario
		nomeEmpresa string
	)
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		perfis, err := perfisPorCodigo(ctx, q, a.EmpresaID, e.Perfis)
		if err != nil {
			return err
		}
		if ee := dominio.ValidarPerfisDoUsuario(tipoDosPerfis(perfis, e.Perfis),
			dominio.Vinculo{Fornecedor: e.FornecedorID != nil, Locatario: e.LocatarioID != nil}); ee != nil {
			return ee
		}
		id := novoID()
		novos := map[string][]string{}
		for c, p := range perfis {
			novos[c] = p.Permissoes
		}
		if ee := dominio.ConferirAtribuicao(a, id, false, novos); ee != nil {
			return ee
		}
		pessoaID := e.PessoaID
		if novaPessoa != nil {
			pessoaID = ptr(novoID())
			if err := criarPessoaTx(ctx, q, a.EmpresaID, *pessoaID, *novaPessoa); err != nil {
				return err
			}
		}
		if pessoaID != nil {
			if err := conferirPessoa(ctx, q, a.EmpresaID, *pessoaID, e.FornecedorID); err != nil {
				return err
			}
		}
		err = q.AdmCriarUsuario(ctx, consultas.AdmCriarUsuarioParams{
			EmpresaID: a.EmpresaID, ID: id, PessoaID: pessoaID, FornecedorID: e.FornecedorID, LocatarioID: e.LocatarioID,
			Nome: nome, Email: &email, Telefone: tel, CustoHoraCentavos: e.CustoHoraCentavos, UnidadePadraoID: e.UnidadePadraoID,
		})
		if err != nil {
			return repo.Traduzir(err)
		}
		for _, c := range e.Perfis {
			if err := q.AdmColocarPerfilUsuario(ctx, consultas.AdmColocarPerfilUsuarioParams{
				EmpresaID: a.EmpresaID, ID: novoID(), UsuarioID: id, PerfilID: perfis[c].ID,
			}); err != nil {
				return repo.Traduzir(err)
			}
		}
		if nomeEmpresa, err = autorizacao.NomeEmpresa(ctx, a.EmpresaID); err != nil {
			return err
		}
		criado, err = lerUsuario(ctx, a.EmpresaID, id)
		return err
	})
	if err != nil {
		return UsuarioCriado{}, err
	}
	enviado := s.convidar(ctx, a.EmpresaID, criado.ID, email, nomeEmpresa, nome) == nil
	return UsuarioCriado{Usuario: criado, ConviteEnviado: enviado}, nil
}

func (s *servico) ObterUsuario(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (Usuario, error) {
	if err := apenasInterno(a, autorizacao.AdmLerUsuarios); err != nil {
		return Usuario{}, err
	}
	var u Usuario
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		var err error
		u, err = lerUsuario(ctx, a.EmpresaID, id)
		return err
	})
	return u, err
}

// ObterEu devolve o próprio usuário, inclusive para o externo.
func (s *servico) ObterEu(ctx context.Context, a autorizacao.Ator) (Usuario, error) {
	var u Usuario
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		var err error
		u, err = lerUsuario(ctx, a.EmpresaID, a.UsuarioID)
		return err
	})
	if err != nil {
		return Usuario{}, err
	}
	if a.Externo() {
		// O externo não lê usuario_perfil (RLS): os perfis vêm do Ator,
		// resolvido nesta requisição.
		u.Perfis = u.Perfis[:0]
		for _, p := range a.Perfis {
			u.Perfis = append(u.Perfis, PerfilAtribuido{Codigo: p, CodigoOrigem: p})
		}
	}
	return u, nil
}

// EditarEu altera só nome e telefone do próprio usuário (ISO-E01-01).
func (s *servico) EditarEu(ctx context.Context, a autorizacao.Ator, e EntradaEu) (Usuario, error) {
	nome, ee := dominio.Texto("nome", e.Nome, 200)
	if ee != nil {
		return Usuario{}, ee
	}
	tel, ee := dominio.TextoOpcional("telefone", e.Telefone, 30)
	if ee != nil {
		return Usuario{}, ee
	}
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		return repo.Traduzir(q.AdmAtualizarEu(ctx, consultas.AdmAtualizarEuParams{
			EmpresaID: a.EmpresaID, ID: a.UsuarioID, Nome: nome, Telefone: tel,
		}))
	})
	if err != nil {
		return Usuario{}, err
	}
	return s.ObterEu(ctx, a)
}

func (s *servico) EditarUsuario(ctx context.Context, a autorizacao.Ator, id uuid.UUID, e EntradaEditarUsuario) (Usuario, error) {
	if err := apenasInterno(a, autorizacao.AdmGerenciarUsuarios); err != nil {
		return Usuario{}, err
	}
	nome, ee := dominio.Texto("nome", e.Nome, 200)
	if ee != nil {
		return Usuario{}, ee
	}
	tel, ee := dominio.TextoOpcional("telefone", e.Telefone, 30)
	if ee != nil {
		return Usuario{}, ee
	}
	if err := validarCusto(e.CustoHoraCentavos); err != nil {
		return Usuario{}, err
	}
	var u Usuario
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		atual, err := q.AdmTravarUsuario(ctx, consultas.AdmTravarUsuarioParams{EmpresaID: a.EmpresaID, ID: id})
		if err != nil {
			return naoEncontrado(err)
		}
		if err := conferirVersao(e.VersaoEsperada, atual.Versao); err != nil {
			return err
		}
		if err := soAdministradorAlteraAdministrador(ctx, q, a, id); err != nil {
			return err
		}
		if e.PessoaID != nil && !mesmoID(e.PessoaID, atual.PessoaID) {
			if err := conferirPessoa(ctx, q, a.EmpresaID, *e.PessoaID, atual.FornecedorID); err != nil {
				return err
			}
		}
		err = q.AdmAtualizarUsuario(ctx, consultas.AdmAtualizarUsuarioParams{
			EmpresaID: a.EmpresaID, ID: id, Nome: nome, Telefone: tel, CustoHoraCentavos: e.CustoHoraCentavos,
			UnidadePadraoID: e.UnidadePadraoID, PessoaID: e.PessoaID,
		})
		if err != nil {
			return repo.Traduzir(err)
		}
		u, err = lerUsuario(ctx, a.EmpresaID, id)
		return err
	})
	return u, err
}

// ehAdministrador informa se o usuário tem algum perfil do tipo administrador.
func ehAdministrador(ctx context.Context, q *consultas.Queries, empresaID, usuarioID uuid.UUID) (bool, error) {
	perfis, err := q.AdmPerfisDoUsuario(ctx, consultas.AdmPerfisDoUsuarioParams{EmpresaID: empresaID, UsuarioID: usuarioID})
	if err != nil {
		return false, err
	}
	for _, p := range perfis {
		if p.CodigoOrigem == autorizacao.PerfilAdministrador {
			return true, nil
		}
	}
	return false, nil
}

func soAdministradorAlteraAdministrador(ctx context.Context, q *consultas.Queries, a autorizacao.Ator, alvo uuid.UUID) error {
	if a.TemPerfil(autorizacao.PerfilAdministrador) {
		return nil
	}
	adm, err := ehAdministrador(ctx, q, a.EmpresaID, alvo)
	if err != nil {
		return err
	}
	if adm {
		return erros.SemPermissao.Com("só administrador altera administrador")
	}
	return nil
}

// garantirOutroAdministrador recusa a operação que deixaria a empresa sem
// administrador ativo. A trava no perfil administrador serializa duas
// operações simultâneas (cada uma veria a outra ainda ativa).
func garantirOutroAdministrador(ctx context.Context, q *consultas.Queries, empresaID, alvo uuid.UUID) error {
	if _, err := q.AdmTravarPerfilAdministrador(ctx, empresaID); err != nil {
		return err
	}
	n, err := q.AdmContarAdministradoresAtivos(ctx, consultas.AdmContarAdministradoresAtivosParams{
		EmpresaID: empresaID, ExcetoUsuarioID: &alvo,
	})
	if err != nil {
		return err
	}
	if n == 0 {
		return ErrUltimoAdministrador.Com("a empresa precisa de pelo menos um administrador ativo")
	}
	return nil
}

// DefinirPerfis troca o conjunto de perfis do usuário (PATCH com `perfis`).
func (s *servico) DefinirPerfis(ctx context.Context, a autorizacao.Ator, id uuid.UUID, codigos []string) (Usuario, error) {
	if err := apenasInterno(a, autorizacao.AdmGerenciarUsuarios); err != nil {
		return Usuario{}, err
	}
	var u Usuario
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		alvo, err := q.AdmTravarUsuario(ctx, consultas.AdmTravarUsuarioParams{EmpresaID: a.EmpresaID, ID: id})
		if err != nil {
			return naoEncontrado(err)
		}
		perfis, err := perfisPorCodigo(ctx, q, a.EmpresaID, codigos)
		if err != nil {
			return err
		}
		if ee := dominio.ValidarPerfisDoUsuario(tipoDosPerfis(perfis, codigos),
			dominio.Vinculo{Fornecedor: alvo.FornecedorID != nil, Locatario: alvo.LocatarioID != nil}); ee != nil {
			return ee
		}
		atuais, err := q.AdmPerfisDoUsuario(ctx, consultas.AdmPerfisDoUsuarioParams{EmpresaID: a.EmpresaID, UsuarioID: id})
		if err != nil {
			return err
		}
		eraAdm, ficaAdm := false, false
		atuaisPorCodigo := map[string]consultas.AdmPerfisDoUsuarioRow{}
		for _, p := range atuais {
			atuaisPorCodigo[p.Codigo] = p
			eraAdm = eraAdm || p.CodigoOrigem == autorizacao.PerfilAdministrador
		}
		novos := map[string][]string{}
		for _, t := range tipoDosPerfis(perfis, codigos) {
			ficaAdm = ficaAdm || t == autorizacao.PerfilAdministrador
		}
		for _, c := range codigos {
			if _, ja := atuaisPorCodigo[c]; !ja {
				novos[c] = perfis[c].Permissoes
			}
		}
		// Anti-escalada sobre o que entra; a regra do próprio usuário e a do
		// administrador valem também para o que sai.
		if ee := dominio.ConferirAtribuicao(a, id, eraAdm, novos); ee != nil {
			return ee
		}
		if eraAdm && !ficaAdm && alvo.Situacao == dominio.UsuarioAtivo {
			if err := garantirOutroAdministrador(ctx, q, a.EmpresaID, id); err != nil {
				return err
			}
		}
		for c, p := range atuaisPorCodigo {
			if !slices.Contains(codigos, c) {
				if err := q.AdmRetirarPerfilUsuario(ctx, consultas.AdmRetirarPerfilUsuarioParams{EmpresaID: a.EmpresaID, ID: p.ID}); err != nil {
					return repo.Traduzir(err)
				}
			}
		}
		for c := range novos {
			if err := q.AdmColocarPerfilUsuario(ctx, consultas.AdmColocarPerfilUsuarioParams{
				EmpresaID: a.EmpresaID, ID: novoID(), UsuarioID: id, PerfilID: perfis[c].ID,
			}); err != nil {
				return repo.Traduzir(err)
			}
		}
		u, err = lerUsuario(ctx, a.EmpresaID, id)
		return err
	})
	return u, err
}

// InativarUsuario: RF-ADM-004. Não inativa o último administrador. As
// sessões e os refresh tokens são revogados pelo trigger de acesso_empresa,
// na mesma transação (ADR-011); app.conferir_escopo() recusa o usuário
// inativo na requisição seguinte, como segunda barreira. A limpeza do cache do
// app (RF-MOB-010) entra com a sincronização.
func (s *servico) InativarUsuario(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (Usuario, error) {
	if err := apenasInterno(a, autorizacao.AdmGerenciarUsuarios); err != nil {
		return Usuario{}, err
	}
	var u Usuario
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		alvo, err := q.AdmTravarUsuario(ctx, consultas.AdmTravarUsuarioParams{EmpresaID: a.EmpresaID, ID: id})
		if err != nil {
			return naoEncontrado(err)
		}
		if alvo.Situacao == dominio.UsuarioInativo {
			return erros.TransicaoInvalida.Com("o usuário já está inativo").ComDado("estado_atual", alvo.Situacao)
		}
		if err := soAdministradorAlteraAdministrador(ctx, q, a, id); err != nil {
			return err
		}
		adm, err := ehAdministrador(ctx, q, a.EmpresaID, id)
		if err != nil {
			return err
		}
		if adm && alvo.Situacao == dominio.UsuarioAtivo {
			if err := garantirOutroAdministrador(ctx, q, a.EmpresaID, id); err != nil {
				return err
			}
		}
		if err := q.AdmDefinirSituacaoUsuario(ctx, consultas.AdmDefinirSituacaoUsuarioParams{
			EmpresaID: a.EmpresaID, ID: id, Situacao: dominio.UsuarioInativo,
		}); err != nil {
			return repo.Traduzir(err)
		}
		u, err = lerUsuario(ctx, a.EmpresaID, id)
		return err
	})
	return u, err
}

// ReativarUsuario: quem já tinha aceitado o convite volta a ativo; quem não
// tinha volta a convidado e recebe um convite novo (proposta).
func (s *servico) ReativarUsuario(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (UsuarioCriado, error) {
	if err := apenasInterno(a, autorizacao.AdmGerenciarUsuarios); err != nil {
		return UsuarioCriado{}, err
	}
	var (
		u           Usuario
		nomeEmpresa string
		reconvidar  bool
	)
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		alvo, err := q.AdmTravarUsuario(ctx, consultas.AdmTravarUsuarioParams{EmpresaID: a.EmpresaID, ID: id})
		if err != nil {
			return naoEncontrado(err)
		}
		if alvo.Situacao != dominio.UsuarioInativo {
			return erros.TransicaoInvalida.Com("só usuário inativo é reativado").ComDado("estado_atual", alvo.Situacao)
		}
		if err := soAdministradorAlteraAdministrador(ctx, q, a, id); err != nil {
			return err
		}
		nova := dominio.UsuarioAtivo
		if alvo.IdentidadeID == nil {
			nova, reconvidar = dominio.UsuarioConvidado, true
		}
		if err := q.AdmDefinirSituacaoUsuario(ctx, consultas.AdmDefinirSituacaoUsuarioParams{
			EmpresaID: a.EmpresaID, ID: id, Situacao: nova,
		}); err != nil {
			return repo.Traduzir(err)
		}
		if nomeEmpresa, err = autorizacao.NomeEmpresa(ctx, a.EmpresaID); err != nil {
			return err
		}
		u, err = lerUsuario(ctx, a.EmpresaID, id)
		return err
	})
	if err != nil {
		return UsuarioCriado{}, err
	}
	res := UsuarioCriado{Usuario: u}
	if reconvidar && u.Email != nil {
		res.ConviteEnviado = s.convidar(ctx, a.EmpresaID, id, *u.Email, nomeEmpresa, u.Nome) == nil
	}
	return res, nil
}

// ReenviarConvite emite um convite novo (o anterior deixa de valer) para o
// usuário convidado. Falha no envio: servico_indisponivel.
func (s *servico) ReenviarConvite(ctx context.Context, a autorizacao.Ator, id uuid.UUID) error {
	if err := apenasInterno(a, autorizacao.AdmGerenciarUsuarios); err != nil {
		return err
	}
	var (
		u           Usuario
		nomeEmpresa string
	)
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		var err error
		if u, err = lerUsuario(ctx, a.EmpresaID, id); err != nil {
			return err
		}
		if u.Situacao != dominio.UsuarioConvidado || u.Email == nil {
			return erros.TransicaoInvalida.Com("só usuário convidado recebe convite").ComDado("estado_atual", u.Situacao)
		}
		nomeEmpresa, err = autorizacao.NomeEmpresa(ctx, a.EmpresaID)
		return err
	})
	if err != nil {
		return err
	}
	return s.convidar(ctx, a.EmpresaID, id, *u.Email, nomeEmpresa, u.Nome)
}

// EncerrarSessoes revoga as sessões e os dispositivos do usuário (POST
// /v1/usuarios/{id}/encerrar-sessoes). Confere no pool da aplicação que o
// usuário é desta empresa; a revogação e o evento são da autenticação.
func (s *servico) EncerrarSessoes(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (int32, error) {
	if err := apenasInterno(a, autorizacao.AdmGerenciarUsuarios); err != nil {
		return 0, err
	}
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		if _, err := q.AdmTravarUsuario(ctx, consultas.AdmTravarUsuarioParams{EmpresaID: a.EmpresaID, ID: id}); err != nil {
			return naoEncontrado(err)
		}
		return soAdministradorAlteraAdministrador(ctx, q, a, id)
	})
	if err != nil {
		return 0, err
	}
	n, err := s.d.Autenticacao.EncerrarSessoesUsuario(ctx, a, id)
	return int32(n), err
}

// ConcederPermissao concede uma permissão avulsa (RF-ADM-005), como
// pessoa.saude.ler. A concessão vai para a trilha pelo trg_auditoria de
// usuario_permissao.
func (s *servico) ConcederPermissao(ctx context.Context, a autorizacao.Ator, id uuid.UUID, permissao string) (Usuario, error) {
	if err := apenasInterno(a, autorizacao.AdmGerenciarUsuarios); err != nil {
		return Usuario{}, err
	}
	var u Usuario
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		alvo, err := q.AdmTravarUsuario(ctx, consultas.AdmTravarUsuarioParams{EmpresaID: a.EmpresaID, ID: id})
		if err != nil {
			return naoEncontrado(err)
		}
		// O tipo do usuário externo é o do vínculo: fornecedor é prestador,
		// locatário é locatario (RN-ADM-022), inclusive com perfil derivado.
		perfilExterno := ""
		switch {
		case alvo.FornecedorID != nil:
			perfilExterno = autorizacao.PerfilPrestador
		case alvo.LocatarioID != nil:
			perfilExterno = autorizacao.PerfilLocatario
		}
		if ee := dominio.ConferirConcessao(a, id, permissao, perfilExterno); ee != nil {
			return ee
		}
		if err := soAdministradorAlteraAdministrador(ctx, q, a, id); err != nil {
			return err
		}
		if err := q.AdmConcederPermissao(ctx, consultas.AdmConcederPermissaoParams{
			EmpresaID: a.EmpresaID, ID: novoID(), UsuarioID: id, Permissao: permissao,
		}); err != nil {
			return repo.Traduzir(err)
		}
		u, err = lerUsuario(ctx, a.EmpresaID, id)
		return err
	})
	return u, err
}

// RetirarPermissao retira a permissão avulsa, com as mesmas regras da
// concessão (ninguém mexe nas próprias; só quem pode conceder retira).
func (s *servico) RetirarPermissao(ctx context.Context, a autorizacao.Ator, id uuid.UUID, permissao string) (Usuario, error) {
	if err := apenasInterno(a, autorizacao.AdmGerenciarUsuarios); err != nil {
		return Usuario{}, err
	}
	if id == a.UsuarioID {
		return Usuario{}, erros.SemPermissao.Com("ninguém altera as próprias permissões")
	}
	if autorizacao.IsencaoAdministrador(permissao) && !a.TemPerfil(autorizacao.PerfilAdministrador) {
		return Usuario{}, erros.SemPermissao.Com("só o administrador retira %s", permissao)
	}
	if !a.PodeConceder(permissao) {
		return Usuario{}, erros.SemPermissao.Com("você não tem a permissão %s", permissao)
	}
	var u Usuario
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		if _, err := q.AdmTravarUsuario(ctx, consultas.AdmTravarUsuarioParams{EmpresaID: a.EmpresaID, ID: id}); err != nil {
			return naoEncontrado(err)
		}
		n, err := q.AdmRetirarPermissao(ctx, consultas.AdmRetirarPermissaoParams{EmpresaID: a.EmpresaID, UsuarioID: id, Permissao: permissao})
		if err != nil {
			return err
		}
		if n == 0 {
			return erros.NaoEncontrado.Com("o usuário não tem a permissão avulsa %s", permissao)
		}
		u, err = lerUsuario(ctx, a.EmpresaID, id)
		return err
	})
	return u, err
}

func (s *servico) ListarUsuarios(ctx context.Context, a autorizacao.Ator, f FiltroUsuarios) (paginacao.Pagina[ResumoUsuario], error) {
	var pag paginacao.Pagina[ResumoUsuario]
	if err := apenasInterno(a, autorizacao.AdmLerUsuarios); err != nil {
		return pag, err
	}
	busca := padraoBusca(f.Busca)
	consulta := paginacao.Assinatura("usuarios",
		textoOuVazio(busca, func(s string) string { return s }),
		textoOuVazio(f.Situacao, func(s string) string { return s }),
		textoOuVazio(f.Perfil, func(s string) string { return s }))
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
		linhas, err := q.AdmListarUsuarios(ctx, consultas.AdmListarUsuariosParams{
			EmpresaID: a.EmpresaID, Busca: busca, Situacao: f.Situacao, Perfil: f.Perfil,
			CursorCriadoEm: cEm, CursorID: cID, Limite: limite + 1,
		})
		if err != nil {
			return err
		}
		res := make([]ResumoUsuario, 0, len(linhas))
		for _, r := range linhas {
			res = append(res, ResumoUsuario{
				ID: r.ID, Nome: r.Nome, Email: r.Email, Telefone: r.Telefone, PessoaID: r.PessoaID,
				FornecedorID: r.FornecedorID, LocatarioID: r.LocatarioID, Situacao: r.Situacao, Perfis: r.Perfis,
				Versao: r.Versao, CriadoEm: r.CriadoEm,
			})
		}
		pag = paginacao.Montar(res, limite, consulta, agora, func(u ResumoUsuario) paginacao.Posicao {
			return paginacao.Posicao{Instante: u.CriadoEm, ID: u.ID}
		})
		return nil
	})
	return pag, err
}

// ------------------------------------------------------- perfis e catálogo

func (s *servico) ListarPerfis(ctx context.Context, a autorizacao.Ator) ([]Perfil, error) {
	if err := apenasInterno(a, autorizacao.AdmLerUsuarios); err != nil {
		return nil, err
	}
	var out []Perfil
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		linhas, err := q.AdmListarPerfis(ctx, a.EmpresaID)
		if err != nil {
			return err
		}
		out = make([]Perfil, 0, len(linhas))
		for _, p := range linhas {
			out = append(out, Perfil{ID: p.ID, Codigo: p.Codigo, Nome: p.Nome, Padrao: p.Padrao,
				PerfilOrigem: p.PerfilOrigem, Permissoes: p.Permissoes, Versao: p.Versao})
		}
		return nil
	})
	return out, err
}

func (s *servico) ListarPermissoes(_ context.Context, a autorizacao.Ator) ([]autorizacao.Permissao, error) {
	if err := apenasInterno(a, autorizacao.AdmGerenciarPerfis); err != nil {
		return nil, err
	}
	return autorizacao.Catalogo(), nil
}
