# E05 — Execução da OS e OS de serviço

**Data do pedido:** 2026-11-12 · **Sprint do roadmap:** S5 · **Marco:** —
**Pedido por:** dono do produto

## Objetivo

O técnico passa a executar a ordem de serviço do começo ao fim, inclusive sem sinal: assume, inicia, aponta horas, cumpre tarefas, registra o código de falha, pausa com motivo e conclui. A subcontratada passa a abrir OS de serviço com o serviço a realizar vindo do catálogo dela.

## Entrega

| Requisito | Observação de escopo nesta etapa |
|---|---|
| RF-OS-006 | completo |
| RF-OS-007 | completo, com a ação de assumir |
| RF-OS-008, RF-OS-009 | completos: tarefas e apontamento de horas |
| RF-OS-011 | Catálogo padrão de códigos de falha: modo, causa e atividade |
| RF-OS-015, RF-OS-018 | completos |
| RF-MOB-003 | Parcial: execução de OS offline |
| RF-MOB-005 | completo |
| RF-SRV-006 | Catálogo próprio de serviços da subcontratada |
| RF-SRV-007 | Abertura de OS de serviço com equipamento, técnico e serviço a realizar |

## Contexto obrigatório

| Documento | Seção | Por quê |
|---|---|---|
| [0_SUMARIO_EXECUTIVO.md](../0_SUMARIO_EXECUTIVO.md) | 3.3, 3.6 | Motivos de pausa, executor e responsável pelo custo, ação de assumir |
| [4_FLUXOS_PROCESSO.md](../4_FLUXOS_PROCESSO.md) | Fluxos 02 e 24 | Ciclo da OS e OS de serviço da subcontratada |
| [4_FLUXOS_PROCESSO.md](../4_FLUXOS_PROCESSO.md) | 29 | Regras de negócio dos módulos OS e SRV |
| [3_ESPECIFICACAO_APIS.md](../3_ESPECIFICACAO_APIS.md) | 8, 22 | Ações de estado da OS e endpoints do módulo SRV |
| [7_PADROES_DESIGN.md](../7_PADROES_DESIGN.md) | 9 | App de campo |

## Agentes previstos

Conduz `tenaz-dominio-manutencao` no comportamento. Apoiam `tenaz-backend-go`, `tenaz-api-contrato`, `tenaz-conexao-contas` (catálogo reaproveitado pelo módulo TER em E07), `tenaz-design-sistema`, `tenaz-frontend-web`, `tenaz-mobile-offline`, `tenaz-qa-automacao`.

## Portão de saída

| Verificação | Como se comprova |
|---|---|
| Assumir restrito | Técnico assume OS `aberta` ou `planejada` da mesma unidade dos seus perfis; OS de outra unidade recusada |
| Pausa com motivo | Pausa sem motivo recusada; cancelamento direto de OS `em_execucao` recusado |
| Apontamento por técnico | Horas somadas por técnico, com início, fim e pausas, sem sobreposição no mesmo técnico |
| Códigos de falha | Catálogo padrão de modo, causa e atividade aplicável por classe de ativo |
| Execução offline completa | OS iniciada, apontada, concluída em modo avião e sincronizada sem perda nem duplicidade |
| Impressão da OS | PDF gerado com os dados da execução |

## Evidência esperada

`docs/screenshots/E05_execucao-da-os/`: fila do técnico no app, execução passo a passo, apontamento, pausa com motivo, registro de código de falha, conclusão, OS impressa, abertura de OS de serviço com catálogo.

## Fora do escopo

Baixa de peça em OS, que é RF-OS-010 e entra em E07. Campos de NR-12 e assinatura (E08). Custo consolidado e liberação, que são RF-OS-013 e RF-OS-014, em E08. Distribuição e reatribuição de OS de serviço, que é RF-SRV-010, em E06.

## Decisões pendentes do dono do produto

Nenhuma.

## Ajustes durante a etapa

Nenhum até agora.
