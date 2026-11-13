// Apoio dos campos de identificação (ADR-013): máscara enquanto digita e
// teclado do celular conforme o tipo. As regras vêm de @tenaz/dominio.

import { mascararEntrada, normalizarIdentificacao, type TipoDeIdentificacao } from '@tenaz/dominio';
import type { ChangeEvent } from 'react';

/** CNPJ e RUC podem ter letra; os demais são só dígitos. */
export function modoDeEntrada(tipo: TipoDeIdentificacao): 'text' | 'numeric' {
  return tipo === 'cnpj' || tipo === 'ruc' ? 'text' : 'numeric';
}

/** Manipulador de `onChange` que aplica a máscara do tipo vigente (lido na hora, pois o tipo pode mudar). */
export function aplicarMascara(tipoAtual: () => TipoDeIdentificacao) {
  return (e: ChangeEvent<HTMLInputElement>): void => {
    e.target.value = mascararEntrada(tipoAtual(), e.target.value);
  };
}

/** Valor do servidor (sem máscara) ou vazio, pronto para o campo. */
export function paraOCampo(tipo: TipoDeIdentificacao, valor: string | null | undefined): string {
  return valor ? mascararEntrada(tipo, valor) : '';
}

/** Valor enviado ao servidor: normalizado (maiúsculas, sem máscara). */
export function paraOEnvio(valor: string): string {
  return normalizarIdentificacao(valor);
}
