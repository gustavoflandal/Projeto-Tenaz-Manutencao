# Modelo de prompt de etapa

Copie o conteúdo abaixo da linha para um arquivo novo `E<NN>_titulo-curto.md` e preencha. Regras em [README.md](README.md). Procedimento na skill [`tenaz-prompt-etapa`](../../.claude/skills/tenaz-prompt-etapa/SKILL.md).

---

```markdown
# E<NN> — <Título da etapa>

**Data do pedido:** AAAA-MM-DD · **Sprint do roadmap:** S<N> · **Marco:** M<N> ou —
**Pedido por:** dono do produto

## Objetivo

Uma a três frases: o que o usuário passa a conseguir fazer no fim desta etapa.

## Entrega

Os requisitos, pelos IDs do documento 1. Nada descrito com outras palavras.

| Requisito | Observação de escopo nesta etapa |
|---|---|
| RF-... | completo, ou a parte que entra aqui |

## Contexto obrigatório

Os documentos a ler antes de começar, com a seção.

| Documento | Seção | Por quê |
|---|---|---|

## Agentes previstos

Quem conduz e quem apoia, conforme docs/execucao/2_EQUIPE_AGENTES.md.

## Portão de saída

Verificações que precisam estar verdes, cada uma com o comando ou a prova que a demonstra.

| Verificação | Como se comprova |
|---|---|

## Evidência esperada

As capturas que vão para docs/screenshots/E<NN>_<titulo>/.

## Fora do escopo

O que não entra, e em que etapa ou fase entra.

## Decisões pendentes do dono do produto

Pendências P1 a P8 ou escolhas que travam a etapa. Vazio, se não houver.

## Ajustes durante a etapa

Acrescentado só durante a execução, com data e motivo. Nunca reescreve o que está acima.
```
