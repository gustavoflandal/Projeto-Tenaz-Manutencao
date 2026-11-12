<div align="center">

<img src="marca/logo/svg/tenaz-horizontal-positivo.svg" alt="Tenaz" width="320">

**Gestão de Manutenção e Frota**

*Cada ativo pronto para operar.*

</div>

---

A Tenaz é um sistema SaaS de gestão de manutenção de equipamentos e frota para empresas de vários segmentos: transportadoras, indústrias, comércio e serviços. Cobre planos preventivos e preditivos, ordens de serviço, estoque de peças, custos, indicadores e aplicativo de campo.

## Estado do projeto

Fase de pesquisa e documentação. Ainda não há código.

| Etapa | Situação |
|---|---|
| Nome e identidade visual | Aprovados em 12/11/2026 |
| Pesquisa de mercado e benchmarking | Concluída em 12/11/2026 ([estudo](docs/pesquisa/reports/)) |
| Documentação do projeto | Versão alinhada ao sumário executivo 0.4 ([índice](docs/README.md)) |
| Registro de etapas | [docs/etapas/](docs/etapas/README.md) |

## Estrutura

```
CLAUDE.md       prompt inicial do projeto: contexto, fontes da verdade e regras obrigatórias
CHANGELOG.md    histórico de versões
marca/          identidade visual: brandbook, logos SVG e PNG
docs/           documentação do projeto (índice em docs/README.md)
docs/etapas/    registro de todas as etapas do projeto
docs/pdf/       versões em PDF dos documentos
docs/pesquisa/  estudo de mercado: notas de pesquisa e relatório final
ferramentas/    conversão de Markdown para PDF e verificador da documentação
```

## Regras de documentação

Toda etapa é registrada em `docs/etapas/`, e toda correção ou implementação atualiza a documentação no mesmo commit. As regras estão em [CLAUDE.md](CLAUDE.md), seção 3. Antes de concluir qualquer etapa:

```bash
python ferramentas/verificar_docs.py
```

## Gerar PDF de um documento

```bash
python ferramentas/md2pdf.py docs/arquivo.md
```

Requer Python com o pacote `markdown` e o Google Chrome ou o Microsoft Edge instalado.

## Premissas adotadas

- Entrega como SaaS multiempresa.
- Foco inicial em empresas de médio porte (50 a 2.000 ativos), com arquitetura preparada para grande porte.
- Stack: Go no backend (monólito modular), React + TypeScript no web, React Native (Expo) + TypeScript no app de campo, PostgreSQL e Docker. A decisão e as alternativas avaliadas estão em `docs/5_ARQUITETURA_SOFTWARE.md`.
