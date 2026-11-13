// Dinheiro (ADR-013, Parte F, documento 5): sem biblioteca, sobre `bigint`.
// Duas espécies, sempre com a moeda ao lado: `Valor` (escala da moeda) e
// `PrecoUnitario` (até seis casas, nunca somado). No contrato, dinheiro é
// `string`; Number(), parseFloat e aritmética de `number` sobre dinheiro são
// proibidos neste pacote. Provado por backend/testes/casos/dinheiro.json.

import { casasDecimaisDaMoeda, MOEDAS } from './paises';

export type CodigoDeMoeda = (typeof MOEDAS)[number];

/** Moeda do catálogo fechado (ISO 4217): código e casas. */
export interface Moeda {
  readonly codigo: CodigoDeMoeda;
  readonly casas: number;
}

export type CodigoDeErroDeDinheiro =
  | 'formato_invalido'
  | 'fora_da_faixa'
  | 'moeda_diferente'
  | 'moeda_desconhecida'
  | 'rateio_invalido'
  | 'quantidade_invalida';

/** Erro da biblioteca. Os internos (moeda_*, rateio_invalido) são defeito, nunca resposta à pessoa. */
export class ErroDeDinheiro extends Error {
  readonly codigo: CodigoDeErroDeDinheiro;
  constructor(codigo: CodigoDeErroDeDinheiro) {
    super(codigo);
    this.name = 'ErroDeDinheiro';
    this.codigo = codigo;
  }
}

/** Constrói a moeda a partir do código; fora do catálogo é `moeda_desconhecida`. */
export function moedaDoCodigo(codigo: string): Moeda {
  const conhecido = MOEDAS.find((m) => m === codigo);
  if (conhecido === undefined) throw new ErroDeDinheiro('moeda_desconhecida');
  return { codigo: conhecido, casas: casasDecimaisDaMoeda(conhecido) };
}

/** Quantia que se paga ou se soma. `unidades` está na unidade menor da moeda. */
export interface Valor {
  readonly especie: 'valor';
  readonly moeda: Moeda;
  readonly unidades: bigint;
}

/** Dinheiro por unidade de uma quantidade. `unidades` está em milionésimos da unidade da moeda. */
export interface PrecoUnitario {
  readonly especie: 'preco_unitario';
  readonly moeda: Moeda;
  readonly unidades: bigint;
}

/** Quantidade como fração exata; o denominador é positivo. */
export interface Quantidade {
  readonly numerador: bigint;
  readonly denominador: bigint;
}

const ZERO = BigInt(0);
const UM = BigInt(1);
const DEZ = BigInt(10);
const CASAS_DO_PRECO = 6;

function potenciaDeDez(n: number): bigint {
  return DEZ ** BigInt(n);
}

/** Maior valor absoluto em unidades: 16 dígitos inteiros no valor, 12 no preço, mais as casas. */
function limite(especie: 'valor' | 'preco_unitario', moeda: Moeda): bigint {
  return potenciaDeDez(especie === 'valor' ? 16 + moeda.casas : 12 + CASAS_DO_PRECO) - UM;
}

function absoluto(n: bigint): bigint {
  return n < ZERO ? -n : n;
}

function dentroDaFaixa(especie: 'valor' | 'preco_unitario', moeda: Moeda, unidades: bigint): boolean {
  return absoluto(unidades) <= limite(especie, moeda);
}

// ------------------------------------------------------------------ leitura

const GRAMATICA = {
  valor: { comSinal: /^(-)?(0|[1-9][0-9]{0,15})(?:\.([0-9]{1,2}))?$/, semSinal: /^()(0|[1-9][0-9]{0,15})(?:\.([0-9]{1,2}))?$/ },
  preco_unitario: {
    comSinal: /^(-)?(0|[1-9][0-9]{0,11})(?:\.([0-9]{1,6}))?$/,
    semSinal: /^()(0|[1-9][0-9]{0,11})(?:\.([0-9]{1,6}))?$/,
  },
} as const;

export interface OpcoesDeLeitura {
  /** Padrão do contrato sem `-?` quando falso (o padrão é falso). */
  readonly aceitaNegativo?: boolean;
}

/** Lê pela gramática exata do padrão do contrato; devolve as unidades na escala `casas`. */
function lerUnidades(
  texto: string,
  especie: 'valor' | 'preco_unitario',
  moeda: Moeda,
  aceitaNegativo: boolean,
): bigint {
  const gramatica = aceitaNegativo ? GRAMATICA[especie].comSinal : GRAMATICA[especie].semSinal;
  const achado = gramatica.exec(texto);
  if (achado === null) throw new ErroDeDinheiro('formato_invalido');
  const sinal = achado[1] ?? '';
  const inteira = achado[2] ?? '';
  const fracao = achado[3] ?? '';
  const casas = especie === 'valor' ? moeda.casas : CASAS_DO_PRECO;
  // Casas além da escala da moeda só valem se forem zeros (valor); nunca arredonda.
  if (fracao.length > casas && /[1-9]/.test(fracao.slice(casas))) throw new ErroDeDinheiro('formato_invalido');
  const unidades = BigInt(inteira + fracao.slice(0, casas).padEnd(casas, '0'));
  return sinal === '-' ? -unidades : unidades;
}

export function lerValor(texto: string, moeda: Moeda, opcoes: OpcoesDeLeitura = {}): Valor {
  return { especie: 'valor', moeda, unidades: lerUnidades(texto, 'valor', moeda, opcoes.aceitaNegativo ?? false) };
}

export function lerPrecoUnitario(texto: string, moeda: Moeda, opcoes: OpcoesDeLeitura = {}): PrecoUnitario {
  return {
    especie: 'preco_unitario',
    moeda,
    unidades: lerUnidades(texto, 'preco_unitario', moeda, opcoes.aceitaNegativo ?? false),
  };
}

// ------------------------------------------------------------ texto canônico

function escrever(unidades: bigint, casas: number): string {
  const negativo = unidades < ZERO;
  const digitos = absoluto(unidades).toString().padStart(casas + 1, '0');
  const inteira = digitos.slice(0, digitos.length - casas);
  const fracao = digitos.slice(digitos.length - casas);
  return `${negativo ? '-' : ''}${inteira}${casas > 0 ? `.${fracao}` : ''}`;
}

/** Valor com exatamente as casas da moeda; nunca `-0`. */
export function textoDoValor(valor: Valor): string {
  return escrever(valor.unidades, valor.moeda.casas);
}

/** Preço com no mínimo as casas da moeda e no máximo seis, sem zeros à direita além das casas da moeda. */
export function textoDoPrecoUnitario(preco: PrecoUnitario): string {
  const completo = escrever(preco.unidades, CASAS_DO_PRECO);
  const minimo = preco.moeda.casas;
  const ponto = completo.indexOf('.');
  let fim = completo.length;
  while (fim - ponto - 1 > minimo && completo.charCodeAt(fim - 1) === 48) fim -= 1;
  const texto = completo.slice(0, fim);
  return minimo === 0 && texto.endsWith('.') ? texto.slice(0, -1) : texto;
}

// ------------------------------------------------------------ soma e subtração

function operar(a: Valor, b: Valor, sinal: bigint): Valor {
  if (a.moeda.codigo !== b.moeda.codigo) throw new ErroDeDinheiro('moeda_diferente');
  const unidades = a.unidades + sinal * b.unidades;
  if (!dentroDaFaixa('valor', a.moeda, unidades)) throw new ErroDeDinheiro('fora_da_faixa');
  return { especie: 'valor', moeda: a.moeda, unidades };
}

/** Só `Valor`, na mesma moeda; o resultado tem de caber na coluna. */
export function somar(a: Valor, b: Valor): Valor {
  return operar(a, b, UM);
}

export function subtrair(a: Valor, b: Valor): Valor {
  return operar(a, b, -UM);
}

// ------------------------------------------------------------- multiplicação

/** Divide arredondando meio para longe do zero; o divisor é positivo. */
function dividirArredondando(dividendo: bigint, divisor: bigint): bigint {
  const modulo = absoluto(dividendo);
  const q = modulo / divisor;
  const resto = modulo % divisor;
  const arredondado = resto * BigInt(2) >= divisor ? q + UM : q;
  return dividendo < ZERO ? -arredondado : arredondado;
}

/** Quantidade a partir de um texto decimal do contrato (`"12.5"`, `"-0.005"`). */
export function quantidadeDeTexto(texto: string): Quantidade {
  const achado = /^(-)?(0|[1-9][0-9]*)(?:\.([0-9]+))?$/.exec(texto);
  if (achado === null) throw new ErroDeDinheiro('quantidade_invalida');
  const fracao = achado[3] ?? '';
  const numerador = BigInt((achado[2] ?? '') + fracao);
  return { numerador: achado[1] === '-' ? -numerador : numerador, denominador: potenciaDeDez(fracao.length) };
}

/**
 * Preço unitário × quantidade, com um só arredondamento no fim (meio para longe
 * do zero), na escala da moeda. A quantidade pode ser negativa (estorno).
 */
export function multiplicar(preco: PrecoUnitario, quantidade: Quantidade): Valor {
  if (quantidade.denominador <= ZERO) throw new ErroDeDinheiro('quantidade_invalida');
  const divisor = quantidade.denominador * potenciaDeDez(CASAS_DO_PRECO - preco.moeda.casas);
  const unidades = dividirArredondando(preco.unidades * quantidade.numerador, divisor);
  if (!dentroDaFaixa('valor', preco.moeda, unidades)) throw new ErroDeDinheiro('fora_da_faixa');
  return { especie: 'valor', moeda: preco.moeda, unidades };
}

// -------------------------------------------------------------------- rateio

/**
 * Reparte o valor em partes proporcionais aos pesos (texto decimal não
 * negativo; lista não vazia, soma positiva). Piso do quinhão sobre o módulo, sobra
 * de unidades menores pelo maior resto (empate para o menor índice), sinal do
 * valor devolvido a todas as partes; a soma das partes é o valor.
 */
export function ratear(valor: Valor, pesos: readonly string[]): Valor[] {
  if (pesos.length === 0) throw new ErroDeDinheiro('rateio_invalido');
  const lidos = pesos.map((p) => {
    const achado = /^(0|[1-9][0-9]*)(?:\.([0-9]+))?$/.exec(p);
    if (achado === null) throw new ErroDeDinheiro('rateio_invalido');
    return { inteira: achado[1] ?? '', fracao: achado[2] ?? '' };
  });
  const casas = Math.max(...lidos.map((p) => p.fracao.length));
  const inteiros = lidos.map((p) => BigInt(p.inteira + p.fracao.padEnd(casas, '0')));
  const soma = inteiros.reduce((acc, p) => acc + p, ZERO);
  if (soma <= ZERO) throw new ErroDeDinheiro('rateio_invalido');

  const modulo = absoluto(valor.unidades);
  const partes = inteiros.map((p) => ({ piso: (modulo * p) / soma, resto: (modulo * p) % soma }));
  let sobra = modulo - partes.reduce((acc, p) => acc + p.piso, ZERO);
  const ordem = partes
    .map((p, indice) => ({ indice, resto: p.resto }))
    .sort((x, y) => (x.resto === y.resto ? x.indice - y.indice : x.resto > y.resto ? -1 : 1));
  const extras = new Set<number>();
  for (const { indice } of ordem) {
    if (sobra <= ZERO) break;
    extras.add(indice);
    sobra -= UM;
  }
  return partes.map((p, indice) => {
    const unidades = p.piso + (extras.has(indice) ? UM : ZERO);
    return { especie: 'valor', moeda: valor.moeda, unidades: valor.unidades < ZERO ? -unidades : unidades };
  });
}
