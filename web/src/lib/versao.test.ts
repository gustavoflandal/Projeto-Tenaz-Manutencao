import { afterEach, describe, expect, it, vi } from 'vitest';

import { mesmaVersao, nomeAcessivelDoSelo, seloDeVersao, versaoDaAPI } from './versao';

describe('selo de versão', () => {
  it('mostra só a versão em produção', () => {
    expect(
      seloDeVersao({ versao: '1.4.2', commit: 'a1b2c3d', dataBuild: '', ambiente: 'producao' }),
    ).toBe('v1.4.2');
  });

  it('mostra o ambiente fora da produção, para captura de tela não ser confundida', () => {
    expect(
      seloDeVersao({ versao: '1.5.0-rc.1', commit: 'a1b2c3d', dataBuild: '', ambiente: 'homologacao' }),
    ).toBe('v1.5.0-rc.1 · Homologação');
  });

  it('tem o nome acessível do componente 6.19', () => {
    expect(
      nomeAcessivelDoSelo({ versao: '1.4.2', commit: 'a1b2c3d', dataBuild: '', ambiente: 'producao' }),
    ).toBe('Versão 1.4.2. Abrir informações sobre o sistema');
    expect(
      nomeAcessivelDoSelo({ versao: '1.5.0-rc.1', commit: 'a1b2c3d', dataBuild: '', ambiente: 'homologacao' }),
    ).toBe('Versão 1.5.0-rc.1, ambiente Homologação. Abrir informações sobre o sistema');
  });

  it('não mostra o metadado de build, que fica no diálogo Sobre', () => {
    const v = { versao: '1.5.0-rc.1+a1b2c3d', commit: 'a1b2c3d', dataBuild: '', ambiente: 'homologacao' };
    expect(seloDeVersao(v)).toBe('v1.5.0-rc.1 · Homologação');
    expect(nomeAcessivelDoSelo(v)).toBe('Versão 1.5.0-rc.1, ambiente Homologação. Abrir informações sobre o sistema');
  });

  it('compara versões ignorando o metadado de build', () => {
    expect(mesmaVersao('1.4.2+a1b2c3d', '1.4.2+e4f5a6b')).toBe(true);
    expect(mesmaVersao('1.4.2', '1.4.3')).toBe(false);
  });

  it('usa a versão injetada no build, nunca texto digitado', () => {
    // Sem TENAZ_VERSAO no ambiente, o vite.config.ts injeta 0.0.0-dev.
    expect(seloDeVersao()).toBe('v0.0.0-dev · Desenvolvimento');
  });
});

describe('versão da API', () => {
  afterEach(() => {
    vi.restoreAllMocks();
  });

  it('consulta GET /v1/versao pelo cliente gerado', async () => {
    const fetch = vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      Response.json({ versao: '1.4.2', commit: 'a1b2c3d', data_build: '2026-11-12T14:05:00Z' }),
    );
    await expect(versaoDaAPI()).resolves.toEqual({
      versao: '1.4.2',
      commit: 'a1b2c3d',
      data_build: '2026-11-12T14:05:00Z',
    });
    const requisicao = fetch.mock.calls[0]?.[0] as Request;
    expect(new URL(requisicao.url).pathname).toBe('/v1/versao');
    expect(requisicao.method).toBe('GET');
  });

  it('falha quando a API responde com erro', async () => {
    vi.spyOn(globalThis, 'fetch').mockResolvedValue(
      new Response(JSON.stringify({ type: 'about:blank', title: 'Indisponível', status: 503, codigo: 'INDISPONIVEL' }), {
        status: 503,
        headers: { 'Content-Type': 'application/problem+json' },
      }),
    );
    await expect(versaoDaAPI()).rejects.toThrow('503');
  });
});
