# 010 — Troca de repositório no GitHub

**Data:** 2026-11-12 · **Tipo:** release
**Pedido por:** dono do produto · **Versão do produto:** —

## Pedido
Fazer o commit para o repositório `https://github.com/gustavoflandal/Projeto-Tenaz-Manutencao.git`.

## Decisões

- **O repositório novo substitui o antigo.** O `origin` apontava para `gustavoflandal/tenaz-manutencao`, que **não existe mais** no GitHub: `git ls-remote` responde "Repository not found". O repositório `Projeto-Tenaz-Manutencao` existe, é público e estava vazio. Não há, portanto, dois repositórios a manter: o antigo foi removido ou renomeado fora deste repositório local.
- **`origin` repontado**, em vez de um segundo remoto. Guardar um remoto para um repositório inexistente só produziria erro em todo `fetch`.
- **URL atualizada em [CLAUDE.md](../../CLAUDE.md), seção 1**, no mesmo commit da troca, porque o arquivo documenta o endereço do repositório e a regra 3.2 exige documentação atualizada junto com a mudança.
- **Histórico preservado.** O commit `8e4b3fc` (concepção) e os seguintes foram enviados ao repositório novo como estão; nada foi reescrito nem recriado.
- Alternativa descartada: recomeçar o histórico no repositório novo. Perderia a rastreabilidade das etapas 001 a 009, que é a própria razão da pasta de etapas.

## Mudanças

| Documento ou arquivo | O que mudou | IDs afetados |
|---|---|---|
| Repositório | `origin` repontado para `gustavoflandal/Projeto-Tenaz-Manutencao`; `main` enviada com todo o histórico | — |
| [CLAUDE.md](../../CLAUDE.md) | Seção 1: URL do repositório | — |
| [docs/etapas/](./) | Este registro e o índice | — |

Commits enviados nesta etapa, além dos já existentes:

| Commit | Conteúdo |
|---|---|
| `1dd31c0` | Equipe de 16 agentes, 12 skills e prompts das etapas E02 a E12, completando a etapa [009](009_2026-11-12_equipe-de-agentes-e-skills.md) |
| este | Troca de repositório e registro desta etapa |

## Verificação

- `python ferramentas/verificar_docs.py`: **0 erros, 0 avisos**.
- `git ls-remote origin` respondendo no repositório novo, e `main` igual a `origin/main` depois do envio.
- Conferido que o commit `ed51d08`, feito antes de esta etapa começar, havia capturado o plano de execução **pela metade**: faltavam `.claude/` inteiro, o registro da etapa 009 e os prompts E02 a E12. O commit `1dd31c0` completou o conjunto.

## Pendências

- O repositório é **público**. Se a intenção for repositório privado antes de haver código e material comercial, a mudança é de um clique nas configurações do GitHub e não depende deste repositório local.
- A pendência **P5** (disponibilidade do nome Tenaz no INPI) segue aberta na seção 6 do [sumário](../0_SUMARIO_EXECUTIVO.md). O nome do repositório não a resolve nem a agrava.
