// Dados do diálogo "Sobre" (componente 6.19, RNF-039): versão da web, versão
// da API, data do build, commit curto e ambiente.

import { useQuery } from '@tanstack/react-query';
import { useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';

import { chaves } from '../../lib/consultas';
import { formatarDataHora } from '../../lib/formatacao';
import {
  COMMIT_DESCONHECIDO,
  mesmaVersao,
  rotuloDoAmbiente,
  versaoDaAPI,
  versaoDaWeb,
  type VersaoDaWeb,
} from '../../lib/versao';

export type EstadoDaAPI = 'consultando' | 'respondeu' | 'sem-resposta';

export interface LinhaSobre {
  readonly chave: 'web' | 'api' | 'dataBuild' | 'commit' | 'ambiente';
  readonly rotulo: string;
  readonly valor: string;
}

export interface InformacoesVersao {
  readonly linhas: readonly LinhaSobre[];
  readonly estadoDaAPI: EstadoDaAPI;
  /** A API está numa versão diferente da web: o navegador ficou com a antiga. */
  readonly versaoNovaDisponivel: boolean;
  readonly textoParaCopiar: string;
}

export function useInformacoesVersao(web: VersaoDaWeb = versaoDaWeb): InformacoesVersao {
  const { t } = useTranslation();
  // A consulta roda a cada abertura do diálogo; sem nova tentativa, para dizer
  // logo que a API não respondeu. O último valor visto fica no cache.
  const consulta = useQuery({ queryKey: chaves.versaoDaAPI, queryFn: versaoDaAPI, retry: false });

  let estadoDaAPI: EstadoDaAPI = 'respondeu';
  if (consulta.isError) estadoDaAPI = 'sem-resposta';
  else if (consulta.isPending) estadoDaAPI = 'consultando';

  let valorDaAPI: string;
  if (estadoDaAPI === 'respondeu' && consulta.data) {
    valorDaAPI = t('sobre.versao', { versao: consulta.data.versao });
  } else if (estadoDaAPI === 'consultando') {
    valorDaAPI = t('sobre.apiConsultando');
  } else if (consulta.data) {
    valorDaAPI = t('sobre.apiUltimaVista', { versao: consulta.data.versao });
  } else {
    valorDaAPI = t('sobre.apiSemResposta');
  }

  const linhas: LinhaSobre[] = [
    { chave: 'web', rotulo: t('sobre.web'), valor: t('sobre.versao', { versao: web.versao }) },
    { chave: 'api', rotulo: t('sobre.api'), valor: valorDaAPI },
    {
      chave: 'dataBuild',
      rotulo: t('sobre.dataBuild'),
      valor: formatarDataHora(web.dataBuild) ?? t('sobre.dataBuildAusente'),
    },
    {
      chave: 'commit',
      rotulo: t('sobre.commit'),
      valor: web.commit === COMMIT_DESCONHECIDO || web.commit === '' ? t('sobre.commitAusente') : web.commit,
    },
    { chave: 'ambiente', rotulo: t('sobre.ambiente'), valor: rotuloDoAmbiente(web.ambiente) },
  ];

  return {
    linhas,
    estadoDaAPI,
    versaoNovaDisponivel:
      estadoDaAPI === 'respondeu' && consulta.data !== undefined && !mesmaVersao(consulta.data.versao, web.versao),
    textoParaCopiar: linhas.map((l) => t('sobre.linhaCopia', { rotulo: l.rotulo, valor: l.valor })).join('\n'),
  };
}

export type EstadoCopia = 'ocioso' | 'copiado' | 'falhou';

/** Copia texto para a área de transferência e informa o resultado. */
export function useCopiar(): { readonly estado: EstadoCopia; readonly copiar: (texto: string) => Promise<void> } {
  const [estado, setEstado] = useState<EstadoCopia>('ocioso');
  const copiar = useCallback(async (texto: string) => {
    try {
      if (!globalThis.navigator.clipboard) throw new Error('área de transferência indisponível');
      await globalThis.navigator.clipboard.writeText(texto);
      setEstado('copiado');
    } catch {
      setEstado('falhou');
    }
  }, []);
  return { estado, copiar };
}
