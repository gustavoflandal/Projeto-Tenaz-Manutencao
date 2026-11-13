// Textos do app: pt-BR e es-419, do pacote @tenaz/i18n (docs/7_PADROES_DESIGN.md,
// 10.7). Nenhum texto de interface fica escrito na tela.
//
// Sem react-i18next: o app só precisa de busca por chave, interpolação
// `{{nome}}` e plural. A lógica é pura para rodar nos testes, no Node.

import { ehIdioma, ESPACO_DE_NOMES, IDIOMA_PADRAO, type Idioma, recursos, type Textos } from '@tenaz/i18n';

type Caminhos<T, P extends string = ''> = {
  [K in keyof T & string]: T[K] extends string ? `${P}${K}` : Caminhos<T[K], `${P}${K}.`>;
}[keyof T & string];

/** Chave de texto do catálogo, no formato `mobile.login.entrar`. Chave que não existe é erro de tipo. */
type SemPlural<S extends string> = S extends `${infer Base}_${'zero' | 'one' | 'two' | 'few' | 'many' | 'other'}` ? Base : never;

/** Também aceita a chave-base de um plural (`formulario.resumoTitulo`), que a tradução resolve por `count`. */
export type ChaveDeTexto = Caminhos<Textos> | SemPlural<Caminhos<Textos>>;

export type Parametros = Record<string, string | number>;

/** Função de tradução já presa a um idioma. */
export type Traduzir = (chave: ChaveDeTexto, parametros?: Parametros) => string;

function buscar(idioma: Idioma, chave: string): string | undefined {
  let no: unknown = recursos[idioma][ESPACO_DE_NOMES];
  for (const parte of chave.split('.')) {
    if (typeof no !== 'object' || no === null) return undefined;
    no = (no as Record<string, unknown>)[parte];
  }
  return typeof no === 'string' ? no : undefined;
}

/**
 * Regra de reserva (CLDR) para aparelho sem Intl.PluralRules. Tem de dar a
 * mesma categoria que o Intl; o teste compara as duas.
 *  - pt-BR: `one` para 0 e 1 (inteiro 0..1, inclui 1,5); `many` para milhões exatos.
 *  - es-419: `one` só para 1; `many` para milhões exatos; o resto é `other`.
 */
export function categoriaDeReserva(idioma: Idioma, n: number): string {
  const x = Math.abs(n);
  const milhoes = Number.isInteger(x) && x !== 0 && x % 1_000_000 === 0;
  if (idioma === 'pt-BR' && x < 2) return 'one';
  if (idioma === 'es-419' && x === 1) return 'one';
  return milhoes ? 'many' : 'other';
}

function categoriaDePlural(idioma: Idioma, n: number): string {
  try {
    return new Intl.PluralRules(idioma).select(n);
  } catch {
    return categoriaDeReserva(idioma, n);
  }
}

export function traduzir(idioma: Idioma, chave: ChaveDeTexto, parametros?: Parametros): string {
  let texto: string | undefined;
  if (parametros && typeof parametros.count === 'number') {
    texto =
      buscar(idioma, `${chave}_${categoriaDePlural(idioma, parametros.count)}`) ?? buscar(idioma, `${chave}_other`);
  }
  // Se faltar no idioma pedido, cai no pt-BR e por último na própria chave:
  // texto de outro idioma é melhor que tela em branco, e o teste de paridade
  // do pacote impede que isso aconteça de verdade.
  texto ??= buscar(idioma, chave) ?? buscar(IDIOMA_PADRAO, chave) ?? chave;
  return texto.replace(/\{\{\s*([^}\s]+)\s*\}\}/g, (_m, nome: string) => String(parametros?.[nome] ?? ''));
}

export function tradutor(idioma: Idioma): Traduzir {
  return (chave, parametros) => traduzir(idioma, chave, parametros);
}

/** Idioma do aparelho, na ordem de preferência do sistema: pt* vira pt-BR, es* vira es-419, o resto cai no padrão. */
export function idiomaDoAparelho(tags: readonly string[]): Idioma {
  for (const tag of tags) {
    if (ehIdioma(tag)) return tag;
    const lingua = tag.toLowerCase().split(/[-_]/)[0];
    if (lingua === 'pt') return 'pt-BR';
    if (lingua === 'es') return 'es-419';
  }
  return IDIOMA_PADRAO;
}

/** A escolha da pessoa vale mais que o aparelho; escolha guardada que não existe mais é ignorada. */
export function escolherIdioma(guardado: string | null, tagsDoAparelho: readonly string[]): Idioma {
  return ehIdioma(guardado) ? guardado : idiomaDoAparelho(tagsDoAparelho);
}

export { IDIOMAS, type Idioma } from '@tenaz/i18n';
