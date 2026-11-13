package textos

import (
	"fmt"
	"regexp"
	"slices"
	"sort"
	"strings"
)

var (
	reVariavel = regexp.MustCompile(`\{\{\s*([A-Za-z0-9_]+)\s*\}\}`)
	rePlural   = regexp.MustCompile(`^(.*)_(zero|one|two|few|many|other)$`)
)

// termosProibidos são os termos que o glossário bilíngue (pacotes/i18n/src/
// glossario.ts) manda não usar em es-419, como o teste de tom de voz da
// interface já confere. Palavra inteira, sem diferenciar maiúscula.
var termosProibidos = map[Idioma][]string{
	Es419: {"e-mail", "email", "mail", "password", "clave", "link", "offline", "ticket", "checklist", "stock", "ingresá", "ingresa", "vos"},
}

// formasObrigatorias são as formas de plural que cada idioma exige da chave
// plural: as dos dois idiomas do MVP têm "one" e "other".
var formasObrigatorias = []string{"one", "other"}

// Conferir compara os catálogos e devolve os problemas, um por linha, ou nada
// se estiver tudo certo. É a tarefa `textos` do CI (ferramentas/tarefas.sh):
//
//	(a) chave faltando ou sobrando em qualquer idioma, tomando o pt-BR como
//	    fonte da verdade;
//	(b) variáveis {{...}} que diferem entre os idiomas;
//	(c) forma de plural exigida pelo idioma que falta;
//	(e) termo proibido pelo glossário em es-419.
//
// O item (d), um código de erro sem título nos dois idiomas, é conferido no
// pacote erros, que conhece os códigos.
func Conferir(c *Catalogo) []string {
	var problemas []string
	origem := c.textos[PtBR]
	for _, id := range Idiomas() {
		alvo := c.textos[id]
		if id != PtBR {
			for _, k := range c.Chaves(PtBR) {
				if _, ok := alvo[k]; !ok {
					problemas = append(problemas, fmt.Sprintf("%s: falta a chave %q", id, k))
				}
			}
			for _, k := range c.Chaves(id) {
				if _, ok := origem[k]; !ok {
					problemas = append(problemas, fmt.Sprintf("%s: a chave %q sobra (não existe em pt-BR)", id, k))
				}
			}
			for _, k := range c.Chaves(PtBR) {
				a, ok := alvo[k]
				if !ok {
					continue
				}
				if !slices.Equal(variaveis(origem[k]), variaveis(a)) {
					problemas = append(problemas, fmt.Sprintf("%s: a chave %q usa as variáveis %v, pt-BR usa %v",
						id, k, variaveis(a), variaveis(origem[k])))
				}
			}
		}
		problemas = append(problemas, conferirPlurais(id, alvo)...)
		for _, k := range c.Chaves(id) {
			if t := proibido(id, alvo[k]); t != "" {
				problemas = append(problemas, fmt.Sprintf("%s: a chave %q usa o termo proibido %q", id, k, t))
			}
		}
	}
	return problemas
}

func variaveis(texto string) []string {
	var l []string
	for _, m := range reVariavel.FindAllStringSubmatch(texto, -1) {
		if !slices.Contains(l, m[1]) {
			l = append(l, m[1])
		}
	}
	sort.Strings(l)
	return l
}

func conferirPlurais(id Idioma, textos map[string]string) []string {
	grupos := map[string][]string{}
	for k := range textos {
		if m := rePlural.FindStringSubmatch(k); m != nil {
			grupos[m[1]] = append(grupos[m[1]], m[2])
		}
	}
	var problemas []string
	bases := make([]string, 0, len(grupos))
	for b := range grupos {
		bases = append(bases, b)
	}
	sort.Strings(bases)
	for _, b := range bases {
		for _, f := range formasObrigatorias {
			if !slices.Contains(grupos[b], f) {
				problemas = append(problemas, fmt.Sprintf("%s: o plural %q não tem a forma %s", id, b, f))
			}
		}
	}
	return problemas
}

func proibido(id Idioma, texto string) string {
	sem := reVariavel.ReplaceAllString(strings.ToLower(texto), " ")
	for _, t := range termosProibidos[id] {
		re := regexp.MustCompile(`(^|[^\p{L}-])` + regexp.QuoteMeta(t) + `($|[^\p{L}-])`)
		if re.MatchString(sem) {
			return t
		}
	}
	return ""
}
