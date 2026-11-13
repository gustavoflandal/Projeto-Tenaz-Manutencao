// Minha conta: dados do próprio usuário, troca de senha e segundo fator.

import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it } from 'vitest';

import { problema, simularAPI } from '../../../testes/api-simulada';
import { sessao, usuario } from '../../../testes/dados';
import { renderizarRota, rotaDaSessao } from '../../../testes/renderizar';

const EU = usuario({ id: 'u-1', nome: 'Ana Souza', email: 'ana@polar.com.br', telefone: '31 99999-0000', versao: 2 });
const BASE = { ...rotaDaSessao(sessao()), 'GET /v1/usuarios/eu': { corpo: EU } };

describe('minha conta', () => {
  it('altera nome e telefone com If-Match, sem mexer no e-mail', async () => {
    const api = simularAPI({ ...BASE, 'PATCH /v1/usuarios/eu': { corpo: { ...EU, nome: 'Ana S. Lima', versao: 3 } } });
    const u = userEvent.setup();
    renderizarRota('/conta');

    const nome = await screen.findByLabelText('Nome');
    await waitFor(() => expect(nome).toHaveValue('Ana Souza'));
    expect(screen.getByLabelText('E-mail')).toHaveAttribute('readonly');
    await u.clear(nome);
    await u.type(nome, 'Ana S. Lima');
    await u.click(screen.getByRole('button', { name: 'Salvar' }));

    expect(await screen.findByText('Dados salvos.')).toBeInTheDocument();
    const [chamada] = api.chamadasA('PATCH /v1/usuarios/eu');
    expect(chamada?.corpo).toEqual({ nome: 'Ana S. Lima', telefone: '31 99999-0000' });
    expect(chamada?.cabecalhos.get('If-Match')).toBe('"2"');
  });

  it('troca a senha; a senha atual errada aparece no campo e não encerra a sessão', async () => {
    simularAPI({
      ...BASE,
      'POST /v1/auth/trocar-senha': problema(401, 'credencial_invalida', {
        instance: '/v1/auth/trocar-senha',
        erros: [{ campo: 'senha_atual', codigo: 'valor_invalido', mensagem: 'A senha atual não confere.' }],
      }),
    });
    const u = userEvent.setup();
    renderizarRota('/conta');

    await u.type(await screen.findByLabelText('Senha atual'), 'errada-errada');
    await u.type(screen.getByLabelText('Senha nova'), 'uma-senha-bem-longa');
    await u.type(screen.getByLabelText('Repita a senha nova'), 'uma-senha-bem-longa');
    await u.click(screen.getByRole('button', { name: 'Trocar senha' }));

    await waitFor(() => expect(screen.getByLabelText('Senha atual')).toHaveAccessibleDescription('A senha atual não confere.'));
    expect(screen.getByRole('heading', { level: 1, name: 'Minha conta' })).toBeInTheDocument();
  });

  it('troca a senha com sucesso e avisa que as outras sessões terminaram', async () => {
    const api = simularAPI({ ...BASE, 'POST /v1/auth/trocar-senha': { status: 204 } });
    const u = userEvent.setup();
    renderizarRota('/conta');

    await u.type(await screen.findByLabelText('Senha atual'), 'senha-antiga-123');
    await u.type(screen.getByLabelText('Senha nova'), 'uma-senha-bem-longa');
    await u.type(screen.getByLabelText('Repita a senha nova'), 'uma-senha-bem-longa');
    await u.click(screen.getByRole('button', { name: 'Trocar senha' }));

    expect(await screen.findByText(/Senha trocada/)).toBeInTheDocument();
    expect(api.chamadasA('POST /v1/auth/trocar-senha')[0]?.corpo).toEqual({ senha_atual: 'senha-antiga-123', nova_senha: 'uma-senha-bem-longa' });
  });

  it('senha fraca do servidor aparece no campo da senha nova', async () => {
    simularAPI({
      ...BASE,
      'POST /v1/auth/trocar-senha': problema(422, 'senha_fraca', {
        erros: [{ campo: 'senha', codigo: 'valor_invalido', mensagem: 'A senha não pode conter o seu e-mail.' }],
      }),
    });
    const u = userEvent.setup();
    renderizarRota('/conta');

    await u.type(await screen.findByLabelText('Senha atual'), 'senha-antiga-123');
    await u.type(screen.getByLabelText('Senha nova'), 'ana@polar-ana-1');
    await u.type(screen.getByLabelText('Repita a senha nova'), 'ana@polar-ana-1');
    await u.click(screen.getByRole('button', { name: 'Trocar senha' }));

    await waitFor(() => expect(screen.getByLabelText('Senha nova')).toHaveAccessibleDescription(/não pode conter o seu e-mail/));
  });

  it('com o segundo fator ativo, desativa com um código', async () => {
    const api = simularAPI({
      ...rotaDaSessao(sessao({ usuario: { id: 'u-1', nome: 'Ana Souza', email: 'ana@polar.com.br', pessoa_id: '', mfa_ativo: true, idioma: 'pt-BR' } })),
      'GET /v1/usuarios/eu': { corpo: EU },
      'POST /v1/auth/mfa/desativar': { status: 204 },
    });
    const u = userEvent.setup();
    renderizarRota('/conta');

    await u.click(await screen.findByRole('button', { name: 'Desativar segundo fator' }));
    await u.type(screen.getByLabelText('Código de 6 dígitos'), '123456');
    await u.click(screen.getAllByRole('button', { name: 'Desativar segundo fator' }).at(-1) as HTMLElement);

    await waitFor(() => expect(api.chamadasA('POST /v1/auth/mfa/desativar')).toHaveLength(1));
    expect(api.chamadasA('POST /v1/auth/mfa/desativar')[0]?.corpo).toEqual({ codigo: '123456' });
  });

  it('se a empresa exige o segundo fator, a desativação é recusada com a explicação', async () => {
    simularAPI({
      ...rotaDaSessao(sessao({ usuario: { id: 'u-1', nome: 'Ana Souza', email: 'ana@polar.com.br', pessoa_id: '', mfa_ativo: true, idioma: 'pt-BR' } })),
      'GET /v1/usuarios/eu': { corpo: EU },
      'POST /v1/auth/mfa/desativar': problema(409, 'transicao_invalida', { estado_atual: 'exigido_pela_empresa', instance: '/v1/auth/mfa/desativar' }),
    });
    const u = userEvent.setup();
    renderizarRota('/conta');

    await u.click(await screen.findByRole('button', { name: 'Desativar segundo fator' }));
    await u.type(screen.getByLabelText('Código de 6 dígitos'), '123456');
    await u.click(screen.getAllByRole('button', { name: 'Desativar segundo fator' }).at(-1) as HTMLElement);

    expect(await screen.findByRole('alert')).toHaveTextContent('A empresa exige o segundo fator');
  });
});

describe('prestador e locatário', () => {
  it('entram com o menu sem a administração e veem o aviso do portal', async () => {
    simularAPI(rotaDaSessao(sessao({ escopo_acesso: 'prestador', permissoes: ['os.ler'], perfis: ['prestador'] })));
    renderizarRota('/');

    expect(await screen.findByText(/portal para prestadores e locatários/)).toBeInTheDocument();
    expect(screen.queryByRole('link', { name: 'Empresa' })).not.toBeInTheDocument();
    expect(screen.queryByRole('link', { name: 'Usuários' })).not.toBeInTheDocument();
  });
});
