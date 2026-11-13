// Fluxo crítico da E01 no navegador, contra a API local:
// criar conta -> convite pelo Mailpit -> login -> MFA (TOTP calculado do
// segredo exibido) -> cadastro de pessoa -> trilha com o evento -> integridade.

import { expect, test, type Page } from '@playwright/test';

import { esperarProximoPasso, gerarCNPJ, gerarCPF, tokenDoEmail, totp } from './apoio';

const SUFIXO = `${Date.now()}`.slice(-8);
const EMAIL = `e2e.${SUFIXO}@tenaz.test`;
const SENHA = `Senha-E2E-${SUFIXO}-forte`;
const NOME_DA_PESSOA = `Pessoa E2E ${SUFIXO}`;

let segredoDoMFA = '';

async function entrar(pagina: Page) {
  await pagina.goto('/entrar');
  await pagina.getByLabel('E-mail').fill(EMAIL);
  await pagina.getByLabel('Senha', { exact: true }).fill(SENHA);
  await pagina.getByRole('button', { name: 'Entrar', exact: true }).click();
}

test.describe.configure({ mode: 'serial' });

test('cria a conta com os turnos do modelo comercial', async ({ page }) => {
  await page.goto('/criar-conta');
  await page.getByLabel('Razão social').fill(`Empresa E2E ${SUFIXO} Ltda`);
  await page.getByLabel('CNPJ').fill(gerarCNPJ());
  await page.getByLabel('Plano', { exact: true }).selectOption('essencial');
  await page.getByLabel('Nome do responsável').fill('Responsável E2E');
  await page.getByLabel('E-mail do responsável').fill(EMAIL);
  await page.getByLabel('Código', { exact: true }).fill(`E${SUFIXO}`.slice(0, 20));
  await page.getByLabel('Nome', { exact: true }).fill('Matriz E2E');
  await page.getByLabel('Estado', { exact: false }).selectOption({ label: 'Minas Gerais' });
  await page.getByLabel('CEP', { exact: false }).fill('30130-100');
  await page.getByRole('button', { name: 'Comercial' }).click();
  await expect(page.getByRole('group', { name: 'Segunda-feira' }).getByLabel('Início')).toHaveValue('08:00');

  await page.getByRole('button', { name: 'Criar conta' }).click();

  await expect(page.getByRole('status')).toContainText(`Enviamos o convite para ${EMAIL}`);
});

test('aceita o convite pelo link do e-mail (token no fragmento)', async ({ page }) => {
  const token = await tokenDoEmail(EMAIL);
  await page.goto(`/convite#token=${token}`);
  // O token sai da barra de endereço depois de lido.
  await expect(page).toHaveURL(/\/convite$/);

  await page.getByLabel('Crie uma senha').fill(SENHA);
  await page.getByLabel('Repita a senha nova').fill(SENHA);
  await page.getByRole('button', { name: 'Aceitar convite' }).click();

  await expect(page.getByRole('status')).toContainText('Convite aceito');
});

test('entra, ativa o segundo fator e entra de novo com o código', async ({ page }) => {
  await entrar(page);
  await expect(page.getByRole('heading', { level: 1, name: 'Painel' })).toBeVisible();
  // O selo de versão está no canto inferior direito.
  await expect(page.getByRole('button', { name: /^Versão / })).toBeVisible();

  await page.goto('/conta');
  await page.getByRole('button', { name: 'Ativar segundo fator' }).click();
  await expect(page.getByRole('img', { name: 'QR code para o aplicativo autenticador' })).toBeVisible();
  segredoDoMFA = (await page.getByText(/^[A-Z2-7]{16,}$/).first().innerText()).trim();
  await page.getByLabel('Código de 6 dígitos').fill(totp(segredoDoMFA));
  await page.getByRole('button', { name: 'Confirmar e ativar' }).click();
  await expect(page.getByText('Segundo fator ativado.')).toBeVisible();

  await page.getByRole('button', { name: /^Menu de / }).click();
  await page.getByRole('menuitem', { name: 'Sair' }).click();
  await expect(page.getByRole('heading', { level: 1, name: 'Entrar' })).toBeVisible();

  await entrar(page);
  await expect(page.getByRole('heading', { level: 1, name: 'Confirme o segundo fator' })).toBeVisible();
  // O servidor recusa um passo de TOTP já usado: espera o próximo.
  await esperarProximoPasso();
  await page.getByLabel('Código', { exact: true }).fill(totp(segredoDoMFA));
  await page.getByRole('button', { name: 'Verificar código' }).click();
  await expect(page.getByRole('heading', { level: 1, name: 'Painel' })).toBeVisible();
});

test('cadastra uma pessoa com o aviso de dado de saúde', async ({ page }) => {
  await entrar(page);
  await esperarProximoPasso();
  await page.getByLabel('Código', { exact: true }).fill(totp(segredoDoMFA));
  await page.getByRole('button', { name: 'Verificar código' }).click();
  await expect(page.getByRole('heading', { level: 1, name: 'Painel' })).toBeVisible();

  await page.goto('/administracao/pessoas/nova');
  await expect(page.getByText('A Tenaz não guarda dado de saúde.')).toBeVisible();
  await page.getByLabel('Nome', { exact: true }).fill(NOME_DA_PESSOA);
  await page.getByLabel('CPF (opcional)').fill(gerarCPF());
  await page.getByRole('button', { name: 'Cadastrar pessoa' }).click();

  await expect(page.getByRole('heading', { level: 1, name: NOME_DA_PESSOA })).toBeVisible();
  await expect(page.getByText('A Tenaz não guarda dado de saúde.')).toBeVisible();

  // A lista mostra o CPF mascarado, nunca o completo.
  await page.goto('/administracao/pessoas');
  const linha = page.getByRole('row', { name: new RegExp(NOME_DA_PESSOA) });
  await expect(linha).toContainText('***.');

  // Trilha: a inclusão da pessoa, com a tela de origem do catálogo e o hash no detalhe.
  await page.goto('/administracao/auditoria?entidade=pessoa');
  // A chave da entidade não traz mais o nome (dado pessoal fora da trilha): acha a inclusão de pessoa mais recente.
  const evento = page.getByRole('row', { name: /Inclusão/ }).filter({ hasText: 'adm.pessoas' }).first();
  await expect(evento).toBeVisible();
  await expect(evento).toContainText('adm.pessoas');
  await evento.getByRole('button', { name: /Ver detalhes/ }).click();
  const detalhe = page.getByRole('dialog', { name: 'Detalhe do evento' });
  await expect(detalhe).toContainText(/[0-9a-f]{64}/);
  await page.keyboard.press('Escape');

  // Endereço da unidade inicial: subdivisão com o nome oficial e CEP normalizado.
  await page.goto('/administracao/unidades');
  await expect(page.getByRole('row', { name: /Matriz E2E.*Minas Gerais/ })).toBeVisible();
  await page.getByRole('link', { name: 'Matriz E2E' }).click();
  await expect(page.locator('select[name="subdivisao"]')).toHaveValue('BR-MG');
  await expect(page.locator('input[name="codigo_postal"]')).toHaveValue('30130100');

  // Integridade: o dia de hoje ainda não está selado; nenhum dia pode aparecer violado.
  await page.goto('/administracao/auditoria/integridade');
  await page.getByRole('button', { name: 'Verificar integridade' }).click();
  const indisponivel = page.getByText('Sem a cópia imutável não é possível afirmar a integridade');
  const tabela = page.getByRole('table', { name: 'Integridade da trilha' });
  await expect(indisponivel.or(tabela)).toBeVisible();
  await expect(page.getByText('Violado', { exact: true })).toHaveCount(0);

  // Capturas: a tela de configuração lista as telas do catálogo.
  await page.goto('/administracao/auditoria/capturas');
  await expect(page.getByRole('heading', { name: 'Módulo adm' })).toBeVisible();
  await expect(page.getByRole('switch', { name: 'Visualização em Pessoas' })).toBeVisible();
});
