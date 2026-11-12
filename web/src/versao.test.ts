import { describe, expect, it } from 'vitest';

import { seloDeVersao } from './versao';

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
});
