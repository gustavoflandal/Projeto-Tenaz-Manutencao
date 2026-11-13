// Unidades, turnos e feriados (RF-ADM-002, RN-ADM-019).

import { screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it } from 'vitest';

import { pagina, resumoUnidade, sessao, unidade, usuario as usuarioDeTeste } from '../../../testes/dados';
import { problema, simularAPI } from '../../../testes/api-simulada';
import { renderizarRota, rotaDaSessao } from '../../../testes/renderizar';

const BASE = { ...rotaDaSessao(sessao()) };

describe('lista de unidades', () => {
  it('mostra código, nome, localização e situação', async () => {
    simularAPI({ ...BASE, 'GET /v1/unidades': { corpo: pagina([resumoUnidade()]) } });
    renderizarRota('/administracao/unidades');

    const linha = await screen.findByRole('row', { name: /Matriz Contagem/ });
    expect(within(linha).getByText('MTZ')).toBeInTheDocument();
    expect(within(linha).getByText('Contagem, Minas Gerais')).toBeInTheDocument();
    expect(within(linha).getByText('Ativa')).toBeInTheDocument();
  });

  it('vazio oferece criar a primeira unidade', async () => {
    simularAPI({ ...BASE, 'GET /v1/unidades': { corpo: pagina([]) } });
    renderizarRota('/administracao/unidades');
    expect(await screen.findByRole('heading', { name: 'Nenhuma unidade cadastrada' })).toBeInTheDocument();
    expect(screen.getAllByRole('button', { name: 'Nova unidade' }).length).toBeGreaterThan(0);
  });

  it('sem adm.gerenciar_empresa, não oferece criar', async () => {
    simularAPI({ ...rotaDaSessao(sessao({ permissoes: ['adm.ler'] })), 'GET /v1/unidades': { corpo: pagina([resumoUnidade()]) } });
    renderizarRota('/administracao/unidades');
    await screen.findByRole('row', { name: /Matriz/ });
    expect(screen.queryByRole('button', { name: 'Nova unidade' })).not.toBeInTheDocument();
  });
});

describe('nova unidade', () => {
  it('exige turno e mostra o erro de plano da segunda unidade (403)', async () => {
    const api = simularAPI({ ...BASE, 'POST /v1/unidades': problema(403, 'plano_nao_inclui_recurso', { plano_necessario: 'profissional' }) });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/unidades/nova');

    await usuario.type(await screen.findByLabelText('Código'), 'FIL');
    await usuario.type(screen.getByLabelText('Nome', { selector: 'input[name="nome"]' }), 'Filial Betim');
    await usuario.click(screen.getByRole('button', { name: 'Criar unidade' }));
    expect((await screen.findAllByText(/Informe ao menos um turno/)).length).toBeGreaterThan(0);
    expect(api.chamadasA('POST /v1/unidades')).toHaveLength(0);

    await usuario.click(screen.getByRole('button', { name: 'Contínuo 24x7' }));
    await usuario.click(screen.getByRole('button', { name: 'Criar unidade' }));

    expect(await screen.findByRole('alert')).toHaveTextContent('O plano da empresa não inclui este recurso');
    const corpo = api.chamadasA('POST /v1/unidades')[0]?.corpo as { turnos: unknown[]; codigo: string };
    expect(corpo.codigo).toBe('FIL');
    expect(corpo.turnos).toHaveLength(7);
  });

  it('código repetido (422 do servidor) aparece no campo', async () => {
    simularAPI({
      ...BASE,
      'POST /v1/unidades': problema(422, 'validacao_falhou', {
        erros: [{ campo: 'codigo', codigo: 'duplicado', mensagem: 'Já existe uma unidade com este código.' }],
      }),
    });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/unidades/nova');

    await usuario.type(await screen.findByLabelText('Código'), 'MTZ');
    await usuario.type(screen.getByLabelText('Nome', { selector: 'input[name="nome"]' }), 'Outra');
    await usuario.click(screen.getByRole('button', { name: 'Comercial' }));
    await usuario.click(screen.getByRole('button', { name: 'Criar unidade' }));

    await waitFor(() => expect(screen.getByLabelText('Código')).toHaveAccessibleDescription(/Já existe uma unidade com este código/));
  });
});

describe('endereço da unidade por país (subdivisão e código postal)', () => {
  const ARGENTINA = {
    ...BASE,
    'GET /v1/usuarios/eu': { corpo: usuarioDeTeste({ id: 'u-1', moeda: 'ARS', pais: 'AR', idioma: 'pt-BR', versao: 1 }) },
  };

  async function preencherBasico(usuario: ReturnType<typeof userEvent.setup>) {
    await usuario.type(await screen.findByLabelText('Código'), 'BA1');
    await usuario.type(screen.getByLabelText('Nome', { selector: 'input[name="nome"]' }), 'Depósito Buenos Aires');
    await usuario.click(screen.getByRole('button', { name: 'Comercial' }));
  }

  it('Argentina: o seletor é "Província" com o nome oficial, e o CPA vai normalizado', async () => {
    const api = simularAPI({ ...ARGENTINA, 'POST /v1/unidades': { status: 201, corpo: unidade({ id: 'un-9' }) }, 'GET /v1/unidades/un-9': { corpo: unidade({ id: 'un-9' }) } });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/unidades/nova');
    await preencherBasico(usuario);

    const provincia = await screen.findByLabelText('Província (opcional)');
    expect(within(provincia).getAllByRole('option')).toHaveLength(25);
    await usuario.selectOptions(provincia, 'Ciudad Autónoma de Buenos Aires');
    await usuario.type(screen.getByLabelText('Código postal (opcional)'), 'c1425 dkf');
    await usuario.click(screen.getByRole('button', { name: 'Criar unidade' }));

    await waitFor(() => expect(api.chamadasA('POST /v1/unidades')).toHaveLength(1));
    expect(api.chamadasA('POST /v1/unidades')[0]?.corpo).toMatchObject({
      endereco: { subdivisao: 'AR-C', codigo_postal: 'C1425DKF' },
    });
  });

  it('Argentina: a letra do CPA tem de ser a da província, sem enviar', async () => {
    const api = simularAPI({ ...ARGENTINA, 'POST /v1/unidades': { status: 201, corpo: unidade() } });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/unidades/nova');
    await preencherBasico(usuario);

    await usuario.selectOptions(await screen.findByLabelText('Província (opcional)'), 'Buenos Aires');
    await usuario.type(screen.getByLabelText('Código postal (opcional)'), 'C1425DKF');
    await usuario.click(screen.getByRole('button', { name: 'Criar unidade' }));

    await waitFor(() =>
      expect(screen.getByLabelText('Código postal (opcional)')).toHaveAccessibleDescription(/letra inicial do CPA/),
    );
    expect(api.chamadasA('POST /v1/unidades')).toHaveLength(0);
  });

  it('Brasil: "Estado", CEP de 8 dígitos e erro do servidor em endereco.subdivisao no campo', async () => {
    simularAPI({
      ...BASE,
      'POST /v1/unidades': problema(422, 'validacao_falhou', {
        erros: [{ campo: 'endereco.subdivisao', codigo: 'valor_invalido', mensagem: 'Subdivisão inválida para o país.' }],
      }),
    });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/unidades/nova');
    await preencherBasico(usuario);

    await usuario.type(screen.getByLabelText('CEP (opcional)'), '1234');
    await usuario.click(screen.getByRole('button', { name: 'Criar unidade' }));
    await waitFor(() => expect(screen.getByLabelText('CEP (opcional)')).toHaveAccessibleDescription(/8 dígitos/));

    await usuario.clear(screen.getByLabelText('CEP (opcional)'));
    await usuario.selectOptions(screen.getByLabelText('Estado (opcional)'), 'São Paulo');
    await usuario.click(screen.getByRole('button', { name: 'Criar unidade' }));
    await waitFor(() => expect(screen.getByLabelText('Estado (opcional)')).toHaveAccessibleDescription(/Subdivisão inválida/));
  });
});

describe('detalhe da unidade', () => {
  const ROTAS = () => ({
    ...BASE,
    'GET /v1/unidades/un-1': { corpo: unidade() },
    'GET /v1/unidades/un-1/feriados': {
      corpo: pagina([{ id: 'f-1', unidade_id: 'un-1', data: '2026-12-25', descricao: 'Natal', versao: 6 }]),
    },
  });

  it('lista turnos por dia e feriados com data dd/mm/aaaa', async () => {
    simularAPI(ROTAS());
    renderizarRota('/administracao/unidades/un-1');

    expect(await screen.findByRole('heading', { level: 1, name: 'Matriz Contagem' })).toBeInTheDocument();
    expect(screen.getByRole('row', { name: /Segunda-feira.*08:00 às 17:00/ })).toBeInTheDocument();
    expect(await screen.findByRole('row', { name: /25\/12\/2026.*Natal/ })).toBeInTheDocument();
  });

  it('inclui um turno; a sobreposição recusada pelo servidor aparece no modal', async () => {
    const api = simularAPI({
      ...ROTAS(),
      'POST /v1/unidades/un-1/turnos': problema(422, 'validacao_falhou', {
        erros: [{ campo: 'hora_inicio', codigo: 'valor_invalido', mensagem: 'O turno se sobrepõe ao de segunda, das 08:00 às 17:00.' }],
      }),
    });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/unidades/un-1');

    await usuario.click(await screen.findByRole('button', { name: 'Adicionar turno' }));
    const dialogo = await screen.findByRole('dialog', { name: 'Novo turno' });
    await usuario.type(within(dialogo).getByLabelText('Início'), '16:00');
    await usuario.type(within(dialogo).getByLabelText('Fim'), '20:00');
    await usuario.click(within(dialogo).getByRole('button', { name: 'Salvar' }));

    await waitFor(() => expect(within(dialogo).getByLabelText('Início')).toHaveAccessibleDescription(/se sobrepõe ao de segunda/));
    expect(api.chamadasA('POST /v1/unidades/un-1/turnos')[0]?.corpo).toEqual({ dia_semana: 1, hora_inicio: '16:00', hora_fim: '20:00' });
  });

  it('altera um turno com o If-Match da versão dele, não da unidade', async () => {
    const api = simularAPI({
      ...ROTAS(),
      'PATCH /v1/unidades/un-1/turnos/t-2': { corpo: { ...unidade().turnos?.[1], hora_fim: '18:00', versao: 5 } },
    });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/unidades/un-1');

    await usuario.click(await screen.findByRole('button', { name: 'Editar turno de Terça-feira, 08:00 às 17:00' }));
    const dialogo = await screen.findByRole('dialog', { name: 'Editar turno' });
    const fim = within(dialogo).getByLabelText('Fim');
    await usuario.clear(fim);
    await usuario.type(fim, '18:00');
    await usuario.click(within(dialogo).getByRole('button', { name: 'Salvar' }));

    await waitFor(() => expect(api.chamadasA('PATCH /v1/unidades/un-1/turnos/t-2')).toHaveLength(1));
    expect(api.chamadasA('PATCH /v1/unidades/un-1/turnos/t-2')[0]?.cabecalhos.get('If-Match')).toBe('"4"');
  });

  it('remove o último turno: o servidor recusa e a mensagem aparece na confirmação', async () => {
    simularAPI({
      ...ROTAS(),
      'DELETE /v1/unidades/un-1/turnos/t-1': problema(422, 'validacao_falhou', {
        erros: [{ campo: 'turnos', codigo: 'tamanho_minimo', mensagem: 'A unidade precisa de pelo menos um turno de trabalho (RN-ADM-019).' }],
      }),
    });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/unidades/un-1');

    await usuario.click(await screen.findByRole('button', { name: 'Remover turno de Segunda-feira, 08:00 às 17:00' }));
    const dialogo = await screen.findByRole('alertdialog', { name: 'Remover este turno?' });
    await usuario.click(within(dialogo).getByRole('button', { name: 'Remover turno' }));

    expect(await within(dialogo).findByRole('alert')).toHaveTextContent('pelo menos um turno');
  });

  it('inclui feriado digitado como dd/mm/aaaa e envia a data em ISO', async () => {
    const api = simularAPI({
      ...ROTAS(),
      'POST /v1/unidades/un-1/feriados': { status: 201, corpo: { id: 'f-2', unidade_id: 'un-1', data: '2026-11-20', descricao: 'Consciência Negra', versao: 1 } },
    });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/unidades/un-1');

    await usuario.click(await screen.findByRole('button', { name: 'Adicionar feriado' }));
    const dialogo = await screen.findByRole('dialog', { name: 'Novo feriado' });
    await usuario.type(within(dialogo).getByLabelText('Data'), '20112026');
    expect(within(dialogo).getByLabelText('Data')).toHaveValue('20/11/2026');
    await usuario.type(within(dialogo).getByLabelText('Descrição'), 'Consciência Negra');
    await usuario.click(within(dialogo).getByRole('button', { name: 'Incluir feriado' }));

    await waitFor(() => expect(api.chamadasA('POST /v1/unidades/un-1/feriados')).toHaveLength(1));
    expect(api.chamadasA('POST /v1/unidades/un-1/feriados')[0]?.corpo).toEqual({ data: '2026-11-20', descricao: 'Consciência Negra' });
  });

  it('recusa data inexistente no campo', async () => {
    const api = simularAPI(ROTAS());
    const usuario = userEvent.setup();
    renderizarRota('/administracao/unidades/un-1');

    await usuario.click(await screen.findByRole('button', { name: 'Adicionar feriado' }));
    const dialogo = await screen.findByRole('dialog', { name: 'Novo feriado' });
    await usuario.type(within(dialogo).getByLabelText('Data'), '31022026');
    await usuario.type(within(dialogo).getByLabelText('Descrição'), 'Teste');
    await usuario.click(within(dialogo).getByRole('button', { name: 'Incluir feriado' }));

    await waitFor(() => expect(within(dialogo).getByLabelText('Data')).toHaveAccessibleDescription(/Data inválida/));
    expect(api.chamadasA('POST /v1/unidades/un-1/feriados')).toHaveLength(0);
  });

  it('remove feriado com o If-Match do feriado', async () => {
    const api = simularAPI({ ...ROTAS(), 'DELETE /v1/unidades/un-1/feriados/f-1': { status: 204 } });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/unidades/un-1');

    await usuario.click(await screen.findByRole('button', { name: 'Remover Natal, 25/12/2026' }));
    await usuario.click(within(await screen.findByRole('alertdialog')).getByRole('button', { name: 'Remover feriado' }));

    await waitFor(() => expect(api.chamadasA('DELETE /v1/unidades/un-1/feriados/f-1')).toHaveLength(1));
    expect(api.chamadasA('DELETE /v1/unidades/un-1/feriados/f-1')[0]?.cabecalhos.get('If-Match')).toBe('"6"');
  });

  it('edita os dados mandando só o que mudou, com If-Match da unidade', async () => {
    const api = simularAPI({ ...ROTAS(), 'PATCH /v1/unidades/un-1': { corpo: unidade({ nome: 'Matriz BH', versao: 3 }) } });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/unidades/un-1');

    const nome = await screen.findByLabelText('Nome', { selector: 'input[name="nome"]' });
    await usuario.clear(nome);
    await usuario.type(nome, 'Matriz BH');
    await usuario.click(screen.getAllByRole('button', { name: 'Salvar' })[0] as HTMLElement);

    await waitFor(() => expect(api.chamadasA('PATCH /v1/unidades/un-1')).toHaveLength(1));
    const [chamada] = api.chamadasA('PATCH /v1/unidades/un-1');
    expect(chamada?.corpo).toEqual({ nome: 'Matriz BH' });
    expect(chamada?.cabecalhos.get('If-Match')).toBe('"2"');
  });

  it('inativa com confirmação; a unidade inativa fica só para leitura', async () => {
    const api = simularAPI({ ...ROTAS(), 'DELETE /v1/unidades/un-1': { status: 204 } });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/unidades/un-1');

    await usuario.click(await screen.findByRole('button', { name: 'Inativar unidade' }));
    const dialogo = await screen.findByRole('alertdialog', { name: 'Inativar Matriz Contagem?' });
    api.definir('GET /v1/unidades/un-1', { corpo: unidade({ situacao: 'inativa', versao: 3 }) });
    await usuario.click(within(dialogo).getByRole('button', { name: 'Inativar unidade' }));

    expect(await screen.findByRole('button', { name: 'Reativar unidade' })).toBeInTheDocument();
    expect(api.chamadasA('DELETE /v1/unidades/un-1')[0]?.cabecalhos.get('If-Match')).toBe('"2"');
    expect(screen.queryByRole('button', { name: 'Adicionar turno' })).not.toBeInTheDocument();
  });

  it('unidade de outra empresa ou inexistente (404) mostra não encontrada', async () => {
    simularAPI({ ...BASE, 'GET /v1/unidades/x': problema(404, 'recurso_nao_encontrado') });
    renderizarRota('/administracao/unidades/x');
    expect(await screen.findByRole('heading', { name: 'Unidade não encontrada' })).toBeInTheDocument();
  });
});
