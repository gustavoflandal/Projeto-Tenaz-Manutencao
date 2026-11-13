// Package textos é o catálogo de textos do servidor (ADR-013, aceito em
// 13/11/2026): título e detalhe dos erros da API, assunto e corpo dos e-mails
// e, mais adiante, notificações e documentos gerados, em pt-BR e es-419.
//
// Convenção única com o catálogo da interface (pacotes/i18n): chaves estáveis
// em JSON aninhado, interpolação {{nome}} e plural por sufixo CLDR (_one,
// _other). Um arquivo por idioma em catalogos/, embutido no binário. A fonte da
// verdade de cada texto é o arquivo pt-BR; o es-419 é tradução dele.
//
// Sem biblioteca de i18n: busca, interpolação e plural cabem aqui, sobre
// golang.org/x/text/feature/plural. O código do erro nunca muda com o idioma;
// só o texto para pessoas.
package textos

import (
	"embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"sort"
	"strings"
	"sync"

	"golang.org/x/text/feature/plural"
	"golang.org/x/text/language"
)

// Idioma é uma etiqueta BCP 47 do catálogo fechado de idiomas suportados.
type Idioma string

// Idiomas suportados. PtBR é o idioma de origem e o de último recurso.
const (
	PtBR  Idioma = "pt-BR"
	Es419 Idioma = "es-419"
)

// Idiomas lista os idiomas suportados, o de origem primeiro.
func Idiomas() []Idioma { return []Idioma{PtBR, Es419} }

// Suportado diz se a etiqueta é um idioma do catálogo.
func Suportado(s string) bool { return s == string(PtBR) || s == string(Es419) }

// Parametros são os valores das variáveis {{nome}} de um texto.
type Parametros map[string]any

//go:embed catalogos/*.json
var arquivos embed.FS

// Catalogo guarda os textos de todos os idiomas, achatados em "a.b.c".
// Depois de montado não muda: pode ser usado por várias goroutines.
type Catalogo struct {
	textos map[Idioma]map[string]string
}

// Novo carrega os catálogos de um sistema de arquivos com um <idioma>.json por
// idioma suportado.
func Novo(sistema fs.FS) (*Catalogo, error) {
	c := &Catalogo{textos: map[Idioma]map[string]string{}}
	for _, id := range Idiomas() {
		b, err := fs.ReadFile(sistema, string(id)+".json")
		if err != nil {
			return nil, fmt.Errorf("catálogo %s: %w", id, err)
		}
		var arvore map[string]any
		if err := json.Unmarshal(b, &arvore); err != nil {
			return nil, fmt.Errorf("catálogo %s inválido: %w", id, err)
		}
		plano := map[string]string{}
		if err := achatar("", arvore, plano); err != nil {
			return nil, fmt.Errorf("catálogo %s: %w", id, err)
		}
		c.textos[id] = plano
	}
	return c, nil
}

func achatar(prefixo string, no map[string]any, saida map[string]string) error {
	for k, v := range no {
		chave := k
		if prefixo != "" {
			chave = prefixo + "." + k
		}
		switch x := v.(type) {
		case string:
			saida[chave] = x
		case map[string]any:
			if err := achatar(chave, x, saida); err != nil {
				return err
			}
		default:
			return fmt.Errorf("a chave %s tem valor que não é texto nem grupo", chave)
		}
	}
	return nil
}

// padrao é o catálogo embutido, montado uma vez e nunca alterado. Um catálogo
// embutido que não carrega é erro de programação: a inicialização para.
var padrao = sync.OnceValue(func() *Catalogo {
	sub, err := fs.Sub(arquivos, "catalogos")
	if err != nil {
		panic(fmt.Sprintf("textos: %v", err))
	}
	c, err := Novo(sub)
	if err != nil {
		panic(fmt.Sprintf("textos: %v", err))
	}
	return c
})

// Padrao devolve o catálogo embutido no binário.
func Padrao() *Catalogo { return padrao() }

// Existe diz se a chave tem texto no idioma, sem cair para o idioma de origem.
func (c *Catalogo) Existe(idioma Idioma, chave string) bool {
	_, ok := c.textos[idioma][chave]
	return ok
}

// Chaves devolve as chaves do idioma, em ordem.
func (c *Catalogo) Chaves(idioma Idioma) []string {
	l := make([]string, 0, len(c.textos[idioma]))
	for k := range c.textos[idioma] {
		l = append(l, k)
	}
	sort.Strings(l)
	return l
}

// Bruto devolve o texto da chave sem interpolar, e se existe.
func (c *Catalogo) Bruto(idioma Idioma, chave string) (string, bool) {
	t, ok := c.textos[idioma][chave]
	return t, ok
}

// T devolve o texto da chave no idioma, com as variáveis trocadas. Chave sem
// texto no idioma cai para o pt-BR; sem texto no pt-BR, devolve a própria
// chave (o teste de conferência impede que isso chegue a produção).
func (c *Catalogo) T(idioma Idioma, chave string, p Parametros) string {
	t, ok := c.textos[idioma][chave]
	if !ok {
		if t, ok = c.textos[PtBR][chave]; !ok {
			return chave
		}
	}
	return c.interpolar(idioma, t, p)
}

// Plural escolhe a forma pela regra CLDR do idioma (em português, zero usa
// _one; em espanhol, _other), busca chave_<forma> e cai para chave_other. A
// contagem também fica disponível como {{n}}.
func (c *Catalogo) Plural(idioma Idioma, chave string, n int, p Parametros) string {
	forma := formaPlural(idioma, n)
	params := Parametros{"n": n}
	for k, v := range p {
		params[k] = v
	}
	for _, f := range []string{forma, "other"} {
		if t, ok := c.textos[idioma][chave+"_"+f]; ok {
			return c.interpolar(idioma, t, params)
		}
	}
	for _, f := range []string{forma, "other"} {
		if t, ok := c.textos[PtBR][chave+"_"+f]; ok {
			return c.interpolar(idioma, t, params)
		}
	}
	return chave
}

func formaPlural(idioma Idioma, n int) string {
	if n < 0 {
		n = -n
	}
	etiqueta := language.MustParse(string(idioma))
	switch plural.Cardinal.MatchPlural(etiqueta, n, 0, 0, 0, 0) {
	case plural.Zero:
		return "zero"
	case plural.One:
		return "one"
	case plural.Two:
		return "two"
	case plural.Few:
		return "few"
	case plural.Many:
		return "many"
	}
	return "other"
}

// Chave, como valor de parâmetro, é traduzida na hora de renderizar: serve a
// valores que são, eles mesmos, texto do catálogo (o estado "ativa" vira
// "activa" em es-419).
type Chave string

func (c *Catalogo) interpolar(idioma Idioma, texto string, p Parametros) string {
	if len(p) == 0 || !strings.Contains(texto, "{{") {
		return texto
	}
	var b strings.Builder
	for {
		i := strings.Index(texto, "{{")
		if i < 0 {
			break
		}
		j := strings.Index(texto[i:], "}}")
		if j < 0 {
			break
		}
		nome := strings.TrimSpace(texto[i+2 : i+j])
		b.WriteString(texto[:i])
		if v, ok := p[nome]; ok {
			if k, ehChave := v.(Chave); ehChave {
				b.WriteString(c.T(idioma, string(k), nil))
			} else {
				fmt.Fprint(&b, v)
			}
		} else {
			b.WriteString(texto[i : i+j+2])
		}
		texto = texto[i+j+2:]
	}
	b.WriteString(texto)
	return b.String()
}

// T é Padrao().T.
func T(idioma Idioma, chave string, p Parametros) string { return Padrao().T(idioma, chave, p) }

// Plural é Padrao().Plural.
func Plural(idioma Idioma, chave string, n int, p Parametros) string {
	return Padrao().Plural(idioma, chave, n, p)
}
