# 008 — Prompt inicial e regras de documentação

**Data:** 2026-11-12 · **Tipo:** decisão
**Pedido por:** dono do produto · **Versão do produto:** —

## Pedido
Criar o prompt inicial do projeto e uma regra de que todas as etapas sejam documentadas e salvas numa subpasta de `docs`, com integridade e referência cruzada entre documentos, e com a documentação atualizada sempre que houver correção ou implementação.

## Decisões
- **Prompt inicial:** [CLAUDE.md](../../CLAUDE.md) na raiz, lido no início de toda sessão: contexto, fontes da verdade, documento dono de cada assunto, regras obrigatórias, decisões tomadas, pendências e estrutura do repositório.
- **Registro de etapas:** pasta `docs/etapas/`, um arquivo por etapa (`NNN_AAAA-MM-DD_titulo.md`), modelo fixo, índice atualizado a cada etapa; etapas 001 a 007 registradas retroativamente.
- **Atualização obrigatória:** toda correção ou implementação atualiza, no mesmo commit, os documentos afetados, o CHANGELOG, o registro da etapa e, quando couber, o sumário.
- **Integridade e referência cruzada:** IDs criados só no documento dono, nunca renumerados; citações por ID e link relativo; verificador automático [ferramentas/verificar_docs.py](../../ferramentas/verificar_docs.py), obrigatório antes de concluir uma etapa.

## Mudanças
| Documento ou arquivo | O que mudou | IDs afetados |
|---|---|---|
| [CLAUDE.md](../../CLAUDE.md) | Criado | — |
| [docs/etapas/](./) | Pasta, regras, modelo e etapas 001 a 008 | — |
| [ferramentas/verificar_docs.py](../../ferramentas/verificar_docs.py) | Criado | — |
| [docs/README.md](../README.md), [README.md](../../README.md), [CHANGELOG.md](../../CHANGELOG.md) | Referências às regras e à pasta de etapas | — |

## Verificação
- Verificador executado sobre toda a documentação: 971 links relativos e 561 IDs definidos (206 RF, 42 RNF, 304 RN, 9 ADR) conferidos; **0 erros, 0 avisos**.
- Teste do próprio verificador com um arquivo de etapa propositalmente errado (link quebrado, IDs inexistentes, nome fora do padrão, ausente do índice): todos os erros foram acusados; arquivo de teste removido.
- PDF do CHANGELOG regerado.

## Pendências
O verificador ainda não roda no CI (não há pipeline). Entra na esteira do sprint S0 — [6_ROADMAP.md](../6_ROADMAP.md).
