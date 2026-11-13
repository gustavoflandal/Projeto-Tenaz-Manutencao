import { screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';

import { apiForaDoAr, apiRespondeVersao, renderizarRota } from '../../testes/renderizar';

// Nos testes, o vite.config.ts injeta a versão 0.0.0-dev, sem commit nem data
// de build, no ambiente de desenvolvimento.
const NOME_DO_SELO = 'Versão 0.0.0-dev, ambiente Desenvolvimento. Abrir informações sobre o sistema';
const API_NA_MESMA_VERSAO = { versao: '0.0.0-dev', commit: 'a1b2c3d', data_build: '2026-11-12T14:05:00Z' };

function linha(dialogo: HTMLElement, chave: string): string {
  return dialogo.querySelector(`[data-linha="${chave}"]`)?.textContent ?? '';
}

function focaveis(): HTMLElement[] {
  return Array.from(
    document.querySelectorAll<HTMLElement>('a[href], button:not([disabled]), input:not([disabled]), [tabindex="0"]'),
  );
}

describe.each(['/entrar', '/'])('selo de versão em %s (RNF-039)', (caminho) => {
  it('aparece com o texto e o nome acessível do componente 6.19', () => {
    renderizarRota(caminho);
    const selo = screen.getByRole('button', { name: NOME_DO_SELO });
    expect(selo).toHaveTextContent('v0.0.0-dev · Desenvolvimento');
  });

  it('é o último elemento focável da página', () => {
    renderizarRota(caminho);
    expect(focaveis().at(-1)).toBe(screen.getByRole('button', { name: NOME_DO_SELO }));
  });
});

describe('diálogo Sobre', () => {
  it('abre pelo teclado e mostra web, API, data do build, commit e ambiente', async () => {
    apiRespondeVersao(API_NA_MESMA_VERSAO);
    const usuario = userEvent.setup();
    renderizarRota('/');

    screen.getByRole('button', { name: NOME_DO_SELO }).focus();
    await usuario.keyboard('{Enter}');

    const dialogo = await screen.findByRole('dialog', { name: 'Sobre a Tenaz' });
    await waitFor(() => expect(linha(dialogo, 'api')).toBe('v0.0.0-dev'));
    expect(within(dialogo).getByText('Web')).toBeInTheDocument();
    expect(linha(dialogo, 'web')).toBe('v0.0.0-dev');
    expect(linha(dialogo, 'dataBuild')).toBe('Build local, sem data');
    expect(linha(dialogo, 'commit')).toBe('Não informado');
    expect(linha(dialogo, 'ambiente')).toBe('Desenvolvimento');
    expect(within(dialogo).queryByText(/versão mais nova/)).not.toBeInTheDocument();
  });

  it('diz que a API não respondeu e mantém a versão da web', async () => {
    apiForaDoAr();
    const usuario = userEvent.setup();
    renderizarRota('/entrar');

    await usuario.click(screen.getByRole('button', { name: NOME_DO_SELO }));

    const dialogo = await screen.findByRole('dialog', { name: 'Sobre a Tenaz' });
    await waitFor(() =>
      expect(linha(dialogo, 'api')).toBe('A API não respondeu. Verifique a conexão e abra de novo.'),
    );
    expect(linha(dialogo, 'web')).toBe('v0.0.0-dev');
  });

  it('avisa quando a API está numa versão diferente da web', async () => {
    apiRespondeVersao({ ...API_NA_MESMA_VERSAO, versao: '1.5.0' });
    const usuario = userEvent.setup();
    renderizarRota('/');

    await usuario.click(screen.getByRole('button', { name: NOME_DO_SELO }));

    const dialogo = await screen.findByRole('dialog', { name: 'Sobre a Tenaz' });
    expect(
      await within(dialogo).findByText('Há uma versão mais nova. Recarregue a página para atualizar.'),
    ).toBeInTheDocument();
    expect(within(dialogo).getByRole('button', { name: 'Recarregar' })).toBeInTheDocument();
  });

  it('fecha com Esc e devolve o foco ao selo', async () => {
    apiRespondeVersao(API_NA_MESMA_VERSAO);
    const usuario = userEvent.setup();
    renderizarRota('/');
    const selo = screen.getByRole('button', { name: NOME_DO_SELO });

    await usuario.click(selo);
    await screen.findByRole('dialog', { name: 'Sobre a Tenaz' });
    await usuario.keyboard('{Escape}');

    await waitFor(() => expect(screen.queryByRole('dialog')).not.toBeInTheDocument());
    expect(selo).toHaveFocus();
  });

  it('copia as informações em texto simples e confirma', async () => {
    apiRespondeVersao(API_NA_MESMA_VERSAO);
    const usuario = userEvent.setup();
    const escrever = vi.spyOn(globalThis.navigator.clipboard, 'writeText').mockResolvedValue();
    renderizarRota('/');

    await usuario.click(screen.getByRole('button', { name: NOME_DO_SELO }));
    const dialogo = await screen.findByRole('dialog', { name: 'Sobre a Tenaz' });
    await waitFor(() => expect(linha(dialogo, 'api')).toBe('v0.0.0-dev'));
    await usuario.click(within(dialogo).getByRole('button', { name: 'Copiar informações' }));

    expect(escrever).toHaveBeenCalledWith(
      [
        'Web: v0.0.0-dev',
        'API: v0.0.0-dev',
        'Data do build: Build local, sem data',
        'Commit: Não informado',
        'Ambiente: Desenvolvimento',
      ].join('\n'),
    );
    expect(await within(dialogo).findByText('Informações da versão copiadas.')).toBeInTheDocument();
  });
});
