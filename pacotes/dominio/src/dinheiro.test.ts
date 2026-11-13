import { describe, expect, it } from 'vitest';

import casos from '../../../backend/testes/casos/dinheiro.json';
import {
  ErroDeDinheiro,
  lerPrecoUnitario,
  lerValor,
  moedaDoCodigo,
  multiplicar,
  quantidadeDeTexto,
  ratear,
  somar,
  subtrair,
  textoDoPrecoUnitario,
  textoDoValor,
  type PrecoUnitario,
  type Quantidade,
  type Valor,
} from './dinheiro';

/** Notação `a/b` da quantidade: existe só no arquivo de casos (ADR-013, Parte F, item 10). */
function lerQuantidade(texto: string): Quantidade {
  const partes = texto.split('/');
  if (partes.length === 2) return { numerador: BigInt(partes[0] ?? ''), denominador: BigInt(partes[1] ?? '') };
  return quantidadeDeTexto(texto);
}

function codigoDoErro(acao: () => unknown): string {
  try {
    acao();
  } catch (erro) {
    if (erro instanceof ErroDeDinheiro) return erro.codigo;
    throw erro;
  }
  return '';
}

const valor = (m: string, t: string): Valor => lerValor(t, moedaDoCodigo(m), { aceitaNegativo: true });

describe('dinheiro (casos compartilhados com o Go, versão 2)', () => {
  it('usa a versão de casos que conhece e nenhum grupo está vazio', () => {
    expect(casos.versao).toBe(2);
    for (const grupo of ['moedas', 'construcao_moeda', 'leitura', 'texto_canonico', 'soma_subtracao', 'multiplicacao', 'rateio'] as const) {
      expect(casos[grupo].length, grupo).toBeGreaterThan(0);
    }
    const total = casos.construcao_moeda.length + casos.leitura.length + casos.texto_canonico.length + casos.soma_subtracao.length + casos.multiplicacao.length + casos.rateio.length;
    expect(total).toBe(266);
  });

  it.each(casos.construcao_moeda)('constrói a moeda: $nome', ({ codigo, valido, casas, erro }) => {
    if (valido) expect(moedaDoCodigo(codigo).casas).toBe(casas);
    else expect(codigoDoErro(() => moedaDoCodigo(codigo))).toBe(erro);
  });

  it.each(casos.leitura)('lê: $nome', (c) => {
    const moeda = moedaDoCodigo(c.moeda);
    const ler = () =>
      c.especie === 'valor'
        ? lerValor(c.entrada, moeda, { aceitaNegativo: c.aceita_negativo })
        : lerPrecoUnitario(c.entrada, moeda, { aceitaNegativo: c.aceita_negativo });
    if (!c.valido) {
      expect(codigoDoErro(ler)).toBe(c.erro);
      return;
    }
    const lido = ler();
    expect(lido.unidades.toString()).toBe(c.unidades);
    expect(lido.especie === 'valor' ? textoDoValor(lido) : textoDoPrecoUnitario(lido)).toBe(c.canonico);
  });

  it.each(casos.texto_canonico)('texto canônico: $nome', (c) => {
    const moeda = moedaDoCodigo(c.moeda);
    const unidades = BigInt(c.unidades);
    const texto =
      c.especie === 'valor'
        ? textoDoValor({ especie: 'valor', moeda, unidades })
        : textoDoPrecoUnitario({ especie: 'preco_unitario', moeda, unidades } satisfies PrecoUnitario);
    expect(texto).toBe(c.esperado);
  });

  it.each(casos.soma_subtracao)('soma e subtração: $nome', (c) => {
    const operar = () => {
      const a = valor(c.a.moeda, c.a.valor);
      const b = valor(c.b.moeda, c.b.valor);
      return c.operacao === 'somar' ? somar(a, b) : subtrair(a, b);
    };
    if (c.erro) expect(codigoDoErro(operar)).toBe(c.erro);
    else expect(textoDoValor(operar())).toBe(c.esperado);
  });

  it.each(casos.multiplicacao)('multiplica: $nome', (c) => {
    const moeda = moedaDoCodigo(c.moeda);
    const calcular = () => multiplicar(lerPrecoUnitario(c.preco, moeda, { aceitaNegativo: true }), lerQuantidade(c.quantidade));
    if (c.erro) {
      expect(codigoDoErro(calcular)).toBe(c.erro);
      return;
    }
    const resultado = calcular();
    expect(textoDoValor(resultado)).toBe(c.esperado);
    expect(resultado.unidades.toString()).toBe(c.unidades);
  });

  it.each(casos.rateio)('rateia: $nome', (c) => {
    const repartir = () => ratear(valor(c.moeda, c.valor), c.pesos);
    if (c.erro) {
      expect(codigoDoErro(repartir)).toBe(c.erro);
      return;
    }
    const partes = repartir();
    expect(partes.map(textoDoValor)).toEqual(c.esperado);
    expect(partes.reduce((s, p) => s + p.unidades, BigInt(0))).toBe(valor(c.moeda, c.valor).unidades);
  });
});

describe('dinheiro: moedas e limites do arquivo de casos', () => {
  it('as moedas do arquivo são exatamente o catálogo do código, com as mesmas casas', () => {
    expect(casos.moedas.map((m) => m.codigo)).toEqual(['BRL', 'ARS', 'PYG']);
    for (const m of casos.moedas) expect(moedaDoCodigo(m.codigo).casas, m.codigo).toBe(m.casas);
  });

  it.each(casos.moedas)('limite do valor em $codigo: o máximo é lido e escrito, e um a mais sai da faixa', ({ codigo }) => {
    const moeda = moedaDoCodigo(codigo);
    const maximo = BigInt(casos.limites.valor_max_unidades[codigo as keyof typeof casos.limites.valor_max_unidades]);
    const nele: Valor = { especie: 'valor', moeda, unidades: maximo };
    const texto = textoDoValor(nele);
    expect(lerValor(texto, moeda).unidades).toBe(maximo);
    expect(lerValor(`-${texto}`, moeda, { aceitaNegativo: true }).unidades).toBe(-maximo);
    const um: Valor = { especie: 'valor', moeda, unidades: BigInt(1) };
    expect(codigoDoErro(() => somar(nele, um))).toBe('fora_da_faixa');
    // Um dígito a mais que o padrão do contrato é formato inválido na leitura.
    expect(codigoDoErro(() => lerValor(`1${texto}`, moeda))).toBe('formato_invalido');
  });

  it.each(casos.moedas)('limite do preço unitário em $codigo: o máximo é lido e escrito, e passa da faixa ao multiplicar', ({ codigo }) => {
    const moeda = moedaDoCodigo(codigo);
    const maximo = BigInt(casos.limites.preco_unitario_max_unidades);
    const preco: PrecoUnitario = { especie: 'preco_unitario', moeda, unidades: maximo };
    const texto = textoDoPrecoUnitario(preco);
    expect(lerPrecoUnitario(texto, moeda).unidades).toBe(maximo);
    expect(codigoDoErro(() => lerPrecoUnitario(`1${texto}`, moeda))).toBe('formato_invalido');
    expect(codigoDoErro(() => multiplicar(preco, { numerador: BigInt(10001), denominador: BigInt(1) }))).toBe('fora_da_faixa');
  });
});

describe('dinheiro: restrições de tipo', () => {
  it('preço unitário não entra em somar (erro de compilação)', () => {
    const preco = lerPrecoUnitario('1.5', moedaDoCodigo('BRL'));
    // @ts-expect-error preço unitário nunca se soma como valor
    expect(() => somar(preco, preco)).not.toThrow();
  });
});
