// Formatadores presos à empresa da sessão (ADR-013, Parte G): o país vem da
// sessão, o fuso da empresa (ou o da unidade, quando o dado tem unidade). Dado
// do servidor, só no cache de consultas; nenhum estado paralelo.

import { PERFIS_DE_PAIS, fusoDoAparelho, type Pais } from '@tenaz/dominio';
import { useCallback, useMemo } from 'react';

import { formatarData, formatarDataHora, formatarPrecoUnitario, simboloDaMoeda } from '../../lib/formatacao';

import { useEmpresa } from './empresa';
import { useSessao } from './sessao';

/** País da empresa da sessão; sem sessão (login, convite), o do idioma de origem. */
export function usePaisDaEmpresa(): Pais {
  const { data } = useSessao();
  return data?.empresa.pais ?? 'BR';
}

export interface Formatos {
  readonly pais: Pais;
  readonly data: (iso: string) => string;
  readonly precoUnitario: (texto: string) => string;
  readonly simboloDaMoeda: string;
}

/** Formatos que dependem só do país da empresa. */
export function useFormatos(): Formatos {
  const pais = usePaisDaEmpresa();
  return useMemo(
    () => ({
      pais,
      data: (iso) => formatarData(iso, pais),
      precoUnitario: (texto) => formatarPrecoUnitario(texto, pais),
      simboloDaMoeda: simboloDaMoeda(pais),
    }),
    [pais],
  );
}

/**
 * Data e hora no fuso da unidade do dado ou, sem unidade, no fuso da empresa.
 * Enquanto a empresa não carrega, vale o fuso padrão do país; nunca o do navegador.
 */
export function useDataHora(): (instante: string, fusoDaUnidade?: string | null) => string | null {
  const pais = usePaisDaEmpresa();
  const { data: empresa } = useEmpresa();
  const fusoDaEmpresa = empresa?.fuso_horario ?? PERFIS_DE_PAIS[pais].fuso_padrao;
  return useCallback(
    (instante, fusoDaUnidade) => formatarDataHora(instante, fusoDaUnidade ?? fusoDaEmpresa, pais),
    [fusoDaEmpresa, pais],
  );
}

/**
 * Data e hora no fuso do aparelho, que ali é explícito: o "Sobre" da web mostra a
 * data do build e a tela de login não tem empresa.
 */
export function formatarNoFusoDoAparelho(instante: string, pais: Pais): string | null {
  return formatarDataHora(instante, fusoDoAparelho(), pais);
}
