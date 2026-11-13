# Tenaz — instruções do projeto

Este arquivo é o prompt inicial do projeto. Vale para toda sessão de trabalho, humana ou de IA, neste repositório. Leia-o inteiro antes de começar.

## 1. O que é a Tenaz

Sistema SaaS multiempresa de gestão de manutenção de equipamentos e frota para o mercado brasileiro: transportadoras, indústrias, comércio, serviços, locadoras e subcontratadas de manutenção. Slogan: **Cada ativo pronto para operar.**

- Estado atual: concepção concluída (marca, estudo de mercado, documentação). Ainda não há código.
- Repositório: `github.com/gustavoflandal/Projeto-Tenaz-Manutencao`, branch principal `main`.
- Idioma de tudo (interface, documentação, commits, código de domínio): **português do Brasil**.

## 2. Fontes da verdade

Em caso de conflito, vale esta ordem:

1. Pedido explícito do dono do produto na sessão.
2. [docs/0_SUMARIO_EXECUTIVO.md](docs/0_SUMARIO_EXECUTIVO.md): módulos, planos, **vocabulário canônico** (nomes de entidades, campos, estados), decisões (seções 3.3 a 3.6) e pendências.
3. O documento dono de cada assunto:

| Assunto | Documento dono |
|---|---|
| Requisitos `RF-<MÓD>-NNN` e `RNF-NNN` | [docs/1_ESPECIFICACAO_REQUISITOS.md](docs/1_ESPECIFICACAO_REQUISITOS.md) |
| Modelo de dados, SQL, isolamento, sincronização | [docs/2_ARQUITETURA_BANCO_DADOS.md](docs/2_ARQUITETURA_BANCO_DADOS.md) |
| Endpoints, eventos de webhook, códigos de erro | [docs/3_ESPECIFICACAO_APIS.md](docs/3_ESPECIFICACAO_APIS.md) |
| Fluxos e regras `RN-<MÓD>-NNN` | [docs/4_FLUXOS_PROCESSO.md](docs/4_FLUXOS_PROCESSO.md) |
| Arquitetura, ADRs, versionamento, testes | [docs/5_ARQUITETURA_SOFTWARE.md](docs/5_ARQUITETURA_SOFTWARE.md) |
| Prazos, sprints, equipe, riscos de cronograma | [docs/6_ROADMAP.md](docs/6_ROADMAP.md) |
| Interface, tokens, componentes, contrastes | [docs/7_PADROES_DESIGN.md](docs/7_PADROES_DESIGN.md) e [marca/README.md](marca/README.md) |
| Segurança, LGPD, SLA, recuperação de desastre | [docs/9_SEGURANCA_CONFIABILIDADE.md](docs/9_SEGURANCA_CONFIABILIDADE.md) |
| Histórico de versões do produto | [CHANGELOG.md](CHANGELOG.md) |
| Histórico de etapas do projeto | [docs/etapas/](docs/etapas/) |
| Entregas consolidadas (etapas aceitas e em `main`) | [docs/entregas/](docs/entregas/README.md) |
| Como o trabalho acontece: agentes, skills, etapas | [docs/execucao/](docs/execucao/README.md) |

Um ID (`RF-`, `RNF-`, `RN-`, `ADR-`) só é **criado** no documento dono. Os outros documentos apenas o citam.

## 3. Regras obrigatórias de documentação

### 3.1 Toda etapa é registrada

Toda etapa do projeto (pesquisa, decisão, mudança de escopo, implementação, correção, release, revisão) gera um registro em `docs/etapas/`, no formato definido em [docs/etapas/README.md](docs/etapas/README.md): `NNN_AAAA-MM-DD_titulo-curto.md`, numeração sequencial, nunca reaproveitada. O registro diz o que foi pedido, o que foi decidido e por quê, o que mudou (arquivos e IDs), como foi verificado e o que ficou pendente.

### 3.2 Nenhuma mudança sem documentação atualizada

Toda correção ou implementação atualiza, **no mesmo commit ou pull request**:

1. Os documentos afetados em `docs/` (requisito, modelo de dados, API, fluxo, regra, design, segurança, roadmap), subindo a versão e acrescentando uma linha no histórico do documento.
2. O [CHANGELOG.md](CHANGELOG.md), na seção `[Não lançado]`, com a referência ao requisito (`RF-`/`RNF-`) ou ao chamado da correção (RNF-040).
3. O registro da etapa em `docs/etapas/`.
4. O sumário executivo, quando a mudança afetar módulo, plano, vocabulário, estado ou pendência.

Uma mudança de comportamento que não aparece na documentação é tratada como defeito.

**A documentação é idempotente ao código** (pedido do dono do produto, etapa E01): em nenhum momento os dois podem divergir. Todo fix, adendo, inclusão, exclusão ou alteração é documentado no mesmo passo em que acontece, e não só no fechamento da etapa: o documento dono, o `CHANGELOG.md` e o registro da etapa em andamento acompanham cada mudança. Agente que altera código sem atualizar o documento dono não terminou a tarefa.

### 3.3 Integridade e referência cruzada

- Use exatamente os nomes canônicos do sumário (entidades, campos, estados, permissões `<módulo>.<ação>`, rotas `/v1/...`, eventos `<entidade>.<fato>`). Nome novo entra primeiro no sumário.
- Cite sempre por ID e por link relativo: `RF-OS-012 em [1_ESPECIFICACAO_REQUISITOS.md](docs/1_ESPECIFICACAO_REQUISITOS.md)`.
- IDs não são renumerados nem reaproveitados. Requisito ou regra que muda de módulo ou é abandonado fica com nota "movido para" ou "substituída por".
- Ao mudar um assunto, procure todas as citações dele nos outros documentos e atualize-as junto.
- Divergência encontrada e não resolvida vai para a seção de inconsistências/pendências do documento e para o registro da etapa. Nunca fica só na conversa.
- Antes de concluir qualquer etapa, rode o verificador e corrija o que ele apontar:

```bash
python ferramentas/verificar_docs.py
```

### 3.4 PDFs

Cada documento tem versão em PDF em `docs/pdf/`, gerada a partir do Markdown (o Markdown é a fonte; o PDF nunca é editado à mão):

```bash
python ferramentas/md2pdf.py docs/<arquivo>.md docs/pdf/<arquivo>.pdf
```

Regere o PDF de todo documento alterado. O conversor precisa de internet para desenhar os diagramas Mermaid. Em diagramas de sequência, não use `;` nem `#` no texto das mensagens.

## 4. Decisões já tomadas

- **Produto:** SaaS multiempresa, foco inicial em empresas de 50 a 2.000 ativos; módulos e planos no sumário, seções 2 e 2.2.
- **Stack:** Go (monólito modular) + React/TypeScript (web) + React Native/Expo (app offline-first, SQLite) + PostgreSQL com Row-Level Security; Docker. Ver ADR-001 a ADR-009 no documento 5.
- **Isolamento:** coluna `empresa_id` em toda tabela de negócio, RLS forçado, escopos `interno`, `prestador`, `locatario`, `link`; conexão entre contas (SRV) só pelo processo de sincronização.
- **Versionamento:** SemVer `MAJOR.MINOR.PATCH` vindo da tag `vX.Y.Z`; versão sempre visível no canto inferior direito da tela principal e do login (RNF-039); `CHANGELOG.md` sempre atualizado (RNF-040).
- **Configuração por segmento:** a empresa configura blocos de catálogo fechado (campos extras, regras de condição e efeito, modelos de projeto, contratos com clientes, pacotes de segmento). Estados, fórmulas de indicadores, regras do offline, isolamento e campos da conexão entre contas não são configuráveis. Condição ou efeito novo entra por requisito com teste, nunca por configuração (sumário 3.7, risco R7).
- **Ferramentas:** só software 100% gratuito ou open source no que o sistema usa diretamente em execução; ferramentas do ambiente de desenvolvimento (Docker Desktop, navegador, editor, GitHub) ficam fora da regra. A hospedagem de produção na AWS é exceção explícita. Ver sumário 3.6 e ADR-012.
- **Marca:** cores Aço `#14202E` e Sinal `#FF6A13`; Archivo Expanded, Inter, JetBrains Mono; regras de contraste em [marca/README.md](marca/README.md). Texto sobre Sinal é sempre Aço.

## 5. Pendências que dependem do dono do produto

Consulte a seção 6 do sumário. Principais: P1 (ICP-Brasil na NR-1), P4 (preços), P5 (INPI), P7 (local da cópia de recuperação de desastre), P8 (CFOPs), tamanho do time do MVP. Não decida essas questões sozinho: apresente as opções.

## 6. Como trabalhar

- **Comece pela skill `tenaz-contexto`**, em [.claude/skills/](.claude/skills/). Ela orienta qualquer tarefa: fontes da verdade, documento dono, vocabulário canônico, pendências e qual agente chamar.
- **Use o agente especializado do assunto**, em [.claude/agents/](.claude/agents/). O catálogo, com a divisão de trabalho e quem revisa quem, está em [docs/execucao/2_EQUIPE_AGENTES.md](docs/execucao/2_EQUIPE_AGENTES.md).
- **Toda etapa segue o ciclo de sete passos** de [docs/execucao/1_MODELO_EXECUCAO.md](docs/execucao/1_MODELO_EXECUCAO.md): prompt, plano, execução, verificação, evidência, registro, aceite. O prompt da etapa é escrito **antes** do trabalho, em [docs/prompts/](docs/prompts/).
- Antes de propor ou alterar algo, leia o sumário e o documento dono do assunto.
- Mudança de escopo pedida pelo dono do produto: atualize primeiro o sumário, depois os documentos donos, depois os que citam, e registre a etapa.
- Quando houver código: testes junto com a funcionalidade, inclusive o teste de vazamento entre empresas, prestadores, locatários e contas conectadas (R2, R4, R5).
- Commits em português, mensagem curta no imperativo ou descritiva, citando os IDs afetados.
- Faça commit e push somente quando o dono do produto pedir.

## 7. Estrutura do repositório

```
CLAUDE.md               este prompt inicial
CHANGELOG.md            histórico de versões do produto (RNF-040)
README.md               apresentação do projeto
backend/                Go: cmd/tenaz, internal/plataforma, internal/modulos,
                        db/migracoes, api/openapi.yaml, testes/vazamento
web/                    React + TypeScript (Vite)
app/                    React Native + Expo, offline-first (src/sync: fila de comandos)
infra/docker/           compose de desenvolvimento (PostgreSQL, Valkey, Mailpit, SeaweedFS)
.github/workflows/      esteira de CI
marca/                  identidade visual (brandbook, logos, regras)
docs/                   documentação do projeto (0 a 9)
docs/etapas/            registro de todas as etapas do projeto
docs/entregas/          cada etapa aceita e consolidada: o que ficou pronto
docs/execucao/          plano de execução: modelo de trabalho, agentes, skills, etapas
docs/prompts/           o prompt de cada etapa, escrito antes do trabalho
docs/screenshots/       evidência visual das telas entregues
docs/pdf/               PDFs gerados dos documentos
docs/pesquisa/          estudo de mercado e notas de pesquisa
ferramentas/            md2pdf.py, verificar_docs.py, versao.sh, tarefas.sh
.claude/agents/         os 16 agentes especializados do projeto
.claude/skills/         as 12 skills com os procedimentos obrigatórios
```

## 8. Comandos do dia a dia

Todas as tarefas passam por um script único, para que "passa na minha máquina" e "passa no CI" signifiquem a mesma coisa. O CI chama exatamente estes subcomandos.

```bash
ferramentas/tarefas.sh ajuda         # lista todas as tarefas
ferramentas/tarefas.sh infra         # sobe PostgreSQL e Valkey
ferramentas/tarefas.sh migrar        # aplica as migrações
ferramentas/tarefas.sh conferir      # formato, análise, contrato, tipos, documentação
ferramentas/tarefas.sh tudo          # conferências + todos os testes
ferramentas/tarefas.sh vazamento     # suíte de isolamento (bloqueia merge e deploy)
```

Requisitos locais: Go 1.27+, Node 22.17+, Docker e Python 3.13+. O `pnpm` é chamado por `npx` porque `corepack enable` precisa de privilégio de administrador no Windows; quem tiver o pnpm instalado globalmente pode usá-lo direto.
