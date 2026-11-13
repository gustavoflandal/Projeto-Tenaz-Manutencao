import { describe, expect, it } from 'vitest';

import { ptBR } from './pt-BR';

function folhas(obj: object, prefixo = ''): Array<[string, string]> {
  return Object.entries(obj).flatMap(([chave, valor]) =>
    typeof valor === 'string'
      ? [[`${prefixo}${chave}`, valor] as [string, string]]
      : folhas(valor as object, `${prefixo}${chave}.`),
  );
}

const textos = folhas(ptBR);

// Regras de tom de voz verificáveis por máquina (seção 10.1 do documento 7).
describe('tom de voz dos textos', () => {
  it.each(textos)('%s não é vazio', (_chave, texto) => {
    expect(texto.trim()).not.toBe('');
  });

  it.each(textos)('%s não usa exclamação', (_chave, texto) => {
    expect(texto).not.toContain('!');
  });

  it.each(textos)('%s não usa termos proibidos', (_chave, texto) => {
    expect(texto.toLowerCase()).not.toMatch(/\b(ticket|chamado|offline|tenant)\b/);
  });
});
