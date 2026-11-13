// Renderiza uma rota da aplicação com os mesmos provedores do navegador, mas
// com roteador em memória e consultas sem nova tentativa.

import { render } from '@testing-library/react';
import { createMemoryRouter, RouterProvider } from 'react-router';
import { vi } from 'vitest';

import { Provedores } from '../app/Provedores';
import { registrarFonteDaTela } from '../lib/api';
import { i18n } from '../lib/i18n';
import { criarClienteDeConsultas, ligarCSRFAoCliente } from '../lib/consultas';
import { telaDoCaminho } from '../lib/telas';
import type { Sessao } from '../lib/tipos';
import { problema, simularAPI, type Manipulador } from './api-simulada';
import { sessao, usuario } from './dados';
import { rotas } from '../app/rotas';

export function renderizarRota(caminho: string) {
  // O mesmo cliente do navegador (com o tratamento de sessão perdida), sem nova tentativa.
  const clienteDeConsultas = criarClienteDeConsultas();
  ligarCSRFAoCliente(clienteDeConsultas);
  clienteDeConsultas.setDefaultOptions({
    ...clienteDeConsultas.getDefaultOptions(),
    queries: { retry: false, refetchOnWindowFocus: false },
  });
  const roteador = createMemoryRouter(rotas, { initialEntries: [caminho] });
  registrarFonteDaTela(() => telaDoCaminho(roteador.state.location.pathname));
  const resultado = render(
    <Provedores clienteDeConsultas={clienteDeConsultas}>
      <RouterProvider router={roteador} />
    </Provedores>,
  );
  return { ...resultado, roteador, clienteDeConsultas };
}

/** A rota da sessão respondendo a sessão dada (ou 401, sem sessão). */
export function rotaDaSessao(sessao: Sessao | null): Record<string, Manipulador> {
  return {
    'GET /v1/auth/sessao': sessao ? { corpo: sessao } : problema(401, 'nao_autenticado'),
    // Com sessão, a interface lê o próprio usuário (idioma e moeda da empresa).
    ...(sessao
      ? {
          // Função, não valor: o idioma do usuário é o que a interface usa na hora da chamada
          // (o teste que troca de idioma antes de entrar o mantém).
          'GET /v1/usuarios/eu': () => ({
            corpo: usuario({
              id: sessao.usuario.id,
              nome: sessao.usuario.nome,
              email: sessao.usuario.email,
              idioma: i18n.language === 'es-419' ? 'es-419' : 'pt-BR',
              versao: 1,
            }),
          }),
        }
      : {}),
  };
}

export interface CorpoVersao {
  readonly versao: string;
  readonly commit: string;
  readonly data_build: string;
}

/** Faz GET /v1/versao responder o corpo dado. */
export function apiRespondeVersao(corpo: CorpoVersao) {
  return simularAPI({ ...rotaDaSessao(sessao()), 'GET /v1/versao': { corpo } });
}

/** Faz toda chamada à API falhar como falha de rede. */
export function apiForaDoAr() {
  return vi.spyOn(globalThis, 'fetch').mockRejectedValue(new TypeError('Failed to fetch'));
}
