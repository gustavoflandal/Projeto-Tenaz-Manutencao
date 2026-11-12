import { defineConfig } from 'vitest/config';

// O motor de sincronização é testado com SQLite real no Node, como manda a
// seção 9.1 de docs/5_ARQUITETURA_SOFTWARE.md. node:sqlite ainda é
// experimental no Node 22 e precisa da permissão explícita.
export default defineConfig({
  test: {
    include: ['src/**/*.test.ts'],
    environment: 'node',
    pool: 'forks',
    poolOptions: {
      forks: { execArgv: ['--experimental-sqlite'] },
    },
  },
});
