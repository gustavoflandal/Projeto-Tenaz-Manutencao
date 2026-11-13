package adm_test

import (
	"errors"
	"reflect"
	"slices"
	"testing"

	"github.com/google/uuid"

	"tenaz/internal/modulos/adm"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/erros"
)

// Isolamento entre empresas nos casos de uso: todo caso de uso que recebe um
// id responde 404 ao id de outra empresa, com a mesma resposta que daria a
// um id que não existe (nenhum oráculo), e o registro da outra empresa não
// muda.
func TestOutraEmpresaNaoEncontra(t *testing.T) {
	a, b := novaConta(t, "profissional"), novaConta(t, "profissional")
	ctx := ctxWeb()

	// Registros de A, em todas as situações que os casos de uso aceitam.
	unidades, err := svc.ListarUnidades(ctx, a.adm)
	if err != nil {
		t.Fatal(err)
	}
	uA, err := svc.ObterUnidade(ctx, a.adm, unidades[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	uInativa, err := svc.CriarUnidade(ctx, a.adm, adm.EntradaUnidade{Codigo: "F2", Nome: "Filial", Turnos: turnosComerciais()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.InativarUnidade(ctx, a.adm, uInativa.ID); err != nil {
		t.Fatal(err)
	}
	feriado, err := svc.IncluirFeriado(ctx, a.adm, uA.ID, adm.EntradaFeriado{Data: "2027-11-15", Descricao: "República"})
	if err != nil {
		t.Fatal(err)
	}
	pessoa, err := svc.CriarPessoa(ctx, a.adm, adm.EntradaPessoa{Nome: "Pessoa de A"})
	if err != nil {
		t.Fatal(err)
	}
	pInativa, err := svc.CriarPessoa(ctx, a.adm, adm.EntradaPessoa{Nome: "Pessoa inativa de A"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.InativarPessoa(ctx, a.adm, pInativa.ID); err != nil {
		t.Fatal(err)
	}
	tec, _ := convidarEEntrar(t, a, adm.EntradaUsuario{Perfis: []string{"tecnico"}})
	if _, err := svc.ConcederPermissao(ctx, a.adm, tec.ID, autorizacao.OsEncerrar); err != nil {
		t.Fatal(err)
	}
	inativo, _ := convidarEEntrar(t, a, adm.EntradaUsuario{Perfis: []string{"tecnico"}})
	if _, err := svc.InativarUsuario(ctx, a.adm, inativo.ID); err != nil {
		t.Fatal(err)
	}
	convidado, err := svc.CriarUsuario(ctx, a.adm, adm.EntradaUsuario{Nome: "Convidado de A", Email: emailAleatorio("conv"), Perfis: []string{"tecnico"}})
	if err != nil {
		t.Fatal(err)
	}
	evs, err := svc.ListarEventosAuditoria(ctx, a.adm, adm.FiltroAuditoria{Limite: 1})
	if err != nil || len(evs.Dados) == 0 {
		t.Fatalf("evento de A: %v", err)
	}
	bUnidades, err := svc.ListarUnidades(ctx, b.adm)
	if err != nil {
		t.Fatal(err)
	}
	bUnidade := bUnidades[0].ID
	turno := adm.EntradaTurno{DiaSemana: 6, HoraInicio: "08:00", HoraFim: "12:00"}

	// Cada caso roda com o id de A e com um id que não existe.
	casos := []struct {
		nome string
		id   uuid.UUID
		usar func(id uuid.UUID) error
	}{
		{"obter unidade", uA.ID, func(id uuid.UUID) error { _, err := svc.ObterUnidade(ctx, b.adm, id); return err }},
		{"editar unidade", uA.ID, func(id uuid.UUID) error {
			_, err := svc.EditarUnidade(ctx, b.adm, id, adm.EntradaUnidade{Codigo: "X", Nome: "Invadida"})
			return err
		}},
		{"inativar unidade", uA.ID, func(id uuid.UUID) error { _, err := svc.InativarUnidade(ctx, b.adm, id); return err }},
		{"reativar unidade", uInativa.ID, func(id uuid.UUID) error { _, err := svc.ReativarUnidade(ctx, b.adm, id); return err }},
		{"incluir turno", uA.ID, func(id uuid.UUID) error { _, err := svc.IncluirTurno(ctx, b.adm, id, turno); return err }},
		{"alterar turno", uA.ID, func(id uuid.UUID) error {
			_, err := svc.AlterarTurno(ctx, b.adm, id, uA.Turnos[0].ID, turno)
			return err
		}},
		{"alterar turno de A pela unidade de B", uA.Turnos[0].ID, func(id uuid.UUID) error {
			_, err := svc.AlterarTurno(ctx, b.adm, bUnidade, id, turno)
			return err
		}},
		{"remover turno", uA.ID, func(id uuid.UUID) error { return svc.RemoverTurno(ctx, b.adm, id, uA.Turnos[0].ID) }},
		{"remover turno de A pela unidade de B", uA.Turnos[0].ID, func(id uuid.UUID) error { return svc.RemoverTurno(ctx, b.adm, bUnidade, id) }},
		{"listar feriados", uA.ID, func(id uuid.UUID) error { _, err := svc.ListarFeriados(ctx, b.adm, id, nil, nil); return err }},
		{"incluir feriado", uA.ID, func(id uuid.UUID) error {
			_, err := svc.IncluirFeriado(ctx, b.adm, id, adm.EntradaFeriado{Data: "2027-12-08", Descricao: "X"})
			return err
		}},
		{"remover feriado", uA.ID, func(id uuid.UUID) error { return svc.RemoverFeriado(ctx, b.adm, id, feriado.ID) }},
		{"remover feriado de A pela unidade de B", feriado.ID, func(id uuid.UUID) error { return svc.RemoverFeriado(ctx, b.adm, bUnidade, id) }},
		{"obter pessoa", pessoa.ID, func(id uuid.UUID) error { _, err := svc.ObterPessoa(ctx, b.adm, id); return err }},
		{"editar pessoa", pessoa.ID, func(id uuid.UUID) error {
			_, err := svc.EditarPessoa(ctx, b.adm, id, adm.EntradaPessoa{Nome: "Invadida"})
			return err
		}},
		{"inativar pessoa", pessoa.ID, func(id uuid.UUID) error { _, err := svc.InativarPessoa(ctx, b.adm, id); return err }},
		{"reativar pessoa", pInativa.ID, func(id uuid.UUID) error { _, err := svc.ReativarPessoa(ctx, b.adm, id); return err }},
		{"obter usuário", tec.ID, func(id uuid.UUID) error { _, err := svc.ObterUsuario(ctx, b.adm, id); return err }},
		{"editar usuário", tec.ID, func(id uuid.UUID) error {
			_, err := svc.EditarUsuario(ctx, b.adm, id, adm.EntradaEditarUsuario{Nome: "Invadido"})
			return err
		}},
		{"definir perfis", tec.ID, func(id uuid.UUID) error {
			_, err := svc.DefinirPerfis(ctx, b.adm, id, []string{"visualizador"})
			return err
		}},
		{"inativar usuário", tec.ID, func(id uuid.UUID) error { _, err := svc.InativarUsuario(ctx, b.adm, id); return err }},
		{"reativar usuário", inativo.ID, func(id uuid.UUID) error { _, err := svc.ReativarUsuario(ctx, b.adm, id); return err }},
		{"reenviar convite", convidado.Usuario.ID, func(id uuid.UUID) error { return svc.ReenviarConvite(ctx, b.adm, id) }},
		{"encerrar sessões", tec.ID, func(id uuid.UUID) error { _, err := svc.EncerrarSessoes(ctx, b.adm, id); return err }},
		{"conceder permissão", tec.ID, func(id uuid.UUID) error {
			_, err := svc.ConcederPermissao(ctx, b.adm, id, autorizacao.EstDecidirDestino)
			return err
		}},
		{"retirar permissão", tec.ID, func(id uuid.UUID) error {
			_, err := svc.RetirarPermissao(ctx, b.adm, id, autorizacao.OsEncerrar)
			return err
		}},
		{"obter evento", evs.Dados[0].ID, func(id uuid.UUID) error { _, err := svc.ObterEventoAuditoria(ctx, b.adm, id); return err }},
		{"trocar para a empresa", a.criada.EmpresaID, func(id uuid.UUID) error { _, err := auth.TrocarEmpresa(ctx, b.adm, id); return err }},
	}
	for _, c := range casos {
		t.Run(c.nome, func(t *testing.T) {
			deA := c.usar(c.id)
			inexistente := c.usar(uuid.New())
			if !errors.Is(deA, erros.NaoEncontrado) {
				t.Fatalf("id de outra empresa: esperado recurso_nao_encontrado, veio %v", deA)
			}
			if inexistente == nil || deA.Error() != inexistente.Error() {
				t.Errorf("a resposta difere da de um id inexistente (oráculo):\n de A:        %v\n inexistente: %v", deA, inexistente)
			}
		})
	}

	// Ligar à pessoa de A um usuário criado em B: a pessoa não existe para B,
	// com a mesma resposta de uma pessoa inexistente.
	criarComPessoa := func(id uuid.UUID) error {
		_, err := svc.CriarUsuario(ctx, b.adm, adm.EntradaUsuario{Nome: "X", Email: emailAleatorio("x"), Perfis: []string{"tecnico"}, PessoaID: &id})
		return err
	}
	deA, inexistente := criarComPessoa(pessoa.ID), criarComPessoa(uuid.New())
	if !errors.Is(deA, erros.EntradaInvalida) || inexistente == nil || deA.Error() != inexistente.Error() ||
		!reflect.DeepEqual(comoErro(t, deA).Campos, comoErro(t, inexistente).Campos) {
		t.Errorf("pessoa de outra empresa: %v / %v", deA, inexistente)
	}

	// Nada de A mudou.
	if u, err := svc.ObterUnidade(ctx, a.adm, uA.ID); err != nil || u.Nome != uA.Nome || u.Situacao != "ativa" || len(u.Turnos) != len(uA.Turnos) {
		t.Errorf("unidade de A mudou: %+v %v", u, err)
	}
	if fs, err := svc.ListarFeriados(ctx, a.adm, uA.ID, nil, nil); err != nil || len(fs) != 1 {
		t.Errorf("feriados de A: %v %v", fs, err)
	}
	if p, err := svc.ObterPessoa(ctx, a.adm, pessoa.ID); err != nil || p.Nome != "Pessoa de A" || p.Situacao != "ativa" {
		t.Errorf("pessoa de A mudou: %+v %v", p, err)
	}
	if u, err := svc.ObterUsuario(ctx, a.adm, tec.ID); err != nil || u.Situacao != "ativo" || u.Nome != tec.Nome ||
		len(u.Perfis) != 1 || !slices.Equal(u.PermissoesAvulsas, []string{autorizacao.OsEncerrar}) {
		t.Errorf("usuário de A mudou: %+v %v", u, err)
	}
	if u, err := svc.ObterUsuario(ctx, a.adm, inativo.ID); err != nil || u.Situacao != "inativo" {
		t.Errorf("inativo de A mudou: %+v %v", u, err)
	}
	if _, err := auth.AutenticarWeb(ctx, a.cookie, "", false); err != nil {
		t.Errorf("a sessão de A continua: %v", err)
	}
}
