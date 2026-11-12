# Gestão de frota e manutenção de frota: benchmark global e Brasil (pesquisa de nov/2026)

> Escopo: requisitos do módulo de frota do Tenaz (SaaS multi-tenant de manutenção de equipamentos e veículos para empresas médias brasileiras). A pesquisa foi feita em 12/11/2026. Muitos dados de preço vêm de agregadores de terceiros, e isso está indicado em cada item. Itens sem fonte confiável estão nas seções "Gaps".

---

## 1. Líderes globais (Fleetio, Samsara, Geotab, Motive, Whip Around, Teletrac Navman, Chevin FleetWave): posicionamento, funcionalidades de manutenção, preço, integrações e avaliações

### Takeaway
O mercado global tem dois polos. De um lado estão as plataformas de telemática (Samsara, Geotab, Motive), que geram dados como odômetro, horímetro, DTC, DVIR e alertas, e desde 2025 vêm adicionando ordens de serviço "leves". Do outro estão os softwares de manutenção de frota (Fleetio, Whip Around, Chevin FleetWave), que são a camada de oficina, peças, pneus e custo e se integram às telemáticas. Para o Tenaz, o modelo a seguir é o Fleetio: preço público por veículo, usuários ilimitados e integrações com telemática e cartão combustível. A ele o Tenaz deve somar a conformidade regulatória local, como faz o FleetWave.

### Cited Findings
**Fleetio (camada de manutenção, referência de produto)**
- Tem três planos por veículo/mês: Essential a US$ 4 (anual) ou US$ 5 (mensal), Professional a US$ 7 e Premium a US$ 10 (só anual). O preço varia com o tamanho da frota e o prazo do contrato, e há preço GSA para governo — [Fleetio Pricing](https://www.fleetio.com/pricing); [PricingSaaS](https://pricingsaas.com/companies/fleetio)
- Todos os planos têm usuários ilimitados, então a cobrança cresce com o número de veículos e não com o de usuários — [SpotSaaS](https://www.spotsaas.com/blog/fleetio-pricing)
- Não há plano gratuito em jul/2026 — [CostBench](https://costbench.com/software/fleet-management/fleetio/free-plan/)
- O empacotamento por tier, conforme a página de preços (algumas marcações estavam ambíguas no texto extraído):
  - Ordens de serviço: ausentes no Essential e limitadas no Professional.
  - Gestão de pneus (atividade, inspeções de pneu, garantia): só no Premium.
  - Peças avançadas (pedidos de compra, ajustes de estoque, alerta de estoque baixo, métodos de valoração, integração PartsTech): Premium.
  - API: sem acesso no Essential, 50 req/min no Professional e 250 req/min no Premium.
  - Webhooks: 10 mil eventos por 30 dias no Professional, ilimitados no Premium.
  - Apontamento de mão de obra (clock in/out, job tracking): Premium.
  - Inspeções móveis offline com templates FMCSA estão nos planos de entrada.
  - Fonte: [Fleetio Pricing](https://www.fleetio.com/pricing)
- O Fleetio está no Geotab Marketplace. Um alerta detectado no dispositivo Geotab é sincronizado no Fleetio, e o gestor dispara a partir dele um fluxo de manutenção interna ou terceirizada. Há integrações com oficinas e com emissores de cartão combustível — [Geotab Marketplace – Fleetio](https://marketplace.geotab.com/solutions/fleetio); [Fleetio App Directory – Geotab](https://www.fleetio.com/app-directory/fleet-telematics/geotab)
- Reviews no GetApp citam que o preço por veículo fica caro em frotas grandes — [GetApp](https://www.getapp.com/communication-software/a/fleetio/)

**Samsara**
- Não publica preço de lista. Terceiros estimam cerca de US$ 27 por veículo/mês ([SelectHub](https://www.selecthub.com/fleet-management-software/samsara-vs-gomotive/)) ou US$ 27 a 33 por veículo/mês com contrato de 3 anos ([Airpinpoint](https://airpinpoint.com/compare/samsara-competitors)).
- No G2 tem nota 4.5/5 em cerca de 4.983 reviews. A crítica recorrente é o preço alto para frotas pequenas — [G2 Motive vs Samsara](https://www.g2.com/compare/motive-technologies-vs-samsara)
- Lançou o recurso "fault code intelligence", que decodifica DTCs e usa IA para criar ordens de manutenção — [CCJ Digital](https://www.ccjdigital.com/technology/artificial-intelligence/article/15749298/samsaras-latest-asset-maintenance-features-utilize-ai) (página bloqueada com 403; informação vista via snippet de busca e [OxMaint](https://oxmaint.com/industries/fleet-management/samsara-geotab-motive-telematics-comparison-fleet-2026))

**Geotab (e Marketplace)**
- Em 26/02/2025, no Geotab Connect, anunciou no Maintenance Center:
  - Work Order Management: converte solicitações em OS, acompanha as OS abertas e analisa custo e downtime, com análise preditiva de veículos em risco de falha.
  - Fault Code Enrichment, feito com a Noregon: dá severidade, descrição atualizada e próximo passo recomendado ao DTC ("parar imediatamente", "na próxima parada/serviço" ou "seguro continuar"), padronizado entre montadoras.
  - O Noregon TripVision é vendido com preço específico no Geotab Order Now Marketplace.
  - Fonte: [Geotab press release](https://www.geotab.com/press-release/maintenance-announcements-connect-2025/)
- O preço depende de hardware, pacote, revendedor e região, e não há lista pública. As estimativas de terceiros variam bastante:
  - GO9+ de US$ 100 a 150, mais US$ 15 a 30 por veículo/mês.
  - Software de US$ 20 a 40 por veículo/mês e hardware de US$ 80 a 120 (TechRadar 2026, via agregador).
  - Add-ons do Marketplace têm assinatura própria.
  - Fontes: [CostBench Geotab](https://costbench.com/software/fleet-management/geotab/); [Geotab software packages](https://www.geotab.com/software-packages/)
- A Geotab tem presença e conteúdo em pt-BR, inclusive sobre gestão de pneus e controle de documentos — [Geotab BR – pneus](https://www.geotab.com/pt-br/blog/gestao-de-pneus-frota/); [Geotab BR – documentos](https://www.geotab.com/pt-br/blog/controle-documentos-conformidade-frota/)

**Motive**
- Os diagnósticos alimentam fluxos de manutenção: problemas detectados viram OS atribuídas. O fluxo é focado em manutenção corretiva, com status e atribuição. O próprio Fleetio, que é concorrente, argumenta que o Motive é menos voltado à oficina — [Fleetio alternatives – Samsara/Motive](https://fleetio.com/alternatives/samsara); [Gitnux](https://gitnux.org/best/trucking-fleet-maintenance-software/)
- Preço: cerca de US$ 35 por veículo/mês ([SelectHub](https://www.selecthub.com/fleet-management-software/samsara-vs-gomotive/)) ou a partir de US$ 25 ([Airpinpoint](https://airpinpoint.com/compare/samsara-competitors)). Segundo a TechRadar, os tiers superiores ficam entre US$ 40 e 50 e o contrato mínimo é de 12 meses — [TechRadar](https://www.techradar.com/nz/reviews/keeptruckin)
- No G2 tem nota 4.4/5 em cerca de 3.257 reviews — [G2](https://www.g2.com/compare/motive-technologies-vs-samsara)

**Whip Around**
- É uma plataforma de inspeção (DVIR digital, formulários customizáveis) e manutenção (OS, preventiva, conformidade DOT, combustível e relatórios) — [SelectHub Whip Around](https://www.selecthub.com/p/fleet-maintenance-software/whip-around/)
- Preço: Standard a US$ 5 e Pro a US$ 9 ou 10 por ativo/mês ([G2](https://www.g2.com/products/whip-around-whip-around/pricing); [PricingSaaS](https://pricingsaas.com/companies/whiparound)).
- Oferece o modelo "FixedUnlimited": preço fixo com ativos ilimitados, incluindo veículos e equipamentos — [Whip Around FixedUnlimited](https://whiparound.com/fixedunlimited-guide/); [Help – planos](https://help.whiparound.com/en/articles/6925626-how-whip-around-s-plans-work)

**Chevin FleetWave**
- Cobre o ciclo de vida completo, da aquisição ao descarte e remarketing, passando por custos, combustível, uso, manutenção, conformidade e acidentes — [Chevin – What is FleetWave](https://www.chevinfleet.com/products/fleetwave/)
- Automatiza lembretes e workflows de conformidade (renovação de licenças, recertificação de equipamentos) e tem relatórios de TCO e de ciclo de vida. O app FleetWave Driver registra defeitos, acidentes, recibos de combustível e formulários de conformidade — [Chevin](https://www.chevinfleet.com); [SoftwareAdvice UK](https://www.softwareadvice.co.uk/software/61705/fleetwave)
- Terceiros destacam o suporte a exigências de conformidade específicas de cada país — [SourceForge](https://www.sourceforge.net/software/product/FleetWave)
- Preço segundo agregadores: Lite a £3 por veículo/mês, Core a £5 e Advanced sob consulta. Não foi confirmado na página oficial — [Chevin pricing](https://www.chevinfleet.com/pricing); [FleetWave Lite – Government Fleet](https://www.government-fleet.com/news/fleetwave-lite-software-launched-to-manage-compliance-headaches)

**Padrão de mercado**
- Várias fontes descrevem a telemática como fonte de dados e o CMMS como camada de OS, peças e custo. Uma delas diz explicitamente que nenhuma das três telemáticas é um CMMS — [OxMaint – Telematics-CMMS](https://oxmaint.com/industries/fleet-management/telematics-cmms-integration-for-fleets-guide) (fonte de fornecedor de CMMS, então há viés)

### Inferences
- O preço de referência para software de manutenção puro, sem hardware, fica entre US$ 4 e 10 por veículo/mês. A telemática com hardware custa de US$ 20 a 50 por veículo/mês. O Tenaz deve se posicionar como camada de manutenção, perto do primeiro grupo, integrado às telemáticas e sem competir com elas.
- O empacotamento do Fleetio mostra o que o mercado trata como "premium": pneus, pedidos de compra e estoque avançado, apontamento de mão de obra e API com volume. O Tenaz pode diferenciar oferecendo pneus no plano base, já que no Brasil o pneu é o segundo maior custo da frota pesada.
- O enriquecimento de DTC (severidade e ação recomendada) virou recurso padrão em 2025 e 2026. O Tenaz deve ter uma regra configurável "DTC → solicitação ou OS" com severidade.
- O modelo FixedUnlimited do Whip Around e o "equipamentos + veículos" se alinham com a proposta do Tenaz de cobrir ambos.

### Gaps
- Teletrac Navman não foi pesquisado: não há dados de funcionalidades nem de preço nesta rodada.
- Não foi possível ler a documentação de API da Samsara. O caminho developers.samsara.com/docs/maintenance retornou só um título, e o índice está em https://developers.samsara.com/llms.txt, que não foi lido.
- Não foi possível ler a documentação de API do Fleetio (developer.fleetio.com/docs/overview retornou 404).
- Não há preço oficial publicado de Samsara, Motive, Geotab nem Chevin. Todos os números são de terceiros.
- Não foram encontrados reviews de G2 ou Capterra de Whip Around e FleetWave.

---

## 2. Players brasileiros: escopo (manutenção vs. telemática pura)

### Takeaway
No Brasil, a maior parte dos players de destaque é de rastreamento e telemática: Sascar/Michelin, Omnilink, Cobli e Vivo, com Maxtrack, Positron e Onixsat sem dados nesta rodada. Eles oferecem manutenção apenas no nível de "lembretes, checklist e histórico". A gestão de manutenção profunda, com OS, almoxarifado, pneus e custo, fica com softwares de gestão de frota (Gestran, Infleet, Prolog, Bsoft) e com o modelo "rede credenciada + cartão" de Ticket Log/Edenred. Há espaço para um CMMS de frota e equipamentos que integre todos esses players.

### Cited Findings
**Cobli**
- Oferece histórico de manutenção por veículo e alertas de preventiva — [Produttivo](https://produttivo.com.br/blog/sistema-de-gestao-de-frotas)
- A API pública tem endpoints de manutenção: avisos de manutenção, registros de manutenção (listar e criar) e cadastro de peças (createpart) — [docs.cobli.co](https://docs.cobli.co/reference/createpart)
- Preço: um artigo diz que atende frotas a partir de 3 veículos com planos a partir de R$ 50 por veículo/mês ([Produttivo](https://produttivo.com.br/blog/sistema-de-gestao-de-frotas)). O B2B Stack mostra "preço inicial R$ 76,00/mensal" e, em outro campo, R$ 1.990/mensal, sem deixar claro a unidade, então é conflitante — [B2B Stack Cobli](https://www.b2bstack.com.br/product/cobli)
- No Reclame Aqui a reputação é "Ótima", com nota média 9.3/10 nos 6 meses anteriores (período exibido até ago/2026). Uma reclamação em destaque fala de rastreamento ineficiente e demora para resolver indisponibilidade — [Reclame Aqui Cobli](https://www.reclameaqui.com.br/empresa/cobli/)

**Sascar → MICHELIN Connected Fleet**
- A Sascar foi rebatizada como "MICHELIN Connected Fleet powered by Sascar". Tem rede técnica nacional 24/7 para instalação e manutenção dos rastreadores — [Michelin Connected Fleet BR](https://go.connectedfleet.michelin.com/pt-br/lp-brasil/pfro-gestao-de-frota)
- As plataformas incluem checklists e lembretes de manutenção. A empresa promove TPMS (pressão e temperatura de pneus) como base da manutenção baseada em dados — [Michelin blog – preventiva](https://connectedfleet.michelin.com/pt-br/blog/como-a-manutencao-preventiva-pode-controlar-melhor-os-gastos-de-sua-frota); [Michelin blog – TPMS](https://connectedfleet.michelin.com/pt-br/blog/o-papel-do-tpms-na-eficiência-operacional-e-segurança-da-frota)
- Na CES 2026, a Michelin e a Sonatus apresentaram modelos preditivos de desgaste de pneu (SmartWear/SmartLoad) a partir de sinais dinâmicos do veículo — [Mecânica Online](https://mecanicaonline.com.br/2025/12/michelin-e-sonatus-apresentam-sistema-preditivo-de-pneus-na-ces-2026/)

**Omnilink**
- Lançou o "Pneu Conectado", com sensor nativo de localização, temperatura e pressão integrado ao pneu — [SETCESP](https://setcesp.org.br/inovacao/omnilink-introduz-pneu-inteligente-com-sensor-nativo-e-amplia-gestao-sobre-um-dos-ativos-mais-criticos-da-frota/)

**Infleet**
- É uma startup baiana. Tem um "superapp" que integra monitoramento, telemetria, videotelemetria, manutenção, checklist e abastecimento — [Startups.com.br](https://startups.com.br/logistica/logtech-infleet-espera-crescer-150-com-seed-de-r-5-mi/); [Startupi](https://startupi.com.br/?p=242407)
- Captou R$ 40 milhões em Série A (set/2025) para expandir uma plataforma IoT com IA — [Teletime](https://teletime.com.br/19/09/2025/infleet-capta-r-40-milhoes-para-ampliar-plataforma-iot-de-gestao-de-frotas/)
- Preço sob consulta — [B2B Stack](https://www.b2bstack.com.br/compare/infleet-o-copiloto-inteligente-para-gestao-de-frotas-vs-coprev-controle-inteligente-de-manutencao)

**Gestran (SaaS de gestão de frota com manutenção profunda)**
- O Gestran Frota tem módulos de abastecimento, pneus, manutenção, checklist, despesas, financeiro e almoxarifado. A empresa lançou o Gestran Data e o OpenFleet, este último para integração com telemetria e dados de pneus, e mirava crescer 40% em 2024 — [Mundo Logística](https://mundologistica.com.br/noticias/gestran-lanca-solucoes-tecnologicas-para-gestao-de-frotas)
- Lançou também um agente de IA para gestão de frotas — [Mundo Logística](https://mundologistica.com.br/noticias/gestran-agente-ia-gestao-frotas)
- O software de pneus mantém histórico por pneu e gera alertas de rodízio, calibragem, envio para reforma e desgaste. A empresa cita até 25% de economia, número do próprio fornecedor. Tem parceria com a Marangoni (Control Tyre) — [Mundo Logística – pneus](https://mundologistica.com.br/noticias/software-de-gestao-de-pneus-pode-gerar-economia-de-ate-25)

**Prolog (gestão de pneus)**
- O app se descreve como "maior sistema de gestão de pneu do Brasil" (alegação de marketing). Faz coleta de km no checklist, aferição e movimentação de pneus, com aferidor conectado via Bluetooth — [App Store Prolog](https://apps.apple.com/bg/app/prolog/id6470036075)

**Ticket Log / Edenred**
- Oferece gestão da manutenção com duas opções: autogestão pelo cliente ou terceirização da operação. Em 2019 tinha mais de 300 mil veículos geridos por essa solução — [Edenred release 2019](https://www.edenred.com.br/wp-content/uploads/2020/01/20190301_Ticket-Log_release_WOW_Manutenção-REV.docx.pdf)
- A Rede Log reunia mais de 10 mil oficinas credenciadas. O app "Ticket Log – Rede de Manutenção" permite que a oficina consulte OS e acompanhe autorizações (orçamento → aprovação) — [SETCESP](https://setcesp.org.br/noticias/ticket-log-lanca-aplicativo-dedicado-para-oficinas-credenciadas/)
- Todos esses dados são de 2018 a 2021 e precisam ser confirmados para 2026.

**Montadoras (portais OEM)**
- Volvo: a Connect API é um serviço por assinatura, sem hardware adicional, que inclui ACEA rFMS e um conjunto de dados estendido (ficha técnica em pt-BR). Via rFMS também é possível obter dados de caminhões de outras marcas — [Volvo Buses – Connect API PT-BR](https://volvobuses.com/content/dam/volvo-buses/markets/master/data-access/connect-api/Fact%20sheet%20Data%20Access%20Services_Connect%20API_PT-BR_VolvoBuses_2024.pdf); [Volvo FMS Gateway PT-BR](https://volvobuses.com/content/dam/volvo-buses/markets/master/data-access/connect-api/Fact-sheet-Data-Access-Services_FMS%20Gateway_PT-BR_VolvoBuses_2024.pdf)
- Scania: entrega dados no padrão rFMS via API servidor a servidor, com portal developer.scania.com — [Scania Data Access](https://scania.com/uk/en/home/services/connected-services/data-access.html); [Scania leaflet](https://scania.com/content/dam/www/market/master/pdf/services/digital-services/Scania_Data_access_Leaflet.pdf)
- Mercedes-Benz do Brasil: lançou em 2025 o FleetBus, telemetria e conectividade para frotas de ônibus — [Truck & Bus Builder](https://truckandbusbuilder.com/article/2025/07/23/mercedes-benz-do-brasil-offers-latest-in-fleet-management-technology)

**Outros**
- Vivo Frota Inteligente: telemática de operadora, com planos por veículo. A instalação não está incluída no valor — [Vivo – planos](https://vivo.com.br/para-empresas/produtos-e-servicos/digitais/iot/frota-inteligente/compare-os-planos)
- Também aparecem no mercado: 3SAT, Sofit View, Frota Certa (preço sob consulta), Mobi7, COPREV (controle de manutenção) e Bsoft — [B2B Stack – categoria](https://www.b2bstack.com.br/categoria/softwares-de-gestao-de-frota); [Bsoft – pneus](https://bsoft.com.br/blog/gestao-de-pneus-reducao-de-custos)

### Inferences
- Telemática brasileira (Sascar, Omnilink, Cobli, Vivo): o foco é rastreamento, segurança e roubo de carga, telemetria e TPMS. A manutenção é secundária, com lembretes, checklist e histórico. O Tenaz deve integrá-las como fontes de km e horímetro, DTC e pressão de pneu, sem competir com elas.
- Gestran e Prolog/Bsoft são concorrentes diretos no módulo de frota, especialmente em pneus e almoxarifado. Ticket Log é concorrente e ao mesmo tempo parceiro potencial (rede credenciada, transações de manutenção e de combustível).
- O padrão rFMS (Volvo, Scania e, via rFMS, outras marcas) é a forma de integrar km, horas e diagnóstico de frotas sem rastreador próprio.

### Gaps
- Não houve dados nesta rodada sobre Maxtrack, Positron/PST, Onixsat, Veloe/Sem Parar (tags de pedágio e gestão de frota), Frota Fácil, nem sobre detalhes dos portais Mercedes-Benz (Fleetboard) e Scania no Brasil.
- Não foram encontrados preços oficiais de Sascar, Omnilink, Onixsat, Infleet ou Gestran.
- Não está confirmado se Ticket Log ainda opera a Rede Log com mais de 10 mil oficinas em 2026.
- Não foram consultadas fontes da imprensa especializada (O Carreteiro, Transporte Moderno, Frota&Cia, Tecnologística).

---

## 3. Escopo funcional específico de frota (cadastro, medições, preventiva, pneus, combustível, checklist, sinistros, multas, documentos, oficinas, garantia, TCO, renovação) e requisitos regulatórios brasileiros

### Takeaway
O módulo de frota do Tenaz precisa de três coisas além de um CMMS genérico:
- Medidores múltiplos (km e horas) com origem rastreável (manual, telemática, cartão combustível).
- Pneu como ativo próprio, com vida, recapagem, sulco e CPK.
- Uma agenda de conformidade regulatória brasileira com vencimentos: CRLV/IPVA, cronotacógrafo a cada 2 anos, CIV/CIPP com prazos variáveis, RNTRC e seguros obrigatórios ANTT, e multas com indicação de condutor em 30 dias via Senatran/SNE.

### Cited Findings
**Cronotacógrafo**
- A verificação Inmetro é obrigatória a cada 2 anos após a verificação inicial (Portaria Inmetro nº 481/2021, subitem 6.3.1) — [Inmetro FAQ – periodicidade](https://www.gov.br/inmetro/pt-br/acesso-a-informacao/perguntas-frequentes/metrologia-legal/cronotacografo/apos-a-verificacao-inicial-qual-e-a-periodicidade-de-verificacao-dos-cronotacografos)
- Veículo 0 km também exige verificação imediata — [Inmetro FAQ – 0 km](https://www.gov.br/inmetro/pt-br/acesso-a-informacao/perguntas-frequentes/metrologia-legal/cronotacografo/veiculo-zero-km-deve-ter-o-cronotacografo-verificado)
- O processo é ensaio em Posto Autorizado de Cronotacógrafo (PAC) mais pagamento de GRU (R$ 90,09 segundo a página do Inmetro). O certificado provisório vale 30 dias. A base legal é o CTB, art. 105, e as resoluções CONTRAN (Res. 938/2022) — [Gov.br – verificação subsequente](https://www.gov.br/pt-br/servicos/solicitar-verificacao-subsequente-de-cronotacografo); [CONTRAN 938/2022](https://www.legisweb.com.br/legislacao/?id=429779)
- A fiscalização é ativa: mais de 300 tacógrafos irregulares em uma operação no ES, e 22 autuados em uma rodovia do MS — [Folha Vitória](https://www.folhavitoria.com.br/geral/mais-de-300-tacografos-apresentam-irregularidades-durante-operacao-no-espirito-santo/); [SEMADESC MS](https://www.semadesc.ms.gov.br/aem-ms-realiza-fiscalizacao-de-cronotacografos-na-rodovia-ms-134-e-autua-22-equipamentos)

**Produtos perigosos (CIV, CIPP, CTPP)**
- O CIV certifica a inspeção do veículo e o CIPP o equipamento (tanque). A validade varia conforme o ano de fabricação e o produto, de 4 a 36 meses (fonte de imprensa, out/2026) — [Portal do Trânsito](https://www.portaldotransito.com.br/noticias/fiscalizacao-e-legislacao/documentacao-do-veiculo/cargas-perigosas-flexibilizacao-na-exigencia-de-documentos-pode-causar-acidentes-nas-estradas/?noamp=mobile)
- O CTPP é emitido para equipamentos novos, pelo Inmetro junto com o fabricante, com a mesma validade legal do CIPP — [Portal do Trânsito](https://www.portaldotransito.com.br/noticias/fiscalizacao-e-legislacao/documentacao-do-veiculo/cargas-perigosas-flexibilizacao-na-exigencia-de-documentos-pode-causar-acidentes-nas-estradas/?noamp=mobile)
- A inspeção é feita por organismos acreditados pelo Inmetro (OIA-PP). Os certificados são enviados à ANTT por integração, conforme manual — [Resolução ANTT 5998/2022](https://depi.unicamp.br/arquivo/uploads/resolucao-5998-3-de-nov-2022); [Manual ANTT–Inmetro](https://www.gov.br/antt/pt-br/assuntos/cargas/arquivos/arquivos-produtos-perigosos/manual-de-integracao-antt-e-inmetro-para-envio-dos-certificados.pdf); [Inmetro – OIA-PP](https://www.gov.br/inmetro/pt-br/acesso-a-informacao/perguntas-frequentes/acreditacao/organismos-de-inspecao-de-produtos-perigosos-oia-pp)
- Portar CIV vencido ou não ter o CIV original no veículo são infrações autuáveis — [Defesa Civil PR – checklist](https://www.defesacivil.pr.gov.br/sites/defesa-civil/arquivos_restritos/files/documento/2021-12/checklist_codigos_autucao.pdf)

**RNTRC e ANTT**
- A Resolução ANTT nº 6.068/2025 torna a comprovação de 3 seguros obrigatórios condição para inscrição ou renovação no RNTRC, com cruzamento de dados com a Susep e possibilidade de suspensão do cadastro. As categorias são TAC, ETC e CTC — [Mundo Logística](https://mundologistica.com.br/noticias/antt-reforca-obrigatoriedade-de-seguros-no-trc); [Mundo Logística – fiscalização eletrônica](https://mundologistica.com.br/noticias/antt-fiscalizacao-eletronica-seguros-obrigatorios)
- Operar sem RNTRC ativo gera multa de R$ 3 mil, segundo fonte setorial — [Mundo Logística](https://mundologistica.com.br/noticias/regularizacao-rntrc-cresce-apos-revalidacao-da-antt)
- O serviço oficial de cadastro tem o RNTRC Digital — [Gov.br – cadastro RNTRC](https://www.gov.br/pt-br/servicos/solicitar-cadastro-de-transportador-no-rntrc)

**Multas e indicação de condutor**
- A empresa com frota (PJ) indica o real infrator pelo Portal de Serviços da Senatran com certificado digital. O condutor confirma no app CDT. O prazo é de 30 dias a partir da notificação — [Serpro – CDT condutor infrator](https://serpro.gov.br/menu/noticias/noticias-2021/cdt-condutor-infrator); [Comércio SP](https://comerciosaopaulo.com.br/?p=12505)
- O SNE (Serpro/Senatran) faz a notificação eletrônica e dá 40% de desconto para quem adere. A Senatran e o Serpro criaram o "e-Frotas", um conjunto de serviços para PJ gerirem sua frota — [Serpro – frotistas](https://serpro.gov.br/menu/noticias/noticias-2023/tecnologia-serpro-facilita-vida-de-frotistas); [Serpro – Senatran](https://serpro.gov.br/clientes/senatran)
- Se a PJ não indica o condutor, há multa NIC (não identificação do condutor) — [Comércio SP](https://comerciosaopaulo.com.br/?p=12505)

**Pneus**
- Sem rodízio, o desgaste é desigual. A recapagem exige registro de calibragem, desgaste e histórico de rodízio por veículo — [Jornal do Brás](https://jornaldobras.com.br/noticia/119497/manutencao-de-pneus-exige-controle-tecnico-na-frota)
- A má gestão pode aumentar em até 20% o gasto com pneus — [Diário do Comércio](https://diariodocomercio.com.br/negocios/entenda-ma-gestao-pode-aumentar-gastos-pneus/)
- Recursos de mercado: histórico por pneu, alertas de rodízio, calibragem, reforma e desgaste ([Gestran](https://mundologistica.com.br/noticias/software-de-gestao-de-pneus-pode-gerar-economia-de-ate-25)), e aferição via Bluetooth ([Prolog](https://apps.apple.com/bg/app/prolog/id6470036075)). No Premium do Fleetio: inspeção de pneu e garantia de pneu ([Fleetio Pricing](https://www.fleetio.com/pricing)). TPMS e pneu com sensor nativo: [Michelin](https://connectedfleet.michelin.com/pt-br/blog/o-papel-do-tpms-na-eficiência-operacional-e-segurança-da-frota); [Omnilink](https://setcesp.org.br/inovacao/omnilink-introduz-pneu-inteligente-com-sensor-nativo-e-amplia-gestao-sobre-um-dos-ativos-mais-criticos-da-frota/)

**Ciclo de vida, TCO, sinistros e checklist**
- Benchmark do FleetWave: aquisição → uso → manutenção → conformidade → acidentes → descarte e remarketing, com relatórios de TCO. O app do motorista registra defeitos, acidentes e recibos de combustível — [Chevin](https://www.chevinfleet.com/products/fleetwave/)
- Inspeções móveis offline com templates configuráveis estão até no plano de entrada do Fleetio — [Fleetio Pricing](https://www.fleetio.com/pricing)

### Inferences (lista consolidada de funcionalidades recomendada para o módulo de frota do Tenaz)
**A. Cadastro do veículo e implemento**
- Placa no padrão Mercosul e no padrão antigo (validação de formato), RENAVAM, chassi/VIN, marca, modelo, ano de fabricação e modelo, combustível e categoria (leve, pesado, implemento, máquina).
- Composição cavalo + carreta (vínculo dinâmico de implementos) e centro de custo e filial.
- Proprietário, que pode ser próprio, locado ou agregado, para separar TAC, ETC e frota de terceiros.
- A estrutura deve compartilhar o modelo de "ativo" com equipamentos industriais, ou seja, um único CMMS.

**B. Medidores**
- Múltiplos medidores por ativo: km, horímetro do motor, horímetro da PTO e da refrigeração.
- Cada leitura guarda a origem (manual, app, telemática, cartão combustível ou OEM rFMS).
- Validação de leituras: rejeitar regressão e saltos implausíveis, e tratar troca de painel ou hodômetro.

**C. Planos de preventiva**
- Gatilho por km, horas, tempo ou o que vencer primeiro, com antecipação configurável.
- Planos por modelo ou família e kits de peças por plano.
- Previsão de próxima data pela média de rodagem.

**D. OS e oficina**
- Oficina interna, com mão de obra apontada, boxes e técnicos.
- Oficina externa: orçamento → aprovação → execução → nota, no estilo Ticket Log.
- Solicitação de serviço originada de checklist, de DTC ou do motorista.
- Garantia de peças e serviços (por km ou prazo) e recall.

**E. Pneus**
- Pneu como ativo serializado com número de fogo e DOT.
- Vida (novo, R1, R2...), recapadora, medições de sulco em mm por posição, pressão, rodízio e mapa de posições por eixo.
- Sucateamento com motivo e garantia.
- CPK (custo por km): R$ do pneu e das reformas ÷ km rodados na vida, calculado por marca, modelo e banda.
- Integração com aferidor Bluetooth e TPMS.

**F. Combustível**
- Abastecimentos manuais e importados do cartão (Ticket Log etc.), consumo km/l e l/h, metas por modelo.
- Alertas de fraude e anomalia: volume acima da capacidade do tanque, km incoerente, abastecimento fora da rota ou do horário, frequência anômala.
- Arla 32 e combustível de tanque interno (comboio e posto próprio).

**G. Checklist do motorista**
- Checklist pré-viagem mobile offline com foto e assinatura.
- Item reprovado gera solicitação de manutenção, e item crítico bloqueia o veículo.
- Coleta de km no próprio checklist (como no Prolog).

**H. Sinistros**
- Registro com fotos, B.O., terceiros, seguradora, franquia e OS de reparo vinculada.

**I. Multas**
- Importação e cadastro da infração (AIT, órgão, enquadramento, valor, pontos).
- Identificação automática do condutor por escala, checklist ou login no app.
- Controle do prazo de 30 dias para indicação via Senatran, do risco de NIC, do desconto SNE e do desconto em folha.

**J. Documentos e vencimentos**
- Por veículo: CRLV-e/licenciamento, IPVA, seguro, cronotacógrafo (2 anos), CIV/CIPP/CTPP (4 a 36 meses), AET, RNTRC e seguros ANTT, laudos (opacidade, inspeção ambiental).
- Por motorista: CNH, MOPP, toxicológico.
- Calendário de alertas, anexos e bloqueio de despacho com documento vencido.

**K. TCO e renovação**
- Custo por km e por hora por ativo, composto por manutenção, pneus, combustível, multas, sinistros, seguro, impostos e depreciação.
- Curva de custo acumulado para apoiar a decisão de renovar ou vender, e registro de baixa e venda.

### Gaps
- Não foi confirmada a tabela oficial de validade de CIV/CIPP por produto e idade (a fonte é de imprensa), nem a regra vigente de AET (DNIT e órgãos estaduais).
- Não foi encontrada documentação técnica oficial de uma API Serpro/Senatran (e-Frotas) para que um software de terceiros consulte multas ou indique condutores, nem os requisitos de credenciamento.
- Não houve fonte nesta rodada para as regras de placa Mercosul, RENAVAM, CRLV-e digital, IPVA por UF, MOPP ou toxicológico. São requisitos conhecidos, mas sem citação aqui.
- Não há benchmark brasileiro de metas de CPK nem de consumo km/l.

---

## 4. Padrões de integração (APIs de telemática, J1939/DTC, cartões combustível, tags de pedágio, portais OEM)

### Takeaway
O padrão consolidado é: a telemática ou OEM envia odômetro, horas, DTC e localização. O CMMS recebe esses dados por API ou webhook e converte alertas em solicitações ou OS. O cartão combustível envia as transações (litros, valor, km informado, posto). No Brasil, as APIs relevantes são Cobli, Geotab e o rFMS das montadoras Volvo e Scania. Ainda faltam dados de Sascar, Omnilink e tags de pedágio.

### Cited Findings
- O fluxo Geotab → Fleetio: alerta do dispositivo, sincronização com o Fleetio e disparo de fluxo de manutenção interna ou externa. O Fleetio também integra com cartões combustível e oficinas — [Fleetio – integração Geotab (PDF)](https://marketing-cdn.fleetio.com/documents/app-directory/geotab-integration-overview.pdf); [Geotab Marketplace – Fleetio](https://marketplace.geotab.com/solutions/fleetio)
- O enriquecimento de DTC da Geotab/Noregon padroniza códigos entre montadoras e dá severidade e ação recomendada — [Geotab](https://www.geotab.com/press-release/maintenance-announcements-connect-2025/)
- Segundo um guia de integração, OS, peças e custos completos normalmente ficam em um CMMS construído sobre o feed de dados, e a escolha é de "propriedade, não capacidade" — [OxMaint](https://oxmaint.com/industries/fleet-management/telematics-cmms-integration-for-fleets-guide); [OxMaint – J1939](https://oxmaint.com/industries/fleet-management/j1939-engine-data-real-time-diagnostics-fleet-guide-2026)
- Limites de API e webhook do Fleetio como referência de mercado: 50 a 250 req/min e 10 mil eventos por 30 dias, ou ilimitados — [Fleetio Pricing](https://www.fleetio.com/pricing)
- Volvo: Connect API e FMS Gateway, com ACEA rFMS mais dados estendidos, sem hardware adicional — [Volvo Buses PT-BR](https://volvobuses.com/content/dam/volvo-buses/markets/master/data-access/connect-api/Fact%20sheet%20Data%20Access%20Services_Connect%20API_PT-BR_VolvoBuses_2024.pdf). Scania: rFMS via API servidor a servidor — [Scania](https://scania.com/uk/en/home/services/connected-services/data-access.html)
- Cobli: API com endpoints de avisos e registros de manutenção e de peças — [docs.cobli.co](https://docs.cobli.co/reference/createpart)
- Gestran OpenFleet: integração com telemetria e dados de pneus — [Mundo Logística](https://mundologistica.com.br/noticias/gestran-lanca-solucoes-tecnologicas-para-gestao-de-frotas)
- Agregadores de telemática OEM multimarca, como a Navixy, existem como alternativa a integrar montadora por montadora — [Navixy OEM telematics](https://navixy.com/learn-more/oem-telematics/brands)

### Inferences
- O Tenaz deve ter uma camada de conectores com modelo canônico. Os eventos seriam: leitura de medidor (ativo, tipo, valor, data/hora, origem), DTC (SPN/FMI no J1939 ou P-code no OBD-II, com severidade), transação de combustível, evento de pneu (pressão e temperatura) e localização.
- A ordem de prioridade dos conectores seria: Cobli e Geotab (têm API documentada), rFMS (Volvo e Scania), Ticket Log (combustível e manutenção), e depois Sascar e Omnilink.
- É recomendável oferecer API pública e webhooks desde o plano intermediário, porque o Fleetio cobra por isso e é um diferencial de venda.
- A importação por CSV ou planilha de abastecimentos, multas e leituras deve existir como fallback, já que muitos fornecedores brasileiros não têm API pública.

### Gaps
- Não foi verificada a existência nem o formato de API pública de Sascar/Michelin, Omnilink, Maxtrack, Positron, Onixsat, Ticket Log, Veloe ou Sem Parar.
- Não houve dados sobre integração de tags de pedágio (Veloe, Sem Parar, ConectCar) com software de frota.
- A documentação da API de manutenção da Samsara e o overview da API do Fleetio não foram lidos (404 e conteúdo vazio).
- Não há dados sobre a disponibilidade de PGN/SPN J1939 específicos nas telemáticas brasileiras.

---

## 5. Preço por veículo (Brasil e global)

### Takeaway
Globalmente, software de manutenção puro custa de US$ 4 a 10 por veículo/mês (Fleetio, Whip Around, FleetWave de £3 a 5). A telemática com hardware custa de US$ 20 a 50 por veículo/mês. No Brasil quase todos trabalham com preço sob consulta. A única referência pública encontrada é de cerca de R$ 50 a 76 por veículo/mês para telemática leve (Cobli), mas com dados conflitantes.

### Cited Findings
- Fleetio: US$ 4/5, 7 e 10 por veículo/mês, com usuários ilimitados — [Fleetio Pricing](https://www.fleetio.com/pricing)
- Whip Around: US$ 5 a 10 por ativo/mês, ou flat FixedUnlimited — [G2](https://www.g2.com/products/whip-around-whip-around/pricing); [Whip Around](https://whiparound.com/fixedunlimited-guide/)
- Chevin FleetWave: Lite £3 e Core £5 por veículo/mês, Advanced sob consulta (agregador) — [Chevin pricing](https://www.chevinfleet.com/pricing)
- Samsara: cerca de US$ 27 a 33 por veículo/mês com contrato de 3 anos (terceiros) — [Airpinpoint](https://airpinpoint.com/compare/samsara-competitors)
- Motive: cerca de US$ 25 a 50 por veículo/mês com contrato de 12 meses (terceiros) — [TechRadar](https://www.techradar.com/nz/reviews/keeptruckin); [SelectHub](https://www.selecthub.com/fleet-management-software/samsara-vs-gomotive/)
- Geotab: US$ 15 a 40 por veículo/mês mais hardware de US$ 80 a 150, com add-ons cobrados à parte (terceiros) — [CostBench](https://costbench.com/software/fleet-management/geotab/)
- Cobli: "a partir de R$ 50 por veículo/mês" ([Produttivo](https://produttivo.com.br/blog/sistema-de-gestao-de-frotas)) contra "R$ 76,00 mensal" e "R$ 1.990 mensal" no B2B Stack, com unidade não clara ([B2B Stack](https://www.b2bstack.com.br/product/cobli)). As fontes são conflitantes.
- Vivo Frota Inteligente: preço por veículo, com instalação à parte (o valor não foi capturado) — [Vivo](https://vivo.com.br/para-empresas/produtos-e-servicos/digitais/iot/frota-inteligente/compare-os-planos)
- Infleet e Frota Certa: sob consulta — [B2B Stack](https://www.b2bstack.com.br/compare/infleet-o-copiloto-inteligente-para-gestao-de-frotas-vs-coprev-controle-inteligente-de-manutencao); [B2B Stack – Frota Certa](https://www.b2bstack.com.br/compare/frota-certa-vs-cobli)

### Inferences
- Com câmbio de cerca de R$ 5 a 5,5/US$ (premissa não verificada nesta pesquisa), o Fleetio custaria aproximadamente R$ 20 a 55 por veículo/mês. Para o Tenaz, uma faixa entre R$ 15 e 45 por veículo/mês para software sem hardware parece competitiva. Isso é hipótese a validar com cotações de Gestran e Infleet.
- Preço público, usuários ilimitados e cobrança por ativo, como no Fleetio, são diferenciais em um mercado brasileiro dominado por "sob consulta". Um modelo flat como o FixedUnlimited pode atrair empresas com muitos equipamentos além de veículos.

### Gaps
- Não foram obtidos preços oficiais em 2026 de Sascar, Omnilink, Onixsat, Maxtrack, Positron, Gestran, Infleet, Ticket Log (taxa de administração da manutenção e do combustível) nem Veloe.
- Não há preço oficial de lista de Samsara, Motive ou Geotab (todos os números são estimativas de terceiros).
