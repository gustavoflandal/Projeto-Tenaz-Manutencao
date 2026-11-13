import { describe, expect, it } from 'vitest';

import { custoHoraOpcional, precoUnitarioParaOContrato } from './formularios';

describe('custo por hora (preço unitário, ADR-013 Parte F)', () => {
  it('aceita vírgula decimal em BRL e converte para o texto com ponto', () => {
    expect(precoUnitarioParaOContrato('45,123456', 'BRL')).toBe('45.123456');
    expect(precoUnitarioParaOContrato(' 125,5 ', 'BRL')).toBe('125.5');
    expect(precoUnitarioParaOContrato('45.50', 'BRL')).toBe('45.50');
    expect(precoUnitarioParaOContrato('0', 'BRL')).toBe('0');
  });

  it('recusa sinal, sétima casa, milhar e zero à esquerda', () => {
    for (const texto of ['-1', '-0', '+1', '45,1234567', '1.234,50', '01', '1,', ',5', 'abc', '1e3', '1000000000000']) {
      expect(precoUnitarioParaOContrato(texto, 'BRL'), texto).toBeNull();
    }
  });

  it('em PYG aceita casas no preço unitário', () => {
    expect(precoUnitarioParaOContrato('75000.5', 'PYG')).toBe('75000.5');
    expect(precoUnitarioParaOContrato('75000,5', 'PYG')).toBe('75000.5');
    expect(precoUnitarioParaOContrato('75000.0000001', 'PYG')).toBeNull();
  });

  it('o esquema aceita vazio e recusa com a chave do catálogo', () => {
    const esquema = custoHoraOpcional('BRL');
    expect(esquema.safeParse('').success).toBe(true);
    expect(esquema.safeParse('45,123456').success).toBe(true);
    const recusado = esquema.safeParse('-5');
    expect(recusado.success).toBe(false);
    expect(recusado.error?.issues[0]?.message).toBe('validacao.custoHora');
  });
});
