# Tecnologia, arquitetura e confiabilidade em CMMS/EAM e plataformas de frota (benchmark para o Tenaz, nov/2026)

> Notas de pesquisa para o documento de arquitetura e os requisitos não funcionais (RNF) do Tenaz (SaaS multi-tenant de manutenção de equipamentos e frota; stack Go + React/TS + PostgreSQL + Docker).
> Data da pesquisa: 12/11/2026. Fontes datadas quando relevante. Afirmações sem fonte foram para "Inferências" ou "Gaps".
> Contexto de mercado relevante: a Autodesk anunciou a compra da MaintainX em 28/05/2026 (cerca de US$ 3,6 bi, em dinheiro) e concluiu a aquisição em 03/08/2026, criando a divisão "Autodesk Operations Solutions" — [SiliconANGLE](https://siliconangle.com/2026/05/28/autodesk-acquire-maintainx-3-6-billion-push-operations/); [TipRanks – conclusão](https://www.tipranks.com/news/company-announcements/autodesk-completes-acquisition-of-maintainx-to-expand-capabilities); [Nasdaq/PR Newswire](https://www.nasdaq.com/press-release/autodesk-acquire-maintainx-advancing-unified-platform-operations-2026-05-28)

---

## 1. Multi-tenancy no PostgreSQL (shared schema + tenant_id/RLS vs schema por tenant vs banco por tenant) e como os líderes fazem

### Takeaway
A AWS chama os três modelos de **pool** (tabelas compartilhadas + RLS), **bridge** (schema ou banco por tenant na mesma instância) e **silo** (instância por tenant). Para um SaaS de mid-market com milhares de tenants pequenos e médios, a opção mais defensável é **pool com `tenant_id` em todas as tabelas + RLS forçado**, com a possibilidade de mover um tenant grande para silo depois ("tiering"). Os fornecedores publicam pouco sobre isso: a Samsara se diz multi-tenant com separação lógica, o Fiix trata isolamento de tenant como "crítico" num modelo de responsabilidade compartilhada, e uma página antiga da UpKeep fala em "banco separado por aplicação".

### Cited Findings
- A AWS Prescriptive Guidance descreve três modelos aplicáveis ao PostgreSQL. **Silo** é uma instância ou cluster por tenant. **Bridge** são bancos ou schemas separados. **Pool** são tabelas compartilhadas numa única instância e schema — [AWS Prescriptive Guidance – Decision matrix](https://docs.aws.eu/prescriptive-guidance/latest/saas-multitenant-managed-postgresql/matrix.html); [AWS – visão geral](https://docs.aws.amazon.com/prescriptive-guidance/latest/saas-multitenant-managed-postgresql/welcome.html)
- No modelo pool a AWS considera o RLS **necessário** para manter o isolamento entre tenants — [AWS Database Blog – Multi-tenant data isolation with PostgreSQL RLS](https://aws.amazon.com/blogs/database/multi-tenant-data-isolation-with-postgresql-row-level-security)
- O modelo bridge "sofre das mesmas preocupações de noisy neighbor e de isolamento de performance que o pool" — [AWS – PostgreSQL bridge model](https://docs.aws.amazon.com/prescriptive-guidance/latest/saas-multitenant-managed-postgresql/bridge.html). Ou seja, schema por tenant traz custo operacional (migrações × N schemas) sem ganho real de isolamento de recursos.
- Por padrão, o dono da tabela **não** é restringido pelas policies de RLS. É preciso `ALTER TABLE ... FORCE ROW LEVEL SECURITY` — [AWS Database Blog](https://aws.amazon.com/blogs/database/multi-tenant-data-isolation-with-postgresql-row-level-security); [Bytebase – RLS footguns](https://www.bytebase.com/blog/postgres-row-level-security-footguns/)
- Armadilhas práticas de RLS com pool de conexões, segundo blogs de praticantes, não verificadas contra a documentação oficial:
  - `SET` em nível de sessão vaza o contexto do tenant para a próxima requisição na mesma conexão do pool.
  - Com PgBouncer em modo transaction, deve-se usar `SET LOCAL` ou `set_config('app.current_tenant_id', x, true)` dentro de transação explícita.
  - O `server_reset_query` é ignorado por padrão em transaction pooling.
  - Atrás de um pooler, `current_user` não identifica o tenant.
  — [Bytebase](https://www.bytebase.com/blog/postgres-row-level-security-footguns/); [Permit.io – RLS guide](https://www.permit.io/blog/postgres-rls-implementation-guide); [PlanetScale – "RLS sounds great until it isn't"](https://planetscale.com/blog/rls-sounds-great-until-it-isnt)
- Armadilhas de performance de RLS, mesmas fontes de praticantes:
  - falta de índice em `tenant_id`;
  - funções não `LEAKPROOF` em policies impedem o uso de índice;
  - subqueries em policies podem ser avaliadas por linha.
  Para diagnóstico: `EXPLAIN (ANALYZE, SETTINGS)` — [Bytebase](https://www.bytebase.com/blog/postgres-row-level-security-footguns/); [DEV – "Why we stopped using RLS after 1M rows"](https://dev.to/canu/why-we-stopped-using-row-level-security-after-1m-rows-4dpb)
- Um artigo da comunidade defende que silo, pool e bridge sejam escolhidos **por camada** (banco, cache, storage, compute), não para o sistema inteiro. É opinião de comunidade, não orientação da AWS — [DEV – Multi-tenancy layer by layer](https://dev.to/bilelsalemdev/multi-tenancy-layer-by-layer-from-row-level-security-to-whole-cluster-isolation-4650)
- **Samsara**: arquitetura multi-tenant com "separação lógica de dados" e VPC na AWS — [Samsara Security](https://www.samsara.com/legal/security)
- **Fiix (Rockwell)**: "como plataforma multitenant, o isolamento de tenant é tratado como crítico", num modelo de responsabilidade compartilhada com o cliente — [Fiix Trust Portal](https://www.fiixsoftware.com/trust-portal)
- **UpKeep**: a página de segurança diz que "cada aplicação recebe seu próprio banco de dados separado e com controle de acesso". A página parece antiga e não fica claro se "aplicação" quer dizer cliente — [UpKeep – How secure is my data](https://help.onupkeep.com/en/articles/118925-how-secure-is-my-data-on-upkeep)
- **MaintainX**: as vagas publicadas citam stack TypeScript, React, Node.js e **MySQL** (não PostgreSQL), mais React Native e GraphQL. Há uma vaga de "Senior Software Developer, Compliance and Multi-Region", o que indica implantação em múltiplas regiões — [MaintainX Careers – Full-Stack Engineer](https://www.getmaintainx.com/jobs/full-stack-engineer-0d6fbb36); [startup.jobs – Compliance & Multi-Region](https://startup.jobs/senior-software-developer-compliance-and-multi-region-maintainx-company-7973177)
- **Fleetio**: monólito Ruby on Rails, que também serve a API do app mobile, escolhido "por velocidade e baixa complexidade". Usa PostgreSQL e React — [Rails Foundation – Fleetio](https://rubyonrails.org/foundation/fleetio); [jobs.rubyonrails.org – Fleetio](https://jobs.rubyonrails.org/jobs/742)

### Inferências
- Recomendação para o Tenaz: **pool + RLS** como padrão. Detalhamento:
  - `tenant_id UUID NOT NULL` em todas as tabelas de negócio.
  - Chaves compostas e índices começando por `tenant_id`.
  - `ENABLE` e `FORCE ROW LEVEL SECURITY`.
  - A aplicação Go conecta com um papel sem `BYPASSRLS` e que não é dono das tabelas.
  - O middleware abre transação e executa `set_config('app.tenant_id', $1, true)`.
  - Defesa em profundidade: um filtro explícito por `tenant_id` no repositório (sqlc/pgx), além do RLS.
  - Teste automatizado de vazamento entre tenants no CI.
  - Migrações (goose/atlas) com papel separado.
- Prever um caminho de **tiering**: um campo `tenant.placement` (pool/silo) e um roteador de conexões no Go permitem mover clientes enterprise ou "noisy" para um banco dedicado sem reescrever código. Isso só funciona se nenhuma query depender de joins entre tenants.
- Schema por tenant não compensa no mid-market. Multiplica o tempo de migração e o catálogo do Postgres, e a AWS reconhece que não resolve noisy neighbor.
- Escala, por inferência e sem benchmark citado: um Postgres gerenciado bem dimensionado (com réplicas de leitura) atende com folga milhares de tenants de mid-market no modelo pool. Citus ou sharding por `tenant_id` só se justificam com volumes muito acima disso. Projetar `tenant_id` como chave de distribuição desde o dia 1 mantém essa porta aberta.
- Dois líderes de mercado (MaintainX em MySQL/Node e Fleetio em monólito Rails/Postgres) mostram que um **monólito modular** bem operado basta para chegar a escala de líder. Microserviços não são pré-requisito de robustez.

### Gaps
- Nenhum fornecedor (MaintainX, UpKeep, Fiix, Fleetio, Samsara) publicou post de engenharia descrevendo o modelo de tenancy do banco. As afirmações acima vêm de páginas de segurança e de vagas.
- A página de best practices da AWS (`.../best-practices.html`) retornou 404 no fetch. Os limites de escala por modelo não foram obtidos de fonte primária.
- Não encontrei benchmark confiável e atual do custo de RLS em Postgres 16/17 com milhões de linhas por tenant.

---

## 2. Mobile offline-first para técnicos de campo

### Takeaway
Os líderes tratam "offline" de forma **limitada e pragmática**. A MaintainX faz cache local apenas de **ordens de serviço** atribuídas ao técnico ou à equipe, com limite de quantidade e sincronização depois. A empresa não documenta resolução de conflitos e o offline não existe na versão web. Para o Tenaz, React Native/Expo com SQLite local e um motor de sync com servidor autoritativo (PowerSync, ou implementação própria por outbox) é a escolha mais robusta. O Realm/Atlas Device Sync chegou ao fim de vida em 30/09/2025 e deve ser descartado.

### Cited Findings
- **MaintainX Offline Mode**, documentado na central de ajuda:
  - Disponível só no app mobile, para Full User e Administrator.
  - Faz cache apenas de work orders: as atribuídas ao usuário, ou ao usuário e às suas equipes.
  - Ordena por data de vencimento até um "Cache Limit".
  - O cache persiste a reboot e bateria descarregada.
  - Ativa automaticamente sem conexão. É preciso estar online para criar o cache.
  - A sincronização é automática ou manual conforme a configuração.
  - Não permite trocar de organização offline.
  - Não documenta resolução de conflitos.
  — [MaintainX Help – Offline mode](https://help.getmaintainx.com/offline-mode)
- MaintainX usa React Native, citado como requisito ou diferencial nas vagas, com TypeScript e GraphQL — [MaintainX Careers](https://www.getmaintainx.com/jobs/full-stack-engineer-0d6fbb36); [startup.jobs](https://startup.jobs/intermediate-full-stack-developer-maintainx-company-7973095)
- **Realm**: a MongoDB anunciou em set/2024 a descontinuação do Atlas Device Sync e dos Atlas Device SDKs (Realm), com fim de vida em **30/09/2025** — [RxDB – Realm alternative](https://rxdb.info/articles/alternatives/mongodb-realm-alternative.html); [Couchbase – alternativa ao Device Sync](https://www.couchbase.com/blog/couchbase-mobile-alternative-to-mongodb-sync/)
- **PowerSync** mantém um **SQLite local** no cliente sincronizado com Postgres, MongoDB ou MySQL no backend. Funciona totalmente offline, e as escritas passam por um backend **autoritativo** (o seu) — [PowerSync – ElectricSQL vs PowerSync](https://powersync.com/blog/electricsql-vs-powersync); [PowerSync + Supabase](https://www.powersync.com/blog/offline-first-apps-made-simple-supabase-powersync). Fontes do próprio fornecedor.
- **Electric (antigo ElectricSQL)** foi reescrito como sync de leitura Postgres → cliente ("shapes"). A persistência local e as escritas offline ficam a cargo do desenvolvedor. Em 2026 a empresa se reposicionou como "plataforma de agentes baseada em sync" — [PowerSync – Electric-next vs PowerSync](https://powersync.com/blog/electricsql-electric-next-vs-powersync); [pistack – sync engines 2026](https://www.pistack.xyz/posts/2026-10-21-local-first-sync-engines-electricsql-powersync-triplit/). Fonte parcialmente concorrente.
- Um guia de 2026 recomenda engines como Electric, PowerSync ou Zero para apps já baseados em Postgres que querem offline **sem** remodelar os dados em CRDTs — [alexcloudstar – Local-first guide 2026](https://alexcloudstar.com/md/blog/local-first-software-developer-guide-2026/)
- Alternativas citadas ao Realm: Couchbase Mobile (Sync Gateway self-hosted ou Capella), ObjectBox Sync e RxDB (JS/TS, bom para React Native e browser). Todas as fontes são dos próprios fornecedores — [RxDB](https://rxdb.info/articles/alternatives/mongodb-realm-alternative.html); [ObjectBox](https://objectbox.io/category/data-sync/page/2/)

### Inferências
- Arquitetura recomendada para o app de campo do Tenaz:
  - React Native (Expo) com SQLite (expo-sqlite ou op-sqlite).
  - **Outbox de comandos idempotentes**: cada ação offline é um comando com UUID gerado no cliente, como "apontar horas", "concluir checklist", "anexar foto" ou "registrar leitura de horímetro/odômetro".
  - O backend Go aplica os comandos com chave de idempotência.
  - **Pull incremental** por `updated_at`/versão (ou PowerSync com "sync rules" por tenant e usuário).
- Resolução de conflitos: evitar CRDT genérico.
  - Usar **last-writer-wins por campo** para campos simples.
  - Usar **append-only** para eventos (apontamentos, comentários, fotos, leituras), que por natureza não conflitam.
  - Usar **regras de domínio** para transições de status da OS, por exemplo: OS cancelada no servidor bloqueia o "concluir" vindo do offline e gera uma notificação.
- Escopo do offline: superar a MaintainX cacheando também ativos, checklists/procedimentos, peças do almoxarifado e manuais (PDF) das OS atribuídas, e permitindo criar OS e solicitações offline. É um diferencial real para operações em áreas rurais, mineração, agro e rodovias no Brasil.
- PWA com IndexedDB serve para o portal web e para solicitantes. No iOS o storage do PWA é menos previsível (inferência; não verificado nesta pesquisa), então o app nativo é preferível para técnicos.
- Fotos e anexos: upload em fila separada, retomável, com compressão no device e URLs pré-assinadas (S3 compatível).

### Gaps
- Não encontrei a stack de sync da UpKeep nem detalhes técnicos internos da MaintainX (biblioteca de storage local, protocolo de sync).
- Não há benchmark independente PowerSync vs alternativas. A maior parte do material comparativo é escrito pelos próprios fornecedores.
- Não verifiquei se o PowerSync tem SDK/serviço self-hosted compatível com backend Go. O backend de escrita é do desenvolvedor, então em princípio é agnóstico, mas isso precisa ser confirmado na documentação.

---

## 3. IoT, monitoramento de condição e telemática de frota

### Takeaway
O padrão de referência é sensor/PLC → **edge gateway** (normaliza e filtra) → **MQTT (Sparkplug B)** ou **OPC UA** → broker na nuvem → banco de séries temporais → regras de alerta → abertura de OS no CMMS. O CMMS fica na zona de TI ou DMZ, nunca ligado diretamente aos controladores. Na frota, o caminho prático é **integrar com as APIs/webhooks dos telemáticos** (Samsara, Geotab), como a Fleetio faz, e não construir hardware. Para o Tenaz, **TimescaleDB** (ou tabelas particionadas no próprio Postgres) é suficiente no começo. ClickHouse só se justifica para analytics de frota em larga escala.

### Cited Findings
- Edge gateway: limpa e filtra dados, converte unidades e publica "para o norte" via MQTT ou OPC UA. OPC UA e MQTT são **complementares**: OPC UA traz o modelo de informação navegável e os perfis de segurança, MQTT traz o transporte pub/sub — [ifactory – OPC UA + MQTT](https://ifactoryapp.com/blog/opc-ua-mqtt-ai-ready-factory-sensor-data); [Automation.com – industrial MQTT networks](https://www.automation.com/article/building-industrial-mqtt-networks-edge-computing)
- **Sparkplug B** (Eclipse) define o namespace de tópicos `spBv1.0/<Group>/<MsgType>/<EdgeNode>/<Device>`, o payload e o controle de estado com mensagens de birth e death — [Opto 22 – Industrial-strength MQTT/Sparkplug B](https://www.opto22.com/articles/industrial-strength-mqtt-sparkplug-b)
- Um guia de CMMS recomenda que o CMMS se conecte ao gateway ou ao servidor OPC UA **na DMZ**, nunca aos controladores de nível 1. Também recomenda começar com vibração em 3 a 5 ativos críticos e detectar perda de dados por *gap* de sequência e timeout de heartbeat — [OxMaint – CMMS + PLC/SCADA/IoT gateway](https://oxmaint.com/blog/post/blog-post-cmms-iot-gateway-plc-scada-integration); [ifactory – PdM pipeline](https://ifactoryapp.com/predictive-maintenance/predictive-maintenance-data-pipeline-sensor-to-work-order). Fontes de fornecedores, usar como padrão de design.
- **Samsara** (referência de escala):
  - ingestão em tempo real com AWS Kinesis;
  - analytics offline com Spark/Databricks;
  - data warehouse em Spark na AWS com AWS Step Functions;
  - um time-series store **próprio** (in-house), segundo as vagas;
  - volume na casa dos **trilhões** de pontos por ano (2 trilhões ao ano num post e mais de 9 trilhões no FY24 em outro).
  — [Samsara – Engineering culture](https://www.samsara.com/blog/samsara-engineering-culture); [Samsara – Data pipelines](https://www.samsara.com/blog/data-pipelines-at-samsara); [Vaga Staff SWE TimeSeries](https://builtin.com/job/staff-software-engineer-timeseries/6781985)
- Os gateways da Samsara **armazenam dados localmente** durante falhas da nuvem e enviam ao restabelecer — [Samsara Security](https://www.samsara.com/legal/security). É o padrão store-and-forward a copiar em qualquer gateway.
- **Samsara API de falhas**:
  - alerta "Fault Code" envia webhook (condição `DeviceHasVehicleFault`, Webhooks 1.0);
  - também há endpoints REST de falhas atuais, histórico e um event feed;
  - evento beta `EngineFaultOn`;
  - recurso de IA que classifica falhas J1939 em severidade alta, média ou baixa.
  — [Samsara – Fault monitoring](https://developers.samsara.com/docs/fault-monitoring); [Samsara – Webhooks 1.0](https://developers.samsara.com/docs/webhooks-v1-reference); [Samsara – EngineFaultOn](https://developers.samsara.com/docs/enginefaulton)
- **Fleetio + Samsara**: cria "faults" a partir dos DTCs da Samsara quando "Engine Faults" está ativo, usa token de API da Samsara e casa os veículos por **VIN**. VINs errados precisam ser corrigidos manualmente — [Fleetio Help – Samsara integration](https://help.fleetio.com/en_US/samsara-integration-5)
- **Geotab**: o GO device lê códigos de falha da ECM e os envia ao MyGeotab — [Geotab – Faults overview](https://support.geotab.com/mygeotab/doc/faults-beta-user-guide). Recomenda-se uma conta de serviço dedicada, porque o acesso via API depende das permissões do usuário — [FleetRabbit – Geotab integration](https://fleetrabbit.com/article/geotab-telematics-integration) (fonte de fornecedor).
- J1939 é o padrão SAE para pesados. Gateways J1939 completos expõem SPN/FMI específicos de fabricante que o OBD-II básico não expõe — [heavyvehicleinspection – J1939 guide](https://heavyvehicleinspection.com/blog/post/ai-powered-fault-code-analysis-guide) (fonte secundária).
- **TimescaleDB vs ClickHouse**: não há benchmark independente de 2026.
  - Testes autopublicados indicam ClickHouse muito melhor em agregações e varreduras grandes (ex.: média horária em 200 ms vs 4 s) e TimescaleDB melhor em point lookups (últimas 100 leituras de um device: 50 ms vs 800 ms).
  - ClickHouse comprime muito mais (4,2 GB vs 453 MB num dataset de 2021).
  - No ClickHouse, UPDATE é uma mutação assíncrona e cara.
  — [DEV – ClickHouse vs TimescaleDB](https://dev.to/heleo/clickhouse-vs-timescaledb-the-real-performance-showdown-for-time-series-2bdo); [DEV – Why I chose ClickHouse](https://dev.to/umarali1/why-i-chose-clickhouse-over-postgresql-for-a-billion-row-analytics-platform-599f); [arXiv SciTS 2204.09795](https://www.arxiv.org/pdf/2204.09795)

### Inferências
- O mid-market brasileiro tem pouca instrumentação. O volume inicial será de **leituras de medidores** (horímetro, odômetro, temperatura, vibração em poucos ativos) e de **eventos de telemática**, não de streams de alta frequência. Proposta para o Tenaz:
  - **Fase 1**: API de ingestão HTTP/MQTT (broker gerenciado, ex.: EMQX ou AWS IoT Core) que grava numa hypertable do TimescaleDB, no mesmo cluster Postgres, com `tenant_id` e retenção e compressão nativas. Mais um motor de regras (limiar, taxa de variação, contador de uso) que gera OS e PM por condição.
  - **Fase 2**: um conector OPC UA via edge gateway parceiro (ex.: Node-RED, Kepware ou HighByte do lado do cliente) publicando Sparkplug B.
  - **ClickHouse** só se analytics de frota ultrapassar a casa dos bilhões de linhas.
- Para frota, priorizar conectores nativos com os telemáticos usados no Brasil:
  - Samsara e Geotab, já que ambos têm API documentada;
  - os players locais (Sascar, Omnilink, Positron, Cobli e outros), cuja existência e qualidade de API **não foram verificadas** nesta pesquisa.
  - Padrão de conector: **polling + webhook**, casamento por **placa/VIN/chassi**, deduplicação de DTC (SPN+FMI+ativo, com janela de tempo) e uma regra "DTC severo vira OS corretiva automaticamente".
- Store-and-forward no gateway e idempotência no endpoint de ingestão (ID de mensagem + sequência) são RNFs obrigatórios.

### Gaps
- Não encontrei documentação das APIs dos rastreadores brasileiros (Sascar, Omnilink, Cobli etc.). Recomenda-se pesquisa dedicada.
- Não encontrei descrição técnica de como MaintainX ou Fiix ingerem dados de sensores (protocolos, armazenamento).
- Não há números de custo comparativo Timescale vs ClickHouse gerenciados na região São Paulo.

---

## 4. Manutenção preditiva e IA (2025–2026): o que é real vs marketing

### Takeaway
O que está **em produção** nos líderes é IA generativa "de produtividade":
- assistente que busca em manuais e no histórico de OS (RAG);
- geração de procedimentos e checklists;
- sugestões dentro da OS;
- transcrição de voz em notas;
- report builder em linguagem natural;
- análise de causa raiz.

Preditiva "de verdade" (modelos sobre sensores) continua dependente de dados de condição que o mid-market raramente tem. Mesmo a IBM posiciona o Maximo Assistant como **somente leitura**: não cria nem edita dados.

### Cited Findings
- **MaintainX CoPilot**: assistente generativo que busca na documentação dos ativos, no histórico de OS e em interações passadas. Ajuda a pesquisar manuais, criar tarefas e diagnosticar problemas. A empresa admite que a eficácia "depende da consistência com que os usuários inserem dados" — [MaintainX Help – AI basics](https://help.getmaintainx.com/ai-basics); [Blog – CoPilot](https://www.getmaintainx.com/blog/maintainx-copilot-ai-assistant-for-maintenance)
- Outros recursos da MaintainX:
  - **Work Order Suggestions**, sugestões em linha na OS, plano Enterprise, final de 2025 — [MaintainX – What's new Nov 2025](https://www.getmaintainx.com/blog/whats-new-at-maintainx-november-2025); [Copilot suggestions](https://webflow.getmaintainx.com/blog/smarter-work-orders-with-copilot-suggestions)
  - **AI Procedure Recommendations**, que anexam um procedimento ao criar a OS, maio/2026 — [MaintainX – What's new May 2026](https://www.getmaintainx.com/blog/whats-new-at-maintainx-may-2026)
  - **Report Builder AI** e **Root Cause Analysis GA**, março/2026 — [MaintainX Help – March 2026](https://help.getmaintainx.com/help-center-updates/march-2026)
  - AI Procedure Generator e voice memos com transcrição, lançados antes de 2026 — [MaintainX – Advanced AI features](https://www.getmaintainx.com/newsroom/maintainx-unveils-advanced-ai-features)
- **IBM Maximo Application Suite 9.1** (GA em jun/2025): assistente generativo sobre watsonx.ai, parte do "Maximo AI Service" — [IBM Newsroom](https://newsroom.ibm.com/blog-enhanced-maximo-streamlines-workforce-efficiency,-investment-planning,-and-facilities-management-introduces-gen-ai-assistant)
- Segundo a documentação da IBM, o assistente **localiza e recupera** OS, solicitações de serviço e ativos aos quais o usuário tem acesso. **Não automatiza tarefas nem cria ou edita dados** — [IBM Docs – Using the AI assistant](https://www.ibm.com/docs/SSLPL8_cd/com.ibm.mam.doc/aiconfiguration/c_max_assist.html)
- Atualizações posteriores do Maximo: schema linking (NL→consulta), suporte ao modelo **gpt-oss-120b** e respostas mais rápidas — [IBM Support 7252056](https://www.ibm.com/support/pages/node/7252056); [IBM Support 7245715](https://www.ibm.com/support/pages/node/7245715)
- **Fiix MAX** (2026): assistente conversacional embutido no CMMS, segundo o fornecedor "treinado" apenas nos dados privados do cliente e no help center do Fiix — [Fiix – AI](https://fiixsoftware.com/cmms/ai/)
- Para o recurso Prescriptive Maintenance, o Fiix usa "apenas modelos de IA privados" na nuvem da Rockwell e afirma que não envia dados de clientes pela API pública da OpenAI nem usa esses dados para treinar outros modelos — [Fiix Helpdesk](https://helpdesk.fiixsoftware.com/hc/en-us/articles/27279719926164)
- **Samsara**: classificação de severidade de falhas J1939 por IA — [Samsara – Fault monitoring](https://developers.samsara.com/docs/fault-monitoring)

### Inferências
- Arquitetura prática de IA para o Tenaz, em Go + Postgres:
  1. **RAG sobre manuais e histórico de OS** com `pgvector` no mesmo Postgres. Cada chunk leva `tenant_id`, e o RLS também vale para os embeddings, o que é crítico para não vazar dados entre clientes.
  2. **Criação de OS a partir de texto, voz ou foto**: transcrição, depois um LLM com *structured output*/tool calling que preenche um rascunho de OS (ativo, sintoma, prioridade, tipo). O humano sempre confirma. Nunca gravar direto.
  3. **Consultas em linguagem natural** via ferramentas tipadas (funções Go com parâmetros validados). Evitar text-to-SQL livre.
  4. **Preditiva "honesta"**: começar por regras e estatística (limiares, tendência de leituras, MTBF/Weibull por família de ativo, DTCs recorrentes) antes de ML. Chamar isso de "manutenção baseada em condição", não de "IA preditiva".
- Requisitos de IA para RFP/LGPD, copiando o discurso do Fiix:
  - provedor com garantia contratual de não treinar com dados do cliente;
  - opção de processamento no Brasil (inferência, a verificar por provedor);
  - log de prompts e respostas por tenant;
  - possibilidade de desligar a IA por tenant.
- Diferencial possível no Brasil: RAG sobre manuais em português e sobre normas de segurança brasileiras relevantes para manutenção (por exemplo NR-12 e NR-10). A disponibilidade e o licenciamento desses textos não foram verificados.

### Gaps
- Não encontrei informação primária sobre os recursos de IA da UpKeep em 2026.
- Não há dados independentes de ROI ou acurácia desses copilotos. Os números de mercado são de marketing.
- Não está claro qual LLM a MaintainX usa nem onde os dados são processados.

---

## 5. Integrações: REST, webhooks, OpenAPI, eventos, ERPs brasileiros, SSO/SCIM

### Takeaway
O padrão dos líderes é **API REST pública documentada + webhooks assinados**, com rate limit por plano e 429 com `Retry-After` (Fleetio), mais conectores prontos (Samsara/Fleetio). SSO via **SAML** costuma ser exclusivo do plano Enterprise (MaintainX). **SCIM** é um ponto fraco até em líderes: segundo um terceiro, a MaintainX não oferece. Nos ERPs brasileiros:
- **TOTVS Protheus** ainda tem integrações de manutenção (SIGAMNT) baseadas em **Mensagem Única/TBC**, além de APIs REST genéricas;
- **Sankhya** exige que tudo passe pelo **API Gateway** com AppKey.

Ou seja, a integração com ERP brasileiro exige uma camada de conectores dedicada.

### Cited Findings
- **Fleetio API**:
  - rate limit depende do plano;
  - 429 com cabeçalho `Retry-After` em segundos;
  - webhooks **assinados** (HMAC com chave), disponíveis só nos planos Professional e Premium;
  - Bulk API limitada a 100 registros por requisição (fonte de terceiros).
  — [Fleetio – Rate limiting](https://developer.fleetio.com/docs/overview/rate-limiting); [Fleetio – Webhooks](https://developer.fleetio.com/docs/overview/webhooks); [Fleetio Help – Webhooks overview](https://help.fleetio.com/en_US/fleetio-webhooks-overview-5); [Supergood – API report card](https://supergood.ai/api-report-card/fleetio)
- Os números de rate limit da Fleetio são conflitantes: 20 req/min num anúncio de 2019 vs 50 req/min no plano Professional segundo um blog de terceiros — [Fleetio Updates 2019](https://updates.fleetio.com/date/2019/8); [usecarly – Fleetio API](https://www.usecarly.com/blog/fleetio-api/)
- **MaintainX SSO**: exige plano Enterprise e papel de Administrator. SAML é o padrão, com verificação de domínio e roteamento por domínio. Há guias para Google, Okta e Microsoft Entra ID — [MaintainX – About SSO](https://help.getmaintainx.com/about-sso-for-maintainx); [SSO setup overview](https://help.getmaintainx.com/sso-setup-overview)
- **MaintainX SCIM**: segundo a Iden, fornecedor de provisionamento com interesse comercial, a MaintainX **não oferece SCIM** em nenhum plano. Não há confirmação oficial — [Iden – MaintainX has no SCIM](https://www.idenhq.com/en/scim-tax/maintainx)
- **TOTVS Protheus SIGAMNT**: as integrações documentadas (MNT x BackOffice RM, mobilidade/parte diária, OS de pneus em frotas) usam **Mensagem Única** (mensagens como `MaintenanceOrder`) e webservice **TBC** configurado no SigaCFG. A documentação é antiga (2013 em diante) — [TDN – MNTxBackOfficeRM](https://tdn.totvs.com.br/download/attachments/547228354/documentacao-MNTxBackOfficeRM-020813.pdf?api=v2); [TDN – Integração Mobilidade SIGAMNT](https://tdn.totvs.com.br/download/attachments/547228028/Integra%C3%A7%C3%A3o_Mobilidade_SIGAMNT.pdf?api=v2); [TDN – Configurar Webservice TBC](https://tdn.totvs.com.br/download/attachments/547227958/BT%20-%20Configurar%20%20Webservice%20TBC%20%20no%20Protheus%20SigaCFG.pdf?api=v2); [TDN – MntFrotas OS Pneus](https://tdn.totvs.com.br/download/attachments/547227985/Especificacao_Projeto_MntFrotas-OrdemServicoPneus.pdf?api=v2)
- O Protheus também expõe APIs REST genéricas. Um guia de terceiros indica que o uso de APIs exige licença TOTVS e patch de API base — [DIO – API REST no Protheus](https://www.dio.me/articles/api-rest-no-protheus-conectando-o-erp-ao-mundo-digital-b4594aa19f6a); [Base – Integração TOTVS Protheus](https://base.com/pt-BR/ajuda/knowledgebase/integracao-com-o-totvs-protheus-guia-completo-de-configuracao-parte-1/)
- **Sankhya**: "todas as requisições passando pelo gateway". Registros inseridos direto no banco **não são reconhecidos**. O integrador precisa de **AppKey** e de Sankhya ID, e a paginação usa `offsetPage`. O portal tem "API de Serviços Gateway" e índice OpenAPI — [Sankhya – Requisitos de Integração Gateway](https://sankhya.com.br/arquivos/Requisitos-de-Integracao-Gateway-Final.pdf); [Sankhya Developer – API de Serviços](https://developer.sankhya.com.br/docs/api-de-servi%C3%A7os)
- **SAP**: a MaintainX e o Fiix estão listados como parceiros na SAP Store, o que mostra que integração certificada com SAP é argumento de venda enterprise — [SAP Store – MaintainX](https://www.sap.com/products/scm/partners/maintainx-inc-maintainx-maintenance-and-asset-management.html); [SAP Store – Fiix](https://www.sap.com/hk/products/scm/partners/rockwell-automation-canada-ltd-fiix-cmms-manage-maintenance-and-asset-performance.html)

### Inferências
- RNFs de integração para o Tenaz:
  - **API REST pública versionada** (`/v1`), especificada em **OpenAPI 3.1** e gerada a partir do código Go.
  - Paginação por cursor.
  - **Idempotency-Key** em POST.
  - Rate limit por tenant e plano, com 429 + `Retry-After`.
  - **Webhooks assinados** (HMAC-SHA256 + timestamp contra replay), com retry exponencial, DLQ e painel de entregas para o cliente.
  - **Transactional outbox** no Postgres para publicar eventos de domínio (OS criada, OS concluída, peça consumida, leitura registrada) sem perder mensagens.
- Camada de conectores ERP: um "integration hub" interno (workers Go consumindo o outbox) com adaptadores por ERP.
  - **TOTVS Protheus**: REST onde houver, ou Mensagem Única para MNT.
  - **TOTVS RM**: webservices TBC/REST.
  - **Sankhya**: Gateway com AppKey.
  - **Senior** e **Omie**: a verificar.
  - **SAP**: OData/BAPI via SAP BTP.
  - Fluxos típicos: cadastro de ativos e centros de custo vindo do ERP, requisição e baixa de peças, custo de OS e notas de serviço indo para o ERP.
- Diferencial: SSO **OIDC + SAML** e **SCIM 2.0** já no plano intermediário, não só no Enterprise. Isso explora uma lacuna dos líderes, já que a ausência de SCIM na MaintainX é alegação de terceiro.

### Gaps
- Não encontrei documentação de APIs REST específicas do SIGAMNT em versões recentes do Protheus (TLPP/REST). A pesquisa só trouxe documentação antiga de Mensagem Única.
- Não obtive detalhes oficiais das APIs da **Omie** (autenticação app_key/app_secret), da **Senior X Platform** nem da autenticação do gateway Sankhya (OAuth ou não).
- Não verifiquei a política oficial de rate limit e versionamento da API da MaintainX nem a da UpKeep.

---

## 6. Confiabilidade e segurança: SLA, RPO/RTO, backups, DR, observabilidade, auditoria, residência de dados, certificações, LGPD

### Takeaway
A régua de mercado é:
- **SLA de 99,9% mensal** com créditos limitados (a MaintainX dá 5% por período, teto de 25%, e exclui manutenção programada);
- **SOC 2 Type II** (e ISO 27001 no caso da MaintainX);
- criptografia em trânsito (TLS 1.2+) e em repouso (AES-256 com KMS);
- pentest anual;
- bug bounty (Samsara/Bugcrowd).

Quase **nenhum** líder publica RPO/RTO. A UpKeep tem snapshots de banco a cada 30 minutos. Para o Brasil, a LGPD **não exige** hospedagem local, mas a Resolução ANPD 19/2024 (cláusulas-padrão obrigatórias desde 23/08/2025 para transferência internacional) faz da **hospedagem em São Paulo** um diferencial comercial e jurídico concreto.

### Cited Findings
- **SLA da MaintainX**:
  - 99,9% de disponibilidade mensal, excluindo manutenção programada, quedas de terceiros e suspensões;
  - créditos só para quedas acima de 60 minutos consecutivos;
  - 5% da mensalidade por período de 30 minutos ou mais adicionais, no máximo 1 crédito por dia e teto de 25%;
  - o cliente precisa solicitar em até 30 dias;
  - os créditos são o remédio exclusivo.
  — [MaintainX – Service Level Agreement](https://getmaintainx.com/service-level-agreement)
- A MaintainX anuncia "99,9% uptime, guaranteed" na página Enterprise e mantém SOC 2 Type 2. Também anunciou ISO 27001 junto com a renovação do SOC 2 (post de 2023; não encontrei anúncios de renovação de 2024 a 2026) — [MaintainX Enterprise](https://www.getmaintainx.com/enterprise); [MaintainX – SOC 2 Type 2 renewal & ISO 27001](https://www.getmaintainx.com/blog/maintainx-announces-soc-2-type-2-compliance-renewal-and-iso-27001-certification)
- Status de terceiros: em 21/10/2026 a MaintainX mostrava 100% de uptime em 30 dias, com dois incidentes curtos de degradação no início de outubro — [Statusfield – MaintainX](https://statusfield.com/services/maintainx) (fonte secundária)
- **Samsara**:
  - AWS, TLS 1.2;
  - criptografia em repouso **FIPS 140-2** em bancos, file stores e backups, com chaves gerenciadas via AWS KMS;
  - dados replicados em múltiplos servidores;
  - sistema "projetado para failover rápido", sem números de RPO/RTO;
  - SOC 2 Type 2 (auditor Schellman) sob solicitação;
  - pentests independentes ao menos anuais, incluindo hardware e gateways;
  - bug bounty via Bugcrowd;
  - logs de auditoria de autenticação para administradores;
  - MFA para acesso de staff à produção.
  — [Samsara Security](https://www.samsara.com/legal/security)
- **UpKeep**:
  - AWS;
  - snapshots dos bancos a cada **30 minutos** e das configurações a cada 3 horas;
  - restauração automática e redeploy em múltiplos data centers;
  - Cloudflare/WAF contra DDoS;
  - acesso de staff aos dados de clientes só para suporte ou obrigação legal;
  - ISO 27001 listado (página antiga e parcialmente confusa).
  — [UpKeep – How secure is my data](https://help.onupkeep.com/en/articles/118925-how-secure-is-my-data-on-upkeep)
  - Trust center no SafeBase — [trust.upkeep.com](https://trust.upkeep.com)
  - SOC 2 Type II segundo diretório de terceiros — [trustlists.org – UpKeep](https://trustlists.org/company/upkeep/)
- **Fiix**: ISO 27017:2015; relatório SOC 2 sob solicitação; política de segurança pública revisada pela última vez em 29/10/2021 — [Fiix Trust Portal](https://www.fiixsoftware.com/trust-portal); [Fiix – Security policy](https://fiixsoftware.com/legal/security-policy/)
- **LGPD e residência**:
  - A LGPD não exige hospedagem nacional. Exige que o controlador saiba e documente para onde vão os dados. Hospedar no Brasil elimina a questão do mecanismo de transferência internacional — [Locaweb – LGPD e Cloud](https://www.locaweb.com.br/blog/temas/seguranca-digital/lgpd-cloud-computing-onde-dados-precisam-estar/)
  - **Resolução CD/ANPD nº 19/2024** (23/08/2024): cláusulas-padrão contratuais (Anexo II) devem ser adotadas **sem alterações**. Prazo de 12 meses, encerrado em **23/08/2025**. A não incorporação pode gerar fiscalização — [Lefosse](https://lefosse.com/noticias/alerta/transferencia-internacional-de-dados-prazo-para-adocao-das-clausulas-padrao-da-anpd-se-encerra-em-23-de-agosto/); [Mattos Filho](https://www.mattosfilho.com.br/unico/regulamentacao-transferencia-internacional-dados/); [Migalhas – fiscalização ANPD](https://www.migalhas.com.br/quentes/449475/envio-de-dados-ao-exterior-esta-no-radar-fiscalizatorio-da-anpd)
  - A transferência internacional de dados está no radar fiscalizatório da ANPD — [Migalhas](https://www.migalhas.com.br/quentes/449475/envio-de-dados-ao-exterior-esta-no-radar-fiscalizatorio-da-anpd)
- Concorrente nacional: não encontrei informação pública sobre hospedagem, SLA ou data center do Engeman — [Capterra – Engeman](https://capterra.com/p/13520/Engeman/)

### Inferências
- Proposta de RNFs para o Tenaz, com o objetivo de igualar ou superar os líderes:
  - **Disponibilidade**: SLO interno de 99,95% e SLA contratual de **99,9%** mensal. Créditos melhores que os da MaintainX (crédito a partir de 99,9% sem a exigência de 60 minutos consecutivos) e **status page pública** com histórico.
  - **RPO ≤ 5 min** (PITR com WAL contínuo do Postgres gerenciado) e **RTO ≤ 1 h** para falha de AZ e ≤ 4 h para falha de região. Isso já supera os snapshots de 30 minutos da UpKeep, e **publicar** RPO/RTO é diferencial, porque ninguém publica.
  - **Backups**: PITR de 7 a 35 dias, mais cópia diária **em outra região** com criptografia e conta separada (proteção contra ransomware e erro humano). Teste de restauração **trimestral** documentado.
  - **Topologia**: Postgres gerenciado Multi-AZ em **AWS sa-east-1 (São Paulo)**, ou Azure Brazil South ou GCP southamerica-east1, com réplica de leitura. DR de dados em us-east-1, coberto pelas cláusulas-padrão da ANPD, ou em uma segunda região brasileira se o provedor tiver (a disponibilidade de 2ª região no Brasil não foi verificada).
  - **Criptografia**: TLS 1.2+ (preferir 1.3), AES-256 em repouso via KMS, segredos em secret manager, e opção futura de chave por tenant.
  - **Observabilidade**: OpenTelemetry (traces, métricas e logs) no backend Go e no app, com `tenant_id` como atributo, SLOs com error budget e alertas por *burn rate*. A adoção de OTel pelos concorrentes não foi verificada; é prática de mercado.
  - **Auditoria**: trilha de auditoria imutável (append-only, quem, o quê, quando, valores antes e depois, IP e dispositivo) para entidades críticas como OS, ativos, estoque e permissões, exportável pelo cliente. Logs de autenticação como os da Samsara.
  - **Segurança**: pentest anual por terceiro com sumário executivo disponível ao cliente, programa de divulgação responsável, MFA obrigatório para staff, SAST/DAST e SCA no CI.
  - **Compliance**: roadmap com LGPD (DPA, RIPD, encarregado/DPO) no dia 1, **SOC 2 Type II** em 12 a 18 meses e **ISO 27001** depois. SOC 2 é o selo que os líderes usam em RFPs.

### Gaps
- Nenhum líder pesquisado publica RPO/RTO formais nem a região de hospedagem de clientes brasileiros.
- Não verifiquei se MaintainX, UpKeep, Fiix ou Fleetio oferecem residência de dados no Brasil. A vaga "Multi-Region" da MaintainX sugere multirregião, sem detalhe.
- O PDF comercial da Senior Cloud não pôde ser extraído, então não há SLA de referência de fornecedor brasileiro.
- Não encontrei o SLA da Samsara ("Hosted Software SLA") nem o da Fleetio.

---

## 7. Relatórios/BI: analytics embarcado, réplicas de leitura, views materializadas, exportação de dados

### Takeaway
Os líderes estão migrando de relatórios fixos para **report builder com IA** (MaintainX Report Builder AI, 2026) e mantêm times dedicados de "Reporting & Analytics" (Fleetio). A Samsara separa o caminho de tempo real (dashboards ao vivo) do warehouse analítico (Spark/Databricks). Para o Tenaz, o caminho incremental é: réplica de leitura + views materializadas por tenant, depois um warehouse/OLAP, mais exportação self-service.

### Cited Findings
- MaintainX lançou **Report Builder AI** em março de 2026 — [MaintainX Help – March 2026](https://help.getmaintainx.com/help-center-updates/march-2026)
- A Fleetio contrata para os times "Reporting & Analytics" e "Strategic Reporting & Tools" (stack Rails/Postgres), sinal de que relatórios são uma área de produto dedicada — [jobs.rubyonrails.org – Fleetio Strategic Reporting](https://jobs.rubyonrails.org/jobs/19415-senior-software-engineer-strategic-reporting-tools-fleetio); [rubyonremote – Reporting & Analytics](https://rubyonremote.com/jobs/68739-senior-software-engineer-reporting-analytics-at-fleetio)
- A Samsara move dados dos backends de produção em tempo real (que alimentam dashboards ao vivo) para um data warehouse Spark na AWS, com transformações definidas em DSL e orquestradas por Step Functions — [Samsara – Data pipelines](https://www.samsara.com/blog/data-pipelines-at-samsara)
- RLS também pode ser usado para analytics multi-tenant em Postgres — [Draxlr – RLS for multi-tenant analytics](https://www.draxlr.com/blogs/row-level-security-multi-tenant-analytics/)
- A ClickHouse publica um guia de arquitetura de SaaS multi-tenant com Postgres + ClickHouse para analytics — [ClickHouse – Multi-tenant SaaS Postgres architecture](https://clickhouse.com/resources/engineering/multi-tenant-saas-postgres-architecture)

### Inferências
- Abordagem para o Tenaz, por fase:
  - **Fase 1**: KPIs de manutenção (MTBF, MTTR, backlog, % preventiva vs corretiva, custo por ativo e por km, disponibilidade) calculados em **views materializadas** com refresh incremental ou agendado, ou em tabelas de agregados atualizadas por eventos do outbox. Consultas pesadas vão para uma **réplica de leitura**.
  - **Fase 2**: CDC (logical replication/Debezium) para ClickHouse ou um warehouse, quando o volume justificar.
  - **Analytics embarcado**: um builder próprio em React com a API de consultas tipadas, ou ferramenta embarcável (Metabase/Superset embed com filtro por tenant). Evitar dar SQL livre ao cliente.
- **Exportação de dados** como RNF de confiança:
  - export CSV/XLSX de qualquer listagem;
  - export completo do tenant (JSON/CSV + anexos), agendável;
  - conector para Power BI (OData ou endpoint de dataset), muito usado no mid-market brasileiro (inferência).
  - Isso também atende à portabilidade da LGPD e reduz o medo de lock-in.

### Gaps
- Não encontrei descrição técnica de como MaintainX, UpKeep ou Fiix implementam relatórios (engine, embarcado ou próprio).
- Não há dados sobre demanda de conector Power BI em CMMS no Brasil.

---

## 8. Preços e billing para SaaS por uso (por ativo/usuário), feature flags e planos

### Takeaway
Os líderes vendem por **usuário e plano**, com recursos premium (SSO, webhooks, offline completo, IA avançada) restritos aos planos superiores. Exemplos: SSO só no Enterprise na MaintainX; webhooks só em Professional e Premium na Fleetio; Work Order Suggestions só no Enterprise. Para billing por uso (ativos, veículos, leituras de IoT, chamadas de IA), o ecossistema de 2026 tem:
- **Stripe Billing**, que comprou a Metronome em jan/2026;
- **Lago**, open source, billing + metering;
- **OpenMeter**, só metering, hoje parte da Kong.

No Brasil, a cobrança em boleto/PIX e a NFS-e influenciam a escolha.

### Cited Findings
- **MaintainX**: SSO exige plano Enterprise — [MaintainX – About SSO](https://help.getmaintainx.com/about-sso-for-maintainx). Work Order Suggestions no Enterprise — [MaintainX – Nov 2025](https://www.getmaintainx.com/blog/whats-new-at-maintainx-november-2025). Offline disponível de Basic a Enterprise, mas não para o perfil Requester — [MaintainX – Offline mode](https://help.getmaintainx.com/offline-mode)
- **Fleetio**: webhooks só nos planos Professional e Premium; rate limit da API varia por plano — [Fleetio Help – Webhooks](https://help.fleetio.com/en_US/fleetio-webhooks-overview-5); [Fleetio – Rate limiting](https://developer.fleetio.com/docs/overview/rate-limiting)
- **Stripe Billing** é descrito como "subscription-first", com limitações em créditos pré-pagos e modelos híbridos. A Stripe adquiriu a Metronome (especialista em usage billing), negócio concluído em 14/01/2026 — [HackerNoon – usage-based billing 2026](https://hackernoon.com/best-usage-based-billing-platforms-for-ai-companies-in-2026-metronome-orb-stripe-and-alternatives); [Solvimon – billing systems 2026](https://www.solvimon.com/blog/best-billing-systems-in-2026) (fontes de terceiros ou concorrentes)
- **Lago**: billing + metering open source (núcleo AGPLv3, há divergência entre fontes sobre a licença) com versão cloud. Alguns recursos ficam no plano pago — [Lago – vs Metronome](https://getlago.com/blog/lago-vs-metronome); [Flexprice – Lago alternatives](https://flexprice.io/blog/lago-alternatives-for-usage-based-pricing)
- **OpenMeter**: metering open source, hoje parte da Kong, normalmente combinado com outro sistema de faturamento — [Dodo Payments – Lago alternatives](https://dodopayments.com/blogs/lago-alternatives/); [youngju.dev – SaaS billing 2026](https://www.youngju.dev/blog/culture/2026-05-16-saas-billing-metering-2026-stripe-billing-paddle-lago-openmeter-orb-chargebee-recurly-metronome-deep-dive.en)
- Pagamentos da UpKeep processados pela Stripe (PCI Level 1) — [UpKeep – Security](https://help.onupkeep.com/en/articles/118925-how-secure-is-my-data-on-upkeep)

### Inferências
- Arquitetura de planos e billing do Tenaz:
  1. **Catálogo de planos e entitlements no próprio banco**: tabelas `plan`, `plan_feature` (flags booleanas e limites numéricos como `max_assets`, `max_vehicles`, `sso_enabled`, `webhooks_enabled`, `ai_credits_month`) e `tenant_subscription` com overrides por tenant.
  2. Um serviço `entitlements.Check(ctx, tenant, feature)` em Go, com cache, usado no backend e exposto ao front para esconder ou mostrar recursos.
  3. **Feature flags operacionais** (rollout gradual, kill switch) separadas dos entitlements comerciais, via OpenFeature com provider (Unleash/Flagsmith self-hosted ou similar). A ferramenta específica não foi pesquisada.
  4. **Metering**: eventos de uso (ativos ativos no mês, veículos com telemática, leituras ingeridas, tokens de IA) agregados diariamente a partir de snapshots ou do outbox.
  5. **Cobrança**: gateway com boleto, PIX e cartão recorrente e **NFS-e** (por exemplo Iugu, Asaas, Vindi ou Pagar.me, nenhum verificado nesta pesquisa), ou Stripe + emissor de NFS-e. Lago self-hosted é opção se o modelo por uso ficar complexo.
- Modelo de preço sugerido para o mid-market brasileiro: **por usuário técnico + faixa de ativos/veículos**, com solicitantes ilimitados e gratuitos. Hipótese de produto: esse é um gargalo nos concorrentes, já que a MaintainX restringe o Requester. **SSO, SCIM, API e webhooks no plano intermediário** contrastariam com os líderes, que travam esses recursos no topo.

### Gaps
- Não pesquisei a tabela de preços atual (2026) de MaintainX, UpKeep, Fiix e Fleetio, nem como cobram por ativo ou veículo. Outro pesquisador pode cobrir isso.
- Não verifiquei suporte a boleto/PIX/NFS-e em Lago, Stripe Billing ou gateways nacionais.
- Não encontrei fontes sobre o sistema de entitlements que os líderes de CMMS usam internamente.
