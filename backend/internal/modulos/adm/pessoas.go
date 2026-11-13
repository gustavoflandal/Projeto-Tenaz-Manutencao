package adm

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"tenaz/internal/modulos/adm/internal/dominio"
	"tenaz/internal/modulos/adm/internal/repo"
	"tenaz/internal/plataforma/auditoria"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/db/consultas"
	"tenaz/internal/plataforma/erros"
	"tenaz/internal/plataforma/paginacao"
)

// Pessoas: RF-ADM-015. A pessoa pode ter ou não usuário; nunca é excluída.

type pessoaValida struct {
	nome                                   string
	cpf, nascimento, matricula, cargo, tel *string
	email, registro                        *string
	vinculo                                string
	fornecedorID, unidadeID                *uuid.UUID
}

func validarPessoa(prefixo string, e EntradaPessoa) (pessoaValida, error) {
	var p pessoaValida
	var ee *erros.Erro
	if p.nome, ee = dominio.Texto(prefixo+"nome", e.Nome, 200); ee != nil {
		return p, ee
	}
	if e.CPF != nil && strings.TrimSpace(*e.CPF) != "" {
		c, ee := dominio.NormalizarCPF(prefixo+"cpf", *e.CPF)
		if ee != nil {
			return p, ee
		}
		p.cpf = &c
	}
	if e.DataNascimento != nil && *e.DataNascimento != "" {
		d, ee := dominio.Data(prefixo+"data_nascimento", *e.DataNascimento)
		if ee != nil {
			return p, ee
		}
		p.nascimento = &d
	}
	textos := []struct {
		campo string
		v     *string
		dst   **string
		max   int
	}{
		{"matricula", e.Matricula, &p.matricula, 40}, {"cargo", e.Cargo, &p.cargo, 120},
		{"telefone", e.Telefone, &p.tel, 30}, {"registro_profissional", e.RegistroProfissional, &p.registro, 60},
	}
	for _, t := range textos {
		v, ee := dominio.TextoOpcional(prefixo+t.campo, t.v, t.max)
		if ee != nil {
			return p, ee
		}
		*t.dst = v
	}
	if e.Email != nil && strings.TrimSpace(*e.Email) != "" {
		m, ee := dominio.Email(prefixo+"email", *e.Email)
		if ee != nil {
			return p, ee
		}
		p.email = &m
	}
	p.vinculo = e.Vinculo
	if p.vinculo == "" {
		p.vinculo = VinculoEmpregado
	}
	switch p.vinculo {
	case VinculoEmpregado, VinculoTerceiro, VinculoAutonomo:
	default:
		return p, erros.EntradaInvalida.ComCampo(prefixo+"vinculo", erros.CampoValorInvalido, "Vínculo: empregado, terceiro ou autonomo.")
	}
	// Técnico de fornecedor é sempre terceiro (ck_pessoa_fornecedor). A
	// existência do fornecedor é conferida quando a tabela existir (E06,
	// pendência E01-1 do documento 2).
	if e.FornecedorID != nil && p.vinculo != VinculoTerceiro {
		return p, erros.EntradaInvalida.ComCampo(prefixo+"vinculo", erros.CampoValorInvalido,
			"Pessoa de fornecedor tem o vínculo terceiro.")
	}
	p.fornecedorID, p.unidadeID = e.FornecedorID, e.UnidadeID
	return p, nil
}

func criarPessoaTx(ctx context.Context, q *consultas.Queries, empresaID, id uuid.UUID, p pessoaValida) error {
	return repo.Traduzir(q.AdmCriarPessoa(ctx, consultas.AdmCriarPessoaParams{
		EmpresaID: empresaID, ID: id, Nome: p.nome, Cpf: p.cpf, DataNascimento: p.nascimento, Matricula: p.matricula,
		Cargo: p.cargo, Telefone: p.tel, Email: p.email, Vinculo: p.vinculo, FornecedorID: p.fornecedorID,
		UnidadeID: p.unidadeID, RegistroProfissional: p.registro,
	}))
}

func lerPessoa(ctx context.Context, empresaID, id uuid.UUID) (Pessoa, error) {
	q, err := repo.Q(ctx)
	if err != nil {
		return Pessoa{}, err
	}
	r, err := q.AdmBuscarPessoa(ctx, consultas.AdmBuscarPessoaParams{EmpresaID: empresaID, ID: id})
	if err != nil {
		return Pessoa{}, naoEncontrado(err)
	}
	return Pessoa{
		ID: r.ID, Nome: r.Nome, CPF: r.Cpf, DataNascimento: vazioParaNil(r.DataNascimento), Matricula: r.Matricula,
		Cargo: r.Cargo, Telefone: r.Telefone, Email: r.Email, Vinculo: r.Vinculo, FornecedorID: r.FornecedorID,
		UnidadeID: r.UnidadeID, RegistroProfissional: r.RegistroProfissional, Situacao: r.Situacao,
		UsuarioID: r.UsuarioID, Versao: r.Versao, CriadoEm: r.CriadoEm, AtualizadoEm: r.AtualizadoEm,
	}, nil
}

// ObterPessoa: adm.ler_pessoas, ou o titular (documento 3, seção 6.1). O CPF
// inteiro é dado pessoal: a leitura dele vai para a trilha como
// visualizacao_sensivel (seção 9.1 do documento 2).
func (s *servico) ObterPessoa(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (Pessoa, error) {
	titular := a.PessoaID != nil && *a.PessoaID == id
	if a.Externo() && !titular {
		return Pessoa{}, erros.NaoEncontrado
	}
	if !titular {
		if err := a.Exigir(autorizacao.AdmLerPessoas); err != nil {
			return Pessoa{}, err
		}
	}
	var p Pessoa
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		var err error
		p, err = lerPessoa(ctx, a.EmpresaID, id)
		if err != nil || p.CPF == nil {
			return err
		}
		return auditoria.Registrar(ctx, auditoria.Evento{
			EmpresaID: a.EmpresaID, AtorTipo: db.AtorUsuario, UsuarioID: ptr(a.UsuarioID),
			Acao: auditoria.AcaoVisualizacaoSensivel, Entidade: "pessoa", EntidadeID: ptr(id),
			Depois: map[string][]string{"campos": {"cpf"}},
		})
	})
	return p, err
}

func (s *servico) CriarPessoa(ctx context.Context, a autorizacao.Ator, e EntradaPessoa) (Pessoa, error) {
	if err := apenasInterno(a, autorizacao.AdmGerenciarPessoas); err != nil {
		return Pessoa{}, err
	}
	pv, err := validarPessoa("", e)
	if err != nil {
		return Pessoa{}, err
	}
	var p Pessoa
	err = s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		id := novoID()
		if err := criarPessoaTx(ctx, q, a.EmpresaID, id, pv); err != nil {
			return err
		}
		p, err = lerPessoa(ctx, a.EmpresaID, id)
		return err
	})
	return p, err
}

func (s *servico) EditarPessoa(ctx context.Context, a autorizacao.Ator, id uuid.UUID, e EntradaPessoa) (Pessoa, error) {
	if err := apenasInterno(a, autorizacao.AdmGerenciarPessoas); err != nil {
		return Pessoa{}, err
	}
	pv, err := validarPessoa("", e)
	if err != nil {
		return Pessoa{}, err
	}
	var p Pessoa
	err = s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		atual, err := q.AdmTravarPessoa(ctx, consultas.AdmTravarPessoaParams{EmpresaID: a.EmpresaID, ID: id})
		if err != nil {
			return naoEncontrado(err)
		}
		if err := conferirVersao(e.VersaoEsperada, atual.Versao); err != nil {
			return err
		}
		// Pessoa ligada a usuário: o fornecedor da pessoa tem de ser o do
		// usuário prestador (RF-ADM-015).
		u, err := q.AdmBuscarUsuarioDaPessoa(ctx, consultas.AdmBuscarUsuarioDaPessoaParams{EmpresaID: a.EmpresaID, PessoaID: &id})
		switch {
		case errors.Is(err, pgx.ErrNoRows):
		case err != nil:
			return err
		case !mesmoID(pv.fornecedorID, u.FornecedorID):
			// Nos dois sentidos: a pessoa de um prestador mantém o fornecedor
			// dele, e a de um usuário interno não ganha fornecedor (RN-ADM-020).
			return erros.EntradaInvalida.ComCampo("fornecedor_id", erros.CampoValorInvalido,
				"A pessoa está ligada a um usuário; o fornecedor dela tem de ser o do usuário (nenhum, se o usuário é interno).")
		}
		err = q.AdmAtualizarPessoa(ctx, consultas.AdmAtualizarPessoaParams{
			EmpresaID: a.EmpresaID, ID: id, Nome: pv.nome, Cpf: pv.cpf, DataNascimento: pv.nascimento,
			Matricula: pv.matricula, Cargo: pv.cargo, Telefone: pv.tel, Email: pv.email, Vinculo: pv.vinculo,
			FornecedorID: pv.fornecedorID, UnidadeID: pv.unidadeID, RegistroProfissional: pv.registro,
		})
		if err != nil {
			return repo.Traduzir(err)
		}
		p, err = lerPessoa(ctx, a.EmpresaID, id)
		return err
	})
	return p, err
}

func mesmoID(a, b *uuid.UUID) bool {
	if a == nil || b == nil {
		return a == b
	}
	return *a == *b
}

func (s *servico) InativarPessoa(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (Pessoa, error) {
	return s.situacaoPessoa(ctx, a, id, "inativa")
}

func (s *servico) ReativarPessoa(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (Pessoa, error) {
	return s.situacaoPessoa(ctx, a, id, "ativa")
}

func (s *servico) situacaoPessoa(ctx context.Context, a autorizacao.Ator, id uuid.UUID, nova string) (Pessoa, error) {
	if err := apenasInterno(a, autorizacao.AdmGerenciarPessoas); err != nil {
		return Pessoa{}, err
	}
	var p Pessoa
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		atual, err := q.AdmTravarPessoa(ctx, consultas.AdmTravarPessoaParams{EmpresaID: a.EmpresaID, ID: id})
		if err != nil {
			return naoEncontrado(err)
		}
		if atual.Situacao == nova {
			return erros.TransicaoInvalida.Com("a pessoa já está %s", nova).ComDado("estado_atual", atual.Situacao)
		}
		// Pessoa ligada a usuário ativo ou convidado não é inativada: o
		// usuário é inativado antes, por quem gerencia usuários (RN-ADM-021).
		if nova == "inativa" {
			u, err := q.AdmBuscarUsuarioDaPessoa(ctx, consultas.AdmBuscarUsuarioDaPessoaParams{EmpresaID: a.EmpresaID, PessoaID: &id})
			switch {
			case errors.Is(err, pgx.ErrNoRows):
			case err != nil:
				return err
			case u.Situacao != "inativo":
				return erros.TransicaoInvalida.Com("a pessoa está ligada a um usuário %s; inative o usuário antes", u.Situacao).
					ComDado("estado_atual", atual.Situacao)
			}
		}
		if err := q.AdmDefinirSituacaoPessoa(ctx, consultas.AdmDefinirSituacaoPessoaParams{
			EmpresaID: a.EmpresaID, ID: id, Situacao: nova,
		}); err != nil {
			return repo.Traduzir(err)
		}
		p, err = lerPessoa(ctx, a.EmpresaID, id)
		return err
	})
	return p, err
}

// padraoBusca monta o ILIKE com os curingas da entrada escapados: "50%" busca
// o texto "50%", não "50 seguido de qualquer coisa".
func padraoBusca(b *string) *string {
	if b == nil || strings.TrimSpace(*b) == "" {
		return nil
	}
	r := strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`)
	v := "%" + r.Replace(strings.TrimSpace(*b)) + "%"
	return &v
}

func cursorParams(p *paginacao.Posicao) (*time.Time, *uuid.UUID) {
	if p == nil {
		return nil, nil
	}
	return &p.Instante, &p.ID
}

func textoOuVazio[T any](v *T, f func(T) string) string {
	if v == nil {
		return ""
	}
	return f(*v)
}

func (s *servico) ListarPessoas(ctx context.Context, a autorizacao.Ator, f FiltroPessoas) (paginacao.Pagina[ResumoPessoa], error) {
	var pag paginacao.Pagina[ResumoPessoa]
	if err := apenasInterno(a, autorizacao.AdmLerPessoas); err != nil {
		return pag, err
	}
	busca := padraoBusca(f.Busca)
	consulta := paginacao.Assinatura("pessoas",
		textoOuVazio(busca, func(s string) string { return s }),
		textoOuVazio(f.Situacao, func(s string) string { return s }),
		textoOuVazio(f.FornecedorID, uuid.UUID.String), textoOuVazio(f.UnidadeID, uuid.UUID.String),
		textoOuVazio(f.ComUsuario, func(b bool) string { return map[bool]string{true: "s", false: "n"}[b] }))
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
		linhas, err := q.AdmListarPessoas(ctx, consultas.AdmListarPessoasParams{
			EmpresaID: a.EmpresaID, Busca: busca, Situacao: f.Situacao, FornecedorID: f.FornecedorID,
			UnidadeID: f.UnidadeID, ComUsuario: f.ComUsuario, CursorCriadoEm: cEm, CursorID: cID, Limite: limite + 1,
		})
		if err != nil {
			return err
		}
		res := make([]ResumoPessoa, 0, len(linhas))
		for _, r := range linhas {
			var mascara *string
			if r.Cpf != nil {
				mascara = ptr(dominio.MascararCPF(*r.Cpf))
			}
			res = append(res, ResumoPessoa{
				ID: r.ID, Nome: r.Nome, CPFMascarado: mascara, Matricula: r.Matricula, Cargo: r.Cargo,
				Vinculo: r.Vinculo, FornecedorID: r.FornecedorID, UnidadeID: r.UnidadeID, Situacao: r.Situacao,
				UsuarioID: r.UsuarioID, Versao: r.Versao, CriadoEm: r.CriadoEm,
			})
		}
		pag = paginacao.Montar(res, limite, consulta, agora, func(p ResumoPessoa) paginacao.Posicao {
			return paginacao.Posicao{Instante: p.CriadoEm, ID: p.ID}
		})
		return nil
	})
	return pag, err
}
