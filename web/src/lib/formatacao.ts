// Formatação de valores para exibição (seção 3.11 de docs/7_PADROES_DESIGN.md,
// ADR-013, Parte G). Tudo passa pelas funções de `@tenaz/dominio` (formatos.ts),
// sem Intl e sem `number` para dinheiro: o formato é do PAÍS DA EMPRESA, lido da
// sessão; o idioma do usuário não muda nada. O ESLint barra Intl e toLocale*.

import {
  ErroDeFormato,
  formatarData as formatarDataDoPais,
  formatarDataHora as formatarDataHoraDoPais,
  formatarPrecoUnitario as formatarPrecoDoPais,
  PERFIS_DE_PAIS,
  simboloDoPais,
  type Pais,
} from '@tenaz/dominio';

/**
 * Instante (RFC 3339) como dd/mm/aaaa hh:mm (24 h) num fuso explícito: o da
 * unidade ou, sem unidade, o da empresa. Devolve null quando o valor não é um
 * instante válido (o chamador mostra o texto como veio).
 */
export function formatarDataHora(instante: string, fuso: string, pais: Pais): string | null {
  try {
    return formatarDataHoraDoPais(instante, fuso, pais);
  } catch (erro) {
    if (erro instanceof ErroDeFormato) return null;
    throw erro;
  }
}

/** Data civil `AAAA-MM-DD` como dd/mm/aaaa, sem fuso; texto que não é data volta como veio. */
export function formatarData(iso: string, pais: Pais): string {
  try {
    return formatarDataDoPais(iso, pais);
  } catch (erro) {
    if (erro instanceof ErroDeFormato) return iso;
    throw erro;
  }
}

/**
 * Preço unitário ou custo por hora, recebido da API (texto canônico), na moeda
 * do país da empresa: de 2 a 6 casas (0 a 6 no guarani), sem passar por `number`.
 */
export function formatarPrecoUnitario(texto: string, pais: Pais): string {
  return formatarPrecoDoPais(texto, PERFIS_DE_PAIS[pais].moeda, pais);
}

/** Símbolo da moeda da empresa (R$, $, Gs.), para rótulos de campo. */
export function simboloDaMoeda(pais: Pais): string {
  return simboloDoPais(pais);
}

/**
 * Data digitada como dd/mm/aaaa para `AAAA-MM-DD`. Devolve null se o texto não
 * for uma data existente (31/02/2026, por exemplo). A digitação é com máscara;
 * nunca só calendário (seção 6.2 do documento 7).
 */
export function dataParaISO(texto: string): string | null {
  const encontrado = /^(\d{2})\/(\d{2})\/(\d{4})$/.exec(texto.trim());
  if (!encontrado) return null;
  const [, dia, mes, ano] = encontrado;
  const data = new Date(Date.UTC(Number(ano), Number(mes) - 1, Number(dia)));
  const confere =
    data.getUTCFullYear() === Number(ano) && data.getUTCMonth() === Number(mes) - 1 && data.getUTCDate() === Number(dia);
  return confere ? `${ano}-${mes}-${dia}` : null;
}

/** Aplica a máscara dd/mm/aaaa enquanto a pessoa digita. */
export function mascararData(texto: string): string {
  const d = texto.replaceAll(/\D/g, '').slice(0, 8);
  if (d.length <= 2) return d;
  if (d.length <= 4) return `${d.slice(0, 2)}/${d.slice(2)}`;
  return `${d.slice(0, 2)}/${d.slice(2, 4)}/${d.slice(4)}`;
}

/**
 * Junta itens numa lista de leitura ("Contagem, Minas Gerais"), com vírgula:
 * nada de `Intl.ListFormat`, que une com espaço ou com "e".
 */
export function listar(itens: ReadonlyArray<string | null | undefined>): string {
  return itens.filter((item): item is string => Boolean(item)).join(', ');
}
