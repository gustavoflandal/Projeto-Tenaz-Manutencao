# Etapas de Execução — Tenaz

**Versão:** 0.2 · **Data:** 12/11/2026 · **Status:** Em elaboração
**Produto:** Tenaz · Gestão de Manutenção e Frota
Alinhado ao roadmap 0.4 e ao sumário executivo 0.5.

**Mudanças da versão 0.1:** o sumário 0.5 acrescentou ao MVP a configuração por segmento. Os três requisitos novos foram acrescentados às entregas das etapas, conforme o [roadmap](../6_ROADMAP.md), seção 4: RF-ATV-013 em E02, E04, E06 e E08; RF-OS-021 em E04, E05 e E09; RF-ADM-016 em E06 e E09. O MVP passa a 129 requisitos e a Fase 2 a 75, com a trilha E (projetos, contratos com clientes e pacotes de segmento). O encaixe não muda o número de etapas; o prazo (13 ou 14 sprints, com ou sem reforço do time) é decisão do dono do produto, no roadmap, seção 2.2.

As etapas E00 a E12 correspondem uma a uma às sprints S0 a S12 do MVP no [roadmap](../6_ROADMAP.md), seção 4. O escopo é o mesmo; o que muda é a forma de executar, descrita em [1_MODELO_EXECUCAO.md](1_MODELO_EXECUCAO.md).

**Etapas não têm duração prevista.** Terminam quando o portão fecha e o dono do produto aceita. Os marcos M0 a M5 do roadmap continuam a referência de prazo para tudo que depende de terceiros (pilotos, pareceres, pentest, lojas).

---

## 1. Mapa das etapas

| Etapa | Sprint | Tema | Marco | Prompt |
|---|---|---|---|---|
| E00 | S0 | Fundação técnica | M0 | [E00](../prompts/E00_fundacao-tecnica.md) |
| E01 | S1 | Conta e acesso | | [E01](../prompts/E01_conta-e-acesso.md) |
| E02 | S2 | Ativos | | [E02](../prompts/E02_ativos.md) |
| E03 | S3 | Medidores e composição | M1 | [E03](../prompts/E03_medidores.md) |
| E04 | S4 | Solicitações, núcleo da OS e fundação da conexão | | [E04](../prompts/E04_solicitacoes-e-os.md) |
| E05 | S5 | Execução da OS e OS de serviço | | [E05](../prompts/E05_execucao-da-os.md) |
| E06 | S6 | Preventiva | M2 | [E06](../prompts/E06_preventiva.md) |
| E07 | S7 | Checklists, estoque básico e fornecedores | | [E07](../prompts/E07_checklists-estoque-fornecedores.md) |
| E08 | S8 | NR-12, documentos e assinatura | M3 | [E08](../prompts/E08_nr12-documentos-assinatura.md) |
| E09 | S9 | Indicadores e terceiros | | [E09](../prompts/E09_indicadores-e-terceiros.md) |
| E10 | S10 | API, SSO, cobrança e portão da conexão | M4 | [E10](../prompts/E10_api-sso-cobranca.md) |
| E11 | S11 | Endurecimento | | [E11](../prompts/E11_endurecimento.md) |
| E12 | S12 | Lançamento | M5 | [E12](../prompts/E12_lancamento.md) |

Fases 2 e 3: seções 4 e 5 deste documento.

---

## 2. Como ler uma etapa

Cada etapa abaixo traz:

- **Entrega:** os requisitos pelos IDs, como no [roadmap](../6_ROADMAP.md), seção 4. A redação de cada um está em [1_ESPECIFICACAO_REQUISITOS.md](../1_ESPECIFICACAO_REQUISITOS.md).
- **Agentes:** quem conduz, pela ordem de [1_MODELO_EXECUCAO.md](1_MODELO_EXECUCAO.md), seção 5.2.
- **Portão de saída:** o que precisa estar verde. Vermelho não fecha etapa.
- **Evidência:** o que vai para [../screenshots/](../screenshots/).
- **Depende de:** etapa anterior e pendências externas.

---

## 3. MVP

### E00 — Fundação técnica

| | |
|---|---|
| **Entrega** | Base de RNF-013, RNF-016, RNF-020 e RNF-032. Esteira de versionamento de RNF-039 e RNF-040: versão gerada da tag, `CHANGELOG.md` verificado pelo CI, release bloqueado em divergência, `GET /v1/versao`. Esquema da plataforma separado do RLS por empresa. Projeto do app com banco local e prova de conceito da fila de comandos. Rascunho de ADR-009. Sem requisito funcional. |
| **Agentes** | `tenaz-plataforma-sre` conduz; `tenaz-arquiteto` (estrutura do repositório, ADR-009); `tenaz-banco-dados` (base do esquema, migração inicial); `tenaz-isolamento` (suíte de vazamento vazia no CI); `tenaz-mobile-offline` (prova de conceito da fila); `tenaz-documentacao` |
| **Portão de saída** | CI verde; infraestrutura em São Paulo criada por código; teste de vazamento entre empresas rodando no CI; esteira de versão ativa com `GET /v1/versao` respondendo; build interno do app instalável; verificador de documentação sem erro agora **no CI** (pendência aberta na etapa 008) |
| **Evidência** | Captura do pipeline verde; captura do app instalado mostrando o selo de versão |
| **Depende de** | Nada no repositório. Fora dele: conta de nuvem e contas de desenvolvedor das lojas |

> Esta é a única etapa que pode começar sem nenhuma decisão de produto pendente. Começar por ela é o que torna todas as outras verificáveis.

### E01 — Conta e acesso

| | |
|---|---|
| **Entrega** | RF-ADM-001 (sem o roteiro de primeiros passos), RF-ADM-002, RF-ADM-004, RF-ADM-005, RF-ADM-006, RF-ADM-011, RF-ADM-015; selo de versão e tela "Sobre" na web e no app (RNF-039); RF-MOB-001; RF-SRV-001 parcial (conta com o módulo SRV, perfis, escolha de empresa no login) |
| **Agentes** | `tenaz-produto`, `tenaz-banco-dados`, `tenaz-isolamento`, `tenaz-backend-go`, `tenaz-api-contrato`, `tenaz-design-sistema`, `tenaz-frontend-web`, `tenaz-mobile-offline`, `tenaz-qa-automacao`, `tenaz-documentacao` |
| **Portão de saída** | Permissões no formato `<módulo>.<ação>` com `os.encerrar` e `est.decidir_destino`; escopos `interno`, `prestador`, `locatario`, `link` existentes no modelo; trilha de auditoria gravando valores antes e depois; teste de vazamento cobrindo as tabelas novas; selo de versão visível no login e na tela principal |
| **Evidência** | Login, tela principal com o selo de versão, "Sobre", cadastro de pessoa, cadastro de usuário com perfis, trilha de auditoria |
| **Depende de** | E00 |

> `pessoa` (RF-ADM-015) entra aqui por dependência crítica 9 do roadmap: criar a entidade depois exigiria migrar usuários, motoristas e documentos pessoais.

### E02 — Ativos

| | |
|---|---|
| **Entrega** | RF-ATV-001, RF-ATV-002, RF-ATV-003, RF-ATV-005, RF-ATV-006, RF-ATV-007, RF-ATV-008; RF-INT-003 para localizações e ativos; RF-MOB-002 (cache de ativos), RF-MOB-004; RF-SRV-004 (clientes atendidos), RF-SRV-002 (especialidades de técnico); RF-ATV-013 parcial (fase de ciclo de vida, ativo `em_implantacao` com status `inativo`, cadastro e importação em implantação, filtro por fase) |
| **Agentes** | `tenaz-dominio-manutencao` (hierarquia, criticidade, status), `tenaz-banco-dados`, `tenaz-isolamento`, `tenaz-backend-go`, `tenaz-api-contrato`, `tenaz-design-sistema`, `tenaz-frontend-web`, `tenaz-mobile-offline`, `tenaz-qa-automacao` |
| **Portão de saída** | Árvore de localizações e de ativos com até 6 níveis, sem ciclo; status operacional do ativo conforme a regra do [sumário](../0_SUMARIO_EXECUTIVO.md), seção 3.3, com `parado` prevalecendo sobre `em_manutencao`; importação de 2.000 ativos por planilha; QR code resolvendo para o ativo |
| **Evidência** | Lista de ativos, ficha do ativo, árvore de localizações, importação por planilha, leitura de QR no app |
| **Depende de** | E01 |

### E03 — Medidores e composição

| | |
|---|---|
| **Entrega** | RF-MED-001 a RF-MED-008; RF-ATV-004, RF-ATV-009, RF-ATV-011; RF-ATV-010 parcial; RF-MOB-003 parcial (leitura offline); RF-SRV-005 (equipamentos de clientes), RF-SRV-003 parcial (técnicos com especialidades e capacitações) |
| **Agentes** | `tenaz-dominio-manutencao` (validação de leitura, troca de medidor), `tenaz-banco-dados`, `tenaz-backend-go`, `tenaz-api-contrato`, `tenaz-frontend-web`, `tenaz-mobile-offline`, `tenaz-qa-automacao` |
| **Portão de saída** | **Marco M1:** empresa de teste com 2.000 ativos importados em menos de 2 minutos; leituras validadas com regressão, teto de salto e troca de medidor; pessoas com e sem usuário; conta de subcontratada de teste com clientes, equipamentos `de_cliente` e técnicos com especialidades. Casos compartilhados de validação de leitura passando nas duas linguagens |
| **Evidência** | Registro de leitura na web e no app, histórico de leituras, troca de medidor, ficha de equipamento de cliente |
| **Depende de** | E02. Dependência crítica 1 do roadmap: sem leitura confiável, a preventiva por medidor gera OS errada |

### E04 — Solicitações, núcleo da OS e fundação da conexão

| | |
|---|---|
| **Entrega** | RF-SOL-001 a RF-SOL-005; RF-OS-001 a RF-OS-005, RF-OS-016, RF-OS-017; RF-ADM-012; RF-MOB-006 e solicitação offline (RF-MOB-003). Fundação da conexão, sem requisito funcional: ADR-009 fechado, entidade de conexão no esquema da plataforma, papel de banco do sincronizador, fila de saída nos dois lados, idempotência, **suíte de três contas rodando no CI**. Configuração por segmento: origens `projeto` e `regra` e `causa_externa` na OS (RF-OS-001); RF-OS-021 parcial (definições de `campo_extra_os`, valores, obrigatoriedade na abertura, web, API e app offline); RF-ATV-013 parcial (tipo `implantacao`, ação `colocar_em_operacao`, também em lote) |
| **Agentes** | `tenaz-arquiteto` (fecha ADR-009), `tenaz-dominio-manutencao` (estados e marcos de tempo), `tenaz-conexao-contas`, `tenaz-isolamento`, `tenaz-banco-dados`, `tenaz-backend-go`, `tenaz-api-contrato`, `tenaz-frontend-web`, `tenaz-mobile-offline`, `tenaz-qa-automacao` |
| **Portão de saída** | Máquina de estados da OS idêntica no servidor e no app pelos casos compartilhados; `data_reconhecimento` gravado uma única vez; transições proibidas recusadas com erro de código estável; suíte de três contas verde (ainda sem fluxo de OS entre contas); primeira versão do app nas lojas em teste fechado |
| **Evidência** | Portal de solicitação, triagem, ficha da OS, OS aberta no app em modo avião |
| **Depende de** | E03. Dependência crítica 16 do roadmap: a suíte de três contas vem **antes** de qualquer fluxo entre contas |

### E05 — Execução da OS e OS de serviço

| | |
|---|---|
| **Entrega** | RF-OS-006, RF-OS-007 (com a ação de assumir), RF-OS-008, RF-OS-009, RF-OS-011, RF-OS-015, RF-OS-018; execução de OS offline (RF-MOB-003), RF-MOB-005; RF-SRV-006 (catálogo próprio), RF-SRV-007 (abertura com equipamento, técnico e serviço a realizar); RF-OS-021 parcial (obrigatoriedade na conclusão, validada no app offline antes de enfileirar o comando; campos extras na impressão) |
| **Agentes** | `tenaz-dominio-manutencao`, `tenaz-backend-go`, `tenaz-api-contrato`, `tenaz-conexao-contas` (catálogo reaproveitado), `tenaz-design-sistema`, `tenaz-frontend-web`, `tenaz-mobile-offline`, `tenaz-qa-automacao` |
| **Portão de saída** | Pausa com motivo obrigatório; apontamento de horas por técnico; catálogo padrão de códigos de falha de modo, causa e atividade; impressão da OS; execução completa offline com reconexão sem perda |
| **Evidência** | Execução da OS no app passo a passo, apontamento, pausa com motivo, OS impressa em PDF, OS de serviço da subcontratada |
| **Depende de** | E04 |

### E06 — Preventiva

| | |
|---|---|
| **Entrega** | RF-PRE-001 a RF-PRE-006, RF-PRE-008, RF-PRE-009, RF-PRE-010; RF-ADM-010; OS preventiva no app; RF-SRV-009 (horas por técnico), RF-SRV-010 (distribuição e reatribuição), RF-TER-001 antecipado (fornecedor com papéis), RF-SRV-011 parcial; RF-ATV-013 parcial (planos suspensos em implantação e ativados na entrada em operação); RF-ADM-016 parcial (motor de regras, rotina diária idempotente, `adm.gerenciar_regras`, condições `leitura_fora_faixa` e `os_concluida_com_campo`, efeitos `notificar`, `abrir_solicitacao` e `bloquear_ativo`) |
| **Agentes** | `tenaz-dominio-manutencao` (gatilhos e antecipação), `tenaz-backend-go`, `tenaz-banco-dados`, `tenaz-api-contrato`, `tenaz-frontend-web`, `tenaz-mobile-offline`, `tenaz-qa-automacao` |
| **Portão de saída** | **Marco M2:** ciclo completo de solicitação a encerramento, com reconhecimento, assumir, desprogramação e liberação recusada; plano por tempo, por medidor e por "o que vencer primeiro" gerando OS `planejada`; **app operando 72 horas offline sem perda**; OS de serviço distribuída, reatribuída e executada offline por técnico de subcontratada |
| **Evidência** | Cadastro de plano com os três tipos de gatilho, calendário de preventivas, OS gerada por plano, distribuição de OS de serviço |
| **Depende de** | E03 (medidores) e E05 |

### E07 — Checklists, estoque básico e fornecedores

| | |
|---|---|
| **Entrega** | RF-CHK-001 a RF-CHK-006; RF-SOL-006; RF-EST-001 a RF-EST-004, RF-EST-006; RF-OS-010; RF-FRO-001; RF-TER-002; checklist offline, RF-MOB-007; RF-SRV-014 (convite, aceite e campos compartilhados) |
| **Agentes** | `tenaz-dominio-manutencao` (reprovação gera solicitação, reserva e baixa), `tenaz-banco-dados`, `tenaz-backend-go`, `tenaz-api-contrato`, `tenaz-conexao-contas`, `tenaz-isolamento`, `tenaz-frontend-web`, `tenaz-mobile-offline`, `tenaz-qa-automacao` |
| **Portão de saída** | Item reprovado em checklist gerando solicitação; reserva e baixa de peça em OS com estoque mínimo e ajuste com motivo; convite de conexão funcionando **só em homologação** até o termo de LGPD da conexão estar aprovado (dependência crítica 18 do roadmap) |
| **Evidência** | Modelo de checklist, execução no app, item reprovado gerando solicitação, saldo de estoque, baixa em OS, convite de conexão |
| **Depende de** | E05, E06 e RF-TER-001 (E06). Fora do repositório: termo da conexão aprovado pelo advogado |

### E08 — NR-12, documentos e assinatura

| | |
|---|---|
| **Entrega** | RF-CON-001 a RF-CON-006, RF-CON-014, RF-CON-016; RF-OS-013, RF-OS-014; campos de NR-12 e assinatura no app; RF-SRV-008 (componentes substituídos), RF-SRV-015 (ativos cobertos pela conexão); RF-ATV-013 completo (documentos críticos exigidos para `colocar_em_operacao`); custo da OS `implantacao` separado do custo de manutenção (RF-OS-013) |
| **Agentes** | `tenaz-conformidade-br` conduz a interpretação da norma; `tenaz-dominio-manutencao`, `tenaz-backend-go`, `tenaz-api-contrato`, `tenaz-conexao-contas`, `tenaz-design-sistema`, `tenaz-frontend-web`, `tenaz-mobile-offline`, `tenaz-qa-automacao` |
| **Portão de saída** | **Marco M3:** pelo menos 3 empresas piloto com dados reais; OS de máquina NR-12 não conclui sem os sete campos do item 12.11.2, com assinatura avançada; alertas de documento em D-60, D-30 e D-7; reserva e baixa de peças; subcontratada piloto registrando OS de serviço com componentes; convite, aceite e ativos sincronizados entre duas contas em homologação |
| **Evidência** | Conclusão de OS de NR-12 com os sete campos, assinatura com manifesto de evidências, agenda de documentos, alerta de vencimento, componentes substituídos no app |
| **Depende de** | E07. **Pendência P1** (parecer sobre ICP-Brasil) com prazo até esta etapa; sem ela, segue a contingência do [roadmap](../6_ROADMAP.md), seção 6 |

### E09 — Indicadores e terceiros

| | |
|---|---|
| **Entrega** | RF-IND-001 a RF-IND-005, RF-IND-007; RF-ATV-010 completo; RF-TER-003, RF-TER-004, RF-TER-005; RF-MOB-010 e o perfil de prestador no app; RF-SRV-012, RF-SRV-013 (com os campos extras da OS), RF-SRV-003 completo, RF-SRV-016 (OS do cliente chega à subcontratada); KPIs sem OS `implantacao` nem ativo `em_implantacao`; RF-ADM-016 completo (condições de documento e de checklist, `bloquear_ativo_ate_documento`, simulação, histórico, efeito na conta conectada só pelo status no teste de três contas); RF-OS-021 completo (filtro, exportação, link de uso único) |
| **Agentes** | `tenaz-dominio-manutencao` (fórmulas), `tenaz-isolamento` (perfil de prestador e link de uso único), `tenaz-backend-go`, `tenaz-api-contrato`, `tenaz-conexao-contas`, `tenaz-design-sistema`, `tenaz-frontend-web`, `tenaz-mobile-offline`, `tenaz-qa-automacao` |
| **Portão de saída** | **Teste de isolamento do prestador (RNF-037) verde no CI** — sem ele o perfil de prestador não é liberado (risco R4); link de uso único com token no fragmento da URL, recusando reuso, outro aparelho e revogação; todo indicador com o texto de "como calculamos" (RF-IND-002) |
| **Evidência** | Painel de indicadores, "como calculamos" aberto, OS atribuída a fornecedor, portal do prestador, página do link de uso único, relatório da OS em PDF com assinatura do cliente |
| **Depende de** | E08 |

### E10 — API, SSO, cobrança e portão da conexão

| | |
|---|---|
| **Entrega** | RF-INT-001, RF-INT-002, RF-INT-003 completo; RF-ADM-007, RF-ADM-009, RF-ADM-013; SSO no app; RF-SRV-017, RF-SRV-018, RF-SRV-019, RF-SRV-020, RF-SRV-021, RF-SRV-001 com cobrança, RF-SRV-011 completo |
| **Agentes** | `tenaz-api-contrato` conduz; `tenaz-conexao-contas`, `tenaz-isolamento`, `tenaz-backend-go`, `tenaz-plataforma-sre` (chaves, limites), `tenaz-frontend-web`, `tenaz-mobile-offline`, `tenaz-qa-automacao` |
| **Portão de saída** | **Marco M4 e portão da conexão:** os 129 requisitos do MVP entregues; RNF-041 verde e atraso p95 de RNF-042 abaixo de 60 s medido em homologação; ciclo completo entre contas, inclusive liberação recusada e encerramento da conexão sem perda de histórico; documentação pública da API publicada; nenhum defeito crítico aberto. Portão vermelho aqui dispara a contingência do [roadmap](../6_ROADMAP.md), seção 6, com aval do dono do produto |
| **Evidência** | Documentação pública da API, tela de chaves e escopos, entrega de webhook, login por SSO, ciclo de OS entre duas contas lado a lado, relatório de serviços por cliente |
| **Depende de** | E09 |

### E11 — Endurecimento

| | |
|---|---|
| **Entrega** | Sem requisito funcional novo: carga (RNF-007 a RNF-009, RNF-030), **7 dias offline** (RNF-010), restauração de backup (RNF-005), acessibilidade (RNF-025), pentest focado em isolamento, acesso do prestador, link de uso único e conexão entre contas, primeiro release candidato em homologação, correções dos pilotos |
| **Agentes** | `tenaz-qa-automacao` conduz; `tenaz-plataforma-sre`, `tenaz-isolamento`, `tenaz-conexao-contas`, `tenaz-design-sistema` (acessibilidade), `tenaz-mobile-offline` (aparelhos de entrada) |
| **Portão de saída** | p95 dentro de RNF-007 com massa de 2.000 ativos; 7 dias offline sem perda; restauração de backup ensaiada com RPO e RTO medidos; WCAG 2.1 AA verificado automática e manualmente; pentest sem achado crítico ou alto aberto; atraso da sincronização entre contas medido **sob carga** |
| **Evidência** | Relatório de carga, relatório de acessibilidade, captura do app após 7 dias offline, painel de observabilidade com o atraso da sincronização |
| **Depende de** | E10 |

### E12 — Lançamento

| | |
|---|---|
| **Entrega** | RF-ADM-001 completo (roteiro de primeiros passos), página de status (RNF-020), documentação pública da API, release `v1.0.0`, migração dos pilotos para produção, RF-SRV-001 completo, publicação nas lojas |
| **Agentes** | `tenaz-plataforma-sre` conduz o release; `tenaz-produto`, `tenaz-design-sistema` (roteiro de primeiros passos), `tenaz-frontend-web`, `tenaz-mobile-offline`, `tenaz-documentacao` |
| **Portão de saída** | **Marco M5:** metas de RNF verificadas; pentest sem achado crítico ou alto aberto; kit LGPD publicado, com o termo da conexão entre contas; página de status no ar; RPO e RTO publicados, com o RTO de região só depois de **P7**; preços publicados após **P4**; marca liberada (**P5**); release `v1.0.0` conforme o [roadmap](../6_ROADMAP.md), seção 3.3 |
| **Evidência** | Primeiros passos guiados, página de status, página de preços, selo `v1.0.0` em produção, app publicado nas lojas |
| **Depende de** | E11 e as pendências P4, P5 e P7 resolvidas pelo dono do produto |

---

## 4. Fase 2 — meses 7 a 12

Os 75 requisitos e as cinco trilhas (a trilha E é proposta, com time e prazo a decidir pelo dono do produto) estão no [roadmap](../6_ROADMAP.md), seção 5, com entrega verificável por mês na seção 5.1. Este plano não os redivide; acrescenta só a atribuição de agentes por trilha.

| Trilha do roadmap | Conteúdo | Agente que conduz | Apoio |
|---|---|---|---|
| **A** — frota e conformidade | Pneus, combustível, multas, sinistros, NR-13, LOTO, aptidão para despacho | `tenaz-conformidade-br` nas normas; `tenaz-dominio-manutencao` na frota | `tenaz-backend-go`, `tenaz-frontend-web`, `tenaz-mobile-offline` |
| **B** — estoque completo e estendido | Inventário, compras, NF-e, vários almoxarifados, estoque em terceiros, peças serializadas e reparáveis, ordem de reparo | `tenaz-dominio-manutencao`; `tenaz-conformidade-br` no fiscal (**P8**) | `tenaz-banco-dados`, `tenaz-backend-go`, `tenaz-isolamento` |
| **C** — integrações, identidade e IA | Multi-unidade, capacitações, telemática, ERP, SAML e SCIM, IA nível 1 | `tenaz-integracoes-dados` | `tenaz-api-contrato`, `tenaz-plataforma-sre` |
| **D** — terceiros, locação e SRV | Contratos de prestação com SLA, avaliação, locação completa, portal do locatário, evolução do SRV pedida pelos pilotos | `tenaz-conexao-contas` e `tenaz-dominio-manutencao` | `tenaz-isolamento` (RNF-038), `tenaz-frontend-web` |
| **E** — projetos, contratos com clientes e pacotes | Projetos de implantação com etapas (PRJ), contratos com clientes com meta de disponibilidade e medição mensal (CTR), pacotes de segmento | `tenaz-dominio-manutencao` | `tenaz-produto` (conteúdo dos pacotes, P10), `tenaz-conformidade-br` (P9), `tenaz-conexao-contas` (OS de etapa na conta conectada), `tenaz-frontend-web`, `tenaz-mobile-offline` |

Portão da fase: **marco M6** do [roadmap](../6_ROADMAP.md), seção 3.2, com os testes de isolamento do locatário e do almoxarifado em terceiro verdes no CI e a janela de observação do SOC 2 iniciada.

---

## 5. Fase 3 — meses 13 a 24

Os 16 requisitos estão no [roadmap](../6_ROADMAP.md), seção 6, um bloco por mês. Conduzem `tenaz-integracoes-dados` (IoT, IA, telemática, BI), `tenaz-dominio-manutencao` (manutenção por condição, confiabilidade, RCM e FMEA, ISO 55001) e `tenaz-plataforma-sre` (banco dedicado, SOC 2 Type II, preparação da ISO 27001).

Portões: **M7** (relatório SOC 2 Type II, mês 19) e **M8** (fase completa, mês 24).

---

## 6. Regras das etapas

1. **Uma etapa por vez.** Etapas não se sobrepõem. Dentro de uma etapa, trilhas independentes podem andar em paralelo.
2. **Etapa não começa sem prompt** escrito em [../prompts/](../prompts/).
3. **Etapa não fecha com portão vermelho.** Não existe fechamento com ressalva. O que falta vira pendência registrada e desloca o que depende dela.
4. **Etapa não fecha sem registro** em [../etapas/](../etapas/) e sem o verificador de documentação limpo.
5. **Escopo de etapa não cresce no meio.** Pedido novo durante a etapa vira etapa seguinte ou pendência, por decisão do dono do produto.
6. **Requisito entregue em partes** conta como pronto só na última etapa em que aparece, como no [roadmap](../6_ROADMAP.md), seção 4.

---

**Histórico de versões**

- v0.1 (12/11/2026): primeira versão. Etapas E00 a E12 correspondentes às sprints S0 a S12, com entrega, agentes, portão de saída e evidência; atribuição de agentes às trilhas da fase 2 e aos blocos da fase 3.
- v0.2 (12/11/2026): alinhada ao roadmap 0.4 e ao sumário executivo 0.5. Configuração por segmento acrescentada às entregas: RF-ATV-013 (E02, E04, E06, E08), RF-OS-021 (E04, E05, E09) e RF-ADM-016 (E06, E09); portão de E10 com 129 requisitos; Fase 2 com 75 requisitos e a trilha E (PRJ, CTR e pacotes), com agentes.
