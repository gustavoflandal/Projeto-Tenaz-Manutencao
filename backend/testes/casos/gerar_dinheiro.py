"""Gera e confere backend/testes/casos/dinheiro.json (ADR-013, Parte F, itens 4 a 10).

O oráculo abaixo implementa a Parte F com Fraction (exato, só biblioteca padrão).
Todo caso escrito à mão com resultado esperado é conferido contra o oráculo;
divergência aborta sem gravar. Os casos são a fonte; o JSON é gerado.

Uso:
    python -I backend/testes/casos/gerar_dinheiro.py            # regrava dinheiro.json
    python -I backend/testes/casos/gerar_dinheiro.py --conferir # falha se o JSON divergir
"""
import json
import re
import sys
from fractions import Fraction
from pathlib import Path

CASAS = {"BRL": 2, "ARS": 2, "PYG": 0}
RE = {
    "valor": re.compile(r"-?(0|[1-9][0-9]{0,15})(\.[0-9]{1,2})?"),
    "preco_unitario": re.compile(r"-?(0|[1-9][0-9]{0,11})(\.[0-9]{1,6})?"),
}
ESCALA_PRECO = 6


def escala(especie, moeda):
    return CASAS[moeda] if especie == "valor" else ESCALA_PRECO


def max_unidades(especie, moeda):
    if especie == "valor":
        return 10 ** (16 + CASAS[moeda]) - 1
    return 10 ** 18 - 1


def ler(entrada, especie, moeda, aceita_negativo):
    if not RE[especie].fullmatch(entrada, 0) or not entrada.isascii():
        return None
    neg = entrada.startswith("-")
    if neg and not aceita_negativo:
        return None
    corpo = entrada[1:] if neg else entrada
    inteiro, _, frac = corpo.partition(".")
    e = escala(especie, moeda)
    if len(frac) > e and set(frac[e:]) - {"0"}:
        return None
    frac = (frac[:e]).ljust(e, "0")
    u = int(inteiro) * 10 ** e + (int(frac) if frac else 0)
    u = -u if neg else u
    if abs(u) > max_unidades(especie, moeda):
        return None
    return u


def canonico(u, especie, moeda):
    e = escala(especie, moeda)
    casas = CASAS[moeda]
    sinal = "-" if u < 0 else ""
    a = abs(u)
    inteiro, frac = divmod(a, 10 ** e)
    s = str(frac).rjust(e, "0") if e else ""
    if especie == "preco_unitario":
        s = s.rstrip("0")
        s = s.ljust(casas, "0")
    return sinal + str(inteiro) + ("." + s if s else "")


def arred(x: Fraction) -> int:
    """Meio para longe do zero."""
    sinal = -1 if x < 0 else 1
    a = abs(x)
    q, r = divmod(a.numerator, a.denominator)
    if 2 * r >= a.denominator:
        q += 1
    return sinal * q


def qtd(texto):
    return Fraction(texto)


def multiplicar(moeda, preco_txt, quantidade_txt):
    p = ler(preco_txt, "preco_unitario", moeda, True)
    assert p is not None, preco_txt
    x = Fraction(p) * qtd(quantidade_txt) / 10 ** (ESCALA_PRECO - CASAS[moeda])
    u = arred(x)
    if abs(u) > max_unidades("valor", moeda):
        return None, "fora_da_faixa"
    return u, ""


def somar(op, a, b):
    if a["moeda"] != b["moeda"]:
        return None, "moeda_diferente"
    ua = ler(a["valor"], "valor", a["moeda"], True)
    ub = ler(b["valor"], "valor", b["moeda"], True)
    assert ua is not None and ub is not None, (a, b)
    u = ua + ub if op == "somar" else ua - ub
    if abs(u) > max_unidades("valor", a["moeda"]):
        return None, "fora_da_faixa"
    return u, ""


RE_PESO = re.compile(r"(0|[1-9][0-9]*)(\.[0-9]+)?")  # decimal sem sinal (Parte F, item 10)


def construir_moeda(codigo):
    if codigo in CASAS:
        return CASAS[codigo], ""
    return None, "moeda_desconhecida"


def ratear(moeda, valor_txt, pesos_txt):
    t = ler(valor_txt, "valor", moeda, True)
    assert t is not None
    if not pesos_txt or not all(p.isascii() and RE_PESO.fullmatch(p) for p in pesos_txt):
        return None, "rateio_invalido"
    pesos = [Fraction(p) for p in pesos_txt]
    w = sum(pesos)
    if w <= 0:
        return None, "rateio_invalido"
    partes = _ratear(t, pesos, w)
    assert _ratear(-t, pesos, w) == [-p for p in partes]  # ratear -v é o oposto de ratear v
    return partes, ""


def _ratear(t, pesos, w):
    sinal = -1 if t < 0 else 1
    a = abs(t)
    ideais = [a * p / w for p in pesos]
    base = [x.numerator // x.denominator for x in ideais]
    restos = [x - b for x, b in zip(ideais, base)]
    sobra = a - sum(base)
    # maior resto; empate: menor índice (ordem de entrada)
    ordem = sorted(range(len(pesos)), key=lambda i: (-restos[i], i))
    for i in ordem[:sobra]:
        base[i] += 1
    partes = [sinal * b for b in base]
    assert sum(partes) == t
    return partes


falhas = []


def confere(nome, obtido, esperado):
    if obtido != esperado:
        falhas.append(f"{nome}: oráculo {obtido!r}, escrito {esperado!r}")


# ---------------------------------------------------------------- leitura
# (nome, especie, moeda, entrada, aceita_negativo, canonico esperado ou None)
V, P = "valor", "preco_unitario"
leitura_src = [
    # valor aceito
    ("BRL valor com duas casas", V, "BRL", "1250.40", False, "1250.40"),
    ("BRL valor com uma casa completa com zero", V, "BRL", "125.5", False, "125.50"),
    ("BRL valor sem casas", V, "BRL", "125", False, "125.00"),
    ("BRL valor zero", V, "BRL", "0", False, "0.00"),
    ("BRL valor zero com duas casas", V, "BRL", "0.00", False, "0.00"),
    ("BRL valor de um centavo", V, "BRL", "0.01", False, "0.01"),
    ("BRL valor negativo onde o campo aceita", V, "BRL", "-1.50", True, "-1.50"),
    ("BRL valor negativo de um centavo", V, "BRL", "-0.01", True, "-0.01"),
    ("BRL -0 vira zero sem sinal", V, "BRL", "-0", True, "0.00"),
    ("BRL -0.00 vira zero sem sinal", V, "BRL", "-0.00", True, "0.00"),
    ("BRL valor no limite superior de numeric(18,2)", V, "BRL", "9999999999999999.99", False, "9999999999999999.99"),
    ("BRL valor no limite inferior de numeric(18,2)", V, "BRL", "-9999999999999999.99", True, "-9999999999999999.99"),
    ("ARS valor com duas casas", V, "ARS", "1234567.89", False, "1234567.89"),
    ("ARS valor com uma casa", V, "ARS", "0.5", False, "0.50"),
    ("ARS valor no limite superior", V, "ARS", "9999999999999999.99", False, "9999999999999999.99"),
    ("PYG valor inteiro", V, "PYG", "75000", False, "75000"),
    ("PYG valor com duas casas zeradas", V, "PYG", "75000.00", False, "75000"),
    ("PYG valor com uma casa zerada", V, "PYG", "75000.0", False, "75000"),
    ("PYG valor zero com casas zeradas", V, "PYG", "0.00", False, "0"),
    ("PYG -0.00 vira zero sem sinal", V, "PYG", "-0.00", True, "0"),
    ("PYG valor negativo onde o campo aceita", V, "PYG", "-75000", True, "-75000"),
    ("PYG valor no limite superior", V, "PYG", "9999999999999999", False, "9999999999999999"),
    ("PYG valor no limite superior com casas zeradas", V, "PYG", "9999999999999999.00", False, "9999999999999999"),
    ("PYG valor no limite inferior", V, "PYG", "-9999999999999999", True, "-9999999999999999"),
    # valor recusado
    ("BRL valor com terceira casa", V, "BRL", "1.505", False, None),
    ("BRL valor de um milésimo", V, "BRL", "0.001", False, None),
    ("ARS valor com terceira casa", V, "ARS", "10.999", False, None),
    ("PYG valor com 50 cêntimos", V, "PYG", "75000.50", False, None),
    ("PYG valor com 5 cêntimos", V, "PYG", "75000.05", False, None),
    ("PYG valor com uma casa não zero", V, "PYG", "0.1", False, None),
    ("BRL notação científica", V, "BRL", "1e3", False, None),
    ("BRL notação científica maiúscula", V, "BRL", "1E3", False, None),
    ("BRL notação científica com casas", V, "BRL", "1.5e2", False, None),
    ("PYG notação científica", V, "PYG", "7.5E4", False, None),
    ("BRL vírgula decimal", V, "BRL", "1,50", False, None),
    ("BRL milhar com ponto e vírgula decimal", V, "BRL", "1.000,50", False, None),
    ("BRL milhar com vírgula", V, "BRL", "1,000.50", False, None),
    ("PYG milhar com ponto", V, "PYG", "75.000", False, None),
    ("BRL espaço como milhar", V, "BRL", "1 000.00", False, None),
    ("BRL espaço à esquerda", V, "BRL", " 1.00", False, None),
    ("BRL espaço à direita", V, "BRL", "1.00 ", False, None),
    ("BRL quebra de linha no fim", V, "BRL", "1.00\n", False, None),
    ("BRL tabulação no início", V, "BRL", "\t1.00", False, None),
    ("BRL espaço entre sinal e número", V, "BRL", "- 1.00", True, None),
    ("BRL sinal de mais", V, "BRL", "+1.00", False, None),
    ("BRL sinal de mais onde aceita negativo", V, "BRL", "+1.00", True, None),
    ("BRL zero à esquerda", V, "BRL", "01.00", False, None),
    ("BRL dois zeros", V, "BRL", "00", False, None),
    ("BRL zero à esquerda negativo", V, "BRL", "-01", True, None),
    ("PYG zero à esquerda", V, "PYG", "075000", False, None),
    ("BRL vazio", V, "BRL", "", False, None),
    ("PYG vazio", V, "PYG", "", False, None),
    ("BRL só ponto", V, "BRL", ".", False, None),
    ("BRL ponto sem casas", V, "BRL", "1.", False, None),
    ("BRL ponto sem parte inteira", V, "BRL", ".50", False, None),
    ("BRL negativo sem parte inteira", V, "BRL", "-.50", True, None),
    ("BRL só sinal", V, "BRL", "-", True, None),
    ("BRL dois sinais", V, "BRL", "--1", True, None),
    ("BRL dois pontos", V, "BRL", "1.0.0", False, None),
    ("BRL sublinhado de literal Go", V, "BRL", "1_000", False, None),
    ("BRL hexadecimal", V, "BRL", "0x1F", False, None),
    ("BRL Infinity", V, "BRL", "Infinity", False, None),
    ("BRL NaN", V, "BRL", "NaN", False, None),
    ("BRL com símbolo da moeda", V, "BRL", "R$ 1,00", False, None),
    ("BRL dígitos de largura total", V, "BRL", "１.00", False, None),
    ("BRL dígito árabe-índico", V, "BRL", "١", False, None),
    ("BRL 17 dígitos inteiros", V, "BRL", "10000000000000000", False, None),
    ("BRL 17 dígitos inteiros com casas", V, "BRL", "10000000000000000.00", False, None),
    ("BRL negativo com 17 dígitos inteiros", V, "BRL", "-10000000000000000", True, None),
    ("BRL 2^63 em texto", V, "BRL", "9223372036854775808", False, None),
    ("BRL 20 dígitos", V, "BRL", "99999999999999999999", False, None),
    ("ARS 17 dígitos inteiros", V, "ARS", "12345678901234567.00", False, None),
    ("PYG 17 dígitos inteiros", V, "PYG", "10000000000000000", False, None),
    ("BRL negativo onde o campo não aceita", V, "BRL", "-1.50", False, None),
    ("BRL negativo de um centavo onde o campo não aceita", V, "BRL", "-0.01", False, None),
    ("PYG negativo onde o campo não aceita", V, "PYG", "-75000", False, None),
    ("BRL -0 onde o campo não aceita negativo", V, "BRL", "-0", False, None),
    ("BRL -0.00 onde o campo não aceita negativo", V, "BRL", "-0.00", False, None),
    ("PYG -0 onde o campo não aceita negativo", V, "PYG", "-0", False, None),
    ("BRL terceira casa zerada passa do padrão do contrato", V, "BRL", "1.500", False, None),
    ("ARS terceira casa zerada passa do padrão do contrato", V, "ARS", "0.000", False, None),
    ("PYG três casas zeradas passam do padrão do contrato", V, "PYG", "75000.000", False, None),
    # preço unitário aceito
    ("BRL preço com duas casas", P, "BRL", "45.50", False, "45.50"),
    ("BRL preço com uma casa completa até a da moeda", P, "BRL", "45.5", False, "45.50"),
    ("BRL preço inteiro", P, "BRL", "45", False, "45.00"),
    ("BRL preço de combustível com três casas", P, "BRL", "6.299", False, "6.299"),
    ("BRL preço com quatro casas", P, "BRL", "37.8462", False, "37.8462"),
    ("BRL preço com zeros à direita além da moeda", P, "BRL", "37.846200", False, "37.8462"),
    ("BRL preço com seis casas zeradas", P, "BRL", "1.000000", False, "1.00"),
    ("BRL preço de um milionésimo", P, "BRL", "0.000001", False, "0.000001"),
    ("BRL preço zero", P, "BRL", "0", False, "0.00"),
    ("BRL preço negativo onde o campo aceita", P, "BRL", "-1.5", True, "-1.50"),
    ("BRL preço -0 vira zero sem sinal", P, "BRL", "-0.000000", True, "0.00"),
    ("BRL preço no limite superior de numeric(18,6)", P, "BRL", "999999999999.999999", False, "999999999999.999999"),
    ("BRL preço no limite inferior de numeric(18,6)", P, "BRL", "-999999999999.999999", True, "-999999999999.999999"),
    ("ARS preço com seis casas", P, "ARS", "1234.567891", False, "1234.567891"),
    ("ARS preço com uma casa", P, "ARS", "1000.1", False, "1000.10"),
    ("PYG preço inteiro", P, "PYG", "25000", False, "25000"),
    ("PYG preço de meio guarani", P, "PYG", "0.5", False, "0.5"),
    ("PYG preço de meio guarani com zero à direita", P, "PYG", "0.50", False, "0.5"),
    ("PYG preço inteiro com seis casas zeradas", P, "PYG", "25000.000000", False, "25000"),
    ("PYG preço com seis casas", P, "PYG", "0.123456", False, "0.123456"),
    ("PYG preço no limite superior", P, "PYG", "999999999999.999999", False, "999999999999.999999"),
    # preço unitário recusado
    ("BRL preço com sétima casa", P, "BRL", "0.1234567", False, None),
    ("BRL preço de um décimo de milionésimo", P, "BRL", "0.0000001", False, None),
    ("PYG preço com sétima casa", P, "PYG", "0.0000005", False, None),
    ("BRL preço com 13 dígitos inteiros", P, "BRL", "1000000000000", False, None),
    ("PYG preço com 13 dígitos inteiros", P, "PYG", "1000000000000", False, None),
    ("BRL preço com 20 dígitos", P, "BRL", "99999999999999999999", False, None),
    ("BRL preço em notação científica", P, "BRL", "1e-6", False, None),
    ("BRL preço com vírgula", P, "BRL", "6,299", False, None),
    ("BRL preço com espaço", P, "BRL", "6.299 ", False, None),
    ("BRL preço com sinal de mais", P, "BRL", "+6.299", False, None),
    ("BRL preço com zero à esquerda", P, "BRL", "06.299", False, None),
    ("BRL preço vazio", P, "BRL", "", False, None),
    ("BRL preço só ponto", P, "BRL", ".", False, None),
    ("BRL preço negativo onde o campo não aceita (custo_hora)", P, "BRL", "-45.50", False, None),
    ("PYG preço negativo onde o campo não aceita", P, "PYG", "-0.5", False, None),
    ("BRL preço -0 onde o campo não aceita negativo", P, "BRL", "-0", False, None),
    ("BRL preço com sétima casa zerada passa do padrão do contrato", P, "BRL", "0.1234560", False, None),
    ("PYG preço com sete casas zeradas passa do padrão do contrato", P, "PYG", "25000.0000000", False, None),
]

leitura = []
for nome, esp, moeda, ent, neg, canon in leitura_src:
    u = ler(ent, esp, moeda, neg)
    if canon is None:
        confere(nome, u, None)
        leitura.append({"nome": nome, "especie": esp, "moeda": moeda, "entrada": ent,
                        "aceita_negativo": neg, "valido": False, "unidades": "",
                        "canonico": "", "erro": "formato_invalido"})
    else:
        if u is None:
            falhas.append(f"{nome}: oráculo recusou")
            continue
        confere(nome, canonico(u, esp, moeda), canon)
        leitura.append({"nome": nome, "especie": esp, "moeda": moeda, "entrada": ent,
                        "aceita_negativo": neg, "valido": True, "unidades": str(u),
                        "canonico": canon, "erro": ""})

# ---------------------------------------------------------------- canônico
canon_src = [
    ("BRL valor de 1,50", V, "BRL", 150, "1.50"),
    ("BRL valor de 1.250,40", V, "BRL", 125040, "1250.40"),
    ("BRL valor zero", V, "BRL", 0, "0.00"),
    ("BRL valor de um centavo", V, "BRL", 1, "0.01"),
    ("BRL valor de dez centavos", V, "BRL", 10, "0.10"),
    ("BRL valor negativo de cinco centavos", V, "BRL", -5, "-0.05"),
    ("BRL valor redondo", V, "BRL", 100000, "1000.00"),
    ("BRL valor máximo", V, "BRL", 999999999999999999, "9999999999999999.99"),
    ("BRL valor mínimo", V, "BRL", -999999999999999999, "-9999999999999999.99"),
    ("ARS valor de 1.234.567,89", V, "ARS", 123456789, "1234567.89"),
    ("ARS valor negativo", V, "ARS", -100, "-1.00"),
    ("PYG valor de 75.000", V, "PYG", 75000, "75000"),
    ("PYG valor zero", V, "PYG", 0, "0"),
    ("PYG valor de um guarani", V, "PYG", 1, "1"),
    ("PYG valor negativo", V, "PYG", -1, "-1"),
    ("PYG valor máximo", V, "PYG", 9999999999999999, "9999999999999999"),
    ("BRL preço de 45,50", P, "BRL", 45500000, "45.50"),
    ("BRL preço de 37,8462", P, "BRL", 37846200, "37.8462"),
    ("BRL preço de combustível 6,299", P, "BRL", 6299000, "6.299"),
    ("BRL preço inteiro", P, "BRL", 45000000, "45.00"),
    ("BRL preço de um milionésimo", P, "BRL", 1, "0.000001"),
    ("BRL preço negativo de um milionésimo", P, "BRL", -1, "-0.000001"),
    ("BRL preço de dez centavos", P, "BRL", 100000, "0.10"),
    ("BRL preço zero", P, "BRL", 0, "0.00"),
    ("BRL preço máximo", P, "BRL", 999999999999999999, "999999999999.999999"),
    ("ARS preço com seis casas", P, "ARS", 1234567891, "1234.567891"),
    ("ARS preço com uma casa significativa", P, "ARS", 1000100000, "1000.10"),
    ("PYG preço de 25.000", P, "PYG", 25000000000, "25000"),
    ("PYG preço de meio guarani", P, "PYG", 500000, "0.5"),
    ("PYG preço com seis casas", P, "PYG", 123456, "0.123456"),
    ("PYG preço com três casas", P, "PYG", 1250000, "1.25"),
    ("PYG preço zero", P, "PYG", 0, "0"),
    ("PYG preço negativo de meio guarani", P, "PYG", -500000, "-0.5"),
]
texto_canonico = []
for nome, esp, moeda, u, esperado in canon_src:
    confere(nome, canonico(u, esp, moeda), esperado)
    texto_canonico.append({"nome": nome, "especie": esp, "moeda": moeda,
                           "unidades": str(u), "esperado": esperado})

# ---------------------------------------------------------------- soma/subtração
def op(nome, operacao, ma, va, mb, vb, esperado, erro=""):
    return (nome, operacao, {"moeda": ma, "valor": va}, {"moeda": mb, "valor": vb}, esperado, erro)

soma_src = [
    op("BRL soma simples", "somar", "BRL", "1.50", "BRL", "2.75", "4.25"),
    op("BRL 0,10 + 0,20 é exatamente 0,30", "somar", "BRL", "0.10", "BRL", "0.20", "0.30"),
    op("BRL soma com negativo", "somar", "BRL", "-0.05", "BRL", "0.05", "0.00"),
    op("BRL soma de dois negativos", "somar", "BRL", "-1.25", "BRL", "-2.50", "-3.75"),
    op("BRL soma de zero no limite", "somar", "BRL", "9999999999999999.99", "BRL", "0", "9999999999999999.99"),
    op("BRL soma que chega exatamente ao limite", "somar", "BRL", "9999999999999999.98", "BRL", "0.01", "9999999999999999.99"),
    op("BRL soma que passa do limite por um centavo", "somar", "BRL", "9999999999999999.99", "BRL", "0.01", "", "fora_da_faixa"),
    op("BRL soma de dois máximos", "somar", "BRL", "9999999999999999.99", "BRL", "9999999999999999.99", "", "fora_da_faixa"),
    op("BRL soma que passa do limite inferior", "somar", "BRL", "-9999999999999999.99", "BRL", "-0.01", "", "fora_da_faixa"),
    op("ARS soma", "somar", "ARS", "1234567.89", "ARS", "0.11", "1234568.00"),
    op("ARS soma que passa do limite", "somar", "ARS", "5000000000000000.00", "ARS", "5000000000000000.00", "", "fora_da_faixa"),
    op("PYG soma", "somar", "PYG", "75000", "PYG", "25000", "100000"),
    op("PYG soma no limite", "somar", "PYG", "9999999999999998", "PYG", "1", "9999999999999999"),
    op("PYG soma que passa do limite por um guarani", "somar", "PYG", "9999999999999999", "PYG", "1", "", "fora_da_faixa"),
    op("BRL mais ARS", "somar", "BRL", "1.00", "ARS", "1.00", "", "moeda_diferente"),
    op("BRL mais PYG", "somar", "BRL", "1.00", "PYG", "1", "", "moeda_diferente"),
    op("zero BRL mais zero ARS também é recusado", "somar", "BRL", "0", "ARS", "0", "", "moeda_diferente"),
    op("BRL subtração simples", "subtrair", "BRL", "10.00", "BRL", "2.51", "7.49"),
    op("BRL subtração que fica negativa", "subtrair", "BRL", "1.00", "BRL", "2.50", "-1.50"),
    op("BRL subtração que zera não deixa -0", "subtrair", "BRL", "1.00", "BRL", "1.00", "0.00"),
    op("BRL subtração de negativo", "subtrair", "BRL", "1.00", "BRL", "-1.00", "2.00"),
    op("BRL subtração que passa do limite inferior", "subtrair", "BRL", "-9999999999999999.99", "BRL", "0.01", "", "fora_da_faixa"),
    op("BRL máximo menos mínimo", "subtrair", "BRL", "9999999999999999.99", "BRL", "-9999999999999999.99", "", "fora_da_faixa"),
    op("BRL subtração no limite inferior exato", "subtrair", "BRL", "-9999999999999999.98", "BRL", "0.01", "-9999999999999999.99"),
    op("PYG subtração", "subtrair", "PYG", "100000", "PYG", "75001", "24999"),
    op("PYG subtração de negativo que passa do limite", "subtrair", "PYG", "9999999999999999", "PYG", "-1", "", "fora_da_faixa"),
    op("ARS menos PYG", "subtrair", "ARS", "10.00", "PYG", "10", "", "moeda_diferente"),
    op("PYG menos BRL", "subtrair", "PYG", "10", "BRL", "10.00", "", "moeda_diferente"),
]
soma_subtracao = []
for nome, operacao, a, b, esperado, erro in soma_src:
    u, e = somar(operacao, a, b)
    confere(nome + " (erro)", e, erro)
    if not erro:
        confere(nome, canonico(u, V, a["moeda"]), esperado)
    soma_subtracao.append({"nome": nome, "operacao": operacao, "a": a, "b": b,
                           "esperado": esperado, "erro": erro})

# ---------------------------------------------------------------- multiplicação
P_MAX = "999999999999.999999"
mult_src = [
    ("BRL preço por quantidade inteira", "BRL", "45.50", "2", "91.00"),
    ("BRL quatro casas no preço, arredonda para cima", "BRL", "37.8462", "3", "113.54"),
    ("BRL empate positivo 2,345 vai para 2,35", "BRL", "2.345", "1", "2.35"),
    ("BRL empate negativo -2,345 vai para -2,35", "BRL", "-2.345", "1", "-2.35"),
    ("BRL empate 0,125 vai para 0,13 (não para o par)", "BRL", "0.125", "1", "0.13"),
    ("BRL empate negativo -0,125 vai para -0,13", "BRL", "-0.125", "1", "-0.13"),
    ("BRL logo abaixo do empate arredonda para baixo", "BRL", "0.124999", "1", "0.12"),
    ("BRL negativo logo abaixo do empate em módulo vai para -0,12", "BRL", "-0.124999", "1", "-0.12"),
    ("BRL empate vindo da quantidade", "BRL", "1.00", "0.005", "0.01"),
    ("BRL empate negativo vindo da quantidade", "BRL", "1.00", "-0.005", "-0.01"),
    ("BRL combustível 6,299 por 45,37 litros", "BRL", "6.299", "45.37", "285.79"),
    ("BRL um só arredondamento: preço abaixo do centavo vezes mil", "BRL", "0.004", "1000", "4.00"),
    ("BRL um só arredondamento: 0,333333 vezes 3", "BRL", "0.333333", "3", "1.00"),
    ("BRL resultado negativo que arredonda a zero não vira -0", "BRL", "-0.004", "1", "0.00"),
    ("BRL quantidade zero", "BRL", "45.50", "0", "0.00"),
    ("BRL preço zero", "BRL", "0", "12.5", "0.00"),
    ("BRL custo-hora por 50 minutos (50/60)", "BRL", "45.50", "50/60", "37.92"),
    ("BRL custo-hora por 20 minutos dá dízima (1/3)", "BRL", "10.00", "20/60", "3.33"),
    ("BRL custo-hora por 40 minutos dá dízima (2/3)", "BRL", "10.00", "40/60", "6.67"),
    ("BRL custo-hora com empate em 5 minutos", "BRL", "0.30", "5/60", "0.03"),
    ("BRL custo-hora de seis casas por 1 minuto", "BRL", "45.123456", "1/60", "0.75"),
    ("BRL custo-hora por 0 minutos", "BRL", "45.50", "0/60", "0.00"),
    ("BRL custo-hora por 600 horas", "BRL", "45.50", "36000/60", "27300.00"),
    ("ARS preço de seis casas por quantidade decimal", "ARS", "1234.567890", "2.5", "3086.42"),
    ("ARS empate", "ARS", "0.015", "1", "0.02"),
    ("ARS custo-hora por 90 minutos", "ARS", "15000.00", "90/60", "22500.00"),
    ("PYG preço inteiro por quantidade inteira", "PYG", "25000", "3", "75000"),
    ("PYG meio guarani por 1 empata e vai para 1", "PYG", "0.5", "1", "1"),
    ("PYG meio guarani por 3 empata e vai para 2", "PYG", "0.5", "3", "2"),
    ("PYG empate negativo -0,5 vai para -1", "PYG", "-0.5", "1", "-1"),
    ("PYG empate negativo -1,5 vai para -2", "PYG", "-0.5", "3", "-2"),
    ("PYG 0,4 arredonda a zero", "PYG", "0.4", "1", "0"),
    ("PYG -0,4 arredonda a zero sem sinal", "PYG", "-0.4", "1", "0"),
    ("PYG preço de grama por 1.250 gramas", "PYG", "0.123456", "1250", "154"),
    ("PYG custo-hora por 50 minutos", "PYG", "75000", "50/60", "62500"),
    ("PYG custo-hora por 7 minutos (dízima)", "PYG", "75000", "7/60", "8750"),
    ("PYG custo-hora quebrado por 1 minuto empata", "PYG", "90", "1/60", "2"),
    ("PYG quantidade negativa de ajuste de inventário", "PYG", "25000", "-3", "-75000"),
    ("BRL quantidade negativa em fração com empate", "BRL", "0.30", "-5/60", "-0.03"),
    ("BRL preço máximo por 10.000 chega exatamente ao limite do valor", "BRL", P_MAX, "10000", "9999999999999999.99"),
    ("BRL preço máximo por 10.001 passa do limite", "BRL", P_MAX, "10001", None),
    ("BRL arredondamento que leva além do limite", "BRL", P_MAX, "1000000000000001/100000000000", None),
    ("BRL preço mínimo por 10.001 passa do limite inferior", "BRL", "-" + P_MAX, "10001", None),
    ("PYG preço máximo por 10.000 passa do limite do guarani", "PYG", P_MAX, "10000", "9999999999999999.99_PYG"),
    ("PYG preço inteiro máximo por 10.001 passa do limite", "PYG", "999999999999", "10001", None),
    ("BRL preço máximo por 60 minutos: produto intermediário acima de int64", "BRL", P_MAX, "60/60", "1000000000000.00"),
]
multiplicacao = []
for nome, moeda, preco, q, esperado in mult_src:
    u, e = multiplicar(moeda, preco, q)
    if esperado == "9999999999999999.99_PYG":
        # P_MAX x 10000 = 9999999999999999.99 PYG -> arredonda a 10^16 -> fora
        esperado = None
    if esperado is None:
        confere(nome, e, "fora_da_faixa")
        multiplicacao.append({"nome": nome, "moeda": moeda, "preco": preco, "quantidade": q,
                              "esperado": "", "unidades": "", "erro": "fora_da_faixa"})
    else:
        if e:
            falhas.append(f"{nome}: oráculo deu {e}")
            continue
        confere(nome, canonico(u, V, moeda), esperado)
        multiplicacao.append({"nome": nome, "moeda": moeda, "preco": preco, "quantidade": q,
                              "esperado": esperado, "unidades": str(u), "erro": ""})

# ---------------------------------------------------------------- rateio
rateio_src = [
    ("BRL 100,00 em três partes iguais: sobra vai para a primeira", "BRL", "100.00", ["1", "1", "1"], ["33.34", "33.33", "33.33"]),
    ("BRL 0,02 em três partes iguais", "BRL", "0.02", ["1", "1", "1"], ["0.01", "0.01", "0.00"]),
    ("BRL 0,03 em cinco partes iguais: mais partes que centavos", "BRL", "0.03", ["1", "1", "1", "1", "1"], ["0.01", "0.01", "0.01", "0.00", "0.00"]),
    ("BRL 10,00 nos pesos 1 e 2: sobra vai para o maior resto", "BRL", "10.00", ["1", "2"], ["3.33", "6.67"]),
    ("BRL empate de restos fora da primeira parte: vai para a de menor índice entre as empatadas", "BRL", "0.05", ["4", "3", "3"], ["0.02", "0.02", "0.01"]),
    ("BRL a ordem das partes decide o empate", "BRL", "0.05", ["3", "3", "4"], ["0.02", "0.01", "0.02"]),
    ("BRL maior resto vence a ordem", "BRL", "0.11", ["1", "3", "3"], ["0.01", "0.05", "0.05"]),
    ("BRL percentuais com casas: 999,99 em 33,33/33,33/33,34", "BRL", "999.99", ["33.33", "33.33", "33.34"], ["333.30", "333.30", "333.39"]),
    ("BRL divisão exata não distribui sobra", "BRL", "90.00", ["1", "2", "3", "3"], ["10.00", "20.00", "30.00", "30.00"]),
    ("BRL peso zero recebe zero e nunca a sobra", "BRL", "1.01", ["0", "1", "1"], ["0.00", "0.51", "0.50"]),
    ("BRL uma parte só recebe o total", "BRL", "123.45", ["7"], ["123.45"]),
    ("BRL valor zero", "BRL", "0.00", ["1", "1", "1"], ["0.00", "0.00", "0.00"]),
    ("BRL valor negativo: rateia o módulo e devolve o sinal", "BRL", "-100.00", ["1", "1", "1"], ["-33.34", "-33.33", "-33.33"]),
    ("BRL valor negativo nos pesos 1 e 2", "BRL", "-10.00", ["1", "2"], ["-3.33", "-6.67"]),
    ("BRL máximo em três partes: produto intermediário acima de int64", "BRL", "9999999999999999.99", ["3333", "3333", "3334"], None),
    ("ARS 1.000,00 em sete partes iguais", "ARS", "1000.00", ["1"] * 7, ["142.86", "142.86", "142.86", "142.86", "142.86", "142.85", "142.85"]),
    ("ARS pesos por horas apontadas com casas", "ARS", "50000.00", ["1.5", "2.25", "0.75"], None),
    ("PYG 100.000 em três partes iguais", "PYG", "100000", ["1", "1", "1"], ["33334", "33333", "33333"]),
    ("PYG 10 em quatro partes iguais: sobra de duas unidades", "PYG", "10", ["1", "1", "1", "1"], ["3", "3", "2", "2"]),
    ("PYG 1 em três partes", "PYG", "1", ["1", "1", "1"], ["1", "0", "0"]),
    ("PYG valor negativo", "PYG", "-10", ["1", "1", "1", "1"], ["-3", "-3", "-2", "-2"]),
    ("PYG máximo em pesos grandes", "PYG", "9999999999999999", ["3333", "3333", "3334"], None),
    ("BRL pesos iguais em texto decimal com casas", "BRL", "0.10", ["1.0", "1.00", "1"], ["0.04", "0.03", "0.03"]),
    ("BRL rateio sem partes", "BRL", "100.00", [], "rateio_invalido"),
    ("BRL pesos todos zero", "BRL", "100.00", ["0", "0", "0"], "rateio_invalido"),
    ("BRL peso único zero", "BRL", "100.00", ["0"], "rateio_invalido"),
    ("BRL peso negativo mesmo com soma positiva", "BRL", "100.00", ["1", "-1", "1"], "rateio_invalido"),
    ("BRL peso -0 tem sinal e é recusado", "BRL", "100.00", ["-0", "1"], "rateio_invalido"),
    ("PYG peso negativo", "PYG", "100000", ["2", "-1"], "rateio_invalido"),
]
rateio = []
for nome, moeda, valor, pesos, esperado in rateio_src:
    u, e = ratear(moeda, valor, pesos)
    if isinstance(esperado, str):
        confere(nome, e, esperado)
        rateio.append({"nome": nome, "moeda": moeda, "valor": valor, "pesos": pesos,
                       "esperado": [], "erro": esperado})
        continue
    if e:
        falhas.append(f"{nome}: oráculo deu {e}")
        continue
    partes = [canonico(x, V, moeda) for x in u]
    if esperado is not None:
        confere(nome, partes, esperado)
    rateio.append({"nome": nome, "moeda": moeda, "valor": valor, "pesos": pesos,
                   "esperado": partes, "erro": ""})

# ---------------------------------------------------------------- construção de Moeda
moeda_src = [
    ("real", "BRL", 2),
    ("peso argentino", "ARS", 2),
    ("guarani", "PYG", 0),
    ("dólar não está no catálogo", "USD", None),
    ("código em minúsculas", "brl", None),
    ("código com espaço", " BRL", None),
    ("código vazio", "", None),
    ("código numérico ISO 4217 do real", "986", None),
    ("país no lugar da moeda", "BR", None),
]
construcao_moeda = []
for nome, codigo, casas in moeda_src:
    c, e = construir_moeda(codigo)
    if casas is None:
        confere(nome, e, "moeda_desconhecida")
        construcao_moeda.append({"nome": nome, "codigo": codigo, "valido": False, "casas": 0,
                                 "erro": "moeda_desconhecida"})
    else:
        confere(nome, c, casas)
        construcao_moeda.append({"nome": nome, "codigo": codigo, "valido": True, "casas": casas,
                                 "erro": ""})

if falhas:
    print("\n".join(falhas))
    sys.exit(1)

doc = {
    "descricao": (
        "Casos compartilhados de dinheiro (ADR-013, Parte F, itens 4 a 10, documento 5). Gerados e conferidos "
        "por gerar_dinheiro.py, ao lado; não edite à mão. Executados pelo Go (backend/internal/plataforma/dinheiro) "
        "e pelo TypeScript (pacotes/dominio/src/dinheiro.ts). "
        "construcao_moeda: só a construção de Moeda a partir do código falha, com moeda_desconhecida (interno); "
        "as demais operações recebem o tipo Moeda já construído. "
        "Espécies: valor (escala da moeda, numeric(18,2)) e preco_unitario (milionésimos, numeric(18,6)). "
        "unidades: inteiro em texto, na unidade menor da moeda para valor e em milionésimos da unidade da moeda "
        "para preco_unitario (texto porque passa de 2^53). "
        "leitura: a gramática é exatamente a do padrão do contrato (sem -? quando aceita_negativo é false), "
        "depois a escala da moeda; toda recusa é formato_invalido, sem arredondar. "
        "texto_canonico: valor com exatamente as casas da moeda; preço com no mínimo as casas da moeda e no "
        "máximo seis, sem zeros à direita além das casas da moeda; nunca -0. "
        "soma_subtracao: operandos em texto canônico, lidos com negativo permitido; erro moeda_diferente "
        "(interno) ou fora_da_faixa (resultado fora de numeric(18,2) na escala da moeda). "
        "multiplicacao: preco (preço unitário) vezes quantidade, uma fração exata que pode ser negativa; no "
        "arquivo vem como texto decimal ou a/b (minutos/60), e o leitor dessa notação fica no código de teste; "
        "um só arredondamento meio para longe do zero na casa da moeda; erro fora_da_faixa. "
        "rateio: pesos em texto decimal sem sinal, lista não vazia e soma positiva, senão rateio_invalido "
        "(interno); cada parte recebe o piso do quinhão sobre o módulo do valor, a sobra de unidades menores vai "
        "uma a uma pelo maior resto, empate pelo menor índice, e o sinal do valor volta a todas as partes; "
        "a soma das partes é sempre o valor. erro vazio quando não há erro."
    ),
    "versao": 2,
    "moedas": [{"codigo": m, "casas": c} for m, c in CASAS.items()],
    "limites": {
        "valor_max_unidades": {m: str(max_unidades(V, m)) for m in CASAS},
        "preco_unitario_max_unidades": str(max_unidades(P, "BRL")),
    },
    "construcao_moeda": construcao_moeda,
    "leitura": leitura,
    "texto_canonico": texto_canonico,
    "soma_subtracao": soma_subtracao,
    "multiplicacao": multiplicacao,
    "rateio": rateio,
}

GRUPOS = ("construcao_moeda", "leitura", "texto_canonico", "soma_subtracao", "multiplicacao", "rateio")
nomes = {}
for grupo in GRUPOS:
    vistos = set()
    for c in doc[grupo]:
        assert c["nome"] not in vistos, (grupo, c["nome"])
        vistos.add(c["nome"])
    nomes[grupo] = len(doc[grupo])

blocos = []
for k, v in doc.items():
    if k in GRUPOS:
        linhas = ",\n".join("    " + json.dumps(c, ensure_ascii=False) for c in v)
        blocos.append(f'  "{k}": [\n{linhas}\n  ]')
    else:
        blocos.append(f'  "{k}": ' + json.dumps(v, ensure_ascii=False))
texto = "{\n" + ",\n".join(blocos) + "\n}\n"

saida = Path(__file__).with_name("dinheiro.json")
if "--conferir" in sys.argv[1:]:
    atual = saida.read_text(encoding="utf-8") if saida.exists() else ""
    if atual != texto:
        print(f"{saida.name} diverge do gerador: rode sem --conferir e revise a diferença")
        sys.exit(1)
    print(f"{saida.name} confere com o gerador:", nomes, "total", sum(nomes.values()))
else:
    with open(saida, "w", encoding="utf-8", newline="\n") as f:
        f.write(texto)
    print("gravado", saida.name, nomes, "total", sum(nomes.values()))
