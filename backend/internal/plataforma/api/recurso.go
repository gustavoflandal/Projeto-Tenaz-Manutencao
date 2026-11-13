package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"tenaz/internal/plataforma/erros"
	"tenaz/internal/plataforma/paginacao"
)

// ------------------------------------------------------------ ETag e If-Match

// ETag devolve o ETag do recurso: a versão entre aspas (seção 2.11).
func ETag(versao int64) *string {
	e := `"` + strconv.FormatInt(versao, 10) + `"`
	return &e
}

// NaoModificado informa se If-None-Match confere com a versão atual (304).
func NaoModificado(ifNoneMatch string, versao int64) bool {
	if ifNoneMatch == "" {
		return false
	}
	for _, parte := range strings.Split(ifNoneMatch, ",") {
		if v, ok := lerETag(parte); ok && v == versao {
			return true
		}
	}
	return false
}

func lerETag(s string) (int64, bool) {
	s = strings.TrimPrefix(strings.TrimSpace(s), "W/")
	s = strings.Trim(s, `"`)
	v, err := strconv.ParseInt(s, 10, 64)
	return v, err == nil && v >= 0
}

// VersaoEsperada lê o If-Match obrigatório de PATCH, DELETE e ações de
// estado: ausente, 428 pre_condicao_obrigatoria; malformado, 400.
func VersaoEsperada(ifMatch string) (int64, error) {
	if strings.TrimSpace(ifMatch) == "" {
		return 0, erros.PreCondicaoObrigatoria.Com("Envie o cabeçalho If-Match com o ETag lido do recurso.")
	}
	v, ok := lerETag(ifMatch)
	if !ok {
		return 0, erros.RequisicaoInvalida.Com(`If-Match deve ser o ETag lido do recurso, como "12".`)
	}
	return v, nil
}

// Desatualizado monta o 412 com o recurso atual (seção 2.11).
func Desatualizado(versaoAtual int64, recursoAtual any) error {
	return erros.VersaoDesatualizada.
		Com("O registro mudou desde a leitura (versão atual %d); leia de novo e refaça a alteração.", versaoAtual).
		ComDado("versao_atual", versaoAtual).ComDado("recurso_atual", recursoAtual)
}

// ConferirVersao compara a versão esperada com a atual.
func ConferirVersao(esperada, atual int64, recursoAtual func() any) error {
	if esperada != atual {
		return Desatualizado(atual, recursoAtual())
	}
	return nil
}

// EhDesatualizado informa se o erro é 412 versao_desatualizada.
func EhDesatualizado(err error) bool { return errors.Is(err, erros.VersaoDesatualizada) }

// ----------------------------------------------------------------- listas

// Paginacao é o bloco de paginação das listas (seção 2.6).
type Paginacao struct {
	ProximoCursor *string `json:"proximo_cursor" doc:"Valor para o parâmetro cursor da próxima página; nulo quando não há mais"`
	TemMais       bool    `json:"tem_mais" doc:"Há mais itens depois desta página"`
	Limite        int32   `json:"limite" doc:"Itens por página" example:"50"`
}

// Lista é o envelope de toda lista: dados e paginacao.
type Lista[T any] struct {
	Dados     []T       `json:"dados"`
	Paginacao Paginacao `json:"paginacao"`
}

// DePagina converte a página do caso de uso, aplicando a conversão de cada
// item.
func DePagina[A, B any](p paginacao.Pagina[A], conv func(A) B) Lista[B] {
	l := Lista[B]{Dados: make([]B, 0, len(p.Dados)), Paginacao: Paginacao{TemMais: p.TemMais, Limite: p.Limite}}
	for _, x := range p.Dados {
		l.Dados = append(l.Dados, conv(x))
	}
	if p.ProximoCursor != "" {
		c := p.ProximoCursor
		l.Paginacao.ProximoCursor = &c
	}
	return l
}

// ListaCompleta embrulha uma lista curta, devolvida inteira, no mesmo
// envelope (tem_mais falso).
func ListaCompleta[A, B any](itens []A, conv func(A) B) Lista[B] {
	l := Lista[B]{Dados: make([]B, 0, len(itens))}
	for _, x := range itens {
		l.Dados = append(l.Dados, conv(x))
	}
	l.Paginacao.Limite = int32(len(l.Dados))
	return l
}

// ----------------------------------------------------- JSON Merge Patch

// MesclarPatch aplica o corpo do PATCH, como JSON Merge Patch (RFC 7396),
// sobre os campos editáveis atuais e devolve o resultado em destino: campo
// ausente fica como está, campo com null é apagado e objeto aninhado é
// mesclado campo a campo. O tipo de cada campo já foi validado pelo contrato.
func MesclarPatch(atual any, corpo []byte, destino any) (mudou bool, err error) {
	var patch map[string]any
	if err := json.Unmarshal(corpo, &patch); err != nil || patch == nil {
		return false, erros.RequisicaoInvalida.Com("O corpo do PATCH deve ser um objeto JSON.")
	}
	b, err := json.Marshal(atual)
	if err != nil {
		return false, err
	}
	var base map[string]any
	if err := json.Unmarshal(b, &base); err != nil {
		return false, err
	}
	mesclado := mesclar(base, patch)
	b, err = json.Marshal(mesclado)
	if err != nil {
		return false, err
	}
	if err := json.Unmarshal(b, destino); err != nil {
		return false, erros.EntradaInvalida.Com("O PATCH tem valor de tipo errado.")
	}
	return len(patch) > 0, nil
}

// PatchTemCampoAlem informa se o PATCH traz algum campo além dos indicados.
func PatchTemCampoAlem(corpo []byte, campos ...string) bool {
	var patch map[string]json.RawMessage
	if json.Unmarshal(corpo, &patch) != nil {
		return false
	}
	for k := range patch {
		outro := true
		for _, c := range campos {
			if k == c {
				outro = false
			}
		}
		if outro {
			return true
		}
	}
	return false
}

func mesclar(alvo any, patch map[string]any) map[string]any {
	base, ok := alvo.(map[string]any)
	if !ok || base == nil {
		base = map[string]any{}
	}
	for k, v := range patch {
		switch pv := v.(type) {
		case nil:
			delete(base, k)
		case map[string]any:
			base[k] = mesclar(base[k], pv)
		default:
			base[k] = v
		}
	}
	return base
}

// ---------------------------------------------------- dinheiro e instantes

// Dinheiro converte centavos em texto decimal com duas casas (seção 2.5).
func Dinheiro(centavos *int64) *string {
	if centavos == nil {
		return nil
	}
	c := *centavos
	sinal := ""
	if c < 0 {
		sinal, c = "-", -c
	}
	s := fmt.Sprintf("%s%d.%02d", sinal, c/100, c%100)
	return &s
}

// Centavos lê texto decimal com até duas casas ("125.5", "125.50").
func Centavos(campo string, texto *string) (*int64, error) {
	if texto == nil {
		return nil, nil
	}
	t := strings.TrimSpace(*texto)
	inteira, frac, temPonto := strings.Cut(t, ".")
	invalido := erros.EntradaInvalida.ComCampo(campo, erros.CampoFormatoInvalido, `Use valor decimal com ponto e até duas casas, como "125.50".`)
	if inteira == "" || len(inteira) > 10 || (temPonto && (frac == "" || len(frac) > 2)) {
		return nil, invalido
	}
	for _, r := range inteira + frac {
		if r < '0' || r > '9' {
			return nil, invalido
		}
	}
	for len(frac) < 2 {
		frac += "0"
	}
	n, err := strconv.ParseInt(inteira+frac, 10, 64)
	if err != nil {
		return nil, invalido
	}
	return &n, nil
}

// Instante normaliza para UTC com precisão de segundos (seção 2.4).
func Instante(t time.Time) time.Time { return t.UTC().Truncate(time.Second) }

// InstanteOpcional normaliza o ponteiro.
func InstanteOpcional(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	v := Instante(*t)
	return &v
}
