// Acessibilidade (axe-core, WCAG 2.1 AA, RNF-025) de todas as telas de E01,
// nos temas claro e escuro. O contraste de cor não é medido no jsdom; os pares
// de tokens são conferidos em pacotes/tokens-ui. Aqui valem nomes, rótulos,
// papéis, ordem de leitura, landmarks e relações ARIA.

import { screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import axe from 'axe-core';
import { describe, expect, it } from 'vitest';

import { i18n } from '../lib/i18n';
import { definirPreferenciaTema } from '../lib/tema';
import { problema, simularAPI, type Manipulador } from '../testes/api-simulada';
import {
  PERFIS,
  captura,
  diaVerificado,
  PERMISSOES,
  empresa,
  evento,
  pagina,
  pessoa,
  resumoEvento,
  resumoPessoa,
  resumoUnidade,
  resumoUsuario,
  sessao,
  unidade,
  usuario,
} from '../testes/dados';
import { renderizarRota, rotaDaSessao } from '../testes/renderizar';

const REGRAS: axe.RunOptions = {
  runOnly: { type: 'tag', values: ['wcag2a', 'wcag2aa', 'wcag21a', 'wcag21aa', 'best-practice'] },
  resultTypes: ['violations'],
};

async function violacoes() {
  const resultado = await axe.run(document.documentElement, REGRAS);
  return resultado.violations.map((v) => ({
    regra: v.id,
    descricao: v.help,
    elementos: v.nodes.map((n) => n.target.join(' ')),
  }));
}

const FERIADO = { id: 'f-1', unidade_id: 'un-1', data: '2026-12-25', descricao: 'Natal', versao: 6 };

const ROTAS: Record<string, Manipulador> = {
  ...rotaDaSessao(sessao()),
  'GET /v1/empresa': { corpo: empresa() },
  'GET /v1/unidades': { corpo: pagina([resumoUnidade()]) },
  'GET /v1/unidades/un-1': { corpo: unidade() },
  'GET /v1/unidades/un-1/feriados': { corpo: pagina([FERIADO]) },
  'GET /v1/pessoas': { corpo: pagina([resumoPessoa()]) },
  'GET /v1/pessoas/pe-1': { corpo: pessoa() },
  'GET /v1/usuarios': { corpo: pagina([resumoUsuario()]) },
  'GET /v1/usuarios/u-2': { corpo: usuario({ permissoes_avulsas: ['os.encerrar'] }) },
  'GET /v1/perfis': { corpo: pagina(PERFIS) },
  'GET /v1/permissoes': { corpo: pagina(PERMISSOES) },
  // O idioma do usuário acompanha o da interface no momento da chamada, como o de quem escolheu o idioma.
  'GET /v1/usuarios/eu': () => ({ corpo: usuario({ id: 'u-1', nome: 'Ana Souza', idioma: i18n.language === 'es-419' ? 'es-419' : 'pt-BR' }) }),
  'GET /v1/eventos-auditoria': { corpo: pagina([resumoEvento()]) },
  'GET /v1/eventos-auditoria/ev-1': { corpo: evento() },
  'GET /v1/capturas-auditoria': { corpo: pagina([captura(), captura({ tela: 'adm.usuarios', nome: 'Usuários', modulo: 'adm' })]) },
  'GET /v1/verificacoes-auditoria': {
    corpo: pagina([diaVerificado(), diaVerificado({ dia: '2026-11-11', estado: 'violado', motivo: 'evento_alterado', evento_id: 'ev-1' }), diaVerificado({ dia: '2026-11-13', estado: 'nao_selado', motivo: 'ainda_nao_selado' })]),
  },
  'POST /v1/auth/mfa/iniciar': {
    corpo: { segredo_base32: 'JBSWY3DPEHPK3PXP', uri_otpauth: 'otpauth://totp/Tenaz:ana?secret=JBSWY3DPEHPK3PXP' },
  },
  'POST /v1/auth/convites/consultar': { corpo: { email: 'novo@polar.com.br', pede_senha: true } },
};

const TOKEN = 't'.repeat(43);
const PEDIDO_APP =
  '/autorizar-app?response_type=code&client_id=tenaz-app&redirect_uri=br.com.tenaz.app%3A%2Foauth&code_challenge=abc&code_challenge_method=S256';

interface Tela {
  readonly nome: string;
  readonly caminho: string;
  /** Texto de um título ou campo que indica que a tela terminou de carregar. */
  readonly pronta: string | RegExp;
  readonly papel?: 'heading' | 'button' | 'table';
  readonly rotas?: Record<string, Manipulador>;
}

const TELAS: readonly Tela[] = [
  { nome: 'entrar', caminho: '/entrar', pronta: 'Entrar' },
  { nome: 'criar conta', caminho: '/criar-conta', pronta: 'Criar conta da empresa' },
  { nome: 'esqueci a senha', caminho: '/esqueci-a-senha', pronta: 'Recuperar senha' },
  { nome: 'redefinir senha', caminho: `/redefinir-senha#token=${TOKEN}`, pronta: 'Criar senha nova' },
  { nome: 'convite', caminho: `/convite#token=${TOKEN}`, pronta: 'Aceitar convite', papel: 'button' },
  { nome: 'ativar segundo fator (sessão restrita)', caminho: '/ativar-mfa', pronta: 'Ative o segundo fator', rotas: rotaDaSessao(sessao({ mfa_pendente: true })) },
  { nome: 'autorizar o app', caminho: PEDIDO_APP, pronta: 'Autorizar o app' },
  { nome: 'minha conta', caminho: '/conta', pronta: 'Minha conta' },
  { nome: 'empresa', caminho: '/administracao/empresa', pronta: 'Empresa' },
  { nome: 'unidades', caminho: '/administracao/unidades', pronta: 'Unidades' },
  { nome: 'nova unidade', caminho: '/administracao/unidades/nova', pronta: 'Nova unidade' },
  { nome: 'detalhe da unidade', caminho: '/administracao/unidades/un-1', pronta: 'Matriz Contagem' },
  { nome: 'pessoas', caminho: '/administracao/pessoas', pronta: 'Pessoas' },
  { nome: 'nova pessoa', caminho: '/administracao/pessoas/nova', pronta: 'Nova pessoa' },
  { nome: 'detalhe da pessoa', caminho: '/administracao/pessoas/pe-1', pronta: 'João Lima' },
  { nome: 'usuários', caminho: '/administracao/usuarios', pronta: 'Usuários' },
  { nome: 'convidar usuário', caminho: '/administracao/usuarios/novo', pronta: 'Convidar usuário' },
  { nome: 'detalhe do usuário', caminho: '/administracao/usuarios/u-2', pronta: 'Carlos Lima' },
  { nome: 'perfis', caminho: '/administracao/perfis', pronta: 'Perfis' },
  { nome: 'auditoria', caminho: '/administracao/auditoria', pronta: 'Auditoria' },
  { nome: 'integridade da trilha', caminho: '/administracao/auditoria/integridade?de=2026-11-10&ate=2026-11-13', pronta: 'Integridade da trilha', papel: 'table' },
  { nome: 'capturas da auditoria', caminho: '/administracao/auditoria/capturas', pronta: 'Módulo adm', papel: 'heading' },
  { nome: 'não encontrada', caminho: '/nao-existe', pronta: 'Página não encontrada' },
  { nome: 'sem permissão', caminho: '/administracao/usuarios', pronta: 'Sem acesso a esta área', rotas: rotaDaSessao(sessao({ permissoes: ['adm.ler'] })) },
  { nome: 'erro de carga', caminho: '/administracao/pessoas', pronta: 'Tentar de novo', papel: 'button', rotas: { 'GET /v1/pessoas': problema(500, 'erro_interno') } },
];

async function esperarPronta(tela: Tela, idioma: string) {
  if (idioma === 'pt-BR') {
    const papel = tela.papel ?? 'heading';
    await screen.findByRole(papel, { name: tela.pronta });
  } else {
    // Nos outros idiomas os nomes mudam: espera a página e o fim das cargas.
    await screen.findByRole('main');
  }
  await waitFor(() => expect(document.querySelector('[aria-busy="true"]')).toBeNull());
  // Deixa terminar as consultas que as telas disparam além do título.
  await new Promise((resolver) => setTimeout(resolver, 40));
}

describe.each([
  ['pt-BR', 'claro'],
  ['pt-BR', 'escuro'],
  ['es-419', 'claro'],
  ['es-419', 'escuro'],
] as const)('acessibilidade em %s no tema %s', (idioma, tema) => {
  it.each(TELAS)('$nome não tem violação', async (tela) => {
    await i18n.changeLanguage(idioma);
    definirPreferenciaTema(tema);
    simularAPI({ ...ROTAS, ...tela.rotas });
    renderizarRota(tela.caminho);
    await esperarPronta(tela, idioma);

    expect(document.documentElement).toHaveAttribute('data-tema', tema);
    expect(document.documentElement).toHaveAttribute('lang', idioma);
    expect(await violacoes()).toEqual([]);
  });
});

describe('acessibilidade de diálogos e estados', () => {
  it('modal de turno, confirmação, detalhe do evento e erro de formulário', async () => {
    simularAPI(ROTAS);
    const u = userEvent.setup();
    renderizarRota('/administracao/unidades/un-1');

    await u.click(await screen.findByRole('button', { name: 'Adicionar turno' }));
    const modal = await screen.findByRole('dialog', { name: 'Novo turno' });
    await u.click(within(modal).getByRole('button', { name: 'Salvar' }));
    await within(modal).findAllByText('Use HH:MM, de 00:00 a 23:59.');
    expect(await violacoes()).toEqual([]);
    await u.keyboard('{Escape}');

    await u.click(await screen.findByRole('button', { name: 'Inativar unidade' }));
    await screen.findByRole('alertdialog');
    expect(await violacoes()).toEqual([]);
  });

  it('detalhe do evento aberto', async () => {
    simularAPI(ROTAS);
    const u = userEvent.setup();
    renderizarRota('/administracao/auditoria');
    await u.click(await screen.findByRole('button', { name: /^Ver detalhes/ }));
    await screen.findByRole('table', { name: 'Valores antes e depois' });
    expect(await violacoes()).toEqual([]);
  });

  it('menu do usuário aberto', async () => {
    simularAPI(ROTAS);
    const u = userEvent.setup();
    renderizarRota('/');
    await u.click(await screen.findByRole('button', { name: 'Menu de Ana Souza' }));
    await screen.findByRole('menu');
    expect(await violacoes()).toEqual([]);
  });

  it('formulário de pessoa com erros por campo e resumo', async () => {
    simularAPI(ROTAS);
    const u = userEvent.setup();
    renderizarRota('/administracao/pessoas/nova');
    await u.type(await screen.findByLabelText('CPF (opcional)'), '111.111.111-11');
    await u.click(screen.getByRole('button', { name: 'Cadastrar pessoa' }));
    await screen.findByRole('alert');
    expect(await violacoes()).toEqual([]);
  });

  it('editor de turnos com modelo, erro e sobreposição', async () => {
    simularAPI(ROTAS);
    const u = userEvent.setup();
    renderizarRota('/administracao/unidades/nova');
    await u.click(await screen.findByRole('button', { name: 'Comercial' }));
    await u.click(screen.getByRole('button', { name: 'Adicionar turno de Segunda-feira' }));
    const segunda = screen.getByRole('group', { name: 'Segunda-feira' });
    await u.type(within(segunda).getAllByLabelText('Início')[1] as HTMLElement, '10:00');
    await u.type(within(segunda).getAllByLabelText('Fim')[1] as HTMLElement, '12:00');
    await u.click(screen.getByRole('button', { name: 'Criar unidade' }));
    await within(segunda).findByText(/se sobrepõe/);
    expect(await violacoes()).toEqual([]);
  });
});
