# Arquitetura de Software — Tenaz

**Versão:** 0.9 · **Data:** 13/11/2026 · **Status:** Em elaboração
**Produto:** Tenaz · Gestão de Manutenção e Frota
Alinhado ao sumário executivo 0.8

Este documento descreve como a Tenaz é construída: stack, contêineres, estrutura do monólito modular, isolamento entre empresas e de acessos externos (prestador, locatário e link de uso único), conexão entre contas (cliente e subcontratada, módulo SRV), repositório, convenções, testes, app offline, telemática, IA, infraestrutura e observabilidade. As decisões estruturais estão registradas como ADRs na seção 15. O alinhamento com o sumário e as pendências abertas estão na seção 16.

**Mudanças da versão 0.3 deste documento:** escopo externo `prestador` no lugar de `fornecedor`; permissão de liberação `os.encerrar`; rotas definitivas do link de uso único (sumário 3.6); módulo SRV (pacote `servicos`) nos diagramas, nas dependências, nos testes e no app; processo de sincronização entre contas (seção 6.8) e [ADR-009](#adr-009--conexão-entre-contas); teste automatizado com três contas (seção 9.3).

**Mudanças da versão 0.4 deste documento:** as funções de contexto e de escopo passam a ser citadas com os nomes do documento dono do assunto, [2_ARQUITETURA_BANCO_DADOS.md](2_ARQUITETURA_BANCO_DADOS.md): esquema `app` em vez de `tenaz`, `app.escopo_atual()` em vez de `tenaz.escopo_tipo()`, e a variável de transação `app.escopo` em vez de `app.escopo_tipo`. A divergência foi encontrada na implementação da fundação técnica (etapa [011](etapas/011_2026-11-12_fundacao-tecnica.md)) e resolvida em favor do documento 2, conforme a ordem de fontes da verdade de [CLAUDE.md](../CLAUDE.md), seção 2. Nenhuma decisão de arquitetura mudou: só a grafia dos nomes citados.

**Complemento na versão 0.5:** os três resíduos daquela troca apontados na seção 16 foram corrigidos — as seções 6.4 e 6.8 passaram a citar `app.empresas_ativas()`, `app.conexao_atual()` e `app.escopo = 'interno'`.

**Mudanças da versão 0.5 deste documento:** configuração por segmento (sumário 0.5, seção 3.7; etapa 012). Novo [ADR-010](#adr-010--configuração-por-segmento-com-catálogos-fechados): núcleo canônico fixo e configuração por blocos de catálogo fechado, regras de condição e efeito avaliadas por evento interno e por rotina diária idempotente na fila River, com registro em `disparo_regra`, campos extras da OS em jsonb validados pela definição (também no app, antes de enfileirar o comando) e pacotes de segmento aplicados por cópia idempotente. Mecânica na nova seção 5.6. Pacotes Go `projetos` (PRJ), `contratos` (CTR) e `configuracao` no mapa de módulos e nas regras de dependência (seções 5.1 e 5.3); eventos citados pelos requisitos novos (seção 5.5); jobs novos na fila (seção 4.3); `campo_extra_os` na lista de tabelas com acesso externo (seção 6.6); repositório (seção 7); fluxos E2E 14 a 17 e verificações novas de isolamento (seções 9.2 e 9.3); cache de configuração, comandos e conflito de campo obrigatório no app offline (seções 10.2, 10.3 e 10.5); entitlements de PRJ e CTR (seção 13.4); alinhamento e pendências novas (seção 16). Fora a linha de versão, a de alinhamento e a frase de abertura da seção 16, nenhum texto das versões anteriores foi reescrito; só acréscimos.

**Mudanças da versão 0.6 deste documento:** etapa E01 (conta e acesso), sumário 0.6. Novo [ADR-011](#adr-011--papel-de-banco-próprio-para-as-credenciais-de-autenticação), que resolve o achado ISO-E01-02 da revisão de isolamento: as tabelas de autenticação do esquema `plataforma` (`identidade`, `acesso_empresa`, `sessao`, `token_renovacao`, `codigo_autorizacao`, `desafio_mfa`, `token_email`) passam a ser alcançadas só pelo papel `tenaz_autenticacao`, com pool exclusivo do pacote de autenticação; o `tenaz_app` e o `tenaz_job` perdem todo privilégio nelas. Mecânica na nova seção 6.9; camadas de defesa (6.1), resolução da sessão no middleware (6.3), credencial por modo de execução (4.1), regra 12 de dependência (5.3), verificações da suíte (9.3) e pooler (13.1). Na mesma passada, a pendência E01-3 do documento 2: `app.escopo_id()` não é criada; a seção 6.6 passa a usar `app.fornecedor_atual()`, `app.locatario_atual()` e `app.link_hash_atual()`, e a política do link confere que o link não está revogado, expirado nem usado. A sessão do app passa a durar 30 dias de inatividade, como em RNF-017 (seção 10.6), e o modelo de identidade do login único no app deixa de ser pendência (seção 10.8). Alinhamento ao sumário 0.6 na seção 16.

**Emenda na versão 0.6 (mesma etapa E01, achado ISO-E01-08):** o ADR-011 passa a fixar também o **dono** das tabelas de credencial. Sem RLS, o dono lê tudo, e com dono `tenaz_dono` o `tenaz_manutencao`, membro dele com herança, lia todas as credenciais. As sete tabelas e as duas funções `SECURITY DEFINER` que as tocam pertencem agora ao papel `tenaz_credencial`, sem membros de runtime (decisões 1 e 9, alternativas, risco residual). Seções 6.1, 6.4, 6.9 e 9.3; pendências E01-14 (`conexao_empresa`) e herança do `tenaz_manutencao` na seção 16. Na mesma emenda, a tabela da seção 9.3 passa a descrever exatamente o que a suíte de isolamento de `backend/testes/vazamento/` verifica: linhas novas "Visões", "Papéis e relações de membro", "Moldes da trilha", "Corrida de atribuição", "Contexto no pool" e "Portão de achados"; linhas "Mapa de privilégios", "Dono das credenciais" e "Ataques do achado ISO-E01-02" reescritas conforme os testes.

**Alinhamento ao código de E01 na versão 0.6:** o layout de módulo validado no ADM vira regra (seção 5.2): casos de uso no pacote raiz, HTTP em `<modulo>/rotas`, composição única em `internal/servidor` (seções 2.5 e 7); o desvio do código sqlc do ADM fica como pendência. A seção 2.5 passa a descrever a geração real do contrato (`tenaz openapi`) e marca `oasdiff` e a documentação pública como ainda não implantados. A seção 8.2 passa ao formato de erro real (códigos `snake_case` do documento 3, `type` em `https://api.tenaz.com.br/erros/<codigo>`, `id_rastreio`, `erros[]`). A seção 9.5 fixa a posição na divergência A-38 (`ambiente` não sai de `GET /v1/versao`).

**Mudanças da versão 0.7 deste documento (13/11/2026, etapa E01, sumário 0.7):** três decisões do dono do produto. (1) **Só ferramentas e recursos 100% gratuitos ou de código aberto** (prompt da etapa E01, "Ajustes durante a etapa", linha "Ferramentas e recursos"): novo [ADR-012](#adr-012--só-software-livre-ou-gratuito-e-seaweedfs-no-lugar-do-minio), com a regra de admissão de software, a auditoria da stack e o **SeaweedFS** no lugar do MinIO como armazenamento compatível com S3 e com bloqueio de objeto, ensaiado na prática; seções 2.3 e 2.4, nota na seção 7 (ambiente local), notas no ADR-004 (plano B) e no ADR-006 (suspenso e, na emenda abaixo, de volta como exceção explícita) e pendências novas na seção 16, que dependem do dono do produto. (2) **A Tenaz não guarda dado de saúde** (sumário 0.7, seção 3.6, "Dados de saúde"): saem as citações das seções 6.8 (o que nunca trafega na conexão), 9.3 (corrida de atribuição, que passa a usar `os.encerrar`, como o teste), 12.3 (o que fica fora do índice de IA) e 14.3 (o que nunca vai para log). (3) A tabela da seção 9.3 volta a descrever exatamente a suíte de `backend/testes/vazamento/`: linhas "Consultas sqlc", "Mapa de privilégios" (privilégios padrão) e "Funções de credencial" reescritas, e linhas novas "NaEmpresa confere o escopo no banco", "Extensões" e "Sem dado de saúde".

**Emenda na versão 0.7 (13/11/2026, mesma etapa):** respostas do dono do produto à auditoria do ADR-012 (sumário 0.7, seção 3.6, linha "Ferramentas e recursos"). A regra vale só para o que o sistema usa diretamente, não para o ambiente de desenvolvimento; a produção roda em nuvem gerenciada (AWS, São Paulo) como exceção explícita, e o ADR-006 volta a valer; o selo da trilha fica, em produção, no S3 com Object Lock em modo de conformidade, e o SeaweedFS fica só no desenvolvimento, numa máquina, endurecido; os serviços pagos por natureza ficam pendentes até a etapa em que entram. Emenda do ADR-012; seções 2.3, 2.4, 13.1 e 16; ADR-006 sem a suspensão. Na mesma passada, a convenção dos testes de integração (seção 8.2: PostgreSQL descartável por pacote, subida serializada e repetição da detecção do Docker) e a verificação estrita do access token do app (seção 10.6). Depois, o SeaweedFS aplicado no compose pelo `tenaz-plataforma-sre` (decisão 3 do ADR-012, auditoria, seções 7 e 16), com a correção de "desenvolvimento e CI" para "desenvolvimento": o CI ainda não o sobe, e o job fica como pendência na seção 16. Por fim, do `tenaz-mobile-offline`: Expo SDK 58 ratificado (seção 2.3 e auditoria do ADR-012), sessão do app com uma só credencial viva por aparelho, banco local por conta preservado ao sair e cofre com itens contados (seções 10.6 e 10.8), sem identificador de atualização enquanto não houver servidor de atualização (seção 9.5), e a pendência da janela de tolerância do refresh rotativo (seção 16).

**Mudanças da versão 0.8 deste documento (13/11/2026, etapa [014](etapas/014_2026-11-13_multipais-e-idiomas.md), sumário 0.8):** a Tenaz passa a ser multipaís (Brasil, Argentina e Paraguai) e bilíngue (`pt-BR` e `es-419`) desde o MVP. Novo [ADR-013](#adr-013--país-e-idioma-perfil-de-país-no-núcleo-textos-por-catálogo-e-conteúdo-legal-por-país-e-jurisdição), **proposta**: país da empresa e perfil de país no núcleo; identificação genérica por tipo e valor no lugar de `cnpj` e `cpf`, com migração nova e sem editar `00001` a `00007`; moeda ao lado de todo valor em dinheiro; fuso IANA validado pelo país; idioma por usuário e padrão da empresa; negociação por `Accept-Language` com `Content-Language`; textos do servidor e da interface em catálogos verificados no CI; conteúdo legal por país e por jurisdição, sem prazo legal no código; o núcleo deixa de citar a NR-12 pelo nome; bloqueio de energia como regra comum aos três países; hospedagem única em São Paulo dependente da P14. Incorpora a pesquisa do `tenaz-conformidade-mercosul` ([Argentina](pesquisa/mercosul_argentina.md), [Paraguai](pesquisa/mercosul_paraguai.md)). Notas nas seções 2.2, 7, 8.1, 8.4 e 13.1; a convenção de dinheiro da seção 8.1 (inteiro em centavos), que divergia dos documentos 2 e 3, fica substituída pela do ADR-013. Tabela de aplicação por dono e perguntas Q1 a Q3 ao dono do produto na seção 16. Na mesma versão, a seção 9.3 passa a descrever a parte 3 da suíte de isolamento, já pronta (`backend/testes/vazamento/api_*.go` e `selo_auditoria_test.go`): a linha "Leitura cruzada na API" vira "Leitura e escrita cruzadas na API", guiada pelo contrato; "Listagens e busca" ganha a verificação na API; a linha do R4 passa a cobrir toda rota do contrato, com web e app; e entra a linha "Trilha configurável e selo". Também na mesma versão, conferido contra o código de E01 do `tenaz-backend-go`: o subcomando `tenaz selar-auditoria` e a configuração do S3 (seção 4.1), o primeiro job real (seção 4.3), o pacote `internal/plataforma/armazenamento` (seção 2.4), as 17 variáveis de `NaEmpresa` e `db.SemTransacao` com o cache da captura (seção 6.3), a divergência da listagem de empresas do selo (seção 6.4), as dependências do cliente S3 na auditoria do ADR-012 e quatro pendências novas na seção 16 (agendamento do selo, credencial do S3 em produção, `go.mod` e listagem de empresas). Nenhum texto das versões anteriores foi reescrito, fora as linhas de versão e de alinhamento e a frase de abertura da seção 16.

**Aceite na versão 0.8 (13/11/2026, mesma etapa 014).** O dono do produto **aceitou o ADR-013** e respondeu às perguntas: um país por empresa no MVP, sem país na unidade (`unidade.pais` fora do MVP, gatilho (h)); uma moeda por empresa; idioma de último recurso pela regra proposta, terminando em `pt-BR` (decidido por padrão, salvo objeção); e o prefixo do número da OS pelo país da empresa (`OS-` no Brasil, `OT-` na Argentina e no Paraguai), gravado na criação e nunca alterado pelo idioma (Parte A, item 9). A migração da identificação passa a ser a `00009`, porque a `00008` é a trava da trilha (`00008_auditoria_insercao_por_coluna.sql`). As notas "(proposta, ADR-013)" das seções 2.2, 7, 8.1, 8.4 e 13.1 passam a "(ADR-013, aceito)", e a seção 16 marca as perguntas como respondidas. Na seção 9.3, do `tenaz-isolamento`: a linha "Trilha configurável e selo" passa a cobrir o evento sem data retroativa (gravação só pelas colunas do `GRANT`), e a linha "Mapa de privilégios" passa a ver também o privilégio concedido por coluna. Depois, conferido contra `app/`: textos e idioma do app sem `react-i18next`, idioma inicial do aparelho e escolha no cofre (`tenaz.idioma`), cores de `@tenaz/tokens-ui/tokens.json` e `resolveJsonModule` (seções 2.3 e 10.6); `expo-localization` e `expo-clipboard` (MIT) na auditoria do ADR-012. Por fim, o defeito do plural de reserva do app (zero em `pt-BR`), apontado na seção 10.6, foi corrigido pelo `tenaz-mobile-offline`, e a nota passou a "resolvido; a reserva é testada contra `Intl.PluralRules`". E, do `tenaz-plataforma-sre`, as tarefas `api` e `selar-auditoria` de `ferramentas/tarefas.sh` com o S3 de desenvolvimento pela identidade `tenaz_auditoria`, sem `Admin` (seção 7, decisão 5 do ADR-012 e seção 16). Na seção 9.3, a suíte adaptada à migração `00009` (`fixtures_test.go`, `pais_test.go`): empresas de três países (A no Brasil, B na Argentina, C no Paraguai), linha nova "País, identificação e idioma" e, na tabela do R4, as colunas que o externo altera na própria linha de `usuario`.

**Mudanças da versão 0.9 deste documento (13/11/2026, etapa E01, sumário 0.8).** Duas coisas. (1) **Registro do que o `tenaz-backend-go` implantou do ADR-013, conferido no código:** pacotes `plataforma/textos` (catálogos `pt-BR` e `es-419` embutidos, interpolação, plural CLDR sobre `golang.org/x/text/feature/plural`, que passou a dependência direta, e negociação do idioma) e `plataforma/pais` (perfil de país, normalização, validação e máscara da identificação, CNPJ alfanumérico, casos em `backend/testes/casos/identificacao.json`); erros por chave em `plataforma/erros`; `Content-Language` e `Vary` nas respostas; tarefa `textos` em `conferir`; redação de `identificacao` nos logs. Seções 2.1, 7, 8.2, 8.4 e 14.3, nota na Parte B, item 8, do ADR-013, linha nova na auditoria do ADR-012 e estado na seção 16. A migração dos textos livres que restavam no Go para o catálogo estava **em andamento na E01** e foi concluída na mesma versão (ver o fim deste parágrafo). (2) **Complemento do ADR-013, Parte F: representação do dinheiro.** Duas espécies (valor, `numeric(18,2)` na escala da moeda; preço unitário, `numeric(18,6)`), escala pela ISO 4217 (`BRL` e `ARS` 2, `PYG` 0) com `CHECK` por linha, texto decimal canônico na API, sem arredondar o que a pessoa digita, arredondamento meio para longe do zero (o mesmo do `round` do PostgreSQL), um só, no fim do cálculo; no Go, pacote `plataforma/dinheiro` com `int64` na unidade menor da moeda e `math/big` no cálculo intermediário; no TypeScript, `bigint` nativo em `pacotes/dominio`; nenhuma dependência nova. `usuario.custo_hora` passa a preço unitário, `numeric(18,6)`, e saem `api.Centavos` e `api.Dinheiro`. Substitui, no item 5 da Parte A, a proposta de biblioteca decimal, a largura delegada e a casa decimal "do perfil de país". Notas nas seções 2.1, 8.1 e 8.4 e aplicação por dono na seção 16. Fora a linha de versão, as notas marcadas e a frase de abertura da seção 16, nenhum texto anterior foi reescrito. Ainda na 0.9, depois da versão 1 de `backend/testes/casos/dinheiro.json` (`tenaz-qa-automacao`, 239 casos): a gramática da leitura passa a ser exatamente a do padrão do contrato (campo sem negativo, padrão sem `-?`), e o novo item 10 da Parte F fecha o rateio (piso sobre o módulo, sobra pelo maior resto, empate no menor índice, sinal devolvido), a quantidade como fração exata, o `-0` recusado em campo sem negativo, `formato_invalido` como único erro de leitura, `fora_da_faixa` só para resultado de operação, `moeda_diferente`, `moeda_desconhecida` e `rateio_invalido` como erros internos (500, sem código novo no documento 3) e o preço unitário sem soma por restrição de tipo. Por fim, também na 0.9: o texto livre dos erros no Go foi todo passado ao catálogo (118 pontos), e a trava do texto livre (`.Com(` e `.ComCampo(` proibidos no código de produção, `erros.Novo` conferido contra o catálogo, `Erro.Tecnico` para o log) roda na tarefa `textos`, com o limite declarado de não cobrir texto montado fora de `erros` (seção 8.2; a dívida deixa de estar em andamento). E, a partir de divergência apontada pelo `tenaz-produto` (a web formatava com `Intl` pela localidade do país, e o documento 7, seção 3.11, o descrevia): **Parte G do ADR-013, formatação para pessoas**. Mantém a Parte A, item 7, e descarta aceitar o `Intl` na web. A formatação vai para `pacotes/dominio/src/formatos.ts`, pelo país da empresa, sobre texto canônico e com fuso explícito, com `Intl.DateTimeFormat` só como motor de fuso, trava no ESLint e casos em `backend/testes/casos/formatos.json`; o formatador do Go espera o primeiro e-mail ou PDF. Notas nas seções 2.2, 8.4 e 16. Por último, na seção 9 (seções 9.1 e 9.4) e na árvore da seção 7: `backend/testes/integracao/` (fluxos HTTP) e o novo `backend/testes/migracoes/` (conversão, reversão e recusa das migrações de dados, do `tenaz-banco-dados`) não rodavam no CI e ganharam o passo "Testes de integração e de migração" em `.github/workflows/ci.yml`. Ainda na 0.9: a seção 9.3 registra a suíte com a empresa do Paraguai no cenário da API, as escritas de rota única executadas, a forma fixada do oráculo de unicidade, a prova do detector dos seis pares e a da FK para `app.empresa` por outra coluna. Também na 0.9: `dinheiro.json` na versão 2 (266 casos), gerado por `gerar_dinheiro.py` e conferido pela tarefa `conferir-casos`, no `conferir` e no CI (Parte F, item 8; seção 7).

Os nomes de módulos (ADM, ATV, OS...) e de entidades (`empresa`, `ativo`, `ordem_servico`...) seguem o [Sumário Executivo](0_SUMARIO_EXECUTIVO.md). As metas de disponibilidade, recuperação e segurança estão em [9_SEGURANCA_CONFIABILIDADE.md](9_SEGURANCA_CONFIABILIDADE.md).

**Convenção deste documento:** o que aparece marcado como **(proposta)** é decisão de projeto ainda não validada em código ou com custos reais. O que não tem marca é premissa já aceita pelo dono do produto.

---

## Índice

1. [Princípios](#1-princípios)
2. [Stack](#2-stack)
3. [Visão de contexto](#3-visão-de-contexto)
4. [Visão de contêineres](#4-visão-de-contêineres)
5. [Monólito modular](#5-monólito-modular)
6. [Isolamento entre empresas no código](#6-isolamento-entre-empresas-no-código)
7. [Estrutura do repositório](#7-estrutura-do-repositório)
8. [Convenções de código](#8-convenções-de-código)
9. [Testes e CI/CD](#9-testes-e-cicd)
10. [App offline-first](#10-app-offline-first)
11. [Telemática e IoT](#11-telemática-e-iot)
12. [Inteligência artificial](#12-inteligência-artificial)
13. [Infraestrutura](#13-infraestrutura)
14. [Observabilidade](#14-observabilidade)
15. [Decisões de arquitetura (ADRs)](#15-decisões-de-arquitetura-adrs)
16. [Alinhamento com o sumário e pendências](#16-alinhamento-com-o-sumário-e-pendências)

---

## 1. Princípios

| # | Princípio | Consequência prática |
|---|---|---|
| 1 | **Isolamento entre empresas acima de tudo** | Toda tabela de dados de cliente tem `empresa_id`, RLS forçado e filtro explícito no repositório. Um build com vazamento não vai para produção. |
| 2 | **Simples antes de distribuído** | Um binário Go, um banco PostgreSQL, uma fila dentro do banco. Um componente novo só entra com um gatilho mensurável (ver ADRs). |
| 3 | **O campo não tem sinal** | O app trata a conexão como exceção. Toda ação de campo funciona offline e converge sem perda. |
| 4 | **Um contrato, vários clientes** | A API pública, o web, o app e os integradores usam a mesma API REST `/v1`, descrita em OpenAPI 3.1. Não existe API "interna" paralela. |
| 5 | **Dados no Brasil** | Dados primários na região São Paulo. Qualquer saída do país é explícita e declarada no DPA. O local da cópia de recuperação de desastre é a pendência P7 (seção 13.5). |
| 6 | **Humano no comando da IA** | A IA lê, sugere e rascunha. Só grava depois de confirmação humana. |
| 7 | **Publicar o que se mede** | SLOs, RPO, RTO e histórico de incidentes são métricas reais, não texto de marketing. |
| 8 | **Acesso externo é exceção estreita** | Prestador, locatário e link de uso único veem só o que lhes foi atribuído, por uma camada de escopo adicional ao RLS por empresa (seção 6.6), com teste automatizado em todo build. |
| 9 | **Cada conta guarda a sua cópia** | Contas conectadas (cliente e subcontratada) nunca leem dados uma da outra. Um processo de sincronização com papel de banco próprio copia só os campos permitidos, gravando em cada conta numa transação separada (seção 6.8, ADR-009). |

---

## 2. Stack

### 2.1 Backend

| Item | Escolha | Observação |
|---|---|---|
| Linguagem | Go (versão estável vigente, 1.26 ou superior) | Política: só usar versões com suporte oficial. Atualização em até 60 dias após cada lançamento. |
| HTTP e contrato | `net/http` padrão + **Huma v2** (proposta) | Huma gera o OpenAPI 3.1 a partir dos tipos Go, valida a entrada contra o schema e responde erros no formato RFC 9457. Alternativa avaliada: `ogen` (spec-first). |
| Acesso a dados | **sqlc** + **pgx v5** | SQL explícito, tipado na compilação. Sem ORM. |
| Migrações | **goose** (proposta), arquivos `.sql` versionados | Padrão expandir e contrair (expand/contract) para não derrubar versões antigas durante o deploy. |
| Jobs e eventos | **River** (fila no PostgreSQL) | Ver [ADR-005](#adr-005--fila-de-jobs-no-postgresql). |
| Logs | `log/slog` em JSON | Campos obrigatórios na seção 14. |
| Telemetria | OpenTelemetry Go SDK | Traces, métricas e correlação com logs. |
| Testes | `testing` + testify + **testcontainers-go** | Postgres real nos testes de integração. |
| Textos do servidor (versão 0.9) | `plataforma/textos` sobre `golang.org/x/text/feature/plural` (BSD-3; dependência direta desde a E01, sem biblioteca de i18n) | [ADR-013](#adr-013--país-e-idioma-perfil-de-país-no-núcleo-textos-por-catálogo-e-conteúdo-legal-por-país-e-jurisdição), Parte B. |
| Dinheiro (versão 0.9) | `plataforma/dinheiro`, só biblioteca padrão (`int64` e `math/big`) | ADR-013, Parte F. Nenhuma biblioteca decimal. |

### 2.2 Web

| Item | Escolha |
|---|---|
| Base | React 19 + TypeScript (modo `strict`), Vite, SPA |
| Dados do servidor | TanStack Query + cliente gerado do OpenAPI (`openapi-typescript` + `openapi-fetch`) |
| Rotas | React Router (proposta) |
| Formulários | React Hook Form + Zod |
| Interface | Tailwind CSS + componentes Radix, com os tokens da [marca](../marca/README.md). Detalhes em `7_PADROES_DESIGN.md`. |
| Internacionalização | i18next (proposta), com `pt-BR` como idioma de origem. **Versão 0.8 (ADR-013, aceito):** `pt-BR` e `es-419`, idioma do usuário enviado em `Accept-Language`, catálogo `es-419` tipado pelo de `pt-BR` e formatação de número, data e moeda por funções próprias de `pacotes/dominio`, pelo país, nunca pelo `Intl` com `es-419`. **Versão 0.9 (ADR-013, Parte G):** mantido; nenhum `Intl` para número nem moeda, nem com a localidade do país; `Intl.DateTimeFormat` só como motor de fuso; casos em `backend/testes/casos/formatos.json` |
| Testes | Vitest + Testing Library; Playwright para E2E |

### 2.3 App de campo (MOB)

| Item | Escolha |
|---|---|
| Base | React Native com **Expo** (SDK estável mais recente) + TypeScript. **Versão 0.7:** SDK 58 (`expo ~58.0.25`), ratificado: o repositório misturava `expo` 54 com `expo-router`, `expo-sqlite` e `expo-secure-store` 58, e o bundle não resolvia; com o 58, `expo export --platform android` gera o bundle. Do 54 para o 58, `newArchEnabled` e `splash` saíram de `app.config.ts`, e a tela de abertura passou ao plugin `expo-splash-screen` |
| Banco local | **SQLite** via `expo-sqlite`, criptografado com SQLCipher (proposta) |
| Acesso ao SQLite | Drizzle ORM com consultas reativas (proposta) |
| Credenciais | `expo-secure-store` (Keychain e Keystore) |
| Recursos nativos | Câmera e QR code, localização sob demanda, notificações push, compressão de foto |
| Distribuição | EAS Build (lojas) e EAS Update (atualização do JavaScript sem loja, por canal). **Versão 0.7:** o EAS é serviço hospedado com cota gratuita e excedente pago, e as lojas cobram a conta de desenvolvedor; a forma de distribuição é decidida em E12 (ADR-012, pendência H5 da seção 16) |
| Testes | Jest + React Native Testing Library; Maestro para E2E em emulador (proposta) |
| Textos e idioma (versão 0.8, ADR-013) | `@tenaz/i18n` **sem** `react-i18next` (`app/src/i18n`): busca por chave tipada (chave inexistente é erro de tipo), interpolação `{{nome}}` e plural por `Intl.PluralRules`, com volta a `pt-BR` e depois à própria chave se o texto faltar. Idioma inicial do aparelho (`expo-localization`): `pt*` vira `pt-BR`, `es*` vira `es-419`, o resto `pt-BR`. A escolha da pessoa fica no cofre local, chave `tenaz.idioma`, e vale antes do login, até a API ter `usuario.idioma` |
| Cores (versão 0.8) | Lidas de `@tenaz/tokens-ui/tokens.json` (`app/src/componentes/tema.ts`), sem transcrição à mão; o `tsconfig` do app liga `resolveJsonModule` |

### 2.4 Dados e plataforma

| Item | Escolha |
|---|---|
| Banco | PostgreSQL gerenciado (versão 17 ou superior suportada pelo provedor), com RLS forçado |
| Busca vetorial | Extensão `pgvector` no mesmo banco |
| Arquivos | Object storage compatível com S3, com bloqueio de objeto. Produção: S3 gerenciado, com Object Lock em modo de conformidade no bucket do selo (ADR-006). Desenvolvimento: **SeaweedFS** (Apache-2.0), por versão e digest, no lugar do MinIO (o CI ainda não o sobe; pendência na seção 16). O backend fala só a API S3, pelo `aws-sdk-go-v2` (Apache-2.0), sem API própria do servidor ([ADR-012](#adr-012--só-software-livre-ou-gratuito-e-seaweedfs-no-lugar-do-minio)). **Versão 0.8 (aplicado em E01):** pacote `internal/plataforma/armazenamento`, interface `Armazenamento` com `GravarImutavel` (Object Lock `COMPLIANCE` com data de retenção), `Existe` (`HEAD`) e `Ler`, implementada sobre `aws-sdk-go-v2`, `service/s3`, `credentials` e `smithy-go`. A repetição é só a da plataforma (`resiliencia`: 10 s por tentativa, 3 tentativas com recuo, disjuntor aberto por 30 s depois de 5 falhas seguidas), com a repetição do SDK desligada; erro de credencial, acesso negado, bucket inexistente e pedido malformado são permanentes, sem nova tentativa. Primeiro uso: o selo da trilha (seção 4.1) |
| Cache | Valkey gerenciado (proposta), nunca como fonte de verdade |
| Contêineres | Docker, imagem mínima (ver 9_SEGURANCA_CONFIABILIDADE.md) |
| Nuvem | AWS `sa-east-1` ([ADR-006](#adr-006--nuvem-aws-na-região-são-paulo)): provedor e região decididos pelo dono do produto em 13/11/2026, como exceção explícita à regra de admissão (ADR-012); serviços e configuração continuam proposta |
| IaC | OpenTofu (proposta) |
| Regra de admissão | O que o sistema usa diretamente (bibliotecas, imagens, servidores e serviços que vão para a produção) é de código aberto (licença aprovada pela OSI; fontes sob OFL) ou gratuito sem condição; licença "source-available" e recurso pago de produto "open core" não entram. O ambiente de desenvolvimento fica fora da regra, e a hospedagem de produção é a exceção explícita ([ADR-012](#adr-012--só-software-livre-ou-gratuito-e-seaweedfs-no-lugar-do-minio)) |

### 2.5 Contrato único

A API é a única fronteira entre os clientes e o backend.

1. Os tipos de requisição e resposta são structs Go anotadas, registradas no Huma (v2). As rotas de todos os pacotes são montadas num ponto só, `internal/servidor` (`servidor.Montar`), usado pelo processo `tenaz api`, pela geração do contrato e pelos testes de rota: o contrato gerado é exatamente o que o servidor atende. O subcomando `tenaz openapi` escreve `backend/api/openapi.yaml` (OpenAPI 3.1); a tarefa `contrato` de `ferramentas/tarefas.sh` o regrava e regera os tipos TypeScript.
2. O arquivo gerado é versionado no repositório. O CI (tarefa `conferir-contrato`) falha se o arquivo gerado diferir do versionado (alguém mudou a API sem regenerar).
3. O CI compara o contrato com o da branch principal usando `oasdiff`, e mudança incompatível em `/v1` bloqueia o merge. **(proposta: ainda não está na esteira; pendência na seção 16)**
4. O pacote `pacotes/api-tipos` gera os tipos TypeScript a partir do arquivo, e o CI confere que os tipos versionados estão em dia. A web importa dali (`@tenaz/api-tipos`); o app passa a importar quando tiver chamadas à API. Um campo renomeado no Go quebra a compilação do TypeScript no mesmo PR.
5. A documentação pública da API é gerada do mesmo arquivo. **(proposta: ainda não implantada)**

---

## 3. Visão de contexto

Quem usa a Tenaz e com quais sistemas ela conversa.

```mermaid
flowchart LR
    tecnico["Técnico e motorista<br/>(app de campo)"]
    gestor["Gestor, planejador,<br/>almoxarife, administrador<br/>(web)"]
    solicitante["Solicitante<br/>(web, app ou QR code)"]
    integrador["Integrador do cliente<br/>(API e webhooks)"]
    prestador["Prestador: subcontratada,<br/>assistência técnica, oficina<br/>(app e web restritos)"]
    ocasional["Prestador ocasional<br/>(link de uso único, sem login)"]
    locatario["Locatário<br/>(portal web restrito)"]
    subcontratada["Subcontratada com conta própria<br/>Tenaz Serviços (SRV)<br/>(gestor no web, técnicos no app)"]
    clientefinal["Cliente da subcontratada<br/>sem conta Tenaz<br/>(recebe o relatório PDF da OS)"]

    tenaz[["Tenaz<br/>Gestão de Manutenção e Frota"]]

    idp["Provedor de identidade do cliente<br/>(OIDC, SAML, SCIM)"]
    telem["Telemáticas e montadoras<br/>(Cobli, Geotab, rFMS Volvo/Scania,<br/>Sascar, Omnilink)"]
    cartao["Cartão combustível<br/>(Ticket Log e outros)"]
    erp["ERPs<br/>(Sankhya, Protheus, Omie)"]
    llm["Provedor de LLM<br/>(sem treinamento com dados)"]
    msg["E-mail e push<br/>(SES, APNs, FCM)"]
    gw["Gateway IoT parceiro<br/>(Fase 3: MQTT Sparkplug B)"]
    status["Página de status<br/>(hospedagem externa)"]

    tecnico --> tenaz
    gestor --> tenaz
    solicitante --> tenaz
    integrador --> tenaz
    prestador --> tenaz
    ocasional --> tenaz
    locatario --> tenaz
    subcontratada --> tenaz
    tenaz -.-> clientefinal
    tenaz --> idp
    tenaz --> telem
    tenaz --> cartao
    tenaz <--> erp
    tenaz --> llm
    tenaz --> msg
    gw --> tenaz
    tenaz -.-> status
```

| Sistema externo | Direção | Fase | Observação |
|---|---|---|---|
| Provedor de identidade | Tenaz consulta | MVP (OIDC), Fase 2 (SAML, SCIM) | Google e Microsoft Entra primeiro |
| Telemáticas e montadoras | Tenaz consulta ou recebe webhook | Fase 2 | Modelo canônico de evento (seção 11) |
| Cartão combustível | Tenaz importa | Fase 2 | Arquivo ou API, conforme o fornecedor |
| ERPs | Bidirecional | Fase 2 | Conector por ERP; o ERP é dono dos itens fiscais, das contas a pagar e da **emissão** de NF-e (remessa e retorno a subcontratada, envio para conserto, fatura de locação e de avarias). O cadastro de `fornecedor` com papéis é do módulo TER e é sincronizado com o ERP; a Tenaz só registra número e chave da NF-e emitida |
| Provedor de LLM | Tenaz chama | Fase 2 | Atrás de uma interface interna; troca de provedor sem mudar módulos |
| E-mail e push | Tenaz envia | MVP | Notificações, alertas de vencimento e envio do link de uso único ao prestador ocasional |
| Prestadores e locatários | Usam a Tenaz | MVP (prestador e link), Fase 2 (locatário) | Não são sistemas: são pessoas de empresas externas com acesso restrito ao que lhes foi atribuído (seção 6.6) |
| Subcontratada com conta própria (SRV) | Usa a Tenaz como empresa | MVP | Não é sistema externo: é outra `empresa` da Tenaz, isolada pelo mesmo RLS. Quando o cliente também usa a Tenaz, as duas contas se conectam e a OS trafega pelo sincronizador entre contas (seção 6.8, ADR-009). Seus técnicos usam o mesmo Tenaz Mobile (seção 10.8) |
| Cliente da subcontratada sem conta Tenaz | Recebe da Tenaz | MVP | Recebe o relatório PDF da OS (`GET /v1/ordens-servico/{id}/relatorio-pdf`), enviado pela subcontratada por e-mail ou baixado e repassado por ela. Não acessa a Tenaz |
| Gateway IoT parceiro | Envia para a Tenaz | Fase 3 | A Tenaz nunca se conecta a CLP ou SCADA do cliente |
| Página de status | Tenaz publica | MVP | Fora da infraestrutura principal, para continuar no ar durante incidentes |

---

## 4. Visão de contêineres

```mermaid
flowchart TB
    subgraph clientes["Clientes"]
        web["Web SPA<br/>React + TypeScript<br/>(CDN + object storage)"]
        portal["Portais do prestador e do locatário<br/>(mesma SPA, layout restrito)"]
        link["Página do link de uso único<br/>(bundle separado, sem login)"]
        app["App de campo Tenaz Mobile<br/>Expo + SQLite local por conta<br/>outbox de comandos"]
        ext["Integrações do cliente<br/>API pública /v1"]
    end

    subgraph borda["Borda"]
        cdn["CDN + WAF"]
        lb["Balanceador de carga<br/>(TLS)"]
    end

    subgraph computacao["Contêineres Go (um binário, quatro modos)"]
        api["API<br/>tenaz api<br/>REST /v1, sync, webhooks de entrada,<br/>rotas públicas do link"]
        worker["Workers<br/>tenaz worker<br/>preventiva, alertas, webhooks, push,<br/>importações, IA, exportações, PDF"]
        sinc["Sincronizador entre contas<br/>tenaz sincronizador<br/>papel de banco tenaz_sinc<br/>(SRV, ADR-009)"]
        ingest["Ingestão de telemetria<br/>tenaz ingestor<br/>(Fase 2)"]
    end

    subgraph dados["Dados"]
        pg[("PostgreSQL primário<br/>RLS forçado, pgvector,<br/>fila River, outbox")]
        replica[("Réplica de leitura<br/>painéis e exportações")]
        cache[("Cache Valkey<br/>rate limit, sessões revogadas,<br/>entitlements")]
        s3[("Object storage<br/>fotos, anexos, manuais,<br/>exportações")]
        ts[("Séries temporais<br/>(Fase 3)")]
    end

    llm["Provedor de LLM"]
    push["Expo Push<br/>APNs e FCM"]
    broker["Broker MQTT gerenciado<br/>(Fase 3)"]
    telem["APIs de telemática"]

    web --> cdn --> lb
    portal --> cdn
    link --> cdn
    app --> lb
    ext --> lb
    lb --> api
    api --> pg
    api --> replica
    api --> cache
    api --> s3
    worker --> pg
    worker --> s3
    worker --> llm
    worker --> push
    push -. notificação de tarefa nova .-> app
    sinc -- "transação na conta de origem" --> pg
    sinc -- "transação na conta de destino" --> pg
    ingest --> telem
    broker --> ingest
    ingest --> pg
    ingest --> ts
    pg -. replicação .-> replica
```

| Contêiner | Responsabilidade | Escala | Fase |
|---|---|---|---|
| Web SPA | Interface do gestor, planejador, almoxarife, comprador e administrador | Arquivos estáticos na CDN | MVP |
| Portais do prestador e do locatário | Mesma SPA com layout e menu restritos, carregados conforme o escopo da sessão (seção 6.6). O prestador vê só as OS do seu fornecedor e o almoxarifado em poder dele; o locatário vê só os ativos locados a ele, suas solicitações e OS | Arquivos estáticos na CDN | MVP (prestador), Fase 2 (locatário) |
| Página do link de uso único | Bundle web separado e mínimo, para celular, sem acesso à sessão da SPA principal; mostra uma única OS e registra a execução (seção 6.7) | Arquivos estáticos na CDN | MVP |
| App de campo (Tenaz Mobile) | Técnico da empresa, técnico da subcontratada com conta própria (SRV), motorista, solicitante e usuário de prestador; offline; troca de conta para quem atua em mais de uma (seção 10.8) | No dispositivo | MVP |
| API (`tenaz api`) | REST `/v1`, autenticação, autorização, escopo externo, sincronização do app, recepção de webhooks de terceiros, rotas públicas do link de uso único, gestão de conexões entre contas (`/v1/conexoes-empresa`) | Horizontal, sem estado; mínimo de 2 tarefas em zonas diferentes | MVP |
| Workers (`tenaz worker`) | Geração de OS preventiva, alertas de vencimento, envio de webhooks (MVP, plano Profissional), notificações push de tarefa nova, envio e expiração de links de uso único, geração do relatório PDF da OS, importação de CSV e NF-e, indexação de documentos para IA, exportação completa da empresa; na Fase 2, prazos de ordem de reparo, verificação de SLA de fornecedor e vencimento de contratos de locação | Horizontal; mínimo de 2 tarefas | MVP |
| Sincronizador (`tenaz sincronizador`) | Único componente que grava dados vindos de outra conta: consome a fila `conexao`, lê o evento na conta de origem e aplica na conta de destino, em transações separadas, com o papel de banco `tenaz_sinc` (seção 6.8, ADR-009). Credencial própria, inacessível às tarefas da API e dos workers | Horizontal; mínimo de 2 tarefas; ordem garantida por vínculo | MVP |
| Ingestão (`tenaz ingestor`) | Conectores de telemática, normalização para o modelo canônico, consumo MQTT na Fase 3 | Separado da API para não disputar recursos | Fase 2 |
| PostgreSQL primário | Fonte de verdade, fila de jobs, outbox, vetores | Vertical; réplica para leitura | MVP |
| Réplica de leitura | Painéis de IND, relatórios, exportações | Uma réplica; mais se o lag ou a carga exigirem | MVP |
| Cache | Rate limit por chave de API, lista de sessões revogadas, cache de entitlements | Nó pequeno com réplica | MVP (proposta) |
| Object storage | Arquivos com chave prefixada por `empresa_id` | Gerenciado | MVP |
| Séries temporais | Dados de sensor de alta frequência | Ver seção 11 | Fase 3 |

### 4.1 Um binário, quatro modos

O backend compila um único binário estático `tenaz` com subcomandos: `api`, `worker`, `sincronizador`, `ingestor` e `migrate`. Os quatro modos de serviço compartilham código e configuração, mas rodam como serviços separados para escalar e falhar de forma independente. **(proposta)**

Motivos: uma imagem só para construir, assinar e escanear; versões sempre coerentes entre API, workers e sincronizador; custo baixo de memória por tarefa.

O `sincronizador` é um modo separado, e não uma fila a mais dentro do `worker`, por um motivo de segurança: ele é o único que conecta ao banco com o papel `tenaz_sinc`. A credencial desse papel fica num segredo que só a tarefa do sincronizador lê (política IAM própria). Uma falha na API ou nos workers não dá acesso ao papel que atravessa contas, e o papel `tenaz_app` nunca grava em duas empresas no mesmo processo de requisição.

**Credencial de banco por modo (versão 0.6).** Cada modo recebe só as credenciais de que precisa, cada uma num segredo próprio lido pela política IAM daquela tarefa:

| Modo | Papéis de banco | Observação |
|---|---|---|
| `api` | `tenaz_app` (pool da aplicação, réplica com `tenaz_leitura`) e `tenaz_autenticacao` (pool pequeno, só do pacote de autenticação) | O segundo pool existe só neste modo ([ADR-011](#adr-011--papel-de-banco-próprio-para-as-credenciais-de-autenticação), seção 6.9) |
| `worker` | `tenaz_app` e `tenaz_job` | Não recebe a credencial de autenticação: a purga das credenciais vencidas é uma função `SECURITY DEFINER` estreita (seção 6.9) |
| `sincronizador` | `tenaz_sinc` | Seção 6.8 |
| `ingestor` | `tenaz_app` | Fase 2 |
| `migrate` | `tenaz_migrador` | Só no pipeline de deploy |
| `selar-auditoria` (versão 0.8) | `tenaz_job` | Subcomando de execução única, não serviço (abaixo) |

**Subcomando `tenaz selar-auditoria` (versão 0.8, etapa E01; código em `cmd/tenaz/main.go` e `internal/plataforma/auditoria/selo.go`).** É o job diário do selo da trilha (RF-ADM-011; documento 2, seção 9.4). Roda como `tenaz_job`, com a credencial `TENAZ_BANCO_URL_JOB` (obrigatória fora do desenvolvimento, diferente das URLs da aplicação e da autenticação, com TLS em produção; no desenvolvimento, vazia, usa `TENAZ_BANCO_URL` com `SET ROLE tenaz_job`). Percorre as empresas e, para cada uma, numa transação `NaEmpresa` com escopo interno: lê `retencao_auditoria_anos`, sela os dias pendentes com `app.selar_dias_pendentes()` (o banco calcula tudo) e lê os selos a conferir. Fora da transação, grava no armazenamento imutável cada selo novo e cada selo dos últimos 7 dias que esteja sem cópia (conferido por `HEAD`), com Object Lock em modo `COMPLIANCE` e retenção até agora + `retencao_auditoria_anos`. Falha numa empresa não interrompe as outras: os erros voltam juntos no fim, e a rodada é idempotente. **Não há agendador ainda** (o River do ADR-005 continua proposta): a infraestrutura dispara o subcomando uma vez por dia, depois de 01:00 UTC (pendência na seção 16). O armazenamento é **obrigatório** no subcomando (não parte sem as variáveis abaixo) e **opcional** na API: sem ele, a verificação de integridade responde `503 servico_indisponivel`, porque sem a cópia não há como afirmar que o período está íntegro.

Configuração do armazenamento (`internal/plataforma/config`): `TENAZ_S3_ENDPOINT` (vazio: o endpoint da AWS para a região), `TENAZ_S3_REGIAO`, `TENAZ_S3_CAMINHO_EXPLICITO` (`true` ou `false`; o SeaweedFS exige `true`), `TENAZ_S3_BUCKET_AUDITORIA`, `TENAZ_S3_CHAVE_ACESSO` e `TENAZ_S3_CHAVE_SECRETA`.

### 4.2 Uso da réplica de leitura

- Vão para a réplica: painéis de IND, relatórios, exportações e consultas da IA que só leem.
- Ficam no primário: qualquer leitura logo após uma escrita do mesmo usuário, a sincronização do app e toda validação de regra de negócio.
- O contexto de empresa (`app.empresa_id`) funciona igual na réplica, porque é definido dentro da transação de leitura.
- Se o atraso da réplica passar de 30 segundos, as leituras voltam para o primário e um alerta dispara. **(proposta)**

### 4.3 Fila de jobs

No início, a fila é uma tabela no próprio PostgreSQL, consumida com `FOR UPDATE SKIP LOCKED`, por meio da biblioteca **River**. O job é inserido **na mesma transação** da mudança de negócio. Se a transação falhar, o job não existe; se ela confirmar, o job existe. Isso elimina a classe de bugs "gravou mas não avisou".

Gatilhos para migrar parte da carga para uma fila dedicada (SQS ou NATS JetStream) estão no [ADR-005](#adr-005--fila-de-jobs-no-postgresql). A telemetria de alta frequência nunca passa pela fila de jobs evento a evento: ela é agrupada em lotes pelo `ingestor`.

**Jobs da configuração por segmento (versão 0.5).** Rodam no `tenaz worker`, na mesma fila River, sem componente novo ([ADR-010](#adr-010--configuração-por-segmento-com-catálogos-fechados), seção 5.6) **(proposta)**:

| Job | Quando | Fase |
|---|---|---|
| Avaliação de regras por evento | Assinante do outbox (seção 5.5) para os eventos que correspondem às condições do catálogo | MVP |
| Rotina diária de regras | Agendada uma vez por dia, um job por empresa ativa (padrão da seção 6.4), depois da rotina de alertas de vencimento | MVP |
| Aplicação de pacote de segmento | Sob demanda, um job único por empresa, pacote e versão | Fase 2 |
| Medição mensal do contrato com cliente | Abertura no primeiro dia do mês seguinte, um job por empresa com CTR | Fase 2 |
| Prazos de contrato com cliente e de etapa de projeto | Assinantes de `ordem_servico.*` e rotina diária por empresa | Fase 2 |

Os jobs de regras usam uma fila River própria, `regras` **(proposta)**, para que uma avaliação lenta não atrase webhooks, push e o sincronizador.

**Primeiro job real (versão 0.8):** o selo diário da trilha roda como subcomando `tenaz selar-auditoria`, disparado pela infraestrutura, e não pela fila, porque o River ainda não está no projeto (seção 4.1). Quando o `worker` existir, o selo passa a ser um job agendado da fila, com o mesmo serviço `auditoria.Selador`.

---

## 5. Monólito modular

### 5.1 Um pacote por módulo

Cada módulo do sumário vira um pacote Go em `backend/internal/modulos/`. Três nomes colidem com identificadores da linguagem (`os`, `int` e `sync`) e viram `ordemservico`, `integracoes` e `mob`, como define o sumário (seção 3.4). Os módulos TER, LOC e SRV usam o nome por extenso, `terceiros`, `locacao` e `servicos`.

| Código | Pacote Go | Motivo do nome |
|---|---|---|
| ADM | `adm` | |
| ATV | `atv` | |
| MED | `med` | |
| SOL | `sol` | |
| OS | `ordemservico` | `os` colide com o pacote `os` da biblioteca padrão |
| PRE | `pre` | |
| CHK | `chk` | |
| EST | `est` | Inclui almoxarifados `proprio`, `em_terceiro` e `quarentena`, peças serializadas, remessa e retorno a terceiro e ordem de reparo |
| TER | `terceiros` | Nome por extenso, como no sumário. Dono de `fornecedor`, `servico_fornecedor`, `contrato_prestacao`, `avaliacao_fornecedor`, `orcamento_fornecedor` e `link_execucao` |
| LOC | `locacao` | Nome por extenso, como no sumário. Dono de `locatario`, `contrato_locacao`, `item_contrato_locacao`, `responsabilidade_locacao` e `vistoria_locacao` |
| SRV | `servicos` | Nome por extenso, como no sumário. Dono de `cliente_servico`, `especialidade_tecnico`, `pessoa_especialidade`, `vinculo_os` e da lógica de `plataforma.conexao_empresa`. Contém o código de aplicação de mensagens que o sincronizador executa (seção 6.8). Reaproveita OS, apontamento, ativos e app; não duplica nenhuma dessas entidades |
| FRO | `fro` | |
| CON | `con` | |
| IND | `ind` | |
| INT | `integracoes` | `int` é tipo predeclarado do Go |
| IA | `ia` | |
| MOB | `mob` | Contém a sincronização do app (o nome `sync` colidiria com a biblioteca padrão) |
| PRJ | `projetos` | Nome por extenso, como no sumário (seção 3.4). Fase 2. Dono de `modelo_projeto`, `modelo_etapa_projeto`, `projeto` e `etapa_projeto`. Não executa trabalho: organiza OS `implantacao` pelos contratos internos de `ordemservico` |
| CTR | `contratos` | Nome por extenso, como no sumário (seção 3.4). Fase 2. Dono de `contrato_cliente`, `item_contrato_cliente`, `faixa_desconto_contrato` e `medicao_contrato`. Os contratos de prestação e de locação continuam em `terceiros` e `locacao` |
| (ADM) | `configuracao` | Sem código de módulo próprio: é o pacote das regras configuráveis e dos pacotes de segmento, que o sumário põe no escopo de ADM (seções 2 e 3.4). Separado de `adm` para não criar ciclo, porque depende de ATV, CON, SOL e OS, que dependem de ADM. Dono de `regra_automacao`, `disparo_regra`, `pacote_segmento` e `aplicacao_pacote`. **Não** é dono de `campo_extra_os`: a definição do campo extra é de OS (RF-OS-021), que valida os valores na própria transação (seção 5.6) |

Código transversal fica em `backend/internal/plataforma/`: banco e transação, contexto de empresa e de escopo externo (seção 6.6), autenticação, autorização, entitlements, outbox e eventos, auditoria, armazenamento de arquivos, erros, observabilidade e configuração. A plataforma não importa nenhum módulo. A autenticação (`plataforma/autenticacao`) é o único pacote com papel de banco e pool próprios dentro da API: só ele lê e grava as tabelas de credencial (seção 6.9, regra 12 da seção 5.3).

### 5.2 Camadas dentro de um módulo

Layout validado em código pelo ADM (etapa E01) e adotado como regra para os módulos seguintes (versão 0.6):

```
backend/internal/modulos/adm/
├── adm.go                 # API pública: interface Servico, DTOs de entrada e saída
├── servico.go             # Novo(Dependencias), o tipo não exportado que implementa Servico
│                          # e os erros próprios do módulo (ex.: ErrUltimoAdministrador)
├── unidades.go, pessoas.go, usuarios.go, conta.go, auditoria.go
│                          # casos de uso, um arquivo por assunto: orquestram domínio,
│                          # repositório, auditoria, autorização e outros módulos
├── eventos.go             # eventos que o módulo publica (quando houver; regra 6 da seção 5.3)
├── rotas/                 # camada HTTP: handlers Huma que mapeiam a requisição para o Servico;
│                          # importado só pela composição em internal/servidor
├── jobs/                  # workers River do módulo, registrados pela composição do modo
│                          # worker (proposta: nenhum módulo tem job ainda)
└── internal/
    ├── dominio/           # entidades, regras, máquina de estados. Sem IO
    └── repo/
        ├── adm.sql        # consultas sqlc do módulo
        ├── repo.go        # consultas na transação de NaEmpresa; violação do banco vira erro de domínio
        └── gerado/        # regra: código gerado pelo sqlc (não editar); no ADM, ainda em
                           # plataforma/db/consultas (desvio abaixo)
```

O diretório `internal/` dentro do módulo é imposto pelo compilador Go: nenhum pacote fora de `adm/` importa `adm/internal/...`. Para outros módulos, o único ponto de entrada é o pacote raiz `adm`.

**Por que os casos de uso ficam no pacote raiz e o HTTP fica em `rotas/`, fora de `internal/`.** Os DTOs da API pública ficam no pacote raiz, e quem implementa o `Servico` precisa deles. Com os casos de uso em `internal/app`, o construtor `Novo`, no pacote raiz, importaria `internal/app`, que importaria o pacote raiz pelos DTOs: ciclo de importação. Com o HTTP em `internal/http`, a composição em `internal/servidor` não poderia importá-lo (o compilador só deixa `adm/...` importar `adm/internal/...`), e o pacote raiz teria de registrá-lo, de novo em ciclo. Alternativas descartadas: pôr os DTOs num subpacote `tipos` (um pacote a mais por módulo, sem ganho de fronteira) e registrar as rotas a partir do pacote raiz (o domínio passaria a depender de HTTP). A implementação continua escondida: só os identificadores exportados do pacote raiz são API, e o tipo que implementa `Servico` não é exportado. Os workers River seguem o mesmo raciocínio em `jobs/`.

| Camada | Pode importar | Não pode importar |
|---|---|---|
| `internal/dominio` | Biblioteca padrão, tipos da plataforma (IDs, erros, autorização, tempo) | Banco, HTTP, outros módulos |
| Pacote raiz (API pública e casos de uso) | `internal/dominio`, `internal/repo`, plataforma, **API pública** de outros módulos | `internal/` de outros módulos, `rotas` e `jobs` de qualquer módulo |
| `internal/repo` | Código gerado do sqlc, `plataforma/db`, `plataforma/erros` | Pacote raiz, `rotas`, `jobs` |
| `rotas` | Pacote raiz do próprio módulo, `plataforma/api`, `plataforma/autorizacao`, `plataforma/erros` e as interfaces de serviços da plataforma (ex.: convites de `plataforma/autenticacao`) | `internal/` do módulo (o compilador já impede), outros módulos |
| `jobs` **(proposta)** | Pacote raiz do próprio módulo, `plataforma/fila` | `internal/` do módulo, outros módulos |

Só `internal/servidor` importa `<modulo>/rotas`, e só a composição do worker importará `<modulo>/jobs`. Hoje isso vale por construção (só existe `adm/rotas`, importado só por `internal/servidor`); a regra do `depguard` entra com o segundo módulo (pendência na seção 16).

**Desvio atual do ADM: código sqlc compartilhado.** O `sqlc.yaml` gera as consultas de `adm/internal/repo/adm.sql` em `internal/plataforma/db/consultas`, o mesmo pacote das consultas de auditoria e de autorização, e `db.Q(ctx)` devolve esse tipo único. Como o pacote fica na plataforma, qualquer módulo poderia importá-lo e chamar uma consulta de ADM, e a regra 2 da seção 5.3 (nenhum módulo lê ou grava tabela de outro) passaria a depender só de revisão. A regra continua sendo a do layout acima, já seguida por `plataforma/autenticacao` (`autenticacao/internal/repo/gerado`): o código gerado de cada módulo fica no `internal/repo/gerado` dele, criado sobre a transação de `db.Tx(ctx)`. Correção antes do primeiro módulo de E02 (pendência na seção 16).

### 5.3 Regras de dependência entre módulos

As setas indicam "chama a API pública de". Ciclos são proibidos. Quando dois módulos precisam reagir um ao outro, um deles chama e o outro escuta um evento.

```mermaid
flowchart BT
    ADM
    ATV --> ADM
    MED --> ATV
    SOL --> ATV
    CON --> ATV
    CON --> ADM
    EST --> ADM
    TER --> ADM
    TER --> ATV
    EST --> TER
    LOC --> ATV
    LOC --> MED
    LOC --> TER
    LOC --> CHK
    SRV --> ADM
    SRV --> ATV
    SRV --> TER
    SRV --> OS
    OS --> ATV
    OS --> MED
    OS --> EST
    OS --> CON
    OS --> TER
    OS --> LOC
    PRE --> OS
    PRE --> MED
    CHK --> ATV
    CHK --> SOL
    FRO --> ATV
    FRO --> MED
    FRO --> EST
    INT --> OS
    INT --> MED
    INT --> FRO
    INT --> TER
    INT --> LOC
    IA --> OS
    IA --> ATV
    MOB --> OS
    MOB --> CHK
    MOB --> MED
    MOB --> SOL
    MOB --> EST
    MOB --> LOC
    PRJ --> OS
    PRJ --> ATV
    PRJ --> TER
    PRJ --> CON
    PRJ --> CTR
    PRJ --> SRV
    CTR --> ATV
    CTR --> OS
    CTR --> SRV
    INT --> PRJ
    CFG["configuracao"] --> ADM
    CFG --> ATV
    CFG --> MED
    CFG --> SOL
    CFG --> CHK
    CFG --> PRE
    CFG --> CON
    CFG --> OS
    CFG --> PRJ
```

Dependências dos módulos novos e o motivo de cada uma:

| De → para | Motivo | Exemplo de chamada |
|---|---|---|
| TER → ADM | Usuário de perfil `prestador` fica ligado ao `fornecedor` (`usuario.fornecedor_id`) | `adm.Servico.UsuariosDoFornecedor` |
| TER → ATV | O catálogo de serviços (`servico_fornecedor`) indica classes de ativo e marcas atendidas | `atv.Servico.ValidarClasses` |
| OS → TER | **Executor da OS:** com `executor_tipo = fornecedor`, a OS valida o `fornecedor_id`, o serviço do catálogo e, na Fase 2, o contrato de prestação e o SLA aplicável; também cria e consome o link de uso único | `terceiros.Servico.ValidarExecutor`, `terceiros.Servico.CriarLinkExecucao`, `terceiros.Servico.ConsumirLink` |
| EST → TER | **Estoque em terceiro e reparo:** almoxarifado `em_terceiro` exige `fornecedor_id` de papel `subcontratada`; ordem de reparo exige fornecedor e registra o orçamento em `orcamento_fornecedor` | `terceiros.Servico.ValidarPapel`, `terceiros.Servico.RegistrarOrcamento` |
| OS → EST | Já existia; passa a incluir baixa no almoxarifado `em_terceiro` do executor e envio de peça defeituosa removida para `quarentena` | `est.Servico.BaixarConsumo`, `est.Servico.ReceberDefeituosa` |
| OS → LOC | **Responsável pelo custo:** ao criar a OS de um ativo com contrato de locação vigente, a OS pede à matriz `responsabilidade_locacao` o `responsavel_custo` padrão para o tipo de serviço; o usuário pode alterar com justificativa | `locacao.Servico.ResponsavelCustoPadrao(ativo, tipo, data)` |
| LOC → ATV | Itens do contrato são ativos; ativo `cedido` em locação passa a ter o local de uso como localização; ativo `tomado` tem `propriedade = de_terceiro` | `atv.Servico.DefinirLocalUso` |
| LOC → MED | A vistoria de entrega e de devolução registra leituras de medidor; franquia de horas ou km usa as leituras | `med.Servico.RegistrarLeitura` |
| LOC → TER | Contrato `tomado` exige fornecedor com papel `locador`; o SLA de substituição do locador vem do cadastro dele | `terceiros.Servico.ValidarPapel` |
| LOC → CHK | O checklist da vistoria usa um `checklist_modelo` | `chk.Servico.ModeloDeVistoria` |
| SRV → ADM | **Técnicos da subcontratada:** técnico é `pessoa` + `usuario` com perfil `tecnico`; SRV liga especialidades (`pessoa_especialidade`) e consulta técnicos ativos para distribuir OS | `adm.Servico.PessoasComPerfil`, `adm.Servico.ValidarPessoa` |
| SRV → ATV | **Equipamentos de clientes:** cria e atualiza `ativo` com `propriedade = de_cliente` e `cliente_servico_id`; na conexão, grava a cópia sincronizada com `conexao_empresa_id` e `ativo_remoto_id` e marca como somente leitura os campos de cadastro que pertencem ao cliente | `atv.Servico.CriarAtivoDeCliente`, `atv.Servico.AplicarCadastroRemoto` |
| SRV → TER | **Conexão do lado do cliente:** o convite parte de um `fornecedor` com papel `subcontratada` ou `assistencia_tecnica`; do lado da subcontratada, os itens `servico_os` referenciam o catálogo `servico_fornecedor` da própria conta | `terceiros.Servico.ValidarPapel`, `terceiros.Servico.ServicoDoCatalogo` |
| SRV → OS | **OS vinculada:** cria a OS correspondente na outra conta (na subcontratada, a partir da OS atribuída ao fornecedor conectado; no cliente, `aberta` com `origem = prestador`) e aplica andamento, apontamentos, componentes, fotos, conclusão, liberação e recusa recebidos pela conexão, sempre pelos casos de uso comuns da OS | `ordemservico.Servico.CriarVinculada`, `ordemservico.Servico.AplicarExecucaoRemota`, `ordemservico.Servico.AplicarLiberacaoRemota` |

Dependências novas da versão 0.5 (configuração por segmento, PRJ e CTR), todas **(proposta)**; os nomes das funções são exemplos:

| De → para | Motivo | Exemplo de chamada |
|---|---|---|
| PRJ → OS | **OS de implantação da etapa:** cria OS `implantacao` com `origem = projeto` e a etapa de origem, pelos contratos internos da OS (a OS guarda a etapa como referência, sem importar `projetos`); consulta a situação das OS da etapa para concluí-la (RF-PRJ-003, RF-PRJ-004) | `ordemservico.Servico.CriarImplantacao`, `ordemservico.Servico.SituacaoDasOS` |
| PRJ → ATV | Ativos do projeto, normalmente `em_implantacao`; a conclusão do projeto executa `colocar_em_operacao` (RF-PRJ-007, RF-ATV-013) | `atv.Servico.ColocarEmOperacao` |
| PRJ → TER | Executor `fornecedor` da etapa e valor da etapa vindo do contrato de prestação (RF-PRJ-006) | `terceiros.Servico.ValidarExecutor`, `terceiros.Servico.ValorDoContrato` |
| PRJ → CON | Documentos exigidos por etapa (tipo de documento e documento registrado no ativo depois do início da etapa) e assinatura do termo de entrega | `con.Servico.DocumentoRegistradoDesde`, `con.Servico.Assinar` |
| PRJ → CTR | Projeto ligado a um contrato com cliente; na conclusão, inclui os ativos como itens cobertos a partir da entrada em operação, na mesma transação (RF-PRJ-007, RF-CTR-002) | `contratos.Servico.ValidarContrato`, `contratos.Servico.IncluirItens` |
| PRJ → SRV | Cliente atendido do projeto (opcional), só leitura do cadastro (regra 9) | `servicos.Servico.ClienteAtendido` |
| CTR → ATV | Itens cobertos: propriedade, fase de ciclo de vida, `data_entrada_operacao` e status do ativo (RF-CTR-002) | `atv.Servico.ValidarItemCoberto` |
| CTR → OS | **Lê as paradas da OS:** registros de parada e `causa_externa` das OS dos itens, por período, para a medição e o painel (RF-CTR-003, RF-CTR-004). CTR não grava na OS: a aprovação da exclusão fica na medição | `ordemservico.Servico.ParadasDoAtivo(ativo, periodo)` |
| CTR → SRV | Cliente atendido do contrato, só leitura do cadastro (regra 9) | `servicos.Servico.ClienteAtendido` |
| INT → PRJ | Criação de projetos em lote por planilha (RF-PRJ-002, RF-INT-003) | `projetos.Servico.CriarEmLote` |
| `configuracao` → ATV | Filtro das regras por classe, unidade e criticidade; o pacote cria classes de ativo com campos extras e medidores | `atv.Servico.AtivosDoFiltro`, `atv.Servico.GarantirClasse` |
| `configuracao` → MED | Condição `leitura_fora_faixa` | `med.Servico.Leitura` |
| `configuracao` → CON | Condições `documento_vencido`, `documento_a_vencer` e, na Fase 2, `capacitacao_vencida`; efeitos `bloquear_ativo`, `bloquear_ativo_ate_documento` e, na Fase 2, `impedir_atribuicao` (CON é dono de `bloqueio_ativo` e das capacitações); o pacote cria tipos de documento | `con.Servico.DocumentosVencidos`, `con.Servico.BloquearAtivo`, `con.Servico.GarantirTipoDocumento` |
| `configuracao` → SOL | Efeito `abrir_solicitacao`, com origem `regra` | `sol.Servico.Abrir` |
| `configuracao` → CHK | Condição `item_checklist_reprovado`; o pacote cria modelos de checklist | `chk.Servico.ItemReprovado`, `chk.Servico.GarantirModelo` |
| `configuracao` → PRE | O pacote cria procedimentos e planos padrão | `pre.Servico.GarantirProcedimento`, `pre.Servico.GarantirPlanoPadrao` |
| `configuracao` → OS | Condição `os_concluida_com_campo` (valor do campo extra `sim_nao`); o pacote cria `campo_extra_os` e códigos de falha sugeridos | `ordemservico.Servico.ValorCampoExtra`, `ordemservico.Servico.GarantirCampoExtra` |
| `configuracao` → PRJ | Condição `etapa_projeto_atrasada` (Fase 2); o pacote cria modelos de projeto | `projetos.Servico.EtapasAtrasadas`, `projetos.Servico.GarantirModelo` |
| `configuracao` → ADM | Destinatários do efeito `notificar` (usuários por perfil e unidade) | `adm.Servico.UsuariosComPerfil` |

Reações que **não** são chamada direta, para não criar ciclo:

- Avaria na vistoria de devolução: LOC publica `vistoria_locacao.registrada` com as avarias; OS escuta e cria a OS corretiva com `responsavel_custo = locatario`.
- SLA de fornecedor (Fase 2): TER escuta `ordem_servico.*` (reconhecimento, início, conclusão) e publica `fornecedor.sla_violado` quando o prazo do contrato estoura.
- Retorno de ordem de reparo: EST publica `ordem_reparo.retornada`; quem precisar (IND, webhooks) escuta.
- OS vinculada: SRV escuta `ordem_servico.*`, `apontamento.*` e `anexo.*` (atribuição a fornecedor conectado, início, pausa, conclusão, liberação, recusa, cancelamento) e, se a OS tiver `vinculo_os` com conexão `ativa`, enfileira a mensagem para o sincronizador (seção 6.8). OS não importa `servicos`: quem chama é sempre SRV → OS.
- Conexão: SRV publica `conexao_empresa.ativada` e `conexao_empresa.encerrada`; TER escuta para mostrar o estado da conexão no cadastro do fornecedor; IND e webhooks escutam `vinculo_os.sincronizado` e `vinculo_os.erro`.
- Entrada em operação (versão 0.5): ATV publica `ativo.colocado_em_operacao`; PRE escuta e ativa os planos a partir de `data_entrada_operacao`. ATV não importa `pre` (PRE → OS → ATV fecharia ciclo).
- Andamento da etapa de projeto (Fase 2): PRJ escuta `ordem_servico.*` das OS com `origem = projeto` (início, encerramento, cancelamento) para passar a etapa a `em_andamento` e conferir a conclusão. OS não importa `projetos`.
- Prazos do contrato com cliente (Fase 2): CTR escuta `ordem_servico.*` das OS dos itens cobertos e publica `contrato_cliente.sla_violado` quando o prazo de atendimento ou de solução estoura, como TER faz com o SLA de fornecedor. OS não importa `contratos`.
- Regras configuráveis: `configuracao` escuta os eventos que correspondem às condições do catálogo (`ordem_servico.concluida`, `checklist_execucao.item_reprovado`, `leitura_medidor.registrada`, `documento_vencimento.vencendo`) e aplica o efeito pela API pública do módulo dono do dado (seção 5.6). Nenhum módulo importa `configuracao`.

Regras:

1. Um módulo só chama outro pela interface exportada no pacote raiz (`ordemservico.Servico`).
2. Um módulo nunca lê ou grava tabelas de outro módulo. Cada tabela tem um módulo dono, registrado em `2_ARQUITETURA_BANCO_DADOS.md`.
3. Exceção controlada: IND lê de visões e visões materializadas de leitura (prefixo `vw_`) que juntam dados de vários módulos. Essas visões são somente leitura e têm o mesmo RLS das tabelas de origem.
4. INT, IA e MOB são "módulos de borda": dependem de muitos, e ninguém depende deles.
5. As regras são verificadas no CI pelo linter `depguard` (no `golangci-lint`) com a lista de importações permitidas por módulo. **(proposta)**
6. Escutar um evento não cria dependência de importação: os tipos de payload ficam em `plataforma/eventos/contratos` e o `eventos.go` de cada módulo só os reexporta. Assim TER pode escutar `ordem_servico.*` sem importar `ordemservico`, embora OS chame TER. **(proposta)**
7. Nenhum módulo de domínio chama SRV. SRV fica acima de OS, ATV, TER e ADM e reage a eles por eventos. Assim, uma empresa sem o módulo SRV e sem conexões não executa nenhum código de `servicos` no caminho das suas requisições.
8. Só o pacote `servicos` (código de aplicação de mensagens) roda dentro do sincronizador com o papel `tenaz_sinc`, e só chama as funções de OS e ATV marcadas como "remotas" (`AplicarExecucaoRemota`, `AplicarCadastroRemoto`...), que aceitam apenas os campos da lista fechada (seção 6.8). O `depguard` impede que outros pacotes importem o cliente de banco do sincronizador.
9. **Exceção estreita à regra 7 (versão 0.5).** O sumário (seção 3.1) faz `cliente_servico` existir nas contas com SRV e nas contas com CTR. Enquanto o dono da tabela continuar sendo SRV (pendência na seção 16), CTR e PRJ podem chamar de `servicos` **só** a leitura e o cadastro de `cliente_servico` (`servicos.Servico.ClienteAtendido`), que não toca conexão, vínculo nem sincronizador. O `depguard` restringe essa importação a esse arquivo de API. **(proposta)**
10. **`configuracao` é pacote de borda**, como INT, IA e MOB (regra 4): depende de ATV, MED, SOL, CHK, PRE, CON, OS, PRJ e ADM, e ninguém depende dele. Reage por eventos e grava sempre pela API pública do módulo dono do dado; nunca grava em tabela de outro módulo (regra 2). **(proposta)**
11. **Verificação sem importação.** Quando um módulo de baixo precisa confirmar uma condição de um módulo de cima antes de uma ação síncrona, ele declara uma interface no seu pacote raiz e o módulo de cima a implementa; a ligação é feita em `cmd/tenaz`, na inicialização. Primeiro caso: `colocar_em_operacao` (RF-ATV-013) exige que não haja OS `implantacao` não encerrada e que os documentos críticos da classe estejam válidos. ATV declara `atv.VerificadorEntradaOperacao`; OS e CON o implementam. As importações continuam OS → ATV e CON → ATV, a chamada roda na mesma transação e o grafo continua sem ciclo. **(proposta, validar em código)**
12. **Credenciais só pela autenticação (versão 0.6).** Só `plataforma/autenticacao` alcança as tabelas `identidade`, `acesso_empresa`, `sessao`, `token_renovacao`, `codigo_autorizacao`, `desafio_mfa` e `token_email`, pelo repositório em `autenticacao/internal/repo` (o compilador impede outro importador) e pelo pool do papel `tenaz_autenticacao`. Módulos de negócio, inclusive ADM, chamam a API do serviço de autenticação (`EmitirConvite`, `EncerrarSessoesUsuario`...) **fora** da própria transação, depois do `COMMIT`; o repositório recusa ser chamado dentro de `NaEmpresa`. Um teste falha se algum arquivo `.sql` fora desse repositório citar uma das sete tabelas ([ADR-011](#adr-011--papel-de-banco-próprio-para-as-credenciais-de-autenticação), seção 6.9). Não é dependência nova entre módulos: módulos já podem importar a plataforma; a regra só fecha o caminho até o banco.

### 5.4 Como um módulo chama outro

Chamada síncrona, no mesmo processo e **na mesma transação**:

- A transação de banco, já com o contexto de empresa definido, viaja no `context.Context`.
- O repositório de qualquer módulo obtém a transação do contexto. Ele nunca abre conexão própria.
- Se o caso de uso de OS chama `est.Servico.ReservarItens(ctx, ...)` e a OS falha depois, a reserva é desfeita no mesmo rollback.

Exemplo de fluxo: técnico conclui uma OS de máquina NR-12.

1. `ordemservico/app` carrega a OS e valida a transição `em_execucao → concluida` no domínio.
2. Chama `con.Servico.ValidarRegistroNR12(ctx, os)` para conferir os sete campos do item 12.11.2.
3. Chama `est.Servico.BaixarConsumo(ctx, itens)`.
4. Grava a OS e publica o evento `ordem_servico.concluida` no outbox.
5. Confirma a transação. Tudo ou nada.

### 5.5 Eventos internos e outbox transacional

Usado quando a reação pode ser assíncrona ou quando a chamada direta criaria ciclo.

```mermaid
sequenceDiagram
    participant UC as Caso de uso (OS)
    participant PG as PostgreSQL
    participant W as Worker River
    participant S as Assinantes

    UC->>PG: BEGIN, altera ordem_servico
    UC->>PG: INSERT evento_dominio (ordem_servico.concluida)
    UC->>PG: INSERT job despachar_evento (mesma transação)
    UC->>PG: COMMIT
    W->>PG: pega job (SKIP LOCKED)
    W->>S: um job por assinante (CON, IND, webhooks, notificações)
    S->>PG: grava efeito + evento_processado (idempotência)
```

| Item | Definição **(proposta)** |
|---|---|
| Tabela | `evento_dominio` (append-only): `empresa_id`, `id` (UUIDv7), `tipo`, `agregado`, `agregado_id`, `versao_schema`, `payload` (jsonb), `ocorrido_em`, `usuario_id`, `origem` |
| Nome do evento | `<entidade>.<fato_no_particípio>`, com o nome completo da entidade (sumário, seção 3.4): `ordem_servico.concluida`, `ordem_servico.encerrada`, `leitura_medidor.registrada`, `plano_manutencao.vencido`, `checklist_execucao.item_reprovado`, `documento_vencimento.vencendo`, `solicitacao.aprovada`. Novos na versão 0.3: `ordem_reparo.enviada`, `ordem_reparo.retornada`, `movimentacao_estoque.registrada`, `contrato_locacao.ativado`, `contrato_locacao.encerrado`, `vistoria_locacao.registrada`, `fornecedor.sla_violado`. Novos na versão 0.4 (SRV): `conexao_empresa.ativada`, `conexao_empresa.encerrada`, `vinculo_os.sincronizado`, `vinculo_os.erro` |
| Entrega | Pelo menos uma vez. Todo assinante é idempotente: grava `evento_processado(empresa_id, evento_id, assinante)` com chave única na mesma transação do efeito |
| Ordem | Garantida por agregado (mesmo `agregado_id`), não global |
| Contexto de empresa | O job carrega `empresa_id`; o worker abre transação com esse contexto antes de qualquer leitura |
| Webhooks públicos | São um assinante. Entram no **MVP**, no plano Profissional (entitlement `webhooks`). Reusam o nome do evento interno quando ele é estável e público; o payload público é versionado à parte. Eventos de um escopo externo (prestador, locatário) nunca são entregues a destinos do prestador ou do locatário: webhooks são só da empresa |
| Retenção | 90 dias na tabela quente; depois, arquivado no object storage |

Exemplos de uso: `checklist_execucao.item_reprovado` → SOL cria solicitação; `ordem_servico.encerrada` → IND atualiza custos consolidados (separados por `responsavel_custo`) e INT dispara webhook; `leitura_medidor.registrada` → PRE reavalia gatilhos por medidor; `vistoria_locacao.registrada` com avaria → OS cria corretiva com custo do locatário; `contrato_locacao.encerrado` → IND fecha o relatório de cobrança de avarias e horas excedentes e INT dispara webhook (a volta do ativo cedido à localização de origem é chamada direta LOC → ATV, na mesma transação do encerramento); `ordem_servico.concluida` numa OS com `vinculo_os` → SRV enfileira a mensagem de conclusão para a conta do cliente; `vinculo_os.erro` → notificação ao gestor das duas contas.

Os eventos de uma conta nunca são entregues a assinantes de outra conta. O que atravessa contas não é o evento, e sim uma **mensagem de conexão** montada pelo sincronizador só com os campos permitidos (seção 6.8).

**Eventos citados pelos requisitos da versão 0.5** (nomes de [1_ESPECIFICACAO_REQUISITOS.md](1_ESPECIFICACAO_REQUISITOS.md); a lista pública de webhooks é de `3_ESPECIFICACAO_APIS.md`): `ativo.colocado_em_operacao` (RF-ATV-013); na Fase 2, `projeto.iniciado`, `projeto.concluido`, `projeto.cancelado` (RF-PRJ-002), `contrato_cliente.ativado`, `contrato_cliente.encerrado`, `contrato_cliente.sla_violado` (RF-CTR-001) e `medicao_contrato.fechada` (RF-CTR-004). O disparo de uma regra não é evento novo: fica em `disparo_regra`, e o efeito publica os eventos do módulo dono (por exemplo, a mudança de status do ativo bloqueado). Eventos gerados por efeito de regra levam `origem = regra` (seção 5.6).

### 5.6 Configuração por segmento

Sumário, seção 3.7; requisitos RF-ADM-016, RF-ADM-017, RF-ATV-013 e RF-OS-021 em [1_ESPECIFICACAO_REQUISITOS.md](1_ESPECIFICACAO_REQUISITOS.md). Decisão e alternativas no [ADR-010](#adr-010--configuração-por-segmento-com-catálogos-fechados). A empresa configura **blocos que o código já conhece**; não cria estados, telas, fluxos nem colunas.

**Núcleo canônico (não configurável).**

| Item | Onde vive | Por que é fixo |
|---|---|---|
| Estados (OS, solicitação, projeto, etapa, contratos, medição, conexão) | Domínio Go e `pacotes/dominio`, com casos JSON compartilhados (seção 7) | O app offline resolve conflitos por estado (seção 10.5) e a conexão entre contas define dono por estado (seção 6.8) |
| Fórmulas dos indicadores | IND, sobre as visões `vw_` (regra 3 da seção 5.3) | Indicadores comparáveis entre empresas e entre segmentos. A exclusão da OS `implantacao` e do ativo `em_implantacao` é regra fixa da visão, não configuração. A disponibilidade contratual de CTR é cálculo próprio de `contratos` e não altera a disponibilidade técnica de IND |
| Regras do app offline | Seções 10.3 a 10.5 | Uma regra configurável não roda no aparelho |
| Isolamento | Seção 6 | Nenhum bloco cria tabela, política ou escopo |
| Campos da conexão entre contas | `campos_conexao.yaml` (seção 6.8) | Campos extras nunca entram na lista no lançamento (sumário 3.6) |

**Blocos configuráveis e dono de cada um.**

| Bloco | Pacote dono | Onde fica | Fase |
|---|---|---|---|
| Classes de ativo | `atv` | Tabelas da empresa (já existiam) | MVP |
| Campos extras da OS | `ordemservico` | Definição em `campo_extra_os`; valores em `ordem_servico.campos_extras` (jsonb) | MVP |
| Tipos de documento | `con` | `tipo_documento` (já existia) | MVP |
| Regras de condição e efeito | `configuracao` | `regra_automacao` e `disparo_regra` | MVP |
| Fase de ciclo de vida | `atv` | `ativo.fase_ciclo_vida` e `data_entrada_operacao` | MVP |
| Modelos de projeto | `projetos` | `modelo_projeto` e `modelo_etapa_projeto` | Fase 2 |
| Contratos com clientes | `contratos` | `contrato_cliente`, `item_contrato_cliente` e `faixa_desconto_contrato` | Fase 2 |
| Pacotes de segmento | `configuracao` | `pacote_segmento` (catálogo global) e `aplicacao_pacote` (por empresa) | Fase 2 |

Toda tabela de configuração da empresa tem `empresa_id`, RLS forçado e a política `escopo_externo` padrão (só `interno`); a única que entra na lista de acesso externo é `campo_extra_os` (seção 6.6). `pacote_segmento` é a única tabela global nova.

**Catálogos fechados no código** **(proposta)**.

- Cada condição e cada efeito é um tipo Go do pacote `configuracao` que implementa a interface `Condicao` (código, validação dos parâmetros, busca de candidatos) ou `Efeito` (código, validação dos parâmetros, aplicação). O registro é uma lista fixa, compilada no binário: não há plugin, carga dinâmica, expressão nem script.
- No banco, `regra_automacao.condicao` e `regra_automacao.efeito` só aceitam os códigos do catálogo (restrição `CHECK` ou tipo enumerado, a decidir no documento 2). Código novo exige migração, requisito em [1_ESPECIFICACAO_REQUISITOS.md](1_ESPECIFICACAO_REQUISITOS.md) e teste no mesmo PR (risco R7).
- Os parâmetros da condição e do efeito ficam em jsonb e são decodificados no struct do tipo, com rejeição de campo desconhecido.
- A matriz de combinações permitidas (RF-ADM-016) é uma tabela no código. O CI falha se uma combinação permitida não tiver caso de teste.
- Os códigos da Fase 2 (`capacitacao_vencida`, `etapa_projeto_atrasada`, `impedir_atribuicao`) só entram no catálogo quando o requisito for implementado.

**Avaliação das regras.**

```mermaid
sequenceDiagram
    participant UC as Caso de uso (OS, CHK, MED, CON)
    participant PG as PostgreSQL
    participant AG as Agendador River
    participant W as Worker (fila regras)
    participant MD as Módulo dono do efeito

    UC->>PG: T1 grava o fato + evento_dominio + job (outbox)
    W->>PG: T2 em NaEmpresa, regras ativas com a condição do evento e filtro que casa
    W->>PG: INSERT disparo_regra (regra, chave do fato) ON CONFLICT DO NOTHING
    W->>MD: só se o disparo foi inserido, aplica o efeito pela API pública (mesma T2)
    MD->>PG: bloqueio_ativo, solicitacao ou notificação, com origem regra
    AG->>W: rotina diária, um job por empresa ativa
    W->>PG: T3 em NaEmpresa, candidatos de cada regra ativa sem disparo_regra
    W->>MD: aplica o efeito dos candidatos novos (mesma T3)
```

- **Por evento.** Para as condições ligadas a um fato (`os_concluida_com_campo`, `item_checklist_reprovado`, `leitura_fora_faixa`), o assinante de `configuracao` avalia a regra no job do outbox, logo depois do fato. O registro em `disparo_regra` e o efeito vão na **mesma transação**: ou os dois existem, ou nenhum.
- **Rotina diária idempotente.** Segue o padrão da seção 6.4 (um job por empresa ativa, cada um em `NaEmpresa`, sem `BYPASSRLS`). Avalia as condições de tempo (`documento_vencido`, `documento_a_vencer` e, na Fase 2, `capacitacao_vencida` e `etapa_projeto_atrasada`) e repassa as condições de fato, para recuperar uma avaliação perdida ou que falhou. Roda depois da rotina de alertas de vencimento e termina no mesmo prazo (RNF-009).
- **Idempotência pela chave do fato.** `disparo_regra` tem chave única `(empresa_id, regra_automacao_id, chave_fato)`. A chave identifica o fato, não o dia: documento (`documento_vencimento_id` e validade), OS (`ordem_servico_id`), item de checklist (`checklist_execucao_id` e item), leitura (`leitura_medidor_id`), capacitação (`capacitacao_id` e validade), etapa (`etapa_projeto_id` e data prevista). Um documento renovado que vence de novo dispara de novo; o mesmo vencimento, nunca. **(proposta; colunas finais no documento 2)**
- **Efeito pelo dono do dado.** `bloquear_ativo` e `bloquear_ativo_ate_documento` chamam CON, dono de `bloqueio_ativo` (`origem = regra`); `abrir_solicitacao` chama SOL (origem `regra`); `notificar` usa o mesmo caminho de envio dos alertas de vencimento. Em `bloquear_ativo_ate_documento`, o bloqueio guarda o tipo de documento exigido e a data do fato; quem levanta o bloqueio é o próprio CON, no caso de uso que registra o documento, sem nova avaliação de regra.
- **Sem cascata.** Os eventos gerados por um efeito levam `origem = regra`, e o assinante de `configuracao` os ignora. Uma regra nunca dispara outra. O teste da matriz de combinações confere isso para cada condição nova.
- **Simular sem gravar.** A simulação (RF-ADM-016) usa a mesma busca de candidatos, numa transação somente leitura na réplica (seção 4.2), sem `disparo_regra` e sem efeito.
- **Inativar.** Inativar a regra não levanta bloqueios; a lista dos bloqueios ainda ativos vem de `disparo_regra` ligado a `bloqueio_ativo` com `origem = regra`.
- **Auditoria.** Todo efeito grava `evento_auditoria` com `origem = regra` e o ID da regra **(proposta)**. Os jobs rodam no escopo `interno`.
- **Conexão entre contas.** A regra é avaliada só na conta onde o evento aconteceu, e o efeito grava só nessa conta. Um evento aplicado pelo sincronizador (`origem = conexao`) é avaliado normalmente na conta de destino, mas os campos extras não trafegam: no lançamento, `os_concluida_com_campo` não dispara na conta do cliente para a OS executada por uma subcontratada conectada (pendência na seção 16).

**Campos extras da OS.**

- A definição (`campo_extra_os`) é de OS; os valores ficam em `ordem_servico.campos_extras`, um objeto jsonb indexado pela chave da definição. Nenhuma empresa ganha coluna própria.
- A validação é **uma função** do domínio da OS em Go e a mesma função em `pacotes/dominio` (TypeScript), com os mesmos casos JSON em `backend/testes/casos/` (seção 7): chave conhecida na empresa, tipo, opções da lista, obrigatoriedade na abertura ou na conclusão para o tipo de OS e a classe do ativo. Chave desconhecida é recusada; valor de definição inativada é preservado.
- O servidor é a autoridade. A web e o app validam antes de enviar; o app valida **antes de enfileirar o comando** (seção 10.3).
- Filtro e exportação por campo extra usam operadores jsonb com índice GIN **(proposta; documento 2)**. O limite de 30 campos ativos por combinação (RF-OS-021) limita o tamanho do payload e do pull.
- Nenhuma fórmula de indicador lê `campos_extras`.
- Escopo externo: o prestador lê as definições que se aplicam às OS do seu fornecedor; o link de uso único, só as marcadas como do prestador **(proposta)**. A projeção para escopo externo leva só os valores desses campos.
- Conexão entre contas: `campos_extras` não está em `campos_conexao.yaml` e o teste com três contas confere (seção 9.3).

**Pacotes de segmento** (Fase 2) **(proposta)**.

- **Fonte versionada.** Cada pacote é um arquivo no repositório (`backend/internal/modulos/configuracao/pacotes/<codigo>.yaml`), revisado em PR e validado no CI: condições e efeitos do catálogo, combinações permitidas, regras nascendo `inativa`, chaves naturais únicas, tipos de campo válidos.
- **Catálogo global.** `tenaz migrate` publica cada versão em `pacote_segmento`, tabela sem `empresa_id`, na lista de tabelas globais do teste da seção 9.3, com só `SELECT` para `tenaz_app`; escrita só pelo papel `tenaz_dono`, nas migrações. Nenhuma conta grava no catálogo.
- **Cópia idempotente que só cria o que falta.** Um job `aplicar_pacote`, único por empresa, pacote e versão, roda numa transação em `NaEmpresa`. Percorre os itens na ordem de dependência (classes, campos extras, tipos de documento, procedimentos, planos, checklists, modelos de projeto, regras) e, para cada um, chama a função `Garantir...` do módulo dono, que cria o item só se a chave natural não existir na conta (chave única por empresa e `ON CONFLICT DO NOTHING`) e devolve "criado" ou "ignorado". Nunca altera nem apaga. Item de módulo fora do plano é "ignorado". O resultado vai para `aplicacao_pacote`.
- **Prévia** ("ver antes de aplicar"): o mesmo percurso numa transação somente leitura.
- **Cópia, não referência.** Depois de aplicado, nada na conta aponta para o pacote em tempo de execução: o que foi criado é configuração comum da empresa. Uma versão nova do pacote não muda nenhuma conta até ser aplicada, e mesmo assim só cria o que falta.

---

## 6. Isolamento entre empresas no código

O modelo é pool (tabelas compartilhadas) com RLS. A decisão e as alternativas estão no [ADR-002](#adr-002--multiempresa-em-pool-com-rls). Esta seção mostra a mecânica.

### 6.1 Camadas de defesa

| Camada | Mecanismo |
|---|---|
| Esquema | Toda tabela de dados de cliente tem `empresa_id uuid NOT NULL`, índice começando por `empresa_id` e chaves estrangeiras compostas `(empresa_id, id)`. Uma OS não consegue apontar para um ativo de outra empresa nem por erro de código. |
| Papéis do banco | Migrações rodam com `tenaz_dono`, dono das tabelas, exceto as de credencial de `plataforma`, que pertencem a `tenaz_credencial`, sem membros de runtime (seção 6.9). A aplicação conecta como `tenaz_app`: sem `BYPASSRLS`, sem superusuário, sem ser dona de tabela. Dois papéis de propósito único ficam fora dos privilégios padrão e têm credencial própria: `tenaz_sinc` (seção 6.8) e `tenaz_autenticacao` (seção 6.9). |
| Credenciais de login (versão 0.6) | As tabelas de autenticação do esquema `plataforma` ficam fora do RLS, porque são lidas antes de a empresa ser conhecida. Sem RLS, o privilégio é a defesa: só o papel `tenaz_autenticacao`, usado só pelo pacote de autenticação, as alcança; `tenaz_app`, `tenaz_job`, `tenaz_leitura` e `tenaz_sinc` recebem "permissão negada". O dono também é defesa: é `tenaz_credencial`, e não `tenaz_dono`, para que o `tenaz_manutencao`, que herda do dono comum, não as leia (achado ISO-E01-08) ([ADR-011](#adr-011--papel-de-banco-próprio-para-as-credenciais-de-autenticação)). |
| RLS | `ENABLE` e `FORCE ROW LEVEL SECURITY` em todas as tabelas com `empresa_id`, com política de leitura (`USING`) e de escrita (`WITH CHECK`). |
| Transação | O contexto de empresa é definido com `set_config('app.empresa_id', $1, true)` dentro da transação. Nunca `SET` de sessão, que vazaria entre requisições no pool. |
| Repositório | Toda consulta sqlc recebe `empresa_id` como parâmetro e filtra por ele, mesmo com RLS. Um teste lê os arquivos `.sql` e falha se alguma consulta sobre tabela de cliente não tiver `empresa_id`. |
| Sem contexto, sem dados | Se `app.empresa_id` não estiver definido, a política retorna zero linhas. Falha fechada. |
| Fora do banco | Chaves de cache e de object storage começam por `empresa_id`. Jobs carregam `empresa_id`. Logs e traces carregam `empresa_id`. |
| Escopo externo | Para prestador, locatário e link de uso único, uma segunda política RLS **restritiva** limita as linhas ao fornecedor, ao locatário ou à OS do link, dentro da empresa (seção 6.6). |
| Conexão entre contas | O sincronizador usa o papel `tenaz_sinc`, também sem `BYPASSRLS`, com uma política restritiva própria (só linhas ligadas à conexão em processamento) e `GRANT` por coluna (só os campos da lista fechada). Cada gravação acontece numa transação de uma única empresa (seção 6.8). |
| Teste | Teste de vazamento em todo build, entre empresas, entre escopos externos e entre contas conectadas, com três contas (seção 9.3). |

### 6.2 Política e função de contexto

```sql
CREATE FUNCTION app.empresa_atual() RETURNS uuid
  LANGUAGE sql STABLE
  AS $$ SELECT NULLIF(current_setting('app.empresa_id', true), '')::uuid $$;

ALTER TABLE ordem_servico ENABLE ROW LEVEL SECURITY;
ALTER TABLE ordem_servico FORCE ROW LEVEL SECURITY;

CREATE POLICY isolamento_empresa ON ordem_servico
  USING      (empresa_id = (SELECT app.empresa_atual()))
  WITH CHECK (empresa_id = (SELECT app.empresa_atual()));
```

Detalhes que importam:

- `NULLIF(..., '')`: depois que uma transação termina, a configuração volta a texto vazio, não a nulo. Sem o `NULLIF`, o cast falharia com erro em vez de devolver zero linhas.
- `(SELECT app.empresa_atual())`: o subselect faz o Postgres avaliar a função uma vez por consulta, e não uma vez por linha.

### 6.3 Transação com contexto de empresa no Go

```go
// NaEmpresa abre uma transação, define a empresa e coloca a transação no contexto.
// Todo caso de uso que toca dados de cliente passa por aqui.
func (db *DB) NaEmpresa(ctx context.Context, empresaID uuid.UUID, fn func(ctx context.Context) error) error {
    return pgx.BeginTxFunc(ctx, db.pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
        if _, err := tx.Exec(ctx, `SELECT set_config('app.empresa_id', $1, true)`, empresaID.String()); err != nil {
            return fmt.Errorf("definir empresa: %w", err)
        }
        return fn(comTransacao(ctx, tx))
    })
}
```

O middleware HTTP extrai a empresa **e o escopo** (seção 6.6) da sessão autenticada e chama `NaEmpresa`. Nenhum handler recebe `empresa_id`, `fornecedor_id` ou `locatario_id` do corpo ou da URL para decidir o escopo. Na versão com escopo, `NaEmpresa` recebe também um `Escopo` (tipo, usuário e, conforme o tipo, fornecedor, locatário ou hash do link) e define `app.escopo`, `app.usuario_id`, `app.fornecedor_id`, `app.locatario_id` e `app.link_hash` na mesma transação, com o mesmo `set_config(..., true)`, seguido de `app.conferir_escopo()` (documento 2, seções 2.4 e 2.9.3). Não existe uma variável genérica `app.escopo_id` (versão 0.6, seção 6.6).

**No código (versão 0.8, `internal/plataforma/db/db.go`).** `NaEmpresa` fixa **17 variáveis** de transação, todas com `set_config(..., true)`: numa instrução, `app.empresa_id`, `app.escopo`, `app.usuario_id`, `app.fornecedor_id`, `app.locatario_id`, `app.link_hash`, `app.ator_tipo`, `app.origem`, `app.ip`, `app.request_id`, `app.user_agent`, `app.token_api_id`, `app.dispositivo_id`, `app.ocorrido_no_dispositivo_em`, `app.modulo` e `app.tela`; numa segunda, `app.usuario_nome`, lido do cadastro do usuário naquela empresa (a leitura passa pelo RLS e precisa da empresa e do escopo já fixados). Em seguida, `app.conferir_escopo()`. Tudo num lote só, uma ida ao banco. `db.SemTransacao(ctx)` devolve o contexto sem a transação, para abrir **outra** transação `NaEmpresa`, em outra conexão, com empresa e escopo próprios: existe só para ler configuração que o escopo da requisição não alcança, com cache. O único uso é a captura da trilha (`captura_auditoria`), lida em transação interna à parte, porque a requisição de prestador, locatário ou link não lê a tabela, e guardada em cache por instância por 30 segundos, invalidado pela alteração feita na própria instância (outra instância vê a mudança em até 30 segundos). Enquanto dura, a leitura segura uma segunda conexão do pool; por isso, só com cache.

**De onde vem a empresa (versão 0.6).** A sessão é resolvida em três passos, nunca na mesma transação, porque o primeiro usa outro papel de banco (seção 6.9):

1. **Pool da autenticação** (`tenaz_autenticacao`), transação curta sem empresa: acha a sessão pelo sha256 do cookie (web) ou pelo `sid` do access token (app), confere canal, revogação, expiração e inatividade e atualiza `ultimo_uso_em`. Devolve só `sessao_id`, `empresa_id`, `usuario_id`, `identidade_id` e `canal`.
2. **Pool da aplicação** (`tenaz_app`), `NaEmpresa(empresa da sessão, escopo interno, sem usuário)`: lê o usuário, os perfis, as permissões e a empresa (documento 2, seção 2.9.2) e confere que `usuario.identidade_id` é a identidade da sessão e que o usuário está `ativo`. Divergência: 401.
3. **Pool da aplicação**, `NaEmpresa(empresa da sessão, escopo do cadastro, com usuário)`: o caso de uso, com `app.conferir_escopo()` conferindo de novo.

O `tenaz_app` nunca lê a tabela de onde a empresa veio; recebe a empresa já resolvida, do mesmo processo.

### 6.4 Rotinas que atravessam empresas

Algumas rotinas precisam percorrer todas as empresas (exemplo: verificar planos preventivos vencidos). Elas **não** usam um papel com `BYPASSRLS`. O padrão é:

1. Um job agendador chama a função `app.empresas_ativas()`, `SECURITY DEFINER`, que devolve **apenas os IDs** das empresas ativas.
2. Para cada empresa, enfileira um job com `empresa_id`.
3. Cada job roda dentro de `NaEmpresa`, com RLS ativo.

As únicas exceções com acesso amplo são as migrações (usuário `tenaz_migrador`, que assume `tenaz_dono` e, para as credenciais, `tenaz_credencial`) e o suporte de banco em emergência, com acesso registrado (ver 9_SEGURANCA_CONFIABILIDADE.md).

**Versão 0.8: o selo da trilha diverge do passo 1.** O `tenaz selar-auditoria` (seção 4.1) não chama `app.empresas_ativas()`: lista as empresas com `SELECT id FROM app.empresa`, fora de qualquer transação de empresa, pela política `empresa_listagem_job` da migração `00004`, que só libera a leitura ao `tenaz_job` quando `app.empresa_atual()` é nulo; e percorre as empresas em sequência, no mesmo processo, em vez de enfileirar um job por empresa (não há fila ainda). Os passos 3 e o "sem `BYPASSRLS`" valem. A diferença é que a política libera a linha inteira de `app.empresa`, e não só o ID, ao `tenaz_job` sem contexto. Pendência na seção 16.

A sincronização entre contas conectadas (SRV) segue o mesmo padrão e **não** é exceção: ela toca duas empresas, mas nunca na mesma transação. Lê na conta de origem numa transação com o contexto da origem e grava na conta de destino em outra transação, com o contexto do destino (seção 6.8).

### 6.5 Banco dedicado (silo) no Enterprise

Previsto para a Fase 3. O mesmo binário e o mesmo esquema rodam contra outro banco. Um registro `empresa.destino_banco` indica o pool de conexões. O RLS continua ativo no banco dedicado. **(proposta)**

### 6.6 Escopo externo: prestador e locatário

Risco R4 do sumário. O RLS por empresa impede que uma empresa veja outra, mas não basta para usuários de empresas externas: dentro da mesma empresa, um prestador não pode ver OS de outro prestador, custos internos ou o resto do cadastro, e um locatário só pode ver os ativos locados a ele. Para isso existe uma **segunda camada de isolamento**, aplicada sempre junto com a da empresa. Decisão no [ADR-008](#adr-008--acesso-externo-por-escopo-restritivo-e-link-de-uso-único).

**Escopos.**

| `app.escopo` | Quem | Variável do escopo (versão 0.6) | De onde vem |
|---|---|---|---|
| `interno` | Usuários com perfis da empresa, jobs, chaves de API | nenhuma (`app.usuario_id` vazio em jobs) | Sessão, job ou chave |
| `prestador` | Usuário com perfil `prestador` | `app.fornecedor_id` = `usuario.fornecedor_id`, lida por `app.fornecedor_atual()` | Cadastro do usuário no banco |
| `locatario` | Usuário com perfil `locatario` | `app.locatario_id` = `usuario.locatario_id`, lida por `app.locatario_atual()` | Cadastro do usuário no banco |
| `link` | Portador de um link de uso único | `app.link_hash` = sha256 do segredo apresentado, lida por `app.link_hash_atual()`; sem usuário | `link_execucao` validado (6.7) |

**Por que não há `app.escopo_id` (versão 0.6).** As versões até 0.5 usavam uma variável genérica `app.escopo_id`. Ela não existe no banco: cada escopo tem a sua variável e a sua função, já nas migrações (documento 2, seção 2.9.2). Motivos: `app.conferir_escopo()` compara cada uma com a coluna própria do cadastro do usuário, e um identificador genérico misturaria fornecedor, locatário e OS no mesmo campo, de modo que um erro de tipo de escopo passaria despercebido; e o escopo `link` é resolvido pelo hash do segredo apresentado, sem um segundo identificador de OS no contexto, que o código poderia definir sem ter o segredo. Decisão aprovada pela revisão de isolamento (pendência E01-3 do documento 2), com duas condições, ambas aplicadas aqui: esta seção cita as três funções, e a política do link confere, na própria subconsulta, que o link não está revogado, expirado nem usado.

Regras:

1. O escopo é resolvido no servidor a partir do banco, a cada requisição, como as permissões. Nunca vem do token, do corpo ou da URL.
2. Perfis `prestador` e `locatario` são exclusivos: o usuário não acumula perfil interno. A regra anti-escalada impede que um administrador conceda permissão interna a usuário externo.
3. A mesma empresa externa que atende duas contratantes tem um cadastro de `fornecedor` (ou de `locatario`) e um usuário em cada empresa contratante, com um token por empresa, como qualquer usuário com acesso a duas empresas. Nada liga os dois cadastros no banco.
4. Sem `app.escopo` definido, as políticas restritivas devolvem zero linhas. Falha fechada, igual à empresa.
5. O escopo `prestador` é para a subcontratada **sem** conta própria. A subcontratada com conta Tenaz Serviços (SRV) não usa escopo externo: seus técnicos e gestores são usuários `interno` da conta dela, e o que vem do cliente chega como cópia pela conexão (seção 6.8). Os dois caminhos convivem: um cliente pode ter uma subcontratada conectada e outra que usa o perfil `prestador`.
6. Os valores de escopo são só estes quatro (sumário 3.6). O sincronizador não é um escopo: é um papel de banco distinto (`tenaz_sinc`), com política própria (seção 6.8).

**No banco.** Toda tabela com `empresa_id` recebe, além da política permissiva `isolamento_empresa`, uma política **restritiva** `escopo_externo`. Como o PostgreSQL combina políticas restritivas com `AND`, a linha precisa passar nas duas. O padrão é "só escopo interno"; apenas as tabelas da lista abaixo têm predicado para escopos externos.

```sql
CREATE FUNCTION app.escopo_atual() RETURNS text
  LANGUAGE sql STABLE
  AS $$ SELECT NULLIF(current_setting('app.escopo', true), '') $$;

-- Versão 0.6: uma função por escopo, como no documento 2 (seção 2.9.2); não existe app.escopo_id()
CREATE FUNCTION app.fornecedor_atual() RETURNS uuid
  LANGUAGE sql STABLE
  AS $$ SELECT NULLIF(current_setting('app.fornecedor_id', true), '')::uuid $$;

CREATE FUNCTION app.locatario_atual() RETURNS uuid
  LANGUAGE sql STABLE
  AS $$ SELECT NULLIF(current_setting('app.locatario_id', true), '')::uuid $$;

CREATE FUNCTION app.link_hash_atual() RETURNS bytea   -- sha256 do segredo do link
  LANGUAGE sql STABLE
  AS $$ SELECT decode(NULLIF(current_setting('app.link_hash', true), ''), 'hex') $$;

-- Padrão para tabelas sem acesso externo (ex.: pedido_compra)
CREATE POLICY escopo_externo ON pedido_compra AS RESTRICTIVE
  USING      ((SELECT app.escopo_atual()) = 'interno')
  WITH CHECK ((SELECT app.escopo_atual()) = 'interno');

-- OS: cada escopo vê só o que lhe foi atribuído
CREATE POLICY escopo_externo ON ordem_servico AS RESTRICTIVE
  USING (
    CASE (SELECT app.escopo_atual())
      WHEN 'interno'    THEN true
      WHEN 'prestador'  THEN executor_tipo = 'fornecedor'
                         AND fornecedor_id = (SELECT app.fornecedor_atual())
      WHEN 'link'       THEN id IN (SELECT l.ordem_servico_id FROM link_execucao l
                                     WHERE l.token_hash  = (SELECT app.link_hash_atual())
                                       AND l.revogado_em IS NULL      -- não revogado
                                       AND l.usado_em    IS NULL      -- não consumido
                                       AND l.expira_em   > now()      -- não expirado
                                       AND l.excluido_em IS NULL)
      WHEN 'locatario'  THEN app.ativo_locado_ao_escopo(ativo_id, data_abertura)
      ELSE false
    END)
  WITH CHECK (
    CASE (SELECT app.escopo_atual())
      WHEN 'interno'    THEN true
      WHEN 'prestador'  THEN executor_tipo = 'fornecedor'
                         AND fornecedor_id = (SELECT app.fornecedor_atual())
      WHEN 'link'       THEN id IN (SELECT l.ordem_servico_id FROM link_execucao l
                                     WHERE l.token_hash  = (SELECT app.link_hash_atual())
                                       AND l.revogado_em IS NULL
                                       AND l.usado_em    IS NULL
                                       AND l.expira_em   > now()
                                       AND l.excluido_em IS NULL)
      ELSE false        -- locatário não grava OS; abre solicitação
    END);
```

- O `WITH CHECK` impede que o prestador reatribua a OS a outro fornecedor ou a converta em execução interna.
- **Link (versão 0.6).** A subconsulta passa também pelas políticas de `link_execucao` (empresa da transação e escopo `link`), e a unicidade `(empresa_id, token_hash)` garante no máximo uma linha. Como a política confere revogação, expiração e consumo a cada consulta, um link revogado deixa de ver a OS na hora, sem depender do código da rota. Consequência para `registrar-execucao`: marcar `link_execucao.usado_em` tem de ser a **última** escrita da transação, depois das escritas na OS, nos apontamentos e nos anexos; antes disso, as escritas seguintes falhariam no `WITH CHECK`. As mesmas quatro condições valem para as demais tabelas da lista abaixo que aceitam o escopo `link`. O texto definitivo da política é do documento 2, na etapa do link (TER).
- `app.ativo_locado_ao_escopo(ativo_id, data)` é uma função SQL `STABLE` (não `SECURITY DEFINER`, portanto sujeita ao RLS) que verifica em `item_contrato_locacao` e `contrato_locacao` se o ativo estava cedido ao locatário do escopo (`app.locatario_atual()`) naquela data, entre a entrega e a devolução. Índice `(empresa_id, ativo_id)` em `item_contrato_locacao`. Se o custo da consulta pesar, a alternativa é gravar o locatário na OS no momento da criação (campo a propor no documento 2). **(proposta)**

**Tabelas com acesso externo** (lista fechada, revisada em code review e conferida pelo teste da seção 9.3; qualquer outra tabela devolve zero linhas a escopo externo):

| Tabela | Prestador (`prestador`) | Locatário (`locatario`) | Link (`link`) |
|---|---|---|---|
| `ordem_servico` e filhas (`tarefa_os`, `apontamento`, anexos da OS, assinaturas) | OS com `executor_tipo = fornecedor` do seu fornecedor; lê e registra execução | OS dos ativos cedidos a ele, no período da locação; só leitura | Só a OS do link; registra execução |
| `ativo`, `medidor`, `leitura_medidor` | Ativos das OS do seu fornecedor; só leitura, campos reduzidos; registra leitura na execução | Ativos dos contratos vigentes; só leitura | O ativo da OS; registra leitura na execução |
| `solicitacao` | Nenhuma | As que abriu e as dos ativos locados; cria | Nenhuma |
| `almoxarifado`, saldos, `movimentacao_estoque`, `peca_serializada` | Só almoxarifado `em_terceiro` do seu fornecedor; registra consumo em OS e devolução | Nenhuma | Nenhuma no MVP |
| `fornecedor`, `servico_fornecedor`, `contrato_prestacao` | O próprio cadastro, catálogo e contratos; só leitura | Nenhuma | Nenhuma |
| `locatario`, `contrato_locacao`, `item_contrato_locacao`, `vistoria_locacao` | Nenhuma | Os próprios; só leitura; assina a vistoria | Nenhuma |
| `usuario`, `pessoa` | Usuários e técnicos do próprio fornecedor | O próprio usuário | Nenhuma |
| `campo_extra_os` (versão 0.5, **proposta**) | Definições ativas que se aplicam às OS do seu fornecedor; só leitura. Os valores estão em `ordem_servico.campos_extras`, já coberta pela linha da OS | Nenhuma | Só as definições marcadas como do prestador, da OS do link; só leitura |

**Fora do banco.**

| Camada | Mecanismo |
|---|---|
| Rotas | Toda rota declara os escopos que aceita; o padrão é só `interno`. As rotas abertas a escopo externo são uma lista fechada (listar e executar OS, consultar almoxarifado em terceiro, abrir e acompanhar solicitação, consultar contrato e vistoria). Rota não declarada responde 404 a escopo externo, para não confirmar que o recurso existe. |
| Projeção | DTOs próprios para escopo externo: sem custos internos, sem preços de outros fornecedores, sem histórico de outras OS do ativo, sem dados de pessoas da empresa além do nome do responsável pela OS. |
| Funções indisponíveis | Busca global, IA, indicadores, exportação, relatórios, webhooks e chaves de API não existem para escopo externo. **(proposta)** |
| Sincronização do app | O escopo de sincronização é calculado sob o mesmo escopo (seção 10.2); o pull roda dentro de `NaEmpresa` com o escopo, então o RLS restritivo também protege o aparelho. |
| Arquivos e cache | URL assinada só depois de ler a linha sob o escopo. Chave de cache inclui `empresa_id`, o tipo do escopo e o identificador correspondente (fornecedor, locatário ou hash do link). |
| Fim do acesso | Fornecedor inativado ou sem o papel, usuário desativado, ou locatário sem contrato vigente há mais de 90 dias **(proposta)**: a versão de sessão muda e os tokens anteriores deixam de valer. |
| Cobrança | Usuário `prestador` que executa OS conta como técnico (entitlement de quantidade, seção 13.4); `locatario` e link de uso único não contam. O técnico da subcontratada conectada, que tem usuário só na conta dela, não conta no cliente (sumário 2.2). |
| Liberação | A permissão `os.encerrar` nunca é concedida a `prestador` nem a `locatario`; o link de uso único também não encerra. A regra anti-escalada recusa a concessão. |

### 6.7 Endpoint público do link de uso único

O prestador ocasional registra a execução de **uma** OS sem login, por um link enviado pelo gestor. O link é a única porta da API que aceita requisição sem sessão de usuário, por isso tem controles próprios. Entidade `link_execucao`, dona: TER. Decisão no [ADR-008](#adr-008--acesso-externo-por-escopo-restritivo-e-link-de-uso-único).

```mermaid
sequenceDiagram
    participant G as Gestor (web)
    participant API as API /v1
    participant P as Prestador (celular)
    participant PUB as API /v1/publico

    G->>API: POST /v1/links-execucao (ordem_servico_id, contato, validade)
    API-->>G: link https://app.tenaz.com.br/l/#35;token (exibido uma vez)
    API->>P: e-mail com o link (ou o gestor envia por outro canal)
    P->>PUB: POST /v1/publico/links-execucao/ativar (token no corpo)
    PUB-->>P: cookie de sessão curta do link (HttpOnly, 2 h)
    P->>PUB: GET /v1/publico/links-execucao/atual (dados reduzidos da OS)
    P->>PUB: POST /v1/publico/links-execucao/atual/anexos (URL pré-assinada para foto)
    P->>PUB: POST /v1/publico/links-execucao/atual/registrar-execucao (horas, peças, leitura, observação, nome, documento, assinatura simples)
    PUB->>API: casos de uso da OS em NaEmpresa(empresa do link, escopo link)
    Note over PUB: link consumido, OS fica concluida aguardando liberação (os.encerrar)
```

| Controle | Definição **(proposta)** |
|---|---|
| Força do token | 256 bits do gerador criptográfico do sistema, em base64url (43 caracteres). Guardado só como SHA-256 em `link_execucao.token_hash`, com índice único; a entropia alta dispensa hash lento. O token aparece uma única vez, na criação, e nunca em log. |
| Transporte | O token fica no **fragmento** da URL (`/l/#<token>`, depois de `#`), nunca no caminho nem na query, que o navegador não envia ao servidor, à CDN nem no cabeçalho `Referer`. A página lê o fragmento, apaga-o da barra de endereço e o envia uma vez no corpo do POST de ativação. A página tem `Referrer-Policy: no-referrer`, CSP restrita e nenhum script de terceiro. |
| Expiração | Padrão de 72 horas, máximo de 7 dias, configurável pela empresa. A sessão ativada dura 2 horas, renováveis até a expiração do link. |
| Uso único | O link é consumido em `registrar-execucao`, que registra e conclui numa só chamada. Depois disso a página mostra só "Execução registrada na OS-2026-004812". Só um link ativo por OS: gerar outro revoga o anterior. |
| Revogação | Manual pelo gestor a qualquer momento. Automática quando a OS é cancelada, encerrada ou muda de executor ou de fornecedor, e quando o fornecedor é inativado. O estado é lido do banco a cada requisição, sem cache, então a revogação tem efeito imediato. |
| Escopo mínimo | Uma OS, sob o escopo `link` (6.6). Leitura: número, descrição, tarefas, ativo (TAG, nome, local), prazo e contato do gestor. Escrita: apontamento de horas, peças usadas (texto livre ou itens do almoxarifado em terceiro do fornecedor, quando houver), leitura de medidor, fotos, observação, assinatura simples e conclusão. Sem listagem, busca, custos, outros IDs ou navegação. |
| Proteção contra encaminhamento | (1) O link fica preso à primeira sessão que o ativou; outro aparelho recebe "Este link já está em uso em outro aparelho. Peça um novo ao gestor." (2) Código de confirmação opcional, por empresa: 6 dígitos enviados ao contato cadastrado do fornecedor por outro canal, exigido na ativação. (3) Nome e documento de quem executou, registrados com a assinatura simples. (4) Aviso ao gestor na ativação, com horário, aparelho e localização aproximada pelo IP. |
| Rate limit | Por IP e por token, no WAF e na API: 10 ativações inválidas por IP em 10 minutos bloqueiam o IP por 1 hora; 5 erros de código de confirmação revogam o link. Token inexistente, expirado, revogado ou consumido recebem a mesma resposta (404, "Link inválido ou expirado. Peça um novo ao gestor."), em tempo constante. |
| Auditoria | `evento_auditoria` com `origem = link`: criação (quem gerou, destino, validade), envio, cada abertura e ativação (IP, aparelho, horário), cada registro, conclusão, revogação e tentativas inválidas. O `usuario_id` é o de quem gerou o link; o executor fica nos dados do registro. |
| Arquivos | Upload por URL pré-assinada limitada à OS do link, só imagem, tamanho máximo definido, validade de 5 minutos. |
| Efeito na OS | O envio aplica os casos de uso comuns da OS: `iniciar` (se a OS ainda não estiver `em_execucao`, com o horário informado) e `concluir`. A OS fica `concluida` aguardando liberação por quem tem a permissão `os.encerrar`; a liberação pode ser recusada (com `motivo_recusa_liberacao`) e voltar a `em_execucao`, como em qualquer OS. |
| Cobrança | Não conta como técnico (sumário 2.1). |

A rota de gestão é `/v1/links-execucao` (autenticada, escopo interno). As rotas públicas são as do sumário (seção 3.6): `POST /v1/publico/links-execucao/ativar` (token no corpo, devolve a sessão curta), `GET /v1/publico/links-execucao/atual`, `POST /v1/publico/links-execucao/atual/anexos` e `POST /v1/publico/links-execucao/atual/registrar-execucao`. Nenhuma rota leva `{token}` no caminho: depois da ativação, o recurso é sempre `atual`, resolvido pela sessão. As rotas sob `/v1/publico/` são as únicas sem autenticação de usuário (além de `GET /v1/versao`); o prefixo permite aplicar regras próprias de WAF, rate limit e CSP.

### 6.8 Conexão entre contas (SRV)

Risco R5 do sumário: a sincronização entre contas conectadas copia dado além do permitido ou para a conta errada. A subcontratada com conta própria (módulo SRV) e o cliente que também usa a Tenaz são **duas empresas**, isoladas pelo RLS como quaisquer outras. A conexão não abre uma conta para a outra: cada conta guarda a sua cópia, e um processo próprio copia só os campos permitidos. Decisão e alternativas no [ADR-009](#adr-009--conexão-entre-contas). As regras de negócio (quem abre, quem libera, o que acontece no encerramento) estão no sumário, seção 3.5, e em `4_FLUXOS_PROCESSO.md`.

**Onde fica cada dado.**

| Dado | Onde | Isolamento |
|---|---|---|
| `plataforma.conexao_empresa` | Esquema `plataforma`, fora do RLS por empresa, porque a linha pertence a duas empresas (`empresa_cliente_id`, `empresa_prestadora_id`) | O papel `tenaz_app` não tem privilégio na tabela. Cada conta lê as próprias conexões por funções `SECURITY DEFINER` com `search_path` fixo, que devolvem só as colunas da sua ponta (2_ARQUITETURA_BANCO_DADOS.md, seção 2.10.3). Só `tenaz_sinc` lê a tabela diretamente. Escrita só pelas funções de ciclo de vida (convite, aceite, recusa, suspensão, reativação, encerramento), que conferem a ponta e o estado. `campos_compartilhados` guarda os grupos opcionais ligados |
| `vinculo_os` | **Uma linha em cada conta (proposta)**, com `conexao_empresa_id`, `os_cliente_id`, `os_prestador_id`, `lado_origem` e `estado_sincronizacao`. A justificativa final fica no documento 2. Motivo da proposta: cada conta consulta e audita o vínculo sob o seu próprio RLS, sem ler a plataforma em cada tela, e o histórico continua em cada conta depois do encerramento | RLS por empresa comum |
| Ativos compartilhados | Seleção feita pelo cliente na conexão (ativos, localizações ou classes; na Fase 2, também os ativos cobertos pelo `contrato_prestacao`), guardada na conta do cliente. Entidade a definir no documento 2 **(proposta)**. Ativo de OS atribuída ao fornecedor conectado é compartilhado junto com a OS | RLS por empresa do cliente |
| Equipamento do cliente na subcontratada | `ativo` com `propriedade = de_cliente`, `cliente_servico_id`, `conexao_empresa_id` e `ativo_remoto_id` | RLS por empresa da subcontratada |
| OS vinculada | Uma `ordem_servico` em cada conta, cada uma com seu número (`OS-AAAA-NNNNNN` sequencial da própria conta) | RLS por empresa |
| Caixa de saída | `envio_conexao` em cada conta **(proposta)**: `empresa_id`, `mensagem_id`, `conexao_empresa_id`, `vinculo_os_id`, `sequencia`, `tipo`, `estado` (`pendente`, `entregue`, `erro`), `tentativas`, `ultimo_erro` | RLS por empresa |
| Caixa de entrada | `recebimento_conexao` em cada conta **(proposta)**: `empresa_id`, `conexao_empresa_id`, `mensagem_id` (chave única com `empresa_id`), `vinculo_os_id`, `sequencia`, `resultado`, `hash_conteudo` | RLS por empresa |

**Papel de banco `tenaz_sinc`.**

| Propriedade | Definição **(proposta)** |
|---|---|
| Atributos | `NOLOGIN` herdado por um usuário de conexão próprio; sem `BYPASSRLS`, sem superusuário, não é dono de tabela |
| Tabelas | `GRANT` só nas tabelas da lista fechada: `ordem_servico`, `tarefa_os`, `apontamento`, itens de material da OS, `servico_os`, `anexo` da OS, assinaturas, `ativo`, `medidor`, `leitura_medidor`, `vinculo_os`, `envio_conexao`, `recebimento_conexao`, `evento_auditoria` (só `INSERT`), `evento_dominio` (só `INSERT`), tabelas da fila River e `SELECT` em `plataforma.conexao_empresa`. Qualquer outra tabela responde "permissão negada", não zero linhas |
| Colunas | `GRANT SELECT (...)`, `INSERT (...)` e `UPDATE (...)` **por coluna**, exatamente as da lista de campos abaixo. Uma consulta que peça coluna fora da lista falha no banco, mesmo que o código Go esteja errado. Validar na prova de conceito a interação entre privilégio por coluna e predicados de RLS |
| Contexto | Cada transação define `app.empresa_id` (uma só empresa), `app.escopo = 'interno'` e `app.conexao_empresa_id` |
| Política restritiva | Toda tabela da lista tem a política `conexao_sinc AS RESTRICTIVE TO tenaz_sinc`: a linha precisa estar ligada a um `vinculo_os` (ou a um ativo compartilhado) da conexão em processamento |
| Credencial | Segredo próprio, lido só pela tarefa do `tenaz sincronizador` (seção 4.1); rotação automática |

```sql
CREATE FUNCTION app.conexao_atual() RETURNS uuid
  LANGUAGE sql STABLE
  AS $$ SELECT NULLIF(current_setting('app.conexao_empresa_id', true), '')::uuid $$;

-- Só o sincronizador é afetado; o papel tenaz_app não passa por esta política
CREATE POLICY conexao_sinc ON ordem_servico AS RESTRICTIVE TO tenaz_sinc
  USING (EXISTS (
    SELECT 1 FROM vinculo_os v
     WHERE v.empresa_id = ordem_servico.empresa_id
       AND v.conexao_empresa_id = (SELECT app.conexao_atual())
       AND ordem_servico.id IN (v.os_cliente_id, v.os_prestador_id)))
  WITH CHECK (EXISTS (
    SELECT 1 FROM vinculo_os v
     WHERE v.empresa_id = ordem_servico.empresa_id
       AND v.conexao_empresa_id = (SELECT app.conexao_atual())
       AND ordem_servico.id IN (v.os_cliente_id, v.os_prestador_id)));

-- Lista fechada por coluna (trecho)
GRANT SELECT (id, empresa_id, numero, tipo, prioridade, descricao, estado,
              data_reconhecimento, data_inicio, data_conclusao, data_encerramento,
              motivo_recusa_liberacao, descricao_servico)
  ON ordem_servico TO tenaz_sinc;
```

Na criação de uma OS vinculada no destino, o `vinculo_os` é gravado antes da OS, na mesma transação, com a chave estrangeira adiável, para que a política `conexao_sinc` aceite o `INSERT`.

**Fluxo de uma mensagem.**

```mermaid
sequenceDiagram
    participant UC as Caso de uso (conta A, tenaz_app)
    participant W as Worker (conta A, tenaz_app)
    participant SN as Sincronizador (tenaz_sinc)
    participant A as Conta A (origem)
    participant B as Conta B (destino)

    UC->>A: T1: altera OS + evento_dominio + job (outbox de A)
    W->>A: T2: assinante SRV: OS tem vínculo ativo?<br/>grava envio_conexao (mensagem_id, sequência) + job na fila conexao
    SN->>A: T3 (contexto A, conexão X): lê envio_conexao e projeta só as colunas da lista
    SN->>SN: destino = a outra ponta de X em plataforma.conexao_empresa<br/>(nunca vem do job nem da mensagem)
    SN->>B: T4 (contexto B, conexão X): conexão ativa? mensagem_id já recebida?<br/>sequência = última + 1? aplica pelos casos de uso remotos<br/>grava recebimento_conexao, evento_auditoria, vinculo_os.sincronizado
    SN->>A: T5 (contexto A): envio_conexao entregue, vinculo_os.estado_sincronizacao = sincronizado
    Note over SN,B: falha em T4: nova tentativa com backoff, esgotadas, envio em erro,<br/>vinculo_os.erro nas duas contas e aviso ao gestor
```

- **Transações separadas por conta.** T3, T4 e T5 são transações distintas, cada uma com uma única `app.empresa_id`. Nenhuma consulta junta linhas de duas empresas. Se o processo cair entre T4 e T5, a repetição chega a T4, que reconhece o `mensagem_id` na caixa de entrada e responde "já aplicada" sem reaplicar.
- **Outbox nos dois lados.** A origem grava `envio_conexao` na mesma transação lógica do evento (via assinante, com `evento_processado`), e o destino grava `recebimento_conexao` na mesma transação da aplicação. Como a OS trafega nos dois sentidos (cliente atribui e libera; subcontratada executa e abre), as duas tabelas existem nas duas contas.
- **Idempotência.** `mensagem_id` é o `id` do `evento_dominio` de origem (UUIDv7). Chave única `(empresa_id, conexao_empresa_id, mensagem_id)` em `recebimento_conexao`.
- **Sem eco.** Mudança aplicada pelo sincronizador grava `origem = conexao` no evento local; o assinante SRV não gera mensagem de volta para a mesma conexão a partir dela.
- **Auditoria nos dois lados.** A origem registra "enviado à conexão X" com a lista de nomes de campos e o `hash_conteudo`; o destino registra "recebido da conexão X", com conta de origem, `mensagem_id` e o mesmo hash. Os valores não vão para log, só para a auditoria de cada conta.

**Lista fechada de campos.** Versionada no repositório (`backend/internal/modulos/servicos/campos_conexao.yaml`, **proposta**), é a mesma usada para gerar os `GRANT` por coluna, os DTOs de mensagem e o teste da seção 9.3. Mudança na lista exige revisão de segurança no PR.

| Grupo | Direção | Campos | Padrão |
|---|---|---|---|
| Cadastro do ativo | cliente → subcontratada | TAG, nome, tipo, classe (nome), fabricante, modelo, número de série, ano, unidade e localização (texto), criticidade, medidores (tipo, unidade, última leitura) | Sempre |
| Documentos técnicos do ativo | cliente → subcontratada | Manuais e documentos marcados pelo cliente como compartilháveis | Opcional, desligado |
| OS atribuída | cliente → subcontratada | Número da OS no cliente, tipo, prioridade, descrição, tarefas, prazo, nome e telefone de trabalho do contato na unidade | Sempre |
| Liberação e cancelamento | cliente → subcontratada | Encerramento (`data_encerramento`), recusa com `motivo_recusa_liberacao`, cancelamento com motivo | Sempre |
| OS aberta pela subcontratada | subcontratada → cliente | Ativo (`ativo_remoto_id`), descrição, tipo e prioridade sugeridos, fotos | Sempre |
| Execução | subcontratada → cliente | Estado, `data_reconhecimento`, `data_inicio`, pausas com motivo, `data_conclusao`, `descricao_servico`, itens `servico_os` (descrição e quantidade), apontamentos (nome do técnico, início, fim, horas), componentes (`tipo_material`, item em texto, quantidade, `numero_serie_retirado`, `numero_serie_instalado`), leituras de medidor, respostas de checklist, fotos, observações, assinatura do cliente (nome, função, imagem, hash) | Sempre |
| Preços | subcontratada → cliente | Preço de referência dos itens `servico_os` e dos componentes | Opcional, desligado; só a subcontratada liga |

Nunca trafegam: CPF, contatos pessoais, documentos e capacitações de técnicos (a Tenaz não guarda dado de saúde, sumário 3.6); custos internos e preços de outros fornecedores do cliente; OS do ativo não vinculadas; dados de outros clientes da subcontratada; usuários e permissões. Cada lado controla os grupos opcionais dos dados que **envia**.

**Dono de cada dado e conflitos.**

| Dado | Dono | Na outra conta |
|---|---|---|
| Cadastro do ativo | Cliente | Somente leitura na subcontratada. Ela pode acrescentar anotações locais e QR code próprio, que não voltam |
| Abertura, descrição, prioridade e prazo da OS atribuída | Cliente | Somente leitura na subcontratada; mudança no cliente sobrescreve |
| OS aberta pela subcontratada | Subcontratada até a triagem do cliente; depois, cliente | Nasce `aberta` no cliente com `origem = prestador` |
| Distribuição ao técnico, especialidade, equipe | Subcontratada | Ao cliente vai só o nome do técnico executor |
| Execução (de `em_execucao` a `concluida`, apontamentos, componentes, fotos, checklist, assinatura) | Subcontratada | Somente leitura no cliente |
| Liberação (`os.encerrar`), recusa e cancelamento | Cliente | Na subcontratada, a OS vinculada encerra, volta a `em_execucao` com o motivo ou recebe a marca "cancelada no cliente" |
| Leitura de medidor registrada na execução | Subcontratada registra; cliente valida | Entra no cliente com origem `os` e passa pelas mesmas regras de validação (`a_revisar` quando fora do padrão) |

- **Ordem.** Sequência por vínculo e por sentido, atribuída na origem. O destino só aplica `sequencia = última aplicada + 1`; menor ou igual é repetição (responde "já aplicada"); maior fica aguardando, com nova tentativa. Lacuna por mais de 15 minutos vira `erro` e pede reenvio à origem. **(proposta)**
- **Campo de outro dono.** Mensagem que tenta alterar campo que não pertence à conta de origem é recusada inteira, registrada como `erro` de conflito e gera alerta de segurança (não deveria acontecer com o código correto).
- **Transição impossível.** Mensagem que exige transição inválida no destino (ex.: conclusão chegando a uma OS já cancelada no cliente) não é descartada: os apontamentos, fotos e componentes são preservados e ligados à OS, e o gestor da conta de destino recebe uma pendência, com a mesma regra do app (seção 10.5).
- **A liberação continua do cliente.** Não existe mensagem de encerramento no sentido subcontratada → cliente. Na conta da subcontratada, a OS vinculada não pode ser encerrada localmente enquanto a conexão estiver `ativa`; sem conexão (cliente fora da Tenaz), quem tem `os.encerrar` na subcontratada encerra.

**Ciclo de vida da conexão.**

| Ação | Quem | Efeito técnico |
|---|---|---|
| `convidar` | Cliente, a partir de um `fornecedor` com papel `subcontratada` ou `assistencia_tecnica` | Conexão `convidada` com código de convite (256 bits, guardado como hash, validade de 14 dias **proposta**) enviado por e-mail ao responsável. A subcontratada sem conta é levada a criar a conta Tenaz Serviços antes de aceitar |
| `aceitar` | Subcontratada, usuário com a permissão `srv.gerenciar_conexao` (**proposta**, padrão `administrador` e `gestor`; existe também na conta do cliente para convidar, suspender e encerrar) | Confere que o CNPJ da conta que aceita é o do `fornecedor` convidado (convite encaminhado a outra empresa não serve); liga ou cria o `cliente_servico`; conexão `ativa`; `conexao_empresa.ativada` nas duas contas; cópia inicial dos ativos compartilhados |
| `recusar` | Subcontratada | Conexão `encerrada` sem dados copiados |
| `suspender` / `reativar` | Qualquer lado | Em `suspensa`, as mensagens se acumulam na caixa de saída (retidas por até 30 dias **proposta**); ao reativar, são entregues em ordem; depois de 30 dias, o vínculo é ressincronizado por estado completo |
| `encerrar` | Qualquer lado | Mensagens geradas antes de `data_encerramento` ainda são entregues por até 24 horas **(proposta)**; depois, nada mais trafega. Cada conta mantém o histórico e as cópias que já recebeu (retenção em 9_SEGURANCA_CONFIABILIDADE.md). OS vinculadas em aberto continuam em cada conta: no cliente, voltam ao fluxo comum de TER (perfil `prestador` ou link); na subcontratada, terminam localmente com o relatório PDF. `conexao_empresa.encerrada` nas duas contas |

**Reprocessamento.**

1. Nova tentativa automática com backoff exponencial (fila River `conexao`), até 24 horas. Esgotadas, `envio_conexao` e `vinculo_os` ficam em `erro` e o evento `vinculo_os.erro` é publicado nas duas contas.
2. Tela "Sincronização da conexão" em cada conta (detalhe em `7_PADROES_DESIGN.md`): vínculos em erro, motivo legível e ação **reprocessar**, que reenfileira a partir da última sequência aplicada.
3. **Ressincronizar vínculo:** a origem envia uma mensagem `estado_completo` com o valor atual de todos os campos que ela possui; o destino sobrescreve só esses campos. Usada para lacunas, suspensão longa e divergência.
4. **Reconciliação diária:** o sincronizador calcula, em cada conta e em transações separadas, o hash dos campos compartilhados de cada vínculo ativo; divergência dispara a ressincronização e conta na métrica de divergência.

**Observabilidade.**

| Sinal | Definição **(proposta)** |
|---|---|
| Atraso | Histograma do tempo entre o commit na origem (`evento_dominio.ocorrido_em`) e o commit em T4. **SLO: p95 abaixo de 60 segundos** em 30 dias (seção 14.5) |
| Fila | Mensagens pendentes por estado; idade da pendente mais antiga; tentativas |
| Erros | Por motivo: conexão não ativa, lacuna de sequência, transição impossível, campo de outro dono, campo fora da lista, ponta inválida |
| Reconciliação | Vínculos verificados e divergentes por dia |
| Logs e traces | `conexao_empresa_id`, `empresa_origem_id`, `empresa_destino_id`, `vinculo_os_id`, `mensagem_id`, `sequencia`. O contexto de trace viaja em `envio_conexao`, ligando T1 a T5 num mesmo trace. Nunca o valor dos campos |
| Visível ao usuário | Cada OS vinculada mostra o estado de sincronização e o horário da última entrega nas duas contas |

### 6.9 Credenciais de autenticação (versão 0.6)

Achado ISO-E01-02 da revisão de isolamento da etapa E01: as tabelas de autenticação do esquema `plataforma` (nomes do sumário 0.6, seção 3.2) não têm RLS, porque são lidas antes de a empresa ser conhecida, e estavam ao alcance do `tenaz_app`. Com o contexto da empresa A, o `tenaz_app` lia e-mail, `senha_hash` e `mfa_segredo_cifrado` de todas as identidades, lia `acesso_empresa` inteiro (que revela quais clientes cada subcontratada atende), lia sessões de outra empresa com IP e user agent, criava sessão web válida para um usuário de outra empresa com o cookie que escolhesse, criava código OAuth e trocava a senha de qualquer identidade. As três regras do desenho inicial (só hash, empresa fixada na criação, chave estrangeira composta) impediam mover uma credencial de empresa, mas não lê-la nem forjá-la. Decisão no [ADR-011](#adr-011--papel-de-banco-próprio-para-as-credenciais-de-autenticação). O DDL e os `GRANT` são do documento 2 (seções 2.2, 2.11 e 5.1).

**Papel `tenaz_autenticacao`.** `NOLOGIN`, `NOBYPASSRLS`, não é dono de objeto, não é membro de `tenaz_dono`, `tenaz_app` nem `tenaz_job`, e nenhum deles é membro dele. A conexão usa o usuário `tenaz_autenticacao_conexao` (`LOGIN`, herda o papel), no mesmo par que `tenaz_sinc` e `tenaz_sinc_conexao`. Fica fora dos privilégios padrão do esquema `app`.

**Dono das credenciais: `tenaz_credencial` (emenda ISO-E01-08).** Numa tabela sem RLS, o dono lê tudo, e quem herda do dono também. Com as sete tabelas de dono `tenaz_dono`, o `tenaz_manutencao` (`LOGIN`, membro de `tenaz_dono` com herança, para partições e `VACUUM`) lia `senha_hash` e o segredo de MFA de todas as identidades. Por isso as sete tabelas e as duas funções `SECURITY DEFINER` que as tocam (`plataforma.tg_acesso_empresa()` e `plataforma.purgar_credenciais_vencidas()`) pertencem a `tenaz_credencial`: `NOLOGIN`, `NOBYPASSRLS`, cujo único membro possível é o papel que aplica migração: a `00004` concede a pertença a quem a executa, se ainda não a tiver. Em homologação e produção é o `tenaz_migrador`, que vira membro; em desenvolvimento e nos testes o `tenaz_migrador` é superusuário (imagem oficial do PostgreSQL), não recebe a pertença e o papel fica sem membros. Não é membro de `tenaz_dono` nem tem como membro `tenaz_dono`, `tenaz_manutencao` ou qualquer papel de runtime. As migrações criam esses objetos com `SET ROLE tenaz_credencial` e, no `Up`, voltam a `tenaz_dono` antes de terminar (documento 2, seção 14.2). Os privilégios cruzados entre os dois donos são só os necessários a criar os objetos, às chaves estrangeiras e ao trigger, todos concedidos na `00004`:

| Concedente → beneficiário | Privilégio | Para quê |
|---|---|---|
| `tenaz_credencial` → `tenaz_dono` | `REFERENCES (id)` em `plataforma.identidade` | Chave estrangeira de `app.usuario.identidade_id`; não dá leitura |
| `tenaz_credencial` → `tenaz_dono` | `EXECUTE` em `plataforma.tg_acesso_empresa()` | Criar o trigger em `app.usuario`; função de trigger não se chama diretamente |
| `tenaz_dono` → `tenaz_credencial` | `REFERENCES (empresa_id, id, identidade_id)` em `app.usuario` | Chaves estrangeiras compostas das tabelas de credencial; não dá leitura |
| Esquemas de `tenaz_dono` → `tenaz_credencial` | `USAGE` e `CREATE` em `plataforma`; `USAGE` em `app` | Criar as tabelas e funções em `plataforma` e resolver `app.usuario` e `app.uuid_v7()`; nenhum privilégio em tabela de `app` além do `REFERENCES` acima |
| Esquema `public` → `tenaz_credencial` | `USAGE` | Tipo `citext` de `identidade.email` |

Nenhum outro privilégio entre os dois donos. A regra vale para toda tabela futura de `plataforma` sem RLS que não seja catálogo global de leitura: o dono não é `tenaz_dono` (pendência E01-14 do documento 2, a primeira é `conexao_empresa`; seção 16).

| Objeto | `tenaz_autenticacao` | `tenaz_app` | `tenaz_job` |
|---|---|---|---|
| `plataforma.identidade` | `SELECT`, `INSERT`, `UPDATE (email, senha_hash, mfa_segredo_cifrado, mfa_ativado_em, mfa_ultimo_passo, email_verificado_em, falhas_login, bloqueada_ate, bloqueada_em, ultimo_login_em, atualizado_em)` | — | — |
| `plataforma.acesso_empresa` | `SELECT` (a escrita continua só do trigger) | — | — |
| `plataforma.sessao` | `SELECT`, `INSERT`, `UPDATE (ultimo_uso_em, expira_em, revogada_em)` | — | — |
| `plataforma.token_renovacao` | `SELECT`, `INSERT`, `UPDATE (usado_em, revogado_em)` | — | — |
| `plataforma.codigo_autorizacao` | `SELECT`, `INSERT`, `UPDATE (usado_em, sessao_id)` | — | — |
| `plataforma.desafio_mfa` | `SELECT`, `INSERT`, `UPDATE (tentativas, concluido_em)` | — | — |
| `plataforma.token_email` | `SELECT`, `INSERT`, `UPDATE (usado_em, revogado_em)` | — | — |
| `app.evento_auditoria` | `INSERT` só nas colunas `empresa_id`, `ator_tipo`, `usuario_id`, `acao`, `entidade`, `entidade_id`, `depois`, `campos_alterados`, `origem`, `ip`, `dispositivo_id`, `user_agent`, `request_id`; sem `SELECT` | inalterado | inalterado |
| Função `plataforma.purgar_credenciais_vencidas()` | — | — | `EXECUTE` |
| Esquemas e funções de apoio | `USAGE` em `plataforma`, `app` e `public` (o tipo `citext` de `email` mora em `public`); `EXECUTE` em `app.uuid_v7()` e `app.empresa_atual()` | `USAGE` em `plataforma` mantido, para as funções seguras da conexão e os catálogos globais | `USAGE` em `plataforma`, para a purga |

Fora o dono `tenaz_credencial`, nenhum papel tem `DELETE` nem `TRUNCATE` nas sete tabelas; a purga apaga com o privilégio do dono, dentro da função. `tenaz_leitura`, `tenaz_sinc`, `tenaz_dono` e `tenaz_manutencao` não têm privilégio nelas, salvo o `REFERENCES (id)` do `tenaz_dono` acima. Qualquer acesso fora da tabela acima responde "permissão negada", não zero linhas.

**Quem usa o pool.** Só o pacote `internal/plataforma/autenticacao`, pelo repositório `autenticacao/internal/repo`, com o código sqlc gerado ali mesmo (regra 12 da seção 5.3). O middleware HTTP de sessão não toca o pool: chama o serviço. O pool é aberto só no modo `tenaz api` (seção 4.1), com poucas conexões. O repositório:

- recusa ser chamado dentro de uma transação de `NaEmpresa` (erro de programação) e nunca chama `NaEmpresa` nem `app.conferir_escopo()`, que lê `app.usuario` e falharia por falta de privilégio;
- abre dois tipos de transação: **sem empresa**, para achar credenciais por hash, e-mail ou id; e **na trilha de uma empresa**, que define só `app.empresa_id`, `app.origem`, `app.ip`, `app.user_agent` e `app.request_id` com `set_config(..., true)`, para gravar o evento de auditoria junto da mudança de credencial;
- só tem consultas com filtro por hash, `id`, `identidade_id`, `(empresa_id, usuario_id)` ou e-mail exato. Um teste falha se alguma consulta do repositório ler uma das sete tabelas sem um desses filtros, para que nem o pacote de autenticação enumere credenciais.

**Resolução da sessão.** Três passos em dois pools, descritos na seção 6.3. A empresa da transação de negócio vem sempre da linha da sessão, lida pelo `tenaz_autenticacao`.

**Auditoria de autenticação.** Login, falha de login, logout, troca de senha e ativação ou desativação de MFA vão para `app.evento_auditoria` de cada empresa em que a identidade tem acesso (documento 2, seção 9.1), gravados pelo próprio `tenaz_autenticacao`, no molde do `tenaz_sinc`: `INSERT` por coluna e sem leitura da trilha; o papel entra na política permissiva de inserção (empresa da transação) e ganha uma política **restritiva** própria, que só aceita `ator_tipo = 'usuario'`, `entidade = 'identidade'`, `acao` em `login`, `login_falho`, `logout` e `alteracao`, e um `usuario_id` que seja, em `acesso_empresa`, o usuário daquela identidade naquela empresa. Com isso:

- a sessão, o zerar das falhas de login e o evento de login na empresa da sessão são uma transação só, no pool da autenticação;
- o evento nas demais empresas da identidade é uma transação por empresa, que não cita as outras (a empresa B não fica sabendo que a pessoa também entra na empresa A);
- o `tenaz_app` não grava evento de autenticação, e o `tenaz_autenticacao` não grava nenhum outro tipo de evento.

**`acesso_empresa` e o trigger.** `plataforma.tg_acesso_empresa()` continua `SECURITY DEFINER`, agora de dono `tenaz_credencial` (roda com o privilégio do dono das credenciais, que o `tenaz_dono` não tem), disparado pelo `INSERT` ou `UPDATE` em `app.usuario` feito pelo `tenaz_app` na transação do ADM; ele nunca dependeu de privilégio do `tenaz_app` em `acesso_empresa`, então continua funcionando sem mudança. Ganha uma responsabilidade: quando o acesso fica desabilitado (usuário `inativo` ou excluído), revoga, na mesma transação, as sessões abertas daquele usuário naquela empresa e os refresh tokens das famílias revogadas. A inativação passa a revogar as credenciais de forma atômica sem que o ADM toque tabela de credencial; `app.conferir_escopo()` continua recusando o usuário inativo na requisição seguinte, como segunda barreira.

**Purga.** `plataforma.purgar_credenciais_vencidas()`, `SECURITY DEFINER`, dono `tenaz_credencial`, `search_path` fixo, sem parâmetro, apaga das cinco tabelas de credencial as linhas com `expira_em` há mais de 30 dias **(proposta**, o prazo da investigação de incidente) e devolve a contagem. Só o `tenaz_job` executa, por um job diário. Não recebe nenhum dado da requisição, por isso é o caso certo para uma função estreita.

**Operações que atravessam os dois papéis.** São sequências de transações, nunca uma, na ordem que torna a repetição segura: a credencial de uso único é consumida por último.

| Operação | Sequência |
|---|---|
| Login web | (1) autenticação: identidade por e-mail, senha, bloqueio, desafio de MFA, lista de acessos; (2) aplicação: `NaEmpresa(empresa escolhida, interno, sem usuário)` lê usuário, perfis e nome da empresa; (3) autenticação, na trilha da empresa: cria a sessão, zera as falhas, grava o evento `login`; (4) autenticação: evento `login` em cada outra empresa da identidade, sem derrubar o login se falhar |
| Troca de empresa | (1) autenticação: acesso à empresa pedida; (2) aplicação: carrega o ator na empresa nova; (3) autenticação, na trilha da empresa nova: revoga a sessão atual e a família, cria a nova, grava o evento |
| Convite de usuário (ADM) | (1) aplicação: o ADM cria o usuário `convidado` e confirma; (2) autenticação: `EmitirConvite` grava `token_email` com a empresa do ator autenticado (a chave estrangeira composta garante que o usuário é daquela empresa) e envia o e-mail. Falha em (2) deixa o usuário convidado sem convite válido; a resposta diz isso e o administrador reenvia |
| Aceite do convite | (1) autenticação: confere o token sem consumir e cria ou acha a identidade; (2) aplicação: liga `usuario.identidade_id` e ativa o usuário (o trigger cria o `acesso_empresa`); (3) autenticação: consome o token e grava o evento. Cada passo é idempotente |
| Encerrar sessões a pedido do administrador | Autenticação, com a empresa do ator autenticado e o usuário alvo; evento na trilha da empresa |
| Inativação de usuário | Só a aplicação; a revogação vem do trigger |

**Configuração.** `TENAZ_BANCO_URL_AUTENTICACAO` (obrigatória no modo `api` quando houver banco) e `TENAZ_BANCO_AUTENTICACAO_MAX_CONN` (padrão 5, **proposta**). Em homologação e produção, a credencial de `tenaz_autenticacao_conexao` fica num segredo próprio do Secrets Manager, lido só pela tarefa da API, com rotação automática; o processo recusa subir se a URL for igual a `TENAZ_BANCO_URL`. Em desenvolvimento e nos testes, a URL pode ser a do migrador, com `SET ROLE tenaz_autenticacao` na conexão, como já acontece com `tenaz_app`.

---

## 7. Estrutura do repositório

Monorepo. As pastas `docs/`, `marca/` e `ferramentas/` já existem e continuam onde estão. **(proposta)**

```
tenaz/
├── backend/
│   ├── cmd/tenaz/                # main.go: subcomandos api, worker, sincronizador, ingestor, migrate
│   ├── internal/
│   │   ├── servidor/             # composição única da API: plataforma/api + <modulo>/rotas;
│   │   │                         # usada por tenaz api, tenaz openapi e testes de rota (seção 2.5)
│   │   ├── plataforma/           # existentes: api, auditoria, autenticacao, autorizacao, config,
│   │   │                         # db, email, erros, migracao, obs, paginacao, resiliencia, versao;
│   │   │                         # pais e textos (ADR-013, versão 0.9; textos/catalogos/: pt-BR.json
│   │   │                         # e es-419.json embutidos); a seguir: dinheiro (ADR-013, Parte F),
│   │   │                         # empresa, entitlements, eventos, fila, arquivos
│   │   └── modulos/              # cada um no layout da seção 5.2 (raiz, rotas/, jobs/, internal/)
│   │                             # adm, atv, med, sol, ordemservico, pre, chk, est,
│   │                             # terceiros, servicos, locacao, fro, con, ind, integracoes, ia, mob
│   │                             # + projetos, contratos, configuracao (versão 0.5;
│   │                             # configuracao/pacotes/: pacotes de segmento em YAML)
│   ├── db/
│   │   ├── migracoes/            # goose: NNNNN_descricao.sql (cinco dígitos)
│   │   └── sementes/             # dados de desenvolvimento e de teste (empresas fictícias)
│   ├── api/openapi.yaml          # contrato gerado e versionado
│   ├── testes/
│   │   ├── integracao/           # testcontainers: fluxos HTTP de ponta a ponta
│   │   ├── migracoes/            # versão 0.9: conversão, reversão e recusa das migrações
│   │   │                         # de dados (migracao.AplicarAte, ReverterAte)
│   │   ├── vazamento/            # isolamento entre empresas, entre escopos externos
│   │   │                         # e entre contas conectadas (três contas)
│   │   └── casos/                # casos compartilhados com o TypeScript (JSON);
│   │                             # identificacao.json (versão 0.9) e dinheiro.json, este
│   │                             # gerado por gerar_dinheiro.py (tarefa conferir-casos)
│   ├── sqlc.yaml
│   ├── .golangci.yml
│   └── Dockerfile
├── web/
│   ├── src/
│   │   ├── app/                  # rotas, layout, provedores
│   │   ├── modulos/              # uma pasta por módulo: os/, atv/, pre/, ter/, srv/, loc/...
│   │   ├── portais/              # layouts restritos do prestador e do locatário
│   │   ├── link/                 # entrada separada do link de uso único (bundle próprio)
│   │   ├── componentes/          # componentes de interface reutilizáveis
│   │   └── lib/                  # cliente da API, formatação, autenticação
│   ├── e2e/                      # Playwright
│   └── vite.config.ts
├── app/
│   ├── app/                      # telas (expo-router)
│   ├── src/
│   │   ├── sync/                 # outbox de comandos, pull incremental, conflitos
│   │   ├── banco/                # esquema e migrações do SQLite
│   │   ├── contas/               # sessões por conta, seletor de conta, um SQLite por conta
│   │   ├── modulos/              # os/, chk/, med/, sol/, est/ (em terceiro), srv/ (componentes,
│   │   │                         # assinatura do cliente), loc/ (vistoria)
│   │   └── componentes/
│   ├── e2e/                      # fluxos Maestro
│   └── app.config.ts
├── pacotes/
│   ├── api-tipos/                # tipos TS gerados do OpenAPI (não editar)
│   ├── dominio/                  # regras compartilhadas web + app: estados da OS,
│   │                             # validação de leitura, formatação de TAG e números
│   ├── tokens-ui/                # tokens de cor, tipografia e espaçamento da marca
│   └── i18n/                     # textos da interface: pt-BR (fonte) e, versão 0.8
│                                 # (ADR-013, aceito), es-419; os do servidor ficam em
│                                 # backend/internal/plataforma/textos/
├── infra/
│   ├── tofu/
│   │   ├── modulos/              # rede, ecs, rds, s3, kms, backup, observabilidade
│   │   └── ambientes/            # dev, hml, prd
│   ├── docker/                   # compose de desenvolvimento local
│   └── observabilidade/          # painéis, SLOs e alertas como código
├── docs/
├── marca/
├── ferramentas/
├── .github/workflows/
├── pnpm-workspace.yaml
└── README.md
```

| Decisão | Escolha **(proposta)** | Motivo |
|---|---|---|
| Gerenciador JS | pnpm workspaces | Monorepo com pacotes compartilhados entre web e app. Se o Metro do Expo exigir, usar `node-linker=hoisted`. |
| Orquestração de tarefas | Scripts do pnpm; Turborepo só se o tempo de CI justificar | Menos ferramentas no início |
| Regras duplicadas Go e TS | `pacotes/dominio` com casos de teste em JSON em `backend/testes/casos/`, executados pelas duas suítes | A máquina de estados da OS e a validação de leitura existem nos dois lados (o app valida offline). Os mesmos casos impedem divergência. |
| Desenvolvimento local | `docker compose` com PostgreSQL, Valkey, Mailpit e SeaweedFS (S3), todos por versão e digest (SeaweedFS no lugar do MinIO desde a versão 0.7, ADR-012; o compose é do `tenaz-plataforma-sre`). `ferramentas/tarefas.sh infra` espera os quatro saudáveis e cria o bucket `tenaz-auditoria` com bloqueio de objeto; `armazenamento-provar` qualifica o endurecimento em 19 verificações. Backend e web rodam no host. **Versão 0.8:** `ferramentas/tarefas.sh api` já sobe a API com o S3 do selo apontado ao bucket `tenaz-auditoria`, pela identidade `tenaz_auditoria` (só de desenvolvimento, `TENAZ_S3_*`, com `TENAZ_S3_CAMINHO_EXPLICITO=true`, exigido pelo SeaweedFS); `ferramentas/tarefas.sh selar-auditoria` roda o job do selo sobre o banco e o bucket locais. O job de desenvolvimento usa `TENAZ_BANCO_URL` com `SET ROLE tenaz_job`; fora do desenvolvimento, `TENAZ_BANCO_URL_JOB` é obrigatória e diferente das outras URLs | Ciclo rápido |
| Validação de campos extras Go e TS (versão 0.5) | Mesma função no domínio da OS e em `pacotes/dominio`, com casos JSON em `backend/testes/casos/` | O app valida os campos extras offline antes de enfileirar o comando; o servidor revalida. Os mesmos casos impedem que o app aceite o que o servidor recusa (seção 5.6) |
| Pacotes de segmento (versão 0.5) | Arquivos YAML em `backend/internal/modulos/configuracao/pacotes/`, validados no CI e publicados por `tenaz migrate` | Revisão em PR e teste antes de chegar a qualquer conta (seção 5.6) |

---

## 8. Convenções de código

### 8.1 Gerais

| Tema | Convenção |
|---|---|
| Idioma do domínio | Português, sem acento, com os nomes canônicos do sumário: `OrdemServico`, `PlanoManutencao`, `Concluir`, `RegistrarLeitura` |
| Idioma técnico | Inglês idiomático da linguagem para termos técnicos: `ctx`, `err`, `Handler`, `New` |
| Banco | `snake_case`, português, sem acento (`ordem_servico.data_abertura`) |
| Rotas | `/v1/<recurso-no-plural>` em kebab-case (`/v1/ordens-servico`, `/v1/clientes-servico`, `/v1/conexoes-empresa`). `GET /v1/versao` é exceção aceita (sumário 3.6) |
| JSON da API | `snake_case`, igual ao banco |
| IDs | UUIDv7 em todas as entidades. Ordenáveis por tempo, não enumeráveis, e podem ser gerados no app offline. Números legíveis (`OS-2026-004812`) são campos à parte. |
| Datas | UTC no banco e na API (RFC 3339); fuso da unidade só na exibição |
| Dinheiro | Inteiro em centavos mais o código da moeda; nunca `float`. **Versão 0.8 (ADR-013, aceito, Parte A, item 5):** substituído: decimal (`numeric` no banco, texto decimal na API, tipo decimal da plataforma no Go), sempre com a moeda ISO 4217 ao lado e arredondado às casas da moeda (o guarani não tem centavos); a regra de centavos divergia dos documentos 2 e 3. **Versão 0.9 (ADR-013, Parte F):** duas espécies, valor (`numeric(18,2)`, na escala da moeda) e preço unitário (`numeric(18,6)`); texto decimal canônico na API; no Go, `plataforma/dinheiro` (`int64` na unidade menor da moeda, `math/big` no cálculo), no TypeScript, `bigint` em `pacotes/dominio`; arredondamento meio para longe do zero, um só, no fim do cálculo |
| Datas e fusos (versão 0.8, ADR-013, aceito) | Fuso sempre por nome IANA, validado pelo país da empresa; `time/tzdata` embutido no binário |
| Textos para pessoas (versão 0.8, ADR-013, aceito) | Nunca escritos no código: vêm do catálogo do idioma (`pacotes/i18n` na interface, `plataforma/textos` no servidor), com chave estável, `{{variavel}}` e plural CLDR |
| Commits | Conventional Commits em português: `feat(os): pausa com motivo obrigatório` |
| Comentários | Explicam o porquê, não o quê |

### 8.2 Go

- `gofmt` e `goimports` obrigatórios; `golangci-lint` com `errcheck`, `staticcheck`, `gosec`, `depguard`, `sqlclosecheck`, `bodyclose`, `contextcheck`.
- `context.Context` é sempre o primeiro parâmetro de funções com IO.
- Sem estado global mutável. Dependências entram pelo construtor.
- Erros com contexto: `fmt.Errorf("baixar consumo da OS %s: %w", id, err)`.
- Erros de domínio são valores `*erros.Erro` criados com `erros.Novo(codigo, titulo, status)`, com código estável em `snake_case` e português, da tabela da seção 3.2 de `3_ESPECIFICACAO_APIS.md` (código novo entra primeiro lá). Os gerais ficam em `plataforma/erros` (ex.: `erros.TransicaoInvalida` → `transicao_invalida`, 409); os de um módulo, no pacote raiz dele (ex.: `adm.ErrUltimoAdministrador` → `ultimo_administrador`). O caso de uso preenche a ocorrência com `.Com(detalhe)`, `.ComCampo(campo, codigo, mensagem)` e `.ComDado(membro, valor)`; a camada HTTP traduz para RFC 9457 (`application/problem+json`).
- **Erros com texto do catálogo (versão 0.9, ADR-013, Parte B, item 5; conferido no código).** O título vem do catálogo pela chave `erros.ChaveTitulo(codigo)` (`erro.<codigo>.titulo`); o detalhe, por `.ComChave(chave, params)`; o erro por campo, por `.ComCampoChave(campo, codigo, chave, params, valoresAceitos...)`, que leva os valores aceitos de uma enumeração em `erros[].valores_aceitos`. `erros.ParaProblemaIdioma` renderiza no idioma negociado, e a resposta leva `Content-Language` e `Vary: Accept-Language`. `Error()` registra o código e a chave, nunca o texto traduzido; `.Tecnico(...)` guarda texto só para o log, fora da resposta. `erros.Todos()` lista os códigos para a tarefa `textos` conferir que todo código tem título nos dois idiomas. Os grupos de chaves do catálogo do servidor são `erro.<codigo>.titulo`, `detalhe.*`, `campo.*`, `validacao.*`, `identificacao.*`, `pais.*`, `tela.<codigo>` (telas da trilha) e `email.*`. **Feito na E01:** o `tenaz-backend-go` passou ao catálogo todo o texto de erro que estava livre (118 pontos na autenticação, no middleware, na idempotência, na auditoria e no ADM); nenhum literal resta no código de produção.
- **Trava do texto livre (versão 0.9, ADR-013, Parte B, item 5; conferida em `erros/trava_test.go` e `ferramentas/tarefas.sh`).** Texto para pessoa não se escreve no código Go. Erro com texto sai por `Erro.ComChave(chave, params)` e `Erro.ComCampoChave(campo, codigo, chave, params, valoresAceitos...)`, que buscam no catálogo de `plataforma/textos` e renderizam no idioma da requisição. Texto só para o log sai por `Erro.Tecnico(...)`: aparece em `Error()`, nunca na resposta. Valores que são texto do catálogo (estado, dia da semana) vão nos parâmetros como `textos.Chave("estado.x")`. `.Com(` e `.ComCampo(` ficam **proibidos no código de produção** (`internal/` e `cmd/`); os métodos continuam só para fontes de teste. A tarefa `textos` roda `TestSemTextoLivreNoCodigoDeProducao` (sobre `go/ast`, biblioteca padrão), que recusa `.Com(`, `.ComCampo(` e `erros.Novo(` com título que não confere com o catálogo `pt-BR` ou sem texto nos dois idiomas; `TestTravaAcusaTextoLivre` prova que ela acusa, sobre fontes falsas (oito formas proibidas acusadas, cinco permitidas aceitas). Código de erro novo entra por `erros.Novo` com título literal, conferido contra `erro.<codigo>.titulo`. **Limite declarado:** a trava não cobre texto montado fora de `erros` (por exemplo, `fmt.Errorf` usado como erro de domínio que chegue à resposta); hoje não há nenhum caso, e erro de domínio continua sendo `*erros.Erro` (item acima).
- **Idioma da requisição (versão 0.9).** `textos.Preferencias.Resolver()` aplica a ordem da Parte B, item 3, do ADR-013: o `Accept-Language` é lido ao envolver a requisição, e o idioma do usuário e o padrão da empresa entram com o `Ator`, no middleware da credencial. `textos.T(idioma, chave, params)` cai para `pt-BR` quando falta a tradução e, sem texto em nenhum idioma, devolve a própria chave.
- `panic` só para erro de programação na inicialização. Um middleware recupera panics em requisições, registra o erro e devolve 500.
- Sem SQL concatenado. Ordenação dinâmica só por lista fechada de colunas.
- Toda chamada externa (LLM, telemática, ERP) tem timeout, retry com backoff e disjuntor (circuit breaker).
- Testes de integração (versão 0.7): cada pacote de teste de integração sobe o seu PostgreSQL descartável (`internal/plataforma/db/dbteste`). A subida é serializada entre os processos do `go test ./...` por um arquivo de trava no diretório temporário, e a falha transitória de detecção do Docker é repetida com recuo.

Formato de erro (o dono é a seção 3.1 de `3_ESPECIFICACAO_APIS.md`; o tipo Go é `erros.Problema`). `type` é `https://api.tenaz.com.br/erros/<codigo>` (o domínio depende da P5); `codigo` é o que o cliente trata, e `title` e `detail` são texto para pessoas; `id_rastreio` repete o `X-Request-Id`; `erros[]` lista os erros por campo, com `campo` em notação de ponto, `codigo` do vocabulário curto (`obrigatorio`, `valor_invalido`, `formato_invalido`...) e `mensagem`; membros adicionais do código (ex.: `estado_atual` em `transicao_invalida`) saem no nível de cima. Exemplo real de E01, reativar uma unidade que já está ativa:

```json
{
  "type": "https://api.tenaz.com.br/erros/transicao_invalida",
  "title": "Ação não permitida no estado atual",
  "status": 409,
  "detail": "a unidade já está ativa",
  "instance": "/v1/unidades/0193f2a1-7c3e-7b10-9a2e-5d1c3f0a8b77/reativar",
  "codigo": "transicao_invalida",
  "id_rastreio": "0193f2a1-80aa-7c55-b1d2-6e3f4a5b6c7d",
  "estado_atual": "ativa"
}
```

E de validação, com erro por campo:

```json
{
  "type": "https://api.tenaz.com.br/erros/validacao_falhou",
  "title": "Os dados enviados não passaram na validação",
  "status": 422,
  "codigo": "validacao_falhou",
  "erros": [
    { "campo": "responsavel.email", "codigo": "obrigatorio", "mensagem": "Informe o e-mail." }
  ]
}
```

### 8.3 API

| Tema | Convenção |
|---|---|
| Paginação | Por cursor (`?cursor=...&limite=50`), nunca por offset em listas grandes |
| Idempotência | Cabeçalho `Idempotency-Key` aceito em todo `POST`; resposta repetida por 24 horas |
| Concorrência | `ETag` e `If-Match` em atualizações de recursos editáveis; 412 em conflito |
| Ações de estado | Sub-recursos verbais: `POST /v1/ordens-servico/{id}/concluir`, `POST /v1/ordens-servico/{id}/assumir`, `POST /v1/ordens-servico/{id}/encerrar` (permissão `os.encerrar`), `POST /v1/conexoes-empresa/{id}/aceitar` |
| Limites | Rate limit por plano; 429 com `Retry-After` |
| Versão | `/v1` estável. Mudança incompatível só em `/v2`, com convivência mínima de 12 meses |

Detalhes em `3_ESPECIFICACAO_APIS.md`.

### 8.4 TypeScript

- `strict: true`, `noUncheckedIndexedAccess: true`. `any` proibido pelo ESLint, salvo com comentário justificando.
- ESLint + Prettier; `typescript-eslint` com regras de tipo.
- Dados do servidor só via TanStack Query (web) ou SQLite local (app). Sem estado global para cópia de dados do servidor.
- Componentes funcionais; lógica de negócio em hooks ou em `pacotes/dominio`, não no JSX.
- Textos de interface só via `pacotes/i18n`. **Versão 0.8 (ADR-013, aceito):** todo texto existe em `pt-BR` e `es-419`; o `es-419.ts` é declarado com `satisfies Textos`, e a tarefa `textos` de `ferramentas/tarefas.sh` quebra o build por chave, variável ou forma de plural faltando em qualquer idioma. **Correção na versão 0.9 (conferida em `tarefas.sh`):** a tarefa `textos` confere só os catálogos do **servidor** (testes Go `TestConferirCatalogos`, `TestTodoCodigoTemTitulo` e `TestTodaTelaTemNome`, e a trava do texto livre da seção 8.2); na interface, chave faltando ou sobrando é apanhada pela checagem de tipos (`satisfies Textos`), e paridade, plural, variáveis e glossário pelos testes de `pacotes/i18n` (`es-419.test.ts`).
- **Dinheiro (versão 0.9, ADR-013, Parte F):** só em `pacotes/dominio/src/dinheiro.ts`, com `bigint` nativo; no contrato continua `string`. Proibidos `Number()`, `parseFloat` e aritmética de `number` sobre dinheiro.
- **Formatação para pessoas (versão 0.9, ADR-013, Parte G):** só pelas funções de `pacotes/dominio/src/formatos.ts`, pelo país da empresa e com fuso explícito. `Intl.NumberFormat`, `Intl.RelativeTimeFormat`, `Intl.ListFormat`, `Intl.DateTimeFormat` como gerador de texto e `toLocale*String` são barrados por regra do ESLint; o `Intl.DateTimeFormat` só aparece em `formatos.ts`, como motor de fuso.
- Acessibilidade: componentes Radix, rótulos em todo campo, navegação por teclado no web. Critérios em `7_PADROES_DESIGN.md`.

---

## 9. Testes e CI/CD

### 9.1 Pirâmide

| Tipo | Ferramenta | Escopo | Quando roda |
|---|---|---|---|
| Unitário Go | `testing` + testify | Domínio e casos de uso com dependências falsas | Todo push |
| Unitário TS | Vitest, Jest | Hooks, `pacotes/dominio`, componentes | Todo push |
| Casos compartilhados | JSON em `backend/testes/casos/` | Máquina de estados da OS, validação de leitura, gatilhos de preventiva | Todo push, nas duas linguagens |
| Integração | testcontainers-go com PostgreSQL real (mesma versão e extensões de produção) | Repositórios, migrações, RLS, fila, outbox | Todo push |
| **Vazamento entre empresas** | Suíte própria (9.3) | Banco e API | **Todo build; falha bloqueia o merge e o deploy** |
| **Isolamento entre contas conectadas** | Suíte própria com três contas (9.3) | Banco, sincronizador e API | **Todo build; falha bloqueia o merge e o deploy** |
| Contrato | Huma + `oasdiff` + Schemathesis | Drift do OpenAPI, mudança incompatível, fuzz da API | Todo PR |
| E2E web | Playwright | Fluxos críticos no navegador contra o ambiente efêmero | Todo PR para `main` |
| Sincronização do app | Testes do motor de sync com SQLite em Node | Offline longo, conflitos, reenvio, idempotência, troca de conta sem mistura de dados | Todo push |
| E2E app | Maestro em emulador Android | Fluxos de campo, inclusive modo avião | Diário e antes de release |
| Carga | k6 | Metas de p95 da API, sync em massa | Semanal em homologação e antes de release grande |

**Onde cada suíte Go roda no CI (versão 0.9, conferido em `.github/workflows/ci.yml`, job do backend).** (1) "Testes unitários": `go test ./internal/... -count=1`, que inclui os testes de integração dentro dos pacotes, com o PostgreSQL descartável de cada pacote (seção 8.2). (2) **"Testes de integração e de migração"**: `go test ./testes/integracao/ ./testes/migracoes/ -count=1`. `backend/testes/integracao/` tem os fluxos HTTP de ponta a ponta no backend (rotas, autenticação e MFA, idioma). `backend/testes/migracoes/` foi criado pelo `tenaz-banco-dados` e testa conversão, reversão e recusa das migrações de dados, aplicando e revertendo até uma versão com `migracao.AplicarAte` e `migracao.ReverterAte` (documento 2, seção 14.2). Até a E01, essas duas pastas rodavam só em `ferramentas/tarefas.sh tudo`, **fora do CI**; o passo entrou na E01 e foi conferido localmente. (3) "Suíte de isolamento": `go test ./testes/vazamento/ -count=1 -v`, que bloqueia merge e deploy (seção 9.3).

Metas de cobertura **(proposta)**: 80% de linhas nos pacotes `dominio` e `app`. A cobertura é sinal, não meta em si; os fluxos da seção 9.2 são obrigatórios independentemente do número.

### 9.2 Fluxos E2E obrigatórios

1. Solicitação aberta por solicitante → triagem → OS → execução → conclusão → encerramento, com custos consolidados.
2. Plano preventivo por "o que vencer primeiro" (tempo ou km) gera OS na hora certa.
3. OS de máquina NR-12 não conclui sem os sete campos do item 12.11.2.
4. Item reprovado em checklist gera solicitação.
5. Documento com vencimento dispara alertas D-60, D-30 e D-7.
6. Técnico executa OS em modo avião, reconecta, e o gestor vê o resultado sem perda.
7. OS atribuída a fornecedor: o usuário `prestador` vê no app e no web só as OS do seu fornecedor, executa, consome peça do almoxarifado `em_terceiro` dele, conclui; o gestor libera.
8. Link de uso único: gestor gera o link, prestador ocasional abre no celular, registra horas, fotos e assinatura, conclui; reabrir o link consumido, abrir em outro aparelho e usar link revogado são recusados.
9. (Fase 2) Remessa de peças a subcontratada e retorno; ordem de reparo `aberta` → `enviada` → `em_orcamento` → `aprovada` → `retornada`, com a peça voltando como recondicionada.
10. (Fase 2) Contrato de locação cedido: ativação, vistoria de entrega, solicitação aberta pelo locatário no portal, vistoria de devolução com avaria gerando OS com `responsavel_custo = locatario`.
11. SRV sem cliente na Tenaz: a subcontratada cadastra técnico com especialidade, cliente atendido e equipamento; abre OS de serviço com técnico e serviço a realizar; o técnico recebe a notificação push, executa em modo avião, registra componente substituído com leitura do número de série, fotos e assinatura do cliente, conclui; o gestor da subcontratada encerra e gera o relatório PDF.
12. SRV com cliente conectado: o cliente convida, a subcontratada aceita; o cliente atribui uma OS ao fornecedor conectado; a OS aparece na subcontratada em menos de 60 s; a subcontratada distribui; o técnico assume e executa; o andamento aparece no cliente com o selo da subcontratada e campos somente leitura; o cliente recusa a liberação, a OS volta a `em_execucao` nas duas contas; o técnico corrige; o cliente encerra e a OS da subcontratada encerra.
13. A subcontratada abre OS num ativo do cliente conectado: nasce `aberta` no cliente com `origem = prestador`; depois, a conexão é encerrada, nada mais trafega e as duas contas mantêm o histórico.
14. Regra e campo extra (versão 0.5): a empresa cria o campo extra `sim_nao` "Rompeu lacre?" obrigatório na conclusão e a regra "OS concluída com o campo → `bloquear_ativo_ate_documento`"; o técnico tenta concluir em modo avião sem preencher e o app recusa antes de enfileirar; preenche "sim", conclui, reconecta; o ativo fica `parado` com bloqueio de `origem = regra` e um `disparo_regra`; reprocessar o mesmo evento não cria segundo bloqueio; o registro do documento exigido, com emissão posterior, levanta o bloqueio sozinho.
15. Implantação (versão 0.5): ativo cadastrado `em_implantacao` fica `inativo` e fora dos KPIs; OS `implantacao` executada no app não entra em MTTR nem no custo de manutenção; `colocar_em_operacao` é recusado enquanto houver OS `implantacao` não encerrada; depois de aceito, o ativo fica `operando` e os planos ativam a partir de `data_entrada_operacao`.
16. (Fase 2) Projeto e contrato com cliente: projeto criado de um modelo libera as etapas pela ordem das dependências, gera OS `implantacao` atribuídas ao prestador, conclui e coloca os ativos em operação, incluídos no contrato com cliente; a medição mensal desconta só a parada com causa externa aceita e aprovada, e a disponibilidade técnica do ativo não muda.
17. (Fase 2) Pacote de segmento aplicado duas vezes à mesma conta: a segunda aplicação não cria nada e não altera nenhum item que a empresa editou depois da primeira; as regras do pacote nascem `inativa`.

### 9.3 Teste de vazamento entre empresas

Roda em todo build, com empresas fictícias de três países (A no Brasil, B na Argentina e C no Paraguai, ADR-013), populadas com dados em todas as tabelas. (Versão 0.8; até a 0.7, duas empresas, A e B.)

| Verificação | Como |
|---|---|
| Cobertura de RLS | Consulta `information_schema` e `pg_class`: toda tabela com coluna `empresa_id` precisa ter RLS habilitado **e** forçado, e políticas `USING` e `WITH CHECK`. Tabela nova sem política quebra o build. |
| Tabelas globais | Tabelas sem `empresa_id` (catálogos globais) ficam numa lista explícita revisada em code review. Tabela fora da lista e sem `empresa_id` quebra o build. |
| Leitura cruzada no banco | Conectado como `tenaz_app` com contexto A, `SELECT` em cada tabela filtrando pelo `empresa_id` de B deve retornar zero linhas. |
| Escrita cruzada no banco | Com contexto A, `INSERT` e `UPDATE` com `empresa_id` de B devem falhar. |
| Sem contexto | Sem `app.empresa_id`, toda tabela retorna zero linhas. |
| Leitura e escrita cruzadas na API (versão 0.8) | Guiadas pelo contrato: para **toda** rota autenticada do `openapi.yaml` (que precisa ser idêntico ao gerado), com a credencial de A pela web e pelo app e IDs reais de B e de C no caminho, inclusive o filho de B sob o pai de A: `404 recurso_nao_encontrado`, com a mesma resposta de um ID inexistente (menos `id_rastreio` e `instance`), e nenhum UUID nem texto de B. Referência a B no corpo e na consulta, busca e cursor de B: mesma resposta do inexistente, com controle positivo. `empresa_id` na consulta ou no corpo: 400. Sem credencial: 401. O retrato do banco de B e o de C antes e depois da varredura são os mesmos. As escritas válidas de A, inclusive as de rota única (`PATCH /v1/empresa`, `PATCH /v1/usuarios/eu`, `POST /v1/unidades` com endereço), são executadas com sucesso, com If-Match lido do GET da própria rota. A exceção da unicidade tem forma fixada: `PATCH /v1/empresa` com a identificação de outra empresa do mesmo país responde 409 só em `identificacao`, sem nada da outra empresa (`oraculoDeIdentificacao`). Rota nova entra sozinha; rota que precisa de caso quebra o build. Cada verificação tem prova de detector com um servidor estragado de propósito (`api_test.go`, `api_casos_test.go`, `api_contrato_test.go`, `api_banco_test.go`, `api_apoio_test.go`, `api_detector_test.go`). |
| Listagens e busca | Listas, filtros, busca textual e busca vetorial com token de A não retornam nada de B. Na API, toda lista e todo `GET` da própria empresa só trazem UUIDs da própria empresa; lista vazia é erro; a captura da trilha de B nunca aparece para A. As unidades das três empresas têm endereço com sentinela no município: um endereço alheio sem UUID é acusado pelo texto. |
| Arquivos | URL de arquivo de B não é assinada para usuário de A. |
| Jobs | Job enfileirado para A não lê nem grava dados de B. |
| Consultas sqlc | Toda consulta sqlc sobre tabela de cliente, sem comentários nem literais, tem `empresa_id = sqlc.arg(empresa_id)` como `AND` no nível de cima de todo `WHERE` (sem `OR` no nível de cima; subconsulta não conta) ou, no `INSERT`, na lista de colunas; nome entre aspas e tabela sem esquema são recusados. (Limite conhecido da suíte: a regra é por consulta, não por tabela da junção.) |
| Catálogos globais só leitura (versão 0.5) | Conectado como `tenaz_app`, `INSERT`, `UPDATE` e `DELETE` em `pacote_segmento` falham. A aplicação de um pacote em A não cria nem altera nada em B. |
| Regras e disparos (versão 0.5) | A rotina diária de regras de A não lê candidatos nem grava `disparo_regra`, bloqueio ou solicitação em B. |
| Visões (versão 0.6) | Em todo esquema que não seja do sistema (fora de `pg_catalog`, `information_schema` e `pg_*`): nenhuma visão materializada, e toda visão com `security_invoker = true`. Uma visão comum roda com os privilégios do dono e, sobre `plataforma.identidade`, entregaria `senha_hash` a quem tivesse `SELECT` nela. |
| Papéis e relações de membro (versão 0.6) | As relações de membro de e para os papéis de execução (`tenaz_app`, `tenaz_job`, `tenaz_leitura`, `tenaz_sinc`, `tenaz_sinc_conexao`, `tenaz_autenticacao`, `tenaz_autenticacao_conexao` e `tenaz_manutencao`) batem, nos dois sentidos e com as opções `INHERIT` e `SET`, com o arquivo versionado da suíte `relacoes_de_membro.txt`: relação a mais ou a menos quebra o build. Fica de fora só a concessão pura de `ADMIN` (sem `INHERIT` nem `SET`), que não transfere privilégio. Hoje o arquivo declara três: `tenaz_sinc_conexao` em `tenaz_sinc`, `tenaz_autenticacao_conexao` em `tenaz_autenticacao` e `tenaz_manutencao` em `tenaz_dono`. Todo papel de execução existe, e nenhum, inclusive o `tenaz_manutencao`, tem `BYPASSRLS` ou é superusuário. O par da autenticação ([ADR-011](#adr-011--papel-de-banco-próprio-para-as-credenciais-de-autenticação)): `tenaz_autenticacao` é `NOLOGIN`, `tenaz_autenticacao_conexao` é `LOGIN` e membro só dele, com herança, e nenhum outro papel o herda ou assume; nenhum dos dois tem `BYPASSRLS`, `CREATEROLE`, `CREATEDB` ou `REPLICATION` nem é superusuário; nenhum é dono de tabela, função, esquema ou tipo; nenhuma relação de membro, em nenhum sentido, com `tenaz_dono`, `tenaz_app`, `tenaz_job`, `tenaz_leitura`, `tenaz_sinc`, `tenaz_sinc_conexao` ou `tenaz_manutencao`; `USAGE` exatamente em `app`, `plataforma` e `public`. |
| Mapa de privilégios (versão 0.6) | O catálogo bate, nos dois sentidos e coluna a coluna, com o arquivo versionado da suíte `privilegios.txt`: privilégio a mais ou a menos quebra o build. Entram: todo privilégio em tabela, visão, sequência ou função **fora** de `app`, em qualquer esquema que não seja do sistema, concedido a um papel `tenaz_*` ou a `PUBLIC` (exceto o `EXECUTE` padrão de `PUBLIC` em função, coberto pela linha "Funções de credencial" e pela lista aprovada de funções `SECURITY DEFINER`); todo privilégio, em qualquer esquema, inclusive `app`, dos três papéis de lista fechada: sincronizador (`tenaz_sinc`, `tenaz_sinc_conexao`), autenticação (`tenaz_autenticacao`, `tenaz_autenticacao_conexao`) e `tenaz_credencial`; e todo privilégio de esquema (`USAGE`, `CREATE`) de papel `tenaz_*` ou de `PUBLIC` que não seja o dono do esquema, inclusive os privilégios padrão (`ALTER DEFAULT PRIVILEGES`). Nas sete tabelas de credencial, o arquivo declara só os privilégios do `tenaz_autenticacao` da seção 6.9 e o `REFERENCES (id)` do `tenaz_dono` em `identidade`; nele estão também os demais privilégios cruzados entre os dois donos (seção 6.9) e as colunas de `INSERT` do `tenaz_autenticacao` e do `tenaz_sinc` em `evento_auditoria`. Fora do mapa, para os papéis da aplicação (`tenaz_app`, `tenaz_job`, `tenaz_leitura` e os pares do sincronizador e da autenticação): nenhum `TRUNCATE`, `TRIGGER` ou `REFERENCES` em tabela de nenhum esquema, e nenhum `CREATE` em esquema nenhum. Toda verificação de privilégio vê também o concedido por coluna (versão 0.8). |
| Dono das credenciais (emenda ISO-E01-08, versão 0.6) | Todo objeto de `app` e `plataforma` (tabela, visão, sequência, tabela externa, índice, tipo de linha e função) pertence a `tenaz_dono`, e os dois esquemas também, com uma exceção exata: as sete tabelas de credencial, com os seus índices e tipos de linha, e as funções `plataforma.tg_acesso_empresa()` e `plataforma.purgar_credenciais_vencidas()` pertencem a `tenaz_credencial`. Objeto a mais ou a menos na exceção, ou objeto de credencial ausente, quebra o build. `tenaz_credencial` existe, não faz login, não tem `BYPASSRLS` nem `CREATEROLE` e não é superusuário; não é membro de nenhum papel; não tem membro, com `INHERIT` ou `SET`, além de superusuário ou do `tenaz_migrador`; e nenhum papel de runtime o alcança, nem por herança transitiva: `tenaz_app`, `tenaz_job`, `tenaz_leitura`, `tenaz_sinc`, `tenaz_sinc_conexao`, `tenaz_autenticacao`, `tenaz_autenticacao_conexao`, `tenaz_manutencao` e `tenaz_dono`. Como `tenaz_manutencao` e como `tenaz_dono`, `SELECT` em cada uma das sete tabelas e a troca da senha da identidade de B recebem "permissão negada". O trigger `trg_acesso_empresa` existe em `app.usuario`, não está desabilitado e chama `plataforma.tg_acesso_empresa()`. Em `plataforma`, tabela da classe `credencial` em `tabelas_globais.txt` tem dono `tenaz_credencial`, e tabela sem RLS habilitado e forçado e de dono `tenaz_dono` só pode ser da classe `catalogo` ou `controle` (regra da pendência E01-14). |
| Ataques do achado ISO-E01-02 (versão 0.6) | Como `tenaz_app`, `tenaz_job`, `tenaz_leitura`, `tenaz_sinc`, `tenaz_sinc_conexao`, `tenaz_dono` e `tenaz_manutencao`, com contexto interno de A: ler a identidade (com `senha_hash` e segredo de MFA), os acessos e as sessões de B, criar sessão web ou código OAuth para o usuário de B, trocar a senha da identidade de B e `SELECT`, `INSERT`, `UPDATE` e `DELETE` em cada uma das sete tabelas recebem "permissão negada", e não zero linhas. |
| Auditoria pelo papel de autenticação (versão 0.6) | Como `tenaz_autenticacao`, com contexto A: evento de outra entidade ou ação, evento com o usuário de B ou com usuário que não é o da identidade em A, e qualquer `SELECT` na trilha são recusados. |
| Consultas de credencial (versão 0.6) | Nenhum arquivo `.sql` fora do repositório da autenticação cita as sete tabelas; toda consulta do repositório filtra por hash, `id`, `identidade_id`, `(empresa_id, usuario_id)` ou e-mail exato. |
| Funções de credencial (versão 0.6) | `plataforma.tg_acesso_empresa()` e `plataforma.purgar_credenciais_vencidas()` na lista aprovada de funções `SECURITY DEFINER`, de dono `tenaz_credencial`, com `search_path` fixo e `EXECUTE` revogado de `PUBLIC` (a purga, só para `tenaz_job`). Inativar um usuário em A revoga as sessões e os refresh tokens dele em A e não toca os do mesmo usuário em B. Nenhuma função `SECURITY DEFINER` fora da lista aprovada em esquema nenhum, sem exceção para extensões; nenhuma função fora de `app` e `plataforma`, exceto as de extensão declarada em `extensoes.txt`, com a impressão das suas funções. |
| Moldes da trilha (versão 0.6) | As três restritivas de `INSERT` em `app.evento_auditoria`, cada uma com o seu molde (o que passa e o que a política recusa, conferido pela mensagem de RLS, e não por outro erro) e a sua prova (sem a política, o teste acusa). `auditoria_sem_autenticacao`: `tenaz_app`, `tenaz_job` e `tenaz_sinc` não gravam evento de `entidade = 'identidade'` (ISO-E01-07), e o `tenaz_app` grava evento de outra entidade. `escopo_externo_insercao`: prestador, locatário e link gravam só em nome de si e no próprio escopo (a coluna `escopo` do evento confere com o escopo da transação), sem atribuir o evento a outro usuário, a colega do mesmo fornecedor ou a ninguém; o interno grava evento de qualquer usuário da empresa. `conexao_sinc`: `tenaz_sinc` e `tenaz_sinc_conexao` gravam só evento da conexão em processamento, com ator `sincronizacao_conexao`, e não gravam evento de outra conexão, sem conexão, de outro ator, sem a conexão no contexto ou de B na transação de A; o sincronizador não lê a trilha. Como o `tenaz_sinc` ainda não tem `USAGE` em `app` (E04), o teste o concede só dentro da transação desfeita, para que a recusa venha da política e não do esquema. |
| Corrida de atribuição (ISO-E01-06, versão 0.6) | Duas conexões reais, no isolamento padrão (READ COMMITTED), como `tenaz_app` com contexto interno de A: a primeira atribui a um usuário interno o perfil `gestor` ou a permissão avulsa `os.encerrar`, vedada ao escopo externo, e não confirma; a segunda tenta trocar o `fornecedor_id` desse usuário com `lock_timeout` de 300 ms e precisa receber `55P03` (`lock_not_available`), pela trava `FOR SHARE` dos triggers de atribuição. A prova troca o trigger de atribuição de perfil por uma versão sem a trava, e a troca passa. |
| Contexto no pool (versão 0.6) | Com o helper real `db.NaEmpresa`, num pool de uma conexão como `tenaz_app`: depois de uma transação confirmada, de uma abortada por erro de SQL, de uma desfeita pelo caso de uso e de uma de prestador, a mesma conexão, fora de `NaEmpresa`, tem `app.empresa_atual()`, `app.escopo_atual()`, `app.usuario_atual()` e `app.fornecedor_atual()` nulos e vê zero linhas em toda tabela de empresa; e `NaEmpresa` de B, logo depois de A na mesma conexão, não vê a unidade de A. Antes, o teste exige que `NaEmpresa` de A veja a unidade de A, para não passar por vacuidade. |
| NaEmpresa confere o escopo no banco (versão 0.7) | Escopo bem formado e incoerente com o cadastro é recusado antes do caso de uso (pool de uma conexão; nenhum contexto sobra na conexão). |
| Extensões (versão 0.7) | Só as de `extensoes.txt`, com versão e impressão das funções. |
| Sem dado de saúde (versão 0.7) | Nenhum objeto com `saude` em `app` e `plataforma`, nenhuma leitura de `app.saude`, `pessoa.saude.ler` recusada pelo formato de permissão. |
| País, identificação e idioma (versão 0.8, `pais_test.go`) | Empresas de países diferentes não se alcançam (os seis pares entre A, B e C, em toda tabela). País e moeda da empresa imutáveis, até para o superusuário. O usuário externo altera na própria linha só nome, telefone e idioma, conferido coluna a coluna de `app.usuario`; coluna nova entra recusada. Cada regra tem prova de detector própria: a dos seis pares desliga a RLS de `app.unidade` numa transação desfeita e exige a acusação dos seis pares e do total (`TestDetectorAcusaPaisesQueSeAlcancam`). As fixtures gravam a moeda do usuário igual à da empresa, e a FK composta da 00010 recusa moeda divergente (`TestBancoRecusaUsuarioEmMoedaDeOutraEmpresa`). |
| FK composta (versão 0.8) | Toda FK para tabela de empresa leva `empresa_id` casado; para `app.empresa`, só `empresa_id` na primeira posição apontando para `id`, sozinha ou com a moeda. Outra coluna para `app.empresa`, ou `empresa_id` para outra coluna, quebra o build (`TestDetectorAcusaFKSimples`, `TestDetectorAcusaFKParaEmpresaPorOutraColuna`), porque a FK ignora o RLS. |
| Trilha configurável e selo (versão 0.8) | Evento sem data retroativa: `tenaz_app` e `tenaz_job` só gravam as colunas do `GRANT` (sem `id`, `ocorrido_em`, `transacao` e `hash`), e o servidor as preenche; hash do evento não forjável (o trigger troca o que o migrador mandar); `selo_auditoria` somente inserção, `INSERT` só do `tenaz_job`, imutável até para o superusuário; `selar_dias_pendentes` só para o `tenaz_job`; selo não forjável; verificação de integridade por empresa; `captura_auditoria` só interna; evento de identidade com `usuario_nome` conferido com `acesso_empresa` (`selo_auditoria_test.go`, `autenticacao_test.go` e os testes gerados do catálogo). |
| Portão de achados (versão 0.6) | Falha de isolamento já demonstrada por um teste e devolvida ao dono da migração fica no mapa `achadosAbertos` da suíte, com o ID do achado. Não é lista de exceções: no CI (variável `CI` definida), o teste do achado aberto falha e a suíte encerra com erro enquanto o mapa não estiver vazio; fora do CI, o teste aparece como SKIP com o ID. Numa execução completa (sem `-run`), chave do mapa que nenhum teste exercitou é erro (achado sem prova), e o achado que deixa de se reproduzir faz o teste falhar pedindo a retirada. |

**Nota (versão 0.7):** a suíte roda igual em contêiner descartável e no banco local compartilhado entre pacotes (`ferramentas/tarefas.sh testar-rapido`). As provas de detector rodam numa transação do superusuário que é desfeita no fim do teste, e nenhum objeto de prova fica visível para outra conexão. O helper `estragar` aplica os comandos num savepoint com `lock_timeout` de 1 segundo e, em `55P03` ou `40P01`, repete com recuo, até 8 vezes: a prova espera a vez e não falha por concorrência alheia. A prova que altera tabela real (dono, trigger ou coluna) segura a trava só até o fim da própria transação.

**Isolamento por fornecedor e por locatário** (risco R4), na mesma suíte e com a mesma regra de bloqueio. As sementes têm, na empresa A, os fornecedores F1 e F2 (com OS, almoxarifados `em_terceiro` e usuários `prestador`) e os locatários L1 e L2 (com contratos, vistorias e solicitações); na empresa B, um fornecedor F3 com o **mesmo CNPJ** de F1, para simular o prestador que atende duas contratantes.

| Verificação | Como |
|---|---|
| Cobertura da política restritiva | Toda tabela com `empresa_id` tem a política `escopo_externo` `AS RESTRICTIVE`. Tabela nova sem ela quebra o build. |
| Lista de acesso externo | As tabelas cujo predicado aceita escopo externo batem com a lista da seção 6.6, guardada no repositório. Divergência quebra o build. |
| Leitura cruzada no banco | Com escopo F1, `SELECT` em cada tabela retorna zero linhas de F2, de OS internas e de tabelas fora da lista; com escopo L1, zero linhas de L2 e de ativos não locados a L1; com escopo `link`, só a OS do link. |
| Escrita cruzada no banco | Com escopo F1, `UPDATE` que troca `fornecedor_id` para F2 ou `executor_tipo` para `interno` falha; movimentação em almoxarifado de F2 ou `proprio` falha. Com escopo L1, qualquer escrita em `ordem_servico` falha. O prestador e o locatário alteram na própria linha de `usuario` só nome, telefone e idioma (versão 0.8). |
| Leitura cruzada na API (versão 0.8) | Para toda rota do contrato, com a credencial (web e app) do prestador F1, do locatário L1 e do prestador de B com o fornecedor de F1: em rota que não declara o escopo, 404 igual ao de ID inexistente e ao de rota inexistente; nas rotas declaradas, nenhum UUID fora do próprio alcance. A escrita válida do externo em rota declarada para ele (`PATCH /v1/usuarios/eu`, trocando o idioma) dá certo pela web e pelo app e responde só o próprio alcance. |
| Empresa cruzada | Token do usuário de F1 na empresa A não lê nada de F3 na empresa B, e vice-versa. |
| Período de locação | L1 não vê OS do ativo abertas antes da entrega nem depois da devolução. |
| Link de uso único | Token inválido, expirado, revogado, consumido ou já ativado em outra sessão: mesma resposta 404; o link não acessa outra OS mesmo trocando IDs na requisição. |
| Sincronização do app | Pull com usuário F1 não traz nenhuma linha de F2 nem de OS internas. |
| Arquivos e cache | Anexo de OS de F2 não recebe URL assinada para F1; chaves de cache distintas por escopo. |
| Campos extras (versão 0.5) | Com escopo F1, `campo_extra_os` devolve só as definições das OS de F1; com escopo `link`, só as marcadas como do prestador; com escopo L1, nenhuma. Escrita em `campo_extra_os` sob escopo externo falha. |

Qualquer falha é tratada como incidente de segurança de severidade alta, mesmo no CI.

**Isolamento entre contas conectadas, com três contas** (risco R5), na mesma suíte e com a mesma regra de bloqueio. As sementes têm três empresas: os clientes **C1** e **C2** (plano Planta, cada um com o fornecedor que representa S) e a subcontratada **S** (plano Serviços), conectada a C1 e a C2 por duas conexões `ativa` distintas. C1 e C2 compartilham com S ativos diferentes, com campos internos preenchidos (custos, histórico, pessoas, documentos) que **não** podem trafegar. S tem um cliente sem conta Tenaz e técnicos com CPF e dados pessoais preenchidos. O teste roda o sincronizador real (papel `tenaz_sinc`) contra o PostgreSQL do testcontainers.

| Verificação | Como |
|---|---|
| Cobertura do papel | Toda tabela com `GRANT` para `tenaz_sinc` tem a política `conexao_sinc AS RESTRICTIVE`. As colunas com `GRANT` batem exatamente com `campos_conexao.yaml`. Tabela ou coluna fora da lista com privilégio quebra o build. `tenaz_sinc` sem `BYPASSRLS` e sem ser dono de tabela |
| Cópia além do permitido | Depois de um ciclo completo (convite, OS atribuída, execução, liberação), as linhas em S vindas de C1 contêm só os campos da lista; nenhum valor sentinela dos campos internos de C1 aparece em nenhuma tabela de S; nenhum CPF, contato pessoal ou documento de técnico de S aparece em C1. Grupo opcional desligado não trafega; ligado, trafega |
| Campo injetado | Mensagem montada à mão com campo extra (custo, CPF) ou com campo de outro dono (S alterando a descrição da OS de C1; S enviando encerramento): recusada inteira, `erro` registrado e alerta de segurança emitido |
| Conta errada | Job com conexão C1–S e `empresa_origem` = C2: recusado (C2 não é ponta da conexão). Nenhuma mensagem da conexão C1–S grava em C2, e vice-versa. O destino é sempre calculado a partir da conexão, nunca do job |
| Sem transitividade | C1 e C2 conectados à mesma S: nada de C1 chega a C2 por meio de S, nem por OS, nem por ativo, nem por relatório |
| Leitura cruzada do sincronizador | Com contexto C1 e conexão C1–S, `tenaz_sinc` não lê nenhuma linha de S nem OS de C1 sem vínculo; consulta a tabela fora da lista recebe "permissão negada" |
| Replay e ordem | A mesma mensagem entregue duas vezes é aplicada uma vez; sequência antiga é ignorada; mensagens fora de ordem são aplicadas na ordem da sequência |
| Falsificação de origem | Linha de `envio_conexao` inserida por `tenaz_app` em C2 apontando para a conexão C1–S é recusada pelo sincronizador; `envio_conexao` com `vinculo_os_id` de outra conexão é recusado |
| Estados da conexão | `convidada`, `suspensa` e `encerrada` não aplicam nada (fora da janela de 24 h do encerramento); reativação entrega o acumulado em ordem; o histórico continua nas duas contas depois do encerramento |
| Liberação | Só C1 encerra a OS vinculada; em S, `encerrar` local de OS com conexão ativa é recusado |
| API e app | Token de usuário de S com IDs de OS, ativo ou conexão de C1 ou C2: 404. Pull de um técnico de S traz só dados de S. `GET /v1/conexoes-empresa` em C2 não mostra a conexão C1–S |
| Cobrança | Técnicos de S não contam como técnicos cobrados em C1 nem em C2 |
| Auditoria | Cada mensagem aplicada tem registro na origem e no destino com o mesmo `hash_conteudo` |
| Campos extras e regras (versão 0.5) | OS de C1 com `campos_extras` preenchidos com valores sentinela, atribuída a S e executada por S com campos extras próprios: nenhum valor sentinela aparece em S e nenhum valor de S aparece em C1. Regra ativa em C1 não grava nada em S, e regra ativa em S não grava nada em C1 |

### 9.4 Pipeline

```mermaid
flowchart LR
    pr["PR"] --> lint["Lint e formatação<br/>Go, TS, SQL, OpenTofu"]
    lint --> seg["Segurança<br/>SAST, SCA, segredos, IaC"]
    seg --> test["Unitários, integração,<br/>vazamento, contrato"]
    test --> build["Build da imagem<br/>SBOM, scan, assinatura"]
    build --> efem["Ambiente efêmero<br/>E2E Playwright"]
    efem --> merge["Merge em main"]
    merge --> dev["Deploy em dev<br/>(automático)"]
    dev --> hml["Deploy em homologação<br/>(automático) + DAST"]
    hml --> aprov["Aprovação manual"]
    aprov --> prd["Produção<br/>rolling com rollback por alarme"]
```

| Etapa | Detalhe **(proposta)** |
|---|---|
| Plataforma | GitHub Actions, com actions fixadas por SHA |
| Testes do backend no CI (versão 0.9, em uso) | Três passos em sequência no job do backend: unitários e integração por pacote (`./internal/...`), integração e migração (`./testes/integracao/` e `./testes/migracoes/`, desde a E01) e isolamento (`./testes/vazamento/`). Detalhe na seção 9.1 |
| Imagem | Construída uma vez, promovida entre ambientes pelo digest; nunca reconstruída para produção. Na promoção para produção recebe também a tag da versão (seção 9.5) |
| Versão | Tag `vX.Y.Z`, CHANGELOG e versão injetada no build conferidos antes do release (seção 9.5) |
| Migrações | Tarefa única (`tenaz migrate`) antes do deploy. Toda migração é compatível com a versão anterior do código (expandir e contrair). Remoção de coluna só no release seguinte ao que parou de usá-la. |
| Deploy da API e dos workers | Rolling no ECS com health check e rollback automático se o alarme de erro ou de latência disparar nos primeiros 10 minutos |
| Web | Upload versionado no object storage + invalidação da CDN; rollback trocando o ponteiro de versão |
| App | EAS Build para as lojas; EAS Update para mudanças só de JavaScript compatíveis com a versão nativa (`runtimeVersion`). Canais `homologacao` e `producao`; liberação gradual por porcentagem. |
| Frequência | Deploy de backend e web várias vezes por semana; app nas lojas a cada 2 a 4 semanas |
| Janela | Sem janela de manutenção para deploy comum. Deploy com risco alto fora do horário comercial. |

Os controles de segurança do pipeline estão em 9_SEGURANCA_CONFIABILIDADE.md.

### 9.5 Versionamento e release

Requisitos do dono do produto RNF-039 (versão visível) e RNF-040 (CHANGELOG.md sempre atualizado), em [1_ESPECIFICACAO_REQUISITOS.md](1_ESPECIFICACAO_REQUISITOS.md): toda entrega tem uma versão única, rastreável até o commit, visível no canto inferior direito da tela principal e da tela de login, no web e no app (sumário 3.4), e registrada no histórico de versões.

**Regras de versão**

| Regra | Definição |
|---|---|
| Formato | **SemVer 2.0.0**: `MAJOR.MINOR.PATCH` (ex.: `1.4.2`). Pré-lançamento `1.5.0-rc.1` só fora da produção. Metadado de build `+a1b2c3d` (commit curto) opcional, nunca usado para comparar versões. |
| Incremento | MAJOR = quebra da API pública (coincide com uma nova versão de caminho, `/v2`, seção 8.3); MINOR = funcionalidade nova compatível; PATCH = correção. |
| Uma versão por release | Web, app, API e workers do mesmo release têm a **mesma versão**. Ela vem da tag git `vX.Y.Z` e **nunca é digitada à mão** em arquivo de código. |
| Artefatos de cada release | Tag anotada `vX.Y.Z` em `main`; seção no `CHANGELOG.md` no formato *Keep a Changelog* (Adicionado, Alterado, Obsoleto, Removido, Corrigido, Segurança), com data; imagem do contêiner com a tag da versão, além do digest. |
| Histórico de versões (RNF-040) | O `CHANGELOG.md` na raiz registra **todas** as implementações (com ID `RF-`/`RNF-`) e correções (com número do chamado) de cada versão. Toda mudança de comportamento entra em `[Não lançado]` no mesmo pull request; um check do CI falha se o PR altera código de produto (`backend/`, `web/`, `app/`, migrações) sem alterar o `CHANGELOG.md`, salvo rótulo `interno` com justificativa (refatoração, testes, ferramentas). No release, `[Não lançado]` vira a seção da versão. O build publica a seção da versão como `novidades.json` (web e app) para o diálogo "Novidades desta versão". |
| Sugestão do incremento | Os Conventional Commits (8.1) desde a última tag sugerem o incremento: `fix` → PATCH, `feat` → MINOR, `!` ou `BREAKING CHANGE` → MAJOR. Uma ferramenta de release (release-please ou git-cliff, **proposta**) abre o PR de release com a versão sugerida e o rascunho do CHANGELOG; uma pessoa revisa o texto e aprova. O merge desse PR cria a tag. |
| Ambientes | `dev` e efêmeros usam a próxima versão com sufixo de pré-lançamento e o commit (`1.5.0-dev.37+a1b2c3d`); `hml` recebe `rc`; `prd` só aceita versão sem pré-lançamento. |

**Injeção da versão no build**

| Componente | Como **(proposta)** |
|---|---|
| Backend Go | `go build -ldflags "-X tenaz/internal/plataforma/versao.Versao=1.4.2 -X tenaz/internal/plataforma/versao.Commit=a1b2c3d -X tenaz/internal/plataforma/versao.DataBuild=2026-11-12T14:05:00Z"`. O subcomando `tenaz versao` imprime os três valores; `debug.ReadBuildInfo` confere o commit gravado pelo próprio Go. A versão vai nos logs (campo `versao`, 14.3), nos traces e no cabeçalho `Tenaz-Versao` das respostas. |
| Web | Variável de build do Vite (`define` com `__TENAZ_VERSAO__`, `__TENAZ_COMMIT__`, `__TENAZ_DATA_BUILD__`, `__TENAZ_AMBIENTE__`), preenchida pelo CI a partir da tag. O bundle publica também `versao.json` junto dos arquivos estáticos, usado para detectar que há versão mais nova. |
| App | `app.config.ts` lê as variáveis do CI: `version` = versão do release que gerou o binário (`1.4.2`); `ios.buildNumber` e `android.versionCode` = inteiro crescente gerado pelo CI (ex.: número da execução), exigido pelas lojas; `runtimeVersion` com a política `fingerprint` do Expo, que muda só quando a parte nativa muda. A versão do release também é embutida no JavaScript (`extra.versaoRelease`), para o selo mostrar a versão que está de fato rodando. |
| Endpoint | `GET /v1/versao`, público e sem autenticação, devolve só `versao`, `commit` (curto) e `data_build` da API. Nada de dependências, sistema operacional, região ou infraestrutura (controles em 9_SEGURANCA_CONFIABILIDADE.md). Não devolve `ambiente`: o ambiente que o selo e o "Sobre" mostram (RNF-039, itens 2 e 3) vem do build do próprio cliente (`__TENAZ_AMBIENTE__` na web, `TENAZ_AMBIENTE` no `app.config.ts`), que é publicado no mesmo ambiente da API. Posição deste documento na divergência A-38 de `3_ESPECIFICACAO_APIS.md` (seção 16). |

**App nas lojas e atualização OTA**

- Mudança só de JavaScript compatível com o binário instalado (mesmo `runtimeVersion`) vai por EAS Update: o aparelho passa a rodar o JavaScript do release novo e o selo mostra a versão nova (`v1.4.3`), enquanto o binário da loja continua `1.4.0`. O diálogo Sobre mostra os dois (versão do release e binário da loja com build) e o identificador da atualização OTA. **(Versão 0.7: o `expo-updates` ainda não está no app, porque falta decidir o servidor de atualização, ligado à pendência H5 da seção 16; até lá não há identificador de atualização, e o Sobre mostra só a versão do release e a do binário. O `runtimeVersion` já usa a política `fingerprint`.)**
- Mudança nativa gera `runtimeVersion` nova e exige binário novo nas lojas; o OTA desse release só chega a quem já instalou o binário novo.
- Como o app trabalha até 7 dias offline, um aparelho pode estar alguns releases atrás. A API aceita apps do mesmo MAJOR; abaixo da versão mínima suportada (informada pela API na sincronização), o app pede a atualização sem descartar a outbox. **(proposta)**

**Verificação no CI** (bloqueia o release)

1. A tag segue `vMAJOR.MINOR.PATCH` (ou `-rc.N` para `hml`) e é maior que a anterior.
2. O `CHANGELOG.md` tem a seção `## [X.Y.Z] - AAAA-MM-DD` correspondente, não vazia.
3. A versão impressa por `tenaz versao`, a do `versao.json` do web e a do `app.config` gerado são iguais à tag, e o commit é o da tag.
4. Se o `oasdiff` apontar mudança incompatível na API pública, o incremento tem de ser MAJOR; senão, MAJOR é recusado.
5. O deploy em `prd` recusa versão com pré-lançamento e versão sem tag.

Depois do deploy, a verificação sintética chama `GET /v1/versao` e confirma que a produção responde a versão liberada.

---

## 10. App offline-first

### 10.1 Metas

Vindas do sumário (seção 5) e do estudo de mercado:

- Operar **7 dias** sem conexão, sem perda de registro.
- Convergir em **menos de 60 segundos** após a reconexão (com rede 4G razoável e até 7 dias de trabalho acumulado).
- Criar OS e solicitações offline, não só executar as existentes.

### 10.2 O que fica no aparelho

O escopo de sincronização é calculado por usuário, a partir das unidades e dos perfis dele. **(proposta)**

| Dado | Escopo |
|---|---|
| OS | Atribuídas ao usuário ou à equipe dele, em aberto, mais as encerradas nos últimos 30 dias para consulta |
| Ativos | Árvore completa das unidades do usuário, com medidores e última leitura |
| Checklists | Modelos aplicáveis às classes de ativo das unidades; execuções do usuário nos últimos 30 dias |
| Procedimentos | Os referenciados pelas OS e planos em escopo |
| Estoque | Itens e saldos dos almoxarifados das unidades (saldo é informativo offline) |
| Códigos de falha | Catálogo das classes em escopo |
| Documentos e manuais | Só os marcados como "disponível offline", com limite de tamanho por empresa |
| Fotos de OS antigas | Miniaturas; original sob demanda |
| Cache de configuração (versão 0.5) | Definições ativas de `campo_extra_os` que se aplicam aos tipos de OS e às classes de ativo em escopo, com a versão de cada uma |

Usuário `prestador` no app: o escopo é calculado sob o escopo `prestador` (seção 6.6). Vão para o aparelho só as OS atribuídas ao fornecedor dele (em aberto e encerradas nos últimos 30 dias), os ativos dessas OS com medidores, os procedimentos referenciados e o saldo do almoxarifado `em_terceiro` do fornecedor. Nenhum dado de outra OS, de outro fornecedor ou de cadastro interno. **(proposta)**

Técnico da subcontratada com conta própria (SRV): é usuário `interno` da conta dela, e o escopo segue a tabela acima, com o significado do lado de quem presta serviço: OS atribuídas a ele ou à equipe dele (inclusive as vinculadas a clientes conectados) e as `aberta` ou `planejada` da equipe que ele pode assumir; os equipamentos de clientes (`propriedade = de_cliente`) dessas OS, com medidores; o `cliente_servico` de cada OS (nome, endereço de atendimento, contato); o catálogo `servico_fornecedor` da própria conta; o catálogo de itens e o saldo do estoque básico. Ele não recebe a carteira inteira de clientes nem os equipamentos de clientes sem OS para ele. **(proposta)**

**Cache de configuração (versão 0.5).** As definições de campos extras da OS (RF-OS-021) vão ao aparelho no **cache de configuração**, pelo mesmo pull incremental das demais tabelas (seção 10.4), sob o RLS e o escopo do usuário: o prestador recebe só as definições das OS do seu fornecedor (seção 6.6). O app não edita configuração. Os ativos `em_implantacao` das unidades do usuário também vão ao aparelho, porque recebem OS `implantacao`. As regras configuráveis (`regra_automacao`) **não** vão ao aparelho e não são avaliadas nele: o servidor avalia quando aplica o comando, e o efeito (por exemplo, o ativo `parado`) chega no pull seguinte (seção 5.6). **(proposta)**

### 10.3 Modelo: comandos para cima, alterações para baixo

```mermaid
sequenceDiagram
    participant UI as Tela do app
    participant DB as SQLite local
    participant OB as Outbox de comandos
    participant API as API /v1/sync
    participant PG as PostgreSQL

    UI->>DB: aplica mudança otimista
    UI->>OB: grava comando (id UUIDv7, tipo, payload, versão base)
    Note over OB: sem conexão: comandos acumulam
    OB->>API: POST /v1/sync/comandos (lote, em ordem)
    API->>PG: aplica cada comando via casos de uso (mesmas regras da API)
    API-->>OB: resultado por comando: aplicado, rejeitado ou ajustado
    OB->>DB: remove aplicados, marca rejeitados para o usuário
    DB->>API: GET /v1/sync/alteracoes?cursor=...
    API-->>DB: linhas alteradas e exclusões desde o cursor + novo cursor
```

**Comandos**, não diffs de linha. O app envia intenções de domínio: `AssumirOS` (técnico se atribui OS `aberta` ou `planejada` da sua equipe; grava `data_reconhecimento`), `IniciarOS`, `PausarOS`, `ApontarHoras`, `RegistrarLeitura`, `ResponderItemChecklist`, `RegistrarComponente` (com `tipo_material` e números de série retirado e instalado), `ConcluirOS`, `AbrirOS` (na subcontratada, com equipamento, técnico e serviço a realizar), `AbrirSolicitacao`, `AnexarArquivo`, `Assinar` (técnico) e `ColetarAssinaturaCliente`. O servidor aplica cada comando pelos mesmos casos de uso da API pública. Não existe caminho de escrita exclusivo do app.

**Campos extras nos comandos (versão 0.5).** Os valores dos campos extras viajam no payload de `AbrirOS` e `ConcluirOS`, e o preenchimento durante a execução usa o comando `RegistrarCamposExtras` **(proposta)**. Antes de **enfileirar** qualquer um desses comandos na outbox, o app valida os valores contra as definições do cache de configuração, com a mesma função compartilhada de `pacotes/dominio` (seção 7): chave, tipo, opções e obrigatoriedade na abertura ou na conclusão. Comando inválido não entra na fila; o técnico corrige na hora, mesmo sem sinal (RF-OS-021). O comando leva a versão das definições usadas, e o servidor revalida tudo ao aplicar.

| Elemento | Definição **(proposta)** |
|---|---|
| ID do comando | UUIDv7 gerado no app; é a chave de idempotência |
| Idempotência | Tabela `comando_processado(empresa_id, comando_id)` com chave única. Reenvio devolve o mesmo resultado sem reaplicar. |
| Ordem | FIFO por agregado (mesma OS, mesmo ativo). Comando com dependência (`ConcluirOS` depois de `IniciarOS`) espera o anterior. |
| Versão base | O comando informa a versão do registro que o usuário viu. Serve para resolver conflito de campo editável. |
| Horários | Cada comando leva o horário do aparelho; o servidor grava também o horário de recebimento. A diferença fica registrada (relógio errado não apaga a evidência). |
| Lote | Até 200 comandos por requisição; o resultado vem por comando |
| Arquivos | Upload separado, por URL pré-assinada, retomável (multipart). O comando referencia o `anexo_id`; o servidor aceita o comando antes do arquivo chegar e marca o anexo como pendente. Fotos comprimidas no aparelho antes do envio. |

### 10.4 Pull incremental

- Cada tabela sincronizável grava, por gatilho, o ID da transação que alterou a linha (`txid_alteracao xid8`).
- O cursor do app é um **snapshot do Postgres** (`pg_current_snapshot()`), não um número de sequência.
- A consulta devolve as linhas alteradas por transações que não estavam visíveis no snapshot anterior. Isso evita o problema clássico do cursor por sequência: perder linhas de uma transação longa que confirmou depois de outra mais nova.
- Exclusões e saídas de escopo (OS reatribuída a outro técnico) chegam como registros de remoção (tombstones) com retenção de 30 dias. Aparelho que ficar mais de 30 dias sem sincronizar faz ressincronização completa.
- Tudo dentro de `NaEmpresa`: o RLS também protege o pull.

**(proposta, validar com prova de conceito antes de fechar o MVP)**

### 10.5 Conflitos resolvidos por regra de domínio

Não há CRDT genérico. Cada tipo de dado tem uma regra explícita.

| Situação | Regra |
|---|---|
| Apontamento de horas, foto, comentário, leitura, resposta de checklist | Append-only. Nunca conflitam: todos são aceitos. |
| Leitura de medidor fora do padrão (regressão, salto acima do teto) | Aceita e marcada como `a_revisar`, nunca descartada. PRE ignora leitura não revisada no cálculo de gatilho. |
| OS cancelada ou encerrada no servidor e concluída offline | `ConcluirOS` é rejeitado. Apontamentos, fotos e consumo de peças são preservados e ligados à OS; o gestor recebe uma pendência para decidir. |
| OS reatribuída a outro técnico enquanto o primeiro trabalhava offline | Os registros do primeiro técnico são aceitos (o trabalho aconteceu). O gestor é avisado da execução dupla. |
| Campo descritivo editado nos dois lados (descrição, prioridade) | Por campo: se o campo não mudou no servidor desde a versão base, vale o do app; se mudou, vale o do servidor e o valor do app vira comentário na OS. |
| Consumo de peça sem saldo suficiente ao sincronizar | Aceito (a peça foi usada). Gera pendência de ajuste para o almoxarife. |
| Transição de estado impossível | Rejeitada com motivo legível; o app mostra na lista "Pendências de sincronização". |
| Assinatura offline | Aceita se o conteúdo assinado (hash) ainda corresponde ao registro. Se o registro mudou no servidor, a assinatura é guardada como evidência e o responsável é chamado a assinar de novo. |
| Campo extra obrigatório criado ou alterado depois do último pull do aparelho (versão 0.5) | `ConcluirOS` validado no aparelho com as definições que ele tinha é aceito (o trabalho aconteceu); a OS fica `concluida` com uma pendência de preenchimento para o gestor antes da liberação. Valor de tipo ou opção inválidos para a definição que o aparelho tinha é rejeitado com motivo. **(proposta; regra a numerar em `4_FLUXOS_PROCESSO.md`)** |

Nenhum comando rejeitado desaparece em silêncio. O app mostra o motivo e o que foi preservado.

### 10.6 Segurança e sessão offline

- Banco local criptografado; chave no Keychain ou Keystore.
- Sessão do app com token de atualização válido por 30 dias de inatividade (RNF-017 em [1_ESPECIFICACAO_REQUISITOS.md](1_ESPECIFICACAO_REQUISITOS.md); até a versão 0.5 este documento dizia 14 dias, corrigido na etapa E01), rotativo, com detecção de reuso e ligado ao aparelho. Sem conexão, o app continua funcionando até o fim desse prazo. As credenciais ficam nas tabelas da seção 6.9.
- Access token do app (versão 0.7, `plataforma/autenticacao/jwt.go`; formato na seção 4.3 de [3_ESPECIFICACAO_APIS.md](3_ESPECIFICACAO_APIS.md)): a verificação aceita um formato só e recusa cabeçalho diferente de `alg` EdDSA, `typ` JWT e `kid` conhecido, `iss` ou `aud` errados, `exp` ausente ou vencido, e `nbf` ou `iat` no futuro. A sessão é conferida depois, no banco.
- Na reconexão, o servidor confere se o usuário e o aparelho continuam ativos. Se não, os comandos pendentes são aceitos ou recusados conforme a política da empresa e o banco local é apagado.
- Apagamento remoto do aparelho pelo administrador.
- Cofre (versão 0.7, `app/src/acesso/sessoes.ts`): o `expo-secure-store` guarda só itens pequenos, porque o SecureStore tem limite de cerca de 2 KB por valor: a credencial (access e refresh token da conta ativa), o snapshot da sessão para abrir sem rede, os nomes da última conta e o UUID do dispositivo, mais a lista de revogações que não saíram por falta de rede (`tenaz.revogacoes-pendentes`), só enquanto houver alguma. As permissões **não** vão para o cofre: ficam no banco local da conta. **Versão 0.8:** mais um item, o idioma escolhido (`tenaz.idioma`, alguns bytes), que não é segredo e vale para o aparelho, não para a conta (seção 2.3). Quando a API tiver `usuario.idioma` (ADR-013), o idioma do usuário passa a valer depois do login, e o item do cofre fica só para a tela de antes do login. Regra de plural de reserva, para aparelho sem `Intl.PluralRules` (`categoriaDeReserva`): o defeito do zero em `pt-BR` foi resolvido; a reserva é testada contra `Intl.PluralRules` nos dois idiomas (`app/src/i18n/i18n.test.ts`).
- Saída e sessão encerrada (versão 0.7): sair apaga as credenciais locais primeiro (funciona sem rede) e revoga a família no servidor, ou guarda a revogação para a próxima abertura; sessão encerrada pelo servidor (`invalid_grant`) apaga só as credenciais. **Nenhum dos dois apaga o banco local nem a fila de comandos**: comando de campo não se perde por sair da conta (princípio 3). O banco só é apagado pelo apagamento remoto e pela desativação reconhecida na reconexão (itens acima).

### 10.7 PowerSync ou sincronização própria

| Critério | PowerSync | Sincronização própria |
|---|---|---|
| Leitura (pull) | Pronta: replicação lógica do Postgres, regras de sincronização, SDK para React Native | A construir: cursor por snapshot, tombstones, escopo por usuário (seção 10.4) |
| Escrita (push) | Não resolve: as escritas vão para o backend da Tenaz por uma fila de upload | A construir, mas igual nos dois casos: outbox de comandos + casos de uso |
| Conflitos | Responsabilidade do backend | Responsabilidade do backend |
| Isolamento entre empresas | Lê pela replicação lógica, **fora do RLS**. O isolamento passa a depender das regras de sincronização do PowerSync: um segundo mecanismo a manter e testar. | O pull roda dentro de `NaEmpresa`, sob o mesmo RLS da API |
| Operação | Serviço adicional (nuvem do fornecedor ou auto-hospedado, com armazenamento próprio) e slot de replicação lógica no primário | Nenhum componente novo |
| Residência de dados | Nuvem do fornecedor em São Paulo a verificar; auto-hospedado resolve, com custo de operação | Dados ficam na infraestrutura da Tenaz |
| Licença e dependência | Edição aberta com licença própria (a verificar); dependência de fornecedor no componente mais crítico do app | Sem dependência |
| Tempo até o MVP | Menor para a leitura | Maior para a leitura (estimativa: 3 a 5 semanas de uma pessoa a mais) |

**Recomendação (proposta): sincronização própria.** A parte difícil do problema, que é escrita, idempotência e conflito por regra de domínio, fica com a Tenaz nos dois cenários. O PowerSync economiza só a leitura e, em troca, cria um segundo caminho de acesso aos dados fora do RLS, que é justamente o risco número 1 do produto (R2). Também acrescenta um serviço para operar e uma dúvida de residência de dados.

Para manter a porta aberta: o esquema SQLite local e a outbox de comandos seguem o formato que o PowerSync usa (tabelas locais espelhando as do servidor, fila de upload separada). Gatilho de revisão: se o pull incremental próprio passar de 2 incidentes de perda ou inconsistência por trimestre depois do MVP, ou se a convergência p95 passar de 60 s sem solução simples, reavaliar o PowerSync auto-hospedado. Ver [ADR-004](#adr-004--app-offline-first-com-sincronização-própria).

### 10.8 Um app para o técnico da empresa e o da subcontratada

O Tenaz Mobile é **o mesmo app** (mesmo binário, mesmas telas) para o técnico da empresa e para o técnico da subcontratada com conta Tenaz Serviços (sumário 2). O que muda é a conta e, com ela, os módulos habilitados (`GET /v1/me/recursos`): numa conta com SRV, a OS mostra cliente atendido, endereço de atendimento, serviço a realizar, componentes com número de série e assinatura do cliente. Não há app separado para subcontratadas. **(proposta)**

**Contas no aparelho.**

| Tema | Definição **(proposta)** |
|---|---|
| Quem tem mais de uma conta | Técnico que trabalha para duas subcontratadas, técnico de subcontratada que também é usuário `prestador` numa conta de cliente sem conexão, gestor de subcontratada que também é técnico. Cada acesso é um `usuario` numa `empresa`, com token próprio (regra 3 da seção 6.6) |
| Seletor de conta | Cabeçalho do app mostra sempre o nome da conta ativa. A troca de conta é uma ação explícita, com a contagem de tarefas novas e de pendências de sincronização de cada conta |
| Dados locais | **Um banco SQLite por conta** (`tenaz-<empresa>-<usuario>.db`, versão 0.7), com sua outbox de comandos, seu cursor de pull e as permissões da conta; a chave de cifragem por banco no Keychain ou Keystore segue proposta, com o SQLCipher (seção 2.3). Nenhuma tabela local mistura duas contas; a troca de conta fecha um banco e abre outro |
| Sincronização | **Versão 0.7:** só a conta ativa sincroniza, porque o aparelho tem uma só credencial viva por vez (linha "Login"). As outboxes das outras contas ficam guardadas no banco de cada uma e saem quando a conta volta a ser a ativa; o seletor mostra as pendências de cada conta. O servidor recusa comando cujo token não seja da empresa daquele banco. A versão 0.6 dizia que todas as contas sincronizavam em segundo plano, cada uma com o próprio token: substituído pela decisão de uma credencial por aparelho |
| Sair de uma conta | **Versão 0.7:** apaga só as credenciais e revoga a família no servidor; **não** apaga o banco local nem a fila (seção 10.6), que continuam lá para a próxima entrada na mesma conta. Desativação do usuário ou apagamento remoto atinge só o banco da conta que o emitiu. A versão 0.6 dizia que sair apagava o banco local: substituído, para não perder comando de campo |
| Login | Login único por e-mail (versão 0.6): a `identidade` é global e `acesso_empresa` lista as contas em que ela entra (documento 2, seção 5.1). **Versão 0.7 ("uma sessão e um banco local por conta"):** um banco por conta e **uma só credencial ativa no aparelho**. Trocar de empresa usa `POST /v1/auth/trocar-empresa`, que revoga no servidor a credencial anterior e devolve a da conta nova; por isso **exige conexão**. Sem rede, o técnico continua na conta ativa. Menos superfície num aparelho perdido (uma credencial viva, não uma por conta), ao custo de as outras contas não sincronizarem em segundo plano |

**Notificações push de tarefa nova.**

| Tema | Definição **(proposta)** |
|---|---|
| Quando | OS atribuída ao técnico (inclusive a que chega do cliente conectado e é distribuída pela subcontratada), OS reatribuída ou retirada dele, OS assumida por outro técnico da equipe, liberação recusada (volta a `em_execucao`), prazo de atendimento próximo |
| Caminho | O evento da OS (`ordem_servico.*`) gera um job no `tenaz worker`, que envia pelo Expo Push (APNs e FCM). O registro do aparelho é por (aparelho, `usuario`, `empresa`) e é apagado ao sair da conta ou desativar o usuário |
| Conteúdo | Título genérico e número: "Nova tarefa: OS-2026-000123", com o nome da conta quando o aparelho tem mais de uma. Sem nome de cliente pessoa física, endereço, descrição ou dados pessoais (seção 13.1) |
| Ao tocar | Abre a conta certa (troca de conta, se preciso) e a OS; se a OS ainda não estiver no aparelho, dispara um pull daquela conta |
| Garantia | O push é só um aviso. A fonte de verdade é o pull: aparelho offline recebe a tarefa na próxima sincronização, com ou sem push. Push perdido não perde tarefa |
| Agrupamento | Várias atribuições em sequência (distribuição em lote) viram uma notificação "5 tarefas novas" |
| Silêncio | O técnico pode silenciar notificações fora do horário de trabalho por conta |

---

## 11. Telemática e IoT

### 11.1 Por fase

| Fase | Entrega | Armazenamento |
|---|---|---|
| **MVP** | Leitura manual e importação por CSV. Modelo de leitura com origem (`manual`, `importacao`, `os`) e validação. | `leitura_medidor` |
| **Fase 2** | Conectores de telemática e cartão combustível com modelo canônico de evento. Regra configurável "DTC → solicitação ou OS". | Eventos brutos em tabela particionada por mês; leituras válidas em `leitura_medidor` (origem `telematica` ou `cartao_combustivel`) |
| **Fase 3** | Sensores via gateway parceiro, MQTT com Sparkplug B, regras de limiar e tendência, manutenção baseada em condição | Banco de séries temporais (11.5) |

### 11.2 Conectores (Fase 2)

Cada conector é um adaptador no módulo INT que traduz a API do fornecedor para o modelo canônico. O resto do sistema só conhece o modelo canônico.

Ordem sugerida, conforme o estudo de mercado: Cobli e Geotab; rFMS (Volvo e Scania); depois Ticket Log, Sascar e Omnilink, cujas APIs não foram verificadas.

| Aspecto | Definição **(proposta)** |
|---|---|
| Forma de coleta | Webhook de entrada quando o fornecedor oferece; senão, consulta periódica por job, respeitando o limite de cada API |
| Credenciais do cliente | Cifradas com chave por empresa (envelope encryption com KMS), nunca em texto aberto no banco nem em log |
| Casamento de veículo | Por VIN, depois por placa. Evento sem casamento vai para a fila "não identificados", visível ao gestor, em vez de ser descartado |
| Redução de volume | Odômetro e horímetro viram leitura de medidor no máximo uma vez por hora por veículo, ou na mudança de dia. O PRE só precisa do último valor confiável. |
| Deduplicação | Chave natural por evento do fornecedor (`fornecedor`, `id_externo`) |
| Retenção do bruto | 90 dias, configurável por plano |

### 11.3 Modelo canônico de evento

| Tipo | Campos principais |
|---|---|
| `leitura` | `ativo_id`, `medidor` (km, horas, horas de PTO, horas de refrigeração), `valor`, `medido_em`, `fonte` |
| `dtc` | `ativo_id`, `protocolo` (J1939, OBD-II), `spn`, `fmi`, `ocorrencias`, `severidade` (`parar_imediatamente`, `proxima_parada`, `seguro_continuar`, `desconhecida`), `descricao`, `ativo_em`, `resolvido_em` |
| `combustivel` | `ativo_id`, `litros`, `produto` (diesel, Arla, gasolina, etanol), `valor_centavos`, `posto`, `hodometro`, `condutor_ref` |
| `pneu` | `ativo_id`, `posicao`, `pressao`, `temperatura`, `sensor_id` |
| `posicao` | `ativo_id`, `latitude`, `longitude`, `velocidade`, `registrado_em` (dado pessoal indireto: ver LGPD em 9_SEGURANCA_CONFIABILIDADE.md) |
| `alerta` | `ativo_id`, `codigo`, `texto`, `severidade` |

Todo evento carrega `empresa_id`, `fornecedor`, `id_externo`, `recebido_em` e o payload original para auditoria.

### 11.4 Regras de evento para manutenção

- Regra configurável por empresa e classe de ativo: "DTC com severidade X → solicitação" ou "→ OS corretiva".
- Janela de deduplicação: o mesmo SPN/FMI no mesmo ativo não abre um segundo chamado enquanto houver um aberto.
- Toda solicitação gerada automaticamente mostra o evento de origem.

### 11.5 IoT (Fase 3)

Padrão de referência: sensor ou CLP → gateway de borda do parceiro → MQTT Sparkplug B (ou OPC UA convertido no gateway) → broker → ingestão → séries temporais → regras → solicitação ou OS. A Tenaz fica fora da rede de automação do cliente e nunca se conecta a CLP ou SCADA.

| Componente | Proposta |
|---|---|
| Broker | MQTT gerenciado na região São Paulo (ex.: AWS IoT Core), com certificado mTLS por gateway e tópicos por empresa |
| Gateway | De parceiro, com armazenamento local durante falha (store-and-forward). Um gateway próprio em Rust fica como opção futura, só se houver demanda. |
| Decodificação | `tenaz ingestor` decodifica Sparkplug B (protobuf) e grava em lotes |
| Séries temporais | **TimescaleDB** é a preferência do estudo, mas **não está disponível no RDS for PostgreSQL**. Opções: (a) particionamento nativo do Postgres com índices BRIN no próprio RDS, suficiente para volume moderado; (b) Timescale gerenciado, se houver oferta em São Paulo (a verificar); (c) ClickHouse via CDC, só se o volume justificar. Recomendação: começar com (a) e medir. |
| Regras | Limiar, tendência e, depois, Weibull por família de ativo, executados em jobs |
| Contador vindo de sensor | Vira `leitura_medidor` com origem `iot` (já prevista no sumário para a Fase 3). |

---

## 12. Inteligência artificial

O módulo IA entra na Fase 2 (nível 1) e evolui para agentes com aprovação na Fase 3. A governança detalhada está em 9_SEGURANCA_CONFIABILIDADE.md; aqui está a mecânica. Ver [ADR-007](#adr-007--ia-com-provedor-externo-rag-isolado-e-ferramentas-tipadas).

### 12.1 Casos de uso

| Nível | Caso | Grava algo? |
|---|---|---|
| 1 (Fase 2) | Perguntar sobre manuais e histórico de OS, com citação da fonte | Não |
| 1 | Rascunho de OS ou solicitação a partir de voz ou foto | Só rascunho; o usuário revisa e confirma |
| 1 | Gerar modelo de checklist ou procedimento a partir de PDF do fabricante | Só rascunho |
| 1 | Sugestão de código de falha ao encerrar OS | Sugestão; o técnico escolhe |
| 2 (Fase 3) | Agentes: triagem, deduplicação, sugestão de atribuição e de programação | Propõe ação; um humano aprova cada uma |
| 2 | Servidor MCP somente leitura | Não |

### 12.2 Provedor

- Acesso a LLM por API, atrás da interface interna `ia.Provedor`. Trocar de provedor não muda os módulos.
- Requisito contratual: **cláusula de não treinamento** com dados do cliente e retenção zero ou mínima.
- Preferência por execução no Brasil (ex.: modelo disponível no Amazon Bedrock em `sa-east-1` sem inferência entre regiões). Se o processamento ocorrer fora do país, ele é transferência internacional: entra no DPA e segue as cláusulas-padrão da Res. CD/ANPD 19/2024.
- Escolha do modelo por tarefa: modelo menor para classificar e sugerir código de falha; maior para rascunho e resposta com RAG.

### 12.3 RAG isolado por empresa

| Elemento | Definição **(proposta)** |
|---|---|
| Fontes | Manuais e documentos técnicos de ativos, procedimentos, histórico de OS (descrição, causa, solução) |
| Fora do índice | Dados sensíveis (CPF completo, dados de vítimas de sinistro, CPF de locatário pessoa física; sumário 3.6), dados de pessoas, credenciais, documentos marcados como restritos. A Tenaz não guarda dado de saúde |
| Armazenamento | Tabela `trecho_documento` com `empresa_id`, `documento_id`, `origem`, `texto`, `embedding vector(n)`, com RLS forçado como qualquer tabela |
| Índice | HNSW do pgvector. Usar a varredura iterativa do pgvector para filtros, para não retornar vazio quando o filtro por empresa descarta os vizinhos mais próximos. |
| Ingestão | Jobs River: extrair texto do PDF, dividir em trechos, gerar embedding, gravar. Reindexação quando o documento muda. |
| Permissão | A busca respeita a permissão do usuário: quem não vê um documento não recebe trechos dele |
| Citação | Toda resposta mostra as fontes (documento, página ou OS). Sem fonte, a resposta diz que não encontrou. |

### 12.4 Ferramentas tipadas

O modelo não gera SQL. Ele chama ferramentas com schema fixo, implementadas sobre a API pública dos módulos e executadas **com as permissões do usuário que perguntou**.

| Ferramenta | Tipo | Exemplo |
|---|---|---|
| `buscar_ativos` | Leitura | "Empilhadeiras da unidade Cajamar paradas" |
| `historico_os_ativo` | Leitura | "Últimas 10 corretivas do CMP-0347" |
| `indicadores_ativo` | Leitura | MTBF, MTTR e disponibilidade de um ativo |
| `buscar_documentos` | Leitura (RAG) | "Torque dos parafusos da roda no manual" |
| `rascunhar_os` | Rascunho | Cria um rascunho visível só ao usuário |
| `rascunhar_checklist` | Rascunho | Idem |

Regras:

1. Nenhuma ferramenta grava dado definitivo. A gravação acontece quando o usuário confirma na interface, que chama o endpoint comum da API, com validação e auditoria normais.
2. Na auditoria, o registro criado a partir de rascunho de IA fica com `origem = ia` e a referência da interação.
3. Conteúdo de manuais e OS é tratado como dado, não como instrução. Uma injeção de prompt num PDF não ganha mais poder que o usuário, porque as ferramentas já estão limitadas às permissões dele.

### 12.5 Governança e custo

| Controle | Mecanismo **(proposta)** |
|---|---|
| Liga e desliga | Por empresa (administrador) e por perfil ou usuário (permissão `ia.usar`). Indisponível para os perfis `prestador` e `locatario` e para o link de uso único |
| Fontes permitidas | O administrador escolhe se o histórico de OS e quais pastas de documentos entram no RAG |
| Registro | Pergunta, fontes, ferramentas chamadas, resposta, modelo e custo por interação, com `empresa_id`. Retenção de 90 dias, configurável; consulta e exportação pelo administrador. |
| Créditos | Razão (ledger) de créditos por empresa: o plano inclui uma cota mensal, pacotes adicionais são pré-pagos. Cada chamada estima custo por tokens e debita. |
| Limites | Teto por empresa e por usuário por dia. Sem créditos, a IA para e o restante do sistema segue normal. |
| Redução de custo | Cache de prompts do provedor, cache de embeddings, modelo menor quando a tarefa permite |
| Qualidade | Conjunto de avaliação (golden set) por caso de uso, rodado a cada mudança de prompt ou de modelo. Prompts versionados no repositório. |
| Isolamento de falha | Timeout e disjuntor. IA fora do ar não afeta OS, preventiva nem app. |

---

## 13. Infraestrutura

### 13.1 Nuvem e região

Proposta: **AWS `sa-east-1` (São Paulo)**, com 3 zonas de disponibilidade. Comparação no [ADR-006](#adr-006--nuvem-aws-na-região-são-paulo).

**Versão 0.7:** a nuvem gerenciada é a exceção explícita do [ADR-012](#adr-012--só-software-livre-ou-gratuito-e-seaweedfs-no-lugar-do-minio) à regra de só software gratuito ou de código aberto (decisão do dono do produto, 13/11/2026). Os serviços da tabela abaixo continuam proposta quanto à configuração; e-mail, push do iOS e lojas, pagos por natureza e ainda sem decisão, seguem a seção 16 (H5 e H6).

**Versão 0.8 (ADR-013, aceito, Parte D):** a mesma região e a mesma célula atendem Brasil, Argentina e Paraguai. Para empresas da Argentina e do Paraguai, a hospedagem em São Paulo é transferência internacional de dados, permitida com garantias contratuais por país (pendência P14 do sumário). Gatilho de célula por país: parecer da P14 que exija armazenamento local, ou cliente que exija residência no país.

| Camada | Serviço **(proposta)** | Configuração inicial |
|---|---|---|
| DNS | Route 53 | Registros com health check |
| CDN e WAF | CloudFront + AWS WAF | Regras gerenciadas, limite por IP, bloqueio geográfico opcional por empresa |
| Web | S3 (privado) atrás do CloudFront | Versões imutáveis |
| Balanceador | Application Load Balancer | TLS 1.2+ com política moderna, redirecionamento HTTP → HTTPS |
| Contêineres | **ECS com Fargate**, arquitetura ARM64 (Graviton) | API: 2 a N tarefas; workers: 2 a N; ingestor: Fase 2. Distribuídas em 3 zonas. Escala por CPU e por profundidade da fila. |
| Banco | **RDS for PostgreSQL Multi-AZ** | Standby síncrono em outra zona; uma réplica de leitura; backups automáticos com PITR de 35 dias em produção |
| Cache | ElastiCache (Valkey) | Nó pequeno com réplica em outra zona |
| Arquivos | S3 com versionamento, SSE-KMS e bloqueio de acesso público | Prefixo por `empresa_id`; ciclo de vida para classes mais baratas |
| Backups | AWS Backup com cópia para conta separada e cofre com bloqueio (Vault Lock) | Detalhes em 9_SEGURANCA_CONFIABILIDADE.md. O local da cópia de recuperação de desastre regional é a pendência P7 (seção 13.5) |
| Segredos e chaves | Secrets Manager + KMS com chaves gerenciadas pelo cliente (CMK) | Rotação automática |
| E-mail | SES na região | DKIM, SPF e DMARC |
| Push | Expo Push → APNs e FCM | Sem dado pessoal no conteúdo da notificação |
| Rede | VPC com sub-redes privadas para tarefas e banco; VPC endpoints para S3, ECR, Secrets Manager e KMS | Banco sem acesso público |

**Por que ECS Fargate e não Kubernetes:** para um monólito com três serviços, o EKS acrescenta um plano de controle, atualizações de cluster, operadores e uma camada de rede para manter, sem benefício no estágio atual. O ECS Fargate não tem servidor para corrigir, escala por serviço e integra com IAM, ALB e CloudWatch. Gatilho de revisão: necessidade de rodar a Tenaz em outra nuvem ou on-premises para um cliente, ou mais de 15 serviços.

**Por que RDS e não Aurora no início:** custo menor na escala inicial. O Aurora PostgreSQL entra se o failover do Multi-AZ (tipicamente 1 a 2 minutos) passar a ameaçar o SLO, ou se o volume pedir várias réplicas. A troca é suportada pela própria AWS e não muda o código.

**Pooler de conexões:** sem RDS Proxy nem PgBouncer no MVP. Cada tarefa mantém um pool `pgx` dimensionado por papel de banco: a API mantém o pool da aplicação e o pool pequeno da autenticação (versão 0.6, seção 6.9); o sincronizador, o seu. Se adotado depois, o pooler precisa funcionar em modo de transação com `set_config(..., true)` e com a fila River (a verificar no momento da adoção).

### 13.2 Infraestrutura como código

- **OpenTofu** (proposta), compatível com Terraform e com licença aberta (MPL). Estado remoto em S3 com travamento.
- Módulos próprios em `infra/tofu/modulos/`, ambientes em `infra/tofu/ambientes/`.
- Toda mudança de infraestrutura passa por PR com `tofu plan` anexado e scan de IaC no CI.
- Nada é criado manualmente no console de produção. Mudança emergencial manual vira PR em até 1 dia útil.

### 13.3 Contas e ambientes

Contas AWS separadas sob AWS Organizations. **(proposta)**

| Conta | Uso |
|---|---|
| Gerência | Organizations, faturamento, SSO da equipe. Sem carga de trabalho. |
| Segurança e logs | CloudTrail organizacional, GuardDuty, Security Hub, logs centralizados, somente leitura para a equipe |
| Backup | Cofres de backup com bloqueio. Ninguém da operação diária tem acesso de escrita. |
| `dev` | Desenvolvimento e ambientes efêmeros de PR |
| `hml` | Homologação, igual à produção em topologia, menor em tamanho |
| `prd` | Produção |

| Ambiente | Dados | Deploy | Observação |
|---|---|---|---|
| Local | Sementes fictícias | Manual | `docker compose` |
| Efêmero de PR | Sementes fictícias | Automático por PR | Destruído no merge; usado no E2E |
| `dev` | Sementes fictícias | Automático no merge | Pode quebrar |
| `hml` | Sintéticos; **nunca dados reais de cliente** | Automático no merge, após `dev` | DAST, carga, ensaio de restauração e de failover |
| `prd` | Reais | Aprovação manual | Acesso humano só por papel temporário, com registro |

### 13.4 Planos, entitlements e feature flags

Duas coisas diferentes, guardadas no banco:

| Conceito | Para quê | Exemplo |
|---|---|---|
| **Entitlement** | O que o plano contratado permite | `api_publica`, `webhooks`, `multi_unidade`, `sso_saml`, `nr13`, `estoque_terceiros`, `ordem_reparo`, `contrato_prestacao`, `locacao`, `srv`, `conexao_prestador`, `ia_creditos_mes = 500`, `rate_limit_rpm = 120` |
| **Feature flag** | Liberação gradual de código novo | `nova_tela_programacao` ligada para 10% das empresas |

Tabelas **(proposta)**: `plano`, `plano_recurso` (recurso, valor), `empresa_assinatura` (plano, vigência, quantidades contratadas), `empresa_recurso_ajuste` (exceções negociadas, com validade), `feature_flag` e `feature_flag_empresa`.

- O backend é a autoridade: `entitlements.Exigir(ctx, "api_publica")` nos handlers e casos de uso. Retorno 403 com código `RECURSO_FORA_DO_PLANO`.
- Web e app leem `GET /v1/me/recursos` para esconder ou mostrar funções, mas nunca são a barreira.
- Quantidades (técnicos, veículos ativos) são conferidas na criação ou ativação, com aviso antes do limite. Técnico cobrado = usuário ativo com perfil `tecnico`, mais usuário `prestador` que executa OS; link de uso único não conta. Veículo cobrado = ativo `veiculo` não `inativo`. Vale o maior número diário do ciclo (sumário 2.1).
- Plano Serviços (sumário 2.2): a conta da subcontratada é uma `empresa` com o entitlement `srv`; técnico cobrado = usuário ativo com perfil `tecnico` nessa conta; `gestor` e `planejador` não pagam. A conta pode somar um plano Planta ou Frota se tiver ativos próprios. O técnico da subcontratada conectada, que só tem usuário na conta dela, nunca entra na contagem do cliente. O desconto de canal por conexão ativa com cliente pagante depende da P4 e, se aprovado, é calculado pela cobrança a partir de `plataforma.conexao_empresa`, sem ler dados de OS das contas. **(proposta)**
- `conexao_prestador` é o lado do cliente: convidar e manter subcontratadas conectadas. Proposta: disponível em todos os planos, como o TER básico; validar com P4.
- O que cada plano inclui vem só da tabela de recursos do sumário (seção 2.1). Prestadores com catálogo, perfil `prestador` e link de uso único estão em todos os planos; estoque em terceiros, ordem de reparo, contratos de prestação e locação, no Profissional e no Enterprise.
- Cache de 60 segundos no Valkey, invalidado na mudança de plano.
- Avaliação de flags pela interface OpenFeature, para permitir trocar o armazenamento por um serviço externo no futuro sem mudar o código.
- Configuração por segmento (versão 0.5, sumário 2.1): campos extras na OS, regras configuráveis, fase de ciclo de vida e pacotes de segmento estão em todos os planos e não têm entitlement de liberação; o custo é limitado pelos limites de RF-ADM-016 (50 regras ativas) e RF-OS-021 (30 campos ativos por combinação), guardados como recursos de quantidade (`regras_ativas_max`, `campos_extras_max`) para permitir ajuste negociado. PRJ e CTR, no Profissional e no Enterprise, usam os entitlements `projetos` e `contratos_cliente`. Na aplicação de um pacote, item de módulo sem entitlement é "ignorado" (seção 5.6). **(proposta)**

### 13.5 Recuperação de desastre regional: pendência P7

A AWS tem uma única região no Brasil. Falha de zona é coberta em São Paulo (Multi-AZ, 3 zonas, RTO de 1 hora). Para falha da **região inteira**, o local da cópia de recuperação é a **pendência P7 do sumário, decisão do dono do produto com apoio jurídico**. Este documento apresenta as opções e não decide.

| Critério | (a) Cópia de DR fora do Brasil, declarada no DPA | (b) Segundo provedor no Brasil | (c) DR só multi-AZ em São Paulo |
|---|---|---|---|
| Como funciona | Backups, logs do banco, buckets, imagens e segredos replicados para outra região AWS (ex.: `us-east-1`), em modelo "pilot light"; infraestrutura recriada por IaC no desastre | Cópia lógica do banco (replicação lógica ou dumps frequentes) e dos arquivos para um provedor com região no Brasil (ex.: Oracle Cloud em São Paulo ou Vinhedo, Magalu Cloud); infraestrutura equivalente mantida por IaC própria desse provedor | Tudo em `sa-east-1`: Multi-AZ, PITR e cópia imutável em outra conta na mesma região. Em falha regional, aguarda-se a volta da região |
| Promessa "dados no Brasil" | Dados primários no Brasil; **cópias cifradas fora do país** e processamento fora só durante o desastre. Transferência internacional declarada no DPA, com as cláusulas-padrão da ANPD | Mantida integralmente | Mantida integralmente |
| RPO de região | ≤ 5 minutos (meta; a medir no ensaio) | Maior: depende do método de cópia (minutos com replicação lógica; horas com dumps) | O da última cópia que sobreviver; sem garantia se a região inteira for perdida |
| RTO de região | ≤ 4 horas (meta do sumário) | Maior que 4 horas: stack diferente, sem serviços equivalentes um a um, ensaio mais caro | Indeterminado: depende do provedor; o RTO de região publicado teria de ser bem maior ou "sem compromisso" |
| Custo e operação | Menor: mesmas ferramentas, replicação gerenciada | Maior: segundo provedor, segunda IaC, segunda equipe de conhecimento, cópia lógica a monitorar | Menor de todos |
| Impacto no código | Nenhum | Nenhum no domínio; adaptações de armazenamento, fila e segredos para o segundo provedor | Nenhum |

Combinações também são possíveis (por exemplo, (c) como padrão e (b) como opção Enterprise). Enquanto a P7 não for decidida, o RTO de região de 4 horas não deve ser prometido em contrato. As camadas de backup e o runbook regional estão em 9_SEGURANCA_CONFIABILIDADE.md, seção 3.

---

## 14. Observabilidade

### 14.1 Instrumentação

- **OpenTelemetry** no backend (HTTP, pgx, River, chamadas externas) e no web (erros e Web Vitals). No app, relatório de falhas e métricas de sincronização.
- Coletor OpenTelemetry (ADOT) como sidecar ou serviço no ECS.
- Destino **(proposta)**: CloudWatch (logs e métricas) e X-Ray (traces) em `sa-east-1` no MVP, para manter a telemetria no Brasil e com pouca operação. Como a instrumentação é OpenTelemetry, trocar para Grafana Cloud, Honeycomb ou Datadog não exige mudar o código; a residência dos dados do destino deve ser avaliada antes.

### 14.2 `empresa_id` em todo sinal, com cuidado de cardinalidade

| Sinal | Uso de `empresa_id` |
|---|---|
| Traces | Atributo em todo span raiz e de banco |
| Logs | Campo obrigatório em toda linha emitida dentro de requisição ou job |
| Métricas | **Não** como rótulo em todas as métricas (milhares de empresas explodiriam a cardinalidade). Só em métricas de negócio agregadas por empresa (uso de créditos de IA, comandos de sync pendentes) e via exemplars ligando a métrica ao trace. |

Responder "a empresa X está lenta?" é feito por consulta a logs e traces filtrados por `empresa_id`.

### 14.3 Logs estruturados

JSON via `slog`. Campos obrigatórios: `ts`, `nivel`, `msg`, `servico`, `versao` (SemVer do release, seção 9.5), `trace_id`, `span_id`, `empresa_id`, `usuario_id`, `request_id`, `rota`, `status`, `duracao_ms`.

Nunca registrar: senha, token, chave de API, segredo de conector, conteúdo de documento sensível, CPF completo. Um filtro de redação no logger mascara campos com nomes conhecidos, e um teste confere. **Versão 0.9 (ADR-013, Parte A, item 3; conferido em `obs.go` e `obs_test.go`):** o filtro redige também `identificacao`, que leva a identificação de qualquer país (CPF, CUIL, DNI, CI e as de empresa).

Retenção **(proposta)**: 30 dias quente, 13 meses arquivado, exceto logs de segurança (ver 9_SEGURANCA_CONFIABILIDADE.md).

### 14.4 Métricas principais

| Área | Métricas |
|---|---|
| API | Taxa, erros e duração por rota (RED); 429 por plano |
| Banco | Conexões em uso, latência de consulta, lag da réplica, deadlocks, crescimento de tabelas da fila e do outbox |
| Jobs | Profundidade da fila por tipo, idade do job mais antigo, falhas e reprocessamentos |
| Sync | Comandos recebidos, rejeitados por motivo, tempo de convergência, idade do comando pendente mais antigo reportado pelos aparelhos |
| Webhooks | Taxa de sucesso, tentativas, itens na fila de mortos (DLQ) |
| Conexões entre contas | Atraso da mensagem (origem até destino), mensagens pendentes, idade da pendente mais antiga, erros por motivo, divergências da reconciliação (seção 6.8) |
| Push | Enviados, falhas por provedor, tokens inválidos removidos |
| Integrações | Erros e latência por conector, eventos não casados |
| IA | Chamadas, tokens, custo, latência, respostas sem fonte |
| Negócio | OS preventivas geradas no prazo, alertas de vencimento enviados |

### 14.5 SLOs

| SLO **(proposta)** | Indicador | Meta | Janela |
|---|---|---|---|
| Disponibilidade da API (Essencial e Profissional) | Requisições válidas sem 5xx nem timeout / total | 99,95% (SLA de 99,9%) | 30 dias |
| Disponibilidade da API (Enterprise) | Idem, para as empresas Enterprise | 99,97% (SLA de 99,95%) | 30 dias |
| Latência de leitura | Leituras abaixo de 300 ms | 95% | 30 dias |
| Latência de escrita | Escritas abaixo de 500 ms | 95% | 30 dias |
| Painéis | Carregamento de painel de IND abaixo de 2 s | 95% | 30 dias |
| Convergência do app | Lotes de sync aplicados em menos de 60 s após a reconexão | 99% | 30 dias |
| Preventiva no prazo | OS preventivas geradas até 15 minutos após o vencimento do gatilho | 99,9% | 30 dias |
| Webhooks | Entregues (primeira tentativa ou retry) em até 5 minutos | 99% | 30 dias |
| Sincronização entre contas | Atraso entre o commit na conta de origem e o commit na conta de destino (seção 6.8) | p95 abaixo de 60 s | 30 dias |

Orçamento de erro em 30 dias: 99,97% são cerca de 13 minutos; 99,95% cerca de 21,6 minutos; 99,9% cerca de 43,2 minutos. As metas por plano vêm do sumário (seção 5). No modelo pool, as empresas Enterprise compartilham a infraestrutura das demais; na prática, enquanto houver contrato Enterprise sem banco dedicado, a plataforma inteira opera com o alvo de 99,97%. Com banco dedicado (Fase 3), a medição passa a ser separada.

### 14.6 Alertas por burn rate

Alertas sobre o consumo do orçamento de erro, com duas janelas para evitar falso positivo:

| Severidade | Taxa de consumo | Janela longa | Janela curta | Ação |
|---|---|---|---|---|
| Página (acorda alguém) | 14,4× | 1 hora | 5 minutos | Plantão responde em até 15 minutos |
| Página | 6× | 6 horas | 30 minutos | Plantão responde em até 15 minutos |
| Ticket | 1× | 3 dias | 6 horas | Tratar no próximo dia útil |

Além dos SLOs, há alertas diretos para sinais que não podem esperar: falha do teste de vazamento em produção (verificação sintética), backup não concluído, réplica com lag acima de 30 s, fila com job mais antigo acima de 15 minutos, certificado TLS vencendo em 14 dias. Na conexão entre contas: atraso p95 acima de 60 s por 15 minutos, mensagem pendente há mais de 5 minutos com a conexão `ativa`, vínculo em `erro` há mais de 30 minutos (ticket) e, como **alerta de segurança imediato**, qualquer recusa por ponta inválida, campo fora da lista ou campo de outro dono.

Verificações sintéticas externas (login, listagem de OS, sync) rodam a cada minuto de pelo menos 3 locais e alimentam a página de status.

---

## 15. Decisões de arquitetura (ADRs)

Formato: contexto, decisão, alternativas consideradas e consequências. Status possíveis: **Aceita** (decidida pelo dono do produto) ou **Proposta** (aguarda validação).

| ADR | Título | Status |
|---|---|---|
| [ADR-001](#adr-001--stack-go-react-react-native-e-postgresql) | Stack: Go, React, React Native e PostgreSQL | Aceita |
| [ADR-002](#adr-002--multiempresa-em-pool-com-rls) | Multiempresa em pool com RLS | Aceita |
| [ADR-003](#adr-003--monólito-modular) | Monólito modular | Aceita |
| [ADR-004](#adr-004--app-offline-first-com-sincronização-própria) | App offline-first com sincronização própria | Proposta |
| [ADR-005](#adr-005--fila-de-jobs-no-postgresql) | Fila de jobs no PostgreSQL | Proposta |
| [ADR-006](#adr-006--nuvem-aws-na-região-são-paulo) | Nuvem AWS na região São Paulo | Proposta; provedor e região aceitos como exceção do ADR-012 |
| [ADR-007](#adr-007--ia-com-provedor-externo-rag-isolado-e-ferramentas-tipadas) | IA com provedor externo, RAG isolado e ferramentas tipadas | Proposta |
| [ADR-008](#adr-008--acesso-externo-por-escopo-restritivo-e-link-de-uso-único) | Acesso externo por escopo restritivo e link de uso único | Proposta |
| [ADR-009](#adr-009--conexão-entre-contas) | Conexão entre contas | Proposta |
| [ADR-010](#adr-010--configuração-por-segmento-com-catálogos-fechados) | Configuração por segmento com catálogos fechados | Proposta |
| [ADR-011](#adr-011--papel-de-banco-próprio-para-as-credenciais-de-autenticação) | Papel de banco próprio para as credenciais de autenticação | Proposta |
| [ADR-012](#adr-012--só-software-livre-ou-gratuito-e-seaweedfs-no-lugar-do-minio) | Só software livre ou gratuito, e SeaweedFS no lugar do MinIO | Regra: Aceita (dono do produto). Substituto do MinIO: Proposta |
| [ADR-013](#adr-013--país-e-idioma-perfil-de-país-no-núcleo-textos-por-catálogo-e-conteúdo-legal-por-país-e-jurisdição) | País e idioma: perfil de país no núcleo, textos por catálogo e conteúdo legal por país e jurisdição | Aceita (dono do produto, 13/11/2026); complementos técnicos da Parte F, representação do dinheiro, e da Parte G, formatação para pessoas (versão 0.9) |

---

### ADR-001 — Stack: Go, React, React Native e PostgreSQL

**Status:** Aceita · **Data:** 12/11/2026

**Contexto.** A Tenaz é um SaaS multiempresa com web, app de campo offline, API pública e, a partir da Fase 2, ingestão de telemetria. O custo de infraestrutura por cliente precisa ser baixo para sustentar preço em reais e plano mensal. O isolamento entre empresas é o risco número 1. O dono do produto já opera um sistema em Go (PCP-Lev).

**Decisão.**

- Backend em **Go**, monólito modular, com **sqlc + pgx** e SQL explícito.
- Web em **React + TypeScript**.
- App em **React Native (Expo) + TypeScript**, offline-first com **SQLite**.
- **PostgreSQL** com RLS como banco único no início.
- **Docker** para empacotar e rodar.
- **Contrato único:** OpenAPI 3.1 gerado e validado no Go; tipos TypeScript gerados a partir dele.

Motivos:

1. Experiência do dono com Go no PCP-Lev: menos risco de execução.
2. Binário estático e baixo consumo de memória: mais clientes por tarefa, menor custo por cliente, imagem mínima com pouca superfície de ataque.
3. Concorrência nativa (goroutines) para ingestão de telemetria, conectores e workers, sem runtime adicional.
4. SQL explícito com sqlc + pgx: cada consulta é visível, revisável e testável quanto ao filtro por `empresa_id`; controle direto de `set_config` na transação.
5. Promessa de compatibilidade do Go 1: código de hoje compila daqui a anos com poucas mudanças.
6. TypeScript em web e app: um ecossistema só no front e regras compartilhadas em `pacotes/dominio`.

**Alternativas consideradas.**

| Alternativa | A favor | Contra | Resultado |
|---|---|---|---|
| TypeScript/Node no backend | Uma linguagem só em toda a stack; usada pela MaintainX | Maior consumo de memória por instância; ORMs populares escondem o SQL e dificultam garantir o filtro por empresa; concorrência para ingestão menos natural | Descartada |
| Java/Kotlin com Spring | Mercado grande de profissionais no Brasil; ecossistema maduro | JVM pesada para o custo por cliente pretendido; JPA/Hibernate atrapalha o controle fino de RLS e de `set_config` por transação; inicialização lenta | Descartada |
| Ruby on Rails | Produtividade alta; usado pelo Fleetio com sucesso | Sem experiência do dono; desempenho e custo de memória piores para ingestão; ActiveRecord pelo mesmo motivo do ORM | Descartada |
| Python | Melhor ecossistema de ciência de dados | Desempenho e tipagem fracos para o núcleo transacional | Descartada para o núcleo. **Reservada** para um serviço futuro de modelos de falha (Weibull, detecção de anomalia), Fase 3 |
| Rust | Desempenho e segurança de memória máximos | Curva de aprendizado e velocidade de entrega piores para um CRUD de domínio rico | Descartada para o núcleo. **Reservada** para um gateway de borda próprio, se houver demanda, Fase 3 |

**Consequências.**

- Positivas: custo de infraestrutura baixo, binário único, controle total do SQL, tipos ponta a ponta via OpenAPI.
- Negativas: duas linguagens (Go e TypeScript) e regras duplicadas entre backend e app offline, mitigadas pelos casos de teste compartilhados (seção 7). Contratação de Go no Brasil é mais difícil que de Java ou Node.
- Gatilho de revisão: nenhum previsto para o núcleo.

---

### ADR-002 — Multiempresa em pool com RLS

**Status:** Aceita · **Data:** 12/11/2026

**Contexto.** Muitos clientes pequenos e médios no mesmo sistema. Vazamento entre empresas é o risco R2, classificado como crítico. A AWS descreve três modelos: pool (tabelas compartilhadas), bridge (esquema por cliente) e silo (instância por cliente), e considera o RLS necessário no pool.

**Decisão.** Modelo **pool**: tabelas compartilhadas com `empresa_id` em toda tabela de cliente, **RLS habilitado e forçado**, papel de aplicação sem `BYPASSRLS`, contexto definido por `set_config(..., true)` dentro da transação, filtro explícito no repositório, chaves estrangeiras compostas e teste de vazamento em todo build (seções 6 e 9.3). Banco dedicado (silo) como opção do plano Enterprise na Fase 3, com o mesmo esquema.

**Alternativas consideradas.**

| Alternativa | A favor | Contra | Resultado |
|---|---|---|---|
| Só filtro na aplicação, sem RLS | Simples | Um único `WHERE` esquecido vaza dados | Descartada |
| Esquema por empresa (bridge) | Separação visível | Migrações multiplicadas por milhares de esquemas; não isola recursos (mesmos problemas de vizinho barulhento do pool, segundo a AWS) | Descartada |
| Banco por empresa (silo) para todos | Isolamento máximo | Custo e operação inviáveis para clientes pequenos | Descartada como padrão; mantida como opção Enterprise |

**Consequências.**

- Positivas: uma migração por release, custo baixo por empresa, defesa em profundidade.
- Negativas: armadilhas conhecidas do RLS (dono da tabela ignora política sem `FORCE`, `SET` de sessão vaza no pool, desempenho sem índice em `empresa_id`), todas tratadas por regra e teste. Vizinho barulhento exige limites por plano e monitoramento por empresa.
- Gatilho de revisão: empresa cujo volume degrade as demais de forma recorrente → migrar essa empresa para silo.

---

### ADR-003 — Monólito modular

**Status:** Aceita · **Data:** 12/11/2026

**Contexto.** Equipe pequena, domínio rico com muitas regras que cruzam módulos (OS consulta CON e baixa EST na mesma operação). Fleetio e MaintainX chegaram à escala com monólitos.

**Decisão.** Um único binário Go com módulos isolados por pacote (seção 5): API pública por módulo, `internal/` imposto pelo compilador, regras de importação verificadas no CI, chamadas síncronas na mesma transação e eventos via outbox transacional para reações assíncronas.

**Alternativas consideradas.**

| Alternativa | A favor | Contra | Resultado |
|---|---|---|---|
| Microsserviços desde o início | Escala e deploy independentes | Transações distribuídas, sagas, rede, observabilidade e custo muito acima do necessário; isolamento por empresa repetido em cada serviço | Descartada |
| Monólito sem fronteiras | Mais rápido no primeiro mês | Acoplamento cresce e impede extrair partes depois | Descartada |

**Consequências.**

- Positivas: transações locais, um deploy, depuração simples, refatoração barata.
- Negativas: disciplina de fronteiras depende de lint e revisão; um bug grave pode afetar todos os módulos (mitigado por modos de execução separados: API, worker, ingestor).
- Gatilho de revisão: um módulo com perfil de carga muito diferente (ex.: ingestão IoT) ou equipe acima de 25 pessoas → extrair esse módulo como serviço, mantendo a API pública dele.

---

### ADR-004 — App offline-first com sincronização própria

**Status:** Proposta · **Data:** 12/11/2026

**Contexto.** Técnicos e motoristas trabalham em áreas sem sinal. Meta: 7 dias offline, sem perda, convergência em menos de 60 s. A queixa recorrente contra os incumbentes brasileiros é a sincronização frágil. O Realm/Atlas Device Sync chegou ao fim de vida em 30/09/2025.

**Decisão.** SQLite local como fonte de leitura do app; **outbox de comandos de domínio idempotentes** para escrita; **pull incremental** por cursor de snapshot do Postgres, executado sob RLS; conflitos resolvidos por regras de domínio explícitas (seção 10). Sincronização construída pela Tenaz, com esquema local compatível com o formato do PowerSync.

**Alternativas consideradas.**

| Alternativa | A favor | Contra | Resultado |
|---|---|---|---|
| PowerSync | Pull pronto; SDK para React Native | Escrita e conflitos continuam com a Tenaz; lê fora do RLS (segundo mecanismo de isolamento); serviço adicional; residência a verificar | Plano B. **Versão 0.7 (ADR-012):** suspenso enquanto a edição autohospedada estiver sob licença "source-available" (FSL), que a regra de admissão não aceita; conferir a licença vigente antes de reabrir |
| ElectricSQL | Sincronização de leitura por "shapes" sobre HTTP | Mesmas objeções de isolamento e operação; escrita também fica com o backend | Descartada |
| WatermelonDB com sync próprio | Banco local reativo maduro | Protocolo de sync ainda a construir; menos alinhado ao SQLite puro | Descartada |
| Realm / Atlas Device Sync | Era a referência | Fim de vida em 2025 | Descartada |
| App online com cache | Simples | Não atende à meta de 7 dias nem à criação de OS offline | Descartada |
| CRDT genérico | Mesclagem automática | Regras de manutenção não são mescláveis (OS cancelada não pode ser concluída); complexidade alta | Descartada |

**Consequências.**

- Positivas: um único caminho de escrita (casos de uso da API), um único mecanismo de isolamento (RLS), dados na infraestrutura própria.
- Negativas: 3 a 5 semanas a mais de desenvolvimento para o pull; responsabilidade integral por bugs de sincronização.
- Gatilho de revisão: mais de 2 incidentes de perda ou inconsistência de sync por trimestre após o MVP, ou convergência p95 acima de 60 s sem correção simples.

---

### ADR-005 — Fila de jobs no PostgreSQL

**Status:** Proposta · **Data:** 12/11/2026

**Contexto.** O sistema precisa de jobs (preventiva, alertas, webhooks, importações, indexação de IA, exportações) e de eventos internos confiáveis. Escrever no banco e publicar numa fila externa em passos separados causa inconsistência ("gravou mas não avisou").

**Decisão.** Fila em tabela do PostgreSQL com `SKIP LOCKED`, usando **River** (biblioteca Go sobre pgx), com inserção do job na mesma transação da mudança de negócio. O outbox de eventos de domínio usa o mesmo mecanismo.

**Alternativas consideradas.**

| Alternativa | A favor | Contra | Resultado |
|---|---|---|---|
| Fila própria com `SKIP LOCKED` | Sem dependência | Reimplementar retry, agendamento, unicidade, painel | Descartada |
| Amazon SQS | Gerenciada, escala alta | Não participa da transação: exige outbox separado mesmo assim; mais um serviço | Futuro, para cargas específicas |
| Redis (Asynq) | Rápida | Outra fonte de verdade; não transacional com o Postgres | Descartada |
| NATS JetStream, RabbitMQ, Kafka | Escala e fan-out | Operação pesada para o estágio | Futuro, se os gatilhos ocorrerem |

**Consequências.**

- Positivas: consistência transacional, nenhuma infraestrutura nova, jobs visíveis por SQL.
- Negativas: carga e bloat no primário; exige vacuum bem configurado e limpeza de jobs concluídos.
- **Quando migrar** (qualquer um dos gatilhos): mais de 500 jobs por segundo sustentados; tabela de jobs afetando a latência do primário (p95 de escrita acima da meta por causa da fila); necessidade de consumidores fora do monólito; ingestão IoT da Fase 3 em volume alto. A migração é parcial: o tipo de carga que disparou o gatilho vai para SQS ou NATS JetStream, alimentado pelo outbox; o restante fica no River.

---

### ADR-006 — Nuvem AWS na região São Paulo

**Status:** Proposta; provedor e região **aceitos** pelo dono do produto em 13/11/2026 · **Data:** 12/11/2026

**Nota de 13/11/2026 ([ADR-012](#adr-012--só-software-livre-ou-gratuito-e-seaweedfs-no-lugar-do-minio)).** O dono do produto fixou a regra de só usar software 100% gratuito ou de código aberto no que o sistema usa diretamente e, no mesmo dia, abriu a **exceção explícita** para esta decisão: a produção roda em nuvem gerenciada (AWS, região São Paulo), com banco, chaves e armazenamento gerenciados. A suspensão registrada mais cedo, na mesma versão 0.7, foi retirada. Continuam proposta os serviços e a configuração (seção 13), e continua aberta a P7.

**Contexto.** Promessa comercial: dados primários hospedados no Brasil, RPO e RTO publicados. Metas do sumário: SLA de 99,9% (Essencial e Profissional) e 99,95% (Enterprise), RPO ≤ 5 min, RTO ≤ 1 h (zona) e ≤ 4 h (região, dependente da P7). Equipe pequena: o máximo de serviços gerenciados.

**Decisão.** AWS `sa-east-1` como região primária, com ECS Fargate, RDS for PostgreSQL Multi-AZ, S3, ElastiCache, KMS, Secrets Manager, AWS Backup com cópia para outra conta, e infraestrutura em OpenTofu. **O local da cópia de recuperação de desastre regional não é decidido aqui:** é a pendência P7, com as opções da seção 13.5 (cópia fora do Brasil declarada no DPA, segundo provedor no Brasil, ou DR só multi-AZ em São Paulo com RTO de região maior).

**Alternativas consideradas.**

| Alternativa | A favor | Contra | Resultado |
|---|---|---|---|
| Google Cloud `southamerica-east1` (São Paulo) | Cloud SQL e Cloud Run simples | Menos serviços gerenciados usados pelo mercado-alvo; menos profissionais | Alternativa viável |
| Azure Brazil South | Forte em clientes Microsoft; Entra ID | Postgres gerenciado com histórico de mudanças de produto; custo semelhante | Alternativa viável |
| Oracle Cloud (São Paulo, Vinhedo) | Duas regiões no Brasil, preço agressivo | Ecossistema menor; menos experiência de mercado | Candidata a segundo provedor na opção (b) da P7 |
| Magalu Cloud | Empresa brasileira | Catálogo de serviços gerenciados ainda limitado | Acompanhar |
| VPS com Docker Compose (modelo PCP-Lev) | Custo mínimo | Sem alta disponibilidade, sem PITR gerenciado; não atinge o SLA | Descartada |

**Consequências.**

- Positivas: 3 zonas em São Paulo, Multi-AZ síncrono (RPO zero em falha de zona), PITR com envio de logs a cada 5 minutos, ferramentas maduras de backup imutável.
- Negativas: `sa-east-1` é uma das regiões mais caras da AWS. A AWS só tem **uma região no Brasil**: recuperar de falha de região exige cópia fora do país ou um segundo provedor no Brasil; sem nenhum dos dois, o RTO de região fica maior (pendência P7, seção 13.5).
- Gatilho de revisão: decisão da P7 pelo dono do produto; cliente ou regulação exigindo que nem o backup saia do Brasil.

---

### ADR-007 — IA com provedor externo, RAG isolado e ferramentas tipadas

**Status:** Proposta · **Data:** 12/11/2026

**Contexto.** O mercado já usa IA generativa em CMMS (busca em manuais, rascunho de OS). Os líderes corporativos limitam a IA à leitura ou exigem controle por usuário. Os riscos são vazamento entre empresas via índice vetorial, alucinação, custo descontrolado e uso de dados do cliente para treinar modelos.

**Decisão.**

- LLM de provedor externo atrás da interface `ia.Provedor`, com cláusula de não treinamento e preferência por execução no Brasil.
- RAG com **pgvector no mesmo PostgreSQL**, com `empresa_id` e RLS forçado como qualquer tabela.
- Ferramentas tipadas sobre a API pública dos módulos, executadas com as permissões do usuário. Sem text-to-SQL.
- **Nenhuma gravação sem confirmação humana.**
- Governança por empresa e por usuário; custo controlado por créditos.

**Alternativas consideradas.**

| Alternativa | A favor | Contra | Resultado |
|---|---|---|---|
| Banco vetorial separado (Pinecone, Qdrant) | Escala vetorial maior | Segundo mecanismo de isolamento; dado fora do Postgres e possivelmente fora do Brasil | Descartada no início |
| Modelo próprio auto-hospedado (Llama, Qwen) | Dados não saem da infraestrutura | GPU cara e ociosa para o volume inicial; qualidade inferior em português técnico | Reavaliar na Fase 3 |
| Text-to-SQL | Flexível | Ignora regras de negócio; risco de expor dados; difícil de auditar | Descartada |
| Agentes que gravam sozinhos | Mais automação | Responsabilidade por OS e documentos de segurança exige humano | Descartada; Fase 3 só com aprovação |

**Consequências.**

- Positivas: isolamento herdado do RLS; IA não ganha mais poder que o usuário; custo previsível.
- Negativas: dependência de provedor externo; possível transferência internacional; o pgvector compartilha recursos com o banco transacional.
- Gatilho de revisão: índice vetorial afetando a latência do primário → mover vetores para réplica dedicada ou banco separado, mantendo `empresa_id` e RLS.

---

### ADR-008 — Acesso externo por escopo restritivo e link de uso único

**Status:** Proposta · **Data:** 12/11/2026

**Contexto.** A versão 0.3 do sumário traz usuários de empresas externas para dentro da Tenaz: o perfil `prestador` (subcontratada, assistência técnica, oficina) vê e executa só as OS atribuídas ao seu fornecedor e movimenta só o almoxarifado em poder dele; o perfil `locatario` abre solicitações e acompanha OS só dos ativos locados a ele; o prestador ocasional registra a execução por link de uso único, sem login. O risco R4 (prestador vendo dados de outra contratante ou de outro prestador) é classificado como crítico. O RLS por empresa (ADR-002) não separa usuários dentro da mesma empresa.

**Decisão.**

- Um **escopo** por sessão (`interno`, `prestador`, `locatario`, `link`), resolvido no servidor a partir do banco e definido na transação com `set_config(..., true)`, junto com a empresa.
- Política RLS **restritiva** `escopo_externo` em toda tabela com `empresa_id`: por padrão só o escopo `interno` passa; uma lista fechada de tabelas tem predicado por fornecedor, por locatário ou pela OS do link (seção 6.6).
- Rotas declaram os escopos aceitos (padrão `interno`), DTOs reduzidos para escopo externo e funções amplas (busca, IA, exportação, relatórios, webhooks, chaves de API) indisponíveis.
- Link de uso único com token de 256 bits guardado como hash, entregue no fragmento da URL (`/l/#<token>`), trocado por uma sessão curta em `POST /v1/publico/links-execucao/ativar` e nunca presente no caminho das rotas, expiração, consumo na conclusão, revogação imediata, vínculo à primeira sessão, código de confirmação opcional, rate limit e auditoria completa (seção 6.7).
- Teste automatizado de isolamento por fornecedor e por locatário em todo build, com a mesma regra de bloqueio do teste entre empresas (seção 9.3).

**Alternativas consideradas.**

| Alternativa | A favor | Contra | Resultado |
|---|---|---|---|
| Só filtro na aplicação por fornecedor e locatário | Simples | Mesmo argumento do ADR-002: um filtro esquecido vaza dados de outro prestador | Descartada |
| Política permissiva por perfil no lugar da política por empresa | Uma política só | Mistura dois eixos; um erro no predicado externo poderia abrir a empresa inteira | Descartada: a política restritiva só pode **reduzir** o que a da empresa permite |
| Portal separado com banco ou serviço próprio para terceiros | Isolamento físico | Duplicação de dados e de regras de OS; sincronização entre bases; contradiz o monólito modular | Descartada |
| Link com token curto (6 a 8 caracteres) ou código numérico | Fácil de digitar | Adivinhável; exige rate limit muito agressivo | Descartada; o código numérico só existe como segundo fator opcional |
| Link sem expiração, reutilizável | Conveniente para prestador recorrente | Vira credencial permanente fora de qualquer controle de usuário | Descartada: prestador recorrente recebe usuário `prestador` |
| Token na query string | Padrão comum | Vai para logs da CDN e do balanceador, histórico e `Referer` | Descartada: fragmento da URL |

**Consequências.**

- Revisão da versão 0.3 deste documento: o valor de escopo do prestador passou de `fornecedor` para `prestador` (sumário 3.6); o predicado continua usando `executor_tipo = 'fornecedor'` e `fornecedor_id`, que são campos da OS. A subcontratada com conta própria não usa este escopo: conecta a conta dela (ADR-009).
- Positivas: R4 tratado com a mesma defesa em profundidade do R2; nenhum caminho de escrita novo (o link usa os casos de uso da OS); o acesso externo é auditável e revogável.
- Negativas: cada tabela nova exige decidir e testar a política `escopo_externo`; consultas de locatário dependem de função com `EXISTS` (custo a medir); o link é a primeira rota sem autenticação de usuário e entra no escopo obrigatório do pentest.
- Gatilho de revisão: latência das consultas sob escopo `locatario` acima da meta → gravar o locatário na OS; abuso do link (tentativas em massa ou encaminhamento recorrente) → tornar o código de confirmação obrigatório.

---

### ADR-009 — Conexão entre contas

**Status:** Proposta · **Data:** 12/11/2026

**Contexto.** A versão 0.4 do sumário cria o módulo SRV, no MVP: a subcontratada tem **conta própria** na Tenaz (uma `empresa` com o plano Serviços), com técnicos, clientes atendidos, equipamentos de clientes e OS de serviço. Quando o cliente também usa a Tenaz, as duas contas se conectam (`plataforma.conexao_empresa`) e a OS trafega entre elas: o cliente atribui e libera; a subcontratada distribui, executa e também pode abrir OS num ativo do cliente. Isso quebra uma premissa do ADR-002: pela primeira vez, um mesmo fato de negócio precisa existir em duas empresas. O risco R5 (copiar dado além do permitido ou para a conta errada) é classificado como crítico, no mesmo nível de R2 e R4. A meta de atraso é p95 abaixo de 60 segundos, igual à convergência do app.

**Decisão.**

1. **Cópia por conta.** Cada conta guarda a sua cópia: a OS existe nas duas contas, cada uma com seu número, ligadas por `vinculo_os` (uma linha em cada conta, proposta); os ativos compartilhados existem na subcontratada como `ativo` `de_cliente` com `conexao_empresa_id` e `ativo_remoto_id`. Nenhuma consulta, tela, relatório ou job lê dados de duas empresas.
2. **Processo de sincronização com papel de banco próprio.** Só o modo `tenaz sincronizador` grava dados vindos de outra conta. Ele conecta com o papel `tenaz_sinc`: sem `BYPASSRLS`, com `GRANT` só nas tabelas e **colunas** da lista fechada, e sujeito à política restritiva `conexao_sinc` (só linhas ligadas à conexão em processamento). A credencial é lida só pela tarefa do sincronizador (seções 4.1 e 6.8).
3. **Transações separadas por conta.** Ler na origem, gravar no destino e confirmar na origem são três transações, cada uma com uma única `app.empresa_id`. A conta de destino é sempre a outra ponta da conexão lida em `plataforma.conexao_empresa`, nunca um valor do job ou da mensagem.
4. **Outbox nos dois lados.** A origem grava `envio_conexao` a partir do evento de domínio (outbox transacional da seção 5.5); o destino grava `recebimento_conexao` na mesma transação em que aplica. As duas tabelas existem nas duas contas, porque a OS trafega nos dois sentidos.
5. **Idempotência.** `mensagem_id` = `id` do evento de origem (UUIDv7), com chave única na caixa de entrada. Repetição devolve "já aplicada". Mudança aplicada pela conexão não gera eco para a mesma conexão.
6. **Lista fechada de campos.** Um arquivo versionado define, por grupo e por sentido, os campos que podem trafegar; dele saem os `GRANT` por coluna, os DTOs de mensagem e o teste. Grupos opcionais (documentos técnicos, preços) começam desligados e são ligados pelo lado que envia. Nunca trafegam dados pessoais de técnicos além do nome, custos internos, OS não vinculadas, dados de outros clientes.
7. **Ordem e conflitos por dono de dado.** Sequência por vínculo e por sentido; o destino aplica só a próxima. Não há "último a gravar vence": cada dado tem um dono. O **cliente** é dono do cadastro do ativo e da abertura, descrição, prioridade e prazo da OS que atribui, e da **liberação** (`os.encerrar`), recusa e cancelamento. A **subcontratada** é dona da execução: distribuição, estados de `em_execucao` a `concluida`, apontamentos, componentes, fotos, checklist e assinatura do cliente. Mensagem que altera campo de outro dono é recusada; transição impossível vira pendência sem perder o que foi registrado.
8. **Reprocessamento.** Nova tentativa automática com backoff por até 24 horas; depois, `erro` visível nas duas contas, com ação de reprocessar, ressincronização por estado completo e reconciliação diária por hash dos campos compartilhados.
9. **Auditoria e observabilidade.** Registro em `evento_auditoria` nas duas contas, com o mesmo hash de conteúdo; métricas de atraso (SLO p95 abaixo de 60 s), fila, erros por motivo e divergência; alerta de segurança imediato para ponta inválida, campo fora da lista ou campo de outro dono.
10. **Ciclo de vida.** `convidada` → `ativa` → `suspensa` ou `encerrada`. Suspender acumula mensagens; encerrar interrompe a sincronização sem apagar o histórico de nenhuma das contas.
11. **Teste com três contas** (dois clientes e uma subcontratada conectada aos dois) em todo build, com a mesma regra de bloqueio do teste de vazamento (seção 9.3).

**Alternativas consideradas.**

| Alternativa | A favor | Contra | Resultado |
|---|---|---|---|
| **Leitura cruzada com `BYPASSRLS`**: um processo ou uma visão lê a OS do cliente e grava na subcontratada numa mesma transação | Simples; sem cópia, sem fila, sem atraso | Cria um papel que enxerga todas as empresas: um filtro esquecido vaza a base inteira, e esse papel vira o alvo mais valioso do sistema. Contradiz o ADR-002 e o princípio 1. O teste de vazamento deixaria de valer para esse caminho | Descartada |
| **Conta compartilhada**: a subcontratada trabalha dentro da conta do cliente com o perfil `prestador` | Já existe (ADR-008); nenhuma sincronização | A subcontratada precisa de uma conta por cliente, sem visão única dos seus técnicos, clientes, agenda e relatórios; os técnicos dela teriam de existir em N contas; clientes fora da Tenaz não são atendidos. Continua valendo para subcontratadas sem conta própria | Descartada como solução do SRV; mantida para quem não tem conta |
| **Tudo na conta do cliente**: a OS existe só no cliente e a subcontratada vê uma visão federada das OS de todos os clientes | Fonte única de verdade, sem conflito | O painel da subcontratada (OS por cliente, por técnico) exigiria ler N empresas a cada tela; o histórico some quando a conexão acaba; o app offline do técnico teria de sincronizar com várias empresas que não são a dele; a conta da subcontratada dependeria do plano e da disponibilidade de cada cliente | Descartada |
| **Tabela de OS compartilhada na plataforma**, fora do RLS, com uma linha para as duas empresas | Sem cópia | Dados de negócio de duas empresas numa linha sem RLS por empresa; toda tela e todo job precisariam de regra especial; nenhuma das defesas do ADR-002 se aplica | Descartada |
| **Sincronização pela API pública**: cada conta chama a API da outra com um token de integração por conexão, alimentado por webhooks | Usa só o contrato público; isolamento forte (cada lado só vê o que a API devolve) | Exige webhooks (plano Profissional) dos dois lados; tokens de longa duração por conexão para guardar e girar; mais pontos de falha (HTTP, rate limit) para cumprir o p95 de 60 s; a lista de campos ficaria espalhada entre DTOs públicos | Plano B, se o sincronizador interno se mostrar difícil de auditar |

**Consequências.**

- Positivas: R5 tratado com defesa em profundidade (papel próprio, privilégio por coluna, política restritiva, destino calculado, idempotência, auditoria nos dois lados, teste com três contas); cada conta continua isolada pelo RLS e audita o que recebeu; encerrar a conexão não apaga história; o SRV reaproveita OS, apontamento, ativos e app (risco R6).
- Negativas: dados duplicados e atraso de até 60 s entre as contas; um componente novo para operar (o sincronizador); toda mudança nos campos de OS ou de ativo exige decidir se entra na lista; regras de dono de dado e de transição remota a manter em `4_FLUXOS_PROCESSO.md`; o sincronizador entra no escopo obrigatório do pentest.
- Gatilho de revisão: atraso p95 acima de 60 s sem correção simples; mais de 1 divergência por 1.000 vínculos na reconciliação mensal; qualquer incidente de dado copiado para a conta errada → revisar o desenho e avaliar o plano B.

---

### ADR-010 — Configuração por segmento com catálogos fechados

**Status:** Proposta · **Data:** 12/11/2026

**Contexto.** A versão 0.5 do sumário torna a Tenaz configurável por segmento (seção 3.7). A origem é o caso de um fornecedor de radares de fiscalização eletrônica, instalados e mantidos por prestadores em várias cidades: precisa registrar lacres na OS, impedir que o radar volte a operar depois de romper o lacre ou vencer a verificação metrológica, implantar pontos por etapas (autorização, energia, fundação e poste, laços, montagem, aferição, entrega) e medir a disponibilidade exigida em contrato com a prefeitura. Outros segmentos do público-alvo (telecomunicações, energia, saneamento, elevadores, além de transportadoras, indústria, comércio e locadoras) têm pedidos do mesmo tipo, cada um com o seu vocabulário. O risco R7 do sumário é a configuração crescer até virar um sistema "faça você mesmo", com estados e fluxos livres: indicadores deixam de ser comparáveis, os testes deixam de cobrir o que roda e a implantação fica cara. Ao mesmo tempo, a configuração não pode abrir brecha no isolamento (R2, R4, R5) nem quebrar o app offline, que valida regras sem conexão (ADR-004).

**Decisão.**

1. **Núcleo canônico fixo.** Estados, fórmulas dos indicadores, regras do app offline, isolamento (RLS, escopos, papel do sincronizador) e a lista de campos da conexão entre contas não são configuráveis. Ficam no código e nos testes (seção 5.6).
2. **Configuração por blocos de catálogo fechado.** A empresa configura só blocos que o código conhece (classes de ativo, campos extras da OS, tipos de documento, regras, fase de ciclo de vida e, na Fase 2, modelos de projeto, contratos com clientes e pacotes), cada um com um pacote Go dono. Nenhum bloco cria tabela, coluna, estado, tela, política nem escopo.
3. **Regras condição → efeito com catálogos no código.** Cada condição e cada efeito é um tipo Go registrado numa lista compilada; o banco só aceita os códigos do catálogo; condição ou efeito novo entra por requisito, migração e teste. Uma regra nunca muda estado de OS, nunca altera fórmula e nunca grava em outra conta, e um efeito de regra nunca dispara outra regra.
4. **Avaliação por evento interno e por rotina diária idempotente, na fila de jobs existente.** As condições de fato são avaliadas pelo assinante do outbox (seção 5.5); as de tempo, e a repescagem das de fato, por uma rotina diária com um job por empresa (padrão da seção 6.4). Tudo roda no River (ADR-005), numa fila própria `regras`, sem componente novo. Cada disparo grava `disparo_regra` com chave única pelo fato, na mesma transação do efeito: repetir o evento ou a rotina não repete o efeito.
5. **Campos extras em jsonb validados pela definição.** Valores em `ordem_servico.campos_extras`, definição em `campo_extra_os` (dono: OS). A mesma função de validação existe em Go e em `pacotes/dominio`, com casos compartilhados; o app valida antes de enfileirar o comando, e o servidor revalida ao aplicar. Campos extras não entram em indicador nem na lista de campos da conexão.
6. **Pacotes como catálogo global aplicado por cópia idempotente.** Pacotes versionados no repositório, validados no CI e publicados no catálogo global `pacote_segmento` (sem `empresa_id`, só leitura para as contas). Aplicar copia para a conta, pela API pública de cada módulo dono, só o que falta pela chave natural; nunca altera nem apaga. Depois da cópia, nada na conta depende do pacote.

**Alternativas consideradas.**

| Alternativa | A favor | Contra | Resultado |
|---|---|---|---|
| **Editor livre de fluxos ou estados (low-code)**: a empresa desenha os estados da OS, as transições e as telas | Atende qualquer segmento sem desenvolvimento; argumento de venda forte | Indicadores deixam de ser comparáveis (MTTR de quem tem um estado a mais não é o MTTR de quem não tem); a máquina de estados duplicada no app offline, as regras de conflito (seção 10.5) e o dono de cada dado na conexão entre contas (seção 6.8) dependem de estados conhecidos; o número de combinações torna os testes impossíveis; a implantação vira projeto de consultoria. É exatamente o risco R7 | Descartada |
| **Motor de regras genérico com expressões ou scripts** (linguagem de expressão como CEL ou expr, scripts em Lua ou JavaScript em sandbox, motor externo) | Combinações ilimitadas sem requisito novo; menos código por condição | Uma expressão pode ler dados que o efeito não deveria ver e escapa da revisão por requisito; custo de avaliação sem limite na rotina diária; sem teste por combinação; sandbox de script é superfície de ataque nova e, como serviço externo, um componente a mais sem gatilho mensurável (princípio 2); erros de regra viram suporte | Descartada. Gatilho de revisão abaixo |
| **Campos extras como colunas criadas por empresa** (DDL dinâmico) **ou em tabela chave-valor** (EAV) | Coluna nativa: tipo e índice do banco. EAV: sem mudar o esquema | DDL dinâmico: tabelas são compartilhadas no modelo pool (ADR-002), então a coluna de uma empresa aparece para todas; `ALTER TABLE` em produção bloqueia a tabela da OS; o sqlc tipado e as migrações expandir e contrair deixam de valer; o esquema diverge por empresa. EAV: uma linha por valor, junções pesadas em toda listagem, validação de tipo fraca, mais uma tabela no pull do app e na política `escopo_externo` | Descartadas: jsonb numa coluna da OS, validado pela definição |
| **Pacote por referência**: a conta aponta para o pacote e herda as mudanças de cada versão | Melhorias do pacote chegam a todos sem ação | Mudança no catálogo global altera a configuração de muitas contas sem aviso; conflita com o que a empresa editou; toda consulta de configuração passaria a ler o catálogo global em tempo de execução, misturando dado global e dado da empresa | Descartada: cópia idempotente |
| **Avaliar regras só por evento** ou **só pela rotina diária** | Um mecanismo só | Só evento: as condições de tempo (vencimento) não têm evento no instante certo, e um evento perdido nunca é recuperado. Só rotina diária: até 24 horas entre o lacre rompido e o bloqueio, com o radar operando nesse intervalo | Descartadas: os dois, com a mesma chave de idempotência |

**Consequências.**

- Positivas: o R7 é tratado por construção, não por disciplina: o que não está no catálogo não existe no banco nem no binário. Os indicadores continuam comparáveis entre empresas e segmentos. O isolamento não muda: toda tabela nova tem `empresa_id`, RLS e `escopo_externo`, e só `pacote_segmento` entra na lista de tabelas globais. O app offline continua validando sem conexão, com a mesma função do servidor. Nenhum componente novo na stack.
- Negativas: cada condição ou efeito novo custa um ciclo de requisito, código, migração e teste, e alguns pedidos de clientes vão esperar por isso. Campos extras em jsonb têm filtro e ordenação menos eficientes que coluna nativa e exigem índice GIN. A rotina diária soma carga ao primário no início do dia. A regra de conflito do campo obrigatório criado com o aparelho offline (seção 10.5) aceita OS concluída sem o campo novo, com pendência para o gestor. Os campos extras não trafegam na conexão entre contas no lançamento: a regra do cliente com `os_concluida_com_campo` não dispara para OS executada por subcontratada conectada (seção 16).
- Métricas **(proposta)**: atraso entre o fato e o efeito na avaliação por evento (p95), duração da rotina diária por empresa, disparos por regra e por dia, efeitos com erro por motivo, itens criados e ignorados por aplicação de pacote.
- Gatilhos de revisão: (a) rotina diária de regras de uma empresa com 2.000 ativos acima de 5 minutos, ou alertas das 07h00 atrasados por causa dela (referência RNF-009) → mover a busca de candidatos para a réplica ou tornar a avaliação incremental; (b) atraso p95 entre o fato e o efeito acima de 60 segundos sem correção simples → revisar a fila `regras`; (c) mais de 6 condições ou efeitos novos por requisito em 12 meses, ou o mesmo pedido fora do catálogo vindo de 3 ou mais clientes pagantes → avaliar uma linguagem de expressão restrita **só para o filtro** da regra, com efeitos ainda fechados, num ADR novo; (d) filtro por campo extra acima da meta de RNF-008 em listas → promover o campo mais usado a coluna canônica, por requisito, para todas as empresas.

---

### ADR-011 — Papel de banco próprio para as credenciais de autenticação

**Status:** Proposta · **Data:** 12/11/2026 · **Emenda:** 12/11/2026

**Emenda de 12/11/2026 (mesma etapa E01).** A revisão do `tenaz-qa-automacao` encontrou o achado ISO-E01-08: as sete tabelas de credencial eram de dono `tenaz_dono`, e o `tenaz_manutencao` (`LOGIN`, membro de `tenaz_dono` com herança) as lia inteiras, com `senha_hash` e segredo de MFA. A decisão original fixava quem tem privilégio, mas não quem é o dono, e numa tabela sem RLS o dono é o privilégio mais amplo de todos. A emenda acrescenta o parágrafo final da decisão 1, a verificação de dono na decisão 9, duas alternativas descartadas e o risco residual do dono do esquema. Correção aplicada nas migrações `00004` e `00005` e no documento 2 (seções 2.2, 2.6, 2.7, 5.1, 14.2 e 17.3) pelo `tenaz-banco-dados`.

**Contexto.** A etapa E01 criou no esquema `plataforma` as tabelas de autenticação: `identidade`, `acesso_empresa`, `sessao`, `token_renovacao`, `codigo_autorizacao`, `desafio_mfa` e `token_email` (nomes do sumário 0.6, seção 3.2). Elas ficam fora do RLS por empresa porque são lidas antes de a empresa ser conhecida: o e-mail do login, o cookie, o refresh token, o código OAuth e o link do e-mail chegam sem empresa, e uma identidade serve a várias empresas. O desenho inicial concedia os privilégios ao `tenaz_app` e se apoiava em três regras: só o hash dos segredos, empresa fixada na criação e chave estrangeira composta para `app.usuario`. A revisão de isolamento (achado ISO-E01-02) demonstrou que isso impede mover uma credencial de empresa, mas não lê-la nem forjá-la: o `tenaz_app`, papel de todo módulo de negócio, com o contexto da empresa A, lia `senha_hash` e `mfa_segredo_cifrado` de todas as identidades, lia `acesso_empresa` inteiro (quais clientes cada subcontratada atende), lia sessões de outra empresa com IP e user agent, criava sessão web válida para usuário de outra empresa com o cookie que escolhesse, criava código OAuth e trocava a senha de qualquer identidade. Um defeito em qualquer consulta de qualquer módulo viraria vazamento entre empresas (R2) e tomada de conta. Pelo princípio 1, um build assim não vai para produção. O precedente é o `tenaz_sinc` (ADR-009): quando a defesa não pode ser o RLS por empresa, ela é um papel próprio, com credencial própria e privilégio exato.

**Decisão.**

1. **Papel próprio `tenaz_autenticacao`.** `NOLOGIN`, `NOBYPASSRLS`, sem objetos, fora dos privilégios padrão de `app`, sem relação de membro com `tenaz_dono`, `tenaz_app` ou `tenaz_job` em nenhum sentido. Conexão pelo usuário `tenaz_autenticacao_conexao`, que herda o papel. Privilégios exatos, tabela a tabela, na seção 6.9: `SELECT` e `INSERT` nas sete tabelas (só `SELECT` em `acesso_empresa`), `UPDATE` só nas colunas de ciclo de vida, nenhum `DELETE`, e `INSERT` por coluna em `app.evento_auditoria`, sem leitura da trilha.

   **Dono das credenciais (emenda ISO-E01-08).** As sete tabelas de credencial e as funções `SECURITY DEFINER` que as tocam (`plataforma.tg_acesso_empresa()` e `plataforma.purgar_credenciais_vencidas()`) pertencem ao papel `tenaz_credencial` (`NOLOGIN`, `NOBYPASSRLS`, sem membros além do papel que aplica migração), e não a `tenaz_dono`. Sem RLS, o dono lê tudo; com dono `tenaz_dono`, o `tenaz_manutencao`, membro do dono com herança, lia todas as credenciais. `tenaz_dono` recebe só `REFERENCES (id)` em `identidade` e `EXECUTE` em `tg_acesso_empresa()`; `tenaz_credencial` recebe `REFERENCES` nas colunas da chave composta de `app.usuario` (`empresa_id`, `id`, `identidade_id`), mais `USAGE` e `CREATE` em `plataforma` e `USAGE` em `app` e `public`, para criar os objetos; nenhum outro privilégio entre os dois (tabela na seção 6.9). A mesma regra vale para toda tabela futura de `plataforma` sem RLS que não seja catálogo global de leitura: o dono é um papel próprio `NOLOGIN` sem membros de runtime, nunca `tenaz_dono`.
2. **Quem perde.** `tenaz_app` e `tenaz_job` perdem todo privilégio nas sete tabelas; `tenaz_leitura` e `tenaz_sinc` nunca tiveram. A purga das credenciais vencidas passa a ser a função `plataforma.purgar_credenciais_vencidas()`, `SECURITY DEFINER`, sem parâmetro, executada só pelo `tenaz_job`.
3. **Quem usa o pool.** Só o pacote `plataforma/autenticacao`, pelo repositório interno a ele (regra 12 da seção 5.3). O pool existe só no modo `tenaz api`; o middleware de sessão chama o serviço e não toca o pool. O repositório recusa rodar dentro de `NaEmpresa` e nunca chama `NaEmpresa` nem `app.conferir_escopo()`.
4. **Sessão primeiro, empresa depois.** A cada requisição: o pool da autenticação resolve a sessão e devolve só `empresa_id`, `usuario_id`, `identidade_id`, `canal` e `sessao_id`; o pool da aplicação carrega o ator em `NaEmpresa` com escopo interno e sem usuário e confere a identidade; o caso de uso roda em `NaEmpresa` com o escopo do cadastro (seção 6.3). Nunca na mesma transação.
5. **Auditoria de autenticação pelo próprio papel.** Os eventos de login, falha, logout, troca de senha e MFA são gravados pelo `tenaz_autenticacao` em `app.evento_auditoria` de cada empresa da identidade, numa transação que define só `app.empresa_id` e os dados da requisição. A política permissiva de inserção passa a incluir o papel, e uma política restritiva própria limita a `entidade = 'identidade'`, às ações `login`, `login_falho`, `logout` e `alteracao` e ao usuário daquela identidade naquela empresa, conferido em `acesso_empresa`. Assim, sessão, zerar das falhas e evento de login na empresa da sessão continuam numa transação só.
6. **`acesso_empresa` continua do trigger.** `plataforma.tg_acesso_empresa()`, `SECURITY DEFINER`, disparado pela escrita do `tenaz_app` em `app.usuario`, não depende de privilégio do `tenaz_app` na tabela e continua igual; passa a revogar também as sessões e os refresh tokens do usuário naquela empresa quando o acesso fica desabilitado.
7. **Operações que atravessam os dois papéis são sequências de transações**, na ordem da seção 6.9 (convite, aceite, troca de empresa, encerrar sessões), com a credencial de uso único consumida por último e cada passo idempotente.
8. **Credencial própria.** `TENAZ_BANCO_URL_AUTENTICACAO` e `TENAZ_BANCO_AUTENTICACAO_MAX_CONN`; segredo próprio lido só pela tarefa da API; recusa de partida se a URL for igual à da aplicação fora do desenvolvimento (seções 4.1 e 6.9).
9. **Teste em todo build**, com a mesma regra de bloqueio da suíte de vazamento: mapa de privilégios nos dois sentidos, reprodução dos ataques do ISO-E01-02 esperando "permissão negada", molde da auditoria, consultas de credencial só no repositório, lista aprovada de funções `SECURITY DEFINER` e, desde a emenda, dono das sete tabelas e das duas funções igual a `tenaz_credencial`, sem membros de runtime, com `tenaz_manutencao` recebendo "permissão negada" nas credenciais (seção 9.3).

**Alternativas consideradas.**

| Alternativa | A favor | Contra | Resultado |
|---|---|---|---|
| **Manter o `tenaz_app`** com privilégio por coluna, só hash e chave estrangeira composta (desenho inicial de E01) | Um pool, nenhum papel novo | Demonstrado pelo ISO-E01-02: impede mover a credencial, não impede ler `senha_hash` e segredo de MFA de todas as identidades nem forjar sessão de outra empresa. Toda consulta de todo módulo vira superfície de tomada de conta | Descartada |
| **RLS por variável `app.identidade_id`** com funções `SECURITY DEFINER` de busca por hash ou e-mail (opção B da revisão de isolamento) | Um pool; transações de negócio, sem a variável, não veem credenciais | A variável é definida pelo mesmo código que se quer conter; login e validação de sessão começam antes de haver identidade, então as funções de busca devolvem a credencial ao `tenaz_app` de qualquer forma; o `INSERT` da sessão continua com o `tenaz_app`, que escolhe o cookie: **não impede forjar sessão** | Descartada |
| **Variante do papel próprio com funções `SECURITY DEFINER` estreitas como única interface** (login por e-mail exato, validar sessão por hash, rotacionar refresh token...), sem `SELECT` direto | Nem o pacote de autenticação enumera a tabela | O serviço inteiro, cerca de vinte operações (falhas e bloqueio, desafio de MFA, sessão web e do app, rotação com detecção de reuso, código PKCE, links de e-mail, troca de empresa), viraria PL/pgSQL rodando como o dono das tabelas; lógica duplicada entre Go e SQL, testes mais caros, e a lista aprovada de funções `SECURITY DEFINER` passaria de duas para mais de vinte, cada uma um ponto de escalada. O ganho é pequeno: "login por e-mail" devolve o hash ao chamador, que é o mesmo pacote; o papel próprio já garante que só esse pacote chama, e a enumeração é coberta por teste de consulta | Descartada como interface geral. **Adotada só onde a interface é fixa e não recebe dado da requisição**: o trigger de `acesso_empresa` e a purga dos vencidos |
| **Credenciais no esquema `app`, com `empresa_id` e RLS** | A mesma defesa de todo o resto | A identidade é global (uma pessoa em várias empresas) e a busca acontece antes de a empresa ser conhecida | Descartada |
| **Gravar a auditoria de login pelo `tenaz_app`**, no pool da aplicação, depois da transação da credencial | O `tenaz_autenticacao` não teria nenhum privilégio em `app` | Sessão e evento de login deixariam de ser atômicos (sessão sem evento, ou evento sem sessão), e o `tenaz_app` passaria a gravar fatos de identidade que não consegue conferir | Descartada: o papel grava, contido pela política restritiva, como o `tenaz_sinc` |
| **Purga pelo `tenaz_autenticacao` no worker** | Nenhuma função nova | A credencial que lê `senha_hash` passaria a existir também na tarefa do worker | Descartada |
| **Emenda ISO-E01-08: RLS forçado nas sete tabelas**, com política só para `tenaz_autenticacao`, mantendo o dono `tenaz_dono` | Nenhum papel novo | `FORCE ROW LEVEL SECURITY` vale para o dono, então as funções `SECURITY DEFINER` do trigger e da purga, que rodam como dono, seriam barradas; dar política ao dono para liberá-las reabriria o acesso ao `tenaz_manutencao`, que herda dele. E o dono pode desligar o `FORCE` com um `ALTER TABLE` | Descartada |
| **Emenda ISO-E01-08: tirar a herança do `tenaz_manutencao`** (`NOINHERIT`) | Nenhum papel novo; corrige a causa imediata | Não resolve: o membro continua podendo `SET ROLE tenaz_dono` e ler como dono, e todo membro futuro de `tenaz_dono` reabre o achado. A defesa ficaria num atributo de papel, não na posse do objeto | Descartada |
| **Provedor de identidade externo** (Cognito, Keycloak) ou serviço de autenticação separado | Isolamento físico; MFA e bloqueio prontos | Componente novo sem gatilho mensurável (princípio 2); segunda fonte de verdade de usuário e empresa; residência e custo. O SSO e o SAML de E10 federam o provedor do cliente e não exigem trocar o nosso | Descartada; gatilho (d) abaixo |

**Consequências.**

- Positivas: o achado ISO-E01-02 deixa de existir por construção, com "permissão negada" e não zero linhas; quem lê ou grava credencial passa de "todo módulo" para um pacote; o `tenaz_job` deixa de ler IP e user agent de todas as sessões; a inativação revoga as credenciais na mesma transação; o padrão é o mesmo do `tenaz_sinc`, já conhecido pela suíte.
- Negativas: um segredo e um pool a mais na API; login, troca de empresa, convite e aceite viram sequências de transações em dois pools, com ordem definida e passos idempotentes; um convite pode ficar sem token se o segundo passo falhar, e o administrador reenvia; o `tenaz_autenticacao` tem `INSERT` em `app.evento_auditoria`, contido pela política restritiva; a política consulta `acesso_empresa` a cada evento (busca pela chave primária). Desde a emenda: um papel dono a mais, migrações que alternam `SET ROLE` entre os dois donos, e `VACUUM` e `ANALYZE` manuais nas sete tabelas só pelo papel de migração (o autovacuum continua valendo, e a purga diária mantém as tabelas pequenas).
- Risco residual: o `tenaz_app` ainda grava `usuario.identidade_id`, que só pode ser preenchido uma vez. Ligar um usuário de A à identidade de outra pessoa não expõe dado de B, só faria A aparecer na lista de empresas dela; e o `tenaz_app` não consegue mais descobrir identificadores de identidade, porque não lê a tabela. A escrita fica restrita ao caso de uso de aceite e coberta por teste.
- Risco residual da emenda: o esquema `plataforma` continua de dono `tenaz_dono`, e o dono do esquema pode apagar tabelas que não são dele. Com a credencial do `tenaz_manutencao` roubada, alguém não lê as credenciais, mas pode apagar as sete tabelas, recriá-las como `tenaz_dono` e conceder privilégio ao `tenaz_autenticacao`, capturando o que for gravado depois (hashes de sessão, nova senha), ou desligar o trigger `trg_acesso_empresa` em `app.usuario`, que é dele, e a inativação deixa de revogar sessões. É um ataque destrutivo e ruidoso (todo login quebra), mas é a mesma classe do achado: o `tenaz_manutencao` herda o poder de DDL do dono comum, e o dono de uma tabela de `app` também pode desligar o `FORCE ROW LEVEL SECURITY` dela. Pendência na seção 16, a resolver antes da primeira implantação em homologação.
- Gatilhos de revisão: (a) outro pacote precisar das tabelas de credencial (SSO e SAML em E10, login por telefone da P11) → o código entra em `plataforma/autenticacao` ou este ADR é revisto; nunca por privilégio ao `tenaz_app`; (b) incidente, achado de pentest ou revisão que encontre leitura ampla pelo próprio pacote de autenticação → adotar a variante de funções estreitas para `identidade`; (c) espera por conexão do pool da autenticação com p95 acima de 50 ms na API **(proposta)** → aumentar o pool, nunca juntar os dois; (d) cliente Enterprise exigindo provedor de identidade próprio que o SSO de E10 não resolva → avaliar provedor externo num ADR novo.

---

### ADR-012 — Só software livre ou gratuito, e SeaweedFS no lugar do MinIO

**Status:** regra de admissão, alcance e exceção **Aceitos** (dono do produto, 13/11/2026); SeaweedFS no desenvolvimento e verificação no CI **Proposta** · **Data:** 13/11/2026 · **Emenda:** 13/11/2026

**Emenda de 13/11/2026 (mesma versão 0.7).** Respostas do dono do produto à auditoria abaixo, já no sumário 0.7 (seção 3.6, "Ferramentas e recursos"), em [CLAUDE.md](../CLAUDE.md) (seção 4) e no prompt de E01:

- **Alcance da regra**, nas palavras dele: "Ambiente de desenvolvimento é diferente de ambiente de produção. A questão de gratuidade se aplica exclusivamente para ferramentas utilizadas diretamente pelo sistema." Ficam **fora** da regra as ferramentas de quem desenvolve: Docker Desktop, navegador (inclusive o do `md2pdf.py`), editor, GitHub e GitHub Actions. Ficam **dentro** o que vai para a produção: dependências de `backend/go.mod`, `web/`, `app/` e `pacotes/`, imagens de produção, servidores e serviços que o sistema chama. Software proprietário gratuito sem condição usado pelo sistema (por exemplo, o Expo Push) atende ao item (b) da decisão 1. Isso resolve H2, H3 e H4.
- **Exceção explícita:** hospedagem de produção em nuvem gerenciada (AWS, São Paulo), com banco, chaves e armazenamento gerenciados. O [ADR-006](#adr-006--nuvem-aws-na-região-são-paulo) volta a valer. Resolve a H1.
- **Selo em produção:** S3 com Object Lock em modo de conformidade, que resiste à própria conta root. Resolve a H10 na produção.
- **Selo e anexos em desenvolvimento:** "por enquanto, em dev, tudo em uma única máquina (docker)": SeaweedFS no compose, endurecido com `security.toml` (decisão 4). O risco residual vale só nesse ambiente, que não tem dado real.
- **Serviços pagos por natureza** (lojas e build de iOS, e-mail transacional, SMS, IA, pagamento) ficam pendentes até a etapa em que entram (seção 16, H5 a H9).

O texto original abaixo fica como registro; onde a emenda muda o efeito, a mudança está marcada.

**Contexto.** Na etapa E01 o dono do produto fixou o requisito: **não usar nenhuma ferramenta ou recurso que não seja 100% gratuito ou open source** (prompt [E01](prompts/E01_conta-e-acesso.md), "Ajustes durante a etapa", linha "Ferramentas e recursos"). O primeiro alvo é o MinIO, que o compose de desenvolvimento previa para os anexos: a edição comunitária deixou de ter binários e imagens publicados, a tag `latest` não aceita mais pull anônimo (comentário do compose, etapa 011) e o repositório `minio/minio` está arquivado (conferido pela API do GitHub em 13/11/2026: `archived = true`, último push em 24/04/2026). O sistema precisa de armazenamento compatível com S3 com **bloqueio de objeto em modo de conformidade**, para o selo diário da trilha (RF-ADM-011 em [1_ESPECIFICACAO_REQUISITOS.md](1_ESPECIFICACAO_REQUISITOS.md), versão 0.5; `selo_auditoria.chave_objeto` em [2_ARQUITETURA_BANCO_DADOS.md](2_ARQUITETURA_BANCO_DADOS.md)), e de URL pré-assinada para os anexos a partir de E02 (seção 9.3, linha "Arquivos"). O requisito alcança também a nuvem do ADR-006 e os serviços das seções 10 a 14, que são pagos por natureza: hospedagem nenhuma é gratuita.

**Decisão.**

1. **Regra de admissão de software.** Biblioteca, ferramenta de desenvolvimento ou de CI, imagem de contêiner e sistema operado pela Tenaz só entram se forem (a) de código aberto, com licença aprovada pela OSI (MIT, BSD, ISC, Apache-2.0, MPL-2.0, LGPL, GPL, AGPL), ou fonte tipográfica sob OFL-1.1; ou (b) gratuitos **sem condição**: sem cota que vire cobrança, sem limite de porte ou de receita da empresa, sem período de teste. **(Emenda: a regra vale só para o que o sistema usa diretamente; ferramenta de desenvolvimento e de CI fica fora, e software proprietário gratuito sem condição usado pelo sistema atende ao item (b).)** **Não entram:** licença "source-available" (BSL, SSPL, FSL, Elastic License, Commons Clause, RSAL; por exemplo Terraform, HashiCorp Vault e Nomad, Redis desde a 7.4, a edição comercial do MinIO) e o recurso pago de produto "open core", mesmo com a base aberta (River Pro, sqlc Cloud, Testcontainers Cloud, Maestro Cloud, edições Enterprise). Software GPL ou AGPL roda como processo separado e sem modificação; modificar o código dele exige revisão antes, porque a AGPL obriga a publicar o código modificado oferecido por rede.
2. **Verificação no CI (proposta, com o alcance da emenda).** Tarefa `licencas` em `ferramentas/tarefas.sh`, chamada por `conferir`, **só sobre o que vai para o produto**: `go-licenses` (Apache-2.0) sobre as dependências do binário `cmd/tenaz` (sem as só de teste), `pnpm licenses list --prod` sobre `web`, `app` e `pacotes`, e as imagens de produção (base do `backend/Dockerfile` e o que for implantado), contra a lista de licenças da decisão 1 e um arquivo versionado com licença e versão fixa (tag concreta, nunca `latest` nem tag flutuante). Ficam de fora as `devDependencies`, as dependências só de teste (testcontainers e afins), o compose de desenvolvimento e as ferramentas do CI. Licença fora da lista quebra o build. A tarefa continua fazendo sentido porque a regra, mesmo menor, precisa de um portão automático: dependência transitiva nova entra sem ninguém ler a licença. Responsável: `tenaz-plataforma-sre`.
3. **SeaweedFS no lugar do MinIO.** SeaweedFS (Apache-2.0), imagem `chrislusf/seaweedfs:4.53`, fixada por versão e digest (`sha256:ccad23e18639c8b98cee17e7100b95830a06077b56d8e5b006d6dcf8ec5e94a0`, conferido em 13/11/2026) e atualizada só depois de passar no ensaio da decisão 5. No desenvolvimento, um contêiner `weed server -s3` (master, volume, filer e gateway S3 num processo), numa máquina só. **(Emenda: em produção o armazenamento é o S3 gerenciado, pela exceção do ADR-006; o SeaweedFS fica no desenvolvimento.)** **(Aplicado em 13/11/2026:** servidor único no compose, no `up` comum, sem perfil; master, volume e filer escutam só dentro do contêiner (`-ip.bind=127.0.0.1`), e só o gateway S3 escuta na rede; telemetria do projeto desligada (`-master.telemetry=false`).**)** Não é componente novo: o object storage compatível com S3 já está na stack (seção 2.4); muda só o servidor.
4. **Endurecimento obrigatório onde o SeaweedFS rodar, inclusive o local.** (a) `security.toml` com chaves JWT de volume e de filer, para leitura e escrita: **sem ele, a API HTTP do filer apaga objeto em retenção de conformidade sem credencial nenhuma** (ensaio abaixo). (b) Só a porta do gateway S3 é publicada; master, volume, filer e gRPC ficam na rede interna do armazenamento. (c) Identidades S3 separadas por uso, com as ações mínimas: a dos anexos sem permissão de bloqueio de objeto; a do selo com gravação e retenção, sem `bypass`; `Admin` só para a operação. (d) Chave de objeto prefixada por `empresa_id` (seção 6.1), versionamento ligado e, no bucket do selo, bloqueio de objeto com retenção padrão em modo `COMPLIANCE`.
5. **O backend fala só S3.** Cliente `aws-sdk-go-v2` (Apache-2.0), com endpoint, região e estilo de caminho por configuração; nenhuma API própria do SeaweedFS, nenhum `minio-go`. Trocar de servidor é trocar configuração. As 15 verificações do ensaio viram teste de integração do armazenamento, contra o contêiner do compose **(proposta**, `tenaz-qa-automacao` com `tenaz-plataforma-sre`), acrescidas do `DELETE` recursivo no filer sem JWT, que precisa receber 401; o teste qualifica toda versão nova da imagem e qualquer substituto. **(Versão 0.8:** a API local e o job do selo usam só a API S3 e a identidade `tenaz_auditoria`, sem permissão de `Admin`.**)**
6. **Serviços pagos por natureza não são decididos aqui.** Hospedagem, domínio, lojas de aplicativo, SMS, e-mail transacional, push do iOS, pagamento e GPU para IA custam dinheiro qualquer que seja o software. A tabela abaixo dá, para cada serviço proposto neste documento, o equivalente autohospedado em software livre e o que continua pago ou não tem substituto. **(Emenda: a última coluna diz o que o dono do produto decidiu; a hospedagem de produção é a exceção, e o resto entra na etapa indicada na seção 16.)**

| Serviço proposto | Equivalente autohospedado em software livre | O que continua pago ou não tem substituto | Situação (emenda) |
|---|---|---|---|
| Computação (ECS Fargate, ALB) | Docker Engine com Compose, ou k3s (Apache-2.0); proxy Caddy (Apache-2.0) ou nginx (BSD-2) | As máquinas (VPS, servidor dedicado ou colocation no Brasil) | Exceção: nuvem gerenciada (ADR-006) |
| Banco gerenciado (RDS Multi-AZ, PITR de 35 dias) | PostgreSQL com Patroni (MIT) e pgBackRest (MIT), ou CloudNativePG (Apache-2.0) sobre k3s | Máquinas em pelo menos dois datacenters para o SLA de 99,9% e o RTO de zona | Exceção: nuvem gerenciada (ADR-006) |
| Cache (ElastiCache) | Valkey (BSD-3) | Nada além da máquina | Exceção: nuvem gerenciada (ADR-006) |
| Arquivos e cópia imutável (S3, Object Lock, AWS Backup com Vault Lock) | SeaweedFS (decisão 3); cópias com restic (BSD-2) e pgBackRest para um segundo SeaweedFS em outro local, com bloqueio de objeto | O segundo local; a imutabilidade não resiste a quem administra o host (risco residual abaixo) | Exceção em produção: S3 com Object Lock em conformidade; SeaweedFS no desenvolvimento |
| Segredos e chaves (Secrets Manager, KMS) | OpenBao (MPL-2.0, fork aberto do Vault), ou SOPS (MPL-2.0) com age (BSD-3) | Sem HSM: a chave-mestra fica protegida por procedimento, não por hardware | Exceção: nuvem gerenciada (ADR-006) |
| DNS (Route 53) | Servidores DNS do registro.br, sem custo para domínio registrado nele, ou PowerDNS (GPL-2.0) | **O domínio**: registro anual no registro.br, sem substituto gratuito | Exceção: nuvem gerenciada; o domínio segue a P5 |
| CDN e WAF (CloudFront, AWS WAF) | Caddy ou nginx com Coraza (Apache-2.0) e OWASP CRS (Apache-2.0); certificados do Let's Encrypt (gratuito) | Sem CDN global; a defesa contra DDoS volumétrico depende do provedor das máquinas | Exceção: nuvem gerenciada (ADR-006) |
| E-mail transacional (SES) | Postfix (IPL-1.0 ou EPL-2.0) com OpenDKIM (BSD); Mailpit só no desenvolvimento | Entregabilidade: IP próprio com reputação e porta 25 liberada pelo provedor; risco real de cair em spam | Pendente: H6, até a primeira homologação (Mailpit no desenvolvimento) |
| Observabilidade (CloudWatch, X-Ray, ADOT) | Coletor OpenTelemetry, Prometheus, Jaeger, VictoriaMetrics e VictoriaLogs (todos Apache-2.0); Grafana (AGPL-3.0) | Máquinas e disco para a retenção de 13 meses | Exceção: nuvem gerenciada (ADR-006) |
| Segurança da conta (CloudTrail, GuardDuty, Security Hub) | Wazuh (GPL-2.0), Falco (Apache-2.0), auditd | Pessoas para olhar os alertas | Exceção: nuvem gerenciada (ADR-006) |
| Página de status (RNF-020) | Uptime Kuma (MIT) ou Gatus (Apache-2.0), numa máquina fora da infraestrutura principal | A máquina externa | Com a configuração da seção 13, em E12 |
| Push (Expo Push para APNs e FCM) | Envio direto do worker ao APNs e ao FCM; no Android, ntfy com UnifiedPush (Apache-2.0) | O APNs exige a conta de desenvolvedor Apple, paga; Expo Push e FCM são gratuitos, mas proprietários (H2) | Expo Push e FCM atendem (gratuitos sem condição); a conta Apple segue a H5, E12 |
| Distribuição do app (EAS Build, EAS Update, lojas) | Build local (`eas build --local`, eas-cli MIT, ou `expo prebuild` com Gradle); servidor próprio do protocolo `expo-updates` para a atualização sem loja | **As lojas**: Apple Developer Program (anual) e Google Play (taxa única), sem substituto para iOS; o build de iOS exige macOS e Xcode (gratuito, proprietário, em hardware Apple) | Pendente: H5, E12 |
| SMS (P11) | Nenhum: a operadora cobra cada mensagem; um gateway num celular Android só troca o provedor pelo plano da operadora | Tudo. Saída sem SMS: e-mail, TOTP, convite por QR code | Pendente: H7 (P11), E12 |
| IA (ADR-007, Fase 2) | Modelos de pesos abertos sob Apache-2.0 servidos por vLLM (Apache-2.0) ou Ollama (MIT) | GPU; qualidade a medir no conjunto de avaliação da seção 12.5 | Pendente: H8, Fase 2 |
| Pagamento (RF-ADM-009) | Faturamento próprio ou Kill Bill (Apache-2.0) | Gateway, banco ou Pix: tarifa por transação, sem substituto | Pendente: H9, E10 |

**Ensaio de 13/11/2026**, fora do repositório, num diretório de rascunho e numa rede Docker isolada, sem tocar no compose. Cliente `amazon/aws-cli:2.40.2` e `curlimages/curl:8.17.0`. As 15 verificações: (1) bucket criado com bloqueio de objeto; (2) versionamento ligado automaticamente; (3) retenção padrão `COMPLIANCE` de 1 dia; (4) gravação do selo com retenção `COMPLIANCE` explícita; (5) retenção lida de volta; (6) apagar a versão protegida: `AccessDenied`; (7) apagar com `--bypass-governance-retention`: `AccessDenied`; (8) encurtar a retenção: `AccessDenied`; (9) trocar para `GOVERNANCE`: `AccessDenied`; (10) sobrescrever a chave: cria versão nova; (11) apagar sem versão: só marcador de exclusão; (12) apagar o bucket: recusado; (13) a versão protegida continua lá com o conteúdo original; (14) URL pré-assinada de 5 minutos: HTTP 200; (15) URL adulterada: HTTP 403.

| Candidata | Licença | Resultado das 15 verificações | Memória em repouso | Imagem |
|---|---|---|---|---|
| SeaweedFS 4.53 | Apache-2.0 | 15 de 15 | cerca de 100 MiB | 792 MB |
| RustFS 1.0.5 | Apache-2.0 | 15 de 15 | cerca de 31 MiB | 346 MB |
| Versity Gateway v1.11.0 | Apache-2.0 | 15 de 15 | cerca de 13 MiB | 82 MB |

Dois achados no SeaweedFS, que moldam as decisões 4 e 5 e o risco residual: **sem `security.toml`**, um `DELETE` recursivo na API HTTP do filer (porta 8888, sem credencial) apagou o bucket inteiro, com a versão em retenção de conformidade (HTTP 204); **com** o JWT do filer, o mesmo pedido recebeu 401 e a versão ficou intacta. E `weed shell` dentro do contêiner, com acesso de operador, apagou o bucket mesmo com o JWT.

**Auditoria da stack contra a regra (13/11/2026).** Dependências reais do repositório e itens citados como software nos documentos 2, 5, 6 e 9:

| Item | Onde | Licença | Veredito |
|---|---|---|---|
| Go 1.27 e imagem `golang:1.27-alpine` | backend, Dockerfile, CI | BSD-3 | ok |
| Huma v2, pgx v5, goose v3, `google/uuid`, `x/crypto`, testify | `go.mod` | MIT; MIT; MIT; BSD-3; BSD-3; MIT | ok |
| testcontainers-go | `go.mod` | MIT | ok; o Testcontainers Cloud, pago, fica de fora |
| OpenTelemetry Go | `go.mod` (indireta) | Apache-2.0 | ok |
| River (ADR-005, ainda fora do `go.mod`) | backend | MPL-2.0 | ok; o River Pro, pago, fica de fora |
| sqlc 1.30.0 pela imagem `sqlc/sqlc` | `sqlc.yaml` | MIT | ok; o sqlc Cloud fica de fora |
| golangci-lint | `.golangci.yml` | GPL-3.0, ferramenta não ligada ao binário | ok |
| `gcr.io/distroless/static-debian12:nonroot` | Dockerfile | Apache-2.0, pacotes livres do Debian | ok |
| PostgreSQL 17.11 (`postgres:17.11-alpine3.24`, por digest) e pgvector | compose, documento 2 | PostgreSQL License | ok; os testes com testcontainers sobem a mesma imagem por digest (`dbteste.ImagemPostgres` e `testes/vazamento/ambiente_test.go`, alinhados em 13/11/2026) |
| Valkey 8.1.11 (`valkey/valkey:8.1.11-alpine3.24`, por digest) | compose | BSD-3 | ok |
| Mailpit v1.27.0 (por digest) | compose | MIT | ok |
| MinIO (`quay.io/minio/minio:latest`) | estava no compose, perfil `anexos` | AGPL-3.0; repositório arquivado e sem imagem comunitária | **substituído**, fora do compose desde 13/11/2026 |
| SeaweedFS 4.53 (`chrislusf/seaweedfs:4.53`, por digest) | compose | Apache-2.0 | ok (decisão 3) |
| aws-cli 2.40.2 (`amazon/aws-cli:2.40.2`, por digest, perfil `ferramentas`) | compose, ferramenta de desenvolvimento | Apache-2.0 | ok; fora da regra pelo alcance da emenda |
| `aws-sdk-go-v2` v1.48.0, `service/s3` v1.114.10, `credentials` v1.21.0 e `smithy-go` v1.28.9 (acréscimo da versão 0.8: cliente do armazenamento, decisão 5) | `go.mod`, `internal/plataforma/armazenamento` | Apache-2.0 | ok; vão para o binário de produção. No `go.mod` aparecem como `// indirect` apesar da importação direta (pendência na seção 16) |
| `golang.org/x/text` v0.44.0 (acréscimo da versão 0.9: plural CLDR de `plataforma/textos`, ADR-013) | `go.mod`, agora dependência direta | BSD-3 | ok; já vinha como indireta, nenhuma biblioteca nova |
| Node 22.17, pnpm 10, TypeScript 5.9 | raiz | MIT; MIT; Apache-2.0 | ok |
| React 19, React Router 7, TanStack Query, i18next e react-i18next, openapi-fetch, openapi-typescript, Radix, lucide-react | web, pacotes | MIT; lucide-react ISC | ok |
| Vite, Vitest, Tailwind CSS 4, ESLint e plugins, typescript-eslint, Testing Library, jsdom | web | MIT; typescript-eslint MIT e BSD-2 | ok |
| axe-core | web | MPL-2.0 | ok |
| Fontsource (Archivo, Inter, JetBrains Mono) | web | pacotes MIT; fontes OFL-1.1 | ok |
| Expo SDK 58, expo-router, expo-sqlite, expo-secure-store, expo-splash-screen, expo-application, expo-constants, expo-crypto, expo-local-authentication, expo-web-browser, React Native, react-native-safe-area-context | app | MIT | ok (SDK 58 desde 13/11/2026) |
| `expo-localization` 58 e `expo-clipboard` 58 (acréscimo da versão 0.8: idioma do aparelho, ADR-013; área de transferência) | app | MIT | ok |
| SQLCipher, Drizzle, Jest, Playwright, Maestro, React Hook Form, Zod, OpenFeature (propostas) | seções 2.2, 2.3, 13.4 | BSD-3 (edição comunitária); Apache-2.0; MIT; Apache-2.0; Apache-2.0 (CLI); MIT; MIT; Apache-2.0 | ok; edições comerciais do SQLCipher e o Maestro Cloud ficam de fora |
| OpenTofu | seção 13.2 | MPL-2.0 | ok (o Terraform, BSL, não entra) |
| gitleaks, Trivy, Syft, Cosign, Checkov, OWASP ZAP, Renovate | documento 9 | MIT; Apache-2.0 nos cinco seguintes; Renovate AGPL-3.0, CLI autohospedada | ok |
| Python 3.13 e o pacote `markdown` | `ferramentas/` | PSF; BSD-3 | ok |
| Mermaid 10.9.1 pelo jsDelivr | `ferramentas/md2pdf.py` | MIT; CDN gratuita | ok |
| Chrome ou Edge no `md2pdf.py` | `ferramentas/md2pdf.py` | proprietário, gratuito | fora da regra (ferramenta de desenvolvimento, emenda) |
| Docker Desktop | máquinas de desenvolvimento | proprietário; gratuito só para empresa com menos de 250 pessoas **e** menos de US$ 10 milhões de receita anual | fora da regra (ferramenta de desenvolvimento, emenda); continua |
| GitHub e GitHub Actions | CI | serviço proprietário; gratuito e sem limite de minutos em repositório público | fora da regra (ferramenta de desenvolvimento, emenda); continua |
| Docker Hub e quay.io | compose | serviço gratuito, com limite de pull anônimo | ok |
| PowerSync (plano B do ADR-004) | ADR-004 | edição autohospedada sob FSL (a confirmar na versão vigente) | suspenso |

**Alternativas consideradas** (para o armazenamento de objetos).

| Alternativa | A favor | Contra | Resultado |
|---|---|---|---|
| **RustFS 1.0.5** | Apache-2.0; 15 de 15 no ensaio; leve; um processo só, sem API paralela à S3; desenho próximo do MinIO | Projeto jovem (1.0 em 2026), mantido por uma empresa só, o mesmo perfil de governança que levou ao fim do MinIO comunitário; pouco histórico de segurança em produção | **Alternativa qualificada**: substituto imediato pelos gatilhos (a) e (b) |
| **Versity Gateway v1.11.0** | Apache-2.0; 15 de 15; o mais leve | Gateway sobre sistema de arquivos POSIX: durabilidade e alta disponibilidade dependem do sistema de arquivos embaixo; o bloqueio fica em atributos de arquivo, que quem tem acesso ao disco apaga | Descartada como padrão; serve para servidor único ou CI |
| **Ceph RGW** | LGPL; Object Lock completo e maduro há anos; governança neutra (Ceph Foundation); o mais robusto em escala | Monitores, gerenciadores, OSDs e RGW; mínimo de três nós em produção; impraticável no compose de desenvolvimento | Descartada agora; gatilho (c) |
| **Zenko CloudServer** (Scality) | Apache-2.0; Object Lock implementado | Backend autohospedado é arquivo local ou memória; o distribuído é o RING, comercial; a imagem do Docker Hub está parada desde 2023 | Descartada |
| **Garage** | AGPL-3.0; leve; feito para vários locais | Sem versionamento e sem bloqueio de objeto (página de compatibilidade S3 do projeto, conferida em 13/11/2026) | Descartada: não atende o selo |
| **OpenStack Swift** com `s3api` | Apache-2.0; maduro | Sem Object Lock; operação pesada | Descartada |
| **MinIO compilado da fonte ou fork comunitário** | Comportamento conhecido | Repositório arquivado, sem correção de segurança; fork de terceiro sem garantia de continuidade | Descartada |
| **Só o S3 gerenciado** | O modo de conformidade resiste até à conta root | Serviço pago; não resolve desenvolvimento nem CI | **Adotada em produção pela emenda** (exceção do ADR-006); o SeaweedFS fica no desenvolvimento, com o mesmo contrato S3 |
| **Disco local e hash no banco, sem S3** | Nenhum componente | Sem bloqueio de objeto nem URL pré-assinada; o selo ficaria no mesmo domínio do banco que ele protege | Descartada |

**Consequências.**

- Positivas: licença permissiva, sem copyleft; imagem com versão fixa publicada; um contêiner no desenvolvimento; código preso só à API S3, com um substituto já qualificado pelo mesmo ensaio; a regra de admissão vira verificação automática no CI.
- Negativas: imagem pesada para o desenvolvimento (792 MB); endurecimento obrigatório que, se esquecido, abre um caminho de exclusão sem credencial, e por isso o teste de integração exige 401 no filer; o Object Lock do SeaweedFS é recente (correções de retenção e de exclusão em setembro e outubro de 2026, issues 11335, 11492, 12006 e 12008 do projeto); governança de um mantenedor principal que vende edição comercial.
- **Risco residual: bloqueio de objeto autohospedado não resiste a quem administra o host.** No ensaio, `weed shell` apagou o bucket com a versão em conformidade; o mesmo vale para RustFS e Versity, bastando apagar o disco. Na nuvem gerenciada, o modo de conformidade resiste à própria conta root. Para o selo de RF-ADM-011, autohospedado, a garantia "guardado fora do banco, imutável" protege contra defeito do código e contra credencial S3 vazada, mas não contra o operador. **(Emenda: em produção o selo fica no S3 com Object Lock em conformidade, que resiste à conta root; o limite vale só para o desenvolvimento, numa máquina e sem dado real, e é aceito.)**
- Gatilhos de revisão: (a) versão nova do SeaweedFS reprovada no ensaio, sem correção em 30 dias → RustFS, qualificado pelo mesmo ensaio; (b) mudança de licença, fim da edição comunitária ou 6 meses sem release → o mesmo; (c) produção autohospedada acima de 50 TB, ou exigência de três locais com erasure coding **(proposta: limiar)** → avaliar Ceph RGW num ADR novo; (d) acionado pela emenda: produção no S3 gerenciado, SeaweedFS só no desenvolvimento; se a produção um dia sair da nuvem gerenciada, o SeaweedFS endurecido e a questão da H10 voltam; (e) componente da auditoria que passe a licença "source-available" → sai pela decisão 1, com ADR se for estrutural.

---

### ADR-013 — País e idioma: perfil de país no núcleo, textos por catálogo e conteúdo legal por país e jurisdição

**Status:** Aceita (dono do produto, 13/11/2026); Partes F e G acrescentadas na versão 0.9 (etapa E01) · **Data:** 13/11/2026 · Etapa [014](etapas/014_2026-11-13_multipais-e-idiomas.md). Depende das pendências P12 a P15 do sumário (seção 6); onde depende, está marcado.

**Nota de aceite (13/11/2026, mesma versão 0.8).** O dono do produto aceitou o ADR e respondeu às perguntas da seção 16: **Q1**, um país por empresa no MVP, e a unidade não tem país próprio (`unidade.pais` sai do MVP, com o gatilho (h)); **Q3**, uma moeda por empresa; **Q2**, vale a regra proposta, terminando em `pt-BR` (decidida por padrão, salvo objeção). Acrescentou o prefixo do número da OS pelo país da empresa (Parte A, item 9). A migração da Parte A, item 4, passa a ser a `00009`, porque a `00008` já existe (`00008_auditoria_insercao_por_coluna.sql`, trava da trilha). O texto abaixo já está com essas mudanças.

**Contexto.** O sumário 0.8 (seção 3.6, linha "Países e idiomas") fez a Tenaz **multipaís desde o MVP** (Brasil, Argentina e Paraguai) e **bilíngue desde a E01** (`pt-BR` e `es-419`) em telas, e-mails, notificações, documentos gerados e textos dos erros da API, com o idioma escolhido por usuário, o do navegador ou do aparelho na primeira vez e um padrão da empresa. O que existe hoje supõe o Brasil em todo lugar: `empresa.cnpj char(14)` com `ck_empresa_cnpj` só de dígitos, `unidade.cnpj`, `pessoa.cpf char(11)`, fuso padrão `America/Sao_Paulo` e lista de fusos só brasileira na web, dinheiro com "moeda implícita BRL" (documento 3, seção 2.5), textos de erro em português escritos no código Go (`erros.Novo(codigo, titulo, status)` e `.Com(detalhe)`), e-mails com assunto e corpo em português no código (`autenticacao/fluxos_email.go`), `pacotes/i18n` só com `pt-BR`, e a NR-12 citada pelo nome em tabelas e colunas do núcleo (`registro_nr12`, `ativo.sujeito_nr12`, `ordem_servico.exige_nr12`, motivo `nr12_inapto`). As migrações `00001` a `00007` já existem e não se editam. Nada foi implantado em homologação nem em produção, e nenhum integrador externo usa `/v1` (o `oasdiff` ainda não está na esteira, seção 2.5).

A pesquisa do `tenaz-conformidade-mercosul` ([mercosul_argentina.md](pesquisa/mercosul_argentina.md) e [mercosul_paraguai.md](pesquisa/mercosul_paraguai.md), 13/11/2026; pesquisa, não parecer) trouxe sete fatos que moldam o desenho: (1) prazos legais variam por **jurisdição**, não só por país: na Argentina, o prazo de inspeção de vaso de pressão é provincial (exemplo: Resolução SPA 231/96 da província de Buenos Aires) e a revisão técnica de veículo depende de adesão provincial (Santa Fe não aderiu à reforma de 2025); no Paraguai, a ITV está em disputa legislativa e a licença é municipal; (2) por isso, tipos de documento, normas e capacitações devem nascer como **catálogo por país, versionado e sem prazo embutido no código**, com os prazos legais preenchidos pelo cliente até haver parecer externo; (3) identificação: CNPJ e CPF (Brasil), CUIT, CUIL e DNI (Argentina), RUC e CI (Paraguai); a regra oficial do dígito verificador do CUIT (resultado 10) e do RUC (RUC antigo com letra) **não foi confirmada**; (4) a Tenaz não emite fatura: guarda a identificação fiscal e, quando houver integração, a referência ao documento fiscal do país (remito e fatura com CAE na Argentina; NRE, KuDE e CDC do SIFEN no Paraguai), sem CFOP fora do Brasil; (5) o guarani (`PYG`) não tem casas decimais; fusos `America/Argentina/Buenos_Aires` e `America/Asuncion`, com o Paraguai em UTC−3 permanente desde outubro de 2024 (Ley 7354/2024, número vindo de imprensa, a conferir); (6) o Brasil **não está** na lista de países adequados da autoridade argentina (AAIP): hospedar em São Paulo dados de empresa argentina é transferência permitida com garantias (cláusulas modelo da Resolução AAIP 198/2023 ou da Disposição 60/2016), com a Tenaz como "encarregada" (art. 25 da Ley 25.326); no Paraguai, a Ley 7593/2025 entra em vigor por volta de 27/11/2027 e o art. 19 exige garantias para a transferência; (7) a regra do bloqueio de máquina em reparo, **chave ou letreiro sob a guarda de quem repara**, retirado só por quem o colocou, com dispositivo de uso múltiplo para vários grupos, vale nos dois países (Decreto 351/79, art. 109; Decreto 14.390/92, art. 122), como já vale no Brasil (NR-12, item 12.11.3).

O risco é duplo: espalhar `if pais == "AR"` pelo código, ou, no outro extremo, abrir a configuração livre que o ADR-010 recusou (R7). E, como em todo o resto, nada disso pode tocar o isolamento (princípio 1).

**Decisão.**

**Parte A — País**

1. **`empresa.pais`**, `char(2)`, ISO 3166-1 alfa-2, obrigatório, num catálogo fechado: `BR`, `AR`, `PY`. Escolhido na criação da conta e **não alterável pela API**: trocar o país muda identificação, moeda, conteúdo legal e o aditivo de transferência de dados (P14), e vira operação assistida do suporte. **Um país por empresa no MVP** (Q1): a unidade não tem país próprio, e tudo o que depende de país (identificação da unidade, fusos aceitos, feriados nacionais, conteúdo legal) vem do país da empresa. Uma empresa com unidades em dois países tem uma conta por país. `unidade.pais` fica fora do MVP; entra só pelo gatilho (h).
2. **Perfil de país no código, não configurável.** Pacote `plataforma/pais` em Go e a mesma tabela em `pacotes/dominio` (TypeScript), com casos JSON compartilhados em `backend/testes/casos/` (seção 7). Para cada país: tipos de identificação aceitos e o validador de cada um, regra de "mesma entidade" para a conexão entre contas (item 4), moeda e casas decimais, fuso padrão e lista de fusos válidos, idioma padrão sugerido, tabela de formatação (separadores, padrão de data, símbolo da moeda) e feriados nacionais por ano. Isso é **núcleo**, como os estados (ADR-010, decisão 1): país novo entra por requisito, código e teste, nunca por configuração. Nenhum módulo testa o país por conta própria: pergunta ao perfil (`pais.Perfil(codigo)`). *(Estado na E01, versão 0.9, conferido no código: `backend/internal/plataforma/pais/` implanta `Perfil(codigo)` → `PerfilPais` com moeda, idioma padrão, fuso padrão, fusos válidos do país (`FusoValido`), tipo de identificação de empresa e tipos de pessoa; `Normalizar`, `ValidarIdentificacao` e `Mascarar`, com o CNPJ alfanumérico pelo algoritmo oficial; casos em `backend/testes/casos/identificacao.json`, por ora lidos só pela suíte Go. Faltam no perfil a formatação e os feriados nacionais; falta em `pacotes/dominio` a mesma tabela, que hoje valida só CPF e CNPJ (`documentos.ts`) e não lê os casos. As casas da moeda ficam em `plataforma/dinheiro`, Parte F.)*
3. **Identificação genérica: tipo e valor.** Em todo cadastro com identificação, o par `identificacao_tipo` e `identificacao`, este último normalizado (maiúsculas, só letras e dígitos, sem ponto, barra nem hífen; a máscara é só de exibição). Tipos, num catálogo fechado do perfil de país:

   | País | Empresa e unidade | Pessoa |
   |---|---|---|
   | BR | `cnpj` | `cpf` |
   | AR | `cuit` | `cuil`, `dni` |
   | PY | `ruc` | `ci` |

   - **Quem aceita o quê.** A empresa e a unidade aceitam só os tipos de pessoa jurídica do próprio país. A pessoa aceita os tipos de pessoa de **qualquer** país suportado, porque uma transportadora brasileira pode ter motorista paraguaio. Fornecedor, locatário e cliente atendido, quando entrarem nas suas etapas, aceitam qualquer tipo de qualquer país (cliente brasileiro com subcontratada paraguaia). Tipos a mais (CUIT de pessoa com atividade, passaporte de estrangeiro) entram só com a resposta da P13.
   - **Validação em duas camadas.** O banco tem `CHECK` de forma por tipo (tamanho e alfabeto) e de coerência (os dois preenchidos ou os dois nulos); o dígito verificador é conferido no Go e no TypeScript, com os mesmos casos JSON. **CNPJ e CPF** têm o dígito conferido como hoje. **CUIT e RUC: só a forma, sem dígito verificador**, até a regra oficial ser confirmada (P13; dúvidas do resultado 10 no CUIT e do RUC antigo com letra); o validador de dígito entra depois por requisito, com teste, sem migração. **DNI e CI** não têm dígito verificador. Nenhum validador consulta serviço externo (ARCA, DNIT, Receita): o dígito pega erro de digitação, não prova que o número existe.
   - **CNPJ alfanumérico.** A Receita Federal passou a atribuir CNPJ com letras nas 12 primeiras posições a partir de julho de 2026 (IN RFB 2.229/2024; a confirmar pelo `tenaz-conformidade-br`). O `ck_empresa_cnpj` atual (`^[0-9]{14}$`) recusaria uma empresa aberta agora; o `CHECK` novo já nasce aceitando `[0-9A-Z]{12}[0-9]{2}`, e o validador de dígito do CNPJ passa a calcular com letras.
   - **Unicidade.** Empresa: única em toda a plataforma por `(identificacao_tipo, identificacao)`, como o `uq_empresa_cnpj` de hoje. Pessoa: única por `(empresa_id, identificacao_tipo, identificacao)`.
   - **Dado sensível.** A classificação interna do sumário 3.6 ("CPF completo") passa a valer para a identificação pessoal completa de qualquer tipo (CPF, CUIL, DNI, CI): exibição mascarada, visualização gravada na trilha com `sensivel`. `cpf_mascarado` da API dá lugar a `identificacao_mascarada`.
4. **O que fazer com as colunas que já existem.** Uma migração nova (a próxima da sequência, `00009`, porque a `00008` é a trava da trilha; o nome é do `tenaz-banco-dados`), sem editar `00001` a `00007`: acrescenta `pais`, `idioma_padrao`, `moeda`, `identificacao_tipo` e `identificacao` em `empresa`; `identificacao_tipo` e `identificacao` em `unidade`; `identificacao_tipo` e `identificacao` em `pessoa`; `idioma` em `usuario`; e `moeda` ao lado de `usuario.custo_hora` (item 5). Copia `cnpj` para (`cnpj`, valor) e `cpf` para (`cpf`, valor), recria as unicidades e **remove na mesma migração** `cnpj`, `cpf`, `uq_empresa_cnpj`, `uq_pessoa_cpf` e os três `CHECK` antigos. Expandir e contrair numa migração só é exceção à regra da seção 9.4 (coluna removida só no release seguinte), justificada porque nenhuma versão antiga do código roda em ambiente nenhum; desde a primeira implantação em homologação, a regra volta a valer sem exceção. Na API, `/v1/empresas`, unidades e pessoas trocam `cnpj`, `cpf` e `cpf_mascarado` pelos campos novos **agora**, antes do primeiro integrador: depois disso, seria mudança incompatível em `/v1`. A conexão entre contas (`plataforma.conexao_empresa.documento_convidado`, ainda não migrada) compara a "mesma entidade" pela regra do perfil de país: raiz de 8 posições no CNPJ, como hoje; o CUIT e o RUC inteiros, porque na Argentina e no Paraguai a filial usa a mesma identificação da matriz. Os nomes das colunas dessa tabela e de `fornecedor`, `locatario` e `cliente_servico` seguem o mesmo par e são fixados nas etapas delas.
5. **Moeda e dinheiro.**
   - **`empresa.moeda`**, `char(3)`, ISO 4217: `BRL`, `ARS` ou `PYG`, dada pelo país na criação e não alterável pela API.
   - **Dinheiro é sempre `numeric` e sempre tem a moeda ao lado**: toda tabela com valor em dinheiro tem a coluna `moeda char(3) NOT NULL`, uma por linha, que vale para todos os valores daquela linha. Nunca `float`; nunca inteiro em centavos: o guarani não tem casas decimais, e a convenção de centavos da seção 8.1 divergia dos documentos 2 e 3 (corrigida aqui).
   - **No MVP, uma moeda por empresa** (decidido pelo dono do produto, Q3)**:** a `moeda` de cada linha é igual a `empresa.moeda`, conferido no serviço e por teste. Indicadores somam uma moeda só. Multimoeda com cotação (compra em dólar, unidade em outro país) fica fora, com gatilho (e).
   - **Casas decimais da moeda:** `BRL` e `ARS` com duas, `PYG` com zero; o arredondamento à casa da moeda é uma função do perfil de país, a mesma no Go e no TypeScript, com casos compartilhados. O banco guarda `numeric` com a escala da coluna, e um valor em `PYG` sempre tem a parte decimal zerada. *(Versão 0.9: a casa é da moeda, não do país, e a regra vale para o valor, não para o preço unitário; ver a Parte F.)*
   - **Largura:** o `numeric(14,2)` dos totais (documento 2, seção 1.4: "até R$ 999 bilhões") deve ser revisto pela escala do guarani e pela inflação do peso argentino; proposta `numeric(18,2)` nos totais. A decisão é do `tenaz-banco-dados`, e agora é barata, porque a única coluna de dinheiro já migrada é `usuario.custo_hora`. *(Versão 0.9: substituído pela Parte F, que fixa a largura aqui, porque o limite do tipo Go depende dela.)*
   - **Na API:** valor em texto decimal (documento 3, seção 2.5), com tantas casas quantas a moeda tiver (`"1250.40"` em BRL, `"75000"` em PYG), e todo recurso que tem dinheiro traz `moeda`. No Go, um tipo da plataforma com valor decimal e moeda, sem `float`; a biblioteca decimal é escolhida pela regra do ADR-012 (candidatas: `shopspring/decimal`, MIT, ou `cockroachdb/apd`, Apache-2.0) **(proposta)**. *(Versão 0.9: substituído pela Parte F: nenhuma biblioteca decimal; regras de entrada e saída do texto.)*
6. **Fuso.** Sempre nome IANA, nunca deslocamento fixo. Padrão pelo país: `America/Sao_Paulo`, `America/Argentina/Buenos_Aires`, `America/Asuncion`. `empresa.fuso_horario` e `unidade.fuso_horario` (já existem) passam a ser validados contra a lista do país da empresa. As regras de horário vêm da base tz, nunca de tabela própria: o binário Go embute `time/tzdata`, para não depender da imagem; o RDS e o sistema do aparelho trazem a sua. A mudança do Paraguai em 2024 é o exemplo do motivo: sistemas com base antiga mostraram a hora errada. Limites de dia (vencimento, preventiva, indicador diário) continuam calculados no fuso da unidade (documento 3, seção 2.4).
7. **Formatos.** Números, datas e moeda são formatados por funções próprias de `pacotes/dominio` e do Go (para os PDFs), com a tabela do perfil de país e casos compartilhados, e **nunca** pelo `Intl` com a etiqueta `es-419`: o padrão regional `es-419` do CLDR usa ponto decimal, e os três países usam vírgula decimal, ponto de milhar e data `dd/mm/aaaa`. A formatação usa o idioma do usuário com o país da empresa. Além disso, o `Intl` do Hermes, motor do app, é incompleto. A moeda aparece com o símbolo do país quando é a moeda do país, e com o código ISO quando não é (o `$` é o símbolo do peso argentino). O desenho de cada formato é do documento 7. *(Versão 0.9: mantido e detalhado na Parte G, depois de a web ter sido feita com `Intl` pela localidade do país.)*
8. **Feriados nacionais** ficam no perfil de país, por ano, e o serviço os grava em `feriado_unidade` segundo o país da empresa, como já faz com os do Brasil (RF-ADM-002). Feriado transferido por decreto (comum na Argentina) chega por release; o calendário da unidade continua editável pela empresa.
9. **Prefixo do número da OS pelo país da empresa** (decisão do dono do produto no aceite): `OS-` no Brasil, `OT-` (*orden de trabajo*) na Argentina e no Paraguai, como em `OT-2026-004812`. O número é gravado na criação da OS e **nunca** muda com o idioma de quem lê: um usuário `pt-BR` numa empresa argentina vê `OT-`. O prefixo vem do perfil de país, não do catálogo de textos. Os números legíveis de outros documentos seguem a mesma regra quando entrarem.

**Parte B — Idioma**

1. **Idiomas suportados**, num catálogo fechado de etiquetas BCP 47: `pt-BR` e `es-419`. `pt-BR` é o idioma de origem e o de último recurso.
2. **Onde fica a escolha.** `empresa.idioma_padrao` é obrigatório, sugerido pelo país (`pt-BR` no Brasil, `es-419` na Argentina e no Paraguai) e alterável pelo `administrador`. `usuario.idioma` é obrigatório: no convite recebe o padrão da empresa; no aceite do convite (o primeiro acesso), a tela já vem no idioma do navegador ou do aparelho, se for suportado, e grava o idioma em que a pessoa aceitou; depois, só o próprio usuário muda, no perfil. O responsável que cria a conta fica com o idioma da tela de cadastro. O idioma é do **usuário em cada empresa**, não da identidade: as tabelas de credencial não mudam (ADR-011), e o dono do dado é o ADM.
3. **Negociação na API**, nesta ordem: (a) `Accept-Language` (RFC 9110, seção 12.5.4, com pesos `q`): vale a primeira etiqueta aceitável que casa com um idioma suportado, com `pt` e `pt-*` levando a `pt-BR` e `es` e `es-*` levando a `es-419`; (b) sem casamento, o `usuario.idioma` da sessão; (c) token de API ou rota sem usuário, `empresa.idioma_padrao`; (d) `pt-BR`. A web e o app **sempre** enviam `Accept-Language` com o idioma do usuário. Toda resposta que leva texto para pessoa (erro, documento gerado) traz `Content-Language` e `Vary: Accept-Language`.
4. **Recursos JSON não levam texto traduzido**: só códigos e o que a empresa digitou. Rótulos (de indicadores, de estados, de permissões) são traduzidos pelo cliente a partir do código. Com isso, o `ETag` (versão do recurso) não depende do idioma, e o cache continua correto.
5. **Erros.** `codigo`, `type` e os membros adicionais não mudam com o idioma; `title`, `detail` e `erros[].mensagem` saem do catálogo do servidor no idioma negociado. Em Go, o título passa a vir do catálogo pela chave do código, e o detalhe livre de `.Com(detalhe)` dá lugar a uma chave com parâmetros; a camada HTTP renderiza no idioma da requisição. Texto para pessoa escrito no código Go passa a ser defeito, apanhado por teste. Logs registram código e chave, nunca o texto traduzido. A web e o app continuam mostrando o próprio texto pelo `codigo` (como `erro.codigos` de `pacotes/i18n` já faz) e usam o `detail` do servidor quando o código não está no catálogo da interface.
6. **Texto do sistema e texto da empresa.** O que o sistema cria a partir de catálogo (os 12 perfis padrão, o catálogo de permissões, as telas da trilha, os tipos de operação, as situações, os modelos de turno) é mostrado pela chave do código, no idioma de quem lê; o nome gravado só vale quando a empresa renomeia. O que a empresa ou o usuário digita (descrição de OS, nome de ativo, motivo) **nunca é traduzido**, nem por máquina. Conteúdo copiado de um catálogo para a empresa (Parte C, pacotes de segmento do ADR-010) é copiado no `idioma_padrao` da empresa e vira dado dela (cópia, não referência).
7. **E-mails, notificações e documentos no idioma do destinatário.**
   - E-mail disparado pela própria pessoa (esqueci a senha, verificação de e-mail): idioma negociado na requisição que o disparou, porque a identidade é global e não guarda idioma.
   - E-mail ou notificação disparado por outra pessoa ou por job (convite, alerta de vencimento, aviso de OS, push): `usuario.idioma` do destinatário.
   - Destinatário sem usuário (contato de fornecedor, cliente final, prestador do link): `empresa.idioma_padrao` da empresa que envia. A página do link de uso único usa o idioma do navegador.
   - Documento gerado (OS impressa, relatório em PDF, exportação CSV ou XLSX, extração da trilha): idioma de quem pede. O processamento assíncrono (documento 3, seção 2.14) grava o idioma no pedido, e todo job que gera texto recebe o idioma **já resolvido** no enfileiramento; o worker não adivinha.
   - O idioma nunca entra no conteúdo com hash da trilha: `evento_auditoria` guarda códigos e dados, e a mesma trilha é verificada (`integro`, `violado`) em qualquer idioma.
8. **Catálogos de texto e fonte da verdade.** Dois catálogos, separados por quem renderiza, com **uma só convenção**: chaves estáveis em JSON aninhado, interpolação `{{nome}}` e plural por sufixo CLDR (`_one`, `_other`). A forma escolhida segue a regra CLDR de cada idioma, que difere no zero (em português, zero usa `_one`; em espanhol, `_other`).
   - **Interface (web e app):** `pacotes/i18n`. `pt-BR.ts` é a fonte e define o tipo `Textos`; o `es-419.ts` é declarado com `satisfies Textos`, então chave faltando ou sobrando quebra a compilação. Os textos do app ficam no bundle (o app troca de idioma sem conexão), e o idioma do aparelho vem do `expo-localization` (MIT; dependência nova, admitida pela regra do ADR-012).
   - **Servidor (erros, e-mails, notificações, documentos):** um arquivo por idioma em `backend/internal/plataforma/textos/`, embutido no binário com `embed`. O pacote `plataforma/textos` faz busca, interpolação e plural com `golang.org/x/text/feature/plural` (BSD-3, já dependência indireta), sem biblioteca nova.
   - **A fonte da verdade de cada texto é o arquivo `pt-BR` do seu catálogo**; o `es-419` é tradução dele. O vocabulário de domínio em português continua sendo o do sumário; o glossário bilíngue e a revisão das traduções são do `tenaz-idiomas`, e todo PR que altera um `es-419` passa pela revisão dele.
   - **Verificação no CI:** tarefa `textos` em `ferramentas/tarefas.sh`, chamada por `conferir`, que quebra o build se (a) alguma chave faltar ou sobrar em qualquer idioma, nos dois catálogos; (b) as variáveis `{{...}}` de uma chave diferirem entre os idiomas; (c) faltar forma de plural exigida pelo idioma; (d) um código de erro registrado no Go ou um modelo de e-mail não tiver texto nos dois idiomas; (e) um texto usar termo proibido do glossário daquele idioma, como já faz o teste de tom de voz de `pt-BR.test.ts`, que passa a valer por idioma. Responsável: `tenaz-plataforma-sre`.
   - **Estado na E01 (versão 0.9, conferido no código).** Implantado pelo `tenaz-backend-go`: `backend/internal/plataforma/textos/` com `catalogos/pt-BR.json` e `catalogos/es-419.json` embutidos, `textos.T(idioma, chave, params)` (cai para `pt-BR`; sem texto, devolve a chave), `textos.Plural` e a negociação em `textos.Preferencias.Resolver()`. A tarefa `textos` roda os testes Go `TestConferirCatalogos` (itens (a), (b), (c) e (e) no catálogo do servidor), `TestTodoCodigoTemTitulo` (item (d), sobre `erros.Todos()`) e `TestTodaTelaTemNome` (telas da trilha, `tela.<codigo>`); os modelos de e-mail estão em `email.*` e são cobertos pelo item (a). No catálogo da interface, os mesmos itens são conferidos pela checagem de tipos e por `es-419.test.ts` (seção 8.4). O texto livre do Go foi todo passado ao catálogo, e a mesma tarefa roda a trava contra texto livre (`TestSemTextoLivreNoCodigoDeProducao` e `TestTravaAcusaTextoLivre`, seção 8.2).

**Parte C — Conteúdo legal por país e por jurisdição**

1. **Três camadas, nenhuma com prazo legal no código.** (a) O **perfil de país** (Parte A, item 2), núcleo. (b) Os **modelos de registro de conformidade**, catálogo fechado no código: só a **forma** de um registro cujos campos a lei fixa, como os sete campos do item 12.11.2 da NR-12. (c) O **conteúdo de país**: tipos de documento, normas, capacitações, prazos, antecedências de alerta e modelos de checklist legais, como **dado versionado no repositório e copiado para a empresa**. Nenhum prazo legal é constante no Go nem no TypeScript, em nenhum país, inclusive no Brasil.
2. **Jurisdição como atributo.** O conteúdo de país e os prazos ganham a jurisdição, em código ISO 3166-1 alfa-2 (nacional, ex.: `AR`) ou ISO 3166-2 (subdivisão: estado, província ou departamento, ex.: `BR-SP`, `AR-B`). O documento registrado ganha a **jurisdição emissora**. O prazo de um tipo de documento é resolvido pela jurisdição emissora do documento: prazo da subdivisão, se houver; senão o nacional; senão nenhum, e então a validade é digitada no documento. Município não é jurisdição de prazo no MVP (gatilho (f)). O tipo de documento também tem país, para uma empresa brasileira ter a licença paraguaia do seu motorista paraguaio.
3. **Prazo legal preenchido pelo cliente até haver parecer.** No pacote de país, o prazo de um item nasce **vazio** quando a pesquisa não o confirmou em fonte oficial ou quando não há parecer externo (P12): a empresa preenche, e a tela diz que o prazo é dela. Prazo confirmado por parecer entra no pacote como dado, com a fonte, numa versão nova. Regras de condição e efeito do pacote que dependem de prazo nascem `inativa`, como as da P9 (ADR-010). O Brasil segue a mesma regra: prazos hoje fixos em `CHECK` ou citados em comentários de tabela (por exemplo `ck_tipo_documento_civ_cipp`, `ck_tipo_documento_mopp` e a `regra_prazo_nr13`) passam a ser dado do pacote `BR`.
4. **Pacote de país na mecânica do ADR-010.** O conteúdo de país usa o formato e a cópia idempotente dos pacotes de segmento (ADR-010, decisão 6; seção 5.6): arquivo versionado em `backend/internal/modulos/configuracao/pacotes/`, validado no CI, aplicado pela função `Garantir...` do módulo dono, criando só o que falta e nunca alterando nem apagando. Diferenças: o pacote de país é aplicado **automaticamente** na criação da empresa; é lido do binário, sem a tabela global `pacote_segmento`, que continua na Fase 2; e, por isso, a cópia idempotente chega ao MVP só para ele. A aplicação é uma reação do pacote `configuracao` ao evento de criação da empresa, então **nenhuma dependência nova entre módulos** (seção 5.3: `configuracao` já depende de ADM, CON, CHK, PRE e os outros donos; ADM não importa ninguém). Os feriados nacionais ficam fora do pacote: são do perfil de país, gravados pelo próprio ADM (Parte A, item 8). Os pacotes de segmento, na Fase 2, ganham a lista de países em que podem ser aplicados.
5. **O núcleo não cita norma pelo nome.** OS, ATV, o app e o sincronizador falam de "conformidade": a OS pergunta a CON quais registros de conformidade o ativo exige e se o registro está completo. CON é dono do catálogo de normas do pacote de país e dos modelos de registro do código. Um registro genérico, `registro_conformidade`, substitui `registro_nr12`: colunas comuns tipadas (OS, norma, versão da norma, indicação conclusiva, restrições, responsável, registro profissional, assinatura) e os campos próprios do modelo em `jsonb`, validados pela mesma função no Go e em `pacotes/dominio`, com casos JSON compartilhados, como os campos extras (seção 5.6). O ativo e a classe passam a marcar as normas aplicáveis, e a OS copia as exigidas na abertura, no lugar de `sujeito_nr12`, `sujeito_nr12_padrao` e `exige_nr12`. O bloqueio por parecer `inapto` usa um motivo genérico com a norma ao lado, no lugar de `nr12_inapto`. Uma norma só se aplica a ativo de empresa do país dela. **A NR-12 é o primeiro modelo de registro, com os mesmos sete campos e as mesmas regras**: muda onde ela mora, não o que ela exige. Modelo de registro novo entra por requisito, código e teste (R7). Na Argentina e no Paraguai a pesquisa não achou registro com campos fixados em lei além do "registro periódico" de caldeira do Paraguai (art. 102, inc. 19), que cabe no histórico de OS e em documento; até a P12, esses países nascem sem modelo de registro próprio.
6. **Bloqueio de energia é núcleo comum, não regra de norma.** A regra "o dispositivo de bloqueio (cadeado, etiqueta ou letreiro) fica sob a guarda de quem repara e só é retirado por quem o colocou, com uso múltiplo para vários grupos" vale nos três países. Os campos de bloqueio saem do registro da norma e ficam no registro comum da OS. No LOTO da Fase 2 (`bloqueio_energia`, `bloqueio_ponto`), a retirada de um dispositivo exige `removido_por` igual a `aplicado_por`, e um mesmo ponto aceita vários dispositivos ao mesmo tempo, um por pessoa: o `uq_bloqueio_ponto` atual impede isso e muda. A regra fica no Go e no trigger de transição, como a de hoje (OS com bloqueio aplicado não vai para `concluida`).
7. **Documento fiscal: só referência.** A Tenaz não emite fatura nem remessa em país nenhum. Onde hoje há `nfe_numero`, `nfe_chave` e `cfop` (EST, P8), o desenho passa a tipo e referência do documento fiscal do país: NF-e no Brasil; remito e fatura com CAE na Argentina; NRE e KuDE com CDC no Paraguai. O CFOP existe só quando o documento é NF-e. Os nomes e o catálogo de tipos são fixados na etapa de EST, com a P8 e a P13.

**Parte D — Hospedagem**

1. **Uma região, uma célula, três países.** A AWS `sa-east-1` (ADR-006) atende Brasil, Argentina e Paraguai no mesmo pool (ADR-002), na célula `pool-01` (`empresa.celula`, documento 2). O princípio 5 continua: os dados primários ficam no Brasil, para os três países. Para uma empresa argentina ou paraguaia, isso é **transferência internacional**: na Argentina, permitida com cláusulas modelo, com a Tenaz como encarregada; no Paraguai, pelo art. 19 da Ley 7593/2025, que entra em vigor por volta de 27/11/2027. O contrato com o cliente precisa de um **aditivo de transferência por país**, com a Tenaz como encarregada (P14, decisão jurídica). O `empresa.pais` é o dado que diz qual aditivo vale para cada conta. A conexão entre contas de países diferentes (cliente brasileiro com subcontratada paraguaia) roda no mesmo banco e não muda tecnicamente; o termo da conexão segue a P14.
2. **Sinais por país.** Logs, métricas e traços ganham o atributo `pais` (três valores, sem problema de cardinalidade, seção 14.2), para medir latência e erro por país.

**Parte E — Impacto no que já existe (E01)**

| Área | O que muda |
|---|---|
| `empresa` | `pais`, `idioma_padrao`, `moeda`, `identificacao_tipo` e `identificacao` no lugar de `cnpj`; fuso padrão pelo país; o CNPJ aceita o formato alfanumérico. A regra do sumário 3.6 "CNPJ da empresa" (o `administrador` altera, com registro na trilha) passa a valer para a identificação |
| `unidade` | Identificação no lugar de `cnpj`, de tipo do país da empresa; fuso validado e feriados nacionais pelo país da empresa; sem país próprio no MVP (Q1) |
| `pessoa` | Identificação no lugar de `cpf`, com tipos de qualquer país; `cpf_mascarado` vira `identificacao_mascarada`; `registro_profissional` continua texto livre |
| `usuario` | `idioma`; `moeda` ao lado de `custo_hora`. *(Versão 0.9: `custo_hora` passa a preço unitário, `numeric(18,6)`, Parte F, item 9.)* |
| `perfil` | Nome dos perfis padrão mostrado pela chave do código; `perfil.nome` só para o perfil renomeado ou personalizado |
| Criação da conta (`/v1/empresas`) | Recebe país, identificação e idioma; o país define moeda, fuso e idioma sugeridos e dispara o pacote de país |
| Autenticação | E-mails de convite, redefinição e verificação no idioma da Parte B, item 7; nenhuma tabela de credencial muda; tela de login no idioma do navegador até haver sessão |
| Auditoria | Nenhuma coluna muda; extração e impressão da trilha no idioma de quem pede; o hash não depende de idioma |
| Erros (`plataforma/erros`) | Título e detalhe pelo catálogo do servidor; os textos de E01 ganham `es-419` |
| Web | `es-419.ts`; detecção do idioma do navegador antes do login; seletor de idioma no aceite do convite e no perfil; `Accept-Language` em toda chamada; fusos da Argentina e do Paraguai; formatação própria por país; textos que citam CPF (por exemplo `erro.codigos.conflito`) passam a falar de "documento" |
| App | `expo-localization`; idioma por conta no banco local; textos no bundle |

**Alternativas consideradas.**

| Alternativa | A favor | Contra | Resultado |
|---|---|---|---|
| **Uma coluna por documento de cada país** (`cnpj`, `cuit`, `ruc`, `cpf`, `dni`, `ci`...) | Tipo e `CHECK` nativos por coluna | Colunas quase sempre nulas; país novo vira migração em toda tabela com identificação; unicidade e busca espalhadas; a conexão entre contas teria uma regra por coluna | Descartada |
| **Manter `cnpj` e `cpf` e acrescentar o par só para os outros países** | Nenhuma mudança no que funciona | Dois caminhos para o mesmo dado, para sempre; a regra de "mesma entidade" e a unicidade teriam dois formatos | Descartada |
| **Tabela de identificações 1:N por cadastro** | Pessoa com DNI e CUIL ao mesmo tempo | Junção em toda listagem e busca; nenhum caso do MVP pede mais de uma; a unicidade da empresa fica mais cara | Descartada; gatilho (d) |
| **Expandir e contrair em duas migrações agora** | Segue a regra da seção 9.4 ao pé da letra | Não há versão antiga rodando em lugar nenhum: a segunda migração só adiaria a remoção, com o código lendo dois formatos até lá | Descartada; a regra volta na primeira implantação |
| **Validar o dígito do CUIT e do RUC já** com o algoritmo de fontes secundárias | Pega erro de digitação desde o primeiro dia | Se a regra estiver errada (resultado 10 do CUIT, RUC antigo com letra), recusa uma empresa real e trava o cadastro dela | Descartada até a P13 |
| **Idioma na identidade** (`plataforma.identidade`) | Uma preferência por pessoa, em todas as empresas | A identidade só nasce no aceite, e o convite precisa de idioma antes; as tabelas de credencial ficam só com o `tenaz_autenticacao` (ADR-011), e o perfil do usuário é do ADM | Descartada |
| **Idioma só da empresa** | Mais simples | Contraria a decisão do dono do produto (idioma por usuário) | Descartada |
| **Servidor só com códigos, tradução só no cliente** | Nenhum catálogo no Go | Contraria a decisão do dono do produto: e-mails, PDFs e o `title` e o `detail` dos erros saem no idioma; integradores não têm o catálogo da web | Descartada |
| **Um catálogo único para servidor e interface**, com o Go gerado a partir de `pacotes/i18n` | Um lugar para o tradutor | O binário Go passaria a depender de um pacote JavaScript e de um passo de geração; o app precisa dos seus textos sem servidor; quem renderiza é diferente. A convenção única já dá ao `tenaz-idiomas` um formato só | Descartada; gatilho (g) |
| **Biblioteca de i18n em Go** (`go-i18n`, `gotext`) | Pronta | Formato próprio, diferente do da interface; o necessário (busca, `{{...}}` e plural CLDR) cabe num pacote pequeno sobre `x/text` | Descartada |
| **Plataforma de tradução hospedada** (Crowdin, Lokalise e afins) | Fluxo de tradução pronto | Serviço pago ou com cota (ADR-012) e componente novo sem gatilho (princípio 2) | Descartada; gatilho (g) |
| **Tradução automática do conteúdo digitado pela empresa** | Usuário de outro idioma lê tudo | Texto técnico e legal traduzido sem revisão vira erro de manutenção; custo por chamada; dado da empresa enviado a terceiro | Descartada |
| **Espanhol por país** (`es-AR`, `es-PY`) | Vocabulário local | O dono do produto pediu espanhol latino-americano neutro; as diferenças de número e data se resolvem pela formatação por país | Descartada |
| **Formatar com `Intl` e a etiqueta do idioma** | Sem código próprio | `es-419` dá ponto decimal; o `Intl` do Hermes é incompleto; o PDF gerado no Go formataria diferente da tela | Descartada |
| **Dinheiro como inteiro em centavos** (seção 8.1 até a versão 0.7) | Sem biblioteca decimal | O guarani não tem centavos; divergia dos documentos 2 e 3 | Descartada |
| **Moeda só em `empresa.moeda`, sem coluna nas linhas** | Uma coluna a menos por tabela | Multimoeda futura exigiria descobrir a moeda de cada valor antigo; o integrador não vê a moeda junto do valor | Descartada |
| **Multimoeda com cotação no MVP** | Atende compra em dólar, comum na Argentina | Cotação, data de conversão e indicadores em várias moedas: escopo grande sem pedido de cliente; a P15 ainda decide a moeda de cobrança | Descartada; gatilho (e) |
| **Prazos legais no código**, por país, como constantes testadas | Prazo certo sem ação do cliente | A pesquisa mostrou prazos provinciais, normas instáveis (ITV paraguaia) e fontes conflitantes; cada mudança de lei viraria release; um prazo errado no código vira descumprimento atribuído à Tenaz | Descartada |
| **Uma tabela tipada por norma** (`registro_nr12`, uma para cada equivalente) | `NOT NULL` e `CHECK` nativos por campo | Cada norma nova traz tabela, política `escopo_externo`, pull do app e grupo na lista da conexão entre contas (R5) | Descartada: registro genérico com modelo validado no código |
| **Norma como campos extras da empresa** (ADR-010) | Nenhum código novo | Campo legal não pode ser apagado nem renomeado pela empresa; é exatamente o R7 | Descartada |
| **Região ou célula por país desde já** | Dado no país | A AWS não tem região na Argentina nem no Paraguai (a conferir na data); seria outro provedor e outra operação, sem gatilho; a pesquisa indica que a transferência é permitida com garantias | Descartada; gatilho (a) |

**Consequências.**

- Positivas: país novo é perfil de código mais pacote de conteúdo, sem `if` por país espalhado; o CNPJ alfanumérico entra sem migração de emergência; nenhum prazo legal fica preso ao release, e prazos provinciais cabem sem modelo novo; os textos têm fonte única por catálogo, com verificação automática dos dois idiomas; os recursos JSON e a trilha não dependem do idioma; o isolamento não muda (as colunas novas ficam em tabelas que já têm `empresa_id` e RLS, e nenhuma tabela global nova entra no MVP); nenhum componente novo na stack.
- Negativas: a migração `00009` mexe em `empresa`, `unidade`, `pessoa` e `usuario` e em todo o código de E01 que lê `cnpj` e `cpf`; todo texto de E01 precisa de tradução revisada antes do aceite da etapa; a empresa preenche prazos legais que o cliente talvez esperasse prontos (a tela precisa dizer isso com clareza, documento 7); o desenho de CON muda antes de ser migrado (o registro genérico perde o `NOT NULL` nativo nos campos da norma e depende de validação em Go, em TypeScript e de teste); sem dígito verificador do CUIT e do RUC, erro de digitação passa até a P13; formatação própria é código a manter.
- Gatilhos de revisão: (a) parecer da P14 que exija armazenamento no país ou vede a transferência para alguma categoria de dado, ou cliente ou contrato público que exija residência local → célula por país (mesmo código, outra região ou provedor) num ADR novo, usando `empresa.celula`; (b) p95 da API medido para empresas da Argentina ou do Paraguai acima da meta de RNF-008 por duas semanas → borda ou CDN antes de mover dados; (c) regra oficial do dígito verificador do CUIT e do RUC confirmada (P13) → validador por requisito; (d) caso real de cadastro que precise de duas identificações ao mesmo tempo → tabela 1:N; (e) primeiro cliente pagante com compra em outra moeda ou com unidade em país de outra moeda, ou decisão da P15 → multimoeda com cotação num ADR novo; (f) prazo legal que varie por município num cliente real → jurisdição municipal; (g) terceiro idioma ou mais de 3.000 chaves por catálogo → avaliar catálogo único gerado ou ferramenta de tradução autohospedada, de código aberto; (h) cliente pagante, ou contratação em negociação, que precise de unidade em outro país na mesma conta (por exemplo, indústria brasileira com fábrica no Paraguai) → `unidade.pais`, com identificação, fuso, feriados e conteúdo legal pela unidade, numa emenda deste ADR; a moeda continua a da empresa até o gatilho (e).

**Nomes canônicos novos.** O ADR foi aceito; os nomes entram no sumário, seção 3.6, pelo `tenaz-produto`, e só então os documentos donos os usam. Os marcados como "nome proposto" ainda podem mudar nesse passo.

| Nome | Onde | Valores ou forma |
|---|---|---|
| `empresa.pais` | ADM | ISO 3166-1 alfa-2: `BR`, `AR`, `PY` (`unidade.pais` fora do MVP, gatilho (h)) |
| Prefixo do número da OS pelo país da empresa | OS | `OS-` (BR), `OT-` (AR, PY); gravado na criação, nunca muda com o idioma |
| `empresa.idioma_padrao`, `usuario.idioma` | ADM | `pt-BR`, `es-419` |
| `empresa.moeda` e `moeda` ao lado de todo valor em dinheiro (primeira: `usuario.moeda`, de `custo_hora`) | Todos os módulos com dinheiro | ISO 4217: `BRL`, `ARS`, `PYG` |
| `identificacao_tipo`, `identificacao` (em `empresa`, `unidade`, `pessoa`; depois em `fornecedor`, `locatario`, `cliente_servico` e na conexão entre contas) | ADM; TER, LOC e SRV nas suas etapas | `cnpj`, `cpf`, `cuit`, `cuil`, `dni`, `ruc`, `ci` (lista final pela P13). Substituem `empresa.cnpj`, `unidade.cnpj`, `pessoa.cpf` e, nas etapas delas, `cnpj_cpf` e `documento_convidado` |
| `identificacao_mascarada` | API de pessoas | Substitui `cpf_mascarado` |
| Jurisdição: `jurisdicao` (no catálogo de tipos de documento e nos prazos) e `documento_vencimento.jurisdicao_emissora` | CON | ISO 3166-1 alfa-2 ou ISO 3166-2 |
| País do tipo de documento e do tipo de capacitação (`tipo_documento.pais`, `tipo_capacitacao.pais`) | CON | ISO 3166-1 alfa-2 |
| Prazo por jurisdição de um tipo de documento (tabela filha; nome proposto `prazo_tipo_documento`) | CON | Validade em meses por jurisdição, preenchida pela empresa ou pelo pacote com fonte |
| `registro_conformidade` (substitui `registro_nr12`), com a norma e a versão da norma | CON | Código da norma no formato `<pais>_<norma>`, minúsculo: `br_nr12`, `br_nr13` |
| Normas aplicáveis do ativo e da classe e normas exigidas na OS (substituem `sujeito_nr12`, `sujeito_nr12_padrao`, `exige_nr12`; nomes propostos `ativo.normas_aplicaveis`, `classe_ativo.normas_aplicaveis_padrao`, `ordem_servico.normas_exigidas`) | ATV, OS | Lista de códigos de norma |
| Motivo de bloqueio por parecer (substitui `nr12_inapto`; nome proposto `conformidade_inapto`) | CON | — |
| Evento de criação da empresa (nome proposto `empresa.criada`) | ADM | Formato `<entidade>.<fato>` da seção 5.5 |
| Documento fiscal por tipo e referência (substitui `nfe_numero`, `nfe_chave`) | EST | Nomes e tipos na etapa de EST, com a P8 e a P13 |

Ficam **fora** do sumário, porque são nomes de código deste documento: os pacotes `plataforma/pais`, `plataforma/textos` e o tipo de dinheiro da plataforma, e a tarefa `textos` de `ferramentas/tarefas.sh`. Nenhum código de erro novo: identificação, idioma, país, fuso ou moeda inválidos usam `validacao_falhou` com `formato_invalido` ou `valor_invalido` por campo. Permissões com norma no nome (por exemplo `con.gerenciar_nr13`, de RF-ADM-005) não mudam aqui: o catálogo de permissões é do dono do produto (sumário 3.6), e a pergunta fica com o `tenaz-produto`.

**Complemento de 13/11/2026 (etapa E01, versão 0.9 deste documento): Parte F — Representação do dinheiro.** Decisão técnica do `tenaz-arquiteto`, sem escolha do dono do produto: não muda o que a Tenaz faz, só como o número é guardado e trafega.

*Contexto.* A Parte A, item 5, fixou `numeric`, moeda ao lado e texto decimal na API, mas deixou em aberto a biblioteca, a largura e o tratamento por moeda. O código de E01 ficou no meio do caminho: `usuario.custo_hora` é `numeric(12,2)` no banco, mas o Go o converte para centavos inteiros (`api.Centavos`, `api.Dinheiro`, `adm.centavos`) e o SQL divide por `100.0` depois de um `::bigint`. Em `PYG` isso aceita `"75000.50"` e devolve `"75000.00"`, contra a Parte A; e a leitura do `numeric` trunca casas a mais em vez de recusar. O documento 2 já separava "valor (total)" (`numeric(14,2)`) de "custo unitário, custo médio" (`numeric(18,6)`), e o preço de combustível no Brasil tem três casas: nem todo número em dinheiro cabe na escala da moeda.

*Decisão.*

1. **Duas espécies de número em dinheiro, ambas com a `moeda` da linha ao lado.**
   - **Valor**: uma quantia que se paga, se lança ou se soma (total da OS, valor de uma linha de peça, valor de contrato, valor de reposição do ativo). Sempre na **escala da moeda**.
   - **Preço unitário**: dinheiro por unidade de uma quantidade (por hora, litro, quilo, peça; inclui o custo médio de estoque). **Não** é limitado à escala da moeda (o grama de um insumo pode custar meio guarani) e nunca é somado como valor: vira valor só multiplicado por uma quantidade, e aí se arredonda (item 7).
   - O documento 2 classifica cada coluna de dinheiro numa das duas espécies. `usuario.custo_hora` é preço unitário (custo de uma hora).
2. **Escala pela moeda, ISO 4217.** Catálogo fechado no código: `BRL` 2, `ARS` 2, `PYG` 0. A casa decimal é propriedade da **moeda**, não do país (substitui "função do perfil de país" do item 5): fica em `plataforma/dinheiro`, e o perfil de país só diz qual moeda a empresa usa. Moeda nova entra com país novo, por requisito, código e teste.
3. **Banco.**
   - Valor em `numeric(18,2)`; preço unitário em `numeric(18,6)`, em qualquer tabela (substitui a "largura" do item 5 e o `numeric(14,2)` dos totais do documento 2). Dezoito dígitos dão até 10^16 unidades da moeda no valor, com folga para o guarani e para a inflação do peso, e fazem a unidade menor caber em `int64` no Go (10^18 − 1 < 2^63). O `numeric` ocupa o espaço do número guardado, não o da precisão declarada: a largura não custa disco.
   - Toda coluna de valor tem `CHECK (<coluna> = round(<coluna>, app.casas_moeda(moeda)))`, com uma função `IMMUTABLE` que devolve as casas do catálogo do item 2 (nome proposto; o nome e a migração são do `tenaz-banco-dados`). Ela nasce com a primeira coluna de valor; a de preço unitário não tem esse `CHECK`.
   - A coluna `moeda` não tem `DEFAULT`: quem grava informa a moeda da empresa. A igualdade com `empresa.moeda` continua conferida no serviço e por teste (item 5 da Parte A).
   - Nenhum SQL converte dinheiro por `float`, `::bigint` ou divisão por 100. Cálculo de dinheiro em SQL (indicadores, relatórios) usa `numeric` e termina em `round(x, app.casas_moeda(moeda))` quando o resultado é valor.
4. **API (JSON).** Texto decimal (documento 3, seção 2.5) com ponto, sem separador de milhar, sem expoente, sem `+` e sem zero à esquerda; o sinal `-` só onde o campo aceita negativo.
   - **Saída canônica.** Valor com exatamente as casas da moeda (`"1250.40"` em `BRL`, `"75000"` em `PYG`). Preço unitário com no mínimo as casas da moeda e no máximo seis, sem zeros à direita além das casas da moeda (`"45.50"`, `"37.8462"`, `"25000"` e `"0.5"` em `PYG`). Nunca `"-0"`.
   - **Entrada.** Aceita menos casas (`"125.5"` vale `125.50`). Aceita casas além da escala só quando são zeros (`"75000.00"` em `PYG`), para não travar o integrador que formata tudo com duas casas. **Recusa, sem arredondar**, o que não é exatamente representável: `"75000.50"` em `PYG`, a terceira casa num valor em `BRL`, a sétima num preço unitário; resposta `validacao_falhou` com `formato_invalido` no campo. O que a pessoa digitou nunca é arredondado.
   - **Forma no contrato.** Valor: `^-?(0|[1-9][0-9]{0,15})(\.[0-9]{1,2})?$`; preço unitário: `^-?(0|[1-9][0-9]{0,11})(\.[0-9]{1,6})?$`. Campo que não aceita negativo usa o mesmo padrão **sem** o `-?`. A regra da escala por moeda é conferida no serviço, que conhece a moeda da empresa. A `moeda` vem ao lado, só leitura.
   - **A gramática da leitura é a do padrão do contrato** (fechamento da versão 0.9): `plataforma/dinheiro` e `dinheiro.ts` aceitam exatamente o que o padrão aceita, e só depois conferem a moeda. Por isso "zeros além da escala" vale só entre as casas da moeda e o limite da espécie: `"75000.00"` em `PYG` passa; `"1.500"` em `BRL`, `"75000.000"` em `PYG` e `"0.1234560"` como preço unitário são recusados, como o contrato já os recusaria. Assim, a resposta é a mesma quer o erro seja apanhado pelo padrão, quer pela biblioteca.
5. **Go: pacote `plataforma/dinheiro`, só biblioteca padrão.**
   - Tipo `Moeda` com as casas; `plataforma/pais` passa a usar `dinheiro.Moeda` (`pais` importa `dinheiro`, nunca o contrário).
   - `Valor`: `int64` na **unidade menor da moeda** (o centavo em `BRL` e `ARS`, o próprio guarani em `PYG`) mais a `Moeda`. Não é a convenção de centavos descartada no item 5: a escala vem da moeda, e o banco e a API continuam decimais.
   - `PrecoUnitario`: `int64` em milionésimos da unidade da moeda mais a `Moeda`.
   - Operações: somar e subtrair só com a mesma moeda, com verificação de faixa (o resultado tem de caber na coluna); multiplicar preço unitário por quantidade (quantidade como `*big.Rat` exato, lida do `numeric` ou de minutos ÷ 60) com o cálculo em `math/big` e **um** arredondamento no fim; ratear um valor em partes distribuindo a sobra de unidades menores pelo maior resto, para a soma das partes ser igual ao valor.
   - Conversão exata com o banco pelo `pgtype.Numeric` (`Int` e `Exp`), sem `float64`; número com mais casas do que a espécie e a moeda permitem, ou fora da faixa, é erro, nunca truncamento.
   - Texto: ler valor e preço unitário de texto pela moeda (item 4) e escrever o texto canônico; os erros de leitura saem pelo catálogo (`ComCampoChave`, seção 8.2).
   - Nenhum `float32` ou `float64` em dinheiro, em nenhuma camada.
6. **TypeScript: `pacotes/dominio/src/dinheiro.ts`, sem biblioteca.** As mesmas espécies, leitura, texto canônico, multiplicação, arredondamento e rateio, sobre `bigint` nativo (o alvo `ES2022` de `pacotes/dominio` já o tem, e o Hermes do Expo SDK 58 também; se o empacotador do app recusar o literal `10n`, usa-se `BigInt(10)`). No contrato e nos tipos gerados, dinheiro continua `string`; `Number()`, `parseFloat` e aritmética de `number` sobre dinheiro são proibidos (seção 8.4). A formatação para pessoas (vírgula decimal, símbolo) é a da Parte A, item 7, aplicada sobre o texto canônico.
7. **Arredondamento.** Meio para longe do zero (`2.345` → `2.35`, `-2.345` → `-2.35`), o mesmo do `round(numeric, n)` do PostgreSQL, para o cálculo em SQL, em Go e em TypeScript dar o mesmo resultado no empate. Um só arredondamento, no fim do cálculo de cada valor (preço × quantidade, minutos × `custo_hora` ÷ 60); nunca nas parcelas intermediárias. Total é a soma dos valores das linhas já arredondados e gravados, nunca recalculado a partir dos preços. Não é regra fiscal: a Tenaz não emite documento fiscal (Parte C, item 7).
8. **Casos compartilhados.** `backend/testes/casos/dinheiro.json`, rodado pelas suítes Go e TypeScript: leitura válida e inválida por moeda e espécie, texto canônico, empates positivos e negativos, `PYG`, limites de faixa, `-0` e rateio com sobra. (Versão 1 do arquivo escrita pelo `tenaz-qa-automacao` em 13/11/2026, 239 casos conferidos por oráculo exato; os pontos que ele levantou estão fechados no item 10. Versão 2, com 266 casos e o grupo de construção da moeda, ajustada ao item 10.) O JSON não é editado à mão: é gerado por `backend/testes/casos/gerar_dinheiro.py`, oráculo de aritmética exata (`fractions` da biblioteca padrão do Python), e a tarefa `conferir-casos` (chamada por `conferir` e no job de documentação do CI) falha se o arquivo divergir do gerador.
9. **O que muda agora em `usuario.custo_hora`.** Passa a preço unitário: `numeric(18,6)`, `CHECK (custo_hora >= 0)` e `usuario.moeda` sem `DEFAULT`, numa migração nova (a `00010`, salvo nome do `tenaz-banco-dados`), sem editar a `00004` nem a `00009`; os valores gravados não mudam (a escala só aumenta). No Go, saem `api.Centavos`, `api.Dinheiro`, `adm.centavos`, o campo `CustoHoraCentavos` e o parâmetro `custo_hora_centavos` do sqlc, e entra `dinheiro.PrecoUnitario`; o erro de faixa passa a apontar o campo `custo_hora` (hoje aponta `custo_hora_centavos`, que não existe no contrato). Na API, `custo_hora` aceita até seis casas, sem sinal (padrão `^(0|[1-9][0-9]{0,11})(\.[0-9]{1,6})?$`), e sai no texto canônico do item 4. Na web, o formulário de usuário valida pelo `dinheiro.ts`, com a moeda da empresa, no lugar da expressão de duas casas.
10. **Fechamentos pedidos pelo `tenaz-qa-automacao` (versão 0.9, depois da versão 1 dos casos).**
    - **Rateio.** Pesos em texto decimal **não negativo** (a mesma gramática da quantidade, sem sinal), lista não vazia e soma positiva. Cada parte recebe o piso do seu quinhão sobre o **módulo** do valor; a sobra de unidades menores vai uma a uma pelo maior resto; **empate de restos vai para o menor índice**, na ordem de entrada; o sinal do valor volta a todas as partes (`-100.00` em três → `-33.34`, `-33.33`, `-33.33`), de modo que ratear `-v` é sempre o oposto de ratear `v`. A soma das partes é sempre o valor. As três propostas do QA ficam confirmadas.
    - **Quantidade na multiplicação.** A operação recebe a quantidade como **fração exata**: `*big.Rat` no Go, par `{ numerador, denominador }` de `bigint` no TypeScript. O código de produção monta a fração (minutos e `60`) ou a lê de um texto decimal do contrato; a notação `a/b` é só do arquivo de casos, e o leitor dela fica no código de teste das duas suítes. A API nunca aceita `a/b`. A quantidade pode ser negativa na biblioteca (estorno, ajuste de inventário: `-0.005` × `1.00` → `-0.01`); se um campo aceita quantidade negativa é regra do campo, nos documentos 2 e 3.
    - **`-0` em campo sem negativo** é recusado com `formato_invalido`: o sinal é recusado onde o campo não aceita negativo, qualquer que seja o número, como faz o padrão sem `-?`. Onde o campo aceita negativo, `-0` vale zero e sai sem sinal.
    - **Erro de leitura é sempre `formato_invalido`**, inclusive para sinal não aceito e número com dígitos demais: o limite da faixa é o número de dígitos do padrão, então todo valor fora da faixa já é erro de forma, e é o que o contrato devolve (`expected string to match pattern` vira `formato_invalido` em `plataforma/api/rota.go`). A chave do texto pode dizer à pessoa qual foi o problema; o código é um só. `fora_da_faixa` fica para o **resultado de uma operação** (soma, subtração ou multiplicação que sai da coluna) e para limite de campo mais estreito que o da espécie, se algum campo o tiver. O `custo_hora` não tem limite próprio: o `fora_da_faixa` de hoje em `usuarios.go` sai junto com `validarCusto`.
    - **Erros internos, que não chegam à API com nome próprio.** `moeda_diferente` (somar ou comparar moedas diferentes), `moeda_desconhecida` e `rateio_invalido` (pesos vazios, soma zero ou peso negativo) são erros da biblioteca, valores sentinela no Go e no TypeScript, com esses nomes nos casos compartilhados. No MVP, todo valor de uma empresa tem a moeda dela, e os pesos de rateio vêm do código: qualquer um dos três é defeito e sai como `erro_interno` (500), registrado no log. Se um dia uma pessoa digitar pesos de rateio, o serviço valida o campo antes (`validacao_falhou` com `valor_invalido`). Nenhum código novo vai para o documento 3. `fora_da_faixa` vindo de uma operação sobre entrada da pessoa sai como `validacao_falhou` com `fora_da_faixa` no campo que a originou (vocabulário que já existe).
    - **Moeda na leitura.** A leitura recebe o tipo `Moeda`, nunca texto: a moeda vem da empresa, gravada e conferida pelo banco. Só a construção de `Moeda` a partir de um código (`BRL`, `ARS`, `PYG`) pode falhar, com `moeda_desconhecida`, e isso é defeito, porque o cliente não envia moeda.
    - **Preço unitário não se soma** por restrição de **tipo**, sem erro em tempo de execução nem caso compartilhado: no Go, `PrecoUnitario` não tem `Somar` nem `Subtrair`; no TypeScript, valor e preço unitário são tipos distintos, com o discriminante `especie`, e `somar` aceita só `Valor`.

*Alternativas consideradas.*

| Alternativa | A favor | Contra | Resultado |
|---|---|---|---|
| **`shopspring/decimal`** (MIT), com o adaptador `pgx-shopspring-decimal` | Pronta e conhecida | Duas dependências novas; precisão arbitrária sem necessidade; precisão de divisão em variável global mutável (`DivisionPrecision`); ainda precisaria de um tipo com moeda e escala por cima | Descartada: a biblioteca padrão basta (ADR-012, preferência por não acrescentar dependência) |
| **`cockroachdb/apd`** (Apache-2.0) | Decimal exato com contexto de arredondamento explícito | Mais pesada; o mesmo envoltório de moeda e escala | Descartada; gatilho (i) |
| **`big.Rat` ou `big.Int` com expoente como o valor guardado** | Exato, biblioteca padrão | Alocação em toda operação, sem faixa, igualdade dependente de normalização | Descartada como representação; usada só no cálculo intermediário |
| **`pgtype.Numeric` no domínio** | Nenhuma conversão | Tipo do driver no domínio, sem aritmética | Descartada |
| **Uma espécie só, tudo na escala da moeda** | Uma regra a menos | Preço de combustível com três casas e custo médio de estoque perderiam precisão; em `PYG` o preço de um grama viraria zero | Descartada |
| **`custo_hora` como valor** (`numeric(18,2)`, escala da moeda) | Nenhuma mudança de escala agora | O custo-hora sai muitas vezes de salário dividido por horas, e o erro de arredondamento se multiplica pelas horas apontadas; abriria exceção na regra das espécies | Descartada |
| **Escala da moeda no tipo da coluna** (`numeric(18,0)` para `PYG`) | O tipo recusa casas no guarani | A escala é da coluna, e a moeda é da linha: exigiria uma coluna por moeda | Descartada: `CHECK` por linha |
| **Número JSON** na API | Sem aspas | `float64` no JavaScript; contraria o documento 3 | Descartada |
| **Inteiro na unidade menor na API** | Sem leitura de decimal no cliente | Todo integrador precisa saber a escala de cada moeda; contraria o documento 3 e a Parte A | Descartada |
| **Arredondar a entrada** à escala da moeda | Nunca recusa | Esconde erro de digitação: `"75000.50"` em `PYG` viraria `75001` sem aviso | Descartada |
| **Arredondamento bancário** (meio para o par) | Menos viés em somas grandes | Difere do `round(numeric)` do PostgreSQL: o cálculo em SQL e em Go divergiria no empate | Descartada |
| **Chave estrangeira composta** `(empresa_id, moeda)` para `empresa (id, moeda)` | O banco garante a moeda da empresa | Índice único a mais em `empresa` e uma chave por tabela, a desfazer no gatilho (e); serviço e teste já garantem | Descartada no MVP |

*Consequências.*

- Positivas: `PYG` correto de ponta a ponta; nenhuma dependência nova; o mesmo resultado no SQL, no Go e no TypeScript, provado pelos mesmos casos; aritmética em `int64`, sem alocação no caminho comum; a recusa de valor não representável pega erro de digitação em vez de escondê-lo; o isolamento não muda (nenhuma tabela nova).
- Negativas: código de dinheiro mantido em duas linguagens (contido pelos casos compartilhados); cada coluna de dinheiro precisa ser classificada como valor ou preço unitário; a faixa precisa ser conferida em toda soma; a migração `00010` e a troca do tipo no ADM e na web entram ainda na E01.
- Gatilho de revisão: (i) cálculo financeiro com divisões encadeadas (juros, depreciação composta, rateio por índice) que o `int64` com `math/big` torne difícil de manter, ou necessidade de mais de 18 dígitos → reavaliar uma biblioteca decimal pela regra do ADR-012, em emenda desta Parte F. A multimoeda continua no gatilho (e).

**Complemento de 13/11/2026 (etapa E01, versão 0.9 deste documento): Parte G — Formatação para pessoas.** Decisão técnica do `tenaz-arquiteto`, a partir de divergência apontada pelo `tenaz-produto`.

*Contexto.* A Parte A, item 7, manda formatar número, data e moeda por funções próprias, com casos compartilhados, e nunca pelo `Intl` com `es-419`. A web foi feita com `Intl` (`web/src/lib/formatacao.ts`), pela localidade `pt-BR`, `es-AR` ou `es-PY` escolhida pela moeda, e o documento 7, seção 3.11, descreve isso. Isso evita o ponto decimal do `es-419` quando a moeda é `ARS` ou `PYG`, mas não o resto: (1) numa empresa brasileira, o usuário em `es-419` cai no `Intl` `es-419` e vê o separador errado; (2) `formatarMoeda` recebe `number`, contra a Parte F; (3) `formatarPrecoUnitario` passa o texto ao `Intl.NumberFormat`, que só o lê exato nos motores com NumberFormat v3, e o Hermes do app não o tem; (4) a saída do `Intl` depende da versão do ICU e do CLDR de cada navegador (o símbolo do guarani, o tipo de espaço antes do símbolo), então a mesma tela formata diferente em máquinas diferentes, e um PDF gerado no Go nunca casaria com ela; (5) não há caso compartilhado. Hoje quase nada formata: na web, `formatarPrecoUnitario` numa tela e datas na auditoria e nas unidades; no app, a data do build no "Sobre"; no Go, nada.

*Decisão: manter o ADR (opção (b)) e mover a formatação para `pacotes/dominio`.*

1. **Formato pelo país da empresa, palavras pelo idioma do usuário.** Separadores, padrão de data e símbolo da moeda vêm do perfil de país da empresa (os três países usam vírgula decimal, ponto de milhar, `dd/mm/aaaa` e 24 h); o idioma só escolhe as palavras, que vêm do catálogo (`há {{n}} min`, "ontem às"). Enquanto a sessão da web não trouxer o país, ele é tirado da moeda (`paisDaMoeda`): no MVP a relação é de um para um (Q3).
2. **Funções próprias em `pacotes/dominio/src/formatos.ts`, sem `Intl` para número nem moeda.** Entram textos canônicos (Parte F, item 4, e os decimais do documento 3, seção 2.5) ou `bigint`, nunca `number` para dinheiro nem para decimal de medidor; contagens inteiras podem ser `number`. A formatação trabalha sobre os dígitos do texto: exata, sem arredondar o que veio do servidor, e igual em qualquer motor. Espaço inseparável é sempre U+00A0.
3. **Datas.** A data civil (`AAAA-MM-DD`) é montada sem fuso. O instante é mostrado num **fuso explícito** (da unidade ou, sem unidade, da empresa; nunca o fuso implícito do navegador): a conversão de fuso usa o `Intl.DateTimeFormat` só como motor da base tz, com localidade fixa e as partes numéricas de `formatToParts`, e o texto é montado pela função própria. No app, o "Sobre" usa o fuso do aparelho, que é explícito ali.
4. **`Intl` permitido só para:** conversão de fuso (item 3), ordenação (`Intl.Collator`) e conferência de plural nos testes (`Intl.PluralRules`, seção 10.6). Proibidos em código de interface: `Intl.NumberFormat`, `Intl.RelativeTimeFormat`, `Intl.ListFormat`, o `Intl.DateTimeFormat` como gerador de texto e `toLocaleString`, `toLocaleDateString` e `toLocaleTimeString`. **Trava:** regra do ESLint (`no-restricted-syntax` ou `no-restricted-properties`) na web, no app e nos pacotes, com exceção só para o arquivo de `formatos.ts` que faz a conversão de fuso; responsável `tenaz-frontend-web`, como o `tenaz-backend-go` fez com o texto livre.
5. **Casos compartilhados:** `backend/testes/casos/formatos.json`, com inteiro, decimal, percentual, medidor com unidade, valor e preço unitário nas três moedas (inclusive negativo e acima de 2^53), símbolo por país, data civil e data e hora em fuso de cada país. Os instantes ficam longe de mudança de horário, para não depender da versão da base tz. Rodado por `pacotes/dominio` (a web e o app usam as mesmas funções) e pelo Go quando houver formatador no Go. Versão 1 escrita pelo `tenaz-qa-automacao` em 13/11/2026: 115 casos, gerados por `backend/testes/casos/gerar_formatos.py` (oráculo só com a biblioteca padrão do Python, que recusa instante fora das janelas sem mudança de horário) e conferidos pela tarefa `conferir-casos`, no `conferir` e no CI, como o `dinheiro.json`; os pontos que ele levantou foram fechados no documento 7, seção 3.11.
6. **O desenho é do documento 7:** símbolo de cada moeda, posição do sinal negativo, casas por grandeza. Os casos registram o que o documento 7 fixar. Lá, a regra da Parte A, item 7, continua: símbolo do país quando a moeda é a do país (`R$`, `$`, o do guarani).
7. **O que entra agora na E01:** `formatos.ts` com número, decimal, percentual, valor, preço unitário, símbolo, data e data e hora; `formatos.json`; a web passa a usar só essas funções (saem `formatarMoeda(number)`, `formatarNumero` com opções do `Intl`, `localidadeDosNumeros` e `simboloDaMoeda` pelo `Intl`; `listar`, `mascararData` e `dataParaISO`, que não usam `Intl`, podem ficar ou subir para `pacotes/dominio`); a regra do ESLint; e a data do build no "Sobre" do app, que é uma função só. **Pode esperar:** o formatador do Go, até o primeiro e-mail ou PDF com número, dinheiro ou data, quando ele nasce passando nos mesmos casos; hora relativa, duração e prazo futuro, até a primeira tela que os use, com as palavras do catálogo; a conversão de fuso no Hermes, a conferir pelo `tenaz-mobile-offline` na primeira tela do app com data de negócio (se o Hermes não der as partes num fuso nomeado, entra uma alternativa nesta Parte G).

*Alternativas consideradas.*

| Alternativa | A favor | Contra | Resultado |
|---|---|---|---|
| **(a) Emendar o ADR e aceitar `Intl` com a localidade do país na web** (condições: localidade pelo país da empresa, nunca `es-419`; dinheiro só por texto; testes por navegador) | Nenhum código de formatação a manter; a web já está feita | Saída dependente do ICU e do CLDR de cada navegador; o app (Hermes, sem NumberFormat v3) e o Go precisariam de outro caminho, e a tela, o app e o PDF divergiriam; dinheiro acima de 2^53 ou com seis casas perderia precisão onde o motor não lê texto; nenhum caso compartilhado possível, porque o resultado varia por motor | Descartada |
| **Biblioteca de formatação** (por exemplo, `date-fns` com localidades, ou uma de número) | Pronta | Dependência nova para três países com as mesmas convenções; não resolve o Go; os formatos do documento 7 são poucos e fixos | Descartada |
| **Formatar no servidor e mandar o texto pronto** | Uma implementação só | Contraria a Parte B, item 4 (recursos JSON sem texto traduzido, `ETag` independente do idioma); o app offline precisa formatar sem servidor | Descartada |

*Consequências.* Positivas: a mesma saída na web, no app e, depois, no PDF, provada pelos mesmos casos; dinheiro e decimais formatados sem passar por `number`; nenhuma dependência nova; o defeito do usuário `es-419` em empresa brasileira some. Negativas: código de formatação a manter (pequeno: três países com as mesmas convenções); a web refaz `formatacao.ts` ainda na E01; o documento 7, seção 3.11, muda.

---

## 16. Alinhamento com o sumário e pendências

Situação deste documento em relação ao sumário executivo 0.8.

**Versão 0.8 (etapa 014, sumário 0.8): [ADR-013](#adr-013--país-e-idioma-perfil-de-país-no-núcleo-textos-por-catálogo-e-conteúdo-legal-por-país-e-jurisdição), aceito pelo dono do produto em 13/11/2026, a aplicar fora deste documento.** Nada abaixo foi aplicado ainda: os nomes entram primeiro no sumário. Cada dono altera o que é seu, no mesmo passo da mudança de código (CLAUDE.md, seção 3.2).

| Dono | O que muda |
|---|---|
| `tenaz-produto` (sumário) | Seção 1.4 ("Idioma: português do Brasil" e "Hospedagem") e seção 3.6: países, idiomas, moeda e os nomes canônicos da tabela do ADR-013 (aceito); regra "CNPJ da empresa" passa a "identificação da empresa"; "dado sensível" passa a cobrir a identificação pessoal de qualquer país; P14 com o aditivo de transferência por país; respostas Q1 a Q3 abaixo; prefixo do número da OS por país (`OS-`, `OT-`) |
| `tenaz-produto` (documento 1) | RNF-029 ("Idioma e formatos", hoje "moeda R$") reescrito para os dois idiomas e os três países; RNF novo de países e de idioma por usuário; RF-ADM-001, RF-ADM-002, RF-ADM-003 e o de usuários (país, identificação, idioma, fusos, feriados por país); RF-ADM-005 (permissões com norma no nome); dados principais de CON (jurisdição, prazo preenchido pela empresa, registro de conformidade); erros "em português" do requisito da API pública passam a "no idioma negociado" |
| `tenaz-banco-dados` (documento 2) | Seção 1.4 (tipos: identificação, moeda, largura dos totais); seção 5.1 e a migração `00009` (Parte A, item 4); prefixo do número da OS por país; `perfil.nome` dos perfis padrão; CON: `registro_conformidade`, jurisdição, prazo por jurisdição, país do tipo de documento, `CHECK` de prazos brasileiros convertidos em dado, unicidade de `bloqueio_ponto` para uso múltiplo e retirada só por quem aplicou; ATV e OS: normas aplicáveis e exigidas; EST: documento fiscal por tipo; `conexao_empresa`: identificação e regra de "mesma entidade" por país (com a E01-13) |
| `tenaz-api-contrato` (documento 3) | Seção 2.4 (fuso por país), 2.5 (moeda em todo recurso com dinheiro, casas pela moeda), 2.15 (negociação, `Content-Language`, `Vary`, recursos sem texto traduzido), 2.16 (cabeçalhos), 3.1 (`title`, `detail` e `mensagem` no idioma negociado; exemplo em espanhol); campos de `/v1/empresas`, unidades, pessoas e usuários; `identificacao_mascarada` |
| `tenaz-dominio-manutencao` (documento 4) | Regras de CON: a OS exige os registros das normas aplicáveis, sem citar a NR-12 no fluxo do núcleo; bloqueio de energia com guarda e retirada por quem aplicou, uso múltiplo; prazo resolvido pela jurisdição emissora; regras de conflito do app para registro de conformidade |
| `tenaz-design-sistema` e `tenaz-idiomas` (documento 7) | Formatos por país (número, data, moeda com símbolo ou código); seletor de idioma no aceite do convite e no perfil; aviso de "prazo definido pela empresa"; glossário bilíngue e termos proibidos em espanhol; tom de voz em `es-419` |
| `tenaz-plataforma-sre`, `tenaz-conformidade-br` e `tenaz-conformidade-mercosul` (documento 9) | Transferência internacional por país (AAIP: cláusulas modelo; Paraguai: art. 19 da Ley 7593/2025, vigência por volta de 27/11/2027), papel de encarregada, aditivo de transferência por país e termo da conexão entre países (P14); identificação pessoal dos três países como dado de acesso restrito; tzdata atualizada como item de operação; CNPJ alfanumérico (confirmar a IN RFB 2.229/2024) |
| `tenaz-backend-go` | `plataforma/pais`, `plataforma/textos`, tipo de dinheiro; `erros` com título e detalhe por chave; e-mails por idioma; negociação de idioma no middleware; migração `00009` com o `tenaz-banco-dados`; `time/tzdata` no binário |
| `tenaz-frontend-web` e `tenaz-mobile-offline` | `es-419.ts` com `satisfies Textos`; detecção do idioma; `Accept-Language`; formatação própria em `pacotes/dominio`; `expo-localization` no app |
| `tenaz-plataforma-sre` | Tarefa `textos` chamada por `conferir`; atributo `pais` nos sinais |
| `tenaz-qa-automacao` | Casos compartilhados de identificação, arredondamento por moeda e formatação; teste de erro renderizado nos dois idiomas; teste de que o hash da trilha não muda com o idioma |

**Perguntas ao dono do produto (ADR-013): respondidas em 13/11/2026.** Não eram P1 a P8; Q3 se liga à P15.

| # | Pergunta | Opções | Resposta |
|---|---|---|---|
| Q1 | Uma empresa pode ter unidade em outro dos três países no MVP (por exemplo, indústria brasileira com fábrica no Paraguai)? | (a) Não no MVP: `unidade.pais` existe e só aceita o país da empresa. (b) Sim: identificação, fuso, feriados e conteúdo legal pela unidade, com a moeda da empresa até a Q3 | **(a), mais restrita:** um país por empresa; a unidade não tem país próprio, e `unidade.pais` sai do MVP (gatilho (h)) |
| Q2 | Idioma de último recurso para quem não fala nenhum dos dois (navegador em inglês, por exemplo) | (a) `pt-BR`, idioma de origem. (b) Pelo país da empresa, quando houver empresa; senão `pt-BR` | **(a)**, a regra proposta, terminando em `pt-BR`: decidida por padrão, salvo objeção |
| Q3 | Multimoeda no MVP (compra em dólar, comum na Argentina; unidade em país de outra moeda) | (a) Não: uma moeda por empresa, gatilho (e). (b) Sim, com cotação | **(a):** uma moeda por empresa |

**Versão 0.9 (etapa E01): estado do ADR-013 no backend e aplicação da Parte F.** Conferido no código em 13/11/2026.

| Dono | Situação |
|---|---|
| `tenaz-backend-go` | **Feito:** `plataforma/textos` (catálogos embutidos, `T`, `Plural`, negociação por `Preferencias.Resolver()`), `plataforma/pais` (perfil, fusos do país, identificação, CNPJ alfanumérico, casos `identificacao.json`), erros por chave com `Content-Language` e `Vary`, e-mails pelo catálogo (`email.*`), redação de `identificacao` nos logs, `time/tzdata` no binário; todo o texto de erro no catálogo (118 pontos) e a trava contra texto livre na tarefa `textos` (seção 8.2). **A fazer na E01 (Parte F):** `plataforma/dinheiro`, `dinheiro.Moeda` em `plataforma/pais`, troca de `custo_hora` (Parte F, item 9). **Falta no perfil de país:** formatação e feriados nacionais (Parte A, itens 7 e 8) |
| `tenaz-banco-dados` | **A fazer na E01 (Parte F, item 9):** migração `00010` em `usuario` (`custo_hora numeric(18,6)`, `CHECK (custo_hora >= 0)`, `moeda` sem `DEFAULT`). Documento 2, seção 1.4: as duas espécies de dinheiro no lugar de "Valor em dinheiro (total) `numeric(14,2)`", e a função `app.casas_moeda` descrita para nascer com a primeira coluna de valor |
| `tenaz-api-contrato` | Documento 3, seção 2.5: "moeda implícita BRL" e "duas casas" dão lugar às regras de entrada e saída da Parte F, item 4; `custo_hora` com até seis casas; `erros[].valores_aceitos` (seção 8.2) no formato de erro. Divergência a resolver: o documento 3 limita custos unitários de estoque a quatro casas, e o documento 2 guarda seis; pela Parte F, a API não pode perder casas do que está gravado, então o limite é seis |
| `tenaz-frontend-web` e `tenaz-mobile-offline` | `pacotes/dominio/src/dinheiro.ts` (Parte F, item 6) e o formulário de usuário validado por ele; a tabela de identificação do perfil de país em `pacotes/dominio`, lendo `identificacao.json` (hoje `documentos.ts` valida só CPF e CNPJ) |
| `tenaz-qa-automacao` | `backend/testes/casos/dinheiro.json` (Parte F, item 8), rodado pelas duas suítes; `backend/testes/casos/formatos.json` (Parte G, item 5), com os formatos do documento 7 |
| `tenaz-frontend-web` (Parte G) | **A fazer na E01:** `pacotes/dominio/src/formatos.ts`; `web/src/lib/formatacao.ts` sem `Intl` para número, moeda e texto de data; fuso explícito em `formatarDataHora`; regra do ESLint contra `Intl` e `toLocale*String` na web, no app e nos pacotes |
| `tenaz-mobile-offline` (Parte G) | **A fazer na E01:** a data do build no "Sobre" (`app/src/acesso/versao.ts`) por `formatos.ts`, com o fuso do aparelho. **Depois:** conferir a conversão de fuso no Hermes na primeira tela com data de negócio |
| `tenaz-design-sistema` (documento 7, Parte G) | Seção 3.11: formatação pelas funções de `pacotes/dominio` e pelo país da empresa, não pelo `Intl`; símbolo de cada moeda e posição do sinal fixados para os casos; o exemplo da moeda alinhado à Parte A, item 7 (símbolo do país, não o código ISO) |
| `tenaz-backend-go` (Parte G) | **Depois:** formatador no Go com o primeiro e-mail ou PDF com número, dinheiro ou data, passando em `formatos.json` |

**Resolvido na versão 0.7 (etapa E01, sumário 0.7) e aplicado aqui**

| Tema | Como ficou |
|---|---|
| Só ferramentas e recursos 100% gratuitos ou de código aberto (requisito do dono do produto) | [ADR-012](#adr-012--só-software-livre-ou-gratuito-e-seaweedfs-no-lugar-do-minio) e a sua emenda: a regra vale só para o que o sistema usa diretamente (o desenvolvimento fica fora); produção em nuvem gerenciada como exceção explícita; verificação de licenças no CI só sobre o que vai para o produto (proposta). Seções 2.3, 2.4, 7 e 13.1; notas no ADR-004 e no ADR-006, que volta a valer. |
| Substituto do MinIO | Desenvolvimento: SeaweedFS 4.53, numa máquina, endurecido com `security.toml`, aplicado no compose em 13/11/2026 (19 de 19 provas de `armazenamento-provar`); RustFS qualificado como alternativa. Produção: S3 com Object Lock em conformidade. Backend preso à API S3. Ensaio prático no ADR-012. |
| Dado de saúde (sumário 3.6, "Dados de saúde") | A Tenaz não guarda dado de saúde: saem as citações das seções 6.8, 9.3 (o teste de corrida usa `os.encerrar`), 12.3 e 14.3. Nenhuma referência a `pessoa.saude.ler` resta neste documento. |
| Tabela da seção 9.3 contra a suíte | Linhas "Consultas sqlc", "Mapa de privilégios" e "Funções de credencial" conferidas com `catalogo_test.go` e `consultas_test.go`; linha nova "NaEmpresa confere o escopo no banco", conferida com `pool_test.go`; linhas novas "Extensões" (`catalogo_test.go`, `extensoes.txt`) e "Sem dado de saúde" (`saude_test.go`). |
| Expo SDK 58 e sessão do app (`tenaz-mobile-offline`) | SDK 58 ratificado (seção 2.3, auditoria do ADR-012). Uma credencial viva por aparelho, troca de empresa só com conexão, banco por conta preservado ao sair, cofre com itens contados e permissões no banco local (seções 10.6 e 10.8, conferidas com `app/src/acesso/sessoes.ts`, `app/src/sync/deposito-expo.ts` e `app/app.config.ts`); sem identificador de atualização até haver servidor de atualização (seção 9.5, H5). |
| Testes de integração e access token | Seção 8.2: PostgreSQL descartável por pacote, subida serializada por arquivo de trava e repetição com recuo da detecção do Docker (conferido com `dbteste.go`). Seção 10.6: verificação estrita do access token (conferida com `jwt.go`). |

**Perguntas da auditoria do ADR-012 (versão 0.7).** H1 a H4 e H10 foram respondidas pelo dono do produto em 13/11/2026 (emenda do ADR-012); H5 a H9 ficam abertas até a etapa em que o serviço entra.

| # | Pergunta | Situação | Etapa |
|---|---|---|---|
| H1 | Onde a Tenaz roda em homologação e produção | **Resolvida:** nuvem gerenciada (AWS, São Paulo), exceção explícita à regra; o ADR-006 volta a valer | — |
| H2 | Software proprietário gratuito atende? | **Resolvida pelo alcance:** ferramenta de desenvolvimento fica fora da regra; usado pelo sistema, o gratuito sem condição atende ao item (b) | — |
| H3 | Docker Desktop | **Resolvida pelo alcance:** ferramenta de desenvolvimento, fora da regra; continua | — |
| H4 | GitHub e GitHub Actions | **Resolvida pelo alcance:** fora da regra; continuam | — |
| H5 | Distribuição do app: lojas (Apple Developer anual e Google Play taxa única, sem substituto para iOS), build de iOS (macOS e Xcode), EAS gratuito com cota ou build local com servidor próprio do `expo-updates`, conta Apple para o APNs | Aberta | E12 (lançamento); o teste fechado nas lojas do risco 15 do roadmap pode antecipar |
| H6 | E-mail transacional em homologação e produção: provedor pago ou Postfix próprio, com risco de entregabilidade. Mailpit no desenvolvimento | Aberta | Até a primeira implantação em homologação; no mais tardar E12 |
| H7 | SMS (P11): não há SMS gratuito | Aberta, com a P11 | E12, ou antes se o login por telefone for antecipado |
| H8 | IA: provedor pago (ADR-007), modelo aberto com GPU paga, ou adiar | Aberta | Fase 2 |
| H9 | Pagamento da assinatura (RF-ADM-009): todo gateway cobra tarifa | Aberta | E10 (API, SSO e cobrança) |
| H10 | Imutabilidade do selo | **Resolvida:** produção no S3 com Object Lock em conformidade, que resiste à conta root; desenvolvimento numa máquina, com SeaweedFS endurecido e o limite aceito | — |

**Aplicação do ADR-012 fora deste documento (versão 0.7).** Cada dono altera o que é seu, no mesmo passo da mudança de código:

| Dono | O que muda |
|---|---|
| `tenaz-plataforma-sre` | **Feito em 13/11/2026:** SeaweedFS 4.53 por digest no lugar do MinIO, com `security.toml` e `s3.json` versionados, só a porta S3 publicada e `healthcheck`; `postgres`, `valkey` e `mailpit` por versão e digest; `tarefas.sh infra` espera os quatro saudáveis e cria o bucket `tenaz-auditoria`; tarefa `armazenamento-provar` com 19 verificações; tarefas `api` e `selar-auditoria` com as variáveis `TENAZ_S3_*` de desenvolvimento (13/11/2026). **Falta:** a tarefa `licencas` (decisão 2), chamada por `conferir`, só sobre as dependências de produção; e o job de CI da linha abaixo |
| `tenaz-plataforma-sre` com `tenaz-qa-automacao` | **Pendência:** o CI ainda não sobe o SeaweedFS. Job que rode `tarefas.sh infra` e `armazenamento-provar` em todo PR, para que a qualificação do endurecimento (decisões 4 e 5) não dependa da máquina de quem desenvolve |
| `tenaz-backend-go` | **Feito para o selo em E01 (versão 0.8):** `internal/plataforma/armazenamento`, interface `Armazenamento` sobre `aws-sdk-go-v2`, sem `minio-go`, com endpoint, região e estilo de caminho por configuração (seções 2.4 e 4.1). Falta para E02 (anexos): URL pré-assinada na mesma interface. Antes: quando o armazenamento entrar (E01, selo; E02, anexos): cliente `aws-sdk-go-v2` atrás de uma interface da plataforma, sem `minio-go`, com endpoint, região e estilo de caminho por configuração. A imagem do PostgreSQL dos testes com testcontainers foi alinhada à do compose em 13/11/2026 (`dbteste.ImagemPostgres`, por digest) |
| `tenaz-qa-automacao` | Teste de integração do armazenamento com as 15 verificações do ensaio e o `DELETE` no filer sem JWT esperando 401 |
| `tenaz-banco-dados` | Documento 2: o bucket do selo e do arquivo da auditoria (`selo_auditoria.chave_objeto`) é, em produção, S3 com Object Lock em modo de conformidade, e no desenvolvimento, SeaweedFS. KMS, Secrets Manager e RDS ficam como estão (exceção do ADR-006) |
| `tenaz-plataforma-sre` (documento 9, parte de operação) | Nada a desfazer: a nuvem gerenciada continua (exceção do ADR-012). Registrar que o selo em produção usa S3 Object Lock em conformidade e que o SeaweedFS do desenvolvimento não protege contra o operador da máquina (sem dado real) |
| `tenaz-produto` (documento 6) | S12 ("Publicação nas lojas") e risco 15 dependem da H5; E10 (cobrança), da H9; a primeira homologação e E12, da H6 e da H7. M0 fica como está |

**Resolvido na versão 0.6 (etapa E01, sumário 0.6) e aplicado aqui**

| Tema | Como ficou |
|---|---|
| Nomes da autenticação (sumário 0.6, seção 3.2) | `identidade`, `acesso_empresa`, `sessao`, `token_renovacao`, `codigo_autorizacao`, `desafio_mfa` e `token_email` citados com os nomes aprovados (seção 6.9 e ADR-011). O calendário da unidade (`turno_unidade`, `feriado_unidade`), `usuario_permissao` e `usuario.situacao` não afetam a arquitetura. |
| Achado ISO-E01-02: `tenaz_app` lia e forjava credenciais de qualquer empresa | [ADR-011](#adr-011--papel-de-banco-próprio-para-as-credenciais-de-autenticação) e seção 6.9: papel `tenaz_autenticacao` com pool exclusivo do pacote de autenticação; `tenaz_app` e `tenaz_job` sem privilégio nas tabelas de credencial; auditoria de login pelo próprio papel, contida por política restritiva; revogação pelo trigger na inativação; purga por função estreita. Seções 4.1, 5.1, 5.3 (regra 12), 6.1, 6.3, 9.3 e 13.1. |
| Achado ISO-E01-08: `tenaz_manutencao`, membro de `tenaz_dono` com herança, lia todas as credenciais, porque o dono das sete tabelas sem RLS era `tenaz_dono` | Emenda do [ADR-011](#adr-011--papel-de-banco-próprio-para-as-credenciais-de-autenticação) (decisões 1 e 9): dono `tenaz_credencial`, `NOLOGIN`, sem membros de runtime, também das duas funções `SECURITY DEFINER`; privilégios cruzados só para as chaves estrangeiras e o trigger. Seções 6.1, 6.4, 6.9 e 9.3; documento 2, seções 2.2, 2.6, 2.7, 5.1, 14.2 e 17.3. |
| Pendência E01-3 do documento 2: `app.escopo_id()` | Não é criada. A seção 6.6 usa `app.fornecedor_atual()`, `app.locatario_atual()` e `app.link_hash_atual()`, e a política do link confere, na subconsulta, revogação, expiração, consumo e exclusão do link. As duas condições da revisão de isolamento estão atendidas. |
| Validade da sessão do app | 30 dias de inatividade, como em RNF-017; a seção 10.6 dizia 14 dias. |
| Modelo de identidade do login único no app | Identidade global e `acesso_empresa` (documento 2, seção 5.1); seção 10.8. Sai da lista de pendências. |
| P11 (login por telefone e código) | Não afeta a arquitetura; quando entrar, fica dentro de `plataforma/autenticacao` (gatilho (a) do ADR-011). |

**Resolvido pelo sumário 0.5 e aplicado aqui**

| Tema | Como ficou |
|---|---|
| Configuração por segmento (sumário 3.7, R7) | [ADR-010](#adr-010--configuração-por-segmento-com-catálogos-fechados) e seção 5.6: núcleo canônico fixo, blocos de catálogo fechado, catálogos de condição e efeito no código e no banco. |
| Regras configuráveis (`regra_automacao`, `disparo_regra`) | Pacote `configuracao`; avaliação por evento e por rotina diária idempotente na fila River `regras`, sem cascata e sem efeito em outra conta (seções 4.3, 5.3 e 5.6). |
| Campos extras da OS (`campo_extra_os`, `ordem_servico.campos_extras`) | Dono OS; jsonb validado pela definição em Go e em `pacotes/dominio`; no cache de configuração do app, validados antes de enfileirar o comando; fora da conexão entre contas e dos indicadores (seções 5.6, 6.6, 7, 9.3, 10.2, 10.3 e 10.5). |
| Fase de ciclo de vida e OS `implantacao` | `ativo.colocado_em_operacao` ativa os planos por evento; verificação de `colocar_em_operacao` por interface, sem ciclo ATV → OS (regra 11 da seção 5.3); exclusão fixa dos KPIs (seção 5.6); fluxo E2E 15. |
| Pacotes Go `projetos`, `contratos` e `configuracao` (sumário 3.4) | Mapa de módulos (seção 5.1), dependências e reações (seção 5.3), repositório (seção 7). |
| PRJ e CTR (Fase 2) | `projetos` usa `ordemservico` pelos contratos internos e reage a `ordem_servico.*`; `contratos` lê as paradas da OS e acompanha prazos por evento; nenhum dos dois é importado por OS (seção 5.3); fluxo E2E 16; entitlements `projetos` e `contratos_cliente` (seção 13.4). |
| Pacotes de segmento (`pacote_segmento`, `aplicacao_pacote`) | Catálogo global versionado no repositório, só leitura para as contas; cópia idempotente que só cria o que falta (seção 5.6); fluxo E2E 17; verificação de catálogo global só leitura (seção 9.3). |

**Resolvido pelo sumário 0.4 e aplicado aqui**

| Tema | Como ficou |
|---|---|
| Escopo externo | Valores `interno`, `prestador`, `locatario`, `link`; `prestador` substitui `fornecedor` na tabela de escopos, no SQL das políticas e no ADR-008 (seção 6.6). |
| Permissão de liberação | `os.encerrar`, nunca concedida a `prestador`, `locatario` nem ao link (seções 6.6, 6.7 e 8.3). |
| Ação `assumir` | Rota `POST /v1/ordens-servico/{id}/assumir` e comando `AssumirOS` no app, gravando `data_reconhecimento` (seções 8.3 e 10.3). |
| Rotas do link de uso único | Token no fragmento (`/l/#<token>`); `ativar` com token no corpo; depois `atual`, `atual/anexos`, `atual/registrar-execucao`. Nenhum `{token}` no caminho (seção 6.7). |
| `GET /v1/versao` | Exceção aceita à regra de plural (seção 8.1). |
| Versão visível e histórico de versões | RNF-039 e RNF-040 citados na seção 9.5. |
| Módulo SRV | Pacote `servicos` nos contêineres, nos pacotes, nas dependências (SRV → ADM, ATV, TER, OS; reações por evento), nos eventos `conexao_empresa.*` e `vinculo_os.*`, nos fluxos E2E 11 a 13 e nos entitlements do plano Serviços (seções 4, 5, 9 e 13.4). |
| Conexão entre contas (R5) | Sincronizador com papel `tenaz_sinc`, transações separadas por conta, outbox nos dois lados, lista fechada de campos, dono de dado, reprocessamento, SLO de atraso p95 abaixo de 60 s e teste com três contas (seções 6.8, 9.3, 14 e ADR-009). |
| App do técnico da subcontratada | Mesmo Tenaz Mobile, um SQLite por conta, troca de conta e push de tarefa nova (seção 10.8). |
| Técnico cobrado | Técnico da subcontratada conectada não conta no cliente; plano Serviços conta só o perfil `tecnico` da conta da subcontratada (seções 6.6 e 13.4). |

**Resolvido na versão anterior e mantido:** nomes dos pacotes Go, nomes dos eventos, webhooks no MVP (Profissional), SLA e SLO por plano, origem `iot`, oficina externa em TER, cadastro de fornecedor no TER, TER e LOC nos contêineres e testes, escopo externo restritivo e link de uso único (ADR-008), controle de versão (seção 9.5), cobrança de técnicos e veículos.

**Continua em aberto**

| # | Pendência | Onde |
|---|---|---|
| P7 | Local da cópia de recuperação de desastre regional. Decisão do dono do produto, com apoio jurídico; opções na seção 13.5. Até a decisão, o RTO de região de 4 horas é meta, não compromisso contratual. | Seções 13.5 e ADR-006 |
| P8 | CFOPs de remessa e retorno a subcontratada, envio para conserto e locação. Afeta validações de EST e LOC, não a arquitetura. | Sumário, seção 6 |
| P4 | Preço do plano Serviços, desconto de canal por conexão ativa e em quais planos do cliente fica o entitlement `conexao_prestador`. | Seção 13.4 |
| — | Decisão final e justificativa de `vinculo_os` (uma linha em cada conta, proposta aqui, ou registro na plataforma) e definição de `envio_conexao`, `recebimento_conexao` e da seleção de ativos compartilhados em `2_ARQUITETURA_BANCO_DADOS.md`. | Seção 6.8 |
| — | Nome da permissão de gerenciar conexões (`srv.gerenciar_conexao`, proposta) e rotas exatas das ações de `/v1/conexoes-empresa` a registrar em `1_ESPECIFICACAO_REQUISITOS.md` e `3_ESPECIFICACAO_APIS.md`. | Seção 6.8 |
| — | Prova de conceito do papel `tenaz_sinc`: privilégio por coluna combinado com RLS, chave estrangeira adiável na criação da OS vinculada, custo da política `conexao_sinc`. | Seção 6.8 |
| — | Regras de transição remota (cancelamento no cliente com OS em execução na subcontratada, janela de 24 h após o encerramento) a fechar em `4_FLUXOS_PROCESSO.md`. | Seção 6.8 |
| — | **ADR-011 a aplicar fora deste documento (versão 0.6):** papel `tenaz_autenticacao`, privilégios, política restritiva de auditoria, trigger de revogação e função de purga no documento 2 (seções 2.2, 2.6, 2.11, 5.1, 9.1, 9.2 e 17.3) e nas migrações; credencial própria, segredo lido só pela API e o novo papel no modelo de ameaças em `9_SEGURANCA_CONFIABILIDADE.md`; resposta de `POST /v1/usuarios` quando o convite não é emitido, em `3_ESPECIFICACAO_APIS.md`. | Seção 6.9 |
| E01-14 | Pendência do documento 2, seção 17.3: toda tabela futura de `plataforma` sem RLS que não seja catálogo global de leitura segue a regra da emenda do ADR-011 (dono próprio `NOLOGIN` sem membros de runtime, nunca `tenaz_dono`). A primeira é `conexao_empresa`, lida pelo `tenaz_sinc` (ADR-009). Em E04 decidir, numa emenda do ADR-009, se ela usa o próprio `tenaz_credencial` ou um dono só da conexão, e estender a ela a verificação "Dono das credenciais" da seção 9.3. Até lá, a verificação de dono já quebra o build se a tabela nascer com dono `tenaz_dono` fora da lista de catálogos | Seções 6.8, 6.9 e 9.3; ADR-009 e ADR-011 |
| — | **Herança do `tenaz_manutencao` (risco residual da emenda do ADR-011).** Membro de `tenaz_dono` com herança, ele não lê mais as credenciais, mas herda o DDL do dono comum: pode apagar e recriar as tabelas de credencial (o esquema `plataforma` é de `tenaz_dono`), desligar o trigger `trg_acesso_empresa` e, nas tabelas de `app`, desligar o `FORCE ROW LEVEL SECURITY` ou criar política para si e ler dado de qualquer empresa. Opções a avaliar pelo `tenaz-banco-dados` com o `tenaz-isolamento`: (a) tirar o `tenaz_manutencao` de `tenaz_dono` e dar a ele só `EXECUTE` em funções `SECURITY DEFINER` estreitas (criar e remover partição) mais o privilégio `MAINTAIN` do PostgreSQL 17 para `VACUUM` e `ANALYZE`; (b) mantê-lo membro, com credencial só no pipeline, como o `tenaz_migrador`, e sem uso fora dele. Prazo: antes da primeira implantação em homologação. Registrar o achado na lista do documento 2 | Seções 6.4 e 6.9; documento 2, seções 2.2 e 2.7 |
| — | **Código sqlc do ADM fora do módulo (desvio da seção 5.2).** As consultas de `adm/internal/repo/adm.sql` são geradas em `internal/plataforma/db/consultas`, junto das de auditoria e de autorização, e qualquer módulo poderia chamá-las. Mover para `adm/internal/repo/gerado`, criado sobre `db.Tx(ctx)`, como já faz `plataforma/autenticacao`, antes do primeiro módulo de E02. Responsáveis: `tenaz-backend-go` com `tenaz-banco-dados` | Seções 5.2 e 5.3 (regra 2) |
| — | Regra do `depguard` para a fronteira de pacotes da seção 5.2: só `internal/servidor` importa `<modulo>/rotas`, só a composição do worker importa `<modulo>/jobs`, e nenhum módulo importa o código gerado de outro. Entra com o segundo módulo, junto da lista de importações por módulo da regra 5 da seção 5.3 | Seções 5.2 e 5.3 |
| — | Comparação do contrato com a branch principal por `oasdiff` e documentação pública gerada do contrato (seção 2.5, itens 3 e 5) ainda não estão na esteira. A primeira precisa estar antes do primeiro integrador externo usar `/v1` | Seção 2.5 |
| A-38 | Divergência de `3_ESPECIFICACAO_APIS.md`: a seção 18.8 dele lista `ambiente` em `GET /v1/versao`. Posição deste documento (seção 9.5) e do servidor: o endpoint público devolve só `versao`, `commit` e `data_build`, como o item 6 de RNF-039; o ambiente do "Sobre" vem do build do cliente. Falta alinhar a seção 18.8 do documento 3 (`tenaz-api-contrato`) e decidir o destino do tipo `versao.Detalhe`, que tem `ambiente` e não é usado por nenhuma rota | Seção 9.5 |
| — | Política do link com as quatro condições (revogado, usado, expirado, excluído) em todas as tabelas que aceitam o escopo `link`, e `usado_em` como última escrita de `registrar-execucao`, a escrever no documento 2 na etapa do link (TER). | Seção 6.6 |
| — | Política `escopo_externo` e índices de `item_contrato_locacao` a registrar em `2_ARQUITETURA_BANCO_DADOS.md`; decidir se o locatário é gravado na OS. | Seção 6.6 |
| — | Contrato de payload de eventos em `plataforma/eventos/contratos` (regra 6 da seção 5.3) a validar em código. | Seção 5.3 |
| — | **Agendamento do selo da trilha (versão 0.8).** Sem o River (ADR-005), o `tenaz selar-auditoria` depende de a infraestrutura dispará-lo uma vez por dia, depois de 01:00 UTC: em produção, uma tarefa agendada que só receba a credencial `TENAZ_BANCO_URL_JOB` e as do S3; no desenvolvimento, à mão ou por `ferramentas/tarefas.sh`. Responsável: `tenaz-plataforma-sre`. Quando o `worker` existir, o selo vira job agendado da fila | Seções 4.1 e 4.3 |
| — | **Credencial do S3 em produção (versão 0.8).** O cliente usa só chave de acesso e chave secreta estáticas (`TENAZ_S3_CHAVE_ACESSO`, `TENAZ_S3_CHAVE_SECRETA`). Na AWS, o caminho preferível é a credencial temporária do papel IAM da tarefa, sem segredo de longa duração. Decidir com o `tenaz-plataforma-sre` antes da primeira implantação em homologação: aceitar a chave estática no Secrets Manager com rotação, ou aceitar também a cadeia de credenciais da AWS quando as duas variáveis estiverem vazias | Seção 2.4; ADR-006 |
| — | **`go.mod` com dependências diretas marcadas como `// indirect` (versão 0.8).** `aws-sdk-go-v2`, `service/s3`, `credentials` e `smithy-go` são importados por `internal/plataforma/armazenamento`, mas o `go.mod` os marca como indiretos: falta `go mod tidy`. Rodar e incluir a conferência (`go mod tidy -diff`) em `conferir`. Responsável: `tenaz-backend-go` | Auditoria do ADR-012 |
| — | **Listagem das empresas pelo `tenaz_job` (versão 0.8).** O selo lista as empresas pela política `empresa_listagem_job` (linha inteira de `app.empresa` sem contexto), e não por `app.empresas_ativas()`, que devolve só IDs (seção 6.4). Recomendação deste documento: a rotina usa a função de IDs, e a política de leitura sem contexto deixa de existir ou fica restrita à coluna `id` por privilégio de coluna. Decisão com o `tenaz-banco-dados` e o `tenaz-isolamento`, antes do segundo job que atravesse empresas. A função devolve só as empresas ativas; o selo precisa também das suspensas e das canceladas dentro da retenção, o que pede um parâmetro ou uma função irmã | Seção 6.4 |
| — | **Janela de tolerância do refresh rotativo (sugestão do `tenaz-mobile-offline`; decisão do `tenaz-arquiteto` com o `tenaz-isolamento`, sem implementar).** Se o app morre entre a resposta de `/oauth/token` e a gravação no cofre, o refresh rotativo já foi consumido no servidor, e o próximo uso do refresh antigo conta como reuso e derruba a família: o técnico precisa entrar de novo, com conexão (a fila não se perde, seção 10.6). Opções: (a) manter o rigor atual; (b) janela curta (proposta: até 60 s) em que o refresh **imediatamente anterior**, apresentado pela mesma sessão e pelo mesmo aparelho, devolve de forma idempotente o mesmo par recém-emitido, sem nova rotação, e qualquer outro reuso continua derrubando a família. O risco de (b) é um refresh roubado valer por mais 60 s; a mitigação é amarrar ao aparelho e registrar o uso na trilha. Afeta a seção 4.3 do documento 3 e o modelo de ameaças do documento 9 | Seções 6.9 e 10.6 |
| — | **Outbox da conta que sai na troca de empresa.** Com uma credencial por aparelho, a fila de uma conta só sai quando ela é a ativa. Avaliar com o `tenaz-mobile-offline` se a troca de empresa (que já exige conexão) envia antes a outbox pendente da conta atual, para não deixar comando parado em conta inativa | Seção 10.8 |
| — | Escolha da ferramenta de release (release-please ou git-cliff) e da janela de versões de app aceitas pela API. | Seção 9.5 |
| — | **Divergência (versão 0.5):** RF-ADM-016 diz que, na conexão entre contas, o efeito de uma regra (o ativo `parado`) chega à outra conta "pelo campo de status já compartilhado", mas o grupo "Cadastro do ativo" de `campos_conexao.yaml` não tem o status do ativo. Decidir se o status entra na lista (revisão de segurança, R5) ou corrigir o requisito. | Seção 6.8 e [1_ESPECIFICACAO_REQUISITOS.md](1_ESPECIFICACAO_REQUISITOS.md) |
| — | Campos extras não trafegam na conexão (sumário 3.6): a regra do cliente com `os_concluida_com_campo` (ex.: "Rompeu lacre?" no caso dos radares) não dispara para OS executada por subcontratada conectada. Só se resolve com um grupo novo na lista fechada, decidido no sumário. | Seção 5.6 e sumário 3.6 |
| — | Dono de `cliente_servico` agora que a tabela existe também em contas com CTR sem SRV: continua SRV com a exceção estreita da regra 9 da seção 5.3 (proposta aqui) ou passa a cadastro comum. Decisão no documento 2. | Seção 5.3 |
| — | Definição em `2_ARQUITETURA_BANCO_DADOS.md` de `regra_automacao`, `disparo_regra` (chave única pelo fato), `campo_extra_os` (dono OS, índice GIN em `campos_extras`), `pacote_segmento` (tabela global), `aplicacao_pacote` e das tabelas de PRJ e CTR, com o módulo dono de cada uma. | Seção 5.6 |
| — | Marca "do prestador" na definição de campo extra (RF-OS-021, proposta) ainda não está nos dados principais do requisito; a lista de acesso externo da seção 6.6 depende dela. | Seções 5.6 e 6.6 |
| — | Regra de conflito do campo extra obrigatório criado com o aparelho offline, a numerar em `4_FLUXOS_PROCESSO.md`; comando `RegistrarCamposExtras` e rotas a registrar em `3_ESPECIFICACAO_APIS.md`. | Seções 10.3 e 10.5 |
| — | Interface `atv.VerificadorEntradaOperacao` (regra 11 da seção 5.3) e `depguard` da exceção da regra 9, a validar em código. | Seção 5.3 |
| — | Nomes dos entitlements `projetos`, `contratos_cliente`, `regras_ativas_max` e `campos_extras_max` (proposta); o preço de PRJ e CTR segue a P4. | Seção 13.4 |
| P9 | Regras metrológicas dos radares: as regras do pacote `fiscalizacao_eletronica` nascem `inativa` até a decisão. Não afeta a arquitetura. | Sumário, seção 6 |
| P10 | Acesso próprio do cliente final do contrato (ex.: prefeitura). Se aprovado, exige um escopo novo além dos quatro do sumário 3.6 e revisão do ADR-008; até lá, a medição sai em PDF e XLSX. | Sumário, seção 6 |
| — | **Resíduos da troca de nomes da versão 0.4: resolvidos na versão 0.5.** As seções 6.4 e 6.8 passaram a citar `app.empresas_ativas()`, `app.conexao_atual()` e `app.escopo = 'interno'`, conforme o documento 2. Continuam abertos dois itens do mesmo bloco: o diagrama de sequência da seção 6.7 usa `#35;` no texto de uma mensagem, o que a regra dos diagramas proíbe; e o sumário 3.4 dá o exemplo `internal/ordemservico`, enquanto a seção 7 usa `backend/internal/modulos/ordemservico`. | Seções 6.7 e 7 |
