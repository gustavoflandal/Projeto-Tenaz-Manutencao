// E2E do Paraguai (ADR-013): conta com RUC, unidade com departamento e código
// postal de 6 dígitos, pessoa com CI, tudo em espanhol; a recusa do servidor
// chega no idioma da resposta (texto exato do catálogo es-419 do servidor).
// Roda contra a API local (veja conta-e-acesso.spec.ts).

import { expect, test, type Page } from '@playwright/test';

import { gerarCI, gerarRUC, tokenDoEmail } from './apoio';

const SUFIXO = `${Date.now()}`.slice(-8);
const EMAIL = `e2e.py.${SUFIXO}@tenaz.test`;
const SENHA = `Senha-E2E-${SUFIXO}-forte`;
const CI = gerarCI();

test.use({ locale: 'es-PY' });
test.describe.configure({ mode: 'serial' });

async function entrar(pagina: Page) {
  await pagina.goto('/entrar');
  await pagina.locator('input[type="email"], input[name="email"]').first().fill(EMAIL);
  await pagina.locator('input[type="password"]').first().fill(SENHA);
  await pagina.locator('form button[type="submit"]').click();
  await expect(pagina).not.toHaveURL(/\/entrar/);
}

test('cria a conta paraguaia: RUC, departamento e código postal de 6 dígitos', async ({ page }) => {
  await page.goto('/criar-conta');
  await page.locator('select[name="pais"]').selectOption('PY');
  await expect(page.locator('select[name="idioma_padrao"]')).toHaveValue('es-419');
  await expect(page.locator('select[name="fuso_horario"]')).toHaveValue('America/Asuncion');

  await page.locator('input[name="razao_social"]').fill(`Transportes E2E ${SUFIXO} SA`);
  await page.getByLabel('RUC').fill(gerarRUC());
  await page.locator('select[name="plano"]').selectOption('essencial');
  await page.locator('input[name="responsavel_nome"]').fill('Responsable E2E');
  await page.locator('input[name="responsavel_email"]').fill(EMAIL);
  await page.locator('input[name="unidade_codigo"]').fill(`P${SUFIXO}`);
  await page.locator('input[name="unidade_nome"]').fill('Depósito Luque');
  // Os turnos primeiro: a mensagem do código postal muda o layout e faria o clique se perder.
  await page.getByRole('button', { name: 'Comercial' }).click();
  await expect(page.getByRole('group', { name: 'Lunes' }).getByText('Sin turno este día.')).toHaveCount(0);

  // O rótulo é "Departamento" e o código postal tem exemplo de 6 dígitos.
  const departamento = page.locator('select[name="unidade_subdivisao"]');
  await expect(departamento).toHaveAccessibleName(/Departamento/);
  await expect(page.locator('select[name="unidade_subdivisao"] option')).toHaveCount(19);
  await departamento.selectOption({ label: 'Central' });
  const postal = page.locator('input[name="unidade_codigo_postal"]');
  await expect(postal).toHaveAccessibleDescription('Ejemplo: 001526');
  await postal.fill('11093');
  await page.getByRole('button', { name: 'Crear cuenta' }).click();
  await expect(postal).toHaveAccessibleDescription(/6 dígitos/);

  await postal.fill('110935');
  // Sai do campo: a mensagem some e o layout se ajusta antes do clique.
  await postal.blur();
  await expect(postal).not.toHaveAccessibleDescription(/6 dígitos, como/);
  await page.getByRole('button', { name: 'Crear cuenta' }).click();
  await expect(page.getByRole('status')).toContainText(`Enviamos la invitación a ${EMAIL}`);
});

test('aceita o convite e entra em espanhol', async ({ page }) => {
  await page.goto(`/convite#token=${await tokenDoEmail(EMAIL)}`);
  for (const campo of await page.locator('input[type="password"]').all()) await campo.fill(SENHA);
  await page.getByRole('button', { name: 'Aceptar invitación' }).click();
  await expect(page.getByRole('status')).toContainText('Invitación aceptada');

  await entrar(page);
  await expect(page.locator('html')).toHaveAttribute('lang', 'es-419');
});

test('unidade: departamento e código postal gravados, e recusa local do código de 5 dígitos', async ({ page }) => {
  await entrar(page);
  await page.goto('/administracao/unidades');
  await expect(page.getByRole('row', { name: /Depósito Luque.*Central/ })).toBeVisible();
  await page.getByRole('link', { name: 'Depósito Luque' }).click();

  const departamento = page.locator('select[name="subdivisao"]');
  const postal = page.locator('input[name="codigo_postal"]');
  await expect(departamento).toHaveValue('PY-11');
  await expect(postal).toHaveValue('110935');

  await postal.fill('00152');
  await page.getByRole('button', { name: 'Guardar' }).click();
  await expect(postal).toHaveAccessibleDescription(/ Use el código postal de 6 dígitos, como 001526\.$/);
});

test('pessoa com CI: máscara e recusa do servidor em espanhol (texto exato do catálogo)', async ({ page }) => {
  await entrar(page);
  await page.goto('/administracao/pessoas/nova');
  await expect(page.locator('select[name="identificacao_tipo"]')).toHaveValue('ci');
  await page.locator('input[name="nome"]').fill(`Persona PY ${SUFIXO}`);
  const campo = page.locator('input[name="identificacao"]');
  await campo.fill(CI);
  await expect(campo).toHaveValue(/^\d(\.\d{3})?(\.\d{3})?$/);
  await page.getByRole('button', { name: 'Registrar persona' }).click();
  await expect(page.getByRole('heading', { level: 1, name: `Persona PY ${SUFIXO}` })).toBeVisible();

  // A mesma CI de novo: o servidor recusa, na língua da resposta (Accept-Language es-419).
  await page.goto('/administracao/pessoas/nova');
  await page.locator('input[name="nome"]').fill(`Duplicada PY ${SUFIXO}`);
  await campo.fill(CI);
  await page.getByRole('button', { name: 'Registrar persona' }).click();
  await expect(campo).toHaveAccessibleDescription(/ Ya existe un registro con este valor\.$/);
});
