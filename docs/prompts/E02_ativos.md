# E02 — Ativos

**Data do pedido:** 2026-11-12 · **Sprint do roadmap:** S2 · **Marco:** —
**Pedido por:** dono do produto

## Objetivo

A empresa cadastra o que tem: localizações em árvore, ativos em árvore, classes, criticidade, QR code e documentos técnicos. O app passa a trazer os ativos em cache. A conta de subcontratada ganha os clientes atendidos e as especialidades de técnico.

## Entrega

| Requisito | Observação de escopo nesta etapa |
|---|---|
| RF-ATV-001, RF-ATV-002, RF-ATV-003 | completos |
| RF-ATV-005, RF-ATV-006, RF-ATV-007, RF-ATV-008 | completos |
| RF-INT-003 | Parcial: importação por planilha de localizações e ativos. Completa em E10 |
| RF-MOB-002 | Cache de ativos |
| RF-MOB-004 | completo |
| RF-SRV-004 | Clientes atendidos |
| RF-SRV-002 | Especialidades de técnico |

## Contexto obrigatório

| Documento | Seção | Por quê |
|---|---|---|
| [0_SUMARIO_EXECUTIVO.md](../0_SUMARIO_EXECUTIVO.md) | 3.1, 3.3, 3.4 | Tipos e propriedade do ativo, status operacional, formato da TAG |
| [2_ARQUITETURA_BANCO_DADOS.md](../2_ARQUITETURA_BANCO_DADOS.md) | 5, 6, 10 | DDL, hierarquia de ativos e localizações, anexos |
| [4_FLUXOS_PROCESSO.md](../4_FLUXOS_PROCESSO.md) | 29 | Regras de negócio do módulo ATV |
| [7_PADROES_DESIGN.md](../7_PADROES_DESIGN.md) | 7, 8 | Cores de status e criticidade; telas de lista e de ficha |

## Agentes previstos

Conduz `tenaz-dominio-manutencao` no comportamento e `tenaz-backend-go` na implementação. Apoiam `tenaz-banco-dados`, `tenaz-isolamento`, `tenaz-api-contrato`, `tenaz-design-sistema`, `tenaz-frontend-web`, `tenaz-mobile-offline`, `tenaz-qa-automacao`.

## Portão de saída

| Verificação | Como se comprova |
|---|---|
| Árvores com limite e sem ciclo | Tentativa de criar o sétimo nível e tentativa de tornar um ativo filho de si mesmo, as duas recusadas |
| Status operacional correto | Teste cobrindo cada causa de `em_manutencao`, `parado` e `inativo`, com `parado` prevalecendo |
| TAG única por empresa | Duas empresas com a mesma TAG aceitas; a mesma TAG repetida na mesma empresa recusada |
| Importação | 2.000 ativos importados por planilha, com relatório de erro linha a linha |
| QR code | Leitura no app abrindo a ficha do ativo correspondente |
| Isolamento | Suíte de vazamento cobrindo as tabelas novas, inclusive anexos com URL assinada |

## Evidência esperada

`docs/screenshots/E02_ativos/`: lista de ativos com filtros, ficha do ativo, árvore de localizações, cadastro de classe, importação por planilha com erros, leitura de QR no app, cadastro de cliente atendido.

## Fora do escopo

Medidores e leituras (E03). Composição de veículos, que é RF-ATV-004 e entra em E03. Dados específicos de veículo, que é RF-FRO-001 e entra em E07. Pneus e combustível (fase 2).

## Decisões pendentes do dono do produto

Nenhuma.

## Ajustes durante a etapa

Nenhum até agora.
