import { describe, expect, it } from 'vitest';

import { ErroDaAPI } from './api';
import { aplicarErrosDoServidor, caminhoDoCampo, mensagemDoErro } from './erros';
import { tokenDoFragmento } from './fragmento';
import { retornoSeguro } from './retorno';

const erro = (status: number, codigo: string, extras: object = {}) =>
  new ErroDaAPI(codigo, status, { type: 'x', title: 'x', status, codigo, ...extras });

describe('mensagemDoErro', () => {
  it('falha de rede vira "sem conexão"', () => {
    expect(mensagemDoErro(new TypeError('Failed to fetch'))).toMatch(/^Sem conexão/);
  });

  it('usa o texto do código estável', () => {
    expect(mensagemDoErro(erro(409, 'ultimo_administrador'))).toMatch(/ao menos um administrador ativo/);
  });

  it('conta as tentativas restantes do segundo fator', () => {
    expect(mensagemDoErro(erro(401, 'mfa_codigo_invalido', { tentativas_restantes: 1 }))).toMatch(/Resta 1 tentativa/);
    expect(mensagemDoErro(erro(401, 'mfa_codigo_invalido', { tentativas_restantes: 4 }))).toMatch(/Restam 4 tentativas/);
  });

  it('5xx traz o código de rastreio para o suporte', () => {
    expect(mensagemDoErro(erro(500, 'erro_interno', { id_rastreio: 'E-4F21' }))).toMatch(/E-4F21/);
  });

  it('código desconhecido cai no texto do status', () => {
    expect(mensagemDoErro(erro(404, 'qualquer_coisa'))).toMatch(/Não encontramos este registro/);
  });

  it('validação sem campo mostra o que o servidor disse', () => {
    expect(
      mensagemDoErro(erro(422, 'validacao_falhou', { erros: [{ campo: 'turnos', codigo: 'x', mensagem: 'Informe um turno.' }] })),
    ).toBe('Informe um turno.');
  });
});

describe('erros por campo', () => {
  it('converte a notação de índice', () => {
    expect(caminhoDoCampo('unidade.turnos[0].hora_fim')).toBe('unidade.turnos.0.hora_fim');
  });

  it('leva cada erro para o campo, com prefixo, mapa e transformação', () => {
    const chamadas: Array<[string, string | undefined]> = [];
    const setError = (nome: string, e: { message?: string }) => chamadas.push([nome, e.message]);
    const problema = erro(422, 'validacao_falhou', {
      erros: [
        { campo: 'senha', codigo: 'x', mensagem: 'A' },
        { campo: 'turnos[1]', codigo: 'x', mensagem: 'B' },
      ],
    });
    expect(aplicarErrosDoServidor(problema, setError as never, { prefixo: 'p', mapa: { senha: 'nova_senha' } })).toBe(true);
    expect(chamadas).toEqual([
      ['p.nova_senha', 'A'],
      ['p.turnos.1', 'B'],
    ]);
    expect(aplicarErrosDoServidor(new Error('x'), setError as never)).toBe(false);
  });
});

describe('retornoSeguro', () => {
  it.each([
    ['/administracao/pessoas?busca=a', '/administracao/pessoas?busca=a'],
    ['//outro.site', '/'],
    ['https://outro.site', '/'],
    ['/\\outro.site', '/'],
    ['/entrar', '/'],
    [null, '/'],
    ['', '/'],
  ])('%s vira %s', (entrada, esperado) => {
    expect(retornoSeguro(entrada)).toBe(esperado);
  });
});

describe('tokenDoFragmento', () => {
  it('lê o token de 43 caracteres', () => {
    expect(tokenDoFragmento(`#token=${'a'.repeat(43)}`)).toBe('a'.repeat(43));
  });

  it('recusa token ausente ou de outro tamanho', () => {
    expect(tokenDoFragmento('')).toBeNull();
    expect(tokenDoFragmento('#token=curto')).toBeNull();
    expect(tokenDoFragmento('#outro=' + 'a'.repeat(43))).toBeNull();
  });
});
