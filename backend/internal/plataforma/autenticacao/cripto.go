package autenticacao

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1" //nolint:gosec // HMAC-SHA1 é o algoritmo do TOTP (RFC 6238) que os aplicativos autenticadores usam
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"
)

// ------------------------------------------------------------ valores aleatórios

// TamanhoSegredo é o tamanho de todo segredo emitido: cookie de sessão, CSRF,
// refresh token, código OAuth, desafio de MFA e link de e-mail. 256 bits.
const TamanhoSegredo = 32

// novoSegredo devolve 256 bits aleatórios em base64url sem preenchimento (43
// caracteres) e o sha256 dele, que é o que o banco guarda.
func novoSegredo() (texto string, hash []byte, err error) {
	b := make([]byte, TamanhoSegredo)
	if _, err := rand.Read(b); err != nil {
		return "", nil, fmt.Errorf("gerar segredo: %w", err)
	}
	texto = base64.RawURLEncoding.EncodeToString(b)
	return texto, hashSegredo(texto), nil
}

// hashSegredo é o sha256 do texto do segredo, como recebido do cliente.
func hashSegredo(texto string) []byte {
	h := sha256.Sum256([]byte(texto))
	return h[:]
}

// segredoBemFormado recusa, sem ir ao banco, o que não pode ter sido emitido
// aqui: tamanho errado ou caractere fora do base64url.
func segredoBemFormado(texto string) bool {
	if len(texto) != base64.RawURLEncoding.EncodedLen(TamanhoSegredo) {
		return false
	}
	_, err := base64.RawURLEncoding.DecodeString(texto)
	return err == nil
}

// iguais compara hashes em tempo constante.
func iguais(a, b []byte) bool {
	return len(a) == len(b) && subtle.ConstantTimeCompare(a, b) == 1
}

// tokenCSRF deriva o token CSRF do identificador da sessão web.
//
// O banco guarda só o hash dos dois, então o CSRF não pode ser lido de volta
// para GET /v1/auth/sessao. Derivá-lo do cookie por HMAC resolve: quem tem o
// cookie (o navegador, que o envia) recebe o mesmo token a cada consulta, e
// um site de outra origem, que não lê o cookie HttpOnly, não consegue
// calculá-lo. O token não revela o cookie (HMAC é de mão única).
func tokenCSRF(cookie string) string {
	m := hmac.New(sha256.New, []byte(cookie))
	m.Write([]byte("tenaz:csrf:v1"))
	return base64.RawURLEncoding.EncodeToString(m.Sum(nil))
}

// ------------------------------------------------------------------- cifra

// Cifrador cifra o segredo TOTP com AES-256-GCM antes de gravá-lo em
// identidade.mfa_segredo_cifrado.
//
// Formato: versão (1 byte, 0x01) | nonce (12 bytes) | texto cifrado + tag.
// O dado associado é o id da identidade: um segredo copiado de uma linha para
// outra não decifra.
//
// Em desenvolvimento a chave vem de TENAZ_CHAVE_MFA (32 bytes em base64). Em
// produção, do cofre (KMS), por envelope; a troca de chave usa o byte de
// versão para conviver com segredos antigos até a recifragem.
type Cifrador struct {
	aead cipher.AEAD
}

const versaoCifra = 0x01

// NovoCifrador recebe a chave de 32 bytes.
func NovoCifrador(chave []byte) (*Cifrador, error) {
	if len(chave) != 32 {
		return nil, fmt.Errorf("a chave do MFA precisa de 32 bytes, recebeu %d", len(chave))
	}
	bloco, err := aes.NewCipher(chave)
	if err != nil {
		return nil, err
	}
	aead, err := cipher.NewGCM(bloco)
	if err != nil {
		return nil, err
	}
	return &Cifrador{aead: aead}, nil
}

// Cifrar devolve o envelope.
func (c *Cifrador) Cifrar(claro, associado []byte) ([]byte, error) {
	nonce := make([]byte, c.aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}
	saida := append([]byte{versaoCifra}, nonce...)
	return c.aead.Seal(saida, nonce, claro, associado), nil
}

// Decifrar abre o envelope.
func (c *Cifrador) Decifrar(envelope, associado []byte) ([]byte, error) {
	n := c.aead.NonceSize()
	if len(envelope) < 1+n+c.aead.Overhead() || envelope[0] != versaoCifra {
		return nil, errors.New("segredo cifrado em formato desconhecido")
	}
	return c.aead.Open(nil, envelope[1:1+n], envelope[1+n:], associado)
}

// -------------------------------------------------------------------- TOTP

// TOTP segue a RFC 6238 com os valores que todo aplicativo autenticador
// aceita: HMAC-SHA1, 6 dígitos, passo de 30 segundos. A verificação aceita um
// passo antes e um depois, para relógio de celular atrasado.
const (
	totpDigitos   = 6
	totpPasso     = 30 * time.Second
	totpJanela    = 1
	totpTamSegred = 20 // 160 bits, o tamanho do HMAC-SHA1 (RFC 4226, seção 4)
)

var base32SemPreenchimento = base32.StdEncoding.WithPadding(base32.NoPadding)

// novoSegredoTOTP devolve o segredo bruto e a forma base32 que o usuário
// digita ou lê pelo QR code.
func novoSegredoTOTP() ([]byte, string, error) {
	s := make([]byte, totpTamSegred)
	if _, err := rand.Read(s); err != nil {
		return nil, "", err
	}
	return s, base32SemPreenchimento.EncodeToString(s), nil
}

// codigoTOTP calcula o código de um passo (RFC 4226, seção 5.3).
func codigoTOTP(segredo []byte, passo uint64) string {
	var msg [8]byte
	binary.BigEndian.PutUint64(msg[:], passo)
	m := hmac.New(sha1.New, segredo)
	m.Write(msg[:])
	soma := m.Sum(nil)
	off := soma[len(soma)-1] & 0x0f
	bin := binary.BigEndian.Uint32(soma[off:off+4]) & 0x7fffffff
	return fmt.Sprintf("%0*d", totpDigitos, bin%1_000_000)
}

// verificarTOTP confere o código na janela de ±1 passo, comparando em tempo
// constante e sem sair cedo.
func verificarTOTP(segredo []byte, codigo string, agora time.Time) bool {
	_, ok := passoTOTP(segredo, codigo, agora)
	return ok
}

// passoTOTP devolve o passo em que o código confere. O passo aceito é gravado
// em identidade.mfa_ultimo_passo, e um código de passo igual ou anterior é
// recusado: o mesmo código não vale duas vezes (RFC 6238, seção 5.2).
func passoTOTP(segredo []byte, codigo string, agora time.Time) (uint64, bool) {
	codigo = strings.TrimSpace(codigo)
	if len(codigo) != totpDigitos {
		return 0, false
	}
	atual := uint64(agora.Unix()) / uint64(totpPasso/time.Second)
	var achado uint64
	ok := 0
	for d := -totpJanela; d <= totpJanela; d++ {
		p := atual + uint64(int64(d))
		igual := subtle.ConstantTimeCompare([]byte(codigoTOTP(segredo, p)), []byte(codigo))
		if igual == 1 {
			achado = p
		}
		ok |= igual
	}
	return achado, ok == 1
}

// uriOTPAuth monta o otpauth:// que o QR code leva.
func uriOTPAuth(emissor, conta, segredoBase32 string) string {
	rotulo := url.PathEscape(emissor + ":" + conta)
	q := url.Values{}
	q.Set("secret", segredoBase32)
	q.Set("issuer", emissor)
	q.Set("algorithm", "SHA1")
	q.Set("digits", fmt.Sprint(totpDigitos))
	q.Set("period", fmt.Sprint(int(totpPasso/time.Second)))
	return "otpauth://totp/" + rotulo + "?" + q.Encode()
}

// ------------------------------------------------------------------- PKCE

// conferePKCE confere o code_verifier contra o code_challenge S256 (RFC 7636,
// seção 4.6). O verificador tem de 43 a 128 caracteres não reservados.
func conferePKCE(verificador, desafio string) bool {
	if len(verificador) < 43 || len(verificador) > 128 {
		return false
	}
	for _, r := range verificador {
		if !(r >= 'A' && r <= 'Z' || r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || strings.ContainsRune("-._~", r)) {
			return false
		}
	}
	h := sha256.Sum256([]byte(verificador))
	calculado := base64.RawURLEncoding.EncodeToString(h[:])
	return subtle.ConstantTimeCompare([]byte(calculado), []byte(desafio)) == 1
}
