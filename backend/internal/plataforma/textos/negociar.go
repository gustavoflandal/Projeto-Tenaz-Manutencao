package textos

import (
	"context"
	"net/http"
	"strings"

	"golang.org/x/text/language"
)

// Negociação do idioma da resposta (ADR-013, Parte B, item 3), nesta ordem:
//
//  1. Accept-Language (RFC 9110, seção 12.5.4, com pesos q): vale a primeira
//     etiqueta aceitável que casa com um idioma suportado, com pt e pt-*
//     levando a pt-BR, e es e es-* levando a es-419;
//  2. sem casamento, o idioma do usuário da sessão;
//  3. token de API ou rota sem usuário, o idioma padrão da empresa;
//  4. pt-BR.
//
// Toda resposta que leva texto para pessoa traz Content-Language e
// Vary: Accept-Language.

// CabecalhoVary é o valor de Vary das respostas com texto para pessoa.
const CabecalhoVary = "Accept-Language"

// DeAcceptLanguage devolve o idioma suportado que o cabeçalho pede, ou falso
// se ele é vazio, malformado ou só pede idiomas fora do catálogo. Pesos q
// ordenam as etiquetas; q=0 recusa a etiqueta.
func DeAcceptLanguage(cabecalho string) (Idioma, bool) {
	if strings.TrimSpace(cabecalho) == "" {
		return "", false
	}
	etiquetas, _, err := language.ParseAcceptLanguage(cabecalho)
	if err != nil {
		return "", false
	}
	for _, e := range etiquetas {
		base, _ := e.Base()
		switch base.String() {
		case "pt":
			return PtBR, true
		case "es":
			return Es419, true
		}
	}
	return "", false
}

// Preferencias reúne o que a requisição sabe para negociar o idioma. O que
// ainda não se sabe (usuário antes da credencial, empresa de rota pública)
// fica vazio.
type Preferencias struct {
	AcceptLanguage string // cabeçalho, como veio
	Usuario        string // usuario.idioma da sessão
	Empresa        string // empresa.idioma_padrao
}

// Resolver aplica a ordem da negociação.
func (p Preferencias) Resolver() Idioma {
	if i, ok := DeAcceptLanguage(p.AcceptLanguage); ok {
		return i
	}
	if Suportado(p.Usuario) {
		return Idioma(p.Usuario)
	}
	if Suportado(p.Empresa) {
		return Idioma(p.Empresa)
	}
	return PtBR
}

type chavePreferencias struct{}

// ComPreferencias guarda as preferências no contexto da requisição.
func ComPreferencias(ctx context.Context, p Preferencias) context.Context {
	return context.WithValue(ctx, chavePreferencias{}, p)
}

// PreferenciasDe devolve as preferências do contexto (vazias se não houver).
func PreferenciasDe(ctx context.Context) Preferencias {
	p, _ := ctx.Value(chavePreferencias{}).(Preferencias)
	return p
}

// ComUsuario acrescenta o idioma do usuário e o da empresa às preferências do
// contexto, depois que a credencial foi resolvida.
func ComUsuario(ctx context.Context, usuario, empresa string) context.Context {
	p := PreferenciasDe(ctx)
	p.Usuario, p.Empresa = usuario, empresa
	return ComPreferencias(ctx, p)
}

// DoContexto resolve o idioma da requisição a partir do contexto.
func DoContexto(ctx context.Context) Idioma { return PreferenciasDe(ctx).Resolver() }

// MarcarResposta escreve Content-Language e Vary na resposta que leva texto
// para pessoa.
func MarcarResposta(h http.Header, idioma Idioma) {
	h.Set("Content-Language", string(idioma))
	h.Add("Vary", CabecalhoVary)
}
