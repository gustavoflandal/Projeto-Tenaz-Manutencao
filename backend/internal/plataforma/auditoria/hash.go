package auditoria

import (
	"crypto/sha256"
	"strings"
	"time"
	"unicode/utf8"
)

// FormatoHash versiona a regra do hash do evento (documento 2, seção 9.2.1).
// Mudar a regra exige formato 2 e uma regra por formato na verificação.
const FormatoHash = 1

// ConteudoEvento é o evento como o hash o enxerga: todos os campos menos o
// próprio hash, cada um já no texto que o banco lhe dá. É o que um
// verificador externo recebe da exportação (RF-ADM-013) e o que o teste usa
// para recalcular o hash no Go e compará-lo com o que app.hash_evento_auditoria
// gravou.
//
// Campos opcionais são ponteiros: nil é null. Os quatro campos jsonb (Antes,
// Depois, Detalhe, CamposAlterados) vêm como o texto do jsonb do PostgreSQL
// (jsonb_out), que o formato 1 adota sem reescrever: chaves ordenadas pelo
// tamanho e depois pelos bytes, ": " e ", " como separadores.
type ConteudoEvento struct {
	Acao                    string
	Antes                   *string
	AtorTipo                string
	CamposAlterados         *string
	ConexaoEmpresaID        *string
	Depois                  *string
	Detalhe                 *string
	DispositivoID           *string
	EmpresaID               string
	Entidade                string
	EntidadeChave           *string
	EntidadeID              *string
	Escopo                  *string
	ID                      string
	IP                      *string // texto do inet, como o PostgreSQL o devolve
	Modulo                  *string
	OcorridoEm              time.Time
	OcorridoNoDispositivoEm *time.Time
	Origem                  *string
	RequestID               *string
	Resultado               string
	Sensivel                bool
	Tela                    *string
	TokenAPIID              *string
	Transacao               *string // número do xid8, como texto
	UserAgent               *string
	UsuarioID               *string
	UsuarioNome             *string
}

// Hash calcula o sha256 do conteúdo (formato 1): o texto UTF-8 de um objeto
// JSON sem espaços, com as 29 chaves em ordem de bytes.
func (e ConteudoEvento) Hash() [32]byte {
	return sha256.Sum256([]byte(e.Serializar()))
}

// Serializar devolve o texto que o hash cobre.
func (e ConteudoEvento) Serializar() string {
	var b strings.Builder
	cad := func(chave string, v *string) {
		b.WriteString(`"` + chave + `":`)
		if v == nil {
			b.WriteString("null")
			return
		}
		escreverCadeia(&b, *v)
	}
	cadObrig := func(chave, v string) { cad(chave, &v) }
	cru := func(chave string, v *string) { // jsonb já em texto
		b.WriteString(`"` + chave + `":`)
		if v == nil {
			b.WriteString("null")
			return
		}
		b.WriteString(*v)
	}
	instante := func(chave string, t *time.Time) {
		if t == nil {
			cad(chave, nil)
			return
		}
		s := t.UTC().Format("2006-01-02T15:04:05.000000Z")
		cad(chave, &s)
	}
	sep := func() { b.WriteByte(',') }

	b.WriteByte('{')
	cadObrig("acao", e.Acao)
	sep()
	cru("antes", e.Antes)
	sep()
	cadObrig("ator_tipo", e.AtorTipo)
	sep()
	cru("campos_alterados", e.CamposAlterados)
	sep()
	cad("conexao_empresa_id", e.ConexaoEmpresaID)
	sep()
	cru("depois", e.Depois)
	sep()
	cru("detalhe", e.Detalhe)
	sep()
	cad("dispositivo_id", e.DispositivoID)
	sep()
	cadObrig("empresa_id", e.EmpresaID)
	sep()
	cadObrig("entidade", e.Entidade)
	sep()
	cad("entidade_chave", e.EntidadeChave)
	sep()
	cad("entidade_id", e.EntidadeID)
	sep()
	cad("escopo", e.Escopo)
	sep()
	b.WriteString(`"formato":1`)
	sep()
	cadObrig("id", e.ID)
	sep()
	cad("ip", e.IP)
	sep()
	cad("modulo", e.Modulo)
	sep()
	t := e.OcorridoEm
	instante("ocorrido_em", &t)
	sep()
	instante("ocorrido_no_dispositivo_em", e.OcorridoNoDispositivoEm)
	sep()
	cad("origem", e.Origem)
	sep()
	cad("request_id", e.RequestID)
	sep()
	cadObrig("resultado", e.Resultado)
	sep()
	if e.Sensivel {
		b.WriteString(`"sensivel":true`)
	} else {
		b.WriteString(`"sensivel":false`)
	}
	sep()
	cad("tela", e.Tela)
	sep()
	cad("token_api_id", e.TokenAPIID)
	sep()
	cad("transacao", e.Transacao)
	sep()
	cad("user_agent", e.UserAgent)
	sep()
	cad("usuario_id", e.UsuarioID)
	sep()
	cad("usuario_nome", e.UsuarioNome)
	b.WriteByte('}')
	return b.String()
}

const hexMinusculo = "0123456789abcdef"

// escreverCadeia escreve a cadeia como o to_jsonb do PostgreSQL: aspas e barra
// invertida com barra, \b \f \n \r \t, \u00xx minúsculo nos demais controles
// (abaixo de 0x20), e o resto em UTF-8 puro. Não é o encoding/json, que também
// escapa <, >, & e U+2028/U+2029.
func escreverCadeia(b *strings.Builder, s string) {
	b.WriteByte('"')
	for i := 0; i < len(s); {
		r, tam := utf8.DecodeRuneInString(s[i:])
		switch {
		case r == '"':
			b.WriteString(`\"`)
		case r == '\\':
			b.WriteString(`\\`)
		case r == '\b':
			b.WriteString(`\b`)
		case r == '\f':
			b.WriteString(`\f`)
		case r == '\n':
			b.WriteString(`\n`)
		case r == '\r':
			b.WriteString(`\r`)
		case r == '\t':
			b.WriteString(`\t`)
		case r < 0x20:
			b.WriteString(`\u00` + string(hexMinusculo[r>>4]) + string(hexMinusculo[r&0x0f]))
		default:
			b.WriteString(s[i : i+tam])
		}
		i += tam
	}
	b.WriteByte('"')
}
