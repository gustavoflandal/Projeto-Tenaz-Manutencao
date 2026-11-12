// Tipos da fila de comandos do app de campo.
//
// O app envia INTENÇÕES de domínio, não diferenças de linha: "ConcluirOS",
// "RegistrarLeitura", "ApontarHoras". O servidor aplica cada comando pelos
// mesmos casos de uso da API pública — não existe caminho de escrita exclusivo
// do app (seção 10.3 de docs/5_ARQUITETURA_SOFTWARE.md).

/** Tipos de comando previstos para o MVP (seção 10.3). */
export const TIPOS_DE_COMANDO = [
  'AssumirOS',
  'IniciarOS',
  'PausarOS',
  'ApontarHoras',
  'RegistrarLeitura',
  'ResponderItemChecklist',
  'RegistrarComponente',
  'ConcluirOS',
  'AbrirOS',
  'AbrirSolicitacao',
  'AnexarArquivo',
  'Assinar',
  'ColetarAssinaturaCliente',
] as const;

export type TipoDeComando = (typeof TIPOS_DE_COMANDO)[number];

/**
 * Situação de um comando na fila.
 *
 * `rejeitado` é terminal e NUNCA desaparece em silêncio: fica na fila com o
 * motivo, para o app mostrar em "Pendências de sincronização" (seção 10.5).
 */
export type SituacaoComando = 'pendente' | 'enviando' | 'aplicado' | 'rejeitado';

/** Comando como o app o grava e como o envia ao servidor. */
export interface Comando {
  /** UUIDv7 gerado no aparelho. É a chave de idempotência. */
  readonly id: string;
  readonly tipo: TipoDeComando;
  /**
   * Agregado a que o comando pertence: a OS, o ativo, a solicitação.
   * A ordem é FIFO POR AGREGADO — comandos de OS diferentes não se bloqueiam.
   */
  readonly agregado: string;
  readonly carga: Readonly<Record<string, unknown>>;
  /**
   * Versão do registro que o usuário viu quando agiu. Serve para resolver
   * conflito de campo editável (seção 10.5).
   */
  readonly versaoBase: number | null;
  /** Horário do APARELHO. O servidor grava também o de recebimento, e a
   * diferença fica registrada: relógio errado não apaga a evidência. */
  readonly horarioAparelho: string;
}

/** Comando com o estado que a fila mantém. */
export interface ComandoNaFila extends Comando {
  readonly situacao: SituacaoComando;
  readonly tentativas: number;
  readonly sequencia: number;
  readonly motivoRejeicao: string | null;
  readonly proximaTentativaEm: string | null;
}

/** Resultado que o servidor devolve para cada comando do lote. */
export type ResultadoComando =
  | { readonly id: string; readonly situacao: 'aplicado' }
  | { readonly id: string; readonly situacao: 'rejeitado'; readonly motivo: string }
  /**
   * Falha de transporte ou erro do servidor: o comando continua pendente e é
   * reenviado. Diferente de `rejeitado`, que é decisão de domínio e terminal.
   */
  | { readonly id: string; readonly situacao: 'retentar'; readonly motivo: string };

/** Limite de comandos por requisição (seção 10.3). */
export const TAMANHO_MAXIMO_DO_LOTE = 200;
