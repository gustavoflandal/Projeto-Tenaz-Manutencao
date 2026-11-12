# 012 — Configuração por segmento

**Data:** 2026-11-12 · **Tipo:** escopo
**Pedido por:** dono do produto · **Versão do produto:** —

## Pedido
O dono do produto apresentou um caso de uso: uma empresa fornece radares de controle de velocidade para várias cidades. Em cada cidade, um prestador instala e mantém os equipamentos com equipe própria. A instalação envolve obra na via, fixação de postes e parte elétrica, e hoje tudo é controlado em planilhas Excel.

A análise mostrou que a Tenaz atendia a manutenção, mas não a implantação (obra), a aferição do Inmetro com lacres, nem o contrato com a prefeitura. O dono do produto pediu que o sistema fosse "flexível e preparado para essas situações", configurável para situações diversas ou determinadas. Depois aprovou a proposta: "Sim, registrar".

## Decisões
- **Configurar blocos conhecidos, não criar fluxos livres** (sumário 3.7, ADR-010). Ficam fora da configuração: os estados canônicos, as fórmulas dos indicadores, as regras do offline, o isolamento e a lista de campos da conexão entre contas. Alternativas descartadas no ADR-010: editor livre de fluxos (low-code), motor de regras com expressões ou scripts, e colunas criadas por empresa.
- **No MVP:**
  - regras de condição e efeito com catálogos fechados (`regra_automacao`, RF-ADM-016);
  - campos extras na OS (`campo_extra_os`, RF-OS-021);
  - fase de ciclo de vida do ativo (`fase_ciclo_vida`, RF-ATV-013). O ativo `em_implantacao` fica com status `inativo` e a interface mostra "Em implantação";
  - tipo de OS `implantacao`, fora dos indicadores de manutenção.
- **Na Fase 2:**
  - pacotes de segmento (RF-ADM-017);
  - módulo **PRJ** de projetos de implantação (RF-PRJ-001 a 007);
  - módulo **CTR** de contratos com clientes (RF-CTR-001 a 005).
- **Nomes canônicos novos** (sumário 3.1 a 3.6):
  - Origem de OS `projeto`. A OS **não** tem origem `regra`: a regra abre solicitação, com origem `regra`.
  - `bloqueio_ativo.origem`, com os valores `manual`, `conformidade` e `regra`.
  - `causa_externa` e `decisao_exclusao`.
  - O estado da etapa é `disponivel`, e não "liberada", para não confundir com a liberação da OS.
  - Numeração `PJ-` e `CT-`.
  - Nove permissões novas.
  - Os campos novos da OS e do ativo nunca são compartilhados na conexão entre contas.
- **Em ativo já em operação, a OS de `melhoria` ou de `implantacao` marcada "retira o ativo de operação" deixa o ativo `em_manutencao`.** É o caso da reinstalação.
- **A causa externa não muda a disponibilidade técnica (Ao).** Ela só desconta na medição do contrato, depois de aprovada.

## Mudanças
| Documento ou arquivo | O que mudou | IDs afetados |
|---|---|---|
| [0_SUMARIO_EXECUTIVO.md](../0_SUMARIO_EXECUTIVO.md) | Versão 0.5. Mudou: proposta de valor 8, público, módulos PRJ e CTR, planos, entidades, estados, numeração, pacotes Go, decisões da 3.6, seção 3.7 nova, escopo por fase | P9, P10, R7 |
| [1_ESPECIFICACAO_REQUISITOS.md](../1_ESPECIFICACAO_REQUISITOS.md) | Versão 0.4. MVP com 129 RF, Fase 2 com 75, total de 220. Pontos em aberto de 35 a 38 | RF-ADM-016, RF-ADM-017, RF-ATV-013, RF-OS-021, RF-PRJ-001 a 007, RF-CTR-001 a 005. Também alterados: RF-ATV-002, RF-ATV-006, RF-SOL-001, RF-OS-001, RF-OS-002, RF-OS-005, RF-CON-005 |
| [2_ARQUITETURA_BANCO_DADOS.md](../2_ARQUITETURA_BANCO_DADOS.md) | Versão 0.4. Tabelas novas de configuração, PRJ e CTR. `plataforma.pacote_segmento` como catálogo global. Isolamento e campos que nunca atravessam a conexão. Pendências S27 a S35 | — |
| [3_ESPECIFICACAO_APIS.md](../3_ESPECIFICACAO_APIS.md) | Versão 0.4. O total de endpoints passou de 482 para 564. Nove eventos novos e 24 códigos de erro. Divergências A-29 a A-36 | — |
| [4_FLUXOS_PROCESSO.md](../4_FLUXOS_PROCESSO.md) | Versão 0.4. Fluxos 29 (regras), 30 (implantação) e 31 (medição), agora 31 fluxos. Disponibilidade contratual na 33.9. Lacunas L-31 a L-39 | RN-ADM-009 a 016, RN-ATV-011 a 015, RN-OS-022 a 026, RN-CON-040 a 044, RN-IND-009 a 011, RN-PRJ-001 a 016, RN-CTR-001 a 015 |
| [5_ARQUITETURA_SOFTWARE.md](../5_ARQUITETURA_SOFTWARE.md) | Versão 0.5. Pacotes `projetos`, `contratos` e `configuracao`, com suas dependências. Seção 5.6 nova. Cache de campos extras no app | ADR-010 |
| [6_ROADMAP.md](../6_ROADMAP.md) | Versão 0.4. Encaixe dos 3 RF novos nas sprints do MVP. Opções de prazo A, B e C. Trilha E na Fase 2. Riscos 26 a 28 | — |
| [7_PADROES_DESIGN.md](../7_PADROES_DESIGN.md) | Versão 0.4. Chip "Em implantação", componentes 6.20 a 6.23, telas 8.29 a 8.34 e 9.10, microcopy das ações novas. Pendências 13.24 a 13.34; parte delas resolvida nesta revisão | — |
| [execucao/4_ETAPAS.md](../execucao/4_ETAPAS.md) | Versão 0.2. Itens novos acrescentados às entregas de E02, E04, E05, E06, E08 e E09. Trilha E da Fase 2 | — |
| [execucao/1_MODELO_EXECUCAO.md](../execucao/1_MODELO_EXECUCAO.md), [execucao/2_EQUIPE_AGENTES.md](../execucao/2_EQUIPE_AGENTES.md), `.claude/agents/tenaz-dominio-manutencao.md`, `.claude/skills/tenaz-indicador/SKILL.md` | Contagens e números de seção atualizados: 220 RF, 31 fluxos, 361 RN, 10 ADRs, indicadores na seção 33 do documento 4 | — |
| [CLAUDE.md](../../CLAUDE.md) | Decisão "configuração por segmento" na seção 4 | — |
| [CHANGELOG.md](../../CHANGELOG.md) | Entrada em `[Não lançado]` → Adicionado | — |

## Verificação
- Seis agentes especialistas atualizaram os documentos 2 a 7 em paralelo, cada um no documento que é dono do assunto.
- No meio do trabalho, o coordenador enviou a cada agente os ajustes canônicos (`disponivel`, origem de OS sem `regra`, `decisao_exclusao`, permissões).
- Em seguida veio a revisão cruzada. Foi feita busca por `liberada`, origem `regra` e pela lista de OS que deixam o ativo `em_manutencao`. As divergências de nome apontadas pelos agentes foram corrigidas no sumário e nos documentos 1 e 7.
- Rodou-se `python ferramentas/verificar_docs.py` e os PDFs dos documentos alterados foram gerados de novo.

## Pendências
**Decisões do dono do produto:**
1. **Prazo e time do MVP** ([6_ROADMAP.md](../6_ROADMAP.md), seção 2.2). O acréscimo é de cerca de 9 a 10 sprints-pessoa:
   - A: absorver nas 13 sprints, com risco;
   - B: somar 1 sprint;
   - C: somar 1 backend de S4 a S10.
2. **Time e prazo da Fase 2:**
   - 1: reforço de 1 backend e 1 frontend;
   - 2: estender até o mês 13 ou 14;
   - 3: mover requisitos de fase.
3. **P9 e P10** ([0_SUMARIO_EXECUTIVO.md](../0_SUMARIO_EXECUTIVO.md), seção 6):
   - P9: regras metrológicas do radar;
   - P10: acesso do cliente final.

**Lacunas de produto com subcontratada conectada** ([1_ESPECIFICACAO_REQUISITOS.md](../1_ESPECIFICACAO_REQUISITOS.md), seção 8.3):
- Ponto 38: campo extra que dispara regra no cliente, como "Rompeu lacre?", não chega quando a OS é executada por subcontratada conectada. Três opções registradas.
- A causa externa marcada pela subcontratada conectada também não chega ao cliente (documento 2, S30; documento 3, A-32).

**Detalhes registrados nos documentos donos, para a próxima revisão:**
- Documento 4, L-31 a L-39:
  - se a regra, ao ser ativada, dispara para fatos que já existem;
  - a permissão para levantar à mão o bloqueio de regra;
  - comparação por dia da emissão do documento;
  - causa externa em OS já encerrada;
  - critérios de arredondamento e de proporção da medição.
- Documento 3: A-24 (rota da ação `devolver`), A-30, A-31, A-35 e A-36.
- Documento 2: S27 a S35, com os nomes propostos de tabelas filhas: `item_medicao_contrato`, `exclusao_parada_contrato`, `ativo_projeto`, `lancamento_custo_projeto`, `medicao_etapa_projeto`, `dependencia_etapa_modelo`.

**Nomes propostos pelos agentes que ainda não estão no sumário:**
- A medição de etapas tem dois nomes: `medicao_etapa` no documento 3 e `medicao_etapa_projeto` no documento 2. É preciso unificar.

**Sobras da troca de nomes da etapa 011, apontadas pelo documento 5 (seção 16):**
- `tenaz.empresas_ativas()`;
- `tenaz.conexao_atual()`;
- `app.escopo_tipo`.

**Lacunas antigas do documento 4 (L-21 a L-30) já decididas no sumário 3.6, mas ainda tratadas como abertas.**
