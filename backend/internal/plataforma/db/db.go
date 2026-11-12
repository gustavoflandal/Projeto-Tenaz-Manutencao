// Package db dá acesso ao PostgreSQL com o contexto de empresa e de escopo
// sempre definido na transação.
//
// Esta é a terceira das quatro camadas de isolamento da seção 1.2 de
// docs/2_ARQUITETURA_BANCO_DADOS.md. As outras três são o RLS forçado no
// banco, o filtro explícito por empresa_id em toda consulta e as chaves
// estrangeiras compostas com empresa_id.
//
// Regra inviolável: nenhum código toca dado de cliente fora de NaEmpresa.
package db

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"tenaz/internal/plataforma/erros"
)

// Tipos de escopo. São só estes quatro, pela seção 3.6 do sumário executivo.
const (
	EscopoInterno   = "interno"
	EscopoPrestador = "prestador"
	EscopoLocatario = "locatario"
	EscopoLink      = "link"
)

// Escopo é a camada de isolamento adicional aplicada dentro da empresa, para
// prestador, locatário e link de uso único (risco R4 do sumário).
//
// O escopo é resolvido no servidor a partir do banco, a cada requisição.
// Nunca vem do token, do corpo da requisição nem da URL.
type Escopo struct {
	Tipo string

	FornecedorID *uuid.UUID // escopo prestador
	LocatarioID  *uuid.UUID // escopo locatario
	LinkHashHex  string     // escopo link: sha256 do segredo, em hexadecimal

	UsuarioID    *uuid.UUID
	PodeLerSaude bool // permissão pessoa.saude.ler
}

// Interno é o escopo de usuários da própria empresa, de jobs e de chaves de API.
func Interno(usuarioID *uuid.UUID) Escopo {
	return Escopo{Tipo: EscopoInterno, UsuarioID: usuarioID}
}

// Valido recusa escopo incoerente antes de chegar ao banco. O banco também
// recusa, por app.conferir_escopo(); esta é a checagem barata que vem antes.
func (e Escopo) Valido() error {
	switch e.Tipo {
	case EscopoInterno:
		if e.FornecedorID != nil || e.LocatarioID != nil || e.LinkHashHex != "" {
			return errors.New("escopo interno não pode ter fornecedor, locatário nem link")
		}
	case EscopoPrestador:
		if e.FornecedorID == nil {
			return errors.New("escopo prestador exige fornecedor_id")
		}
		if e.UsuarioID == nil {
			return errors.New("escopo prestador exige usuário")
		}
	case EscopoLocatario:
		if e.LocatarioID == nil {
			return errors.New("escopo locatario exige locatario_id")
		}
		if e.UsuarioID == nil {
			return errors.New("escopo locatario exige usuário")
		}
	case EscopoLink:
		if e.LinkHashHex == "" {
			return errors.New("escopo link exige o hash do segredo")
		}
	default:
		return fmt.Errorf("escopo inválido: %q", e.Tipo)
	}
	return nil
}

// DB é o pool de conexões.
type DB struct {
	pool *pgxpool.Pool
}

// Abrir monta o pool.
func Abrir(ctx context.Context, url string, maxConn int32) (*DB, error) {
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		return nil, fmt.Errorf("interpretar TENAZ_BANCO_URL: %w", err)
	}
	if maxConn > 0 {
		cfg.MaxConns = maxConn
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("abrir pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("alcançar o banco: %w", err)
	}
	return &DB{pool: pool}, nil
}

// Fechar devolve as conexões.
func (db *DB) Fechar() {
	if db != nil && db.pool != nil {
		db.pool.Close()
	}
}

// Pool expõe o pool para usos que não tocam dado de empresa: migração,
// verificação de saúde e as suítes de vazamento.
func (db *DB) Pool() *pgxpool.Pool { return db.pool }

type chaveTx struct{}

// NaEmpresa abre uma transação, fixa o contexto de empresa e de escopo e
// coloca a transação no contexto. Todo caso de uso que toca dado de cliente
// passa por aqui.
//
// O contexto é definido com set_config(..., is_local => true), local à
// transação. Nunca SET de sessão: no pool, um SET de sessão vazaria para a
// requisição seguinte, de outra empresa.
func (db *DB) NaEmpresa(ctx context.Context, empresaID uuid.UUID, escopo Escopo, fn func(context.Context) error) error {
	if empresaID == uuid.Nil {
		return erros.SemContextoEmpresa.Com("empresa_id vazio em NaEmpresa")
	}
	if err := escopo.Valido(); err != nil {
		return erros.SemContextoEmpresa.Por(err)
	}

	return pgx.BeginTxFunc(ctx, db.pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		if err := fixarContexto(ctx, tx, empresaID, escopo); err != nil {
			return err
		}
		return fn(context.WithValue(ctx, chaveTx{}, tx))
	})
}

func fixarContexto(ctx context.Context, tx pgx.Tx, empresaID uuid.UUID, escopo Escopo) error {
	definir := func(chave, valor string) error {
		if _, err := tx.Exec(ctx, `SELECT set_config($1, $2, true)`, chave, valor); err != nil {
			return fmt.Errorf("definir %s: %w", chave, err)
		}
		return nil
	}

	if err := definir("app.empresa_id", empresaID.String()); err != nil {
		return err
	}
	if err := definir("app.escopo", escopo.Tipo); err != nil {
		return err
	}
	if err := definir("app.usuario_id", textoUUID(escopo.UsuarioID)); err != nil {
		return err
	}
	if err := definir("app.fornecedor_id", textoUUID(escopo.FornecedorID)); err != nil {
		return err
	}
	if err := definir("app.locatario_id", textoUUID(escopo.LocatarioID)); err != nil {
		return err
	}
	if err := definir("app.link_hash", escopo.LinkHashHex); err != nil {
		return err
	}
	if err := definir("app.saude", simNao(escopo.PodeLerSaude)); err != nil {
		return err
	}

	// O banco confere o escopo contra o cadastro e recusa a transação se
	// divergirem. Falha fechada.
	if _, err := tx.Exec(ctx, `SELECT app.conferir_escopo()`); err != nil {
		return fmt.Errorf("conferir escopo: %w", err)
	}
	return nil
}

func textoUUID(id *uuid.UUID) string {
	if id == nil {
		return ""
	}
	return id.String()
}

func simNao(v bool) string {
	if v {
		return "sim"
	}
	return "nao"
}

// Tx recupera a transação do contexto. Devolve erro quando o chamador está
// fora de NaEmpresa, o que é falha de programação, não condição de execução.
func Tx(ctx context.Context) (pgx.Tx, error) {
	tx, ok := ctx.Value(chaveTx{}).(pgx.Tx)
	if !ok {
		return nil, erros.SemContextoEmpresa.Com("acesso a dado de cliente fora de NaEmpresa")
	}
	return tx, nil
}
