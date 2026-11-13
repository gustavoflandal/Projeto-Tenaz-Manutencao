// Package dinheiro é a representação do dinheiro da Tenaz (ADR-013, Parte F):
// duas espécies, valor e preço unitário, ambas com a moeda da linha ao lado,
// sobre int64 e math/big, só com a biblioteca padrão.
//
//   - Valor: uma quantia que se paga, se lança ou se soma. int64 na unidade
//     menor da moeda (o centavo em BRL e ARS, o próprio guarani em PYG), na
//     escala da moeda. Cabe em numeric(18,2) (BRL e ARS) e em numeric(18,0)
//     efetivo (PYG: 16 dígitos inteiros).
//   - PrecoUnitario: dinheiro por unidade de uma quantidade (hora, litro,
//     quilo). int64 em milionésimos da unidade da moeda; não é limitado à
//     escala da moeda e não se soma: vira Valor só multiplicado por uma
//     quantidade, com um único arredondamento (meio para longe do zero).
//
// A escala é propriedade da moeda, não do país. Nenhum float32 ou float64 em
// dinheiro, em nenhuma camada. Os mesmos casos de
// backend/testes/casos/dinheiro.json valem para este pacote e para
// pacotes/dominio/src/dinheiro.ts.
package dinheiro

import "errors"

// Erros da biblioteca. Nenhum chega à API com nome próprio: no MVP todo valor
// de uma empresa tem a moeda dela e os pesos de rateio vêm do código, então
// MoedaDiferente, MoedaDesconhecida e RateioInvalido são defeito e saem como
// erro_interno. ForaDaFaixa, vindo de uma operação sobre entrada da pessoa, vira
// validacao_falhou com fora_da_faixa no campo (ErroDaOperacao).
var (
	ErrMoedaDesconhecida = errors.New("moeda_desconhecida")
	ErrMoedaDiferente    = errors.New("moeda_diferente")
	ErrForaDaFaixa       = errors.New("fora_da_faixa")
	ErrRateioInvalido    = errors.New("rateio_invalido")
	// ErrConversao: o numeric não cabe exatamente na espécie (NaN, infinito,
	// casas a mais ou fora da faixa). Nunca se trunca.
	ErrConversao = errors.New("conversao_numeric_inexata")
)

// Moeda é uma moeda do catálogo fechado (ISO 4217). O valor zero não é moeda
// válida: só ParaMoeda e as variáveis BRL, ARS e PYG a constroem.
type Moeda struct{ codigo string }

// As moedas do MVP. Moeda nova entra com país novo, por requisito, código e
// teste.
var (
	BRL = Moeda{"BRL"}
	ARS = Moeda{"ARS"}
	PYG = Moeda{"PYG"}
)

// ParaMoeda constrói a moeda de um código. Código fora do catálogo é defeito,
// porque o cliente nunca envia moeda: ErrMoedaDesconhecida.
func ParaMoeda(codigo string) (Moeda, error) {
	switch codigo {
	case "BRL":
		return BRL, nil
	case "ARS":
		return ARS, nil
	case "PYG":
		return PYG, nil
	}
	return Moeda{}, ErrMoedaDesconhecida
}

// Codigo devolve o código ISO 4217 ("" na moeda zero).
func (m Moeda) Codigo() string { return m.codigo }

// String devolve o código.
func (m Moeda) String() string { return m.codigo }

// Valida diz se é uma moeda do catálogo.
func (m Moeda) Valida() bool { return m == BRL || m == ARS || m == PYG }

// Casas devolve as casas decimais da moeda: 2 em BRL e ARS, 0 em PYG. Na moeda
// zero, 0.
func (m Moeda) Casas() int {
	switch m {
	case BRL, ARS:
		return 2
	}
	return 0
}

// Limites das colunas: valor numeric(18,2), preço unitário numeric(18,6).
// maxPreco é 10^18 - 1 milionésimos; o do valor depende da moeda
// (10^(16+casas) - 1 unidades menores: 16 dígitos inteiros).
const (
	maxPreco = 999999999999999999
	// casasPreco é a escala fixa do preço unitário.
	casasPreco = 6
)

func pot10(n int) int64 {
	r := int64(1)
	for i := 0; i < n; i++ {
		r *= 10
	}
	return r
}

// maxValor é o maior valor em unidades menores da moeda.
func maxValor(m Moeda) int64 { return pot10(16+m.Casas()) - 1 }

// Valor é uma quantia na escala da moeda.
type Valor struct {
	unidades int64
	moeda    Moeda
}

// NovoValor constrói o valor de uma quantidade de unidades menores da moeda.
// Moeda desconhecida: ErrMoedaDesconhecida; fora da coluna: ErrForaDaFaixa.
func NovoValor(unidades int64, m Moeda) (Valor, error) {
	if !m.Valida() {
		return Valor{}, ErrMoedaDesconhecida
	}
	if unidades > maxValor(m) || unidades < -maxValor(m) {
		return Valor{}, ErrForaDaFaixa
	}
	return Valor{unidades, m}, nil
}

// Unidades devolve o valor em unidades menores da moeda.
func (v Valor) Unidades() int64 { return v.unidades }

// Moeda devolve a moeda do valor.
func (v Valor) Moeda() Moeda { return v.moeda }

// PrecoUnitario é o dinheiro de uma unidade de quantidade, em milionésimos da
// unidade da moeda. Não tem Somar nem Subtrair: preço não se soma, só se
// multiplica por uma quantidade e vira Valor.
type PrecoUnitario struct {
	micros int64
	moeda  Moeda
}

// NovoPrecoUnitario constrói o preço de uma quantidade de milionésimos.
func NovoPrecoUnitario(micros int64, m Moeda) (PrecoUnitario, error) {
	if !m.Valida() {
		return PrecoUnitario{}, ErrMoedaDesconhecida
	}
	if micros > maxPreco || micros < -maxPreco {
		return PrecoUnitario{}, ErrForaDaFaixa
	}
	return PrecoUnitario{micros, m}, nil
}

// Micros devolve o preço em milionésimos da unidade da moeda.
func (p PrecoUnitario) Micros() int64 { return p.micros }

// Moeda devolve a moeda do preço.
func (p PrecoUnitario) Moeda() Moeda { return p.moeda }
