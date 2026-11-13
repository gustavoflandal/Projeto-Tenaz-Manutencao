// Escolha do idioma da interface (pt-BR ou es-419).
//
// Antes do login: a escolha da pessoa, guardada no navegador; depois o idioma
// do navegador (qualquer variante de espanhol vira es-419, de português vira
// pt-BR); por último o padrão, pt-BR. Depois do login vale `usuario.idioma`
// (modulos/adm/idioma.ts), e a escolha da pessoa é gravada nele.

import type { Pais } from '@tenaz/dominio';
import { ehIdioma, IDIOMA_PADRAO, type Idioma } from '@tenaz/i18n';

const CHAVE = 'tenaz.idioma';

export function lerIdiomaSalvo(): Idioma | null {
  try {
    const salvo = globalThis.localStorage.getItem(CHAVE);
    return ehIdioma(salvo) ? salvo : null;
  } catch {
    return null;
  }
}

export function guardarIdioma(idioma: Idioma): void {
  try {
    globalThis.localStorage.setItem(CHAVE, idioma);
  } catch {
    // Sem armazenamento, a escolha vale até recarregar.
  }
}

/** Idioma de uma lista de preferências do navegador (`navigator.languages`). */
export function idiomaDoNavegador(preferidos: readonly string[]): Idioma | null {
  for (const preferido of preferidos) {
    const base = preferido.toLowerCase().split('-')[0];
    if (base === 'es') return 'es-419';
    if (base === 'pt') return 'pt-BR';
  }
  return null;
}

export function idiomaInicial(): Idioma {
  const navegador = typeof globalThis.navigator === 'undefined' ? [] : (globalThis.navigator.languages ?? []);
  return lerIdiomaSalvo() ?? idiomaDoNavegador(navegador) ?? IDIOMA_PADRAO;
}

/** País sugerido pela região do navegador (es-AR, es-PY, pt-BR); sem região reconhecida, nenhum. */
export function paisDoNavegador(preferidos: readonly string[]): Pais | null {
  for (const preferido of preferidos) {
    const regiao = preferido.toUpperCase().split('-')[1];
    if (regiao === 'BR' || regiao === 'AR' || regiao === 'PY') return regiao;
  }
  return null;
}
