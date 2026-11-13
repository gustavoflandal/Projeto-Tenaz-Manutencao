// Versão da web em tempo de execução.
//
// Os valores são substituídos no build pelo Vite, a partir das variáveis que o
// CI deriva da tag git. Nenhum deles é digitado à mão (RNF-039).

import type { VersaoDaAPI } from '@tenaz/api-tipos';

import { api, ErroDaAPI } from './api';
import { i18n } from './i18n';

declare const __TENAZ_VERSAO__: string;
declare const __TENAZ_COMMIT__: string;
declare const __TENAZ_DATA_BUILD__: string;
declare const __TENAZ_AMBIENTE__: string;

export type { VersaoDaAPI };

export interface VersaoDaWeb {
  readonly versao: string;
  readonly commit: string;
  readonly dataBuild: string;
  readonly ambiente: string;
}

export const versaoDaWeb: VersaoDaWeb = {
  versao: __TENAZ_VERSAO__,
  commit: __TENAZ_COMMIT__,
  dataBuild: __TENAZ_DATA_BUILD__,
  ambiente: __TENAZ_AMBIENTE__,
};

/** Valor que o build usa quando o CI não informa o commit. */
export const COMMIT_DESCONHECIDO = 'desconhecido';

/**
 * Versão sem o metadado de build (`+a1b2c3d`). O metadado não aparece no selo
 * e não conta para dizer se uma versão é outra (SemVer 2.0.0, item 10); ele
 * fica só no diálogo Sobre (componente 6.19).
 */
export function versaoSemMetadado(versao: string): string {
  const inicio = versao.indexOf('+');
  return inicio === -1 ? versao : versao.slice(0, inicio);
}

/** Web e API estão no mesmo release, ignorando o metadado de build. */
export function mesmaVersao(a: string, b: string): boolean {
  return versaoSemMetadado(a) === versaoSemMetadado(b);
}

/**
 * Texto do selo de versão, exibido no canto inferior direito da tela principal
 * e da tela de login (RNF-039).
 *
 * Fora da produção, acompanha o ambiente, para que nenhuma captura de tela de
 * homologação seja confundida com produção.
 */
export function seloDeVersao(v: VersaoDaWeb = versaoDaWeb): string {
  const versao = versaoSemMetadado(v.versao);
  return v.ambiente === 'producao'
    ? i18n.t('selo.texto', { versao })
    : i18n.t('selo.textoComAmbiente', { versao, ambiente: rotuloDoAmbiente(v.ambiente) });
}

/** Nome acessível do selo, como manda o componente 6.19 do documento 7. */
export function nomeAcessivelDoSelo(v: VersaoDaWeb = versaoDaWeb): string {
  const versao = versaoSemMetadado(v.versao);
  return v.ambiente === 'producao'
    ? i18n.t('selo.nomeAcessivel', { versao })
    : i18n.t('selo.nomeAcessivelComAmbiente', { versao, ambiente: rotuloDoAmbiente(v.ambiente) });
}

export function rotuloDoAmbiente(ambiente: string): string {
  switch (ambiente) {
    case 'producao':
      return i18n.t('ambiente.producao');
    case 'homologacao':
      return i18n.t('ambiente.homologacao');
    case 'desenvolvimento':
      return i18n.t('ambiente.desenvolvimento');
    default:
      return ambiente;
  }
}

/** Consulta a versão publicada da API (GET /v1/versao). */
export async function versaoDaAPI(): Promise<VersaoDaAPI> {
  const { data, response } = await api.GET('/v1/versao');
  if (!data) {
    throw new ErroDaAPI(`GET /v1/versao devolveu ${response.status}`, response.status);
  }
  return data;
}
