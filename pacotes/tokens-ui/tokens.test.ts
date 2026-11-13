import { readFileSync } from 'node:fs';

import { describe, expect, it } from 'vitest';

import { ARQUIVOS, gerar, lerTokens } from './gerar.mjs';

type Tema = Record<string, string>;

const tokens = lerTokens() as { temas: Record<'claro' | 'escuro' | 'sol', Tema> };

describe('arquivos gerados', () => {
  it('tokens.css e tailwind.css estão em dia com tokens.json', () => {
    const gerado = gerar(tokens);
    const ler = (caminho: string) => readFileSync(caminho, 'utf8').replace(/\r\n/g, '\n');
    expect(ler(ARQUIVOS.tokensCss), 'rode: npx pnpm --filter @tenaz/tokens-ui gerar').toBe(gerado.tokensCss);
    expect(ler(ARQUIVOS.tailwindCss), 'rode: npx pnpm --filter @tenaz/tokens-ui gerar').toBe(gerado.tailwindCss);
  });

  it('os três temas definem exatamente os mesmos tokens', () => {
    const chaves = (t: Tema) => Object.keys(t).sort();
    expect(chaves(tokens.temas.escuro)).toEqual(chaves(tokens.temas.claro));
    expect(chaves(tokens.temas.sol)).toEqual(chaves(tokens.temas.claro));
  });
});

// Luminância relativa e razão de contraste da WCAG 2.1, a mesma fórmula da
// seção 4 de docs/7_PADROES_DESIGN.md.
function luminancia(hex: string): number {
  const canal = (i: number) => {
    const c = Number.parseInt(hex.slice(i, i + 2), 16) / 255;
    return c <= 0.03928 ? c / 12.92 : ((c + 0.055) / 1.055) ** 2.4;
  };
  return 0.2126 * canal(1) + 0.7152 * canal(3) + 0.0722 * canal(5);
}

function contraste(a: string, b: string): number {
  const [claro, escuro] = [luminancia(a), luminancia(b)].sort((x, y) => y - x) as [number, number];
  return (claro + 0.05) / (escuro + 0.05);
}

// Pares que a interface usa, com o mínimo do critério: 4,5 para texto
// (1.4.3) e 3 para contorno de componente e foco (1.4.11).
const PARES: ReadonlyArray<readonly [string, string, number]> = [
  ['texto', 'fundo', 4.5],
  ['texto', 'superficie', 4.5],
  ['texto', 'superficie-elevada', 4.5],
  ['texto-secundario', 'fundo', 4.5],
  ['texto-secundario', 'superficie', 4.5],
  ['texto-secundario', 'superficie-elevada', 4.5],
  ['texto-terciario', 'fundo', 4.5],
  ['texto-terciario', 'superficie', 4.5],
  ['acao-primaria-texto', 'acao-primaria', 4.5],
  ['acao-primaria-texto', 'acao-primaria-hover', 4.5],
  ['acao-primaria-texto', 'acao-primaria-pressionada', 4.5],
  ['perigo-texto', 'perigo', 4.5],
  ['nav-texto', 'nav-fundo', 4.5],
  ['nav-texto-secundario', 'nav-fundo', 4.5],
  ['nav-texto', 'nav-item-ativo', 4.5],
  ['toast-texto', 'toast-fundo', 4.5],
  ['foco', 'fundo', 3],
  ['foco', 'superficie', 3],
  ['foco', 'superficie-elevada', 3],
  ['borda-controle', 'superficie', 3],
  ['status-operando-texto', 'status-operando-fundo', 4.5],
  ['status-preventiva-texto', 'status-preventiva-fundo', 4.5],
  ['status-atencao-texto', 'status-atencao-fundo', 4.5],
  ['status-parado-texto', 'status-parado-fundo', 4.5],
  ['status-inativo-texto', 'status-inativo-fundo', 4.5],
];

describe.each(['claro', 'escuro', 'sol'] as const)('contraste WCAG 2.1 AA no tema %s', (nome) => {
  const tema = tokens.temas[nome];
  it.each(PARES)('%s sobre %s atinge %d:1', (frente, fundo, minimo) => {
    const a = tema[frente];
    const b = tema[fundo];
    expect(a, frente).toMatch(/^#[0-9A-F]{6}$/i);
    expect(b, fundo).toMatch(/^#[0-9A-F]{6}$/i);
    expect(contraste(a as string, b as string)).toBeGreaterThanOrEqual(minimo);
  });

  it('texto sobre Sinal é sempre Aço', () => {
    expect(tema['acao-primaria']).toBe('#FF6A13');
    expect(tema['acao-primaria-texto']).toBe('#14202E');
  });
});
