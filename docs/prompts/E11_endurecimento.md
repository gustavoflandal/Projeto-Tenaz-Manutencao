# E11 — Endurecimento

**Data do pedido:** 2026-11-12 · **Sprint do roadmap:** S11 · **Marco:** —
**Pedido por:** dono do produto

## Objetivo

Nenhum requisito novo. Esta etapa existe para provar, com número, o que as etapas anteriores afirmaram: desempenho, 7 dias offline, restauração de backup, acessibilidade e segurança. É também quando as correções pedidas pelos pilotos entram.

## Entrega

| Requisito | Observação de escopo nesta etapa |
|---|---|
| RNF-007, RNF-008, RNF-009 | Verificados em carga, não implementados |
| RNF-030 | Verificado em carga |
| RNF-010 | 7 dias de operação offline sem perda de registro |
| RNF-005 | Restauração de backup ensaiada, com RPO e RTO medidos |
| RNF-025 | WCAG 2.1 AA verificado automática e manualmente em todas as telas |
| RNF-027 | App testado em aparelhos de entrada |
| RNF-042 | Atraso da sincronização entre contas medido **sob carga** |

Também entra: pentest focado em isolamento entre empresas, acesso do prestador, link de uso único e conexão entre contas; primeiro release candidato em homologação; reconciliação diária da sincronização; correções dos pilotos.

## Contexto obrigatório

| Documento | Seção | Por quê |
|---|---|---|
| [1_ESPECIFICACAO_REQUISITOS.md](../1_ESPECIFICACAO_REQUISITOS.md) | 5 | Metas dos requisitos não funcionais, com o número de cada uma |
| [9_SEGURANCA_CONFIABILIDADE.md](../9_SEGURANCA_CONFIABILIDADE.md) | 1, 3, 11 | Metas de confiabilidade, backup e recuperação, checklist de segurança por release |
| [5_ARQUITETURA_SOFTWARE.md](../5_ARQUITETURA_SOFTWARE.md) | 9, 14 | Pirâmide de testes, carga, observabilidade |
| [7_PADROES_DESIGN.md](../7_PADROES_DESIGN.md) | 4, 11 | Contraste WCAG e acessibilidade |

## Agentes previstos

Conduz `tenaz-qa-automacao`. Apoiam `tenaz-plataforma-sre` (carga, backup, observabilidade), `tenaz-isolamento` e `tenaz-conexao-contas` (pentest e carga da sincronização), `tenaz-design-sistema` (acessibilidade), `tenaz-mobile-offline` (7 dias offline e aparelhos de entrada).

## Portão de saída

| Verificação | Como se comprova |
|---|---|
| Desempenho | p95 abaixo de 300 ms em leitura e 500 ms em escrita, com massa de 2.000 ativos; segunda rodada com 10.000 ativos para conhecer o limite |
| **7 dias offline** | App operando 7 dias sem conexão, com ações acumuladas, fotos e conflitos, sincronizado sem perda |
| Backup restaurado | Restauração completa executada, com RPO e RTO medidos e registrados |
| Acessibilidade | Relatório automático sem falha de nível AA, mais verificação manual de teclado e leitor de tela |
| Pentest | Nenhum achado crítico ou alto aberto nos quatro focos |
| Sincronização sob carga | Atraso p95 dentro da meta com volume de várias conexões simultâneas |
| Release candidato | Versão `-rc` publicada em homologação pela esteira, sem passo manual |

## Evidência esperada

`docs/screenshots/E11_endurecimento/`: relatório de carga, relatório de acessibilidade, app depois de 7 dias offline, painel de observabilidade com o atraso da sincronização, selo de versão `-rc` em homologação.

## Fora do escopo

Requisito funcional novo. Qualquer coisa das fases 2 e 3. Pentest externo completo, que é do mês 12, na fase 2.

## Decisões pendentes do dono do produto

- Contratação da empresa de pentest, se ainda não feita.
- **P7**, se ainda aberta: sem a decisão, o lançamento publica só o RTO de zona.

## Ajustes durante a etapa

Nenhum até agora.
