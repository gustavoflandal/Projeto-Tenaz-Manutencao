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
| Documentação do projeto | Em elaboração ([índice](docs/README.md)) |

## Estrutura

```
marca/          identidade visual: brandbook, logos SVG e PNG
docs/           documentação do projeto
docs/pesquisa/  estudo de mercado: notas de pesquisa e relatório final
ferramentas/    utilitários (conversão de Markdown para PDF)
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
