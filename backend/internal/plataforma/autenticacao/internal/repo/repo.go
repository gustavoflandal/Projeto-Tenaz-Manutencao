// Package repo é o único acesso às tabelas de credencial do esquema
// plataforma: identidade, acesso_empresa, sessao, token_renovacao,
// codigo_autorizacao, desafio_mfa e token_email.
//
// Usa um pool próprio, injetado, e nunca o pool da aplicação nem NaEmpresa
// (achado ISO-E01-02 da revisão de isolamento): quem lê estas tabelas lê
// hashes de senha e segredos de MFA de todas as identidades e pode criar
// sessão para qualquer usuário, então o papel de banco que as alcança tem de
// ser outro que não o tenaz_app. O pacote fica em autenticacao/internal: o
// compilador impede qualquer outro pacote de importá-lo.
//
// O pool é o do papel tenaz_autenticacao (ADR-011; seção 6.9 do documento 5),
// aberto por TENAZ_BANCO_URL_AUTENTICACAO só no modo api. Duas transações:
// Tx, sem empresa, para achar credenciais por hash, id, identidade_id,
// (empresa_id, usuario_id) ou e-mail exato; e NaTrilha, para gravar o evento
// de autenticação junto da mudança de credencial.
package repo

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"tenaz/internal/plataforma/autenticacao/internal/repo/gerado"
	"tenaz/internal/plataforma/db"
)

// Repo abre transações no pool da autenticação.
type Repo struct {
	pool *pgxpool.Pool
}

// Novo recebe o pool da autenticação.
func Novo(pool *pgxpool.Pool) *Repo {
	return &Repo{pool: pool}
}

// ErrDentroDeTransacaoDaAplicacao: o repositório foi chamado de dentro de
// NaEmpresa. Além de misturar os dois pools, a transação da autenticação não
// enxergaria as linhas ainda não confirmadas da aplicação (um usuário recém
// criado, por exemplo) e a chave estrangeira falharia. Chame depois do COMMIT.
var ErrDentroDeTransacaoDaAplicacao = errors.New("autenticação chamada dentro de transação da aplicação")

// Tx executa fn numa transação do pool da autenticação.
func (r *Repo) Tx(ctx context.Context, fn func(q *gerado.Queries) error) error {
	if db.EmTransacao(ctx) {
		return ErrDentroDeTransacaoDaAplicacao
	}
	return pgx.BeginTxFunc(ctx, r.pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		return fn(gerado.New(tx))
	})
}

// NaTrilha executa fn numa transação do pool da autenticação com o contexto
// mínimo para gravar na trilha da empresa (seção 6.9 do documento 5): só
// app.empresa_id e os dados da requisição (origem, IP, user agent,
// request_id), com set_config(..., true). Sem escopo, sem usuário e sem
// app.conferir_escopo(), que lê app.usuario e falharia por falta de
// privilégio. A mudança de credencial e o evento vão na mesma transação.
func (r *Repo) NaTrilha(ctx context.Context, empresaID uuid.UUID, fn func(q *gerado.Queries) error) error {
	if db.EmTransacao(ctx) {
		return ErrDentroDeTransacaoDaAplicacao
	}
	if empresaID == uuid.Nil {
		return errors.New("trilha sem empresa")
	}
	req := db.RequisicaoDe(ctx).Saneada()
	return pgx.BeginTxFunc(ctx, r.pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `SELECT
  set_config('app.empresa_id', $1, true),
  set_config('app.origem',     $2, true),
  set_config('app.ip',         $3, true),
  set_config('app.user_agent', $4, true),
  set_config('app.request_id', $5, true)`,
			empresaID.String(), req.Origem, req.IP, req.UserAgent, req.RequestID)
		if err != nil {
			return fmt.Errorf("fixar contexto da trilha: %w", err)
		}
		return fn(gerado.New(tx))
	})
}
