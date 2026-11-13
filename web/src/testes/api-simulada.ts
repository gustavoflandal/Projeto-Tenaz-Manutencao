// API simulada para os testes de componente: responde às rotas do contrato
// com dados de teste e registra cada chamada (método, caminho, corpo e
// cabeçalhos), para os testes conferirem If-Match, X-CSRF-Token e o corpo.

import { vi } from 'vitest';

import type { Problema } from '../lib/api';

export interface Chamada {
  readonly metodo: string;
  readonly caminho: string;
  readonly consulta: URLSearchParams;
  readonly corpo: unknown;
  readonly cabecalhos: Headers;
}

export interface RespostaSimulada {
  readonly status?: number;
  readonly corpo?: unknown;
  /** Corpo de erro (application/problem+json). */
  readonly problema?: Partial<Problema> & { readonly codigo: string };
  readonly cabecalhos?: Record<string, string>;
}

export type Manipulador = RespostaSimulada | ((chamada: Chamada) => RespostaSimulada | Promise<RespostaSimulada>);

export interface ApiSimulada {
  readonly chamadas: Chamada[];
  /** Chamadas a uma rota que o teste não simulou (não deveriam existir). */
  readonly semRota: string[];
  chamadasA: (metodoECaminho: string) => Chamada[];
  /** Troca ou acrescenta o manipulador de uma rota no meio do teste. */
  definir: (metodoECaminho: string, manipulador: Manipulador) => void;
}

/** Cria um problema no formato RFC 9457 do contrato. */
export function problema(status: number, codigo: string, extras: Partial<Problema> = {}): RespostaSimulada {
  return {
    status,
    problema: {
      type: `https://api.tenaz.com.br/erros/${codigo}`,
      title: codigo,
      status,
      codigo,
      ...extras,
    },
  };
}

export function simularAPI(rotas: Record<string, Manipulador>): ApiSimulada {
  const manipuladores = new Map<string, Manipulador>(Object.entries(rotas));
  const chamadas: Chamada[] = [];
  const semRota: string[] = [];

  vi.spyOn(globalThis, 'fetch').mockImplementation(async (entrada) => {
    const requisicao = entrada instanceof Request ? entrada : new Request(String(entrada));
    const url = new URL(requisicao.url);
    const texto = await requisicao.clone().text();
    const chamada: Chamada = {
      metodo: requisicao.method,
      caminho: url.pathname,
      consulta: url.searchParams,
      corpo: texto ? (JSON.parse(texto) as unknown) : undefined,
      cabecalhos: requisicao.headers,
    };
    chamadas.push(chamada);

    const manipulador = manipuladores.get(`${chamada.metodo} ${chamada.caminho}`);
    if (manipulador === undefined) {
      semRota.push(`${chamada.metodo} ${chamada.caminho}`);
      return Response.json({ type: 'x', title: 'sem rota', status: 599, codigo: 'sem_rota' }, { status: 599 });
    }
    const resposta = typeof manipulador === 'function' ? await manipulador(chamada) : manipulador;
    const status = resposta.status ?? 200;
    if (resposta.problema) {
      return new Response(JSON.stringify(resposta.problema), {
        status,
        headers: { 'Content-Type': 'application/problem+json', ...resposta.cabecalhos },
      });
    }
    if (status === 204 || status === 202 || resposta.corpo === undefined) {
      return new Response(null, { status, headers: resposta.cabecalhos });
    }
    return Response.json(resposta.corpo, { status, headers: resposta.cabecalhos });
  });

  return {
    chamadas,
    semRota,
    chamadasA: (chave) => chamadas.filter((c) => `${c.metodo} ${c.caminho}` === chave),
    definir: (chave, manipulador) => manipuladores.set(chave, manipulador),
  };
}
