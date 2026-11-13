# E00 — Fundação técnica

**Etapa:** E00 · **Sprint do roadmap:** S0 · **Marco:** M0 · **Consolidada em:** 2026-11-12
**Commits:** 294d7f1 (fundação), fc9cbbc (correção do CI), bd62bee (registro da primeira execução do CI)
**Registro da etapa:** [011](../etapas/011_2026-11-12_fundacao-tecnica.md) · **Prompt:** [E00](../prompts/E00_fundacao-tecnica.md)

## O que foi entregue

O repositório saiu da documentação pura e passou a ter o esqueleto do produto, pronto para receber requisitos funcionais. Nenhuma tela de negócio nem requisito funcional entrou nesta etapa, por desenho.

- **Backend em Go**, monólito modular, com os subcomandos `api`, `migrate`, `migrate-status`, `openapi` e `versao`. Configuração, log estruturado com redação de segredos e erros no formato RFC 9457 com código estável.
- **Banco PostgreSQL com isolamento por empresa** desde a primeira migração: esquemas `app` (dados de empresa, sempre com RLS forçado) e `plataforma` (dados globais), papéis de banco separados por função, funções de contexto e de escopo, UUID v7 e o trigger das colunas padrão. Todo objeto tem dono `tenaz_dono`, para que o RLS forçado valha também para o dono.
- **Acesso ao banco só por `NaEmpresa`**, o helper que fixa a empresa e o escopo dentro da transação e falha fechado sem eles.
- **Contrato OpenAPI 3.1 gerado do código** e versionado em `backend/api/openapi.yaml`, com `GET /v1/versao` e o modelo de erro da Tenaz.
- **Esteira de versão** (RNF-039, RNF-040): a versão vem da tag `vX.Y.Z` e é injetada no build do backend, da web e do app; `GET /v1/versao` devolve versão, commit e data do build; produção recusa versão de pré-lançamento; o CI bloqueia o release quando tag, CHANGELOG e versão embutida divergem, e recusa pull request que altera código de produto sem registrar a mudança no CHANGELOG.
- **Suíte de isolamento entre empresas** (RNF-013) rodando em todo build: cobertura de RLS de toda tabela com `empresa_id`, lista explícita de tabelas globais, dono dos objetos, papéis sem contorno de RLS e a prova de que o próprio detector funciona.
- **App de campo**: projeto Expo com banco local e a fila de comandos offline (RF-MOB-002, RF-MOB-003 em prova de conceito), testada contra SQLite real.
- **Web**: projeto Vite com injeção de versão e `versao.json` publicado no build.
- **Infraestrutura local em contêiner** ([infra/docker/](../../infra/docker/)): PostgreSQL 17 e Valkey 8; MinIO atrás de um perfil.
- **Orquestração de tarefas** por um script único, [ferramentas/tarefas.sh](../../ferramentas/tarefas.sh), que o CI chama com os mesmos subcomandos.
- **CI no GitHub** com cinco verificações: documentação, CHANGELOG, backend com a suíte de isolamento, versão, e web e app com a fila offline.

## Requisitos

| Requisito | Situação nesta entrega |
|---|---|
| RNF-013 (isolamento entre empresas) | Base entregue: esquema, RLS forçado, `NaEmpresa` e suíte de vazamento no CI. Cobertura cresce a cada etapa com as tabelas novas |
| RNF-016 (autenticação e segredos) | Base: configuração e redação de segredos no log. A autenticação em si é da E01 |
| RNF-020 (observabilidade) | Base: log estruturado. Trace e métrica completos entram com a infraestrutura de produção |
| RNF-032 (proteção de dados) | Base em desenvolvimento. Criptografia em repouso e TLS gerenciado entram com a infraestrutura de produção |
| RNF-039 (versão visível e esteira de versão) | Esteira completa: versão da tag, injeção no build, `GET /v1/versao`, bloqueio de release em divergência. O selo na tela entra com as telas da E01 |
| RNF-040 (CHANGELOG) | Verificação no CI que recusa mudança de código de produto sem entrada no CHANGELOG |
| RF-MOB-002, RF-MOB-003 | Prova de conceito da fila de comandos; o cache e a execução offline completos entram nas etapas seguintes |

## Portão de saída

| Verificação | Resultado | Prova |
|---|---|---|
| CI verde de ponta a ponta | **Verde** na segunda execução, depois da correção do bit de executável dos scripts | Commits fc9cbbc e bd62bee; [registro 011](../etapas/011_2026-11-12_fundacao-tecnica.md), "Correção em 2026-11-12" |
| Teste de vazamento entre empresas no CI | **Verde**, com a verificação de cobertura quebrando o build se faltar política | Registro 011, seção Verificação |
| Esteira de versão ativa | **Verde**: release bloqueado em divergência; `GET /v1/versao` responde | Registro 011, seção Verificação |
| Verificador de documentação no CI | **Verde**, com 0 erros | Registro 011, seção Verificação |
| Infraestrutura criada por código, na nuvem | **Deslocado** para a fase de produção, por decisão do dono do produto | [Prompt E00](../prompts/E00_fundacao-tecnica.md), "Ajustes durante a etapa" |
| Build do app instalável | **Deslocado** para a fase de produção, junto com as contas de loja, por decisão do dono do produto | Prompt E00, "Ajustes durante a etapa" |

Os dois itens deslocados não são dívida silenciosa: dependem de conta de nuvem, contas de loja e da pendência **P7** (local da cópia de recuperação de desastre), e estão registrados no prompt e no registro da etapa.

## Evidência

- Execução do CI no GitHub, registrada no [registro 011](../etapas/011_2026-11-12_fundacao-tecnica.md).
- Não há tela de negócio nesta etapa, então não há captura em [../screenshots/](../screenshots/README.md). A captura do pipeline verde e do app instalado com o selo de versão fica para quando as contas externas existirem.

## O que ficou para depois

| Item | Para onde foi |
|---|---|
| Infraestrutura de nuvem (`infra/tofu/`) e build de loja do app | Fase de produção |
| `pg_partman` | Resolvido na E01 para a auditoria, com partições criadas por função própria; `leitura_medidor` decide em E03 |
| Versão fixa da imagem do MinIO | Antes do módulo de anexos (E02) |
| `app.escopo_id()` | Resolvido na E01: a função não é criada (ADR-011, documento 5, seção 6.6) |
| P7 (cópia de recuperação de desastre) | Pendência do dono do produto, prazo da semana 20 do roadmap |

## Notas posteriores

Nenhuma até agora. Mudanças da E01 sobre o código desta entrega (por exemplo, o `.gitattributes` com fim de linha LF, o `status` do erro em 64 bits e os códigos de erro alinhados ao documento 3) serão anotadas aqui quando a E01 for consolidada.
