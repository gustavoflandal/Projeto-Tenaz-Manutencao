// E2E multipaís (ADR-013): empresa argentina criada num navegador em espanhol,
// convite, login, troca de idioma gravada no usuário e pessoa com DNI, com a
// recusa do servidor no idioma. Roda contra a API local (veja conta-e-acesso.spec.ts).

import { expect, test, type Page } from '@playwright/test';

import { gerarCUIT, gerarDNI, tokenDoEmail } from './apoio';

const SUFIXO = `${Date.now()}`.slice(-8);
const EMAIL = `e2e.ar.${SUFIXO}@tenaz.test`;
const SENHA = `Senha-E2E-${SUFIXO}-forte`;
const DNI = gerarDNI();

test.use({ locale: 'es-AR' });
test.describe.configure({ mode: 'serial' });

/** Entra pela tela de login, que está no idioma do navegador (espanhol) ou no que a pessoa guardou. */
async function entrar(pagina: Page) {
  await pagina.goto('/entrar');
  await pagina.locator('input[type="email"], input[name="email"]').first().fill(EMAIL);
  await pagina.locator('input[type="password"]').first().fill(SENHA);
  await pagina.locator('form button[type="submit"]').click();
  await expect(pagina).not.toHaveURL(/\/entrar/);
}

test('cria a conta argentina em espanhol: o país define CUIT, idioma e fuso', async ({ page }) => {
  await page.goto('/criar-conta');
  await expect(page.getByRole('heading', { level: 1, name: 'Crear cuenta de la empresa' })).toBeVisible();
  await expect(page.locator('select[name="pais"]')).toHaveValue('AR');
  await expect(page.locator('select[name="idioma_padrao"]')).toHaveValue('es-419');
  await expect(page.locator('select[name="fuso_horario"]')).toHaveValue('America/Argentina/Buenos_Aires');

  await page.locator('input[name="razao_social"]').fill(`Frigorífico E2E ${SUFIXO} SA`);
  await page.getByLabel('CUIT').fill(gerarCUIT());
  await page.locator('select[name="plano"]').selectOption('essencial');
  await page.locator('input[name="responsavel_nome"]').fill('Responsable E2E');
  await page.locator('input[name="responsavel_email"]').fill(EMAIL);
  await page.locator('input[name="unidade_codigo"]').fill(`A${SUFIXO}`);
  await page.locator('input[name="unidade_nome"]').fill('Planta Buenos Aires');
  // Argentina: a subdivisão é a província (nome oficial) e o CPA vai normalizado.
  await page.locator('select[name="unidade_subdivisao"]').selectOption({ label: 'Ciudad Autónoma de Buenos Aires' });
  await page.locator('input[name="unidade_codigo_postal"]').fill('c1425 dkf');
  await page.getByRole('button', { name: /Comercial/ }).click();
  await page.getByRole('button', { name: 'Crear cuenta' }).click();

  await expect(page.getByRole('status')).toContainText(`Enviamos la invitación a ${EMAIL}`);
});

test('aceita o convite e entra em espanhol (idioma padrão da empresa)', async ({ page }) => {
  await page.goto(`/convite#token=${await tokenDoEmail(EMAIL)}`);
  for (const campo of await page.locator('input[type="password"]').all()) await campo.fill(SENHA);
  await page.getByRole('button', { name: 'Aceptar invitación' }).click();
  await expect(page.getByRole('status')).toContainText('Invitación aceptada');

  await entrar(page);
  await expect(page.getByRole('heading', { level: 1, name: /Panel/ })).toBeVisible();
  await expect(page.locator('html')).toHaveAttribute('lang', 'es-419');
});

test('troca o idioma e a escolha fica no usuário, não só no navegador', async ({ page }) => {
  await entrar(page);

  await page.goto('/conta');
  await page.getByRole('radio', { name: 'Português (Brasil)' }).click();
  await expect(page.getByRole('heading', { level: 1, name: 'Minha conta' })).toBeVisible();
  await expect(page.locator('html')).toHaveAttribute('lang', 'pt-BR');

  // Sem o que o navegador guardou, vale o idioma gravado no usuário.
  await page.evaluate(() => globalThis.localStorage.clear());
  await page.reload();
  await expect(page.getByRole('heading', { level: 1, name: 'Minha conta' })).toBeVisible();
});

test('empresa: identificação CUIT com máscara, país e moeda só leitura', async ({ page }) => {
  await entrar(page);
  await page.goto('/administracao/empresa');
  await expect(page.locator('input[name="identificacao"]')).toHaveValue(/^\d\d-\d{8}-\d$/);
  await expect(page.getByLabel('CUIT')).toBeVisible();
  await expect(page.getByLabel('País')).toHaveAttribute('readonly', '');
  await expect(page.getByLabel('Moeda')).toHaveValue('Peso argentino (ARS)');
  await expect(page.locator('select[name="idioma_padrao"]')).toHaveValue('es-419');
});

test('unidade argentina: província, CPA normalizado e recusa de CPA de outra província', async ({ page }) => {
  await entrar(page);
  await page.goto('/administracao/unidades');
  await page.getByRole('link', { name: 'Planta Buenos Aires' }).click();
  const provincia = page.locator('select[name="subdivisao"]');
  const cpa = page.locator('input[name="codigo_postal"]');
  await expect(provincia).toHaveValue('AR-C');
  await expect(cpa).toHaveValue('C1425DKF');

  // Muda para uma província cuja letra não é a do CPA: a recusa é local, sem enviar.
  await provincia.selectOption('AR-X');
  await page.getByRole('button', { name: 'Salvar' }).click();
  await expect(cpa).toHaveAccessibleDescription(/letra inicial do CPA/);

  await cpa.fill('X5000ABC');
  // Sai do campo: a mensagem some e o layout se ajusta antes do clique (senão o botão se move sob o mouse).
  await cpa.blur();
  await expect(cpa).not.toHaveAccessibleDescription(/letra inicial do CPA/);
  await page.getByRole('button', { name: 'Salvar' }).click();
  await expect(page.getByRole('status').filter({ hasText: 'Unidade salva' })).toBeVisible();
});

test('pessoa com DNI: máscara, recusa local e recusa do servidor em português', async ({ page }) => {
  await entrar(page); // o idioma gravado no usuário (português) vale

  await page.goto('/administracao/pessoas/nova');
  await expect(page.locator('select[name="identificacao_tipo"]')).toHaveValue('dni');
  await page.locator('input[name="nome"]').fill(`Pessoa AR ${SUFIXO}`);

  const campo = page.locator('input[name="identificacao"]');
  await campo.fill('123456');
  await page.getByRole('button', { name: 'Cadastrar pessoa' }).click();
  await expect(campo).toHaveAccessibleDescription(/DNI inválido/);

  await campo.fill(DNI);
  await expect(campo).toHaveValue(/^\d{1,2}\.\d{3}\.\d{3}$/);
  await page.getByRole('button', { name: 'Cadastrar pessoa' }).click();
  await expect(page.getByRole('heading', { level: 1, name: `Pessoa AR ${SUFIXO}` })).toBeVisible();

  await page.goto('/administracao/pessoas');
  await expect(page.getByRole('row', { name: new RegExp(`Pessoa AR ${SUFIXO}`) })).toContainText('DNI *****');

  // O mesmo DNI de novo: o servidor recusa e a mensagem vem no idioma da interface (Accept-Language).
  await page.goto('/administracao/pessoas/nova');
  await page.locator('input[name="nome"]').fill(`Duplicada ${SUFIXO}`);
  await campo.fill(DNI);
  await page.getByRole('button', { name: 'Cadastrar pessoa' }).click();
  // O servidor recusa a identificação repetida; o texto exato do catálogo pt-BR vem no campo
  // (aqui o usuário já está em português; a prova em espanhol está em paraguai.spec.ts).
  await expect(campo).toHaveAccessibleDescription(/ Já existe um registro com este valor\.$/);
});
