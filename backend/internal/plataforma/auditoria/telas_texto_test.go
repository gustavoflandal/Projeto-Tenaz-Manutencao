package auditoria_test

import (
	"testing"

	"tenaz/internal/plataforma/auditoria"
	"tenaz/internal/plataforma/textos"
)

// Toda tela do catálogo fechado tem nome nos dois idiomas, e o do pt-BR é o
// escrito no código (ADR-013): a tela nova que esquecer o texto quebra aqui.
func TestTodaTelaTemNomeNosDoisIdiomas(t *testing.T) {
	c := textos.Padrao()
	for _, tela := range auditoria.Telas() {
		chave := "tela." + tela.Codigo
		for _, id := range textos.Idiomas() {
			if !c.Existe(id, chave) {
				t.Errorf("tela %s sem nome em %s (chave %s)", tela.Codigo, id, chave)
			}
		}
		if pt, _ := c.Bruto(textos.PtBR, chave); pt != tela.Nome {
			t.Errorf("tela %s: nome do catálogo %q difere do escrito no código %q", tela.Codigo, pt, tela.Nome)
		}
	}
}
