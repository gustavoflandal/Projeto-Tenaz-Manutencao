// Tipos do contrato da API da Tenaz. Gerados; não escreva tipo de API à mão.

export type { components, operations, paths } from './esquema';

import type { paths } from './esquema';

/** Corpo da resposta 200 de uma operação GET. */
export type RespostaGet<Rota extends keyof paths> = paths[Rota] extends {
  get: { responses: { 200: { content: { 'application/json': infer Corpo } } } };
}
  ? Corpo
  : never;

/** Resposta de GET /v1/versao (RNF-039). */
export type VersaoDaAPI = RespostaGet<'/v1/versao'>;
