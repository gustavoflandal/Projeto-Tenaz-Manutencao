package vazamento

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

// País, identificação, moeda e idioma (00009, ADR-013). Três verificações:
//
//   - empresas de países diferentes não se alcançam: A (Brasil), B (Argentina)
//     e C (Paraguai), em toda tabela de empresa. A contra B, com escrita, já é
//     a matriz de empresas_test.go; aqui entram os seis pares;
//   - o usuário externo (prestador, locatário) altera na própria linha só
//     nome, telefone e idioma: coluna a coluna, gerado do catálogo, para que
//     coluna nova de app.usuario entre recusada (ISO-E01-01);
//   - o país e a moeda da empresa não mudam, nem pelo superusuário.

// Empresas de países diferentes: no contexto de cada uma, nenhuma linha das
// outras duas, e o total é exatamente o da própria.
func TestEmpresasDePaisesDiferentesNaoSeAlcancam(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	empresas := []empresaFixture{f.A, f.B, f.C}
	for _, tb := range tabelasDeEmpresa(t, ctx) {
		for _, x := range empresas {
			t.Run(tb.Nome+"/"+x.Pais, func(t *testing.T) {
				emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
					p, err := problemasEntrePaises(ctx, tx, tb, x, empresas)
					for _, m := range p {
						t.Error(m)
					}
					if err != nil {
						t.Errorf("ler %s no contexto da empresa de %s: %v", tb.Nome, x.Pais, err)
					}
				})
			})
		}
	}
}

// problemasEntrePaises: na tabela tb, no contexto da empresa x, o tenaz_app
// não lê linha de nenhuma outra empresa de empresas, e o total é exatamente o
// de x. Roda na transação tx (o superusuário conta as linhas próprias; o
// tenaz_app lê num savepoint desfeito), para que a prova do detector estrague
// o banco na mesma transação e nada da prova chegue a outra conexão.
func problemasEntrePaises(ctx context.Context, tx pgx.Tx, tb tabelaDeEmpresa, x empresaFixture, empresas []empresaFixture) ([]string, error) {
	col := pgx.Identifier{tb.Coluna}.Sanitize()
	var p []string
	proprias, err := contar(ctx, tx, `SELECT count(*) FROM `+tb.Nome+` WHERE `+col+` = $1`, x.ID)
	if err != nil {
		return p, err
	}
	if proprias == 0 {
		p = append(p, fmt.Sprintf("a empresa de %s não tem linha em %s nas fixtures; o teste passaria por vacuidade", x.Pais, tb.Nome))
	}
	sub, err := tx.Begin(ctx)
	if err != nil {
		return p, err
	}
	defer func() { _ = sub.Rollback(ctx) }()
	if err := trocarPapel(ctx, sub, papelApp, contexto{Empresa: x.ID, Escopo: "interno"}); err != nil {
		return p, err
	}
	for _, y := range empresas {
		if y.ID == x.ID {
			continue
		}
		n, err := contar(ctx, sub, `SELECT count(*) FROM `+tb.Nome+` WHERE `+col+` = $1`, y.ID)
		if err != nil {
			return p, err
		}
		if n != 0 {
			p = append(p, fmt.Sprintf("no contexto da empresa de %s, o tenaz_app leu %d linha(s) da empresa de %s em %s", x.Pais, n, y.Pais, tb.Nome))
		}
	}
	total, err := contar(ctx, sub, `SELECT count(*) FROM `+tb.Nome)
	if err != nil {
		return p, err
	}
	if total != proprias {
		p = append(p, fmt.Sprintf("no contexto da empresa de %s, o tenaz_app viu %d linha(s) em %s; ela tem %d", x.Pais, total, tb.Nome, proprias))
	}
	return p, nil
}

// Prova do detector entre países (pendência da etapa 014): com a RLS de
// app.unidade desligada, o teste acusa os seis pares ordenados (cada país
// lendo cada um dos outros dois) e o total de cada país; e não acusa nenhuma
// outra tabela, que continua defendida.
func TestDetectorAcusaPaisesQueSeAlcancam(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	empresas := []empresaFixture{f.A, f.B, f.C}
	const alvo = "app.unidade"
	tabelas := tabelasDeEmpresa(t, ctx)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		estragar(t, ctx, tx, "ALTER TABLE "+alvo+" DISABLE ROW LEVEL SECURITY")
		var pares int
		for _, tb := range tabelas {
			for _, x := range empresas {
				p, err := problemasEntrePaises(ctx, tx, tb, x, empresas)
				if err != nil {
					t.Fatalf("ler %s no contexto da empresa de %s: %v", tb.Nome, x.Pais, err)
				}
				if tb.Nome != alvo {
					if len(p) > 0 {
						t.Errorf("a prova desligou só a RLS de %s, e o detector acusou %s: %q", alvo, tb.Nome, p)
					}
					continue
				}
				for _, y := range empresas {
					if y.ID == x.ID {
						continue
					}
					inicio := fmt.Sprintf("no contexto da empresa de %s, o tenaz_app leu ", x.Pais)
					fim := fmt.Sprintf(" linha(s) da empresa de %s em %s", y.Pais, alvo)
					achou := false
					for _, m := range p {
						achou = achou || (strings.HasPrefix(m, inicio) && strings.HasSuffix(m, fim))
					}
					if !achou {
						t.Errorf("sem a RLS de %s, o detector não acusou a leitura de %s por %s: %q", alvo, y.Pais, x.Pais, p)
						continue
					}
					pares++
				}
				total := fmt.Sprintf("no contexto da empresa de %s, o tenaz_app viu ", x.Pais)
				if !strings.Contains(strings.Join(p, " | "), total) {
					t.Errorf("sem a RLS de %s, o detector não acusou o total de %s: %q", alvo, x.Pais, p)
				}
			}
		}
		if pares != 6 {
			t.Errorf("a prova acusou %d par(es) de países; esperava os seis", pares)
		}
	})
}

// As fixtures gravam a moeda do usuário igual à da empresa (A BRL, B ARS,
// C PYG). Até a 00010, o DEFAULT 'BRL' fazia os usuários de B e C nascerem em
// BRL sem que nada acusasse; sem o DEFAULT o INSERT esquecido falha, e este
// teste cobre o INSERT que informa a moeda errada.
func TestFixturesGravamUsuarioNaMoedaDaEmpresa(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	for _, e := range []empresaFixture{f.A, f.B, f.C} {
		var total, divergentes int64
		if err := pool.QueryRow(ctx, `
			SELECT count(*), count(*) FILTER (WHERE u.moeda IS DISTINCT FROM e.moeda)
			FROM app.usuario u JOIN app.empresa e ON e.id = u.empresa_id
			WHERE u.empresa_id = $1`, e.ID).Scan(&total, &divergentes); err != nil {
			t.Fatalf("contar usuários da empresa de %s: %v", e.Pais, err)
		}
		if total == 0 {
			t.Errorf("a empresa de %s não tem usuário nas fixtures; o teste passaria por vacuidade", e.Pais)
		}
		if divergentes != 0 {
			t.Errorf("%d usuário(s) da empresa de %s com moeda diferente de %s", divergentes, e.Pais, e.Moeda)
		}
	}
}

// app.usuario.moeda é a da empresa: a FK composta fk_usuario_moeda_empresa
// (empresa_id, moeda) -> app.empresa (id, moeda), da 00010, recusa com 23503
// um usuário da empresa brasileira em pesos argentinos. O custo do apontamento
// multiplica custo_hora por horas na moeda da empresa; moeda errada daria custo
// errado em ordem de grandeza. A moeda da empresa é imutável, então a FK não
// precisa de ON UPDATE. O controle exige que a moeda certa seja aceita, para o
// teste não passar por uma recusa que viesse de outra causa.
func TestBancoRecusaUsuarioEmMoedaDeOutraEmpresa(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	inserir := func(moeda string) error {
		return comoPapel(ctx, papelApp, f.internoA(), func(tx pgx.Tx) error {
			sub, err := tx.Begin(ctx)
			if err != nil {
				return err
			}
			defer func() { _ = sub.Rollback(ctx) }()
			_, err = sub.Exec(ctx, `INSERT INTO app.usuario (empresa_id, nome, moeda) VALUES ($1, 'Teste de moeda', $2)`,
				f.A.ID, moeda)
			return err
		})
	}
	if err := inserir(f.A.Moeda); err != nil {
		t.Errorf("controle: usuário da empresa de %s em %s foi recusado (SQLSTATE %s: %v)", f.A.Pais, f.A.Moeda, sqlstate(err), err)
	}
	err := inserir(f.B.Moeda)
	if err == nil {
		t.Fatalf("o banco aceitou usuário da empresa de %s em %s", f.A.Pais, f.B.Moeda)
	}
	if sqlstate(err) != "23503" {
		t.Errorf("usuário da empresa de %s em %s: SQLSTATE %s, esperado 23503 (fk_usuario_moeda_empresa): %v",
			f.A.Pais, f.B.Moeda, sqlstate(err), err)
	}
}

// mudouComo executa o UPDATE ... RETURNING <valor> como papel, num savepoint
// desfeito, e diz se o valor gravado ficou diferente do anterior. Um UPDATE
// que um trigger devolve ao valor antigo (criado_em, pelo trg_padrao) não
// mudou nada.
func mudouComo(ctx context.Context, tx pgx.Tx, papel string, c contexto, antes *string, sql string, args ...any) (bool, error) {
	sub, err := tx.Begin(ctx)
	if err != nil {
		return false, err
	}
	defer func() { _ = sub.Rollback(ctx) }()
	if err := trocarPapel(ctx, sub, papel, c); err != nil {
		return false, erroPreparacao{fmt.Sprintf("preparar %s: %v", papel, err)}
	}
	var depois *string
	if err := sub.QueryRow(ctx, sql, args...).Scan(&depois); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	if antes == nil || depois == nil {
		return (antes == nil) != (depois == nil), nil
	}
	return *antes != *depois, nil
}

// alterouComo executa o comando como papel, num savepoint desfeito, e devolve
// as linhas alteradas.
func alterouComo(ctx context.Context, tx pgx.Tx, papel string, c contexto, sql string, args ...any) (int64, error) {
	sub, err := tx.Begin(ctx)
	if err != nil {
		return 0, err
	}
	defer func() { _ = sub.Rollback(ctx) }()
	if err := trocarPapel(ctx, sub, papel, c); err != nil {
		return 0, erroPreparacao{fmt.Sprintf("preparar %s: %v", papel, err)}
	}
	return alterou(ctx, sub, sql, args...)
}

// Colunas que o usuário externo altera na própria linha (00009: idioma), e as
// que o trg_padrao preenche sozinho.
var (
	colunasLivresDoExterno = map[string]string{
		"nome": "'Nome alterado pela suíte'", "telefone": "'11999990000'", "idioma": "'es-419'",
	}
	colunasDoSistema = map[string]bool{"versao": true, "alterado_xid": true, "atualizado_em": true, "atualizado_por": true}
	// valoresDeTeste: coluna com CHECK de lista fechada recebe um valor válido,
	// para que só a defesa do escopo possa recusar a troca.
	valoresDeTeste = map[string]string{"moeda": "CASE WHEN moeda = 'BRL' THEN 'ARS' ELSE 'BRL' END"}
	reCaractere    = regexp.MustCompile(`^character\((\d+)\)$`)
)

// valorDiferente devolve uma expressão SQL que muda o valor da coluna, pelo
// tipo; "" para tipo que a suíte ainda não sabe mudar.
func valorDiferente(coluna, tipo string) string {
	c := pgx.Identifier{coluna}.Sanitize()
	switch {
	case tipo == "text" || tipo == "citext" || strings.HasPrefix(tipo, "character varying"):
		return "coalesce(" + c + "::text, '') || 'x'"
	case reCaractere.MatchString(tipo):
		n := reCaractere.FindStringSubmatch(tipo)[1]
		return "CASE WHEN " + c + " = repeat('Z', " + n + ") THEN repeat('Y', " + n + ") ELSE repeat('Z', " + n + ") END"
	case tipo == "uuid":
		return "gen_random_uuid()"
	case tipo == "integer" || tipo == "bigint" || tipo == "smallint" || strings.HasPrefix(tipo, "numeric"):
		return "coalesce(" + c + ", 0) + 1"
	case tipo == "boolean":
		return "NOT coalesce(" + c + ", false)"
	case tipo == "timestamp with time zone":
		return "coalesce(" + c + ", now()) + interval '1 day'"
	case tipo == "date":
		return "coalesce(" + c + ", current_date) + 1"
	case tipo == "jsonb":
		return "coalesce(" + c + ", '{}'::jsonb) || '{\"vazamento\": 1}'::jsonb"
	case tipo == "text[]":
		return "coalesce(" + c + ", '{}'::text[]) || 'x'::text"
	}
	return ""
}

// problemasDasColunasDoExterno: o prestador F1 e o locatário L1, cada um na
// própria linha, mudam cada coluna de app.usuario. As livres mudam; toda outra
// é recusada com 42501 (o trigger tg_usuario_colunas_externo, ou outra
// defesa antes dele), e recusa por outro motivo não prova nada.
func problemasDasColunasDoExterno(t *testing.T, ctx context.Context, tx pgx.Tx, f *fixturesE01) []string {
	t.Helper()
	linhas, err := tx.Query(ctx, `SELECT attname, format_type(atttypid, atttypmod) FROM pg_attribute
	                               WHERE attrelid = 'app.usuario'::regclass AND attnum > 0 AND NOT attisdropped AND attgenerated = ''
	                               ORDER BY attnum`)
	if err != nil {
		t.Fatalf("colunas de app.usuario: %v", err)
	}
	type coluna struct{ nome, tipo string }
	var colunas []coluna
	for linhas.Next() {
		var c coluna
		if err := linhas.Scan(&c.nome, &c.tipo); err != nil {
			t.Fatal(err)
		}
		colunas = append(colunas, c)
	}
	linhas.Close()

	var p []string
	for _, e := range []struct {
		rotulo  string
		c       contexto
		usuario string
	}{{"prestador F1", f.prestadorF1(), f.UsuarioF1}, {"locatário L1", f.locatarioL1(), f.UsuarioL1}} {
		for _, c := range colunas {
			if colunasDoSistema[c.nome] {
				continue
			}
			valor, livre := colunasLivresDoExterno[c.nome]
			if v, ok := valoresDeTeste[c.nome]; ok && !livre {
				valor = v
			} else if !livre {
				if valor = valorDiferente(c.nome, c.tipo); valor == "" {
					p = append(p, fmt.Sprintf("coluna %s (%s) de app.usuario sem valor de teste: acrescente o tipo em valorDiferente", c.nome, c.tipo))
					continue
				}
			}
			col := pgx.Identifier{c.nome}.Sanitize()
			var antes *string
			if err := tx.QueryRow(ctx, `SELECT `+col+`::text FROM app.usuario WHERE id = $1`, e.usuario).Scan(&antes); err != nil {
				t.Fatalf("ler %s de %s: %v", c.nome, e.rotulo, err)
			}
			mudou, err := mudouComo(ctx, tx, papelApp, e.c, antes,
				`UPDATE app.usuario SET `+col+` = `+valor+` WHERE id = $1 RETURNING `+col+`::text`, e.usuario)
			var pgErr *pgconn.PgError
			switch {
			case livre && (err != nil || !mudou):
				p = append(p, fmt.Sprintf("%s não alterou a própria coluna livre %s (erro: %v)", e.rotulo, c.nome, err))
			case livre:
			case err == nil && mudou:
				p = append(p, fmt.Sprintf("%s alterou a própria coluna %s", e.rotulo, c.nome))
			case err == nil:
			case !errors.As(err, &pgErr) || pgErr.Code != erroPrivilegio:
				p = append(p, fmt.Sprintf("%s, coluna %s: recusado por outro motivo, que não prova a defesa: %v", e.rotulo, c.nome, err))
			}
		}
	}
	return p
}

func TestUsuarioExternoSoAlteraNomeTelefoneEIdioma(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		for _, p := range problemasDasColunasDoExterno(t, ctx, tx, f) {
			t.Errorf("%s (ISO-E01-01; 00009)", p)
		}
	})
}

// A prova libera o e-mail ao externo. A moeda não serve mais de prova: sem o
// trigger, a FK fk_usuario_moeda_empresa (00010) ainda a recusaria, com 23503.
func TestDetectorAcusaColunaLiberadaAoExterno(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		estragar(t, ctx, tx, `
			CREATE OR REPLACE FUNCTION app.tg_usuario_colunas_externo() RETURNS trigger LANGUAGE plpgsql AS $$
			DECLARE
			  v_livres text[] := ARRAY['nome', 'telefone', 'idioma', 'email',
			                           'versao', 'alterado_xid', 'atualizado_em', 'atualizado_por'];
			BEGIN
			  IF app.escopo_atual() IN ('prestador', 'locatario', 'link')
			     AND (to_jsonb(NEW) - v_livres) IS DISTINCT FROM (to_jsonb(OLD) - v_livres) THEN
			    RAISE EXCEPTION 'prova' USING ERRCODE = '42501';
			  END IF;
			  RETURN NEW;
			END
			$$`)
		p := strings.Join(problemasDasColunasDoExterno(t, ctx, tx, f), " | ")
		for _, esperado := range []string{"prestador F1 alterou a própria coluna email", "locatário L1 alterou a própria coluna email"} {
			if !strings.Contains(p, esperado) {
				t.Errorf("com o e-mail liberado ao externo, o detector não acusou %q: %q", esperado, p)
			}
		}
	})
}

// problemasDoPais: a empresa brasileira A não vira argentina nem troca de
// moeda, pelo tenaz_app nem pelo superusuário; e o idioma padrão, que não é
// travado, muda (controle positivo: a linha é alterável).
func problemasDoPais(t *testing.T, ctx context.Context, tx pgx.Tx, f *fixturesE01) []string {
	t.Helper()
	var p []string
	paraArgentina := `UPDATE app.empresa SET pais = 'AR', moeda = 'ARS', identificacao_tipo = 'cuit', identificacao = '` +
		documentoAleatorio("cuit") + `' WHERE id = $1`
	soMoeda := `UPDATE app.empresa SET moeda = 'ARS' WHERE id = $1`
	for _, papel := range []string{papelApp, ""} {
		rotulo := papel
		if rotulo == "" {
			rotulo = "superusuário"
		}
		for _, sql := range []string{paraArgentina, soMoeda} {
			n, err := alterouComo(ctx, tx, papel, f.internoA(), sql, f.A.ID)
			if err == nil || sqlstate(err) != erroPrivilegio {
				p = append(p, fmt.Sprintf("%s: %q não foi recusado pela trava do país (n=%d, erro: %v)", rotulo, sql, n, err))
			}
		}
	}
	if n, err := alterouComo(ctx, tx, papelApp, f.internoA(), `UPDATE app.empresa SET idioma_padrao = 'es-419' WHERE id = $1`, f.A.ID); err != nil || n != 1 {
		p = append(p, fmt.Sprintf("controle: o tenaz_app não alterou o idioma padrão de A (n=%d, erro: %v)", n, err))
	}
	return p
}

func TestPaisEMoedaDaEmpresaNaoMudam(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		for _, p := range problemasDoPais(t, ctx, tx, f) {
			t.Errorf("%s (ADR-013; 00009)", p)
		}
	})
}

func TestDetectorAcusaPaisMutavel(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		estragar(t, ctx, tx, "DROP TRIGGER trg_empresa_pais_imutavel ON app.empresa")
		p := strings.Join(problemasDoPais(t, ctx, tx, f), " | ")
		for _, papel := range []string{papelApp, "superusuário"} {
			if !strings.Contains(p, papel+`: "UPDATE app.empresa SET pais = 'AR'`) {
				t.Errorf("sem a trava, a troca de país por %s não foi acusada: %q", papel, p)
			}
		}
		if strings.Contains(p, "controle:") {
			t.Errorf("a prova quebrou o controle positivo: %q", p)
		}
	})
}
