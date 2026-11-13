// Região da empresa (ADR-013): país, moeda e perfil de identificação.
//
// País e moeda vêm do próprio usuário (GET /v1/usuarios/eu, que serve a todos
// os escopos); a interface lê o país e nunca o deriva da moeda (documento 3,
// seção 2.17). Dado do servidor: vive só no cache de consultas.

import { casasDecimaisDaMoeda, PERFIS_DE_PAIS, type Moeda, type Pais, type PerfilDePais } from '@tenaz/dominio';
import { useTranslation } from 'react-i18next';

import { useFormatos } from './formatos';
import { useUsuarioEu } from './usuarios';

export interface Regiao {
  /** Falso enquanto o usuário não carregou. */
  readonly pronta: boolean;
  readonly moeda: Moeda | null;
  readonly pais: Pais | null;
  readonly perfil: PerfilDePais | null;
}

export function useRegiao(): Regiao {
  const { data } = useUsuarioEu();
  const moeda = data?.moeda ?? null;
  const pais = data?.pais ?? null;
  return { pronta: moeda !== null, moeda, pais, perfil: pais ? PERFIS_DE_PAIS[pais] : null };
}

/** Símbolo da moeda da empresa para o rótulo dos campos de valor, e se ela tem centavos. */
export function useMoedaDeCusto(): { readonly simbolo: string; readonly semCentavos: boolean } {
  const { t } = useTranslation();
  const { moeda } = useRegiao();
  const formatos = useFormatos();
  return {
    simbolo: moeda ? formatos.simboloDaMoeda : t('moeda.rotulo'),
    semCentavos: moeda !== null && casasDecimaisDaMoeda(moeda) === 0,
  };
}
