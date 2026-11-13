// Cadastro da empresa pelo site (RF-ADM-001): turnos informados pelo
// responsável, sem padrão imposto; os modelos só preenchem o formulário.

import { screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it } from 'vitest';

import { problema, simularAPI } from '../../../testes/api-simulada';
import { renderizarRota } from '../../../testes/renderizar';

type Usuario = ReturnType<typeof userEvent.setup>;

function dia(nome: string): HTMLElement {
  return screen.getByRole('group', { name: nome });
}

async function preencherDadosBasicos(usuario: Usuario) {
  await usuario.type(await screen.findByLabelText('Razão social'), 'Refrigeração Polar Ltda');
  await usuario.type(screen.getByLabelText('CNPJ'), '11.222.333/0001-81');
  await usuario.selectOptions(screen.getByLabelText('Plano'), 'essencial');
  await usuario.type(screen.getByLabelText('Nome do responsável'), 'Ana Souza');
  await usuario.type(screen.getByLabelText('E-mail do responsável'), 'ana@polar.com.br');
  await usuario.type(screen.getByLabelText('Código'), 'MTZ');
  await usuario.type(screen.getByLabelText('Nome', { selector: 'input[name="unidade_nome"]' }), 'Matriz Contagem');
}

describe('criar conta', () => {
  it('nasce sem turnos: nenhum calendário é imposto', async () => {
    renderizarRota('/criar-conta');
    await screen.findByLabelText('Razão social');
    for (const nome of ['Segunda-feira', 'Domingo']) {
      expect(within(dia(nome)).getByText('Sem turno neste dia.')).toBeInTheDocument();
    }
  });

  it('o modelo comercial preenche segunda a sexta e continua editável', async () => {
    const usuario = userEvent.setup();
    renderizarRota('/criar-conta');

    await usuario.click(await screen.findByRole('button', { name: 'Comercial' }));

    for (const nome of ['Segunda-feira', 'Terça-feira', 'Quarta-feira', 'Quinta-feira', 'Sexta-feira']) {
      expect(within(dia(nome)).getByLabelText('Início')).toHaveValue('08:00');
      expect(within(dia(nome)).getByLabelText('Fim')).toHaveValue('17:00');
    }
    expect(within(dia('Sábado')).getByText('Sem turno neste dia.')).toBeInTheDocument();

    const inicio = within(dia('Segunda-feira')).getByLabelText('Início');
    await usuario.clear(inicio);
    await usuario.type(inicio, '07:30');
    expect(inicio).toHaveValue('07:30');
  });

  it('o modelo contínuo preenche os sete dias de 00:00 a 24:00', async () => {
    const usuario = userEvent.setup();
    renderizarRota('/criar-conta');

    await usuario.click(await screen.findByRole('button', { name: 'Contínuo 24x7' }));

    for (const nome of ['Segunda-feira', 'Sábado', 'Domingo']) {
      expect(within(dia(nome)).getByLabelText('Início')).toHaveValue('00:00');
      expect(within(dia(nome)).getByLabelText('Fim')).toHaveValue('24:00');
    }
  });

  it('recusa conta sem turno e turnos sobrepostos, com o erro no campo e sem chamar a API', async () => {
    const api = simularAPI({});
    const usuario = userEvent.setup();
    renderizarRota('/criar-conta');
    await preencherDadosBasicos(usuario);

    await usuario.click(screen.getByRole('button', { name: 'Criar conta' }));
    expect(await screen.findAllByText('Informe ao menos um turno de trabalho. Use um modelo para começar.')).not.toHaveLength(0);

    await usuario.click(screen.getByRole('button', { name: 'Comercial' }));
    await usuario.click(screen.getByRole('button', { name: 'Adicionar turno de Segunda-feira' }));
    const segunda = dia('Segunda-feira');
    await usuario.type(within(segunda).getAllByLabelText('Início')[1] as HTMLElement, '16:00');
    await usuario.type(within(segunda).getAllByLabelText('Fim')[1] as HTMLElement, '20:00');
    await usuario.click(screen.getByRole('button', { name: 'Criar conta' }));

    expect(await within(segunda).findByText('Este turno se sobrepõe a outro da unidade. Ajuste os horários.')).toBeInTheDocument();
    expect(api.chamadas).toHaveLength(0);
  });

  it('valida CNPJ e horário antes de enviar', async () => {
    const usuario = userEvent.setup();
    renderizarRota('/criar-conta');
    await usuario.type(await screen.findByLabelText('CNPJ'), '11.222.333/0001-82');
    await usuario.click(screen.getByRole('button', { name: 'Adicionar turno de Segunda-feira' }));
    await usuario.type(within(dia('Segunda-feira')).getByLabelText('Início'), '25:00');
    await usuario.click(screen.getByRole('button', { name: 'Criar conta' }));

    expect(await screen.findByLabelText('CNPJ')).toHaveAccessibleDescription(/CNPJ inválido\. Confira os 14 caracteres\./);
    expect(within(dia('Segunda-feira')).getByLabelText('Início')).toHaveAccessibleDescription('Use HH:MM, de 00:00 a 23:59.');
  });

  it('envia o corpo do contrato e avisa do convite por e-mail', async () => {
    const api = simularAPI({ 'POST /v1/empresas': { status: 201, corpo: { convite_enviado: true } } });
    const usuario = userEvent.setup();
    renderizarRota('/criar-conta');
    await preencherDadosBasicos(usuario);
    await usuario.click(screen.getByRole('button', { name: 'Comercial' }));

    await usuario.click(screen.getByRole('button', { name: 'Criar conta' }));

    expect(await screen.findByRole('status')).toHaveTextContent('Enviamos o convite para ana@polar.com.br');
    const [chamada] = api.chamadasA('POST /v1/empresas');
    expect(chamada?.corpo).toMatchObject({
      pais: 'BR',
      identificacao_tipo: 'cnpj',
      identificacao: '11222333000181',
      idioma_padrao: 'pt-BR',
      razao_social: 'Refrigeração Polar Ltda',
      fuso_horario: 'America/Sao_Paulo',
      plano: 'essencial',
      responsavel: { nome: 'Ana Souza', email: 'ana@polar.com.br' },
      unidade: { codigo: 'MTZ', nome: 'Matriz Contagem' },
    });
    const unidade = (chamada?.corpo as { unidade: { turnos: unknown[] } }).unidade;
    expect(unidade.turnos).toHaveLength(5);
    expect(unidade.turnos[0]).toEqual({ dia_semana: 1, hora_inicio: '08:00', hora_fim: '17:00' });
    expect(chamada?.cabecalhos.get('Idempotency-Key')).toBeTruthy();
  });

  it('o país define o tipo de identificação, o idioma padrão e o fuso: Argentina pede CUIT', async () => {
    const api = simularAPI({ 'POST /v1/empresas': { status: 201, corpo: { convite_enviado: true } } });
    const usuario = userEvent.setup();
    renderizarRota('/criar-conta');
    await usuario.type(await screen.findByLabelText('Razão social'), 'Frigorífico Pampa SA');
    expect(screen.getByLabelText('País')).toHaveValue('BR');
    expect(screen.getByLabelText('Idioma padrão')).toHaveValue('pt-BR');

    await usuario.selectOptions(screen.getByLabelText('País'), 'AR');

    expect(screen.queryByLabelText('CNPJ')).not.toBeInTheDocument();
    expect(screen.getByLabelText('Idioma padrão')).toHaveValue('es-419');
    expect(screen.getByLabelText('Fuso horário')).toHaveValue('America/Argentina/Buenos_Aires');
    const cuit = screen.getByLabelText('CUIT');
    await usuario.type(cuit, '33693450239');
    expect(cuit).toHaveValue('33-69345023-9');

    await usuario.type(screen.getByLabelText('Nome do responsável'), 'Lucía Fernández');
    await usuario.type(screen.getByLabelText('E-mail do responsável'), 'lucia@pampa.com.ar');
    await usuario.type(screen.getByLabelText('Código'), 'BA1');
    await usuario.type(screen.getByLabelText('Nome', { selector: 'input[name="unidade_nome"]' }), 'Planta Buenos Aires');
    await usuario.selectOptions(screen.getByLabelText('Plano'), 'essencial');
    await usuario.click(screen.getByRole('button', { name: 'Comercial' }));
    await usuario.click(screen.getByRole('button', { name: 'Criar conta' }));

    await screen.findByRole('status');
    expect(api.chamadasA('POST /v1/empresas')[0]?.corpo).toMatchObject({
      pais: 'AR',
      identificacao_tipo: 'cuit',
      identificacao: '33693450239',
      idioma_padrao: 'es-419',
      fuso_horario: 'America/Argentina/Buenos_Aires',
    });
  });

  it('o idioma padrão escolhido à mão não muda quando se troca o país', async () => {
    const usuario = userEvent.setup();
    renderizarRota('/criar-conta');
    await usuario.selectOptions(await screen.findByLabelText('Idioma padrão'), 'es-419');
    await usuario.selectOptions(screen.getByLabelText('País'), 'AR');
    await usuario.selectOptions(screen.getByLabelText('País'), 'BR');
    expect(screen.getByLabelText('Idioma padrão')).toHaveValue('es-419');
  });

  it('CNPJ alfanumérico: aceita letras em minúsculas, mascara e envia em maiúsculas', async () => {
    const api = simularAPI({ 'POST /v1/empresas': { status: 201, corpo: { convite_enviado: true } } });
    const usuario = userEvent.setup();
    renderizarRota('/criar-conta');
    await preencherDadosBasicos(usuario);
    const campo = screen.getByLabelText('CNPJ');
    await usuario.clear(campo);
    await usuario.type(campo, '12abc34501de35');
    expect(campo).toHaveValue('12.ABC.345/01DE-35');
    await usuario.click(screen.getByRole('button', { name: 'Comercial' }));
    await usuario.click(screen.getByRole('button', { name: 'Criar conta' }));

    await screen.findByRole('status');
    expect(api.chamadasA('POST /v1/empresas')[0]?.corpo).toMatchObject({ identificacao: '12ABC34501DE35' });
  });

  it('mostra a recusa de sobreposição do servidor no turno apontado', async () => {
    simularAPI({
      'POST /v1/empresas': problema(422, 'validacao_falhou', {
        erros: [{ campo: 'unidade.turnos[1]', codigo: 'valor_invalido', mensagem: 'O turno se sobrepõe ao de segunda, das 08:00 às 17:00.' }],
      }),
    });
    const usuario = userEvent.setup();
    renderizarRota('/criar-conta');
    await preencherDadosBasicos(usuario);
    await usuario.click(screen.getByRole('button', { name: 'Comercial' }));
    await usuario.click(screen.getByRole('button', { name: 'Criar conta' }));

    // turnos[1] é a terça no modelo comercial: o erro aparece na linha dela.
    expect(await within(dia('Terça-feira')).findByText('O turno se sobrepõe ao de segunda, das 08:00 às 17:00.')).toBeInTheDocument();
  });

  it('CNPJ já cadastrado (409) aparece como erro geral', async () => {
    simularAPI({ 'POST /v1/empresas': problema(409, 'conflito') });
    const usuario = userEvent.setup();
    renderizarRota('/criar-conta');
    await preencherDadosBasicos(usuario);
    await usuario.click(screen.getByRole('button', { name: 'Comercial' }));
    await usuario.click(screen.getByRole('button', { name: 'Criar conta' }));

    expect(await screen.findByRole('alert')).toHaveTextContent('Já existe um registro com estes dados');
    await waitFor(() => expect(screen.getByRole('button', { name: 'Criar conta' })).toBeEnabled());
  });

  it('avisa quando a conta foi criada, mas o e-mail do convite não saiu', async () => {
    simularAPI({ 'POST /v1/empresas': { status: 201, corpo: { convite_enviado: false } } });
    const usuario = userEvent.setup();
    renderizarRota('/criar-conta');
    await preencherDadosBasicos(usuario);
    await usuario.click(screen.getByRole('button', { name: 'Comercial' }));
    await usuario.click(screen.getByRole('button', { name: 'Criar conta' }));

    expect(await screen.findByRole('status')).toHaveTextContent('não saiu');
  });
});
