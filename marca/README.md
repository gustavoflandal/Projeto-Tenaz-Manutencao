# Marca Tenaz

Identidade visual aprovada em 12/11/2026. A referência completa, com construção do símbolo, usos incorretos e aplicações, está em [brandbook.html](brandbook.html). Abra o arquivo no navegador.

## Nome

**Tenaz**: firme, persistente, que não cede. Também é o nome de uma ferramenta de oficina e de forja.

- No texto corrido: **Tenaz**, com inicial maiúscula. Caixa alta só no logotipo.
- Assinatura completa: **Tenaz · Gestão de Manutenção e Frota**.
- Slogan: **Cada ativo pronto para operar.**
- Módulos: Tenaz Frota, Tenaz Planta, Tenaz Mobile.

## Arquivos do logo

| Pasta | Conteúdo |
|---|---|
| `logo/svg/` | Vetores: `simbolo`, `horizontal`, `vertical` e `wordmark` nas variantes `positivo`, `negativo`, `mono`, `mono-branco` e `sobre-sinal`. Também `favicon.svg`, que se adapta ao tema claro ou escuro, e `tenaz-app-icon.svg`. |
| `logo/png/` | Bitmaps com fundo transparente, ícones de app de 1024, 512, 192 e 180 px e favicons de 32 a 128 px. |

Use a versão SVG sempre que possível. Os arquivos foram gerados por script a partir das mesmas formas geométricas do brandbook.

## Símbolo

- Porca sextavada: a peça comum a caminhão, prensa, empilhadeira e gerador.
- T vazado: a inicial do nome e o desenho de uma chave T vista de frente.
- Ponto laranja: a luz de status do painel. É o único elemento de cor do símbolo.

Área de proteção: `2x` em todos os lados, onde `x` é a espessura da barra do T. Tamanho mínimo: símbolo a partir de 16 px; logo horizontal a partir de 112 px na tela e 25 mm no impresso.

## Cores

| Token | HEX | RGB | Uso |
|---|---|---|---|
| Aço | `#14202E` | 20 32 46 | Logo, títulos, menus, faixas |
| Sinal | `#FF6A13` | 255 106 19 | Ponto do símbolo, botão de ação principal |
| Névoa | `#F2F4F6` | 242 244 246 | Fundo de telas e documentos |
| Concreto | `#9AA6B2` | 154 166 178 | Linhas e elementos decorativos (nunca texto) |

Escala do Aço para interface: 900 `#14202E`, 700 `#2B3A4C`, 500 `#4A5868`, 400 `#5F6C7B`, 300 `#9AA6B2`, 100 `#D9DFE6`, 50 `#F2F4F6`. O Aço 400 é o menor tom aceito para texto (texto terciário).

### Regras de acessibilidade (WCAG 2.1 AA)

Contrastes calculados em [docs/7_PADROES_DESIGN.md](../docs/7_PADROES_DESIGN.md).

- Texto e ícone sobre Sinal são sempre Aço (5,74:1). Branco sobre Sinal não passa (2,87:1) e é proibido.
- Sinal nunca é cor de texto, de contorno de foco, de status nem de série de gráfico.
- O contorno de foco é Aço no tema claro.
- Concreto não é usado em texto (2,48:1 sobre branco).
- Chip de status no tema claro: texto Aço e ícone na cor do status. No tema escuro, texto e ícone na cor do status.
- O âmbar de Atenção usado sozinho em gráfico ou ícone recebe contorno Aço.

Proporção de uso: Névoa 60%, Aço 28%, Concreto 8%, Sinal 4%.

Os valores CMYK e Pantone que aparecem no brandbook são aproximados. Confirme em prova de cor antes de imprimir.

### Cores de status do ativo

São cores semânticas, separadas das cores da marca. Sempre aparecem com ícone e texto.

| Status | Claro | Escuro | Significado |
|---|---|---|---|
| Operando | `#1E9E6A` | `#3CC28B` | Disponível e dentro dos planos |
| Preventiva | `#2F6FDB` | `#6E9BF0` | Parada programada em execução |
| Atenção | `#C98C00` | `#F0B83A` | Plano vencendo ou leitura fora da faixa |
| Parado | `#D33A2C` | `#F06A5C` | Corretiva aberta, ativo indisponível |
| Inativo | `#7E8A97` | `#97A3B0` | Baixado, vendido ou em reserva |

O laranja Sinal nunca é usado como status, para não ser confundido com Atenção.

## Tipografia

| Papel | Família | Pesos |
|---|---|---|
| Títulos | Archivo Expanded (largura 125) | 700 a 900 |
| Texto e interface | Inter | 400 a 700 |
| Dados e códigos | JetBrains Mono | 400 e 600 |

As três famílias estão no Google Fonts sob licença OFL.

## Tom de voz

Direto, técnico sem jargão, calmo nos alertas. Diz o que fazer e quando.

- Faça: "Preventiva do FRT-0112 vence em 2.300 km."
- Evite: "ATENÇÃO!!! Seu veículo precisa urgentemente de manutenção!"
