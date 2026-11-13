// Prova que a trava do ESLint acusa (ADR-013, Parte G, item 4): Intl gerando
// texto e toLocale*String são barrados na web, no app e nos pacotes; o
// Intl.DateTimeFormat só passa em formatos.ts, como motor de fuso, sem `format`.

import { resolve } from 'node:path';

import { ESLint } from 'eslint';
import { describe, expect, it } from 'vitest';

const RAIZ = resolve(import.meta.dirname, '../../..');

const daWeb = new ESLint({ cwd: resolve(RAIZ, 'web') });
const dosPacotes = new ESLint({ cwd: RAIZ, overrideConfigFile: resolve(RAIZ, 'web/eslint.pacotes.config.js') });

async function acusa(eslint: ESLint, arquivo: string, codigo: string): Promise<string[]> {
  const [resultado] = await eslint.lintText(codigo, { filePath: resolve(eslint === daWeb ? resolve(RAIZ, 'web') : RAIZ, arquivo) });
  return (resultado?.messages ?? []).filter((m) => m.ruleId === 'no-restricted-syntax').map((m) => m.message);
}

const PROIBIDOS: ReadonlyArray<readonly [string, string]> = [
  ['Intl.NumberFormat', "const f = new Intl.NumberFormat('pt-BR').format(1);"],
  ['Intl.NumberFormat sem new', "const f = Intl.NumberFormat('pt-BR');"],
  ['Intl.RelativeTimeFormat', "const f = new Intl.RelativeTimeFormat('pt-BR');"],
  ['Intl.ListFormat', "const f = new Intl.ListFormat('pt-BR');"],
  ['Intl.DateTimeFormat', "const f = new Intl.DateTimeFormat('pt-BR').format(new Date());"],
  ['globalThis.Intl.NumberFormat', "const f = new globalThis.Intl.NumberFormat('pt-BR');"],
  ['toLocaleString', 'const t = (1).toLocaleString();'],
  ['toLocaleDateString', 'const t = new Date().toLocaleDateString("pt-BR");'],
  ['toLocaleTimeString', 'const t = new Date().toLocaleTimeString("pt-BR");'],
];

const PERMITIDOS: ReadonlyArray<readonly [string, string]> = [
  ['Intl.Collator', "const c = new Intl.Collator('pt-BR').compare('a', 'b');"],
  ['Intl.PluralRules', "const p = new Intl.PluralRules('pt-BR').select(1);"],
  ['tipo Intl.NumberFormatOptions', 'let o: Intl.NumberFormatOptions | undefined;\nexport { o };'],
];

describe.each([
  ['web', daWeb, 'src/modulos/qualquer.ts'],
  ['app', dosPacotes, 'app/src/qualquer.ts'],
  ['pacotes', dosPacotes, 'pacotes/dominio/src/qualquer.ts'],
] as const)('trava de formatação: %s', (_nome, eslint, arquivo) => {
  it.each(PROIBIDOS)('acusa %s', async (_caso, codigo) => {
    expect((await acusa(eslint, arquivo, codigo)).length).toBeGreaterThan(0);
  });

  it.each(PERMITIDOS)('deixa passar %s', async (_caso, codigo) => {
    expect(await acusa(eslint, arquivo, codigo)).toEqual([]);
  });
});

describe('trava de formatação: o motor de fuso de formatos.ts', () => {
  const formatos = 'pacotes/dominio/src/formatos.ts';

  it('aceita o Intl.DateTimeFormat com formatToParts', async () => {
    expect(await acusa(dosPacotes, formatos, "const f = new Intl.DateTimeFormat('en-US'); f.formatToParts(new Date());")).toEqual([]);
  });

  it('acusa format e formatRange do DateTimeFormat', async () => {
    expect((await acusa(dosPacotes, formatos, "new Intl.DateTimeFormat('en-US').format(new Date());")).length).toBeGreaterThan(0);
    expect((await acusa(dosPacotes, formatos, "new Intl.DateTimeFormat('en-US').formatRange(1, 2);")).length).toBeGreaterThan(0);
  });

  it('continua acusando o NumberFormat e o toLocaleString em formatos.ts', async () => {
    expect((await acusa(dosPacotes, formatos, "new Intl.NumberFormat('pt-BR');")).length).toBeGreaterThan(0);
    expect((await acusa(dosPacotes, formatos, '(1).toLocaleString();')).length).toBeGreaterThan(0);
  });

  it('o Intl.DateTimeFormat acusa em qualquer outro arquivo dos pacotes', async () => {
    expect((await acusa(dosPacotes, 'pacotes/dominio/src/outro.ts', "new Intl.DateTimeFormat('en-US');")).length).toBeGreaterThan(0);
  });
});
