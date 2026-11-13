// Package paginacao implementa o cursor opaco das listas (seção 2.6 do
// documento 3): limite padrão 50 e máximo 200, cursor válido por 24 horas e
// preso à consulta que o gerou (mudar filtro ou ordem exige começar de novo).
//
// O cursor não é segredo nem autorização: toda consulta filtra por empresa_id
// e passa pelo RLS, então um cursor adulterado só muda a posição na lista da
// própria empresa.
package paginacao

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"tenaz/internal/plataforma/erros"
)

// Limites da seção 2.6 do documento 3.
const (
	LimitePadrao = 50
	LimiteMaximo = 200
	Validade     = 24 * time.Hour
)

// Posicao é o último item da página: a chave de ordenação e o id, que
// desempata.
type Posicao struct {
	Instante time.Time
	ID       uuid.UUID
}

type cursor struct {
	Instante string `json:"t"`
	ID       string `json:"i"`
	Consulta string `json:"q"`
	Emitido  int64  `json:"e"`
}

// Limite aplica o padrão e o teto.
func Limite(pedido int32) int32 {
	switch {
	case pedido <= 0:
		return LimitePadrao
	case pedido > LimiteMaximo:
		return LimiteMaximo
	}
	return pedido
}

// Assinatura resume os filtros e a ordem da consulta. Quem lista monta a
// assinatura com os mesmos valores na ida e na volta.
func Assinatura(partes ...string) string {
	h := sha256.New()
	for _, p := range partes {
		h.Write([]byte(p))
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil)[:8])
}

// Codificar monta o cursor da próxima página.
func Codificar(p Posicao, consulta string, agora time.Time) string {
	b, _ := json.Marshal(cursor{
		Instante: p.Instante.UTC().Format(time.RFC3339Nano), ID: p.ID.String(), Consulta: consulta, Emitido: agora.Unix(),
	})
	return base64.RawURLEncoding.EncodeToString(b)
}

// Decodificar lê o cursor; vazio devolve nil (primeira página). Cursor
// corrompido, de outra consulta ou vencido: cursor_invalido.
func Decodificar(s, consulta string, agora time.Time) (*Posicao, error) {
	if s == "" {
		return nil, nil
	}
	invalido := erros.CursorInvalido.Com("cursor corrompido, de outra consulta ou vencido; recomece sem cursor")
	if len(s) > 512 {
		return nil, invalido
	}
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return nil, invalido
	}
	var c cursor
	if err := json.Unmarshal(b, &c); err != nil || c.Consulta != consulta {
		return nil, invalido
	}
	emitido := time.Unix(c.Emitido, 0)
	if agora.Sub(emitido) > Validade || emitido.After(agora.Add(time.Minute)) {
		return nil, invalido
	}
	t, err := time.Parse(time.RFC3339Nano, c.Instante)
	if err != nil {
		return nil, invalido
	}
	id, err := uuid.Parse(c.ID)
	if err != nil {
		return nil, invalido
	}
	return &Posicao{Instante: t, ID: id}, nil
}

// Pagina é o envelope de uma lista (dados + paginacao, seção 2.6).
type Pagina[T any] struct {
	Dados         []T
	ProximoCursor string // vazio quando não há mais
	TemMais       bool
	Limite        int32
}

// Montar corta a página: a consulta busca limite+1 linhas, e a linha extra
// só diz se há mais. posicao extrai a chave do último item.
func Montar[T any](linhas []T, limite int32, consulta string, agora time.Time, posicao func(T) Posicao) Pagina[T] {
	p := Pagina[T]{Dados: linhas, Limite: limite}
	if int32(len(linhas)) > limite {
		p.Dados = linhas[:limite]
		p.TemMais = true
		p.ProximoCursor = Codificar(posicao(p.Dados[len(p.Dados)-1]), consulta, agora)
	}
	if p.Dados == nil {
		p.Dados = []T{}
	}
	return p
}
