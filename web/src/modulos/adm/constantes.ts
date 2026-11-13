// Valores fechados do contrato usados em seletores. Os nomes de exibição ficam
// em @tenaz/i18n (`planos.*`, `fusos.*`, `vinculos.*`).

import { PERFIS_DE_PAIS, type Pais } from '@tenaz/dominio';

/** Planos de `NovaEmpresa.plano` (RF-ADM-001). */
export const PLANOS = ['essencial', 'profissional', 'enterprise'] as const;

/** Fusos horários IANA oferecidos por país (ADR-013): a lista que o servidor aceita, vinda do perfil de país. */
export const FUSOS_POR_PAIS: Readonly<Record<Pais, readonly string[]>> = {
  BR: PERFIS_DE_PAIS.BR.fusos,
  AR: PERFIS_DE_PAIS.AR.fusos,
  PY: PERFIS_DE_PAIS.PY.fusos,
};

/**
 * Fusos para um seletor: os do país e, se a unidade já tem um fuso fora da
 * lista (cadastro antigo), ele também, para não ser trocado em silêncio.
 */
export function fusosParaSeletor(pais: Pais | null, atual?: string | null): readonly string[] {
  const lista = FUSOS_POR_PAIS[pais ?? 'BR'];
  return atual && !lista.includes(atual) ? [...lista, atual] : lista;
}

/** Vínculos de `pessoa.vinculo` (sumário 3.6). */
export const VINCULOS = ['empregado', 'terceiro', 'autonomo'] as const;

/** Tipos de operação da trilha (`evento.acao`, lista fechada do sumário 3.6). */
export const ACOES_DE_AUDITORIA = [
  'login',
  'logout',
  'inclusao',
  'alteracao',
  'exclusao',
  'delecao',
  'visualizacao',
  'impressao',
  'extracao',
] as const;

/** Chaves de captura por tela: o que a empresa liga ou desliga (RF-ADM-011). */
export const CAPTURAS_CONFIGURAVEIS = ['visualizacao', 'impressao', 'extracao'] as const;
export type CapturaConfiguravel = (typeof CAPTURAS_CONFIGURAVEIS)[number];
