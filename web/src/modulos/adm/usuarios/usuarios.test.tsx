// Usuários e perfis (RF-ADM-004 e RF-ADM-005).

import { screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it } from 'vitest';

import { problema, simularAPI } from '../../../testes/api-simulada';
import { PERFIS, PERMISSOES, pagina, pessoa, resumoPessoa, resumoUnidade, resumoUsuario, sessao, unidade, usuario } from '../../../testes/dados';
import { renderizarRota, rotaDaSessao } from '../../../testes/renderizar';

const BASE = {
  ...rotaDaSessao(sessao()),
  'GET /v1/perfis': { corpo: pagina(PERFIS) },
  'GET /v1/unidades': { corpo: pagina([resumoUnidade(unidade())]) },
  'GET /v1/permissoes': { corpo: pagina(PERMISSOES) },
};

describe('lista de usuários', () => {
  it('mostra nome, e-mail, perfil pelo nome e situação', async () => {
    simularAPI({ ...BASE, 'GET /v1/usuarios': { corpo: pagina([resumoUsuario(), resumoUsuario(usuario({ id: 'u-3', nome: 'Bia Reis', situacao: 'convidado' }))]) } });
    renderizarRota('/administracao/usuarios');

    const linha = await screen.findByRole('row', { name: /Carlos Lima/ });
    expect(within(linha).getByText('Técnico')).toBeInTheDocument();
    expect(within(linha).getByText('Ativo')).toBeInTheDocument();
    expect(within(screen.getByRole('row', { name: /Bia Reis/ })).getByText('Convidado')).toBeInTheDocument();
  });

  it('filtra por situação e perfil no servidor', async () => {
    const api = simularAPI({ ...BASE, 'GET /v1/usuarios': { corpo: pagina([resumoUsuario()]) } });
    const usuarioEvento = userEvent.setup();
    renderizarRota('/administracao/usuarios');

    await usuarioEvento.selectOptions(await screen.findByLabelText('Situação'), 'inativo');
    await usuarioEvento.selectOptions(screen.getByLabelText('Perfil'), 'tecnico');

    await waitFor(() => {
      const ultima = api.chamadasA('GET /v1/usuarios').at(-1);
      expect(ultima?.consulta.get('situacao')).toBe('inativo');
      expect(ultima?.consulta.get('perfil')).toBe('tecnico');
    });
  });
});

describe('convite de usuário', () => {
  async function preencher(u: ReturnType<typeof userEvent.setup>) {
    await u.type(await screen.findByLabelText('Nome'), 'Bia Reis');
    await u.type(screen.getByLabelText('E-mail'), 'bia@polar.com.br');
  }

  it('não oferece os perfis de prestador e de locatário', async () => {
    simularAPI(BASE);
    renderizarRota('/administracao/usuarios/novo');
    expect(await screen.findByRole('checkbox', { name: 'Técnico' })).toBeInTheDocument();
    expect(screen.queryByRole('checkbox', { name: 'Prestador' })).not.toBeInTheDocument();
  });

  it('exige ao menos um perfil', async () => {
    const api = simularAPI(BASE);
    const u = userEvent.setup();
    renderizarRota('/administracao/usuarios/novo');
    await preencher(u);

    await u.click(screen.getByRole('button', { name: 'Convidar usuário' }));

    expect((await screen.findAllByText('Escolha ao menos um perfil.')).length).toBeGreaterThan(0);
    expect(api.chamadasA('POST /v1/usuarios')).toHaveLength(0);
  });

  it('convida criando a pessoa junto e mostra o aviso do convite', async () => {
    const criado = { ...usuario({ id: 'u-9', nome: 'Bia Reis', email: 'bia@polar.com.br', situacao: 'convidado' }), convite_enviado: true };
    const api = simularAPI({
      ...BASE,
      'POST /v1/usuarios': { status: 201, corpo: criado },
      'GET /v1/usuarios/u-9': { corpo: usuario({ id: 'u-9', nome: 'Bia Reis', email: 'bia@polar.com.br', situacao: 'convidado' }) },
    });
    const u = userEvent.setup();
    renderizarRota('/administracao/usuarios/novo');
    await preencher(u);
    await u.click(await screen.findByRole('checkbox', { name: 'Técnico' }));
    await u.click(screen.getByRole('radio', { name: /Criar a pessoa junto/ }));
    await u.click(screen.getByRole('button', { name: 'Convidar usuário' }));

    expect(await screen.findByRole('heading', { level: 1, name: 'Bia Reis' })).toBeInTheDocument();
    expect(await screen.findByText(/Convite enviado/)).toBeInTheDocument();
    expect(api.chamadasA('POST /v1/usuarios')[0]?.corpo).toEqual({
      nome: 'Bia Reis',
      email: 'bia@polar.com.br',
      perfis: ['tecnico'],
      nova_pessoa: { nome: 'Bia Reis', email: 'bia@polar.com.br' },
    });
    expect(api.chamadasA('POST /v1/usuarios')[0]?.cabecalhos.get('Idempotency-Key')).toBeTruthy();
  });

  it('vincula a uma pessoa existente buscada no servidor', async () => {
    const api = simularAPI({
      ...BASE,
      'GET /v1/pessoas': { corpo: pagina([resumoPessoa(pessoa())]) },
      'POST /v1/usuarios': { status: 201, corpo: { ...usuario({ id: 'u-9' }), convite_enviado: false } },
      'GET /v1/usuarios/u-9': { corpo: usuario({ id: 'u-9' }) },
    });
    const u = userEvent.setup();
    renderizarRota('/administracao/usuarios/novo');
    await preencher(u);
    await u.click(await screen.findByRole('checkbox', { name: 'Técnico' }));
    await u.click(screen.getByRole('radio', { name: /Vincular a uma pessoa já cadastrada/ }));
    await u.selectOptions(await screen.findByRole('combobox', { name: 'Pessoa' }), 'pe-1');
    await u.click(screen.getByRole('button', { name: 'Convidar usuário' }));

    expect(await screen.findByText(/e-mail do convite não saiu/)).toBeInTheDocument();
    expect(api.chamadasA('POST /v1/usuarios')[0]?.corpo).toMatchObject({ pessoa_id: 'pe-1' });
    expect(api.chamadasA('GET /v1/pessoas')[0]?.consulta.get('com_usuario')).toBe('false');
  });

  it('e-mail repetido na empresa (409) aparece como erro geral', async () => {
    simularAPI({ ...BASE, 'POST /v1/usuarios': problema(409, 'conflito') });
    const u = userEvent.setup();
    renderizarRota('/administracao/usuarios/novo');
    await preencher(u);
    await u.click(await screen.findByRole('checkbox', { name: 'Técnico' }));
    await u.click(screen.getByRole('button', { name: 'Convidar usuário' }));
    expect(await screen.findByRole('alert')).toHaveTextContent('Já existe um registro com estes dados');
  });

  it('erro de campo do servidor aparece no campo', async () => {
    simularAPI({
      ...BASE,
      'POST /v1/usuarios': problema(422, 'validacao_falhou', {
        erros: [{ campo: 'email', codigo: 'duplicado', mensagem: 'Este e-mail já tem usuário na empresa.' }],
      }),
    });
    const u = userEvent.setup();
    renderizarRota('/administracao/usuarios/novo');
    await preencher(u);
    await u.click(await screen.findByRole('checkbox', { name: 'Técnico' }));
    await u.click(screen.getByRole('button', { name: 'Convidar usuário' }));
    await waitFor(() => expect(screen.getByLabelText('E-mail')).toHaveAccessibleDescription(/já tem usuário/));
  });
});

describe('detalhe do usuário', () => {
  const ROTAS = { ...BASE, 'GET /v1/usuarios/u-2': { corpo: usuario() } };

  it('altera perfis e dados com If-Match, mandando só o que mudou', async () => {
    const api = simularAPI({ ...ROTAS, 'PATCH /v1/usuarios/u-2': { corpo: usuario({ versao: 8 }) } });
    const u = userEvent.setup();
    renderizarRota('/administracao/usuarios/u-2');

    await u.click(await screen.findByRole('checkbox', { name: 'Administrador' }));
    await u.click(screen.getAllByRole('button', { name: 'Salvar' })[0] as HTMLElement);

    await waitFor(() => expect(api.chamadasA('PATCH /v1/usuarios/u-2')).toHaveLength(1));
    const [chamada] = api.chamadasA('PATCH /v1/usuarios/u-2');
    expect(chamada?.corpo).toEqual({ perfis: ['tecnico', 'administrador'] });
    expect(chamada?.cabecalhos.get('If-Match')).toBe('"7"');
  });

  it('desativar o último administrador é recusado e a mensagem explica o que fazer', async () => {
    simularAPI({ ...ROTAS, 'POST /v1/usuarios/u-2/desativar': problema(409, 'ultimo_administrador') });
    const u = userEvent.setup();
    renderizarRota('/administracao/usuarios/u-2');

    await u.click(await screen.findByRole('button', { name: 'Desativar' }));
    const dialogo = await screen.findByRole('alertdialog', { name: 'Desativar Carlos Lima?' });
    await u.click(within(dialogo).getByRole('button', { name: 'Desativar usuário' }));

    expect(await within(dialogo).findByRole('alert')).toHaveTextContent('ao menos um administrador ativo');
  });

  it('desativa com If-Match e depois permite reativar', async () => {
    const api = simularAPI({
      ...ROTAS,
      'POST /v1/usuarios/u-2/desativar': { corpo: usuario({ situacao: 'inativo', versao: 8 }) },
      'POST /v1/usuarios/u-2/reativar': { corpo: { ...usuario({ versao: 9 }), convite_enviado: false } },
    });
    const u = userEvent.setup();
    renderizarRota('/administracao/usuarios/u-2');

    await u.click(await screen.findByRole('button', { name: 'Desativar' }));
    api.definir('GET /v1/usuarios/u-2', { corpo: usuario({ situacao: 'inativo', versao: 8 }) });
    await u.click(within(await screen.findByRole('alertdialog')).getByRole('button', { name: 'Desativar usuário' }));

    const reativar = await screen.findByRole('button', { name: 'Reativar' });
    expect(api.chamadasA('POST /v1/usuarios/u-2/desativar')[0]?.cabecalhos.get('If-Match')).toBe('"7"');
    await u.click(reativar);
    await waitFor(() => expect(api.chamadasA('POST /v1/usuarios/u-2/reativar')).toHaveLength(1));
    expect(api.chamadasA('POST /v1/usuarios/u-2/reativar')[0]?.cabecalhos.get('If-Match')).toBe('"8"');
  });

  it('reenvia convite de usuário convidado, sem If-Match', async () => {
    const api = simularAPI({ ...BASE, 'GET /v1/usuarios/u-2': { corpo: usuario({ situacao: 'convidado' }) }, 'POST /v1/usuarios/u-2/reenviar-convite': { status: 204 } });
    const u = userEvent.setup();
    renderizarRota('/administracao/usuarios/u-2');

    await u.click(await screen.findByRole('button', { name: 'Reenviar convite' }));

    expect(await screen.findByText(/Convite reenviado para carlos@polar.com.br/)).toBeInTheDocument();
    expect(api.chamadasA('POST /v1/usuarios/u-2/reenviar-convite')[0]?.cabecalhos.has('If-Match')).toBe(false);
  });

  it('encerra as sessões e informa quantas foram', async () => {
    simularAPI({ ...ROTAS, 'POST /v1/usuarios/u-2/encerrar-sessoes': { corpo: { sessoes_encerradas: 3 } } });
    const u = userEvent.setup();
    renderizarRota('/administracao/usuarios/u-2');

    await u.click(await screen.findByRole('button', { name: 'Encerrar sessões' }));
    await u.click(within(await screen.findByRole('alertdialog')).getByRole('button', { name: 'Encerrar sessões' }));

    expect(await screen.findByText('Sessões encerradas: 3.')).toBeInTheDocument();
  });

  it('concede e retira permissão avulsa', async () => {
    const api = simularAPI({
      ...ROTAS,
      'POST /v1/usuarios/u-2/permissoes': { status: 201, corpo: usuario({ permissoes_avulsas: ['os.encerrar'] }) },
      'DELETE /v1/usuarios/u-2/permissoes/os.encerrar': { status: 204 },
    });
    const u = userEvent.setup();
    renderizarRota('/administracao/usuarios/u-2');

    await screen.findByRole('option', { name: 'os.encerrar' });
    await u.selectOptions(screen.getByLabelText('Conceder permissão'), 'os.encerrar');
    api.definir('GET /v1/usuarios/u-2', { corpo: usuario({ permissoes_avulsas: ['os.encerrar'], versao: 7 }) });
    await u.click(screen.getByRole('button', { name: 'Conceder' }));

    await waitFor(() => expect(api.chamadasA('POST /v1/usuarios/u-2/permissoes')[0]?.corpo).toEqual({ permissao: 'os.encerrar' }));
    const retirar = await screen.findByRole('button', { name: 'Retirar a permissão os.encerrar' });
    await u.click(retirar);
    await waitFor(() => expect(api.chamadasA('DELETE /v1/usuarios/u-2/permissoes/os.encerrar')).toHaveLength(1));
  });

  it('permissão avulsa recusada pelo servidor mostra o motivo', async () => {
    simularAPI({ ...ROTAS, 'POST /v1/usuarios/u-2/permissoes': problema(403, 'permissao_negada') });
    const u = userEvent.setup();
    renderizarRota('/administracao/usuarios/u-2');

    await screen.findByRole('option', { name: 'os.ler' });
    await u.selectOptions(screen.getByLabelText('Conceder permissão'), 'os.ler');
    await u.click(screen.getByRole('button', { name: 'Conceder' }));
    expect(await screen.findByRole('alert')).toHaveTextContent('Seu perfil não pode fazer isso');
  });

  it('no próprio usuário, ninguém concede nem retira permissão de si mesmo', async () => {
    simularAPI({ ...BASE, 'GET /v1/usuarios/u-1': { corpo: usuario({ id: 'u-1', permissoes_avulsas: ['os.ler'] }) } });
    renderizarRota('/administracao/usuarios/u-1');
    expect(await screen.findByText(/Ninguém concede nem retira permissão de si mesmo/)).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Conceder' })).not.toBeInTheDocument();
    expect(screen.queryByRole('button', { name: /Retirar a permissão/ })).not.toBeInTheDocument();
  });

  it('usuário de outra empresa (404) mostra não encontrado', async () => {
    simularAPI({ ...BASE, 'GET /v1/usuarios/zz': problema(404, 'recurso_nao_encontrado') });
    renderizarRota('/administracao/usuarios/zz');
    expect(await screen.findByRole('heading', { name: 'Usuário não encontrado' })).toBeInTheDocument();
  });

  it('só leitura: sem permissão de gerenciar, não há ações', async () => {
    simularAPI({ ...rotaDaSessao(sessao({ permissoes: ['adm.ler_usuarios'] })), 'GET /v1/perfis': { corpo: pagina(PERFIS) }, 'GET /v1/unidades': { corpo: pagina([]) }, 'GET /v1/usuarios/u-2': { corpo: usuario() } });
    renderizarRota('/administracao/usuarios/u-2');
    await screen.findByRole('heading', { level: 1, name: 'Carlos Lima' });
    expect(screen.queryByRole('button', { name: 'Desativar' })).not.toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Salvar' })).not.toBeInTheDocument();
  });
});

describe('perfis', () => {
  it('lista os perfis com a quantidade de permissões e destaca as sensíveis', async () => {
    simularAPI({
      ...BASE,
      'GET /v1/perfis': { corpo: pagina(PERFIS.slice(0, 1).map((p) => ({ ...p, permissoes: ['os.encerrar', 'os.ler'] }))) },
    });
    const u = userEvent.setup();
    renderizarRota('/administracao/perfis');

    const resumo = await screen.findByText('Administrador');
    expect(screen.getByText(/Permissões: 2/)).toBeInTheDocument();
    await u.click(resumo);
    const sensivel = await screen.findByText('os.encerrar');
    expect(within(sensivel.closest('li') as HTMLElement).getByText('Permissão sensível:')).toBeInTheDocument();
    expect(within(screen.getByText('os.ler').closest('li') as HTMLElement).queryByText('Permissão sensível:')).not.toBeInTheDocument();
  });

  it('sem o catálogo (403), lista os perfis sem destacar', async () => {
    simularAPI({ ...BASE, 'GET /v1/permissoes': problema(403, 'permissao_negada') });
    renderizarRota('/administracao/perfis');
    expect(await screen.findByText('Administrador')).toBeInTheDocument();
  });
});
