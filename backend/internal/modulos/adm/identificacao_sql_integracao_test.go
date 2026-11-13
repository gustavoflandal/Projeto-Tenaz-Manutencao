package adm_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"strings"
	"testing"

	"tenaz/internal/plataforma/pais"
)

// A função app.identificacao_valida (migração 00009) repete em SQL a forma de
// cada tipo de identificação, que o Go e o TypeScript conferem pelos casos
// compartilhados. Este teste roda os mesmos casos contra a função e contra o
// Go, e exige que digam o mesmo sobre a forma.
//
// O que cada lado promete: o Go confere a forma e, só para cnpj e cpf, o dígito
// verificador; a função SQL confere só a forma, sobre o valor já normalizado
// (maiúsculas, só letras e dígitos). Logo, para um caso do arquivo:
//
//	SQL aceita  <=>  o Go aceita ou recusa só pelo dígito verificador.
const arquivoCasosIdentificacao = "../../../testes/casos/identificacao.json"

type casoIdentificacao struct {
	Nome, Tipo, Entrada, Motivo string
	Valido                      bool
}

func carregarCasosDeIdentificacao(t *testing.T) []casoIdentificacao {
	t.Helper()
	b, err := os.ReadFile(arquivoCasosIdentificacao)
	if err != nil {
		t.Fatalf("ler casos: %v", err)
	}
	var c struct {
		Versao    int
		Validacao []casoIdentificacao
	}
	if err := json.Unmarshal(b, &c); err != nil {
		t.Fatalf("casos inválidos: %v", err)
	}
	if c.Versao != 3 || len(c.Validacao) != 34 {
		t.Fatalf("casos de identificação: versão %d com %d casos de validação, esperados versão 3 e 34", c.Versao, len(c.Validacao))
	}
	return c.Validacao
}

// formaNoSQL pergunta à função do banco, sobre o valor normalizado.
func formaNoSQL(t *testing.T, tipo, entrada string) bool {
	t.Helper()
	var ok bool
	if err := amb.Migrador.QueryRow(context.Background(), `SELECT app.identificacao_valida($1, $2)`, tipo, pais.Normalizar(entrada)).Scan(&ok); err != nil {
		t.Fatalf("app.identificacao_valida(%q, %q): %v", tipo, entrada, err)
	}
	return ok
}

// formaNoGo diz se o Go aceita a forma: valido, ou recusado só pelo dígito
// verificador (a forma passou).
func formaNoGo(tipo, entrada string) bool {
	_, err := pais.ValidarIdentificacao(tipo, entrada)
	return err == nil || errors.Is(err, pais.ErrDigitoVerificador)
}

func TestIdentificacaoSQLSegueOsCasosCompartilhados(t *testing.T) {
	var divergencias []string
	for _, c := range carregarCasosDeIdentificacao(t) {
		esperado := c.Valido || c.Motivo == "digito_verificador"
		sql, gol := formaNoSQL(t, c.Tipo, c.Entrada), formaNoGo(c.Tipo, c.Entrada)
		if gol != esperado {
			t.Errorf("%s: o Go diverge do arquivo de casos (forma %v, esperada %v)", c.Nome, gol, esperado)
		}
		if sql != gol {
			divergencias = append(divergencias, c.Nome+": SQL="+map[bool]string{true: "aceita", false: "recusa"}[sql]+", Go="+map[bool]string{true: "aceita", false: "recusa"}[gol])
		}
	}
	if len(divergencias) > 0 {
		t.Errorf("a função app.identificacao_valida diverge do Go em %d caso(s):\n  %s", len(divergencias), strings.Join(divergencias, "\n  "))
	}
}

// Além dos casos do arquivo, as fronteiras de tamanho e de alfabeto de cada
// tipo: dígitos de 0 a 15 posições, com uma letra no começo, no meio e no fim.
// É onde duas cópias da mesma regra costumam divergir.
func TestIdentificacaoSQLSegueOGoNasFronteiras(t *testing.T) {
	var divergencias []string
	n := 0
	for _, tipo := range []string{"cnpj", "cpf", "cuit", "cuil", "dni", "ruc", "ci", "passaporte"} {
		for tam := 0; tam <= 15; tam++ {
			base := strings.Repeat("1", tam)
			variantes := []string{base}
			if tam > 0 {
				variantes = append(variantes, "A"+base[1:], base[:tam-1]+"A", base[:tam/2]+"A"+base[tam/2+1:])
			}
			for _, v := range variantes {
				n++
				if sql, gol := formaNoSQL(t, tipo, v), formaNoGo(tipo, v); sql != gol {
					divergencias = append(divergencias, tipo+" "+v)
				}
			}
		}
	}
	if len(divergencias) > 0 {
		t.Errorf("SQL e Go divergem em %d de %d valores de fronteira:\n  %s", len(divergencias), n, strings.Join(divergencias, "\n  "))
	}
}
