// Verificação automática de acessibilidade (axe-core) das telas desta etapa.
// Meta WCAG 2.1 AA (RNF-025, seção 11 do documento 7). O contraste de cor não
// é medido no jsdom; os pares de tokens são conferidos em pacotes/tokens-ui.

import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';

import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import axe from 'axe-core';
import { describe, expect, it } from 'vitest';

import { apiRespondeVersao, renderizarRota } from '../testes/renderizar';

const REGRAS: axe.RunOptions = {
  runOnly: { type: 'tag', values: ['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa', 'best-practice'] },
  resultTypes: ['violations'],
};

async function violacoes() {
  const resultado = await axe.run(document.documentElement, REGRAS);
  return resultado.violations.map((v) => ({
    regra: v.id,
    descricao: v.help,
    elementos: v.nodes.map((n) => n.target.join(' ')),
  }));
}

describe('acessibilidade WCAG 2.1 AA', () => {
  it.each([
    ['tela de login', '/entrar'],
    ['tela principal', '/'],
    ['página não encontrada', '/nao-existe'],
  ])('%s não tem violação', async (_nome, caminho) => {
    renderizarRota(caminho);
    expect(await violacoes()).toEqual([]);
  });

  it('diálogo Sobre aberto não tem violação', async () => {
    apiRespondeVersao({ versao: '0.0.0-dev', commit: 'a1b2c3d', data_build: '2026-11-12T14:05:00Z' });
    const usuario = userEvent.setup();
    renderizarRota('/');
    await usuario.click(screen.getByRole('button', { name: /^Versão / }));
    await screen.findByRole('dialog');
    expect(await violacoes()).toEqual([]);
  });

  it('menu de tema aberto não tem violação', async () => {
    const usuario = userEvent.setup();
    renderizarRota('/');
    await usuario.click(screen.getByRole('button', { name: /^Tema: / }));
    await screen.findByRole('menu');
    expect(await violacoes()).toEqual([]);
  });

  it('o index.html declara o idioma pt-BR (critério 3.1.1)', () => {
    const html = readFileSync(resolve(__dirname, '..', '..', 'index.html'), 'utf8');
    expect(html).toMatch(/<html lang="pt-BR">/);
  });
});
