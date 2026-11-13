package adm

import (
	"context"
	"encoding/json"
	"fmt"
	"tenaz/internal/plataforma/textos"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"tenaz/internal/modulos/adm/internal/dominio"
	"tenaz/internal/modulos/adm/internal/repo"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/db/consultas"
	"tenaz/internal/plataforma/erros"
	"tenaz/internal/plataforma/pais"
)

// Unidade, turnos e feriados: RF-ADM-002.

func unidadeDeLinha(r consultas.AdmBuscarUnidadeRow) Unidade {
	u := Unidade{
		ID: r.ID, Codigo: r.Codigo, Nome: r.Nome, IdentificacaoTipo: r.IdentificacaoTipo, Identificacao: r.Identificacao,
		FusoHorario: r.FusoHorario,
		Latitude:    coordenada(r.Latitude), Longitude: coordenada(r.Longitude), Situacao: r.Situacao, Versao: r.Versao,
		CriadoEm: r.CriadoEm, AtualizadoEm: r.AtualizadoEm,
	}
	if len(r.Endereco) > 0 {
		var e Endereco
		if json.Unmarshal(r.Endereco, &e) == nil {
			u.Endereco = &e
		}
	}
	return u
}

func turnoDeLinha(t consultas.AdmListarTurnosRow) Turno {
	return Turno{ID: t.ID, UnidadeID: t.UnidadeID, DiaSemana: int32(t.DiaSemana), HoraInicio: t.HoraInicio,
		HoraFim: t.HoraFim, Nome: t.Nome, Versao: t.Versao}
}

// coordenada converte numeric(9,6) em float64; coordenada não é dinheiro.
func coordenada(n pgtype.Numeric) *float64 {
	if !n.Valid {
		return nil
	}
	f, err := n.Float64Value()
	if err != nil || !f.Valid {
		return nil
	}
	return &f.Float64
}

// lerUnidade lê a unidade com os turnos.
func lerUnidade(ctx context.Context, empresaID, id uuid.UUID) (Unidade, error) {
	q, err := repo.Q(ctx)
	if err != nil {
		return Unidade{}, err
	}
	r, err := q.AdmBuscarUnidade(ctx, consultas.AdmBuscarUnidadeParams{EmpresaID: empresaID, ID: id})
	if err != nil {
		return Unidade{}, naoEncontrado(err)
	}
	u := unidadeDeLinha(r)
	turnos, err := q.AdmListarTurnos(ctx, consultas.AdmListarTurnosParams{EmpresaID: empresaID, UnidadeID: id})
	if err != nil {
		return Unidade{}, err
	}
	u.Turnos = make([]Turno, 0, len(turnos))
	for _, t := range turnos {
		u.Turnos = append(u.Turnos, turnoDeLinha(t))
	}
	return u, nil
}

func (s *servico) ListarUnidades(ctx context.Context, a autorizacao.Ator) ([]Unidade, error) {
	if err := apenasInterno(a, autorizacao.AdmLer); err != nil {
		return nil, err
	}
	var out []Unidade
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		linhas, err := q.AdmListarUnidades(ctx, a.EmpresaID)
		if err != nil {
			return err
		}
		out = make([]Unidade, 0, len(linhas))
		for _, r := range linhas {
			out = append(out, unidadeDeLinha(consultas.AdmBuscarUnidadeRow(r)))
		}
		return nil
	})
	return out, err
}

func (s *servico) ObterUnidade(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (Unidade, error) {
	if err := apenasInterno(a, autorizacao.AdmLer); err != nil {
		return Unidade{}, err
	}
	var u Unidade
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		var err error
		u, err = lerUnidade(ctx, a.EmpresaID, id)
		if err != nil {
			return err
		}
		return s.visualizou(ctx, a, "unidade", id, false)
	})
	return u, err
}

func (s *servico) CriarUnidade(ctx context.Context, a autorizacao.Ator, e EntradaUnidade) (Unidade, error) {
	if err := apenasInterno(a, autorizacao.AdmGerenciarEmpresa); err != nil {
		return Unidade{}, err
	}
	un, err := validarUnidade("", e, true)
	if err != nil {
		return Unidade{}, err
	}
	var u Unidade
	err = s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		// A trava na empresa serializa duas criações simultâneas: sem ela, as
		// duas contariam zero unidades e o Essencial ficaria com duas.
		emp, err := q.AdmTravarEmpresa(ctx, a.EmpresaID)
		if err != nil {
			return naoEncontrado(err)
		}
		n, err := q.AdmContarUnidades(ctx, a.EmpresaID)
		if err != nil {
			return err
		}
		if ee := dominio.PodeCriarUnidade(emp.Plano, n); ee != nil {
			return ee
		}
		// A unidade só aceita o tipo de identificação do país da empresa.
		perfil, ok := pais.Perfil(emp.Pais)
		if !ok {
			return fmt.Errorf("empresa com país %q fora do perfil de país", emp.Pais)
		}
		if ee := dominio.ConferirTipoDaUnidade("", perfil, un.ident); ee != nil {
			return ee
		}
		if ee := conferirFusoDaUnidade("", perfil, un.fuso); ee != nil {
			return ee
		}
		if ee := conferirEnderecoDaUnidade("", perfil, &un, e.Endereco); ee != nil {
			return ee
		}
		id := novoID()
		if err := criarUnidadeTx(ctx, q, a.EmpresaID, id, un); err != nil {
			return err
		}
		u, err = lerUnidade(ctx, a.EmpresaID, id)
		return err
	})
	return u, err
}

func (s *servico) EditarUnidade(ctx context.Context, a autorizacao.Ator, id uuid.UUID, e EntradaUnidade) (Unidade, error) {
	if err := apenasInterno(a, autorizacao.AdmGerenciarEmpresa); err != nil {
		return Unidade{}, err
	}
	if len(e.Turnos) > 0 {
		return Unidade{}, erros.EntradaInvalida.ComCampoChave("turnos", erros.CampoValorInvalido, "campo.turnos_na_edicao", nil)
	}
	un, err := validarUnidade("", e, false)
	if err != nil {
		return Unidade{}, err
	}
	var u Unidade
	err = s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		atual, err := q.AdmTravarUnidade(ctx, consultas.AdmTravarUnidadeParams{EmpresaID: a.EmpresaID, ID: id})
		if err != nil {
			return naoEncontrado(err)
		}
		if err := conferirVersao(e.VersaoEsperada, atual.Versao); err != nil {
			return err
		}
		emp, err := q.AdmBuscarEmpresa(ctx, a.EmpresaID)
		if err != nil {
			return naoEncontrado(err)
		}
		perfil, ok := pais.Perfil(emp.Pais)
		if !ok {
			return fmt.Errorf("empresa com país %q fora do perfil de país", emp.Pais)
		}
		if ee := dominio.ConferirTipoDaUnidade("", perfil, un.ident); ee != nil {
			return ee
		}
		if ee := conferirFusoDaUnidade("", perfil, un.fuso); ee != nil {
			return ee
		}
		if ee := conferirEnderecoDaUnidade("", perfil, &un, e.Endereco); ee != nil {
			return ee
		}
		err = q.AdmAtualizarUnidade(ctx, consultas.AdmAtualizarUnidadeParams{
			EmpresaID: a.EmpresaID, ID: id, Codigo: un.codigo, Nome: un.nome, IdentificacaoTipo: un.identTipo(),
			Identificacao: un.identValor(), Endereco: un.endereco,
			FusoHorario: un.fuso, Latitude: un.lat, Longitude: un.lon,
		})
		if err != nil {
			return repo.Traduzir(err)
		}
		u, err = lerUnidade(ctx, a.EmpresaID, id)
		return err
	})
	return u, err
}

// InativarUnidade: a unidade com ativos, OS ou movimentações não é excluída,
// só inativada (RF-ADM-002). A última unidade ativa não é inativada
// (proposta): a empresa ficaria sem onde cadastrar. A conferência de ativos e
// OS entra com ATV (E02).
func (s *servico) InativarUnidade(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (Unidade, error) {
	return s.situacaoUnidade(ctx, a, id, "inativa")
}

func (s *servico) ReativarUnidade(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (Unidade, error) {
	return s.situacaoUnidade(ctx, a, id, "ativa")
}

func (s *servico) situacaoUnidade(ctx context.Context, a autorizacao.Ator, id uuid.UUID, nova string) (Unidade, error) {
	if err := apenasInterno(a, autorizacao.AdmGerenciarEmpresa); err != nil {
		return Unidade{}, err
	}
	var u Unidade
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		if _, err := q.AdmTravarEmpresa(ctx, a.EmpresaID); err != nil {
			return naoEncontrado(err)
		}
		atual, err := q.AdmTravarUnidade(ctx, consultas.AdmTravarUnidadeParams{EmpresaID: a.EmpresaID, ID: id})
		if err != nil {
			return naoEncontrado(err)
		}
		if atual.Situacao == nova {
			return erros.TransicaoInvalida.ComChave("detalhe.unidade_ja_situacao", textos.Parametros{"situacao": textos.Chave("situacao." + nova)}).ComDado("estado_atual", atual.Situacao)
		}
		if nova == "inativa" {
			n, err := q.AdmContarUnidadesAtivas(ctx, a.EmpresaID)
			if err != nil {
				return err
			}
			if n <= 1 {
				return erros.TransicaoInvalida.ComChave("detalhe.empresa_uma_unidade_ativa", nil).
					ComDado("estado_atual", atual.Situacao)
			}
		}
		if err := q.AdmDefinirSituacaoUnidade(ctx, consultas.AdmDefinirSituacaoUnidadeParams{
			EmpresaID: a.EmpresaID, ID: id, Situacao: nova,
		}); err != nil {
			return repo.Traduzir(err)
		}
		u, err = lerUnidade(ctx, a.EmpresaID, id)
		return err
	})
	return u, err
}

// ------------------------------------------------------------------ turnos

func validarEntradaTurno(e EntradaTurno) (dominio.Turno, *string, error) {
	t := dominio.Turno{DiaSemana: e.DiaSemana, HoraInicio: e.HoraInicio, HoraFim: e.HoraFim}
	if ee := dominio.ValidarTurno("", t); ee != nil {
		return t, nil, ee
	}
	nome, ee := dominio.TextoOpcional("nome", e.Nome, 60)
	if ee != nil {
		return t, nil, ee
	}
	return t, nome, nil
}

// conferirSobreposicao compara o turno com os demais da unidade (RN-ADM-019).
// exceto é o próprio turno, na alteração. Roda com a unidade travada, para
// que duas inclusões simultâneas não passem uma sem ver a outra.
func conferirSobreposicao(ctx context.Context, q *consultas.Queries, empresaID, unidadeID uuid.UUID, t dominio.Turno, exceto *uuid.UUID) error {
	atuais, err := q.AdmListarTurnos(ctx, consultas.AdmListarTurnosParams{EmpresaID: empresaID, UnidadeID: unidadeID})
	if err != nil {
		return err
	}
	for _, x := range atuais {
		if exceto != nil && x.ID == *exceto {
			continue
		}
		outro := dominio.Turno{DiaSemana: int32(x.DiaSemana), HoraInicio: x.HoraInicio, HoraFim: x.HoraFim}
		if dominio.Sobrepoem(t, outro) {
			return dominio.ErrSobreposicao("hora_inicio", outro)
		}
	}
	return nil
}

// unidadeDaEmpresa trava a unidade, conferindo que existe nesta empresa.
func unidadeDaEmpresa(ctx context.Context, q *consultas.Queries, empresaID, unidadeID uuid.UUID) error {
	_, err := q.AdmTravarUnidade(ctx, consultas.AdmTravarUnidadeParams{EmpresaID: empresaID, ID: unidadeID})
	return naoEncontrado(err)
}

func (s *servico) IncluirTurno(ctx context.Context, a autorizacao.Ator, unidadeID uuid.UUID, e EntradaTurno) (Turno, error) {
	if err := apenasInterno(a, autorizacao.AdmGerenciarEmpresa); err != nil {
		return Turno{}, err
	}
	t, nome, err := validarEntradaTurno(e)
	if err != nil {
		return Turno{}, err
	}
	var out Turno
	err = s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		if err := unidadeDaEmpresa(ctx, q, a.EmpresaID, unidadeID); err != nil {
			return err
		}
		if err := conferirSobreposicao(ctx, q, a.EmpresaID, unidadeID, t, nil); err != nil {
			return err
		}
		id := novoID()
		err = q.AdmCriarTurno(ctx, consultas.AdmCriarTurnoParams{
			EmpresaID: a.EmpresaID, ID: id, UnidadeID: unidadeID, DiaSemana: int16(t.DiaSemana),
			HoraInicio: t.HoraInicio, HoraFim: t.HoraFim, Nome: nome,
		})
		if err != nil {
			return repo.Traduzir(err)
		}
		out, err = lerTurno(ctx, q, a.EmpresaID, unidadeID, id)
		return err
	})
	return out, err
}

func lerTurno(ctx context.Context, q *consultas.Queries, empresaID, unidadeID, id uuid.UUID) (Turno, error) {
	r, err := q.AdmBuscarTurno(ctx, consultas.AdmBuscarTurnoParams{EmpresaID: empresaID, ID: id})
	if err != nil {
		return Turno{}, naoEncontrado(err)
	}
	if r.UnidadeID != unidadeID {
		return Turno{}, erros.NaoEncontrado
	}
	return turnoDeLinha(consultas.AdmListarTurnosRow(r)), nil
}

func (s *servico) AlterarTurno(ctx context.Context, a autorizacao.Ator, unidadeID, turnoID uuid.UUID, e EntradaTurno) (Turno, error) {
	if err := apenasInterno(a, autorizacao.AdmGerenciarEmpresa); err != nil {
		return Turno{}, err
	}
	t, nome, err := validarEntradaTurno(e)
	if err != nil {
		return Turno{}, err
	}
	var out Turno
	err = s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		if err := unidadeDaEmpresa(ctx, q, a.EmpresaID, unidadeID); err != nil {
			return err
		}
		if _, err := lerTurno(ctx, q, a.EmpresaID, unidadeID, turnoID); err != nil {
			return err
		}
		if err := conferirSobreposicao(ctx, q, a.EmpresaID, unidadeID, t, &turnoID); err != nil {
			return err
		}
		err = q.AdmAtualizarTurno(ctx, consultas.AdmAtualizarTurnoParams{
			EmpresaID: a.EmpresaID, ID: turnoID, DiaSemana: int16(t.DiaSemana), HoraInicio: t.HoraInicio,
			HoraFim: t.HoraFim, Nome: nome,
		})
		if err != nil {
			return repo.Traduzir(err)
		}
		out, err = lerTurno(ctx, q, a.EmpresaID, unidadeID, turnoID)
		return err
	})
	return out, err
}

// RemoverTurno recusa o último: "pelo menos um turno" é regra do serviço
// (RF-ADM-002; decisão da seção 5.1 do documento 2). A trava na unidade
// serializa remoções simultâneas.
func (s *servico) RemoverTurno(ctx context.Context, a autorizacao.Ator, unidadeID, turnoID uuid.UUID) error {
	if err := apenasInterno(a, autorizacao.AdmGerenciarEmpresa); err != nil {
		return err
	}
	return s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		if err := unidadeDaEmpresa(ctx, q, a.EmpresaID, unidadeID); err != nil {
			return err
		}
		if _, err := lerTurno(ctx, q, a.EmpresaID, unidadeID, turnoID); err != nil {
			return err
		}
		n, err := q.AdmContarTurnos(ctx, consultas.AdmContarTurnosParams{EmpresaID: a.EmpresaID, UnidadeID: unidadeID})
		if err != nil {
			return err
		}
		if n <= 1 {
			return dominio.ErrUltimoTurno()
		}
		return q.AdmRemoverTurno(ctx, consultas.AdmRemoverTurnoParams{EmpresaID: a.EmpresaID, ID: turnoID})
	})
}

// ---------------------------------------------------------------- feriados

func (s *servico) ListarFeriados(ctx context.Context, a autorizacao.Ator, unidadeID uuid.UUID, de, ate *string) ([]Feriado, error) {
	if err := apenasInterno(a, autorizacao.AdmLer); err != nil {
		return nil, err
	}
	for campo, v := range map[string]*string{"de": de, "ate": ate} {
		if v != nil {
			if _, ee := dominio.Data(campo, *v); ee != nil {
				return nil, ee
			}
		}
	}
	var out []Feriado
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		if _, err := q.AdmBuscarUnidade(ctx, consultas.AdmBuscarUnidadeParams{EmpresaID: a.EmpresaID, ID: unidadeID}); err != nil {
			return naoEncontrado(err)
		}
		linhas, err := q.AdmListarFeriados(ctx, consultas.AdmListarFeriadosParams{
			EmpresaID: a.EmpresaID, UnidadeID: unidadeID, De: de, Ate: ate,
		})
		if err != nil {
			return err
		}
		out = make([]Feriado, 0, len(linhas))
		for _, f := range linhas {
			out = append(out, Feriado{ID: f.ID, UnidadeID: f.UnidadeID, Data: f.DataFeriado, Descricao: f.Descricao, Versao: f.Versao})
		}
		return nil
	})
	return out, err
}

func (s *servico) IncluirFeriado(ctx context.Context, a autorizacao.Ator, unidadeID uuid.UUID, e EntradaFeriado) (Feriado, error) {
	if err := apenasInterno(a, autorizacao.AdmGerenciarEmpresa); err != nil {
		return Feriado{}, err
	}
	data, ee := dominio.Data("data", e.Data)
	if ee != nil {
		return Feriado{}, ee
	}
	desc, ee := dominio.Texto("descricao", e.Descricao, 120)
	if ee != nil {
		return Feriado{}, ee
	}
	var out Feriado
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		if err := unidadeDaEmpresa(ctx, q, a.EmpresaID, unidadeID); err != nil {
			return err
		}
		id := novoID()
		if err := q.AdmCriarFeriado(ctx, consultas.AdmCriarFeriadoParams{
			EmpresaID: a.EmpresaID, ID: id, UnidadeID: unidadeID, DataFeriado: data, Descricao: desc,
		}); err != nil {
			return repo.Traduzir(err)
		}
		out = Feriado{ID: id, UnidadeID: unidadeID, Data: data, Descricao: desc, Versao: 1}
		return nil
	})
	return out, err
}

func (s *servico) RemoverFeriado(ctx context.Context, a autorizacao.Ator, unidadeID, feriadoID uuid.UUID) error {
	if err := apenasInterno(a, autorizacao.AdmGerenciarEmpresa); err != nil {
		return err
	}
	return s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		n, err := q.AdmRemoverFeriado(ctx, consultas.AdmRemoverFeriadoParams{EmpresaID: a.EmpresaID, UnidadeID: unidadeID, ID: feriadoID})
		if err != nil {
			return err
		}
		if n == 0 {
			return erros.NaoEncontrado
		}
		return nil
	})
}
