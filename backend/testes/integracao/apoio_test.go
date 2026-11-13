// Testes de rota da API (etapa E01), em nível HTTP e contra PostgreSQL real.
//
// Sobem a API inteira (servidor.Montar, o mesmo do processo e do contrato)
// num httptest.Server e falam com ela só por HTTP: cookie, X-CSRF-Token,
// Origin, Authorization, If-Match e Idempotency-Key, como a web e o app.
// O e-mail vai para uma caixa em memória.
package integracao_test

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha1" //nolint:gosec // TOTP da RFC 6238
	"crypto/sha256"
	"encoding/base32"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"

	"tenaz/internal/modulos/adm"
	"tenaz/internal/plataforma/api"
	"tenaz/internal/plataforma/auditoria"
	"tenaz/internal/plataforma/autenticacao"
	"tenaz/internal/plataforma/db/dbteste"
	"tenaz/internal/plataforma/email"
	"tenaz/internal/servidor"
)

const (
	origemWeb = "https://app.tenaz.teste"
	senhaBoa  = "Manutenção preditiva 7x"
)

var (
	amb   *dbteste.Ambiente
	aut   *autenticacao.Servico
	srv   *httptest.Server
	caixa = &caixaDeEntrada{}
)

type caixaDeEntrada struct {
	mu   sync.Mutex
	msgs []email.Mensagem
}

func (c *caixaDeEntrada) Enviar(_ context.Context, m email.Mensagem) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.msgs = append(c.msgs, m)
	return nil
}

var reToken = regexp.MustCompile(`#token=([A-Za-z0-9_-]{43})`)

// token devolve o segredo do último e-mail com link para o destino.
func (c *caixaDeEntrada) token(t *testing.T, para string) string {
	t.Helper()
	aut.Aguardar()
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := len(c.msgs) - 1; i >= 0; i-- {
		if strings.EqualFold(c.msgs[i].Para, para) {
			if m := reToken.FindStringSubmatch(c.msgs[i].Texto); m != nil {
				return m[1]
			}
		}
	}
	t.Fatalf("nenhum e-mail com link para %s", para)
	return ""
}

func TestMain(m *testing.M) {
	ctx := context.Background()
	var err error
	amb, err = dbteste.Preparar(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	chave := make([]byte, 32)
	_, _ = rand.Read(chave)
	cifrador, _ := autenticacao.NovoCifrador(chave)
	semente := make([]byte, 32)
	_, _ = rand.Read(semente)
	assinador, _ := autenticacao.NovoAssinador("teste", semente, nil)
	hasher, err := autenticacao.NovoHasher(autenticacao.ParametrosArgon2{MemoriaKiB: 8 * 1024, Iteracoes: 1, Paralelo: 1, TamSal: 16, TamChave: 32}, 4)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	aut, err = autenticacao.Novo(autenticacao.ConfigPadrao(origemWeb), autenticacao.Dependencias{
		Banco: amb.App, PoolAutenticacao: amb.Auth, Hasher: hasher, Politica: autenticacao.NovaPolitica(),
		Cifrador: cifrador, Assinador: assinador, Email: caixa,
	})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	capturas, err := auditoria.Novo(auditoria.Dependencias{Banco: amb.App})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	modADM, err := adm.Novo(adm.Dependencias{Banco: amb.App, Autenticacao: aut, Auditoria: capturas})
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	montado := servidor.Montar(api.Config{CookieSeguro: true, OrigensWeb: []string{origemWeb}, URLWeb: origemWeb},
		servidor.Dependencias{Autenticacao: aut, ADM: modADM})
	srv = httptest.NewServer(montado.Handler())

	codigo := m.Run()
	srv.Close()
	aut.Aguardar()
	amb.Encerrar()
	os.Exit(codigo)
}

// ------------------------------------------------------------- cliente HTTP

// cliente é um navegador (cookie, CSRF e Origin) ou o app (Bearer).
type cliente struct {
	t      *testing.T
	cookie string
	csrf   string
	bearer string
	origem string
}

func anonimo(t *testing.T) *cliente { return &cliente{t: t, origem: origemWeb} }

type resposta struct {
	t      *testing.T
	status int
	cab    http.Header
	corpo  []byte
}

// cab monta pares chave, valor.
type cab map[string]string

func (c *cliente) req(metodo, caminho string, corpo any, extras ...cab) *resposta {
	c.t.Helper()
	var leitor io.Reader
	tipo := ""
	switch v := corpo.(type) {
	case nil:
	case url.Values:
		leitor, tipo = strings.NewReader(v.Encode()), "application/x-www-form-urlencoded"
	case string:
		leitor, tipo = strings.NewReader(v), "application/json"
	default:
		b, err := json.Marshal(v)
		if err != nil {
			c.t.Fatal(err)
		}
		leitor, tipo = bytes.NewReader(b), "application/json"
	}
	r, err := http.NewRequest(metodo, srv.URL+caminho, leitor)
	if err != nil {
		c.t.Fatal(err)
	}
	if tipo != "" {
		r.Header.Set("Content-Type", tipo)
	}
	if c.cookie != "" {
		r.AddCookie(&http.Cookie{Name: api.CookieSessao, Value: c.cookie})
	}
	if c.csrf != "" && metodo != http.MethodGet {
		r.Header.Set("X-CSRF-Token", c.csrf)
	}
	if c.bearer != "" {
		r.Header.Set("Authorization", "Bearer "+c.bearer)
	}
	if c.origem != "" && metodo != http.MethodGet {
		r.Header.Set("Origin", c.origem)
	}
	for _, e := range extras {
		for k, v := range e {
			if v == "" {
				r.Header.Del(k)
			} else {
				r.Header.Set(k, v)
			}
		}
	}
	cli := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	res, err := cli.Do(r)
	if err != nil {
		c.t.Fatal(err)
	}
	defer res.Body.Close()
	b, _ := io.ReadAll(res.Body)
	return &resposta{t: c.t, status: res.StatusCode, cab: res.Header, corpo: b}
}

func (c *cliente) get(caminho string, extras ...cab) *resposta {
	return c.req(http.MethodGet, caminho, nil, extras...)
}

func (c *cliente) post(caminho string, corpo any, extras ...cab) *resposta {
	return c.req(http.MethodPost, caminho, corpo, extras...)
}

// ok confere o status de sucesso e devolve o corpo decodificado.
func (r *resposta) ok(status int) map[string]any {
	r.t.Helper()
	if r.status != status {
		r.t.Fatalf("esperado %d, veio %d: %s", status, r.status, r.corpo)
	}
	if len(r.corpo) == 0 {
		return nil
	}
	var m map[string]any
	if err := json.Unmarshal(r.corpo, &m); err != nil {
		r.t.Fatalf("corpo não é JSON: %s", r.corpo)
	}
	return m
}

// problema confere o formato RFC 9457 com o código estável.
func (r *resposta) problema(status int, codigo string) map[string]any {
	r.t.Helper()
	if r.status != status {
		r.t.Fatalf("esperado %d %s, veio %d: %s", status, codigo, r.status, r.corpo)
	}
	if ct := r.cab.Get("Content-Type"); !strings.HasPrefix(ct, "application/problem+json") {
		r.t.Fatalf("erro sem application/problem+json (%q): %s", ct, r.corpo)
	}
	var m map[string]any
	if err := json.Unmarshal(r.corpo, &m); err != nil {
		r.t.Fatalf("problema não é JSON: %s", r.corpo)
	}
	if m["codigo"] != codigo {
		r.t.Fatalf("esperado codigo %q, veio %v: %s", codigo, m["codigo"], r.corpo)
	}
	if m["type"] != "https://api.tenaz.com.br/erros/"+codigo || m["title"] == "" || m["status"] != float64(status) {
		r.t.Fatalf("problema fora do formato da seção 3.1: %s", r.corpo)
	}
	if m["id_rastreio"] == nil || m["id_rastreio"] != r.cab.Get("X-Request-Id") {
		r.t.Fatalf("id_rastreio ausente ou diferente de X-Request-Id: %s", r.corpo)
	}
	if r.cab.Get("Tenaz-Versao") == "" {
		r.t.Fatal("resposta de erro sem Tenaz-Versao")
	}
	return m
}

func (r *resposta) etag() string {
	r.t.Helper()
	e := r.cab.Get("ETag")
	if e == "" {
		r.t.Fatalf("resposta sem ETag: %s", r.corpo)
	}
	return e
}

func texto(m map[string]any, chave string) string {
	s, _ := m[chave].(string)
	return s
}

// ------------------------------------------------------------- cenários

func cnpjAleatorio() string {
	d := make([]int, 14)
	for i := 0; i < 12; i++ {
		n, _ := rand.Int(rand.Reader, big.NewInt(10))
		d[i] = int(n.Int64())
	}
	dv := func(pesos []int) int {
		s := 0
		for i, p := range pesos {
			s += d[i] * p
		}
		if r := s % 11; r >= 2 {
			return 11 - r
		}
		return 0
	}
	d[12] = dv([]int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2})
	d[13] = dv([]int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2})
	var b strings.Builder
	for _, x := range d {
		b.WriteByte(byte('0' + x))
	}
	return b.String()
}

func emailAleatorio(prefixo string) string {
	return prefixo + "." + strings.ReplaceAll(uuid.NewString()[:13], "-", "") + "@teste.tenaz.dev"
}

// novaConta cria a conta por POST /v1/empresas, aceita o convite pelo link
// do e-mail e devolve o e-mail do administrador.
func novaConta(t *testing.T, plano, emailAdm string) string {
	t.Helper()
	if emailAdm == "" {
		emailAdm = emailAleatorio("adm")
	}
	c := anonimo(t)
	c.post("/v1/empresas", map[string]any{
		"pais": "BR", "identificacao_tipo": "cnpj", "identificacao": cnpjAleatorio(), "razao_social": "Transportes Teste Ltda", "plano": plano,
		"responsavel": map[string]any{"nome": "Ana Administradora", "email": emailAdm},
		"unidade":     map[string]any{"codigo": "MTZ", "nome": "Matriz", "turnos": turnoUtil},
	}).ok(http.StatusCreated)
	tok := caixa.token(t, emailAdm)
	c.post("/v1/auth/convites/aceitar", map[string]any{"token": tok, "senha": senhaBoa}).ok(http.StatusOK)
	return emailAdm
}

// entrar faz login web (sem MFA) e devolve o navegador com cookie e CSRF.
func entrar(t *testing.T, em string, empresa string) *cliente {
	t.Helper()
	c := anonimo(t)
	corpo := map[string]any{"email": em, "senha": senhaBoa}
	if empresa != "" {
		corpo["empresa_escolhida_id"] = empresa
	}
	r := c.post("/v1/auth/login", corpo)
	s := r.ok(http.StatusOK)
	c.cookie = cookieDaResposta(t, r)
	c.csrf = texto(s, "csrf_token")
	if c.csrf == "" {
		t.Fatalf("login sem csrf_token: %s", r.corpo)
	}
	return c
}

var reCookie = regexp.MustCompile(`^__Host-tenaz_sessao=([^;]+);`)

func cookieDaResposta(t *testing.T, r *resposta) string {
	t.Helper()
	sc := r.cab.Get("Set-Cookie")
	m := reCookie.FindStringSubmatch(sc)
	if m == nil {
		t.Fatalf("Set-Cookie fora do formato: %q", sc)
	}
	for _, attr := range []string{"Path=/", "HttpOnly", "Secure", "SameSite=Lax"} {
		if !strings.Contains(sc, attr) {
			t.Fatalf("Set-Cookie sem %s: %q", attr, sc)
		}
	}
	if strings.Contains(strings.ToLower(sc), "domain=") {
		t.Fatalf("Set-Cookie com Domain: %q", sc)
	}
	return m[1]
}

// empresaAtual lê o id da empresa atual na sessão.
func empresaAtual(t *testing.T, c *cliente) string {
	t.Helper()
	s := c.get("/v1/auth/sessao").ok(http.StatusOK)
	for _, e := range s["empresas"].([]any) {
		em := e.(map[string]any)
		if em["atual"] == true {
			return texto(em, "id")
		}
	}
	t.Fatal("sessão sem empresa atual")
	return ""
}

// convidar cria o usuário por POST /v1/usuarios, aceita o convite e entra.
func convidar(t *testing.T, admin *cliente, perfis []string, extra map[string]any) (*cliente, string, string) {
	t.Helper()
	em := emailAleatorio(strings.Join(perfis, "-"))
	corpo := map[string]any{"nome": "Usuário " + strings.Join(perfis, " "), "email": em, "perfis": perfis}
	for k, v := range extra {
		corpo[k] = v
	}
	u := admin.post("/v1/usuarios", corpo).ok(http.StatusCreated)
	if u["situacao"] != "convidado" || u["convite_enviado"] != true {
		t.Fatalf("usuário novo: %v", u)
	}
	anonimo(t).post("/v1/auth/convites/aceitar", map[string]any{"token": caixa.token(t, em), "senha": senhaBoa}).ok(http.StatusOK)
	return entrar(t, em, empresaAtual(t, admin)), texto(u, "id"), em
}

// codigoTOTP calcula o código atual (RFC 6238).
func codigoTOTP(t *testing.T, segredoBase32 string, quando time.Time) string {
	t.Helper()
	s, err := base32.StdEncoding.WithPadding(base32.NoPadding).DecodeString(segredoBase32)
	if err != nil {
		t.Fatal(err)
	}
	var msg [8]byte
	binary.BigEndian.PutUint64(msg[:], uint64(quando.Unix()/30))
	h := hmac.New(sha1.New, s)
	h.Write(msg[:])
	soma := h.Sum(nil)
	off := soma[len(soma)-1] & 0x0f
	return fmt.Sprintf("%06d", (binary.BigEndian.Uint32(soma[off:off+4])&0x7fffffff)%1_000_000)
}

// pkce devolve verificador e desafio S256.
func pkce() (string, string) {
	b := make([]byte, 32)
	_, _ = rand.Read(b)
	v := base64.RawURLEncoding.EncodeToString(b)
	h := sha256.Sum256([]byte(v))
	return v, base64.RawURLEncoding.EncodeToString(h[:])
}
