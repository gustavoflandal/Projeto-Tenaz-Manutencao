// Fila de comandos (outbox) do app de campo da Tenaz.
//
// É o componente mais crítico do app e o maior risco técnico do projeto
// (risco 3 do cronograma, em docs/6_ROADMAP.md). A meta é operar 7 dias sem
// conexão sem perder registro, e convergir em menos de 60 s depois da
// reconexão (seção 10.1 de docs/5_ARQUITETURA_SOFTWARE.md).
//
// Invariantes que esta fila garante, e que os testes trancam:
//
//  1. Nenhum comando é perdido. Nem por falha de rede, nem por rejeição, nem
//     por fechamento do app no meio do envio.
//  2. Reenviar não duplica: o id do comando é a chave de idempotência, e
//     gravar o mesmo id duas vezes é sem efeito.
//  3. Ordem FIFO POR AGREGADO. Comandos da mesma OS saem na ordem em que
//     foram criados; comandos de OS diferentes não se bloqueiam.
//  4. Um comando que falha bloqueia os seguintes DO MESMO agregado, nunca os
//     de outro.
//  5. Comando rejeitado é terminal, guarda o motivo e não some.
//
// O depósito é uma interface para que a lógica seja testada contra SQLite real
// no Node (seção 9.1) e rode sobre expo-sqlite no aparelho.

import {
  type Comando,
  type ComandoNaFila,
  type ResultadoComando,
  TAMANHO_MAXIMO_DO_LOTE,
} from './comando';

/** Armazenamento da fila. Implementado por SQLite no app e nos testes. */
export interface DepositoDaFila {
  /** Grava o comando. Se o id já existe, não faz nada (idempotência local). */
  inserir(c: Comando): void;
  /** Comandos pendentes, em ordem de criação. */
  pendentes(): ComandoNaFila[];
  /** Um comando pelo id. */
  buscar(id: string): ComandoNaFila | null;
  /** Troca a situação de um comando. */
  marcar(id: string, mudanca: MudancaDeSituacao): void;
  /** Todos os comandos, em ordem de criação. Para a tela de pendências. */
  todos(): ComandoNaFila[];
  /** Quantos comandos ainda não terminaram. */
  contarNaoTerminados(): number;
}

export interface MudancaDeSituacao {
  readonly situacao: ComandoNaFila['situacao'];
  readonly motivoRejeicao?: string | null;
  readonly tentativas?: number;
  readonly proximaTentativaEm?: string | null;
}

/** Quem envia o lote ao servidor. */
export interface Transporte {
  enviar(lote: readonly Comando[]): Promise<readonly ResultadoComando[]>;
}

export interface OpcoesDaFila {
  readonly deposito: DepositoDaFila;
  readonly transporte: Transporte;
  /** Relógio injetável, para o teste controlar o tempo. */
  readonly agora?: () => Date;
  readonly tamanhoDoLote?: number;
  /** Tentativas antes de o comando ficar em espera longa. */
  readonly maximoDeTentativas?: number;
}

export interface ResultadoDaDescarga {
  readonly enviados: number;
  readonly aplicados: number;
  readonly rejeitados: number;
  readonly aRetentar: number;
  readonly bloqueadosPorAgregado: number;
}

export class FilaDeComandos {
  private readonly deposito: DepositoDaFila;
  private readonly transporte: Transporte;
  private readonly agora: () => Date;
  private readonly tamanhoDoLote: number;
  private readonly maximoDeTentativas: number;
  private descarregando = false;

  constructor(opcoes: OpcoesDaFila) {
    this.deposito = opcoes.deposito;
    this.transporte = opcoes.transporte;
    this.agora = opcoes.agora ?? (() => new Date());
    this.tamanhoDoLote = Math.min(
      opcoes.tamanhoDoLote ?? TAMANHO_MAXIMO_DO_LOTE,
      TAMANHO_MAXIMO_DO_LOTE,
    );
    this.maximoDeTentativas = opcoes.maximoDeTentativas ?? 10;
  }

  /**
   * Enfileira um comando. Chamado logo depois de a tela aplicar a mudança
   * otimista no banco local.
   *
   * Gravar o mesmo id duas vezes é sem efeito: é o que torna seguro repetir a
   * ação quando o app é fechado entre a gravação local e a da fila.
   */
  enfileirar(c: Comando): void {
    this.deposito.inserir(c);
  }

  /** Comandos que o usuário precisa ver em "Pendências de sincronização". */
  pendenciasParaOUsuario(): ComandoNaFila[] {
    return this.deposito.todos().filter((c) => c.situacao === 'rejeitado');
  }

  /** Quantos comandos ainda não chegaram ao servidor. */
  quantidadeNaoSincronizada(): number {
    return this.deposito.contarNaoTerminados();
  }

  /**
   * Envia o que está pendente, respeitando a ordem por agregado.
   *
   * Não roda duas vezes em paralelo: o app chama a descarga na reconexão, ao
   * voltar do segundo plano e por temporizador, e essas chamadas se cruzam.
   * Dois envios simultâneos do mesmo comando não causariam duplicidade no
   * servidor, por causa da idempotência, mas embaralhariam a ordem.
   */
  async descarregar(): Promise<ResultadoDaDescarga> {
    if (this.descarregando) {
      return { enviados: 0, aplicados: 0, rejeitados: 0, aRetentar: 0, bloqueadosPorAgregado: 0 };
    }
    this.descarregando = true;
    try {
      return await this.descarregarUmaVez();
    } finally {
      this.descarregando = false;
    }
  }

  private async descarregarUmaVez(): Promise<ResultadoDaDescarga> {
    const pendentes = this.deposito.pendentes();
    const { lote, bloqueados } = this.montarLote(pendentes);

    if (lote.length === 0) {
      return {
        enviados: 0,
        aplicados: 0,
        rejeitados: 0,
        aRetentar: 0,
        bloqueadosPorAgregado: bloqueados,
      };
    }

    for (const c of lote) {
      this.deposito.marcar(c.id, { situacao: 'enviando' });
    }

    let resultados: readonly ResultadoComando[];
    try {
      resultados = await this.transporte.enviar(lote.map(semEstado));
    } catch (erro) {
      // Falha de transporte: tudo volta para pendente. Nada se perde.
      const motivo = erro instanceof Error ? erro.message : String(erro);
      for (const c of lote) {
        this.registrarRetentativa(c, motivo);
      }
      return {
        enviados: lote.length,
        aplicados: 0,
        rejeitados: 0,
        aRetentar: lote.length,
        bloqueadosPorAgregado: bloqueados,
      };
    }

    let aplicados = 0;
    let rejeitados = 0;
    let aRetentar = 0;

    const porID = new Map(resultados.map((r) => [r.id, r]));

    for (const c of lote) {
      const r = porID.get(c.id);

      // Comando enviado sem resultado correspondente: o servidor pode ter
      // aplicado. Volta para pendente e é reenviado; a idempotência do
      // servidor garante que não aplique duas vezes.
      if (!r) {
        this.registrarRetentativa(c, 'servidor não respondeu sobre este comando');
        aRetentar += 1;
        continue;
      }

      switch (r.situacao) {
        case 'aplicado':
          this.deposito.marcar(c.id, { situacao: 'aplicado', motivoRejeicao: null });
          aplicados += 1;
          break;
        case 'rejeitado':
          // Decisão de domínio: terminal, com o motivo guardado para o usuário.
          this.deposito.marcar(c.id, { situacao: 'rejeitado', motivoRejeicao: r.motivo });
          rejeitados += 1;
          break;
        case 'retentar':
          this.registrarRetentativa(c, r.motivo);
          aRetentar += 1;
          break;
      }
    }

    return {
      enviados: lote.length,
      aplicados,
      rejeitados,
      aRetentar,
      bloqueadosPorAgregado: bloqueados,
    };
  }

  /**
   * Monta o lote preservando a ordem por agregado.
   *
   * Para cada agregado, entra no lote apenas a sequência inicial contínua de
   * comandos prontos. Se o primeiro comando de uma OS está em espera de
   * retentativa, nenhum comando posterior daquela OS entra — senão um
   * `ConcluirOS` chegaria antes do `IniciarOS` que o precede.
   */
  private montarLote(pendentes: readonly ComandoNaFila[]): {
    lote: ComandoNaFila[];
    bloqueados: number;
  } {
    const agora = this.agora().getTime();
    const lote: ComandoNaFila[] = [];
    const agregadoTrancado = new Set<string>();
    let bloqueados = 0;

    for (const c of pendentes) {
      if (agregadoTrancado.has(c.agregado)) {
        bloqueados += 1;
        continue;
      }

      const esperando =
        c.proximaTentativaEm !== null && new Date(c.proximaTentativaEm).getTime() > agora;

      if (esperando) {
        // Este comando ainda está de castigo: tranca o agregado dele.
        agregadoTrancado.add(c.agregado);
        bloqueados += 1;
        continue;
      }

      if (lote.length >= this.tamanhoDoLote) {
        agregadoTrancado.add(c.agregado);
        bloqueados += 1;
        continue;
      }

      lote.push(c);
    }

    return { lote, bloqueados };
  }

  private registrarRetentativa(c: ComandoNaFila, motivo: string): void {
    const tentativas = c.tentativas + 1;
    const espera = this.esperaDaTentativa(tentativas);
    const proxima = new Date(this.agora().getTime() + espera).toISOString();

    this.deposito.marcar(c.id, {
      situacao: 'pendente',
      tentativas,
      proximaTentativaEm: proxima,
      motivoRejeicao: motivo,
    });
  }

  /**
   * Recuo exponencial com teto de 5 minutos.
   *
   * Nunca desiste: um comando com muitas tentativas continua pendente, só
   * espera mais. Desistir seria perder trabalho de campo, o que esta fila não
   * faz em nenhuma circunstância.
   */
  private esperaDaTentativa(tentativas: number): number {
    const limitada = Math.min(tentativas, this.maximoDeTentativas);
    const base = 1000 * 2 ** (limitada - 1);
    return Math.min(base, 5 * 60 * 1000);
  }
}

function semEstado(c: ComandoNaFila): Comando {
  return {
    id: c.id,
    tipo: c.tipo,
    agregado: c.agregado,
    carga: c.carga,
    versaoBase: c.versaoBase,
    horarioAparelho: c.horarioAparelho,
  };
}
