import { describe, expect, it } from 'vitest';

import casos from '../../../backend/testes/casos/formatos.json';
import {
  ErroDeFormato,
  formatarData,
  formatarDataHora,
  formatarDecimal,
  formatarInteiro,
  formatarMedidor,
  formatarPercentual,
  formatarPrecoUnitario,
  formatarValor,
  simboloDoPais,
} from './formatos';
import { ehPais, type Moeda, type Pais } from './paises';

function pais(texto: string): Pais {
  if (!ehPais(texto)) throw new Error(`país fora do catálogo nos casos: ${texto}`);
  return texto;
}

function moeda(texto: string): Moeda {
  if (texto !== 'BRL' && texto !== 'ARS' && texto !== 'PYG') throw new Error(`moeda fora do catálogo nos casos: ${texto}`);
  return texto;
}

describe('formatos (casos compartilhados, versão 1)', () => {
  it('trava a versão, o tamanho de cada grupo e o total', () => {
    expect(casos.versao).toBe(1);
    const tamanhos = {
      inteiro: casos.inteiro.length,
      decimal: casos.decimal.length,
      percentual: casos.percentual.length,
      medidor: casos.medidor.length,
      valor: casos.valor.length,
      preco_unitario: casos.preco_unitario.length,
      data: casos.data.length,
      data_hora: casos.data_hora.length,
    };
    expect(tamanhos).toEqual({ inteiro: 14, decimal: 13, percentual: 9, medidor: 11, valor: 21, preco_unitario: 16, data: 9, data_hora: 22 });
    expect(Object.values(tamanhos).reduce((a, b) => a + b, 0)).toBe(115);
  });

  it('o símbolo de cada país é o dos casos', () => {
    for (const [codigo, simbolo] of Object.entries(casos.simbolos)) expect(simboloDoPais(pais(codigo))).toBe(simbolo);
  });

  it.each(casos.inteiro)('inteiro: $nome', (c) => {
    expect(formatarInteiro(c.entrada, pais(c.pais))).toBe(c.esperado);
  });
  it.each(casos.decimal)('decimal: $nome', (c) => {
    expect(formatarDecimal(c.entrada, pais(c.pais))).toBe(c.esperado);
  });
  it.each(casos.percentual)('percentual: $nome', (c) => {
    expect(formatarPercentual(c.entrada, pais(c.pais))).toBe(c.esperado);
  });
  it.each(casos.medidor)('medidor: $nome', (c) => {
    expect(formatarMedidor(c.entrada, c.unidade, pais(c.pais))).toBe(c.esperado);
  });
  it.each(casos.valor)('valor: $nome', (c) => {
    expect(formatarValor(c.entrada, moeda(c.moeda), pais(c.pais))).toBe(c.esperado);
  });
  it.each(casos.preco_unitario)('preço unitário: $nome', (c) => {
    expect(formatarPrecoUnitario(c.entrada, moeda(c.moeda), pais(c.pais))).toBe(c.esperado);
  });
  it.each(casos.data)('data: $nome', (c) => {
    expect(formatarData(c.entrada, pais(c.pais))).toBe(c.esperado);
  });
  it.each(casos.data_hora)('data e hora: $nome', (c) => {
    expect(formatarDataHora(c.instante, c.fuso, pais(c.pais))).toBe(c.esperado);
  });
});

describe('formatos: o que não é texto canônico lança erro', () => {
  it('recusa entrada fora do texto canônico', () => {
    for (const ruim of ['', ' 1', '1,5', '1.', '.5', '01', '1e3', '+1', '1_000', 'abc']) {
      expect(() => formatarDecimal(ruim, 'BR'), ruim).toThrow(ErroDeFormato);
    }
    for (const ruim of ['1.5', '', '01', '1 000']) expect(() => formatarInteiro(ruim, 'BR'), ruim).toThrow(ErroDeFormato);
    expect(() => formatarInteiro(1.5, 'BR')).toThrow(ErroDeFormato);
    expect(() => formatarInteiro(Number.NaN, 'BR')).toThrow(ErroDeFormato);
    expect(() => formatarInteiro(2 ** 60, 'BR')).toThrow(ErroDeFormato);
  });

  it('recusa -0 em todo formato numérico', () => {
    expect(() => formatarInteiro('-0', 'BR')).toThrow(ErroDeFormato);
    expect(() => formatarInteiro(-0, 'BR')).toThrow(ErroDeFormato);
    expect(() => formatarDecimal('-0', 'BR')).toThrow(ErroDeFormato);
    expect(() => formatarDecimal('-0.00', 'BR')).toThrow(ErroDeFormato);
    expect(() => formatarPercentual('-0', 'BR')).toThrow(ErroDeFormato);
    expect(() => formatarValor('-0.00', 'BRL', 'BR')).toThrow(ErroDeFormato);
    expect(() => formatarPrecoUnitario('-0.00', 'BRL', 'BR')).toThrow(ErroDeFormato);
  });

  it('recusa dinheiro fora do texto canônico da Parte F', () => {
    expect(() => formatarValor('10,00', 'BRL', 'BR')).toThrow();
    expect(() => formatarValor('10.005', 'BRL', 'BR')).toThrow();
    expect(() => formatarPrecoUnitario('1.1234567', 'BRL', 'BR')).toThrow();
  });

  it('recusa medidor sem unidade', () => {
    expect(() => formatarMedidor('10', ' ', 'BR')).toThrow(ErroDeFormato);
  });

  it('recusa data inexistente ou fora do formato', () => {
    for (const ruim of ['2026-02-29', '2026-13-01', '2026-00-10', '2026-04-31', '0000-01-01', '2026-1-1', '12/11/2026', '']) {
      expect(() => formatarData(ruim, 'BR'), ruim).toThrow(ErroDeFormato);
    }
    expect(formatarData('2028-02-29', 'BR')).toBe('29/02/2028');
  });

  it('recusa instante inválido e fuso desconhecido', () => {
    for (const ruim of ['', 'ontem', '2026-07-15', '2026-07-15T17:05', '2026-02-30T10:00:00Z', '2026-07-15T24:00:00Z', '2026-07-15T17:60:00Z']) {
      expect(() => formatarDataHora(ruim, 'America/Sao_Paulo', 'BR'), ruim).toThrow(ErroDeFormato);
    }
    expect(() => formatarDataHora('2026-07-15T17:05:00Z', 'Marte/Olimpo', 'BR')).toThrow(ErroDeFormato);
    expect(() => formatarDataHora('2026-07-15T17:05:00Z', '', 'BR')).toThrow(ErroDeFormato);
  });

  it('corta os segundos e aceita fração e deslocamento', () => {
    expect(formatarDataHora('2026-07-15T17:05:59Z', 'America/Sao_Paulo', 'BR')).toBe('15/07/2026 14:05');
    expect(formatarDataHora('2026-07-15T17:05:59.999Z', 'America/Sao_Paulo', 'BR')).toBe('15/07/2026 14:05');
    expect(formatarDataHora('2026-07-15T14:05:00-03:00', 'UTC', 'BR')).toBe('15/07/2026 17:05');
  });
});

describe('formatos: o idioma do usuário não entra', () => {
  it('as funções não recebem idioma: o formato é do país', () => {
    expect(formatarValor('1234.5', 'BRL', 'BR')).toBe('R$ 1.234,50');
    expect(formatarDecimal('1234.5', 'BR')).toBe('1.234,5');
  });
});
