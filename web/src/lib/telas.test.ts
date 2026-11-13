import { readFileSync } from 'node:fs';
import { resolve } from 'node:path';

import { describe, expect, it } from 'vitest';

import { TELAS_DO_CATALOGO, telaDoCaminho } from './telas';

describe('catálogo de telas', () => {
  it('é igual ao do backend (catalogo.go), na mesma ordem', () => {
    const go = readFileSync(resolve(__dirname, '../../../backend/internal/plataforma/auditoria/catalogo.go'), 'utf8');
    const codigos = [...go.matchAll(/\{"([a-z]+\.[a-z_]+)", "[a-z]+", "/g)].map((m) => m[1]);
    expect(codigos).toEqual([...TELAS_DO_CATALOGO]);
  });
});

describe('telaDoCaminho', () => {
  it.each([
    ['/', 'ind.painel'],
    ['/administracao/empresa', 'adm.empresa'],
    ['/administracao/unidades', 'adm.unidades'],
    ['/administracao/unidades/nova', 'adm.unidades'],
    ['/administracao/unidades/abc', 'adm.unidade_detalhe'],
    ['/administracao/pessoas/pe-1', 'adm.pessoa_detalhe'],
    ['/administracao/usuarios/novo', 'adm.usuarios'],
    ['/administracao/usuarios/u-1', 'adm.usuario_detalhe'],
    ['/administracao/perfis', 'adm.perfis'],
    ['/administracao/auditoria', 'adm.auditoria'],
    ['/administracao/auditoria/integridade', 'adm.auditoria'],
    ['/administracao/auditoria/capturas', 'adm.parametros'],
    ['/conta', null],
    ['/entrar', null],
  ])('%s é %s', (caminho, esperado) => {
    expect(telaDoCaminho(caminho)).toBe(esperado);
  });
});
