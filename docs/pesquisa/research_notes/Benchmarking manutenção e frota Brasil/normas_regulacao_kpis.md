# Normas, regulação brasileira e KPIs para um CMMS de equipamentos + frota (Brasil, nov/2026)

> Convenção de confiabilidade usada nestas notas:
> - **[verificado]**: lido em fonte nesta pesquisa (texto da norma, FAQ oficial ou compilação fiel do texto legal).
> - **[fonte primária indicada, não re-lida]**: fato de conhecimento consolidado, com a URL oficial onde está o texto; o redator deve conferir antes de virar regra contratual.
> - Benchmarks de fornecedores (Tractian, Oxmaint etc.) são "regra de bolso", não valores SMRP oficiais.

---

## 1. Normas técnicas e terminologia (NBR 5462, ISO 14224, ISO 55000, EN 13306, EN 15341, SMRP, RCM, FMEA): o que o CMMS deve implementar

### Takeaway
A base terminológica brasileira continua sendo a ABNT NBR 5462:1994 (vigente, não substituída), enquanto a gestão de ativos foi revisada em 2024 (ISO 55000/55001:2024, adotadas pela ABNT em 09/2024, transição de certificação até 31/07/2027). O CMMS deve modelar a hierarquia de ativos e a codificação de falhas segundo a ISO 14224 (modo/mecanismo/causa/método de detecção/atividade), adotar os conceitos de tempo da NBR 5462/EN 13306 para calcular disponibilidade/MTBF/MTTR e expor KPIs nomeados conforme EN 15341:2019+A1:2022 e SMRP.

### Cited Findings

**ABNT NBR 5462 (Confiabilidade e mantenabilidade — terminologia)**
- Edição de 11/1994 aparece como **"Vigente"** em catálogo comercial de normas; substituiu a NBR 5462/1981; escopo: "define os termos relacionados com a confiabilidade e a mantenabilidade" — [Normas.com.br](https://normas.com.br/visualizar/abnt-nbr-nm/8044/abnt-nbr5462-confiabilidade-e-mantenabilidade); [Target](https://www.target.com.br/produtos/normas-tecnicas/29220/nbr5462-confiabilidade-e-mantenabilidade); [cópia do texto (UFSB)](https://ufsb.edu.br/propa/images/dinfra/coman/Legislações/NBR-5462.pdf)
- Não encontrei evidência de cancelamento/substituição pela IEC 60050-192 (busca específica sem resultado) — [busca: catálogos ABNT/Target]( https://www.target.com.br/produtos/normas-tecnicas/29220/nbr5462-confiabilidade-e-mantenabilidade)
- A NBR 5462 continua citada como referência normativa por normas mais novas (ex.: ABNT NBR ISO 14224 e ABNT NBR ISO 8927 — máquinas rodoviárias, disponibilidade de máquina, vocabulário) — [Normas.com.br ISO 14224](https://www.normas.com.br/visualizar/abnt-nbr-nm/31226/abnt-nbriso14224-industrias-de-petroleo-coleta-e-intercambio-de-dados-de-confiabilidade-e-manutencao-para-equipamentos); [Target NBR ISO 8927](https://www.target.com.br/produtos/normas-tecnicas/39118/nbriso8927-maquinas-rodoviarias-disponibilidade-de-maquina-vocabulario)
- Norma correlata: ABNT NBR 13533 — coleta de dados de campo relativos a confiabilidade, mantenabilidade, disponibilidade e suporte à manutenção — [Target](https://www.target.com.br/produtos/normas-tecnicas/28414/nbr13533-coleta-de-dados-de-campo-relativos-a-confiabilidade-mantenabilidade-disponibilidade-e-suporte-a-manutencao)
- Conceitos usados no mercado a partir da NBR 5462: manutenção corretiva (após a pane), preventiva (intervalos predeterminados ou critérios prescritos) e preditiva/"controlada" (baseada em condição, com técnicas de análise e supervisão) — [FM2S, resumo NBR 5462](https://www.fm2s.com.br/blog/nbr-5462)
- Para edificações: ABNT NBR 5674 (manutenção de edificações) foi revisada em 2024 ("NBR 5674 (ABNT, 2024b)") — relevante para clientes de varejo/serviços com manutenção predial — [Monografia Poli-USP 2025/2026](https://poli-integra.poli.usp.br/wp-content/uploads/2026/02/2025_Igor-Gama-Galante.pdf)

**ISO 55000 / 55001 / 55002 (gestão de ativos)**
- ISO 55001:2024 é a 2ª edição (julho/2024) — requisitos do sistema de gestão de ativos — [ISO/TC 251](https://committee.iso.org/sites/tc251/home/projects/published/iso-55001.html); [preview ANSI](https://webstore.ansi.org/preview-pages/ISO/preview_ISO+55001-2024.pdf)
- ISO 55000:2024 traz princípios e definições usados na 55001:2024; comparativo 2014 x 2024 — [Assetivity](https://www.assetivity.com.au/articles/asset-management/iso-550012024-vs-iso-550012014-a-comparison/)
- A revisão amplia escopo para todos os tipos de ativos/organizações e enfatiza sustentabilidade, dados e conhecimento — [artigo 2025 (CSPress)](https://science.cspress.com.cn/uploadfile/202509/628f219cd77e5c8.pdf)
- ABNT NBR ISO 55000 e NBR ISO 55001 publicadas em **09/2024**, situação vigente; existe ABNT NBR ISO 55002 (diretrizes de aplicação) — [Normas.com.br CEE-251](https://www.normas.com.br/produto/normas-brasileiras-e-mercosul/pesquisar/cee-251/gestao-de-ativos/palavra/6820/1); [NBR ISO 55002](https://www.normas.com.br/visualizar/abnt-nbr-nm/33982/abnt-nbriso55002-gestao-de-ativos-sistemas-de-gestao-diretrizes-para-a-aplicacao-da-abnt-nbr-iso-55001)
- Inmetro/Cgcre (Ofício Circular 19/2026): transição de 3 anos, até **31/07/2027**, para certificadoras migrarem à ISO 55001:2024 — [Inmetro](https://www.gov.br/inmetro/pt-br/assuntos/acreditacao-reconhecimento-bpl/comunicados/dicor-alerta/2026/oficio-circular-19-orientacoes-para-transicao-da-norma-iso-55001-2024.pdf)

**ISO 14224 (coleta e intercâmbio de dados de confiabilidade e manutenção)**
- Existe adoção brasileira ABNT NBR ISO 14224 ("Indústrias de petróleo... coleta e intercâmbio de dados de confiabilidade e manutenção para equipamentos") — [Normas.com.br](https://www.normas.com.br/visualizar/abnt-nbr-nm/31226/abnt-nbriso14224-industrias-de-petroleo-coleta-e-intercambio-de-dados-de-confiabilidade-e-manutencao-para-equipamentos)
- Edição internacional vigente: ISO 14224:2016 (3ª ed.). Estrutura (conhecimento consolidado): taxonomia de 9 níveis — 1 indústria, 2 categoria de negócio, 3 instalação, 4 planta/unidade, 5 seção/sistema, 6 unidade de equipamento (classe/tipo), 7 subunidade, 8 componente/item manutenível, 9 peça; três grupos de dados: equipamento, falha, manutenção; códigos de **modo de falha** por classe de equipamento (3 letras, ex.: FTS falha na partida, STP falha na parada, ELP vazamento externo de fluido de processo, VIB vibração, NOI ruído, OHE superaquecimento, BRD quebra, SER pequenos problemas em serviço); **mecanismo de falha** em grupos numerados (1 mecânico — 1.1 vazamento, 1.2 vibração, 1.3 folga/desalinhamento, 1.4 deformação, 1.5 afrouxamento, 1.6 travamento; 2 material — 2.1 cavitação, 2.2 corrosão, 2.3 erosão, 2.4 desgaste, 2.5 ruptura, 2.6 fadiga, 2.7 superaquecimento, 2.8 explosão; 3 instrumentação; 4 elétrico; 5 influência externa; 6 diversos); **causa de falha** (1 projeto, 2 fabricação/instalação, 3 operação/manutenção, 4 gestão, 5 diversos); **método de detecção** (manutenção periódica, teste funcional, inspeção, monitoramento de condição, observação casual, produção/demanda etc.); **atividade de manutenção** (substituir, reparar, modificar, ajustar, reformar, verificar, servir, testar, inspecionar, revisão geral, combinação, outro) — [ISO 14224:2016 (catálogo ISO)](https://www.iso.org/standard/64076.html) [fonte primária indicada, não re-lida]

**EN 13306 / EN 15341 (CEN/TC 319)**
- EN 15341 (KPIs da função manutenção): edição vigente **EN 15341:2019+A1:2022**; a edição 2019 "pura" consta como retirada em 15/06/2022 em listagem nacional; preparada pelo CEN/TC 319 — [Genorma](https://genorma.com/en/standards/en-15341-2019-a1-2022); [iTeh EN 15341:2019](https://iteh.es/catalog/standards/cen/7d046ef9-0294-41a0-8cd5-d2c6c02ff742/en-15341-2019)
- Estrutura da EN 15341:2019: KPIs organizados em **8 grupos** — 1 de gestão de ativos físicos + 6 subfunções da manutenção (o 8º não confirmado no snippet); a edição 2007 era organizada em indicadores econômicos, técnicos e organizacionais — [Reliamag KPI reference](https://reliamag.com/guides/maintenance-reliability-kpi-reference/)
- Literatura brasileira usa EN 15341:2019 e ISO 55000:2024 como base para KPIs e classifica indicadores em econômicos (ex.: custo de manutenção / valor de reposição do ativo, custo por tipo de manutenção, R$/m²), técnicos (MTBF, MTTR, disponibilidade, % cumprimento do plano preventivo) e organizacionais (backlog em horas, relação horas planejadas x não planejadas) — [Monografia Poli-USP](https://poli-integra.poli.usp.br/wp-content/uploads/2026/02/2025_Igor-Gama-Galante.pdf)
- EN 13306 (terminologia de manutenção) — edição de referência 2017 — [catálogo CEN/iTeh EN 13306](https://standards.iteh.ai/catalog/standards/cen/5af77559-ad83-4f2a-b875-e0d2e2b01e1a/en-13306-2017) [fonte primária indicada, não re-lida; URL de catálogo não verificada]

**SMRP Best Practices**
- SMRP define "Maintenance cost as % of RAV" (custo anual total de manutenção / valor de reposição do ativo); métricas de gestão do trabalho incluem wrench time, trabalho planejado e PM/PdM compliance (métrica **5.4.14**); edições 5ª/6ª em circulação — [Fiix MC/RAV](https://fiixsoftware.com/maintenance-metrics/mcrav/); [SMRP Best Practices Metrics 5ª ed. (cópia)](https://tegnum.edu.pe/wp-content/uploads/2023/09/Standard-SMRP-Best-Practices-Metrics-5ta-Edition-1.pdf); [SMRP workshop 2018 (PEMAC)](https://pemac.org/sites/default/files/2018_SMRP_BPMW_Full_Day.pdf)
- SMRP organiza métricas em 5 pilares do corpo de conhecimento (negócio e gestão; confiabilidade de processo de produção; confiabilidade de equipamentos; organização e liderança; gestão do trabalho) — [Open Exam Prep CMRP](https://open-exam-prep.com/study-guides/cmrp-smrp/strategic-planning-metrics/measure-performance-metrics)

**RCM e FMEA** [fonte primária indicada, não re-lida]
- SAE JA1011 (critérios para processo RCM, ed. 1999/rev. 2009) exige que o processo responda 7 perguntas: funções e padrões de desempenho; falhas funcionais; modos de falha; efeitos; consequências; tarefas proativas e intervalos; ações default quando não há tarefa adequada; SAE JA1012 é o guia — [SAE JA1011](https://www.sae.org/standards/content/ja1011_200908/)
- FMEA: IEC 60812:2018 (FMEA/FMECA) e Manual AIAG-VDA FMEA (2019), que substitui o RPN pela "Action Priority (AP)" baseada em Severidade/Ocorrência/Detecção — [IEC 60812:2018](https://webstore.iec.ch/en/publication/26359); [AIAG-VDA](https://www.aiag.org/quality/automotive-core-tools/fmea)

### Inferences
- **Modelo de dados mínimo (requisito):** Ativo com hierarquia configurável inspirada nos níveis 3–9 da ISO 14224 (Empresa > Unidade/Planta > Sistema/Área > Equipamento (classe/tipo) > Subconjunto > Componente/Item manutenível > Peça); para frota, o "equipamento" é o veículo (placa/RENAVAM/chassi) com subconjuntos padrão (motor, transmissão, freios, suspensão, pneus por posição, implementos, cronotacógrafo, tanque/equipamento PP).
- **Codificação de falha em toda OS corretiva:** campos obrigatórios Modo de falha (lista por classe), Mecanismo, Causa raiz (grupo ISO), Método de detecção, Atividade executada, além de "parou o ativo? (S/N)" e criticidade. Sem esses campos não há MTBF por modo nem análise de Pareto/Weibull.
- **Conceitos de tempo (NBR 5462/EN 13306):** o CMMS precisa registrar timestamps distintos: falha/abertura, reconhecimento (MTTA), início do reparo, fim do reparo, liberação para operação; e estados do ativo (operando, disponível parado, indisponível por manutenção preventiva, corretiva, aguardando peça/logística, causa externa) para suportar disponibilidade inerente x operacional.
- **ISO 55001:2024** não é exigência legal, mas clientes maiores (e auditorias até 07/2027) vão pedir: registro de ativos com valor de reposição (RAV), criticidade/risco, custo de ciclo de vida, planos vinculados a riscos e trilha de decisão — o CMMS deve armazenar RAV e custo de aquisição por ativo.
- **RCM/FMEA:** módulo opcional de "análise de modos de falha" por classe de ativo cujas tarefas geram planos preventivos/preditivos (vínculo tarefa ↔ modo de falha), com Severidade/Ocorrência/Detecção e AP.

### Gaps
- Não confirmei o status oficial da NBR 5462 no catálogo da própria ABNT (abntcatalogo.com.br), só em revendedores; não achei projeto de revisão.
- Não confirmei a lista nominal dos 8 grupos de KPIs da EN 15341:2019+A1:2022 nem o número exato de indicadores.
- Não verifiquei se há revisão em andamento da EN 13306:2017 ou da ISO 14224:2016, nem o ano da ABNT NBR ISO 14224 vigente.
- Códigos numéricos exatos das métricas SMRP (exceto 5.4.14) não foram verificados.

---

## 2. Normas Regulamentadoras (NR-1, 10, 11, 12, 13, 18, 20, 35): registros, intervalos e evidências

### Takeaway
As NRs com maior impacto em CMMS são a NR-12 (registro obrigatório de cada manutenção com 7 campos e cronograma de preventivas, inventário de máquinas) e a NR-13 (prazos máximos de inspeção por categoria, prontuário, Registro de Segurança, relatório em até 60/90 dias). A NR-1 (GRO/PGR) passou a incluir riscos psicossociais em 26/05/2026, com fiscalização orientativa até ~ago/2026. O CMMS deve gerar evidências auditáveis (quem, quando, o quê, conclusão de segurança, assinatura do responsável) e alertar vencimentos legais.

### Cited Findings

**NR-12 — Máquinas e equipamentos** (redação Portaria SEPRT 916/2019; última alteração listada Portaria MTP 4.219/2022) — [Guia Trabalhista NR-12](https://guiatrabalhista.com.br/legislacao/nr/nr12.htm); [texto MTE NR-12 atualizada 2023](https://www.gov.br/trabalho-e-emprego/pt-br/acesso-a-informacao/participacao-social/conselhos-e-orgaos-colegiados/comissao-tripartite-partitaria-permanente/normas-regulamentadora/normas-regulamentadoras-vigentes/NR12atualizada2023.pdf) [verificado]
- 12.11.1: manutenção na forma e periodicidade do fabricante, por profissional legalmente habilitado ou qualificado (ou conforme normas técnicas).
- **12.11.2: registro em livro próprio, ficha ou sistema informatizado interno**, com: (a) intervenções realizadas; (b) data; (c) serviço realizado; (d) peças reparadas/substituídas; (e) condições de segurança do equipamento; (f) **indicação conclusiva quanto às condições de segurança da máquina**; (g) nome do responsável pela execução.
- 12.11.2.1: registros disponíveis aos trabalhadores de operação/manutenção, CIPA (CIPA e de Assédio, desde MTP 4.219/2022), SESMT e Auditoria-Fiscal do Trabalho.
- 12.11.2.2: manutenções relacionadas à segurança devem ter: (a) **cronograma de execução** (preventiva); (b) descrição das técnicas de análise e meios de supervisão centralizada ou amostragem (preditiva).
- 12.11.3: intervenção com máquina parada, isolamento/dissipação de energias, **bloqueio e etiquetagem (LOTO) com horário, data, motivo e nome do responsável**, medidas para energia residual, retenção mecânica.
- 12.11.4: ENDs (ensaios não destrutivos) quando indicado pelo fabricante em partes estruturais; 12.11.5: defeito que comprometa a segurança → reparo/substituição imediata por peça original ou equivalente.
- 12.13: manual de instruções em português; se perdido, reconstituído sob responsabilidade de profissional qualificado/habilitado; ME/EPP podem ter "ficha de informação" com periodicidade de inspeção e manutenção.
- 12.14.2: inspeção rotineira pelo operador no início de cada turno; **12.14.2.1: registro dessa inspeção do operador não é obrigatório**.
- 12.16: operação/manutenção por trabalhadores capacitados e formalmente autorizados; função registrada na CTPS.
- **12.18.1 (inventário): manter relação atualizada das máquinas e equipamentos** à disposição da Auditoria-Fiscal.
- Anexo X (máquinas usadas) com prazo prorrogado até 02/01/2025 (Portaria MTE 224/2024).

**NR-13 — Caldeiras, vasos de pressão, tubulações e tanques metálicos** (redação Portaria MTP 1.846/2022, DOU 04/07/2022; alterada pela Portaria MTP 4.219/2022; retificação 20/10/2022) — [Guia Trabalhista NR-13](https://www.guiatrabalhista.com.br/legislacao/nr/nr13.htm); [retificação (IBP)](https://www.ibp.org.br/personalizado/uploads/2022/11/portaria-mtp-n-1846-retificacao-20102022.pdf) [verificado]
- Caldeiras: Categoria A (pressão de operação ≥ 1.960 kPa) e B (> 60 kPa e < 1.960 kPa) (13.4.1.1).
- Inspeção de segurança periódica de caldeiras **sem SPIE** (13.4.4.4): 12 meses (A e B); 18 meses (recuperação de álcalis); 24 meses (cat. A com teste de pressão de abertura das válvulas a cada 12 meses); 30 meses (cat. B com sistema de gerenciamento de combustão conforme Anexo IV).
- **Com SPIE** (Serviço Próprio de Inspeção de Equipamentos) (13.4.4.5): 24 meses (recuperação de álcalis); 24 meses (cat. B); 30 meses (cat. A); 48 meses (cat. A com sistema instrumentado de segurança, Anexo IV).
- Vasos de pressão — Tabela 2 (13.5.4.5), prazos máximos exame externo / interno:

| Categoria | Sem SPIE ext./int. | Com SPIE ext./int. |
|---|---|---|
| I | 1 / 3 anos | 3 / 6 anos |
| II | 2 / 4 anos | 4 / 8 anos |
| III | 3 / 6 anos | 5 / 10 anos |
| IV | 4 / 8 anos | 6 / 12 anos |
| V | 5 / 10 anos | 7 anos / critério do PLH |

  (Categoria do vaso definida pela classe de fluido e grupo de potencial de risco, Tabela 1, 13.5.1.1.) Observação: fonte secundária (Produttivo) apresenta valores "com SPIE" diferentes para III–V; usar o texto consolidado/DOU — [Produttivo](https://www.produttivo.com.br/blog/inspecao-de-vasos-de-pressao)
- Tubulações (13.6.2.2): intervalo não pode exceder o maior prazo de inspeção interna do vaso/caldeira mais crítico interligado; pode ser dobrado com justificativa do PLH até 10 anos.
- Tanques metálicos (13.7.3.2): sem prazo fixo — programa de inspeção elaborado por responsável técnico conforme códigos/normas; enquadramento de tanques (13.2.1.1 f) vigente a partir de **04/07/2026**; inspeção obrigatória para tanques instalados a partir de 20/12/2022.
- Válvulas de segurança: caldeiras — desmontagem/inspeção/teste em prazo ≤ inspeção periódica; cat. B — teste manual (alavanca) ao menos mensal com água não tratada; vasos — prazo ≤ inspeção interna (13.4.4.7/13.4.4.8/13.5.4.9).
- Documentação: **Prontuário** (código de construção, materiais, procedimentos, PMTA, teste hidrostático, desenhos, dados de dispositivos de segurança, ano, categoria); **Registro de Segurança** (ocorrências e inspeções com nome e assinatura do PLH) para caldeiras (13.4.1.8), vasos (13.5.1.7), tanques (13.7.1.3); **Relatório de inspeção emitido em até 60 dias (90 dias em parada geral)** (13.3.8); projeto de instalação e projeto de alteração/reparo por PLH; prazo para adequações de vasos (13.5.1.6.3) até 20/12/2028.
- PLH = profissional legalmente habilitado (engenheiro com atribuição) (13.3.2); certificação voluntária via organismo acreditado Inmetro (Anexo III).

**NR-1 — GRO/PGR**
- Portaria MTE 1.419/2024 incluiu fatores de riscos psicossociais no GRO/PGR; Portaria MTE 765/2025 adiou a vigência para **26/05/2026**; MTE publicou manual de interpretação do Capítulo 1.5; fiscalização com dupla visita/orientativa nos ~90 dias iniciais, autuação plena a partir do fim de agosto/2026; multas pela NR-28 — [Contábeis](https://www.contabeis.com.br/noticias/76947/nr-1-passa-a-valer-na-proxima-terca-26-mas-ainda-gera-duvidas/); [FINDES — manual MTE](https://findes.com.br/wp-content/uploads/2026/03/Informe-Estrategico-–-MTE-lanca-Manual-de-Interpretacao-e-Aplicacao-do-Capitulo-1.5-da-NR-1.pdf); [Guia Trabalhista](https://www.guiatrabalhista.com.br/tematicas/nr1-riscos-psicossociais.htm); data divergente 25/05 em [Wagner Adv.](https://wagner.adv.br/mte-prorroga-por-um-ano-a-obrigacao-das-empresas-indicarem-riscos-psicossociais/)
- PGR = inventário de riscos + plano de ação; revisão a cada 2 anos (3 anos com certificação em SST); histórico de atualizações guardado por no mínimo 20 anos; NR-1 item 1.6 permite documentos em meio digital — [NR-1 (MTE)](https://www.gov.br/trabalho-e-emprego/pt-br/acesso-a-informacao/participacao-social/conselhos-e-orgaos-colegiados/comissao-tripartite-partitaria-permanente/normas-regulamentadora/normas-regulamentadoras-vigentes/nr-01-atualizada-2025-i-3.pdf) [fonte primária indicada, não re-lida; URL exata do PDF pode ter mudado]

**NR-10 — Instalações elétricas** [fonte primária indicada, não re-lida] — [NR-10 (MTE)](https://www.gov.br/trabalho-e-emprego/pt-br/acesso-a-informacao/participacao-social/conselhos-e-orgaos-colegiados/comissao-tripartite-partitaria-permanente/normas-regulamentadora/normas-regulamentadoras-vigentes/norma-regulamentadora-no-10-nr-10)
- Esquemas unifilares atualizados (10.2.3); **Prontuário de Instalações Elétricas** obrigatório para carga instalada > 75 kW (10.2.4), contendo: procedimentos e instruções; documentação das **inspeções e medições do SPDA e aterramentos**; especificação de EPC/EPI e ferramental; documentação de qualificação/habilitação/capacitação/autorização dos trabalhadores; resultados de testes de isolação de EPC/EPI; certificações de equipamentos em áreas classificadas; **relatório técnico das inspeções atualizadas com recomendações e cronogramas de adequação**.
- Treinamento básico 40 h (Anexo III) e reciclagem bienal (10.8.8).
- SPDA (ABNT NBR 5419-3:2015): inspeção completa periódica — 1 ano para estruturas com atmosfera explosiva/serviços essenciais, 3 anos para demais — [ABNT NBR 5419-3 (catálogo)](https://www.normas.com.br/visualizar/abnt-nbr-nm/33573/abnt-nbr5419-3-protecao-contra-descargas-atmosfericas-parte-3-danos-fisicos-a-estruturas-e-perigos-a-vida) [não re-lido; confirmar prazos]

**NR-11 — Transporte e movimentação de materiais (empilhadeiras)** [fonte primária indicada, não re-lida] — [NR-11 (Guia Trabalhista)](https://www.guiatrabalhista.com.br/legislacao/nr/nr11.htm)
- Operador de equipamento motorizado com treinamento e **cartão de identificação com nome e foto, validade de 1 ano** (revalidação com exame de saúde) (11.1.6/11.1.6.1); carga máxima indicada no equipamento (11.1.3.2); sinal sonoro de advertência (11.1.8); equipamentos de transporte permanentemente inspecionados, com peças defeituosas substituídas imediatamente (11.1.3.1).

**NR-20 (inflamáveis e combustíveis), NR-35 (trabalho em altura), NR-18 (construção — PTA, gruas)** [fonte primária indicada, não re-lida]
- NR-20 exige plano de inspeção e manutenção das instalações (classes I–III), com registros, e prontuário da instalação — [NR-20 (Guia Trabalhista)](https://www.guiatrabalhista.com.br/legislacao/nr/nr20.htm)
- NR-35: inspeção de EPI/acessórios/sistemas de ancoragem antes do uso e periódica, com registro das recusas; treinamento periódico bienal (8 h) — [NR-35 (Guia Trabalhista)](https://www.guiatrabalhista.com.br/legislacao/nr/nr35.htm)
- NR-18 (Portaria SEPRT 3.733/2020): plataformas de trabalho aéreo e equipamentos de guindar exigem inspeção diária pelo operador, manutenção registrada e manual em português — [NR-18 (Guia Trabalhista)](https://www.guiatrabalhista.com.br/legislacao/nr/nr18.htm)

### Inferences
- **Regra de negócio NR-12:** toda OS concluída em ativo do tipo "máquina NR-12" deve exigir os 7 campos do 12.11.2 (incluindo um campo obrigatório "parecer conclusivo de segurança: apto / apto com restrição / inapto" e responsável), e o plano preventivo deve gerar o "cronograma de execução" exportável (PDF) por máquina. Gerar "Inventário NR-12" como relatório padrão (todas as máquinas ativas com localização, fabricante, modelo, nº série, ano, manual disponível S/N).
- **LOTO:** checklist de bloqueio com registro de data/hora/motivo/responsável por ponto de bloqueio → ganho de diferencial em manufatura.
- **NR-13:** cadastro de equipamento com atributos: tipo (caldeira/vaso/tubulação/tanque), categoria (A/B; I–V), SPIE (S/N), datas de última inspeção externa/interna e teste de válvulas; motor de cálculo de vencimento conforme as tabelas acima (parametrizável, pois a empresa pode ter prazos menores definidos pelo PLH); relatório com prazo de emissão (60/90 dias) monitorado; anexos obrigatórios: prontuário, Registro de Segurança (livro digital com assinatura do PLH), relatórios.
- **NR-10:** "Prontuário digital" como pasta estruturada por instalação com checklist dos 7 componentes e vencimentos (SPDA, testes de isolação, treinamentos).
- **Treinamentos/habilitações** (NR-10, NR-11, NR-12, NR-35, MOPP) devem ser entidades do sistema vinculadas ao técnico/motorista; OS de atividade regulada deve bloquear atribuição a quem não tem capacitação válida.
- **NR-1/PGR:** o CMMS não é o PGR, mas pode exportar evidências de controles de engenharia (manutenção de proteções, LOTO) para o inventário de riscos.

### Gaps
- Não li o texto oficial da NR-10, NR-11, NR-20, NR-35 e NR-18 nesta pesquisa; itens e números acima vêm de conhecimento consolidado e devem ser conferidos.
- Não confirmei se a NR-1 (item 1.6) exige ICP-Brasil para documentos digitais de SST ou aceita outros meios (gov.br/assinatura avançada) após as alterações de 2024–2025 — **ponto crítico** para assinatura de OS/relatórios NR-12/NR-13.
- Divergência na tabela NR-13 "com SPIE" categorias III–V entre fonte compilada e blog secundário; validar no DOU.

---

## 3. Frota e transporte: CTB/CONTRAN, cronotacógrafo, Lei 13.103/2015, ANTT/RNTRC, produtos perigosos, licenciamento, emissões, multas

### Takeaway
Para frota, as obrigações que geram vencimentos e documentos no CMMS são: verificação do cronotacógrafo a cada 2 anos (Portaria Inmetro 481/2021), CIV/CIPP para produtos perigosos (Res. ANTT 5.998/2022 + Portarias Inmetro 127 e 128/2022, validade variável conforme idade do equipamento, 4–12 meses para CIV), licenciamento anual/CRLV-e, RNTRC, controle de jornada (Lei 13.103/2015) e conformidade PROCONVE/Arla 32. Inspeção técnica veicular periódica obrigatória continua não implantada nacionalmente.

### Cited Findings

**Cronotacógrafo**
- Obrigatoriedade da verificação metrológica prevista na Resolução CONTRAN 92/1999 — [Inmetro FAQ](https://www.gov.br/inmetro/pt-br/acesso-a-informacao/perguntas-frequentes/metrologia-legal/cronotacografo) [verificado via snippet]
- Regulamento técnico metrológico consolidado: **Portaria Inmetro 481, de 06/12/2021**; Portaria Inmetro 91/2022 trata das atividades materiais/acessórias das verificações subsequentes — [Inmetro FAQ — portarias](https://www.gov.br/inmetro/pt-br/acesso-a-informacao/perguntas-frequentes/metrologia-legal/cronotacografo/quais-as-portarias-que-tratam-de-cronotacografos); [Normas Brasil Portaria 481/2021](https://www.normasbrasil.com.br/norma/portaria-481-2021_423946.html)
- **Periodicidade: a cada 2 anos após a verificação inicial e instalação** (subitem 6.3.1 do RTM aprovado pela Portaria 481/2021); certificado com validade de até 2 anos — [Inmetro FAQ periodicidade](https://www.gov.br/inmetro/pt-br/acesso-a-informacao/perguntas-frequentes/metrologia-legal/cronotacografo/apos-a-verificacao-inicial-qual-e-a-periodicidade-de-verificacao-dos-cronotacografos) (FAQ atualizada em 08/01/2026)
- GRU de verificação subsequente: **R$ 90,09**; se reprovado, refaz todo o serviço com nova taxa; selagem em posto autorizado gera declaração de selagem válida por 7 dias para fazer o ensaio — [Inmetro FAQ regularização](https://www.gov.br/inmetro/pt-br/acesso-a-informacao/perguntas-frequentes/metrologia-legal/cronotacografo/o-que-devo-fazer-para-regularizar-meu-cronotacografo); [Paraná — serviço](https://www.parana.pr.gov.br/servicos/Servicos/Solicitacao/Solicitar-verificacao-de-cronotacografo-aPo4WVrm)
- Resolução CONTRAN 938/2022 (consolida regras do registrador instantâneo e inalterável de velocidade e tempo) — [Legisweb](https://www.legisweb.com.br/legislacao/?id=429779); [DOU (Asmetro)](https://asmetro.org.br/portalsn/storage/2022/04/RESOLUCAO-CONTRAN-No-938-DE-28-DE-MARCO-DE-2022-RESOLUCAO-CONTRAN-No-938-DE-28-DE-MARCO-DE-2022-DOU-Imprensa-Nacional.pdf)
- CTB art. 105, II: registrador instantâneo inalterável de velocidade e tempo obrigatório para veículos de escolares, de passageiros com mais de 10 lugares e de carga com PBT > 4.536 kg — [CTB, Lei 9.503/1997](https://www.planalto.gov.br/ccivil_03/leis/l9503compilado.htm) [fonte primária indicada, não re-lida]

**Produtos perigosos (ANTT / Inmetro)**
- Resolução ANTT **5.998, de 03/11/2022** (regulamento do transporte rodoviário de produtos perigosos), vigente desde **01/06/2023**, revogou a Res. 5.947/2021; exige que veículos/equipamentos portem dispositivos de identificação válidos conforme Inmetro; inspeções por OIA geram **CIV** (Certificado de Inspeção Veicular) e **CIPP** (Certificado de Inspeção para o Transporte de Produtos Perigosos) — [Lex — Res. 5.998](https://www.lex.com.br/resolucao-antt-no-5-998-de-3-de-novembro-de-2022/); [texto (Unicamp)](https://depi.unicamp.br/arquivo/uploads/resolucao-5998-3-de-nov-2022)
- Portaria Inmetro 127/2022 (veículos — CIV) e **Portaria Inmetro 128/2022** (requisitos de avaliação da conformidade para inspeção de equipamentos rodoviários de PP; Anexo D = tanques para líquidos); retificação DOU 19/09/2022 (inserção de novos nºs de CIV durante a validade do CIPP); Portaria Inmetro 262/2023 e Portaria 8/2025 autorizam/ajustam inspeções periódicas por OIA (anexos B, C, D, F) — [Normas Brasil 128/2022](https://www.normasbrasil.com.br/norma/portaria-128-2022_436431.html); [Lex 262/2023](https://www.lex.com.br/portaria-inmetro-no-262-de-3-de-julho-de-2023/); [Normas Brasil 8/2025](https://www.normasbrasil.com.br/norma/portaria-8-2025_471880.html)
- Integração ANTT–Inmetro para envio eletrônico dos certificados — [Manual ANTT](https://www.gov.br/antt/pt-br/assuntos/cargas/arquivos/arquivos-produtos-perigosos/manual-de-integracao-antt-e-inmetro-para-envio-dos-certificados.pdf)
- **Validades:** CIV varia com o ano de fabricação do chassi — mínimo 4 meses, máximo 12 meses; CIPP varia com o ano de fabricação do equipamento (faixa citada 4–36 meses, redação ambígua); CTPP (equipamento novo) tem mesma validade legal do CIPP — [Portal do Trânsito (2021, atualizado 2022)](https://www.portaldotransito.com.br/noticias/fiscalizacao-e-legislacao/documentacao-do-veiculo/cargas-perigosas-flexibilizacao-na-exigencia-de-documentos-pode-causar-acidentes-nas-estradas/?noamp=mobile) (fonte jornalística; conferir anexos das Portarias 127/128)
- Infrações: veículo com CIV vencido (código 127/2224), CIPP/CTPP vencido (código 1210/2227) em checklist de fiscalização — [Defesa Civil PR](https://www.defesacivil.pr.gov.br/sites/defesa-civil/arquivos_restritos/files/documento/2021-12/checklist_codigos_autucao.pdf)
- Contêineres-tanque: certificado internacional pode valer até 90 dias após vencimento (CSC/72) — [Portaria Inmetro 329 (UDOP)](https://www.udop.com.br/download/legislacao/seguranca/institucional_site_juridico/portaria_329_transporte_tanque_inmetro.pdf)

**Lei 13.103/2015 (motorista profissional) e CTB** [fonte primária indicada, não re-lida]
- CTB art. 67-C: motorista de carga não pode dirigir mais de 5h30 ininterruptas; 30 min de descanso a cada 6 h (fracionável); 11 h de descanso a cada 24 h (fracionável, garantidas 8 h ininterruptas no 1º período); controle de jornada por diário de bordo, papeleta ou **meios eletrônicos instalados no veículo**; exame toxicológico (art. 148-A) — [Lei 13.103/2015](https://www.planalto.gov.br/ccivil_03/_ato2015-2018/2015/lei/l13103.htm); [CTB compilado](https://www.planalto.gov.br/ccivil_03/leis/l9503compilado.htm)
- Licenciamento anual (CTB art. 130); conduzir veículo não licenciado é infração gravíssima com remoção (art. 230, V); CRLV-e digital; identificação do condutor infrator em **30 dias** (art. 257 §7º, redação da Lei 14.071/2020) — [CTB compilado](https://www.planalto.gov.br/ccivil_03/leis/l9503compilado.htm); [Lei 14.071/2020](https://www.planalto.gov.br/ccivil_03/_ato2019-2022/2020/lei/l14071.htm)
- Inspeção de segurança veicular periódica (CTB art. 104) depende de regulamentação CONTRAN e não está implantada nacionalmente (sem fonte de 2026 verificada).

**ANTT/RNTRC** [fonte primária indicada, não re-lida]
- RNTRC regulado pela Resolução ANTT 5.982/2022; transportador deve manter veículos cadastrados (placa/RENAVAM) e registro válido; CIOT/pagamento de frete eletrônico — [ANTT RNTRC](https://www.gov.br/antt/pt-br/assuntos/cargas/rntrc)

**Emissões — PROCONVE / Arla 32 / fumaça preta** [fonte primária indicada, não re-lida]
- PROCONVE fase P8 (equivalente Euro VI) obrigatória para veículos pesados novos desde 01/01/2022 (Resolução CONAMA 490/2018); fase P7 (2012) introduziu SCR com Arla 32 — [CONAMA 490/2018 (Ibama)](https://www.gov.br/ibama/pt-br/assuntos/emissoes-e-residuos/emissoes/programa-de-controle-de-emissoes-veiculares-proconve)
- Portaria IBAMA 85/1996: empresas com frota a diesel devem manter Programa Interno de Autofiscalização da Correta Manutenção da Frota quanto à emissão de fumaça preta (opacidade/escala Ringelmann), com registros — [Ibama](https://www.gov.br/ibama/pt-br/assuntos/emissoes-e-residuos/emissoes)

### Inferences
- **Entidade "Documento/Certificado com vencimento" genérica** por veículo/equipamento: CRLV-e (exercício), IPVA (parcelas por UF), seguro/RCTR-C, RNTRC, cronotacógrafo (2 anos, nº certificado, GRU), CIV, CIPP/CTPP (validade lida do certificado, não calculada — varia por idade), AET (autorização especial de trânsito, se aplicável), laudo de opacidade (programa Ibama), tacógrafo/telemetria. Alertas D-60/D-30/D-7 e bloqueio de despacho do veículo vencido (integração com TMS).
- **Não calcular validade CIV/CIPP automaticamente** — exigir upload e data de validade do certificado do OIA; o sistema só sugere janela com base no ano de fabricação.
- Multas: importação RENAINF/Detran via integradores; workflow de indicação de condutor com prazo de 30 dias, cruzando jornada/telemetria (quem dirigia) — reduz pontuação indevida no CNPJ e NIC (multa por não indicação).
- Arla 32: registrar consumo de Arla por veículo (razão Arla/diesel ~3–5% é heurística de mercado, não verificada) para detectar fraude de SCR; OS de "emulador/desativação SCR" como não conformidade grave.
- Controle de jornada (Lei 13.103) normalmente é do TMS/telemetria, mas o CMMS deve registrar o motorista do checklist pré-viagem e vincular ocorrências de manutenção.

### Gaps
- Não obtive a tabela oficial de validade CIV/CIPP por idade do equipamento (Anexos das Portarias Inmetro 127/128/2022).
- Não verifiquei se em 2025–2026 houve nova resolução CONTRAN sobre inspeção técnica veicular periódica, nem mudanças de 2026 no CRLV-e/licenciamento.
- Não confirmei a vigência atual da Portaria IBAMA 85/1996 nem a numeração da Resolução ANTT do RNTRC.
- Não verifiquei a Lei 14.599/2023 e atualizações de 2024–2026 sobre exame toxicológico.

---

## 4. Outros: LGPD, fiscal (NF-e, CIAP, depreciação), assinatura eletrônica, ANVISA RDC 509/2021

### Takeaway
Todos os itens abaixo vêm de legislação consolidada que não foi re-lida nesta sessão; servem como direção de requisitos e precisam de validação jurídica/contábil. Os pontos de projeto são: base legal e minimização para geolocalização/biometria de motoristas e técnicos (LGPD); vincular NF-e/NFS-e de peças e serviços a OS e ativos (com campos IBS/CBS da reforma tributária a partir de 2026); distinguir despesa x capitalização (reparos que estendem vida útil > 1 ano); e assinatura eletrônica com evidências de autoria/integridade (MP 2.200-2 art. 10 §2º) para checklists e encerramento de OS.

### Cited Findings [fonte primária indicada, não re-lida]
- **LGPD (Lei 13.709/2018):** bases legais (art. 7º) aplicáveis — cumprimento de obrigação legal (ex.: controle de jornada, Lei 13.103), execução de contrato de trabalho, legítimo interesse (segurança patrimonial/rastreamento, exige teste de balanceamento); dado biométrico e de saúde (exame toxicológico, ASO) são **dados sensíveis** (art. 5º, II; art. 11); direitos do titular (art. 18); relatório de impacto (art. 38) — [Lei 13.709/2018](https://www.planalto.gov.br/ccivil_03/_ato2015-2018/2018/lei/l13709compilado.htm)
- ANPD: Regulamento de comunicação de incidente de segurança (Res. CD/ANPD 15/2024 — prazo de 3 dias úteis) e Regulamento do Encarregado (Res. CD/ANPD 18/2024) — [ANPD](https://www.gov.br/anpd/pt-br)
- **Assinatura eletrônica:** MP 2.200-2/2001 art. 10 §1º (documentos com certificado ICP-Brasil presumem-se verdadeiros) e §2º (outros meios de comprovação de autoria e integridade são válidos se admitidos pelas partes ou aceitos pela pessoa a quem for oposto) — [MP 2.200-2/2001](https://www.planalto.gov.br/ccivil_03/mpv/antigas_2001/2200-2.htm); Lei 14.063/2020 classifica assinaturas em **simples, avançada e qualificada (ICP-Brasil)** — aplica-se às interações com entes públicos, mas a taxonomia é usada como referência — [Lei 14.063/2020](https://www.planalto.gov.br/ccivil_03/_ato2019-2022/2020/lei/l14063.htm)
- **Fiscal:** crédito de ICMS de ativo imobilizado apropriado em **48 parcelas mensais** (LC 87/1996 art. 20 §5º — controle via CIAP/bloco G do SPED) — [LC 87/1996](https://www.planalto.gov.br/ccivil_03/leis/lcp/lcp87.htm); taxas de depreciação fiscal (IN RFB 1.700/2017, Anexo III): caminhões e veículos de carga 25% a.a. (4 anos), automóveis 20% a.a., máquinas/equipamentos em geral 10% a.a. — [IN RFB 1.700/2017](http://normas.receita.fazenda.gov.br/sijut2consulta/link.action?idAto=81268); RIR/2018 (Decreto 9.580/2018): gastos com reparos/conservação que aumentem a vida útil em mais de 1 ano devem ser capitalizados — [Decreto 9.580/2018](https://www.planalto.gov.br/ccivil_03/_ato2015-2018/2018/decreto/d9580.htm)
- **Reforma tributária (LC 214/2025):** 2026 é ano de teste com destaque de CBS 0,9% e IBS 0,1% nos documentos fiscais e NFS-e em padrão nacional — [LC 214/2025](https://www.planalto.gov.br/ccivil_03/leis/lcp/lcp214.htm)
- **ANVISA RDC 509/2021:** gerenciamento de tecnologias em saúde em estabelecimentos de saúde (substituiu a RDC 2/2010): plano de gerenciamento, inventário, registros de manutenção preventiva/corretiva, rastreabilidade e responsável técnico — [ANVISA RDC 509/2021](https://www.in.gov.br/en/web/dou/-/resolucao-de-diretoria-colegiada-rdc-n-509-de-27-de-maio-de-2021-323002775)

### Inferences
- **LGPD:** geolocalização contínua apenas em horário de trabalho/veículo da empresa; política de retenção configurável; perfis de acesso (mecânico não vê dado de saúde do motorista); log de acesso; RIPD modelo fornecido pela Tenaz; Tenaz atua como operadora e o cliente como controlador (DPA no contrato).
- **Assinatura de OS/checklist:** assinatura eletrônica avançada própria (login + geolocalização + timestamp + hash do documento + foto/selfie opcional + trilha de auditoria imutável) atende o art. 10 §2º para uso interno; oferecer integração ICP-Brasil (certificado A1/A3 ou nuvem) para documentos NR-13/laudos do PLH e quando a NR-1 exigir.
- **Fiscal:** importação de XML de NF-e (modelo 55) e NFS-e (padrão nacional) com rateio de itens para OS/ativo/centro de custo; flag "capitalizável" para reformas/overhaul; cadastro do ativo com data de aquisição, valor, vida útil fiscal e societária para custo total de propriedade (TCO).
- **RDC 509/2021:** relevante somente se a Tenaz atender clínicas/hospitais (segmento "serviços"); pode virar template de plano.

### Gaps
- Nenhum dos itens desta seção foi re-lido em 2026 nesta pesquisa; números de resoluções ANPD, alíquotas-teste e taxas de depreciação devem ser validados por contador/jurídico.
- Não verifiquei orientações ANPD específicas sobre telemetria veicular/rastreamento de empregados.

---

## 5. KPIs: fórmulas e valores de referência

### Takeaway
Os KPIs devem ser calculados a partir de eventos e tempos registrados na OS e no estado do ativo, com definições alinhadas a NBR 5462/EN 15341/SMRP. Benchmarks públicos confiáveis são escassos: ABRAMAN (2022) indica custo de manutenção/faturamento concentrado até 4% (comum < 2%; 5–7% em setores intensivos em ativos); para os demais, só há valores "classe mundial" divulgados por fornecedores (PM compliance ≥ 90%, trabalho planejado ≥ 85%, wrench time 50–55%, custo/RAV 1–3%).

### Cited Findings
- ABRAMAN, Documento Nacional 2022: empresas destinam parcela pequena do faturamento bruto à manutenção, concentrada **até 4%**, sendo comum **< 2%**; setores intensivos em ativos chegam a **5–7%**; **> 10% é pouco frequente** — [Monografia Poli-USP citando ABRAMAN (2022)](https://poli-integra.poli.usp.br/wp-content/uploads/2026/02/2025_Igor-Gama-Galante.pdf)
- Documento Nacional ABRAMAN é bienal (edições 2022 e 2024 analisadas em estudo comparativo, mostrando avanço em diagnóstico de maturidade e uso de indicadores) — [mesma monografia / busca](https://poli-integra.poli.usp.br/wp-content/uploads/2026/02/2025_Igor-Gama-Galante.pdf)
- Indicador CMF = (Custo total de manutenção / Faturamento bruto) × 100 — [FM2S](https://www.fm2s.com.br/blog/indicadores-de-manutencao); [Tractian](https://tractian.com/blog/indicadores-de-manutencao)
- MC/RAV (SMRP) = custo anual total de manutenção / valor de reposição dos ativos × 100 — [Fiix](https://fiixsoftware.com/maintenance-metrics/mcrav/)
- Benchmarks (fornecedor Tractian): PM compliance típico 60–75%, melhor ≥ 90%; manutenção planejada típico 55–65%, melhor ≥ 85%; wrench time típico 25–35%, melhor 50–55%; custo/RAV típico 3–5%, melhor 1–3% — [Tractian maintenance statistics](https://tractian.com/en/glossary/maintenance-statistics)
- Oxmaint: custo/RAV ótimo 2–3%; PM compliance ≥ 95% e wrench time ≥ 55% como "world class" em FMCG — [Oxmaint manufatura](https://oxmaint.com/industries/manufacturing-plant/manufacturing-plant-maintenance-metrics-that-matter-2026); [Oxmaint FMCG](https://oxmaint.com/industries/fmcg/maintenance-kpis-fmcg-plants-metrics-track)
- Reliamag adverte que benchmarks de schedule/PM compliance, wrench time e PMP são "commonly cited: verify against your standard" e recomenda baseline próprio — [Reliamag KPI reference](https://reliamag.com/guides/maintenance-reliability-kpi-reference/)
- Literatura acadêmica brasileira usa ABRAMAN como benchmark para disponibilidade, MTBF, MTTR e CMF — [Tractian glossário KPI](https://tractian.com/glossario/maintenance-kpi); [TCC UFC 2022](https://repositorio.ufc.br/bitstream/riufc/70869/1/2022_tcc_sneto.pdf)

### Fórmulas recomendadas (definições de mercado/normas; formalização para regra de negócio)
| KPI | Fórmula | Observações para o CMMS |
|---|---|---|
| Disponibilidade inerente (Ai) | MTBF / (MTBF + MTTR) | Só tempo de reparo ativo; ignora logística e preventiva |
| Disponibilidade técnica/alcançada (Aa) | MTBM / (MTBM + M̄) — inclui preventiva e corretiva | M̄ = tempo médio de manutenção ativa |
| Disponibilidade operacional (Ao) | Tempo disponível / Tempo requerido = (TR − tempo indisponível total) / TR | Inclui espera de peça, logística, aprovação; é a de frota |
| Disponibilidade física (frota/mineração) | (HC − HM) / HC × 100, HC = horas calendário, HM = horas em manutenção | Termo usado em frota/equipamentos móveis |
| MTBF | Tempo total de operação / nº de falhas (que pararam o ativo) | Por ativo, classe e modo de falha; para frota pode ser **MKBF** (km entre falhas) |
| MTTR | Σ tempo de reparo / nº de reparos | Definir se inclui espera de peça (MTTR "total") ou só reparo ativo |
| MTTA | Σ (início do atendimento − abertura) / nº de chamados | Requer timestamp de reconhecimento |
| Confiabilidade R(t) | e^(−t/MTBF) (exponencial) ou e^(−(t/η)^β) (Weibull) | Weibull por modo de falha |
| OEE | Disponibilidade × Performance × Qualidade | Só manufatura; requer apontamento de produção (integração MES/PCP) |
| Backlog (semanas) | Σ HH estimadas de OS liberadas pendentes / HH disponíveis por semana da equipe | Referência de mercado 2–4 semanas (não verificada) |
| Aderência/PM compliance | OS preventivas concluídas no prazo (janela, ex.: ±10% do intervalo) / OS preventivas programadas × 100 | SMRP 5.4.14 |
| Schedule compliance | HH executadas da programação semanal / HH programadas × 100 | |
| % planejada vs. corretiva | HH (ou nº OS) de preventiva+preditiva+planejada / total × 100; corretiva emergencial / total | |
| Custo de manutenção / RAV | Custo anual total / valor de reposição × 100 | Requer RAV cadastrado |
| CMF (custo/faturamento) | Custo total manutenção / faturamento bruto × 100 | ABRAMAN |
| Custo por km | (Manutenção + pneus + combustível + Arla + lubrificantes [+ fixos]) / km rodado | Separar custo variável de manutenção (R$/km) |
| Consumo | km rodados / litros abastecidos (km/l) | Requer abastecimentos + hodômetro/telemetria |
| CPK pneu | (Custo de compra + recapagens + consertos − valor residual) / km total rodado pelo pneu (geralmente R$/1.000 km) | Controle por nº de fogo, vida (nova, R1, R2), posição e sulco |
| Wrench time | Tempo em execução direta / tempo total pago × 100 | Requer apontamento/amostragem |
| First-time fix rate | OS resolvidas sem retorno em X dias (ex.: 30) / OS concluídas × 100 | Requer vínculo de retrabalho |

### Inferences
- O motor de KPIs só é confiável se a OS distinguir: tipo (preventiva, preditiva, corretiva planejada, corretiva emergencial, melhoria, inspeção legal), "falha com parada" S/N, e tempos (abertura, reconhecimento, início, pausas com motivo, fim, liberação).
- Exibir sempre a definição usada ao lado do KPI (tooltip "como calculamos") e permitir parametrizar janela de aderência e inclusão/exclusão de espera de peça no MTTR — é onde clientes divergem.
- Benchmarks na UI devem ser rotulados com fonte (ABRAMAN 2022 para CMF; "referência de mercado" para os demais).

### Gaps
- Não obtive os números do Documento Nacional ABRAMAN 2024 (disponibilidade operacional média, distribuição corretiva/preventiva/preditiva, pessoal próprio x terceiros); o PDF da ABRAMAN não é público.
- Não acessei o texto oficial da SMRP para valores-alvo; os benchmarks acima são de fornecedores.
- Benchmarks brasileiros de custo/km, km/l e CPK por categoria de veículo (ex.: CNT, NTC&Logística, Cepea/ESALQ-LOG) não foram pesquisados nesta rodada.
- Referência de backlog de 2–4 semanas e razão Arla/diesel não foram verificadas.
