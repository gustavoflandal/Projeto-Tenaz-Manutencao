import { screen, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it } from 'vitest';

import { simularAPI } from '../testes/api-simulada';
import { sessao } from '../testes/dados';
import { renderizarRota, rotaDaSessao } from '../testes/renderizar';

describe('casca da tela principal', () => {
  it('tem os marcos de navegação, cabeçalho, conteúdo e o link para pular ao conteúdo', async () => {
    simularAPI(rotaDaSessao(sessao()));
    renderizarRota('/');
    expect(await screen.findByRole('navigation', { name: 'Menu principal' })).toBeInTheDocument();
    expect(screen.getByRole('banner')).toBeInTheDocument();
    expect(screen.getByRole('main')).toHaveAttribute('id', 'conteudo');
    expect(screen.getByRole('link', { name: 'Pular para o conteúdo' })).toHaveAttribute('href', '#conteudo');
    expect(screen.getByRole('heading', { level: 1, name: 'Painel' })).toBeInTheDocument();
  });

  it('marca o item ativo do menu com aria-current', async () => {
    simularAPI(rotaDaSessao(sessao()));
    renderizarRota('/');
    expect(await screen.findByRole('link', { name: 'Painel' })).toHaveAttribute('aria-current', 'page');
  });

  it('recolhe o menu, mantém o nome do item e guarda a escolha', async () => {
    simularAPI(rotaDaSessao(sessao()));
    const usuario = userEvent.setup();
    renderizarRota('/');

    await usuario.click(await screen.findByRole('button', { name: 'Recolher menu' }));

    expect(screen.getByRole('button', { name: 'Expandir menu' })).toHaveAttribute('aria-expanded', 'false');
    expect(screen.getByRole('link', { name: 'Painel' })).toBeInTheDocument();
    expect(globalThis.localStorage.getItem('tenaz.menu')).toBe('recolhido');
  });

  it('mostra no menu só o que a permissão libera', async () => {
    simularAPI(rotaDaSessao(sessao({ permissoes: ['adm.ler', 'adm.ler_auditoria'] })));
    renderizarRota('/');
    const menu = await screen.findByRole('navigation', { name: 'Menu principal' });

    expect(within(menu).getByRole('link', { name: 'Empresa' })).toBeInTheDocument();
    expect(within(menu).getByRole('link', { name: 'Auditoria' })).toBeInTheDocument();
    expect(within(menu).queryByRole('link', { name: 'Usuários' })).not.toBeInTheDocument();
    expect(within(menu).queryByRole('link', { name: 'Pessoas' })).not.toBeInTheDocument();
    expect(within(menu).getByRole('link', { name: 'Integridade da trilha' })).toHaveAttribute('href', '/administracao/auditoria/integridade');
    // Sem adm.gerenciar_parametros, as capturas nem aparecem.
    expect(within(menu).queryByRole('link', { name: 'Capturas da auditoria' })).not.toBeInTheDocument();
  });

  it('mostra a empresa e o selo de versão', async () => {
    simularAPI(rotaDaSessao(sessao()));
    renderizarRota('/');
    expect(await screen.findByText('Refrigeração Polar')).toBeInTheDocument();
    expect(screen.getByRole('button', { name: /^Versão / })).toBeInTheDocument();
  });
});
