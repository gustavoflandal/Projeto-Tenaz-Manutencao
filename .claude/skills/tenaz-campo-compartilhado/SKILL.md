---
name: tenaz-campo-compartilhado
description: Use sempre que nascer ou mudar um campo de ordem de serviço ou de ativo na Tenaz, e em qualquer trabalho na sincronização entre contas conectadas do módulo SRV. Classifica o campo como compartilhável ou nunca compartilhado, atualiza a lista fechada e o privilégio do papel do sincronizador, e acrescenta o caso na suíte de três contas.
---

# Campo que trafega entre contas conectadas

Esta skill sustenta o risco **R5** do [sumário](../../../docs/0_SUMARIO_EXECUTIVO.md): a sincronização entre contas conectadas copiar dado além do permitido ou para a conta errada. É crítico, e é o único risco do projeto sem paralelo em nenhuma outra parte do produto.

O critério de pronto do [roadmap](../../../docs/6_ROADMAP.md), item 2 da seção 3.1, é explícito: **todo campo novo de OS ou de ativo é classificado como "compartilhável" ou "nunca compartilhado"**. A classificação vem antes de o campo existir.

## 1. Classificar

Para cada campo, responda:

| Pergunta | Se a resposta for sim |
|---|---|
| É dado pessoal? CPF, contato pessoal, documento, dado de saúde | **Nunca compartilhado.** Sem discussão |
| É custo, preço ou margem de uma das contas? | Nunca compartilhado por padrão. Preço da subcontratada só trafega em grupo opcional que **ela** liga |
| É histórico interno, observação interna, dado de outro fornecedor? | Nunca compartilhado |
| A outra conta precisa dele para **executar** o trabalho? | Candidato a compartilhável |
| A outra conta precisa dele para **prestar contas** do trabalho? | Candidato a compartilhável |

Na dúvida, **nunca compartilhado**. Acrescentar um campo à lista depois é fácil; retirar um campo que já vazou é impossível.

## 2. Registrar na lista fechada

1. Acrescente o campo à lista fechada de campos da conexão, guardada no repositório.
2. Indique o **sentido**: do cliente para a subcontratada, da subcontratada para o cliente, ou os dois.
3. Indique o **grupo**: obrigatório da conexão ou opcional. Grupo opcional só trafega se ligado; o grupo de valores só trafega se a subcontratada também o ligar.
4. Atualize o privilégio do papel do sincronizador para **exatamente** essas colunas. Nunca privilégio na tabela inteira.

A suíte verifica que as colunas com privilégio **batem exatamente** com a lista. Coluna fora da lista com privilégio quebra o build.

## 3. Validar nas duas pontas

| Regra | Detalhe |
|---|---|
| Validação no envio **e** no recebimento | Confiar só em quem envia não protege de erro de programação na outra ponta |
| Mensagem inválida é recusada **inteira** | Campo extra, campo de outro dono, campo não classificado: recusa total, erro registrado, alerta de segurança. **Nunca** se aplica "a parte válida" |
| Destino vem da conexão | Nunca do job. Job cuja empresa de origem não é ponta da conexão é recusado |
| Dono do campo | Cada campo tem um lado que pode alterá-lo. A subcontratada não altera a descrição da OS do cliente nem envia encerramento |
| Transação separada por conta | Uma gravação por conta, nunca uma transação abrangendo as duas |
| Auditoria nos dois lados | Mesma identificação de conteúdo na origem e no destino |

## 4. Acrescentar o caso na suíte de três contas

A suíte tem dois clientes (C1 e C2) e uma subcontratada (S) conectada aos dois. Para o campo novo, acrescente:

1. **Valor sentinela** no campo, quando ele for "nunca compartilhado": o teste procura o sentinela em **todas** as tabelas da outra conta e falha se o encontrar.
2. **Campo injetado:** mensagem montada à mão com o campo, enviada por quem não é dono dele, recusada inteira.
3. **Sem transitividade:** nada de C1 chegando a C2 por meio de S.
4. **Conta errada:** mensagem da conexão de C1 com S tentando gravar em C2, recusada.
5. **Replay e ordem:** mesma mensagem duas vezes aplicada uma; sequência antiga ignorada; fora de ordem aplicada na ordem da sequência.
6. **Estados da conexão:** convidada, suspensa e encerrada não aplicam nada, fora da janela de 24 horas do encerramento.

## 5. Verificar

| Verificação | Prova |
|---|---|
| RNF-041 | Suíte de três contas verde no CI, com todos os casos da seção 9.3 de [docs/5_ARQUITETURA_SOFTWARE.md](../../../docs/5_ARQUITETURA_SOFTWARE.md) |
| RNF-042 | Atraso p95 da sincronização abaixo de 60 segundos em homologação, e sob carga no endurecimento |
| Cobertura do papel | Colunas com privilégio iguais à lista fechada; papel sem privilégio de contornar RLS e sem ser dono de tabela |

**Quem declara o portão verde é `tenaz-isolamento`**, não quem implementou.

## 6. Fechar

1. Atualize [docs/2_ARQUITETURA_BANCO_DADOS.md](../../../docs/2_ARQUITETURA_BANCO_DADOS.md) na parte de sincronização e [docs/3_ESPECIFICACAO_APIS.md](../../../docs/3_ESPECIFICACAO_APIS.md), seção 22, quando o campo aparecer na API.
2. Se a mudança afetar o vocabulário ou a lista de grupos, atualize o [sumário](../../../docs/0_SUMARIO_EXECUTIVO.md), seções 3.5 e 3.6.
3. [`tenaz-fechar-etapa`](../tenaz-fechar-etapa/SKILL.md).

## Lembrete

O termo de LGPD da conexão entre contas precisa estar aprovado pelo advogado **antes do primeiro convite fora de homologação**. É dependência crítica 18 do [roadmap](../../../docs/6_ROADMAP.md).
