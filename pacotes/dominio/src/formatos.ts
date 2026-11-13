// Formatação para pessoas (ADR-013, Parte G, documento 5; documento 7, seção
// 3.11): funções próprias, sem Intl gerando texto, sobre texto canônico ou
// `bigint`. O formato é do PAÍS DA EMPRESA; o idioma do usuário só escolhe as
// palavras e não muda nada aqui. Dinheiro e decimal nunca passam por `number`.
// Provado por backend/testes/casos/formatos.json.
//
// Única exceção ao Intl: o motor de fuso (`partesNoFuso`), com localidade fixa
// e só as partes numéricas de `formatToParts`; o texto é montado aqui.

import { lerPrecoUnitario, lerValor, moedaDoCodigo, textoDoPrecoUnitario, textoDoValor } from './dinheiro';
import { PERFIS_DE_PAIS, type Moeda, type Pais } from './paises';

/** Entrada fora do texto canônico, `-0`, data inexistente ou fuso desconhecido: defeito do chamador. */
export class ErroDeFormato extends Error {
  constructor(motivo: string) {
    super(motivo);
    this.name = 'ErroDeFormato';
  }
}

/** Espaço inseparável (U+00A0), escrito por código para não depender do editor. */
const NBSP = String.fromCodePoint(0xa0);

/**
 * Convenções por país. Os três usam hoje as mesmas (ponto de milhar, vírgula
 * decimal), mas a tabela fica por país para o dia em que um deles divergir.
 */
interface ConvencaoDoPais {
  readonly milhar: string;
  readonly decimal: string;
  /** Separador da data dd/mm/aaaa. */
  readonly dataSeparador: string;
  /** Símbolo da moeda do país: o `₲` (U+20B2) está descartado (documento 7, 3.11). */
  readonly simbolo: string;
}

const CONVENCOES: Readonly<Record<Pais, ConvencaoDoPais>> = {
  BR: { milhar: '.', decimal: ',', dataSeparador: '/', simbolo: 'R$' },
  AR: { milhar: '.', decimal: ',', dataSeparador: '/', simbolo: '$' },
  PY: { milhar: '.', decimal: ',', dataSeparador: '/', simbolo: 'Gs.' },
};

/** Símbolo da moeda do país da empresa (R$, $, Gs.), para rótulos de campo e para o dinheiro. */
export function simboloDoPais(pais: Pais): string {
  return CONVENCOES[pais].simbolo;
}

// ------------------------------------------------------------------- números

const INTEIRO = /^(-)?(0|[1-9][0-9]*)$/;
const DECIMAL = /^(-)?(0|[1-9][0-9]*)(?:\.([0-9]+))?$/;

function agrupar(digitos: string, separador: string): string {
  let saida = '';
  for (let i = 0; i < digitos.length; i += 1) {
    if (i > 0 && (digitos.length - i) % 3 === 0) saida += separador;
    saida += digitos.charAt(i);
  }
  return saida;
}

/** Sinal, parte inteira agrupada e parte decimal como veio, a partir de um texto canônico. */
function partesDoNumero(sinal: string | undefined, inteira: string, fracao: string | undefined, pais: Pais): { negativo: boolean; texto: string } {
  const c = CONVENCOES[pais];
  const negativo = sinal === '-';
  if (negativo && /^0*$/.test(inteira + (fracao ?? ''))) throw new ErroDeFormato('menos_zero');
  const decimais = fracao === undefined ? '' : c.decimal + fracao;
  return { negativo, texto: agrupar(inteira, c.milhar) + decimais };
}

function comSinal(negativo: boolean, texto: string): string {
  return negativo ? `-${texto}` : texto;
}

/**
 * Inteiro: milhar com ponto a partir de 4 dígitos. Aceita texto canônico ou
 * `bigint`; `number` só para contagem pequena (inteiro seguro).
 */
export function formatarInteiro(valor: string | bigint | number, pais: Pais): string {
  let texto: string;
  if (typeof valor === 'number') {
    if (!Number.isSafeInteger(valor) || Object.is(valor, -0)) throw new ErroDeFormato('inteiro_invalido');
    texto = String(valor);
  } else {
    texto = String(valor);
  }
  const achado = INTEIRO.exec(texto);
  if (achado === null) throw new ErroDeFormato('inteiro_invalido');
  const { negativo, texto: miolo } = partesDoNumero(achado[1], achado[2] ?? '', undefined, pais);
  return comSinal(negativo, miolo);
}

/**
 * Decimal do texto canônico da API: vírgula decimal, milhar com ponto, e as
 * casas exatamente como vieram (não completa, não corta, não arredonda).
 */
export function formatarDecimal(texto: string, pais: Pais): string {
  const achado = DECIMAL.exec(texto);
  if (achado === null) throw new ErroDeFormato('decimal_invalido');
  const { negativo, texto: miolo } = partesDoNumero(achado[1], achado[2] ?? '', achado[3], pais);
  return comSinal(negativo, miolo);
}

/** Percentual em pontos percentuais ("97.4" vira "97,4 %"), com U+00A0 antes do `%`. */
export function formatarPercentual(texto: string, pais: Pais): string {
  return `${formatarDecimal(texto, pais)}${NBSP}%`;
}

/** Medidor: valor, U+00A0 e a unidade (km, h, ciclos, °C, un.). */
export function formatarMedidor(texto: string, unidade: string, pais: Pais): string {
  if (unidade.trim() === '') throw new ErroDeFormato('unidade_vazia');
  return `${formatarDecimal(texto, pais)}${NBSP}${unidade}`;
}

// -------------------------------------------------------------------- dinheiro

/**
 * Símbolo antes do número, separado por U+00A0, e o sinal negativo antes do
 * símbolo ("-R$ 10,00"). Moeda que não é a do país da empresa leva o código ISO
 * no lugar do símbolo (Parte A, item 7); não ocorre com uma moeda por empresa.
 */
function montarDinheiro(canonico: string, moeda: Moeda, pais: Pais): string {
  const negativo = canonico.startsWith('-');
  const sem = negativo ? canonico.slice(1) : canonico;
  const ponto = sem.indexOf('.');
  const c = CONVENCOES[pais];
  const inteira = ponto === -1 ? sem : sem.slice(0, ponto);
  const fracao = ponto === -1 ? '' : c.decimal + sem.slice(ponto + 1);
  const rotulo = PERFIS_DE_PAIS[pais].moeda === moeda ? c.simbolo : moeda;
  return `${negativo ? '-' : ''}${rotulo}${NBSP}${agrupar(inteira, c.milhar)}${fracao}`;
}

/** Valor em dinheiro, com as casas da moeda (2 em BRL e ARS, 0 em PYG), a partir do texto canônico. */
export function formatarValor(texto: string, moeda: Moeda, pais: Pais): string {
  const lido = lerValor(texto, moedaDoCodigo(moeda), { aceitaNegativo: true });
  if (texto.startsWith('-') && lido.unidades === BigInt(0)) throw new ErroDeFormato('menos_zero');
  return montarDinheiro(textoDoValor(lido), moeda, pais);
}

/** Preço unitário e custo por hora: de 2 a 6 casas (0 a 6 no guarani), sem arredondar o que veio. */
export function formatarPrecoUnitario(texto: string, moeda: Moeda, pais: Pais): string {
  const lido = lerPrecoUnitario(texto, moedaDoCodigo(moeda), { aceitaNegativo: true });
  if (texto.startsWith('-') && lido.unidades === BigInt(0)) throw new ErroDeFormato('menos_zero');
  return montarDinheiro(textoDoPrecoUnitario(lido), moeda, pais);
}

// ----------------------------------------------------------------------- datas

function bissexto(ano: number): boolean {
  return (ano % 4 === 0 && ano % 100 !== 0) || ano % 400 === 0;
}

function diasDoMes(ano: number, mes: number): number {
  if (mes === 2) return bissexto(ano) ? 29 : 28;
  return [4, 6, 9, 11].includes(mes) ? 30 : 31;
}

function existe(ano: number, mes: number, dia: number): boolean {
  return ano >= 1 && mes >= 1 && mes <= 12 && dia >= 1 && dia <= diasDoMes(ano, mes);
}

const DATA_CIVIL = /^([0-9]{4})-([0-9]{2})-([0-9]{2})$/;

function dataMontada(ano: string, mes: string, dia: string, pais: Pais): string {
  const s = CONVENCOES[pais].dataSeparador;
  return `${dia}${s}${mes}${s}${ano}`;
}

/** Data civil `AAAA-MM-DD` como dd/mm/aaaa, sem fuso: a data não é um instante. */
export function formatarData(iso: string, pais: Pais): string {
  const achado = DATA_CIVIL.exec(iso);
  if (achado === null) throw new ErroDeFormato('data_invalida');
  const [, ano = '', mes = '', dia = ''] = achado;
  if (!existe(Number(ano), Number(mes), Number(dia))) throw new ErroDeFormato('data_inexistente');
  return dataMontada(ano, mes, dia, pais);
}

const INSTANTE = /^([0-9]{4})-([0-9]{2})-([0-9]{2})T([0-9]{2}):([0-9]{2})(?::([0-9]{2})(?:\.[0-9]+)?)?(Z|[+-][0-9]{2}:[0-9]{2})$/;

const formatadoresDeFuso = new Map<string, Intl.DateTimeFormat>();

/**
 * Motor de fuso: converte o instante em campos civis no fuso IANA dado. É o
 * único uso do `Intl.DateTimeFormat`, com localidade fixa e partes numéricas.
 */
function partesNoFuso(instante: Date, fuso: string): { ano: string; mes: string; dia: string; hora: string; minuto: string } {
  let formatador = formatadoresDeFuso.get(fuso);
  if (formatador === undefined) {
    try {
      formatador = new Intl.DateTimeFormat('en-US', {
        timeZone: fuso,
        calendar: 'gregory',
        numberingSystem: 'latn',
        hourCycle: 'h23',
        year: 'numeric',
        month: '2-digit',
        day: '2-digit',
        hour: '2-digit',
        minute: '2-digit',
      });
    } catch {
      throw new ErroDeFormato('fuso_desconhecido');
    }
    formatadoresDeFuso.set(fuso, formatador);
  }
  const campos = new Map<string, string>();
  for (const parte of formatador.formatToParts(instante)) campos.set(parte.type, parte.value);
  const campo = (tipo: string): string => {
    const valor = campos.get(tipo);
    if (valor === undefined) throw new ErroDeFormato('fuso_sem_partes');
    return valor;
  };
  const hora = campo('hour');
  return {
    ano: campo('year').padStart(4, '0'),
    mes: campo('month').padStart(2, '0'),
    dia: campo('day').padStart(2, '0'),
    // Alguns motores devolvem "24" à meia-noite; o formato é 00:00.
    hora: hora === '24' ? '00' : hora.padStart(2, '0'),
    minuto: campo('minute').padStart(2, '0'),
  };
}

/**
 * Instante (RFC 3339) como dd/mm/aaaa hh:mm, 24 h, no fuso IANA dado (o da
 * unidade ou, sem unidade, o da empresa; nunca o implícito do navegador).
 * Segundos são cortados, nunca arredondados; um espaço comum entre data e hora.
 */
export function formatarDataHora(instante: string, fuso: string, pais: Pais): string {
  const achado = INSTANTE.exec(instante);
  if (achado === null) throw new ErroDeFormato('instante_invalido');
  const [, ano = '', mes = '', dia = '', hora = '', minuto = '', segundo = '00'] = achado;
  if (!existe(Number(ano), Number(mes), Number(dia)) || Number(hora) > 23 || Number(minuto) > 59 || Number(segundo) > 59) {
    throw new ErroDeFormato('instante_inexistente');
  }
  if (fuso.trim() === '') throw new ErroDeFormato('fuso_desconhecido');
  const data = new Date(instante);
  if (Number.isNaN(data.getTime())) throw new ErroDeFormato('instante_invalido');
  const p = partesNoFuso(data, fuso);
  return `${dataMontada(p.ano, p.mes, p.dia, pais)} ${p.hora}:${p.minuto}`;
}

/**
 * Fuso IANA do aparelho ou do navegador, para telas sem empresa (o "Sobre" do
 * app e da web). Ali o fuso é explícito: foi o aparelho quem o declarou.
 */
export function fusoDoAparelho(): string {
  const fuso = new Intl.DateTimeFormat('en-US').resolvedOptions().timeZone;
  return fuso === '' ? 'UTC' : fuso;
}
