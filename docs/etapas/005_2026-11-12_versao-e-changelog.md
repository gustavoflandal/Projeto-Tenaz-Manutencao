# 005 — Versão visível e histórico de versões

**Data:** 2026-11-12 · **Tipo:** escopo
**Pedido por:** dono do produto · **Versão do produto:** —

## Pedido
1. Obrigar o controle de versão, mostrando sempre no canto inferior direito da tela principal a versão atual no formato usual.
2. Manter um documento sempre atualizado com todas as implementações e correções de cada versão.

## Decisões
- SemVer 2.0.0 (`MAJOR.MINOR.PATCH`), vindo da tag git `vX.Y.Z`, nunca digitado à mão; selo `v1.4.2` no canto inferior direito da tela principal e do login, web e app; diálogo "Sobre"; endpoint `GET /v1/versao`.
- `CHANGELOG.md` na raiz, Keep a Changelog em português; atualizado no mesmo pull request; CI rejeita mudança de código sem registro; "Novidades desta versão" no "Sobre".

## Mudanças
| Documento ou arquivo | O que mudou | IDs afetados |
|---|---|---|
| [1_ESPECIFICACAO_REQUISITOS.md](../1_ESPECIFICACAO_REQUISITOS.md) | Requisitos novos | RNF-039, RNF-040 |
| [3_ESPECIFICACAO_APIS.md](../3_ESPECIFICACAO_APIS.md) | `GET /v1/versao`, cabeçalho `Tenaz-Versao` | — |
| [5_ARQUITETURA_SOFTWARE.md](../5_ARQUITETURA_SOFTWARE.md) | Seção 9.5 de versionamento e release | — |
| [6_ROADMAP.md](../6_ROADMAP.md) | Esteira em S0, selo em S1, critério de release | — |
| [7_PADROES_DESIGN.md](../7_PADROES_DESIGN.md) | Componente "Selo de versão" e diálogo "Sobre" | — |
| [CHANGELOG.md](../../CHANGELOG.md) | Criado | — |

## Verificação
Requisito citado de forma igual nos documentos 1, 3, 5, 6, 7 e 9.

## Pendências
Escolha da ferramenta de release (release-please ou git-cliff) — documento 5.
