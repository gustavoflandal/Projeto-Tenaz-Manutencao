"""Converte um documento Markdown da Tenaz em PDF com a identidade visual da marca.

Uso:
    python ferramentas/md2pdf.py docs/arquivo.md [saida.pdf]

Requer o pacote Python `markdown` e o Google Chrome (ou Edge) instalado.
O PDF é gerado pelo modo headless do navegador (--print-to-pdf).
"""
import base64
import html
import pathlib
import re
import subprocess
import sys
import tempfile

import markdown

RAIZ = pathlib.Path(__file__).resolve().parent.parent
LOGO = RAIZ / "marca" / "logo" / "svg" / "tenaz-horizontal-positivo.svg"
NAVEGADORES = [
    r"C:\Program Files\Google\Chrome\Application\chrome.exe",
    r"C:\Program Files (x86)\Google\Chrome\Application\chrome.exe",
    r"C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe",
    "/usr/bin/google-chrome",
    "/usr/bin/chromium",
]

# Mermaid carregado do CDN (versão fixa) para desenhar diagramas; exige internet na geração.
MERMAID_JS = "https://cdn.jsdelivr.net/npm/mermaid@10.9.1/dist/mermaid.min.js"

CSS = """
@page { size: A4; margin: 22mm 18mm 20mm 18mm;
  @bottom-left { content: "Tenaz · Gestão de Manutenção e Frota"; font: 8pt 'Segoe UI', Arial, sans-serif; color: #6B7887; }
  @bottom-right { content: counter(page) " / " counter(pages); font: 8pt 'Segoe UI', Arial, sans-serif; color: #6B7887; }
}
@page :first { margin: 0; @bottom-left { content: none; } @bottom-right { content: none; } }
* { box-sizing: border-box; }
html { -webkit-print-color-adjust: exact; print-color-adjust: exact; }
body { margin: 0; font-family: 'Segoe UI', 'Inter', Arial, sans-serif; font-size: 10pt; line-height: 1.55; color: #14202E; }
.capa { height: 297mm; width: 210mm; padding: 28mm 22mm; display: flex; flex-direction: column; justify-content: space-between; background: #F2F4F6; page-break-after: always; position: relative; }
.capa::after { content: ""; position: absolute; left: 0; right: 0; bottom: 0; height: 8mm; background: #FF6A13; }
.capa img { width: 90mm; }
.capa .tipo { font: 600 9pt Consolas, 'Cascadia Mono', monospace; letter-spacing: .12em; text-transform: uppercase; color: #B54500; }
.capa h1 { font-size: 30pt; line-height: 1.1; margin: 4mm 0 6mm; border: 0; padding: 0; font-weight: 800; letter-spacing: -.01em; }
.capa .sub { font-size: 12pt; color: #4A5868; max-width: 150mm; }
.capa .meta { font: 9pt Consolas, 'Cascadia Mono', monospace; color: #4A5868; line-height: 1.8; margin-bottom: 10mm; }
main { padding: 0; }
h1, h2, h3, h4 { line-height: 1.2; break-after: avoid; color: #14202E; }
h1 { font-size: 20pt; font-weight: 800; margin: 0 0 6mm; padding-bottom: 3mm; border-bottom: 2px solid #14202E; }
h2 { font-size: 14pt; font-weight: 800; margin: 9mm 0 3mm; padding-top: 2mm; border-top: 1px solid #D9DFE6; }
h3 { font-size: 11.5pt; font-weight: 700; margin: 6mm 0 2mm; }
h4 { font-size: 10pt; font-weight: 700; margin: 4mm 0 1.5mm; color: #4A5868; }
p { margin: 0 0 2.5mm; orphans: 3; widows: 3; }
ul, ol { margin: 0 0 3mm; padding-left: 6mm; }
li { margin-bottom: 1mm; }
a { color: #B54500; text-decoration: none; word-break: break-word; }
strong { font-weight: 700; }
code { font-family: Consolas, 'Cascadia Mono', monospace; font-size: 8.8pt; background: #EEF1F4; padding: .2mm 1mm; border-radius: 1mm; }
pre { background: #14202E; color: #EEF2F6; padding: 4mm; border-radius: 2mm; overflow: hidden; white-space: pre-wrap; font-size: 8.5pt; break-inside: avoid; }
pre code { background: none; padding: 0; color: inherit; }
pre.mermaid { background: #FFFFFF; color: #14202E; border: 1px solid #D9DFE6; text-align: center; white-space: pre; }
pre.mermaid svg { max-width: 100%; height: auto; }
blockquote { margin: 3mm 0; padding: 2mm 4mm; border-left: 3px solid #FF6A13; background: #FFF4EC; color: #14202E; break-inside: avoid; }
blockquote p:last-child { margin: 0; }
table { width: 100%; border-collapse: collapse; margin: 3mm 0 5mm; font-size: 8.6pt; break-inside: auto; }
thead { display: table-header-group; }
tr { break-inside: avoid; }
th { background: #14202E; color: #FFFFFF; text-align: left; font-weight: 600; padding: 1.8mm 2mm; }
td { padding: 1.6mm 2mm; border-bottom: 1px solid #D9DFE6; vertical-align: top; }
tbody tr:nth-child(even) td { background: #F7F8FA; }
hr { border: 0; border-top: 1px solid #D9DFE6; margin: 6mm 0; }
.toc { background: #F2F4F6; padding: 4mm 6mm; border-radius: 2mm; margin-bottom: 6mm; break-inside: avoid; }
.toc ul { list-style: none; padding-left: 4mm; margin: 0; }
.toc > ul { padding-left: 0; }
"""


def navegador():
    for caminho in NAVEGADORES:
        if pathlib.Path(caminho).exists():
            return caminho
    sys.exit("Chrome ou Edge não encontrado.")


def converter(origem: pathlib.Path, destino: pathlib.Path):
    texto = origem.read_text(encoding="utf-8")
    # O primeiro título vira a capa
    titulo = origem.stem
    m = re.match(r"\s*#\s+(.+)\n", texto)
    if m:
        titulo = m.group(1).strip()
        texto = texto[m.end():]
    corpo = markdown.markdown(
        texto,
        extensions=["tables", "fenced_code", "sane_lists", "toc", "attr_list", "def_list"],
        extension_configs={"toc": {"permalink": False}},
        output_format="html5",
    )
    # Blocos ```mermaid viram <pre class="mermaid"> para o Mermaid desenhar
    corpo, n_mermaid = re.subn(
        r'<pre><code class="language-mermaid">(.*?)</code></pre>',
        r'<pre class="mermaid">\1</pre>',
        corpo,
        flags=re.S,
    )
    script_mermaid = ""
    if n_mermaid:
        script_mermaid = (
            f'<script src="{MERMAID_JS}"></script>'
            "<script>mermaid.initialize({startOnLoad:true,theme:'base',securityLevel:'strict',"
            "flowchart:{useMaxWidth:true},er:{useMaxWidth:true},state:{useMaxWidth:true},"
            "themeVariables:{fontFamily:'Segoe UI, Arial, sans-serif',fontSize:'13px',"
            "primaryColor:'#F2F4F6',primaryBorderColor:'#14202E',primaryTextColor:'#14202E',"
            "lineColor:'#4A5868',secondaryColor:'#FFF4EC',tertiaryColor:'#FFFFFF'}});</script>"
        )
    logo = ""
    if LOGO.exists():
        logo = "data:image/svg+xml;base64," + base64.b64encode(LOGO.read_bytes()).decode()
    doc = f"""<!doctype html><html lang="pt-BR"><head><meta charset="utf-8">
<title>{html.escape(titulo)}</title><style>{CSS}</style></head><body>
<section class="capa">
  <div>{f'<img src="{logo}" alt="Tenaz">' if logo else ''}</div>
  <div>
    <div class="tipo">Documentação do projeto</div>
    <h1>{html.escape(titulo)}</h1>
    <div class="sub">Tenaz · Gestão de Manutenção e Frota</div>
  </div>
  <div class="meta">Arquivo: {html.escape(origem.name)}</div>
</section>
<main>{corpo}</main>{script_mermaid}</body></html>"""
    with tempfile.TemporaryDirectory() as tmp:
        arq = pathlib.Path(tmp) / "doc.html"
        arq.write_text(doc, encoding="utf-8")
        subprocess.run(
            [navegador(), "--headless=new", "--disable-gpu", "--no-pdf-header-footer",
             "--virtual-time-budget=30000", "--run-all-compositor-stages-before-draw",
             f"--print-to-pdf={destino}", arq.as_uri()],
            check=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL, timeout=180,
        )
    print(f"PDF gerado: {destino}")


if __name__ == "__main__":
    if len(sys.argv) < 2:
        sys.exit(__doc__)
    origem = pathlib.Path(sys.argv[1]).resolve()
    destino = pathlib.Path(sys.argv[2]).resolve() if len(sys.argv) > 2 else origem.with_suffix(".pdf")
    converter(origem, destino)
