// Package repo dá ao módulo ADM as consultas geradas pelo sqlc (adm.sql) na
// transação do contexto e traduz as violações do banco para erros de domínio.
package repo

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"

	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/db/consultas"
	"tenaz/internal/plataforma/erros"
)

// Q devolve as consultas na transação aberta por NaEmpresa.
func Q(ctx context.Context) (*consultas.Queries, error) {
	return db.Q(ctx)
}

// campoUnico liga cada índice único ao campo da entrada.
var campoUnico = map[string]string{
	"uq_empresa_cnpj":         "cnpj",
	"uq_unidade_codigo":       "codigo",
	"uq_turno_unidade_inicio": "hora_inicio",
	"uq_feriado_unidade_data": "data",
	"uq_pessoa_cpf":           "cpf",
	"uq_usuario_email":        "email",
	"uq_usuario_pessoa":       "pessoa_id",
	"uq_usuario_identidade":   "email",
	"uq_usuario_permissao":    "permissao",
	"uq_perfil_codigo":        "codigo",
}

// campoChave liga cada chave estrangeira ao campo da entrada. Um id de outra
// empresa também cai aqui: a FK composta com empresa_id não o encontra, e a
// resposta é "não encontrado", nunca "de outra empresa".
var campoChave = map[string]string{
	"fk_unidade_empresa":         "empresa_id",
	"fk_turno_unidade_unidade":   "unidade_id",
	"fk_feriado_unidade_unidade": "unidade_id",
	"fk_pessoa_unidade":          "unidade_id",
	"fk_usuario_pessoa":          "pessoa_id",
	"fk_usuario_unidade":         "unidade_padrao_id",
	"fk_usuario_perfil_perfil":   "perfis",
	"fk_usuario_perfil_unidade":  "unidade_id",
}

// Traduzir converte violação de unicidade em conflito (409), de chave
// estrangeira e de CHECK em validacao_falhou (422), com o campo. Outros erros
// voltam como vieram.
func Traduzir(err error) error {
	var pg *pgconn.PgError
	if !errors.As(err, &pg) {
		return err
	}
	switch pg.Code {
	case "23505":
		campo := campoUnico[pg.ConstraintName]
		if campo == "" {
			campo = pg.ConstraintName
		}
		return erros.Duplicado.ComCampo(campo, erros.CampoDuplicado, "Já existe um registro com este valor.").Por(err)
	case "23503":
		campo := campoChave[pg.ConstraintName]
		if campo == "" {
			campo = pg.ConstraintName
		}
		return erros.EntradaInvalida.ComCampo(campo, erros.CampoNaoEncontrado, "Registro referenciado não encontrado.").Por(err)
	case "23514":
		// CHECK ou gatilho de regra (perfil externo, permissão vedada): a
		// mensagem do banco é técnica e não vai para o cliente.
		return erros.EntradaInvalida.Com("a combinação de dados não é permitida").Por(err)
	}
	return err
}
