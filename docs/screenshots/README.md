# Evidência visual

Capturas de tela da Tenaz, organizadas por etapa. Servem a quatro propósitos: aceite da etapa pelo dono do produto, ilustração da documentação, base do manual de operação (documento 8, reservado) e comparação entre versões.

O procedimento de captura está na skill [`tenaz-evidencia-visual`](../../.claude/skills/tenaz-evidencia-visual/SKILL.md). A regra de quando capturar está em [../execucao/1_MODELO_EXECUCAO.md](../execucao/1_MODELO_EXECUCAO.md), passo 5 do ciclo da etapa.

## Estrutura

```
docs/screenshots/
├── E00_fundacao-tecnica/
├── E01_conta-e-acesso/
├── E02_ativos/
└── ...                      uma pasta por etapa, com o mesmo nome do prompt
```

## Nome do arquivo

```
<superfície>_<tela>_<estado>_<tema>.png
```

| Parte | Valores | Exemplo |
|---|---|---|
| `superfície` | `web`, `app` | `web` |
| `tela` | nome curto em minúsculas, sem acento, com hífens | `ficha-ativo` |
| `estado` | `padrao`, `vazio`, `carregando`, `erro`, `validacao`, ou um estado do domínio | `em-execucao` |
| `tema` | `claro`, `escuro` | `escuro` |

Exemplos:

```
web_lista-ativos_padrao_claro.png
web_lista-ativos_vazio_claro.png
app_execucao-os_em-execucao_escuro.png
app_execucao-os_offline_claro.png
```

## Regras

1. **Uma captura por tela entregue**, no estado padrão, mais uma por estado que o requisito exija (vazio, erro, validação, offline, bloqueado).
2. **Os dois temas** quando a tela existir nos dois.
3. **Massa de dados plausível.** Nunca "teste 1", "aaa", "Lorem ipsum". Use ativos, TAGs, números de OS e nomes do mesmo tipo que um cliente teria. A semente de desenvolvimento serve a isso.
4. **Nenhum dado pessoal real.** Nome, CPF, CNPJ, telefone, e-mail, placa e número de documento são fictícios. Captura de piloto com dado real não entra no repositório.
5. **Selo de versão visível** nas telas principais e de login, por RNF-039. Captura que o esconde não serve como evidência.
6. **Largura fixa:** 1440 px no web e 390 px no app, para que versões diferentes sejam comparáveis. Captura de tela inteira, não de recorte, salvo quando o recorte for o assunto.
7. **PNG**, sem moldura de navegador nem de aparelho.
8. **Captura substituída, não acumulada.** Quando a tela muda, a captura da pasta da etapa onde ela nasceu é mantida como histórico e a nova vai para a pasta da etapa que a alterou.
9. **A captura é citada** no registro da etapa em [../etapas/](../etapas/), por link relativo.

## O que não vai aqui

- Diagramas: ficam nos próprios documentos, em Mermaid.
- Imagens da marca: ficam em [../../marca/](../../marca/).
- Gravações de tela e vídeo: não entram no repositório.
