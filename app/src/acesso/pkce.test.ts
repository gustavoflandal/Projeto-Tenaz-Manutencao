import { describe, expect, it } from 'vitest';

import { base64Url, desafioS256, gerarPKCE, uuidV4 } from './pkce';
import { CLIENT_ID, interpretarRetorno, lerConsulta, URI_DE_RETORNO, urlDeAutorizacao } from './retorno';
import { criptoDeTeste } from './testes/servidor-falso';

describe('PKCE S256', () => {
  it('reproduz o vetor da RFC 7636, apêndice B', async () => {
    const desafio = await desafioS256(criptoDeTeste, 'dBjftJeZ4CVP-mB92K27uhbUJU1p1r_wW1gFWFOEjXk');
    expect(desafio).toBe('E9Melhoa2OwvFrEMTJguCHaoeK1t8URWbuGJSstw-cM');
  });

  it('gera verificador de 43 caracteres do alfabeto permitido, e state diferente a cada tentativa', async () => {
    const a = await gerarPKCE(criptoDeTeste);
    const b = await gerarPKCE(criptoDeTeste);
    expect(a.verificador).toMatch(/^[A-Za-z0-9_-]{43}$/);
    expect(a.desafio).toMatch(/^[A-Za-z0-9_-]{43}$/);
    expect(a.verificador).not.toBe(b.verificador);
    expect(a.state).not.toBe(b.state);
  });

  it('codifica base64url sem preenchimento, igual ao Node', () => {
    for (const tamanho of [1, 2, 3, 4, 16, 32]) {
      const bytes = criptoDeTeste.bytesAleatorios(tamanho);
      expect(base64Url(bytes)).toBe(Buffer.from(bytes).toString('base64url'));
    }
  });

  it('gera UUID v4 válido para o dispositivo', () => {
    expect(uuidV4(criptoDeTeste)).toMatch(/^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$/);
  });
});

describe('URL de autorização e retorno', () => {
  it('monta a URL com os parâmetros da seção 4.3', () => {
    const url = urlDeAutorizacao('https://api.tenaz.com.br/', 'DESAFIO', 'st');
    const q = lerConsulta(url);
    expect(url.startsWith('https://api.tenaz.com.br/oauth/authorize?')).toBe(true);
    expect(q.get('response_type')).toBe('code');
    expect(q.get('client_id')).toBe(CLIENT_ID);
    expect(q.get('redirect_uri')).toBe('br.com.tenaz.app:/oauth');
    expect(q.get('code_challenge')).toBe('DESAFIO');
    expect(q.get('code_challenge_method')).toBe('S256');
    expect(q.get('state')).toBe('st');
  });

  it('aceita o retorno com código e o state da tentativa', () => {
    expect(interpretarRetorno(`${URI_DE_RETORNO}?code=abc%2B1&state=st`, 'st')).toEqual({ tipo: 'codigo', codigo: 'abc+1' });
  });

  it('recusa retorno com state de outra tentativa', () => {
    expect(interpretarRetorno(`${URI_DE_RETORNO}?code=abc&state=outro`, 'st').tipo).toBe('erro');
  });

  it('recusa retorno sem state', () => {
    expect(interpretarRetorno(`${URI_DE_RETORNO}?code=abc`, 'st').tipo).toBe('erro');
  });

  it('recusa URI de retorno que não é a do app', () => {
    expect(interpretarRetorno('https://mal.example/oauth?code=abc&state=st', 'st').tipo).toBe('erro');
    expect(interpretarRetorno('br.com.tenaz.app:/outro?code=abc&state=st', 'st').tipo).toBe('erro');
  });

  it('entende access_denied do locatário', () => {
    expect(interpretarRetorno(`${URI_DE_RETORNO}?error=access_denied&state=st`, 'st')).toEqual({ tipo: 'negado' });
  });

  it('recusa retorno sem código nem erro', () => {
    expect(interpretarRetorno(`${URI_DE_RETORNO}?state=st`, 'st').tipo).toBe('erro');
  });
});
