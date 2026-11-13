package integracao_test

import (
	"net/http"
	"strings"
	"testing"
)

// Negociação do idioma da API (ADR-013): Accept-Language, depois o idioma do
// usuário, depois o padrão da empresa, depois pt-BR. O código do erro nunca
// muda; title, detail e erros[].mensagem saem no idioma; Content-Language e
// Vary acompanham a resposta com texto.

func corpoContaInvalida() map[string]any {
	return map[string]any{
		"pais": "BR", "identificacao_tipo": "cnpj", "identificacao": "11222333000180", "razao_social": "X",
		"responsavel": map[string]any{"nome": "A", "email": "a@teste.tenaz.dev"},
		"unidade":     map[string]any{"codigo": "M", "nome": "M", "turnos": turnoUtil},
	}
}

func TestIdiomaDoErroPeloAcceptLanguage(t *testing.T) {
	c := anonimo(t)
	pt := c.post("/v1/empresas", corpoContaInvalida())
	m := pt.problema(http.StatusUnprocessableEntity, "validacao_falhou")
	if m["title"] != "Os dados enviados não passaram na validação" || pt.cab.Get("Content-Language") != "pt-BR" ||
		!strings.Contains(pt.cab.Get("Vary"), "Accept-Language") {
		t.Errorf("sem Accept-Language: %v %v", m["title"], pt.cab)
	}

	es := c.post("/v1/empresas", corpoContaInvalida(), cab{"Accept-Language": "es-AR,es;q=0.9,pt;q=0.5"})
	m = es.problema(http.StatusUnprocessableEntity, "validacao_falhou")
	if m["title"] != "Los datos enviados no pasaron la validación" || es.cab.Get("Content-Language") != "es-419" ||
		!strings.Contains(es.cab.Get("Vary"), "Accept-Language") {
		t.Errorf("com Accept-Language es: %v %v", m["title"], es.cab)
	}
	campo := m["erros"].([]any)[0].(map[string]any)
	if campo["campo"] != "identificacao" || campo["codigo"] != "formato_invalido" ||
		!strings.Contains(texto(campo, "mensagem"), "dígito verificador no coincide") {
		t.Errorf("erro de campo em es-419: %v", campo)
	}

	// Idioma que a Tenaz não tem: pt-BR.
	en := c.post("/v1/empresas", corpoContaInvalida(), cab{"Accept-Language": "en-US"})
	en.problema(http.StatusUnprocessableEntity, "validacao_falhou")
	if en.cab.Get("Content-Language") != "pt-BR" {
		t.Errorf("en cai para pt-BR: %v", en.cab.Get("Content-Language"))
	}

	// Erro da validação do contrato: plural e valores aceitos no idioma.
	corpo := corpoContaInvalida()
	corpo["pais"] = "US"
	v := c.post("/v1/empresas", corpo, cab{"Accept-Language": "es"}).problema(http.StatusUnprocessableEntity, "validacao_falhou")
	if !strings.Contains(texto(v, "detail"), "campo con problema") {
		t.Errorf("detail em es-419: %v", v["detail"])
	}
}

func TestIdiomaDoUsuarioEDaEmpresaNaNegociacao(t *testing.T) {
	em := novaConta(t, "essencial", "")
	adm := entrar(t, em, "")

	emp := adm.get("/v1/empresa").ok(http.StatusOK)
	if emp["pais"] != "BR" || emp["moeda"] != "BRL" || emp["idioma_padrao"] != "pt-BR" ||
		emp["identificacao_tipo"] != "cnpj" || emp["cnpj"] != nil {
		t.Fatalf("empresa: %v", emp)
	}
	eu := adm.get("/v1/usuarios/eu").ok(http.StatusOK)
	if eu["idioma"] != "pt-BR" || eu["moeda"] != "BRL" {
		t.Fatalf("eu: %v", eu)
	}

	// País e moeda não mudam pela API.
	r := adm.get("/v1/empresa")
	for _, campo := range []string{"pais", "moeda"} {
		adm.req("PATCH", "/v1/empresa", map[string]any{campo: "AR"}, cab{"If-Match": r.etag()}).
			problema(http.StatusUnprocessableEntity, "campo_somente_leitura")
	}
	// O administrador não altera o idioma de outro usuário por PATCH /v1/usuarios/{id}.
	r = adm.get("/v1/usuarios/" + texto(eu, "id"))
	adm.req("PATCH", "/v1/usuarios/"+texto(eu, "id"), map[string]any{"idioma": "es-419"}, cab{"If-Match": r.etag()}).
		problema(http.StatusUnprocessableEntity, "campo_somente_leitura")

	// Sem Accept-Language vale o idioma do usuário; depois de trocá-lo, o erro sai em es-419.
	r = adm.get("/v1/usuarios/eu")
	adm.req("PATCH", "/v1/usuarios/eu", map[string]any{"idioma": "es-419"}, cab{"If-Match": r.etag()}).ok(http.StatusOK)
	nao := adm.get("/v1/pessoas/00000000-0000-7000-8000-000000000000")
	nao.problema(http.StatusNotFound, "recurso_nao_encontrado")
	if nao.cab.Get("Content-Language") != "es-419" || !strings.Contains(string(nao.corpo), "Recurso no encontrado") {
		t.Errorf("idioma do usuário: %v %s", nao.cab.Get("Content-Language"), nao.corpo)
	}
	// O cabeçalho vale mais que o idioma do usuário.
	pt := adm.get("/v1/pessoas/00000000-0000-7000-8000-000000000000", cab{"Accept-Language": "pt-BR"})
	pt.problema(http.StatusNotFound, "recurso_nao_encontrado")
	if pt.cab.Get("Content-Language") != "pt-BR" || !strings.Contains(string(pt.corpo), "Recurso não encontrado") {
		t.Errorf("Accept-Language sobre o usuário: %v %s", pt.cab.Get("Content-Language"), pt.corpo)
	}
	// Idioma inválido no perfil.
	r = adm.get("/v1/usuarios/eu")
	adm.req("PATCH", "/v1/usuarios/eu", map[string]any{"idioma": "fr"}, cab{"If-Match": r.etag()}).
		problema(http.StatusUnprocessableEntity, "validacao_falhou")
}

// O nome da tela em GET /v1/capturas-auditoria vem no idioma da resposta; o
// código da tela não muda.
func TestNomeDaTelaNoIdioma(t *testing.T) {
	adm := entrar(t, novaConta(t, "essencial", ""), "")
	nome := func(acceptLanguage string) (string, string) {
		extras := []cab{}
		if acceptLanguage != "" {
			extras = append(extras, cab{"Accept-Language": acceptLanguage})
		}
		r := adm.get("/v1/capturas-auditoria", extras...)
		m := r.ok(http.StatusOK)
		for _, x := range m["dados"].([]any) {
			c := x.(map[string]any)
			if c["tela"] == "adm.pessoas" {
				return texto(c, "nome"), r.cab.Get("Content-Language")
			}
		}
		t.Fatal("tela adm.pessoas ausente")
		return "", ""
	}
	if n, l := nome(""); n != "Pessoas" || l != "pt-BR" {
		t.Errorf("pt-BR: %q %q", n, l)
	}
	if n, l := nome("es"); n != "Personas" || l != "es-419" {
		t.Errorf("es-419: %q %q", n, l)
	}
}

// Os textos que saíam em português fixo, agora pelo catálogo: o mesmo código de
// erro, o detail no idioma pedido.
func TestTextosMigradosNoIdioma(t *testing.T) {
	es := cab{"Accept-Language": "es"}
	anon := anonimo(t)

	// Sem sessão.
	r := anon.get("/v1/empresa", es)
	m := r.problema(http.StatusUnauthorized, "nao_autenticado")
	if m["detail"] != "Inicie sesión en Tenaz para continuar." {
		t.Errorf("sem sessão: %v", m["detail"])
	}
	pt := anon.get("/v1/empresa").problema(http.StatusUnauthorized, "nao_autenticado")
	if pt["detail"] != "Entre na Tenaz para continuar." {
		t.Errorf("sem sessão em pt-BR: %v", pt["detail"])
	}

	// Login recusado.
	l := anon.post("/v1/auth/login", map[string]any{"email": "ninguem@teste.tenaz.dev", "senha": "x"}, es).
		problema(http.StatusUnauthorized, "credencial_invalida")
	if !strings.Contains(texto(l, "detail"), "contraseña incorrectos") {
		t.Errorf("login: %v", l["detail"])
	}

	adm := entrar(t, novaConta(t, "essencial", ""), "")
	// If-Match ausente.
	p := adm.get("/v1/pessoas?limite=1").ok(http.StatusOK)
	_ = p
	f := adm.post("/v1/pessoas", map[string]any{"nome": "Ana"}).ok(http.StatusCreated)
	sem := adm.req("PATCH", "/v1/pessoas/"+texto(f, "id"), map[string]any{"cargo": "x"}, es)
	sem.problema(http.StatusPreconditionRequired, "pre_condicao_obrigatoria")
	if !strings.Contains(string(sem.corpo), "encabezado If-Match") {
		t.Errorf("If-Match: %s", sem.corpo)
	}
	// Cursor inválido.
	c := adm.get("/v1/pessoas?cursor=lixo", es)
	c.problema(http.StatusBadRequest, "cursor_invalido")
	if !strings.Contains(string(c.corpo), "Cursor dañado") {
		t.Errorf("cursor: %s", c.corpo)
	}
	// Estado atual, com o estado traduzido: inativar a pessoa duas vezes.
	r1 := adm.get("/v1/pessoas/" + texto(f, "id"))
	adm.req("DELETE", "/v1/pessoas/"+texto(f, "id"), nil, cab{"If-Match": r1.etag()}).ok(http.StatusNoContent)
	r2 := adm.get("/v1/pessoas/" + texto(f, "id"))
	d := adm.req("DELETE", "/v1/pessoas/"+texto(f, "id"), nil, cab{"If-Match": r2.etag(), "Accept-Language": "es"})
	d.problema(http.StatusConflict, "transicao_invalida")
	if !strings.Contains(string(d.corpo), "La persona ya está inactiva") {
		t.Errorf("transição: %s", d.corpo)
	}
	// Fuso de outro país: valor_invalido com valores_aceitos, nos dois idiomas.
	corpo := corpoContaInvalida()
	corpo["identificacao"] = cnpjAleatorio()
	corpo["fuso_horario"] = "America/Asuncion"
	fz := anon.post("/v1/empresas", corpo, es).problema(http.StatusUnprocessableEntity, "validacao_falhou")
	campo := fz["erros"].([]any)[0].(map[string]any)
	if campo["campo"] != "fuso_horario" || campo["codigo"] != "valor_invalido" || len(campo["valores_aceitos"].([]any)) != 16 ||
		!strings.Contains(texto(campo, "mensagem"), "zona horaria") {
		t.Errorf("fuso: %v", campo)
	}
}

// custo_hora na API: texto decimal sem sinal, até seis casas. O padrão do
// contrato e a leitura do serviço dão a mesma resposta: validacao_falhou com
// formato_invalido no campo custo_hora.
func TestCustoHoraNaApi(t *testing.T) {
	adm := entrar(t, novaConta(t, "essencial", ""), "")
	novo := func(custo any) *resposta {
		return adm.post("/v1/usuarios", map[string]any{"nome": "Custo", "email": emailAleatorio("custo"),
			"perfis": []string{"tecnico"}, "custo_hora": custo})
	}
	for _, ruim := range []string{"-1", "0.1234567", "1,5", "1e3", "+1", "01", ""} {
		m := novo(ruim).problema(http.StatusUnprocessableEntity, "validacao_falhou")
		campo := m["erros"].([]any)[0].(map[string]any)
		if campo["campo"] != "custo_hora" || campo["codigo"] != "formato_invalido" {
			t.Errorf("%q: %v", ruim, campo)
		}
	}
	u := novo("85.5").ok(http.StatusCreated)
	if u["custo_hora"] != "85.50" || u["moeda"] != "BRL" {
		t.Errorf("custo_hora: %v %v", u["custo_hora"], u["moeda"])
	}
	if u = novo("6.299").ok(http.StatusCreated); u["custo_hora"] != "6.299" {
		t.Errorf("três casas: %v", u["custo_hora"])
	}
	r := adm.get("/v1/usuarios/" + texto(u, "id"))
	adm.req("PATCH", "/v1/usuarios/"+texto(u, "id"), map[string]any{"custo_hora": "-3"}, cab{"If-Match": r.etag()}).
		problema(http.StatusUnprocessableEntity, "validacao_falhou")
	p := adm.req("PATCH", "/v1/usuarios/"+texto(u, "id"), map[string]any{"custo_hora": "45.123456"}, cab{"If-Match": r.etag()}).ok(http.StatusOK)
	if p["custo_hora"] != "45.123456" {
		t.Errorf("PATCH com seis casas: %v", p["custo_hora"])
	}
	r = adm.get("/v1/usuarios/" + texto(u, "id"))
	if p = adm.req("PATCH", "/v1/usuarios/"+texto(u, "id"), map[string]any{"custo_hora": nil}, cab{"If-Match": r.etag()}).ok(http.StatusOK); p["custo_hora"] != nil {
		t.Errorf("PATCH com null: %v", p["custo_hora"])
	}
}
