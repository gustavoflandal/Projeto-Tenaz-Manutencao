# E08 — NR-12, documentos e assinatura

**Data do pedido:** 2026-11-12 · **Sprint do roadmap:** S8 · **Marco:** M3
**Pedido por:** dono do produto

## Objetivo

Entra o diferencial do produto: conformidade brasileira nativa. OS de máquina NR-12 não conclui sem os campos que a norma exige, documentos com validade avisam antes de vencer, e a assinatura eletrônica avançada dá valor de evidência ao que foi registrado. Fecha o marco de piloto fechado, com empresas usando dados reais.

## Entrega

| Requisito | Observação de escopo nesta etapa |
|---|---|
| RF-CON-001 a RF-CON-004 | NR-12: registro do item 12.11.2, inventário, bloqueio e etiquetagem do item 12.11.3 |
| RF-CON-005, RF-CON-006 | Documentos com vencimento e alertas |
| RF-CON-014 | Assinatura eletrônica avançada, com arquitetura plugável |
| RF-CON-016 | Dossiê de fiscalização |
| RF-OS-013, RF-OS-014 | Custo consolidado e conclusão com liberação, inclusive liberação recusada |
| RF-SRV-008 | Componentes substituídos, consumidos e fornecidos pelo cliente |
| RF-SRV-015 | Ativos cobertos pela conexão |

## Contexto obrigatório

| Documento | Seção | Por quê |
|---|---|---|
| [4_FLUXOS_PROCESSO.md](../4_FLUXOS_PROCESSO.md) | Fluxos 06 e 08 | Encerramento de OS em máquina NR-12; documentos com vencimento |
| [9_SEGURANCA_CONFIABILIDADE.md](../9_SEGURANCA_CONFIABILIDADE.md) | 8 | Assinatura eletrônica de OS e checklists |
| [1_ESPECIFICACAO_REQUISITOS.md](../1_ESPECIFICACAO_REQUISITOS.md) | 4 (CON) | Redação dos requisitos de conformidade |
| [0_SUMARIO_EXECUTIVO.md](../0_SUMARIO_EXECUTIVO.md) | 3.1, 3.6 | Documento com vencimento, capacitação, bloqueio de ativo |
| [docs/pesquisa/](../pesquisa/) | normas e regulação | Base da interpretação das normas |

## Agentes previstos

Conduz `tenaz-conformidade-br` na interpretação da norma. Apoiam `tenaz-dominio-manutencao`, `tenaz-backend-go`, `tenaz-api-contrato`, `tenaz-conexao-contas`, `tenaz-design-sistema`, `tenaz-frontend-web`, `tenaz-mobile-offline`, `tenaz-qa-automacao`.

## Portão de saída

| Verificação | Como se comprova |
|---|---|
| **Marco M3** | Pelo menos 3 empresas piloto com dados reais, e pelo menos uma subcontratada piloto com conta própria |
| NR-12 bloqueante | Conclusão de OS de máquina NR-12 recusada sem os sete campos do item 12.11.2 |
| Alertas de documento | Alertas disparando em D-60, D-30 e D-7, com teste de calendário |
| Assinatura avançada | Reautenticação, hash do conteúdo, manifesto e PDF de evidências; troca do método sem migrar dado |
| Liberação recusada | OS `concluida` voltando a `em_execucao` com motivo registrado |
| Ativos cobertos | Lista de ativos escolhida pelo cliente sincronizando para a conta da subcontratada, só com os campos permitidos |
| Isolamento | Suíte de três contas provando que nenhum campo interno do cliente trafega junto com o ativo |
| Revisão de SST | Requisitos de conformidade revisados pelo consultor de SST |

## Evidência esperada

`docs/screenshots/E08_nr12-documentos-assinatura/`: conclusão de OS de NR-12 com os sete campos, recusa por campo faltando, assinatura com manifesto, agenda de documentos, alerta de vencimento, dossiê de fiscalização, componentes substituídos no app, ativos cobertos pela conexão.

## Fora do escopo

NR-13 e LOTO por ponto de bloqueio, que são da fase 2. Capacitações que bloqueiam atribuição, que é RF-CON-013, da fase 2 — no módulo SRV, o cadastro e os alertas de capacitação já entraram em E03, sem bloqueio. Assinatura ICP-Brasil, que é RF-CON-015, da fase 3.

## Decisões pendentes do dono do produto

- **P1** — parecer jurídico sobre o item 1.6 da NR-1 e a exigência de ICP-Brasil em documentos de SST. Prazo do roadmap: antes desta etapa. Sem o parecer, vale a contingência: o MVP sai com assinatura avançada e aviso ao cliente.
- **P2** — tabela oficial de validade de CIV e CIPP. Sem ela, a validade continua digitada a partir do certificado.

## Ajustes durante a etapa

Nenhum até agora.
