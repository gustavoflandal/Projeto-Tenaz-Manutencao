# Documentação do projeto Tenaz

Este conjunto de documentos especifica a **Tenaz**, sistema SaaS de gestão de manutenção de equipamentos e frota. Cada documento existe em `.md`, e uma versão `.pdf` pode ser gerada com `python ferramentas/md2pdf.py docs/<arquivo>.md`.

## Pesquisa

| Documento | Conteúdo |
|---|---|
| [pesquisa/reports/](pesquisa/reports/) | Estudo de mercado e benchmarking: sistemas de referência, funcionalidades, normas, indicadores, arquitetura e oportunidades |
| [pesquisa/research_notes/](pesquisa/research_notes/) | Notas de pesquisa por tema, com as fontes |

## Documentos do projeto

| # | Documento | Conteúdo | Situação |
|---|---|---|---|
| 0 | [0_SUMARIO_EXECUTIVO.md](0_SUMARIO_EXECUTIVO.md) | Visão do produto, público, módulos, planos, vocabulário canônico, estados, pendências | Em elaboração |
| 1 | [1_ESPECIFICACAO_REQUISITOS.md](1_ESPECIFICACAO_REQUISITOS.md) | Requisitos funcionais por módulo, requisitos não funcionais com metas, prioridades | Em elaboração |
| 2 | [2_ARQUITETURA_BANCO_DADOS.md](2_ARQUITETURA_BANCO_DADOS.md) | Modelo de dados multiempresa, DDL, isolamento por empresa, prestador e locatário, sincronização entre contas | Em elaboração |
| 3 | [3_ESPECIFICACAO_APIS.md](3_ESPECIFICACAO_APIS.md) | API REST pública, webhooks, autenticação, limites, endpoints por módulo | Em elaboração |
| 4 | [4_FLUXOS_PROCESSO.md](4_FLUXOS_PROCESSO.md) | Fluxos de negócio, estados, catálogo de regras de negócio, cálculo dos KPIs | Em elaboração |
| 5 | [5_ARQUITETURA_SOFTWARE.md](5_ARQUITETURA_SOFTWARE.md) | Stack, módulos, offline, telemetria, IA, versionamento, decisões de arquitetura (ADR) | Em elaboração |
| 6 | [6_ROADMAP.md](6_ROADMAP.md) | MVP e fases seguintes, sprints, marcos, equipe, riscos | Em elaboração |
| 7 | [7_PADROES_DESIGN.md](7_PADROES_DESIGN.md) | Tokens, componentes, contrastes, telas, app de campo | Em elaboração |
| 8 | Manual de operação | Reservado; escrito quando houver telas | — |
| 9 | [9_SEGURANCA_CONFIABILIDADE.md](9_SEGURANCA_CONFIABILIDADE.md) | Segurança, LGPD, SLA, backup, recuperação de desastre, observabilidade | Em elaboração |
| — | [../CHANGELOG.md](../CHANGELOG.md) | Histórico de versões: todas as implementações e correções de cada versão | Sempre atualizado |
| — | [etapas/](etapas/README.md) | Registro de todas as etapas do projeto: pedido, decisões, mudanças, verificação, pendências | Sempre atualizado |
| — | [entregas/](entregas/README.md) | Entregas consolidadas: cada etapa aceita e em `main`, com o que ficou pronto, o portão item a item e o que foi para depois | A cada etapa aceita |
| — | [pdf/](pdf/) | Versões em PDF dos documentos, geradas do Markdown | Regeradas a cada mudança |

## Execução

Como o projeto sai da concepção e vira software: o modelo de trabalho, a equipe de agentes especializados, as skills e as etapas.

| Documento | Conteúdo |
|---|---|
| [execucao/](execucao/README.md) | Plano de execução: modelo de trabalho, equipe de 16 agentes, 12 skills, etapas E00 a E12 do MVP |
| [prompts/](prompts/README.md) | O prompt de cada etapa, escrito antes do trabalho e versionado |
| [screenshots/](screenshots/README.md) | Evidência visual de cada tela entregue, usada no aceite e na documentação |

Os agentes e as skills estão definidos em [../.claude/agents/](../.claude/agents/) e [../.claude/skills/](../.claude/skills/).

## Regras de documentação

As regras completas estão em [CLAUDE.md](../CLAUDE.md), seção 3. Em resumo:

1. **Toda etapa é registrada** em [etapas/](etapas/README.md), um arquivo por etapa, com índice.
2. **Nenhuma correção ou implementação sem documentação atualizada:** no mesmo commit entram os documentos afetados (com versão e histórico), o [CHANGELOG](../CHANGELOG.md), o registro da etapa e, quando couber, o sumário executivo.
3. **Integridade e referência cruzada:** nomes canônicos do sumário; IDs criados só no documento dono e nunca renumerados; citações por ID e link relativo; divergência não resolvida vai para a seção de pendências.
4. **Verificação obrigatória** antes de concluir uma etapa:

```bash
python ferramentas/verificar_docs.py
```

5. **PDFs** regerados do Markdown para todo documento alterado.

O manual de operação, o roteiro de implantação, a homologação e o treinamento entram depois, quando houver telas implementadas.
