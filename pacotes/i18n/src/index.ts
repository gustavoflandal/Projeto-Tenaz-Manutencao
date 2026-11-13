// Textos da Tenaz, compartilhados pela web e pelo app.
//
// O idioma de origem é pt-BR. O espaço de nomes único `tenaz` mantém as
// chaves estáveis quando outros idiomas entrarem.

import { ptBR } from './pt-BR';

export { ptBR, type Textos } from './pt-BR';

export const IDIOMA_PADRAO = 'pt-BR';
export const ESPACO_DE_NOMES = 'tenaz';

export const recursos = {
  'pt-BR': { [ESPACO_DE_NOMES]: ptBR },
} as const;
