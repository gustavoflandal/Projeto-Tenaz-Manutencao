package email

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"
)

func TestMontarRecusaInjecaoDeCabecalho(t *testing.T) {
	s, err := NovoSMTP(ConfigSMTP{Endereco: "127.0.0.1:1", Remetente: "nao-responda@tenaz.local"})
	if err != nil {
		t.Fatal(err)
	}
	err = s.Enviar(context.Background(), Mensagem{Para: "a@b.com", Assunto: "oi\r\nBcc: x@y.com", Texto: "t"})
	if err == nil || !strings.Contains(err.Error(), "quebra de linha") {
		t.Fatalf("assunto com CRLF: %v", err)
	}
	if _, err := NovoSMTP(ConfigSMTP{Endereco: "x:25", Remetente: "a@b.com", Usuario: "u"}); err == nil {
		t.Error("autenticação sem TLS é recusada")
	}
}

// Envio real pelo Mailpit do compose local. Opcional: só roda com
// TENAZ_TESTE_SMTP=127.0.0.1:1025 (e a API do Mailpit em
// TENAZ_TESTE_MAILPIT_API, padrão http://127.0.0.1:8025).
func TestEnvioPeloMailpit(t *testing.T) {
	endereco := os.Getenv("TENAZ_TESTE_SMTP")
	if endereco == "" {
		t.Skip("TENAZ_TESTE_SMTP não definida")
	}
	apiMailpit := os.Getenv("TENAZ_TESTE_MAILPIT_API")
	if apiMailpit == "" {
		apiMailpit = "http://127.0.0.1:8025"
	}
	s, err := NovoSMTP(ConfigSMTP{Endereco: endereco, Remetente: "nao-responda@tenaz.local", NomeRemetente: "Tenaz"})
	if err != nil {
		t.Fatal(err)
	}
	assunto := "Teste de integração " + time.Now().Format(time.RFC3339Nano)
	ctx, cancelar := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelar()
	if err := s.Enviar(ctx, Mensagem{Para: "destino@teste.tenaz.dev", Assunto: assunto, Texto: "Olá, ação concluída.\n"}); err != nil {
		t.Fatalf("enviar: %v", err)
	}
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, apiMailpit+"/api/v1/messages?limit=20", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("API do Mailpit: %v", err)
	}
	defer resp.Body.Close()
	var corpo struct {
		Messages []struct{ Subject string } `json:"messages"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&corpo); err != nil {
		t.Fatal(err)
	}
	for _, m := range corpo.Messages {
		if m.Subject == assunto {
			return
		}
	}
	t.Fatalf("mensagem não chegou ao Mailpit")
}
