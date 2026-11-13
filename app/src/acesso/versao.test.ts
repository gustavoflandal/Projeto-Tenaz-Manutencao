import { describe, expect, it } from 'vitest';

import { tradutor } from '../i18n';
import { ErroAcessoNegado, ErroApi, ErroContaMudou, ErroDeRede, ErroSessaoEncerrada, chaveDoErro } from './erros';
import {
  formatarDataDoBuild,
  lerInformacoes,
  linhasDoSobre,
  rotuloAcessivelDoSelo,
  textoDoSelo,
  textoParaCopiar,
} from './versao';

const pt = tradutor('pt-BR');
const es = tradutor('es-419');
const extra = { versaoRelease: '1.5.0-rc.1', commit: 'a1b2c3d', dataBuild: '2026-11-12T14:05:00Z', ambiente: 'homologacao' };

describe('selo e Sobre', () => {
  it('mostra a versão do release embutido, com sufixo de pré-lançamento', () => {
    const i = lerInformacoes(extra, { versao: '1.4.0', build: '87' }, 'upd-123');
    expect(textoDoSelo(pt, i)).toBe('v1.5.0-rc.1'); // a que roda, não a do binário da loja
    expect(i.binario).toEqual({ versao: '1.4.0', build: '87' });
    expect(i.atualizacaoId).toBe('upd-123');
  });

  it('nome acessível em produção não cita o ambiente', () => {
    const i = lerInformacoes({ ...extra, versaoRelease: '1.4.2', ambiente: 'producao' }, { versao: null, build: null }, null);
    expect(rotuloAcessivelDoSelo(pt, i)).toBe('Versão 1.4.2. Abrir informações sobre o sistema');
    expect(rotuloAcessivelDoSelo(es, i)).toBe('Versión 1.4.2. Abrir información sobre el sistema');
  });

  it('nome acessível fora da produção cita o ambiente, no idioma da tela', () => {
    const i = lerInformacoes(extra, { versao: null, build: null }, null);
    expect(rotuloAcessivelDoSelo(pt, i)).toBe('Versão 1.5.0-rc.1, ambiente Homologação. Abrir informações sobre o sistema');
    expect(rotuloAcessivelDoSelo(es, i)).toBe('Versión 1.5.0-rc.1, ambiente Homologación. Abrir información sobre el sistema');
  });

  it('sem extra, usa valores explícitos de desenvolvimento e nunca inventa versão', () => {
    const i = lerInformacoes(undefined, { versao: null, build: null }, null);
    expect(i).toMatchObject({ versaoRelease: '0.0.0-dev', commit: 'desconhecido', ambiente: 'desenvolvimento', binario: null });
  });

  it('formata a data do build pelo país e pelo fuso, nunca pelo idioma, e trata data ausente ou ilegível', () => {
    const local = { pais: 'BR', fuso: 'America/Sao_Paulo' } as const;
    expect(formatarDataDoBuild(pt, '2026-11-12T14:05:00Z', local)).toBe('12/11/2026 11:05');
    expect(formatarDataDoBuild(es, '2026-11-12T14:05:00Z', local)).toBe('12/11/2026 11:05');
    expect(formatarDataDoBuild(pt, '2026-11-12T14:05:00Z', { pais: 'BR', fuso: 'UTC' })).toBe('12/11/2026 14:05');
    expect(formatarDataDoBuild(pt, '')).toBe('Não informada');
    expect(formatarDataDoBuild(es, '')).toBe('No informada');
    expect(formatarDataDoBuild(pt, 'ontem')).toBe('Não informada');
  });

  it('monta as linhas do Sobre na ordem do documento 7 e o texto de "Copiar informações"', () => {
    const i = lerInformacoes(extra, { versao: '1.4.0', build: '87' }, null);
    const linhas = linhasDoSobre(pt, i, { estado: 'ok', versao: '1.5.0-rc.1' });
    expect(linhas.map((l) => l.rotulo)).toEqual([
      'App',
      'Binário da loja',
      'Atualização sem loja',
      'API',
      'Data do build',
      'Commit',
      'Ambiente',
    ]);
    expect(linhas[1]?.valor).toBe('v1.4.0 (build 87)');
    expect(linhas[2]?.valor).toBe('Nenhuma, roda o pacote do binário');
    expect(linhas[3]?.valor).toBe('v1.5.0-rc.1');
    const copia = textoParaCopiar(pt, linhas);
    expect(copia.split('\n')).toHaveLength(7);
    expect(copia).toContain('Commit: a1b2c3d');
    expect(copia).toContain('Ambiente: Homologação');
  });

  it('sem rede a API aparece como indisponível, não como erro', () => {
    const i = lerInformacoes(extra, { versao: null, build: null }, null);
    expect(linhasDoSobre(pt, i, { estado: 'indisponivel' })[3]?.valor).toBe('Indisponível sem conexão');
    expect(linhasDoSobre(es, i, { estado: 'indisponivel' })[3]?.valor).toBe('No disponible sin conexión');
    expect(linhasDoSobre(pt, i, { estado: 'consultando' })[3]?.valor).toBe('Consultando a API');
  });
});

describe('chave do texto de erro', () => {
  it('orienta o segundo fator na web quando há mfa_pendente, nos dois idiomas', () => {
    const k = chaveDoErro(new ErroApi(401, 'nao_autenticado', null, true, 'x'));
    expect(k).toBe('mobile.erros.mfaPendente');
    expect(pt(k)).toContain('segundo fator');
    expect(pt(k)).not.toMatch(/MFA/);
    expect(es(k)).toContain('segundo factor');
  });

  it('mapeia cada erro conhecido para uma chave que existe nos dois idiomas', () => {
    const erros = [
      new ErroDeRede(),
      new ErroAcessoNegado(),
      new ErroSessaoEncerrada(),
      new ErroContaMudou(),
      new ErroApi(404, 'x', null, false, ''),
      new ErroApi(429, 'x', null, false, ''),
      new ErroApi(503, 'x', null, false, ''),
      new ErroApi(400, 'x', null, false, ''),
      new Error('qualquer'),
    ];
    for (const e of erros) {
      const k = chaveDoErro(e);
      expect(pt(k), k).not.toBe(k);
      expect(es(k), k).not.toBe(k);
      expect(es(k), k).not.toBe(pt(k));
    }
  });
});
