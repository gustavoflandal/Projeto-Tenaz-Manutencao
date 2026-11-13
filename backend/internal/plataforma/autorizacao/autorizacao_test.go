package autorizacao

import (
	"bufio"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"tenaz/internal/plataforma/db"
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
	if len(cat) != 117 {
		t.Fatalf("o documento tem %d permissões no catálogo; o sumário 3.6 cita 117", len(cat))
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
	if got := quemTem(PessoaSaudeLer); len(got) != 0 {
		t.Errorf("pessoa.saude.ler não pode estar em perfil padrão; está em %v", got)
	}
	if got := quemTem(SrvGerenciarConexao); !slices.Equal(got, []string{PerfilAdministrador, PerfilGestor}) {
		t.Errorf("srv.gerenciar_conexao deve ser padrão de administrador e gestor; está em %v", got)
	}
	if slices.Contains(PermissoesDoPerfil(PerfilLocatario), MobUsarApp) {
		t.Error("locatario não entra no app (sumário 3.6, RF-MOB-001)")
	}

	// Administrador: todas, menos quatro.
	adm := PermissoesDoPerfil(PerfilAdministrador)
	for _, p := range Catalogo() {
		fora := p.Codigo == PessoaSaudeLer || p.Codigo == OsEncerrar || p.Codigo == EstDecidirDestino ||
			p.Codigo == TerManterCatalogoProprio
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
		"os.encerrar": true, "pessoa.saude.ler": true, "adm.ler": true, "est.decidir_destino": true,
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
	if len(v) != 26 {
		t.Errorf("esperadas 26 vedadas (1 + 9 + srv + saúde + 14 do grupo Admin); são %d: %v", len(v), v)
	}
	for _, p := range []string{OsEncerrar, PessoaSaudeLer, SrvGerenciarConexao, AdmLerAuditoria, AdmGerenciarUsuarios,
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
	interno := Ator{Escopo: db.Escopo{Tipo: db.EscopoInterno}, Permissoes: NovoConjunto(AdmLer)}
	if err := interno.Exigir(AdmLer); err != nil {
		t.Fatalf("tem adm.ler: %v", err)
	}
	if err := interno.Exigir(AdmGerenciarEmpresa); err == nil || !strings.Contains(err.Error(), "permissao_negada") {
		t.Errorf("interno sem a permissão: esperado permissao_negada, veio %v", err)
	}
	externo := Ator{Escopo: db.Escopo{Tipo: db.EscopoPrestador}, Permissoes: NovoConjunto("os.ler")}
	if err := externo.Exigir(AdmLerPessoas); err == nil || !strings.Contains(err.Error(), "recurso_nao_encontrado") {
		t.Errorf("externo sem a permissão: esperado recurso_nao_encontrado, veio %v", err)
	}

	adm := Ator{Perfis: []string{PerfilAdministrador}, Permissoes: NovoConjunto(PermissoesDoPerfil(PerfilAdministrador)...)}
	for _, p := range []string{OsEncerrar, EstDecidirDestino, PessoaSaudeLer} {
		if !adm.PodeConceder(p) {
			t.Errorf("administrador concede %s a outros", p)
		}
	}
	for _, perfil := range PerfisPadrao() {
		if ok, falta := adm.PodeAtribuirPerfil(PermissoesDoPerfil(perfil.Codigo)); !ok {
			t.Errorf("administrador deveria atribuir %s; falta %s", perfil.Codigo, falta)
		}
	}
	gestor := Ator{Perfis: []string{PerfilGestor}, Permissoes: NovoConjunto(PermissoesDoPerfil(PerfilGestor)...)}
	if gestor.PodeConceder(PessoaSaudeLer) {
		t.Error("gestor não concede pessoa.saude.ler")
	}
	if ok, _ := gestor.PodeAtribuirPerfil(PermissoesDoPerfil(PerfilAdministrador)); ok {
		t.Error("gestor não atribui o perfil administrador")
	}
}
