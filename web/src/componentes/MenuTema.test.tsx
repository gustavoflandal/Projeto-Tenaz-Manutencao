import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it } from 'vitest';

import { renderizarRota } from '../testes/renderizar';

const raiz = () => document.documentElement;

async function escolherTema(usuario: ReturnType<typeof userEvent.setup>, atual: string, opcao: string) {
  await usuario.click(screen.getByRole('button', { name: `Tema: ${atual}. Trocar tema` }));
  await usuario.click(await screen.findByRole('menuitemradio', { name: opcao }));
  await waitFor(() => expect(screen.queryByRole('menu')).not.toBeInTheDocument());
}

describe('troca de tema', () => {
  it('começa igual ao sistema, sem atributo de tema', () => {
    renderizarRota('/');
    expect(screen.getByRole('button', { name: 'Tema: Igual ao sistema. Trocar tema' })).toBeInTheDocument();
    expect(raiz()).not.toHaveAttribute('data-tema');
  });

  it('aplica o tema escolhido e guarda a escolha', async () => {
    const usuario = userEvent.setup();
    renderizarRota('/');

    await escolherTema(usuario, 'Igual ao sistema', 'Escuro');

    expect(raiz()).toHaveAttribute('data-tema', 'escuro');
    expect(globalThis.localStorage.getItem('tenaz.tema')).toBe('escuro');
    expect(screen.getByRole('button', { name: 'Tema: Escuro. Trocar tema' })).toBeInTheDocument();
  });

  it('volta a seguir o sistema', async () => {
    const usuario = userEvent.setup();
    renderizarRota('/');

    await escolherTema(usuario, 'Igual ao sistema', 'Sol (alto contraste)');
    expect(raiz()).toHaveAttribute('data-tema', 'sol');

    await escolherTema(usuario, 'Sol (alto contraste)', 'Igual ao sistema');
    expect(raiz()).not.toHaveAttribute('data-tema');
    expect(globalThis.localStorage.getItem('tenaz.tema')).toBeNull();
  });

  it('marca a opção atual no menu', async () => {
    const usuario = userEvent.setup();
    renderizarRota('/');

    await escolherTema(usuario, 'Igual ao sistema', 'Claro');
    await usuario.click(screen.getByRole('button', { name: 'Tema: Claro. Trocar tema' }));

    expect(await screen.findByRole('menuitemradio', { name: 'Claro' })).toHaveAttribute('aria-checked', 'true');
    expect(screen.getByRole('menuitemradio', { name: 'Escuro' })).toHaveAttribute('aria-checked', 'false');
  });
});
