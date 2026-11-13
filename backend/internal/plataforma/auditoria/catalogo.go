package auditoria

import (
	"slices"
	"strings"
)

// Tela é uma tela do sistema que pode ter captura configurada (RF-ADM-011).
//
// O catálogo é uma lista fechada no Go, versionada com a web e a API, e não
// uma tabela: tela é conceito de interface e muda junto com o release
// (documento 2, seção 9.2.2). O banco só confere o formato
// '<modulo>.<tela>'; o Go recusa a tela que não está aqui.
type Tela struct {
	Codigo string // '<modulo>.<tela>', prefixo = código do módulo (sumário, seção 2)
	Modulo string
	Nome   string // nome de exibição
}

// telas: as de E01 (rotas da web de Administração e do painel) e as previstas
// nas etapas seguintes, que já têm código reservado para que a captura
// configurada nunca precise ser migrada. Ordem: a da tela de configuração, Administração primeiro.
var telas = []Tela{
	{"adm.empresa", "adm", "Empresa"},
	{"adm.unidades", "adm", "Unidades"},
	{"adm.unidade_detalhe", "adm", "Detalhe da unidade"},
	{"adm.pessoas", "adm", "Pessoas"},
	{"adm.pessoa_detalhe", "adm", "Detalhe da pessoa"},
	{"adm.usuarios", "adm", "Usuários"},
	{"adm.usuario_detalhe", "adm", "Detalhe do usuário"},
	{"adm.perfis", "adm", "Perfis e permissões"},
	{"adm.auditoria", "adm", "Trilha de auditoria"},
	{"adm.parametros", "adm", "Parâmetros"},
	{"ind.painel", "ind", "Painel"},
	{"atv.ativos", "atv", "Ativos"},
	{"atv.ativo_detalhe", "atv", "Detalhe do ativo"},
	{"med.medidores", "med", "Medidores e leituras"},
	{"sol.solicitacoes", "sol", "Solicitações"},
	{"sol.solicitacao_detalhe", "sol", "Detalhe da solicitação"},
	{"os.ordens", "os", "Ordens de serviço"},
	{"os.ordem_detalhe", "os", "Detalhe da ordem de serviço"},
	{"pre.planos", "pre", "Planos de manutenção"},
	{"chk.modelos", "chk", "Modelos de checklist"},
	{"est.itens", "est", "Itens de estoque"},
	{"est.movimentacoes", "est", "Movimentações de estoque"},
	{"ter.fornecedores", "ter", "Fornecedores e prestadores"},
	{"srv.conexoes", "srv", "Conexões entre contas"},
	{"prj.projetos", "prj", "Projetos"},
	{"ctr.contratos", "ctr", "Contratos"},
	{"loc.locatarios", "loc", "Locatários"},
	{"loc.contratos", "loc", "Contratos de locação"},
	{"fro.veiculos", "fro", "Veículos da frota"},
	{"con.conformidade", "con", "Conformidade e documentos"},
	{"int.integracoes", "int", "Integrações"},
}

// Telas devolve o catálogo fechado, na ordem da tela de configuração.
func Telas() []Tela { return slices.Clone(telas) }

// BuscarTela devolve a tela do catálogo pelo código.
func BuscarTela(codigo string) (Tela, bool) {
	for _, t := range telas {
		if t.Codigo == codigo {
			return t, true
		}
	}
	return Tela{}, false
}

// TelaValida informa se o código está no catálogo.
func TelaValida(codigo string) bool {
	_, ok := BuscarTela(codigo)
	return ok
}

// ModuloDaTela devolve o código do módulo de um código de tela.
func ModuloDaTela(codigo string) string {
	m, _, _ := strings.Cut(codigo, ".")
	return m
}
