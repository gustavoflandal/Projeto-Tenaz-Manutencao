// Renderiza uma rota da aplicação com os mesmos provedores do navegador, mas
// com roteador em memória e consultas sem nova tentativa.

import { QueryClient } from '@tanstack/react-query';
import { render } from '@testing-library/react';
import { createMemoryRouter, RouterProvider } from 'react-router';
import { vi } from 'vitest';

import { Provedores } from '../app/Provedores';
import { rotas } from '../app/rotas';

export function renderizarRota(caminho: string) {
  const clienteDeConsultas = new QueryClient({ defaultOptions: { queries: { retry: false } } });
  const roteador = createMemoryRouter(rotas, { initialEntries: [caminho] });
  return render(
    <Provedores clienteDeConsultas={clienteDeConsultas}>
      <RouterProvider router={roteador} />
    </Provedores>,
  );
}

export interface CorpoVersao {
  readonly versao: string;
  readonly commit: string;
  readonly data_build: string;
}

/** Faz GET /v1/versao responder o corpo dado. */
export function apiRespondeVersao(corpo: CorpoVersao) {
  return vi.spyOn(globalThis, 'fetch').mockImplementation(async () => Response.json(corpo));
}

/** Faz toda chamada à API falhar como falha de rede. */
export function apiForaDoAr() {
  return vi.spyOn(globalThis, 'fetch').mockRejectedValue(new TypeError('Failed to fetch'));
}
