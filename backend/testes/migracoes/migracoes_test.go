// Package migracoes testa as migrações que mexem em dado (documento 2, seção
// 14.2): aplica até N-1, grava dado no formato antigo COMO MIGRADOR (que passa
// pelo RLS forçado), aplica N e confere a conversão; depois o Down, e os
// caminhos de recusa (conferência que aborta a migração com violador).
//
// Cada cenário usa um banco novo dentro de um só contêiner descartável. Os
// papéis do bootstrap (00001) são do cluster e a migração os cria só se faltar,
// então vários bancos convivem.
//
// Toda migração de dados nova ganha aqui um teste de conversão e um de recusa.
package migracoes

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"tenaz/internal/plataforma/migracao"
)

// Mesma imagem da suíte de vazamento e do compose (ADR-012, decisão 2).
const imagemPostgres = "postgres:17.11-alpine3.24@sha256:e5d2edd28b1f8d40a0a6c0c9cdc545b19b56a0d5d88f94bc11a97d970d4700d0"

var (
	urlBase   string
	contadorB atomic.Int64
)

func TestMain(m *testing.M) {
	ctx := context.Background()
	cont, err := postgres.Run(ctx, imagemPostgres,
		postgres.WithDatabase("tenaz"),
		postgres.WithUsername("tenaz_migrador"),
		postgres.WithPassword("tenaz_teste"),
		postgres.BasicWaitStrategies(),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).WithStartupTimeout(90*time.Second)),
	)
	if err != nil {
		fmt.Fprintf(os.Stderr, "subir contêiner: %v\n", err)
		os.Exit(1)
	}
	urlBase, err = cont.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = testcontainers.TerminateContainer(cont)
		fmt.Fprintf(os.Stderr, "obter URL: %v\n", err)
		os.Exit(1)
	}
	codigo := m.Run()
	_ = testcontainers.TerminateContainer(cont)
	os.Exit(codigo)
}

// novoBanco cria um banco vazio e devolve a URL e uma conexão de migrador.
func novoBanco(t *testing.T) (string, *pgx.Conn) {
	t.Helper()
	ctx := context.Background()
	adm, err := pgx.Connect(ctx, urlBase)
	if err != nil {
		t.Fatalf("conectar: %v", err)
	}
	defer adm.Close(ctx)
	nome := fmt.Sprintf("mig_%d", contadorB.Add(1))
	if _, err := adm.Exec(ctx, "CREATE DATABASE "+nome); err != nil {
		t.Fatalf("criar banco: %v", err)
	}
	u, err := url.Parse(urlBase)
	if err != nil {
		t.Fatal(err)
	}
	u.Path = "/" + nome
	c, err := pgx.Connect(ctx, u.String())
	if err != nil {
		t.Fatalf("conectar ao banco novo: %v", err)
	}
	t.Cleanup(func() { c.Close(context.Background()) })
	return u.String(), c
}

func subirAte(t *testing.T, u string, v int64) {
	t.Helper()
	if err := migracao.AplicarAte(context.Background(), u, v); err != nil {
		t.Fatalf("aplicar até %d: %v", v, err)
	}
}

func descerAte(t *testing.T, u string, v int64) {
	t.Helper()
	if err := migracao.ReverterAte(context.Background(), u, v); err != nil {
		t.Fatalf("reverter até %d: %v", v, err)
	}
}

func exige(t *testing.T, c *pgx.Conn, sql string, args ...any) {
	t.Helper()
	if _, err := c.Exec(context.Background(), sql, args...); err != nil {
		t.Fatalf("%s: %v", sql, err)
	}
}

func um[T any](t *testing.T, c *pgx.Conn, sql string, args ...any) T {
	t.Helper()
	var v T
	if err := c.QueryRow(context.Background(), sql, args...).Scan(&v); err != nil {
		t.Fatalf("%s: %v", sql, err)
	}
	return v
}

func versao(t *testing.T, u string) int64 {
	t.Helper()
	v, err := migracao.Versao(context.Background(), u)
	if err != nil {
		t.Fatal(err)
	}
	return v
}

// exigeRecusa confere que a migração abortou com o SQLSTATE esperado e que o
// banco ficou na versão anterior.
func exigeRecusa(t *testing.T, u string, destino int64, codigo string) {
	t.Helper()
	err := migracao.AplicarAte(context.Background(), u, destino)
	if err == nil {
		t.Fatalf("a migração %d devia ser recusada (%s)", destino, codigo)
	}
	var pg *pgconn.PgError
	if errors.As(err, &pg) {
		if pg.Code != codigo {
			t.Fatalf("SQLSTATE %s, esperado %s: %v", pg.Code, codigo, err)
		}
	} else if !strings.Contains(err.Error(), codigo) {
		t.Fatalf("erro sem o código %s: %v", codigo, err)
	}
	if v := versao(t, u); v != destino-1 {
		t.Fatalf("versão após a recusa = %d, esperado %d", v, destino-1)
	}
}

// empresaAntiga insere uma empresa no formato da 00004 (com cnpj).
func empresaAntiga(t *testing.T, c *pgx.Conn, cnpj string) string {
	t.Helper()
	return um[string](t, c, `INSERT INTO app.empresa (razao_social, cnpj) VALUES ('Empresa '||$1, $1) RETURNING id::text`, cnpj)
}

func TestMigracao00009_IdentificacaoEPais(t *testing.T) {
	u, c := novoBanco(t)
	subirAte(t, u, 8)

	emp := empresaAntiga(t, c, "11222333000181")
	exige(t, c, `INSERT INTO app.unidade (empresa_id, codigo, nome, cnpj) VALUES ($1, 'U1', 'Matriz', '11222333000262')`, emp)
	exige(t, c, `INSERT INTO app.unidade (empresa_id, codigo, nome) VALUES ($1, 'U2', 'Sem CNPJ')`, emp)
	exige(t, c, `INSERT INTO app.pessoa (empresa_id, nome, cpf) VALUES ($1, 'Ana', '12345678909')`, emp)
	exige(t, c, `INSERT INTO app.pessoa (empresa_id, nome) VALUES ($1, 'Sem CPF')`, emp)
	exige(t, c, `INSERT INTO app.usuario (empresa_id, nome, custo_hora) VALUES ($1, 'Beto', 55.5)`, emp)

	subirAte(t, u, 9)

	if got := um[string](t, c, `SELECT pais||moeda||idioma_padrao||identificacao_tipo||identificacao FROM app.empresa`); got != "BRBRLpt-BRcnpj11222333000181" {
		t.Errorf("empresa convertida = %q", got)
	}
	if n := um[int](t, c, `SELECT count(*) FROM app.unidade WHERE identificacao_tipo = 'cnpj' AND identificacao = '11222333000262'`); n != 1 {
		t.Errorf("unidade com CNPJ convertida: %d", n)
	}
	if n := um[int](t, c, `SELECT count(*) FROM app.unidade WHERE identificacao IS NULL AND identificacao_tipo IS NULL`); n != 1 {
		t.Errorf("unidade sem CNPJ deve ficar sem identificação: %d", n)
	}
	if n := um[int](t, c, `SELECT count(*) FROM app.pessoa WHERE identificacao_tipo = 'cpf' AND identificacao = '12345678909'`); n != 1 {
		t.Errorf("pessoa convertida: %d", n)
	}
	if got := um[string](t, c, `SELECT idioma||moeda FROM app.usuario`); got != "pt-BRBRL" {
		t.Errorf("usuario = %q", got)
	}

	// Down: volta ao formato antigo sem perder o dado.
	descerAte(t, u, 8)
	if got := um[string](t, c, `SELECT cnpj FROM app.empresa`); got != "11222333000181" {
		t.Errorf("Down empresa.cnpj = %q", got)
	}
	if n := um[int](t, c, `SELECT count(*) FROM app.unidade WHERE cnpj = '11222333000262'`); n != 1 {
		t.Errorf("Down unidade.cnpj: %d", n)
	}
	if n := um[int](t, c, `SELECT count(*) FROM app.pessoa WHERE cpf = '12345678909'`); n != 1 {
		t.Errorf("Down pessoa.cpf: %d", n)
	}

	// Up de novo com dado, e a recusa do Down com empresa de outro país.
	subirAte(t, u, 9)
	exige(t, c, `INSERT INTO app.empresa (razao_social, pais, moeda, idioma_padrao, identificacao_tipo, identificacao)
	             VALUES ('Argentina', 'AR', 'ARS', 'es-419', 'cuit', '20123456786')`)
	if err := migracao.ReverterAte(context.Background(), u, 8); err == nil {
		t.Fatal("o Down da 00009 devia ser recusado com empresa argentina")
	}
	if v := versao(t, u); v != 9 {
		t.Errorf("versão após o Down recusado = %d", v)
	}
}

func TestMigracao00010_MoedaDoUsuario(t *testing.T) {
	u, c := novoBanco(t)
	subirAte(t, u, 9)

	br := empresaAntiga9(t, c, "BR", "BRL", "cnpj", "11222333000181")
	exige(t, c, `INSERT INTO app.usuario (empresa_id, nome, custo_hora, moeda) VALUES ($1, 'Beto', 55.55, 'BRL')`, br)

	subirAte(t, u, 10)
	if got := um[string](t, c, `SELECT custo_hora::text FROM app.usuario`); got != "55.550000" {
		t.Errorf("custo_hora = %q (esperado 55.550000)", got)
	}
	// A FK composta barra moeda diferente da da empresa.
	_, err := c.Exec(context.Background(), `INSERT INTO app.usuario (empresa_id, nome, moeda) VALUES ($1, 'Errado', 'ARS')`, br)
	var pg *pgconn.PgError
	if !errors.As(err, &pg) || pg.Code != "23503" {
		t.Errorf("moeda diferente da empresa devia dar 23503, deu %v", err)
	}

	descerAte(t, u, 9)
	if got := um[string](t, c, `SELECT custo_hora::text FROM app.usuario`); got != "55.55" {
		t.Errorf("Down custo_hora = %q", got)
	}
	subirAte(t, u, 10)
}

func TestMigracao00010_RecusaUsuarioEmMoedaDiferente(t *testing.T) {
	u, c := novoBanco(t)
	subirAte(t, u, 9)
	// Empresa argentina com usuário que ficou em BRL (o DEFAULT da 00009): violador
	// que a validação da FK, sob RLS forçado, não enxergaria; a 00010 confere antes.
	ar := empresaAntiga9(t, c, "AR", "ARS", "cuit", "20123456786")
	exige(t, c, `INSERT INTO app.usuario (empresa_id, nome) VALUES ($1, 'Violador')`, ar)

	exigeRecusa(t, u, 10, "23503")

	// Corrigido o dado, a migração passa.
	exige(t, c, `UPDATE app.usuario SET moeda = 'ARS'`)
	subirAte(t, u, 10)
}

func empresaAntiga9(t *testing.T, c *pgx.Conn, pais, moeda, tipo, ident string) string {
	t.Helper()
	idioma := "pt-BR"
	if pais != "BR" {
		idioma = "es-419"
	}
	return um[string](t, c, `INSERT INTO app.empresa (razao_social, pais, moeda, idioma_padrao, identificacao_tipo, identificacao)
	  VALUES ('Empresa '||$1, $1, $2, $3, $4, $5) RETURNING id::text`, pais, moeda, idioma, tipo, ident)
}

func TestMigracao00011_EnderecoSubdivisaoCodigoPostal(t *testing.T) {
	u, c := novoBanco(t)
	subirAte(t, u, 10)

	br := empresaAntiga9(t, c, "BR", "BRL", "cnpj", "11222333000181")
	unidade := func(cod, endereco string) {
		exige(t, c, `INSERT INTO app.unidade (empresa_id, codigo, nome, endereco) VALUES ($1, $2, $2, $3::jsonb)`, br, cod, nullSeVazio(endereco))
	}
	unidade("A", `{"uf":"sp","cep":"13010-111","municipio":"Campinas","logradouro":"Rua A"}`)
	unidade("B", `{"uf":"RJ"}`)
	unidade("C", `{"cep":"abc","municipio":"Sem dígito"}`)
	unidade("D", `{"municipio":"Sem uf nem cep"}`)
	unidade("E", ``)

	subirAte(t, u, 11)

	fim := func(cod string) string {
		return um[string](t, c, `SELECT coalesce(endereco::text, 'NULL') FROM app.unidade WHERE codigo = $1`, cod)
	}
	casos := map[string]string{
		"A": `{"municipio": "Campinas", "logradouro": "Rua A", "subdivisao": "BR-SP", "codigo_postal": "13010111"}`,
		"B": `{"subdivisao": "BR-RJ"}`,
		"C": `{"municipio": "Sem dígito"}`,
		"D": `{"municipio": "Sem uf nem cep"}`,
		"E": `NULL`,
	}
	for cod, esperado := range casos {
		if got := fim(cod); got != esperado {
			t.Errorf("unidade %s após o Up = %s, esperado %s", cod, got, esperado)
		}
	}
	if n := um[int](t, c, `SELECT count(*) FROM app.unidade WHERE endereco ? 'uf' OR endereco ? 'cep' OR endereco ? 'pais'`); n != 0 {
		t.Errorf("sobraram chaves antigas ou pais: %d", n)
	}

	// Empresa argentina já na forma nova: o Down devolve o sufixo e o código postal.
	ar := empresaAntiga9(t, c, "AR", "ARS", "cuit", "20123456786")
	exige(t, c, `INSERT INTO app.unidade (empresa_id, codigo, nome, endereco) VALUES ($1, 'AR1', 'AR1', '{"subdivisao":"AR-C","codigo_postal":"C1425DKF"}')`, ar)

	descerAte(t, u, 10)
	if got := fim("A"); got != `{"uf": "SP", "cep": "13010111", "municipio": "Campinas", "logradouro": "Rua A"}` {
		t.Errorf("Down A = %s", got)
	}
	if got := fim("AR1"); got != `{"uf": "C", "cep": "C1425DKF"}` {
		t.Errorf("Down AR1 = %s", got)
	}

	// Up de novo sobre o dado devolvido (a AR com uf e cep é recusada, ver abaixo).
	exige(t, c, `DELETE FROM app.unidade WHERE codigo = 'AR1'`)
	subirAte(t, u, 11)
	if got := fim("A"); got != casos["A"] {
		t.Errorf("segundo Up A = %s", got)
	}
}

func TestMigracao00011_RecusaEnderecoComUfEmEmpresaNaoBrasileira(t *testing.T) {
	for _, tc := range []struct{ pais, moeda, tipo, ident, endereco string }{
		{"AR", "ARS", "cuit", "20123456786", `{"uf":"C"}`},
		{"PY", "PYG", "ruc", "80012345", `{"cep":"1234"}`},
	} {
		t.Run(tc.pais, func(t *testing.T) {
			u, c := novoBanco(t)
			subirAte(t, u, 10)
			emp := empresaAntiga9(t, c, tc.pais, tc.moeda, tc.tipo, tc.ident)
			exige(t, c, `INSERT INTO app.unidade (empresa_id, codigo, nome, endereco) VALUES ($1, 'X', 'X', $2::jsonb)`, emp, tc.endereco)

			exigeRecusa(t, u, 11, "23514")

			// Corrigido, passa.
			exige(t, c, `UPDATE app.unidade SET endereco = '{"municipio":"Cidade"}'`)
			subirAte(t, u, 11)
		})
	}
}

func nullSeVazio(s string) any {
	if s == "" {
		return nil
	}
	return s
}
