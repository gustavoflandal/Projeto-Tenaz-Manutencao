// Token no fragmento da URL (`#token=...`), como nos links de e-mail
// (seção 4.2 do documento 3). O fragmento nunca vai ao servidor nem a log; a
// página lê o token, tira-o da barra de endereço e o envia no corpo da chamada.

import { useEffect, useState } from 'react';
import { useLocation, useNavigate } from 'react-router';

const TAMANHO_DO_TOKEN = 43;

/** Lê `token` de um fragmento (`#token=abc`). Devolve null se faltar ou tiver tamanho errado. */
export function tokenDoFragmento(fragmento: string): string | null {
  const parametros = new URLSearchParams(fragmento.startsWith('#') ? fragmento.slice(1) : fragmento);
  const token = parametros.get('token');
  return token !== null && token.length === TAMANHO_DO_TOKEN ? token : null;
}

/**
 * Token do fragmento da página atual. Guarda o valor em estado e limpa o
 * fragmento da URL, para que ele não fique no histórico nem na barra.
 */
export function useTokenDoFragmento(): string | null {
  const { hash, pathname, search } = useLocation();
  const navegar = useNavigate();
  const [token] = useState(() => tokenDoFragmento(hash));
  useEffect(() => {
    if (hash !== '') void navegar({ pathname, search }, { replace: true });
  }, [hash, pathname, search, navegar]);
  return token;
}
