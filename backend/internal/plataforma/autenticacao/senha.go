package autenticacao

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/subtle"
	_ "embed"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/crypto/argon2"

	"tenaz/internal/plataforma/erros"
	"tenaz/internal/plataforma/textos"
)

// ParametrosArgon2 são os parâmetros do Argon2id (documento 9, seção 5.3).
//
// Padrão: m = 19 MiB, t = 2, p = 1, sal de 16 bytes e chave de 32 bytes, a
// configuração mínima recomendada pela OWASP para Argon2id. O custo é por
// login, e a memória é por login simultâneo: o Hasher limita quantos hashes
// rodam ao mesmo tempo, para que uma rajada de logins não esgote a memória.
//
// O hash guarda os parâmetros (formato PHC), então subir o custo depois não
// invalida senhas antigas: o login certo regrava o hash com os parâmetros
// novos (PrecisaRehash).
type ParametrosArgon2 struct {
	MemoriaKiB uint32
	Iteracoes  uint32
	Paralelo   uint8
	TamSal     uint32
	TamChave   uint32
}

// ParametrosPadrao é a configuração de produção.
func ParametrosPadrao() ParametrosArgon2 {
	return ParametrosArgon2{MemoriaKiB: 19 * 1024, Iteracoes: 2, Paralelo: 1, TamSal: 16, TamChave: 32}
}

// Hasher calcula e confere hashes de senha.
type Hasher struct {
	p        ParametrosArgon2
	vagas    chan struct{}
	ficticio string // hash de uma senha aleatória, para tempo constante
}

// NovoHasher monta o hasher. maxSimultaneos limita os cálculos em paralelo.
func NovoHasher(p ParametrosArgon2, maxSimultaneos int) (*Hasher, error) {
	if p.MemoriaKiB < 8*1024 || p.Iteracoes < 1 || p.Paralelo < 1 || p.TamSal < 16 || p.TamChave < 16 {
		return nil, fmt.Errorf("parâmetros do Argon2id abaixo do mínimo: %+v", p)
	}
	if maxSimultaneos < 1 {
		maxSimultaneos = 1
	}
	h := &Hasher{p: p, vagas: make(chan struct{}, maxSimultaneos)}
	aleatoria := make([]byte, 32)
	if _, err := rand.Read(aleatoria); err != nil {
		return nil, fmt.Errorf("gerar senha fictícia: %w", err)
	}
	ficticio, err := h.Gerar(context.Background(), base64.RawStdEncoding.EncodeToString(aleatoria))
	if err != nil {
		return nil, err
	}
	h.ficticio = ficticio
	return h, nil
}

func (h *Hasher) ocupar(ctx context.Context) (func(), error) {
	select {
	case h.vagas <- struct{}{}:
		return func() { <-h.vagas }, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Gerar devolve o hash no formato PHC: $argon2id$v=19$m=...,t=...,p=...$sal$chave.
func (h *Hasher) Gerar(ctx context.Context, senha string) (string, error) {
	sal := make([]byte, h.p.TamSal)
	if _, err := rand.Read(sal); err != nil {
		return "", fmt.Errorf("gerar sal: %w", err)
	}
	liberar, err := h.ocupar(ctx)
	if err != nil {
		return "", err
	}
	defer liberar()
	chave := argon2.IDKey([]byte(senha), sal, h.p.Iteracoes, h.p.MemoriaKiB, h.p.Paralelo, h.p.TamChave)
	return fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, h.p.MemoriaKiB, h.p.Iteracoes,
		h.p.Paralelo, base64.RawStdEncoding.EncodeToString(sal), base64.RawStdEncoding.EncodeToString(chave)), nil
}

var errHashInvalido = errors.New("hash de senha em formato desconhecido")

func lerHash(hash string) (ParametrosArgon2, []byte, []byte, error) {
	partes := strings.Split(hash, "$")
	if len(partes) != 6 || partes[1] != "argon2id" {
		return ParametrosArgon2{}, nil, nil, errHashInvalido
	}
	var versao int
	if _, err := fmt.Sscanf(partes[2], "v=%d", &versao); err != nil || versao != argon2.Version {
		return ParametrosArgon2{}, nil, nil, errHashInvalido
	}
	var p ParametrosArgon2
	if _, err := fmt.Sscanf(partes[3], "m=%d,t=%d,p=%d", &p.MemoriaKiB, &p.Iteracoes, &p.Paralelo); err != nil {
		return ParametrosArgon2{}, nil, nil, errHashInvalido
	}
	// Teto contra hash adulterado que pediria memória absurda.
	if p.MemoriaKiB > 1024*1024 || p.Iteracoes > 64 || p.Paralelo < 1 || p.Iteracoes < 1 {
		return ParametrosArgon2{}, nil, nil, errHashInvalido
	}
	sal, err := base64.RawStdEncoding.DecodeString(partes[4])
	if err != nil {
		return ParametrosArgon2{}, nil, nil, errHashInvalido
	}
	chave, err := base64.RawStdEncoding.DecodeString(partes[5])
	if err != nil || len(chave) < 16 {
		return ParametrosArgon2{}, nil, nil, errHashInvalido
	}
	p.TamSal = uint32(len(sal))
	p.TamChave = uint32(len(chave))
	return p, sal, chave, nil
}

// Verificar confere a senha contra o hash, em tempo constante na comparação.
// precisaRehash indica que o hash usa parâmetros diferentes dos atuais.
func (h *Hasher) Verificar(ctx context.Context, senha, hash string) (ok, precisaRehash bool, err error) {
	p, sal, esperada, err := lerHash(hash)
	if err != nil {
		return false, false, err
	}
	liberar, err := h.ocupar(ctx)
	if err != nil {
		return false, false, err
	}
	defer liberar()
	obtida := argon2.IDKey([]byte(senha), sal, p.Iteracoes, p.MemoriaKiB, p.Paralelo, p.TamChave)
	ok = subtle.ConstantTimeCompare(obtida, esperada) == 1
	precisaRehash = p.MemoriaKiB != h.p.MemoriaKiB || p.Iteracoes != h.p.Iteracoes || p.Paralelo != h.p.Paralelo ||
		p.TamChave != h.p.TamChave
	return ok, precisaRehash, nil
}

// GastarTempo faz o mesmo trabalho de uma verificação contra um hash
// fictício. É o que o login chama quando o e-mail não existe, a conta está
// bloqueada ou não tem senha, para que o tempo de resposta não revele qual
// dos casos aconteceu (documento 9, seção 5.3).
func (h *Hasher) GastarTempo(ctx context.Context, senha string) {
	_, _, _ = h.Verificar(ctx, senha, h.ficticio)
}

// Limites da política de senha (RF-ADM-006, proposta).
const (
	SenhaMinimo = 10
	SenhaMaximo = 128
)

//go:embed senhas_comuns.txt
var senhasComunsArquivo string

// Politica confere a senha nova. A lista de senhas comuns e vazadas é local:
// a senha nunca sai do servidor.
type Politica struct {
	comuns map[string]struct{}
}

// NovaPolitica carrega a lista embarcada.
func NovaPolitica() *Politica {
	comuns := make(map[string]struct{}, 16000)
	s := bufio.NewScanner(strings.NewReader(senhasComunsArquivo))
	for s.Scan() {
		l := s.Text()
		// Comentário começa por "# "; senha da lista nunca tem espaço, então
		// "#1princess" é senha, não comentário.
		if l == "" || strings.HasPrefix(l, "# ") {
			continue
		}
		comuns[l] = struct{}{}
	}
	return &Politica{comuns: comuns}
}

// Validar devolve senha_fraca com o motivo no campo "senha", ou nil.
func (p *Politica) Validar(senha, email string) error {
	n := utf8.RuneCountInString(senha)
	switch {
	case !utf8.ValidString(senha):
		return erros.SenhaFraca.ComCampoChave("senha", erros.CampoFormatoInvalido, "campo.senha_caracteres", nil)
	case n < SenhaMinimo:
		return erros.SenhaFraca.ComCampoChave("senha", erros.CampoTamanhoMinimo, "campo.senha_minimo",
			textos.Parametros{"min": SenhaMinimo})
	case n > SenhaMaximo:
		return erros.SenhaFraca.ComCampoChave("senha", erros.CampoTamanhoMaximo, "campo.tamanho_maximo",
			textos.Parametros{"max": SenhaMaximo})
	}
	normal := strings.ToLower(strings.TrimSpace(senha))
	if _, comum := p.comuns[normal]; comum {
		return erros.SenhaFraca.ComCampoChave("senha", erros.CampoValorInvalido, "campo.senha_vazada", nil)
	}
	if repetitiva(normal) {
		return erros.SenhaFraca.ComCampoChave("senha", erros.CampoValorInvalido, "campo.senha_sequencia", nil)
	}
	if local, _, ok := strings.Cut(strings.ToLower(email), "@"); ok && len(local) >= 4 && strings.Contains(normal, local) {
		return erros.SenhaFraca.ComCampoChave("senha", erros.CampoValorInvalido, "campo.senha_email", nil)
	}
	return nil
}

// repetitiva pega "aaaaaaaaaa" e sequências como "abcdefghij" ou "9876543210",
// que nenhuma lista cobre por inteiro.
func repetitiva(s string) bool {
	r := []rune(s)
	if len(r) < 2 {
		return true
	}
	igual, sobe, desce := true, true, true
	for i := 1; i < len(r); i++ {
		igual = igual && r[i] == r[0]
		sobe = sobe && r[i] == r[i-1]+1
		desce = desce && r[i] == r[i-1]-1
	}
	if igual || sobe || desce {
		return true
	}
	// Dois caracteres alternados ou repetidos: "1212121212", "abababab".
	distintos := map[rune]struct{}{}
	for _, c := range r {
		distintos[c] = struct{}{}
	}
	return len(distintos) <= 2 && !unicode.IsSpace(r[0])
}
