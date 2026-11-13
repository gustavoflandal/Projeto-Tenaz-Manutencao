// Perfil de país (ADR-013): o que muda de um país para outro e é regra do
// núcleo. A mesma tabela de backend/internal/plataforma/pais, provada pelos
// casos compartilhados (backend/testes/casos/identificacao.json, "perfis").

import { SUBDIVISOES } from './enderecos';
import { TIPOS_DE_PESSOA, type TipoDeEmpresa, type TipoDePessoa } from './identificacao';

export const PAISES = ['BR', 'AR', 'PY'] as const;
export type Pais = (typeof PAISES)[number];

export const MOEDAS = ['BRL', 'ARS', 'PYG'] as const;
export type Moeda = (typeof MOEDAS)[number];

export type IdiomaDoPerfil = 'pt-BR' | 'es-419';

export interface PerfilDePais {
  readonly pais: Pais;
  readonly moeda: Moeda;
  readonly idioma_padrao: IdiomaDoPerfil;
  readonly fuso_padrao: string;
  /** Fusos IANA aceitos pelo servidor para o país (espelho de identificacao.json). */
  readonly fusos: readonly string[];
  /** Tipo de identificação de pessoa jurídica: o único que a empresa e a unidade aceitam. */
  readonly tipo_empresa: TipoDeEmpresa;
  /** Tipos de pessoa física do país; a pessoa aceita os de qualquer país. */
  readonly tipos_pessoa: readonly TipoDePessoa[];
  /** Subdivisões do país (ISO 3166-2 completo), na ordem do catálogo. */
  readonly subdivisoes: readonly string[];
}

const FUSOS: Readonly<Record<Pais, readonly string[]>> = {
  BR: [
    'America/Noronha',
    'America/Belem',
    'America/Fortaleza',
    'America/Recife',
    'America/Araguaina',
    'America/Maceio',
    'America/Bahia',
    'America/Sao_Paulo',
    'America/Campo_Grande',
    'America/Cuiaba',
    'America/Santarem',
    'America/Porto_Velho',
    'America/Boa_Vista',
    'America/Manaus',
    'America/Eirunepe',
    'America/Rio_Branco',
  ],
  AR: [
    'America/Argentina/Buenos_Aires',
    'America/Argentina/Cordoba',
    'America/Argentina/Salta',
    'America/Argentina/Jujuy',
    'America/Argentina/Tucuman',
    'America/Argentina/Catamarca',
    'America/Argentina/La_Rioja',
    'America/Argentina/San_Juan',
    'America/Argentina/Mendoza',
    'America/Argentina/San_Luis',
    'America/Argentina/Rio_Gallegos',
    'America/Argentina/Ushuaia',
  ],
  PY: [
    'America/Asuncion',
  ],
};

export const PERFIS_DE_PAIS: Readonly<Record<Pais, PerfilDePais>> = {
  BR: { pais: 'BR', moeda: 'BRL', idioma_padrao: 'pt-BR', fuso_padrao: 'America/Sao_Paulo', fusos: FUSOS.BR, tipo_empresa: 'cnpj', tipos_pessoa: ['cpf'], subdivisoes: SUBDIVISOES.BR },
  AR: {
    pais: 'AR',
    moeda: 'ARS',
    idioma_padrao: 'es-419',
    fuso_padrao: 'America/Argentina/Buenos_Aires', fusos: FUSOS.AR,
    tipo_empresa: 'cuit',
    tipos_pessoa: ['cuil', 'dni'],
    subdivisoes: SUBDIVISOES.AR,
  },
  PY: { pais: 'PY', moeda: 'PYG', idioma_padrao: 'es-419', fuso_padrao: 'America/Asuncion', fusos: FUSOS.PY, tipo_empresa: 'ruc', tipos_pessoa: ['ci'], subdivisoes: SUBDIVISOES.PY },
};

export function ehPais(valor: unknown): valor is Pais {
  return typeof valor === 'string' && (PAISES as readonly string[]).includes(valor);
}

/** País de uma moeda (uma moeda por país): para quem só sabe a moeda da empresa. */
export function paisDaMoeda(moeda: string): Pais | null {
  return PAISES.find((p) => PERFIS_DE_PAIS[p].moeda === moeda) ?? null;
}

/**
 * Tipos de pessoa em ordem de oferta: o padrão do país da empresa primeiro
 * (BR cpf; AR dni; PY ci), depois os outros do país e os dos demais países.
 */
export function tiposDePessoaParaOPais(pais: Pais | null): readonly TipoDePessoa[] {
  const padrao: Record<Pais, TipoDePessoa> = { BR: 'cpf', AR: 'dni', PY: 'ci' };
  const primeiro = pais === null ? 'cpf' : padrao[pais];
  return [primeiro, ...TIPOS_DE_PESSOA.filter((t) => t !== primeiro)];
}

/** Casas decimais da moeda: o guarani paraguaio não tem centavos. */
export function casasDecimaisDaMoeda(moeda: Moeda): number {
  return moeda === 'PYG' ? 0 : 2;
}
