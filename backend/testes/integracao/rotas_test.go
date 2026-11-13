package integracao_test

import (
	"net/http"
	"strings"
	"testing"

	"github.com/google/uuid"
)

// cenarioADM tem duas empresas: A, com administrador, visualizador e
// prestador, e B, com administrador.
type cenarioADM struct {
	admA, visA, presA, admB *cliente
	unidade, turno, feriado string
	pessoa, usuario, evento string
}

func montarCenario(t *testing.T) cenarioADM {
	t.Helper()
	var c cenarioADM
	c.admA = entrar(t, novaConta(t, "profissional", ""), "")
	c.admB = entrar(t, novaConta(t, "profissional", ""), "")

	var idVis string
	c.visA, idVis, _ = convidar(t, c.admA, []string{"visualizador"}, nil)
	c.usuario = idVis
	c.presA, _, _ = convidar(t, c.admA, []string{"prestador"}, map[string]any{"fornecedor_id": uuid.NewString()})

	u := c.admA.post("/v1/unidades", map[string]any{"codigo": "FIL", "nome": "Filial Norte", "turnos": turnoUtil}).ok(http.StatusCreated)
	c.unidade = texto(u, "id")
	c.turno = texto(u["turnos"].([]any)[0].(map[string]any), "id")
	f := c.admA.post("/v1/unidades/"+c.unidade+"/feriados", map[string]any{"data": "2027-04-21", "descricao": "Tiradentes"}).ok(http.StatusCreated)
	c.feriado = texto(f, "id")
	p := c.admA.post("/v1/pessoas", map[string]any{"nome": "Carlos Motorista"}).ok(http.StatusCreated)
	c.pessoa = texto(p, "id")
	ev := c.admA.get("/v1/eventos-auditoria?limite=1").ok(http.StatusOK)
	c.evento = texto(ev["dados"].([]any)[0].(map[string]any), "id")
	return c
}

var turnoUtil = []any{map[string]any{"dia_semana": 1, "hora_inicio": "08:00", "hora_fim": "17:00"}}

// rotaTeste descreve uma rota para a matriz de acesso.
type rotaTeste struct {
	metodo, caminho string
	corpo           any
	ifMatch         bool // exige If-Match
	visNega         bool // o visualizador não tem a permissão: 403
	comID           bool // tem identificador de registro da empresa A
	proprio         bool // rota do próprio usuário: externo também usa
}

func (c cenarioADM) rotas() []rotaTeste {
	un, tu, fe, pe, us, evt := "/v1/unidades/"+c.unidade, "/turnos/"+c.turno, "/feriados/"+c.feriado, "/v1/pessoas/"+c.pessoa, "/v1/usuarios/"+c.usuario, "/v1/eventos-auditoria/"+c.evento
	return []rotaTeste{
		{metodo: "GET", caminho: "/v1/empresa"},
		{metodo: "PATCH", caminho: "/v1/empresa", corpo: map[string]any{"nome_fantasia": "X"}, ifMatch: true, visNega: true},
		{metodo: "GET", caminho: "/v1/unidades"},
		{metodo: "POST", caminho: "/v1/unidades", corpo: map[string]any{"codigo": "Z", "nome": "Z", "turnos": turnoUtil}, visNega: true},
		{metodo: "GET", caminho: un, comID: true},
		{metodo: "PATCH", caminho: un, corpo: map[string]any{"nome": "N"}, ifMatch: true, visNega: true, comID: true},
		{metodo: "DELETE", caminho: un, ifMatch: true, visNega: true, comID: true},
		{metodo: "POST", caminho: un + "/reativar", ifMatch: true, visNega: true, comID: true},
		{metodo: "POST", caminho: un + "/turnos", corpo: map[string]any{"dia_semana": 6, "hora_inicio": "08:00", "hora_fim": "12:00"}, visNega: true, comID: true},
		{metodo: "PATCH", caminho: un + tu, corpo: map[string]any{"hora_fim": "18:00"}, ifMatch: true, visNega: true, comID: true},
		{metodo: "DELETE", caminho: un + tu, ifMatch: true, visNega: true, comID: true},
		{metodo: "GET", caminho: un + "/feriados", comID: true},
		{metodo: "POST", caminho: un + "/feriados", corpo: map[string]any{"data": "2027-09-07", "descricao": "Independência"}, visNega: true, comID: true},
		{metodo: "DELETE", caminho: un + fe, ifMatch: true, visNega: true, comID: true},
		{metodo: "GET", caminho: "/v1/pessoas", visNega: true},
		{metodo: "POST", caminho: "/v1/pessoas", corpo: map[string]any{"nome": "P"}, visNega: true},
		{metodo: "GET", caminho: pe, visNega: true, comID: true},
		{metodo: "PATCH", caminho: pe, corpo: map[string]any{"cargo": "Mecânico"}, ifMatch: true, visNega: true, comID: true},
		{metodo: "DELETE", caminho: pe, ifMatch: true, visNega: true, comID: true},
		{metodo: "POST", caminho: pe + "/reativar", ifMatch: true, visNega: true, comID: true},
		{metodo: "GET", caminho: "/v1/usuarios", visNega: true},
		{metodo: "POST", caminho: "/v1/usuarios", corpo: map[string]any{"nome": "U", "email": emailAleatorio("u"), "perfis": []string{"tecnico"}}, visNega: true},
		{metodo: "GET", caminho: us, visNega: true, comID: true},
		{metodo: "PATCH", caminho: us, corpo: map[string]any{"telefone": "11999990000"}, ifMatch: true, visNega: true, comID: true},
		{metodo: "POST", caminho: us + "/desativar", ifMatch: true, visNega: true, comID: true},
		{metodo: "POST", caminho: us + "/reativar", ifMatch: true, visNega: true, comID: true},
		{metodo: "POST", caminho: us + "/reenviar-convite", visNega: true, comID: true},
		{metodo: "POST", caminho: us + "/encerrar-sessoes", visNega: true, comID: true},
		{metodo: "POST", caminho: us + "/permissoes", corpo: map[string]any{"permissao": "os.encerrar"}, visNega: true, comID: true},
		{metodo: "DELETE", caminho: us + "/permissoes/os.encerrar", visNega: true, comID: true},
		{metodo: "GET", caminho: "/v1/perfis", visNega: true},
		{metodo: "GET", caminho: "/v1/permissoes", visNega: true},
		{metodo: "GET", caminho: "/v1/eventos-auditoria", visNega: true},
		{metodo: "GET", caminho: evt, visNega: true, comID: true},
		{metodo: "GET", caminho: "/v1/usuarios/eu", proprio: true},
		{metodo: "PATCH", caminho: "/v1/usuarios/eu", corpo: map[string]any{"nome": "Eu"}, ifMatch: true, proprio: true},
		{metodo: "GET", caminho: "/v1/auth/sessao", proprio: true},
		{metodo: "POST", caminho: "/v1/auth/trocar-empresa", corpo: map[string]any{"empresa_escolhida_id": uuid.NewString()}, proprio: true},
		{metodo: "POST", caminho: "/v1/auth/mfa/desativar", corpo: map[string]any{"codigo": "123456"}, proprio: true},
	}
}

// Matriz de acesso de toda rota autenticada de E01: 401 sem credencial,
// 404 para perfil externo fora das rotas permitidas e para id de outra
// empresa, 403 sem permissão, 403 sem CSRF, 428 sem If-Match e 412 com
// If-Match errado, sempre em problem+json com o código estável.
func TestMatrizDeAcesso(t *testing.T) {
	c := montarCenario(t)
	for _, r := range c.rotas() {
		t.Run(r.metodo+" "+r.caminho, func(t *testing.T) {
			ifm := cab{"If-Match": `"1"`}
			anonimo(t).req(r.metodo, r.caminho, r.corpo, ifm).problema(http.StatusUnauthorized, "nao_autenticado")

			comT := func(cl *cliente) *cliente { x := *cl; x.t = t; return &x }
			if !r.proprio {
				comT(c.presA).req(r.metodo, r.caminho, r.corpo, ifm).problema(http.StatusNotFound, "recurso_nao_encontrado")
			}
			if r.visNega {
				comT(c.visA).req(r.metodo, r.caminho, r.corpo, ifm).problema(http.StatusForbidden, "permissao_negada")
			}
			if r.comID {
				comT(c.admB).req(r.metodo, r.caminho, r.corpo, ifm).problema(http.StatusNotFound, "recurso_nao_encontrado")
			}
			if r.metodo != http.MethodGet {
				semCSRF := comT(c.admA)
				semCSRF.csrf = ""
				semCSRF.req(r.metodo, r.caminho, r.corpo, ifm).problema(http.StatusForbidden, "csrf_invalido")
				// Cookie válido com o CSRF de outra sessão válida, ou com o
				// próprio cookie no lugar do token: também csrf_invalido.
				for _, errado := range []string{c.admB.csrf, c.visA.csrf, c.admA.cookie} {
					outro := comT(c.admA)
					outro.csrf = errado
					outro.req(r.metodo, r.caminho, r.corpo, ifm).problema(http.StatusForbidden, "csrf_invalido")
				}
				outraOrigem := comT(c.admA)
				outraOrigem.origem = "https://atacante.exemplo"
				outraOrigem.req(r.metodo, r.caminho, r.corpo, ifm).problema(http.StatusForbidden, "csrf_invalido")
			}
			if r.ifMatch {
				comT(c.admA).req(r.metodo, r.caminho, r.corpo).problema(http.StatusPreconditionRequired, "pre_condicao_obrigatoria")
				p := comT(c.admA).req(r.metodo, r.caminho, r.corpo, cab{"If-Match": `"999"`}).problema(http.StatusPreconditionFailed, "versao_desatualizada")
				if p["recurso_atual"] == nil || p["versao_atual"] == nil {
					t.Fatalf("412 sem recurso_atual: %v", p)
				}
			}
		})
	}
}

// Caminho feliz de cada rota do ADM, pela API.
func TestRotasADMCaminhoFeliz(t *testing.T) {
	adm := entrar(t, novaConta(t, "profissional", ""), "")

	// Empresa.
	r := adm.get("/v1/empresa")
	e := r.ok(http.StatusOK)
	if e["situacao"] == nil || e["plano"] != "profissional" || e["nome_fantasia"] != nil {
		t.Fatalf("empresa: %v", e)
	}
	e = adm.req("PATCH", "/v1/empresa", map[string]any{"nome_fantasia": "Transteste"}, cab{"If-Match": r.etag()}).ok(http.StatusOK)
	if e["nome_fantasia"] != "Transteste" || e["razao_social"] != "Transportes Teste Ltda" {
		t.Fatalf("PATCH empresa: %v", e)
	}
	r = adm.get("/v1/empresa")
	if e = adm.req("PATCH", "/v1/empresa", map[string]any{"nome_fantasia": nil}, cab{"If-Match": r.etag(), "Content-Type": "application/merge-patch+json"}).ok(http.StatusOK); e["nome_fantasia"] != nil {
		t.Fatalf("merge patch com null: %v", e)
	}

	// Unidades, turnos e feriados.
	lista := adm.get("/v1/unidades").ok(http.StatusOK)
	if len(lista["dados"].([]any)) != 1 || lista["paginacao"].(map[string]any)["tem_mais"] != false {
		t.Fatalf("lista de unidades: %v", lista)
	}
	cr := adm.post("/v1/unidades", map[string]any{"codigo": "FIL", "nome": "Filial", "endereco": map[string]any{"municipio": "Campinas", "uf": "SP"},
		"turnos": []any{map[string]any{"dia_semana": 1, "hora_inicio": "07:00", "hora_fim": "16:00"}}})
	u := cr.ok(http.StatusCreated)
	id := texto(u, "id")
	if cr.cab.Get("Location") != "/v1/unidades/"+id || len(u["turnos"].([]any)) != 1 {
		t.Fatalf("criar unidade: %v %v", cr.cab, u)
	}
	r = adm.get("/v1/unidades/" + id)
	u = adm.req("PATCH", "/v1/unidades/"+id, map[string]any{"endereco": map[string]any{"bairro": "Centro"}}, cab{"If-Match": r.etag()}).ok(http.StatusOK)
	end := u["endereco"].(map[string]any)
	if end["municipio"] != "Campinas" || end["bairro"] != "Centro" {
		t.Fatalf("merge do endereço: %v", end)
	}
	tr := adm.post("/v1/unidades/"+id+"/turnos", map[string]any{"dia_semana": 2, "hora_inicio": "07:00", "hora_fim": "16:00"})
	turno := tr.ok(http.StatusCreated)
	tid := texto(turno, "id")
	tp := adm.req("PATCH", "/v1/unidades/"+id+"/turnos/"+tid, map[string]any{"hora_fim": "17:00", "nome": "Manhã"}, cab{"If-Match": tr.etag()})
	if turno = tp.ok(http.StatusOK); turno["hora_fim"] != "17:00" || turno["nome"] != "Manhã" {
		t.Fatalf("turno: %v", turno)
	}
	adm.req("DELETE", "/v1/unidades/"+id+"/turnos/"+tid, nil, cab{"If-Match": tp.etag()}).ok(http.StatusNoContent)
	fr := adm.post("/v1/unidades/"+id+"/feriados", map[string]any{"data": "2027-11-15", "descricao": "Proclamação"})
	f := fr.ok(http.StatusCreated)
	fl := adm.get("/v1/unidades/" + id + "/feriados?data[gte]=2027-01-01&data[lte]=2027-12-31").ok(http.StatusOK)
	if len(fl["dados"].([]any)) != 1 {
		t.Fatalf("feriados: %v", fl)
	}
	adm.req("DELETE", "/v1/unidades/"+id+"/feriados/"+texto(f, "id"), nil, cab{"If-Match": fr.etag()}).ok(http.StatusNoContent)
	r = adm.get("/v1/unidades/" + id)
	adm.req("DELETE", "/v1/unidades/"+id, nil, cab{"If-Match": r.etag()}).ok(http.StatusNoContent)
	r = adm.get("/v1/unidades/" + id)
	if u = r.ok(http.StatusOK); u["situacao"] != "inativa" {
		t.Fatalf("unidade inativada: %v", u)
	}
	u = adm.req("POST", "/v1/unidades/"+id+"/reativar", nil, cab{"If-Match": r.etag()}).ok(http.StatusOK)
	if u["situacao"] != "ativa" {
		t.Fatalf("unidade reativada: %v", u)
	}

	// Pessoas.
	pr := adm.post("/v1/pessoas", map[string]any{"nome": "Joana Técnica", "cpf": "529.982.247-25", "vinculo": "empregado"})
	p := pr.ok(http.StatusCreated)
	pid := texto(p, "id")
	pl := adm.get("/v1/pessoas?busca=Joana&com_usuario=false&limite=10").ok(http.StatusOK)
	pd := pl["dados"].([]any)
	if len(pd) != 1 || strings.Contains(texto(pd[0].(map[string]any), "cpf_mascarado"), "529") {
		t.Fatalf("lista de pessoas: %v", pl)
	}
	r = adm.get("/v1/pessoas/" + pid)
	if p = r.ok(http.StatusOK); p["cpf"] != "52998224725" {
		t.Fatalf("detalhe com CPF inteiro: %v", p)
	}
	p = adm.req("PATCH", "/v1/pessoas/"+pid, map[string]any{"cargo": "Eletricista"}, cab{"If-Match": r.etag()}).ok(http.StatusOK)
	r = adm.get("/v1/pessoas/" + pid)
	adm.req("DELETE", "/v1/pessoas/"+pid, nil, cab{"If-Match": r.etag()}).ok(http.StatusNoContent)
	r = adm.get("/v1/pessoas/" + pid)
	adm.req("POST", "/v1/pessoas/"+pid+"/reativar", nil, cab{"If-Match": r.etag()}).ok(http.StatusOK)

	// Paginação por cursor.
	for i := 0; i < 3; i++ {
		adm.post("/v1/pessoas", map[string]any{"nome": "Pessoa paginada"}).ok(http.StatusCreated)
	}
	p1 := adm.get("/v1/pessoas?busca=paginada&limite=2").ok(http.StatusOK)
	pag := p1["paginacao"].(map[string]any)
	if len(p1["dados"].([]any)) != 2 || pag["tem_mais"] != true || pag["proximo_cursor"] == nil {
		t.Fatalf("primeira página: %v", p1)
	}
	p2 := adm.get("/v1/pessoas?busca=paginada&limite=2&cursor=" + texto(pag, "proximo_cursor")).ok(http.StatusOK)
	if len(p2["dados"].([]any)) != 1 || p2["paginacao"].(map[string]any)["proximo_cursor"] != nil {
		t.Fatalf("segunda página: %v", p2)
	}
	adm.get("/v1/pessoas?busca=outra&limite=2&cursor="+texto(pag, "proximo_cursor")).problema(http.StatusBadRequest, "cursor_invalido")

	// Usuários: com pessoa nova, perfis, permissão avulsa, sessões.
	tec, tecID, _ := convidar(t, adm, []string{"tecnico"}, map[string]any{"custo_hora": "85.5", "nova_pessoa": map[string]any{"nome": "Téc"}})
	ul := adm.get("/v1/usuarios?perfil=tecnico").ok(http.StatusOK)
	if len(ul["dados"].([]any)) != 1 {
		t.Fatalf("lista de usuários: %v", ul)
	}
	r = adm.get("/v1/usuarios/" + tecID)
	us := r.ok(http.StatusOK)
	if us["custo_hora"] != "85.50" || us["situacao"] != "ativo" || us["pessoa_id"] == nil || us["escopo_acesso"] != "interno" {
		t.Fatalf("usuário: %v", us)
	}
	us = adm.req("PATCH", "/v1/usuarios/"+tecID, map[string]any{"perfis": []string{"tecnico", "almoxarife"}, "telefone": "11988887777"}, cab{"If-Match": r.etag()}).ok(http.StatusOK)
	if len(us["perfis"].([]any)) != 2 || us["telefone"] != "11988887777" {
		t.Fatalf("PATCH usuário: %v", us)
	}
	pc := adm.post("/v1/usuarios/"+tecID+"/permissoes", map[string]any{"permissao": "os.encerrar"})
	if us = pc.ok(http.StatusCreated); len(us["permissoes_avulsas"].([]any)) != 1 || !strings.HasSuffix(pc.cab.Get("Location"), "/permissoes/os.encerrar") {
		t.Fatalf("permissão avulsa: %v", us)
	}
	adm.req("DELETE", "/v1/usuarios/"+tecID+"/permissoes/os.encerrar", nil).ok(http.StatusNoContent)
	es := adm.post("/v1/usuarios/"+tecID+"/encerrar-sessoes", nil).ok(http.StatusOK)
	if es["sessoes_encerradas"] != float64(1) {
		t.Fatalf("encerrar sessões: %v", es)
	}
	tec.get("/v1/usuarios/eu").problema(http.StatusUnauthorized, "credencial_invalida")
	r = adm.get("/v1/usuarios/" + tecID)
	adm.req("POST", "/v1/usuarios/"+tecID+"/desativar", nil, cab{"If-Match": r.etag()}).ok(http.StatusOK)
	r = adm.get("/v1/usuarios/" + tecID)
	adm.req("POST", "/v1/usuarios/"+tecID+"/desativar", nil, cab{"If-Match": r.etag()}).problema(http.StatusConflict, "transicao_invalida")
	adm.req("POST", "/v1/usuarios/"+tecID+"/reativar", nil, cab{"If-Match": r.etag()}).ok(http.StatusOK)
	novo := adm.post("/v1/usuarios", map[string]any{"nome": "Convidado", "email": emailAleatorio("c"), "perfis": []string{"solicitante"}}).ok(http.StatusCreated)
	adm.post("/v1/usuarios/"+texto(novo, "id")+"/reenviar-convite", nil).ok(http.StatusNoContent)

	// Último administrador.
	eu := adm.get("/v1/usuarios/eu")
	euID := texto(eu.ok(http.StatusOK), "id")
	adm.req("POST", "/v1/usuarios/"+euID+"/desativar", nil, cab{"If-Match": eu.etag()}).problema(http.StatusConflict, "ultimo_administrador")
	eu = adm.get("/v1/usuarios/eu")
	if x := adm.req("PATCH", "/v1/usuarios/eu", map[string]any{"nome": "Ana A."}, cab{"If-Match": eu.etag()}).ok(http.StatusOK); x["nome"] != "Ana A." {
		t.Fatalf("PATCH eu: %v", x)
	}

	// Perfis, permissões e trilha.
	pf := adm.get("/v1/perfis").ok(http.StatusOK)
	if len(pf["dados"].([]any)) != 12 {
		t.Fatalf("perfis: %v", len(pf["dados"].([]any)))
	}
	cat := adm.get("/v1/permissoes").ok(http.StatusOK)
	achou := false
	for _, x := range cat["dados"].([]any) {
		achou = achou || x.(map[string]any)["codigo"] == "os.encerrar"
	}
	if !achou {
		t.Fatal("catálogo sem os.encerrar")
	}
	tl := adm.get("/v1/eventos-auditoria?entidade=pessoa&entidade_id=" + pid).ok(http.StatusOK)
	if len(tl["dados"].([]any)) == 0 {
		t.Fatalf("trilha da pessoa: %v", tl)
	}
	ev := adm.get("/v1/eventos-auditoria/" + texto(tl["dados"].([]any)[0].(map[string]any), "id")).ok(http.StatusOK)
	if _, tem := ev["antes"]; !tem {
		t.Fatalf("evento sem antes: %v", ev)
	}
}

// Regras transversais: empresa_id, campo desconhecido, campo somente
// leitura, filtro não declarado, origem de outro site, idempotência.
func TestRegrasTransversais(t *testing.T) {
	adm := entrar(t, novaConta(t, "essencial", ""), "")

	adm.get("/v1/pessoas?empresa_id="+uuid.NewString()).problema(http.StatusBadRequest, "campo_nao_permitido")
	adm.post("/v1/pessoas", map[string]any{"nome": "X", "empresa_id": uuid.NewString()}).problema(http.StatusBadRequest, "campo_nao_permitido")
	adm.post("/v1/pessoas", map[string]any{"nome": "X", "apelido": "Y"}).problema(http.StatusBadRequest, "requisicao_invalida")
	adm.post("/v1/pessoas", `{"nome":`).problema(http.StatusBadRequest, "requisicao_invalida")
	v := adm.post("/v1/pessoas", map[string]any{"vinculo": "outro"}).problema(http.StatusUnprocessableEntity, "validacao_falhou")
	if len(v["erros"].([]any)) != 2 {
		t.Fatalf("erros por campo: %v", v)
	}
	for _, x := range v["erros"].([]any) {
		m := x.(map[string]any)
		if m["campo"] == "" || m["codigo"] == "" || m["mensagem"] == "" {
			t.Fatalf("erro de campo fora do formato: %v", m)
		}
	}
	r := adm.get("/v1/empresa")
	adm.req("PATCH", "/v1/empresa", map[string]any{"situacao": "ativa"}, cab{"If-Match": r.etag()}).problema(http.StatusUnprocessableEntity, "campo_somente_leitura")
	adm.get("/v1/pessoas?ordenar=-nome").problema(http.StatusBadRequest, "filtro_nao_suportado")
	adm.get("/v1/nao-existe").problema(http.StatusNotFound, "recurso_nao_encontrado")
	adm.get("/v1/pessoas/nao-e-uuid").problema(http.StatusNotFound, "recurso_nao_encontrado")

	// Login vindo de outro site é recusado (login CSRF).
	outro := anonimo(t)
	outro.origem = "https://atacante.exemplo"
	outro.post("/v1/auth/login", map[string]any{"email": "a@b.c", "senha": "x"}).problema(http.StatusForbidden, "csrf_invalido")

	// Segunda unidade no Essencial: plano_necessario.
	p := adm.post("/v1/unidades", map[string]any{"codigo": "F2", "nome": "Filial", "turnos": turnoUtil}).problema(http.StatusForbidden, "plano_nao_inclui_recurso")
	if p["plano_necessario"] != "profissional" {
		t.Fatalf("plano_necessario: %v", p)
	}

	// Idempotência: mesmo corpo, uma pessoa; corpo diferente, 422.
	chave := cab{"Idempotency-Key": uuid.NewString()}
	a := adm.post("/v1/pessoas", map[string]any{"nome": "Idempotente"}, chave)
	pa := a.ok(http.StatusCreated)
	b := adm.post("/v1/pessoas", map[string]any{"nome": "Idempotente"}, chave)
	pb := b.ok(http.StatusCreated)
	if pa["id"] != pb["id"] || b.cab.Get("Idempotency-Replayed") != "true" || b.cab.Get("Location") != a.cab.Get("Location") {
		t.Fatalf("repetição: %v %v %v", pa["id"], pb["id"], b.cab)
	}
	l := adm.get("/v1/pessoas?busca=Idempotente").ok(http.StatusOK)
	if len(l["dados"].([]any)) != 1 {
		t.Fatalf("idempotência criou duas pessoas: %v", l)
	}
	adm.post("/v1/pessoas", map[string]any{"nome": "Outra"}, chave).problema(http.StatusUnprocessableEntity, "idempotencia_conflito")

	// Versão do produto em toda resposta.
	if anonimo(t).get("/v1/versao").cab.Get("Tenaz-Versao") == "" {
		t.Fatal("sem Tenaz-Versao")
	}
}
