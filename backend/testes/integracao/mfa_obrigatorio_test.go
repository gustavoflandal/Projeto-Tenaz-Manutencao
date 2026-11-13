package integracao_test

import (
	"net/http"
	"testing"
	"time"
)

// RN-ADM-027: com mfa_obrigatorio ligado, quem não tem MFA fica com uma
// sessão que só serve para ler a sessão, sair, trocar de empresa e ativar o
// MFA; as demais rotas respondem 401 nao_autenticado com mfa_pendente.
// Ativado, a mesma sessão volta a valer. Enquanto a empresa exige, o MFA não
// é desativado.
func TestSessaoRestritaPorMFAObrigatorio(t *testing.T) {
	em := novaConta(t, "essencial", "")
	adm := entrar(t, em, "")

	r := adm.get("/v1/empresa")
	emp := adm.req("PATCH", "/v1/empresa", map[string]any{"mfa_obrigatorio": true}, cab{"If-Match": r.etag()}).ok(http.StatusOK)
	if emp["mfa_obrigatorio"] != true {
		t.Fatalf("PATCH /v1/empresa com mfa_obrigatorio: %v", emp)
	}

	p := adm.get("/v1/empresa").problema(http.StatusUnauthorized, "nao_autenticado")
	if p["mfa_pendente"] != true {
		t.Fatalf("rota comum com MFA pendente: %v", p)
	}
	if s := adm.get("/v1/auth/sessao").ok(http.StatusOK); s["mfa_pendente"] != true {
		t.Fatalf("a sessão diz que o MFA está pendente: %v", s)
	}

	// Novo login nessa empresa também nasce restrito.
	outro := entrar(t, em, "")
	outro.get("/v1/unidades").problema(http.StatusUnauthorized, "nao_autenticado")

	cfg := adm.post("/v1/auth/mfa/iniciar", nil).ok(http.StatusOK)
	segredo := texto(cfg, "segredo_base32")
	adm.post("/v1/auth/mfa/confirmar", map[string]any{"codigo": codigoTOTP(t, segredo, time.Now())}).ok(http.StatusNoContent)

	if s := adm.get("/v1/auth/sessao").ok(http.StatusOK); s["mfa_pendente"] != false {
		t.Fatalf("depois de ativar, a sessão deixa de ser restrita: %v", s)
	}
	adm.get("/v1/empresa").ok(http.StatusOK)

	d := adm.post("/v1/auth/mfa/desativar", map[string]any{"codigo": codigoTOTP(t, segredo, time.Now().Add(30*time.Second))}).
		problema(http.StatusConflict, "transicao_invalida")
	if d["estado_atual"] != "exigido_pela_empresa" {
		t.Fatalf("desativar com a empresa exigindo: %v", d)
	}
}
