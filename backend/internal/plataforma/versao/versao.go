// Package versao guarda a identificação da versão em execução.
//
// Os valores são injetados pelo build a partir da tag git `vX.Y.Z`, conforme
// RNF-039 e a seção 9.5 de docs/5_ARQUITETURA_SOFTWARE.md. Nunca são digitados
// à mão: o valor padrão existe só para desenvolvimento local.
package versao

import (
	"runtime/debug"
	"strings"
	"sync"
)

// Injetados com -ldflags -X. Ver ferramentas/versao.sh.
var (
	Versao    = "0.0.0-dev"
	Commit    = "desconhecido"
	DataBuild = ""
	Ambiente  = "desenvolvimento"
)

// Dados é o que o endpoint público expõe. Nada além destes três campos:
// sem dependência, sistema operacional, região ou infraestrutura.
type Dados struct {
	Versao    string `json:"versao" doc:"Versão do release em execução (SemVer)" example:"1.4.2"`
	Commit    string `json:"commit" doc:"Commit curto que gerou o build" example:"a1b2c3d"`
	DataBuild string `json:"data_build" doc:"Instante do build em UTC (RFC 3339)" example:"2026-11-12T14:05:00Z"`
}

// Detalhe é o que a tela "Sobre" mostra, só para quem está autenticado.
type Detalhe struct {
	Dados
	Ambiente string `json:"ambiente" example:"producao"`
}

var (
	umaVez        sync.Once
	commitDoBuild string
	sujo          bool
)

func lerBuildInfo() {
	umaVez.Do(func() {
		info, ok := debug.ReadBuildInfo()
		if !ok {
			return
		}
		for _, s := range info.Settings {
			switch s.Key {
			case "vcs.revision":
				commitDoBuild = s.Value
			case "vcs.modified":
				sujo = s.Value == "true"
			}
		}
	})
}

// CommitCurto devolve o commit em sete caracteres.
func CommitCurto() string {
	c := Commit
	if c == "" || c == "desconhecido" {
		lerBuildInfo()
		if commitDoBuild != "" {
			c = commitDoBuild
		}
	}
	if len(c) > 7 {
		return c[:7]
	}
	return c
}

// Coerente informa se o commit injetado pelo build confere com o que o próprio
// Go gravou no binário. Divergência significa build montado fora da esteira.
// Um binário com a árvore de trabalho suja nunca é coerente.
func Coerente() (bool, string) {
	lerBuildInfo()
	if commitDoBuild == "" {
		return true, "sem informação de controle de versão no binário"
	}
	if sujo {
		return false, "árvore de trabalho suja no momento do build"
	}
	if Commit == "" || Commit == "desconhecido" {
		return false, "commit não injetado pelo build"
	}
	if !strings.HasPrefix(commitDoBuild, Commit) && !strings.HasPrefix(Commit, commitDoBuild) {
		return false, "commit injetado (" + Commit + ") difere do gravado pelo Go (" + commitDoBuild + ")"
	}
	return true, ""
}

// Publico devolve os dados do endpoint GET /v1/versao.
func Publico() Dados {
	return Dados{Versao: Versao, Commit: CommitCurto(), DataBuild: DataBuild}
}

// Interno devolve os dados da tela "Sobre".
func Interno() Detalhe {
	return Detalhe{Dados: Publico(), Ambiente: Ambiente}
}

// PreLancamento informa se a versão tem sufixo de pré-lançamento (-rc.N, -dev).
// Produção recusa versão de pré-lançamento.
func PreLancamento() bool {
	return strings.Contains(Versao, "-")
}
