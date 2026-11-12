# E04 — Solicitações, núcleo da OS e fundação da conexão

**Data do pedido:** 2026-11-12 · **Sprint do roadmap:** S4 · **Marco:** —
**Pedido por:** dono do produto

## Objetivo

Qualquer pessoa passa a pedir manutenção, a triagem converte o pedido em ordem de serviço e a OS existe com todos os seus estados e marcos de tempo. Em paralelo, a fundação da conexão entre contas entra com a suíte de três contas já no CI, antes de existir qualquer fluxo real entre contas.

## Entrega

| Requisito | Observação de escopo nesta etapa |
|---|---|
| RF-SOL-001 a RF-SOL-005 | completos |
| RF-OS-001 a RF-OS-005 | completos: criação, campos, estados, reconhecimento, paradas |
| RF-OS-016, RF-OS-017 | completos |
| RF-ADM-012 | completo |
| RF-MOB-006 | completo |
| RF-MOB-003 | Parcial: solicitação offline |

Fundação da conexão, sem requisito funcional: ADR-009 fechado em [5_ARQUITETURA_SOFTWARE.md](../5_ARQUITETURA_SOFTWARE.md); entidade de conexão no esquema da plataforma, fora do RLS por empresa; papel de banco próprio do sincronizador, sem privilégio de contornar RLS; fila de saída nos dois lados; idempotência por identificador de mensagem e ordem por sequência; **suíte de três contas rodando no CI**, ainda sem fluxo de OS.

## Contexto obrigatório

| Documento | Seção | Por quê |
|---|---|---|
| [0_SUMARIO_EXECUTIVO.md](../0_SUMARIO_EXECUTIVO.md) | 3.3, 3.5, 3.6 | Estados canônicos da solicitação e da OS, marcos de tempo, conexão entre contas |
| [4_FLUXOS_PROCESSO.md](../4_FLUXOS_PROCESSO.md) | Fluxos 01 e 02 | Solicitação, triagem e ciclo de vida da OS |
| [5_ARQUITETURA_SOFTWARE.md](../5_ARQUITETURA_SOFTWARE.md) | 6, 15 | Isolamento no código e ADR-009 |
| [5_ARQUITETURA_SOFTWARE.md](../5_ARQUITETURA_SOFTWARE.md) | 9.3 | Teste de isolamento entre contas conectadas, com três contas |
| [3_ESPECIFICACAO_APIS.md](../3_ESPECIFICACAO_APIS.md) | 7, 8 | Solicitações e ordens de serviço |

## Agentes previstos

Conduzem `tenaz-arquiteto` (fecha ADR-009) e `tenaz-dominio-manutencao` (estados e marcos). Apoiam `tenaz-conexao-contas`, `tenaz-isolamento`, `tenaz-banco-dados`, `tenaz-backend-go`, `tenaz-api-contrato`, `tenaz-frontend-web`, `tenaz-mobile-offline`, `tenaz-qa-automacao`.

## Portão de saída

| Verificação | Como se comprova |
|---|---|
| Máquina de estados fiel | Toda transição do diagrama do sumário permitida e toda transição fora dele recusada, com erro de código estável; mesma bateria de casos no Go e no TypeScript |
| Reconhecimento | `data_reconhecimento` gravado uma única vez, também quando implícito ao programar, iniciar ou assumir |
| Desistência e duplicidade | Solicitação indo de `nova` direto para `recusada` com motivo de desistência; `duplicada` final, com reabertura como nova |
| **Suíte de três contas verde no CI** | Saída da suíte com a verificação de cobertura do papel do sincronizador e das colunas da lista fechada |
| Papel do sincronizador | Papel sem privilégio de contornar RLS e sem ser dono de tabela; consulta a tabela fora da lista recebendo permissão negada |
| App nas lojas | Primeira versão em teste fechado nas duas lojas |

## Evidência esperada

`docs/screenshots/E04_solicitacoes-e-os/`: portal de solicitação, fila de triagem, conversão em OS, ficha da OS com os estados, solicitação aberta no app em modo avião.

## Fora do escopo

Execução da OS, apontamento de horas, peças e códigos de falha (E05). Qualquer fluxo real de OS entre contas conectadas (E07 a E10). Assinatura (E08).

## Decisões pendentes do dono do produto

Nenhuma trava a etapa. O **termo da conexão entre contas** (LGPD) precisa estar aprovado antes de E07, quando o convite aparece.

## Ajustes durante a etapa

Nenhum até agora.
