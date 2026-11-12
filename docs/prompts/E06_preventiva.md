# E06 — Preventiva

**Data do pedido:** 2026-11-12 · **Sprint do roadmap:** S6 · **Marco:** M2
**Pedido por:** dono do produto

## Objetivo

O sistema passa a trabalhar sozinho: planos por tempo, por medidor ou pelo que vencer primeiro geram ordens de serviço na hora certa. É o primeiro valor que o cliente percebe sem fazer nada. Fecha o marco de alpha interno, com o app provado em 72 horas offline.

## Entrega

| Requisito | Observação de escopo nesta etapa |
|---|---|
| RF-PRE-001 a RF-PRE-006 | completos |
| RF-PRE-008, RF-PRE-009, RF-PRE-010 | completos |
| RF-ADM-010 | completo |
| RF-SRV-009 | Início, fim e apontamentos por técnico na OS de serviço |
| RF-SRV-010 | Distribuição e reatribuição de OS de serviço |
| RF-TER-001 | Antecipado do sprint seguinte: fornecedor com papéis, porque o convite de conexão parte dele |
| RF-SRV-011 | Parcial: tarefas do técnico da subcontratada no app. Completo em E10 |

## Contexto obrigatório

| Documento | Seção | Por quê |
|---|---|---|
| [4_FLUXOS_PROCESSO.md](../4_FLUXOS_PROCESSO.md) | Fluxo 03 | Manutenção preventiva: geração, antecipação, atraso |
| [4_FLUXOS_PROCESSO.md](../4_FLUXOS_PROCESSO.md) | 31 | Parâmetros configuráveis que afetam a geração |
| [0_SUMARIO_EXECUTIVO.md](../0_SUMARIO_EXECUTIVO.md) | 3.1 | Plano, gatilho e procedimento |
| [2_ARQUITETURA_BANCO_DADOS.md](../2_ARQUITETURA_BANCO_DADOS.md) | 5 | DDL do MVP para planos e gatilhos |
| [5_ARQUITETURA_SOFTWARE.md](../5_ARQUITETURA_SOFTWARE.md) | 4.3, 15 | Fila de jobs e ADR-005 |

## Agentes previstos

Conduz `tenaz-dominio-manutencao`. Apoiam `tenaz-backend-go`, `tenaz-banco-dados`, `tenaz-api-contrato`, `tenaz-frontend-web`, `tenaz-mobile-offline`, `tenaz-qa-automacao`, `tenaz-conexao-contas` (fornecedor e distribuição).

## Portão de saída

| Verificação | Como se comprova |
|---|---|
| **Marco M2** | Ciclo completo de solicitação a encerramento, com reconhecimento, assumir, desprogramação e liberação recusada |
| Três tipos de gatilho | Plano por tempo, por medidor e pelo que vencer primeiro, cada um gerando OS `planejada` na data certa, com teste de borda |
| Geração idempotente | O job rodando duas vezes no mesmo dia não duplica OS |
| **72 horas offline** | App operando 72 horas sem conexão, com ações acumuladas e sincronizadas sem perda |
| OS de serviço distribuída | OS de serviço aberta, distribuída, reatribuída e executada offline por técnico de subcontratada, com horas por técnico |
| Biblioteca de planos padrão | Plano criado a partir da biblioteca, aplicável por classe de ativo |

## Evidência esperada

`docs/screenshots/E06_preventiva/`: cadastro de plano com os três gatilhos, calendário de preventivas, OS gerada por plano, biblioteca de planos padrão, distribuição de OS de serviço, app depois de 72 horas offline.

## Fora do escopo

Planos aninhados, que são RF-PRE-007, da fase 2. Manutenção por condição (fase 3). Checklists (E07).

## Decisões pendentes do dono do produto

Nenhuma. Revisão de escopo prevista no roadmap ao fim deste sprint, com os dados de velocidade das etapas anteriores.

## Ajustes durante a etapa

Nenhum até agora.
