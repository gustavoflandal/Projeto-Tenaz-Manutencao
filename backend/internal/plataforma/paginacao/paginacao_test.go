package paginacao

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"tenaz/internal/plataforma/erros"
)

func TestCursorIdaEVolta(t *testing.T) {
	agora := time.Now()
	p := Posicao{Instante: agora.Add(-time.Hour).Round(0), ID: uuid.New()}
	q := Assinatura("pessoas", "busca=ana")
	c := Codificar(p, q, agora)

	got, err := Decodificar(c, q, agora.Add(time.Hour))
	if err != nil || got == nil || !got.Instante.Equal(p.Instante) || got.ID != p.ID {
		t.Fatalf("ida e volta: %+v %v", got, err)
	}
	if got, err := Decodificar("", q, agora); got != nil || err != nil {
		t.Error("cursor vazio é a primeira página")
	}
	for nome, caso := range map[string]func() (*Posicao, error){
		"outra consulta": func() (*Posicao, error) { return Decodificar(c, Assinatura("pessoas", "busca=bia"), agora) },
		"vencido":        func() (*Posicao, error) { return Decodificar(c, q, agora.Add(25*time.Hour)) },
		"corrompido":     func() (*Posicao, error) { return Decodificar(c[:len(c)-3]+"xyz", q, agora) },
		"lixo":           func() (*Posicao, error) { return Decodificar("%%%", q, agora) },
	} {
		if _, err := caso(); !errors.Is(err, erros.CursorInvalido) {
			t.Errorf("%s: esperado CURSOR_INVALIDO, veio %v", nome, err)
		}
	}
}

func TestLimiteEMontar(t *testing.T) {
	if Limite(0) != 50 || Limite(500) != 200 || Limite(10) != 10 {
		t.Error("limites da seção 2.6")
	}
	linhas := []int{1, 2, 3}
	p := Montar(linhas, 2, "q", time.Now(), func(int) Posicao { return Posicao{Instante: time.Now(), ID: uuid.New()} })
	if len(p.Dados) != 2 || !p.TemMais || p.ProximoCursor == "" {
		t.Errorf("página cortada: %+v", p)
	}
	p = Montar([]int{}, 2, "q", time.Now(), nil)
	if p.TemMais || p.ProximoCursor != "" || p.Dados == nil {
		t.Errorf("página vazia: %+v", p)
	}
}
