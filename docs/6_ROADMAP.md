# Roadmap — Tenaz

**Versão:** 0.6 · **Data:** 13/11/2026 · **Status:** Em elaboração
**Produto:** Tenaz · Gestão de Manutenção e Frota
Alinhado ao sumário executivo 0.8.

**Mudanças da versão 0.6:** o sumário 0.8 tornou a Tenaz **multipaís e bilíngue desde o MVP** (Brasil, Argentina e Paraguai; `pt-BR` e `es-419`; [etapa 014](etapas/014_2026-11-13_multipais-e-idiomas.md), ADR-013) e o documento 1 ganhou RNF-043 a RNF-046 (os não funcionais passam de 42 para **46**). A seção 2.3 estima o impacto no prazo e no time (cerca de 17 a 20 sprints-pessoa no MVP, além do que a E01 já gastou, e cerca de 10 a 13 na Fase 2) e apresenta as opções A a D para o MVP e três para a Fase 2. O roadmap **não decide** a pendência **P15** (preço, moeda de cobrança e impacto no prazo e no time): a tabela da seção 2 e as sprints **não mudam até a decisão do dono do produto**. Entram a tabela de onde o multipaís entra em cada etapa (seção 4), as atividades de apoio externo da Argentina e do Paraguai (seção 7), as dependências 26 e 27, a contingência do portão por país (seção 6), os riscos 29 a 34 e o critério do marco M5. Sem mudança de sprint, prazo, time ou escopo dos requisitos funcionais.

**Mudanças da versão 0.5:** a Tenaz não guarda dado de saúde (sumário 0.7, seção 3.6, decisão do dono do produto em 13/11/2026); o risco 24 deixa de citar saúde entre os dados pessoais fora da lista fechada da conexão. Sem mudança de sprint, prazo ou escopo do cronograma.

Este documento organiza a entrega das três fases definidas no [sumário executivo](0_SUMARIO_EXECUTIVO.md) (seção 4): MVP (meses 0 a 6), Fase 2 (meses 6 a 12) e Fase 3 (meses 12 a 24). O MVP está quebrado em sprints de 2 semanas; as Fases 2 e 3, em blocos mensais. Cada bloco lista os requisitos que entrega, com os IDs de [1_ESPECIFICACAO_REQUISITOS.md](1_ESPECIFICACAO_REQUISITOS.md) (versão 0.7).

A sequência, a duração dos blocos, a equipe e as datas são **proposta**, salvo onde o texto indica a fonte. Nenhuma estimativa vem de velocidade medida: o time ainda não existe.

**Mudanças da versão 0.1:** o sumário 0.3 acrescentou o TER básico ao MVP e o TER completo, o EST estendido e a LOC à Fase 2, e moveu a assinatura avançada, os webhooks e a reserva de estoque para o MVP. O MVP passou de 99 para 105 requisitos funcionais e a Fase 2 de 43 para 62. Para manter os prazos do sumário, o time cresce: o terceiro backend é antecipado para o início do MVP (12 pessoas) e a Fase 2 ganha uma quarta trilha (15 pessoas). Entrou também o RNF-039 (controle de versão e versão visível na interface).

**Mudanças da versão 0.2:** o sumário 0.4 colocou no MVP o módulo **SRV** completo (conta própria da subcontratada, técnicos, clientes, equipamentos de clientes, OS de serviço, componentes, app do técnico, relatório em PDF e conexão com clientes Tenaz). O MVP passa de 105 para **126** requisitos funcionais e os não funcionais de 40 para **42** (RNF-041 e RNF-042, risco R5). O reaproveitamento de OS, apontamento, ativos e app reduz o esforço, mas o acréscimo (cerca de 16 sprints-pessoa de desenvolvimento, seção 2.1) não cabe no time de 12 pessoas em 13 sprints. **Proposta:** antecipar para S1 o quarto backend e o terceiro frontend que já estavam previstos para a Fase 2 (MVP com **14 pessoas**, Fase 2 continua com 15), manter as **13 sprints** e o lançamento na semana 26, abrir uma **trilha SRV** paralela a partir de S1 e liberar a conexão entre contas só com o teste de três contas de RNF-041 verde (riscos R5 e R6 do sumário). Também foram aplicadas as decisões de nomes do sumário 3.6 (por exemplo, `os.encerrar` e o link de uso único com token no fragmento da URL).

**Mudanças da versão 0.3:** o sumário 0.5 tornou a Tenaz **configurável por segmento** (sumário 3.7). Entram no MVP três requisitos: regras configuráveis de condição e efeito (RF-ADM-016), fase de ciclo de vida do ativo com o tipo de OS `implantacao` (RF-ATV-013) e campos extras da OS (RF-OS-021); o MVP passa de 126 para **129** requisitos funcionais. Entram na Fase 2 os pacotes de segmento (RF-ADM-017) e os módulos PRJ (RF-PRJ-001 a RF-PRJ-007) e CTR (RF-CTR-001 a RF-CTR-005); a Fase 2 passa de 62 para **75**. O acréscimo do MVP (cerca de 9 a 10 sprints-pessoa, seção 2.2) consome a folga pequena da versão 0.3; o roadmap **não decide** o tamanho do time nem o prazo e apresenta as opções (absorver nas 13 sprints com o risco descrito, somar uma sprint ou somar um backend). Os três requisitos foram encaixados nas sprints S2 a S9, depois do que reaproveitam. Na Fase 2, PRJ, CTR e os pacotes formam a **trilha E** (proposta), com as opções de equipe e de prazo da seção 2.2. Entram as pendências P9 e P10 nas atividades não técnicas e o risco R7 (configuração virar sistema "faça você mesmo") nos riscos do cronograma.

---

## Índice

1. [Visão geral](#1-visão-geral)
2. [Premissa de equipe](#2-premissa-de-equipe)
3. [Critério de pronto](#3-critério-de-pronto)
4. [MVP — meses 0 a 6](#4-mvp--meses-0-a-6)
5. [Fase 2 — meses 7 a 12](#5-fase-2--meses-7-a-12)
6. [Fase 3 — meses 13 a 24](#6-fase-3--meses-13-a-24)
7. [Atividades não técnicas](#7-atividades-não-técnicas)
8. [Dependências críticas](#8-dependências-críticas)
9. [Riscos do cronograma](#9-riscos-do-cronograma)

---

## 1. Visão geral

```
Mês   0 ─────────────── 6 ─────────────── 12 ────────────────────────────── 24
      MVP (S0 a S12)      Fase 2              Fase 3
      núcleo CMMS,        frota completa,     condição e IoT, agentes de IA,
      NR-12, documentos,  compras e NF-e,     RCM/FMEA, TCO, ISO 55001,
      KPIs, app offline,  estoque em terceiro programação por capacidade,
      assinatura,         e peças reparáveis, MCP, banco dedicado,
      TER básico,         TER completo, LOC,  SOC 2 Type II
      SRV e conexão       NR-13, LOTO,
      entre contas,       capacitações,
      API v1, webhooks,   telemática, ERP,
      SSO OIDC,           SAML/SCIM, IA 1,
      regras, campos      PRJ, CTR e
      extras da OS e      pacotes de
      fase de ciclo       segmento
      de vida
Marcos:
  M0 fundação (sem 2) · M1 cadastros (sem 8) · M2 alpha interno (sem 14)
  M3 piloto fechado (sem 18) · M4 MVP completo (sem 22) · M5 lançamento (sem 26)
  M6 Fase 2 completa (mês 12) · M7 relatório SOC 2 Type II (mês 19) · M8 Fase 3 completa (mês 24)
```

| Fase | Período | Requisitos funcionais | Fonte do escopo |
|---|---|---|---|
| MVP | Meses 0 a 6 (13 sprints; ver as opções da seção 2.2) | 129 (105 + 21 do SRV + 3 da configuração por segmento) | Sumário, seção 4 |
| Fase 2 | Meses 7 a 12 (ver as opções da seção 2.2) | 75 (62 + 13 da configuração por segmento) | Sumário, seção 4 |
| Fase 3 | Meses 13 a 24 | 16 | Sumário, seção 4 |

Requisitos não funcionais: 46 (RNF-001 a RNF-046), verificados nos marcos da seção 3.2. RNF-041 e RNF-042 (conexão entre contas) valem desde o MVP, assim como RNF-043 a RNF-046 (idiomas, países, dinheiro e conteúdo legal por país), que nascem na E01; RNF-046 só tem código a partir da E08.

O mês 0 começa quando o núcleo do time (seção 2) estiver contratado. As datas de calendário serão fixadas nesse momento (proposta).

---

## 2. Premissa de equipe

**Os prazos deste roadmap só valem com o time abaixo.** É um time enxuto para o escopo: 14 pessoas em tempo integral no MVP e 15 a partir da Fase 2 (proposta). Em relação à versão 0.2, o quarto backend e o terceiro frontend, que entrariam na Fase 2, são antecipados para S1 por causa do SRV (seção 2.1); o time da Fase 2 não muda.

A configuração por segmento (sumário 0.5) acrescenta trabalho ao MVP e à Fase 2 que este time cobre só com risco ou com mais prazo; as opções de reforço estão na seção 2.2. **A tabela abaixo não muda até a decisão do dono do produto sobre o tamanho do time.** O escopo multipaís e bilíngue (sumário 0.8) soma-se à configuração por segmento: as opções estão na seção 2.3.

| Papel | MVP | Fase 2 | Fase 3 | Responsabilidade principal |
|---|---|---|---|---|
| Gerente de produto | 1 | 1 | 1 | Prioridades, entrevistas (P4), pilotos, aceite |
| Designer de produto | 1 | 1 | 1 | Fluxos, design system, testes com usuários (RNF-028), fluxos da subcontratada (SRV) |
| Líder técnico (Go) | 1 | 1 | 1 | Arquitetura, revisão, RLS, decisões técnicas, ADR-009 (conexão entre contas) |
| Backend Go | 4 | 4 | 4 | Domínio, API, rotinas, integrações, isolamento do prestador e do locatário; no MVP, um backend dedicado à trilha SRV e à sincronização entre contas |
| Frontend React/TypeScript | 3 | 3 | 3 | Web de gestão, painéis, portal de solicitação, página do link de execução, portal do locatário; no MVP, um frontend dedicado às telas do SRV e da conexão |
| Mobile React Native (Expo) | 2 | 2 | 2 | App offline, sincronização, câmera, QR, perfil prestador, técnico da subcontratada, assinatura do cliente, vistorias |
| QA com automação | 1 | 1 | 1 | Testes de ponta a ponta, offline, isolamento entre empresas, prestadores, locatários e contas conectadas (RNF-041) |
| Plataforma/SRE | 1 | 1 | 1 | Infra em São Paulo, CI/CD, esteira de versionamento (RNF-039), backup, observabilidade (inclusive atraso da sincronização entre contas, RNF-042), SOC 2 |
| Engenharia de dados e IA | 0 | 1 | 1 | Telemática, RAG, IoT, KPIs avançados |
| **Total em tempo integral** | **14** | **15** | **15** | |

Apoio externo em tempo parcial (proposta): advogado (P1, P7, contratos, LGPD) e, para a Argentina e o Paraguai, advogado local (P14); revisor de espanhol falante nativo (seção 2.3); consultor de SST para validar NR-12, NR-13 e LOTO; contador para NF-e, CBS e IBS e para os CFOPs de remessa, conserto e locação (P8), e contador de cada país para o fiscal da Argentina e do Paraguai (P13); consultor de SST de cada país (P12); empresa de pentest; auditor SOC 2; encarregado de dados como serviço.

### 2.1 O MVP cabe em 13 sprints?

**Com o SRV no MVP, não cabe no time de 12 pessoas da versão 0.2. Cabe em 13 sprints com 14 pessoas** (quarto backend e terceiro frontend antecipados para S1), com folga pequena e linha de corte revista (proposta, sem velocidade medida). Com os três requisitos da configuração por segmento (sumário 0.5), essa folga é consumida; as opções estão na seção 2.2. Nas estimativas abaixo, "sprint" é o trabalho de uma pessoa durante uma sprint de 2 semanas.

#### Acréscimo do SRV (sumário 0.4)

O SRV reaproveita a OS (estados, marcos, execução, assinatura, impressão), o apontamento, o cadastro de ativos e medidores, o catálogo de serviços do TER, o estoque básico e o app offline. O que é novo é o lado da subcontratada (cliente atendido, equipamento `de_cliente`, especialidades, serviço a realizar, componentes, distribuição, assinatura do cliente, relatório) e, principalmente, a **conexão entre contas**, que não tem paralelo no restante do produto.

| Bloco | Requisitos | Reaproveita | Backend | Frontend | Mobile |
|---|---|---|---|---|---|
| Conta e plano Serviços | RF-SRV-001 | RF-ADM-001, RF-ADM-009 | 0,5 | 0,5 | — |
| Técnicos, especialidades e clientes | RF-SRV-002 a RF-SRV-004 | RF-ADM-004, RF-ADM-015, RF-CON-005; padrão do cadastro de locatário | 1 | 1 | — |
| Equipamentos de clientes | RF-SRV-005 | Todo o ATV e o MED (só a propriedade `de_cliente` e o vínculo com o cliente são novos) | 0,5 | 0,5 | 0,25 |
| Catálogo, OS de serviço, componentes e horas | RF-SRV-006 a RF-SRV-009 | RF-TER-002, RF-OS-001 a RF-OS-010 | 1,5 | 1 | 0,75 |
| Distribuição, app do técnico e assinatura do cliente | RF-SRV-010 a RF-SRV-012 | RF-OS-007, RF-OS-017, RF-MOB-001 a RF-MOB-006, RF-CON-014 | 0,5 | 0,5 | 1 |
| Relatório em PDF, relatório por cliente e indicadores | RF-SRV-013, RF-SRV-020, RF-SRV-021 | RF-OS-018, RF-IND-001, RF-IND-002, RF-IND-007 | 0,75 | 1 | — |
| **Conexão entre contas** | RF-SRV-014 a RF-SRV-019, RNF-041, RNF-042 | Outbox e idempotência do app (RF-MOB-002), auditoria | **3 a 4** | 1 | 0,5 |
| **Total SRV** | 21 RF, 2 RNF | | **cerca de 8 a 9** | **cerca de 5,5** | **cerca de 2,5** |

Somam-se cerca de 1,5 a 2 sprints de QA (teste de três contas e cenários de sincronização), 1 a 1,5 de design (fluxos da subcontratada e da conexão) e 0,5 de plataforma (papel de banco e monitoração do atraso da sincronização). Sem o reaproveitamento, o mesmo escopo custaria cerca do dobro.

**Efeito.** O backend é o gargalo: 8 a 9 sprints divididas por 3 backends dariam cerca de 3 sprints a mais. Com o time de 12 pessoas, o MVP iria para **15 a 16 sprints** (lançamento entre as semanas 30 e 32), contra as 13 do sumário. Por isso a proposta é:

1. **Time:** antecipar para S1 o quarto backend (dono da trilha SRV e da sincronização entre contas, com o líder técnico no desenho do ADR-009) e o terceiro frontend (telas do SRV e da conexão). Os dois já estavam no plano da Fase 2; o custo adicional é de cerca de 6 meses de duas pessoas. O MVP fica com 14 pessoas; a Fase 2 continua com 15. O terceiro frontend, que conhece React, apoia as telas do app do técnico da subcontratada em React Native (proposta), aliviando o mobile.
2. **Prazo:** mantido em 13 sprints (S0 a S12), com lançamento na semana 26. A capacidade nova (cerca de 10 sprints de backend e 9 de frontend) cobre o SRV e devolve a folga de frontend que a versão 0.2 já não tinha.
3. **Linha de corte:** revista com três itens do SRV que podem ir para o início da Fase 2 sem tirar o valor principal (seção abaixo).
4. **Portão da conexão (R5):** a conexão entre contas (RF-SRV-014 a RF-SRV-019) só é liberada com o teste de três contas de RNF-041 verde e o atraso de RNF-042 medido em homologação. Se o portão não estiver verde no fim de S10 (M4), a conexão sai do lançamento e entra no mês 7, na trilha D, e o SRV é lançado no modo "sem cliente na Tenaz" (relatório da OS em PDF), o que exige aval do dono do produto, porque o sumário 4 põe a conexão no MVP.

**Se o quarto backend e o terceiro frontend não forem contratados até S1**, as alternativas são: (a) MVP em 15 a 16 sprints com 12 pessoas, lançamento entre as semanas 30 e 32; (b) só o quarto backend: MVP em 14 sprints (lançamento na semana 28), ou 13 sprints com a linha de corte aplicada até o item 10; (c) nenhum dos dois e prazo fixo: 13 sprints com a conexão entre contas inteira (RF-SRV-014 a RF-SRV-019) no mês 7 da Fase 2, decisão que contraria o sumário 4 e precisa do dono do produto. Essas alternativas foram calculadas antes da configuração por segmento; sem reforço, cada uma ganha cerca de 1 a 1,5 sprint (seção 2.2).

#### Acréscimo da versão 0.2 (sumário 0.3), já absorvido

O que o sumário 0.3 acrescentou ao MVP, em estimativa grosseira de esforço:

| Mudança | Efeito no MVP | Esforço adicional estimado |
|---|---|---|
| TER básico (RF-TER-001 a RF-TER-005) | 4 requisitos novos (RF-TER-001 substitui RF-EST-010); filtro por fornecedor no servidor e na sincronização; teste automatizado de isolamento do prestador (RNF-037, risco R4); página pública do link de uso único | 1 a 1,5 sprint de backend; 0,5 sprint de frontend; pouco no mobile, que reaproveita o fluxo do técnico |
| Pessoa (RF-ADM-015) | Entidade nova que muda o cadastro de usuários, motoristas e documentos pessoais desde S1 | 0,5 sprint de backend |
| Assinatura avançada (RF-CON-014) no MVP | Antes só a base plugável estava em S8; agora entram reautenticação, hash, manifesto e PDF de evidências | 0,5 sprint de backend |
| Reserva no MVP (RF-EST-006) e inventário na Fase 2 (RF-EST-005) | Troca de um requisito por outro de tamanho parecido | Neutro |
| Controle de versão (RNF-039) | Esteira em S0 e selo de versão na primeira tela | Menos de 0,5 sprint de plataforma e frontend |

Total: cerca de 2 a 2,5 sprints de backend e 0,5 a 1 sprint de frontend a mais, num cronograma que já não tinha folga (risco 2). O terceiro backend, antecipado para S0 na versão 0.2, absorveu o acréscimo de backend; o frontend ficou sem folga em S7 a S10, o que o terceiro frontend antecipado nesta versão resolve.

#### Linha de corte

**Se o time for menor que o da seção 2**, o prazo cresce ou o escopo cai. Com 6 desenvolvedores em vez de 10 (por exemplo, líder técnico, 2 backend, 1 frontend e 2 mobile), estima-se de 11 a 12 meses para o escopo de 126 requisitos da versão 0.3, e cerca de meio mês a mais com os três da configuração por segmento, 129 no total (proposta). Para manter 6 meses com time menor, a linha de corte abaixo sai do MVP e vai para o início da Fase 2, nesta ordem:

| Ordem de corte | Requisito | Impacto de adiar |
|---|---|---|
| 1 | RF-OS-018 Impressão da OS | Baixo; o relatório da OS de serviço (RF-SRV-013) passa a ter gerador próprio de PDF, o que aumenta um pouco o esforço do SRV |
| 2 | RF-PRE-010 Biblioteca de planos padrão | Baixo; onboarding mais lento |
| 3 | RF-SRV-020 Relatório de serviços por cliente | Baixo; o relatório de cada OS (RF-SRV-013) continua; a prestação de contas do período sai da exportação da lista de OS |
| 4 | RF-MED-006 Medidor derivado | Implemento sem km automático |
| 5 | RF-MED-008 Leitura em lote | Leituras uma a uma |
| 6 | RF-TER-002 Catálogo de serviços do fornecedor | O planejador do cliente escolhe o prestador sem busca por classe e marca; o custo previsto da OS de terceiro é digitado. O catálogo próprio da subcontratada (RF-SRV-006) fica |
| 7 | RF-CON-016 Dossiê de fiscalização | Relatórios avulsos continuam disponíveis |
| 8 | RF-SRV-017 OS aberta pela subcontratada nasce no cliente | A subcontratada registra a falha encontrada na observação da OS e no relatório; o cliente abre a OS por conta própria |
| 9 | RF-TER-005 Link de execução de uso único | Prestador ocasional não registra a execução; um usuário interno lança por ele, ou o prestador recebe login e passa a contar como técnico |
| 10 | RF-SRV-019, só suspensão e reativação | A conexão só pode ser encerrada, não suspensa; o encerramento nunca sai do MVP |
| 11 | RF-SOL-004 Portal de solicitação por QR | Solicitante precisa de login |
| 12 | RF-SOL-005 Prazos de atendimento | MTTA calculado, sem alerta de prazo |
| 13 | RF-INT-002 Webhooks | Integrações por consulta periódica à API |
| 14 | RF-ATV-004 Composição de veículos | Afeta transportadoras; cortar só se o piloto não tiver carretas |

Nunca saem do MVP: isolamento entre empresas (RNF-013), auditoria (RF-ADM-011), offline (RF-MOB-002 e RF-MOB-003), códigos de falha com catálogo padrão (RF-OS-011), assinatura avançada (RF-CON-014), NR-12 (RF-CON-001 a RF-CON-004), documentos com alertas (RF-CON-005 e RF-CON-006), "como calculamos" (RF-IND-002), o controle de versão (RNF-039 e RNF-040) e o núcleo do SRV que funciona sem cliente na Tenaz (RF-SRV-001 a RF-SRV-013), que sustenta o plano Serviços. São os diferenciais do sumário (seção 1.2), condição de operação ou linha de receita. Dois portões de isolamento valem como regra: o perfil `prestador` (RF-TER-004) só é liberado com o teste de RNF-037 verde (sem ele, as OS de terceiro são registradas por usuário interno), e a conexão entre contas (RF-SRV-014 a RF-SRV-019) só é liberada com o teste de três contas de RNF-041 verde e o atraso de RNF-042 medido (sem eles, a conexão vai para o mês 7 e o SRV é lançado sem conexão, com aval do dono do produto).

Os três requisitos da configuração por segmento (RF-ADM-016, RF-ATV-013 e RF-OS-021) **não entram na linha de corte** nesta versão: o sumário 4 os põe no MVP. Pôr na linha de corte partes deles (seção 2.2, opção A) é decisão do dono do produto.

### 2.2 Configuração por segmento (sumário 0.5)

O sumário 0.5 (seção 3.7) acrescenta ao MVP as regras configuráveis, a fase de ciclo de vida com o tipo de OS `implantacao` e os campos extras da OS, e à Fase 2 os módulos PRJ e CTR e os pacotes de segmento. Tudo nesta seção é **(proposta)**: estimativa sem velocidade medida, como na seção 2.1. O roadmap mostra as opções e **não decide** o tamanho do time nem o prazo, que são pendência do dono do produto (sumário 6 e [CLAUDE.md](../CLAUDE.md), seção 5).

#### Acréscimo no MVP

| Bloco | Requisitos | Reaproveita | Backend | Frontend | Mobile |
|---|---|---|---|---|---|
| Fase de ciclo de vida e OS de implantação | RF-ATV-013; tipo `implantacao` de RF-OS-002 | Status do ativo e `inativo` (RF-ATV-006), transferência (RF-ATV-011), importação (RF-INT-003), suspensão de planos (RF-PRE-009), documentos críticos (RF-CON-005), exclusão dos KPIs (RF-IND-003, RF-IND-005) | 0,75 a 1 | 0,5 | 0,25 |
| Campos extras da OS | RF-OS-021 | Padrão dos campos extras por classe (RF-ATV-005), OS (RF-OS-001), impressão (RF-OS-018), link de uso único (RF-TER-005), relatório da OS de serviço (RF-SRV-013), cache de configuração do app (RF-MOB-002) | 1 | 1 | 0,75 |
| Regras configuráveis | RF-ADM-016 | `bloqueio_ativo` (RF-ATV-006), notificações (RF-ADM-012), solicitação com origem (RF-SOL-001), alertas de documento (RF-CON-006), item reprovado (RF-CHK-003), validação de leitura (RF-MED-003), rotinas diárias idempotentes, auditoria (RF-ADM-011) | 2 a 2,5 | 1 | — |
| **Total** | 3 RF | | **cerca de 3,75 a 4,5** | **cerca de 2,5** | **cerca de 1** |

Somam-se cerca de 1 a 1,5 sprint de QA (combinações de condição e efeito, idempotência da rotina diária, campo obrigatório na conclusão offline, efeito da regra chegando à conta conectada só pelo status), 0,5 de design (editor de regra por catálogo, campos dinâmicos na OS) e quase nada de plataforma (a rotina diária usa o agendador existente). Total: **cerca de 9 a 10 sprints-pessoa**. As regras rodam no servidor; o app só valida os campos extras e mostra o ativo `parado` e a fase "Em implantação".

Pelo critério de pronto 3.1, item 2, todo campo novo de OS ou de ativo é classificado na lista fechada da conexão (RF-SRV-014). `fase_ciclo_vida` e `data_entrada_operacao` (RF-ATV-013) e `ordem_servico.campos_extras` (RF-OS-021, sumário 3.6) já estão definidos como "nunca compartilhado" no lançamento. **`ordem_servico.causa_externa` ainda não tem classificação** no sumário nem nos requisitos: precisa dela antes de RF-OS-001 fechar com o campo, em S4; até lá, vale "nunca compartilhado", o padrão mais seguro (proposta), e a decisão volta ao sumário. O efeito de `bloquear_ativo` chega à conta conectada só pelo status já compartilhado, verificado no teste de três contas de RNF-041.

**Efeito.** O trabalho cai nos três backends do núcleo (o quarto backend é da trilha SRV e o terceiro está no TER de S7 a S9): cerca de 1,25 a 1,5 sprint a mais por backend do núcleo, mais do que a folga pequena da versão 0.3. O mobile, já carregado pelo SRV (risco 23), ganha cerca de meia sprint por pessoa.

**Opções para o MVP (a decidir pelo dono do produto):**

| Opção | Time | Prazo | Risco e custo |
|---|---|---|---|
| **A. Absorver nas sprints atuais** | 14 pessoas | 13 sprints, lançamento na semana 26 | Sem folga em S6, S8 e S9; qualquer atraso leva à linha de corte, cuja aplicação é do dono do produto; mobile e QA no limite. Para reduzir o risco, partes dos três requisitos poderiam entrar na linha de corte, sem tirar o valor principal: simulação de regra e histórico de disparos por ativo (RF-ADM-016), colocação em operação em lote (RF-ATV-013), filtro e exportação por campo extra e preenchimento pelo link de uso único (RF-OS-021). Economia estimada: cerca de 1 sprint de backend e 0,5 de frontend. Incluir esses itens na linha de corte exige aval do dono do produto |
| **B. Somar 1 sprint** | 14 pessoas | 14 sprints, lançamento na semana 28 (M4 na semana 24, M5 na semana 28) | Custo de duas semanas do time inteiro. A sprint a mais entra entre S9 e S10 e recebe o que S6 a S9 não fecharem, a começar por RF-ADM-016 completo; S10 continua sendo a última sprint com RF novo, como pede a linha de corte. Se a revisão de escopo do fim de S6 (risco 2) mostrar velocidade abaixo da prevista, pode virar **2 sprints** (semana 30) |
| **C. Somar 1 backend** | 15 pessoas no MVP | 13 sprints, lançamento na semana 26 | Um backend a mais de S4 a S10 no núcleo (cerca de 3,5 meses de uma pessoa); entrada no meio do MVP custa integração nas primeiras sprints; frontend e mobile continuam absorvendo 2,5 e 1 sprint. Se a Fase 2 receber a trilha E com reforço (abaixo), a mesma pessoa segue nela |

As opções se combinam (por exemplo, B com os cortes parciais de A). Tirar do MVP algum dos três requisitos inteiros contraria o sumário 4 e o roadmap não propõe. Em qualquer opção, o encaixe nas sprints é o da seção 4 e o das etapas E02 a E09 em [execucao/4_ETAPAS.md](execucao/4_ETAPAS.md); na opção B, o que deslizar para a sprint extra é decidido na revisão do fim de S6.

#### Acréscimo na Fase 2

| Bloco | Requisitos | Reaproveita | Backend | Frontend | Mobile |
|---|---|---|---|---|---|
| Projetos de implantação | RF-PRJ-001 a RF-PRJ-007 | Fase de ciclo de vida e `colocar_em_operacao` (RF-ATV-013), OS `implantacao` e campos extras (RF-OS-002, RF-OS-021), fornecedor e OS de terceiro (RF-TER-003), OS na conta conectada (RF-SRV-016), documentos (RF-CON-005), assinatura (RF-CON-014), "como calculamos" (RF-IND-002) | 3 a 3,5 | 2,5 | 0,25 |
| Contratos com clientes | RF-CTR-001 a RF-CTR-005 | Cliente atendido (`cliente_servico`, RF-SRV-004), registros de parada (RF-OS-005), `causa_externa` gravada desde o MVP (RF-OS-001), disponibilidade técnica (RF-IND-003), exportação em PDF e XLSX | 2 a 2,5 | 1,5 | 0,25 |
| Pacotes de segmento | RF-ADM-017 | Regras (RF-ADM-016), campos extras (RF-OS-021), classes (RF-ATV-005), planos padrão (RF-PRE-010), modelos de projeto (RF-PRJ-001), roteiro de primeiros passos (RF-ADM-001) | 1 | 0,5 | — |
| **Total** | 13 RF | | **cerca de 6 a 7** | **cerca de 4,5** | **cerca de 0,5** |

Somam-se cerca de 1 a 1,5 sprint de QA, 1 de design (cronograma do projeto, painel do contrato) e cerca de 1 sprint do gerente de produto com o domínio de manutenção para montar o conteúdo dos seis pacotes. Total: **cerca de 13 a 15 sprints-pessoa**. A Fase 2 já está apertada (risco 19).

**Opções para a Fase 2 (a decidir pelo dono do produto):**

| Opção | Time | Prazo | Risco e custo |
|---|---|---|---|
| **1. Trilha E com reforço** | 17 pessoas (mais 1 backend e 1 frontend a partir do mês 7; se a opção C do MVP for escolhida, o backend já está no time) | M6 mantido no mês 12 | Custo de duas pessoas por 6 meses; a capacidade que sobrar na trilha E absorve a evolução da configuração pedida pelos pilotos e reforça a trilha D |
| **2. Sem reforço, com mais prazo** | 15 pessoas | Fase 2 até o mês 13 ou 14; M6, a Fase 3 e o M8 deslizam o mesmo tanto | PRJ, CTR e pacotes divididos entre as trilhas A e D, na mesma ordem da seção 5; a janela do SOC 2 não muda, porque não depende desses módulos |
| **3. Sem reforço e prazo fixo** | 15 pessoas | M6 no mês 12 | Exige mover requisitos da Fase 2 para a Fase 3 além da válvula do risco 19 (RF-LOC-007 e RF-LOC-009); qualquer mudança de fase é decisão do dono do produto |

### 2.3 Multipaís e idiomas (sumário 0.8, ADR-013)

O sumário 0.8 (seção 3.6, "Países e idiomas") torna a Tenaz **multipaís desde o MVP** (Brasil, Argentina e Paraguai) e **bilíngue desde a E01** (`pt-BR` e `es-419`), e a [etapa 014](etapas/014_2026-11-13_multipais-e-idiomas.md) os transforma nos requisitos RNF-043 a RNF-046 ([documento 1](1_ESPECIFICACAO_REQUISITOS.md), seção 5.12). Nenhum requisito funcional foi acrescentado nem mudou de fase: o multipaís é um custo **transversal**, pago em toda etapa que tem tela, texto, dinheiro ou conteúdo legal. Tudo nesta seção é **(proposta)**, estimativa sem velocidade medida, como nas seções 2.1 e 2.2. O roadmap mostra a estimativa e as opções e **não decide** o prazo nem o time: é a pendência **P15** do sumário, do dono do produto ([CLAUDE.md](../CLAUDE.md), seção 5).

#### O que já foi feito na E01

A base está no código e nos testes da E01 em andamento: país, identificação por tipo, moeda e dinheiro (Parte F do ADR-013), idioma por usuário e por empresa, textos do servidor e da interface nos dois idiomas, fuso e endereço por país, suíte de isolamento com empresas dos três países (rastreio em RNF-043 a RNF-045). Esse esforço foi real e **saiu da folga da E01**, que já não era grande (estimativa grosseira: cerca de 5 a 7 sprints-pessoa, quase todo em backend e frontend). O consumo real da E01 é o primeiro dado de velocidade do multipaís e deve ser registrado pelo gerente de produto antes da decisão abaixo.

#### Acréscimo nas etapas E02 a E12

| Bloco | Etapas | Backend | Frontend | Mobile |
|---|---|---|---|---|
| **Textos em dois idiomas** em toda tela, e-mail, notificação e erro novos: catálogo `es-419`, glossário, teste de paridade, revisão pelo glossário bilíngue; o espanhol ocupa mais espaço na tela e pede formulários por país | E02 a E12 | 1 | 2 | 1 |
| **Dinheiro com moeda** em toda tabela e indicador novos (RNF-045): valor do ativo, custos da OS, estoque e compras, KPIs em dinheiro | E02, E04, E05, E07, E09 | 1 | 0,5 | 0,25 |
| **Documentos gerados e notificações no idioma** (OS impressa, relatório da OS de serviço, exportação, push) | E05, E09, E10 | 1 | 0,25 | 0,25 |
| **Prefixo `OS-` ou `OT-`, feriados nacionais e formatação no perfil de país** | E02, E04 | 0,5 | 0,25 | — |
| **Conformidade por país** (RNF-046): `registro_conformidade` no lugar de `registro_nr12`, mecanismo de pacote de país aplicado na criação da conta, pacote `BR` com os prazos hoje em `CHECK` e comentários, jurisdição emissora no documento | E08 | 1,5 a 2 | 0,5 | 0,25 |
| **Identificação de qualquer país** em fornecedor, cliente atendido, locatário e conexão entre contas; "mesma entidade" pelo perfil do país; teste de três contas com países diferentes | E07, E09, E10 | 0,75 | 0,25 | — |
| **Cobrança em `ARS` e `PYG`** (RF-ADM-009; depende de **P15**) | E10 | 0 a 1,5 | 0 a 0,5 | — |
| **Total** | | **cerca de 6 a 8** | **cerca de 4 a 4,5** | **cerca de 1,75 a 2** |

Somam-se cerca de 2,5 a 3 sprints de QA (testes de ponta a ponta nos dois idiomas, vazamento com os três países em toda tabela nova, carga e pentest nos dois idiomas), cerca de 1 de design (formulários por país, espanhol mais longo, as telas de idioma), cerca de 0,5 de plataforma (tarefa `textos` no CI, sinais por país na observabilidade, aditivo por país) e cerca de 1 do gerente de produto (glossário, revisão, conteúdo dos pacotes). Total: **cerca de 17 a 20 sprints-pessoa**, mais que a configuração por segmento (9 a 10, seção 2.2) e na mesma ordem do acréscimo do SRV já descontado o reaproveitamento (cerca de 16, seção 2.1). O trabalho cai nos mesmos gargalos de sempre: os backends do núcleo, o QA e o mobile. Como a folga já foi consumida pela configuração por segmento, **nada deste acréscimo cabe de graça nas 13 sprints com 14 pessoas**.

**Apoio externo novo (proposta, tempo parcial):** advogado com experiência na Argentina e no Paraguai (**P14**, aditivo de transferência de dados por país, termo da conexão entre contas de países diferentes); consultor de SST de cada país (**P12**, equivalentes da NR-12 e da NR-13, prazos por jurisdição); contador de cada país (**P13**, identificação fiscal e fatura eletrônica, antes do fiscal da Fase 2); revisor de espanhol **falante nativo** (um argentino e um paraguaio), por release, para os termos que o glossário bilíngue marcou. A Tenaz não tem hoje nenhum deles contratado; o prazo de contratação entra no cronograma (risco 29).

**Opções para o MVP (a decidir pelo dono do produto, junto com a **P15**):**

| Opção | Time | Prazo | Risco e custo |
|---|---|---|---|
| **A. Absorver nas sprints atuais** | 14 pessoas | 13 sprints, semana 26 | Não recomendável: a folga já é zero (seção 2.2). Só fecha se partes do multipaís e da configuração por segmento forem para a linha de corte, o que é decisão do dono do produto, e com QA e mobile além do limite. Sem aval, qualquer atraso cai na linha de corte da seção 2.1 |
| **B. Somar sprints** | 14 pessoas | Mais **2 sprints** sobre a opção escolhida em 2.2: com a opção A de 2.2, 15 sprints (semana 30); com a B de 2.2, 16 sprints (semana 32) | Custo de quatro semanas do time inteiro. As 17 a 20 sprints-pessoa dariam cerca de 1,3 sprint se tudo corresse em paralelo; o gargalo do backend e do QA leva a 2. O que escorrega é decidido na revisão do fim de S6 (risco 2) |
| **C. Somar 1 backend e 1 QA** | 16 pessoas no MVP (17 com a opção C de 2.2) | 13 sprints, semana 26 | O backend cobre cerca de 11 sprints-pessoa do núcleo (textos, dinheiro, conformidade por país); o QA dobra uma capacidade que hoje é de uma pessoa para dois idiomas e três países. Custo de duas pessoas por cerca de 5 meses, com entrada em S2 ou S3 e custo de integração. Decisão até o fim de S3, como a da configuração por segmento |
| **D. Portão por país** (combina com A, B ou C) | O da opção escolhida | O Brasil lança em M5; Argentina e Paraguai saem à venda quando o portão do país fechar | O código multipaís e o espanhol, já decididos e feitos na E01, ficam inteiros; a **venda** em cada país só começa com: parecer jurídico da **P14**, o que couber da **P12** (documentos sem prazo legal preenchido, como o ADR-013 já prevê), preço e moeda da **P15**, e revisão por falante nativo concluída. Tira do caminho crítico de M5 os pacotes de país e a cobrança em `ARS` e `PYG` (itens 5 e 7 da tabela, cerca de 2 a 4 sprints-pessoa). **Contraria a leitura literal de "MVP já multipaís" na data de venda, não no código:** precisa de aval do dono do produto |

As opções se combinam (por exemplo, C com D). Cortes parciais candidatos, sem aplicação automática: documento gerado e exportação em espanhol só na etapa seguinte ao M5; notificação por push só no idioma do usuário, sem texto específico do país; cobrança em `ARS` e `PYG` por fatura manual no lançamento. Tirar o idioma ou o país do MVP inteiro contraria a decisão do dono do produto e o roadmap não propõe.

#### Acréscimo na Fase 2

| Bloco | Requisitos | Backend | Frontend | Mobile |
|---|---|---|---|---|
| Fiscal da Argentina e do Paraguai: referência do documento fiscal do país no lugar de NF-e, CFOP e chave de 44 dígitos (**P13**) | RF-EST-009, RF-EST-015, RF-EST-019, RF-LOC-002, RF-LOC-009 | 1,5 a 2 | 0,5 | — |
| Pacotes de país de AR e PY: documentos, normas, capacitações e prazos por jurisdição, só depois da **P12** | RF-CON-005 a RF-CON-013, RF-FRO-001 a RF-FRO-011 | 2 a 3 | 1 | 0,25 |
| Textos em dois idiomas nas telas da Fase 2 | Toda a Fase 2 | 0,5 | 1,5 a 2 | 0,5 |
| **Total** | | **cerca de 4 a 5,5** | **cerca de 3 a 3,5** | **cerca de 0,75** |

Somam-se cerca de 1,5 a 2 sprints de QA e cerca de 1,5 do gerente de produto com os consultores dos países: **cerca de 10 a 13 sprints-pessoa**, sobre uma Fase 2 que já está apertada (risco 19). O conteúdo legal da Argentina e do Paraguai depende de pesquisa e parecer externo; sem eles, os pacotes nascem sem prazo preenchido (RNF-046). **Opções:** (1) um reforço de 1 backend e 0,5 de QA na Fase 2, dedicado aos países; (2) mais prazo (M6 no mês 13 ou 14), somando à opção 2 da seção 2.2; (3) Argentina e Paraguai entrarem na Fase 2 país a país, nos portões da opção D. A escolha é do dono do produto e anda junto com a do MVP.

#### Onde entra cada parte do multipaís (proposta)

| Etapa | O que o multipaís pede |
|---|---|
| E01 | Base entregue em parte (acima); fecha com o aceite do dono do produto, a revisão por falante nativo dos termos marcados e o app com o idioma do usuário |
| E02 | Valor do ativo com moeda (RNF-045); feriados nacionais e formatação no perfil de país; telas em dois idiomas|
| E04 | Prefixo `OS-` ou `OT-` do número da OS (RF-OS-001, ponto 49 do documento 1); custo da OS em moeda; solicitações e erros nos dois idiomas |
| E05 | OS impressa e notificações no idioma de quem lê; app com os textos no bundle |
| E07 | Fornecedor com identificação de qualquer país; estoque e compras em moeda; termo da conexão entre contas de países diferentes (**P14**) |
| E08 | Conformidade por país (RNF-046): pacote `BR`, `registro_conformidade`, jurisdição emissora; consultor de SST do Brasil em S8, e os de AR e PY pelo portão por país |
| E09 | Indicadores em dinheiro com a moeda; relatório da OS de serviço no idioma do cliente atendido |
| E10 | Cobrança em `ARS` e `PYG` (**P15**); exportação e extração da trilha no idioma; documentação pública da API em espanhol **a decidir** (o pedido do dono do produto cobre telas, e-mails, notificações, documentos gerados e erros, não a documentação) |
| E11 | Carga, acessibilidade, 7 dias offline e pentest com os dois idiomas e os três países; revisão final do espanhol |
| E12 | Kit LGPD com aditivo por país (**P14**), página de preços por país (**P15**), base de conhecimento e suporte em espanhol **a decidir**, pilotos por país conforme o portão |

---

## 3. Critério de pronto

### 3.1 Para cada requisito

Um requisito só conta como entregue quando:

1. Validações e regras descritas no requisito têm testes automatizados.
2. Toda tabela nova tem RLS forçado e entra no teste de vazamento entre empresas que roda a cada build (RNF-013). Tabela lida pelos perfis `prestador` ou `locatario` entra também nos testes de RNF-037 ou RNF-038. Tabela ou campo lido ou gravado pelo processo de sincronização entre contas entra no teste de três contas de RNF-041, e todo campo novo de OS ou de ativo é classificado como "compartilhável" ou "nunca compartilhado" na lista fechada da conexão (RF-SRV-014).
3. Endpoints estão na especificação OpenAPI e são os mesmos usados pela web e pelo app.
4. Alterações relevantes geram eventos de auditoria (RF-ADM-011).
5. Se o requisito prevê uso offline, há teste de cenário com desconexão, ação, reconexão e conflito.
6. As telas novas passam por verificação automática e manual de acessibilidade (RNF-025) e não escondem o selo de versão (RNF-039).
7. O p95 fica dentro de RNF-007 em homologação com massa de 2.000 ativos.
8. Há ajuda contextual mínima na tela.
9. O gerente de produto aceitou em homologação. Requisitos de CON também passam pelo consultor de SST; requisitos com NF-e de remessa, conserto ou locação passam pelo contador (P8).

### 3.2 Para cada marco

| Marco | Quando | Critério de pronto |
|---|---|---|
| M0 Fundação | Fim de S0 (semana 2) | CI/CD verde; infraestrutura em São Paulo criada por código; teste de vazamento entre empresas rodando no CI; esteira de versionamento de RNF-039 ativa (versão gerada da tag, CHANGELOG verificado, release bloqueado em divergência, `GET /v1/versao` respondendo); build interno do app instalável |
| M1 Cadastros | Fim de S3 (semana 8) | Empresa de teste com 2.000 ativos importados por planilha em menos de 2 minutos; leituras validadas com regressão, teto de salto e troca de medidor; pessoas com e sem usuário cadastradas; conta de subcontratada de teste com clientes, equipamentos de clientes (`de_cliente`) e técnicos com especialidades |
| M2 Alpha interno | Fim de S6 (semana 14) | Ciclo completo solicitação → OS → execução offline → encerramento, com reconhecimento, `assumir`, desprogramação e liberação recusada; preventiva por tempo, medidor e "o que vencer primeiro" gerando OS `planejada`; app operando 72 horas offline sem perda; OS de serviço aberta com serviço a realizar, distribuída, reatribuída e executada offline por técnico de subcontratada, com horas por técnico; teste de três contas de RNF-041 rodando no CI desde S4 |
| M3 Piloto fechado | Fim de S8 (semana 18) | Pelo menos 3 empresas piloto com dados reais; OS de máquina NR-12 com os sete campos e assinatura avançada; documentos com alertas D-60, D-30 e D-7 funcionando; reserva e baixa de peças em OS; pelo menos uma subcontratada piloto com conta própria registrando OS de serviço com componentes substituídos; convite, aceite e ativos sincronizados entre duas contas em homologação |
| M4 MVP completo | Fim de S10 (semana 22; semana 24 na opção B da seção 2.2) | Os 129 requisitos do MVP entregues; ativo cadastrado `em_implantacao`, fora dos KPIs, posto em operação por `colocar_em_operacao` depois da OS `implantacao` encerrada; regra `bloquear_ativo_ate_documento` deixando o ativo `parado` e levantada sozinha pelo documento novo; campo extra obrigatório na conclusão impedindo `concluir` também no app offline; OS executada por usuário `prestador` e por link de uso único; teste de isolamento do prestador (RNF-037) no CI sem vazamento; relatório da OS em PDF com assinatura do cliente; ciclo completo entre contas (OS do cliente chega à subcontratada, execução offline, andamento de volta, liberação e recusa no cliente; OS aberta pela subcontratada nasce no cliente; encerramento da conexão sem perda de histórico); **portão da conexão:** teste de três contas de RNF-041 verde e atraso p95 < 60 s de RNF-042 medido em homologação; documentação pública da API; nenhum defeito crítico aberto |
| M5 Lançamento | Fim de S12 (semana 26; semana 28 na opção B da seção 2.2) | Metas de RNF verificadas (carga, 7 dias offline, restauração de backup testada, atraso da sincronização entre contas sob carga); pentest sem achado crítico ou alto aberto, incluindo acesso do prestador, link de uso único e conexão entre contas; kit LGPD publicado, com o termo da conexão entre contas e, para cada país liberado à venda, o aditivo de transferência de dados (**P14**); RNF-043 a RNF-046 verificados nos dois idiomas e nos três países, com a revisão por falante nativo dos termos marcados concluída e o portão do país fechado (seção 2.3, opção D, se for a escolhida); página de status no ar; RPO e RTO publicados, com o RTO de região só depois da decisão de **P7**; preços publicados após P4; marca liberada (P5); release de produção `v1.0.0` conforme a seção 3.3 |
| M6 Fase 2 completa | Fim do mês 12 (mês 13 ou 14 na opção 2 da seção 2.2) | Os 75 requisitos da Fase 2 entregues; projeto de implantação concluído colocando os ativos em operação; medição mensal de contrato com cliente fechada com memória de cálculo; pacote de segmento aplicado sem sobrescrever o que a conta já tinha; testes de isolamento do locatário (RNF-038) e do almoxarifado em terceiro (RNF-037) no CI; pentest externo completo, incluindo portal do locatário; controles SOC 2 em operação e janela de observação iniciada |
| M7 SOC 2 Type II | Mês 19 | Relatório Type II emitido pelo auditor |
| M8 Fase 3 completa | Fim do mês 24 | Os 16 requisitos da Fase 3 entregues; preparação ISO 27001 iniciada |

### 3.3 Para cada release em produção

Vale para o lançamento (M5) e para toda entrega posterior em produção, conforme **RNF-039** e **RNF-040**:

1. O release tem tag git `vX.Y.Z`, com incremento conforme SemVer 2.0.0: MAJOR para quebra de compatibilidade da API pública, MINOR para funcionalidade nova, PATCH para correção.
2. O `CHANGELOG.md`, no formato Keep a Changelog, tem a seção da versão com **todas** as implementações (com ID do requisito) e correções (com número do chamado) incluídas no release, e as notas aparecem em "Novidades desta versão".
3. Web, app, API e workers do release têm a mesma versão, gerada pelo build a partir da tag, nunca digitada; o CI bloqueia o release se tag, CHANGELOG e versão embutida divergirem.
4. A versão aparece no canto inferior direito da tela principal e da tela de login da web e do app, e "Sobre" mostra versão da web ou do app, versão da API, data do build, commit curto e ambiente.
5. `GET /v1/versao` devolve a versão, o commit e a data do build do release publicado.
6. Versões de pré-lançamento (`-rc.N`) só existem fora da produção.

---

## 4. MVP — meses 0 a 6

Treze sprints de 2 semanas (S0 a S12). Web e app avançam em paralelo; a coluna "App" mostra o que o time mobile entrega na mesma sprint. O terceiro backend conduz, a partir de S7, a trilha de terceiros (TER) e a assinatura avançada e, em S10, reforça a conexão entre contas. A **trilha SRV** (quarto backend e terceiro frontend, com o líder técnico no desenho da conexão) corre em paralelo desde S1, na coluna própria.

| Sprint | Semanas | Tema | Requisitos entregues (web e servidor) | Trilha SRV (web e servidor) | App (MOB) | Depende de |
|---|---|---|---|---|---|---|
| S0 | 1–2 | Fundação técnica | Base de RNF-013, RNF-016, RNF-020 e RNF-032; **esteira de versionamento de RNF-039 e RNF-040** (versão gerada da tag `vX.Y.Z`, `CHANGELOG.md` Keep a Changelog, CI rejeita pull request sem entrada no CHANGELOG, bloqueio de release no CI, `GET /v1/versao`) (sem RF de negócio) | Esquema `plataforma`, separado do RLS por empresa, previsto na infraestrutura como código; rascunho do ADR-009 (líder técnico) | Projeto Expo, SQLite, prova de conceito da fila de comandos (outbox), versão do app gerada pelo build | Time contratado |
| S1 | 3–4 | Conta e acesso | RF-ADM-001 (sem o roteiro de primeiros passos), RF-ADM-002, RF-ADM-004, RF-ADM-005 (com `os.encerrar`, `est.decidir_destino` e escopos `interno`, `prestador`, `locatario`, `link`), RF-ADM-006, RF-ADM-011, RF-ADM-015; **selo de versão e tela "Sobre" na web** (RNF-039), na tela de login e na tela principal | **Conta da subcontratada:** RF-SRV-001 parcial (empresa com o módulo SRV, perfis da conta, escolha de empresa no login) | RF-MOB-001, com selo de versão e "Sobre" no app (RNF-039) | S0; quarto backend e terceiro frontend contratados |
| S2 | 5–6 | Ativos | RF-ATV-001, RF-ATV-002, RF-ATV-003, RF-ATV-005, RF-ATV-006, RF-ATV-007, RF-ATV-008; RF-INT-003 para localizações e ativos; **RF-ATV-013 parcial** (`fase_ciclo_vida`, ativo `em_implantacao` com status `inativo` e "Em implantação" na interface, cadastro e importação já em implantação, filtro por fase) | **Clientes:** RF-SRV-004 (clientes atendidos); RF-SRV-002 (especialidades de técnico) | RF-MOB-002 (cache de ativos, com a fase do ativo), RF-MOB-004 | S1 |
| S3 | 7–8 | Medidores e composição | RF-MED-001 a RF-MED-008, RF-ATV-004, RF-ATV-009, RF-ATV-011; RF-ATV-010 parcial | **Equipamentos e técnicos:** RF-SRV-005 (equipamentos de clientes, `de_cliente`); RF-SRV-003 parcial (técnicos com especialidades e capacitações; documentos em S9) | RF-MOB-003 parcial (leitura offline); cache de equipamentos de clientes | S2 (ativos) |
| S4 | 9–10 | Solicitações e núcleo da OS | RF-SOL-001 a RF-SOL-005, RF-OS-001 a RF-OS-005, RF-OS-016, RF-OS-017, RF-ADM-012; origens `projeto` e `regra` e `causa_externa` gravados na OS (RF-OS-001), sem efeito em indicador; **RF-OS-021 parcial** (definições de `campo_extra_os`, valores em `ordem_servico.campos_extras`, obrigatoriedade na abertura, web e API); **RF-ATV-013 parcial** (tipo `implantacao` de RF-OS-002, ação `colocar_em_operacao`, também em lote, com o evento `ativo.colocado_em_operacao`) | **Fundação da conexão** (sem RF): ADR-009 fechado; `conexao_empresa` no esquema da plataforma; papel de banco do processo de sincronização; outbox nos dois lados; idempotência; teste de três contas rodando no CI (base de RNF-041 e RNF-042) | RF-MOB-006; solicitação offline (RF-MOB-003); campos extras na abertura de OS offline, com as definições no cache de configuração; primeira versão nas lojas em teste fechado | S3 |
| S5 | 11–12 | Execução da OS | RF-OS-006, RF-OS-007 (com a ação `assumir`), RF-OS-008, RF-OS-009, RF-OS-011, RF-OS-015, RF-OS-018; **RF-OS-021 parcial** (obrigatoriedade na conclusão, campos extras na impressão da OS) | **OS de serviço:** RF-SRV-006 (catálogo próprio, entidade `servico_fornecedor` que RF-TER-002 reaproveita em S7); RF-SRV-007 (abertura com equipamento, técnico e serviço a realizar) | Execução de OS offline (RF-MOB-003), RF-MOB-005; OS de serviço com serviço a realizar no app; campos extras na execução, com a obrigatoriedade na conclusão validada antes de enfileirar o comando | S4 |
| S6 | 13–14 | Preventiva | RF-PRE-001 a RF-PRE-006, RF-PRE-008, RF-PRE-009, RF-PRE-010, RF-ADM-010; **RF-ATV-013 parcial** (planos suspensos no ativo `em_implantacao` e ativados a partir da entrada em operação); **RF-ADM-016 parcial** (motor de `regra_automacao` e `disparo_regra`, rotina diária idempotente, permissão `adm.gerenciar_regras`, condições `leitura_fora_faixa` e `os_concluida_com_campo`, efeitos `notificar`, `abrir_solicitacao` e `bloquear_ativo`) | RF-SRV-009 (início, fim e apontamentos por técnico); RF-SRV-010 (distribuição e reatribuição); RF-TER-001 antecipado para a trilha SRV, porque o convite parte do fornecedor | OS preventiva no app; teste de 72 h offline; tarefas do técnico da subcontratada (RF-SRV-011 parcial) | S3 (medidores), S5 |
| S7 | 15–16 | Checklists, estoque básico e fornecedores | RF-CHK-001 a RF-CHK-006, RF-SOL-006, RF-EST-001 a RF-EST-004, RF-EST-006, RF-OS-010, RF-FRO-001, RF-TER-002 | **Conexão:** RF-SRV-014 (convite, aceite e campos compartilhados) | Checklist offline, RF-MOB-007 | S5, S6; RF-TER-001 (S6) |
| S8 | 17–18 | NR-12, documentos e assinatura | RF-CON-001 a RF-CON-006, RF-CON-014, RF-CON-016, RF-OS-013 (com o custo da OS `implantacao` separado do custo de manutenção), RF-OS-014; **RF-ATV-013 completo** (documentos críticos da classe exigidos para `colocar_em_operacao`, depois de RF-CON-005) | **Componentes:** RF-SRV-008 (substituído, consumido, fornecido pelo cliente, com números de série); **conexão:** RF-SRV-015 (ativos cobertos) | Campos NR-12 e assinatura avançada no app; componentes offline | S7 (RF-OS-010, RF-SRV-014); **parecer P1** |
| S9 | 19–20 | Indicadores e terceiros | RF-IND-001 a RF-IND-005 (sem OS `implantacao` nem ativo `em_implantacao` nos KPIs de manutenção), RF-IND-007; RF-ATV-010 completo; RF-TER-003, RF-TER-004, RF-TER-005; teste de isolamento do prestador no CI (RNF-037); **RF-ADM-016 completo** (condições `documento_vencido`, `documento_a_vencer` e `item_checklist_reprovado`, efeito `bloquear_ativo_ate_documento`, simulação e histórico de disparos, efeito na conta conectada só pelo status no teste de RNF-041); **RF-OS-021 completo** (filtro e exportação por campo extra, preenchimento pelo link de uso único) | **Relatório PDF:** RF-SRV-012 (assinatura do cliente), RF-SRV-013 (relatório da OS em PDF, com os campos extras da OS); RF-SRV-003 completo (documentos dos técnicos); **conexão:** RF-SRV-016 (OS do cliente chega à subcontratada) | RF-MOB-010; perfil `prestador` no app, com cache filtrado por fornecedor; assinatura do cliente offline; OS recebida de cliente conectado | S8 (custos, marcos de tempo, assinatura, documentos); RF-OS-007 (atribuição a fornecedor); RF-CON-005 e RF-CON-006 (S8) para as condições e o efeito de documento; RF-CHK-003 (S7) para `item_checklist_reprovado` |
| S10 | 21–22 | API, SSO e cobrança | RF-INT-001, RF-INT-002, RF-INT-003 completo, RF-ADM-007, RF-ADM-009 (com prestador cobrado como técnico e plano Serviços), RF-ADM-013 | **Conexão:** RF-SRV-017 (OS aberta pela subcontratada nasce no cliente), RF-SRV-018 (andamento e liberação pelo cliente), RF-SRV-019 (suspensão e encerramento); RF-SRV-020 (relatório por cliente), RF-SRV-021 (indicadores); RF-SRV-001 com cobrança do plano Serviços; **portão da conexão** (RNF-041 verde, RNF-042 medido) | SSO no app; RF-SRV-011 completo (abrir OS em equipamento de cliente conectado, liberação recusada) | S9 |
| S11 | 23–24 | Endurecimento | Sem RF novo: testes de carga (RNF-007 a RNF-009, RNF-030), teste de 7 dias offline (RNF-010), restauração de backup (RNF-005), acessibilidade (RNF-025), pentest focado em isolamento entre empresas, acesso do prestador, link de uso único e conexão entre contas (proposta), primeiro release candidato `-rc` em homologação (RNF-039), correções do piloto | Carga da sincronização entre contas (RNF-042 sob carga), reconciliação diária, correções do piloto de subcontratada | Correções do piloto, teste em aparelhos de entrada (RNF-027), selo de versão na área segura dos aparelhos | S10 |
| S12 | 25–26 | Lançamento | RF-ADM-001 completo (roteiro de primeiros passos), página de status (RNF-020), documentação pública da API, release `v1.0.0` (seção 3.3), migração dos pilotos para produção | RF-SRV-001 completo (roteiro de primeiros passos da conta Serviços) | Publicação nas lojas | S11; P4, P5 e P7 resolvidas |

Verificação de cobertura: as 13 sprints entregam os 129 requisitos do MVP (105 da versão 0.2, 21 do SRV e 3 da configuração por segmento). RF-ADM-001, RF-ADM-016, RF-ATV-010, RF-ATV-013, RF-INT-003, RF-MOB-003, RF-OS-021, RF-SRV-001, RF-SRV-003 e RF-SRV-011 são entregues em partes e contam como prontos na última sprint em que aparecem. RF-TER-001 passa de S7 para S6, na trilha SRV, sem mudar a contagem. O encaixe da configuração por segmento é o mesmo nas três opções da seção 2.2; na opção B, a sprint extra entra entre S9 e S10.

**Onde entra cada parte da configuração por segmento (proposta):**

| Requisito | Sprint | Partes |
|---|---|---|
| RF-ATV-013 Fase de ciclo de vida e entrada em operação | S2, S4, S6, S8 | S2: fase, status `inativo`, cadastro e importação em implantação, filtro; S4: tipo `implantacao` (RF-OS-002) e `colocar_em_operacao`; S6: planos suspensos e ativados na entrada em operação; S8: documentos críticos exigidos (completo). A exclusão dos KPIs entra com RF-IND-003 e RF-IND-005 em S9 |
| RF-OS-021 Campos extras da OS | S4, S5, S9 | S4: definições, valores, obrigatoriedade na abertura, web, API e app; S5: obrigatoriedade na conclusão, validada no app offline, e impressão; S9: filtro, exportação, link de uso único e relatório da OS de serviço em PDF (completo) |
| RF-ADM-016 Regras configuráveis de condição e efeito | S6, S9 | S6: motor, rotina diária, condições de leitura e de campo da OS, efeitos `notificar`, `abrir_solicitacao` e `bloquear_ativo`; S9: condições de documento e de checklist, `bloquear_ativo_ate_documento`, simulação, histórico e verificação no teste de três contas (completo) |

As partes da Fase 2 do catálogo de RF-ADM-016 (`capacitacao_vencida`, `etapa_projeto_atrasada`, `impedir_atribuicao`) entram com os requisitos que as tornam possíveis (seção 5) e não mudam a contagem.

**Onde entra cada parte do SRV:**

| Parte do SRV | Sprint | Requisitos |
|---|---|---|
| Conta da subcontratada e plano Serviços | S1 (conta), S10 (cobrança), S12 (roteiro) | RF-SRV-001 |
| Técnicos e clientes | S2 (clientes, especialidades), S3 (técnicos e equipamentos de clientes), S9 (documentos dos técnicos) | RF-SRV-002 a RF-SRV-005 |
| OS de serviço e componentes | S5 (catálogo e abertura), S6 (horas, distribuição e app), S8 (componentes) | RF-SRV-006 a RF-SRV-011 |
| Relatório PDF | S9 (assinatura do cliente e relatório da OS), S10 (relatório por cliente e indicadores) | RF-SRV-012, RF-SRV-013, RF-SRV-020, RF-SRV-021 |
| Conexão entre contas | S4 (fundação e teste de três contas), S7 (convite), S8 (ativos), S9 (OS do cliente), S10 (OS da subcontratada, andamento, liberação, encerramento; portão) | RF-SRV-014 a RF-SRV-019, RNF-041, RNF-042 |

**Onde entra cada parte do multipaís e dos idiomas (proposta; a decisão da P15 pode mudar o encaixe, seção 2.3):** E01 entrega a base (país, identificação, moeda, idioma, textos, fuso, endereço) e fecha com o aceite e a revisão por falante nativo; E02 e E04 acrescentam moeda, feriados, formatação e o prefixo `OS-` ou `OT-`; E05 e E09 levam o documento gerado e a notificação ao idioma de quem lê; E07 traz a identificação de qualquer país no fornecedor e o termo da conexão entre países; **E08** traz o conteúdo legal por país (RNF-046); E10 traz a cobrança em `ARS` e `PYG`; E11 e E12 verificam os dois idiomas e os três países e publicam o aditivo por país. O detalhe por etapa está na tabela "Onde entra cada parte do multipaís" da seção 2.3.

### 4.1 Por que esta ordem

- **Versionamento e selo de versão desde o início:** a esteira de RNF-039 fica pronta em S0, junto com o CI, e o selo de versão entra na primeira entrega de interface (S1). Assim, todo build de homologação e de piloto já é identificável.
- **Pessoas junto com usuários (S1):** CNH, capacitações e documentos pessoais ligam-se à `pessoa`; criar a entidade depois exigiria migrar usuários e motoristas.
- **Ativos e medidores antes da OS:** a OS, a preventiva por medidor e os KPIs dependem do ativo, da hierarquia e de leituras confiáveis.
- **Offline desde o início:** a sincronização é o maior risco técnico. O app começa em S0 e cada sprint acrescenta operações offline, em vez de deixar o offline para o fim.
- **Preventiva no meio do MVP:** é o primeiro valor percebido pelo cliente e precisa de OS e medidores prontos.
- **Fornecedores em S7, junto com o estoque:** o cadastro com papéis (RF-TER-001) serve ao estoque e à atribuição de OS; o catálogo (RF-TER-002) é cadastro simples.
- **NR-12 e assinatura em S8:** dependem da OS completa e do parecer P1 sobre assinatura, que tem até a semana 14 para chegar.
- **Prestador em S9:** o acesso do prestador reutiliza a execução de OS (S5), a assinatura (S8) e o modelo de permissões (S1); o filtro por fornecedor só é liberado com o teste de isolamento verde (risco R4).
- **KPIs depois dos dados:** só fazem sentido com marcos de tempo, tipos, custos e falhas já gravados.
- **API em S10:** a API pública expõe o domínio estável. A web e o app usam os mesmos endpoints desde S1, então S10 trata de chaves, escopos, limites, webhooks e documentação.
- **Trilha SRV em paralelo desde S1:** o SRV reaproveita quase tudo do núcleo, mas não pode esperar o núcleo terminar. A trilha anda uma sprint atrás do que reaproveita: equipamentos de clientes em S3, depois do ativo (S2); OS de serviço em S5, depois do núcleo da OS (S4); horas e distribuição em S6, depois do apontamento e da atribuição (S5); componentes em S8, depois da baixa de peças (S7); assinatura do cliente e relatório em S9, depois da assinatura avançada (S8).
- **Fundação da conexão cedo (S4):** a conexão é o maior risco novo (R5, crítico). O teste de três contas entra no CI em S4, vazio, e cresce a cada sprint, como o teste de vazamento entre empresas desde S0. Assim o isolamento é verificado antes de existir fluxo de OS entre contas.
- **Fornecedor em S6, na trilha SRV:** o convite de conexão parte do `fornecedor` do cliente (RF-SRV-014, S7), por isso RF-TER-001 sai de S7 e é feito pela trilha SRV em S6; RF-TER-002 continua em S7 e usa a entidade `servico_fornecedor` criada para o catálogo próprio da subcontratada em S5.
- **Fluxos de OS entre contas em S9 e S10:** dependem da atribuição a fornecedor (RF-OS-007, S5), da conclusão e liberação completas (RF-OS-014, S8) e dos ativos já sincronizados (S8). O terceiro backend, livre do TER em S10, reforça a conexão. Os itens do SRV que estão na linha de corte (RF-SRV-017, RF-SRV-019 parcial, RF-SRV-020) ficam em S10, a última sprint com RF novo, para que o corte não desmonte o que veio antes.
- **Fase de ciclo de vida junto com o ativo (S2):** a fase muda o status e o denominador dos KPIs; acrescentá-la depois exigiria migrar ativos já cadastrados pelos pilotos. `colocar_em_operacao` vem com o tipo `implantacao`, no núcleo da OS (S4), e só fecha em S8, quando existem os documentos críticos que a entrada em operação pode exigir.
- **Campos extras com o núcleo da OS (S4 e S5):** são parte do registro da OS e da validação offline da conclusão; entrar depois obrigaria a refazer o formulário da OS na web e no app. O que depende de outras entregas (link de uso único e relatório em PDF) fecha em S9.
- **Regras depois dos fatos que as disparam (S6 e S9):** o motor entra em S6, quando já existem leituras (S3), campos extras (S4 e S5), notificações e solicitação (S4) e bloqueio de ativo (S2); as condições de documento e de checklist e o bloqueio levantado pelo documento esperam RF-CON-005 e RF-CHK-003 e fecham em S9. Nenhuma regra muda o estado da OS (sumário 3.7), por isso o motor não toca a máquina de estados validada em S4.

---

## 5. Fase 2 — meses 7 a 12

Quatro trilhas em paralelo (proposta): **A** frota e conformidade, **B** estoque (completo e estendido), **C** integrações, identidade e IA, **D** terceiros, locação e SRV. A trilha D ocupa o quarto backend e o terceiro frontend, que vêm da trilha SRV do MVP e acumulam a evolução do SRV pedida pelos pilotos de subcontratadas. A Fase 2 não acrescenta requisitos ao SRV; os ganchos com a Fase 2 são: o contrato de prestação (RF-TER-006, mês 8) passa a definir os ativos cobertos da conexão (RF-SRV-015); as capacitações (RF-CON-012, mês 7) reaproveitam a entidade `capacitacao` já entregue para os técnicos do SRV e o bloqueio de atribuição (RF-CON-013) passa a valer também na conta da subcontratada; o desempenho do fornecedor (RF-TER-007, mês 9) usa os marcos que chegam pela conexão.

A configuração por segmento (sumário 0.5) acrescenta a **trilha E**, projetos, contratos com clientes e pacotes (proposta). A coluna mostra a ordem sugerida com a opção 1 da seção 2.2 (1 backend e 1 frontend a mais); nas opções 2 e 3, a ordem e as dependências são as mesmas, e o trabalho é dividido entre as trilhas A e D, com mais prazo ou com requisitos movidos para a Fase 3, por decisão do dono do produto. Dependências da trilha E:

- **PRJ** depende de RF-ATV-013 e RF-OS-021, entregues no MVP; a OS de etapa em subcontratada conectada usa RF-SRV-016 (MVP); a medição de etapas (RF-PRJ-006) usa o valor do contrato de prestação (RF-TER-006, mês 8) quando houver.
- **CTR** depende dos registros de parada da OS (RF-OS-005) e da `causa_externa` gravada desde o MVP (RF-OS-001), e do cliente atendido (`cliente_servico`, RF-SRV-004); a inclusão automática dos ativos de projeto concluído (RF-CTR-002) anda junto com RF-PRJ-007.
- **Pacotes de segmento** (RF-ADM-017) dependem das regras (RF-ADM-016) e dos campos extras (RF-OS-021) do MVP e dos modelos de projeto (RF-PRJ-001); o pacote `fiscalizacao_eletronica` depende de **P9**: sem ela, sai com a validade da verificação metrológica digitada e as regras metrológicas desligadas (sumário 6).
- As condições e o efeito da Fase 2 de RF-ADM-016 entram com o que os torna possíveis: `capacitacao_vencida` e `impedir_atribuicao` com RF-CON-012 e RF-CON-013 (mês 7, trilha C) e `etapa_projeto_atrasada` com RF-PRJ-003 (mês 8, trilha E).

| Mês | Trilha A — frota e conformidade | Trilha B — estoque | Trilha C — integrações, identidade e IA | Trilha D — terceiros, locação e SRV | Trilha E — projetos, contratos com clientes e pacotes (proposta) | Depende de |
|---|---|---|---|---|---|---|
| 7 | RF-PRE-007 (planos aninhados), RF-OS-012 (catálogo de falhas) | RF-EST-005 (inventário), RF-EST-007 (requisição), RF-EST-008 (compras), RF-EST-012 (vários almoxarifados e transferência) | RF-ADM-003 (multi-unidade), RF-IND-009, RF-CON-012, RF-CON-013 (capacitações e bloqueio; com `capacitacao_vencida` e `impedir_atribuicao` de RF-ADM-016) | RF-TER-009 (orçamento de fornecedor), RF-TER-008 (oficina externa e assistência técnica) | RF-PRJ-001 (modelos de projeto), RF-PRJ-002 (projeto de implantação) | MVP; feedback do lançamento; RF-ATV-013 e RF-OS-021 (MVP) para PRJ |
| 8 | RF-FRO-002 a RF-FRO-006 (pneus, CPK), RF-MOB-009 | RF-EST-009 (NF-e e NFS-e), RF-EST-011, RF-MOB-008 | RF-IA-001; infraestrutura de busca vetorial isolada por empresa | RF-TER-006 (contratos de prestação com SLA, que passam a definir os ativos cobertos da conexão, RF-SRV-015), RF-OS-020 (garantia) | RF-PRJ-003 (etapas e dependências, com `etapa_projeto_atrasada`), RF-PRJ-004 (OS de implantação da etapa, também no app offline e na subcontratada conectada) | RF-EST-008 para NF-e; contador para CBS e IBS; RF-SRV-016 (MVP) para RF-PRJ-004 |
| 9 | RF-FRO-007, RF-FRO-008, RF-FRO-009 (combustível), RF-INT-007, RF-FRO-010 (multas), RF-FRO-011 (sinistros), RF-IND-006 | RF-EST-013, RF-EST-014 (almoxarifado em terceiro, saldo geral e por subcontratada), RF-EST-015 (remessa e retorno), RF-EST-019 (NF-e de remessa e retorno) | RF-IA-002, RF-IA-003 | RF-TER-007 (avaliação), RF-LOC-001 (locatários), RF-LOC-002 (contrato de locação) | RF-PRJ-005 (acompanhamento), RF-PRJ-006 (medição de etapas para pagamento do prestador), RF-CTR-001 (contrato com cliente) | **P8** para RF-EST-015 e RF-EST-019; RF-TER-004 (MVP); RF-TER-006 (mês 8) para RF-PRJ-006 |
| 10 | RF-CON-008, RF-CON-009, RF-CON-010 (NR-13) | RF-EST-016 (peças serializadas), RF-EST-017 (reparáveis e quarentena) | RF-IA-004, RF-INT-004 (modelo canônico de telemática), RF-ADM-008 (SAML e SCIM) | RF-LOC-003 (itens), RF-LOC-004 (matriz de responsabilidade), RF-LOC-005 (vistorias, também no app offline), RF-LOC-006 (avaria gera OS) | RF-PRJ-007 (entrega e entrada em operação), RF-CTR-002 (itens cobertos), RF-CTR-003 (paradas por causa externa) | **P6** para RF-CON-009; consultor de SST; RF-CHK-001 para vistorias; RF-OS-005 (MVP) para CTR |
| 11 | RF-CON-011 (LOTO), RF-CON-007 (aptidão e bloqueio de despacho) | RF-EST-018 (ordem de reparo) | RF-INT-005 (Cobli, Geotab, rFMS), RF-INT-006 (DTC), RF-INT-008 (ERP: Sankhya, Protheus, Omie), RF-IND-008 | RF-LOC-007 (ativo tomado e SLA de substituição), RF-LOC-008 (portal do locatário), RF-LOC-009 (relatório de cobrança) | RF-CTR-004 (medição mensal), RF-CTR-005 (painel do contrato), RF-ADM-017 (pacotes de segmento) | RF-INT-004; RF-TER-009 e RF-EST-017 para a ordem de reparo; acesso a ambientes de teste dos ERPs e telemáticas; RF-PRJ-001 para RF-ADM-017; **P9** para o pacote `fiscalizacao_eletronica`; **P10** para o acesso do cliente ao painel |
| 12 | Endurecimento e piloto de frota e NR-13 | Endurecimento e piloto de estoque em terceiro | Pentest externo completo; início da janela de observação SOC 2 | Endurecimento, teste de isolamento do locatário (RNF-038) e piloto de locação | Endurecimento, primeira medição mensal real fechada e piloto de configuração por segmento | Todos os anteriores |

Verificação de cobertura: os seis blocos entregam os 75 requisitos da Fase 2 (trilha A: 20; B: 14; C: 14; D: 14; E: 13).

**Por que esta ordem na trilha E:** PRJ começa no mês 7 porque só depende do MVP e é o que o caso de origem (fiscalização eletrônica) usa primeiro, na instalação. CTR começa no mês 9 e fecha a medição no mês 11, para que a primeira medição real (de um mês inteiro) seja fechada no mês 12, ainda dentro do piloto. Os pacotes vêm por último porque empacotam o que já existe, inclusive os modelos de projeto, e porque o conteúdo de `fiscalizacao_eletronica` depende de P9.

### 5.1 Entregas e critério de pronto por mês

| Mês | Entrega verificável |
|---|---|
| 7 | OS enviada a oficina externa com orçamento aprovado por alçada e nota lançada; empresa com duas unidades e usuário restrito a uma delas, testado também pela API; transferência entre dois almoxarifados com recebimento confirmado; técnico sem capacitação válida bloqueado na atribuição; modelo de projeto com etapas e dependências sem ciclo e projeto criado a partir dele, com as datas previstas calculadas |
| 8 | Veículo com mapa de pneus, rodízio no app offline e CPK calculado por pneu; NF-e real importada com itens casados e custo médio atualizado; contrato de prestação com SLA estourado gerando `fornecedor.sla_violado`; contrato de prestação com fornecedor conectado definindo os ativos sincronizados na conta da subcontratada; etapa liberada só com as predecessoras concluídas ou dispensadas; OS `implantacao` gerada pela etapa, executada offline por subcontratada conectada e fora dos KPIs de manutenção |
| 9 | Abastecimentos importados de cartão com alertas de anomalia; multa com prazo de 30 dias e condutor identificado pelo vínculo; remessa a subcontratada com NF-e registrada, consumo pelo `prestador` em OS e retorno das sobras, com saldo geral e por subcontratada conferidos; contrato de locação cedido ativado; cronograma previsto e real do projeto com avanço físico; medição de etapas por fornecedor; contrato com cliente ativado com meta, calendário e faixas de desconto |
| 10 | Prazos NR-13 calculados para caldeira e vaso, com a tabela conferida no DOU; peça serializada instalada e removida pela OS, indo para a quarentena; vistoria de devolução offline com avaria gerando OS com custo do locatário; projeto concluído colocando os ativos em operação e incluindo-os como itens do contrato com cliente; parada com causa externa e evidência registrada pelo app |
| 11 | OS com LOTO que não inicia sem bloqueios verificados; aptidão do veículo consultável pela API; ordem de reparo completa (envio, orçamento, aprovação, retorno como recondicionada com garantia); leituras e DTCs de pelo menos um conector de telemática chegando a uma empresa piloto; sincronização com pelo menos um ERP; login SAML e SCIM testados; locatário abrindo solicitação pelo portal e relatório de cobrança de avarias e horas excedentes gerado; medição mensal do contrato com memória de cálculo, parada com causa externa aprovada descontada só na medição e disponibilidade técnica inalterada; pacote de segmento aplicado a uma conta nova e reaplicado em versão nova sem sobrescrever o que a conta editou; `fiscalizacao_eletronica` com as regras metrológicas inativas enquanto **P9** não for resolvida |
| 12 | Pentest externo sem achado crítico ou alto aberto, incluindo a conexão entre contas; testes de isolamento do prestador, do locatário e das contas conectadas (RNF-041) verdes; controles SOC 2 documentados e em operação; métricas dos pilotos de frota, NR-13, estoque em terceiro, locação e subcontratadas com conta própria revisadas; primeira medição mensal real de contrato com cliente fechada no piloto de configuração por segmento |

---

## 6. Fase 3 — meses 13 a 24

| Mês | Entrega | Requisitos | Depende de |
|---|---|---|---|
| 13 | Programação por capacidade, primeira parte (calendário de técnicos, carga semanal) | RF-OS-019 (início) | RF-OS-007 |
| 14 | Programação congelada, schedule compliance e wrench time | RF-OS-019, RF-IND-011 | Mês 13 |
| 15 | Ingestão de IoT por MQTT e HTTP, séries temporais, gateway de parceiro para OPC UA e Sparkplug B | RF-INT-010 | Parceiro de gateway definido |
| 16 | Manutenção baseada em condição; confiabilidade com Weibull | RF-PRE-011, RF-IND-010 | RF-INT-010; histórico de falhas da Fase 2 |
| 17 | Agentes de IA com aprovação; assistente de análise de falha | RF-IA-006, RF-IA-005 | RF-IA-001 a RF-IA-004 |
| 18 | Relatórios em linguagem natural; servidor MCP somente leitura | RF-IA-007, RF-INT-011 | RF-INT-001 |
| 19 | TCO e renovação; conector de BI; **relatório SOC 2 Type II** | RF-IND-012, RF-INT-012 | Custos de frota da Fase 2; janela SOC 2 de 6 meses |
| 20 | RCM e FMEA com Action Priority | RF-PRE-012 | RF-OS-012; leitura de SAE JA1011 e AIAG-VDA |
| 21 | Gestão de ativos ISO 55001 | RF-ATV-012 | RF-PRE-012 |
| 22 | Assinatura ICP-Brasil; importação automática de multas | RF-CON-015, RF-INT-009 | **P1** e **P3** |
| 23 | Banco dedicado opcional | RF-ADM-014 | Primeiro cliente Enterprise que exija |
| 24 | Consolidação, pentest anual, preparação ISO 27001 | Sem RF novo | — |

Verificação de cobertura: os doze blocos entregam os 16 requisitos da Fase 3. Na Fase 3, o quarto backend e o terceiro frontend permanecem e absorvem a evolução de TER, SRV, EST e LOC pedida pelos clientes (proposta).

Contingências:
- Se o **portão da conexão** (RNF-041 verde e RNF-042 medido) não passar em M4, RF-SRV-014 a RF-SRV-019 saem do lançamento e entram no mês 7, na trilha D, com aval do dono do produto (o sumário 4 põe a conexão no MVP). O restante da trilha D desliza um mês e a válvula é a do risco 19: RF-LOC-007 e RF-LOC-009 passam para o início da Fase 3.
- Se o parecer **P1** exigir ICP-Brasil para documentos de SST, o MVP sai com a assinatura avançada (RF-CON-014) e aviso ao cliente sobre a pendência (proposta); RF-CON-015 sobe para o mês 9 da Fase 2, na trilha C, e RF-IA-003 desce para o mês 12.
- Se **P3** confirmar que não há API acessível, RF-INT-009 sai do roadmap e as multas continuam por digitação e planilha.
- Se **P8** não estiver resolvida até o mês 9, RF-EST-015 e RF-EST-019 entram com CFOP informativo e nota obrigatória por parâmetro, como prevê RF-EST-019, e a validação fiscal vira correção no mês 12.
- Se **P9** não estiver resolvida até o mês 11, RF-ADM-017 sai com o pacote `fiscalizacao_eletronica` trazendo o tipo de documento "Verificação metrológica" com validade digitada e as regras metrológicas inativas, como prevê o requisito; as regras são ligadas por versão nova do pacote quando P9 fechar, sem sobrescrever o que a conta editou.
- Se **P10** não estiver decidida até o mês 11, RF-CTR-005 sai como painel interno e a medição vai ao cliente final em PDF e XLSX (RF-CTR-004); acesso próprio do cliente final é escopo novo, com requisito, isolamento e preço próprios.
- Se o dono do produto escolher a opção 2 da seção 2.2 para a Fase 2, a trilha E é dividida entre as trilhas A e D, M6 passa para o mês 13 ou 14 e os blocos da Fase 3 deslizam o mesmo tanto, exceto a janela do SOC 2.
- Se o parecer da **P14** ou o mínimo da **P12** de um país não estiver pronto até o M5, o país **não é liberado à venda** (portão por país, opção D da seção 2.3, com aval do dono do produto): o código e o espanhol continuam no lançamento, e o cliente do país entra quando o portão fechar. Se a **P15** não estiver decidida até a semana 10, a página de preços (S10) sai só para o Brasil e a cobrança em `ARS` e `PYG` segue por fatura manual, sem mudar o código de moeda (RNF-045).
- Se a **P13** não estiver resolvida até o mês 8, o fiscal da Argentina e do Paraguai (RF-EST-009, RF-EST-015, RF-EST-019) entra com a referência do documento fiscal digitada, sem validação do formato do país, como o ADR-013 prevê, e a validação vira correção posterior.
- CDC para ClickHouse só entra se o volume de leituras e telemetria justificar, como recomenda o estudo; não há bloco reservado.

---

## 7. Atividades não técnicas

| Atividade | Início | Prazo | Responsável | Resultado esperado | Requisitos afetados |
|---|---|---|---|---|---|
| **P1** Parecer jurídico sobre o item 1.6 da NR-1 (ICP-Brasil em documentos de SST) | S0 | Semana 14 (fim de S6) | Gerente de produto e advogado | Forma de assinatura aceita para OS de NR-12 e registros de NR-13 | RF-OS-014, RF-CON-002, RF-CON-010, RF-CON-014, RF-CON-015 |
| **P2** Leitura dos anexos das Portarias Inmetro 127 e 128/2022 (CIV e CIPP) | S0 | Semana 6 | Gerente de produto e consultor de SST | Confirmar ou descartar o cálculo sugerido de validade | RF-CON-005 |
| **P3** Consulta ao Serpro sobre API Senatran (e-Frotas) | S0 | Mês 6 | Gerente de produto | Decidir RF-INT-009 | RF-FRO-010, RF-INT-009 |
| **P4** 10 a 15 entrevistas com gestores de manutenção e frota, mais 5 a 8 com subcontratadas e assistências técnicas (proposta) | S1 | Semana 10; síntese na semana 12 | Gerente de produto e designer | Preços e mínimos validados antes da página de preços (S10), inclusive o plano Serviços e o desconto de canal; prioridades da Fase 2; interesse em TER, SRV e LOC; recrutamento de pilotos | RF-ADM-009, RF-SRV-001, planos marcados "(proposta)"; RF-ADM-016 e RF-OS-021 no plano Serviços (ponto 37 da seção 8 de 1_ESPECIFICACAO_REQUISITOS.md) |
| **P5** Busca no INPI (classes 9 e 42), depósito de marca e registro de domínio | S0, semana 1 | Busca na semana 1; depósito até a semana 4 | Gerente de produto e advogado | Nome liberado ou troca decidida antes de qualquer material do piloto (S8) | Todo material público |
| **P6** Conferência no DOU dos prazos NR-13 com SPIE, categorias III a V | Mês 7 | Mês 9 | Consultor de SST | Tabela de prazos confirmada | RF-CON-009 |
| **P7** Decisão sobre o local da cópia de recuperação de desastre | S0 | Semana 20 (fim de S9) | Dono do produto, com advogado e Plataforma/SRE | Uma das opções da seção 5.1 do documento 1 escolhida; contrato de tratamento de dados e RTO de região ajustados antes do lançamento | RNF-004, RNF-005, RNF-006 |
| **P8** Validação com o contador dos CFOPs e do tratamento fiscal de remessa e retorno a subcontratada, conserto e locação | Mês 7 | Mês 8 | Contador e gerente de produto | Lista de CFOPs e obrigatoriedade da NF-e por tipo de movimentação | RF-EST-015, RF-EST-018, RF-EST-019, RF-LOC-002, RF-LOC-009 |
| Piloto fechado do MVP | Recrutamento em S2–S5 | Uso em S8–S11 | Gerente de produto | 3 a 5 empresas, com pelo menos uma indústria com NR-12, uma transportadora e uma operação logística ou de varejo, e pelo menos uma que use subcontratada ou assistência técnica (proposta) | Todo o MVP |
| Piloto de subcontratada com conta própria (SRV) | Recrutamento em S2–S5 | Uso em S8–S11 | Gerente de produto | 1 a 2 subcontratadas ou assistências técnicas com conta própria; pelo menos uma conectada a uma das empresas piloto, para exercitar a conexão entre contas em S10 e S11 (proposta) | RF-SRV-001 a RF-SRV-021, RNF-041, RNF-042 |
| Termo da conexão entre contas (LGPD) | S4 | S7 (antes do primeiro convite) | Advogado e gerente de produto | Papéis de cada conta no compartilhamento (controlador e operador), base legal, aceite no convite e na aceitação, cláusula no DPA | RF-SRV-014, RNF-021 |
| **P9** Leitura da portaria do Inmetro e da resolução do Contran vigentes sobre medidores de velocidade (periodicidade da verificação, reparo que exige nova verificação, controle de lacres) | Mês 7 | Mês 10 (antes de RF-ADM-017, mês 11) | Gerente de produto, com especialista em metrologia legal (proposta) | Validade sugerida da verificação metrológica e regras do pacote confirmadas; sem resposta, contingência da seção 6 | RF-ADM-016, RF-ADM-017, RF-CON-005 |
| **P10** Decisão sobre o acesso próprio do cliente final de um contrato (ex.: prefeitura) para abrir chamados e ver a medição | Mês 7 | Mês 9 (antes do desenho de RF-CTR-005) | Dono do produto | Manter painel interno com medição em PDF e XLSX, ou abrir escopo novo com requisito, isolamento e preço | RF-CTR-004, RF-CTR-005 |
| **P12** Pesquisa confirmada por consultor de SST de cada país: equivalentes da NR-12 e da NR-13, inspeção de vasos de pressão por jurisdição, documentos de veículo e de condutor, multas | S0 (pesquisa feita em 13/11/2026, a validar) | Mínimo antes de E08 (semana 14) para o desenho dos pacotes; parecer completo antes de liberar cada país | Gerente de produto e consultor de SST da Argentina e do Paraguai | Conteúdo dos pacotes de país confirmado, com a fonte; sem parecer, os pacotes nascem sem prazo legal (RNF-046) | RNF-046, RF-CON-001 a RF-CON-013, RF-FRO-001 a RF-FRO-011 |
| **P13** Contador de cada país: identificação fiscal (CUIT, CUIL, DNI, RUC e CI) e fatura eletrônica no lugar da NF-e | Mês 7 (a regra do dígito do CUIT e do RUC antes, se o cadastro for validar) | Mês 8 | Contador da Argentina e do Paraguai e gerente de produto | Regra oficial do dígito verificador; tipo e referência do documento fiscal do país | RNF-044, RF-EST-009, RF-EST-015, RF-EST-019, RF-LOC-002, RF-LOC-009 |
| **P14** Parecer jurídico sobre a proteção de dados da Argentina e do Paraguai e a transferência internacional para a nuvem em São Paulo | S4 | S11 (antes do lançamento) | Advogado local de cada país e gerente de produto | Aditivo de transferência por país, com a Tenaz como encarregada; termo da conexão entre contas de países diferentes; kit LGPD por país | RNF-006, RNF-021, RF-SRV-014 |
| **P15** Decisão do dono do produto sobre preço, moeda de cobrança e meio de pagamento na Argentina e no Paraguai, e sobre as opções de prazo e time da seção 2.3 | S2 (junto com a análise do consumo real da E01) | Fim de S3 para time e prazo; semana 10 para preço (junto com P4) | Dono do produto, com o gerente de produto | Opção A a D da seção 2.3 escolhida; tabela de preços por país; moeda de cobrança | RF-ADM-009, RF-SRV-001, RNF-045 |
| Revisão do espanhol por falante nativo (um argentino e um paraguaio) | E01, depois a cada release | Antes do piloto (S8) e do lançamento (S12) | Gerente de produto e revisor externo | Termos marcados no glossário bilíngue revisados; `es-419` neutro aceito nos dois países | RNF-043 |
| Piloto por país (Argentina e Paraguai), se o dono do produto quiser | Recrutamento em S5 | Uso em S9 a S11 | Gerente de produto | Pelo menos uma empresa em cada país, para exercitar a identificação, a moeda e o espanhol com dado real (proposta; depende do portão por país) | RNF-043 a RNF-045 |
| Conteúdo dos pacotes de segmento | Mês 9 | Mês 11 | Gerente de produto, com o domínio de manutenção e o consultor de SST nas partes de conformidade | Seis pacotes iniciais (proposta) revisados, com regras inativas por padrão | RF-ADM-017 |
| Piloto de configuração por segmento (Fase 2) | Recrutamento no mês 8 | Mês 12 | Gerente de produto | Empresa com implantação em campo e contrato com meta de disponibilidade (por exemplo, o fornecedor de radares do caso de origem) usando projetos, contrato com cliente e o pacote do segmento (proposta) | RF-PRJ-001 a RF-PRJ-007, RF-CTR-001 a RF-CTR-005, RF-ADM-017 |
| Piloto de frota (Fase 2) | Mês 8 | Mês 12 | Gerente de produto | Transportadora usando pneus, combustível e multas | RF-FRO-002 a RF-FRO-011 |
| Piloto NR-13 (Fase 2) | Mês 10 | Mês 12 | Gerente de produto e consultor de SST | Indústria com caldeira ou vasos usando prazos e Registro de Segurança | RF-CON-008 a RF-CON-010 |
| Piloto de estoque em terceiro (Fase 2) | Mês 9 | Mês 12 | Gerente de produto | Empresa com subcontratada mantendo estoque e com peças reparáveis enviadas para conserto | RF-EST-013 a RF-EST-019 |
| Piloto de locação (Fase 2) | Mês 9 | Mês 12 | Gerente de produto | Locadora de equipamentos ou empresa com ativos tomados em locação usando contratos, vistorias e portal do locatário | RF-LOC-001 a RF-LOC-009 |
| Kit LGPD (DPA, RIPD modelo, política de privacidade, termos de uso, encarregado) | S8 | S11 | Advogado e gerente de produto | Publicado no lançamento | RNF-021 |
| Pentest | Focado em isolamento, autenticação, acesso do prestador, link de uso único e conexão entre contas em S11 (proposta); completo no mês 12, incluindo portal do locatário; depois anual | — | Plataforma/SRE e empresa contratada | Sumário de pentest disponível a clientes | RNF-018, RNF-037, RNF-038, RNF-041 |
| SOC 2 | Avaliação de lacunas no mês 7; controles e ferramenta de conformidade nos meses 8 a 11 | Janela de observação dos meses 12 a 18; relatório Type II no mês 19 (proposta) | Plataforma/SRE e auditor | Controles e janela iniciados na Fase 2 e relatório Type II na Fase 3, conforme o sumário | RNF-036 |
| ISO 27001 | Mês 21 | Depois do mês 24 | Plataforma/SRE | Preparação iniciada | RNF-036 |
| Validação fiscal da importação de NF-e (CBS e IBS) | Mês 7 | Mês 8 | Contador | Leiaute aceito pelo leitor de XML | RF-EST-009 |
| Página de preços, documentação da API, base de conhecimento | S10 | S12 | Gerente de produto e designer | Publicadas no lançamento | RF-ADM-009, RF-INT-001 |

---

## 8. Dependências críticas

1. **Medidores → preventiva por medidor → KPIs de frota.** Sem RF-MED-003 confiável, a preventiva por km gera OS errada e o custo por km fica sem sentido.
2. **OS completa → NR-12.** RF-CON-002 é uma validação na conclusão da OS (RF-OS-014); não pode vir antes.
3. **P1 → assinatura de OS de NR-12 e de registros de NR-13.** A assinatura avançada (RF-CON-014) entra em S8 com arquitetura plugável, para trocar o método sem migrar dados.
4. **Compras → NF-e → custo de oficina externa.** RF-EST-009 depende de RF-EST-008; RF-TER-008 usa RF-EST-009 para a nota e RF-TER-009 para o orçamento.
5. **Capacitações → bloqueio de atribuição → LOTO.** RF-CON-013 precisa de RF-CON-012; o LOTO usa as capacitações para saber quem pode bloquear.
6. **Modelo canônico de telemática → conectores → DTC.** RF-INT-005 e RF-INT-006 dependem de RF-INT-004.
7. **IA nível 1 → agentes.** RF-IA-006 depende da governança (RF-IA-001) e da busca (RF-IA-002).
8. **SOC 2 Type II → janela de observação.** O relatório não sai antes de alguns meses de controles operando; não é possível comprimir no fim.
9. **Pessoa → usuários, motoristas, capacitações.** RF-ADM-015 entra em S1, antes de qualquer cadastro que dependa da pessoa.
10. **Perfil prestador e isolamento → estoque em terceiro.** RF-EST-013 depende de RF-TER-004 e do teste de RNF-037 já em operação desde o MVP.
11. **Peças serializadas → reparáveis e quarentena → ordem de reparo.** RF-EST-018 depende de RF-EST-016, RF-EST-017 e do orçamento de fornecedor (RF-TER-009).
12. **P8 → remessa, retorno e conserto com NF-e.** RF-EST-015, RF-EST-018 e RF-EST-019 dependem da lista de CFOPs; a contingência da seção 6 evita bloqueio.
13. **Contrato de locação → itens e matriz → vistorias → avaria → cobrança.** RF-LOC-009 só faz sentido com RF-LOC-004 a RF-LOC-006 prontos; a vistoria reutiliza o motor de checklists (RF-CHK-001).
14. **Esteira de versionamento → todo release.** RNF-039 precisa estar ativo desde S0; um release sem tag, CHANGELOG e versão coerentes é bloqueado no CI.
15. **Ativo e OS → SRV.** O equipamento de cliente é um `ativo` (RF-ATV-002, S2) e a OS de serviço é uma `ordem_servico` (RF-OS-001 a RF-OS-009, S4 e S5); a trilha SRV anda uma sprint atrás do que reaproveita.
16. **Fundação da conexão → fornecedor e convite → ativos → OS entre contas → liberação.** O teste de três contas (S4) vem antes de qualquer fluxo; o convite (RF-SRV-014, S7) depende de RF-TER-001 (S6); a OS do cliente (RF-SRV-016, S9) depende dos ativos sincronizados (RF-SRV-015, S8) e da atribuição a fornecedor (RF-OS-007); a liberação de volta (RF-SRV-018, S10) depende de RF-OS-014 (S8).
17. **Assinatura avançada → assinatura do cliente e relatório.** RF-SRV-012 e RF-SRV-013 (S9) usam a assinatura e o gerador de PDF de RF-CON-014 e RF-OS-018.
18. **Termo da conexão (LGPD) → primeiro convite.** Sem o termo aprovado até S7, o convite fica restrito a homologação.
19. **Contrato de prestação (Fase 2) → ativos cobertos.** No MVP, a conexão usa a seleção manual de ativos; a partir do mês 8, o contrato vigente define a lista (RF-SRV-015).
20. **Ativo e núcleo da OS → fase de ciclo de vida e campos extras → regras.** RF-ATV-013 começa com o ativo (S2) e RF-OS-021 com o núcleo da OS (S4); o motor de RF-ADM-016 (S6) usa os campos extras `sim_nao` como condição.
21. **Documentos com vencimento → bloqueio até documento e entrada em operação.** `bloquear_ativo_ate_documento` (RF-ADM-016, S9) e a exigência de documento crítico em `colocar_em_operacao` (RF-ATV-013, S8) dependem de RF-CON-005 (S8).
22. **Fase de ciclo de vida e campos extras (MVP) → PRJ (Fase 2).** RF-PRJ-002 cria ativos `em_implantacao`, RF-PRJ-004 gera OS `implantacao` com campos extras obrigatórios e RF-PRJ-007 executa `colocar_em_operacao`.
23. **Paradas da OS → medição do contrato.** RF-CTR-004 calcula a disponibilidade contratual com os registros de parada (RF-OS-005) e a `causa_externa` (RF-OS-001, RF-CTR-003); sem marcos de parada confiáveis desde o MVP, a primeira medição não fecha.
24. **Regras, campos extras e modelos de projeto → pacotes; P9 → `fiscalizacao_eletronica`.** RF-ADM-017 só empacota o que já existe; as regras metrológicas do pacote esperam P9 (contingência da seção 6).
25. **Lista fechada da conexão → todo campo novo.** Os campos da configuração por segmento são classificados antes de fechar o requisito que os cria (critério de pronto 3.1, item 2); `causa_externa` ainda depende de classificação no sumário (seção 2.2).
26. **Perfil de país e moeda → toda tela, tabela e documento novos.** RNF-043 a RNF-045 valem desde a E01: toda etapa seguinte nasce com texto nos dois idiomas, `moeda` ao lado de todo valor e o teste de vazamento com os três países; deixar para depois obrigaria a refazer telas e migrar colunas.
27. **P12, P13 e P14 → liberar um país à venda.** O conteúdo legal (RNF-046), o fiscal e o aditivo de transferência de dados dependem de pesquisa e parecer externos; o código multipaís não depende deles, a **venda** no país sim (portão por país, seção 2.3).

---

## 9. Riscos do cronograma

| # | Risco | Impacto | Probabilidade | Mitigação |
|---|---|---|---|---|
| 1 | Contratação atrasada ou time menor que o da seção 2, em especial o terceiro backend (S0), o quarto backend e o terceiro frontend antecipados para S1 por causa do SRV | Alto | Alta | Iniciar com o núcleo; aplicar a linha de corte da seção 2; mês 0 só conta com o núcleo contratado; sem os reforços do SRV, as alternativas da seção 2.1 (15 a 16 sprints, 14 sprints com só o quarto backend, ou conexão no mês 7 com aval do dono do produto); as alternativas da seção 2.1 não contam a configuração por segmento, que soma cerca de 1 a 1,5 sprint a cada uma delas se não houver reforço (seção 2.2) |
| 2 | Escopo do MVP (129 requisitos, com a configuração por segmento) grande para 6 meses | Alto | Alta | Quarto backend e terceiro frontend desde S1; trilha SRV paralela; linha de corte pré-definida, com os itens do SRV concentrados em S10; opções A, B e C da seção 2.2 levadas ao dono do produto; revisão de escopo no fim de S6 com os dados de velocidade das primeiras sprints, que confirma a opção escolhida |
| 3 | Sincronização offline mais difícil que o previsto (conflitos, numeração provisória, fotos, cache filtrado do prestador e do técnico da subcontratada, comando offline sobre OS ligada a cliente conectado) | Alto | Média | Dois desenvolvedores mobile desde S0; prova de conceito em S0; teste de 72 h em S6 e de 7 dias em S11; escopo de cache fechado em RF-MOB-002; o comando offline é aplicado primeiro na conta da subcontratada e só depois sincronizado com o cliente (RF-SRV-011) |
| 4 | Parecer P1 atrasado ou exigindo ICP-Brasil | Alto | Média | Contratar em S0; assinatura avançada plugável; contingência da seção 6 |
| 5 | Desempenho de RLS com volume real, somado aos filtros por fornecedor e por locatário e ao processo de sincronização entre contas | Médio | Média | Índices por empresa e por fornecedor desde S0; testes de carga a partir de S6 com 2.000 e 10.000 ativos; carga da sincronização entre contas em S11 com meta de RNF-042 |
| 6 | Pilotos indisponíveis ou pouco engajados | Médio | Média | Recrutar 5 empresas para manter 3; gratuidade durante o piloto (proposta); contato semanal |
| 7 | Nome Tenaz indisponível no INPI (P5) | Médio | Baixa | Busca na semana 1; decisão de troca antes de S8 |
| 8 | Preço não validado (P4) ou upgrade fraco para o Profissional no lançamento (R3 do sumário), já que pneus, NR-13, multi-unidade, estoque em terceiro e locação só chegam na Fase 2 | Médio | Média | Entrevistas até a semana 10; preço de lançamento do Profissional com desconto até a Fase 2, ou lançar comercialmente o Profissional só na Fase 2, conforme o sumário |
| 9 | APIs de telemática e cartão combustível não verificadas (Ticket Log, Sascar, Omnilink) | Médio | Alta | Começar por Cobli, Geotab e rFMS; importação de arquivo como alternativa |
| 10 | Conectores de ERP exigem ambientes de teste e parcerias | Médio | Média | Iniciar contato com Sankhya, TOTVS e Omie no mês 8 |
| 11 | Mudança regulatória (tabela NR-13, leiaute de NF-e com CBS e IBS) | Médio | Média | Tabelas parametrizáveis e versionadas; contador e consultor de SST no ciclo |
| 12 | Custo de modelos de IA acima do previsto | Médio | Média | Créditos por plano; limites por empresa; IA só na Fase 2 |
| 13 | "Price match" da Tractian anula vantagem de preço (R1 do sumário) | Médio | Média | Priorizar diferenciais (frota unificada, offline, conformidade, terceiros e locação) nos pilotos e no material de venda |
| 14 | Dependência de poucas pessoas (líder técnico, mobile, backend da sincronização entre contas) | Médio | Média | Revisão em par; decisões registradas em 5_ARQUITETURA_SOFTWARE.md (inclusive o ADR-009); dois mobile desde o início; o terceiro backend acompanha a conexão desde S7 e assume parte dela em S10 |
| 15 | Aprovação nas lojas de apps e atualizações do Expo | Baixo | Média | Teste fechado nas lojas desde S4; folga de uma semana antes do lançamento |
| 16 | **R4 do sumário:** prestador com acesso ao sistema vê dados de outra empresa contratante ou de outro prestador; o mesmo vale para o locatário | Crítico | Média | Filtro por fornecedor e por locatário no servidor, além do RLS; testes automatizados de RNF-037 (MVP, S9) e RNF-038 (Fase 2, mês 12) a cada build; perfil `prestador` só liberado com o teste verde; pentest focado em S11 e completo no mês 12 |
| 17 | **P7** sem decisão até o lançamento | Alto | Média | Decisão até a semana 20; sem ela, o lançamento publica só o RTO de zona e mantém a cópia de recuperação dentro do que o contrato com o cliente permitir |
| 18 | **P8** sem resposta do contador antes do estoque em terceiro e da locação | Médio | Média | Contador contratado no mês 7; CFOP informativo e nota obrigatória por parâmetro como contingência (seção 6) |
| 19 | Escopo da Fase 2 (75 requisitos, com PRJ, CTR e pacotes) grande para 6 meses, com mobile sobrecarregado por vistorias, peças serializadas e prestador movimentando estoque | Alto | Alta | Quarta trilha com backend e frontend adicionais; trilha E conforme a opção escolhida pelo dono do produto (seção 2.2); vistoria reaproveita o motor de checklists e a OS de implantação reaproveita a OS; revisão de escopo no fim do mês 9; se preciso, e com aval do dono do produto, RF-LOC-007 e RF-LOC-009 passam para o início da Fase 3 |
| 20 | Esteira de versionamento bloqueando releases urgentes de correção | Baixo | Baixa | Release de PATCH com o mesmo fluxo automatizado (tag, CHANGELOG, build); nenhum atalho manual de versão (RNF-039) |
| 21 | **R5 do sumário:** a sincronização entre contas conectadas copia dado além do permitido ou para a conta errada | Crítico | Média | Fundação e teste de três contas no CI desde S4 (RNF-041); lista fechada de campos validada no envio e no recebimento; transações separadas por conta; papel de banco próprio sem `BYPASSRLS`; auditoria nos dois lados; conexão liberada só com o portão de M4; pentest focado em S11 e completo no mês 12; classificação de todo campo novo como compartilhável ou não (critério de pronto 3.1) |
| 22 | **R6 do sumário:** SRV no MVP aumenta o escopo do lançamento | Alto | Alta | Reaproveitar OS, apontamento, ativos, app e catálogo (seção 2.1); quarto backend e terceiro frontend antecipados para S1; trilha SRV paralela; itens do SRV na linha de corte concentrados em S10; portão da conexão com saída para o mês 7 |
| 23 | Mobile e QA sobrecarregados pelo SRV (componentes, assinatura do cliente, OS ligada a cliente, cenários de sincronização entre contas) e pela configuração por segmento (campos extras dinâmicos validados offline, combinações de condição e efeito) | Médio | Média | O terceiro frontend apoia as telas do app do técnico da subcontratada (proposta); o teste de três contas é escrito pelo backend da trilha SRV, com revisão do QA; QA focado na conexão em S10 e S11 |
| 24 | Dúvida jurídica sobre os papéis de cada conta no compartilhamento de dados entre contas conectadas (LGPD) | Médio | Média | Termo da conexão aprovado até S7 (seção 7); sem ele, convites só em homologação; dados pessoais fora da lista fechada (CPF, documentos) |
| 25 | Conexão entre contas pouco usada no lançamento, porque depende de cliente e subcontratada estarem na Tenaz ao mesmo tempo | Médio | Média | Piloto com pelo menos uma subcontratada conectada a uma empresa piloto; o SRV tem valor sem conexão (relatório em PDF); desconto de canal (P4) como incentivo |
| 26 | **R7 do sumário:** a configuração cresce até virar um sistema "faça você mesmo", com estados e fluxos livres; indicadores deixam de ser comparáveis, os testes não cobrem as combinações e a implantação fica cara | Alto | Média | Catálogos fechados de condição e efeito e estados, fórmulas, regras do app offline, isolamento e lista da conexão fixos (sumário 3.7); condição ou efeito novo só por requisito em 1_ESPECIFICACAO_REQUISITOS.md, com teste, entrando na etapa seguinte e nunca no meio de uma sprint; limites de 50 regras ativas e 30 campos extras por combinação (RF-ADM-016, RF-OS-021); pacotes mantidos pela Tenaz, sem editor de pacotes para o cliente; pedido de piloto por "mais flexibilidade" tratado como mudança de escopo, levada ao dono do produto |
| 27 | Configuração por segmento aumenta o escopo do MVP (3 requisitos, cerca de 9 a 10 sprints-pessoa) e da Fase 2 (13 requisitos, cerca de 13 a 15 sprints-pessoa) sem decisão sobre time e prazo | Alto | Alta | Encaixe depois do que cada requisito reaproveita (seção 4); opções A, B e C para o MVP e 1, 2 e 3 para a Fase 2 apresentadas ao dono do produto (seção 2.2); decisão até o fim de S3, para que um eventual backend a mais entre em S4; cortes parciais candidatos listados, sem aplicação automática |
| 28 | **P9** ou **P10** sem resposta antes da trilha E | Médio | Média | P9 até o mês 10, com o pacote saindo com validade digitada e regras inativas; P10 até o mês 9, com painel interno e medição em PDF e XLSX (contingências da seção 6) |
| 29 | **Multipaís e idiomas aumentam o escopo do MVP** (cerca de 17 a 20 sprints-pessoa além da E01) e da Fase 2 (cerca de 10 a 13), sobre uma folga que já foi consumida pela configuração por segmento e pelo SRV; a contratação do apoio externo de dois países leva tempo | Alto | Alta | Encaixe nas etapas em que cada parte aparece (seção 4); opções A a D para o MVP e três para a Fase 2 levadas ao dono do produto (seção 2.3, **P15**), com decisão até o fim de S3; consumo real da E01 registrado como primeiro dado de velocidade; cortes parciais listados, sem aplicação automática |
| 30 | **Conteúdo legal da Argentina e do Paraguai sem parecer** (P12, P13): prazos provinciais, norma em disputa (revisão técnica de veículo no Paraguai), regra oficial do dígito do CUIT e do RUC não confirmada; um prazo errado no produto vira descumprimento atribuído à Tenaz | Alto | Alta | Nenhum prazo legal no código (RNF-046); prazo não confirmado nasce vazio e a empresa o preenche; CUIT e RUC só pela forma; consultor de cada país antes de liberar o país; portão por país |
| 31 | **Qualidade do espanhol:** cada tela, e-mail e erro novo duplica o texto; deriva entre catálogos, termo errado do domínio, voseo, ou espanhol "neutro" que soa estranho num dos dois países | Médio | Média | Glossário bilíngue como fonte única; tarefa `textos` e teste de paridade quebram o build; revisão do `es-419` a cada mudança e por falante nativo (um argentino e um paraguaio) antes do piloto e do lançamento |
| 32 | **Proteção de dados e transferência internacional** (P14): o Brasil não está na lista de países adequados da autoridade argentina e a lei paraguaia de 2025 entra em vigor por volta de novembro de 2027; sem parecer, não se vende no país | Alto | Média | Advogado local contratado em S4; aditivo por país no contrato; portão por país; hospedagem só em São Paulo com cláusulas modelo, como o ADR-013 descreve |
| 33 | **Preço, moeda e meio de pagamento** (P15): inflação e câmbio do peso argentino, guarani sem casas decimais, provedor de pagamento nos dois países | Médio | Alta | Decisão da P15 junto com a P4, até a semana 10; fatura manual como contingência (seção 6); valor e moeda sempre juntos no código (RNF-045) |
| 34 | **Erro de dinheiro entre Go e TypeScript** (arredondamento, rateio, casas da moeda) e **divergência de formatação** entre web, app e PDF (ponto 48 do documento 1) | Médio | Baixa | Os 266 casos compartilhados e o oráculo no CI (RNF-045); a divergência de formatação vai ao `tenaz-arquiteto` e ao `tenaz-design-sistema` antes de E02 |

---

**Histórico de versões**

- v0.1 (12/11/2026): primeira versão, em elaboração.
- v0.2 (12/11/2026): alinhada ao sumário executivo 0.3. TER básico no MVP; TER completo, EST estendido e LOC na Fase 2, com quarta trilha; MVP com 105 requisitos e Fase 2 com 62; time de 12 pessoas no MVP e 15 a partir da Fase 2; RNF-039 com esteira de versionamento em S0, selo de versão em S1 e critério de pronto de release (seção 3.3); riscos R4, P7 e P8; linha de corte revista.
- v0.3 (12/11/2026): alinhada ao sumário executivo 0.4. SRV completo no MVP, com trilha SRV paralela desde S1 e a tabela de onde entra cada parte; MVP com 126 requisitos funcionais e 42 não funcionais; estimativa do acréscimo do SRV (cerca de 16 sprints-pessoa, já descontado o reaproveitamento); quarto backend e terceiro frontend antecipados para S1 (MVP com 14 pessoas, Fase 2 com 15), prazo mantido em 13 sprints; RF-TER-001 antecipado para S6; portão da conexão em M4 (RNF-041 e RNF-042); linha de corte com itens do SRV; marcos, dependências, piloto de subcontratada, termo LGPD da conexão e riscos R5, R6 e 23 a 25.
- v0.4 (12/11/2026): alinhada ao sumário executivo 0.5. Configuração por segmento: RF-ADM-016, RF-ATV-013 (com o tipo de OS `implantacao`) e RF-OS-021 encaixados em S2 a S9, com a tabela de onde entra cada parte; MVP com 129 requisitos e Fase 2 com 75; estimativa do acréscimo (seção 2.2: cerca de 9 a 10 sprints-pessoa no MVP e 13 a 15 na Fase 2) e opções de time e prazo sem decisão (A absorver, B mais 1 sprint, C mais 1 backend; na Fase 2, 1 trilha E com reforço, 2 mais prazo, 3 requisitos movidos para a Fase 3); trilha E na Fase 2 com PRJ, CTR e pacotes de segmento, ordem e dependências; P9 e P10 nas atividades e nas contingências; classificação de `causa_externa` na lista da conexão apontada como pendente; marcos M4 a M6, dependências 20 a 25 e riscos 26 (R7), 27 e 28.
- v0.5 (13/11/2026): alinhada ao sumário executivo 0.7. Risco 24 sem a menção a dados de saúde, que a Tenaz não guarda (decisão do dono do produto). Sem mudança de cronograma.
- v0.6 (13/11/2026): alinhada ao sumário executivo 0.8 (etapa 014). Multipaís e bilíngue desde o MVP: seção 2.3 com a estimativa do impacto (cerca de 17 a 20 sprints-pessoa no MVP, além da E01, e 10 a 13 na Fase 2) e as opções A a D para o MVP e três para a Fase 2, sem decisão (pendência P15); tabela de onde entra cada parte por etapa; apoio externo da Argentina e do Paraguai; marco M5 com RNF-043 a RNF-046; atividades P12 a P15, revisão do espanhol e piloto por país; dependências 26 e 27; contingência do portão por país; riscos 29 a 34. Os não funcionais passam de 42 para 46. Sem mudança de sprint, prazo ou time da tabela da seção 2.
