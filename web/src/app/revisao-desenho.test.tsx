// Correções da revisão de desenho (documento 7): tooltip do menu recolhido,
// conflito de versão com as alterações não salvas, "Não autorizar" do OAuth,
// confirmação ao trocar turnos por um modelo, rolagem até o alerta e foco ao
// fim da lista.

import { screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';

import { problema, simularAPI } from '../testes/api-simulada';
import { empresa, pagina, pessoa, resumoPessoa, resumoUnidade, sessao, unidade } from '../testes/dados';
import { renderizarRota, rotaDaSessao } from '../testes/renderizar';

describe('menu recolhido', () => {
  it('mostra o nome do item em tooltip ao focar, além do aria-label', async () => {
    simularAPI(rotaDaSessao(sessao()));
    const u = userEvent.setup();
    renderizarRota('/');

    await u.click(await screen.findByRole('button', { name: 'Recolher menu' }));
    const item = screen.getByRole('link', { name: 'Empresa' });
    item.focus();

    expect(await screen.findByRole('tooltip')).toHaveTextContent('Empresa');
  });

  it('expandido, não precisa de tooltip e não tem item "em breve"', async () => {
    simularAPI(rotaDaSessao(sessao()));
    renderizarRota('/');
    const menu = await screen.findByRole('navigation', { name: 'Menu principal' });
    expect(within(menu).queryByText(/em breve/i)).not.toBeInTheDocument();
    expect(within(menu).queryAllByRole('button', { name: /em breve/i })).toHaveLength(0);
  });
});

describe('conflito de versão (412)', () => {
  it('lista o que a pessoa digitou e não salvou, campo e valor', async () => {
    simularAPI({
      ...rotaDaSessao(sessao()),
      'GET /v1/unidades': { corpo: pagina([resumoUnidade(unidade())]) },
      'GET /v1/pessoas/pe-1': { corpo: pessoa() },
      'PATCH /v1/pessoas/pe-1': problema(412, 'versao_desatualizada', { versao_atual: 6 }),
    });
    const u = userEvent.setup();
    renderizarRota('/administracao/pessoas/pe-1');

    const cargo = await screen.findByLabelText('Cargo (opcional)');
    await u.clear(cargo);
    await u.type(cargo, 'Líder');
    await u.click(screen.getByRole('button', { name: 'Salvar' }));

    const aviso = await screen.findByText('Suas alterações não salvas, para digitar de novo:');
    expect(within(aviso.closest('[role="status"]') as HTMLElement).getByText('Cargo: Líder')).toBeInTheDocument();
  });

  it('na empresa, o campo vazio aparece como "vazio"', async () => {
    simularAPI({
      ...rotaDaSessao(sessao()),
      'GET /v1/empresa': { corpo: empresa() },
      'PATCH /v1/empresa': problema(412, 'versao_desatualizada'),
    });
    const u = userEvent.setup();
    renderizarRota('/administracao/empresa');

    await u.clear(await screen.findByLabelText('Nome fantasia (opcional)'));
    await u.click(screen.getByRole('button', { name: 'Salvar' }));

    expect(await screen.findByText('Nome fantasia: vazio')).toBeInTheDocument();
  });
});

describe('/autorizar-app: Não autorizar', () => {
  it('devolve ao app o erro access_denied pela URI de retorno, com o state, e não vai ao painel', async () => {
    const api = simularAPI(rotaDaSessao(sessao()));
    const assign = vi.fn();
    vi.spyOn(globalThis, 'location', 'get').mockReturnValue({ ...globalThis.location, origin: 'http://localhost:3000', assign } as unknown as Location);
    const u = userEvent.setup();
    const { roteador } = renderizarRota(
      '/autorizar-app?response_type=code&client_id=tenaz-app&redirect_uri=br.com.tenaz.app%3A%2Foauth&code_challenge=abc&code_challenge_method=S256&state=xyz',
    );

    await u.click(await screen.findByRole('button', { name: 'Não autorizar' }));

    expect(assign).toHaveBeenCalledTimes(1);
    const destino = new URL(String(assign.mock.calls[0]?.[0]));
    expect(destino.protocol).toBe('br.com.tenaz.app:');
    expect(destino.searchParams.get('error')).toBe('access_denied');
    expect(destino.searchParams.get('state')).toBe('xyz');
    expect(destino.searchParams.has('code')).toBe(false);
    expect(roteador.state.location.pathname).toBe('/autorizar-app');
    expect(api.chamadasA('POST /oauth/authorize')).toHaveLength(0);
    expect(await screen.findByText(/Autorização negada/)).toBeInTheDocument();
  });
});

describe('editor de turnos', () => {
  it('pergunta antes de trocar turnos já digitados por um modelo', async () => {
    simularAPI({});
    const u = userEvent.setup();
    renderizarRota('/criar-conta');

    await u.click(await screen.findByRole('button', { name: 'Comercial' }));
    // Já há turnos: trocar pelo contínuo pede confirmação.
    await u.click(screen.getByRole('button', { name: 'Contínuo 24x7' }));
    const dialogo = await screen.findByRole('alertdialog', { name: 'Substituir os turnos pelo modelo?' });
    await u.click(within(dialogo).getByRole('button', { name: 'Voltar' }));
    expect(within(screen.getByRole('group', { name: 'Segunda-feira' })).getByLabelText('Fim')).toHaveValue('17:00');

    await u.click(screen.getByRole('button', { name: 'Contínuo 24x7' }));
    await u.click(within(await screen.findByRole('alertdialog')).getByRole('button', { name: 'Substituir turnos' }));
    await waitFor(() => expect(within(screen.getByRole('group', { name: 'Segunda-feira' })).getByLabelText('Fim')).toHaveValue('24:00'));
  });
});

describe('alertas e listas', () => {
  it('o alerta de sucesso no topo do formulário rola até ele', async () => {
    const rolar = vi.spyOn(Element.prototype, 'scrollIntoView');
    simularAPI({
      ...rotaDaSessao(sessao()),
      'GET /v1/empresa': { corpo: empresa() },
      'PATCH /v1/empresa': { corpo: empresa({ nome_fantasia: 'Novo', versao: 4 }) },
    });
    const u = userEvent.setup();
    renderizarRota('/administracao/empresa');

    const campo = await screen.findByLabelText('Nome fantasia (opcional)');
    await u.clear(campo);
    await u.type(campo, 'Novo');
    await u.click(screen.getByRole('button', { name: 'Salvar' }));

    await screen.findByText('Dados da empresa salvos.');
    expect(rolar).toHaveBeenCalled();
  });

  it('ao sumir o "Carregar mais", o foco vai para a contagem da lista', async () => {
    simularAPI({
      ...rotaDaSessao(sessao()),
      'GET /v1/unidades': { corpo: pagina([]) },
      'GET /v1/pessoas': (c) =>
        c.consulta.get('cursor') ? { corpo: pagina([resumoPessoa(pessoa({ id: 'pe-2', nome: 'Maria' }))]) } : { corpo: pagina([resumoPessoa()], 'C2') },
    });
    const u = userEvent.setup();
    renderizarRota('/administracao/pessoas');

    await u.click(await screen.findByRole('button', { name: 'Carregar mais' }));

    await waitFor(() => expect(screen.getByText('Fim da lista · 2 itens')).toHaveFocus());
  });
});
