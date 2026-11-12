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
- Plano de execução do projeto em `docs/execucao/`: modelo de trabalho com ciclo de sete passos por etapa, catálogo de 16 agentes especializados correspondentes aos papéis do time de 15 pessoas do roadmap, catálogo de 12 skills com os procedimentos obrigatórios, e as etapas E00 a E12 do MVP com entrega, agentes, portão de saída e evidência de cada uma.
- Equipe de agentes especializados em `.claude/agents/` e skills do projeto em `.claude/skills/`, com os portões de isolamento entre empresas (RNF-013), de acesso externo (RNF-037, RNF-038), entre contas conectadas (RNF-041, RNF-042), de operação offline (RNF-010), de acessibilidade (RNF-025) e de versão visível (RNF-039) embutidos nos procedimentos.
- Pasta `docs/prompts/` com o prompt de cada etapa do MVP e o modelo, e pasta `docs/screenshots/` com a convenção de evidência visual das telas entregues.
- **Fundação técnica (etapa E00).** Primeiro código do produto: monólito modular em Go com os subcomandos `api`, `migrate`, `migrate-status`, `openapi` e `versao`; esquemas `app` e `plataforma` com os papéis de banco, as funções de contexto e de escopo, o UUID v7 e o trigger das colunas padrão; contrato OpenAPI 3.1 gerado dos tipos Go e versionado; infraestrutura local em contêiner com PostgreSQL 17 e Valkey; esteira de CI com cinco verificações; projeto da web em Vite e do app em Expo.
- Esteira de versão e histórico (RNF-039 e RNF-040): a versão vem da tag `vX.Y.Z` e é injetada no build do backend, da web e do app, nunca digitada em arquivo de código; `GET /v1/versao` devolve versão, commit curto e data do build e nada mais; produção recusa versão de pré-lançamento; o CI bloqueia o release quando a tag, o CHANGELOG e a versão embutida divergem, e recusa pull request que altere código de produto sem registrar a mudança aqui.
- Suíte de isolamento entre empresas (RNF-013), rodando a cada build e bloqueando merge e deploy: verifica a cobertura de RLS de toda tabela com `empresa_id`, a lista explícita de catálogos globais, a propriedade de todo objeto por `tenaz_dono`, a ausência de contorno de RLS nos papéis da aplicação e a ausência de privilégio amplo do papel de sincronização (RNF-041). Inclui a prova de que o detector funciona, criando de propósito uma tabela sem política e exigindo que seja acusada.
- Fila de comandos (outbox) do app de campo, com testes contra SQLite real: nenhum comando se perde por falha de rede, rejeição ou fechamento do app; reenvio não duplica; a ordem é FIFO por ordem de serviço, sem uma OS bloquear outra; comando rejeitado guarda o motivo e aparece em "Pendências de sincronização"; cenário de sete dias de campo acumulado convergindo na reconexão (RF-MOB-002, RNF-010).
- Verificador de documentação passa a rodar no CI, resolvendo a pendência aberta na etapa 008.
- **Configuração por segmento** (escopo, sumário 0.5, seção 3.7). A Tenaz se ajusta a segmentos diferentes por blocos de catálogo fechado, sem estados nem fluxos livres. No MVP: regras de condição e efeito, como "verificação metrológica vencida bloqueia o ativo até a nova verificação" (RF-ADM-016); fase de ciclo de vida do ativo, com "Em implantação" e a ação "Colocar em operação" (RF-ATV-013); campos extras na OS, como números de lacre (RF-OS-021); tipo de OS `implantacao`, fora dos indicadores de manutenção. Na Fase 2: pacotes de segmento prontos, como o de fiscalização eletrônica (RF-ADM-017); módulo PRJ de projetos de implantação com etapas, dependências, cronograma, avanço e medição de etapas para pagamento do prestador (RF-PRJ-001 a RF-PRJ-007); módulo CTR de contratos com clientes, com meta de disponibilidade, faixas de desconto, paradas por causa externa e medição mensal (RF-CTR-001 a RF-CTR-005).

### Corrigido

- Divergência de nomes entre os documentos 2 e 5: as funções de contexto e de escopo eram citadas no documento 5 com o esquema `tenaz` e o nome `escopo_tipo`, e no documento 2, dono do assunto, com o esquema `app` e `escopo_atual`. O documento 5 passou a usar os nomes do documento 2. Encontrada ao implementar a fundação técnica.
- Prompt inicial do projeto em `CLAUDE.md`, com as regras obrigatórias de documentação.
- Registro de etapas do projeto em `docs/etapas/`, com as etapas 001 a 008.
- Verificador de integridade e referências cruzadas da documentação em `ferramentas/verificar_docs.py` (links, IDs RF/RNF/RN/ADR, versões, registro de etapas, PDFs desatualizados).

### Alterado

- Regras de acessibilidade aplicadas à marca: texto Aço sobre o laranja Sinal, contorno de foco em Aço, chip de status com texto Aço no tema claro e texto terciário em Aço 400.
