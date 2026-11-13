package autenticacao

import (
	"crypto/ed25519"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Access token do app (documento 3, seção 4.3): JWT EdDSA (Ed25519) de 15
// minutos, com sub (usuário), emp (empresa), sid (sessão do aparelho), aud,
// exp. As permissões não vão no token.
//
// Implementação própria, e não uma biblioteca de JWT: o servidor só emite e só
// aceita um formato (alg EdDSA, typ JWT, kid conhecido), e a verificação
// estrita é curta. Isso elimina por construção as falhas clássicas das
// bibliotecas genéricas — "alg: none", troca de algoritmo, chave escolhida
// pelo token — e não traz dependência nova. Ed25519 é da biblioteca padrão.
//
// Chaves: em desenvolvimento, a semente de 32 bytes vem de TENAZ_CHAVE_JWT
// (base64) com o identificador em TENAZ_CHAVE_JWT_ID. Rotação: a chave nova
// passa a assinar e as anteriores ficam só para verificar, por pelo menos 15
// minutos (a vida do access token), em TENAZ_CHAVES_JWT_ANTERIORES
// (kid=chave_publica_base64, separadas por vírgula). Em produção as chaves
// vêm do cofre.

const (
	audienciaAPI     = "tenaz-api"
	emissorTenaz     = "tenaz"
	validadeAccess   = 15 * time.Minute
	toleranciaRelogi = 30 * time.Second
	tamanhoMaxJWT    = 2048
)

// ClaimsAcesso são as declarações do access token.
type ClaimsAcesso struct {
	Emissor   string `json:"iss"`
	Sujeito   string `json:"sub"` // usuario_id
	Empresa   string `json:"emp"` // empresa_id
	Sessao    string `json:"sid"` // sessao.id
	Audiencia string `json:"aud"`
	EmitidoEm int64  `json:"iat"`
	ExpiraEm  int64  `json:"exp"`
	ID        string `json:"jti"`
}

type cabecalhoJWT struct {
	Alg string `json:"alg"`
	Typ string `json:"typ"`
	Kid string `json:"kid"`
}

// Assinador emite e confere access tokens.
type Assinador struct {
	kid         string
	privada     ed25519.PrivateKey
	verificacao map[string]ed25519.PublicKey
}

// NovoAssinador monta o assinador com a semente atual e as chaves públicas
// anteriores, ainda aceitas na verificação.
func NovoAssinador(kid string, semente []byte, anteriores map[string]ed25519.PublicKey) (*Assinador, error) {
	if kid == "" || strings.ContainsAny(kid, `"\`) {
		return nil, errors.New("identificador da chave JWT inválido")
	}
	if len(semente) != ed25519.SeedSize {
		return nil, fmt.Errorf("a semente Ed25519 precisa de %d bytes, recebeu %d", ed25519.SeedSize, len(semente))
	}
	priv := ed25519.NewKeyFromSeed(semente)
	ver := map[string]ed25519.PublicKey{kid: priv.Public().(ed25519.PublicKey)}
	for k, pub := range anteriores {
		if k == kid || len(pub) != ed25519.PublicKeySize {
			return nil, fmt.Errorf("chave JWT anterior inválida: %q", k)
		}
		ver[k] = pub
	}
	return &Assinador{kid: kid, privada: priv, verificacao: ver}, nil
}

var b64 = base64.RawURLEncoding

// Emitir assina o access token.
func (a *Assinador) Emitir(usuarioID, empresaID, sessaoID uuid.UUID, agora time.Time) (string, time.Time, error) {
	exp := agora.Add(validadeAccess)
	c := ClaimsAcesso{
		Emissor: emissorTenaz, Sujeito: usuarioID.String(), Empresa: empresaID.String(), Sessao: sessaoID.String(),
		Audiencia: audienciaAPI, EmitidoEm: agora.Unix(), ExpiraEm: exp.Unix(), ID: uuid.NewString(),
	}
	cab, err := json.Marshal(cabecalhoJWT{Alg: "EdDSA", Typ: "JWT", Kid: a.kid})
	if err != nil {
		return "", time.Time{}, err
	}
	corpo, err := json.Marshal(c)
	if err != nil {
		return "", time.Time{}, err
	}
	assinado := b64.EncodeToString(cab) + "." + b64.EncodeToString(corpo)
	sig := ed25519.Sign(a.privada, []byte(assinado))
	return assinado + "." + b64.EncodeToString(sig), exp, nil
}

var errTokenInvalido = errors.New("access token inválido")

// Verificar confere assinatura, cabeçalho, emissor, audiência e validade.
// Não confere a sessão: isso é do serviço, no banco.
func (a *Assinador) Verificar(token string, agora time.Time) (ClaimsAcesso, error) {
	if len(token) > tamanhoMaxJWT {
		return ClaimsAcesso{}, errTokenInvalido
	}
	partes := strings.Split(token, ".")
	if len(partes) != 3 {
		return ClaimsAcesso{}, errTokenInvalido
	}
	cabBruto, err := b64.DecodeString(partes[0])
	if err != nil {
		return ClaimsAcesso{}, errTokenInvalido
	}
	var cab cabecalhoJWT
	if err := json.Unmarshal(cabBruto, &cab); err != nil || cab.Alg != "EdDSA" || cab.Typ != "JWT" {
		return ClaimsAcesso{}, errTokenInvalido
	}
	pub, ok := a.verificacao[cab.Kid]
	if !ok {
		return ClaimsAcesso{}, errTokenInvalido
	}
	sig, err := b64.DecodeString(partes[2])
	if err != nil || len(sig) != ed25519.SignatureSize {
		return ClaimsAcesso{}, errTokenInvalido
	}
	if !ed25519.Verify(pub, []byte(partes[0]+"."+partes[1]), sig) {
		return ClaimsAcesso{}, errTokenInvalido
	}
	corpo, err := b64.DecodeString(partes[1])
	if err != nil {
		return ClaimsAcesso{}, errTokenInvalido
	}
	var c ClaimsAcesso
	if err := json.Unmarshal(corpo, &c); err != nil {
		return ClaimsAcesso{}, errTokenInvalido
	}
	if c.Emissor != emissorTenaz || c.Audiencia != audienciaAPI {
		return ClaimsAcesso{}, errTokenInvalido
	}
	if agora.After(time.Unix(c.ExpiraEm, 0).Add(toleranciaRelogi)) ||
		time.Unix(c.EmitidoEm, 0).After(agora.Add(toleranciaRelogi)) {
		return ClaimsAcesso{}, errTokenInvalido
	}
	return c, nil
}
