# Entregas consolidadas

Esta pasta registra cada etapa do projeto que já foi **entregue e consolidada**: o portão de saída fechou, o dono do produto aceitou e o resultado está na branch principal (`main`). É a visão do que existe de fato, para quem precisa saber o que já pode ser usado ou construído em cima, sem ler o histórico inteiro.

Ela não substitui o [registro de etapas](../etapas/README.md). O registro conta **como** cada etapa aconteceu: pedido, decisões, mudanças, verificação e pendências, inclusive das etapas em andamento. A entrega diz **o que ficou pronto** e em que condição: o que foi entregue, o estado de cada item do portão, onde está a evidência, de qual commit se trata e o que ficou pendente para depois.

## Regras

1. **Só entra etapa consolidada.** Três condições, todas obrigatórias: portão de saída fechado (item adiado só conta se o dono do produto o deslocou por decisão registrada no prompt da etapa), aceite do dono do produto e commit em `main`. Etapa em andamento ou com portão vermelho não entra.
2. **Um arquivo por entrega**, com o identificador da etapa do [plano de execução](../execucao/4_ETAPAS.md): `E00_fundacao-tecnica.md`, `E01_conta-e-acesso.md` e assim por diante. A concepção, anterior às etapas de construção, é `C0_concepcao.md`.
3. **A entrega é escrita no fechamento da etapa**, no mesmo commit do registro da etapa, pelo procedimento da skill `tenaz-fechar-etapa`.
4. **Entrega publicada não é reescrita.** Se algo entregue mudar depois (correção, regressão, item adiado que finalmente entra), acrescenta-se uma nota datada no fim do arquivo, citando a etapa que mudou, e o índice abaixo é atualizado.
5. **Toda afirmação aponta para a prova:** commit, registro da etapa, evidência em [../screenshots/](../screenshots/README.md), saída de teste ou documento dono. Nada que o código ou o CI não sustentem.

## Modelo

```markdown
# E00 — Título da entrega

**Etapa:** E00 · **Marco:** M0 · **Consolidada em:** AAAA-MM-DD
**Commits:** abc1234, def5678 · **Registro da etapa:** [NNN](../etapas/NNN_...md) · **Prompt:** [E00](../prompts/E00_...md)

## O que foi entregue
Em linguagem de quem vai usar ou construir em cima.

## Requisitos
| Requisito | Situação nesta entrega |
|---|---|

## Portão de saída
| Verificação | Resultado | Prova |
|---|---|---|

## Evidência
Onde estão as capturas e as saídas.

## O que ficou para depois
Item, motivo, para onde foi.

## Notas posteriores
(vazio até que algo entregue mude)
```

## Índice

| Entrega | Etapas do registro | Consolidada em | Commits | Situação |
|---|---|---|---|---|
| [C0 — Concepção](C0_concepcao.md) | 001 a 010 e 012 | 2026-11-12 | 8e4b3fc, ed51d08, 1dd31c0, 2011feb, 294d7f1 | Consolidada |
| [E00 — Fundação técnica](E00_fundacao-tecnica.md) | 011 | 2026-11-12 | 294d7f1, fc9cbbc, bd62bee | Consolidada, com dois itens do portão deslocados para a fase de produção |

A etapa **E01 — Conta e acesso** está em andamento ([registro 013](../etapas/013_2026-11-12_conta-e-acesso.md)) e entra aqui quando for aceita.
