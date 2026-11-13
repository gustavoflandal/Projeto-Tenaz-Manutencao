// Package autorizacao é a fonte única, no código, do catálogo de permissões e
// da matriz dos 12 perfis padrão de RF-ADM-005, e da verificação de permissão
// a cada requisição (seção 4.8 do documento 3; seção 5.4 do documento 9).
//
// As permissões não vão no token: são lidas do banco a cada requisição
// (perfis do usuário + permissões avulsas de usuario_permissao), para que uma
// revogação valha na hora. Este pacote só guarda o resultado dessa leitura,
// no Ator, e responde se ele pode.
package autorizacao

import (
	"slices"
	"strings"

	"github.com/google/uuid"

	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/erros"
)

// Fase em que a permissão passa a ter efeito (coluna Fase do catálogo).
type Fase string

// Fases do catálogo.
const (
	FaseMVP Fase = "MVP"
	Fase2   Fase = "F2"
	Fase3   Fase = "F3"
)

// Permissao é uma linha do catálogo de RF-ADM-005.
type Permissao struct {
	Codigo    string // <módulo>.<ação>; pessoa.saude.ler é a exceção aceita
	Descricao string
	Grupo     string // grupo no documento 3, seção 4.8
	Fase      Fase
	Sensivel  bool // concessão e retirada destacadas na tela e auditadas
}

// Códigos dos perfis padrão (sumário 3.1). Também são o valor de
// perfil.codigo e de perfil.perfil_origem.
const (
	PerfilAdministrador = "administrador"
	PerfilGestor        = "gestor"
	PerfilPlanejador    = "planejador"
	PerfilTecnico       = "tecnico"
	PerfilAlmoxarife    = "almoxarife"
	PerfilComprador     = "comprador"
	PerfilPLH           = "plh"
	PerfilMotorista     = "motorista"
	PerfilSolicitante   = "solicitante"
	PerfilVisualizador  = "visualizador"
	PerfilPrestador     = "prestador"
	PerfilLocatario     = "locatario"
)

// Permissões citadas pelo código. Permissão usada em regra tem constante, para
// que um erro de digitação seja erro de compilação.
const (
	AdmLer                   = "adm.ler"
	AdmGerenciarEmpresa      = "adm.gerenciar_empresa"
	AdmLerPessoas            = "adm.ler_pessoas"
	AdmGerenciarPessoas      = "adm.gerenciar_pessoas"
	AdmLerUsuarios           = "adm.ler_usuarios"
	AdmGerenciarUsuarios     = "adm.gerenciar_usuarios"
	AdmGerenciarPerfis       = "adm.gerenciar_perfis"
	AdmLerAuditoria          = "adm.ler_auditoria"
	OsEncerrar               = "os.encerrar"
	EstDecidirDestino        = "est.decidir_destino"
	PessoaSaudeLer           = "pessoa.saude.ler"
	TerManterCatalogoProprio = "ter.manter_catalogo_proprio"
	MobUsarApp               = "mob.usar_app"
	SrvGerenciarConexao      = "srv.gerenciar_conexao"
	AdmGerenciarRegras       = "adm.gerenciar_regras"
	AdmGerenciarConfiguracao = "adm.gerenciar_configuracao"
	AtvColocarEmOperacao     = "atv.colocar_em_operacao"
	PrjGerenciar             = "prj.gerenciar"
	PrjDispensarEtapa        = "prj.dispensar_etapa"
	PrjAjustarMedicao        = "prj.ajustar_medicao"
	CtrGerenciar             = "ctr.gerenciar"
	CtrFecharMedicao         = "ctr.fechar_medicao"
	CtrReabrirMedicao        = "ctr.reabrir_medicao"
	permissaoExcecaoFormat   = PessoaSaudeLer
)

// PerfilPadrao descreve um dos 12 perfis padrão.
type PerfilPadrao struct {
	Codigo  string
	Nome    string
	Externo bool // prestador e locatario: escopo externo, alcance da seção 4.9 do documento 3
}

// perfisPadrao na ordem da matriz de RF-ADM-005.
var perfisPadrao = [...]PerfilPadrao{
	{PerfilAdministrador, "Administrador", false},
	{PerfilGestor, "Gestor", false},
	{PerfilPlanejador, "Planejador", false},
	{PerfilTecnico, "Técnico", false},
	{PerfilAlmoxarife, "Almoxarife", false},
	{PerfilComprador, "Comprador", false},
	{PerfilPLH, "PLH", false},
	{PerfilMotorista, "Motorista", false},
	{PerfilSolicitante, "Solicitante", false},
	{PerfilVisualizador, "Visualizador", false},
	{PerfilPrestador, "Prestador", true},
	{PerfilLocatario, "Locatário", true},
}

// isencoesAdministrador são as três permissões internas que o administrador
// não tem por padrão e mesmo assim concede a OUTROS usuários, nunca a si
// mesmo (RF-ADM-005; documento 9, seção 5.4, regras anti-escalada).
var isencoesAdministrador = [...]string{OsEncerrar, EstDecidirDestino, PessoaSaudeLer}

// configuracaoPRJCTR são as nove permissões da configuração, de PRJ e de CTR
// do sumário 3.6, todas vedadas a usuário externo.
var configuracaoPRJCTR = [...]string{
	AdmGerenciarRegras, AdmGerenciarConfiguracao, AtvColocarEmOperacao,
	PrjGerenciar, PrjDispensarEtapa, PrjAjustarMedicao,
	CtrGerenciar, CtrFecharMedicao, CtrReabrirMedicao,
}

// Catalogo devolve uma cópia do catálogo, na ordem do documento.
func Catalogo() []Permissao {
	return slices.Clone(catalogo[:])
}

// Buscar devolve a permissão do catálogo pelo código.
func Buscar(codigo string) (Permissao, bool) {
	for _, p := range catalogo {
		if p.Codigo == codigo {
			return p, true
		}
	}
	return Permissao{}, false
}

// Existe informa se o código está no catálogo.
func Existe(codigo string) bool {
	_, ok := Buscar(codigo)
	return ok
}

// PerfisPadrao devolve os 12 perfis padrão, na ordem da matriz.
func PerfisPadrao() []PerfilPadrao {
	return slices.Clone(perfisPadrao[:])
}

// BuscarPerfilPadrao devolve o perfil padrão pelo código.
func BuscarPerfilPadrao(codigo string) (PerfilPadrao, bool) {
	for _, p := range perfisPadrao {
		if p.Codigo == codigo {
			return p, true
		}
	}
	return PerfilPadrao{}, false
}

// PerfilExterno informa se o código (de perfil padrão ou de origem de um
// personalizado) é de escopo externo.
func PerfilExterno(codigo string) bool {
	p, ok := BuscarPerfilPadrao(codigo)
	return ok && p.Externo
}

// PermissoesDoPerfil devolve as permissões do perfil padrão, na ordem do
// catálogo. Perfil desconhecido: nil.
func PermissoesDoPerfil(codigo string) []string {
	return slices.Clone(matriz[codigo])
}

// VedadasExterno são as permissões que usuário externo nunca recebe, nem por
// perfil nem avulsa: os.encerrar, as nove da configuração, de PRJ e de CTR,
// srv.gerenciar_conexao, pessoa.saude.ler e as do grupo Admin (RF-ADM-005,
// regras dos perfis padrão).
//
// O banco repete esta lista em app.permissoes_vedadas_externo() para os CHECKs
// e gatilhos; TestVedadasExternoIgualAoBanco confere que as duas são iguais
// (pendência E01-10 do documento 2).
func VedadasExterno() []string {
	var v []string
	for _, p := range catalogo {
		if p.Codigo == OsEncerrar || p.Codigo == PessoaSaudeLer || p.Codigo == SrvGerenciarConexao ||
			slices.Contains(configuracaoPRJCTR[:], p.Codigo) || strings.HasPrefix(p.Grupo, "Admin") {
			v = append(v, p.Codigo)
		}
	}
	return v
}

// VedadaExterno informa se a permissão é vedada a usuário externo.
func VedadaExterno(codigo string) bool {
	return slices.Contains(VedadasExterno(), codigo)
}

// IsencaoAdministrador informa se a permissão é uma das três que o
// administrador concede a outros sem tê-las.
func IsencaoAdministrador(codigo string) bool {
	return slices.Contains(isencoesAdministrador[:], codigo)
}

// Conjunto é um conjunto de permissões efetivas.
type Conjunto map[string]struct{}

// NovoConjunto monta um conjunto com os códigos dados.
func NovoConjunto(codigos ...string) Conjunto {
	c := make(Conjunto, len(codigos))
	for _, p := range codigos {
		c[p] = struct{}{}
	}
	return c
}

// Tem informa se a permissão está no conjunto.
func (c Conjunto) Tem(codigo string) bool {
	_, ok := c[codigo]
	return ok
}

// Lista devolve os códigos em ordem alfabética.
func (c Conjunto) Lista() []string {
	l := make([]string, 0, len(c))
	for p := range c {
		l = append(l, p)
	}
	slices.Sort(l)
	return l
}

// Ator é quem faz a requisição, resolvido a partir da credencial e do banco a
// cada requisição. Nenhum campo vem do corpo, da URL nem do token sem
// conferência no servidor.
type Ator struct {
	EmpresaID    uuid.UUID
	UsuarioID    uuid.UUID
	IdentidadeID uuid.UUID
	SessaoID     uuid.UUID
	Canal        string // web ou app

	// Escopo de acesso com o fornecedor ou locatário do cadastro e
	// PodeLerSaude derivado das permissões efetivas.
	Escopo db.Escopo

	PessoaID   *uuid.UUID
	Perfis     []string // códigos de origem dos perfis ativos do usuário
	Permissoes Conjunto // perfis + permissões avulsas

	Plano         *string // nulo: conta só do plano Serviços
	SrvHabilitado bool

	// MFAObrigatorio: a empresa exige segundo fator (empresa.mfa_obrigatorio).
	MFAObrigatorio bool
	// SomenteAtivarMFA: a empresa exige MFA e a identidade ainda não o
	// ativou. A sessão serve só para ativá-lo (RN-ADM-027): as demais rotas e
	// o app respondem como não autenticado. Preenchido pela autenticação, que
	// lê a identidade.
	SomenteAtivarMFA bool
}

// Pode informa se o ator tem a permissão.
func (a Ator) Pode(codigo string) bool {
	return a.Permissoes.Tem(codigo)
}

// TemPerfil informa se o ator tem o perfil (pelo código de origem).
func (a Ator) TemPerfil(codigo string) bool {
	return slices.Contains(a.Perfis, codigo)
}

// Externo informa se o ator está em escopo prestador ou locatario.
func (a Ator) Externo() bool {
	return a.Escopo.Tipo == db.EscopoPrestador || a.Escopo.Tipo == db.EscopoLocatario
}

// Exigir confere que o ator tem todas as permissões.
//
// Usuário externo sem a permissão recebe recurso_nao_encontrado, e não permissao_negada:
// os perfis externos só acessam as rotas em que aparecem nomeados, e qualquer
// outra responde 404 para não revelar a estrutura da empresa contratante
// (documento 3, seção 4.8).
func (a Ator) Exigir(codigos ...string) error {
	for _, c := range codigos {
		if !a.Pode(c) {
			if a.Externo() {
				return erros.NaoEncontrado
			}
			return erros.SemPermissao.Com("falta a permissão %s", c)
		}
	}
	return nil
}

// ExigirAlguma confere que o ator tem pelo menos uma das permissões.
func (a Ator) ExigirAlguma(codigos ...string) error {
	for _, c := range codigos {
		if a.Pode(c) {
			return nil
		}
	}
	if a.Externo() {
		return erros.NaoEncontrado
	}
	return erros.SemPermissao.Com("falta uma das permissões %s", strings.Join(codigos, ", "))
}

// PodeConceder aplica a regra anti-escalada (documento 9, seção 5.4): só
// concede uma permissão quem a tem, com a exceção do administrador, que
// concede as três permissões internas que não tem por padrão. Não confere se
// o alvo é o próprio ator; isso é do caso de uso.
func (a Ator) PodeConceder(codigo string) bool {
	if a.Pode(codigo) {
		return true
	}
	return a.TemPerfil(PerfilAdministrador) && IsencaoAdministrador(codigo)
}

// PodeAtribuirPerfil aplica a regra anti-escalada a um perfil inteiro: o
// editor precisa ter todas as permissões do perfil. O administrador é isento
// das três permissões internas e de ter.manter_catalogo_proprio, exclusiva do
// prestador, que ele não tem por padrão e que acompanham os perfis gestor,
// planejador, almoxarife e prestador (RF-ADM-005). Devolve a primeira
// permissão que falta.
func (a Ator) PodeAtribuirPerfil(permissoesDoPerfil []string) (bool, string) {
	for _, p := range permissoesDoPerfil {
		if a.Pode(p) {
			continue
		}
		if a.TemPerfil(PerfilAdministrador) && (IsencaoAdministrador(p) || p == TerManterCatalogoProprio) {
			continue
		}
		return false, p
	}
	return true, ""
}

// FormatoValido confere o formato <módulo>.<ação> (sumário 3.6), com a
// exceção aceita de pessoa.saude.ler. Módulo de 2 ou 3 letras minúsculas;
// ação começando por letra, com letras, dígitos e sublinhado.
func FormatoValido(codigo string) bool {
	if codigo == permissaoExcecaoFormat {
		return true
	}
	modulo, acao, ok := strings.Cut(codigo, ".")
	if !ok || len(modulo) < 2 || len(modulo) > 3 || len(acao) < 2 {
		return false
	}
	for _, r := range modulo {
		if r < 'a' || r > 'z' {
			return false
		}
	}
	// A ação aceita dígito depois da primeira letra: con.gerenciar_nr13 está
	// no catálogo.
	for i, r := range acao {
		minuscula := r >= 'a' && r <= 'z'
		digito := r >= '0' && r <= '9'
		switch {
		case i == 0 && !minuscula:
			return false
		case r == '_' && i == len(acao)-1:
			return false
		case !minuscula && !digito && r != '_':
			return false
		}
	}
	return true
}
