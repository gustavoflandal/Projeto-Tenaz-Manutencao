// Trilha de auditoria (RF-ADM-011): consulta com filtros, exportação em CSV,
// detalhe com antes, depois e hash, verificação de integridade e configuração
// das capturas por tela.

import { screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';

import { problema, simularAPI } from '../../../testes/api-simulada';
import { captura, diaVerificado, evento, pagina, resumoEvento, sessao } from '../../../testes/dados';
import { renderizarRota, rotaDaSessao } from '../../../testes/renderizar';

const BASE = { ...rotaDaSessao(sessao()) };
const EVENTOS = { 'GET /v1/eventos-auditoria': { corpo: pagina([resumoEvento()]) } };

describe('trilha de eventos', () => {
  it('lista data, usuário pelo nome, ação, tela, entidade e IP', async () => {
    simularAPI({ ...BASE, ...EVENTOS });
    renderizarRota('/administracao/auditoria');

    const linha = await screen.findByRole('row', { name: /Alteração/ });
    expect(within(linha).getByText('Ana Souza')).toBeInTheDocument();
    expect(within(linha).getByText('adm.pessoa_detalhe')).toBeInTheDocument();
    expect(within(linha).getByText('pessoa, João Lima')).toBeInTheDocument();
    expect(within(linha).getByText('203.0.113.7')).toBeInTheDocument();
  });

  it('marca falha de login e dado sensível por texto, não só por ícone', async () => {
    simularAPI({
      ...BASE,
      'GET /v1/eventos-auditoria': {
        corpo: pagina([resumoEvento({ acao: 'login', resultado: 'falha', usuario_nome: null }), resumoEvento({ id: 'ev-2', acao: 'visualizacao', sensivel: true })]),
      },
    });
    renderizarRota('/administracao/auditoria');

    const falha = await screen.findByRole('row', { name: /Entrada.*Falha/ });
    expect(within(falha).getByText('Falha')).toBeInTheDocument();
    expect(within(screen.getByRole('row', { name: /Visualização/ })).getByText('Dado sensível')).toBeInTheDocument();
  });

  it('o histórico de um registro chega pela URL e vai para a consulta', async () => {
    const api = simularAPI({ ...BASE, ...EVENTOS });
    renderizarRota('/administracao/auditoria?entidade=pessoa&entidade_id=pe-1&usuario_id=u-1');

    await screen.findByRole('row', { name: /Alteração/ });
    const consulta = api.chamadasA('GET /v1/eventos-auditoria')[0]?.consulta;
    expect(consulta?.get('entidade')).toBe('pessoa');
    expect(consulta?.get('entidade_id')).toBe('pe-1');
    expect(consulta?.get('usuario_id')).toBe('u-1');
    expect(screen.getByText('Registro pe-1')).toBeInTheDocument();
  });

  it('aplica os filtros novos juntos: módulo, tela, ação, período e IP', async () => {
    const api = simularAPI({ ...BASE, ...EVENTOS });
    const u = userEvent.setup();
    renderizarRota('/administracao/auditoria');

    await u.selectOptions(await screen.findByLabelText('Módulo'), 'adm');
    await u.selectOptions(screen.getByLabelText('Tela'), 'adm.pessoas');
    await u.selectOptions(screen.getByLabelText('Ação'), 'visualizacao');
    await u.type(screen.getByLabelText('De'), '01112026');
    await u.type(screen.getByLabelText('Até'), '10112026');
    await u.type(screen.getByLabelText('IP'), '203.0.113.7');
    await u.click(screen.getByRole('button', { name: 'Buscar' }));

    await waitFor(() => {
      const consulta = api.chamadasA('GET /v1/eventos-auditoria').at(-1)?.consulta;
      expect(consulta?.get('modulo')).toBe('adm');
      expect(consulta?.get('tela')).toBe('adm.pessoas');
      expect(consulta?.get('acao')).toBe('visualizacao');
      expect(consulta?.get('ip')).toBe('203.0.113.7');
      expect(consulta?.get('ocorrido_em[gte]')).toBe(new Date(2026, 10, 1).toISOString());
      // O fim é inclusivo na tela e exclusivo na API.
      expect(consulta?.get('ocorrido_em[lt]')).toBe(new Date(2026, 10, 11).toISOString());
    });
  });

  it('o filtro por IP exige a data inicial e não consulta sem ela', async () => {
    const api = simularAPI({ ...BASE, ...EVENTOS });
    const u = userEvent.setup();
    renderizarRota('/administracao/auditoria');
    await screen.findByRole('row', { name: /Alteração/ });
    const antes = api.chamadasA('GET /v1/eventos-auditoria').length;

    await u.type(screen.getByLabelText('IP'), '203.0.113.7');
    await u.click(screen.getByRole('button', { name: 'Buscar' }));

    await waitFor(() => expect(screen.getByLabelText('IP')).toHaveAccessibleDescription(/exige a data inicial/));
    expect(api.chamadasA('GET /v1/eventos-auditoria')).toHaveLength(antes);
  });

  it('período invertido é recusado sem consultar', async () => {
    const api = simularAPI({ ...BASE, ...EVENTOS });
    const u = userEvent.setup();
    renderizarRota('/administracao/auditoria');
    await screen.findByRole('row', { name: /Alteração/ });
    const antes = api.chamadasA('GET /v1/eventos-auditoria').length;

    await u.type(screen.getByLabelText('De'), '10112026');
    await u.type(screen.getByLabelText('Até'), '01112026');
    await u.click(screen.getByRole('button', { name: 'Buscar' }));

    await waitFor(() => expect(screen.getByLabelText('Até')).toHaveAccessibleDescription(/depois do início/));
    expect(api.chamadasA('GET /v1/eventos-auditoria')).toHaveLength(antes);
  });

  it('pagina por cursor', async () => {
    const api = simularAPI({
      ...BASE,
      'GET /v1/eventos-auditoria': (c) =>
        c.consulta.get('cursor') === 'PROX'
          ? { corpo: pagina([resumoEvento({ id: 'ev-2', acao: 'inclusao', entidade_chave: 'Ana' })]) }
          : { corpo: pagina([resumoEvento()], 'PROX') },
    });
    const u = userEvent.setup();
    renderizarRota('/administracao/auditoria');

    await u.click(await screen.findByRole('button', { name: 'Carregar mais' }));

    expect(await screen.findByRole('row', { name: /Inclusão/ })).toBeInTheDocument();
    expect(api.chamadasA('GET /v1/eventos-auditoria')).toHaveLength(2);
  });

  it('o detalhe mostra hash, tela, antes e depois, e o detalhe de uma extração', async () => {
    simularAPI({
      ...BASE,
      ...EVENTOS,
      'GET /v1/eventos-auditoria/ev-1': {
        corpo: evento({ acao: 'extracao', detalhe: { formato: 'csv', linhas: 12 }, antes: null, depois: null }),
      },
    });
    const u = userEvent.setup();
    renderizarRota('/administracao/auditoria');

    await u.click(await screen.findByRole('button', { name: /^Ver detalhes: Alteração em / }));

    const dialogo = await screen.findByRole('dialog', { name: 'Detalhe do evento' });
    expect(await within(dialogo).findByText('ab'.repeat(32))).toBeInTheDocument();
    expect(within(dialogo).getByText('adm.pessoa_detalhe')).toBeInTheDocument();
    expect(within(dialogo).getByText(/"formato": "csv"/)).toBeInTheDocument();
  });

  it('o detalhe de uma alteração mostra os valores antes e depois', async () => {
    simularAPI({ ...BASE, ...EVENTOS, 'GET /v1/eventos-auditoria/ev-1': { corpo: evento() } });
    const u = userEvent.setup();
    renderizarRota('/administracao/auditoria');

    await u.click(await screen.findByRole('button', { name: /^Ver detalhes/ }));

    const tabela = await screen.findByRole('table', { name: 'Valores antes e depois' });
    const linha = within(tabela).getByRole('row', { name: /cargo/ });
    expect(within(linha).getByText('Auxiliar')).toBeInTheDocument();
    expect(within(linha).getByText('Técnico')).toBeInTheDocument();
    expect(within(linha).getByText('alterado')).toBeInTheDocument();
  });

  it('sem eventos, mostra o estado vazio', async () => {
    simularAPI({ ...BASE, 'GET /v1/eventos-auditoria': { corpo: pagina([]) } });
    renderizarRota('/administracao/auditoria');
    expect(await screen.findByRole('heading', { name: 'Nenhum evento' })).toBeInTheDocument();
  });

  it('sem adm.ler_auditoria, a tela nem consulta', async () => {
    const api = simularAPI(rotaDaSessao(sessao({ permissoes: ['adm.ler'] })));
    renderizarRota('/administracao/auditoria');
    expect(await screen.findByRole('heading', { name: 'Sem acesso a esta área' })).toBeInTheDocument();
    expect(api.chamadasA('GET /v1/eventos-auditoria')).toHaveLength(0);
  });

  it('403 do servidor mostra falta de permissão', async () => {
    simularAPI({ ...BASE, 'GET /v1/eventos-auditoria': problema(403, 'permissao_negada') });
    renderizarRota('/administracao/auditoria');
    expect(await screen.findByRole('heading', { name: 'Sem acesso a esta área' })).toBeInTheDocument();
  });
});

describe('exportação em CSV', () => {
  function prepararDownload() {
    const criar = vi.fn(() => 'blob:csv');
    globalThis.URL.createObjectURL = criar;
    globalThis.URL.revokeObjectURL = vi.fn();
    const clicar = vi.spyOn(HTMLAnchorElement.prototype, 'click').mockImplementation(() => undefined);
    return { criar, clicar };
  }

  it('exige o período inteiro e não chama a API sem ele', async () => {
    const api = simularAPI({ ...BASE, ...EVENTOS });
    const u = userEvent.setup();
    renderizarRota('/administracao/auditoria');

    await u.click(await screen.findByRole('button', { name: 'Exportar CSV' }));

    expect(await screen.findByRole('alert')).toHaveTextContent('Informe o período inteiro');
    expect(api.chamadasA('GET /v1/exportacoes-eventos-auditoria')).toHaveLength(0);
  });

  it('baixa o CSV com os mesmos filtros da tela', async () => {
    const { criar, clicar } = prepararDownload();
    const api = simularAPI({
      ...BASE,
      ...EVENTOS,
      'GET /v1/exportacoes-eventos-auditoria': { corpo: undefined, status: 200, cabecalhos: { 'Content-Type': 'text/csv' } },
    });
    const u = userEvent.setup();
    renderizarRota('/administracao/auditoria?acao=extracao&de=2026-11-01&ate=2026-11-10');

    await u.click(await screen.findByRole('button', { name: 'Exportar CSV' }));

    await waitFor(() => expect(clicar).toHaveBeenCalled());
    expect(criar).toHaveBeenCalled();
    const consulta = api.chamadasA('GET /v1/exportacoes-eventos-auditoria')[0]?.consulta;
    expect(consulta?.get('acao')).toBe('extracao');
    expect(consulta?.get('ocorrido_em[gte]')).toBe(new Date(2026, 10, 1).toISOString());
    expect(consulta?.get('ocorrido_em[lt]')).toBe(new Date(2026, 10, 11).toISOString());
    expect(await screen.findByText(/Arquivo CSV gerado/)).toBeInTheDocument();
  });

  it('acima do teto de 50.000 linhas, mostra o motivo do servidor', async () => {
    prepararDownload();
    simularAPI({
      ...BASE,
      ...EVENTOS,
      'GET /v1/exportacoes-eventos-auditoria': problema(422, 'validacao_falhou', {
        erros: [{ campo: 'ocorrido_em[lt]', codigo: 'valor_invalido', mensagem: 'A exportação passa de 50.000 eventos. Reduza o período.' }],
      }),
    });
    const u = userEvent.setup();
    renderizarRota('/administracao/auditoria?de=2026-01-01&ate=2026-11-10');

    await u.click(await screen.findByRole('button', { name: 'Exportar CSV' }));

    expect(await screen.findByRole('alert')).toHaveTextContent('passa de 50.000 eventos');
  });
});

describe('verificação de integridade', () => {
  const ROTA = 'GET /v1/verificacoes-auditoria';

  it('não consulta enquanto ninguém pede, e pede o período em dias UTC', async () => {
    const api = simularAPI({ ...BASE, [ROTA]: { corpo: pagina([diaVerificado()]) } });
    const u = userEvent.setup();
    renderizarRota('/administracao/auditoria/integridade');

    await u.clear(await screen.findByLabelText('Primeiro dia'));
    await u.type(screen.getByLabelText('Primeiro dia'), '10112026');
    await u.clear(screen.getByLabelText('Último dia'));
    await u.type(screen.getByLabelText('Último dia'), '12112026');
    expect(api.chamadasA(ROTA)).toHaveLength(0);
    await u.click(screen.getByRole('button', { name: 'Verificar integridade' }));

    await waitFor(() => expect(api.chamadasA(ROTA)).toHaveLength(1));
    expect(api.chamadasA(ROTA)[0]?.consulta.get('de')).toBe('2026-11-10');
    expect(api.chamadasA(ROTA)[0]?.consulta.get('ate')).toBe('2026-11-12');
  });

  it('mostra íntegro, violado com motivo e evento, e não selado, por texto e ícone', async () => {
    simularAPI({
      ...BASE,
      [ROTA]: {
        corpo: pagina([
          diaVerificado({ dia: '2026-11-10' }),
          diaVerificado({ dia: '2026-11-11', estado: 'violado', motivo: 'evento_alterado', evento_id: 'ev-1' }),
          diaVerificado({ dia: '2026-11-12', estado: 'nao_selado', motivo: 'ainda_nao_selado' }),
        ]),
      },
      'GET /v1/eventos-auditoria/ev-1': { corpo: evento() },
    });
    const u = userEvent.setup();
    renderizarRota('/administracao/auditoria/integridade?de=2026-11-10&ate=2026-11-12');

    const violado = await screen.findByRole('row', { name: /11\/11\/2026/ });
    expect(within(violado).getByText('Violado')).toBeInTheDocument();
    expect(within(violado).getByText('Um evento foi alterado depois de gravado.')).toBeInTheDocument();
    expect(within(screen.getByRole('row', { name: /10\/11\/2026/ })).getByText('Íntegro')).toBeInTheDocument();
    const naoSelado = screen.getByRole('row', { name: /12\/11\/2026/ });
    expect(within(naoSelado).getByText('Não selado')).toBeInTheDocument();
    expect(within(naoSelado).getByText('O dia ainda não foi selado.')).toBeInTheDocument();
    expect(screen.getByRole('alert')).toHaveTextContent('dias violados');

    await u.click(within(violado).getByRole('button', { name: /Ver o evento apontado em 11\/11\/2026/ }));
    expect(await screen.findByRole('dialog', { name: 'Detalhe do evento' })).toBeInTheDocument();
  });

  it('se o armazenamento do selo não responde (503), nada se afirma', async () => {
    simularAPI({ ...BASE, [ROTA]: problema(503, 'servico_indisponivel') });
    renderizarRota('/administracao/auditoria/integridade?de=2026-11-10&ate=2026-11-12');

    expect(await screen.findByRole('alert')).toHaveTextContent('Sem a cópia imutável não é possível afirmar a integridade');
    expect(screen.queryByText('Íntegro')).not.toBeInTheDocument();
  });

  it('recusa período invertido ou acima de 400 dias antes de chamar a API', async () => {
    const api = simularAPI({ ...BASE, [ROTA]: { corpo: pagina([]) } });
    const u = userEvent.setup();
    renderizarRota('/administracao/auditoria/integridade');

    await u.clear(await screen.findByLabelText('Primeiro dia'));
    await u.type(screen.getByLabelText('Primeiro dia'), '01012025');
    await u.clear(screen.getByLabelText('Último dia'));
    await u.type(screen.getByLabelText('Último dia'), '10112026');
    await u.click(screen.getByRole('button', { name: 'Verificar integridade' }));

    await waitFor(() => expect(screen.getByLabelText('Último dia')).toHaveAccessibleDescription('O período vai até 400 dias.'));
    expect(api.chamadasA(ROTA)).toHaveLength(0);
  });

  it('sem adm.ler_auditoria, nem abre', async () => {
    const api = simularAPI(rotaDaSessao(sessao({ permissoes: ['adm.ler'] })));
    renderizarRota('/administracao/auditoria/integridade');
    expect(await screen.findByRole('heading', { name: 'Sem acesso a esta área' })).toBeInTheDocument();
    expect(api.chamadasA(ROTA)).toHaveLength(0);
  });
});

describe('capturas da auditoria', () => {
  const CAPTURAS = {
    'GET /v1/capturas-auditoria': {
      corpo: pagina([captura(), captura({ tela: 'adm.usuarios', nome: 'Usuários', visualizacao: true, configurada: true }), captura({ tela: 'os.ordens', modulo: 'os', nome: 'Ordens de serviço' })]),
    },
  };

  it('lista por módulo, com as três chaves de cada tela e o que é sempre gravado, sem chave', async () => {
    simularAPI({ ...BASE, ...CAPTURAS });
    renderizarRota('/administracao/auditoria/capturas');

    expect(await screen.findByRole('heading', { name: 'Módulo adm' })).toBeInTheDocument();
    expect(screen.getByRole('heading', { name: 'Módulo os' })).toBeInTheDocument();
    expect(screen.getByText(/Entrada, saída, inclusão, alteração, exclusão, deleção/)).toBeInTheDocument();
    // Padrão: impressão e extração ligadas; visualização desligada.
    expect(screen.getByRole('switch', { name: 'Visualização em Pessoas' })).toHaveAttribute('aria-checked', 'false');
    expect(screen.getByRole('switch', { name: 'Impressão em Pessoas' })).toHaveAttribute('aria-checked', 'true');
    expect(screen.getByRole('switch', { name: 'Visualização em Usuários' })).toHaveAttribute('aria-checked', 'true');
    // Só as três chaves configuráveis por tela.
    expect(screen.getAllByRole('switch')).toHaveLength(9);
  });

  it('liga a visualização de uma tela com PUT das três chaves e confirma', async () => {
    const api = simularAPI({
      ...BASE,
      ...CAPTURAS,
      'PUT /v1/capturas-auditoria/adm.pessoas': { corpo: captura({ visualizacao: true, configurada: true }) },
    });
    const u = userEvent.setup();
    renderizarRota('/administracao/auditoria/capturas');

    await u.click(await screen.findByRole('switch', { name: 'Visualização em Pessoas' }));

    await waitFor(() => expect(api.chamadasA('PUT /v1/capturas-auditoria/adm.pessoas')).toHaveLength(1));
    expect(api.chamadasA('PUT /v1/capturas-auditoria/adm.pessoas')[0]?.corpo).toEqual({ visualizacao: true, impressao: true, extracao: true });
    expect(await screen.findByText('Captura de Visualização ligada em Pessoas.')).toBeInTheDocument();
  });

  it('recusa da API (403) aparece como erro', async () => {
    simularAPI({ ...BASE, ...CAPTURAS, 'PUT /v1/capturas-auditoria/adm.pessoas': problema(403, 'permissao_negada') });
    const u = userEvent.setup();
    renderizarRota('/administracao/auditoria/capturas');

    await u.click(await screen.findByRole('switch', { name: 'Extração em Pessoas' }));

    expect(await screen.findByRole('alert')).toHaveTextContent('Seu perfil não pode fazer isso');
  });

  it('sem adm.gerenciar_parametros, nem abre nem chama a API', async () => {
    const api = simularAPI(rotaDaSessao(sessao({ permissoes: ['adm.ler_auditoria'] })));
    renderizarRota('/administracao/auditoria/capturas');
    expect(await screen.findByRole('heading', { name: 'Sem acesso a esta área' })).toBeInTheDocument();
    expect(api.chamadasA('GET /v1/capturas-auditoria')).toHaveLength(0);
  });
});

describe('cabeçalho X-Tenaz-Tela', () => {
  it('toda chamada leva o código da tela do catálogo', async () => {
    const api = simularAPI({ ...BASE, ...EVENTOS });
    renderizarRota('/administracao/auditoria');
    await screen.findByRole('row', { name: /Alteração/ });
    expect(api.chamadasA('GET /v1/eventos-auditoria')[0]?.cabecalhos.get('X-Tenaz-Tela')).toBe('adm.auditoria');
  });

  it('a tela de configuração das capturas é adm.parametros', async () => {
    const api = simularAPI({ ...BASE, 'GET /v1/capturas-auditoria': { corpo: pagina([captura()]) } });
    renderizarRota('/administracao/auditoria/capturas');
    await screen.findByRole('heading', { name: 'Módulo adm' });
    expect(api.chamadasA('GET /v1/capturas-auditoria')[0]?.cabecalhos.get('X-Tenaz-Tela')).toBe('adm.parametros');
  });

  it('rota fora do catálogo (Minha conta) não envia o cabeçalho', async () => {
    const api = simularAPI({ ...BASE, 'GET /v1/usuarios/eu': { corpo: { id: 'u-1', nome: 'Ana', email: 'a@b.c', telefone: null, versao: 1 } } });
    renderizarRota('/conta');
    await screen.findByLabelText('Nome');
    expect(api.chamadasA('GET /v1/usuarios/eu')[0]?.cabecalhos.has('X-Tenaz-Tela')).toBe(false);
  });
});
