// Apoio dos testes E2E: documentos válidos, TOTP (RFC 6238) e leitura do
// e-mail no Mailpit. Nada daqui vai para o produto.

import { createHmac, randomInt } from 'node:crypto';

export const MAILPIT = process.env.E2E_MAILPIT ?? 'http://localhost:8025';

function digitoVerificador(base: number[], pesos: number[]): number {
  const soma = base.reduce((acc, d, i) => acc + d * (pesos[i] ?? 0), 0);
  const resto = soma % 11;
  return resto < 2 ? 0 : 11 - resto;
}

/** CNPJ válido, aleatório, só com dígitos (14). */
export function gerarCNPJ(): string {
  const base = Array.from({ length: 8 }, () => randomInt(0, 10)).concat([0, 0, 0, 1]);
  const d1 = digitoVerificador(base, [5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2]);
  const d2 = digitoVerificador([...base, d1], [6, 5, 4, 3, 2, 9, 8, 7, 6, 5, 4, 3, 2]);
  return [...base, d1, d2].join('');
}

/** CUIT de 11 dígitos, aleatório (o servidor confere só a forma, até a P13). */
export function gerarCUIT(): string {
  return `30${Array.from({ length: 9 }, () => randomInt(0, 10)).join('')}`;
}

/** RUC de 9 dígitos, aleatório (o servidor confere só a forma, até a P13). */
export function gerarRUC(): string {
  return `80${Array.from({ length: 7 }, () => randomInt(0, 10)).join('')}`;
}

/** CI de 7 dígitos, aleatória (primeiro dígito diferente de zero). */
export function gerarCI(): string {
  return `${randomInt(1, 10)}${Array.from({ length: 6 }, () => randomInt(0, 10)).join('')}`;
}

/** DNI de 8 dígitos, aleatório (primeiro dígito diferente de zero). */
export function gerarDNI(): string {
  return `${randomInt(1, 10)}${Array.from({ length: 7 }, () => randomInt(0, 10)).join('')}`;
}

/** CPF válido, aleatório, só com dígitos (11). */
export function gerarCPF(): string {
  const base = Array.from({ length: 9 }, () => randomInt(0, 10));
  const dv = (digitos: number[]) => {
    const soma = digitos.reduce((acc, d, i) => acc + d * (digitos.length + 1 - i), 0);
    const r = (soma * 10) % 11;
    return r === 10 ? 0 : r;
  };
  const d1 = dv(base);
  const d2 = dv([...base, d1]);
  return [...base, d1, d2].join('');
}

const ALFABETO_BASE32 = 'ABCDEFGHIJKLMNOPQRSTUVWXYZ234567';

function base32ParaBytes(texto: string): Buffer {
  let bits = '';
  for (const c of texto.replaceAll(/[\s=]/g, '').toUpperCase()) {
    const valor = ALFABETO_BASE32.indexOf(c);
    if (valor < 0) throw new Error(`Caractere fora do base32: ${c}`);
    bits += valor.toString(2).padStart(5, '0');
  }
  const bytes: number[] = [];
  for (let i = 0; i + 8 <= bits.length; i += 8) bytes.push(Number.parseInt(bits.slice(i, i + 8), 2));
  return Buffer.from(bytes);
}

const PASSO_DO_TOTP = 30;

/** Código TOTP de 6 dígitos (HMAC-SHA1, passo de 30 s) do segredo em base32. */
export function totp(segredoBase32: string, agora = Date.now()): string {
  const contador = Math.floor(agora / 1000 / PASSO_DO_TOTP);
  const mensagem = Buffer.alloc(8);
  mensagem.writeBigUInt64BE(BigInt(contador));
  const hmac = createHmac('sha1', base32ParaBytes(segredoBase32)).update(mensagem).digest();
  const deslocamento = (hmac[hmac.length - 1] ?? 0) & 0x0f;
  const codigo = (hmac.readUInt32BE(deslocamento) & 0x7fffffff) % 1_000_000;
  return String(codigo).padStart(6, '0');
}

/** Espera o próximo passo do TOTP: o servidor recusa um passo já usado. */
export async function esperarProximoPasso(): Promise<void> {
  const passoAtual = Math.floor(Date.now() / 1000 / PASSO_DO_TOTP);
  while (Math.floor(Date.now() / 1000 / PASSO_DO_TOTP) === passoAtual) {
    await new Promise((resolver) => setTimeout(resolver, 500));
  }
}

interface MensagemDoMailpit {
  readonly ID: string;
}

/** Espera o e-mail enviado a `para` e devolve o `#token=` do link dele. */
export async function tokenDoEmail(para: string, tempoMaximoMs = 30_000): Promise<string> {
  const limite = Date.now() + tempoMaximoMs;
  while (Date.now() < limite) {
    const busca = await fetch(`${MAILPIT}/api/v1/search?query=${encodeURIComponent(`to:${para}`)}`);
    const { messages } = (await busca.json()) as { messages?: MensagemDoMailpit[] };
    const primeira = messages?.[0];
    if (primeira) {
      const mensagem = (await (await fetch(`${MAILPIT}/api/v1/message/${primeira.ID}`)).json()) as { Text?: string; HTML?: string };
      const corpo = `${mensagem.Text ?? ''}\n${mensagem.HTML ?? ''}`;
      const achado = /#token=([A-Za-z0-9_-]{43})/.exec(corpo);
      if (achado?.[1]) return achado[1];
    }
    await new Promise((resolver) => setTimeout(resolver, 500));
  }
  throw new Error(`Nenhum e-mail com link para ${para} no Mailpit (${MAILPIT}).`);
}
