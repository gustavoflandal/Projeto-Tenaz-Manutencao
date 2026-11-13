package vazamento

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
)

// A Tenaz não guarda dado de saúde (decisão do dono do produto, sumário 0.7,
// seção 3.6; migração 00006). Não basta não ter a tabela: a permissão, a
// função de contexto e o CHECK que a mantinham fora dos perfis também saem,
// para que ninguém volte a modelar o dado "porque a permissão já existe".
//
// Verificação: nenhum objeto de app ou plataforma com "saude" no nome
// (tabela, coluna, função, constraint, policy), nenhuma função que leia a
// variável app.saude, e o código pessoa.saude.ler recusado pelo formato de
// permissão.
const sqlRastrosDeSaude = `
SELECT 'tabela ' || c.oid::regclass::text FROM pg_class c JOIN pg_namespace n ON n.oid = c.relnamespace
WHERE n.nspname IN ('app', 'plataforma') AND c.relname ILIKE '%saude%'
UNION ALL
SELECT 'coluna ' || a.attrelid::regclass::text || '.' || a.attname
FROM pg_attribute a JOIN pg_class c ON c.oid = a.attrelid JOIN pg_namespace n ON n.oid = c.relnamespace
WHERE n.nspname IN ('app', 'plataforma') AND a.attnum > 0 AND NOT a.attisdropped AND a.attname ILIKE '%saude%'
UNION ALL
SELECT 'função ' || p.oid::regprocedure::text FROM pg_proc p JOIN pg_namespace n ON n.oid = p.pronamespace
WHERE n.nspname IN ('app', 'plataforma') AND (p.proname ILIKE '%saude%' OR p.prosrc ILIKE '%app.saude%')
UNION ALL
SELECT 'constraint ' || con.conname FROM pg_constraint con JOIN pg_namespace n ON n.oid = con.connamespace
WHERE n.nspname IN ('app', 'plataforma') AND (con.conname ILIKE '%saude%' OR pg_get_constraintdef(con.oid) ILIKE '%saude%')
UNION ALL
SELECT 'policy ' || pol.polname || ' em ' || pol.polrelid::regclass::text FROM pg_policy pol
WHERE pol.polname ILIKE '%saude%' OR coalesce(pg_get_expr(pol.polqual, pol.polrelid), '') ILIKE '%saude%'
   OR coalesce(pg_get_expr(pol.polwithcheck, pol.polrelid), '') ILIKE '%saude%'
ORDER BY 1`

func TestNaoHaDadoDeSaude(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	for _, r := range listarTexto(t, ctx, pool, sqlRastrosDeSaude) {
		t.Errorf("rastro de dado de saúde no banco: %s; a Tenaz não guarda dado de saúde (sumário 0.7, 3.6)", r)
	}
	err := comoPapel(ctx, papelApp, f.internoA(), func(tx pgx.Tx) error {
		_, err := tx.Exec(ctx, `INSERT INTO app.usuario_permissao (empresa_id, usuario_id, permissao)
		                        VALUES ($1, $2, 'pessoa.saude.ler')`, f.A.ID, f.A.UsuarioInativo)
		return err
	})
	exigirErro(t, "conceder pessoa.saude.ler", err, erroCheck)
}

func TestDetectorAcusaRastroDeSaude(t *testing.T) {
	ctx := context.Background()
	// Só objetos novos, da própria transação: nenhuma trava em tabela real
	// (app.pessoa, app.perfil), que outros pacotes podem estar gravando quando
	// todos usam o mesmo banco local (tarefas.sh testar-rapido).
	tx := criarTabelasDeProva(t, ctx, []string{
		"CREATE FUNCTION app.pode_ler_saude() RETURNS boolean LANGUAGE sql STABLE " +
			"AS $$ SELECT coalesce(current_setting('app.saude', true), '') = 'sim' $$",
		"CREATE TABLE app.tabela_prova_saude (empresa_id uuid NOT NULL, aso_saude date, " +
			"permissoes text[] NOT NULL DEFAULT '{}', padrao boolean NOT NULL DEFAULT false, " +
			"CONSTRAINT ck_prova_perfil_saude CHECK (NOT padrao OR NOT ('pessoa.saude.ler' = ANY (permissoes))))",
	})
	r := listarTexto(t, ctx, tx, sqlRastrosDeSaude)
	for _, esperado := range []string{"função app.pode_ler_saude()", "tabela app.tabela_prova_saude",
		"coluna app.tabela_prova_saude.aso_saude", "constraint ck_prova_perfil_saude"} {
		if !contem(r, esperado) {
			t.Errorf("não acusou %q; acusou %v", esperado, r)
		}
	}
}
