import { describe, expect, it } from 'vitest';

import casos from '../../../backend/testes/casos/identificacao.json';
import {
  cpaConfereComSubdivisao,
  mascararEntrada,
  mascararParaLista,
  normalizarIdentificacao,
  PERFIS_DE_PAIS,
  primeiraSobreposicao,
  senhaTemTamanho,
  SUBDIVISOES,
  subdivisaoValida,
  sobrepoem,
  turnosDoModelo,
  validarCodigoPostal,
  validarIdentificacao,
  validarTurno,
} from './index';

describe('identificação (casos compartilhados com o Go, versão 3)', () => {
  it('usa a versão de casos que conhece e o total de cada grupo (como o Go)', () => {
    expect(casos.versao).toBe(3);
    expect(casos.normalizacao).toHaveLength(7);
    expect(casos.validacao).toHaveLength(34);
    expect(casos.mascara).toHaveLength(6);
    expect(casos.perfis).toHaveLength(3);
  });

  it.each(casos.normalizacao)('normaliza: $nome', ({ entrada, esperado }) => {
    expect(normalizarIdentificacao(entrada)).toBe(esperado);
  });

  it.each(casos.validacao)('valida: $nome', ({ tipo, entrada, normalizado, valido, motivo }) => {
    const r = validarIdentificacao(tipo, entrada);
    expect(r.valido).toBe(valido);
    if (r.valido) expect(r.normalizado).toBe(normalizado);
    else expect(r.motivo).toBe(motivo);
  });

  it.each(casos.mascara)('máscara de lista: $nome', ({ tipo, entrada, esperado }) => {
    expect(mascararParaLista(tipo, entrada)).toBe(esperado);
  });

  it('os perfis por país são os do servidor', () => {
    for (const caso of casos.perfis) {
      const perfil = PERFIS_DE_PAIS[caso.pais as keyof typeof PERFIS_DE_PAIS];
      expect(perfil.moeda).toBe(caso.moeda);
      expect(perfil.idioma_padrao).toBe(caso.idioma_padrao);
      expect(perfil.fuso_padrao).toBe(caso.fuso_padrao);
      expect(perfil.tipo_empresa).toBe(caso.tipo_empresa);
      expect([...perfil.tipos_pessoa]).toEqual(caso.tipos_pessoa);
      expect([...perfil.fusos]).toEqual(caso.fusos);
    }
    expect(Object.keys(PERFIS_DE_PAIS)).toHaveLength(casos.perfis.length);
  });
});

describe('máscara de entrada', () => {
  it('desenha cada tipo enquanto a pessoa digita', () => {
    expect(mascararEntrada('cnpj', '12abc34501de35')).toBe('12.ABC.345/01DE-35');
    expect(mascararEntrada('cnpj', '1122')).toBe('11.22');
    expect(mascararEntrada('cpf', '52998224725')).toBe('529.982.247-25');
    expect(mascararEntrada('cuit', '33693450239')).toBe('33-69345023-9');
    expect(mascararEntrada('cuil', '2012')).toBe('20-12');
    expect(mascararEntrada('dni', '12345678')).toBe('12.345.678');
    expect(mascararEntrada('dni', '1234567')).toBe('1.234.567');
    expect(mascararEntrada('ci', '1234')).toBe('1.234');
    expect(mascararEntrada('ruc', '800123456')).toBe('80012345-6');
  });

  it('corta o que passa do tamanho e é idempotente', () => {
    expect(mascararEntrada('dni', '1234567890')).toBe('12.345.678');
    const m = mascararEntrada('cnpj', '11222333000181');
    expect(mascararEntrada('cnpj', m)).toBe(m);
  });
});

describe('turnos (RN-ADM-019)', () => {
  const t = (dia: number, ini: string, fim: string) => ({ dia_semana: dia, hora_inicio: ini, hora_fim: fim });

  it('valida horários', () => {
    expect(validarTurno(t(1, '08:00', '17:00'))).toBeNull();
    expect(validarTurno(t(1, '00:00', '24:00'))).toBeNull();
    expect(validarTurno(t(1, '24:00', '08:00'))).toBe('hora_inicio_invalida');
    expect(validarTurno(t(1, '8:00', '17:00'))).toBe('hora_inicio_invalida');
    expect(validarTurno(t(1, '08:00', '25:00'))).toBe('hora_fim_invalida');
    expect(validarTurno(t(1, '08:00', '08:00'))).toBe('fim_igual_ao_inicio');
  });

  it('turnos encostados não se sobrepõem', () => {
    expect(sobrepoem(t(1, '08:00', '12:00'), t(1, '12:00', '17:00'))).toBe(false);
  });

  it('turno que atravessa a meia-noite invade o dia seguinte', () => {
    expect(sobrepoem(t(1, '22:00', '06:00'), t(2, '05:00', '09:00'))).toBe(true);
    expect(sobrepoem(t(1, '22:00', '06:00'), t(2, '06:00', '09:00'))).toBe(false);
  });

  it('o de domingo termina na segunda', () => {
    expect(sobrepoem(t(7, '22:00', '06:00'), t(1, '05:00', '09:00'))).toBe(true);
  });

  it('acha o primeiro par que se sobrepõe', () => {
    expect(primeiraSobreposicao([t(1, '08:00', '12:00'), t(2, '08:00', '12:00'), t(1, '11:00', '13:00')])).toEqual([0, 2]);
    expect(primeiraSobreposicao([t(1, '08:00', '12:00')])).toBeNull();
  });

  it('os modelos não têm sobreposição e têm os dias previstos', () => {
    const comercial = turnosDoModelo('comercial');
    const continuo = turnosDoModelo('continuo');
    expect(comercial).toHaveLength(5);
    expect(continuo).toHaveLength(7);
    expect(primeiraSobreposicao(comercial)).toBeNull();
    expect(primeiraSobreposicao(continuo)).toBeNull();
    expect(continuo.every((x) => validarTurno(x) === null)).toBe(true);
  });
});

describe('senha', () => {
  it('exige 10 caracteres', () => {
    expect(senhaTemTamanho('123456789')).toBe(false);
    expect(senhaTemTamanho('1234567890')).toBe(true);
  });
});

describe('endereço multipaís (casos compartilhados, versão 3)', () => {
  const paisDe = (p: string) => p as keyof typeof SUBDIVISOES;

  it('nenhum grupo vem vazio', () => {
    expect(casos.subdivisoes).toHaveLength(3);
    expect(casos.subdivisao_invalida).toHaveLength(12);
    expect(casos.codigo_postal).toHaveLength(16);
    expect(casos.codigo_postal_cpa).toHaveLength(12);
  });

  it.each(casos.subdivisoes)('lista de subdivisões de $pais', ({ pais, lista }) => {
    expect(SUBDIVISOES[paisDe(pais)]).toEqual(lista);
    expect(PERFIS_DE_PAIS[paisDe(pais)].subdivisoes).toEqual(lista);
  });

  it.each(casos.subdivisao_invalida)('recusa a subdivisão %j em qualquer país', (codigo) => {
    for (const pais of ['BR', 'AR', 'PY'] as const) expect(subdivisaoValida(pais, codigo)).toBe(false);
  });

  it.each(casos.codigo_postal)('código postal: $nome', ({ pais, entrada, normalizado, valido }) => {
    const r = validarCodigoPostal(paisDe(pais), entrada);
    expect(r.valido).toBe(valido);
    if (r.valido) expect(r.normalizado).toBe(normalizado);
  });

  it.each(casos.codigo_postal_cpa)('CPA e subdivisão: $nome', ({ subdivisao, codigo_postal, confere }) => {
    expect(cpaConfereComSubdivisao(subdivisao, codigo_postal)).toBe(confere);
  });
});
