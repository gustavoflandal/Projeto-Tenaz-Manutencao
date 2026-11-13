import { describe, expect, it } from 'vitest';

import { es419 } from './es-419';
import { glossario, regrasDoGlossario } from './glossario';
import { ptBR } from './pt-BR';

function folhas(obj: object, prefixo = ''): Array<[string, string]> {
  return Object.entries(obj).flatMap(([chave, valor]) =>
    typeof valor === 'string'
      ? [[`${prefixo}${chave}`, valor] as [string, string]]
      : folhas(valor as object, `${prefixo}${chave}.`),
  );
}

// Categorias de plural: o espanhol acrescenta `many` às do português.
const SUFIXO_PLURAL = /_(zero|one|two|few|many|other)$/;
const base = (chave: string): string => chave.replace(SUFIXO_PLURAL, '');
const marcadores = (texto: string): string[] =>
  [...texto.matchAll(/\{\{\s*([^}\s]+)\s*\}\}/g)].map((m) => m[1] ?? '').sort();

// `\b` do JavaScript não reconhece letra acentuada; troca por olhar-em-volta.
function compilar(padrao: string): RegExp {
  const convertido = padrao.replace(/\\b/g, (_m: string, pos: number) => {
    const antes = pos === 0 ? '(' : (padrao[pos - 1] ?? '(');
    return antes === '(' || antes === '|' ? '(?<![\\p{L}\\p{N}_])' : '(?![\\p{L}\\p{N}_])';
  });
  return new RegExp(convertido, 'iu');
}

// Chaves que existem só em pt-BR, por decisão do dono do produto: o aviso fixo de
// dado de saúde não vai em espanhol (fica como customização, se o cliente pedir).
// Qualquer outra chave faltando em es-419 continua quebrando o build.
const SO_PT_BR: readonly string[] = ['pessoa.avisoSaude'];

const pt = folhas(ptBR);
// Marcadores {{...}} são variáveis, não texto: ficam fora das regras de glossário.
const semMarcadores = (texto: string): string => texto.replace(/\{\{[^}]*\}\}/g, ' ');
const es = folhas(es419);
const ptPorBase = new Map<string, string[]>();
for (const [k, v] of pt) ptPorBase.set(base(k), [...(ptPorBase.get(base(k)) ?? []), v]);
const esPorChave = new Map(es);
const ptPorChave = new Map(pt);

describe('paridade de chaves entre pt-BR e es-419', () => {
  it('a lista de chaves só do pt-BR é pequena e todas existem em pt-BR e faltam em es-419', () => {
    expect(SO_PT_BR).toEqual(['pessoa.avisoSaude']);
    for (const k of SO_PT_BR) {
      expect(ptPorChave.has(k)).toBe(true);
      expect(esPorChave.has(k)).toBe(false);
    }
  });

  it('toda chave de pt-BR existe em es-419', () => {
    const faltam = pt.map(([k]) => k).filter((k) => !SO_PT_BR.includes(k)).filter((k) => !esPorChave.has(k) && !esPorChave.has(`${base(k)}_many`) && SUFIXO_PLURAL.test(k) === false);
    expect(faltam).toEqual([]);
  });

  it('toda chave de es-419 existe em pt-BR (exceto _many de plural)', () => {
    const sobram = es
      .map(([k]) => k)
      .filter((k) => !ptPorChave.has(k) && !(k.endsWith('_many') && ptPorChave.has(`${base(k)}_other`)));
    expect(sobram).toEqual([]);
  });

  it('todo plural de pt-BR tem one, many e other em es-419', () => {
    const faltam: string[] = [];
    for (const [k] of pt) {
      if (!k.endsWith('_other')) continue;
      for (const s of ['one', 'many', 'other']) {
        if (!esPorChave.has(`${base(k)}_${s}`)) faltam.push(`${base(k)}_${s}`);
      }
    }
    expect(faltam).toEqual([]);
  });

  it.each(es)('%s mantém os mesmos marcadores {{...}} do pt-BR', (chave, texto) => {
    const origem = ptPorChave.get(chave) ?? ptPorChave.get(`${base(chave)}_other`);
    expect(origem).toBeDefined();
    expect(marcadores(texto)).toEqual(marcadores(origem ?? ''));
  });
});

describe('tom de voz em espanhol (seção 10.5 do documento 7)', () => {
  it.each(es)('%s não é vazio', (_chave, texto) => {
    expect(texto.trim()).not.toBe('');
  });

  it.each(es)('%s não usa exclamação', (_chave, texto) => {
    expect(texto).not.toMatch(/[!¡]/);
  });

  it.each(es)('%s não usa anglicismos proibidos', (_chave, texto) => {
    expect(texto.toLowerCase()).not.toMatch(/(?<![\p{L}])(ticket|offline|tenant)(?![\p{L}])/u);
  });
});

describe('glossário es-419', () => {
  it('cada termo tem fonte e traduções', () => {
    for (const g of glossario) {
      expect(g.fonte.length, g.id).toBeGreaterThan(0);
      expect(g.es.trim(), g.id).not.toBe('');
      expect(g.pt.trim(), g.id).not.toBe('');
    }
  });

  it('os ids do glossário são únicos', () => {
    const ids = glossario.map((g) => g.id);
    expect(new Set(ids).size).toBe(ids.length);
  });

  const proibicoes = regrasDoGlossario.flatMap((g) =>
    (g.esProibidos ?? []).map((p) => [g.id, p] as [string, string]),
  );
  it.each(proibicoes)('%s: nenhum texto usa %s', (_id, padrao) => {
    const re = compilar(padrao);
    const achados = es.filter(([, texto]) => re.test(semMarcadores(texto))).map(([k, t]) => `${k}: ${t}`);
    expect(achados).toEqual([]);
  });

  const consistencias = regrasDoGlossario.filter((g) => g.ptDetecta && g.esExige);
  it.each(consistencias.map((g) => [g.id, g] as const))('%s: onde o pt-BR usa o termo, o es-419 usa o termo aprovado', (_id, g) => {
    const detecta = compilar(g.ptDetecta ?? '');
    const exige = compilar(g.esExige ?? '');
    const erros: string[] = [];
    for (const [chave, textoPt] of pt) {
      if (!detecta.test(textoPt)) continue;
      const textoEs = esPorChave.get(chave);
      if (textoEs !== undefined && !exige.test(textoEs)) erros.push(`${chave}: ${textoEs}`);
    }
    expect(erros).toEqual([]);
  });
});
