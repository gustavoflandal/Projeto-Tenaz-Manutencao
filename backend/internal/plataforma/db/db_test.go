package db

import (
	"context"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestEscopoValido(t *testing.T) {
	u, f := uuid.New(), uuid.New()
	hash := strings.Repeat("ab", 32)
	validos := []Escopo{
		Interno(nil), Interno(&u), {Tipo: EscopoInterno, UsuarioID: &u},
		{Tipo: EscopoPrestador, FornecedorID: &f, UsuarioID: &u},
		{Tipo: EscopoLocatario, LocatarioID: &f, UsuarioID: &u},
		{Tipo: EscopoLink, LinkHashHex: hash},
	}
	for _, e := range validos {
		if err := e.Valido(); err != nil {
			t.Errorf("%+v: %v", e, err)
		}
	}
	invalidos := []Escopo{
		{Tipo: ""}, {Tipo: "fornecedor"}, {Tipo: EscopoInterno, FornecedorID: &f},
		{Tipo: EscopoPrestador, UsuarioID: &u}, {Tipo: EscopoPrestador, FornecedorID: &f},
		{Tipo: EscopoLink, LinkHashHex: "abc"}, {Tipo: EscopoLink, LinkHashHex: hash, UsuarioID: &u},
	}
	for _, e := range invalidos {
		if e.Valido() == nil {
			t.Errorf("%+v deveria ser recusado", e)
		}
	}
}

func TestRequisicaoNormalizada(t *testing.T) {
	u := uuid.New()
	r, err := Requisicao{IP: "::ffff:203.0.113.9", UserAgent: strings.Repeat("ã", 400), Origem: OrigemWeb}.normalizada(Interno(&u))
	if err != nil || r.AtorTipo != AtorUsuario || r.IP != "203.0.113.9" || len(r.UserAgent) > maxUserAgent {
		t.Fatalf("normalizada: %+v %v", r, err)
	}
	r, _ = Requisicao{IP: "não é ip"}.normalizada(Interno(nil))
	if r.AtorTipo != AtorSistema || r.IP != "" {
		t.Errorf("sem usuário é sistema; IP inválido some: %+v", r)
	}
	if _, err := (Requisicao{Origem: "telepatia"}).normalizada(Interno(nil)); err == nil {
		t.Error("origem fora da lista")
	}
	if _, err := (Requisicao{AtorTipo: AtorSincronizacaoConexao}).normalizada(Interno(nil)); err == nil {
		t.Error("sincronização só no pool do sincronizador")
	}
	if (Requisicao{Origem: "x"}).Saneada().Origem != "" {
		t.Error("Saneada descarta origem desconhecida")
	}
	agora := time.Date(2026, 11, 12, 10, 0, 0, 0, time.FixedZone("BRT", -3*3600))
	if textoInstante(&agora) != "2026-11-12T13:00:00Z" {
		t.Errorf("instante do aparelho em UTC: %s", textoInstante(&agora))
	}
}

func TestTxForaDeTransacao(t *testing.T) {
	if _, err := Tx(context.Background()); err == nil {
		t.Error("Tx fora de NaEmpresa é erro")
	}
	if EmTransacao(context.Background()) {
		t.Error("contexto vazio não está em transação")
	}
}
