# E12 — Lançamento

**Data do pedido:** 2026-11-12 · **Sprint do roadmap:** S12 · **Marco:** M5
**Pedido por:** dono do produto

## Objetivo

Publicar a versão `v1.0.0`: primeiros passos guiados para quem entra sozinho, página de status no ar, preços publicados, app nas lojas e os pilotos migrados para produção. Nada aqui é técnico difícil; o que trava são as pendências que dependem de terceiros.

## Entrega

| Requisito | Observação de escopo nesta etapa |
|---|---|
| RF-ADM-001 | Completo: roteiro de primeiros passos |
| RF-SRV-001 | Completo: roteiro de primeiros passos da conta de serviços |
| RNF-020 | Página de status com histórico |
| RNF-039, RNF-040 | Release `v1.0.0` conforme a seção 3.3 do [roadmap](../6_ROADMAP.md) |

Também entra: documentação pública da API publicada, kit de LGPD publicado com o termo da conexão entre contas, RPO e RTO publicados, página de preços, migração dos pilotos para produção, publicação nas lojas.

## Contexto obrigatório

| Documento | Seção | Por quê |
|---|---|---|
| [6_ROADMAP.md](../6_ROADMAP.md) | 3.2, 3.3, 7 | Critério do marco M5, critério de release, atividades não técnicas |
| [9_SEGURANCA_CONFIABILIDADE.md](../9_SEGURANCA_CONFIABILIDADE.md) | 1, 7, 10 | Metas publicadas, LGPD, certificações |
| [0_SUMARIO_EXECUTIVO.md](../0_SUMARIO_EXECUTIVO.md) | 2.1, 2.2, 5, 6 | Planos e preços, metas de qualidade, pendências |
| [7_PADROES_DESIGN.md](../7_PADROES_DESIGN.md) | 8, 10 | Telas de primeiros passos e tom de voz |

## Agentes previstos

Conduz `tenaz-plataforma-sre` no release. Apoiam `tenaz-produto` (conferência do escopo), `tenaz-design-sistema` (primeiros passos e tom de voz), `tenaz-frontend-web`, `tenaz-mobile-offline`, `tenaz-documentacao` (CHANGELOG da versão e registro da etapa), `tenaz-conformidade-br` (kit de LGPD).

## Portão de saída

| Verificação | Como se comprova |
|---|---|
| **Marco M5** | Todos os itens da seção 3.2 do [roadmap](../6_ROADMAP.md) para M5 conferidos um a um |
| Release íntegro | Tag `v1.0.0`; seção correspondente no `CHANGELOG.md`; web, app, API e workers com a mesma versão, gerada pela tag; `GET /v1/versao` respondendo a versão publicada depois do deploy |
| Primeiros passos | Empresa criada do zero, sem treinamento, chegando a abrir a primeira OS |
| Página de status | No ar, com histórico de incidentes |
| Metas publicadas | RPO e RTO publicados; RTO de região só depois de **P7** |
| Preços publicados | Depois de **P4** |
| Marca liberada | Depois de **P5** |
| App publicado | Aprovado nas duas lojas |

## Evidência esperada

`docs/screenshots/E12_lancamento/`: primeiros passos guiados passo a passo, página de status, página de preços, tela principal com o selo `v1.0.0` em produção, ficha do app nas lojas, "Novidades desta versão".

## Fora do escopo

Qualquer requisito das fases 2 e 3. SOC 2, que começa no mês 7.

## Decisões pendentes do dono do produto

Esta etapa só fecha com **P4**, **P5** e **P7** resolvidas. São as três pendências que o roadmap coloca como condição do lançamento. Nenhuma delas se resolve dentro do repositório.

## Ajustes durante a etapa

Nenhum até agora.
