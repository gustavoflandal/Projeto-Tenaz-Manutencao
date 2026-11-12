# 004 — Terceiros, estoque estendido e locação

**Data:** 2026-11-12 · **Tipo:** escopo
**Pedido por:** dono do produto · **Versão do produto:** —

## Pedido
Controlar empresas subcontratadas de manutenção, com estoque geral e por subcontratada, transferência e retorno de peças; peças com defeito, envio para conserto e retorno; cadastro de assistências técnicas com descritivo de serviços; manutenção de equipamentos locados e controle de locatários. Locação nos dois sentidos (confirmado pelo dono do produto).

## Decisões
- Módulo **TER**: cadastro único de fornecedores com papéis, catálogo de serviços, contratos com SLA, OS executada por prestador, perfil `prestador`, link de uso único. A oficina externa saiu de FRO para TER.
- Módulo **EST** estendido: almoxarifado `em_terceiro` e `quarentena`, peças serializadas e reparáveis, ordem de reparo `OR-`, remessa e retorno com NF-e registrada.
- Módulo **LOC**: contratos `cedido` e `tomado`, locatários, matriz de responsabilidade, vistorias, portal do locatário.
- Sumário 0.3; riscos R4 e pendência P8 (CFOPs).

## Mudanças
| Documento ou arquivo | O que mudou | IDs afetados |
|---|---|---|
| [0_SUMARIO_EXECUTIVO.md](../0_SUMARIO_EXECUTIVO.md) | Versão 0.3 | P8, R4 |
| Documentos 1 a 9 | Alinhados ao sumário 0.3 | RF-TER-*, RF-LOC-*, RF-EST-012 a 019, RNF-037, RNF-038, RN-TER-*, RN-LOC-*, ADR-008 |

## Verificação
Busca por termos antigos em todos os documentos; checagem por script dos IDs RN e RF citados.

## Pendências
P8 (CFOPs de remessa, conserto e locação) — validar com contador.
