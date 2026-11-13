// Análise estática da web (seção 8.4 de docs/5_ARQUITETURA_SOFTWARE.md):
// TypeScript estrito, `any` proibido, regras de hooks, acessibilidade em JSX
// (seção 11.5 do documento 7) e nenhum texto de interface fora do pacote de
// internacionalização, mais a trava de formatação do ADR-013, Parte G.

import js from '@eslint/js';
import { defineConfig } from 'eslint/config';
import i18next from 'eslint-plugin-i18next';
import jsxA11y from 'eslint-plugin-jsx-a11y';
import reactHooks from 'eslint-plugin-react-hooks';
import globals from 'globals';
import tseslint from 'typescript-eslint';

import { regrasDeFormatacao } from '../eslint.formatacao.mjs';

export default defineConfig(
  { ignores: ['dist/**', 'node_modules/**', 'test-results/**', 'playwright-report/**'] },
  js.configs.recommended,
  tseslint.configs.strict,
  {
    files: ['**/*.{ts,tsx}'],
    languageOptions: { globals: { ...globals.browser } },
    rules: {
      '@typescript-eslint/no-explicit-any': 'error',
      '@typescript-eslint/consistent-type-imports': ['error', { fixStyle: 'inline-type-imports' }],
      // Formatação para pessoas só por pacotes/dominio/src/formatos.ts (ADR-013, Parte G).
      ...regrasDeFormatacao,
    },
  },
  {
    files: ['src/**/*.tsx'],
    extends: [jsxA11y.flatConfigs.strict, reactHooks.configs.flat.recommended],
  },
  {
    files: ['src/**/*.tsx'],
    ignores: ['src/**/*.test.tsx', 'src/testes/**'],
    plugins: { i18next },
    rules: {
      // Texto visível e nomes acessíveis só pelo i18n; atributos técnicos
      // (variante, align, aria-hidden) ficam de fora.
      'i18next/no-literal-string': [
        'error',
        {
          mode: 'jsx-only',
          // Identificadores técnicos (nome de campo, rota, chave de texto) não são texto de interface.
          words: { exclude: ['[0-9\\s]*', '^[a-z][a-z0-9_.]*$', '^[A-Z_-]+$'] },
          'jsx-attributes': { include: ['alt', 'title', 'placeholder', 'aria-label', 'aria-description', 'aria-valuetext'] },
        },
      ],
    },
  },
  {
    files: ['vite.config.ts', 'eslint.config.js', 'playwright.config.ts', 'e2e/**/*.ts'],
    languageOptions: { globals: { ...globals.node } },
  },
);
