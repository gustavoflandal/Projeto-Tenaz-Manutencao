package dominio

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/erros"
)

func TestCPF(t *testing.T) {
	casos := map[string]bool{
		"52998224725": true, "11144477735": true,
		"52998224724": false, "11111111111": false, "1234567890": false, "5299822472a": false, "": false,
	}
	for cpf, ok := range casos {
		if CPFValido(cpf) != ok {
			t.Errorf("CPFValido(%q) = %v", cpf, !ok)
		}
	}
	if d, e := NormalizarCPF("cpf", "529.982.247-25"); e != nil || d != "52998224725" {
		t.Errorf("normalizar com máscara: %q %v", d, e)
	}
	if MascararCPF("52998224725") != "***.982.247-**" {
		t.Errorf("máscara: %s", MascararCPF("52998224725"))
	}
}

func TestCNPJ(t *testing.T) {
	casos := map[string]bool{
		"11222333000181": true, "45723174000110": true,
		"11222333000180": false, "00000000000000": false, "1122233300018": false,
	}
	for cnpj, ok := range casos {
		if CNPJValido(cnpj) != ok {
			t.Errorf("CNPJValido(%q) = %v", cnpj, !ok)
		}
	}
	if _, e := NormalizarCNPJ("cnpj", "11.222.333/0001-81"); e != nil {
		t.Errorf("CNPJ com máscara: %v", e)
	}
}

func TestPlanoEUnidade(t *testing.T) {
	ess, pro := PlanoEssencial, PlanoProfissional
	if PodeCriarUnidade(&ess, 0) != nil {
		t.Error("Essencial cria a primeira unidade")
	}
	e := PodeCriarUnidade(&ess, 1)
	if e == nil || !errors.Is(e, erros.PlanoNaoIncluiRecurso) || e.Dados["plano_necessario"] != PlanoProfissional {
		t.Errorf("segunda unidade no Essencial: esperado PLANO_NAO_INCLUI_RECURSO com plano_necessario, veio %v", e)
	}
	if PodeCriarUnidade(nil, 1) == nil {
		t.Error("conta só do plano Serviços segue o Essencial")
	}
	if PodeCriarUnidade(&pro, 5) != nil {
		t.Error("Profissional cria mais unidades")
	}
	if PlanoValido(nil, false) == nil || PlanoValido(nil, true) != nil {
		t.Error("sem plano principal só com SRV")
	}
	x := "ouro"
	if PlanoValido(&x, false) == nil {
		t.Error("plano desconhecido")
	}
}

func TestTurno(t *testing.T) {
	validos := []Turno{{1, "08:00", "17:00"}, {7, "22:00", "06:00"}, {3, "00:00", "24:00"}}
	for _, v := range validos {
		if e := ValidarTurno("", v); e != nil {
			t.Errorf("%+v deveria valer: %v", v, e)
		}
	}
	invalidos := []Turno{{0, "08:00", "17:00"}, {8, "08:00", "17:00"}, {1, "8:00", "17:00"}, {1, "24:00", "06:00"},
		{1, "08:00", "08:00"}, {1, "08:00", "25:00"}, {1, "08:60", "17:00"}}
	for _, v := range invalidos {
		if ValidarTurno("turnos[0]", v) == nil {
			t.Errorf("%+v deveria ser recusado", v)
		}
	}
}

// RN-ADM-019: os casos de sobreposição da tabela do tenaz-dominio-manutencao.
func TestSobreposicaoDeTurnos(t *testing.T) {
	seg, ter, dom := int32(1), int32(2), int32(7)
	casos := []struct {
		nome   string
		a, b   Turno
		aceito bool
	}{
		{"seg 08-12 e 13-17", Turno{seg, "08:00", "12:00"}, Turno{seg, "13:00", "17:00"}, true},
		{"seg 08-12 e 12-17 (encostados)", Turno{seg, "08:00", "12:00"}, Turno{seg, "12:00", "17:00"}, true},
		{"seg 08-17 e 16-20", Turno{seg, "08:00", "17:00"}, Turno{seg, "16:00", "20:00"}, false},
		{"seg 22-06 e ter 05-14", Turno{seg, "22:00", "06:00"}, Turno{ter, "05:00", "14:00"}, false},
		{"seg 22-06 e ter 06-14", Turno{seg, "22:00", "06:00"}, Turno{ter, "06:00", "14:00"}, true},
		{"dom 22-06 e seg 05-08", Turno{dom, "22:00", "06:00"}, Turno{seg, "05:00", "08:00"}, false},
		{"seg 00:00-24:00 e seg 08-17", Turno{seg, "00:00", "24:00"}, Turno{seg, "08:00", "17:00"}, false},
		{"seg 18-02 e seg 01-03", Turno{seg, "18:00", "02:00"}, Turno{seg, "01:00", "03:00"}, true},
	}
	for _, c := range casos {
		for _, par := range [][2]Turno{{c.a, c.b}, {c.b, c.a}} {
			if got := !Sobrepoem(par[0], par[1]); got != c.aceito {
				t.Errorf("%s: aceito = %v, esperado %v", c.nome, got, c.aceito)
			}
		}
	}

	// Alterar um turno sem mexer no horário: ele é tirado da lista antes da
	// comparação, e o resto do calendário não conflita.
	calendario := []Turno{{seg, "08:00", "12:00"}, {seg, "13:00", "17:00"}}
	outros := calendario[1:]
	alterado := calendario[0]
	if i, _ := PrimeiraSobreposicao(append([]Turno{alterado}, outros...)); i != -1 {
		t.Error("alterar sem mexer no horário é aceito")
	}
	if i, j := PrimeiraSobreposicao([]Turno{{seg, "08:00", "17:00"}, {ter, "08:00", "17:00"}, {seg, "16:00", "18:00"}}); i != 0 || j != 2 {
		t.Errorf("primeira sobreposição: %d, %d", i, j)
	}
}

func TestTextosEFuso(t *testing.T) {
	if _, e := FusoHorario("fuso_horario", "America/Sao_Paulo"); e != nil {
		t.Errorf("fuso válido: %v", e)
	}
	if _, e := FusoHorario("fuso_horario", "Marte/Olimpo"); e == nil {
		t.Error("fuso inválido")
	}
	if v, e := Email("email", " Ana@Exemplo.com.br "); e != nil || v != "ana@exemplo.com.br" {
		t.Errorf("e-mail normalizado: %q %v", v, e)
	}
	if _, e := Email("email", "Ana <ana@x.com>"); e == nil {
		t.Error("e-mail com nome não é aceito")
	}
	if _, e := Data("data", "2026-02-30"); e == nil {
		t.Error("data inexistente")
	}
}

func TestPerfisDoUsuario(t *testing.T) {
	casos := []struct {
		nome   string
		perfis []string
		v      Vinculo
		ok     bool
	}{
		{"interno com dois perfis", []string{"gestor", "tecnico"}, Vinculo{}, true},
		{"sem perfil", nil, Vinculo{}, false},
		{"repetido", []string{"tecnico", "tecnico"}, Vinculo{}, false},
		{"prestador com fornecedor", []string{"prestador"}, Vinculo{Fornecedor: true}, true},
		{"prestador sem fornecedor", []string{"prestador"}, Vinculo{}, false},
		{"fornecedor com perfil interno", []string{"prestador", "tecnico"}, Vinculo{Fornecedor: true}, false},
		{"fornecedor com perfil locatario", []string{"locatario"}, Vinculo{Fornecedor: true}, false},
		{"locatario com locatário", []string{"locatario"}, Vinculo{Locatario: true}, true},
		{"interno com locatario", []string{"gestor", "locatario"}, Vinculo{}, false},
		{"os dois vínculos", []string{"prestador"}, Vinculo{Fornecedor: true, Locatario: true}, false},
	}
	for _, c := range casos {
		if got := ValidarPerfisDoUsuario(c.perfis, c.v) == nil; got != c.ok {
			t.Errorf("%s: esperado ok=%v", c.nome, c.ok)
		}
	}
}

func atorComPerfis(perfis ...string) autorizacao.Ator {
	var perms []string
	for _, p := range perfis {
		perms = append(perms, autorizacao.PermissoesDoPerfil(p)...)
	}
	id := uuid.New()
	return autorizacao.Ator{UsuarioID: id, Perfis: perfis, Permissoes: autorizacao.NovoConjunto(perms...),
		Escopo: db.Interno(&id)}
}

func TestAntiEscalada(t *testing.T) {
	adm := atorComPerfis(autorizacao.PerfilAdministrador)
	gestor := atorComPerfis(autorizacao.PerfilGestor)
	alvo := uuid.New()
	perfil := func(cods ...string) map[string][]string {
		m := map[string][]string{}
		for _, c := range cods {
			m[c] = autorizacao.PermissoesDoPerfil(c)
		}
		return m
	}

	if e := ConferirAtribuicao(adm, alvo, false, perfil("gestor", "tecnico")); e != nil {
		t.Errorf("administrador atribui gestor (com os.encerrar): %v", e)
	}
	if e := ConferirAtribuicao(adm, adm.UsuarioID, true, perfil("administrador")); e == nil {
		t.Error("ninguém altera os próprios perfis")
	}
	if e := ConferirAtribuicao(gestor, alvo, false, perfil("administrador")); e == nil {
		t.Error("gestor não dá o perfil administrador")
	}
	if e := ConferirAtribuicao(gestor, alvo, true, perfil("tecnico")); e == nil {
		t.Error("só administrador altera administrador")
	}

	// Permissões avulsas.
	if e := ConferirConcessao(adm, alvo, autorizacao.PessoaSaudeLer, ""); e != nil {
		t.Errorf("administrador concede pessoa.saude.ler a outro: %v", e)
	}
	if e := ConferirConcessao(adm, adm.UsuarioID, autorizacao.OsEncerrar, ""); e == nil {
		t.Error("o administrador nunca concede a si mesmo")
	}
	if e := ConferirConcessao(gestor, alvo, autorizacao.OsEncerrar, ""); e == nil {
		t.Error("só o administrador concede os.encerrar, mesmo quem a tem")
	}
	if e := ConferirConcessao(adm, alvo, autorizacao.OsEncerrar, "prestador"); e == nil || !errors.Is(e, erros.EntradaInvalida) {
		t.Errorf("os.encerrar a usuário externo: esperado ENTRADA_INVALIDA, veio %v", e)
	}
	if e := ConferirConcessao(adm, alvo, "os.executar", "prestador"); e != nil {
		t.Errorf("prestador recebe avulsa do conjunto do seu perfil: %v", e)
	}
	if e := ConferirConcessao(adm, alvo, "ind.ler", "prestador"); !errors.Is(e, erros.EntradaInvalida) {
		t.Errorf("prestador não recebe avulsa fora do conjunto do seu perfil (RN-ADM-026): %v", e)
	}
	if e := ConferirConcessao(adm, alvo, "os.voar", ""); e == nil {
		t.Error("permissão fora do catálogo")
	}
	if e := ConferirConcessao(gestor, alvo, "ind.exportar", ""); e != nil {
		t.Errorf("gestor concede permissão que tem: %v", e)
	}
	if e := ConferirConcessao(gestor, alvo, "int.gerenciar_tokens", ""); e == nil {
		t.Error("gestor não concede o que não tem")
	}
}
