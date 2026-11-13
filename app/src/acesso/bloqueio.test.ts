import { beforeEach, describe, expect, it } from 'vitest';

import { type AutenticadorLocal, BloqueioLocal, TEMPO_PARA_BLOQUEAR_MS } from './bloqueio';
import { CofreEmMemoria } from './cofre';

let agora = 0;
let disponivel = true;
let aprovar = true;
let cofre: CofreEmMemoria;
const autenticador: AutenticadorLocal = {
  disponivel: async () => disponivel,
  autenticar: async () => aprovar,
};

beforeEach(() => {
  agora = 1_000_000;
  disponivel = true;
  aprovar = true;
  cofre = new CofreEmMemoria();
});

const novo = () => new BloqueioLocal(cofre, autenticador, () => agora);

describe('bloqueio local', () => {
  it('desligado por padrão: nunca bloqueia', async () => {
    const b = novo();
    await b.iniciar();
    b.aoIrParaSegundoPlano();
    agora += TEMPO_PARA_BLOQUEAR_MS * 10;
    b.aoVoltar();
    expect(b.estaBloqueado()).toBe(false);
  });

  it('não liga sem biometria nem PIN cadastrado no aparelho', async () => {
    disponivel = false;
    const b = novo();
    await b.iniciar();
    expect(await b.ativar('m')).toBe(false);
    expect(b.ativo()).toBe(false);
  });

  it('não liga se a pessoa não passa na verificação', async () => {
    aprovar = false;
    const b = novo();
    await b.iniciar();
    expect(await b.ativar('m')).toBe(false);
  });

  it('ligado, a abertura a frio já nasce bloqueada', async () => {
    const b = novo();
    await b.iniciar();
    await b.ativar('m');
    const reaberto = novo();
    await reaberto.iniciar();
    expect(reaberto.estaBloqueado()).toBe(true);
  });

  it('voltar rápido do segundo plano não bloqueia; demorar bloqueia', async () => {
    const b = novo();
    await b.iniciar();
    await b.ativar('m');

    b.aoIrParaSegundoPlano();
    agora += TEMPO_PARA_BLOQUEAR_MS - 1;
    b.aoVoltar();
    expect(b.estaBloqueado()).toBe(false);

    b.aoIrParaSegundoPlano();
    agora += TEMPO_PARA_BLOQUEAR_MS;
    b.aoVoltar();
    expect(b.estaBloqueado()).toBe(true);
  });

  it('desbloqueia só quando a verificação passa', async () => {
    const b = novo();
    await b.iniciar();
    await b.ativar('m');
    b.aoIrParaSegundoPlano();
    agora += TEMPO_PARA_BLOQUEAR_MS;
    b.aoVoltar();

    aprovar = false;
    expect(await b.desbloquear('m')).toBe(false);
    expect(b.estaBloqueado()).toBe(true);
    aprovar = true;
    expect(await b.desbloquear('m')).toBe(true);
    expect(b.estaBloqueado()).toBe(false);
  });
});
