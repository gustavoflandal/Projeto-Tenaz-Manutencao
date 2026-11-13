package dinheiro

import (
	"errors"
	"math/big"
	"strings"

	"tenaz/internal/plataforma/erros"
	"tenaz/internal/plataforma/textos"
)

// Gramática da leitura: exatamente a do padrão do contrato (documento 5, Parte
// F, item 4), e só depois a conferência da moeda.
//
//	valor:           -?(0|[1-9][0-9]{0,15})(\.[0-9]{1,2})?
//	preço unitário:  -?(0|[1-9][0-9]{0,11})(\.[0-9]{1,6})?
//	quantidade:      -?(0|[1-9][0-9]{0,11})(\.[0-9]{1,6})?
//
// O sinal - só onde o campo aceita negativo; em campo sem negativo, qualquer
// sinal é recusado, inclusive -0. Só dígitos ASCII; nada de espaço, +, vírgula,
// separador de milhar, expoente ou zero à esquerda. Erro de leitura é sempre
// formato_invalido; a chave do texto diz à pessoa qual foi o problema.

type problema int

const (
	pFormato problema = iota
	pSemNegativo
	pCasasDaMoeda
)

// partes é o número lido: sinal, parte inteira e parte fracionária (sem ponto).
type partes struct {
	negativo  bool
	inteira   string
	fracao    string
	problema  problema
	reconhece bool
}

func digitos(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return len(s) > 0
}

// decompor aplica a gramática do padrão: maxInt dígitos inteiros e maxFrac
// casas. Devolve reconhece falso se o texto não casa.
func decompor(texto string, maxInt, maxFrac int, aceitaNegativo bool) partes {
	r := partes{problema: pFormato}
	t := texto
	if strings.HasPrefix(t, "-") {
		if !aceitaNegativo {
			r.problema = pSemNegativo
			return r
		}
		r.negativo = true
		t = t[1:]
	}
	inteira, fracao, temPonto := strings.Cut(t, ".")
	if !digitos(inteira) || len(inteira) > maxInt || (inteira[0] == '0' && len(inteira) > 1) {
		return r
	}
	if temPonto && (!digitos(fracao) || len(fracao) > maxFrac) {
		return r
	}
	r.inteira, r.fracao, r.reconhece = inteira, fracao, true
	return r
}

func erroLeitura(campo string, p problema, moeda Moeda, chaveFormato string) *erros.Erro {
	chave := chaveFormato
	var params textos.Parametros
	switch p {
	case pSemNegativo:
		chave = "dinheiro.sem_negativo"
	case pCasasDaMoeda:
		chave, params = "dinheiro.casas_da_moeda", textos.Parametros{"moeda": moeda.Codigo()}
	}
	return erros.EntradaInvalida.ComCampoChave(campo, erros.CampoFormatoInvalido, chave, params)
}

// ParaInt64 junta inteira e fração (já preenchida até casas) em unidades.
func juntar(inteira, fracao string, casas int) int64 {
	for len(fracao) < casas {
		fracao += "0"
	}
	var n int64
	for _, c := range inteira + fracao[:casas] {
		n = n*10 + int64(c-'0')
	}
	return n
}

// LerValor lê o texto decimal de um valor na moeda dada, recusando (sem
// arredondar) o que não é exatamente representável: casas além da escala que
// não sejam zeros, sinal onde o campo não aceita negativo, formato fora da
// gramática. A moeda é um Moeda, nunca texto: vem da empresa. Erro:
// validacao_falhou com formato_invalido no campo.
func LerValor(campo, texto string, m Moeda, aceitaNegativo bool) (Valor, *erros.Erro) {
	if !m.Valida() {
		return Valor{}, erros.Interno.Tecnico("moeda desconhecida na leitura de %s", campo)
	}
	p := decompor(texto, 16, 2, aceitaNegativo)
	if !p.reconhece {
		return Valor{}, erroLeitura(campo, p.problema, m, "dinheiro.formato_valor")
	}
	casas := m.Casas()
	// Casas além da escala da moeda só valem se forem zeros ("75000.00" em PYG).
	if len(p.fracao) > casas && strings.Trim(p.fracao[casas:], "0") != "" {
		return Valor{}, erroLeitura(campo, pCasasDaMoeda, m, "")
	}
	n := juntar(p.inteira, p.fracao, casas)
	if p.negativo {
		n = -n
	}
	v, err := NovoValor(n, m)
	if err != nil {
		// A faixa do valor é o número de dígitos do padrão: fora dela já é
		// erro de forma.
		return Valor{}, erroLeitura(campo, pFormato, m, "dinheiro.formato_valor")
	}
	return v, nil
}

// LerPrecoUnitario lê o texto decimal de um preço unitário, com até seis
// casas. Mesmas recusas de LerValor, sem a escala da moeda.
func LerPrecoUnitario(campo, texto string, m Moeda, aceitaNegativo bool) (PrecoUnitario, *erros.Erro) {
	if !m.Valida() {
		return PrecoUnitario{}, erros.Interno.Tecnico("moeda desconhecida na leitura de %s", campo)
	}
	p := decompor(texto, 12, casasPreco, aceitaNegativo)
	if !p.reconhece {
		return PrecoUnitario{}, erroLeitura(campo, p.problema, m, "dinheiro.formato_preco")
	}
	n := juntar(p.inteira, p.fracao, casasPreco)
	if p.negativo {
		n = -n
	}
	pr, err := NovoPrecoUnitario(n, m)
	if err != nil {
		return PrecoUnitario{}, erroLeitura(campo, pFormato, m, "dinheiro.formato_preco")
	}
	return pr, nil
}

// LerQuantidade lê o texto decimal de uma quantidade (horas, litros, quilos)
// como fração exata. A notação a/b não existe na API: a fração de minutos
// sobre 60 o código monta com big.NewRat.
func LerQuantidade(campo, texto string, aceitaNegativo bool) (*big.Rat, *erros.Erro) {
	p := decompor(texto, 12, casasPreco, aceitaNegativo)
	if !p.reconhece {
		return nil, erroLeitura(campo, p.problema, Moeda{}, "dinheiro.formato_quantidade")
	}
	num := new(big.Int)
	num.SetString(p.inteira+p.fracao, 10)
	if p.negativo {
		num.Neg(num)
	}
	den := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(len(p.fracao))), nil)
	return new(big.Rat).SetFrac(num, den), nil
}

// LerPesos lê os pesos de um rateio: textos decimais sem sinal (a mesma
// gramática da quantidade). Lista vazia, peso com sinal (inclusive -0) ou fora
// da gramática é ErrRateioInvalido, defeito de quem monta os pesos, e não erro
// de campo; a soma positiva é conferida por Ratear.
func LerPesos(textos []string) ([]*big.Rat, error) {
	if len(textos) == 0 {
		return nil, ErrRateioInvalido
	}
	pesos := make([]*big.Rat, len(textos))
	for i, t := range textos {
		q, ee := LerQuantidade("peso", t, false)
		if ee != nil {
			return nil, ErrRateioInvalido
		}
		pesos[i] = q
	}
	return pesos, nil
}

func formatar(n int64, casas, minCasas int) string {
	neg := n < 0
	var u uint64
	if neg {
		u = uint64(-(n + 1)) + 1 // evita o estouro de -MinInt64
	} else {
		u = uint64(n)
	}
	d := uint64(pot10(casas))
	inteira, frac := u/d, u%d
	var b strings.Builder
	if neg && u != 0 {
		b.WriteByte('-')
	}
	b.WriteString(uitoa(inteira))
	if casas == 0 {
		return b.String()
	}
	f := uitoa(frac)
	f = strings.Repeat("0", casas-len(f)) + f
	// Corta zeros à direita até o mínimo de casas (as da moeda, no preço).
	for len(f) > minCasas && f[len(f)-1] == '0' {
		f = f[:len(f)-1]
	}
	if f != "" {
		b.WriteByte('.')
		b.WriteString(f)
	}
	return b.String()
}

func uitoa(u uint64) string {
	if u == 0 {
		return "0"
	}
	var buf [20]byte
	i := len(buf)
	for u > 0 {
		i--
		buf[i] = byte('0' + u%10)
		u /= 10
	}
	return string(buf[i:])
}

// String devolve o texto canônico do valor: exatamente as casas da moeda
// ("1250.40" em BRL, "75000" em PYG), nunca "-0".
func (v Valor) String() string { return formatar(v.unidades, v.moeda.Casas(), v.moeda.Casas()) }

// String devolve o texto canônico do preço: no mínimo as casas da moeda e no
// máximo seis, sem zeros à direita além das casas da moeda ("45.50",
// "37.8462", "25000" e "0.5" em PYG), nunca "-0".
func (p PrecoUnitario) String() string { return formatar(p.micros, casasPreco, p.moeda.Casas()) }

// ErroDaOperacao traduz o erro de uma operação sobre entrada da pessoa:
// ErrForaDaFaixa vira validacao_falhou com fora_da_faixa no campo que a
// originou; os demais (moeda diferente, moeda desconhecida, rateio inválido)
// são defeito e saem como erro_interno, registrado no log.
func ErroDaOperacao(campo string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, ErrForaDaFaixa) {
		return erros.EntradaInvalida.ComCampoChave(campo, erros.CampoForaDaFaixa, "dinheiro.fora_da_faixa", nil)
	}
	return erros.Interno.Tecnico("operação de dinheiro: %v", err).Por(err)
}
