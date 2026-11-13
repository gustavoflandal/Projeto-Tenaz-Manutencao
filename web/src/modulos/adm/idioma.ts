// Idioma do usuário (ADR-013): depois do login vale `usuario.idioma`; a escolha
// feita na interface é gravada nele por PATCH /v1/usuarios/eu. Sem sessão (login,
// criar conta), a escolha fica só no navegador (lib/idioma.ts).

import { useQueryClient } from '@tanstack/react-query';
import type { Idioma } from '@tenaz/i18n';
import { useCallback, useEffect } from 'react';

import { chaves } from '../../lib/consultas';
import { definirIdioma, i18n } from '../../lib/i18n';
import type { Sessao } from '../../lib/tipos';

import { chaveUsuarioEu, lerUsuarioEu, useAlterarUsuarioEu, useUsuarioEu } from './usuarios';

/** Traz o idioma do usuário para a interface quando ele muda no servidor (login, outra aba). */
export function useAplicarIdiomaDoUsuario(): void {
  const { data } = useUsuarioEu();
  const idioma = data?.idioma;
  useEffect(() => {
    if (idioma && idioma !== i18n.language) void definirIdioma(idioma);
  }, [idioma]);
}

/** Troca o idioma da interface e, havendo sessão plena, grava a escolha no usuário. */
export function useEscolherIdioma(): (idioma: Idioma) => Promise<void> {
  const cliente = useQueryClient();
  const { mutate } = useAlterarUsuarioEu();

  return useCallback(
    async (idioma: Idioma) => {
      await definirIdioma(idioma);
      // Lê o cache na hora do clique, sem observar a sessão: nas telas sem sessão isso não dispara GET /v1/auth/sessao.
      const sessao = cliente.getQueryData<Sessao | null>(chaves.sessao);
      if (!sessao || sessao.mfa_pendente) return;
      try {
        const eu = await cliente.ensureQueryData({ queryKey: chaveUsuarioEu, queryFn: lerUsuarioEu });
        if (eu.idioma !== idioma) mutate({ versao: eu.versao, corpo: { idioma } });
      } catch {
        // Sem o usuário ou sem rede, a escolha vale neste navegador; a próxima carga do usuário a corrige.
      }
    },
    [cliente, mutate],
  );
}

