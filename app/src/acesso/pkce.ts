// PKCE (RFC 7636) do login do app, método S256.
//
// A criptografia vem de fora (interface Cripto): no aparelho é expo-crypto; nos
// testes, node:crypto. Assim a lógica é a mesma nos dois lados e o vetor de
// teste da RFC 7636 (apêndice B) prova o cálculo.

/** O mínimo de criptografia de que o login precisa. */
export interface Cripto {
  bytesAleatorios(quantidade: number): Uint8Array;
  sha256(entrada: string): Promise<Uint8Array>;
}

const ALFABETO = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_';

/** Base64 URL sem preenchimento (RFC 4648, seção 5). Não usa Buffer: não existe no Hermes. */
export function base64Url(bytes: Uint8Array): string {
  let saida = '';
  for (let i = 0; i < bytes.length; i += 3) {
    const b0 = bytes[i] ?? 0;
    const b1 = bytes[i + 1];
    const b2 = bytes[i + 2];
    saida += ALFABETO[b0 >> 2];
    saida += ALFABETO[((b0 & 0x03) << 4) | ((b1 ?? 0) >> 4)];
    if (b1 !== undefined) saida += ALFABETO[((b1 & 0x0f) << 2) | ((b2 ?? 0) >> 6)];
    if (b2 !== undefined) saida += ALFABETO[b2 & 0x3f];
  }
  return saida;
}

export interface ParametrosPKCE {
  /** Segredo de uso único que fica no aparelho até a troca do código. */
  verificador: string;
  /** SHA-256 do verificador em base64url: é o que vai na URL de autorização. */
  desafio: string;
  /** Valor opaco que volta na URL de retorno e prova que a resposta é desta tentativa. */
  state: string;
}

/** Verificador com 32 bytes aleatórios: 43 caracteres, o mínimo da RFC e dentro do limite de 128 do contrato. */
export async function gerarPKCE(cripto: Cripto): Promise<ParametrosPKCE> {
  const verificador = base64Url(cripto.bytesAleatorios(32));
  const desafio = await desafioS256(cripto, verificador);
  const state = base64Url(cripto.bytesAleatorios(16));
  return { verificador, desafio, state };
}

export async function desafioS256(cripto: Cripto, verificador: string): Promise<string> {
  return base64Url(await cripto.sha256(verificador));
}

/** UUID v4 do dispositivo, para o cabeçalho X-Tenaz-Dispositivo-Id. */
export function uuidV4(cripto: Cripto): string {
  const b = cripto.bytesAleatorios(16).slice();
  b[6] = ((b[6] ?? 0) & 0x0f) | 0x40;
  b[8] = ((b[8] ?? 0) & 0x3f) | 0x80;
  const h = Array.from(b, (x) => x.toString(16).padStart(2, '0')).join('');
  return `${h.slice(0, 8)}-${h.slice(8, 12)}-${h.slice(12, 16)}-${h.slice(16, 20)}-${h.slice(20)}`;
}
