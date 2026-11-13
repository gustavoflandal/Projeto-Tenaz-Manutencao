// Package pais é o perfil de país da Tenaz (ADR-013, aceito em 13/11/2026):
// o que muda de um país para outro e é regra do núcleo, não configuração.
//
// Hoje: moeda, idioma padrão sugerido e os tipos de identificação aceitos de
// empresa, unidade e pessoa. Fuso, formatação e feriados entram no perfil
// quando as etapas que dependem deles chegarem.
//
// Nenhum módulo testa o país por conta própria: pergunta ao perfil
// (pais.Perfil). País novo entra por requisito, código e teste, nunca por
// configuração. A tabela é a mesma de pacotes/dominio (TypeScript), com os
// casos JSON compartilhados de backend/testes/casos/identificacao.json.
package pais

import (
	"slices"

	"tenaz/internal/plataforma/dinheiro"
)

// Códigos ISO 3166-1 alfa-2 dos países suportados (ck_empresa_pais).
const (
	BR = "BR"
	AR = "AR"
	PY = "PY"
)

// Idiomas suportados, etiquetas BCP 47 (ck_empresa_idioma). Os mesmos de
// plataforma/textos.
const (
	IdiomaPtBR  = "pt-BR"
	IdiomaEs419 = "es-419"
)

// PerfilPais é o perfil de um país.
type PerfilPais struct {
	Codigo string
	// Moeda é a única da empresa, dada pelo país (as casas decimais são da
	// moeda: plataforma/dinheiro).
	Moeda        dinheiro.Moeda
	IdiomaPadrao string
	// FusoPadrao é o fuso IANA sugerido quando a conta não informa um.
	FusoPadrao string
	// Fusos são os fusos IANA válidos para empresa e unidade do país (a base tz
	// dá as regras de horário; aqui só a lista de quais nomes o país aceita).
	Fusos []string
	// TipoEmpresa é o tipo de identificação de pessoa jurídica do país: o único
	// que a empresa e a unidade aceitam.
	TipoEmpresa string
	// TiposPessoa são os tipos de pessoa física do país. A pessoa aceita os de
	// qualquer país (TiposPessoaDeTodos).
	TiposPessoa []string
	// Subdivisoes são os códigos ISO 3166-2 válidos para o endereço do país
	// (27, 24 e 18), na ordem do contrato.
	Subdivisoes []string
}

// Fusos do país: as zonas do arquivo zone1970.tab da base tz com o país. Nome
// IANA, nunca deslocamento fixo.
var (
	fusosBR = []string{"America/Noronha", "America/Belem", "America/Fortaleza", "America/Recife", "America/Araguaina",
		"America/Maceio", "America/Bahia", "America/Sao_Paulo", "America/Campo_Grande", "America/Cuiaba",
		"America/Santarem", "America/Porto_Velho", "America/Boa_Vista", "America/Manaus", "America/Eirunepe",
		"America/Rio_Branco"}
	fusosAR = []string{"America/Argentina/Buenos_Aires", "America/Argentina/Cordoba", "America/Argentina/Salta",
		"America/Argentina/Jujuy", "America/Argentina/Tucuman", "America/Argentina/Catamarca",
		"America/Argentina/La_Rioja", "America/Argentina/San_Juan", "America/Argentina/Mendoza",
		"America/Argentina/San_Luis", "America/Argentina/Rio_Gallegos", "America/Argentina/Ushuaia"}
	fusosPY = []string{"America/Asuncion"}
)

// FusoValido diz se o fuso IANA está na lista do país.
func FusoValido(p PerfilPais, fuso string) bool { return slices.Contains(p.Fusos, fuso) }

// perfis na ordem em que o contrato os lista. Não é alterado depois do
// carregamento do pacote.
var perfis = [...]PerfilPais{
	{Codigo: BR, Moeda: dinheiro.BRL, IdiomaPadrao: IdiomaPtBR, FusoPadrao: "America/Sao_Paulo", Fusos: fusosBR, TipoEmpresa: TipoCNPJ, TiposPessoa: []string{TipoCPF}, Subdivisoes: subdivisoesBR},
	{Codigo: AR, Moeda: dinheiro.ARS, IdiomaPadrao: IdiomaEs419, FusoPadrao: "America/Argentina/Buenos_Aires", Fusos: fusosAR, TipoEmpresa: TipoCUIT, TiposPessoa: []string{TipoCUIL, TipoDNI}, Subdivisoes: subdivisoesAR},
	{Codigo: PY, Moeda: dinheiro.PYG, IdiomaPadrao: IdiomaEs419, FusoPadrao: "America/Asuncion", Fusos: fusosPY, TipoEmpresa: TipoRUC, TiposPessoa: []string{TipoCI}, Subdivisoes: subdivisoesPY},
}

// Perfil devolve o perfil do país, ou falso se o país não é suportado.
func Perfil(codigo string) (PerfilPais, bool) {
	for _, p := range perfis {
		if p.Codigo == codigo {
			p.TiposPessoa = slices.Clone(p.TiposPessoa)
			p.Fusos = slices.Clone(p.Fusos)
			p.Subdivisoes = slices.Clone(p.Subdivisoes)
			return p, true
		}
	}
	return PerfilPais{}, false
}

// Codigos lista os países suportados, na ordem do contrato.
func Codigos() []string {
	l := make([]string, 0, len(perfis))
	for _, p := range perfis {
		l = append(l, p.Codigo)
	}
	return l
}

// IdiomaValido diz se a etiqueta é um idioma suportado.
func IdiomaValido(idioma string) bool {
	return idioma == IdiomaPtBR || idioma == IdiomaEs419
}

// TiposEmpresa lista os tipos de pessoa jurídica de todos os países: o que a
// unidade aceita em tese (a coerência com o país da empresa é do serviço).
func TiposEmpresa() []string {
	l := make([]string, 0, len(perfis))
	for _, p := range perfis {
		l = append(l, p.TipoEmpresa)
	}
	return l
}

// TiposPessoaDeTodos lista os tipos de pessoa física de todos os países
// suportados: a pessoa aceita qualquer um deles.
func TiposPessoaDeTodos() []string {
	var l []string
	for _, p := range perfis {
		l = append(l, p.TiposPessoa...)
	}
	return l
}

// PaisDaMoeda devolve o país cuja moeda é a informada (ck_empresa_moeda faz a
// relação ser de um para um), ou falso se nenhum perfil usa a moeda. Serve ao
// usuário, que guarda a moeda da empresa e não o país.
func PaisDaMoeda(moeda string) (string, bool) {
	for _, p := range perfis {
		if p.Moeda.Codigo() == moeda {
			return p.Codigo, true
		}
	}
	return "", false
}
