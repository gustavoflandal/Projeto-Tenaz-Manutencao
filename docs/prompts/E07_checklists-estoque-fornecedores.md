# E07 — Checklists, estoque básico e fornecedores

**Data do pedido:** 2026-11-12 · **Sprint do roadmap:** S7 · **Marco:** —
**Pedido por:** dono do produto

## Objetivo

Três frentes que se apoiam: o checklist no app, que vira solicitação quando um item reprova; o estoque básico, que permite reservar e baixar peça na OS; e o cadastro de fornecedores com catálogo de serviços, que abre o caminho do convite de conexão entre contas.

## Entrega

| Requisito | Observação de escopo nesta etapa |
|---|---|
| RF-CHK-001 a RF-CHK-006 | completos |
| RF-SOL-006 | completo |
| RF-EST-001 a RF-EST-004 | Estoque básico: itens, um almoxarifado por unidade, saldo, estoque mínimo, ajuste com motivo |
| RF-EST-006 | Reserva de estoque |
| RF-OS-010 | Baixa de peça em OS |
| RF-FRO-001 | Dados específicos de veículo |
| RF-TER-002 | Catálogo de serviços do fornecedor |
| RF-MOB-007 | completo |
| RF-SRV-014 | Convite de conexão, aceite e campos compartilhados |

## Contexto obrigatório

| Documento | Seção | Por quê |
|---|---|---|
| [4_FLUXOS_PROCESSO.md](../4_FLUXOS_PROCESSO.md) | Fluxos 05, 10 e 25 | Checklist no app; estoque; conexão entre contas |
| [0_SUMARIO_EXECUTIVO.md](../0_SUMARIO_EXECUTIVO.md) | 2.1, 3.1, 3.6 | Fronteira entre estoque básico e completo; almoxarifado e movimentação; grupos compartilhados na conexão |
| [2_ARQUITETURA_BANCO_DADOS.md](../2_ARQUITETURA_BANCO_DADOS.md) | 5 | DDL de checklist e de estoque |
| [3_ESPECIFICACAO_APIS.md](../3_ESPECIFICACAO_APIS.md) | 10, 22 | Executar checklist; conexão entre contas |

## Agentes previstos

Conduz `tenaz-dominio-manutencao` no estoque e no checklist, e `tenaz-conexao-contas` no convite. Apoiam `tenaz-banco-dados`, `tenaz-backend-go`, `tenaz-api-contrato`, `tenaz-isolamento`, `tenaz-frontend-web`, `tenaz-mobile-offline`, `tenaz-qa-automacao`.

## Portão de saída

| Verificação | Como se comprova |
|---|---|
| Item reprovado gera solicitação | Checklist com item reprovado criando solicitação ligada ao ativo e ao checklist |
| Checklist offline | Checklist preenchido em modo avião, com foto, sincronizado depois |
| Reserva e baixa | Peça reservada na OS indisponível para outra OS; baixa reduzindo o saldo; devolução repondo |
| Saldo nunca negativo | Baixa acima do saldo recusada; ajuste exigindo motivo |
| Convite de conexão | Convite enviado, lista de grupos compartilhados visível antes do aceite, aceite e recusa funcionando, estados da conexão respeitados |
| **Convite restrito a homologação** | Enquanto o termo de LGPD da conexão não estiver aprovado, o convite não funciona em produção |
| Isolamento | Suíte de três contas cobrindo o convite e os campos compartilhados; suíte de vazamento cobrindo as tabelas de estoque |

## Evidência esperada

`docs/screenshots/E07_checklists-estoque-fornecedores/`: modelo de checklist, execução no app, item reprovado gerando solicitação, lista de itens de estoque, saldo, reserva e baixa em OS, cadastro de fornecedor com papéis, catálogo de serviços, convite de conexão com os grupos compartilhados.

## Fora do escopo

Vários almoxarifados, transferência, inventário, requisição, compras e NF-e, que são da fase 2. Estoque em poder de subcontratada e peças serializadas (fase 2). Pneus e combustível (fase 2). Ativos cobertos pela conexão, que é RF-SRV-015, em E08.

## Decisões pendentes do dono do produto

- **Termo da conexão entre contas (LGPD)** aprovado pelo advogado. Sem ele, o convite fica restrito a homologação. Dependência crítica 18 do [roadmap](../6_ROADMAP.md).

## Ajustes durante a etapa

Nenhum até agora.
