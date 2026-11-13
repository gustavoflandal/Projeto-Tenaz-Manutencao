package erros_test

import (
	"testing"

	"tenaz/internal/plataforma/erros"
	"tenaz/internal/plataforma/textos"
)

// (d) da conferência dos catálogos (ADR-013): todo código de erro registrado
// tem título nos dois idiomas, e o título do catálogo pt-BR é o escrito no
// código.
func TestTodoCodigoTemTituloNosDoisIdiomas(t *testing.T) {
	c := textos.Padrao()
	for _, e := range erros.Todos() {
		chave := erros.ChaveTitulo(e.Codigo)
		for _, id := range textos.Idiomas() {
			if !c.Existe(id, chave) {
				t.Errorf("código %s sem título em %s (chave %s)", e.Codigo, id, chave)
			}
		}
		if pt, _ := c.Bruto(textos.PtBR, chave); pt != e.Titulo {
			t.Errorf("código %s: título do catálogo %q difere do escrito no código %q", e.Codigo, pt, e.Titulo)
		}
	}
}

func TestProblemaNoIdiomaNegociado(t *testing.T) {
	err := erros.EntradaInvalida.ComChave("detalhe.campos_problema_other", textos.Parametros{"n": 2}).
		ComCampoChave("cnpj", erros.CampoFormatoInvalido, "identificacao.digito", textos.Parametros{"tipo": "CNPJ"}).
		ComDado("estado_atual", "x")

	pt := erros.ParaProblemaIdioma(err, "/v1/x", textos.PtBR)
	es := erros.ParaProblemaIdioma(err, "/v1/x", textos.Es419)

	if pt.Title != "Os dados enviados não passaram na validação" || es.Title != "Los datos enviados no pasaron la validación" {
		t.Errorf("títulos: %q / %q", pt.Title, es.Title)
	}
	if pt.Detail != "2 campos com problema." || es.Detail != "2 campos con problema." {
		t.Errorf("detalhes: %q / %q", pt.Detail, es.Detail)
	}
	if pt.Erros[0].Mensagem != "CNPJ inválido: o dígito verificador não confere." ||
		es.Erros[0].Mensagem != "CNPJ no válido: el dígito verificador no coincide." {
		t.Errorf("mensagens: %q / %q", pt.Erros[0].Mensagem, es.Erros[0].Mensagem)
	}
	// O código, o tipo, o status e os membros adicionais não mudam com o idioma.
	if pt.Codigo != es.Codigo || pt.Type != es.Type || pt.Status != es.Status ||
		pt.Erros[0].Codigo != es.Erros[0].Codigo || *pt.EstadoAtual != *es.EstadoAtual {
		t.Errorf("campo estável mudou com o idioma: %+v / %+v", pt, es)
	}
}

// Texto livre, sem chave, sai como foi escrito (dívida de tradução registrada).
func TestDetalheSemChaveSaiComoEscrito(t *testing.T) {
	p := erros.ParaProblemaIdioma(erros.NaoEncontrado.Com("detalhe livre"), "", textos.Es419)
	if p.Detail != "detalhe livre" || p.Title != "Recurso no encontrado" {
		t.Errorf("%+v", p)
	}
}

// Logs registram código e chave, nunca o texto traduzido.
func TestErrorUsaAChave(t *testing.T) {
	e := erros.EntradaInvalida.ComChave("detalhe.corpo_json", nil)
	if got := e.Error(); got != "validacao_falhou: detalhe.corpo_json" {
		t.Errorf("Error() = %q", got)
	}
}

// Texto técnico é só para o log: aparece em Error() e nunca na resposta.
func TestTecnicoNaoVaiParaAResposta(t *testing.T) {
	e := erros.Interno.Tecnico("falha %d", 7)
	if e.Error() != "erro_interno: falha 7" {
		t.Errorf("Error() = %q", e.Error())
	}
	n := erros.NaoEncontrado.Tecnico("usuário não encontrado")
	p := erros.ParaProblemaIdioma(n, "/x", textos.Es419)
	if p.Detail != "" || p.Title != "Recurso no encontrado" {
		t.Errorf("a resposta trouxe o texto técnico: %+v", p)
	}
}
