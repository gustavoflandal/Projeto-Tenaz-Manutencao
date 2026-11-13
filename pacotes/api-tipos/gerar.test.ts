import { readFileSync } from 'node:fs';

import { describe, expect, it } from 'vitest';

import { DESTINO, gerarTipos } from './gerar.mjs';

describe('tipos da API', () => {
  it('src/esquema.ts está em dia com backend/api/openapi.yaml', async () => {
    const versionado = readFileSync(DESTINO, 'utf8').replace(/\r\n/g, '\n');
    expect(versionado, 'rode: bash ferramentas/tarefas.sh api-tipos').toBe(await gerarTipos());
  });
});
