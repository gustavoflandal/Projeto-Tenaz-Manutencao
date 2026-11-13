# Plano de execução da Tenaz

**Versão:** 0.1 · **Data:** 12/11/2026 · **Status:** Em elaboração
**Produto:** Tenaz · Gestão de Manutenção e Frota

Este conjunto de documentos diz **como** a Tenaz sai da concepção e vira software. O [roadmap](../6_ROADMAP.md) diz o que se entrega e em que ordem, pressupondo um time de 15 pessoas em dedicação integral. Este plano trata da execução real: o dono do produto conduzindo uma **equipe de agentes especializados**, com os portões de qualidade do roadmap mantidos intactos.

O roadmap continua sendo a fonte da verdade de escopo, sequência, dependências e marcos. Este plano não altera nenhum requisito nem nenhuma data do roadmap; organiza o trabalho de quem vai executá-lo.

## Documentos

| # | Documento | Conteúdo |
|---|---|---|
| 1 | [1_MODELO_EXECUCAO.md](1_MODELO_EXECUCAO.md) | Como se trabalha: papéis, ciclo de uma etapa, orquestração dos agentes, portões que bloqueiam, o que não se delega |
| 2 | [2_EQUIPE_AGENTES.md](2_EQUIPE_AGENTES.md) | Os 18 agentes: especialidade, documentos donos, quando chamar, quem revisa quem |
| 3 | [3_SKILLS.md](3_SKILLS.md) | As 12 skills: o procedimento obrigatório de cada tarefa recorrente |
| 4 | [4_ETAPAS.md](4_ETAPAS.md) | As etapas E00 a E12 do MVP e os blocos das fases 2 e 3, com agentes, prompt e evidência de cada uma |

## Pastas de apoio

| Pasta | Para que serve |
|---|---|
| [../prompts/](../prompts/) | O prompt de cada etapa, versionado: o que se pediu, com que contexto e com que critério de aceite |
| [../screenshots/](../screenshots/) | Evidência visual de cada tela entregue, usada no aceite, na documentação e no manual de operação |
| [../etapas/](../etapas/) | Registro do que de fato aconteceu em cada etapa (regra de [CLAUDE.md](../../CLAUDE.md), seção 3.1) |
| [../../.claude/agents/](../../.claude/agents/) | Definição executável dos agentes |
| [../../.claude/skills/](../../.claude/skills/) | Definição executável das skills |

## Como usar, na prática

1. Leia [1_MODELO_EXECUCAO.md](1_MODELO_EXECUCAO.md) uma vez. Ele define o ciclo que vale para toda etapa.
2. Abra a etapa da vez em [4_ETAPAS.md](4_ETAPAS.md) e o prompt dela em [../prompts/](../prompts/).
3. Trabalhe com os agentes indicados na etapa. A skill `tenaz-contexto` orienta qualquer sessão nova.
4. Só feche a etapa quando o portão dela estiver verde e o registro em [../etapas/](../etapas/) estiver escrito.

---

**Histórico de versões**

- v0.1 (12/11/2026): primeira versão. Modelo de execução com equipe de agentes, catálogo de 16 agentes e 12 skills, etapas E00 a E12 do MVP e pastas de prompts e de evidência visual.
