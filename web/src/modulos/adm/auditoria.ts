// Trilha de auditoria (RF-ADM-011): consulta, detalhe, exportação em CSV,
// verificação de integridade e configuração das capturas por tela. Só a
// configuração escreve (`adm.gerenciar_parametros`); nenhum perfil altera evento.

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { components } from '@tenaz/api-tipos';

import { api, ErroDaAPI, obter } from '../../lib/api';
import { chaves } from '../../lib/consultas';
import { LIMITE_DA_PAGINA, useListaPorCursor } from '../../lib/lista';

export interface FiltroDeEventos {
  readonly entidade?: string;
  readonly entidade_id?: string;
  readonly usuario_id?: string;
  readonly acao?: string;
  readonly modulo?: string;
  readonly tela?: string;
  /** Endereço exato; a API exige `de` junto. */
  readonly ip?: string;
  /** A partir deste instante, inclusive (RFC 3339 com fuso). */
  readonly de?: string;
  /** Antes deste instante (RFC 3339 com fuso). */
  readonly ate?: string;
}


function consultaDoFiltro(filtro: FiltroDeEventos) {
  const { de, ate, acao, ...resto } = filtro;
  return {
    ...resto,
    // `acao` é enum no contrato; o valor vem do seletor de tipos fechados.
    ...(acao === undefined ? {} : { acao: acao as components['schemas']['Evento']['acao'] }),
    ...(de === undefined ? {} : { 'ocorrido_em[gte]': de }),
    ...(ate === undefined ? {} : { 'ocorrido_em[lt]': ate }),
  };
}

export function useEventos(filtro: FiltroDeEventos) {
  return useListaPorCursor(chaves.eventos(filtro), (cursor) =>
    obter(
      api.GET('/v1/eventos-auditoria', {
        params: {
          query: { ...consultaDoFiltro(filtro), limite: LIMITE_DA_PAGINA, ...(cursor === undefined ? {} : { cursor }) },
        },
      }),
    ),
  );
}

export function useEvento(id: string | null) {
  return useQuery({
    queryKey: chaves.evento(id ?? ''),
    queryFn: () => obter(api.GET('/v1/eventos-auditoria/{id}', { params: { path: { id: id ?? '' } } })),
    enabled: id !== null,
  });
}

/** Baixa o CSV da consulta: período obrigatório, no máximo 50.000 eventos. */
export async function exportarEventos(filtro: FiltroDeEventos): Promise<Blob> {
  const { data, error, response } = await api.GET('/v1/exportacoes-eventos-auditoria', {
    params: { query: consultaDoFiltro(filtro) as never },
    parseAs: 'blob',
  });
  if (!response.ok || data === undefined) {
    const problema = typeof error === 'object' && error !== null && 'codigo' in error ? (error as never) : undefined;
    throw new ErroDaAPI(`A exportação devolveu ${response.status}`, response.status, problema);
  }
  return data;
}

// ------------------------------------------------------------ integridade

export function useVerificacao(de: string | null, ate: string | null) {
  return useQuery({
    queryKey: chaves.verificacao(de ?? '', ate ?? ''),
    queryFn: () =>
      obter(api.GET('/v1/verificacoes-auditoria', { params: { query: { de: de ?? '', ate: ate ?? '' } } })),
    enabled: de !== null && ate !== null,
    retry: false,
    // Cada verificação recalcula no servidor: nada de reaproveitar resultado antigo.
    gcTime: 0,
    staleTime: 0,
  });
}

// ---------------------------------------------------------------- capturas

export function useCapturas() {
  return useQuery({ queryKey: chaves.capturas, queryFn: () => obter(api.GET('/v1/capturas-auditoria')) });
}

export function useDefinirCaptura() {
  const cliente = useQueryClient();
  return useMutation({
    mutationFn: ({ tela, chaves: valores }: { tela: string; chaves: components['schemas']['EntradaCapturaAuditoria'] }) =>
      obter(api.PUT('/v1/capturas-auditoria/{tela}', { params: { path: { tela } }, body: valores })),
    onSuccess: () => void cliente.invalidateQueries({ queryKey: chaves.capturas }),
  });
}
