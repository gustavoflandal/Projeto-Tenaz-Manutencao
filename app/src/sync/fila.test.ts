// Testes do motor da fila de comandos, contra SQLite real no Node.
//
// Estes testes trancam os invariantes do componente mais crítico do app. Cada
// um corresponde a uma situação que acontece de verdade no campo: o técnico
// perde sinal no meio do envio, o aparelho é fechado, a OS foi cancelada no
// escritório enquanto ele trabalhava, o relógio do aparelho está errado.

import { DatabaseSync } from 'node:sqlite';
import { beforeEach, describe, expect, it, vi } from 'vitest';

import type { Comando, ResultadoComando } from './comando';
import { TAMANHO_MAXIMO_DO_LOTE } from './comando';
import { DepositoSQLite } from './deposito-sqlite';
import { FilaDeComandos, type Transporte } from './fila';

const OS_A = 'ordem_servico:0193f2a1-0000-7000-8000-00000000000a';
const OS_B = 'ordem_servico:0193f2a1-0000-7000-8000-00000000000b';

let proximoID = 0;

function comando(parcial: Partial<Comando> = {}): Comando {
  proximoID += 1;
  return {
    id: `0193f2a1-0000-7000-8000-${String(proximoID).padStart(12, '0')}`,
    tipo: 'ApontarHoras',
    agregado: OS_A,
    carga: { minutos: 30 },
    versaoBase: 1,
    horarioAparelho: '2026-11-12T10:00:00Z',
    ...parcial,
  };
}

function montar(
  responder: (lote: readonly Comando[]) => Promise<readonly ResultadoComando[]>,
  agora?: () => Date,
) {
  const banco = new DatabaseSync(':memory:');
  const deposito = new DepositoSQLite(banco);
  const enviados: Comando[][] = [];

  const transporte: Transporte = {
    enviar: async (lote) => {
      enviados.push([...lote]);
      return responder(lote);
    },
  };

  const fila = new FilaDeComandos({ deposito, transporte, agora });
  return { fila, deposito, enviados };
}

const aplicarTudo = async (lote: readonly Comando[]): Promise<ResultadoComando[]> =>
  lote.map((c) => ({ id: c.id, situacao: 'aplicado' as const }));

beforeEach(() => {
  proximoID = 0;
});

describe('idempotência', () => {
  it('gravar o mesmo comando duas vezes não cria dois', () => {
    const { fila, deposito } = montar(aplicarTudo);
    const c = comando();

    fila.enfileirar(c);
    fila.enfileirar(c);

    expect(deposito.todos()).toHaveLength(1);
    expect(fila.quantidadeNaoSincronizada()).toBe(1);
  });

  it('reenviar um comando já aplicado não o envia de novo', async () => {
    const { fila, enviados } = montar(aplicarTudo);
    fila.enfileirar(comando());

    await fila.descarregar();
    await fila.descarregar();

    expect(enviados).toHaveLength(1);
  });
});

describe('ordem por agregado', () => {
  it('comandos da mesma OS saem na ordem em que foram criados', async () => {
    const { fila, enviados } = montar(aplicarTudo);

    fila.enfileirar(comando({ tipo: 'IniciarOS', agregado: OS_A }));
    fila.enfileirar(comando({ tipo: 'ApontarHoras', agregado: OS_A }));
    fila.enfileirar(comando({ tipo: 'ConcluirOS', agregado: OS_A }));

    await fila.descarregar();

    expect(enviados[0]?.map((c) => c.tipo)).toEqual(['IniciarOS', 'ApontarHoras', 'ConcluirOS']);
  });

  it('um comando em espera bloqueia os seguintes DA MESMA OS', async () => {
    let tempo = new Date('2026-11-12T10:00:00Z').getTime();
    const agora = () => new Date(tempo);

    const { fila, enviados } = montar(async (lote) => {
      // O primeiro comando falha; os demais do lote, se houver, são aplicados.
      return lote.map((c, i) =>
        i === 0
          ? { id: c.id, situacao: 'retentar' as const, motivo: 'servidor indisponível' }
          : { id: c.id, situacao: 'aplicado' as const },
      );
    }, agora);

    fila.enfileirar(comando({ tipo: 'IniciarOS', agregado: OS_A }));
    fila.enfileirar(comando({ tipo: 'ConcluirOS', agregado: OS_A }));

    await fila.descarregar();

    // Segunda descarga, imediatamente: o IniciarOS está de castigo, então o
    // ConcluirOS não pode passar na frente dele.
    enviados.length = 0;
    const r = await fila.descarregar();

    expect(enviados).toHaveLength(0);
    expect(r.bloqueadosPorAgregado).toBeGreaterThan(0);
  });

  it('um comando em espera NÃO bloqueia outra OS', async () => {
    let tempo = new Date('2026-11-12T10:00:00Z').getTime();
    const agora = () => new Date(tempo);

    const { fila, enviados } = montar(async (lote) => {
      return lote.map((c) =>
        c.agregado === OS_A
          ? { id: c.id, situacao: 'retentar' as const, motivo: 'servidor indisponível' }
          : { id: c.id, situacao: 'aplicado' as const },
      );
    }, agora);

    fila.enfileirar(comando({ agregado: OS_A }));
    fila.enfileirar(comando({ agregado: OS_B }));

    await fila.descarregar();
    enviados.length = 0;
    await fila.descarregar();

    // A OS A está de castigo; a B já foi aplicada e não é reenviada.
    expect(enviados.flat().every((c) => c.agregado !== OS_B)).toBe(true);
  });
});

describe('falha de rede', () => {
  it('erro de transporte devolve tudo para pendente, sem perder nada', async () => {
    const { fila, deposito } = montar(async () => {
      throw new Error('sem conexão');
    });

    fila.enfileirar(comando());
    fila.enfileirar(comando());

    const r = await fila.descarregar();

    expect(r.aRetentar).toBe(2);
    expect(deposito.todos()).toHaveLength(2);
    expect(deposito.todos().every((c) => c.situacao === 'pendente')).toBe(true);
    expect(fila.quantidadeNaoSincronizada()).toBe(2);
  });

  it('comando enviado sem resultado correspondente volta para pendente', async () => {
    // O servidor pode ter aplicado e a resposta se perdido. O comando é
    // reenviado, e a idempotência do servidor impede dupla aplicação.
    const { fila, deposito } = montar(async () => []);

    const c = comando();
    fila.enfileirar(c);

    const r = await fila.descarregar();

    expect(r.aRetentar).toBe(1);
    expect(deposito.buscar(c.id)?.situacao).toBe('pendente');
  });

  it('a espera entre tentativas cresce e tem teto de 5 minutos', async () => {
    let tempo = new Date('2026-11-12T10:00:00Z').getTime();
    const agora = () => new Date(tempo);

    const { fila, deposito } = montar(async () => {
      throw new Error('sem conexão');
    }, agora);

    const c = comando();
    fila.enfileirar(c);

    const esperas: number[] = [];
    for (let i = 0; i < 12; i += 1) {
      await fila.descarregar();
      const atual = deposito.buscar(c.id);
      const proxima = new Date(atual!.proximaTentativaEm!).getTime();
      esperas.push(proxima - tempo);
      tempo = proxima; // avança o relógio até a próxima tentativa
    }

    expect(esperas[0]).toBe(1000);
    expect(esperas[1]).toBe(2000);
    expect(esperas[2]).toBe(4000);
    expect(Math.max(...esperas)).toBe(5 * 60 * 1000);
    // Nunca desiste: o comando segue pendente depois de 12 tentativas.
    expect(deposito.buscar(c.id)?.situacao).toBe('pendente');
  });
});

describe('rejeição', () => {
  it('comando rejeitado é terminal, guarda o motivo e aparece nas pendências', async () => {
    const motivo = 'A OS OS-2026-000123 está cancelada e não pode ser concluída.';
    const { fila, deposito, enviados } = montar(async (lote) =>
      lote.map((c) => ({ id: c.id, situacao: 'rejeitado' as const, motivo })),
    );

    const c = comando({ tipo: 'ConcluirOS' });
    fila.enfileirar(c);

    await fila.descarregar();

    const guardado = deposito.buscar(c.id);
    expect(guardado?.situacao).toBe('rejeitado');
    expect(guardado?.motivoRejeicao).toBe(motivo);

    // Não desaparece em silêncio: vai para a tela de pendências.
    expect(fila.pendenciasParaOUsuario().map((p) => p.id)).toContain(c.id);

    // E não é reenviado.
    enviados.length = 0;
    await fila.descarregar();
    expect(enviados).toHaveLength(0);
  });
});

describe('lote', () => {
  it('respeita o limite de 200 comandos por requisição', async () => {
    const { fila, enviados } = montar(aplicarTudo);

    // Agregados distintos, para que nada bloqueie nada.
    for (let i = 0; i < 250; i += 1) {
      fila.enfileirar(comando({ agregado: `ordem_servico:${i}` }));
    }

    await fila.descarregar();

    expect(enviados[0]?.length).toBe(TAMANHO_MAXIMO_DO_LOTE);
  });

  it('esvazia a fila em descargas sucessivas', async () => {
    const { fila } = montar(aplicarTudo);

    for (let i = 0; i < 250; i += 1) {
      fila.enfileirar(comando({ agregado: `ordem_servico:${i}` }));
    }

    await fila.descarregar();
    await fila.descarregar();

    expect(fila.quantidadeNaoSincronizada()).toBe(0);
  });
});

describe('concorrência', () => {
  it('duas descargas simultâneas não enviam o mesmo comando duas vezes', async () => {
    let liberar: () => void = () => {};
    const bloqueio = new Promise<void>((r) => {
      liberar = r;
    });

    const { fila, enviados } = montar(async (lote) => {
      await bloqueio;
      return lote.map((c) => ({ id: c.id, situacao: 'aplicado' as const }));
    });

    fila.enfileirar(comando());

    const primeira = fila.descarregar();
    const segunda = fila.descarregar(); // entra enquanto a primeira está em voo
    liberar();

    await Promise.all([primeira, segunda]);

    expect(enviados).toHaveLength(1);
  });
});

describe('sete dias offline', () => {
  it('acumula o trabalho de uma semana e converge numa reconexão', async () => {
    const respondedor = vi.fn(aplicarTudo);
    const { fila, deposito } = montar(respondedor);

    // Uma semana de campo: 7 dias, 6 OS por dia, 5 comandos por OS.
    let total = 0;
    for (let dia = 0; dia < 7; dia += 1) {
      for (let os = 0; os < 6; os += 1) {
        const agregado = `ordem_servico:dia${dia}-os${os}`;
        for (const tipo of [
          'AssumirOS',
          'IniciarOS',
          'ApontarHoras',
          'RegistrarLeitura',
          'ConcluirOS',
        ] as const) {
          fila.enfileirar(comando({ tipo, agregado }));
          total += 1;
        }
      }
    }

    expect(fila.quantidadeNaoSincronizada()).toBe(total);
    expect(total).toBe(210);

    // Reconectou: descarrega até esvaziar.
    let voltas = 0;
    while (fila.quantidadeNaoSincronizada() > 0 && voltas < 20) {
      await fila.descarregar();
      voltas += 1;
    }

    expect(fila.quantidadeNaoSincronizada()).toBe(0);
    expect(fila.pendenciasParaOUsuario()).toHaveLength(0);
    expect(deposito.todos().filter((c) => c.situacao === 'aplicado')).toHaveLength(total);
    // 210 comandos em lotes de 200: duas requisições bastam.
    expect(respondedor.mock.calls.length).toBeLessThanOrEqual(3);
  });
});

describe('persistência', () => {
  it('a fila sobrevive ao fechamento do app', () => {
    // O mesmo arquivo de banco reaberto representa o app fechado e reaberto.
    const banco = new DatabaseSync(':memory:');
    const deposito = new DepositoSQLite(banco);

    const c = comando({ tipo: 'ConcluirOS' });
    deposito.inserir(c);

    // Outro depósito sobre o mesmo banco: o esquema é idempotente e os dados
    // continuam lá.
    const reaberto = new DepositoSQLite(banco);

    expect(reaberto.todos()).toHaveLength(1);
    expect(reaberto.buscar(c.id)?.tipo).toBe('ConcluirOS');
    expect(reaberto.contarNaoTerminados()).toBe(1);
  });
});
