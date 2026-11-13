// Escolha de idioma (pt-BR e es-419): detecção, menu do usuário, guarda no
// navegador, atributo lang do documento e formatos de data.

import { screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it } from 'vitest';

import { idiomaDoNavegador, idiomaInicial, lerIdiomaSalvo } from '../lib/idioma';
import { formatarDataHora } from '../lib/formatacao';
import { i18n } from '../lib/i18n';
import { simularAPI } from '../testes/api-simulada';
import { sessao } from '../testes/dados';
import { renderizarRota, rotaDaSessao } from '../testes/renderizar';

describe('detecção do idioma', () => {
  it('qualquer variante de espanhol vira es-419 e de português vira pt-BR', () => {
    expect(idiomaDoNavegador(['es-AR', 'en'])).toBe('es-419');
    expect(idiomaDoNavegador(['es'])).toBe('es-419');
    expect(idiomaDoNavegador(['en-US', 'pt-PT'])).toBe('pt-BR');
    expect(idiomaDoNavegador(['en-US', 'fr'])).toBeNull();
  });

  it('a escolha guardada vence o idioma do navegador', () => {
    globalThis.localStorage.setItem('tenaz.idioma', 'es-419');
    expect(lerIdiomaSalvo()).toBe('es-419');
    expect(idiomaInicial()).toBe('es-419');
  });

  it('valor guardado inválido é ignorado', () => {
    globalThis.localStorage.setItem('tenaz.idioma', 'fr');
    expect(lerIdiomaSalvo()).toBeNull();
  });
});

describe('Minha conta', () => {
  const EU = { id: 'u-1', nome: 'Ana', email: 'a@b.c', telefone: null, versao: 1 };

  it('troca para espanhol, guarda a escolha e muda lang e título do documento', async () => {
    simularAPI({ ...rotaDaSessao(sessao()), 'GET /v1/usuarios/eu': { corpo: EU } });
    const u = userEvent.setup();
    renderizarRota('/conta');

    await u.click(await screen.findByRole('radio', { name: 'Español' }));

    await waitFor(() => expect(document.documentElement).toHaveAttribute('lang', 'es-419'));
    expect(await screen.findByRole('heading', { level: 1, name: 'Mi cuenta' })).toBeInTheDocument();
    expect(globalThis.localStorage.getItem('tenaz.idioma')).toBe('es-419');
    expect(document.title).toMatch(/Gestión de Mantenimiento/);
  });

  it('marca o idioma atual e volta ao português', async () => {
    await i18n.changeLanguage('es-419');
    simularAPI({ ...rotaDaSessao(sessao()), 'GET /v1/usuarios/eu': { corpo: EU } });
    const u = userEvent.setup();
    renderizarRota('/conta');

    expect(await screen.findByRole('radio', { name: 'Español' })).toBeChecked();
    await u.click(screen.getByRole('radio', { name: 'Português (Brasil)' }));

    expect(await screen.findByRole('heading', { level: 1, name: 'Minha conta' })).toBeInTheDocument();
    expect(document.documentElement).toHaveAttribute('lang', 'pt-BR');
  });
});

describe('telas sem sessão', () => {
  it('têm links de idioma no pé do cartão, com o nome de cada idioma no próprio idioma', async () => {
    const u = userEvent.setup();
    renderizarRota('/entrar');

    const links = await screen.findByRole('navigation', { name: 'Idioma da interface' });
    expect(within(links).getByRole('button', { name: 'Português (Brasil)' })).toHaveAttribute('aria-current', 'true');
    await u.click(within(links).getByRole('button', { name: 'Español' }));

    expect(await screen.findByRole('heading', { level: 1, name: 'Iniciar sesión' })).toBeInTheDocument();
    expect(globalThis.localStorage.getItem('tenaz.idioma')).toBe('es-419');
  });
});

describe('telas em espanhol', () => {
  it('o login aparece em espanhol, com o termo do glossário', async () => {
    await i18n.changeLanguage('es-419');
    renderizarRota('/entrar');
    expect(await screen.findByRole('heading', { level: 1, name: 'Iniciar sesión' })).toBeInTheDocument();
    expect(screen.getByLabelText('Contraseña')).toBeInTheDocument();
  });

  it('data e hora seguem o país da empresa, não o idioma da interface', async () => {
    expect(formatarDataHora('2026-11-12T14:05:00Z', 'UTC', 'BR')).toBe('12/11/2026 14:05');
    await i18n.changeLanguage('es-419');
    expect(formatarDataHora('2026-11-12T14:05:00Z', 'UTC', 'BR')).toBe('12/11/2026 14:05');
  });

});
