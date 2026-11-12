# E10 — API, SSO, cobrança e portão da conexão

**Data do pedido:** 2026-11-12 · **Sprint do roadmap:** S10 · **Marco:** M4
**Pedido por:** dono do produto

## Objetivo

Fecha o escopo funcional do MVP: a API pública ganha chaves, escopos, limites, webhooks e documentação; entra o login por provedor de identidade; entra a cobrança por técnico e por veículo; e a conexão entre contas completa o ciclo e passa pelo portão que a libera para produção.

## Entrega

| Requisito | Observação de escopo nesta etapa |
|---|---|
| RF-INT-001 | API pública v1 com chaves, escopos, limites e documentação |
| RF-INT-002 | Webhooks |
| RF-INT-003 | Completo |
| RF-ADM-007 | SSO por OIDC |
| RF-ADM-009 | Cobrança, com prestador contando como técnico e com o plano Serviços |
| RF-ADM-013 | completo |
| RF-SRV-017 | OS aberta pela subcontratada nascendo na conta do cliente |
| RF-SRV-018 | Andamento de volta e liberação pelo cliente |
| RF-SRV-019 | Suspensão e encerramento da conexão |
| RF-SRV-020 | Relatório de serviços por cliente |
| RF-SRV-021 | Indicadores da conta de serviços |
| RF-SRV-001 | Completo no que se refere à cobrança do plano Serviços |
| RF-SRV-011 | Completo |

## Contexto obrigatório

| Documento | Seção | Por quê |
|---|---|---|
| [3_ESPECIFICACAO_APIS.md](../3_ESPECIFICACAO_APIS.md) | 2, 4, 5, 13, 18, 22 | Convenções, autenticação, limites, webhooks, versionamento, módulo SRV |
| [4_FLUXOS_PROCESSO.md](../4_FLUXOS_PROCESSO.md) | Fluxos 25, 26, 27 e 28 | Conexão, OS entre contas nos dois sentidos, falha de sincronização |
| [0_SUMARIO_EXECUTIVO.md](../0_SUMARIO_EXECUTIVO.md) | 2.1, 2.2, 3.6 | Base de cobrança, plano Serviços, devolução de OS pela subcontratada |
| [5_ARQUITETURA_SOFTWARE.md](../5_ARQUITETURA_SOFTWARE.md) | 14, 15 | Observabilidade do atraso da sincronização; ADR-009 |

## Agentes previstos

Conduzem `tenaz-api-contrato` (API e SSO) e `tenaz-conexao-contas` (ciclo completo entre contas). Apoiam `tenaz-isolamento`, `tenaz-backend-go`, `tenaz-plataforma-sre`, `tenaz-frontend-web`, `tenaz-mobile-offline`, `tenaz-qa-automacao`, `tenaz-produto` (contagem dos 126 requisitos).

## Portão de saída

| Verificação | Como se comprova |
|---|---|
| **Marco M4** | Os 126 requisitos funcionais do MVP entregues, conferidos um a um contra o [roadmap](../6_ROADMAP.md) |
| **Portão da conexão: RNF-041 verde** | Suíte de três contas verde no CI, com todos os casos da seção 9.3 da arquitetura |
| **Portão da conexão: RNF-042 medido** | Atraso p95 da sincronização abaixo de 60 segundos medido em homologação |
| Ciclo completo entre contas | OS do cliente chegando, executada offline, andamento de volta, liberação recusada e aceita; OS aberta pela subcontratada nascendo no cliente; encerramento da conexão sem perda de histórico |
| Contrato estável | Contrato gerado igual ao versionado; nenhuma mudança incompatível em `/v1` |
| Cobrança | Maior número diário de técnicos e de veículos no ciclo; prestador que executa contando; link de uso único não contando; técnico de subcontratada com conta própria não contando no cliente |
| Documentação pública da API | Publicada e gerada do mesmo contrato |
| Nenhum defeito crítico aberto | Lista de defeitos revisada |

Portão vermelho dispara a contingência do [roadmap](../6_ROADMAP.md), seção 6: a conexão sai do lançamento e vai para o mês 7, com o módulo SRV lançado no modo sem cliente na Tenaz. **Exige aval do dono do produto**, porque contraria a seção 4 do sumário.

## Evidência esperada

`docs/screenshots/E10_api-sso-cobranca/`: documentação pública da API, tela de chaves e escopos, entrega de webhook com tentativa e repetição, login por SSO, ciclo de OS entre duas contas lado a lado, suspensão e encerramento da conexão, relatório de serviços por cliente, tela de cobrança.

## Fora do escopo

SAML e SCIM, que são RF-ADM-008, da fase 2. Conectores de ERP e de telemática (fase 2). Servidor MCP (fase 3). Qualquer requisito das fases 2 e 3.

## Decisões pendentes do dono do produto

- **P4** — validação de preço, mínimos e desconto de canal, antes da página de preços de E12.
- Aval sobre a contingência, se o portão da conexão não fechar.

## Ajustes durante a etapa

Nenhum até agora.
