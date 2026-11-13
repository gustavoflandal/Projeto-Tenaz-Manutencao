package dinheiro_test

import (
	"encoding/json"
	"errors"
	"math/big"
	"os"
	"strconv"
	"testing"

	"github.com/jackc/pgx/v5/pgtype"

	"tenaz/internal/plataforma/dinheiro"
	"tenaz/internal/plataforma/erros"
)

// Os casos são os mesmos que pacotes/dominio/src/dinheiro.ts executa: um
// arquivo, duas suítes (documento 5, Parte F, item 8).
const arquivoCasos = "../../../testes/casos/dinheiro.json"

type casos struct {
	Versao int
	Moedas []struct {
		Codigo string
		Casas  int
	}
	ConstrucaoMoeda []struct {
		Nome, Codigo, Erro string
		Valido             bool
		Casas              int
	} `json:"construcao_moeda"`
	Limites struct {
		ValorMaxUnidades         map[string]string `json:"valor_max_unidades"`
		PrecoUnitarioMaxUnidades string            `json:"preco_unitario_max_unidades"`
	}
	Leitura []struct {
		Nome, Especie, Moeda, Entrada string
		AceitaNegativo                bool `json:"aceita_negativo"`
		Valido                        bool
		Unidades, Canonico, Erro      string
	}
	TextoCanonico []struct {
		Nome, Especie, Moeda, Unidades, Esperado string
	} `json:"texto_canonico"`
	SomaSubtracao []struct {
		Nome, Operacao string
		A, B           struct{ Moeda, Valor string }
		Esperado, Erro string
	} `json:"soma_subtracao"`
	Multiplicacao []struct {
		Nome, Moeda, Preco, Quantidade, Esperado, Unidades, Erro string
	}
	Rateio []struct {
		Nome, Moeda, Valor, Erro string
		Pesos, Esperado          []string
	}
}

func carregar(t *testing.T) casos {
	t.Helper()
	b, err := os.ReadFile(arquivoCasos)
	if err != nil {
		t.Fatalf("ler casos: %v", err)
	}
	var c casos
	if err := json.Unmarshal(b, &c); err != nil {
		t.Fatalf("casos inválidos: %v", err)
	}
	return c
}

// A trava do arquivo de casos: se um grupo for renomeado no JSON, o Go não pode
// rodar zero casos e passar. Versão e tamanho de cada grupo são os que o
// TypeScript também espera (dinheiro.test.ts); mudar o arquivo exige mudar os dois.
func TestCasosTravados(t *testing.T) {
	c := carregar(t)
	if c.Versao != 2 {
		t.Fatalf("versão dos casos %d, esperada 2", c.Versao)
	}
	grupos := map[string][2]int{
		"moedas": {len(c.Moedas), 3}, "construcao_moeda": {len(c.ConstrucaoMoeda), 9}, "leitura": {len(c.Leitura), 121},
		"texto_canonico": {len(c.TextoCanonico), 33}, "soma_subtracao": {len(c.SomaSubtracao), 28},
		"multiplicacao": {len(c.Multiplicacao), 46}, "rateio": {len(c.Rateio), 29},
		"limites.valor_max_unidades": {len(c.Limites.ValorMaxUnidades), 3},
	}

	for nome, g := range grupos {
		if g[0] == 0 || g[0] != g[1] {
			t.Errorf("grupo %s: %d casos, esperados %d (e nenhum grupo pode ficar vazio)", nome, g[0], g[1])
		}
	}
	total := len(c.ConstrucaoMoeda) + len(c.Leitura) + len(c.TextoCanonico) + len(c.SomaSubtracao) + len(c.Multiplicacao) + len(c.Rateio)
	if total != 266 {
		t.Errorf("total de casos %d, esperado 266", total)
	}
	if c.Limites.PrecoUnitarioMaxUnidades == "" {
		t.Error("limites.preco_unitario_max_unidades ausente")
	}
}

func moeda(t *testing.T, cod string) dinheiro.Moeda {
	t.Helper()
	m, err := dinheiro.ParaMoeda(cod)
	if err != nil {
		t.Fatalf("moeda %s: %v", cod, err)
	}
	return m
}

func inteiro(t *testing.T, s string) int64 {
	t.Helper()
	n, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		t.Fatalf("inteiro %q: %v", s, err)
	}
	return n
}

// quantidade lê a notação do arquivo de casos: decimal ou fração a/b. A
// notação a/b é só daqui (código de teste); a API nunca a aceita.
func quantidade(t *testing.T, s string) *big.Rat {
	t.Helper()
	r, ok := new(big.Rat).SetString(s)
	if !ok {
		t.Fatalf("quantidade %q", s)
	}
	return r
}

func TestMoedasELimites(t *testing.T) {
	c := carregar(t)
	for _, m := range c.Moedas {
		if got := moeda(t, m.Codigo).Casas(); got != m.Casas {
			t.Errorf("%s: %d casas, esperado %d", m.Codigo, got, m.Casas)
		}
	}
	for cod, max := range c.Limites.ValorMaxUnidades {
		m, n := moeda(t, cod), inteiro(t, max)
		if _, err := dinheiro.NovoValor(n, m); err != nil {
			t.Errorf("%s: o máximo %d deveria caber: %v", cod, n, err)
		}
		if _, err := dinheiro.NovoValor(-n, m); err != nil {
			t.Errorf("%s: o mínimo %d deveria caber: %v", cod, -n, err)
		}
		if _, err := dinheiro.NovoValor(n+1, m); !errors.Is(err, dinheiro.ErrForaDaFaixa) {
			t.Errorf("%s: %d deveria estar fora da faixa: %v", cod, n+1, err)
		}
	}
	max := inteiro(t, c.Limites.PrecoUnitarioMaxUnidades)
	if _, err := dinheiro.NovoPrecoUnitario(max, dinheiro.BRL); err != nil {
		t.Error(err)
	}
	if _, err := dinheiro.NovoPrecoUnitario(max+1, dinheiro.BRL); !errors.Is(err, dinheiro.ErrForaDaFaixa) {
		t.Errorf("preço acima do máximo: %v", err)
	}
}

func TestConstrucaoDaMoeda(t *testing.T) {
	for _, c := range carregar(t).ConstrucaoMoeda {
		m, err := dinheiro.ParaMoeda(c.Codigo)
		if c.Valido {
			if err != nil || m.Codigo() != c.Codigo || m.Casas() != c.Casas {
				t.Errorf("%s: (%v, %v)", c.Nome, m, err)
			}
			continue
		}
		if !errors.Is(err, sentinela(c.Erro)) {
			t.Errorf("%s: %q deveria dar %s, veio %v", c.Nome, c.Codigo, c.Erro, err)
		}
	}
	if _, err := dinheiro.NovoValor(1, dinheiro.Moeda{}); !errors.Is(err, dinheiro.ErrMoedaDesconhecida) {
		t.Errorf("moeda zero: %v", err)
	}
}

func TestLeitura(t *testing.T) {
	for _, c := range carregar(t).Leitura {
		m := moeda(t, c.Moeda)
		var (
			canonico string
			unidades int64
			ee       *erros.Erro
		)
		switch c.Especie {
		case "valor":
			var v dinheiro.Valor
			v, ee = dinheiro.LerValor("campo", c.Entrada, m, c.AceitaNegativo)
			canonico, unidades = v.String(), v.Unidades()
		case "preco_unitario":
			var p dinheiro.PrecoUnitario
			p, ee = dinheiro.LerPrecoUnitario("campo", c.Entrada, m, c.AceitaNegativo)
			canonico, unidades = p.String(), p.Micros()
		default:
			t.Fatalf("espécie %q", c.Especie)
		}
		if !c.Valido {
			if ee == nil || ee.Codigo != "validacao_falhou" || len(ee.Campos) != 1 ||
				ee.Campos[0].Campo != "campo" || ee.Campos[0].Codigo != c.Erro {
				t.Errorf("%s: %q deveria dar %s, veio %v", c.Nome, c.Entrada, c.Erro, ee)
			}
			continue
		}
		if ee != nil || canonico != c.Canonico || unidades != inteiro(t, c.Unidades) {
			t.Errorf("%s: %q = (%q, %d, %v), esperado (%q, %s)", c.Nome, c.Entrada, canonico, unidades, ee, c.Canonico, c.Unidades)
		}
	}
}

func TestTextoCanonico(t *testing.T) {
	for _, c := range carregar(t).TextoCanonico {
		m, n := moeda(t, c.Moeda), inteiro(t, c.Unidades)
		var got string
		if c.Especie == "valor" {
			v, err := dinheiro.NovoValor(n, m)
			if err != nil {
				t.Fatalf("%s: %v", c.Nome, err)
			}
			got = v.String()
		} else {
			p, err := dinheiro.NovoPrecoUnitario(n, m)
			if err != nil {
				t.Fatalf("%s: %v", c.Nome, err)
			}
			got = p.String()
		}
		if got != c.Esperado {
			t.Errorf("%s: %q, esperado %q", c.Nome, got, c.Esperado)
		}
	}
}

func valorDe(t *testing.T, cod, texto string) dinheiro.Valor {
	t.Helper()
	v, ee := dinheiro.LerValor("v", texto, moeda(t, cod), true)
	if ee != nil {
		t.Fatalf("valor %s %q: %v", cod, texto, ee)
	}
	return v
}

func sentinela(nome string) error {
	switch nome {
	case "moeda_diferente":
		return dinheiro.ErrMoedaDiferente
	case "fora_da_faixa":
		return dinheiro.ErrForaDaFaixa
	case "moeda_desconhecida":
		return dinheiro.ErrMoedaDesconhecida
	case "rateio_invalido":
		return dinheiro.ErrRateioInvalido
	}
	return nil
}

func TestSomaESubtracao(t *testing.T) {
	for _, c := range carregar(t).SomaSubtracao {
		a, b := valorDe(t, c.A.Moeda, c.A.Valor), valorDe(t, c.B.Moeda, c.B.Valor)
		var r dinheiro.Valor
		var err error
		if c.Operacao == "somar" {
			r, err = a.Somar(b)
		} else {
			r, err = a.Subtrair(b)
		}
		if c.Erro != "" {
			if !errors.Is(err, sentinela(c.Erro)) {
				t.Errorf("%s: erro %v, esperado %s", c.Nome, err, c.Erro)
			}
			continue
		}
		if err != nil || r.String() != c.Esperado {
			t.Errorf("%s: (%s, %v), esperado %s", c.Nome, r, err, c.Esperado)
		}
	}
}

func TestMultiplicacao(t *testing.T) {
	for _, c := range carregar(t).Multiplicacao {
		m := moeda(t, c.Moeda)
		p, ee := dinheiro.LerPrecoUnitario("preco", c.Preco, m, true)
		if ee != nil {
			t.Fatalf("%s: preço %q: %v", c.Nome, c.Preco, ee)
		}
		r, err := p.Multiplicar(quantidade(t, c.Quantidade))
		if c.Erro != "" {
			if !errors.Is(err, sentinela(c.Erro)) {
				t.Errorf("%s: erro %v, esperado %s", c.Nome, err, c.Erro)
			}
			continue
		}
		if err != nil || r.String() != c.Esperado || r.Unidades() != inteiro(t, c.Unidades) {
			t.Errorf("%s: (%s, %d, %v), esperado (%s, %s)", c.Nome, r, r.Unidades(), err, c.Esperado, c.Unidades)
		}
	}
}

func TestRateio(t *testing.T) {
	for _, c := range carregar(t).Rateio {
		v := valorDe(t, c.Moeda, c.Valor)
		pesos, err := dinheiro.LerPesos(c.Pesos)
		var partes []dinheiro.Valor
		if err == nil {
			partes, err = v.Ratear(pesos)
		}
		if c.Erro != "" {
			if !errors.Is(err, sentinela(c.Erro)) {
				t.Errorf("%s: erro %v, esperado %s", c.Nome, err, c.Erro)
			}
			continue
		}
		if err != nil || len(partes) != len(c.Esperado) {
			t.Errorf("%s: (%v, %v)", c.Nome, partes, err)
			continue
		}
		soma := valorDe(t, c.Moeda, "0")
		for i, p := range partes {
			if p.String() != c.Esperado[i] {
				t.Errorf("%s: parte %d = %s, esperado %s", c.Nome, i, p, c.Esperado[i])
			}
			soma, _ = soma.Somar(p)
		}
		if soma != v {
			t.Errorf("%s: a soma das partes %s difere do valor %s", c.Nome, soma, v)
		}
		// Ratear -v é o oposto de ratear v.
		if v.Unidades() != 0 {
			oposto, _ := dinheiro.NovoValor(-v.Unidades(), v.Moeda())
			op, err := oposto.Ratear(pesos)
			if err != nil {
				t.Fatal(err)
			}
			for i := range partes {
				if op[i].Unidades() != -partes[i].Unidades() {
					t.Errorf("%s: ratear o oposto diverge na parte %d", c.Nome, i)
				}
			}
		}
	}
}

func TestRateioInvalido(t *testing.T) {
	v := valorDe(t, "BRL", "10.00")
	r := func(s string) *big.Rat { return quantidade(t, s) }
	for nome, pesos := range map[string][]*big.Rat{
		"vazio": nil, "soma zero": {r("0"), r("0")}, "peso negativo": {r("1"), r("-1")}, "peso nulo": {r("1"), nil},
	} {
		if _, err := v.Ratear(pesos); !errors.Is(err, dinheiro.ErrRateioInvalido) {
			t.Errorf("%s: %v", nome, err)
		}
	}
}

// -0 vira zero sem sinal onde o campo aceita negativo; onde não aceita, o sinal
// é recusado, qualquer que seja o número.
func TestMenosZero(t *testing.T) {
	if v, ee := dinheiro.LerValor("v", "-0", dinheiro.BRL, true); ee != nil || v.String() != "0.00" || v.Unidades() != 0 {
		t.Errorf("-0 aceito: (%s, %v)", v, ee)
	}
	if _, ee := dinheiro.LerValor("v", "-0", dinheiro.BRL, false); ee == nil || ee.Campos[0].Codigo != "formato_invalido" {
		t.Errorf("-0 sem negativo: %v", ee)
	}
	if _, ee := dinheiro.LerPrecoUnitario("p", "-0.000000", dinheiro.PYG, false); ee == nil {
		t.Error("-0.000000 sem negativo foi aceito")
	}
}

func TestQuantidade(t *testing.T) {
	q, ee := dinheiro.LerQuantidade("q", "45.37", false)
	if ee != nil || q.Cmp(big.NewRat(4537, 100)) != 0 {
		t.Errorf("45.37: %v %v", q, ee)
	}
	if q, ee = dinheiro.LerQuantidade("q", "-0.005", true); ee != nil || q.Cmp(big.NewRat(-5, 1000)) != 0 {
		t.Errorf("-0.005: %v %v", q, ee)
	}
	for _, ruim := range []string{"", "1/2", "1e3", "1,5", " 1", "01", "-1", "1.1234567", "1000000000000"} {
		if _, ee := dinheiro.LerQuantidade("q", ruim, false); ee == nil || ee.Campos[0].Codigo != "formato_invalido" {
			t.Errorf("%q aceito: %v", ruim, ee)
		}
	}
}

func TestNumericExato(t *testing.T) {
	v, _ := dinheiro.NovoValor(-125040, dinheiro.BRL)
	n := v.ParaNumeric()
	if n.Int.Int64() != -125040 || n.Exp != -2 || !n.Valid {
		t.Errorf("numeric do valor: %+v", n)
	}
	volta, err := dinheiro.ValorDeNumeric(n, dinheiro.BRL)
	if err != nil || volta != v {
		t.Errorf("ida e volta: %v %v", volta, err)
	}
	p, _ := dinheiro.NovoPrecoUnitario(37846200, dinheiro.PYG)
	pn := p.ParaNumeric()
	if pn.Exp != -6 {
		t.Errorf("exp do preço: %d", pn.Exp)
	}
	if pv, err := dinheiro.PrecoDeNumeric(pn, dinheiro.PYG); err != nil || pv != p {
		t.Errorf("ida e volta do preço: %v %v", pv, err)
	}

	// Escalas diferentes do mesmo número: o banco devolve numeric(18,6) com Exp -6.
	seis := pgtype.Numeric{Int: big.NewInt(45500000), Exp: -6, Valid: true}
	if pv, err := dinheiro.PrecoDeNumeric(seis, dinheiro.BRL); err != nil || pv.String() != "45.50" {
		t.Errorf("45.500000: %v %v", pv, err)
	}
	// Zeros além da escala da moeda valem; casas que não são zeros, não: erro, nunca truncamento.
	if v, err := dinheiro.ValorDeNumeric(pgtype.Numeric{Int: big.NewInt(7500000), Exp: -2, Valid: true}, dinheiro.PYG); err != nil || v.String() != "75000" {
		t.Errorf("75000.00 em PYG: %v %v", v, err)
	}
	for nome, n := range map[string]pgtype.Numeric{
		"casas a mais":     {Int: big.NewInt(7500050), Exp: -2, Valid: true},
		"nulo":             {},
		"NaN":              {NaN: true, Valid: true},
		"infinito":         {InfinityModifier: pgtype.Infinity, Valid: true},
		"fora da faixa":    {Int: big.NewInt(1), Exp: 17, Valid: true},
		"estoura int64":    {Int: new(big.Int).Lsh(big.NewInt(1), 70), Exp: 0, Valid: true},
		"expoente absurdo": {Int: big.NewInt(1), Exp: 1000, Valid: true},
	} {
		if _, err := dinheiro.ValorDeNumeric(n, dinheiro.PYG); !errors.Is(err, dinheiro.ErrConversao) {
			t.Errorf("%s: %v", nome, err)
		}
	}
	if _, err := dinheiro.ValorDeNumeric(n, dinheiro.Moeda{}); !errors.Is(err, dinheiro.ErrMoedaDesconhecida) {
		t.Errorf("moeda zero: %v", err)
	}
}

// ErroDaOperacao: fora_da_faixa vira erro de campo; o resto é defeito.
func TestErroDaOperacao(t *testing.T) {
	err := dinheiro.ErroDaOperacao("total", dinheiro.ErrForaDaFaixa)
	var e *erros.Erro
	if !errors.As(err, &e) || e.Codigo != "validacao_falhou" || e.Campos[0].Campo != "total" || e.Campos[0].Codigo != "fora_da_faixa" {
		t.Errorf("fora da faixa: %v", err)
	}
	for _, s := range []error{dinheiro.ErrMoedaDiferente, dinheiro.ErrMoedaDesconhecida, dinheiro.ErrRateioInvalido} {
		if !errors.As(dinheiro.ErroDaOperacao("x", s), &e) || e.Codigo != "erro_interno" {
			t.Errorf("%v deveria ser erro_interno", s)
		}
	}
	if dinheiro.ErroDaOperacao("x", nil) != nil {
		t.Error("nil")
	}
}
