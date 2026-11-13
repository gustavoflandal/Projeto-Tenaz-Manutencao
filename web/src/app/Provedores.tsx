import { QueryClientProvider, type QueryClient } from '@tanstack/react-query';
import type { ReactNode } from 'react';

import '../lib/i18n';

export interface ProvedoresProps {
  readonly clienteDeConsultas: QueryClient;
  readonly children: ReactNode;
}

/** Provedores comuns à aplicação e aos testes de componente. */
export function Provedores({ clienteDeConsultas, children }: ProvedoresProps) {
  return <QueryClientProvider client={clienteDeConsultas}>{children}</QueryClientProvider>;
}
