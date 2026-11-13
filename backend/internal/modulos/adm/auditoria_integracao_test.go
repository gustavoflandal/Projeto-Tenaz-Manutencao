package adm_test

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"tenaz/internal/modulos/adm"
	"tenaz/internal/plataforma/auditoria"
	"tenaz/internal/plataforma/autenticacao"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/db/dbteste"
	"tenaz/internal/plataforma/erros"
)

// Trilha de auditoria configurável (RF-ADM-011, ampliado na etapa E01): um
// teste por tipo de operação, captura por tela, filtros e CSV, verificação de
// integridade, permissão e isolamento. Os critérios de aceite citados são os
// do requisito.

// ctxCom é o contexto de uma requisição com IP e tela de origem (como o
// middleware o montaria a partir de X-Tenaz-Tela).
func ctxCom(ip, tela string) context.Context {
	return db.ComRequisicao(context.Background(), db.Requisicao{
		Origem: db.OrigemWeb, IP: ip, UserAgent: "teste-auditoria", RequestID: uuid.NewString(), Tela: tela,
	})
}

func ipAleatorio() string {
	return fmt.Sprintf("198.51.100.%d", 1+uuid.New()[0]%250)
}

func ptrTexto(p *string) string {
	if p == nil {
		return "<nulo>"
	}
	return *p
}

// eventosFiltrados lê todas as páginas da consulta, em páginas pequenas, e
// devolve só o resumo.
func eventosFiltrados(t *testing.T, c conta, f adm.FiltroAuditoria) []adm.ResumoEvento {
	t.Helper()
	f.Limite = 3
	var out []adm.ResumoEvento
	for {
		pag, err := svc.ListarEventosAuditoria(ctxWeb(), c.adm, f)
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, pag.Dados...)
		if pag.ProximoCursor == "" {
			return out
		}
		f.Cursor = pag.ProximoCursor
	}
}

func primeiroEvento(t *testing.T, c conta, f adm.FiltroAuditoria) adm.Evento {
	t.Helper()
	l := eventosFiltrados(t, c, f)
	if len(l) == 0 {
		t.Fatalf("nenhum evento para o filtro %+v", f)
	}
	ev, err := svc.ObterEventoAuditoria(ctxWeb(), c.adm, l[0].ID)
	if err != nil {
		t.Fatal(err)
	}
	return ev
}

type esperado struct {
	acao    string
	usuario autorizacao.Ator
	nome    string
	ip      string
	modulo  string
	tela    *string
}

func conferirEvento(t *testing.T, nome string, e adm.Evento, q esperado, de, ate time.Time) {
	t.Helper()
	if e.OcorridoEm.IsZero() || e.OcorridoEm.Before(de) || e.OcorridoEm.After(ate) {
		t.Errorf("%s: data e hora %v fora de [%v, %v]", nome, e.OcorridoEm, de, ate)
	}
	if e.UsuarioID == nil || *e.UsuarioID != q.usuario.UsuarioID {
		t.Errorf("%s: usuário %v, esperado %s", nome, e.UsuarioID, q.usuario.UsuarioID)
	}
	if e.UsuarioNome == nil || *e.UsuarioNome != q.nome {
		t.Errorf("%s: nome do usuário %s, esperado %s", nome, ptrTexto(e.UsuarioNome), q.nome)
	}
	if e.IP == nil || *e.IP != q.ip {
		t.Errorf("%s: IP %s, esperado %s", nome, ptrTexto(e.IP), q.ip)
	}
	if e.Acao != q.acao {
		t.Errorf("%s: tipo %s, esperado %s", nome, e.Acao, q.acao)
	}
	if e.Modulo == nil || *e.Modulo != q.modulo {
		t.Errorf("%s: módulo %s, esperado %s", nome, ptrTexto(e.Modulo), q.modulo)
	}
	if ptrTexto(e.Tela) != ptrTexto(q.tela) {
		t.Errorf("%s: tela %s, esperada %s", nome, ptrTexto(e.Tela), ptrTexto(q.tela))
	}
	if len(e.Hash) != 32 || bytes.Equal(e.Hash, make([]byte, 32)) {
		t.Errorf("%s: hash %x", nome, e.Hash)
	}
	if e.Resultado != "sucesso" {
		t.Errorf("%s: resultado %s", nome, e.Resultado)
	}
}

// Critério 1: para cada tipo de operação, o evento traz data e hora, usuário
// (id e nome), IP, tipo, tela ou módulo e hash; e o hash do banco é o do
// formato 1 recalculado no Go.
func TestEventoPorTipoDeOperacao(t *testing.T) {
	c := novaConta(t, "profissional")
	ip := ipAleatorio()
	const nome = "Ana Administradora"
	empresa := c.criada.EmpresaID
	inicio := relogioDoBanco(t).Add(-time.Second)
	telaUnidades, telaPessoas, telaAuditoria := ptr("adm.unidades"), ptr("adm.pessoas"), ptr("adm.auditoria")
	buscar := func(acao, entidade string, tela *string, id *uuid.UUID) adm.Evento {
		t.Helper()
		f := adm.FiltroAuditoria{Acao: ptr(acao), Entidade: ptr(entidade), IP: ptr(ip), De: &inicio, EntidadeID: id}
		if tela != nil {
			f.Tela = tela
		}
		return primeiroEvento(t, c, f)
	}
	fim := func() time.Time { return relogioDoBanco(t).Add(time.Second) }

	// login e logout (gravados pelo papel da autenticação, na trilha da empresa)
	r, err := auth.Login(ctxCom(ip, ""), autenticacao.EntradaLogin{Email: c.email, Senha: senhaBoa, EmpresaID: &empresa})
	if err != nil || r.Sessao == nil {
		t.Fatalf("login: %v", err)
	}
	conferirEvento(t, "login", buscar("login", "identidade", nil, nil),
		esperado{acao: "login", usuario: c.adm, nome: nome, ip: ip, modulo: "adm"}, inicio, fim())
	sessao, err := auth.AutenticarWeb(ctxCom(ip, ""), r.Sessao.Cookie, r.Sessao.CSRF, true)
	if err != nil {
		t.Fatal(err)
	}
	if err := auth.Logout(ctxCom(ip, ""), sessao); err != nil {
		t.Fatal(err)
	}
	conferirEvento(t, "logout", buscar("logout", "identidade", nil, nil),
		esperado{acao: "logout", usuario: c.adm, nome: nome, ip: ip, modulo: "adm"}, inicio, fim())

	// inclusão e alteração (trg_auditoria), com a tela de origem
	u, err := svc.CriarUnidade(ctxCom(ip, "adm.unidades"), c.adm, adm.EntradaUnidade{Codigo: "EVT", Nome: "Filial dos eventos", Turnos: turnosComerciais()})
	if err != nil {
		t.Fatal(err)
	}
	conferirEvento(t, "inclusao", buscar("inclusao", "unidade", telaUnidades, &u.ID),
		esperado{acao: "inclusao", usuario: c.adm, nome: nome, ip: ip, modulo: "adm", tela: telaUnidades}, inicio, fim())
	if _, err := svc.EditarUnidade(ctxCom(ip, "adm.unidades"), c.adm, u.ID, adm.EntradaUnidade{Codigo: "EVT", Nome: "Filial renomeada"}); err != nil {
		t.Fatal(err)
	}
	conferirEvento(t, "alteracao", buscar("alteracao", "unidade", telaUnidades, &u.ID),
		esperado{acao: "alteracao", usuario: c.adm, nome: nome, ip: ip, modulo: "adm", tela: telaUnidades}, inicio, fim())

	// exclusão (lógica): a retirada de uma permissão avulsa
	alvo, _ := convidarEEntrar(t, c, adm.EntradaUsuario{Perfis: []string{"tecnico"}})
	if _, err := svc.ConcederPermissao(ctxCom(ip, "adm.usuarios"), c.adm, alvo.ID, autorizacao.OsEncerrar); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.RetirarPermissao(ctxCom(ip, "adm.usuarios"), c.adm, alvo.ID, autorizacao.OsEncerrar); err != nil {
		t.Fatal(err)
	}
	conferirEvento(t, "exclusao", buscar("exclusao", "usuario_permissao", ptr("adm.usuarios"), nil),
		esperado{acao: "exclusao", usuario: c.adm, nome: nome, ip: ip, modulo: "adm", tela: ptr("adm.usuarios")}, inicio, fim())

	// deleção (física): o trigger genérico grava a linha inteira em antes
	f, err := svc.IncluirFeriado(ctxCom(ip, "adm.unidades"), c.adm, u.ID, adm.EntradaFeriado{Data: "2026-12-25", Descricao: "Natal"})
	if err != nil {
		t.Fatal(err)
	}
	err = amb.App.NaEmpresa(db.ComModulo(ctxCom(ip, "adm.unidades"), "adm"), empresa, c.adm.Escopo, func(ctx context.Context) error {
		tx, err := db.Tx(ctx)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx, `DELETE FROM app.feriado_unidade WHERE empresa_id = $1 AND id = $2`, empresa, f.ID)
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	del := buscar("delecao", "feriado_unidade", telaUnidades, &f.ID)
	conferirEvento(t, "delecao", del,
		esperado{acao: "delecao", usuario: c.adm, nome: nome, ip: ip, modulo: "adm", tela: telaUnidades}, inicio, fim())
	if campo(t, del.Antes, "descricao") != "Natal" || len(del.Depois) != 0 {
		t.Errorf("deleção: antes %s, depois %s", del.Antes, del.Depois)
	}

	// visualização: o CPF completo é sempre gravado
	cpf := cpfAleatorio()
	p, err := svc.CriarPessoa(ctxCom(ip, "adm.pessoas"), c.adm, adm.EntradaPessoa{Nome: "Carlos Silva", IdentificacaoTipo: ptr("cpf"), Identificacao: &cpf})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.ObterPessoa(ctxCom(ip, "adm.pessoas"), c.adm, p.ID); err != nil {
		t.Fatal(err)
	}
	vis := buscar("visualizacao", "pessoa", telaPessoas, &p.ID)
	conferirEvento(t, "visualizacao", vis,
		esperado{acao: "visualizacao", usuario: c.adm, nome: nome, ip: ip, modulo: "adm", tela: telaPessoas}, inicio, fim())
	if !vis.Sensivel {
		t.Error("a visualização do CPF completo é sensível")
	}

	// impressão: segue a captura da tela (padrão: ligada)
	err = amb.App.NaEmpresa(db.ComModulo(ctxCom(ip, "adm.pessoas"), "adm"), empresa, c.adm.Escopo, func(ctx context.Context) error {
		gravou, err := aud.Imprimiu(ctx, auditoria.Operacao{EmpresaID: empresa, UsuarioID: ptr(c.adm.UsuarioID), Entidade: "pessoa",
			EntidadeID: &p.ID, Detalhe: map[string]any{"formato": "pdf"}})
		if err == nil && !gravou {
			t.Error("a impressão deveria ter sido gravada (captura padrão ligada)")
		}
		return err
	})
	if err != nil {
		t.Fatal(err)
	}
	imp := buscar("impressao", "pessoa", telaPessoas, &p.ID)
	conferirEvento(t, "impressao", imp,
		esperado{acao: "impressao", usuario: c.adm, nome: nome, ip: ip, modulo: "adm", tela: telaPessoas}, inicio, fim())
	if campo(t, imp.Detalhe, "formato") != "pdf" {
		t.Errorf("detalhe da impressão: %s", imp.Detalhe)
	}

	// extração: a exportação da própria trilha
	ate := relogioDoBanco(t).Add(time.Minute)
	x, err := svc.ExportarEventosAuditoria(ctxCom(ip, "adm.auditoria"), c.adm, adm.FiltroAuditoria{De: &inicio, Ate: &ate})
	if err != nil || x.Quantidade < 8 {
		t.Fatalf("exportação: %d linhas, %v", x.Quantidade, err)
	}
	ext := buscar("extracao", "evento_auditoria", telaAuditoria, nil)
	conferirEvento(t, "extracao", ext,
		esperado{acao: "extracao", usuario: c.adm, nome: nome, ip: ip, modulo: "adm", tela: telaAuditoria}, inicio, fim())
	if campo(t, ext.Detalhe, "formato") != "csv" || campo(t, ext.Detalhe, "quantidade") != float64(x.Quantidade) {
		t.Errorf("detalhe da extração: %s", ext.Detalhe)
	}

	// o hash gravado pelo banco é o do formato 1 recalculado no Go, em todos
	// os eventos da empresa (inclusive os do trigger, sem tela, com antes e depois)
	n, div, err := dbteste.ConferirHashesDaTrilha(context.Background(), amb.Migrador, empresa)
	if err != nil || n < 20 {
		t.Fatalf("hashes conferidos: %d, %v", n, err)
	}
	for _, d := range div {
		t.Error(d)
	}
}

// Critério 8: o evento da visualização do CPF completo não contém o CPF.
func TestVisualizacaoDoCPFNaoGravaOCPF(t *testing.T) {
	c := novaConta(t, "profissional")
	cpf := cpfAleatorio()
	p, err := svc.CriarPessoa(ctxWeb(), c.adm, adm.EntradaPessoa{Nome: "Maria Souza", IdentificacaoTipo: ptr("cpf"), Identificacao: &cpf})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.ObterPessoa(ctxWeb(), c.adm, p.ID); err != nil {
		t.Fatal(err)
	}
	var n int
	// o CPF da criação aparece só no evento de inclusão; nenhuma linha de visualização o contém
	if err := amb.Migrador.QueryRow(context.Background(), `SELECT count(*) FROM app.evento_auditoria
WHERE empresa_id = $1 AND acao = 'visualizacao' AND entidade_id = $2
  AND (coalesce(antes::text, '') || coalesce(depois::text, '') || coalesce(detalhe::text, '') || coalesce(entidade_chave, '') ||
       coalesce(user_agent, '') || coalesce(request_id, '')) LIKE '%' || $3 || '%'`, c.criada.EmpresaID, p.ID, cpf).Scan(&n); err != nil || n != 0 {
		t.Fatalf("o CPF apareceu em evento de visualização: %d %v", n, err)
	}
	if err := amb.Migrador.QueryRow(context.Background(), `SELECT count(*) FROM app.evento_auditoria
WHERE empresa_id = $1 AND acao = 'visualizacao' AND entidade_id = $2 AND sensivel`, c.criada.EmpresaID, p.ID).Scan(&n); err != nil || n != 1 {
		t.Fatalf("visualização sensível gravada: %d %v", n, err)
	}
}

// Critérios 2, 3 e 4: o que sempre grava, a captura de visualização por tela
// e a de impressão e extração, inclusive na exportação da própria trilha.
func TestCapturaPorTela(t *testing.T) {
	c := novaConta(t, "profissional")
	empresa := c.criada.EmpresaID
	ctx := ctxWeb()
	todas := auditoria.Telas()

	caps, err := svc.ListarCapturasAuditoria(ctx, c.adm)
	if err != nil || len(caps) != len(todas) {
		t.Fatalf("lista de capturas: %d telas, %v", len(caps), err)
	}
	for _, cp := range caps {
		if cp.Visualizacao || !cp.Impressao || !cp.Extracao || cp.Configurada {
			t.Fatalf("padrão da criação (impressão e extração ligadas, visualização desligada): %+v", cp)
		}
	}
	if caps[0].Tela != "adm.empresa" || caps[0].Modulo != "adm" || caps[0].Nome == "" {
		t.Errorf("catálogo: %+v", caps[0])
	}

	contar := func(acao, entidade string, id *uuid.UUID) int {
		t.Helper()
		return len(eventosFiltrados(t, c, adm.FiltroAuditoria{Acao: ptr(acao), Entidade: ptr(entidade), EntidadeID: id}))
	}
	u1, err := svc.CriarUnidade(ctx, c.adm, adm.EntradaUnidade{Codigo: "U1", Nome: "Unidade 1", Turnos: turnosComerciais()})
	if err != nil {
		t.Fatal(err)
	}
	u2, err := svc.CriarUnidade(ctx, c.adm, adm.EntradaUnidade{Codigo: "U2", Nome: "Unidade 2", Turnos: turnosComerciais()})
	if err != nil {
		t.Fatal(err)
	}
	const tela = "adm.unidade_detalhe"

	// visualização desligada (padrão): abrir o registro não gera evento
	if _, err := svc.ObterUnidade(ctxCom(ipAleatorio(), tela), c.adm, u1.ID); err != nil {
		t.Fatal(err)
	}
	if n := contar("visualizacao", "unidade", &u1.ID); n != 0 {
		t.Fatalf("visualização desligada gravou %d eventos", n)
	}

	// ligada: um evento por registro aberto
	r, err := svc.AlterarCapturaAuditoria(ctx, c.adm, tela, adm.EntradaCaptura{Visualizacao: true, Impressao: true, Extracao: true})
	if err != nil || !r.Visualizacao || !r.Configurada || r.Tela != tela {
		t.Fatalf("ligar a visualização: %+v %v", r, err)
	}
	for _, id := range []uuid.UUID{u1.ID, u1.ID, u2.ID} {
		if _, err := svc.ObterUnidade(ctxCom(ipAleatorio(), tela), c.adm, id); err != nil {
			t.Fatal(err)
		}
	}
	if n1, n2 := contar("visualizacao", "unidade", &u1.ID), contar("visualizacao", "unidade", &u2.ID); n1 != 2 || n2 != 1 {
		t.Fatalf("um evento por registro aberto: u1 %d (esperados 2), u2 %d (esperado 1)", n1, n2)
	}
	ev := primeiroEvento(t, c, adm.FiltroAuditoria{Acao: ptr("visualizacao"), EntidadeID: &u2.ID})
	if ev.Sensivel || ptrTexto(ev.Tela) != tela || ptrTexto(ev.Modulo) != "adm" || ev.UsuarioNome == nil {
		t.Errorf("evento de visualização: %+v", ev.ResumoEvento)
	}
	// sem tela (API) ou com a tela de outro módulo, vale o padrão: nada é gravado
	if _, err := svc.ObterUnidade(ctxWeb(), c.adm, u2.ID); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.ObterUnidade(ctxCom(ipAleatorio(), "ind.painel"), c.adm, u2.ID); err != nil {
		t.Fatal(err)
	}
	if n := contar("visualizacao", "unidade", &u2.ID); n != 1 {
		t.Fatalf("sem tela ou com tela de outro módulo não grava: %d eventos", n)
	}

	// a troca da chave aparece na trilha como inclusão e, depois, alteração
	inc := eventosFiltrados(t, c, adm.FiltroAuditoria{Entidade: ptr("captura_auditoria"), Acao: ptr("inclusao")})
	if len(inc) != 1 || ptrTexto(inc[0].EntidadeChave) != tela || inc[0].UsuarioNome == nil {
		t.Fatalf("a primeira configuração é uma inclusão com a tela na chave: %+v", inc)
	}
	if _, err := svc.AlterarCapturaAuditoria(ctx, c.adm, tela, adm.EntradaCaptura{Visualizacao: false, Impressao: true, Extracao: true}); err != nil {
		t.Fatal(err)
	}
	alt := primeiroEvento(t, c, adm.FiltroAuditoria{Entidade: ptr("captura_auditoria"), Acao: ptr("alteracao")})
	if ptrTexto(alt.EntidadeChave) != tela || !slices.Equal(alt.CamposAlterados, []string{"visualizacao"}) ||
		campo(t, alt.Antes, "visualizacao") != true || campo(t, alt.Depois, "visualizacao") != false {
		t.Errorf("a troca da chave: %+v antes %s depois %s", alt.ResumoEvento, alt.Antes, alt.Depois)
	}
	// desligada de novo: os 3 eventos antigos ficam, e o novo não é gravado
	if _, err := svc.ObterUnidade(ctxCom(ipAleatorio(), tela), c.adm, u2.ID); err != nil {
		t.Fatal(err)
	}
	if n := contar("visualizacao", "unidade", &u2.ID); n != 1 {
		t.Fatalf("depois de desligar, não grava: %d", n)
	}

	// Critério 4: impressão e extração seguem a chave da tela, inclusive na
	// exportação da própria trilha.
	exportar := func() int {
		t.Helper()
		de, ate := relogioDoBanco(t).Add(-time.Hour), relogioDoBanco(t).Add(time.Hour)
		before := contar("extracao", "evento_auditoria", nil)
		if _, err := svc.ExportarEventosAuditoria(ctxCom(ipAleatorio(), "adm.auditoria"), c.adm, adm.FiltroAuditoria{De: &de, Ate: &ate}); err != nil {
			t.Fatal(err)
		}
		return contar("extracao", "evento_auditoria", nil) - before
	}
	if n := exportar(); n != 1 {
		t.Fatalf("extração ligada (padrão): %d eventos", n)
	}
	if _, err := svc.AlterarCapturaAuditoria(ctx, c.adm, "adm.auditoria", adm.EntradaCaptura{Visualizacao: false, Impressao: true, Extracao: false}); err != nil {
		t.Fatal(err)
	}
	if n := exportar(); n != 0 {
		t.Fatalf("extração desligada na tela de auditoria gravou %d eventos", n)
	}
	imprimir := func(telaImp string) bool {
		t.Helper()
		var gravou bool
		err := amb.App.NaEmpresa(db.ComModulo(ctxCom(ipAleatorio(), telaImp), "adm"), empresa, c.adm.Escopo, func(ctx context.Context) error {
			var err error
			gravou, err = aud.Imprimiu(ctx, auditoria.Operacao{EmpresaID: empresa, UsuarioID: ptr(c.adm.UsuarioID), Entidade: "unidade", EntidadeID: &u1.ID})
			return err
		})
		if err != nil {
			t.Fatal(err)
		}
		return gravou
	}
	if !imprimir("adm.auditoria") {
		t.Error("impressão ligada na tela de auditoria deveria gravar")
	}
	if _, err := svc.AlterarCapturaAuditoria(ctx, c.adm, "adm.auditoria", adm.EntradaCaptura{}); err != nil {
		t.Fatal(err)
	}
	if imprimir("adm.auditoria") {
		t.Error("impressão desligada não deveria gravar")
	}
	if !imprimir("adm.unidades") {
		t.Error("a chave de uma tela não vale para outra: adm.unidades continua no padrão")
	}

	// Critério 2: com todas as capturas desligadas, o que sempre grava continua
	// gravando. A visualização do CPF completo é sensível e não tem chave.
	for _, tl := range todas {
		if _, err := svc.AlterarCapturaAuditoria(ctx, c.adm, tl.Codigo, adm.EntradaCaptura{}); err != nil {
			t.Fatal(err)
		}
	}
	caps, err = svc.ListarCapturasAuditoria(ctx, c.adm)
	if err != nil {
		t.Fatal(err)
	}
	for _, cp := range caps {
		if cp.Visualizacao || cp.Impressao || cp.Extracao || !cp.Configurada {
			t.Fatalf("tudo desligado: %+v", cp)
		}
	}
	inicio := relogioDoBanco(t).Add(-time.Second)
	cpf := cpfAleatorio()
	p, err := svc.CriarPessoa(ctxCom(ipAleatorio(), "adm.pessoas"), c.adm, adm.EntradaPessoa{Nome: "Sem Captura", IdentificacaoTipo: ptr("cpf"), Identificacao: &cpf})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.EditarPessoa(ctxCom(ipAleatorio(), "adm.pessoas"), c.adm, p.ID, adm.EntradaPessoa{Nome: "Sem Captura 2", IdentificacaoTipo: ptr("cpf"), Identificacao: &cpf}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.ObterPessoa(ctxCom(ipAleatorio(), "adm.pessoa_detalhe"), c.adm, p.ID); err != nil {
		t.Fatal(err)
	}
	alvo, _ := convidarEEntrar(t, c, adm.EntradaUsuario{Perfis: []string{"tecnico"}})
	_, _ = svc.ConcederPermissao(ctx, c.adm, alvo.ID, autorizacao.OsEncerrar)
	_, _ = svc.RetirarPermissao(ctx, c.adm, alvo.ID, autorizacao.OsEncerrar)
	if _, err := auth.Login(ctxWeb(), autenticacao.EntradaLogin{Email: c.email, Senha: senhaBoa, EmpresaID: &empresa}); err != nil {
		t.Fatal(err)
	}
	for _, q := range []struct {
		acao, entidade string
		id             *uuid.UUID
	}{
		{"inclusao", "pessoa", &p.ID}, {"alteracao", "pessoa", &p.ID}, {"visualizacao", "pessoa", &p.ID},
		{"exclusao", "usuario_permissao", nil}, {"login", "identidade", nil},
	} {
		l := eventosFiltrados(t, c, adm.FiltroAuditoria{Acao: ptr(q.acao), Entidade: ptr(q.entidade), EntidadeID: q.id, De: &inicio})
		if len(l) == 0 {
			t.Errorf("com todas as capturas desligadas, %s de %s não foi gravado", q.acao, q.entidade)
		}
	}
	// uma visualização comum, sem dado sensível, não grava
	if _, err := svc.ObterUnidade(ctxCom(ipAleatorio(), tela), c.adm, u1.ID); err != nil {
		t.Fatal(err)
	}
	if n := contar("visualizacao", "unidade", &u1.ID); n != 2 {
		t.Errorf("visualização comum com a chave desligada: %d eventos (esperados os 2 de antes)", n)
	}

	// tela fora do catálogo é recusada
	_, err = svc.AlterarCapturaAuditoria(ctx, c.adm, "adm.inventada", adm.EntradaCaptura{})
	exigirErro(t, err, erros.EntradaInvalida)
	_, err = svc.AlterarCapturaAuditoria(ctx, c.adm, "xxx.painel", adm.EntradaCaptura{})
	exigirErro(t, err, erros.EntradaInvalida)
}

// A requisição de escopo externo usa o cache da configuração da empresa,
// carregada num contexto interno: o prestador não lê captura_auditoria, e
// mesmo assim a captura de visualização da tela vale para ele.
func TestCapturaNoEscopoExternoEMCache(t *testing.T) {
	c := novaConta(t, "profissional")
	empresa := c.criada.EmpresaID
	forn := uuid.New()
	prest, aPrest := convidarEEntrar(t, c, adm.EntradaUsuario{Perfis: []string{"prestador"}, FornecedorID: &forn})
	const tela = "adm.perfis"
	visualizar := func() bool {
		t.Helper()
		var gravou bool
		err := amb.App.NaEmpresa(db.ComModulo(ctxCom(ipAleatorio(), tela), "adm"), empresa, aPrest.Escopo, func(ctx context.Context) error {
			var err error
			gravou, err = aud.Visualizou(ctx, auditoria.Operacao{EmpresaID: empresa, UsuarioID: ptr(aPrest.UsuarioID), Entidade: "perfil", EntidadeID: ptr(uuid.New())})
			return err
		})
		if err != nil {
			t.Fatal(err)
		}
		return gravou
	}
	if visualizar() {
		t.Fatal("visualização desligada (padrão): o prestador não grava")
	}
	if _, err := svc.AlterarCapturaAuditoria(ctxWeb(), c.adm, tela, adm.EntradaCaptura{Visualizacao: true, Impressao: true, Extracao: true}); err != nil {
		t.Fatal(err)
	}
	if !visualizar() {
		t.Fatal("o prestador (escopo externo) tem de enxergar a captura ligada pelo administrador")
	}
	ev := primeiroEvento(t, c, adm.FiltroAuditoria{Acao: ptr("visualizacao"), UsuarioID: &prest.ID})
	if ev.AtorTipo != db.AtorUsuario || ev.UsuarioNome == nil || *ev.UsuarioNome != prest.Nome {
		t.Errorf("evento do prestador: %+v", ev.ResumoEvento)
	}
	var escopo *string
	if err := amb.Migrador.QueryRow(context.Background(), `SELECT escopo FROM app.evento_auditoria WHERE id = $1`, ev.ID).Scan(&escopo); err != nil || ptrTexto(escopo) != "prestador" {
		t.Errorf("escopo do evento: %v %v", escopo, err)
	}

	// Cache: mudança feita direto no banco só é vista depois da validade ou da
	// invalidação; mudança pelo serviço vale na hora.
	if _, err := amb.Migrador.Exec(context.Background(), `UPDATE app.captura_auditoria SET visualizacao = false WHERE empresa_id = $1 AND tela = $2`, empresa, tela); err != nil {
		t.Fatal(err)
	}
	if !visualizar() {
		t.Error("o cache deveria segurar o valor antigo até a invalidação")
	}
	aud.InvalidarCapturas(empresa)
	if visualizar() {
		t.Error("depois de invalidar, vale o banco")
	}
}

// Critério 6: os filtros por usuário, tela ou módulo, tipo, período e IP
// devolvem só os eventos que atendem a todos, e a exportação em CSV traz as
// mesmas linhas, na mesma ordem.
func TestFiltrosEExportacaoCSV(t *testing.T) {
	c := novaConta(t, "profissional")
	empresa := c.criada.EmpresaID
	adm2, a2 := convidarEEntrar(t, c, adm.EntradaUsuario{Nome: "Segundo Admin", Perfis: []string{"administrador"}})
	ip1, ip2 := "198.51.100.11", "198.51.100.22"
	inicio := relogioDoBanco(t).Add(-time.Second)

	u, err := svc.CriarUnidade(ctxCom(ip1, "adm.unidades"), c.adm, adm.EntradaUnidade{Codigo: "F1", Nome: "Filial F1", Turnos: turnosComerciais()})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.EditarUnidade(ctxCom(ip2, "adm.unidades"), a2, u.ID, adm.EntradaUnidade{Codigo: "F1", Nome: "Filial F1 B"}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.EditarUnidade(ctxCom(ip2, "adm.empresa"), a2, u.ID, adm.EntradaUnidade{Codigo: "F1", Nome: "Filial =F1 C"}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.EditarEu(ctxCom(ip1, ""), a2, adm.EntradaEu{Nome: "=HYPERLINK(\"http://x\")"}); err != nil {
		t.Fatal(err)
	}
	// depois da troca de nome, o evento carrega o nome novo (o antigo ficou nos anteriores)
	if _, err := svc.EditarUnidade(ctxCom(ip1, "adm.unidades"), a2, u.ID, adm.EntradaUnidade{Codigo: "F1", Nome: "Filial F1 D"}); err != nil {
		t.Fatal(err)
	}
	cpf := cpfAleatorio()
	p, err := svc.CriarPessoa(ctxCom(ip1, "adm.pessoas"), c.adm, adm.EntradaPessoa{Nome: "Pessoa do filtro", IdentificacaoTipo: ptr("cpf"), Identificacao: &cpf})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.ObterPessoa(ctxCom(ip2, "adm.pessoas"), a2, p.ID); err != nil {
		t.Fatal(err)
	}
	fim := relogioDoBanco(t).Add(time.Second)

	// o oráculo: a mesma pergunta em SQL, direto no banco, em ordem decrescente
	esperados := func(usuario *uuid.UUID, modulo, tela, acao, ip *string, de, ate *time.Time) []uuid.UUID {
		t.Helper()
		rows, err := amb.Migrador.Query(context.Background(), `SELECT id FROM app.evento_auditoria
WHERE empresa_id = $1 AND ($2::uuid IS NULL OR usuario_id = $2) AND ($3::text IS NULL OR modulo = $3)
  AND ($4::text IS NULL OR tela = $4) AND ($5::text IS NULL OR acao = $5) AND ($6::text IS NULL OR host(ip) = $6)
  AND ($7::timestamptz IS NULL OR ocorrido_em >= $7) AND ($8::timestamptz IS NULL OR ocorrido_em < $8)
ORDER BY ocorrido_em DESC, id DESC`, empresa, usuario, modulo, tela, acao, ip, de, ate)
		if err != nil {
			t.Fatal(err)
		}
		defer rows.Close()
		var ids []uuid.UUID
		for rows.Next() {
			var id uuid.UUID
			if err := rows.Scan(&id); err != nil {
				t.Fatal(err)
			}
			ids = append(ids, id)
		}
		return ids
	}
	idsDe := func(l []adm.ResumoEvento) []uuid.UUID {
		out := make([]uuid.UUID, len(l))
		for i, e := range l {
			out[i] = e.ID
		}
		return out
	}
	idsDoCSV := func(b []byte) []uuid.UUID {
		t.Helper()
		linhas, err := csv.NewReader(bytes.NewReader(b)).ReadAll()
		if err != nil {
			t.Fatal(err)
		}
		if strings.Join(linhas[0], ",") != "id,ocorrido_em,usuario_id,usuario_nome,ator_tipo,ip,acao,resultado,modulo,tela,entidade,entidade_id,entidade_chave,campos_alterados,sensivel,origem,request_id,hash" {
			t.Fatalf("cabeçalho do CSV: %v", linhas[0])
		}
		var ids []uuid.UUID
		for _, l := range linhas[1:] {
			ids = append(ids, uuid.MustParse(l[0]))
		}
		return ids
	}
	casos := []struct {
		nome               string
		f                  adm.FiltroAuditoria
		usuario            *uuid.UUID
		modulo, tela, acao *string
		ip                 *string
		quer               int // mínimo de eventos que o caso precisa ter
	}{
		{nome: "só o período", f: adm.FiltroAuditoria{}, quer: 8},
		{nome: "usuário", f: adm.FiltroAuditoria{UsuarioID: &adm2.ID}, usuario: &adm2.ID, quer: 3},
		{nome: "módulo", f: adm.FiltroAuditoria{Modulo: ptr("adm")}, modulo: ptr("adm"), quer: 8},
		{nome: "módulo sem eventos", f: adm.FiltroAuditoria{Modulo: ptr("os")}, modulo: ptr("os"), quer: 0},
		{nome: "tela", f: adm.FiltroAuditoria{Tela: ptr("adm.unidades")}, tela: ptr("adm.unidades"), quer: 2},
		{nome: "tipo", f: adm.FiltroAuditoria{Acao: ptr("alteracao")}, acao: ptr("alteracao"), quer: 3},
		{nome: "IP", f: adm.FiltroAuditoria{IP: ptr(ip2)}, ip: ptr(ip2), quer: 3},
		{nome: "usuário e tela", f: adm.FiltroAuditoria{UsuarioID: &adm2.ID, Tela: ptr("adm.unidades")}, usuario: &adm2.ID, tela: ptr("adm.unidades"), quer: 1},
		{nome: "tipo, IP e usuário", f: adm.FiltroAuditoria{Acao: ptr("alteracao"), IP: ptr(ip2), UsuarioID: &adm2.ID}, usuario: &adm2.ID, acao: ptr("alteracao"), ip: ptr(ip2), quer: 2},
		{nome: "tudo junto", f: adm.FiltroAuditoria{UsuarioID: &adm2.ID, Modulo: ptr("adm"), Tela: ptr("adm.pessoas"), Acao: ptr("visualizacao"), IP: ptr(ip2)},
			usuario: &adm2.ID, modulo: ptr("adm"), tela: ptr("adm.pessoas"), acao: ptr("visualizacao"), ip: ptr(ip2), quer: 1},
		{nome: "combinação vazia", f: adm.FiltroAuditoria{UsuarioID: &adm2.ID, Tela: ptr("adm.pessoas"), Acao: ptr("inclusao")},
			usuario: &adm2.ID, tela: ptr("adm.pessoas"), acao: ptr("inclusao"), quer: 0},
	}
	for _, k := range casos {
		t.Run(k.nome, func(t *testing.T) {
			f := k.f
			f.De, f.Ate = &inicio, &fim
			quer := esperados(k.usuario, k.modulo, k.tela, k.acao, k.ip, &inicio, &fim)
			if len(quer) < k.quer {
				t.Fatalf("o cenário deveria ter ao menos %d eventos para este filtro, tem %d", k.quer, len(quer))
			}
			lista := idsDe(eventosFiltrados(t, c, f))
			if !slices.Equal(lista, quer) {
				t.Errorf("consulta paginada: %d eventos, o banco diz %d\n  consulta %v\n  banco    %v", len(lista), len(quer), lista, quer)
			}
			x, err := svc.ExportarEventosAuditoria(ctxWeb(), c.adm, f)
			if err != nil {
				t.Fatal(err)
			}
			csvIDs := idsDoCSV(x.CSV)
			if !slices.Equal(csvIDs, quer) || x.Quantidade != len(quer) {
				t.Errorf("CSV: %d linhas (quantidade %d), o banco diz %d", len(csvIDs), x.Quantidade, len(quer))
			}
		})
	}

	// texto livre que começa por = é neutralizado na planilha
	x, err := svc.ExportarEventosAuditoria(ctxWeb(), c.adm, adm.FiltroAuditoria{UsuarioID: &adm2.ID, Acao: ptr("alteracao"), Entidade: ptr("unidade"), De: &inicio, Ate: &fim})
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(x.CSV), ",=HYPERLINK") || !strings.Contains(string(x.CSV), `'=HYPERLINK`) {
		t.Errorf("injeção de fórmula no CSV:\n%s", x.CSV)
	}

	// filtros inválidos e período
	for nome, f := range map[string]adm.FiltroAuditoria{
		"tipo desconhecido":     {Acao: ptr("login_falho")},
		"módulo desconhecido":   {Modulo: ptr("xx")},
		"tela fora do catálogo": {Tela: ptr("adm.inventada")},
		"IP malformado":         {IP: ptr("999.1.1.1"), De: &inicio},
		"IP sem período":        {IP: ptr(ip1)},
		"período invertido":     {De: &fim, Ate: &inicio},
	} {
		if _, err := svc.ListarEventosAuditoria(ctxWeb(), c.adm, f); !errorsIs(err, erros.EntradaInvalida) {
			t.Errorf("%s: %v", nome, err)
		}
	}
	if _, err := svc.ExportarEventosAuditoria(ctxWeb(), c.adm, adm.FiltroAuditoria{}); !errorsIs(err, erros.EntradaInvalida) {
		t.Errorf("exportação sem período: %v", err)
	}
	if _, err := svc.ExportarEventosAuditoria(ctxWeb(), c.adm, adm.FiltroAuditoria{De: &inicio}); !errorsIs(err, erros.EntradaInvalida) {
		t.Errorf("exportação sem fim do período: %v", err)
	}
}

func errorsIs(err error, alvo *erros.Erro) bool {
	e := comoErroOuNil(err)
	return e != nil && e.Codigo == alvo.Codigo
}

func comoErroOuNil(err error) *erros.Erro {
	for err != nil {
		if e, ok := err.(*erros.Erro); ok {
			return e
		}
		u, ok := err.(interface{ Unwrap() error })
		if !ok {
			return nil
		}
		err = u.Unwrap()
	}
	return nil
}

// Critério 7: sem adm.ler_auditoria, a consulta, o detalhe, a exportação e a
// verificação respondem 403; sem adm.gerenciar_parametros, a configuração das
// capturas também. Usuário externo não vê as rotas (404).
func TestAuditoriaExigePermissao(t *testing.T) {
	c := novaConta(t, "profissional")
	_, tec := convidarEEntrar(t, c, adm.EntradaUsuario{Perfis: []string{"tecnico"}})
	forn := uuid.New()
	_, ext := convidarEEntrar(t, c, adm.EntradaUsuario{Perfis: []string{"prestador"}, FornecedorID: &forn})
	ctx := ctxWeb()
	agora := time.Now()
	ontem := agora.AddDate(0, 0, -1)
	id := uuid.New()

	chamadas := map[string]func(a autorizacao.Ator) error{
		"listar": func(a autorizacao.Ator) error {
			_, err := svc.ListarEventosAuditoria(ctx, a, adm.FiltroAuditoria{})
			return err
		},
		"detalhe": func(a autorizacao.Ator) error { _, err := svc.ObterEventoAuditoria(ctx, a, id); return err },
		"exportar": func(a autorizacao.Ator) error {
			_, err := svc.ExportarEventosAuditoria(ctx, a, adm.FiltroAuditoria{De: &ontem, Ate: &agora})
			return err
		},
		"verificar": func(a autorizacao.Ator) error {
			_, err := svc.VerificarIntegridadeAuditoria(ctx, a, ontem, agora)
			return err
		},
		"listar capturas": func(a autorizacao.Ator) error { _, err := svc.ListarCapturasAuditoria(ctx, a); return err },
		"alterar captura": func(a autorizacao.Ator) error {
			_, err := svc.AlterarCapturaAuditoria(ctx, a, "adm.pessoas", adm.EntradaCaptura{})
			return err
		},
	}
	for nome, chamar := range chamadas {
		exigirErro(t, chamar(tec), erros.SemPermissao)
		exigirErro(t, chamar(ext), erros.NaoEncontrado)
		_ = nome
	}
	// o administrador tem as duas permissões
	if _, err := svc.ListarCapturasAuditoria(ctx, c.adm); err != nil {
		t.Errorf("administrador lista as capturas: %v", err)
	}
	if _, err := svc.ListarEventosAuditoria(ctx, c.adm, adm.FiltroAuditoria{}); err != nil {
		t.Errorf("administrador consulta a trilha: %v", err)
	}
}

// Critério 9: a trilha de uma empresa não mostra eventos de outra, a
// exportação também não, e a verificação usa só os selos da própria empresa.
func TestAuditoriaIsolaEmpresas(t *testing.T) {
	a, b := novaConta(t, "profissional"), novaConta(t, "profissional")
	ctx := ctxWeb()
	inicio := relogioDoBanco(t).Add(-time.Hour)
	fim := relogioDoBanco(t).Add(time.Hour)
	if _, err := svc.CriarUnidade(ctx, a.adm, adm.EntradaUnidade{Codigo: "ISO", Nome: "Só da A", Turnos: turnosComerciais()}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CriarUnidade(ctx, b.adm, adm.EntradaUnidade{Codigo: "ISO", Nome: "Só da B", Turnos: turnosComerciais()}); err != nil {
		t.Fatal(err)
	}
	for _, par := range []struct{ de, outra conta }{{a, b}, {b, a}} {
		l := eventosFiltrados(t, par.de, adm.FiltroAuditoria{})
		if len(l) == 0 {
			t.Fatal("trilha vazia")
		}
		ids := make([]uuid.UUID, len(l))
		for i, e := range l {
			ids[i] = e.ID
		}
		var deOutra int
		if err := amb.Migrador.QueryRow(context.Background(), `SELECT count(*) FROM app.evento_auditoria WHERE id = ANY($1) AND empresa_id <> $2`,
			ids, par.de.criada.EmpresaID).Scan(&deOutra); err != nil || deOutra != 0 {
			t.Fatalf("a trilha mostrou %d eventos de outra empresa: %v", deOutra, err)
		}
		x, err := svc.ExportarEventosAuditoria(ctx, par.de.adm, adm.FiltroAuditoria{De: &inicio, Ate: &fim})
		if err != nil {
			t.Fatal(err)
		}
		if strings.Contains(string(x.CSV), "Só da "+map[bool]string{true: "B", false: "A"}[par.de.criada.EmpresaID == a.criada.EmpresaID]) {
			t.Error("o CSV trouxe dado da outra empresa")
		}
		// evento de outra empresa pelo id: 404
		_, err = svc.ObterEventoAuditoria(ctx, par.outra.adm, l[0].ID)
		exigirErro(t, err, erros.NaoEncontrado)
	}
}

// Critério 5, pela API do módulo: período intacto é íntegro; adulterar um
// evento direto no banco, com o papel dono, faz a verificação mostrar violado
// e apontar o evento; o selo existe também no armazenamento, e a verificação
// de uma empresa não vê os selos da outra.
func TestVerificacaoDeIntegridadePelaAPI(t *testing.T) {
	ctx := context.Background()
	a, b := novaConta(t, "profissional"), novaConta(t, "profissional")
	arm := dbteste.NovaMemoria()
	ver, err := auditoria.NovoVerificador(amb.App, arm)
	if err != nil {
		t.Fatal(err)
	}
	svcMem, err := adm.Novo(adm.Dependencias{Banco: amb.App, Autenticacao: auth, Auditoria: aud, Verificador: ver})
	if err != nil {
		t.Fatal(err)
	}
	selador, err := auditoria.NovoSelador(auditoria.DependenciasSelador{Banco: amb.Job, Armazenamento: arm})
	if err != nil {
		t.Fatal(err)
	}
	d0 := time.Now().UTC().AddDate(0, 0, -3)
	dia := func(n int) time.Time { return time.Date(d0.Year(), d0.Month(), d0.Day()+n, 0, 0, 0, 0, time.UTC) }
	// eventos antigos, semeados como o migrador: dois dias com 2 eventos em A, um em B
	semear := func(c conta, d time.Time, n int) []uuid.UUID {
		var ids []uuid.UUID
		for i := 0; i < n; i++ {
			id := uuid.Must(uuid.NewV7())
			if _, err := amb.Migrador.Exec(ctx, `INSERT INTO app.evento_auditoria (empresa_id, id, ocorrido_em, ator_tipo, acao, entidade, modulo)
VALUES ($1, $2, $3, 'sistema', 'inclusao', 'unidade', 'adm')`, c.criada.EmpresaID, id, d.Add(time.Duration(8+i)*time.Hour)); err != nil {
				t.Fatal(err)
			}
			ids = append(ids, id)
		}
		return ids
	}
	// A e B têm eventos de hoje (da criação da conta) que não entram: só o passado é selado
	idsA := semear(a, dia(0), 2)
	semear(a, dia(1), 2)
	semear(b, dia(0), 1)
	// o selo de cada empresa vai para o armazenamento
	for _, c := range []conta{a, b} {
		if _, _, err := selador.Empresa(ctx, c.criada.EmpresaID); err != nil {
			t.Fatal(err)
		}
	}
	verificar := func(c conta) map[string]adm.DiaVerificado {
		t.Helper()
		dias, err := svcMem.VerificarIntegridadeAuditoria(ctxWeb(), c.adm, dia(-1), dia(1))
		if err != nil {
			t.Fatal(err)
		}
		m := map[string]adm.DiaVerificado{}
		for _, d := range dias {
			m[d.Dia.Format("2006-01-02")] = d
		}
		return m
	}
	rA := verificar(a)
	for _, n := range []int{0, 1} {
		d := rA[dia(n).Format("2006-01-02")]
		if d.Estado != "integro" {
			t.Errorf("A, dia %s: %s %v", dia(n).Format("2006-01-02"), d.Estado, d.Motivo)
		}
	}
	if d := rA[dia(-1).Format("2006-01-02")]; d.Estado != "nao_selado" {
		t.Errorf("A, dia anterior ao primeiro selo: %s", d.Estado)
	}
	// o selo está no armazenamento, na chave da empresa
	if ok, _ := arm.Existe(ctx, "selos/empresas/"+a.criada.EmpresaID.String()+"/"+dia(0).Format("2006-01-02")+".selo"); !ok {
		t.Error("o selo de A não foi para o armazenamento imutável")
	}

	// adulteração direta, como o papel dono: desliga o trigger, muda o evento, religa
	tx, err := amb.Migrador.Begin(ctx)
	if err != nil {
		t.Fatal(err)
	}
	for _, sql := range []string{
		`ALTER TABLE app.evento_auditoria DISABLE TRIGGER trg_somente_insercao`,
		fmt.Sprintf(`UPDATE app.evento_auditoria SET entidade = 'adulterada' WHERE id = '%s'`, idsA[1]),
		`ALTER TABLE app.evento_auditoria ENABLE TRIGGER trg_somente_insercao`,
	} {
		if _, err := tx.Exec(ctx, sql); err != nil {
			_ = tx.Rollback(ctx)
			t.Fatal(err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		t.Fatal(err)
	}
	rA = verificar(a)
	d := rA[dia(0).Format("2006-01-02")]
	if d.Estado != "violado" || d.Motivo == nil || *d.Motivo != "evento_alterado" || d.EventoID == nil || *d.EventoID != idsA[1] {
		t.Errorf("A depois da adulteração: %+v (esperado violado/evento_alterado apontando %s)", d, idsA[1])
	}
	if d := rA[dia(1).Format("2006-01-02")]; d.Estado != "integro" {
		t.Errorf("o dia intacto continua íntegro: %s", d.Estado)
	}
	// B não é afetada, e a verificação dela não enxerga os selos de A
	rB := verificar(b)
	if d := rB[dia(0).Format("2006-01-02")]; d.Estado != "integro" {
		t.Errorf("B: %+v", d)
	}
	if d := rB[dia(1).Format("2006-01-02")]; d.Estado != "integro" {
		t.Errorf("B, dia sem eventos de B: %+v (o selo de A não é dela)", d)
	}
	// sem o armazenamento a verificação não afirma nada
	_, err = svc.VerificarIntegridadeAuditoria(ctxWeb(), a.adm, dia(0), dia(1))
	if bucket == nil {
		exigirErro(t, err, erros.ServicoIndisponivel)
	}
}
