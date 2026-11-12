# E00 — Fundação técnica

**Data do pedido:** 2026-11-12 · **Sprint do roadmap:** S0 · **Marco:** M0
**Pedido por:** dono do produto

## Objetivo

Deixar o repositório pronto para receber requisito: pipeline verde, infraestrutura criada por código na região São Paulo, isolamento entre empresas verificado automaticamente desde o primeiro build e versão rastreável em tudo que for publicado. Nenhum requisito funcional entra nesta etapa.

## Entrega

| Requisito | Observação de escopo nesta etapa |
|---|---|
| RNF-013 | Base: esquema com `empresa_id`, RLS forçado e a suíte de vazamento rodando no CI, ainda com poucas tabelas |
| RNF-016 | Base de autenticação e gestão de segredos |
| RNF-020 | Base de observabilidade: log estruturado, trace e métrica |
| RNF-032 | Base de proteção de dados em trânsito e em repouso |
| RNF-039 | Esteira completa: versão gerada da tag, injeção no build do backend, da web e do app, `GET /v1/versao`, bloqueio de release em divergência |
| RNF-040 | Verificação no CI que recusa alteração de código de produto sem entrada no `CHANGELOG.md` |

Também entra, sem requisito funcional: estrutura do monorepo conforme [5_ARQUITETURA_SOFTWARE.md](../5_ARQUITETURA_SOFTWARE.md), seção 7; esquema da plataforma separado do RLS por empresa; projeto do app com banco local e prova de conceito da fila de comandos; rascunho de ADR-009; execução do verificador de documentação no CI.

## Contexto obrigatório

| Documento | Seção | Por quê |
|---|---|---|
| [5_ARQUITETURA_SOFTWARE.md](../5_ARQUITETURA_SOFTWARE.md) | 2, 4, 5, 7, 9 | Stack, contêineres, monólito modular, estrutura do repositório, testes e pipeline |
| [5_ARQUITETURA_SOFTWARE.md](../5_ARQUITETURA_SOFTWARE.md) | 9.5 | Versionamento e release, regra por regra |
| [2_ARQUITETURA_BANCO_DADOS.md](../2_ARQUITETURA_BANCO_DADOS.md) | 1, 2, 4 | Convenções, multiempresa e isolamento, base comum do esquema |
| [9_SEGURANCA_CONFIABILIDADE.md](../9_SEGURANCA_CONFIABILIDADE.md) | 2, 3, 11 | Alta disponibilidade, backup, checklist de segurança por release |
| [CLAUDE.md](../../CLAUDE.md) | inteiro | Regras do projeto |

## Agentes previstos

Conduz `tenaz-plataforma-sre`. Apoiam `tenaz-arquiteto` (estrutura e ADR-009), `tenaz-banco-dados` (esquema base e primeira migração), `tenaz-isolamento` (suíte de vazamento no CI), `tenaz-mobile-offline` (projeto do app e fila de comandos), `tenaz-documentacao`.

## Portão de saída

| Verificação | Como se comprova |
|---|---|
| CI verde de ponta a ponta | Execução do pipeline no repositório |
| Infraestrutura criada por código | Plano e aplicação da infraestrutura em ambiente de desenvolvimento |
| Teste de vazamento entre empresas no CI | Saída da suíte, com a verificação de cobertura de RLS quebrando o build se faltar política |
| Esteira de versão ativa | Release de teste bloqueado quando tag, CHANGELOG e versão injetada divergem; `GET /v1/versao` respondendo |
| Build do app instalável | Instalação em um aparelho, com o selo de versão na tela |
| Verificador de documentação no CI | `python ferramentas/verificar_docs.py` como etapa obrigatória do pipeline |

## Evidência esperada

`docs/screenshots/E00_fundacao-tecnica/`: pipeline verde, app instalado com o selo de versão, saída de `GET /v1/versao`.

## Fora do escopo

Qualquer requisito funcional. Ambiente de produção. SSO (E10). Suíte de três contas, que entra em E04 com a fundação da conexão. Telas de negócio.

## Decisões pendentes do dono do produto

- **P7** (local da cópia de recuperação de desastre) afeta a infraestrutura de backup. Prazo do roadmap: semana 20. Nesta etapa, a cópia fica dentro do Brasil e a decisão é registrada como pendência na infraestrutura.
- Contas de nuvem e contas de desenvolvedor das lojas precisam existir antes da etapa começar.

## Ajustes durante a etapa

**2026-11-12 — Trabalho segue em ambiente de desenvolvimento.** Por decisão do dono do produto, os dois itens do portão que dependem de conta externa saem desta etapa e passam para a fase de produção:

| Item do portão | Novo destino |
|---|---|
| Infraestrutura em São Paulo criada por código | Fase de produção. `infra/tofu/` não é escrito agora; o ambiente é o compose local de `infra/docker/` |
| Build interno do app instalável | Fase de produção, junto com as contas de desenvolvedor das lojas |

Motivo: não há conta de nuvem nem conta de loja, e a pendência **P7** (local da cópia de recuperação de desastre) segue aberta — provisionar produção antes dela exigiria refazer a configuração de backup e republicar o RTO.

Os demais itens do portão continuam exigidos e foram cumpridos em ambiente de desenvolvimento. O que **não** foi afrouxado: a esteira de versão, a suíte de isolamento, o contrato OpenAPI e o verificador de documentação valem igual. A etapa fecha com esses dois itens registrados como pendência da fase de produção, não como dívida silenciosa.
