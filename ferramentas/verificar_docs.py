"""Verifica a integridade e as referências cruzadas da documentação da Tenaz.

Uso:
    python ferramentas/verificar_docs.py

Confere:
  1. Links relativos em Markdown apontam para arquivos que existem.
  2. Todo ID citado (RF-, RNF-, RN-, ADR-) está definido no documento dono:
       RF e RNF  -> docs/1_ESPECIFICACAO_REQUISITOS.md
       RN        -> docs/4_FLUXOS_PROCESSO.md
       ADR       -> docs/5_ARQUITETURA_SOFTWARE.md
     e nenhum ID é definido duas vezes.
  3. Os documentos 0 a 9 têm a linha de versão no cabeçalho.
  4. O registro de etapas (docs/etapas) segue o padrão de nome, a numeração é
     contínua e o índice lista todas as etapas.
  5. Aviso: PDF em docs/pdf mais antigo que o Markdown de origem.

Termina com código 1 se houver erro (avisos não falham a verificação).
Regra do projeto: CLAUDE.md, seção 3.3.
"""
import pathlib
import re
import sys
import urllib.parse

RAIZ = pathlib.Path(__file__).resolve().parent.parent
DOCS = RAIZ / "docs"
DONO = {
    "RF": DOCS / "1_ESPECIFICACAO_REQUISITOS.md",
    "RNF": DOCS / "1_ESPECIFICACAO_REQUISITOS.md",
    "RN": DOCS / "4_FLUXOS_PROCESSO.md",
    "ADR": DOCS / "5_ARQUITETURA_SOFTWARE.md",
}
# Como cada documento dono define um ID: título de seção ou primeira célula de tabela.
DEFINICAO = {
    "RF": re.compile(r"^#{2,6}\s+(RF-[A-Z]{2,3}-\d{3})\b", re.M),
    "RNF": re.compile(r"^\|\s*(RNF-\d{3})\s*\|", re.M),
    "RN": re.compile(r"^\|\s*(RN-[A-Z]{2,3}-\d{3})\s*\|", re.M),
    "ADR": re.compile(r"^#{2,6}\s+(ADR-\d{3})\b", re.M),
}
CITACAO = re.compile(r"\b(RNF-\d{3}|RF-[A-Z]{2,3}-\d{3}|RN-[A-Z]{2,3}-\d{3}|ADR-\d{3})\b")
INTERVALO = re.compile(r"\b((?:RF|RN)-[A-Z]{2,3})-(\d{3})\s+a\s+(?:\1-)?(\d{3})\b|\b(RNF|ADR)-(\d{3})\s+a\s+(?:\4-)?(\d{3})\b")
LINK = re.compile(r"(?<!!)\[[^\]]*\]\(([^)\s]+)\)")
IGNORAR_PASTAS = {".git", "node_modules"}

erros: list[str] = []
avisos: list[str] = []


def rel(p: pathlib.Path) -> str:
    return p.relative_to(RAIZ).as_posix()


def arquivos_md():
    for p in RAIZ.rglob("*.md"):
        if not any(parte in IGNORAR_PASTAS for parte in p.parts):
            yield p


def sem_codigo(texto: str) -> str:
    """Remove blocos de código cercados, para não validar exemplos."""
    return re.sub(r"```.*?```", "", texto, flags=re.S)


def verificar_links():
    for arq in arquivos_md():
        texto = sem_codigo(arq.read_text(encoding="utf-8"))
        for n, linha in enumerate(texto.splitlines(), 1):
            for alvo in LINK.findall(linha):
                if re.match(r"^(https?:|mailto:|#)", alvo):
                    continue
                caminho = urllib.parse.unquote(alvo.split("#", 1)[0])
                if not caminho:
                    continue
                destino = (arq.parent / caminho).resolve()
                if not destino.exists():
                    erros.append(f"{rel(arq)}:{n}: link quebrado -> {alvo}")


def ids_definidos():
    definidos: dict[str, set[str]] = {}
    for tipo, dono in DONO.items():
        texto = dono.read_text(encoding="utf-8")
        encontrados = DEFINICAO[tipo].findall(texto)
        vistos = set()
        for i in encontrados:
            if i in vistos:
                erros.append(f"{rel(dono)}: {i} definido mais de uma vez")
            vistos.add(i)
        definidos[tipo] = vistos
    return definidos


def tipo_de(id_: str) -> str:
    if id_.startswith("RNF-"):
        return "RNF"
    if id_.startswith("RF-"):
        return "RF"
    if id_.startswith("RN-"):
        return "RN"
    return "ADR"


def verificar_ids():
    definidos = ids_definidos()
    todos = set().union(*definidos.values())
    for arq in arquivos_md():
        if "research_notes" in arq.parts:
            continue
        texto = arq.read_text(encoding="utf-8")
        citados = set(CITACAO.findall(texto))
        # Intervalos "RF-TER-001 a RF-TER-009": confere as pontas.
        for m in INTERVALO.finditer(texto):
            if m.group(1):
                citados.update({f"{m.group(1)}-{m.group(2)}", f"{m.group(1)}-{m.group(3)}"})
            else:
                citados.update({f"{m.group(4)}-{m.group(5)}", f"{m.group(4)}-{m.group(6)}"})
        for i in sorted(citados):
            if i not in todos:
                erros.append(f"{rel(arq)}: {i} citado, mas não definido em {rel(DONO[tipo_de(i)])}")


def verificar_versoes():
    for arq in sorted(DOCS.glob("[0-9]_*.md")):
        cabecalho = "\n".join(arq.read_text(encoding="utf-8").splitlines()[:15])
        if not re.search(r"Vers[aã]o\**:?\**\s*\d+\.\d+", cabecalho):
            erros.append(f"{rel(arq)}: sem linha de versão no cabeçalho")


def verificar_etapas():
    pasta = DOCS / "etapas"
    if not pasta.exists():
        erros.append("docs/etapas não existe")
        return
    padrao = re.compile(r"^(\d{3})_(\d{4}-\d{2}-\d{2})_[a-z0-9-]+\.md$")
    numeros = []
    for p in sorted(pasta.glob("*.md")):
        if p.name == "README.md":
            continue
        m = padrao.match(p.name)
        if not m:
            erros.append(f"{rel(p)}: nome fora do padrão NNN_AAAA-MM-DD_titulo.md")
            continue
        numeros.append(int(m.group(1)))
        primeira = p.read_text(encoding="utf-8").splitlines()[0]
        if not primeira.startswith(f"# {m.group(1)} "):
            erros.append(f"{rel(p)}: o título deve começar com '# {m.group(1)} —'")
    esperados = list(range(1, len(numeros) + 1))
    if sorted(numeros) != esperados:
        erros.append(f"docs/etapas: numeração não é contínua: {sorted(numeros)}")
    indice = (pasta / "README.md").read_text(encoding="utf-8")
    for p in pasta.glob("[0-9]*.md"):
        if p.name not in indice:
            erros.append(f"docs/etapas/README.md: índice não lista {p.name}")


def verificar_pdfs():
    for pdf in (DOCS / "pdf").glob("*.pdf"):
        origem = DOCS / f"{pdf.stem}.md"
        if not origem.exists():
            origem = RAIZ / f"{pdf.stem}.md"
        if origem.exists() and origem.stat().st_mtime > pdf.stat().st_mtime + 1:
            avisos.append(f"{rel(pdf)}: mais antigo que {rel(origem)}; regere com ferramentas/md2pdf.py")


def main():
    verificar_links()
    verificar_ids()
    verificar_versoes()
    verificar_etapas()
    verificar_pdfs()
    for a in avisos:
        print(f"AVISO  {a}")
    for e in erros:
        print(f"ERRO   {e}")
    print(f"\n{len(erros)} erro(s), {len(avisos)} aviso(s).")
    sys.exit(1 if erros else 0)


if __name__ == "__main__":
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8")
    main()
