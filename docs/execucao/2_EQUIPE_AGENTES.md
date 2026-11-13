# Equipe de Agentes — Tenaz

**Versão:** 0.1 · **Data:** 12/11/2026 · **Status:** Em elaboração
**Produto:** Tenaz · Gestão de Manutenção e Frota

Catálogo dos 16 agentes especializados que executam o projeto. A definição executável de cada um está em [../../.claude/agents/](../../.claude/agents/); este documento explica a divisão de trabalho, o que cada agente é dono e quando chamá-lo. As regras de orquestração estão em [1_MODELO_EXECUCAO.md](1_MODELO_EXECUCAO.md), seção 5.

---

## 1. Visão geral

```
CONDUÇÃO          tenaz-arquiteto · tenaz-produto · tenaz-dominio-manutencao
DADOS             tenaz-banco-dados · tenaz-isolamento
SERVIDOR          tenaz-backend-go · tenaz-api-contrato · tenaz-conexao-contas
INTERFACE         tenaz-design-sistema · tenaz-frontend-web · tenaz-mobile-offline
QUALIDADE         tenaz-qa-automacao · tenaz-plataforma-sre
ESPECIALIDADE     tenaz-conformidade-br · tenaz-integracoes-dados
MEMÓRIA           tenaz-documentacao
```

| Agente | Especialidade em uma linha |
|---|---|
| [`tenaz-arquiteto`](../../.claude/agents/tenaz-arquiteto.md) | Decide arquitetura, escreve ADR, guarda as fronteiras entre módulos |
| [`tenaz-produto`](../../.claude/agents/tenaz-produto.md) | Dono dos requisitos, do aceite e da rastreabilidade requisito → código → teste |
| [`tenaz-dominio-manutencao`](../../.claude/agents/tenaz-dominio-manutencao.md) | Especialista em manutenção e frota: estados, regras de negócio, indicadores |
| [`tenaz-banco-dados`](../../.claude/agents/tenaz-banco-dados.md) | Modelo físico, DDL, migrações, consultas tipadas, desempenho |
| [`tenaz-isolamento`](../../.claude/agents/tenaz-isolamento.md) | Isolamento entre empresas, escopos externos e contas conectadas |
| [`tenaz-backend-go`](../../.claude/agents/tenaz-backend-go.md) | Domínio e casos de uso em Go no monólito modular |
| [`tenaz-api-contrato`](../../.claude/agents/tenaz-api-contrato.md) | Contrato único: OpenAPI, erros, paginação, idempotência, tipos gerados |
| [`tenaz-conexao-contas`](../../.claude/agents/tenaz-conexao-contas.md) | Sincronização entre contas conectadas do módulo SRV |
| [`tenaz-design-sistema`](../../.claude/agents/tenaz-design-sistema.md) | Tokens da marca, componentes, contraste, microcopy, telas |
| [`tenaz-frontend-web`](../../.claude/agents/tenaz-frontend-web.md) | Web de gestão, portais restritos, página do link de uso único |
| [`tenaz-mobile-offline`](../../.claude/agents/tenaz-mobile-offline.md) | App de campo e a sincronização offline |
| [`tenaz-qa-automacao`](../../.claude/agents/tenaz-qa-automacao.md) | Pirâmide de testes, fluxos E2E obrigatórios, carga |
| [`tenaz-plataforma-sre`](../../.claude/agents/tenaz-plataforma-sre.md) | Infraestrutura como código, CI/CD, esteira de versão, observabilidade, backup |
| [`tenaz-conformidade-br`](../../.claude/agents/tenaz-conformidade-br.md) | NR-12, NR-13, LOTO, documentos, LGPD, assinatura e tratamento fiscal |
| [`tenaz-integracoes-dados`](../../.claude/agents/tenaz-integracoes-dados.md) | Telemática, cartão combustível, ERP, IoT, IA e indicadores avançados |
| [`tenaz-documentacao`](../../.claude/agents/tenaz-documentacao.md) | Integridade da documentação, CHANGELOG, registro de etapas, PDFs |

---

## 2. Correspondência com o time do roadmap

O [roadmap](../6_ROADMAP.md), seção 2, prevê 15 pessoas no MVP. A tabela mostra como os papéis viraram agentes. Papéis de apoio externo não viram agente: continuam sendo pessoas contratadas.

| Papel do roadmap | Qtd. | Agente correspondente |
|---|---|---|
| Gerente de produto | 1 | `tenaz-produto` (parte; aceite e pendências ficam com o dono do produto) |
| Designer de produto | 1 | `tenaz-design-sistema` |
| Líder técnico (Go) | 1 | `tenaz-arquiteto` |
| Backend Go | 4 | `tenaz-backend-go`, `tenaz-banco-dados`, `tenaz-api-contrato`, `tenaz-conexao-contas` |
| Frontend React/TypeScript | 3 | `tenaz-frontend-web` |
| Mobile React Native | 2 | `tenaz-mobile-offline` |
| QA com automação | 1 | `tenaz-qa-automacao`, com `tenaz-isolamento` nas suítes de vazamento |
| Plataforma/SRE | 1 | `tenaz-plataforma-sre` |
| Engenharia de dados e IA (fase 2) | 1 | `tenaz-integracoes-dados` |
| — (sem equivalente no roadmap) | — | `tenaz-dominio-manutencao`, `tenaz-conformidade-br`, `tenaz-documentacao` |
| Advogado, consultor de SST, contador, pentest, auditor SOC 2 | parcial | **Não vira agente.** `tenaz-conformidade-br` prepara o material e as perguntas; a decisão é da pessoa. |

Os três agentes sem equivalente existem porque este projeto tem exigências que, num time humano, ficariam distribuídas: o domínio de manutenção com 376 regras de negócio e fórmulas de indicador, a conformidade brasileira nativa que é a proposta de valor 4 do sumário, e a disciplina documental de [CLAUDE.md](../../CLAUDE.md), seção 3, que é condição de fechamento de qualquer etapa.

---

## 3. Condução

### 3.1 `tenaz-arquiteto`

| | |
|---|---|
| **Dono de** | [5_ARQUITETURA_SOFTWARE.md](../5_ARQUITETURA_SOFTWARE.md), inclusive os `ADR-`; estrutura do repositório; regras de dependência entre módulos |
| **Decide** | Fronteira de módulo, forma de um módulo chamar outro, eventos internos e outbox, uso de réplica de leitura, entrada de qualquer componente novo na stack |
| **Chamar quando** | Há decisão técnica sem ADR; um módulo precisa de dado de outro; aparece proposta de componente novo; há dúvida sobre a camada onde algo mora |
| **Não faz** | Implementação de módulo, escolha de escopo, decisão de preço |
| **Portão que guarda** | Nenhum componente novo entra sem gatilho mensurável (princípio 2 da arquitetura) |

### 3.2 `tenaz-produto`

| | |
|---|---|
| **Dono de** | [1_ESPECIFICACAO_REQUISITOS.md](../1_ESPECIFICACAO_REQUISITOS.md) (`RF-`, `RNF-`), [6_ROADMAP.md](../6_ROADMAP.md) |
| **Decide** | Redação de requisito, critério de aceite, rastreabilidade requisito → código → teste → CHANGELOG |
| **Chamar quando** | Começa uma etapa; um requisito está ambíguo; é preciso saber se algo está no MVP ou na fase 2; o escopo parece crescer |
| **Não faz** | Aceite final (é do dono do produto), decisão sobre P1 a P8, aplicação da linha de corte |
| **Portão que guarda** | Os 9 itens do critério de pronto ([roadmap](../6_ROADMAP.md), seção 3.1) |

### 3.3 `tenaz-dominio-manutencao`

| | |
|---|---|
| **Dono de** | [4_FLUXOS_PROCESSO.md](../4_FLUXOS_PROCESSO.md): os 31 fluxos, as regras `RN-`, o cálculo dos indicadores e os parâmetros configuráveis |
| **Decide** | Comportamento esperado no domínio: transições de estado, efeito de uma falha na disponibilidade, fórmula de indicador, gatilho de plano preventivo |
| **Chamar quando** | Implementa-se qualquer regra de OS, plano, medidor, checklist ou indicador; há dúvida sobre o que significa um termo do domínio; um KPI precisa do "como calculamos" (RF-IND-002) |
| **Não faz** | Código de produção, modelo físico de dados |
| **Portão que guarda** | Máquina de estados da OS e validação de leitura idênticas no servidor e no app, pelos casos compartilhados em JSON |

---

## 4. Dados

### 4.1 `tenaz-banco-dados`

| | |
|---|---|
| **Dono de** | [2_ARQUITETURA_BANCO_DADOS.md](../2_ARQUITETURA_BANCO_DADOS.md), migrações, consultas tipadas, índices |
| **Decide** | Colunas, tipos, chaves, índices, estratégia de migração expandir e contrair, numeração de documentos, particionamento |
| **Chamar quando** | Nasce ou muda uma tabela; uma consulta está lenta; é preciso migrar dado; há dúvida sobre volume ou crescimento |
| **Não faz** | Aprovar a própria política de isolamento (revisão é de `tenaz-isolamento`) |
| **Portão que guarda** | Toda migração compatível com a versão anterior do código; remoção de coluna só no release seguinte |

### 4.2 `tenaz-isolamento`

| | |
|---|---|
| **Dono de** | Suítes de vazamento ([arquitetura](../5_ARQUITETURA_SOFTWARE.md), seção 9.3), políticas de RLS, lista de tabelas globais, lista de acesso externo, lista fechada de campos da conexão |
| **Decide** | Se uma tabela pode ser lida por escopo externo; se um campo pode trafegar entre contas; se o portão está verde |
| **Chamar quando** | Nasce qualquer tabela; nasce qualquer rota; muda qualquer campo que o prestador, o locatário ou uma conta conectada possa ver; antes de fechar qualquer etapa |
| **Não faz** | Implementação de regra de negócio |
| **Portão que guarda** | Os três portões de isolamento: entre empresas (RNF-013), de acesso externo (RNF-037, RNF-038) e entre contas conectadas (RNF-041) |

> Este é o agente com poder de veto mais amplo do projeto, por decisão de risco: R2, R4 e R5 do [sumário](../0_SUMARIO_EXECUTIVO.md) são os únicos classificados como críticos.

---

## 5. Servidor

### 5.1 `tenaz-backend-go`

| | |
|---|---|
| **Dono de** | `backend/internal/modulos/`, domínio e casos de uso, erros tipados, jobs |
| **Decide** | Organização interna de um módulo dentro das camadas definidas pela arquitetura |
| **Chamar quando** | Há regra de negócio para implementar, job para agendar, evento interno para publicar |
| **Não faz** | Criar tabela, definir rota pública, decidir fronteira entre módulos |
| **Portão que guarda** | Nenhum SQL concatenado, nenhum estado global mutável, toda chamada externa com timeout, retry e disjuntor |

### 5.2 `tenaz-api-contrato`

| | |
|---|---|
| **Dono de** | [3_ESPECIFICACAO_APIS.md](../3_ESPECIFICACAO_APIS.md), `backend/api/openapi.yaml`, `pacotes/api-tipos`, códigos de erro, eventos de webhook |
| **Decide** | Forma da rota, nome do campo no JSON, código de erro, paginação, idempotência, concorrência por ETag, escopo de token |
| **Chamar quando** | Nasce ou muda um endpoint ou um evento de webhook; o contrato divergiu do código; é preciso saber se uma mudança quebra compatibilidade |
| **Não faz** | Regra de negócio, consulta ao banco |
| **Portão que guarda** | Contrato gerado igual ao versionado; mudança incompatível em `/v1` bloqueia o merge; web e app usam os mesmos endpoints da API pública |

### 5.3 `tenaz-conexao-contas`

| | |
|---|---|
| **Dono de** | ADR-009, esquema da plataforma, processo de sincronização entre contas, vínculo de OS, suíte de três contas junto com `tenaz-isolamento` |
| **Decide** | Formato da mensagem entre contas, idempotência e ordem, tratamento de falha e reprocessamento, efeito de suspensão e encerramento |
| **Chamar quando** | Qualquer trabalho em RF-SRV-014 a RF-SRV-019; um campo novo de OS ou de ativo precisa ser classificado; a sincronização atrasa ou falha |
| **Não faz** | Liberar o portão da conexão sozinho (é de `tenaz-isolamento`) |
| **Portão que guarda** | RNF-041 verde e RNF-042 medido antes de qualquer convite fora de homologação |

---

## 6. Interface

### 6.1 `tenaz-design-sistema`

| | |
|---|---|
| **Dono de** | [7_PADROES_DESIGN.md](../7_PADROES_DESIGN.md), `pacotes/tokens-ui`, `pacotes/i18n`, [marca/README.md](../../marca/README.md) na aplicação à interface |
| **Decide** | Token, componente, cor de status, contraste, ícone, microcopy, layout de tela |
| **Chamar quando** | Nasce ou muda qualquer tela; aparece cor, espaçamento ou texto fora dos tokens; há dúvida de acessibilidade |
| **Não faz** | Implementação de tela |
| **Portão que guarda** | WCAG 2.1 AA (RNF-025); texto sobre o laranja Sinal sempre em Aço; selo de versão visível (RNF-039) |

### 6.2 `tenaz-frontend-web`

| | |
|---|---|
| **Dono de** | `web/` |
| **Decide** | Estrutura de rota e de estado da web, dentro dos padrões da arquitetura |
| **Chamar quando** | Há tela de gestão, painel, portal do prestador ou do locatário, ou a página do link de uso único para construir |
| **Não faz** | Inventar campo que não está no contrato gerado; criar token de design |
| **Portão que guarda** | Dados do servidor só pelo cliente gerado do OpenAPI; textos só pelo pacote de internacionalização |

### 6.3 `tenaz-mobile-offline`

| | |
|---|---|
| **Dono de** | `app/`, fila de comandos, sincronização, banco local, sessões por conta |
| **Decide** | O que entra no cache, como o conflito é resolvido, numeração provisória, política de atualização sem loja |
| **Chamar quando** | Uma operação precisa funcionar offline; a sincronização apresenta conflito ou perda; o app precisa de recurso nativo |
| **Não faz** | Mudar regra de domínio para facilitar o offline |
| **Portão que guarda** | 7 dias offline sem perda de registro (RNF-010); troca de conta sem mistura de dados |

---

## 7. Qualidade e operação

### 7.1 `tenaz-qa-automacao`

| | |
|---|---|
| **Dono de** | A pirâmide de testes, os 13 fluxos E2E obrigatórios ([arquitetura](../5_ARQUITETURA_SOFTWARE.md), seção 9.2), os casos compartilhados em JSON, os testes de carga |
| **Decide** | Nível de teste adequado, massa de dados, cenário de falha, critério de aprovação |
| **Chamar quando** | Qualquer requisito é implementado; antes de fechar qualquer etapa; há defeito para reproduzir |
| **Não faz** | Aprovar o próprio teste de isolamento (revisão cruzada com `tenaz-isolamento`) |
| **Portão que guarda** | p95 dentro de RNF-007 com massa de 2.000 ativos; fluxo E2E obrigatório sempre verde |

### 7.2 `tenaz-plataforma-sre`

| | |
|---|---|
| **Dono de** | `infra/`, `.github/workflows/`, esteira de versionamento (RNF-039, RNF-040), observabilidade, backup e recuperação, [9_SEGURANCA_CONFIABILIDADE.md](../9_SEGURANCA_CONFIABILIDADE.md) na parte de operação |
| **Decide** | Pipeline, promoção de imagem por digest, injeção de versão no build, alarmes, SLO, rotina de backup e de ensaio de restauração |
| **Chamar quando** | A etapa toca infraestrutura, pipeline, release, monitoração ou segredo; é preciso medir RPO, RTO ou atraso de sincronização |
| **Não faz** | Decidir o local da cópia de recuperação de desastre (pendência P7, do dono do produto) |
| **Portão que guarda** | Release sem tag, sem CHANGELOG ou com versão divergente é bloqueado no CI |

---

## 8. Especialidade

### 8.1 `tenaz-conformidade-br`

| | |
|---|---|
| **Dono de** | Os requisitos do módulo CON na interpretação da norma; LGPD e assinatura eletrônica em [9_SEGURANCA_CONFIABILIDADE.md](../9_SEGURANCA_CONFIABILIDADE.md); tratamento fiscal de remessa, conserto e locação |
| **Decide** | Como a norma se traduz em campo, validação e evidência; que pergunta fazer ao apoio externo |
| **Chamar quando** | A etapa toca NR-12, NR-13, LOTO, documentos com vencimento, capacitações, CNH, MOPP, assinatura, LGPD, NF-e ou CFOP |
| **Não faz** | Dar parecer jurídico, contábil ou de SST. Prepara o material e a pergunta; a resposta é da pessoa contratada (P1, P2, P6, P8) |
| **Portão que guarda** | OS de máquina NR-12 não conclui sem os sete campos do item 12.11.2; dado de saúde só com a permissão própria |

### 8.2 `tenaz-integracoes-dados`

| | |
|---|---|
| **Dono de** | Módulos INT e IA; modelo canônico de telemática; busca vetorial isolada por empresa; indicadores avançados |
| **Decide** | Forma do conector, normalização de evento externo, governança da IA, ingestão de IoT |
| **Chamar quando** | A etapa toca telemática, cartão combustível, ERP, importação de planilha, IoT, assistente de IA ou indicador avançado |
| **Não faz** | Deixar a IA gravar sem confirmação humana (princípio 6 da arquitetura) |
| **Portão que guarda** | Toda integração com timeout, retry, disjuntor e idempotência; nenhum dado de empresa cruza a busca vetorial |

---

## 9. Memória

### 9.1 `tenaz-documentacao`

| | |
|---|---|
| **Dono de** | [../etapas/](../etapas/), [../../CHANGELOG.md](../../CHANGELOG.md), a integridade de referência cruzada e o verificador |
| **Decide** | Se uma etapa pode ser fechada do ponto de vista documental |
| **Chamar quando** | Qualquer etapa termina; um documento muda; aparece divergência entre documentos; é preciso regerar PDF |
| **Não faz** | Mudar o conteúdo técnico de um documento sem o agente dono |
| **Portão que guarda** | `python ferramentas/verificar_docs.py` sem erro; mudança de comportamento sem documentação é defeito |

---

## 10. Como manter este catálogo

1. Agente novo ou papel alterado: atualize a definição em [../../.claude/agents/](../../.claude/agents/) **e** este documento, no mesmo commit.
2. Dois agentes disputando o mesmo assunto é sinal de fronteira mal desenhada: `tenaz-arquiteto` decide e registra.
3. Agente que nunca é chamado em três etapas seguidas entra em revisão: ou a especialidade não existe, ou está escondida dentro de outro agente.

---

**Histórico de versões**

- v0.1 (12/11/2026): primeira versão. Dezesseis agentes em seis grupos, correspondência com o time de 15 pessoas do roadmap 0.3, donos, portões e regra de manutenção do catálogo.
