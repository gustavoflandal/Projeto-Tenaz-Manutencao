package dinheiro

import (
	"math/big"
	"sort"

	"github.com/jackc/pgx/v5/pgtype"
)

// Somar soma dois valores da mesma moeda. Moedas diferentes:
// ErrMoedaDiferente. Resultado fora da coluna: ErrForaDaFaixa. A soma de dois
// int64 dentro da faixa cabe em int64 (10^18 + 10^18 < 2^63).
func (v Valor) Somar(o Valor) (Valor, error) {
	if v.moeda != o.moeda {
		return Valor{}, ErrMoedaDiferente
	}
	return NovoValor(v.unidades+o.unidades, v.moeda)
}

// Subtrair subtrai dois valores da mesma moeda, com as mesmas verificações de
// Somar.
func (v Valor) Subtrair(o Valor) (Valor, error) {
	if v.moeda != o.moeda {
		return Valor{}, ErrMoedaDiferente
	}
	return NovoValor(v.unidades-o.unidades, v.moeda)
}

// ArredondarMeioLongeDeZero divide num por den (den > 0) e arredonda o
// quociente para o inteiro mais próximo, com o meio indo para longe do zero
// (2,345 vira 2,35 e -2,345 vira -2,35), como o round(numeric, n) do
// PostgreSQL.
func arredondarMeioLongeDeZero(num, den *big.Int) *big.Int {
	neg := num.Sign() < 0
	abs := new(big.Int).Abs(num)
	q, r := new(big.Int).QuoRem(abs, den, new(big.Int))
	if r.Lsh(r, 1).Cmp(den) >= 0 {
		q.Add(q, big.NewInt(1))
	}
	if neg {
		q.Neg(q)
	}
	return q
}

// Multiplicar multiplica o preço unitário por uma quantidade exata (fração:
// minutos e 60, ou o texto decimal lido por LerQuantidade) e devolve o Valor na
// escala da moeda. O cálculo é exato em math/big e tem UM arredondamento, no
// fim, meio para longe do zero; nunca nas parcelas intermediárias. A
// quantidade pode ser negativa (estorno). Resultado fora da coluna:
// ErrForaDaFaixa.
func (p PrecoUnitario) Multiplicar(quantidade *big.Rat) (Valor, error) {
	if !p.moeda.Valida() {
		return Valor{}, ErrMoedaDesconhecida
	}
	// valor em unidades menores = micros * q / 10^(6 - casas)
	num := new(big.Int).Mul(big.NewInt(p.micros), quantidade.Num())
	den := new(big.Int).Mul(quantidade.Denom(), new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(casasPreco-p.moeda.Casas())), nil))
	r := arredondarMeioLongeDeZero(num, den)
	if !r.IsInt64() {
		return Valor{}, ErrForaDaFaixa
	}
	return NovoValor(r.Int64(), p.moeda)
}

// Ratear reparte o valor em partes proporcionais aos pesos (decimais não
// negativos, lista não vazia e soma positiva; senão ErrRateioInvalido). Cada
// parte recebe o piso do seu quinhão sobre o MÓDULO do valor; a sobra de
// unidades menores vai uma a uma pelo maior resto, e o empate de restos vai
// para o menor índice (ordem de entrada); o sinal do valor volta a todas as
// partes. A soma das partes é sempre o valor, e ratear -v é o oposto de ratear
// v.
func (v Valor) Ratear(pesos []*big.Rat) ([]Valor, error) {
	if len(pesos) == 0 {
		return nil, ErrRateioInvalido
	}
	soma := new(big.Rat)
	for _, w := range pesos {
		if w == nil || w.Sign() < 0 {
			return nil, ErrRateioInvalido
		}
		soma.Add(soma, w)
	}
	if soma.Sign() <= 0 {
		return nil, ErrRateioInvalido
	}
	mod := v.unidades
	if mod < 0 {
		mod = -mod
	}
	base := make([]int64, len(pesos))
	restos := make([]*big.Rat, len(pesos))
	var dado int64
	for i, w := range pesos {
		quinhao := new(big.Rat).Mul(big.NewRat(mod, 1), w)
		quinhao.Quo(quinhao, soma)
		piso := new(big.Int).Quo(quinhao.Num(), quinhao.Denom()) // quinhão >= 0: piso
		base[i] = piso.Int64()
		dado += base[i]
		restos[i] = quinhao.Sub(quinhao, new(big.Rat).SetInt(piso))
	}
	sobra := mod - dado // menor que o número de partes
	ordem := make([]int, len(pesos))
	for i := range ordem {
		ordem[i] = i
	}
	sort.SliceStable(ordem, func(a, b int) bool { return restos[ordem[a]].Cmp(restos[ordem[b]]) > 0 })
	for k := int64(0); k < sobra; k++ {
		base[ordem[k]]++
	}
	partes := make([]Valor, len(pesos))
	for i, n := range base {
		if v.unidades < 0 {
			n = -n
		}
		partes[i] = Valor{n, v.moeda}
	}
	return partes, nil
}

// ParaNumeric converte o valor para o pgtype.Numeric, exato: Int são as
// unidades menores e Exp é -casas da moeda. Sem float.
func (v Valor) ParaNumeric() pgtype.Numeric {
	return pgtype.Numeric{Int: big.NewInt(v.unidades), Exp: int32(-v.moeda.Casas()), Valid: true}
}

// ParaNumeric converte o preço para o pgtype.Numeric, exato (Exp = -6).
func (p PrecoUnitario) ParaNumeric() pgtype.Numeric {
	return pgtype.Numeric{Int: big.NewInt(p.micros), Exp: -casasPreco, Valid: true}
}

// unidadesDe converte o numeric para um inteiro na escala dada, sem truncar:
// NaN, infinito, nulo, casas a mais que não sejam zeros ou fora de int64 são
// ErrConversao.
func unidadesDe(n pgtype.Numeric, casas int) (int64, error) {
	if !n.Valid || n.NaN || n.InfinityModifier != pgtype.Finite || n.Int == nil {
		return 0, ErrConversao
	}
	// valor = Int * 10^Exp; unidades = valor * 10^casas = Int * 10^(Exp+casas)
	desloc := int(n.Exp) + casas
	r := new(big.Int).Set(n.Int)
	switch {
	case desloc > 40:
		return 0, ErrConversao // evita 10^huge; nenhum valor da coluna chega aqui
	case desloc >= 0:
		r.Mul(r, new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(desloc)), nil))
	default:
		div := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(-desloc)), nil)
		q, resto := new(big.Int).QuoRem(r, div, new(big.Int))
		if resto.Sign() != 0 {
			return 0, ErrConversao
		}
		r = q
	}
	if !r.IsInt64() {
		return 0, ErrConversao
	}
	return r.Int64(), nil
}

// ValorDeNumeric lê o numeric do banco como valor da moeda, de modo exato: casas
// além da escala da moeda (que não sejam zeros), fora da coluna ou não finito
// são erro, nunca truncamento.
func ValorDeNumeric(n pgtype.Numeric, m Moeda) (Valor, error) {
	if !m.Valida() {
		return Valor{}, ErrMoedaDesconhecida
	}
	u, err := unidadesDe(n, m.Casas())
	if err != nil {
		return Valor{}, err
	}
	v, err := NovoValor(u, m)
	if err != nil {
		return Valor{}, ErrConversao
	}
	return v, nil
}

// PrecoDeNumeric lê o numeric do banco como preço unitário, com até seis
// casas, de modo exato.
func PrecoDeNumeric(n pgtype.Numeric, m Moeda) (PrecoUnitario, error) {
	if !m.Valida() {
		return PrecoUnitario{}, ErrMoedaDesconhecida
	}
	u, err := unidadesDe(n, casasPreco)
	if err != nil {
		return PrecoUnitario{}, err
	}
	p, err := NovoPrecoUnitario(u, m)
	if err != nil {
		return PrecoUnitario{}, ErrConversao
	}
	return p, nil
}
