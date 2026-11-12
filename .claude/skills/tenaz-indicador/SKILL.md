---
name: tenaz-indicador
description: Use sempre que um indicador da Tenaz for calculado, exibido ou alterado - disponibilidade, MTBF, MTTR, MTTA, cumprimento de preventiva, custo por ativo, custo por km, custo por hora, consumo, e os indicadores avançados. Garante a fórmula do documento 4, a base de tempo certa e o texto de "como calculamos" exigido por RF-IND-002.
---

# Indicador na Tenaz

Indicador errado é pior que indicador ausente: a pessoa decide com base nele. E indicador sem explicação visível não é indicador entregue — RF-IND-002 exige o "como calculamos" na tela.

## 1. A fórmula vem do documento, não da memória

A seção 30 de [docs/4_FLUXOS_PROCESSO.md](../../../docs/4_FLUXOS_PROCESSO.md) tem o cálculo de cada indicador. Leia a fórmula exata antes de implementar. Indicador que não está lá **não existe ainda**: proponha, registre com `tenaz-dominio-manutencao`, e só então implemente.

Base normativa: NBR 5462 para confiabilidade, EN 13306 e EN 15341 e SMRP para indicadores, ISO 14224 para a estrutura de falhas. Os valores de referência de mercado estão no [estudo de mercado](../../../docs/pesquisa/).

## 2. As cinco armadilhas deste domínio

| Armadilha | Regra correta |
|---|---|
| Contar tempo até o encerramento | KPIs de tempo usam `data_inicio` e `data_conclusao` e os registros de parada, **nunca** `data_encerramento`. Liberação administrativa tardia não pode inflar o MTTR |
| Incluir OS sem ativo nos KPIs por ativo | OS só de localização (manutenção predial) **não entra** nos KPIs por ativo |
| Ignorar o regime de operação | Disponibilidade depende das horas por dia e dos dias por semana do ativo, não de 24 por 7 por padrão |
| Confundir MTTA com MTTR | MTTA usa `data_reconhecimento`; MTTR usa início e conclusão |
| Misturar tipos de OS | Corretiva emergencial, corretiva planejada, preventiva, preditiva e inspeção legal entram em indicadores diferentes. Confira o tipo na fórmula |

Outros cuidados: dinheiro sempre em inteiro de centavos; divisão por zero tratada com resultado indefinido explícito, não com zero; período com fuso da unidade na exibição, UTC no cálculo.

## 3. Implementar

1. **Calcule no servidor**, nunca na interface. Web e app consomem o número pronto pela API.
2. **Use a réplica de leitura** quando o cálculo for pesado, conforme a seção 4.2 de [docs/5_ARQUITETURA_SOFTWARE.md](../../../docs/5_ARQUITETURA_SOFTWARE.md).
3. **Materialize só com motivo medido.** Visão materializada é otimização: entra quando o p95 exigir, com o número que o justifique.
4. **Respeite o isolamento.** O cálculo roda com contexto de empresa. Indicador é consulta como qualquer outra e entra na suíte de vazamento, inclusive na verificação de listagens e agregações.
5. **Escopo externo não vê indicador interno.** Prestador e locatário veem só o que lhes cabe.

## 4. O "como calculamos"

Obrigatório por RF-IND-002, visível na tela, junto do número. Precisa dizer:

1. A **fórmula**, em palavras que um gestor de manutenção entenda.
2. O que **entra** e o que **fica de fora** (por exemplo: OS só de localização não entra).
3. A **base de tempo** usada.
4. O **período** e o **filtro** aplicados.

Texto no pacote de internacionalização, no tom de voz da seção 10 do documento 7. Sem jargão de software.

## 5. Testar

| Teste | Como |
|---|---|
| Valor esperado | Massa pequena montada à mão, com o resultado calculado à mão. É o único jeito de provar a fórmula |
| Borda | Período sem OS, ativo sem leitura, OS sem ativo, divisão por zero, OS que atravessa a virada do período |
| Isolamento | Indicador de uma empresa não enxerga dado de outra; agregação não soma linha alheia |
| Desempenho | p95 dentro de RNF-007 com massa de 2.000 ativos |
| Coerência com a tela | O número da tela é o número da API, sem recálculo no cliente |

## 6. Fechar

1. Atualize [docs/4_FLUXOS_PROCESSO.md](../../../docs/4_FLUXOS_PROCESSO.md), seção 30, se a fórmula ou a exclusão mudou, com versão do documento e linha no histórico.
2. Capture a tela com o "como calculamos" aberto, por [`tenaz-evidencia-visual`](../tenaz-evidencia-visual/SKILL.md).
3. [`tenaz-fechar-etapa`](../tenaz-fechar-etapa/SKILL.md).
