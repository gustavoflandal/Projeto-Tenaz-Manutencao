package vazamento

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
)

// Escopo externo dentro da empresa A (R4, RNF-037, RNF-038; seção 2.9 do
// documento 2): prestador de F1, locatário de L1 e link de uso único.

type escopoDeTeste struct {
	nome string
	ctx  contexto
}

func escoposExternos(f *fixturesE01) []escopoDeTeste {
	return []escopoDeTeste{
		{"prestador_F1", f.prestadorF1()},
		{"locatario_L1", f.locatarioL1()},
		{"link", f.linkA()},
	}
}

func lerAcessoExternoDeclarado(t *testing.T) map[string]map[string]bool {
	t.Helper()
	l, err := lerLista(arquivoAcessoExterno)
	if err != nil {
		t.Fatalf("ler %s: %v", arquivoAcessoExterno, err)
	}
	saida := map[string]map[string]bool{}
	for _, e := range l {
		saida[e.Nome] = map[string]bool{}
		for _, cmd := range strings.Split(strings.Join(e.Campos, ","), ",") {
			saida[e.Nome][cmd] = true
		}
	}
	return saida
}

// Orientado a tabela: em toda tabela, o comando que não está na lista de
// acesso externo não devolve nem altera nada para nenhum escopo externo. Tabela
// nova entra sozinha, com a regra "só interno".
func TestEscopoExternoForaDaListaNaoVeNemGrava(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	lista := lerAcessoExternoDeclarado(t)

	for _, tb := range tabelasDeEmpresa(t, ctx) {
		a := amostrar(t, ctx, tb, f)
		if a.LinhasA == 0 {
			continue // TestFixturesPovoamTodaTabela acusa
		}
		col := pgx.Identifier{tb.Coluna}.Sanitize()
		permitido := lista[tb.Nome]

		for _, e := range escoposExternos(f) {
			t.Run(tb.Nome+"/"+e.nome, func(t *testing.T) {
				executar := func(sql string, args ...any) (int64, error) {
					var n int64
					err := comoPapel(ctx, papelApp, e.ctx, func(tx pgx.Tx) error {
						if _, err := tx.Exec(ctx, "SELECT app.conferir_escopo()"); err != nil {
							return fmt.Errorf("conferir_escopo recusou o contexto da fixture: %w", err)
						}
						var err error
						if strings.HasPrefix(sql, "SELECT") {
							n, err = contar(ctx, tx, sql, args...)
						} else {
							n, err = alterou(ctx, tx, sql, args...)
						}
						return err
					})
					return n, err
				}

				if !permitido["SELECT"] {
					n, err := executar("SELECT count(*) FROM " + tb.Nome)
					if err != nil {
						t.Errorf("ler %s: %v", tb.Nome, err)
					} else if n != 0 {
						t.Errorf("%s leu %d linha(s) de %s, que é só interna", e.nome, n, tb.Nome)
					}
				}
				if !permitido["INSERT"] {
					// Cópia de uma linha de A com id novo: o que barra é o escopo.
					linha := a.LinhaA
					if tb.TemID && tb.Coluna != "id" {
						linha = comIDNovo(t, linha)
					}
					_, err := executar("INSERT INTO "+tb.Nome+
						" SELECT * FROM jsonb_populate_record(NULL::"+tb.Nome+", $1::jsonb)", string(linha))
					exigirRLS(t, e.nome+": INSERT em "+tb.Nome, err)
				}
				if !permitido["UPDATE"] {
					n, err := executar("UPDATE " + tb.Nome + " SET " + col + " = " + col)
					naoMudouNada(t, e.nome+": UPDATE em "+tb.Nome, n, err)
				}
				if !permitido["DELETE"] {
					n, err := executar("DELETE FROM " + tb.Nome)
					naoMudouNada(t, e.nome+": DELETE em "+tb.Nome, n, err)
				}
			})
		}
	}
}

func comIDNovo(t *testing.T, linha []byte) []byte {
	t.Helper()
	var m map[string]any
	if err := json.Unmarshal(linha, &m); err != nil {
		t.Fatalf("ler linha de amostra: %v", err)
	}
	m["id"] = novoID()
	b, err := json.Marshal(m)
	if err != nil {
		t.Fatalf("montar linha: %v", err)
	}
	return b
}

// idsVisiveis devolve os ids que o escopo enxerga numa tabela, ordenados.
func idsVisiveis(t *testing.T, ctx context.Context, c contexto, tabela string) []string {
	t.Helper()
	var ids []string
	err := comoPapel(ctx, papelApp, c, func(tx pgx.Tx) error {
		if _, err := tx.Exec(ctx, "SELECT app.conferir_escopo()"); err != nil {
			return err
		}
		linhas, err := tx.Query(ctx, "SELECT id::text FROM "+tabela+" ORDER BY 1")
		if err != nil {
			return err
		}
		ids, err = pgx.CollectRows(linhas, pgx.RowTo[string])
		return err
	})
	if err != nil {
		t.Fatalf("ler %s: %v", tabela, err)
	}
	return ids
}

func conjunto(ids ...string) []string {
	s := append([]string(nil), ids...)
	sort.Strings(s)
	return s
}

func mesmoConjunto(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Linhas exatas que cada escopo enxerga nas tabelas da lista de acesso
// externo. "Só as dele" e "todas as dele": ver de menos também é defeito,
// porque empurra o desenvolvedor a contornar a policy.
func TestEscopoExternoVeExatamenteOQueLheCabe(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)

	casos := []struct {
		escopo   string
		c        contexto
		tabela   string
		esperado []string
	}{
		{"prestador_F1", f.prestadorF1(), "app.usuario", conjunto(f.UsuarioF1, f.UsuarioF1Colega)},
		{"prestador_F1", f.prestadorF1(), "app.pessoa", conjunto(f.PessoaUsuarioF1, f.PessoaF1SemLogin)},
		{"prestador_F1", f.prestadorF1(), "app.unidade", conjunto(f.A.Unidade)},
		{"prestador_F1", f.prestadorF1(), "app.empresa", conjunto(f.A.ID)},
		{"locatario_L1", f.locatarioL1(), "app.usuario", conjunto(f.UsuarioL1, f.UsuarioL1Colega)},
		{"locatario_L1", f.locatarioL1(), "app.pessoa", conjunto(f.PessoaUsuarioL1)},
		{"locatario_L1", f.locatarioL1(), "app.unidade", conjunto(f.A.Unidade)},
		{"locatario_L1", f.locatarioL1(), "app.empresa", conjunto(f.A.ID)},
		{"link", f.linkA(), "app.usuario", nil},
		{"link", f.linkA(), "app.pessoa", nil},
		{"link", f.linkA(), "app.unidade", conjunto(f.A.Unidade)},
		{"link", f.linkA(), "app.empresa", conjunto(f.A.ID)},
	}
	for _, c := range casos {
		t.Run(c.escopo+"/"+c.tabela, func(t *testing.T) {
			got := idsVisiveis(t, ctx, c.c, c.tabela)
			if !mesmoConjunto(got, c.esperado) {
				t.Errorf("%s em %s viu %v; esperado %v", c.escopo, c.tabela, got, c.esperado)
			}
		})
	}
}

// Escrita do prestador em app.usuario: só a própria linha, sem trocar o
// vínculo. É a única escrita externa em tabela de E01 (fora a auditoria).
func TestPrestadorEscreveSoAPropriaLinhaSemTrocarVinculo(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)

	exec := func(sql string, args ...any) (int64, error) {
		var n int64
		err := comoPapel(ctx, papelApp, f.prestadorF1(), func(tx pgx.Tx) error {
			var err error
			n, err = alterou(ctx, tx, sql, args...)
			return err
		})
		return n, err
	}

	n, err := exec(`UPDATE app.usuario SET nome = 'Novo nome' WHERE id = $1`, f.UsuarioF1)
	if err != nil || n != 1 {
		t.Errorf("o prestador não conseguiu alterar o próprio nome (n=%d, err=%v); a policy corta demais", n, err)
	}

	_, err = exec(`UPDATE app.usuario SET fornecedor_id = $2 WHERE id = $1`, f.UsuarioF1, f.F2)
	exigirErro(t, "prestador troca o próprio fornecedor para F2", err, erroPrivilegio)

	_, err = exec(`UPDATE app.usuario SET fornecedor_id = NULL WHERE id = $1`, f.UsuarioF1)
	exigirErro(t, "prestador apaga o próprio fornecedor e vira interno", err, erroPrivilegio)

	n, err = exec(`UPDATE app.usuario SET nome = 'x' WHERE id = $1`, f.UsuarioF1Colega)
	naoMudouNada(t, "prestador altera o colega", n, err)

	n, err = exec(`UPDATE app.usuario SET nome = 'x' WHERE id = $1`, f.UsuarioF2)
	naoMudouNada(t, "prestador altera usuário de F2", n, err)

	n, err = exec(`DELETE FROM app.usuario WHERE id = $1`, f.UsuarioF1)
	naoMudouNada(t, "prestador apaga a si mesmo", n, err)

	_, err = exec(`INSERT INTO app.usuario (empresa_id, nome, fornecedor_id) VALUES ($1, 'Novo', $2)`, f.A.ID, f.F1)
	exigirRLS(t, "prestador cria usuário no próprio fornecedor", err)

	// Achado ISO-E01-01: as colunas que não são vínculo também não são do
	// prestador. Com pessoa_id, ele passa a enxergar a pessoa interna (CPF).
	t.Run("pessoa_id", func(t *testing.T) {
		var leu string
		err := comoPapel(ctx, papelApp, f.prestadorF1(), func(tx pgx.Tx) error {
			if _, err := tx.Exec(ctx, `UPDATE app.usuario SET pessoa_id = $2 WHERE id = $1`, f.UsuarioF1, f.A.PessoaInterna); err != nil {
				return err
			}
			return tx.QueryRow(ctx, `SELECT coalesce(max(cpf), '') FROM app.pessoa WHERE id = $1`, f.A.PessoaInterna).Scan(&leu)
		})
		if err != nil && sqlstate(err) != erroPrivilegio {
			t.Fatalf("erro inesperado: %v", err)
		}
		confirmarAchado(t, "ISO-E01-01", err == nil,
			fmt.Sprintf("o prestador ligou a si a pessoa interna e leu o CPF %q", leu))
	})
	t.Run("custo_hora", func(t *testing.T) {
		n, err := exec(`UPDATE app.usuario SET custo_hora = 999 WHERE id = $1`, f.UsuarioF1)
		if err != nil && sqlstate(err) != erroPrivilegio {
			t.Fatalf("erro inesperado: %v", err)
		}
		confirmarAchado(t, "ISO-E01-01", err == nil && n == 1, "o prestador alterou o próprio custo_hora")
	})
	t.Run("locatario_pessoa_id", func(t *testing.T) {
		err := comoPapel(ctx, papelApp, f.locatarioL1(), func(tx pgx.Tx) error {
			_, err := tx.Exec(ctx, `UPDATE app.usuario SET pessoa_id = $2 WHERE id = $1`, f.UsuarioL1, f.A.PessoaInterna)
			return err
		})
		if err != nil && sqlstate(err) != erroPrivilegio {
			t.Fatalf("erro inesperado: %v", err)
		}
		confirmarAchado(t, "ISO-E01-01", err == nil, "o locatário ligou a si a pessoa interna")
	})
}

// Regra anti-escalada (RF-ADM-005; seção 2.9.4): usuário externo nunca recebe
// perfil interno nem permissão vedada; interno não recebe perfil externo; e o
// vínculo de quem já tem perfil não muda. Quem tenta aqui é o administrador
// interno — o banco precisa recusar mesmo assim.
func TestUsuarioExternoNaoEscala(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)

	casos := []struct {
		rotulo, sql string
		args        []any
		codigo      string
	}{
		{"perfil derivado de prestador com os.encerrar",
			`INSERT INTO app.perfil (empresa_id, codigo, nome, perfil_origem, permissoes)
			 VALUES ($1, 'prest_x', 'x', 'prestador', '{os.encerrar}')`, []any{f.A.ID}, erroCheck},
		{"perfil derivado de locatário com permissão de Admin",
			`INSERT INTO app.perfil (empresa_id, codigo, nome, perfil_origem, permissoes)
			 VALUES ($1, 'loc_x', 'x', 'locatario', '{adm.gerenciar_usuarios}')`, []any{f.A.ID}, erroCheck},
		{"perfil padrão de prestador com os.encerrar",
			`UPDATE app.perfil SET permissoes = permissoes || '{os.encerrar}' WHERE id = $1`, []any{f.A.PerfilPrestador}, erroCheck},
		{"perfil gestor para o prestador",
			`INSERT INTO app.usuario_perfil (empresa_id, usuario_id, perfil_id) VALUES ($1, $2, $3)`,
			[]any{f.A.ID, f.UsuarioF1, f.A.PerfilGestor}, erroCheck},
		{"perfil prestador para o interno",
			`INSERT INTO app.usuario_perfil (empresa_id, usuario_id, perfil_id) VALUES ($1, $2, $3)`,
			[]any{f.A.ID, f.A.UsuarioInterno, f.A.PerfilPrestador}, erroCheck},
		{"perfil locatário para o prestador",
			`INSERT INTO app.usuario_perfil (empresa_id, usuario_id, perfil_id) VALUES ($1, $2, $3)`,
			[]any{f.A.ID, f.UsuarioF1, f.A.PerfilLocatario}, erroCheck},
		{"os.encerrar avulsa para o prestador",
			`INSERT INTO app.usuario_permissao (empresa_id, usuario_id, permissao) VALUES ($1, $2, 'os.encerrar')`,
			[]any{f.A.ID, f.UsuarioF1}, erroCheck},
		{"permissão de Admin avulsa para o locatário",
			`INSERT INTO app.usuario_permissao (empresa_id, usuario_id, permissao) VALUES ($1, $2, 'adm.gerenciar_usuarios')`,
			[]any{f.A.ID, f.UsuarioL1}, erroCheck},
		{"prestador com perfil vira interno",
			`UPDATE app.usuario SET fornecedor_id = NULL WHERE id = $1`, []any{f.UsuarioF1}, erroCheck},
		{"interno com perfil vira prestador",
			`UPDATE app.usuario SET fornecedor_id = $2 WHERE id = $1`, []any{f.A.UsuarioInterno, f.F1}, erroCheck},
		{"identidade do usuário trocada",
			`UPDATE app.usuario SET identidade_id = $2 WHERE id = $1`, []any{f.A.UsuarioInterno, f.IdentidadeF1}, erroPrivilegio},
	}
	for _, c := range casos {
		t.Run(c.rotulo, func(t *testing.T) {
			err := comoPapel(ctx, papelApp, f.internoA(), func(tx pgx.Tx) error {
				_, err := tx.Exec(ctx, c.sql, c.args...)
				return err
			})
			exigirErro(t, c.rotulo, err, c.codigo)
		})
	}

	// Achado ISO-E01-03: o tipo do perfil (perfil_origem, ou o código) muda
	// depois de atribuído, e o prestador passa a ter os.encerrar.
	t.Run("perfil atribuído ao prestador muda de tipo", func(t *testing.T) {
		var permissoes []string
		err := comoPapel(ctx, papelApp, f.internoA(), func(tx pgx.Tx) error {
			perfil := novoID()
			for _, p := range []struct {
				sql  string
				args []any
			}{
				{`INSERT INTO app.perfil (empresa_id, id, codigo, nome, perfil_origem, permissoes)
				  VALUES ($1, $2, 'prest_derivado', 'Prestador derivado', 'prestador', '{os.executar}')`, []any{f.A.ID, perfil}},
				{`INSERT INTO app.usuario_perfil (empresa_id, usuario_id, perfil_id) VALUES ($1, $2, $3)`,
					[]any{f.A.ID, f.UsuarioF1, perfil}},
				{`UPDATE app.perfil SET perfil_origem = NULL, permissoes = '{os.executar,os.encerrar}' WHERE id = $1`,
					[]any{perfil}},
			} {
				if _, err := tx.Exec(ctx, p.sql, p.args...); err != nil {
					return err
				}
			}
			return tx.QueryRow(ctx, `
				SELECT array_agg(DISTINCT x) FROM app.usuario_perfil up
				JOIN app.perfil p ON p.empresa_id = up.empresa_id AND p.id = up.perfil_id, unnest(p.permissoes) x
				WHERE up.usuario_id = $1 AND up.excluido_em IS NULL`, f.UsuarioF1).Scan(&permissoes)
		})
		if err != nil && sqlstate(err) != erroCheck && sqlstate(err) != erroPrivilegio {
			t.Fatalf("erro inesperado: %v", err)
		}
		confirmarAchado(t, "ISO-E01-03", err == nil && contem(permissoes, "os.encerrar"),
			fmt.Sprintf("o prestador ficou com as permissões %v", permissoes))
	})
}
