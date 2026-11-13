package auditoria_test

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"slices"
	"testing"
	"time"

	"github.com/google/uuid"

	"tenaz/internal/plataforma/auditoria"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/db/dbteste"
	"tenaz/internal/plataforma/erros"
)

// eventoGravado é a linha de app.evento_auditoria, lida pelo migrador.
type eventoGravado struct {
	OcorridoEm               time.Time
	AtorTipo, Acao, Entidade string
	UsuarioID, EntidadeID    *uuid.UUID
	EntidadeChave            *string
	Antes, Depois            []byte
	Campos                   []string
	Origem, IP, UserAgent    *string
	RequestID, Escopo        *string
}

// Registrar grava, na transação de NaEmpresa, antes, depois, autor, instante
// e o contexto técnico da requisição (RF-ADM-011; seção 9.1 do documento 2);
// desfeita a transação, o evento some; fora de NaEmpresa ou em outra empresa,
// nada é gravado.
func TestRegistrar(t *testing.T) {
	ctx := context.Background()
	amb, err := dbteste.Preparar(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(amb.Encerrar)

	ler := func(t *testing.T, empresa uuid.UUID, entidadeID uuid.UUID) []eventoGravado {
		t.Helper()
		linhas, err := amb.Migrador.Query(ctx, `SELECT ocorrido_em, ator_tipo, acao, entidade, usuario_id, entidade_id,
  entidade_chave, antes, depois, campos_alterados, origem, host(ip), user_agent, request_id, escopo
FROM app.evento_auditoria WHERE empresa_id = $1 AND entidade_id = $2 ORDER BY ocorrido_em, id`, empresa, entidadeID)
		if err != nil {
			t.Fatal(err)
		}
		defer linhas.Close()
		var out []eventoGravado
		for linhas.Next() {
			var e eventoGravado
			if err := linhas.Scan(&e.OcorridoEm, &e.AtorTipo, &e.Acao, &e.Entidade, &e.UsuarioID, &e.EntidadeID,
				&e.EntidadeChave, &e.Antes, &e.Depois, &e.Campos, &e.Origem, &e.IP, &e.UserAgent, &e.RequestID, &e.Escopo); err != nil {
				t.Fatal(err)
			}
			out = append(out, e)
		}
		if err := linhas.Err(); err != nil {
			t.Fatal(err)
		}
		return out
	}
	relogioDoBanco := func(t *testing.T) time.Time {
		t.Helper()
		var agora time.Time
		if err := amb.Migrador.QueryRow(ctx, `SELECT clock_timestamp()`).Scan(&agora); err != nil {
			t.Fatal(err)
		}
		return agora
	}
	jsonIgual := func(t *testing.T, nome string, got []byte, quer any) {
		t.Helper()
		var g, q any
		if err := json.Unmarshal(got, &g); err != nil {
			t.Fatalf("%s não é JSON: %s", nome, got)
		}
		b, _ := json.Marshal(quer)
		_ = json.Unmarshal(b, &q)
		if !reflect.DeepEqual(g, q) {
			t.Errorf("%s: gravado %s, esperado %s", nome, got, b)
		}
	}

	empresa := uuid.Must(uuid.NewV7())
	autor := uuid.Must(uuid.NewV7())
	req := db.Requisicao{Origem: db.OrigemWeb, IP: "203.0.113.9", UserAgent: "teste-auditoria", RequestID: uuid.NewString()}
	ctxReq := db.ComRequisicao(ctx, req)

	t.Run("antes, depois, autor e instante", func(t *testing.T) {
		alvo := uuid.Must(uuid.NewV7())
		antes := map[string]any{"nome": "Matriz", "situacao": "ativa"}
		depois := map[string]any{"nome": "Matriz Sul", "situacao": "ativa"}
		inicio := relogioDoBanco(t)
		err := amb.App.NaEmpresa(ctxReq, empresa, db.Interno(nil), func(ctx context.Context) error {
			return auditoria.Registrar(ctx, auditoria.Evento{
				EmpresaID: empresa, AtorTipo: db.AtorUsuario, UsuarioID: &autor, Acao: auditoria.AcaoAlteracao,
				Entidade: "unidade", EntidadeID: &alvo, Antes: antes, Depois: depois, CamposAlterados: []string{"nome"},
			})
		})
		fim := relogioDoBanco(t)
		if err != nil {
			t.Fatal(err)
		}
		ev := ler(t, empresa, alvo)
		if len(ev) != 1 {
			t.Fatalf("esperado 1 evento, vieram %d", len(ev))
		}
		e := ev[0]
		jsonIgual(t, "antes", e.Antes, antes)
		jsonIgual(t, "depois", e.Depois, depois)
		if e.AtorTipo != db.AtorUsuario || e.UsuarioID == nil || *e.UsuarioID != autor {
			t.Errorf("autor: %s %v", e.AtorTipo, e.UsuarioID)
		}
		if e.Acao != auditoria.AcaoAlteracao || e.Entidade != "unidade" || !slices.Equal(e.Campos, []string{"nome"}) {
			t.Errorf("ação, entidade e campos: %s %s %v", e.Acao, e.Entidade, e.Campos)
		}
		// O instante é o do banco (now() da transação), comparado com o
		// relógio do próprio banco antes e depois: o relógio da máquina do
		// teste pode se afastar do da VM do banco.
		if e.OcorridoEm.Before(inicio) || e.OcorridoEm.After(fim) {
			t.Errorf("instante %v fora de [%v, %v]", e.OcorridoEm, inicio, fim)
		}
		if e.Origem == nil || *e.Origem != db.OrigemWeb || e.IP == nil || *e.IP != req.IP ||
			e.UserAgent == nil || *e.UserAgent != req.UserAgent || e.RequestID == nil || *e.RequestID != req.RequestID {
			t.Errorf("contexto técnico: origem %v ip %v ua %v request %v", e.Origem, e.IP, e.UserAgent, e.RequestID)
		}
		if e.Escopo != nil && *e.Escopo != db.EscopoInterno {
			t.Errorf("escopo: %v", *e.Escopo)
		}
	})

	t.Run("sem antes nem depois grava nulo; sem ator é sistema", func(t *testing.T) {
		alvo := uuid.Must(uuid.NewV7())
		err := amb.App.NaEmpresa(ctx, empresa, db.Interno(nil), func(ctx context.Context) error {
			return auditoria.Registrar(ctx, auditoria.Evento{
				EmpresaID: empresa, Acao: auditoria.AcaoVisualizacaoSensivel, Entidade: "pessoa", EntidadeID: &alvo,
			})
		})
		if err != nil {
			t.Fatal(err)
		}
		ev := ler(t, empresa, alvo)
		if len(ev) != 1 || ev[0].Antes != nil || ev[0].Depois != nil || ev[0].AtorTipo != db.AtorSistema || ev[0].UsuarioID != nil {
			t.Fatalf("evento mínimo: %+v", ev)
		}
	})

	t.Run("rollback desfaz o evento", func(t *testing.T) {
		alvo := uuid.Must(uuid.NewV7())
		falha := errors.New("o fato foi desfeito")
		err := amb.App.NaEmpresa(ctx, empresa, db.Interno(nil), func(ctx context.Context) error {
			if err := auditoria.Registrar(ctx, auditoria.Evento{EmpresaID: empresa, Acao: auditoria.AcaoInclusao,
				Entidade: "unidade", EntidadeID: &alvo, Depois: map[string]string{"nome": "X"}}); err != nil {
				return err
			}
			return falha
		})
		if !errors.Is(err, falha) {
			t.Fatal(err)
		}
		if ev := ler(t, empresa, alvo); len(ev) != 0 {
			t.Errorf("evento sobreviveu ao rollback: %+v", ev)
		}
	})

	t.Run("fora de NaEmpresa, em outra empresa ou com valor não serializável, nada é gravado", func(t *testing.T) {
		alvo := uuid.Must(uuid.NewV7())
		if err := auditoria.Registrar(ctx, auditoria.Evento{EmpresaID: empresa, Acao: auditoria.AcaoInclusao,
			Entidade: "unidade", EntidadeID: &alvo}); !errors.Is(err, erros.SemContextoEmpresa) {
			t.Errorf("fora de NaEmpresa: %v", err)
		}
		outra := uuid.Must(uuid.NewV7())
		err := amb.App.NaEmpresa(ctx, empresa, db.Interno(nil), func(ctx context.Context) error {
			return auditoria.Registrar(ctx, auditoria.Evento{EmpresaID: outra, Acao: auditoria.AcaoInclusao,
				Entidade: "unidade", EntidadeID: &alvo})
		})
		if err == nil {
			t.Error("evento de outra empresa na transação desta: a policy auditoria_insercao recusa")
		}
		err = amb.App.NaEmpresa(ctx, empresa, db.Interno(nil), func(ctx context.Context) error {
			return auditoria.Registrar(ctx, auditoria.Evento{EmpresaID: empresa, Acao: auditoria.AcaoInclusao,
				Entidade: "unidade", EntidadeID: &alvo, Depois: map[string]any{"canal": make(chan int)}})
		})
		if err == nil {
			t.Error("valor não serializável")
		}
		if ev := ler(t, empresa, alvo); len(ev) != 0 {
			t.Errorf("gravado mesmo assim: %+v", ev)
		}
		if ev := ler(t, outra, alvo); len(ev) != 0 {
			t.Errorf("gravado na outra empresa: %+v", ev)
		}
	})
}
