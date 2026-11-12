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

O manual de operação, o roteiro de implantação, a homologação e o treinamento entram depois, quando houver telas implementadas.
