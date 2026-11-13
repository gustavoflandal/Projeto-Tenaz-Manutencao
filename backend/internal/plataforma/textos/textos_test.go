package textos_test

import (
	"net/http"
	"strings"
	"testing"
	"testing/fstest"

	"tenaz/internal/plataforma/textos"
)

// TestConferirCatalogos é a tarefa `textos` do CI: chave faltando ou sobrando,
// variável diferente, plural incompleto e termo proibido quebram o build.
func TestConferirCatalogos(t *testing.T) {
	if p := textos.Conferir(textos.Padrao()); len(p) > 0 {
		t.Fatalf("catálogos com problema:\n  %s", strings.Join(p, "\n  "))
	}
}

func catalogoDe(t *testing.T, pt, es string) *textos.Catalogo {
	t.Helper()
	c, err := textos.Novo(fstest.MapFS{
		"pt-BR.json":  {Data: []byte(pt)},
		"es-419.json": {Data: []byte(es)},
	})
	if err != nil {
		t.Fatal(err)
	}
	return c
}

func TestConferirApanhaOsDefeitos(t *testing.T) {
	c := catalogoDe(t,
		`{"a":{"b":"Olá, {{nome}}."},"so_pt":"x","p_one":"um","p_other":"{{n}}","q_other":"{{n}}"}`,
		`{"a":{"b":"Hola, {{nombre}}."},"so_es":"y","p_one":"uno","p_other":"{{n}}","q_other":"{{n}} link"}`)
	problemas := strings.Join(textos.Conferir(c), "\n")
	for _, esperado := range []string{
		`es-419: falta a chave "so_pt"`,
		`es-419: a chave "so_es" sobra`,
		`a chave "a.b" usa as variáveis [nombre]`,
		`o plural "q" não tem a forma one`,
		`usa o termo proibido "link"`,
	} {
		if !strings.Contains(problemas, esperado) {
			t.Errorf("não apanhou %q em:\n%s", esperado, problemas)
		}
	}
}

func TestInterpolacaoEBusca(t *testing.T) {
	c := catalogoDe(t, `{"s":"Olá, {{nome}}! {{x}}","so_pt":"só pt"}`, `{"s":"Hola, {{nome}}! {{x}}"}`)
	if got := c.T(textos.Es419, "s", textos.Parametros{"nome": "Ana"}); got != "Hola, Ana! {{x}}" {
		t.Errorf("T es-419 = %q", got)
	}
	// Chave sem texto no idioma cai para pt-BR; sem texto em nenhum, a chave.
	if got := c.T(textos.Es419, "so_pt", nil); got != "só pt" {
		t.Errorf("queda para pt-BR = %q", got)
	}
	if got := c.T(textos.Es419, "nada", nil); got != "nada" {
		t.Errorf("chave ausente = %q", got)
	}
}

func TestPluralPelaRegraCLDR(t *testing.T) {
	c := catalogoDe(t, `{"k_one":"{{n}} item","k_other":"{{n}} itens"}`, `{"k_one":"{{n}} elemento","k_other":"{{n}} elementos"}`)
	casos := []struct {
		idioma textos.Idioma
		n      int
		quer   string
	}{
		// Em português, zero usa a forma _one; em espanhol, _other.
		{textos.PtBR, 0, "0 item"}, {textos.PtBR, 1, "1 item"}, {textos.PtBR, 2, "2 itens"},
		{textos.Es419, 0, "0 elementos"}, {textos.Es419, 1, "1 elemento"}, {textos.Es419, 2, "2 elementos"},
	}
	for _, x := range casos {
		if got := c.Plural(x.idioma, "k", x.n, nil); got != x.quer {
			t.Errorf("Plural(%s, %d) = %q, esperado %q", x.idioma, x.n, got, x.quer)
		}
	}
}

func TestAcceptLanguage(t *testing.T) {
	casos := []struct {
		cab  string
		quer textos.Idioma
		ok   bool
	}{
		{"", "", false},
		{"pt", textos.PtBR, true},
		{"pt-PT", textos.PtBR, true},
		{"es", textos.Es419, true},
		{"es-AR", textos.Es419, true},
		{"es-PY,es;q=0.9", textos.Es419, true},
		{"en-US,en;q=0.9", "", false},
		{"en-US,es;q=0.5,pt;q=0.4", textos.Es419, true},
		{"pt;q=0.3,es;q=0.8", textos.Es419, true},
		{"es;q=0", "", false},
		{"*", "", false},
		{"lixo;;;", "", false},
	}
	for _, x := range casos {
		got, ok := textos.DeAcceptLanguage(x.cab)
		if got != x.quer || ok != x.ok {
			t.Errorf("DeAcceptLanguage(%q) = (%q, %v), esperado (%q, %v)", x.cab, got, ok, x.quer, x.ok)
		}
	}
}

func TestOrdemDaNegociacao(t *testing.T) {
	casos := []struct {
		nome string
		p    textos.Preferencias
		quer textos.Idioma
	}{
		{"nada: pt-BR", textos.Preferencias{}, textos.PtBR},
		{"cabeçalho vence o usuário", textos.Preferencias{AcceptLanguage: "es", Usuario: "pt-BR", Empresa: "pt-BR"}, textos.Es419},
		{"sem casamento, o usuário", textos.Preferencias{AcceptLanguage: "en", Usuario: "es-419", Empresa: "pt-BR"}, textos.Es419},
		{"sem usuário, a empresa", textos.Preferencias{AcceptLanguage: "en", Empresa: "es-419"}, textos.Es419},
		{"usuário inválido cai para a empresa", textos.Preferencias{Usuario: "fr", Empresa: "es-419"}, textos.Es419},
		{"tudo inválido: pt-BR", textos.Preferencias{AcceptLanguage: "fr", Usuario: "fr", Empresa: "fr"}, textos.PtBR},
	}
	for _, x := range casos {
		if got := x.p.Resolver(); got != x.quer {
			t.Errorf("%s: %q, esperado %q", x.nome, got, x.quer)
		}
	}
}

func TestMarcarResposta(t *testing.T) {
	h := http.Header{}
	textos.MarcarResposta(h, textos.Es419)
	if h.Get("Content-Language") != "es-419" || h.Get("Vary") != "Accept-Language" {
		t.Errorf("cabeçalhos: %v", h)
	}
}

func TestChaveComoParametroETraduzidaNoIdioma(t *testing.T) {
	c := catalogoDe(t, `{"s":"A pessoa já está {{estado}}.","e":{"ativa":"ativa"}}`, `{"s":"La persona ya está {{estado}}.","e":{"ativa":"activa"}}`)
	p := textos.Parametros{"estado": textos.Chave("e.ativa")}
	if got := c.T(textos.PtBR, "s", p); got != "A pessoa já está ativa." {
		t.Errorf("pt-BR: %q", got)
	}
	if got := c.T(textos.Es419, "s", p); got != "La persona ya está activa." {
		t.Errorf("es-419: %q", got)
	}
}
