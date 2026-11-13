// Empresa (RF-ADM-001) e minha conta: PATCH só do que mudou, If-Match com a
// versão lida, conflito de versão e o interruptor do MFA obrigatório.

import { screen, waitFor } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it } from 'vitest';

import { problema, simularAPI } from '../../../testes/api-simulada';
import { empresa, sessao } from '../../../testes/dados';
import { renderizarRota, rotaDaSessao } from '../../../testes/renderizar';

describe('empresa', () => {
  it('mostra os dados, com o CNPJ formatado', async () => {
    simularAPI({ ...rotaDaSessao(sessao()), 'GET /v1/empresa': { corpo: empresa() } });
    renderizarRota('/administracao/empresa');

    expect(await screen.findByLabelText('Razão social')).toHaveValue('Refrigeração Polar Ltda');
    expect(screen.getByLabelText('CNPJ')).toHaveValue('11.222.333/0001-81');
  });

  it('país e moeda são só leitura e o idioma padrão é editável', async () => {
    const api = simularAPI({
      ...rotaDaSessao(sessao()),
      'GET /v1/empresa': { corpo: empresa() },
      'PATCH /v1/empresa': { corpo: empresa({ idioma_padrao: 'es-419', versao: 4 }) },
    });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/empresa');

    expect(await screen.findByLabelText('País')).toHaveValue('Brasil');
    expect(screen.getByLabelText('País')).toHaveAttribute('readonly');
    expect(screen.getByLabelText('Moeda')).toHaveValue('Real brasileiro (BRL)');
    expect(screen.getByLabelText('Moeda')).toHaveAttribute('readonly');
    expect(screen.getByLabelText('Idioma padrão')).toHaveValue('pt-BR');

    await usuario.selectOptions(screen.getByLabelText('Idioma padrão'), 'es-419');
    await usuario.click(screen.getByRole('button', { name: 'Salvar' }));

    await waitFor(() => expect(api.chamadasA('PATCH /v1/empresa')).toHaveLength(1));
    expect(api.chamadasA('PATCH /v1/empresa')[0]?.corpo).toEqual({ idioma_padrao: 'es-419' });
  });

  it('empresa paraguaia: RUC com máscara, moeda em guarani e fuso de Assunção', async () => {
    simularAPI({
      ...rotaDaSessao(sessao()),
      'GET /v1/empresa': {
        corpo: empresa({
          pais: 'PY',
          moeda: 'PYG',
          idioma_padrao: 'es-419',
          identificacao_tipo: 'ruc',
          identificacao: '800123456',
          fuso_horario: 'America/Asuncion',
        }),
      },
    });
    renderizarRota('/administracao/empresa');

    expect(await screen.findByLabelText('RUC')).toHaveValue('80012345-6');
    expect(screen.queryByLabelText('CNPJ')).not.toBeInTheDocument();
    expect(screen.getByLabelText('País')).toHaveValue('Paraguai');
    expect(screen.getByLabelText('Moeda')).toHaveValue('Guarani paraguaio (PYG)');
    expect(screen.getByLabelText('Fuso horário')).toHaveValue('America/Asuncion');
  });

  it('envia só o campo alterado, com If-Match da versão lida', async () => {
    const api = simularAPI({
      ...rotaDaSessao(sessao()),
      'GET /v1/empresa': { corpo: empresa() },
      'PATCH /v1/empresa': { corpo: empresa({ nome_fantasia: 'Polar Frio', versao: 4 }) },
    });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/empresa');

    const campo = await screen.findByLabelText('Nome fantasia (opcional)');
    await usuario.clear(campo);
    await usuario.type(campo, 'Polar Frio');
    await usuario.click(screen.getByRole('button', { name: 'Salvar' }));

    expect(await screen.findByRole('status')).toHaveTextContent('Dados da empresa salvos.');
    const [chamada] = api.chamadasA('PATCH /v1/empresa');
    expect(chamada?.corpo).toEqual({ nome_fantasia: 'Polar Frio' });
    expect(chamada?.cabecalhos.get('If-Match')).toBe('"3"');
    expect(chamada?.cabecalhos.get('X-CSRF-Token')).toBe('csrf-de-teste');
  });

  it('a troca de CNPJ vai só com os dígitos', async () => {
    const api = simularAPI({
      ...rotaDaSessao(sessao()),
      'GET /v1/empresa': { corpo: empresa() },
      'PATCH /v1/empresa': { corpo: empresa({ identificacao: '11444777000161', versao: 4 }) },
    });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/empresa');

    const campo = await screen.findByLabelText('CNPJ');
    await usuario.clear(campo);
    await usuario.type(campo, '11.444.777/0001-61');
    await usuario.click(screen.getByRole('button', { name: 'Salvar' }));

    await waitFor(() => expect(api.chamadasA('PATCH /v1/empresa')).toHaveLength(1));
    expect(api.chamadasA('PATCH /v1/empresa')[0]?.corpo).toEqual({ identificacao_tipo: 'cnpj', identificacao: '11444777000161' });
  });

  it('CNPJ com dígito errado é recusado no campo, sem enviar', async () => {
    const api = simularAPI({ ...rotaDaSessao(sessao()), 'GET /v1/empresa': { corpo: empresa() } });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/empresa');

    const campo = await screen.findByLabelText('CNPJ');
    await usuario.clear(campo);
    await usuario.type(campo, '11.222.333/0001-82');
    await usuario.click(screen.getByRole('button', { name: 'Salvar' }));

    await waitFor(() => expect(campo).toHaveAccessibleDescription(/CNPJ inválido/));
    expect(api.chamadasA('PATCH /v1/empresa')).toHaveLength(0);
  });

  it('CNPJ de outra empresa (409/422 do servidor) aparece no campo', async () => {
    simularAPI({
      ...rotaDaSessao(sessao()),
      'GET /v1/empresa': { corpo: empresa() },
      'PATCH /v1/empresa': problema(422, 'validacao_falhou', {
        erros: [{ campo: 'identificacao', codigo: 'duplicado', mensagem: 'Já existe uma empresa com este CNPJ.' }],
      }),
    });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/empresa');

    const campo = await screen.findByLabelText('CNPJ');
    await usuario.clear(campo);
    await usuario.type(campo, '11.444.777/0001-61');
    await usuario.click(screen.getByRole('button', { name: 'Salvar' }));

    await waitFor(() => expect(campo).toHaveAccessibleDescription(/Já existe uma empresa com este CNPJ\./));
  });

  it('versão desatualizada (412) avisa e recarrega a empresa', async () => {
    simularAPI({
      ...rotaDaSessao(sessao()),
      'GET /v1/empresa': { corpo: empresa() },
      'PATCH /v1/empresa': problema(412, 'versao_desatualizada', { versao_atual: 4 }),
    });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/empresa');

    const campo = await screen.findByLabelText('Razão social');
    await usuario.type(campo, ' S.A.');
    await usuario.click(screen.getByRole('button', { name: 'Salvar' }));

    expect(await screen.findByText(/alterado por outra pessoa/)).toBeInTheDocument();
  });

  it('exigir o segundo fator pede confirmação e liga o interruptor no servidor', async () => {
    const api = simularAPI({
      ...rotaDaSessao(sessao()),
      'GET /v1/empresa': { corpo: empresa() },
      'PATCH /v1/empresa': { corpo: empresa({ mfa_obrigatorio: true, versao: 4 }) },
    });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/empresa');

    const chave = await screen.findByRole('switch', { name: 'Exigir segundo fator de todos os usuários' });
    expect(chave).toHaveAttribute('aria-checked', 'false');
    await usuario.click(chave);

    const dialogo = await screen.findByRole('alertdialog', { name: 'Exigir segundo fator de todos os usuários?' });
    expect(dialogo).toHaveTextContent('inclusive você');
    // O que sai sem agir é "Voltar", nunca "Cancelar".
    expect(screen.getByRole('button', { name: 'Voltar' })).toBeInTheDocument();
    await usuario.click(screen.getByRole('button', { name: 'Exigir segundo fator' }));

    await waitFor(() => expect(api.chamadasA('PATCH /v1/empresa')).toHaveLength(1));
    expect(api.chamadasA('PATCH /v1/empresa')[0]?.corpo).toEqual({ mfa_obrigatorio: true });
    await waitFor(() => expect(screen.getByRole('switch', { name: /Exigir segundo fator/ })).toHaveAttribute('aria-checked', 'true'));
  });

  it('desligar o interruptor não pede confirmação', async () => {
    const api = simularAPI({
      ...rotaDaSessao(sessao()),
      'GET /v1/empresa': { corpo: empresa({ mfa_obrigatorio: true }) },
      'PATCH /v1/empresa': { corpo: empresa({ mfa_obrigatorio: false, versao: 4 }) },
    });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/empresa');

    await usuario.click(await screen.findByRole('switch', { name: /Exigir segundo fator/ }));

    await waitFor(() => expect(api.chamadasA('PATCH /v1/empresa')[0]?.corpo).toEqual({ mfa_obrigatorio: false }));
  });

  it('sem adm.gerenciar_empresa, mostra só leitura', async () => {
    simularAPI({ ...rotaDaSessao(sessao({ permissoes: ['adm.ler'] })), 'GET /v1/empresa': { corpo: empresa() } });
    renderizarRota('/administracao/empresa');

    expect(await screen.findByLabelText('Razão social')).toHaveAttribute('readonly');
    expect(screen.queryByRole('button', { name: 'Salvar' })).not.toBeInTheDocument();
    expect(screen.getByRole('switch', { name: /Exigir segundo fator/ })).toBeDisabled();
  });

  it('erro de carga oferece tentar de novo', async () => {
    const api = simularAPI({ ...rotaDaSessao(sessao()), 'GET /v1/empresa': problema(500, 'erro_interno', { id_rastreio: 'E-4F21' }) });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/empresa');

    expect(await screen.findByRole('alert')).toHaveTextContent('E-4F21');
    api.definir('GET /v1/empresa', { corpo: empresa() });
    await usuario.click(screen.getByRole('button', { name: 'Tentar de novo' }));
    expect(await screen.findByLabelText('Razão social')).toBeInTheDocument();
  });
});
