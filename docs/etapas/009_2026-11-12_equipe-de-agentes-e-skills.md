# 009 — Plano de execução, equipe de agentes e skills

**Data:** 2026-11-12 · **Tipo:** decisão
**Pedido por:** dono do produto · **Versão do produto:** —

## Pedido
Ler o projeto, criar um plano de execução, criar uma subpasta em `docs` para ele, mais a pasta de screenshots e uma para os prompts. Depois, à pergunta sobre quem executaria o projeto: criar uma equipe de agentes super especializados e um conjunto de skills que ajudem nesse projeto, salvos na pasta `.claude` do projeto.

## Decisões

- **Quem executa.** O [roadmap](../6_ROADMAP.md) pressupõe 15 pessoas em dedicação integral. A execução real é o dono do produto conduzindo uma equipe de agentes. O plano registra o que isso muda (capacidade, calendário, gargalo, risco dominante) e o que não muda (sequência, dependências críticas e critério de pronto).
- **Plano em `docs/execucao/`**, em quatro documentos: modelo de execução, equipe de agentes, catálogo de skills e etapas. O roadmap continua a fonte da verdade de escopo e sequência; o plano trata só do método.
- **Dezesseis agentes**, em `.claude/agents/`, correspondentes aos papéis do roadmap mais três que o projeto exige e que num time humano ficariam distribuídos: `tenaz-dominio-manutencao` (304 regras de negócio e as fórmulas de indicador), `tenaz-conformidade-br` (conformidade brasileira, que é a proposta de valor 4 do sumário) e `tenaz-documentacao` (a disciplina documental da seção 3 de [CLAUDE.md](../../CLAUDE.md)). Apoio externo — advogado, consultor de SST, contador, pentest, auditor — **não** vira agente.
- **`tenaz-isolamento` com poder de veto amplo**, porque R2, R4 e R5 do [sumário](../0_SUMARIO_EXECUTIVO.md) são os únicos riscos classificados como críticos. Quem cria política de isolamento não aprova a própria política.
- **Doze skills**, em `.claude/skills/`, cada uma o procedimento de uma tarefa recorrente, terminando sempre em verificação. Skill é procedimento; o conteúdo técnico fica nos documentos de `docs/`, e a skill aponta para eles.
- **Etapas E00 a E12 correspondentes às sprints S0 a S12**, sem duração prevista: terminam quando o portão fecha e o dono do produto aceita. Os marcos M0 a M5 seguem como referência de prazo para o que depende de terceiros.
- **Seis portões que bloqueiam** o fechamento de etapa: isolamento entre empresas, isolamento de acesso externo, isolamento entre contas conectadas, offline sem perda, acessibilidade com identidade e integridade documental. Não existe fechamento com ressalva.
- **Prompt de cada etapa versionado** em `docs/prompts/`, escrito antes do trabalho e nunca reescrito depois: é o que permite separar, meses depois, "a decisão estava errada" de "o pedido estava ambíguo".
- **Evidência visual em `docs/screenshots/`**, uma pasta por etapa, com convenção de nome, largura fixa, dois temas, massa de dados plausível e proibição de dado pessoal real.
- Alternativa descartada: redividir o escopo do roadmap em etapas próprias. Manter a correspondência um a um com as sprints preserva as 19 dependências críticas e a rastreabilidade dos requisitos.

## Mudanças

| Documento ou arquivo | O que mudou | IDs afetados |
|---|---|---|
| [docs/execucao/](./) | Pasta criada, com README, [1_MODELO_EXECUCAO.md](../execucao/1_MODELO_EXECUCAO.md), [2_EQUIPE_AGENTES.md](../execucao/2_EQUIPE_AGENTES.md), [3_SKILLS.md](../execucao/3_SKILLS.md) e [4_ETAPAS.md](../execucao/4_ETAPAS.md) | Cita RNF-005, RNF-007, RNF-010, RNF-013, RNF-025, RNF-037, RNF-038, RNF-039, RNF-040, RNF-041, RNF-042 |
| [docs/prompts/](../prompts/) | Pasta criada, com README, modelo e os prompts das etapas E00 a E12 | Requisitos do MVP, por etapa |
| [docs/screenshots/](../screenshots/) | Pasta criada, com a convenção de captura, nome e guarda | RNF-039 |
| [.claude/agents/](../../.claude/agents/) | 16 definições de agente criadas | — |
| [.claude/skills/](../../.claude/skills/) | 12 definições de skill criadas | — |
| [CLAUDE.md](../../CLAUDE.md) | Seção 2 com o plano de execução entre as fontes; seção 6 com a skill de contexto, os agentes e o ciclo de sete passos; seção 7 com as pastas novas | — |
| [docs/README.md](../README.md) | Seção "Execução" com as três pastas novas | — |
| [CHANGELOG.md](../../CHANGELOG.md) | Quatro entradas em Adicionado | — |

## Verificação

- `python ferramentas/verificar_docs.py`: **0 erros, 0 avisos**, com os 38 arquivos novos no escopo da verificação. O verificador cobre os `.md` de `.claude/`, então links relativos e IDs citados nas definições de agente e de skill foram conferidos como os dos demais documentos.
- Cobertura de requisitos das etapas E00 a E12 conferida contra a seção 4 do [roadmap](../6_ROADMAP.md), sprint por sprint.
- Identificadores citados conferidos contra o documento dono: 206 RF (ADM 1–15, ATV 1–12, MED 1–8, SOL 1–6, OS 1–20, PRE 1–12, CHK 1–6, EST 1–19, TER 1–9, SRV 1–21, LOC 1–9, FRO 1–12, CON 1–16, IND 1–12, INT 1–12, IA 1–7, MOB 1–10), 42 RNF e 9 ADR.
- Nenhum PDF a regerar: os documentos 0 a 9 não foram alterados.

## Pendências

- As pastas `docs/execucao/`, `docs/prompts/` e `docs/screenshots/` não têm PDF, por serem material de processo e não especificação de produto. Se o dono do produto quiser PDF do plano, entra como etapa nova.
- O plano não altera nenhuma pendência de produto. **P1** a **P8** seguem como estão na seção 6 do [sumário](../0_SUMARIO_EXECUTIVO.md), e o plano registra em que etapa cada uma trava o trabalho.
- A pendência aberta na etapa [008](008_2026-11-12_prompt-e-regras-de-documentacao.md) — verificador sem execução no CI — entra formalmente no portão de saída da etapa E00, em [docs/prompts/E00_fundacao-tecnica.md](../prompts/E00_fundacao-tecnica.md).
- O catálogo de agentes prevê revisão: agente não chamado em três etapas seguidas entra em análise, porque ou a especialidade não existe ou está escondida dentro de outro agente.
