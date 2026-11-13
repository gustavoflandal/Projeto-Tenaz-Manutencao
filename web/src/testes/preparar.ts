// Preparação comum dos testes de componente (Vitest + Testing Library, jsdom).

import '@testing-library/jest-dom/vitest';

import { cleanup } from '@testing-library/react';
import { afterEach, vi } from 'vitest';

import { i18n } from '../lib/i18n';
import { iniciarTema } from '../lib/tema';

// O jsdom não implementa estas APIs, usadas pelos primitivos do Radix.
class ObservadorDeTamanho {
  observe(): void {}
  unobserve(): void {}
  disconnect(): void {}
}
globalThis.ResizeObserver ??= ObservadorDeTamanho;
Element.prototype.scrollIntoView ??= function scrollIntoView() {};
Element.prototype.hasPointerCapture ??= () => false;
Element.prototype.releasePointerCapture ??= () => {};

// O jsdom não carrega o index.html nem o principal.tsx: idioma e título são
// os mesmos que eles declaram.
document.documentElement.lang = 'pt-BR';
document.title = i18n.t('app.titulo');

afterEach(() => {
  cleanup();
  vi.restoreAllMocks();
  globalThis.localStorage.clear();
  // Volta o tema a "igual ao sistema", como num navegador sem escolha salva.
  iniciarTema();
});
