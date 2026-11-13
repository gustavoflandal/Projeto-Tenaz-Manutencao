package autorizacao

import (
	"bufio"
	"errors"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/erros"
)

// caminhoDocumento1 parte do diretório do pacote até a raiz do repositório.
var caminhoDocumento1 = filepath.Join("..", "..", "..", "..", "docs", "1_ESPECIFICACAO_REQUISITOS.md")

type linhaDoc struct {
	codigo, descricao, grupo, fase string
	sensivel                       bool
}

// lerTabelasDoDocumento relê o catálogo e a matriz de RF-ADM-005.
func lerTabelasDoDocumento(t *testing.T) ([]linhaDoc, map[string][]bool) {
	t.Helper()
	f, err := os.Open(caminhoDocumento1)
	if err != nil {
		t.Fatalf("abrir o documento 1: %v", err)
	}
	defer f.Close()

	var cat []linhaDoc
	mat := map[string][]bool{}
	modo := ""
	s := bufio.NewScanner(f)
	s.Buffer(make([]byte, 1<<20), 1<<20)
	for s.Scan() {
		l := s.Text()
		switch {
		case strings.HasPrefix(l, "**Catálogo**"):
			modo = "cat"
			continue
		case strings.HasPrefix(l, "**Matriz dos perfis padrão**"):
			modo = "mat"
			continue
		case strings.HasPrefix(l, "**Totais"):
			modo = ""
		}
		if !strings.HasPrefix(l, "| `") {
			continue
		}
		cols := strings.Split(strings.Trim(strings.TrimSpace(l), "|"), "|")
		for i := range cols {
			cols[i] = strings.TrimSpace(cols[i])
		}
		cod := strings.Trim(cols[0], "`")
		switch modo {
		case "cat":
			cat = append(cat, linhaDoc{cod, cols[1], cols[2], cols[3], cols[4] == "sim"})
		case "mat":
			var m []bool
			for _, c := range cols[1:13] {
				m = append(m, c == "●")
			}
			mat[cod] = m
		}
	}
	if err := s.Err(); err != nil {
		t.Fatal(err)
	}
	return cat, mat
}

// O catálogo e a matriz do código são iguais aos de RF-ADM-005: mesmo código,
// descrição, grupo, fase e marcação de sensível, e os mesmos ● por perfil.
func TestCatalogoIgualAoDocumento(t *testing.T) {
	cat, mat := lerTabelasDoDocumento(t)
	if len(cat) != 116 {
		t.Fatalf("o documento tem %d permissões no catálogo; o sumário 3.6 cita 116", len(cat))
	}
	codigo := Catalogo()
	if len(codigo) != len(cat) {
		t.Fatalf("código tem %d permissões, documento tem %d", len(codigo), len(cat))
	}
	for i, d := range cat {
		c := codigo[i]
		if c.Codigo != d.codigo || c.Descricao != d.descricao || c.Grupo != d.grupo ||
			string(c.Fase) != d.fase || c.Sensivel != d.sensivel {
			t.Errorf("linha %d diverge:\n código:    %+v\n documento: %+v", i, c, d)
		}
	}
	for i, p := range PerfisPadrao() {
		var doc []string
		for _, d := range cat {
			if mat[d.codigo][i] {
				doc = append(doc, d.codigo)
			}
		}
		if got := PermissoesDoPerfil(p.Codigo); !slices.Equal(got, doc) {
			t.Errorf("perfil %s diverge do documento:\n código:    %v\n documento: %v", p.Codigo, got, doc)
		}
	}
}

// Totais por perfil do documento (linha "Totais por perfil padrão").
func TestTotaisPorPerfil(t *testing.T) {
	esperado := map[string]int{
		PerfilAdministrador: 113, PerfilGestor: 102, PerfilPlanejador: 87, PerfilAlmoxarife: 41,
		PerfilTecnico: 38, PerfilComprador: 32, PerfilPLH: 24, PerfilVisualizador: 23,
		PerfilPrestador: 18, PerfilMotorista: 14, PerfilLocatario: 11, PerfilSolicitante: 6,
	}
	for perfil, n := range esperado {
		if got := len(PermissoesDoPerfil(perfil)); got != n {
			t.Errorf("%s: %d permissões, documento diz %d", perfil, got, n)
		}
	}
}

// Regras do sumário executivo (seções 3.1 e 3.6) sobre a matriz.
func TestMatrizRespeitaRegrasDoSumario(t *testing.T) {
	quemTem := func(p string) []string {
		var l []string
		for _, perfil := range PerfisPadrao() {
			if slices.Contains(PermissoesDoPerfil(perfil.Codigo), p) {
				l = append(l, perfil.Codigo)
			}
		}
		return l
	}

	if got := quemTem(OsEncerrar); !slices.Equal(got, []string{PerfilGestor, PerfilPlanejador}) {
		t.Errorf("os.encerrar deve estar só em gestor e planejador; está em %v", got)
	}
	if got := quemTem(EstDecidirDestino); !slices.Equal(got, []string{PerfilGestor, PerfilPlanejador, PerfilAlmoxarife}) {
		t.Errorf("est.decidir_destino deve estar só em gestor, planejador e almoxarife; está em %v", got)
	}
	// A Tenaz não guarda dado de saúde: a permissão não existe mais.
	for _, p := range Catalogo() {
		if strings.Contains(p.Codigo, "saude") {
			t.Errorf("permissão de saúde no catálogo: %s", p.Codigo)
		}
	}
	if got := quemTem(SrvGerenciarConexao); !slices.Equal(got, []string{PerfilAdministrador, PerfilGestor}) {
		t.Errorf("srv.gerenciar_conexao deve ser padrão de administrador e gestor; está em %v", got)
	}
	if slices.Contains(PermissoesDoPerfil(PerfilLocatario), MobUsarApp) {
		t.Error("locatario não entra no app (sumário 3.6, RF-MOB-001)")
	}

	// Administrador: todas, menos três.
	adm := PermissoesDoPerfil(PerfilAdministrador)
	for _, p := range Catalogo() {
		fora := p.Codigo == OsEncerrar || p.Codigo == EstDecidirDestino || p.Codigo == TerManterCatalogoProprio
		if slices.Contains(adm, p.Codigo) == fora {
			t.Errorf("administrador e %s: esperado ter = %v", p.Codigo, !fora)
		}
	}

	// Nenhuma vedada em perfil externo.
	for _, externo := range []string{PerfilPrestador, PerfilLocatario} {
		for _, p := range PermissoesDoPerfil(externo) {
			if VedadaExterno(p) {
				t.Errorf("%s tem a permissão vedada %s", externo, p)
			}
		}
	}
}

func TestCatalogoBemFormado(t *testing.T) {
	vistos := map[string]bool{}
	for _, p := range Catalogo() {
		if vistos[p.Codigo] {
			t.Errorf("código repetido: %s", p.Codigo)
		}
		vistos[p.Codigo] = true
		if !FormatoValido(p.Codigo) {
			t.Errorf("formato inválido: %s", p.Codigo)
		}
	}
	for perfil, lista := range matriz {
		if _, ok := BuscarPerfilPadrao(perfil); !ok {
			t.Errorf("perfil desconhecido na matriz: %s", perfil)
		}
		for _, p := range lista {
			if !vistos[p] {
				t.Errorf("%s tem permissão fora do catálogo: %s", perfil, p)
			}
		}
	}
	if len(PerfisPadrao()) != 12 {
		t.Fatalf("são 12 perfis padrão")
	}
}

func TestFormatoValido(t *testing.T) {
	casos := map[string]bool{
		"os.encerrar": true, "pessoa.saude.ler": false, "adm.ler": true, "est.decidir_destino": true,
		"OS.encerrar": false, "os": false, "os.": false, "o.ler": false, "abcd.ler": false,
		"os._ler": false, "os.ler_": false, "con.gerenciar_nr13": true, "os.1ler": false, "os.ler-x": false, "pessoa.saude": false, "": false,
	}
	for c, esperado := range casos {
		if got := FormatoValido(c); got != esperado {
			t.Errorf("FormatoValido(%q) = %v, esperado %v", c, got, esperado)
		}
	}
}

func TestVedadasExterno(t *testing.T) {
	v := VedadasExterno()
	if len(v) != 25 {
		t.Errorf("esperadas 25 vedadas (1 + 9 + srv + 14 do grupo Admin); são %d: %v", len(v), v)
	}
	for _, p := range []string{OsEncerrar, SrvGerenciarConexao, AdmLerAuditoria, AdmGerenciarUsuarios,
		CtrReabrirMedicao, "ia.gerenciar", "con.gerenciar_tipos_documento"} {
		if !slices.Contains(v, p) {
			t.Errorf("%s deveria ser vedada a externo", p)
		}
	}
	if VedadaExterno("os.executar") {
		t.Error("os.executar não é vedada: o prestador executa OS")
	}
}

func TestExigirEAntiEscalada(t *testing.T) {
	// O erro é conferido pela identidade (errors.Is), não pelo texto: a
	// mensagem pode mudar; o código estável, não.
	interno := Ator{Escopo: db.Escopo{Tipo: db.EscopoInterno}, Permissoes: NovoConjunto(AdmLer)}
	if err := interno.Exigir(AdmLer); err != nil {
		t.Fatalf("tem adm.ler: %v", err)
	}
	if err := interno.ExigirAlguma(AdmGerenciarEmpresa, AdmLer); err != nil {
		t.Fatalf("tem uma das permissões: %v", err)
	}
	for nome, err := range map[string]error{
		"Exigir":        interno.Exigir(AdmLer, AdmGerenciarEmpresa),
		"ExigirAlguma":  interno.ExigirAlguma(AdmGerenciarEmpresa, AdmGerenciarUsuarios),
		"sem permissão": (Ator{Escopo: db.Escopo{Tipo: db.EscopoInterno}}).Exigir(AdmLer),
	} {
		if !errors.Is(err, erros.SemPermissao) || errors.Is(err, erros.NaoEncontrado) {
			t.Errorf("%s: interno sem a permissão recebe SemPermissao, veio %v", nome, err)
		}
	}
	for _, tipo := range []string{db.EscopoPrestador, db.EscopoLocatario} {
		externo := Ator{Escopo: db.Escopo{Tipo: tipo}, Permissoes: NovoConjunto("os.ler")}
		for nome, err := range map[string]error{
			"Exigir":       externo.Exigir(AdmLerPessoas),
			"ExigirAlguma": externo.ExigirAlguma(AdmLerPessoas, AdmLerUsuarios),
		} {
			if !errors.Is(err, erros.NaoEncontrado) || errors.Is(err, erros.SemPermissao) {
				t.Errorf("%s %s: externo sem a permissão recebe NaoEncontrado, veio %v", tipo, nome, err)
			}
		}
		if err := externo.Exigir("os.ler"); err != nil {
			t.Errorf("%s com a permissão: %v", tipo, err)
		}
	}

	adm := Ator{Perfis: []string{PerfilAdministrador}, Permissoes: NovoConjunto(PermissoesDoPerfil(PerfilAdministrador)...)}
	for _, p := range []string{OsEncerrar, EstDecidirDestino} {
		if !adm.PodeConceder(p) {
			t.Errorf("administrador concede %s a outros", p)
		}
	}
	if adm.PodeConceder(TerManterCatalogoProprio) {
		t.Error("ter.manter_catalogo_proprio não é isenção do administrador na concessão avulsa")
	}
	for _, perfil := range PerfisPadrao() {
		if ok, falta := adm.PodeAtribuirPerfil(PermissoesDoPerfil(perfil.Codigo)); !ok {
			t.Errorf("administrador deveria atribuir %s; falta %s", perfil.Codigo, falta)
		}
	}
	gestor := Ator{Perfis: []string{PerfilGestor}, Permissoes: NovoConjunto(PermissoesDoPerfil(PerfilGestor)...)}
	if gestor.PodeConceder(AdmGerenciarUsuarios) || IsencaoAdministrador(OsEncerrar) && !gestor.PodeConceder(OsEncerrar) {
		t.Error("gestor concede o que tem (os.encerrar) e não o que não tem (adm.gerenciar_usuarios)")
	}
	if ok, _ := gestor.PodeAtribuirPerfil(PermissoesDoPerfil(PerfilAdministrador)); ok {
		t.Error("gestor não atribui o perfil administrador")
	}
}
