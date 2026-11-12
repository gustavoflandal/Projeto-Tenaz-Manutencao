package obs

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"
)

// A seção 14.3 de docs/5_ARQUITETURA_SOFTWARE.md exige um filtro de redação e
// um teste que o confira. Este é esse teste.
func TestRedacaoDeCamposSensiveis(t *testing.T) {
	casos := []struct {
		campo  string
		sigilo bool
	}{
		{"senha", true},
		{"senha_atual", true},
		{"Password", true},
		{"token", true},
		{"refresh_token", true},
		{"Authorization", true},
		{"chave_api", true},
		{"X-Api-Key", true},
		{"cpf", true},
		{"cpf_condutor", true},
		{"exame_toxicologico", true},
		{"aso_resultado", true},
		{"segredo_conector", true},
		{"empresa_id", false},
		{"usuario_id", false},
		{"rota", false},
		{"status", false},
		{"duracao_ms", false},
		{"tag_ativo", false},
		{"numero_os", false},
	}

	for _, c := range casos {
		if got := Sensivel(c.campo); got != c.sigilo {
			t.Errorf("Sensivel(%q) = %v, esperado %v", c.campo, got, c.sigilo)
		}
	}
}

func TestLoggerRedigeValorSensivel(t *testing.T) {
	var buf bytes.Buffer
	l := Novo(&buf, "api", "info", true)

	l.Info("tentativa de login",
		"usuario_id", "0193f2a1-7c3e-7b10-9a2e-5d1c3f0a8b77",
		"senha", "LaranjaSinal2026",
		"chave_api", "tnz_live_abc123",
		CampoEmpresaID, "0193f2a1-0000-7b10-9a2e-000000000001",
	)

	saida := buf.String()

	if strings.Contains(saida, "LaranjaSinal2026") {
		t.Fatal("a senha apareceu no log")
	}
	if strings.Contains(saida, "tnz_live_abc123") {
		t.Fatal("a chave de API apareceu no log")
	}

	var linha map[string]any
	if err := json.Unmarshal(buf.Bytes(), &linha); err != nil {
		t.Fatalf("log não é JSON válido: %v", err)
	}

	if linha["senha"] != Redigido {
		t.Errorf("campo senha = %v, esperado %s", linha["senha"], Redigido)
	}
	if linha["usuario_id"] == nil {
		t.Error("usuario_id deveria ter sido preservado")
	}
	for _, obrigatorio := range []string{"ts", "nivel", "msg", CampoServico, CampoVersao} {
		if _, ok := linha[obrigatorio]; !ok {
			t.Errorf("campo obrigatório ausente no log: %s", obrigatorio)
		}
	}
}
