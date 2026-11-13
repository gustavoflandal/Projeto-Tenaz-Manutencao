import { describe, expect, it } from 'vitest';

import { categoriaDeReserva, escolherIdioma, idiomaDoAparelho, traduzir } from './index';

describe('idioma do aparelho e escolha da pessoa', () => {
  it('segue a primeira língua do aparelho que a Tenaz atende', () => {
    expect(idiomaDoAparelho(['pt-BR'])).toBe('pt-BR');
    expect(idiomaDoAparelho(['pt-PT'])).toBe('pt-BR');
    expect(idiomaDoAparelho(['es-MX', 'pt-BR'])).toBe('es-419');
    expect(idiomaDoAparelho(['es-419'])).toBe('es-419');
    expect(idiomaDoAparelho(['es_AR'])).toBe('es-419');
  });

  it('cai no português quando nenhuma língua do aparelho é atendida', () => {
    expect(idiomaDoAparelho(['en-US', 'fr-FR'])).toBe('pt-BR');
    expect(idiomaDoAparelho([])).toBe('pt-BR');
  });

  it('a escolha guardada vale mais que o aparelho; escolha inválida é ignorada', () => {
    expect(escolherIdioma('pt-BR', ['es-MX'])).toBe('pt-BR');
    expect(escolherIdioma('es-419', ['pt-BR'])).toBe('es-419');
    expect(escolherIdioma('fr', ['es-MX'])).toBe('es-419');
    expect(escolherIdioma(null, ['es-MX'])).toBe('es-419');
  });
});

describe('tradução', () => {
  it('troca o texto inteiro ao trocar o idioma', () => {
    expect(traduzir('pt-BR', 'mobile.login.entrar')).toBe('Entrar');
    expect(traduzir('es-419', 'mobile.login.entrar')).toBe('Iniciar sesión');
    expect(traduzir('es-419', 'mobile.login.slogan')).toBe('Cada activo listo para operar.');
  });

  it('interpola variáveis', () => {
    expect(traduzir('pt-BR', 'mobile.empresas.entrarRotulo', { nome: 'Alfa' })).toBe('Entrar em Alfa');
    expect(traduzir('es-419', 'mobile.empresas.entrarRotulo', { nome: 'Alfa' })).toBe('Entrar en Alfa');
  });

  it('escolhe a forma de plural do idioma', () => {
    expect(traduzir('pt-BR', 'formulario.resumoTitulo', { count: 1 })).toBe('Corrija 1 campo para continuar');
    expect(traduzir('pt-BR', 'formulario.resumoTitulo', { count: 3 })).toBe('Corrija 3 campos para continuar');
  });

  it('o texto do app diz "segundo fator" e "Minha conta", não "MFA" nem "verificação em duas etapas"', () => {
    const texto = traduzir('pt-BR', 'mobile.mfa.texto');
    expect(texto).toContain('segundo fator');
    expect(texto).toContain('Minha conta');
    expect(texto).not.toMatch(/MFA|duas etapas|Meu perfil/);
    expect(traduzir('es-419', 'mobile.mfa.texto')).toContain('Mi cuenta');
  });

  it('nenhum texto do app mostra número de requisito', () => {
    for (const idioma of ['pt-BR', 'es-419'] as const) {
      for (const k of ['mobile.hoje.vazioTitulo', 'mobile.hoje.vazioTexto'] as const) {
        expect(traduzir(idioma, k)).not.toMatch(/RF-|RNF-/);
      }
    }
  });
});

describe('regra de plural de reserva (sem Intl.PluralRules)', () => {
  it('dá a mesma categoria que o Intl nos dois idiomas', () => {
    for (const idioma of ['pt-BR', 'es-419'] as const) {
      for (const n of [0, 1, 2, 1.5, 0.5, 21, 1_000_000, 2_000_000]) {
        expect(categoriaDeReserva(idioma, n), `${idioma} ${n}`).toBe(new Intl.PluralRules(idioma).select(n));
      }
    }
  });

  it('zero é one em português e other em espanhol', () => {
    expect(categoriaDeReserva('pt-BR', 0)).toBe('one');
    expect(categoriaDeReserva('es-419', 0)).toBe('other');
  });
});
