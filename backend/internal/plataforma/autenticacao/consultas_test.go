package autenticacao

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// As sete tabelas de credencial (ADR-011; seção 6.9 do documento 5).
var tabelasCredencial = []string{
	"plataforma.identidade", "plataforma.acesso_empresa", "plataforma.sessao", "plataforma.token_renovacao",
	"plataforma.codigo_autorizacao", "plataforma.desafio_mfa", "plataforma.token_email",
}

const dirRepo = "internal/repo"

// Nenhum .sql do backend fora do repositório da autenticação cita as tabelas
// de credencial: o pool da aplicação não as alcança, e uma consulta que as
// cite em outro módulo falharia em produção por falta de privilégio.
func TestCredenciaisSoNoRepositorioDaAutenticacao(t *testing.T) {
	raiz := filepath.Join("..", "..") // backend/internal
	err := filepath.WalkDir(raiz, func(caminho string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() || !strings.HasSuffix(caminho, ".sql") {
			return err
		}
		if strings.Contains(filepath.ToSlash(caminho), "plataforma/autenticacao/"+dirRepo+"/") {
			return nil
		}
		b, err := os.ReadFile(caminho)
		if err != nil {
			return err
		}
		texto := strings.ToLower(string(b))
		for _, tab := range tabelasCredencial {
			if regexp.MustCompile(regexp.QuoteMeta(tab) + `\b`).MatchString(texto) {
				t.Errorf("%s cita %s, que só o repositório da autenticação lê", caminho, tab)
			}
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// Toda consulta do repositório que lê ou altera uma tabela de credencial
// filtra por chave exata: hash, id, identidade_id, (empresa_id, usuario_id),
// família (ids de sessão) ou e-mail. Nem o pacote de autenticação enumera
// credenciais (seção 6.9 do documento 5).
func TestConsultasDoRepositorioFiltramPorChaveExata(t *testing.T) {
	b, err := os.ReadFile(filepath.Join(dirRepo, "autenticacao.sql"))
	if err != nil {
		t.Fatal(err)
	}
	permitidos := []string{
		"token_hash = sqlc.arg(token_hash)", "codigo_hash = sqlc.arg(codigo_hash)", "id = sqlc.arg(id)",
		"identidade_id = sqlc.arg(identidade_id)", "email = sqlc.arg(email)",
		"empresa_id = sqlc.arg(empresa_id) and usuario_id = sqlc.arg(usuario_id)",
		"familia_id = any (sqlc.arg(familias)::uuid[])",
	}
	consultas := strings.Split(strings.ToLower(string(b)), "-- name: ")[1:]
	if len(consultas) < 20 {
		t.Fatalf("esperadas as consultas do repositório; achadas %d", len(consultas))
	}
	for _, c := range consultas {
		nome := strings.Fields(c)[0]
		corpo := c[strings.Index(c, "\n"):]
		if strings.HasPrefix(strings.TrimSpace(corpo), "insert") {
			continue
		}
		toca := false
		for _, tab := range tabelasCredencial {
			toca = toca || strings.Contains(corpo, tab)
		}
		if !toca {
			continue
		}
		i := strings.Index(corpo, "where")
		if i < 0 {
			t.Errorf("%s não tem WHERE", nome)
			continue
		}
		onde := corpo[i:]
		ok := false
		for _, p := range permitidos {
			ok = ok || strings.Contains(strings.Join(strings.Fields(onde), " "), p)
		}
		if !ok {
			t.Errorf("%s não filtra por chave exata", nome)
		}
	}
}
