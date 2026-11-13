package pais

import (
	"errors"
	"strings"
)

// Tipos de identificação (identificacao_tipo). Catálogo fechado do perfil de
// país; o banco confere a forma de cada um (app.identificacao_valida) e o
// dígito verificador é conferido aqui e em pacotes/dominio.
const (
	TipoCNPJ = "cnpj"
	TipoCPF  = "cpf"
	TipoCUIT = "cuit"
	TipoCUIL = "cuil"
	TipoDNI  = "dni"
	TipoRUC  = "ruc"
	TipoCI   = "ci"
)

// Erros da validação de identificação. O chamador os traduz em erro de campo.
var (
	// ErrTipoDesconhecido: o tipo não está no catálogo.
	ErrTipoDesconhecido = errors.New("tipo de identificação desconhecido")
	// ErrFormato: tamanho ou alfabeto fora da forma do tipo.
	ErrFormato = errors.New("forma da identificação inválida")
	// ErrDigitoVerificador: a forma está certa e o dígito verificador não bate.
	ErrDigitoVerificador = errors.New("dígito verificador inválido")
)

// Normalizar põe em maiúsculas e mantém só letras e dígitos ASCII: a máscara
// (pontos, hífens, barras, espaços) é só de exibição.
func Normalizar(valor string) string {
	var b strings.Builder
	for _, r := range strings.ToUpper(valor) {
		if (r >= '0' && r <= '9') || (r >= 'A' && r <= 'Z') {
			b.WriteRune(r)
		}
	}
	return b.String()
}

// ValidarIdentificacao normaliza o valor e confere a forma do tipo, e o dígito
// verificador onde ele é conferido (CNPJ, inclusive o alfanumérico, e CPF).
//
// CUIT e RUC só pela forma, até a pendência P13 (a regra do dígito do CUIT
// quando o resultado é 10 e a do RUC antigo com letra não estão confirmadas em
// fonte oficial); CUIL, DNI e CI não têm dígito conferido. Nenhum validador
// consulta serviço externo: o dígito pega erro de digitação, não prova que o
// número existe.
func ValidarIdentificacao(tipo, valor string) (string, error) {
	if !tipoConhecido(tipo) {
		return "", ErrTipoDesconhecido
	}
	n := Normalizar(valor)
	if !formaValida(tipo, n) {
		return "", ErrFormato
	}
	switch tipo {
	case TipoCNPJ:
		if !cnpjDigitosConferem(n) {
			return "", ErrDigitoVerificador
		}
	case TipoCPF:
		if !cpfDigitosConferem(n) {
			return "", ErrDigitoVerificador
		}
	}
	return n, nil
}

func tipoConhecido(tipo string) bool {
	switch tipo {
	case TipoCNPJ, TipoCPF, TipoCUIT, TipoCUIL, TipoDNI, TipoRUC, TipoCI:
		return true
	}
	return false
}

// formaValida repete, em Go, app.identificacao_valida (migração 00009): o
// mesmo alfabeto e os mesmos tamanhos, para recusar antes de ir ao banco.
func formaValida(tipo, n string) bool {
	switch tipo {
	case TipoCNPJ:
		// Letras só nas 12 primeiras posições (CNPJ alfanumérico, IN RFB
		// 2.229/2024); os dois dígitos verificadores são numéricos.
		return len(n) == 14 && alfabeto(n[:12], true) && alfabeto(n[12:], false)
	case TipoCPF, TipoCUIT, TipoCUIL:
		return len(n) == 11 && alfabeto(n, false)
	case TipoDNI:
		return len(n) >= 7 && len(n) <= 8 && alfabeto(n, false)
	case TipoRUC:
		return len(n) >= 6 && len(n) <= 10 && alfabeto(n, true) // RUC antigo pode ter letra (P13)
	case TipoCI:
		return len(n) >= 4 && len(n) <= 8 && alfabeto(n, false)
	}
	return false
}

func alfabeto(s string, comLetras bool) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			continue
		}
		if comLetras && c >= 'A' && c <= 'Z' {
			continue
		}
		return false
	}
	return true
}

// repetido: todos os caracteres iguais (000...0, 111...1) passam no cálculo do
// dígito e não são documentos.
func repetido(s string) bool {
	return strings.Count(s, s[:1]) == len(s)
}

func cpfDigitosConferem(cpf string) bool {
	if repetido(cpf) {
		return false
	}
	dv := func(n int) int {
		soma := 0
		for i := 0; i < n; i++ {
			soma += int(cpf[i]-'0') * (n + 1 - i)
		}
		r := soma * 10 % 11
		if r == 10 {
			r = 0
		}
		return r
	}
	return dv(9) == int(cpf[9]-'0') && dv(10) == int(cpf[10]-'0')
}

// cnpjDigitosConferem calcula o dígito com o valor de cada caractere igual ao
// código ASCII menos 48 ('0' a '9' valem 0 a 9, 'A' vale 17, 'Z' vale 42), os
// mesmos pesos do CNPJ numérico e módulo 11 (IN RFB 2.229/2024; método
// divulgado pelo Serpro). Para CNPJ só de dígitos o resultado é o do cálculo
// clássico. O exemplo oficial 12ABC34501DE35 está nos casos compartilhados.
func cnpjDigitosConferem(cnpj string) bool {
	if repetido(cnpj) {
		return false
	}
	pesos1 := []int{5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	pesos2 := []int{6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2}
	dv := func(pesos []int) int {
		soma := 0
		for i, p := range pesos {
			soma += int(cnpj[i]-'0') * p
		}
		r := soma % 11
		if r < 2 {
			return 0
		}
		return 11 - r
	}
	return dv(pesos1) == int(cnpj[12]-'0') && dv(pesos2) == int(cnpj[13]-'0')
}

// Mascarar esconde a identificação pessoal para a lista (a completa é dado
// sensível e só sai no detalhe, com a visualização gravada na trilha). O CPF
// mostra o miolo, ***.456.789-**; os demais tipos mostram os 3 últimos
// caracteres e escondem o resto, sem mudar o tamanho. Valor vazio dá vazio.
func Mascarar(tipo, valor string) string {
	if valor == "" {
		return ""
	}
	if tipo == TipoCPF && len(valor) == 11 {
		return "***." + valor[3:6] + "." + valor[6:9] + "-**"
	}
	if len(valor) <= 3 {
		return strings.Repeat("*", len(valor))
	}
	return strings.Repeat("*", len(valor)-3) + valor[len(valor)-3:]
}
