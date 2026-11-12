# 011 — Fundação técnica (etapa E00)

**Data:** 2026-11-12 · **Tipo:** implementação
**Pedido por:** dono do produto · **Versão do produto:** —

## Pedido
Iniciar o projeto. Executada a etapa **E00 — Fundação técnica**, conforme [docs/prompts/E00_fundacao-tecnica.md](../prompts/E00_fundacao-tecnica.md). Durante a etapa, o dono do produto levantou o Docker e pediu que a infraestrutura do sistema fosse hospedada nele.

## Decisões

- **Primeiro código do produto.** Até aqui o repositório só tinha documentação. Entraram `backend/`, `web/`, `app/`, `infra/docker/` e `.github/workflows/`, na estrutura da seção 7 de [5_ARQUITETURA_SOFTWARE.md](../5_ARQUITETURA_SOFTWARE.md).
- **Migração dividida em duas**, e não uma só: a `00001_bootstrap` cria papéis, extensões, esquemas e privilégios padrão rodando como o papel migrador; a `00002_funcoes_base` começa com `SET ROLE tenaz_dono` e cria as funções. O motivo é de segurança, não de organização: `ALTER DEFAULT PRIVILEGES FOR ROLE tenaz_dono` só alcança objetos criados **por** `tenaz_dono`, e `FORCE ROW LEVEL SECURITY` **não vale para o dono da tabela** — uma tabela de dono superusuário contornaria o RLS silenciosamente. A primeira tentativa deixou as 11 funções com dono `tenaz_migrador`, e a suíte passou a cobrar isso.
- **Tabela de controle do goose passa para `tenaz_dono`** na migração de bootstrap, uma única vez. A alternativa era exigir `RESET ROLE` no fim de cada migração futura — disciplina que uma hora falha.
- **`app.conferir_escopo()` entra em duas etapas.** Agora valida o valor do escopo, exige empresa e exige o segredo do link. A conferência contra `app.usuario` entra com a migração de ADM, quando a tabela existir (RF-ADM-004, RF-ADM-005), por `CREATE OR REPLACE`.
- **`pg_partman` sai da base.** A seção 4.1 do documento 2 a põe entre as extensões iniciais, mas ela não acompanha a imagem oficial do PostgreSQL. Entra na migração que criar a primeira tabela particionada (`leitura_medidor`). Pendência registrada abaixo.
- **Suíte de isolamento com prova do próprio detector.** Em E00 não existe tabela de negócio, então uma verificação de cobertura passaria por vacuidade — e passaria também se a consulta estivesse errada. O teste `TestDetectorAcusaTabelaSemPolitica` cria de propósito uma tabela com `empresa_id` e sem política e exige que a verificação a acuse.
- **`info.version` do contrato OpenAPI fica fixo em `"1"`**, a versão do caminho `/v1`, e não a versão do build. Com a versão do build, o contrato gerado mudaria a cada compilação e a verificação de divergência do CI acusaria falso positivo em todo pull request. A versão do build é exposta por `GET /v1/versao`.
- **Modelo de erro próprio no contrato**, com o campo `codigo` estável da seção 3 do documento 3, em `application/problem+json`, inclusive em rota desconhecida. O modelo padrão da biblioteca não tem `codigo`.
- **Orquestração por script, não por `make`.** A seção 7 do documento 5 decide "scripts do pnpm"; `make` não existe por padrão no Windows. Um `Makefile` foi escrito e descartado na mesma etapa, em favor de [ferramentas/tarefas.sh](../../ferramentas/tarefas.sh), que o CI chama nos mesmos subcomandos. Duas fontes de verdade para tarefas é pior que uma.
- **Fila de comandos testada contra SQLite real no Node**, com a interface de depósito separada, e adaptador para `expo-sqlite` no aparelho. O motor é o mesmo nos dois: se o teste e o aparelho usassem motores diferentes, o teste não provaria nada sobre o campo.
- **MinIO atrás de um perfil do compose.** A MinIO deixou de publicar a tag `latest` e nem o Docker Hub nem o quay.io aceitam pull anônimo dela. Anexos só aparecem a partir de E02, então o serviço não sobe no `up` comum. Pendência registrada.
- **`pnpm` chamado por `npx`.** `corepack enable` exige privilégio de administrador no Windows e falhou com `EPERM`. No CI, que é Linux, o `corepack enable` é usado normalmente.
- Alternativa descartada: adotar o PowerSync para a leitura do app. Mantida a recomendação do [ADR-004](../5_ARQUITETURA_SOFTWARE.md): sincronização própria, porque a parte difícil (escrita, idempotência e conflito) fica com a Tenaz nos dois casos e o PowerSync criaria um segundo caminho de acesso aos dados fora do RLS, que é o risco R2.

## Mudanças

| Documento ou arquivo | O que mudou | IDs afetados |
|---|---|---|
| `backend/cmd/tenaz/` | Comando com os subcomandos `api`, `migrate`, `migrate-status`, `openapi`, `versao` | RNF-039 |
| `backend/internal/plataforma/versao/` | Versão injetada pelo build, commit conferido contra o que o Go grava no binário, recusa de pré-lançamento em produção | RNF-039 |
| `backend/internal/plataforma/{config,obs,erros,db,api,migracao}/` | Configuração, log estruturado com redação, erros RFC 9457 com código estável, acesso ao banco com `NaEmpresa`, roteador e migrador | RNF-013, RNF-016, RNF-020, RNF-032 |
| `backend/db/migracoes/00001_bootstrap.sql` | Extensões, papéis, esquemas `app` e `plataforma`, privilégios padrão | RNF-013 |
| `backend/db/migracoes/00002_funcoes_base.sql` | Funções de contexto e de escopo, `app.uuid_v7()`, `app.tg_padrao()`, `app.conferir_escopo()` | RNF-013, RNF-037 |
| `backend/api/openapi.yaml` | Contrato gerado e versionado, com `GET /v1/versao` e o modelo de erro da Tenaz | RF-INT-001 |
| `backend/testes/vazamento/` | Suíte de isolamento, lista de catálogos globais e prova do detector | RNF-013, RNF-037, RNF-041 |
| `backend/{Dockerfile,sqlc.yaml,.golangci.yml}` | Imagem distroless com versão conferida no build, geração de código e analisadores | — |
| `infra/docker/compose.yaml` | PostgreSQL 17 e Valkey 8 locais; MinIO atrás do perfil `anexos` | — |
| `.github/workflows/ci.yml` | Cinco verificações: documentação, CHANGELOG, backend, versão e JavaScript | RNF-039, RNF-040 |
| `ferramentas/versao.sh`, `ferramentas/tarefas.sh` | Derivação da versão a partir da tag e orquestração das tarefas | RNF-039 |
| `app/src/sync/` | Fila de comandos com depósito em SQLite, adaptador do Expo e 14 testes | RF-MOB-002, RF-MOB-003, RNF-010 |
| `app/app.config.ts` | Versão do release, número de build e versão de runtime vindos do CI | RNF-039 |
| `web/` | Projeto Vite com injeção de versão e `versao.json` publicado no build | RNF-039 |
| [5_ARQUITETURA_SOFTWARE.md](../5_ARQUITETURA_SOFTWARE.md) | Versão 0.4: nomes das funções de contexto alinhados ao documento 2 | — |
| [CLAUDE.md](../../CLAUDE.md) | Estrutura do repositório e seção 8, com os comandos do dia a dia | — |
| [CHANGELOG.md](../../CHANGELOG.md) | Entradas em Adicionado e Corrigido | — |

## Verificação

Todos os comandos abaixo foram executados, com a saída conferida.

| Verificação | Resultado |
|---|---|
| `ferramentas/tarefas.sh conferir` | formato, análise estática, contrato e documentação: **tudo verde** |
| Suíte de isolamento, banco local | **7 testes de catálogo + 5 de falha fechada, todos verdes** |
| Suíte de isolamento, contêiner descartável (modo do CI) | **verde em 7,6 s** |
| Testes de JavaScript | **16 verdes**: 14 da fila de comandos, 2 da web |
| Tipos de JavaScript | web e app **sem erro** |
| `python ferramentas/verificar_docs.py` | **0 erros, 0 avisos** |
| Propriedade dos objetos no banco | 11 funções e 2 esquemas, todos com dono `tenaz_dono` |
| Papéis de banco | só `tenaz_migrador` com contorno de RLS; `tenaz_app`, `tenaz_job`, `tenaz_leitura`, `tenaz_sinc` e `tenaz_sinc_conexao` sem |
| `GET /v1/versao` | devolve exatamente `versao`, `commit` e `data_build` |
| Rota desconhecida | 404 em `application/problem+json`, com `codigo` |
| Produção com versão de pré-lançamento | recusada na partida |
| Contrato gerado duas vezes | idêntico: a verificação de divergência do CI não dá falso positivo |
| Build da web com a versão da tag | `dist/versao.json` com a versão, o commit e o instante do build |

Defeitos encontrados pelos próprios testes e corrigidos nesta etapa:

1. O filtro de redação do log não pegava `X-Api-Key`, porque a lista registra `api_key` e o cabeçalho usa hífen. Passou a normalizar hífen, ponto e espaço.
2. As 11 funções ficaram com dono `tenaz_migrador` em vez de `tenaz_dono`, com a consequência de segurança descrita nas decisões.
3. `ferramentas/versao.sh` considerava limpa uma árvore com arquivo não rastreado, enquanto o Go a considera suja, e a conferência de coerência acusava divergência sem motivo aparente. Passou a usar `git status --porcelain`.
4. O campo `$schema` aparecia nas respostas, fora do que a especificação prevê. Vinha de um gancho de criação da biblioteca, não de configuração — descoberto lendo a fonte, depois de duas tentativas erradas.
5. Erro saía como `application/json` em vez de `application/problem+json`, e rota desconhecida saía como texto puro.
6. Dois dos meus próprios testes afirmavam mais do que a especificação garante: que `app.escopo_interno()` devolveria `false` sem escopo (devolve NULL, que também falha fechado) e que o `app.uuid_v7()` seria estritamente ordenado (é ordenado por milissegundo; dentro do mesmo milissegundo os bits são aleatórios). Os dois foram corrigidos para afirmar a propriedade real, e um teste novo confirma no banco que predicado nulo não libera linha.
7. **Teste instável.** A verificação de variação do prefixo de tempo do UUID v7 usava `generate_series(...) , pg_sleep(...)`, em que a pausa é avaliada uma única vez para o cruzamento inteiro: as 200 amostras saíam no mesmo milissegundo e o resultado passava a depender da velocidade da máquina. Passou na primeira execução, contra o banco local, e falhou na execução em contêiner. A pausa foi movida para fora da consulta, entre três consultas separadas, e o teste foi executado três vezes seguidas para confirmar a estabilidade. Teste que passa às vezes é pior que teste ausente, porque ensina a ignorar vermelho.

## Pendências

- **Infraestrutura de nuvem e build de loja passam para a fase de produção**, por decisão do dono do produto tomada ao fim desta etapa: o trabalho segue em ambiente de **desenvolvimento**, sobre o compose local de [infra/docker/](../../infra/docker/). `infra/tofu/` não é escrito agora. O ajuste está registrado em [docs/prompts/E00_fundacao-tecnica.md](../prompts/E00_fundacao-tecnica.md), seção "Ajustes durante a etapa". Nenhum outro item do portão foi afrouxado.
- ~~**O CI foi escrito mas ainda não executou no GitHub.**~~ Resolvido; ver a correção abaixo.
- **`pg_partman`:** a seção 4.1 do documento 2 a lista entre as extensões iniciais, mas ela não acompanha a imagem oficial do PostgreSQL. Decidir, antes de `leitura_medidor`, entre imagem própria de desenvolvimento, extensão oferecida pelo provedor gerenciado ou particionamento sem a extensão. Dono: `tenaz-plataforma-sre` com `tenaz-banco-dados`.
- **MinIO:** fixar uma release concreta da imagem antes de o módulo de anexos entrar, ou trocar por outro serviço compatível com S3 em desenvolvimento.
- **`app.escopo_id()`:** o documento 5 passou a citar esta função depois do alinhamento de nomes, mas o documento 2 define `app.fornecedor_atual()` e `app.locatario_atual()` separadamente, sem uma função única de identificador de escopo. A função foi citada e **não** existe na migração. Fechar o nome no documento 2 antes da política `escopo_externo` de E09. Dono: `tenaz-banco-dados` com `tenaz-isolamento`.
- **Telas:** nenhuma tela de negócio existe ainda, então não há captura em [docs/screenshots/](../screenshots/). A pasta `E00_fundacao-tecnica/` recebe a evidência do pipeline verde e do app instalado quando as duas pendências de conta externa estiverem resolvidas.
- **P7** (local da cópia de recuperação de desastre) segue aberta e afeta a infraestrutura de backup. Prazo do roadmap: semana 20.

---

## Correção em 2026-11-12 — primeira execução do CI

A pendência "o CI foi escrito mas ainda não executou no GitHub" foi resolvida no mesmo dia, depois do push.

**A primeira execução falhou**, com `Permission denied` em todos os passos que chamam `ferramentas/tarefas.sh` e `ferramentas/versao.sh`. Causa: o `chmod +x` foi feito na máquina de desenvolvimento, mas o Git no Windows não rastreia o bit de executável por padrão, e os dois arquivos foram commitados como `100644`. É uma falha que nenhuma verificação local poderia pegar, porque localmente o bit existe no sistema de arquivos.

Correção, em três camadas:

1. Bit de executável acrescentado ao índice, de `100644` para `100755`.
2. Os passos do CI passam a invocar os scripts por `bash`, o que os torna imunes ao modo do arquivo — a camada que impede o problema de voltar.
3. Verificação nova no CI recusa qualquer `.sh` sem o bit no índice, para a regressão não voltar em silêncio.

**Resultado da segunda execução:** sucesso nos quatro trabalhos que rodam em push — esteira de versão (RNF-039), documentação, backend com a suíte de isolamento, e web e app com a fila offline. O trabalho do CHANGELOG (RNF-040) é pulado em push, por desenho: ele só se aplica a pull request.

Lição que fica para as etapas seguintes: a esteira só está provada depois de executar no CI. "Verde na minha máquina" não é verde.
