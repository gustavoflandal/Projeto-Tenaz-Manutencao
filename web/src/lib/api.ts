// Cliente da API da Tenaz, gerado do contrato OpenAPI.
//
// Toda chamada ao servidor passa por aqui, com os tipos de @tenaz/api-tipos.
// Os componentes não chamam `fetch` direto: usam hooks com TanStack Query.

import type { components, paths } from '@tenaz/api-tipos';
import createClient, { type Middleware } from 'openapi-fetch';

import { i18n } from './i18n';

export type Problema = components['schemas']['Problema'];

export const api = createClient<paths>({
  // A web é servida no mesmo domínio da API (em desenvolvimento, pelo proxy
  // do Vite). O endereço absoluto evita depender de URL relativa no Request.
  baseUrl: globalThis.location.origin,
  // Resolvido a cada chamada, e não na criação do cliente, para que os testes
  // possam substituir o fetch global.
  fetch: (requisicao) => globalThis.fetch(requisicao),
});

/** Falha de uma chamada à API, com o status HTTP e o problema quando houve resposta. */
export class ErroDaAPI extends Error {
  readonly status: number | undefined;
  readonly problema: Problema | undefined;

  constructor(mensagem: string, status?: number, problema?: Problema) {
    super(mensagem);
    this.name = 'ErroDaAPI';
    this.status = status;
    this.problema = problema;
  }

  /** Código estável do erro (seção 3.2 do documento 3), o que o cliente trata. */
  get codigo(): string | undefined {
    return this.problema?.codigo;
  }
}

// -------------------------------------------------------------------- CSRF

type FonteDoCSRF = () => string | null | undefined;

let fonteDoCSRF: FonteDoCSRF = () => null;

/**
 * Registra de onde sai o token CSRF da sessão (seção 4.2 do documento 3). A
 * fonte é o cache de consultas: este módulo não guarda cópia do dado.
 */
export function registrarFonteDoCSRF(fonte: FonteDoCSRF): void {
  fonteDoCSRF = fonte;
}

type FonteDaTela = () => string | null | undefined;

let fonteDaTela: FonteDaTela = () => null;

/**
 * Registra de onde sai o código da tela atual (X-Tenaz-Tela). A fonte é o
 * roteador: este módulo não guarda cópia do endereço.
 */
export function registrarFonteDaTela(fonte: FonteDaTela): void {
  fonteDaTela = fonte;
}

const METODOS_QUE_ALTERAM = new Set(['POST', 'PATCH', 'PUT', 'DELETE']);

const middleware: Middleware = {
  onRequest({ request }) {
    // O servidor negocia o idioma das mensagens e dos nomes pelo idioma em uso na interface.
    if (i18n.language) request.headers.set('Accept-Language', i18n.language);
    // A trilha grava a tela de origem de toda chamada (RF-ADM-011).
    const tela = fonteDaTela();
    if (tela) request.headers.set('X-Tenaz-Tela', tela);
    if (METODOS_QUE_ALTERAM.has(request.method)) {
      const csrf = fonteDoCSRF();
      if (csrf) request.headers.set('X-CSRF-Token', csrf);
    }
    // A web envia sempre Idempotency-Key nos POST (seção 2.9 do documento 3).
    if (request.method === 'POST' && !request.headers.has('Idempotency-Key')) {
      request.headers.set('Idempotency-Key', globalThis.crypto.randomUUID());
    }
    return request;
  },
};
api.use(middleware);

// ----------------------------------------------------------------- resultados

interface Resultado<D> {
  readonly data?: D;
  readonly error?: unknown;
  readonly response: Response;
}

function comoProblema(erro: unknown): Problema | undefined {
  if (typeof erro === 'object' && erro !== null && 'codigo' in erro && typeof erro.codigo === 'string') {
    return erro as Problema;
  }
  return undefined;
}

function falhaDe(resultado: Resultado<unknown>): ErroDaAPI {
  const problema = comoProblema(resultado.error);
  return new ErroDaAPI(
    problema?.detail ?? problema?.title ?? `A API devolveu ${resultado.response.status}`,
    resultado.response.status,
    problema,
  );
}

/** Espera a resposta com corpo; qualquer falha vira ErroDaAPI. */
export async function obter<D>(chamada: Promise<Resultado<D>>): Promise<D> {
  const resultado = await chamada;
  if (!resultado.response.ok || resultado.error !== undefined || resultado.data === undefined) {
    throw falhaDe(resultado);
  }
  return resultado.data;
}

/** Espera a resposta sem corpo (204, 202); qualquer falha vira ErroDaAPI. */
export async function executar(chamada: Promise<Resultado<unknown>>): Promise<void> {
  const resultado = await chamada;
  if (!resultado.response.ok || resultado.error !== undefined) throw falhaDe(resultado);
}

/** If-Match com o ETag lido: a versão do recurso entre aspas (seção 2.11 do documento 3). */
export function seVersao(versao: number): string {
  return `"${versao}"`;
}
