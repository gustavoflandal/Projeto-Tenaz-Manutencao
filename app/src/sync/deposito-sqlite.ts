// Depósito da fila de comandos em SQLite.
//
// O mesmo SQL roda nos testes, contra `node:sqlite`, e no aparelho, contra
// `expo-sqlite`: a interface Banco abaixo é o pouco que as duas têm em comum.
// Testar o motor de sincronização com SQLite real no Node é o que a seção 9.1
// de docs/5_ARQUITETURA_SOFTWARE.md prescreve.
//
// Um banco por conta (empresa + usuário), com outbox e cursor próprios: a
// troca de conta fecha um banco e abre outro, e nenhuma tabela local mistura
// duas contas (seção 10.8).

import type { Comando, ComandoNaFila, SituacaoComando, TipoDeComando } from './comando';
import type { DepositoDaFila, MudancaDeSituacao } from './fila';

/** O mínimo que precisamos de um SQLite, presente nas duas implementações. */
export interface Banco {
  exec(sql: string): void;
  prepare(sql: string): Declaracao;
}

export interface Declaracao {
  run(...parametros: unknown[]): unknown;
  all(...parametros: unknown[]): unknown[];
  get(...parametros: unknown[]): unknown;
}

export const ESQUEMA_DA_FILA = `
CREATE TABLE IF NOT EXISTS comando_pendente (
  id                  TEXT    PRIMARY KEY,
  sequencia           INTEGER NOT NULL,
  tipo                TEXT    NOT NULL,
  agregado            TEXT    NOT NULL,
  carga               TEXT    NOT NULL,
  versao_base         INTEGER,
  horario_aparelho    TEXT    NOT NULL,
  situacao            TEXT    NOT NULL DEFAULT 'pendente',
  tentativas          INTEGER NOT NULL DEFAULT 0,
  motivo_rejeicao     TEXT,
  proxima_tentativa_em TEXT,
  CHECK (situacao IN ('pendente', 'enviando', 'aplicado', 'rejeitado'))
);

-- A sequência é a ordem de criação e define o FIFO.
CREATE UNIQUE INDEX IF NOT EXISTS ix_comando_sequencia
  ON comando_pendente (sequencia);

-- Consulta do caminho quente: pendentes em ordem, por agregado.
CREATE INDEX IF NOT EXISTS ix_comando_pendente
  ON comando_pendente (situacao, agregado, sequencia);
`;

interface LinhaComando {
  id: string;
  sequencia: number;
  tipo: string;
  agregado: string;
  carga: string;
  versao_base: number | null;
  horario_aparelho: string;
  situacao: string;
  tentativas: number;
  motivo_rejeicao: string | null;
  proxima_tentativa_em: string | null;
}

export class DepositoSQLite implements DepositoDaFila {
  constructor(private readonly banco: Banco) {
    this.banco.exec(ESQUEMA_DA_FILA);
  }

  inserir(c: Comando): void {
    // INSERT OR IGNORE é a idempotência local: gravar o mesmo id duas vezes
    // não cria dois comandos. É o que torna seguro repetir a ação quando o
    // app fecha entre a gravação na tela e a gravação na fila.
    this.banco
      .prepare(
        `INSERT OR IGNORE INTO comando_pendente
           (id, sequencia, tipo, agregado, carga, versao_base, horario_aparelho)
         VALUES
           (?, (SELECT coalesce(max(sequencia), 0) + 1 FROM comando_pendente), ?, ?, ?, ?, ?)`,
      )
      .run(
        c.id,
        c.tipo,
        c.agregado,
        JSON.stringify(c.carga),
        c.versaoBase,
        c.horarioAparelho,
      );
  }

  pendentes(): ComandoNaFila[] {
    const linhas = this.banco
      .prepare(
        `SELECT * FROM comando_pendente
          WHERE situacao IN ('pendente', 'enviando')
          ORDER BY sequencia`,
      )
      .all() as LinhaComando[];
    return linhas.map(paraComando);
  }

  buscar(id: string): ComandoNaFila | null {
    const linha = this.banco
      .prepare(`SELECT * FROM comando_pendente WHERE id = ?`)
      .get(id) as LinhaComando | undefined;
    return linha ? paraComando(linha) : null;
  }

  marcar(id: string, mudanca: MudancaDeSituacao): void {
    const atual = this.buscar(id);
    if (!atual) return;

    this.banco
      .prepare(
        `UPDATE comando_pendente
            SET situacao = ?,
                motivo_rejeicao = ?,
                tentativas = ?,
                proxima_tentativa_em = ?
          WHERE id = ?`,
      )
      .run(
        mudanca.situacao,
        mudanca.motivoRejeicao === undefined ? atual.motivoRejeicao : mudanca.motivoRejeicao,
        mudanca.tentativas === undefined ? atual.tentativas : mudanca.tentativas,
        mudanca.proximaTentativaEm === undefined
          ? atual.proximaTentativaEm
          : mudanca.proximaTentativaEm,
        id,
      );
  }

  todos(): ComandoNaFila[] {
    const linhas = this.banco
      .prepare(`SELECT * FROM comando_pendente ORDER BY sequencia`)
      .all() as LinhaComando[];
    return linhas.map(paraComando);
  }

  contarNaoTerminados(): number {
    const linha = this.banco
      .prepare(
        `SELECT count(*) AS total FROM comando_pendente
          WHERE situacao IN ('pendente', 'enviando')`,
      )
      .get() as { total: number };
    return linha.total;
  }
}

function paraComando(l: LinhaComando): ComandoNaFila {
  return {
    id: l.id,
    sequencia: l.sequencia,
    tipo: l.tipo as TipoDeComando,
    agregado: l.agregado,
    carga: JSON.parse(l.carga) as Record<string, unknown>,
    versaoBase: l.versao_base,
    horarioAparelho: l.horario_aparelho,
    situacao: l.situacao as SituacaoComando,
    tentativas: l.tentativas,
    motivoRejeicao: l.motivo_rejeicao,
    proximaTentativaEm: l.proxima_tentativa_em,
  };
}
