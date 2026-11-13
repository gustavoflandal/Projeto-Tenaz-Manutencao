// Breakpoints da seção 3.10 do documento 7. CSS não aceita variável em
// @media; os valores vêm do mesmo tokens.json que gera o tema do Tailwind.

import tokens from '@tenaz/tokens-ui/tokens.json';
import { useSyncExternalStore } from 'react';

export type Breakpoint = keyof typeof tokens.breakpoints;

function consulta(bp: Breakpoint): MediaQueryList | undefined {
  return typeof globalThis.matchMedia === 'function'
    ? globalThis.matchMedia(`(min-width: ${tokens.breakpoints[bp]})`)
    : undefined;
}

/**
 * A janela tem pelo menos a largura do breakpoint. Sem `matchMedia` (testes
 * em jsdom), considera a tela de mesa.
 */
export function useAcimaDe(bp: Breakpoint): boolean {
  return useSyncExternalStore(
    (avisar) => {
      const lista = consulta(bp);
      lista?.addEventListener('change', avisar);
      return () => lista?.removeEventListener('change', avisar);
    },
    () => consulta(bp)?.matches ?? true,
    () => true,
  );
}
