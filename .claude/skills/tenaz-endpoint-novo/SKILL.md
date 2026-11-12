---
name: tenaz-endpoint-novo
description: Use sempre que nascer ou mudar uma rota, um campo de JSON, um código de erro, um evento de webhook ou um escopo de token na API da Tenaz. Garante o contrato único - OpenAPI gerado e versionado, sem mudança incompatível em v1, erro no formato RFC 9457, paginação por cursor, idempotência, ETag e tipos TypeScript regenerados.
---

# Endpoint novo na API da Tenaz

O princípio é o 4 da arquitetura: **um contrato, vários clientes**. A API pública, a web, o app e os integradores usam o mesmo `/v1`. Criar rota que só a web usa quebra esse princípio e cria uma API interna paralela, que não existe neste projeto.

## 1. Antes de criar

1. Procure a rota no catálogo de [docs/3_ESPECIFICACAO_APIS.md](../../../docs/3_ESPECIFICACAO_APIS.md), seção 6, e no detalhamento por módulo. Quase tudo já está especificado.
2. Confira o nome do recurso no [sumário](../../../docs/0_SUMARIO_EXECUTIVO.md), seção 3.4.
3. Pergunte-se: é recurso novo ou **ação de estado** de um recurso existente? Ação de estado é sub-recurso verbal, não rota nova.

## 2. Forma da rota

| Item | Regra | Exemplo |
|---|---|---|
| Caminho | `/v1/<recurso-no-plural>` em kebab-case | `/v1/ordens-servico`, `/v1/leituras-medidor`, `/v1/clientes-servico` |
| Ação de estado | Sub-recurso verbal | `POST /v1/ordens-servico/{id}/concluir`, `/assumir`, `/encerrar` |
| JSON | `snake_case`, igual ao banco | `data_abertura` |
| Exceção aceita | `GET /v1/versao`, no singular | — |

## 3. Comportamento obrigatório

| Tema | Regra |
|---|---|
| Paginação | Por cursor, nunca por deslocamento em lista grande |
| Idempotência | Cabeçalho de chave de idempotência aceito em **todo** `POST`; resposta repetida por 24 horas |
| Concorrência | ETag e `If-Match` em atualização de recurso editável; 412 em conflito |
| Erro | RFC 9457 em `application/problem+json`, com campo de código estável e mensagem em português que diz o que fazer |
| Recurso de outra empresa | **404, nunca 403.** 403 confirma existência |
| Limite | Por plano; 429 com `Retry-After` |
| Escopo | Toda rota declara se aceita `interno`, `prestador`, `locatario` ou `link` |

## 4. Escopo de acesso — o passo mais esquecido

1. Declare explicitamente os escopos aceitos.
2. Rota não declarada para escopo externo responde **404** para ele. A suíte de `tenaz-isolamento` varre todas as rotas do contrato automaticamente; rota nova sem declaração aparece como falha.
3. Rota que o prestador usa: o filtro por fornecedor fica **no servidor**, além do RLS. Nunca só na interface.
4. Rota que expõe dado ligado a outra conta conectada: **nunca** exponha identificador da outra conta. Use o número legível da OS remota.

## 5. Contrato

1. Escreva os tipos de requisição e resposta como structs Go anotadas. O contrato OpenAPI 3.1 é **gerado** a partir delas.
2. Regenere `backend/api/openapi.yaml` e **versione o arquivo gerado**. O CI falha se o gerado divergir do versionado.
3. Rode a comparação com a branch principal. **Mudança incompatível em `/v1` bloqueia o merge** e exige `/v2`, com convivência mínima de 12 meses.
4. Regenere os tipos TypeScript de `pacotes/api-tipos`. Campo renomeado no Go deve quebrar a compilação do TypeScript no mesmo pull request — isso é o mecanismo funcionando, não um problema.
5. **Nunca edite à mão** o contrato gerado nem os tipos gerados.

## 6. Webhook, quando houver

Nome no padrão `<entidade>.<fato_no_particípio>`, com o nome **completo** da entidade: `ordem_servico.encerrada`, `documento_vencimento.vencendo`, `leitura_medidor.registrada`. Entrega com repetição, assinatura e registro de tentativa, conforme a seção 13 do documento 3.

## 7. Verificar

| Verificação | Como |
|---|---|
| Contrato sem divergência | Geração no CI, comparando com o versionado |
| Sem quebra de compatibilidade | Comparação com a branch principal |
| Isolamento por rota | Suíte de vazamento: token de uma empresa com identificador de outra, 404 e nenhum identificador alheio no corpo |
| Escopo externo | Rota não declarada respondendo 404 para prestador, locatário e link |
| Idempotência | Mesmo `POST` com a mesma chave duas vezes, um só efeito |
| Erro | Formato RFC 9457 com código estável |

## 8. Fechar

Atualize [docs/3_ESPECIFICACAO_APIS.md](../../../docs/3_ESPECIFICACAO_APIS.md) com a rota, os erros e os escopos, suba a versão do documento e a linha do histórico. Depois, [`tenaz-fechar-etapa`](../tenaz-fechar-etapa/SKILL.md).
