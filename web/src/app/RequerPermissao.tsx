// Rota que exige uma permissão. Sem ela, mostra o estado "sem permissão" em vez
// de chamar uma API que responderia 403. A decisão final é do servidor, que
// confere a permissão em cada requisição.

import { Outlet } from 'react-router';

import { EstadoSemPermissao } from '../componentes/Estados';
import { useTemPermissao } from '../modulos/adm/sessao';

export function RequerPermissao({ permissao }: { readonly permissao: string }) {
  const tem = useTemPermissao(permissao);
  return tem ? <Outlet /> : <EstadoSemPermissao />;
}
