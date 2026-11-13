import { screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it } from 'vitest';

import { renderizarRota } from '../testes/renderizar';

describe('casca da tela principal', () => {
  it('tem os marcos de navegação, cabeçalho, conteúdo e o link para pular ao conteúdo', () => {
    renderizarRota('/');
    expect(screen.getByRole('navigation', { name: 'Menu principal' })).toBeInTheDocument();
    expect(screen.getByRole('banner')).toBeInTheDocument();
    expect(screen.getByRole('main')).toHaveAttribute('id', 'conteudo');
    expect(screen.getByRole('link', { name: 'Pular para o conteúdo' })).toHaveAttribute('href', '#conteudo');
    expect(screen.getByRole('heading', { level: 1, name: 'Painel' })).toBeInTheDocument();
  });

  it('marca o item ativo do menu com aria-current', () => {
    renderizarRota('/');
    expect(screen.getByRole('link', { name: 'Painel' })).toHaveAttribute('aria-current', 'page');
  });

  it('recolhe o menu, mantém o nome do item e guarda a escolha', async () => {
    const usuario = userEvent.setup();
    renderizarRota('/');

    await usuario.click(screen.getByRole('button', { name: 'Recolher menu' }));

    expect(screen.getByRole('button', { name: 'Expandir menu' })).toHaveAttribute('aria-expanded', 'false');
    expect(screen.getByRole('link', { name: 'Painel' })).toBeInTheDocument();
    expect(globalThis.localStorage.getItem('tenaz.menu')).toBe('recolhido');
  });
});
