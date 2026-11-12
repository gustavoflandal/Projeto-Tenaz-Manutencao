package vazamento

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"strings"

	esquema "tenaz/db"
)

// arquivoListaGlobais guarda os catálogos globais permitidos sem empresa_id.
// A lista fica no repositório de propósito: acrescentar uma linha a ela é uma
// decisão visível em revisão de código, não um efeito colateral de migração.
const arquivoListaGlobais = "tabelas_globais.txt"

// lerListaDeGlobais devolve o conjunto de tabelas permitidas sem empresa_id.
func lerListaDeGlobais() (map[string]bool, error) {
	// O teste roda com o diretório do pacote como diretório de trabalho.
	dados, err := os.ReadFile(arquivoListaGlobais)
	if err != nil {
		return nil, err
	}

	permitidas := map[string]bool{}
	leitor := bufio.NewScanner(strings.NewReader(string(dados)))
	for leitor.Scan() {
		linha := strings.TrimSpace(leitor.Text())
		if linha == "" || strings.HasPrefix(linha, "#") {
			continue
		}
		// Formato: esquema.tabela  # justificativa
		if i := strings.Index(linha, "#"); i >= 0 {
			linha = strings.TrimSpace(linha[:i])
		}
		if linha == "" {
			continue
		}
		if !strings.Contains(linha, ".") {
			return nil, fmt.Errorf("linha %q sem esquema; use esquema.tabela", linha)
		}
		permitidas[linha] = true
	}
	return permitidas, leitor.Err()
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
