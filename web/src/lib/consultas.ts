// Cliente do TanStack Query: a única camada que guarda dado do servidor na
// web. Nenhum estado global mantém cópia de dado do servidor (seção 8.4 de
// docs/5_ARQUITETURA_SOFTWARE.md).

import { QueryClient } from '@tanstack/react-query';

export function criarClienteDeConsultas(): QueryClient {
  return new QueryClient({
    defaultOptions: {
      queries: {
        retry: 1,
      },
    },
  });
}

/** Chaves de consulta, num só lugar, para invalidação previsível. */
export const chaves = {
  versaoDaAPI: ['plataforma', 'versao'] as const,
};
