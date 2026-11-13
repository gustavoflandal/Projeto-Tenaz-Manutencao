// Pessoas (RF-ADM-015): quem executa, conduz ou responde por algo, com ou sem
// usuário. Lista paginada por cursor, com CPF mascarado; o detalhe traz o CPF
// completo (e a API grava a visualização na trilha).

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { components } from '@tenaz/api-tipos';

import { api, executar, obter, seVersao } from '../../lib/api';
import { chaves } from '../../lib/consultas';
import { LIMITE_DA_PAGINA, useListaPorCursor } from '../../lib/lista';

type Esquemas = components['schemas'];

export interface FiltroDePessoas {
  readonly busca?: string;
  readonly situacao?: 'ativa' | 'inativa';
  readonly com_usuario?: 'true' | 'false';
  readonly unidade_id?: string;
}

export function usePessoas(filtro: FiltroDePessoas) {
  return useListaPorCursor(chaves.pessoas(filtro), (cursor) =>
    obter(
      api.GET('/v1/pessoas', {
        params: { query: { ...filtro, limite: LIMITE_DA_PAGINA, ...(cursor === undefined ? {} : { cursor }) } },
      }),
    ),
  );
}

export function usePessoa(id: string) {
  return useQuery({
    queryKey: chaves.pessoa(id),
    queryFn: () => obter(api.GET('/v1/pessoas/{id}', { params: { path: { id } } })),
  });
}

export function useCriarPessoa() {
  const cliente = useQueryClient();
  return useMutation({
    mutationFn: (corpo: Esquemas['EntradaPessoa']) => obter(api.POST('/v1/pessoas', { body: corpo })),
    onSuccess: () => void cliente.invalidateQueries({ queryKey: chaves.pessoas() }),
  });
}

export function useAlterarPessoa() {
  const cliente = useQueryClient();
  return useMutation({
    mutationFn: ({ id, versao, corpo }: { id: string; versao: number; corpo: Esquemas['PatchPessoa'] }) =>
      obter(
        api.PATCH('/v1/pessoas/{id}', { params: { path: { id }, header: { 'If-Match': seVersao(versao) } }, body: corpo }),
      ),
    onSuccess: () => void cliente.invalidateQueries({ queryKey: chaves.pessoas() }),
  });
}

export function useInativarPessoa() {
  const cliente = useQueryClient();
  return useMutation({
    mutationFn: ({ id, versao }: { id: string; versao: number }) =>
      executar(api.DELETE('/v1/pessoas/{id}', { params: { path: { id }, header: { 'If-Match': seVersao(versao) } } })),
    onSuccess: () => void cliente.invalidateQueries({ queryKey: chaves.pessoas() }),
  });
}

export function useReativarPessoa() {
  const cliente = useQueryClient();
  return useMutation({
    mutationFn: ({ id, versao }: { id: string; versao: number }) =>
      obter(
        api.POST('/v1/pessoas/{id}/reativar', { params: { path: { id }, header: { 'If-Match': seVersao(versao) } } }),
      ),
    onSuccess: () => void cliente.invalidateQueries({ queryKey: chaves.pessoas() }),
  });
}
