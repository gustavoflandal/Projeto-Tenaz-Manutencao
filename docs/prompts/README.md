# Prompts das etapas

O pedido de cada etapa, escrito antes de a etapa começar e versionado junto com o código. Um prompt diz o que se quer, com que contexto, com que critério de aceite e o que está fora do escopo.

Por que versionar o prompt: meses depois, o código mostra o que foi feito e o registro em [../etapas/](../etapas/) mostra o que mudou, mas só o prompt mostra **o que se pediu e com que informação**. É o que permite separar "a decisão estava errada" de "o pedido estava ambíguo".

O procedimento de escrita está na skill [`tenaz-prompt-etapa`](../../.claude/skills/tenaz-prompt-etapa/SKILL.md).

## Arquivos

| Arquivo | Conteúdo |
|---|---|
| [000_modelo.md](000_modelo.md) | Modelo a copiar |
| `E00` a `E12` | As etapas do MVP, correspondentes às sprints S0 a S12 do [roadmap](../6_ROADMAP.md) |

| # | Prompt | Etapa |
|---|---|---|
| E00 | [E00_fundacao-tecnica.md](E00_fundacao-tecnica.md) | Fundação técnica |
| E01 | [E01_conta-e-acesso.md](E01_conta-e-acesso.md) | Conta e acesso |
| E02 | [E02_ativos.md](E02_ativos.md) | Ativos |
| E03 | [E03_medidores.md](E03_medidores.md) | Medidores e composição |
| E04 | [E04_solicitacoes-e-os.md](E04_solicitacoes-e-os.md) | Solicitações, núcleo da OS e fundação da conexão |
| E05 | [E05_execucao-da-os.md](E05_execucao-da-os.md) | Execução da OS e OS de serviço |
| E06 | [E06_preventiva.md](E06_preventiva.md) | Preventiva |
| E07 | [E07_checklists-estoque-fornecedores.md](E07_checklists-estoque-fornecedores.md) | Checklists, estoque básico e fornecedores |
| E08 | [E08_nr12-documentos-assinatura.md](E08_nr12-documentos-assinatura.md) | NR-12, documentos e assinatura |
| E09 | [E09_indicadores-e-terceiros.md](E09_indicadores-e-terceiros.md) | Indicadores e terceiros |
| E10 | [E10_api-sso-cobranca.md](E10_api-sso-cobranca.md) | API, SSO, cobrança e portão da conexão |
| E11 | [E11_endurecimento.md](E11_endurecimento.md) | Endurecimento |
| E12 | [E12_lancamento.md](E12_lancamento.md) | Lançamento |

## Regras

1. **Nome:** `E<NN>_titulo-curto.md` para etapa do MVP; `F2-M<NN>_titulo.md` e `F3-M<NN>_titulo.md` para as fases 2 e 3. Título em minúsculas, sem acento, com hífens. O nome é o mesmo da pasta de evidência em [../screenshots/](../screenshots/).
2. **O prompt é escrito antes do trabalho.** Prompt escrito depois é relatório, não pedido.
3. **Prompt publicado não é reescrito.** Ajuste de rumo no meio da etapa entra como seção "Ajustes durante a etapa", com data e motivo. É o que preserva a diferença entre o pedido original e o que ele se tornou.
4. **Cite por ID e por link relativo.** Requisitos pelos IDs de [../1_ESPECIFICACAO_REQUISITOS.md](../1_ESPECIFICACAO_REQUISITOS.md); nunca descreva um requisito com outras palavras em vez de citá-lo.
5. **Diga o que está fora.** A seção de fora do escopo evita que a etapa cresça no meio, risco 2 do cronograma no [roadmap](../6_ROADMAP.md), seção 9.
6. **Um portão por etapa, verificável por comando.** "Funcionando bem" não é portão; "suíte de três contas verde no CI" é.
7. **Nenhum segredo no prompt:** sem senha, token, chave ou dado pessoal real.
