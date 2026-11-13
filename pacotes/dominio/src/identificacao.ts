// Identificação por país (ADR-013): as mesmas regras de
// backend/internal/plataforma/pais, provadas pelos mesmos casos
// (backend/testes/casos/identificacao.json). O servidor decide; aqui o
// cliente só avisa antes de enviar e desenha a máscara.

export const TIPOS_DE_IDENTIFICACAO = ['cnpj', 'cpf', 'cuit', 'cuil', 'dni', 'ruc', 'ci'] as const;
export type TipoDeIdentificacao = (typeof TIPOS_DE_IDENTIFICACAO)[number];

/** Tipos de pessoa jurídica: o que a empresa e a unidade aceitam (um por país). */
export const TIPOS_DE_EMPRESA = ['cnpj', 'cuit', 'ruc'] as const;
export type TipoDeEmpresa = (typeof TIPOS_DE_EMPRESA)[number];

/** Tipos de pessoa física: a pessoa aceita os de qualquer país. */
export const TIPOS_DE_PESSOA = ['cpf', 'dni', 'cuil', 'ci'] as const;
export type TipoDePessoa = (typeof TIPOS_DE_PESSOA)[number];

export type MotivoDeRecusa = 'formato' | 'digito_verificador' | 'tipo_desconhecido';

export function ehTipoDeIdentificacao(valor: unknown): valor is TipoDeIdentificacao {
  return typeof valor === 'string' && (TIPOS_DE_IDENTIFICACAO as readonly string[]).includes(valor);
}

export function ehTipoDeEmpresa(valor: unknown): valor is TipoDeEmpresa {
  return typeof valor === 'string' && (TIPOS_DE_EMPRESA as readonly string[]).includes(valor);
}

export function ehTipoDePessoa(valor: unknown): valor is TipoDePessoa {
  return typeof valor === 'string' && (TIPOS_DE_PESSOA as readonly string[]).includes(valor);
}

/** Maiúsculas e só letras e dígitos ASCII: a máscara é só de exibição. */
export function normalizarIdentificacao(valor: string): string {
  return valor.replaceAll(/[^A-Za-z0-9]/g, '').toUpperCase();
}

interface Forma {
  readonly minimo: number;
  readonly maximo: number;
  /** Letras só até esta posição (exclusiva); 0 = só dígitos. */
  readonly letrasAte: number;
}

const FORMAS: Readonly<Record<TipoDeIdentificacao, Forma>> = {
  // Letras só nas 12 primeiras posições; os dois dígitos verificadores são numéricos.
  cnpj: { minimo: 14, maximo: 14, letrasAte: 12 },
  cpf: { minimo: 11, maximo: 11, letrasAte: 0 },
  cuit: { minimo: 11, maximo: 11, letrasAte: 0 },
  cuil: { minimo: 11, maximo: 11, letrasAte: 0 },
  dni: { minimo: 7, maximo: 8, letrasAte: 0 },
  // RUC antigo pode ter letra (P13).
  ruc: { minimo: 6, maximo: 10, letrasAte: 10 },
  ci: { minimo: 4, maximo: 8, letrasAte: 0 },
};

function formaValida(tipo: TipoDeIdentificacao, n: string): boolean {
  const forma = FORMAS[tipo];
  if (n.length < forma.minimo || n.length > forma.maximo) return false;
  return [...n].every((c, i) => /[0-9]/.test(c) || (i < forma.letrasAte && /[A-Z]/.test(c)));
}

function repetido(n: string): boolean {
  return [...n].every((c) => c === n[0]);
}

/** Valor de um caractere no cálculo: ASCII menos 48 ('0' a '9' valem 0 a 9, 'A' vale 17). */
function valorDoCaractere(c: string | undefined): number {
  return (c?.codePointAt(0) ?? 48) - 48;
}

function cpfConfere(cpf: string): boolean {
  if (repetido(cpf)) return false;
  const dv = (n: number): number => {
    let soma = 0;
    for (let i = 0; i < n; i++) soma += valorDoCaractere(cpf[i]) * (n + 1 - i);
    const r = (soma * 10) % 11;
    return r === 10 ? 0 : r;
  };
  return dv(9) === valorDoCaractere(cpf[9]) && dv(10) === valorDoCaractere(cpf[10]);
}

const PESOS_CNPJ_1 = [5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2] as const;
const PESOS_CNPJ_2 = [6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2] as const;

function cnpjConfere(cnpj: string): boolean {
  if (repetido(cnpj)) return false;
  const dv = (pesos: readonly number[]): number => {
    const soma = pesos.reduce((acc, peso, i) => acc + valorDoCaractere(cnpj[i]) * peso, 0);
    const r = soma % 11;
    return r < 2 ? 0 : 11 - r;
  };
  return dv(PESOS_CNPJ_1) === valorDoCaractere(cnpj[12]) && dv(PESOS_CNPJ_2) === valorDoCaractere(cnpj[13]);
}

export type ResultadoDaValidacao =
  | { readonly valido: true; readonly normalizado: string }
  | { readonly valido: false; readonly motivo: MotivoDeRecusa };

/**
 * Normaliza e confere a forma do tipo e, para CNPJ (inclusive o alfanumérico)
 * e CPF, o dígito verificador. CUIT e RUC só pela forma, até a pendência P13.
 */
export function validarIdentificacao(tipo: string, valor: string): ResultadoDaValidacao {
  if (!ehTipoDeIdentificacao(tipo)) return { valido: false, motivo: 'tipo_desconhecido' };
  const n = normalizarIdentificacao(valor);
  if (!formaValida(tipo, n)) return { valido: false, motivo: 'formato' };
  if (tipo === 'cnpj' && !cnpjConfere(n)) return { valido: false, motivo: 'digito_verificador' };
  if (tipo === 'cpf' && !cpfConfere(n)) return { valido: false, motivo: 'digito_verificador' };
  return { valido: true, normalizado: n };
}

export function identificacaoValida(tipo: string, valor: string): boolean {
  return validarIdentificacao(tipo, valor).valido;
}

function agruparMilhares(n: string): string {
  return n.replaceAll(/\B(?=(\d{3})+(?!\d))/g, '.');
}

/**
 * Máscara de entrada: o texto digitado (com ou sem pontuação) como o tipo o
 * escreve. Descarta o que não é letra ou dígito e o que passa do tamanho.
 * CNPJ 12.ABC.345/01DE-35, CPF 529.982.247-25, CUIT e CUIL 33-69345023-9,
 * DNI e CI com milhares (12.345.678), RUC com hífen antes do último caractere.
 */
export function mascararEntrada(tipo: string, texto: string): string {
  if (!ehTipoDeIdentificacao(tipo)) return texto;
  const n = normalizarIdentificacao(texto).slice(0, FORMAS[tipo].maximo);
  switch (tipo) {
    case 'cnpj': {
      const base = [n.slice(0, 2), n.slice(2, 5), n.slice(5, 8)].filter(Boolean).join('.');
      return base + (n.length > 8 ? `/${n.slice(8, 12)}` : '') + (n.length > 12 ? `-${n.slice(12)}` : '');
    }
    case 'cpf': {
      const base = [n.slice(0, 3), n.slice(3, 6), n.slice(6, 9)].filter(Boolean).join('.');
      return base + (n.length > 9 ? `-${n.slice(9)}` : '');
    }
    case 'cuit':
    case 'cuil':
      if (n.length <= 2) return n;
      return n.length <= 10 ? `${n.slice(0, 2)}-${n.slice(2)}` : `${n.slice(0, 2)}-${n.slice(2, 10)}-${n.slice(10)}`;
    case 'dni':
    case 'ci':
      return agruparMilhares(n);
    case 'ruc':
      return n.length <= 1 ? n : `${n.slice(0, -1)}-${n.slice(-1)}`;
  }
}

/**
 * Máscara de lista: a mesma que o servidor devolve em `identificacao_mascarada`.
 * CPF mostra o miolo (***.982.247-**); os demais, os 3 últimos caracteres.
 */
export function mascararParaLista(tipo: string, valor: string): string {
  if (valor === '') return '';
  if (tipo === 'cpf' && valor.length === 11) return `***.${valor.slice(3, 6)}.${valor.slice(6, 9)}-**`;
  if (valor.length <= 3) return '*'.repeat(valor.length);
  return '*'.repeat(valor.length - 3) + valor.slice(-3);
}
