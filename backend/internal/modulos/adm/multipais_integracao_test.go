package adm_test

import (
	"context"
	crand "crypto/rand"
	"errors"
	"math/big"
	"strings"
	"testing"

	"github.com/google/uuid"

	"tenaz/internal/modulos/adm"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/erros"
	"tenaz/internal/plataforma/textos"
)

// Multipaís (ADR-013): país, moeda, idioma e identificação por tipo.

func digitosAleatorios(n int) string {
	var b strings.Builder
	for i := 0; i < n; i++ {
		d, _ := crand.Int(crand.Reader, big.NewInt(10))
		b.WriteByte(byte('0' + d.Int64()))
	}
	return b.String()
}

// cnpjAlfanumericoAleatorio sorteia 12 caracteres de [0-9A-Z] e calcula os dois
// dígitos pelo método da IN RFB 2.229/2024 (valor = código ASCII menos 48).
func cnpjAlfanumericoAleatorio() string {
	const alfabeto = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	base := make([]byte, 12)
	for i := range base {
		n, _ := crand.Int(crand.Reader, big.NewInt(int64(len(alfabeto))))
		base[i] = alfabeto[n.Int64()]
	}
	dv := func(corpo []byte, pesos []int) byte {
		soma := 0
		for i, p := range pesos {
			soma += int(corpo[i]-'0') * p
		}
		if r := soma % 11; r >= 2 {
			return byte('0' + 11 - r)
		}
		return '0'
	}
	base = append(base, dv(base, []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}))
	base = append(base, dv(base, []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}))
	return string(base)
}

// campoDoErro devolve o campo do primeiro erro de campo do erro de domínio.
func campoDoErro(t *testing.T, err error) string {
	t.Helper()
	var e *erros.Erro
	if !errors.As(err, &e) || len(e.Campos) == 0 {
		t.Fatalf("esperado erro de campo, veio %v", err)
	}
	return e.Campos[0].Campo
}

// ultimaMensagem devolve o assunto do último e-mail para o destino.
func (c *caixaDeEntrada) assunto(t *testing.T, para string) string {
	t.Helper()
	auth.Aguardar()
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := len(c.msgs) - 1; i >= 0; i-- {
		if c.msgs[i].Para == para {
			return c.msgs[i].Assunto
		}
	}
	t.Fatalf("nenhum e-mail para %s", para)
	return ""
}

func entradaConta(pais, tipo, valor, email string) adm.EntradaCriarConta {
	return adm.EntradaCriarConta{
		Pais: pais, IdentificacaoTipo: tipo, Identificacao: valor, RazaoSocial: "Empresa " + pais,
		Responsavel: adm.EntradaResponsavel{Nome: "Resp " + pais, Email: email},
		Unidade:     adm.EntradaUnidade{Codigo: "MTZ", Nome: "Matriz", Turnos: turnosComerciais()},
	}
}

// A criação da conta deriva moeda, idioma e fuso do país, grava tudo
// explicitamente e manda o convite no idioma do responsável.
func TestCriarContaPorPais(t *testing.T) {
	ctx := ctxWeb()
	casos := []struct {
		pais, tipo, valor, moeda, idioma, fuso, assunto string
	}{
		{"BR", "cnpj", cnpjAleatorio(), "BRL", "pt-BR", "America/Sao_Paulo", "Convite para a Tenaz"},
		{"AR", "cuit", "30-" + digitosAleatorios(8) + "-1", "ARS", "es-419", "America/Argentina/Buenos_Aires", "Invitación a Tenaz"},
		{"PY", "ruc", digitosAleatorios(8) + "-1", "PYG", "es-419", "America/Asuncion", "Invitación a Tenaz"},
	}
	for _, c := range casos {
		t.Run(c.pais, func(t *testing.T) {
			em := emailAleatorio("pais")
			criada, err := svc.CriarConta(ctx, entradaConta(c.pais, c.tipo, c.valor, em))
			if err != nil || !criada.ConviteEnviado {
				t.Fatalf("criar conta %s: %v", c.pais, err)
			}
			if a := caixa.assunto(t, em); !strings.HasPrefix(a, c.assunto) {
				t.Errorf("assunto do convite %q, esperado começar por %q", a, c.assunto)
			}
			if _, err := svc.AceitarConvite(ctx, caixa.token(t, em), senhaBoa, nil); err != nil {
				t.Fatal(err)
			}
			ator, _ := entrar(t, em, senhaBoa, &criada.EmpresaID)
			emp, err := svc.ObterEmpresa(ctx, ator)
			if err != nil {
				t.Fatal(err)
			}
			if emp.Pais != c.pais || emp.Moeda != c.moeda || emp.IdiomaPadrao != c.idioma || emp.FusoHorario != c.fuso ||
				emp.IdentificacaoTipo == nil || *emp.IdentificacaoTipo != c.tipo ||
				emp.Identificacao == nil || *emp.Identificacao != strings.ReplaceAll(strings.ReplaceAll(c.valor, "-", ""), ".", "") {
				t.Errorf("empresa: %+v", emp)
			}
			eu, err := svc.ObterEu(ctx, ator)
			if err != nil || eu.Idioma != c.idioma || eu.Moeda != c.moeda {
				t.Errorf("responsável: idioma %q, moeda %q (%v)", eu.Idioma, eu.Moeda, err)
			}
		})
	}
}

func TestCriarContaRecusaPaisETipoErrados(t *testing.T) {
	ctx := ctxWeb()
	casos := []struct {
		nome  string
		e     adm.EntradaCriarConta
		campo string
	}{
		{"sem país", entradaConta("", "cnpj", cnpjAleatorio(), emailAleatorio("x")), "pais"},
		{"país fora do catálogo", entradaConta("US", "cnpj", cnpjAleatorio(), emailAleatorio("x")), "pais"},
		{"cuit numa empresa brasileira", entradaConta("BR", "cuit", digitosAleatorios(11), emailAleatorio("x")), "identificacao_tipo"},
		{"cnpj numa empresa argentina", entradaConta("AR", "cnpj", cnpjAleatorio(), emailAleatorio("x")), "identificacao_tipo"},
		{"cpf não é de empresa", entradaConta("BR", "cpf", cpfAleatorio(), emailAleatorio("x")), "identificacao_tipo"},
		{"cnpj com dígito errado", entradaConta("BR", "cnpj", "11222333000180", emailAleatorio("x")), "identificacao"},
		{"ruc fora da forma", entradaConta("PY", "ruc", "123", emailAleatorio("x")), "identificacao"},
		{"sem identificação", entradaConta("AR", "", "", emailAleatorio("x")), "identificacao_tipo"},
	}
	for _, c := range casos {
		_, err := svc.CriarConta(ctx, c.e)
		exigirErro(t, err, erros.EntradaInvalida)
		if got := campoDoErro(t, err); got != c.campo {
			t.Errorf("%s: campo %q, esperado %q", c.nome, got, c.campo)
		}
	}

	// A unidade inicial só aceita o tipo do país da empresa.
	e := entradaConta("BR", "cnpj", cnpjAleatorio(), emailAleatorio("x"))
	e.Unidade.IdentificacaoTipo, e.Unidade.Identificacao = ptr("ruc"), ptr("800123456")
	_, err := svc.CriarConta(ctx, e)
	exigirErro(t, err, erros.EntradaInvalida)
	if got := campoDoErro(t, err); got != "unidade.identificacao_tipo" {
		t.Errorf("campo da unidade: %q", got)
	}

	// CNPJ alfanumérico (IN RFB 2.229/2024): o do exemplo oficial é válido (nos
	// casos compartilhados); aqui, um sorteado a cada execução.
	alfa := cnpjAlfanumericoAleatorio()
	if _, err := svc.CriarConta(ctx, entradaConta("BR", "cnpj", alfa, emailAleatorio("x"))); err != nil {
		t.Errorf("CNPJ alfanumérico válido recusado: %v", err)
	}
	if _, err := svc.CriarConta(ctx, entradaConta("BR", "cnpj", strings.ToLower(alfa), emailAleatorio("y"))); err == nil {
		t.Error("o mesmo CNPJ alfanumérico foi aceito duas vezes (unicidade na plataforma)")
	} else {
		exigirErro(t, err, erros.Duplicado)
	}
}

// País e moeda não mudam: a API não os aceita e o banco recusa a troca.
func TestPaisEMoedaNaoMudam(t *testing.T) {
	c := novaConta(t, "")
	_, err := amb.Migrador.Exec(context.Background(), `UPDATE app.empresa SET pais = 'AR', moeda = 'ARS' WHERE id = $1`, c.criada.EmpresaID)
	if err == nil || !strings.Contains(err.Error(), "não mudam") {
		t.Errorf("o banco deveria recusar a troca de país e moeda: %v", err)
	}
	emp, _ := svc.ObterEmpresa(ctxWeb(), c.adm)
	if emp.Pais != "BR" || emp.Moeda != "BRL" {
		t.Errorf("país e moeda mudaram: %+v", emp)
	}
}

// A identificação da empresa e da unidade é só do tipo do país; a da pessoa
// aceita o tipo de qualquer país.
func TestIdentificacaoDeEmpresaUnidadeEPessoa(t *testing.T) {
	c := novaConta(t, "profissional")
	ctx := ctxWeb()
	emp, _ := svc.ObterEmpresa(ctx, c.adm)

	// Empresa: trocar para o tipo de outro país é recusado.
	_, err := svc.EditarEmpresa(ctx, c.adm, adm.EntradaEmpresa{
		RazaoSocial: "X", IdentificacaoTipo: "cuit", Identificacao: digitosAleatorios(11),
		IdiomaPadrao: emp.IdiomaPadrao, FusoHorario: emp.FusoHorario, VersaoEsperada: &emp.Versao})
	exigirErro(t, err, erros.EntradaInvalida)
	novo := cnpjAleatorio()
	editada, err := svc.EditarEmpresa(ctx, c.adm, adm.EntradaEmpresa{
		RazaoSocial: "X", IdentificacaoTipo: "cnpj", Identificacao: novo,
		IdiomaPadrao: emp.IdiomaPadrao, FusoHorario: emp.FusoHorario, VersaoEsperada: &emp.Versao})
	if err != nil || editada.Identificacao == nil || *editada.Identificacao != novo || editada.Pais != "BR" {
		t.Fatalf("editar identificação: %+v %v", editada, err)
	}

	// Unidade: o tipo do país da empresa, junto ou nenhum.
	_, err = svc.CriarUnidade(ctx, c.adm, adm.EntradaUnidade{Codigo: "U2", Nome: "U2", Turnos: turnosComerciais(),
		IdentificacaoTipo: ptr("ruc"), Identificacao: ptr("800123456")})
	exigirErro(t, err, erros.EntradaInvalida)
	_, err = svc.CriarUnidade(ctx, c.adm, adm.EntradaUnidade{Codigo: "U2", Nome: "U2", Turnos: turnosComerciais(),
		IdentificacaoTipo: ptr("cnpj")})
	exigirErro(t, err, erros.EntradaInvalida)
	if got := campoDoErro(t, err); got != "identificacao" {
		t.Errorf("campo: %q", got)
	}
	un, err := svc.CriarUnidade(ctx, c.adm, adm.EntradaUnidade{Codigo: "U2", Nome: "U2", Turnos: turnosComerciais(),
		IdentificacaoTipo: ptr("cnpj"), Identificacao: ptr("11.222.333/0001-81")})
	if err != nil || un.Identificacao == nil || *un.Identificacao != "11222333000181" {
		t.Fatalf("unidade com CNPJ: %+v %v", un, err)
	}
	_, err = svc.EditarUnidade(ctx, c.adm, un.ID, adm.EntradaUnidade{Codigo: "U2", Nome: "U2",
		IdentificacaoTipo: ptr("cuit"), Identificacao: ptr("30693450239"), VersaoEsperada: &un.Versao})
	exigirErro(t, err, erros.EntradaInvalida)

	// Pessoa: cpf, cuil, dni e ci na mesma empresa brasileira; lista mascarada.
	casos := []struct{ tipo, valor, mascara string }{
		{"cpf", cpfAleatorio(), ""},
		{"cuil", "20" + digitosAleatorios(8) + "6", ""},
		{"dni", digitosAleatorios(8), ""},
		{"ci", digitosAleatorios(7), ""},
	}
	for _, x := range casos {
		p, err := svc.CriarPessoa(ctx, c.adm, adm.EntradaPessoa{Nome: "Pessoa " + x.tipo,
			IdentificacaoTipo: ptr(x.tipo), Identificacao: ptr(x.valor)})
		if err != nil {
			t.Fatalf("pessoa com %s: %v", x.tipo, err)
		}
		det, err := svc.ObterPessoa(ctx, c.adm, p.ID)
		if err != nil || det.Identificacao == nil || *det.Identificacao != x.valor || *det.IdentificacaoTipo != x.tipo {
			t.Errorf("detalhe com %s: %+v %v", x.tipo, det, err)
		}
		// A identificação repetida na empresa é conflito; em outro tipo, não.
		_, err = svc.CriarPessoa(ctx, c.adm, adm.EntradaPessoa{Nome: "Repetida", IdentificacaoTipo: ptr(x.tipo), Identificacao: ptr(x.valor)})
		exigirErro(t, err, erros.Duplicado)
	}
	pag, err := svc.ListarPessoas(ctx, c.adm, adm.FiltroPessoas{Busca: ptr("Pessoa dni")})
	if err != nil || len(pag.Dados) != 1 || pag.Dados[0].IdentificacaoMascarada == nil ||
		!strings.HasPrefix(*pag.Dados[0].IdentificacaoMascarada, "*****") {
		t.Fatalf("lista mascarada: %+v %v", pag.Dados, err)
	}
	// Tipo de empresa na pessoa, e tipo sem número, são recusados.
	_, err = svc.CriarPessoa(ctx, c.adm, adm.EntradaPessoa{Nome: "Y", IdentificacaoTipo: ptr("cnpj"), Identificacao: ptr(cnpjAleatorio())})
	exigirErro(t, err, erros.EntradaInvalida)
	_, err = svc.CriarPessoa(ctx, c.adm, adm.EntradaPessoa{Nome: "Y", IdentificacaoTipo: ptr("dni")})
	exigirErro(t, err, erros.EntradaInvalida)
}

// Idioma: o convite recebe o padrão da empresa; o aceite grava o escolhido; só
// o próprio usuário o altera; a moeda do usuário é a da empresa.
func TestIdiomaDoUsuario(t *testing.T) {
	c := novaConta(t, "")
	ctx := ctxWeb()

	u1, err := svc.CriarUsuario(ctx, c.adm, adm.EntradaUsuario{Nome: "Um", Email: emailAleatorio("um"), Perfis: []string{"tecnico"}})
	if err != nil || u1.Usuario.Idioma != "pt-BR" || u1.Usuario.Moeda != "BRL" {
		t.Fatalf("convite com o padrão da empresa: %+v %v", u1.Usuario, err)
	}
	if a := caixa.assunto(t, *u1.Usuario.Email); !strings.HasPrefix(a, "Convite para a Tenaz") {
		t.Errorf("assunto: %q", a)
	}

	// O administrador muda o padrão da empresa: vale para os convites seguintes.
	emp, _ := svc.ObterEmpresa(ctx, c.adm)
	if _, err := svc.EditarEmpresa(ctx, c.adm, adm.EntradaEmpresa{RazaoSocial: "X", IdentificacaoTipo: "cnpj",
		Identificacao: *emp.Identificacao, IdiomaPadrao: "es-419", FusoHorario: emp.FusoHorario, VersaoEsperada: &emp.Versao}); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.EditarEmpresa(ctx, c.adm, adm.EntradaEmpresa{RazaoSocial: "X", IdentificacaoTipo: "cnpj",
		Identificacao: *emp.Identificacao, IdiomaPadrao: "en", FusoHorario: emp.FusoHorario}); err == nil {
		t.Error("idioma fora do catálogo aceito")
	}
	em2 := emailAleatorio("dois")
	u2, err := svc.CriarUsuario(ctx, c.adm, adm.EntradaUsuario{Nome: "Dois", Email: em2, Perfis: []string{"tecnico"}})
	if err != nil || u2.Usuario.Idioma != "es-419" {
		t.Fatalf("convite com o novo padrão: %+v %v", u2.Usuario, err)
	}
	if a := caixa.assunto(t, em2); !strings.HasPrefix(a, "Invitación a Tenaz") {
		t.Errorf("e-mail do convite em es-419: %q", a)
	}

	// O aceite grava o idioma em que a pessoa aceitou.
	if _, err := svc.AceitarConvite(ctx, caixa.token(t, em2), senhaBoa, ptr("pt-BR")); err != nil {
		t.Fatal(err)
	}
	a2, _ := entrar(t, em2, senhaBoa, &c.criada.EmpresaID)
	eu, _ := svc.ObterEu(ctx, a2)
	if eu.Idioma != "pt-BR" {
		t.Errorf("idioma do aceite: %q", eu.Idioma)
	}
	if _, err := svc.AceitarConvite(ctx, caixa.token(t, *u1.Usuario.Email), senhaBoa, ptr("en")); err == nil {
		t.Error("idioma inválido no aceite")
	}

	// Só o próprio usuário altera o idioma.
	eu, err = svc.EditarEu(ctx, a2, adm.EntradaEu{Nome: eu.Nome, Idioma: ptr("es-419")})
	if err != nil || eu.Idioma != "es-419" {
		t.Fatalf("o próprio usuário troca o idioma: %+v %v", eu, err)
	}
	if _, err := svc.EditarEu(ctx, a2, adm.EntradaEu{Nome: eu.Nome, Idioma: ptr("fr")}); err == nil {
		t.Error("idioma fora do catálogo aceito em EditarEu")
	}
	eu, _ = svc.EditarEu(ctx, a2, adm.EntradaEu{Nome: "Dois Silva"})
	if eu.Idioma != "es-419" {
		t.Errorf("idioma ausente mantém o atual: %q", eu.Idioma)
	}
	// Quem administra usuários edita nome, telefone e custo; o idioma segue.
	u2e, err := svc.EditarUsuario(ctx, c.adm, u2.Usuario.ID, adm.EntradaEditarUsuario{Nome: "Dois por admin"})
	if err != nil || u2e.Idioma != "es-419" || u2e.Moeda != "BRL" {
		t.Errorf("EditarUsuario não mexe no idioma nem na moeda: %+v %v", u2e, err)
	}
}

// O e-mail de redefinição de senha sai no idioma negociado na requisição.
func TestRedefinicaoDeSenhaNoIdiomaDaRequisicao(t *testing.T) {
	c := novaConta(t, "")
	ctxEs := textos.ComPreferencias(ctxWeb(), textos.Preferencias{AcceptLanguage: "es-AR,es;q=0.9"})
	if err := auth.RecuperarSenha(ctxEs, c.email); err != nil {
		t.Fatal(err)
	}
	if a := caixa.assunto(t, c.email); !strings.HasPrefix(a, "Restablecimiento de contraseña") {
		t.Errorf("assunto em es-419: %q", a)
	}
	if err := auth.RecuperarSenha(ctxWeb(), c.email); err != nil {
		t.Fatal(err)
	}
	if a := caixa.assunto(t, c.email); !strings.HasPrefix(a, "Redefinição de senha") {
		t.Errorf("assunto sem Accept-Language: %q", a)
	}
}

// O fuso de empresa e unidade é um dos fusos do país da empresa (perfil de país).
func TestFusoValidadoContraOPais(t *testing.T) {
	ctx := ctxWeb()
	// Conta brasileira com fuso argentino: valor_invalido com os valores aceitos.
	e := entradaConta("BR", "cnpj", cnpjAleatorio(), emailAleatorio("fuso"))
	e.FusoHorario = "America/Argentina/Buenos_Aires"
	_, err := svc.CriarConta(ctx, e)
	exigirErro(t, err, erros.EntradaInvalida)
	var ee *erros.Erro
	if !errors.As(err, &ee) || ee.Campos[0].Campo != "fuso_horario" || ee.Campos[0].Codigo != erros.CampoValorInvalido ||
		len(ee.Campos[0].ValoresAceitos) != 16 {
		t.Fatalf("erro de fuso: %v", err)
	}
	// Nome que não é IANA continua recusado, e o de outro país também na unidade inicial.
	e.FusoHorario = "UTC"
	_, err = svc.CriarConta(ctx, e)
	exigirErro(t, err, erros.EntradaInvalida)
	e.FusoHorario = "America/Manaus"
	e.Unidade.FusoHorario = ptr("America/Asuncion")
	_, err = svc.CriarConta(ctx, e)
	if got := campoDoErro(t, err); got != "unidade.fuso_horario" {
		t.Errorf("fuso da unidade: %q", got)
	}
	e.Unidade.FusoHorario = ptr("America/Cuiaba")
	if _, err = svc.CriarConta(ctx, e); err != nil {
		t.Fatalf("fusos do Brasil: %v", err)
	}

	// Edição de empresa e de unidade, numa conta paraguaia (um só fuso).
	em := emailAleatorio("py")
	criada, err := svc.CriarConta(ctx, entradaConta("PY", "ruc", digitosAleatorios(8)+"-1", em))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.AceitarConvite(ctx, caixa.token(t, em), senhaBoa, nil); err != nil {
		t.Fatal(err)
	}
	ator, _ := entrar(t, em, senhaBoa, &criada.EmpresaID)
	emp, _ := svc.ObterEmpresa(ctx, ator)
	_, err = svc.EditarEmpresa(ctx, ator, adm.EntradaEmpresa{RazaoSocial: "X", IdentificacaoTipo: "ruc",
		Identificacao: *emp.Identificacao, IdiomaPadrao: emp.IdiomaPadrao, FusoHorario: "America/Sao_Paulo", VersaoEsperada: &emp.Versao})
	if got := campoDoErro(t, err); got != "fuso_horario" {
		t.Errorf("editar empresa com fuso de outro país: %q", got)
	}
	un, _ := svc.ObterUnidade(ctx, ator, criada.UnidadeID)
	_, err = svc.EditarUnidade(ctx, ator, un.ID, adm.EntradaUnidade{Codigo: un.Codigo, Nome: un.Nome,
		FusoHorario: ptr("America/Sao_Paulo"), VersaoEsperada: &un.Versao})
	if got := campoDoErro(t, err); got != "fuso_horario" {
		t.Errorf("editar unidade com fuso de outro país: %q", got)
	}
	if _, err := svc.EditarUnidade(ctx, ator, un.ID, adm.EntradaUnidade{Codigo: un.Codigo, Nome: un.Nome,
		FusoHorario: ptr("America/Asuncion"), VersaoEsperada: &un.Versao}); err != nil {
		t.Errorf("fuso do Paraguai: %v", err)
	}
}

// O endereço da unidade é conferido pelo perfil do país da empresa
// (documento 3, 2.18): subdivisão em lista fechada, código postal por forma, e
// a regra da Tenaz do CPA argentino.
func TestEnderecoDaUnidadeConferidoPeloPais(t *testing.T) {
	ctx := ctxWeb()
	turno := []adm.EntradaTurno{{DiaSemana: 1, HoraInicio: "08:00", HoraFim: "17:00"}}
	nova := func(pais, tipo, ident string) (autorizacao.Ator, adm.Unidade) {
		em := emailAleatorio("end")
		entrada := entradaConta(pais, tipo, ident, em)
		entrada.Plano = ptr("profissional") // o Essencial tem uma unidade só
		criada, err := svc.CriarConta(ctx, entrada)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := svc.AceitarConvite(ctx, caixa.token(t, em), senhaBoa, nil); err != nil {
			t.Fatal(err)
		}
		ator, _ := entrar(t, em, senhaBoa, &criada.EmpresaID)
		un, err := svc.ObterUnidade(ctx, ator, criada.UnidadeID)
		if err != nil {
			t.Fatal(err)
		}
		return ator, un
	}
	editar := func(a autorizacao.Ator, un adm.Unidade, e *adm.Endereco) (adm.Unidade, error) {
		atual, err := svc.ObterUnidade(ctx, a, un.ID) // a versão sobe a cada edição certa
		if err != nil {
			t.Fatal(err)
		}
		return svc.EditarUnidade(ctx, a, un.ID, adm.EntradaUnidade{Codigo: un.Codigo, Nome: un.Nome, Endereco: e, VersaoEsperada: &atual.Versao})
	}
	criar := func(a autorizacao.Ator, cod string, e *adm.Endereco) (adm.Unidade, error) {
		return svc.CriarUnidade(ctx, a, adm.EntradaUnidade{Codigo: cod, Nome: "Filial", Endereco: e, Turnos: turno})
	}
	campo := func(err error) (string, string, int) {
		var ee *erros.Erro
		if !errors.As(err, &ee) || len(ee.Campos) != 1 {
			t.Fatalf("esperava um erro de campo, veio %v", err)
		}
		return ee.Campos[0].Campo, ee.Campos[0].Codigo, len(ee.Campos[0].ValoresAceitos)
	}

	// Argentina: CPA normalizado, subdivisão da lista, letra do CPA = a da subdivisão.
	ar, unAR := nova("AR", "cuit", "30-"+digitosAleatorios(8)+"-1")
	u, err := editar(ar, unAR, &adm.Endereco{Municipio: "CABA", Subdivisao: "AR-C", CodigoPostal: "c1425 dkf"})
	if err != nil {
		t.Fatalf("endereço argentino válido: %v", err)
	}
	if u.Endereco == nil || u.Endereco.CodigoPostal != "C1425DKF" || u.Endereco.Subdivisao != "AR-C" {
		t.Errorf("endereço gravado: %+v", u.Endereco)
	}
	if _, err := editar(ar, u, &adm.Endereco{Subdivisao: "AR-B", CodigoPostal: "C1425DKF"}); err == nil {
		t.Error("CPA de AR-C com AR-B passou")
	} else if c, cod, _ := campo(err); c != "endereco.codigo_postal" || cod != erros.CampoValorInvalido {
		t.Errorf("CPA x subdivisão: %s %s", c, cod)
	}
	if _, err := editar(ar, u, &adm.Endereco{Subdivisao: "AR-B", CodigoPostal: "1425"}); err != nil {
		t.Errorf("código antigo de 4 dígitos não cruza: %v", err)
	}
	if _, err := editar(ar, u, &adm.Endereco{Subdivisao: "BR-SP"}); err == nil {
		t.Error("subdivisão brasileira numa empresa argentina passou")
	} else if c, cod, n := campo(err); c != "endereco.subdivisao" || cod != erros.CampoValorInvalido || n != 24 {
		t.Errorf("subdivisão de outro país: %s %s %d", c, cod, n)
	}
	if _, err := criar(ar, "F2", &adm.Endereco{CodigoPostal: "I1425DKF"}); err == nil {
		t.Error("CPA com a letra I passou")
	} else if c, cod, _ := campo(err); c != "endereco.codigo_postal" || cod != erros.CampoFormatoInvalido {
		t.Errorf("CPA inválido: %s %s", c, cod)
	}

	// Paraguai: sem zero à esquerda, seis dígitos.
	py, unPY := nova("PY", "ruc", digitosAleatorios(8)+"-1")
	if _, err := editar(py, unPY, &adm.Endereco{Subdivisao: "PY-ASU", CodigoPostal: "001-526"}); err != nil {
		t.Errorf("endereço paraguaio válido: %v", err)
	}
	if _, err := editar(py, unPY, &adm.Endereco{Subdivisao: "PY-01"}); err == nil {
		t.Error("PY-01 passou")
	} else if _, _, n := campo(err); n != 18 {
		t.Errorf("valores_aceitos do Paraguai: %d", n)
	}
	if _, err := editar(py, unPY, &adm.Endereco{CodigoPostal: "00152"}); err == nil {
		t.Error("código postal de cinco dígitos passou")
	}

	// Brasil: CEP com hífen normaliza; subdivisão fora da lista falha; sem endereço, nada a conferir.
	br := novaConta(t, "")
	unBR, err := svc.ObterUnidade(ctx, br.adm, br.criada.UnidadeID)
	if err != nil {
		t.Fatal(err)
	}
	u, err = editar(br.adm, unBR, &adm.Endereco{Subdivisao: "BR-SP", CodigoPostal: "01310-100", Municipio: "São Paulo"})
	if err != nil || u.Endereco.CodigoPostal != "01310100" {
		t.Fatalf("endereço brasileiro: %v %+v", err, u.Endereco)
	}
	if _, err := editar(br.adm, u, &adm.Endereco{Subdivisao: "SP"}); err == nil {
		t.Error("UF sem prefixo passou")
	}
	if _, err := editar(br.adm, u, nil); err != nil {
		t.Errorf("sem endereço: %v", err)
	}
	// Conta nova com endereço na unidade inicial: o campo leva o prefixo unidade.
	e := entradaConta("BR", "cnpj", cnpjAleatorio(), emailAleatorio("endini"))
	e.Unidade.Endereco = &adm.Endereco{Subdivisao: "BR-XX"}
	_, err = svc.CriarConta(ctx, e)
	if c, _, _ := campo(err); c != "unidade.endereco.subdivisao" {
		t.Errorf("campo da unidade inicial: %s", c)
	}
}

// custo_hora é preço unitário (documento 5, Parte F): até seis casas, sem
// sinal, lido com a moeda da empresa e nunca arredondado.
func TestCustoHoraComAMoedaDaEmpresa(t *testing.T) {
	ctx := ctxWeb()
	criar := func(a autorizacao.Ator, custo *string) (adm.UsuarioCriado, error) {
		return svc.CriarUsuario(ctx, a, adm.EntradaUsuario{Nome: "Custo", Email: emailAleatorio("custo"),
			Perfis: []string{"tecnico"}, CustoHora: custo})
	}
	texto := func(u adm.Usuario) string {
		if u.CustoHora == nil {
			return "<nulo>"
		}
		return u.CustoHora.String()
	}
	campoFormato := func(err error) bool {
		var e *erros.Erro
		return errors.As(err, &e) && len(e.Campos) == 1 && e.Campos[0].Campo == "custo_hora" &&
			e.Campos[0].Codigo == erros.CampoFormatoInvalido
	}

	// Empresa brasileira (BRL): mínimo de duas casas na saída, até seis na entrada.
	br := novaConta(t, "")
	for entrada, esperado := range map[string]string{"85.5": "85.50", "85": "85.00", "6.299": "6.299", "37.846200": "37.8462", "0.000001": "0.000001", "0": "0.00"} {
		u, err := criar(br.adm, &entrada)
		if err != nil || texto(u.Usuario) != esperado || u.Usuario.Moeda != "BRL" {
			t.Errorf("BRL %q: (%s, %v), esperado %s", entrada, texto(u.Usuario), err, esperado)
		}
	}
	for _, ruim := range []string{"0.1234567", "-1", "-0", "+1", "1,5", "1e3", " 1", "", "01", "1.", ".5", "1000000000000"} {
		ruim := ruim
		if _, err := criar(br.adm, &ruim); !campoFormato(err) {
			t.Errorf("BRL %q deveria ser formato_invalido em custo_hora: %v", ruim, err)
		}
	}
	// No banco, numeric(18,6) exato: nada truncado, nada arredondado.
	u, err := criar(br.adm, ptr("45.123456"))
	if err != nil {
		t.Fatal(err)
	}
	var cru string
	if err := amb.Migrador.QueryRow(context.Background(), `SELECT custo_hora::text FROM app.usuario WHERE id = $1`, u.Usuario.ID).Scan(&cru); err != nil || cru != "45.123456" {
		t.Errorf("custo_hora no banco: %q %v", cru, err)
	}
	// Edição: troca, e nulo apaga.
	ed, err := svc.EditarUsuario(ctx, br.adm, u.Usuario.ID, adm.EntradaEditarUsuario{Nome: "Custo", CustoHora: ptr("50")})
	if err != nil || texto(ed) != "50.00" {
		t.Errorf("editar: %s %v", texto(ed), err)
	}
	if ed, err = svc.EditarUsuario(ctx, br.adm, u.Usuario.ID, adm.EntradaEditarUsuario{Nome: "Custo"}); err != nil || ed.CustoHora != nil {
		t.Errorf("custo nulo: %v %v", ed.CustoHora, err)
	}
	if _, err = svc.EditarUsuario(ctx, br.adm, u.Usuario.ID, adm.EntradaEditarUsuario{Nome: "Custo", CustoHora: ptr("1.0000001")}); !campoFormato(err) {
		t.Errorf("editar com sétima casa: %v", err)
	}

	// Empresa paraguaia (PYG, sem casas): zeros além da escala valem no valor,
	// mas o custo-hora é preço unitário, então "0.5" guarani é legítimo.
	em := emailAleatorio("pycusto")
	criada, err := svc.CriarConta(ctx, entradaConta("PY", "ruc", digitosAleatorios(8)+"-1", em))
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.AceitarConvite(ctx, caixa.token(t, em), senhaBoa, nil); err != nil {
		t.Fatal(err)
	}
	py, _ := entrar(t, em, senhaBoa, &criada.EmpresaID)
	for entrada, esperado := range map[string]string{"75000": "75000", "75000.00": "75000", "0.5": "0.5", "37.8462": "37.8462", "0": "0"} {
		entrada := entrada
		u, err := criar(py, &entrada)
		if err != nil || texto(u.Usuario) != esperado || u.Usuario.Moeda != "PYG" {
			t.Errorf("PYG %q: (%s, %v), esperado %s", entrada, texto(u.Usuario), err, esperado)
		}
	}
	// "75000.50" é um preço unitário válido (meio guarani); o que a escala da
	// moeda recusa é o VALOR, e nenhum valor entra neste cadastro.
	if u, err := criar(py, ptr("75000.50")); err != nil || texto(u.Usuario) != "75000.5" {
		t.Errorf("PYG preço de 75000.50: (%s, %v)", texto(u.Usuario), err)
	}
}

// Controle positivo do externo: o prestador e o locatário trocam o idioma da
// própria linha por PATCH /v1/usuarios/eu (EditarEu), e nenhuma outra coluna
// muda: nem as que o serviço não expõe (perfis, vínculo, situação, custo,
// moeda), nem o que a rota de gestão permitiria a um interno. A suíte de
// isolamento prova o contrário (as demais colunas recusadas no banco); aqui
// se prova que o caminho legítimo funciona.
func TestExternoAlteraSoOProprioIdioma(t *testing.T) {
	c := novaConta(t, "")
	ctx := ctxWeb()
	forn, loc := uuid.New(), uuid.New()
	for _, e := range []struct {
		rotulo string
		in     adm.EntradaUsuario
	}{
		{"prestador", adm.EntradaUsuario{Perfis: []string{"prestador"}, FornecedorID: &forn}},
		{"locatario", adm.EntradaUsuario{Perfis: []string{"locatario"}, LocatarioID: &loc}},
	} {
		u, a := convidarEEntrar(t, c, e.in)
		if u.Idioma != "pt-BR" {
			t.Fatalf("%s: idioma inicial %q", e.rotulo, u.Idioma)
		}
		// Tudo que não é nome, telefone nem idioma, nem os campos que o gatilho de
		// padrões mantém (versão, atualização).
		fotografia := func() string {
			var j string
			if err := amb.Migrador.QueryRow(context.Background(), `SELECT (to_jsonb(u) - ARRAY['nome','telefone','idioma','versao','atualizado_em','atualizado_por','alterado_xid'])::text
			    FROM app.usuario u WHERE id = $1`, u.ID).Scan(&j); err != nil {
				t.Fatal(err)
			}
			return j
		}
		antes := fotografia()

		eu, err := svc.EditarEu(ctx, a, adm.EntradaEu{Nome: u.Nome, Idioma: ptr("es-419")})
		if err != nil || eu.Idioma != "es-419" {
			t.Fatalf("%s troca o próprio idioma: %+v %v", e.rotulo, eu, err)
		}
		var gravado string
		if err := amb.Migrador.QueryRow(context.Background(), `SELECT idioma FROM app.usuario WHERE id = $1`, u.ID).Scan(&gravado); err != nil || gravado != "es-419" {
			t.Fatalf("%s: idioma gravado %q %v", e.rotulo, gravado, err)
		}
		// Idioma ausente mantém o atual; fora do catálogo é recusado e não grava.
		if eu, err = svc.EditarEu(ctx, a, adm.EntradaEu{Nome: "Renomeado " + e.rotulo}); err != nil || eu.Idioma != "es-419" || eu.Nome != "Renomeado "+e.rotulo {
			t.Fatalf("%s: nome sem idioma: %+v %v", e.rotulo, eu, err)
		}
		if _, err := svc.EditarEu(ctx, a, adm.EntradaEu{Nome: u.Nome, Idioma: ptr("fr")}); campoDoErro(t, err) != "idioma" {
			t.Errorf("%s: idioma fora do catálogo: %v", e.rotulo, err)
		}
		if depois := fotografia(); depois != antes {
			t.Errorf("%s: outra coluna mudou\n antes: %s\n depois: %s", e.rotulo, antes, depois)
		}
		// A gestão de usuários continua vedada ao externo: não troca o idioma de ninguém.
		if _, err := svc.EditarUsuario(ctx, a, c.adm.UsuarioID, adm.EntradaEditarUsuario{Nome: "Invasão"}); err == nil {
			t.Errorf("%s editou outro usuário", e.rotulo)
		}
	}
}
