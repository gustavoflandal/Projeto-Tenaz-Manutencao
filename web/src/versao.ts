// Versão da web em tempo de execução.
//
// Os valores são substituídos no build pelo Vite, a partir das variáveis que o
// CI deriva da tag git. Nenhum deles é digitado à mão (RNF-039).

declare const __TENAZ_VERSAO__: string;
declare const __TENAZ_COMMIT__: string;
declare const __TENAZ_DATA_BUILD__: string;
declare const __TENAZ_AMBIENTE__: string;

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

/**
 * Texto do selo de versão, exibido no canto inferior direito da tela principal
 * e da tela de login (RNF-039).
 *
 * Fora da produção, acompanha o ambiente, para que nenhuma captura de tela de
 * homologação seja confundida com produção.
 */
export function seloDeVersao(v: VersaoDaWeb = versaoDaWeb): string {
  const base = `v${v.versao}`;
  return v.ambiente === 'producao' ? base : `${base} · ${rotuloDoAmbiente(v.ambiente)}`;
}

function rotuloDoAmbiente(ambiente: string): string {
  switch (ambiente) {
    case 'homologacao':
      return 'Homologação';
    case 'desenvolvimento':
      return 'Desenvolvimento';
    default:
      return ambiente;
  }
}

/** Resposta de GET /v1/versao. */
export interface VersaoDaAPI {
  readonly versao: string;
  readonly commit: string;
  readonly data_build: string;
}

/** Consulta a versão publicada da API. */
export async function versaoDaAPI(): Promise<VersaoDaAPI> {
  const resposta = await fetch('/v1/versao');
  if (!resposta.ok) {
    throw new Error(`GET /v1/versao devolveu ${resposta.status}`);
  }
  return (await resposta.json()) as VersaoDaAPI;
}
