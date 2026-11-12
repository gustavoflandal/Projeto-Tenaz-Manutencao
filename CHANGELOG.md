# Histórico de versões

Todas as implementações e correções da Tenaz são registradas neste arquivo, versão a versão.

O formato segue o [Keep a Changelog 1.1.0](https://keepachangelog.com/pt-BR/1.1.0/) e as versões seguem o [Versionamento Semântico 2.0.0](https://semver.org/lang/pt-BR/). As regras de uso estão nos requisitos RNF-039 e RNF-040 de [docs/1_ESPECIFICACAO_REQUISITOS.md](docs/1_ESPECIFICACAO_REQUISITOS.md).

## Como registrar

- Toda mudança de comportamento entra em `[Não lançado]` no mesmo pull request que a implementa.
- Grupos: **Adicionado**, **Alterado**, **Corrigido**, **Removido**, **Obsoleto**, **Segurança**.
- Cada item descreve a mudança em linguagem de usuário e cita a referência: `RF-...`/`RNF-...` para implementação, número do chamado ou da issue para correção.
- Quebra de compatibilidade da API pública é marcada com **QUEBRA** e exige versão MAJOR.
- No release, `[Não lançado]` vira `## [X.Y.Z] - AAAA-MM-DD`, igual à tag `vX.Y.Z`.
- Entradas publicadas não são reescritas; erros são corrigidos com nota na própria entrada.

## [Não lançado]

Fase de concepção. Ainda não há versão do software publicada.

### Adicionado

- Nome e identidade visual da Tenaz: símbolo, logotipo, cores, tipografia, brandbook e arquivos SVG e PNG em `marca/`.
- Estudo de mercado e benchmarking de sistemas de manutenção e frota em `docs/pesquisa/`.
- Documentação do projeto em `docs/`: sumário executivo, requisitos, banco de dados, APIs, fluxos e regras de negócio, arquitetura, roadmap, padrões de design, segurança e confiabilidade.
- Módulos de terceiros e assistência técnica (TER) e de locação (LOC), estoque em poder de subcontratadas, peças serializadas e reparáveis com envio para conserto.
- Módulo de gestão de serviços do prestador (SRV, "Tenaz Serviços"): conta própria da subcontratada com cadastro de técnicos, clientes e equipamentos de clientes, OS vinculada a equipamento e técnico com serviço a realizar, componentes substituídos e horas, app do técnico e conexão com as contas de clientes que usam a Tenaz.
- Requisito de versão visível no canto inferior direito da tela principal (RNF-039) e deste histórico de versões (RNF-040).
- Ferramenta de conversão de documentos Markdown para PDF com a identidade da marca em `ferramentas/md2pdf.py`.

### Alterado

- Regras de acessibilidade aplicadas à marca: texto Aço sobre o laranja Sinal, contorno de foco em Aço, chip de status com texto Aço no tema claro e texto terciário em Aço 400.
