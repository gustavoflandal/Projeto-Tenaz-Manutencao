package pais_test

import (
	"encoding/json"
	"errors"
	"os"
	"slices"
	"testing"
	"time"
	_ "time/tzdata"

	"tenaz/internal/plataforma/pais"
)

// Os casos são os mesmos que pacotes/dominio (TypeScript) executa: um arquivo,
// duas suítes, para a web e o servidor não divergirem.
const arquivoCasos = "../../../testes/casos/identificacao.json"

type casos struct {
	Versao       int
	Normalizacao []struct{ Nome, Entrada, Esperado string }
	Validacao    []struct {
		Nome, Tipo, Entrada, Normalizado string
		Valido                           bool
		Motivo                           string
	}
	Mascara     []struct{ Nome, Tipo, Entrada, Esperado string }
	Subdivisoes []struct {
		Pais  string
		Lista []string
	}
	SubdivisaoInvalida []string `json:"subdivisao_invalida"`
	CodigoPostal       []struct {
		Nome, Pais, Entrada, Normalizado string
		Valido                           bool
	} `json:"codigo_postal"`
	CodigoPostalCPA []struct {
		Nome, Subdivisao string
		CodigoPostal     string `json:"codigo_postal"`
		Confere          bool
	} `json:"codigo_postal_cpa"`
	Perfis []struct {
		Pais         string
		Moeda        string
		IdiomaPadrao string   `json:"idioma_padrao"`
		FusoPadrao   string   `json:"fuso_padrao"`
		Fusos        []string `json:"fusos"`
		TipoEmpresa  string   `json:"tipo_empresa"`
		TiposPessoa  []string `json:"tipos_pessoa"`
	}
}

func carregar(t *testing.T) casos {
	t.Helper()
	b, err := os.ReadFile(arquivoCasos)
	if err != nil {
		t.Fatalf("ler casos: %v", err)
	}
	var c casos
	if err := json.Unmarshal(b, &c); err != nil {
		t.Fatalf("casos inválidos: %v", err)
	}
	return c
}

func TestNormalizar(t *testing.T) {
	for _, c := range carregar(t).Normalizacao {
		if got := pais.Normalizar(c.Entrada); got != c.Esperado {
			t.Errorf("%s: Normalizar(%q) = %q, esperado %q", c.Nome, c.Entrada, got, c.Esperado)
		}
	}
}

func TestValidarIdentificacao(t *testing.T) {
	motivos := map[string]error{
		"formato":            pais.ErrFormato,
		"digito_verificador": pais.ErrDigitoVerificador,
		"tipo_desconhecido":  pais.ErrTipoDesconhecido,
	}
	for _, c := range carregar(t).Validacao {
		got, err := pais.ValidarIdentificacao(c.Tipo, c.Entrada)
		if c.Valido {
			if err != nil || got != c.Normalizado {
				t.Errorf("%s: = (%q, %v), esperado (%q, nil)", c.Nome, got, err, c.Normalizado)
			}
			continue
		}
		if !errors.Is(err, motivos[c.Motivo]) {
			t.Errorf("%s: erro %v, esperado %v", c.Nome, err, motivos[c.Motivo])
		}
	}
}

func TestMascarar(t *testing.T) {
	for _, c := range carregar(t).Mascara {
		if got := pais.Mascarar(c.Tipo, c.Entrada); got != c.Esperado {
			t.Errorf("%s: Mascarar = %q, esperado %q", c.Nome, got, c.Esperado)
		}
	}
}

func TestPerfis(t *testing.T) {
	c := carregar(t)
	if got := pais.Codigos(); len(got) != len(c.Perfis) {
		t.Fatalf("países: %v, esperados %d", got, len(c.Perfis))
	}
	for _, e := range c.Perfis {
		p, ok := pais.Perfil(e.Pais)
		if !ok {
			t.Fatalf("país %s sem perfil", e.Pais)
		}
		if p.Moeda.Codigo() != e.Moeda || p.IdiomaPadrao != e.IdiomaPadrao || p.FusoPadrao != e.FusoPadrao || !slices.Equal(p.Fusos, e.Fusos) || p.TipoEmpresa != e.TipoEmpresa ||
			!slices.Equal(p.TiposPessoa, e.TiposPessoa) {
			t.Errorf("perfil %s = %+v, esperado %+v", e.Pais, p, e)
		}
	}
	if _, ok := pais.Perfil("US"); ok {
		t.Error("US não é suportado")
	}
	if !slices.Equal(pais.TiposPessoaDeTodos(), []string{"cpf", "cuil", "dni", "ci"}) {
		t.Errorf("tipos de pessoa de todos: %v", pais.TiposPessoaDeTodos())
	}
}

// O perfil devolvido é cópia: quem altera a lista não muda o perfil.
func TestPerfilEhCopia(t *testing.T) {
	p, _ := pais.Perfil(pais.AR)
	p.TiposPessoa[0] = "x"
	q, _ := pais.Perfil(pais.AR)
	if q.TiposPessoa[0] != pais.TipoCUIL {
		t.Errorf("perfil alterado por fora: %v", q.TiposPessoa)
	}
}

// Todo fuso da lista existe na base tz, o padrão está na lista, e o que não é
// do país não vale.
func TestFusosDoPais(t *testing.T) {
	for _, cod := range pais.Codigos() {
		p, _ := pais.Perfil(cod)
		for _, f := range p.Fusos {
			if _, err := time.LoadLocation(f); err != nil {
				t.Errorf("%s: fuso %s fora da base tz: %v", cod, f, err)
			}
		}
		if !pais.FusoValido(p, p.FusoPadrao) {
			t.Errorf("%s: o fuso padrão %s não está na lista", cod, p.FusoPadrao)
		}
	}
	br, _ := pais.Perfil(pais.BR)
	for _, f := range []string{"America/Argentina/Buenos_Aires", "America/Asuncion", "UTC", "Europe/Paris", "", "america/sao_paulo", "-03:00"} {
		if pais.FusoValido(br, f) {
			t.Errorf("BR aceitou %q", f)
		}
	}
}

// As listas de subdivisão do perfil são as do arquivo de casos (27, 24 e 18) e
// o que não está nelas, inclusive PY com zero à esquerda, não vale.
func TestSubdivisoes(t *testing.T) {
	c := carregar(t)
	if len(c.Subdivisoes) != len(pais.Codigos()) {
		t.Fatalf("listas no arquivo: %d", len(c.Subdivisoes))
	}
	tamanhos := map[string]int{pais.BR: 27, pais.AR: 24, pais.PY: 18}
	for _, e := range c.Subdivisoes {
		p, ok := pais.Perfil(e.Pais)
		if !ok {
			t.Fatalf("país %s sem perfil", e.Pais)
		}
		if !slices.Equal(p.Subdivisoes, e.Lista) || len(p.Subdivisoes) != tamanhos[e.Pais] {
			t.Errorf("%s: subdivisões %v, esperadas %v", e.Pais, p.Subdivisoes, e.Lista)
		}
		for _, s := range e.Lista {
			if !pais.SubdivisaoValida(p, s) {
				t.Errorf("%s: %s recusada", e.Pais, s)
			}
		}
		for _, s := range c.SubdivisaoInvalida {
			if pais.SubdivisaoValida(p, s) {
				t.Errorf("%s: %q aceita", e.Pais, s)
			}
		}
	}
	// o código de um país não vale em outro
	br, _ := pais.Perfil(pais.BR)
	if pais.SubdivisaoValida(br, "AR-C") || pais.SubdivisaoValida(br, "PY-ASU") {
		t.Error("BR aceitou subdivisão de outro país")
	}
}

func TestValidarCodigoPostal(t *testing.T) {
	for _, c := range carregar(t).CodigoPostal {
		p, ok := pais.Perfil(c.Pais)
		if !ok {
			t.Fatalf("%s: país %s sem perfil", c.Nome, c.Pais)
		}
		got, err := pais.ValidarCodigoPostal(p, c.Entrada)
		if c.Valido {
			if err != nil || got != c.Normalizado {
				t.Errorf("%s: = (%q, %v), esperado (%q, nil)", c.Nome, got, err, c.Normalizado)
			}
			continue
		}
		if !errors.Is(err, pais.ErrCodigoPostalFormato) {
			t.Errorf("%s: erro %v, esperado ErrCodigoPostalFormato", c.Nome, err)
		}
	}
}

// Uma regra só nas duas suítes: o prefixo da subdivisão decide e o código postal
// chega como a pessoa digitou (a função normaliza).
func TestCPAConfereComSubdivisao(t *testing.T) {
	for _, c := range carregar(t).CodigoPostalCPA {
		if got := pais.CPAConfereComSubdivisao(c.Subdivisao, c.CodigoPostal); got != c.Confere {
			t.Errorf("%s: = %v, esperado %v", c.Nome, got, c.Confere)
		}
	}
}

// A trava do arquivo de casos: se um grupo for renomeado no JSON, o Go não pode
// rodar zero casos e passar. Versão e tamanho de cada grupo são os que o
// TypeScript também espera; mudar o arquivo exige mudar os dois lados.
func TestCasosTravados(t *testing.T) {
	c := carregar(t)
	if c.Versao != 3 {
		t.Fatalf("versão dos casos %d, esperada 3", c.Versao)
	}
	for nome, g := range map[string][2]int{
		"normalizacao": {len(c.Normalizacao), 7}, "validacao": {len(c.Validacao), 34}, "mascara": {len(c.Mascara), 6},
		"perfis": {len(c.Perfis), 3}, "subdivisoes": {len(c.Subdivisoes), 3}, "subdivisao_invalida": {len(c.SubdivisaoInvalida), 12},
		"codigo_postal": {len(c.CodigoPostal), 16}, "codigo_postal_cpa": {len(c.CodigoPostalCPA), 12},
	} {
		if g[0] == 0 || g[0] != g[1] {
			t.Errorf("grupo %s: %d casos, esperados %d (e nenhum grupo pode ficar vazio)", nome, g[0], g[1])
		}
	}
}
