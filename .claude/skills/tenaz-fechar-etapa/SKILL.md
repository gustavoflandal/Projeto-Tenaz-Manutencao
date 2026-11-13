---
name: tenaz-fechar-etapa
description: Use ao terminar qualquer etapa, correção ou implementação na Tenaz. Executa as quatro atualizações obrigatórias de CLAUDE.md - documentos afetados com versão subida, CHANGELOG, registro da etapa em docs/etapas e sumário quando couber - escreve a entrega consolidada em docs/entregas quando a etapa é aceita, roda o verificador e regera os PDFs.
---

# Fechar uma etapa da Tenaz

A regra do projeto é dura: **mudança de comportamento que não aparece na documentação é tratada como defeito** ([CLAUDE.md](../../../CLAUDE.md), seção 3.2). Esta skill é o procedimento que a cumpre.

Tudo isto acontece **no mesmo commit ou pull request** das mudanças que descreve. Não é uma tarefa de limpeza para depois.

## 1. Documentos afetados

Para cada documento de `docs/` que a etapa alterou:

1. Aplique a mudança com o **agente dono** do assunto. Requisitos com `tenaz-produto`, banco com `tenaz-banco-dados`, API com `tenaz-api-contrato`, fluxos e regras com `tenaz-dominio-manutencao`, arquitetura e ADR com `tenaz-arquiteto`, design com `tenaz-design-sistema`, segurança com `tenaz-plataforma-sre` ou `tenaz-conformidade-br`.
2. **Suba a versão** do documento no cabeçalho.
3. Acrescente a linha no **histórico de versões** do fim do documento.
4. Procure as **citações desse assunto nos outros documentos** e atualize-as junto. Citação órfã é o defeito documental mais comum deste projeto.

## 2. CHANGELOG

Em [CHANGELOG.md](../../../CHANGELOG.md), na seção de não lançado:

| Regra | Detalhe |
|---|---|
| Grupo | Adicionado, Alterado, Corrigido, Removido, Obsoleto ou Segurança |
| Linguagem | De usuário, não de desenvolvedor. Quem lê é cliente |
| Referência | O ID do requisito para implementação; o número do chamado ou da issue para correção |
| Quebra de compatibilidade | Marcada, e exige versão maior |
| Entrada publicada | Nunca reescrita. Erro se corrige com nota na própria entrada |

Omitir isto viola RNF-040, e um check do CI reprova o pull request que altera código de produto sem alterar o CHANGELOG.

## 3. Registro da etapa

Em [docs/etapas/](../../../docs/etapas/), arquivo `NNN_AAAA-MM-DD_titulo-curto.md`:

1. **`NNN` sequencial de três dígitos, nunca reaproveitado.** Olhe o último número da pasta.
2. Título em minúsculas, sem acento, com hífens.
3. A primeira linha começa com `# NNN ` e o travessão. O verificador cobra isso.
4. Seções fixas: **Pedido** (fiel ao que foi pedido), **Decisões** (cada uma com o motivo, e as alternativas descartadas), **Mudanças** (tabela de arquivo e IDs afetados), **Verificação** (como foi conferido, com a saída dos comandos), **Pendências** (o que ficou aberto e onde está registrado).
5. **Atualize o índice** no README da pasta.
6. Cite as capturas de [docs/screenshots/](../../../docs/screenshots/) por link relativo.

Registro publicado **não é reescrito**: erro se corrige com nota de correção datada no fim do próprio registro, ou com etapa nova que o cite.

## 3-b. Entrega consolidada

Quando a etapa é **aceita pelo dono do produto** e vai para `main`, escreva a entrega em [docs/entregas/](../../../docs/entregas/README.md), no formato e nas regras do README da pasta: arquivo `<ID da etapa>_<titulo>.md` (ex.: `E01_conta-e-acesso.md`), com o que foi entregue, a situação de cada requisito, o portão item a item com a prova, a evidência e o que ficou para depois. Atualize o índice da pasta. Etapa em andamento ou com portão vermelho não entra. Se a etapa mudou algo de uma entrega anterior, acrescente uma nota datada em "Notas posteriores" daquela entrega.

## 4. Sumário executivo

Atualize [docs/0_SUMARIO_EXECUTIVO.md](../../../docs/0_SUMARIO_EXECUTIVO.md) quando a mudança afetar **módulo, plano, vocabulário, estado ou pendência**. Nome canônico novo entra aqui primeiro, e só depois nos outros documentos.

## 5. Pendências

Pendência aberta na etapa vai para **dois lugares**: a seção de pendências do documento dono (ou do sumário) e o registro da etapa. Nunca fica só na conversa. A etapa que resolver uma pendência cita a etapa que a abriu.

## 6. Verificar

```bash
python ferramentas/verificar_docs.py
```

**Zero erro.** Ele confere links relativos de todos os `.md` do repositório (inclusive `.claude/`), IDs citados contra o documento dono, linha de versão nos documentos 0 a 9, padrão e numeração das etapas com o índice, e avisa sobre PDF desatualizado.

Atenção: o verificador também lê IDs **dentro de blocos de código**. Identificador de exemplo inventado quebra a verificação; use a forma genérica `RF-<MÓDULO>-<NNN>`.

## 7. PDFs

Para cada documento alterado:

```bash
python ferramentas/md2pdf.py docs/<arquivo>.md docs/pdf/<arquivo>.pdf
```

O Markdown é a fonte; **o PDF nunca é editado à mão**. A conversão precisa de internet para desenhar os diagramas. Em diagrama de sequência, não use ponto e vírgula nem cerquilha no texto das mensagens.

## 8. Lista de conferência final

- [ ] Documentos afetados atualizados, com versão subida e linha no histórico
- [ ] Citações em outros documentos atualizadas
- [ ] CHANGELOG com a entrada e a referência
- [ ] Registro da etapa escrito, com numeração nova e índice atualizado
- [ ] Entrega consolidada escrita em `docs/entregas/`, se a etapa foi aceita, e notas nas entregas anteriores que ela alterou
- [ ] Sumário atualizado, se o escopo ou o vocabulário mudou
- [ ] Pendências nos dois lugares
- [ ] Evidência visual capturada e citada
- [ ] Verificador com zero erro, saída colada
- [ ] PDFs regerados
- [ ] Portão da etapa verde, com a prova de cada item

## 9. Commit

**Só quando o dono do produto pedir** ([CLAUDE.md](../../../CLAUDE.md), seção 6). Mensagem em português, citando os IDs afetados, no padrão de commits convencionais.

O commit vai para o **branch da etapa** (`etapa/eNN-titulo-curto`), nunca direto para `main`. No fechamento, abra o pull request para `main` com o resumo da etapa e o link do registro; a integração só acontece com o CI verde e a pedido do dono do produto ([docs/execucao/1_MODELO_EXECUCAO.md](../../../docs/execucao/1_MODELO_EXECUCAO.md), seção 4.1).
