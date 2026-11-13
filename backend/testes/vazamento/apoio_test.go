package vazamento

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"strings"

	esquema "tenaz/db"
)

// Listas versionadas da suíte. Ficam no repositório de propósito: acrescentar
// uma linha a qualquer uma delas é uma decisão visível em revisão de código,
// não um efeito colateral de migração, e o revisor é o tenaz-isolamento.
const (
	// Tabelas fora do RLS por empresa (seção 2.11 do documento 2).
	arquivoListaGlobais = "tabelas_globais.txt"
	// Tabelas de app com empresa_id cuja defesa tem outra forma (seção 2.8).
	arquivoExcecoesCobertura = "excecoes_cobertura.txt"
	// Tabelas e comandos alcançáveis por escopo externo (seção 2.9.4 do
	// documento 2; seção 6.6 do documento 5).
	arquivoAcessoExterno = "tabelas_acesso_externo.txt"
	// Funções SECURITY DEFINER aprovadas e quem as executa (consulta 6).
	arquivoSecurityDefiner = "funcoes_security_definer.txt"
	// Privilégios fora do RLS: fora de app, do sincronizador e da autenticação.
	arquivoPrivilegios = "privilegios.txt"
	// Relações de membro entre os papéis de execução.
	arquivoRelacoesDeMembro = "relacoes_de_membro.txt"
)

// entrada é uma linha de lista: nome, campos opcionais e justificativa.
//
// Formato: nome [campo ...]  # justificativa
type entrada struct {
	Nome          string
	Campos        []string
	Justificativa string
}

type lista []entrada

func (l lista) nomes() []string {
	n := make([]string, 0, len(l))
	for _, e := range l {
		n = append(n, e.Nome)
	}
	return n
}

// lerLista lê um arquivo de lista. Linha sem justificativa é erro: a lista é
// a memória da decisão, e uma decisão sem porquê não se revisa.
func lerLista(arquivo string) (lista, error) {
	return lerListaDe(arquivo, true)
}

// lerListaDe é lerLista com a regra do esquema opcional: a lista de relações
// de membro nomeia papéis, que não têm esquema.
func lerListaDe(arquivo string, exigirEsquema bool) (lista, error) {
	// O teste roda com o diretório do pacote como diretório de trabalho.
	dados, err := os.ReadFile(arquivo)
	if err != nil {
		return nil, err
	}

	var saida lista
	vistos := map[string]bool{}
	leitor := bufio.NewScanner(strings.NewReader(string(dados)))
	for n := 1; leitor.Scan(); n++ {
		linha := strings.TrimSpace(leitor.Text())
		if linha == "" || strings.HasPrefix(linha, "#") {
			continue
		}
		var justificativa string
		if i := strings.Index(linha, "#"); i >= 0 {
			justificativa = strings.TrimSpace(linha[i+1:])
			linha = strings.TrimSpace(linha[:i])
		}
		if justificativa == "" {
			return nil, fmt.Errorf("%s:%d: linha sem justificativa (use `# porquê`)", arquivo, n)
		}
		campos := strings.Fields(linha)
		if len(campos) == 0 {
			return nil, fmt.Errorf("%s:%d: linha sem nome", arquivo, n)
		}
		if exigirEsquema && !strings.Contains(campos[0], ".") && !strings.HasPrefix(campos[0], "esquema:") {
			return nil, fmt.Errorf("%s:%d: %q sem esquema; use esquema.nome (ou esquema:nome para o próprio esquema)",
				arquivo, n, campos[0])
		}
		chave := strings.Join(campos, " ")
		if vistos[chave] {
			return nil, fmt.Errorf("%s:%d: linha repetida %q", arquivo, n, chave)
		}
		vistos[chave] = true
		saida = append(saida, entrada{Nome: campos[0], Campos: campos[1:], Justificativa: justificativa})
	}
	return saida, leitor.Err()
}

// classesDeGlobais são as classes aceitas em tabelas_globais.txt.
var classesDeGlobais = map[string]bool{"tenant": true, "credencial": true, "catalogo": true, "controle": true}

// lerListaDeGlobais devolve o conjunto de tabelas permitidas sem empresa_id.
func lerListaDeGlobais() (map[string]bool, error) {
	classes, err := lerClassesDeGlobais()
	if err != nil {
		return nil, err
	}
	permitidas := map[string]bool{}
	for nome := range classes {
		permitidas[nome] = true
	}
	return permitidas, nil
}

// lerClassesDeGlobais devolve tabela -> classe.
func lerClassesDeGlobais() (map[string]string, error) {
	l, err := lerLista(arquivoListaGlobais)
	if err != nil {
		return nil, err
	}
	classes := map[string]string{}
	for _, e := range l {
		if len(e.Campos) != 1 || !classesDeGlobais[e.Campos[0]] {
			return nil, fmt.Errorf("%s: %s precisa de exatamente uma classe (tenant, credencial, catalogo, controle), veio %v",
				arquivoListaGlobais, e.Nome, e.Campos)
		}
		classes[e.Nome] = e.Campos[0]
	}
	return classes, nil
}

// Migracao é um arquivo de migração lido do binário.
type Migracao struct {
	Nome     string
	Conteudo string
}

// listarMigracoes devolve as migrações embutidas, em ordem de nome.
func listarMigracoes() ([]Migracao, error) {
	entradas, err := fs.ReadDir(esquema.Migracoes, esquema.DiretorioMigracoes)
	if err != nil {
		return nil, err
	}

	var saida []Migracao
	for _, e := range entradas {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".sql") {
			continue
		}
		b, err := fs.ReadFile(esquema.Migracoes, esquema.DiretorioMigracoes+"/"+e.Name())
		if err != nil {
			return nil, err
		}
		saida = append(saida, Migracao{Nome: e.Name(), Conteudo: string(b)})
	}
	return saida, nil
}
