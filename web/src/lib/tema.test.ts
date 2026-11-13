import { describe, expect, it } from 'vitest';

import { aplicarTema, iniciarTema, lerPreferenciaSalva, resolverTema } from './tema';

describe('tema', () => {
  it('resolve "sistema" pela preferência do sistema operacional', () => {
    expect(resolverTema('sistema', true)).toBe('escuro');
    expect(resolverTema('sistema', false)).toBe('claro');
    expect(resolverTema('sol', true)).toBe('sol');
  });

  it('aplica o tema pelo atributo data-tema, e o remove para seguir o sistema', () => {
    const raiz = document.createElement('html');
    aplicarTema('escuro', raiz);
    expect(raiz.getAttribute('data-tema')).toBe('escuro');
    aplicarTema('sistema', raiz);
    expect(raiz.hasAttribute('data-tema')).toBe(false);
  });

  it('ignora valor salvo desconhecido', () => {
    globalThis.localStorage.setItem('tenaz.tema', 'roxo');
    expect(lerPreferenciaSalva()).toBe('sistema');
  });

  it('aplica a preferência salva antes da primeira renderização', () => {
    globalThis.localStorage.setItem('tenaz.tema', 'sol');
    iniciarTema();
    expect(document.documentElement.getAttribute('data-tema')).toBe('sol');
  });
});
