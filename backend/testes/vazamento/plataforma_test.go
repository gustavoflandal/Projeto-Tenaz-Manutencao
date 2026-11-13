package vazamento

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
)

// Tabelas de autenticação em plataforma, sem RLS (seção 5.1 do documento 2).
// O desenho se apoia em quatro regras: só o papel tenaz_autenticacao as
// alcança (ADR-011; testes em autenticacao_test.go), só hash, empresa fixada
// na criação sem UPDATE e FK composta para app.usuario. Este arquivo prova as
// três últimas NO PAPEL QUE TEM O PRIVILÉGIO: provar no tenaz_app, que já não
// tem nenhum, passaria por vacuidade.

// acesso_empresa só é escrita pelo trigger SECURITY DEFINER.
func TestAppNaoEscreveEmAcessoEmpresa(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)

	for _, c := range []struct {
		rotulo, sql string
		args        []any
	}{
		{"INSERT", `INSERT INTO plataforma.acesso_empresa (identidade_id, empresa_id, usuario_id) VALUES ($1, $2, $3)`,
			[]any{f.IdentidadeF1, f.B.ID, f.B.UsuarioInterno}},
		{"UPDATE", `UPDATE plataforma.acesso_empresa SET habilitado = true WHERE identidade_id = $1`, []any{f.A.Identidade}},
		{"DELETE", `DELETE FROM plataforma.acesso_empresa WHERE identidade_id = $1`, []any{f.A.Identidade}},
	} {
		for _, papel := range []string{papelAutenticacao, papelApp, papelJob, papelLeitura, papelSinc} {
			err := comoPapel(ctx, papel, f.internoA(), func(tx pgx.Tx) error {
				_, err := tx.Exec(ctx, c.sql, c.args...)
				return err
			})
			exigirErro(t, papel+": "+c.rotulo+" em acesso_empresa", err, erroPrivilegio)
		}
	}

	// O trigger cria o acesso a partir da linha real de app.usuario.
	var n int64
	if err := pool.QueryRow(ctx, `SELECT count(*) FROM plataforma.acesso_empresa
	                              WHERE identidade_id = $1 AND empresa_id = $2 AND usuario_id = $3 AND habilitado`,
		f.A.Identidade, f.A.ID, f.A.UsuarioInterno).Scan(&n); err != nil {
		t.Fatalf("ler acesso_empresa: %v", err)
	}
	if n != 1 {
		t.Errorf("o trigger não criou o acesso da identidade do usuário interno de A (achou %d)", n)
	}

	// Nem anexando o trigger SECURITY DEFINER a uma tabela temporária.
	err := comoPapel(ctx, papelApp, f.internoA(), func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, `CREATE TEMP TABLE isca (identidade_id uuid, empresa_id uuid, id uuid,
		                           situacao text, excluido_em timestamptz)`); err != nil {
			return err
		}
		_, err := tx.Exec(ctx, `CREATE TRIGGER isca AFTER INSERT ON isca FOR EACH ROW
		                         EXECUTE FUNCTION plataforma.tg_acesso_empresa()`)
		return err
	})
	exigirErro(t, "tenaz_app pendura o trigger SECURITY DEFINER numa tabela temporária", err, erroPrivilegio)
}

// Uma credencial nunca muda de empresa, de usuário, de identidade nem de
// segredo: nem o tenaz_autenticacao, o único com UPDATE nas tabelas, tem
// UPDATE nessas colunas. É o que impede uma sessão de passar a servir outra
// empresa.
func TestCredencialNaoMudaDeDono(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)

	casos := []struct {
		tabela, coluna, id string
	}{
		{"plataforma.sessao", "empresa_id", f.A.SessaoWeb},
		{"plataforma.sessao", "usuario_id", f.A.SessaoWeb},
		{"plataforma.sessao", "identidade_id", f.A.SessaoWeb},
		{"plataforma.sessao", "token_hash", f.A.SessaoWeb},
		{"plataforma.sessao", "csrf_hash", f.A.SessaoWeb},
		{"plataforma.sessao", "canal", f.A.SessaoWeb},
		{"plataforma.codigo_autorizacao", "empresa_id", f.A.CodigoAutoriz},
		{"plataforma.codigo_autorizacao", "usuario_id", f.A.CodigoAutoriz},
		{"plataforma.codigo_autorizacao", "codigo_hash", f.A.CodigoAutoriz},
		{"plataforma.codigo_autorizacao", "uri_retorno", f.A.CodigoAutoriz},
		{"plataforma.codigo_autorizacao", "desafio_pkce", f.A.CodigoAutoriz},
		{"plataforma.token_renovacao", "familia_id", f.A.TokenRenovacao},
		{"plataforma.token_renovacao", "token_hash", f.A.TokenRenovacao},
		{"plataforma.token_email", "empresa_id", f.A.TokenEmail},
		{"plataforma.token_email", "usuario_id", f.A.TokenEmail},
		{"plataforma.token_email", "email", f.A.TokenEmail},
		{"plataforma.token_email", "finalidade", f.A.TokenEmail},
		{"plataforma.desafio_mfa", "identidade_id", f.A.DesafioMFA},
		{"plataforma.desafio_mfa", "token_hash", f.A.DesafioMFA},
		{"plataforma.identidade", "id", f.A.Identidade},
	}
	for _, papel := range []string{papelAutenticacao, papelAutenticacaoConexao, papelApp} {
		for _, c := range casos {
			t.Run(papel+"/"+c.tabela+"."+c.coluna, func(t *testing.T) {
				// SET coluna = coluna: o privilégio de coluna é conferido antes de
				// qualquer linha, e o valor não importa.
				col := pgx.Identifier{c.coluna}.Sanitize()
				sql := fmt.Sprintf(`UPDATE %s SET %s = %s WHERE id = $1`, c.tabela, col, col)
				err := comoPapel(ctx, papel, f.internoA(), func(tx pgx.Tx) error {
					_, err := tx.Exec(ctx, sql, c.id)
					return err
				})
				exigirErro(t, papel+" altera "+c.tabela+"."+c.coluna, err, erroPrivilegio)
			})
		}

		// Nem apaga: só a função de purga, como tenaz_job.
		for _, tabela := range []string{"plataforma.sessao", "plataforma.token_renovacao", "plataforma.codigo_autorizacao",
			"plataforma.desafio_mfa", "plataforma.token_email", "plataforma.identidade", "plataforma.acesso_empresa"} {
			err := comoPapel(ctx, papel, f.internoA(), func(tx pgx.Tx) error {
				_, err := tx.Exec(ctx, "DELETE FROM "+tabela+" WHERE false")
				return err
			})
			exigirErro(t, papel+" apaga em "+tabela, err, erroPrivilegio)
		}
	}

	// A prova de que o teste não passa por vacuidade: as colunas de ciclo de
	// vida o papel das credenciais altera.
	err := comoPapel(ctx, papelAutenticacao, contexto{}, func(tx pgx.Tx) error {
		n, err := alterou(ctx, tx, `UPDATE plataforma.sessao SET ultimo_uso_em = now() WHERE id = $1`, f.A.SessaoWeb)
		if err == nil && n != 1 {
			err = fmt.Errorf("alterou %d linha(s)", n)
		}
		return err
	})
	if err != nil {
		t.Errorf("tenaz_autenticacao não conseguiu atualizar ultimo_uso_em da sessão: %v; o privilégio de ciclo de vida sumiu", err)
	}
}

// A FK composta (empresa_id, usuario_id, identidade_id) só aceita a
// identidade que é de fato a do usuário naquela empresa.
func TestCredencialComIdentidadeAlheiaERecusada(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)

	casos := []struct {
		rotulo, sql string
		args        []any
	}{
		{"sessão de usuário de A com a identidade de B",
			`INSERT INTO plataforma.sessao (canal, identidade_id, empresa_id, usuario_id, expira_em)
			 VALUES ('app', $1, $2, $3, now() + interval '1 day')`, []any{f.B.Identidade, f.A.ID, f.A.UsuarioInterno}},
		{"sessão de usuário de B na empresa A",
			`INSERT INTO plataforma.sessao (canal, identidade_id, empresa_id, usuario_id, expira_em)
			 VALUES ('app', $1, $2, $3, now() + interval '1 day')`, []any{f.B.Identidade, f.A.ID, f.B.UsuarioInterno}},
		{"sessão de usuário sem identidade",
			`INSERT INTO plataforma.sessao (canal, identidade_id, empresa_id, usuario_id, expira_em)
			 VALUES ('app', $1, $2, $3, now() + interval '1 day')`, []any{f.A.Identidade, f.A.ID, f.A.UsuarioInativo}},
		{"código OAuth com a identidade do prestador para o interno",
			`INSERT INTO plataforma.codigo_autorizacao (codigo_hash, identidade_id, empresa_id, usuario_id,
			   cliente_oauth, uri_retorno, desafio_pkce, expira_em)
			 VALUES ($1, $2, $3, $4, 'tenaz-app', 'tenaz://retorno', 'abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQ',
			         now() + interval '60 seconds')`, []any{hashAleatorio(), f.IdentidadeF1, f.A.ID, f.A.UsuarioInterno}},
		{"convite de A para usuário de B",
			`INSERT INTO plataforma.token_email (token_hash, finalidade, email, empresa_id, usuario_id, expira_em)
			 VALUES ($1, 'convite', 'x@teste.tenaz', $2, $3, now() + interval '1 day')`,
			[]any{hashAleatorio(), f.A.ID, f.B.UsuarioInativo}},
	}
	for _, c := range casos {
		t.Run(c.rotulo, func(t *testing.T) {
			// Só o papel da autenticação insere credenciais (ADR-011); é nele
			// que a FK composta precisa recusar a identidade alheia.
			err := comoPapel(ctx, papelAutenticacao, contexto{}, func(tx pgx.Tx) error {
				_, err := tx.Exec(ctx, c.sql, c.args...)
				return err
			})
			exigirErro(t, c.rotulo, err, erroFK)
		})
	}
}
