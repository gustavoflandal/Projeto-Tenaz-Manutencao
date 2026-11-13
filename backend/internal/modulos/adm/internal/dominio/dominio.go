// Package dominio tem as regras puras do módulo ADM: identificação por país,
// calendário da unidade, plano e as regras de perfis e permissões de usuário
// (RF-ADM-001, 002, 004, 005, 015). Sem IO: tudo aqui é testável sem banco.
package dominio

import (
	"errors"
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
	"tenaz/internal/plataforma/pais"
	"tenaz/internal/plataforma/textos"
)

// ------------------------------------------------------------ identificação

// Identificacao é o par tipo e valor já validado e normalizado (ADR-013): o
// valor em maiúsculas, só letras e dígitos.
type Identificacao struct {
	Tipo  string
	Valor string
}

// IdentificacaoDaEmpresa valida a identificação da empresa: obrigatória, e só
// do tipo de pessoa jurídica do país (cnpj, cuit ou ruc). prefixo é o do campo
// no corpo ("" na empresa, "unidade." na unidade inicial).
func IdentificacaoDaEmpresa(prefixo string, perfil pais.PerfilPais, tipo, valor string) (Identificacao, *erros.Erro) {
	id, ee := validarPar(prefixo, &tipo, &valor, []string{perfil.TipoEmpresa}, "identificacao.tipo_da_empresa",
		textos.Parametros{"tipo": strings.ToUpper(perfil.TipoEmpresa)}, true)
	if ee != nil {
		return Identificacao{}, ee
	}
	return *id, nil
}

// IdentificacaoDaUnidade valida o par da unidade, opcional: os dois ausentes
// ou os dois presentes, de um tipo de pessoa jurídica. A coerência com o país
// da empresa é de ConferirTipoDaUnidade, que precisa do país.
func IdentificacaoDaUnidade(prefixo string, tipo, valor *string) (*Identificacao, *erros.Erro) {
	return validarPar(prefixo, tipo, valor, pais.TiposEmpresa(), "identificacao.tipo_da_pessoa", nil, false)
}

// ConferirTipoDaUnidade confere que o tipo da unidade é o de pessoa jurídica
// do país da empresa: a unidade não tem país próprio (Q1).
func ConferirTipoDaUnidade(prefixo string, perfil pais.PerfilPais, id *Identificacao) *erros.Erro {
	if id == nil || id.Tipo == perfil.TipoEmpresa {
		return nil
	}
	return erros.EntradaInvalida.ComCampoChave(prefixo+"identificacao_tipo", erros.CampoValorInvalido,
		"identificacao.tipo_da_empresa", textos.Parametros{"tipo": strings.ToUpper(perfil.TipoEmpresa)}, perfil.TipoEmpresa)
}

// IdentificacaoDaPessoa valida o par da pessoa, opcional: aceita os tipos de
// pessoa física de qualquer país suportado (cpf, cuil, dni, ci), porque uma
// transportadora brasileira pode ter motorista paraguaio.
func IdentificacaoDaPessoa(prefixo string, tipo, valor *string) (*Identificacao, *erros.Erro) {
	return validarPar(prefixo, tipo, valor, pais.TiposPessoaDeTodos(), "identificacao.tipo_da_pessoa", nil, false)
}

// validarPar é a regra comum: tipo e valor andam juntos; o tipo é um dos
// aceitos; o valor passa pela forma e, onde há, pelo dígito verificador.
func validarPar(prefixo string, tipo, valor *string, aceitos []string, chaveTipo string, paramsTipo textos.Parametros,
	obrigatorio bool) (*Identificacao, *erros.Erro) {
	campoTipo, campoValor := prefixo+"identificacao_tipo", prefixo+"identificacao"
	t, v := textoOuVazio(tipo), textoOuVazio(valor)
	switch {
	case t == "" && v == "":
		if obrigatorio {
			return nil, erros.EntradaInvalida.
				ComCampoChave(campoTipo, erros.CampoObrigatorio, "identificacao.tipo_obrigatorio", nil).
				ComCampoChave(campoValor, erros.CampoObrigatorio, "identificacao.valor_obrigatorio", nil)
		}
		return nil, nil
	case t == "":
		return nil, erros.EntradaInvalida.ComCampoChave(campoTipo, erros.CampoObrigatorio, "identificacao.tipo_obrigatorio", nil)
	case v == "":
		return nil, erros.EntradaInvalida.ComCampoChave(campoValor, erros.CampoObrigatorio, "identificacao.valor_obrigatorio", nil)
	}
	if !slices.Contains(aceitos, t) {
		chave := chaveTipo
		if !slices.Contains(tiposConhecidos(), t) {
			chave, paramsTipo = "identificacao.tipo_desconhecido", nil
		}
		return nil, erros.EntradaInvalida.ComCampoChave(campoTipo, erros.CampoValorInvalido, chave, paramsTipo, aceitos...)
	}
	n, err := pais.ValidarIdentificacao(t, v)
	switch {
	case errors.Is(err, pais.ErrFormato):
		return nil, erros.EntradaInvalida.ComCampoChave(campoValor, erros.CampoFormatoInvalido, "identificacao.formato",
			textos.Parametros{"tipo": strings.ToUpper(t)})
	case errors.Is(err, pais.ErrDigitoVerificador):
		return nil, erros.EntradaInvalida.ComCampoChave(campoValor, erros.CampoFormatoInvalido, "identificacao.digito",
			textos.Parametros{"tipo": strings.ToUpper(t)})
	case err != nil:
		return nil, erros.EntradaInvalida.ComCampoChave(campoTipo, erros.CampoValorInvalido, "identificacao.tipo_desconhecido", nil, aceitos...)
	}
	return &Identificacao{Tipo: t, Valor: n}, nil
}

func tiposConhecidos() []string {
	return append(pais.TiposEmpresa(), pais.TiposPessoaDeTodos()...)
}

func textoOuVazio(s *string) string {
	if s == nil {
		return ""
	}
	return strings.TrimSpace(*s)
}

// MascararIdentificacao esconde a identificação pessoal na lista de pessoas; a
// completa sai só no detalhe, com a visualização gravada na trilha como
// sensível (seção 9.1 do documento 2).
func MascararIdentificacao(tipo, valor string) string {
	return pais.Mascarar(tipo, valor)
}

// ------------------------------------------------------------- país e idioma

// PerfilDoPais confere o país da criação da conta e devolve o perfil, de onde
// saem a moeda e o idioma padrão sugerido.
func PerfilDoPais(campo, codigo string) (pais.PerfilPais, *erros.Erro) {
	c := strings.TrimSpace(codigo)
	if c == "" {
		return pais.PerfilPais{}, erros.EntradaInvalida.ComCampoChave(campo, erros.CampoObrigatorio, "pais.obrigatorio", nil, pais.Codigos()...)
	}
	p, ok := pais.Perfil(c)
	if !ok {
		return pais.PerfilPais{}, erros.EntradaInvalida.ComCampoChave(campo, erros.CampoValorInvalido, "pais.desconhecido", nil, pais.Codigos()...)
	}
	return p, nil
}

// Idioma confere a etiqueta de idioma (pt-BR ou es-419).
func Idioma(campo, valor string) (string, *erros.Erro) {
	v := strings.TrimSpace(valor)
	if !pais.IdiomaValido(v) {
		return "", erros.EntradaInvalida.ComCampoChave(campo, erros.CampoValorInvalido, "pais.idioma_invalido", nil,
			pais.IdiomaPtBR, pais.IdiomaEs419)
	}
	return v, nil
}

// ------------------------------------------------------------ textos simples

// Texto confere obrigatório e tamanho máximo, e devolve o texto aparado.
func Texto(campo, valor string, max int) (string, *erros.Erro) {
	v := strings.TrimSpace(valor)
	switch {
	case v == "":
		return "", erros.EntradaInvalida.ComCampoChave(campo, erros.CampoObrigatorio, "campo.informe", textos.Parametros{"campo": campo})
	case utf8.RuneCountInString(v) > max:
		return "", erros.EntradaInvalida.ComCampoChave(campo, erros.CampoTamanhoMaximo, "campo.tamanho_maximo", textos.Parametros{"max": max})
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
		return "", erros.EntradaInvalida.ComCampoChave(campo, erros.CampoObrigatorio, "campo.email_obrigatorio", nil)
	}
	a, err := mail.ParseAddress(v)
	if err != nil || a.Address != v || len(v) > 254 {
		return "", erros.EntradaInvalida.ComCampoChave(campo, erros.CampoFormatoInvalido, "campo.email_invalido", nil)
	}
	return v, nil
}

// FusoHorario confere o nome IANA (America/Sao_Paulo).
func FusoHorario(campo, valor string) (string, *erros.Erro) {
	v := strings.TrimSpace(valor)
	if v == "" || v == "Local" {
		return "", erros.EntradaInvalida.ComCampoChave(campo, erros.CampoValorInvalido, "campo.fuso_obrigatorio", nil)
	}
	if _, err := time.LoadLocation(v); err != nil {
		return "", erros.EntradaInvalida.ComCampoChave(campo, erros.CampoValorInvalido, "campo.fuso_desconhecido", nil)
	}
	return v, nil
}

// FusoDoPais confere o fuso contra a lista de fusos do país da empresa (perfil
// de país): nome IANA válido E aceito pelo país. Erro valor_invalido com
// valores_aceitos.
func FusoDoPais(campo string, perfil pais.PerfilPais, valor string) (string, *erros.Erro) {
	v, ee := FusoHorario(campo, valor)
	if ee != nil {
		return "", ee
	}
	if !pais.FusoValido(perfil, v) {
		return "", erros.EntradaInvalida.ComCampoChave(campo, erros.CampoValorInvalido, "campo.fuso_fora_do_pais",
			textos.Parametros{"pais": perfil.Codigo}, perfil.Fusos...)
	}
	return v, nil
}

// SubdivisaoECodigoPostal confere o par do endereço contra o perfil do país da
// empresa (documento 3, 2.18). Cada um é opcional (vazio = não informado) e
// devolve o código postal normalizado. Erros em prefixo+"subdivisao" e
// prefixo+"codigo_postal": valor_invalido (com valores_aceitos, no primeiro),
// formato_invalido, e valor_invalido quando o CPA argentino de 8 caracteres não
// começa com a letra da subdivisão (regra da Tenaz).
func SubdivisaoECodigoPostal(prefixo string, perfil pais.PerfilPais, subdivisao, codigoPostal string) (string, *erros.Erro) {
	if subdivisao != "" && !pais.SubdivisaoValida(perfil, subdivisao) {
		return "", erros.EntradaInvalida.ComCampoChave(prefixo+"subdivisao", erros.CampoValorInvalido, "campo.subdivisao_fora_do_pais",
			textos.Parametros{"pais": perfil.Codigo}, perfil.Subdivisoes...)
	}
	if strings.TrimSpace(codigoPostal) == "" {
		return "", nil
	}
	n, err := pais.ValidarCodigoPostal(perfil, codigoPostal)
	if err != nil {
		return "", erros.EntradaInvalida.ComCampoChave(prefixo+"codigo_postal", erros.CampoFormatoInvalido, "campo.codigo_postal_formato",
			textos.Parametros{"pais": perfil.Codigo})
	}
	if !pais.CPAConfereComSubdivisao(subdivisao, n) {
		return "", erros.EntradaInvalida.ComCampoChave(prefixo+"codigo_postal", erros.CampoValorInvalido, "campo.codigo_postal_cpa_subdivisao",
			textos.Parametros{"subdivisao": subdivisao})
	}
	return n, nil
}

// Data confere AAAA-MM-DD.
func Data(campo, valor string) (string, *erros.Erro) {
	v := strings.TrimSpace(valor)
	if _, err := time.Parse(time.DateOnly, v); err != nil {
		return "", erros.EntradaInvalida.ComCampoChave(campo, erros.CampoFormatoInvalido, "campo.data_formato", nil)
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
			return erros.EntradaInvalida.ComCampoChave("plano", erros.CampoObrigatorio, "campo.plano_obrigatorio", nil)
		}
		return nil
	}
	switch *plano {
	case PlanoEssencial, PlanoProfissional, PlanoEnterprise:
		return nil
	}
	return erros.EntradaInvalida.ComCampoChave("plano", erros.CampoValorInvalido, "campo.plano_desconhecido", nil)
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
	return erros.PlanoNaoIncluiRecurso.ComChave("detalhe.plano_segunda_unidade", nil).
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
		return erros.EntradaInvalida.ComCampoChave(campo("dia_semana"), erros.CampoForaDaFaixa, "campo.dia_semana", nil)
	case !reHora.MatchString(t.HoraInicio):
		return erros.EntradaInvalida.ComCampoChave(campo("hora_inicio"), erros.CampoFormatoInvalido, "campo.hora_inicio", nil)
	case !reHora.MatchString(t.HoraFim) && t.HoraFim != "24:00":
		return erros.EntradaInvalida.ComCampoChave(campo("hora_fim"), erros.CampoFormatoInvalido, "campo.hora_fim", nil)
	case t.HoraFim == t.HoraInicio:
		return erros.EntradaInvalida.ComCampoChave(campo("hora_fim"), erros.CampoValorInvalido, "campo.hora_fim_igual", nil)
	}
	return nil
}

// ErrUltimoTurno: a unidade ficaria sem calendário. Vale para criar unidade
// sem turno, inclusive a inicial da conta, e para remover o último
// (RN-ADM-019). O sistema não impõe calendário: os modelos "comercial" e
// "contínuo" são da web e só preenchem o formulário.
func ErrUltimoTurno() *erros.Erro {
	return erros.EntradaInvalida.ComCampoChave("turnos", erros.CampoTamanhoMinimo, "campo.ultimo_turno", nil)
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
	return erros.EntradaInvalida.ComCampoChave(campo, erros.CampoValorInvalido, "campo.turno_sobrepoe", textos.Parametros{
		"dia": textos.Chave(fmt.Sprintf("dia.%d", existente.DiaSemana)), "inicio": existente.HoraInicio, "fim": existente.HoraFim})
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
		return erros.EntradaInvalida.ComCampoChave("locatario_id", erros.CampoValorInvalido, "campo.fornecedor_ou_locatario", nil)
	}
	if len(perfis) == 0 {
		return erros.EntradaInvalida.ComCampoChave("perfis", erros.CampoTamanhoMinimo, "campo.perfil_obrigatorio", nil)
	}
	vistos := map[string]bool{}
	for i, p := range perfis {
		if vistos[p] {
			return erros.EntradaInvalida.ComCampoChave(fmt.Sprintf("perfis[%d]", i), erros.CampoDuplicado, "campo.perfil_repetido", nil)
		}
		vistos[p] = true
	}
	switch {
	case v.Fornecedor:
		if !slices.Equal(perfis, []string{autorizacao.PerfilPrestador}) {
			return erros.EntradaInvalida.ComCampoChave("perfis", erros.CampoValorInvalido, "campo.perfil_fornecedor", nil)
		}
	case v.Locatario:
		if !slices.Equal(perfis, []string{autorizacao.PerfilLocatario}) {
			return erros.EntradaInvalida.ComCampoChave("perfis", erros.CampoValorInvalido, "campo.perfil_locatario", nil)
		}
	default:
		for i, p := range perfis {
			if autorizacao.PerfilExterno(p) {
				campo := "fornecedor_id"
				if p == autorizacao.PerfilLocatario {
					campo = "locatario_id"
				}
				return erros.EntradaInvalida.ComCampoChave(campo, erros.CampoObrigatorio, "campo.perfil_exige_vinculo", textos.Parametros{
					"perfil": p, "indice": i, "alvo": textos.Chave("alvo." + strings.TrimSuffix(campo, "_id"))})
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
		return erros.SemPermissao.ComChave("detalhe.ninguem_altera_proprios_perfis", nil)
	}
	_, ganhaAdm := perfisNovos[autorizacao.PerfilAdministrador]
	if (alvoEhAdministrador || ganhaAdm) && !ator.TemPerfil(autorizacao.PerfilAdministrador) {
		return erros.SemPermissao.ComChave("detalhe.so_administrador_altera", nil)
	}
	codigos := make([]string, 0, len(perfisNovos))
	for c := range perfisNovos {
		codigos = append(codigos, c)
	}
	slices.Sort(codigos)
	for _, c := range codigos {
		if ok, falta := ator.PodeAtribuirPerfil(perfisNovos[c]); !ok {
			return erros.SemPermissao.ComChave("detalhe.perfil_da_permissao", textos.Parametros{"perfil": c, "permissao": falta})
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
		return erros.EntradaInvalida.ComCampoChave("permissao", erros.CampoValorInvalido, "campo.permissao_fora_catalogo", nil)
	}
	if alvoID.String() == ator.UsuarioID.String() {
		return erros.SemPermissao.ComChave("detalhe.ninguem_concede_a_si", nil)
	}
	if autorizacao.IsencaoAdministrador(permissao) && !ator.TemPerfil(autorizacao.PerfilAdministrador) {
		return erros.SemPermissao.ComChave("detalhe.so_admin_concede", textos.Parametros{"permissao": permissao})
	}
	if !ator.PodeConceder(permissao) {
		return erros.SemPermissao.ComChave("detalhe.sem_a_permissao", textos.Parametros{"permissao": permissao})
	}
	if perfilExterno != "" {
		if autorizacao.VedadaExterno(permissao) {
			return erros.EntradaInvalida.ComCampoChave("permissao", erros.CampoValorInvalido, "campo.permissao_externo", nil)
		}
		if !slices.Contains(autorizacao.PermissoesDoPerfil(perfilExterno), permissao) {
			return erros.EntradaInvalida.ComCampoChave("permissao", erros.CampoValorInvalido, "campo.permissao_do_perfil",
				textos.Parametros{"perfil": perfilExterno})
		}
	}
	return nil
}
