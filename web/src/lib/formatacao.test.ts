import { afterEach, describe, expect, it } from 'vitest';

import { dataParaISO, formatarData, formatarDataHora, formatarPrecoUnitario, listar, mascararData, simboloDaMoeda } from './formatacao';
import { i18n } from './i18n';

describe('formatarDataHora', () => {
  it('usa dd/mm/aaaa hh:mm em 24 h, no fuso dado', () => {
    expect(formatarDataHora('2026-11-12T14:05:00Z', 'UTC', 'BR')).toBe('12/11/2026 14:05');
    expect(formatarDataHora('2026-11-12T03:07:00Z', 'America/Sao_Paulo', 'BR')).toBe('12/11/2026 00:07');
  });

  it('devolve null para valor vazio ou inválido', () => {
    expect(formatarDataHora('', 'UTC', 'BR')).toBeNull();
    expect(formatarDataHora('ontem', 'UTC', 'BR')).toBeNull();
  });
});

describe('datas civis', () => {
  it('formata e converte sem mudar o dia por causa de fuso', () => {
    expect(formatarData('2026-11-13', 'BR')).toBe('13/11/2026');
    expect(dataParaISO('13/11/2026')).toBe('2026-11-13');
  });

  it('texto que não é data volta como veio', () => {
    expect(formatarData('2026-02-31', 'BR')).toBe('2026-02-31');
  });

  it('recusa data inexistente ou fora do formato', () => {
    expect(dataParaISO('31/02/2026')).toBeNull();
    expect(dataParaISO('2026-11-13')).toBeNull();
    expect(dataParaISO('1/1/2026')).toBeNull();
  });

  it('aplica a máscara enquanto digita', () => {
    expect(mascararData('1311')).toBe('13/11');
    expect(mascararData('13112026')).toBe('13/11/2026');
    expect(mascararData('13/11/20269')).toBe('13/11/2026');
  });
});

describe('dinheiro pelo país da empresa, nunca pelo idioma da interface', () => {
  afterEach(async () => {
    await i18n.changeLanguage('pt-BR');
  });

  it('símbolo de cada país', () => {
    expect(simboloDaMoeda('BR')).toBe('R$');
    expect(simboloDaMoeda('AR')).toBe('$');
    expect(simboloDaMoeda('PY')).toBe('Gs.');
  });

  it('preço unitário do texto da API mostra de duas a seis casas, sem passar por number', () => {
    expect(formatarPrecoUnitario('45.50', 'BR')).toBe('R$ 45,50');
    expect(formatarPrecoUnitario('45.123456', 'BR')).toBe('R$ 45,123456');
    expect(formatarPrecoUnitario('75000.5', 'PY')).toBe('Gs. 75.000,5');
  });

  it('usuário em es-419 numa empresa brasileira vê vírgula decimal e R$', async () => {
    await i18n.changeLanguage('es-419');
    expect(formatarPrecoUnitario('1234.5', 'BR')).toBe('R$ 1.234,50');
    expect(formatarData('2026-11-13', 'BR')).toBe('13/11/2026');
  });
});

describe('listar', () => {
  it('junta com vírgula e ignora vazios', () => {
    expect(listar(['Contagem', null, 'Minas Gerais', ''])).toBe('Contagem, Minas Gerais');
  });
});
