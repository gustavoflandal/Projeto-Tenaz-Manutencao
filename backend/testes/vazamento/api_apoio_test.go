package vazamento

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"tenaz/internal/modulos/adm"
	"tenaz/internal/plataforma/api"
	"tenaz/internal/plataforma/armazenamento"
	"tenaz/internal/plataforma/auditoria"
	"tenaz/internal/plataforma/autenticacao"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/email"
	"tenaz/internal/servidor"
)

// Ambiente da parte 3: a API inteira (servidor.Montar, a mesma do processo e
// do contrato) num httptest.Server, sobre o banco da suíte, com o tenaz_app
// e o tenaz_autenticacao de produção. O cenário (empresas A e B, prestadores,
// locatários) é montado PELA API, como um cliente faria: a suíte não fabrica
// sessão no banco, para que o teste passe pelo mesmo caminho do ataque.

const (
	origemWebAPI    = "https://app.tenaz.teste"
	senhaAPI        = "Manutenção preditiva 7x"
	clienteOAuthApp = "tenaz-app"
	retornoOAuthApp = "br.com.tenaz.app:/oauth"
)

// caixaAPI recebe os e-mails do serviço de autenticação.
type caixaAPI struct {
	mu   sync.Mutex
	msgs []email.Mensagem
}

func (c *caixaAPI) Enviar(_ context.Context, m email.Mensagem) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.msgs = append(c.msgs, m)
	return nil
}

var reTokenEmail = regexp.MustCompile(`#token=([A-Za-z0-9_-]{43})`)

func (c *caixaAPI) token(para string) (string, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	for i := len(c.msgs) - 1; i >= 0; i-- {
		if strings.EqualFold(c.msgs[i].Para, para) {
			if m := reTokenEmail.FindStringSubmatch(c.msgs[i].Texto); m != nil {
				return m[1], nil
			}
		}
	}
	return "", fmt.Errorf("nenhum e-mail com link para %s", para)
}

type ambienteAPI struct {
	srv     *httptest.Server
	handler http.Handler
	rotas   *huma.OpenAPI // rotas do servidor montado
	aut     *autenticacao.Servico
	banco   *db.DB
	poolAut *pgxpool.Pool
	caixa   *caixaAPI
	cen     *cenarioAPI
}

var (
	apiOnce sync.Once
	apiAmb  *ambienteAPI
	apiErro error
)

// ambienteDaAPI sobe a API e monta o cenário uma vez por execução.
func ambienteDaAPI(t *testing.T) *ambienteAPI {
	t.Helper()
	apiOnce.Do(func() {
		apiAmb = &ambienteAPI{caixa: &caixaAPI{}}
		apiErro = apiAmb.subir(context.Background())
	})
	if apiErro != nil {
		t.Fatalf("subir a API e montar o cenário da parte 3: %v", apiErro)
	}
	return apiAmb
}

// encerrarAPI é chamado pelo TestMain antes de fechar o banco.
func encerrarAPI() {
	a := apiAmb
	if a == nil {
		return
	}
	if a.srv != nil {
		a.srv.Close()
	}
	if a.aut != nil {
		a.aut.Aguardar()
	}
	if a.banco != nil {
		a.banco.Fechar()
	}
	if a.poolAut != nil {
		a.poolAut.Close()
	}
}

func (a *ambienteAPI) subir(ctx context.Context) error {
	var err error
	if a.banco, err = db.AbrirComOpcoes(ctx, db.Opcoes{URL: urlBanco, MaxConn: 8, Papel: papelApp}); err != nil {
		return fmt.Errorf("pool da aplicação: %w", err)
	}
	if a.poolAut, err = db.AbrirPool(ctx, db.Opcoes{URL: urlBanco, MaxConn: 4, Papel: papelAutenticacao}); err != nil {
		return fmt.Errorf("pool da autenticação: %w", err)
	}
	chave, semente := make([]byte, 32), make([]byte, 32)
	if _, err := rand.Read(chave); err != nil {
		return err
	}
	if _, err := rand.Read(semente); err != nil {
		return err
	}
	cifrador, err := autenticacao.NovoCifrador(chave)
	if err != nil {
		return err
	}
	assinador, err := autenticacao.NovoAssinador("vazamento", semente, nil)
	if err != nil {
		return err
	}
	hasher, err := autenticacao.NovoHasher(autenticacao.ParametrosArgon2{MemoriaKiB: 8 * 1024, Iteracoes: 1, Paralelo: 1, TamSal: 16, TamChave: 32}, 4)
	if err != nil {
		return err
	}
	if a.aut, err = autenticacao.Novo(autenticacao.ConfigPadrao(origemWebAPI), autenticacao.Dependencias{
		Banco: a.banco, PoolAutenticacao: a.poolAut, Hasher: hasher, Politica: autenticacao.NovaPolitica(),
		Cifrador: cifrador, Assinador: assinador, Email: a.caixa,
	}); err != nil {
		return err
	}
	aud, err := auditoria.Novo(auditoria.Dependencias{Banco: a.banco})
	if err != nil {
		return err
	}
	verificador, err := auditoria.NovoVerificador(a.banco, &armazenamentoEmMemoria{})
	if err != nil {
		return err
	}
	modADM, err := adm.Novo(adm.Dependencias{Banco: a.banco, Autenticacao: a.aut, Auditoria: aud, Verificador: verificador})
	if err != nil {
		return err
	}
	montado := servidor.Montar(api.Config{CookieSeguro: true, OrigensWeb: []string{origemWebAPI}, URLWeb: origemWebAPI},
		servidor.Dependencias{Autenticacao: a.aut, ADM: modADM})
	a.handler = montado.Handler()
	a.rotas = montado.API().OpenAPI()
	a.srv = httptest.NewServer(a.handler)
	a.cen, err = montarCenarioAPI(a)
	return err
}

// armazenamentoEmMemoria é o armazenamento do selo da parte 3: a verificação
// de integridade compara o selo com a cópia, e sem armazenamento ela responde
// 503. A chave é selos/empresas/<empresa_id>/<dia>.selo, tirada do selo que o
// RLS deixa a empresa ler.
type armazenamentoEmMemoria struct {
	mu      sync.Mutex
	objetos map[string][]byte
}

func (m *armazenamentoEmMemoria) GravarImutavel(_ context.Context, chave string, corpo []byte, _ time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.objetos == nil {
		m.objetos = map[string][]byte{}
	}
	m.objetos[chave] = append([]byte(nil), corpo...)
	return nil
}

func (m *armazenamentoEmMemoria) Existe(_ context.Context, chave string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.objetos[chave]
	return ok, nil
}

func (m *armazenamentoEmMemoria) Ler(_ context.Context, chave string) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	b, ok := m.objetos[chave]
	if !ok {
		return nil, armazenamento.ErrNaoEncontrado
	}
	return b, nil
}

// ----------------------------------------------------------- cliente HTTP

// clienteAPI é um navegador (cookie, CSRF e Origin) ou o app (Bearer).
type clienteAPI struct {
	rotulo string
	cookie string
	csrf   string
	bearer string
}

func (c *clienteAPI) canal() string {
	if c != nil && c.bearer != "" {
		return "app"
	}
	return "web"
}

type respostaAPI struct {
	metodo  string
	caminho string // sem a consulta
	alvo    string // com a consulta
	status  int
	cab     http.Header
	corpo   []byte
	enviado []byte // corpo da requisição
}

func (r *respostaAPI) String() string {
	c := string(r.corpo)
	if len(c) > 600 {
		c = c[:600] + "..."
	}
	return fmt.Sprintf("%d %s", r.status, c)
}

// fazer envia a requisição. Erro só de transporte: status é do chamador.
func (c *clienteAPI) fazer(base, metodo, alvo string, corpo any, extras map[string]string) (*respostaAPI, error) {
	var (
		leitor  io.Reader
		tipo    string
		enviado []byte
	)
	switch v := corpo.(type) {
	case nil:
	case url.Values:
		enviado = []byte(v.Encode())
		tipo = "application/x-www-form-urlencoded"
	default:
		b, err := json.Marshal(v)
		if err != nil {
			return nil, err
		}
		enviado, tipo = b, "application/json"
	}
	if enviado != nil {
		leitor = bytes.NewReader(enviado)
	}
	r, err := http.NewRequest(metodo, base+alvo, leitor)
	if err != nil {
		return nil, err
	}
	if tipo != "" {
		r.Header.Set("Content-Type", tipo)
	}
	app := c != nil && c.bearer != ""
	if c != nil && c.cookie != "" {
		r.AddCookie(&http.Cookie{Name: api.CookieSessao, Value: c.cookie})
	}
	if c != nil && c.csrf != "" && metodo != http.MethodGet {
		r.Header.Set("X-CSRF-Token", c.csrf)
	}
	if app {
		r.Header.Set("Authorization", "Bearer "+c.bearer)
	} else if metodo != http.MethodGet {
		r.Header.Set("Origin", origemWebAPI)
	}
	for k, v := range extras {
		r.Header.Set(k, v)
	}
	cli := &http.Client{CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse }}
	res, err := cli.Do(r)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	caminho := alvo
	if i := strings.IndexByte(caminho, '?'); i >= 0 {
		caminho = caminho[:i]
	}
	return &respostaAPI{metodo: metodo, caminho: caminho, alvo: alvo, status: res.StatusCode, cab: res.Header,
		corpo: b, enviado: enviado}, nil
}

// json devolve o corpo como objeto (nil se vazio ou não for objeto).
func (r *respostaAPI) json() map[string]any {
	var m map[string]any
	if len(r.corpo) == 0 || json.Unmarshal(r.corpo, &m) != nil {
		return nil
	}
	return m
}

// codigo devolve o código estável do problem+json, ou "".
func (r *respostaAPI) codigo() string {
	if !strings.HasPrefix(r.cab.Get("Content-Type"), "application/problem+json") {
		return ""
	}
	return textoDe(r.json(), "codigo")
}

func textoDe(m map[string]any, k string) string {
	s, _ := m[k].(string)
	return s
}

// dados devolve a lista de uma resposta {dados, paginacao}, e se é lista.
func (r *respostaAPI) dados() ([]any, bool) {
	m := r.json()
	if m == nil {
		return nil, false
	}
	d, ok := m["dados"].([]any)
	return d, ok
}

// linhas conta os itens de uma lista: {dados} em JSON, ou as linhas de dado
// de um CSV (sem o cabeçalho). ehLista é falso para resposta que não é lista.
func (r *respostaAPI) linhas() (n int, ehLista bool) {
	if strings.HasPrefix(r.cab.Get("Content-Type"), "text/csv") {
		for _, l := range strings.Split(strings.TrimSpace(string(r.corpo)), "\n") {
			if strings.TrimSpace(l) != "" {
				n++
			}
		}
		return max(n-1, 0), true
	}
	d, ok := r.dados()
	return len(d), ok
}

// --------------------------------------------------------------- cenário

// empresaAPI é uma empresa do cenário, criada por POST /v1/empresas.
type empresaAPI struct {
	rotulo    string
	ID        string
	sentinela string // marca todo registro interno da empresa
	razao     string // razão social, única por execução
	admEmail  string
	adm       *clienteAPI // sessão web do administrador
	admApp    *clienteAPI // access token do app do administrador
	// ids: segmento do caminho -> um registro desta empresa.
	ids map[string]string
	// pessoa devolve a identificação de uma pessoa nova, no tipo do país.
	pessoa func() map[string]any
	pa     paisAPI
}

// endereco devolve um endereço válido no país da empresa, com o município
// dado (subdivisão e código postal do país).
func (e *empresaAPI) endereco(municipio string) map[string]any {
	return map[string]any{"logradouro": "Rua da varredura", "numero": "100", "municipio": municipio,
		"subdivisao": e.pa.subdivisao, "codigo_postal": e.pa.codigoPostal}
}

// comCampos junta os campos de b em a.
func comCampos(a, b map[string]any) map[string]any {
	for k, v := range b {
		a[k] = v
	}
	return a
}

// externoAPI é um usuário prestador ou locatário.
type externoAPI struct {
	rotulo  string
	escopo  string // prestador ou locatario
	empresa *empresaAPI
	usuario string
	email   string
	web     *clienteAPI
	app     *clienteAPI // nil quando o perfil não usa o app (locatario)
}

// cenarioAPI: A, B e C, todas no Profissional, cada uma com unidade (com
// endereço), turno, feriado, unidade inativa (com endereço), pessoas (ativa,
// inativa, sem usuário), técnico com permissão avulsa, usuário inativo,
// convidado e trilha. Em A, os prestadores F1 e F2 e os locatários L1 e L2;
// em B, o prestador B-F1, com o MESMO fornecedor_id de F1 (o prestador que
// atende duas contratantes). B é o alvo de todos os ataques; C, a terceira
// empresa, é alvo do identificador no caminho, e o texto e o retrato dela são
// conferidos como os de B.
type cenarioAPI struct {
	A, B, C             *empresaAPI
	F1, F2, L1, L2, BF1 *externoAPI
	fornecedorF1        string
}

// outra é a empresa alvo dos ataques de e (B para A; A para B e C).
func (c *cenarioAPI) outra(e *empresaAPI) *empresaAPI {
	if e == c.A {
		return c.B
	}
	return c.A
}

// outras são todas as empresas do cenário menos e: o texto delas é alheio a e.
func (c *cenarioAPI) outras(e *empresaAPI) []*empresaAPI {
	var l []*empresaAPI
	for _, x := range []*empresaAPI{c.A, c.B, c.C} {
		if x != e {
			l = append(l, x)
		}
	}
	return l
}

// falhaCenario interrompe a montagem do cenário.
type falhaCenario struct{ err error }

type montador struct {
	amb  *ambienteAPI
	base string
}

func montarCenarioAPI(a *ambienteAPI) (c *cenarioAPI, err error) {
	defer func() {
		if v := recover(); v != nil {
			f, ok := v.(falhaCenario)
			if !ok {
				panic(v)
			}
			c, err = nil, f.err
		}
	}()
	m := &montador{amb: a, base: a.srv.URL}
	c = &cenarioAPI{fornecedorF1: uuid.NewString()}
	c.A = m.empresa("A")
	c.B = m.empresa("B")
	c.C = m.empresa("C")
	c.F1 = m.externo(c.A, "F1", "prestador", "fornecedor_id", c.fornecedorF1)
	c.F2 = m.externo(c.A, "F2", "prestador", "fornecedor_id", uuid.NewString())
	c.L1 = m.externo(c.A, "L1", "locatario", "locatario_id", uuid.NewString())
	c.L2 = m.externo(c.A, "L2", "locatario", "locatario_id", uuid.NewString())
	c.BF1 = m.externo(c.B, "B-F1", "prestador", "fornecedor_id", c.fornecedorF1)
	for _, e := range []*empresaAPI{c.A, c.B, c.C} {
		_, ev := m.req(e.adm, http.MethodGet, "/v1/eventos-auditoria?limite=1", nil, nil, http.StatusOK)
		d, _ := ev["dados"].([]any)
		if len(d) == 0 {
			panic(falhaCenario{fmt.Errorf("empresa %s sem evento na trilha", e.rotulo)})
		}
		e.ids["eventos-auditoria"] = textoDe(d[0].(map[string]any), "id")
	}
	a.aut.Aguardar()
	return c, nil
}

func (m *montador) req(c *clienteAPI, metodo, alvo string, corpo any, extras map[string]string, status int) (*respostaAPI, map[string]any) {
	r, err := c.fazer(m.base, metodo, alvo, corpo, extras)
	if err != nil {
		panic(falhaCenario{fmt.Errorf("cenário: %s %s: %w", metodo, alvo, err)})
	}
	if r.status != status {
		rot := "anônimo"
		if c != nil {
			rot = c.rotulo
		}
		panic(falhaCenario{fmt.Errorf("cenário: %s %s (%s): esperado %d, veio %s", metodo, alvo, rot, status, r)})
	}
	return r, r.json()
}

func idDe(m map[string]any) string { return textoDe(m, "id") }

// paisAPI é o país de uma empresa do cenário (00009, ADR-013): A é
// brasileira, B argentina e C paraguaia, para que todo ataque cruze também o
// país (identificação, moeda, idioma e endereço diferentes).
type paisAPI struct {
	pais, idioma, tipoEmpresa, tipoPessoa string
	subdivisao, codigoPostal              string // endereço da unidade
	empresa, pessoa                       func() string
}

var paisesDoCenario = map[string]paisAPI{
	"A": {pais: "BR", idioma: "pt-BR", tipoEmpresa: "cnpj", tipoPessoa: "cpf", empresa: cnpjValido, pessoa: cpfValido,
		subdivisao: "BR-SP", codigoPostal: "01310-100"},
	"B": {pais: "AR", idioma: "es-419", tipoEmpresa: "cuit", tipoPessoa: "dni", empresa: cuitValido,
		pessoa: func() string { return digitosSemZeroNaFrente(8) },
		// CPA: a letra é a da subdivisão (AR-C, Cidade de Buenos Aires).
		subdivisao: "AR-C", codigoPostal: "C1425ABC"},
	"C": {pais: "PY", idioma: "es-419", tipoEmpresa: "ruc", tipoPessoa: "ci",
		empresa: func() string { return digitosSemZeroNaFrente(8) }, pessoa: func() string { return digitosSemZeroNaFrente(7) },
		subdivisao: "PY-ASU", codigoPostal: "001001"},
}

// digitosSemZeroNaFrente: n dígitos aleatórios, o primeiro de 1 a 9 (DNI, CI
// e RUC só pela forma, até a pendência P13).
func digitosSemZeroNaFrente(n int) string {
	return juntarDigitos(append([]int{1 + digitosAleatorios(1)[0]%9}, digitosAleatorios(n-1)...))
}

func (m *montador) empresa(rotulo string) *empresaAPI {
	e := &empresaAPI{rotulo: rotulo, sentinela: "VAZ" + rotulo + aleatorioHex(5), ids: map[string]string{}}
	e.razao = "Empresa " + rotulo + " RZ" + aleatorioHex(5) + " Ltda"
	e.admEmail = emailAPI("adm" + strings.ToLower(rotulo))
	pa := paisesDoCenario[rotulo]
	e.pa = pa
	e.pessoa = func() map[string]any {
		return map[string]any{"identificacao_tipo": pa.tipoPessoa, "identificacao": pa.pessoa()}
	}
	s := e.sentinela
	m.req(nil, http.MethodPost, "/v1/empresas", map[string]any{
		"pais": pa.pais, "idioma_padrao": pa.idioma, "identificacao_tipo": pa.tipoEmpresa, "identificacao": pa.empresa(),
		"razao_social": e.razao, "plano": "profissional",
		"responsavel": map[string]any{"nome": "Administrador " + s, "email": e.admEmail},
		"unidade":     map[string]any{"codigo": "MTZ", "nome": "Matriz " + s, "turnos": turnoDeTeste},
	}, nil, http.StatusCreated)
	m.aceitarConvite(e.admEmail)
	e.adm = m.entrar(e.admEmail, "administrador de "+rotulo)
	_, emp := m.req(e.adm, http.MethodGet, "/v1/empresa", nil, nil, http.StatusOK)
	e.ID = idDe(emp)
	e.admApp = m.app(e.adm, e.ID, "administrador de "+rotulo+" (app)")
	if e.admApp == nil {
		panic(falhaCenario{fmt.Errorf("o administrador de %s não recebeu access token do app", rotulo)})
	}

	// O município leva a sentinela: um endereço de outra empresa que vaze sem
	// UUID nem nome é pego pela varredura de texto.
	_, u := m.req(e.adm, http.MethodPost, "/v1/unidades", map[string]any{"codigo": "FIL", "nome": "Filial " + s, "turnos": turnoDeTeste,
		"endereco": e.endereco("Cidade " + s)}, nil, http.StatusCreated)
	fil := idDe(u)
	e.ids["unidades"] = fil
	e.ids["turnos"] = idDe(u["turnos"].([]any)[0].(map[string]any))
	_, f := m.req(e.adm, http.MethodPost, "/v1/unidades/"+fil+"/feriados", map[string]any{"data": "2027-04-21", "descricao": "Feriado " + s}, nil, http.StatusCreated)
	e.ids["feriados"] = idDe(f)
	// Captura da trilha: só B configura a tela, com uma combinação que o
	// padrão não tem; A, que não configurou, nunca pode vê-la.
	if rotulo == "B" {
		m.req(e.adm, http.MethodPut, "/v1/capturas-auditoria/"+telaDeCaptura, capturaDeB, nil, http.StatusOK)
	}
	r, u2 := m.req(e.adm, http.MethodPost, "/v1/unidades", map[string]any{"codigo": "INA", "nome": "Inativa " + s, "turnos": turnoDeTeste,
		"endereco": e.endereco("Cidade inativa " + s)}, nil, http.StatusCreated)
	m.req(e.adm, http.MethodDelete, "/v1/unidades/"+idDe(u2), nil, map[string]string{"If-Match": r.cab.Get("ETag")}, http.StatusNoContent)

	_, p := m.req(e.adm, http.MethodPost, "/v1/pessoas", comCampos(e.pessoa(), map[string]any{"nome": "Pessoa " + s, "unidade_id": fil}), nil, http.StatusCreated)
	e.ids["pessoas"] = idDe(p)
	r, p2 := m.req(e.adm, http.MethodPost, "/v1/pessoas", map[string]any{"nome": "Pessoa inativa " + s}, nil, http.StatusCreated)
	m.req(e.adm, http.MethodDelete, "/v1/pessoas/"+idDe(p2), nil, map[string]string{"If-Match": r.cab.Get("ETag")}, http.StatusNoContent)
	m.req(e.adm, http.MethodPost, "/v1/pessoas", comCampos(e.pessoa(), map[string]any{"nome": "Pessoa livre " + s, "unidade_id": fil}), nil, http.StatusCreated)

	tec, _, _ := m.convidar(e, []string{"tecnico"}, nil, "Técnico "+s, true, false)
	m.req(e.adm, http.MethodPost, "/v1/usuarios/"+tec+"/permissoes", map[string]any{"permissao": "os.encerrar"}, nil, http.StatusCreated)
	e.ids["usuarios"] = tec
	ina, _, _ := m.convidar(e, []string{"tecnico"}, nil, "Técnico inativo "+s, true, false)
	r, _ = m.req(e.adm, http.MethodGet, "/v1/usuarios/"+ina, nil, nil, http.StatusOK)
	m.req(e.adm, http.MethodPost, "/v1/usuarios/"+ina+"/desativar", nil, map[string]string{"If-Match": r.cab.Get("ETag")}, http.StatusOK)
	m.convidar(e, []string{"solicitante"}, nil, "Convidado "+s, false, false)
	return e
}

func (m *montador) externo(e *empresaAPI, rotulo, escopo, campo, vinculo string) *externoAPI {
	nome := map[string]string{"prestador": "Prestador ", "locatario": "Locatário "}[escopo] + rotulo
	id, em, web := m.convidar(e, []string{escopo}, map[string]any{campo: vinculo}, nome, true, true)
	x := &externoAPI{rotulo: rotulo, escopo: escopo, empresa: e, usuario: id, email: em, web: web}
	x.app = m.app(web, e.ID, nome+" (app)")
	// O prestador usa o app; o locatario, só o portal web (seção 4.9 do
	// documento 3). Sem o app do prestador, a varredura pelo Bearer não rodaria.
	if (escopo == "prestador") != (x.app != nil) {
		panic(falhaCenario{fmt.Errorf("%s (%s): app = %v, contrário da seção 4.9 do documento 3", rotulo, escopo, x.app != nil)})
	}
	return x
}

// convidar cria o usuário por POST /v1/usuarios; aceita o convite e entra,
// se pedido.
func (m *montador) convidar(e *empresaAPI, perfis []string, extra map[string]any, nome string, aceitar, entrar bool) (string, string, *clienteAPI) {
	em := emailAPI(strings.Join(perfis, "-"))
	corpo := map[string]any{"nome": nome, "email": em, "perfis": perfis}
	for k, v := range extra {
		corpo[k] = v
	}
	_, u := m.req(e.adm, http.MethodPost, "/v1/usuarios", corpo, nil, http.StatusCreated)
	if !aceitar {
		return idDe(u), em, nil
	}
	m.aceitarConvite(em)
	if !entrar {
		return idDe(u), em, nil
	}
	return idDe(u), em, m.entrar(em, nome)
}

func (m *montador) aceitarConvite(em string) {
	m.amb.aut.Aguardar()
	tok, err := m.amb.caixa.token(em)
	if err != nil {
		panic(falhaCenario{err})
	}
	m.req(nil, http.MethodPost, "/v1/auth/convites/aceitar", map[string]any{"token": tok, "senha": senhaAPI}, nil, http.StatusOK)
}

var reCookieSessao = regexp.MustCompile(`__Host-tenaz_sessao=([^;]+);`)

// entrar faz login web e devolve o navegador.
func (m *montador) entrar(em, rotulo string) *clienteAPI {
	c, err := entrarNaAPI(m.base, em, rotulo)
	if err != nil {
		panic(falhaCenario{err})
	}
	return c
}

// entrarNaAPI é o login web, fora do cenário: as provas de controle abrem
// uma sessão descartável para não gastar a do cenário.
func entrarNaAPI(base, em, rotulo string) (*clienteAPI, error) {
	r, err := (*clienteAPI)(nil).fazer(base, http.MethodPost, "/v1/auth/login", map[string]any{"email": em, "senha": senhaAPI}, nil)
	if err != nil {
		return nil, err
	}
	if r.status != http.StatusOK {
		return nil, fmt.Errorf("login de %s: %s", rotulo, r)
	}
	cm := reCookieSessao.FindStringSubmatch(r.cab.Get("Set-Cookie"))
	csrf := textoDe(r.json(), "csrf_token")
	if cm == nil || csrf == "" {
		return nil, fmt.Errorf("login de %s sem cookie ou csrf_token: %q %s", rotulo, r.cab.Get("Set-Cookie"), r)
	}
	return &clienteAPI{rotulo: rotulo, cookie: cm[1], csrf: csrf}, nil
}

// app obtém o access token do app pelo OAuth com PKCE, a partir da sessão
// web; nil quando o perfil não usa o app (access_denied).
func (m *montador) app(web *clienteAPI, empresa, rotulo string) *clienteAPI {
	verificador, desafio := desafioPKCE()
	_, a := m.req(web, http.MethodPost, "/oauth/authorize", map[string]any{
		"response_type": "code", "client_id": clienteOAuthApp, "redirect_uri": retornoOAuthApp,
		"code_challenge": desafio, "code_challenge_method": "S256", "state": "vazamento", "empresa_escolhida_id": empresa,
	}, nil, http.StatusOK)
	retorno, err := url.Parse(textoDe(a, "uri_redirecionamento"))
	if err != nil {
		panic(falhaCenario{err})
	}
	if retorno.Query().Get("error") == "access_denied" {
		return nil
	}
	_, tk := m.req(&clienteAPI{rotulo: rotulo}, http.MethodPost, "/oauth/token", url.Values{
		"grant_type": {"authorization_code"}, "code": {retorno.Query().Get("code")}, "redirect_uri": {retornoOAuthApp},
		"client_id": {clienteOAuthApp}, "code_verifier": {verificador},
	}, map[string]string{"Origin": ""}, http.StatusOK)
	return &clienteAPI{rotulo: rotulo, bearer: textoDe(tk, "access_token")}
}

// ------------------------------------------------------------ geradores

func aleatorioHex(n int) string {
	b := make([]byte, n)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	return hex.EncodeToString(b)
}

func emailAPI(prefixo string) string {
	return prefixo + "." + aleatorioHex(6) + "@vazamento.tenaz.dev"
}

func desafioPKCE() (string, string) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	v := base64.RawURLEncoding.EncodeToString(b)
	h := sha256.Sum256([]byte(v))
	return v, base64.RawURLEncoding.EncodeToString(h[:])
}

func digitosAleatorios(n int) []int {
	d := make([]int, n)
	for i := range d {
		x, err := rand.Int(rand.Reader, big.NewInt(10))
		if err != nil {
			panic(err)
		}
		d[i] = int(x.Int64())
	}
	return d
}

func digitoVerificador(d []int, pesos []int) int {
	s := 0
	for i, p := range pesos {
		s += d[i] * p
	}
	if r := s % 11; r >= 2 {
		return 11 - r
	}
	return 0
}

func juntarDigitos(d []int) string {
	var b strings.Builder
	for _, x := range d {
		b.WriteByte(byte('0' + x))
	}
	return b.String()
}

func cnpjValido() string {
	d := append(digitosAleatorios(12), 0, 0)
	d[12] = digitoVerificador(d, []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2})
	d[13] = digitoVerificador(d, []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2})
	return juntarDigitos(d)
}

// cuitValido: CUIT argentino com dígito verificador (módulo 11, pesos
// 5,4,3,2,7,6,5,4,3,2); resto que daria 10 é descartado.
func cuitValido() string {
	for {
		d := append([]int{2, 0}, digitosAleatorios(8)...)
		s := 0
		for i, p := range []int{5, 4, 3, 2, 7, 6, 5, 4, 3, 2} {
			s += d[i] * p
		}
		v := 11 - s%11
		if v == 10 {
			continue
		}
		if v == 11 {
			v = 0
		}
		return juntarDigitos(append(d, v))
	}
}

func cpfValido() string {
	for {
		d := append(digitosAleatorios(9), 0, 0)
		d[9] = digitoVerificador(d, []int{10, 9, 8, 7, 6, 5, 4, 3, 2})
		d[10] = digitoVerificador(d, []int{11, 10, 9, 8, 7, 6, 5, 4, 3, 2})
		if s := juntarDigitos(d); strings.Count(s, s[:1]) != len(s) {
			return s
		}
	}
}
