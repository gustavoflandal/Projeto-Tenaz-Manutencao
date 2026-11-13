// Unidades, turnos e feriados (RF-ADM-002, RN-ADM-019).

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { components } from '@tenaz/api-tipos';

import { api, executar, obter, seVersao } from '../../lib/api';
import { chaves } from '../../lib/consultas';

type Esquemas = components['schemas'];

export function useUnidades() {
  return useQuery({ queryKey: chaves.unidades, queryFn: () => obter(api.GET('/v1/unidades')) });
}

export function useUnidade(id: string) {
  return useQuery({
    queryKey: chaves.unidade(id),
    queryFn: () => obter(api.GET('/v1/unidades/{id}', { params: { path: { id } } })),
  });
}

function useInvalidarUnidade() {
  const cliente = useQueryClient();
  return (id: string) => {
    void cliente.invalidateQueries({ queryKey: chaves.unidade(id) });
    void cliente.invalidateQueries({ queryKey: chaves.unidades, exact: true });
  };
}

export function useCriarUnidade() {
  const cliente = useQueryClient();
  return useMutation({
    mutationFn: (corpo: Esquemas['EntradaUnidade']) => obter(api.POST('/v1/unidades', { body: corpo })),
    onSuccess: () => void cliente.invalidateQueries({ queryKey: chaves.unidades }),
  });
}

export function useAlterarUnidade() {
  const invalidar = useInvalidarUnidade();
  return useMutation({
    mutationFn: ({ id, versao, corpo }: { id: string; versao: number; corpo: Esquemas['PatchUnidade'] }) =>
      obter(
        api.PATCH('/v1/unidades/{id}', {
          params: { path: { id }, header: { 'If-Match': seVersao(versao) } },
          body: corpo,
        }),
      ),
    onSuccess: (_u, { id }) => invalidar(id),
  });
}

export function useInativarUnidade() {
  const invalidar = useInvalidarUnidade();
  return useMutation({
    mutationFn: ({ id, versao }: { id: string; versao: number }) =>
      executar(api.DELETE('/v1/unidades/{id}', { params: { path: { id }, header: { 'If-Match': seVersao(versao) } } })),
    onSuccess: (_v, { id }) => invalidar(id),
  });
}

export function useReativarUnidade() {
  const invalidar = useInvalidarUnidade();
  return useMutation({
    mutationFn: ({ id, versao }: { id: string; versao: number }) =>
      obter(
        api.POST('/v1/unidades/{id}/reativar', { params: { path: { id }, header: { 'If-Match': seVersao(versao) } } }),
      ),
    onSuccess: (_u, { id }) => invalidar(id),
  });
}

// ------------------------------------------------------------------- turnos

export function useIncluirTurno() {
  const invalidar = useInvalidarUnidade();
  return useMutation({
    mutationFn: ({ unidadeId, corpo }: { unidadeId: string; corpo: Esquemas['EntradaTurno'] }) =>
      obter(api.POST('/v1/unidades/{id}/turnos', { params: { path: { id: unidadeId } }, body: corpo })),
    onSuccess: (_t, { unidadeId }) => invalidar(unidadeId),
  });
}

export function useAlterarTurno() {
  const invalidar = useInvalidarUnidade();
  return useMutation({
    mutationFn: ({
      unidadeId,
      turnoId,
      versao,
      corpo,
    }: {
      unidadeId: string;
      turnoId: string;
      versao: number;
      corpo: Esquemas['PatchTurno'];
    }) =>
      obter(
        api.PATCH('/v1/unidades/{id}/turnos/{turno_id}', {
          params: { path: { id: unidadeId, turno_id: turnoId }, header: { 'If-Match': seVersao(versao) } },
          body: corpo,
        }),
      ),
    onSuccess: (_t, { unidadeId }) => invalidar(unidadeId),
  });
}

export function useRemoverTurno() {
  const invalidar = useInvalidarUnidade();
  return useMutation({
    mutationFn: ({ unidadeId, turnoId, versao }: { unidadeId: string; turnoId: string; versao: number }) =>
      executar(
        api.DELETE('/v1/unidades/{id}/turnos/{turno_id}', {
          params: { path: { id: unidadeId, turno_id: turnoId }, header: { 'If-Match': seVersao(versao) } },
        }),
      ),
    onSuccess: (_v, { unidadeId }) => invalidar(unidadeId),
  });
}

// ----------------------------------------------------------------- feriados

export function useFeriados(unidadeId: string) {
  return useQuery({
    queryKey: chaves.feriados(unidadeId),
    queryFn: () => obter(api.GET('/v1/unidades/{id}/feriados', { params: { path: { id: unidadeId } } })),
  });
}

export function useIncluirFeriado() {
  const cliente = useQueryClient();
  return useMutation({
    mutationFn: ({ unidadeId, corpo }: { unidadeId: string; corpo: Esquemas['EntradaFeriado'] }) =>
      obter(api.POST('/v1/unidades/{id}/feriados', { params: { path: { id: unidadeId } }, body: corpo })),
    onSuccess: (_f, { unidadeId }) => void cliente.invalidateQueries({ queryKey: chaves.feriados(unidadeId) }),
  });
}

export function useRemoverFeriado() {
  const cliente = useQueryClient();
  return useMutation({
    mutationFn: ({ unidadeId, feriadoId, versao }: { unidadeId: string; feriadoId: string; versao: number }) =>
      executar(
        api.DELETE('/v1/unidades/{id}/feriados/{feriado_id}', {
          params: { path: { id: unidadeId, feriado_id: feriadoId }, header: { 'If-Match': seVersao(versao) } },
        }),
      ),
    onSuccess: (_v, { unidadeId }) => void cliente.invalidateQueries({ queryKey: chaves.feriados(unidadeId) }),
  });
}
