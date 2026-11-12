---
name: tenaz-operacao-offline
description: Use sempre que uma operação da Tenaz precisar funcionar sem conexão - executar OS, registrar leitura, preencher checklist, abrir solicitação, registrar componente, assinar, fazer vistoria. Garante fila de comandos, identificador gerado no aparelho, idempotência, numeração provisória, regra de conflito definida antes do código e teste de desconexão, ação, reconexão e conflito.
---

# Operação offline no app da Tenaz

O princípio 3 da arquitetura: **o campo não tem sinal.** A meta é 7 dias de operação sem conexão sem perda de registro (RNF-010). Offline não é um recurso que se acrescenta no fim: é a forma como a operação nasce.

## 1. Antes de escrever código

Responda estas seis perguntas **por escrito**. Operação offline sem elas respondidas vira perda de dado em campo.

| # | Pergunta | Por que importa |
|---|---|---|
| 1 | Que comando o app envia? | A fila é de intenção ("concluir OS"), não de estado. Comando para cima, alteração para baixo |
| 2 | O registro nasce no aparelho? | Se sim, identificador UUIDv7 gerado localmente, definitivo |
| 3 | Tem número legível? | Então é **provisório** até o servidor confirmar, e a tela mostra isso |
| 4 | O que acontece se o servidor já mudou? | A regra de conflito, decidida agora, por escrito |
| 5 | Que dados precisam estar em cache para a operação funcionar? | O escopo do cache vem do requisito, não do seu julgamento |
| 6 | A operação envolve foto ou arquivo? | Compressão, fila separada, interface que não bloqueia |

## 2. Implementar

1. **Identificador no aparelho:** UUIDv7, ordenável por tempo. Permite o registro nascer offline com identidade definitiva.
2. **Validação local idêntica à do servidor:** use `pacotes/dominio` e os casos compartilhados em JSON. A máquina de estados da OS e a validação de leitura decidem igual nos dois lados. Divergência faz o app aceitar o que o servidor recusa depois — defeito crítico.
3. **Fila de comandos durável:** cada comando com chave de idempotência própria. Reenvio não duplica; o servidor repete a resposta.
4. **Cache filtrado por escopo:** app de prestador traz só o do fornecedor dele; app de técnico de subcontratada, só o da conta dele. Verificado em teste.
5. **Um banco local por conta:** troca de conta não mistura dado.
6. **Numeração provisória visível:** o usuário precisa saber que o número vai mudar.
7. **Fotos:** comprimidas, enviadas em fila própria, sem travar a tela.

## 3. Caso especial: OS ligada a cliente conectado

O comando é aplicado **primeiro na conta da subcontratada** e só depois sincronizado com a conta do cliente. O app **nunca** escreve direto na conta do cliente. Risco 3 do cronograma no [roadmap](../../../docs/6_ROADMAP.md).

## 4. Conflito

Escolha a regra e registre-a:

| Regra | Quando serve |
|---|---|
| O comando vence | Registro de campo que o servidor não tinha: leitura, apontamento, foto |
| O servidor vence, com aviso | O estado mudou de forma incompatível: a OS foi cancelada enquanto o técnico a executava |
| Pergunta ao usuário | Divergência que só a pessoa resolve, com os dois valores à vista |

**Nunca** descarte em silêncio. **Nunca** sobrescreva sem avisar. Conflito sem regra escrita é perda de dado disfarçada de sincronização.

## 5. Testar

Cenário mínimo, automatizado, para **toda** operação offline:

1. Desconectar.
2. Executar a operação, mais de uma vez quando fizer sentido.
3. Reconectar.
4. Verificar que o resultado no servidor é o esperado, uma única vez.
5. Provocar o conflito: alterar no servidor antes de reconectar, e verificar a regra escolhida.
6. Reenviar o mesmo comando e verificar que não duplica.

Mais: fila longa, troca de conta sem mistura, foto grande, atualização do app **sem perder a fila**.

Nos marcos: 72 horas offline em E06, 7 dias em E11.

## 6. O que nunca se faz

- Perder, truncar ou descartar a fila de comandos, em nenhuma circunstância, inclusive em atualização do app.
- Mudar regra de domínio para o offline ficar mais simples.
- Inventar a numeração definitiva no aparelho.
- Deixar dado de saúde ou dado pessoal sensível no cache.
- Guardar credencial fora do armazenamento seguro do sistema.

## 7. Fechar

Critério 5 dos nove critérios de pronto exige o teste de desconexão, ação, reconexão e conflito. Depois, [`tenaz-fechar-etapa`](../tenaz-fechar-etapa/SKILL.md).
