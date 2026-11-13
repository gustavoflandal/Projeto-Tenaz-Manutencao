// Guarda do portão da web: nenhum valor de cor, espaço ou tipo escrito à mão
// nos componentes (seção 3 do documento 7). Os valores vêm só dos tokens de
// @tenaz/tokens-ui, pelas classes geradas em tailwind.css.

import { readdirSync, readFileSync } from 'node:fs';
import { join, relative, resolve } from 'node:path';

import { describe, expect, it } from 'vitest';

const SRC = resolve(__dirname, '..');

function arquivosDeCodigo(pasta: string): string[] {
  return readdirSync(pasta, { withFileTypes: true }).flatMap((item) => {
    const caminho = join(pasta, item.name);
    if (item.isDirectory()) return item.name === 'testes' ? [] : arquivosDeCodigo(caminho);
    return /\.(ts|tsx|css)$/.test(item.name) && !/\.test\.tsx?$/.test(item.name) ? [caminho] : [];
  });
}

const PROIBIDOS: ReadonlyArray<readonly [string, RegExp]> = [
  ['cor hexadecimal', /#[0-9a-fA-F]{3,8}\b(?![-\w])/],
  ['cor rgb ou hsl', /\b(rgba?|hsla?)\(/],
  ['valor arbitrário do Tailwind', /\b[a-z][a-z0-9-]*-\[[^\]]+\]/],
  ['estilo em linha', /\bstyle=\{\{/],
];

describe('nenhum valor de estilo fora dos tokens', () => {
  const arquivos = arquivosDeCodigo(SRC);

  it('encontra os arquivos de código', () => {
    expect(arquivos.length).toBeGreaterThan(5);
  });

  it.each(PROIBIDOS)('sem %s', (_nome, padrao) => {
    const ocorrencias = arquivos.flatMap((arquivo) =>
      readFileSync(arquivo, 'utf8')
        .split('\n')
        .map((linha, i) => ({ linha, i }))
        .filter(({ linha }) => padrao.test(linha))
        .map(({ linha, i }) => `${relative(SRC, arquivo)}:${i + 1}: ${linha.trim()}`),
    );
    expect(ocorrencias).toEqual([]);
  });
});
