// Endereço multipaís (documento 3, seção 2.18): subdivisão (ISO 3166-2 completo,
// lista fechada por país) e código postal. Espelha backend/internal/plataforma/pais,
// provado pelos casos compartilhados (backend/testes/casos/identificacao.json, versão 2).

import type { Pais } from './paises';

export const SUBDIVISOES: Readonly<Record<Pais, readonly string[]>> = {
  BR: [
    'BR-AC', 'BR-AL', 'BR-AP', 'BR-AM', 'BR-BA', 'BR-CE', 'BR-DF', 'BR-ES', 'BR-GO', 'BR-MA', 'BR-MT', 'BR-MS', 'BR-MG',
    'BR-PA', 'BR-PB', 'BR-PR', 'BR-PE', 'BR-PI', 'BR-RJ', 'BR-RN', 'BR-RS', 'BR-RO', 'BR-RR', 'BR-SC', 'BR-SP', 'BR-SE', 'BR-TO',
  ],
  AR: [
    'AR-A', 'AR-B', 'AR-C', 'AR-D', 'AR-E', 'AR-F', 'AR-G', 'AR-H', 'AR-J', 'AR-K', 'AR-L', 'AR-M',
    'AR-N', 'AR-P', 'AR-Q', 'AR-R', 'AR-S', 'AR-T', 'AR-U', 'AR-V', 'AR-W', 'AR-X', 'AR-Y', 'AR-Z',
  ],
  PY: [
    'PY-ASU', 'PY-1', 'PY-2', 'PY-3', 'PY-4', 'PY-5', 'PY-6', 'PY-7', 'PY-8', 'PY-9', 'PY-10', 'PY-11', 'PY-12', 'PY-13',
    'PY-14', 'PY-15', 'PY-16', 'PY-19',
  ],
};

/** Só a forma exata do código (sem espaço, sem minúscula, sem zero à esquerda). */
export function subdivisaoValida(pais: Pais, codigo: string): boolean {
  return SUBDIVISOES[pais].includes(codigo);
}

const FORMATO_DO_CODIGO_POSTAL: Readonly<Record<Pais, RegExp>> = {
  BR: /^\d{8}$/,
  AR: /^([A-HJ-NP-Z]\d{4}[A-Z]{3}|\d{4})$/,
  PY: /^\d{6}$/,
};

/** Maiúsculas, sem espaço nem hífen. */
export function normalizarCodigoPostal(entrada: string): string {
  return entrada.replaceAll(/[\s-]/g, '').toUpperCase();
}

export type ResultadoDoCodigoPostal =
  | { readonly valido: true; readonly normalizado: string }
  | { readonly valido: false };

export function validarCodigoPostal(pais: Pais, entrada: string): ResultadoDoCodigoPostal {
  const normalizado = normalizarCodigoPostal(entrada);
  return FORMATO_DO_CODIGO_POSTAL[pais].test(normalizado) ? { valido: true, normalizado } : { valido: false };
}

/**
 * Argentina: o CPA de 8 caracteres começa pela letra da província, que deve ser a
 * da subdivisão informada. Sem subdivisão, com código de 4 dígitos ou fora da
 * Argentina não há cruzamento.
 */
export function cpaConfereComSubdivisao(subdivisao: string, codigoPostal: string): boolean {
  const cpa = normalizarCodigoPostal(codigoPostal);
  if (subdivisao === '' || cpa.length !== 8 || !subdivisao.startsWith('AR-')) return true;
  return cpa.startsWith(subdivisao.slice(3));
}
