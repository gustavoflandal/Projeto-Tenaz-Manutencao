package adm_test

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"tenaz/internal/modulos/adm"
	"tenaz/internal/plataforma/autenticacao"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/erros"
)

func exigirErro(t *testing.T, err error, alvo *erros.Erro) {
	t.Helper()
	if !errors.Is(err, alvo) {
		t.Fatalf("esperado %s, veio %v", alvo.Codigo, err)
	}
}

// RF-ADM-001: a empresa nasce com a unidade, os turnos, os 12 perfis com a
// matriz de RF-ADM-005 e o responsável administrador; a trilha registra.
func TestCriarConta(t *testing.T) {
	c := novaConta(t, "")
	ctx := ctxWeb()

	emp, err := svc.ObterEmpresa(ctx, c.adm)
	if err != nil || emp.Plano == nil || *emp.Plano != "essencial" || emp.Situacao != "teste" {
		t.Fatalf("empresa: %+v %v", emp, err)
	}
	perfis, err := svc.ListarPerfis(ctx, c.adm)
	if err != nil || len(perfis) != 12 {
		t.Fatalf("perfis: %d %v", len(perfis), err)
	}
	for _, p := range perfis {
		if !p.Padrao || !slices.Equal(p.Permissoes, autorizacao.PermissoesDoPerfil(p.Codigo)) {
			t.Errorf("perfil %s diverge da matriz", p.Codigo)
		}
	}
	unidades, err := svc.ListarUnidades(ctx, c.adm)
	if err != nil || len(unidades) != 1 {
		t.Fatalf("unidades: %v %v", unidades, err)
	}
	u, err := svc.ObterUnidade(ctx, c.adm, unidades[0].ID)
	if err != nil || len(u.Turnos) != 5 {
		t.Fatalf("turno padrão de segunda a sexta: %+v %v", u.Turnos, err)
	}
	eu, err := svc.ObterEu(ctx, c.adm)
	if err != nil || eu.Situacao != "ativo" || !eu.IdentidadeLigada || len(eu.Perfis) != 1 || eu.Perfis[0].Codigo != "administrador" {
		t.Fatalf("responsável: %+v %v", eu, err)
	}
	if c.adm.Pode(autorizacao.OsEncerrar) || !c.adm.Pode(autorizacao.AdmGerenciarUsuarios) || c.adm.Escopo.PodeLerSaude {
		t.Error("permissões efetivas do administrador")
	}

	// A trilha tem a criação da empresa (trigger) e o login (serviço).
	for _, f := range []adm.FiltroAuditoria{{Entidade: ptr("empresa"), Acao: ptr("inclusao")}, {Acao: ptr("login")}} {
		pag, err := svc.ListarEventosAuditoria(ctx, c.adm, f)
		if err != nil || len(pag.Dados) == 0 {
			t.Errorf("trilha %+v: %v %v", f, pag.Dados, err)
		}
	}

	// CNPJ repetido.
	_, err = svc.CriarConta(ctx, adm.EntradaCriarConta{CNPJ: *emp.CNPJ, RazaoSocial: "Outra",
		Responsavel: adm.EntradaResponsavel{Nome: "B", Email: emailAleatorio("b")}, Unidade: adm.EntradaUnidade{Codigo: "U", Nome: "U"}})
	exigirErro(t, err, erros.Duplicado)
	_, err = svc.CriarConta(ctx, adm.EntradaCriarConta{CNPJ: "11222333000180", RazaoSocial: "X",
		Responsavel: adm.EntradaResponsavel{Nome: "B", Email: emailAleatorio("b")}, Unidade: adm.EntradaUnidade{Codigo: "U", Nome: "U"}})
	exigirErro(t, err, erros.EntradaInvalida)
}

func ptr[T any](v T) *T { return &v }

// Pendência E01-10: a lista de vedadas a externo do banco é a do catálogo.
func TestVedadasExternoIgualAoBanco(t *testing.T) {
	var banco []string
	if err := amb.Migrador.QueryRow(context.Background(), `SELECT app.permissoes_vedadas_externo()`).Scan(&banco); err != nil {
		t.Fatal(err)
	}
	codigo := autorizacao.VedadasExterno()
	slices.Sort(banco)
	slices.Sort(codigo)
	if !slices.Equal(banco, codigo) {
		t.Fatalf("app.permissoes_vedadas_externo() diverge do catálogo:\n banco:  %v\n código: %v", banco, codigo)
	}
}

// RF-ADM-002: Essencial tem uma unidade; pelo menos um turno; feriado único.
func TestUnidadesTurnosFeriados(t *testing.T) {
	c := novaConta(t, "essencial")
	ctx := ctxWeb()
	_, err := svc.CriarUnidade(ctx, c.adm, adm.EntradaUnidade{Codigo: "F2", Nome: "Filial",
		Turnos: []adm.EntradaTurno{{DiaSemana: 1, HoraInicio: "08:00", HoraFim: "17:00"}}})
	exigirErro(t, err, erros.PlanoNaoIncluiRecurso)

	pro := novaConta(t, "profissional")
	f2, err := svc.CriarUnidade(ctx, pro.adm, adm.EntradaUnidade{Codigo: "F2", Nome: "Filial",
		Turnos: []adm.EntradaTurno{{DiaSemana: 6, HoraInicio: "22:00", HoraFim: "06:00"}}})
	if err != nil {
		t.Fatalf("Profissional cria a segunda unidade: %v", err)
	}
	_, err = svc.CriarUnidade(ctx, pro.adm, adm.EntradaUnidade{Codigo: "F2", Nome: "Repetida",
		Turnos: []adm.EntradaTurno{{DiaSemana: 1, HoraInicio: "08:00", HoraFim: "17:00"}}})
	exigirErro(t, err, erros.Duplicado)

	err = svc.RemoverTurno(ctx, pro.adm, f2.ID, f2.Turnos[0].ID)
	exigirErro(t, err, erros.EntradaInvalida)
	t2, err := svc.IncluirTurno(ctx, pro.adm, f2.ID, adm.EntradaTurno{DiaSemana: 7, HoraInicio: "00:00", HoraFim: "24:00"})
	if err != nil || t2.HoraFim != "24:00" {
		t.Fatalf("turno até a meia-noite: %+v %v", t2, err)
	}
	if err := svc.RemoverTurno(ctx, pro.adm, f2.ID, f2.Turnos[0].ID); err != nil {
		t.Fatalf("com dois turnos, remove um: %v", err)
	}

	if _, err := svc.IncluirFeriado(ctx, pro.adm, f2.ID, adm.EntradaFeriado{Data: "2026-12-25", Descricao: "Natal"}); err != nil {
		t.Fatal(err)
	}
	_, err = svc.IncluirFeriado(ctx, pro.adm, f2.ID, adm.EntradaFeriado{Data: "2026-12-25", Descricao: "De novo"})
	exigirErro(t, err, erros.Duplicado)
	fs, err := svc.ListarFeriados(ctx, pro.adm, f2.ID, ptr("2026-12-01"), ptr("2026-12-31"))
	if err != nil || len(fs) != 1 {
		t.Fatalf("feriados: %v %v", fs, err)
	}

	// Versão otimista.
	_, err = svc.EditarUnidade(ctx, pro.adm, f2.ID, adm.EntradaUnidade{Codigo: "F2", Nome: "Filial 2", VersaoEsperada: ptr(int64(99))})
	exigirErro(t, err, erros.VersaoDesatualizada)
	ed, err := svc.EditarUnidade(ctx, pro.adm, f2.ID, adm.EntradaUnidade{Codigo: "F2", Nome: "Filial 2", VersaoEsperada: &f2.Versao})
	if err != nil || ed.Nome != "Filial 2" {
		t.Fatalf("editar unidade: %v", err)
	}

	// A outra empresa não enxerga a unidade.
	_, err = svc.ObterUnidade(ctx, c.adm, f2.ID)
	exigirErro(t, err, erros.NaoEncontrado)
}

// RF-ADM-015: pessoa com e sem usuário, CPF válido e único, CPF mascarado na
// lista e auditado no detalhe.
func TestPessoas(t *testing.T) {
	c := novaConta(t, "")
	ctx := ctxWeb()
	cpf := cpfAleatorio()
	_, err := svc.CriarPessoa(ctx, c.adm, adm.EntradaPessoa{Nome: "Carlos", CPF: ptr("123.456.789-00")})
	exigirErro(t, err, erros.EntradaInvalida)

	semUsuario, err := svc.CriarPessoa(ctx, c.adm, adm.EntradaPessoa{Nome: "Carlos Motorista", CPF: &cpf})
	if err != nil || semUsuario.UsuarioID != nil {
		t.Fatalf("pessoa sem usuário: %+v %v", semUsuario, err)
	}
	_, err = svc.CriarPessoa(ctx, c.adm, adm.EntradaPessoa{Nome: "Outro", CPF: &cpf})
	exigirErro(t, err, erros.Duplicado)

	// Pessoa com usuário, nas duas direções: usuário criado junto com a
	// pessoa, e usuário ligado a pessoa existente.
	u1, _ := convidarEEntrar(t, c, adm.EntradaUsuario{Perfis: []string{"tecnico"},
		NovaPessoa: &adm.EntradaPessoa{Nome: "Técnica Bia", CPF: ptr(cpfAleatorio())}})
	if u1.PessoaID == nil {
		t.Fatal("usuário criado com a pessoa")
	}
	u2, err := svc.CriarUsuario(ctx, c.adm, adm.EntradaUsuario{Nome: "Carlos", Email: emailAleatorio("carlos"),
		Perfis: []string{"motorista"}, PessoaID: &semUsuario.ID})
	if err != nil || u2.Usuario.PessoaID == nil || *u2.Usuario.PessoaID != semUsuario.ID {
		t.Fatalf("ligar a pessoa existente: %v", err)
	}
	_, err = svc.CriarUsuario(ctx, c.adm, adm.EntradaUsuario{Nome: "Outro", Email: emailAleatorio("x"),
		Perfis: []string{"motorista"}, PessoaID: &semUsuario.ID})
	exigirErro(t, err, erros.Duplicado)

	pag, err := svc.ListarPessoas(ctx, c.adm, adm.FiltroPessoas{Busca: ptr("carlos")})
	if err != nil || len(pag.Dados) != 1 || pag.Dados[0].CPFMascarado == nil || strings.Contains(*pag.Dados[0].CPFMascarado, cpf[:3]) {
		t.Fatalf("lista com CPF mascarado: %+v %v", pag.Dados, err)
	}
	if pag.Dados[0].UsuarioID == nil {
		t.Error("a lista mostra o usuário ligado")
	}
	det, err := svc.ObterPessoa(ctx, c.adm, semUsuario.ID)
	if err != nil || det.CPF == nil || *det.CPF != cpf {
		t.Fatalf("detalhe com CPF: %v", err)
	}
	ev, err := svc.ListarEventosAuditoria(ctx, c.adm, adm.FiltroAuditoria{Acao: ptr("visualizacao_sensivel"), EntidadeID: &semUsuario.ID})
	if err != nil || len(ev.Dados) != 1 {
		t.Fatalf("visualização do CPF na trilha: %v %v", ev.Dados, err)
	}

	in, err := svc.InativarPessoa(ctx, c.adm, semUsuario.ID)
	if err != nil || in.Situacao != "inativa" {
		t.Fatalf("inativar: %v", err)
	}
	_, err = svc.InativarPessoa(ctx, c.adm, semUsuario.ID)
	exigirErro(t, err, erros.TransicaoInvalida)

	// Técnico não lê pessoas.
	_, tec := convidarEEntrar(t, c, adm.EntradaUsuario{Perfis: []string{"tecnico"}})
	_, err = svc.ListarPessoas(ctx, tec, adm.FiltroPessoas{})
	exigirErro(t, err, erros.SemPermissao)
}

// RF-ADM-004 e RF-ADM-005: perfis, externos, último administrador,
// anti-escalada, permissão avulsa e pessoa.saude.ler.
func TestUsuariosEPermissoes(t *testing.T) {
	c := novaConta(t, "")
	ctx := ctxWeb()

	// Último administrador.
	_, err := svc.InativarUsuario(ctx, c.adm, c.adm.UsuarioID)
	exigirErro(t, err, adm.ErrUltimoAdministrador)
	_, err = svc.DefinirPerfis(ctx, c.adm, c.adm.UsuarioID, []string{"gestor"})
	exigirErro(t, err, erros.SemPermissao) // ninguém altera os próprios perfis

	gestor, aGestor := convidarEEntrar(t, c, adm.EntradaUsuario{Perfis: []string{"gestor"}})
	if !aGestor.Pode(autorizacao.OsEncerrar) || aGestor.Pode(autorizacao.AdmGerenciarUsuarios) {
		t.Error("gestor libera OS e não gerencia usuários")
	}
	_, err = svc.CriarUsuario(ctx, aGestor, adm.EntradaUsuario{Nome: "X", Email: emailAleatorio("x"), Perfis: []string{"tecnico"}})
	exigirErro(t, err, erros.SemPermissao)

	// Externo: prestador só com fornecedor e só com o perfil prestador.
	forn := uuid.New()
	_, err = svc.CriarUsuario(ctx, c.adm, adm.EntradaUsuario{Nome: "P", Email: emailAleatorio("p"), Perfis: []string{"prestador"}})
	exigirErro(t, err, erros.EntradaInvalida)
	_, err = svc.CriarUsuario(ctx, c.adm, adm.EntradaUsuario{Nome: "P", Email: emailAleatorio("p"), Perfis: []string{"prestador", "tecnico"}, FornecedorID: &forn})
	exigirErro(t, err, erros.EntradaInvalida)
	prest, aPrest := convidarEEntrar(t, c, adm.EntradaUsuario{Perfis: []string{"prestador"}, FornecedorID: &forn})
	if aPrest.Escopo.Tipo != "prestador" || aPrest.Escopo.FornecedorID == nil || *aPrest.Escopo.FornecedorID != forn {
		t.Fatalf("escopo do prestador: %+v", aPrest.Escopo)
	}
	_, err = svc.ListarUsuarios(ctx, aPrest, adm.FiltroUsuarios{})
	exigirErro(t, err, erros.NaoEncontrado) // externo não vê a rota: 404
	eu, err := svc.EditarEu(ctx, aPrest, adm.EntradaEu{Nome: "Prestador Renomeado"})
	if err != nil || eu.Nome != "Prestador Renomeado" || len(eu.Perfis) != 1 {
		t.Fatalf("externo edita nome: %+v %v", eu, err)
	}
	_, err = svc.ConcederPermissao(ctx, c.adm, prest.ID, autorizacao.OsEncerrar)
	exigirErro(t, err, erros.EntradaInvalida) // vedada a externo

	// Permissão avulsa: só o administrador concede pessoa.saude.ler, e não a si.
	_, err = svc.ConcederPermissao(ctx, c.adm, c.adm.UsuarioID, autorizacao.PessoaSaudeLer)
	exigirErro(t, err, erros.SemPermissao)
	_, err = svc.ConcederPermissao(ctx, aGestor, c.adm.UsuarioID, autorizacao.PessoaSaudeLer)
	exigirErro(t, err, erros.SemPermissao)
	u, err := svc.ConcederPermissao(ctx, c.adm, gestor.ID, autorizacao.PessoaSaudeLer)
	if err != nil || !slices.Contains(u.PermissoesAvulsas, autorizacao.PessoaSaudeLer) {
		t.Fatalf("conceder pessoa.saude.ler: %v", err)
	}
	_, err = svc.ConcederPermissao(ctx, c.adm, gestor.ID, autorizacao.PessoaSaudeLer)
	exigirErro(t, err, erros.Duplicado)
	recarregado, err := auth.CarregarAtor(ctx, c.criada.EmpresaID, gestor.ID)
	if err != nil || !recarregado.Escopo.PodeLerSaude {
		t.Fatalf("PodeLerSaude vem das permissões efetivas: %v", err)
	}
	if _, err := svc.RetirarPermissao(ctx, c.adm, gestor.ID, autorizacao.PessoaSaudeLer); err != nil {
		t.Fatal(err)
	}
	recarregado, _ = auth.CarregarAtor(ctx, c.criada.EmpresaID, gestor.ID)
	if recarregado.Escopo.PodeLerSaude {
		t.Error("retirada vale na próxima requisição")
	}

	// Dois administradores: um pode deixar de ser.
	_, err = svc.DefinirPerfis(ctx, c.adm, gestor.ID, []string{"gestor", "administrador"})
	if err != nil {
		t.Fatalf("promover a administrador: %v", err)
	}
	_, err = svc.DefinirPerfis(ctx, c.adm, gestor.ID, []string{})
	exigirErro(t, err, erros.EntradaInvalida) // pelo menos um perfil

	// Inativar encerra o acesso: a sessão deixa de valer.
	tec, _ := convidarEEntrar(t, c, adm.EntradaUsuario{Perfis: []string{"tecnico"}})
	r, err := auth.Login(ctxWeb(), autenticacao.EntradaLogin{Email: *tec.Email, Senha: senhaBoa})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := svc.InativarUsuario(ctx, c.adm, tec.ID); err != nil {
		t.Fatal(err)
	}
	_, err = auth.AutenticarWeb(ctxWeb(), r.Sessao.Cookie, r.Sessao.CSRF, false)
	exigirErro(t, err, erros.CredencialInvalida)
	var vivas int
	if err := amb.Migrador.QueryRow(context.Background(),
		`SELECT count(*) FROM plataforma.sessao WHERE usuario_id = $1 AND revogada_em IS NULL`, tec.ID).Scan(&vivas); err != nil || vivas != 0 {
		t.Errorf("a inativação revoga as sessões na mesma transação (trigger de acesso_empresa): %d vivas, %v", vivas, err)
	}
	re, err := svc.ReativarUsuario(ctx, c.adm, tec.ID)
	if err != nil || re.Usuario.Situacao != "ativo" {
		t.Fatalf("reativar: %v", err)
	}

	// Lista com cursor.
	p1, err := svc.ListarUsuarios(ctx, c.adm, adm.FiltroUsuarios{Limite: 2})
	if err != nil || len(p1.Dados) != 2 || !p1.TemMais {
		t.Fatalf("primeira página: %v %v", p1, err)
	}
	p2, err := svc.ListarUsuarios(ctx, c.adm, adm.FiltroUsuarios{Limite: 2, Cursor: p1.ProximoCursor})
	if err != nil || len(p2.Dados) == 0 || p2.Dados[0].ID == p1.Dados[1].ID {
		t.Fatalf("segunda página: %v %v", p2, err)
	}
	_, err = svc.ListarUsuarios(ctx, c.adm, adm.FiltroUsuarios{Limite: 2, Cursor: p1.ProximoCursor, Situacao: ptr("ativo")})
	exigirErro(t, err, erros.CursorInvalido)
}

// Isolamento entre empresas nos casos de uso: id de outra empresa é 404.
func TestOutraEmpresaNaoEncontra(t *testing.T) {
	a, b := novaConta(t, ""), novaConta(t, "")
	ctx := ctxWeb()
	if _, err := svc.ObterUsuario(ctx, b.adm, a.adm.UsuarioID); !errors.Is(err, erros.NaoEncontrado) {
		t.Errorf("usuário de outra empresa: %v", err)
	}
	if _, err := svc.InativarUsuario(ctx, b.adm, a.adm.UsuarioID); !errors.Is(err, erros.NaoEncontrado) {
		t.Errorf("inativar usuário de outra empresa: %v", err)
	}
	if _, err := svc.ConcederPermissao(ctx, b.adm, a.adm.UsuarioID, "ind.exportar"); !errors.Is(err, erros.NaoEncontrado) {
		t.Errorf("permissão a usuário de outra empresa: %v", err)
	}
	if n, err := svc.EncerrarSessoes(ctx, b.adm, a.adm.UsuarioID); !errors.Is(err, erros.NaoEncontrado) || n != 0 {
		t.Errorf("encerrar sessões de outra empresa: %v", err)
	}
	if _, err := auth.TrocarEmpresa(ctx, b.adm, a.criada.EmpresaID); !errors.Is(err, erros.NaoEncontrado) {
		t.Errorf("trocar para empresa sem acesso: %v", err)
	}
	if _, err := auth.AutenticarWeb(ctx, a.cookie, "", false); err != nil {
		t.Errorf("a sessão de A continua: %v", err)
	}
}

// RF-ADM-006: falhas, bloqueio temporário, tempo constante e trilha.
func TestLoginFalhasEBloqueio(t *testing.T) {
	c := novaConta(t, "")
	ctx := ctxWeb()
	_, err := auth.Login(ctx, autenticacao.EntradaLogin{Email: emailAleatorio("ninguem"), Senha: senhaBoa})
	exigirErro(t, err, erros.CredencialInvalida)
	for i := 0; i < 5; i++ {
		_, err = auth.Login(ctx, autenticacao.EntradaLogin{Email: c.email, Senha: "errada-errada-1"})
		exigirErro(t, err, erros.CredencialInvalida)
	}
	_, err = auth.Login(ctx, autenticacao.EntradaLogin{Email: c.email, Senha: senhaBoa})
	exigirErro(t, err, erros.CredencialInvalida) // bloqueada por 15 min
	ev, err := svc.ListarEventosAuditoria(ctx, c.adm, adm.FiltroAuditoria{Acao: ptr("login_falho")})
	if err != nil || len(ev.Dados) < 6 {
		t.Fatalf("falhas na trilha: %d %v", len(ev.Dados), err)
	}
	// Desbloqueio pelo tempo (simulado no banco).
	if _, err := amb.Migrador.Exec(context.Background(),
		`UPDATE plataforma.identidade SET bloqueada_ate = now() - interval '1 minute' WHERE email = $1`, c.email); err != nil {
		t.Fatal(err)
	}
	entrar(t, c.email, senhaBoa, nil)
}

// MFA por TOTP: ativação com confirmação, desafio com tentativas.
func TestMFA(t *testing.T) {
	c := novaConta(t, "")
	ctx := ctxWeb()
	cfg, err := auth.IniciarMFA(ctx, c.adm)
	if err != nil || !strings.HasPrefix(cfg.URIOTPAuth, "otpauth://totp/") {
		t.Fatalf("iniciar MFA: %v", err)
	}
	exigirErro(t, auth.ConfirmarMFA(ctx, c.adm, "000000"), erros.MFACodigoInvalido)
	if err := auth.ConfirmarMFA(ctx, c.adm, codigoTOTP(t, cfg.SegredoBase32, time.Now())); err != nil {
		t.Fatalf("confirmar MFA: %v", err)
	}
	r, err := auth.Login(ctx, autenticacao.EntradaLogin{Email: c.email, Senha: senhaBoa})
	if err != nil || r.DesafioMFA == nil || r.Sessao != nil {
		t.Fatalf("login com MFA devolve o desafio: %+v %v", r, err)
	}
	_, err = auth.VerificarMFA(ctx, autenticacao.EntradaVerificarMFA{DesafioID: r.DesafioMFA.DesafioID, Codigo: "000000"})
	var e *erros.Erro
	if !errors.As(err, &e) || e.Codigo != erros.MFACodigoInvalido.Codigo || e.Dados["tentativas_restantes"] != 4 {
		t.Fatalf("código errado consome uma tentativa: %v", err)
	}
	r2, err := auth.VerificarMFA(ctx, autenticacao.EntradaVerificarMFA{DesafioID: r.DesafioMFA.DesafioID,
		Codigo: codigoTOTP(t, cfg.SegredoBase32, time.Now())})
	if err != nil || r2.Sessao == nil {
		t.Fatalf("código certo conclui o login: %v", err)
	}
	_, err = auth.VerificarMFA(ctx, autenticacao.EntradaVerificarMFA{DesafioID: r.DesafioMFA.DesafioID,
		Codigo: codigoTOTP(t, cfg.SegredoBase32, time.Now())})
	exigirErro(t, err, erros.CredencialInvalida) // desafio é de uso único
}

// OAuth 2.1 do app: PKCE, rotação do refresh token e detecção de reuso.
func TestOAuthApp(t *testing.T) {
	c := novaConta(t, "")
	ctx := ctxWeb()
	verificador := strings.Repeat("v3rif1cad0r-", 5)
	h := sha256.Sum256([]byte(verificador))
	desafio := base64.RawURLEncoding.EncodeToString(h[:])
	entrada := autenticacao.EntradaAutorizacao{ClienteID: "tenaz-app", URIRetorno: "br.com.tenaz.app:/oauth",
		DesafioPKCE: desafio, MetodoPKCE: "S256", EmpresaID: c.criada.EmpresaID}
	_, err := auth.EmitirCodigoAutorizacao(ctx, c.adm, autenticacao.EntradaAutorizacao{ClienteID: "tenaz-app",
		URIRetorno: "https://malicioso/oauth", DesafioPKCE: desafio, MetodoPKCE: "S256", EmpresaID: c.criada.EmpresaID})
	exigirErro(t, err, erros.EntradaInvalida)

	codigo, err := auth.EmitirCodigoAutorizacao(ctx, c.adm, entrada)
	if err != nil {
		t.Fatal(err)
	}
	_, err = auth.TrocarCodigo(ctx, autenticacao.EntradaTrocaCodigo{Codigo: codigo, ClienteID: "tenaz-app",
		URIRetorno: "br.com.tenaz.app:/oauth", VerificadorPKCE: verificador + "x"})
	exigirErro(t, err, erros.CredencialInvalida) // PKCE errado gasta o código

	codigo, _ = auth.EmitirCodigoAutorizacao(ctx, c.adm, entrada)
	tok, err := auth.TrocarCodigo(ctx, autenticacao.EntradaTrocaCodigo{Codigo: codigo, ClienteID: "tenaz-app",
		URIRetorno: "br.com.tenaz.app:/oauth", VerificadorPKCE: verificador})
	if err != nil || tok.ExpiraEmSegs != 900 || !strings.HasPrefix(tok.RefreshToken, "tnz_rt_") {
		t.Fatalf("troca do código: %+v %v", tok, err)
	}
	a, err := auth.AutenticarApp(ctx, tok.AccessToken)
	if err != nil || a.EmpresaID != c.criada.EmpresaID || a.Canal != "app" {
		t.Fatalf("access token: %v", err)
	}
	novo, err := auth.Renovar(ctx, tok.RefreshToken, "tenaz-app")
	if err != nil || novo.RefreshToken == tok.RefreshToken {
		t.Fatalf("renovar: %v", err)
	}
	// Reuso do antigo: a família inteira cai, inclusive o novo.
	_, err = auth.Renovar(ctx, tok.RefreshToken, "tenaz-app")
	exigirErro(t, err, erros.CredencialInvalida)
	_, err = auth.Renovar(ctx, novo.RefreshToken, "tenaz-app")
	exigirErro(t, err, erros.CredencialInvalida)
	_, err = auth.AutenticarApp(ctx, novo.AccessToken)
	exigirErro(t, err, erros.CredencialInvalida)

	// Locatário não entra no app.
	_, err = auth.TrocarCodigo(ctx, autenticacao.EntradaTrocaCodigo{Codigo: codigo, ClienteID: "tenaz-app",
		URIRetorno: "br.com.tenaz.app:/oauth", VerificadorPKCE: verificador})
	exigirErro(t, err, erros.CredencialInvalida) // código reapresentado
	loc := uuid.New()
	_, aLoc := convidarEEntrar(t, c, adm.EntradaUsuario{Perfis: []string{"locatario"}, LocatarioID: &loc})
	entrada.EmpresaID = aLoc.EmpresaID
	_, err = auth.EmitirCodigoAutorizacao(ctx, aLoc, entrada)
	exigirErro(t, err, erros.SemPermissao)
}

// Escolha e troca de empresa: uma identidade em duas empresas.
func TestTrocaDeEmpresa(t *testing.T) {
	a, b := novaConta(t, ""), novaConta(t, "")
	ctx := ctxWeb()
	// O administrador de A convida o e-mail do administrador de B.
	if _, err := svc.CriarUsuario(ctx, a.adm, adm.EntradaUsuario{Nome: "Consultor", Email: b.email, Perfis: []string{"visualizador"}}); err != nil {
		t.Fatal(err)
	}
	aceito, err := svc.AceitarConvite(ctx, caixa.token(t, b.email), "")
	if err != nil || aceito.IdentidadeNova {
		t.Fatalf("convite para e-mail que já entra noutra empresa não pede senha: %+v %v", aceito, err)
	}
	r, err := auth.Login(ctx, autenticacao.EntradaLogin{Email: b.email, Senha: senhaBoa, EmpresaID: &b.criada.EmpresaID})
	if err != nil || len(r.Empresas) != 2 || r.Empresas[0].EmpresaID != b.criada.EmpresaID {
		t.Fatalf("login com duas empresas: %+v %v", r.Empresas, err)
	}
	atorB, err := auth.AutenticarWeb(ctx, r.Sessao.Cookie, r.Sessao.CSRF, true)
	if err != nil {
		t.Fatal(err)
	}
	troca, err := auth.TrocarEmpresa(ctx, atorB, a.criada.EmpresaID)
	if err != nil || troca.Web == nil {
		t.Fatalf("trocar empresa: %v", err)
	}
	if _, err := auth.AutenticarWeb(ctx, r.Sessao.Cookie, "", false); !errors.Is(err, erros.CredencialInvalida) {
		t.Errorf("a sessão anterior é revogada: %v", err)
	}
	atorA, err := auth.AutenticarWeb(ctx, troca.Web.Cookie, troca.Web.CSRF, true)
	if err != nil || atorA.EmpresaID != a.criada.EmpresaID || atorA.Pode(autorizacao.AdmGerenciarUsuarios) {
		t.Fatalf("na empresa A é só visualizador: %+v %v", atorA.Perfis, err)
	}
	// A trilha de A não cita a empresa B.
	ev, err := svc.ListarEventosAuditoria(ctx, a.adm, adm.FiltroAuditoria{Acao: ptr("login"), EntidadeID: &atorA.IdentidadeID})
	if err != nil || len(ev.Dados) == 0 {
		t.Fatalf("login na trilha de A: %v", err)
	}
	for _, e := range ev.Dados {
		det, err := svc.ObterEventoAuditoria(ctx, a.adm, e.ID)
		if err != nil || strings.Contains(string(det.Depois), b.criada.EmpresaID.String()) {
			t.Errorf("evento de A cita a empresa B: %s %v", det.Depois, err)
		}
	}
}

// Recuperação e redefinição de senha.
func TestRedefinirSenha(t *testing.T) {
	c := novaConta(t, "")
	ctx := ctxWeb()
	if err := auth.RecuperarSenha(ctx, emailAleatorio("ninguem")); err != nil {
		t.Fatal("e-mail inexistente responde igual")
	}
	if err := auth.RecuperarSenha(ctx, strings.ToUpper(c.email)); err != nil {
		t.Fatal(err)
	}
	token := caixa.token(t, c.email)
	exigirErro(t, auth.RedefinirSenha(ctx, token, "1234567890"), erros.SenhaFraca)
	if err := auth.RedefinirSenha(ctx, token, "Outra senha boa 2026!"); err != nil {
		t.Fatalf("redefinir: %v", err)
	}
	exigirErro(t, auth.RedefinirSenha(ctx, token, "Mais uma senha 2026!"), erros.LinkEmailInvalido)
	if _, err := auth.AutenticarWeb(ctx, c.cookie, "", false); !errors.Is(err, erros.CredencialInvalida) {
		t.Errorf("redefinir revoga as sessões: %v", err)
	}
	_, err := auth.Login(ctx, autenticacao.EntradaLogin{Email: c.email, Senha: senhaBoa})
	exigirErro(t, err, erros.CredencialInvalida)
	entrar(t, c.email, "Outra senha boa 2026!", nil)
}

// ADR-011: o pool da aplicação não alcança as credenciais; o da
// autenticação roda como tenaz_autenticacao.
func TestPoolsComPapeisSeparados(t *testing.T) {
	ctx := context.Background()
	var papel string
	if err := amb.Auth.QueryRow(ctx, `SELECT current_user`).Scan(&papel); err != nil || papel != "tenaz_autenticacao" {
		t.Fatalf("pool da autenticação: %q %v", papel, err)
	}
	c := novaConta(t, "")
	err := amb.App.NaEmpresa(ctx, c.criada.EmpresaID, c.adm.Escopo, func(ctx context.Context) error {
		tx, err := db.Tx(ctx)
		if err != nil {
			return err
		}
		var n int
		return tx.QueryRow(ctx, `SELECT count(*) FROM plataforma.identidade`).Scan(&n)
	})
	if err == nil || !strings.Contains(err.Error(), "permission denied") {
		t.Fatalf("tenaz_app lendo identidade: esperado permissão negada, veio %v", err)
	}
}
