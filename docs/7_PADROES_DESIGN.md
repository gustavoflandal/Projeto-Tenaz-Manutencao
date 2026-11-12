# Padrões de Design — Tenaz

**Versão:** 0.3 · **Data:** 12/11/2026 · **Status:** Em elaboração
**Produto:** Tenaz · Gestão de Manutenção e Frota
Alinhado ao sumário executivo 0.4
**Base:** identidade visual aprovada em [marca/README.md](../marca/README.md) e [marca/brandbook.html](../marca/brandbook.html); vocabulário e estados em [0_SUMARIO_EXECUTIVO.md](0_SUMARIO_EXECUTIVO.md)
**Plataformas:** web (React + TypeScript) e app de campo (React Native + Expo, offline)

Este documento traduz a marca em regras de interface. Define os tokens de design, os componentes, as telas principais (inclusive as da conta da subcontratada, módulo SRV, e da conexão entre contas), o app de campo, o tom das mensagens e os requisitos de acessibilidade. Quando a marca e este documento divergirem, vale a marca para cores, fontes e símbolo; vale o sumário executivo para nomes de entidades e estados. As divergências encontradas estão na seção 13.

---

## Índice

1. [Princípios de design](#1-princípios-de-design)
2. [A marca na interface](#2-a-marca-na-interface)
3. [Design tokens](#3-design-tokens)
4. [Contraste WCAG 2.1](#4-contraste-wcag-21)
5. [Layout do web](#5-layout-do-web)
6. [Componentes](#6-componentes)
7. [Status, prioridade e criticidade](#7-status-prioridade-e-criticidade)
8. [Telas principais](#8-telas-principais)
9. [App de campo](#9-app-de-campo)
10. [Microcopy e tom de voz](#10-microcopy-e-tom-de-voz)
11. [Acessibilidade](#11-acessibilidade)
12. [Iconografia](#12-iconografia)
13. [Inconsistências e pendências](#13-inconsistências-e-pendências)

---

## 1. Princípios de design

| # | Princípio | O que significa na prática |
|---|---|---|
| P1 | **Usável sem treinamento formal** | Cada tela tem uma ação principal evidente. Rótulos usam o vocabulário do sumário executivo, sem siglas internas. Primeiros passos guiados no primeiro acesso. Nenhuma função essencial fica só em atalho de teclado, clique direito ou gesto escondido. |
| P2 | **Feito para luva, sol e pressa** | No app, alvos de toque de no mínimo 48 dp, ações principais na metade inferior da tela, texto de no mínimo 16 dp, tema de alto contraste para uso a céu aberto. Nada depende de precisão fina, gesto de pinça ou toque longo. |
| P3 | **Denso para quem planeja** | No web, o planejador vê muita informação sem rolar: tabelas com densidade ajustável, colunas fixas, filtros salvos, ações em lote e atalhos. Densidade alta não significa letra pequena: o mínimo continua 13 px. |
| P4 | **Status sempre com ícone e texto** | Nenhuma informação depende só de cor. Todo status tem cor, ícone e palavra. Isso vale para daltonismo, impressão em preto e branco e leitura ao sol. |
| P5 | **Número sempre com unidade** | 248.315 km, 12.480 h, R$ 3.418,90, 97,4 %, 3 dias. Números em JetBrains Mono com algarismos tabulares, alinhados à direita em tabelas. Nunca um número solto. |
| P6 | **Diz o que fazer e quando** | Alerta é uma frase com prazo ou ação. Erro explica o que aconteceu e o próximo passo. Sem exclamação, sem caixa alta, sem emoji. |
| P7 | **Offline é o normal no campo** | O app nunca bloqueia o trabalho por falta de sinal. O usuário sempre sabe o que já foi enviado e o que está só no aparelho. |
| P8 | **O laranja Sinal é raro** | Sinal marca a ação principal de cada tela e o ponto do símbolo. Nunca indica status, alerta, seleção ou gráfico. Se tudo é laranja, nada é. |

---

## 2. A marca na interface

### 2.1 Papéis das cores da marca

| Cor | HEX | Papel na interface | Não usar para |
|---|---|---|---|
| Aço | `#14202E` | Texto principal, menu lateral, títulos, foco, cabeçalho de tabela em destaque, toasts | Fundo de áreas grandes de conteúdo no tema claro |
| Sinal | `#FF6A13` | Fundo do botão primário (uma ação principal por tela), ponto do símbolo | Status, alerta, texto, link, ícone, borda de foco, seleção, gráfico, badge |
| Névoa | `#F2F4F6` | Fundo das telas, linha em hover | Texto sobre fundo claro |
| Concreto | `#9AA6B2` | Bordas decorativas, ícones decorativos, texto desabilitado | Texto legível sobre fundo claro (contraste 2,48:1 em branco) |

Proporção de uso de referência da marca: Névoa 60 %, Aço 28 %, Concreto 8 %, Sinal 4 %. Numa tela típica, Sinal aparece em no máximo um botão e no ponto do logo.

### 2.2 Regras do Sinal

1. Uma ação primária em Sinal por tela ou por modal. Ações primárias de cartões dentro de listas usam o botão secundário.
2. Texto e ícone sobre Sinal são sempre Aço `#14202E` (5,74:1). Texto branco sobre Sinal é proibido (2,87:1, reprova inclusive para texto grande).
3. Sinal nunca é cor de status. O âmbar de Atenção `#C98C00` existe para isso, e os dois não podem aparecer como sinônimos.
4. Sinal não aparece como texto nem como contorno sobre fundo claro (2,87:1 em branco, 2,60:1 em Névoa).
5. No tema escuro o botão primário continua Sinal com texto Aço.

### 2.3 Logo no produto

- Menu lateral expandido: logo horizontal negativo, altura 24 px, a partir de 112 px de largura. Menu recolhido: só o símbolo, 32 px.
- Favicon: `marca/logo/svg/favicon.svg`, que acompanha o tema.
- Ícone do app: `marca/logo/png/` (1024, 512, 192 e 180 px).
- Tela de carregamento do app: símbolo positivo de 96 dp centralizado sobre Névoa, sem animação de giro (a marca não permite girar o símbolo).

---

## 3. Design tokens

Os tokens têm dois níveis:

- **Primitivos:** os valores da marca (`--tz-aco-900`, `--tz-sinal`). Não variam com o tema e não são usados direto nos componentes.
- **Semânticos:** o papel na interface (`--tz-texto`, `--tz-superficie`, `--tz-status-parado`). Variam com o tema. Componentes só usam tokens semânticos.

Prefixo `tz-`, nomes em português sem acento, em kebab-case. Os temas são `claro` (padrão no web), `escuro` e `sol` (alto contraste, padrão sugerido no app quando houver luz ambiente forte).

### 3.1 Cores primitivas

| Token | HEX | Origem |
|---|---|---|
| `aco-900` | `#14202E` | Marca (Aço) |
| `aco-700` | `#2B3A4C` | Marca |
| `aco-500` | `#4A5868` | Marca |
| `aco-400` | `#5F6C7B` | **Extensão proposta** (ver 13.6) |
| `aco-300` | `#9AA6B2` | Marca (Concreto) |
| `aco-100` | `#D9DFE6` | Marca |
| `aco-50` | `#F2F4F6` | Marca (Névoa) |
| `sinal` | `#FF6A13` | Marca |
| `sinal-hover` | `#FF7A2E` | Derivado: Sinal mais claro, mantém texto Aço em 6,33:1 |
| `sinal-pressionado` | `#F25F0A` | Derivado: Sinal mais escuro, texto Aço em 5,03:1 |
| `branco` | `#FFFFFF` | Superfície |
| `escuro-fundo` | `#0E151E` | Brandbook, tema escuro |
| `escuro-superficie` | `#16202C` | Brandbook, tema escuro |
| `escuro-elevada` | `#1F2B3A` | Derivado, tema escuro |
| `escuro-faixa` | `#182330` | Brandbook, tema escuro |
| `escuro-linha` | `#263241` | Brandbook, tema escuro |
| `cinza-controle` | `#8592A0` | Brandbook (texto terciário escuro), usado como borda de campo no claro |
| `cinza-controle-escuro` | `#6B7887` | Brandbook (texto terciário claro), usado como borda de campo no escuro |

Cores de status (marca, semânticas):

| Status (marca) | Claro | Fundo claro | Escuro | Fundo escuro |
|---|---|---|---|---|
| Operando | `#1E9E6A` | `#E3F4EC` | `#3CC28B` | `#12301F` |
| Preventiva | `#2F6FDB` | `#E4ECFB` | `#6E9BF0` | `#17264A` |
| Atenção | `#C98C00` | `#FBF1D9` | `#F0B83A` | `#3A2D0C` |
| Parado | `#D33A2C` | `#FBE5E2` | `#F06A5C` | `#3D1714` |
| Inativo | `#7E8A97` | `#ECEFF2` | `#97A3B0` | `#222C37` |

Os fundos de status vêm do brandbook. Os derivados de perigo: `perigo-hover` claro `#B8301F` (texto branco 6,03:1) e escuro `#F48577` (texto Aço 6,64:1).

### 3.2 Cores semânticas por tema

| Token semântico | Claro | Escuro | Sol |
|---|---|---|---|
| `fundo` | `#F2F4F6` | `#0E151E` | `#FFFFFF` |
| `superficie` | `#FFFFFF` | `#16202C` | `#FFFFFF` |
| `superficie-elevada` (menus, modais) | `#FFFFFF` | `#1F2B3A` | `#FFFFFF` |
| `superficie-hover` | `#F2F4F6` | `#1F2B3A` | `#F2F4F6` |
| `superficie-selecionada` | `#D9DFE6` | `#2B3A4C` | `#D9DFE6` |
| `superficie-esqueleto` | `#D9DFE6` | `#1F2B3A` | `#D9DFE6` |
| `texto` | `#14202E` | `#E8EDF2` | `#14202E` |
| `texto-secundario` | `#4A5868` | `#A9B5C2` | `#2B3A4C` |
| `texto-terciario` | `#5F6C7B` | `#8592A0` | `#4A5868` |
| `texto-desabilitado` | `#9AA6B2` | `#4A5868` | `#9AA6B2` |
| `texto-inverso` | `#F2F4F6` | `#14202E` | `#F2F4F6` |
| `borda` (decorativa) | `#D9DFE6` | `#263241` | `#9AA6B2` |
| `borda-controle` (campos) | `#8592A0` | `#6B7887` | `#4A5868` |
| `borda-controle-hover` | `#4A5868` | `#A9B5C2` | `#14202E` |
| `foco` | `#14202E` | `#E8EDF2` | `#14202E` |
| `acao-primaria` | `#FF6A13` | `#FF6A13` | `#FF6A13` |
| `acao-primaria-hover` | `#FF7A2E` | `#FF7A2E` | `#FF7A2E` |
| `acao-primaria-pressionada` | `#F25F0A` | `#F25F0A` | `#F25F0A` |
| `acao-primaria-texto` | `#14202E` | `#14202E` | `#14202E` |
| `perigo` | `#D33A2C` | `#F06A5C` | `#D33A2C` |
| `perigo-hover` | `#B8301F` | `#F48577` | `#B8301F` |
| `perigo-texto` | `#FFFFFF` | `#14202E` | `#FFFFFF` |
| `nav-fundo` | `#14202E` | `#182330` | `#14202E` |
| `nav-texto` | `#F2F4F6` | `#EEF2F6` | `#FFFFFF` |
| `nav-texto-secundario` | `#A9B6C3` | `#A9B6C3` | `#D9DFE6` |
| `nav-item-ativo` | `#2B3A4C` | `#2B3A4C` | `#2B3A4C` |
| `toast-fundo` | `#14202E` | `#2B3A4C` | `#14202E` |
| `toast-texto` | `#F2F4F6` | `#E8EDF2` | `#FFFFFF` |
| `grade-grafico` | `#D9DFE6` | `#263241` | `#9AA6B2` |
| `serie-1` | `#14202E` | `#E8EDF2` | `#14202E` |
| `serie-2` | `#4A5868` | `#A9B5C2` | `#4A5868` |
| `serie-3` | `#9AA6B2` | `#4A5868` | `#9AA6B2` |
| `veu` (fundo de modal) | `rgba(20,32,46,.56)` | `rgba(0,0,0,.64)` | `rgba(20,32,46,.64)` |

Status por tema (`status-<nome>`, `status-<nome>-fundo`, `status-<nome>-texto`):

| Token | Claro | Escuro | Sol |
|---|---|---|---|
| `status-<nome>` (ícone, barra, gráfico) | cor clara da marca | cor escura da marca | cor clara da marca |
| `status-<nome>-fundo` (chip) | fundo claro da marca | fundo escuro da marca | cor clara da marca (chip sólido) |
| `status-<nome>-texto` (texto do chip) | `#14202E` | cor escura da marca | ver 7.1 (Aço ou branco) |

No tema claro, o texto do chip de status é Aço e só o ícone leva a cor. O padrão do brandbook, com texto na cor do status sobre o fundo tingido, reprova AA no claro (ver seção 4) e só é mantido no escuro, onde passa.

### 3.3 Tokens em CSS (web)

```css
/* tokens.css — Tenaz v0.1 */
:root {
  /* Primitivos */
  --tz-aco-900: #14202E;
  --tz-aco-700: #2B3A4C;
  --tz-aco-500: #4A5868;
  --tz-aco-400: #5F6C7B; /* extensão proposta */
  --tz-aco-300: #9AA6B2; /* Concreto */
  --tz-aco-100: #D9DFE6;
  --tz-aco-50:  #F2F4F6; /* Névoa */
  --tz-sinal: #FF6A13;
  --tz-sinal-hover: #FF7A2E;
  --tz-sinal-pressionado: #F25F0A;
  --tz-branco: #FFFFFF;

  /* Semânticos — tema claro (padrão) */
  --tz-fundo: var(--tz-aco-50);
  --tz-superficie: var(--tz-branco);
  --tz-superficie-elevada: var(--tz-branco);
  --tz-superficie-hover: var(--tz-aco-50);
  --tz-superficie-selecionada: var(--tz-aco-100);
  --tz-superficie-esqueleto: var(--tz-aco-100);
  --tz-texto: var(--tz-aco-900);
  --tz-texto-secundario: var(--tz-aco-500);
  --tz-texto-terciario: var(--tz-aco-400);
  --tz-texto-desabilitado: var(--tz-aco-300);
  --tz-texto-inverso: var(--tz-aco-50);
  --tz-borda: var(--tz-aco-100);
  --tz-borda-controle: #8592A0;
  --tz-borda-controle-hover: var(--tz-aco-500);
  --tz-foco: var(--tz-aco-900);
  --tz-acao-primaria: var(--tz-sinal);
  --tz-acao-primaria-hover: var(--tz-sinal-hover);
  --tz-acao-primaria-pressionada: var(--tz-sinal-pressionado);
  --tz-acao-primaria-texto: var(--tz-aco-900);
  --tz-perigo: #D33A2C;
  --tz-perigo-hover: #B8301F;
  --tz-perigo-texto: var(--tz-branco);
  --tz-nav-fundo: var(--tz-aco-900);
  --tz-nav-texto: var(--tz-aco-50);
  --tz-nav-texto-secundario: #A9B6C3;
  --tz-nav-item-ativo: var(--tz-aco-700);
  --tz-toast-fundo: var(--tz-aco-900);
  --tz-toast-texto: var(--tz-aco-50);
  --tz-grade-grafico: var(--tz-aco-100);
  --tz-serie-1: var(--tz-aco-900);
  --tz-serie-2: var(--tz-aco-500);
  --tz-serie-3: var(--tz-aco-300);
  --tz-veu: rgba(20, 32, 46, .56);

  --tz-status-operando: #1E9E6A;   --tz-status-operando-fundo: #E3F4EC;   --tz-status-operando-texto: var(--tz-aco-900);
  --tz-status-preventiva: #2F6FDB; --tz-status-preventiva-fundo: #E4ECFB; --tz-status-preventiva-texto: var(--tz-aco-900);
  --tz-status-atencao: #C98C00;    --tz-status-atencao-fundo: #FBF1D9;    --tz-status-atencao-texto: var(--tz-aco-900);
  --tz-status-parado: #D33A2C;     --tz-status-parado-fundo: #FBE5E2;     --tz-status-parado-texto: var(--tz-aco-900);
  --tz-status-inativo: #7E8A97;    --tz-status-inativo-fundo: #ECEFF2;    --tz-status-inativo-texto: var(--tz-aco-900);
  /* Cores de status para usar sobre superfícies Aço (toast, menu): variantes escuras */
  --tz-status-operando-sobre-aco: #3CC28B;
  --tz-status-preventiva-sobre-aco: #6E9BF0;
  --tz-status-atencao-sobre-aco: #F0B83A;
  --tz-status-parado-sobre-aco: #F06A5C;
  --tz-status-inativo-sobre-aco: #97A3B0;

  /* Tipografia */
  --tz-fonte-titulo: 'Archivo', 'Arial Black', 'Helvetica Neue', Arial, sans-serif;
  --tz-fonte-texto: 'Inter', system-ui, -apple-system, 'Segoe UI', Roboto, sans-serif;
  --tz-fonte-dado: 'JetBrains Mono', ui-monospace, 'Cascadia Mono', Consolas, monospace;
  --tz-largura-titulo: 125%; /* Archivo Expanded */

  --tz-texto-display:   800 48px/52px var(--tz-fonte-titulo);
  --tz-texto-h1:        800 32px/38px var(--tz-fonte-titulo);
  --tz-texto-h2:        700 22px/28px var(--tz-fonte-titulo);
  --tz-texto-h3:        700 18px/24px var(--tz-fonte-texto);
  --tz-texto-h4:        600 16px/22px var(--tz-fonte-texto);
  --tz-texto-corpo:     400 16px/24px var(--tz-fonte-texto);
  --tz-texto-corpo-p:   400 14px/20px var(--tz-fonte-texto);
  --tz-texto-rotulo:    600 14px/20px var(--tz-fonte-texto);
  --tz-texto-botao:     600 14px/20px var(--tz-fonte-texto);
  --tz-texto-legenda:   500 13px/18px var(--tz-fonte-texto);
  --tz-texto-sobrelinha: 600 12px/16px var(--tz-fonte-dado); /* caixa alta, espaçamento .08em */
  --tz-texto-dado:      400 14px/20px var(--tz-fonte-dado);
  --tz-texto-dado-forte: 600 14px/20px var(--tz-fonte-dado);
  --tz-texto-kpi:       600 28px/32px var(--tz-fonte-dado);

  /* Espaçamento — base 4 px */
  --tz-espaco-0: 0;
  --tz-espaco-1: 4px;
  --tz-espaco-2: 8px;
  --tz-espaco-3: 12px;
  --tz-espaco-4: 16px;
  --tz-espaco-5: 20px;
  --tz-espaco-6: 24px;
  --tz-espaco-8: 32px;
  --tz-espaco-10: 40px;
  --tz-espaco-12: 48px;
  --tz-espaco-16: 64px;
  --tz-espaco-20: 80px;

  /* Alturas de controle */
  --tz-controle-p: 32px;
  --tz-controle-m: 40px;
  --tz-controle-g: 48px;

  /* Raios */
  --tz-raio-xs: 4px;    /* badge, checkbox */
  --tz-raio-sm: 6px;    /* chip de status (brandbook) */
  --tz-raio-md: 8px;    /* botão, campo (brandbook) */
  --tz-raio-lg: 10px;   /* cartão, painel (brandbook) */
  --tz-raio-xl: 16px;   /* folha inferior e modal no app */
  --tz-raio-total: 999px;

  /* Sombras */
  --tz-sombra-1: 0 1px 2px rgba(20, 32, 46, .06);
  --tz-sombra-2: 0 1px 2px rgba(20, 32, 46, .06), 0 2px 6px rgba(20, 32, 46, .08);
  --tz-sombra-3: 0 1px 2px rgba(20, 32, 46, .08), 0 8px 24px rgba(20, 32, 46, .10);
  --tz-sombra-4: 0 2px 4px rgba(20, 32, 46, .08), 0 12px 32px rgba(20, 32, 46, .18);

  /* Camadas */
  --tz-z-base: 0;
  --tz-z-fixo: 100;      /* cabeçalho de tabela, colunas fixas */
  --tz-z-nav: 200;       /* menu lateral, barra superior */
  --tz-z-suspenso: 300;  /* menus, popovers, seletores */
  --tz-z-veu: 400;
  --tz-z-modal: 500;
  --tz-z-toast: 600;
  --tz-z-dica: 700;      /* tooltip */

  /* Movimento */
  --tz-duracao-rapida: 100ms;  /* hover, pressionado */
  --tz-duracao-base: 160ms;    /* menus, tooltips, troca de estado */
  --tz-duracao-media: 240ms;   /* modais, painéis laterais */
  --tz-duracao-lenta: 320ms;   /* folhas inferiores, transições de tela */
  --tz-curva-padrao: cubic-bezier(.2, 0, 0, 1);
  --tz-curva-entrada: cubic-bezier(0, 0, .2, 1);
  --tz-curva-saida: cubic-bezier(.4, 0, 1, 1);

  /* Foco */
  --tz-foco-largura: 2px;
  --tz-foco-afastamento: 2px;

  color-scheme: light;
}

/* Tema escuro: automático, salvo se o usuário escolheu claro */
@media (prefers-color-scheme: dark) {
  :root:not([data-tema="claro"]):not([data-tema="sol"]) { /* mesmos valores do bloco abaixo */ }
}
:root[data-tema="escuro"] {
  --tz-fundo: #0E151E;
  --tz-superficie: #16202C;
  --tz-superficie-elevada: #1F2B3A;
  --tz-superficie-hover: #1F2B3A;
  --tz-superficie-selecionada: var(--tz-aco-700);
  --tz-superficie-esqueleto: #1F2B3A;
  --tz-texto: #E8EDF2;
  --tz-texto-secundario: #A9B5C2;
  --tz-texto-terciario: #8592A0;
  --tz-texto-desabilitado: var(--tz-aco-500);
  --tz-texto-inverso: var(--tz-aco-900);
  --tz-borda: #263241;
  --tz-borda-controle: #6B7887;
  --tz-borda-controle-hover: #A9B5C2;
  --tz-foco: #E8EDF2;
  --tz-perigo: #F06A5C;
  --tz-perigo-hover: #F48577;
  --tz-perigo-texto: var(--tz-aco-900);
  --tz-nav-fundo: #182330;
  --tz-nav-texto: #EEF2F6;
  --tz-toast-fundo: var(--tz-aco-700);
  --tz-toast-texto: #E8EDF2;
  --tz-grade-grafico: #263241;
  --tz-serie-1: #E8EDF2;
  --tz-serie-2: #A9B5C2;
  --tz-serie-3: var(--tz-aco-500);
  --tz-veu: rgba(0, 0, 0, .64);

  --tz-status-operando: #3CC28B;   --tz-status-operando-fundo: #12301F;   --tz-status-operando-texto: #3CC28B;
  --tz-status-preventiva: #6E9BF0; --tz-status-preventiva-fundo: #17264A; --tz-status-preventiva-texto: #6E9BF0;
  --tz-status-atencao: #F0B83A;    --tz-status-atencao-fundo: #3A2D0C;    --tz-status-atencao-texto: #F0B83A;
  --tz-status-parado: #F06A5C;     --tz-status-parado-fundo: #3D1714;     --tz-status-parado-texto: #F06A5C;
  --tz-status-inativo: #97A3B0;    --tz-status-inativo-fundo: #222C37;    --tz-status-inativo-texto: #97A3B0;

  --tz-sombra-1: 0 1px 2px rgba(0, 0, 0, .40);
  --tz-sombra-2: 0 2px 6px rgba(0, 0, 0, .45);
  --tz-sombra-3: 0 8px 24px rgba(0, 0, 0, .50);
  --tz-sombra-4: 0 12px 32px rgba(0, 0, 0, .60);
  color-scheme: dark;
}

/* Movimento reduzido: sem deslocamento, sem brilho de esqueleto, fades curtos */
@media (prefers-reduced-motion: reduce) {
  :root {
    --tz-duracao-base: 0ms;
    --tz-duracao-media: 0ms;
    --tz-duracao-lenta: 0ms;
  }
  *, *::before, *::after {
    animation-duration: 1ms !important;
    animation-iteration-count: 1 !important;
    scroll-behavior: auto !important;
  }
}

/* Foco visível padrão */
:focus-visible {
  outline: var(--tz-foco-largura) solid var(--tz-foco);
  outline-offset: var(--tz-foco-afastamento);
}

/* Títulos usam Archivo Expanded */
h1, h2, .tz-titulo { font-family: var(--tz-fonte-titulo); font-stretch: var(--tz-largura-titulo); letter-spacing: -0.01em; text-wrap: balance; }

/* Dados */
.tz-dado, td.numero { font-family: var(--tz-fonte-dado); font-variant-numeric: tabular-nums; }
```

O tema Sol também existe no web (`:root[data-tema="sol"]`, valores da coluna Sol em 3.2), para tablets usados em pátio e oficina. No CSS final, o bloco `prefers-color-scheme` repete os valores do tema escuro (CSS não permite herdar um seletor de outro). Gere os dois blocos a partir de um único arquivo de tokens (`tokens.json`) com Style Dictionary ou script equivalente, para não haver divergência.

### 3.4 Tokens em React Native (app)

```ts
// tokens.ts — Tenaz v0.1 (app de campo)
import { Platform } from 'react-native';

const marca = {
  aco900: '#14202E', aco700: '#2B3A4C', aco500: '#4A5868', aco400: '#5F6C7B',
  aco300: '#9AA6B2', aco100: '#D9DFE6', aco50: '#F2F4F6',
  sinal: '#FF6A13', sinalHover: '#FF7A2E', sinalPressionado: '#F25F0A',
  branco: '#FFFFFF',
} as const;

const statusClaro = {
  operando:   { cor: '#1E9E6A', fundo: '#E3F4EC', texto: marca.aco900 },
  preventiva: { cor: '#2F6FDB', fundo: '#E4ECFB', texto: marca.aco900 },
  atencao:    { cor: '#C98C00', fundo: '#FBF1D9', texto: marca.aco900 },
  parado:     { cor: '#D33A2C', fundo: '#FBE5E2', texto: marca.aco900 },
  inativo:    { cor: '#7E8A97', fundo: '#ECEFF2', texto: marca.aco900 },
};
const statusEscuro = {
  operando:   { cor: '#3CC28B', fundo: '#12301F', texto: '#3CC28B' },
  preventiva: { cor: '#6E9BF0', fundo: '#17264A', texto: '#6E9BF0' },
  atencao:    { cor: '#F0B83A', fundo: '#3A2D0C', texto: '#F0B83A' },
  parado:     { cor: '#F06A5C', fundo: '#3D1714', texto: '#F06A5C' },
  inativo:    { cor: '#97A3B0', fundo: '#222C37', texto: '#97A3B0' },
};
// Tema sol: chip sólido, texto conforme contraste (ver seção 7.1)
const statusSol = {
  operando:   { cor: '#1E9E6A', fundo: '#1E9E6A', texto: marca.aco900 },
  preventiva: { cor: '#2F6FDB', fundo: '#2F6FDB', texto: marca.branco },
  atencao:    { cor: '#C98C00', fundo: '#C98C00', texto: marca.aco900 },
  parado:     { cor: '#D33A2C', fundo: '#D33A2C', texto: marca.branco },
  inativo:    { cor: '#7E8A97', fundo: '#7E8A97', texto: marca.aco900 },
};

export const cores = {
  claro: {
    fundo: marca.aco50, superficie: marca.branco, superficieElevada: marca.branco,
    superficieSelecionada: marca.aco100, esqueleto: marca.aco100,
    texto: marca.aco900, textoSecundario: marca.aco500, textoTerciario: marca.aco400,
    textoDesabilitado: marca.aco300, textoInverso: marca.aco50,
    borda: marca.aco100, bordaControle: '#8592A0', foco: marca.aco900,
    acaoPrimaria: marca.sinal, acaoPrimariaPressionada: marca.sinalPressionado, acaoPrimariaTexto: marca.aco900,
    perigo: '#D33A2C', perigoTexto: marca.branco,
    barraOffline: marca.aco900, barraOfflineTexto: marca.aco50,
    veu: 'rgba(20,32,46,0.56)',
    status: statusClaro,
  },
  escuro: {
    fundo: '#0E151E', superficie: '#16202C', superficieElevada: '#1F2B3A',
    superficieSelecionada: marca.aco700, esqueleto: '#1F2B3A',
    texto: '#E8EDF2', textoSecundario: '#A9B5C2', textoTerciario: '#8592A0',
    textoDesabilitado: marca.aco500, textoInverso: marca.aco900,
    borda: '#263241', bordaControle: '#6B7887', foco: '#E8EDF2',
    acaoPrimaria: marca.sinal, acaoPrimariaPressionada: marca.sinalPressionado, acaoPrimariaTexto: marca.aco900,
    perigo: '#F06A5C', perigoTexto: marca.aco900,
    barraOffline: marca.aco700, barraOfflineTexto: '#E8EDF2',
    veu: 'rgba(0,0,0,0.64)',
    status: statusEscuro,
  },
  sol: {
    fundo: marca.branco, superficie: marca.branco, superficieElevada: marca.branco,
    superficieSelecionada: marca.aco100, esqueleto: marca.aco100,
    texto: marca.aco900, textoSecundario: marca.aco700, textoTerciario: marca.aco500,
    textoDesabilitado: marca.aco300, textoInverso: marca.branco,
    borda: marca.aco300, bordaControle: marca.aco500, foco: marca.aco900,
    acaoPrimaria: marca.sinal, acaoPrimariaPressionada: marca.sinalPressionado, acaoPrimariaTexto: marca.aco900,
    perigo: '#D33A2C', perigoTexto: marca.branco,
    barraOffline: marca.aco900, barraOfflineTexto: marca.branco,
    veu: 'rgba(20,32,46,0.64)',
    status: statusSol,
  },
} as const;

// Fontes: Inter e JetBrains Mono via @expo-google-fonts; Archivo Expanded como
// instâncias estáticas (largura 125; pesos 700, 800, 900) em assets/fonts.
export const fontes = {
  tituloBold: 'ArchivoExpanded-Bold',
  tituloExtraBold: 'ArchivoExpanded-ExtraBold',
  texto: 'Inter_400Regular',
  textoMedio: 'Inter_500Medium',
  textoSemi: 'Inter_600SemiBold',
  textoBold: 'Inter_700Bold',
  dado: 'JetBrainsMono_400Regular',
  dadoSemi: 'JetBrainsMono_600SemiBold',
};

export const tipografia = {
  display:    { fontFamily: fontes.tituloExtraBold, fontSize: 32, lineHeight: 36 },
  titulo1:    { fontFamily: fontes.tituloExtraBold, fontSize: 26, lineHeight: 32 },
  titulo2:    { fontFamily: fontes.tituloBold,      fontSize: 20, lineHeight: 26 },
  titulo3:    { fontFamily: fontes.textoBold,       fontSize: 18, lineHeight: 24 },
  corpo:      { fontFamily: fontes.texto,           fontSize: 17, lineHeight: 24 },
  corpoForte: { fontFamily: fontes.textoSemi,       fontSize: 17, lineHeight: 24 },
  rotulo:     { fontFamily: fontes.textoSemi,       fontSize: 15, lineHeight: 20 },
  botao:      { fontFamily: fontes.textoSemi,       fontSize: 17, lineHeight: 24 },
  legenda:    { fontFamily: fontes.textoMedio,      fontSize: 14, lineHeight: 20 },
  sobrelinha: { fontFamily: fontes.dadoSemi,        fontSize: 13, lineHeight: 16, letterSpacing: 1, textTransform: 'uppercase' as const },
  dado:       { fontFamily: fontes.dadoSemi,        fontSize: 17, lineHeight: 24, fontVariant: ['tabular-nums' as const] },
  dadoKpi:    { fontFamily: fontes.dadoSemi,        fontSize: 28, lineHeight: 32, fontVariant: ['tabular-nums' as const] },
};

export const espaco = { 0: 0, 1: 4, 2: 8, 3: 12, 4: 16, 5: 20, 6: 24, 8: 32, 10: 40, 12: 48, 16: 64 } as const;
export const raio = { xs: 4, sm: 6, md: 8, lg: 10, xl: 16, total: 999 } as const;

export const alvoToque = { minimo: 48, primario: 56, espacoEntreAlvos: 8 } as const;

export const sombra = {
  1: Platform.select({
    ios: { shadowColor: '#14202E', shadowOffset: { width: 0, height: 1 }, shadowOpacity: 0.08, shadowRadius: 2 },
    android: { elevation: 1 },
  }),
  2: Platform.select({
    ios: { shadowColor: '#14202E', shadowOffset: { width: 0, height: 2 }, shadowOpacity: 0.10, shadowRadius: 6 },
    android: { elevation: 3 },
  }),
  3: Platform.select({
    ios: { shadowColor: '#14202E', shadowOffset: { width: 0, height: 8 }, shadowOpacity: 0.12, shadowRadius: 24 },
    android: { elevation: 8 },
  }),
  4: Platform.select({
    ios: { shadowColor: '#14202E', shadowOffset: { width: 0, height: 12 }, shadowOpacity: 0.18, shadowRadius: 32 },
    android: { elevation: 16 },
  }),
};

export const movimento = {
  rapida: 100, base: 160, media: 240, lenta: 320,
  // Easing equivalente ao web, para Reanimated: Easing.bezier(0.2, 0, 0, 1)
  curvaPadrao: [0.2, 0, 0, 1] as const,
  curvaEntrada: [0, 0, 0.2, 1] as const,
  curvaSaida: [0.4, 0, 1, 1] as const,
};
```

Regras do app:

- O tema segue o sistema (`useColorScheme`) e pode ser fixado pelo usuário em Claro, Escuro ou Sol. Em Sol a escolha é manual; não há detecção automática de luz.
- Respeite `AccessibilityInfo.isReduceMotionEnabled()` (ou `useReducedMotion` do Reanimated): com movimento reduzido, durações viram 0 e o esqueleto fica estático.
- `allowFontScaling` sempre ativo. Teste com escala de fonte de 200 % do sistema; limite `maxFontSizeMultiplier` a 2 só em dados de KPI que quebrariam o layout.

### 3.5 Tipografia

**Famílias.** Archivo Expanded (largura 125) só para títulos de 22 px para cima, pesos 700 a 900. Inter para interface e texto, 400 a 700. JetBrains Mono para dados, códigos, TAG, placa, número de OS, valores e unidades, 400 e 600.

**Carregamento no web.** Hospede as fontes junto com a aplicação (pacotes Fontsource com eixos variáveis, incluindo o eixo de largura da Archivo) em vez do CDN do Google. Isso evita enviar o IP do usuário a terceiros (LGPD) e mantém a interface correta em redes corporativas bloqueadas. Use `font-display: swap` e pré-carregue Inter 400 e 600.

**Escala web** (os valores marcados com "marca" vêm do brandbook):

| Token | Família | Tamanho / altura de linha | Peso | Uso |
|---|---|---|---|---|
| `display` | Archivo Expanded | 48 / 52 | 800 | Número de destaque em tela de boas-vindas, página de status. Marca. |
| `h1` | Archivo Expanded | 32 / 38 | 800 | Título do painel e de telas de primeiro uso. Marca. |
| `h2` | Archivo Expanded | 22 / 28 | 700 | Título de página nas telas de trabalho. Marca. |
| `h3` | Inter | 18 / 24 | 700 | Título de seção, de cartão e de modal |
| `h4` | Inter | 16 / 22 | 600 | Subtítulo de grupo de campos |
| `corpo` | Inter | 16 / 24 | 400 | Texto corrido, descrições. Marca. |
| `corpo-p` | Inter | 14 / 20 | 400 | Células de tabela, formulários densos, menus |
| `rotulo` | Inter | 14 / 20 | 600 | Rótulo de campo, cabeçalho de tabela |
| `botao` | Inter | 14 / 20 | 600 | Texto de botão (brandbook usa 0,92 rem, arredondado para 14) |
| `legenda` | Inter | 13 / 18 | 500 | Texto de apoio, "Atualizado há 4 min". Marca. |
| `sobrelinha` | JetBrains Mono | 12 / 16 | 600 | Rótulo curto em caixa alta, espaçamento 0,08 em, como no brandbook |
| `dado` | JetBrains Mono | 14 / 20 | 400 | Valores, códigos, datas em tabela. Marca. |
| `dado-forte` | JetBrains Mono | 14 / 20 | 600 | Valores em destaque, TAG, número de OS |
| `kpi` | JetBrains Mono | 28 / 32 | 600 | Valor de indicador em cartão de KPI |

**Escala mobile (dp)**:

| Token | Família | Tamanho / altura | Peso | Uso |
|---|---|---|---|---|
| `display` | Archivo Expanded | 32 / 36 | 800 | Número de destaque (OS do dia) |
| `titulo1` | Archivo Expanded | 26 / 32 | 800 | Título de tela raiz |
| `titulo2` | Archivo Expanded | 20 / 26 | 700 | Título de tela interna |
| `titulo3` | Inter | 18 / 24 | 700 | Título de cartão e de seção |
| `corpo` | Inter | 17 / 24 | 400 | Texto padrão |
| `corpo-forte` | Inter | 17 / 24 | 600 | Nome do ativo, item de checklist |
| `rotulo` | Inter | 15 / 20 | 600 | Rótulo de campo |
| `botao` | Inter | 17 / 24 | 600 | Texto de botão |
| `legenda` | Inter | 14 / 20 | 500 | Texto de apoio (menor tamanho permitido no app) |
| `sobrelinha` | JetBrains Mono | 13 / 16 | 600 | Rótulo em caixa alta |
| `dado` | JetBrains Mono | 17 / 24 | 600 | Leituras, TAG, placa |
| `dado-kpi` | JetBrains Mono | 28 / 32 | 600 | Leitura em destaque (hodômetro) |

Regras:

- Archivo Expanded nunca abaixo de 20 dp no app nem de 22 px no web. Abaixo disso, títulos usam Inter 700.
- Caixa alta só na `sobrelinha`, com no máximo três palavras. Nunca em frases, botões ou alertas.
- Texto corrido com no máximo 72 caracteres por linha.
- Números com `font-variant-numeric: tabular-nums` (Inter) ou JetBrains Mono, sempre alinhados à direita em colunas.

### 3.6 Espaçamento

Escala de 4 px. Valores fora da escala só para bordas (1 px e 2 px) e ajustes ópticos de ícone.

| Token | Valor | Uso típico |
|---|---|---|
| `espaco-1` | 4 px | Entre ícone e texto em chip |
| `espaco-2` | 8 px | Entre ícone e texto em botão, entre chips, mínimo entre alvos de toque |
| `espaco-3` | 12 px | Padding vertical de célula densa, entre rótulo e campo |
| `espaco-4` | 16 px | Padding de cartão no app, margem lateral do app, gutter da grade abaixo de 1280 px |
| `espaco-5` | 20 px | Padding de cartão no web |
| `espaco-6` | 24 px | Gutter da grade a partir de 1280 px, margem da página no web, entre seções de formulário |
| `espaco-8` | 32 px | Entre blocos do painel |
| `espaco-10` | 40 px | Altura do controle médio |
| `espaco-12` | 48 px | Alvo de toque mínimo, altura de linha confortável |
| `espaco-16` | 64 px | Margem de estado vazio |
| `espaco-20` | 80 px | Espaço reservado ao botão flutuante no app |

### 3.7 Raios

| Token | Valor | Uso |
|---|---|---|
| `raio-xs` | 4 px | Badge de criticidade, checkbox, barra de progresso |
| `raio-sm` | 6 px | Chip de status (brandbook) |
| `raio-md` | 8 px | Botão, campo, seletor, toast (brandbook) |
| `raio-lg` | 10 px | Cartão, painel, tabela, modal no web (brandbook) |
| `raio-xl` | 16 px | Folha inferior e modal no app (só cantos superiores na folha) |
| `raio-total` | 999 px | Avatar, contador numérico, interruptor |

Cantos arredondados combinam com o traço arredondado dos ícones. Nada com raio maior que 16 px, exceto pílulas.

### 3.8 Sombras e elevação

A interface é plana. A separação vem primeiro de borda e de cor de superfície; a sombra só indica que algo está acima do conteúdo.

| Nível | Elementos | Claro | Escuro | z-index |
|---|---|---|---|---|
| 0 | Página, seções | sem sombra | sem sombra | `z-base` |
| 1 | Cartão, tabela, painel | `sombra-1` + `borda` | `borda`, sem sombra perceptível | `z-base` |
| 2 | Cabeçalho fixo de tabela, cartão sendo arrastado, colunas fixas com rolagem | `sombra-2` | `sombra-2` + `superficie-elevada` | `z-fixo` |
| 3 | Menu suspenso, seletor, popover, "como calculamos" | `sombra-3` | `sombra-3` + `superficie-elevada` | `z-suspenso` |
| 4 | Modal, toast, folha inferior | `sombra-4` | `sombra-4` + `superficie-elevada` + `borda` | `z-modal`, `z-toast` |

No tema escuro a elevação é marcada pela superfície mais clara (`#16202C` para `#1F2B3A`), porque a sombra quase não aparece.

### 3.9 Movimento

| Token | Duração | Curva | Uso |
|---|---|---|---|
| `duracao-rapida` | 100 ms | padrão | Hover, pressionado, troca de cor |
| `duracao-base` | 160 ms | entrada / saída | Menu, tooltip, chip de filtro, expandir linha |
| `duracao-media` | 240 ms | entrada / saída | Modal, painel lateral, toast |
| `duracao-lenta` | 320 ms | padrão | Folha inferior e troca de tela no app |

Regras:

- Movimento explica a origem de algo (o painel vem da direita porque abriu de uma linha). Não há animação decorativa, efeito elástico ou parallax.
- Saída é mais rápida que entrada (use 70 % da duração).
- Com `prefers-reduced-motion: reduce` (web) ou "Reduzir movimento" ativo (iOS e Android): sem deslocamento nem escala; só troca de opacidade de até 100 ms; esqueleto sem brilho; rolagem sem suavização.
- Indicador de progresso giratório é permitido mesmo com movimento reduzido, porque transmite estado; prefira barra de progresso quando o tempo for previsível.

### 3.10 Breakpoints e grade

| Token | Largura mínima | Comportamento |
|---|---|---|
| `bp-sm` | 640 px | Formulários em uma coluna. Web móvel só para solicitação e consulta rápida. |
| `bp-md` | 768 px | Tablet. Menu lateral vira gaveta. Tabelas mostram colunas prioritárias e o resto em linha expansível. |
| `bp-lg` | 1024 px | Menu lateral recolhido (64 px) por padrão. Detalhe em duas colunas. |
| `bp-xl` | 1280 px | Menu lateral expandido (240 px). Grade com gutter de 24 px. |
| `bp-2xl` | 1536 px | Painel com 4 cartões de KPI por linha; quadro de programação mostra 14 dias. |

CSS não aceita variáveis em `@media`; os valores ficam no `tokens.json` e são exportados para o CSS e para o hook `useBreakpoint`.

Grade do web: 12 colunas fluidas, gutter de 16 px abaixo de 1280 px e de 24 px a partir daí, margem da página de 24 px (32 px a partir de 1536 px). Tabelas e quadros usam a largura total. Formulários de cadastro limitam a 720 px de largura útil. Telas de detalhe usam 8 + 4 colunas (conteúdo e painel de contexto).

### 3.11 Números, unidades e datas

| Tipo | Formato | Exemplo |
|---|---|---|
| Inteiro | Ponto como separador de milhar | 248.315 km |
| Decimal | Vírgula, casas fixas por grandeza | 3,2 h · 97,4 % · 0,41 R$/km |
| Moeda | `R$` + espaço inseparável | R$ 3.418,90 |
| Percentual | Espaço inseparável antes de `%` | 97,4 % |
| Medidor | Valor + espaço inseparável + unidade | 12.480 h · 248.315 km · 1.204 ciclos |
| Data | dd/mm/aaaa | 12/11/2026 |
| Data e hora | dd/mm/aaaa hh:mm (24 h) | 12/11/2026 14:32 |
| Hora relativa | Até 24 h; depois, data | há 4 min · há 3 h · ontem às 16:05 |
| Duração | h e min | 3 h 20 min |
| Prazo futuro | Unidade do gatilho e, se possível, estimativa em dias | em 2.300 km (cerca de 9 dias) |
| Placa | Caixa alta, sem hífen no Mercosul | RTZ4F82 · ABC-1234 |
| OS / SS / OR / CL | Formato canônico completo `<PREFIXO>-<AAAA>-<NNNNNN>`, sequencial por empresa, reinicia a cada ano | OS-2026-004812 · SS-2026-000731 · OR-2026-000045 · CL-2026-000012 |

- Unidades em minúsculas conforme o SI (km, h, min, kg, L, °C, bar). "un." para unidades de estoque.
- Em tabela, coluna homogênea traz a unidade no cabeçalho ("Hodômetro (km)") e a célula só o número. Coluna mista (medidores de tipos diferentes) traz a unidade em cada célula.
- Formatação sempre com `Intl.NumberFormat('pt-BR')` e `Intl.DateTimeFormat('pt-BR')`, nunca concatenação manual.

---

## 4. Contraste WCAG 2.1

Razões calculadas com a fórmula de luminância relativa da WCAG 2.1. Critérios: texto normal 4,5:1 (1.4.3), texto grande (a partir de 24 px, ou 18,66 px em negrito) 3:1, componentes de interface e gráficos 3:1 (1.4.11).

### 4.1 Tema claro

| Primeiro plano | Fundo | Razão | Resultado | Decisão |
|---|---|---|---|---|
| Aço `#14202E` | Névoa `#F2F4F6` | 14,93 | AAA | Texto principal |
| Aço `#14202E` | Branco | 16,46 | AAA | Texto principal em superfície |
| Aço 700 `#2B3A4C` | Branco | 11,59 | AAA | Texto secundário no tema sol |
| Aço 700 `#2B3A4C` | Névoa | 10,51 | AAA | — |
| Aço 500 `#4A5868` | Branco | 7,27 | AAA | Texto secundário |
| Aço 500 `#4A5868` | Névoa | 6,60 | AA | Texto secundário |
| Aço 500 `#4A5868` | Aço 100 `#D9DFE6` | 5,42 | AA | Texto secundário em linha selecionada |
| Aço 400 `#5F6C7B` (proposto) | Branco | 5,36 | AA | Texto terciário |
| Aço 400 `#5F6C7B` | Névoa | 4,86 | AA | Texto terciário |
| Aço 400 `#5F6C7B` | Aço 100 | 3,99 | Reprova | Não usar em linha selecionada |
| `#6B7887` (ink-3 do brandbook) | Branco | 4,50 | AA no limite | Substituído por Aço 400 |
| `#6B7887` | Névoa | 4,09 | Reprova | Substituído por Aço 400 |
| Concreto `#9AA6B2` | Branco | 2,48 | Reprova | Só decorativo e desabilitado |
| Concreto `#9AA6B2` | Névoa | 2,25 | Reprova | Só decorativo e desabilitado |
| **Aço `#14202E`** | **Sinal `#FF6A13`** | **5,74** | **AA** | **Texto do botão primário** |
| **Branco** | **Sinal `#FF6A13`** | **2,87** | **Reprova** (até para texto grande) | **Proibido** |
| Névoa | Sinal | 2,60 | Reprova | Proibido |
| Sinal (texto ou contorno) | Branco | 2,87 | Reprova (também para 3:1 de UI) | Sinal não é texto, ícone nem foco no claro |
| Sinal (texto ou contorno) | Névoa | 2,60 | Reprova | Idem |
| Aço | Sinal hover `#FF7A2E` | 6,33 | AA | Botão primário em hover |
| Aço | Sinal pressionado `#F25F0A` | 5,03 | AA | Botão primário pressionado |
| `#B54500` (accent-text do brandbook) | Branco | 5,50 | AA | Não usado no produto (ver 13.7) |
| Branco | Parado `#D33A2C` | 4,76 | AA | Botão de perigo |
| Branco | `#B8301F` | 6,03 | AA | Botão de perigo em hover |
| Névoa | Aço `#14202E` | 14,93 | AAA | Texto do menu lateral e do toast |
| `#A9B6C3` | Aço `#14202E` | 7,97 | AAA | Texto secundário do menu |
| Concreto | Aço `#14202E` | 6,64 | AA | Ícones do menu |
| Branco | Aço 700 `#2B3A4C` | 11,59 | AAA | Item ativo do menu |
| Sinal | Aço `#14202E` | 5,74 | AA | Ponto do símbolo no menu |
| Borda de campo `#8592A0` | Branco | 3,17 | UI 3:1 | Contorno de campo |
| Borda decorativa Aço 100 | Branco | 1,34 | Decorativa | Divisórias, nunca contorno de campo |
| Foco Aço `#14202E` | Névoa / Branco | 14,93 / 16,46 | UI 3:1 | Anel de foco |

**Status no tema claro**

| Par | Operando | Preventiva | Atenção | Parado | Inativo |
|---|---|---|---|---|---|
| Cor do status como texto sobre o fundo do status (padrão do brandbook) | 2,99 reprova | 4,00 reprova | 2,58 reprova | 3,95 reprova | 3,05 reprova |
| Aço como texto sobre o fundo do status (padrão adotado) | 14,44 AAA | 13,87 AAA | 14,65 AAA | 13,64 AAA | 14,26 AAA |
| Cor do status como ícone ou gráfico sobre branco (mínimo 3) | 3,41 passa | 4,75 passa | **2,90 reprova** | 4,76 passa | 3,52 passa |
| Cor do status como ícone ou gráfico sobre Névoa | 3,09 passa | 4,31 passa | **2,63 reprova** | 4,32 passa | 3,19 passa |
| Chip sólido com texto Aço | 4,82 AA | 3,46 reprova | 5,67 AA | 3,46 reprova | 4,68 AA |
| Chip sólido com texto branco | 3,41 reprova | 4,75 AA | 2,90 reprova | 4,76 AA | 3,52 reprova |
| Variante escura do status sobre Aço (toast, menu) | 7,28 | 5,96 | 9,11 | 5,42 | 6,41 |

### 4.2 Tema escuro

| Primeiro plano | Fundo | Razão | Resultado | Uso |
|---|---|---|---|---|
| `#E8EDF2` | `#0E151E` | 15,57 | AAA | Texto principal sobre fundo |
| `#E8EDF2` | `#16202C` | 13,96 | AAA | Texto principal sobre superfície |
| `#E8EDF2` | `#1F2B3A` | 12,17 | AAA | Texto em modal e menu |
| `#E8EDF2` | `#2B3A4C` | 9,84 | AAA | Linha selecionada, toast |
| `#A9B5C2` | `#16202C` | 7,89 | AAA | Texto secundário |
| `#A9B5C2` | `#1F2B3A` | 6,88 | AA | Texto secundário em modal |
| `#A9B5C2` | `#2B3A4C` | 5,56 | AA | Texto secundário em linha selecionada |
| `#8592A0` | `#0E151E` | 5,78 | AA | Texto terciário |
| `#8592A0` | `#16202C` | 5,18 | AA | Texto terciário |
| Aço | Sinal | 5,74 | AA | Botão primário |
| Sinal | `#0E151E` | 6,40 | UI | Botão primário se destaca do fundo |
| Aço | `#F06A5C` | 5,42 | AA | Botão de perigo |
| Aço | `#F48577` | 6,64 | AA | Botão de perigo em hover |
| Borda de campo `#6B7887` | `#16202C` | 3,65 | UI 3:1 | Contorno de campo |
| Borda decorativa `#263241` | `#16202C` | 1,26 | Decorativa | Divisórias |
| `#EEF2F6` | `#182330` | 14,12 | AAA | Menu lateral |
| `#A9B6C3` | `#182330` | 7,69 | AAA | Texto secundário do menu |

| Par (escuro) | Operando | Preventiva | Atenção | Parado | Inativo |
|---|---|---|---|---|---|
| Cor do status sobre superfície `#16202C` | 7,27 | 5,96 | 9,10 | 5,42 | 6,41 |
| Cor do status como texto sobre o fundo do status (brandbook) | 6,32 AA | 5,39 AA | 7,45 AAA | 5,20 AA | 5,52 AA |
| `#E8EDF2` sobre o fundo do status | 12,14 | 12,62 | 11,42 | 13,39 | 12,02 |
| Ícone de status sobre toast `#2B3A4C` | 5,12 | — | 6,41 | 3,82 (só ícone) | — |

### 4.3 Regras derivadas

1. **Texto sobre Sinal é sempre Aço.** Branco sobre Sinal (2,87:1) não atinge nem os 3:1 de texto grande.
2. **Sinal não é texto, ícone, contorno ou anel de foco sobre fundo claro.** O anel de foco é Aço no claro e `#E8EDF2` no escuro.
3. **Chip de status no claro: fundo tingido, ícone na cor do status, texto Aço.** No escuro vale o padrão do brandbook (texto e ícone na cor do status), que passa.
4. **Atenção no claro precisa de reforço gráfico.** O âmbar `#C98C00` tem 2,90:1 sobre branco. Como marca isolada (barra de gráfico, ponto, ícone sem texto), recebe contorno de 1 px em Aço ou fica sobre o fundo `#FBF1D9` com o texto ao lado. Ícone de Atenção desenhado com triângulo âmbar preenchido e ponto de exclamação em Aço (5,67:1) quando precisar aparecer sozinho, por exemplo em marcador de mapa.
5. **Concreto não é texto legível.** Texto secundário usa Aço 500; terciário usa Aço 400.
6. **Contorno de campo usa `borda-controle`**, nunca a borda decorativa.
7. **Texto desabilitado** é isento de contraste (WCAG 1.4.3), mas o controle desabilitado tem de explicar por que está desabilitado (tooltip ou texto de apoio).

---

## 5. Layout do web

### 5.1 Estrutura

```
┌──────────┬──────────────────────────────────────────────────────────────────┐
│ [logo]   │ [Unidade: Planta Contagem ▾]  [ Buscar TAG, placa, OS ou peça  /]│
│          │                                    [+ Criar ▾] [sino 3] [?] [GF]│
│ Painel   ├──────────────────────────────────────────────────────────────────┤
│          │ Ordens de serviço › OS-2026-004812                               │
│ OPERAÇÃO │ ┌──────────────────────────────────────────────────────────────┐ │
│ Solicit. │ │ Título da página (h2)        [Secundária] [Ação primária]    │ │
│ OS       │ └──────────────────────────────────────────────────────────────┘ │
│ Program. │                                                                  │
│ ATIVOS   │   Conteúdo                                                       │
│ Ativos   │                                                                  │
│ Planos   │                                                                  │
│ Checkl.  │                                                                  │
│ FROTA    │                                                                  │
│ ...      │                                                                  │
│──────────│                                                                  │
│ Admin.   │                                                                  │
│ [<<]     │                                                          v1.4.2  │
└──────────┴──────────────────────────────────────────────────────────────────┘
```

O selo de versão ("v1.4.2", componente 6.19) fica fixo no canto inferior direito da área de conteúdo, em todas as páginas do web, inclusive nos portais do prestador e do locatário.

### 5.2 Menu lateral por módulo

- Fundo `nav-fundo` (Aço), 240 px expandido e 64 px recolhido. A escolha fica salva por usuário.
- Agrupado por área, com `sobrelinha` em `nav-texto-secundario`:

| Grupo | Itens (código do módulo) |
|---|---|
| — | Painel (IND) |
| Operação | Solicitações (SOL), Ordens de serviço (OS), Programação (OS) |
| Ativos | Ativos (ATV), Medidores (MED), Planos de manutenção (PRE), Checklists (CHK) |
| Frota | Veículos (FRO, ATV), Pneus, Abastecimentos, Multas, Sinistros |
| Suprimentos | Estoque (EST), Ordens de reparo (EST), Compras (EST) |
| Terceiros | Fornecedores (TER), Contratos de prestação (TER), Links de execução (TER), Conexões com prestadores (lado do cliente, 8.27) |
| Locação | Contratos de locação (LOC), Locatários (LOC), Vistorias (LOC) |
| Serviços | Painel de serviços (SRV), Distribuição (SRV), Clientes atendidos (SRV), Equipamentos de clientes (SRV), Técnicos e especialidades (SRV), Conexões com clientes (SRV). Só em contas com o plano Serviços (8.20 a 8.28) |
| Conformidade | Vencimentos (CON), NR-12, NR-13, Capacitações |
| Análise | Indicadores e relatórios (IND), Assistente (IA) |
| Rodapé | Integrações (INT), Administração (ADM), recolher menu |

- Item ativo: fundo `nav-item-ativo` (Aço 700), barra de 3 px à esquerda em Névoa, texto branco, `aria-current="page"`. O item ativo **não** usa Sinal.
- Item recolhido mostra só o ícone, com tooltip e `aria-label` com o nome.
- Os perfis `prestador` e `locatario` não usam este menu: entram nos portais restritos (8.17 e 8.19).
- **Conta só com o plano Serviços** (subcontratada sem plano Planta ou Frota): o menu mostra Painel de serviços no lugar do Painel de manutenção, o grupo Operação (OS e Distribuição), o grupo Serviços logo abaixo, Estoque (estoque básico) e Indicadores. Os rótulos usam o vocabulário de quem presta serviço ("Clientes atendidos", "Equipamentos de clientes"), mas as telas de OS, apontamento e ativo são as mesmas do restante do produto.
- Grupos e itens aparecem conforme perfil e plano. Item sem permissão não aparece. Recurso de plano superior aparece só em Administração > Plano, nunca como item bloqueado no menu do dia a dia.
- Contadores de pendência ao lado do item (ex.: "Solicitações 7") em pílula `superficie` com texto Aço, nunca vermelha.

### 5.3 Barra superior

- Altura 56 px, fundo `superficie`, borda inferior `borda`.
- **Seletor de unidade:** à esquerda, mostra a unidade atual em `rotulo`. Opções: "Todas as unidades" (gestor e administrador) e a lista com busca quando houver mais de 8. Trocar de unidade mantém a página e os filtros e recarrega os dados. Só aparece quando a empresa tem mais de uma unidade.
- **Busca global:** campo central, até 560 px, atalhos `/` e `Ctrl+K`. Reconhece o tipo pelo padrão e mostra resultados agrupados:
  - número de OS, SS, OR ou CL completo ou parcial (`4812`, `OS-2026-004812`, `OR-2026-000045`), com Enter abrindo direto quando houver um único resultado;
  - número de série de peça serializada; fornecedor e locatário por nome ou CNPJ/CPF (CPF só com os dígitos mascarados no resultado);
  - placa antiga e Mercosul, com ou sem hífen, em qualquer caixa;
  - TAG de ativo (`cmp0347` encontra `CMP-0347`);
  - item de estoque por código ou descrição; pessoa por nome.
  Cada resultado mostra ícone do tipo, identificador em `dado-forte`, descrição e status. Setas navegam, Enter abre, Esc fecha. Busca vazia mostra os últimos 5 itens abertos.
- **Criar:** botão secundário com menu (Solicitação, OS, Ativo, Leitura de medidor). Não usa Sinal, para não competir com a ação primária da página.
- Notificações (sino com contador), ajuda e menu do usuário (perfil, tema, densidade, sair).

### 5.4 Cabeçalho da página

Trilha de navegação (`legenda`), título em `h2`, linha de metadados opcional (status, TAG, responsável) e ações à direita: no máximo uma primária (Sinal) e duas secundárias; o resto vai para o menu "Mais ações" (ícone `ellipsis`).

### 5.5 Densidade

| Densidade | Altura da linha de tabela | Altura do controle | Padding de célula | Padrão para |
|---|---|---|---|---|
| Compacta | 32 px | 32 px | 4 × 8 px | Planejador e almoxarife que pedirem |
| Padrão | 40 px | 40 px | 8 × 12 px | Todos no desktop |
| Confortável | 48 px | 48 px | 12 × 16 px | Tablets e telas de toque (automático com `pointer: coarse`) |

A densidade é preferência do usuário (atributo `data-densidade` no `<html>`). O tamanho da letra não muda com a densidade.

### 5.6 Tela de login

```
┌──────────────────────────────────────────────────────────────────────────┐
│                                                                          │
│                     [logo horizontal positivo]                           │
│                                                                          │
│               ┌──────────────────────────────────────┐                   │
│               │ Entrar (h2)                          │                   │
│               │ E-mail  [______________________]     │                   │
│               │ Senha   [______________________] [o] │                   │
│               │ [            Entrar            ]     │  Sinal            │
│               │ Esqueci a senha                      │                   │
│               │ ──────────── ou ────────────         │                   │
│               │ [Entrar com Google] [Entrar com Microsoft] │             │
│               └──────────────────────────────────────┘                   │
│                                                                          │
│                                          v1.5.0-rc.1 · Homologação      │
└──────────────────────────────────────────────────────────────────────────┘
```

- Fundo `fundo`, cartão `superficie` de 400 px centralizado. Uma ação primária (Entrar).
- O selo de versão (6.19) fica no canto inferior direito da página, fora do cartão. Em produção mostra só "v1.4.2"; fora da produção, a versão e o ambiente, como no exemplo.

---

## 6. Componentes

Estados obrigatórios em todo componente interativo: **padrão, hover, foco, pressionado, desabilitado, carregando, erro**. Os componentes são implementados sobre primitivas acessíveis sem estilo (Radix UI ou React Aria no web) e estilizados só com tokens.

### 6.1 Botões

| Variante | Fundo | Texto e ícone | Borda | Quando usar |
|---|---|---|---|---|
| Primário | `acao-primaria` (Sinal) | Aço | — | A ação principal da tela: "Abrir OS", "Encerrar OS", "Salvar plano". Uma por tela ou modal. |
| Secundário | transparente | `texto` | 1 px `borda-controle` | Ações alternativas: "Ver histórico", "Exportar", "Voltar" |
| Terciário | transparente | `texto`, sublinhado no hover | — | Ações de baixo peso em linha: "Limpar filtros", "Como calculamos" |
| Perigo | `perigo` | `perigo-texto` | — | Ação destrutiva confirmada: "Excluir item", "Cancelar OS" (dentro da confirmação) |
| Perigo secundário | transparente | `texto` com ícone em `perigo` | 1 px `borda-controle` | Ação destrutiva fora da confirmação, para não chamar atenção indevida |
| Só ícone | transparente | `texto-secundario` | — | Barras de ferramenta, ações de linha. Exige `aria-label` e tooltip. |

| Estado | Primário | Secundário | Perigo |
|---|---|---|---|
| Padrão | Sinal, texto Aço | Borda `#8592A0` | `#D33A2C`, texto branco |
| Hover | `#FF7A2E` | Borda Aço 500, fundo `superficie-hover` | `#B8301F` |
| Foco | Anel de 2 px em `foco`, afastado 2 px | Idem | Idem |
| Pressionado | `#F25F0A` | Fundo Aço 100 | `#B8301F` |
| Desabilitado | Fundo Aço 100, texto Concreto, sem sombra | Borda Aço 100, texto Concreto | Fundo Aço 100, texto Concreto |
| Carregando | Mantém largura, troca o ícone por indicador giratório, texto vira gerúndio ("Salvando"), `aria-busy="true"`, clique ignorado | Idem | Idem |
| Erro | O botão não fica vermelho. O erro aparece no campo ou em alerta, e o botão volta ao padrão. | — | — |

Tamanhos web: P 32 px (tabelas e barras de ferramenta), M 40 px (padrão), G 48 px (estados vazios, telas de toque). Padding horizontal 12/16/20 px, ícone 16/20/20 px, raio `raio-md`. No app, mínimo 48 dp e 56 dp para a ação primária, largura total na base da tela.

Regras:

- O texto do botão é um verbo com o objeto: "Encerrar OS", "Gerar requisição". Nunca "OK", "Sim", "Enviar" sozinhos.
- Ordem: no web, ações à direita, primária por último (mais à direita). No app, primária embaixo e com largura total, secundária acima dela.
- Desabilitar o botão de salvar enquanto o formulário tem erro esconde o motivo. Prefira manter habilitado e mostrar os erros ao clicar.

### 6.2 Campos de formulário

Anatomia: rótulo (`rotulo`, acima do campo, sempre visível) → campo → texto de apoio ou erro (`legenda`). Marcação "(opcional)" nos campos opcionais em vez de asterisco nos obrigatórios, porque a maioria é obrigatória.

| Estado | Visual |
|---|---|
| Padrão | Fundo `superficie`, borda 1 px `borda-controle`, raio `raio-md`, altura do controle |
| Hover | Borda `borda-controle-hover` |
| Foco | Borda 1 px `foco` + anel de 2 px `foco` afastado 2 px |
| Preenchido | Igual ao padrão; valor em `texto` |
| Desabilitado | Fundo `fundo`, borda `borda`, texto `texto-desabilitado`, `cursor: not-allowed` |
| Somente leitura | Sem borda, fundo transparente, valor selecionável. Usado em OS encerrada. |
| Carregando (validação remota) | Indicador de 16 px à direita, `aria-busy` |
| Erro | Borda 2 px `status-parado`, ícone `circle-x` à esquerda da mensagem, mensagem em `texto` (não em vermelho), `aria-invalid="true"` e `aria-describedby` ligado à mensagem |
| Aviso | Borda 2 px `status-atencao`, ícone `triangle-alert`, mensagem que não impede salvar ("Leitura 18 % acima da média dos últimos 30 dias. Confirme o valor.") |

Tipos específicos:

- **Medidor:** campo numérico com unidade fixa à direita dentro do campo (km, h), teclado numérico no app, última leitura mostrada abaixo ("Última: 248.315 km em 10/11/2026, pelo app"). Valor menor que a última leitura é erro; salto acima do esperado é aviso.
- **Moeda:** prefixo "R$" fixo, formatação ao sair do campo.
- **Data:** digitação dd/mm/aaaa com máscara e calendário opcional. Nunca só calendário.
- **TAG e placa:** JetBrains Mono, conversão automática para caixa alta.
- **Texto longo:** contador de caracteres só quando houver limite; ditado por voz no app.
- Validação ao sair do campo e ao salvar, nunca a cada tecla. Ao salvar com erros, o foco vai para o primeiro campo com erro e um resumo aparece no topo do formulário com links para cada campo.

### 6.3 Seletores

| Componente | Quando usar | Regras |
|---|---|---|
| Seleção simples (lista) | 5 a 15 opções | Opções na ordem lógica ou alfabética; opção atual com ícone `check` |
| Combobox com busca | Mais de 15 opções (ativos, peças, técnicos, fornecedores) | Busca por código e descrição; mostra os 5 recentes antes de digitar; resultado com TAG em `dado-forte` + descrição + status |
| Seleção múltipla | Filtros, atribuição a vários técnicos | Itens marcados aparecem como chips removíveis no campo; "3 selecionados" quando passar de 2 chips |
| Grupo de opções (radio) | 2 a 5 opções visíveis ao mesmo tempo | Preferir a seletor quando a escolha for frequente (motivo de pausa) |
| Segmentado | 2 a 4 modos de visualização (Lista, Quadro, Calendário) | Não usar para dados de formulário |
| Interruptor | Liga/desliga de efeito imediato (plano ativo) | Rótulo à esquerda descreve o estado ligado |
| Árvore | Localização e hierarquia de ativos | Expandir com seta direita, recolher com seta esquerda (padrão ARIA `tree`); mostra o caminho completo no valor escolhido ("Planta Contagem › Linha 2 › CMP-0347") |

Estados iguais aos de campo. A lista suspensa abre em `superficie-elevada`, `sombra-3`, altura máxima de 320 px, e fecha com Esc devolvendo o foco ao campo. Opção desabilitada mostra o motivo ("João Lima · NR-10 vencida").

### 6.4 Tabela de dados

É o componente mais usado pelo planejador. Implementação sugerida: TanStack Table com virtualização a partir de 200 linhas.

**Anatomia**

```
┌────────────────────────────────────────────────────────────────────────────┐
│ [Visão: Minhas OS abertas ▾]  [Buscar nesta lista]   [Filtros 3] [Colunas] │
│ Status: Em execução, Pausada ×   Unidade: Contagem ×   Limpar filtros      │
├──┬─────────────────┬──────────────────────┬────────────┬──────────┬───┬────┤
│[ ] │ OS ▲            │ Ativo                │ Status     │ Custo (R$)│...│ ⋯ │
├──┼─────────────────┼──────────────────────┼────────────┼──────────┼───┼────┤
│[ ] │ OS-2026-004812  │ FRT-0112 Cavalo 6x2  │ [|| Pausada]│  3.418,90 │   │ ⋯ │
│[x] │ OS-2026-004809  │ CMP-0347 Compressor  │ [> Em exec]│    812,00 │   │ ⋯ │
└──┴─────────────────┴──────────────────────┴────────────┴──────────┴───┴────┘
  1–50 de 1.284                         [25|50|100]   [‹] 1 2 3 … 26 [›]
```

**Regras**

- **Ordenação:** clique no cabeçalho alterna crescente, decrescente e sem ordem. Seta `arrow-up` ou `arrow-down` ao lado do rótulo, `aria-sort` no `th`. Shift + clique adiciona ordenação secundária. A ordenação padrão de cada lista é documentada (OS: prioridade, depois data prevista).
- **Filtros:** barra acima da tabela com busca local, botão "Filtros" com contador e chips dos filtros ativos, removíveis um a um. Filtros e ordenação ficam na URL, para compartilhar e voltar. **Visões salvas** por usuário e visões padrão da empresa ("Atrasadas", "Aguardando peça").
- **Seleção:** checkbox na primeira coluna; o do cabeçalho seleciona a página e mostra "50 selecionadas nesta página. Selecionar todas as 1.284." Shift + clique seleciona intervalo.
- **Ações em lote:** com seleção ativa, uma barra fixa substitui a barra de filtros: "12 OS selecionadas · Programar · Atribuir técnico · Alterar prioridade · Exportar · Mais". Ações destrutivas em lote pedem confirmação com o número de itens. Ações que não se aplicam a parte da seleção informam antes: "3 das 12 OS estão encerradas e não serão alteradas."
- **Colunas fixas:** a coluna de seleção e a do identificador ficam fixas à esquerda; a de ações, à direita. Ao rolar na horizontal, a borda da coluna fixa ganha `sombra-2`. Cabeçalho fixo na rolagem vertical.
- **Colunas configuráveis:** mostrar, ocultar, reordenar e redimensionar; preferência salva por usuário e por lista.
- **Números:** alinhados à direita, JetBrains Mono ou Inter com `tabular-nums`, mesmas casas decimais em toda a coluna, unidade no cabeçalho. Totais em linha de rodapé fixa com `dado-forte`.
- **Identificadores:** TAG, placa e número de OS em `dado-forte`, nunca truncados. Descrições longas truncam com reticências e mostram o texto completo em tooltip.
- **Status:** sempre o chip (seção 7), nunca só uma bolinha colorida.
- **Linha:** a linha inteira é um link para o detalhe (Ctrl + clique abre em nova aba). Hover em `superficie-hover`; selecionada em `superficie-selecionada` com barra de 3 px em Aço à esquerda.
- **Ações por linha:** menu `ellipsis` na última coluna, também aberto com a tecla de contexto ou Shift + F10.
- **Edição em linha:** só para campos simples (prioridade, data prevista, responsável), com Enter para salvar e Esc para descartar.
- **Estados:** carregando (5 linhas de esqueleto com as larguras das colunas), vazio (seção 6.14), erro (linha única com mensagem e botão "Tentar de novo"), atualizando (barra de progresso fina no topo, dados atuais continuam visíveis).
- **Teclado:** padrão ARIA `grid`: setas movem entre células, Enter abre, Espaço seleciona, Home e End vão ao início e fim da linha.
- **Exportar:** CSV e XLSX respeitando filtros e colunas visíveis, com números em formato numérico (não texto).

### 6.5 Chip de status

Anatomia: ícone de 14 px + texto `legenda` peso 600, padding 3 × 10 × 3 × 8 px, raio `raio-sm`, fundo `status-<x>-fundo`, ícone `status-<x>`, texto `status-<x>-texto`. Tamanho grande (app e cabeçalho de detalhe): ícone 16 px, texto 15 px, altura 28 px.

- O texto é a palavra do estado, não um código (`Em execução`, nunca `em_execucao`).
- Pode ter um complemento depois de um ponto médio: "Pausada · Aguardando peça", "Atenção · Preventiva em 300 km".
- Chip de status não é botão. Se o status puder ser alterado ali, o controle é um seletor com o chip dentro e a seta `chevron-down`.
- O mapeamento completo está na seção 7.

### 6.6 Badge de prioridade e criticidade

Badges não usam as cores de status, para que a tela não pareça ter dois status diferentes.

**Criticidade do ativo (A, B, C)** — quadrado de 20 px, raio `raio-xs`, letra em JetBrains Mono 600:

| Criticidade | Visual | Texto acessível |
|---|---|---|
| A | Fundo Aço, letra Névoa (14,93:1) | "Criticidade A" |
| B | Contorno 1,5 px Aço, letra Aço | "Criticidade B" |
| C | Contorno 1 px `borda-controle`, letra Aço 500 | "Criticidade C" |

**Prioridade da solicitação e da OS** — ícone de barras + palavra. Valores canônicos do sumário (seção 3.1): `baixa`, `media`, `alta`, `urgente`.

| Valor canônico | Rótulo | Ícone | Cor do ícone | Texto acessível |
|---|---|---|---|---|
| `urgente` | Urgente | `siren` | `status-parado` | "Prioridade urgente" |
| `alta` | Alta | `signal-high` | `texto` | "Prioridade alta" |
| `media` | Média | `signal-medium` | `texto` | "Prioridade média" |
| `baixa` | Baixa | `signal-low` | `texto-secundario` | "Prioridade baixa" |

- Só Urgente usa cor, porque significa ativo parado ou risco à segurança. Em tabelas, a coluna de prioridade pode ser ordenada (Urgente primeiro) e o ícone fica alinhado à esquerda.
- OS do tipo `corretiva_emergencial` nasce Urgente. Nas demais, a prioridade padrão vem da criticidade do ativo: A → Alta, B → Média, C → Baixa. O formulário mostra a origem: "Alta · padrão da criticidade A".
- Nunca usar "Crítica" ou "Emergência" como prioridade: "Crítica" confunde com criticidade do ativo e "Emergencial" é tipo de OS.

### 6.7 Cartão de ativo

Segue a aplicação do brandbook (seção 07, "Cavalo mecânico 6x2").

```
┌───────────────────────────────────────────────────┐
│ Cavalo mecânico 6x2                 [✓ Operando]  │  h3 Inter 700 + chip
│ FRT-0112 · Placa RTZ4F82          [A]             │  dado + criticidade
│                                                   │
│ HODÔMETRO            CUSTO/KM 90D                 │  sobrelinha
│ 248.315 km           R$ 0,41                      │  dado-forte
│ ÚLTIMA OS            BASE                         │
│ OS-2026-004812       Contagem/MG                  │
│ ┌───────────────────────────────────────────────┐ │
│ │ PRÓXIMA PREVENTIVA                            │ │
│ │ Troca de óleo e filtros em 2.300 km           │ │
│ │ (cerca de 9 dias)                             │ │
│ └───────────────────────────────────────────────┘ │
│ Disponibilidade 30 dias                  97,4 %   │
│ [██████████████████████████████████████████░]     │  barra status-operando
│ [Abrir OS]  [Ver histórico]                       │  primário + secundário
└───────────────────────────────────────────────────┘
```

- Fundo `superficie`, borda `borda`, raio `raio-lg`, padding 20 px, `sombra-1`.
- Pares rótulo/valor em grade de 2 colunas (4 no cartão largo). Rótulos em `sobrelinha` `texto-terciario`, valores em `dado-forte`.
- A barra de disponibilidade usa `status-operando` acima da meta, `status-atencao` (com contorno Aço) entre 90 % da meta e a meta, e `status-parado` abaixo; o valor numérico fica sempre ao lado.
- Em listas de cartões, "Abrir OS" vira secundário (regra de uma ação primária por tela).
- Variante compacta (lista do app): nome, TAG, chip e próxima preventiva, sem botões; o cartão inteiro é tocável.
- Ativos do tipo equipamento trocam hodômetro por horímetro e placa por localização.

### 6.8 Linha do tempo da OS

Lista vertical em ordem cronológica, mais recente no topo (alternável).

```
  ●  14:32  Ana Souza pausou a OS                       [|| Pausada]
  │         Motivo: aguardando peça · Filtro de óleo PSL-118
  │         Registrado no app às 14:32 · enviado às 16:05
  │
  ●  11:10  Ana Souza iniciou a execução                [> Em execução]
  │
  ○  10:45  Comentário de Carlos Lima
  │         "Verificar também a correia do alternador."
  │         [foto] [foto]
  │
  ●  09:02  Carlos Lima programou para 12/11/2026, Ana Souza
```

- Marcador de 24 px com o ícone do evento; mudança de estado leva o chip do novo estado à direita.
- Hora em `dado` (data completa como separador de grupo: "Hoje", "Ontem", "10/11/2026").
- Tipos de evento: mudança de estado, apontamento de horas, peça consumida, leitura de medidor, foto, comentário, assinatura, alteração de campo (com valor antes e depois, vindo da trilha de auditoria).
- Eventos registrados offline mostram o horário real de registro e o horário de envio.
- Filtro por tipo de evento acima da lista. Campo de comentário no topo, com menção a pessoas (`@`).
- A linha do tempo é somente leitura: um evento não é editado nem excluído; correções viram novo evento.

### 6.9 Leitor de QR

Usado no app (principal) e no web em tablets.

- Tela cheia com a câmera, moldura quadrada de 240 dp centralizada com cantos em Névoa, véu escuro fora da moldura.
- Texto acima da moldura: "Aponte para a etiqueta do ativo". Abaixo: botões de 56 dp "Lanterna" e "Digitar TAG".
- Leitura confirmada com vibração curta, som opcional e moldura que muda para `status-operando` por 300 ms; abre o ativo direto.
- Funciona offline: resolve o código pela base local do aparelho.
- Erros:
  - QR de outro sistema: "Este QR não é de um ativo da Tenaz. Leia a etiqueta de patrimônio ou digite a TAG."
  - Ativo de outra empresa ou sem acesso: "Este ativo não está na sua lista. Peça acesso ao gestor da unidade."
  - Ativo fora do aparelho, sem conexão: "Ativo não encontrado no aparelho. Conecte-se para buscar ou digite a TAG."
  - Câmera negada: "Permita o uso da câmera para ler QR. Abrir configurações" + "Digitar TAG".
- Leitura contínua (inventário, conferência de pneus): lista dos lidos na parte inferior, com contador.

### 6.10 Assinatura

- Área de desenho de largura total, altura de 200 dp (ou tela cheia em paisagem), fundo `superficie`, linha de base em `borda-controle` e texto "Assine com o dedo" em `texto-terciario` sobre a linha.
- Acima: nome (preenchido com o usuário logado ou digitado para terceiros) e função. Abaixo: "Limpar" (terciário) e "Confirmar assinatura" (primário, habilitado só com traço).
- Registra data, hora, usuário do aparelho e, se permitido, localização; tudo aparece no documento gerado.
- Alternativa acessível: "Confirmar sem desenhar" para quem não consegue desenhar, registrando aceite com identificação do usuário logado.
- A assinatura no app é eletrônica simples. Para OS de NR-12 e relatórios de NR-13, a exigência de certificado ICP-Brasil depende da pendência P1 do sumário executivo.

### 6.11 Upload de foto

- No app: botões "Tirar foto" (câmera do app) e "Galeria". No web: área de soltar arquivos com o texto "Arraste fotos ou PDFs aqui ou escolha do computador", aceitando JPG, PNG, HEIC e PDF de até 20 MB cada.
- Miniaturas de 88 dp em grade, com legenda opcional e botão remover (com "Desfazer" no toast por 8 s).
- Quando o item exige fotos: contador "2 de 3 fotos" e o item fica incompleto até atingir o mínimo.
- Compressão no aparelho para 1.600 px no maior lado, JPEG qualidade 80, mantendo data e hora da captura.
- Estado por foto: salva no aparelho (`smartphone`), enviando (barra de progresso na miniatura), enviada (`cloud-check`), falha (`triangle-alert` + "Tocar para reenviar"). Fotos pendentes não bloqueiam concluir a OS; vão na fila de sincronização.
- Erros: "Arquivo de 34 MB. O limite é 20 MB por arquivo. Reduza a foto ou envie em PDF." · "Formato .mov não aceito. Envie JPG, PNG, HEIC ou PDF."

### 6.12 Alertas e toasts

**Alerta em linha (banner)** — dentro da página, persiste até a condição mudar.

| Tipo | Ícone | Fundo | Barra lateral 4 px | Exemplo |
|---|---|---|---|---|
| Informação | `info` | `superficie` | Aço | "Esta OS veio do plano PRE-0042. Alterações valem só para esta OS." |
| Sucesso | `circle-check` | `status-operando-fundo` | `status-operando` | "Importação concluída: 312 ativos criados, 4 ignorados." |
| Atenção | `triangle-alert` | `status-atencao-fundo` | `status-atencao` | "CRLV-e do FRT-0112 vence em 12 dias." |
| Erro | `circle-x` | `status-parado-fundo` | `status-parado` | "Não foi possível consolidar os custos. Verifique as 2 peças sem valor." |

Estrutura: ícone + título opcional em `h4` + texto em `corpo-p` + até duas ações terciárias. Texto sempre em `texto`. Fecha com `x` só quando for informativo. `role="status"` para informação e sucesso, `role="alert"` para erro.

**Toast** — confirma uma ação que acabou de acontecer.

- Fundo `toast-fundo`, texto `toast-texto`, ícone na variante "sobre Aço" do status, raio `raio-md`, `sombra-4`, largura de 360 a 480 px.
- Posição: canto inferior esquerdo no web (não cobre a ação primária, que fica à direita); acima da barra de abas no app.
- Duração: 5 s; 8 s quando tem ação "Desfazer"; erro não some sozinho. Pausa no hover e no foco. No máximo 3 empilhados.
- Uma frase, sem título: "OS-2026-004812 encerrada. Ativo liberado para operação."
- `aria-live="polite"` para sucesso, `assertive` para erro.

### 6.13 Modais e confirmações

- Larguras: P 400 px (confirmação), M 560 px (formulário curto), G 720 px (formulário com seções). Acima disso, use página ou painel lateral.
- Painel lateral (480 a 640 px, vem da direita) para ver ou editar um item sem perder a lista.
- Foco preso dentro do modal; foco inicial no primeiro campo ou, em confirmação, no botão menos destrutivo. Esc e clique no véu fecham, exceto durante o envio ou com dados não salvos (aí pergunta "Descartar alterações?").
- No app, modais viram folha inferior com alça e cantos `raio-xl`, ou tela cheia para formulários.

**Confirmação** — só para ações irreversíveis ou de grande alcance. Ações reversíveis usam "Desfazer" no toast em vez de perguntar.

```
┌───────────────────────────────────────────────┐
│ Encerrar OS-2026-004812?                      │  h3
│                                               │
│ Os custos serão consolidados no centro de     │
│ custo do FRT-0112 e a OS não poderá mais ser  │
│ editada.                                      │
│                                               │
│                        [Voltar] [Encerrar OS] │
└───────────────────────────────────────────────┘
```

- Título é a pergunta com o objeto. O texto diz a consequência. O botão repete o verbo.
- O botão que fecha sem agir é "Voltar", nunca "Cancelar", porque "Cancelar" é também uma ação sobre a OS ("Cancelar OS?" com botões "Cancelar" e "Cancelar OS" confunde).
- Ação destrutiva usa o botão de perigo; ação importante mas não destrutiva (encerrar) usa o primário.
- Exclusão em lote acima de 20 itens pede digitar o número de itens para confirmar.

### 6.14 Estados vazios

| Tipo | Estrutura | Exemplo |
|---|---|---|
| Primeiro uso | Ícone de 48 px em `texto-terciario`, título `h3`, uma frase, ação primária e link de ajuda | "Nenhum ativo cadastrado. Cadastre o primeiro ativo ou importe uma planilha." [Cadastrar ativo] [Importar planilha] |
| Sem resultado de filtro | Ícone `search`, frase com o filtro, ação terciária | "Nenhuma OS com status Pausada em Contagem. Limpar filtros" |
| Tudo em dia | Ícone `circle-check` em `status-operando`, frase calma, sem ação | "Nenhuma preventiva atrasada." |
| Sem permissão | Ícone `lock`, frase com o motivo e quem procurar | "Seu perfil não acessa o estoque. Fale com o administrador da empresa." |
| Erro de carga | Ícone `cloud-off` ou `circle-x`, frase e "Tentar de novo" | "Não foi possível carregar as OS. Verifique a conexão e tente de novo." |

Sem ilustrações coloridas. Estados vazios usam ícones da biblioteca no estilo da marca.

### 6.15 Esqueleto de carregamento

- Blocos em `superficie-esqueleto` com o formato do conteúdo final (linhas de tabela, cartões, KPI). Raio igual ao do componente.
- Aparece só se a carga passar de 300 ms, para não piscar; uma vez visível, fica pelo menos 300 ms.
- Brilho em movimento de 1,6 s em loop; com movimento reduzido, bloco estático.
- Acima de 10 s, mostra texto: "Ainda carregando. Listas com mais de 10 mil itens podem levar alguns segundos."
- O contêiner recebe `aria-busy="true"` e um texto oculto "Carregando ordens de serviço".
- Atualização de dados já exibidos não volta ao esqueleto: usa barra de progresso fina no topo da área.

### 6.16 Paginação

```
1–50 de 1.284 OS         Itens por página [50 ▾]      [‹ Anterior] 1 2 3 … 26 [Próxima ›]
```

- Contagem em `dado` à esquerda; tamanho da página (25, 50, 100) no centro; navegação à direita.
- Página atual com fundo Aço e texto Névoa, `aria-current="page"`. Botões de 32 px no web e 48 dp no app.
- Página e tamanho ficam na URL. Mudar filtro volta para a página 1.
- No app, listas usam rolagem infinita com carga incremental e indicador ao fim ("Fim da lista · 128 OS").
- Acima de 10 mil resultados, a contagem é aproximada ("mais de 10 mil") e sugere filtrar.

### 6.17 Gráficos de KPI

**Cartão de KPI**

```
┌──────────────────────────────────────┐
│ DISPONIBILIDADE          ⓘ Como calculamos
│ 97,4 %                               │  kpi
│ ↑ 1,2 p.p. vs. outubro · Meta 96 %   │  legenda + ícone
│ ╭──────────────────────────────╮     │
│ │    ___/‾‾\__/‾‾‾‾‾\___/‾‾‾●  │     │  sparkline 30 dias
│ ╰──────────────────────────────╯     │
└──────────────────────────────────────┘
```

- Rótulo em `sobrelinha`, valor em `kpi` com unidade, variação com ícone de seta e texto (`trending-up`, `trending-down`). A cor da variação só aparece quando o sentido bom é inequívoco (disponibilidade subindo: `status-operando`; MTTR subindo: `status-parado`); o texto sempre diz o sentido.
- Período sempre explícito no cartão ou no filtro do painel ("Últimos 30 dias", "Novembro de 2026").

**Gráfico de linha e área**

- Série principal em `serie-1` (Aço no claro, `#E8EDF2` no escuro), linha de 2 px, área com gradiente da cor da série de 12 % a 0 % de opacidade.
- **Grade suave:** só linhas horizontais, 1 px em `grade-grafico`, 4 a 5 divisões com valores redondos. Sem grade vertical, sem borda no gráfico.
- Eixos: rótulos em `dado` 12 px, `texto-secundario`; eixo Y com unidade no primeiro rótulo ou no título ("h", "%").
- **Ponto final destacado:** círculo de 8 px na cor da série, halo de 2 px na cor da superfície, valor ao lado em `dado-forte`. O ponto final não usa Sinal.
- Meta: linha tracejada (4 × 4) em `texto-secundario` com rótulo "Meta 96 %" na ponta direita.
- Até 3 séries, em `serie-1`, `serie-2` e `serie-3`, diferenciadas também por tracejado (contínua, tracejada, pontilhada) e com legenda direta na ponta de cada linha. Mais de 3 séries vira pequenos múltiplos.
- Cores de status só quando a categoria é um status (ex.: horas por status do ativo, barras empilhadas Operando/Em manutenção/Parado). Nesses casos, barras em Atenção recebem contorno de 1 px em Aço no tema claro.
- Sinal nunca aparece em gráficos.
- Tooltip no hover e no foco do teclado com data, valor e unidade; linha vertical de referência em `borda-controle`.
- Cada gráfico tem "Ver como tabela" (alternativa acessível e de conferência) e exportação de dados.

**Como calculamos**

Botão terciário com ícone `info` no canto do cartão. Abre um popover (`sombra-3`, até 400 px) com:

1. Fórmula em linguagem simples e em notação: "Disponibilidade = (tempo do período − tempo parado) ÷ tempo do período".
2. O que entra e o que não entra: "Conta paradas por corretiva e preventiva. Não conta ativos inativos nem horas fora do turno cadastrado."
3. Fonte dos dados: "Estados do ativo registrados nas OS."
4. Período e filtros aplicados, e horário da última atualização: "Calculado às 06:00 de 12/11/2026."
5. Link "Ver os registros" que abre a lista filtrada que gerou o número.

Indicadores com fórmula padronizada no produto: disponibilidade, MTBF, MTTR, preventivas no prazo, backlog em horas, custo de manutenção por ativo, custo por km, consumo (km/L), percentual de corretiva emergencial. As fórmulas exatas ficam em [1_ESPECIFICACAO_REQUISITOS.md](1_ESPECIFICACAO_REQUISITOS.md) (IND).

### 6.18 Outros componentes de apoio

| Componente | Regras principais |
|---|---|
| Abas | Sublinhado de 2 px em Aço na aba ativa; contador entre parênteses; setas navegam entre abas (padrão ARIA `tablist`) |
| Tooltip | Só para complementar; nunca para informação essencial. Abre no hover após 400 ms e no foco. Fundo Aço, texto Névoa, `legenda`. |
| Checkbox e radio | 20 px no web com área clicável de 32 px; marcado em Aço com símbolo Névoa |
| Barra de progresso | 8 px, raio `raio-xs`, trilha `borda`, preenchimento Aço (progresso) ou cor de status (indicador). Valor em texto ao lado. |
| Avatar | Iniciais em Inter 600 sobre Aço 100; 24, 32 e 40 px |
| Etiqueta de unidade | Texto `legenda` com ícone `map-pin` quando a visão é "Todas as unidades" |

### 6.19 Selo de versão

Mostra a versão do release em uso (SemVer 2.0.0, requisito RNF-039 em [1_ESPECIFICACAO_REQUISITOS.md](1_ESPECIFICACAO_REQUISITOS.md); o histórico de versões é o RNF-040). Serve para o usuário dizer ao suporte qual versão usa e para a equipe distinguir produção de homologação num relance.

**Onde aparece:** sempre no canto inferior direito da tela principal do web (todas as páginas do layout, 5.1, e dos portais 8.17 e 8.19), da tela de login do web (5.6), da tela principal do app (Hoje, 9.1) e da tela de login do app. Também no rodapé da página do link de uso único (8.18), sem ser fixo.

**Texto:**

| Ambiente | Exibição | Exemplo |
|---|---|---|
| Produção | `v` + versão | v1.4.2 |
| Fora da produção | `v` + versão de pré-lançamento + ponto médio + ambiente | v1.5.0-rc.1 · Homologação |

O metadado de build (`+a1b2c3d`) não aparece no selo; fica no diálogo Sobre. A versão vem do build (tag git), nunca de texto digitado.

**Visual:**

| Propriedade | Web | App |
|---|---|---|
| Fonte | JetBrains Mono 400, 12 px / 16 px, algarismos tabulares | JetBrains Mono 400, 12 dp / 16 dp; acompanha a escala de fonte do sistema |
| Cor do texto | `texto-terciario`: Aço 400 `#5F6C7B` no claro (4,86:1 em Névoa, 5,36:1 em branco), `#8592A0` no escuro (5,9:1 em `#0E151E`), Aço 500 `#4A5868` no tema Sol (7,27:1 em branco). Todos AA | Igual |
| Fundo | `fundo` da tela, opaco, com padding 2 × 6 px e raio `raio-xs`, para manter o contraste se algo rolar por baixo | Igual |
| Posição | `position: fixed`, 8 px da borda direita e 8 px da inferior da área de conteúdo, camada `z-fixo` (abaixo de menus, modais e toasts) | 8 dp da borda direita, logo acima da barra de abas e dentro da área segura (`useSafeAreaInsets`); no login, acima da área segura inferior |
| Área de toque | Altura mínima de 24 px | 48 × 48 dp, com o texto alinhado à direita e embaixo dentro da área |

**Regras de layout:**

- Nunca cobre conteúdo nem botões. No web, a área de conteúdo reserva 32 px de respiro inferior, e barras de ação fixas no rodapé (barra de lote da tabela, rodapé de formulário) terminam antes do selo ou o deslocam para cima delas. No app, o selo só aparece em telas sem botão fixo inferior (Hoje e login); nas telas de execução, que têm ação principal fixa embaixo, a versão fica em Mais > Sobre.
- Toasts e o indicador de sincronização do app não sobrepõem o selo: a pilha de toasts começa acima dele.
- Não usa Sinal, cor de status nem ícone. Não pisca nem anima quando a versão muda.

**Interação e acessibilidade:**

- É um botão (`<button>` no web, `Pressable` com `accessibilityRole="button"` no app). Nome acessível: "Versão 1.4.2. Abrir informações sobre o sistema". Fora da produção: "Versão 1.5.0-rc.1, ambiente Homologação. Abrir informações sobre o sistema".
- Foco visível padrão (anel de 2 px em `foco`, 11.1). Alcançável por Tab, por último na ordem da página, e pelo leitor de tela no app.
- Clique, toque, Enter ou Espaço abrem o diálogo **Sobre**:

```
┌───────────────────────────────────────────────┐
│ Sobre a Tenaz                              [×]│  h3
│                                               │
│ WEB             v1.4.2                        │  sobrelinha + dado
│ API             v1.4.2                        │
│ DATA DO BUILD   12/11/2026 14:05              │
│ COMMIT          a1b2c3d                       │
│ AMBIENTE        Produção                      │
│                                               │
│ Novidades desta versão  →                     │  link
│ Histórico de versões    →                     │  link
│                                               │
│ [Copiar informações]          [Fechar]        │
└───────────────────────────────────────────────┘
```

- **Novidades desta versão** (RNF-040) abre a seção do `CHANGELOG.md` da versão em uso, agrupada em Adicionado, Alterado, Corrigido, Removido, Obsoleto e Segurança, em linguagem de usuário. **Histórico de versões** abre a lista de todas as versões, da mais recente para a mais antiga. Na primeira vez que o usuário entra depois de uma versão MINOR ou MAJOR, um toast discreto oferece "Versão 1.5.0 instalada. Ver novidades".

- No app, a primeira linha é APP com a versão do release que está rodando ("v1.4.3", a mesma do selo), seguida de BINÁRIO DA LOJA com versão e build ("1.4.0 (140)") e ATUALIZAÇÃO com o identificador curto da atualização OTA aplicada, quando houver. As duas versões diferem quando o release chegou por atualização OTA (5_ARQUITETURA_SOFTWARE.md, seção 9.5).
- A versão da API vem de `GET /v1/versao`. Sem conexão, a linha mostra "Sem conexão · última vista v1.4.2".
- Se a versão do web for diferente da versão da API (navegador com versão antiga após um deploy), o diálogo mostra a faixa informativa "Há uma versão mais nova. Recarregue a página para atualizar." com o botão [Recarregar].
- "Copiar informações" copia as linhas em texto simples, para colar no chamado de suporte; o toast confirma "Informações da versão copiadas."

---

## 7. Status, prioridade e criticidade

### 7.1 Semântica das cores de status

As cinco cores da marca têm um significado fixo, que vale para qualquer entidade:

| Cor (marca) | Significado geral | Leitura |
|---|---|---|
| Operando (verde) | Disponível, concluído, aprovado | Nada a fazer |
| Preventiva (azul) | Trabalho em andamento ou parada programada | Acompanhar |
| Atenção (âmbar) | Exige ação em breve, interrompido, vencendo | Agir logo |
| Parado (vermelho) | Indisponível, vencido, bloqueado, urgente | Agir agora |
| Inativo (cinza) | Fora de operação, finalizado, cancelado | Histórico |
| Neutro (Aço, sem cor de status) | Etapas de preparação | Sem urgência |

Estados de preparação (aberta, planejada, programada, nova) usam o chip neutro: fundo `superficie-hover`, borda `borda`, ícone e texto em `texto`. Assim a cor fica reservada para o que pede atenção.

Chip sólido (tema sol e etiquetas impressas): Operando, Atenção e Inativo com texto Aço; Preventiva e Parado com texto branco, conforme a seção 4.

### 7.2 Status operacional do ativo

O sumário (seção 3.3) define quatro estados (`operando`, `em_manutencao`, `parado`, `inativo`) e a condição calculada Atenção; a marca define cinco cores. O mapeamento:

| Estado canônico | Rótulo na interface | Cor | Ícone (marca / Lucide) | Quando (sumário 3.3) |
|---|---|---|---|---|
| `operando` | Operando | Operando | `check` | Nenhuma das causas abaixo |
| `operando` + Atenção (calculada) | Atenção | Atenção | `triangle-alert` | Operando, mas com plano vencendo, documento vencendo ou leitura fora da faixa. Não é estado gravado. O complemento diz o motivo: "Atenção · Preventiva em 300 km". |
| `em_manutencao` | Em manutenção | Preventiva | `calendar` (ícone do brandbook) | Há OS `em_execucao` ou `pausada` e o ativo foi retirado de operação de forma programada (preventiva, preditiva, inspeção, corretiva planejada) |
| `parado` | Parado | Parado | `circle-x` | Falha que parou o ativo ("O ativo parou?" = Sim), corretiva emergencial em curso, parecer NR-12 `inapto` ou bloqueio de conformidade. O complemento diz a causa: "Parado · Bloqueio NR-12" |
| `inativo` | Inativo | Inativo | `circle-minus` | Baixado, vendido ou em reserva |

Regras de exibição:

- **Parado prevalece** sobre Em manutenção: se as duas causas existem, o chip é Parado.
- O status volta a Operando no **encerramento** da OS, se não restar outra causa. Enquanto a OS está `concluida` aguardando liberação, o chip continua no estado anterior e o detalhe mostra "Trabalho concluído · aguardando liberação".
- "Em manutenção" usa a cor e o ícone de Preventiva com o rótulo canônico, porque "Preventiva" também é tipo de OS e o ativo pode estar em manutenção por uma corretiva planejada.
- **Apto para despacho** (veículo) não é status: é uma verificação separada, mostrada como linha própria no detalhe e no checklist ("Apto para despacho" com `check` ou "Impedido · CRLV-e vencido" com `octagon-x`). Um veículo Operando pode estar impedido.
- **Locação** não muda o status. Ativo cedido e locado continua Operando; o complemento e a localização mostram o local de uso: "Operando · Locado a Construtora Alfa". Ativo tomado em locação (`propriedade = de_terceiro`) leva o selo neutro "De terceiro" ao lado do nome, com tooltip "Tomado em locação de Locadora Beta · CL-2026-000012".

### 7.3 Estados da ordem de serviço

| Estado | Rótulo | Cor | Ícone (Lucide) | Complemento no chip |
|---|---|---|---|---|
| `aberta` | Aberta | Neutro | `circle-dashed` | — |
| `planejada` | Planejada | Neutro | `clipboard-list` | — |
| `programada` | Programada | Neutro | `calendar-clock` | Data: "Programada · 14/11" |
| `em_execucao` | Em execução | Preventiva (azul) | `wrench` | Técnico: "Em execução · Ana" |
| `pausada` | Pausada | Atenção | `circle-pause` | Motivo obrigatório: "Pausada · Aguardando peça" |
| `concluida` | Concluída | Operando | `circle-check` | "Concluída · aguardando liberação" no detalhe |
| `encerrada` | Encerrada | Inativo | `lock` | — |
| `cancelada` | Cancelada | Inativo | `ban` | Motivo no detalhe e no tooltip |

Rótulos dos motivos de pausa: Aguardando peça, Aguardando aprovação, Aguardando terceiro, Fim do turno, Outro.

**Transições e o controle que as aciona** (diagrama do sumário 3.3). O quadro kanban (8.4) e o menu "Mover para" só oferecem estas; a coluna de destino inválida fica esmaecida.

| De → para | Ação na interface | Observação |
|---|---|---|
| (nova) → `aberta` | Criar OS | |
| (gerada por plano) → `planejada` | — | OS de plano nasce planejada; a linha do tempo mostra "Gerada pelo plano PRE-0042" |
| `aberta` → `planejada` | Planejar | |
| `aberta` → `programada` | Programar | Programar sem planejar |
| `aberta` → `em_execucao` | Iniciar | Só corretiva emergencial; o botão aparece direto na OS urgente |
| `planejada` → `programada` | Programar | |
| `programada` → `planejada` | Desprogramar | Item do menu "Mais"; remove data e responsável |
| `programada` → `em_execucao` | Iniciar | |
| `em_execucao` ⇄ `pausada` | Pausar / Retomar | Pausa pede motivo |
| `em_execucao` → `concluida` | Concluir | |
| `concluida` → `em_execucao` | Recusar liberação | Botão secundário ao lado de "Encerrar"; motivo obrigatório (`motivo_recusa_liberacao`), que aparece na linha do tempo e notifica o executor |
| `concluida` → `encerrada` | Encerrar | Exige a permissão `os.encerrar` (padrão em gestor e planejador; nunca em prestador nem locatário). Em OS vinculada a uma subcontratada conectada, só a conta do cliente encerra (8.28) |
| `aberta`, `planejada`, `programada`, `pausada` → `cancelada` | Cancelar OS | Motivo obrigatório. Em `em_execucao` o item "Cancelar OS" aparece desabilitado com a dica "Pause a OS antes de cancelar"; em `concluida` e `encerrada` não aparece |

**Reconhecimento** não é estado. É o campo `data_reconhecimento`, gravado uma vez quando alguém assume a OS (botão "Reconhecer", ou implicitamente ao programar, iniciar ou assumir). Na lista de OS, uma OS `aberta` ainda não reconhecida leva o indicador "Não reconhecida há 40 min" (ícone `bell-ring` em `texto`, ou em `status-atencao` quando passar do prazo configurado para a prioridade). Depois, a linha do tempo mostra "Reconhecida por Ana Souza às 09:02". Os marcos de tempo exibidos no detalhe usam os rótulos: Aberta em (`data_abertura`), Reconhecida em (`data_reconhecimento`), Iniciada em (`data_inicio`), Concluída em (`data_conclusao`), Encerrada em (`data_encerramento`), Cancelada em (`data_cancelamento`).

**Assumir** também não é estado. O técnico se atribui uma OS `aberta` ou `planejada` da sua equipe pelo botão "Assumir" (secundário na lista, no detalhe e nos cartões da seção "Disponíveis para assumir" do app, 9.5 e 9.8). Assumir torna o técnico responsável e grava `data_reconhecimento` se ainda estiver vazia. A linha do tempo mostra "Assumida por Ana Souza às 09:02". Se outro técnico assumir antes, o app mostra "Esta OS já foi assumida por João Lima." e a retira da lista.

**Executor e responsável pelo custo** são atributos, não status: texto com ícone em `texto-secundario`.

| Campo | Valores e rótulos | Ícone |
|---|---|---|
| Executor (`executor_tipo`) | Interno · "Fornecedor: Refrigeração Polar" · em OS vinculada, o selo "Executada por Refrigeração Polar" (8.28) | `users` / `building-2` / `cable` |
| Responsável pelo custo (`responsavel_custo`) | Empresa · Locatário · Locador · Garantia · Fornecedor | `wallet` |

Quando o responsável pelo custo foi alterado em relação ao padrão da matriz de locação, aparece o complemento "alterado por Carlos Lima" com a justificativa (`justificativa_responsavel_custo`) no tooltip.

**Origem da OS** (`origem`), mostrada no painel lateral do detalhe como ORIGEM, em `texto-secundario`:

| Valor | Rótulo |
|---|---|
| `manual` | Aberta por Carlos Lima |
| `solicitacao` | Solicitação SS-2026-000731 (link) |
| `plano` | Plano PRE-0042 (link) |
| `checklist` | Checklist reprovado (link para a execução) |
| `telematica` | Evento de telemática (link para o evento) |
| `rascunho_ia` | Rascunho do assistente, confirmado por Carlos Lima |
| `vistoria_locacao` | Vistoria de devolução do CL-2026-000012 (link) |
| `prestador` | Aberta por Refrigeração Polar (subcontratada conectada), com o número da OS na conta dela |

**Indicador de prazo** (derivado, aparece ao lado do chip de estado, nunca no lugar dele):

| Condição | Visual |
|---|---|
| Vence hoje | Ícone `clock` em `status-atencao` + "Vence hoje" |
| Atrasada | Ícone `clock-alert` em `status-parado` + "Atrasada 3 dias" |
| No prazo | Nada (sem ruído) |

### 7.4 Tipo de OS

Tipo é atributo, não status: aparece como texto com ícone em `texto-secundario`, sem fundo colorido.

| Tipo | Rótulo | Ícone |
|---|---|---|
| `preventiva` | Preventiva | `calendar-check` (ícone de Preventiva do brandbook) |
| `preditiva` | Preditiva | `activity` |
| `corretiva_planejada` | Corretiva planejada | `wrench` |
| `corretiva_emergencial` | Corretiva emergencial | `siren` |
| `inspecao_legal` | Inspeção legal | `shield-check` |
| `melhoria` | Melhoria | `trending-up` |

### 7.5 Solicitação

| Estado | Rótulo | Cor | Ícone |
|---|---|---|---|
| `nova` | Nova | Neutro | `inbox` |
| `em_triagem` | Em triagem | Preventiva | `search` |
| `aprovada` | Aprovada | Operando | `circle-check` (com link para a OS gerada) |
| `recusada` | Recusada | Inativo | `circle-x` |
| `duplicada` | Duplicada | Inativo | `copy` (com link para a original) |

`duplicada` pode ser marcada a partir de Nova ou Em triagem e é final: se a marcação estiver errada, a ação "Reabrir como nova" cria outra solicitação ligada à primeira, sem desfazer a marca. Desistência do solicitante pode ir de Nova direto para Recusada e aparece como "Recusada · Desistência". Solicitação aberta no portal do locatário mostra a origem "Portal do locatário" (`portal_locatario`).

**Ordem de reparo** (EST, Fase 2):

| Estado | Rótulo | Cor | Ícone |
|---|---|---|---|
| `aberta` | Aberta | Neutro | `circle-dashed` |
| `enviada` | Enviada | Preventiva | `send` |
| `em_orcamento` | Em orçamento | Atenção | `receipt-text` |
| `aprovada` | Aprovada · em conserto | Preventiva | `hammer` |
| `retornada` | Retornada | Operando | `circle-check` |
| `irreparavel` | Irreparável | Inativo | `package-x` |
| `cancelada` | Cancelada | Inativo | `ban` |

Em orçamento usa Atenção porque espera uma decisão da empresa. Prazo previsto vencido usa o indicador de prazo (7.3).

**Contrato de locação** (LOC, Fase 2):

| Estado | Rótulo | Cor | Ícone |
|---|---|---|---|
| `rascunho` | Rascunho | Neutro | `file-pen` |
| `vigente` | Vigente | Operando | `circle-check` |
| `encerrado` | Encerrado | Inativo | `lock` |
| `cancelado` | Cancelado | Inativo | `ban` |

Contrato vigente com fim em até 15 dias recebe o indicador "Termina em 9 dias" ao lado do chip, como o indicador de prazo da OS. Só o Rascunho pode ser cancelado; o Vigente termina pela ação "Encerrar contrato", que pede motivo quando o encerramento é antecipado.

**Conexão entre contas** (SRV):

| Estado | Rótulo | Cor | Ícone |
|---|---|---|---|
| `convidada` | Convite enviado (no cliente) · Convite recebido (na subcontratada) | Neutro | `mail` |
| `ativa` | Conectada | Operando | `cable` |
| `suspensa` | Suspensa | Atenção | `circle-pause` |
| `encerrada` | Encerrada | Inativo | `lock` |

**Sincronização da OS vinculada** (`vinculo_os.estado_sincronizacao`), mostrada como texto pequeno com ícone ao lado do selo da OS vinculada, nunca como chip principal:

| Estado | Visual | Texto |
|---|---|---|
| `sincronizado` | `cloud-check` em `texto-secundario` | "Sincronizada às 14:32" |
| `pendente` | `refresh-cw` em `texto-secundario` | "Enviando para Distribuidora Norte" |
| `erro` | `triangle-alert` em `status-atencao` | "Não sincronizada. Ver motivo" (link para a sincronização da conexão, 8.27) |

### 7.6 Vencimentos (documentos, capacitações, NR)

| Situação | Cor | Ícone | Texto |
|---|---|---|---|
| Válido, mais de 30 dias | Operando | `check` | "Válido até 10/03/2027" |
| Vence em até 30 dias (antecedência configurável) | Atenção | `triangle-alert` | "Vence em 12 dias" |
| Vencido | Parado | `circle-x` | "Vencido há 3 dias" |
| Bloqueante vencido (impede atribuir OS ou liberar o veículo) | Parado | `octagon-x` | "Vencido · bloqueia atribuição" |
| Sem documento | Neutro | `file-question` | "Não cadastrado" |

---

## 8. Telas principais

Wireframes textuais: blocos em ordem de leitura e hierarquia. Os módulos de fase 2 (pneus, combustível) seguem a mesma estrutura.

### 8.1 Painel do gestor

```
┌──────────────────────────────────────────────────────────────────────────┐
│ Painel de manutenção (h1)       [Últimos 30 dias ▾] [Todas as unidades ▾]│
├──────────────────────────────────────────────────────────────────────────┤
│ [KPI Disponibilidade] [KPI Preventivas no prazo] [KPI MTTR] [KPI Backlog]│
├───────────────────────────────────────────┬──────────────────────────────┤
│ ATIVOS AGORA                              │ PRECISA DE VOCÊ              │
│ ███████████████████▒▒▒▓▓░                 │ 4 OS concluídas aguardando   │
│ Operando 412 · Em manutenção 18 ·         │   liberação          [Ver]   │
│ Atenção 23 · Parado 7 · Inativo 12        │ 3 documentos vencem em 7 d   │
│ (barra empilhada + legenda com números)   │ 2 solicitações há mais de 24h│
├───────────────────────────────────────────┤ 5 peças abaixo do mínimo     │
│ PARADOS AGORA (tabela compacta)           │   [Gerar requisição]         │
│ TAG · Ativo · Parado há · OS · Responsável├──────────────────────────────┤
│ CMP-0347 Compressor  6 h  OS-2026-004809  │ PREVENTIVAS DA SEMANA        │
│ ...                                       │ seg ▓▓▓ ter ▓▓ qua ▓▓▓▓ ...  │
├───────────────────────────────────────────┴──────────────────────────────┤
│ [Gráfico: disponibilidade diária x meta]  [Gráfico: custo por mês, R$]   │
└──────────────────────────────────────────────────────────────────────────┘
```

Hierarquia: indicadores → o que está parado → o que depende do gestor → tendências. Cada número é um link para a lista filtrada. O painel não tem ação primária em Sinal; as ações ficam nos itens de "Precisa de você".

### 8.2 Lista de ativos

```
Ativos (h2)                                   [Importar] [Novo ativo]
[Visão: Todos ▾] [Buscar TAG, nome, placa] [Tipo ▾] [Status ▾] [Criticidade ▾] [Localização ▾]
[Lista | Árvore | Mapa]
┌──┬──────────┬─────────────────────┬──────────┬────┬──────────────┬──────────────┬────────────┐
│[ ] │ TAG      │ Ativo               │ Status   │Crit│ Localização  │ Próx. prev.  │ Medidor    │
│[ ] │ FRT-0112 │ Cavalo mecânico 6x2 │ Operando │ A  │ Base Contagem│ em 2.300 km  │ 248.315 km │
└──┴──────────┴─────────────────────┴──────────┴────┴──────────────┴──────────────┴────────────┘
```

A visão Árvore mostra localizações e ativos pai e filho com expansão; a contagem de filhos parados sobe para o pai ("2 componentes parados").

### 8.3 Detalhe do ativo

```
Ativos › FRT-0112
┌──────────────────────────────────────────────────────────────────────────┐
│ Cavalo mecânico 6x2 (h2)  [✓ Operando] [A]          [Mais ▾] [Abrir OS] │
│ FRT-0112 · Placa RTZ4F82 · Base Contagem/MG · Classe Cavalo 6x2         │
├──────────────────────────────────────────────────────────────────────────┤
│ [Resumo] [OS (34)] [Planos (6)] [Medidores] [Checklists] [Documentos (5)]│
│ [Peças] [Custos] [Pneus] [Filhos (3)] [Histórico]                        │
├───────────────────────────────────────────────┬──────────────────────────┤
│ RESUMO                                        │ PRÓXIMOS VENCIMENTOS     │
│ Hodômetro 248.315 km  ·  Disponib. 97,4 %     │ Preventiva em 2.300 km   │
│ Custo/km 90d R$ 0,41  ·  MTBF 412 h           │ CRLV-e vence 12/12       │
│ OS ABERTAS (lista curta com chips)            │ Tacógrafo vence 03/2027  │
│ ÚLTIMAS LEITURAS (mini gráfico)               │ QR DO ATIVO [imprimir]   │
│ DADOS CADASTRAIS (campos em 2 colunas)        │ DOCUMENTOS TÉCNICOS      │
└───────────────────────────────────────────────┴──────────────────────────┘
```

Abas mudam com o tipo: equipamento mostra NR-12 e horímetro; veículo mostra pneus, abastecimentos e multas. Abas sem conteúdo ficam ocultas, não vazias.

### 8.4 Lista de OS

Igual à estrutura de tabela (6.4), com segmentado [Lista | Quadro]. Colunas padrão: OS, Ativo (TAG + nome), Tipo, Status, Prioridade, Responsável, Executor, Data prevista, Prazo, Custo (R$). Colunas opcionais: Responsável pelo custo, Reconhecida em. Visões padrão: "Minhas OS", "Abertas sem planejamento", "Não reconhecidas", "Atrasadas", "Aguardando peça", "Com fornecedor", "Concluídas aguardando liberação". O Quadro é kanban por estado, com cartões arrastáveis só entre as transições da tabela em 7.3; a coluna de destino inválida fica esmaecida durante o arraste e o mesmo movimento está disponível pelo menu "Mover para". Arrastar para Cancelada ou para Em execução a partir de Concluída abre o pedido de motivo antes de gravar.

### 8.5 Detalhe da OS

```
Ordens de serviço › OS-2026-004812
┌──────────────────────────────────────────────────────────────────────────┐
│ OS-2026-004812 · Troca de óleo e filtros (h2)                            │
│ [|| Pausada · Aguardando peça] [Preventiva] [Alta] [Atrasada 1 dia]      │
│                                  [Mais ▾] [Reprogramar] [Retomar OS]    │
├──────────────────────────────────────────────────────────────────────────┤
│ FLUXO: Aberta ─ Planejada ─ Programada ─ [Em execução ⇄ Pausada] ─      │
│        Concluída ─ Encerrada        (etapa atual destacada)             │
├───────────────────────────────────────────────┬──────────────────────────┤
│ [Tarefas 3/5] [Peças] [Horas] [Falha] [Fotos] │ ATIVO                    │
│                                               │ FRT-0112 Cavalo 6x2      │
│ TAREFAS                                       │ [Operando → Em manut.]   │
│ [x] Drenar óleo do motor            0,5 h  Ana  │ RESPONSÁVEL Ana Souza    │
│ [x] Trocar filtro de óleo           0,3 h  Ana  │ EXECUTOR Interno         │
│ [ ] Trocar filtro de combustível                │ CUSTO DE Empresa         │
│                                               │ PROGRAMADA 12/11 08:00   │
│                                               │ RECONHECIDA 11/11 09:02  │
│                                               │ ESTIMADO 2,5 h · REAL 1,8│
│                                               │ CUSTO ATÉ AGORA R$ 812,00│
│ LINHA DO TEMPO (6.8)                          │ ORIGEM Plano PRE-0042    │
└───────────────────────────────────────────────┴──────────────────────────┘
```

A ação primária muda com o estado: Planejar → Programar → (técnico) Iniciar → Concluir → (quem tem a permissão `os.encerrar`) Encerrar. Em Concluída, "Recusar liberação" fica como botão secundário ao lado de "Encerrar". OS aberta ainda não reconhecida mostra "Reconhecer" como secundário; para o técnico da equipe, OS aberta ou planejada mostra "Assumir". OS encerrada vira somente leitura, com faixa informativa "OS encerrada em 12/11/2026 por Carlos Lima. Não pode mais ser editada."

**OS com executor fornecedor.** O painel lateral troca RESPONSÁVEL por EXECUTOR "Refrigeração Polar · serviço Manutenção de câmara fria (SLA 4 h)" e mostra o bloco LINK DE EXECUÇÃO quando houver: "Enviado a joao@polar.com.br em 12/11 10:15 · válido até 15/11 10:15 · aberto em 12/11 11:02" com as ações [Copiar link] (só na criação) [Gerar novo] [Revogar]. Na aba Horas, apontamentos feitos pelo prestador ou pelo link levam o selo "Prestador" ou "Via link". Em ativo locado, CUSTO DE mostra o padrão da matriz ("Locatário · avaria por mau uso, matriz do CL-2026-000012") e o seletor para alterar com justificativa.

### 8.6 Quadro de programação

```
Programação (h2)      [Semana 46 ‹ ›] [Hoje] [Dia|Semana|Mês] [Equipe ▾] [Unidade ▾]
┌────────────────────┬──────────┬──────────┬──────────┬──────────┬──────────┐
│ BACKLOG (filtrável)│ Seg 10/11│ Ter 11/11│ Qua 12/11│ Qui 13/11│ Sex 14/11│
│ OS sem data        ├──────────┼──────────┼──────────┼──────────┼──────────┤
│ [cartão OS]        │Ana  6/8 h│Ana  9/8 h│          │          │          │
│ [cartão OS]        │[OS][OS]  │[OS][OS] ! │          │          │          │
│ [cartão OS]        ├──────────┼──────────┼──────────┼──────────┼──────────┤
│ ...                │João 4/8 h│...       │          │          │          │
└────────────────────┴──────────┴──────────┴──────────┴──────────┴──────────┘
```

- Linhas por técnico ou equipe; colunas por dia. Célula mostra a carga "6/8 h" e uma barra de capacidade; acima de 100 % a barra fica em Atenção e o texto diz "9/8 h · 1 h acima".
- Cartão de OS: número, TAG, tipo (ícone), prioridade, horas estimadas. Arrastar do backlog para a célula programa; alternativa por teclado e menu "Programar para…".
- Técnico com capacitação vencida para a OS aparece com a célula bloqueada e motivo no tooltip.
- Mudanças são rascunho até "Publicar programação" (primário), que notifica os técnicos.

### 8.7 Plano preventivo

```
Planos de manutenção › PRE-0042 Troca de óleo e filtros
[● Ativo]  Aplica-se a: classe Cavalo 6x2 (14 ativos)          [Gerar OS agora] [Salvar plano]
┌─────────────────────────────────────────────────────────────────────────┐
│ GATILHO   ( ) Tempo  ( ) Medidor  (•) O que vencer primeiro             │
│ A cada [20.000] km   ou   a cada [6] meses                              │
│ Antecedência para gerar OS: [1.000] km ou [7] dias                      │
├─────────────────────────────────────────────────────────────────────────┤
│ PROCEDIMENTO: tarefas, peças, ferramentas, tempo estimado (2,5 h)       │
├─────────────────────────────────────────────────────────────────────────┤
│ PREVISÃO (14 ativos)                                                    │
│ TAG       Última execução     Próximo vencimento        Situação        │
│ FRT-0112  15/06/2026 · 226.010 km  em 2.300 km (≈ 9 dias)  [Operando]   │
│ FRT-0120  ...                 vencido há 400 km          [Atenção]      │
└─────────────────────────────────────────────────────────────────────────┘
```

A frase de resumo acima dos campos é gerada a partir do gatilho: "Gera OS a cada 20.000 km ou 6 meses, o que vier primeiro, 1.000 km ou 7 dias antes."

### 8.8 Checklist (modelo e execução no web)

- Modelo: lista de itens arrastáveis com tipo (Sim/Não, Número com faixa, Foto, Texto), obrigatório, "Item crítico: bloqueia o uso se reprovado" e "Reprovar gera solicitação". Pré-visualização do app ao lado.
- Execução vista no web: cabeçalho com ativo, executor, data e hora, resultado ("Aprovado com 2 ressalvas"); itens com resposta, foto e solicitação gerada linkada. Itens reprovados no topo.

### 8.9 Estoque geral e por subcontratada

O **estoque geral** é a soma de todos os almoxarifados: próprios, em poder de subcontratadas (`em_terceiro`) e de quarentena. O seletor de almoxarifado tem três grupos e a opção "Estoque geral".

```
Estoque (h2)        [Almoxarifado: Estoque geral ▾]          [Movimentar ▾] [Novo item]
                     ├ Estoque geral (todos)
                     ├ PRÓPRIOS: Central, Base Contagem
                     ├ EM TERCEIROS: Refrigeração Polar, Hidráulica Sul
                     └ QUARENTENA: Quarentena Central
[Abaixo do mínimo (5)] [Reservados] [Em terceiros] [Em conserto (3)] [Sem giro 180 dias] [Todos]
┌──────────┬────────────────────┬──────────┬──────────┬───────────┬──────────┬──────────┬──────────┐
│ Código   │ Item               │ Geral    │ Próprios │ Em terceir│ Quarent. │ Mínimo   │ Custo méd│
│ PSL-118  │ Filtro de óleo     │ 14 un.   │ 2 un.  ! │ 12 un.    │ 0 un.    │ 6 un.    │ R$ 48,90 │
│ CMP-V12  │ Válvula compressor │ 4 un. Nº │ 1 un.    │ 2 un.     │ 1 un.    │ 2 un.    │ R$ 812,00│
└──────────┴────────────────────┴──────────┴──────────┴───────────┴──────────┴──────────┴──────────┘
Barra de lote: "5 itens selecionados · Gerar requisição"
```

- Saldo abaixo do mínimo leva ícone `triangle-alert` e o texto "abaixo do mínimo" no tooltip e no leitor de tela. O mínimo é avaliado sobre os almoxarifados próprios, salvo configuração do item; o estoque em terceiros aparece ao lado para o planejador decidir se pede devolução ou compra.
- Item que controla série leva o selo `barcode` "Nº"; item reparável leva o selo `refresh-ccw` "Reparável". Os dois são neutros (texto e ícone em `texto-secundario`).
- "Movimentar" é um menu com os tipos canônicos: Entrada, Transferência, Remessa a terceiro, Retorno de terceiro, Envio para conserto, Retorno de conserto, Ajuste, Baixa para sucata. Saída e devolução de OS acontecem dentro da OS.

**Visão por subcontratada** (almoxarifado `em_terceiro` selecionado):

```
Estoque › Em terceiros › Refrigeração Polar (h2)   [Registrar remessa] [Registrar retorno]
Fornecedor: Refrigeração Polar · CNPJ 12.345.678/0001-90 · Subcontratada
Saldo em poder do fornecedor: 38 itens · R$ 9.412,30         Última conferência 01/11/2026
┌──────────┬────────────────────┬──────────┬───────────────┬────────────────────────────┐
│ Código   │ Item               │ Saldo    │ Consumo 30 d  │ Última remessa             │
│ PSL-118  │ Filtro de óleo     │ 12 un.   │ 7 un. (OS)    │ 04/11 · NF-e 45.812        │
└──────────┴────────────────────┴──────────┴───────────────┴────────────────────────────┘
[Movimentações] [Peças com série (4)] [Divergências de conferência]
```

- Registrar remessa e retorno pedem o **número e a chave da NF-e emitida pelo ERP** (campo de chave com 44 dígitos, validação de dígito verificador e leitura do código de barras do DANFE no tablet). A Tenaz registra a nota, não a emite; o texto de ajuda diz isso.
- O consumo do prestador nas OS baixa deste almoxarifado e aparece com o selo "Prestador".
- **Peça em trânsito** (remessa registrada, recebimento ainda não confirmado pela subcontratada) aparece na coluna "Em trânsito" e não pode ser consumida: no seletor de peças da OS ela fica desabilitada com o texto "Em trânsito · aguardando confirmação de recebimento". O prestador confirma o recebimento no portal (8.17) com o botão "Confirmar recebimento".
- **Quarentena:** peça com defeito aguarda a decisão de destino (Enviar para conserto, Baixar para sucata, Devolver ao fornecedor em garantia). O botão "Decidir destino" só aparece para quem tem a permissão `est.decidir_destino` (padrão gestor, planejador e almoxarife); para os demais, a linha mostra "Aguardando decisão de destino".
- O detalhe do item mostra saldos por almoxarifado em três grupos (próprios, em terceiros, quarentena), reservas por OS, histórico de movimentações e gráfico de consumo mensal. Item com série lista cada peça: número de série, situação (Em estoque, Instalada em FRT-0112, Em trânsito, Em conserto, Baixada), selo "Recondicionada" e garantia.

### 8.10 Frota: veículo e pneus por eixo

O detalhe do veículo segue o detalhe do ativo (8.3) com as abas Pneus, Abastecimentos, Multas, Sinistros e Documentos.

**Aba Pneus — diagrama por eixo**

```
Pneus de FRT-0112 (Cavalo 6x2)         [Rodízio] [Montar] [Desmontar]

 Dianteiro      [1E]──────────────[1D]
                 7,2 mm            6,9 mm
 Tração         [2EE][2EI]────[2DI][2DE]
                 9,1  9,0       5,1 ! 8,8
 Truck          [3EE][3EI]────[3DI][3DE]
 Estepe         [E1]

┌───────────────────────────────────────────────────────────────┐
│ POSIÇÃO 2DI · Fogo 004512 · Michelin X Multi D 295/80R22.5   │
│ Vida 2 (1ª recapagem) · Sulco 5,1 mm (mín. 3,0 mm)           │
│ Km nesta vida 48.200 km · Custo/km R$ 0,018                  │
│ [Registrar sulco] [Ver histórico do pneu]                    │
└───────────────────────────────────────────────────────────────┘
```

- Cada pneu é um botão com posição, número de fogo e sulco em `dado`; o estado do sulco usa as cores de status com texto ("5,1 mm · perto do mínimo"). Os layouts de eixo vêm da classe do veículo (4x2, 6x2, 6x4, carreta de 3 eixos).
- Rodízio e troca por arrastar entre posições, com alternativa "Mover para posição…" e confirmação mostrando de/para.
- A versão de tabela (posição, fogo, marca, vida, sulco, km) está sempre disponível e é a alternativa acessível.

### 8.11 Conformidade: vencimentos

```
Vencimentos (h2)                                          [Novo documento]
[Vencidos (3)] [Até 30 dias (11)] [Até 90 dias (27)] [Todos]
[Tipo ▾: CRLV-e, Tacógrafo, CIV, CIPP, Laudo NR-13, Seguro, CNH, MOPP, Capacitações]
[Vínculo ▾: Ativo, Pessoa, Empresa]
┌─────────────────┬──────────────────┬──────────────┬───────────────────────────┬──────────┐
│ Documento       │ Vinculado a      │ Vencimento   │ Situação                  │ Anexo    │
│ CRLV-e          │ FRT-0112         │ 24/11/2026   │ [ ! Vence em 12 dias]      │ PDF      │
│ NR-10           │ João Lima        │ 03/11/2026   │ [⊘ Vencido · bloqueia]    │ PDF      │
└─────────────────┴──────────────────┴──────────────┴───────────────────────────┴──────────┘
[Calendário dos próximos 90 dias, por semana]
```

Os itens bloqueantes explicam o efeito: "Enquanto vencido, João Lima não pode receber OS de serviço elétrico." As telas de NR-12 e NR-13 seguem a mesma estrutura, com os campos obrigatórios das normas.

### 8.12 Indicadores

```
Indicadores (h2)        [Período ▾] [Unidade ▾] [Classe ▾] [Criticidade ▾]  [Exportar] [Salvar visão]
[Disponibilidade] [MTBF] [MTTR] [Preventivas no prazo] [Custo] [Backlog] [Frota]
┌────────────────────────────────────────────────────────────────────────┐
│ DISPONIBILIDADE · 97,4 %  · Meta 96 %          ⓘ Como calculamos       │
│ [gráfico de linha diário com meta e ponto final destacado]             │
│                                               [Ver como tabela]        │
├───────────────────────────────────┬────────────────────────────────────┤
│ PIORES 10 ATIVOS (barras horizont.)│ POR CLASSE (tabela com sparklines) │
└───────────────────────────────────┴────────────────────────────────────┘
```

Cada visão de indicador tem o mesmo esqueleto: número e meta, tendência, quebra (piores ativos, por classe, por unidade), lista de registros. Relatórios agendados por e-mail usam os mesmos componentes renderizados em PDF.

### 8.13 Cadastro de fornecedor: papéis e catálogo de serviços

O `fornecedor` é o cadastro único de empresas externas. Uma empresa pode ter vários papéis; as abas aparecem conforme os papéis marcados.

```
Fornecedores › Refrigeração Polar
┌──────────────────────────────────────────────────────────────────────────┐
│ Refrigeração Polar Ltda (h2)   [● Ativo]               [Mais ▾] [Salvar] │
│ CNPJ 12.345.678/0001-90 · Contagem/MG · contato@polar.com.br             │
├──────────────────────────────────────────────────────────────────────────┤
│ PAPÉIS (marque todos que se aplicam)                                     │
│ [x] Subcontratada         executa manutenção com equipe própria         │
│ [x] Assistência técnica   autorizada ou especializada por marca         │
│ [ ] Fornecedor de peças   [ ] Oficina   [ ] Locador                     │
│ [ ] Posto de combustível  [ ] Seguradora                                │
├──────────────────────────────────────────────────────────────────────────┤
│ [Dados] [Serviços (3)] [Usuários (2)] [Conexão] [OS (41)]                │
│ [Estoque em poder (38)] [Contratos] [Avaliação] [Documentos]             │
├──────────────────────────────────────────────────────────────────────────┤
│ SERVIÇOS                                               [Novo serviço]    │
│ Descrição              Classes / marcas      Área     Atend. Solução Preço│
│ Manutenção de câmara   Câmara fria · Bitzer  MG, SP   4 h    24 h  R$ 180,00/h│
│ fria                                                                     │
│ Visita técnica         Compressor · todas    Grande BH 8 h   —    R$ 350,00/visita│
└──────────────────────────────────────────────────────────────────────────┘
```

- Papéis são caixas de seleção com uma linha de explicação cada. Pelo menos um é obrigatório. Desmarcar um papel em uso (ex.: Subcontratada com estoque em poder) é bloqueado com explicação: "Este fornecedor tem 38 itens em estoque. Registre o retorno antes de remover o papel Subcontratada."
- Formulário de serviço (`servico_fornecedor`): descrição, classes de ativo (seletor múltiplo), marcas atendidas (etiquetas), área de atendimento (texto e UFs), prazo de atendimento e de solução em horas (`dado`, com unidade "h"), preço de referência em R$ e unidade (Hora, Visita, Serviço). A frase de resumo segue o padrão do plano preventivo: "Atende em até 4 h e resolve em até 24 h, a R$ 180,00 por hora."
- Aba **Usuários**: usuários com perfil `prestador` desse fornecedor, com o aviso fixo "Estes usuários veem só as OS atribuídas à Refrigeração Polar e o estoque em poder dela." Convidar usuário pede nome, e-mail e se executa OS (conta como técnico na cobrança).
- Aba **Conexão** (papéis Subcontratada e Assistência técnica): mostra se o fornecedor tem conta Tenaz Serviços conectada a esta empresa, com o chip de estado da conexão (7.5) e a ação "Convidar para conectar" ou "Ver conexão" (8.27). Com conexão ativa, a aba Usuários mostra o aviso "Este fornecedor tem conta própria conectada. Os técnicos dele trabalham na conta dele e não contam como técnicos na sua cobrança."
- Aba **Avaliação** (Fase 2): SLA cumprido em %, retrabalho em %, nota média; cada número com "Como calculamos".
- Na OS, o seletor de executor sugere primeiro os fornecedores cujo catálogo atende a classe e a marca do ativo, com o prazo de atendimento ao lado.

### 8.14 Ordem de reparo

Envio de peças com defeito para conserto (EST, Fase 2). Lista igual à de OS (6.4) com colunas OR, Fornecedor, Itens, Estado, Enviada em, Prazo previsto, Custo (R$).

```
Ordens de reparo › OR-2026-000045
┌──────────────────────────────────────────────────────────────────────────┐
│ OR-2026-000045 · Válvula do compressor (h2)                              │
│ [Em orçamento] [Prazo previsto 20/11]                                    │
│                                   [Mais ▾] [Declarar irreparável] [Aprovar orçamento] │
├──────────────────────────────────────────────────────────────────────────┤
│ FLUXO: Aberta ─ Enviada ─ Em orçamento ─ Aprovada ─ Retornada            │
│        (saídas: Irreparável, Cancelada)                                  │
├───────────────────────────────────────────────┬──────────────────────────┤
│ ITENS                                         │ FORNECEDOR               │
│ CMP-V12 Válvula · Nº série VX-88213           │ Hidráulica Sul           │
│   Retirada na OS-2026-004809 de CMP-0347      │ Assistência técnica      │
│   Defeito: vazamento na sede                  │ ORIGEM Quarentena Central│
│                                               │ NF-E DE REMESSA          │
│ ORÇAMENTO                                     │ 45.901 · chave ...4471   │
│ Recebido 14/11 · R$ 420,00 · prazo 5 dias     │ NF-E DE RETORNO          │
│ Garantia do reparo: 90 dias   [anexo PDF]     │ (pendente)               │
│                                               │ CUSTO DO REPARO R$ 420,00│
│ LINHA DO TEMPO (6.8)                          │ GARANTIA 90 dias         │
└───────────────────────────────────────────────┴──────────────────────────┘
```

- A ação primária segue o estado: Enviar → Registrar orçamento → Aprovar orçamento → Registrar retorno. Recusar o orçamento abre a escolha "Declarar irreparável" ou "Abrir nova ordem com outro fornecedor".
- Enviar e Registrar retorno pedem número e chave da NF-e emitida pelo ERP, como a remessa a terceiro.
- Registrar retorno confirma o almoxarifado de destino e mostra o efeito: "A peça VX-88213 volta ao Almoxarifado Central como recondicionada, com custo de R$ 420,00 e garantia até 19/02/2027."
- Declarar irreparável pede o destino (Sucata ou Devolvida sem reparo) e gera a baixa para sucata quando for o caso.

### 8.15 Contrato de locação

LOC, Fase 2. Um contrato tem um **sentido**: Cedido (ativo próprio locado a um locatário) ou Tomado (ativo de um locador usado pela empresa). O sentido é a primeira escolha do formulário e muda os campos.

```
Locação › CL-2026-000012
┌──────────────────────────────────────────────────────────────────────────┐
│ CL-2026-000012 · Construtora Alfa (h2)   [Vigente] [Termina em 9 dias]   │
│ Cedido · 01/09/2026 a 21/11/2026 · Obra BR-381 km 412                    │
│                                       [Mais ▾] [Registrar devolução]    │
├──────────────────────────────────────────────────────────────────────────┤
│ [Ativos (3)] [Responsabilidade] [Vistorias (4)] [OS (6)] [Cobrança]      │
├──────────────────────────────────────────────────────────────────────────┤
│ ATIVOS                                                                   │
│ TAG      Ativo              Entrega     Devolução   Franquia   Uso       │
│ ESC-0031 Escavadeira 20 t   01/09 [✓]   —           600 h      548 h     │
│ GER-0007 Gerador 150 kVA    01/09 [✓]   —           600 h      612 h  !  │
├──────────────────────────────────────────────────────────────────────────┤
│ RESPONSABILIDADE DA MANUTENÇÃO (matriz)                                  │
│ Serviço                     Quem paga                                    │
│ Preventiva                  [Empresa ▾]                                  │
│ Corretiva por desgaste      [Empresa ▾]                                  │
│ Avaria por mau uso          [Locatário ▾]                                │
│ Pneus                       [Locatário ▾]                                │
│ Deslocamento                [Empresa ▾]                                  │
└──────────────────────────────────────────────────────────────────────────┘
```

- Estados do contrato e cores em 7.5. A ação primária segue o estado: Ativar (rascunho) → Registrar devolução (vigente). Ativar exige a vistoria de entrega de cada ativo. "Cancelar contrato" só existe no Rascunho; no Vigente, o menu "Mais" oferece "Encerrar contrato", com motivo obrigatório quando for antes do fim previsto.
- Campos do cabeçalho: local de uso (`local_uso`), franquia de horas (`franquia_horas`) e de km (`franquia_km`), vigência.
- **Sem sublocação:** no contrato Cedido, o seletor de ativos não oferece ativos "De terceiro"; quem procurar um deles vê "FRT-0140 é tomado em locação e não pode ser cedido a outro locatário."
- No contrato Tomado, a devolução do ativo ao locador o deixa Inativo; a confirmação diz: "Depois da devolução, FRT-0140 passa a Inativo e sai da contagem de veículos."
- A matriz de responsabilidade é uma tabela de seletores com os valores canônicos (Empresa, Locatário, Locador, Garantia, Fornecedor). A frase abaixo da tabela explica o efeito: "OS de ativos deste contrato já nascem com o responsável pelo custo desta matriz. Quem abrir a OS pode alterar com justificativa."
- Franquia excedida usa o indicador de saldo do estoque (`triangle-alert` + texto "12 h acima da franquia"), nunca só cor.
- Contrato **tomado**: o cabeçalho mostra o locador (fornecedor com papel Locador) e o SLA de substituição; a aba Ativos lista ativos com o selo "De terceiro".
- Aba **Cobrança**: relatório de avarias e horas excedentes com custo do locatário, com [Exportar]. Texto fixo: "A fatura é emitida pelo ERP."

### 8.16 Vistoria de entrega e de devolução

Feita no app ou no tablet, ao lado do ativo, com as duas partes presentes. Usa o modo guiado do checklist (9.6).

```
Vistoria de devolução · ESC-0031 · CL-2026-000012          Etapa 3 de 5
┌──────────────────────────────────────────────────────────────┐
│ 1 Leituras  2 Checklist  [3 Avarias]  4 Fotos  5 Assinaturas │
├──────────────────────────────────────────────────────────────┤
│ AVARIAS ENCONTRADAS                                          │
│ ┌──────────────────────────────────────────────────────────┐ │
│ │ [foto] Cilindro da lança com vazamento                    │ │
│ │        Comparação: sem avaria na entrega (01/09)          │ │
│ │        Custo: Locatário (mau uso, matriz)                 │ │
│ └──────────────────────────────────────────────────────────┘ │
│ [ + Registrar avaria ]                                       │
│                                                              │
│ [ Anterior ]                              [   Próximo   ]    │
└──────────────────────────────────────────────────────────────┘

Etapa 5 · Assinaturas
┌──────────────────────────────────────────────────────────────┐
│ Resumo: 214 h no horímetro (entrega 0 h) · 18 itens · 1 avaria│
│ "Declaro que conferi o estado do equipamento descrito acima." │
│ PELA EMPRESA       Carlos Lima       [✓ Assinado 14:20]       │
│ PELO LOCATÁRIO     Nome [__________]  Documento [__________] │
│                    [ área de assinatura 120 dp ]             │
│                    [ Limpar ]                                │
│ [            Registrar vistoria            ]                 │
└──────────────────────────────────────────────────────────────┘
```

- Etapas: leituras de medidor (com a da entrega como referência), checklist do modelo de vistoria, avarias, fotos e assinaturas. Indicador de etapa no topo, como o progresso do checklist.
- **Fotos de avaria:** cada avaria tem pelo menos uma foto obrigatória, com anotação simples (seta e círculo, 9.7) e descrição. Na devolução, a tela mostra ao lado a foto da mesma região tirada na entrega, quando houver, com o rótulo "Entrega 01/09" e "Devolução 21/11". Fotos gerais do ativo (frente, laterais, traseira, painel) são uma grade de 6 posições com contorno tracejado até serem preenchidas.
- **Assinatura das duas partes:** o representante da empresa assina com o próprio usuário (assinatura avançada, reautenticação); o representante do locatário assina na tela com nome, documento e rubrica (assinatura simples), ou com o próprio usuário `locatario`, se tiver. "Registrar vistoria" só fica ativo com as duas assinaturas; a dica explica o que falta.
- Avaria na devolução com custo do locatário gera OS automaticamente; o toast diz: "Vistoria registrada. OS-2026-004851 criada para o cilindro da lança, com custo do locatário."
- Funciona offline como o checklist; o PDF da vistoria (com fotos, leituras, avarias e manifesto de assinaturas) é gerado na sincronização e enviado ao contato do locatário.

### 8.17 Portal do prestador (app e web)

Usuário com perfil `prestador` vê **só as OS atribuídas ao seu fornecedor** e o almoxarifado em poder dele. A interface deixa o limite claro e não mostra menus ou funções que ele não pode usar.

**Web** — layout restrito, sem o menu lateral completo:

```
┌──────────────────────────────────────────────────────────────────────────┐
│ [símbolo] Tenaz · Refrigeração Polar para Distribuidora Norte  [Ana ▾]   │
│ [OS (7)] [Estoque em poder] [Meus dados]                                 │
├──────────────────────────────────────────────────────────────────────────┤
│ OS da Distribuidora Norte atribuídas à Refrigeração Polar (h2)           │
│ [A fazer (5)] [Concluídas aguardando liberação (2)] [Encerradas]         │
│ OS              Ativo                      Prioridade  Prazo    Estado   │
│ OS-2026-004812  CFR-0021 Câmara fria 2     Urgente     Hoje     Programada│
└──────────────────────────────────────────────────────────────────────────┘
```

- O cabeçalho sempre diz **para qual empresa** o prestador está trabalhando ("Refrigeração Polar para Distribuidora Norte"). Quem atende duas contratantes troca pelo menu do usuário, como na troca de empresa, e o cabeçalho muda.
- O detalhe da OS é o mesmo de 8.5, sem custos da empresa, sem histórico de outras OS do ativo e sem o seletor de executor. As ações são as do técnico: Reconhecer, Iniciar, Pausar, Concluir, apontar horas, registrar peças do estoque em poder, fotos e assinatura.
- "Estoque em poder" é a visão por subcontratada (8.9), só leitura, mais o registro de consumo e de devolução dentro da OS.
- Sem busca global, sem indicadores, sem IA, sem exportação.
- O portal é para a subcontratada **sem** conta própria (escopo `prestador`). A subcontratada com conta Tenaz Serviços conectada não usa o portal: trabalha na conta dela (8.20 a 8.28), e a OS chega lá pela conexão.

**App** — mesma navegação por abas (9.1), coluna Prestador: Hoje (OS do fornecedor atribuídas a ele), Estoque em poder, Ler QR, Mais. Ler QR de um ativo sem OS do fornecedor mostra "Este ativo não tem OS atribuída à Refrigeração Polar." e não mostra dados do ativo.

### 8.18 Página do link de uso único (mobile, sem login)

Para o prestador ocasional, que recebe o link por e-mail ou mensagem. Página web leve, feita para celular, sem login, sem menu e sem acesso a nada além de uma OS. O endereço tem a forma `app.tenaz.com.br/l/#<token>`: o código fica depois do `#`, a página o troca por uma sessão curta na abertura e o apaga da barra de endereço. A interface nunca mostra nem pede para copiar o código.

```
┌───────────────────────────────┐
│ [símbolo] Tenaz               │
│ Distribuidora Norte pede o    │  titulo2
│ registro deste serviço        │
│ OS-2026-004812                │  dado-forte
│ Troca do ventilador do        │  corpo-forte
│ evaporador                    │
│ CFR-0021 Câmara fria 2        │
│ Rua A, 100 · Contagem/MG      │
│ Válido até 15/11 às 10:15     │  legenda
├───────────────────────────────┤
│ 1 Horas  2 Peças  3 Fotos     │
│ 4 Observação  5 Assinatura    │
├───────────────────────────────┤
│ Início   [12/11 13:00]        │
│ Fim      [12/11 15:30]        │
│ Executado por [__________]    │
│ Documento     [__________]    │
│                               │
│ [        Próximo        ]     │  56 dp, Sinal
│ Dúvidas: Carlos Lima          │
│ (31) 3333-0000                │
└───────────────────────────────┘
```

- Primeiro acesso: se a empresa exigir código de confirmação, a página pede "Informe o código de 6 dígitos enviado para jo***@polar.com.br" antes de mostrar a OS.
- Etapas iguais ao fluxo do app (9.5), com alvos de 48 dp e ação principal de 56 dp em Sinal. Fotos pela câmera do navegador; peças em texto livre ou, se o fornecedor tiver estoque em poder, por lista.
- Assinatura simples na tela com o texto de declaração. "Enviar e concluir" mostra o resumo antes de confirmar; é um único envio, que registra a execução e conclui a OS (`registrar-execucao`). Fotos são enviadas durante o preenchimento, uma a uma, com o estado de cada uma.
- Depois do envio: "Execução registrada na OS-2026-004812. Obrigado. Este link não pode ser usado de novo." Sem botão de voltar à edição.
- Link inválido, expirado, revogado ou já usado: uma única mensagem neutra, sem dizer qual caso é: "Link inválido ou expirado. Peça um novo a quem enviou." Link aberto em outro aparelho: "Este link já está em uso em outro aparelho. Peça um novo a quem enviou."
- Sem conexão no meio do preenchimento: o que foi digitado fica salvo no navegador e a faixa diz "Sem conexão. Seus dados estão salvos neste aparelho; envie quando houver sinal."
- A página não tem rastreadores, fontes de CDN externas nem links para outras áreas da Tenaz.

### 8.19 Portal do locatário

Usuário com perfil `locatario` (gratuito) vê só os ativos locados a ele, abre solicitações e acompanha OS desses ativos. No lançamento do LOC o locatário usa **só o portal web**; não há aba de locatário no app (9.1). O portal funciona no navegador do celular, com o layout de 320 px (11.5).

```
┌──────────────────────────────────────────────────────────────────────────┐
│ [símbolo] Tenaz · Locadora Norte para Construtora Alfa        [Ana ▾]   │
│ [Equipamentos (3)] [Solicitações (2)] [Contratos] [Vistorias]            │
├──────────────────────────────────────────────────────────────────────────┤
│ Equipamentos locados (h2)                         [Abrir solicitação]    │
│ ┌───────────────────────────┐ ┌───────────────────────────┐              │
│ │ Escavadeira 20 t [Operando]│ │ Gerador 150 kVA [Parado]  │              │
│ │ ESC-0031                  │ │ GER-0007                  │              │
│ │ CL-2026-000012 até 21/11  │ │ OS-2026-004840 Em execução│              │
│ │ Horímetro 548 h de 600 h  │ │ Previsão: 13/11           │              │
│ │ Próxima preventiva 600 h  │ │ [Ver OS]                  │              │
│ └───────────────────────────┘ └───────────────────────────┘              │
└──────────────────────────────────────────────────────────────────────────┘
```

- Cartões de ativo na variante compacta (6.7), sem custos, sem medidores além do usado na franquia e sem histórico anterior à entrega.
- "Abrir solicitação" é a ação primária: escolher o equipamento (ou ler o QR), descrever por texto, foto ou voz, responder "O equipamento parou?". Mesmo formulário do solicitante. A solicitação nasce com origem `portal_locatario`, mostrada como "Portal do locatário" na triagem.
- OS dos ativos locados aparecem com estado, prazo e responsável pelo custo ("Custo: Locatário · avaria por mau uso"), sem valores internos. O valor cobrado aparece só no relatório de cobrança, quando a empresa publicar.
- Contratos e vistorias: consulta e download do PDF da vistoria assinada.
- Sem busca global, indicadores, IA ou exportação. Depois do fim do último contrato, o portal fica disponível só para consulta pelo prazo definido pela empresa.

### 8.20 Painel de serviços (conta da subcontratada)

Tela inicial da conta com o plano Serviços (SRV). Mesma estrutura do painel do gestor (8.1), com o recorte de quem presta serviço: por cliente, por técnico e o que está atrasado.

```
┌──────────────────────────────────────────────────────────────────────────┐
│ Painel de serviços (h1)            [Esta semana ▾] [Todos os clientes ▾] │
├──────────────────────────────────────────────────────────────────────────┤
│ [KPI OS abertas 38] [KPI Atrasadas 5] [KPI No prazo 91 %] [KPI 312 h]    │
├───────────────────────────────────────────┬──────────────────────────────┤
│ OS POR CLIENTE                            │ PRECISA DE VOCÊ              │
│ Cliente               Abertas Atras. Mês  │ 4 OS sem técnico             │
│ Distribuidora Norte ⇄    14     3    22   │   [Distribuir]               │
│ Frigorífico Sul ⇄         9     1     8   │ 2 liberações recusadas       │
│ Padaria Pão Bom           3     0     5   │   pelo cliente     [Ver]     │
│ (barras horizontais com número ao lado)   │ 1 convite de conexão [Ver]   │
├───────────────────────────────────────────┤ 1 OS não sincronizada [Ver]  │
│ OS POR TÉCNICO (hoje)                     │ 3 relatórios a enviar [Ver]  │
│ Técnico       Em exec. A fazer  Carga     ├──────────────────────────────┤
│ Ana Souza        1       3     5,5/8 h    │ ATRASOS                      │
│ João Lima        0       5     2,0/8 h    │ OS-2026-000311 · 2 dias      │
│ Marcos Reis      1       2     9,0/8 h !  │  Distribuidora Norte·CFR-0021│
│                                           │ OS-2026-000298 · 1 dia       │
│                                           │  Frigorífico Sul · CMP-0012  │
└───────────────────────────────────────────┴──────────────────────────────┘
```

- `⇄` representa o ícone `cable` (cliente conectado, com conta Tenaz). Cliente sem conta Tenaz não tem ícone. O leitor de tela lê "Distribuidora Norte, conectada".
- **OS por cliente:** abertas, atrasadas e concluídas no mês; cada número leva à lista de OS filtrada. Ordenada por atrasadas, depois por abertas.
- **OS por técnico:** em execução, a fazer e carga do dia ("5,5/8 h"), com a mesma barra de capacidade do quadro de programação (8.6); acima de 100 % a barra fica em Atenção e o texto diz "9,0/8 h · 1 h acima".
- **Atrasos:** OS com prazo vencido, da mais atrasada para a menos, com o indicador de prazo (7.3). Nada é só cor.
- Sem ação primária em Sinal, como no painel do gestor. As ações ficam em "Precisa de você".

### 8.21 Técnicos e especialidades

```
Técnicos (h2)                                     [Especialidades] [Novo técnico]
[Buscar nome] [Especialidade ▾] [Equipe ▾] [Situação ▾]
┌──────────────┬──────────────────────────┬────────────────────────┬────────┬────────┐
│ Técnico      │ Especialidades           │ Capacitações           │ App    │ Situação│
│ Ana Souza    │ Refrigeração · Elétrica  │ NR-10 até 03/2027      │ v1.4.2 │ Ativo  │
│ João Lima    │ Refrigeração             │ ⊘ NR-10 vencida · bloq.│ v1.4.0 │ Ativo  │
│ Marcos Reis  │ Hidráulica               │ NR-35 até 08/2027      │ —      │ Convite│
└──────────────┴──────────────────────────┴────────────────────────┴────────┴────────┘
5 técnicos ativos contam na cobrança do plano Serviços.

Técnicos › Novo técnico
┌──────────────────────────────────────────────────────────────────────────┐
│ DADOS DA PESSOA                                                          │
│ Nome [_______________________]   CPF [___.___.___-__]                    │
│ Telefone de trabalho [_______________]                                   │
├──────────────────────────────────────────────────────────────────────────┤
│ ACESSO AO APP                                                            │
│ E-mail [_______________________]   [x] Enviar convite para o Tenaz Mobile│
│ Equipe [Refrigeração ▾]                                                  │
│ ⓘ Com Marcos, a conta passa a 6 técnicos cobrados.                      │
├──────────────────────────────────────────────────────────────────────────┤
│ ESPECIALIDADES                                                           │
│ [x] Refrigeração  [x] Elétrica  [ ] Hidráulica  [ ] Caldeiras            │
│ [+ Nova especialidade]                                                   │
├──────────────────────────────────────────────────────────────────────────┤
│ CAPACITAÇÕES E DOCUMENTOS                                                │
│ NR-10   válida até [__/__/____]  [Anexar]           [+ Adicionar]        │
├──────────────────────────────────────────────────────────────────────────┤
│ ⓘ Clientes conectados veem só o nome do técnico nas OS.                 │
│                                         [Voltar] [Salvar técnico]        │
└──────────────────────────────────────────────────────────────────────────┘
```

- Técnico é `pessoa` + `usuario` com perfil `tecnico`. As especialidades vêm do catálogo da própria conta (`especialidade_tecnico`), editado no diálogo [Especialidades] (nome, descrição, classes de ativo sugeridas); a ligação é `pessoa_especialidade`.
- Especialidades aparecem como etiquetas neutras (texto e borda em `texto-secundario`), nunca coloridas.
- O CPF aparece mascarado na lista e na busca (`***.456.789-**`); completo só no formulário, para quem tem permissão de editar pessoas.
- Capacitação vencida que bloqueia atribuição usa o visual de 7.6 ("Vencido · bloqueia atribuição") e é respeitada na abertura de OS (8.23) e na distribuição (8.24).
- A coluna App mostra a versão do app no último acesso, para o suporte; "—" quando o técnico ainda não entrou.
- A frase de cobrança fica sempre visível abaixo da lista e no formulário, antes de salvar um técnico novo (plano Serviços: só o perfil `tecnico` conta).

### 8.22 Clientes atendidos e equipamentos de clientes

```
Clientes atendidos (h2)                                  [Importar] [Novo cliente]
[Buscar nome, CNPJ ou CPF] [Conexão ▾: Todos, Conectados, Sem conta Tenaz]
┌──────────────────────┬──────────────────────┬────────┬──────────┬───────────────┐
│ Cliente              │ Documento            │ Equip. │ OS abert.│ Conexão       │
│ Distribuidora Norte  │ 12.345.678/0001-90   │   42   │    14    │ [⇄ Conectada] │
│ Frigorífico Sul      │ 23.456.789/0001-01   │   17   │     9    │ [⇄ Conectada] │
│ Padaria Pão Bom      │ 98.765.432/0001-10   │    3   │     3    │ —             │
│ Maria Oliveira       │ ***.456.789-**       │    1   │     0    │ —             │
└──────────────────────┴──────────────────────┴────────┴──────────┴───────────────┘

Clientes atendidos › Padaria Pão Bom
┌──────────────────────────────────────────────────────────────────────────┐
│ Padaria Pão Bom Ltda (h2)   Pessoa jurídica              [Mais ▾] [Abrir OS]│
│ CNPJ 98.765.432/0001-10 · Sem conta Tenaz          [Convidar para a Tenaz]│
├──────────────────────────────────────────────────────────────────────────┤
│ [Equipamentos (3)] [OS (8)] [Endereços (1)] [Contatos (2)] [Relatórios]  │
├──────────────────────────────────────────────────────────────────────────┤
│ EQUIPAMENTOS                                        [Novo equipamento]   │
│ TAG        Equipamento              Local de atendimento    Últ. serviço  │
│ PPB-0001   Forno lastro 4 esteiras  Loja Centro             02/11/2026   │
│ PPB-0002   Câmara de fermentação    Loja Centro             15/10/2026   │
└──────────────────────────────────────────────────────────────────────────┘
```

- **Novo cliente:** tipo de pessoa (Jurídica ou Física, controle segmentado) que troca a máscara do documento (CNPJ ou CPF, com dígito verificador); nome; contatos (nome, função, telefone, e-mail e a opção "Recebe o relatório da OS"); endereços de atendimento (rótulo, endereço, referência, horário de acesso). Para pessoa física, a ajuda diz: "Guarde só o necessário para o atendimento."
- **Equipamento de cliente** (`ativo` com `propriedade = de_cliente`): é o mesmo formulário de ativo (8.3), com o cliente fixo, o endereço de atendimento no lugar da localização e a TAG sugerida com o prefixo do cliente (PPB-0001). Número de série com botão de leitura (`scan-barcode`). Na lista geral de ativos da conta, leva o selo neutro "De cliente · Padaria Pão Bom".
- **QR do equipamento:** etiqueta com o nome da subcontratada e a TAG, para colar no equipamento do cliente; ler o QR no app abre o equipamento e as OS dele.
- **Cliente conectado** (aba Equipamentos da Distribuidora Norte): os equipamentos vêm da conta do cliente, com a etiqueta `cable` "Do cliente". Os campos de cadastro (TAG, nome, fabricante, modelo, série, local) ficam somente leitura, com ícone `lock` e o texto de ajuda "Cadastro mantido pela Distribuidora Norte". "Novo equipamento" fica desabilitado com a dica "Os equipamentos da Distribuidora Norte vêm da conta dela. Peça ao cliente para compartilhar o equipamento." A subcontratada pode acrescentar anotações próprias e imprimir o QR dela.
- "Convidar para a Tenaz" (cliente sem conta) envia um convite comercial por e-mail; não cria conexão nem compartilha dados.

### 8.23 Abertura de OS de serviço

Formulário de uma página, com três blocos na ordem da pergunta que o atendente faz ao cliente: onde, o quê, quem.

```
Ordens de serviço › Nova OS de serviço
┌──────────────────────────────────────────────────────────────────────────┐
│ CLIENTE E EQUIPAMENTO                                                    │
│ Cliente      [Distribuidora Norte ▾]   ⇄ Conectada                       │
│ Equipamento  [CFR-0021 Câmara fria 2 ▾]                       [Ler QR]   │
│ Local        CD Contagem · Rua A, 100 · Contagem/MG                      │
├──────────────────────────────────────────────────────────────────────────┤
│ SERVIÇO A REALIZAR                                                       │
│ Descrição  [Troca do ventilador do evaporador e limpeza do condensador ] │
│ Tipo [Corretiva planejada ▾]  Prioridade [Alta ▾]  Alta · padrão crit. A │
│ Serviços do catálogo                                     [+ Adicionar]   │
│   Manutenção de câmara fria     2,0 h      R$ 180,00/h     R$ 360,00     │
│   Visita técnica                1 visita   R$ 350,00       R$ 350,00     │
├──────────────────────────────────────────────────────────────────────────┤
│ TÉCNICO E PRAZO                                                          │
│ Técnico  [Ana Souza ▾]   Refrigeração · 5,5/8 h hoje                     │
│ Ajudantes [+ Adicionar]                                                  │
│ Data prevista [13/11/2026 08:00]      Prazo [13/11/2026 18:00]           │
├──────────────────────────────────────────────────────────────────────────┤
│ ⓘ A OS também será criada na conta da Distribuidora Norte, como Aberta, │
│   para a triagem do cliente.                                             │
│                       [Voltar] [Criar sem técnico] [Criar e enviar ao técnico]│
└──────────────────────────────────────────────────────────────────────────┘
```

- **Serviço a realizar** = descrição livre (`descricao_servico`) + itens do catálogo `servico_fornecedor` da própria conta (`servico_os`: serviço, quantidade, preço de referência). O total aparece em `dado`, alinhado à direita. A descrição é obrigatória; os itens, não.
- **Seletor de técnico:** primeiro os técnicos com a especialidade sugerida pelo serviço ou pela classe do equipamento, com a carga do dia; depois os demais, sob o separador "Sem a especialidade". Técnico com capacitação bloqueante vencida aparece desabilitado com o motivo.
- Ação primária: "Criar e enviar ao técnico" (cria a OS `programada` e notifica o técnico por push). "Criar sem técnico" é secundária e manda a OS para a coluna Sem técnico da distribuição (8.24).
- **Cliente conectado:** a faixa informativa avisa que a OS nasce também no cliente, como `aberta` e com origem "Aberta por Refrigeração Polar". Preços só vão ao cliente se o grupo Preços estiver ligado na conexão; o formulário diz qual é o caso ("Valores não são enviados à Distribuidora Norte").
- OS que chega do cliente conectado não passa por este formulário: aparece direto na coluna Sem técnico da distribuição, com o selo da OS do cliente (8.28).

### 8.24 Quadro de distribuição de tarefas por técnico

```
Distribuição (h2)  [Hoje ‹ ›] [Dia|Semana] [Cliente ▾] [Especialidade ▾] [Publicar distribuição]
┌──────────────────────┬──────────────────┬──────────────────┬──────────────────┐
│ SEM TÉCNICO (4)      │ Ana Souza        │ João Lima        │ Marcos Reis      │
│ por prazo            │ Refrig. · Elétr. │ Refrigeração     │ Hidráulica       │
│                      │ 5,5/8 h ▓▓▓▓▓░░  │ 2,0/8 h ▓▓░░░░░  │ 9,0/8 h ▓▓▓▓▓▓▓! │
├──────────────────────┼──────────────────┼──────────────────┼──────────────────┤
│┌────────────────────┐│┌────────────────┐│┌────────────────┐│┌────────────────┐│
││OS-2026-000315   ⇄  │││OS-2026-000311 ⇄│││OS-2026-000309  │││OS-2026-000302  ││
││Distribuidora Norte │││Em execução     │││Programada 14:00│││Pausada · peça  ││
││CFR-0030 · Alta     │││CFR-0021        │││PPB-0001        │││Atrasada 1 dia  ││
││Refrigeração · 2 h  ││└────────────────┘│└────────────────┘│└────────────────┘│
││Vence hoje          ││┌────────────────┐│                  │                  │
│└────────────────────┘││OS-2026-000314  ││                  │                  │
│ ...                  ││Programada 10:00││                  │                  │
└──────────────────────┴──────────────────┴──────────────────┴──────────────────┘
Rascunho: 3 mudanças não publicadas                        [Descartar]
```

- Colunas: Sem técnico (backlog, ordenado por prazo) e uma coluna por técnico ativo, com especialidades e carga do dia. Na visão Semana, cada coluna de técnico vira uma linha e os dias viram colunas, como o quadro de programação (8.6).
- Cartão: número da OS, cliente, TAG, prioridade, especialidade exigida, horas estimadas, estado (chip de 7.3) e indicador de prazo. OS vinculada a cliente conectado leva o ícone `cable` no canto, com nome acessível "OS da Distribuidora Norte".
- Arrastar o cartão para a coluna de um técnico atribui; alternativa por teclado e pelo menu "Atribuir a…". Técnico sem a especialidade aceita com aviso no cartão ("Marcos não tem a especialidade Refrigeração"); técnico com capacitação bloqueante vencida recusa, e a coluna mostra o motivo durante o arraste.
- As mudanças são rascunho até "Publicar distribuição" (primário), que notifica cada técnico afetado por push ("3 tarefas novas"). O rodapé mostra quantas mudanças faltam publicar.
- Técnicos podem assumir OS da coluna Sem técnico pelo app (9.5), se a empresa permitir; nesse caso o cartão some do backlog e aparece na coluna do técnico com "Assumida por João Lima às 09:02".

### 8.25 App do técnico da subcontratada

As telas do app estão na seção 9.8 (lista de tarefas e execução com componentes, número de série, fotos e assinatura do cliente) e 9.9 (troca de conta e notificações). São as mesmas telas do técnico da empresa, com os blocos do SRV.

### 8.26 Relatório PDF da OS

Gerado por `GET /v1/ordens-servico/{id}/relatorio-pdf` depois da conclusão. É o documento que a subcontratada entrega ao cliente, principalmente ao que não usa a Tenaz. Formato A4 retrato, margens de 15 mm, impressão em preto e branco sem perda de informação.

```
┌──────────────────────────────────────────────────────────────┐
│ [logo da subcontratada]                RELATÓRIO DE SERVIÇO  │  sobrelinha
│ Refrigeração Polar Ltda                OS-2026-000311        │  dado-forte
│ CNPJ 12.345.678/0001-90                Concluída 13/11/2026  │
├──────────────────────────────────────────────────────────────┤
│ CLIENTE      Distribuidora Norte Ltda · CNPJ 12.345.678/...  │
│ LOCAL        CD Contagem · Rua A, 100 · Contagem/MG          │
│ EQUIPAMENTO  CFR-0021 Câmara fria 2 · Bitzer 4FES · Nº 77812 │
│ OS DO CLIENTE OS-2026-004812 (quando conectado)              │
├──────────────────────────────────────────────────────────────┤
│ SERVIÇO REALIZADO                                            │
│ Troca do ventilador do evaporador e limpeza do condensador.  │
│ Manutenção de câmara fria · 2,5 h  ·  Visita técnica · 1     │
├──────────────────────────────────────────────────────────────┤
│ COMPONENTES                                                  │
│ Tipo          Item                 Qtd   Nº retirado  Nº instalado│
│ Substituído   Motor do ventilador  1 un  VT-88213     VT-90455    │
│ Consumido     Fita isolante        1 un  —            —           │
│ Do cliente    Filtro secador       1 un  —            —           │
├──────────────────────────────────────────────────────────────┤
│ HORAS                                                        │
│ Técnico      Início         Fim            Horas             │
│ Ana Souza    13/11 08:12    13/11 10:42    2,5 h             │
│ Início 13/11 08:12 · Conclusão 13/11 10:42 · Pausas 0 min    │
├──────────────────────────────────────────────────────────────┤
│ FOTOS (6)   [Antes 1] [Antes 2] [Depois 1] [Depois 2] [Peça] │
├──────────────────────────────────────────────────────────────┤
│ ASSINATURAS                                                  │
│ Cliente   [rubrica]  José Pereira · Encarregado · 13/11 10:45│
│ Técnico   Ana Souza · assinatura eletrônica · 13/11 10:44    │
│ Código de verificação 7F3A-21C9                         [QR] │
├──────────────────────────────────────────────────────────────┤
│ [símbolo] Gerado pela Tenaz · Emitido 13/11/2026 10:50 · 1/2 │  legenda
└──────────────────────────────────────────────────────────────┘
```

- **Marca:** o cabeçalho é da conta emissora (logo cadastrado em Administração, até 40 mm × 15 mm; sem logo, o nome em Inter 700). A Tenaz aparece só no rodapé, com o símbolo positivo de 4 mm e "Gerado pela Tenaz". Sinal só no ponto do símbolo; nenhuma outra cor da marca em texto. Títulos de bloco em `sobrelinha` (JetBrains Mono 600, caixa alta), texto em Inter 10 pt, dados (números, séries, TAG, horários) em JetBrains Mono 9 pt. Cabeçalhos de tabela com fundo Névoa e texto Aço.
- **Componentes:** a coluna Tipo usa os rótulos de `tipo_material` (Substituído, Consumido, Do cliente). Números de série em JetBrains Mono, nunca abreviados.
- **Fotos:** grade de 3 colunas, com legenda (Antes, Depois, Peça retirada, Placa), data e hora abaixo de cada uma. Mais de 9 fotos: as 9 marcadas como "incluir no relatório" e a frase "Outras 4 fotos disponíveis com a subcontratada".
- **Valores:** fora do relatório por padrão. A opção "Incluir valores" na geração acrescenta a coluna de preço de referência aos serviços e componentes e o total.
- **Assinaturas:** assinatura do cliente com nome, função e data e hora; assinatura eletrônica do técnico. O código de verificação e o QR levam a uma página pública que confirma número, data e hash do relatório, sem mostrar o conteúdo **(proposta)**. A natureza jurídica da assinatura está em 9_SEGURANCA_CONFIABILIDADE.md.
- **Acessibilidade e arquivo:** PDF marcado (tags de título, tabela e imagem com texto alternativo), idioma `pt-BR`, texto selecionável. Nome do arquivo: `OS-2026-000311_Distribuidora-Norte_2026-11-13.pdf`.
- **Ações no detalhe da OS:** [Baixar PDF] (secundário) e [Enviar ao cliente] (primário depois da conclusão), que envia aos contatos marcados "Recebe o relatório da OS" e registra o envio na linha do tempo.

### 8.27 Conexão com cliente (convite, campos compartilhados, estado)

A conexão aparece nas duas contas: no cliente, em Terceiros > Conexões com prestadores e na aba Conexão do fornecedor (8.13); na subcontratada, em Serviços > Conexões com clientes e na aba Conexão do cliente atendido. O texto sempre diz **o que vai e o que não vai** para a outra empresa.

**Convite (conta do cliente):**

```
Convidar Refrigeração Polar para conectar                           [×]
┌──────────────────────────────────────────────────────────────────────┐
│ Com a conexão, as OS que você atribuir à Refrigeração Polar chegam à │
│ conta dela na Tenaz, e o andamento volta para cá automaticamente.    │
│ Enviar convite para [joao@polar.com.br______]                        │
│ CNPJ do fornecedor   12.345.678/0001-90 (do cadastro)                │
├──────────────────────────────────────────────────────────────────────┤
│ EQUIPAMENTOS COMPARTILHADOS                                          │
│ (•) Escolher equipamentos  ( ) Por localização  ( ) Por classe       │
│ [x] CFR-0021 Câmara fria 2   [x] CFR-0030 Câmara fria 3   [ ] ...    │
│ Equipamentos de OS atribuídas à Refrigeração Polar são compartilhados│
│ junto com a OS.                                                      │
├──────────────────────────────────────────────────────────────────────┤
│ O QUE VAI PARA A REFRIGERAÇÃO POLAR                                  │
│ Sempre: cadastro dos equipamentos escolhidos (TAG, nome, fabricante, │
│ modelo, série, local), OS atribuídas a ela, liberação e recusa.      │
│ [ ] Documentos técnicos marcados como compartilháveis                │
│ Nunca: custos, outras OS, pessoas, documentos da empresa.            │
│ [Ver a lista completa de campos]                                     │
│                                        [Voltar] [Enviar convite]     │
└──────────────────────────────────────────────────────────────────────┘
```

**Convite recebido (conta da subcontratada):**

```
Conexões com clientes › Convite da Distribuidora Norte
┌──────────────────────────────────────────────────────────────────────┐
│ Distribuidora Norte Ltda (CNPJ 12.345.678/0001-90) quer conectar com │
│ a Refrigeração Polar. Enviado por Carlos Lima em 12/11 · válido até  │
│ 26/11/2026                                                           │
├──────────────────────────────────────────────────────────────────────┤
│ VOCÊ VAI RECEBER   42 equipamentos e as OS atribuídas a vocês        │
│ VOCÊ VAI ENVIAR    Andamento, horas, componentes, fotos, assinatura  │
│                    do cliente e o nome do técnico de cada OS         │
│ [ ] Enviar também os preços de referência de serviços e componentes  │
│ Nunca: CPF e contatos dos técnicos, outros clientes, outras OS.      │
├──────────────────────────────────────────────────────────────────────┤
│ CLIENTE ATENDIDO  (•) Ligar ao cliente "Distribuidora Norte" (42 OS) │
│                   ( ) Criar um cliente novo                          │
│                                     [Recusar] [Aceitar conexão]      │
└──────────────────────────────────────────────────────────────────────┘
```

**Detalhe da conexão (as duas contas):**

```
Conexões › Distribuidora Norte   [⇄ Conectada]            [Mais ▾]
[Resumo] [O que é compartilhado] [Equipamentos (42)] [OS vinculadas (36)] [Sincronização (1)] [Histórico]
┌──────────────────────────────────────────────────────────────────────────┐
│ RESUMO                                                                   │
│ Conectada desde 12/11/2026 · aceita por João Polar                       │
│ OS vinculadas em aberto 14 · Última sincronização há 12 s                │
├──────────────────────────────────────────────────────────────────────────┤
│ SINCRONIZAÇÃO                                                            │
│ OS               Sentido          Situação              Motivo           │
│ OS-2026-000298   Para o cliente   [! Não sincronizada]  Aguardando       │
│                                                         mensagem anterior│
│                                                         [Reprocessar]    │
└──────────────────────────────────────────────────────────────────────────┘
```

- Chip de estado da conexão conforme 7.5. "Mais" oferece Suspender, Reativar (quando suspensa) e Encerrar, para quem tem a permissão de gerenciar conexões.
- Aba **O que é compartilhado:** a lista fechada de campos em linguagem de usuário, por grupo e por sentido, com os grupos opcionais como interruptores. Cada lado só liga ou desliga os grupos dos dados que **envia**; o outro lado vê o estado, sem poder alterar.
- Aba **Sincronização:** só mostra problemas. Motivo em linguagem de usuário ("Aguardando mensagem anterior", "A OS foi cancelada no cliente", "Conexão suspensa") e a ação "Reprocessar". Sem problemas: "Tudo sincronizado. Última entrega às 14:32."
- Aba **Histórico:** convite, aceite, mudanças nos grupos compartilhados, suspensões e encerramento, com quem e quando.
- **Encerrar:** título "Encerrar a conexão com a Distribuidora Norte?" · texto "As OS deixam de ser sincronizadas. As duas empresas mantêm o histórico do que já receberam. OS em aberto continuam em cada conta." · botões [Voltar] [Encerrar conexão]. Suspender usa o mesmo padrão, com "As mudanças ficam guardadas e são enviadas quando a conexão for reativada."

### 8.28 OS vinculada vista pelo cliente e pela subcontratada

**Na conta do cliente**, a OS executada pela subcontratada conectada:

```
Ordens de serviço › OS-2026-004812
┌──────────────────────────────────────────────────────────────────────────┐
│ OS-2026-004812 · Troca do ventilador do evaporador (h2)                  │
│ [Em execução · Ana Souza] [Corretiva planejada] [Alta]                   │
│ [⇄ Executada por Refrigeração Polar]  OS-2026-000311 na conta dela       │
│ [cloud-check] Sincronizada às 10:12             [Mais ▾] [Ver relatório] │
├──────────────────────────────────────────────────────────────────────────┤
│ [Tarefas 3/5] [Componentes (2)] [Horas] [Fotos (6)] [Assinatura]         │
│                                                                          │
│ COMPONENTES                [lock] Registrado pela Refrigeração Polar     │
│ Substituído  Motor do ventilador   1 un   VT-88213 → VT-90455            │
│ Consumido    Fita isolante         1 un                                  │
│                                                ┌──────────────────────┐  │
│ HORAS                                          │ EXECUTOR             │  │
│ Ana Souza  13/11 08:12 – 10:42   2,5 h         │ Refrigeração Polar   │  │
│                                                │ Técnico: Ana Souza   │  │
│                                                │ CUSTO DE Empresa     │  │
│                                                │ ORIGEM Manual        │  │
│                                                └──────────────────────┘  │
└──────────────────────────────────────────────────────────────────────────┘
```

- **Selo "Executada por <subcontratada>":** etiqueta neutra (fundo `superficie`, borda `borda`, texto Aço, ícone `cable`), logo abaixo dos chips de estado. Não usa Sinal nem cor de status. Mostra ao lado o número da OS na conta da subcontratada, em `dado`. Leitor de tela: "Executada por Refrigeração Polar, subcontratada conectada. OS 2026-000311 na conta dela."
- **Campos somente leitura:** tudo o que é da execução (estado de execução, tarefas, horas, componentes, fotos, checklist, assinatura do cliente) aparece com o ícone `lock` no título do bloco e o texto "Registrado pela Refrigeração Polar". Não há botões de Iniciar, Pausar, Concluir, apontar horas ou registrar peça. Os campos não ficam esmaecidos: são legíveis com contraste normal; só não são editáveis.
- **O que o cliente faz:** Encerrar e Recusar liberação (com `os.encerrar`), Cancelar (o diálogo avisa "A Refrigeração Polar será avisada e a OS dela será marcada como cancelada no cliente") e editar descrição, prioridade e prazo enquanto a OS não estiver em execução.
- Linha do tempo com o prefixo da origem: "Iniciada por Ana Souza (Refrigeração Polar) às 08:12". Do técnico, só o nome.
- Sincronização com erro: a linha de sincronização fica em Atenção, "Não sincronizada há 35 min. Ver motivo", com link para a conexão (8.27).
- OS aberta pela subcontratada num equipamento do cliente aparece na lista como Aberta, com ORIGEM "Aberta por Refrigeração Polar" e na visão "Abertas sem planejamento", para a triagem normal.

**Na conta da subcontratada**, a mesma OS:

- Selo neutro "OS da Distribuidora Norte · OS-2026-004812" com o ícone `cable`.
- Descrição, prioridade, prazo e equipamento ficam somente leitura, com `lock` e "Definido pela Distribuidora Norte". A execução é editável normalmente.
- Não existe o botão Encerrar. Em Concluída, a faixa informativa diz "Aguardando a liberação da Distribuidora Norte. Ela foi avisada às 10:42." Quando o cliente recusa: faixa em Atenção "Liberação recusada pela Distribuidora Norte: falta foto do condensador." e a OS volta a Em execução.
- Quando o cliente encerra, a OS da subcontratada fica Encerrada com "Encerrada pela Distribuidora Norte em 13/11/2026" e libera o envio do relatório.

---

## 9. App de campo

### 9.1 Navegação por abas

Barra inferior com no máximo 5 abas, ícone de 24 dp + rótulo de 12 dp sempre visível, altura de 64 dp mais a área segura. O conjunto depende do perfil:

| Técnico | Técnico de subcontratada (SRV) | Motorista | Solicitante | Prestador |
|---|---|---|---|---|
| Hoje (minhas OS e disponíveis para assumir) | Hoje (minhas tarefas e disponíveis para assumir) | Meu veículo | Minhas solicitações | Hoje (OS do meu fornecedor atribuídas a mim) |
| Ativos | Equipamentos (dos clientes das minhas OS) | Checklist | **Ler QR** | Estoque em poder |
| **Ler QR** (central, destaque) | **Ler QR** | **Ler QR** | Perfil | **Ler QR** |
| Solicitações | Clientes (das minhas OS, com endereço e contato) | Ocorrências | | Mais |
| Mais (sincronização, perfil, ajuda, Sobre) | Mais (sincronização, contas, notificações, Sobre) | Mais | | |

O técnico de subcontratada é o perfil `tecnico` numa conta com SRV; as abas mudam pelos módulos da conta, não por um perfil novo. O locatário não usa o app no lançamento (8.19).

A aba central "Ler QR" é o único elemento em Sinal da barra: círculo de 56 dp com o ícone em Aço. Aba ativa: ícone e rótulo em `texto` com peso 600 e barra superior de 3 dp em Aço; inativa: `texto-secundario`. O portal do prestador está em 8.17. Vistorias de locação (8.16) entram pela aba Ativos do técnico ou por Ler QR.

**Tela principal (Hoje) e login do app com o selo de versão** (6.19):

```
Hoje                                   Login
┌───────────────────────────────┐      ┌───────────────────────────────┐
│ Hoje · 12/11      [cloud-check]│      │                               │
│ 3 OS para hoje                │      │      [símbolo 96 dp]          │
│ [OS-2026-004812]              │      │                               │
│ [OS-2026-004809]              │      │ E-mail [__________________]   │
│ [OS-2026-004790]              │      │ Senha  [__________________]   │
│                               │      │ [          Entrar          ]  │
│                               │      │ Entrar com Google · Microsoft │
│                         v1.4.2│      │                         v1.4.2│
├───────────────────────────────┤      │ (área segura inferior)        │
│ Hoje  Ativos  (QR)  Solic. Mais│     └───────────────────────────────┘
└───────────────────────────────┘
```

### 9.2 Alvos de toque e operação com luva

- Mínimo de 48 × 48 dp para qualquer alvo; 56 dp para a ação primária e para respostas de checklist.
- 8 dp de espaço mínimo entre alvos vizinhos.
- Itens de lista com altura mínima de 64 dp; o item inteiro é o alvo.
- Ações principais na metade inferior da tela, ao alcance do polegar.
- Nada exige gesto de pinça, deslizar como único caminho, toque duplo ou toque longo. Deslizar para ações em lista é atalho; a mesma ação existe em botão.
- Teclado numérico grande para medidores; botões "−" e "+" para quantidades de peças.
- Confirmação de ação crítica por botão grande, não por toque preciso em diálogo pequeno.

### 9.3 Alto contraste ao sol

- Tema **Sol**: fundo branco, texto Aço 900, texto secundário Aço 700, bordas Aço 500 de 2 dp, chips de status sólidos (seção 7.1). Atalho no menu Mais e no próprio indicador de sincronização.
- Pesos 500 ou mais em todo texto de dado; nenhum texto abaixo de 14 dp.
- Sem informação em cinza claro, sem texto sobre foto.
- Tema escuro não é recomendado ao sol; o app sugere o tema Sol se o usuário estiver no escuro e abrir a câmera durante o dia (sugestão única, dispensável).

### 9.4 Offline e sincronização

O app funciona 7 dias sem conexão sem perder registros (meta do sumário). A interface tem de deixar isso claro e calmo.

**Indicador de conexão e fila** — faixa abaixo do cabeçalho, sempre no mesmo lugar:

| Situação | Visual | Texto |
|---|---|---|
| Online, tudo enviado | Ícone `cloud-check` discreto no cabeçalho, sem faixa | (acessível) "Tudo enviado. Atualizado às 14:32." |
| Online, enviando | Faixa `superficie` com `refresh-cw` girando | "Enviando 3 itens" |
| Offline | Faixa `barraOffline` (Aço) com `cloud-off` | "Sem conexão · 12 itens salvos no aparelho" |
| Falha de envio | Faixa com fundo `status-atencao-fundo` e `triangle-alert` | "2 itens não foram enviados. Toque para ver." |
| Offline há muito tempo (5 dias) | Faixa Atenção | "Último envio há 5 dias. Conecte o aparelho para não passar do limite de 7 dias." |

Offline usa Aço, não vermelho nem âmbar: estar sem sinal é normal no campo e não é erro.

**Estado por registro** — ícone pequeno no canto de cada OS, checklist, leitura ou foto: `smartphone` (salvo no aparelho), `cloud-upload` (enviando), `cloud-check` (enviado), `triangle-alert` (falha). A tela **Sincronização** (em Mais) lista os itens pendentes por tipo, com horário de registro, tamanho das fotos e botão "Enviar agora".

**Conflito** — quando o mesmo registro mudou no servidor: "Esta OS foi alterada por Carlos Lima às 15:10. Revise antes de enviar." Tela de comparação com "Manter o meu" e "Usar o do servidor" por campo. Apontamentos, fotos e leituras nunca conflitam: são adicionados.

**Dados disponíveis offline** — o técnico vê o que está no aparelho: "Ativos da Planta Contagem e suas OS dos próximos 14 dias estão no aparelho." Itens fora desse escopo mostram "Disponível só com conexão".

### 9.5 Fluxo de OS em poucos toques

Meta: iniciar uma OS programada em 2 toques a partir da tela Hoje, e registrar uma corretiva a partir do QR em até 5 toques mais o texto.

```
Hoje                         OS-2026-004812              Executar
┌───────────────────┐  toque ┌───────────────────┐ toque ┌───────────────────┐
│ 3 OS para hoje    │ ─────► │ Troca de óleo...  │ ────► │ Tempo 00:12:40    │
│ [OS-2026-004812]  │        │ FRT-0112 Operando │       │ Tarefas 1/5       │
│ [OS-2026-004809]  │        │ Tarefas · Peças   │       │ [x] Drenar óleo   │
│ [OS-2026-004790]  │        │                   │       │ [ ] Trocar filtro │
│                   │        │ [   Iniciar OS   ]│       │ [Pausar][Concluir]│
└───────────────────┘        └───────────────────┘       └───────────────────┘
```

1. **Hoje**: lista das OS do técnico ordenadas por prioridade e horário, com a OS em andamento fixa no topo e o cronômetro visível. Abaixo, a seção "Disponíveis para assumir" com as OS `aberta` ou `planejada` da equipe sem responsável, cada uma com o botão "Assumir" (7.3).
2. **Iniciar OS** (primário, 56 dp): muda para Em execução e inicia o apontamento de horas.
3. **Execução**: tarefas como linhas tocáveis de 64 dp; peças por leitura de código ou busca com recentes; foto e voz em botões fixos.
4. **Pausar**: folha inferior com os motivos como botões grandes (um toque); "Outro" abre campo de texto.
5. **Concluir**: tela de resumo com o que falta destacado no topo (tarefas obrigatórias, fotos mínimas, código de falha em corretiva, leitura do medidor). Código de falha com os 5 mais usados para a classe do ativo primeiro. Assinatura quando o plano exigir. Botão "Concluir OS".

Abrir corretiva a partir do QR: Ler QR → ativo → "Registrar falha" → descrição por voz ou texto + foto → "Ativo parou?" (Sim/Não em botões grandes) → "Enviar solicitação".

### 9.6 Checklist pré-viagem do motorista

```
Checklist pré-viagem · FRT-0112                            Item 4 de 18
┌──────────────────────────────────────────────────────────────┐
│ ████████████░░░░░░░░░░░░░░░░░ (progresso)                    │
│                                                              │
│ Pneus sem cortes ou bolhas aparentes                         │  titulo3
│ Verifique os 10 pneus e o estepe.                            │  corpo
│ [foto de referência]                                         │
│                                                              │
│ ┌──────────────────────┐  ┌──────────────────────┐           │
│ │      ✓  OK           │  │    ✕  Não OK         │  56 dp    │
│ └──────────────────────┘  └──────────────────────┘           │
│ [ Anterior ]                                [ Pular ]        │
└──────────────────────────────────────────────────────────────┘
```

1. **Veículo**: o veículo atribuído ao motorista aparece pronto; outro veículo por QR ou placa.
2. **Hodômetro**: campo grande com teclado numérico e a última leitura como referência. Valor menor que a última leitura é bloqueado com explicação.
3. **Itens**: um item por tela (modo guiado, padrão) ou lista completa (modo rápido, para quem já conhece). Respostas "OK" e "Não OK" em botões de 56 dp lado a lado, com ícone e texto; a resposta selecionada fica com borda de 2 dp em Aço e fundo de status tingido.
4. **Não OK**: pede foto (obrigatória se o modelo exigir) e descrição por voz ou texto. Gera solicitação automaticamente, como define o módulo CHK.
5. **Item crítico reprovado**: tela de bloqueio com fundo `status-parado-fundo`, ícone `octagon-x` e texto direto: "Não saia com o veículo. Freio reprovado. O gestor da frota foi avisado às 06:42." Ação: "Ligar para o gestor" (se cadastrado).
6. **Resumo e envio**: "18 itens verificados · 1 com problema". Assinatura opcional. "Enviar checklist" funciona offline.

### 9.7 Captura de foto e voz

- **Foto**: câmera dentro do app, botão de disparo de 72 dp, lanterna, alternância de câmera; após a captura, "Usar foto" e "Tirar outra". Anotação simples opcional (seta e círculo em Aço com borda branca). Data, hora e TAG gravadas nos metadados, não sobre a imagem.
- **Voz**: botão de microfone de 56 dp. Toque para começar e toque para parar (não exige segurar). Onda de nível em Aço e cronômetro em `dado`. Limite de 2 min por gravação.
- O áudio é salvo no aparelho e anexado ao registro. A transcrição é feita pelo assistente (IA) quando houver conexão e aparece como texto editável marcado "Transcrito automaticamente · confira". Offline: "Áudio salvo. A transcrição aparece quando o aparelho se conectar."
- Ditado do teclado do sistema continua disponível em todo campo de texto.

### 9.8 App do técnico da subcontratada (SRV)

Mesmo Tenaz Mobile, mesmas regras de toque, contraste e offline (9.2 a 9.4). A conta com SRV acrescenta à OS o cliente atendido, o endereço de atendimento, o serviço a realizar, os componentes com número de série e a assinatura do cliente.

**Lista de tarefas (Hoje):**

```
┌───────────────────────────────┐
│ Hoje · 13/11      [cloud-check]│
│ Refrigeração Polar         [▾]│  legenda; seletor só com 2+ contas
│ EM ANDAMENTO                  │  sobrelinha
│ ┌───────────────────────────┐ │
│ │ OS-2026-000311   00:42:10 │ │  dado-forte
│ │ Distribuidora Norte  ⇄    │ │
│ │ CFR-0021 Câmara fria 2    │ │
│ └───────────────────────────┘ │
│ A FAZER (3)                   │
│ ┌───────────────────────────┐ │
│ │ OS-2026-000314 · 10:00    │ │
│ │ Frigorífico Sul · Alta    │ │
│ │ CMP-0012 · Rua B, 200     │ │
│ │ [Como chegar]             │ │
│ └───────────────────────────┘ │
│ DISPONÍVEIS PARA ASSUMIR (2)  │
│ ┌───────────────────────────┐ │
│ │ OS-2026-000316 · Vence hoje│ │
│ │ Padaria Pão Bom · PPB-0001│ │
│ │ [        Assumir        ] │ │
│ └───────────────────────────┘ │
│                         v1.4.2│
├───────────────────────────────┤
│ Hoje  Equip. (QR) Clientes Mais│
└───────────────────────────────┘
```

- Ordem: em andamento (fixa no topo, com cronômetro), a fazer por horário e prioridade, disponíveis para assumir. Cada cartão mostra cliente, equipamento, endereço e prazo; "Como chegar" abre o app de mapas do aparelho.
- "Assumir" é um botão de 56 dp no estilo secundário, porque há vários cartões na lista (regra 1 de 2.2). Offline, a ação fica pendente e a OS mostra "Assumida · aguardando envio"; se outro técnico tiver assumido antes, a sincronização devolve "Esta OS já foi assumida por João Lima." e a remove.

**Execução:** as abas da OS são Tarefas, Serviço, Componentes, Fotos e Assinatura, com [Pausar] e [Concluir] fixos embaixo, como em 9.5.

```
Componentes · OS-2026-000311            Registrar componente
┌───────────────────────────────┐       ┌───────────────────────────────┐
│ SUBSTITUÍDOS (1)              │       │ Tipo                          │
│ Motor do ventilador · 1 un    │       │ (•) Substituído               │
│ VT-88213 → VT-90455           │       │ ( ) Consumido                 │
│ CONSUMIDOS (1)                │       │ ( ) Fornecido pelo cliente    │
│ Fita isolante · 1 un          │ toque │ Item [Motor do ventilador ▾]  │
│ DO CLIENTE (0)                │ ────► │ Quantidade  [−]  1  [+]       │
│                               │       │ PEÇA RETIRADA                 │
│                               │       │ Nº de série [VT-88213 ] [scan]│
│                               │       │ PEÇA INSTALADA                │
│ [ + Registrar componente ]    │       │ Nº de série [VT-90455 ] [scan]│
├───────────────────────────────┤       │ [   Registrar componente   ]  │
│ [ Pausar ]      [ Concluir ]  │       └───────────────────────────────┘
└───────────────────────────────┘
```

- **Tipo** (`tipo_material`) em botões de rádio grandes: Substituído (pede série retirada e instalada quando o item controla série), Consumido (só quantidade) e Fornecido pelo cliente (quantidade e descrição; não baixa estoque da subcontratada).
- **Leitura do número de série:** o botão `scan-barcode` abre a câmera para código de barras ou QR da etiqueta. O valor lido aparece grande, em `dado-kpi`, para o técnico conferir antes de aceitar ("Usar VT-90455" ou "Ler de novo"). Digitação manual sempre disponível, com teclado alfanumérico em caixa alta. Avisos sem bloqueio: série instalada igual à retirada; série já registrada como instalada em outro equipamento. A leitura de texto impresso por OCR fica como proposta, a validar em campo.
- **Item:** busca no catálogo de itens com os recentes primeiro; item que não existe pode ser registrado como texto livre, marcado para o escritório revisar.
- **Fotos:** a aba Fotos sugere as posições "Antes", "Depois", "Peça retirada" e "Placa do equipamento" como uma grade com contorno tracejado, como na vistoria (8.16); a quantidade mínima vem do serviço ou da empresa. Cada foto pode ser marcada "Incluir no relatório".

**Assinatura do cliente e conclusão:**

```
Assinatura do cliente · OS-2026-000311
┌───────────────────────────────┐
│ Resumo                        │
│ Troca do ventilador · 2,5 h   │
│ 2 componentes · 6 fotos       │
│ "Declaro que o serviço        │
│  descrito foi executado."     │
│ Nome    [José Pereira______]  │
│ Função  [Encarregado_______]  │
│ ┌───────────────────────────┐ │
│ │  área de assinatura 120 dp│ │
│ └───────────────────────────┘ │
│ [ Limpar ]                    │
│ [ ] Cliente ausente           │
│ [       Concluir OS        ]  │  56 dp, Sinal
└───────────────────────────────┘
```

- O resumo vem antes da assinatura, para o cliente ler o que assina: serviço, horas, componentes e fotos. "Ver detalhes" abre a lista completa.
- "Cliente ausente" troca a área de assinatura por um motivo obrigatório (Ninguém no local, Recusou assinar, Outro) e uma foto do local; o relatório mostra "Sem assinatura do cliente: ninguém no local".
- "Concluir OS" só fica ativo com o que for obrigatório; a dica diz o que falta ("Falta a foto Depois"), como em 9.5.
- Tudo funciona offline. A conclusão fica no aparelho até o envio; o relatório PDF é gerado no servidor depois da sincronização, e o app mostra "Relatório disponível quando o aparelho se conectar".

### 9.9 Troca de conta e notificações de tarefa nova

**Troca de conta** (técnico que trabalha para mais de uma empresa na Tenaz):

```
Trocar de conta
┌───────────────────────────────┐
│ Refrigeração Polar        ●   │  conta ativa
│ 3 tarefas · tudo enviado      │
├───────────────────────────────┤
│ Clima Frio Serviços           │
│ 1 tarefa nova · 2 itens não   │
│ enviados                      │
├───────────────────────────────┤
│ [ + Entrar em outra conta ]   │
└───────────────────────────────┘
```

- O seletor só aparece com duas ou mais contas no aparelho: o nome da conta ativa fica sob o título da tela Hoje, com o ícone `chevrons-up-down`, e abre uma folha inferior com as contas, a contagem de tarefas novas e as pendências de envio de cada uma.
- Trocar de conta mostra por um instante a tela "Abrindo Clima Frio Serviços" e recarrega Hoje. Nenhuma lista mistura contas, e o nome da conta ativa fica sempre visível no cabeçalho. A cor do app não muda por conta.
- As contas sincronizam em segundo plano; a faixa de sincronização (9.4) é da conta ativa, e o seletor mostra as pendências das outras.
- Sair de uma conta (Mais > Contas) avisa o que ainda não foi enviado: "2 itens de Clima Frio Serviços ainda não foram enviados. Conecte o aparelho antes de sair ou eles serão perdidos." Botões [Voltar] [Sair mesmo assim].

**Notificações push:**

| Situação | Texto da notificação |
|---|---|
| OS atribuída | "Nova tarefa: OS-2026-000314" (com o nome da conta quando houver mais de uma: "Refrigeração Polar · Nova tarefa: OS-2026-000314") |
| Várias atribuições | "3 tarefas novas" |
| OS retirada do técnico | "OS-2026-000314 foi passada para outro técnico" |
| Liberação recusada | "OS-2026-000311 voltou para execução. Veja o motivo" |
| Prazo próximo | "OS-2026-000314 vence em 2 h" |

- O texto não leva nome de cliente pessoa física, endereço nem descrição do serviço; o detalhe aparece só dentro do app.
- Tocar na notificação abre a conta certa (trocando, se preciso) e a OS. Se a OS ainda não estiver no aparelho, a tela mostra "Buscando a OS-2026-000314" enquanto sincroniza; sem conexão, "A tarefa aparece aqui quando o aparelho se conectar."
- Em Mais > Notificações, o técnico escolhe por conta: receber sempre ou silenciar fora do horário de trabalho.

---

## 10. Microcopy e tom de voz

### 10.1 Princípios

Do brandbook: direto, técnico sem jargão, calmo nos alertas. Diz o que fazer e quando.

1. **Comece pelo fato.** Objeto, situação, prazo: "Preventiva do FRT-0112 vence em 2.300 km."
2. **Diga o próximo passo.** Toda mensagem de erro termina numa ação possível.
3. **Use o vocabulário do produto.** Ativo, OS, solicitação, plano, leitura, técnico. Nunca "ticket", "chamado", "item" genérico ou "registro" quando há nome melhor.
4. **Use números com unidade.** "3 peças abaixo do estoque mínimo", não "algumas peças".
5. **Sem exclamação, sem caixa alta, sem emoji, sem culpa.** Não "Você digitou errado", e sim "Leitura menor que a anterior".
6. **Frases curtas.** Um alerta é uma frase. Uma confirmação, duas no máximo.
7. **Trate o usuário por "você" só quando necessário.** Prefira o imperativo ("Informe o motivo da pausa").
8. **Caixa de frase** em títulos, botões e menus: "Nova ordem de serviço", não "Nova Ordem de Serviço".

### 10.2 Termos

| Use | Não use | Por quê |
|---|---|---|
| Encerrar OS | Fechar, finalizar | `encerrada` é o estado canônico |
| Concluir OS | Terminar | `concluida` é o estado canônico, diferente de encerrar |
| Cancelar OS | Excluir OS | OS não é excluída; é cancelada com motivo |
| Inativar ativo | Excluir ativo | Ativo com histórico não é excluído |
| Solicitação | Chamado, ticket | Entidade `solicitacao` |
| Leitura | Medição, apontamento de km | Entidade `leitura_medidor`; "apontamento" é de horas |
| Voltar | Cancelar (em diálogos) | Evita ambiguidade com "Cancelar OS" |
| Sem conexão | Offline, erro de rede | Linguagem do usuário de campo |
| Urgente (prioridade) | Crítica, Emergência | `urgente` é o valor canônico; "crítica" confunde com criticidade e "emergencial" é tipo de OS |
| Reconhecer OS | Aceitar, dar ciência | Ação canônica `reconhecer`, grava `data_reconhecimento` |
| Recusar liberação | Reabrir, reprovar OS | Volta a OS de Concluída para Em execução, com motivo |
| Fornecedor | Terceiro, parceiro, empresa externa | Entidade `fornecedor`; o papel diz o resto (Subcontratada, Assistência técnica, Oficina, Locador) |
| Prestador | Terceirizado | Perfil `prestador`: usuário de um fornecedor |
| Ordem de reparo | Conserto, OS de peça | Entidade `ordem_reparo`, número OR- |
| Remessa a terceiro / Retorno de terceiro | Envio, empréstimo | Tipos canônicos de movimentação |
| Locatário | Cliente da locação, inquilino | Entidade `locatario` |
| Cedido / Tomado | Alugado para / alugado de | Sentidos canônicos do contrato de locação |
| Vistoria de entrega / de devolução | Check-in, check-out | Entidade `vistoria_locacao` |
| Assumir OS | Pegar, puxar, aceitar | Ação canônica `assumir`: o técnico se atribui a OS e grava `data_reconhecimento` |
| Encerrar contrato | Cancelar contrato (quando vigente) | Contrato vigente termina por `encerrar`; só o rascunho é cancelado |
| Reabrir como nova | Desfazer duplicada | `duplicada` é final; a correção cria outra solicitação |
| Cliente atendido | Cliente final, contratante | Entidade `cliente_servico`, na conta da subcontratada |
| Equipamento de cliente | Ativo do cliente, máquina do cliente | `ativo` com `propriedade = de_cliente` |
| Serviço a realizar | Escopo, chamado | `descricao_servico` e itens do catálogo na OS |
| Componente substituído / consumido / fornecido pelo cliente | Peça trocada, material | Valores de `tipo_material` |
| Distribuir (OS aos técnicos) | Despachar | "Despacho" já é termo de frota ("Apto para despacho") |
| Conexão, conectada | Integração, vínculo, parceria | Entidade `conexao_empresa`; "Integrações" é o módulo INT |
| Executada por (subcontratada) | Terceirizada para | Selo da OS vinculada (8.28) |
| Conta | Empresa, tenant, ambiente | Palavra do app para cada empresa em que o usuário trabalha (9.9) |

### 10.3 Mensagens de erro

Formato: **o que aconteceu** + **o que fazer**. Dados do contexto sempre que possível.

| Situação | Evite | Prefira |
|---|---|---|
| Leitura inválida | "Valor inválido." | "Hodômetro menor que a última leitura (248.315 km em 10/11/2026). Confira o valor ou registre a troca do painel." |
| Campo obrigatório | "Campo obrigatório!" | "Informe o motivo da pausa." |
| Sem conexão | "Erro de rede. Tente novamente mais tarde." | "Sem conexão. A OS foi salva no aparelho e será enviada quando houver sinal." |
| Permissão | "Acesso negado." | "Seu perfil não pode encerrar OS. Peça a liberação ao planejador da unidade." |
| Erro do servidor | "Ocorreu um erro inesperado." | "Não foi possível salvar agora. Seus dados continuam no formulário. Tente de novo em alguns minutos; se continuar, informe o código E-4F21 ao suporte." |
| Conflito | "Registro desatualizado." | "Esta OS foi alterada por Ana Souza às 14:20. Revise as diferenças antes de salvar." |
| Saldo | "Estoque insuficiente." | "Saldo de 2 un. no Almoxarifado Central. Reduza a quantidade ou faça uma requisição." |
| Capacitação | "Técnico inválido." | "João Lima não tem NR-10 válida (venceu em 03/11/2026). Escolha outro técnico ou atualize a capacitação." |
| Importação | "Arquivo com erros." | "4 de 316 linhas não foram importadas. Baixe o relatório para ver o motivo de cada uma." |
| Sessão | "Sessão expirada." | "Sua sessão terminou por segurança. Entre de novo; o que você digitou foi mantido." |
| Campo de outra conta | "Campo bloqueado." | "A descrição é definida pela Distribuidora Norte. Peça a alteração ao cliente." |
| Encerrar OS vinculada (subcontratada) | "Ação não permitida." | "A liberação desta OS é feita pela Distribuidora Norte. Ela foi avisada da conclusão às 10:42." |
| Convite de conexão | "Convite inválido." | "Este convite venceu ou foi enviado para outro CNPJ. Peça um novo à Distribuidora Norte." |
| Conexão suspensa | "Erro de sincronização." | "A conexão com a Distribuidora Norte está suspensa. As mudanças desta OS ficam guardadas e serão enviadas quando a conexão for reativada." |
| Número de série | "Valor inválido." | "O número de série instalado é igual ao retirado (VT-88213). Confira a etiqueta da peça nova." |
| OS já assumida | "Conflito." | "Esta OS já foi assumida por João Lima às 09:02." |
| Peça em trânsito | "Saldo indisponível." | "O filtro PSL-118 está em trânsito. Confirme o recebimento antes de usar na OS." |

### 10.4 Confirmações e sucesso

| Situação | Texto |
|---|---|
| OS encerrada (toast) | "OS-2026-004812 encerrada. Ativo liberado para operação." |
| OS criada | "OS-2026-004830 criada para CMP-0347." [Abrir] |
| Programação publicada | "Programação da semana 46 publicada. 6 técnicos notificados." |
| Leitura registrada | "Leitura de 248.315 km registrada para FRT-0112." |
| Ação desfeita | "Foto removida." [Desfazer] |
| Pergunta de ação | "3 peças abaixo do estoque mínimo. Gerar requisição?" |
| Confirmação de cancelamento | Título: "Cancelar OS-2026-004812?" · Texto: "A OS não será executada e o motivo ficará no histórico." · Botões: [Voltar] [Cancelar OS] |
| Confirmação de inativação | Título: "Inativar FRT-0112?" · Texto: "Os planos do ativo param de gerar OS. O histórico continua disponível." · Botões: [Voltar] [Inativar ativo] |
| OS assumida | "OS-2026-000316 assumida. Ela está na sua lista de hoje." |
| Distribuição publicada | "Distribuição de 13/11 publicada. 3 técnicos notificados." |
| Convite de conexão enviado | "Convite enviado a joao@polar.com.br. A conexão fica ativa quando a Refrigeração Polar aceitar." |
| Conexão aceita | "Conexão com a Distribuidora Norte ativa. 42 equipamentos recebidos." |
| Relatório enviado | "Relatório da OS-2026-000311 enviado a jose@norte.com.br." |
| Componente registrado | "Motor do ventilador registrado: VT-88213 retirado, VT-90455 instalado." [Desfazer] |

Evite: "Sucesso!", "Tudo certo por aqui", "Operação realizada com sucesso".

### 10.5 Estados vazios

| Tela | Texto |
|---|---|
| Ativos, primeiro uso | "Nenhum ativo cadastrado. Cadastre o primeiro ou importe uma planilha com TAG, nome e localização." |
| OS do técnico, dia livre | "Nenhuma OS programada para hoje. As novas OS aparecem aqui assim que forem programadas." |
| Preventivas atrasadas | "Nenhuma preventiva atrasada." |
| Vencimentos | "Nenhum documento vence nos próximos 30 dias." |
| Busca | "Nenhum resultado para \"FRT-0199\". Confira a TAG ou busque pela placa." |
| Solicitações do solicitante | "Você ainda não abriu solicitações. Leia o QR de um ativo para registrar um problema." |
| Clientes atendidos, primeiro uso | "Nenhum cliente cadastrado. Cadastre o primeiro cliente e os equipamentos que você atende, ou importe uma planilha." |
| Técnicos, primeiro uso | "Nenhum técnico cadastrado. Cadastre os técnicos e envie o convite para o Tenaz Mobile." |
| Distribuição sem OS | "Nenhuma OS sem técnico. As novas OS, inclusive as que chegam de clientes conectados, aparecem aqui." |
| Conexões | "Nenhuma conexão. Quando um cliente que usa a Tenaz convidar sua empresa, o convite aparece aqui." |
| App, disponíveis para assumir | (a seção não aparece quando vazia) |

### 10.6 Alertas e notificações

| Situação | Texto |
|---|---|
| Plano vencendo | "Preventiva do FRT-0112 vence em 2.300 km." |
| Plano vencido | "Preventiva do CMP-0347 venceu há 40 h. Programe a OS-2026-004801." |
| Documento | "CRLV-e do FRT-0112 vence em 12 dias (24/11/2026)." |
| Leitura fora da faixa | "Temperatura do CMP-0347 em 104 °C, acima da faixa de 60 a 95 °C." |
| Ativo parado | "CMP-0347 parado há 6 h. OS-2026-004809 aguardando peça." |
| Multa | "Indique o condutor da multa de FRT-0112 até 20/11/2026." |
| Link de execução aberto | "Link da OS-2026-004812 aberto por um aparelho em Contagem/MG às 11:02." |
| Execução por link recebida | "Refrigeração Polar concluiu a OS-2026-004812 pelo link. Revise e encerre." |
| Ordem de reparo atrasada | "OR-2026-000045 passou do prazo previsto (20/11). Cobre a Hidráulica Sul." |
| SLA de fornecedor (Fase 2) | "Refrigeração Polar não atendeu a OS-2026-004812 no prazo de 4 h do contrato." |
| Contrato de locação terminando | "CL-2026-000012 termina em 9 dias (21/11/2026). Agende a vistoria de devolução." |
| Avaria na devolução | "Avaria no ESC-0031 registrada na devolução. OS-2026-004851 criada com custo do locatário." |
| Nova tarefa (push, técnico) | "Nova tarefa: OS-2026-000314" (textos completos em 9.9) |
| OS concluída pela subcontratada (cliente) | "Refrigeração Polar concluiu a OS-2026-004812. Revise e encerre." |
| OS aberta pela subcontratada (cliente) | "Refrigeração Polar abriu a OS-2026-004860 no CMP-0347. Faça a triagem." |
| Liberação recusada (subcontratada) | "Distribuidora Norte recusou a liberação da OS-2026-000311: falta foto do condensador." |
| Convite de conexão (subcontratada) | "A Distribuidora Norte convidou sua empresa para conectar. Responda até 26/11/2026." |
| OS não sincronizada (as duas contas) | "A OS-2026-000298 não foi sincronizada com a Distribuidora Norte há 35 min. Veja o motivo." |
| Conexão encerrada (as duas contas) | "A conexão com a Refrigeração Polar foi encerrada por Carlos Lima. O histórico continua disponível." |

---

## 11. Acessibilidade

Meta: WCAG 2.1 nível AA no web e no app, verificada em cada componente antes de entrar na biblioteca.

### 11.1 Foco visível

- Anel de 2 px em `foco` (Aço no claro, `#E8EDF2` no escuro), afastado 2 px, em todo elemento focável, via `:focus-visible`. Sobre o menu lateral e superfícies Aço, o anel é Névoa.
- Nunca `outline: none` sem substituto.
- O foco nunca fica escondido atrás do cabeçalho fixo (`scroll-margin-top` igual à altura do cabeçalho).
- Ao abrir modal, o foco entra; ao fechar, volta ao elemento que abriu. Ao excluir um item de lista, o foco vai para o item seguinte.

### 11.2 Navegação por teclado

| Tecla | Ação |
|---|---|
| Tab / Shift + Tab | Próximo e anterior elemento |
| `/` ou Ctrl + K | Busca global |
| Esc | Fecha menu, popover, modal; limpa a busca |
| Setas | Navegam em menus, abas, árvores, grade da tabela e quadro de programação |
| Enter | Abre o item focado; confirma |
| Espaço | Marca checkbox, seleciona linha |
| Shift + F10 | Menu de ações da linha |
| `?` | Lista de atalhos |

- Link "Pular para o conteúdo" como primeiro elemento da página.
- Toda operação de arrastar (kanban, programação, rodízio de pneus, ordem de itens do checklist) tem alternativa por teclado e por menu.
- Ordem de tabulação segue a ordem visual.

### 11.3 Leitores de tela

- Marcos: `header`, `nav` (com `aria-label="Menu principal"`), `main`, `aside`. Um `h1` ou `h2` de página por tela, hierarquia sem saltos.
- Ícones decorativos com `aria-hidden="true"`; botões só de ícone com `aria-label`.
- Chip de status lido como texto: "Status: Pausada, aguardando peça". Criticidade: "Criticidade A".
- Números lidos com unidade por extenso quando a abreviação for ambígua: `<abbr title="horas">h</abbr>` ou texto oculto ("12.480 horas").
- Tabelas com `th` e `scope`, `aria-sort` nos cabeçalhos ordenáveis, `caption` oculto descrevendo a lista e o filtro.
- Mudanças assíncronas anunciadas em regiões `aria-live`: resultado de filtro ("128 OS encontradas"), toasts, estado de sincronização.
- Gráficos com `role="img"` e `aria-label` com o resumo ("Disponibilidade de 97,4 % em novembro, acima da meta de 96 %") e a tabela de dados como alternativa.
- No app: `accessibilityLabel`, `accessibilityRole`, `accessibilityState` (`selected`, `disabled`, `busy`, `checked`) e `accessibilityHint` nas ações não óbvias. Teste com TalkBack e VoiceOver em cada fluxo da seção 9.

### 11.4 Não depender só de cor

- Status: cor + ícone + texto (seção 7).
- Erro em campo: borda + ícone + mensagem.
- Gráficos: tracejado diferente por série, legenda direta, valor no ponto final.
- Seleção em tabela: fundo + barra lateral + checkbox marcado.
- Itens obrigatórios e opcionais: texto, não cor.
- Links no meio de texto: sublinhados.

### 11.5 Outros critérios

| Critério | Regra |
|---|---|
| 1.4.3 / 1.4.11 Contraste | Tabelas da seção 4; nenhum par novo sem cálculo |
| 1.4.4 Redimensionar texto | Interface funcional com zoom de 200 % no navegador e fonte do sistema em 200 % no app |
| 1.4.10 Refluxo | Páginas funcionam em 320 px de largura sem rolagem horizontal, exceto tabelas e quadros (conteúdo bidimensional permitido pela norma) |
| 1.4.12 Espaçamento de texto | Layout não quebra com espaçamento de linha 1,5 e de letras 0,12 em |
| 1.4.13 Conteúdo em hover | Tooltip dispensável com Esc, não some ao mover o ponteiro sobre ele |
| 2.2.1 Tempo ajustável | Sessão avisa 2 min antes de expirar e permite estender; toasts de erro não somem sozinhos |
| 2.5.3 Rótulo no nome | O nome acessível começa com o texto visível |
| 2.5.5 Tamanho do alvo (AAA, adotado no app) | 48 dp no app; 32 px no web com área clicável de 40 px em tabelas |
| 3.1.1 Idioma | `lang="pt-BR"` no documento |
| 3.3.1 a 3.3.4 Erros | Mensagens da seção 10.3; confirmação para ações irreversíveis |

Ferramentas: axe-core nos testes de componente e de ponta a ponta, eslint-plugin-jsx-a11y, verificação manual com NVDA (Windows), VoiceOver (macOS e iOS) e TalkBack (Android) a cada release.

---

## 12. Iconografia

### 12.1 Biblioteca

**Lucide** (licença ISC): `lucide-react` no web e `lucide-react-native` no app (depende de `react-native-svg`). O estilo coincide com o da marca: grade de 24, traço de 2 px, pontas e junções arredondadas, sem preenchimento.

| Tamanho | Uso | Traço |
|---|---|---|
| 16 px | Chips, tabelas densas, botões P | `absoluteStrokeWidth` com 1,75 px, para não ficar fino demais |
| 20 px | Botões M, campos, menu lateral | 2 px na grade de 24 (padrão) |
| 24 px | App, barra de abas, botões G | 2 px |
| 32 a 48 px | Estados vazios, leitor de QR | 2 px na grade de 24 (fica proporcional, como no brandbook) |

Cor do ícone: `texto` ou `texto-secundario`; cor de status só em ícone de status. Ícone nunca aparece em Sinal, exceto sobre o botão primário, onde é Aço.

### 12.2 Mapeamento de ícones do domínio

| Conceito | Ícone Lucide |
|---|---|
| Frota, veículo (brandbook: Frota) | `truck` |
| Planta, equipamento (brandbook: Planta) | `factory` |
| Ordem de serviço (brandbook) | `clipboard-list` |
| Preventiva (brandbook) | `calendar-check` |
| Medidor, leitura (brandbook: Medição) | `gauge` |
| Peças, estoque (brandbook) | `package` |
| Checklist (brandbook) | `list-checks` |
| Indicadores (brandbook) | `chart-line` |
| Solicitação | `message-square-warning` |
| Programação | `calendar-range` |
| Plano de manutenção | `repeat` |
| Localização | `map-pin` |
| Almoxarifado | `warehouse` |
| Compras | `shopping-cart` |
| Fornecedor (qualquer papel), executor fornecedor | `building-2` |
| Subcontratada | `hard-hat` |
| Assistência técnica | `badge-check` |
| Prestador (usuário) | `user-cog` |
| Link de uso único | `link` |
| Remessa e retorno de terceiro | `arrow-right-left` |
| Ordem de reparo | `hammer` |
| Quarentena, irreparável | `package-x` |
| Peça serializada | `barcode` |
| Locação, contrato de locação | `key-round` |
| Locatário | `handshake` |
| Cliente atendido (SRV) | `briefcase-business` |
| Equipamento de cliente | `factory` com o selo "De cliente" |
| Técnico, especialidade | `user-round` / `award` |
| Distribuição de OS | `columns-3` |
| Componente substituído | `replace` |
| Leitura de número de série | `scan-barcode` |
| Relatório PDF da OS | `file-text` |
| Conexão entre contas, OS vinculada | `cable` |
| Campo somente leitura de outra conta | `lock` |
| Seletor de conta (app) | `chevrons-up-down` |
| Notificação | `bell` |
| Vistoria de locação | `clipboard-check` |
| Responsável pelo custo | `wallet` |
| Reconhecimento da OS | `bell-ring` |
| Prioridade urgente | `siren` |
| Combustível, abastecimento | `fuel` |
| Multa | `receipt-text` |
| Sinistro | `shield-alert` |
| Conformidade | `shield-check` |
| Documento com vencimento | `file-clock` |
| Capacitação | `graduation-cap` |
| Empilhadeira | `forklift` |
| Câmara fria, refrigeração | `snowflake` |
| Ler QR | `scan-qr-code` |
| Foto | `camera` |
| Voz | `mic` |
| Assinatura | `signature` |
| Sem conexão | `cloud-off` |
| Enviando | `cloud-upload` |
| Enviado | `cloud-check` |
| Salvo no aparelho | `smartphone` |
| Assistente (IA) | `sparkles` |
| Administração | `settings` |
| Integrações | `plug` |

Os ícones de status seguem a seção 7. Os nomes acima correspondem à versão atual do Lucide; ao atualizar a biblioteca, confira renomeações no changelog.

### 12.3 Ícones próprios necessários

Conceitos que o Lucide não cobre ou cobre de forma genérica. Desenhar na mesma grade de 24, traço de 2 px, pontas e junções arredondadas, área útil de 20 px (margem de 2 px), sem preenchimento, raio de canto de 2 px. Entregar em SVG em `marca/icones/` e como componentes no pacote de UI.

| Ícone | Uso |
|---|---|
| Símbolo Tenaz (porca, T e ponto) | Logo; não é ícone de função e não pode ser redesenhado |
| Cavalo mecânico | Distinguir de caminhão rígido |
| Carreta / semirreboque | Implemento rodoviário |
| Pneu com posição | Módulo de pneus |
| Diagrama de eixo (2, 3 e 4 pneus por eixo) | Mapa de pneus |
| Sulco do pneu | Leitura de profundidade |
| Horímetro | Medidor de horas (diferente de relógio) |
| Hodômetro | Medidor de km (diferente de velocímetro) |
| Cronotacógrafo | Documento e aferição |
| Vaso de pressão / caldeira | NR-13 |
| Proteção de máquina | NR-12 |
| Bloqueio e etiquetagem (LOTO) | Cadeado com etiqueta |
| Arla 32 | Diferenciar de diesel |
| Compressor | Classe de ativo comum |
| Gerador | Classe de ativo comum |
| Código de falha (modo, causa, atividade) | Catálogo ISO 14224 |
| Ativo pai e filhos | Hierarquia de ativos |

Regras: ícone próprio passa pela mesma revisão de contraste e legibilidade a 16 px; nenhum ícone usa a porca sextavada, para não confundir com o símbolo da marca.

---

## 13. Inconsistências e pendências

Encontradas ao cruzar a marca (`marca/README.md`, `brandbook.html`) com o sumário executivo. As decisões provisórias deste documento estão indicadas; todas precisam de validação. A numeração da versão 0.1 foi mantida para não quebrar referências.

**Resolvidas pelo sumário executivo 0.4** (aplicadas nesta versão)

| Tema | Como ficou |
|---|---|
| Permissão de liberação | `os.encerrar`, padrão em gestor e planejador, nunca em prestador ou locatário (7.3, 8.5) |
| Ação Assumir | Botão "Assumir" no web e no app; grava `data_reconhecimento` (7.3, 9.5, 9.8) |
| Link de uso único | Endereço `/l/#<token>`, código trocado por sessão curta, envio único que registra e conclui (8.18) |
| Solicitação | Desistência de Nova para Recusada; `duplicada` final com "Reabrir como nova"; origem "Portal do locatário" (7.5, 8.19) |
| Contrato de locação | Cancelar só no rascunho; vigente termina por "Encerrar contrato"; sem sublocação; ativo de terceiro devolvido fica Inativo (7.5, 8.15) |
| Campos da OS e da locação | `justificativa_responsavel_custo`, `motivo_recusa_liberacao`, rótulos de `origem` (inclui `prestador`), `local_uso`, `franquia_horas`, `franquia_km` (7.3, 8.15) |
| Papéis de fornecedor | Posto de combustível e Seguradora (8.13) |
| Peça em trânsito e quarentena | Peça em trânsito não é consumível; "Decidir destino" só com `est.decidir_destino` (8.9) |
| Locatário | Só portal web no lançamento (8.19, 9.1) |
| Versão visível e histórico de versões | RNF-039 e RNF-040 citados no selo de versão (6.19) |
| Módulo SRV | Menu (5.2), estados de conexão e de sincronização (7.5), telas 8.20 a 8.28, app 9.8 e 9.9, microcopy (10) e ícones (12.2) |

**Resolvidas pelo sumário executivo 0.3**

| # | Inconsistência | Como foi resolvida |
|---|---|---|
| 13.1 | O sumário tinha 4 estados do ativo e a marca 5 cores; Atenção não era estado gravado e "Preventiva" não correspondia a `em_manutencao`. | Sumário 3.3: Atenção é condição calculada sobre `operando`, com a cor de Atenção; "Em manutenção" usa a cor de Preventiva. Aplicado em 7.2. Resta só a decisão de marca, opcional, de renomear a cor "Preventiva" no brandbook (ver 13.15). |
| 13.2 | Não estava definido se corretiva planejada em execução deixa o ativo `em_manutencao` ou `parado`. | Sumário 3.3: `em_manutencao` quando a retirada de operação é programada (inclui corretiva planejada); `parado` quando houve falha que parou o ativo, corretiva emergencial em curso, parecer NR-12 `inapto` ou bloqueio; `parado` prevalece; volta a `operando` no encerramento. Aplicado em 7.2. |
| 13.9 | O brandbook abrevia o número de OS (`OS-004812`). | Sumário 3.4 fixa `OS-`, `SS-`, `OR-` e `CL-` + `AAAA-NNNNNN`, por empresa, reiniciando por ano. A interface usa sempre o formato completo (3.11). Corrigir os exemplos do brandbook continua como tarefa de marca (13.15). |
| 13.11 | O sumário não definia os valores de prioridade. | Sumário 3.1: `baixa`, `media`, `alta`, `urgente`; `corretiva_emergencial` nasce `urgente`; criticidade A → alta, B → media, C → baixa. A proposta da versão 0.1 foi abandonada. Aplicado em 6.6. |
| 13.12 | Nomes comerciais da marca (Tenaz Frota, Planta, Mobile) não apareciam no sumário. | Sumário 2: *Tenaz Planta* (ATV, MED, SOL, OS, PRE, CHK, EST, TER, CON), *Tenaz Frota* (FRO), *Tenaz Locação* (LOC), *Tenaz Mobile* (MOB). O menu continua com os nomes funcionais. |
| 13.14 | "Preventiva" era ao mesmo tempo cor de status do ativo e tipo de OS. | Resolvida com 13.1: rótulo "Em manutenção" com a cor de Preventiva. |

**Continuam abertas** (questões de marca e de acessibilidade, fora do escopo do sumário)

| # | Inconsistência | Onde | Decisão provisória |
|---|---|---|---|
| 13.3 | O chip de status do brandbook (texto na cor do status sobre fundo tingido) **reprova AA no tema claro**: 2,99 / 4,00 / 2,58 / 3,95 / 3,05. | Brandbook, `.chip-st` | Texto Aço no claro; padrão do brandbook mantido no escuro, onde passa. |
| 13.4 | O brandbook usa **Sinal como anel de foco** (`.copy:focus-visible`). Sinal tem 2,87:1 sobre branco e 2,60:1 sobre Névoa, abaixo dos 3:1 exigidos para indicadores de foco. Também conflita com a regra "Sinal só para ação principal". | Brandbook, CSS | Foco em Aço (claro) e `#E8EDF2` (escuro). |
| 13.5 | O brandbook descreve Concreto para "rótulos"; o README não. Concreto como texto tem 2,48:1 em branco e 2,25:1 em Névoa. | Brandbook × README | Concreto só decorativo e desabilitado. Rótulos em Aço 500 ou Aço 400. |
| 13.6 | O brandbook usa `#6B7887` como texto terciário, fora da escala oficial do Aço, e esse tom **reprova sobre Névoa** (4,09:1). | Brandbook, `--ink-3` | Proposta de **Aço 400 `#5F6C7B`** (4,86:1 em Névoa, 5,36:1 em branco). Precisa de aprovação para entrar na escala da marca. |
| 13.7 | O brandbook usa laranja como texto (`--accent-text` `#B54500`, unidades "%" do painel em Sinal), o que contraria "Sinal só no ponto e na ação principal". | Brandbook | Produto não usa laranja em texto. O tom `#B54500` fica restrito a material de marca. |
| 13.8 | O âmbar de Atenção `#C98C00` tem **2,90:1 sobre branco**, abaixo de 3:1 para gráficos e ícones isolados. | README | Reforço com contorno Aço ou fundo tingido (regra 4.3.4). Avaliar um âmbar mais escuro só para gráficos. |
| 13.10 | Percentual escrito de dois jeitos no brandbook: "97,4 %" (faixa do topo) e "97,4%" (cartão). | Brandbook | "97,4 %" com espaço inseparável, conforme o SI. |
| 13.13 | A marca lista 7 segmentos (inclui construção e mineração, agronegócio, hospitais e facilities, locadoras); o sumário agrupa em 4. | Brandbook × sumário 1.3 | Sem impacto na interface; alinhar no material comercial. |
| 13.15 | O brandbook ainda chama de "Preventiva" a cor usada para Em manutenção e abrevia o número de OS. | Brandbook | Decisão de marca: renomear a cor para "Em manutenção" (ou manter o nome da cor e só o rótulo na interface) e corrigir os exemplos de número. |
| 13.16 | O selo de versão (6.19) usa 12 dp no app, abaixo do mínimo de 14 dp da escala mobile (3.5 e 9.3). | Requisito de versão visível × regra do app | Exceção aceita só para o selo: é metadado, não conteúdo de trabalho, acompanha a escala de fonte do sistema, tem contraste AA e área de toque de 48 dp. Validar com os testes de campo. |
| 13.17 | A assinatura do representante do locatário na vistoria (8.16) é simples quando ele não tem usuário; a do representante da empresa é avançada. Falta confirmar se a assinatura simples basta como prova da avaria para a cobrança. | Sumário 3.1 × 9_SEGURANCA_CONFIABILIDADE.md, seção 8 | Manter as duas, com fotos, leituras e manifesto no PDF. Validar com jurídico junto com a P1. |
| 13.18 | A assinatura do cliente na OS de serviço (9.8) é simples, na tela do técnico, e é a prova de execução que a subcontratada entrega no relatório PDF (8.26). Mesma dúvida da 13.17, agora como base de faturamento entre empresas. | SRV × 9_SEGURANCA_CONFIABILIDADE.md | Resumo antes da assinatura, nome e função, hash e código de verificação no PDF. Validar com jurídico junto com a P1. |
| 13.19 | O relatório PDF (8.26) leva o logo de outra empresa (a subcontratada) no cabeçalho e a marca Tenaz só no rodapé. A marca não tem regra de convivência com logo de cliente. | Marca × 8.26 | Logo da conta no cabeçalho, símbolo Tenaz de 4 mm com "Gerado pela Tenaz" no rodapé. Precisa de aprovação de marca. |
| 13.20 | Não há ícone de marca para "conexão entre contas"; `cable` do Lucide foi escolhido para não confundir com `link` (link de uso único) nem com `plug` (Integrações). | 12.2 | Usar `cable` até teste com usuários; avaliar ícone próprio se houver confusão. |
| 13.21 | A leitura do número de série por OCR de texto impresso (9.8) não foi validada; etiquetas de peças variam muito. | 9.8 | Código de barras e QR no lançamento; OCR como proposta, a testar em campo com 5 técnicos. |
| 13.22 | Login único com várias contas no app (9.9) depende do modelo de identidade do ADM, ainda aberto no documento 5 (seção 10.8). | 9.9 × 5_ARQUITETURA_SOFTWARE.md | Por enquanto, uma entrada por conta ("Entrar em outra conta"). |
| 13.23 | A página pública de verificação do relatório (código e QR, 8.26) não está prevista no documento 3. | 8.26 × 3_ESPECIFICACAO_APIS.md | Mantida como proposta; sem ela, o código de verificação sai do PDF. |

**Outras pendências do documento**

- Validar com 5 técnicos e 5 motoristas o modo guiado do checklist e o tema Sol, em campo e com luva.
- Gerar e testar as instâncias estáticas de Archivo Expanded para o app (largura 125, pesos 700, 800 e 900).
- Definir a biblioteca de gráficos (proposta: Recharts no web; Victory Native ou componentes próprios sobre `react-native-svg` no app) no documento 5.
- Montar a biblioteca de componentes com documentação viva (Storybook) e testes de contraste e de acessibilidade automatizados.

---

## Histórico

| Versão | Data | Alteração |
|---|---|---|
| 0.1 | 12/11/2026 | Primeira versão: princípios, tokens, contraste, componentes, telas, app de campo, microcopy, acessibilidade, iconografia |
| 0.3 | 12/11/2026 | Alinhada ao sumário executivo 0.4: permissão `os.encerrar`, ação Assumir, link `/l/#<token>`, regras de solicitação e de contrato de locação, rótulos de origem da OS, papéis Posto de combustível e Seguradora, peça em trânsito e `est.decidir_destino`, locatário só no web, RNF-039 e RNF-040; módulo SRV: menu, estados de conexão e de sincronização, painel da subcontratada, técnicos e especialidades, clientes e equipamentos de clientes, abertura de OS de serviço, quadro de distribuição, app do técnico com componentes, número de série, fotos e assinatura do cliente, relatório PDF, conexão com cliente, OS vinculada nas duas contas, troca de conta e push; pendências 13.18 a 13.23 |
| 0.2 | 12/11/2026 | Alinhada ao sumário executivo 0.3: prioridade `urgente`, status do ativo com as regras do sumário, transições e reconhecimento da OS, executor e responsável pelo custo, numeração OS/SS/OR/CL, estados de ordem de reparo e contrato de locação; telas de fornecedor, estoque geral e por subcontratada, ordem de reparo, contrato de locação, vistoria, portal do prestador, link de uso único e portal do locatário; selo de versão e tela de login; inconsistências 13.1, 13.2, 13.9, 13.11, 13.12 e 13.14 resolvidas |
