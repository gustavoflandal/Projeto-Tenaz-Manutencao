---
name: tenaz-contexto
description: Use no começo de QUALQUER tarefa no repositório da Tenaz, inclusive em tarefas que pareçam triviais - orienta sobre as fontes da verdade, o documento dono de cada assunto, o vocabulário canônico que não se improvisa, o que já está decidido, o que é pendência do dono do produto e qual agente chamar. Leia antes de responder, explorar arquivos ou propor qualquer coisa.
---

# Contexto da Tenaz

Esta skill existe porque o custo de errar um nome neste projeto é alto: um termo fora do vocabulário canônico se propaga por documentação, banco, API, web e app ao mesmo tempo, e o verificador só pega parte disso.

## Passo 1 — Leia, nesta ordem

1. [CLAUDE.md](../../../CLAUDE.md) inteiro. É o prompt inicial do projeto e vale para toda sessão.
2. [docs/0_SUMARIO_EXECUTIVO.md](../../../docs/0_SUMARIO_EXECUTIVO.md), seção 3 (vocabulário, estados, convenções) e seção 6 (pendências).
3. O documento dono do assunto da sua tarefa, pela tabela abaixo.

## Passo 2 — Ordem das fontes da verdade

Em caso de conflito, vale esta ordem, sem exceção:

1. Pedido explícito do dono do produto na sessão.
2. O sumário executivo.
3. O documento dono do assunto.

## Passo 3 — O documento dono de cada assunto

| Assunto | Documento dono |
|---|---|
| Requisitos funcionais e não funcionais | [docs/1_ESPECIFICACAO_REQUISITOS.md](../../../docs/1_ESPECIFICACAO_REQUISITOS.md) |
| Modelo de dados, SQL, isolamento, sincronização | [docs/2_ARQUITETURA_BANCO_DADOS.md](../../../docs/2_ARQUITETURA_BANCO_DADOS.md) |
| Endpoints, eventos de webhook, códigos de erro | [docs/3_ESPECIFICACAO_APIS.md](../../../docs/3_ESPECIFICACAO_APIS.md) |
| Fluxos, regras de negócio, cálculo de indicadores | [docs/4_FLUXOS_PROCESSO.md](../../../docs/4_FLUXOS_PROCESSO.md) |
| Arquitetura, ADRs, versionamento, testes | [docs/5_ARQUITETURA_SOFTWARE.md](../../../docs/5_ARQUITETURA_SOFTWARE.md) |
| Prazos, sprints, equipe, riscos de cronograma | [docs/6_ROADMAP.md](../../../docs/6_ROADMAP.md) |
| Interface, tokens, componentes, contrastes | [docs/7_PADROES_DESIGN.md](../../../docs/7_PADROES_DESIGN.md) e [marca/README.md](../../../marca/README.md) |
| Segurança, LGPD, SLA, recuperação de desastre | [docs/9_SEGURANCA_CONFIABILIDADE.md](../../../docs/9_SEGURANCA_CONFIABILIDADE.md) |
| Como o trabalho acontece, agentes, etapas | [docs/execucao/](../../../docs/execucao/) |
| Histórico de versões do produto | [CHANGELOG.md](../../../CHANGELOG.md) |
| Histórico de etapas do projeto | [docs/etapas/](../../../docs/etapas/) |

Um identificador (`RF-`, `RNF-`, `RN-`, `ADR-`) só **nasce** no documento dono. Os outros apenas o citam, por ID e link relativo.

## Passo 4 — Vocabulário que não se improvisa

Erros recorrentes, todos já decididos no sumário:

| Erro comum | Correto |
|---|---|
| "ativo em atenção" como status | `operando`, `em_manutencao`, `parado`, `inativo`. "Atenção" é condição **calculada**, não status gravado |
| MTTR contado até o encerramento | KPIs de tempo usam `data_inicio` e `data_conclusao` e as paradas, nunca `data_encerramento` |
| "reconhecimento" como estado da OS | É o campo `data_reconhecimento`, gravado uma única vez |
| MOPP como documento | MOPP é **sempre** capacitação; CNH é documento da pessoa |
| "apto para despacho" como status | Verificação calculada, separada do status |
| LOTO e bloqueio de ativo como a mesma coisa | Bloqueio de ativo deixa o ativo `parado`; o registro LOTO é de bloqueio de energia de uma intervenção |
| Rota no singular | `/v1/<recurso-no-plural>` em kebab-case. `GET /v1/versao` é a única exceção aceita |
| Nome em inglês no domínio | Português sem acento: `ordem_servico.data_abertura`, `OrdemServico`, `Concluir` |

Precisa de um nome que não existe? Ele entra **primeiro no sumário**, por decisão do dono do produto. Não invente no código.

## Passo 5 — O que não se decide sem o dono do produto

Pendências **P1** (ICP-Brasil na NR-1), **P2** (validade de CIV e CIPP), **P3** (API de multas), **P4** (preços e prioridades), **P5** (marca no INPI), **P6** (prazos da NR-13), **P7** (local da cópia de recuperação de desastre), **P8** (CFOPs). Também: preço, mudança de escopo ou de fase, nome canônico novo, aplicação da linha de corte, liberação de portão com teste vermelho, commit e push.

Apareceu uma dessas? Apresente as opções e pare.

## Passo 6 — Qual agente chamar

O catálogo está em [docs/execucao/2_EQUIPE_AGENTES.md](../../../docs/execucao/2_EQUIPE_AGENTES.md). Atalho:

- Decisão estrutural → `tenaz-arquiteto`
- Escopo, requisito, aceite → `tenaz-produto`
- Comportamento do domínio, regra, indicador → `tenaz-dominio-manutencao`
- Tabela, migração, consulta → `tenaz-banco-dados`
- Qualquer coisa que toque isolamento → `tenaz-isolamento`
- Regra em Go → `tenaz-backend-go`
- Rota, erro, webhook → `tenaz-api-contrato`
- Conexão entre contas → `tenaz-conexao-contas`
- Tela, token, contraste → `tenaz-design-sistema`
- Web → `tenaz-frontend-web`
- App e offline → `tenaz-mobile-offline`
- Teste → `tenaz-qa-automacao`
- Infra, pipeline, release → `tenaz-plataforma-sre`
- NR, LGPD, fiscal → `tenaz-conformidade-br`
- Telemática, ERP, IA → `tenaz-integracoes-dados`
- Fechar etapa, CHANGELOG, PDF → `tenaz-documentacao`

## Passo 7 — Antes de dizer que terminou

```bash
python ferramentas/verificar_docs.py
```

Zero erro. E confira as quatro atualizações obrigatórias de [CLAUDE.md](../../../CLAUDE.md), seção 3.2: documento afetado, CHANGELOG, registro da etapa e sumário quando couber.

## Estado atual do projeto

Concepção concluída: marca, estudo de mercado e documentação. **Ainda não há código.** A primeira etapa de construção é E00, em [docs/execucao/4_ETAPAS.md](../../../docs/execucao/4_ETAPAS.md).
