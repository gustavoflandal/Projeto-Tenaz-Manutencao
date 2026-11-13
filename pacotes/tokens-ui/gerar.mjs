// Gera tokens.css e tailwind.css a partir de tokens.json.
//
// O documento 7 (seção 3.3) pede um único arquivo de tokens para que o bloco
// do tema escuro automático (prefers-color-scheme) e o do tema escuro
// escolhido à mão nunca divirjam: CSS não deixa um seletor herdar de outro.
//
// Uso: node gerar.mjs           escreve os dois arquivos
//      node gerar.mjs --conferir falha se os arquivos versionados divergirem

import { readFileSync, writeFileSync } from 'node:fs';
import { dirname, resolve } from 'node:path';
import { fileURLToPath } from 'node:url';

const AQUI = dirname(fileURLToPath(import.meta.url));

const CABECALHO =
  '/* Gerado por pacotes/tokens-ui/gerar.mjs a partir de tokens.json. Não edite à mão. */\n';

/** Chaves do tema que não são cor e por isso não viram utilitário de cor. */
const NAO_COR = (chave) => chave === 'color-scheme' || chave.startsWith('sombra-');

function declaracoesDoTema(tema, recuo) {
  return Object.entries(tema)
    .map(([chave, valor]) =>
      chave === 'color-scheme' ? `${recuo}color-scheme: ${valor};` : `${recuo}--tz-${chave}: ${valor};`,
    )
    .join('\n');
}

/**
 * @param {any} t conteúdo de tokens.json (JSON lido do arquivo, sem tipo estático; a forma é conferida por gerar.test.ts)
 * @returns {{ tokensCss: string, tailwindCss: string }}
 */
export function gerar(t) {
  const linhas = [];
  linhas.push(CABECALHO);
  linhas.push(':root {');
  linhas.push('  /* Primitivos: não variam com o tema e não são usados direto nos componentes */');
  for (const [chave, valor] of Object.entries(t.primitivos)) linhas.push(`  --tz-${chave}: ${valor};`);
  linhas.push('');
  linhas.push('  /* Variantes de status sobre superfícies Aço (toast, menu) */');
  for (const [chave, valor] of Object.entries(t.invariantes)) linhas.push(`  --tz-${chave}: ${valor};`);
  linhas.push('');
  linhas.push('  /* Tipografia */');
  for (const [chave, valor] of Object.entries(t.fontes)) linhas.push(`  --tz-fonte-${chave}: ${valor};`);
  linhas.push(`  --tz-largura-titulo: ${t['largura-titulo']};`);
  linhas.push(`  --tz-titulo-espacamento-letras: ${t['titulo-espacamento-letras']};`);
  for (const [nome, def] of Object.entries(t.tipografia)) {
    linhas.push(`  --tz-texto-${nome}: ${def.valor} var(--tz-fonte-${def.fonte});`);
  }
  linhas.push('');
  linhas.push('  /* Espaçamento, base 4 px */');
  for (const [chave, valor] of Object.entries(t.espaco)) linhas.push(`  --tz-espaco-${chave}: ${valor};`);
  linhas.push('');
  linhas.push('  /* Alturas de controle */');
  for (const [chave, valor] of Object.entries(t.controle)) linhas.push(`  --tz-controle-${chave}: ${valor};`);
  linhas.push('');
  linhas.push('  /* Raios */');
  for (const [chave, valor] of Object.entries(t.raio)) linhas.push(`  --tz-raio-${chave}: ${valor};`);
  linhas.push('');
  linhas.push('  /* Camadas */');
  for (const [chave, valor] of Object.entries(t.camada)) linhas.push(`  --tz-z-${chave}: ${valor};`);
  linhas.push('');
  linhas.push('  /* Movimento */');
  for (const [chave, valor] of Object.entries(t.movimento)) linhas.push(`  --tz-${chave}: ${valor};`);
  linhas.push('');
  linhas.push('  /* Foco */');
  for (const [chave, valor] of Object.entries(t.foco)) linhas.push(`  --tz-foco-${chave}: ${valor};`);
  linhas.push('');
  linhas.push('  /* Semânticos: tema claro (padrão) */');
  linhas.push(declaracoesDoTema(t.temas.claro, '  '));
  linhas.push('}');
  linhas.push('');
  linhas.push('/* Tema escuro automático, salvo se o usuário escolheu claro ou sol */');
  linhas.push('@media (prefers-color-scheme: dark) {');
  linhas.push('  :root:not([data-tema="claro"]):not([data-tema="sol"]) {');
  linhas.push(declaracoesDoTema(t.temas.escuro, '    '));
  linhas.push('  }');
  linhas.push('}');
  linhas.push('');
  linhas.push(':root[data-tema="escuro"] {');
  linhas.push(declaracoesDoTema(t.temas.escuro, '  '));
  linhas.push('}');
  linhas.push('');
  linhas.push('/* Tema sol: alto contraste para tablets em pátio e oficina */');
  linhas.push(':root[data-tema="sol"] {');
  linhas.push(declaracoesDoTema(t.temas.sol, '  '));
  linhas.push('}');
  linhas.push('');
  linhas.push('/* Movimento reduzido: sem deslocamento, sem brilho de esqueleto, fades curtos */');
  linhas.push('@media (prefers-reduced-motion: reduce) {');
  linhas.push('  :root {');
  linhas.push('    --tz-duracao-base: 0ms;');
  linhas.push('    --tz-duracao-media: 0ms;');
  linhas.push('    --tz-duracao-lenta: 0ms;');
  linhas.push('  }');
  linhas.push('  *, *::before, *::after {');
  linhas.push('    animation-duration: 1ms !important;');
  linhas.push('    animation-iteration-count: 1 !important;');
  linhas.push('    scroll-behavior: auto !important;');
  linhas.push('  }');
  linhas.push('}');
  const tokensCss = `${linhas.join('\n')}\n`;

  // Tema do Tailwind: zera as escalas padrão (nenhuma cor, tamanho ou sombra
  // fora dos tokens fica disponível como classe) e aponta cada utilitário
  // para a variável do token.
  const tw = [];
  tw.push(CABECALHO);
  tw.push('@theme inline {');
  tw.push('  --color-*: initial;');
  const cores = Object.keys(t.temas.claro).filter((c) => !NAO_COR(c));
  for (const chave of [...cores, ...Object.keys(t.invariantes)]) tw.push(`  --color-${chave}: var(--tz-${chave});`);
  tw.push('');
  tw.push('  /* Espaçamento: múltiplos da base de 4 px (espaco-1) */');
  tw.push('  --spacing: var(--tz-espaco-1);');
  tw.push('');
  tw.push('  --radius-*: initial;');
  for (const chave of Object.keys(t.raio)) tw.push(`  --radius-${chave}: var(--tz-raio-${chave});`);
  tw.push('');
  tw.push('  --shadow-*: initial;');
  tw.push('  --inset-shadow-*: initial;');
  tw.push('  --drop-shadow-*: initial;');
  tw.push('  --text-shadow-*: initial;');
  for (const chave of Object.keys(t.temas.claro).filter((c) => c.startsWith('sombra-'))) {
    tw.push(`  --shadow-${chave.replace('sombra-', '')}: var(--tz-${chave});`);
  }
  tw.push('');
  tw.push('  --font-*: initial;');
  for (const chave of Object.keys(t.fontes)) tw.push(`  --font-${chave}: var(--tz-fonte-${chave});`);
  tw.push('');
  tw.push('  /* Tamanho de letra só pelos utilitários texto-* abaixo */');
  tw.push('  --text-*: initial;');
  tw.push('  --leading-*: initial;');
  tw.push('  --tracking-*: initial;');
  tw.push('  --container-*: initial;');
  tw.push('');
  tw.push('  --breakpoint-*: initial;');
  for (const [chave, valor] of Object.entries(t.breakpoints)) tw.push(`  --breakpoint-${chave}: ${valor};`);
  tw.push('');
  tw.push('  --ease-*: initial;');
  for (const chave of Object.keys(t.movimento).filter((c) => c.startsWith('curva-'))) {
    tw.push(`  --ease-${chave.replace('curva-', '')}: var(--tz-${chave});`);
  }
  tw.push('}');
  tw.push('');
  for (const [nome, def] of Object.entries(t.tipografia)) {
    tw.push(`@utility texto-${nome} {`);
    tw.push(`  font: var(--tz-texto-${nome});`);
    if (def.fonte === 'titulo') {
      tw.push('  font-stretch: var(--tz-largura-titulo);');
      tw.push('  letter-spacing: var(--tz-titulo-espacamento-letras);');
      tw.push('  text-wrap: balance;');
    }
    if (def.fonte === 'dado') tw.push('  font-variant-numeric: tabular-nums;');
    if (def['caixa-alta']) tw.push('  text-transform: uppercase;');
    if (def['espacamento-letras']) tw.push(`  letter-spacing: ${def['espacamento-letras']};`);
    tw.push('}');
    tw.push('');
  }
  const tailwindCss = `${tw.join('\n').trimEnd()}\n`;

  return { tokensCss, tailwindCss };
}

export function lerTokens() {
  return JSON.parse(readFileSync(resolve(AQUI, 'tokens.json'), 'utf8'));
}

export const ARQUIVOS = {
  tokensCss: resolve(AQUI, 'tokens.css'),
  tailwindCss: resolve(AQUI, 'tailwind.css'),
};

const executadoDireto = process.argv[1] && resolve(process.argv[1]) === fileURLToPath(import.meta.url);
if (executadoDireto) {
  const gerado = gerar(lerTokens());
  if (process.argv.includes('--conferir')) {
    let divergente = false;
    for (const [nome, caminho] of Object.entries(ARQUIVOS)) {
      if (readFileSync(caminho, 'utf8').replace(/\r\n/g, '\n') !== gerado[nome]) {
        console.error(`${caminho} difere do gerado a partir de tokens.json`);
        divergente = true;
      }
    }
    if (divergente) {
      console.error('rode: npx pnpm --filter @tenaz/tokens-ui gerar');
      process.exit(1);
    }
    console.log('tokens em dia');
  } else {
    for (const [nome, caminho] of Object.entries(ARQUIVOS)) writeFileSync(caminho, gerado[nome], 'utf8');
    console.log('tokens.css e tailwind.css gerados');
  }
}
