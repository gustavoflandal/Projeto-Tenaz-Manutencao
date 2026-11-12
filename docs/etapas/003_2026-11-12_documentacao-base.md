# 003 — Documentação base do projeto

**Data:** 2026-11-12 · **Tipo:** escopo
**Pedido por:** dono do produto · **Versão do produto:** —

## Pedido
Desenvolver a documentação do projeto após a validação da marca e do estudo. Pergunta do dono do produto sobre a melhor stack.

## Decisões
- Estrutura numerada no padrão do PCP-Lev (documentos 0 a 9; o 8, manual de operação, reservado).
- O sumário executivo é a referência comum: módulos com código, vocabulário canônico, estados, convenções de ID.
- **Stack:** Go no backend, TypeScript no web (React) e no app (React Native/Expo), PostgreSQL com RLS (ADR-001). Alternativas descartadas: Node no backend, Java/Kotlin, Rails, Python, Rust.
- Sumário 0.2 consolidou as inconsistências apontadas na primeira redação: entidade `pessoa`, perfis `comprador` e `plh`, prioridade `urgente`, `data_reconhecimento`, transições completas da OS, regras do status do ativo, definição de técnico e veículo cobrados, conteúdo de cada plano.
- Correções de acessibilidade na marca (chip de status, contorno de foco, Aço 400).

## Mudanças
| Documento ou arquivo | O que mudou | IDs afetados |
|---|---|---|
| [0_SUMARIO_EXECUTIVO.md](../0_SUMARIO_EXECUTIVO.md) | Versões 0.1 e 0.2 | P1 a P7, R1 a R3 |
| [1_ESPECIFICACAO_REQUISITOS.md](../1_ESPECIFICACAO_REQUISITOS.md) | Primeira redação | RF-*, RNF-001 a RNF-036 |
| [2_ARQUITETURA_BANCO_DADOS.md](../2_ARQUITETURA_BANCO_DADOS.md) | Primeira redação | — |
| [3_ESPECIFICACAO_APIS.md](../3_ESPECIFICACAO_APIS.md) | Primeira redação | — |
| [4_FLUXOS_PROCESSO.md](../4_FLUXOS_PROCESSO.md) | Primeira redação | RN-* |
| [5_ARQUITETURA_SOFTWARE.md](../5_ARQUITETURA_SOFTWARE.md) | Primeira redação | ADR-001 a ADR-007 |
| [6_ROADMAP.md](../6_ROADMAP.md) | Primeira redação | — |
| [7_PADROES_DESIGN.md](../7_PADROES_DESIGN.md) | Primeira redação | — |
| [9_SEGURANCA_CONFIABILIDADE.md](../9_SEGURANCA_CONFIABILIDADE.md) | Primeira redação | — |
| [marca/](../../marca/) | Regras de contraste aplicadas | — |

## Verificação
Cada documento listou as inconsistências encontradas no sumário; todas foram tratadas no sumário 0.2. Contrastes WCAG calculados no documento 7.

## Pendências
P7 (cópia de recuperação de desastre fora do Brasil) criada nesta etapa — decisão do dono do produto.
