---
name: tenaz-evidencia-visual
description: Use para capturar, nomear e guardar screenshots das telas da Tenaz em docs/screenshots. Dispara ao terminar qualquer tela nova ou alterada, e sempre que uma etapa precisar de prova visual para o aceite do dono do produto.
---

# Evidência visual da Tenaz

A captura serve a quatro coisas: o aceite da etapa pelo dono do produto, a ilustração da documentação, a base do manual de operação (documento 8, reservado) e a comparação entre versões. Captura descuidada não serve a nenhuma das quatro.

## 1. O que capturar

| Sempre | Além disso |
|---|---|
| Estado padrão de toda tela entregue | Um por estado que o requisito exija: vazio, erro, validação, offline, bloqueado, somente leitura, sem permissão |
| Nos dois temas, quando a tela existir nos dois | Telas de portal restrito do prestador e do locatário, que são superfícies próprias |
| Telas principais com o selo de versão visível | Fluxo com mais de um passo: uma captura por passo |

## 2. Preparar a massa de dados

A captura é tão boa quanto os dados nela:

1. **Dados plausíveis.** Nunca "teste 1", "aaa", "Lorem ipsum", "asdf". Use ativos, TAGs, números de OS, nomes e datas do tipo que um cliente teria. A semente de desenvolvimento serve exatamente a isso.
2. **Nenhum dado pessoal real.** Nome, CPF, CNPJ, telefone, e-mail, placa e número de documento são fictícios. **Captura de piloto com dado real não entra no repositório.**
3. **Volume realista na lista.** Lista com uma linha não mostra o problema da lista.
4. **Estado coerente.** Uma OS `em_execucao` tem técnico, início e apontamento. Captura com estado incoerente confunde quem revisa.

## 3. Capturar

| Item | Regra |
|---|---|
| Largura | 1440 px na web, 390 px no app. Fixas, para que versões diferentes sejam comparáveis |
| Enquadramento | Tela inteira, sem moldura de navegador nem de aparelho. Recorte só quando o recorte for o assunto |
| Formato | PNG |
| Tema | Claro e escuro, quando a tela existir nos dois |
| Cursor e foco | Sem cursor sobre elemento, salvo quando o foco é o assunto |

## 4. Nomear

```
<superfície>_<tela>_<estado>_<tema>.png
```

| Parte | Valores |
|---|---|
| `superfície` | `web`, `app` |
| `tela` | nome curto em minúsculas, sem acento, com hífens |
| `estado` | `padrao`, `vazio`, `carregando`, `erro`, `validacao`, `offline`, ou um estado do domínio |
| `tema` | `claro`, `escuro` |

Exemplos:

```
web_lista-ativos_padrao_claro.png
web_lista-ativos_vazio_claro.png
app_execucao-os_em-execucao_escuro.png
app_execucao-os_offline_claro.png
```

## 5. Guardar

1. Uma pasta por etapa, com o **mesmo nome do prompt** da etapa: `docs/screenshots/E02_ativos/`.
2. Quando a tela muda numa etapa posterior, a captura antiga **fica onde está**, como histórico, e a nova vai para a pasta da etapa que a alterou.
3. Cite as capturas no registro da etapa em [docs/etapas/](../../../docs/etapas/), por link relativo.

As regras completas da pasta estão em [docs/screenshots/README.md](../../../docs/screenshots/README.md).

## 6. O que não vai para a pasta

- Diagramas: ficam nos documentos, em Mermaid.
- Imagens da marca: ficam em [marca/](../../../marca/).
- Gravação de tela e vídeo: não entram no repositório.
- Captura com dado real de piloto.

## 7. Verificação

| Verificação | Como |
|---|---|
| Selo de versão visível nas telas principais | Olhe o canto inferior direito. Encoberto reprova a captura (RNF-039) |
| Nenhum dado pessoal real | Revisão da imagem antes de salvar |
| Nome no padrão | Confira cada parte do nome |
| Estados exigidos presentes | Compare com a lista de estados do requisito |
