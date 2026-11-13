// Tema da interface: claro, escuro, sol ou igual ao sistema.
//
// Os valores de cor vêm de @tenaz/tokens-ui; aqui só se escolhe qual bloco
// vale, pelo atributo `data-tema` no <html> (seção 3.3 do documento 7). Sem o
// atributo, o CSS segue `prefers-color-scheme`.
//
// A preferência é do navegador (localStorage), não dado do servidor. Quando o
// menu do usuário existir, ela passa a acompanhar o perfil.

import { useSyncExternalStore } from 'react';

export const PREFERENCIAS_DE_TEMA = ['sistema', 'claro', 'escuro', 'sol'] as const;
export type PreferenciaTema = (typeof PREFERENCIAS_DE_TEMA)[number];
export type TemaResolvido = Exclude<PreferenciaTema, 'sistema'>;

const CHAVE = 'tenaz.tema';
const CONSULTA_ESCURO = '(prefers-color-scheme: dark)';

function ehPreferencia(valor: unknown): valor is PreferenciaTema {
  return typeof valor === 'string' && (PREFERENCIAS_DE_TEMA as readonly string[]).includes(valor);
}

export function lerPreferenciaSalva(): PreferenciaTema {
  try {
    const salvo = globalThis.localStorage.getItem(CHAVE);
    return ehPreferencia(salvo) ? salvo : 'sistema';
  } catch {
    // Armazenamento bloqueado pelo navegador: segue o sistema.
    return 'sistema';
  }
}

function consultaEscuro(): MediaQueryList | undefined {
  return typeof globalThis.matchMedia === 'function' ? globalThis.matchMedia(CONSULTA_ESCURO) : undefined;
}

export function resolverTema(preferencia: PreferenciaTema, sistemaEscuro: boolean): TemaResolvido {
  if (preferencia !== 'sistema') return preferencia;
  return sistemaEscuro ? 'escuro' : 'claro';
}

export function aplicarTema(preferencia: PreferenciaTema, raiz: HTMLElement = document.documentElement): void {
  if (preferencia === 'sistema') raiz.removeAttribute('data-tema');
  else raiz.setAttribute('data-tema', preferencia);
}

// Estado mínimo de preferência, fora do React, para que o tema seja aplicado
// antes da primeira renderização e todos os componentes vejam o mesmo valor.
let preferenciaAtual: PreferenciaTema | undefined;
const ouvintes = new Set<() => void>();

function preferencia(): PreferenciaTema {
  preferenciaAtual ??= lerPreferenciaSalva();
  return preferenciaAtual;
}

function avisar(): void {
  for (const ouvinte of ouvintes) ouvinte();
}

/** Aplica a preferência salva. Chamado uma vez, antes de montar a aplicação. */
export function iniciarTema(): void {
  preferenciaAtual = lerPreferenciaSalva();
  aplicarTema(preferenciaAtual);
}

export function definirPreferenciaTema(nova: PreferenciaTema): void {
  preferenciaAtual = nova;
  try {
    if (nova === 'sistema') globalThis.localStorage.removeItem(CHAVE);
    else globalThis.localStorage.setItem(CHAVE, nova);
  } catch {
    // Sem armazenamento, a escolha vale só até recarregar a página.
  }
  aplicarTema(nova);
  avisar();
}

function assinarPreferencia(ouvinte: () => void): () => void {
  ouvintes.add(ouvinte);
  return () => ouvintes.delete(ouvinte);
}

function assinarSistema(ouvinte: () => void): () => void {
  const consulta = consultaEscuro();
  consulta?.addEventListener('change', ouvinte);
  return () => consulta?.removeEventListener('change', ouvinte);
}

function sistemaEscuro(): boolean {
  return consultaEscuro()?.matches ?? false;
}

export interface EstadoTema {
  readonly preferencia: PreferenciaTema;
  readonly resolvido: TemaResolvido;
  readonly definir: (preferencia: PreferenciaTema) => void;
}

export function useTema(): EstadoTema {
  const atual = useSyncExternalStore(assinarPreferencia, preferencia, () => 'sistema' as const);
  const escuro = useSyncExternalStore(assinarSistema, sistemaEscuro, () => false);
  return { preferencia: atual, resolvido: resolverTema(atual, escuro), definir: definirPreferenciaTema };
}
