import { useState } from 'react';
import { createBrowserRouter, RouterProvider } from 'react-router';

import { criarClienteDeConsultas } from '../lib/consultas';

import { Provedores } from './Provedores';
import { rotas } from './rotas';

export function Aplicacao() {
  const [clienteDeConsultas] = useState(criarClienteDeConsultas);
  const [roteador] = useState(() => createBrowserRouter(rotas));
  return (
    <Provedores clienteDeConsultas={clienteDeConsultas}>
      <RouterProvider router={roteador} />
    </Provedores>
  );
}
