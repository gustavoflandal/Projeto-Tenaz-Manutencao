package adm_test

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"slices"
	"strings"
	"sync"
	"testing"

	"github.com/google/uuid"

	"tenaz/internal/modulos/adm"
	"tenaz/internal/plataforma/autenticacao"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/erros"
)

// Testes de regra do ADM pedidos pelo tenaz-qa-automacao na etapa E01.

// ------------------------------------------------- unicidade por empresa

// CPF e e-mail são únicos por empresa, não no sistema (RF-ADM-004,
// RF-ADM-015): a empresa A aceita o CPF e o e-mail que já existem em B, e
// nada do que A recebe revela que eles existem em B.
func TestMesmoCPFEEmailEmOutraEmpresa(t *testing.T) {
	a, b := novaConta(t, ""), novaConta(t, "")
	ctx := ctxWeb()

	cpf := cpfAleatorio()
	pb, err := svc.CriarPessoa(ctx, b.adm, adm.EntradaPessoa{Nome: "Carlos em B", CPF: &cpf})
	if err != nil {
		t.Fatal(err)
	}
	pa, err := svc.CriarPessoa(ctx, a.adm, adm.EntradaPessoa{Nome: "Carlos em A", CPF: &cpf})
	if err != nil {
		t.Fatalf("o CPF de B é aceito em A: %v", err)
	}
	_, err = svc.CriarPessoa(ctx, a.adm, adm.EntradaPessoa{Nome: "Repetido", CPF: &cpf})
	exigirErro(t, err, erros.Duplicado) // controle: dentro de A, é único

	// O e-mail já entra em B (usuário ativo, com identidade).
	ub, _ := convidarEEntrar(t, b, adm.EntradaUsuario{Perfis: []string{"tecnico"}})
	jaExiste, err := svc.CriarUsuario(ctx, a.adm, adm.EntradaUsuario{Nome: "Mesmo e-mail", Email: *ub.Email, Perfis: []string{"tecnico"}})
	if err != nil {
		t.Fatalf("o e-mail de B é aceito em A: %v", err)
	}
	inedito, err := svc.CriarUsuario(ctx, a.adm, adm.EntradaUsuario{Nome: "E-mail novo", Email: emailAleatorio("novo"), Perfis: []string{"tecnico"}})
	if err != nil {
		t.Fatal(err)
	}
	_, err = svc.CriarUsuario(ctx, a.adm, adm.EntradaUsuario{Nome: "Repetido", Email: *ub.Email, Perfis: []string{"tecnico"}})
	exigirErro(t, err, erros.Duplicado) // controle: dentro de A, é único

	// O administrador de A recebe a mesma resposta para o e-mail que existe
	// em B e para um e-mail inédito.
	type visao struct {
		IdentidadeLigada  bool
		Situacao, Escopo  string
		Versao            int64
		Perfis            int
		Avulsas           []string
		Pessoa            *uuid.UUID
		ConviteEnviado    bool
		Fornecedor, Locat *uuid.UUID
	}
	ver := func(u adm.UsuarioCriado) visao {
		return visao{u.Usuario.IdentidadeLigada, u.Usuario.Situacao, u.Usuario.EscopoAcesso, u.Usuario.Versao,
			len(u.Usuario.Perfis), u.Usuario.PermissoesAvulsas, u.Usuario.PessoaID, u.ConviteEnviado,
			u.Usuario.FornecedorID, u.Usuario.LocatarioID}
	}
	if !reflect.DeepEqual(ver(jaExiste), ver(inedito)) {
		t.Errorf("a resposta revela o e-mail de outra empresa:\n existe em B: %+v\n inédito:     %+v", ver(jaExiste), ver(inedito))
	}
	for _, id := range []uuid.UUID{jaExiste.Usuario.ID, inedito.Usuario.ID} {
		u, err := svc.ObterUsuario(ctx, a.adm, id)
		if err != nil || u.IdentidadeLigada || u.Situacao != "convidado" {
			t.Errorf("detalhe em A antes do aceite: %+v %v", u, err)
		}
	}
	lista, err := svc.ListarUsuarios(ctx, a.adm, adm.FiltroUsuarios{Busca: ub.Email})
	if err != nil || len(lista.Dados) != 1 || lista.Dados[0].ID != jaExiste.Usuario.ID {
		t.Errorf("a lista de A traz só o usuário de A: %+v %v", lista.Dados, err)
	}

	// Nada na trilha de A cita B: nem a empresa, nem o usuário, nem a pessoa.
	proibidos := []string{b.criada.EmpresaID.String(), ub.ID.String(), pb.ID.String(), b.adm.UsuarioID.String()}
	pag, err := svc.ListarEventosAuditoria(ctx, a.adm, adm.FiltroAuditoria{Limite: 100})
	if err != nil || len(pag.Dados) == 0 {
		t.Fatalf("trilha de A: %v", err)
	}
	for _, r := range pag.Dados {
		ev, err := svc.ObterEventoAuditoria(ctx, a.adm, r.ID)
		if err != nil {
			t.Fatal(err)
		}
		bruto, _ := json.Marshal(ev)
		for _, p := range proibidos {
			if strings.Contains(string(bruto)+string(ev.Antes)+string(ev.Depois), p) {
				t.Errorf("evento %s/%s de A cita %s, de B", ev.Entidade, ev.Acao, p)
			}
		}
	}
	if det, err := svc.ObterPessoa(ctx, a.adm, pa.ID); err != nil || det.CPF == nil || *det.CPF != cpf {
		t.Errorf("pessoa de A com o CPF: %v", err)
	}
}

// ---------------------------------------------- último administrador

// Dois administradores rebaixam um ao outro ao mesmo tempo: a trava no perfil
// administrador serializa as duas transações, e a segunda falha com
// ultimo_administrador (RF-ADM-004). Nunca as duas passam.
func TestUltimoAdministradorConcorrente(t *testing.T) {
	for rodada := 0; rodada < 4; rodada++ {
		c := novaConta(t, "")
		outro, aOutro := convidarEEntrar(t, c, adm.EntradaUsuario{Perfis: []string{"administrador", "gestor"}})
		if !aOutro.TemPerfil(autorizacao.PerfilAdministrador) {
			t.Fatal("o segundo usuário é administrador")
		}
		var (
			inicio = make(chan struct{})
			wg     sync.WaitGroup
			errs   [2]error
		)
		rebaixar := func(i int, ator autorizacao.Ator, alvo uuid.UUID) {
			defer wg.Done()
			<-inicio
			_, errs[i] = svc.DefinirPerfis(ctxWeb(), ator, alvo, []string{"gestor"})
		}
		wg.Add(2)
		go rebaixar(0, c.adm, outro.ID)
		go rebaixar(1, aOutro, c.adm.UsuarioID)
		close(inicio)
		wg.Wait()

		ok, ultimo := 0, 0
		for _, err := range errs {
			switch {
			case err == nil:
				ok++
			case errors.Is(err, adm.ErrUltimoAdministrador):
				ultimo++
			default:
				t.Fatalf("rodada %d: erro inesperado (bloqueio mútuo?): %v", rodada, err)
			}
		}
		if ok != 1 || ultimo != 1 {
			t.Fatalf("rodada %d: esperado um sucesso e um ultimo_administrador; vieram %v", rodada, errs)
		}
		// Sobra exatamente um administrador ativo.
		admins := 0
		for _, id := range []uuid.UUID{c.adm.UsuarioID, outro.ID} {
			u, err := svc.ObterUsuario(ctxWeb(), c.adm, id)
			if err != nil {
				t.Fatal(err)
			}
			for _, p := range u.Perfis {
				if p.CodigoOrigem == autorizacao.PerfilAdministrador {
					admins++
				}
			}
		}
		if admins != 1 {
			t.Fatalf("rodada %d: %d administradores depois das duas operações", rodada, admins)
		}
	}
}

// Inativar e rebaixar ao mesmo tempo: as duas tiram um administrador; uma
// tem de falhar. A falha pode ser ultimo_administrador ou credencial_invalida
// (a inativação derruba a credencial de quem rebaixa), então este teste não
// prova sozinho a trava no perfil administrador: essa prova é
// TestUltimoAdministradorConcorrente, em que só ultimo_administrador é aceito.
func TestUltimoAdministradorInativarERebaixarJuntos(t *testing.T) {
	c := novaConta(t, "")
	outro, aOutro := convidarEEntrar(t, c, adm.EntradaUsuario{Perfis: []string{"administrador", "gestor"}})
	var (
		inicio = make(chan struct{})
		wg     sync.WaitGroup
		errs   [2]error
	)
	wg.Add(2)
	go func() {
		defer wg.Done()
		<-inicio
		_, errs[0] = svc.InativarUsuario(ctxWeb(), c.adm, outro.ID)
	}()
	go func() {
		defer wg.Done()
		<-inicio
		_, errs[1] = svc.DefinirPerfis(ctxWeb(), aOutro, c.adm.UsuarioID, []string{"gestor"})
	}()
	close(inicio)
	wg.Wait()
	falhas := 0
	for _, err := range errs {
		if err != nil {
			// A inativação do outro pode derrubar a credencial dele no meio
			// (credencial_invalida, pelo conferir_escopo) antes da trava.
			if !errors.Is(err, adm.ErrUltimoAdministrador) && !errors.Is(err, erros.CredencialInvalida) {
				t.Fatalf("erro inesperado: %v", err)
			}
			falhas++
		}
	}
	if falhas != 1 {
		t.Fatalf("esperada exatamente uma falha: %v", errs)
	}
}

// --------------------------------------------- situação do usuário

// Tabela de situação do usuário (RN-ADM-023; `usuario.situacao`: convidado,
// ativo, inativo).
func TestSituacaoDoUsuario(t *testing.T) {
	c := novaConta(t, "")
	ctx := ctxWeb()
	estadoAtual := func(t *testing.T, err error, esperado string) {
		t.Helper()
		e := comoErro(t, err)
		if !errors.Is(e, erros.TransicaoInvalida) || e.Dados["estado_atual"] != esperado {
			t.Fatalf("esperado transicao_invalida com estado_atual %q, veio %v %v", esperado, err, e.Dados)
		}
	}
	situacao := func(t *testing.T, id uuid.UUID) adm.Usuario {
		t.Helper()
		u, err := svc.ObterUsuario(ctx, c.adm, id)
		if err != nil {
			t.Fatal(err)
		}
		return u
	}

	// convidado -> aceitar -> ativo
	em := emailAleatorio("situacao")
	criado, err := svc.CriarUsuario(ctx, c.adm, adm.EntradaUsuario{Nome: "Situação", Email: em, Perfis: []string{"tecnico"}})
	if err != nil || criado.Usuario.Situacao != "convidado" || criado.Usuario.IdentidadeLigada {
		t.Fatalf("nasce convidado: %+v %v", criado.Usuario, err)
	}
	id := criado.Usuario.ID
	if _, err := svc.AceitarConvite(ctx, caixa.token(t, em), senhaBoa); err != nil {
		t.Fatal(err)
	}
	if u := situacao(t, id); u.Situacao != "ativo" || !u.IdentidadeLigada {
		t.Fatalf("aceitar ativa: %+v", u)
	}

	// ativo: reativar 409, reenviar convite 409
	_, err = svc.ReativarUsuario(ctx, c.adm, id)
	estadoAtual(t, err, "ativo")
	estadoAtual(t, svc.ReenviarConvite(ctx, c.adm, id), "ativo")

	// ativo -> inativar -> inativo, com as sessões revogadas
	_, sessao := entrar(t, em, senhaBoa, &c.criada.EmpresaID)
	if u, err := svc.InativarUsuario(ctx, c.adm, id); err != nil || u.Situacao != "inativo" {
		t.Fatalf("inativar: %+v %v", u, err)
	}
	_, err = auth.AutenticarWeb(ctxWeb(), sessao, "", false)
	exigirErro(t, err, erros.CredencialInvalida)
	var vivas int
	if err := amb.Migrador.QueryRow(context.Background(),
		`SELECT count(*) FROM plataforma.sessao WHERE empresa_id = $1 AND usuario_id = $2 AND revogada_em IS NULL`,
		c.criada.EmpresaID, id).Scan(&vivas); err != nil || vivas != 0 {
		t.Fatalf("sessões vivas do inativo: %d %v", vivas, err)
	}
	_, err = auth.Login(ctx, autenticacao.EntradaLogin{Email: em, Senha: senhaBoa})
	exigirErro(t, err, erros.CredencialInvalida) // sem empresa habilitada

	// inativo: inativar 409, reenviar convite 409
	_, err = svc.InativarUsuario(ctx, c.adm, id)
	estadoAtual(t, err, "inativo")
	estadoAtual(t, svc.ReenviarConvite(ctx, c.adm, id), "inativo")

	// inativo com identidade -> reativar -> ativo, sem convite
	re, err := svc.ReativarUsuario(ctx, c.adm, id)
	if err != nil || re.Usuario.Situacao != "ativo" || re.ConviteEnviado {
		t.Fatalf("reativar quem já aceitou: %+v %v", re, err)
	}
	entrar(t, em, senhaBoa, &c.criada.EmpresaID)

	// convidado (sem identidade) -> inativar -> o link deixa de valer
	em2 := emailAleatorio("nunca-aceitou")
	nunca, err := svc.CriarUsuario(ctx, c.adm, adm.EntradaUsuario{Nome: "Nunca aceitou", Email: em2, Perfis: []string{"tecnico"}})
	if err != nil {
		t.Fatal(err)
	}
	link1 := caixa.token(t, em2)
	if _, err := svc.InativarUsuario(ctx, c.adm, nunca.Usuario.ID); err != nil {
		t.Fatal(err)
	}
	_, err = svc.AceitarConvite(ctx, link1, senhaBoa)
	exigirErro(t, err, erros.LinkEmailInvalido)
	if u := situacao(t, nunca.Usuario.ID); u.Situacao != "inativo" || u.IdentidadeLigada {
		t.Fatalf("link de usuário inativado não ativa: %+v", u)
	}

	// inativo sem identidade -> reativar -> convidado, com convite novo
	re, err = svc.ReativarUsuario(ctx, c.adm, nunca.Usuario.ID)
	if err != nil || re.Usuario.Situacao != "convidado" || !re.ConviteEnviado {
		t.Fatalf("reativar quem nunca aceitou: %+v %v", re, err)
	}
	link2 := caixa.token(t, em2)
	if link2 == link1 {
		t.Fatal("a reativação emite convite novo")
	}
	_, err = svc.AceitarConvite(ctx, link1, senhaBoa)
	exigirErro(t, err, erros.LinkEmailInvalido) // o anterior continua inválido
	if _, err := svc.AceitarConvite(ctx, link2, senhaBoa); err != nil {
		t.Fatalf("o convite novo vale: %v", err)
	}
	if u := situacao(t, nunca.Usuario.ID); u.Situacao != "ativo" {
		t.Fatalf("aceitar o convite novo ativa: %+v", u)
	}
}

// --------------------------------------------------------- auditoria

// eventosDe devolve o detalhe dos eventos da entidade, do mais novo ao mais
// antigo, filtrando pelo id do registro ou, em tabela de ligação, pelo
// usuario_id dentro de antes ou depois.
func eventosDe(t *testing.T, c conta, entidade, acao string, entidadeID *uuid.UUID, usuarioLigado *uuid.UUID) []adm.Evento {
	t.Helper()
	ctx := ctxWeb()
	pag, err := svc.ListarEventosAuditoria(ctx, c.adm, adm.FiltroAuditoria{Entidade: &entidade, Acao: &acao, EntidadeID: entidadeID, Limite: 100})
	if err != nil {
		t.Fatal(err)
	}
	var out []adm.Evento
	for _, r := range pag.Dados {
		ev, err := svc.ObterEventoAuditoria(ctx, c.adm, r.ID)
		if err != nil {
			t.Fatal(err)
		}
		if usuarioLigado != nil {
			var x struct {
				UsuarioID *uuid.UUID `json:"usuario_id"`
			}
			lado := ev.Depois
			if acao == "exclusao" {
				lado = ev.Antes
			}
			if json.Unmarshal(lado, &x) != nil || x.UsuarioID == nil || *x.UsuarioID != *usuarioLigado {
				continue
			}
		}
		out = append(out, ev)
	}
	return out
}

// campo lê um campo de antes ou depois.
func campo(t *testing.T, bruto []byte, nome string) any {
	t.Helper()
	if len(bruto) == 0 {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(bruto, &m); err != nil {
		t.Fatalf("JSON inválido na trilha: %s", bruto)
	}
	return m[nome]
}

// A trilha registra antes, depois, autor e instante das alterações de
// unidade, pessoa, usuário, perfil e permissão avulsa (RF-ADM-011; seção 9.3
// do documento 2), com o contexto técnico da requisição.
func TestAuditoriaDasAlteracoes(t *testing.T) {
	c := novaConta(t, "profissional")
	reqID := uuid.NewString()
	ctx := db.ComRequisicao(context.Background(), db.Requisicao{
		Origem: db.OrigemWeb, IP: "198.51.100.23", UserAgent: "teste-auditoria-adm", RequestID: reqID,
	})
	autor := c.adm.UsuarioID
	inicio := relogioDoBanco(t)

	// conferir confere autor, instante e contexto técnico.
	conferir := func(t *testing.T, nome string, evs []adm.Evento, quantos int) adm.Evento {
		t.Helper()
		if len(evs) != quantos {
			t.Fatalf("%s: %d eventos, esperado %d", nome, len(evs), quantos)
		}
		e := evs[0]
		if e.AtorTipo != db.AtorUsuario || e.UsuarioID == nil || *e.UsuarioID != autor {
			t.Errorf("%s: autor %s %v, esperado o administrador %s", nome, e.AtorTipo, e.UsuarioID, autor)
		}
		if e.OcorridoEm.Before(inicio) || e.OcorridoEm.After(relogioDoBanco(t)) {
			t.Errorf("%s: instante %v fora do teste", nome, e.OcorridoEm)
		}
		if e.Origem == nil || *e.Origem != db.OrigemWeb || e.IP == nil || *e.IP != "198.51.100.23" ||
			e.RequestID == nil || *e.RequestID != reqID || e.UserAgent == nil || *e.UserAgent != "teste-auditoria-adm" {
			t.Errorf("%s: contexto técnico %v %v %v %v", nome, e.Origem, e.IP, e.RequestID, e.UserAgent)
		}
		return e
	}
	var conferirRetirada func(t *testing.T, nome string, inclusao adm.Evento)
	antesDepois := func(t *testing.T, nome string, e adm.Evento, chave string, antes, depois any) {
		t.Helper()
		if a := campo(t, e.Antes, chave); !reflect.DeepEqual(a, antes) {
			t.Errorf("%s: antes.%s = %v, esperado %v", nome, chave, a, antes)
		}
		if d := campo(t, e.Depois, chave); !reflect.DeepEqual(d, depois) {
			t.Errorf("%s: depois.%s = %v, esperado %v", nome, chave, d, depois)
		}
	}

	// A retirada de perfil e de permissão avulsa é exclusão lógica
	// (excluido_em): a trilha a grava como alteração de excluido_em, ligada à
	// inclusão pelo entidade_id (o id da linha). Pendência registrada: o
	// tg_auditar não a mostra como exclusao nem repete o código retirado.
	conferirRetirada = func(t *testing.T, nome string, inclusao adm.Evento) {
		t.Helper()
		e := conferir(t, nome, eventosDe(t, c, inclusao.Entidade, "alteracao", inclusao.EntidadeID, nil), 1)
		if !slices.Equal(e.CamposAlterados, []string{"excluido_em"}) || campo(t, e.Antes, "excluido_em") != nil ||
			campo(t, e.Depois, "excluido_em") == nil {
			t.Errorf("%s: esperada a alteração de excluido_em; campos %v, antes %s, depois %s", nome, e.CamposAlterados, e.Antes, e.Depois)
		}
	}

	t.Run("unidade", func(t *testing.T) {
		u, err := svc.CriarUnidade(ctx, c.adm, adm.EntradaUnidade{Codigo: "AUD", Nome: "Filial Auditada",
			Turnos: []adm.EntradaTurno{{DiaSemana: 2, HoraInicio: "07:00", HoraFim: "16:00"}}})
		if err != nil {
			t.Fatal(err)
		}
		e := conferir(t, "inclusão de unidade", eventosDe(t, c, "unidade", "inclusao", &u.ID, nil), 1)
		antesDepois(t, "inclusão de unidade", e, "nome", nil, "Filial Auditada")
		if _, err := svc.EditarUnidade(ctx, c.adm, u.ID, adm.EntradaUnidade{Codigo: "AUD", Nome: "Filial Renomeada"}); err != nil {
			t.Fatal(err)
		}
		e = conferir(t, "alteração de unidade", eventosDe(t, c, "unidade", "alteracao", &u.ID, nil), 1)
		antesDepois(t, "alteração de unidade", e, "nome", "Filial Auditada", "Filial Renomeada")
		if !slices.Contains(e.CamposAlterados, "nome") || slices.Contains(e.CamposAlterados, "versao") {
			t.Errorf("campos alterados: %v", e.CamposAlterados)
		}
		if _, err := svc.InativarUnidade(ctx, c.adm, u.ID); err != nil {
			t.Fatal(err)
		}
		e = conferir(t, "inativação de unidade", eventosDe(t, c, "unidade", "alteracao", &u.ID, nil), 2)
		antesDepois(t, "inativação de unidade", e, "situacao", "ativa", "inativa")
		// O turno também é auditado.
		conferir(t, "turno", eventosDe(t, c, "turno_unidade", "inclusao", &u.Turnos[0].ID, nil), 1)
	})

	t.Run("pessoa", func(t *testing.T) {
		p, err := svc.CriarPessoa(ctx, c.adm, adm.EntradaPessoa{Nome: "Joana Auditada"})
		if err != nil {
			t.Fatal(err)
		}
		e := conferir(t, "inclusão de pessoa", eventosDe(t, c, "pessoa", "inclusao", &p.ID, nil), 1)
		antesDepois(t, "inclusão de pessoa", e, "nome", nil, "Joana Auditada")
		if _, err := svc.EditarPessoa(ctx, c.adm, p.ID, adm.EntradaPessoa{Nome: "Joana Auditada", Cargo: ptr("Mecânica")}); err != nil {
			t.Fatal(err)
		}
		e = conferir(t, "alteração de pessoa", eventosDe(t, c, "pessoa", "alteracao", &p.ID, nil), 1)
		antesDepois(t, "alteração de pessoa", e, "cargo", nil, "Mecânica")
		if !slices.Equal(e.CamposAlterados, []string{"cargo"}) {
			t.Errorf("só o cargo mudou: %v", e.CamposAlterados)
		}
		if _, err := svc.InativarPessoa(ctx, c.adm, p.ID); err != nil {
			t.Fatal(err)
		}
		e = conferir(t, "inativação de pessoa", eventosDe(t, c, "pessoa", "alteracao", &p.ID, nil), 2)
		antesDepois(t, "inativação de pessoa", e, "situacao", "ativa", "inativa")
	})

	// O usuário-alvo: técnico, criado pelo administrador.
	tec, err := svc.CriarUsuario(ctx, c.adm, adm.EntradaUsuario{Nome: "Técnico Auditado", Email: emailAleatorio("aud"), Perfis: []string{"tecnico"}})
	if err != nil {
		t.Fatal(err)
	}
	alvo := tec.Usuario.ID

	t.Run("usuário", func(t *testing.T) {
		e := conferir(t, "inclusão de usuário", eventosDe(t, c, "usuario", "inclusao", &alvo, nil), 1)
		antesDepois(t, "inclusão de usuário", e, "situacao", nil, "convidado")
		if _, err := svc.EditarUsuario(ctx, c.adm, alvo, adm.EntradaEditarUsuario{Nome: "Técnico Auditado", Telefone: ptr("11977776666")}); err != nil {
			t.Fatal(err)
		}
		e = conferir(t, "alteração de usuário", eventosDe(t, c, "usuario", "alteracao", &alvo, nil), 1)
		antesDepois(t, "alteração de usuário", e, "telefone", nil, "11977776666")
		if _, err := svc.InativarUsuario(ctx, c.adm, alvo); err != nil {
			t.Fatal(err)
		}
		e = conferir(t, "inativação de usuário", eventosDe(t, c, "usuario", "alteracao", &alvo, nil), 2)
		antesDepois(t, "inativação de usuário", e, "situacao", "convidado", "inativo")
		if _, err := svc.ReativarUsuario(ctx, c.adm, alvo); err != nil {
			t.Fatal(err)
		}
		// Nenhum segredo do convite vai para a trilha.
		for _, ev := range eventosDe(t, c, "usuario", "alteracao", &alvo, nil) {
			if strings.Contains(string(ev.Antes)+string(ev.Depois), "hash") {
				t.Errorf("hash na trilha do usuário: %s %s", ev.Antes, ev.Depois)
			}
		}
	})

	t.Run("perfil", func(t *testing.T) {
		// Os 12 perfis padrão nascem com a conta, gravados pelo sistema.
		pag, err := svc.ListarEventosAuditoria(ctxWeb(), c.adm, adm.FiltroAuditoria{Entidade: ptr("perfil"), Acao: ptr("inclusao"), Limite: 50})
		if err != nil || len(pag.Dados) != 12 {
			t.Fatalf("inclusão dos perfis padrão: %d %v", len(pag.Dados), err)
		}
		ev, err := svc.ObterEventoAuditoria(ctxWeb(), c.adm, pag.Dados[0].ID)
		if err != nil || ev.AtorTipo != db.AtorSistema || ev.UsuarioID != nil || campo(t, ev.Depois, "codigo") == nil ||
			campo(t, ev.Depois, "permissoes") == nil {
			t.Errorf("perfil padrão na trilha: %+v %v", ev.ResumoEvento, err)
		}

		// Perfis do usuário: o que entra e o que sai.
		perfis, err := svc.ListarPerfis(ctxWeb(), c.adm)
		if err != nil {
			t.Fatal(err)
		}
		idPerfil := map[string]string{}
		for _, p := range perfis {
			idPerfil[p.Codigo] = p.ID.String()
		}
		if _, err := svc.DefinirPerfis(ctx, c.adm, alvo, []string{"almoxarife"}); err != nil {
			t.Fatal(err)
		}
		entradas := eventosDe(t, c, "usuario_perfil", "inclusao", nil, &alvo)
		e := conferir(t, "perfil que entra", entradas, 2)
		antesDepois(t, "perfil que entra", e, "perfil_id", nil, idPerfil["almoxarife"])
		// entradas[1] é o tecnico, colocado na criação; é ele que sai.
		antesDepois(t, "perfil da criação", entradas[1], "perfil_id", nil, idPerfil["tecnico"])
		conferirRetirada(t, "perfil que sai", entradas[1])
	})

	t.Run("permissão avulsa", func(t *testing.T) {
		if _, err := svc.ConcederPermissao(ctx, c.adm, alvo, autorizacao.OsEncerrar); err != nil {
			t.Fatal(err)
		}
		e := conferir(t, "concessão", eventosDe(t, c, "usuario_permissao", "inclusao", nil, &alvo), 1)
		antesDepois(t, "concessão", e, "permissao", nil, autorizacao.OsEncerrar)
		if _, err := svc.RetirarPermissao(ctx, c.adm, alvo, autorizacao.OsEncerrar); err != nil {
			t.Fatal(err)
		}
		conferirRetirada(t, "retirada", e)
	})

	t.Run("operação recusada não deixa evento", func(t *testing.T) {
		antes := len(eventosDe(t, c, "usuario_permissao", "inclusao", nil, &alvo))
		_, err := svc.ConcederPermissao(ctx, c.adm, alvo, "nao.existe")
		exigirErro(t, err, erros.EntradaInvalida)
		if depois := len(eventosDe(t, c, "usuario_permissao", "inclusao", nil, &alvo)); depois != antes {
			t.Errorf("evento de concessão recusada: %d -> %d", antes, depois)
		}
	})
}
