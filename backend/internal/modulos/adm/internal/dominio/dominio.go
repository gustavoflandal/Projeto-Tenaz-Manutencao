// Package dominio tem as regras puras do módulo ADM: documentos (CPF e CNPJ),
// calendário da unidade, plano e as regras de perfis e permissões de usuário
// (RF-ADM-001, 002, 004, 005, 015). Sem IO: tudo aqui é testável sem banco.
package dominio

import (
	"fmt"
	"net/mail"
	"regexp"
	"slices"
	"strings"
	"time"
	_ "time/tzdata" // fuso horário validado igual em qualquer máquina, inclusive Windows e imagem mínima
	"unicode/utf8"

	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/erros"
)

// --------------------------------------------------------------- documentos

// SoDigitos tira pontos, traços, barras e espaços.
func SoDigitos(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}

func repetido(d string) bool {
	return strings.Count(d, d[:1]) == len(d)
}

// CPFValido confere os dois dígitos verificadores (11 dígitos, sem máscara).
func CPFValido(cpf string) bool {
	if len(cpf) != 11 || SoDigitos(cpf) != cpf || repetido(cpf) {
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

// CNPJValido confere os dois dígitos verificadores (14 dígitos, sem máscara).
func CNPJValido(cnpj string) bool {
	if len(cnpj) != 14 || SoDigitos(cnpj) != cnpj || repetido(cnpj) {
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

// MascararCPF esconde o começo e o fim: ***.456.789-**. A lista de pessoas
// mostra só isto; o CPF inteiro sai no detalhe, com registro na trilha
// (seção 9.1 do documento 2, visualização de dado sensível).
func MascararCPF(cpf string) string {
	if len(cpf) != 11 {
		return ""
	}
	return "***." + cpf[3:6] + "." + cpf[6:9] + "-**"
}

// NormalizarCPF devolve o CPF só com dígitos ou o erro de campo.
func NormalizarCPF(campo, cpf string) (string, *erros.Erro) {
	d := SoDigitos(cpf)
	if !CPFValido(d) {
		return "", erros.EntradaInvalida.ComCampo(campo, erros.CampoFormatoInvalido, "CPF inválido.")
	}
	return d, nil
}

// NormalizarCNPJ devolve o CNPJ só com dígitos ou o erro de campo.
func NormalizarCNPJ(campo, cnpj string) (string, *erros.Erro) {
	d := SoDigitos(cnpj)
	if !CNPJValido(d) {
		return "", erros.EntradaInvalida.ComCampo(campo, erros.CampoFormatoInvalido, "CNPJ inválido.")
	}
	return d, nil
}

// ------------------------------------------------------------ textos simples

// Texto confere obrigatório e tamanho máximo, e devolve o texto aparado.
func Texto(campo, valor string, max int) (string, *erros.Erro) {
	v := strings.TrimSpace(valor)
	switch {
	case v == "":
		return "", erros.EntradaInvalida.ComCampo(campo, erros.CampoObrigatorio, "Informe "+campo+".")
	case utf8.RuneCountInString(v) > max:
		return "", erros.EntradaInvalida.ComCampo(campo, erros.CampoTamanhoMaximo, fmt.Sprintf("Use no máximo %d caracteres.", max))
	}
	return v, nil
}

// TextoOpcional apara; vazio vira nulo.
func TextoOpcional(campo string, valor *string, max int) (*string, *erros.Erro) {
	if valor == nil || strings.TrimSpace(*valor) == "" {
		return nil, nil
	}
	v, e := Texto(campo, *valor, max)
	if e != nil {
		return nil, e
	}
	return &v, nil
}

// Email confere o formato e devolve em minúsculas.
func Email(campo, valor string) (string, *erros.Erro) {
	v := strings.ToLower(strings.TrimSpace(valor))
	if v == "" {
		return "", erros.EntradaInvalida.ComCampo(campo, erros.CampoObrigatorio, "Informe o e-mail.")
	}
	a, err := mail.ParseAddress(v)
	if err != nil || a.Address != v || len(v) > 254 {
		return "", erros.EntradaInvalida.ComCampo(campo, erros.CampoFormatoInvalido, "E-mail inválido.")
	}
	return v, nil
}

// FusoHorario confere o nome IANA (America/Sao_Paulo).
func FusoHorario(campo, valor string) (string, *erros.Erro) {
	v := strings.TrimSpace(valor)
	if v == "" || v == "Local" {
		return "", erros.EntradaInvalida.ComCampo(campo, erros.CampoValorInvalido, "Informe um fuso horário IANA, como America/Sao_Paulo.")
	}
	if _, err := time.LoadLocation(v); err != nil {
		return "", erros.EntradaInvalida.ComCampo(campo, erros.CampoValorInvalido, "Fuso horário desconhecido.")
	}
	return v, nil
}

// Data confere AAAA-MM-DD.
func Data(campo, valor string) (string, *erros.Erro) {
	v := strings.TrimSpace(valor)
	if _, err := time.Parse(time.DateOnly, v); err != nil {
		return "", erros.EntradaInvalida.ComCampo(campo, erros.CampoFormatoInvalido, "Use a data no formato AAAA-MM-DD.")
	}
	return v, nil
}

// ------------------------------------------------------------------- plano

// Planos da empresa (ck_empresa_plano). Nulo: conta só do plano Serviços.
const (
	PlanoEssencial    = "essencial"
	PlanoProfissional = "profissional"
	PlanoEnterprise   = "enterprise"
)

// PlanoValido confere o plano e a regra de contratação: sem plano principal,
// a conta precisa do SRV (ck_empresa_contratacao; sumário 2.2).
func PlanoValido(plano *string, srv bool) *erros.Erro {
	if plano == nil {
		if !srv {
			return erros.EntradaInvalida.ComCampo("plano", erros.CampoObrigatorio,
				"Escolha um plano ou habilite o módulo SRV (plano Serviços).")
		}
		return nil
	}
	switch *plano {
	case PlanoEssencial, PlanoProfissional, PlanoEnterprise:
		return nil
	}
	return erros.EntradaInvalida.ComCampo("plano", erros.CampoValorInvalido, "Plano desconhecido.")
}

// PodeCriarUnidade aplica "o Essencial tem uma unidade; a segunda exige o
// Profissional" (documento 3, seção 6.1; RF-ADM-002). A conta só do plano
// Serviços segue o Essencial (proposta: o sumário 2.2 não lista multi-unidade
// no plano Serviços). unidadesExistentes conta as não excluídas, inclusive as
// inativas.
func PodeCriarUnidade(plano *string, unidadesExistentes int64) *erros.Erro {
	if unidadesExistentes < 1 {
		return nil
	}
	if plano != nil && (*plano == PlanoProfissional || *plano == PlanoEnterprise) {
		return nil
	}
	return erros.PlanoNaoIncluiRecurso.Com("o plano tem uma unidade; a segunda exige o Profissional").
		ComDado("plano_necessario", PlanoProfissional)
}

// --------------------------------------------------------------- calendário

var reHora = regexp.MustCompile(`^([01][0-9]|2[0-3]):[0-5][0-9]$`)

// Turno é um turno de trabalho da unidade (turno_unidade).
type Turno struct {
	DiaSemana  int32 // ISO 8601: 1 = segunda ... 7 = domingo
	HoraInicio string
	HoraFim    string // menor que o início: termina no dia seguinte; "24:00" = meia-noite
}

// ValidarTurno confere dia e horários, com o prefixo de campo dado
// ("turnos[0]" na criação da unidade, vazio na rota do turno).
func ValidarTurno(prefixo string, t Turno) *erros.Erro {
	campo := func(c string) string {
		if prefixo == "" {
			return c
		}
		return prefixo + "." + c
	}
	switch {
	case t.DiaSemana < 1 || t.DiaSemana > 7:
		return erros.EntradaInvalida.ComCampo(campo("dia_semana"), erros.CampoForaDaFaixa, "Dia da semana de 1 (segunda) a 7 (domingo).")
	case !reHora.MatchString(t.HoraInicio):
		return erros.EntradaInvalida.ComCampo(campo("hora_inicio"), erros.CampoFormatoInvalido, "Use HH:MM, de 00:00 a 23:59.")
	case !reHora.MatchString(t.HoraFim) && t.HoraFim != "24:00":
		return erros.EntradaInvalida.ComCampo(campo("hora_fim"), erros.CampoFormatoInvalido, "Use HH:MM, de 00:00 a 24:00.")
	case t.HoraFim == t.HoraInicio:
		return erros.EntradaInvalida.ComCampo(campo("hora_fim"), erros.CampoValorInvalido, "O fim do turno tem de ser diferente do início.")
	}
	return nil
}

// ErrUltimoTurno: a unidade ficaria sem calendário. Vale para criar unidade
// sem turno, inclusive a inicial da conta, e para remover o último
// (RN-ADM-019). O sistema não impõe calendário: os modelos "comercial" e
// "contínuo" são da web e só preenchem o formulário.
func ErrUltimoTurno() *erros.Erro {
	return erros.EntradaInvalida.ComCampo("turnos", erros.CampoTamanhoMinimo,
		"A unidade precisa de pelo menos um turno de trabalho (RN-ADM-019).")
}

// minutosSemana é a linha da semana: segunda 00:00 é 0, domingo 24:00 é 10080.
const minutosSemana = 7 * 24 * 60

func minutos(hhmm string) int {
	if hhmm == "24:00" {
		return 24 * 60
	}
	var h, m int
	_, _ = fmt.Sscanf(hhmm, "%d:%d", &h, &m)
	return h*60 + m
}

// intervalo é um trecho [inicio, fim) da linha da semana, em minutos.
type intervalo struct{ inicio, fim int }

// intervalos põe o turno na linha da semana. O turno que atravessa a
// meia-noite invade o dia seguinte; o de domingo invade a segunda, e por isso
// vira dois trechos: o fim do domingo e o começo da segunda.
func (t Turno) intervalos() []intervalo {
	base := int(t.DiaSemana-1) * 24 * 60
	ini, fim := minutos(t.HoraInicio), minutos(t.HoraFim)
	if fim <= ini {
		fim += 24 * 60 // termina no dia seguinte
	}
	a, b := base+ini, base+fim
	if b <= minutosSemana {
		return []intervalo{{a, b}}
	}
	return []intervalo{{a, minutosSemana}, {0, b - minutosSemana}}
}

// Sobrepoem diz se dois turnos se sobrepõem na semana. Turnos encostados (o
// fim de um igual ao início do outro) não se sobrepõem.
func Sobrepoem(x, y Turno) bool {
	for _, a := range x.intervalos() {
		for _, b := range y.intervalos() {
			if a.inicio < b.fim && b.inicio < a.fim {
				return true
			}
		}
	}
	return false
}

// PrimeiraSobreposicao procura dois turnos que se sobrepõem e devolve os
// índices, ou -1, -1.
func PrimeiraSobreposicao(turnos []Turno) (int, int) {
	for i := range turnos {
		for j := i + 1; j < len(turnos); j++ {
			if Sobrepoem(turnos[i], turnos[j]) {
				return i, j
			}
		}
	}
	return -1, -1
}

var nomesDia = [...]string{"segunda", "terça", "quarta", "quinta", "sexta", "sábado", "domingo"}

// ErrSobreposicao descreve o turno existente que conflita (RN-ADM-019).
func ErrSobreposicao(campo string, existente Turno) *erros.Erro {
	return erros.EntradaInvalida.ComCampo(campo, erros.CampoValorInvalido, fmt.Sprintf(
		"O turno se sobrepõe ao de %s, das %s às %s.", nomesDia[existente.DiaSemana-1], existente.HoraInicio, existente.HoraFim))
}

// --------------------------------------------------- usuários e perfis

// Situações do usuário (ck_usuario_situacao).
const (
	UsuarioConvidado = "convidado"
	UsuarioAtivo     = "ativo"
	UsuarioInativo   = "inativo"
)

// Vinculo externo do usuário: fornecedor (prestador) ou locatário.
type Vinculo struct {
	Fornecedor bool
	Locatario  bool
}

// ValidarPerfisDoUsuario aplica as regras de RF-ADM-004 sobre o conjunto de
// perfis (pelo código de origem): pelo menos um; usuário com fornecedor tem
// só o perfil prestador; com locatário, só o locatario; interno nunca recebe
// perfil externo.
func ValidarPerfisDoUsuario(perfis []string, v Vinculo) *erros.Erro {
	if v.Fornecedor && v.Locatario {
		return erros.EntradaInvalida.ComCampo("locatario_id", erros.CampoValorInvalido,
			"O usuário é de um fornecedor ou de um locatário, não dos dois.")
	}
	if len(perfis) == 0 {
		return erros.EntradaInvalida.ComCampo("perfis", erros.CampoTamanhoMinimo, "O usuário precisa de pelo menos um perfil.")
	}
	vistos := map[string]bool{}
	for i, p := range perfis {
		if vistos[p] {
			return erros.EntradaInvalida.ComCampo(fmt.Sprintf("perfis[%d]", i), erros.CampoDuplicado, "Perfil repetido.")
		}
		vistos[p] = true
	}
	switch {
	case v.Fornecedor:
		if !slices.Equal(perfis, []string{autorizacao.PerfilPrestador}) {
			return erros.EntradaInvalida.ComCampo("perfis", erros.CampoValorInvalido,
				"Usuário de fornecedor tem só o perfil prestador.")
		}
	case v.Locatario:
		if !slices.Equal(perfis, []string{autorizacao.PerfilLocatario}) {
			return erros.EntradaInvalida.ComCampo("perfis", erros.CampoValorInvalido,
				"Usuário de locatário tem só o perfil locatario.")
		}
	default:
		for i, p := range perfis {
			if autorizacao.PerfilExterno(p) {
				campo := "fornecedor_id"
				if p == autorizacao.PerfilLocatario {
					campo = "locatario_id"
				}
				return erros.EntradaInvalida.ComCampo(campo, erros.CampoObrigatorio,
					fmt.Sprintf("O perfil %s (perfis[%d]) exige o vínculo com %s.", p, i, strings.TrimSuffix(campo, "_id")))
			}
		}
	}
	return nil
}

// ConferirAtribuicao aplica as regras anti-escalada (documento 9, seção 5.4)
// à troca de perfis de um usuário: ninguém altera os próprios perfis; só
// administrador mexe em administrador; o editor precisa ter as permissões dos
// perfis que atribui (com as isenções do administrador).
func ConferirAtribuicao(ator autorizacao.Ator, alvoID fmt.Stringer, alvoEhAdministrador bool, perfisNovos map[string][]string) *erros.Erro {
	if alvoID.String() == ator.UsuarioID.String() {
		return erros.SemPermissao.Com("ninguém altera os próprios perfis")
	}
	_, ganhaAdm := perfisNovos[autorizacao.PerfilAdministrador]
	if (alvoEhAdministrador || ganhaAdm) && !ator.TemPerfil(autorizacao.PerfilAdministrador) {
		return erros.SemPermissao.Com("só administrador altera administrador")
	}
	codigos := make([]string, 0, len(perfisNovos))
	for c := range perfisNovos {
		codigos = append(codigos, c)
	}
	slices.Sort(codigos)
	for _, c := range codigos {
		if ok, falta := ator.PodeAtribuirPerfil(perfisNovos[c]); !ok {
			return erros.SemPermissao.Com("o perfil %s dá a permissão %s, que você não tem", c, falta)
		}
	}
	return nil
}

// ConferirConcessao aplica as regras da permissão avulsa (RF-ADM-005;
// documento 9, seção 5.4). alvoExterno: usuário prestador ou locatário.
//
// perfilExterno é o tipo (perfil de origem) do usuário externo alvo,
// prestador ou locatario; vazio para usuário interno. O externo recebe
// avulsa só dentro do conjunto padrão do seu perfil, que é fechado
// (RN-ADM-026), e nunca uma permissão vedada a externo.
func ConferirConcessao(ator autorizacao.Ator, alvoID fmt.Stringer, permissao string, perfilExterno string) *erros.Erro {
	if !autorizacao.Existe(permissao) {
		return erros.EntradaInvalida.ComCampo("permissao", erros.CampoValorInvalido, "Permissão fora do catálogo de RF-ADM-005.")
	}
	if alvoID.String() == ator.UsuarioID.String() {
		return erros.SemPermissao.Com("ninguém concede permissão a si mesmo")
	}
	if autorizacao.IsencaoAdministrador(permissao) && !ator.TemPerfil(autorizacao.PerfilAdministrador) {
		return erros.SemPermissao.Com("só o administrador concede %s", permissao)
	}
	if !ator.PodeConceder(permissao) {
		return erros.SemPermissao.Com("você não tem a permissão %s", permissao)
	}
	if perfilExterno != "" {
		if autorizacao.VedadaExterno(permissao) {
			return erros.EntradaInvalida.ComCampo("permissao", erros.CampoValorInvalido,
				"Esta permissão nunca é concedida a usuário prestador ou locatário.")
		}
		if !slices.Contains(autorizacao.PermissoesDoPerfil(perfilExterno), permissao) {
			return erros.EntradaInvalida.ComCampo("permissao", erros.CampoValorInvalido,
				"Usuário "+perfilExterno+" só recebe permissões do conjunto padrão do seu perfil.")
		}
	}
	return nil
}
