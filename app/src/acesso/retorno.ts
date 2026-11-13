// Montagem da URL de autorização e leitura da URL de retorno (seção 4.3 de
// docs/3_ESPECIFICACAO_APIS.md). Sem `URL`: o do Hermes não implementa tudo, e
// o esquema `br.com.tenaz.app:/oauth` não tem host.

export const CLIENT_ID = 'tenaz-app';
export const URI_DE_RETORNO = 'br.com.tenaz.app:/oauth';

export function urlDeAutorizacao(apiUrl: string, desafio: string, state: string): string {
  const pares: [string, string][] = [
    ['response_type', 'code'],
    ['client_id', CLIENT_ID],
    ['redirect_uri', URI_DE_RETORNO],
    ['code_challenge', desafio],
    ['code_challenge_method', 'S256'],
    ['state', state],
  ];
  const q = pares.map(([k, v]) => `${k}=${encodeURIComponent(v)}`).join('&');
  return `${apiUrl.replace(/\/+$/, '')}/oauth/authorize?${q}`;
}

/** Lê `a=b&c=d` com decodificação; chave repetida vale a primeira. */
export function lerConsulta(url: string): Map<string, string> {
  const m = new Map<string, string>();
  const i = url.indexOf('?');
  if (i < 0) return m;
  const resto = url.slice(i + 1).split('#')[0] ?? '';
  for (const par of resto.split('&')) {
    if (!par) continue;
    const j = par.indexOf('=');
    const k = decodificar(j < 0 ? par : par.slice(0, j));
    const v = j < 0 ? '' : decodificar(par.slice(j + 1));
    if (!m.has(k)) m.set(k, v);
  }
  return m;
}

function decodificar(s: string): string {
  try {
    return decodeURIComponent(s.replace(/\+/g, ' '));
  } catch {
    return s;
  }
}

export type RetornoDoLogin =
  | { tipo: 'codigo'; codigo: string }
  | { tipo: 'negado' }
  | { tipo: 'erro'; motivo: string };

/**
 * Interpreta a URL de retorno. O `state` precisa ser o desta tentativa; sem
 * isso, uma resposta de outra tentativa (ou forjada) trocaria o código por
 * tokens de outra conta.
 */
export function interpretarRetorno(url: string, stateEsperado: string): RetornoDoLogin {
  if (!url.startsWith(`${URI_DE_RETORNO}?`)) {
    return { tipo: 'erro', motivo: 'URL de retorno desconhecida.' };
  }
  const q = lerConsulta(url);
  if (q.get('state') !== stateEsperado) {
    return { tipo: 'erro', motivo: 'state não confere.' };
  }
  const erro = q.get('error');
  if (erro === 'access_denied') return { tipo: 'negado' };
  if (erro) return { tipo: 'erro', motivo: erro };
  const codigo = q.get('code');
  if (!codigo) return { tipo: 'erro', motivo: 'sem código.' };
  return { tipo: 'codigo', codigo };
}
