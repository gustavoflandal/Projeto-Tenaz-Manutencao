package integracao_test

import (
	"encoding/json"
	"net/http"
	"reflect"
	"testing"
)

// POST /v1/auth/login responde exatamente igual — status, corpo e
// cabeçalhos de autenticação — para senha errada, conta bloqueada (com a
// senha certa ou errada) e e-mail inexistente (documento 9, seção 5.3). Só
// id_rastreio difere, porque é o da requisição.
func TestLoginMesmaRespostaPelaAPI(t *testing.T) {
	em := novaConta(t, "essencial", "")
	tentar := func(email, senha string) (int, map[string]any, http.Header) {
		r := anonimo(t).post("/v1/auth/login", map[string]any{"email": email, "senha": senha})
		r.problema(http.StatusUnauthorized, "credencial_invalida")
		var m map[string]any
		if err := json.Unmarshal(r.corpo, &m); err != nil {
			t.Fatal(err)
		}
		if m["id_rastreio"] == nil {
			t.Fatalf("sem id_rastreio: %s", r.corpo)
		}
		delete(m, "id_rastreio")
		delete(m, "instance")
		return r.status, m, r.cab
	}
	stErrada, errada, cabErrada := tentar(em, "errada-errada-1")
	stInex, inexistente, _ := tentar(emailAleatorio("ninguem"), senhaBoa)
	for i := 0; i < 4; i++ {
		tentar(em, "errada-errada-1")
	}
	stCerta, bloqueadaCerta, cabBloq := tentar(em, senhaBoa)
	stBloqErrada, bloqueadaErrada, _ := tentar(em, "errada-errada-1")

	for nome, caso := range map[string]struct {
		status int
		corpo  map[string]any
	}{
		"e-mail inexistente":           {stInex, inexistente},
		"bloqueada com a senha certa":  {stCerta, bloqueadaCerta},
		"bloqueada com a senha errada": {stBloqErrada, bloqueadaErrada},
	} {
		if caso.status != stErrada || !reflect.DeepEqual(caso.corpo, errada) {
			t.Errorf("%s responde diferente da senha errada:\n %d %v\n %d %v", nome, caso.status, caso.corpo, stErrada, errada)
		}
	}
	// Nenhum cabeçalho denuncia o bloqueio (Retry-After, cookie).
	for _, c := range []string{"Retry-After", "Set-Cookie", "WWW-Authenticate"} {
		if cabBloq.Get(c) != cabErrada.Get(c) {
			t.Errorf("cabeçalho %s difere: bloqueada %q, senha errada %q", c, cabBloq.Get(c), cabErrada.Get(c))
		}
	}
}
