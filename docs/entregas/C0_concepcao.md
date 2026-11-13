# C0 — Concepção

**Etapa:** concepção, anterior às etapas de construção E00 a E12 · **Consolidada em:** 2026-11-12
**Commits:** 8e4b3fc (marca, estudo de mercado e documentação), ed51d08 (prompt inicial, regras de documentação e plano de execução), 1dd31c0 (equipe de agentes, skills e etapas E02 a E12), 2011feb (troca de repositório), 294d7f1 (configuração por segmento, junto com a E00)
**Registros das etapas:** [001](../etapas/001_2026-11-12_nome-e-identidade-visual.md) a [010](../etapas/010_2026-11-12_troca-de-repositorio.md) e [012](../etapas/012_2026-11-12_configuracao-por-segmento.md)

## O que foi entregue

O produto ficou definido antes de qualquer código: nome, marca, mercado, escopo, especificação técnica e o método de trabalho.

| Entrega | Onde está | Registro |
|---|---|---|
| Nome Tenaz e identidade visual: símbolo, logotipo, cores Aço e Sinal, tipografia, brandbook, regras de contraste | [marca/](../../marca/README.md) | [001](../etapas/001_2026-11-12_nome-e-identidade-visual.md) |
| Estudo de mercado e benchmarking de sistemas de manutenção e frota | [pesquisa/](../pesquisa/) | [002](../etapas/002_2026-11-12_estudo-de-mercado.md) |
| Documentação base: sumário executivo, requisitos, banco de dados, APIs, fluxos e regras, arquitetura, roadmap, padrões de design, segurança | [docs/](../README.md) | [003](../etapas/003_2026-11-12_documentacao-base.md) |
| Módulos de terceiros (TER) e de locação (LOC) e estoque estendido | Sumário e documentos 1 a 9 | [004](../etapas/004_2026-11-12_terceiros-estoque-locacao.md) |
| Versão visível na tela e histórico de versões (RNF-039, RNF-040) | Documento 1 e [CHANGELOG](../../CHANGELOG.md) | [005](../etapas/005_2026-11-12_versao-e-changelog.md) |
| Módulo de gestão de serviços do prestador (SRV) e conexão entre contas | Sumário 3.5 e documentos 1 a 9 | [006](../etapas/006_2026-11-12_gestao-servicos-prestador.md) |
| Primeiro commit no GitHub | Repositório | [007](../etapas/007_2026-11-12_primeiro-commit.md) |
| Prompt inicial do projeto e regras obrigatórias de documentação, com o verificador | [CLAUDE.md](../../CLAUDE.md), [ferramentas/verificar_docs.py](../../ferramentas/verificar_docs.py) | [008](../etapas/008_2026-11-12_prompt-e-regras-de-documentacao.md) |
| Plano de execução: modelo de trabalho, 16 agentes, 12 skills, etapas E00 a E12 e seus prompts | [execucao/](../execucao/README.md), [prompts/](../prompts/README.md) | [009](../etapas/009_2026-11-12_equipe-de-agentes-e-skills.md) |
| Repositório definitivo `Projeto-Tenaz-Manutencao` | GitHub | [010](../etapas/010_2026-11-12_troca-de-repositorio.md) |
| Configuração por segmento (catálogos fechados), módulos PRJ e CTR na Fase 2, ADR-010 | Sumário 3.7 e documentos 1 a 7 | [012](../etapas/012_2026-11-12_configuracao-por-segmento.md) |

## Situação

A concepção não tem portão de saída de software. Ela é consolidada porque cada documento passou pelo verificador sem erro, está em `main` e foi a base aceita pelo dono do produto para iniciar a E00.

Os documentos continuam vivos: cada etapa de construção os atualiza, e o histórico de versões de cada um diz o que mudou e quando. Esta entrega registra o ponto de partida, não um retrato congelado.

## O que ficou para depois

As decisões que dependem do dono do produto ou de terceiros estão no [sumário](../0_SUMARIO_EXECUTIVO.md), seção 6: P1 a P11. A decisão de prazo e de time do MVP com a configuração por segmento está no [roadmap](../6_ROADMAP.md), seção 2.2, e no [registro 012](../etapas/012_2026-11-12_configuracao-por-segmento.md).

## Notas posteriores

Nenhuma até agora.
