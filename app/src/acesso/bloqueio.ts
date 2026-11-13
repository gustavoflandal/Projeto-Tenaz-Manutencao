// Desbloqueio local por biometria ou PIN do aparelho ("proposta" em RF-MOB-001).
//
// Só protege o app aberto no aparelho de quem está logado; não substitui o
// login nem altera a credencial. A decisão de quando pedir é pura e testada; a
// pergunta ao sistema (expo-local-authentication) fica em bloqueio-expo.ts.

import type { Cofre } from './cofre';

const CHAVE = 'tenaz.bloqueio-local';

/** Depois de tanto tempo em segundo plano, o app pede o desbloqueio de novo. */
export const TEMPO_PARA_BLOQUEAR_MS = 5 * 60_000;

export interface AutenticadorLocal {
  /** Há biometria ou PIN cadastrado no aparelho? */
  disponivel(): Promise<boolean>;
  /** Mostra o prompt do sistema. Verdadeiro se a pessoa passou. */
  autenticar(motivo: string): Promise<boolean>;
}

export class BloqueioLocal {
  private ativoNoCofre = false;
  private bloqueado = false;
  private foiParaSegundoPlanoEm: number | null = null;

  constructor(
    private readonly cofre: Cofre,
    private readonly autenticador: AutenticadorLocal,
    private readonly agora: () => number = Date.now,
  ) {}

  /** Lê a preferência e já parte bloqueado na abertura a frio, quando ativa. */
  async iniciar(): Promise<void> {
    this.ativoNoCofre = (await this.cofre.ler(CHAVE)) === '1';
    this.bloqueado = this.ativoNoCofre;
  }

  ativo(): boolean {
    return this.ativoNoCofre;
  }

  estaBloqueado(): boolean {
    return this.ativoNoCofre && this.bloqueado;
  }

  /** Liga só se o aparelho tem com o que desbloquear e a pessoa passa uma vez agora. */
  async ativar(motivo: string): Promise<boolean> {
    if (!(await this.autenticador.disponivel())) return false;
    if (!(await this.autenticador.autenticar(motivo))) return false;
    await this.cofre.gravar(CHAVE, '1');
    this.ativoNoCofre = true;
    this.bloqueado = false;
    return true;
  }

  async desativar(): Promise<void> {
    await this.cofre.apagar(CHAVE);
    this.ativoNoCofre = false;
    this.bloqueado = false;
  }

  /** O app foi para segundo plano. */
  aoIrParaSegundoPlano(): void {
    this.foiParaSegundoPlanoEm = this.agora();
  }

  /** O app voltou ao primeiro plano. */
  aoVoltar(): void {
    if (this.foiParaSegundoPlanoEm !== null && this.agora() - this.foiParaSegundoPlanoEm >= TEMPO_PARA_BLOQUEAR_MS) {
      this.bloqueado = true;
    }
    this.foiParaSegundoPlanoEm = null;
  }

  async desbloquear(motivo: string): Promise<boolean> {
    if (!this.ativoNoCofre) return true;
    const ok = await this.autenticador.autenticar(motivo);
    if (ok) this.bloqueado = false;
    return ok;
  }
}
