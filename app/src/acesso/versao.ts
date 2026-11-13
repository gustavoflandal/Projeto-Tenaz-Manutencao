// Versão do app para o selo e para a tela "Sobre" (RNF-039 e seção 9.5 de
// docs/5_ARQUITETURA_SOFTWARE.md).
//
// O selo mostra a versão que está DE FATO rodando: a do release embutido no
// JavaScript (`extra.versaoRelease`, escrita por app.config.ts a partir da tag
// git). Ela pode ser mais nova que o binário da loja, quando chegou uma
// atualização sem loja; por isso a tela "Sobre" mostra também o binário e o
// identificador da atualização.

import { formatarDataHora, fusoDoAparelho, type Pais } from '@tenaz/dominio';

import type { Traduzir } from '../i18n';

export interface InformacoesDoApp {
  /** Release rodando agora, SemVer completo (1.4.2 ou 1.5.0-rc.1). */
  versaoRelease: string;
  commit: string;
  /** RFC 3339 em UTC; vazio em build local. */
  dataBuild: string;
  /** producao, homologacao, sandbox ou desenvolvimento. */
  ambiente: string;
  /** Binário instalado pela loja. Nulo no Expo Go e na web. */
  binario: { versao: string; build: string } | null;
  /** Identificador da atualização sem loja em execução; nulo quando roda o JavaScript embutido no binário. */
  atualizacaoId: string | null;
}

const AMBIENTES = ['producao', 'homologacao', 'sandbox', 'desenvolvimento'] as const;

/** Nome do ambiente no idioma da tela; ambiente desconhecido aparece como veio. */
export function rotuloDoAmbiente(t: Traduzir, ambiente: string): string {
  return (AMBIENTES as readonly string[]).includes(ambiente) ? t(`ambiente.${ambiente as (typeof AMBIENTES)[number]}`) : ambiente;
}

/** Lê o bloco `extra` do app.config.ts. Campo ausente cai em valor explícito, nunca em versão inventada. */
export function lerInformacoes(
  extra: Record<string, unknown> | undefined,
  nativo: { versao: string | null; build: string | null },
  atualizacaoId: string | null,
): InformacoesDoApp {
  const texto = (v: unknown, padrao: string): string => (typeof v === 'string' && v !== '' ? v : padrao);
  return {
    versaoRelease: texto(extra?.versaoRelease, '0.0.0-dev'),
    commit: texto(extra?.commit, 'desconhecido'),
    dataBuild: texto(extra?.dataBuild, ''),
    ambiente: texto(extra?.ambiente, 'desenvolvimento'),
    binario: nativo.versao ? { versao: nativo.versao, build: nativo.build ?? '?' } : null,
    atualizacaoId,
  };
}

/** Texto do selo: "v1.4.2". */
export function textoDoSelo(t: Traduzir, i: InformacoesDoApp): string {
  return t('selo.texto', { versao: i.versaoRelease });
}

/**
 * Nome acessível do selo (docs/7_PADROES_DESIGN.md, 6.19). Em produção:
 * "Versão 1.4.2. Abrir informações sobre o sistema". Fora dela, o ambiente
 * entra: "Versão 1.5.0-rc.1, ambiente Homologação. Abrir informações...".
 */
export function rotuloAcessivelDoSelo(t: Traduzir, i: InformacoesDoApp): string {
  if (i.ambiente === 'producao') return t('selo.nomeAcessivel', { versao: i.versaoRelease });
  return t('selo.nomeAcessivelComAmbiente', { versao: i.versaoRelease, ambiente: rotuloDoAmbiente(t, i.ambiente) });
}

/** País e fuso para mostrar uma data sem empresa na tela (ADR-013, Parte G): o fuso é o do aparelho, que ali é explícito. */
export interface LocalDeExibicao {
  readonly pais: Pais;
  readonly fuso: string;
}

/**
 * O app ainda não conhece a empresa do usuário no "Sobre", que abre também antes
 * do login: vale o país de origem (BR) e o fuso do aparelho. Quando a sessão do
 * app trouxer o país da empresa, ele entra aqui.
 */
export function localDoAparelho(): LocalDeExibicao {
  return { pais: 'BR', fuso: fusoDoAparelho() };
}

/** "12/11/2026 14:05" (formato do país, 24 h) no fuso do aparelho; "não informada" quando o build não registrou a data. */
export function formatarDataDoBuild(t: Traduzir, iso: string, local: LocalDeExibicao = localDoAparelho()): string {
  if (iso === '') return t('mobile.sobre.dataNaoInformada');
  try {
    return formatarDataHora(iso, local.fuso, local.pais);
  } catch {
    // Data do build ilegível é o mesmo que ausente: nunca mostra texto estranho.
    return t('mobile.sobre.dataNaoInformada');
  }
}

export interface LinhaDoSobre {
  rotulo: string;
  valor: string;
}

/** Estado da consulta a GET /v1/versao na tela Sobre. */
export type VersaoDaApi = { estado: 'consultando' } | { estado: 'indisponivel' } | { estado: 'ok'; versao: string };

/** As linhas da tela Sobre, na ordem do documento 7 (9.11). A mesma lista serve à tela e a "Copiar informações". */
export function linhasDoSobre(
  t: Traduzir,
  i: InformacoesDoApp,
  api: VersaoDaApi,
  local: LocalDeExibicao = localDoAparelho(),
): LinhaDoSobre[] {
  const valorDaApi =
    api.estado === 'ok'
      ? t('sobre.versao', { versao: api.versao })
      : api.estado === 'consultando'
        ? t('sobre.apiConsultando')
        : t('mobile.sobre.apiIndisponivel');
  return [
    { rotulo: t('mobile.sobre.app'), valor: t('sobre.versao', { versao: i.versaoRelease }) },
    {
      rotulo: t('mobile.sobre.binario'),
      valor: i.binario
        ? t('mobile.sobre.binarioValor', { versao: i.binario.versao, build: i.binario.build })
        : t('mobile.sobre.binarioNaoSeAplica'),
    },
    { rotulo: t('mobile.sobre.atualizacao'), valor: i.atualizacaoId ?? t('mobile.sobre.atualizacaoNenhuma') },
    { rotulo: t('sobre.api'), valor: valorDaApi },
    { rotulo: t('sobre.dataBuild'), valor: formatarDataDoBuild(t, i.dataBuild, local) },
    { rotulo: t('sobre.commit'), valor: i.commit },
    { rotulo: t('sobre.ambiente'), valor: rotuloDoAmbiente(t, i.ambiente) },
  ];
}

/** Texto de "Copiar informações" para o usuário colar numa conversa com o suporte. */
export function textoParaCopiar(t: Traduzir, linhas: readonly LinhaDoSobre[]): string {
  return linhas.map((l) => t('sobre.linhaCopia', { rotulo: l.rotulo, valor: l.valor })).join('\n');
}
