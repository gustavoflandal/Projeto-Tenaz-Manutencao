// Política de senha (RF-ADM-006): o servidor também confere senha comum ou
// vazada e senha igual ao e-mail (422 senha_fraca). Aqui só o tamanho.

export const SENHA_TAMANHO_MINIMO = 10;

export function senhaTemTamanho(senha: string): boolean {
  return [...senha].length >= SENHA_TAMANHO_MINIMO;
}
