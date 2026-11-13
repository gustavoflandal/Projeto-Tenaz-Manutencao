import { useEffect, useState } from 'react';
import { createBrowserRouter, RouterProvider } from 'react-router';

import { registrarFonteDaTela } from '../lib/api';
import { criarClienteDeConsultas, ligarCSRFAoCliente } from '../lib/consultas';
import { telaDoCaminho } from '../lib/telas';

import { Provedores } from './Provedores';
import { rotas } from './rotas';

export function Aplicacao() {
  const [clienteDeConsultas] = useState(criarClienteDeConsultas);
  const [roteador] = useState(() => createBrowserRouter(rotas));
  // Fontes do cabeçalho de cada chamada: o token CSRF da sessão e a tela atual.
  useEffect(() => {
    ligarCSRFAoCliente(clienteDeConsultas);
    registrarFonteDaTela(() => telaDoCaminho(roteador.state.location.pathname));
  }, [clienteDeConsultas, roteador]);
  return (
    <Provedores clienteDeConsultas={clienteDeConsultas}>
      <RouterProvider router={roteador} />
    </Provedores>
  );
}
