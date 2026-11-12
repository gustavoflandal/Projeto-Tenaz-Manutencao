---
name: tenaz-release
description: Use para fechar uma versão da Tenaz em produção, conforme RNF-039 e RNF-040 - tag, CHANGELOG da versão, versão igual e injetada pelo build em web, app, API e workers, verificações que bloqueiam o release no CI, publicação do app, novidades da versão e confirmação depois do deploy. Também para release de correção.
---

# Release em produção da Tenaz

Dois requisitos do dono do produto governam isto: **RNF-039** (toda entrega tem versão única, rastreável até o commit, visível no canto inferior direito) e **RNF-040** (histórico de versões sempre atualizado). A seção 3.3 do [roadmap](../../../docs/6_ROADMAP.md) e a 9.5 de [docs/5_ARQUITETURA_SOFTWARE.md](../../../docs/5_ARQUITETURA_SOFTWARE.md) têm o detalhe.

## 1. Escolher o incremento

Versionamento semântico, a partir da tag:

| Incremento | Quando |
|---|---|
| **MAIOR** | Quebra de compatibilidade da API pública. Coincide com uma nova versão de caminho (`/v2`) |
| **MENOR** | Funcionalidade nova compatível |
| **CORREÇÃO** | Correção |

Os commits convencionais desde a última tag sugerem o incremento. **Uma pessoa revisa o texto e aprova**; a ferramenta propõe, não decide. Pré-lançamento `-rc.N` existe só fora da produção.

## 2. Preparar

1. **CHANGELOG:** a seção de não lançado vira `## [X.Y.Z] - AAAA-MM-DD`, com **todas** as implementações (com ID de requisito) e correções (com número do chamado) do release. Não pode estar vazia.
2. **Documentação:** todo documento alterado no ciclo com versão subida e PDF regerado.
3. **Verificador:** `python ferramentas/verificar_docs.py` com zero erro.
4. **Portões abertos:** suítes de vazamento e de três contas verdes; fluxos E2E obrigatórios verdes; p95 dentro da meta; nenhum defeito crítico aberto.
5. **Tag anotada** `vX.Y.Z` na branch principal, maior que a anterior.

## 3. Build e injeção da versão

**A versão nunca é digitada à mão em arquivo de código.** Vem da tag:

| Componente | Como |
|---|---|
| Backend | Versão, commit e data do build injetados no build. O subcomando de versão imprime os três; a versão vai nos logs, nos traces e no cabeçalho de resposta |
| Web | Variáveis de build preenchidas pelo CI a partir da tag; o pacote publica o arquivo de versão junto dos estáticos |
| App | Versão do release no arquivo de configuração; número de build inteiro e crescente gerado pelo CI, exigido pelas lojas; versão de runtime pela política de impressão digital, que muda só quando a parte nativa muda |
| Imagem | Construída **uma vez** e promovida pelo digest; na promoção para produção, recebe também a tag da versão |

Web, app, API e workers do mesmo release têm a **mesma versão**.

## 4. Verificações que bloqueiam o release no CI

1. A tag segue o formato e é maior que a anterior.
2. O CHANGELOG tem a seção da versão, não vazia.
3. A versão impressa pelo backend, a do arquivo de versão da web e a da configuração do app são **iguais à tag**, e o commit é o da tag.
4. Se a comparação de contrato apontar mudança incompatível na API pública, o incremento tem de ser MAIOR; caso contrário, MAIOR é recusado.
5. O deploy em produção **recusa** versão de pré-lançamento e versão sem tag.

Qualquer uma vermelha: o release não sai. Não há atalho manual de versão (risco 20 do cronograma).

## 5. Deploy

1. Migrações como tarefa única antes do deploy, compatíveis com a versão anterior do código.
2. Deploy contínuo com verificação de saúde e **reversão automática** se o alarme de erro ou de latência disparar nos primeiros 10 minutos.
3. Web por publicação versionada mais invalidação de cache; reversão trocando o ponteiro de versão.
4. App: build para as lojas; atualização sem loja só para mudança de JavaScript compatível com o binário, por canal, com liberação gradual.

## 6. Depois do deploy

1. A verificação sintética chama `GET /v1/versao` e confirma que a produção responde a versão liberada.
2. O selo de versão aparece no canto inferior direito da tela principal e do login, na web e no app. A tela "Sobre" mostra versão, versão da API, data do build, commit curto e ambiente; no app, também o binário da loja e o identificador da atualização sem loja.
3. "Novidades desta versão" mostra a seção do CHANGELOG publicada pelo build.
4. Registre a etapa de release com [`tenaz-fechar-etapa`](../tenaz-fechar-etapa/SKILL.md), tipo release.

## 7. Cuidado com o app offline

O app opera até 7 dias sem conexão, então um aparelho pode estar alguns releases atrás. A API aceita apps da mesma versão maior; abaixo da versão mínima suportada, o app pede atualização **sem descartar a fila de comandos**. Nunca force atualização que perca a fila.

## 8. Release de correção

Mesmo fluxo, sem exceção: tag, CHANGELOG, build, verificações. Urgência não autoriza atalho — é exatamente quando o atalho causa o incidente seguinte.
