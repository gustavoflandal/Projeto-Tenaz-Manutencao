// Testes do acesso: login, renovação serializada, separação por conta e erros.
// O servidor falso implementa o contrato da seção 4.3, inclusive a detecção de
// reuso do refresh token rotativo, que é o que a serialização evita.

import { DatabaseSync } from 'node:sqlite';
import { beforeEach, describe, expect, it } from 'vitest';

import { arquivoDaConta } from '../sync/deposito-expo';
import { DepositoSQLite } from '../sync/deposito-sqlite';
import { FilaDeComandos } from '../sync/fila';
import { CofreEmMemoria } from './cofre';
import { ErroAcessoNegado, ErroApi, ErroContaMudou, ErroDeRede, ErroRetornoInvalido, ErroSessaoEncerrada } from './erros';
import { GerenteDeSessoes } from './sessoes';
import { criptoDeTeste, ServidorFalso, type UsuarioFalso } from './testes/servidor-falso';

const EMP_A = '0193f2a1-0000-7000-8000-00000000a001';
const EMP_B = '0193f2a1-0000-7000-8000-00000000b002';
const EMP_C = '0193f2a1-0000-7000-8000-00000000c003';

const ANA: UsuarioFalso = {
  id: '0193f2a1-0000-7000-8000-0000000000a1',
  nome: 'Ana Técnica',
  empresas: [
    { id: EMP_A, nome: 'Transportes Alfa' },
    { id: EMP_B, nome: 'Indústria Beta' },
  ],
};
const LUCAS: UsuarioFalso = { id: '0193f2a1-0000-7000-8000-0000000000a2', nome: 'Lucas Locatário', empresas: [{ id: EMP_C, nome: 'Locadora Gama' }], locatario: true };
const MARTA: UsuarioFalso = {
  id: '0193f2a1-0000-7000-8000-0000000000a3',
  nome: 'Marta Gestora',
  empresas: [{ id: EMP_A, nome: 'Transportes Alfa' }],
  mfaPendenteEm: [EMP_A],
};

let servidor: ServidorFalso;
let cofre: CofreEmMemoria;
let gerente: GerenteDeSessoes;

function criarGerente(): GerenteDeSessoes {
  return new GerenteDeSessoes({
    cofre,
    cripto: criptoDeTeste,
    apiUrl: 'https://api.tenaz.test',
    http: servidor.fetch as typeof fetch,
    agora: () => servidor.relogio,
  });
}

/** Faz o login completo como a tela faz: prepara, o servidor autoriza, o app conclui. */
async function entrar(usuario: UsuarioFalso, empresaId: string) {
  const tentativa = await gerente.prepararLogin();
  const retorno = servidor.autorizar(usuario.id, empresaId, tentativa.url);
  return gerente.concluirLogin(tentativa, retorno);
}

beforeEach(() => {
  servidor = new ServidorFalso(ANA, LUCAS, MARTA);
  cofre = new CofreEmMemoria();
  gerente = criarGerente();
});

describe('login', () => {
  it('troca o código por tokens, lê a sessão e guarda tudo no cofre', async () => {
    const conta = await entrar(ANA, EMP_A);

    expect(conta).toMatchObject({ empresaId: EMP_A, usuarioId: ANA.id, empresaNome: 'Transportes Alfa', mfaPendente: false });
    expect(conta.empresas.map((e) => e.id)).toEqual([EMP_A, EMP_B]);
    expect(gerente.estado()).toEqual({ tipo: 'entrou', conta });
    const guardada = JSON.parse((await cofre.ler('tenaz.credencial')) ?? '{}');
    expect(guardada.refreshToken).toMatch(/^tnz_rt_/);
    // Nada de token fora do cofre: o estado publicado não carrega credencial.
    expect(JSON.stringify(gerente.estado())).not.toContain('tnz_rt_');
  });

  it('não repete o login com o mesmo código (uso único)', async () => {
    const tentativa = await gerente.prepararLogin();
    const retorno = servidor.autorizar(ANA.id, EMP_A, tentativa.url);
    await gerente.concluirLogin(tentativa, retorno);
    await gerente.sair();
    await expect(gerente.concluirLogin(tentativa, retorno)).rejects.toMatchObject({ error: 'invalid_grant' });
  });

  it('rejeita retorno com state de outra tentativa sem chamar o servidor', async () => {
    const tentativa = await gerente.prepararLogin();
    const retorno = servidor.autorizar(ANA.id, EMP_A, tentativa.url).replace(/state=[^&]+/, 'state=forjado');
    const antes = servidor.pedidos.length;
    await expect(gerente.concluirLogin(tentativa, retorno)).rejects.toBeInstanceOf(ErroRetornoInvalido);
    expect(servidor.pedidos.length).toBe(antes);
    expect(gerente.estado().tipo).not.toBe('entrou');
  });

  it('rejeita código trocado com verificador PKCE diferente', async () => {
    const tentativa = await gerente.prepararLogin();
    const retorno = servidor.autorizar(ANA.id, EMP_A, tentativa.url);
    await expect(gerente.concluirLogin({ ...tentativa, verificador: 'x'.repeat(43) }, retorno)).rejects.toMatchObject({ error: 'invalid_grant' });
    expect(await cofre.ler('tenaz.credencial')).toBeNull();
  });

  it('trata access_denied (locatário) sem tentar trocar código', async () => {
    const tentativa = await gerente.prepararLogin();
    const retorno = servidor.autorizar(LUCAS.id, EMP_C, tentativa.url);
    const antes = servidor.pedidos.length;
    await expect(gerente.concluirLogin(tentativa, retorno)).rejects.toBeInstanceOf(ErroAcessoNegado);
    expect(servidor.pedidos.length).toBe(antes);
    expect(await cofre.ler('tenaz.credencial')).toBeNull();
  });

  it('com mfa_pendente, entra restrito e não tenta renovar token por causa do 401', async () => {
    const conta = await entrar(MARTA, EMP_A);
    expect(conta.mfaPendente).toBe(true);

    const antes = servidor.refreshChamadas;
    const erro = await gerente.autenticado(conta, '/v1/empresa').catch((e: unknown) => e);
    expect(erro).toBeInstanceOf(ErroApi);
    expect(erro).toMatchObject({ status: 401, codigo: 'nao_autenticado', mfaPendente: true });
    expect(servidor.refreshChamadas).toBe(antes); // 401 de MFA não é token vencido
    expect(gerente.estado().tipo).toBe('entrou'); // e a pessoa continua logada para ativar o MFA / sair
  });

  it('abre offline: carregar lê só o cofre, sem rede', async () => {
    await entrar(ANA, EMP_A);
    servidor.semRede = true;
    const reaberto = criarGerente();
    const antes = servidor.pedidos.length;
    const estado = await reaberto.carregar();
    expect(estado.tipo).toBe('entrou');
    expect(servidor.pedidos.length).toBe(antes);
  });
});

describe('renovação serializada', () => {
  it('cinco renovações simultâneas viram UMA ida ao servidor', async () => {
    const conta = await entrar(ANA, EMP_A);
    let liberar!: () => void;
    servidor.segurarRefresh = new Promise<void>((r) => (liberar = r));
    const antes = servidor.refreshChamadas;
    const tokenVelho = (JSON.parse((await cofre.ler('tenaz.credencial')) ?? '{}') as { accessToken: string }).accessToken;

    const todas = [1, 2, 3, 4, 5].map(() => gerente.renovar(tokenVelho));
    await Promise.resolve();
    liberar();
    const resultados = await Promise.all(todas);

    expect(servidor.refreshChamadas - antes).toBe(1);
    expect(new Set(resultados.map((r) => r.refreshToken)).size).toBe(1);
    expect(gerente.estado()).toMatchObject({ tipo: 'entrou' }); // sem reuso, a sessão não caiu
    await expect(gerente.autenticado(conta, '/v1/empresa')).resolves.toMatchObject({ ok: true });
  });

  it('chamadas autenticadas concorrentes com access vencido renovam uma vez só', async () => {
    const conta = await entrar(ANA, EMP_A);
    servidor.avancar(16 * 60_000); // 15 min de validade: venceu
    const antes = servidor.refreshChamadas;

    const respostas = await Promise.all([1, 2, 3, 4].map(() => gerente.autenticado<{ ok: boolean }>(conta, '/v1/empresa')));

    expect(respostas.every((r) => r.ok)).toBe(true);
    expect(servidor.refreshChamadas - antes).toBe(1);
  });

  it('um 401 chegando depois da renovação usa a credencial nova, sem outra ida', async () => {
    await entrar(ANA, EMP_A);
    const velho = (JSON.parse((await cofre.ler('tenaz.credencial')) ?? '{}') as { accessToken: string }).accessToken;
    await gerente.renovar(velho);
    const antes = servidor.refreshChamadas;
    await gerente.renovar(velho); // ainda segurando o token antigo
    expect(servidor.refreshChamadas).toBe(antes);
  });

  it('o refresh token novo já está no cofre quando a renovação responde', async () => {
    await entrar(ANA, EMP_A);
    const antes = JSON.parse((await cofre.ler('tenaz.credencial')) ?? '{}') as { refreshToken: string };
    const nova = await gerente.renovar();
    const depois = JSON.parse((await cofre.ler('tenaz.credencial')) ?? '{}') as { refreshToken: string };
    expect(depois.refreshToken).toBe(nova.refreshToken);
    expect(depois.refreshToken).not.toBe(antes.refreshToken);
  });

  it('sem rede, a renovação falha mas a sessão e o refresh token ficam', async () => {
    const conta = await entrar(ANA, EMP_A);
    const antes = await cofre.ler('tenaz.credencial');
    servidor.avancar(16 * 60_000);
    servidor.semRede = true;

    await expect(gerente.autenticado(conta, '/v1/empresa')).rejects.toBeInstanceOf(ErroDeRede);

    expect(gerente.estado().tipo).toBe('entrou');
    expect(await cofre.ler('tenaz.credencial')).toBe(antes);
    // Voltou o sinal: a mesma credencial renova normalmente.
    servidor.semRede = false;
    await expect(gerente.autenticado(conta, '/v1/empresa')).resolves.toMatchObject({ ok: true });
  });

  it('sete dias depois, sem rede no meio, o refresh ainda renova (a sessão sobrevive ao offline)', async () => {
    const conta = await entrar(ANA, EMP_A);
    servidor.avancar(7 * 24 * 3600_000);
    servidor.semRede = true;
    await expect(gerente.autenticado(conta, '/v1/empresa')).rejects.toBeInstanceOf(ErroDeRede);
    servidor.semRede = false;
    await expect(gerente.autenticado(conta, '/v1/empresa')).resolves.toMatchObject({ ok: true });
  });

  it('invalid_grant encerra a sessão e apaga só as credenciais, NUNCA a fila local', async () => {
    const conta = await entrar(ANA, EMP_A);
    const banco = new DatabaseSync(':memory:'); // o "arquivo" da conta
    const fila = new FilaDeComandos({ deposito: new DepositoSQLite(banco), transporte: { enviar: async () => [] } });
    fila.enfileirar({
      id: '0193f2a1-0000-7000-8000-000000000001',
      tipo: 'ApontarHoras',
      agregado: 'ordem_servico:0193f2a1-0000-7000-8000-00000000000a',
      carga: { minutos: 30 },
      versaoBase: 1,
      horarioAparelho: '2026-11-12T10:00:00Z',
    });

    // O servidor revoga a família (por exemplo, o administrador removeu o aparelho).
    const refresh = (JSON.parse((await cofre.ler('tenaz.credencial')) ?? '{}') as { refreshToken: string }).refreshToken;
    await servidor.fetch('https://x/oauth/revoke', { method: 'POST', body: `token=${refresh}` });
    servidor.avancar(16 * 60_000);

    await expect(gerente.autenticado(conta, '/v1/empresa')).rejects.toBeInstanceOf(ErroSessaoEncerrada);

    expect(gerente.estado()).toEqual({ tipo: 'saiu', aviso: 'sessao_encerrada' });
    expect(await cofre.ler('tenaz.credencial')).toBeNull();
    expect(fila.quantidadeNaoSincronizada()).toBe(1); // a fila continua inteira
  });

  it('reuso do refresh token derruba a família: a prova de que serializar é obrigatório', async () => {
    await entrar(ANA, EMP_A);
    const rt = (JSON.parse((await cofre.ler('tenaz.credencial')) ?? '{}') as { refreshToken: string }).refreshToken;
    const usar = () => servidor.fetch('https://x/oauth/token', { method: 'POST', body: `grant_type=refresh_token&client_id=tenaz-app&refresh_token=${rt}` });
    expect((await usar()).status).toBe(200);
    expect((await usar()).status).toBe(401); // segundo uso do mesmo token: reuso
  });
});

describe('troca de empresa e separação por conta', () => {
  it('troca sem sair: tokens novos, outra empresa, outro arquivo de banco', async () => {
    const a = await entrar(ANA, EMP_A);
    const b = await gerente.trocarEmpresa(EMP_B);

    expect(b).toMatchObject({ empresaId: EMP_B, empresaNome: 'Indústria Beta', usuarioId: ANA.id });
    expect(b.arquivoBanco).toBe(arquivoDaConta(EMP_B, ANA.id));
    expect(b.arquivoBanco).not.toBe(a.arquivoBanco);
    const guardada = JSON.parse((await cofre.ler('tenaz.credencial')) ?? '{}');
    expect(guardada.empresaId).toBe(EMP_B);
    expect(gerente.estado()).toEqual({ tipo: 'entrou', conta: b });
    await expect(gerente.autenticado(b, '/v1/empresa')).resolves.toMatchObject({ empresa: EMP_B });
  });

  it('cada empresa tem a sua fila: nada escrito em A aparece em B, e a volta encontra tudo', async () => {
    const bancos = new Map<string, DatabaseSync>();
    const abrir = (arquivo: string) => {
      let b = bancos.get(arquivo);
      if (!b) bancos.set(arquivo, (b = new DatabaseSync(':memory:')));
      return new FilaDeComandos({ deposito: new DepositoSQLite(b), transporte: { enviar: async () => [] } });
    };
    const ids = (arquivo: string) => new DepositoSQLite(bancos.get(arquivo) as DatabaseSync).todos().map((c) => c.id);
    const cmd = (n: number) => ({
      id: `0193f2a1-0000-7000-8000-${String(n).padStart(12, '0')}`,
      tipo: 'ApontarHoras' as const,
      agregado: 'ordem_servico:0193f2a1-0000-7000-8000-00000000000a',
      carga: { minutos: n },
      versaoBase: 1,
      horarioAparelho: '2026-11-12T10:00:00Z',
    });

    const a = await entrar(ANA, EMP_A);
    abrir(a.arquivoBanco).enfileirar(cmd(1));
    const b = await gerente.trocarEmpresa(EMP_B);
    expect(abrir(b.arquivoBanco).quantidadeNaoSincronizada()).toBe(0);
    abrir(b.arquivoBanco).enfileirar(cmd(2));
    const deVolta = await gerente.trocarEmpresa(EMP_A);

    expect(ids(deVolta.arquivoBanco)).toEqual([cmd(1).id]);
    expect(ids(b.arquivoBanco)).toEqual([cmd(2).id]);
  });

  it('usuários diferentes na mesma empresa e no mesmo aparelho também têm bancos diferentes', async () => {
    const ana = await entrar(ANA, EMP_A);
    await gerente.sair();
    const marta = await entrar(MARTA, EMP_A);
    expect(marta.arquivoBanco).not.toBe(ana.arquivoBanco);
  });

  it('um pedido da conta A nunca sai com a credencial da conta B', async () => {
    const a = await entrar(ANA, EMP_A);
    await gerente.trocarEmpresa(EMP_B);
    const antes = servidor.pedidos.length;

    await expect(gerente.autenticado(a, '/v1/ordens-servico')).rejects.toBeInstanceOf(ErroContaMudou);

    expect(servidor.pedidos.length).toBe(antes); // nem chegou ao servidor
  });

  it('troca no meio de um envio: o pedido de A que tomou 401 não é reenviado como B', async () => {
    const a = await entrar(ANA, EMP_A);
    servidor.avancar(16 * 60_000); // o access de A venceu
    let soltar!: () => void;
    servidor.segurarTroca = new Promise<void>((r) => (soltar = r));

    // A troca começa primeiro e segura a trava; o envio de A chega depois.
    const troca = gerente.trocarEmpresa(EMP_B);
    const envio = gerente.autenticado(a, '/v1/ordens-servico').catch((e: unknown) => e);
    await new Promise((r) => setTimeout(r, 0));
    soltar();
    await troca;
    const resultado = await envio;

    expect(resultado).toBeInstanceOf(ErroContaMudou);
    const pedidosDeOS = servidor.pedidos.filter((p) => p.caminho === '/v1/ordens-servico');
    expect(pedidosDeOS.every((p) => p.empresaDoToken !== EMP_B)).toBe(true);
  });

  it('empresa fora do acesso (404) mantém a conta atual intacta', async () => {
    const a = await entrar(ANA, EMP_A);
    const erro = await gerente.trocarEmpresa(EMP_C).catch((e: unknown) => e);
    expect(erro).toMatchObject({ status: 404 });
    expect(gerente.estado()).toEqual({ tipo: 'entrou', conta: a });
    await expect(gerente.autenticado(a, '/v1/empresa')).resolves.toMatchObject({ empresa: EMP_A });
  });

  it('sem rede a troca falha e a empresa atual continua', async () => {
    const a = await entrar(ANA, EMP_A);
    servidor.semRede = true;
    await expect(gerente.trocarEmpresa(EMP_B)).rejects.toBeInstanceOf(ErroDeRede);
    expect(gerente.estado()).toEqual({ tipo: 'entrou', conta: a });
  });

  it('a troca vale também com mfa_pendente (a rota aceita a sessão restrita)', async () => {
    const m = await entrar(MARTA, EMP_A);
    const antes = servidor.pedidos.length;
    await expect(gerente.trocarEmpresa(EMP_A)).resolves.toMatchObject({ empresaId: EMP_A });
    expect(servidor.pedidos.length).toBeGreaterThan(antes);
    expect(m.mfaPendente).toBe(true);
  });
});

describe('sair', () => {
  it('apaga as credenciais, revoga no servidor e preserva a fila local', async () => {
    const conta = await entrar(ANA, EMP_A);
    const banco = new DatabaseSync(':memory:');
    const fila = new FilaDeComandos({ deposito: new DepositoSQLite(banco), transporte: { enviar: async () => [] } });
    fila.enfileirar({
      id: '0193f2a1-0000-7000-8000-000000000009',
      tipo: 'ApontarHoras',
      agregado: 'ordem_servico:0193f2a1-0000-7000-8000-00000000000a',
      carga: {},
      versaoBase: 1,
      horarioAparelho: '2026-11-12T10:00:00Z',
    });
    const refresh = (JSON.parse((await cofre.ler('tenaz.credencial')) ?? '{}') as { refreshToken: string }).refreshToken;

    await gerente.sair();

    expect(gerente.estado()).toEqual({ tipo: 'saiu', aviso: null });
    expect(await cofre.ler('tenaz.credencial')).toBeNull();
    expect(await cofre.ler('tenaz.sessao')).toBeNull();
    // Revogado: o refresh antigo não renova mais.
    const r = await servidor.fetch('https://x/oauth/token', { method: 'POST', body: `grant_type=refresh_token&client_id=tenaz-app&refresh_token=${refresh}` });
    expect(r.status).toBe(401);
    expect(fila.quantidadeNaoSincronizada()).toBe(1);
    // Para onde a fila voltaria: o mesmo arquivo da mesma conta.
    expect(conta.arquivoBanco).toBe(arquivoDaConta(EMP_A, ANA.id));
    expect(await gerente.ultimaConta()).toEqual({ empresaNome: 'Transportes Alfa', usuarioNome: 'Ana Técnica' });
  });

  it('sem rede sai do mesmo jeito e repete a revogação quando houver sinal', async () => {
    await entrar(ANA, EMP_A);
    const refresh = (JSON.parse((await cofre.ler('tenaz.credencial')) ?? '{}') as { refreshToken: string }).refreshToken;
    servidor.semRede = true;

    await gerente.sair();

    expect(await cofre.ler('tenaz.credencial')).toBeNull();
    expect(JSON.parse((await cofre.ler('tenaz.revogacoes-pendentes')) ?? '[]')).toEqual([refresh]);

    servidor.semRede = false;
    await gerente.tentarRevogacoesPendentes();
    expect(await cofre.ler('tenaz.revogacoes-pendentes')).toBeNull();
    const r = await servidor.fetch('https://x/oauth/token', { method: 'POST', body: `grant_type=refresh_token&client_id=tenaz-app&refresh_token=${refresh}` });
    expect(r.status).toBe(401);
  });
});

describe('atualizar sessão', () => {
  it('relê a sessão e atualiza o estado (por exemplo, MFA ativado na web)', async () => {
    const m = await entrar(MARTA, EMP_A);
    expect(m.mfaPendente).toBe(true);
    MARTA.mfaPendenteEm = [];
    try {
      const nova = await gerente.atualizarSessao();
      expect(nova.mfaPendente).toBe(false);
      expect(gerente.estado()).toEqual({ tipo: 'entrou', conta: nova });
    } finally {
      MARTA.mfaPendenteEm = [EMP_A];
    }
  });
});
