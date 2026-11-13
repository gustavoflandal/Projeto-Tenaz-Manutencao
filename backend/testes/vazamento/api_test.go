package vazamento

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"testing"
)

// Parte 3 da suíte: a API, rota a rota (seção 2.8 do documento 2; seção 9.3
// do documento 5; seções 4.7 a 4.9 do documento 3).
//
// Para TODA operação autenticada do openapi.yaml:
//
//   - sem credencial: 401 nao_autenticado;
//   - identificador de B no caminho, com a credencial de A (web e app): 404
//     recurso_nao_encontrado, com a MESMA resposta (menos id_rastreio e
//     instance) de um identificador que não existe, e nenhum UUID nem texto
//     de B no corpo. Em sub-recurso, também o filho de B debaixo do pai de A;
//   - GET com os registros da própria empresa: só UUIDs da própria empresa,
//     e lista vazia é erro (a verificação passaria por vacuidade);
//   - referência a registro de B no corpo e na consulta, busca pelo texto de
//     B e cursor de B: mesma resposta do inexistente, nada de B; com controle
//     positivo (a referência da própria empresa é aceita; o filtro com o
//     registro de B, pela própria B, traz linhas);
//   - empresa_id na consulta e no corpo: 400 campo_nao_permitido;
//   - prestador e locatário em rota que não declara o escopo deles: 404,
//     com a mesma resposta de um id inexistente e de uma rota que não existe
//     (não revela a estrutura); nas rotas declaradas, só o próprio alcance;
//   - as escritas válidas de A (inclusive as de rota única, com o If-Match
//     do GET da mesma rota) e as do perfil externo nas rotas dele dão certo;
//   - no fim, o retrato do banco de B e o de C (toda tabela com empresa_id,
//     as identidades e credenciais dos usuários) são os mesmos de antes:
//     nenhuma requisição de A escreveu em B nem em C.

// sondador faz as requisições e julga as respostas. O detector usa o mesmo
// sondador apontado para um servidor estragado de propósito.
type sondador struct {
	ctx  context.Context
	base string
	amb  *ambienteAPI
	cen  *cenarioAPI
	cls  *classificador
	cob  *coberturaAPI
}

type problemas []string

func (p *problemas) add(f string, a ...any) { *p = append(*p, fmt.Sprintf(f, a...)) }

// falhou registra a falha da própria varredura: ela não prova nada.
func (p *problemas) falhou(err error) bool {
	if err != nil {
		p.add("falha da varredura (o teste não prova a recusa): %v", err)
		return true
	}
	return false
}

func relatar(t *testing.T, p problemas) {
	t.Helper()
	for _, x := range p {
		t.Error(x)
	}
}

// coberturaAPI registra o que cada rota recebeu, para que nenhum caminho do
// código pule uma rota em silêncio.
type coberturaAPI struct {
	mu    sync.Mutex
	feito map[string]map[string]bool
}

func (c *coberturaAPI) marcar(op operacao, rotulo string) {
	if c == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.feito == nil {
		c.feito = map[string]map[string]bool{}
	}
	if c.feito[op.nome()] == nil {
		c.feito[op.nome()] = map[string]bool{}
	}
	c.feito[op.nome()][rotulo] = true
}

// esperado é o que toda rota do plano precisa ter recebido.
func esperado(p planoDaRota, externos []*externoAPI) []string {
	e := []string{"sem_credencial", "empresa_id"}
	if p.op.aceita("interno") {
		if p.get() {
			e = append(e, "conteudo")
		}
		if p.op.temID() {
			e = append(e, "id_outra_empresa")
		}
		for _, r := range p.refsCorpo {
			e = append(e, "ref_corpo:"+r.Campo)
		}
		for _, r := range p.refsConsulta {
			e = append(e, "ref_consulta:"+r.Campo)
		}
		if p.busca {
			e = append(e, "busca")
		}
		if p.cursor {
			e = append(e, "cursor")
		}
		if p.caso != nil && p.caso.escritaPropria {
			e = append(e, "escrita_propria")
		}
	}
	for _, x := range externos {
		e = append(e, "externo:"+x.rotulo)
		if p.caso != nil && p.caso.escritaPropria && p.op.aceita(x.escopo) && len(clientes(p.op, x.web, x.app)) > 0 {
			e = append(e, "escrita_propria:"+x.rotulo)
		}
	}
	return e
}

func (c *coberturaAPI) faltando(planos []planoDaRota, externos []*externoAPI) []string {
	c.mu.Lock()
	defer c.mu.Unlock()
	var p []string
	for _, pl := range planos {
		for _, r := range esperado(pl, externos) {
			if !c.feito[pl.op.nome()][r] {
				p = append(p, fmt.Sprintf("%s: a varredura não fez %q nesta rota", pl.op.nome(), r))
			}
		}
	}
	sort.Strings(p)
	return p
}

// ------------------------------------------------------------- caminhos

func repetir(e *empresaAPI, n int) []*empresaAPI {
	l := make([]*empresaAPI, n)
	for i := range l {
		l[i] = e
	}
	return l
}

// idDoRecurso devolve um registro da empresa para o parâmetro i do caminho:
// o do cenário, ou o primeiro do banco (filho do pai escolhido).
func (s *sondador) idDoRecurso(p planoDaRota, i int, e *empresaAPI) (string, error) {
	seg := p.op.Caminhos[i].Segmento
	if id := e.ids[seg]; id != "" {
		return id, nil
	}
	r := p.recursos[i]
	sql := `SELECT id::text FROM ` + r.Tabela + ` WHERE empresa_id = $1`
	args := []any{e.ID}
	if r.ColunaPai != "" && i > 0 {
		pai, err := s.idDoRecurso(p, i-1, e)
		if err != nil {
			return "", err
		}
		sql += ` AND ` + r.ColunaPai + ` = $2`
		args = append(args, pai)
	}
	var id string
	if err := pool.QueryRow(s.ctx, sql+` ORDER BY id LIMIT 1`, args...).Scan(&id); err != nil {
		return "", fmt.Errorf("%s: nenhum registro de %s na empresa %s para {%s}; povoe o cenário (%w)",
			p.op.nome(), r.Tabela, e.rotulo, p.op.Caminhos[i].Nome, err)
	}
	return id, nil
}

// idDaReferencia devolve um registro da empresa para um campo de referência.
func (s *sondador) idDaReferencia(recurso string, e *empresaAPI) (string, error) {
	if recurso == "empresa" {
		return e.ID, nil
	}
	if id := e.ids[recurso]; id != "" {
		return id, nil
	}
	r := recursosDeCaminho[recurso]
	var id string
	if err := pool.QueryRow(s.ctx, `SELECT id::text FROM `+r.Tabela+` WHERE empresa_id = $1 ORDER BY id LIMIT 1`, e.ID).Scan(&id); err != nil {
		return "", fmt.Errorf("nenhum registro de %s na empresa %s; povoe o cenário (%w)", r.Tabela, e.rotulo, err)
	}
	return id, nil
}

// montarCaminho troca os parâmetros do caminho: de[i] é a empresa do
// registro do parâmetro i; nil, um identificador que não existe. Devolve
// também o If-Match do último identificador, quando a rota o pede e o
// registro existe.
func (s *sondador) montarCaminho(p planoDaRota, de []*empresaAPI) (string, string, error) {
	alvo, etag := p.op.Caminho, ""
	for i, pc := range p.op.Caminhos {
		r := p.recursos[i]
		var v string
		switch {
		case r.Valor != "":
			v = r.Valor
		case de[i] == nil:
			v, etag = idInexistente(), ""
		default:
			id, err := s.idDoRecurso(p, i, de[i])
			if err != nil {
				return "", "", err
			}
			v, etag = id, ""
			if p.op.IfMatch {
				if etag, err = etagDoRegistro(s.ctx, r.Tabela, id); err != nil {
					return "", "", err
				}
			}
		}
		alvo = strings.Replace(alvo, "{"+pc.Nome+"}", v, 1)
	}
	if f := consultasPadrao[p.op.ID]; f != nil {
		alvo += "?" + f().Encode()
	}
	return alvo, etag, nil
}

// montarParaEscrever é montarCaminho com o If-Match da rota única. Na rota
// sem parâmetro de caminho que pede If-Match (PATCH /v1/empresa, PATCH
// /v1/usuarios/eu), o recurso é o da credencial: o ETag vem do GET da mesma
// rota, com a mesma credencial. Sem ele, a escrita pararia no 428 e não
// provaria nada.
func (s *sondador) montarParaEscrever(p planoDaRota, de []*empresaAPI, c *clienteAPI) (string, string, error) {
	alvo, etag, err := s.montarCaminho(p, de)
	if err != nil || !p.op.IfMatch || len(p.op.Caminhos) > 0 || c == nil {
		return alvo, etag, err
	}
	r, err := s.enviar(c, http.MethodGet, alvo, nil, "")
	if err != nil {
		return "", "", err
	}
	if r.status != http.StatusOK || r.cab.Get("ETag") == "" {
		return "", "", fmt.Errorf("%s: o GET da rota única por %s, para o If-Match, respondeu %s sem ETag", p.op.nome(), c.rotulo, r)
	}
	return alvo, r.cab.Get("ETag"), nil
}

func (s *sondador) enviar(c *clienteAPI, metodo, alvo string, corpo any, etag string) (*respostaAPI, error) {
	var extras map[string]string
	if etag != "" {
		extras = map[string]string{"If-Match": etag}
	}
	return c.fazer(s.base, metodo, alvo, corpo, extras)
}

func comConsulta(alvo string, v url.Values) string {
	if strings.Contains(alvo, "?") {
		return alvo + "&" + v.Encode()
	}
	return alvo + "?" + v.Encode()
}

// ------------------------------------------------------------ julgamento

// assinatura é a resposta sem o que muda de uma requisição para outra
// (id_rastreio, e instance, que repete o caminho pedido): respostas de mesma
// assinatura são indistinguíveis para quem chama.
func assinatura(r *respostaAPI) (string, []string) {
	var p []string
	tipo := r.cab.Get("Content-Type")
	if i := strings.IndexByte(tipo, ';'); i >= 0 {
		tipo = tipo[:i]
	}
	corpo := string(r.corpo)
	var v any
	if len(r.corpo) > 0 && json.Unmarshal(r.corpo, &v) == nil {
		if m, ok := v.(map[string]any); ok {
			delete(m, "id_rastreio")
			if inst, tem := m["instance"]; tem {
				if inst != r.caminho {
					p = append(p, fmt.Sprintf("instance %v diferente do caminho pedido %s", inst, r.caminho))
				}
				delete(m, "instance")
			}
		}
		b, _ := json.Marshal(v)
		corpo = string(b)
	}
	return fmt.Sprintf("%d %s etag=%q location=%q %s", r.status, tipo, r.cab.Get("ETag"), r.cab.Get("Location"), corpo), p
}

// mesmaResposta exige respostas indistinguíveis (nenhum oráculo).
func mesmaResposta(p *problemas, rotulo string, a, b *respostaAPI, nomeA, nomeB string) {
	sa, pa := assinatura(a)
	sb, pb := assinatura(b)
	for _, x := range append(pa, pb...) {
		p.add("%s: %s", rotulo, x)
	}
	if sa != sb {
		p.add("%s: oráculo, a resposta a %s difere da resposta a %s:\n  %s: %s\n  %s: %s", rotulo, nomeA, nomeB, nomeA, sa, nomeB, sb)
	}
}

// exigirNaoEncontrado exige 404 recurso_nao_encontrado em problem+json.
func exigirNaoEncontrado(p *problemas, rotulo string, r *respostaAPI) {
	if r.status != http.StatusNotFound || r.codigo() != "recurso_nao_encontrado" {
		p.add("%s: esperado 404 recurso_nao_encontrado, veio %s", rotulo, r)
	}
}

// semVazamento exige que a resposta não traga UUID nem texto alheio ao ator.
// Os UUIDs que a própria requisição levou (caminho, consulta e corpo) não
// contam: repetir o que o cliente mandou não revela nada.
func (s *sondador) semVazamento(p *problemas, rotulo string, a atorAPI, r *respostaAPI) {
	texto := string(r.corpo)
	var m map[string]any
	if json.Unmarshal(r.corpo, &m) == nil {
		delete(m, "id_rastreio")
		semDetalhe(m)
		b, _ := json.Marshal(m)
		texto = string(b)
	}
	texto += "\n" + r.cab.Get("Location")
	enviados := map[string]bool{}
	for _, u := range uuidsNoTexto(r.alvo + string(r.enviado)) {
		enviados[u] = true
	}
	var candidatos []string
	for _, u := range uuidsNoTexto(texto) {
		if !enviados[u] {
			candidatos = append(candidatos, u)
		}
	}
	fora, err := s.cls.alheios(s.ctx, a, candidatos)
	if p.falhou(err) {
		return
	}
	for _, u := range fora {
		p.add("%s: UUID %s alheio a %s na resposta: %s", rotulo, u, a.rotulo, r)
	}
	for _, t := range a.textos {
		if strings.Contains(texto, t) {
			p.add("%s: texto %q alheio a %s na resposta: %s", rotulo, t, a.rotulo, r)
		}
	}
}

func sucesso(r *respostaAPI) bool { return r.status >= 200 && r.status < 300 }

// semDetalhe tira o detalhe dos eventos da trilha: ele repete os filtros que
// o próprio usuário mandou numa exportação (eco do pedido). Um evento de outra
// empresa continua acusado pelo id, pelo usuario_id e pelo entidade_id.
func semDetalhe(v any) {
	switch x := v.(type) {
	case map[string]any:
		if _, ok := x["ocorrido_em"]; ok {
			delete(x, "detalhe")
		}
		for _, y := range x {
			semDetalhe(y)
		}
	case []any:
		for _, y := range x {
			semDetalhe(y)
		}
	}
}

// ---------------------------------------------------------------- atores

// textosDe: a sentinela e a razão social de cada empresa.
func textosDe(es []*empresaAPI) []string {
	var l []string
	for _, o := range es {
		l = append(l, o.sentinela, o.razao)
	}
	return l
}

func (s *sondador) atorInterno(e *empresaAPI) atorAPI {
	return atorAPI{rotulo: "administrador de " + e.rotulo, empresa: e.ID, textos: textosDe(s.cen.outras(e))}
}

func (s *sondador) atorExterno(x *externoAPI) (atorAPI, error) {
	pr, err := propriosDoExterno(s.ctx, x)
	if err != nil {
		return atorAPI{}, err
	}
	return atorAPI{rotulo: x.rotulo, empresa: x.empresa.ID, externo: true, proprios: pr,
		textos: append([]string{x.empresa.sentinela}, textosDe(s.cen.outras(x.empresa))...)}, nil
}

// clientes devolve as credenciais que a rota aceita, dentre web e app.
func clientes(op operacao, web, app *clienteAPI) []*clienteAPI {
	var l []*clienteAPI
	if op.Sessao && web != nil {
		l = append(l, web)
	}
	if op.Bearer && app != nil {
		l = append(l, app)
	}
	return l
}

// ----------------------------------------------------------- verificações

// semCredencial: 401 nao_autenticado.
func (s *sondador) semCredencial(p planoDaRota) problemas {
	var pr problemas
	s.cob.marcar(p.op, "sem_credencial")
	alvo, etag, err := s.montarCaminho(p, repetir(s.cen.A, len(p.op.Caminhos)))
	if pr.falhou(err) {
		return pr
	}
	r, err := s.enviar(nil, p.op.Metodo, alvo, p.corpo(s.cen.A), etag)
	if pr.falhou(err) {
		return pr
	}
	if r.status != http.StatusUnauthorized || r.codigo() != "nao_autenticado" {
		pr.add("%s sem credencial: esperado 401 nao_autenticado, veio %s", p.op.nome(), r)
	}
	return pr
}

// idDeOutraEmpresa: identificador de alvo no caminho, com a credencial do
// atacante. Variantes: todos os identificadores de alvo; e, em sub-recurso,
// um de cada vez de alvo com os outros do atacante (o turno de B debaixo da
// unidade de A).
func (s *sondador) idDeOutraEmpresa(p planoDaRota, c *clienteAPI, atacante, alvo *empresaAPI, a atorAPI) problemas {
	var pr problemas
	s.cob.marcar(p.op, "id_outra_empresa")
	n := len(p.op.Caminhos)
	var idx []int
	for i, r := range p.recursos {
		if r.Valor == "" {
			idx = append(idx, i)
		}
	}
	variantes := [][]*empresaAPI{repetir(alvo, n)}
	if len(idx) > 1 {
		for _, k := range idx {
			v := repetir(atacante, n)
			v[k] = alvo
			variantes = append(variantes, v)
		}
	}
	for _, v := range variantes {
		sem := make([]*empresaAPI, n)
		var partes []string
		for i := range v {
			if v[i] != alvo {
				sem[i] = v[i]
			}
			if p.recursos[i].Valor == "" {
				partes = append(partes, fmt.Sprintf("{%s} de %s", p.op.Caminhos[i].Nome, v[i].rotulo))
			}
		}
		rotulo := fmt.Sprintf("%s por %s (%s), %s", p.op.nome(), c.rotulo, c.canal(), strings.Join(partes, ", "))
		alvoX, etag, err := s.montarCaminho(p, v)
		if pr.falhou(err) {
			continue
		}
		alvoN, _, err := s.montarCaminho(p, sem)
		if pr.falhou(err) {
			continue
		}
		rX, err := s.enviar(c, p.op.Metodo, alvoX, p.corpo(atacante), etag)
		if pr.falhou(err) {
			continue
		}
		rN, err := s.enviar(c, p.op.Metodo, alvoN, p.corpo(atacante), etag)
		if pr.falhou(err) {
			continue
		}
		exigirNaoEncontrado(&pr, rotulo, rX)
		mesmaResposta(&pr, rotulo, rX, rN, "id de "+alvo.rotulo, "id inexistente")
		s.semVazamento(&pr, rotulo, a, rX)
	}
	return pr
}

// conteudoProprio: GET com os registros da própria empresa só traz o que é
// dela, e lista vazia é erro.
func (s *sondador) conteudoProprio(p planoDaRota, c *clienteAPI, e *empresaAPI, a atorAPI) problemas {
	var pr problemas
	s.cob.marcar(p.op, "conteudo")
	rotulo := fmt.Sprintf("%s por %s (%s)", p.op.nome(), c.rotulo, c.canal())
	alvo, _, err := s.montarCaminho(p, repetir(e, len(p.op.Caminhos)))
	if pr.falhou(err) {
		return pr
	}
	r, err := s.enviar(c, http.MethodGet, alvo, nil, "")
	if pr.falhou(err) {
		return pr
	}
	if r.status != http.StatusOK {
		pr.add("%s: esperado 200 com os registros da própria empresa, veio %s", rotulo, r)
		return pr
	}
	s.semVazamento(&pr, rotulo, a, r)
	if n, lista := r.linhas(); lista && n == 0 {
		pr.add("%s: lista vazia; a verificação de vazamento passaria por vacuidade: povoe %s no cenário", rotulo, e.rotulo)
	}
	if f := verificacoesDeLeitura[p.op.ID]; f != nil {
		pr = append(pr, f(s, e, r)...)
	}
	return pr
}

// escritaPropria: a escrita válida de A, que precisa dar certo; o retrato de
// B, no fim da varredura, prova que ela não alcançou B.
func (s *sondador) escritaPropria(p planoDaRota, c *clienteAPI, e *empresaAPI) problemas {
	var pr problemas
	s.cob.marcar(p.op, "escrita_propria")
	alvo, etag, err := s.montarParaEscrever(p, repetir(e, len(p.op.Caminhos)), c)
	if pr.falhou(err) {
		return pr
	}
	r, err := s.enviar(c, p.op.Metodo, alvo, p.corpo(e), etag)
	if pr.falhou(err) {
		return pr
	}
	if !sucesso(r) {
		pr.add("%s por %s: a escrita válida da própria empresa respondeu %s", p.op.nome(), c.rotulo, r)
	}
	return pr
}

// oraculoDeIdentificacao fixa a forma da exceção declarada da unicidade
// (sqlUnicidadeGlobal, em catalogo_test.go): uq_empresa_identificacao é única
// na plataforma, e o PATCH /v1/empresa com a identificação de outra empresa do
// mesmo país responde 409 conflito no campo identificacao. A exceção aceita
// esse bit (a identificação já é cliente da Tenaz), e só ele: nada da outra
// empresa na resposta e nada muda em A. D é cadastrada aqui, pelo cadastro
// público, porque A, B e C são de países diferentes e o tipo de
// identificação da empresa é o do país (00009).
func (s *sondador) oraculoDeIdentificacao(a *empresaAPI, aa atorAPI) problemas {
	var pr problemas
	const rotulo = "PATCH /v1/empresa por administrador de A com a identificação de D"
	ident, razaoD := a.pa.empresa(), "Empresa D RZ"+aleatorioHex(5)+" Ltda"
	r, err := (*clienteAPI)(nil).fazer(s.base, http.MethodPost, "/v1/empresas", map[string]any{
		"pais": a.pa.pais, "idioma_padrao": a.pa.idioma, "identificacao_tipo": a.pa.tipoEmpresa, "identificacao": ident,
		"razao_social": razaoD, "plano": "profissional",
		"responsavel": map[string]any{"nome": "Administrador de D", "email": emailAPI("admd")},
		"unidade":     map[string]any{"codigo": "MTZ", "nome": "Matriz de D", "turnos": turnoDeTeste},
	}, nil)
	if pr.falhou(err) {
		return pr
	}
	if r.status != http.StatusCreated {
		pr.add("%s: o cadastro de D (%s) respondeu %s", rotulo, a.pa.pais, r)
		return pr
	}
	const sqlA = `SELECT to_jsonb(e)::text FROM app.empresa e WHERE id = $1`
	var antes, depois string
	if pr.falhou(pool.QueryRow(s.ctx, sqlA, a.ID).Scan(&antes)) {
		return pr
	}
	g, err := a.adm.fazer(s.base, http.MethodGet, "/v1/empresa", nil, nil)
	if pr.falhou(err) {
		return pr
	}
	r, err = a.adm.fazer(s.base, http.MethodPatch, "/v1/empresa",
		map[string]any{"identificacao_tipo": a.pa.tipoEmpresa, "identificacao": ident},
		map[string]string{"If-Match": g.cab.Get("ETag")})
	if pr.falhou(err) {
		return pr
	}
	var campos []string
	if l, ok := r.json()["erros"].([]any); ok {
		for _, x := range l {
			if m, ok := x.(map[string]any); ok {
				campos = append(campos, textoDe(m, "campo"))
			}
		}
	}
	if r.status != http.StatusConflict || r.codigo() != "conflito" || len(campos) != 1 || campos[0] != "identificacao" {
		pr.add("%s: a exceção declarada é 409 conflito só no campo identificacao, veio %s", rotulo, r)
	}
	ad := aa
	ad.textos = append(append([]string(nil), aa.textos...), razaoD)
	s.semVazamento(&pr, rotulo, ad, r)
	if pr.falhou(pool.QueryRow(s.ctx, sqlA, a.ID).Scan(&depois)) {
		return pr
	}
	if antes != depois {
		pr.add("%s: a empresa A mudou com a escrita recusada:\n  antes:  %s\n  depois: %s", rotulo, antes, depois)
	}
	return pr
}

// consultaAlheia: filtro de consulta com o registro de B (ou busca pelo
// texto de B) responde como o filtro com valor que não existe, sem nada de B.
func (s *sondador) consultaAlheia(p planoDaRota, c *clienteAPI, e *empresaAPI, a atorAPI, campo, valorAlheio, valorNenhum, rotuloCob string) problemas {
	var pr problemas
	s.cob.marcar(p.op, rotuloCob)
	rotulo := fmt.Sprintf("%s?%s= por %s (%s)", p.op.nome(), campo, c.rotulo, c.canal())
	alvo, _, err := s.montarCaminho(p, repetir(e, len(p.op.Caminhos)))
	if pr.falhou(err) {
		return pr
	}
	rO, err := s.enviar(c, http.MethodGet, comConsulta(alvo, url.Values{campo: {valorAlheio}}), nil, "")
	if pr.falhou(err) {
		return pr
	}
	rN, err := s.enviar(c, http.MethodGet, comConsulta(alvo, url.Values{campo: {valorNenhum}}), nil, "")
	if pr.falhou(err) {
		return pr
	}
	if rO.status != http.StatusOK {
		pr.add("%s: esperado 200 com lista vazia, veio %s", rotulo, rO)
	}
	mesmaResposta(&pr, rotulo, rO, rN, "valor de "+s.cen.outra(e).rotulo, "valor inexistente")
	s.semVazamento(&pr, rotulo, a, rO)
	return pr
}

// controleDaConsulta: o mesmo filtro, pela própria dona do registro, traz
// linhas. Sem isso, "A não vê nada" passaria por vacuidade.
func (s *sondador) controleDaConsulta(p planoDaRota, c *clienteAPI, dona *empresaAPI, campo, valor string) problemas {
	var pr problemas
	rotulo := fmt.Sprintf("controle de %s?%s= por %s", p.op.nome(), campo, c.rotulo)
	alvo, _, err := s.montarCaminho(p, repetir(dona, len(p.op.Caminhos)))
	if pr.falhou(err) {
		return pr
	}
	r, err := s.enviar(c, http.MethodGet, comConsulta(alvo, url.Values{campo: {valor}}), nil, "")
	if pr.falhou(err) {
		return pr
	}
	if n, lista := r.linhas(); r.status != http.StatusOK || !lista || n == 0 {
		pr.add("%s: a dona do registro não encontra linha nenhuma (%s); a verificação de A passaria por vacuidade", rotulo, r)
	}
	return pr
}

// cursorDaDona pega um cursor de segunda página da dona.
func (s *sondador) cursorDaDona(p planoDaRota, c *clienteAPI, dona *empresaAPI) (string, problemas) {
	var pr problemas
	alvo, _, err := s.montarCaminho(p, repetir(dona, len(p.op.Caminhos)))
	if pr.falhou(err) {
		return "", pr
	}
	r, err := s.enviar(c, http.MethodGet, comConsulta(alvo, url.Values{"limite": {"1"}}), nil, "")
	if pr.falhou(err) {
		return "", pr
	}
	pag, _ := r.json()["paginacao"].(map[string]any)
	cur := textoDe(pag, "proximo_cursor")
	if cur == "" {
		pr.add("%s: %s não tem segunda página com limite=1 (%s); povoe o cenário para provar o cursor", p.op.nome(), dona.rotulo, r)
	}
	return cur, pr
}

// cursorAlheio: o cursor de B, usado por A, não traz nada de B.
func (s *sondador) cursorAlheio(p planoDaRota, c *clienteAPI, e *empresaAPI, a atorAPI, cursor string) problemas {
	var pr problemas
	s.cob.marcar(p.op, "cursor")
	rotulo := fmt.Sprintf("%s?cursor= de %s por %s (%s)", p.op.nome(), s.cen.outra(e).rotulo, c.rotulo, c.canal())
	alvo, _, err := s.montarCaminho(p, repetir(e, len(p.op.Caminhos)))
	if pr.falhou(err) {
		return pr
	}
	r, err := s.enviar(c, http.MethodGet, comConsulta(alvo, url.Values{"limite": {"1"}, "cursor": {cursor}}), nil, "")
	if pr.falhou(err) {
		return pr
	}
	if r.status != http.StatusOK && r.status != http.StatusBadRequest {
		pr.add("%s: esperado 200 ou 400 cursor_invalido, veio %s", rotulo, r)
	}
	s.semVazamento(&pr, rotulo, a, r)
	return pr
}

// referenciaNoCorpo: campo do corpo com o registro da outra empresa responde
// como o campo com registro inexistente, sem sucesso e sem nada alheio. O
// controle positivo, numa sessão descartável, manda o registro da própria
// empresa e exige sucesso: sem ele, um corpo inválido por outro motivo
// passaria por recusa.
func (s *sondador) referenciaNoCorpo(p planoDaRota, c *clienteAPI, e *empresaAPI, a atorAPI, controle func() (*clienteAPI, error)) problemas {
	var pr problemas
	outra := s.cen.outra(e)
	n := len(p.op.Caminhos)
	for _, ref := range p.refsCorpo {
		s.cob.marcar(p.op, "ref_corpo:"+ref.Campo)
		rotulo := fmt.Sprintf("%s com %s de %s por %s (%s)", p.op.nome(), ref.Campo, outra.rotulo, c.rotulo, c.canal())
		idO, err := s.idDaReferencia(ref.Recurso, outra)
		if pr.falhou(err) {
			continue
		}
		alvo, etag, err := s.montarParaEscrever(p, repetir(e, n), c)
		if pr.falhou(err) {
			continue
		}
		rO, err := s.enviar(c, p.op.Metodo, alvo, p.corpoComReferencia(e, ref.Campo, idO), etag)
		if pr.falhou(err) {
			continue
		}
		alvo, etag, err = s.montarParaEscrever(p, repetir(e, n), c)
		if pr.falhou(err) {
			continue
		}
		rN, err := s.enviar(c, p.op.Metodo, alvo, p.corpoComReferencia(e, ref.Campo, idInexistente()), etag)
		if pr.falhou(err) {
			continue
		}
		if sucesso(rO) {
			pr.add("%s: aceitou a referência a registro de outra empresa: %s", rotulo, rO)
		}
		mesmaResposta(&pr, rotulo, rO, rN, "registro de "+outra.rotulo, "registro inexistente")
		s.semVazamento(&pr, rotulo, a, rO)

		if controle == nil {
			continue
		}
		cc, err := controle()
		if pr.falhou(err) {
			continue
		}
		idP, err := s.idDeControle(ref.Recurso, e)
		if pr.falhou(err) {
			continue
		}
		alvo, etag, err = s.montarParaEscrever(p, repetir(e, n), cc)
		if pr.falhou(err) {
			continue
		}
		rP, err := s.enviar(cc, p.op.Metodo, alvo, p.corpoComReferencia(e, ref.Campo, idP), etag)
		if pr.falhou(err) {
			continue
		}
		if !sucesso(rP) {
			pr.add("controle de %s: com %s de %s a rota respondeu %s; o caso não prova a recusa", rotulo, ref.Campo, e.rotulo, rP)
		}
	}
	return pr
}

// idDeControle: registro da própria empresa para o controle positivo. A
// pessoa é sempre nova, porque o controle a liga a um usuário.
func (s *sondador) idDeControle(recurso string, e *empresaAPI) (string, error) {
	if recurso != "pessoas" {
		return s.idDaReferencia(recurso, e)
	}
	r, err := e.adm.fazer(s.base, http.MethodPost, "/v1/pessoas", map[string]any{"nome": "Pessoa de controle " + e.sentinela}, nil)
	if err != nil {
		return "", err
	}
	if r.status != http.StatusCreated {
		return "", fmt.Errorf("criar pessoa de controle em %s: %s", e.rotulo, r)
	}
	return idDe(r.json()), nil
}

// empresaIDNaRequisicao: empresa_id na consulta e no corpo, 400
// campo_nao_permitido (seção 4.7 do documento 3).
func (s *sondador) empresaIDNaRequisicao(p planoDaRota, c *clienteAPI, e *empresaAPI) problemas {
	var pr problemas
	s.cob.marcar(p.op, "empresa_id")
	outra := s.cen.outra(e)
	alvo, _, err := s.montarCaminho(p, repetir(e, len(p.op.Caminhos)))
	if pr.falhou(err) {
		return pr
	}
	r, err := s.enviar(c, p.op.Metodo, comConsulta(alvo, url.Values{"empresa_id": {outra.ID}}), p.corpo(e), "")
	if pr.falhou(err) {
		return pr
	}
	if r.status != http.StatusBadRequest || r.codigo() != "campo_nao_permitido" {
		pr.add("%s?empresa_id= por %s: esperado 400 campo_nao_permitido, veio %s", p.op.nome(), c.rotulo, r)
	}
	if p.caso == nil {
		return pr
	}
	corpo := p.caso.corpo(e)
	corpo["empresa_id"] = outra.ID
	r, err = s.enviar(c, p.op.Metodo, alvo, corpo, "")
	if pr.falhou(err) {
		return pr
	}
	if r.status != http.StatusBadRequest || r.codigo() != "campo_nao_permitido" {
		pr.add("%s com empresa_id no corpo por %s: esperado 400 campo_nao_permitido, veio %s", p.op.nome(), c.rotulo, r)
	}
	return pr
}

// rotaInexistente é o caminho que nenhuma rota atende: a resposta do perfil
// externo fora das rotas dele precisa ser igual à desta.
const rotaInexistente = "/v1/rota-inexistente-da-varredura"

// externoFora: perfil externo em rota que não declara o escopo dele. 404,
// igual ao de um id inexistente e ao de uma rota que não existe, com os
// registros reais da própria empresa no caminho.
func (s *sondador) externoFora(p planoDaRota, x *externoAPI, c *clienteAPI, a atorAPI) problemas {
	var pr problemas
	s.cob.marcar(p.op, "externo:"+x.rotulo)
	n := len(p.op.Caminhos)
	rotulo := fmt.Sprintf("%s por %s (%s), rota fora do escopo %s", p.op.nome(), x.rotulo, c.canal(), x.escopo)
	alvoP, etag, err := s.montarCaminho(p, repetir(x.empresa, n))
	if pr.falhou(err) {
		return pr
	}
	alvoN, _, err := s.montarCaminho(p, repetir(nil, n))
	if pr.falhou(err) {
		return pr
	}
	r1, err := s.enviar(c, p.op.Metodo, alvoP, p.corpo(x.empresa), etag)
	if pr.falhou(err) {
		return pr
	}
	r2, err := s.enviar(c, p.op.Metodo, alvoN, p.corpo(x.empresa), etag)
	if pr.falhou(err) {
		return pr
	}
	r3, err := s.enviar(c, p.op.Metodo, rotaInexistente, p.corpo(x.empresa), etag)
	if pr.falhou(err) {
		return pr
	}
	exigirNaoEncontrado(&pr, rotulo, r1)
	if n > 0 {
		mesmaResposta(&pr, rotulo, r1, r2, "registro de "+x.empresa.rotulo, "id inexistente")
	}
	mesmaResposta(&pr, rotulo, r1, r3, "rota existente", "rota inexistente")
	s.semVazamento(&pr, rotulo, a, r1)
	return pr
}

// externoDentro: perfil externo em rota que declara o escopo dele. GET sem
// identificador: sucesso e só o próprio alcance. Com identificador: o de
// outra empresa e o interno da própria empresa respondem como o inexistente.
// Referência no corpo: como a do interno, com controle. Escrita válida sem
// identificador (PATCH /v1/usuarios/eu): dá certo e responde só o próprio
// alcance. Ação sem alvo (sair, MFA, trocar a senha) age só sobre a própria
// credencial.
func (s *sondador) externoDentro(p planoDaRota, x *externoAPI, c *clienteAPI, a atorAPI) problemas {
	var pr problemas
	s.cob.marcar(p.op, "externo:"+x.rotulo)
	outra := s.cen.outra(x.empresa)
	if p.op.temID() {
		pr = append(pr, s.idDeOutraEmpresa(p, c, x.empresa, outra, a)...)
		pr = append(pr, s.idDeOutraEmpresa(p, c, x.empresa, x.empresa, a)...)
	} else if p.get() {
		rotulo := fmt.Sprintf("%s por %s (%s)", p.op.nome(), x.rotulo, c.canal())
		r, err := s.enviar(c, http.MethodGet, p.op.Caminho, nil, "")
		if pr.falhou(err) {
			return pr
		}
		if !sucesso(r) {
			pr.add("%s: rota declarada para %s respondeu %s", rotulo, x.escopo, r)
		}
		s.semVazamento(&pr, rotulo, a, r)
	}
	if len(p.refsCorpo) > 0 {
		var controle func() (*clienteAPI, error)
		if c.bearer == "" {
			controle = func() (*clienteAPI, error) { return entrarNaAPI(s.base, x.email, x.rotulo+" (controle)") }
		}
		pr = append(pr, s.referenciaNoCorpo(p, c, x.empresa, a, controle)...)
	}
	if p.caso != nil && p.caso.escritaPropria {
		pr = append(pr, s.escritaDoExterno(p, x, c, a)...)
	}
	return pr
}

// escritaDoExterno: a escrita válida do perfil externo numa rota declarada
// para ele precisa dar certo e responder só o próprio alcance; o retrato de B
// e o de C (para F1 e L1) provam que ela não alcançou outra empresa.
func (s *sondador) escritaDoExterno(p planoDaRota, x *externoAPI, c *clienteAPI, a atorAPI) problemas {
	var pr problemas
	s.cob.marcar(p.op, "escrita_propria:"+x.rotulo)
	rotulo := fmt.Sprintf("%s por %s (%s), escrita própria", p.op.nome(), x.rotulo, c.canal())
	if p.op.temID() {
		pr.add("%s: escritaPropria em rota com identificador declarada para %s; a varredura não sabe escolher o "+
			"registro do alcance dele: acrescente o caso (api_casos_test.go)", rotulo, x.escopo)
		return pr
	}
	alvo, etag, err := s.montarParaEscrever(p, repetir(x.empresa, len(p.op.Caminhos)), c)
	if pr.falhou(err) {
		return pr
	}
	r, err := s.enviar(c, p.op.Metodo, alvo, p.corpo(x.empresa), etag)
	if pr.falhou(err) {
		return pr
	}
	if !sucesso(r) {
		pr.add("%s: a escrita válida respondeu %s", rotulo, r)
	}
	s.semVazamento(&pr, rotulo, a, r)
	return pr
}

// externoEmToda: aplica externoFora ou externoDentro a toda rota, com todas
// as credenciais do externo que a rota aceita.
func (s *sondador) externoEmToda(t *testing.T, planos []planoDaRota, x *externoAPI) {
	a, err := s.atorExterno(x)
	if err != nil {
		t.Fatalf("alcance de %s: %v", x.rotulo, err)
	}
	for _, p := range planos {
		t.Run(p.op.nome(), func(t *testing.T) {
			cs := clientes(p.op, x.web, x.app)
			if len(cs) == 0 {
				// A rota aceita só uma credencial que o perfil não tem (o
				// locatario não usa o app): o escopo não alcança a rota.
				s.cob.marcar(p.op, "externo:"+x.rotulo)
				return
			}
			for _, c := range cs {
				if p.op.aceita(x.escopo) {
					relatar(t, s.externoDentro(p, x, c, a))
				} else {
					relatar(t, s.externoFora(p, x, c, a))
				}
			}
		})
	}
}

// ------------------------------------------------------------------ teste

// TestAPIIsolamentoPorRota é a parte 3 inteira, num teste só: o retrato de B
// é tirado antes do primeiro ataque e comparado depois do último, e -shuffle
// não pode pôr nada no meio.
func TestAPIIsolamentoPorRota(t *testing.T) {
	ctx := context.Background()
	amb := ambienteDaAPI(t)
	cen := amb.cen

	oapi, err := lerContratoVersionado()
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range mesmasRotas(oapi, amb.rotas) {
		t.Error(p)
	}
	ops, err := operacoesDoContrato(oapi)
	if err != nil {
		t.Fatal(err)
	}
	planos, faltas := planejar(ops)
	for _, f := range faltas {
		t.Error(f)
	}
	if len(faltas) > 0 {
		t.FailNow()
	}
	if len(planos) == 0 {
		t.Fatal("o contrato não tem rota autenticada: a varredura passaria por vacuidade")
	}
	cls, err := novoClassificador(ctx)
	if err != nil {
		t.Fatal(err)
	}
	s := &sondador{ctx: ctx, base: amb.srv.URL, amb: amb, cen: cen, cls: cls, cob: &coberturaAPI{}}
	internos := func(f func(p planoDaRota) bool) []planoDaRota {
		var l []planoDaRota
		for _, p := range planos {
			if p.op.aceita("interno") && f(p) {
				l = append(l, p)
			}
		}
		return l
	}
	gets := internos(func(p planoDaRota) bool { return p.get() })
	cursores := map[string]string{}

	// Fase 1: tudo o que os atores de B fazem, antes do retrato de B (uma
	// leitura de B pela própria B grava na trilha de B).
	t.Run("B: conteúdo próprio e controles", func(t *testing.T) {
		b, ab := cen.B, s.atorInterno(cen.B)
		for _, p := range gets {
			t.Run(p.op.nome(), func(t *testing.T) {
				for _, c := range clientes(p.op, b.adm, b.admApp) {
					relatar(t, s.conteudoProprio(p, c, b, ab))
				}
				for _, ref := range p.refsConsulta {
					id, err := s.idDaReferencia(ref.Recurso, b)
					if err != nil {
						t.Error(err)
						continue
					}
					relatar(t, s.controleDaConsulta(p, b.adm, b, ref.Campo, id))
				}
				if p.busca {
					relatar(t, s.controleDaConsulta(p, b.adm, b, "busca", b.sentinela))
				}
				if p.cursor {
					cur, pr := s.cursorDaDona(p, b.adm, b)
					relatar(t, pr)
					cursores[p.op.nome()] = cur
				}
			})
		}
	})
	t.Run("B-F1: prestador de B com o fornecedor de F1", func(t *testing.T) {
		s.externoEmToda(t, planos, cen.BF1)
	})

	amb.aut.Aguardar()
	antes, err := retratoDaEmpresa(ctx, pool, cen.B.ID)
	if err != nil {
		t.Fatal(err)
	}
	antesC, err := retratoDaEmpresa(ctx, pool, cen.C.ID)
	if err != nil {
		t.Fatal(err)
	}

	// Fase 2: só atores de A (e o anônimo).
	a, aa := cen.A, s.atorInterno(cen.A)
	t.Run("sem credencial", func(t *testing.T) {
		for _, p := range planos {
			relatar(t, s.semCredencial(p))
		}
	})
	t.Run("A: conteúdo próprio", func(t *testing.T) {
		for _, p := range gets {
			for _, c := range clientes(p.op, a.adm, a.admApp) {
				relatar(t, s.conteudoProprio(p, c, a, aa))
			}
		}
	})
	t.Run("A: id de B e de C no caminho", func(t *testing.T) {
		for _, p := range internos(func(p planoDaRota) bool { return p.op.temID() }) {
			t.Run(p.op.nome(), func(t *testing.T) {
				for _, c := range clientes(p.op, a.adm, a.admApp) {
					relatar(t, s.idDeOutraEmpresa(p, c, a, cen.B, aa))
					relatar(t, s.idDeOutraEmpresa(p, c, a, cen.C, aa))
				}
			})
		}
	})
	t.Run("A: referência de B na consulta, busca e cursor", func(t *testing.T) {
		for _, p := range gets {
			for _, c := range clientes(p.op, a.adm, a.admApp) {
				for _, ref := range p.refsConsulta {
					id, err := s.idDaReferencia(ref.Recurso, cen.B)
					if err != nil {
						t.Error(err)
						continue
					}
					relatar(t, s.consultaAlheia(p, c, a, aa, ref.Campo, id, idInexistente(), "ref_consulta:"+ref.Campo))
				}
				if p.busca {
					relatar(t, s.consultaAlheia(p, c, a, aa, "busca", cen.B.sentinela, "NADA"+aleatorioHex(4), "busca"))
				}
				if cur := cursores[p.op.nome()]; p.cursor && cur != "" {
					relatar(t, s.cursorAlheio(p, c, a, aa, cur))
				}
			}
		}
	})
	t.Run("A: referência de B no corpo", func(t *testing.T) {
		controle := func() (*clienteAPI, error) { return entrarNaAPI(s.base, a.admEmail, "administrador de A (controle)") }
		for _, p := range internos(func(p planoDaRota) bool { return len(p.refsCorpo) > 0 }) {
			t.Run(p.op.nome(), func(t *testing.T) {
				for _, c := range clientes(p.op, a.adm, a.admApp) {
					ctl := controle
					if c.bearer != "" {
						ctl = nil // o controle da web já provou o corpo
					}
					relatar(t, s.referenciaNoCorpo(p, c, a, aa, ctl))
				}
			})
		}
	})
	t.Run("A: escrita da própria empresa", func(t *testing.T) {
		for _, p := range internos(func(p planoDaRota) bool { return p.caso != nil && p.caso.escritaPropria }) {
			relatar(t, s.escritaPropria(p, a.adm, a))
		}
	})
	t.Run("A: identificação de outra empresa no PATCH /v1/empresa (exceção declarada)", func(t *testing.T) {
		relatar(t, s.oraculoDeIdentificacao(a, aa))
	})
	t.Run("A: empresa_id na requisição", func(t *testing.T) {
		for _, p := range planos {
			relatar(t, s.empresaIDNaRequisicao(p, a.adm, a))
		}
	})
	for _, x := range []*externoAPI{cen.F1, cen.L1} {
		t.Run(x.rotulo+": "+x.escopo+" de A", func(t *testing.T) { s.externoEmToda(t, planos, x) })
	}

	amb.aut.Aguardar()
	depois, err := retratoDaEmpresa(ctx, pool, cen.B.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range diferencasDoRetrato(antes, depois) {
		t.Errorf("uma requisição de A escreveu em B: %s", d)
	}
	depoisC, err := retratoDaEmpresa(ctx, pool, cen.C.ID)
	if err != nil {
		t.Fatal(err)
	}
	for _, d := range diferencasDoRetrato(antesC, depoisC) {
		t.Errorf("uma requisição de A escreveu em C: %s", d)
	}
	for _, f := range s.cob.faltando(planos, []*externoAPI{cen.BF1, cen.F1, cen.L1}) {
		t.Error(f)
	}
}
