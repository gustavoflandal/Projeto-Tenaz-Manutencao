# Especificação de Requisitos — Tenaz

**Versão:** 0.7 · **Data:** 13/11/2026 · **Status:** Em elaboração
**Produto:** Tenaz · Gestão de Manutenção e Frota
Alinhado ao sumário executivo 0.8.

**Mudanças da versão 0.7:** a Tenaz é **multipaís e bilíngue desde o MVP** (decisão do dono do produto, sumário 0.8, seção 3.6, linhas "Países e idiomas" e "País, identificação, moeda e idioma"; ADR-013 do documento 5; [etapa 014](etapas/014_2026-11-13_multipais-e-idiomas.md)). **Requisitos novos:** RNF-043 (idiomas suportados e escolha do idioma), RNF-044 (países suportados e perfil de país), RNF-045 (dinheiro e moedas) e RNF-046 (conteúdo legal por país e jurisdição), na nova seção 5.12, cada um com critério de aceite e rastreio para o código e os testes existentes. Os não funcionais passam de 42 para **46**. **Redação ajustada, sem renumerar:** RNF-029 (idioma e formatos, que deixa de ser só `pt-BR` e R$), RNF-006 e RNF-013 (nota de país), RNF-021 (identificação pessoal de qualquer país e proteção de dados de AR e PY); RF-ADM-001, RF-ADM-002, RF-ADM-004, RF-ADM-006, RF-ADM-009, RF-ADM-010, RF-ADM-011, RF-ADM-012, RF-ADM-013, RF-ADM-015, RF-TER-001, RF-SRV-001, RF-SRV-003, RF-SRV-004, RF-SRV-012, RF-SRV-014 e RF-LOC-001 passam de CNPJ, CPF, CEP, UF e R$ fixos para a forma por país; RF-OS-001 e RF-SRV-007 (prefixo do número da OS pelo país). A seção 3 ganha a convenção de país, idioma e dinheiro. **Notas de país**, sem reescrever conteúdo legal, nas introduções dos módulos EST, TER, SRV, LOC, FRO e CON e em RF-EST-009, RF-EST-015, RF-FRO-010, RF-INT-009 e RF-CON-005, apontando as pendências **P12**, **P13**, **P14** e **P15**. Nenhum requisito funcional novo, nenhum ID apagado ou renumerado; sem mudança de fase, de plano ou de contagem dos RF.

**Mudanças da versão 0.6:** decisão do dono do produto (sumário 0.7, seção 3.6, "Dados de saúde"): **a Tenaz não guarda dado de saúde**. ASO, exame toxicológico, a permissão `pessoa.saude.ler` e as regras de dado sensível de saúde foram removidos do escopo em 2026-11-13 por decisão do dono do produto. O catálogo de RF-ADM-005 passa a ter 116 códigos, todos no formato `<módulo>.<ação>`, sem exceção. Ajustados RF-ADM-005, RF-ADM-011, RF-ADM-015, RF-SRV-003, RF-SRV-014, RF-FRO-001, RF-CON-005, RNF-021 e as seções 2 e 3. Nenhum requisito ficou inteiro sem objeto; nenhum ID foi apagado nem renumerado. Reforço da vedação (sumário 3.6, "Dado de saúde recebido por engano"): aviso fixo em RF-ADM-015, RF-FRO-001 e RNF-033 e recusa de tipo de documento ligado a saúde em RF-CON-005, com critérios de aceite. Seção 7: serviços que custam por natureza são decididos na etapa em que entram.

**Mudanças da versão 0.5:** catálogo de permissões e matriz dos doze perfis padrão, como subseção de RF-ADM-005, a que o documento 3 (seção 4.8) remete; regra explícita do perfil `administrador` e da permissão `srv.gerenciar_conexao` na conta da subcontratada. Redação de RF-ADM-001, RF-ADM-002, RF-ADM-004, RF-ADM-005, RF-ADM-006, RF-ADM-011, RF-ADM-015, RF-MOB-001, RF-SRV-001 e RF-SRV-003 ajustada às decisões da etapa E01: calendário da unidade em `turno_unidade` e `feriado_unidade` e login por telefone e código adiado pela pendência **P11**. Decisões do dono do produto no sumário 0.6, seção 3.6 (catálogo de permissões canônico, situação do usuário, criação da conta, MFA da empresa, CNPJ): responsável nasce `administrador` e `gestor`, turnos da unidade inicial informados na criação sem padrão imposto, CNPJ único entre todas as empresas com reativação da cancelada (proposta), `usuario.situacao`, permissão avulsa em `usuario_permissao` com a regra RN-ADM-026 para usuários externos, `empresa.mfa_obrigatorio`; vínculo (`empregado`, `terceiro`, `autonomo`) e situação (`ativa`, `inativa`, sem `desligada`) da `pessoa` em RF-ADM-015. Nenhum requisito mudou de fase nem de plano; o único escopo ampliado é o de RF-ADM-011, por decisão do dono do produto dentro da etapa E01. Trilha de auditoria ampliada (sumário 0.6, seção 3.6, "Trilha de auditoria configurável"): RF-ADM-011 com tela e módulo, nove tipos de operação, captura configurável por tela (`captura_auditoria`), hash por evento e selo encadeado (`selo_auditoria`) em armazenamento imutável, telas de consulta, verificação e configuração, e critérios de aceite; coerência em RF-ADM-013 e RNF-019. `srv.gerenciar_conexao` no grupo Conexão e pontos 39 e 42 a 44 resolvidos pelo documento 3, versão 0.5; pontos 39 a 44, 46 e 47 resolvidos (seção 8.4); fica em aberto só o ponto 45 (seção 8.3).

**Mudanças da versão 0.4:** configuração por segmento (sumário 3.7). No MVP: regras configuráveis de condição e efeito (RF-ADM-016), fase de ciclo de vida do ativo e entrada em operação (RF-ATV-013), campos extras da OS (RF-OS-021) e o tipo de OS `implantacao`. Na Fase 2: pacotes de segmento (RF-ADM-017) e os módulos novos PRJ, projetos de implantação (RF-PRJ-001 a RF-PRJ-007), e CTR, contratos com clientes (RF-CTR-001 a RF-CTR-005). Origem: caso de um fornecedor de radares instalados e mantidos por prestadores em várias cidades ([etapa 012](etapas/012_2026-11-12_configuracao-por-segmento.md)).

**Mudanças da versão 0.2:** novo módulo **SRV** (gestão de serviços do prestador, MVP, plano Serviços), com 21 requisitos (RF-SRV-001 a RF-SRV-021) e dois requisitos não funcionais para a conexão entre contas (RNF-041 e RNF-042, risco R5). Aplicadas as decisões do sumário 3.6: permissão `os.encerrar`, escopos de acesso externo, ação `assumir`, token do link de uso único no fragmento da URL, desistência e "reabrir como nova" na solicitação, fim do contrato de locação por `encerrar`, categorias canônicas da matriz de responsabilidade, campos e origens da OS, papéis `posto_combustivel` e `seguradora`, peça em trânsito, permissão `est.decidir_destino`, distinção entre `bloqueio_ativo`, aptidão para despacho e `bloqueio_energia`, locatário só no portal web.

Este documento detalha os requisitos funcionais e não funcionais da Tenaz. Ele segue o [sumário executivo](0_SUMARIO_EXECUTIVO.md) para códigos de módulo, nomes de entidade, estados, planos e fases, e o [estudo de mercado](pesquisa/reports/Benchmarking%20manuten%C3%A7%C3%A3o%20e%20frota%20Brasil.md) para normas, indicadores e metas. As regras de negócio aparecem em texto dentro de cada requisito. A numeração formal das regras (`RN-<MÓDULO>-<NNN>`) e os fluxos ficam em [4_FLUXOS_PROCESSO.md](4_FLUXOS_PROCESSO.md).

---

## Índice

1. [Visão geral](#1-visão-geral)
2. [Atores](#2-atores)
3. [Convenções deste documento](#3-convenções-deste-documento)
4. [Requisitos funcionais](#4-requisitos-funcionais)
   - [ADM — Administração](#41-adm--administração)
   - [ATV — Ativos](#42-atv--ativos)
   - [MED — Medidores](#43-med--medidores)
   - [SOL — Solicitações](#44-sol--solicitações)
   - [OS — Ordens de serviço](#45-os--ordens-de-serviço)
   - [PRE — Planos de manutenção](#46-pre--planos-de-manutenção)
   - [CHK — Checklists](#47-chk--checklists)
   - [EST — Estoque e compras](#48-est--estoque-e-compras)
   - [TER — Terceiros](#49-ter--terceiros)
   - [SRV — Gestão de serviços do prestador](#410-srv--gestão-de-serviços-do-prestador)
   - [LOC — Locação](#411-loc--locação)
   - [FRO — Frota](#412-fro--frota)
   - [CON — Conformidade](#413-con--conformidade)
   - [IND — Indicadores e relatórios](#414-ind--indicadores-e-relatórios)
   - [INT — Integrações](#415-int--integrações)
   - [IA — Assistente](#416-ia--assistente)
   - [MOB — App de campo](#417-mob--app-de-campo)
   - [PRJ — Projetos e implantação](#418-prj--projetos-e-implantação)
   - [CTR — Contratos com clientes](#419-ctr--contratos-com-clientes)
5. [Requisitos não funcionais](#5-requisitos-não-funcionais)
6. [Resumo por módulo e fase](#6-resumo-por-módulo-e-fase)
7. [Pendências citadas](#7-pendências-citadas)
8. [Pontos a alinhar com o sumário executivo](#8-pontos-a-alinhar-com-o-sumário-executivo)

---

## 1. Visão geral

A Tenaz é um SaaS multiempresa de gestão de manutenção de equipamentos e frota. Uma empresa cliente cadastra máquinas, instalações, veículos, implementos e componentes numa única base hierárquica, com medidores múltiplos. Sobre essa base, o sistema gera e controla solicitações, ordens de serviço, planos preventivos, checklists, estoque, custos, documentos com vencimento e indicadores. Técnicos e motoristas usam um app que opera sem conexão por pelo menos 7 dias. Prestadores externos (subcontratadas, assistências técnicas, oficinas) executam OS com acesso restrito, mantêm estoque da empresa sob controle e consertam peças rastreadas; ativos locados, cedidos ou tomados, têm contrato, vistoria e responsável pelo custo da manutenção. A subcontratada ou assistência técnica também pode ter **conta própria** na Tenaz (módulo SRV, plano Serviços): cadastra técnicos, clientes e equipamentos de clientes, distribui OS, registra serviço, componentes e horas no app e, quando o cliente também usa a Tenaz, conecta as duas contas para que a OS trafegue entre elas (sumário 3.5). Empresas que instalam e mantêm equipamentos próprios em clientes, como radares de fiscalização eletrônica em várias cidades, controlam a implantação como projeto com etapas (PRJ) e o contrato com o cliente, com meta de disponibilidade e medição mensal (CTR). Cada empresa ajusta a Tenaz ao seu segmento com campos extras, regras de condição e efeito e pacotes prontos, sempre de catálogos fechados (sumário 3.7). O produto é **multipaís e bilíngue desde o MVP**: empresas do Brasil, da Argentina e do Paraguai, com interface, e-mails, notificações, documentos gerados e mensagens de erro da API em português do Brasil ou em espanhol latino-americano neutro, por escolha de cada usuário (RNF-043 a RNF-046; sumário 3.6).

Três ideias orientam os requisitos:

1. **Um só modelo de ativo** para equipamento e veículo. Pneu, combustível, multa e sinistro são extensões de frota sobre esse modelo, não um segundo sistema.
2. **Dados que sustentam KPIs e fiscalização** nascem obrigatórios na OS: tipo, timestamps, "o ativo parou?", falha, peças, horas e, em máquina NR-12, os sete campos do item 12.11.2.
3. **Nada é gravado por IA sem confirmação humana**, e todo indicador mostra como é calculado.

O porte-alvo é de 50 a 2.000 ativos e de 2 a 50 técnicos por empresa (sumário, seção 1.3).

---

## 2. Atores

Os perfis são os definidos no sumário (seção 3.1). Um usuário pode ter mais de um perfil. Quem executa, conduz ou responde por algo é uma `pessoa` (RF-ADM-015), com ou sem `usuario`: um condutor eventual ou um técnico de fornecedor pode existir no sistema sem login.

| Perfil | Quem é | Cobrado? | Uso principal |
|---|---|---|---|
| `administrador` | Responsável pela conta na empresa cliente | Não, salvo se também tiver perfil `tecnico` | Empresa, unidades, pessoas, usuários, perfis, plano, parâmetros, integrações |
| `gestor` | Gerente de manutenção ou de frota | Não, salvo se também tiver perfil `tecnico` | Indicadores, aprovações, liberação de OS, custos, conformidade |
| `planejador` | PCM, programador de manutenção | Não, salvo se também tiver perfil `tecnico` | Triagem, planejamento e programação de OS, liberação de OS, planos preventivos |
| `tecnico` | Mecânico, eletricista, executante | **Sim**, por técnico/mês | Execução de OS no app, apontamento, peças, checklists |
| `almoxarife` | Responsável pelo almoxarifado | Não | Itens, saldos, movimentações, reservas, requisições, inventário |
| `comprador` | Responsável por compras | Não, gratuito e ilimitado | Requisições, cotações, pedidos de compra, entrada de notas |
| `plh` | Profissional legalmente habilitado (NR-13) | Não | Prazos de inspeção, Registro de Segurança, relatórios de inspeção |
| `motorista` | Condutor de veículo da frota | Não, gratuito e ilimitado | Checklist do veículo, leitura de km, solicitação, abastecimento |
| `solicitante` | Qualquer pessoa que pede manutenção | Não, gratuito e ilimitado | Abrir e acompanhar solicitações |
| `visualizador` | Diretoria, auditoria, cliente interno | Não, gratuito e ilimitado | Consulta de painéis, ativos e OS, sem edição |
| `prestador` | Usuário de um fornecedor (subcontratada, assistência técnica, oficina), ligado por `usuario.fornecedor_id` | **Sim**, como técnico, quando executa OS | Ver e executar só as OS atribuídas ao seu fornecedor; movimentar só o almoxarifado `em_terceiro` desse fornecedor (Fase 2) |
| `locatario` | Usuário de um cliente de locação, ligado ao `locatario` | Não, gratuito e ilimitado | Abrir solicitações e acompanhar OS só dos ativos locados a ele (Fase 2) |

A liberação (encerramento) de OS é a permissão `os.encerrar`, presente por padrão em `gestor` e `planejador`, nunca em `prestador` nem `locatario`, e pode ser dada a perfis personalizados (sumário 3.6). A Tenaz não guarda dado de saúde da pessoa (sumário 3.6).

**Na conta de uma subcontratada (SRV).** A conta da subcontratada é uma `empresa` com o módulo SRV e usa os mesmos perfis, com o significado do lado de quem presta o serviço: `administrador` e `gestor` cuidam da conta, dos clientes e das conexões; `planejador` abre, distribui e reatribui OS; `tecnico` executa no app. Só o `tecnico` é cobrado, pelo preço do plano Serviços (RF-SRV-001). O técnico da subcontratada com conta própria **não** conta como técnico cobrado na conta do cliente contratante; o perfil `prestador` dentro da conta do cliente continua existindo para subcontratadas sem conta própria (sumário 2.2).

Atores sem perfil de usuário:

| Ator | Descrição |
|---|---|
| Prestador ocasional | Técnico de fornecedor sem login que registra a execução de uma OS por link seguro de uso único (RF-TER-005). Não é cobrado. |
| Solicitante anônimo | Pessoa que abre solicitação pelo portal QR sem login (RF-SOL-004) |
| Representante do cliente atendido | Pessoa do `cliente_servico` que acompanha o serviço da subcontratada e assina a OS no aparelho do técnico, com nome e documento, sem login (RF-SRV-012). Recebe o relatório da OS em PDF (RF-SRV-013) |

Atores não humanos:

| Ator | Descrição |
|---|---|
| Sistema | Rotinas agendadas: geração de OS preventiva, alertas de vencimento, consolidação de KPIs, prazos de multa |
| Integrador | Sistema externo autenticado na API pública (ERP, TMS, BI, telemática) |
| Dispositivo de telemática ou IoT | Fonte de leituras, DTCs e eventos por conector |
| Processo de sincronização entre contas | Rotina da plataforma, com papel de banco próprio, que copia entre duas contas conectadas só os campos permitidos pela `conexao_empresa`, gravando em cada conta numa transação separada (RF-SRV-014 a RF-SRV-019, RNF-041, RNF-042) |

---

## 3. Convenções deste documento

- **ID:** `RF-<MÓDULO>-<NNN>` para requisitos funcionais e `RNF-<NNN>` para não funcionais, conforme o sumário (seção 3.4).
- **Fase:** MVP (meses 0 a 6), Fase 2 (6 a 12) ou Fase 3 (12 a 24), conforme o sumário (seção 4).
- **Plano mínimo:** o plano mais barato que inclui o requisito (Essencial, Profissional ou Enterprise), conforme a tabela de recursos do sumário (seção 2.1), que é a única fonte de qual plano inclui o quê. Os planos superiores herdam tudo dos inferiores. Quando a tabela não cita o recurso, o plano vem marcado como "(proposta)". Os requisitos do módulo SRV têm plano mínimo **Serviços**, o plano das contas de subcontratadas (sumário 2.2), que não herda nem é herdado pelos outros três; a parte de um requisito SRV executada na conta do cliente contratante indica também o plano exigido do cliente.
- **Permissões:** formato `<módulo>.<ação>` (ex.: `os.encerrar`, `est.decidir_destino`), sem exceção. A antiga exceção `pessoa.saude.ler` foi removida do escopo em 2026-11-13 por decisão do dono do produto (sumário 3.6).
- **(proposta):** valor, regra ou escolha que não vem do sumário nem do estudo de mercado. Precisa de validação com clientes ou com o time.
- **Pendências:** P1 a P15 são as pendências do sumário (seção 6). O requisito que depende de uma delas cita a pendência. **P12** (normas e documentos da Argentina e do Paraguai), **P13** (fiscal), **P14** (proteção de dados e transferência internacional) e **P15** (preço, moeda de cobrança e impacto no prazo e no time) vêm do escopo multipaís.
- **País, idioma e dinheiro:** onde um requisito cita CNPJ, CPF, CEP, UF, R$ ou o prefixo `OS-`, vale a forma por país do ADR-013: identificação por tipo e valor (`identificacao_tipo` e `identificacao`), `subdivisao` e `codigo_postal`, `moeda` ao lado de todo valor, e o prefixo pelo país da empresa. O nome do tipo (por exemplo, CNPJ) só aparece quando o requisito é, de fato, daquele país. Requisito de módulo futuro que depende de país só ganha nota apontando a pendência; o conteúdo legal não é reescrito antes da pesquisa e do parecer externo (RNF-046).
- **Nomes técnicos** entre crases (`ordem_servico`, `em_execucao`) são os canônicos do sumário. O identificador de empresa é sempre `empresa_id`. Rotas seguem `/v1/<recurso-no-plural>`; `GET /v1/versao` é a exceção aceita (sumário 3.6).
- **Estrutura de cada requisito:** descrição, dados principais, operações, validações e regras. Campos sem conteúdo relevante são omitidos.
- **IDs movidos não são reutilizados.** Um requisito que mudou de módulo fica com uma entrada curta no lugar antigo, apontando para o novo ID, e não entra nas contagens.

Requisitos movidos na versão 0.2:

| ID antigo | ID novo | Motivo |
|---|---|---|
| RF-EST-010 · Fornecedores | RF-TER-001 · Fornecedor e papéis | `fornecedor` passa a ser o cadastro único de empresas externas, com papéis, dono do módulo TER (sumário 3.1) |
| RF-FRO-012 · Oficina externa | RF-TER-008 · Oficina externa e assistência técnica | A oficina externa sai de FRO e passa a ser um caso de TER (sumário, mudanças da versão 0.3) |

---

## 4. Requisitos funcionais

### 4.1 ADM — Administração

#### RF-ADM-001 · Cadastro da empresa e primeiros passos

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** criação da conta da empresa cliente (tenant) e roteiro guiado de primeiros passos, para uso sem treinamento formal.
- **Dados principais:** país (`empresa.pais`: `BR`, `AR` ou `PY`), identificação fiscal (`identificacao_tipo` e `identificacao`: `cnpj` no Brasil, `cuit` na Argentina, `ruc` no Paraguai), razão social, nome fantasia, logo, moeda (`empresa.moeda`, dada pelo país), idioma padrão (`empresa.idioma_padrao`, sugerido pelo país), fuso horário (nome IANA da lista do país), plano contratado, ciclo de cobrança (mensal ou anual), contato do responsável, unidade inicial com os seus turnos de trabalho, MFA obrigatório para a empresa (`empresa.mfa_obrigatorio`, RF-ADM-006).
- **Operações:** criar conta pelo site, escolhendo o país (sugerido pela região do navegador) e informando os turnos da unidade inicial (modelos prontos, como comercial ou 24 horas por 7 dias, só preenchem o formulário); editar dados, inclusive a identificação fiscal e o idioma padrão, pelo `administrador`, com registro na trilha (sumário 3.6); o país e a moeda não são editáveis (RNF-044, RNF-045); seguir o roteiro de primeiros passos (cadastrar ou importar ativos, convidar técnicos, criar o primeiro plano preventivo, abrir a primeira OS); carregar dados de exemplo opcionais e apagá-los depois.
- **Validações:** país entre os suportados (RNF-044); identificação do tipo da empresa naquele país (`cnpj`, `cuit` ou `ruc`), normalizada e com a forma e o dígito verificador conferidos como o perfil do país manda (CNPJ, inclusive o alfanumérico, com dígito; CUIT e RUC só a forma, até a **P13**); identificação única entre todas as empresas da Tenaz por `(identificacao_tipo, identificacao)`, inclusive as canceladas: a empresa cancelada que volta é reativada, não recriada (proposta, coerente com o modelo de dados e com RN-ADM-017 do documento 4); fuso da lista do país; idioma padrão entre os suportados (RNF-043); unidade inicial com pelo menos um turno; e-mail do responsável confirmado antes do primeiro acesso.
- **Regras:** todo dado do sistema pertence a uma empresa e só é visível a ela. A empresa nasce com os doze perfis padrão, já com as permissões da matriz de RF-ADM-005, e o responsável recebe os perfis `administrador` e `gestor`, para que uma empresa com uma só pessoa já libere OS (sumário 3.6). O sistema não impõe calendário padrão à unidade inicial. A troca da identificação fiscal de uma conta com conexões entre contas é tratada no requisito da conexão (RF-SRV-014). Uma empresa com operação em dois países abre uma conta por país (um país por empresa no MVP, ADR-013). O roteiro de primeiros passos mostra o progresso e pode ser retomado. Dados de exemplo ficam marcados e são excluídos por inteiro com um comando.

#### RF-ADM-002 · Unidade

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** cadastro da unidade (planta, filial, base ou loja) onde ficam localizações, ativos, almoxarifados e usuários.
- **Dados principais:** nome, código curto, endereço (rua, número, complemento, bairro, cidade, `subdivisao` em ISO 3166-2 e `codigo_postal`, conferidos pelo país da empresa, RNF-044), identificação da filial (opcional, `identificacao_tipo` e `identificacao` de pessoa jurídica do país da empresa: CNPJ, CUIT ou RUC), fuso horário, calendário de trabalho, situação (ativa ou inativa). O calendário fica em duas entidades próprias da unidade: os turnos por dia da semana (`turno_unidade`) e os feriados, com data (`feriado_unidade`) (sumário 3.2).
- **Operações:** criar, editar, inativar; incluir, alterar e remover turnos e feriados.
- **Validações:** código único na empresa; subdivisão da lista do país da empresa e código postal na forma do país (no Brasil, UF e CEP; na Argentina, província e CPA, com a letra inicial igual à da província; no Paraguai, departamento e código de 6 dígitos); fuso da lista do país; calendário com pelo menos um turno (`turno_unidade`). Os turnos são configuráveis pela empresa: o sistema não impõe calendário padrão e só oferece modelos que preenchem o formulário (sumário 3.6).
- **Regras:** no plano Essencial a empresa tem uma unidade. A unidade não é excluída se tiver ativos, OS ou movimentações; só é inativada. O calendário da unidade é a base do tempo requerido nos KPIs quando o ativo não tem calendário próprio. A unidade não tem país próprio no MVP: vale o da empresa (RNF-044); os feriados nacionais vêm do perfil do país.

#### RF-ADM-003 · Multi-unidade

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** várias unidades na mesma empresa, com visão consolidada e acesso restrito por unidade.
- **Dados principais:** lista de unidades; escopo de unidades por usuário.
- **Operações:** criar unidades adicionais; definir a quais unidades cada usuário tem acesso; trocar a unidade ativa na interface; ver painéis consolidados ou por unidade; transferir ativo entre unidades (RF-ATV-011).
- **Validações:** usuário com escopo restrito não lê nem grava dados de outras unidades, inclusive pela API.
- **Regras:** a restrição por unidade é aplicada no servidor, além da restrição por empresa. Ao rebaixar do Profissional para o Essencial, a empresa escolhe a unidade que continua ativa; as demais ficam somente leitura (proposta).

#### RF-ADM-004 · Usuários

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** gestão dos logins de acesso (`usuario`), inclusive solicitantes, motoristas, visualizadores, compradores e locatários gratuitos e usuários de prestadores.
- **Dados principais:** nome, e-mail, telefone celular (opcional), idioma (`usuario.idioma`, RNF-043), pessoa vinculada (`usuario.pessoa_id`, opcional), fornecedor vinculado (`usuario.fornecedor_id`, obrigatório no perfil `prestador`), locatário vinculado (obrigatório no perfil `locatario`), perfis, unidades, custo por hora (para técnicos, preço unitário na moeda da empresa, RNF-045), situação (`usuario.situacao`: `convidado`, `ativo`, `inativo`; sumário 3.6), permissões avulsas (`usuario_permissao`, RF-ADM-005), foto.
- **Operações:** convidar por e-mail; importar por planilha; atribuir perfis; vincular a uma pessoa existente ou criar a pessoa junto; inativar e reativar; redefinir acesso. O convite por link enviado ao telefone depende da pendência **P11**. Conceder e retirar permissões avulsas (RF-ADM-005).
- **Validações:** e-mail obrigatório e único na empresa enquanto a **P11** estiver aberta; pelo menos um perfil; não inativar o último `administrador`; usuário com perfil `prestador` tem `fornecedor_id` de um fornecedor ativo e só tem esse perfil; usuário com perfil `locatario` só tem esse perfil; usuário interno nunca recebe perfil externo (proposta; coerente com o documento 3, seção 4.9, e com o modelo de dados, seção 5.1).
- **Regras:** o convite recebe o idioma padrão da empresa; no primeiro acesso a tela vem no idioma do navegador ou do aparelho, se suportado, e depois só o próprio usuário o muda, no perfil; o usuário externo (`prestador`, `locatario`) altera na própria linha apenas nome, telefone e idioma (RNF-043). Identificação pessoal (CPF, CUIL, DNI, CI), CNH, documentos pessoais e capacitações ficam na `pessoa` (RF-ADM-015), não no usuário. Usuário com histórico (OS, apontamento, assinatura, leitura) nunca é excluído, só inativado, para preservar a trilha. Inativar um usuário encerra suas sessões e marca o cache do app para limpeza na próxima conexão (RF-MOB-010). O login de motoristas e técnicos sem e-mail, por telefone e código de verificação (proposta), porque parte deles não usa e-mail corporativo, fica adiado até a escolha do provedor de SMS (**P11**); até lá, quem não tem e-mail é cadastrado como pessoa sem usuário (RF-ADM-015). Uma mesma pessoa pode ter acesso a mais de uma empresa e escolhe a empresa ao entrar (proposta). O bloqueio temporário por falhas de login é da identidade (`identidade.bloqueada_ate`, RF-ADM-006), não uma situação do usuário (sumário 3.6).

#### RF-ADM-005 · Perfis e permissões

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** controle de acesso por perfil, com matriz de módulo por ação.
- **Dados principais:** os doze perfis padrão do sumário (`administrador`, `gestor`, `planejador`, `tecnico`, `almoxarife`, `comprador`, `plh`, `motorista`, `solicitante`, `visualizador`, `prestador`, `locatario`); matriz de permissões no formato `<módulo>.<ação>` (ler, criar, editar, aprovar, encerrar, excluir, exportar e ações específicas, como `os.encerrar` e `est.decidir_destino`), com o catálogo e a matriz dos perfis padrão na subseção abaixo; permissões avulsas; perfis personalizados; permissões avulsas de um único usuário em `usuario_permissao` (sumário 3.6); **escopo de acesso** do usuário: `interno`, `prestador`, `locatario` ou `link` (sessão curta aberta pelo link de uso único, RF-TER-005) (sumário 3.6).
- **Operações:** consultar a matriz dos perfis padrão e o catálogo de permissões (`GET /v1/permissoes`); criar perfil personalizado a partir de um padrão (Profissional, proposta); atribuir perfis a usuários; conceder e retirar permissões avulsas (`usuario_permissao`), como `os.encerrar` para um usuário que não é `gestor` nem `planejador`.
- **Validações:** perfis padrão não podem ser alterados nem excluídos; permissão avulsa é sempre do catálogo abaixo; ninguém concede nem retira permissão de si mesmo; perfil personalizado ou permissão avulsa só é concedido por quem tem a permissão concedida e a de gerenciar usuários ou perfis, salvo `os.encerrar` e `est.decidir_destino`, que só o `administrador` concede ou retira, sempre de outro usuário, embora não as tenha por padrão (proposta, RN-ADM-026 do documento 4); perfil personalizado derivado de `prestador` ou `locatario` mantém o filtro de escopo do perfil de origem; usuário `prestador` ou `locatario` não recebe permissão avulsa fora do conjunto padrão do seu perfil, nem as permissões vedadas a externos (proposta, RN-ADM-026).
- **Regras:** a permissão é verificada no servidor em cada requisição, inclusive na API e na sincronização do app. Permissões sensíveis (liberar OS, aprovar compra, decidir o destino de peça com defeito, exportar a base inteira) são separadas das de edição comum. A liberação de OS é a permissão `os.encerrar`, padrão em `gestor` e `planejador` e nunca concedida a `prestador` ou `locatario`, nem a perfil personalizado derivado deles. O destino de peça com defeito é a permissão `est.decidir_destino`, padrão em `gestor`, `planejador` e `almoxarife` (RF-EST-017). A Tenaz não guarda dado de saúde; a permissão `pessoa.saude.ler` foi removida do escopo em 2026-11-13 por decisão do dono do produto. Os escopos `prestador`, `locatario` e `link` têm, além da restrição por empresa, um filtro próprio aplicado no servidor (RNF-037 e RNF-038): o `prestador` vê só OS atribuídas ao seu fornecedor, o `locatario` só ativos locados a ele e o `link` só a OS do link. O usuário `interno` não tem filtro adicional além das unidades.

##### Catálogo de permissões e matriz dos perfis padrão (RF-ADM-005)

Esta é a matriz completa de permissões a que o documento 3 remete (3_ESPECIFICACAO_APIS.md, seção 4.8). Respeita o sumário (seções 3.1 e 3.6), os grupos e as doze permissões finas do documento 3 (seções 4.8, 4.9 e o catálogo de endpoints da seção 6) e as permissões citadas no modelo de dados (2_ARQUITETURA_BANCO_DADOS.md, seção 5.1). O que não vem de nenhuma dessas fontes está marcado "(proposta)". O sumário 3.6 (linha "Catálogo de permissões") declara canônicos os códigos deste catálogo e a matriz dos perfis padrão. Permissão nova entra aqui por decisão do dono do produto, com nova versão deste documento, antes de entrar no código.

**Regras do catálogo**

1. **Formato.** `<módulo>.<ação>`, com o código do módulo do sumário (seção 2) em minúsculas, sem exceção (sumário 3.6). A ação `ler` dá leitura. A ação `gerenciar` reúne criar, editar e excluir ou inativar quando os três têm sempre os mesmos perfis padrão, como `prj.gerenciar` e `ctr.gerenciar` no sumário 3.6. Ação de decisão (aprovar, validar, liberar, cancelar, decidir, reabrir) e ação cujos perfis padrão diferem dos do cadastro têm permissão própria (proposta).
2. **O catálogo inteiro existe desde o MVP.** `GET /v1/permissoes` devolve todas as permissões abaixo, e os doze perfis padrão nascem com a empresa (RF-ADM-001, RF-SRV-001) já com a lista completa, inclusive as permissões de recursos das Fases 2 e 3, para que a empresa não precise rever perfis quando o recurso chegar. A coluna **Fase** diz quando a permissão passa a ter efeito; antes disso, a rota não existe (proposta).
3. **Permissão não substitui plano nem alcance.** A ordem é credencial, plano, permissão e alcance (documento 3, seção 4.8). Ter a permissão não libera recurso fora do plano (`403 plano_nao_inclui_recurso`), nem registro fora do alcance do usuário.
4. **Rotas sem permissão fina.** As rotas de sessão e de infraestrutura do próprio usuário exigem só autenticação: `/v1/auth/*`, `/oauth/*`, `GET` e `PATCH /v1/usuarios/eu`, anexos, processamentos, notificações e os próprios dispositivos (documento 3, seção 4.8). `GET /v1/versao` e `GET /v1/configuracao-app` são públicas.
5. **Rota com mais de uma permissão.** Consumir e estornar peça na OS aceita `os.executar` ou `est.movimentar`. Concluir projeto com ativos `em_implantacao` exige `prj.gerenciar` e `atv.colocar_em_operacao`. Os rascunhos de IA exigem `ia.usar` e a permissão do registro (`os.criar`, `sol.criar` ou `chk.gerenciar`).
6. **Sensível.** Marcadas "sim" as permissões que este requisito separa da edição comum (liberar OS, aprovar compra, decidir destino de peça, exportar a base inteira), as ações administrativas sensíveis do documento 3, seção 4.4 (tokens, SSO, perfis e usuários, cobrança e parâmetros, conexões entre contas, aplicação de pacote, regras configuráveis, reabertura de medição) e, por proposta, `con.bloquear_ativo` e `ter.aprovar_orcamento`. A tela de perfis mostra as sensíveis destacadas, e cada concessão ou retirada gera evento de auditoria (RF-ADM-011) (proposta).

**Regras dos perfis padrão**

- **Grupos do documento 3.** A matriz expande os grupos do catálogo de endpoints: Todos (todos os perfis internos, sem `prestador` e `locatario`); Leitura (`administrador`, `gestor`, `planejador`, `tecnico`, `almoxarife`, `comprador`, `plh`, `visualizador`); Gestão (`administrador`, `gestor`, `planejador`); Execução (Gestão e `tecnico`); Estoque (`administrador`, `gestor`, `almoxarife`); Compras (Estoque e `comprador`); Admin (`administrador`); Externos (`prestador` e `locatario`, com o alcance da seção 4.9 do documento 3). Os grupos de uma permissão só (Liberação, Destino de peça, Regras, Configuração, Entrada em operação, Projetos, Dispensa de etapa, Ajuste de medição, Contratos, Fechamento de medição, Reabertura de medição) são as próprias permissões, com os perfis do sumário 3.6.
- **`administrador`.** Tem todas as permissões do catálogo, exceto três: `os.encerrar` e `est.decidir_destino`, decisões da operação que o sumário 3.6 e o documento 3 põem só nos perfis da operação; e `ter.manter_catalogo_proprio`, exclusiva do `prestador`. O administrador que também libera OS ou decide destino de peça recebe o perfil `gestor` ou `planejador` (um usuário pode ter vários perfis, seção 2). Essa regra é decisão do sumário 3.6. Só o administrador concede essas duas permissões internas a outros usuários, por perfil personalizado ou permissão avulsa (`usuario_permissao`), e nunca a si mesmo (proposta, RN-ADM-026).
- **Perfis externos (`prestador`, `locatario`).** Têm só as permissões marcadas na matriz, sempre com o alcance da seção 4.9 do documento 3 (RNF-037, RNF-038), e só nas rotas em que aparecem nomeados no catálogo de endpoints; nas demais, a resposta é `404`. Nunca recebem `os.encerrar`, as nove permissões da configuração, de PRJ e de CTR (sumário 3.6), `srv.gerenciar_conexao` nem permissão do grupo Admin. Perfil personalizado derivado de um deles, e permissão avulsa dada a um deles, só contém permissões do conjunto padrão do perfil de origem (proposta, RN-ADM-026). O `locatario` não tem `mob.usar_app`: o app recusa o login dele (sumário 3.6, RF-MOB-001).
- **Alcance interno.** O alcance não é permissão: o `solicitante`, o `motorista` e o `locatario` veem só as próprias solicitações (proposta do documento 3); o técnico aponta horas só em OS em que está atribuído; o titular vê a própria pessoa e as próprias capacitações; o autor vê a própria execução de checklist e o próprio sinistro.
- **Conta da subcontratada (SRV).** Usa os mesmos doze perfis e a mesma matriz, com o significado do lado de quem presta o serviço (seção 2). A permissão `srv.gerenciar_conexao` (sumário 3.6) vale nos dois lados da conexão e tem como padrão `administrador` e `gestor` (RF-SRV-014, RF-SRV-019); o servidor confere de que lado a conta está: só o cliente convida e altera ativos e campos compartilhados, só a subcontratada aceita ou recusa, e só quem suspendeu reativa. Na conta da subcontratada, `os.encerrar` vale só para OS sem vínculo com conta de cliente; numa OS vinculada, a liberação é sempre do cliente (RF-SRV-018). Nenhuma permissão de uma conta vale na outra: técnicos e gestores da subcontratada não têm credencial na conta do cliente (documento 3, seção 4.7).

**Catálogo** (Fase: MVP, F2 ou F3, quando a permissão passa a ter efeito)

| Permissão | Dá direito a | Grupo no documento 3 | Fase | Sensível |
|---|---|---|---|---|
| **ADM** | | | | |
| `adm.ler` | Ler os dados da empresa e das unidades, com turnos e feriados | Leitura | MVP | não |
| `adm.gerenciar_empresa` | Alterar os dados cadastrais da empresa; criar, alterar e inativar unidades, com turnos (`turno_unidade`) e feriados (`feriado_unidade`) (RF-ADM-001, RF-ADM-002) | Admin | MVP | não |
| `adm.gerenciar_plano` | Consultar plano, recursos e uso; trocar de plano e de ciclo; baixar faturas (RF-ADM-009) | Admin | MVP | sim |
| `adm.ler_pessoas` | Listar e detalhar pessoas (RF-ADM-015) | Gestão | MVP | não |
| `adm.gerenciar_pessoas` | Cadastrar, alterar, inativar e importar pessoas; ligar a pessoa a um usuário (RF-ADM-015) | Gestão | MVP | não |
| `adm.ler_usuarios` | Listar e detalhar usuários e perfis (RF-ADM-004, RF-ADM-005) | Gestão | MVP | não |
| `adm.gerenciar_usuarios` | Criar e convidar usuários; alterar perfis, unidades e dados; desativar e reativar; encerrar sessões e revogar dispositivos de outro usuário (RF-ADM-004, RF-MOB-010) | Admin | MVP | sim |
| `adm.gerenciar_perfis` | Consultar o catálogo de permissões (`GET /v1/permissoes`); criar, alterar e excluir perfis personalizados (Fase 2, Profissional, proposta) (RF-ADM-005) | Admin | MVP | sim |
| `adm.ler_parametros` | Ler os parâmetros da empresa (RF-ADM-010) | Gestão | MVP | não |
| `adm.gerenciar_parametros` | Alterar os parâmetros da empresa (RF-ADM-010) e ligar ou desligar, por tela, a captura de visualização, impressão e extração da trilha (`captura_auditoria`, RF-ADM-011) | Admin | MVP | sim |
| `adm.ler_auditoria` | Consultar, exportar e verificar a integridade da trilha de auditoria (RF-ADM-011) | Admin, `gestor` | MVP | não |
| `adm.gerenciar_sso` | Configurar SSO OIDC (MVP) e SAML, e emitir token SCIM (Fase 2) (RF-ADM-007, RF-ADM-008) | Admin | MVP | sim |
| `adm.exportar` | Exportação completa da empresa, agendamento e pedido de encerramento da conta (RF-ADM-013) | Admin | MVP | sim |
| `adm.ler_configuracao` | Ler e simular regras configuráveis e ler os disparos (RF-ADM-016); ler o catálogo de pacotes de segmento e as aplicações (Fase 2, RF-ADM-017) | Gestão | MVP | não |
| `adm.gerenciar_regras` | Criar, alterar, ativar, inativar e excluir regras configuráveis (RF-ADM-016) | Regras | MVP | sim |
| `adm.gerenciar_configuracao` | Manter as definições de campos extras da OS (RF-OS-021, MVP); pré-visualizar e aplicar pacotes de segmento (RF-ADM-017, Fase 2) | Configuração | MVP | sim |
| **ATV** | | | | |
| `atv.ler` | Ler localizações e ativos, a árvore e o QR code (RF-ATV-001, RF-ATV-002, RF-ATV-007) | Todos, `prestador`, `locatario` | MVP | não |
| `atv.ler_tecnico` | Ler classes de ativo, histórico do ativo e documentos técnicos (RF-ATV-005, RF-ATV-008, RF-ATV-010) | Leitura | MVP | não |
| `atv.gerenciar` | Criar, alterar, mover, inativar e reativar localizações, classes e ativos; emitir etiquetas QR; vincular documentos técnicos | Gestão | MVP | não |
| `atv.colocar_em_operacao` | Executar `colocar_em_operacao`, unitária ou em lote (RF-ATV-013); concluir projeto com ativos `em_implantacao` (RF-PRJ-007, junto com `prj.gerenciar`) | Entrada em operação | MVP | não |
| **MED** | | | | |
| `med.ler` | Ler medidores com a última leitura válida (RF-MED-001) | Leitura, `motorista`, `prestador`, `locatario` | MVP | não |
| `med.ler_leituras` | Ler o histórico de leituras (RF-MED-002) | Leitura, `locatario` | MVP | não |
| `med.registrar_leitura` | Registrar leitura de medidor (RF-MED-002) | Execução, `motorista`, `prestador` | MVP | não |
| `med.gerenciar` | Criar, alterar e inativar medidores; registrar troca ou virada do contador (RF-MED-001, RF-MED-004) | Gestão | MVP | não |
| `med.validar_leitura` | Aceitar ou invalidar leitura `suspeita`, com motivo (RF-MED-003, RF-MED-007) | Gestão | MVP | não |
| **SOL** | | | | |
| `sol.ler` | Ler solicitações e comentários (RF-SOL-003) | Todos, `locatario` | MVP | não |
| `sol.criar` | Abrir solicitação, comentar, alterar a própria enquanto `nova` e desistir da própria (RF-SOL-001) | Todos, `locatario` | MVP | não |
| `sol.triar` | Iniciar triagem, aprovar (cria a OS), recusar, marcar duplicada e reabrir como nova; alterar qualquer solicitação (RF-SOL-002) | Gestão | MVP | não |
| **OS** | | | | |
| `os.ler` | Ler OS, histórico, tarefas, apontamentos, serviços, peças, comentários, campos extras, códigos de falha, registro NR-12, registros de bloqueio e etiquetagem e o relatório em PDF | Leitura, `prestador`, `locatario` | MVP | não |
| `os.criar` | Criar OS (RF-OS-001) | Execução | MVP | não |
| `os.editar` | Alterar os campos editáveis da OS; incluir tarefas e procedimentos; manter os serviços a realizar (RF-OS-006, RF-SRV-007) | Execução | MVP | não |
| `os.executar` | Reconhecer, assumir, iniciar, pausar, retomar e concluir; marcar tarefas; apontar e corrigir horas; consumir e estornar peças; registrar assinaturas e valores de campos extras; registrar e retirar bloqueio e etiquetagem (RF-OS-007 a RF-OS-011, RF-OS-014, RF-CON-004) | Execução, `prestador` | MVP | não |
| `os.comentar` | Comentar a OS, com anexos (RF-OS-016) | Execução, `prestador`, `locatario` | MVP | não |
| `os.planejar` | Planejar, programar e desprogramar; remover tarefa não iniciada; excluir apontamento antes do encerramento (RF-OS-006, RF-OS-007) | Gestão | MVP | não |
| `os.cancelar` | Cancelar OS `aberta`, `planejada`, `programada` ou `pausada`, com motivo (RF-OS-003) | Gestão | MVP | não |
| `os.alterar_responsavel_custo` | Trocar o `responsavel_custo`, com justificativa (RF-OS-013) | Gestão | MVP | não |
| `os.ler_custos` | Ler os custos da OS (RF-OS-013) | Gestão, `locatario` | MVP | não |
| `os.encerrar` | Liberar (`encerrar`) a OS e recusar a liberação (`recusar-liberacao`) (RF-OS-014) | Liberação | MVP | sim |
| `os.gerenciar_codigos_falha` | Criar, alterar e inativar códigos de falha do catálogo personalizável (RF-OS-012) | Gestão | F2 | não |
| **PRE** | | | | |
| `pre.ler` | Ler planos, gatilhos, previsão, vencimentos e procedimentos (RF-PRE-001 a RF-PRE-008) | Leitura | MVP | não |
| `pre.gerenciar` | Criar, alterar, ativar, suspender e excluir planos e gatilhos; manter procedimentos; gerar OS antecipada (RF-PRE-001 a RF-PRE-009) | Gestão | MVP | não |
| **CHK** | | | | |
| `chk.ler` | Ler modelos de checklist e seus itens (RF-CHK-001) | Todos (e `prestador`, 4.9) | MVP | não |
| `chk.ler_execucoes` | Ler execuções de checklist, relatórios e versões de modelo (RF-CHK-006) | Leitura (e `prestador`, 4.9) | MVP | não |
| `chk.executar` | Registrar execução de checklist (RF-CHK-002, RF-CHK-004, RF-CHK-005) | Execução, `motorista` (e `prestador`, 4.9) | MVP | não |
| `chk.gerenciar` | Criar, alterar, publicar e inativar modelos de checklist (RF-CHK-001) | Gestão | MVP | não |
| **EST** | | | | |
| `est.ler` | Ler itens, almoxarifados, saldos, movimentações e reservas; na Fase 2, saldo consolidado, peças serializadas, ordens de reparo, inventários, requisições, pedidos e NF-e | Leitura, `prestador` | MVP | não |
| `est.gerenciar` | Criar, alterar e inativar itens e almoxarifados; na Fase 2, cadastrar peças serializadas e abrir inventário (RF-EST-001, RF-EST-002, RF-EST-005, RF-EST-016) | Estoque | MVP | não |
| `est.movimentar` | Registrar `entrada` e, conforme o plano, `transferencia`, `remessa_terceiro` e `retorno_terceiro`; confirmar recebimento de remessa e entrada de NF-e (Fase 2); consumir e estornar peça na OS (RF-EST-003, RF-EST-015) | Estoque, `prestador` | MVP | não |
| `est.ajustar_saldo` | Registrar `ajuste` com motivo, estornar movimentação e concluir inventário com ajustes (RF-EST-003, RF-EST-005) | Estoque | MVP | não |
| `est.reservar` | Reservar item para OS e cancelar reserva (RF-EST-006) | Execução, Estoque | MVP | não |
| `est.decidir_destino` | Decidir o destino de peça com defeito: abrir ordem de reparo, baixar como sucata (`baixa_sucata`), declarar irreparável, devolver ao fornecedor (RF-EST-017, RF-EST-018) | Destino de peça | F2 | sim |
| `est.contar_inventario` | Registrar contagens de inventário (RF-EST-005) | Estoque, `prestador` | F2 | não |
| `est.instalar_peca` | Instalar e remover peça serializada em ativo, ligada a uma OS (RF-EST-016) | Execução, `prestador` | F2 | não |
| `est.gerenciar_reparo` | Alterar e enviar ordem de reparo; registrar o retorno do conserto (RF-EST-018) | Estoque | F2 | não |
| `est.cancelar_reparo` | Cancelar ordem de reparo `aberta`, com motivo (RF-EST-018) | Estoque, Gestão | F2 | não |
| `est.requisitar` | Criar requisição de material (RF-EST-007) | Execução, Compras | F2 | não |
| `est.aprovar_compra` | Aprovar requisição; aprovar e cancelar pedido de compra (RF-EST-007, RF-EST-008) | Gestão | F2 | sim |
| `est.criar_pedido` | Criar e alterar pedido de compra em rascunho (RF-EST-008) | Compras, Gestão | F2 | não |
| `est.comprar` | Atender requisição; enviar pedido ao fornecedor; importar e conciliar NF-e (RF-EST-007 a RF-EST-009) | Compras | F2 | não |
| **TER** | | | | |
| `ter.ler` | Ler fornecedores e o catálogo de serviços (RF-TER-001, RF-TER-002) | Leitura, `prestador` | MVP | não |
| `ter.gerenciar` | Criar, alterar e inativar fornecedores, papéis e serviços do fornecedor (RF-TER-001, RF-TER-002) | Compras, Gestão | MVP | não |
| `ter.manter_catalogo_proprio` | Criar e alterar os serviços do próprio fornecedor, se a empresa permitir (parâmetro, proposta) (RF-TER-004) | `prestador` (4.9) | MVP | não |
| `ter.atribuir` | Atribuir e retirar a OS de um fornecedor; gerar, listar e revogar links de uso único (RF-TER-003, RF-TER-005) | Gestão | MVP | não |
| `ter.ler_contratos` | Ler contratos de prestação, cumprimento de SLA, avaliações e orçamentos de fornecedor (RF-TER-006, RF-TER-007, RF-TER-009) | Gestão, Compras, `prestador` | F2 | não |
| `ter.gerenciar_contratos` | Criar, alterar e inativar contratos de prestação (RF-TER-006) | Gestão, Compras | F2 | não |
| `ter.avaliar` | Registrar avaliação de fornecedor (RF-TER-007) | Gestão | F2 | não |
| `ter.registrar_orcamento` | Registrar e alterar orçamento de oficina, assistência ou conserto, inclusive o da ordem de reparo (RF-TER-009, RF-EST-018) | Gestão, Compras, `prestador` | F2 | não |
| `ter.aprovar_orcamento` | Aprovar ou recusar orçamento de fornecedor, inclusive o da ordem de reparo (RF-TER-009, RF-EST-018) | Gestão, `comprador` | F2 | sim |
| **SRV** | | | | |
| `srv.ler` | Ler especialidades de técnico e clientes atendidos (RF-SRV-002, RF-SRV-004) | Leitura | MVP | não |
| `srv.gerenciar` | Criar, alterar e inativar especialidades de técnico e clientes atendidos (RF-SRV-002, RF-SRV-004) | Gestão | MVP | não |
| `srv.acompanhar_conexoes` | Ler as conexões entre contas, os ativos compartilhados e os vínculos de OS; reprocessar vínculo em `erro` (RF-SRV-014 a RF-SRV-016) | Gestão | MVP | não |
| `srv.gerenciar_conexao` | Convidar, aceitar, recusar, suspender, reativar e encerrar conexão entre contas; alterar ativos e campos compartilhados (RF-SRV-014, RF-SRV-019) | Conexão | MVP | sim |
| **LOC** | | | | |
| `loc.ler` | Ler locatários, contratos de locação, itens, matriz de responsabilidade e relatório de cobrança (RF-LOC-001 a RF-LOC-004, RF-LOC-009) | Gestão, `locatario` | F2 | não |
| `loc.gerenciar` | Manter locatários e contratos de locação, itens e matriz; `ativar`, `encerrar`, `cancelar`; registrar substituição (RF-LOC-001 a RF-LOC-004, RF-LOC-007) | Gestão | F2 | não |
| `loc.ler_vistorias` | Ler vistorias de locação e o termo em PDF (RF-LOC-005) | Gestão, Execução, `locatario` | F2 | não |
| `loc.vistoriar` | Registrar vistoria de entrega ou de devolução (RF-LOC-005) | Gestão, Execução | F2 | não |
| `loc.assinar_vistoria` | Registrar a assinatura de uma das partes na vistoria (RF-LOC-005) | Gestão, Execução, `locatario` | F2 | não |
| **FRO** | | | | |
| `fro.ler_veiculo` | Ler dados do veículo e, na Fase 2, o mapa de pneus (RF-ATV-003, RF-FRO-003) | Leitura, `motorista` | MVP | não |
| `fro.ler` | Ler pneus, histórico de pneus e abastecimentos (RF-FRO-002 a RF-FRO-007) | Leitura | F2 | não |
| `fro.gerenciar` | Alterar dados do veículo (MVP); na Fase 2, alterar e sucatear pneu, corrigir e excluir abastecimento, ler e manter multas e sinistros | Gestão | MVP | não |
| `fro.acoplar` | Acoplar e desacoplar implemento (RF-ATV-004) | Gestão, `motorista` | MVP | não |
| `fro.gerenciar_pneus` | Cadastrar pneu; enviar e receber recapagem (RF-FRO-002, RF-FRO-005) | Gestão, Estoque | F2 | não |
| `fro.operar_pneus` | Montar e desmontar pneu (RF-FRO-003) | Execução | F2 | não |
| `fro.aferir_pneu` | Registrar sulco e pressão (RF-FRO-004) | Execução, `motorista` | F2 | não |
| `fro.abastecer` | Registrar abastecimento ou Arla (RF-FRO-007) | Execução, `motorista` | F2 | não |
| `fro.registrar_sinistro` | Registrar sinistro (RF-FRO-011) | Gestão, `motorista` | F2 | não |
| **CON** | | | | |
| `con.ler` | Ler tipos de documento, agenda de documentos, inventário e cronograma NR-12, bloqueios de ativo e a consulta de despacho; na Fase 2, NR-13 e pontos de bloqueio | Leitura | MVP | não |
| `con.gerenciar_tipos_documento` | Criar e alterar tipos de documento (RF-CON-005) | Admin | MVP | não |
| `con.gerenciar_documentos` | Registrar, alterar, renovar e inativar documentos com vencimento (RF-CON-005) | Gestão | MVP | não |
| `con.bloquear_ativo` | Criar bloqueio manual de ativo e liberar bloqueio de origem `manual` ou `regra`, com motivo (RF-CON-007, RF-ADM-016) | Gestão | MVP | sim |
| `con.gerenciar_capacitacoes` | Ler e manter capacitações de pessoas (RF-CON-012; no MVP, só nas contas com SRV, RF-SRV-003) | Gestão | F2 | não |
| `con.gerenciar_nr13` | Definir enquadramento e prazos adotados da NR-13 e manter o Registro de Segurança (RF-CON-008 a RF-CON-010) | Gestão (e `plh`, proposta) | F2 | não |
| `con.gerenciar_loto` | Cadastrar, alterar e inativar pontos de bloqueio LOTO (RF-CON-011) | Gestão | F2 | não |
| **IND** | | | | |
| `ind.ler` | Ler indicadores, composição, painéis e o catálogo de relatórios (RF-IND-001 a RF-IND-005) | Leitura | MVP | não |
| `ind.gerenciar_paineis` | Criar, alterar e excluir painéis (RF-IND-001) | Gestão | MVP | não |
| `ind.exportar` | Gerar relatórios e exportar listas em CSV ou XLSX (RF-IND-007) | Leitura | MVP | não |
| **INT** | | | | |
| `int.gerenciar_tokens` | Criar, alterar, rotacionar e revogar tokens de API (RF-INT-001) | Admin | MVP | sim |
| `int.gerenciar_webhooks` | Manter assinaturas de webhook, segredos, entregas e reentregas (RF-INT-002) | Admin | MVP | não |
| `int.importar` | Enviar, validar, confirmar e cancelar importações por planilha (RF-INT-003) | Gestão, Estoque | MVP | não |
| `int.gerenciar_conexoes_telematica` | Criar, alterar e desativar conexões de telemática, com credenciais (RF-INT-005) | Admin | F2 | não |
| `int.operar_telematica` | Ler conexões e eventos de telemática; forçar coleta; casar veículos; manter regras de DTC (RF-INT-004 a RF-INT-006) | Gestão | F2 | não |
| `int.gerenciar_erp` | Manter conectores de ERP, mapeamentos e execuções (RF-INT-008) | Admin | F2 | não |
| **IA** | | | | |
| `ia.usar` | Usar o assistente; sem ela, nenhuma rota de IA responde ao usuário. Os rascunhos exigem também a permissão do registro: `os.criar`, `sol.criar` ou `chk.gerenciar` (RF-IA-001, RF-IA-003, RF-IA-004) | Todos | F2 | não |
| `ia.consultar` | Consultar manuais e histórico e avaliar a resposta (RF-IA-002) | Leitura | F2 | não |
| `ia.gerenciar` | Ligar e desligar a IA por perfil e usuário, retenção e consumo de créditos (RF-IA-001) | Admin | F2 | não |
| `ia.aprovar_acoes` | Aprovar ou rejeitar ações propostas por agentes (RF-IA-006) | Gestão | F3 | não |
| **MOB** | | | | |
| `mob.usar_app` | Entrar no app e sincronizar o escopo offline (RF-MOB-001, RF-MOB-002) | Todos, `prestador` | MVP | não |
| **PRJ** | | | | |
| `prj.ler` | Ler projetos, etapas, acompanhamento e termo de entrega (RF-PRJ-002 a RF-PRJ-005, RF-PRJ-007) | Leitura | F2 | não |
| `prj.ler_medicoes` | Ler medições de etapas e o relatório (RF-PRJ-006) | Gestão, Compras | F2 | não |
| `prj.gerenciar` | Manter modelos de projeto e projetos e ler modelos e lançamentos de custo; `iniciar`, `suspender`, `retomar`, `concluir`, `cancelar`; manter etapas, gerar OS, iniciar e concluir etapas; lançar custo; assinar o termo; gerar, apagar e fechar medição de etapas (RF-PRJ-001 a RF-PRJ-007) | Projetos | F2 | não |
| `prj.dispensar_etapa` | Dispensar etapa `pendente` ou `disponivel`, com motivo (RF-PRJ-003) | Dispensa de etapa | F2 | não |
| `prj.ajustar_medicao` | Ajustar, com motivo, o valor de uma linha da medição de etapas (RF-PRJ-006) | Ajuste de medição | F2 | não |
| **CTR** | | | | |
| `ctr.ler` | Ler contratos com clientes, itens, painel, medições e relatório (RF-CTR-001 a RF-CTR-005) | Gestão | F2 | não |
| `ctr.gerenciar` | Manter contratos com clientes e itens; `ativar`, `encerrar`, `cancelar`; abrir medição (RF-CTR-001, RF-CTR-002, RF-CTR-004) | Contratos | F2 | não |
| `ctr.fechar_medicao` | Decidir a exclusão de paradas com causa externa e fechar a medição (RF-CTR-003, RF-CTR-004) | Fechamento de medição | F2 | não |
| `ctr.reabrir_medicao` | Reabrir medição `fechada`, com motivo (RF-CTR-004) | Reabertura de medição | F2 | sim |

**Matriz dos perfis padrão** (● = a permissão está no perfil padrão). Colunas: Adm `administrador`, Ges `gestor`, Pla `planejador`, Téc `tecnico`, Alm `almoxarife`, Com `comprador`, PLH `plh`, Mot `motorista`, Sol `solicitante`, Vis `visualizador`, Pre `prestador`, Loc `locatario`.

| Permissão | Adm | Ges | Pla | Téc | Alm | Com | PLH | Mot | Sol | Vis | Pre | Loc |
|---|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|:-:|
| **ADM** | | | | | | | | | | | | |
| `adm.ler` | ● | ● | ● | ● | ● | ● | ● |  |  | ● |  |  |
| `adm.gerenciar_empresa` | ● |  |  |  |  |  |  |  |  |  |  |  |
| `adm.gerenciar_plano` | ● |  |  |  |  |  |  |  |  |  |  |  |
| `adm.ler_pessoas` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `adm.gerenciar_pessoas` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `adm.ler_usuarios` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `adm.gerenciar_usuarios` | ● |  |  |  |  |  |  |  |  |  |  |  |
| `adm.gerenciar_perfis` | ● |  |  |  |  |  |  |  |  |  |  |  |
| `adm.ler_parametros` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `adm.gerenciar_parametros` | ● |  |  |  |  |  |  |  |  |  |  |  |
| `adm.ler_auditoria` | ● | ● |  |  |  |  |  |  |  |  |  |  |
| `adm.gerenciar_sso` | ● |  |  |  |  |  |  |  |  |  |  |  |
| `adm.exportar` | ● |  |  |  |  |  |  |  |  |  |  |  |
| `adm.ler_configuracao` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `adm.gerenciar_regras` | ● | ● |  |  |  |  |  |  |  |  |  |  |
| `adm.gerenciar_configuracao` | ● | ● |  |  |  |  |  |  |  |  |  |  |
| **ATV** | | | | | | | | | | | | |
| `atv.ler` | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● |
| `atv.ler_tecnico` | ● | ● | ● | ● | ● | ● | ● |  |  | ● |  |  |
| `atv.gerenciar` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `atv.colocar_em_operacao` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| **MED** | | | | | | | | | | | | |
| `med.ler` | ● | ● | ● | ● | ● | ● | ● | ● |  | ● | ● | ● |
| `med.ler_leituras` | ● | ● | ● | ● | ● | ● | ● |  |  | ● |  | ● |
| `med.registrar_leitura` | ● | ● | ● | ● |  |  |  | ● |  |  | ● |  |
| `med.gerenciar` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `med.validar_leitura` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| **SOL** | | | | | | | | | | | | |
| `sol.ler` | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● |  | ● |
| `sol.criar` | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● |  | ● |
| `sol.triar` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| **OS** | | | | | | | | | | | | |
| `os.ler` | ● | ● | ● | ● | ● | ● | ● |  |  | ● | ● | ● |
| `os.criar` | ● | ● | ● | ● |  |  |  |  |  |  |  |  |
| `os.editar` | ● | ● | ● | ● |  |  |  |  |  |  |  |  |
| `os.executar` | ● | ● | ● | ● |  |  |  |  |  |  | ● |  |
| `os.comentar` | ● | ● | ● | ● |  |  |  |  |  |  | ● | ● |
| `os.planejar` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `os.cancelar` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `os.alterar_responsavel_custo` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `os.ler_custos` | ● | ● | ● |  |  |  |  |  |  |  |  | ● |
| `os.encerrar` |  | ● | ● |  |  |  |  |  |  |  |  |  |
| `os.gerenciar_codigos_falha` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| **PRE** | | | | | | | | | | | | |
| `pre.ler` | ● | ● | ● | ● | ● | ● | ● |  |  | ● |  |  |
| `pre.gerenciar` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| **CHK** | | | | | | | | | | | | |
| `chk.ler` | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● |  |
| `chk.ler_execucoes` | ● | ● | ● | ● | ● | ● | ● |  |  | ● | ● |  |
| `chk.executar` | ● | ● | ● | ● |  |  |  | ● |  |  | ● |  |
| `chk.gerenciar` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| **EST** | | | | | | | | | | | | |
| `est.ler` | ● | ● | ● | ● | ● | ● | ● |  |  | ● | ● |  |
| `est.gerenciar` | ● | ● |  |  | ● |  |  |  |  |  |  |  |
| `est.movimentar` | ● | ● |  |  | ● |  |  |  |  |  | ● |  |
| `est.ajustar_saldo` | ● | ● |  |  | ● |  |  |  |  |  |  |  |
| `est.reservar` | ● | ● | ● | ● | ● |  |  |  |  |  |  |  |
| `est.decidir_destino` |  | ● | ● |  | ● |  |  |  |  |  |  |  |
| `est.contar_inventario` | ● | ● |  |  | ● |  |  |  |  |  | ● |  |
| `est.instalar_peca` | ● | ● | ● | ● |  |  |  |  |  |  | ● |  |
| `est.gerenciar_reparo` | ● | ● |  |  | ● |  |  |  |  |  |  |  |
| `est.cancelar_reparo` | ● | ● | ● |  | ● |  |  |  |  |  |  |  |
| `est.requisitar` | ● | ● | ● | ● | ● | ● |  |  |  |  |  |  |
| `est.aprovar_compra` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `est.criar_pedido` | ● | ● | ● |  | ● | ● |  |  |  |  |  |  |
| `est.comprar` | ● | ● |  |  | ● | ● |  |  |  |  |  |  |
| **TER** | | | | | | | | | | | | |
| `ter.ler` | ● | ● | ● | ● | ● | ● | ● |  |  | ● | ● |  |
| `ter.gerenciar` | ● | ● | ● |  | ● | ● |  |  |  |  |  |  |
| `ter.manter_catalogo_proprio` |  |  |  |  |  |  |  |  |  |  | ● |  |
| `ter.atribuir` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `ter.ler_contratos` | ● | ● | ● |  | ● | ● |  |  |  |  | ● |  |
| `ter.gerenciar_contratos` | ● | ● | ● |  | ● | ● |  |  |  |  |  |  |
| `ter.avaliar` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `ter.registrar_orcamento` | ● | ● | ● |  | ● | ● |  |  |  |  | ● |  |
| `ter.aprovar_orcamento` | ● | ● | ● |  |  | ● |  |  |  |  |  |  |
| **SRV** | | | | | | | | | | | | |
| `srv.ler` | ● | ● | ● | ● | ● | ● | ● |  |  | ● |  |  |
| `srv.gerenciar` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `srv.acompanhar_conexoes` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `srv.gerenciar_conexao` | ● | ● |  |  |  |  |  |  |  |  |  |  |
| **LOC** | | | | | | | | | | | | |
| `loc.ler` | ● | ● | ● |  |  |  |  |  |  |  |  | ● |
| `loc.gerenciar` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `loc.ler_vistorias` | ● | ● | ● | ● |  |  |  |  |  |  |  | ● |
| `loc.vistoriar` | ● | ● | ● | ● |  |  |  |  |  |  |  |  |
| `loc.assinar_vistoria` | ● | ● | ● | ● |  |  |  |  |  |  |  | ● |
| **FRO** | | | | | | | | | | | | |
| `fro.ler_veiculo` | ● | ● | ● | ● | ● | ● | ● | ● |  | ● |  |  |
| `fro.ler` | ● | ● | ● | ● | ● | ● | ● |  |  | ● |  |  |
| `fro.gerenciar` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `fro.acoplar` | ● | ● | ● |  |  |  |  | ● |  |  |  |  |
| `fro.gerenciar_pneus` | ● | ● | ● |  | ● |  |  |  |  |  |  |  |
| `fro.operar_pneus` | ● | ● | ● | ● |  |  |  |  |  |  |  |  |
| `fro.aferir_pneu` | ● | ● | ● | ● |  |  |  | ● |  |  |  |  |
| `fro.abastecer` | ● | ● | ● | ● |  |  |  | ● |  |  |  |  |
| `fro.registrar_sinistro` | ● | ● | ● |  |  |  |  | ● |  |  |  |  |
| **CON** | | | | | | | | | | | | |
| `con.ler` | ● | ● | ● | ● | ● | ● | ● |  |  | ● |  |  |
| `con.gerenciar_tipos_documento` | ● |  |  |  |  |  |  |  |  |  |  |  |
| `con.gerenciar_documentos` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `con.bloquear_ativo` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `con.gerenciar_capacitacoes` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `con.gerenciar_nr13` | ● | ● | ● |  |  |  | ● |  |  |  |  |  |
| `con.gerenciar_loto` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| **IND** | | | | | | | | | | | | |
| `ind.ler` | ● | ● | ● | ● | ● | ● | ● |  |  | ● |  |  |
| `ind.gerenciar_paineis` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `ind.exportar` | ● | ● | ● | ● | ● | ● | ● |  |  | ● |  |  |
| **INT** | | | | | | | | | | | | |
| `int.gerenciar_tokens` | ● |  |  |  |  |  |  |  |  |  |  |  |
| `int.gerenciar_webhooks` | ● |  |  |  |  |  |  |  |  |  |  |  |
| `int.importar` | ● | ● | ● |  | ● |  |  |  |  |  |  |  |
| `int.gerenciar_conexoes_telematica` | ● |  |  |  |  |  |  |  |  |  |  |  |
| `int.operar_telematica` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `int.gerenciar_erp` | ● |  |  |  |  |  |  |  |  |  |  |  |
| **IA** | | | | | | | | | | | | |
| `ia.usar` | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● |  |  |
| `ia.consultar` | ● | ● | ● | ● | ● | ● | ● |  |  | ● |  |  |
| `ia.gerenciar` | ● |  |  |  |  |  |  |  |  |  |  |  |
| `ia.aprovar_acoes` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| **MOB** | | | | | | | | | | | | |
| `mob.usar_app` | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● | ● |  |
| **PRJ** | | | | | | | | | | | | |
| `prj.ler` | ● | ● | ● | ● | ● | ● | ● |  |  | ● |  |  |
| `prj.ler_medicoes` | ● | ● | ● |  | ● | ● |  |  |  |  |  |  |
| `prj.gerenciar` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `prj.dispensar_etapa` | ● | ● |  |  |  |  |  |  |  |  |  |  |
| `prj.ajustar_medicao` | ● | ● |  |  |  |  |  |  |  |  |  |  |
| **CTR** | | | | | | | | | | | | |
| `ctr.ler` | ● | ● | ● |  |  |  |  |  |  |  |  |  |
| `ctr.gerenciar` | ● | ● |  |  |  |  |  |  |  |  |  |  |
| `ctr.fechar_medicao` | ● | ● |  |  |  |  |  |  |  |  |  |  |
| `ctr.reabrir_medicao` | ● | ● |  |  |  |  |  |  |  |  |  |  |

**Totais por perfil padrão:** `administrador` 113, `gestor` 102, `planejador` 87, `almoxarife` 41, `tecnico` 38, `comprador` 32, `plh` 24, `visualizador` 23, `prestador` 18, `motorista` 14, `locatario` 11, `solicitante` 6, num catálogo de 116 permissões. A permissão `pessoa.saude.ler` saiu do catálogo e da matriz em 2026-11-13 (a Tenaz não guarda dado de saúde, sumário 3.6); o código não é reaproveitado.

#### RF-ADM-006 · Autenticação

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** acesso por e-mail e senha, com segundo fator. O acesso por telefone e código fica adiado até a escolha do provedor de SMS (pendência **P11**); até lá, todo login é por e-mail e senha ou por SSO (RF-ADM-007).
- **Dados principais:** credenciais (e-mail e senha), fator TOTP, sessões ativas por dispositivo e por canal (web ou app), tentativas de login com falha e bloqueio temporário, links enviados por e-mail (convite, redefinição de senha, verificação de e-mail), enviados no idioma da requisição que os disparou (redefinição e verificação) ou no idioma do usuário convidado (convite) (RNF-043).
- **Operações:** entrar, sair, recuperar senha por link enviado ao e-mail, ativar MFA por TOTP, tornar MFA obrigatório para a empresa (`empresa.mfa_obrigatorio`), encerrar sessões.
- **Validações:** política de senha (mínimo de 10 caracteres e verificação contra senhas vazadas) (proposta); bloqueio temporário após 5 tentativas erradas (proposta).
- **Regras:** MFA é opcional para clientes e pode ser exigido pelo `administrador` no campo `empresa.mfa_obrigatorio`: ligado, o usuário sem MFA só entra depois de ativá-lo (sumário 3.6). O bloqueio temporário após falhas de login fica na identidade (`identidade.bloqueada_ate`) e não muda a situação do usuário. Para a equipe interna da Tenaz é obrigatório (RNF-017). Todo login, falha de login e troca de senha gera evento de auditoria. Quando a **P11** for decidida, o login por telefone e código entra neste requisito, sem mudar a fase.

#### RF-ADM-007 · SSO OIDC

**Fase:** MVP · **Plano mínimo:** Profissional

- **Descrição:** login único com provedores OIDC, no mínimo Google e Microsoft.
- **Dados principais:** provedor, domínios de e-mail autorizados, perfil padrão de novos usuários.
- **Operações:** configurar provedor; exigir SSO para domínios da empresa; criar usuário no primeiro login (just-in-time) com perfil padrão.
- **Validações:** domínio verificado pela empresa antes de exigir SSO (proposta).
- **Regras:** com SSO exigido, o login por senha fica bloqueado para os domínios cobertos, exceto para um `administrador` de contingência.

#### RF-ADM-008 · SAML e SCIM

**Fase:** Fase 2 · **Plano mínimo:** Enterprise

- **Descrição:** login por SAML 2.0 e provisionamento automático de usuários e grupos por SCIM 2.0.
- **Dados principais:** metadados do IdP, mapeamento de grupos para perfis e unidades, token SCIM.
- **Operações:** configurar IdP; criar, atualizar e desativar usuários pelo SCIM; mapear grupos.
- **Validações:** certificado do IdP válido; mapeamento de grupo para perfil existente.
- **Regras:** desativação via SCIM tem o mesmo efeito de RF-ADM-004 (sessões encerradas, cache do app limpo). Usuários geridos por SCIM não podem ter perfil alterado pela interface.

#### RF-ADM-009 · Plano contratado, limites e cobrança

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** controle do plano, dos recursos liberados e da base de cobrança.
- **Dados principais:** plano (Essencial, Profissional, Enterprise e, para contas com o módulo SRV, Serviços), ciclo (mensal ou anual), moeda de cobrança (a da empresa, `BRL`, `ARS` ou `PYG`; preço, moeda e meio de pagamento na Argentina e no Paraguai dependem da **P15**), número de técnicos e de veículos cobrados, créditos de IA, faturas, forma de pagamento (tratada pelo provedor de pagamento, nunca armazenada pela Tenaz).
- **Operações:** consultar uso; trocar de plano; mudar o ciclo; baixar faturas; comprar créditos de IA (Fase 2).
- **Validações:** downgrade bloqueado enquanto houver recurso em uso que o plano novo não inclui, com lista do que precisa ser desligado (proposta).
- **Regras:** cobra-se por **técnico** e por **veículo**, conforme o sumário (seção 2.1):
  - **Técnico cobrado:** usuário ativo com o perfil `tecnico`. `gestor` ou `planejador` que também executa OS recebe o perfil `tecnico` e passa a contar. Usuário com perfil `prestador` que executa OS também conta como técnico. A execução registrada por link de uso único (RF-TER-005) não conta. O técnico de uma subcontratada com conta própria, que recebe a OS por conexão entre contas (RF-SRV-016), não conta na conta do cliente: é cobrado na conta da subcontratada, pelo plano Serviços (RF-SRV-001).
  - **Veículo cobrado:** ativo do tipo `veiculo` com status diferente de `inativo`, inclusive o tomado em locação (`propriedade = de_terceiro`).
  - A base de cobrança do ciclo é o **maior número diário** de técnicos e de veículos no ciclo.
  - Solicitantes, motoristas, visualizadores, compradores, locatários e ativos de equipamento são gratuitos e ilimitados.
  - Mínimos: Essencial 2 técnicos; Profissional 3 técnicos; Enterprise 10 técnicos ou 200 veículos. Os preços e mínimos são hipótese a validar em **P4** e, fora do Brasil, em **P15**.
  - **Plano Serviços** (sumário 2.2): cobrado por técnico ativo com perfil `tecnico` na conta da subcontratada (R$ 59 anual · R$ 69 mensal, proposta, valores em reais, só para empresas do Brasil; a tabela para `ARS` e `PYG` é a **P15**); `gestor` e `planejador` não pagam. A conta pode contratar também Essencial, Profissional ou Enterprise para frota ou equipamentos próprios; nesse caso, um mesmo técnico é cobrado uma vez só, pelo plano de maior preço (proposta). O desconto de canal para subcontratada conectada a cliente Tenaz pagante (ex.: 30%) depende de **P4** (RF-SRV-001).
  - Recursos fora do plano aparecem na interface com indicação do plano que os libera, sem ocultar dados já gravados. Recursos das Fases 2 e 3 só são vendidos quando estiverem prontos.

#### RF-ADM-010 · Parâmetros da empresa

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** valores que alteram regras e cálculos sem mudança de código.
- **Dados principais (exemplos):** antecedência padrão de alertas de documento (60, 30 e 7 dias); janela de tolerância padrão da preventiva; MTTR com ou sem espera de peça; prazo para considerar retrabalho (padrão de 30 dias, proposta); dias sem leitura que disparam alerta de medidor; obrigatoriedade de foto na conclusão de OS; obrigatoriedade de registro de falha em corretiva com parada; prioridade padrão de OS por criticidade; faturamento bruto anual (para o CMF, na moeda da empresa, RNF-045); idioma padrão da empresa (RNF-043). A moeda e o país não são parâmetros: vêm da criação da conta e não mudam (RNF-044, RNF-045); os formatos de número, data e moeda vêm do perfil do país (RNF-029).
- **Operações:** consultar e editar com histórico de alterações.
- **Validações:** cada parâmetro tem tipo, faixa e valor padrão.
- **Regras:** mudança de parâmetro de cálculo vale a partir da data da alteração e fica registrada; os KPIs mostram qual valor estava vigente no período (RF-IND-002).

#### RF-ADM-011 · Trilha de auditoria

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** registro imutável e verificável de quem fez o quê, quando, de onde e em qual tela, com valores antes e depois (`evento_auditoria`), captura configurável por tela da visualização, da impressão e da extração (`captura_auditoria`) e controle de integridade por hash e selo encadeado (`selo_auditoria`) (sumário 3.2 e 3.6). Ampliado na etapa E01 por pedido do dono do produto. Continua um requisito só: evento, captura, selo e telas são partes da mesma trilha e são entregues juntos.
- **Dados de todo evento:**
  - data e hora do servidor e, na ação feita offline, também a do aparelho;
  - usuário, com id e nome no instante do evento, guardados no próprio evento para que a troca de nome não reescreva o histórico (proposta); sem usuário, a integração (token de API) ou o sistema;
  - IP de origem visto pelo servidor; atrás do balanceador, o informado por ele, nunca um cabeçalho enviado pelo cliente (proposta);
  - tipo de operação (tabela abaixo);
  - tela de origem (`evento_auditoria.tela`) e módulo (`evento_auditoria.modulo`, código do sumário, seção 2, em minúsculas); chamada sem tela (API, job, sincronização) grava só o módulo;
  - escopo de acesso da credencial (`interno`, `prestador`, `locatario`, `link`; RF-ADM-005), origem (web, app, API, sistema) e aparelho;
  - entidade e id do registro e, em inclusão, alteração, exclusão e deleção, os valores antes e depois;
  - hash do próprio conteúdo (`evento_auditoria.hash`).
- **Tipos de operação** (lista fechada, sumário 3.6):

  | Tipo | O que é | Gravação |
  |---|---|---|
  | `login` | Entrada na conta, com sucesso ou falha | Sempre |
  | `logout` | Saída, inclusive por expiração ou revogação da sessão | Sempre |
  | `inclusao` | Cadastro de registro | Sempre |
  | `alteracao` | Mudança de campos de um registro existente | Sempre |
  | `exclusao` | Exclusão lógica ou inativação; o registro continua no banco | Sempre |
  | `delecao` | Remoção física de linha, como desfazer um vínculo | Sempre |
  | `visualizacao` | Abrir o detalhe de um registro numa tela | Sempre quando mostra dado sensível; nas demais telas, só com a captura ligada |
  | `impressao` | Gerar PDF ou imprimir um registro ou uma tela | Com a captura ligada na tela |
  | `extracao` | Exportar lista, gerar relatório, baixar exportação | Com a captura ligada na tela |

- **Dado sensível:** identificação pessoal completa, de qualquer tipo (CPF, CUIL, DNI, CI), dados de vítimas de sinistro e identificação de locatário pessoa física (RNF-021). O idioma nunca entra no conteúdo com hash do evento: a mesma trilha é verificada (`integro`, `violado`) em qualquer idioma; a extração e a impressão saem no idioma de quem pede (RNF-043). Dados de saúde foram removidos do escopo em 2026-11-13 por decisão do dono do produto: a Tenaz não os guarda. A visualização deles é sempre gravada, em qualquer tela, sem opção de desligar.
- **Operações:**
  - **Consulta da trilha:** filtros por usuário, tela ou módulo, tipo de operação, período e IP, além de entidade e registro; histórico de um registro na tela dele; detalhe do evento com antes e depois.
  - **Exportação** da consulta em CSV, com os mesmos filtros; a própria exportação é uma `extracao`.
  - **Verificação de integridade:** o usuário escolhe um período, e a tela mostra cada período selado como **íntegro** ou **violado**, apontando o primeiro evento ou selo divergente; o período ainda não selado aparece como **ainda não selado** (proposta).
  - **Configuração das capturas por tela:** lista das telas do sistema por módulo, com uma chave para visualização, impressão e extração em cada tela; as operações sempre gravadas aparecem fixas, sem chave.
- **Validações:** nenhum perfil, inclusive `administrador` e o suporte da Tenaz, altera ou exclui eventos ou selos; tipo de operação só da lista fechada; a configuração só aceita visualização, impressão e extração; consultar, exportar e verificar exige `adm.ler_auditoria`; configurar as capturas exige `adm.gerenciar_parametros` (RF-ADM-005).
- **Regras:**
  - Cobre no mínimo OS, ativos, medidores, estoque, documentos, conformidade, empresa e unidades (inclusive turnos e feriados), pessoas, usuários, perfis e permissões (inclusive as permissões avulsas em `usuario_permissao`), login e falha de login (RF-ADM-006), parâmetros, configurações de integração e conexões entre contas.
  - Os eventos de inclusão, alteração, exclusão e deleção são gravados na mesma transação do fato: se o fato é desfeito, o evento também é.
  - Cada gravação feita pelo processo de sincronização entre contas gera evento nas duas contas, com origem `sistema`, a conexão e o registro correspondente na outra conta (RF-SRV-014, RNF-041).
  - Ligar ou desligar uma captura é uma `alteracao` sempre gravada e só vale para eventos futuros. Padrão na criação da empresa: impressão e extração ligadas em todas as telas, visualização desligada (proposta).
  - **Integridade:** o hash de cada evento cobre todos os campos dele. Um job encadeia, por empresa e por período, os hashes dos eventos do período, em ordem, com o selo do período anterior, e grava o resultado em `selo_auditoria` e fora do banco, em armazenamento imutável. Período do selo: um dia UTC (proposta, como no modelo de dados, seção 9.4). Armazenamento imutável: em desenvolvimento, o SeaweedFS do compose local, bucket `tenaz-auditoria` com bloqueio de objeto em modo de conformidade; em produção, o S3 da AWS com Object Lock em modo de conformidade (ADR-012 e ADR-006 no documento 5; sumário 3.6, "Ferramentas e recursos"), com retenção igual à da trilha (RNF-019). A verificação recalcula os hashes a partir dos eventos e compara com o selo do banco e com a cópia imutável; evento alterado, apagado ou incluído fora de ordem dá **violado**.
  - **Dados pessoais e LGPD:** a visualização, a impressão e a extração guardam o que foi acessado (entidade, registro, tela, filtros), nunca o conteúdo mostrado. A trilha é dado pessoal mantido para cumprir obrigação legal e defender direitos; pedido de eliminação do titular não apaga eventos dentro da retenção (RNF-019, RNF-021).
  - **Retenção:** mínimo de 5 anos, configurável para mais por empresa, para eventos e selos (RNF-019). Depois de 24 meses os eventos saem do banco para o arquivo imutável (modelo de dados, seção 9.5), e a verificação de um período arquivado roda de forma assíncrona (proposta).
  - **Exportação completa (RF-ADM-013):** o pacote inclui os eventos e os selos do período retido, para que o cliente verifique a integridade fora da Tenaz (proposta). Pedir e baixar a exportação completa é sempre gravado como `extracao`, mesmo com a captura desligada (proposta).
- **Critérios de aceite:**
  1. Para cada tipo de operação, o evento traz data e hora, usuário (id e nome), IP, tipo, tela ou módulo e hash; um teste automatizado por tipo.
  2. Com todas as capturas desligadas, login, logout, inclusão, alteração, exclusão, deleção e visualização de dado sensível continuam gerando evento, e a tela de configuração não oferece chave para eles.
  3. Com a captura de visualização desligada numa tela, abrir um registro nela não gera evento; ligada, gera um evento por registro aberto; a troca da chave aparece na trilha como `alteracao`.
  4. Impressão e extração seguem a chave da tela, inclusive na exportação da própria trilha.
  5. Alterar ou apagar um evento direto no banco, com o papel dono das tabelas, faz a verificação do período mostrar **violado** e apontar o evento; período intacto mostra **íntegro**. O selo existe também no armazenamento imutável, que recusa apagá-lo ou sobrescrevê-lo.
  6. Os filtros por usuário, tela ou módulo, tipo, período e IP devolvem só os eventos que atendem a todos, e a exportação em CSV traz as mesmas linhas.
  7. Nenhum perfil altera nem exclui evento ou selo pela API; usuário sem `adm.ler_auditoria` recebe `403`; sem `adm.gerenciar_parametros`, a configuração das capturas também.
  8. O evento de visualização da identificação pessoal completa de uma pessoa (por exemplo, o CPF) não contém o número.
  9. A trilha de uma empresa não mostra eventos de outra (RNF-013), e a verificação de integridade usa só os selos da própria empresa.

#### RF-ADM-012 · Notificações

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** central de notificações na web e no app, com envio por push e e-mail.
- **Dados principais:** evento, destinatário, canal, situação (enviada, lida), preferências por usuário; texto no idioma do destinatário (`usuario.idioma`), ou no idioma padrão da empresa quando o destinatário não tem usuário (RNF-043).
- **Operações:** ver e marcar notificações; configurar canais por tipo de evento; receber resumo diário por e-mail.
- **Validações:** notificações obrigatórias (OS emergencial atribuída, documento crítico vencido) não podem ser desligadas pelo usuário (proposta).
- **Regras:** eventos mínimos: OS atribuída, OS emergencial aberta, mudança de estado da solicitação para o solicitante, preventiva gerada, documento e capacitação a vencer, leitura suspeita, saldo abaixo do mínimo; nas contas conectadas (SRV), convite de conexão recebido, OS recebida de cliente, OS aberta pela subcontratada aguardando triagem no cliente, liberação recusada pelo cliente e erro de sincronização. Outros canais, como WhatsApp, ficam para depois da validação de **P4** (proposta).

#### RF-ADM-013 · Exportação completa e encerramento da conta

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** exportação de todos os dados da empresa, sob demanda ou agendada, e processo de saída do cliente.
- **Dados principais:** pacote em JSON e CSV com anexos, manifesto com contagens e hashes, no idioma de quem pede (rótulos e cabeçalhos; o conteúdo digitado pela empresa não é traduzido, RNF-043); inclui a trilha de auditoria e os selos de integridade do período retido, para verificação fora da Tenaz (proposta, RF-ADM-011).
- **Operações:** solicitar exportação; agendar exportação periódica; baixar pacote por link temporário; solicitar encerramento da conta.
- **Validações:** só `administrador` solicita; link expira em 7 dias (proposta); pedir e baixar o pacote sempre gera evento `extracao` na trilha (proposta, RF-ADM-011).
- **Regras:** o pacote fica pronto em até 24 horas (RNF-023). Após o encerramento, os dados ficam disponíveis para exportação por 90 dias e depois são excluídos, com comprovante (RNF-024). O encerramento da conta encerra todas as suas conexões com outras contas (RF-SRV-019); a outra conta mantém as cópias que já recebeu.

#### RF-ADM-014 · Banco dedicado opcional

**Fase:** Fase 3 · **Plano mínimo:** Enterprise

- **Descrição:** opção de hospedar os dados da empresa num banco exclusivo (modelo silo), em vez do banco compartilhado com RLS.
- **Operações:** provisionar banco dedicado; migrar a empresa do banco compartilhado sem perda de dados; operar com as mesmas versões de schema.
- **Regras:** a aplicação é a mesma para as duas modalidades; só muda o roteamento de conexão. A migração tem janela agendada e relatório de conferência de contagens.

#### RF-ADM-015 · Pessoas

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** cadastro de quem executa, conduz ou responde por algo (`pessoa`): técnico, motorista, condutor eventual, técnico de fornecedor, PLH e, na conta de uma subcontratada, os técnicos dela (RF-SRV-003). A pessoa pode ter ou não um `usuario`.
- **Dados principais:** nome, identificação pessoal (`identificacao_tipo` e `identificacao`, protegida: `cpf` no Brasil, `cuil` ou `dni` na Argentina, `ci` no Paraguai; a pessoa aceita os tipos de qualquer país suportado, porque uma transportadora brasileira pode ter motorista paraguaio, RNF-044), matrícula, vínculo (`pessoa.vinculo`: `empregado`, `terceiro` ou `autonomo`), fornecedor (quando for técnico de fornecedor), unidade, contato, foto, usuário vinculado (opcional), situação (`pessoa.situacao`: `ativa` ou `inativa`; não há `desligada`) (sumário 3.6). Ligados à pessoa: CNH e demais documentos pessoais (RF-CON-005), capacitações (RF-CON-012) e vínculos com veículo (RF-FRO-001). Dados de saúde (ASO, exame toxicológico) foram removidos do escopo em 2026-11-13 por decisão do dono do produto: a Tenaz não os guarda.
- **Operações:** criar, editar, inativar, importar por planilha; criar o usuário a partir da pessoa ou ligar a um usuário existente; consultar documentos e capacitações da pessoa. Rota `/v1/pessoas`.
- **Validações:** identificação, quando informada, com o tipo e o valor juntos, na forma do tipo e com o dígito verificador conferido quando o tipo tem (CPF; CUIL, DNI e CI só a forma, até a **P13**), e única na empresa por `(identificacao_tipo, identificacao)`; uma pessoa tem no máximo um usuário por empresa; pessoa com `fornecedor` informado tem vínculo `terceiro` e só pode ser ligada a usuário com perfil `prestador` daquele fornecedor.
- **Regras:** pessoa sem login não é cobrada. A Tenaz não guarda dado de saúde da pessoa (sumário 3.6). Pessoa com histórico (condução, capacitação, assinatura, execução) não é excluída, só inativada. A pessoa não é modelada como usuário sem login. Vínculo `autonomo` é o profissional sem empregador nem fornecedor, como o condutor eventual contratado por serviço (sumário 3.6). Pessoa inativa nunca é excluída: continua no banco, ligada ao histórico dos serviços de que participou (por exemplo, o motorista temporário de um transporte de equipamento), e volta por reativação, não por novo cadastro. Enquanto a **P11** estiver aberta, técnico ou motorista sem e-mail fica como pessoa sem usuário e é ligado a um usuário quando passar a ter e-mail ou quando o login por telefone existir (RF-ADM-004).
- **Aviso de dado de saúde** (sumário 3.6, "Dado de saúde recebido por engano"): o cadastro e o detalhe da pessoa mostram, sempre, junto aos documentos e anexos dela, o aviso fixo "A Tenaz não guarda dado de saúde. Não registre nem anexe ASO, exames, atestados ou qualquer informação médica." (texto aprovado pelo dono do produto). O aviso não pode ser fechado nem desligado por usuário ou por parâmetro.
- **Critério de aceite do aviso:** nas telas de criar, editar e detalhar pessoa, na web e no app, o aviso aparece em toda carga da tela, para qualquer perfil, sem botão de fechar, e passa na verificação de acessibilidade (RNF-025).

#### RF-ADM-016 · Regras configuráveis de condição e efeito

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** regras da empresa (`regra_automacao`) no formato **condição → efeito**, com condição e efeito escolhidos de catálogos fechados (sumário 3.7). Permitem que cada segmento ligue o comportamento de que precisa sem código novo: por exemplo, "verificação metrológica vencida bloqueia o radar até a nova verificação".
- **Dados principais:** nome, descrição, condição (`documento_vencido`, `documento_a_vencer`, `os_concluida_com_campo`, `item_checklist_reprovado`, `leitura_fora_faixa`; na Fase 2, `capacitacao_vencida` e `etapa_projeto_atrasada`), parâmetros da condição (tipo de documento, dias de antecedência, campo extra sim/não da OS, item de checklist, medidor), filtro (classes de ativo, unidades, criticidades), efeito (`notificar`, `abrir_solicitacao`, `bloquear_ativo`, `bloquear_ativo_ate_documento`; na Fase 2, `impedir_atribuicao`), parâmetros do efeito (destinatários e canal, prioridade da solicitação, tipo de documento que levanta o bloqueio, motivo do bloqueio), situação (`ativa`, `inativa`), quem criou e alterou. Cada disparo gera um `disparo_regra` (regra, fato que disparou, ativo, efeito aplicado, data e hora). Rota `/v1/regras-automacao`.
- **Operações:** criar a partir de modelos (incluídos nos pacotes de segmento, RF-ADM-017); simular sem gravar, mostrando quais ativos seriam afetados hoje; ativar e inativar; consultar o histórico de disparos por regra e por ativo; exportar.
- **Validações:** condição e efeito do catálogo; combinação permitida (ex.: `bloquear_ativo_ate_documento` exige o tipo de documento que levanta o bloqueio; `abrir_solicitacao` não é permitido com `documento_a_vencer` para não abrir solicitação repetida a cada dia) (proposta); limite de 50 regras ativas por empresa (proposta); criar, alterar e ativar exigem a permissão `adm.gerenciar_regras`, padrão em `administrador` e `gestor`.
- **Regras:**
  - A regra é avaliada no fato que a dispara (documento vence, OS é concluída, checklist é enviado, leitura é registrada) e numa rotina diária idempotente, que não repete o efeito para o mesmo fato.
  - `bloquear_ativo` cria um `bloqueio_ativo` com `origem = regra`, que deixa o ativo `parado` (RF-ATV-006) e só é levantado por quem tem a permissão de bloqueio, com motivo.
  - `bloquear_ativo_ate_documento` cria o mesmo bloqueio e o levanta sozinho quando um documento do tipo exigido é registrado no ativo com emissão posterior ao fato que disparou a regra (sumário 3.6).
  - `abrir_solicitacao` cria solicitação com origem `regra` (RF-SOL-001), que segue a triagem normal.
  - Uma regra **nunca** muda o estado de uma OS, nunca altera fórmula de indicador e nunca grava dado em outra conta; o efeito vale só na conta onde a regra existe (o status do ativo não é campo compartilhado na conexão entre contas). Na OS executada por uma subcontratada conectada, os campos extras não trafegam (sumário 3.6), então a condição `os_concluida_com_campo` não dispara com dados da subcontratada; a solução fica em aberto (seção 8.3, ponto 38).
  - Os comportamentos padrão da Tenaz (alertas de RF-CON-006, item reprovado de checklist gerando solicitação em RF-CHK) continuam valendo sem regra; a regra acrescenta efeitos, não desliga os padrões.
  - Inativar a regra não levanta bloqueios já criados por ela; o sistema lista esses bloqueios para decisão.
  - Condição ou efeito novo entra por requisito neste documento, com teste, nunca por configuração (risco R7).

#### RF-ADM-017 · Pacotes de segmento

**Fase:** Fase 2 · **Plano mínimo:** Essencial

- **Descrição:** conjunto pronto de configuração por segmento (`pacote_segmento`), mantido pela Tenaz e aplicado à conta na criação ou depois, para que a empresa comece com a operação quase montada.
- **Dados principais:** pacote (código, nome, segmento, versão, descrição), conteúdo: classes de ativo com campos extras e medidores, campos extras da OS (RF-OS-021), procedimentos, planos padrão (RF-PRE-010), modelos de checklist, tipos de documento, regras configuráveis (RF-ADM-016), modelos de projeto (RF-PRJ-001), códigos de falha sugeridos e módulos recomendados. A aplicação gera `aplicacao_pacote` (pacote, versão, data, quem aplicou, itens criados e itens ignorados por já existirem).
- **Operações:** escolher um ou mais pacotes no roteiro de primeiros passos (RF-ADM-001); aplicar depois, pela administração; ver antes de aplicar o que será criado; aplicar uma versão nova do pacote, que só cria o que falta.
- **Pacotes iniciais (proposta):** `fiscalizacao_eletronica`, `transportadora`, `industria`, `comercio`, `locadora`, `prestador_servicos`. O `fiscalizacao_eletronica` traz: classes "Radar fixo", "Câmera de fiscalização", "Gabinete", "Poste", "Laço indutivo", "Nobreak" e "Link de dados"; campos extras da OS "Rompeu lacre?", "Lacres retirados" e "Lacres instalados"; tipo de documento "Verificação metrológica", crítico; as regras "verificação vencida → `bloquear_ativo_ate_documento`" e "OS concluída com 'Rompeu lacre?' = sim → `bloquear_ativo_ate_documento`", entregues inativas até **P9**; modelo de projeto "Implantação de ponto de fiscalização" (RF-PRJ-001); planos de limpeza de lentes, poda de vegetação, teste de laços e verificação do nobreak.
- **Validações:** pacote publicado pela Tenaz e compatível com o plano da conta (itens de módulos não contratados não são criados e aparecem como ignorados).
- **Regras:** a aplicação é atômica: cria tudo o que falta ou nada. Aplicar exige `adm.gerenciar_configuracao`. A escolha de pacote no roteiro de primeiros passos aparece a partir da Fase 2; no MVP, o roteiro segue sem ela. Aplicar um pacote **só cria o que falta**; nunca altera nem apaga o que a empresa já tem ou editou, identificado pela chave natural do item (sigla da classe, código do tipo de documento, nome da regra). Tudo o que o pacote cria é configuração comum da empresa, editável como se ela tivesse criado. Regras vindas de pacote nascem `inativa` e a empresa as ativa depois de revisar (proposta). O catálogo de pacotes é global, só leitura para as contas e fica fora do RLS por empresa, como os demais catálogos globais (RNF-013).

---

### 4.2 ATV — Ativos

#### RF-ATV-001 · Localizações

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** árvore de lugares físicos dentro da unidade (prédio, linha, setor, doca, pátio).
- **Dados principais:** nome, código, localização pai, unidade, descrição, coordenadas (opcional), QR code.
- **Operações:** criar, editar, mover na árvore, inativar; ver ativos por localização.
- **Validações:** código único na unidade; profundidade máxima de 6 níveis (sumário 3.1); sem ciclos.
- **Regras:** localização com ativos ou OS não é excluída, só inativada. Uma solicitação ou OS pode apontar para uma localização quando não há ativo específico (manutenção predial).

#### RF-ATV-002 · Cadastro de ativo

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** cadastro único para qualquer bem que recebe manutenção, com tipo `equipamento`, `veiculo`, `implemento` ou `componente`.
- **Dados principais:** TAG, nome, tipo, classe, criticidade (A, B ou C), status operacional, localização, ativo pai, fabricante, modelo, número de série, ano de fabricação, data e valor de aquisição, **valor de reposição** (`valor_reposicao`, RAV), **regime de operação** (`horas_operacao_dia`, `dias_operacao_semana`), **propriedade** (`proprio`, `de_terceiro` ou `de_cliente`), cliente atendido (`cliente_servico_id`, só em `de_cliente`), **fase de ciclo de vida** (`fase_ciclo_vida`: `em_implantacao` ou `em_operacao`, padrão `em_operacao`) e data de entrada em operação (`data_entrada_operacao`) (RF-ATV-013), garantia até, centro de custo, responsável, fotos, campos extras da classe, marcação "sujeito à NR-12", marcação "sujeito à NR-13".
- **Operações:** criar, editar, duplicar, mover na árvore, inativar, consultar a árvore e o histórico.
- **Validações:** TAG única na empresa; ativo pai na mesma empresa; profundidade máxima de 6 níveis; sem ciclos; `componente` sempre tem ativo pai (proposta); marcação NR-12 só em `equipamento` ou `implemento` (proposta); RAV maior que zero quando informado; `horas_operacao_dia` entre 0 e 24 e `dias_operacao_semana` entre 0 e 7.
- **Regras:** ativo com histórico nunca é excluído, só inativado. A sugestão de TAG segue `<SIGLA>-<NNNN>` da classe. Os campos obrigatórios extras vêm da classe (RF-ATV-005). O regime de operação é a base do tempo requerido nos KPIs quando informado; sem ele, vale o calendário da unidade. Ativo `de_terceiro` é o tomado em locação ou em comodato; quando tomado, liga-se a um contrato de locação `tomado` (RF-LOC-002). Ativo `de_cliente` só existe na conta de uma subcontratada com o módulo SRV: é o equipamento de um cliente atendido, com `cliente_servico_id` obrigatório; quando vem de um cliente conectado, é cópia sincronizada e guarda `conexao_empresa_id` e `ativo_remoto_id` (RF-SRV-005, RF-SRV-015). Pneu não é ativo: é a entidade própria `pneu` (RF-FRO-002).

#### RF-ATV-003 · Dados de veículo

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** campos específicos para ativos do tipo `veiculo` e `implemento` rodoviário.
- **Dados principais:** placa, RENAVAM, chassi/VIN, ano de fabricação e ano modelo, combustível principal, consome Arla 32 (sim/não), capacidade do tanque, PBT ou CMT, configuração de eixos (usada no mapa de pneus), categoria de CNH exigida, RNTRC vinculado, transporta produto perigoso (sim/não), possui cronotacógrafo (sim/não), número do cronotacógrafo.
- **Operações:** preencher no cadastro do ativo; buscar por placa, RENAVAM ou VIN.
- **Validações:** placa no padrão Mercosul ou no padrão anterior; placa única entre ativos não inativos da empresa; RENAVAM com 11 dígitos e dígito verificador; VIN com 17 caracteres, sem I, O e Q; capacidade do tanque maior que zero.
- **Regras:** veículo com "transporta produto perigoso" exige documentos CIV e CIPP na agenda (RF-CON-005). Veículo com cronotacógrafo exige documento de verificação Inmetro. O VIN e a placa são as chaves de casamento com telemática (RF-INT-004).

#### RF-ATV-004 · Composição de veículos (cavalo e carreta)

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** registro de engate e desengate de implementos (carreta, semirreboque, dolly) num veículo trator.
- **Dados principais:** veículo trator, implemento, posição na composição, data e hora de engate e desengate, leitura de km no engate e no desengate.
- **Operações:** engatar, desengatar, consultar composição atual e histórico.
- **Validações:** implemento só em uma composição por vez; desengate posterior ao engate; implemento e trator não inativos.
- **Regras:** o km do implemento é derivado do km do trator durante o período engatado (RF-MED-006). Checklist de motorista feito no trator oferece os itens do implemento engatado (RF-CHK-004).

#### RF-ATV-005 · Classes de ativo

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** famílias de ativos com características comuns (`classe_ativo`).
- **Dados principais:** nome, tipo de ativo, sigla de TAG, campos extras (texto, número, data, lista, sim/não, com obrigatoriedade), medidores padrão, planos padrão, checklists padrão, códigos de falha aplicáveis, teto de salto de medidor, configuração de eixos (veículos), capacitações exigidas para intervenção (Fase 2).
- **Operações:** criar, editar, copiar de uma biblioteca inicial de classes comuns (empilhadeira elétrica, compressor de parafuso, cavalo mecânico 6x2, carreta 3 eixos, gerador, câmara fria) (proposta).
- **Validações:** nome único na empresa; campo extra obrigatório adicionado a uma classe com ativos existentes não bloqueia os ativos antigos, mas sinaliza pendência de preenchimento.
- **Regras:** ao criar um ativo de uma classe, o sistema propõe medidores, planos e checklists da classe. Alterar a classe não altera planos já instanciados nos ativos, salvo comando explícito "aplicar aos ativos".

#### RF-ATV-006 · Criticidade e status operacional

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** criticidade A, B ou C e status operacional `operando`, `em_manutencao`, `parado` ou `inativo`, com histórico.
- **Dados principais:** criticidade, meta de disponibilidade por criticidade (parâmetro), status atual, histórico de status (início, fim, motivo, OS relacionada, usuário).
- **Operações:** alterar criticidade com motivo; alterar status manualmente com motivo; ver linha do tempo de status.
- **Validações:** retorno a `operando` bloqueado enquanto restar uma causa de `parado` ou de `em_manutencao` (abaixo).
- **Regras:** o status segue o sumário (seção 3.3) e muda automaticamente pelas OS (RF-OS-005) e pela conformidade:

  | Status | Quando |
  |---|---|
  | `operando` | Nenhuma causa abaixo |
  | `em_manutencao` | Há OS em `em_execucao` ou `pausada` e o ativo foi retirado de operação de forma programada (OS `preventiva`, `preditiva`, `inspecao_legal`, `corretiva_planejada`, `melhoria` ou `implantacao` em ativo `em_operacao`, marcada "retira o ativo de operação") |
  | `parado` | Falha que parou o ativo ("o ativo parou?" = sim), `corretiva_emergencial` em curso, parecer NR-12 `inapto` (RF-CON-002) ou bloqueio de conformidade (`bloqueio_ativo`) |
  | `inativo` | Baixado, vendido ou em reserva (RF-ATV-011), ou em implantação (`fase_ciclo_vida = em_implantacao`, RF-ATV-013) |

  - `parado` prevalece sobre `em_manutencao`.
  - O status volta a `operando` no **encerramento** da OS, se não restar outra causa.
  - Os KPIs de tempo (MTTR, disponibilidade) usam `data_inicio`, `data_conclusao` e os registros de parada, não o encerramento administrativo; uma liberação tardia não infla o MTTR.
  - **Atenção** não é status gravado: é condição calculada sobre `operando` (plano vencendo, documento vencendo, leitura fora da faixa), exibida com a cor de Atenção da marca.
  - **Apto para despacho** (veículo) é verificação calculada, separada do status e sem entidade própria (RF-CON-007): um veículo `operando` pode estar impedido por documento crítico vencido.
  - `bloqueio_ativo` é a entidade que deixa o ativo `parado` (bloqueio de conformidade); não se confunde com `bloqueio_energia`, o registro de bloqueio e etiquetagem (LOTO) de uma intervenção (RF-CON-004, RF-CON-011) (sumário 3.6).
  - Ativo cedido e em locação continua `operando`; sua localização passa a ser o local de uso do locatário enquanto o item do contrato estiver entregue (RF-LOC-003).
  - O `bloqueio_ativo` tem `origem` `manual`, `conformidade` ou `regra` (RF-ADM-016, sumário 3.6).
  - `inativo` sai do denominador dos KPIs. A criticidade define a prioridade padrão das OS do ativo: A → `alta`, B → `media`, C → `baixa`.

#### RF-ATV-007 · QR code e etiquetas

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** identificação física de ativos e localizações por QR code ou código de barras.
- **Dados principais:** código vinculado (gerado ou existente), modelo de etiqueta.
- **Operações:** gerar etiquetas em lote em PDF para impressão; vincular código já colado no ativo; ler pelo app para abrir o ativo; ler pela câmera do celular sem app para abrir o portal de solicitação (RF-SOL-004).
- **Validações:** código único na empresa.
- **Regras:** o QR não contém dado do ativo, só um identificador opaco; o conteúdo é resolvido no servidor de acordo com quem lê.

#### RF-ATV-008 · Documentos técnicos e anexos

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** manuais, desenhos, catálogos de peças, fotos e procedimentos do fabricante ligados ao ativo ou à classe.
- **Dados principais:** arquivo, título, tipo, versão, ativo ou classe, marcação "disponível offline".
- **Operações:** anexar, substituir por nova versão, baixar, marcar para uso offline no app.
- **Validações:** tamanho máximo por arquivo (RNF-033); tipos permitidos (PDF, imagens, vídeo curto, planilhas).
- **Regras:** documento ligado à classe aparece em todos os ativos da classe. Versões anteriores ficam acessíveis. Na Fase 2, os documentos alimentam a busca por IA (RF-IA-002), respeitando as permissões do usuário.

#### RF-ATV-009 · Lista de peças do ativo

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** relação de itens de estoque aplicáveis a um ativo ou classe.
- **Dados principais:** item de estoque, quantidade instalada, posição ou subconjunto, observação.
- **Operações:** incluir, remover, consultar saldo a partir do ativo; sugerir peças ao registrar consumo na OS.
- **Validações:** item existente e ativo.

#### RF-ATV-010 · Histórico do ativo

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** linha do tempo única do ativo.
- **Dados principais:** OS, solicitações, leituras, checklists, mudanças de status, documentos, movimentações de peças, custos acumulados e, na Fase 2, abastecimentos, pneus, multas e sinistros.
- **Operações:** filtrar por tipo de evento e período; exportar em PDF; abrir offline no app (últimos 12 meses, proposta).
- **Regras:** o histórico respeita permissões; o `solicitante` vê só as próprias solicitações.

#### RF-ATV-011 · Transferência e baixa de ativo

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** mudança de localização e inativação definitiva (venda, sucata, perda total).
- **Dados principais:** data, origem, destino, motivo, valor de venda ou residual, documento comprobatório.
- **Operações:** transferir entre localizações; transferir entre unidades (exige RF-ADM-003); inativar com motivo.
- **Validações:** inativação bloqueada com OS não encerrada ou planos ativos, até que sejam encerrados ou cancelados; filhos acompanham o pai ou são realocados antes.
- **Regras:** inativar um ativo suspende seus planos, documentos e cobranças de veículo a partir da data informada.

#### RF-ATV-012 · Gestão de ativos ISO 55001

**Fase:** Fase 3 · **Plano mínimo:** Enterprise

- **Descrição:** apoio à ISO 55001:2024 com risco, custo de ciclo de vida e plano de ativo.
- **Dados principais:** matriz de risco (probabilidade e consequência), RAV, custo de ciclo de vida, vida útil estimada, plano de ativo com ações ligadas a riscos.
- **Operações:** avaliar risco por ativo; ligar riscos a planos de manutenção; gerar relatório de evidências para auditoria de certificação.
- **Regras:** a ISO 55001 não é exigência legal; o recurso atende clientes que buscam certificação. A transição das certificadoras vai até 31/07/2027, segundo o estudo.

#### RF-ATV-013 · Fase de ciclo de vida e entrada em operação

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** separação entre o ativo que está sendo instalado e o ativo que já opera, para que a implantação não conte como indisponibilidade nem como custo de manutenção.
- **Dados principais:** `fase_ciclo_vida` (`em_implantacao`, `em_operacao`), `data_entrada_operacao`, quem colocou em operação, projeto de origem (RF-PRJ-002, Fase 2), documento de entrega ou termo de recebimento (opcional).
- **Operações:** cadastrar o ativo já `em_implantacao` (manualmente, por importação ou por um projeto); abrir OS `implantacao` no ativo; ação `colocar_em_operacao` (`POST /v1/ativos/{id}/colocar-em-operacao`), pela web ou pela API; colocar em operação vários ativos de uma vez; filtrar ativos por fase.
- **Validações:** `colocar_em_operacao` só em ativo `em_implantacao`; exige que não haja OS `implantacao` não encerrada no ativo, salvo cancelada (proposta); exige os campos obrigatórios da classe preenchidos e, quando a classe marcar, os documentos críticos válidos (ex.: verificação metrológica) (proposta); exige a permissão `atv.colocar_em_operacao`, padrão em `gestor` e `planejador`.
- **Regras:** o ativo `em_implantacao` fica com status `inativo`, sai dos KPIs e da base de veículos cobrados, e seus planos ficam suspensos; a interface mostra "Em implantação" (sumário 3.3). `colocar_em_operacao` muda a fase para `em_operacao`, o status para `operando`, grava `data_entrada_operacao`, ativa os planos a partir dessa data e gera o evento `ativo.colocado_em_operacao`. A fase **não volta** para `em_implantacao`: a reinstalação de um ativo em operação em outro ponto é transferência (RF-ATV-011) com OS `melhoria` marcada "retira o ativo de operação" ou, na Fase 2, projeto cuja etapa gera OS `implantacao` (RF-PRJ-004), com o ativo `em_manutencao` enquanto durar a OS. Ativo `de_cliente` (SRV) e ativo tomado (`de_terceiro`) também podem ter fase `em_implantacao`. Na conexão entre contas, a fase não é campo compartilhado no lançamento.

---

### 4.3 MED — Medidores

#### RF-MED-001 · Cadastro de medidor

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** contadores de uso de um ativo (`medidor`), do tipo `km`, `horas`, `ciclos` ou `outro`.
- **Dados principais:** ativo, tipo, nome (ex.: "Hodômetro", "Horímetro do motor", "Horímetro da PTO", "Horímetro da refrigeração"), unidade, valor inicial e data inicial, teto de salto (por dia), dias sem leitura para alerta, medidor principal do tipo (sim/não), situação.
- **Operações:** criar, editar, inativar; herdar da classe.
- **Validações:** nome único no ativo; um único medidor principal por tipo no ativo; todo `veiculo` tem um medidor `km` principal (proposta).
- **Regras:** os planos preventivos por medidor apontam para um medidor específico. O teto de salto vem da classe e pode ser ajustado por ativo.

#### RF-MED-002 · Registro de leitura

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** registro do valor de um medidor num instante (`leitura_medidor`), com origem rastreável.
- **Dados principais:** medidor, valor, data e hora da leitura, origem (`manual`, `telematica`, `cartao_combustivel`, `os`, `importacao`; `iot` na Fase 3), usuário ou integração, documento de referência (OS, checklist, abastecimento, lote de importação), foto do painel, geolocalização, situação (`valida`, `suspeita`, `rejeitada`) (proposta).
- **Operações:** registrar na web, no app (inclusive offline), na OS, no checklist, no abastecimento, por API e por importação.
- **Validações:** as de RF-MED-003.
- **Regras:** a data e hora válida é a do momento da leitura, não a da sincronização. Conforme o sumário (seção 3.1): leitura vinda de checklist ou de abastecimento digitado usa a origem `manual` com referência ao documento; leitura vinda de CSV ou da API de um integrador usa `importacao`; a de conector de telemática usa `telematica`; a de sensor usa `iot` (Fase 3, RF-INT-010). Leitura registrada na vistoria de locação usa `manual` com referência à vistoria (RF-LOC-005). Foto do painel pode ser obrigatória por classe.

#### RF-MED-003 · Validação de leitura

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** bloqueio de leituras impossíveis e quarentena de leituras improváveis.
- **Validações:**
  - **Sem regressão:** uma leitura não pode ser menor que a leitura válida anterior nem maior que a leitura válida seguinte na ordem cronológica, exceto após troca de medidor (RF-MED-004).
  - **Teto de salto:** o incremento dividido pelos dias decorridos não pode passar do teto do medidor. Para `horas`, o incremento nunca pode passar das horas de relógio decorridas.
  - **Data futura:** leitura com data mais de 5 minutos à frente do relógio do servidor é rejeitada (proposta).
  - **Duplicidade:** mesma origem, mesmo valor e mesma data e hora são tratados como a mesma leitura (idempotência).
- **Regras:** leitura manual que viola o teto pode ser gravada com confirmação e justificativa, e fica `suspeita`. Leitura automática que viola o teto fica `suspeita` sem confirmação. Leitura `suspeita` não dispara preventiva nem entra em KPI até ser aprovada (RF-MED-007). Leitura com regressão é `rejeitada` e fica registrada para diagnóstico da origem.

#### RF-MED-004 · Troca de medidor

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** troca de hodômetro, painel ou horímetro, e virada de contador.
- **Dados principais:** data e hora, leitura final do medidor antigo, leitura inicial do novo, motivo (troca de painel, defeito, virada), OS relacionada.
- **Operações:** registrar troca; consultar segmentos.
- **Validações:** data posterior à última leitura válida; leitura final não menor que a última leitura válida.
- **Regras:** o medidor mantém um **valor acumulado** (soma dos segmentos) além do valor exibido no painel. Planos, KPIs e km de pneus usam o acumulado. Leituras posteriores à troca são validadas contra o novo segmento.

#### RF-MED-005 · Média de uso e projeção

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** cálculo da média diária de uso e alerta de medidor sem leitura.
- **Dados principais:** média diária móvel dos últimos 30 dias (proposta), data da última leitura válida.
- **Operações:** exibir média e data prevista de vencimento das preventivas por medidor; alertar medidor sem leitura há mais dias que o parâmetro.
- **Regras:** sem leituras suficientes (menos de 2 nos últimos 30 dias), a projeção usa a média de toda a vida do medidor e é marcada como estimativa.

#### RF-MED-006 · Medidor derivado

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** medidor cujo valor é calculado a partir de outro ativo.
- **Operações:** configurar o implemento para herdar o km do trator enquanto engatado; configurar o componente para herdar o medidor do ativo pai enquanto instalado.
- **Regras:** o valor derivado é a soma dos incrementos do medidor de origem nos períodos de vínculo (engate ou instalação). Leituras derivadas não aceitam lançamento manual.

#### RF-MED-007 · Revisão de leituras suspeitas

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** fila de leituras `suspeitas` para decisão humana.
- **Operações:** aprovar, corrigir o valor (gera nova leitura e marca a antiga como `rejeitada`), rejeitar.
- **Validações:** só `gestor` ou `planejador` decidem; motivo obrigatório na correção e na rejeição.
- **Regras:** ao aprovar, o sistema reavalia os planos por medidor afetados. Toda decisão vai para a auditoria.

#### RF-MED-008 · Leitura em lote

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** registro rápido de muitas leituras de uma vez.
- **Operações:** grade na web por localização ou classe; ronda de leituras no app com leitura de QR e próximo ativo automático.
- **Validações:** cada linha passa por RF-MED-003; o lote mostra resultado por linha.

---

### 4.4 SOL — Solicitações

#### RF-SOL-001 · Abertura de solicitação

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** pedido de manutenção (`solicitacao`) aberto por qualquer usuário, inclusive `solicitante` e `motorista` gratuitos.
- **Dados principais:** número `SS-<AAAA>-<NNNNNN>`, ativo ou localização, descrição, fotos, vídeo curto ou áudio, "o ativo parou?", prioridade sugerida, solicitante, contato, data e hora, origem (web, app, portal QR, checklist, telemática, API, `regra` para a solicitação aberta por uma regra configurável (RF-ADM-016) e, na Fase 2, `portal_locatario` para a solicitação aberta pelo locatário, RF-LOC-008).
- **Operações:** abrir pela web, pelo app (inclusive offline), pelo portal QR, pela API; anexar mídia; acompanhar.
- **Validações:** ativo ou localização obrigatório; descrição com pelo menos 10 caracteres ou um anexo de mídia (proposta).
- **Regras:** a solicitação nasce `nova`. Ao digitar o ativo, o sistema mostra solicitações e OS abertas do mesmo ativo para evitar duplicidade. Solicitação offline recebe número provisório e o número definitivo na sincronização.

#### RF-SOL-002 · Triagem

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** análise da solicitação e decisão.
- **Operações:** assumir (`nova` → `em_triagem`); aprovar e gerar OS (`aprovada`), definindo tipo, prioridade e responsável; recusar com motivo (`recusada`); marcar como `duplicada` a partir de `nova` ou `em_triagem`, vinculando a outra solicitação ou OS; aprovar vinculando a uma OS já existente sem criar nova (proposta); registrar a desistência do solicitante como `recusada` com motivo `desistencia`, direto de `nova` ou de `em_triagem`; **reabrir como nova** uma solicitação marcada como `duplicada` por engano.
- **Validações:** motivo obrigatório na recusa; solicitação ou OS de referência obrigatória na duplicidade; só `planejador`, `gestor` ou perfil com permissão de triagem decidem.
- **Regras:** abrir a triagem grava o instante de reconhecimento da solicitação. A OS gerada herda ativo, descrição, mídia e "o ativo parou?" da solicitação. Solicitação aberta pelo `locatario` sobre ativo locado (RF-LOC-008) traz o contrato de locação, e a OS gerada recebe o responsável pelo custo padrão da matriz do contrato (RF-LOC-004). Regras de aprovação automática configuráveis (por exemplo, ativo de criticidade A parado gera OS `corretiva_emergencial` sem triagem manual) (proposta). Estados finais (`aprovada`, `recusada`, `duplicada`) não voltam atrás; uma nova solicitação deve ser aberta. `nova` → `recusada` só é permitido com motivo `desistencia`; as demais recusas passam pela triagem. O próprio solicitante pode desistir enquanto a solicitação estiver `nova` (proposta). "Reabrir como nova" não altera a solicitação `duplicada`: cria uma solicitação `nova` com os mesmos dados e mídia, ligada à original, e registra o motivo na auditoria (sumário 3.6).

#### RF-SOL-003 · Acompanhamento pelo solicitante

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** retorno ao solicitante sobre o andamento.
- **Operações:** ver as próprias solicitações e o estado da OS gerada; comentar; receber notificação de triagem, início e conclusão; avaliar o atendimento (proposta).
- **Regras:** o solicitante vê só as próprias solicitações e um resumo da OS (estado, data prevista, conclusão), sem custos nem apontamentos.

#### RF-SOL-004 · Portal de solicitação por QR code

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** abertura de solicitação sem login, lendo o QR do ativo ou da localização com a câmera do celular.
- **Dados principais:** nome e contato do solicitante (configurável como obrigatório), descrição, foto.
- **Operações:** ligar ou desligar o portal por empresa e por unidade; exigir identificação.
- **Validações:** limite de solicitações por origem e por intervalo contra abuso (RNF-014); anexos limitados em tamanho.
- **Regras:** o portal mostra só o nome e a TAG do ativo, nunca histórico ou custos.

#### RF-SOL-005 · Prazos de atendimento

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** prazos-alvo de triagem e de reconhecimento por prioridade.
- **Dados principais:** prioridade, prazo de triagem, prazo de reconhecimento da OS (proposta: urgente 1 h, alta 4 h, média 1 dia útil, baixa 3 dias úteis, todos configuráveis).
- **Operações:** configurar prazos; ver solicitações e OS com prazo estourado; receber alerta antes do vencimento.
- **Regras:** os prazos usam o calendário da unidade. O cumprimento entra no MTTA (RF-IND-003).

#### RF-SOL-006 · Solicitação gerada pelo sistema

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** criação automática de solicitação por eventos.
- **Operações:** gerar solicitação por item reprovado de checklist (RF-CHK-003); por DTC de telemática (RF-INT-006, Fase 2); por alerta de condição (RF-PRE-011, Fase 3).
- **Regras:** a solicitação automática registra a origem e o evento que a gerou. Eventos repetidos do mesmo ativo e da mesma causa, com solicitação ainda aberta, incrementam um contador em vez de abrir outra.

---
### 4.5 OS — Ordens de serviço

#### RF-OS-001 · Criação de OS

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** registro de um trabalho de manutenção (`ordem_servico`).
- **Dados principais:** número `<prefixo>-<AAAA>-<NNNNNN>`, com o prefixo pelo país da empresa (`OS-` no Brasil; `OT-`, de *orden de trabajo*, na Argentina e no Paraguai), sequencial por empresa e reiniciado a cada ano, gravado na criação e nunca alterado pelo idioma de quem lê (RNF-043, RNF-044), tipo, prioridade (`baixa`, `media`, `alta`, `urgente`), ativo ou localização, unidade, título e descrição, `origem` (`manual`, `solicitacao`, `plano`, `checklist`, `telematica`, `rascunho_ia`, `vistoria_locacao`, `prestador`, `cliente_conectado`, `projeto`), solicitação, plano e etapa de projeto de origem, campos extras (`campos_extras`, RF-OS-021), causa externa da parada (`causa_externa`, opcional: `vandalismo`, `acidente_terceiro`, `falta_energia`, `obra_terceiro`, `solicitacao_cliente`, `caso_fortuito`), "o ativo parou?", "retira o ativo de operação" (OS programadas), data e hora da parada, data-limite, responsável, equipe, **executor** (`executor_tipo`: `interno` ou `fornecedor`, com `fornecedor_id`), **responsável pelo custo** (`responsavel_custo`: `empresa`, `locatario`, `locador`, `garantia` ou `fornecedor`), `justificativa_responsavel_custo` (quando o responsável difere do padrão), contrato de locação ou de prestação aplicável, anexos. Na conta de uma subcontratada (SRV), também `descricao_servico` e os itens `servico_os` do serviço a realizar (RF-SRV-007).
- **Operações:** criar na web, no app (inclusive offline), pela API, a partir de solicitação, de plano ou de vistoria de locação; receber a OS aberta por uma subcontratada conectada (origem `prestador`, RF-SRV-017).
- **Validações:** tipo obrigatório; ativo ou localização obrigatório; data-limite não anterior à abertura; `fornecedor_id` obrigatório quando `executor_tipo = fornecedor`; `responsavel_custo = locatario` só em ativo com contrato de locação `cedido` vigente na data da abertura ou da avaria; `responsavel_custo = locador` só em ativo com contrato `tomado`.
- **Regras:** a OS nasce `aberta`, salvo as geradas por plano, que nascem `planejada` (RF-PRE-009). O número é atribuído pelo servidor; OS criada offline recebe número provisório e o definitivo na sincronização, sem lacuna visível ao usuário de campo (proposta). A prioridade padrão vem da criticidade do ativo (A → `alta`, B → `media`, C → `baixa`) e pode ser alterada; `corretiva_emergencial` nasce `urgente`. O responsável pelo custo padrão é `empresa`; em ativo locado vem da matriz do contrato (RF-LOC-004) e pode ser alterado com `justificativa_responsavel_custo`; em ativo ou peça em garantia o sistema sugere `garantia` (RF-OS-020). OS só de localização não entra nos KPIs por ativo. OS de origem `prestador` nasce `aberta` e passa pela triagem normal da empresa (sumário 3.5).

#### RF-OS-002 · Tipos de OS

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** classificação obrigatória que alimenta os KPIs: `preventiva`, `preditiva`, `corretiva_planejada`, `corretiva_emergencial`, `inspecao_legal`, `melhoria`, `implantacao`.
- **Validações:** `preventiva` e `inspecao_legal` geradas por plano mantêm o vínculo com o plano; `inspecao_legal` exige referência à obrigação atendida (documento, NR-12, NR-13) (proposta); `implantacao` só em ativo `em_implantacao`, ligada a uma etapa de projeto (RF-PRJ-004) ou, na conta da subcontratada, recebida de um cliente conectado (origem `cliente_conectado`), e nunca é gerada por plano.
- **Regras:** `corretiva_emergencial` pode ir de `aberta` direto para `em_execucao`. A OS `implantacao` segue o mesmo ciclo de vida, fica fora dos KPIs de manutenção (MTBF, MTTR, disponibilidade, % planejada, cumprimento de preventiva, backlog de manutenção) e seu custo é apresentado como custo de implantação, separado do custo de manutenção (RF-OS-013, RF-IND-005). A troca de tipo depois de `concluida` exige permissão de `gestor`, motivo e gera auditoria, porque altera % planejada e PM compliance.

#### RF-OS-003 · Ciclo de vida da OS

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** estados canônicos `aberta`, `planejada`, `programada`, `em_execucao`, `pausada`, `concluida`, `encerrada` e `cancelada`.
- **Operações:** reconhecer, assumir, planejar, programar, desprogramar, iniciar, pausar, retomar, concluir, encerrar (liberar), recusar a liberação, cancelar.
- **Validações:** `pausada` exige motivo da lista canônica (`aguardando_peca`, `aguardando_aprovacao`, `aguardando_terceiro`, `fim_turno`, `outro`), com texto obrigatório em `outro`; `cancelada` exige motivo; recusa da liberação exige `motivo_recusa_liberacao`; encerrar e recusar a liberação exigem a permissão `os.encerrar`; `encerrada` não aceita edição; OS cancelada no servidor bloqueia a conclusão enviada depois por um app offline (RNF-011).
- **Regras:** as transições permitidas são as do diagrama do sumário (seção 3.3):

  | De | Para | Observação |
  |---|---|---|
  | (criação) | `aberta` | OS manual, de solicitação, de checklist, de telemática, de vistoria, de rascunho de IA ou aberta por subcontratada conectada (origem `prestador`) |
  | (criação) | `planejada` | OS gerada por plano (RF-PRE-009) |
  | `aberta` | `planejada` | Planejar |
  | `aberta` | `programada` | Programar sem planejar |
  | `aberta` | `em_execucao` | Só `corretiva_emergencial` |
  | `planejada` | `programada` | Programar |
  | `programada` | `planejada` | Desprogramar |
  | `programada` | `em_execucao` | Iniciar |
  | `em_execucao` | `pausada` | Pausar, com motivo |
  | `pausada` | `em_execucao` | Retomar |
  | `em_execucao` | `concluida` | Concluir (RF-OS-014) |
  | `concluida` | `em_execucao` | Liberação recusada, com `motivo_recusa_liberacao` e a permissão `os.encerrar` |
  | `concluida` | `encerrada` | Liberar, com a permissão `os.encerrar` |
  | `aberta`, `planejada`, `programada`, `pausada` | `cancelada` | Cancelar, com motivo |

  OS em `em_execucao` não é cancelada direto: é pausada antes. `concluida` e `encerrada` nunca são canceladas. Quem pode executar cada transição fica em 4_FLUXOS_PROCESSO.md. Correção de OS encerrada é feita por OS nova vinculada como retrabalho (RF-OS-015).

#### RF-OS-004 · Marcos de tempo da OS

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** registro dos instantes que sustentam MTTA, MTTR, disponibilidade e backlog.
- **Dados principais:** `data_abertura`; `data_reconhecimento` (quando um técnico, planejador ou prestador assume a OS); `data_inicio` (primeira entrada em `em_execucao`); pausas (início, fim, motivo, observação); `data_conclusao` (última entrada em `concluida`); `data_encerramento` (entrada em `encerrada`); `data_cancelamento`. Para cada marco: hora do dispositivo e hora de recebimento no servidor.
- **Operações:** marcar automaticamente nas transições; ação `reconhecer`, que grava `data_reconhecimento` sem mudar o estado; ajustar manualmente com permissão.
- **Validações:** ordem cronológica (`data_abertura` ≤ `data_reconhecimento` ≤ `data_inicio` ≤ `data_conclusao` ≤ `data_encerramento`); pausas dentro do intervalo de execução e sem sobreposição.
- **Regras:** reconhecimento não é estado. `data_reconhecimento` é gravada uma única vez: pela ação `reconhecer` ou, implicitamente, ao programar, iniciar ou assumir a OS, o que ocorrer primeiro. É a base do MTTA. Quando a liberação é recusada e a OS volta a `em_execucao`, `data_conclusao` passa a ser a da nova conclusão e o período entre as duas conta no MTTR (proposta). Ação offline usa a hora do dispositivo. Se o relógio do dispositivo divergir mais de 5 minutos do servidor na sincronização, os marcos ficam sinalizados para revisão (proposta). Ajuste manual de marco exige `gestor`, justificativa e gera auditoria.

#### RF-OS-005 · "O ativo parou?"

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** indicação de que o trabalho envolve indisponibilidade do ativo, por falha ou por retirada programada.
- **Dados principais:** "o ativo parou?" (falha que parou o ativo: sim ou não); "retira o ativo de operação" (OS `preventiva`, `preditiva`, `inspecao_legal`, `corretiva_planejada`, `melhoria` e `implantacao` em ativo `em_operacao`: sim ou não); data e hora da parada (padrão: abertura, editável para antes); data e hora do retorno à operação.
- **Operações:** informar na solicitação, na abertura, no planejamento ou durante a execução.
- **Regras:** as mudanças de status seguem RF-ATV-006 e o sumário (seção 3.3):
  - "o ativo parou?" = sim, ou `corretiva_emergencial` em curso: o ativo fica `parado` desde a parada até o encerramento, inclusive durante a execução e as pausas.
  - OS programada marcada "retira o ativo de operação": o ativo fica `em_manutencao` enquanto a OS estiver em `em_execucao` ou `pausada`.
  - `parado` prevalece sobre `em_manutencao`.
  - O ativo volta a `operando` no **encerramento** da OS, se não restar outra causa. Com várias OS no mesmo ativo, o retorno só ocorre quando a última causa terminar.
  - Em máquina NR-12 com parecer `inapto` ou com `bloqueio_ativo` vigente, o ativo não volta a `operando` (RF-CON-002).
  - Os KPIs de tempo usam `data_inicio`, `data_conclusao` e os registros de parada (RF-IND-003), por isso a liberação tardia não infla o MTTR.

#### RF-OS-006 · Planejamento da OS

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** definição do trabalho antes da execução.
- **Dados principais:** tarefas (`tarefa_os`, manuais ou de um `procedimento`), HH estimadas, peças previstas, serviços de terceiros previstos, ferramentas, medidas de segurança, capacitações exigidas (Fase 2), anexos.
- **Operações:** incluir procedimento, tarefas e peças; estimar HH; marcar como `planejada`.
- **Validações:** para ir a `planejada`, pelo menos uma tarefa e HH estimadas maiores que zero (proposta).
- **Regras:** a OS guarda a versão do procedimento usada; mudanças futuras no procedimento não alteram OS existentes.

#### RF-OS-007 · Programação e atribuição

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** definição de data e responsável.
- **Dados principais:** data e hora programadas, duração prevista, responsável (técnico interno ou fornecedor), auxiliares, técnicos do fornecedor (pessoas, quando informados).
- **Operações:** programar e atribuir pela lista, pelo calendário ou pelo quadro; reatribuir; desprogramar (`programada` → `planejada`); **assumir**: o técnico se atribui uma OS `aberta` ou `planejada` da sua equipe, pela web ou pelo app (`POST /v1/ordens-servico/{id}/assumir`), o que grava `data_reconhecimento` se ainda não houver (sumário 3.6); atribuir a um fornecedor (RF-TER-003), inclusive a um fornecedor com conta conectada, o que envia a OS à conta dele (RF-SRV-016).
- **Validações:** responsável interno ativo, com perfil `tecnico` e acesso à unidade; fornecedor ativo com papel `subcontratada`, `assistencia_tecnica` ou `oficina`; a partir da Fase 2, capacitação válida para OS regulada (RF-CON-013), também para técnicos de fornecedor informados como pessoas.
- **Regras:** a atribuição notifica o técnico, ou os usuários `prestador` do fornecedor, e envia a OS ao cache do app. Assumir não muda o estado da OS: define o técnico como responsável; a OS passa a `programada` quando recebe data (proposta). Assumir só vale para OS ainda sem responsável (proposta) e só para quem tem perfil `tecnico` com acesso à unidade; OS atribuída a fornecedor não pode ser assumida por técnico interno. A programação por capacidade com arrastar e soltar vem na Fase 3 (RF-OS-019).

#### RF-OS-008 · Execução

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** registro do trabalho em campo, na web ou no app offline.
- **Dados principais:** situação de cada tarefa (feita, não aplicável, com valor medido), fotos antes e depois, observações, leitura de medidor na conclusão (origem `os`), checklist anexado.
- **Operações:** iniciar, pausar com motivo, retomar, marcar tarefas, registrar medições, anexar mídia, concluir.
- **Validações:** tarefa com valor fora da faixa exige observação; foto obrigatória quando a tarefa ou o parâmetro exigir; leitura de medidor obrigatória na conclusão de preventiva por medidor (proposta).
- **Regras:** vários técnicos podem trabalhar na mesma OS; cada um tem seu apontamento. A leitura informada na conclusão passa pela validação de RF-MED-003 e rebaseia o plano de origem.

#### RF-OS-009 · Apontamento de horas

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** horas trabalhadas por técnico (`apontamento`).
- **Dados principais:** técnico, OS, início e fim (ou duração), tipo de hora (normal, extra) (proposta), custo por hora vigente.
- **Operações:** cronômetro automático a partir de iniciar e pausar; lançamento manual; ajuste com justificativa.
- **Validações:** fim depois do início; o mesmo técnico não tem apontamentos sobrepostos em OS diferentes, salvo justificativa aprovada por `gestor` (proposta).
- **Regras:** o custo de mão de obra é horas × custo por hora vigente na data do apontamento. Apontamentos de OS `encerrada` não são alterados.

#### RF-OS-010 · Peças e materiais na OS

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** consumo de itens de estoque e de materiais não estocados.
- **Dados principais:** item, almoxarifado, quantidade, custo unitário (custo médio do item), número de série (item que controla série, Fase 2), material avulso com descrição e valor.
- **Operações:** baixar do estoque (`saida_os`); devolver ao estoque (`devolucao_os`); lançar material avulso; na Fase 2, instalar peça serializada no ativo e retirar a peça substituída (RF-EST-016). Na conta de uma subcontratada, cada item de material tem `tipo_material` (`consumido`, `substituido`, `fornecido_pelo_cliente`) e, quando substituído, `numero_serie_retirado` e `numero_serie_instalado` (RF-SRV-008).
- **Validações:** quantidade maior que zero; saldo suficiente para baixa feita online; peça `em_transito` não é baixada antes da confirmação de recebimento (RF-EST-015); OS executada por fornecedor só baixa do almoxarifado `em_terceiro` daquele fornecedor ou de almoxarifado `proprio` com autorização do `almoxarife` (Fase 2, RF-EST-013).
- **Regras:** cada baixa gera `movimentacao_estoque` de saída com referência à OS. Baixa feita offline é aceita mesmo sem saldo e gera pendência de ajuste para o `almoxarife`, em vez de saldo negativo silencioso (proposta). Peças reservadas (RF-EST-006) são baixadas da reserva. Peça reparável retirada do ativo vai para o almoxarifado `quarentena` (Fase 2, RF-EST-017).

#### RF-OS-011 · Registro de falha

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** codificação simplificada da falha, inspirada na ISO 14224, com **modo**, **causa** e **atividade** (`codigo_falha`).
- **Dados principais:** sintoma (texto), modo de falha, causa, atividade executada, componente afetado.
- **Operações:** selecionar do catálogo padrão da Tenaz; buscar por texto.
- **Validações:** obrigatório em OS `corretiva_planejada` ou `corretiva_emergencial` com "o ativo parou?" = sim, quando o parâmetro estiver ligado (padrão ligado, proposta).
- **Regras:** o catálogo padrão é fixo no Essencial. O catálogo personalizável por classe é RF-OS-012. Sem a falha registrada, a OS conta no MTBF, mas fica fora das análises por modo de falha.

#### RF-OS-012 · Catálogo de códigos de falha personalizável

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** catálogo próprio de modos, causas e atividades por classe de ativo.
- **Operações:** criar, editar, inativar códigos; ligar a classes; importar de planilha; mapear para o catálogo padrão.
- **Validações:** código único por tipo e classe; código usado em OS não é excluído, só inativado.
- **Regras:** o mapeamento para o catálogo padrão mantém comparáveis os indicadores entre classes.

#### RF-OS-013 · Custos da OS

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** custo total do trabalho, consolidado no encerramento.
- **Dados principais:** mão de obra própria, peças de estoque, materiais avulsos, serviços de terceiros (preço do catálogo, do contrato de prestação, do orçamento aprovado ou da nota), outros custos; centro de custo; responsável pelo custo (`responsavel_custo`); marcação "reparo capitalizável".
- **Operações:** ver custos parciais durante a execução; consolidar no encerramento; ratear por ativo quando a OS for de localização (proposta); alterar o responsável pelo custo até o encerramento, com `justificativa_responsavel_custo` obrigatória.
- **Regras:** depois de `encerrada`, o custo não muda. O custo da OS é registrado inteiro, qualquer que seja o responsável; o responsável define quem paga: OS com `responsavel_custo = locatario` alimenta o relatório de cobrança (RF-LOC-009); `locador` e `fornecedor` são custos a recuperar ou não incorridos pela empresa; `garantia` liga a OS à garantia (RF-OS-020). Os KPIs de custo mostram o custo total e o custo líquido da empresa (RF-IND-005). A fatura ao locatário ou a cobrança ao fornecedor são emitidas pelo ERP. Reparo que aumenta a vida útil em mais de um ano pode ser marcado como capitalizável para a contabilidade (regra fiscal do Decreto 9.580/2018 citada no estudo, não re-lida; o sistema só marca, não contabiliza).

#### RF-OS-014 · Conclusão, assinatura e liberação

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** fechamento técnico (`concluida`) pelo executante e liberação (`encerrada`) pelo responsável.
- **Dados principais:** resumo do serviço, assinatura do técnico (ou do prestador), assinatura de quem libera, data e hora, geolocalização, hash do conteúdo assinado; em máquina NR-12, os sete campos de RF-CON-002.
- **Operações:** concluir; liberar (`encerrar`); recusar a liberação com `motivo_recusa_liberacao` (a OS volta a `em_execucao`); emitir PDF assinado.
- **Validações:** todas as tarefas obrigatórias resolvidas; falha registrada quando exigida; leitura de medidor quando exigida; em máquina NR-12, os sete campos preenchidos; liberação e recusa só por quem tem a permissão `os.encerrar` (padrão em `gestor` e `planejador`); o `prestador` conclui, mas nunca libera. Em OS ligada a uma subcontratada conectada, a liberação é sempre da conta do cliente (RF-SRV-018).
- **Regras:** a assinatura é a eletrônica avançada de RF-CON-014, já no MVP e em todos os planos, com identificação do usuário autenticado, reautenticação, carimbo de tempo, geolocalização e hash. A execução registrada por link de uso único é assinada pelo prestador ocasional com nome, documento e aceite no próprio link (RF-TER-005). Se o parecer de **P1** exigir ICP-Brasil para OS de NR-12, essas OS passam a usar RF-CON-015. Ao encerrar, custos são consolidados e a OS fica imutável.

#### RF-OS-015 · Retrabalho

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** vínculo entre uma OS nova e uma anterior do mesmo ativo cujo problema voltou.
- **Operações:** sugerir vínculo quando uma corretiva for aberta no mesmo ativo dentro do prazo de retrabalho (padrão de 30 dias, proposta) após outra encerrada; confirmar ou descartar a sugestão.
- **Regras:** o vínculo alimenta o first-time fix (RF-IND-004).

#### RF-OS-016 · Comentários, menções e histórico

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** comunicação dentro da OS.
- **Operações:** comentar, mencionar usuários (gera notificação), anexar mídia, ver o histórico de estados e alterações.
- **Regras:** comentários não são editáveis depois de 15 minutos (proposta); exclusão é lógica e auditada.

#### RF-OS-017 · Listas, quadro e calendário

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** visões de trabalho para planejador, gestor e técnico.
- **Operações:** lista com filtros (estado, tipo, prioridade, ativo, responsável, unidade, atraso); quadro por estado; calendário por técnico; busca por número, TAG e texto; visões salvas por usuário; exportação CSV/XLSX.
- **Regras:** cada usuário vê apenas OS das unidades a que tem acesso; o `tecnico` vê por padrão as suas e as da sua equipe; o `prestador` vê só as OS atribuídas ao seu fornecedor (RNF-037); o `locatario` vê só o resumo das OS dos ativos locados a ele (RNF-038). Filtros adicionais: executor (interno ou fornecedor) e responsável pelo custo.

#### RF-OS-018 · Impressão da OS

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** OS em PDF para oficinas que ainda usam papel e para arquivo.
- **Operações:** imprimir OS em branco para execução (com QR para lançar depois no sistema) e OS encerrada com assinaturas.
- **Regras:** o PDF de OS encerrada traz o hash de verificação e a data de emissão. O relatório da OS de serviço para o cliente atendido (RF-SRV-013) usa o mesmo gerador de PDF.

#### RF-OS-019 · Programação por capacidade

**Fase:** Fase 3 · **Plano mínimo:** Profissional (proposta)

- **Descrição:** programação semanal com base nas HH disponíveis.
- **Dados principais:** calendário de cada técnico (turno, férias, folgas), HH disponíveis por semana, backlog priorizado.
- **Operações:** arrastar e soltar OS no calendário; ver carga por técnico e por especialidade; congelar a programação da semana.
- **Validações:** alerta de sobrecarga acima de 100% das HH disponíveis.
- **Regras:** a programação congelada é a base do schedule compliance (RF-IND-011).

#### RF-OS-020 · Garantia de peças e serviços

**Fase:** Fase 2 · **Plano mínimo:** Profissional (proposta)

- **Descrição:** controle de garantia do ativo, de peças instaladas (inclusive peças serializadas e recondicionadas) e de serviços de terceiros.
- **Dados principais:** item, peça serializada ou serviço, fornecedor, início, prazo em dias ou em medidor, documento, origem (compra, serviço de oficina ou assistência, ordem de reparo).
- **Operações:** registrar garantia na entrada da peça, no retorno da oficina (RF-TER-008) ou no retorno do conserto (RF-EST-018); avisar ao abrir OS em ativo ou peça em garantia.
- **Regras:** OS em garantia recebe `responsavel_custo = garantia` como sugestão, é marcada para cobrança do fornecedor e o custo pode ser lançado como recuperável (proposta). Falha de peça recondicionada dentro da garantia do reparo abre nova ordem de reparo sem custo para o mesmo fornecedor (proposta).

#### RF-OS-021 · Campos extras da OS

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** campos definidos pela empresa (`campo_extra_os`) para registrar na OS o que é próprio do segmento: números de lacre, protocolo da concessionária, número da autorização do órgão, leitura de um instrumento.
- **Dados principais:** definição: chave, rótulo, tipo (`texto`, `numero`, `data`, `lista`, `sim_nao`), opções (lista), unidade (número), tipos de OS e classes de ativo em que aparece, momento da obrigatoriedade (`nenhum`, `abertura`, `conclusao`), marcação "preenchido pelo prestador" (aparece para o perfil `prestador` e no link de uso único) (proposta), texto de ajuda, ordem, situação. Valores: `ordem_servico.campos_extras`, gravados com a chave da definição.
- **Operações:** criar, editar, ordenar e inativar definições; preencher na web, no app (inclusive offline), pelo link de uso único (só campos marcados como do prestador, proposta) e pela API; filtrar e exportar OS pelos campos; usar campo `sim_nao` como condição de regra (`os_concluida_com_campo`, RF-ADM-016).
- **Validações:** chave única na empresa, em minúsculas, sem acento; tipo não muda depois de haver valor gravado (cria-se um campo novo); valor conforme o tipo; campo obrigatório na conclusão impede `concluir` enquanto vazio, inclusive no app offline, que valida antes de enfileirar o comando; limite de 30 campos ativos por combinação de tipo de OS e classe (proposta).
- **Regras:** inativar a definição não apaga valores já gravados; a OS mostra o rótulo vigente na data do registro. Campo novo obrigatório não afeta OS já concluídas, nem a conclusão feita no app offline antes de o aparelho receber a definição nova: o servidor aceita essa conclusão sem o campo e a sinaliza para revisão (proposta). Criar e editar definições exige `adm.gerenciar_configuracao` (sumário 3.6). Os campos extras **não** trafegam na conexão entre contas no lançamento (sumário 3.6) e não entram em fórmula de indicador; aparecem no relatório da OS em PDF (RF-SRV-013) e na impressão (RF-OS-018). As definições viajam ao app junto com o cache de configuração.

---

### 4.6 PRE — Planos de manutenção

#### RF-PRE-001 · Procedimentos

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** lista padrão reutilizável de tarefas, peças e ferramentas (`procedimento`).
- **Dados principais:** nome, tarefas ordenadas (tipo de resposta: feito, sim/não, número com faixa aceitável, texto, foto), peças e quantidades, ferramentas, HH estimadas, medidas de segurança, capacitações exigidas (Fase 2), versão.
- **Operações:** criar, editar (gera nova versão), duplicar, inativar; anexar a planos e OS.
- **Validações:** pelo menos uma tarefa; faixa mínima menor ou igual à máxima.
- **Regras:** OS e planos guardam a versão usada.

#### RF-PRE-002 · Plano de manutenção

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** regra que gera OS periodicamente para um ativo ou para todos os ativos de uma classe (`plano_manutencao`).
- **Dados principais:** nome, escopo (ativo ou classe), tipo de OS gerada (`preventiva`, `preditiva` ou `inspecao_legal`), procedimento, prioridade, responsável padrão, gatilhos (`gatilho_plano`), antecipação, tolerância, marcação "ligado à segurança" (NR-12), situação, data de início.
- **Operações:** criar, editar, suspender, reativar, encerrar; aplicar plano de classe aos ativos da classe com ajustes por ativo.
- **Validações:** pelo menos um gatilho; ativo do escopo não inativo; medidor do gatilho pertence ao ativo.
- **Regras:** plano de classe cria uma instância por ativo, cada uma com seu próprio vencimento. Edição do plano vale para as próximas OS, não para as já geradas.

#### RF-PRE-003 · Gatilho por tempo

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** vencimento por intervalo de calendário.
- **Dados principais:** intervalo (dias, semanas, meses, anos); base fixa (calendário, ex.: todo dia 10) ou flutuante (a partir da última execução).
- **Validações:** intervalo maior que zero.
- **Regras:** na base fixa, o próximo vencimento não depende da data da execução. Na flutuante, conta a partir da data de fim da última OS do plano. Vencimento em dia sem expediente é antecipado para o último dia útil (proposta).

#### RF-PRE-004 · Gatilho por medidor

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** vencimento por uso (ex.: a cada 10.000 km ou 250 horas).
- **Dados principais:** medidor, intervalo, base fixa (múltiplos do intervalo) ou flutuante (leitura da última execução mais o intervalo).
- **Validações:** medidor ativo do mesmo ativo; intervalo maior que zero.
- **Regras:** a avaliação acontece a cada leitura válida e usa o valor acumulado (RF-MED-004). A data prevista de vencimento é projetada pela média de uso (RF-MED-005). Leitura `suspeita` não dispara geração.

#### RF-PRE-005 · Gatilho combinado ("o que vencer primeiro")

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** plano com gatilho de tempo e de medidor, vencendo pelo que ocorrer antes (ex.: 10.000 km ou 6 meses).
- **Regras:** a OS é gerada uma vez, pelo primeiro gatilho que vencer. Ao executar, os dois gatilhos são rebaseados. A OS registra qual gatilho venceu.

#### RF-PRE-006 · Antecipação e tolerância

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** quando a OS é gerada e qual janela conta como "em dia".
- **Dados principais:** antecipação (dias ou unidades de medidor antes do vencimento); tolerância (dias ou unidades antes e depois do vencimento).
- **Validações:** antecipação maior ou igual a zero; tolerância maior ou igual a zero.
- **Regras:** a OS é gerada quando falta, para o vencimento, menos que a antecipação. Executada dentro da tolerância, conta como cumprida no PM compliance e, na base fixa, o próximo vencimento não muda. Executada fora da janela, conta como fora do prazo. Cada plano e ativo tem no máximo uma OS aberta; se um novo vencimento chegar com a anterior ainda aberta, o sistema não cria outra, registra a ocorrência como perdida para o PM compliance e alerta o planejador.

#### RF-PRE-007 · Planos aninhados

**Fase:** Fase 2 · **Plano mínimo:** Essencial (proposta)

- **Descrição:** revisões em que a maior inclui a menor (ex.: revisão de 40.000 km inclui as de 20.000 e 10.000).
- **Operações:** declarar que um plano engloba outro.
- **Regras:** quando dois planos aninhados vencem na mesma janela, só o maior gera OS e o menor é rebaseado como executado por ela.

#### RF-PRE-008 · Previsão de preventivas

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** calendário projetado das próximas preventivas.
- **Operações:** ver previsão para até 12 meses por ativo, classe, unidade ou técnico; ver carga de HH por semana; exportar.
- **Regras:** gatilhos de medidor usam a projeção da média de uso e são marcados como estimados. Filtrada por planos "ligados à segurança", a previsão é o cronograma exigido pelo item 12.11.2.2 da NR-12 (RF-CON-003).

#### RF-PRE-009 · Geração automática de OS

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** rotina que transforma vencimentos em OS.
- **Operações:** avaliação de gatilhos de tempo em rotina periódica (pelo menos de hora em hora, proposta) e de medidor a cada leitura válida; geração manual antecipada pelo planejador.
- **Regras:** a geração é idempotente: a mesma ocorrência de vencimento nunca gera duas OS, mesmo com reprocessamento. A OS gerada herda tipo, procedimento, peças, HH, prioridade, responsável e executor (interno ou fornecedor) do plano e nasce `planejada` (sumário 3.3); a programação é um passo seguinte. Em ativo locado, o responsável pelo custo vem da matriz do contrato (RF-LOC-004); em ativo coberto por contrato de prestação, o executor padrão pode vir do contrato (RF-TER-006).

#### RF-PRE-010 · Biblioteca de planos padrão

**Fase:** MVP · **Plano mínimo:** Essencial (proposta)

- **Descrição:** modelos iniciais de planos e procedimentos para classes comuns.
- **Operações:** copiar modelo para a empresa e ajustar.
- **Regras:** o modelo copiado vira dado da empresa; atualizações da biblioteca não alteram cópias.

#### RF-PRE-011 · Manutenção baseada em condição

**Fase:** Fase 3 · **Plano mínimo:** Enterprise

- **Descrição:** gatilho por condição medida (vibração, temperatura, pressão, sulco, análise de óleo).
- **Dados principais:** variável, origem (inspeção manual ou sensor via RF-INT-010), limiar, regra de tendência (inclinação em janela de tempo), ação (solicitação ou OS `preditiva`).
- **Regras:** alertas repetidos com OS aberta não geram nova OS. Com IoT, o volume de dados segue os limites do plano Enterprise (RF-INT-010).

#### RF-PRE-012 · RCM e FMEA

**Fase:** Fase 3 · **Plano mínimo:** Enterprise

- **Descrição:** análise de modos de falha e efeitos que liga cada modo de falha a uma tarefa de manutenção.
- **Dados principais:** função, falha funcional, modo de falha (do catálogo), efeito, severidade, ocorrência, detecção, Action Priority, tarefa recomendada, plano gerado.
- **Operações:** montar análise por classe; gerar planos a partir das tarefas; revisar periodicamente.
- **Regras:** segue a lógica de Action Priority do manual AIAG-VDA (2019) e as sete perguntas da SAE JA1011. As duas fontes não foram re-lidas no estudo e devem ser conferidas antes da especificação detalhada.

---

### 4.7 CHK — Checklists

#### RF-CHK-001 · Modelos de checklist

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** listas de verificação reutilizáveis (`checklist_modelo`).
- **Dados principais:** nome, seções, itens com tipo (sim/não, número com faixa aceitável, foto, texto), item crítico, foto obrigatória em caso de reprovação, instrução e imagem de referência, exige assinatura, exige leitura de medidor, classes e ativos aplicáveis, frequência sugerida (por turno, diária, semanal, antes de viagem), versão.
- **Operações:** criar, editar (gera versão), duplicar, inativar.
- **Validações:** pelo menos um item; faixa válida nos itens numéricos.
- **Regras:** execuções guardam a versão do modelo usada.

#### RF-CHK-002 · Execução de checklist

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** preenchimento de um checklist num ativo e numa data (`checklist_execucao`), na web ou no app offline.
- **Dados principais:** modelo e versão, ativo, executante, início e fim, respostas, fotos com data, hora e geolocalização, leitura de medidor, assinatura.
- **Operações:** iniciar lendo o QR do ativo; responder; salvar parcial; concluir.
- **Validações:** itens obrigatórios respondidos; número fora da faixa vale como reprovado; foto exigida na reprovação quando configurada; leitura passa por RF-MED-003.
- **Regras:** execução concluída não é editada. Itens não aplicáveis exigem justificativa (proposta).

#### RF-CHK-003 · Item reprovado gera solicitação

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** abertura automática de solicitação a partir de reprovação.
- **Regras:** cada item reprovado gera uma solicitação com o texto do item, a resposta, as fotos e o ativo, ou uma solicitação agrupada por execução, conforme configuração do modelo. Item crítico reprovado marca a solicitação com "o ativo parou?" = sim e o ativo como `parado`, quando o modelo assim definir (proposta). Repetição do mesmo item reprovado com solicitação aberta segue a deduplicação de RF-SOL-006.

#### RF-CHK-004 · Checklist do motorista

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** inspeção de pré-viagem e pós-viagem pelo `motorista` gratuito.
- **Operações:** escolher o veículo por QR ou placa; responder o checklist do veículo e dos implementos engatados (RF-ATV-004); registrar km; assinar.
- **Regras:** checklist de veículo com item crítico reprovado deixa o veículo inapto até a liberação, e a situação aparece para o despacho (RF-CON-007, Fase 2). Pelo item 12.14.2.1 da NR-12, o registro da inspeção rotineira do operador de máquina não é obrigatório; o sistema oferece o recurso sem apresentá-lo como exigência legal.

#### RF-CHK-005 · Checklist dentro da OS

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** modelo de checklist anexado a uma OS, a um procedimento ou a um plano.
- **Regras:** a OS só é concluída com os checklists obrigatórios concluídos. Reprovação dentro de OS não gera solicitação automática; vira observação e, se o técnico decidir, OS nova vinculada (proposta).

#### RF-CHK-006 · Acompanhamento de checklists

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** controle de checklists esperados e realizados.
- **Operações:** ver execuções por ativo, executante e período; ver checklists esperados e não realizados conforme a frequência; ver índice de reprovação por item.
- **Regras:** "esperado" é calculado pela frequência do modelo e pelo calendário da unidade.

---
### 4.8 EST — Estoque e compras

**Nota de país (sumário 0.8, ADR-013).** A Tenaz não emite fatura nem remessa em país nenhum: guarda o tipo e a referência do documento fiscal do país da empresa (NF-e no Brasil; remito e fatura com CAE na Argentina; NRE e KuDE com CDC no Paraguai). Chave de 44 dígitos, CFOP e o leitor de XML da NF-e são do Brasil; o equivalente da Argentina e do Paraguai e os nomes dos campos dependem da **P13** e são fixados na etapa de EST. Valores em dinheiro na moeda da empresa, com a coluna `moeda` ao lado (RNF-045).

A fronteira entre os planos segue a tabela de recursos do sumário (seção 2.1):

| Conjunto | Plano mínimo | Fase | Requisitos |
|---|---|---|---|
| **Estoque básico:** itens, um almoxarifado por unidade, saldo, reserva e baixa em OS, estoque mínimo, ajuste com motivo | Essencial | MVP | RF-EST-001 a RF-EST-004, RF-EST-006 |
| **Estoque e compras completos:** vários almoxarifados, transferência, inventário, requisição, pedido de compra, entrada por NF-e | Profissional | Fase 2 | RF-EST-005, RF-EST-007 a RF-EST-009, RF-EST-011, RF-EST-012 |
| **Estoque em terceiros:** almoxarifado em poder de subcontratada, saldo geral e por subcontratada, remessa e retorno de peças | Profissional | Fase 2 | RF-EST-013 a RF-EST-015, RF-EST-019 |
| **Peças serializadas e reparáveis:** número de série, instalação em ativo, peça com defeito, envio para conserto e retorno | Profissional | Fase 2 | RF-EST-016 a RF-EST-018 |

O cadastro de fornecedores saiu deste módulo e passou a ser RF-TER-001.

Valores canônicos usados nesta seção (sumário 3.1):

- **Tipo de almoxarifado** (`almoxarifado.tipo`): `proprio`; `em_terceiro` (estoque da empresa em poder de uma subcontratada, com `fornecedor_id` obrigatório); `quarentena` (peças com defeito aguardando decisão). O **estoque geral** é a soma de todos os almoxarifados.
- **Tipo de movimentação** (`movimentacao_estoque.tipo`): `entrada`, `saida_os`, `devolucao_os`, `transferencia`, `remessa_terceiro`, `retorno_terceiro`, `envio_conserto`, `retorno_conserto`, `ajuste`, `baixa_sucata`.
- **NF-e em remessas e retornos:** a Tenaz registra o número (`nfe_numero`) e a chave (`nfe_chave`) da NF-e emitida pelo ERP; não emite nota. CFOPs e tratamento fiscal dependem de **P8**.

#### RF-EST-001 · Item de estoque

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** peça, material ou insumo controlado (`item_estoque`).
- **Dados principais:** código, descrição, unidade de medida, categoria, fabricante, part number, código de barras, NCM, estoque mínimo, estoque máximo, ponto de pedido, lead time de reposição, fornecedores, ativos e classes aplicáveis, item crítico, foto, custo médio, situação; a partir da Fase 2, **controla número de série** (`controla_serie`) e **reparável** (`reparavel`).
- **Operações:** criar, editar, inativar, importar por planilha, imprimir etiqueta.
- **Validações:** código único na empresa; mínimo ≤ ponto de pedido ≤ máximo, quando informados; NCM com 8 dígitos quando informado; `controla_serie` só é ligado num item com saldo depois de informar os números de série do saldo existente (proposta).
- **Regras:** item com movimentação não é excluído, só inativado. O custo médio é ponderado e recalculado a cada entrada (proposta). Item reparável sem controle de série é tratado por quantidade na quarentena e na ordem de reparo.

#### RF-EST-002 · Almoxarifados

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** locais de armazenagem (`almoxarifado`) ligados a uma unidade, com endereços opcionais.
- **Dados principais:** nome, unidade, tipo (`proprio`, `em_terceiro`, `quarentena`), fornecedor (no tipo `em_terceiro`), responsável, endereços (rua, prateleira, posição).
- **Operações:** criar, editar, inativar; definir endereço padrão por item.
- **Validações:** almoxarifado com saldo não é inativado; no Essencial, um único almoxarifado do tipo `proprio` por unidade; os tipos `em_terceiro` e `quarentena` e o segundo almoxarifado `proprio` da unidade exigem o Profissional (Fase 2).
- **Regras:** no Essencial, o almoxarifado da unidade é criado junto com ela. Vários almoxarifados por unidade e transferência entre eles: RF-EST-012. Almoxarifado em poder de terceiro: RF-EST-013. Quarentena: RF-EST-017.

#### RF-EST-003 · Movimentações

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** registro de toda mudança de saldo (`movimentacao_estoque`).
- **Dados principais:** tipo (lista canônica acima), item, almoxarifado de origem e de destino, quantidade, custo unitário, números de série (itens com série), data e hora, usuário, motivo, referência (OS, NF, pedido, inventário, centro de custo, ordem de reparo), número e chave da NF-e (remessas e retornos).
- **Operações:** no MVP, entrada manual (`entrada`), saída para OS (`saida_os`, RF-OS-010), devolução de OS (`devolucao_os`), ajuste com motivo (`ajuste`), inclusive saída avulsa para centro de custo (proposta), e estorno. Na Fase 2: `transferencia` (RF-EST-012), `remessa_terceiro` e `retorno_terceiro` (RF-EST-015), `envio_conserto` e `retorno_conserto` (RF-EST-018), `baixa_sucata` (RF-EST-017). Rota `/v1/movimentacoes-estoque`.
- **Validações:** quantidade maior que zero; saldo não fica negativo, exceto pela pendência de baixa offline descrita em RF-OS-010; ajuste exige motivo e perfil `almoxarife` ou `gestor`; tipos de Fase 2 recusados no Essencial.
- **Regras:** movimentação não é editada nem excluída; corrige-se por estorno, que gera movimentação inversa ligada à original. Toda movimentação gera o evento `movimentacao_estoque.registrada`.

#### RF-EST-004 · Saldo e ponto de pedido

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** posição de estoque por item e almoxarifado.
- **Dados principais:** saldo físico, reservado, disponível (físico menos reservado), valor em estoque.
- **Operações:** consultar saldo pela web, pelo app e a partir do ativo; ver itens abaixo do ponto de pedido; exportar.
- **Regras:** saldo abaixo do ponto de pedido gera alerta ao `almoxarife`. Item crítico abaixo do mínimo gera alerta também ao `gestor`. Com vários almoxarifados (Fase 2), a visão por almoxarifado e o estoque geral seguem RF-EST-014.

#### RF-EST-005 · Inventário

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** contagem física e ajuste de saldo.
- **Operações:** abrir inventário geral ou por almoxarifado, endereço ou categoria; contar pela web ou pelo app com leitura de código de barras; recontar divergências; aprovar ajustes; inventariar almoxarifado `em_terceiro`, com contagem feita pelo `prestador` ou por pessoa da empresa (proposta).
- **Validações:** aprovação por perfil diferente de quem contou (proposta); contagem feita pelo `prestador` sempre aprovada por usuário interno.
- **Regras:** o saldo considerado é o do momento da abertura; movimentações durante a contagem são tratadas na aprovação. Os ajustes aprovados geram movimentações de ajuste com referência ao inventário. Itens com série são contados por número de série. No Essencial, o acerto de saldo é feito por ajuste com motivo (RF-EST-003).

#### RF-EST-006 · Reserva para OS

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** separação lógica de saldo para OS planejadas ou programadas.
- **Operações:** reservar automaticamente as peças previstas ao programar a OS e manualmente ao planejar; liberar reserva no cancelamento e no encerramento; baixar da reserva na execução.
- **Validações:** reserva limitada ao saldo disponível; disponível insuficiente gera reserva parcial e pendência de falta; saldo em `quarentena` ou em conserto não é reservável.
- **Regras:** reserva de OS `programada` tem prioridade sobre reserva de OS `planejada` (proposta). Reserva sem baixa em 30 dias após a data programada gera alerta (proposta). A reserva não movimenta o saldo físico. As regras numeradas ficam em 4_FLUXOS_PROCESSO.md (RN-EST).

#### RF-EST-007 · Requisição de material

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** pedido de retirada de peças feito pelo técnico ao almoxarifado.
- **Operações:** requisitar a partir da OS; separar; entregar com confirmação de recebimento no app; devolver sobras.
- **Validações:** requisição ligada a OS ou a centro de custo.
- **Regras:** a entrega gera a saída do estoque com referência à OS.

#### RF-EST-008 · Compras

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** ciclo de solicitação de compra, cotação, aprovação e pedido de compra, conduzido pelo `comprador`.
- **Dados principais:** solicitação de compra (origem: ponto de pedido, OS, requisição, orçamento de fornecedor, manual), cotações por fornecedor, pedido de compra (fornecedor, itens, preços, prazo, condição de pagamento, almoxarifado de entrega, inclusive `em_terceiro`), alçadas de aprovação por valor, situação.
- **Operações:** sugerir compras a partir do ponto de pedido e das OS; cotar com vários fornecedores de papel `fornecedor_pecas`; comparar; aprovar conforme alçada; emitir pedido em PDF ou e-mail; acompanhar entrega; receber parcial ou total.
- **Validações:** pedido aprovado não muda de valor sem nova aprovação; recebimento não excede o pedido sem tolerância configurada (proposta).
- **Regras:** a contratação de serviço de oficina, assistência técnica ou conserto segue o orçamento de fornecedor (RF-TER-009). Com conector de ERP (RF-INT-008), o pedido pode ser emitido no ERP e acompanhado na Tenaz.

#### RF-EST-009 · Entrada por NF-e e NFS-e

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** entrada de mercadorias e serviços a partir do XML do documento fiscal.
- **Dados principais:** XML da NF-e (modelo 55) ou da NFS-e, chave de acesso, emitente, itens, quantidades, valores, impostos destacados, pedido de compra vinculado.
- **Operações:** importar XML por upload ou por caixa de e-mail dedicada (proposta); casar itens da nota com itens de estoque (por código do fornecedor, part number ou NCM, com memória dos casamentos anteriores); vincular a pedido de compra; ratear valor de NFS-e por OS e por ativo; confirmar entrada.
- **Validações:** (NF-e, só no Brasil; o documento fiscal da Argentina e do Paraguai depende da **P13**) chave de acesso com 44 dígitos e dígito verificador válido; destinatário com CNPJ da empresa ou de uma unidade; nota já importada é recusada; divergência de quantidade ou preço com o pedido exige aprovação.
- **Regras:** o custo de entrada atualiza o custo médio. O leitor de XML deve aceitar os novos grupos de CBS e IBS, já que 2026 é ano de teste desses tributos nos documentos fiscais (LC 214/2025, citada no estudo e não re-lida). A Tenaz não faz escrituração fiscal. NF-e de remessa, de retorno e de conserto não passam por esta entrada: são registradas em RF-EST-019.

#### RF-EST-010 · Fornecedores (movido)

Movido para **RF-TER-001** na versão 0.2. O ID não é reutilizado e não entra nas contagens.

#### RF-EST-011 · Curva ABC e inventário cíclico

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** classificação de itens por valor consumido e contagem periódica por classe.
- **Operações:** calcular curva ABC dos últimos 12 meses; definir frequência de contagem por classe (ex.: A mensal, B trimestral, C semestral) (proposta); gerar contagens cíclicas automáticas; ver giro e itens sem movimento.
- **Regras:** contagens cíclicas usam o fluxo de RF-EST-005.

#### RF-EST-012 · Vários almoxarifados e transferência

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** mais de um almoxarifado `proprio` por unidade (central, oficina, veículo de apoio) e transferência de saldo entre almoxarifados.
- **Dados principais:** almoxarifado de origem e de destino, itens, quantidades, números de série, responsável pelo envio e pelo recebimento, situação da transferência (`em_transito`, `recebida`) (proposta).
- **Operações:** criar almoxarifados adicionais; transferir entre almoxarifados da mesma unidade; transferir entre unidades (exige RF-ADM-003); confirmar o recebimento.
- **Validações:** origem e destino diferentes e ativos; saldo disponível na origem; transferência entre unidades exige acesso às duas; envio para `em_terceiro` não é transferência, é remessa (RF-EST-015).
- **Regras:** a transferência gera duas movimentações `transferencia` ligadas (saída e entrada), ao custo médio da origem. Em trânsito, o saldo não está disponível em nenhum dos dois almoxarifados (proposta). Ao rebaixar para o Essencial, os almoxarifados adicionais ficam somente leitura até o saldo ser transferido (proposta).

#### RF-EST-013 · Almoxarifado em poder de subcontratada

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** estoque da empresa guardado e consumido por uma subcontratada, controlado como almoxarifado do tipo `em_terceiro`.
- **Dados principais:** almoxarifado `em_terceiro`, fornecedor (`fornecedor_id`, com papel `subcontratada`), endereço onde o material fica, pessoa responsável no fornecedor, itens autorizados, contrato de prestação (RF-TER-006, opcional), valor máximo em poder do terceiro (proposta).
- **Operações:** criar o almoxarifado ao habilitar estoque para a subcontratada; o `prestador` consulta saldo e baixa peças nas OS atribuídas ao seu fornecedor (RF-OS-010); o `almoxarife` acompanha saldo, consumos e divergências; inventariar (RF-EST-005).
- **Validações:** `fornecedor_id` obrigatório, de fornecedor com papel `subcontratada`; o `prestador` só lê e movimenta o almoxarifado `em_terceiro` do seu fornecedor (RNF-037); saída para OS desse almoxarifado só em OS com `executor_tipo = fornecedor` do mesmo fornecedor, ou com autorização do `almoxarife` (proposta); valor acima do máximo gera alerta, não bloqueio (proposta).
- **Regras:** o material continua sendo da empresa: o saldo em poder do terceiro entra no estoque geral e no valor em estoque. O saldo entra por remessa (RF-EST-015) ou por compra entregue direto ao terceiro, com NF-e registrada (proposta). Divergência de inventário gera ajuste com motivo e pode ser marcada para cobrança do fornecedor (só registro). O fornecedor só é inativado com o almoxarifado zerado, por retorno ou ajuste.

#### RF-EST-014 · Saldo geral e por subcontratada

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** visão do estoque somando todos os almoxarifados e separando o que está em poder de cada subcontratada.
- **Dados principais:** por item e almoxarifado: saldo físico, reservado e disponível; estoque geral (soma de todos os almoxarifados); saldo em `quarentena`, em trânsito e em conserto; valor por almoxarifado e por fornecedor; data da última movimentação e do último inventário.
- **Operações:** consultar por item, almoxarifado, tipo de almoxarifado e fornecedor; ver a posição consolidada de cada subcontratada (itens, quantidades, valor); exportar.
- **Regras:** o ponto de pedido e o estoque mínimo de RF-EST-004 consideram por padrão o disponível dos almoxarifados `proprio` e `em_terceiro` da unidade, sem `quarentena` e sem itens em conserto (proposta, configurável por item). O `prestador` vê só a posição do seu fornecedor.

#### RF-EST-015 · Remessa e retorno de peças a subcontratada

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** envio de itens de um almoxarifado próprio para o almoxarifado em poder de uma subcontratada e devolução de sobras ao almoxarifado próprio.
- **Dados principais:** almoxarifado de origem e de destino, fornecedor, itens, quantidades, números de série, custo, NF-e de remessa ou de retorno (RF-EST-019), datas de envio e de recebimento, situação (`em_transito`, `recebida`) (proposta).
- **Operações:** registrar remessa (`remessa_terceiro`); confirmar o recebimento pelo `prestador` ou pelo `almoxarife`; registrar retorno (`retorno_terceiro`); consultar remessas com saldo ainda em poder do terceiro.
- **Validações:** na remessa, destino do tipo `em_terceiro`; no retorno, origem `em_terceiro` e destino `proprio`; saldo disponível na origem; quantidade recebida não maior que a enviada; divergência no recebimento exige motivo; NF-e conforme RF-EST-019; peça `em_transito` não é consumida em OS, nem reservada, antes da confirmação de recebimento pela subcontratada (sumário 3.6).
- **Regras:** remessa e retorno não alteram o custo médio: o material muda de lugar, não de dono. Enquanto `em_transito`, o saldo não está disponível nem na origem nem no destino. Saldo remetido sem consumo nem retorno dentro do prazo de retorno gera alerta; o prazo depende do tratamento fiscal (**P8**) e é parametrizável. As movimentações geram `movimentacao_estoque.registrada`.

#### RF-EST-016 · Peças serializadas

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** controle individual das unidades de um item com `controla_serie` (`peca_serializada`): onde está, em qual ativo foi instalada, quantas vezes foi reparada.
- **Dados principais:** item, número de série, situação (`em_estoque`, `instalada`, `em_transito`, `em_conserto`, `baixada`), almoxarifado ou ativo e posição, recondicionada (sim/não), garantia (fornecedor, início, prazo em dias ou em medidor), custo atual, histórico de instalações (ativo, posição, OS, data, leitura de medidor na instalação e na remoção) e de reparos (ordem de reparo, fornecedor, custo, resultado).
- **Operações:** cadastrar na entrada (manual ou pela NF-e); instalar no ativo pela OS, com posição; remover do ativo pela OS, com destino (almoxarifado `proprio`, `quarentena` ou sucata); consultar o histórico da peça e as peças instaladas num ativo. Rota `/v1/pecas-serializadas`.
- **Validações:** número de série único por item na empresa; a peça está num único lugar (almoxarifado ou ativo); toda movimentação de item com série informa os números de série; peça `baixada` não volta a uso; peça `em_conserto` ou em `quarentena` não é instalada; peça `em_transito` não é instalada nem consumida antes da confirmação de recebimento.
- **Regras:** a peça carrega o custo de entrada e, depois de reparada, o custo do reparo (RF-EST-018). O uso da peça é calculado pelo medidor do ativo entre instalação e remoção, base para vida média e MTBF por componente (proposta). Peça instalada com garantia vigente gera aviso ao abrir OS no ativo (RF-OS-020). Pneu não é peça serializada: tem controle próprio (RF-FRO-002).

#### RF-EST-017 · Itens reparáveis e quarentena

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** peça com defeito retirada de um ativo, que pode voltar ao estoque depois de consertada, guardada num almoxarifado do tipo `quarentena` até a decisão.
- **Dados principais:** item `reparavel`, peça serializada ou quantidade, origem (OS, ativo, data, defeito relatado, sintoma), almoxarifado `quarentena` da unidade, decisão (`enviar_conserto`, `sucatear`, `devolver_ao_estoque`) (proposta), responsável pela decisão (usuário com a permissão `est.decidir_destino`).
- **Operações:** receber em quarentena na remoção pela OS (RF-OS-010); decidir; abrir ordem de reparo a partir da quarentena (RF-EST-018); sucatear com motivo (`baixa_sucata`); devolver ao estoque quando a peça for reaproveitável sem reparo (com motivo).
- **Validações:** só itens `reparavel` entram em quarentena pela remoção; saldo de quarentena não é disponível para OS nem para reserva; a decisão de destino (enviar para conserto, sucatear, devolver ao estoque) exige a permissão `est.decidir_destino`, padrão em `gestor`, `planejador` e `almoxarife` (sumário 3.6); sucateamento e devolução ao estoque exigem motivo.
- **Regras:** a peça em quarentena entra com valor zero no valor em estoque, ou com valor residual parametrizado (proposta); o valor real é definido no retorno do conserto. Peça em quarentena há mais dias que o parâmetro (proposta: 30) gera alerta ao `almoxarife`.

#### RF-EST-018 · Ordem de reparo

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** envio de uma ou mais peças com defeito a um fornecedor para conserto (`ordem_reparo`), com orçamento, prazo, custo, NF-e de remessa e de retorno e garantia do reparo.
- **Dados principais:** número `OR-<AAAA>-<NNNNNN>` (sequencial por empresa, reinicia a cada ano), fornecedor, itens e peças (com números de série), defeito relatado, orçamento (RF-TER-009), custo aprovado, prazo previsto, datas de envio e de retorno, NF-e de remessa e de retorno, `garantia_dias`, situação.
- **Estados:**

  | De | Para | Ação |
  |---|---|---|
  | (criação) | `aberta` | Abrir a partir da quarentena |
  | `aberta` | `enviada` | `enviar`: movimentação `envio_conserto` saindo da quarentena; peça `em_conserto` |
  | `enviada` | `em_orcamento` | `registrar-orcamento` |
  | `em_orcamento` | `aprovada` | `aprovar`, conforme alçada |
  | `aprovada` | `retornada` | `registrar-retorno`: movimentação `retorno_conserto` para almoxarifado `proprio`; peça recondicionada |
  | `enviada`, `em_orcamento` | `irreparavel` | `declarar-irreparavel`: peça vai para sucata (`baixa_sucata`) ou é devolvida sem reparo |
  | `aberta`, `enviada` | `cancelada` | `cancelar`, com motivo |

- **Operações:** as ações acima, pela web e pela rota `/v1/ordens-reparo` (ações `enviar`, `registrar-orcamento`, `aprovar`, `registrar-retorno`, `declarar-irreparavel`, `cancelar`); acompanhar prazo; consultar ordens por fornecedor, item e situação.
- **Validações:** só peças em `quarentena` entram na ordem; fornecedor ativo, com papel `assistencia_tecnica`, `oficina` ou `subcontratada` (proposta); NF-e de remessa para enviar e de retorno para registrar o retorno, quando o parâmetro da empresa exigir (RF-EST-019); orçamento acima de um percentual do custo de reposição do item (proposta: 60%) exige confirmação explícita de quem aprova; cancelamento de ordem `enviada` exige a peça de volta à quarentena.
- **Regras:**
  - Orçamento recusado leva a `irreparavel` ou a nova ordem com outro fornecedor.
  - Ao retornar, a peça entra no almoxarifado como recondicionada, com o custo do reparo como custo de entrada e com a garantia do reparo (`garantia_dias`, RF-OS-020). O efeito do custo do reparo no custo médio do item está em aberto (seção 8).
  - Prazo previsto vencido gera alerta ao `almoxarife`; com contrato de prestação que tenha SLA de conserto, gera também `fornecedor.sla_violado` (RF-TER-006).
  - Eventos `ordem_reparo.enviada` e `ordem_reparo.retornada`.
  - Falha da peça dentro da garantia do reparo sugere nova ordem para o mesmo fornecedor, sem custo (proposta).

#### RF-EST-019 · Registro de NF-e de remessa e retorno

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** registro da NF-e emitida pelo ERP que acompanha remessas e retornos a subcontratadas e envios e retornos de conserto. A Tenaz registra a nota; não a emite.
- **Dados principais:** número (`nfe_numero`), chave de acesso (`nfe_chave`), emitente, destinatário, data de emissão, CFOP (lista parametrizável), XML ou DANFE anexado (opcional), movimentação ou ordem de reparo vinculada.
- **Operações:** informar número e chave na movimentação; importar o XML para preencher e conferir itens e quantidades (proposta); consultar notas de remessa sem retorno correspondente; receber número e chave do ERP pelo conector (RF-INT-008, Enterprise).
- **Validações:** chave com 44 dígitos e dígito verificador válido; a mesma chave não é usada em duas movimentações do mesmo sentido; (NF-e, só no Brasil, **P13**) emitente ou destinatário com CNPJ da empresa ou de uma unidade; itens e quantidades do XML, quando importado, iguais aos da movimentação ou divergência justificada.
- **Regras:** os CFOPs aceitos e a obrigatoriedade da nota por tipo de movimentação dependem de **P8** (referência inicial do sumário: 5.915/5.916 para conserto). Até a validação com o contador, o CFOP é informativo e a obrigatoriedade da nota é parâmetro da empresa (proposta).

---

### 4.9 TER — Terceiros

**Nota de país (sumário 0.8, ADR-013).** O fornecedor aceita identificação de qualquer tipo de qualquer país suportado (`cnpj`, `cpf`, `cuit`, `cuil`, `dni`, `ruc`, `ci`), porque uma empresa brasileira pode ter fornecedor paraguaio (RNF-044). Onde este módulo diz CNPJ ou CPF, leia-se a identificação do tipo e do país do cadastro. Normas, documentos e fiscal por país: **P12** e **P13**.

`fornecedor` é o cadastro único de empresas externas. Uma mesma empresa pode ter vários papéis (`papeis`): `fornecedor_pecas`, `subcontratada` (executa manutenção com equipe própria e pode manter estoque da empresa), `assistencia_tecnica` (autorizada ou especializada por marca), `oficina`, `locador`, `posto_combustivel` e `seguradora`. Toda OS tem **executor** (`executor_tipo`: `interno` ou `fornecedor`, com `fornecedor_id`) e **responsável pelo custo** (`responsavel_custo`: `empresa`, `locatario`, `locador`, `garantia` ou `fornecedor`).

O TER básico (RF-TER-001 a RF-TER-005) entra no MVP, em todos os planos. Contratos com SLA e avaliação (RF-TER-006 e RF-TER-007) são Fase 2, plano Profissional; oficina externa e orçamento (RF-TER-008 e RF-TER-009) são Fase 2, Profissional (proposta), porque a tabela de planos do sumário não os cita. O isolamento do prestador (RNF-037) responde ao risco R4 do sumário.

Uma subcontratada ou assistência técnica pode chegar à OS de três formas: usuário `prestador` dentro da conta da empresa (RF-TER-004), link de uso único (RF-TER-005) ou **conta própria conectada** (módulo SRV, RF-SRV-014 a RF-SRV-019). Na terceira, o `fornecedor` da empresa fica ligado a uma `conexao_empresa`, a OS atribuída a ele é criada também na conta da subcontratada e o andamento volta para a OS da empresa.

#### RF-TER-001 · Fornecedor e papéis

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** cadastro único das empresas externas (`fornecedor`) que vendem peças, executam manutenção, mantêm estoque da empresa, consertam peças ou locam ativos. Substitui o antigo RF-EST-010.
- **Dados principais:** identificação fiscal (`identificacao_tipo` e `identificacao`, de qualquer tipo de qualquer país suportado), razão social, nome fantasia, papéis (um ou mais de `fornecedor_pecas`, `subcontratada`, `assistencia_tecnica`, `oficina`, `locador`, `posto_combustivel`, `seguradora`), categoria livre para outros fornecimentos (recapadora, transportadora) (proposta), marcas de que é autorizado (assistência técnica), contatos, endereço, área de atendimento, condições comerciais, conexão com a conta do fornecedor na Tenaz (`conexao_empresa`, opcional, RF-SRV-014), situação.
- **Operações:** criar, editar, inativar, importar por planilha; incluir e retirar papéis; convidar o fornecedor `subcontratada` ou `assistencia_tecnica` a conectar a conta dele (RF-SRV-014); ver o histórico do fornecedor (OS executadas, ordens de reparo, compras, contratos, avaliações). Rota `/v1/fornecedores`.
- **Validações:** identificação válida na forma do tipo e única na empresa; pelo menos um papel; um papel não é retirado enquanto houver uso que dependa dele: OS não encerrada atribuída, almoxarifado `em_terceiro` com saldo, contrato de locação ou de prestação vigente (proposta).
- **Regras:** fornecedor com pedido, nota, OS, ordem de reparo ou contrato não é excluído, só inativado. Inativar o fornecedor desativa seus usuários `prestador` e revoga os links de execução abertos. Os papéis definem onde o fornecedor aparece: compras (`fornecedor_pecas`), atribuição de OS (`subcontratada`, `assistencia_tecnica`, `oficina`), estoque em terceiro (`subcontratada`), contrato de locação tomado (`locador`), abastecimento (`posto_combustivel`, RF-FRO-007), sinistro e apólice (`seguradora`, RF-FRO-011). Inativar o fornecedor com conexão `ativa` ou `suspensa` encerra a conexão (RF-SRV-019).

#### RF-TER-002 · Catálogo de serviços do fornecedor

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** lista do que cada fornecedor faz (`servico_fornecedor`), usada para escolher o prestador e estimar o custo da OS.
- **Dados principais:** fornecedor, descrição, classes de ativo atendidas, marcas atendidas, área de atendimento (`area_atendimento`), prazo de atendimento (`prazo_atendimento_horas`), prazo de solução (`prazo_solucao_horas`), preço de referência (`preco_referencia`), unidade de preço (`unidade_preco`: `hora`, `visita`, `servico`), vigência do preço, situação.
- **Operações:** criar, editar, inativar, importar por planilha; ao atribuir uma OS, buscar prestadores por classe do ativo, marca e área de atendimento. Rota `/v1/servicos-fornecedor`.
- **Validações:** fornecedor com papel `subcontratada`, `assistencia_tecnica` ou `oficina`; preço maior ou igual a zero; prazo de solução maior ou igual ao de atendimento.
- **Regras:** o preço de referência alimenta o custo previsto da OS executada por fornecedor; o custo real vem do contrato, do orçamento aprovado ou da nota. Os prazos do catálogo são referência; só o contrato de prestação (RF-TER-006) tem SLA cobrável. Mudança de preço não altera OS já programadas.

#### RF-TER-003 · OS executada por fornecedor

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** atribuição de uma OS a um fornecedor, que executa com equipe própria, com registro de quem executou, quanto custou e quem paga.
- **Dados principais:** `executor_tipo = fornecedor`, `fornecedor_id`, serviço do catálogo (opcional), técnicos do fornecedor (pessoas, opcional), forma de registro da execução (usuário `prestador`, link de uso único, conta própria conectada (RF-SRV-016) ou lançamento por usuário interno), custo previsto, custo realizado, número do documento fiscal do serviço (opcional no MVP), `responsavel_custo`.
- **Operações:** atribuir a OS ao fornecedor na programação (RF-OS-007); notificar os usuários `prestador` do fornecedor ou gerar link de uso único (RF-TER-005); acompanhar reconhecimento, início e conclusão; lançar o custo do serviço; liberar ou recusar a liberação.
- **Validações:** fornecedor ativo com papel compatível; a conclusão tem as mesmas exigências da OS interna (tarefas, falha, leitura, os sete campos NR-12); em máquina NR-12, o responsável pela manutenção (campo 7) é a pessoa do fornecedor que executou.
- **Regras:** a OS segue o mesmo ciclo de vida (RF-OS-003). O `prestador` reconhece, inicia, pausa e conclui, mas nunca libera nem cancela. Enquanto o fornecedor não inicia, a empresa pode pausar a OS com `aguardando_terceiro`. Horas do fornecedor podem ser apontadas por pessoa, sem compor mão de obra própria; o custo entra em serviços de terceiros (RF-OS-013). A OS entra no MTTA e no MTTR do ativo e, na Fase 2, no desempenho do fornecedor (RF-TER-007). Com fornecedor conectado, os marcos, apontamentos, componentes, fotos e a conclusão chegam pela sincronização entre contas (RF-SRV-018), e a OS mostra o estado da sincronização (`vinculo_os.estado_sincronizacao`).

#### RF-TER-004 · Perfil prestador e acesso restrito

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** acesso de usuários de um fornecedor à web e ao app com o perfil `prestador`, limitado às OS atribuídas ao seu fornecedor.
- **Dados principais:** usuário com `usuario.fornecedor_id`, pessoa vinculada, perfil `prestador`.
- **Operações:** `administrador` ou `gestor` convida o usuário do fornecedor; o prestador vê as OS atribuídas ao seu fornecedor e as executa na web ou no app, inclusive offline, com apontamento, fotos, peças, falha, leitura de medidor e assinatura; consulta do ativo da OS só os dados técnicos, os documentos técnicos e as OS que o próprio fornecedor executou (proposta); na Fase 2, consulta e movimenta o almoxarifado `em_terceiro` do seu fornecedor (RF-EST-013).
- **Validações:** filtro por fornecedor aplicado no servidor em toda leitura e escrita, inclusive API e sincronização do app (RNF-037); o prestador não vê custos internos, outros fornecedores, outras OS do ativo nem pessoas fora do seu fornecedor; só abre solicitação sobre o ativo de uma OS atribuída a ele (proposta).
- **Regras:** o usuário `prestador` que executa OS conta como técnico cobrado (RF-ADM-009). Quando a OS é reatribuída, o prestador perde o acesso a ela na próxima sincronização; comandos feitos antes da reatribuição são aceitos e os posteriores são recusados com registro (proposta). Um técnico que presta serviço a duas empresas clientes da Tenaz tem um usuário em cada uma, sem nada compartilhado entre elas (RNF-013). Quando a subcontratada tem conta própria (SRV) e está conectada, o caminho preferido é a conexão (RF-SRV-014): o técnico usa só a conta da subcontratada e não precisa de usuário `prestador` no cliente, nem é cobrado nele.

#### RF-TER-005 · Link de execução de uso único

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** link seguro (`link_execucao`) para o prestador ocasional registrar a execução de uma OS sem login e sem cobrança.
- **Dados principais:** OS, fornecedor, destinatário (nome, telefone ou e-mail), validade, situação (`ativo`, `usado`, `expirado`, `revogado`) (proposta), quem gerou; dados enviados: início e fim, descrição do serviço, fotos, peças informadas, leitura de medidor, nome e documento de quem executou, aceite.
- **Operações:** gerar o link numa OS atribuída a fornecedor; enviar por e-mail ou copiar para mensagem; o prestador abre o link no navegador do celular, registra e envia uma única vez; o `planejador` revisa e incorpora o registro à OS; revogar. Gestão do link pela rota autenticada `/v1/links-execucao`. Uso pelo prestador (sumário 3.6):
  1. o link tem o formato `/l/#<token>`: o token fica no **fragmento** da URL, que o navegador não envia ao servidor nem grava em registros de acesso;
  2. a página lê o fragmento e envia o token no corpo de `POST /v1/publico/links-execucao/ativar`, que devolve uma sessão curta com escopo `link`, limitada à OS do link;
  3. com essa sessão, a página usa `GET /v1/publico/links-execucao/atual` (dados mínimos da OS), `POST /v1/publico/links-execucao/atual/anexos` (fotos) e `POST /v1/publico/links-execucao/atual/registrar-execucao` (envio único).
- **Validações:** um link vale para uma OS e é consumido no primeiro envio; validade padrão de 7 dias e máxima de 30 (proposta); token opaco e aleatório, guardado só como hash; o token **nunca** aparece no caminho nem na query de nenhuma rota; a sessão curta expira em 60 minutos sem uso e no envio (proposta); a página mostra só o mínimo da OS (ativo, TAG, localização, descrição, tarefas), sem custos nem histórico (RNF-037); limite de tentativas por origem.
- **Regras:** o envio grava os dados como "informados pelo prestador" e não conclui a OS sozinho: um usuário interno confere e conclui (proposta). Peças informadas não baixam estoque; viram sugestão para o `almoxarife` ou o `planejador`. Leitura de medidor informada entra com origem `manual` e passa por RF-MED-003. O uso do link vai para a auditoria com IP e dispositivo. O link não conta como técnico cobrado e exige conexão.

#### RF-TER-006 · Contratos de prestação com SLA

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** contrato com subcontratada ou assistência técnica (`contrato_prestacao`) que define escopo, cobertura, SLA e preços; base da avaliação de desempenho.
- **Dados principais:** fornecedor, número do contrato, vigência (início, fim, renovação), escopo (texto e serviços do catálogo), ativos ou classes cobertos, unidades, SLA de atendimento (horas até o reconhecimento) e de solução (horas até a conclusão) por prioridade, calendário do SLA (24×7 ou horário comercial), preços (por serviço, hora, visita ou valor mensal), multas contratuais (só registro), anexos, situação. Rota `/v1/contratos-prestacao`.
- **Operações:** criar, editar, renovar, encerrar; vincular ativos e classes; ao atribuir OS de ativo coberto, sugerir o fornecedor do contrato; acompanhar o SLA de cada OS; alertar D-60, D-30 e D-7 antes do fim da vigência (proposta). Com fornecedor conectado, os ativos cobertos pelo contrato vigente definem quais ativos são compartilhados com a conta dele (RF-SRV-015).
- **Validações:** fim da vigência posterior ao início; fornecedor com papel `subcontratada` ou `assistencia_tecnica`; um ativo coberto por no máximo um contrato vigente para o mesmo serviço (proposta).
- **Regras:** o relógio do SLA começa na atribuição ao fornecedor e para nas pausas por motivo atribuível à empresa (`aguardando_aprovacao`, ou `aguardando_peca` quando a peça é fornecida pela empresa) (proposta). SLA estourado gera o evento `fornecedor.sla_violado` e notifica o `gestor`. O preço do contrato prevalece sobre o do catálogo no custo previsto da OS.

#### RF-TER-007 · Avaliação de desempenho do prestador

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** avaliação de cada fornecedor por OS e por período (`avaliacao_fornecedor`), com indicadores calculados e nota dada pela empresa.
- **Dados principais:** fornecedor, OS ou período, SLA de atendimento e de solução cumpridos (sim/não, desvio em horas), retrabalho (vínculo de RF-OS-015 dentro do prazo), nota de 1 a 5 e comentário de quem liberou a OS, custo realizado versus previsto.
- **Operações:** avaliar na liberação da OS (opcional; obrigatório por parâmetro, proposta); ver painel por fornecedor com % de SLA cumprido, taxa de retrabalho, MTTR das OS do fornecedor, nota média e custo; comparar fornecedores do mesmo serviço; exportar. Rota `/v1/avaliacoes-fornecedor`.
- **Validações:** nota entre 1 e 5; uma avaliação por OS, editável por 7 dias (proposta).
- **Regras:** os indicadores mostram como são calculados (RF-IND-002). O `prestador` não vê as avaliações (proposta). A nota média e o % de SLA ordenam a sugestão de fornecedores ao atribuir OS.

#### RF-TER-008 · Oficina externa e assistência técnica

**Fase:** Fase 2 · **Plano mínimo:** Profissional (proposta)

- **Descrição:** manutenção executada fora da empresa, em oficina ou assistência técnica, para veículos e equipamentos, com o fluxo envio → orçamento → aprovação → execução → retorno → nota → garantia. Substitui o antigo RF-FRO-012.
- **Dados principais:** OS (`executor_tipo = fornecedor`), fornecedor (papel `oficina` ou `assistencia_tecnica`), data e leitura de medidor na saída, pessoa que levou o ativo, orçamento (RF-TER-009), aprovação, execução, data e leitura de medidor no retorno, nota fiscal ou NFS-e, garantia do serviço.
- **Operações:** enviar o ativo; registrar orçamento; aprovar, recusar ou pedir novo orçamento; acompanhar a execução; registrar o retorno; lançar a nota (RF-EST-009); comparar orçado com realizado; registrar a garantia do serviço (RF-OS-020).
- **Validações:** aprovação conforme alçada de valor; nota acima do orçamento aprovado em mais de 5% exige nova aprovação (proposta); leitura de medidor no retorno validada por RF-MED-003.
- **Regras:** enquanto o ativo está na oficina, a OS fica `pausada` com motivo `aguardando_terceiro` e o tempo conta como indisponível. O custo da nota entra em serviços de terceiros da OS. A garantia do serviço alimenta RF-OS-020. Com contrato de prestação, o SLA de solução é acompanhado como em RF-TER-006.

#### RF-TER-009 · Orçamento de fornecedor

**Fase:** Fase 2 · **Plano mínimo:** Profissional (proposta)

- **Descrição:** orçamento de oficina, assistência técnica ou conserto (`orcamento_fornecedor`), ligado a uma OS ou a uma ordem de reparo, com aprovação por alçada.
- **Dados principais:** fornecedor, OS ou ordem de reparo, itens (peças e serviços, quantidades, valores), prazo de execução, validade do orçamento, condições, anexo, versão, situação (`recebido`, `aprovado`, `recusado`, `substituido`, `expirado`) (proposta), aprovação (quem, quando, alçada, justificativa). Rota `/v1/orcamentos-fornecedor`.
- **Operações:** registrar o orçamento (digitado ou com o PDF do fornecedor anexado); comparar orçamentos de fornecedores diferentes para o mesmo serviço; aprovar, recusar com motivo ou pedir revisão (nova versão); converter peças do orçamento em compra quando a empresa fornecer as peças (RF-EST-008).
- **Validações:** total igual à soma dos itens; aprovação conforme alçada de valor (a mesma de compras, proposta); orçamento vencido não é aprovado.
- **Regras:** só um orçamento aprovado por OS ou ordem de reparo de cada vez; nova versão substitui a anterior e exige nova aprovação. O orçamento aprovado define o custo previsto; o realizado vem da nota. Aprovação em ordem de reparo leva a ordem a `aprovada` (RF-EST-018).

---

### 4.10 SRV — Gestão de serviços do prestador

**Nota de país (sumário 0.8, ADR-013).** Onde este módulo diz CNPJ ou CPF, leia-se identificação fiscal ou pessoal do tipo e do país do cadastro (RNF-044). A conexão entre contas de países diferentes (cliente brasileiro com subcontratada paraguaia) roda no mesmo banco e não muda tecnicamente; o termo da conexão segue a **P14**. Preço e moeda de cobrança do plano Serviços fora do Brasil: **P15**.

Módulo da **conta própria da subcontratada** ou assistência técnica (sumário 2.2 e 3.5), vendido como *Tenaz Serviços*. A conta é uma `empresa` com o módulo SRV e usa as entidades do restante da Tenaz com o significado do lado de quem presta o serviço. Todo o módulo é **MVP**, plano **Serviços**. O pacote Go é `servicos`.

| Necessidade | Entidade (sumário 3.5) | Reaproveita |
|---|---|---|
| Técnicos | `pessoa` + `usuario` com perfil `tecnico`; `especialidade_tecnico` e `pessoa_especialidade` | RF-ADM-004, RF-ADM-015, RF-CON-005, `capacitacao` |
| Clientes atendidos | `cliente_servico` | Padrão de cadastro de `locatario` (RF-LOC-001) |
| Equipamentos de clientes | `ativo` com `propriedade = de_cliente` e `cliente_servico_id` | RF-ATV-001 a RF-ATV-010, RF-MED |
| Serviço a realizar | `descricao_servico` na OS e itens `servico_os` (referência a `servico_fornecedor` da própria conta) | RF-TER-002 |
| OS de serviço | `ordem_servico` da conta da subcontratada | RF-OS-001 a RF-OS-018 (estados, marcos, execução, assinatura) |
| Componentes | Itens de material da OS com `tipo_material` (`consumido`, `substituido`, `fornecido_pelo_cliente`), `numero_serie_retirado`, `numero_serie_instalado` | RF-OS-010, RF-EST (estoque básico) |
| Horas | `data_inicio`, `data_conclusao` e `apontamento` por técnico | RF-OS-004, RF-OS-009 |
| App do técnico | Tenaz Mobile com perfil `tecnico` | RF-MOB-001 a RF-MOB-006, RF-MOB-010 |
| Conexão com cliente Tenaz | `conexao_empresa` (esquema da plataforma, fora do RLS por empresa) e `vinculo_os` | — (novo; ADR-009 em 5_ARQUITETURA_SOFTWARE.md) |

**Conexão entre contas.** Quando o cliente também usa a Tenaz, ele cadastra a subcontratada como `fornecedor` e a convida; a subcontratada aceita e nasce uma `conexao_empresa` (`convidada` → `ativa` → `suspensa` ou `encerrada`). Cada conta guarda a sua cópia; nenhuma consulta lê dados de duas empresas. Um processo de sincronização com papel de banco próprio grava em cada conta numa transação separada, só com os campos da lista fechada da conexão, com outbox nos dois lados, idempotência e auditoria nos dois lados (RNF-041, RNF-042, risco **R5** do sumário). A **liberação** (encerramento) da OS continua sendo do cliente.

**Planos.** Na conta da subcontratada, todo requisito SRV exige o plano Serviços. Na conta do cliente, convidar, receber OS aberta pela subcontratada e receber o andamento (RF-SRV-014 a RF-SRV-019) ficam disponíveis a partir do **Essencial (proposta)**, junto com o TER básico, porque a tabela do sumário 2.1 não cita a conexão do lado do cliente.

**Sem cliente na Tenaz.** A subcontratada usa o SRV sozinha: cadastra o cliente e os equipamentos, executa as OS e envia o relatório da OS em PDF, com assinatura (RF-SRV-013).

#### RF-SRV-001 · Conta da subcontratada e plano Serviços

**Fase:** MVP · **Plano mínimo:** Serviços

- **Descrição:** criação da conta de uma subcontratada ou assistência técnica (`empresa` com o módulo SRV), roteiro de primeiros passos próprio e cobrança pelo plano Serviços (sumário 2.2).
- **Dados principais:** os de RF-ADM-001 (país, identificação fiscal, razão social, nome fantasia, logo, moeda, idioma padrão, fuso, contato do responsável, unidade inicial); módulo SRV habilitado; plano Serviços e ciclo (mensal ou anual); planos adicionais contratados (Essencial, Profissional ou Enterprise, opcionais); número de técnicos cobrados; desconto de canal vigente (**P4**); conexões ativas.
- **Operações:** criar a conta pelo site escolhendo *Tenaz Serviços*; seguir o roteiro de primeiros passos (cadastrar especialidades e técnicos, cadastrar ou importar clientes e equipamentos, montar o catálogo de serviços, abrir e distribuir a primeira OS, emitir o primeiro relatório, aceitar o primeiro convite de cliente quando houver); habilitar o SRV numa conta já existente da mesma empresa; contratar ou cancelar plano adicional; consultar uso e faturas (RF-ADM-009).
- **Validações:** identificação fiscal válida (RF-ADM-001); o SRV é habilitado na conta existente daquela identificação, nunca numa segunda conta da mesma empresa (proposta, coerente com a identificação única de RF-ADM-001); recursos fora do plano Serviços ficam visíveis com a indicação do plano que os libera.
- **Regras:**
  - A conta nasce com os mesmos doze perfis padrão e a mesma matriz de permissões das demais contas (RF-ADM-005), com o significado do lado de quem presta o serviço (seção 2). Conexões com clientes exigem `srv.gerenciar_conexao`, padrão em `administrador` e `gestor`. Quem tem acesso à conta da subcontratada e a outra empresa escolhe a empresa ao entrar (RF-ADM-004, RF-MOB-001).
  - O plano Serviços inclui (sumário 2.2): SRV completo (técnicos, clientes, equipamentos de clientes, OS de serviço, catálogo de serviços, componentes, app offline, assinatura do cliente na OS, relatório de serviços por cliente, conexão com clientes Tenaz), estoque básico (RF-EST-001 a RF-EST-004, RF-EST-006) e indicadores essenciais (RF-IND-001 a RF-IND-005, RF-IND-007, aplicados à conta). Inclui também o que a OS e o app já trazem: checklists (RF-CHK), documentos com vencimento (RF-CON-005, RF-CON-006), assinatura eletrônica avançada (RF-CON-014) e importação e exportação (RF-INT-003, RF-ADM-013).
  - Não inclui API pública, webhooks e SSO (Profissional) nem planos preventivos (PRE) sobre equipamentos de clientes, que o sumário 2.2 não cita (seção 8).
  - **Técnico cobrado:** usuário ativo com perfil `tecnico` na conta da subcontratada, pelo maior número diário no ciclo; `gestor` e `planejador` não pagam. Preço proposto R$ 59 (anual) e R$ 69 (mensal) por técnico/mês, em reais, para empresa do Brasil; mínimo de 1 técnico (proposta). Preço e mínimo dependem de **P4**; o preço para `ARS` e `PYG` e a moeda de cobrança são a **P15**.
  - **Desconto de canal (proposta do sumário):** subcontratada conectada a pelo menos um cliente Tenaz pagante ganha desconto (ex.: 30%) enquanto a conexão estiver `ativa`; depende de **P4**.
  - No cliente contratante, o técnico da subcontratada com conta própria **não** conta como técnico cobrado (RF-ADM-009).
  - Ativo `de_cliente` não conta como veículo cobrado em nenhuma das duas contas: na subcontratada é equipamento de terceiro, e no cliente já é cobrado pelo plano dele (proposta).

#### RF-SRV-002 · Especialidades de técnico

**Fase:** MVP · **Plano mínimo:** Serviços

- **Descrição:** catálogo por empresa das especialidades dos técnicos (`especialidade_tecnico`), como mecânica diesel, elétrica, refrigeração, hidráulica, empilhadeiras, ar-condicionado.
- **Dados principais:** nome, descrição, classes de ativo relacionadas (opcional), capacitações recomendadas (opcional), situação. Rota `/v1/especialidades-tecnico`.
- **Operações:** criar, editar, inativar, importar por planilha; copiar uma lista inicial sugerida (proposta); ver os técnicos de cada especialidade.
- **Validações:** nome único na empresa; especialidade usada por técnico ou por serviço do catálogo não é excluída, só inativada.
- **Regras:** a especialidade é usada para sugerir técnicos na distribuição (RF-SRV-010), como exigência opcional de um serviço do catálogo (RF-SRV-006) e como agrupamento nos indicadores de horas (RF-SRV-021). A especialidade é dado interno da subcontratada e não é compartilhada com clientes conectados.

#### RF-SRV-003 · Técnicos da subcontratada

**Fase:** MVP · **Plano mínimo:** Serviços

- **Descrição:** cadastro dos técnicos como `pessoa` com `usuario` de perfil `tecnico`, com especialidades, capacitações e documentos.
- **Dados principais:** pessoa (RF-ADM-015: nome, identificação pessoal protegida, matrícula, contato, foto, situação); usuário com perfil `tecnico` (RF-ADM-004); especialidades (`pessoa_especialidade`: especialidade e, opcionalmente, nível `auxiliar`, `pleno` ou `especialista`, proposta); capacitações (`capacitacao`: tipo, data de realização, validade, carga horária, instituição, certificado); documentos pessoais com vencimento (`documento_vencimento`: CNH e outros, RF-CON-005; ASO removido do escopo em 2026-11-13 por decisão do dono do produto); custo por hora (para custo e indicadores internos); região de atendimento (proposta); calendário de trabalho (o da unidade, por padrão).
- **Operações:** criar, editar, inativar; criar o usuário e convidar para o app por e-mail (o convite por telefone depende da **P11**, RF-ADM-004); vincular e desvincular especialidades; cadastrar e renovar capacitações e documentos, com alertas D-60, D-30 e D-7 (RF-CON-006); importar por planilha; ver agenda, OS e horas do técnico; ver a matriz de técnicos por especialidade e capacitação.
- **Validações:** identificação pessoal válida e única na conta, quando informada (RF-ADM-015); técnico precisa de pelo menos uma especialidade para aparecer como sugestão na distribuição (proposta); capacitação com validade posterior à realização.
- **Regras:** técnico ativo com perfil `tecnico` é cobrado (RF-SRV-001); técnico inativado perde o acesso e tem o cache do app limpo (RF-MOB-010). No MVP, capacitação vencida ou ausente gera **aviso** na atribuição, sem bloquear; o bloqueio de atribuição por capacitação (RF-CON-013) chega na Fase 2 (proposta). O cadastro de capacitações do técnico usa a mesma entidade `capacitacao` de RF-CON-012, entregue no MVP só para contas com SRV (proposta). Para o cliente conectado, a conexão compartilha só o nome do técnico e, se o campo estiver liberado, as capacitações exigidas pela OS e sua validade (RF-SRV-014); identificação pessoal, documentos pessoais e custo por hora nunca saem da conta.

#### RF-SRV-004 · Clientes atendidos

**Fase:** MVP · **Plano mínimo:** Serviços

- **Descrição:** cadastro das empresas e pessoas para quem a subcontratada presta serviço (`cliente_servico`).
- **Dados principais:** `tipo_pessoa` (`juridica`, `fisica`), `documento` (identificação fiscal de qualquer tipo suportado: CNPJ ou CPF; CUIT, CUIL ou DNI; RUC ou CI; par `identificacao_tipo` e `identificacao`, nomes fixados na etapa do SRV, ADR-013), `nome`, `contatos` (nome, função, e-mail, telefone, recebe relatório da OS: sim ou não), `enderecos_atendimento` (endereço, coordenadas, responsável no local, horário de acesso, instruções de entrada), `conexao_empresa_id` (opcional, quando o cliente usa a Tenaz e está conectado), prazos de atendimento e de solução por prioridade (proposta), parâmetros do relatório (mostrar preços: sim ou não; envio automático na conclusão: sim ou não) (proposta), assinatura do cliente obrigatória na conclusão (proposta), observações, situação. Rota `/v1/clientes-servico`.
- **Operações:** criar, editar, inativar, importar por planilha; ver equipamentos, OS, relatórios, horas e indicadores do cliente; ver o estado da conexão; ao aceitar um convite, ligar a conexão a um cliente existente ou criar o cliente a partir dos dados do convite (RF-SRV-014).
- **Validações:** documento com dígito verificador válido e único na conta; pelo menos um contato e um endereço de atendimento; cliente com OS não encerrada não é inativado; um cliente tem no máximo uma conexão não encerrada com a conta (proposta).
- **Regras:** cliente com OS ou equipamento não é excluído, só inativado. A identificação de cliente pessoa física é dado pessoal protegido, mascarada na interface (RNF-021). Na conexão, razão social e identificação fiscal podem ser preenchidos com o que o cliente compartilha, sem sobrescrever contatos e endereços locais. Os prazos de atendimento e de solução do cliente são a base do SLA dos indicadores (RF-SRV-021); com contrato de prestação do cliente conectado (Fase 2), prevalecem os do contrato.

#### RF-SRV-005 · Equipamentos de clientes

**Fase:** MVP · **Plano mínimo:** Serviços

- **Descrição:** cadastro dos equipamentos dos clientes atendidos como `ativo` com `propriedade = de_cliente`, com hierarquia, QR code, medidores e documentos técnicos como qualquer ativo.
- **Dados principais:** os de RF-ATV-002 (TAG, nome, tipo, classe, fabricante, modelo, número de série, ano, criticidade, fotos, campos extras da classe); `cliente_servico_id` (obrigatório); endereço de atendimento do cliente onde o equipamento está; garantia do serviço da subcontratada (proposta); na cópia sincronizada de cliente conectado, `conexao_empresa_id` e `ativo_remoto_id` (RF-SRV-015).
- **Operações:** criar, editar, inativar, importar por planilha; gerar e ler QR code (RF-ATV-007); cadastrar medidores e leituras (RF-MED-001, RF-MED-002); anexar documentos técnicos (RF-ATV-008); ver o histórico de serviços do equipamento (RF-ATV-010); mudar o equipamento de endereço de atendimento dentro do mesmo cliente.
- **Validações:** cliente ativo; endereço de atendimento do mesmo cliente; TAG única na conta: na cópia sincronizada, se a TAG do cliente colidir com uma TAG da conta, a cópia recebe um sufixo e a TAG original continua visível (proposta); na cópia sincronizada, os campos vindos do cliente são somente leitura.
- **Regras:** equipamento `de_cliente` não usa `localizacao` da subcontratada: a posição é o endereço de atendimento (proposta). O status operacional do equipamento é informativo na conta da subcontratada: na cópia sincronizada, espelha o do cliente; no equipamento sem conexão, é o informado pelo técnico. Equipamento `de_cliente` não entra nos KPIs de disponibilidade da subcontratada, nem na cobrança de veículos (RF-SRV-001), nem em contrato de locação (RF-LOC-003). Dados locais acrescentados à cópia sincronizada (observações, fotos próprias) não voltam ao cliente (proposta). Com a conexão encerrada, a cópia continua como histórico, sem atualização, e pode ser convertida em equipamento local do mesmo cliente (proposta).

#### RF-SRV-006 · Catálogo de serviços da subcontratada

**Fase:** MVP · **Plano mínimo:** Serviços

- **Descrição:** lista do que a subcontratada vende, usada como itens do serviço a realizar na OS.
- **Dados principais:** `servico_fornecedor` da própria conta (RF-TER-002): código, descrição, especialidade exigida (opcional, RF-SRV-002), classes de ativo e marcas atendidas, `prazo_atendimento_horas`, `prazo_solucao_horas`, `preco_referencia`, `unidade_preco` (`hora`, `visita`, `servico`), vigência do preço, procedimento ou checklist padrão (opcional, RF-PRE-001, RF-CHK-001), situação.
- **Operações:** criar, editar, inativar, importar por planilha; usar como item `servico_os` ao abrir ou planejar a OS (RF-SRV-007); definir preço por cliente (proposta).
- **Validações:** as de RF-TER-002, exceto o papel do fornecedor; código único na conta; preço maior ou igual a zero.
- **Regras:** na conta da subcontratada, o catálogo próprio usa a mesma entidade `servico_fornecedor`, sem fornecedor vinculado (o modelo fica em 2_ARQUITETURA_BANCO_DADOS.md; seção 8). Mudança de preço não altera OS já abertas: o item `servico_os` guarda o preço do momento. Preços do catálogo só são enviados ao cliente conectado se a conexão incluir esse campo (RF-SRV-014).

#### RF-SRV-007 · OS de serviço: abertura com equipamento, técnico e serviço a realizar

**Fase:** MVP · **Plano mínimo:** Serviços

- **Descrição:** abertura de `ordem_servico` na conta da subcontratada, vinculada a um equipamento de cliente e a um técnico, com o serviço a realizar.
- **Dados principais:** os de RF-OS-001 (número da conta, com o prefixo pelo país da empresa como em RF-OS-001, tipo, prioridade, "o ativo parou?", data-limite, anexos); equipamento `de_cliente` (o cliente atendido e o endereço de atendimento vêm dele); técnico responsável e auxiliares; data programada; contato do cliente no local; **serviço a realizar**: `descricao_servico` (texto livre) e itens `servico_os` (referência a `servico_fornecedor` da própria conta, quantidade, preço de referência copiado); tarefas e checklist do procedimento do serviço; vínculo com a OS do cliente conectado (`vinculo_os`, quando houver).
- **Operações:** abrir na web ou no app, inclusive offline, pelo `planejador`, `gestor` ou `tecnico` (falha encontrada em campo); abrir a partir de outra OS (serviço adicional no mesmo equipamento); planejar, programar e atribuir (RF-SRV-010); receber automaticamente a OS de um cliente conectado (RF-SRV-016).
- **Validações:** equipamento `de_cliente` ativo de cliente ativo; serviço a realizar com `descricao_servico` ou pelo menos um item `servico_os`; responsável com perfil `tecnico` ativo, quando informado; quantidade dos itens maior que zero; data-limite não anterior à abertura; em equipamento sincronizado de cliente conectado, a abertura segue RF-SRV-017.
- **Regras:** a OS segue o mesmo ciclo de vida, marcos e regras da OS da empresa (RF-OS-003, RF-OS-004), inclusive a ação `assumir`. A prioridade padrão vem da criticidade do equipamento e, sem ela, é `media` (proposta). OS aberta pela própria subcontratada tem `origem = manual`; a OS criada a partir da OS de um cliente conectado tem `origem = manual` e `vinculo_os.lado_origem = cliente` (proposta; seção 8). Na OS sem vínculo com cliente conectado, a liberação (`os.encerrar`) é feita pela própria subcontratada, depois da assinatura do cliente ou da recusa registrada (RF-SRV-012).

#### RF-SRV-008 · Componentes substituídos e materiais

**Fase:** MVP · **Plano mínimo:** Serviços

- **Descrição:** registro das peças e materiais usados na OS de serviço, distinguindo o que foi consumido, o que foi substituído e o que o cliente forneceu.
- **Dados principais:** itens de material da OS (RF-OS-010) com `tipo_material`:
  - `consumido`: insumo ou peça da subcontratada (item de estoque ou material avulso), quantidade, custo;
  - `substituido`: peça retirada e peça instalada, com descrição, item de estoque ou avulso, posição no equipamento, `numero_serie_retirado` e `numero_serie_instalado` quando houver, defeito observado, destino da peça retirada (`devolvida_ao_cliente`, `descartada`, `levada_para_conserto`) (proposta), fotos;
  - `fornecido_pelo_cliente`: peça entregue pelo cliente, com descrição, quantidade e número de série quando houver, sem custo nem baixa de estoque para a subcontratada.
- **Operações:** registrar na web e no app, inclusive offline; ler o código de barras do item (RF-MOB-004); fotografar peça retirada e instalada; baixar do estoque básico da subcontratada quando for item de estoque (`saida_os`); devolver sobra ao estoque.
- **Validações:** `substituido` exige a descrição da peça retirada e da instalada; número de série obrigatório quando o item controla série ou quando o serviço do catálogo exigir (proposta); quantidade maior que zero; `fornecido_pelo_cliente` não gera movimentação nem custo; saldo e baixa offline seguem RF-OS-010.
- **Regras:** os componentes entram no relatório da OS (RF-SRV-013) e, em máquina NR-12 de cliente conectado, preenchem a sugestão do campo 4 (peças reparadas ou substituídas, RF-CON-002). No cliente conectado, os componentes chegam como "informados pelo prestador" (RF-SRV-018): não baixam o estoque do cliente sozinhos; `fornecido_pelo_cliente` vira sugestão de baixa para o `almoxarife` do cliente, e os números de série permitem atualizar as peças serializadas do cliente na Fase 2 (RF-EST-016), após confirmação (proposta). O custo dos materiais da subcontratada nunca vai ao cliente; o preço só vai se a conexão incluir esse campo.

#### RF-SRV-009 · Hora de início e fim e apontamentos por técnico

**Fase:** MVP · **Plano mínimo:** Serviços

- **Descrição:** registro de quando o serviço começou e terminou e das horas de cada técnico na OS.
- **Dados principais:** `data_inicio`, pausas (início, fim, motivo canônico), `data_conclusao` da OS (RF-OS-004); `apontamento` por técnico: técnico, início, fim, duração, tipo de hora (normal, extra) (proposta), natureza (`execucao` ou `deslocamento`) (proposta), geolocalização do início e do fim (proposta), custo por hora vigente.
- **Operações:** cronômetro automático por técnico ao iniciar, pausar, retomar e concluir no app; cada auxiliar inicia e encerra o próprio apontamento; registrar deslocamento até o cliente (proposta); lançar ou ajustar apontamento na web com justificativa, pelo `planejador` ou `gestor`.
- **Validações:** as de RF-OS-004 e RF-OS-009: fim depois do início, ordem cronológica dos marcos, sem apontamentos sobrepostos do mesmo técnico em OS diferentes, salvo justificativa aprovada; apontamento de execução dentro do intervalo entre `data_inicio` e `data_conclusao`; marcos com divergência de relógio acima de 5 minutos sinalizados (RF-OS-004).
- **Regras:** as horas alimentam o custo de mão de obra da subcontratada, o relatório da OS (horas por técnico) e os indicadores (RF-SRV-021). No cliente conectado, se a conexão incluir os campos, chegam `data_inicio`, pausas, `data_conclusao` e os apontamentos (nome do técnico, início, fim, duração, natureza), que entram no MTTR do ativo do cliente; o custo por hora nunca é compartilhado. Apontamento de OS `encerrada` não é alterado.

#### RF-SRV-010 · Distribuição e reatribuição de tarefas

**Fase:** MVP · **Plano mínimo:** Serviços

- **Descrição:** distribuição das OS entre os técnicos da subcontratada e troca de técnico quando necessário.
- **Dados principais:** OS sem técnico, por cliente, prioridade, data-limite e endereço de atendimento; técnicos com especialidades, capacitações, carga do dia e da semana (HH programadas), OS em andamento; histórico de atribuições da OS (quem atribuiu, de quem para quem, quando, motivo).
- **Operações:** atribuir responsável e auxiliares pela lista, pelo quadro ou pelo calendário (RF-OS-007, RF-OS-017); sugerir técnicos com a especialidade exigida pelos itens do serviço e capacitações válidas, ordenados por carga (proposta); distribuir várias OS de uma vez; reatribuir com motivo; desatribuir; o técnico **assume** uma OS `aberta` ou `planejada` da sua equipe (`POST /v1/ordens-servico/{id}/assumir`), o que grava `data_reconhecimento`.
- **Validações:** responsável e auxiliares com perfil `tecnico` ativos; reatribuição exige motivo; aviso, sem bloqueio no MVP, quando faltar especialidade ou capacitação válida (RF-SRV-003); OS `concluida`, `encerrada` ou `cancelada` não é reatribuída.
- **Regras:** a atribuição e a reatribuição notificam por push os técnicos envolvidos (RF-MOB-006). Reatribuir uma OS `em_execucao` encerra o apontamento aberto do técnico anterior na hora da reatribuição (proposta). A OS reatribuída sai do cache do técnico anterior na sincronização seguinte; comandos que ele fez antes da reatribuição são aceitos e os posteriores são recusados com registro (como em RF-TER-004). Na OS ligada a cliente conectado, o cliente vê só o técnico responsável atual, se o campo estiver na conexão, e não o histórico de distribuição.

#### RF-SRV-011 · App do técnico da subcontratada

**Fase:** MVP · **Plano mínimo:** Serviços

- **Descrição:** uso do app Tenaz Mobile pelo técnico da subcontratada, com o perfil `tecnico` na conta dela. Reaproveita o app do técnico da empresa: acesso (RF-MOB-001), cache e sincronização (RF-MOB-002), operações offline (RF-MOB-003), QR e código de barras (RF-MOB-004), fotos (RF-MOB-005), push (RF-MOB-006) e gestão do aparelho (RF-MOB-010). Este requisito acrescenta só o que é próprio do serviço em cliente.
- **Operações:**
  - **receber tarefas:** lista do dia e da semana com cliente, endereço de atendimento, equipamento, serviço a realizar, prioridade e data-limite; abrir o endereço no aplicativo de mapas do aparelho;
  - **notificações:** OS atribuída, reatribuída, cancelada pelo cliente e liberação recusada (RF-MOB-006);
  - **executar:** iniciar, pausar com motivo, retomar e concluir; tarefas e checklist do serviço (RF-CHK-005); fotos antes e depois; leitura de medidor; falha (RF-OS-011); campos NR-12 quando o equipamento de cliente conectado exigir (RF-CON-002);
  - **componentes e horas:** RF-SRV-008 e RF-SRV-009;
  - **assinatura do cliente:** RF-SRV-012;
  - **abrir OS** para falha encontrada na visita (RF-SRV-007 e RF-SRV-017);
  - **consultar** o equipamento, seus documentos técnicos marcados como offline e as últimas OS da subcontratada nele;
  - tudo isso **offline**, com sincronização posterior.
- **Validações:** concluir exige tarefas obrigatórias resolvidas, checklists obrigatórios concluídos, fotos exigidas e a assinatura do cliente ou a recusa registrada, quando o cliente exigir (RF-SRV-004); em equipamento de cliente conectado, também as exigências da OS do cliente (RF-SRV-018).
- **Regras:** o cache do técnico contém só as OS atribuídas a ele e o necessário para executá-las (RF-MOB-002). O técnico não vê custos; preços dos itens só aparecem se o parâmetro da conta permitir (proposta). Autonomia offline e conflitos seguem RNF-010 e RNF-011; metas de uso em RNF-028. O comando offline sobre OS ligada a cliente conectado é aplicado primeiro na conta da subcontratada e só depois sincronizado com o cliente (RF-SRV-018).

#### RF-SRV-012 · Assinatura do cliente na OS

**Fase:** MVP · **Plano mínimo:** Serviços

- **Descrição:** aceite do serviço pelo representante do cliente atendido, colhido no aparelho do técnico, sem login.
- **Dados principais:** nome, documento (identificação pessoal do país ou matrícula, conforme parâmetro do cliente) e função do assinante; assinatura desenhada na tela; data, hora e geolocalização; resumo exibido no momento (serviço executado, componentes, horas); hash SHA-256 do conteúdo; ou recusa de assinatura com motivo.
- **Operações:** mostrar o resumo ao cliente; colher a assinatura, inclusive offline; registrar a recusa; enviar cópia ao e-mail do assinante junto com o relatório (RF-SRV-013).
- **Validações:** nome obrigatório; documento conforme o parâmetro do cliente; assinatura colhida na conclusão, depois do último apontamento de execução; recusa exige motivo e foto do local (proposta).
- **Regras:** o técnico assina com a assinatura eletrônica avançada (RF-CON-014: usuário autenticado, reautenticação, carimbo de tempo, geolocalização, hash). O cliente assina como aceite com evidências (identificação declarada, imagem da assinatura, geolocalização, carimbo de tempo e hash), sem autenticação, como o prestador ocasional do link (RF-TER-005). Conteúdo assinado não é editado: alterar componentes ou horas depois da assinatura exige nova assinatura (proposta). Os dados do assinante são dados pessoais (RNF-021). No cliente conectado, a evidência da assinatura chega à OS se o campo estiver na conexão, mas não substitui a liberação do cliente (RF-SRV-018).

#### RF-SRV-013 · Relatório da OS em PDF

**Fase:** MVP · **Plano mínimo:** Serviços

- **Descrição:** documento da OS de serviço para o cliente atendido, com serviço, componentes, horas, fotos e assinatura do cliente. Rota `GET /v1/ordens-servico/{id}/relatorio-pdf`.
- **Dados principais:** dados e logo da subcontratada; cliente e endereço de atendimento; equipamento (TAG, nome, fabricante, modelo, número de série); número, tipo e prioridade da OS; serviço a realizar e itens do catálogo, com ou sem preços conforme o parâmetro do cliente; serviço executado e observações; componentes (`consumido`, `substituido` com séries retirada e instalada, `fornecido_pelo_cliente`); `data_inicio`, `data_conclusao` e horas por técnico; fotos selecionadas; checklist; assinaturas do técnico e do cliente (ou a recusa); data de emissão, versão do relatório, hash e QR de verificação.
- **Operações:** pré-visualizar durante a execução (marcado "prévia"); gerar na conclusão ou depois; enviar por e-mail aos contatos do cliente marcados para receber relatório, por link temporário; envio automático na conclusão, conforme o parâmetro do cliente (RF-SRV-004); baixar; reemitir, gerando nova versão com histórico.
- **Validações:** relatório definitivo só de OS `concluida` ou `encerrada`; só usuários da conta da subcontratada geram e enviam; o link enviado expira em 30 dias (proposta).
- **Regras:** usa o mesmo gerador de PDF da OS (RF-OS-018). Nunca traz custos internos; traz preços só quando o parâmetro do cliente permitir. O hash permite verificar a integridade (RF-CON-014). Todo envio fica na auditoria. No cliente conectado, o relatório é anexado à OS do cliente se o campo estiver na conexão. Gerado em p95 < 10 s para OS com até 20 fotos (proposta).

#### RF-SRV-014 · Conexão com cliente Tenaz: convite, aceite e campos compartilhados

**Fase:** MVP · **Plano mínimo:** Serviços (subcontratada); Essencial no cliente (proposta)

- **Descrição:** ligação entre a conta de um cliente que usa a Tenaz e a conta da subcontratada (`conexao_empresa`), com a lista fechada do que cada lado envia.
- **Dados principais:** `empresa_cliente_id`, `empresa_prestadora_id`, `fornecedor_id` (na conta do cliente), `cliente_servico_id` (na conta da subcontratada), `estado` (`convidada`, `ativa`, `suspensa`, `encerrada`), `campos_compartilhados`, `data_convite`, `data_aceite`, `data_encerramento`; e-mail de destino do convite, quem convidou e quem aceitou; ativos ou classes compartilhados (RF-SRV-015). Rota `/v1/conexoes-empresa`, com as ações `convidar`, `aceitar`, `recusar`, `suspender`, `reativar` e `encerrar`.
- **Campos compartilhados (lista fechada; a conexão só liga ou desliga itens dela):**

  | Sentido | Campos | Padrão |
  |---|---|---|
  | Cliente → subcontratada | Ativo: TAG, nome, tipo, classe, fabricante, modelo, número de série, ano, criticidade, status operacional, endereço da unidade e localização em texto | Ligado |
  | Cliente → subcontratada | Ativo: medidores e última leitura válida; documentos técnicos marcados como compartilháveis | Desligado |
  | Cliente → subcontratada | OS atribuída: número, tipo, prioridade, título, descrição, data-limite, "o ativo parou?", tarefas e procedimento, anexos da abertura, contato no local | Ligado |
  | Subcontratada → cliente | Estado e marcos (reconhecimento, início, pausas com motivo, conclusão), técnico responsável (nome), apontamentos (nome, início, fim, duração), componentes, fotos, checklist, falha, campos NR-12, leituras de medidor, observações, assinatura do cliente, relatório PDF | Ligado |
  | Subcontratada → cliente | Capacitações do técnico exigidas pela OS e validade; preço dos itens do serviço | Desligado |
  | Nunca | Custos internos de qualquer lado, ativos não compartilhados, outros clientes ou fornecedores, usuários, identificação pessoal e documentos pessoais, avaliações do fornecedor, especialidades e custo por hora dos técnicos | — |

- **Operações:** no cliente, `administrador` ou `gestor` convida a partir do fornecedor (RF-TER-001), com e-mail, ativos e campos; na subcontratada, `administrador` ou `gestor` recebe a notificação e aceita (ligando a um `cliente_servico` existente ou criando um) ou recusa com motivo; reenviar e cancelar convite pendente; consultar a conexão, a lista de campos e o histórico de mudanças; cada lado liga ou desliga os campos que **envia**; ampliar o que o cliente envia não exige novo aceite, e ampliar o que a subcontratada envia não exige ação do cliente (proposta).
- **Validações:** fornecedor ativo com papel `subcontratada` ou `assistencia_tecnica`; convite expira em 14 dias (proposta); no máximo uma conexão não encerrada por par de empresas; o aceite mostra a identificação fiscal do fornecedor cadastrada pelo cliente e exige confirmação quando ela difere da identificação da conta que aceita (proposta); a "mesma entidade" é conferida pela regra do perfil de país: raiz de 8 posições no CNPJ, CUIT e RUC inteiros (ADR-013, Parte A, item 4); só a conta destinatária aceita ou recusa; campos fora da lista fechada não existem na configuração.
- **Regras:** `conexao_empresa` fica no esquema da plataforma, fora do RLS por empresa; cada conta só vê as conexões em que é parte, por um serviço da plataforma, nunca por consulta direta (RNF-041). O aceite leva a `ativa`, gera o evento `conexao_empresa.ativada` nas duas contas e inicia a carga dos ativos compartilhados (RF-SRV-015). Mudança de campos vale para dados futuros e não apaga o que já foi copiado. Convite, aceite, recusa e mudança de campos geram auditoria nas duas contas.

#### RF-SRV-015 · Sincronização dos ativos cobertos

**Fase:** MVP · **Plano mínimo:** Serviços (subcontratada); Essencial no cliente (proposta)

- **Descrição:** cópia, na conta da subcontratada, dos ativos do cliente cobertos pela conexão, como equipamentos `de_cliente`.
- **Dados principais:** ativos ou classes selecionados na conexão (no MVP) ou cobertos por contrato de prestação vigente com o fornecedor (Fase 2, RF-TER-006); na cópia: `propriedade = de_cliente`, `cliente_servico_id`, `conexao_empresa_id`, `ativo_remoto_id` e só os campos ligados na conexão.
- **Operações:** carga inicial no aceite; atualização incremental quando um campo compartilhado muda no cliente; incluir e retirar ativos da conexão (cliente); propor a inclusão de um ativo ao abrir OS sobre ativo ainda não compartilhado (proposta); ver na subcontratada a data da última atualização de cada cópia.
- **Validações:** só ativos não inativos do cliente; campo fora da lista ligada nunca é copiado; os campos sincronizados são somente leitura na subcontratada.
- **Regras:** o cliente é o dono do ativo. Ativo retirado da conexão ou inativado no cliente fica inativo na subcontratada, com motivo e histórico preservado. Hierarquia: o ativo pai só é copiado se também for compartilhado; caso contrário, a cópia fica sem pai (proposta). A sincronização é idempotente (reaplicar o mesmo evento não muda o resultado) e respeita a meta de atraso de RNF-042.

#### RF-SRV-016 · OS do cliente chega à subcontratada

**Fase:** MVP · **Plano mínimo:** Serviços (subcontratada); Essencial no cliente (proposta)

- **Descrição:** quando o cliente atribui uma OS ao fornecedor conectado (RF-OS-007, RF-TER-003), a Tenaz cria a OS correspondente na conta da subcontratada e liga as duas por `vinculo_os`.
- **Dados principais:** `vinculo_os`: `conexao_empresa_id`, `os_cliente_id`, `os_prestador_id`, `lado_origem = cliente`, `estado_sincronizacao` (`pendente`, `sincronizado`, `erro`); na OS da subcontratada, os campos de OS ligados na conexão, o equipamento `de_cliente` correspondente e o número da OS do cliente.
- **Operações:** criar a OS na subcontratada como `aberta`, sem técnico; notificar `planejador` e `gestor` da subcontratada; reconhecer, distribuir (RF-SRV-010) e executar; devolver ao cliente, com motivo, uma OS ainda não iniciada (proposta), o que retira o fornecedor da OS no cliente e notifica o `planejador` do cliente; acompanhar o estado do vínculo nas duas contas e reprocessar um vínculo em `erro`.
- **Validações:** conexão `ativa`; ativo da OS compartilhado na conexão (sem isso, o cliente é convidado a incluí-lo antes de atribuir); OS do cliente em estado que permite atribuição.
- **Regras:** a numeração é independente em cada conta, e cada OS mostra o número da outra. Mudanças do cliente em campos compartilhados (descrição, prioridade, data-limite, tarefas) chegam à OS da subcontratada. Cancelamento no cliente cancela a OS da subcontratada com motivo "cancelada pelo cliente"; como no ciclo canônico, OS em execução é pausada antes. Reatribuição pelo cliente a outro executor cancela a OS da subcontratada com motivo "reatribuída pelo cliente", preservando os apontamentos já feitos. Falha na criação ou na atualização deixa o vínculo em `erro`, gera `vinculo_os.erro` e é retentada automaticamente; o sucesso gera `vinculo_os.sincronizado`.

#### RF-SRV-017 · OS aberta pela subcontratada nasce no cliente

**Fase:** MVP · **Plano mínimo:** Serviços (subcontratada); Essencial no cliente (proposta)

- **Descrição:** a subcontratada abre uma OS num ativo de cliente conectado (por exemplo, falha encontrada na visita) e a OS nasce também na conta do cliente.
- **Dados principais:** equipamento sincronizado, descrição, fotos, "o ativo parou?", tipo e prioridade sugeridos, serviço a realizar proposto (`descricao_servico` e itens, com preço só se o campo estiver na conexão); `vinculo_os` com `lado_origem = prestador`.
- **Operações:** abrir na web ou no app, inclusive offline (número provisório até a sincronização); acompanhar a triagem do cliente; ser notificado do aceite, da atribuição ou do cancelamento.
- **Validações:** só em equipamento sincronizado de conexão `ativa`; descrição ou mídia obrigatória, como na solicitação (RF-SOL-001).
- **Regras:** no cliente, a OS nasce `aberta`, com `origem = prestador`, executor sugerido igual ao fornecedor da conexão, e segue a triagem normal do cliente (sumário 3.5). O "o ativo parou?" informado pela subcontratada fica como sugestão até a triagem do cliente e só então muda o status do ativo (proposta). Na subcontratada, a OS fica `aberta` e não pode ser iniciada enquanto o cliente não a atribuir ao fornecedor, salvo se a conexão permitir execução imediata de `corretiva_emergencial` (proposta; seção 8). Se o cliente atribuir a outro executor ou cancelar, a OS da subcontratada é cancelada com o motivo informado pelo cliente.

#### RF-SRV-018 · Andamento, conclusão e liberação pelo cliente

**Fase:** MVP · **Plano mínimo:** Serviços (subcontratada); Essencial no cliente (proposta)

- **Descrição:** retorno da execução da subcontratada para a OS do cliente e liberação, que continua sendo do cliente.
- **Dados principais:** o que a conexão liga no sentido subcontratada → cliente (RF-SRV-014): reconhecimento, início, pausas com motivo, retomadas, conclusão, apontamentos, componentes, fotos, checklist, falha, campos NR-12, leituras de medidor, observações, assinatura do cliente e relatório PDF.
- **Operações:** enviar cada mudança da OS da subcontratada à OS do cliente; no cliente, quem tem `os.encerrar` libera ou recusa a liberação com `motivo_recusa_liberacao`; a recusa volta à subcontratada, reabre a OS dela em `em_execucao` com o motivo e notifica o técnico; o encerramento no cliente leva a OS da subcontratada a `encerrada`.
- **Validações:** a subcontratada só conclui a OS ligada se atender também às exigências da OS do cliente que foram compartilhadas (tarefas obrigatórias, leitura de medidor, em máquina NR-12 os sete campos de RF-CON-002); a OS ligada não é encerrada pela subcontratada; conclusão recusada pela validação do cliente deixa o vínculo em `erro` com a mensagem, visível nas duas contas.
- **Regras:**
  - Cada campo tem um único dono: execução (marcos, apontamentos, componentes, fotos, checklist, falha, NR-12) é da subcontratada; planejamento, custo, responsável pelo custo e liberação são do cliente. Não há "última escrita vence" entre contas (RNF-042).
  - No cliente, componentes chegam como "informados pelo prestador" e não baixam estoque sozinhos (como em RF-TER-005); leituras de medidor entram com origem `os`, referência à OS do cliente e validação de RF-MED-003.
  - O custo da OS no cliente segue RF-OS-013 (serviços de terceiros pelo contrato, orçamento ou nota do cliente); o preço da subcontratada só aparece se o campo estiver na conexão.
  - Os marcos recebidos entram no MTTA e no MTTR do ativo do cliente e, na Fase 2, no desempenho do fornecedor (RF-TER-007).
  - O técnico da subcontratada não é usuário do cliente e não conta como técnico cobrado nele (RF-ADM-009).

#### RF-SRV-019 · Suspensão, reativação e encerramento da conexão

**Fase:** MVP · **Plano mínimo:** Serviços (subcontratada); Essencial no cliente (proposta)

- **Descrição:** interrupção temporária ou definitiva da sincronização entre as duas contas, sem apagar histórico.
- **Dados principais:** estado, motivo, quem suspendeu, reativou ou encerrou, `data_encerramento`, OS ligadas não encerradas no momento.
- **Operações:** suspender (qualquer lado, com motivo); reativar (só o lado que suspendeu, proposta); encerrar (qualquer lado, com motivo); ver as OS ligadas afetadas.
- **Validações:** motivo obrigatório; só `administrador` ou `gestor` de uma das contas; `encerrada` é final: nova ligação exige novo convite; só conexão `ativa` é suspensa e só `suspensa` é reativada.
- **Regras:**
  - **Suspensa:** nenhuma OS nova é enviada e nada é gravado na outra conta; as mudanças ficam no outbox de cada lado e são entregues em ordem na reativação (proposta). Cada conta continua trabalhando nas suas cópias.
  - **Encerrada:** a sincronização para de vez e gera `conexao_empresa.encerrada`; nada é apagado: cada conta mantém o histórico do que já recebeu (sumário 3.5). Equipamentos sincronizados ficam como histórico (RF-SRV-005); OS ligadas não encerradas continuam em cada conta, marcadas "conexão encerrada": o cliente decide se conclui com registro interno, reatribui ou cancela; a subcontratada pode concluir localmente e enviar o relatório PDF (RF-SRV-013).
  - Inativar o fornecedor no cliente, inativar o cliente na subcontratada ou encerrar uma das contas encerra a conexão (RF-TER-001, RF-ADM-013).
  - O desconto de canal da subcontratada (RF-SRV-001) acompanha o estado `ativa`.
  - Toda mudança de estado gera auditoria nas duas contas.

#### RF-SRV-020 · Relatório de serviços por cliente

**Fase:** MVP · **Plano mínimo:** Serviços

- **Descrição:** prestação de contas periódica da subcontratada a cada cliente atendido (sumário 2.2).
- **Dados principais:** cliente, período, OS executadas por estado e tipo, equipamentos atendidos, serviços do catálogo, componentes substituídos e fornecidos pelo cliente, horas por técnico, SLA de atendimento e de solução, recusas de liberação, links para os relatórios de cada OS; preços conforme o parâmetro do cliente.
- **Operações:** gerar por cliente e período; exportar em PDF e XLSX; enviar por e-mail aos contatos do cliente; agendar envio mensal (proposta).
- **Validações:** período de até 12 meses por relatório (proposta); só OS da própria conta.
- **Regras:** nunca traz custos internos. Os números seguem as mesmas definições de RF-SRV-021 e mostram como são calculados (RF-IND-002). O envio fica na auditoria.

#### RF-SRV-021 · Indicadores da subcontratada

**Fase:** MVP · **Plano mínimo:** Serviços

- **Descrição:** painel e indicadores da operação de serviços, sobre os dados da própria conta.
- **Dados principais e fórmulas (proposta):**

  | Indicador | Cálculo |
  |---|---|
  | OS por cliente | OS abertas, em andamento, concluídas, encerradas, canceladas e atrasadas (data-limite vencida) por cliente e período |
  | Horas por técnico | Σ apontamentos por técnico, por cliente e por especialidade, separando execução e deslocamento; utilização = horas apontadas / horas disponíveis no calendário do técnico |
  | SLA de atendimento cumprido | OS com `data_reconhecimento` − `data_abertura` ≤ prazo de atendimento / OS do período × 100 |
  | SLA de solução cumprido | OS com `data_conclusao` − `data_abertura`, descontadas pausas atribuíveis ao cliente (`aguardando_aprovacao`, `aguardando_peca` com peça do cliente), ≤ prazo de solução / OS do período × 100 |
  | Tempo médio de atendimento e de execução | MTTA e média de `data_conclusao` − `data_inicio` por cliente |
  | Retrabalho | OS no mesmo equipamento dentro do prazo de retrabalho (RF-OS-015) / OS concluídas × 100 |
  | Recusas de liberação | OS com liberação recusada pelo cliente conectado / OS ligadas concluídas × 100 |
  | Componentes | Peças substituídas por modelo de equipamento e por cliente |

- **Operações:** painel inicial da conta (RF-IND-001); filtrar por cliente, técnico, especialidade, tipo de OS e período; abrir as OS que compõem cada número; exportar.
- **Validações:** OS sem prazo definido ficam fora do SLA e são listadas à parte.
- **Regras:** o prazo de SLA vem do cliente atendido (RF-SRV-004) por prioridade, ou dos prazos do serviço do catálogo quando o cliente não tiver prazo; com contrato de prestação do cliente conectado (Fase 2), vale o SLA do contrato. Todo indicador mostra "como calculamos" (RF-IND-002). Os indicadores usam só dados da conta da subcontratada; a avaliação do fornecedor pelo cliente é RF-TER-007, na conta do cliente.

---

### 4.11 LOC — Locação

**Nota de país (sumário 0.8, ADR-013).** Onde este módulo diz CNPJ ou CPF, leia-se identificação do tipo e do país do cadastro (RNF-044). CFOP e tratamento fiscal da locação são do Brasil (**P8**); o equivalente da Argentina e do Paraguai depende da **P13**.

Contratos de locação nos dois sentidos (sumário 3.1): `cedido` (ativo próprio locado a um `locatario`) e `tomado` (ativo de um fornecedor com papel `locador`, usado pela empresa). Todo o módulo é Fase 2, plano Profissional. A Tenaz controla manutenção, vistoria e responsabilidade de custo; a fatura da locação e a cobrança de avarias e excedentes são emitidas pelo ERP. CFOP e tratamento fiscal da locação dependem de **P8**. O isolamento do locatário está em RNF-038.

#### RF-LOC-001 · Locatários

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** cadastro dos clientes que recebem ativos em locação (`locatario`).
- **Dados principais:** tipo (empresa ou pessoa), identificação (`identificacao_tipo` e `identificacao`, de qualquer tipo de qualquer país suportado), razão social ou nome, contatos (nome, e-mail, telefone, função: financeiro, operação, vistoria), locais de uso (endereço, coordenadas, responsável no local), situação. Rota `/v1/locatarios`.
- **Operações:** criar, editar, inativar, importar por planilha; convidar usuários com perfil `locatario` (RF-LOC-008); ver contratos, ativos locados, OS e avarias do locatário.
- **Validações:** identificação válida na forma do tipo e única na empresa; pelo menos um contato; local de uso com endereço.
- **Regras:** locatário com contrato não é excluído, só inativado. Inativar exige nenhum contrato vigente e desativa os usuários do locatário. A identificação de locatário pessoa física é dado pessoal protegido (RNF-021).

#### RF-LOC-002 · Contrato de locação

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** contrato de locação de um ou mais ativos (`contrato_locacao`), cedido ou tomado.
- **Dados principais:** número `CL-<AAAA>-<NNNNNN>` (sequencial por empresa, reinicia a cada ano), sentido (`cedido`, com `locatario_id`; `tomado`, com `fornecedor_id` de papel `locador`), vigência (início, fim, renovação), local de uso (`local_uso`), franquia por período (`franquia_horas`, `franquia_km`), valor de referência da hora ou km excedente, matriz de responsabilidade (RF-LOC-004), SLA de substituição (contrato tomado, RF-LOC-007), número do contrato no ERP, anexos, situação (`rascunho`, `vigente`, `encerrado`, `cancelado`). Rota `/v1/contratos-locacao`, com ações `ativar`, `encerrar` e `cancelar`.
- **Operações:** criar em rascunho; incluir itens (RF-LOC-003); preencher a matriz; ativar; aditar (incluir ou devolver itens, prorrogar) com histórico (proposta); encerrar, inclusive antes do fim da vigência; cancelar.
- **Validações:** `rascunho` → `vigente` exige pelo menos um item, matriz preenchida e vigência; `vigente` → `encerrado` exige todos os itens devolvidos com vistoria de devolução e, quando antes do fim da vigência, motivo; `cancelado` só a partir de `rascunho`: contrato `vigente` termina sempre por `encerrar` (sumário 3.6); no `cedido`, `locatario_id` obrigatório e sem `fornecedor_id`, e o inverso no `tomado`; um ativo em no máximo um contrato vigente de cada sentido no mesmo período; **sublocação não é permitida**: ativo `de_terceiro` (tomado) não entra em contrato `cedido`.
- **Regras:** a ativação gera `contrato_locacao.ativado` e o encerramento gera `contrato_locacao.encerrado`. Contrato perto do fim da vigência gera alertas D-30 e D-7 (proposta). Contrato `encerrado` ou `cancelado` não é editado.

#### RF-LOC-003 · Itens do contrato

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** ativos do contrato (`item_contrato_locacao`), com datas de entrega e devolução por item.
- **Dados principais:** contrato, ativo, datas previstas e reais de entrega e de devolução, leituras de medidor na entrega e na devolução (das vistorias), franquia do item, local de uso do item quando diferente do contrato, situação (`a_entregar`, `entregue`, `devolvido`) (proposta).
- **Operações:** incluir e retirar itens em rascunho ou por aditivo; registrar entrega (exige vistoria de entrega, RF-LOC-005); registrar devolução (exige vistoria de devolução); substituir item (devolução de um ativo e entrega de outro no mesmo contrato).
- **Validações:** no `cedido`, ativo `proprio` e não `inativo` (sem sublocação de ativo `de_terceiro`; ativo `de_cliente` nunca entra em contrato de locação); no `tomado`, ativo `de_terceiro`; ativo `parado` só é entregue com justificativa do `gestor` (proposta); devolução posterior à entrega.
- **Regras:** ativo `cedido` entregue continua `operando` e sua localização passa a ser o local de uso do locatário até a devolução, quando volta à localização de origem (proposta). Ativo `tomado` é cadastrado com `propriedade = de_terceiro` e, se for `veiculo`, entra na base de veículos cobrados (RF-ADM-009); devolvido ao locador, passa a `inativo` com motivo "devolvido ao locador" (sumário 3.6). Os planos preventivos do ativo continuam ativos durante a locação; a matriz define quem executa e quem paga.

#### RF-LOC-004 · Matriz de responsabilidade

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** linhas do contrato (`responsabilidade_locacao`) que dizem quem paga e quem executa cada tipo de manutenção e definem o `responsavel_custo` padrão das OS dos ativos do contrato.
- **Dados principais:** contrato, categoria (`preventiva`, `corretiva_desgaste`, `avaria_mau_uso`, `pneus`, `deslocamento`, sumário 3.6), responsável pelo custo (`empresa`, `locatario` ou `locador`), executor (`interno` ou `fornecedor`), observação, limite de valor por evento (proposta).
- **Operações:** preencher a partir de um modelo da empresa (proposta); editar em rascunho; alterar no contrato vigente por aditivo, com data de vigência da linha.
- **Validações:** uma linha por categoria; as cinco categorias canônicas são obrigatórias e a lista é fechada (outros casos são tratados na observação da linha mais próxima); no `cedido`, responsável `empresa` ou `locatario`; no `tomado`, responsável `empresa` ou `locador`.
- **Regras:** ao criar OS em ativo de contrato vigente, o sistema escolhe a linha pela categoria derivada do tipo de OS e da causa (`preventiva` → `preventiva`; corretiva com causa de desgaste → `corretiva_desgaste`; avaria por mau uso indicada na OS ou na vistoria → `avaria_mau_uso`; OS de pneu → `pneus`; deslocamento do técnico ou do ativo → `deslocamento`) (proposta) e preenche `responsavel_custo` e o executor; o usuário altera com `justificativa_responsavel_custo`, auditada. Mudança da matriz não altera OS já encerradas.

#### RF-LOC-005 · Vistorias de entrega e devolução

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** checklist de entrega ou de devolução de um ativo locado (`vistoria_locacao`), com leituras de medidor, fotos, avarias e assinatura das duas partes.
- **Dados principais:** contrato, item, tipo (`entrega`, `devolucao`), modelo de checklist (RF-CHK-001), respostas, leituras de medidor, fotos com data, hora e geolocalização, avarias encontradas (descrição, local no ativo, fotos, gravidade), assinatura do representante da empresa e do locatário ou locador (nome e documento), data e hora. Rota `/v1/vistorias-locacao`.
- **Operações:** realizar na web ou no app, inclusive offline; comparar a devolução com a entrega lado a lado; emitir PDF assinado; registrar contestação da outra parte (proposta).
- **Validações:** vistoria de entrega antes de marcar o item como entregue e de devolução antes de marcá-lo como devolvido; leitura de medidor obrigatória quando o ativo tiver medidor, validada por RF-MED-003; assinatura das duas partes, ou recusa de assinatura registrada com testemunha (proposta).
- **Regras:** vistoria concluída não é editada e gera `vistoria_locacao.registrada`. A leitura entra com origem `manual` e referência à vistoria. A assinatura segue RF-CON-014.

#### RF-LOC-006 · Avaria gera OS com custo do locatário

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** abertura de OS corretiva para as avarias encontradas na vistoria de devolução de contrato `cedido`, com custo do locatário.
- **Operações:** gerar OS a partir das avarias da vistoria, uma por avaria ou agrupadas (proposta); ajustar tipo e prioridade; anexar as fotos da entrega e da devolução; registrar contestação do locatário.
- **Validações:** só avaria que não constava na vistoria de entrega gera OS com custo do locatário; OS de avaria nasce com `origem = vistoria_locacao` e `responsavel_custo = locatario`, alterável só com `justificativa_responsavel_custo` e permissão de `gestor`.
- **Regras:** a OS segue o ciclo normal; ao ser encerrada, o custo entra no relatório de cobrança (RF-LOC-009). Avaria durante a locação, aberta por solicitação do locatário ou por inspeção, segue a matriz (avaria por mau uso → `locatario`). Em contrato `tomado`, avaria encontrada na devolução ao locador é registrada para conferência, sem OS automática (proposta).

#### RF-LOC-007 · Ativo tomado: manutenção do locador e SLA de substituição

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** controle da manutenção a cargo do locador e do prazo para o locador substituir ou reparar um ativo tomado que parou.
- **Dados principais:** contrato `tomado`, SLA de substituição (horas até o ativo substituto ou o reparo), contato do locador para chamados, OS a cargo do locador.
- **Operações:** abrir OS com `executor_tipo = fornecedor` (o locador) e `responsavel_custo = locador`, conforme a matriz; notificar o locador por e-mail ou pelos seus usuários `prestador`; acompanhar o prazo; registrar a substituição (devolução do ativo parado e entrega do substituto, RF-LOC-003).
- **Validações:** só em ativo de contrato `tomado` vigente; o substituto entra no mesmo contrato.
- **Regras:** o relógio do SLA de substituição começa quando o ativo fica `parado` com OS a cargo do locador e termina na entrega do substituto ou no encerramento da OS; SLA estourado gera `fornecedor.sla_violado` e alerta ao `gestor`. O tempo parado conta na disponibilidade do ativo e num indicador de disponibilidade por locador (proposta). Preventiva a cargo do locador gera OS para o locador; a empresa confirma a execução (proposta).

#### RF-LOC-008 · Portal do locatário

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** acesso do locatário, com o perfil gratuito `locatario` e escopo `locatario`, para abrir solicitações e acompanhar OS só dos ativos locados a ele. No lançamento do LOC, o locatário usa **só o portal web**, sem o app (sumário 3.6).
- **Operações:** ver os ativos locados (TAG, nome, local de uso, datas do contrato); abrir solicitação sobre um ativo locado, com fotos e "o ativo parou?"; acompanhar o estado das solicitações e o resumo das OS desses ativos (estado, data prevista, conclusão); ver as vistorias assinadas e o relatório de avarias e excedentes do seu contrato (proposta); informar leituras de medidor (proposta, origem `manual`).
- **Validações:** filtro por locatário aplicado no servidor em toda leitura e escrita (RNF-038); só ativos com item entregue em contrato vigente; o locatário não vê custos internos, histórico anterior ao contrato, OS de outros ativos nem dados de outros locatários.
- **Regras:** o perfil `locatario` é gratuito e ilimitado. Ao fim da locação, o acesso ao ativo cessa; vistorias e relatórios do contrato continuam visíveis por 90 dias (proposta). A solicitação do locatário tem origem `portal_locatario` e gera OS com o responsável pelo custo da matriz (RF-SOL-002).

#### RF-LOC-009 · Relatório de cobrança de avarias e horas excedentes

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** relatório por contrato e período com o que o locatário deve além da locação: OS com `responsavel_custo = locatario` e horas ou km acima da franquia.
- **Dados principais:** contrato, período, OS encerradas com custo do locatário (mão de obra, peças, terceiros), uso por item (leituras válidas no início e no fim do período, ou na entrega e na devolução), franquia, excedente, valor de referência do excedente, total.
- **Operações:** gerar por contrato ou por locatário e período; marcar itens como enviados ao ERP para faturamento (proposta); exportar em PDF e XLSX e pela API; ver a memória de cálculo (como em RF-IND-002).
- **Validações:** só OS `encerrada` entram; leituras `suspeita` ou `rejeitada` não entram no cálculo do excedente.
- **Regras:** a fatura é emitida pelo ERP; a Tenaz não fatura. Em contrato `tomado`, o mesmo relatório mostra o que a empresa deve ao locador por excedente e o que o locador deve por OS com `responsavel_custo = locador` (proposta).

---

### 4.12 FRO — Frota

**Nota de país (sumário 0.8, ADR-013).** Os requisitos deste módulo descrevem o Brasil (placa, RENAVAM, CRLV-e, RNTRC e ANTT, CNH, cronotacógrafo, Senatran, SNE). Os documentos de veículo e de condutor da Argentina e do Paraguai (por exemplo, revisão técnica, cédula e habilitação) e as multas desses países dependem da **P12**, e o fiscal, da **P13**. O catálogo de tipos de documento e os prazos viram pacote por país, sem prazo no código (RNF-046). Este documento não reescreve conteúdo legal de outro país antes da pesquisa e do parecer externo.

Os dados cadastrais de veículo e a composição ficam em ATV (RF-ATV-003 e RF-ATV-004). Este módulo trata dos registros próprios da operação de frota.

#### RF-FRO-001 · Motoristas e vínculo com veículo

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** cadastro do motorista e registro de quem conduz cada veículo em cada período.
- **Dados principais:** `pessoa` (RF-ADM-015) que conduz, com ou sem usuário de perfil `motorista` (o condutor eventual não precisa de login), número e categoria da CNH, validade da CNH (documento com vencimento da pessoa, RF-CON-005), vínculo com veículo (início, fim, km no início e no fim).
- **Operações:** cadastrar motorista; vincular e desvincular veículo pela web ou pelo app (o próprio motorista assume o veículo ao ler o QR); registrar condutor eventual sem login; consultar quem conduzia um veículo numa data e hora.
- **Validações:** categoria da CNH compatível com a exigida pelo veículo; um motorista com um veículo trator por vez (proposta); CNH vencida gera aviso no vínculo e, a partir da Fase 2, bloqueio se o documento for marcado como crítico (RF-CON-007).
- **Regras:** o histórico de vínculo é a base para identificar o condutor de uma multa (RF-FRO-010). O exame toxicológico do motorista foi removido do escopo em 2026-11-13 por decisão do dono do produto: a Tenaz não guarda dado de saúde (sumário 3.6). O dossiê do motorista (cadastro, CNH, documentos e vínculos) mostra o mesmo aviso fixo de RF-ADM-015, sem opção de fechar.
- **Critério de aceite do aviso:** na tela do dossiê do motorista, na web e no modo motorista do app (RF-MOB-007), o aviso aparece em toda carga, para qualquer perfil, sem botão de fechar.

#### RF-FRO-002 · Cadastro de pneus

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** controle individual do pneu (`pneu`) como ativo próprio, pelo número de fogo.
- **Dados principais:** número de fogo, marca, modelo, medida, DOT (semana e ano de fabricação), tipo de banda (direcional, tração, misto), vida atual (novo, R1, R2, R3…), sulco inicial de cada vida, pressão recomendada, valor de compra, fornecedor, nota fiscal, situação (`estoque`, `montado`, `em_recapagem`, `em_conserto`, `descartado`) (proposta), localização atual (almoxarifado ou veículo e posição).
- **Operações:** cadastrar individualmente ou em lote a partir da NF-e; consultar histórico do pneu.
- **Validações:** número de fogo único na empresa; medida e tipo compatíveis com a posição ao montar.
- **Regras:** pneu descartado não volta a uso. O histórico do pneu acompanha todas as vidas.

#### RF-FRO-003 · Montagem, desmontagem e rodízio

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** movimentação do pneu entre estoque e posições do veículo, com mapa por eixo.
- **Dados principais:** veículo, posição (padrão de nomenclatura por eixo e lado, como 1E, 1D, 2EI, 2EE, 2DI, 2DE, estepe) (proposta), data e hora, leitura de km do veículo, sulco no evento, motivo (montagem, rodízio, desgaste, avaria, recapagem), OS relacionada.
- **Operações:** montar, desmontar, rodiziar (trocar posições) num único registro; ver mapa visual do veículo com sulco e pressão por posição.
- **Validações:** cada posição tem no máximo um pneu; cada pneu está em no máximo uma posição; km do evento passa por RF-MED-003; a configuração de eixos vem da classe ou do veículo.
- **Regras:** o km do pneu é a soma dos km do veículo enquanto montado. Rodízio não muda a vida do pneu.

#### RF-FRO-004 · Inspeção de sulco e pressão

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** aferição periódica por posição.
- **Dados principais:** data, veículo, km, por posição: sulco em milímetros, pressão, observação, foto.
- **Operações:** inspecionar pelo app sobre o mapa do veículo (RF-MOB-009); importar aferições de dispositivo Bluetooth ou de pneu com sensor quando houver integração (proposta, sem fornecedor definido).
- **Validações:** sulco não maior que o sulco inicial da vida; sulco não maior que o da inspeção anterior na mesma vida (tolerância de medição de 0,5 mm, proposta).
- **Regras:** sulco no limite de retirada para recapagem (parâmetro, proposta de 3 mm) gera solicitação. Sulco no limite legal gera alerta crítico; o valor padrão do limite legal deve ser confirmado na regulamentação de trânsito antes da implementação, porque o estudo não o cita. Pressão fora da faixa gera alerta.

#### RF-FRO-005 · Recapagem, conserto e descarte

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** ciclo do pneu fora do veículo.
- **Operações:** enviar à recapadora ou ao conserto (fornecedor de RF-TER-001, data, custo previsto; pneu não usa a ordem de reparo de RF-EST-018, porque não é item de estoque); receber de volta com nova vida e novo sulco inicial; registrar recusa da carcaça; descartar com motivo e valor residual.
- **Validações:** recapagem incrementa a vida em uma unidade; pneu em recapagem não pode ser montado.
- **Regras:** os custos de recapagem e conserto entram no CPK do pneu.

#### RF-FRO-006 · CPK e vida útil do pneu

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** custo por quilômetro e projeção de vida.
- **Regras:** CPK = (compra + recapagens + consertos − residual) / km do pneu, conforme o estudo. O cálculo é feito por pneu e agregado por marca, modelo, vida e posição. A projeção de km restante usa a taxa de desgaste (mm por 1.000 km) das inspeções. Ver RF-IND-006.

#### RF-FRO-007 · Abastecimento

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** registro de combustível e Arla 32 (`abastecimento`), manual ou importado.
- **Dados principais:** data e hora, veículo, motorista, posto (fornecedor com papel `posto_combustivel`) ou tanque interno, produto (diesel S10, diesel S500, gasolina, etanol, GNV, Arla 32), litros, preço por litro, valor total, leitura de km ou horas, tanque cheio (sim/não), nota ou cupom, origem (manual, cartão, tanque interno).
- **Operações:** lançar pela web ou pelo app do motorista (inclusive offline); importar de cartão combustível (RF-INT-007); calcular consumo.
- **Validações:** litros não maiores que a capacidade do tanque mais 5% (proposta); produto compatível com o veículo; leitura de medidor passa por RF-MED-003; valor total coerente com litros × preço (tolerância de R$ 0,05, proposta).
- **Regras:** o consumo (km/l ou l/h) é calculado entre dois abastecimentos de tanque cheio. A leitura do abastecimento usa origem `cartao_combustivel` quando importada e `manual` quando digitada.

#### RF-FRO-008 · Anomalias de consumo

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** alertas sobre abastecimentos fora do padrão.
- **Regras:** alertas mínimos (proposta): consumo fora da média do veículo ou da classe em mais de 15%; litros acima da capacidade; dois abastecimentos do mesmo veículo em intervalo menor que o parâmetro; km regressivo; produto incompatível; abastecimento sem motorista vinculado. Cada alerta pode ser justificado ou convertido em ocorrência para apuração.

#### RF-FRO-009 · Tanque interno

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** controle de combustível de posto próprio ou comboio.
- **Dados principais:** tanque, produto, capacidade, entradas (com nota), saídas (abastecimentos), medições de régua, encerrante da bomba.
- **Operações:** registrar entradas e medições; conciliar saldo teórico com medido.
- **Regras:** diferença entre saldo teórico e medido acima do parâmetro gera alerta de perda.

#### RF-FRO-010 · Multas

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** registro de infrações de trânsito (`multa`) e controle do prazo de indicação do condutor. Regra do Brasil (Senatran, SNE, NIC); multas da Argentina e do Paraguai dependem da **P12** (nota do módulo FRO).
- **Dados principais:** número do auto de infração, órgão autuador, data, hora e local da infração, código e descrição do enquadramento, gravidade, pontos, valor, data da notificação da autuação, prazo de indicação, veículo, condutor identificado, adesão ao SNE, valor com desconto, vencimento, situação (`recebida`, `condutor_identificado`, `indicada`, `em_recurso`, `paga`, `encerrada`) (proposta), anexos.
- **Operações:** lançar manualmente ou por planilha; identificar o condutor pelo vínculo de RF-FRO-001 na data e hora da infração; gerar termo de indicação para assinatura do condutor; registrar que a indicação foi feita no Portal Senatran; registrar recurso e pagamento; descontar de responsável (só registro, sem integração com folha).
- **Validações:** prazo de indicação = data da notificação + 30 dias; veículo existente pela placa.
- **Regras:** alertas ao `gestor` em D-10, D-5 e D-1 do prazo de indicação (proposta). Sem indicação no prazo, a multa é marcada com risco de NIC (multa por não identificação do condutor). O desconto de 40% do SNE é exibido quando a empresa aderiu. A indicação em si é feita pela empresa no Portal Senatran, com certificado digital; a Tenaz não a envia. A importação automática de multas depende de **P3** (RF-INT-009).

#### RF-FRO-011 · Sinistros

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** registro de acidente, avaria, roubo ou furto envolvendo um ativo (`sinistro`).
- **Dados principais:** data, hora e local, tipo (colisão, tombamento, avaria, incêndio, roubo, furto, outro), ativos envolvidos, condutor, vítimas (sim/não), terceiros, boletim de ocorrência, fotos, seguradora (fornecedor com papel `seguradora`), apólice, número do sinistro, franquia, valores (reparo, franquia, indenização, ressarcimento de terceiro), situação (`aberto`, `em_regulacao`, `em_reparo`, `encerrado`) (proposta).
- **Operações:** abrir pelo app (inclusive offline) ou pela web; gerar OS corretiva de reparo vinculada; registrar regulação e indenização; encerrar.
- **Validações:** pelo menos um ativo; condutor obrigatório quando o ativo for veículo em movimento (proposta).
- **Regras:** perda total leva o ativo a `inativo` (RF-ATV-011) com motivo "sinistro". Custos de OS vinculadas somam no custo do sinistro. Dados de vítimas são dados pessoais com acesso restrito (RNF-021).

#### RF-FRO-012 · Movido para TER

Movido para **RF-TER-008** (oficina externa e assistência técnica) na versão 0.2, com o orçamento em RF-TER-009. O ID não é reutilizado e não entra nas contagens.

---

### 4.13 CON — Conformidade

**Nota de país (sumário 0.8, ADR-013).** Os requisitos deste módulo descrevem o Brasil (NR-12, NR-13, LOTO, Inmetro). O núcleo fala de "conformidade", e a NR-12 é o primeiro modelo de `registro_conformidade`, com os mesmos sete campos e as mesmas regras (ADR-013, Parte C, item 5). Os equivalentes da Argentina e do Paraguai, os prazos por jurisdição (nacional ou subdivisão ISO 3166-2) e os documentos de cada país dependem da **P12** e entram como pacote de país, com o prazo vazio até a empresa preencher ou haver parecer (RNF-046). O bloqueio de energia (RF-CON-004, RF-CON-011) vale nos três países. Nenhum prazo legal fica constante no código, nem o do Brasil.

#### RF-CON-001 · Máquinas NR-12 e inventário

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** marcação das máquinas sujeitas à NR-12 e relatório de inventário atualizado, exigido pelo item 12.18.1 à disposição da fiscalização.
- **Dados principais:** marcação "sujeito à NR-12" no ativo; dados do ativo (TAG, tipo, fabricante, modelo, número de série, ano, localização); campos complementares da classe.
- **Operações:** marcar e desmarcar com motivo; emitir inventário NR-12 em PDF e XLSX, por unidade e localização, com data de emissão.
- **Regras:** o conteúdo mínimo exato do inventário deve ser conferido no texto do item 12.18.1 antes da implementação; o estudo confirma a exigência, não a lista de campos. Desmarcar uma máquina é auditado.

#### RF-CON-002 · Registro de manutenção NR-12 (item 12.11.2)

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** os sete elementos do item 12.11.2 tornam-se campos obrigatórios em toda OS de máquina marcada como sujeita à NR-12.
- **Dados principais (os sete campos):**
  1. intervenções realizadas;
  2. data da manutenção;
  3. serviço realizado;
  4. peças reparadas ou substituídas;
  5. condições de segurança do equipamento;
  6. indicação conclusiva sobre a segurança da máquina: `apto`, `apto com restrição` ou `inapto`;
  7. nome do responsável pela manutenção.
- **Operações:** preencher na conclusão da OS, na web ou no app; reaproveitar dados já registrados (tarefas, peças, datas, técnico) como sugestão editável.
- **Validações:** a OS não passa a `concluida` sem os sete campos; "apto com restrição" exige o texto da restrição; "inapto" exige descrição do risco.
- **Regras:** "inapto" impede que o ativo volte a `operando` e gera solicitação de corretiva vinculada, até que nova OS registre parecer "apto" ou "apto com restrição". O registro assinado não é editável. A forma de assinatura do responsável depende de **P1**.

#### RF-CON-003 · Cronograma de preventivas de segurança (item 12.11.2.2)

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** cronograma de execução das manutenções preventivas ligadas à segurança das máquinas.
- **Operações:** emitir o cronograma a partir dos planos marcados como "ligados à segurança" (RF-PRE-002 e RF-PRE-008), com situação de cada ocorrência (prevista, executada no prazo, executada fora do prazo, atrasada).
- **Regras:** máquina NR-12 sem nenhum plano ligado à segurança aparece como pendência no painel de conformidade.

#### RF-CON-004 · Registro de bloqueio e etiquetagem (item 12.11.3)

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** registro simples de bloqueio e etiquetagem durante a manutenção de máquina NR-12, com horário, data, motivo e responsável, como exige o item 12.11.3.
- **Dados principais:** `bloqueio_energia` da OS: data e hora de aplicação e de retirada, motivo, responsável.
- **Operações:** registrar aplicação e retirada do bloqueio dentro da OS.
- **Validações:** os quatro dados obrigatórios; retirada posterior à aplicação.
- **Regras:** a OS de máquina NR-12 não é concluída com bloqueio aplicado e não retirado. O controle completo por ponto de bloqueio é RF-CON-011. O `bloqueio_energia` é o registro de bloqueio e etiquetagem de uma intervenção e não muda o status do ativo; não se confunde com `bloqueio_ativo`, que deixa o ativo `parado` (sumário 3.6).

#### RF-CON-005 · Documentos com vencimento

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** agenda genérica de documentos com validade ligados a ativo, pessoa ou empresa (`documento_vencimento`).
- **Dados principais:** tipo de documento, vínculo (ativo, pessoa ou empresa), número, emissor, data de emissão, validade, arquivo, documento crítico para operação (sim/não), situação calculada (`valido`, `a_vencer`, `vencido`) (proposta), histórico de renovações.
- **Catálogo inicial de tipos (pacote `BR`; os pacotes da Argentina e do Paraguai dependem da **P12**, RNF-046):** CRLV-e; verificação do cronotacógrafo (Inmetro); CIV; CIPP; RNTRC; apólices de seguro (inclusive os seguros exigidos para o RNTRC pela Resolução ANTT 6.068/2025); laudo ou relatório de inspeção NR-13; CNH; prontuário NR-10; inspeção de SPDA; verificação metrológica (Inmetro ou órgão da RBMLQ, para medidores de velocidade, balanças e instrumentos regulamentados); licenças diversas. A empresa pode criar tipos.
- **Operações:** cadastrar; anexar o arquivo; renovar (a nova versão substitui a vigente e a anterior vai para o histórico); consultar por vencimento; exportar.
- **Validações:** validade posterior à emissão; arquivo obrigatório para tipos marcados como críticos (proposta); tipo de documento com nome ligado a saúde é recusado (regra abaixo).
- **Regras:**
  - **CIV e CIPP:** a validade é sempre digitada a partir do certificado, nunca calculada. O cálculo sugerido só será oferecido depois da confirmação da tabela oficial (**P2**).
  - **Cronotacógrafo:** o sistema sugere validade de 2 anos a partir da verificação (Portaria Inmetro 481/2021, subitem 6.3.1), editável.
  - **CRLV-e:** validade digitada; o licenciamento anual (CTB, art. 105) não foi re-lido no estudo.
  - **NR-13 no Essencial:** o laudo é tratado como documento com validade digitada; o cálculo dos prazos é RF-CON-009 (Profissional).
  - Prazos de NR-10 e SPDA são digitados; o estudo não re-leu essas normas.
  - **Verificação metrológica:** validade digitada a partir do certificado até a conferência da norma de cada instrumento (**P9**). O efeito de bloqueio vem de uma regra configurável (RF-ADM-016), não do tipo de documento.
  - **Tipo de documento ligado a saúde é recusado** (sumário 3.6, "Dado de saúde recebido por engano"). Vale ao criar e ao renomear um tipo, pela tela, pela API, pela importação por planilha e pela aplicação de pacote de segmento (RF-ADM-017), e é conferida no servidor.
    - **Comparação (proposta):** o nome e o código do tipo são normalizados (sem acento, em minúsculas) e quebrados em palavras por qualquer caractere que não seja letra ou número; há recusa quando um termo de uma palavra é igual a uma palavra inteira, ou quando um termo de duas palavras aparece como duas palavras seguidas. Parte de palavra não conta ("medidor" não casa com "medico"; "Exames médicos" casa com "exames medicos", mas "Examinador" não casa com nada).
    - **Termos (proposta):** `aso`, `atestado`, `atestados`, `toxicologico`, `toxicologica`, `saude`, `medico`, `medica`, `medicos`, `medicas`, `audiometria`, `audiometrico`, `espirometria`, `hemograma`, `vacina`, `vacinacao`, `pcmso`, `cid`, `exame medico`, `exames medicos`, `laudo medico`, `prontuario medico`, `ficha medica`. "Exame", "clínica", "laudo" e "prontuário" sozinhos não são termos, para não recusar nomes legítimos como o exame de qualificação de soldador, o alvará da clínica, o laudo NR-13 e o prontuário NR-10 (decisão do dono do produto). O nome do arquivo anexado não é verificado; contra ele vale só o aviso fixo (decisão do dono do produto). A lista só muda com nova versão deste documento.
    - **Erro (pedido ao documento 3):** `422 tipo_documento_saude_vedado`, com o termo encontrado no `detail` e o aviso de que a Tenaz não guarda dado de saúde. Na importação e na aplicação de pacote, a linha ou o item é recusado com o mesmo código, e o restante segue a regra de cada um.
    - **Critério de aceite:** teste automatizado com os nomes "ASO", "Exame toxicológico", "ATESTADO médico", "Ficha Médica" e "Saúde ocupacional" recusados com `422 tipo_documento_saude_vedado`; e "Laudo NR-13", "Prontuário NR-10", "Verificação do medidor", "Examinador credenciado", "Exame de qualificação de soldador" e "Alvará da clínica" aceitos; o mesmo resultado pela tela, pela API, pela importação e pelo pacote; renomear um tipo aceito para um nome vedado também é recusado.
    - **Tipos que já existem** (decisão do dono do produto): quando a regra entrar (etapa E08), os tipos de documento já cadastrados com nome vedado não são inativados nem alterados; ficam listados para revisão do `administrador`, com o termo encontrado. Aceite: um tipo vedado criado antes da regra continua ativo e aparece na lista de revisão.

#### RF-CON-006 · Alertas de vencimento D-60, D-30 e D-7

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** avisos antes e depois do vencimento de documentos e capacitações.
- **Dados principais:** antecedências (padrão 60, 30 e 7 dias, configuráveis por tipo), destinatários (responsável do ativo, `gestor`, a própria pessoa no caso de documento pessoal), canais.
- **Operações:** receber alertas por e-mail, push e central de notificações; ver painel de vencimentos (vencidos, 7, 30, 60 dias); confirmar ciência.
- **Regras:** alerta no dia do vencimento e, depois, aviso semanal enquanto vencido (proposta). Os alertas são gerados por rotina diária idempotente e não se repetem para o mesmo marco. A renovação encerra os alertas pendentes do documento anterior.

#### RF-CON-007 · Aptidão do ativo e bloqueio de despacho

**Fase:** Fase 2 · **Plano mínimo:** Profissional (proposta)

- **Descrição:** situação "apto" ou "inapto" para operar, calculada a partir de documentos críticos, OS com parada, parecer NR-12, `bloqueio_ativo` vigente e checklist crítico reprovado. "Apto para despacho" é uma verificação calculada a cada consulta, não uma entidade gravada (sumário 3.6).
- **Operações:** consultar a aptidão na web, no app e pela API e webhook, para que um TMS ou sistema de despacho bloqueie o veículo; liberar exceção com justificativa e prazo, por `gestor` (proposta).
- **Regras:** veículo com documento crítico vencido fica inapto para despacho e a situação é exposta por API, conforme recomendação do estudo; o status do ativo não muda por isso, salvo se houver `bloqueio_ativo`. A exceção é auditada e expira sozinha. No MVP, a situação aparece como aviso, sem API de bloqueio.

#### RF-CON-008 · Cadastro e enquadramento NR-13

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** cadastro dos equipamentos sujeitos à NR-13 com os dados que definem prazos.
- **Dados principais:** tipo (caldeira, vaso de pressão, tubulação, tanque metálico), categoria (caldeira A ou B; vaso I a V), existência de SPIE, condições especiais (caldeira de recuperação de álcalis; teste anual das válvulas; gerenciamento de combustão do Anexo IV; sistema instrumentado de segurança), equipamentos interligados (para tubulações), PLH responsável, prontuário (documentos anexados).
- **Operações:** cadastrar no ativo; anexar prontuário; ligar tubulações ao vaso ou caldeira mais crítico.
- **Validações:** categoria obrigatória para caldeira e vaso.

#### RF-CON-009 · Motor de prazos NR-13

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** cálculo dos prazos máximos de inspeção pela Portaria MTP 1.846/2022 e dos prazos adotados pela empresa.
- **Regras:**
  - A tabela de prazos máximos é **parametrizável** e vem pré-carregada com a tabela do estudo (caldeiras, vasos por categoria com e sem SPIE, inspeções externa e interna).
  - Os prazos de vasos das categorias III a V com SPIE ficam marcados "a confirmar" até a conferência no DOU (**P6**).
  - O prazo adotado pode ser menor que o máximo, definido pelo PLH, nunca maior.
  - Tubulações: prazo até o maior prazo interno do vaso ou caldeira mais crítico interligado, podendo ser dobrado com justificativa, até 10 anos.
  - Tanques metálicos: sem prazo fixo; o prazo é o do programa do responsável técnico, digitado.
  - Cada vencimento gera plano `inspecao_legal` e alertas D-60, D-30 e D-7.
  - Mudança na tabela por atualização da norma é versionada; prazos já calculados mostram a versão usada.

#### RF-CON-010 · Registro de Segurança e relatório de inspeção NR-13

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** livro de Registro de Segurança eletrônico e controle do prazo de emissão do relatório.
- **Dados principais:** ocorrências e inspeções registradas, nome e assinatura do PLH, data da inspeção, data de emissão do relatório, parada geral (sim/não).
- **Operações:** lançar ocorrências; registrar inspeção; anexar relatório; emitir o livro em PDF.
- **Validações:** relatório emitido em até 60 dias após a inspeção, ou 90 dias em parada geral; o sistema alerta antes do fim do prazo.
- **Regras:** registro assinado não é editável. A forma de assinatura do PLH depende de **P1** (RF-CON-014 e RF-CON-015).

#### RF-CON-011 · Bloqueio e etiquetagem por ponto (LOTO)

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** controle completo de bloqueio de energias por ponto de bloqueio.
- **Dados principais:** pontos de bloqueio do ativo (tipo de energia: elétrica, pneumática, hidráulica, mecânica, térmica, química, gravitacional), procedimento de bloqueio, foto de cada ponto, cadeado e etiqueta usados, responsável por ponto, verificação de energia zero, data e hora de aplicação e retirada.
- **Operações:** aplicar bloqueios pelo app ponto a ponto; verificar energia zero; retirar; registrar retirada excepcional (cadeado de pessoa ausente) com aprovação do `gestor` e justificativa.
- **Validações:** OS com LOTO exigido não entra em `em_execucao` antes de todos os pontos bloqueados e verificados; não é concluída com cadeado aplicado; cada cadeado é retirado por quem aplicou, salvo retirada excepcional.
- **Regras:** o registro atende e amplia RF-CON-004: cada intervenção tem um `bloqueio_energia` com os pontos (`ponto_bloqueio`) bloqueados. O estudo não encontrou módulo forte de LOTO em nenhum CMMS do mid-market, o que torna este requisito um diferencial. Permissão de trabalho formal fica fora do escopo desta versão.

#### RF-CON-012 · Capacitações

**Fase:** Fase 2 · **Plano mínimo:** Essencial (proposta)

- **Descrição:** treinamentos e habilitações de pessoas com validade (`capacitacao`), como NR-10, NR-11, NR-12, NR-35 e MOPP.
- **Dados principais:** `pessoa` (RF-ADM-015; inclusive técnico de fornecedor e técnico da conta de uma subcontratada, RF-SRV-003), tipo, data de realização, validade, carga horária, instituição, instrutor, certificado anexado.
- **Operações:** cadastrar, renovar, importar por planilha; receber alertas D-60, D-30 e D-7 (RF-CON-006); ver matriz de pessoas por capacitação.
- **Regras:** as validades típicas (reciclagem bienal da NR-10, cartão anual do operador de empilhadeira da NR-11) podem ser sugeridas, mas o estudo não re-leu essas normas; a validade é digitada a partir do certificado até essa conferência.

#### RF-CON-013 · Capacitação vencida bloqueia atribuição

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** nenhuma OS regulada é atribuída a técnico sem a capacitação exigida válida.
- **Dados principais:** capacitações exigidas na classe do ativo, no procedimento ou no plano; data programada da OS.
- **Validações:** atribuição, programação e início bloqueados se o técnico não tiver todas as capacitações exigidas válidas na data programada e na data de início. Em OS atribuída a fornecedor, a regra vale para as pessoas do fornecedor informadas na OS e para o usuário `prestador` que a inicia.
- **Regras:** não há exceção pela interface (regra do estudo). No app offline, o bloqueio usa as capacitações em cache; o servidor revalida na sincronização e sinaliza ao `gestor` qualquer início feito com capacitação que venceu durante o período offline. OS "regulada" é a que exige ao menos uma capacitação.

#### RF-CON-014 · Assinatura eletrônica avançada

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** assinatura eletrônica própria para OS, checklists, pareceres NR-12, vistorias de locação (Fase 2) e registros NR-13 (Fase 2). Está no MVP e em todos os planos (sumário 2.1 e 4).
- **Dados principais:** usuário autenticado (com MFA quando exigido), reautenticação no momento da assinatura, geolocalização, carimbo de tempo, hash SHA-256 do conteúdo, evento na trilha de auditoria, manifesto de verificação.
- **Operações:** assinar; verificar a integridade de um documento assinado pelo hash; emitir PDF com página de evidências.
- **Regras:** apoia-se na MP 2.200-2, art. 10, §2º, que admite meios de comprovação de autoria e integridade fora da ICP-Brasil quando aceitos pelas partes (citada no estudo, não re-lida). Se o parecer de **P1** concluir que documentos de SST exigem ICP-Brasil, esses documentos passam a exigir RF-CON-015.

#### RF-CON-015 · Assinatura ICP-Brasil

**Fase:** Fase 3 · **Plano mínimo:** Profissional (proposta)

- **Descrição:** assinatura com certificado ICP-Brasil (A1, A3 ou em nuvem) para laudos de PLH e documentos de SST.
- **Operações:** assinar com certificado do usuário; validar cadeia e revogação; anexar a assinatura ao PDF.
- **Regras:** condicionado à resposta de **P1**. Se P1 exigir ICP-Brasil antes da Fase 3, o requisito sobe de fase (ver roadmap).

#### RF-CON-016 · Dossiê de fiscalização

**Fase:** MVP · **Plano mínimo:** Essencial (proposta)

- **Descrição:** pacote de evidências por máquina ou por unidade para apresentar à fiscalização ou à auditoria.
- **Operações:** gerar PDF com inventário NR-12, registros de manutenção do período com os sete campos, cronograma de preventivas de segurança, registros de bloqueio e documentos vigentes; na Fase 2, prontuário e Registro de Segurança NR-13.
- **Regras:** o dossiê traz data de emissão, emissor e hash de cada documento incluído.

---
### 4.14 IND — Indicadores e relatórios

#### Tabela de KPIs e "como calculamos"

As fórmulas vêm do estudo de mercado. A coluna "Dados na Tenaz" diz de onde sai cada termo; as escolhas dessa coluna são propostas de implementação. Toda tela de KPI mostra esta explicação (RF-IND-002).

| KPI | Fórmula | Dados na Tenaz | Fase |
|---|---|---|---|
| MTBF | Tempo de operação / nº de falhas com parada | Tempo de operação = tempo requerido − tempo indisponível; falha = OS `corretiva_planejada` ou `corretiva_emergencial` com "o ativo parou?" = sim | MVP |
| MKBF | km rodados / nº de falhas com parada | Incremento do medidor `km` acumulado no período | MVP |
| MTTR | Σ tempo de reparo / nº de reparos | Tempo de reparo = `data_conclusao` − `data_inicio` − pausas; parâmetro inclui ou exclui pausas `aguardando_peca`, `aguardando_aprovacao` e `aguardando_terceiro` | MVP |
| MTTA | Σ (início do atendimento − abertura) / nº de chamados | `data_reconhecimento` − `data_abertura` das OS corretivas | MVP |
| Disponibilidade operacional (Ao) | (Tempo requerido − tempo indisponível total) / tempo requerido | Indisponível = tempo em `parado` ou `em_manutencao`, inclusive espera de peça e aprovação; tempo requerido pelo calendário do ativo ou da unidade | MVP |
| Disponibilidade inerente (Ai) | MTBF / (MTBF + MTTR) | MTTR sem esperas | MVP |
| Disponibilidade física | (Horas calendário − horas em manutenção) / horas calendário × 100 | Horas calendário 24×7; usada em frota e equipamentos móveis | MVP |
| Backlog (semanas) | Σ HH estimadas pendentes / HH disponíveis por semana | OS não encerradas nem canceladas; HH disponíveis dos técnicos ativos pelo calendário | MVP |
| PM compliance | Preventivas concluídas na janela / preventivas programadas × 100 | Janela = tolerância do plano (RF-PRE-006); ocorrências perdidas contam como não cumpridas | MVP |
| % planejada | HH de preventiva, preditiva e planejada / HH totais × 100 | Numerador: `preventiva`, `preditiva`, `corretiva_planejada` e `inspecao_legal` (proposta); denominador: todas as HH apontadas | MVP |
| First-time fix | OS sem retorno em X dias / OS concluídas × 100 | Retorno = vínculo de retrabalho (RF-OS-015); X = parâmetro (30 dias) | MVP |
| Custo / RAV | Custo anual de manutenção / valor de reposição × 100 | RAV do ativo (RF-ATV-002); ativos sem RAV ficam fora e são listados | MVP |
| CMF | Custo de manutenção / faturamento bruto × 100 | Faturamento informado em parâmetro (RF-ADM-010) | MVP |
| Custo por km | (Manutenção + pneus + combustível + Arla + lubrificantes) / km | No MVP, só manutenção; completo com os módulos de frota | MVP parcial; completo na Fase 2 |
| km/l | km entre abastecimentos de tanque cheio / litros | RF-FRO-007 | Fase 2 |
| CPK do pneu | (Compra + recapagens + consertos − residual) / km do pneu | RF-FRO-006 | Fase 2 |
| Confiabilidade R(t) | e^(−t/MTBF) ou Weibull e^(−(t/η)^β) | Weibull por modo de falha, com mínimo de falhas para ajuste | Fase 3 |
| Schedule compliance | HH executadas da programação / HH programadas × 100 | Programação congelada (RF-OS-019) | Fase 3 |
| Wrench time | Tempo em execução direta / tempo pago × 100 | Apontamentos e jornada do técnico | Fase 3 |

Referências exibidas ao lado dos KPIs, sempre com a fonte: CMF até 4% do faturamento (ABRAMAN 2022); PM compliance ≥ 90%, trabalho planejado ≥ 85%, wrench time de 50% a 55% e custo/RAV de 1% a 3% como "referência de mercado"; backlog de 2 a 4 semanas como "referência não verificada". OEE fica fora do escopo, porque exige apontamento de produção.

#### RF-IND-001 · Painel inicial por perfil

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** primeira tela com o que cada perfil precisa agir hoje.
- **Operações:** `gestor`: ativos parados, disponibilidade, custo do mês, conformidade; `planejador`: solicitações em triagem, OS sem programação, preventivas vencidas, backlog, OS com fornecedor aguardando atendimento; `tecnico`: OS do dia; `prestador`: OS atribuídas ao seu fornecedor; `almoxarife`: itens abaixo do ponto de pedido, reservas pendentes e, na Fase 2, peças em quarentena e ordens de reparo atrasadas; `motorista`: veículo atual, checklist pendente; `locatario` (Fase 2): ativos locados e solicitações abertas. Na conta de uma subcontratada (SRV): `gestor` e `planejador` veem OS sem técnico, OS atrasadas por cliente, OS recebidas de clientes conectados aguardando distribuição, liberações recusadas e erros de sincronização (RF-SRV-021).
- **Regras:** cada número do painel abre a lista que o compõe.

#### RF-IND-002 · "Como calculamos"

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** explicação de cada KPI dentro do produto.
- **Operações:** abrir a explicação de qualquer KPI; ver fórmula, dados usados, filtros aplicados, parâmetros vigentes no período (por exemplo, MTTR com ou sem espera) e lista de registros que entraram no cálculo; ver a referência e sua fonte.
- **Regras:** o valor exibido é reproduzível a partir da lista de registros. Referências sem fonte verificável não são exibidas.

#### RF-IND-003 · KPIs de confiabilidade e disponibilidade

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** MTBF, MKBF, MTTR, MTTA, Ao, Ai e disponibilidade física, conforme a tabela acima.
- **Operações:** filtrar por ativo, classe, criticidade, localização, unidade e período; comparar períodos; ver Pareto de falhas por modo, causa e ativo; ver a meta de disponibilidade por criticidade.
- **Regras:** ativos `inativo` saem do denominador durante o período inativo. Leituras `suspeitas` não entram no MKBF.

#### RF-IND-004 · KPIs de gestão do trabalho

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** backlog em semanas, PM compliance, % planejada e first-time fix.
- **Operações:** filtrar como em RF-IND-003; ver evolução mensal; ver OS que compõem cada número.
- **Regras:** o tipo de OS é obrigatório (RF-OS-002), o que garante o cálculo de % planejada.

#### RF-IND-005 · KPIs de custo

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** custo por ativo, classe, centro de custo, tipo de OS, categoria (mão de obra, peças, terceiros), executor (interno ou fornecedor) e responsável pelo custo, custo/RAV, CMF e custo de manutenção por km.
- **Regras:** só OS `encerrada` entram no custo consolidado; OS em andamento aparecem como custo parcial, separado. Cada visão mostra o custo total e o custo líquido da empresa (sem OS de responsabilidade de `locatario`, `locador`, `garantia` ou `fornecedor`); custo/RAV e CMF usam o custo líquido (proposta).

#### RF-IND-006 · KPIs de frota

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** custo por km completo, km/l, CPK e ranking de veículos, motoristas, marcas de pneu e postos.
- **Regras:** veículos sem leituras válidas suficientes no período aparecem sinalizados e fora das médias. O estudo não traz benchmarks brasileiros de custo por km, km/l e CPK; o sistema não exibe referência para eles.

#### RF-IND-007 · Relatórios padrão e exportação

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** relatórios prontos e exportação de listas.
- **Operações:** OS por período, histórico do ativo, preventivas vencidas e previstas, consumo de peças, movimentações de estoque, documentos a vencer, inventário NR-12, leituras suspeitas, checklists realizados; exportar em CSV, XLSX e PDF.
- **Regras:** toda exportação respeita permissões e é registrada na auditoria.

#### RF-IND-008 · Relatórios agendados e personalizados

**Fase:** Fase 2 · **Plano mínimo:** Profissional (proposta)

- **Descrição:** relatórios montados pelo usuário e enviados por e-mail em horário definido.
- **Operações:** escolher entidade, colunas, filtros e agrupamentos; salvar; agendar envio diário, semanal ou mensal.
- **Regras:** o destinatário precisa ser usuário da empresa com permissão para os dados do relatório.

#### RF-IND-009 · Visão consolidada multi-unidade

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** KPIs e custos lado a lado por unidade.
- **Regras:** usuário com acesso parcial vê só as suas unidades no consolidado.

#### RF-IND-010 · Confiabilidade avançada

**Fase:** Fase 3 · **Plano mínimo:** Profissional (proposta)

- **Descrição:** R(t) exponencial e Weibull por modo de falha e classe.
- **Regras:** o ajuste de Weibull só é exibido com número mínimo de falhas (proposta: 5) e mostra o intervalo de confiança.

#### RF-IND-011 · Schedule compliance e wrench time

**Fase:** Fase 3 · **Plano mínimo:** Profissional (proposta)

- **Descrição:** aderência à programação congelada e proporção do tempo em execução direta.
- **Regras:** depende de RF-OS-019 (programação) e da jornada dos técnicos.

#### RF-IND-012 · TCO e renovação de frota e máquinas

**Fase:** Fase 3 · **Plano mínimo:** Enterprise

- **Descrição:** custo total de propriedade e ponto ótimo de substituição.
- **Dados principais:** aquisição, manutenção, combustível, pneus, residual estimado, depreciação de referência.
- **Operações:** ver custo acumulado e custo anual equivalente por ativo; simular a substituição; comparar com ativo novo.
- **Regras:** as taxas de depreciação de referência (caminhões 25% ao ano e máquinas 10% ao ano, IN RFB 1.700/2017, citada no estudo e não re-lida) são parametrizáveis e servem só para análise gerencial.

---

### 4.15 INT — Integrações

#### RF-INT-001 · API REST pública v1

**Fase:** MVP · **Plano mínimo:** Profissional

- **Descrição:** API de leitura e escrita em `/v1/<recurso-no-plural>`, descrita em OpenAPI 3.1.
- **Dados principais:** recursos mínimos no MVP: ativos, localizações, classes, medidores, leituras, solicitações, ordens de serviço, planos, itens de estoque, movimentações (`/v1/movimentacoes-estoque`), documentos com vencimento, pessoas (`/v1/pessoas`), fornecedores (`/v1/fornecedores`), serviços de fornecedor (`/v1/servicos-fornecedor`), links de execução (`/v1/links-execucao`), usuários (leitura), versão (`GET /v1/versao`, exceção aceita à regra do plural); do módulo SRV: `/v1/clientes-servico`, `/v1/especialidades-tecnico`, `/v1/conexoes-empresa` (ações `convidar`, `aceitar`, `recusar`, `suspender`, `reativar`, `encerrar`), `/v1/ordens-servico/{id}/assumir` e `/v1/ordens-servico/{id}/relatorio-pdf`. Na Fase 2, com os módulos: `/v1/contratos-prestacao`, `/v1/avaliacoes-fornecedor`, `/v1/orcamentos-fornecedor`, `/v1/pecas-serializadas`, `/v1/ordens-reparo`, `/v1/locatarios`, `/v1/contratos-locacao`, `/v1/vistorias-locacao`.
- **Operações:** autenticar por chave de API com escopos, criada pelo `administrador`; listar com paginação por cursor e filtro por data de alteração; criar e alterar com `Idempotency-Key`; consultar o estado das cotas de uso.
- **Validações:** as mesmas regras de negócio da interface; erros padronizados com código e mensagem em português.
- **Regras:** a API nunca contorna permissão, estado ou validação. Limites por plano e resposta 429 com `Retry-After` (RNF-014). Detalhes em 3_ESPECIFICACAO_APIS.md.

#### RF-INT-002 · Webhooks

**Fase:** MVP · **Plano mínimo:** Profissional

- **Descrição:** notificação de eventos a sistemas externos.
- **Dados principais:** URL de destino, eventos assinados, segredo, situação. Os eventos usam o nome completo da entidade (`<entidade>.<fato_no_particípio>`, sumário 3.4). Mínimo no MVP: `ordem_servico.criada`, `ordem_servico.status_alterado`, `ordem_servico.encerrada`, `solicitacao.criada`, `leitura_medidor.registrada`, `ativo.status_alterado`, `documento_vencimento.vencendo`, `documento_vencimento.vencido`, `movimentacao_estoque.registrada`; do módulo SRV: `conexao_empresa.ativada`, `conexao_empresa.encerrada`, `vinculo_os.sincronizado`, `vinculo_os.erro`. Na Fase 2: `ordem_reparo.enviada`, `ordem_reparo.retornada`, `contrato_locacao.ativado`, `contrato_locacao.encerrado`, `vistoria_locacao.registrada`, `fornecedor.sla_violado`. A lista completa fica em 3_ESPECIFICACAO_APIS.md.
- **Operações:** cadastrar, testar, pausar; ver painel de entregas com resposta; reenviar manualmente.
- **Regras:** webhooks estão no MVP, no plano Profissional (sumário 2.1 e 4). Assinatura HMAC-SHA256 com carimbo de tempo, retentativa exponencial por 24 horas e fila de mensagens mortas (RNF-015). O conteúdo do evento respeita os mesmos escopos do token que cadastrou o webhook.

#### RF-INT-003 · Importação por planilha

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** carga inicial e atualizações em massa por CSV, e por XLSX (proposta).
- **Dados principais:** localizações, ativos, classes, medidores, leituras, itens de estoque, saldos iniciais, fornecedores e serviços de fornecedor, pessoas, planos, documentos com vencimento, usuários; na conta de uma subcontratada, clientes atendidos, equipamentos de clientes e especialidades de técnico (RF-SRV-002, RF-SRV-004, RF-SRV-005); na Fase 2, peças serializadas e locatários.
- **Operações:** baixar modelo; enviar arquivo; validar sem gravar (simulação) com relatório de erros por linha; confirmar; desfazer a última importação enquanto não houver movimentação dependente (proposta).
- **Validações:** as mesmas do cadastro; chave natural (TAG, código, placa) para atualizar em vez de duplicar.
- **Regras:** a importação é atômica por arquivo: ou grava tudo, ou nada (proposta). Leituras importadas usam origem `importacao`.

#### RF-INT-004 · Modelo canônico de telemática

**Fase:** Fase 2 · **Plano mínimo:** Profissional (proposta)

- **Descrição:** formato único para eventos vindos de qualquer conector: leitura de medidor (odômetro, horas), DTC (SPN, FMI e severidade), combustível, pneu (pressão, temperatura) e posição.
- **Operações:** casar o dispositivo com o ativo por placa, VIN ou identificador do rastreador; fila de dispositivos não casados para correção manual.
- **Validações:** leitura passa por RF-MED-003 com origem `telematica`.
- **Regras:** leituras de alta frequência são reduzidas antes de gravar (proposta: no máximo uma leitura por medidor por hora, mais a última do dia). VIN errado no rastreador vai para a fila de não casados, sem criar ativo.

#### RF-INT-005 · Conectores de telemática

**Fase:** Fase 2 · **Plano mínimo:** Profissional (proposta)

- **Descrição:** conectores prontos, na ordem recomendada pelo estudo: Cobli e Geotab, rFMS (Volvo e Scania), depois Ticket Log, Sascar e Omnilink.
- **Operações:** autorizar o conector com credenciais do cliente; escolher veículos; ver situação da última sincronização.
- **Regras:** a Cobli tem API pública verificada pelo estudo; as APIs de Ticket Log, Sascar e Omnilink não foram verificadas e podem mudar a ordem ou o escopo (ver riscos do roadmap).

#### RF-INT-006 · Regra de DTC para solicitação ou OS

**Fase:** Fase 2 · **Plano mínimo:** Profissional (proposta)

- **Descrição:** conversão configurável de código de falha do veículo em solicitação ou OS.
- **Dados principais:** filtro por severidade ("parar imediatamente", "na próxima parada", "seguro continuar"), SPN, FMI, classe; ação (solicitação, OS `corretiva_emergencial`, só registro).
- **Regras:** o mesmo DTC ativo no mesmo veículo com solicitação ou OS aberta não gera outra (RF-SOL-006). DTC com severidade "parar imediatamente" notifica o `gestor` na hora.

#### RF-INT-007 · Importação de cartão combustível

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** carga de transações de cartão combustível como abastecimentos.
- **Operações:** importar arquivo do fornecedor do cartão com mapeamento de colunas salvo; usar API quando o fornecedor oferecer (não verificado no estudo).
- **Validações:** transação já importada é ignorada; veículo casado por placa ou número do cartão.
- **Regras:** leituras de km importadas usam origem `cartao_combustivel`.

#### RF-INT-008 · Conectores de ERP

**Fase:** Fase 2 · **Plano mínimo:** Enterprise

- **Descrição:** sincronização com Sankhya (via API Gateway com AppKey), TOTVS Protheus (REST ou Mensagem Única) e Omie.
- **Dados principais:** fornecedores, itens, centros de custo, pedidos de compra, notas de entrada, lançamento do custo das OS, número e chave das NF-e de remessa e retorno (RF-EST-019), dados para faturar avarias e excedentes de locação (RF-LOC-009).
- **Regras:** o ERP é a fonte de fornecedores, itens e centros de custo; a Tenaz é a fonte de OS e custos de manutenção (proposta). A Sankhya não reconhece registros gravados direto no banco, por isso toda escrita passa pelo Gateway. Erros de sincronização ficam numa fila visível com reenvio.

#### RF-INT-009 · Importação automática de multas

**Fase:** Fase 3 · **Plano mínimo:** Profissional (proposta)

- **Descrição:** carga automática de autuações por API oficial.
- **Regras:** só para o Brasil; condicionado a **P3** (existência de API Senatran ou Serpro acessível a terceiros). Até lá, multas entram por digitação ou planilha (RF-FRO-010).

#### RF-INT-010 · Ingestão de IoT

**Fase:** Fase 3 · **Plano mínimo:** Enterprise

- **Descrição:** recebimento de medições de sensores para manutenção baseada em condição.
- **Operações:** receber por MQTT ou HTTP; integrar OPC UA e MQTT Sparkplug B por gateway de parceiro; definir variáveis por ativo; ver série temporal.
- **Regras:** a Tenaz fica fora da rede de controle e nunca se conecta direto a controladores; o gateway guarda dados localmente durante falhas e reenvia depois. Volume por empresa limitado por contrato.

#### RF-INT-011 · Servidor MCP

**Fase:** Fase 3 · **Plano mínimo:** Enterprise

- **Descrição:** servidor MCP somente leitura para uso com assistentes de IA de terceiros.
- **Regras:** usa as mesmas permissões e escopos da API; nenhuma escrita; registro de cada consulta na auditoria.

#### RF-INT-012 · Conector de BI

**Fase:** Fase 3 · **Plano mínimo:** Enterprise (proposta)

- **Descrição:** acesso a dados analíticos por Power BI e ferramentas equivalentes.
- **Regras:** dados expostos por visões somente leitura, isoladas por empresa, com atualização ao menos diária (proposta).

---

### 4.16 IA — Assistente

#### RF-IA-001 · Governança de IA e créditos

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** controles exigidos por clientes corporativos para uso de IA.
- **Operações:** ligar e desligar a IA por empresa e por usuário; consultar o log de pedidos e respostas da empresa; ver consumo de créditos; comprar pacotes adicionais.
- **Regras:** nenhuma gravação feita por IA sem confirmação humana; toda resposta cita as fontes; o provedor de modelo tem cláusula de não treinar com dados do cliente; o log fica isolado por empresa (RNF-022). Créditos incluídos no Profissional e pacotes pré-pagos adicionais seguem a recomendação do estudo; valores dependem de **P4**.

#### RF-IA-002 · Busca em manuais e histórico

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** perguntas em linguagem natural sobre documentos técnicos e histórico de OS.
- **Operações:** perguntar na web ou no app (com conexão); abrir a fonte citada (documento e página, ou OS).
- **Regras:** busca vetorial isolada por empresa e restrita ao que o usuário pode ver. Sem fonte encontrada, o assistente diz que não sabe em vez de responder.

#### RF-IA-003 · Rascunho de OS ou solicitação por voz ou foto

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** criação de rascunho a partir de áudio em português ou de foto.
- **Operações:** gravar áudio ou tirar foto; receber rascunho com ativo identificado (pela TAG lida ou pelo QR), descrição, sintoma, tipo e prioridade sugeridos; editar; confirmar.
- **Regras:** o rascunho só vira registro após confirmação do usuário. Offline, o áudio ou a foto ficam guardados e o rascunho é gerado na sincronização.

#### RF-IA-004 · Checklist, procedimento e plano a partir de PDF

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** extração de tarefas e intervalos de um manual em PDF.
- **Operações:** enviar o manual; receber rascunho de modelo de checklist, procedimento ou plano com gatilhos; revisar item a item; aprovar.
- **Regras:** cada item proposto mostra a página de origem. Nada é criado sem aprovação.

#### RF-IA-005 · Assistente de análise de falha

**Fase:** Fase 3 · **Plano mínimo:** Profissional (proposta)

- **Descrição:** apoio à análise de causa raiz a partir do histórico do ativo e da classe.
- **Regras:** o assistente sugere hipóteses com as OS que as sustentam; a conclusão é sempre registrada por uma pessoa.

#### RF-IA-006 · Agentes com aprovação

**Fase:** Fase 3 · **Plano mínimo:** Enterprise

- **Descrição:** agentes que propõem ações: triagem de solicitações, deduplicação, atribuição e sugestão de programação.
- **Operações:** ver fila de propostas; aprovar, editar ou rejeitar cada uma; definir quais tipos de proposta cada agente pode fazer.
- **Regras:** toda ação proposta passa por aprovação humana antes de ser gravada, e a aprovação fica na auditoria com a identificação do agente.

#### RF-IA-007 · Relatórios em linguagem natural

**Fase:** Fase 3 · **Plano mínimo:** Profissional (proposta)

- **Descrição:** perguntas sobre dados ("quais ativos de criticidade A mais pararam no trimestre?") respondidas por ferramentas tipadas.
- **Regras:** o assistente não gera SQL livre; chama consultas pré-definidas com parâmetros validados, que respeitam o isolamento por empresa e as permissões do usuário.

---

### 4.17 MOB — App de campo

#### RF-MOB-001 · Acesso ao app

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** login e perfis no app Android e iOS.
- **Operações:** entrar por e-mail e senha, com o segundo fator TOTP quando ativado, ou por SSO (quando contratado), sempre pelo navegador do sistema (RF-ADM-006); escolher a empresa; desbloquear rapidamente por PIN ou biometria do aparelho (proposta). A entrada por telefone e código fica adiada até a escolha do provedor de SMS (**P11**).
- **Regras:** entra no app quem tem a permissão `mob.usar_app` (RF-ADM-005). O app mostra as telas conforme o perfil: `tecnico` (inclusive o técnico da conta de uma subcontratada, RF-SRV-011), `prestador` (só OS do seu fornecedor), `motorista`, `solicitante`, `almoxarife` (Fase 2) e `gestor` (consulta e aprovações). O `locatario` usa só o portal web (RF-LOC-008), sem acesso ao app (sumário 3.6). Quem tem acesso a mais de uma empresa (por exemplo, à conta da subcontratada e à conta própria de um cliente) escolhe a empresa ao entrar e troca sem sair; o cache de cada empresa fica separado no aparelho. A versão do app aparece na tela de login e na tela principal (RNF-039). O desbloqueio local funciona offline; o token é renovado ao reconectar.

#### RF-MOB-002 · Cache offline e sincronização

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** base local que permite trabalhar sem conexão por pelo menos 7 dias.
- **Dados em cache:** OS atribuídas ao usuário e à sua equipe e OS abertas da unidade (até um limite configurável); para o `prestador`, só as OS atribuídas ao seu fornecedor e os ativos dessas OS (RNF-037); para o técnico de uma subcontratada (SRV), as OS atribuídas a ele, os equipamentos de clientes dessas OS, o nome, o endereço de atendimento e o contato do cliente atendido, o catálogo de serviços e os itens de estoque da própria conta (RF-SRV-011); modelos de vistoria e contratos de locação da unidade (Fase 2); ativos e localizações da unidade com hierarquia; medidores e últimas leituras; modelos de checklist, procedimentos e planos aplicáveis; itens de estoque e último saldo conhecido; documentos técnicos marcados como offline; capacitações dos técnicos (Fase 2); veículo e documentos do motorista.
- **Operações:** primeira carga com conexão; sincronização incremental automática quando houver rede; sincronização manual; indicador de itens pendentes de envio e da última sincronização.
- **Regras:** cada ação offline vira um comando com identificador único numa fila local (outbox), reenviado até confirmação; o servidor aplica cada comando uma única vez. Conflitos seguem RNF-011. O cache sobrevive ao fechamento e à reinicialização do aparelho.

#### RF-MOB-003 · Operações offline

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** o que pode ser feito sem conexão.
- **Operações:** criar solicitação; criar OS; assumir OS da equipe (confirmada na sincronização); executar OS inteira (iniciar, pausar, retomar, tarefas, apontamento, peças, fotos, falha, campos NR-12, assinatura, conclusão); na conta de uma subcontratada, registrar componentes substituídos e coletar a assinatura do cliente atendido (RF-SRV-008, RF-SRV-012); executar checklist; registrar leitura; consultar ativo e histórico em cache; engatar e desengatar implemento; registrar abastecimento e sinistro (Fase 2); realizar vistoria de locação (Fase 2, RF-LOC-005); instalar e remover peça serializada na OS (Fase 2).
- **Regras:** o que exige servidor (IA, consulta fora do cache, aprovação de compra) aparece desabilitado com o motivo. Registros criados offline recebem número provisório até a sincronização.

#### RF-MOB-004 · Leitura de QR e código de barras

**Fase:** MVP · **Plano mínimo:** Essencial

- **Operações:** ler QR de ativo e localização (abre o ativo, ou o checklist ou OS pendente dele); ler código de barras de item de estoque; funcionar offline contra o cache.
- **Regras:** código desconhecido no cache oferece buscar no servidor quando houver conexão.

#### RF-MOB-005 · Fotos, vídeo e áudio

**Fase:** MVP · **Plano mínimo:** Essencial

- **Operações:** capturar foto, vídeo curto (até 30 segundos, proposta) e áudio; anotar sobre a foto (proposta); anexar a OS, tarefa, checklist, solicitação, leitura.
- **Regras:** fotos são comprimidas no aparelho; carregam data, hora e geolocalização da captura; o envio ocorre em segundo plano, inclusive em rede móvel, com opção de enviar só em Wi-Fi.

#### RF-MOB-006 · Notificações push

**Fase:** MVP · **Plano mínimo:** Essencial

- **Operações:** receber push de OS atribuída, OS reatribuída a outro técnico, OS emergencial, mudança em solicitação própria, vencimento de documento e checklist pendente; abrir o registro pela notificação.
- **Regras:** push não carrega dado sensível no texto; o conteúdo é carregado após abrir o app.

#### RF-MOB-007 · Modo motorista

**Fase:** MVP · **Plano mínimo:** Essencial

- **Descrição:** fluxo simplificado para o `motorista` gratuito.
- **Operações:** assumir veículo pelo QR; fazer checklist de pré e pós-viagem; informar km; abrir solicitação com foto ou áudio; consultar documentos do veículo (CRLV-e e outros em cache); registrar abastecimento e sinistro (Fase 2).
- **Regras:** o motorista não vê custos, OS de outros veículos nem dados de outros motoristas.

#### RF-MOB-008 · Modo almoxarife

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Operações:** atender requisições; separar com leitura de código de barras; registrar entrega com confirmação do técnico; contar inventário; consultar saldo.
- **Regras:** contagem offline é reconciliada na aprovação do inventário (RF-EST-005).

#### RF-MOB-009 · Pneus no app

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Operações:** ver o mapa de eixos do veículo; registrar inspeção de sulco e pressão posição por posição; montar, desmontar e rodiziar arrastando pneus entre posições.
- **Regras:** valida as mesmas regras de RF-FRO-003 e RF-FRO-004, também offline.

#### RF-MOB-010 · Gestão do aparelho

**Fase:** MVP · **Plano mínimo:** Essencial

- **Operações:** apagar o cache do aparelho remotamente na próxima conexão (usuário inativado ou aparelho perdido); exigir versão mínima do app; enviar diagnóstico de sincronização ao suporte com consentimento do usuário.
- **Regras:** a limpeza ocorre na próxima conexão, depois do envio dos comandos pendentes. O servidor aceita comandos feitos antes da inativação do usuário e recusa, com registro, os feitos depois (proposta). Os dados locais do app ficam criptografados no aparelho.

### 4.18 PRJ — Projetos e implantação

Módulo da Fase 2, plano Profissional. Controla a implantação de ativos que exige uma sequência de etapas com dependências, como a instalação de um radar (autorização, energia, fundação e poste, laços, montagem, aferição, entrega), a instalação de uma máquina ou a montagem de uma câmara fria. Reaproveita OS, fornecedor, documentos com vencimento, campos extras e o app; o projeto não executa trabalho, organiza OS.

#### RF-PRJ-001 · Modelos de projeto

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** roteiro reutilizável de implantação (`modelo_projeto`), com as etapas (`modelo_etapa_projeto`) e o que cada uma exige.
- **Dados principais:** modelo: nome, descrição, classes de ativo a que se aplica, situação, versão. Etapa do modelo: nome, ordem, etapas predecessoras, duração prevista em dias, peso no avanço físico, executor padrão (`interno` ou `fornecedor`), gera OS (sim/não) e, quando gera, procedimento ou tarefas da OS `implantacao`, tipos de documento exigidos para concluir, campos extras exigidos para concluir, valor de referência da etapa (para pagamento do prestador), marcação "etapa de entrega" (a última, que coloca os ativos em operação).
- **Operações:** criar, copiar, editar (gera versão nova), inativar; criar a partir de um pacote de segmento (RF-ADM-017); ver as dependências em diagrama. Rota `/v1/modelos-projeto`.
- **Validações:** pelo menos uma etapa; predecessoras do mesmo modelo, sem ciclos; peso maior ou igual a zero e soma dos pesos maior que zero; no máximo 40 etapas por modelo (proposta); no máximo uma etapa de entrega.
- **Regras:** projeto já criado guarda a versão do modelo usada; mudar o modelo não altera projetos existentes.

#### RF-PRJ-002 · Projeto de implantação

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** implantação de um ou mais ativos (`projeto`), criada a partir de um modelo.
- **Dados principais:** número `PJ-<AAAA>-<NNNNNN>`, modelo e versão, nome, unidade, localização, ativos implantados (normalmente `em_implantacao`, RF-ATV-013), cliente atendido e contrato com cliente (opcionais, RF-CTR-001), responsável, datas previstas e reais de início e fim, orçamento previsto por etapa, situação (`rascunho`, `em_andamento`, `suspenso`, `concluido`, `cancelado`, sumário 3.3), anexos. Rota `/v1/projetos`, com ações `iniciar`, `suspender`, `retomar`, `concluir` e `cancelar`.
- **Operações:** criar a partir de um modelo, com as etapas copiadas e as datas previstas calculadas pela duração e pelas dependências; criar vários projetos de uma vez (um por ponto ou por ativo) a partir de uma planilha (RF-INT-003) (proposta); incluir e retirar ativos em `rascunho`; iniciar; suspender e retomar com motivo; cancelar com motivo.
- **Validações:** `iniciar` exige pelo menos um ativo ou uma localização e responsável; `suspender` e `cancelar` exigem motivo; cancelar com OS em `em_execucao` exige pausá-las antes; ativo em no máximo um projeto não concluído (proposta).
- **Regras:** iniciar o projeto torna `disponivel` as etapas sem predecessoras. Suspender pausa os prazos das etapas, não as OS (o planejador decide sobre cada OS). Projeto `concluido` ou `cancelado` não é editado. Cancelar não apaga os ativos `em_implantacao`: o sistema pede o destino de cada um (manter em implantação ou inativar com motivo). Gera os eventos `projeto.iniciado`, `projeto.concluido` e `projeto.cancelado`.

#### RF-PRJ-003 · Etapas e dependências

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** execução de cada etapa (`etapa_projeto`) do projeto, na ordem das dependências.
- **Dados principais:** projeto, etapa do modelo, executor (`interno` ou `fornecedor`, com `fornecedor_id`), responsável, datas previstas e reais, situação (`pendente`, `disponivel`, `em_andamento`, `concluida`, `dispensada`, sumário 3.3), OS da etapa, documentos e campos exigidos com a situação de cada um, valor da etapa, motivo de dispensa.
- **Operações:** atribuir executor e responsável; gerar a OS da etapa (RF-PRJ-004); registrar documentos exigidos (RF-CON-005) e campos exigidos; concluir; dispensar com motivo.
- **Validações:** a etapa só passa a `disponivel` quando todas as predecessoras estão `concluida` ou `dispensada`; `em_andamento` quando a primeira OS da etapa inicia ou por ação manual em etapa sem OS; `concluida` exige as OS da etapa `encerrada` ou `cancelada`, pelo menos uma `encerrada` quando a etapa gera OS, e os documentos e campos exigidos registrados; `dispensada` só a partir de `pendente` ou `disponivel`, com motivo e a permissão `prj.dispensar_etapa`, padrão em `gestor`. Criar, iniciar, suspender, retomar, concluir e cancelar projeto e modelo exigem `prj.gerenciar` (sumário 3.6).
- **Regras:** concluir ou dispensar uma etapa torna `disponivel` as sucessoras cujas predecessoras terminaram e notifica seus responsáveis. Etapa com data prevista vencida e não concluída fica "atrasada" (condição calculada, não estado) e pode disparar regra `etapa_projeto_atrasada` (RF-ADM-016). Documento exigido registrado no ativo vale para a etapa quando emitido depois do início da etapa (proposta).

#### RF-PRJ-004 · OS de implantação da etapa

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** cada etapa que "gera OS" produz uma ou mais OS do tipo `implantacao`, executadas como qualquer OS.
- **Dados principais:** OS com `tipo = implantacao`, `origem = projeto`, etapa de origem, ativo ou localização do projeto, tarefas do procedimento da etapa, executor da etapa.
- **Operações:** gerar a OS quando a etapa fica `disponivel` (automático, proposta) ou manualmente; gerar uma OS por ativo do projeto quando a etapa se repete por ativo (proposta); atribuir a fornecedor (RF-TER-003), inclusive com conta conectada (RF-SRV-016), ou a técnico interno; executar no app offline.
- **Validações:** as da OS (RF-OS-001 a RF-OS-014) e os campos extras obrigatórios da OS (RF-OS-021).
- **Regras:** a OS de implantação fica fora dos KPIs de manutenção (RF-OS-002) e entra no custo realizado do projeto. Na conta da subcontratada conectada, a OS chega como qualquer OS do cliente; o projeto e suas etapas não são compartilhados no lançamento, só a OS.

#### RF-PRJ-005 · Acompanhamento do projeto

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** visão de prazo, avanço e custo de um projeto e da carteira de projetos.
- **Dados principais:** cronograma (barras por etapa, previsto e real, com dependências), avanço físico (soma dos pesos das etapas `concluida` e `dispensada` sobre a soma dos pesos), etapas atrasadas, custo previsto (orçamento) e realizado (custos das OS de implantação, RF-OS-013, mais lançamentos avulsos do projeto) por etapa e total, desvio de prazo em dias.
- **Operações:** ver o projeto em cronograma, lista de etapas e quadro; ver a carteira por unidade, cliente, fornecedor e situação; filtrar atrasados; exportar em PDF e XLSX e pela API.
- **Regras:** os números mostram como são calculados (RF-IND-002). A etapa dispensada conta no avanço físico e não conta no custo.

#### RF-PRJ-006 · Medição de etapas para pagamento do prestador

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** relatório por fornecedor e período com as etapas concluídas por ele e o valor de cada uma, base para o pagamento da empreitada.
- **Dados principais:** fornecedor, período, etapas `concluida` no período com projeto, ativo, data de conclusão e valor da etapa; ajustes com motivo; total.
- **Operações:** gerar por fornecedor e período; ajustar valor de uma linha com motivo e a permissão `prj.ajustar_medicao`; marcar como enviada ao ERP (proposta); exportar em PDF e XLSX e pela API.
- **Validações:** só etapas `concluida` entram; etapa dispensada não entra; uma etapa entra em uma única medição (proposta).
- **Regras:** a Tenaz não paga nem fatura: o pagamento sai do ERP. O valor padrão vem do valor da etapa no projeto, que vem do modelo ou do contrato de prestação (RF-TER-006) quando houver.

#### RF-PRJ-007 · Entrega e entrada em operação

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** fechamento do projeto, com termo de entrega e entrada dos ativos em operação.
- **Dados principais:** etapa de entrega, termo de entrega (gerado em PDF com os ativos, as etapas, os documentos e as fotos), assinatura do responsável e, quando houver, do representante do cliente atendido (RF-CON-014), data de entrada em operação.
- **Operações:** concluir a etapa de entrega; gerar e assinar o termo; concluir o projeto.
- **Validações:** `concluir` o projeto exige todas as etapas `concluida` ou `dispensada`; a etapa de entrega não pode ser dispensada (proposta).
- **Regras:** concluir o projeto executa `colocar_em_operacao` em todos os ativos do projeto que ainda estão `em_implantacao` (RF-ATV-013), ativa os planos e, quando o projeto tem contrato com cliente, inclui os ativos como itens cobertos a partir da data de entrada em operação (RF-CTR-002) (proposta).

### 4.19 CTR — Contratos com clientes

Módulo da Fase 2, plano Profissional. Controla o contrato de serviço da empresa com um cliente atendido que exige disponibilidade dos equipamentos, como o contrato de fiscalização eletrônica com uma prefeitura, de elevadores com um condomínio ou de equipamentos médicos com um hospital. Usa o cadastro `cliente_servico` (sumário 3.1), as OS e os registros de parada, sem mudar o cálculo técnico da disponibilidade (RF-IND-003).

#### RF-CTR-001 · Contrato com cliente

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** contrato de serviço (`contrato_cliente`) com um cliente atendido.
- **Dados principais:** número `CT-<AAAA>-<NNNNNN>`, cliente atendido (`cliente_servico`), número do contrato no cliente (ex.: número do processo licitatório) (proposta), vigência (início, fim, renovação), unidades, meta de disponibilidade do contrato (%), calendário de apuração (`continuo_24x7` ou calendário do contrato, com dias e horários), prazos de atendimento e de solução por prioridade, faixas de desconto (`faixa_desconto_contrato`: de % de disponibilidade, até %, desconto em % do valor do item), causas externas aceitas (subconjunto da lista canônica de `causa_externa`), valor mensal de referência, anexos, situação (`rascunho`, `vigente`, `encerrado`, `cancelado`). Rota `/v1/contratos-cliente`, com ações `ativar`, `encerrar` e `cancelar`.
- **Operações:** criar em rascunho; incluir itens (RF-CTR-002); ativar; aditar com histórico (proposta); encerrar; cancelar; alertas D-60, D-30 e D-7 antes do fim da vigência (proposta).
- **Validações:** `rascunho` → `vigente` exige pelo menos um item, meta, calendário e faixas sem sobreposição; só `rascunho` vai para `cancelado`; meta entre 0 e 100; faixas contíguas, da meta para baixo.
- **Regras:** a ativação gera `contrato_cliente.ativado` e o encerramento gera `contrato_cliente.encerrado`. Os prazos de atendimento e de solução valem para as OS dos itens cobertos; o estouro gera `contrato_cliente.sla_violado` e notifica o `gestor`. Contrato `encerrado` ou `cancelado` não é editado.

#### RF-CTR-002 · Itens cobertos

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** ativos cobertos pelo contrato (`item_contrato_cliente`), com datas de inclusão e exclusão.
- **Dados principais:** contrato, ativo, data de inclusão, data de exclusão, meta própria do item (opcional, sobrepõe a do contrato), valor mensal do item, situação.
- **Operações:** incluir e excluir itens em rascunho ou por aditivo; incluir automaticamente os ativos de um projeto concluído (RF-PRJ-007) (proposta); ver a disponibilidade corrente de cada item.
- **Validações:** ativo `proprio`, `de_terceiro` ou `de_cliente`, não `inativo` na inclusão, salvo `em_implantacao` com inclusão a partir da entrada em operação; um ativo em no máximo um contrato com cliente vigente no mesmo período.
- **Regras:** o item só conta na medição entre a inclusão (ou a entrada em operação, o que for depois) e a exclusão.

#### RF-CTR-003 · Paradas por causa externa

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** marcação, na OS, de que a parada teve causa fora do controle da empresa, para que a medição do contrato possa desconsiderá-la.
- **Dados principais:** `ordem_servico.causa_externa` (`vandalismo`, `acidente_terceiro`, `falta_energia`, `obra_terceiro`, `solicitacao_cliente`, `caso_fortuito`), evidências (boletim de ocorrência, protocolo da concessionária, fotos), quem marcou, aprovação da exclusão na medição (quem, quando, motivo).
- **Operações:** marcar na abertura ou na execução, pela web, pelo app ou pelo prestador; anexar evidência; aprovar ou recusar a exclusão ao fechar a medição.
- **Validações:** causa da lista canônica; evidência obrigatória por parâmetro do contrato (proposta); só causa aceita no contrato pode ser excluída.
- **Regras:** a causa externa **não muda** o MTTR, o MTBF nem a disponibilidade técnica do ativo (RF-IND-003): só a medição do contrato desconta a parada, e só depois de aprovada. A mesma marcação alimenta o relatório de danos causados por terceiros (proposta).

#### RF-CTR-004 · Medição mensal do contrato

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** apuração por período (`medicao_contrato`) da disponibilidade contratual de cada item e do desconto devido.
- **Dados principais:** contrato, período, por item: horas requeridas (pelo calendário de apuração), horas paradas (registros de parada das OS do item, RF-OS-005), paradas excluídas aprovadas, disponibilidade contratual, meta, faixa atingida, desconto em % e em valor; total do contrato; memória de cálculo; situação (`em_apuracao`, `fechada`, sumário 3.3); quem fechou e quando.
- **Operações:** abrir a medição do período (automática no primeiro dia do mês seguinte, proposta); revisar paradas e aprovar exclusões; fechar; reabrir com motivo e a permissão `ctr.reabrir_medicao`; exportar em PDF e XLSX e pela API; marcar como enviada ao cliente ou ao ERP (proposta).
- **Validações:** fechar exige todas as paradas com causa externa decididas (`decisao_exclusao` `aprovada` ou `recusada`) e a permissão `ctr.fechar_medicao`; medição `fechada` não é editada; criar, ativar, encerrar e cancelar contrato exigem `ctr.gerenciar` (sumário 3.6).
- **Regras:** disponibilidade contratual do item = (horas requeridas − horas paradas não excluídas) ÷ horas requeridas, dentro do calendário de apuração; a fórmula e o tratamento de paradas que cruzam o período estão em 4_FLUXOS_PROCESSO.md. A medição mostra como é calculada (RF-IND-002). A Tenaz não fatura: o valor com desconto vai para o ERP. Mudança em OS depois do fechamento não altera a medição fechada; aparece como divergência na reabertura (proposta). Gera `medicao_contrato.fechada`.

#### RF-CTR-005 · Painel do contrato

**Fase:** Fase 2 · **Plano mínimo:** Profissional

- **Descrição:** acompanhamento corrente do contrato antes do fechamento.
- **Dados principais:** disponibilidade parcial do mês por item e do contrato, itens abaixo da meta, desconto projetado, OS abertas dos itens com prazo de atendimento e de solução, paradas com causa externa aguardando decisão.
- **Operações:** ver por contrato e por cliente; ordenar itens pelo impacto no desconto; abrir a OS a partir do item; exportar.
- **Regras:** a projeção usa o mesmo cálculo da medição, com o período até o momento. O acesso do cliente final ao painel depende de **P10**; até lá, o painel é interno e a medição sai em PDF e XLSX.

---
## 5. Requisitos não funcionais

As metas da seção 5 do sumário e as da seção de requisitos não funcionais do estudo são obrigatórias. As demais estão marcadas como "(proposta)". O detalhamento de segurança, backup e recuperação fica em 9_SEGURANCA_CONFIABILIDADE.md.

### 5.1 Disponibilidade e recuperação

| ID | Requisito | Meta | Origem |
|---|---|---|---|
| RNF-001 | Disponibilidade do serviço (web, API e sincronização do app) | Essencial e Profissional: SLA contratual de **99,9%** ao mês, com SLO interno de **99,95%**; Enterprise: SLA de **99,95%**, com SLO interno de **99,97%** | Sumário 2.1 e 5 |
| RNF-002 | Créditos de SLA | Crédito em qualquer mês abaixo do SLA, sem exigir queda mínima de 60 minutos seguidos | Estudo |
| RNF-003 | Perda máxima de dados (RPO) | **≤ 5 minutos**, com recuperação a um ponto no tempo (PITR) a partir de WAL contínuo | Sumário 5; estudo |
| RNF-004 | Tempo de recuperação (RTO) | **≤ 1 hora** em falha de zona; **≤ 4 horas** em falha de região, condicionado à decisão de **P7** | Sumário 5; estudo |
| RNF-005 | Backups e restauração | PITR com retenção de **14 a 35 dias**; cópia diária criptografada em outra região e outra conta; **teste de restauração trimestral** documentado | Estudo |
| RNF-006 | Residência de dados | Dados primários, anexos e backups primários na região São Paulo, **para os três países** (ADR-013, Parte D); para empresa da Argentina ou do Paraguai isso é transferência internacional, com aditivo por país no contrato, a Tenaz como encarregada (**P14**); RPO e RTO publicados no site | Sumário 1.2, 1.4 e 3.6; estudo; ADR-013 |

RNF-004 (falha de região) e RNF-005 (cópia em outra região) exigem uma segunda região, enquanto RNF-006 promete dados no Brasil. O local da cópia de recuperação de desastre é a pendência **P7** do sumário, decisão do dono do produto com apoio jurídico. Este documento não decide; registra as opções:

| Opção | Efeito em RNF-004 e RNF-006 | Custo e condição |
|---|---|---|
| Cópia em região da AWS fora do Brasil | RTO de região de 4 horas viável; dados primários continuam no Brasil, mas há transferência internacional | Cláusulas-padrão da Res. CD/ANPD 19/2024 e declaração no contrato de tratamento de dados |
| Segundo provedor de nuvem no Brasil | Dados sempre no Brasil; RTO de região depende da automação entre provedores e pode passar de 4 horas no início | Maior custo e complexidade operacional |
| Sem cópia fora da região São Paulo | Promessa "dados no Brasil" mantida; não há RTO de região a publicar | Falha de região fica sem recuperação garantida; RNF-004 de região e a cópia de RNF-005 caem |

Até a decisão, vale a regra: dados primários sempre em São Paulo; cópia de recuperação de desastre só onde o contrato com o cliente permitir; o RTO de região só é publicado depois de P7.

### 5.2 Desempenho

| ID | Requisito | Meta | Origem |
|---|---|---|---|
| RNF-007 | Tempo de resposta da API | p95 **< 300 ms** em leitura e **< 500 ms** em escrita, medido no servidor, excluídas exportações e relatórios assíncronos | Sumário 5; estudo (meta de projeto) |
| RNF-008 | Desempenho da web | Painéis com p95 **< 2 s** (estudo); listas de até 50 linhas com p95 < 1 s e busca global com p95 < 1 s (proposta); primeira carga da aplicação com LCP < 2,5 s em banda larga (proposta) | Estudo; proposta |
| RNF-009 | Processamentos em lote | Avaliação de gatilhos de tempo de uma empresa com 2.000 ativos em < 5 min; importação de 10.000 linhas validada em < 2 min; rotina diária de alertas de vencimento concluída até 07h00 no fuso da unidade | Proposta |

### 5.3 App offline e sincronização

| ID | Requisito | Meta | Origem |
|---|---|---|---|
| RNF-010 | Autonomia offline | Operar **≥ 7 dias** sem conexão, sem perda de registro; convergência **< 60 s** após reconexão para até 500 comandos pendentes (o volume é proposta) | Sumário 5; estudo |
| RNF-011 | Conflitos e integridade | Última escrita vence **por campo** em dados cadastrais; registros de evento (apontamentos, fotos, leituras, comentários) são só acrescentados, nunca sobrescritos; regras de domínio prevalecem em estados (OS cancelada no servidor recusa "concluir" vindo do app, com aviso ao técnico); comandos idempotentes pelo identificador gerado no aparelho | Estudo |
| RNF-012 | Desempenho e recursos do app | Primeira carga de uma unidade com 2.000 ativos e 500 OS em < 3 min em 4G; abertura do app em < 3 s em aparelho intermediário; cache sem anexos ≤ 300 MB; consumo de dados em segundo plano ≤ 20 MB por dia sem fotos | Proposta |

### 5.4 Segurança e isolamento

| ID | Requisito | Meta | Origem |
|---|---|---|---|
| RNF-013 | Isolamento entre empresas | **Zero** vazamento em testes automatizados a cada build; RLS com `FORCE ROW LEVEL SECURITY` em toda tabela com dado de empresa; filtro explícito por empresa no repositório; papel de aplicação sem `BYPASSRLS`; contexto da empresa definido por transação (`set_config(..., true)`), nunca por sessão; índice por empresa em toda tabela com RLS. `conexao_empresa`, que liga duas empresas, fica no esquema da plataforma, fora do RLS por empresa, e é protegida por RNF-041. O teste de vazamento usa empresas dos **três países** (Brasil, Argentina e Paraguai) e cobre os seis pares entre elas (RNF-044) | Sumário 5 e R2; estudo |
| RNF-014 | Limites de uso da API e do portal | Profissional: **≥ 120 req/min** por empresa (estudo); Enterprise: 600 req/min (proposta); portal QR: 10 solicitações por hora por origem (proposta); resposta 429 com `Retry-After` | Estudo; proposta |
| RNF-015 | Entrega de webhooks | Assinatura HMAC-SHA256 com carimbo de tempo; retentativa exponencial por **24 h**; fila de mensagens mortas; painel de entregas; primeira tentativa em p95 < 60 s após o evento (proposta) | Estudo; proposta |
| RNF-016 | Criptografia | TLS 1.2 ou superior em trânsito (preferir 1.3); AES-256 em repouso com chaves gerenciadas por KMS; dados locais do app criptografados | Estudo |
| RNF-017 | Autenticação e sessão | MFA obrigatório para a equipe interna da Tenaz e opcional (exigível pelo cliente) para usuários; sessão web expira após 12 h de inatividade e token do app é renovável por até 30 dias offline (proposta); senhas com hash resistente (Argon2id ou bcrypt) (proposta) | Estudo; proposta |
| RNF-018 | Segurança no desenvolvimento | SAST, DAST e SCA no CI; **pentest externo anual** com sumário disponível ao cliente; programa de divulgação responsável; correção de vulnerabilidade crítica em até 7 dias e alta em até 30 dias (proposta) | Estudo; proposta |
| RNF-019 | Trilha de auditoria | Imutável (só acréscimo), com quem, o quê, quando, antes e depois, IP, dispositivo, tela ou módulo e tipo de operação em OS, ativos, estoque e permissões; integridade verificável por hash de cada evento e selo encadeado guardado fora do banco (RF-ADM-011); exportável; retenção mínima de **5 anos**, configurável por empresa | Estudo (o prazo é proposta do estudo) |
| RNF-020 | Observabilidade | OpenTelemetry com identificador da empresa (`empresa_id`) em traces, métricas e logs; alertas por consumo do orçamento de erro do SLO; **página de status pública com histórico de 90 dias** | Sumário 1.2; estudo |
| RNF-037 | Isolamento do prestador | Usuário `prestador` lê e grava só OS com `fornecedor_id` igual ao seu `usuario.fornecedor_id`, os ativos dessas OS (campos permitidos) e o almoxarifado `em_terceiro` do seu fornecedor; filtro aplicado no servidor **além** do RLS por empresa, em web, API, sincronização do app, exportação, busca e webhooks; teste automatizado a cada build com dois fornecedores na mesma empresa e o mesmo fornecedor em duas empresas, com **zero** vazamento; OS reatribuída sai do cache do prestador na sincronização seguinte. Link de uso único: token aleatório de pelo menos 128 bits guardado só como hash, levado no fragmento da URL (`/l/#<token>`) e enviado só no corpo de `POST /v1/publico/links-execucao/ativar`, nunca no caminho nem na query; sessão curta com escopo `link`; uso único, expiração (padrão 7 dias), revogável, limitado a uma OS, sem custos nem histórico, com limite de tentativas por origem | Sumário 2.1 e R4; proposta (valores) |
| RNF-038 | Isolamento do locatário | Usuário `locatario` lê só ativos com item entregue em contrato `cedido` vigente do seu `locatario`, o resumo das OS e solicitações desses ativos e as vistorias e relatórios do seu contrato; cria só solicitações sobre esses ativos; nunca vê custos internos, outros locatários ou histórico anterior ao contrato; filtro no servidor além do RLS por empresa; teste automatizado a cada build com dois locatários do mesmo ativo em períodos diferentes, com **zero** vazamento | Sumário 3.1; proposta (casos de teste) |

### 5.5 Privacidade, IA e portabilidade

| ID | Requisito | Meta | Origem |
|---|---|---|---|
| RNF-021 | LGPD | Kit no lançamento: DPA, RIPD modelo e encarregado nomeado; dados de saúde da `pessoa` (exame toxicológico, ASO) removidos do escopo em 2026-11-13 por decisão do dono do produto: a Tenaz não os guarda e não tem a permissão `pessoa.saude.ler`; dados de vítimas de sinistro e identificação de locatário pessoa física com acesso restrito; identificação pessoal completa de qualquer tipo (CPF, CUIL, DNI, CI) mascarada na interface e com visualização gravada na trilha; para empresa da Argentina (Ley 25.326) e do Paraguai (Ley 7.593/2025, em vigor por volta de novembro de 2027), a Tenaz como encarregada e as garantias da transferência internacional dependem da **P14**; fluxo de incidente compatível com o prazo de comunicação à ANPD de 3 dias úteis (Res. 15/2024, a validar juridicamente) | Estudo |
| RNF-022 | Governança de IA | Provedor com cláusula de não treinar com dados do cliente; log de pedidos por empresa; desligamento por empresa e por usuário; citação de fontes; **nenhuma gravação sem confirmação humana**; processamento preferencialmente no Brasil ou com cláusulas-padrão da Res. CD/ANPD 19/2024 quando houver transferência internacional | Estudo |
| RNF-023 | Portabilidade | Exportação completa da empresa (JSON/CSV + anexos), agendável, entregue em **≤ 24 h** | Estudo |
| RNF-024 | Retenção e exclusão | Após o encerramento do contrato, 90 dias para exportação; exclusão definitiva em até 30 dias depois disso, incluindo réplicas, e expiração natural nos backups dentro da retenção de RNF-005; comprovante de exclusão | Proposta |

### 5.6 Usabilidade, acessibilidade e compatibilidade

| ID | Requisito | Meta | Origem |
|---|---|---|---|
| RNF-025 | Acessibilidade | Web em conformidade com **WCAG 2.1 nível AA**; app com rótulos para leitor de tela (TalkBack e VoiceOver), contraste AA, texto ampliável até 200% sem perda de função e alvos de toque ≥ 48 dp, pensados para uso com luvas (proposta) | Solicitação do produto; proposta |
| RNF-026 | Navegadores | Duas últimas versões estáveis de Chrome, Edge, Firefox e Safari (16.4 ou superior); telas de gestão otimizadas para ≥ 1280 px; portal de solicitação e consultas utilizáveis a partir de 360 px de largura | Proposta |
| RNF-027 | Aparelhos móveis | Android 8.0 (API 26) ou superior e iOS 15.1 ou superior, respeitando o mínimo do SDK do Expo adotado; aparelhos Android de entrada com 3 GB de RAM; coletores robustos Android com leitor de código de barras integrado | Proposta |
| RNF-028 | Usabilidade sem treinamento | Novo `administrador` cria o primeiro ativo e a primeira OS em ≤ 15 min sem ajuda; técnico conclui uma OS simples no app em ≤ 6 toques após abri-la; motorista conclui checklist de 20 itens em ≤ 5 min (todas propostas, a medir em teste com usuários) | Sumário 1.2; proposta |
| RNF-029 | Idioma e formatos | Interface nos idiomas de RNF-043, com todo texto externalizado em catálogo; nos três países, data `dd/mm/aaaa`, vírgula decimal e ponto de milhar; moeda pela empresa (RNF-045: R$ no Brasil, `$` no peso argentino, `₲` no guarani, sem casas, e o código ISO quando a moeda não é a do país); fuso por nome IANA, da lista do país, calculado no fuso da unidade. A formatação usa o idioma do usuário com o país da empresa | Sumário 1.4 e 3.6; ADR-013 |

### 5.7 Escalabilidade

| ID | Requisito | Meta | Origem |
|---|---|---|---|
| RNF-030 | Escala por empresa | Atender, com as metas de desempenho acima, uma empresa com **2.000 ativos**, 50 técnicos, 500 veículos, 1.000 usuários gratuitos, 50.000 OS por ano, 2 milhões de leituras de medidor por ano e 500 GB de anexos acumulados; testes de carga com 5 vezes esse volume (10.000 ativos) sem violar RNF-007 | Porte do sumário 1.3; números são proposta |
| RNF-031 | Escala da plataforma | Arquitetura de banco compartilhado com RLS dimensionada para **10.000 empresas** e 500 mil usuários; 1.500 req/s de pico na API com escala horizontal da aplicação; caminho documentado para banco dedicado por empresa (RF-ADM-014) e para particionamento de tabelas de alto volume (leituras, auditoria, telemetria) | Proposta |

### 5.8 Contrato técnico, anexos e evidências

| ID | Requisito | Meta | Origem |
|---|---|---|---|
| RNF-032 | Versionamento da API | Contrato em OpenAPI 3.1; paginação por cursor; `Idempotency-Key` em criações; mudança incompatível só em nova versão (`/v2`), com a anterior mantida por no mínimo 12 meses após o aviso (proposta) | Estudo; proposta |
| RNF-033 | Anexos | Até 25 MB por arquivo e 100 MB por manual em PDF (proposta); armazenamento de objetos na região São Paulo, criptografado; URLs de acesso temporárias; varredura de malware no envio; toda tela de envio de anexo, na web e no app, mostra o aviso fixo de RF-ADM-015 (a Tenaz não guarda dado de saúde), sem opção de fechar (sumário 3.6); o nome do arquivo não é verificado; o aviso entra na etapa E02, com o componente de anexos. Aceite: o componente de envio de anexo exibe o aviso em toda tela que o usa, verificado por teste do componente e na verificação de acessibilidade (RNF-025) | Proposta; decisão do dono do produto (aviso) |
| RNF-034 | Integridade de registros assinados | Hash SHA-256 do conteúdo assinado guardado com a assinatura; carimbo de tempo do servidor; verificação de integridade disponível ao cliente; forma de assinatura dos documentos de SST conforme **P1** | Estudo; P1 |
| RNF-035 | Notificações | Push e notificação na central em p95 < 1 min após o evento (inclusive OS recebida de cliente conectado, contado a partir da gravação na conta da subcontratada); e-mails transacionais em p95 < 5 min; alertas de vencimento sem duplicidade (proposta) | Proposta |

### 5.9 Certificações

| ID | Requisito | Meta | Origem |
|---|---|---|---|
| RNF-036 | Certificações e selos | LGPD desde o lançamento; SOC 2: controles e janela de observação começam na **Fase 2**, relatório **Type II na Fase 3**; ISO 27001 em seguida | Sumário 4; estudo |

### 5.10 Controle de versão

| ID | Requisito | Meta | Origem |
|---|---|---|---|
| RNF-039 | Controle de versão e versão visível na interface | Todos os planos, desde o MVP. (1) **Versionamento semântico** (SemVer 2.0.0) no formato `MAJOR.MINOR.PATCH` (ex.: `1.4.2`); pré-lançamento só fora da produção (ex.: `1.5.0-rc.1`); metadados de build opcionais (ex.: `+a1b2c3d`, hash curto do commit). (2) A versão atual aparece **sempre** no canto inferior direito da tela principal da web e do app, e também na tela de login, no formato `v1.4.2`; fora da produção, acompanha o nome do ambiente (ex.: `v1.5.0-rc.1 · Homologação`). (3) Clicar ou tocar na versão abre "Sobre", com versão da web ou do app, versão da API, data do build, commit curto e ambiente. (4) A versão **nunca é digitada à mão**: vem do build, a partir da tag git `vX.Y.Z` do release; web, app, API e workers do mesmo release têm a mesma versão. (5) Toda entrega em produção tem tag `vX.Y.Z`, CHANGELOG atualizado no formato Keep a Changelog e incremento conforme SemVer: MAJOR para quebra de compatibilidade da API pública, MINOR para funcionalidade nova, PATCH para correção; o CI **bloqueia o release** se a tag, o CHANGELOG e a versão embutida no build divergirem. (6) Endpoint público `GET /v1/versao`, sem autenticação, devolve versão, commit e data do build. (7) A versão não cobre conteúdo nem botões, respeita a área segura do celular e segue o padrão visual de 7_PADROES_DESIGN.md | Decisão do dono do produto |
| RNF-040 | Histórico de versões sempre atualizado | Todos os planos, desde o MVP. (1) O arquivo **`CHANGELOG.md`**, na raiz do repositório, é o registro oficial de **todas as implementações e correções de cada versão**. (2) Formato Keep a Changelog 1.1.0 em português, com uma seção por versão (`## [1.4.2] - 2026-11-12`) e os grupos **Adicionado**, **Alterado**, **Corrigido**, **Removido**, **Obsoleto** e **Segurança**. (3) Cada item diz o que mudou em linguagem que o usuário entende e cita a referência: ID do requisito (`RF-...`, `RNF-...`) para implementação, número do chamado ou da issue para correção. Migrações de banco e mudanças na API pública são registradas; quebra de compatibilidade aparece destacada com **QUEBRA**. (4) Toda mudança de comportamento entra na seção `## [Não lançado]` **no mesmo pull request** que a implementa; o CI rejeita o pull request que altera código de produto sem alterar o `CHANGELOG.md`, exceto quando marcado como interno (refatoração, testes, ferramentas), com justificativa. (5) No release, a seção `[Não lançado]` vira a versão da tag `vX.Y.Z` (RNF-039); versão, data e tag têm de coincidir. (6) Entradas de versões publicadas não são reescritas; um erro é corrigido com nota na própria entrada. (7) As notas da versão em uso aparecem para o usuário em "Novidades desta versão", aberto pelo diálogo "Sobre" (RNF-039), e o histórico completo fica disponível no sistema e no portal de documentação | Decisão do dono do produto |

### 5.11 Conexão entre contas (SRV)

Respondem ao risco **R5** do sumário: a sincronização entre contas conectadas copia dado além do permitido ou para a conta errada. Valem desde o MVP. A decisão de arquitetura é o ADR-009 (conexão entre contas), em 5_ARQUITETURA_SOFTWARE.md.

| ID | Requisito | Meta | Origem |
|---|---|---|---|
| RNF-041 | Isolamento entre contas conectadas | (1) **Sem leitura cruzada:** nenhuma consulta, da aplicação ou do processo de sincronização, lê dados de duas empresas na mesma transação; o processo lê da conta de origem e grava na conta de destino em **transações separadas**, cada uma com o contexto da própria empresa (RNF-013). (2) O processo de sincronização usa um **papel de banco próprio**, sem `BYPASSRLS`, que só grava nas tabelas e colunas que a sincronização pode alterar; o papel da aplicação não lê `conexao_empresa` diretamente, só por um serviço da plataforma que devolve as conexões da empresa do contexto. (3) **Lista fechada de campos:** só os campos ligados em `campos_compartilhados` saem da conta de origem; a lista é validada no envio e de novo no recebimento, e campo desconhecido é recusado. (4) O destino é sempre a outra parte da `conexao_empresa` do registro, conferida no envio e no recebimento; conexão `suspensa` ou `encerrada` não recebe gravação. (5) **Auditoria nos dois lados:** cada gravação gera evento em cada conta, com a conexão, o registro de origem e o de destino (RF-ADM-011). (6) **Teste automatizado a cada build com três contas** (um cliente conectado a duas subcontratadas, uma delas também conectada a um segundo cliente): **zero** dado fora da lista de campos da conexão e **zero** gravação na conta errada; o teste cobre web, API, sincronização do app, exportação, busca, webhooks e o próprio processo de sincronização. (7) Pentest do lançamento inclui a conexão entre contas | Sumário 3.5 e R5; proposta (desenho do teste) |
| RNF-042 | Integridade e atraso da sincronização entre contas | (1) **Atraso:** com as duas contas online, uma mudança confirmada no servidor da conta de origem aparece na conta de destino em **p95 < 60 s** (medido do commit na origem ao commit no destino); comando que chega de um app offline conta a partir da chegada ao servidor. (2) **Outbox nos dois lados:** a mudança e o evento a sincronizar são gravados na mesma transação da conta de origem; nada se perde com queda do processo, do banco ou da rede. (3) **Idempotência:** cada evento tem identificador único e o destino aplica cada um uma única vez; reprocessar não duplica OS, apontamento, componente, foto ou leitura. (4) **Ordem:** eventos do mesmo registro são aplicados na ordem em que ocorreram na origem. (5) **Dono por campo:** cada campo tem um único lado dono (execução na subcontratada; planejamento, custo e liberação no cliente, RF-SRV-018); não há "última escrita vence" entre contas. (6) **Erro visível:** falha em 3 tentativas seguidas (proposta) deixa `vinculo_os.estado_sincronizacao = erro`, gera `vinculo_os.erro` e aparece nas duas contas; as retentativas automáticas continuam, exponenciais, por 24 h (proposta), e depois resta o reprocessamento manual; o sucesso volta o vínculo a `sincronizado` e gera `vinculo_os.sincronizado`. (7) **Reconciliação:** rotina diária compara origem e destino de cada vínculo ativo e abre erro para toda divergência, com meta de **zero** divergência não sinalizada (proposta). (8) Encerrar a conexão para a sincronização em até 60 s, sem apagar dados de nenhum lado | Sumário 3.5 e R5; proposta (retentativa e reconciliação) |

### 5.12 Multipaís, idiomas e dinheiro

Decisão do dono do produto de 13/11/2026 (sumário 0.8, seção 3.6, linhas "Países e idiomas" e "País, identificação, moeda e idioma"; [ADR-013](5_ARQUITETURA_SOFTWARE.md) e [etapa 014](etapas/014_2026-11-13_multipais-e-idiomas.md)): o MVP é **multipaís** (Brasil, Argentina e Paraguai) e **bilíngue** (português do Brasil e espanhol latino-americano neutro). Valem desde o MVP, a partir da etapa E01, nos três planos. Nenhum deles muda o isolamento (RNF-013): as colunas novas ficam em tabelas que já têm `empresa_id` e RLS. Normas, fiscal, proteção de dados, preço e prazo dependem das pendências **P12** a **P15** (seção 7).

| ID | Requisito | Meta | Origem |
|---|---|---|---|
| RNF-043 | Idiomas suportados e escolha do idioma | (1) **Idiomas** num catálogo fechado de etiquetas BCP 47: `pt-BR` (idioma de origem e de último recurso) e `es-419` (espanhol latino-americano neutro, sem voseo); idioma novo entra por requisito, código e teste, nunca por configuração. (2) **Escolha:** `usuario.idioma` é obrigatório; `empresa.idioma_padrao` é obrigatório, sugerido pelo país (`pt-BR` no Brasil, `es-419` na Argentina e no Paraguai) e alterável pelo `administrador`; no primeiro acesso a tela já vem no idioma do navegador ou do aparelho, se suportado; depois, só o próprio usuário muda o seu. (3) **Onde vale:** telas da web e do app, e-mails e notificações, documentos gerados (OS impressa, relatório em PDF, exportação, extração da trilha) e o texto (`title`, `detail`, `erros[].mensagem`) dos erros da API; o `codigo` do erro e os recursos JSON **não** mudam com o idioma e não levam texto traduzido. (4) **Negociação na API**, nesta ordem: `Accept-Language`, `usuario.idioma` da sessão, `empresa.idioma_padrao`, `pt-BR`; toda resposta com texto para pessoa traz `Content-Language` e `Vary: Accept-Language`; web e app sempre enviam `Accept-Language`. (5) **Destinatário:** e-mail disparado pela própria pessoa sai no idioma da requisição; e-mail ou notificação disparado por outra pessoa ou por rotina, no `usuario.idioma` do destinatário; sem usuário, no `idioma_padrao` da empresa; documento gerado, no idioma de quem pede, resolvido no enfileiramento. (6) **O que não se traduz:** o que a empresa ou o usuário digita nunca é traduzido, nem por máquina; o que o sistema cria a partir de catálogo (perfis padrão, permissões, telas da trilha, situações) aparece pela chave do código, no idioma de quem lê; o número legível do documento (`OS-` ou `OT-`) nunca muda com o idioma; o idioma não entra no conteúdo com hash da trilha. (7) **Qualidade verificada por máquina:** o build quebra se faltar ou sobrar chave em qualquer idioma, se as variáveis `{{...}}` ou as formas de plural diferirem, se um código de erro ou modelo de e-mail não tiver texto nos dois idiomas, se um texto usar termo proibido do glossário bilíngue, ou se texto para pessoa for escrito direto no código Go. (8) **Revisão humana:** toda mudança em `es-419` passa pela revisão do glossário bilíngue; os termos marcados no glossário seguem para **revisão por falante nativo antes da produção**. (9) O aviso de dado de saúde (RF-ADM-015, RNF-033) aparece só com a interface em português; em espanhol, fica como customização se o cliente pedir (sumário 3.6) | Decisão do dono do produto (etapa 014); ADR-013, Parte B |
| RNF-044 | Países suportados e perfil de país | (1) **Países** num catálogo fechado: `BR`, `AR` e `PY`; país novo entra por requisito, código e teste (perfil de país no código, não configuração), e nenhum módulo testa o país por conta própria: pergunta ao perfil. (2) **Um país por empresa** no MVP: `empresa.pais` é escolhido na criação da conta e não é alterável pela API nem pelo papel da aplicação (trava no banco); a unidade não tem país próprio; grupo com operação em dois países abre uma conta por país. (3) **Identificação** pelo par `identificacao_tipo` e `identificacao`, normalizado (maiúsculas, só letras e dígitos; a máscara é só de exibição): empresa e unidade aceitam só `cnpj` (BR), `cuit` (AR) ou `ruc` (PY); pessoa aceita `cpf`, `cuil`, `dni`, `ruc` ou `ci`, de **qualquer** país suportado; fornecedor, locatário e cliente atendido aceitam qualquer tipo de qualquer país. Forma conferida no banco e no Go; **dígito verificador** conferido no Go e no TypeScript pelos mesmos casos JSON para CNPJ (inclusive o alfanumérico, atribuído desde julho de 2026) e CPF; **CUIT, RUC, CUIL, DNI e CI só pela forma** até a **P13** confirmar a regra oficial; nenhum validador consulta serviço externo. Unicidade: da empresa em toda a plataforma por tipo e valor; da pessoa por empresa. Identificação pessoal completa é dado sensível (RNF-021). (4) **Perfil de país:** moeda (RNF-045), idioma padrão sugerido, fuso padrão e lista de fusos válidos (BR 16, AR 12, PY 1, nome IANA, nunca deslocamento fixo), subdivisões em ISO 3166-2 (BR 27, AR 24, PY 18), forma do código postal (CEP; CPA com a letra inicial igual à da província; código de 6 dígitos no Paraguai), regra de "mesma entidade" da conexão entre contas (CNPJ pela raiz de 8 posições; CUIT e RUC inteiros), prefixo do número da OS (`OS-` no Brasil; `OT-` na Argentina e no Paraguai) e, quando entrarem, formatação e feriados nacionais por ano. (5) **Isolamento:** empresas de países diferentes não se alcançam; o teste de vazamento usa uma empresa de cada país, cobre os seis pares e prova que o detector acusa o vazamento. (6) **Hospedagem** só em São Paulo para os três países, com aditivo de transferência por país (**P14**); `empresa.pais` diz qual aditivo vale | Decisão do dono do produto (etapa 014); ADR-013, Partes A e D |
| RNF-045 | Dinheiro e moedas | (1) **Moeda** por empresa: `BRL`, `ARS` ou `PYG`, dada pelo país na criação e não alterável pela API; **uma moeda por empresa** no MVP, e multimoeda com cotação fora do escopo. (2) **Armazenamento:** valor em `numeric(18,2)` e preço unitário em `numeric(18,6)`, nunca `float` e nunca inteiro em centavos; toda tabela com dinheiro tem a coluna `moeda` ao lado, igual à da empresa, garantida pelo banco (chave estrangeira composta com a empresa). (3) **Casas da moeda:** `BRL` e `ARS` com duas, `PYG` com zero; um valor em `PYG` nunca tem parte decimal. (4) **API:** valor em texto decimal canônico com as casas da moeda (`"1250.40"` em BRL, `"75000"` em PYG), `moeda` na resposta (uma vez por objeto, só leitura); a entrada é lida pela gramática exata do contrato e recusada (`formato_invalido`) se não for exatamente representável, **sem arredondar**; `-0` e sinal onde não cabe são recusados. (5) **Cálculo:** um único arredondamento, meio para longe do zero, na multiplicação; soma e subtração só de mesma moeda; **rateio pelo maior resto**, com a soma das partes sempre igual ao total e o empate para o menor índice; preço unitário não se soma. (6) **Mesmas regras no Go e no TypeScript**, provadas pelos mesmos casos JSON, gerados por um oráculo de aritmética exata, com o CI falhando se o arquivo divergir do gerador; sem biblioteca decimal nova. (7) **Exibição** pela moeda da empresa e pelo idioma do usuário, com o símbolo do país quando é a moeda do país e o código ISO quando não é. (8) **Cobrança do plano:** moeda, preço e meio de pagamento na Argentina e no Paraguai são a **P15** (RF-ADM-009) | Decisão do dono do produto (etapa 014); ADR-013, Partes A e F |
| RNF-046 | Conteúdo legal por país e jurisdição | (1) **Nenhum prazo legal é constante no código**, nem o do Brasil: tipos de documento, normas, capacitações, prazos, antecedências de alerta e modelos de checklist legais são **dado de pacote por país**, versionado no repositório, validado no CI e copiado para a empresa na criação da conta, sem nunca alterar nem apagar o que a empresa editou. (2) **Jurisdição:** o conteúdo e o prazo carregam a jurisdição (ISO 3166-1 para o nacional, ISO 3166-2 para estado, província ou departamento); o documento registrado guarda a **jurisdição emissora**; o prazo é o da subdivisão, se houver, senão o nacional, senão nenhum, e então a validade é digitada no documento; município não é jurisdição de prazo no MVP. (3) **Prazo não confirmado nasce vazio:** sem fonte oficial ou parecer externo, a empresa preenche e a tela diz que o prazo é dela; regra de condição e efeito que dependa de prazo nasce `inativa`. (4) **O núcleo não cita norma pelo nome:** OS, ATV, app e sincronização falam de "conformidade"; a NR-12 é o primeiro modelo de `registro_conformidade`, com os mesmos sete campos e as mesmas regras; norma só se aplica a ativo de empresa do país dela; modelo de registro novo entra por requisito, código e teste. (5) O bloqueio de energia (chave ou letreiro sob a guarda de quem repara, retirado só por quem o colocou) é núcleo comum aos três países. (6) **Documento fiscal:** só a referência do documento do país; a Tenaz não emite fatura nem remessa em país nenhum (**P13**). Dependem de **P12**, **P13** e do parecer de consultor de cada país; **sem o parecer, a Argentina e o Paraguai nascem sem modelo de registro próprio e sem prazo legal preenchido** | Decisão do dono do produto (etapa 014); ADR-013, Parte C; pesquisa de AR e PY |

#### Critério de aceite e rastreio (RNF-043 a RNF-046)

Cada linha diz o que prova o requisito e onde a prova está. A coluna "Existe" é o que o repositório tem em 13/11/2026, na etapa E01; **não é aceite**: o requisito só conta como entregue com os nove itens da seção 3.1 do [roadmap](6_ROADMAP.md), inclusive o aceite do gerente de produto em homologação. A coluna "A fazer" é o que o ADR-013 prevê e ainda não tem código.

| Requisito | Critério de aceite verificável | Existe (código e teste) | A fazer |
|---|---|---|---|
| Aceite de RNF-043 | Para cada rota com erro, `title`, `detail` e `erros[].mensagem` saem no idioma negociado, na ordem do item (4), e o `codigo` é o mesmo nos dois; `Content-Language` e `Vary` presentes; usuário e empresa definem o idioma quando não há `Accept-Language`; e-mail de redefinição no idioma da requisição; catálogo com chave, variável, plural ou termo proibido em falta quebra o build; texto livre no código Go quebra o build; a web mostra a conta argentina em espanhol e a brasileira em português | `backend/testes/integracao/idioma_test.go` (`TestIdiomaDoErroPeloAcceptLanguage`, `TestIdiomaDoUsuarioEDaEmpresaNaNegociacao`, `TestNomeDaTelaNoIdioma`, `TestTextosMigradosNoIdioma`); `backend/internal/modulos/adm/multipais_integracao_test.go` (`TestIdiomaDoUsuario`, `TestRedefinicaoDeSenhaNoIdiomaDaRequisicao`); `backend/internal/plataforma/textos/textos_test.go` (`TestConferirCatalogos`, `TestAcceptLanguage`, `TestOrdemDaNegociacao`, `TestMarcarResposta`); `backend/internal/plataforma/erros/trava_test.go` (`TestSemTextoLivreNoCodigoDeProducao`, `TestTravaAcusaTextoLivre`); `pacotes/i18n/src/es-419.test.ts`; `web/e2e/multipais.spec.ts`; `backend/testes/vazamento/pais_test.go` (`TestUsuarioExternoSoAlteraNomeTelefoneEIdioma`); `app/src/i18n/contexto.tsx` e `app/src/acesso/versao.test.ts` (app em `es-419`) | Documento gerado (PDF, exportação, trilha) e notificação por push no idioma de quem pede ou do destinatário, nas etapas que os criam (E04 em diante); idioma do usuário no banco local do app (E01, `tenaz-mobile-offline`; o idioma do aparelho por `expo-localization` já existe); revisão por falante nativo dos termos marcados no glossário, antes da produção; plural das mensagens com `{{min}}`, `{{max}}` e `{{dias}}` (etapa 014, não bloqueante) |
| Aceite de RNF-044 | Criar conta em cada país com a identificação do tipo certo e recusar tipo de outro país, tipo errado, dígito de CNPJ ou CPF inválido, fuso de outro país, subdivisão fora da lista e código postal de forma errada; `pais` e `moeda` não mudam; os casos de identificação, fuso e código postal rodam no Go e no TypeScript; as empresas A (Brasil), B (Argentina) e C (Paraguai) não se alcançam nos seis pares, e o detector acusa o vazamento quando a RLS é desligada | `backend/testes/casos/identificacao.json` (versão 2); `backend/internal/plataforma/pais/pais_test.go` (`TestValidarIdentificacao`, `TestPerfis`, `TestFusosDoPais`, `TestSubdivisoes`, `TestValidarCodigoPostal`, `TestCPAConfereComSubdivisao`); `backend/internal/modulos/adm/multipais_integracao_test.go` (`TestCriarContaPorPais`, `TestCriarContaRecusaPaisETipoErrados`, `TestPaisEMoedaNaoMudam`, `TestIdentificacaoDeEmpresaUnidadeEPessoa`, `TestFusoValidadoContraOPais`, `TestEnderecoDaUnidadeConferidoPeloPais`); `backend/testes/vazamento/pais_test.go` (`TestEmpresasDePaisesDiferentesNaoSeAlcancam`, `TestDetectorAcusaPaisesQueSeAlcancam`, `TestPaisEMoedaDaEmpresaNaoMudam`, `TestDetectorAcusaPaisMutavel`); `pacotes/dominio/src/dominio.test.ts`; `web/e2e/multipais.spec.ts` | Prefixo `OS-` ou `OT-` no número da OS (E04, RF-OS-001); formatação e feriados nacionais no perfil de país; validador de dígito do CUIT e do RUC depois da **P13**; "mesma entidade" por país na conexão entre contas (E07, RF-SRV-014); aditivo de transferência de dados por país (**P14**) |
| Aceite de RNF-045 | Os 266 casos compartilhados (construção da moeda, leitura, texto canônico, soma, multiplicação com um arredondamento, rateio pelo maior resto) passam no Go e no TypeScript em BRL, ARS e PYG; o CI falha se `dinheiro.json` divergir do gerador; o banco recusa `moeda` diferente da empresa; a API devolve o custo por hora em texto decimal com a moeda e recusa o que não é exatamente representável | `backend/testes/casos/dinheiro.json` e `gerar_dinheiro.py` (tarefa `conferir-casos`); `backend/internal/plataforma/dinheiro/dinheiro_test.go`; `pacotes/dominio/src/dinheiro.test.ts`; `backend/internal/modulos/adm/multipais_integracao_test.go` (`TestCustoHoraComAMoedaDaEmpresa`); `backend/testes/integracao/idioma_test.go` (`TestCustoHoraNaApi`); `backend/testes/vazamento/pais_test.go` (`TestBancoRecusaUsuarioEmMoedaDeOutraEmpresa`, `TestFixturesGravamUsuarioNaMoedaDaEmpresa`) | Toda coluna de dinheiro das etapas seguintes nasce com `moeda` e a chave composta (documento 2, seção 1.4); teste de indicador em dinheiro com a moeda; cobrança do plano em `ARS` e `PYG` (**P15**, E10); teste de exibição do guarani sem casas na web |
| Aceite de RNF-046 | Pacote de país aplicado na criação da empresa, idempotente (reaplicar não altera nem apaga o que a empresa editou); nenhum prazo legal como constante no Go nem no TypeScript, apanhado por teste; prazo sem fonte nasce vazio; documento guarda a jurisdição emissora; o registro da NR-12 continua com os sete campos e as mesmas regras depois da troca para `registro_conformidade` | Nenhum ainda: o conteúdo legal entra com CON (E08, RF-CON-001 a RF-CON-006). A pesquisa que o sustenta existe em [mercosul_argentina.md](pesquisa/mercosul_argentina.md), [mercosul_paraguai.md](pesquisa/mercosul_paraguai.md) e [enderecos_br_ar_py.md](pesquisa/enderecos_br_ar_py.md) | Tudo, em E08: pacote `BR` com os prazos hoje em `CHECK` e em comentários; `registro_conformidade` no lugar de `registro_nr12`; pacotes de AR e PY só depois da **P12** |

Pontos de atenção do escopo, registrados na seção 8.3: o desenho de formatação do ADR-013 (Parte A, item 7: funções próprias, nunca o `Intl` com `es-419`) e o que a web faz hoje (`Intl` com a localidade do país da moeda) divergem (ponto 48); e o número da OS com prefixo pelo país ainda não está no código (ponto 49).

---

## 6. Resumo por módulo e fase

Contagem de requisitos funcionais por fase e por plano mínimo. As colunas de plano somam o mesmo total das colunas de fase. RF-EST-010 e RF-FRO-012, movidos para TER, não entram na contagem. A coluna Serviços conta os requisitos do plano das contas de subcontratadas (sumário 2.2); a parte de RF-SRV-014 a RF-SRV-019 executada na conta do cliente exige dele o Essencial (proposta), sem mudar a contagem.

| Módulo | MVP | Fase 2 | Fase 3 | Total | Essencial | Profissional | Enterprise | Serviços |
|---|---|---|---|---|---|---|---|---|
| ADM | 13 | 3 | 1 | 17 | 13 | 2 | 2 | 0 |
| ATV | 12 | 0 | 1 | 13 | 12 | 0 | 1 | 0 |
| MED | 8 | 0 | 0 | 8 | 8 | 0 | 0 | 0 |
| SOL | 6 | 0 | 0 | 6 | 6 | 0 | 0 | 0 |
| OS | 18 | 2 | 1 | 21 | 18 | 3 | 0 | 0 |
| PRE | 9 | 1 | 2 | 12 | 10 | 0 | 2 | 0 |
| CHK | 6 | 0 | 0 | 6 | 6 | 0 | 0 | 0 |
| EST | 5 | 13 | 0 | 18 | 5 | 13 | 0 | 0 |
| TER | 5 | 4 | 0 | 9 | 5 | 4 | 0 | 0 |
| SRV | 21 | 0 | 0 | 21 | 0 | 0 | 0 | 21 |
| LOC | 0 | 9 | 0 | 9 | 0 | 9 | 0 | 0 |
| FRO | 1 | 10 | 0 | 11 | 1 | 10 | 0 | 0 |
| CON | 8 | 7 | 1 | 16 | 9 | 7 | 0 | 0 |
| IND | 6 | 3 | 3 | 12 | 6 | 5 | 1 | 0 |
| INT | 3 | 5 | 4 | 12 | 1 | 7 | 4 | 0 |
| IA | 0 | 4 | 3 | 7 | 0 | 6 | 1 | 0 |
| MOB | 8 | 2 | 0 | 10 | 8 | 2 | 0 | 0 |
| PRJ | 0 | 7 | 0 | 7 | 0 | 7 | 0 | 0 |
| CTR | 0 | 5 | 0 | 5 | 0 | 5 | 0 | 0 |
| **Total** | **129** | **75** | **16** | **220** | **108** | **80** | **11** | **21** |

Mudanças em relação à versão 0.3 (126 / 62 / 16 = 204):

- **MVP +3:** RF-ADM-016 (regras configuráveis), RF-ATV-013 (fase de ciclo de vida) e RF-OS-021 (campos extras da OS), todos no Essencial.
- **Fase 2 +13:** RF-ADM-017 (pacotes de segmento, Essencial), RF-PRJ-001 a RF-PRJ-007 e RF-CTR-001 a RF-CTR-005 (Profissional).
- **Sem mudar contagens:** RF-ATV-002, RF-ATV-006, RF-SOL-001, RF-OS-001, RF-OS-002 e RF-CON-005 receberam a fase de ciclo de vida, a origem do bloqueio, as origens `projeto` e `regra`, os campos extras, a causa externa, o tipo `implantacao` e o tipo de documento "verificação metrológica".

Mudanças em relação à versão 0.2 (105 / 62 / 16 = 183):

- **MVP +21:** entra o módulo SRV inteiro (RF-SRV-001 a RF-SRV-021). Nenhum requisito existente mudou de fase ou de plano; os existentes que o SRV reaproveita (OS, apontamento, ativos, app, TER) receberam regras para a conta da subcontratada e para a conexão.
- **Decisões do sumário 3.6** aplicadas sem mudar contagens: RF-ADM-005, RF-SOL-001, RF-SOL-002, RF-OS-001, RF-OS-003, RF-OS-007, RF-OS-013, RF-OS-014, RF-EST-015 a RF-EST-017, RF-TER-001, RF-TER-005, RF-LOC-002 a RF-LOC-004, RF-LOC-006, RF-LOC-008, RF-FRO-007, RF-FRO-011, RF-CON-004, RF-CON-007, RF-CON-011, RF-MOB-001.

Mudanças da versão 0.1 para a 0.2 (99 / 43 / 16 = 158 → 105 / 62 / 16 = 183):

- **MVP +6:** entram RF-ADM-015 (pessoas), RF-TER-002 a RF-TER-005 (TER básico), RF-EST-006 (reserva, parte do estoque básico) e RF-CON-014 (assinatura avançada); sai RF-EST-005 (inventário, parte do estoque completo). RF-EST-010 virou RF-TER-001, sem mudar a contagem do MVP.
- **Fase 2 +19:** entram RF-EST-012 a RF-EST-019 (8), RF-TER-006, RF-TER-007 e RF-TER-009 (3), RF-LOC-001 a RF-LOC-009 (9) e RF-EST-005; saem RF-EST-006 e RF-CON-014 para o MVP. RF-FRO-012 virou RF-TER-008, sem mudar a contagem da Fase 2.
- **Planos:** RF-PRE-011, RF-IND-012 e RF-INT-011 passam de Profissional (proposta) para Enterprise, conforme a tabela de recursos do sumário.

Requisitos não funcionais: 46 (RNF-001 a RNF-046), válidos desde o MVP, exceto onde a meta indica outra data (SOC 2, ISO 27001, SLA do Enterprise) ou o módulo da Fase 2 a que se aplicam (RNF-038 e a parte de RNF-037 sobre almoxarifado em terceiro). Novos na versão 0.7: RNF-043 (idiomas suportados e escolha do idioma), RNF-044 (países suportados e perfil de país), RNF-045 (dinheiro e moedas) e RNF-046 (conteúdo legal por país e jurisdição), do escopo multipaís e bilíngue (sumário 0.8, ADR-013); RNF-046 só tem código e teste a partir da etapa E08, com CON. Novos na versão 0.3: RNF-041 (isolamento entre contas conectadas) e RNF-042 (integridade e atraso da sincronização entre contas), ambos do risco R5. Novos na versão 0.2: RNF-037 (isolamento do prestador), RNF-038 (isolamento do locatário), RNF-039 (controle de versão e versão visível na interface) e RNF-040 (histórico de versões sempre atualizado).

---

## 7. Pendências citadas

| Pendência | Requisitos afetados |
|---|---|
| P1 — assinatura ICP-Brasil para documentos de SST (NR-1, item 1.6) | RF-OS-014, RF-CON-002, RF-CON-010, RF-CON-014, RF-CON-015, RF-SRV-012 (OS de máquina NR-12 de cliente conectado), RNF-034 |
| P2 — tabela oficial de validade de CIV e CIPP | RF-CON-005 |
| P3 — API Senatran ou Serpro para multas | RF-FRO-010, RF-INT-009 |
| P4 — validação de preço e de prioridades | RF-ADM-009, RF-ADM-012, RF-IA-001, RF-SRV-001 (preço do plano Serviços e desconto de canal), RF-SRV-019 e todos os planos marcados "(proposta)" |
| P5 — disponibilidade do nome Tenaz | Nenhum requisito funcional; afeta o roadmap |
| P6 — prazos NR-13 com SPIE, categorias III a V | RF-CON-009 |
| P7 — local da cópia de recuperação de desastre | RNF-004, RNF-005, RNF-006 (opções na seção 5.1) |
| P8 — CFOPs e tratamento fiscal de remessa e retorno a subcontratada, envio para conserto e locação | RF-EST-015, RF-EST-018, RF-EST-019, RF-LOC-002, RF-LOC-009 |
| P9 — regras metrológicas de medidores de velocidade (verificação, reparo, lacres) | RF-ADM-017 (pacote `fiscalizacao_eletronica`), RF-CON-005 |
| P10 — acesso do cliente final de um contrato | RF-CTR-005 |
| Serviços que custam por natureza (sumário 3.6, "Ferramentas e recursos"): lojas de aplicativo, e-mail transacional, SMS (**P11**), IA e pagamento | Decididos na etapa em que cada um entra, não antes: lojas de aplicativo em RF-MOB-001 e RNF-027; e-mail transacional em RF-ADM-001, RF-ADM-004, RF-ADM-006, RF-ADM-012 e RNF-035; SMS em RF-ADM-004, RF-ADM-006 e RF-MOB-001; IA em RF-IA-001 a RF-IA-007; pagamento em RF-ADM-009 e RF-SRV-001. O resto do sistema usa só software gratuito ou open source |
| P11 — provedor de SMS para o login por telefone e código | RF-ADM-004, RF-ADM-006, RF-ADM-015, RF-MOB-001, RF-SRV-003 (até a decisão, login só por e-mail e senha, com MFA por TOTP, ou SSO) |
| P12 — normas e documentos obrigatórios da Argentina e do Paraguai (equivalentes da NR-12 e da NR-13, inspeção de vasos de pressão por jurisdição, documentos de veículo e de condutor, multas) | RNF-046; notas de país dos módulos CON e FRO; RF-CON-005, RF-FRO-010 |
| P13 — fiscal da Argentina e do Paraguai (identificação fiscal CUIT, CUIL, DNI, RUC e CI; fatura eletrônica ARCA e SIFEN no lugar da NF-e; equivalentes do CFOP) | RNF-044 (dígito verificador de CUIT e RUC), RNF-046; nota de país do módulo EST; RF-EST-009, RF-EST-015, RF-EST-019 |
| P14 — proteção de dados da Argentina (Ley 25.326) e do Paraguai (Ley 7.593/2025) e a transferência internacional para a nuvem em São Paulo | RNF-006, RNF-021, RNF-044 (aditivo por país); RF-SRV-014 (termo da conexão entre contas de países diferentes) |
| P15 — preço, moeda de cobrança e meio de pagamento na Argentina e no Paraguai; impacto do escopo multipaís no prazo e no time do MVP | RNF-045; RF-ADM-009, RF-SRV-001; documento 6 (seção 2.3) |

---

## 8. Pontos a alinhar com o sumário executivo

### 8.1 Resolvidos pelo sumário 0.3

| # | Ponto da versão 0.1 | Como foi resolvido |
|---|---|---|
| 1 | Webhooks fora do escopo do MVP | Sumário 4 inclui webhooks no MVP, plano Profissional (RF-INT-002), com eventos de nome completo da entidade |
| 2 | Códigos de falha como recurso só do Profissional | Catálogo padrão no Essencial e no MVP (RF-OS-011); personalizável no Profissional, Fase 2 (RF-OS-012) |
| 3 | Origens de leitura não cobriam API de integrador, checklist e abastecimento | `manual` inclui checklist e abastecimento digitado; `importacao` inclui CSV e API; `iot` na Fase 3 (RF-MED-002) |
| 4 | Ativo voltar a operar só na `encerrada` inflaria MTTR e Ao | O ativo volta a `operando` no encerramento, mas os KPIs de tempo usam `data_inicio`, `data_conclusao` e as paradas (RF-ATV-006, RF-OS-005). A opção "ativo liberado para operar" foi retirada |
| 5 | Sem transição de `concluida` para `em_execucao` | Transição de liberação recusada, com motivo (RF-OS-003, RF-OS-014). A "liberação com ressalva" foi retirada |
| 6 | Estados a partir dos quais se cancela | `aberta`, `planejada`, `programada`, `pausada`; nunca `em_execucao`, `concluida`, `encerrada` (RF-OS-003) |
| 7 | Reconhecimento sem estado | Campo `data_reconhecimento` e ação `reconhecer`, também implícito ao programar, iniciar ou assumir (RF-OS-004) |
| 8 | MOPP em documento e em capacitação | MOPP é sempre capacitação; CNH é documento da pessoa (RF-CON-005, RF-CON-012) |
| 9 | Fronteira entre estoque básico e completo | Tabela de recursos do sumário 2.1: reserva passou para o MVP (RF-EST-006); inventário, transferência e vários almoxarifados foram para a Fase 2 (RF-EST-005, RF-EST-012) |
| 10 | Recursos sem plano comercial | A tabela de recursos do sumário 2.1 define a maioria. IoT e condição, servidor MCP, TCO, RCM/FMEA e ISO 55001 são Enterprise; telemática, frota e cartão, relatórios agendados e capacitações que bloqueiam são Profissional; assinatura avançada é Essencial. O que ainda não consta da tabela está em 8.3 |
| 12 | "IND básico" e "IND essencial" | Sumário usa "indicadores essenciais" nas duas seções (RF-IND-001 a 005 e 007) |
| 13 | NR-12 no Essencial e LOTO no Profissional | Registro de bloqueio do item 12.11.3 no Essencial e no MVP (RF-CON-004); LOTO por ponto no Profissional, Fase 2 (RF-CON-011) |
| 14 | Cobrança de gestor ou planejador que executa | Só conta quem tem o perfil `tecnico`; `prestador` que executa conta; link de uso único não conta (RF-ADM-009) |
| 15 | SOC 2 Type II e janela de observação | Controles e janela começam na Fase 2; relatório Type II na Fase 3 (RNF-036) |
| 16 | Mapa da documentação sem o número 8 | Número 8 reservado ao manual de operação |
| — | Prioridade, status do ativo, perfis, pessoa, SLA | Valores e regras do sumário 3.1 e 3.3 aplicados em RF-OS-001, RF-ATV-006, RF-ADM-005, RF-ADM-015 e RNF-001 |

### 8.2 Resolvidos pelo sumário 0.4

| # | Ponto da versão 0.2 | Como foi resolvido |
|---|---|---|
| 22 | O sumário não dizia se o `locatario` usa o app | Só o portal web no lançamento do LOC (sumário 3.6; RF-LOC-008, RF-MOB-001) |
| 23 | A rota `/v1/versao` de RNF-039 fugia da convenção de rotas no plural | Exceção aceita (sumário 3.6; seção 3 e RF-INT-001) |
| — | Permissão de liberação sem nome | `os.encerrar`, no formato `<módulo>.<ação>`, nunca em `prestador` nem `locatario` (RF-ADM-005, RF-OS-003, RF-OS-014) |
| — | Escopo do acesso externo | Valores `interno`, `prestador`, `locatario`, `link` (RF-ADM-005) |
| — | "Assumir" a OS só para OS `aberta` | Técnico assume OS `aberta` ou `planejada` da sua equipe, gravando `data_reconhecimento` (RF-OS-007, RF-SRV-010) |
| — | Token do link de uso único sem posição definida na URL | Token no fragmento (`/l/#<token>`), ativação por `POST /v1/publico/links-execucao/ativar` e sessão curta nas rotas `/atual` (RF-TER-005, RNF-037) |
| — | Desistência e duplicidade na solicitação | `nova` → `recusada` com motivo `desistencia`; `duplicada` é final, com "reabrir como nova"; origem `portal_locatario` (RF-SOL-001, RF-SOL-002) |
| — | Cancelamento de contrato `vigente` sem item entregue (proposta) e destino do ativo tomado devolvido (proposta) | Só `rascunho` é cancelado; `vigente` termina por `encerrar`; sem sublocação; ativo `de_terceiro` devolvido passa a `inativo` (RF-LOC-002, RF-LOC-003) |
| — | Matriz de responsabilidade aberta a categorias criadas pela empresa | Lista fechada: `preventiva`, `corretiva_desgaste`, `avaria_mau_uso`, `pneus`, `deslocamento` (RF-LOC-004) |
| — | Campos da OS e da locação sem nome canônico | `justificativa_responsavel_custo`, `motivo_recusa_liberacao`, valores de `origem` (inclui `prestador`); `local_uso`, `franquia_horas`, `franquia_km` (RF-OS-001, RF-OS-013, RF-OS-014, RF-LOC-002) |
| — | Posto e seguradora como "categoria livre" do fornecedor (proposta) | Papéis `posto_combustivel` e `seguradora` (RF-TER-001, RF-FRO-007, RF-FRO-011) |
| — | Consumo de peça em trânsito | Não pode ser consumida nem reservada antes da confirmação de recebimento (RF-EST-015, RF-EST-016, RF-OS-010) |
| — | Quem decide o destino da peça com defeito (`almoxarife` ou `gestor`) | Permissão `est.decidir_destino`, padrão em `gestor`, `planejador` e `almoxarife` (RF-EST-017) |
| — | Bloqueio do ativo, aptidão e LOTO misturados | `bloqueio_ativo` deixa o ativo `parado`; apto para despacho é verificação calculada; `bloqueio_energia` é o registro LOTO (RF-ATV-006, RF-CON-004, RF-CON-007, RF-CON-011) |

### 8.3 Ainda em aberto

| # | Ponto | Solução adotada aqui |
|---|---|---|
| 11 | Local da cópia de recuperação de desastre (P7) | Opções registradas na seção 5.1, sem decisão; RTO de região só é publicado depois de P7 |
| 17 | O sumário define `em_manutencao` como "retirado de operação de forma programada", mas não nomeia o campo da OS que registra isso | Campo "retira o ativo de operação" nas OS programadas (RF-OS-001, RF-OS-005) (proposta); nome canônico a definir no modelo de dados |
| 18 | A tabela de planos do sumário não cita a oficina externa nem o orçamento de fornecedor | Profissional (proposta), junto com o restante do TER da Fase 2 (RF-TER-008, RF-TER-009) |
| 19 | Efeito do custo do reparo no custo médio quando a peça recondicionada volta ao estoque | Proposta: o custo do reparo entra como custo de entrada da peça e no custo médio do item; alternativa a avaliar: código de item separado para recondicionados (RF-EST-018) |
| 20 | CFOPs e obrigatoriedade de NF-e por tipo de movimentação (P8) | CFOP informativo e obrigatoriedade da nota por parâmetro da empresa até a validação com o contador (RF-EST-019) |
| 21 | Conserto de pneu fica fora da ordem de reparo, porque pneu não é item de estoque | Mantido em RF-FRO-005 (proposta) |
| 24 | O sumário 2.1 não diz que plano o **cliente** precisa ter para convidar uma subcontratada e trocar OS com ela | Essencial (proposta), junto com o TER básico (RF-SRV-014 a RF-SRV-019) |
| 25 | O plano Serviços (sumário 2.2) não cita planos preventivos (PRE) sobre equipamentos de clientes, API, webhooks nem SSO | Fora do plano Serviços no lançamento; a subcontratada que precisar contrata um plano adicional (RF-SRV-001) (proposta). Alternativa a decidir com P4: incluir PRE sobre equipamentos `de_cliente` no Serviços |
| 26 | Conta com plano Serviços e plano adicional: o mesmo técnico seria cobrado duas vezes | Cobrado uma vez, pelo plano de maior preço (RF-ADM-009) (proposta) |
| 27 | Preço do plano Serviços e desconto de canal estão como proposta no sumário | Valores do sumário 2.2 usados como proposta, dependentes de **P4** (RF-SRV-001) |
| 28 | `servico_fornecedor` exige fornecedor; na conta da subcontratada o catálogo é da própria conta | Mesma entidade sem fornecedor vinculado (RF-SRV-006) (proposta); modelo a fechar em 2_ARQUITETURA_BANCO_DADOS.md |
| 29 | O sumário não define a `origem` da OS criada na subcontratada a partir da OS do cliente | `origem = manual` com `vinculo_os.lado_origem = cliente` (RF-SRV-007) (proposta); alternativa: novo valor de origem no sumário |
| 30 | `vinculo_os` com uma linha em cada conta ou um registro na plataforma | Os requisitos valem para as duas opções; a escolha e a justificativa ficam em 2_ARQUITETURA_BANCO_DADOS.md (ficha de alinhamento) |
| 31 | No MVP não há contrato de prestação (Fase 2) para dizer quais ativos do cliente são "cobertos" | O cliente seleciona ativos ou classes na própria conexão; na Fase 2, o contrato vigente passa a definir a lista (RF-SRV-015) |
| 32 | OS aberta pela subcontratada só pode ser executada depois da triagem do cliente, o que atrasa emergências em campo | Execução imediata só de `corretiva_emergencial` e só se a conexão permitir (RF-SRV-017) (proposta) |
| 33 | Capacitações de técnicos (RF-CON-012) são Fase 2, mas o sumário 3.5 as inclui no cadastro de técnicos do SRV (MVP) | Cadastro e alertas de `capacitacao` entregues no MVP para contas com SRV, sem bloqueio de atribuição; RF-CON-012 completa a Fase 2 para as demais contas (RF-SRV-003) (proposta) |
| 34 | Prazos de SLA do cliente atendido não têm nome canônico | Prazos de atendimento e de solução por prioridade no `cliente_servico` (RF-SRV-004, RF-SRV-021) (proposta); nomes a definir no modelo de dados. Na Fase 2, quando houver `contrato_cliente`, os prazos do contrato prevalecem (RF-CTR-001) |
| 35 | Dano causado por terceiro em equipamento fixo (abalroamento, vandalismo, furto de cabos) não tem registro de sinistro: `sinistro` é de veículo (RF-FRO-011) | A causa externa da OS registra o fato e a evidência (RF-CTR-003); a cobrança do custo do causador fica como proposta a detalhar |
| 36 | O documento de design prevê a visão "Mapa" na lista de ativos, sem requisito que a garanta | Necessária para pontos espalhados em cidades; requisito a criar quando a tela de ativos for desenhada (proposta) |
| 37 | Regras configuráveis no plano Serviços (contas de subcontratadas) | Campos extras da OS e regras disponíveis também no plano Serviços (proposta), dependente de P4 |
| 38 | Campo extra da OS que dispara regra no cliente (ex.: "Rompeu lacre?") não chega quando a OS é executada por subcontratada conectada, porque campos extras não trafegam na conexão | Em aberto, para decisão no sumário. Opções: (a) grupo novo da conexão para campos extras marcados pelo cliente como compartilháveis, com revisão do risco R5; (b) o cliente confirma o campo ao liberar a OS; (c) a regra fica restrita a OS executadas sem conexão |
| 45 | A tabela de permissões típicas do documento 9 (seção 5.4) difere dos grupos do documento 3: deixa o `planejador` sem executar OS nem checklist, o `visualizador` sem abrir solicitação e o `tecnico` e o `almoxarife` sem ver indicadores | A matriz de RF-ADM-005 segue o documento 3 (Execução inclui `planejador`; Todos abre solicitação; Leitura vê indicadores). O documento 9 precisa remeter à matriz em vez de manter tabela própria |
| 48 | O ADR-013 (Parte A, item 7) manda formatar número, data e moeda por funções próprias e **nunca** pelo `Intl` com `es-419`; a web (`web/src/lib/formatacao.ts`) usa `Intl`, com a localidade do país da moeda (`es-AR`, `es-PY`) quando a moeda é a de um deles, e o documento 7 (seção 3.11) descreve o uso de `Intl` | Em aberto, para o `tenaz-arquiteto` e o `tenaz-design-sistema`: ou o ADR é emendado para aceitar `Intl` com a localidade do país, ou a web passa às funções de `pacotes/dominio`. Vale para o app (Hermes, `Intl` incompleto) e para o PDF gerado no Go, que precisam do mesmo resultado. O critério de aceite de RNF-029 e de RNF-045 (item 7) é por resultado, não por mecanismo |
| 49 | O prefixo do número da OS pelo país (`OS-` ou `OT-`, ADR-013, Parte A, item 9) está nos requisitos (RF-OS-001, RF-SRV-007, RNF-044) e ainda não tem código | Entra na etapa E04, com o núcleo da OS; o gatilho do prefixo é o `empresa.pais` e o número nunca muda com o idioma |
| 50 | Os documentos 3 (seção 26.2) e 4 ainda citam `cnpj` e `cpf` em fornecedor, locatário, cliente atendido e conexão entre contas | Cada etapa que implementa o cadastro fixa o nome do par `identificacao_tipo` e `identificacao` e atualiza o documento dono; os requisitos deste documento já dizem "identificação" (nota de país dos módulos TER, SRV e LOC) |

### 8.4 Resolvidos pelo sumário 0.6 (etapa E01)

| # | Ponto | Como foi resolvido |
|---|---|---|
| 40 | O `administrador` tem `os.encerrar` e `est.decidir_destino`? | Não: tem todas as permissões, menos `os.encerrar`, `est.decidir_destino` e `ter.manter_catalogo_proprio` (a quarta exceção, `pessoa.saude.ler`, saiu do catálogo em 2026-11-13); quem administra e libera OS recebe também `gestor` ou `planejador`, e o responsável pela conta nasce `administrador` e `gestor` (sumário 3.6; RF-ADM-001, RF-ADM-005). Só o administrador concede as três permissões internas que não tem, sempre a outro usuário (proposta, RN-ADM-026). O documento 9 ainda precisa alinhar a tabela da seção 5.4 e a regra anti-escalada |
| 41 | Permissão avulsa sem lugar no modelo de dados | Entidade `usuario_permissao` (sumário 3.6; RF-ADM-004, RF-ADM-005) |
| 46 | Códigos de permissão fora do sumário | O sumário 3.6 declara canônico o catálogo de RF-ADM-005, sem repetir a lista |
| — | Situação do usuário divergente entre os documentos 1 e 2 | `usuario.situacao`: `convidado`, `ativo`, `inativo`; o bloqueio por falhas de login é da identidade (RF-ADM-004, RF-ADM-006) |
| — | MFA obrigatório por empresa sem campo | `empresa.mfa_obrigatorio` (RF-ADM-006) |
| — | Calendário padrão da unidade | Sem padrão imposto: o responsável informa os turnos na criação; modelos só preenchem o formulário (RF-ADM-001, RF-ADM-002) |
| 39 | Perfis padrão de `srv.gerenciar_conexao`: o documento 3 punha as ações de conexão no grupo Admin | Documento 3, versão 0.5: grupo **Conexão** (`srv.gerenciar_conexao`, padrão `administrador` e `gestor`), com o lado de cada ação conferido pelo servidor |
| 42 | O catálogo de endpoints do documento 3 não nomeava o `prestador` nas rotas de checklist | Documento 3, versão 0.5: `prestador` nomeado nas rotas de modelos e execuções de checklist, com o alcance das OS do seu fornecedor (`chk.ler`, `chk.ler_execucoes`, `chk.executar`) |
| 43 | O `plh` só aparecia no grupo Leitura do documento 3 | Documento 3, versão 0.5: `plh` nomeado em `PUT /v1/ativos/{id}/enquadramento-nr13` (`con.gerenciar_nr13`, Fase 2) |
| 44 | Grupo divergente no registro do orçamento da ordem de reparo | Documento 3, versão 0.5: `registrar-orcamento` e `aprovar` da ordem de reparo usam os grupos de `ter.registrar_orcamento` e `ter.aprovar_orcamento` |
| 47 | Vínculo e situação da `pessoa` sem nome no sumário; `desligada` sem comportamento definido (lacuna L-42 do documento 4) | `pessoa.vinculo`: `empregado`, `terceiro`, `autonomo`; `pessoa.situacao`: só `ativa` ou `inativa`, sem `desligada`; pessoa inativa nunca é excluída e continua ligada ao histórico (sumário 3.6; RF-ADM-015). O documento 2 já tirou `desligada` de `ck_pessoa_situacao` (ponto S39) |
| — | Troca de CNPJ só pelo suporte (proposta do documento 4) | Recusada: o `administrador` altera o CNPJ, com registro na trilha; a troca numa conta com conexões é tratada na etapa da conexão (RF-ADM-001) |


### 8.5 Resolvidos pelo sumário 0.8 (etapa 014, multipaís e idiomas)

| # | Ponto | Como foi resolvido |
|---|---|---|
| — | RNF-029 fixava a interface em português e a moeda R$ | RNF-029 reescrito (formatos por país) e RNF-043 (idiomas e escolha do idioma), sem renumerar |
| — | CNPJ, CPF, CEP, UF e R$ fixos nos requisitos de ADM, TER, SRV e LOC | Par `identificacao_tipo` e `identificacao`, `subdivisao` e `codigo_postal` e `moeda` ao lado do valor (RNF-044, RNF-045); convenção na seção 3; notas de país nos módulos |
| — | Nenhum requisito para país, moeda e conteúdo legal fora do Brasil | RNF-044, RNF-045 e RNF-046, com critério de aceite e rastreio na seção 5.12 |
| — | Normas, fiscal, proteção de dados e preço fora do Brasil | Pendências **P12** a **P15** (seção 7), sem decisão neste documento |
| — | Aviso de dado de saúde em espanhol | Só em português; em espanhol fica como customização se o cliente pedir (RNF-043, item 9; sumário 3.6) |

---

**Histórico de versões**

- v0.1 (12/11/2026): primeira versão, em elaboração.
- v0.2 (12/11/2026): alinhada ao sumário executivo 0.3. Novos módulos TER e LOC; EST estendido (estoque em terceiro, remessa e retorno, peças serializadas, reparáveis, quarentena, ordem de reparo, NF-e de remessa e retorno); RF-ADM-015 (pessoas); RF-EST-010 movido para RF-TER-001 e RF-FRO-012 para RF-TER-008; reserva no MVP e inventário na Fase 2; assinatura avançada no MVP; prioridades, reconhecimento, transições da OS, status do ativo, perfis, cobrança, webhooks, SLA e SOC 2 conforme o sumário; RNF-037 a RNF-039.
- v0.3 (12/11/2026): alinhada ao sumário executivo 0.4. Novo módulo SRV (RF-SRV-001 a RF-SRV-021, MVP, plano Serviços) e conexão entre contas; RNF-041 e RNF-042 (risco R5); decisões do sumário 3.6 aplicadas (`os.encerrar`, escopos, `assumir`, link com token no fragmento, solicitação, locação, matriz, campos da OS, papéis de fornecedor, peça em trânsito, `est.decidir_destino`, bloqueios, locatário só na web); MVP com 126 requisitos e total de 204; seção 8 dividida em resolvidos pelo 0.3, resolvidos pelo 0.4 e em aberto.
- v0.4 (12/11/2026): alinhada ao sumário executivo 0.5. Configuração por segmento: RF-ADM-016 (regras configuráveis), RF-ATV-013 (fase de ciclo de vida), RF-OS-021 (campos extras da OS) e o tipo de OS `implantacao` no MVP; RF-ADM-017 (pacotes de segmento) e os módulos PRJ (RF-PRJ-001 a RF-PRJ-007) e CTR (RF-CTR-001 a RF-CTR-005) na Fase 2; origens `projeto` e `regra`, causa externa da parada, origem do bloqueio de ativo e tipo de documento "verificação metrológica"; MVP com 129 requisitos e total de 220; pendências P9 e P10; pontos 35 a 37 em aberto.
- v0.5 (12/11/2026): alinhada ao sumário executivo 0.6 (etapa E01). Catálogo de 117 permissões e matriz dos doze perfis padrão como subseção de RF-ADM-005, com as regras do `administrador`, dos perfis externos e da conta da subcontratada (`srv.gerenciar_conexao`); calendário da unidade em `turno_unidade` e `feriado_unidade` (RF-ADM-002); login por telefone e código adiado pela P11 (RF-ADM-004, RF-ADM-006, RF-ADM-015, RF-MOB-001, RF-SRV-003); perfis padrão criados com a empresa (RF-ADM-001, RF-SRV-001); escopo de acesso e cadastros de acesso na auditoria (RF-ADM-011); usuário externo só com o perfil do seu tipo (RF-ADM-004); P11 na seção 7. Decisões do sumário 0.6, seção 3.6: catálogo de permissões canônico, regra do `administrador`, responsável `administrador` e `gestor`, turnos da unidade inicial sem padrão imposto, CNPJ único entre todas as empresas com reativação (proposta) e alterável pelo `administrador`, `usuario.situacao`, `usuario_permissao` e RN-ADM-026, `empresa.mfa_obrigatorio`; vínculo e situação da `pessoa`, sem `desligada` (RF-ADM-015). Trilha de auditoria configurável (RF-ADM-011 ampliado, com critérios de aceite; RF-ADM-013 e RNF-019 coerentes). `srv.gerenciar_conexao` no grupo Conexão. Pontos 39, 40, 41, 42, 43, 44, 46 e 47 resolvidos (nova seção 8.4); fica em aberto o 45. Escopo ampliado só em RF-ADM-011, por decisão do dono do produto; sem mudança de fase, de plano ou de contagem.
- v0.6 (13/11/2026): alinhada ao sumário executivo 0.7. A Tenaz não guarda dado de saúde (decisão do dono do produto, etapa E01): ASO, exame toxicológico, a permissão `pessoa.saude.ler` e as regras de dado sensível de saúde removidos do escopo em 2026-11-13; catálogo de RF-ADM-005 com 116 códigos, sem exceção ao formato; `administrador` sem três permissões; ajustes em RF-ADM-005, RF-ADM-011, RF-ADM-015, RF-SRV-003, RF-SRV-014, RF-FRO-001, RF-CON-005, RNF-021 e nas seções 2 e 3. Nenhum ID apagado nem renumerado; sem mudança de fase, de plano ou de contagem. Reforço da vedação (sumário 0.7, "Dado de saúde recebido por engano"): aviso fixo na pessoa (RF-ADM-015), no dossiê do motorista (RF-FRO-001) e no envio de anexo (RNF-033); recusa de tipo de documento com nome ligado a saúde, com termos, comparação por palavra inteira e erro `422 tipo_documento_saude_vedado` propostos (RF-CON-005). Serviços que custam por natureza registrados na seção 7. Respostas do dono do produto: texto do aviso aprovado; `exame`, `exames`, `clinico` e `clinica` fora da lista de termos (entram `exame medico` e `exames medicos`); nome do arquivo não verificado; aviso de anexo em E02; tipos vedados já existentes só listados para revisão do administrador em E08. Armazenamento imutável do selo de auditoria (RF-ADM-011): SeaweedFS em desenvolvimento, bucket `tenaz-auditoria`, e S3 da AWS com Object Lock em produção, ambos em modo de conformidade (ADR-012, ADR-006); o MinIO saiu do documento.
- v0.7 (13/11/2026): alinhada ao sumário executivo 0.8 (etapa 014, dentro da E01). Multipaís e bilíngue desde o MVP: RNF-043 (idiomas suportados e escolha do idioma), RNF-044 (países suportados e perfil de país), RNF-045 (dinheiro e moedas) e RNF-046 (conteúdo legal por país e jurisdição), na seção 5.12, com critério de aceite e rastreio para os testes existentes; RNF-029, RNF-006, RNF-013 e RNF-021 ajustados; RF-ADM-001, RF-ADM-002, RF-ADM-004, RF-ADM-006, RF-ADM-009, RF-ADM-010, RF-ADM-011, RF-ADM-012, RF-ADM-013, RF-ADM-015, RF-OS-001, RF-TER-001, RF-SRV-001, RF-SRV-003, RF-SRV-004, RF-SRV-007, RF-SRV-012, RF-SRV-014 e RF-LOC-001 para a forma por país; notas de país nos módulos EST, TER, SRV, LOC, FRO e CON e em RF-EST-009, RF-EST-015, RF-FRO-010, RF-INT-009 e RF-CON-005, apontando P12 a P15; convenção de país, idioma e dinheiro na seção 3; P12 a P15 na seção 7; seção 8.5 (resolvidos) e pontos 48 a 50 em aberto na 8.3. 46 requisitos não funcionais; contagem dos 220 RF, fases e planos sem mudança; nenhum ID apagado nem renumerado.
