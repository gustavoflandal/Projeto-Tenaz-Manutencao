# 006 — Gestão de serviços do prestador e alinhamento 0.4

**Data:** 2026-11-12 · **Tipo:** escopo
**Pedido por:** dono do produto · **Versão do produto:** —

## Pedido
Módulo para as subcontratadas gerirem seus serviços: cadastro de técnicos, abertura de OS vinculada a equipamento e técnico, serviço a realizar, componentes substituídos, hora de início e fim, e app em smartphone ou tablet para o técnico receber tarefas e apontar o andamento. Escolhas do dono do produto: **conta própria da subcontratada** e entrada no **MVP**.

## Decisões
- Módulo **SRV** ("Tenaz Serviços"): a subcontratada é uma `empresa` com SRV; reaproveita OS, apontamento, ativos e o app; `cliente_servico`, ativo `de_cliente`, componentes por `tipo_material`, relatório PDF; plano Serviços por técnico (proposta).
- **Conexão entre contas:** `plataforma.conexao_empresa`, `vinculo_os` com uma linha em cada conta, sincronização pelo papel `tenaz_sinc` em transações separadas, lista fechada de 15 grupos de campos, idempotência por `mensagem_id` e `sequencia` (ADR-009, risco R5).
- Divergências de nomes entre documentos fechadas no sumário 3.6 (permissão `os.encerrar`, escopo `prestador`, link com token no fragmento, ação `assumir`, origens da OS, categorias da matriz e outras).
- Na conferência final: aplicação lê a conexão só por funções seguras; o cliente escolhe os grupos compartilhados; ação `devolver` da subcontratada; 8 diagramas de sequência corrigidos no documento 5.

## Mudanças
| Documento ou arquivo | O que mudou | IDs afetados |
|---|---|---|
| [0_SUMARIO_EXECUTIVO.md](../0_SUMARIO_EXECUTIVO.md) | Versão 0.4 (seções 2.2, 3.5, 3.6) | R5, R6 |
| Documentos 1 a 9 | Alinhados ao sumário 0.4 (versões 0.3) | RF-SRV-001 a 021, RNF-041, RNF-042, RN-SRV-*, ADR-009 |
| [docs/pdf/](../pdf/) | PDFs gerados de todos os documentos | — |

## Verificação
Busca por termos antigos; checagem de IDs RN (304 definidos, nenhum repetido); renderização individual de todos os diagramas Mermaid; PDFs sem erro de diagrama.

## Pendências
- Tamanho do time do MVP: 14 pessoas para manter 13 sprints, ou estender o prazo — [6_ROADMAP.md](../6_ROADMAP.md).
- A ação `devolver` da subcontratada está decidida no sumário 3.6, mas ainda não detalhada nos documentos 3 e 4.
- P4 (preço do plano Serviços e desconto de canal), P7, P8.
