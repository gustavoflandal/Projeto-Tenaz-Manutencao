// Textos da Tenaz, compartilhados pela web e pelo app.
//
// O idioma de origem é pt-BR; es-419 (espanhol latino-americano neutro) é o
// segundo. O espaço de nomes único `tenaz` mantém as chaves estáveis.

import { es419 } from './es-419';
import { ptBR } from './pt-BR';

export { es419 } from './es-419';
export { ptBR, type Textos } from './pt-BR';

export const IDIOMA_PADRAO = 'pt-BR';
export const IDIOMAS = ['pt-BR', 'es-419'] as const;
export type Idioma = (typeof IDIOMAS)[number];
export const ESPACO_DE_NOMES = 'tenaz';

/**
 * Forma de um catálogo: as mesmas chaves do pt-BR com texto. O espanhol
 * acrescenta `_many` às formas de plural; fora isso, chave faltando é erro de tipo.
 */
export type Forma<T> = {
  [K in keyof T]: T[K] extends string ? string : Forma<T[K]>;
} & { [chave: `${string}_many`]: string };

export const recursos = {
  'pt-BR': { [ESPACO_DE_NOMES]: ptBR },
  'es-419': { [ESPACO_DE_NOMES]: es419 },
} as const;

export function ehIdioma(valor: unknown): valor is Idioma {
  return typeof valor === 'string' && (IDIOMAS as readonly string[]).includes(valor);
}
