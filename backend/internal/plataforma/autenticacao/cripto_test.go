package autenticacao

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"tenaz/internal/plataforma/erros"
)

// RFC 6238, apêndice B (SHA-1), com os 6 últimos dígitos dos códigos de 8.
func TestTOTPVetoresRFC6238(t *testing.T) {
	segredo := []byte("12345678901234567890")
	casos := map[int64]string{59: "287082", 1111111109: "081804", 1111111111: "050471", 1234567890: "005924",
		2000000000: "279037", 20000000000: "353130"}
	for unix, esperado := range casos {
		if got := codigoTOTP(segredo, uint64(unix)/30); got != esperado {
			t.Errorf("T=%d: %s, esperado %s", unix, got, esperado)
		}
	}
	agora := time.Unix(1111111109, 0)
	if !verificarTOTP(segredo, "081804", agora) {
		t.Error("código do passo atual")
	}
	if !verificarTOTP(segredo, codigoTOTP(segredo, uint64(agora.Unix())/30-1), agora) {
		t.Error("um passo atrás é aceito")
	}
	if verificarTOTP(segredo, codigoTOTP(segredo, uint64(agora.Unix())/30-2), agora) {
		t.Error("dois passos atrás é recusado")
	}
	if verificarTOTP(segredo, "08180", agora) || verificarTOTP(segredo, "", agora) {
		t.Error("tamanho errado é recusado")
	}
	if !strings.HasPrefix(uriOTPAuth("Tenaz", "ana@x.com", "ABC"), "otpauth://totp/Tenaz:ana@x.com?") {
		t.Errorf("URI otpauth: %s", uriOTPAuth("Tenaz", "ana@x.com", "ABC"))
	}
}

func TestCifraDoSegredoMFA(t *testing.T) {
	chave := make([]byte, 32)
	_, _ = rand.Read(chave)
	c, err := NovoCifrador(chave)
	if err != nil {
		t.Fatal(err)
	}
	id1, id2 := uuid.New(), uuid.New()
	env, err := c.Cifrar([]byte("segredo"), id1[:])
	if err != nil {
		t.Fatal(err)
	}
	if claro, err := c.Decifrar(env, id1[:]); err != nil || string(claro) != "segredo" {
		t.Fatalf("ida e volta: %q %v", claro, err)
	}
	if _, err := c.Decifrar(env, id2[:]); err == nil {
		t.Error("segredo copiado para outra identidade não decifra")
	}
	env[len(env)-1] ^= 1
	if _, err := c.Decifrar(env, id1[:]); err == nil {
		t.Error("envelope adulterado não decifra")
	}
	if _, err := NovoCifrador(chave[:16]); err == nil {
		t.Error("chave curta é recusada")
	}
}

func TestJWT(t *testing.T) {
	semente := make([]byte, 32)
	_, _ = rand.Read(semente)
	a, err := NovoAssinador("k1", semente, nil)
	if err != nil {
		t.Fatal(err)
	}
	agora := time.Now()
	u, e, s := uuid.New(), uuid.New(), uuid.New()
	tok, exp, err := a.Emitir(u, e, s, agora)
	if err != nil || exp.Sub(agora) != 15*time.Minute {
		t.Fatalf("emitir: %v %v", exp, err)
	}
	c, err := a.Verificar(tok, agora.Add(time.Minute))
	if err != nil || c.Sujeito != u.String() || c.Empresa != e.String() || c.Sessao != s.String() {
		t.Fatalf("verificar: %+v %v", c, err)
	}
	if _, err := a.Verificar(tok, agora.Add(16*time.Minute)); err == nil {
		t.Error("vencido é recusado")
	}
	partes := strings.Split(tok, ".")
	none := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none","typ":"JWT","kid":"k1"}`))
	if _, err := a.Verificar(none+"."+partes[1]+".", agora); err == nil {
		t.Error("alg none é recusado")
	}
	corpo, _ := base64.RawURLEncoding.DecodeString(partes[1])
	outra := strings.Replace(string(corpo), e.String(), uuid.NewString(), 1)
	if _, err := a.Verificar(partes[0]+"."+base64.RawURLEncoding.EncodeToString([]byte(outra))+"."+partes[2], agora); err == nil {
		t.Error("troca de empresa no corpo quebra a assinatura")
	}

	// Rotação: a chave nova assina; a anterior ainda verifica.
	semente2 := make([]byte, 32)
	_, _ = rand.Read(semente2)
	pub1 := ed25519.NewKeyFromSeed(semente).Public().(ed25519.PublicKey)
	b, err := NovoAssinador("k2", semente2, map[string]ed25519.PublicKey{"k1": pub1})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := b.Verificar(tok, agora); err != nil {
		t.Errorf("token da chave anterior ainda vale: %v", err)
	}
	if _, err := a.Verificar(func() string { t2, _, _ := b.Emitir(u, e, s, agora); return t2 }(), agora); err == nil {
		t.Error("kid desconhecido é recusado")
	}
}

func TestPKCEeSegredos(t *testing.T) {
	verificador := strings.Repeat("aB3-", 12)
	h := sha256.Sum256([]byte(verificador))
	desafio := base64.RawURLEncoding.EncodeToString(h[:])
	if !conferePKCE(verificador, desafio) {
		t.Error("par válido")
	}
	if conferePKCE(verificador+"x", desafio) || conferePKCE("curto", desafio) || conferePKCE(strings.Repeat("a", 50)+"!", desafio) {
		t.Error("verificador errado, curto ou com caractere inválido")
	}
	texto, hash, err := novoSegredo()
	if err != nil || len(texto) != 43 || !iguais(hash, hashSegredo(texto)) || !segredoBemFormado(texto) {
		t.Fatalf("segredo: %q %v", texto, err)
	}
	if segredoBemFormado(texto[:42]) || segredoBemFormado(texto[:42]+"*") {
		t.Error("segredo malformado")
	}
	outro, _, _ := novoSegredo()
	if tokenCSRF(texto) != tokenCSRF(texto) || tokenCSRF(texto) == tokenCSRF(outro) || tokenCSRF(texto) == texto {
		t.Error("CSRF derivado do cookie: estável, diferente por sessão e diferente do cookie")
	}
}

func hasherDeTeste(t *testing.T) *Hasher {
	t.Helper()
	h, err := NovoHasher(ParametrosArgon2{MemoriaKiB: 8 * 1024, Iteracoes: 1, Paralelo: 1, TamSal: 16, TamChave: 32}, 2)
	if err != nil {
		t.Fatal(err)
	}
	return h
}

func TestArgon2id(t *testing.T) {
	ctx := context.Background()
	h := hasherDeTeste(t)
	hash, err := h.Gerar(ctx, "cavalo-correto-bateria")
	if err != nil || !strings.HasPrefix(hash, "$argon2id$v=19$m=8192,t=1,p=1$") {
		t.Fatalf("hash: %s %v", hash, err)
	}
	if ok, rehash, err := h.Verificar(ctx, "cavalo-correto-bateria", hash); !ok || rehash || err != nil {
		t.Errorf("senha certa: %v %v %v", ok, rehash, err)
	}
	if ok, _, _ := h.Verificar(ctx, "cavalo-correto-bateriA", hash); ok {
		t.Error("senha errada")
	}
	padrao, err := NovoHasher(ParametrosPadrao(), 1)
	if err != nil {
		t.Fatal(err)
	}
	if ok, rehash, _ := padrao.Verificar(ctx, "cavalo-correto-bateria", hash); !ok || !rehash {
		t.Error("hash com parâmetros antigos confere e pede regravação")
	}
	if _, _, err := h.Verificar(ctx, "x", "$argon2id$v=19$m=99999999,t=1,p=1$YQ$YQ"); err == nil {
		t.Error("hash adulterado pedindo memória absurda é recusado")
	}
	if _, err := NovoHasher(ParametrosArgon2{MemoriaKiB: 1024, Iteracoes: 1, Paralelo: 1, TamSal: 16, TamChave: 32}, 1); err == nil {
		t.Error("parâmetros abaixo do mínimo")
	}
}

func TestPoliticaDeSenha(t *testing.T) {
	p := NovaPolitica()
	if len(p.comuns) < 10000 {
		t.Fatalf("lista embarcada com %d senhas", len(p.comuns))
	}
	recusadas := []string{"curta123", "1234567890", "qwertyuiop", "senha12345", "Password123", "aaaaaaaaaaaa",
		"abcdefghijkl", "1212121212", "ana.silva.2026", strings.Repeat("x9", 70)}
	for _, s := range recusadas {
		if err := p.Validar(s, "ana.silva@tenaz.com.br"); !errors.Is(err, erros.SenhaFraca) {
			t.Errorf("%q deveria ser recusada, veio %v", s, err)
		}
	}
	for _, s := range []string{"cavalo-correto-bateria", "Manutenção preditiva 7x", "#1princesa da oficina"} {
		if err := p.Validar(s, "ana.silva@tenaz.com.br"); err != nil {
			t.Errorf("%q deveria passar: %v", s, err)
		}
	}
}

func TestBloqueioProgressivo(t *testing.T) {
	s := &Servico{cfg: ConfigPadrao("http://x")}
	casos := map[int]time.Duration{1: 0, 4: 0, 5: 15 * time.Minute, 6: 0, 10: 30 * time.Minute, 15: time.Hour,
		50: 24 * time.Hour, 100: 24 * time.Hour}
	for falhas, esperado := range casos {
		if got := s.duracaoBloqueio(falhas); got != esperado {
			t.Errorf("%d falhas: %v, esperado %v", falhas, got, esperado)
		}
	}
}
