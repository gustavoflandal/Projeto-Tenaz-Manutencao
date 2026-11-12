---
name: tenaz-prompt-etapa
description: Use antes de começar qualquer etapa da Tenaz, para escrever o prompt dela em docs/prompts. Também para registrar um ajuste de rumo no meio de uma etapa já em andamento, sem reescrever o pedido original.
---

# Escrever o prompt de uma etapa

O prompt é escrito **antes** do trabalho. Prompt escrito depois é relatório, não pedido.

Por que ele existe: meses depois, o código mostra o que foi feito e o registro em [docs/etapas/](../../../docs/etapas/) mostra o que mudou, mas só o prompt mostra **o que se pediu e com que informação**. É o que permite separar "a decisão estava errada" de "o pedido estava ambíguo".

## 1. Antes de escrever

1. Leia a etapa em [docs/execucao/4_ETAPAS.md](../../../docs/execucao/4_ETAPAS.md): entrega, agentes, portão, evidência, dependências.
2. Leia a sprint correspondente no [roadmap](../../../docs/6_ROADMAP.md), seção 4, e o critério do marco, se a etapa fechar um.
3. Confira as **dependências críticas** da seção 8 do roadmap que tocam esta etapa.
4. Confira que a etapa anterior está fechada. Etapas não se sobrepõem.

## 2. Escrever

Copie [docs/prompts/000_modelo.md](../../../docs/prompts/000_modelo.md) para `E<NN>_titulo-curto.md` e preencha as nove seções.

| Seção | O que não pode faltar |
|---|---|
| **Objetivo** | Uma a três frases, em termos de **o que o usuário passa a conseguir fazer**. Não em termos de componente técnico |
| **Entrega** | Os requisitos **pelos IDs**, numa tabela, com a observação de escopo de cada um nesta etapa (completo, ou qual parte). Nunca descreva um requisito com outras palavras em vez de citá-lo |
| **Contexto obrigatório** | Documento **e seção**. "Leia o documento 4" não serve: ele tem 3.300 linhas |
| **Agentes previstos** | Quem conduz e quem apoia, pelos nomes de [docs/execucao/2_EQUIPE_AGENTES.md](../../../docs/execucao/2_EQUIPE_AGENTES.md) |
| **Portão de saída** | Cada verificação com **o comando ou a prova** que a demonstra. "Funcionando bem" não é portão; "suíte de três contas verde no CI" é |
| **Evidência esperada** | As capturas, e a pasta de destino em [docs/screenshots/](../../../docs/screenshots/) |
| **Fora do escopo** | O que não entra, **e em que etapa ou fase entra**. É esta seção que evita o crescimento de escopo, risco 2 do cronograma |
| **Decisões pendentes** | Pendências P1 a P8 ou escolhas que travam a etapa. Vazio, se não houver |
| **Ajustes durante a etapa** | Vazio no começo |

## 3. Regras

1. **Nome:** `E<NN>_titulo-curto.md` no MVP; `F2-M<NN>_titulo.md` e `F3-M<NN>_titulo.md` nas fases 2 e 3. Minúsculas, sem acento, com hífens. **O mesmo nome da pasta de evidência.**
2. **Cite por ID e link relativo.** Todo ID citado precisa existir no documento dono, ou o verificador reprova.
3. **Prompt publicado não é reescrito.** Mudança de rumo entra em "Ajustes durante a etapa", com data e motivo.
4. **Nenhum segredo:** sem senha, token, chave ou dado pessoal real.
5. **Atualize o índice** de [docs/prompts/README.md](../../../docs/prompts/README.md).

## 4. Teste do prompt

Antes de dar o prompt por pronto, responda:

- [ ] Alguém que não esteve na conversa consegue executar a etapa só com isto?
- [ ] Todo requisito está citado por ID?
- [ ] Todo item do portão tem um comando ou uma prova associada?
- [ ] A seção de fora do escopo diz onde cada item excluído entra?
- [ ] As pendências que travam a etapa estão nomeadas?
- [ ] O verificador passa?

```bash
python ferramentas/verificar_docs.py
```

Qualquer "não" significa prompt incompleto. Um prompt vago produz uma etapa vaga, e o custo aparece no aceite.
