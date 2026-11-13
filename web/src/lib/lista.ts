// Lista paginada por cursor (seção 2.6 do documento 3): cada página vem do
// servidor, e a próxima só é pedida quando a pessoa pede. Nunca se carrega
// tudo para filtrar no navegador (seção 6.4 do documento 7).

import { useInfiniteQuery, type QueryKey } from '@tanstack/react-query';

import type { PaginaDeLista } from './tipos';

export const LIMITE_DA_PAGINA = 50;

export interface ListaPorCursor<T> {
  readonly itens: readonly T[];
  readonly carregando: boolean;
  readonly erro: unknown;
  readonly temMais: boolean;
  readonly carregandoMais: boolean;
  readonly carregarMais: () => void;
  readonly recarregar: () => void;
}

export function useListaPorCursor<T>(
  chave: QueryKey,
  buscar: (cursor: string | undefined) => Promise<PaginaDeLista<T>>,
): ListaPorCursor<T> {
  const consulta = useInfiniteQuery({
    queryKey: chave,
    queryFn: ({ pageParam }) => buscar(pageParam),
    initialPageParam: undefined as string | undefined,
    getNextPageParam: (ultima) => (ultima.paginacao.tem_mais ? (ultima.paginacao.proximo_cursor ?? undefined) : undefined),
  });
  return {
    itens: consulta.data?.pages.flatMap((pagina) => pagina.dados ?? []) ?? [],
    carregando: consulta.isPending,
    erro: consulta.error,
    temMais: consulta.hasNextPage,
    carregandoMais: consulta.isFetchingNextPage,
    carregarMais: () => void consulta.fetchNextPage(),
    recarregar: () => void consulta.refetch(),
  };
}
