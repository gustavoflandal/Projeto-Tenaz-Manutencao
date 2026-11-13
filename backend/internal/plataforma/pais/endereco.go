package pais

import (
	"errors"
	"regexp"
	"slices"
	"strings"
)

// Endereço multipaís (documento 3, seção 2.18; nota docs/pesquisa/enderecos_br_ar_py.md):
// o país do endereço é sempre o da empresa, e subdivisão e código postal são
// conferidos pelo perfil dele. Os nomes das subdivisões ficam na interface; o
// servidor só confere o código.

// ErrCodigoPostalFormato é o código postal fora da forma do país.
var ErrCodigoPostalFormato = errors.New("forma do código postal inválida")

// Subdivisões ISO 3166-2, na ordem do contrato. Paraguai sem zero à esquerda
// (PY-1 a PY-16, PY-19, PY-ASU); Argentina de A a Z sem I nem O.
var (
	subdivisoesBR = []string{"BR-AC", "BR-AL", "BR-AP", "BR-AM", "BR-BA", "BR-CE", "BR-DF", "BR-ES", "BR-GO", "BR-MA",
		"BR-MT", "BR-MS", "BR-MG", "BR-PA", "BR-PB", "BR-PR", "BR-PE", "BR-PI", "BR-RJ", "BR-RN", "BR-RS", "BR-RO",
		"BR-RR", "BR-SC", "BR-SP", "BR-SE", "BR-TO"}
	subdivisoesAR = []string{"AR-A", "AR-B", "AR-C", "AR-D", "AR-E", "AR-F", "AR-G", "AR-H", "AR-J", "AR-K", "AR-L",
		"AR-M", "AR-N", "AR-P", "AR-Q", "AR-R", "AR-S", "AR-T", "AR-U", "AR-V", "AR-W", "AR-X", "AR-Y", "AR-Z"}
	subdivisoesPY = []string{"PY-ASU", "PY-1", "PY-2", "PY-3", "PY-4", "PY-5", "PY-6", "PY-7", "PY-8", "PY-9", "PY-10",
		"PY-11", "PY-12", "PY-13", "PY-14", "PY-15", "PY-16", "PY-19"}
)

// Forma do código postal normalizado, por país.
var (
	formaCEP = regexp.MustCompile(`^\d{8}$`)
	// CPA (letra, quatro dígitos, três letras) ou o código antigo de quatro dígitos.
	formaCPAouAntigo = regexp.MustCompile(`^([A-HJ-NP-Z]\d{4}[A-Z]{3}|\d{4})$`)
	formaPY          = regexp.MustCompile(`^\d{6}$`)
)

// formaCodigoPostal devolve a expressão do país, ou nil se não há perfil.
func formaCodigoPostal(codigo string) *regexp.Regexp {
	switch codigo {
	case BR:
		return formaCEP
	case AR:
		return formaCPAouAntigo
	case PY:
		return formaPY
	}
	return nil
}

// SubdivisaoValida diz se o código ISO 3166-2 está na lista fechada do país.
func SubdivisaoValida(p PerfilPais, subdivisao string) bool {
	return slices.Contains(p.Subdivisoes, subdivisao)
}

// ValidarCodigoPostal normaliza (maiúsculas, sem espaço nem hífen) e confere a
// forma do país. Devolve o valor normalizado, ou ErrCodigoPostalFormato.
func ValidarCodigoPostal(p PerfilPais, valor string) (string, error) {
	re := formaCodigoPostal(p.Codigo)
	n := Normalizar(valor)
	if re == nil || !re.MatchString(n) {
		return "", ErrCodigoPostalFormato
	}
	return n, nil
}

// CPAConfereComSubdivisao é a regra da Tenaz para a Argentina, a mesma do
// TypeScript (cpaConfereComSubdivisao em pacotes/dominio): quem decide é o
// prefixo da subdivisão, não o país da empresa. Com subdivisão "AR-" e um CPA
// de 8 caracteres, a primeira letra do CPA é a da subdivisão (C1425DKF só com
// AR-C). Subdivisão vazia, de outro país ou sem prefixo, e o código antigo de
// 4 dígitos, não cruzam e o resultado é verdadeiro. O código postal chega como
// a pessoa digitou: a função o normaliza sozinha.
func CPAConfereComSubdivisao(subdivisao, codigoPostal string) bool {
	letra, ok := strings.CutPrefix(subdivisao, "AR-")
	cpa := Normalizar(codigoPostal)
	if !ok || len(cpa) != 8 {
		return true
	}
	return strings.HasPrefix(cpa, letra)
}
