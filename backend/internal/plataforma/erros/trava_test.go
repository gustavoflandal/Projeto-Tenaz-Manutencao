package erros_test

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"testing/fstest"

	"tenaz/internal/plataforma/erros"
	"tenaz/internal/plataforma/textos"
)

// A trava do texto livre (ADR-013, Parte B, item 5): texto para pessoa não se
// escreve no código Go. Erro com texto sai por ComChave e ComCampoChave, que
// buscam no catálogo do servidor; texto só para o log sai por Tecnico. Este
// teste é parte da tarefa `textos` do CI e recusa, no código de produção:
//
//   - qualquer chamada a .Com( e .ComCampo(, que aceitam texto livre;
//   - erros.Novo( com um código sem título nos dois catálogos, ou com título
//     diferente do catálogo pt-BR;
//   - ComChave( e ComCampoChave( com chave literal que falta em um dos dois
//     catálogos (chave ausente mostraria a própria chave à pessoa), e
//     textos.Chave( idem; chave montada em tempo de execução não é literal e
//     fica a cargo dos testes de catálogo;
//   - import com apelido dos pacotes erros e textos é seguido pelo apelido; o
//     import com ponto desses pacotes é proibido.
//
// O mesmo olhar é aplicado a fontes falsas, para provar que a trava acusa.

type achado struct {
	posicao string
	motivo  string
}

const (
	caminhoErros  = "tenaz/internal/plataforma/erros"
	caminhoTextos = "tenaz/internal/plataforma/textos"
)

// nomeLocal devolve o nome pelo qual o arquivo enxerga o pacote: o apelido do
// import, ou o nome padrão quando o import não existe (fontes falsas sem
// imports) ou não tem apelido.
func nomeLocal(f *ast.File, caminho, padrao string) (nome string, ponto bool) {
	for _, im := range f.Imports {
		c, err := strconv.Unquote(im.Path.Value)
		if err != nil || c != caminho {
			continue
		}
		if im.Name == nil {
			return padrao, false
		}
		return im.Name.Name, im.Name.Name == "."
	}
	return padrao, false
}

// chaveFaltante confere uma chave literal nos dois catálogos.
func chaveFaltante(catalogo *textos.Catalogo, chave string) []string {
	var falta []string
	for _, id := range textos.Idiomas() {
		if !catalogo.Existe(id, chave) {
			falta = append(falta, string(id))
		}
	}
	return falta
}

// procurar devolve os problemas de um arquivo Go já interpretado.
func procurar(fset *token.FileSet, f *ast.File, catalogo *textos.Catalogo) []achado {
	var res []achado
	nomeErros, pontoErros := nomeLocal(f, caminhoErros, "erros")
	nomeTextos, pontoTextos := nomeLocal(f, caminhoTextos, "textos")
	if pontoErros || pontoTextos {
		res = append(res, achado{fset.Position(f.Package).String(), "import com ponto de erros ou textos esconde as chamadas da trava"})
	}
	chaveLiteral := func(pos string, arg ast.Expr, quem string) {
		chave, ok := literal(arg)
		if !ok {
			return
		}
		if falta := chaveFaltante(catalogo, chave); len(falta) > 0 {
			res = append(res, achado{pos, quem + ": chave " + chave + " sem texto em " + strings.Join(falta, ", ")})
		}
	}
	ast.Inspect(f, func(n ast.Node) bool {
		call, ok := n.(*ast.CallExpr)
		if !ok {
			return true
		}
		sel, ok := call.Fun.(*ast.SelectorExpr)
		if !ok {
			return true
		}
		pos := fset.Position(call.Pos()).String()
		switch sel.Sel.Name {
		case "Com", "ComCampo":
			res = append(res, achado{pos, "." + sel.Sel.Name + "( aceita texto livre; use ComChave, ComCampoChave ou Tecnico"})
		case "ComChave":
			if len(call.Args) >= 1 {
				chaveLiteral(pos, call.Args[0], "ComChave")
			}
		case "ComCampoChave":
			if len(call.Args) >= 3 {
				chaveLiteral(pos, call.Args[2], "ComCampoChave")
			}
		case "Chave":
			if pkg, ok := sel.X.(*ast.Ident); ok && pkg.Name == nomeTextos && len(call.Args) == 1 {
				chaveLiteral(pos, call.Args[0], "textos.Chave")
			}
		case "Novo":
			pkg, ok := sel.X.(*ast.Ident)
			if !ok || pkg.Name != nomeErros || len(call.Args) < 2 {
				return true
			}
			codigo, okC := literal(call.Args[0])
			titulo, okT := literal(call.Args[1])
			if !okC || !okT {
				res = append(res, achado{pos, "erros.Novo( com código e título que não são literais"})
				return true
			}
			chave := erros.ChaveTitulo(codigo)
			for _, id := range textos.Idiomas() {
				if !catalogo.Existe(id, chave) {
					res = append(res, achado{pos, "código " + codigo + " sem título em " + string(id) + " (chave " + chave + ")"})
				}
			}
			if pt, _ := catalogo.Bruto(textos.PtBR, chave); pt != titulo {
				res = append(res, achado{pos, "título de " + codigo + " difere do catálogo pt-BR"})
			}
		}
		return true
	})
	return res
}

func literal(e ast.Expr) (string, bool) {
	l, ok := e.(*ast.BasicLit)
	if !ok || l.Kind != token.STRING {
		return "", false
	}
	s, err := strconv.Unquote(l.Value)
	return s, err == nil
}

func achados(t *testing.T, fonte string, catalogo *textos.Catalogo) []achado {
	t.Helper()
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "falso.go", fonte, 0)
	if err != nil {
		t.Fatal(err)
	}
	return procurar(fset, f, catalogo)
}

// A prova de que a trava acusa: fontes falsas com cada forma proibida.
func TestTravaAcusaTextoLivre(t *testing.T) {
	c := textos.Padrao()
	proibidos := map[string]string{
		"Com com literal":                 `package x; func f() { _ = erros.NaoEncontrado.Com("recurso sumiu") }`,
		"Com com formato":                 `package x; func f() { _ = erros.NaoEncontrado.Com("falta %s", p) }`,
		"Com com variável":                `package x; func f() { _ = erros.NaoEncontrado.Com(msg) }`,
		"ComCampo com literal":            `package x; func f() { _ = erros.EntradaInvalida.ComCampo("a", "obrigatorio", "Informe a.") }`,
		"ComCampo encadeado":              "package x; func f() { _ = erros.X.\n\tComChave(\"detalhe.corpo_json\", nil).\n\tComCampo(c, k, \"texto\") }",
		"Novo com código sem texto":       `package x; var e = erros.Novo("codigo_inexistente", "Qualquer", 400)`,
		"Novo com título errado":          `package x; var e = erros.Novo("ultimo_administrador", "Outro título", 409)`,
		"Novo sem literal":                `package x; var e = erros.Novo(cod, tit, 400)`,
		"ComChave com chave ausente":      `package x; func f() { _ = erros.NaoEncontrado.ComChave("detalhe.nao_existe_em_lugar_nenhum", nil) }`,
		"ComCampoChave com chave ausente": `package x; func f() { _ = erros.EntradaInvalida.ComCampoChave("a", "obrigatorio", "campo.nao_existe", nil) }`,
		"textos.Chave com chave ausente":  `package x; func f() { _ = textos.Parametros{"s": textos.Chave("situacao.nao_existe")} }`,
		"Novo com apelido do pacote": `package x; import e "tenaz/internal/plataforma/erros"
var v = e.Novo("codigo_inexistente", "Qualquer", 400)`,
		"Com através de apelido": `package x; import e "tenaz/internal/plataforma/erros"
func f() { _ = e.NaoEncontrado.Com("texto livre") }`,
		"ComChave com chave ausente e apelido": `package x; import e "tenaz/internal/plataforma/erros"
func f() { _ = e.NaoEncontrado.ComChave("detalhe.nao_existe", nil) }`,
		"textos.Chave com apelido": `package x; import tx "tenaz/internal/plataforma/textos"
func f() { _ = tx.Chave("situacao.nao_existe") }`,
		"import com ponto": `package x; import . "tenaz/internal/plataforma/erros"
var v = Novo("codigo_inexistente", "Qualquer", 400)`,
	}
	for nome, fonte := range proibidos {
		if len(achados(t, fonte, c)) == 0 {
			t.Errorf("a trava não acusou: %s", nome)
		}
	}
	permitidos := map[string]string{
		"ComChave":                  `package x; func f() { _ = erros.NaoEncontrado.ComChave("detalhe.rota_inexistente", nil) }`,
		"ComCampoChave":             `package x; func f() { _ = erros.EntradaInvalida.ComCampoChave("a", "obrigatorio", "campo.informe", nil) }`,
		"Tecnico":                   `package x; func f() { _ = erros.Interno.Tecnico("falha de programação %d", 1) }`,
		"Novo com texto":            `package x; var e = erros.Novo("ultimo_administrador", "A empresa ficaria sem administrador", 409)`,
		"chamada qualquer":          `package x; func f() { _ = strings.Join(a, ",") }`,
		"chave montada em execução": `package x; func f() { _ = erros.NaoEncontrado.ComChave("detalhe."+x, nil); _ = textos.Chave("situacao." + y) }`,
		"Novo com apelido": `package x; import e "tenaz/internal/plataforma/erros"
var v = e.Novo("ultimo_administrador", "A empresa ficaria sem administrador", 409)`,
		"ComChave com chave presente e apelido": `package x; import e "tenaz/internal/plataforma/erros"
func f() { _ = e.NaoEncontrado.ComChave("detalhe.rota_inexistente", nil) }`,
	}
	for nome, fonte := range permitidos {
		if a := achados(t, fonte, c); len(a) != 0 {
			t.Errorf("a trava acusou o que é permitido (%s): %v", nome, a)
		}
	}
}

// Chave que existe em um idioma e falta no outro também é acusada.
func TestTravaAcusaChaveFaltandoEmUmIdioma(t *testing.T) {
	cat, err := textos.Novo(fstest.MapFS{
		"pt-BR.json":  {Data: []byte(`{"detalhe":{"so_pt":"Só em português","nos_dois":"Nos dois"}}`)},
		"es-419.json": {Data: []byte(`{"detalhe":{"nos_dois":"En los dos"}}`)},
	})
	if err != nil {
		t.Fatal(err)
	}
	a := achados(t, `package x; func f() { _ = erros.NaoEncontrado.ComChave("detalhe.so_pt", nil) }`, cat)
	if len(a) != 1 || !strings.Contains(a[0].motivo, "es-419") {
		t.Errorf("deveria acusar a falta em es-419: %v", a)
	}
	if a := achados(t, `package x; func f() { _ = erros.NaoEncontrado.ComChave("detalhe.nos_dois", nil) }`, cat); len(a) != 0 {
		t.Errorf("chave presente nos dois idiomas acusada: %v", a)
	}
}

// A trava no código de produção: nenhum .Com( ou .ComCampo( fora dos testes.
func TestSemTextoLivreNoCodigoDeProducao(t *testing.T) {
	c := textos.Padrao()
	var problemas []string
	for _, raiz := range []string{"../..", "../../../cmd"} {
		err := filepath.WalkDir(raiz, func(caminho string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() || !strings.HasSuffix(caminho, ".go") || strings.HasSuffix(caminho, "_test.go") {
				return nil
			}
			fset := token.NewFileSet()
			f, err := parser.ParseFile(fset, caminho, nil, 0)
			if err != nil {
				return err
			}
			for _, a := range procurar(fset, f, c) {
				problemas = append(problemas, a.posicao+": "+a.motivo)
			}
			return nil
		})
		if err != nil {
			t.Fatal(err)
		}
	}
	if len(problemas) > 0 {
		t.Fatalf("texto livre no código de produção (ADR-013, a trava do documento 5):\n  %s", strings.Join(problemas, "\n  "))
	}
}
