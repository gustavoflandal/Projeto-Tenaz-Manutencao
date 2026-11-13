"""Gera e confere backend/testes/casos/formatos.json (ADR-013, Parte G; documento 7, seção 3.11).

O oráculo abaixo monta o texto para pessoas a partir dos dígitos do texto
canônico, sem float, como formatos.ts e o formatador do Go devem fazer. Todo
caso escrito à mão é conferido contra o oráculo; divergência aborta sem gravar.

Fusos: o deslocamento de cada fuso fica numa tabela fechada abaixo, válida só
para os instantes escolhidos (longe de mudança de horário). Os instantes do
Brasil ficam entre maio e agosto (fora da janela histórica do horário de verão,
caso ele volte); os do Paraguai, entre dezembro e fevereiro, onde a regra antiga
(-04 no inverno, até a base tz 2024b) e a atual (-03 o ano todo) dão o mesmo
-03. Se o zoneinfo do Python tiver a base tz, ela também é conferida.

Uso:
    python -I backend/testes/casos/gerar_formatos.py            # regrava formatos.json
    python -I backend/testes/casos/gerar_formatos.py --conferir # falha se o JSON divergir
"""
import datetime as dt
import json
import re
import sys
from pathlib import Path

NBSP = " "
SIMBOLO = {"BR": "R$", "AR": "$", "PY": "Gs."}
MOEDA_DO_PAIS = {"BR": "BRL", "AR": "ARS", "PY": "PYG"}
CASAS = {"BRL": 2, "ARS": 2, "PYG": 0}
IDIOMAS = ("pt-BR", "es-419")

# Deslocamento em minutos, válido para os instantes deste arquivo.
DESLOCAMENTO = {
    "America/Sao_Paulo": -180,
    "America/Manaus": -240,
    "America/Rio_Branco": -300,
    "America/Noronha": -120,
    "America/Cuiaba": -240,
    "America/Argentina/Buenos_Aires": -180,
    "America/Argentina/Cordoba": -180,
    "America/Argentina/Ushuaia": -180,
    "America/Asuncion": -180,
}
PAIS_DO_FUSO = {f: ("PY" if f == "America/Asuncion" else "AR" if "/Argentina/" in f else "BR")
                for f in DESLOCAMENTO}

RE_DECIMAL = re.compile(r"-?(0|[1-9][0-9]*)(\.[0-9]+)?")
RE_INTEIRO = re.compile(r"-?(0|[1-9][0-9]*)")
RE_VALOR = re.compile(r"-?(0|[1-9][0-9]{0,15})(\.[0-9]{1,2})?")
RE_PRECO = re.compile(r"-?(0|[1-9][0-9]{0,11})(\.[0-9]{1,6})?")
RE_INSTANTE = re.compile(r"(\d{4})-(\d{2})-(\d{2})T(\d{2}):(\d{2}):00Z")

falhas = []


def confere(nome, obtido, esperado):
    if obtido != esperado:
        falhas.append(f"{nome}: oráculo {obtido!r}, escrito {esperado!r}")


def agrupar(digitos):
    if len(digitos) < 4:
        return digitos
    partes = []
    while digitos:
        partes.append(digitos[-3:])
        digitos = digitos[:-3]
    return ".".join(reversed(partes))


def numero(texto):
    assert texto.isascii() and RE_DECIMAL.fullmatch(texto), texto
    assert texto.lstrip("-").strip("0.") or not texto.startswith("-"), f"-0 não é canônico: {texto}"
    sinal = "-" if texto.startswith("-") else ""
    inteiro, _, frac = texto.lstrip("-").partition(".")
    return sinal + agrupar(inteiro) + ("," + frac if frac else "")


def dinheiro(pais, moeda, texto, especie):
    assert MOEDA_DO_PAIS[pais] == moeda, (pais, moeda)
    casas = CASAS[moeda]
    frac = texto.partition(".")[2]
    if especie == "valor":
        assert RE_VALOR.fullmatch(texto) and len(frac) == casas, f"valor não canônico: {texto}"
    else:
        assert RE_PRECO.fullmatch(texto) and casas <= len(frac) <= 6, f"preço não canônico: {texto}"
        assert len(frac) == casas or not frac.endswith("0"), f"preço não canônico: {texto}"
    n = numero(texto)
    sinal = "-" if n.startswith("-") else ""
    return sinal + SIMBOLO[pais] + NBSP + n.lstrip("-")


def data(texto):
    assert re.fullmatch(r"\d{4}-\d{2}-\d{2}", texto), texto
    d = dt.date.fromisoformat(texto)
    return f"{d.day:02d}/{d.month:02d}/{d.year:04d}"


def data_hora(fuso, instante):
    m = RE_INSTANTE.fullmatch(instante)
    assert m, f"instante fora da forma da API com segundos zerados: {instante}"
    utc = dt.datetime(*map(int, m.groups()), tzinfo=dt.timezone.utc)
    pais = PAIS_DO_FUSO[fuso]
    if pais == "BR":
        assert 5 <= utc.month <= 8, f"instante do Brasil fora de maio a agosto: {instante}"
    if pais == "PY":
        assert utc.month in (12, 1, 2), f"instante do Paraguai fora de dezembro a fevereiro: {instante}"
    local = utc + dt.timedelta(minutes=DESLOCAMENTO[fuso])
    try:
        import zoneinfo
        z = zoneinfo.ZoneInfo(fuso)
    except Exception:  # sem base tz no Python (Windows sem o pacote tzdata): vale a tabela
        z = None
    if z is not None:
        real = utc.astimezone(z)
        assert real.utcoffset() == dt.timedelta(minutes=DESLOCAMENTO[fuso]), (fuso, instante, real)
    return f"{local.day:02d}/{local.month:02d}/{local.year:04d} {local.hour:02d}:{local.minute:02d}"


# ---------------------------------------------------------------- casos
# Cada fonte: (nome, pais, idioma, ...entradas, esperado). Em esperado, "~" é U+00A0.
def nb(s):
    return s.replace("~", NBSP)


inteiro_src = [
    ("zero", "BR", "pt-BR", "0", "0"),
    ("três dígitos sem ponto", "BR", "pt-BR", "999", "999"),
    ("quatro dígitos já agrupam", "BR", "pt-BR", "1000", "1.000"),
    ("quatro dígitos agrupam também na Argentina (não segue o mínimo de 2 grupos do CLDR es)", "AR", "es-419", "1000", "1.000"),
    ("quatro dígitos agrupam também no Paraguai", "PY", "es-419", "4812", "4.812"),
    ("seis dígitos", "BR", "pt-BR", "248315", "248.315"),
    ("sete dígitos", "AR", "es-419", "1234567", "1.234.567"),
    ("milhão redondo", "PY", "es-419", "1000000", "1.000.000"),
    ("negativo de três dígitos", "BR", "pt-BR", "-999", "-999"),
    ("negativo de quatro dígitos", "BR", "pt-BR", "-1000", "-1.000"),
    ("negativo grande", "AR", "es-419", "-1234567", "-1.234.567"),
    ("acima de 2^53 em texto", "BR", "pt-BR", "9007199254740993", "9.007.199.254.740.993"),
    ("usuário es-419 em empresa brasileira vê o formato do Brasil", "BR", "es-419", "12480", "12.480"),
    ("usuário pt-BR em empresa paraguaia vê o formato do Paraguai", "PY", "pt-BR", "12480", "12.480"),
]

decimal_src = [
    ("uma casa", "BR", "pt-BR", "3.2", "3,2"),
    ("zero com casas", "BR", "pt-BR", "0.41", "0,41"),
    ("zeros à direita do texto são mantidos", "BR", "pt-BR", "3.20", "3,20"),
    ("milhar com uma casa", "BR", "pt-BR", "1234.5", "1.234,5"),
    ("negativo com milhar", "BR", "pt-BR", "-1234.5", "-1.234,5"),
    ("negativo menor que um", "AR", "es-419", "-0.5", "-0,5"),
    ("parte decimal nunca é agrupada", "PY", "es-419", "0.123456", "0,123456"),
    ("parte decimal longa nunca é agrupada", "BR", "pt-BR", "1234567.891011", "1.234.567,891011"),
    ("decimal sem parte decimal", "AR", "es-419", "1000", "1.000"),
    ("acima de 2^53 com casas, sem passar por float", "BR", "pt-BR", "9007199254740993.25", "9.007.199.254.740.993,25"),
    ("muitas casas não são arredondadas", "PY", "es-419", "0.1000000000000000055511151231257827", "0,1000000000000000055511151231257827"),
    ("usuário es-419 em empresa brasileira: vírgula, não ponto", "BR", "es-419", "1234.5", "1.234,5"),
    ("usuário pt-BR em empresa argentina", "AR", "pt-BR", "1234.5", "1.234,5"),
]

percentual_src = [
    ("uma casa", "BR", "pt-BR", "97.4", "97,4~%"),
    ("inteiro", "BR", "pt-BR", "100", "100~%"),
    ("zero", "AR", "es-419", "0", "0~%"),
    ("negativo", "BR", "pt-BR", "-2.5", "-2,5~%"),
    ("milhar no percentual", "BR", "pt-BR", "1250", "1.250~%"),
    ("menor que um", "PY", "es-419", "0.5", "0,5~%"),
    ("duas casas", "AR", "es-419", "33.33", "33,33~%"),
    ("usuário es-419 em empresa brasileira", "BR", "es-419", "97.4", "97,4~%"),
    ("usuário pt-BR em empresa paraguaia", "PY", "pt-BR", "-2.5", "-2,5~%"),
]

medidor_src = [
    ("hodômetro em km", "BR", "pt-BR", "248315", "km", "248.315~km"),
    ("horímetro em h", "BR", "pt-BR", "12480", "h", "12.480~h"),
    ("ciclos", "AR", "es-419", "1204", "ciclos", "1.204~ciclos"),
    ("horímetro com casa decimal", "PY", "es-419", "1234.5", "h", "1.234,5~h"),
    ("três dígitos", "BR", "pt-BR", "850", "km", "850~km"),
    ("zero", "AR", "es-419", "0", "km", "0~km"),
    ("litros com L maiúsculo", "BR", "pt-BR", "45.37", "L", "45,37~L"),
    ("temperatura negativa", "AR", "es-419", "-12.5", "°C", "-12,5~°C"),
    ("pressão em bar", "PY", "es-419", "8.5", "bar", "8,5~bar"),
    ("unidade de estoque", "BR", "pt-BR", "1500", "un.", "1.500~un."),
    ("usuário es-419 em empresa brasileira", "BR", "es-419", "248315", "km", "248.315~km"),
]

valor_src = [
    ("real", "BR", "pt-BR", "BRL", "3418.90", "R$~3.418,90"),
    ("peso argentino", "AR", "es-419", "ARS", "3418.90", "$~3.418,90"),
    ("guarani sem casas", "PY", "es-419", "PYG", "3419", "Gs.~3.419"),
    ("real zero", "BR", "pt-BR", "BRL", "0.00", "R$~0,00"),
    ("guarani zero", "PY", "es-419", "PYG", "0", "Gs.~0"),
    ("real de cinco centavos", "BR", "pt-BR", "BRL", "0.05", "R$~0,05"),
    ("real abaixo do milhar", "BR", "pt-BR", "BRL", "999.99", "R$~999,99"),
    ("real de mil exato", "BR", "pt-BR", "BRL", "1000.00", "R$~1.000,00"),
    ("guarani de três dígitos", "PY", "es-419", "PYG", "500", "Gs.~500"),
    ("real negativo: sinal antes do símbolo", "BR", "pt-BR", "BRL", "-10.00", "-R$~10,00"),
    ("peso argentino negativo", "AR", "es-419", "ARS", "-10.00", "-$~10,00"),
    ("guarani negativo", "PY", "es-419", "PYG", "-10", "-Gs.~10"),
    ("peso argentino negativo com milhar", "AR", "es-419", "ARS", "-1234567.89", "-$~1.234.567,89"),
    ("real negativo de um centavo", "BR", "pt-BR", "BRL", "-0.01", "-R$~0,01"),
    ("real no máximo de numeric(18,2), acima de 2^53 em centavos", "BR", "pt-BR", "BRL", "9999999999999999.99", "R$~9.999.999.999.999.999,99"),
    ("real no mínimo", "BR", "pt-BR", "BRL", "-9999999999999999.99", "-R$~9.999.999.999.999.999,99"),
    ("guarani no máximo", "PY", "es-419", "PYG", "9999999999999999", "Gs.~9.999.999.999.999.999"),
    ("peso argentino em inflação", "AR", "es-419", "ARS", "125000000.00", "$~125.000.000,00"),
    ("usuário es-419 em empresa brasileira vê R$ e vírgula", "BR", "es-419", "BRL", "3418.90", "R$~3.418,90"),
    ("usuário pt-BR em empresa argentina vê $", "AR", "pt-BR", "ARS", "3418.90", "$~3.418,90"),
    ("usuário pt-BR em empresa paraguaia vê Gs.", "PY", "pt-BR", "PYG", "3419", "Gs.~3.419"),
]

preco_src = [
    ("real com duas casas", "BR", "pt-BR", "BRL", "45.50", "R$~45,50"),
    ("real com três casas (combustível)", "BR", "pt-BR", "BRL", "6.299", "R$~6,299"),
    ("real com quatro casas", "BR", "pt-BR", "BRL", "37.8462", "R$~37,8462"),
    ("real com seis casas", "BR", "pt-BR", "BRL", "0.000001", "R$~0,000001"),
    ("real zero", "BR", "pt-BR", "BRL", "0.00", "R$~0,00"),
    ("peso argentino com seis casas e milhar", "AR", "es-419", "ARS", "1234.123456", "$~1.234,123456"),
    ("peso argentino com duas casas", "AR", "es-419", "ARS", "15000.00", "$~15.000,00"),
    ("guarani inteiro", "PY", "es-419", "PYG", "25000", "Gs.~25.000"),
    ("meio guarani", "PY", "es-419", "PYG", "0.5", "Gs.~0,5"),
    ("guarani com seis casas", "PY", "es-419", "PYG", "0.123456", "Gs.~0,123456"),
    ("guarani zero", "PY", "es-419", "PYG", "0", "Gs.~0"),
    ("real negativo", "BR", "pt-BR", "BRL", "-1.50", "-R$~1,50"),
    ("guarani negativo com casas", "PY", "es-419", "PYG", "-0.5", "-Gs.~0,5"),
    ("real no máximo de numeric(18,6), acima de 2^53 em milionésimos", "BR", "pt-BR", "BRL", "999999999999.999999", "R$~999.999.999.999,999999"),
    ("usuário es-419 em empresa brasileira", "BR", "es-419", "BRL", "37.8462", "R$~37,8462"),
    ("usuário pt-BR em empresa paraguaia", "PY", "pt-BR", "PYG", "0.5", "Gs.~0,5"),
]

data_src = [
    ("dia e mês com dois dígitos", "BR", "pt-BR", "2026-11-12", "12/11/2026"),
    ("zero à esquerda no dia", "BR", "pt-BR", "2026-11-03", "03/11/2026"),
    ("zero à esquerda no dia e no mês", "AR", "es-419", "2026-07-05", "05/07/2026"),
    ("primeiro do ano não recua para o ano anterior (data sem fuso)", "BR", "pt-BR", "2027-01-01", "01/01/2027"),
    ("último do ano", "PY", "es-419", "2026-12-31", "31/12/2026"),
    ("29 de fevereiro", "PY", "es-419", "2028-02-29", "29/02/2028"),
    ("vencimento de documento", "AR", "es-419", "2027-03-31", "31/03/2027"),
    ("usuário es-419 em empresa brasileira: dd/mm/aaaa", "BR", "es-419", "2026-11-12", "12/11/2026"),
    ("usuário pt-BR em empresa argentina", "AR", "pt-BR", "2026-11-12", "12/11/2026"),
]

SP, MAN, RB, NOR, CUI = ("America/Sao_Paulo", "America/Manaus", "America/Rio_Branco",
                         "America/Noronha", "America/Cuiaba")
BA, COR, USH, ASU = ("America/Argentina/Buenos_Aires", "America/Argentina/Cordoba",
                     "America/Argentina/Ushuaia", "America/Asuncion")
data_hora_src = [
    ("São Paulo à tarde", "BR", "pt-BR", SP, "2026-07-15T17:05:00Z", "15/07/2026 14:05"),
    ("São Paulo: meia-noite é 00:00", "BR", "pt-BR", SP, "2026-07-01T03:00:00Z", "01/07/2026 00:00"),
    ("São Paulo: um minuto antes da meia-noite", "BR", "pt-BR", SP, "2026-07-01T02:59:00Z", "30/06/2026 23:59"),
    ("São Paulo: madrugada em UTC ainda é o dia anterior", "BR", "pt-BR", SP, "2026-07-01T02:30:00Z", "30/06/2026 23:30"),
    ("Manaus: o mesmo instante uma hora antes", "BR", "pt-BR", MAN, "2026-07-01T02:30:00Z", "30/06/2026 22:30"),
    ("Rio Branco: o mesmo instante duas horas antes", "BR", "pt-BR", RB, "2026-07-01T02:30:00Z", "30/06/2026 21:30"),
    ("Noronha: o mesmo instante já é o dia seguinte", "BR", "pt-BR", NOR, "2026-07-01T02:30:00Z", "01/07/2026 00:30"),
    ("Cuiabá à meia-noite", "BR", "pt-BR", CUI, "2026-08-10T04:00:00Z", "10/08/2026 00:00"),
    ("hora e minuto com zero à esquerda", "BR", "pt-BR", SP, "2026-07-05T08:07:00Z", "05/07/2026 05:07"),
    ("São Paulo: fim do mês na virada", "BR", "pt-BR", SP, "2026-06-01T02:59:00Z", "31/05/2026 23:59"),
    ("Buenos Aires à tarde", "AR", "es-419", BA, "2026-07-15T17:05:00Z", "15/07/2026 14:05"),
    ("Córdoba (província) à meia-noite", "AR", "es-419", COR, "2026-07-09T03:00:00Z", "09/07/2026 00:00"),
    ("Ushuaia (província) de manhã", "AR", "es-419", USH, "2026-06-15T12:34:00Z", "15/06/2026 09:34"),
    ("Buenos Aires: virada do ano", "AR", "es-419", BA, "2027-01-01T02:59:00Z", "31/12/2026 23:59"),
    ("Assunção no verão à tarde", "PY", "es-419", ASU, "2026-12-10T17:05:00Z", "10/12/2026 14:05"),
    ("Assunção: meia-noite da virada do ano", "PY", "es-419", ASU, "2027-01-01T03:00:00Z", "01/01/2027 00:00"),
    ("Assunção: um minuto antes da virada do ano", "PY", "es-419", ASU, "2027-01-01T02:59:00Z", "31/12/2026 23:59"),
    ("Assunção em 29 de fevereiro", "PY", "es-419", ASU, "2028-02-29T15:00:00Z", "29/02/2028 12:00"),
    ("Assunção em janeiro de manhã", "PY", "es-419", ASU, "2027-01-20T11:15:00Z", "20/01/2027 08:15"),
    ("usuário es-419 em empresa brasileira: 24 h e dd/mm/aaaa", "BR", "es-419", SP, "2026-07-15T17:05:00Z", "15/07/2026 14:05"),
    ("usuário pt-BR em empresa paraguaia", "PY", "pt-BR", ASU, "2026-12-10T17:05:00Z", "10/12/2026 14:05"),
    ("unidade em Manaus de empresa com fuso de São Paulo: vale o fuso dado", "BR", "pt-BR", MAN, "2026-07-15T17:05:00Z", "15/07/2026 13:05"),
]


def montar():
    grupos = {}

    def simples(nome_grupo, fonte, f):
        casos = []
        for nome, pais, idioma, entrada, esperado in fonte:
            esperado = nb(esperado)
            confere(f"{nome_grupo}: {nome}", f(entrada), esperado)
            casos.append({"nome": nome, "pais": pais, "idioma": idioma, "entrada": entrada,
                          "esperado": esperado})
        grupos[nome_grupo] = casos

    def inteiro(t):
        assert RE_INTEIRO.fullmatch(t), t
        return numero(t)

    simples("inteiro", inteiro_src, inteiro)
    simples("decimal", decimal_src, numero)
    simples("percentual", percentual_src, lambda t: numero(t) + NBSP + "%")

    casos = []
    for nome, pais, idioma, entrada, unidade, esperado in medidor_src:
        esperado = nb(esperado)
        confere(f"medidor: {nome}", numero(entrada) + NBSP + unidade, esperado)
        casos.append({"nome": nome, "pais": pais, "idioma": idioma, "entrada": entrada,
                      "unidade": unidade, "esperado": esperado})
    grupos["medidor"] = casos

    for grupo, fonte, especie in (("valor", valor_src, "valor"),
                                  ("preco_unitario", preco_src, "preco_unitario")):
        casos = []
        for nome, pais, idioma, moeda, entrada, esperado in fonte:
            esperado = nb(esperado)
            confere(f"{grupo}: {nome}", dinheiro(pais, moeda, entrada, especie), esperado)
            casos.append({"nome": nome, "pais": pais, "idioma": idioma, "moeda": moeda,
                          "entrada": entrada, "esperado": esperado})
        grupos[grupo] = casos

    simples("data", data_src, data)

    casos = []
    for nome, pais, idioma, fuso, instante, esperado in data_hora_src:
        assert PAIS_DO_FUSO[fuso] == pais, (nome, fuso)
        confere(f"data_hora: {nome}", data_hora(fuso, instante), esperado)
        casos.append({"nome": nome, "pais": pais, "idioma": idioma, "fuso": fuso,
                      "instante": instante, "esperado": esperado})
    grupos["data_hora"] = casos

    for g, cs in grupos.items():
        nomes = [c["nome"] for c in cs]
        assert len(nomes) == len(set(nomes)), f"nome repetido em {g}"
        for c in cs:
            assert c["idioma"] in IDIOMAS and c["pais"] in SIMBOLO, c
    return grupos


grupos = montar()
if falhas:
    print("\n".join(falhas))
    sys.exit(1)

doc = {
    "descricao": (
        "Casos compartilhados de formatação para pessoas (ADR-013, Parte G, documento 5; documento 7, "
        "seção 3.11). Gerados e conferidos por gerar_formatos.py, ao lado; não edite à mão. Executados por "
        "pacotes/dominio/src/formatos.ts (web e app) e pelo formatador do Go quando houver. O formato vem do "
        "pais da empresa; idioma é o do usuário e não muda o resultado (há casos com idioma trocado para "
        "provar). Entradas em texto: inteiro e decimal no texto canônico da API (ponto decimal, sem milhar, "
        "sem expoente), valor e preco_unitario no texto canônico da Parte F com a moeda ao lado, percentual "
        "em pontos percentuais (97.4 é 97,4 %), data civil AAAA-MM-DD sem fuso, instante ISO 8601 em UTC "
        "com Z e fuso IANA explícito. Saída: milhar com ponto a partir de 4 dígitos, vírgula decimal, parte "
        "decimal nunca agrupada e copiada como veio (sem arredondar nem completar), negativo com U+002D na "
        "frente (antes do símbolo na moeda), espaço inseparável U+00A0 entre símbolo e número, antes de % e "
        "antes da unidade do medidor; símbolo R$, $ ou Gs. pelo país; data dd/mm/aaaa; data e hora "
        "dd/mm/aaaa hh:mm em 24 h, sem segundos, com U+0020 entre data e hora. Instantes longe de mudança "
        "de horário: Brasil de maio a agosto, Paraguai de dezembro a fevereiro (a mesma saída na base tz "
        "antiga e na atual)."
    ),
    "versao": 1,
    "simbolos": SIMBOLO,
    **grupos,
}

blocos = []
for k, v in doc.items():
    if k in grupos:
        linhas = ",\n".join("    " + json.dumps(c, ensure_ascii=False) for c in v)
        blocos.append(f'  "{k}": [\n{linhas}\n  ]')
    else:
        blocos.append(f'  "{k}": ' + json.dumps(v, ensure_ascii=False))
# Caracteres invisíveis sempre escapados, para não sumirem num editor.
texto = ("{\n" + ",\n".join(blocos) + "\n}\n").replace(NBSP, "\\u00a0")
assert json.loads(texto)["valor"][0]["esperado"] == "R$" + NBSP + "3.418,90"

contagem = {g: len(c) for g, c in grupos.items()}
saida = Path(__file__).with_name("formatos.json")
if "--conferir" in sys.argv[1:]:
    atual = saida.read_text(encoding="utf-8") if saida.exists() else ""
    if atual != texto:
        print(f"{saida.name} diverge do gerador: rode sem --conferir e revise a diferença")
        sys.exit(1)
    print(f"{saida.name} confere com o gerador:", contagem, "total", sum(contagem.values()))
else:
    with open(saida, "w", encoding="utf-8", newline="\n") as f:
        f.write(texto)
    print("gravado", saida.name, contagem, "total", sum(contagem.values()))
