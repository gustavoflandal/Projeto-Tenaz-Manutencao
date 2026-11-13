package auditoria_test

import (
	"encoding/json"
	"sort"
	"strings"
	"testing"
	"time"

	"tenaz/internal/plataforma/auditoria"
)

// O texto do hash tem as 29 chaves do formato 1, em ordem de bytes, sem
// espaços, e é JSON válido (documento 2, seção 9.2.1).
func TestSerializacaoDoFormato1(t *testing.T) {
	nome := "Zé \"da Silva\" \\ \n\t\x01 日本 <&>"
	e := auditoria.ConteudoEvento{
		Acao: "alteracao", AtorTipo: "usuario", EmpresaID: "01a1d4eb-cbb3-71ec-95c3-1fea241f3dee", Entidade: "unidade",
		ID: "01a1d4eb-cbb3-71ec-95c3-1fea241f3def", Resultado: "sucesso", UsuarioNome: &nome,
		OcorridoEm: time.Date(2026, 11, 13, 9, 30, 15, 123456000, time.FixedZone("BRT", -3*3600)),
	}
	s := e.Serializar()
	if strings.ContainsAny(strings.ReplaceAll(s, nome, ""), " \n") && !strings.Contains(s, `"usuario_nome":"Zé`) {
		t.Errorf("espaço fora das cadeias: %s", s)
	}
	var m map[string]json.RawMessage
	if err := json.Unmarshal([]byte(s), &m); err != nil {
		t.Fatalf("não é JSON: %v\n%s", err, s)
	}
	if len(m) != 29 {
		t.Errorf("esperadas 29 chaves, vieram %d", len(m))
	}
	// ordem de bytes das chaves, na ordem em que aparecem no texto
	dec := json.NewDecoder(strings.NewReader(s))
	if _, err := dec.Token(); err != nil {
		t.Fatal(err)
	}
	var ordem []string
	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			t.Fatal(err)
		}
		ordem = append(ordem, tok.(string))
		var valor json.RawMessage
		if err := dec.Decode(&valor); err != nil {
			t.Fatal(err)
		}
	}
	if len(ordem) != 29 || !sort.StringsAreSorted(ordem) {
		t.Errorf("chaves fora da ordem de bytes ou em número errado: %v", ordem)
	}
	if string(m["formato"]) != "1" || string(m["ocorrido_em"]) != `"2026-11-13T12:30:15.123456Z"` {
		t.Errorf("formato %s, ocorrido_em %s", m["formato"], m["ocorrido_em"])
	}
	// o escape é o do PostgreSQL: só aspas, barra, controles; < > & e U+2028 ficam como estão
	esperado := `"Zé \"da Silva\" \\ \n\t\u0001 日本 <&>"`
	if string(m["usuario_nome"]) != esperado {
		t.Errorf("usuario_nome: %s, esperado %s", m["usuario_nome"], esperado)
	}
	for _, nula := range []string{"antes", "depois", "detalhe", "campos_alterados", "ip", "tela", "modulo", "transacao", "usuario_id"} {
		if string(m[nula]) != "null" {
			t.Errorf("%s ausente deve ser null: %s", nula, m[nula])
		}
	}
	if h1, h2 := e.Hash(), e.Hash(); h1 != h2 {
		t.Error("hash não determinístico")
	}
	e2 := e
	e2.Sensivel = true
	if e.Hash() == e2.Hash() {
		t.Error("o hash tem de cobrir sensivel")
	}
}
