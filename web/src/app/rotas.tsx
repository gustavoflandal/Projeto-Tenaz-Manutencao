// Rotas da aplicação principal. O link de uso único (/l/#<token>) não está
// aqui: tem entrada e pacote próprios, em src/link/ (seção 6.7 do documento 5).

import type { RouteObject } from 'react-router';

import { TelaEntrar } from '../modulos/adm/TelaEntrar';
import { TelaPainel } from '../modulos/ind/TelaPainel';

import { LayoutPrincipal } from './LayoutPrincipal';
import { TelaNaoEncontrada } from './TelaNaoEncontrada';

export const rotas: RouteObject[] = [
  { path: '/entrar', element: <TelaEntrar /> },
  {
    path: '/',
    element: <LayoutPrincipal />,
    children: [
      { index: true, element: <TelaPainel /> },
      { path: '*', element: <TelaNaoEncontrada /> },
    ],
  },
];
