// Adaptador de `expo-sqlite` para a interface Banco da fila de comandos.
//
// O motor da fila (fila.ts) não conhece Expo nem Node: ele fala com a
// interface Banco. Isso é o que permite testar a lógica contra SQLite real no
// Node, como manda a seção 9.1 de docs/5_ARQUITETURA_SOFTWARE.md, e rodar o
// MESMO código no aparelho. Se o teste e o aparelho usassem motores
// diferentes, o teste não provaria nada sobre o campo.
//
// Um banco por conta (empresa + usuário), cada um com a própria chave no
// Keychain ou Keystore, a própria fila e o próprio cursor. Nenhuma tabela
// local mistura duas contas (seção 10.8).

import type { SQLiteDatabase } from 'expo-sqlite';

import type { Banco, Declaracao } from './deposito-sqlite';

/**
 * Nome do arquivo de banco de uma conta.
 *
 * O nome carrega empresa e usuário de propósito: dois acessos do mesmo
 * aparelho nunca compartilham arquivo. Sair da conta e sessão encerrada NÃO
 * apagam este arquivo: só as credenciais saem, e a fila de comandos fica
 * guardada até a mesma conta entrar de novo. Hoje não existe operação que
 * apague o banco local; se ela vier, deve recusar enquanto houver comandos
 * não sincronizados (FilaDeComandos.quantidadeNaoSincronizada() > 0).
 */
export function arquivoDaConta(empresaID: string, usuarioID: string): string {
  return `tenaz-${empresaID}-${usuarioID}.db`;
}

class DeclaracaoExpo implements Declaracao {
  constructor(
    private readonly banco: SQLiteDatabase,
    private readonly sql: string,
  ) {}

  run(...parametros: unknown[]): unknown {
    return this.banco.runSync(this.sql, parametros as never);
  }

  all(...parametros: unknown[]): unknown[] {
    return this.banco.getAllSync(this.sql, parametros as never);
  }

  get(...parametros: unknown[]): unknown {
    return this.banco.getFirstSync(this.sql, parametros as never) ?? undefined;
  }
}

/** Envolve um banco do expo-sqlite na interface que a fila espera. */
export class BancoExpo implements Banco {
  constructor(private readonly banco: SQLiteDatabase) {}

  exec(sql: string): void {
    this.banco.execSync(sql);
  }

  prepare(sql: string): Declaracao {
    return new DeclaracaoExpo(this.banco, sql);
  }
}
