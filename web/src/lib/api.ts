// Cliente da API da Tenaz, gerado do contrato OpenAPI.
//
// Toda chamada ao servidor passa por aqui, com os tipos de @tenaz/api-tipos.
// Os componentes não chamam `fetch` direto: usam hooks com TanStack Query.

import type { paths } from '@tenaz/api-tipos';
import createClient from 'openapi-fetch';

export const api = createClient<paths>({
  // A web é servida no mesmo domínio da API (em desenvolvimento, pelo proxy
  // do Vite). O endereço absoluto evita depender de URL relativa no Request.
  baseUrl: globalThis.location.origin,
  // Resolvido a cada chamada, e não na criação do cliente, para que os testes
  // possam substituir o fetch global.
  fetch: (requisicao) => globalThis.fetch(requisicao),
});

/** Falha de uma chamada à API, com o status HTTP quando houve resposta. */
export class ErroDaAPI extends Error {
  readonly status: number | undefined;

  constructor(mensagem: string, status?: number) {
    super(mensagem);
    this.name = 'ErroDaAPI';
    this.status = status;
  }
}
