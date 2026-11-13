// Cofre das credenciais: o armazenamento seguro do sistema (Keychain no iOS,
// Keystore no Android), por trás de uma interface para os testes usarem memória.
//
// Credencial nunca vai para AsyncStorage, SQLite, arquivo ou log. Só o cofre.

export interface Cofre {
  ler(chave: string): Promise<string | null>;
  gravar(chave: string, valor: string): Promise<void>;
  apagar(chave: string): Promise<void>;
}

/** Cofre em memória, só para testes. */
export class CofreEmMemoria implements Cofre {
  readonly itens = new Map<string, string>();

  async ler(chave: string): Promise<string | null> {
    return this.itens.get(chave) ?? null;
  }
  async gravar(chave: string, valor: string): Promise<void> {
    this.itens.set(chave, valor);
  }
  async apagar(chave: string): Promise<void> {
    this.itens.delete(chave);
  }
}
