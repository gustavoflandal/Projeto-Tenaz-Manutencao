# Registro de etapas do projeto

Toda etapa do projeto Tenaz é documentada nesta pasta: pesquisa, decisão, mudança de escopo, implementação, correção, release e revisão. É a memória do projeto: permite saber, a qualquer momento, o que foi feito, por quê, por quem e o que ficou pendente. A regra está em [CLAUDE.md](../../CLAUDE.md), seção 3.

## Regras

1. Um arquivo por etapa: `NNN_AAAA-MM-DD_titulo-curto.md`. `NNN` é sequencial, com três dígitos, e nunca é reaproveitado. Título curto em minúsculas, sem acento, com hífens.
2. O registro é criado **no mesmo commit** das mudanças que descreve.
3. Registro publicado não é reescrito. Erro se corrige com uma nota "Correção em AAAA-MM-DD" no fim do próprio registro, ou com uma etapa nova que o cite.
4. Toda mudança citada aponta para o documento e o ID afetados, com link relativo.
5. Pendências abertas na etapa também vão para a seção de pendências do documento dono ou do sumário executivo. A etapa que as resolver cita a etapa que as abriu.
6. A tabela de índice abaixo é atualizada a cada etapa nova.

## Modelo

```markdown
# NNN — Título da etapa

**Data:** AAAA-MM-DD · **Tipo:** pesquisa | decisão | escopo | implementação | correção | release | revisão
**Pedido por:** dono do produto | equipe · **Versão do produto:** — ou vX.Y.Z

## Pedido
O que foi solicitado, nas palavras de quem pediu (resumo fiel).

## Decisões
Cada decisão com o motivo. Alternativas descartadas, quando houver.

## Mudanças
| Documento ou arquivo | O que mudou | IDs afetados |
|---|---|---|

## Verificação
Como foi conferido (verificador de documentação, testes, revisão, PDFs gerados).

## Pendências
O que ficou em aberto e onde está registrado.
```

## Índice

| # | Data | Tipo | Etapa |
|---|---|---|---|
| 001 | 2026-11-12 | decisão | [Nome e identidade visual](001_2026-11-12_nome-e-identidade-visual.md) |
| 002 | 2026-11-12 | pesquisa | [Estudo de mercado e benchmarking](002_2026-11-12_estudo-de-mercado.md) |
| 003 | 2026-11-12 | escopo | [Documentação base do projeto](003_2026-11-12_documentacao-base.md) |
| 004 | 2026-11-12 | escopo | [Terceiros, estoque estendido e locação](004_2026-11-12_terceiros-estoque-locacao.md) |
| 005 | 2026-11-12 | escopo | [Versão visível e histórico de versões](005_2026-11-12_versao-e-changelog.md) |
| 006 | 2026-11-12 | escopo | [Gestão de serviços do prestador e alinhamento 0.4](006_2026-11-12_gestao-servicos-prestador.md) |
| 007 | 2026-11-12 | release | [Primeiro commit no GitHub](007_2026-11-12_primeiro-commit.md) |
| 008 | 2026-11-12 | decisão | [Prompt inicial e regras de documentação](008_2026-11-12_prompt-e-regras-de-documentacao.md) |
| 009 | 2026-11-12 | decisão | [Plano de execução, equipe de agentes e skills](009_2026-11-12_equipe-de-agentes-e-skills.md) |
