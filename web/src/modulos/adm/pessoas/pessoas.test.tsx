// Pessoas (RF-ADM-015): identificação mascarada na lista e completa no detalhe, aviso
// fixo de dado de saúde sem botão de fechar, paginação por cursor.

import { screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it } from 'vitest';

import { problema, simularAPI } from '../../../testes/api-simulada';
import { pagina, pessoa, resumoPessoa, resumoUnidade, sessao, unidade, usuario as usuarioDeTeste } from '../../../testes/dados';
import { renderizarRota, rotaDaSessao } from '../../../testes/renderizar';

const AVISO = 'A Tenaz não guarda dado de saúde. Não registre nem anexe ASO, exames, atestados ou qualquer informação médica.';
const BASE = {
  ...rotaDaSessao(sessao()),
  'GET /v1/unidades': { corpo: pagina([resumoUnidade(unidade())]) },
};

describe('lista de pessoas', () => {
  it('mostra o CPF mascarado e nunca o completo', async () => {
    simularAPI({ ...BASE, 'GET /v1/pessoas': { corpo: pagina([resumoPessoa()]) } });
    renderizarRota('/administracao/pessoas');

    const linha = await screen.findByRole('row', { name: /João Lima/ });
    expect(within(linha).getByText('CPF ***.982.247-**')).toBeInTheDocument();
    expect(screen.getByRole('columnheader', { name: 'Documento' })).toBeInTheDocument();
    expect(document.body).not.toHaveTextContent('529.982.247-25');
  });

  it('busca no servidor e guarda o filtro na URL', async () => {
    const api = simularAPI({ ...BASE, 'GET /v1/pessoas': { corpo: pagina([resumoPessoa()]) } });
    const usuario = userEvent.setup();
    const { roteador } = renderizarRota('/administracao/pessoas');

    await usuario.type(await screen.findByLabelText('Buscar por nome'), 'joão{Enter}');

    await waitFor(() => expect(api.chamadasA('GET /v1/pessoas').at(-1)?.consulta.get('busca')).toBe('joão'));
    expect(roteador.state.location.search).toContain('busca=jo');
  });

  it('pagina por cursor: a segunda página usa o proximo_cursor', async () => {
    const api = simularAPI({
      ...BASE,
      'GET /v1/pessoas': (c) =>
        c.consulta.get('cursor') === 'CURSOR-2'
          ? { corpo: pagina([resumoPessoa(pessoa({ id: 'pe-2', nome: 'Maria Dias' }))]) }
          : { corpo: pagina([resumoPessoa()], 'CURSOR-2') },
    });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/pessoas');

    expect(await screen.findByText('1 item carregado')).toBeInTheDocument();
    await usuario.click(screen.getByRole('button', { name: 'Carregar mais' }));

    expect(await screen.findByRole('row', { name: /Maria Dias/ })).toBeInTheDocument();
    expect(screen.getByRole('row', { name: /João Lima/ })).toBeInTheDocument();
    expect(screen.getByText('Fim da lista · 2 itens')).toBeInTheDocument();
    expect(api.chamadasA('GET /v1/pessoas')[0]?.consulta.get('limite')).toBe('50');
  });

  it('filtro sem resultado explica e oferece limpar', async () => {
    simularAPI({ ...BASE, 'GET /v1/pessoas': { corpo: pagina([]) } });
    renderizarRota('/administracao/pessoas?situacao=inativa');
    expect(await screen.findByRole('heading', { name: 'Nenhum resultado' })).toBeInTheDocument();
    expect(screen.getByRole('button', { name: 'Limpar filtros' })).toBeInTheDocument();
  });

  it('erro de carga oferece tentar de novo', async () => {
    simularAPI({ ...BASE, 'GET /v1/pessoas': problema(503, 'indisponivel') });
    renderizarRota('/administracao/pessoas');
    expect(await screen.findByRole('button', { name: 'Tentar de novo' })).toBeInTheDocument();
  });
});

describe('aviso de dado de saúde', () => {
  it.each([
    ['criar', '/administracao/pessoas/nova', {}],
    ['detalhar e editar', '/administracao/pessoas/pe-1', { 'GET /v1/pessoas/pe-1': { corpo: pessoa() } }],
  ])('aparece ao %s, sem botão de fechar', async (_nome, caminho, extras) => {
    simularAPI({ ...BASE, ...extras });
    renderizarRota(caminho);

    const aviso = await screen.findByText(AVISO);
    const caixa = aviso.closest('[role="status"]') as HTMLElement;
    expect(caixa).not.toBeNull();
    expect(within(caixa).queryByRole('button')).not.toBeInTheDocument();
  });

  it('aparece também enquanto o detalhe carrega e para quem só pode ler', async () => {
    simularAPI({ ...rotaDaSessao(sessao({ permissoes: ['adm.ler_pessoas'] })), ...{ 'GET /v1/unidades': { corpo: pagina([]) } }, 'GET /v1/pessoas/pe-1': { corpo: pessoa() } });
    renderizarRota('/administracao/pessoas/pe-1');
    expect(await screen.findByText(AVISO)).toBeInTheDocument();
    expect(await screen.findByLabelText('Nome')).toHaveAttribute('readonly');
  });
});

describe('cadastro de pessoa', () => {
  it('cria com CPF só em dígitos e vai para o detalhe', async () => {
    const api = simularAPI({
      ...BASE,
      'POST /v1/pessoas': { status: 201, corpo: pessoa({ id: 'pe-9', nome: 'Rita Alves' }) },
      'GET /v1/pessoas/pe-9': { corpo: pessoa({ id: 'pe-9', nome: 'Rita Alves' }) },
    });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/pessoas/nova');

    await usuario.type(await screen.findByLabelText('Nome'), 'Rita Alves');
    await usuario.type(screen.getByLabelText('CPF (opcional)'), '529.982.247-25');
    await usuario.type(screen.getByLabelText('Data de nascimento (opcional)'), '15031990');
    await usuario.selectOptions(screen.getByLabelText('Vínculo'), 'terceiro');
    await usuario.click(screen.getByRole('button', { name: 'Cadastrar pessoa' }));

    expect(await screen.findByRole('heading', { level: 1, name: 'Rita Alves' })).toBeInTheDocument();
    expect(api.chamadasA('POST /v1/pessoas')[0]?.corpo).toEqual({
      nome: 'Rita Alves',
      vinculo: 'terceiro',
      identificacao_tipo: 'cpf',
      identificacao: '52998224725',
      data_nascimento: '1990-03-15',
    });
  });

  it('cadastra sem CPF: pessoa sem login e sem documento é válida', async () => {
    const api = simularAPI({ ...BASE, 'POST /v1/pessoas': { status: 201, corpo: pessoa({ id: 'pe-9' }) }, 'GET /v1/pessoas/pe-9': { corpo: pessoa({ id: 'pe-9' }) } });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/pessoas/nova');
    await usuario.type(await screen.findByLabelText('Nome'), 'Sem Documento');
    await usuario.click(screen.getByRole('button', { name: 'Cadastrar pessoa' }));
    await waitFor(() => expect(api.chamadasA('POST /v1/pessoas')).toHaveLength(1));
    expect(api.chamadasA('POST /v1/pessoas')[0]?.corpo).toEqual({ nome: 'Sem Documento', vinculo: 'empregado' });
  });

  it('recusa nome vazio e CPF inválido no campo, sem enviar', async () => {
    const api = simularAPI(BASE);
    const usuario = userEvent.setup();
    renderizarRota('/administracao/pessoas/nova');

    await usuario.type(await screen.findByLabelText('CPF (opcional)'), '111.111.111-11');
    await usuario.click(screen.getByRole('button', { name: 'Cadastrar pessoa' }));

    await waitFor(() => expect(screen.getByLabelText('Nome')).toHaveAccessibleDescription('Informe o nome.'));
    expect(screen.getByLabelText('CPF (opcional)')).toHaveAccessibleDescription(/CPF inválido/);
    expect(api.chamadasA('POST /v1/pessoas')).toHaveLength(0);
    expect(screen.getByRole('link', { name: 'Nome: Informe o nome.' })).toBeInTheDocument();
  });

  it('CPF já usado na empresa (422 do servidor) aparece no campo', async () => {
    simularAPI({
      ...BASE,
      'POST /v1/pessoas': problema(422, 'validacao_falhou', {
        erros: [{ campo: 'identificacao', codigo: 'duplicado', mensagem: 'Já existe uma pessoa com este CPF.' }],
      }),
    });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/pessoas/nova');

    await usuario.type(await screen.findByLabelText('Nome'), 'Rita');
    await usuario.type(screen.getByLabelText('CPF (opcional)'), '52998224725');
    await usuario.click(screen.getByRole('button', { name: 'Cadastrar pessoa' }));

    await waitFor(() => expect(screen.getByLabelText('CPF (opcional)')).toHaveAccessibleDescription(/Já existe uma pessoa com este CPF/));
  });
});

describe('pessoa de empresa argentina', () => {
  const ARGENTINA = {
    ...BASE,
    'GET /v1/usuarios/eu': { corpo: usuarioDeTeste({ id: 'u-1', moeda: 'ARS', pais: 'AR', idioma: 'pt-BR', versao: 1 }) },
  };

  it('o tipo padrão é o DNI, com máscara de milhares e envio normalizado', async () => {
    const api = simularAPI({ ...ARGENTINA, 'POST /v1/pessoas': { status: 201, corpo: pessoa({ id: 'pe-9' }) }, 'GET /v1/pessoas/pe-9': { corpo: pessoa({ id: 'pe-9' }) } });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/pessoas/nova');

    await usuario.type(await screen.findByLabelText('Nome'), 'Martín Gómez');
    await waitFor(() => expect(screen.getByLabelText('Tipo de documento')).toHaveValue('dni'));
    await usuario.type(screen.getByLabelText('DNI (opcional)'), '12345678');
    expect(screen.getByLabelText('DNI (opcional)')).toHaveValue('12.345.678');
    await usuario.click(screen.getByRole('button', { name: 'Cadastrar pessoa' }));

    await waitFor(() => expect(api.chamadasA('POST /v1/pessoas')).toHaveLength(1));
    expect(api.chamadasA('POST /v1/pessoas')[0]?.corpo).toEqual({
      nome: 'Martín Gómez',
      vinculo: 'empregado',
      identificacao_tipo: 'dni',
      identificacao: '12345678',
    });
  });

  it('aceita os quatro tipos e reescreve o número na máscara do tipo novo', async () => {
    simularAPI(ARGENTINA);
    const usuario = userEvent.setup();
    renderizarRota('/administracao/pessoas/nova');

    await waitFor(() => expect(screen.getByLabelText('Tipo de documento')).toHaveValue('dni'));
    const opcoes = within(screen.getByLabelText('Tipo de documento')).getAllByRole('option').map((o) => o.getAttribute('value'));
    expect(opcoes).toEqual(['dni', 'cpf', 'cuil', 'ci']);

    await usuario.type(screen.getByLabelText('DNI (opcional)'), '20123456786');
    expect(screen.getByLabelText('DNI (opcional)')).toHaveValue('20.123.456');
    await usuario.clear(screen.getByLabelText('DNI (opcional)'));
    await usuario.type(screen.getByLabelText('DNI (opcional)'), '1234567');
    await usuario.selectOptions(screen.getByLabelText('Tipo de documento'), 'ci');
    expect(screen.getByLabelText('CI (opcional)')).toHaveValue('1.234.567');
  });

  it('recusa DNI com tamanho errado no campo, sem enviar', async () => {
    const api = simularAPI(ARGENTINA);
    const usuario = userEvent.setup();
    renderizarRota('/administracao/pessoas/nova');

    await usuario.type(await screen.findByLabelText('Nome'), 'Martín');
    await usuario.type(await screen.findByLabelText('DNI (opcional)'), '123456');
    await usuario.click(screen.getByRole('button', { name: 'Cadastrar pessoa' }));

    await waitFor(() => expect(screen.getByLabelText('DNI (opcional)')).toHaveAccessibleDescription(/DNI inválido/));
    expect(api.chamadasA('POST /v1/pessoas')).toHaveLength(0);
  });

  it('a lista mostra o tipo e a máscara dos 3 últimos caracteres', async () => {
    simularAPI({
      ...ARGENTINA,
      'GET /v1/pessoas': { corpo: pagina([{ ...resumoPessoa(), identificacao_tipo: 'dni', identificacao_mascarada: '*****678' }]) },
    });
    renderizarRota('/administracao/pessoas');
    expect(await screen.findByText('DNI *****678')).toBeInTheDocument();
  });
});

describe('detalhe da pessoa', () => {
  const ROTAS = { ...BASE, 'GET /v1/pessoas/pe-1': { corpo: pessoa() } };

  it('mostra o CPF completo formatado', async () => {
    simularAPI(ROTAS);
    renderizarRota('/administracao/pessoas/pe-1');
    expect(await screen.findByLabelText('CPF (opcional)')).toHaveValue('529.982.247-25');
    expect(screen.getByLabelText('Data de nascimento (opcional)')).toHaveValue('15/03/1990');
  });

  it('edita mandando só o que mudou, com If-Match', async () => {
    const api = simularAPI({ ...ROTAS, 'PATCH /v1/pessoas/pe-1': { corpo: pessoa({ cargo: 'Líder', versao: 6 }) } });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/pessoas/pe-1');

    const cargo = await screen.findByLabelText('Cargo (opcional)');
    await usuario.clear(cargo);
    await usuario.type(cargo, 'Líder');
    await usuario.click(screen.getByRole('button', { name: 'Salvar' }));

    await waitFor(() => expect(api.chamadasA('PATCH /v1/pessoas/pe-1')).toHaveLength(1));
    const [chamada] = api.chamadasA('PATCH /v1/pessoas/pe-1');
    expect(chamada?.corpo).toEqual({ cargo: 'Líder' });
    expect(chamada?.cabecalhos.get('If-Match')).toBe('"5"');
  });

  it('limpar um campo opcional manda null', async () => {
    const api = simularAPI({ ...ROTAS, 'PATCH /v1/pessoas/pe-1': { corpo: pessoa({ cargo: null, versao: 6 }) } });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/pessoas/pe-1');

    await usuario.clear(await screen.findByLabelText('Cargo (opcional)'));
    await usuario.click(screen.getByRole('button', { name: 'Salvar' }));

    await waitFor(() => expect(api.chamadasA('PATCH /v1/pessoas/pe-1')[0]?.corpo).toEqual({ cargo: null }));
  });

  it('conflito de versão (412) avisa e recarrega', async () => {
    const api = simularAPI({ ...ROTAS, 'PATCH /v1/pessoas/pe-1': problema(412, 'versao_desatualizada', { versao_atual: 6 }) });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/pessoas/pe-1');

    const cargo = await screen.findByLabelText('Cargo (opcional)');
    await usuario.type(cargo, ' sênior');
    await usuario.click(screen.getByRole('button', { name: 'Salvar' }));

    expect(await screen.findByText(/alterado por outra pessoa/)).toBeInTheDocument();
    await waitFor(() => expect(api.chamadasA('GET /v1/pessoas/pe-1').length).toBeGreaterThan(1));
  });

  it('inativa com confirmação e If-Match; pessoa inativa pode ser reativada', async () => {
    const api = simularAPI({ ...ROTAS, 'DELETE /v1/pessoas/pe-1': { status: 204 }, 'POST /v1/pessoas/pe-1/reativar': { corpo: pessoa({ versao: 7 }) } });
    const usuario = userEvent.setup();
    renderizarRota('/administracao/pessoas/pe-1');

    await usuario.click(await screen.findByRole('button', { name: 'Inativar pessoa' }));
    const dialogo = await screen.findByRole('alertdialog', { name: 'Inativar João Lima?' });
    api.definir('GET /v1/pessoas/pe-1', { corpo: pessoa({ situacao: 'inativa', versao: 6 }) });
    await usuario.click(within(dialogo).getByRole('button', { name: 'Inativar pessoa' }));

    const reativar = await screen.findByRole('button', { name: 'Reativar pessoa' });
    expect(api.chamadasA('DELETE /v1/pessoas/pe-1')[0]?.cabecalhos.get('If-Match')).toBe('"5"');
    expect(screen.getByLabelText('Nome')).toHaveAttribute('readonly');
    await usuario.click(reativar);
    await waitFor(() => expect(api.chamadasA('POST /v1/pessoas/pe-1/reativar')).toHaveLength(1));
    expect(api.chamadasA('POST /v1/pessoas/pe-1/reativar')[0]?.cabecalhos.get('If-Match')).toBe('"6"');
  });

  it('mostra o link do usuário vinculado e o histórico do registro', async () => {
    simularAPI({ ...BASE, 'GET /v1/pessoas/pe-1': { corpo: pessoa({ usuario_id: 'u-2' }) } });
    renderizarRota('/administracao/pessoas/pe-1');
    expect(await screen.findByRole('link', { name: 'Ver usuário' })).toHaveAttribute('href', '/administracao/usuarios/u-2');
    expect(screen.getByRole('link', { name: 'Histórico do registro' })).toHaveAttribute(
      'href',
      '/administracao/auditoria?entidade=pessoa&entidade_id=pe-1',
    );
  });

  it('pessoa de outra empresa (404) mostra não encontrada', async () => {
    simularAPI({ ...BASE, 'GET /v1/pessoas/zz': problema(404, 'recurso_nao_encontrado') });
    renderizarRota('/administracao/pessoas/zz');
    expect(await screen.findByRole('heading', { name: 'Pessoa não encontrada' })).toBeInTheDocument();
  });
});
