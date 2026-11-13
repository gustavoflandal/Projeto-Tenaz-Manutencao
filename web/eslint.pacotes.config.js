// Trava de formatação do ADR-013, Parte G (documento 5) no app e nos pacotes. A
// web tem a sua, em eslint.config.js. Aqui só a trava: nenhuma outra regra
// (a análise geral do app e dos pacotes é assunto à parte). Rodado a partir da
// raiz do repositório, com as dependências de lint da web, por `lint:pacotes`.

import { defineConfig } from 'eslint/config';
import tseslint from 'typescript-eslint';

import { regrasDeFormatacao, regrasDoMotorDeFuso } from '../eslint.formatacao.mjs';

export default defineConfig(
  {
    ignores: [
      '**/node_modules/**',
      '**/dist/**',
      '**/.expo/**',
      'web/**',
      'backend/**',
      // Gerado do contrato e dos tokens: não se edita nem se analisa.
      'pacotes/api-tipos/src/esquema.ts',
    ],
  },
  {
    files: ['app/**/*.{ts,tsx}', 'pacotes/**/*.{ts,tsx}'],
    languageOptions: { parser: tseslint.parser },
    // Comentários `eslint-disable` de regras que esta configuração não carrega não são erro.
    linterOptions: { reportUnusedDisableDirectives: 'off' },
    rules: regrasDeFormatacao,
  },
  {
    // O motor de fuso: o único Intl.DateTimeFormat do repositório.
    files: ['pacotes/dominio/src/formatos.ts'],
    rules: regrasDoMotorDeFuso,
  },
);
