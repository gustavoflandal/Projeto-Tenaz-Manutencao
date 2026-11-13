package autenticacao

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PurgarCredenciaisVencidas apaga as credenciais vencidas há mais de 30 dias
// (sessão, refresh token, código OAuth, desafio de MFA, link de e-mail) e
// devolve quantas linhas saíram. É o corpo do job diário do worker (seção 6.9
// do documento 5): roda no pool do tenaz_job, que não tem privilégio nenhum
// nas tabelas e só executa plataforma.purgar_credenciais_vencidas(),
// SECURITY DEFINER e sem parâmetro.
//
// O agendador de jobs ainda não existe no projeto (River, ADR-005, entra com o
// worker); até lá esta função fica pronta e sem chamador.
func PurgarCredenciaisVencidas(ctx context.Context, poolJob *pgxpool.Pool) (int64, error) {
	var n int64
	if err := poolJob.QueryRow(ctx, `SELECT plataforma.purgar_credenciais_vencidas()`).Scan(&n); err != nil {
		return 0, fmt.Errorf("purgar credenciais vencidas: %w", err)
	}
	return n, nil
}
