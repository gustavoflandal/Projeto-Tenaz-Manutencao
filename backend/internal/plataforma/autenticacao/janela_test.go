package autenticacao

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

// RFC 6238, seção 5.2: a janela aceita um passo antes e um depois do atual,
// e nada além. O passo devolvido é o que vai para mfa_ultimo_passo.
func TestTOTPJanelaDeUmPasso(t *testing.T) {
	segredo := []byte("12345678901234567890")
	agora := time.Unix(1111111109, 0)
	atual := uint64(agora.Unix()) / 30
	casos := []struct {
		desvio int64
		aceito bool
	}{{-2, false}, {-1, true}, {0, true}, {1, true}, {2, false}, {3, false}}
	for _, c := range casos {
		p := uint64(int64(atual) + c.desvio)
		passo, ok := passoTOTP(segredo, codigoTOTP(segredo, p), agora)
		if ok != c.aceito {
			t.Errorf("passo %+d: aceito=%v, esperado %v", c.desvio, ok, c.aceito)
		}
		if ok && passo != p {
			t.Errorf("passo %+d: devolveu o passo %d, esperado %d", c.desvio, passo, p)
		}
	}
	// Espaço em volta é tolerado; letra no lugar de dígito, não.
	if _, ok := passoTOTP(segredo, " "+codigoTOTP(segredo, atual)+" ", agora); !ok {
		t.Error("código com espaço em volta")
	}
	if _, ok := passoTOTP(segredo, "08180a", agora); ok {
		t.Error("código com letra")
	}
}

// Tabela do bloqueio progressivo (documento 9, seção 5.3; RF-ADM-006): a cada
// 5 falhas seguidas, bloqueio que dobra a partir de 15 minutos, com teto de
// 24 horas. Entre os múltiplos de 5, nenhum bloqueio novo.
func TestBloqueioProgressivoTabela(t *testing.T) {
	s := &Servico{cfg: ConfigPadrao("http://x")}
	casos := map[int]time.Duration{
		0: 0, 1: 0, 2: 0, 3: 0, 4: 0,
		5: 15 * time.Minute, 6: 0, 9: 0,
		10: 30 * time.Minute, 11: 0,
		15: time.Hour, 20: 2 * time.Hour, 25: 4 * time.Hour, 30: 8 * time.Hour,
		35: 16 * time.Hour, 36: 0,
		40: 24 * time.Hour, 41: 0, 45: 24 * time.Hour, 50: 24 * time.Hour, 32000: 24 * time.Hour,
	}
	for falhas, esperado := range casos {
		if got := s.duracaoBloqueio(falhas); got != esperado {
			t.Errorf("%d falhas: %v, esperado %v", falhas, got, esperado)
		}
	}
}

// assinarClaims assina um corpo arbitrário com a chave do assinador, para
// testar o que o servidor recusa mesmo com assinatura válida.
func assinarClaims(t *testing.T, a *Assinador, cab, corpo any) string {
	t.Helper()
	cb, err := json.Marshal(cab)
	if err != nil {
		t.Fatal(err)
	}
	co, err := json.Marshal(corpo)
	if err != nil {
		t.Fatal(err)
	}
	assinado := b64.EncodeToString(cb) + "." + b64.EncodeToString(co)
	return assinado + "." + b64.EncodeToString(ed25519.Sign(a.privada, []byte(assinado)))
}

// Access token com assinatura válida e declarações fora da regra é recusado:
// sem exp, nbf no futuro, iat no futuro, emissor ou audiência errados,
// cabeçalho fora do único formato emitido.
func TestJWTRecusaDeclaracoesForaDaRegra(t *testing.T) {
	semente := make([]byte, 32)
	_, _ = rand.Read(semente)
	a, err := NovoAssinador("k1", semente, nil)
	if err != nil {
		t.Fatal(err)
	}
	agora := time.Now()
	cabOK := map[string]string{"alg": "EdDSA", "typ": "JWT", "kid": "k1"}
	base := func() map[string]any {
		return map[string]any{
			"iss": emissorTenaz, "sub": uuid.NewString(), "emp": uuid.NewString(), "sid": uuid.NewString(),
			"aud": audienciaAPI, "iat": agora.Unix(), "exp": agora.Add(10 * time.Minute).Unix(), "jti": uuid.NewString(),
		}
	}
	if _, err := a.Verificar(assinarClaims(t, a, cabOK, base()), agora); err != nil {
		t.Fatalf("controle: o token bem formado é aceito: %v", err)
	}
	if _, err := a.Verificar(assinarClaims(t, a, cabOK, base()), agora.Add(15*time.Minute)); !errors.Is(err, errTokenInvalido) {
		t.Error("vencido além da tolerância é recusado")
	}
	casos := map[string]func(c map[string]any){
		"exp ausente":       func(c map[string]any) { delete(c, "exp") },
		"exp zero":          func(c map[string]any) { c["exp"] = 0 },
		"exp em texto":      func(c map[string]any) { c["exp"] = "9999999999" },
		"nbf no futuro":     func(c map[string]any) { c["nbf"] = agora.Add(5 * time.Minute).Unix() },
		"iat no futuro":     func(c map[string]any) { c["iat"] = agora.Add(5 * time.Minute).Unix() },
		"emissor errado":    func(c map[string]any) { c["iss"] = "outro" },
		"audiência errada":  func(c map[string]any) { c["aud"] = "outra-api" },
		"audiência ausente": func(c map[string]any) { delete(c, "aud") },
	}
	for nome, mudar := range casos {
		c := base()
		mudar(c)
		if _, err := a.Verificar(assinarClaims(t, a, cabOK, c), agora); !errors.Is(err, errTokenInvalido) {
			t.Errorf("%s: esperado errTokenInvalido, veio %v", nome, err)
		}
	}
	// nbf no passado (ou agora) não impede.
	c := base()
	c["nbf"] = agora.Add(-time.Minute).Unix()
	if _, err := a.Verificar(assinarClaims(t, a, cabOK, c), agora); err != nil {
		t.Errorf("nbf no passado: %v", err)
	}
	for nome, cab := range map[string]map[string]string{
		"alg HS256":        {"alg": "HS256", "typ": "JWT", "kid": "k1"},
		"typ ausente":      {"alg": "EdDSA", "kid": "k1"},
		"kid desconhecido": {"alg": "EdDSA", "typ": "JWT", "kid": "k9"},
	} {
		if _, err := a.Verificar(assinarClaims(t, a, cab, base()), agora); !errors.Is(err, errTokenInvalido) {
			t.Errorf("%s: esperado errTokenInvalido, veio %v", nome, err)
		}
	}
	// Assinatura truncada ou de outro corpo.
	tok := assinarClaims(t, a, cabOK, base())
	if _, err := a.Verificar(tok[:len(tok)-4], agora); !errors.Is(err, errTokenInvalido) {
		t.Error("assinatura truncada")
	}
	outro := assinarClaims(t, a, cabOK, base())
	if _, err := a.Verificar(tok[:len(tok)-86]+outro[len(outro)-86:], agora); !errors.Is(err, errTokenInvalido) {
		t.Error("assinatura de outro token")
	}
}
