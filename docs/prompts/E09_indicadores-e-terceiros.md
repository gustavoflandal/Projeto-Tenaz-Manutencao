# E09 — Indicadores e terceiros

**Data do pedido:** 2026-11-12 · **Sprint do roadmap:** S9 · **Marco:** —
**Pedido por:** dono do produto

## Objetivo

Os dados acumulados nas etapas anteriores passam a responder perguntas: disponibilidade, MTBF, MTTR, MTTA, custo por ativo. E o prestador externo ganha acesso ao sistema, restrito ao que é dele — o ponto mais sensível do MVP depois da conexão entre contas.

## Entrega

| Requisito | Observação de escopo nesta etapa |
|---|---|
| RF-IND-001 a RF-IND-005 | Indicadores essenciais |
| RF-IND-007 | completo |
| RF-ATV-010 | Completo |
| RF-TER-003 | OS executada por terceiro |
| RF-TER-004 | Perfil de prestador com acesso restrito |
| RF-TER-005 | Link de execução de uso único |
| RF-MOB-010 | completo |
| RF-SRV-012 | Assinatura do cliente na OS |
| RF-SRV-013 | Relatório da OS em PDF |
| RF-SRV-003 | Completo: documentos dos técnicos |
| RF-SRV-016 | OS do cliente chegando à conta da subcontratada |

## Contexto obrigatório

| Documento | Seção | Por quê |
|---|---|---|
| [4_FLUXOS_PROCESSO.md](../4_FLUXOS_PROCESSO.md) | 30 | Cálculo dos indicadores, fórmula por fórmula |
| [4_FLUXOS_PROCESSO.md](../4_FLUXOS_PROCESSO.md) | Fluxos 18 e 26 | OS executada por prestador; OS do cliente executada pela subcontratada conectada |
| [5_ARQUITETURA_SOFTWARE.md](../5_ARQUITETURA_SOFTWARE.md) | 6, 9.3, 15 | Escopo externo no código, suíte de isolamento, ADR-008 |
| [0_SUMARIO_EXECUTIVO.md](../0_SUMARIO_EXECUTIVO.md) | 3.6 | Link de uso único com token no fragmento da URL |
| [3_ESPECIFICACAO_APIS.md](../3_ESPECIFICACAO_APIS.md) | 11, 19 | Consultar KPIs; prestadores e link de uso único |

## Agentes previstos

Conduzem `tenaz-dominio-manutencao` (fórmulas) e `tenaz-isolamento` (prestador e link). Apoiam `tenaz-backend-go`, `tenaz-api-contrato`, `tenaz-conexao-contas`, `tenaz-design-sistema`, `tenaz-frontend-web`, `tenaz-mobile-offline`, `tenaz-qa-automacao`.

## Portão de saída

| Verificação | Como se comprova |
|---|---|
| **Teste de isolamento do prestador verde no CI (RNF-037)** | Saída da suíte com prestador de dois fornecedores na mesma empresa e o mesmo CNPJ em duas empresas. **Sem este teste verde, o perfil de prestador não é liberado** (risco R4) |
| Link de uso único | Token no fragmento da URL; reuso, uso em outro aparelho, link revogado, expirado e consumido recebendo a mesma resposta de não encontrado |
| Indicadores corretos | Cada fórmula com teste de valor esperado; KPIs de tempo usando início e conclusão, não o encerramento administrativo |
| "Como calculamos" | Todo indicador com a explicação visível ao usuário (RF-IND-002) |
| OS entre contas | OS atribuída a fornecedor conectado aparecendo na conta da subcontratada, ligada por vínculo, sem expor identificador da outra conta |
| Relatório da OS | PDF com assinatura do cliente, gerado também para cliente sem conta na Tenaz |

## Evidência esperada

`docs/screenshots/E09_indicadores-e-terceiros/`: painel de indicadores, "como calculamos" aberto, OS atribuída a fornecedor, portal do prestador, cache filtrado no app do prestador, página do link de uso único, relatório da OS em PDF com assinatura do cliente, OS recebida de cliente conectado.

## Fora do escopo

Indicadores avançados e relatórios agendados (fase 2). Contratos de prestação com SLA e avaliação de desempenho, que são RF-TER-006 e RF-TER-007, da fase 2. Oficina externa e orçamento, que são RF-TER-008 e RF-TER-009, da fase 2. Estoque em poder de subcontratada (fase 2).

## Decisões pendentes do dono do produto

- **P7** — local da cópia de recuperação de desastre. Prazo do roadmap: fim desta etapa, porque afeta o RTO publicado no lançamento.

## Ajustes durante a etapa

Nenhum até agora.
