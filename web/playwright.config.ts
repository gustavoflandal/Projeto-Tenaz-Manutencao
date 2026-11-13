// E2E de navegador (Playwright) contra a API local e o servidor do Vite.
//
// Pré-requisito: `bash ferramentas/tarefas.sh infra` (PostgreSQL, Mailpit...)
// e `bash ferramentas/tarefas.sh api` (API em :8080). O servidor do Vite sobe
// sozinho (ou reaproveita o que já estiver em :5173). Rodar: `pnpm e2e`.
//
// Os testes criam dados com sufixo único e não limpam o banco local: cada
// execução é uma empresa nova, sem tocar nas existentes.

import { defineConfig, devices } from '@playwright/test';

const PORTA = 5173;

export default defineConfig({
  testDir: './e2e',
  // Um fluxo só, em ordem: criar conta, convite, login, MFA, pessoa, trilha.
  fullyParallel: false,
  workers: 1,
  retries: 0,
  timeout: 120_000,
  expect: { timeout: 15_000 },
  reporter: [['list']],
  use: {
    baseURL: `http://localhost:${PORTA}`,
    locale: 'pt-BR',
    trace: 'retain-on-failure',
    screenshot: 'only-on-failure',
  },
  projects: [{ name: 'chromium', use: { ...devices['Desktop Chrome'] } }],
  webServer: {
    command: 'npx vite --port 5173 --strictPort',
    url: `http://localhost:${PORTA}`,
    reuseExistingServer: true,
    timeout: 60_000,
  },
});
