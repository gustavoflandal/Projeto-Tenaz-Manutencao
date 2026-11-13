package vazamento

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/danielgtaylor/huma/v2"

	"tenaz/internal/servidor"
)

// Provas de detector da parte 3. Cada prova monta, no próprio teste, um
// servidor que embrulha a API real com um defeito de isolamento (o repositório
// sem o filtro da empresa, a camada HTTP sem a conferência do escopo, um
// oráculo no erro, uma escrita silenciosa em B) e exige que a MESMA
// verificação da varredura o acuse. Sem isto, uma verificação errada passaria
// por vacuidade.

// sondadorEstragado aponta o sondador para a API real embrulhada no defeito.
// defeito devolve verdadeiro quando já respondeu.
func sondadorEstragado(t *testing.T, defeito func(amb *ambienteAPI, w http.ResponseWriter, r *http.Request) bool) (*sondador, *ambienteAPI) {
	t.Helper()
	amb := ambienteDaAPI(t)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !defeito(amb, w, r) {
			amb.handler.ServeHTTP(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	cls, err := novoClassificador(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	return &sondador{ctx: context.Background(), base: srv.URL, amb: amb, cen: amb.cen, cls: cls}, amb
}

// repassarComo leva a requisição à API real com a credencial de outro
// cliente e devolve a resposta gravada.
func repassarComo(amb *ambienteAPI, r *http.Request, c *clienteAPI, corpo []byte) *httptest.ResponseRecorder {
	r2 := r.Clone(r.Context())
	r2.Header.Del("Cookie")
	r2.Header.Del("Authorization")
	r2.Header.Del("X-CSRF-Token")
	if c != nil && c.cookie != "" {
		r2.AddCookie(&http.Cookie{Name: "__Host-tenaz_sessao", Value: c.cookie})
		r2.Header.Set("X-CSRF-Token", c.csrf)
	}
	if corpo != nil {
		r2.Body = io.NopCloser(bytes.NewReader(corpo))
		r2.ContentLength = int64(len(corpo))
	}
	g := httptest.NewRecorder()
	amb.handler.ServeHTTP(g, r2)
	return g
}

// escreverGravado copia a resposta gravada, com o corpo trocado.
func escreverGravado(w http.ResponseWriter, g *httptest.ResponseRecorder, corpo []byte) {
	for k, v := range g.Header() {
		if k != "Content-Length" {
			w.Header()[k] = v
		}
	}
	w.WriteHeader(g.Code)
	_, _ = w.Write(corpo)
}

func lerCorpo(r *http.Request) []byte {
	if r.Body == nil {
		return nil
	}
	b, _ := io.ReadAll(r.Body)
	return b
}

func temCookie(r *http.Request, c *clienteAPI) bool {
	k, err := r.Cookie("__Host-tenaz_sessao")
	return err == nil && c != nil && k.Value == c.cookie
}

// planoDe devolve o plano de uma rota do contrato.
func planoDe(t *testing.T, metodo, caminho string) planoDaRota {
	t.Helper()
	oapi, err := lerContratoVersionado()
	if err != nil {
		t.Fatal(err)
	}
	ops, err := operacoesDoContrato(oapi)
	if err != nil {
		t.Fatal(err)
	}
	planos, faltas := planejar(ops)
	if len(faltas) > 0 {
		t.Fatalf("plano com faltas: %v", faltas)
	}
	for _, p := range planos {
		if p.op.Metodo == metodo && p.op.Caminho == caminho {
			return p
		}
	}
	t.Fatalf("rota %s %s fora do contrato", metodo, caminho)
	return planoDaRota{}
}

// exigirAcusacao confere que a verificação acusou o defeito, com o trecho
// esperado na mensagem.
func exigirAcusacao(t *testing.T, defeito string, p problemas, trechos ...string) {
	t.Helper()
	if len(p) == 0 {
		t.Fatalf("o detector não acusou: %s", defeito)
	}
	tudo := strings.Join(p, "\n")
	for _, tr := range trechos {
		if !strings.Contains(tudo, tr) {
			t.Errorf("%s: a acusação não diz %q:\n%s", defeito, tr, tudo)
		}
	}
	if strings.Contains(tudo, "falha da varredura") {
		t.Errorf("%s: a prova falhou por outro motivo, não pelo defeito:\n%s", defeito, tudo)
	}
}

// O repositório sem o filtro da empresa: o detalhe da pessoa de B sai para A.
func TestDetectorAPIAcusaDetalheDeOutraEmpresa(t *testing.T) {
	s, amb := sondadorEstragado(t, func(amb *ambienteAPI, w http.ResponseWriter, r *http.Request) bool {
		if r.Method == http.MethodGet && r.URL.Path == "/v1/pessoas/"+amb.cen.B.ids["pessoas"] {
			g := repassarComo(amb, r, amb.cen.B.adm, nil)
			escreverGravado(w, g, g.Body.Bytes())
			return true
		}
		return false
	})
	cen := amb.cen
	p := s.idDeOutraEmpresa(planoDe(t, "GET", "/v1/pessoas/{id}"), cen.A.adm, cen.A, cen.B, s.atorInterno(cen.A))
	exigirAcusacao(t, "detalhe de B devolvido a A", p, "esperado 404", "oráculo", "alheio a administrador de A")
}

// Um 403 para o id de B confirma que ele existe.
func TestDetectorAPIAcusaProibidoParaIdDeOutraEmpresa(t *testing.T) {
	s, amb := sondadorEstragado(t, func(amb *ambienteAPI, w http.ResponseWriter, r *http.Request) bool {
		if strings.HasSuffix(r.URL.Path, "/"+amb.cen.B.ids["unidades"]) {
			w.Header().Set("Content-Type", "application/problem+json")
			w.WriteHeader(http.StatusForbidden)
			_, _ = w.Write([]byte(`{"type":"https://api.tenaz.com.br/erros/permissao_negada","title":"Permissão negada","status":403,"codigo":"permissao_negada"}`))
			return true
		}
		return false
	})
	cen := amb.cen
	for _, rota := range [][2]string{{"GET", "/v1/unidades/{id}"}, {"DELETE", "/v1/unidades/{id}"}} {
		p := s.idDeOutraEmpresa(planoDe(t, rota[0], rota[1]), cen.A.adm, cen.A, cen.B, s.atorInterno(cen.A))
		exigirAcusacao(t, "403 para id de B em "+rota[0], p, "esperado 404", "oráculo")
	}
}

// O 404 certo, mas com um detalhe que só o id de B recebe.
func TestDetectorAPIAcusaOraculoNoDetalhe(t *testing.T) {
	s, amb := sondadorEstragado(t, func(amb *ambienteAPI, w http.ResponseWriter, r *http.Request) bool {
		if strings.Contains(r.URL.Path, amb.cen.B.ids["turnos"]) {
			g := repassarComo(amb, r, amb.cen.A.adm, lerCorpo(r))
			var m map[string]any
			_ = json.Unmarshal(g.Body.Bytes(), &m)
			m["detail"] = "O turno pertence a outra empresa."
			b, _ := json.Marshal(m)
			escreverGravado(w, g, b)
			return true
		}
		return false
	})
	cen := amb.cen
	p := s.idDeOutraEmpresa(planoDe(t, "PATCH", "/v1/unidades/{id}/turnos/{turno_id}"), cen.A.adm, cen.A, cen.B, s.atorInterno(cen.A))
	exigirAcusacao(t, "detalhe diferente para o turno de B", p, "oráculo", "{turno_id} de B")
}

// A lista de A com uma linha de B no meio.
func TestDetectorAPIAcusaListaComLinhaDeOutraEmpresa(t *testing.T) {
	s, amb := sondadorEstragado(t, func(amb *ambienteAPI, w http.ResponseWriter, r *http.Request) bool {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/unidades" || !temCookie(r, amb.cen.A.adm) {
			return false
		}
		g := repassarComo(amb, r, amb.cen.A.adm, nil)
		gB := repassarComo(amb, r, amb.cen.B.adm, nil)
		var lA, lB map[string]any
		_ = json.Unmarshal(g.Body.Bytes(), &lA)
		_ = json.Unmarshal(gB.Body.Bytes(), &lB)
		lA["dados"] = append(lA["dados"].([]any), lB["dados"].([]any)[0])
		b, _ := json.Marshal(lA)
		escreverGravado(w, g, b)
		return true
	})
	cen := amb.cen
	p := s.conteudoProprio(planoDe(t, "GET", "/v1/unidades"), cen.A.adm, cen.A, s.atorInterno(cen.A))
	exigirAcusacao(t, "linha de B na lista de A", p, "alheio a administrador de A")
}

// A escrita que alcança B e responde o 404 de um id inexistente: a resposta
// é perfeita, e só o retrato do banco a denuncia.
func TestDetectorAPIAcusaEscritaSilenciosaEmOutraEmpresa(t *testing.T) {
	amb := ambienteDaAPI(t)
	cen := amb.cen
	ctx := context.Background()
	r, err := cen.B.adm.fazer(amb.srv.URL, http.MethodPost, "/v1/pessoas", map[string]any{"nome": "Pessoa da prova " + cen.B.sentinela}, nil)
	if err != nil || r.status != http.StatusCreated {
		t.Fatalf("pessoa da prova em B: %v %v", err, r)
	}
	alvo := idDe(r.json())

	s, _ := sondadorEstragado(t, func(amb *ambienteAPI, w http.ResponseWriter, r *http.Request) bool {
		if r.Method != http.MethodDelete || r.URL.Path != "/v1/pessoas/"+alvo {
			return false
		}
		_ = repassarComo(amb, r, amb.cen.B.adm, nil) // a escrita em B
		r2 := r.Clone(r.Context())
		r2.URL.Path = "/v1/pessoas/" + idInexistente()
		g := repassarComo(amb, r2, amb.cen.A.adm, nil)
		var m map[string]any
		_ = json.Unmarshal(g.Body.Bytes(), &m)
		m["instance"] = r.URL.Path
		b, _ := json.Marshal(m)
		escreverGravado(w, g, b)
		return true
	})
	bProva := *cen.B
	bProva.ids = map[string]string{}
	for k, v := range cen.B.ids {
		bProva.ids[k] = v
	}
	bProva.ids["pessoas"] = alvo
	cenProva := *cen
	cenProva.B = &bProva
	s.cen = &cenProva

	amb.aut.Aguardar()
	antes, err := retratoDaEmpresa(ctx, pool, cen.B.ID)
	if err != nil {
		t.Fatal(err)
	}
	if p := s.idDeOutraEmpresa(planoDe(t, "DELETE", "/v1/pessoas/{id}"), cen.A.adm, cen.A, &bProva, s.atorInterno(cen.A)); len(p) > 0 {
		t.Fatalf("a prova precisa de uma resposta perfeita, para que só o retrato acuse; veio:\n%s", strings.Join(p, "\n"))
	}
	amb.aut.Aguardar()
	depois, err := retratoDaEmpresa(ctx, pool, cen.B.ID)
	if err != nil {
		t.Fatal(err)
	}
	exigirAcusacao(t, "inativação silenciosa da pessoa de B", diferencasDoRetrato(antes, depois), "app.pessoa")
}

// A configuração de captura de B devolvida a A: booleanos, sem UUID nem
// texto, que só a verificação específica da rota vê.
func TestDetectorAPIAcusaCapturaDeOutraEmpresa(t *testing.T) {
	s, amb := sondadorEstragado(t, func(amb *ambienteAPI, w http.ResponseWriter, r *http.Request) bool {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/capturas-auditoria" || !temCookie(r, amb.cen.A.adm) {
			return false
		}
		g := repassarComo(amb, r, amb.cen.B.adm, nil)
		escreverGravado(w, g, g.Body.Bytes())
		return true
	})
	cen := amb.cen
	p := s.conteudoProprio(planoDe(t, "GET", "/v1/capturas-auditoria"), cen.A.adm, cen.A, s.atorInterno(cen.A))
	exigirAcusacao(t, "captura de B devolvida a A", p, "mostra a configuração de captura de B")
}

// Referência a registro de B aceita no corpo.
func TestDetectorAPIAcusaReferenciaDeOutraEmpresaAceita(t *testing.T) {
	s, amb := sondadorEstragado(t, func(amb *ambienteAPI, w http.ResponseWriter, r *http.Request) bool {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/usuarios" {
			return false
		}
		corpo := lerCorpo(r)
		if !bytes.Contains(corpo, []byte(amb.cen.B.ids["pessoas"])) {
			r.Body = io.NopCloser(bytes.NewReader(corpo))
			return false
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{"id":"` + idInexistente() + `","pessoa_id":"` + amb.cen.B.ids["pessoas"] + `"}`))
		return true
	})
	cen := amb.cen
	p := s.referenciaNoCorpo(planoDe(t, "POST", "/v1/usuarios"), cen.A.adm, cen.A, s.atorInterno(cen.A), nil)
	exigirAcusacao(t, "pessoa de B aceita em POST /v1/usuarios", p, "aceitou a referência", "oráculo")
}

// empresa_id descartado em silêncio, em vez de recusado.
func TestDetectorAPIAcusaEmpresaIDAceito(t *testing.T) {
	s, amb := sondadorEstragado(t, func(amb *ambienteAPI, w http.ResponseWriter, r *http.Request) bool {
		if r.Method != http.MethodPost || r.URL.Path != "/v1/pessoas" {
			return false
		}
		q := r.URL.Query()
		q.Del("empresa_id")
		r.URL.RawQuery = q.Encode()
		var m map[string]any
		corpo := lerCorpo(r)
		if json.Unmarshal(corpo, &m) == nil {
			delete(m, "empresa_id")
			corpo, _ = json.Marshal(m)
		}
		g := repassarComo(amb, r, amb.cen.A.adm, corpo)
		escreverGravado(w, g, g.Body.Bytes())
		return true
	})
	cen := amb.cen
	p := s.empresaIDNaRequisicao(planoDe(t, "POST", "/v1/pessoas"), cen.A.adm, cen.A)
	exigirAcusacao(t, "empresa_id descartado", p, "?empresa_id= por", "com empresa_id no corpo")
}

// Rota interna atendida para o prestador (escopo não conferido).
func TestDetectorAPIAcusaExternoEmRotaInterna(t *testing.T) {
	s, amb := sondadorEstragado(t, func(amb *ambienteAPI, w http.ResponseWriter, r *http.Request) bool {
		if r.URL.Path != "/v1/pessoas" || !temCookie(r, amb.cen.F1.web) {
			return false
		}
		g := repassarComo(amb, r, amb.cen.A.adm, lerCorpo(r))
		escreverGravado(w, g, g.Body.Bytes())
		return true
	})
	cen := amb.cen
	a, err := s.atorExterno(cen.F1)
	if err != nil {
		t.Fatal(err)
	}
	p := s.externoFora(planoDe(t, "GET", "/v1/pessoas"), cen.F1, cen.F1.web, a)
	exigirAcusacao(t, "prestador lê /v1/pessoas", p, "esperado 404", "rota inexistente", "alheio a F1")
}

// Rota do próprio usuário que traz um colega do mesmo fornecedor... de outro
// fornecedor (F2).
func TestDetectorAPIAcusaExternoQueVeAlheio(t *testing.T) {
	s, amb := sondadorEstragado(t, func(amb *ambienteAPI, w http.ResponseWriter, r *http.Request) bool {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/usuarios/eu" || !temCookie(r, amb.cen.F1.web) {
			return false
		}
		g := repassarComo(amb, r, amb.cen.F1.web, nil)
		var m map[string]any
		_ = json.Unmarshal(g.Body.Bytes(), &m)
		m["colega_id"] = amb.cen.F2.usuario
		b, _ := json.Marshal(m)
		escreverGravado(w, g, b)
		return true
	})
	cen := amb.cen
	a, err := s.atorExterno(cen.F1)
	if err != nil {
		t.Fatal(err)
	}
	p := s.externoDentro(planoDe(t, "GET", "/v1/usuarios/eu"), cen.F1, cen.F1.web, a)
	exigirAcusacao(t, "F1 vê o usuário de F2", p, cen.F2.usuario)
}

// Rota autenticada que responde sem credencial.
func TestDetectorAPIAcusaRotaSemCredencial(t *testing.T) {
	s, amb := sondadorEstragado(t, func(amb *ambienteAPI, w http.ResponseWriter, r *http.Request) bool {
		if r.URL.Path != "/v1/unidades" || r.Header.Get("Cookie") != "" {
			return false
		}
		g := repassarComo(amb, r, amb.cen.A.adm, nil)
		escreverGravado(w, g, g.Body.Bytes())
		return true
	})
	exigirAcusacao(t, "GET /v1/unidades sem credencial", s.semCredencial(planoDe(t, "GET", "/v1/unidades")), "esperado 401")
	_ = amb
}

// Rota nova que a suíte não sabe testar quebra o build, com a instrução.
func TestDetectorAPIAcusaRotaSemCaso(t *testing.T) {
	corpo := &huma.Schema{Type: "object"}
	ops := []operacao{
		{ID: "obterCoisa", Metodo: "GET", Caminho: "/v1/coisas/{id}", Escopos: []string{"interno"}, Sessao: true,
			Caminhos: []parametroCaminho{{Nome: "id", Segmento: "coisas"}}},
		{ID: "criarCoisa", Metodo: "POST", Caminho: "/v1/coisas", Escopos: []string{"interno"}, Sessao: true, Corpo: corpo},
		{ID: "criarPessoa", Metodo: "POST", Caminho: "/v1/pessoas", Escopos: []string{"interno"}, Sessao: true, Corpo: corpo,
			CamposUUID: []string{"ativo_id"}},
		{ID: "listarCoisas", Metodo: "GET", Caminho: "/v1/coisas", Escopos: []string{"interno"}, Sessao: true,
			Consulta: map[string]*huma.Schema{"ativo_id": {Type: "string", Format: "uuid"}}},
	}
	_, faltas := planejar(ops)
	exigirAcusacao(t, "rota sem caso", faltas,
		"vem depois de /coisas/, que não está em recursosDeCaminho",
		"POST /v1/coisas (criarCoisa): rota autenticada com corpo e sem caso",
		`campo uuid "ativo_id" do corpo sem classificação`,
		`parâmetro de consulta uuid "ativo_id" sem classificação`,
		`casosDeCorpo["alterarPessoa"]: nenhuma rota autenticada com corpo tem esse operationId`)
}

// Contrato versionado desatualizado, e rota do servidor fora do contrato.
func TestDetectorAPIAcusaContratoDesatualizadoERotaForaDoContrato(t *testing.T) {
	dados, err := os.ReadFile(arquivoContrato)
	if err != nil {
		t.Fatal(err)
	}
	velho := filepath.Join(t.TempDir(), "openapi.yaml")
	if err := os.WriteFile(velho, bytes.Replace(dados, []byte("/v1/pessoas/{id}/reativar:"), []byte("/v1/pessoas/{id}/reativar-velho:"), 1), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := lerContrato(velho); err == nil || !strings.Contains(err.Error(), "difere do contrato gerado") {
		t.Fatalf("contrato desatualizado não acusado: %v", err)
	}

	oapi, err := lerContratoVersionado()
	if err != nil {
		t.Fatal(err)
	}
	montado := servidor.Contrato()
	huma.Get(montado.API(), "/v1/rota-extra-da-prova", func(context.Context, *struct{}) (*struct{}, error) { return nil, nil })
	exigirAcusacao(t, "rota fora do contrato", mesmasRotas(oapi, montado.API().OpenAPI()),
		"rota atendida pelo servidor e ausente do contrato: GET /v1/rota-extra-da-prova")
}

// O classificador: UUID de B é alheio a A; o da própria pessoa não. Para o
// externo, até o registro interno da própria empresa é alheio.
func TestDetectorAPIClassificaUUIDs(t *testing.T) {
	amb := ambienteDaAPI(t)
	cen := amb.cen
	ctx := context.Background()
	s := &sondador{ctx: ctx, cen: cen}
	cls, err := novoClassificador(ctx)
	if err != nil {
		t.Fatal(err)
	}
	f1, err := s.atorExterno(cen.F1)
	if err != nil {
		t.Fatal(err)
	}
	casos := []struct {
		rotulo string
		a      atorAPI
		u      string
		alheio bool
	}{
		{"pessoa de B para A", s.atorInterno(cen.A), cen.B.ids["pessoas"], true},
		{"pessoa de A para A", s.atorInterno(cen.A), cen.A.ids["pessoas"], false},
		{"fornecedor de F1, também usado em B, para A", s.atorInterno(cen.A), cen.fornecedorF1, false},
		{"pessoa interna de A para F1", f1, cen.A.ids["pessoas"], true},
		{"usuário de F2 para F1", f1, cen.F2.usuario, true},
		{"o próprio F1", f1, cen.F1.usuario, false},
		{"a empresa de F1", f1, cen.A.ID, false},
		{"o fornecedor de F1", f1, cen.fornecedorF1, false},
	}
	for _, c := range casos {
		fora, err := cls.alheios(ctx, c.a, []string{strings.ToLower(c.u)})
		if err != nil {
			t.Fatal(err)
		}
		if (len(fora) > 0) != c.alheio {
			t.Errorf("%s: alheio = %v, esperado %v", c.rotulo, len(fora) > 0, c.alheio)
		}
	}
}
