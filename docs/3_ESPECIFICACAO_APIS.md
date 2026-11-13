# Especificação de APIs — Tenaz

**Versão:** 0.6 · **Data:** 13/11/2026 · **Status:** Em elaboração
**Alinhado ao sumário executivo 0.7.**
**Produto:** Tenaz · Gestão de Manutenção e Frota
**Protocolo:** HTTPS (TLS 1.2 ou superior, preferência por 1.3) · **Formato:** JSON · **Contrato:** OpenAPI 3.1

**Mudanças da versão 0.4:** configuração por segmento (sumário 0.5, seção 3.7). Entram no MVP as regras configuráveis de condição e efeito (`/v1/regras-automacao`, `/v1/disparos-regra`), os campos extras da OS (`/v1/campos-extras-os` e `campos_extras` na OS), a fase de ciclo de vida do ativo (`fase_ciclo_vida`, `data_entrada_operacao`, ação `colocar-em-operacao`, unitária e em lote) e o tipo de OS `implantacao`. Entram na Fase 2 os pacotes de segmento (`/v1/pacotes-segmento`, `/v1/aplicacoes-pacote`) e dois módulos novos: **PRJ** (`/v1/modelos-projeto`, `/v1/projetos`, etapas e `/v1/medicoes-etapa`) e **CTR** (`/v1/contratos-cliente`, itens, painel e `/v1/medicoes-contrato`). A OS ganha a origem `projeto` (e passa a usar `cliente_conectado`, que o sumário já tinha), `causa_externa` e `etapa_projeto_id`; a solicitação ganha a origem `regra`; o bloqueio de ativo ganha `origem`. Nove permissões novas, 9 eventos de webhook novos (um deles proposta) e 24 códigos de erro novos. Catálogo de 482 para 564 endpoints. Seções novas 23 a 25; as pendências passam para a seção 26.

**Mudanças da versão 0.6:** a Tenaz não guarda dado de saúde (sumário 0.7, seção 3.6, "Dados de saúde"). Sai a permissão `pessoa.saude.ler`, e com ela a única exceção ao formato `<módulo>.<ação>`; saem as menções a documentos de saúde na pessoa, na agenda de documentos, nos escopos de token, na exportação e nos exemplos. A seção 4.8 passa a listar doze permissões finas.

**Mudanças da versão 0.5:** contrato de E01 (conta e acesso) implementado e gerado do código. Códigos de erro: o código Go passa a usar exatamente os desta seção 3.2 (snake_case), e entram `senha_fraca`, `mfa_codigo_invalido`, `link_email_invalido`, `ultimo_administrador` e `pre_condicao_obrigatoria`; o problema traz `desafio_id`, `desafio_expira_em`, `tentativas_restantes` e `versao_atual`. Rotas novas de ADM: `POST /v1/empresas` (cadastro pelo site), convites (`/v1/auth/convites/consultar` e `/aceitar`), MFA do próprio usuário (`/v1/auth/mfa/iniciar`, `/confirmar`, `/desativar`), `POST /v1/auth/trocar-senha`, `POST /oauth/authorize` (decisão da tela de autorização), turnos e feriados da unidade, reativação de unidade e de pessoa, reenvio de convite e permissões avulsas do usuário; `DELETE /v1/unidades/{id}` passa para o MVP. Cookie de desenvolvimento, conferência de Origin, fluxo do authorize, `empresa_escolhida_id` e extensões `x-escopos-acesso` e `x-permissoes` no contrato. Divergências do documento 1 (seção 8.3, pontos 39, 42, 43 e 44): grupo **Conexão** (`srv.gerenciar_conexao`), `prestador` nas rotas de checklist, `plh` no enquadramento NR-13 e grupo da rota de orçamento da ordem de reparo. Catálogo de 564 para 583 endpoints.

Este documento especifica a API REST da Tenaz: convenções, autenticação, limites, catálogo de endpoints por módulo, sincronização do app offline, webhooks, integrações, terceiros, estoque estendido, locação, gestão de serviços do prestador com conexão entre contas, configuração por segmento, projetos de implantação, contratos com clientes e política de versões. Usa os códigos de módulo, os nomes de entidade e os estados canônicos do [Sumário Executivo](0_SUMARIO_EXECUTIVO.md). Quando este documento e o sumário divergirem, vale o sumário, e a divergência deve ser registrada na seção 26.

Tudo marcado como **proposta** é decisão de projeto ainda não validada. O restante decorre do sumário executivo ou do [estudo de mercado](pesquisa/reports/Benchmarking%20manuten%C3%A7%C3%A3o%20e%20frota%20Brasil.md). As regras de negócio citadas aqui são resumidas; a fonte definitiva será o [4_FLUXOS_PROCESSO.md](4_FLUXOS_PROCESSO.md), e os nomes de tabela e campo serão confirmados no [2_ARQUITETURA_BANCO_DADOS.md](2_ARQUITETURA_BANCO_DADOS.md).

---

## Índice

1. [Visão geral](#1-visão-geral)
2. [Princípios e convenções](#2-princípios-e-convenções)
3. [Erros](#3-erros)
4. [Autenticação e autorização](#4-autenticação-e-autorização)
5. [Rate limit](#5-rate-limit)
6. [Catálogo de endpoints por módulo](#6-catálogo-de-endpoints-por-módulo)
7. [Solicitações: detalhamento](#7-solicitações-detalhamento)
8. [Ordens de serviço: detalhamento](#8-ordens-de-serviço-detalhamento)
9. [Medidores: registrar leitura](#9-medidores-registrar-leitura)
10. [Checklists: executar](#10-checklists-executar)
11. [Indicadores: consultar KPIs](#11-indicadores-consultar-kpis)
12. [Sincronização do app offline](#12-sincronização-do-app-offline)
13. [Webhooks](#13-webhooks)
14. [Integrações de entrada](#14-integrações-de-entrada)
15. [Integrações de saída: conectores de ERP](#15-integrações-de-saída-conectores-de-erp)
16. [Bloqueio de despacho de veículo](#16-bloqueio-de-despacho-de-veículo)
17. [Exportação completa dos dados da empresa](#17-exportação-completa-dos-dados-da-empresa)
18. [Versionamento, depreciação e compatibilidade](#18-versionamento-depreciação-e-compatibilidade)
19. [Terceiros: prestadores, OS atribuída e link de uso único](#19-terceiros-prestadores-os-atribuída-e-link-de-uso-único)
20. [Estoque estendido: estoque em terceiros, peças serializadas e conserto](#20-estoque-estendido-estoque-em-terceiros-peças-serializadas-e-conserto)
21. [Locação: contratos, vistorias e responsabilidade de custo](#21-locação-contratos-vistorias-e-responsabilidade-de-custo)
22. [SRV: gestão de serviços do prestador e conexão entre contas](#22-srv-gestão-de-serviços-do-prestador-e-conexão-entre-contas)
23. [Configuração por segmento: regras, campos extras, fase de ciclo de vida e pacotes](#23-configuração-por-segmento-regras-campos-extras-fase-de-ciclo-de-vida-e-pacotes)
24. [PRJ: projetos de implantação](#24-prj-projetos-de-implantação)
25. [CTR: contratos com clientes](#25-ctr-contratos-com-clientes)
26. [Pendências e decisões em aberto](#26-pendências-e-decisões-em-aberto)

---

## 1. Visão geral

### 1.1 Uma API para todos os clientes

A Tenaz tem **uma única API**. O web (React), o app de campo (React Native) e os integradores dos clientes usam os mesmos endpoints. O que muda entre eles é o método de autenticação (seção 4) e o limite de requisições (seção 5).

| Cliente | Autenticação | Uso |
|---|---|---|
| Web | Sessão por cookie httpOnly + token CSRF | Toda a interface administrativa e operacional |
| App de campo | OAuth 2.1 com PKCE + refresh token rotativo | Execução offline, sincronização |
| Integração do cliente | Token de API da empresa, com escopos | ERP, TMS, BI, scripts |
| Provedor de identidade | Token SCIM dedicado | Provisionamento de usuários (Enterprise) |
| Prestador ocasional | Link seguro de uso único, sem login: token no fragmento da URL, trocado por uma sessão curta (19.4) | Registro da execução de uma única OS atribuída ao seu fornecedor |

Usuários de prestadores (perfil `prestador`) e de locatários (perfil `locatario`) usam o web e o app com a mesma autenticação dos demais, mas com alcance restrito (4.9).

A **conta de uma subcontratada** (módulo SRV, seção 22) é uma `empresa` como outra qualquer: seus usuários, técnicos e integradores usam a mesma API, com as mesmas credenciais. A troca de OS entre a conta do cliente e a da subcontratada conectada não passa por nenhum cliente da API: é feita por um processo interno de sincronização (22.6), e nenhum endpoint lê dados das duas contas.

**Disponibilidade por plano.** Web e app funcionam em todos os planos, inclusive no plano Serviços (sumário 2.2). **Tokens de API e webhooks** exigem o plano **Profissional** ou superior, conforme a seção 2.1 do sumário; o plano Serviços não os inclui (22.10). No catálogo (seção 6), a coluna "Plano" indica o plano mínimo do **recurso de negócio**; o acesso por token de API exige, além disso, o Profissional.

### 1.2 Endereços e ambientes (proposta)

| Ambiente | URL base | Dados |
|---|---|---|
| Produção | `https://api.tenaz.com.br/v1` | Reais, região São Paulo |
| Sandbox | `https://api.sandbox.tenaz.com.br/v1` | Fictícios, empresa de teste por cliente, reiniciável |

O domínio depende da pendência P5 do sumário (disponibilidade do nome Tenaz). O web chama a API pelo mesmo domínio do aplicativo (`https://app.tenaz.com.br/api/v1`, proposta) para que o cookie de sessão seja de primeira parte.

### 1.3 OpenAPI 3.1 como fonte da verdade

- O contrato OpenAPI 3.1 é **gerado** a cada build a partir dos tipos Go anotados das rotas (biblioteca Huma) e versionado em `backend/api/openapi.yaml`. O CI falha se o arquivo gerado divergir do versionado; a revisão do contrato acontece no pull request, sobre a diferença desse arquivo (seção 2.5 do [documento 5](5_ARQUITETURA_SOFTWARE.md)).
- O servidor valida toda requisição contra o esquema da rota antes do caso de uso; o erro de validação sai em `problem+json` (seção 3). O processo (`tenaz api`), a geração do contrato (`tenaz openapi`) e os testes de rota montam a API pelo mesmo ponto (`internal/servidor`), para que o contrato seja exatamente o que o servidor atende.
- Os tipos TypeScript do web e do app são gerados do contrato (proposta: `openapi-typescript`). Nenhum tipo de API é escrito à mão no frontend.
- O CI roda lint do contrato (proposta: Spectral, com as regras deste documento: snake_case, rotas no plural em kebab-case, `operationId` obrigatório, exemplos obrigatórios) e detecta mudança incompatível contra a última versão publicada (proposta: `oasdiff`). Mudança incompatível em `/v1` quebra o build.
- A documentação pública para integradores é gerada do mesmo arquivo e publicada com exemplos em curl, Go e TypeScript.
- Extensões próprias do contrato: `x-modulo` (código do módulo), `x-plano-minimo`, `x-fase`, `x-perfis`, `x-escopos` (escopos de token de API, 4.4), `x-escopos-acesso` (escopos de acesso aceitos pela rota, 4.9: `interno`, `prestador`, `locatario`, `link`, ou `publico` para rota sem credencial), `x-permissoes` (permissões de RF-ADM-005 exigidas pela rota), `x-estabilidade` (`estavel`, `beta`, `depreciado`) e `x-enum-extensivel` (ver 18.3). Desde a versão 0.5, toda operação gerada traz `x-modulo`, `x-fase`, `x-plano-minimo` e `x-escopos-acesso`, e as que exigem permissão trazem `x-permissoes`; a suíte de isolamento lê `x-escopos-acesso` para varrer as rotas. `x-escopos` entra com os tokens de API.

---

## 2. Princípios e convenções

### 2.1 REST e recursos

- Recursos são substantivos no plural, em português, sem acento, em kebab-case: `/v1/ordens-servico`, `/v1/leituras-medidor`.
- **Regra de plural (proposta):** a rota é o nome canônico da entidade com a **primeira palavra no plural** e `_` trocado por `-`. Exemplos: `ordem_servico` → `ordens-servico`; `plano_manutencao` → `planos-manutencao`; `classe_ativo` → `classes-ativo`; `leitura_medidor` → `leituras-medidor`; `documento_vencimento` → `documentos-vencimento`; `checklist_modelo` → `checklists-modelo`; `servico_fornecedor` → `servicos-fornecedor`; `contrato_locacao` → `contratos-locacao`; `cliente_servico` → `clientes-servico`; `especialidade_tecnico` → `especialidades-tecnico`; `conexao_empresa` → `conexoes-empresa`; `vinculo_os` → `vinculos-os`; `regra_automacao` → `regras-automacao`; `disparo_regra` → `disparos-regra`; `pacote_segmento` → `pacotes-segmento`; `aplicacao_pacote` → `aplicacoes-pacote`; `modelo_projeto` → `modelos-projeto`; `contrato_cliente` → `contratos-cliente`; `medicao_contrato` → `medicoes-contrato`. Adjetivo concorda com o substantivo: `nota_fiscal_entrada` → `notas-fiscais-entrada`; `peca_serializada` → `pecas-serializadas`; `campo_extra_os` → `campos-extras-os`.
- **Exceção aceita (sumário 3.6):** `GET /v1/versao` é recurso único e público e fica no singular.
- Sub-recursos expressam pertença: `/v1/ordens-servico/{id}/apontamentos`.
- **Transições de estado são sub-recursos de ação**, sempre com `POST` e verbo no infinitivo: `POST /v1/ordens-servico/{id}/iniciar`. O campo `status` nunca é alterado por `PATCH`; a API devolve `422 campo_somente_leitura` se o cliente tentar.
- Identificadores na rota são o `id` (UUID). A exceção são os catálogos globais mantidos pela Tenaz, identificados pelo código estável: `/v1/indicadores/{codigo}`, `/v1/relatorios/{codigo}/gerar` e `/v1/pacotes-segmento/{codigo}`. Números legíveis (`OS-2026-004812`, `SS-2026-000731`, `OR-2026-000045`, `CL-2026-000012`, `PJ-2026-000031`, `CT-2026-000004`) e TAGs servem como filtro: `GET /v1/ordens-servico?numero=OS-2026-004812`, `GET /v1/ativos?tag=CMP-0347`. **Nenhuma rota recebe token ou segredo no caminho nem na query string**: o token do link de uso único vai no corpo de `POST /v1/publico/links-execucao/ativar`, e o código de convite de conexão, no corpo de `aceitar` (19.4 e 22.5). Caminhos e queries aparecem em logs, histórico do navegador e cabeçalho `Referer`; corpos não.

| Método | Uso | Corpo | Sucesso |
|---|---|---|---|
| `GET` | Ler recurso ou lista | — | `200` (ou `304` com `If-None-Match`) |
| `POST` | Criar recurso ou executar ação | JSON | `201` com `Location` (criação), `200` (ação), `202` (assíncrono) |
| `PATCH` | Alterar campos | JSON Merge Patch (RFC 7396) | `200` |
| `DELETE` | Excluir ou inativar | — | `204` |

`PUT` não é usado nos recursos de negócio (proposta). As exceções são SCIM (RFC 7644) e configurações que são substituídas por inteiro, como `PUT /v1/ativos/{id}/enquadramento-nr13`.

### 2.2 JSON

- Campos em `snake_case`, **iguais aos nomes das colunas do banco** (`data_abertura`, `ativo_id`, `motivo_pausa`).
- Valores de enumeração em `snake_case`, português, sem acento, iguais aos estados canônicos: `em_execucao`, `corretiva_emergencial`, `aguardando_peca`.
- Chaves estrangeiras terminam em `_id`; listas de chaves, em `_ids`.
- Campos sem valor vêm como `null`, não são omitidos (exceto quando o cliente pede campos parciais, 2.8).
- `Content-Type: application/json; charset=utf-8`. `PATCH` aceita `application/merge-patch+json` e `application/json` com a mesma semântica.
- Campos comuns a todo recurso (proposta, a confirmar no documento 2):

| Campo | Tipo | Descrição |
|---|---|---|
| `id` | UUID (v7, proposta) | Identificador. O app pode gerá-lo offline (seção 12). |
| `versao` | inteiro | Incrementa a cada alteração. Base do ETag (2.11). |
| `criado_em` | data-hora | Criação no servidor |
| `atualizado_em` | data-hora | Última alteração |
| `criado_por_id` | UUID | Usuário ou token que criou |

- `empresa_id` **não aparece** no corpo das respostas e **nunca é aceito** no corpo ou na URL (ver 4.7): enviado, responde `400 campo_nao_permitido`. O identificador da empresa só aparece como `id` em `GET /v1/empresa`, no envelope de webhooks e, para o seletor de empresa, em `empresas[].id` da sessão (4.2). As rotas que escolhem a empresa da credencial (login, `mfa/verificar`, `trocar-empresa` e `POST /oauth/authorize`) recebem o campo `empresa_escolhida_id`, conferido contra os acessos da pessoa (fora deles, `404`).

### 2.3 Versionamento na URL

A versão principal está no caminho: `/v1`. Só muda com mudança incompatível (seção 18). Mudanças compatíveis entram na mesma versão sem aviso prévio, registradas no changelog. O prefixo da rota é independente da versão do produto, que segue SemVer 2.0.0 e sai em `GET /v1/versao` e no cabeçalho `Tenaz-Versao` (18.8).

### 2.4 Datas e horas

- Data-hora em **ISO 8601, sempre em UTC, com `Z`**: `"2026-11-12T13:45:00Z"`. Precisão de segundos; milissegundos são aceitos e descartados (proposta).
- Data sem hora em `YYYY-MM-DD`: `"2027-03-31"` (vencimento de documento, data de validade).
- O servidor rejeita data-hora sem fuso (`422 data_sem_fuso`).
- A interface exibe a data no fuso da unidade (`unidade.fuso_horario`, padrão `America/Sao_Paulo`). Indicadores agregados por dia usam o fuso da unidade; o parâmetro `fuso_horario` permite outro.

### 2.5 Dinheiro e números decimais

- **Dinheiro é string decimal** com ponto e duas casas: `"1250.40"`. Custos unitários de estoque aceitam até quatro casas: `"12.3456"`. Moeda implícita BRL; o campo `moeda` existe só em documentos fiscais e importações.
- **Proposta:** quantidades de estoque, litros de combustível e valores de medidor também são strings decimais (`"125430.5"`), para não perder precisão em `float` no JavaScript. Contagens inteiras (número de falhas, quantidade de OS) são `integer`.
- Percentuais em indicadores são número decimal de 0 a 100 com até duas casas (`92.35`), por serem só leitura e exibição.

### 2.6 Paginação por cursor

Toda lista é paginada por cursor opaco.

| Parâmetro | Padrão | Máximo | Descrição |
|---|---|---|---|
| `limite` | 50 | 200 (500 em `/v1/sincronizacao`) | Itens por página |
| `cursor` | — | — | Valor de `proximo_cursor` da página anterior |
| `incluir_total` | `false` | — | Quando `true`, devolve `total` (consulta mais cara; limite próprio, seção 5) |

```json
{
  "dados": [ { "id": "0193a1b2-7c4d-7e8f-9a0b-1c2d3e4f5a6b", "numero": "OS-2026-004812" } ],
  "paginacao": {
    "proximo_cursor": "eyJrIjpbIjIwMjYtMTEtMTJUMTM6NDU6MDBaIiwiMDE5M2ExYjIiXX0",
    "tem_mais": true,
    "limite": 50
  }
}
```

- O cursor codifica a ordenação e o último item; é válido por 24 horas (proposta). Cursor inválido ou de outra consulta: `400 cursor_invalido`.
- Mudar filtros ou ordenação exige começar sem cursor.
- **Proposta:** recurso único volta sem envelope (o próprio objeto); listas sempre com `dados` e `paginacao`.
- Listas curtas e fechadas (unidades, feriados da unidade, perfis, catálogo de permissões) voltam inteiras, no mesmo envelope, com `tem_mais = false` e `proximo_cursor = null`. As listas abertas (pessoas, usuários, eventos de auditoria) são paginadas por cursor, com `limite` de 1 a 200 (padrão 50) e ordem `-criado_em` (na trilha, `-ocorrido_em`).

### 2.7 Filtros e ordenação

| Forma | Exemplo | Significado |
|---|---|---|
| Igualdade | `status=aberta` | Campo igual ao valor |
| Lista (OU) | `status=aberta,planejada` | Campo em qualquer dos valores |
| Faixa | `data_abertura[gte]=2026-11-01T00:00:00Z&data_abertura[lt]=2026-12-01T00:00:00Z` | Operadores `gt`, `gte`, `lt`, `lte` |
| Nulo | `responsavel_id[nulo]=true` | Campo sem valor |
| Texto | `busca=compressor` | Busca textual nos campos indicados no contrato de cada recurso |
| Alterados desde | `atualizado_em[gte]=2026-11-12T00:00:00Z` | Sincronização simples de integradores |
| Ordenação | `ordenar=-prioridade,data_abertura` | Campos separados por vírgula; `-` = decrescente |

- Cada recurso declara no OpenAPI os campos filtráveis e ordenáveis. Filtro ou ordenação não declarada devolve `400 filtro_nao_suportado`.
- A ordenação padrão é `-criado_em`. Todo cursor desempata por `id`.
- Filtros por unidade (`unidade_id`) respeitam as unidades a que o usuário tem acesso (4.8): pedir uma unidade sem acesso devolve lista vazia, não erro.

### 2.8 Campos parciais e expansão

- `campos=id,numero,status,ativo_id` devolve só os campos pedidos (mais `id` e `versao`, sempre presentes).
- `expandir=ativo,responsavel` embute o objeto relacionado (resumo) no lugar de só o `_id`. Cada recurso declara as expansões possíveis; profundidade máxima 1 (proposta).
- Os dois podem ser combinados: `campos=id,numero,ativo&expandir=ativo`.

### 2.9 Idempotência (`Idempotency-Key`)

- Todo `POST` aceita o cabeçalho `Idempotency-Key` (UUID ou string de até 64 caracteres). É **obrigatório** para tokens de API e para o app (proposta); o web envia sempre.
- O servidor guarda a chave, o hash do corpo e a resposta por **24 horas** (proposta), por empresa e por credencial.
- Repetição com o mesmo corpo devolve a resposta original com o cabeçalho `Idempotency-Replayed: true`.
- Mesma chave com corpo diferente: `422 idempotencia_conflito`.
- Mesma chave enquanto a primeira requisição ainda está em processamento: `409 idempotencia_em_processamento`, com `Retry-After: 1`.
- A sincronização do app (seção 12) usa o `id_comando` como chave, com retenção maior (90 dias).
- A chave vale por empresa e por sessão da credencial (no app, a sessão do aparelho, que sobrevive à renovação do access token). A repetição devolve o status, o corpo e os cabeçalhos `Content-Type`, `Location` e `ETag` da primeira resposta; resposta 5xx não é guardada, para que a repetição tente de novo. Chave com mais de 64 caracteres ou com caractere não visível: `400 requisicao_invalida`.
- Rotas que emitem credencial (`/v1/auth/*`, `/oauth/*`) e rotas públicas aceitam o cabeçalho e o ignoram: a resposta traz cookie ou token novo e não é guardada para repetição.
- **E01:** o depósito das chaves fica na memória de cada instância da API, o que só é correto com uma instância. O depósito compartilhado é pendência A-37 (seção 26). A obrigatoriedade para o app e para tokens de API ainda não é exigida.

### 2.10 Campos de auditoria e motivo

Toda escrita gera um `evento_auditoria` imutável (quem, o quê, quando, antes e depois, IP, dispositivo). O cabeçalho opcional `X-Tenaz-Motivo` (até 500 caracteres) é gravado no evento, para alterações administrativas que não têm campo de motivo próprio (proposta).

### 2.11 Concorrência otimista (ETag e If-Match)

- Toda resposta de recurso único traz `ETag: "<versao>"`, por exemplo `ETag: "12"`.
- `PATCH`, `DELETE` e **todas as ações de estado** exigem `If-Match` com o ETag lido.
  - Sem `If-Match`: `428 pre_condicao_obrigatoria`.
  - Com ETag desatualizado: `412 versao_desatualizada`, com o recurso atual em `recurso_atual` no corpo do erro, para o cliente mostrar a diferença.
- `GET` com `If-None-Match` igual à versão atual devolve `304 Not Modified`.
- Ações que só **acrescentam** registros (criar apontamento, registrar leitura, comentar) não exigem `If-Match` no recurso pai.
- `If-Match` aceita `"12"`, `W/"12"` ou `12`; valor malformado responde `400 requisicao_invalida`. O `412` traz também `versao_atual`.
- Sub-recurso com versão própria (turno e feriado da unidade) usa o ETag dele, não o da unidade. A concessão e a retirada de permissão avulsa (`/v1/usuarios/{id}/permissoes`) não exigem `If-Match`: acrescentam ou retiram uma linha de `usuario_permissao`, sem editar o usuário. Encerrar sessões e reenviar convite também não: não mudam o estado do usuário.
- **E01:** no `PATCH`, a versão é conferida na gravação, na mesma transação. Nas ações de estado e no `DELETE`, a camada HTTP lê o registro, compara a versão e só então chama a ação; a transição continua validada pela máquina de estados (`409 transicao_invalida`). `If-None-Match` e `304` ainda não estão implementados (pendência A-37).

### 2.12 Exclusão

- Registros com histórico operacional (ativo, OS, item de estoque, peça serializada, medidor, pessoa, usuário, fornecedor, locatário) **não são apagados**: `DELETE` inativa (`ativo=false` ou `status=inativo`) e devolve `204`. A lista padrão omite inativos; `incluir_inativos=true` mostra.
- Registros sem dependentes (rascunho de plano, comentário recém-criado) são apagados de fato.
- OS `encerrada` e evento de auditoria nunca são alterados nem apagados: `409 os_encerrada_imutavel`.
- A exclusão definitiva de dados pessoais por pedido de titular (LGPD) é um fluxo administrativo à parte, descrito no documento 9.

### 2.13 Anexos (fotos, PDFs, assinaturas)

Arquivos não trafegam em base64 no JSON. O fluxo é:

1. `POST /v1/anexos` com `nome_arquivo`, `tipo_conteudo`, `tamanho_bytes` e `sha256`. O cliente pode enviar o `id` (gerado offline). A resposta traz `url_upload` pré-assinada (válida por 15 minutos) e o `id` do anexo.
2. O cliente envia o binário com `PUT` direto na `url_upload`.
3. O cliente referencia o anexo pelo id (`anexos_ids`) na OS, no checklist, no comentário.

Limites (proposta): 25 MB por arquivo; imagens são redimensionadas para 2560 px no lado maior, com o original guardado. Tipos aceitos: JPEG, PNG, HEIC, WebP, PDF, MP4 (até 60 s), XML, CSV, XLSX. Anexo não confirmado em 24 horas é descartado.

### 2.14 Operações assíncronas

Operações demoradas (importação, exportação, relatório, reprocessamento) devolvem `202 Accepted` com `Location: /v1/processamentos/{id}`.

```json
{
  "id": "0193a1c0-1111-7aaa-8bbb-000000000001",
  "tipo": "exportacao_empresa",
  "status": "em_andamento",
  "progresso_percentual": 35,
  "criado_em": "2026-11-12T13:00:00Z",
  "concluido_em": null,
  "resultado": null,
  "erro": null
}
```

Estados: `na_fila`, `em_andamento`, `concluido`, `falhou`, `cancelado`. Ao concluir, `resultado` traz links de download com validade. Se o cliente assinou webhooks, recebe o evento correspondente (`exportacao.concluida`, `importacao.concluida`).

### 2.15 Idioma

- A API, os nomes de recurso, os campos e os enums são em **português do Brasil, sem acento**. Não há aliases em inglês. A exceção é SCIM, que segue o esquema padrão da RFC 7643.
- Textos para pessoas (`title`, `detail` dos erros, `mensagem` de validação, rótulos de indicadores) saem em pt-BR. O servidor respeita `Accept-Language`, mas no lançamento só existe `pt-BR`; outro idioma cai em `pt-BR` e a resposta traz `Content-Language: pt-BR`.
- **Códigos de erro e enums são estáveis e nunca traduzidos.** O cliente decide pelo código, nunca pelo texto.

### 2.16 Cabeçalhos padrão

| Cabeçalho | Direção | Uso |
|---|---|---|
| `Authorization: Bearer <token>` | Requisição | App e tokens de API |
| `Idempotency-Key` | Requisição | `POST` (2.9) |
| `If-Match` / `If-None-Match` | Requisição | Concorrência e cache (2.11) |
| `X-CSRF-Token` | Requisição | Web, em métodos que alteram dados (4.2) |
| `X-Request-Id` | Ambos | O cliente pode enviar; o servidor sempre devolve. Aparece nos erros como `id_rastreio`. |
| `traceparent` | Ambos | W3C Trace Context, para correlação com OpenTelemetry |
| `X-Tenaz-App-Versao` | Requisição | Versão do app de campo (18.5) |
| `X-Tenaz-Dispositivo-Id` | Requisição | Dispositivo do app (seção 12) |
| `ETag`, `Location` | Resposta | Versão do recurso, recurso criado |
| `Idempotency-Replayed: true` | Resposta | Resposta repetida pela chave de idempotência (2.9) |
| `Cache-Control: no-store` | Resposta | Toda resposta que traz credencial, `csrf_token` ou segredo de MFA |
| `Set-Cookie` | Resposta | Cookie de sessão web (4.2) |
| `Tenaz-Versao` | Resposta | Versão do produto em SemVer (`Tenaz-Versao: 1.4.2`), em toda resposta, inclusive de erro (18.8) |
| `RateLimit-*`, `Retry-After` | Resposta | Seção 5 |
| `Deprecation`, `Sunset`, `Link` | Resposta | Seção 18 |

---

## 3. Erros

### 3.1 Formato RFC 9457 (problem+json)

Todo erro usa `Content-Type: application/problem+json` e o formato da RFC 9457, com membros de extensão em português.

```json
{
  "type": "https://api.tenaz.com.br/erros/validacao_falhou",
  "title": "Os dados enviados não passaram na validação",
  "status": 422,
  "detail": "2 campos com problema.",
  "instance": "/v1/ordens-servico",
  "codigo": "validacao_falhou",
  "id_rastreio": "01JCZ4W8N0Q3R6S9T2V5X8Y1Z4",
  "erros": [
    { "campo": "ativo_id", "codigo": "obrigatorio", "mensagem": "Informe o ativo." },
    { "campo": "tipo", "codigo": "valor_invalido", "mensagem": "Tipo de OS inválido.", "valores_aceitos": ["preventiva", "preditiva", "corretiva_planejada", "corretiva_emergencial", "inspecao_legal", "melhoria", "implantacao"] }
  ]
}
```

| Membro | Descrição |
|---|---|
| `type` | URI estável do tipo de erro, que leva à página de documentação do código |
| `title` | Resumo legível, fixo por código |
| `status` | Código HTTP |
| `detail` | Explicação específica desta ocorrência |
| `instance` | Caminho da requisição |
| `codigo` | **Código estável** em snake_case, português. É o que o cliente deve tratar. |
| `id_rastreio` | Igual ao `X-Request-Id`; o suporte usa para localizar o log |
| `erros` | Lista de erros por campo, com `campo` em notação de ponto (`registro_nr12.indicacao_conclusiva`, `respostas[3].valor_numero`), `codigo` do vocabulário curto (abaixo), `mensagem` e, em enumeração, `valores_aceitos` |

Membros adicionais aparecem quando úteis: `recurso_atual` e `versao_atual` (412), `desafio_id` e `desafio_expira_em` (`401 mfa_obrigatorio`), `tentativas_restantes` (`401 mfa_codigo_invalido`), `mfa_pendente` (`401 nao_autenticado` de sessão restrita à ativação do MFA, 4.2), `error` (erros de `POST /oauth/token`, com o valor da RFC 6749, como `invalid_grant`), `estado_atual` e `transicoes_permitidas` (409 de transição), `plano_necessario` (403 de plano), `campos_faltantes` (422 de transição e de campo extra obrigatório), `estado_conexao`, `lado_dono`, `campos` e `grupos` (erros de conexão entre contas, 22.8), `tipos_documento_faltantes`, `ordens_servico_pendentes`, `predecessoras_pendentes`, `etapas_pendentes` e `ativos` (erros de entrada em operação e de projeto, seções 23 e 24), `faixas` (contrato com cliente, 25.2), `limite` (erros de limite de configuração).

### 3.2 Códigos de erro estáveis

| Código | HTTP | Quando |
|---|---|---|
| `requisicao_invalida` | 400 | JSON malformado, tipo errado, parâmetro desconhecido |
| `filtro_nao_suportado` | 400 | Filtro, ordenação ou expansão não declarada |
| `cursor_invalido` | 400 | Cursor corrompido, de outra consulta ou vencido |
| `campo_nao_permitido` | 400 | Campo proibido no corpo, como `empresa_id` |
| `data_sem_fuso` | 422 | Data-hora sem `Z` ou deslocamento |
| `nao_autenticado` | 401 | Sem credencial, ou credencial de empresa que exige segundo fator para quem ainda não o ativou; neste caso traz `mfa_pendente: true` (RN-ADM-027, 4.2) |
| `credencial_invalida` | 401 | Token ou sessão inválidos, expirados ou revogados |
| `mfa_obrigatorio` | 401 | Login exige segundo fator; traz `desafio_id` e `desafio_expira_em` (4.2) |
| `mfa_codigo_invalido` | 401 | Código TOTP errado ao concluir o login, ativar ou desativar o segundo fator; no login traz `tentativas_restantes` (cada erro consome uma das 5 tentativas do desafio) |
| `csrf_invalido` | 403 | Token CSRF ausente ou errado (web) |
| `permissao_negada` | 403 | O perfil não permite a operação |
| `escopo_insuficiente` | 403 | O token de API não tem o escopo exigido |
| `plano_nao_inclui_recurso` | 403 | Recurso fora do plano contratado; traz `plano_necessario` |
| `empresa_suspensa` | 403 | Empresa com contrato suspenso: só leitura e exportação |
| `recurso_nao_encontrado` | 404 | Não existe, **pertence a outra empresa** ou está **fora do alcance de prestador ou locatário** (nunca 403 nesses casos; ver 4.9) |
| `conflito` | 409 | Conflito genérico de unicidade |
| `tag_duplicada` | 409 | TAG de ativo já usada na empresa |
| `placa_duplicada` | 409 | Placa já cadastrada em veículo ativo |
| `transicao_invalida` | 409 | Ação não permitida no estado atual; traz `estado_atual` e `transicoes_permitidas` |
| `os_encerrada_imutavel` | 409 | Alteração em OS encerrada |
| `idempotencia_em_processamento` | 409 | Mesma `Idempotency-Key` ainda em processamento |
| `ultimo_administrador` | 409 | Desativar, ou tirar o perfil `administrador`, do último administrador ativo da empresa (RF-ADM-004) |
| `cursor_expirado` | 410 | Cursor de sincronização fora da retenção; refazer instantâneo |
| `versao_descontinuada` | 410 | Endpoint após a data de `Sunset` (seção 18) |
| `link_email_invalido` | 410 | Link de convite ou de redefinição de senha inexistente, vencido, já usado ou substituído por outro envio |
| `versao_desatualizada` | 412 | `If-Match` não confere; traz `recurso_atual` |
| `payload_muito_grande` | 413 | Corpo ou arquivo acima do limite |
| `tipo_midia_nao_suportado` | 415 | `Content-Type` não aceito |
| `pre_condicao_obrigatoria` | 428 | `PATCH`, `DELETE` ou ação de estado sem `If-Match` (2.11) |
| `validacao_falhou` | 422 | Erros de campo; detalhes em `erros` |
| `campo_somente_leitura` | 422 | Tentativa de alterar `status`, `numero`, `versao` ou campo calculado |
| `campos_obrigatorios_transicao` | 422 | Faltam campos exigidos pela transição; traz `campos_faltantes` |
| `motivo_obrigatorio` | 422 | Pausa, cancelamento, recusa de solicitação, recusa de liberação ou alteração do responsável pelo custo sem motivo |
| `idempotencia_conflito` | 422 | Mesma `Idempotency-Key` com corpo diferente |
| `senha_fraca` | 422 | Senha fora da política: menos de 10 caracteres, senha comum ou vazada, ou igual ao e-mail (RF-ADM-006) |
| `nr12_registro_incompleto` | 422 | Encerramento de OS em máquina NR-12 sem os sete campos do item 12.11.2 |
| `nr12_maquina_inapta` | 422 | Tentativa de devolver à operação, por alteração manual, máquina com indicação conclusiva `inapto` vigente |
| `leitura_regressiva` | 422 | Leitura menor que a anterior sem troca de medidor registrada |
| `leitura_salto_excessivo` | 422 | Leitura acima do teto de salto sem confirmação |
| `leitura_futura` | 422 | Data da leitura no futuro |
| `hierarquia_profundidade_excedida` | 422 | Ativo além do 6º nível |
| `hierarquia_ciclica` | 422 | Ativo ou localização definida como pai de si mesma |
| `estoque_insuficiente` | 422 | Saída ou reserva acima do saldo disponível |
| `estoque_em_transito` | 422 | Consumo, reserva ou movimentação de saldo ou de peça serializada ainda `em_transito`, antes da confirmação de recebimento pela subcontratada (20.2); traz `quantidade_em_transito` |
| `capacitacao_vencida` | 422 | Atribuição de OS regulada a técnico com capacitação vencida (Fase 2) |
| `ativo_inativo` | 422 | Operação em ativo inativo |
| `checklist_modelo_desatualizado` | 422 | Execução enviada para versão de modelo retirada |
| `dependencia_rejeitada` | 422 | Comando de sincronização que depende de outro comando rejeitado no mesmo lote (seção 12) |
| `importacao_arquivo_invalido` | 422 | Arquivo de importação ilegível ou fora do modelo |
| `nfe_destinatario_divergente` | 422 | CNPJ destinatário da NF-e não é da empresa |
| `nfe_ja_importada` | 409 | Chave de acesso já importada |
| `nfe_chave_invalida` | 422 | Chave de NF-e de remessa ou retorno sem 44 dígitos ou com dígito verificador inválido |
| `fornecedor_papel_incompativel` | 422 | O fornecedor não tem o papel exigido pela operação (por exemplo, OS atribuída a fornecedor sem `subcontratada`, `assistencia_tecnica` ou `oficina`; contrato `tomado` com fornecedor sem `locador`; almoxarifado `em_terceiro` de fornecedor sem `subcontratada`) |
| `fornecedor_inativo` | 422 | Operação nova com fornecedor inativo |
| `os_sem_fornecedor` | 422 | Ação que exige OS com `executor_tipo = fornecedor` (gerar link de uso único, registrar orçamento de fornecedor) |
| `link_execucao_expirado` | 410 | Link de uso único fora da validade |
| `link_execucao_utilizado` | 410 | Link de uso único já usado para registrar a execução |
| `link_execucao_revogado` | 410 | Link de uso único revogado pela empresa ou invalidado porque a OS mudou de executor, foi cancelada ou encerrada |
| `cliente_servico_inativo` | 422 | Operação nova (OS, equipamento, conexão) com cliente atendido inativo (SRV, 22.3) |
| `conexao_inativa` | 409 | Operação que exige conexão entre contas `ativa` feita com conexão `convidada`, `suspensa` ou `encerrada` (atribuir OS ao fornecedor conectado, abrir OS sincronizada em equipamento de cliente conectado, reprocessar vínculo); traz `estado_conexao` (22.8) |
| `conexao_duplicada` | 409 | Já existe conexão não encerrada entre as duas contas para o mesmo fornecedor (22.5) |
| `convite_conexao_invalido` | 422 | Aceite de convite com código errado ou vencido, ou feito por conta cuja raiz de CNPJ difere da do fornecedor convidado; traz `motivo` (`codigo_incorreto`, `expirado`, `cnpj_divergente`) |
| `campo_nao_compartilhado` | 422 | A requisição pede que um dado vá para a conta conectada, mas o grupo de campos não está em `campos_compartilhados` da conexão; traz `grupos` e `campos` (22.8) |
| `os_vinculada_somente_leitura` | 409 | Alteração ou ação, numa OS com vínculo entre contas, sobre campo ou etapa que pertence ao outro lado (por exemplo, a subcontratada tenta encerrar, ou o cliente tenta corrigir um apontamento da subcontratada); traz `lado_dono` e `campos` (22.8) |
| `ativo_sincronizado_somente_leitura` | 409 | Alteração, na conta da subcontratada, de campo de equipamento de cliente que vem da conta do cliente pela conexão; traz `campos` (22.4) |
| `sublocacao_nao_permitida` | 422 | Ativo com `propriedade = de_terceiro` incluído em contrato de locação `cedido` (21.1) |
| `almoxarifado_tipo_incompativel` | 422 | Tipo de movimentação incompatível com o tipo de almoxarifado de origem ou destino (por exemplo, `remessa_terceiro` para almoxarifado que não é `em_terceiro`) |
| `numero_serie_obrigatorio` | 422 | Movimentação ou consumo de item com `controla_serie = true` sem as peças serializadas |
| `numero_serie_duplicado` | 409 | Número de série já cadastrado para o mesmo item |
| `peca_situacao_invalida` | 409 | Peça serializada fora da situação exigida pela operação; traz `situacao_atual` |
| `item_nao_reparavel` | 422 | Envio para conserto de item com `reparavel = false` |
| `ativo_ja_locado` | 409 | Ativo já está em outro contrato de locação vigente no mesmo período |
| `contrato_locacao_sem_itens` | 422 | Ativação de contrato de locação sem ativos |
| `vistoria_entrega_pendente` | 422 | Vistoria de devolução sem vistoria de entrega do mesmo item |
| `vistoria_devolucao_pendente` | 422 | Encerramento de contrato com item sem vistoria de devolução; traz `itens_pendentes` |
| `registro_imutavel` | 409 | Alteração de registro em estado final que não aceita edição: projeto `concluido` ou `cancelado`, contrato com cliente `encerrado` ou `cancelado`, medição `fechada` (de contrato ou de etapas); traz `estado_atual`. A OS mantém o seu código próprio, `os_encerrada_imutavel` |
| `regra_combinacao_invalida` | 422 | Condição e efeito da regra configurável que não combinam, ou parâmetro exigido pela combinação ausente (por exemplo, `bloquear_ativo_ate_documento` sem o tipo de documento que levanta o bloqueio; `abrir_solicitacao` com `documento_a_vencer`) (23.1) |
| `limite_regras_ativas_excedido` | 422 | Ativação de regra além do limite de regras ativas por empresa (proposta: 50); traz `limite` |
| `limite_campos_extras_excedido` | 422 | Campo extra ativo além do limite por combinação de tipo de OS e classe de ativo (proposta: 30); traz `limite` |
| `campo_extra_obrigatorio` | 422 | Campo extra obrigatório vazio na abertura ou na conclusão da OS, ou campo exigido pela etapa de projeto vazio; traz `campos_faltantes` com as chaves (`campos_extras.rompeu_lacre`) (23.2) |
| `ativo_nao_em_implantacao` | 422 | `colocar-em-operacao` em ativo que já está `em_operacao`, ou OS `implantacao` em ativo `em_operacao` sem etapa de projeto (23.3) |
| `ativo_em_implantacao` | 422 | Operação que exige ativo `em_operacao` feita em ativo `em_implantacao`, como reativá-lo por `alterar-status-operacional`; use `colocar-em-operacao` (23.3) |
| `os_implantacao_pendente` | 422 | `colocar-em-operacao` com OS `implantacao` do ativo ainda não `encerrada` nem `cancelada`; traz `ordens_servico_pendentes` (números) (proposta) |
| `documentos_exigidos_pendentes` | 422 | `colocar-em-operacao` sem os documentos críticos válidos exigidos pela classe, ou conclusão de etapa de projeto sem os documentos exigidos; traz `tipos_documento_faltantes` |
| `ativos_nao_aptos_operacao` | 422 | Conclusão de projeto em que algum ativo `em_implantacao` não pode ser posto em operação; traz `ativos`, cada um com `tag` e o código do impedimento (24.6) |
| `dependencia_ciclica` | 422 | Etapas do modelo de projeto com predecessoras em ciclo, ou predecessora de outro modelo |
| `etapa_predecessora_pendente` | 409 | Ação em etapa de projeto `pendente` cujas predecessoras não estão `concluida` nem `dispensada`; traz `predecessoras_pendentes` |
| `etapa_os_pendente` | 422 | Conclusão de etapa com OS que não está `encerrada` nem `cancelada`, ou sem nenhuma OS `encerrada` quando a etapa gera OS; traz `ordens_servico_pendentes` |
| `etapa_entrega_nao_dispensavel` | 422 | Tentativa de dispensar a etapa de entrega do projeto (proposta) |
| `projeto_etapas_pendentes` | 422 | Conclusão de projeto com etapa que não está `concluida` nem `dispensada`; traz `etapas_pendentes` |
| `projeto_com_os_em_execucao` | 422 | Cancelamento de projeto com OS em `em_execucao`; pause-as antes; traz `ordens_servico_pendentes` |
| `ativo_em_outro_projeto` | 409 | Ativo incluído num projeto quando já está em outro projeto não concluído nem cancelado (proposta) |
| `etapa_ja_medida` | 409 | Etapa de projeto incluída numa medição de etapas quando já está em outra (proposta) |
| `contrato_cliente_sem_itens` | 422 | Ativação de contrato com cliente sem ativos cobertos |
| `faixas_desconto_invalidas` | 422 | Faixas de desconto do contrato com cliente sobrepostas, com lacuna ou fora da ordem da meta para baixo; traz `faixas` |
| `ativo_em_outro_contrato_cliente` | 409 | Ativo já coberto por outro contrato com cliente vigente no mesmo período |
| `medicao_periodo_duplicado` | 409 | Abertura de medição de contrato para período que já tem medição |
| `medicao_exclusoes_pendentes` | 422 | Fechamento de medição com parada de causa externa em `decisao_exclusao = pendente`; traz `ordens_servico_pendentes` |
| `causa_externa_nao_aceita` | 422 | Aprovação da exclusão de uma parada cuja causa externa não está entre as aceitas pelo contrato |
| `versao_app_nao_suportada` | 400 | App abaixo da versão mínima (18.5) |
| `limite_requisicoes_excedido` | 429 | Rate limit; ver `Retry-After` |
| `erro_interno` | 500 | Falha inesperada; o `id_rastreio` localiza o log |
| `servico_indisponivel` | 503 | Manutenção ou dependência fora; com `Retry-After` |

Regras:

- Um código publicado **nunca muda de significado** nem é removido em `/v1`. Códigos novos podem surgir; o cliente deve tratar código desconhecido pelo `status` HTTP.
- Erros 5xx nunca expõem detalhes internos (pilha, SQL, nome de tabela). Todo erro inesperado sai como `500 erro_interno`; o código interno e a causa ficam só no log, localizados pelo `id_rastreio`.
- Os códigos são os desta tabela, com a mesma grafia no código Go (`internal/plataforma/erros`), que é a fonte do contrato gerado. A URI de `type` é `https://api.tenaz.com.br/erros/<codigo>`.
- Erros que a biblioteca de rotas gera sozinha também saem com estes códigos: corpo que não é JSON, valor de tipo errado ou campo desconhecido, `400 requisicao_invalida`; `empresa_id` no corpo ou na consulta, `400 campo_nao_permitido`; campo somente leitura no corpo (`id`, `versao`, `situacao`, `status`, `numero`, `criado_em`, `atualizado_em`, `plano`, `email` do usuário e semelhantes), `422 campo_somente_leitura`; parâmetro de consulta não declarado, `400 filtro_nao_suportado`; rota ou método inexistente e identificador malformado no caminho, `404 recurso_nao_encontrado`.
- Erros por campo usam um vocabulário curto e estável em `erros[].codigo`: `obrigatorio`, `valor_invalido`, `tamanho_maximo`, `tamanho_minimo`, `formato_invalido`, `fora_da_faixa`, `nao_encontrado`, `duplicado`, `inativo`.

---

## 4. Autenticação e autorização

### 4.1 Resumo

| Método | Cliente | Plano | Fase | Credencial |
|---|---|---|---|---|
| E-mail e senha + MFA opcional | Web e app | Todos | MVP | Sessão (web) ou tokens OAuth (app) |
| SSO OIDC (Google, Microsoft Entra ID, outros OIDC) | Web e app | Profissional | MVP | Idem, após login no provedor |
| SSO SAML 2.0 | Web e app | Enterprise | Fase 2 | Idem |
| SCIM 2.0 | Provedor de identidade | Enterprise | Fase 2 | Token SCIM |
| Token de API | Integrações | Profissional | MVP | `Authorization: Bearer tnz_live_...` |

As senhas seguem as regras do documento 9 (hash Argon2id, proposta). A Tenaz não guarda senha de usuário que entra por SSO.

### 4.2 Web: sessão por cookie httpOnly + CSRF

- `POST /v1/auth/login` com `email`, `senha` e, opcional, `empresa_escolhida_id` (sem ele, a empresa da última sessão). Toda falha de credencial (e-mail inexistente, senha errada, conta bloqueada) responde igual, `401 credencial_invalida`, no mesmo tempo. Se o usuário tiver MFA, a resposta é `401 mfa_obrigatorio` com `desafio_id` e `desafio_expira_em`, e o login termina em `POST /v1/auth/mfa/verificar` com `desafio_id`, `codigo` (6 dígitos) e `empresa_escolhida_id` opcional (TOTP; WebAuthn em fase posterior, proposta).
- **`desafio_id` é segredo:** 256 bits aleatórios, de uso único, válido por 5 minutos e por 5 tentativas; o banco guarda só o sha256. Vai sempre no corpo, nunca na URL nem em log. Código errado responde `401 mfa_codigo_invalido` com `tentativas_restantes`; desafio vencido, gasto ou esgotado, `401 credencial_invalida`, e o login recomeça pela senha.
- O servidor cria a sessão e devolve o cookie `__Host-tenaz_sessao` com `HttpOnly`, `Secure`, `SameSite=Lax`, `Path=/`, sem `Domain`, e `Expires` no teto absoluto. O valor é um identificador aleatório de 256 bits; o estado da sessão fica no servidor.
- **Desenvolvimento sobre http:** o navegador recusa o prefixo `__Host-` sem `Secure`. Com `TENAZ_COOKIE_SEGURO=false`, aceito **só** com `TENAZ_AMBIENTE=desenvolvimento` (o processo recusa subir com ele em homologação e produção), o cookie se chama `tenaz_sessao` e sai sem `Secure`; os demais atributos não mudam. A tarefa `ferramentas/tarefas.sh api` sobe assim.
- A resposta do login, de `mfa/verificar` e de `GET /v1/auth/sessao` é a **sessão**: `canal` (`web` ou `app`), `csrf_token` (nulo no app), `expira_em` (só no login e na troca de empresa), `escopo_acesso`, `usuario` (`id`, `nome`, `email`, `pessoa_id`, `mfa_ativo`), `empresa` (`nome`, `plano`, `srv_habilitado`, sem identificador), `perfis`, `permissoes` (efetivas, em ordem alfabética) e `empresas` (`id`, `nome`, `atual`), com a atual primeiro, e `mfa_pendente` (verdadeiro na sessão restrita à ativação do MFA, abaixo). A troca de empresa devolve o mesmo objeto, com `tokens` (nulo na web). Respostas com credencial trazem `Cache-Control: no-store`.
- **Segundo fator exigido pela empresa (RN-ADM-027 em [4_FLUXOS_PROCESSO.md](4_FLUXOS_PROCESSO.md); sumário 3.6, "MFA da empresa").** `empresa.mfa_obrigatorio` sai em `GET /v1/empresa` e é ligado ou desligado pelo `PATCH /v1/empresa` (`adm.gerenciar_empresa`; a troca vai para a trilha, RN-ADM-030). Com ele ligado, quem entra nessa empresa sem MFA ativado (no login ou na troca de empresa pela web) recebe uma **sessão web restrita**, com `mfa_pendente = true`, que só atende `GET /v1/auth/sessao`, `POST /v1/auth/logout`, `POST /v1/auth/trocar-empresa`, `POST /v1/auth/mfa/iniciar` e `POST /v1/auth/mfa/confirmar`. Qualquer outra rota, inclusive `/v1/usuarios/eu`, `trocar-senha`, `mfa/desativar` e `POST /oauth/authorize`, responde `401 nao_autenticado` com `mfa_pendente: true`. Confirmado o MFA, a mesma sessão volta a valer, sem novo login. A restrição é conferida a cada requisição: se a empresa liga o interruptor com a sessão aberta, a sessão aberta fica restrita já na requisição seguinte.
- No app não existe sessão restrita: com a empresa exigindo e a pessoa sem MFA, o access token, a troca do código e a renovação em `POST /oauth/token` e a troca de empresa pelo app respondem `401 nao_autenticado` com `mfa_pendente: true`; na renovação, a sessão do aparelho é revogada, e o app entra de novo depois de a pessoa ativar o MFA pela web.
- Enquanto a pessoa tiver acesso a alguma empresa que exige MFA, `POST /v1/auth/mfa/desativar` responde `409 transicao_invalida` com `estado_atual = "exigido_pela_empresa"`. Código TOTP errado em `mfa/desativar` soma uma falha ao bloqueio progressivo da identidade, como no login, e um código de passo TOTP já usado vale como errado em qualquer rota (RN-ADM-027).
- Validade (proposta): 12 horas de inatividade, 7 dias no máximo. O administrador pode reduzir em `/v1/parametros-empresa`.
- **CSRF:** a resposta do login e `GET /v1/auth/sessao` trazem `csrf_token`, derivado da sessão por HMAC (o mesmo valor em toda consulta da mesma sessão; o banco guarda só o hash). Todo `POST`, `PATCH`, `PUT` e `DELETE` feito com cookie exige o cabeçalho `X-CSRF-Token` igual ao token da sessão; senão, `403 csrf_invalido`.
- **Origin:** o servidor confere `Origin` (na falta dele, a origem de `Referer`) contra a lista de origens da web (`TENAZ_ORIGENS_WEB`; vazia, a origem de `TENAZ_URL_WEB`). Requisição com cookie que altera dados sem nenhum dos dois cabeçalhos, ou de outra origem, responde `403 csrf_invalido`. Rota pública que altera dados (login, cadastro da conta, convite) também recusa `Origin` de outra origem, contra login CSRF; sem `Origin` (cliente que não é navegador), ela é aceita.
- **Escopo e permissão:** depois da credencial, a camada HTTP confere o escopo de acesso da rota (`x-escopos-acesso`; fora dele, `404`) e as permissões de `x-permissoes` (`403 permissao_negada` para interno, `404` para externo), antes de ler o corpo. A sessão é resolvida em dois pools, nunca na mesma transação (seção 6.3 do documento 5).
- **Rotas do próprio usuário** (todas as credenciais, inclusive `prestador` e `locatario`): `POST /v1/auth/trocar-senha` (`senha_atual`, `nova_senha`; encerra as outras sessões da pessoa e mantém a atual; senha atual errada conta como falha de login); `POST /v1/auth/mfa/iniciar` (devolve `segredo_base32` e `uri_otpauth`, uma vez), `POST /v1/auth/mfa/confirmar` e `POST /v1/auth/mfa/desativar` (`codigo`).
- **Senha esquecida e convite:** `POST /v1/auth/recuperar-senha` (`email`) responde `202` sempre, exista ou não o e-mail, e envia o link (1 hora). O link de e-mail leva o token no fragmento (`#token=...`); a página lê o fragmento e o envia no corpo de `POST /v1/auth/redefinir-senha` (`token`, `nova_senha`; encerra todas as sessões da pessoa) ou de `POST /v1/auth/convites/consultar` (`token`; devolve `email` e `pede_senha`) e `POST /v1/auth/convites/aceitar` (`token`, `senha` quando `pede_senha`). Link inválido: `410 link_email_invalido`; senha fora da política: `422 senha_fraca`.
- Requisições autenticadas por `Authorization: Bearer` ignoram o cookie e não exigem CSRF.
- `POST /v1/auth/logout` invalida a sessão no servidor (web ou app) e, na web, apaga o cookie (`Max-Age=0`). O administrador pode encerrar todas as sessões de um usuário (`POST /v1/usuarios/{id}/encerrar-sessoes`).

### 4.3 App de campo: OAuth 2.1 com PKCE e refresh token rotativo

O app é um cliente público OAuth 2.1, sem segredo embutido.

1. O app gera `code_verifier` e `code_challenge` (S256) e abre o navegador do sistema (não WebView) em `GET /oauth/authorize?response_type=code&client_id=tenaz-app&redirect_uri=br.com.tenaz.app:/oauth&code_challenge=...&code_challenge_method=S256&state=...`.
2. A API confere `client_id` e `redirect_uri` (byte a byte) e responde `302` para a tela de autorização da web, `/autorizar-app`, com os mesmos parâmetros. Cliente ou URI de retorno inválidos respondem `422 validacao_falhou` e nunca redirecionam. Na tela, o usuário entra com a sessão web normal (senha e MFA; SSO quando existir) e, se tiver acesso a mais de uma empresa, escolhe uma (lista de `GET /v1/auth/sessao`). A web então chama `POST /oauth/authorize` com a sessão web (cookie, `X-CSRF-Token` e `Origin`; o access token do app não serve), em JSON: `response_type`, `client_id`, `redirect_uri`, `code_challenge`, `code_challenge_method`, `state` e `empresa_escolhida_id`. A resposta é `{"uri_redirecionamento": "br.com.tenaz.app:/oauth?code=...&state=..."}`, com código de uso único válido por 60 segundos; a web navega até ela e o sistema entrega ao app. Perfil sem `mob.usar_app` (o `locatario`) recebe a URI com `error=access_denied`; empresa fora do acesso, `404`. Assim o app nunca vê a senha e o login do app segue as mesmas regras do login da web.
3. O app troca o código em `POST /oauth/token`, em `application/x-www-form-urlencoded` (RFC 6749): `grant_type=authorization_code`, `code`, `redirect_uri`, `client_id` e `code_verifier`, com o cabeçalho opcional `X-Tenaz-Dispositivo-Id`. Código, URI de retorno, cliente ou PKCE que não conferem gastam o código e respondem `401 credencial_invalida` com `error=invalid_grant`; código reapresentado depois de usado revoga a sessão que ele gerou. O app recebe:

```json
{
  "access_token": "eyJhbGciOiJFZERTQSIsImtpZCI6IjIwMjYtMTEifQ...",
  "token_type": "Bearer",
  "expires_in": 900,
  "refresh_token": "tnz_rt_6m1Qx...",
  "scope": "app_campo"
}
```

- **Access token:** JWT assinado (EdDSA, proposta), validade de **15 minutos**, com `sub` (usuário), `emp` (empresa), `sid` (sessão do dispositivo), `aud`, `exp`. As permissões não vão no token; são lidas a cada requisição para que revogações valham na hora.
- **Verificação do access token:** o servidor recusa o access token sem `exp`, com `exp` vencido além da tolerância de 30 s, com `iat` ou `nbf` no futuro além da mesma tolerância, com `iss` diferente de `tenaz` ou `aud` diferente de `tenaz-api`, ou com cabeçalho diferente de `alg: EdDSA`, `typ: JWT` e `kid` conhecido. A Tenaz não emite `nbf`; quando presente, é respeitado (RFC 7519, seção 4.1.5). Token recusado responde `401 credencial_invalida`.
- **Refresh token rotativo:** cada uso em `POST /oauth/token` (`grant_type=refresh_token`) devolve um refresh token novo e invalida o anterior. Se um refresh token já usado for reapresentado, **toda a família é revogada** e o usuário precisa entrar de novo (detecção de reuso).
- **Compatível com 7 dias offline:** o refresh token vale **30 dias** de inatividade (janela deslizante, proposta), acima da meta de 7 dias sem conexão do sumário. Os comandos feitos offline ficam na fila local e são enviados depois de renovar o token; access token expirado nunca causa perda de registro.
- O refresh token é guardado no armazenamento seguro do sistema (Keychain/Keystore).
- `POST /oauth/revoke` (formulário, `token`, RFC 7009) revoga a família de refresh tokens e a sessão do aparelho (logout no app) e responde `200`, também para token desconhecido. O refresh token no corpo é a credencial; a rota não exige outra. O administrador pode revogar um dispositivo em `DELETE /v1/dispositivos/{id}`.
- Os erros de `/oauth/token` saem em `problem+json` com `codigo`, como toda a API, e com o membro `error` da RFC 6749; o status de credencial inválida é `401`, e não o `400` da RFC (pendência A-37). Duas renovações simultâneas com o mesmo refresh token contam como reuso: o app serializa as renovações.

### 4.4 Tokens de API por empresa

- Criados pelo `administrador` em `POST /v1/tokens-api`, com `nome`, `escopos`, `expira_em` (opcional, proposta: máximo de 2 anos) e `ips_permitidos` (opcional, lista CIDR).
- O valor é exibido **uma única vez** na criação. A Tenaz guarda só o hash (SHA-256) e os 4 últimos caracteres.
- Formato (proposta): `tnz_live_<32 caracteres base62>` em produção e `tnz_test_...` no sandbox. O prefixo permite que scanners de segredo detectem vazamentos.
- O token pertence à **empresa**, não a um usuário. As ações aparecem na auditoria como "token: <nome>".
- `POST /v1/tokens-api/{id}/rotacionar` emite um token novo e mantém o anterior válido por até 24 horas, para troca sem parada.
- O token registra `ultimo_uso_em` e `ultimo_ip`.

**Escopos (proposta).** Cada endpoint declara no OpenAPI os escopos exigidos (`x-escopos`). `escrita` inclui `leitura`.

| Escopo | Cobre |
|---|---|
| `empresa:leitura` | Empresa, unidades, usuários (sem dados sensíveis), perfis |
| `usuarios:escrita` | Criar, alterar e desativar usuários |
| `pessoas:leitura` / `pessoas:escrita` | Pessoas |
| `ativos:leitura` / `ativos:escrita` | Localizações, classes, ativos, documentos técnicos |
| `medidores:leitura` / `medidores:escrita` | Medidores e leituras |
| `solicitacoes:leitura` / `solicitacoes:escrita` | Solicitações |
| `ordens_servico:leitura` / `ordens_servico:escrita` | OS, tarefas, apontamentos, peças, códigos de falha |
| `planos:leitura` / `planos:escrita` | Planos de manutenção, gatilhos, procedimentos |
| `checklists:leitura` / `checklists:escrita` | Modelos e execuções |
| `estoque:leitura` / `estoque:escrita` | Itens, almoxarifados, saldos, movimentações (inclusive remessa e retorno a terceiros), peças serializadas, ordens de reparo, inventários, compras, NF-e |
| `terceiros:leitura` / `terceiros:escrita` | Fornecedores, serviços de fornecedor, contratos de prestação, avaliações, orçamentos de fornecedor, links de execução |
| `locacao:leitura` / `locacao:escrita` | Locatários, contratos de locação, itens, matriz de responsabilidade, vistorias, relatório de cobrança |
| `servicos:leitura` / `servicos:escrita` | SRV: clientes atendidos, especialidades de técnico, serviços a realizar da OS (`servico_os`) e relatório da OS em PDF. Equipamentos de clientes usam `ativos:*`; OS, apontamentos e componentes, `ordens_servico:*` |
| `conexoes:leitura` | Conexões entre contas, ativos compartilhados e vínculos de OS, só leitura (22.10) |
| `configuracao:leitura` / `configuracao:escrita` | Regras configuráveis e seus disparos, definições de campos extras da OS, catálogo de pacotes de segmento e aplicações de pacote (só leitura). Os valores dos campos extras numa OS usam `ordens_servico:*`; a entrada em operação do ativo usa `ativos:escrita` |
| `projetos:leitura` / `projetos:escrita` | PRJ: modelos de projeto, projetos, etapas, lançamentos de custo, termo de entrega e medições de etapas. As OS de implantação usam `ordens_servico:*` |
| `contratos_cliente:leitura` / `contratos_cliente:escrita` | CTR: contratos com clientes, itens cobertos, painel e medições do contrato. O cadastro de clientes atendidos usa `servicos:*` |
| `frota:leitura` / `frota:escrita` | Dados de veículo, pneus, abastecimentos, multas, sinistros |
| `conformidade:leitura` / `conformidade:escrita` | Documentos com vencimento, capacitações, NR-12, NR-13, bloqueios |
| `despacho:consulta` | Só a consulta de bloqueio de despacho (seção 16), para TMS |
| `indicadores:leitura` | KPIs, painéis e relatórios |
| `telematica:ingestao` | Envio de eventos de telemática (14.1) |
| `importacoes:escrita` | Importações CSV, cartão combustível e NF-e |
| `webhooks:gerenciar` | Criar e manter webhooks |
| `auditoria:leitura` | Eventos de auditoria |
| `exportacao:empresa` | Exportação completa (seção 17) |

Ações administrativas sensíveis (criar tokens, alterar SSO, alterar perfis, alterar parâmetros de cobrança, convidar, aceitar, recusar, suspender, reativar ou encerrar conexão entre contas e alterar `campos_compartilhados`, aplicar pacote de segmento, ativar regra configurável, reabrir medição de contrato) **não** estão disponíveis para tokens de API (proposta). A Tenaz não guarda dado de saúde (sumário 3.6): nenhuma rota, campo, anexo nem escopo o recebe ou o devolve.

### 4.5 SSO OIDC (Profissional)

- Configurado em `/v1/provedores-identidade` pelo `administrador`: `tipo=oidc`, `emissor`, `client_id`, `client_secret` (guardado cifrado), domínios de e-mail associados.
- Fluxo Authorization Code com PKCE também no web. O login por e-mail descobre o provedor pelo domínio.
- **Provisionamento no primeiro login (JIT):** cria o usuário com o perfil padrão definido no provedor. Mapeamento de grupos do provedor para perfis Tenaz é opcional (`mapeamento_grupos`).
- O administrador pode exigir SSO para o domínio (`sso_obrigatorio=true`), bloqueando a senha local, com exceção de uma conta de emergência (proposta).

### 4.6 SAML 2.0 e SCIM 2.0 (Enterprise)

- **SAML:** a Tenaz é SP. Metadados em `GET /saml/{provedor_id}/metadata`, ACS em `POST /saml/{provedor_id}/acs`. Asserções assinadas são obrigatórias; cifradas, opcionais. Atributos mapeados: e-mail, nome, grupos.
- **SCIM 2.0 (RFC 7643 e 7644):** base `https://api.tenaz.com.br/scim/v2`, autenticada por token SCIM próprio (`tnz_scim_...`), criado em `POST /v1/provedores-identidade/{id}/token-scim`. Recursos `Users` e `Groups`, com `filter` (`userName eq`), `PATCH` e paginação SCIM. Grupos SCIM mapeiam para perfis. Desativar um usuário no provedor (`active=false`) revoga sessões e refresh tokens na hora.
- SCIM segue os nomes em inglês do padrão. É a única exceção à regra de idioma (2.15).

### 4.7 Como a empresa (tenant) é determinada

Regra central de isolamento: **a empresa vem sempre da credencial, nunca de dado enviado pelo cliente.**

| Credencial | Origem da empresa |
|---|---|
| Sessão web | Registro da sessão no servidor (`sessao.empresa_id`) |
| Access token do app | Claim `emp`, conferida contra a sessão do dispositivo |
| Token de API | Registro do token (`token_api.empresa_id`) |
| Token SCIM | Registro do provedor de identidade |
| Webhook de provedor de telemática | Segredo da conexão (`conexao_telematica.empresa_id`), conferido pela assinatura |
| Link de uso único: ativação | Registro do link (`link_execucao.empresa_id`), localizado pelo hash SHA-256 do token recebido **no corpo** de `POST /v1/publico/links-execucao/ativar` |
| Link de uso único: sessão curta | Registro da sessão do link no servidor, que guarda `link_execucao_id` e `empresa_id` (19.4) |
| Processo de sincronização entre contas (SRV) | Não é credencial de cliente da API. Um worker interno, com papel próprio no banco, lê a `conexao_empresa` (esquema da plataforma, fora do RLS por empresa) e grava em cada conta numa transação separada, com `app.empresa_id` daquela conta (22.6) |

Consequências:

- Não existe parâmetro, cabeçalho ou subdomínio de empresa. `empresa_id` enviado no corpo ou na consulta devolve `400 campo_nao_permitido`. A escolha da empresa ao entrar ou ao trocar usa `empresa_escolhida_id` (2.2), conferido contra os acessos da pessoa e só emitindo credencial nova.
- A camada de acesso a dados abre uma transação por requisição e executa `set_config('app.empresa_id', <id>, true)` (escopo local da transação) antes de qualquer consulta. O PostgreSQL aplica Row-Level Security forçado em todas as tabelas de negócio, e o repositório também filtra por `empresa_id` (dupla proteção, conforme o risco R2 do sumário).
- Um `id` de outra empresa se comporta como inexistente: `404 recurso_nao_encontrado`, nunca `403`, para não revelar que o recurso existe.
- Usuário com acesso a mais de uma empresa (consultor, grupo econômico) troca de empresa em `POST /v1/auth/trocar-empresa`, que emite **nova** sessão ou novos tokens. Uma credencial nunca vale para duas empresas. Isso vale também para contas conectadas (SRV): o gestor da subcontratada não entra na conta do cliente, e o do cliente não entra na da subcontratada; cada lado vê só a própria cópia da OS e do equipamento.
- Testes automatizados de vazamento entre empresas rodam em todo build sobre todos os endpoints do contrato (meta de isolamento do sumário).

### 4.8 Autorização: perfis, unidades e plano

Cada requisição passa por quatro verificações, nesta ordem:

1. **Credencial válida** → senão `401`.
2. **Plano**: o recurso está no plano da empresa → senão `403 plano_nao_inclui_recurso`.
3. **Perfil ou escopo**: o perfil do usuário (ou o escopo do token) permite a operação → senão `403 permissao_negada` ou `403 escopo_insuficiente`.
4. **Alcance**: o registro está numa unidade a que o usuário tem acesso, e regras de propriedade (o técnico só aponta horas em OS em que está atribuído; o solicitante vê só as próprias solicitações, proposta; prestador e locatário só veem o que a seção 4.9 permite) → senão `404` para leitura e `403 permissao_negada` para escrita em registro visível.

Perfis padrão do sumário (seção 3.1): `administrador`, `gestor`, `planejador`, `tecnico`, `almoxarife`, `comprador`, `plh`, `motorista`, `solicitante`, `visualizador`, `prestador`, `locatario`. Grupos usados no catálogo (seção 6):

| Grupo no catálogo | Perfis |
|---|---|
| **Todos** | Qualquer usuário interno autenticado, inclusive `solicitante` e `motorista`. **Não** inclui `prestador` nem `locatario` (ver abaixo) |
| **Leitura** | `administrador`, `gestor`, `planejador`, `tecnico`, `almoxarife`, `comprador`, `plh`, `visualizador` |
| **Gestão** | `administrador`, `gestor`, `planejador` |
| **Execução** | Gestão + `tecnico` |
| **Estoque** | `administrador`, `gestor`, `almoxarife` |
| **Compras** | Estoque + `comprador` |
| **Liberação** | Quem tem a permissão `os.encerrar` (padrão em `gestor` e `planejador`; o `administrador` pode concedê-la a outros perfis internos, nunca a `prestador` nem a `locatario`) |
| **Destino de peça** | Quem tem a permissão `est.decidir_destino` (padrão em `gestor`, `planejador` e `almoxarife`): decide o destino de peça com defeito na quarentena (conserto, sucata, devolução) |
| **Regras** | Quem tem a permissão `adm.gerenciar_regras` (padrão em `administrador` e `gestor`): cria, altera, ativa, inativa e exclui regras configuráveis |
| **Configuração** | Quem tem a permissão `adm.gerenciar_configuracao` (padrão em `administrador` e `gestor`): mantém as definições de campos extras da OS e aplica pacotes de segmento |
| **Entrada em operação** | Quem tem a permissão `atv.colocar_em_operacao` (padrão em `administrador`, `gestor` e `planejador`): coloca em operação ativo `em_implantacao` |
| **Projetos** | Quem tem a permissão `prj.gerenciar` (padrão em `administrador`, `gestor` e `planejador`): mantém modelos e projetos e executa as ações do projeto e das etapas |
| **Dispensa de etapa** | Quem tem a permissão `prj.dispensar_etapa` (padrão em `administrador` e `gestor`): dispensa etapa de projeto |
| **Ajuste de medição** | Quem tem a permissão `prj.ajustar_medicao` (padrão em `administrador` e `gestor`): ajusta o valor de uma linha da medição de etapas |
| **Contratos** | Quem tem a permissão `ctr.gerenciar` (padrão em `administrador` e `gestor`): mantém contratos com clientes e itens e executa `ativar`, `encerrar` e `cancelar` |
| **Fechamento de medição** | Quem tem a permissão `ctr.fechar_medicao` (padrão em `administrador` e `gestor`): decide as exclusões de paradas e fecha a medição do contrato |
| **Reabertura de medição** | Quem tem a permissão `ctr.reabrir_medicao` (padrão em `administrador` e `gestor`): reabre medição de contrato `fechada` |
| **Conexão** | Quem tem a permissão `srv.gerenciar_conexao` (padrão em `administrador` e `gestor`), nos dois lados da conexão entre contas. O servidor confere o lado: só o cliente convida e altera ativos e campos compartilhados, só a subcontratada aceita ou recusa, e só quem suspendeu reativa |
| **Admin** | `administrador` |
| **Externos** | `prestador` e `locatario`, sempre com o alcance da seção 4.9 |
| **Integração** | Somente token de API com o escopo indicado |

Quando um perfil específico aparece no catálogo (por exemplo `motorista`, `plh`, `prestador`), ele se soma ao grupo. Os perfis externos (`prestador`, `locatario`) só acessam as rotas em que aparecem nomeados e as rotas de sessão e infraestrutura do próprio usuário (`/v1/auth/*`, `/oauth/*`, `GET /v1/usuarios/eu`, `PATCH /v1/usuarios/eu`, anexos, processamentos, notificações, dispositivos e, para `prestador`, sincronização). Qualquer outra rota devolve `404 recurso_nao_encontrado` para eles, não `403`, para não revelar a estrutura de dados da empresa contratante.

O catálogo de permissões e a matriz dos doze perfis padrão são os de RF-ADM-005, em [1_ESPECIFICACAO_REQUISITOS.md](1_ESPECIFICACAO_REQUISITOS.md), canônicos por decisão do sumário 3.6; os grupos acima são a forma resumida desse catálogo no catálogo de endpoints, e a permissão exigida por rota sai no contrato em `x-permissoes`. Perfis personalizados (proposta: plano Profissional) combinam permissões finas no formato `<módulo>.<ação>` (sumário 3.6), com o código do módulo em minúsculas, listadas em `GET /v1/permissoes`. Doze permissões finas aparecem neste documento. Toda permissão segue o formato `<módulo>.<ação>`, sem exceção (sumário 3.6):

| Permissão | Dá direito a | Perfis padrão |
|---|---|---|
| `os.encerrar` | Liberar (`encerrar`) a OS e recusar a liberação (`recusar-liberacao`). Substitui os nomes `os.liberar` e `ordem_servico.liberar` das versões anteriores | `gestor`, `planejador`; nunca `prestador` nem `locatario` |
| `est.decidir_destino` | Decidir o destino de peça com defeito: abrir ordem de reparo, baixar como sucata, declarar irreparável, devolver ao fornecedor | `gestor`, `planejador`, `almoxarife` |
| `adm.gerenciar_regras` | Criar, alterar, ativar, inativar e excluir regras configuráveis (`/v1/regras-automacao`, RF-ADM-016). Ler regras e disparos fica com o grupo Gestão | `administrador`, `gestor` |
| `adm.gerenciar_configuracao` | Criar, alterar, ordenar e inativar definições de campos extras da OS (`/v1/campos-extras-os`, RF-OS-021) e pré-visualizar e aplicar pacotes de segmento (RF-ADM-017) | `administrador`, `gestor` |
| `atv.colocar_em_operacao` | Executar `colocar-em-operacao`, unitário ou em lote, e concluir projeto com ativos `em_implantacao` (RF-ATV-013, RF-PRJ-007) | `administrador`, `gestor`, `planejador` |
| `prj.gerenciar` | Criar e manter modelos de projeto e projetos; `iniciar`, `suspender`, `retomar`, `concluir` e `cancelar` o projeto; manter etapas, gerar OS, iniciar e concluir etapas, gerar e fechar medição de etapas (as três últimas por extensão, proposta) | `administrador`, `gestor`, `planejador` |
| `prj.dispensar_etapa` | Dispensar etapa de projeto `pendente` ou `disponivel`, com motivo (RF-PRJ-003) | `administrador`, `gestor` |
| `prj.ajustar_medicao` | Ajustar, com motivo, o valor de uma linha da medição de etapas (RF-PRJ-006) | `administrador`, `gestor` |
| `ctr.gerenciar` | Criar e manter contratos com clientes e itens; `ativar`, `encerrar` e `cancelar`; abrir medição (por extensão, proposta) (RF-CTR-001, RF-CTR-002) | `administrador`, `gestor` |
| `ctr.fechar_medicao` | Decidir a exclusão de paradas com causa externa (`decisao_exclusao`) e fechar a medição do contrato (RF-CTR-004) | `administrador`, `gestor` |
| `ctr.reabrir_medicao` | Reabrir medição de contrato `fechada`, com motivo (RF-CTR-004) | `administrador`, `gestor` |
| `srv.gerenciar_conexao` | Convidar, aceitar, recusar, suspender, reativar e encerrar conexão entre contas; alterar ativos e campos compartilhados (grupo Conexão, 6.10; RF-SRV-014, RF-SRV-019) | `administrador`, `gestor`; nunca `prestador` nem `locatario` |

As nove permissões acrescentadas na versão 0.4 (`adm.*`, `atv.*`, `prj.*` e `ctr.*` acima, sumário 3.6) nunca são concedidas a `prestador` nem a `locatario` (proposta).

Na conta de uma subcontratada (SRV), `os.encerrar` vale só para OS sem vínculo com conta de cliente; numa OS vinculada, a liberação é sempre do cliente (22.8).

### 4.9 Perfis externos: prestador e locatário

O `prestador` é um usuário ligado a um fornecedor (`usuario.fornecedor_id`); o `locatario` é um usuário ligado a um locatário (`usuario.locatario_id`). Os dois são criados pela empresa contratante em `POST /v1/usuarios`, que exige o vínculo correspondente ao perfil (`422 validacao_falhou` sem ele). Um usuário externo não pode ter perfil interno ao mesmo tempo (proposta).

O filtro de alcance desses perfis é **adicional** ao isolamento por empresa (RLS e filtro do repositório, 4.7): o repositório acrescenta a condição de fornecedor ou de locatário a toda consulta feita com essas credenciais. Uma suíte de testes automatizados percorre todos os endpoints do contrato com usuários `prestador` e `locatario` de dois fornecedores e dois locatários da mesma empresa e de empresas diferentes, e quebra o build se algum registro fora do alcance aparecer (risco R4 do sumário).

**Escopo de acesso.** Toda credencial tem um escopo de acesso, com os valores do sumário 3.6: `interno` (usuários dos perfis internos e tokens de API), `prestador` (usuário com perfil `prestador`), `locatario` (usuário com perfil `locatario`) e `link` (sessão curta de um link de uso único, 19.4). O valor sai em `escopo_acesso` de `GET /v1/auth/sessao` e de `GET /v1/usuarios/eu` e é gravado em todo `evento_auditoria`. O repositório escolhe o filtro de alcance pelo escopo. O nome `fornecedor` não é mais usado como valor de escopo. Técnicos e gestores da conta de uma subcontratada (SRV) são `interno` **na conta dela**; na conta do cliente eles não têm credencial nenhuma.

**Regra geral:** registro fora do alcance devolve `404 recurso_nao_encontrado` tanto em leitura quanto em escrita. Registro visível com ação não permitida devolve `403 permissao_negada`. Listas devolvem só o que está no alcance, sem erro.

| Recurso | `prestador` vê e faz | `locatario` vê e faz | Fora do alcance (404) |
|---|---|---|---|
| `ordens-servico` | Lista e detalha OS com `executor_tipo = fornecedor` e `fornecedor_id` igual ao do usuário, exceto as OS com vínculo para a conta própria do fornecedor (essas são executadas na conta dele, 22.8). Executa `reconhecer`, `assumir` (OS do seu fornecedor sem responsável), `iniciar`, `pausar`, `retomar`, `concluir`; nunca `encerrar` nem `recusar-liberacao`; aponta horas, comenta, anexa fotos, registra assinatura, consome peças do almoxarifado `em_terceiro` do seu fornecedor; marca `causa_externa`. Vê e preenche só os campos extras com `preenchido_pelo_prestador = true` (23.2). Não vê `custos` da empresa, só o valor do próprio orçamento, nem `etapa_projeto_id` (vem `null`) | Lista e detalha OS dos ativos locados a ele em contrato `vigente` (ou encerrado há até 90 dias, proposta): número, tipo, status, datas, descrição, `responsavel_custo`. Comenta. Não vê apontamentos, peças, custos internos nem dados de técnicos | OS de outro fornecedor, OS retirada do fornecedor, OS de ativo não locado a ele |
| `ordens-servico/{id}/custos` | 404 | Só o valor das OS com `responsavel_custo = locatario` (base da cobrança) | Qualquer outra OS |
| `ordens-servico/{id}/campos-extras` | Lê as definições com `preenchido_pelo_prestador = true` e grava os valores delas, nas OS no seu alcance (23.2) | Não acessa; `campos_extras` não sai na OS para o `locatario` (proposta) | Campos extras sem `preenchido_pelo_prestador`, OS fora do alcance |
| `solicitacoes` | Não acessa | Abre solicitação para ativo locado a ele (`origem = portal_locatario`); lista e detalha as próprias; comenta; desiste com `recusar` e motivo `desistencia` | Solicitações de outros usuários e de ativos não locados a ele |
| `ativos` | Detalha (resumo: TAG, nome, classe, marca, modelo, localização) só os ativos de OS no seu alcance; lê medidores e registra leituras desses ativos durante a execução | Lista e detalha (resumo) os ativos locados a ele, com medidores e últimas leituras | Demais ativos; `historico`, `documentos-tecnicos` (salvo os marcados como compartilháveis, proposta) |
| `leituras-medidor` | Registra leitura de ativo de OS no seu alcance (origem `os`) | Lê as leituras dos ativos locados a ele | Leituras de outros ativos |
| `checklists-execucao` | Registra e lê execuções ligadas a OS no seu alcance | Não acessa (as vistorias ficam em `vistorias-locacao`) | Demais execuções |
| `almoxarifados`, `saldos-estoque` | Só o almoxarifado `em_terceiro` do seu fornecedor e o saldo dele | Não acessa | Almoxarifados `proprio`, `quarentena` e de outros fornecedores |
| `movimentacoes-estoque` | Lista as do seu almoxarifado; registra `saida_os` e `devolucao_os` em OS no seu alcance e confirma recebimento de `remessa_terceiro` (20.2) | Não acessa | Movimentações de outros almoxarifados |
| `itens-estoque` | Resumo (código, descrição, unidade) dos itens com saldo no seu almoxarifado | Não acessa | Demais itens |
| `fornecedores` | Só o próprio fornecedor (leitura), com seus serviços | Não acessa | Outros fornecedores |
| `servicos-fornecedor` | Lê e mantém o catálogo do próprio fornecedor, se a empresa permitir (parâmetro, proposta) | Não acessa | Serviços de outros fornecedores |
| `orcamentos-fornecedor` | Cria e lê orçamentos do próprio fornecedor para OS no seu alcance | Não acessa | Orçamentos de outros fornecedores |
| `contratos-prestacao`, `avaliacoes-fornecedor` | Lê os contratos e avaliações do próprio fornecedor (Fase 2) | Não acessa | Os de outros fornecedores |
| `contratos-locacao` | Não acessa | Lê os próprios contratos (itens, vigência, local de uso, franquia, matriz de responsabilidade) e o relatório de cobrança | Contratos de outros locatários e todos os `tomado` |
| `vistorias-locacao` | Não acessa | Lê as vistorias dos seus contratos e registra a assinatura da sua parte | Vistorias de outros contratos |
| `pessoas`, `usuarios` (exceto `eu`), `unidades`, `localizacoes`, `classes-ativo`, planos, indicadores, conformidade, frota, integrações, auditoria | Não acessa | Não acessa | Sempre 404 |
| Configuração e módulos PRJ e CTR: `regras-automacao`, `disparos-regra`, `campos-extras-os`, `pacotes-segmento`, `aplicacoes-pacote`, `modelos-projeto`, `projetos` (e etapas), `medicoes-etapa`, `contratos-cliente` (e itens e painel), `medicoes-contrato`; ação `colocar-em-operacao` | Não acessa | Não acessa | Sempre 404, também para o escopo `link` |
| `ordens-reparo` | Lê as ordens de reparo enviadas ao seu fornecedor e registra orçamento (Fase 2, proposta) | Não acessa | Ordens de outros fornecedores |

Webhooks, tokens de API e exportação nunca são acessíveis a perfis externos. O app de campo (seção 12) monta o escopo offline do `prestador` com os mesmos filtros; o `locatario` usa só o portal web no lançamento do LOC (sumário 3.6), e o app recusa o login dele.

**Perfil `prestador` ou conta própria (SRV).** O perfil `prestador` dentro da conta do cliente continua existindo para subcontratadas sem conta na Tenaz. Quando o fornecedor tem conexão `ativa` com a conta própria dele (seção 22), as OS atribuídas a ele são executadas lá, pelos técnicos dele, e não aparecem para os usuários `prestador` desse fornecedor. O técnico da subcontratada com conta própria **não** conta como técnico cobrado no cliente (sumário 2.2); o usuário `prestador` que executa OS na conta do cliente continua contando (sumário 2.1).

---

## 5. Rate limit

### 5.1 Limites por plano (proposta)

O estudo de mercado recomenda **pelo menos 120 requisições por minuto no plano intermediário**, acima dos 50 req/min do Fleetio Professional. A proposta abaixo supera essa meta. As colunas seguem a tabela de recursos da seção 2.1 do sumário: API v1 e webhooks a partir do Profissional; web e app em todos os planos.

| Plano | Tokens de API (por empresa) | Rajada | Sessões web e app (por usuário) | Sincronização do app (por dispositivo) | Webhooks configuráveis |
|---|---|---|---|---|---|
| Essencial | Sem tokens de API | — | 300 req/min | 60 req/min | — |
| Profissional | **300 req/min** | 50 | 300 req/min | 60 req/min | 10 endpoints |
| Enterprise | **1.200 req/min** (ampliável em contrato) | 200 | 600 req/min | 120 req/min | 50 endpoints |
| Sandbox | 60 req/min | 20 | 120 req/min | 30 req/min | 5 endpoints |

Limites específicos, em todos os planos (proposta):

| Operação | Limite |
|---|---|
| `POST /v1/auth/login`, `POST /oauth/token` | 10 por minuto por IP e por e-mail |
| Listas com `incluir_total=true` | 30 por minuto por credencial |
| Ingestão de telemática (`POST /v1/telematica/eventos`) | 600 req/min por conexão, até 500 eventos por requisição |
| Consulta de bloqueio de despacho | 600 req/min por token, separado do limite geral |
| Ativação do link de uso único (`POST /v1/publico/links-execucao/ativar`) | 10 req/min por IP; 10 tokens inválidos seguidos do mesmo IP bloqueiam o IP por 15 minutos |
| Demais rotas públicas do link (`/v1/publico/links-execucao/atual/*`) | 60 req/min por sessão de link e 30 req/min por IP sem sessão válida |
| Criação de links de uso único | 100 por hora por empresa |
| Convites de conexão entre contas (`POST /v1/conexoes-empresa`) | 20 por dia por empresa; 5 tentativas de aceite com código errado por convite bloqueiam o convite, que precisa ser reenviado |
| Relatório da OS em PDF (`GET /v1/ordens-servico/{id}/relatorio-pdf`) | 30 por minuto por usuário |
| `GET /v1/versao` (público) | 60 req/min por IP |
| Simulação de regra configurável (`POST /v1/regras-automacao/{id}/simular`) | 10 por minuto por usuário |
| Pré-visualização e aplicação de pacote de segmento | 20 pré-visualizações por hora e 5 aplicações por dia, por empresa |
| Entrada em operação em lote (`POST /v1/ativos/colocar-em-operacao`) | Até 100 ativos por requisição; 10 requisições por minuto por usuário |
| Painel do contrato com cliente (`GET /v1/contratos-cliente/{id}/painel`) | 30 por minuto por usuário |
| Importações e relatórios assíncronos | 5 simultâneos por empresa |
| Exportação completa | 1 em andamento por empresa; 4 por dia |
| IA (consultas e rascunhos) | Pelos créditos do plano (seção 6.16) e 20 req/min por usuário |

### 5.2 Algoritmo e cabeçalhos

- Balde de fichas (token bucket) por chave (`empresa+token`, `usuario`, `dispositivo`), com capacidade igual à rajada e reposição contínua (proposta).
- Toda resposta traz os cabeçalhos do rascunho IETF "RateLimit header fields for HTTP":

```
RateLimit-Policy: 300;w=60;burst=50
RateLimit-Limit: 300
RateLimit-Remaining: 287
RateLimit-Reset: 41
```

`RateLimit-Reset` é o número de segundos até a reposição completa da janela.

### 5.3 Resposta 429

```
HTTP/1.1 429 Too Many Requests
Content-Type: application/problem+json
Retry-After: 12
RateLimit-Limit: 300
RateLimit-Remaining: 0
RateLimit-Reset: 12
```

```json
{
  "type": "https://api.tenaz.com.br/erros/limite_requisicoes_excedido",
  "title": "Limite de requisições excedido",
  "status": 429,
  "detail": "O limite de 300 requisições por minuto do token \"ERP Protheus\" foi atingido. Tente de novo em 12 segundos.",
  "codigo": "limite_requisicoes_excedido",
  "id_rastreio": "01JCZ5A1B2C3D4E5F6G7H8J9K0"
}
```

- `Retry-After` vem sempre em **segundos**. O cliente deve esperar pelo menos esse tempo, com variação aleatória (jitter).
- Requisições recusadas por limite não contam no balde.
- O web e o app tratam 429 com espera automática; o app nunca descarta comandos por 429.

---

## 6. Catálogo de endpoints por módulo

Legenda das colunas:

- **Perfil:** grupo da seção 4.8 (Todos, Leitura, Gestão, Execução, Estoque, Compras, Liberação, Destino de peça, Regras, Configuração, Entrada em operação, Projetos, Dispensa de etapa, Ajuste de medição, Contratos, Fechamento de medição, Reabertura de medição, Admin, Externos, Integração) ou perfil nomeado. `prestador` e `locatario` aparecem nomeados onde têm acesso, sempre com o alcance da seção 4.9. "Público" = sem autenticação; "Link" = sem login, autorizado pela sessão curta do link de uso único (19.4), salvo a ativação, autorizada pelo token no corpo.
- **Plano:** plano mínimo do recurso. `E` = Essencial, `P` = Profissional, `Ent` = Enterprise, `S` = Serviços (conta de subcontratada, sumário 2.2). `E, S` = disponível nos dois. As rotas `E` de ADM, ATV, MED, SOL, OS, CHK, estoque básico, IND essencial e MOB também valem no plano Serviços (22.2). Por token de API, o mínimo é sempre `P` (1.1).
- **Fase:** `MVP` (meses 0 a 6), `F2` (6 a 12), `F3` (12 a 24), conforme a seção 4 do sumário.
- Rotas em `/v1/` salvo indicação. Os métodos de recurso seguem a seção 2: `GET` lista com cursor e filtros, `POST` cria com `Idempotency-Key`, `PATCH` exige `If-Match`, `DELETE` inativa quando há histórico.

### 6.1 ADM — Administração

| Método | Rota | Descrição | Perfil | Plano | Fase |
|---|---|---|---|---|---|
| POST | `/v1/auth/login` | Login por e-mail e senha; cria sessão web (4.2) | Público | E | MVP |
| POST | `/v1/auth/mfa/verificar` | Conclui o login com o segundo fator | Público | E | MVP |
| POST | `/v1/auth/logout` | Encerra a sessão | Todos, Externos | E | MVP |
| GET | `/v1/auth/sessao` | Sessão: usuário, empresa, escopo, perfis, permissões, empresas do seletor e `csrf_token` | Todos, Externos | E | MVP |
| POST | `/v1/auth/trocar-empresa` | Nova sessão (web) ou novos tokens (app) em outra empresa a que o usuário tem acesso | Todos, Externos | E | MVP |
| POST | `/v1/auth/recuperar-senha` | Envia link de redefinição de senha; `202` sempre | Público | E | MVP |
| POST | `/v1/auth/redefinir-senha` | Define nova senha com o token do link | Público | E | MVP |
| POST | `/v1/auth/trocar-senha` | Troca a própria senha, conferindo a atual | Todos, Externos | E | MVP |
| POST | `/v1/auth/mfa/iniciar` | Gera o segredo TOTP pendente e a URI do QR code | Todos, Externos | E | MVP |
| POST | `/v1/auth/mfa/confirmar` | Ativa o segundo fator com o primeiro código certo | Todos, Externos | E | MVP |
| POST | `/v1/auth/mfa/desativar` | Desativa o segundo fator, com um código válido | Todos, Externos | E | MVP |
| POST | `/v1/auth/convites/consultar` | Confere o link do convite sem gastá-lo: `email` e `pede_senha` | Público | E | MVP |
| POST | `/v1/auth/convites/aceitar` | Aceita o convite: liga o login ao usuário convidado e o ativa | Público | E | MVP |
| GET | `/oauth/authorize` | Início do fluxo OAuth 2.1 com PKCE (app): `302` para a tela `/autorizar-app` da web | Público | E | MVP |
| POST | `/oauth/authorize` | Decisão da tela de autorização: emite o código para a empresa escolhida (só sessão web) | Todos, Externos | E | MVP |
| POST | `/oauth/token` | Troca de código e renovação rotativa de refresh token (formulário) | Público | E | MVP |
| POST | `/oauth/revoke` | Revoga refresh token; o token no corpo é a credencial | Público | E | MVP |
| GET | `/v1/versao` | Versão do produto em SemVer, commit, data do build e ambiente (18.8) | Público | E | MVP |
| POST | `/v1/empresas` | Cadastro da conta pelo site: empresa, unidade inicial, 12 perfis padrão e responsável como administrador convidado (RF-ADM-001, RF-SRV-001); responde só `convite_enviado`, sem identificadores; `409 conflito` em CNPJ já cadastrado | Público | E, S | MVP |
| GET | `/v1/empresa` | Dados da empresa, inclusive `mfa_obrigatorio` | Leitura | E | MVP |
| PATCH | `/v1/empresa` | Altera dados cadastrais e `mfa_obrigatorio` (segundo fator exigido de todos, RN-ADM-027); plano e SRV não mudam aqui (`422 campo_somente_leitura`) | Admin | E | MVP |
| GET | `/v1/empresa/plano` | Plano, recursos habilitados e uso (técnicos, veículos ativos, créditos de IA) | Admin | E | MVP |
| GET | `/v1/unidades` | Lista unidades | Leitura | E | MVP |
| POST | `/v1/unidades` | Cria unidade com pelo menos um turno. O Essencial tem 1 unidade; a segunda exige o Profissional (`403 plano_nao_inclui_recurso`) | Admin | E | MVP |
| GET | `/v1/unidades/{id}` | Detalha unidade | Leitura | E | MVP |
| PATCH | `/v1/unidades/{id}` | Altera unidade (nome, endereço, fuso, CNPJ) | Admin | E | MVP |
| DELETE | `/v1/unidades/{id}` | Inativa unidade; a última unidade ativa não é inativada (a regra "sem ativos ativos" entra com ATV, em E02) | Admin | E | MVP |
| POST | `/v1/unidades/{id}/reativar` | Reativa unidade inativa, respeitando o limite de unidades do plano | Admin | E | MVP |
| POST | `/v1/unidades/{id}/turnos` | Inclui turno (`turno_unidade`): `dia_semana` (1 = segunda ... 7 = domingo), `hora_inicio`, `hora_fim`, `nome` | Admin | E | MVP |
| PATCH | `/v1/unidades/{id}/turnos/{turno_id}` | Altera turno (ETag do turno) | Admin | E | MVP |
| DELETE | `/v1/unidades/{id}/turnos/{turno_id}` | Remove turno; o último turno da unidade não é removido (`422 validacao_falhou`) | Admin | E | MVP |
| GET | `/v1/unidades/{id}/feriados` | Feriados da unidade (`feriado_unidade`), com filtro `data[gte]` e `data[lte]` | Leitura | E | MVP |
| POST | `/v1/unidades/{id}/feriados` | Inclui feriado (`data`, `descricao`); `409 conflito` em data repetida | Admin | E | MVP |
| DELETE | `/v1/unidades/{id}/feriados/{feriado_id}` | Remove feriado (ETag do feriado) | Admin | E | MVP |
| GET | `/v1/pessoas` | Lista pessoas (`com_usuario`, `fornecedor_id`, `unidade_id`, `situacao`, `busca`; `especialidade_tecnico_id` entra com o SRV), com CPF mascarado | Gestão | E | MVP |
| POST | `/v1/pessoas` | Cadastra pessoa, com ou sem usuário (técnico, motorista, condutor eventual, técnico de fornecedor, PLH); aceita `especialidade_tecnico_ids` (22.3) | Gestão | E | MVP |
| GET | `/v1/pessoas/{id}` | Detalha pessoa, com o CPF completo (a leitura vai para a trilha) | Gestão, titular | E | MVP |
| PATCH | `/v1/pessoas/{id}` | Altera pessoa | Gestão | E | MVP |
| DELETE | `/v1/pessoas/{id}` | Inativa pessoa; ligada a usuário `ativo` ou `convidado`, responde `409 transicao_invalida`: inative o usuário antes (RN-ADM-021) | Gestão | E | MVP |
| POST | `/v1/pessoas/{id}/reativar` | Reativa pessoa inativa | Gestão | E | MVP |
| GET | `/v1/usuarios` | Lista usuários | Gestão | E | MVP |
| POST | `/v1/usuarios` | Cria usuário `convidado` e envia convite (7 dias); aceita `pessoa_id` ou `nova_pessoa` (cria a pessoa junto; nunca os dois), `fornecedor_id` (obrigatório para `prestador`), `locatario_id` (obrigatório para `locatario`), `custo_hora` e `unidade_padrao_id`; responde com `convite_enviado`. A pessoa ligada tem de estar ativa (`422`, `erros[].codigo = "inativo"`) e ter o mesmo fornecedor do usuário, nos dois sentidos: pessoa de fornecedor só com `prestador` dele, e usuário interno só com pessoa sem fornecedor; `nova_pessoa` de um `prestador` herda o `fornecedor_id` dele e o vínculo `terceiro` (RN-ADM-020, RN-ADM-021) | Admin | E | MVP |
| GET | `/v1/usuarios/eu` | Próprio usuário, perfis, unidades e preferências | Todos, Externos | E | MVP |
| PATCH | `/v1/usuarios/eu` | Altera os dados próprios (no MVP, `nome` e `telefone`) | Todos, Externos | E | MVP |
| GET | `/v1/usuarios/{id}` | Detalha usuário | Gestão | E | MVP |
| PATCH | `/v1/usuarios/{id}` | Altera perfis (lista completa em `perfis`), unidades e dados | Admin | E | MVP |
| POST | `/v1/usuarios/{id}/desativar` | `situacao` → `inativo`; revoga sessões e tokens do usuário; `409 ultimo_administrador` | Admin | E | MVP |
| POST | `/v1/usuarios/{id}/reativar` | `inativo` → `ativo` (ou `convidado`, com convite novo, se nunca aceitou) | Admin | E | MVP |
| POST | `/v1/usuarios/{id}/reenviar-convite` | Revoga o link anterior e envia outro (7 dias); só para `convidado` | Admin | E | MVP |
| POST | `/v1/usuarios/{id}/encerrar-sessoes` | Encerra todas as sessões e dispositivos | Admin | E | MVP |
| POST | `/v1/usuarios/{id}/permissoes` | Concede permissão avulsa (`usuario_permissao`), como `os.encerrar` a um técnico; nunca a si mesmo | Admin | E | MVP |
| DELETE | `/v1/usuarios/{id}/permissoes/{permissao}` | Retira permissão avulsa | Admin | E | MVP |
| GET | `/v1/perfis` | Lista perfis padrão e personalizados | Gestão | E | MVP |
| POST | `/v1/perfis` | Cria perfil personalizado (proposta) | Admin | P | F2 |
| PATCH | `/v1/perfis/{id}` | Altera perfil personalizado | Admin | P | F2 |
| DELETE | `/v1/perfis/{id}` | Exclui perfil personalizado sem usuários | Admin | P | F2 |
| GET | `/v1/permissoes` | Catálogo de permissões finas | Admin | E | MVP |
| GET | `/v1/parametros-empresa` | Parâmetros da empresa (MTTR com ou sem espera, janela de aderência, alertas D-60/D-30/D-7, sessão, validade do link de uso único) | Gestão | E | MVP |
| PATCH | `/v1/parametros-empresa` | Altera parâmetros | Admin | E | MVP |
| GET | `/v1/eventos-auditoria` | Trilha de auditoria com filtros por entidade, registro, usuário e período | Admin, `gestor` | E | MVP |
| GET | `/v1/eventos-auditoria/{id}` | Evento com valores antes e depois | Admin, `gestor` | E | MVP |
| GET | `/v1/provedores-identidade` | Lista configurações de SSO | Admin | P | MVP |
| POST | `/v1/provedores-identidade` | Configura OIDC (P) ou SAML (Ent) | Admin | P | MVP |
| PATCH | `/v1/provedores-identidade/{id}` | Altera configuração de SSO | Admin | P | MVP |
| DELETE | `/v1/provedores-identidade/{id}` | Remove SSO (exige conta de emergência ativa) | Admin | P | MVP |
| POST | `/v1/provedores-identidade/{id}/token-scim` | Emite token SCIM | Admin | Ent | F2 |
| GET | `/saml/{provedor_id}/metadata` | Metadados do SP | Público | Ent | F2 |
| POST | `/saml/{provedor_id}/acs` | Recebe asserção SAML | Público | Ent | F2 |
| GET | `/scim/v2/Users` | Lista e filtra usuários (SCIM) | Token SCIM | Ent | F2 |
| POST | `/scim/v2/Users` | Cria usuário (SCIM) | Token SCIM | Ent | F2 |
| GET | `/scim/v2/Users/{id}` | Detalha usuário (SCIM) | Token SCIM | Ent | F2 |
| PUT | `/scim/v2/Users/{id}` | Substitui usuário (SCIM) | Token SCIM | Ent | F2 |
| PATCH | `/scim/v2/Users/{id}` | Altera usuário, inclusive `active` (SCIM) | Token SCIM | Ent | F2 |
| DELETE | `/scim/v2/Users/{id}` | Desativa usuário (SCIM; nunca apaga histórico) | Token SCIM | Ent | F2 |
| GET | `/scim/v2/Groups` | Lista grupos mapeados a perfis | Token SCIM | Ent | F2 |
| POST | `/scim/v2/Groups` | Cria grupo | Token SCIM | Ent | F2 |
| PATCH | `/scim/v2/Groups/{id}` | Altera membros do grupo | Token SCIM | Ent | F2 |
| DELETE | `/scim/v2/Groups/{id}` | Remove grupo | Token SCIM | Ent | F2 |
| POST | `/v1/anexos` | Registra anexo e devolve URL de upload | Todos | E | MVP |
| POST | `/v1/anexos/{id}/confirmar` | Confirma upload e confere `sha256` | Todos | E | MVP |
| GET | `/v1/anexos/{id}` | Metadados e URL temporária de download | Todos | E | MVP |
| GET | `/v1/processamentos/{id}` | Situação de operação assíncrona | Todos | E | MVP |
| GET | `/v1/notificacoes` | Notificações do usuário | Todos | E | MVP |
| POST | `/v1/notificacoes/{id}/marcar-lida` | Marca notificação como lida | Todos | E | MVP |
| POST | `/v1/exportacoes-empresa` | Inicia exportação completa (seção 17) | Admin | E | MVP |
| GET | `/v1/exportacoes-empresa` | Lista exportações | Admin | E | MVP |
| GET | `/v1/exportacoes-empresa/{id}` | Situação e links de download | Admin | E | MVP |
| DELETE | `/v1/exportacoes-empresa/{id}` | Cancela ou apaga o arquivo gerado | Admin | E | MVP |
| PUT | `/v1/exportacoes-empresa/agendamento` | Define exportação periódica | Admin | E | F2 |
| GET | `/v1/regras-automacao` | Lista regras configuráveis (`situacao`, `condicao`, `efeito`, `classe_ativo_id`, `unidade_id`, `origem_pacote`, `busca`) | Gestão | E | MVP |
| POST | `/v1/regras-automacao` | Cria regra `inativa`: condição e efeito dos catálogos fechados, parâmetros e filtro (23.1) | Regras | E | MVP |
| GET | `/v1/regras-automacao/{id}` | Detalha regra, com contagem de disparos e bloqueios vigentes criados por ela | Gestão | E | MVP |
| PATCH | `/v1/regras-automacao/{id}` | Altera nome, descrição, parâmetros e filtro; a alteração vale para os próximos fatos | Regras | E | MVP |
| DELETE | `/v1/regras-automacao/{id}` | Apaga regra que nunca disparou; regra com disparos é só inativada (2.12) | Regras | E | MVP |
| POST | `/v1/regras-automacao/{id}/ativar` | `inativa` → `ativa`; respeita o limite de regras ativas (`422 limite_regras_ativas_excedido`) | Regras | E | MVP |
| POST | `/v1/regras-automacao/{id}/inativar` | `ativa` → `inativa`; não levanta bloqueios já criados e devolve a lista deles em `bloqueios_vigentes` | Regras | E | MVP |
| POST | `/v1/regras-automacao/{id}/simular` | Avalia a regra agora, sem gravar nada nem aplicar efeito: ativos e fatos que disparariam (23.1) | Gestão | E | MVP |
| GET | `/v1/regras-automacao/{id}/disparos` | Histórico de disparos da regra | Gestão | E | MVP |
| GET | `/v1/disparos-regra` | Histórico de disparos de todas as regras (`ativo_id`, `regra_automacao_id`, `efeito`, `condicao`, período) | Gestão | E | MVP |
| GET | `/v1/disparos-regra/{id}` | Detalha disparo: fato, ativo, efeito aplicado e registros criados (solicitação, bloqueio, notificação) | Gestão | E | MVP |
| GET | `/v1/pacotes-segmento` | Catálogo global de pacotes de segmento (`segmento`, `busca`), só leitura | Admin, Gestão | E | F2 |
| GET | `/v1/pacotes-segmento/{codigo}` | Conteúdo do pacote na versão publicada mais recente (`versao` escolhe outra) | Admin, Gestão | E | F2 |
| POST | `/v1/pacotes-segmento/{codigo}/pre-visualizar` | Mostra, sem gravar, o que a aplicação criaria e o que ignoraria (já existe, módulo não contratado) (23.4) | Configuração | E | F2 |
| POST | `/v1/pacotes-segmento/{codigo}/aplicar` | Aplica o pacote: só cria o que falta, numa transação; devolve a `aplicacao_pacote` | Configuração | E | F2 |
| GET | `/v1/aplicacoes-pacote` | Aplicações de pacote da conta (`pacote_codigo`, período) | Admin, Gestão | E | F2 |
| GET | `/v1/aplicacoes-pacote/{id}` | Detalha a aplicação: versão, quem aplicou, itens criados e itens ignorados com o motivo | Admin, Gestão | E | F2 |

**Permissão por rota (E01).** O grupo da coluna Perfil é a forma resumida; a permissão conferida pelo servidor é a de RF-ADM-005 e sai no contrato em `x-permissoes`: `GET /v1/empresa`, `GET /v1/unidades`, `GET /v1/unidades/{id}` e `GET .../feriados`, `adm.ler`; `PATCH /v1/empresa` e a escrita de unidades, turnos e feriados, `adm.gerenciar_empresa`; leitura de pessoas, `adm.ler_pessoas` (o titular lê a própria pessoa sem ela); escrita de pessoas, `adm.gerenciar_pessoas`; `GET /v1/usuarios`, `GET /v1/usuarios/{id}` e `GET /v1/perfis`, `adm.ler_usuarios`; escrita de usuários, ações e permissões avulsas, `adm.gerenciar_usuarios`; `GET /v1/permissoes`, `adm.gerenciar_perfis`; trilha, `adm.ler_auditoria`. As rotas de `/v1/auth/*`, `/oauth/*` e `/v1/usuarios/eu` exigem só a credencial e aceitam os escopos `interno`, `prestador` e `locatario` ("Todos, Externos"); todas as demais rotas de ADM aceitam só `interno` e respondem `404` aos externos. **Situação do usuário:** `usuario.situacao` é `convidado`, `ativo` ou `inativo` (sumário 3.6); o bloqueio temporário por falhas de login é da identidade, não do usuário. **Situação e vínculo da pessoa:** `pessoa.situacao` é `ativa` ou `inativa`, sem `desligada`; `pessoa.vinculo` é `empregado`, `terceiro` ou `autonomo` (sumário 3.6). **Sessão restrita:** em empresa com `mfa_obrigatorio`, a pessoa sem MFA só usa as cinco rotas de sessão e ativação listadas na seção 4.2; as demais rotas desta tabela respondem `401 nao_autenticado` com `mfa_pendente`. **Dinheiro:** `custo_hora` do usuário é texto decimal (2.5).

**Pessoa e usuário.** `pessoa` é quem executa, conduz ou responde por algo, com ou sem login; `usuario` é o login e pode apontar para uma pessoa (`usuario.pessoa_id`, opcional). CNH, capacitações e documentos pessoais ligam-se à pessoa (`pessoa_id` em `documentos-vencimento` e `capacitacoes`), nunca ao usuário. Condutor eventual sem login é uma pessoa sem usuário, não um usuário sem acesso. A Tenaz não guarda dado de saúde (sumário 3.6): não existe tipo de documento, campo nem anexo para ASO, exame toxicológico ou outro dado de saúde da pessoa.

**Regras configuráveis e pacotes de segmento** (RF-ADM-016, RF-ADM-017; sumário 3.6 e 3.7). Ficam em ADM, no pacote Go `configuracao`. Regras exigem `adm.gerenciar_regras`; a aplicação de pacote, `adm.gerenciar_configuracao`. As regras e os campos extras são do Essencial, no MVP; os pacotes são do Essencial, na Fase 2. As rotas `E` também valem no plano Serviços (22.2). O detalhamento está na seção 23.

### 6.2 ATV — Ativos

| Método | Rota | Descrição | Perfil | Plano | Fase |
|---|---|---|---|---|---|
| GET | `/v1/localizacoes` | Lista localizações (`pai_id`, `unidade_id`) | Todos | E | MVP |
| POST | `/v1/localizacoes` | Cria localização | Gestão | E | MVP |
| GET | `/v1/localizacoes/{id}` | Detalha localização | Todos | E | MVP |
| PATCH | `/v1/localizacoes/{id}` | Altera localização | Gestão | E | MVP |
| DELETE | `/v1/localizacoes/{id}` | Inativa localização sem ativos | Gestão | E | MVP |
| GET | `/v1/classes-ativo` | Lista classes de ativo | Leitura | E | MVP |
| POST | `/v1/classes-ativo` | Cria classe (campos extras, planos padrão, códigos de falha) | Gestão | E | MVP |
| GET | `/v1/classes-ativo/{id}` | Detalha classe | Leitura | E | MVP |
| PATCH | `/v1/classes-ativo/{id}` | Altera classe | Gestão | E | MVP |
| DELETE | `/v1/classes-ativo/{id}` | Inativa classe sem ativos | Gestão | E | MVP |
| GET | `/v1/ativos` | Lista ativos (`tipo`, `classe_ativo_id`, `localizacao_id`, `criticidade`, `status_operacional`, `fase_ciclo_vida`, `propriedade`, `em_locacao`, `contrato_locacao_id`, `cliente_servico_id`, `conexao_empresa_id`, `sincronizado`, `projeto_id`, `pai_id`, `tag`, `placa`, `nr12`) | Todos, `prestador`, `locatario` | E | MVP |
| POST | `/v1/ativos` | Cria ativo; com `propriedade = de_cliente`, exige `cliente_servico_id` (conta SRV, 22.4); aceita `fase_ciclo_vida = em_implantacao` (23.3) | Gestão | E | MVP |
| GET | `/v1/ativos/{id}` | Detalha ativo | Todos, `prestador`, `locatario` | E | MVP |
| PATCH | `/v1/ativos/{id}` | Altera ativo | Gestão | E | MVP |
| DELETE | `/v1/ativos/{id}` | Inativa ativo (status `inativo`) | Gestão | E | MVP |
| GET | `/v1/ativos/{id}/arvore` | Descendentes até o 6º nível | Todos | E | MVP |
| POST | `/v1/ativos/{id}/mover` | Troca pai ou localização (valida profundidade e ciclo) | Gestão | E | MVP |
| POST | `/v1/ativos/{id}/alterar-status-operacional` | Inativa (`inativo`, com motivo: baixa, venda, reserva) ou reativa o ativo; `em_manutencao` e `parado` são derivados das causas e não aceitam alteração manual; ativo `em_implantacao` não é reativado por aqui (`422 ativo_em_implantacao`) | Gestão | E | MVP |
| POST | `/v1/ativos/{id}/colocar-em-operacao` | Ação `colocar_em_operacao`: `em_implantacao` → `em_operacao`, status `operando`, grava `data_entrada_operacao` e ativa os planos (23.3) | Entrada em operação | E | MVP |
| POST | `/v1/ativos/colocar-em-operacao` | A mesma ação em lote, até 100 ativos, com resultado por ativo (23.3) | Entrada em operação | E | MVP |
| GET | `/v1/ativos/{id}/historico` | Linha do tempo: OS, leituras, checklists, documentos, peças | Leitura | E | MVP |
| GET | `/v1/ativos/qr/{codigo_qr}` | Resolve o QR code para o ativo | Todos | E | MVP |
| GET | `/v1/ativos/{id}/etiqueta-qr` | Etiqueta em PNG ou PDF | Gestão | E | MVP |
| POST | `/v1/ativos/etiquetas-qr` | Etiquetas em lote, em PDF (assíncrono) | Gestão | E | MVP |
| GET | `/v1/ativos/{id}/documentos-tecnicos` | Manuais, desenhos e fichas do ativo | Leitura | E | MVP |
| POST | `/v1/ativos/{id}/documentos-tecnicos` | Vincula documento técnico (anexo) | Gestão | E | MVP |
| DELETE | `/v1/ativos/{id}/documentos-tecnicos/{documento_id}` | Remove o vínculo | Gestão | E | MVP |

**Status operacional** (sumário 3.3). O campo `status_operacional` é só leitura no `PATCH` e assume `operando`, `em_manutencao`, `parado` ou `inativo`:

| Status | Quando a API o atribui |
|---|---|
| `operando` | Nenhuma das causas abaixo |
| `em_manutencao` | Há OS em `em_execucao` ou `pausada` com o ativo retirado de operação de forma programada (preventiva, preditiva, inspeção legal, corretiva planejada com `ativo_parado = true`) |
| `parado` | Falha que parou o ativo (solicitação ou OS com `ativo_parado = true` por falha), `corretiva_emergencial` em curso, registro NR-12 com indicação `inapto` ou bloqueio de conformidade (`bloqueio_ativo`) |
| `inativo` | Baixado, vendido ou em reserva (`alterar-status-operacional`), ou em implantação (`fase_ciclo_vida = em_implantacao`, 23.3) |

- `parado` prevalece sobre `em_manutencao`.
- **Em implantação** (sumário 3.3): o ativo `em_implantacao` tem `status_operacional = inativo`, fica fora dos KPIs e da contagem de veículos cobrados e tem os planos suspensos; a interface mostra "Em implantação". Só `colocar-em-operacao` o leva a `operando`. A fase não volta para `em_implantacao`.
- `bloqueio_ativo` com `origem = regra` (efeito `bloquear_ativo` ou `bloquear_ativo_ate_documento` de uma regra configurável, 23.1) deixa o ativo `parado` como o bloqueio manual.
- O status volta a `operando` no **encerramento** da OS, se não restar outra causa. A ação `encerrar` não recebe mais o status do ativo; a resposta informa o status resultante em `status_operacional_ativo`.
- MTTR e disponibilidade usam `data_inicio`, `data_conclusao` e os registros de parada, não o encerramento administrativo.
- **Atenção** não é status: é calculada sobre `operando` e vem no campo só leitura `alertas` (`plano_vencendo`, `documento_vencendo`, `leitura_fora_faixa`).
- **Apto para despacho** (veículo) é verificação separada do status (seção 16).
- Ativo **cedido** em locação vigente continua `operando`; sua localização exibida passa a ser o local de uso do contrato (`contrato_locacao.local_uso`, sumário 3.6), e a localização de cadastro continua em `localizacao_id`.
- Ativo `de_terceiro` devolvido ao locador no encerramento do contrato `tomado` passa a `inativo` (sumário 3.6).
- `bloqueio_ativo` é a entidade de bloqueio de conformidade que deixa o ativo `parado` (`/v1/bloqueios-ativo`, 6.13). "Apto para despacho" não é entidade: é a verificação calculada da seção 16. `bloqueio_energia` é o registro LOTO de uma intervenção, ligado à OS.

**Campos do ativo incluídos nesta versão:**

| Campo | Tipo | Regra |
|---|---|---|
| `propriedade` | enum | `proprio` (padrão), `de_terceiro` (tomado em locação ou em comodato) ou `de_cliente` (equipamento de um cliente atendido, na conta de uma subcontratada). Ativo de contrato `tomado` recebe `de_terceiro`. `de_cliente` só existe em conta com o módulo SRV |
| `cliente_servico_id` | UUID | Obrigatório com `propriedade = de_cliente`; proibido nas demais (`422 validacao_falhou`) |
| `conexao_empresa_id` | UUID, só leitura | Na cópia sincronizada de um ativo do cliente conectado: a conexão de origem. `null` em equipamento cadastrado pela própria subcontratada |
| `ativo_remoto_id` | — | Coluna da cópia sincronizada com o id do ativo na conta do cliente. Usada só pelo processo de sincronização: **não sai na API** e não é aceita em filtro nem em rota (nenhum id de outra conta é exposto) |
| `sincronizado` | booleano, só leitura | `true` quando `conexao_empresa_id` está preenchido e a conexão está `ativa`; os campos vindos do cliente ficam só leitura (`409 ativo_sincronizado_somente_leitura`, 22.4) |
| `valor_reposicao` | dinheiro (string) | Valor de reposição (RAV), base do indicador `custo_rav` |
| `horas_operacao_dia` | decimal (string) | Regime de operação, de 0 a 24; base de MTBF e disponibilidade |
| `dias_operacao_semana` | inteiro | Regime de operação, de 1 a 7 |
| `contrato_locacao_id` | UUID, só leitura | Contrato de locação vigente que inclui o ativo, cedido ou tomado; `null` fora de locação |
| `alertas` | lista, só leitura | Condições de Atenção calculadas |
| `fase_ciclo_vida` | enum | `em_implantacao` ou `em_operacao` (padrão). Aceito só na criação (`POST /v1/ativos`, importação ou projeto); depois é só leitura e muda apenas por `colocar-em-operacao` (`422 campo_somente_leitura` no `PATCH`). Vale também para ativo `de_cliente` e `de_terceiro` (RF-ATV-013) |
| `data_entrada_operacao` | data-hora, só leitura | Gravada por `colocar-em-operacao`; `null` enquanto `em_implantacao`. Ativo criado já `em_operacao` recebe a data de criação (proposta) |
| `colocado_em_operacao_por_id` | UUID, só leitura | Usuário que executou a ação (proposta) |
| `projeto_id` | UUID, só leitura | Projeto de implantação de origem (PRJ, Fase 2); `null` sem projeto |
| `termo_entrega_anexo_id` | UUID | Documento de entrega ou termo de recebimento, opcional, informado em `colocar-em-operacao` (proposta) |

A fase de ciclo de vida **não** é campo compartilhado na conexão entre contas no lançamento (RF-ATV-013): a cópia `de_cliente` na conta da subcontratada fica com `fase_ciclo_vida = em_operacao` (ver 26.1, A-31).

A árvore de ativos e a de localizações têm até 6 níveis (`422 hierarquia_profundidade_excedida`). Pneu é entidade própria (`/v1/pneus`), não ativo.

### 6.3 MED — Medidores

| Método | Rota | Descrição | Perfil | Plano | Fase |
|---|---|---|---|---|---|
| GET | `/v1/medidores` | Lista medidores (`ativo_id`, `tipo`) | Leitura, `motorista`, `prestador`, `locatario` | E | MVP |
| POST | `/v1/medidores` | Cria medidor (`km`, `horas`, `ciclos`, `outro`) | Gestão | E | MVP |
| GET | `/v1/medidores/{id}` | Detalha medidor com última leitura válida | Leitura, `motorista`, `prestador`, `locatario` | E | MVP |
| PATCH | `/v1/medidores/{id}` | Altera medidor (teto de salto, unidade, casas decimais) | Gestão | E | MVP |
| DELETE | `/v1/medidores/{id}` | Inativa medidor | Gestão | E | MVP |
| POST | `/v1/medidores/{id}/trocar` | Registra troca ou virada do contador (leitura final e inicial) | Gestão | E | MVP |
| GET | `/v1/leituras-medidor` | Lista leituras (`medidor_id`, `ativo_id`, `origem`, `status_validacao`, período) | Leitura, `locatario` | E | MVP |
| POST | `/v1/leituras-medidor` | Registra leitura (seção 9) | Execução, `motorista`, `prestador` | E | MVP |
| POST | `/v1/leituras-medidor/lote` | Registra até 500 leituras de integração | Integração (`medidores:escrita`) | P | MVP |
| GET | `/v1/leituras-medidor/{id}` | Detalha leitura | Leitura | E | MVP |
| POST | `/v1/leituras-medidor/{id}/validar` | Aceita leitura `suspeita` | Gestão | E | MVP |
| POST | `/v1/leituras-medidor/{id}/invalidar` | Invalida leitura com motivo; recalcula vencimentos de planos | Gestão | E | MVP |

### 6.4 SOL — Solicitações

| Método | Rota | Descrição | Perfil | Plano | Fase |
|---|---|---|---|---|---|
| GET | `/v1/solicitacoes` | Lista solicitações (solicitante, motorista e locatário veem só as próprias, proposta) | Todos, `locatario` | E | MVP |
| POST | `/v1/solicitacoes` | Abre solicitação (seção 7); o locatário só para ativo locado a ele | Todos, `locatario` | E | MVP |
| GET | `/v1/solicitacoes/{id}` | Detalha solicitação | Todos, `locatario` | E | MVP |
| PATCH | `/v1/solicitacoes/{id}` | Altera descrição e anexos; o autor só enquanto `nova` | Todos, Gestão, `locatario` | E | MVP |
| POST | `/v1/solicitacoes/{id}/iniciar-triagem` | `nova` → `em_triagem`; grava o instante de reconhecimento | Gestão | E | MVP |
| POST | `/v1/solicitacoes/{id}/aprovar` | → `aprovada` e cria a OS | Gestão | E | MVP |
| POST | `/v1/solicitacoes/{id}/recusar` | `nova` ou `em_triagem` → `recusada`, com motivo; o autor (inclusive o `locatario`) pode desistir da própria, direto de `nova`, com motivo `desistencia` | Gestão, autor | E | MVP |
| POST | `/v1/solicitacoes/{id}/marcar-duplicada` | `nova` ou `em_triagem` → `duplicada` (estado final), vinculada à solicitação original ou a uma OS | Gestão | E | MVP |
| POST | `/v1/solicitacoes/{id}/reabrir-como-nova` | Corrige marcação de duplicada errada: cria solicitação nova (`nova`) com os dados da duplicada; a original continua `duplicada` (7.1) | Gestão | E | MVP |
| GET | `/v1/solicitacoes/{id}/comentarios` | Lista comentários | Todos, `locatario` | E | MVP |
| POST | `/v1/solicitacoes/{id}/comentarios` | Comenta (com anexos) | Todos, `locatario` | E | MVP |

### 6.5 OS — Ordens de serviço

| Método | Rota | Descrição | Perfil | Plano | Fase |
|---|---|---|---|---|---|
| GET | `/v1/ordens-servico` | Lista OS (`status`, `tipo`, `origem`, `ativo_id`, `responsavel_id`, `prioridade`, `unidade_id`, `executor_tipo`, `fornecedor_id`, `responsavel_custo`, `contrato_locacao_id`, `cliente_servico_id`, `vinculada`, `estado_sincronizacao`, `sem_responsavel`, `causa_externa`, `projeto_id`, `etapa_projeto_id`, `campos_extras.<chave>`, períodos) | Leitura, `prestador`, `locatario` | E, S | MVP |
| POST | `/v1/ordens-servico` | Cria OS (seção 8.1; OS de serviço do SRV em 22.7) | Execução | E, S | MVP |
| GET | `/v1/ordens-servico/{id}` | Detalha OS; OS com vínculo entre contas traz o objeto `vinculo` (22.8) | Leitura, `prestador`, `locatario` | E, S | MVP |
| PATCH | `/v1/ordens-servico/{id}` | Altera campos editáveis no estado atual; em OS vinculada, só os campos do próprio lado (`409 os_vinculada_somente_leitura`) | Execução | E, S | MVP |
| POST | `/v1/ordens-servico/{id}/planejar` | `aberta` → `planejada` | Gestão | E | MVP |
| POST | `/v1/ordens-servico/{id}/programar` | `planejada` ou `aberta` → `programada` (data e responsável) | Gestão | E | MVP |
| POST | `/v1/ordens-servico/{id}/desprogramar` | `programada` → `planejada` | Gestão | E | MVP |
| POST | `/v1/ordens-servico/{id}/reconhecer` | Registra `data_reconhecimento` sem mudar estado | Execução, `prestador` | E | MVP |
| POST | `/v1/ordens-servico/{id}/assumir` | O técnico se atribui OS `aberta` ou `planejada` da sua equipe, sem mudar estado; grava `data_reconhecimento` (8.2) | Execução, `prestador` | E, S | MVP |
| POST | `/v1/ordens-servico/{id}/iniciar` | `programada` → `em_execucao`; `aberta` → `em_execucao` só em `corretiva_emergencial` | Execução, `prestador` | E | MVP |
| POST | `/v1/ordens-servico/{id}/pausar` | `em_execucao` → `pausada`, com motivo | Execução, `prestador` | E | MVP |
| POST | `/v1/ordens-servico/{id}/retomar` | `pausada` → `em_execucao` | Execução, `prestador` | E | MVP |
| POST | `/v1/ordens-servico/{id}/concluir` | `em_execucao` → `concluida` | Execução, `prestador` | E | MVP |
| POST | `/v1/ordens-servico/{id}/encerrar` | Libera: `concluida` → `encerrada` (NR-12 quando aplicável). Exige `os.encerrar`; em OS vinculada, só o lado do cliente | Liberação | E, S | MVP |
| POST | `/v1/ordens-servico/{id}/recusar-liberacao` | `concluida` → `em_execucao`, com `motivo_recusa_liberacao`. Exige `os.encerrar`; em OS vinculada, só o lado do cliente | Liberação | E, S | MVP |
| POST | `/v1/ordens-servico/{id}/cancelar` | `aberta`, `planejada`, `programada` ou `pausada` → `cancelada`, com motivo | Gestão | E | MVP |
| POST | `/v1/ordens-servico/{id}/alterar-responsavel-custo` | Troca o `responsavel_custo`, com justificativa (8.6) | Gestão | E | MVP |
| GET | `/v1/ordens-servico/{id}/historico` | Transições de estado com autor, data e motivo | Leitura, `prestador` | E | MVP |
| GET | `/v1/ordens-servico/{id}/tarefas` | Lista tarefas | Leitura, `prestador` | E | MVP |
| POST | `/v1/ordens-servico/{id}/tarefas` | Adiciona tarefa ou aplica procedimento | Execução | E | MVP |
| PATCH | `/v1/ordens-servico/{id}/tarefas/{tarefa_id}` | Altera ou marca tarefa como concluída | Execução, `prestador` | E | MVP |
| DELETE | `/v1/ordens-servico/{id}/tarefas/{tarefa_id}` | Remove tarefa não iniciada | Gestão | E | MVP |
| GET | `/v1/ordens-servico/{id}/apontamentos` | Lista apontamentos de horas | Leitura, `prestador` | E | MVP |
| POST | `/v1/ordens-servico/{id}/apontamentos` | Aponta horas | Execução, `prestador` | E | MVP |
| PATCH | `/v1/ordens-servico/{id}/apontamentos/{apontamento_id}` | Corrige apontamento (antes do encerramento) | Execução, `prestador` | E | MVP |
| DELETE | `/v1/ordens-servico/{id}/apontamentos/{apontamento_id}` | Exclui apontamento (antes do encerramento) | Gestão | E | MVP |
| GET | `/v1/ordens-servico/{id}/servicos` | Serviços a realizar da OS (`servico_os`) e `descricao_servico` | Leitura, `prestador` | E, S | MVP |
| POST | `/v1/ordens-servico/{id}/servicos` | Inclui serviço a realizar: `servico_fornecedor_id` do catálogo da própria conta, quantidade, preço de referência (22.7) | Execução | E, S | MVP |
| PATCH | `/v1/ordens-servico/{id}/servicos/{servico_os_id}` | Altera quantidade, preço ou marca o serviço como realizado | Execução | E, S | MVP |
| DELETE | `/v1/ordens-servico/{id}/servicos/{servico_os_id}` | Remove serviço ainda não realizado | Execução | E, S | MVP |
| GET | `/v1/ordens-servico/{id}/pecas` | Peças e componentes previstos, consumidos, substituídos e fornecidos pelo cliente | Leitura, `prestador` | E, S | MVP |
| POST | `/v1/ordens-servico/{id}/pecas` | Prevê, consome ou registra componente com `tipo_material` (`consumido`, `substituido`, `fornecido_pelo_cliente`), `numero_serie_retirado` e `numero_serie_instalado` (22.7). Consumo de item de estoque gera `saida_os`; item serializado exige `pecas_serializadas_ids`; saldo `em_transito` não pode ser consumido (`422 estoque_em_transito`) | Execução, Estoque, `prestador` | E, S | MVP |
| POST | `/v1/ordens-servico/{id}/pecas/{peca_id}/estornar` | Devolve peça ao estoque (`devolucao_os`) | Execução, Estoque, `prestador` | E | MVP |
| GET | `/v1/ordens-servico/{id}/comentarios` | Lista comentários e fotos | Leitura, `prestador`, `locatario` | E | MVP |
| POST | `/v1/ordens-servico/{id}/comentarios` | Comenta (com anexos) | Execução, `prestador`, `locatario` | E | MVP |
| POST | `/v1/ordens-servico/{id}/assinaturas` | Registra assinatura (executante, responsável, cliente interno, prestador, cliente atendido) | Execução, `prestador` | E, S | MVP |
| GET | `/v1/ordens-servico/{id}/relatorio-pdf` | Relatório da OS em PDF: serviço, componentes, horas, fotos e assinatura do cliente (22.7) | Leitura, `prestador` | E, S | MVP |
| GET | `/v1/ordens-servico/{id}/custos` | Custos de mão de obra, peças e terceiros, com `responsavel_custo` | Gestão, `locatario` (4.9) | E | MVP |
| GET | `/v1/ordens-servico/{id}/registro-nr12` | Registro do item 12.11.2 em JSON ou PDF | Leitura | E | MVP |
| GET | `/v1/codigos-falha` | Catálogo de modo, causa e atividade (catálogo padrão em todos os planos) | Leitura, `prestador` | E | MVP |
| POST | `/v1/codigos-falha` | Cria código de falha (catálogo personalizável) | Gestão | P | F2 |
| PATCH | `/v1/codigos-falha/{id}` | Altera código | Gestão | P | F2 |
| DELETE | `/v1/codigos-falha/{id}` | Inativa código | Gestão | P | F2 |
| GET | `/v1/campos-extras-os` | Definições de campos extras da OS (`tipo_os`, `classe_ativo_id`, `situacao`, `obrigatorio_em`, `preenchido_pelo_prestador`) | Leitura | E, S | MVP |
| POST | `/v1/campos-extras-os` | Cria definição: chave, rótulo, tipo, opções, unidade, tipos de OS, classes de ativo, obrigatoriedade, ajuda (23.2) | Configuração | E, S | MVP |
| GET | `/v1/campos-extras-os/{id}` | Detalha definição | Leitura | E, S | MVP |
| PATCH | `/v1/campos-extras-os/{id}` | Altera rótulo, opções, ajuda, aplicação e obrigatoriedade; `chave` e `tipo` ficam só leitura depois de haver valor gravado | Configuração | E, S | MVP |
| DELETE | `/v1/campos-extras-os/{id}` | Inativa definição; os valores gravados nas OS são mantidos | Configuração | E, S | MVP |
| PUT | `/v1/campos-extras-os/ordem` | Substitui a ordem de exibição por inteiro (lista de ids) (proposta) | Configuração | E, S | MVP |
| GET | `/v1/ordens-servico/{id}/campos-extras` | Definições que se aplicam à OS (tipo e classe do ativo), com os valores atuais e o rótulo vigente na data do registro | Leitura, `prestador` | E, S | MVP |
| PATCH | `/v1/ordens-servico/{id}/campos-extras` | Grava valores de campos extras (JSON Merge Patch sobre `campos_extras`, `If-Match` da OS); é a única escrita na OS aberta ao `prestador` além das ações de execução | Execução, `prestador` | E, S | MVP |

As ações de atribuição a prestador (`atribuir-fornecedor`, `retirar-fornecedor`) também ficam em `/v1/ordens-servico/{id}/`, mas são listadas e contadas no módulo TER (6.9). `assumir`, `servicos` e `relatorio-pdf` servem a qualquer conta e são contados aqui, embora nasçam da necessidade do SRV (seção 22).

**Campos extras, tipo `implantacao` e causa externa** (RF-OS-021, RF-OS-002, RF-OS-001; sumário 3.6 e 3.7). As definições ficam em `/v1/campos-extras-os` e os valores em `ordem_servico.campos_extras`, gravados pela chave da definição (23.2). A OS `implantacao` instala um ativo ou executa uma etapa de projeto, fica fora dos KPIs de manutenção e tem custo de implantação (23.3 e 24.5). `causa_externa` é opcional e só tem efeito na medição de contrato com cliente (25.4). A ação que gera a OS de uma etapa de projeto (`gerar-os`) fica em `/v1/projetos/{id}/etapas/{etapa_id}/` e é contada em PRJ (6.18).

### 6.6 PRE — Planos de manutenção

| Método | Rota | Descrição | Perfil | Plano | Fase |
|---|---|---|---|---|---|
| GET | `/v1/planos-manutencao` | Lista planos (`ativo_id`, `classe_ativo_id`, `situacao`) | Leitura | E | MVP |
| POST | `/v1/planos-manutencao` | Cria plano (em rascunho) para ativo ou classe | Gestão | E | MVP |
| GET | `/v1/planos-manutencao/{id}` | Detalha plano | Leitura | E | MVP |
| PATCH | `/v1/planos-manutencao/{id}` | Altera plano | Gestão | E | MVP |
| DELETE | `/v1/planos-manutencao/{id}` | Exclui rascunho ou inativa plano | Gestão | E | MVP |
| POST | `/v1/planos-manutencao/{id}/ativar` | Passa a gerar OS | Gestão | E | MVP |
| POST | `/v1/planos-manutencao/{id}/suspender` | Para de gerar OS, com motivo | Gestão | E | MVP |
| GET | `/v1/planos-manutencao/{id}/gatilhos` | Lista gatilhos | Leitura | E | MVP |
| POST | `/v1/planos-manutencao/{id}/gatilhos` | Adiciona gatilho (tempo, medidor, o que vencer primeiro) | Gestão | E | MVP |
| PATCH | `/v1/planos-manutencao/{id}/gatilhos/{gatilho_id}` | Altera gatilho | Gestão | E | MVP |
| DELETE | `/v1/planos-manutencao/{id}/gatilhos/{gatilho_id}` | Remove gatilho | Gestão | E | MVP |
| GET | `/v1/planos-manutencao/{id}/previsao` | Próximos vencimentos projetados pelo uso médio | Leitura | E | MVP |
| POST | `/v1/planos-manutencao/{id}/gerar-os` | Gera OS antecipada manualmente | Gestão | E | MVP |
| GET | `/v1/planos-manutencao/vencimentos` | Calendário consolidado de vencimentos | Leitura | E | MVP |
| GET | `/v1/procedimentos` | Lista procedimentos | Leitura | E | MVP |
| POST | `/v1/procedimentos` | Cria procedimento (tarefas, peças, ferramentas) | Gestão | E | MVP |
| GET | `/v1/procedimentos/{id}` | Detalha procedimento | Leitura | E | MVP |
| PATCH | `/v1/procedimentos/{id}` | Altera procedimento | Gestão | E | MVP |
| DELETE | `/v1/procedimentos/{id}` | Inativa procedimento | Gestão | E | MVP |

### 6.7 CHK — Checklists

| Método | Rota | Descrição | Perfil | Plano | Fase |
|---|---|---|---|---|---|
| GET | `/v1/checklists-modelo` | Lista modelos (filtro por classe de ativo) | Todos, `prestador` | E | MVP |
| POST | `/v1/checklists-modelo` | Cria modelo em rascunho | Gestão | E | MVP |
| GET | `/v1/checklists-modelo/{id}` | Detalha modelo e itens | Todos, `prestador` | E | MVP |
| PATCH | `/v1/checklists-modelo/{id}` | Altera rascunho | Gestão | E | MVP |
| DELETE | `/v1/checklists-modelo/{id}` | Inativa modelo | Gestão | E | MVP |
| POST | `/v1/checklists-modelo/{id}/publicar` | Publica nova versão imutável | Gestão | E | MVP |
| GET | `/v1/checklists-modelo/{id}/versoes` | Versões publicadas | Leitura, `prestador` | E | MVP |
| GET | `/v1/checklists-execucao` | Lista execuções (`ativo_id`, `resultado`, período) | Leitura, `prestador` | E | MVP |
| POST | `/v1/checklists-execucao` | Registra execução completa (seção 10) | Execução, `motorista`, `prestador` | E | MVP |
| GET | `/v1/checklists-execucao/{id}` | Detalha execução | Leitura, autor, `prestador` | E | MVP |
| GET | `/v1/checklists-execucao/{id}/relatorio` | Relatório em PDF com fotos e assinatura | Leitura, autor, `prestador` | E | MVP |

### 6.8 EST — Estoque e compras

| Método | Rota | Descrição | Perfil | Plano | Fase |
|---|---|---|---|---|---|
| GET | `/v1/itens-estoque` | Lista itens (`controla_serie`, `reparavel`) | Leitura, `prestador` | E | MVP |
| POST | `/v1/itens-estoque` | Cria item; `controla_serie` e `reparavel` exigem o Profissional | Estoque | E | MVP |
| GET | `/v1/itens-estoque/{id}` | Detalha item com saldos | Leitura, `prestador` | E | MVP |
| PATCH | `/v1/itens-estoque/{id}` | Altera item (mínimo, máximo, ponto de pedido) | Estoque | E | MVP |
| DELETE | `/v1/itens-estoque/{id}` | Inativa item | Estoque | E | MVP |
| GET | `/v1/almoxarifados` | Lista almoxarifados (`tipo`, `fornecedor_id`, `unidade_id`) | Leitura, `prestador` | E | MVP |
| POST | `/v1/almoxarifados` | Cria almoxarifado ligado a uma unidade; `tipo` `proprio`, `em_terceiro` (com `fornecedor_id`) ou `quarentena`. No Essencial, um `proprio` por unidade | Estoque | E | MVP |
| GET | `/v1/almoxarifados/{id}` | Detalha almoxarifado | Leitura, `prestador` | E | MVP |
| PATCH | `/v1/almoxarifados/{id}` | Altera almoxarifado | Estoque | E | MVP |
| DELETE | `/v1/almoxarifados/{id}` | Inativa almoxarifado sem saldo | Estoque | E | MVP |
| GET | `/v1/saldos-estoque` | Saldo físico, reservado e disponível por item e almoxarifado (`tipo_almoxarifado`, `fornecedor_id`) | Leitura, `prestador` | E | MVP |
| GET | `/v1/saldos-estoque/consolidado` | Estoque geral por item (soma de todos os almoxarifados) com abertura por tipo de almoxarifado e por subcontratada | Leitura | P | F2 |
| GET | `/v1/movimentacoes-estoque` | Lista movimentações (`tipo`, `item_estoque_id`, `almoxarifado_id`, `fornecedor_id`, `ordem_reparo_id`, `nfe_chave`, período) | Leitura, `prestador` | E | MVP |
| POST | `/v1/movimentacoes-estoque` | Registra movimentação pelo `tipo` (20.1): `entrada`, `ajuste`, `transferencia`, `remessa_terceiro`, `retorno_terceiro`, `baixa_sucata` (este exige `est.decidir_destino`). `saida_os` e `devolucao_os` nascem em `/ordens-servico/{id}/pecas`; `envio_conserto` e `retorno_conserto`, em `/ordens-reparo` | Estoque, `prestador` (4.9) | E | MVP |
| GET | `/v1/movimentacoes-estoque/{id}` | Detalha movimentação | Leitura, `prestador` | E | MVP |
| POST | `/v1/movimentacoes-estoque/{id}/estornar` | Estorna movimentação com motivo | Estoque | E | MVP |
| POST | `/v1/movimentacoes-estoque/{id}/confirmar-recebimento` | Subcontratada confirma o recebimento de uma remessa (20.2) | Estoque, `prestador` | P | F2 |
| GET | `/v1/reservas-estoque` | Lista reservas | Leitura | E | MVP |
| POST | `/v1/reservas-estoque` | Reserva item para OS | Execução, Estoque | E | MVP |
| POST | `/v1/reservas-estoque/{id}/cancelar` | Cancela reserva | Execução, Estoque | E | MVP |
| GET | `/v1/inventarios-estoque` | Lista inventários | Leitura | P | F2 |
| POST | `/v1/inventarios-estoque` | Abre inventário de um almoxarifado (congela o saldo de referência) | Estoque | P | F2 |
| GET | `/v1/inventarios-estoque/{id}` | Detalha inventário com contagens e divergências | Leitura | P | F2 |
| POST | `/v1/inventarios-estoque/{id}/contagens` | Registra contagens por item (e número de série) | Estoque, `prestador` | P | F2 |
| POST | `/v1/inventarios-estoque/{id}/concluir` | Conclui e gera ajustes com motivo | Estoque | P | F2 |
| GET | `/v1/pecas-serializadas` | Lista peças (`item_estoque_id`, `numero_serie`, `situacao`, `almoxarifado_id`, `ativo_id`, `recondicionada`) | Leitura, `prestador` | P | F2 |
| POST | `/v1/pecas-serializadas` | Cadastra peça com número de série (saldo inicial ou entrada sem NF-e) | Estoque | P | F2 |
| GET | `/v1/pecas-serializadas/{id}` | Detalha peça com garantia vigente | Leitura, `prestador` | P | F2 |
| PATCH | `/v1/pecas-serializadas/{id}` | Altera dados cadastrais (garantia, observação) | Estoque | P | F2 |
| GET | `/v1/pecas-serializadas/{id}/historico` | Instalações, remoções, movimentações e reparos | Leitura | P | F2 |
| POST | `/v1/pecas-serializadas/{id}/instalar` | Instala em ativo, ligada a uma OS (`em_estoque` → `instalada`) | Execução, `prestador` | P | F2 |
| POST | `/v1/pecas-serializadas/{id}/remover` | Remove do ativo para um almoxarifado, `quarentena` quando defeituosa (`instalada` → `em_estoque`) | Execução, `prestador` | P | F2 |
| POST | `/v1/pecas-serializadas/{id}/baixar` | Baixa como sucata, com motivo (`baixada`; gera `baixa_sucata`) | Destino de peça | P | F2 |
| GET | `/v1/ordens-reparo` | Lista ordens de reparo (`status`, `fornecedor_id`, `item_estoque_id`, prazo) | Leitura, `prestador` | P | F2 |
| POST | `/v1/ordens-reparo` | Abre ordem de reparo com peças da `quarentena` (20.3) | Destino de peça | P | F2 |
| GET | `/v1/ordens-reparo/{id}` | Detalha ordem de reparo | Leitura, `prestador` | P | F2 |
| PATCH | `/v1/ordens-reparo/{id}` | Altera ordem `aberta` (fornecedor, itens, prazo) | Estoque | P | F2 |
| POST | `/v1/ordens-reparo/{id}/enviar` | `aberta` → `enviada`; gera `envio_conserto` com NF-e de remessa | Estoque | P | F2 |
| POST | `/v1/ordens-reparo/{id}/registrar-orcamento` | `enviada` → `em_orcamento`; cria `orcamento_fornecedor` (permissão `ter.registrar_orcamento`, a mesma das rotas de orçamento de fornecedor) | Gestão, Compras, `prestador` | P | F2 |
| POST | `/v1/ordens-reparo/{id}/aprovar` | `em_orcamento` → `aprovada` | Gestão, `comprador` | P | F2 |
| POST | `/v1/ordens-reparo/{id}/registrar-retorno` | `aprovada` → `retornada`; gera `retorno_conserto`; peça volta recondicionada, com custo e garantia | Estoque | P | F2 |
| POST | `/v1/ordens-reparo/{id}/declarar-irreparavel` | `enviada`, `em_orcamento` ou `aprovada` → `irreparavel`; peça vai para sucata ou volta sem reparo | Destino de peça | P | F2 |
| POST | `/v1/ordens-reparo/{id}/cancelar` | `aberta` → `cancelada`, com motivo | Estoque, Gestão | P | F2 |
| GET | `/v1/requisicoes-compra` | Lista requisições | Leitura | P | F2 |
| POST | `/v1/requisicoes-compra` | Cria requisição de compra | Execução, Compras | P | F2 |
| POST | `/v1/requisicoes-compra/{id}/aprovar` | Aprova requisição | Gestão | P | F2 |
| POST | `/v1/requisicoes-compra/{id}/atender` | Atende com saída de estoque ou gera pedido de compra | Compras | P | F2 |
| GET | `/v1/pedidos-compra` | Lista pedidos de compra | Leitura | P | F2 |
| POST | `/v1/pedidos-compra` | Cria pedido de compra | Compras, Gestão | P | F2 |
| GET | `/v1/pedidos-compra/{id}` | Detalha pedido | Leitura | P | F2 |
| PATCH | `/v1/pedidos-compra/{id}` | Altera pedido em rascunho | Compras, Gestão | P | F2 |
| POST | `/v1/pedidos-compra/{id}/aprovar` | Aprova pedido | Gestão | P | F2 |
| POST | `/v1/pedidos-compra/{id}/enviar` | Marca como enviado ao fornecedor | Compras | P | F2 |
| POST | `/v1/pedidos-compra/{id}/cancelar` | Cancela pedido com motivo | Gestão | P | F2 |
| GET | `/v1/notas-fiscais-entrada` | Lista NF-e importadas | Leitura | P | F2 |
| POST | `/v1/notas-fiscais-entrada` | Importa XML de NF-e (14.4) | Compras, Integração | P | F2 |
| GET | `/v1/notas-fiscais-entrada/{id}` | Detalha NF-e, itens e conciliação | Leitura | P | F2 |
| POST | `/v1/notas-fiscais-entrada/{id}/conciliar` | Vincula itens da nota a itens de estoque, pedido, OS e ativo | Compras | P | F2 |
| POST | `/v1/notas-fiscais-entrada/{id}/confirmar-entrada` | Gera entradas de estoque e custos | Estoque | P | F2 |

Planos, conforme a seção 2.1 do sumário:

- **Estoque básico (Essencial):** itens, um almoxarifado `proprio` por unidade, saldo, reserva e baixa em OS (`saida_os`, `devolucao_os`), estoque mínimo, `entrada` e `ajuste` com motivo.
- **Estoque e compras completos (Profissional):** vários almoxarifados por unidade, `transferencia`, inventário, requisição, pedido de compra, entrada por NF-e.
- **Estoque em terceiros (Profissional, Fase 2):** almoxarifado `em_terceiro`, saldo geral e por subcontratada, `remessa_terceiro` e `retorno_terceiro`.
- **Peças serializadas e reparáveis (Profissional, Fase 2):** `controla_serie`, `reparavel`, almoxarifado `quarentena`, ordem de reparo, `envio_conserto`, `retorno_conserto`, `baixa_sucata`.

Tipo de movimentação ou de almoxarifado fora do plano devolve `403 plano_nao_inclui_recurso` com `plano_necessario`. Fornecedores saíram deste módulo e estão em TER (6.9).

### 6.9 TER — Terceiros

| Método | Rota | Descrição | Perfil | Plano | Fase |
|---|---|---|---|---|---|
| GET | `/v1/fornecedores` | Lista fornecedores (`papel`, `cnpj`, `municipio`, `uf`, `busca`) | Leitura, `prestador` | E | MVP |
| POST | `/v1/fornecedores` | Cria fornecedor com `papeis` (`fornecedor_pecas`, `subcontratada`, `assistencia_tecnica`, `oficina`, `locador`, `posto_combustivel`, `seguradora`) | Compras, Gestão | E | MVP |
| GET | `/v1/fornecedores/{id}` | Detalha fornecedor | Leitura, `prestador` | E | MVP |
| PATCH | `/v1/fornecedores/{id}` | Altera fornecedor e papéis | Compras, Gestão | E | MVP |
| DELETE | `/v1/fornecedores/{id}` | Inativa fornecedor; revoga links abertos e desativa usuários `prestador` ligados a ele | Compras, Gestão | E | MVP |
| GET | `/v1/servicos-fornecedor` | Catálogo de serviços; busca de quem atende (`fornecedor_id`, `classe_ativo_id`, `marca`, `municipio`, `uf`, `unidade_preco`) | Leitura, `prestador` | E | MVP |
| POST | `/v1/servicos-fornecedor` | Cadastra serviço do fornecedor | Compras, Gestão, `prestador` (4.9) | E | MVP |
| GET | `/v1/servicos-fornecedor/{id}` | Detalha serviço | Leitura, `prestador` | E | MVP |
| PATCH | `/v1/servicos-fornecedor/{id}` | Altera serviço (prazos, preço de referência, área de atendimento) | Compras, Gestão, `prestador` (4.9) | E | MVP |
| DELETE | `/v1/servicos-fornecedor/{id}` | Inativa serviço | Compras, Gestão | E | MVP |
| POST | `/v1/ordens-servico/{id}/atribuir-fornecedor` | Atribui a OS a um prestador: `executor_tipo = fornecedor` (19.3). Fornecedor com conexão `ativa`: cria o vínculo e a OS na conta da subcontratada (22.8) | Gestão | E | MVP |
| POST | `/v1/ordens-servico/{id}/retirar-fornecedor` | Volta a OS para `executor_tipo = interno`; revoga links abertos; em OS vinculada, cancela a cópia na subcontratada e encerra o vínculo | Gestão | E | MVP |
| GET | `/v1/links-execucao` | Lista links de uso único (`ordem_servico_id`, `fornecedor_id`, `situacao`) | Gestão | E | MVP |
| POST | `/v1/links-execucao` | Gera link de uso único para OS atribuída a fornecedor (19.4) | Gestão | E | MVP |
| GET | `/v1/links-execucao/{id}` | Detalha link (nunca devolve o token) | Gestão | E | MVP |
| POST | `/v1/links-execucao/{id}/revogar` | Revoga link não utilizado | Gestão | E | MVP |
| POST | `/v1/publico/links-execucao/ativar` | Recebe o token **no corpo** e devolve a sessão curta do link (19.4) | Público (token no corpo) | E | MVP |
| GET | `/v1/publico/links-execucao/atual` | Dados mínimos da OS do link da sessão | Link | E | MVP |
| POST | `/v1/publico/links-execucao/atual/anexos` | Registra foto ou assinatura e devolve URL de upload restrita ao link | Link | E | MVP |
| POST | `/v1/publico/links-execucao/atual/registrar-execucao` | Registra a execução e consome o link | Link | E | MVP |
| GET | `/v1/contratos-prestacao` | Lista contratos de prestação (`fornecedor_id`, `vigente_em`, `ativo_id`, `classe_ativo_id`) | Gestão, Compras, `prestador` | P | F2 |
| POST | `/v1/contratos-prestacao` | Cria contrato (vigência, escopo, ativos ou classes cobertos, SLA de atendimento e de solução, preços) | Gestão, Compras | P | F2 |
| GET | `/v1/contratos-prestacao/{id}` | Detalha contrato | Gestão, Compras, `prestador` | P | F2 |
| PATCH | `/v1/contratos-prestacao/{id}` | Altera contrato | Gestão, Compras | P | F2 |
| DELETE | `/v1/contratos-prestacao/{id}` | Inativa contrato | Gestão, Compras | P | F2 |
| GET | `/v1/contratos-prestacao/{id}/sla` | Cumprimento do SLA de atendimento e de solução nas OS cobertas | Gestão, Compras, `prestador` | P | F2 |
| GET | `/v1/avaliacoes-fornecedor` | Lista avaliações (`fornecedor_id`, período) | Gestão, Compras, `prestador` | P | F2 |
| POST | `/v1/avaliacoes-fornecedor` | Registra avaliação (nota, comentário); SLA cumprido e retrabalho são calculados | Gestão | P | F2 |
| GET | `/v1/avaliacoes-fornecedor/{id}` | Detalha avaliação | Gestão, Compras, `prestador` | P | F2 |
| GET | `/v1/orcamentos-fornecedor` | Lista orçamentos (`ordem_servico_id`, `ordem_reparo_id`, `fornecedor_id`, `status`) | Gestão, Compras, `prestador` | P | F2 |
| POST | `/v1/orcamentos-fornecedor` | Registra orçamento de oficina, assistência ou conserto | Gestão, Compras, `prestador` | P | F2 |
| GET | `/v1/orcamentos-fornecedor/{id}` | Detalha orçamento | Gestão, Compras, `prestador` | P | F2 |
| PATCH | `/v1/orcamentos-fornecedor/{id}` | Altera orçamento pendente | Gestão, Compras, `prestador` | P | F2 |
| POST | `/v1/orcamentos-fornecedor/{id}/aprovar` | Aprova orçamento; o valor vira custo previsto de terceiros na OS | Gestão, `comprador` | P | F2 |
| POST | `/v1/orcamentos-fornecedor/{id}/recusar` | Recusa orçamento com motivo | Gestão, `comprador` | P | F2 |

TER básico (fornecedores com papéis, catálogo de serviços, OS executada por prestador, perfil `prestador`, link de uso único) está em todos os planos desde o MVP. Contratos de prestação com SLA, avaliação e oficina externa com orçamento são do Profissional, na Fase 2. A oficina externa deixou de ser sub-recurso de FRO: é uma OS atribuída a fornecedor com papel `oficina`, com orçamento em `/v1/orcamentos-fornecedor`.

Com o SRV (seção 22), um fornecedor pode ter conta própria na Tenaz e se conectar à conta do cliente. Nesse caso `atribuir-fornecedor` não exige usuário `prestador`: a OS vai para a conta da subcontratada pela conexão (22.8).

### 6.10 SRV — Gestão de serviços do prestador

Rotas próprias do módulo. Na conta da subcontratada, "Leitura", "Gestão" e "Execução" são os grupos da seção 4.8 aplicados aos usuários **dela**. As rotas de conexão existem nos dois lados: o cliente convida e administra o compartilhamento; a subcontratada aceita ou recusa; qualquer lado suspende, reativa (quem suspendeu) ou encerra. Detalhamento na seção 22.

| Método | Rota | Descrição | Perfil | Plano | Fase |
|---|---|---|---|---|---|
| GET | `/v1/especialidades-tecnico` | Catálogo de especialidades da conta (`busca`, `ativo`) | Leitura | E, S | MVP |
| POST | `/v1/especialidades-tecnico` | Cria especialidade (nome, descrição, classes de ativo relacionadas) | Gestão | E, S | MVP |
| GET | `/v1/especialidades-tecnico/{id}` | Detalha especialidade, com a contagem de técnicos | Leitura | E, S | MVP |
| PATCH | `/v1/especialidades-tecnico/{id}` | Altera especialidade | Gestão | E, S | MVP |
| DELETE | `/v1/especialidades-tecnico/{id}` | Inativa especialidade; as ligações `pessoa_especialidade` ficam no histórico | Gestão | E, S | MVP |
| GET | `/v1/clientes-servico` | Lista clientes atendidos (`tipo_pessoa`, `documento`, `municipio`, `uf`, `conectado`, `busca`) | Leitura | S, P | MVP |
| POST | `/v1/clientes-servico` | Cadastra cliente atendido (pessoa jurídica ou física, CNPJ ou CPF, contatos, endereços de atendimento) | Gestão | S, P | MVP |
| GET | `/v1/clientes-servico/{id}` | Detalha cliente, com resumo de equipamentos, OS abertas e conexão | Leitura | S, P | MVP |
| PATCH | `/v1/clientes-servico/{id}` | Altera cliente; os dados vindos da conexão (razão social, CNPJ) ficam só leitura | Gestão | S, P | MVP |
| DELETE | `/v1/clientes-servico/{id}` | Inativa cliente sem OS em aberto | Gestão | S, P | MVP |
| GET | `/v1/conexoes-empresa` | Lista as conexões em que a conta é parte, com o papel dela (`cliente` ou `prestador`) e `estado` | Gestão | E, S | MVP |
| POST | `/v1/conexoes-empresa` | Ação `convidar`: o cliente convida a conta do fornecedor; cria a conexão em `convidada` e envia o código por e-mail (22.5) | Conexão (cliente) | E | MVP |
| GET | `/v1/conexoes-empresa/{id}` | Detalha a conexão (`campos_compartilhados`, ativos compartilhados, datas) | Gestão | E, S | MVP |
| PATCH | `/v1/conexoes-empresa/{id}` | O cliente altera os ativos compartilhados e reduz `campos_compartilhados`; ampliar campos exige novo aceite | Conexão (cliente) | E | MVP |
| POST | `/v1/conexoes-empresa/{id}/aceitar` | `convidada` → `ativa`, com o código do convite no corpo; liga a um `cliente_servico` ou cria um | Conexão (prestador) | S | MVP |
| POST | `/v1/conexoes-empresa/{id}/recusar` | `convidada` → `encerrada`, com motivo | Conexão (prestador) | S | MVP |
| POST | `/v1/conexoes-empresa/{id}/suspender` | `ativa` → `suspensa`, com motivo; a sincronização para e fica pendente | Conexão (cliente ou prestador) | E, S | MVP |
| POST | `/v1/conexoes-empresa/{id}/reativar` | `suspensa` → `ativa`, pelo lado que suspendeu; as pendências são sincronizadas | Conexão (quem suspendeu) | E, S | MVP |
| POST | `/v1/conexoes-empresa/{id}/encerrar` | `convidada`, `ativa` ou `suspensa` → `encerrada` (final), com motivo; para a sincronização sem apagar histórico | Conexão (cliente ou prestador) | E, S | MVP |
| GET | `/v1/conexoes-empresa/{id}/ativos` | Ativos compartilhados pela conexão: no cliente, os ativos originais; na subcontratada, as cópias `de_cliente` | Gestão | E, S | MVP |
| GET | `/v1/vinculos-os` | Vínculos de OS da conta (`conexao_empresa_id`, `lado_origem`, `estado_sincronizacao`, período) | Gestão | E, S | MVP |
| GET | `/v1/vinculos-os/{id}` | Detalha o vínculo, com o histórico de sincronizações e o último erro | Gestão | E, S | MVP |
| POST | `/v1/vinculos-os/{id}/reprocessar` | Reenfileira a sincronização de um vínculo em `erro`; exige conexão `ativa` (`409 conexao_inativa`) | Gestão | E, S | MVP |

**Clientes atendidos fora do SRV.** `cliente_servico` também existe nas contas com CTR (sumário 3.1): numa conta do plano Profissional, `/v1/clientes-servico` fica disponível a partir da Fase 2, junto com o CTR, para cadastrar o cliente do contrato (por exemplo, a prefeitura). Daí o plano `S, P` nas cinco rotas acima; o `P` vale só a partir da Fase 2.

Rotas de outros módulos usadas pelo SRV, contadas nos seus módulos: `/v1/pessoas` e `/v1/usuarios` (técnicos), `/v1/ativos` com `propriedade = de_cliente` (equipamentos de clientes), `/v1/ordens-servico` com `assumir`, `servicos`, `pecas`, `apontamentos`, `assinaturas` e `relatorio-pdf`, `/v1/servicos-fornecedor` (catálogo próprio da subcontratada, 22.7), `POST /v1/relatorios/servicos_por_cliente/gerar` (relatório de serviços por cliente) e a sincronização do app (seção 12).

### 6.11 LOC — Locação

| Método | Rota | Descrição | Perfil | Plano | Fase |
|---|---|---|---|---|---|
| GET | `/v1/locatarios` | Lista locatários (`cnpj`, `cpf`, `busca`) | Gestão | P | F2 |
| POST | `/v1/locatarios` | Cadastra locatário (empresa ou pessoa, CNPJ ou CPF, contatos, locais de uso) | Gestão | P | F2 |
| GET | `/v1/locatarios/{id}` | Detalha locatário | Gestão, `locatario` (o próprio) | P | F2 |
| PATCH | `/v1/locatarios/{id}` | Altera locatário | Gestão | P | F2 |
| DELETE | `/v1/locatarios/{id}` | Inativa locatário sem contrato vigente; desativa usuários `locatario` ligados a ele | Gestão | P | F2 |
| GET | `/v1/contratos-locacao` | Lista contratos (`sentido`, `status`, `locatario_id`, `fornecedor_id`, `ativo_id`, `vigente_em`) | Gestão, `locatario` | P | F2 |
| POST | `/v1/contratos-locacao` | Cria contrato em `rascunho`, `cedido` ou `tomado` (21.2) | Gestão | P | F2 |
| GET | `/v1/contratos-locacao/{id}` | Detalha contrato com itens e matriz | Gestão, `locatario` | P | F2 |
| PATCH | `/v1/contratos-locacao/{id}` | Altera contrato; em `vigente`, só data de fim, `franquia_horas`, `franquia_km` e contatos | Gestão | P | F2 |
| POST | `/v1/contratos-locacao/{id}/ativar` | `rascunho` → `vigente` | Gestão | P | F2 |
| POST | `/v1/contratos-locacao/{id}/encerrar` | `vigente` → `encerrado` (todos os itens devolvidos e vistoriados); antes de `data_fim`, exige `motivo_encerramento` | Gestão | P | F2 |
| POST | `/v1/contratos-locacao/{id}/cancelar` | Só `rascunho` → `cancelado`, com motivo; contrato `vigente` termina por `encerrar` | Gestão | P | F2 |
| GET | `/v1/contratos-locacao/{id}/itens` | Ativos do contrato com datas de entrega e devolução | Gestão, `locatario` | P | F2 |
| POST | `/v1/contratos-locacao/{id}/itens` | Inclui ativo (em `rascunho`, ou em `vigente` como aditivo); ativo `de_terceiro` em contrato `cedido` devolve `422 sublocacao_nao_permitida` | Gestão | P | F2 |
| PATCH | `/v1/contratos-locacao/{id}/itens/{item_id}` | Altera datas previstas do item | Gestão | P | F2 |
| DELETE | `/v1/contratos-locacao/{id}/itens/{item_id}` | Remove item ainda não entregue | Gestão | P | F2 |
| POST | `/v1/contratos-locacao/{id}/itens/{item_id}/substituir` | Registra a substituição do ativo pelo locador (contrato `tomado`), com medição do SLA de substituição | Gestão | P | F2 |
| GET | `/v1/contratos-locacao/{id}/responsabilidades` | Matriz de responsabilidade da manutenção | Gestão, `locatario` | P | F2 |
| PUT | `/v1/contratos-locacao/{id}/responsabilidades` | Substitui a matriz por inteiro | Gestão | P | F2 |
| GET | `/v1/contratos-locacao/{id}/cobranca` | Relatório de cobrança: avarias com custo do locatário e horas ou km excedentes à franquia | Gestão, `locatario` | P | F2 |
| GET | `/v1/vistorias-locacao` | Lista vistorias (`contrato_locacao_id`, `ativo_id`, `tipo`) | Gestão, Execução, `locatario` | P | F2 |
| POST | `/v1/vistorias-locacao` | Registra vistoria de entrega ou de devolução (21.4) | Gestão, Execução | P | F2 |
| GET | `/v1/vistorias-locacao/{id}` | Detalha vistoria | Gestão, Execução, `locatario` | P | F2 |
| POST | `/v1/vistorias-locacao/{id}/assinaturas` | Registra a assinatura de uma das partes | Gestão, Execução, `locatario` | P | F2 |
| GET | `/v1/vistorias-locacao/{id}/relatorio` | Termo de vistoria em PDF com fotos, leituras e assinaturas | Gestão, Execução, `locatario` | P | F2 |

Locação é do plano Profissional, na Fase 2. O perfil `locatario` é gratuito (sumário 2.1) e acessa só o portal web do locatário (sumário 3.6), com o alcance da seção 4.9. O relatório consolidado de cobrança de todos os contratos sai por `POST /v1/relatorios/locacao_cobranca/gerar`; a fatura é emitida pelo ERP.

### 6.12 FRO — Frota

| Método | Rota | Descrição | Perfil | Plano | Fase |
|---|---|---|---|---|---|
| GET | `/v1/ativos/{id}/dados-veiculo` | Placa, RENAVAM, chassi/VIN, marca, modelo, ano, capacidade do tanque, eixos | Leitura, `motorista` | E | MVP |
| PATCH | `/v1/ativos/{id}/dados-veiculo` | Altera dados do veículo | Gestão | E | MVP |
| POST | `/v1/ativos/{id}/acoplar` | Acopla implemento ao veículo (cavalo + carreta) (proposta) | Gestão, `motorista` | E | MVP |
| POST | `/v1/ativos/{id}/desacoplar` | Desfaz a composição (proposta) | Gestão, `motorista` | E | MVP |
| GET | `/v1/pneus` | Lista pneus (`situacao`, `ativo_id`, `vida`) | Leitura | P | F2 |
| POST | `/v1/pneus` | Cadastra pneu pelo número de fogo | Gestão, Estoque | P | F2 |
| GET | `/v1/pneus/{id}` | Detalha pneu, vida, sulco e CPK | Leitura | P | F2 |
| PATCH | `/v1/pneus/{id}` | Altera dados cadastrais | Gestão | P | F2 |
| POST | `/v1/pneus/{id}/montar` | Monta em veículo e posição, com leitura de km | Execução | P | F2 |
| POST | `/v1/pneus/{id}/desmontar` | Desmonta, com destino (estoque, recapagem, conserto, sucata) | Execução | P | F2 |
| POST | `/v1/pneus/{id}/afericoes` | Registra sulco por ponto e pressão | Execução, `motorista` | P | F2 |
| POST | `/v1/pneus/{id}/enviar-recapagem` | Envia para recapagem (fornecedor, OS) | Gestão, Estoque | P | F2 |
| POST | `/v1/pneus/{id}/receber-recapagem` | Recebe com nova vida (R1, R2...) e custo | Gestão, Estoque | P | F2 |
| POST | `/v1/pneus/{id}/sucatear` | Baixa definitiva com motivo | Gestão | P | F2 |
| GET | `/v1/pneus/{id}/historico` | Montagens, aferições, recapagens e custos | Leitura | P | F2 |
| GET | `/v1/ativos/{id}/mapa-pneus` | Posições por eixo com pneu montado e último sulco | Leitura, `motorista` | P | F2 |
| GET | `/v1/abastecimentos` | Lista abastecimentos (`ativo_id`, `origem`, `anomalia`) | Leitura | P | F2 |
| POST | `/v1/abastecimentos` | Registra abastecimento ou Arla manual (com hodômetro); `fornecedor_id` opcional, com papel `posto_combustivel` | Execução, `motorista` | P | F2 |
| GET | `/v1/abastecimentos/{id}` | Detalha abastecimento | Leitura | P | F2 |
| PATCH | `/v1/abastecimentos/{id}` | Corrige abastecimento | Gestão | P | F2 |
| DELETE | `/v1/abastecimentos/{id}` | Exclui abastecimento manual com motivo | Gestão | P | F2 |
| GET | `/v1/multas` | Lista multas (`prazo_indicacao_ate`, `situacao`) | Gestão | P | F2 |
| POST | `/v1/multas` | Registra multa | Gestão | P | F2 |
| GET | `/v1/multas/{id}` | Detalha multa | Gestão, condutor | P | F2 |
| PATCH | `/v1/multas/{id}` | Altera multa | Gestão | P | F2 |
| POST | `/v1/multas/{id}/indicar-condutor` | Registra a indicação feita no Portal Senatran (condutor, protocolo, data) | Gestão | P | F2 |
| POST | `/v1/multas/{id}/registrar-pagamento` | Registra pagamento (com ou sem desconto SNE) | Gestão | P | F2 |
| GET | `/v1/sinistros` | Lista sinistros | Gestão | P | F2 |
| POST | `/v1/sinistros` | Registra sinistro (com fotos); `seguradora_id` opcional, fornecedor com papel `seguradora` | Gestão, `motorista` | P | F2 |
| GET | `/v1/sinistros/{id}` | Detalha sinistro | Gestão, autor | P | F2 |
| PATCH | `/v1/sinistros/{id}` | Altera sinistro | Gestão | P | F2 |
| POST | `/v1/sinistros/{id}/encerrar` | Encerra com custos e ressarcimentos | Gestão | P | F2 |

A oficina externa saiu de FRO: veículo enviado a oficina é uma OS atribuída a fornecedor com papel `oficina` (6.9), com orçamento em `/v1/orcamentos-fornecedor`.

### 6.13 CON — Conformidade

| Método | Rota | Descrição | Perfil | Plano | Fase |
|---|---|---|---|---|---|
| GET | `/v1/tipos-documento` | Catálogo de tipos (CRLV-e, cronotacógrafo, CIV, CIPP, laudo NR-13, seguro, CNH, MOPP...), com `critico_despacho` | Leitura | E | MVP |
| POST | `/v1/tipos-documento` | Cria tipo personalizado | Admin | E | MVP |
| PATCH | `/v1/tipos-documento/{id}` | Altera tipo (antecedência de alertas, criticidade para despacho) | Admin | E | MVP |
| GET | `/v1/documentos-vencimento` | Agenda de documentos (`situacao`: `vigente`, `vencendo`, `vencido`; `vence_ate`; dono: `ativo_id`, `pessoa_id` ou empresa) | Leitura | E | MVP |
| POST | `/v1/documentos-vencimento` | Registra documento ligado a ativo, pessoa ou empresa | Gestão | E | MVP |
| GET | `/v1/documentos-vencimento/{id}` | Detalha documento | Leitura | E | MVP |
| PATCH | `/v1/documentos-vencimento/{id}` | Altera documento | Gestão | E | MVP |
| DELETE | `/v1/documentos-vencimento/{id}` | Inativa documento | Gestão | E | MVP |
| POST | `/v1/documentos-vencimento/{id}/renovar` | Registra o documento novo e arquiva o anterior | Gestão | E | MVP |
| GET | `/v1/inventario-nr12` | Inventário atualizado das máquinas NR-12 (item 12.18.1), JSON ou PDF | Leitura | E | MVP |
| GET | `/v1/cronograma-nr12` | Cronograma das preventivas ligadas à segurança (item 12.11.2.2) | Leitura | E | MVP |
| GET | `/v1/bloqueios-ativo` | Bloqueios de conformidade (`bloqueio_ativo`) vigentes e históricos (`ativo_id`, `origem`, `regra_automacao_id`, `vigente`) | Leitura | E | MVP |
| POST | `/v1/bloqueios-ativo` | Cria bloqueio manual (`origem = manual`); o ativo passa a `parado` | Gestão | E | MVP |
| POST | `/v1/bloqueios-ativo/{id}/liberar` | Libera bloqueio de origem `manual` ou `regra`, com motivo; o status do ativo é recalculado | Gestão | E | MVP |
| GET | `/v1/bloqueios-despacho/consulta` | Verificação calculada "apto para despacho" por placa ou VIN, para TMS (seção 16); não é entidade | Leitura, Integração (`despacho:consulta`) | P | MVP |
| POST | `/v1/bloqueios-despacho/consultas` | Verificação em lote, até 100 veículos | Integração (`despacho:consulta`) | P | MVP |
| GET | `/v1/capacitacoes` | Lista capacitações de pessoas (`pessoa_id`; NR-10, NR-11, NR-12, NR-35, MOPP) | Gestão | P | F2 |
| POST | `/v1/capacitacoes` | Registra capacitação com validade | Gestão | P | F2 |
| GET | `/v1/capacitacoes/{id}` | Detalha capacitação | Gestão, titular | P | F2 |
| PATCH | `/v1/capacitacoes/{id}` | Altera capacitação | Gestão | P | F2 |
| DELETE | `/v1/capacitacoes/{id}` | Inativa capacitação | Gestão | P | F2 |
| GET | `/v1/ativos/{id}/enquadramento-nr13` | Categoria, SPIE e prazos de inspeção do equipamento | Leitura | P | F2 |
| PUT | `/v1/ativos/{id}/enquadramento-nr13` | Define enquadramento; prazos menores do PLH são aceitos | Gestão, `plh` | P | F2 |
| GET | `/v1/prazos-nr13` | Calendário de inspeções NR-13 calculado | Leitura | P | F2 |
| GET | `/v1/ordens-servico/{id}/bloqueios-energia` | Registros de bloqueio e etiquetagem da OS | Leitura, `prestador` | E | MVP |
| POST | `/v1/ordens-servico/{id}/bloqueios-energia` | Registra bloqueio e etiquetagem (fonte, horário, data, motivo, responsável; item 12.11.3); `ponto_bloqueio_id` opcional no Profissional | Execução, `prestador` | E | MVP |
| POST | `/v1/ordens-servico/{id}/bloqueios-energia/{bloqueio_id}/retirar` | Retira bloqueio | Execução, `prestador` | E | MVP |
| GET | `/v1/pontos-bloqueio` | Pontos de bloqueio LOTO por ativo, com procedimento | Leitura | P | F2 |
| POST | `/v1/pontos-bloqueio` | Cadastra ponto de bloqueio (tipo de energia, localização, dispositivo, procedimento, foto) | Gestão | P | F2 |
| PATCH | `/v1/pontos-bloqueio/{id}` | Altera ponto de bloqueio | Gestão | P | F2 |
| DELETE | `/v1/pontos-bloqueio/{id}` | Inativa ponto de bloqueio | Gestão | P | F2 |

**Origem do bloqueio de ativo** (sumário 3.6). `bloqueio_ativo.origem` é definida pelo servidor e não é aceita no corpo: `manual` (`POST /v1/bloqueios-ativo`), `conformidade` (gerado pelas regras padrão de conformidade, como o parecer NR-12) ou `regra` (efeito de uma regra configurável, 23.1). Campos só leitura acrescentados: `origem`, `regra_automacao_id` e `disparo_regra_id` (preenchidos com `origem = regra`), `levanta_com_tipo_documento_id` (tipo de documento que levanta o bloqueio sozinho, no efeito `bloquear_ativo_ate_documento`) e `levantado_automaticamente` (proposta). Bloqueio de origem `regra` também pode ser liberado à mão, com motivo, por quem tem o perfil Gestão (proposta). Bloqueio de origem `conformidade` segue a regra do seu documento (por exemplo, nova indicação `apto` no registro NR-12).

O registro de bloqueio e etiquetagem do item 12.11.3 da NR-12 é do Essencial, no MVP. O LOTO com pontos de bloqueio cadastrados e procedimento é do Profissional, na Fase 2 (sumário 2.1).

### 6.14 IND — Indicadores e relatórios

| Método | Rota | Descrição | Perfil | Plano | Fase |
|---|---|---|---|---|---|
| GET | `/v1/indicadores` | Catálogo de KPIs com fórmula, "como calculamos" e fonte de referência | Leitura | E | MVP |
| GET | `/v1/indicadores/valores` | Valores de vários KPIs com os mesmos filtros (seção 11) | Leitura | E | MVP |
| GET | `/v1/indicadores/{codigo}` | Valor, série temporal e agrupamento de um KPI | Leitura | E | MVP |
| GET | `/v1/indicadores/{codigo}/composicao` | Registros que entram no cálculo (detalhamento) | Leitura | E | MVP |
| GET | `/v1/paineis` | Lista painéis | Leitura | E | MVP |
| POST | `/v1/paineis` | Cria painel | Gestão | E | MVP |
| GET | `/v1/paineis/{id}` | Detalha painel | Leitura | E | MVP |
| PATCH | `/v1/paineis/{id}` | Altera painel | Gestão | E | MVP |
| DELETE | `/v1/paineis/{id}` | Exclui painel | Gestão | E | MVP |
| GET | `/v1/relatorios` | Catálogo de relatórios | Leitura | E | MVP |
| POST | `/v1/relatorios/{codigo}/gerar` | Gera relatório em PDF, XLSX ou CSV (assíncrono) | Leitura | E | MVP |
| POST | `/v1/exportacoes-lista` | Exporta uma lista com os mesmos filtros da tela (CSV ou XLSX) | Leitura | E | MVP |

### 6.15 INT — Integrações

| Método | Rota | Descrição | Perfil | Plano | Fase |
|---|---|---|---|---|---|
| GET | `/v1/tokens-api` | Lista tokens de API | Admin | P | MVP |
| POST | `/v1/tokens-api` | Cria token (valor exibido uma vez) | Admin | P | MVP |
| GET | `/v1/tokens-api/{id}` | Detalha token (sem o valor) | Admin | P | MVP |
| PATCH | `/v1/tokens-api/{id}` | Altera nome, escopos e IPs permitidos | Admin | P | MVP |
| DELETE | `/v1/tokens-api/{id}` | Revoga token | Admin | P | MVP |
| POST | `/v1/tokens-api/{id}/rotacionar` | Emite token novo; o antigo vale por até 24 h | Admin | P | MVP |
| GET | `/v1/tipos-evento-webhook` | Catálogo de eventos com esquema de `dados` | Admin, Integração | P | MVP |
| GET | `/v1/webhooks` | Lista assinaturas de webhook | Admin, Integração (`webhooks:gerenciar`) | P | MVP |
| POST | `/v1/webhooks` | Cria assinatura (URL, eventos, filtros) | Admin, Integração (`webhooks:gerenciar`) | P | MVP |
| GET | `/v1/webhooks/{id}` | Detalha assinatura e indicadores de entrega | Admin, Integração (`webhooks:gerenciar`) | P | MVP |
| PATCH | `/v1/webhooks/{id}` | Altera assinatura | Admin, Integração (`webhooks:gerenciar`) | P | MVP |
| DELETE | `/v1/webhooks/{id}` | Exclui assinatura | Admin, Integração (`webhooks:gerenciar`) | P | MVP |
| POST | `/v1/webhooks/{id}/testar` | Envia evento `webhook.teste` | Admin, Integração (`webhooks:gerenciar`) | P | MVP |
| POST | `/v1/webhooks/{id}/rotacionar-segredo` | Novo segredo; o antigo assina junto por 24 h | Admin | P | MVP |
| POST | `/v1/webhooks/{id}/reativar` | Reativa assinatura suspensa | Admin, Integração (`webhooks:gerenciar`) | P | MVP |
| GET | `/v1/webhooks/{id}/entregas` | Painel de entregas (`status`, inclusive `dlq`; período; evento) | Admin, Integração (`webhooks:gerenciar`) | P | MVP |
| GET | `/v1/entregas-webhook/{id}` | Entrega com corpo, tentativas e respostas | Admin, Integração (`webhooks:gerenciar`) | P | MVP |
| POST | `/v1/entregas-webhook/{id}/reentregar` | Reentrega manual de uma entrega | Admin, Integração (`webhooks:gerenciar`) | P | MVP |
| POST | `/v1/webhooks/{id}/reentregar-dlq` | Reentrega em lote as entregas da DLQ de um período | Admin | P | MVP |
| GET | `/v1/conexoes-telematica` | Lista conexões (Cobli, Geotab, rFMS Volvo/Scania, genérica) | Gestão | P | F2 |
| POST | `/v1/conexoes-telematica` | Cria conexão com credenciais do provedor | Admin | P | F2 |
| GET | `/v1/conexoes-telematica/{id}` | Detalha conexão e última coleta | Gestão | P | F2 |
| PATCH | `/v1/conexoes-telematica/{id}` | Altera conexão | Admin | P | F2 |
| DELETE | `/v1/conexoes-telematica/{id}` | Desativa conexão | Admin | P | F2 |
| POST | `/v1/conexoes-telematica/{id}/sincronizar` | Força coleta imediata | Gestão | P | F2 |
| GET | `/v1/conexoes-telematica/{id}/veiculos-externos` | Veículos do provedor e situação de casamento | Gestão | P | F2 |
| POST | `/v1/conexoes-telematica/{id}/veiculos-externos/{id_externo}/vincular` | Casamento manual com um ativo | Gestão | P | F2 |
| POST | `/v1/telematica/eventos` | Ingestão de eventos no modelo canônico (14.1) | Integração (`telematica:ingestao`) | P | F2 |
| POST | `/v1/telematica/receptor/{conexao_id}` | Recebe webhook do provedor, com assinatura própria dele | Provedor (assinatura) | P | F2 |
| GET | `/v1/eventos-telematica` | Eventos recebidos e o que geraram | Gestão | P | F2 |
| GET | `/v1/regras-telematica` | Regras "DTC ou alerta → solicitação ou OS" | Gestão | P | F2 |
| POST | `/v1/regras-telematica` | Cria regra | Gestão | P | F2 |
| PATCH | `/v1/regras-telematica/{id}` | Altera regra | Gestão | P | F2 |
| DELETE | `/v1/regras-telematica/{id}` | Remove regra | Gestão | P | F2 |
| GET | `/v1/importacoes` | Lista importações | Gestão, Estoque | E | MVP |
| POST | `/v1/importacoes` | Envia arquivo e valida sem gravar (14.3) | Gestão, Estoque, Integração | E | MVP |
| GET | `/v1/importacoes/{id}` | Situação e resumo da validação | Gestão, Estoque | E | MVP |
| GET | `/v1/importacoes/{id}/erros` | Linhas com erro (JSON ou CSV) | Gestão, Estoque | E | MVP |
| POST | `/v1/importacoes/{id}/confirmar` | Grava as linhas válidas (assíncrono) | Gestão, Estoque | E | MVP |
| POST | `/v1/importacoes/{id}/cancelar` | Descarta importação validada | Gestão, Estoque | E | MVP |
| GET | `/v1/importacoes/modelos/{tipo}` | Baixa o modelo CSV de um tipo | Gestão, Estoque | E | MVP |
| GET | `/v1/conectores-erp` | Lista conectores (Sankhya, Protheus, Omie) | Admin | Ent | F2 |
| POST | `/v1/conectores-erp` | Cria conector | Admin | Ent | F2 |
| GET | `/v1/conectores-erp/{id}` | Detalha conector e situação de funcionamento | Admin | Ent | F2 |
| PATCH | `/v1/conectores-erp/{id}` | Altera conector e fluxos ativos | Admin | Ent | F2 |
| DELETE | `/v1/conectores-erp/{id}` | Desativa conector | Admin | Ent | F2 |
| POST | `/v1/conectores-erp/{id}/testar` | Testa credenciais e conectividade | Admin | Ent | F2 |
| POST | `/v1/conectores-erp/{id}/sincronizar` | Força sincronização de um fluxo | Admin | Ent | F2 |
| GET | `/v1/conectores-erp/{id}/mapeamentos` | De-para (centro de custo, item, fornecedor, unidade) | Admin | Ent | F2 |
| PUT | `/v1/conectores-erp/{id}/mapeamentos` | Substitui o de-para | Admin | Ent | F2 |
| GET | `/v1/conectores-erp/{id}/execucoes` | Histórico de execuções com erros | Admin | Ent | F2 |
| POST | `/v1/conectores-erp/{id}/execucoes/{execucao_id}/reprocessar` | Reprocessa execução com falha | Admin | Ent | F2 |

### 6.16 IA — Assistente

| Método | Rota | Descrição | Perfil | Plano | Fase |
|---|---|---|---|---|---|
| POST | `/v1/ia/consultas` | Pergunta sobre manuais e histórico; resposta com citações; `Accept: text/event-stream` para streaming | Leitura | P | F2 |
| GET | `/v1/ia/consultas` | Histórico de consultas do usuário | Leitura | P | F2 |
| POST | `/v1/ia/consultas/{id}/avaliacao` | Avalia a resposta (útil ou não, comentário) | Leitura | P | F2 |
| POST | `/v1/ia/rascunhos-os` | Gera rascunho de OS a partir de voz, foto ou texto; **não grava a OS** | Execução | P | F2 |
| POST | `/v1/ia/rascunhos-solicitacao` | Gera rascunho de solicitação a partir de voz ou foto | Todos | P | F2 |
| POST | `/v1/ia/rascunhos-checklist` | Gera rascunho de modelo de checklist a partir de PDF | Gestão | P | F2 |
| GET | `/v1/ia/configuracao` | Governança: IA ligada por empresa, perfil e usuário; retenção de prompts | Admin | P | F2 |
| PATCH | `/v1/ia/configuracao` | Altera governança | Admin | P | F2 |
| GET | `/v1/ia/uso` | Créditos usados e disponíveis | Admin | P | F2 |
| GET | `/v1/ia/acoes-propostas` | Ações sugeridas por agentes (triagem, deduplicação, atribuição, programação) | Gestão | Ent | F3 |
| POST | `/v1/ia/acoes-propostas/{id}/aprovar` | Aprova e executa a ação, em nome do aprovador | Gestão | Ent | F3 |
| POST | `/v1/ia/acoes-propostas/{id}/rejeitar` | Rejeita com motivo | Gestão | Ent | F3 |
| POST | `/mcp` | Servidor MCP somente leitura (Streamable HTTP), autenticado por token de API | Integração | Ent | F3 |

Regras de IA para a API: nenhum endpoint de IA grava dado de negócio sem confirmação humana (o rascunho é confirmado por `POST /v1/ordens-servico` ou equivalente); toda resposta de consulta traz `citacoes` (anexo, página, OS); cada chamada consome créditos do plano. Quando os créditos acabam, a API devolve `403 plano_nao_inclui_recurso` com o membro `motivo: "creditos_esgotados"` (proposta; evita usar o `402`, que não tem semântica padronizada).

### 6.17 MOB — App de campo e sincronização

| Método | Rota | Descrição | Perfil | Plano | Fase |
|---|---|---|---|---|---|
| GET | `/v1/configuracao-app` | Versão mínima do app, recursos ligados e URLs | Público | E | MVP |
| POST | `/v1/dispositivos` | Registra dispositivo e token de notificação push | Todos, Externos | E | MVP |
| GET | `/v1/dispositivos` | Lista dispositivos (próprios; administrador vê todos) | Todos, Externos, Admin | E | MVP |
| DELETE | `/v1/dispositivos/{id}` | Revoga dispositivo e seus refresh tokens | Todos (próprio), Externos (próprio), Admin | E | MVP |
| GET | `/v1/sincronizacao/escopo` | O que entra no escopo offline do usuário e tamanho estimado | Todos, `prestador` | E | MVP |
| GET | `/v1/sincronizacao/instantaneo` | Carga inicial paginada do escopo (seção 12.2) | Todos, `prestador` | E | MVP |
| GET | `/v1/sincronizacao/alteracoes` | Pull incremental por cursor (seção 12.3) | Todos, `prestador` | E | MVP |
| POST | `/v1/sincronizacao/comandos` | Push de lote de comandos idempotentes (seção 12.4) | Todos, `prestador` | E | MVP |
| GET | `/v1/sincronizacao/comandos/{id_comando}` | Resultado de um comando já enviado | Todos, `prestador` | E | MVP |

### 6.18 PRJ — Projetos e implantação

Módulo da Fase 2, plano Profissional (sumário 2.1 e 4; RF-PRJ-001 a RF-PRJ-007). Pacote Go `projetos`. O projeto não executa trabalho: organiza OS `implantacao`, que seguem as rotas de OS (6.5). Detalhamento na seção 24.

| Método | Rota | Descrição | Perfil | Plano | Fase |
|---|---|---|---|---|---|
| GET | `/v1/modelos-projeto` | Lista modelos de projeto (`classe_ativo_id`, `situacao`, `busca`) | Gestão | P | F2 |
| POST | `/v1/modelos-projeto` | Cria modelo com as etapas (`modelo_etapa_projeto`), dependências e exigências (24.2) | Projetos | P | F2 |
| GET | `/v1/modelos-projeto/{id}` | Detalha modelo na versão mais recente (`numero_versao` escolhe outra) | Gestão | P | F2 |
| PATCH | `/v1/modelos-projeto/{id}` | Altera modelo; mudança nas etapas gera `numero_versao` nova, sem afetar projetos existentes | Projetos | P | F2 |
| DELETE | `/v1/modelos-projeto/{id}` | Inativa modelo | Projetos | P | F2 |
| POST | `/v1/modelos-projeto/{id}/copiar` | Cria modelo novo a partir de outro | Projetos | P | F2 |
| GET | `/v1/modelos-projeto/{id}/versoes` | Versões do modelo | Gestão | P | F2 |
| GET | `/v1/projetos` | Lista projetos (`status`, `unidade_id`, `modelo_projeto_id`, `cliente_servico_id`, `contrato_cliente_id`, `responsavel_id`, `fornecedor_id`, `ativo_id`, `atrasado`, `numero`, períodos); traz avanço físico e desvio de prazo | Leitura | P | F2 |
| POST | `/v1/projetos` | Cria projeto `rascunho` a partir de um modelo, com etapas copiadas e datas previstas calculadas (24.3) | Projetos | P | F2 |
| GET | `/v1/projetos/{id}` | Detalha projeto, com etapas resumidas | Leitura | P | F2 |
| PATCH | `/v1/projetos/{id}` | Altera projeto; em `rascunho`, também ativos e localização; `concluido` ou `cancelado` devolve `409 registro_imutavel` | Projetos | P | F2 |
| POST | `/v1/projetos/{id}/iniciar` | `rascunho` → `em_andamento`; torna `disponivel` as etapas sem predecessoras | Projetos | P | F2 |
| POST | `/v1/projetos/{id}/suspender` | `em_andamento` → `suspenso`, com motivo; pausa os prazos das etapas, não as OS | Projetos | P | F2 |
| POST | `/v1/projetos/{id}/retomar` | `suspenso` → `em_andamento` | Projetos | P | F2 |
| POST | `/v1/projetos/{id}/concluir` | `em_andamento` → `concluido`; coloca em operação os ativos ainda `em_implantacao` (24.6) | Projetos + Entrada em operação | P | F2 |
| POST | `/v1/projetos/{id}/cancelar` | `rascunho`, `em_andamento` ou `suspenso` → `cancelado`, com motivo e o destino de cada ativo `em_implantacao` | Projetos | P | F2 |
| GET | `/v1/projetos/{id}/acompanhamento` | Cronograma previsto e real com dependências, avanço físico, etapas atrasadas, custo previsto e realizado por etapa e total, desvio de prazo (RF-PRJ-005) | Leitura | P | F2 |
| GET | `/v1/projetos/{id}/lancamentos-custo` | Lançamentos avulsos de custo do projeto (proposta) | Gestão | P | F2 |
| POST | `/v1/projetos/{id}/lancamentos-custo` | Registra lançamento avulso (etapa, descrição, valor, data, anexo) (proposta) | Projetos | P | F2 |
| DELETE | `/v1/projetos/{id}/lancamentos-custo/{lancamento_id}` | Exclui lançamento de projeto não concluído (proposta) | Projetos | P | F2 |
| GET | `/v1/projetos/{id}/termo-entrega` | Termo de entrega em PDF: ativos, etapas, documentos, fotos e assinaturas (RF-PRJ-007) | Leitura | P | F2 |
| POST | `/v1/projetos/{id}/assinaturas` | Registra a assinatura do responsável ou do representante do cliente atendido no termo de entrega | Projetos | P | F2 |
| GET | `/v1/projetos/{id}/etapas` | Etapas do projeto, com situação, predecessoras, exigências, OS e atraso | Leitura | P | F2 |
| GET | `/v1/projetos/{id}/etapas/{etapa_id}` | Detalha etapa | Leitura | P | F2 |
| PATCH | `/v1/projetos/{id}/etapas/{etapa_id}` | Altera executor (`executor_tipo`, `fornecedor_id`), responsável, datas previstas, valor e campos exigidos de etapa sem OS | Projetos | P | F2 |
| POST | `/v1/projetos/{id}/etapas/{etapa_id}/gerar-os` | Gera a OS `implantacao` da etapa, uma só ou uma por ativo do projeto (24.5) | Projetos | P | F2 |
| POST | `/v1/projetos/{id}/etapas/{etapa_id}/iniciar` | `disponivel` → `em_andamento` em etapa sem OS; com OS, a etapa inicia sozinha quando a primeira OS inicia | Projetos | P | F2 |
| POST | `/v1/projetos/{id}/etapas/{etapa_id}/concluir` | → `concluida`; torna `disponivel` as sucessoras cujas predecessoras terminaram e notifica seus responsáveis (24.4) | Projetos | P | F2 |
| POST | `/v1/projetos/{id}/etapas/{etapa_id}/dispensar` | `pendente` ou `disponivel` → `dispensada`, com motivo | Dispensa de etapa | P | F2 |
| GET | `/v1/medicoes-etapa` | Medições de etapas por fornecedor (`fornecedor_id`, `situacao`, período) (proposta) | Gestão, Compras | P | F2 |
| POST | `/v1/medicoes-etapa` | Gera a medição de um fornecedor e período com as etapas `concluida` ainda não medidas (24.7) | Projetos | P | F2 |
| GET | `/v1/medicoes-etapa/{id}` | Detalha medição: linhas (projeto, etapa, ativo, data de conclusão, valor, ajuste), total | Gestão, Compras | P | F2 |
| DELETE | `/v1/medicoes-etapa/{id}` | Apaga medição `em_apuracao`; as etapas voltam a poder entrar em outra medição | Projetos | P | F2 |
| POST | `/v1/medicoes-etapa/{id}/ajustar` | Ajusta o valor de uma linha, com motivo | Ajuste de medição | P | F2 |
| POST | `/v1/medicoes-etapa/{id}/fechar` | `em_apuracao` → `fechada` (proposta) | Projetos | P | F2 |
| GET | `/v1/medicoes-etapa/{id}/relatorio` | Medição em PDF ou XLSX | Gestão, Compras | P | F2 |

A criação de vários projetos de uma vez, um por ponto ou por ativo, usa a importação por CSV com o tipo `projetos` (14.3), não uma rota própria.

### 6.19 CTR — Contratos com clientes

Módulo da Fase 2, plano Profissional (sumário 2.1 e 4; RF-CTR-001 a RF-CTR-005). Pacote Go `contratos`. Usa o cadastro `cliente_servico` (`/v1/clientes-servico`, 6.10), as OS e os registros de parada, sem mudar o cálculo técnico da disponibilidade. Detalhamento na seção 25.

| Método | Rota | Descrição | Perfil | Plano | Fase |
|---|---|---|---|---|---|
| GET | `/v1/contratos-cliente` | Lista contratos com clientes (`cliente_servico_id`, `status`, `vigente_em`, `ativo_id`, `unidade_id`, `numero`) | Gestão | P | F2 |
| POST | `/v1/contratos-cliente` | Cria contrato em `rascunho`: vigência, meta, calendário de apuração, prazos por prioridade, faixas de desconto, causas externas aceitas (25.2) | Contratos | P | F2 |
| GET | `/v1/contratos-cliente/{id}` | Detalha contrato | Gestão | P | F2 |
| PATCH | `/v1/contratos-cliente/{id}` | Altera contrato; em `vigente`, só data de fim, renovação, contatos e anexos; `encerrado` ou `cancelado` devolve `409 registro_imutavel` | Contratos | P | F2 |
| POST | `/v1/contratos-cliente/{id}/ativar` | `rascunho` → `vigente` | Contratos | P | F2 |
| POST | `/v1/contratos-cliente/{id}/encerrar` | `vigente` → `encerrado`, com motivo quando antes do fim da vigência | Contratos | P | F2 |
| POST | `/v1/contratos-cliente/{id}/cancelar` | Só `rascunho` → `cancelado`, com motivo | Contratos | P | F2 |
| GET | `/v1/contratos-cliente/{id}/itens` | Ativos cobertos (`item_contrato_cliente`), com datas de inclusão e exclusão e disponibilidade corrente | Gestão | P | F2 |
| POST | `/v1/contratos-cliente/{id}/itens` | Inclui ativo (em `rascunho`, ou em `vigente` como aditivo) (25.3) | Contratos | P | F2 |
| PATCH | `/v1/contratos-cliente/{id}/itens/{item_id}` | Altera meta própria, valor mensal ou registra a data de exclusão | Contratos | P | F2 |
| DELETE | `/v1/contratos-cliente/{id}/itens/{item_id}` | Remove item de contrato em `rascunho`; em `vigente`, use a data de exclusão | Contratos | P | F2 |
| GET | `/v1/contratos-cliente/{id}/painel` | Disponibilidade parcial do mês por item e do contrato, itens abaixo da meta, desconto projetado, OS abertas com prazos, paradas aguardando decisão (25.6) | Gestão | P | F2 |
| GET | `/v1/medicoes-contrato` | Medições (`contrato_cliente_id`, `periodo`, `situacao`) | Gestão | P | F2 |
| POST | `/v1/medicoes-contrato` | Abre a medição de um período (também aberta sozinha no primeiro dia do mês seguinte, proposta) | Contratos | P | F2 |
| GET | `/v1/medicoes-contrato/{id}` | Detalha medição: por item, horas requeridas, paradas, exclusões, disponibilidade contratual, faixa e desconto; memória de cálculo (25.5) | Gestão | P | F2 |
| POST | `/v1/medicoes-contrato/{id}/decidir-exclusao` | Grava `decisao_exclusao` (`aprovada` ou `recusada`) de uma parada com causa externa, com motivo | Fechamento de medição | P | F2 |
| POST | `/v1/medicoes-contrato/{id}/fechar` | `em_apuracao` → `fechada` | Fechamento de medição | P | F2 |
| POST | `/v1/medicoes-contrato/{id}/reabrir` | `fechada` → `em_apuracao`, com motivo; mostra as divergências desde o fechamento | Reabertura de medição | P | F2 |
| GET | `/v1/medicoes-contrato/{id}/relatorio` | Medição em PDF ou XLSX, com a memória de cálculo | Gestão | P | F2 |

### 6.20 Resumo do catálogo

| Módulo | Endpoints | MVP | F2 | F3 | Variação desde a versão 0.3 |
|---|---|---|---|---|---|
| ADM | 106 | 83 | 23 | 0 | +17 na versão 0.4: regras configuráveis (9) e disparos (2), no MVP; pacotes de segmento (4) e aplicações (2), na Fase 2. +19 na versão 0.5 (E01): `POST /v1/empresas`, convites (2), MFA do próprio usuário (3), `trocar-senha`, `POST /oauth/authorize`, turnos (3), feriados (3), reativação de unidade e de pessoa, reenvio de convite e permissões avulsas (2); `DELETE /v1/unidades/{id}` passa da Fase 2 ao MVP |
| ATV | 27 | 27 | 0 | 0 | +2: `colocar-em-operacao`, unitária e em lote (campos `fase_ciclo_vida`, `data_entrada_operacao`, `projeto_id`) |
| MED | 12 | 12 | 0 | 0 | — |
| SOL | 11 | 11 | 0 | 0 | — (origem `regra`) |
| OS | 51 | 48 | 3 | 0 | +8: definições de campos extras (6) e campos extras da OS (2) (tipo `implantacao`, origens `projeto` e `cliente_conectado`, `causa_externa`, `etapa_projeto_id`) |
| PRE | 19 | 19 | 0 | 0 | — |
| CHK | 11 | 11 | 0 | 0 | — |
| EST | 59 | 18 | 41 | 0 | — |
| TER | 35 | 20 | 15 | 0 | — |
| SRV | 23 | 23 | 0 | 0 | — (clientes atendidos também no Profissional com CTR) |
| PRJ | 36 | 0 | 36 | 0 | Novo módulo: modelos (7), projetos com ações (9), acompanhamento (1), lançamentos de custo (3), termo de entrega (2), etapas (7), medições de etapas (7) |
| CTR | 19 | 0 | 19 | 0 | Novo módulo: contratos com ações (7), itens (4), painel (1), medições do contrato (7) |
| LOC | 25 | 0 | 25 | 0 | — |
| FRO | 32 | 4 | 28 | 0 | — |
| CON | 31 | 19 | 12 | 0 | — (`bloqueio_ativo.origem`) |
| IND | 12 | 12 | 0 | 0 | — |
| INT | 52 | 26 | 26 | 0 | — |
| IA | 13 | 0 | 9 | 4 | — |
| MOB | 9 | 9 | 0 | 0 | — |
| **Total** | **583** | **342** | **237** | **4** | **+82 na versão 0.4; +19 na versão 0.5** |

ADM inclui autenticação (19 rotas, entre elas as 4 de OAuth fora de `/v1`), SCIM (10 rotas em `/scim/v2`), pessoas, versão do produto, anexos, processamentos, notificações e exportação completa. Os sub-recursos de frota e de conformidade pendurados em `ativos` e `ordens-servico` (dados de veículo, pneus no mapa, registro de bloqueio, NR-13) são contados nos módulos FRO e CON; as ações `atribuir-fornecedor` e `retirar-fornecedor` e as 4 rotas públicas do link de uso único, em TER; `assumir`, `servicos` e `relatorio-pdf`, em OS. O SRV conta só as rotas próprias (22.2); a conta da subcontratada usa também as rotas de ADM, ATV, MED, OS, CHK, EST, IND e MOB. As regras configuráveis e os pacotes de segmento contam em ADM; as definições e os valores de campos extras, em OS; `colocar-em-operacao`, em ATV; `gerar-os` da etapa, em PRJ; as medições de etapas (`/v1/medicoes-etapa`), em PRJ; as medições do contrato (`/v1/medicoes-contrato`), em CTR.

---

## 7. Solicitações: detalhamento

### 7.1 Estados e ações

Estados canônicos (sumário 3.3 e 3.6): `nova` → `em_triagem` → `aprovada` (gera OS) ou `recusada`; `nova` → `recusada` direto, na desistência; `duplicada` a partir de `nova` ou `em_triagem`, vinculada a outra solicitação ou a uma OS. `aprovada`, `recusada` e `duplicada` são finais.

| Ação | Rota | De | Para | Campos obrigatórios |
|---|---|---|---|---|
| Abrir | `POST /v1/solicitacoes` | — | `nova` | `ativo_id` ou `localizacao_id`; `descricao` |
| Iniciar triagem | `POST /v1/solicitacoes/{id}/iniciar-triagem` | `nova` | `em_triagem` | — (grava `data_reconhecimento` da solicitação) |
| Aprovar | `POST /v1/solicitacoes/{id}/aprovar` | `nova`, `em_triagem` | `aprovada` | `ordem_servico.tipo` |
| Recusar | `POST /v1/solicitacoes/{id}/recusar` | `nova`, `em_triagem` | `recusada` | `motivo_recusa` |
| Desistir (autor) | `POST /v1/solicitacoes/{id}/recusar` com `motivo_recusa = desistencia` | `nova` (direto, sem triagem), `em_triagem` | `recusada` | `motivo_recusa = desistencia` |
| Marcar duplicada | `POST /v1/solicitacoes/{id}/marcar-duplicada` | `nova`, `em_triagem` | `duplicada` | `solicitacao_original_id` (não pode estar `recusada` nem `duplicada`) ou `ordem_servico_id` |
| Reabrir como nova | `POST /v1/solicitacoes/{id}/reabrir-como-nova` | `duplicada` | Original continua `duplicada`; nasce outra em `nova` | `motivo` |

Aprovar direto de `nova`, sem passar por `em_triagem`, é permitido (proposta): a triagem explícita serve a equipes com fila de análise. Nesse caso a aprovação grava `data_reconhecimento`. Desistência do solicitante não é estado próprio: é `recusada` com motivo `desistencia`, vai direto de `nova` para `recusada` (sumário 3.6), e o autor só pode usar esse motivo.

`duplicada` é final (sumário 3.6). Marcação errada se corrige com `reabrir-como-nova`: a API cria uma solicitação nova, com número novo, mesmo ativo, descrição, anexos e solicitante, `origem` igual à da original e `solicitacao_reaberta_de_id` apontando para a duplicada; a duplicada ganha `reaberta_como_id`. Resposta `201` com a solicitação nova; webhook `solicitacao.criada`. Reabrir de novo a mesma duplicada devolve `409 transicao_invalida`. O locatário abre solicitação só para ativo locado a ele (4.9); a OS gerada recebe `contrato_locacao_id` e o `responsavel_custo` padrão da matriz do contrato (21.3).

### 7.2 Criar solicitação

`POST /v1/solicitacoes`

```
POST /v1/solicitacoes HTTP/1.1
Authorization: Bearer eyJhbGciOiJFZERTQSIs...
Idempotency-Key: 6f1d2c9e-8a47-4b1e-9c3a-7d5e2f0a1b34
Content-Type: application/json
```

```json
{
  "id": "0193a2e1-4c5d-7e6f-8a9b-0c1d2e3f4a5b",
  "ativo_id": "0192f7c4-1a2b-7c3d-8e4f-5a6b7c8d9e01",
  "titulo": "Vazamento de óleo no compressor",
  "descricao": "Poça de óleo embaixo da unidade compressora, lado do motor. Começou no turno da manhã.",
  "ativo_parado": false,
  "prioridade_sugerida": "alta",
  "anexos_ids": ["0193a2e0-9f8e-7d6c-8b5a-4f3e2d1c0b9a"]
}
```

- `id` é opcional. O app envia o id gerado offline; a web e as integrações podem omitir.
- `localizacao_id` substitui `ativo_id` quando o problema é predial e não há ativo cadastrado.
- `origem` é definida pelo servidor (`web`, `app`, `qr_code`, `checklist`, `telematica`, `api`, `regra`, `portal_locatario`) e não é aceita no corpo. Solicitação aberta por usuário `locatario` recebe sempre `portal_locatario` (sumário 3.6). A de origem `regra` é criada pelo efeito `abrir_solicitacao` de uma regra configurável (RF-SOL-001, RF-ADM-016), traz `disparo_regra_id` (só leitura, proposta), tem como solicitante o autor de sistema `{ "tipo": "regra_automacao", "id": "<id da regra>" }` e segue a triagem normal (23.1).
- `ativo_parado: true` ("o ativo parou?") por falha coloca o ativo em `parado` desde a abertura (6.2).
- O solicitante é sempre o usuário da credencial; integrações registram o token e podem informar `solicitante_nome` e `solicitante_contato`.

**Resposta `201 Created`**

```
Location: /v1/solicitacoes/0193a2e1-4c5d-7e6f-8a9b-0c1d2e3f4a5b
ETag: "1"
```

```json
{
  "id": "0193a2e1-4c5d-7e6f-8a9b-0c1d2e3f4a5b",
  "numero": "SS-2026-000731",
  "versao": 1,
  "status": "nova",
  "ativo_id": "0192f7c4-1a2b-7c3d-8e4f-5a6b7c8d9e01",
  "localizacao_id": "0192f7a0-0000-7aaa-8bbb-00000000000c",
  "unidade_id": "0192f790-0000-7aaa-8bbb-000000000001",
  "titulo": "Vazamento de óleo no compressor",
  "descricao": "Poça de óleo embaixo da unidade compressora, lado do motor. Começou no turno da manhã.",
  "ativo_parado": false,
  "prioridade_sugerida": "alta",
  "origem": "app",
  "solicitante_id": "0192f800-aaaa-7bbb-8ccc-000000000042",
  "anexos_ids": ["0193a2e0-9f8e-7d6c-8b5a-4f3e2d1c0b9a"],
  "data_abertura": "2026-11-12T11:02:17Z",
  "data_reconhecimento": null,
  "motivo_recusa": null,
  "solicitacao_original_id": null,
  "ordem_servico_id": null,
  "criado_em": "2026-11-12T11:02:19Z",
  "atualizado_em": "2026-11-12T11:02:19Z",
  "criado_por_id": "0192f800-aaaa-7bbb-8ccc-000000000042"
}
```

`data_abertura` é o momento informado pelo dispositivo (o registro pode ter sido feito offline); `criado_em` é o momento em que o servidor recebeu. Efeitos: notificação para os perfis de triagem da unidade e webhook `solicitacao.criada`.

### 7.3 Aprovar e gerar OS

`POST /v1/solicitacoes/{id}/aprovar` com `If-Match: "1"`:

```json
{
  "ordem_servico": {
    "tipo": "corretiva_planejada",
    "prioridade": "alta",
    "responsavel_id": "0192f800-aaaa-7bbb-8ccc-000000000017",
    "data_limite": "2026-11-14T21:00:00Z"
  },
  "comentario": "Aprovada. Verificar retentor do eixo."
}
```

**Resposta `200 OK`**

```json
{
  "solicitacao": {
    "id": "0193a2e1-4c5d-7e6f-8a9b-0c1d2e3f4a5b",
    "numero": "SS-2026-000731",
    "versao": 2,
    "status": "aprovada",
    "ordem_servico_id": "0193a2f0-5b6c-7d8e-9f01-23456789abcd"
  },
  "ordem_servico": {
    "id": "0193a2f0-5b6c-7d8e-9f01-23456789abcd",
    "numero": "OS-2026-004812",
    "versao": 1,
    "status": "aberta",
    "tipo": "corretiva_planejada",
    "solicitacao_id": "0193a2e1-4c5d-7e6f-8a9b-0c1d2e3f4a5b"
  }
}
```

A OS herda ativo, descrição, anexos e `ativo_parado` da solicitação. Webhooks: `solicitacao.aprovada` e `ordem_servico.criada`.

---

## 8. Ordens de serviço: detalhamento

### 8.1 Criar OS

`POST /v1/ordens-servico`

```json
{
  "tipo": "corretiva_planejada",
  "ativo_id": "0192f7c4-1a2b-7c3d-8e4f-5a6b7c8d9e01",
  "titulo": "Substituir retentor do eixo da unidade compressora",
  "descricao": "Vazamento de óleo identificado na SS-2026-000731.",
  "prioridade": "alta",
  "ativo_parado": false,
  "data_limite": "2026-11-14T21:00:00Z",
  "procedimento_id": null,
  "executor_tipo": "interno",
  "fornecedor_id": null,
  "responsavel_custo": null,
  "campos_extras": { "numero_protocolo_concessionaria": null },
  "causa_externa": null,
  "tarefas": [
    { "descricao": "Aplicar bloqueio e etiquetagem da fonte elétrica", "obrigatoria": true },
    { "descricao": "Substituir retentor do eixo principal", "obrigatoria": true },
    { "descricao": "Testar proteções e dispositivos de parada de emergência", "obrigatoria": true }
  ],
  "pecas_previstas": [
    { "item_estoque_id": "0192fa10-3333-7ddd-8eee-000000000101", "quantidade": "1", "reservar": true }
  ]
}
```

| Campo | Obrigatório | Regra |
|---|---|---|
| `tipo` | Sim | `preventiva`, `preditiva`, `corretiva_planejada`, `corretiva_emergencial`, `inspecao_legal`, `melhoria`, `implantacao`. `implantacao` só em ativo `em_implantacao` ou em OS gerada por etapa de projeto (`422 ativo_nao_em_implantacao`); plano nunca gera OS `implantacao` (RF-OS-002) |
| `ativo_id` | Sim | Ativo ativo da empresa, em unidade acessível ao usuário. Ativo `em_implantacao` (status `inativo`) é aceito só em OS `implantacao` (`422 ativo_em_implantacao` nos demais tipos, proposta) |
| `titulo` | Sim | Até 200 caracteres |
| `descricao` | Não | Até 10.000 caracteres |
| `prioridade` | Não | `baixa`, `media`, `alta`, `urgente`. Padrão pela criticidade do ativo: A → `alta`, B → `media`, C → `baixa`; `corretiva_emergencial` nasce `urgente` |
| `ativo_parado` | Não | "Parou o ativo?". Padrão `false` |
| `procedimento_id` | Não | Copia tarefas, peças e ferramentas do procedimento |
| `executor_tipo` | Não | `interno` (padrão) ou `fornecedor`. Com `fornecedor`, exige `fornecedor_id` com papel `subcontratada`, `assistencia_tecnica` ou `oficina` (ou `locador`, em ativo tomado); senão `422 fornecedor_papel_incompativel` |
| `fornecedor_id` | Quando `executor_tipo = fornecedor` | Prestador que executa a OS. Também pode ser definido depois por `atribuir-fornecedor` (19.3) |
| `responsavel_custo` | Não | `empresa`, `locatario`, `locador`, `garantia` ou `fornecedor`. Padrão: `empresa`; em ativo com `contrato_locacao_id`, o valor da matriz do contrato para o tipo de OS (21.3). Informar valor diferente do padrão exige `justificativa_responsavel_custo` (`422 motivo_obrigatorio`) |
| `descricao_servico` | Não | Serviço a realizar, em texto livre (até 10.000 caracteres). Usado sobretudo no SRV (22.7) |
| `servicos` | Não | Itens `servico_os`: `servico_fornecedor_id` do catálogo da própria conta, `quantidade`, `preco_referencia` (padrão: o do catálogo) |
| `responsavel_id`, `executores_ids` | Não | Técnico responsável e demais executores; na criação, só para OS que nasce `aberta` (programar continua sendo a ação `programar`) |
| `especialidade_tecnico_ids` | Não | Especialidades exigidas; atribuir a técnico sem elas gera aviso `tecnico_sem_especialidade` (22.3) |
| `enviar_ao_cliente` | Não | Só em equipamento sincronizado de cliente conectado (conta SRV). Padrão `true`: cria a OS também na conta do cliente (22.8) |
| `tarefas`, `pecas_previstas` | Não | Criadas junto com a OS, na mesma transação |
| `solicitacao_id`, `plano_manutencao_id` | Não | Só leitura quando a OS nasce de solicitação ou plano. OS gerada por plano nasce `planejada` |
| `contrato_locacao_id` | — | Só leitura: preenchido pelo servidor quando o ativo está em contrato de locação vigente, cedido ou tomado |
| `campos_extras` | Quando há campo obrigatório na abertura | Objeto com os valores dos campos extras, pela `chave` da definição (23.2). Campo obrigatório na abertura vazio: `422 campo_extra_obrigatorio`. Chave desconhecida, inativa ou que não se aplica ao tipo e à classe do ativo: `422 validacao_falhou` (`campos_extras.<chave>`, `nao_encontrado` ou `inativo`) |
| `causa_externa` | Não | Causa externa da parada (sumário 3.6): `vandalismo`, `acidente_terceiro`, `falta_energia`, `obra_terceiro`, `solicitacao_cliente`, `caso_fortuito`. Pode ser marcada também durante a execução, por `PATCH`; não muda MTTR nem a disponibilidade técnica (25.4) |
| `etapa_projeto_id` | — | Só leitura: etapa de projeto de origem, preenchida quando a OS nasce de `gerar-os` da etapa (24.5). `null` para `prestador` e `locatario` (4.9) |

O tipo `implantacao` fica fora dos KPIs de manutenção (MTBF, MTTR, disponibilidade, % planejada, cumprimento de preventiva, backlog) e seu custo sai como custo de implantação, separado do custo de manutenção (RF-OS-002, RF-IND-005).

**Resposta `201 Created`** (`Location: /v1/ordens-servico/0193a2f0-...`, `ETag: "1"`)

```json
{
  "id": "0193a2f0-5b6c-7d8e-9f01-23456789abcd",
  "numero": "OS-2026-004812",
  "versao": 1,
  "tipo": "corretiva_planejada",
  "status": "aberta",
  "prioridade": "alta",
  "titulo": "Substituir retentor do eixo da unidade compressora",
  "descricao": "Vazamento de óleo identificado na SS-2026-000731.",
  "descricao_servico": null,
  "origem": "manual",
  "vinculo": null,
  "ativo_id": "0192f7c4-1a2b-7c3d-8e4f-5a6b7c8d9e01",
  "unidade_id": "0192f790-0000-7aaa-8bbb-000000000001",
  "localizacao_id": "0192f7a0-0000-7aaa-8bbb-00000000000c",
  "solicitacao_id": null,
  "plano_manutencao_id": null,
  "ativo_parado": false,
  "exige_registro_nr12": true,
  "executor_tipo": "interno",
  "fornecedor_id": null,
  "responsavel_custo": "empresa",
  "justificativa_responsavel_custo": null,
  "contrato_locacao_id": null,
  "etapa_projeto_id": null,
  "campos_extras": { "numero_protocolo_concessionaria": null },
  "causa_externa": null,
  "responsavel_id": null,
  "executores_ids": [],
  "tempo_estimado_min": null,
  "data_abertura": "2026-11-12T11:15:40Z",
  "data_reconhecimento": null,
  "data_programada": null,
  "data_limite": "2026-11-14T21:00:00Z",
  "data_inicio": null,
  "data_conclusao": null,
  "data_encerramento": null,
  "data_cancelamento": null,
  "motivo_pausa": null,
  "motivo_cancelamento": null,
  "motivo_recusa_liberacao": null,
  "servico_realizado": null,
  "codigo_falha_modo_id": null,
  "codigo_falha_causa_id": null,
  "codigo_falha_atividade_id": null,
  "custo_mao_obra": "0.00",
  "custo_pecas": "0.00",
  "custo_terceiros": "0.00",
  "custo_total": "0.00",
  "registro_nr12": null,
  "transicoes_permitidas": ["planejar", "programar", "reconhecer", "assumir", "cancelar"],
  "criado_em": "2026-11-12T11:15:41Z",
  "atualizado_em": "2026-11-12T11:15:41Z",
  "criado_por_id": "0192f800-aaaa-7bbb-8ccc-000000000017"
}
```

- `exige_registro_nr12` é calculado: `true` quando o ativo está marcado como máquina sujeita à NR-12.
- `transicoes_permitidas` lista as ações que **este usuário** pode executar agora, considerando estado e perfil. O web e o app usam esse campo para mostrar botões, sem replicar a regra.
- `origem` é definida pelo servidor e nunca aceita no corpo, com os valores do sumário 3.6:

| `origem` | Quando |
|---|---|
| `manual` | Criada por `POST /v1/ordens-servico` (web, app ou token de API) |
| `solicitacao` | Gerada pela aprovação de uma solicitação (7.3) |
| `plano` | Gerada por plano de manutenção; nasce `planejada` |
| `checklist` | Gerada por item reprovado de checklist com regra de abertura de OS |
| `telematica` | Gerada por regra de telemática (14.1) |
| `rascunho_ia` | Criada a partir de um rascunho do assistente confirmado pelo usuário (6.16) |
| `vistoria_locacao` | Avaria encontrada na vistoria de devolução (21.4) |
| `prestador` | Na conta do cliente, OS aberta pela subcontratada conectada num ativo compartilhado; nasce `aberta` e segue a triagem do cliente (22.8) |
| `cliente_conectado` | Na conta da subcontratada, OS criada pela sincronização a partir de uma OS do cliente conectado (22.8) |
| `projeto` | Gerada pela ação `gerar-os` de uma etapa de projeto, com `tipo = implantacao` e `etapa_projeto_id` (24.5) |

O enum `origem` é extensível (`x-enum-extensivel`, 18.3). Não existe origem `regra` na OS (sumário 3.6): nenhum efeito de regra configurável abre OS; a regra abre solicitação (origem `regra` da solicitação, 7.2), e a OS gerada pela aprovação dela tem `origem = solicitacao`.

- `vinculo` é `null` em OS sem vínculo entre contas; nas vinculadas, traz o resumo da seção 22.8.
- Webhook: `ordem_servico.criada`.

### 8.2 Máquina de estados e transições

Conforme o diagrama da seção 3.3 do sumário:

```
[criada] ──────────────────────→ aberta
[gerada por plano] ────────────→ planejada

aberta ──planejar──→ planejada ──programar──→ programada ──iniciar──→ em_execucao ⇄ pausada
aberta ──programar (sem planejar)──→ programada
programada ──desprogramar──→ planejada
aberta ──iniciar (só corretiva_emergencial)──→ em_execucao
em_execucao ──concluir──→ concluida ──encerrar (liberação)──→ encerrada
concluida ──recusar-liberacao (com motivo)──→ em_execucao
aberta | planejada | programada | pausada ──cancelar──→ cancelada
```

`em_execucao`, `concluida` e `encerrada` **não** podem ser cancelados: de `em_execucao`, pausar antes. Toda ação abaixo é `POST /v1/ordens-servico/{id}/<acao>`, exige `If-Match` e `Idempotency-Key` e devolve `200` com a OS atualizada e novo `ETag`. Ação fora do estado permitido devolve `409 transicao_invalida`; campo exigido ausente devolve `422 campos_obrigatorios_transicao`.

| Ação | De | Para | Campos obrigatórios | Campos opcionais | Efeitos |
|---|---|---|---|---|---|
| `planejar` | `aberta` | `planejada` | `tempo_estimado_min`; ao menos uma tarefa existente ou `procedimento_id` | `pecas_previstas`, `ferramentas` | Reserva peças marcadas com `reservar` |
| `programar` | `planejada`, `aberta` | `programada` | `data_programada`, `responsavel_id` | `executores_ids`, `data_limite` | Grava `data_reconhecimento` se vazia; notifica executores; na Fase 2, valida capacitações (`422 capacitacao_vencida`) |
| `desprogramar` | `programada` | `planejada` | — | `observacao` | Limpa `data_programada`; mantém responsável e reservas; notifica executores |
| `reconhecer` | `aberta`, `planejada`, `programada` | Mesmo estado | — | — | Grava `data_reconhecimento` uma única vez (base do MTTA) |
| `assumir` | `aberta`, `planejada` | Mesmo estado | — (o responsável é sempre o usuário da credencial) | `observacao` | Define `responsavel_id` = usuário e o inclui em `executores_ids`; grava `data_reconhecimento` se vazia. Só para OS da equipe do técnico (unidade e equipe a que ele pertence) sem responsável; OS com outro responsável devolve `409 transicao_invalida`. O `prestador` assume só OS do seu fornecedor. Notifica o planejador |
| `iniciar` | `programada`; `aberta` só em `corretiva_emergencial` | `em_execucao` | `data_inicio` (padrão: agora) | `ativo_parado`, `leituras` | Grava `data_reconhecimento` se vazia; com `ativo_parado`, recalcula o status do ativo (`parado` em corretiva emergencial ou falha; `em_manutencao` em retirada programada; 6.2); o usuário entra em `executores_ids`; abre apontamento automático (proposta) |
| `pausar` | `em_execucao` | `pausada` | `motivo_pausa`; `observacao_pausa` quando o motivo é `outro` | `data_pausa` | Fecha o apontamento aberto |
| `retomar` | `pausada` | `em_execucao` | — | `data_retomada` | Abre apontamento automático (proposta) |
| `concluir` | `em_execucao` | `concluida` | `data_conclusao`, `servico_realizado`, tarefas obrigatórias concluídas, ao menos um apontamento (proposta); `leituras` quando a OS veio de plano por medidor; códigos de falha em corretivas, quando o parâmetro exigir; campos extras obrigatórios na conclusão (`422 campo_extra_obrigatorio`, 23.2) | `registro_nr12` (rascunho), `anexos_ids`, `assinatura_anexo_id`, `campos_extras`, `causa_externa` | Fecha apontamentos abertos; registra leituras; avalia as regras com condição `os_concluida_com_campo` (23.1); notifica quem tem `os.encerrar` |
| `encerrar` (liberação) | `concluida` | `encerrada` | Permissão `os.encerrar`; `registro_nr12` completo quando `exige_registro_nr12` (8.4) | `observacao_encerramento`, `assinatura_anexo_id` | Consolida custos; recalcula o status do ativo (volta a `operando` se não restar outra causa) e o devolve em `status_operacional_ativo`; plano recalcula o próximo vencimento; OS fica imutável |
| `recusar-liberacao` | `concluida` | `em_execucao` | Permissão `os.encerrar`; `motivo_recusa_liberacao` | `anexos_ids` | Limpa `data_conclusao`; notifica o responsável e os executores; conta como retrabalho no first-time fix (proposta) |
| `cancelar` | `aberta`, `planejada`, `programada`, `pausada` | `cancelada` | `motivo_cancelamento` | — | Grava `data_cancelamento`; libera reservas; mantém apontamentos para custo; revoga links de uso único abertos |

Valores de `motivo_pausa`: `aguardando_peca`, `aguardando_aprovacao`, `aguardando_terceiro`, `fim_turno`, `outro`.

**Reconhecimento.** Não é estado: é o campo `data_reconhecimento`, gravado uma única vez pela ação `reconhecer` ou implicitamente ao programar, iniciar ou assumir a OS (ação `assumir`, sumário 3.6). Marcos de tempo da OS: `data_abertura`, `data_reconhecimento`, `data_inicio`, `data_conclusao`, `data_encerramento`, `data_cancelamento`.

**Assumir** — `POST /v1/ordens-servico/0193a2f0-.../assumir`, `If-Match: "1"`, corpo vazio ou `{ "observacao": "Estou no setor, atendo agora." }`. Resposta `200` com a OS (estado inalterado), `responsavel_id` igual ao usuário, `data_reconhecimento` preenchida e `transicoes_permitidas` atualizadas. Pelo app offline, é o comando `ordem_servico.assumir` (12.4); se dois técnicos assumirem a mesma OS offline, vence o primeiro a chegar ao servidor e o segundo recebe `conflito`.

**Executor fornecedor.** Em OS com `executor_tipo = fornecedor`, o `prestador` executa `reconhecer`, `assumir`, `iniciar`, `pausar`, `retomar` e `concluir`; a liberação (`encerrar` ou `recusar-liberacao`) é sempre da empresa contratante, por quem tem `os.encerrar`. O mesmo vale para OS vinculada a uma conta de subcontratada (22.8).

**Campos extras obrigatórios em OS vinculada (proposta).** Os campos extras não trafegam na conexão entre contas (sumário 3.6). Numa OS vinculada, a conclusão que chega da subcontratada não é recusada por campo extra obrigatório vazio na conta do cliente: a OS fica `concluida` com o aviso `campos_extras_pendentes`, e quem libera preenche os campos antes de `encerrar`, que então exige os obrigatórios na conclusão (`422 campo_extra_obrigatorio`). Ver 26.1, A-30.

Datas de transição (`data_inicio`, `data_pausa`, `data_conclusao`) aceitam o horário do dispositivo, para registros offline. O servidor recusa data mais de 5 minutos no futuro em relação ao recebimento e data anterior à transição anterior da mesma OS (`422 validacao_falhou`).

### 8.3 Exemplos de transição

**Iniciar** — `POST /v1/ordens-servico/0193a2f0-.../iniciar`, `If-Match: "3"`:

```json
{
  "data_inicio": "2026-11-12T13:05:00Z",
  "ativo_parado": true,
  "leituras": [
    { "medidor_id": "0192f7d1-2222-7ccc-8ddd-000000000031", "valor": "18342.6" }
  ]
}
```

Resposta `200`, `ETag: "4"`:

```json
{
  "id": "0193a2f0-5b6c-7d8e-9f01-23456789abcd",
  "numero": "OS-2026-004812",
  "versao": 4,
  "status": "em_execucao",
  "ativo_parado": true,
  "data_inicio": "2026-11-12T13:05:00Z",
  "responsavel_id": "0192f800-aaaa-7bbb-8ccc-000000000017",
  "executores_ids": ["0192f800-aaaa-7bbb-8ccc-000000000017"],
  "transicoes_permitidas": ["pausar", "concluir"]
}
```

(Resposta abreviada; a real traz o recurso completo.)

**Pausar:**

```json
{
  "motivo_pausa": "aguardando_peca",
  "observacao_pausa": "Retentor 45x62x8 sem saldo no almoxarifado central; pedido urgente aberto.",
  "data_pausa": "2026-11-12T14:20:00Z"
}
```

**Concluir:**

```json
{
  "data_conclusao": "2026-11-13T10:40:00Z",
  "servico_realizado": "Substituído retentor do eixo principal e limpa a carcaça. Teste de 30 minutos sem vazamento.",
  "codigo_falha_modo_id": "0192fb00-0000-7aaa-8bbb-000000000201",
  "codigo_falha_causa_id": "0192fb00-0000-7aaa-8bbb-000000000305",
  "codigo_falha_atividade_id": "0192fb00-0000-7aaa-8bbb-000000000410",
  "leituras": [
    { "medidor_id": "0192f7d1-2222-7ccc-8ddd-000000000031", "valor": "18343.4" }
  ],
  "anexos_ids": ["0193a3b0-1234-7abc-8def-0000000000f1"]
}
```

**Cancelar:**

```json
{
  "motivo_cancelamento": "Ativo será substituído por equipamento novo na semana 47."
}
```

**Recusar liberação** — `POST /v1/ordens-servico/0193a2f0-.../recusar-liberacao`, `If-Match: "8"`, por usuário com `os.encerrar`:

```json
{
  "motivo_recusa_liberacao": "Teste de estanqueidade não registrado. Repetir o teste com a máquina em carga e anexar a foto do manômetro.",
  "anexos_ids": []
}
```

Resposta `200`, `ETag: "9"` (trecho):

```json
{
  "id": "0193a2f0-5b6c-7d8e-9f01-23456789abcd",
  "numero": "OS-2026-004812",
  "versao": 9,
  "status": "em_execucao",
  "data_conclusao": null,
  "motivo_recusa_liberacao": "Teste de estanqueidade não registrado. Repetir o teste com a máquina em carga e anexar a foto do manômetro.",
  "transicoes_permitidas": ["pausar", "concluir"]
}
```

Webhook: `ordem_servico.status_alterado` com `dados_anteriores.status = "concluida"`.

**Transição inválida** — `POST .../concluir` numa OS `pausada`:

```json
{
  "type": "https://api.tenaz.com.br/erros/transicao_invalida",
  "title": "Ação não permitida no estado atual da OS",
  "status": 409,
  "detail": "A OS-2026-004812 está pausada. Retome a OS antes de concluir.",
  "codigo": "transicao_invalida",
  "estado_atual": "pausada",
  "transicoes_permitidas": ["retomar", "cancelar"],
  "id_rastreio": "01JCZ6P0R1S2T3V4W5X6Y7Z8A9"
}
```

**Campos faltantes** — `POST .../programar` sem responsável:

```json
{
  "type": "https://api.tenaz.com.br/erros/campos_obrigatorios_transicao",
  "title": "Faltam campos exigidos para esta ação",
  "status": 422,
  "detail": "Para programar a OS, informe a data programada e o responsável.",
  "codigo": "campos_obrigatorios_transicao",
  "campos_faltantes": ["responsavel_id"],
  "id_rastreio": "01JCZ6Q2S3T4V5W6X7Y8Z9A0B1"
}
```

### 8.4 Encerrar OS de máquina NR-12

Quando `exige_registro_nr12` é `true`, `encerrar` exige o bloco `registro_nr12` com os **sete elementos do item 12.11.2 da NR-12**. A OS não encerra sem eles (regra de negócio derivada do estudo de mercado; código definitivo no documento 4).

| Item 12.11.2 | Campo | Tipo | Regra |
|---|---|---|---|
| Intervenções realizadas | `intervencoes_realizadas[].descricao` | lista de objetos | Ao menos uma |
| Data | `intervencoes_realizadas[].data_realizacao` | data | Uma por intervenção; entre `data_inicio` e `data_conclusao` da OS |
| Serviço | `servico_realizado` | texto | Pré-preenchido com o `servico_realizado` da conclusão; mínimo de 20 caracteres (proposta) |
| Peças reparadas ou substituídas | `pecas_reparadas_substituidas[]` | lista | Pré-preenchida com as peças consumidas na OS; `acao` = `substituida` ou `reparada`. Lista vazia exige `sem_pecas: true` |
| Condições de segurança | `condicoes_seguranca` | texto | Estado das proteções, dispositivos de segurança e intertravamentos |
| Indicação conclusiva sobre a segurança da máquina | `indicacao_conclusiva` | enum | `apto`, `apto_com_restricao`, `inapto`. `apto_com_restricao` exige `restricoes` |
| Nome do responsável | `responsavel_nome` | texto | Nome de quem executou; `responsavel_usuario_id` quando é usuário da Tenaz |

`POST /v1/ordens-servico/0193a2f0-5b6c-7d8e-9f01-23456789abcd/encerrar`

```
If-Match: "9"
Idempotency-Key: 2b7c4d1e-0f9a-4c3b-8e2d-5a6f7b8c9d0e
```

```json
{
  "observacao_encerramento": "Liberado após teste funcional acompanhado pelo líder de turno.",
  "registro_nr12": {
    "intervencoes_realizadas": [
      { "descricao": "Bloqueio e etiquetagem da alimentação elétrica (ponto QGBT-03)", "data_realizacao": "2026-11-12" },
      { "descricao": "Substituição do retentor do eixo principal", "data_realizacao": "2026-11-13" },
      { "descricao": "Teste das proteções móveis, do intertravamento e da parada de emergência", "data_realizacao": "2026-11-13" }
    ],
    "servico_realizado": "Correção de vazamento de óleo com troca do retentor do eixo principal e limpeza da carcaça.",
    "pecas_reparadas_substituidas": [
      { "item_estoque_id": "0192fa10-3333-7ddd-8eee-000000000101", "descricao": "Retentor 45x62x8 NBR", "quantidade": "1", "acao": "substituida" },
      { "item_estoque_id": null, "descricao": "Proteção móvel lateral (dobradiça)", "quantidade": "1", "acao": "reparada" }
    ],
    "sem_pecas": false,
    "condicoes_seguranca": "Proteções fixas e móveis reinstaladas. Intertravamento da porta lateral e botão de emergência testados e funcionando.",
    "indicacao_conclusiva": "apto",
    "restricoes": null,
    "responsavel_nome": "Carlos Henrique Souza",
    "responsavel_usuario_id": "0192f800-aaaa-7bbb-8ccc-000000000017"
  },
  "assinatura_anexo_id": "0193a3c0-5678-7abc-8def-0000000000a2"
}
```

**Resposta `200 OK`** (`ETag: "10"`), trecho:

```json
{
  "id": "0193a2f0-5b6c-7d8e-9f01-23456789abcd",
  "numero": "OS-2026-004812",
  "versao": 10,
  "status": "encerrada",
  "data_encerramento": "2026-11-13T11:02:00Z",
  "encerrada_por_id": "0192f800-aaaa-7bbb-8ccc-000000000009",
  "status_operacional_ativo": "operando",
  "responsavel_custo": "empresa",
  "custo_mao_obra": "412.50",
  "custo_pecas": "86.90",
  "custo_terceiros": "0.00",
  "custo_total": "499.40",
  "registro_nr12": {
    "intervencoes_realizadas": [ "..." ],
    "indicacao_conclusiva": "apto",
    "responsavel_nome": "Carlos Henrique Souza",
    "registrado_em": "2026-11-13T11:02:00Z",
    "hash_sha256": "9c1f0b7e4a...e2d4",
    "url_pdf": "/v1/ordens-servico/0193a2f0-5b6c-7d8e-9f01-23456789abcd/registro-nr12?formato=pdf"
  },
  "transicoes_permitidas": []
}
```

Regras adicionais:

- `indicacao_conclusiva = inapto` é causa de `parado` (sumário 3.3): a OS encerra e o ativo fica `parado`, com `status_operacional_ativo = "parado"` na resposta. A máquina só volta a `operando` por uma nova OS encerrada com `apto` ou `apto_com_restricao` (proposta); tentativa manual de devolvê-la à operação antes disso devolve `422 nr12_maquina_inapta`.
- O registro é guardado com hash SHA-256 e passa a fazer parte do livro de manutenção da máquina, disponível em `GET /v1/ordens-servico/{id}/registro-nr12` e no inventário NR-12.
- A forma de assinatura depende da pendência P1 do sumário (ICP-Brasil no item 1.6 da NR-1). Até o parecer, a assinatura é eletrônica avançada própria (login, carimbo de tempo, geolocalização quando houver, hash e trilha de auditoria), conforme o estudo de mercado.

**Registro incompleto:**

```json
{
  "type": "https://api.tenaz.com.br/erros/nr12_registro_incompleto",
  "title": "Registro NR-12 incompleto",
  "status": 422,
  "detail": "Esta OS é de máquina sujeita à NR-12. O item 12.11.2 exige os sete elementos do registro antes do encerramento.",
  "codigo": "nr12_registro_incompleto",
  "campos_faltantes": ["registro_nr12.condicoes_seguranca", "registro_nr12.indicacao_conclusiva"],
  "erros": [
    { "campo": "registro_nr12.condicoes_seguranca", "codigo": "obrigatorio", "mensagem": "Descreva as condições de segurança da máquina." },
    { "campo": "registro_nr12.indicacao_conclusiva", "codigo": "obrigatorio", "mensagem": "Informe se a máquina está apta, apta com restrição ou inapta." }
  ],
  "id_rastreio": "01JCZ7B3C4D5E6F7G8H9J0K1L2"
}
```

### 8.5 Apontamento de horas

`POST /v1/ordens-servico/{id}/apontamentos`

```json
{
  "id": "0193a3a0-0000-7aaa-8bbb-0000000000c1",
  "usuario_id": "0192f800-aaaa-7bbb-8ccc-000000000017",
  "data_inicio": "2026-11-13T08:00:00Z",
  "data_fim": "2026-11-13T10:30:00Z",
  "tipo_hora": "normal",
  "observacao": "Troca do retentor"
}
```

- `usuario_id` diferente do próprio exige perfil de Gestão.
- `tipo_hora` (proposta): `normal`, `extra`, `deslocamento`. O custo usa o custo-hora do técnico vigente na data.
- Sobreposição com outro apontamento do mesmo técnico devolve `422 validacao_falhou` (`erros[].codigo = "sobreposicao"`).
- Apontamento de usuário `prestador` registra o tempo do terceiro (base do SLA de solução e do retrabalho), mas não entra em `custo_mao_obra`: o custo do prestador vai para `custo_terceiros`, pelo orçamento aprovado, pelo preço do contrato de prestação ou pelo valor informado no encerramento (proposta).
- Na conta de uma subcontratada (SRV), o apontamento do técnico dela é mão de obra própria (`custo_mao_obra` na conta dela). Numa OS vinculada, os apontamentos voltam para a OS do cliente quando o grupo `apontamentos_horas` está compartilhado (22.8) e, lá, seguem a regra do prestador: tempo do terceiro, sem `custo_mao_obra`. O cliente não altera nem exclui esses apontamentos (`409 os_vinculada_somente_leitura`).

### 8.6 Executor e responsável pelo custo

Toda OS tem **executor** (`executor_tipo`: `interno` ou `fornecedor`, com `fornecedor_id`) e **responsável pelo custo** (`responsavel_custo`: `empresa`, `locatario`, `locador`, `garantia` ou `fornecedor`), conforme a seção 3.3 do sumário.

| Situação | `executor_tipo` padrão | `responsavel_custo` padrão |
|---|---|---|
| Ativo próprio, fora de locação | `interno` | `empresa` |
| OS atribuída a prestador (19.3) | `fornecedor` | `empresa` (ou `garantia`, se a peça ou o serviço estiver em garantia do fornecedor) |
| Ativo `cedido` em contrato vigente | `interno` | Valor da matriz do contrato para o tipo de manutenção (21.3) |
| Ativo `tomado` em contrato vigente, manutenção a cargo do locador pela matriz | `fornecedor` (o locador) | `locador` |
| Avaria encontrada na vistoria de devolução | `interno` | `locatario` |
| Retrabalho de serviço de prestador dentro da garantia | `fornecedor` | `fornecedor` |
| OS aberta pela subcontratada conectada (`origem = prestador`, 22.8) | `fornecedor` (o da conexão) | Padrão da conta do cliente (`empresa`, ou o da matriz em ativo locado) |

O padrão pode ser trocado com justificativa. Na criação, pelo campo `justificativa_responsavel_custo`; depois, pela ação `POST /v1/ordens-servico/{id}/alterar-responsavel-custo` (perfil Gestão, `If-Match`, qualquer estado antes de `encerrada`):

```json
{
  "responsavel_custo": "empresa",
  "justificativa_responsavel_custo": "Desgaste natural do rolamento confirmado pelo técnico; pela matriz, corretiva por desgaste é da empresa, não do locatário."
}
```

Resposta `200` com a OS atualizada. A troca gera evento de auditoria com o valor anterior e a justificativa, e o webhook `ordem_servico.atualizada`. OS com `responsavel_custo = locatario` entra no relatório de cobrança do contrato (`GET /v1/contratos-locacao/{id}/cobranca`); a fatura é emitida pelo ERP.

---

## 9. Medidores: registrar leitura

`POST /v1/leituras-medidor`

```json
{
  "id": "0193a3d0-9abc-7def-8012-3456789abcde",
  "medidor_id": "0192f7d1-2222-7ccc-8ddd-000000000044",
  "valor": "412873",
  "data_leitura": "2026-11-12T10:47:00Z",
  "ordem_servico_id": null,
  "observacao": null,
  "anexos_ids": ["0193a3cf-1111-7aaa-8bbb-0000000000d4"],
  "confirmar_salto": false,
  "justificativa_salto": null
}
```

- Integrações podem trocar `medidor_id` por `ativo_id` + `tipo_medidor` (`km`, `horas`, `ciclos`, `outro`) quando o ativo tem um único medidor daquele tipo.
- `origem` é definida pelo servidor, com os valores do sumário 3.1: `manual` para usuários (inclui a leitura feita em checklist, em vistoria de locação e no abastecimento digitado); `os` quando há `ordem_servico_id`; `telematica` pela ingestão de telemática (14.1); `cartao_combustivel` pela importação do cartão (14.2); `importacao` pela importação CSV e por `POST /v1/leituras-medidor/lote` (API de integrador); `iot` pela ingestão de sensores (Fase 3, Enterprise). No lote, o integrador pode informar `telematica` ou `cartao_combustivel` quando repassa dados dessas fontes; o padrão é `importacao`.

**Resposta `201 Created`**

```json
{
  "id": "0193a3d0-9abc-7def-8012-3456789abcde",
  "versao": 1,
  "medidor_id": "0192f7d1-2222-7ccc-8ddd-000000000044",
  "ativo_id": "0192f7c4-9999-7aaa-8bbb-000000000112",
  "tipo_medidor": "km",
  "valor": "412873",
  "valor_anterior": "412311",
  "delta": "562",
  "data_leitura": "2026-11-12T10:47:00Z",
  "origem": "manual",
  "status_validacao": "valida",
  "salto_confirmado": false,
  "registrado_por_id": "0192f800-aaaa-7bbb-8ccc-000000000088",
  "planos_afetados": [
    {
      "plano_manutencao_id": "0192fc00-0000-7aaa-8bbb-000000000501",
      "nome": "Preventiva 20.000 km — cavalo 6x2",
      "proximo_vencimento_valor": "420000",
      "restante": "7127",
      "situacao": "em_dia"
    }
  ],
  "criado_em": "2026-11-12T10:47:03Z"
}
```

Regras de validação (resumo; regra completa no documento 4):

| Regra | Leitura manual | Leitura de integração |
|---|---|---|
| **Sem regressão:** valor menor que a leitura válida anterior, sem troca de medidor registrada | Recusa: `422 leitura_regressiva` | Grava com `status_validacao = suspeita`; não afeta planos até validação |
| **Teto de salto:** delta acima de `medidor.salto_maximo_por_dia` × dias desde a leitura anterior (mínimo 1 dia) | Recusa `422 leitura_salto_excessivo`, salvo com `confirmar_salto: true` e `justificativa_salto` | Grava como `suspeita` |
| **Sem futuro:** `data_leitura` mais de 5 minutos à frente do recebimento | `422 leitura_futura` | `422 leitura_futura` no item do lote |
| **Retroativa:** leitura com data anterior à última | Aceita se couber entre a leitura anterior e a posterior | Idem |
| **Duplicada:** mesmo medidor, mesma `data_leitura` e mesmo valor | Devolve a leitura existente (`200`) | Idem |
| **Troca de medidor:** via `POST /v1/medidores/{id}/trocar` | O contador recomeça; o acumulado do ativo é preservado | — |

Leituras `suspeita` aparecem numa fila para Gestão (`GET /v1/leituras-medidor?status_validacao=suspeita`), que aceita (`/validar`) ou invalida (`/invalidar`). Webhook: `leitura_medidor.registrada` (inclusive suspeitas, com o status).

**Erro de regressão:**

```json
{
  "type": "https://api.tenaz.com.br/erros/leitura_regressiva",
  "title": "Leitura menor que a anterior",
  "status": 422,
  "detail": "A leitura informada (41.287 km) é menor que a última leitura válida (412.311 km, em 10/11/2026). Confira o valor ou registre a troca do hodômetro.",
  "codigo": "leitura_regressiva",
  "ultima_leitura": { "id": "0193a1f0-0000-7aaa-8bbb-0000000000e9", "valor": "412311", "data_leitura": "2026-11-10T18:30:00Z" },
  "id_rastreio": "01JCZ7M4N5P6Q7R8S9T0V1W2X3"
}
```

---

## 10. Checklists: executar

A execução é enviada **de uma vez**, já preenchida, porque normalmente é feita offline no app. Não há edição posterior (proposta): correção é uma nova execução, e a anterior fica no histórico.

Atributos de cada item do modelo que afetam a execução: `tipo` (`sim_nao`, `numero`, `foto`, `texto`), `obrigatorio`, `resposta_esperada` (para `sim_nao`), `faixa_minima` e `faixa_maxima` (para `numero`), `exige_foto_se_reprovado`, `gera_solicitacao` e `critico_despacho` (proposta: item crítico reprovado bloqueia o despacho do veículo).

`POST /v1/checklists-execucao`

```json
{
  "id": "0193a3e0-7777-7aaa-8bbb-0000000000b1",
  "checklist_modelo_id": "0192fd00-0000-7aaa-8bbb-000000000601",
  "checklist_modelo_versao": 3,
  "ativo_id": "0192f7c4-9999-7aaa-8bbb-000000000112",
  "data_inicio": "2026-11-12T10:40:00Z",
  "data_fim": "2026-11-12T10:52:00Z",
  "posicao": { "latitude": -23.5329, "longitude": -46.7920, "precisao_m": 12 },
  "leituras": [
    { "medidor_id": "0192f7d1-2222-7ccc-8ddd-000000000044", "valor": "412873" }
  ],
  "respostas": [
    { "item_id": "0192fd01-0000-7aaa-8bbb-000000000001", "valor_booleano": true },
    { "item_id": "0192fd01-0000-7aaa-8bbb-000000000002", "valor_booleano": false, "observacao": "Lanterna traseira esquerda queimada.", "anexos_ids": ["0193a3df-2222-7aaa-8bbb-0000000000d5"] },
    { "item_id": "0192fd01-0000-7aaa-8bbb-000000000003", "valor_numero": "6.8", "observacao": "Sulco dianteiro esquerdo" },
    { "item_id": "0192fd01-0000-7aaa-8bbb-000000000004", "anexos_ids": ["0193a3df-3333-7aaa-8bbb-0000000000d6"] },
    { "item_id": "0192fd01-0000-7aaa-8bbb-000000000005", "nao_aplicavel": true },
    { "item_id": "0192fd01-0000-7aaa-8bbb-000000000006", "valor_texto": "Carga lonada e amarrada." }
  ],
  "assinatura_anexo_id": "0193a3df-4444-7aaa-8bbb-0000000000d7"
}
```

**Resposta `201 Created`**

```json
{
  "id": "0193a3e0-7777-7aaa-8bbb-0000000000b1",
  "versao": 1,
  "checklist_modelo_id": "0192fd00-0000-7aaa-8bbb-000000000601",
  "checklist_modelo_versao": 3,
  "ativo_id": "0192f7c4-9999-7aaa-8bbb-000000000112",
  "executado_por_id": "0192f800-aaaa-7bbb-8ccc-000000000088",
  "data_inicio": "2026-11-12T10:40:00Z",
  "data_fim": "2026-11-12T10:52:00Z",
  "resultado": "reprovado",
  "total_itens": 6,
  "itens_reprovados": 1,
  "itens_nao_aplicaveis": 1,
  "respostas": [
    { "item_id": "0192fd01-0000-7aaa-8bbb-000000000002", "reprovado": true, "solicitacao_id": "0193a3e1-0000-7aaa-8bbb-0000000000b2" }
  ],
  "solicitacoes_geradas": [
    { "id": "0193a3e1-0000-7aaa-8bbb-0000000000b2", "numero": "SS-2026-000732", "item_id": "0192fd01-0000-7aaa-8bbb-000000000002" }
  ],
  "leituras_registradas": [
    { "id": "0193a3e1-1111-7aaa-8bbb-0000000000b3", "medidor_id": "0192f7d1-2222-7ccc-8ddd-000000000044", "status_validacao": "valida" }
  ],
  "bloqueio_despacho": null,
  "criado_em": "2026-11-12T10:58:31Z"
}
```

Regras:

- O servidor calcula `reprovado` por item (resposta diferente da esperada ou número fora da faixa) e o `resultado` da execução: `aprovado`, `aprovado_com_ressalvas` (só itens não críticos reprovados sem `gera_solicitacao`) ou `reprovado`.
- **Item reprovado com `gera_solicitacao` gera uma solicitação** com origem `checklist`, vinculada ao item, com a observação e as fotos. Se já houver solicitação aberta para o mesmo ativo e o mesmo item do modelo, a Tenaz acrescenta um comentário em vez de duplicar (proposta).
- Item obrigatório sem resposta, ou reprovado sem foto quando `exige_foto_se_reprovado`: `422 validacao_falhou` com `erros[].campo = "respostas[1].anexos_ids"`.
- Item crítico reprovado cria bloqueio de despacho (seção 16) e devolve o objeto em `bloqueio_despacho`.
- A versão do modelo é gravada na execução. Execução feita offline numa versão já substituída é aceita se `data_inicio` for anterior à publicação da nova versão; senão, `422 checklist_modelo_desatualizado`.
- Leituras informadas passam pelas mesmas regras da seção 9.
- Webhooks: `checklist_execucao.registrada` e, se houver, `checklist_execucao.item_reprovado` e `solicitacao.criada`.

---

## 11. Indicadores: consultar KPIs

### 11.1 Catálogo

`GET /v1/indicadores` devolve o catálogo com código, nome, fórmula, "como calculamos", unidade de medida, parâmetros e referência com fonte. Toda referência exibida traz a fonte (ABRAMAN 2022 para o CMF; "referência de mercado" para os demais), conforme o estudo de mercado.

| Código | Indicador | Unidade | Plano | Fase |
|---|---|---|---|---|
| `disponibilidade_operacional` | Disponibilidade operacional (Ao) | % | E | MVP |
| `disponibilidade_inerente` | Disponibilidade inerente (Ai) | % | E | MVP |
| `disponibilidade_fisica` | Disponibilidade física | % | E | MVP |
| `mtbf` | Tempo médio entre falhas | horas | E | MVP |
| `mttr` | Tempo médio de reparo (parâmetro: com ou sem espera de peça) | horas | E | MVP |
| `mtta` | Tempo médio até o atendimento | horas | E | MVP |
| `backlog_semanas` | Backlog | semanas | E | MVP |
| `pm_compliance` | Aderência à preventiva (janela parametrizável) | % | E | MVP |
| `percentual_planejada` | % de horas planejadas | % | E | MVP |
| `custo_manutencao` | Custo de manutenção (por ativo, classe, unidade) | BRL | E | MVP |
| `mkbf` | Km médio entre falhas | km | P | F2 |
| `custo_por_km` | Custo por km (manutenção, pneus, combustível, Arla, lubrificantes) | BRL/km | P | F2 |
| `cpk_pneu` | Custo por km do pneu | BRL/km | P | F2 |
| `first_time_fix` | Resolução sem retorno | % | P | F2 |
| `wrench_time` | Tempo em execução direta | % | P | F2 |
| `cmf` | Custo de manutenção sobre faturamento | % | P | F2 |
| `custo_rav` | Custo sobre valor de reposição do ativo | % | P | F3 |
| `schedule_compliance` | Aderência à programação | % | P | F3 |

### 11.2 Consultar vários KPIs

`GET /v1/indicadores/valores?codigos=mtbf,mttr,disponibilidade_operacional,backlog_semanas,pm_compliance,percentual_planejada,custo_manutencao&inicio=2026-10-01&fim=2026-10-31&unidade_id=0192f790-0000-7aaa-8bbb-000000000001&classe_ativo_id=0192f7b0-0000-7aaa-8bbb-000000000021&comparar_periodo_anterior=true`

| Parâmetro | Descrição |
|---|---|
| `codigos` | Lista de códigos (até 20) |
| `inicio`, `fim` | Datas (`YYYY-MM-DD`), inclusivas, no fuso da unidade ou em `fuso_horario` |
| `unidade_id`, `localizacao_id`, `classe_ativo_id`, `ativo_id`, `tipo_ativo`, `criticidade` | Filtros, combináveis |
| `comparar_periodo_anterior` | Inclui a variação contra o período anterior de mesmo tamanho |
| `parametros[mttr_inclui_espera_peca]`, `parametros[pm_janela_dias]` | Sobrepõem os parâmetros da empresa só nesta consulta |

**Resposta `200 OK`**

```json
{
  "periodo": { "inicio": "2026-10-01", "fim": "2026-10-31", "fuso_horario": "America/Sao_Paulo" },
  "filtros": {
    "unidade_id": "0192f790-0000-7aaa-8bbb-000000000001",
    "classe_ativo_id": "0192f7b0-0000-7aaa-8bbb-000000000021"
  },
  "calculado_em": "2026-11-12T13:58:02Z",
  "dados_ate": "2026-11-12T13:55:00Z",
  "indicadores": [
    {
      "codigo": "mtbf",
      "nome": "Tempo médio entre falhas",
      "valor": 312.4,
      "unidade_medida": "horas",
      "situacao": "calculado",
      "amostra": { "falhas_com_parada": 9, "horas_operacao": "2811.6", "ativos": 14 },
      "variacao_percentual": 8.2,
      "referencia": null,
      "como_calculamos": "Soma das horas de operação dos ativos no período dividida pelo número de OS corretivas com o ativo parado."
    },
    {
      "codigo": "mttr",
      "nome": "Tempo médio de reparo",
      "valor": 3.6,
      "unidade_medida": "horas",
      "situacao": "calculado",
      "parametros": { "mttr_inclui_espera_peca": false },
      "amostra": { "reparos": 9 },
      "variacao_percentual": -12.5,
      "referencia": null,
      "como_calculamos": "Soma do tempo entre início e conclusão das OS corretivas, sem os intervalos pausados por espera de peça, dividida pelo número de reparos."
    },
    {
      "codigo": "disponibilidade_operacional",
      "nome": "Disponibilidade operacional",
      "valor": 96.8,
      "unidade_medida": "percentual",
      "situacao": "calculado",
      "variacao_percentual": 0.9,
      "referencia": null
    },
    {
      "codigo": "backlog_semanas",
      "nome": "Backlog",
      "valor": 2.7,
      "unidade_medida": "semanas",
      "situacao": "calculado",
      "referencia": { "texto": "2 a 4 semanas", "fonte": "Referência de mercado, não verificada" }
    },
    {
      "codigo": "pm_compliance",
      "nome": "Aderência à preventiva",
      "valor": 91.3,
      "unidade_medida": "percentual",
      "situacao": "calculado",
      "parametros": { "pm_janela_dias": 3 },
      "amostra": { "programadas": 46, "concluidas_na_janela": 42 },
      "referencia": { "texto": "90% ou mais", "fonte": "Referência de mercado" }
    },
    {
      "codigo": "percentual_planejada",
      "nome": "Horas planejadas",
      "valor": null,
      "unidade_medida": "percentual",
      "situacao": "dados_insuficientes",
      "motivo": "Menos de 20 horas apontadas no período.",
      "referencia": { "texto": "85% ou mais", "fonte": "Referência de mercado" }
    },
    {
      "codigo": "custo_manutencao",
      "nome": "Custo de manutenção",
      "valor": "48210.77",
      "unidade_medida": "BRL",
      "situacao": "calculado",
      "composicao": { "mao_obra": "21980.00", "pecas": "19410.77", "terceiros": "6820.00" },
      "variacao_percentual": 4.1
    }
  ]
}
```

- Valores monetários seguem a regra 2.5 (string decimal); os demais são números.
- `situacao` pode ser `calculado`, `dados_insuficientes` (com `motivo`) ou `nao_aplicavel`.
- Os KPIs são calculados a partir de views materializadas ou de réplica de leitura, atualizadas a cada 5 minutos (proposta); `dados_ate` informa o instante dos dados.

### 11.3 Série e agrupamento

`GET /v1/indicadores/mtbf?inicio=2026-01-01&fim=2026-10-31&granularidade=mes&agrupar_por=classe_ativo&limite_grupos=5`

```json
{
  "codigo": "mtbf",
  "unidade_medida": "horas",
  "granularidade": "mes",
  "agrupar_por": "classe_ativo",
  "series": [
    {
      "grupo": { "classe_ativo_id": "0192f7b0-0000-7aaa-8bbb-000000000021", "nome": "Compressor de parafuso" },
      "pontos": [
        { "periodo": "2026-09", "valor": 288.0, "situacao": "calculado" },
        { "periodo": "2026-10", "valor": 312.4, "situacao": "calculado" }
      ]
    }
  ]
}
```

`granularidade`: `dia`, `semana`, `mes`, `trimestre`, `ano`. `agrupar_por`: `ativo`, `classe_ativo`, `unidade`, `localizacao`, `tipo_os`, `modo_falha`. `GET /v1/indicadores/{codigo}/composicao` devolve a lista paginada das OS, leituras e apontamentos que entraram no cálculo, para auditoria do número.

---

## 12. Sincronização do app offline

### 12.1 Princípios

- O app guarda os dados do seu escopo num SQLite local e grava toda ação do usuário numa **fila de saída (outbox) de comandos**. A interface lê do SQLite, nunca espera a rede.
- O **servidor é a autoridade**. O app envia comandos; o servidor valida com as mesmas regras dos endpoints REST e devolve o resultado de cada um.
- **Ids gerados no dispositivo** (UUID v7) para tudo que é criado offline: solicitação, OS, apontamento, leitura, execução de checklist, anexo. Não há tabela de tradução de ids. Números legíveis (`OS-2026-...`, `SS-2026-...`) só são atribuídos pelo servidor; até lá o app mostra "provisório".
- Cada comando é **idempotente** pelo `id_comando`. Reenviar um lote inteiro é sempre seguro.
- Metas do sumário e do estudo de mercado: **7 dias** de operação sem conexão sem perda de registro; convergência em **menos de 60 segundos** após a reconexão (proposta).
- Ciclo: `push` dos comandos pendentes → `pull` das alterações desde o último cursor → repetir enquanto `tem_mais`. O push vem primeiro para que o pull já traga o resultado dos próprios comandos.

```
App                                   API
 |-- POST /v1/sincronizacao/comandos -->|  aplica comandos em ordem
 |<-- resultados por comando -----------|
 |-- GET /v1/sincronizacao/alteracoes ->|  cursor do último pull
 |<-- alterações + proximo_cursor ------|
 |   (repete o pull enquanto tem_mais)  |
```

### 12.2 Escopo e carga inicial

`GET /v1/sincronizacao/escopo` informa o que o usuário recebe offline e o tamanho estimado. Escopo padrão (proposta):

| Entidade | Critério |
|---|---|
| `ordem_servico` | Atribuídas ao usuário ou à sua equipe, em estado não terminal; encerradas e canceladas dos últimos 30 dias |
| `tarefa_os`, `apontamento`, peças, comentários | Das OS no escopo |
| `ativo`, `localizacao` | Das unidades do usuário, em formato resumido; para `motorista`, os veículos vinculados a ele |
| `medidor` | Dos ativos no escopo, com a última leitura válida |
| `checklist_modelo` | Versões publicadas aplicáveis às classes no escopo |
| `procedimento`, `codigo_falha` | Os usados pelas OS e classes no escopo |
| `item_estoque` e saldo | Do almoxarifado da unidade do técnico |
| `solicitacao` | As abertas pelo usuário nos últimos 30 dias |
| Documentos técnicos | Metadados; o arquivo só baixa quando marcado "disponível offline" |
| `bloqueio_ativo` e a verificação "apto para despacho" | Dos veículos do motorista (a verificação vem calculada pelo servidor, 16) |
| `cliente_servico` (conta SRV) | Resumo (nome, endereço de atendimento, contato) dos clientes das OS no escopo |
| `especialidade_tecnico` (conta SRV) | Catálogo da conta, para filtrar a fila de OS que o técnico pode assumir |
| `campo_extra_os` | Definições ativas que se aplicam aos tipos de OS e às classes de ativo no escopo, no cache de configuração (RF-OS-021); para o `prestador`, só as com `preenchido_pelo_prestador = true` |

Para o `prestador`, o escopo segue o alcance da seção 4.9: só OS com `executor_tipo = fornecedor` do seu fornecedor, os ativos dessas OS em formato resumido, seus medidores, e o saldo do almoxarifado `em_terceiro` do fornecedor. Quando a OS é retirada do fornecedor, o pull traz `remocao` com `motivo: "saiu_do_escopo"`. O `locatario` não usa o app de campo no lançamento (sumário 3.6).

Para o técnico da conta de uma subcontratada (SRV), o escopo é o de qualquer técnico na conta dela: as OS atribuídas a ele ou à sua equipe, inclusive as vinculadas a contas de clientes, os equipamentos `de_cliente` dessas OS e o estoque básico da unidade. Os campos de OS vinculada que pertencem ao cliente vêm marcados em `somente_leitura_campos`, e o app não os deixa editar offline.

`GET /v1/sincronizacao/instantaneo?cursor=<pagina>&limite=500` devolve o escopo inteiro em páginas, na mesma forma do pull. A última página traz `cursor_alteracoes`, a partir do qual começa o pull incremental. O instantâneo é usado no primeiro acesso, na troca de usuário e após `410 cursor_expirado`.

### 12.3 Pull incremental

`GET /v1/sincronizacao/alteracoes?cursor=eyJzIjo4ODEyMzQ1LCJlIjoiYzNmIn0&limite=500`

```json
{
  "alteracoes": [
    {
      "entidade": "ordem_servico",
      "id": "0193a2f0-5b6c-7d8e-9f01-23456789abcd",
      "operacao": "upsert",
      "versao": 7,
      "dados": { "id": "0193a2f0-5b6c-7d8e-9f01-23456789abcd", "numero": "OS-2026-004812", "status": "pausada", "versao": 7 }
    },
    {
      "entidade": "ordem_servico",
      "id": "0193a1aa-0000-7aaa-8bbb-000000000777",
      "operacao": "remocao",
      "motivo": "saiu_do_escopo"
    },
    {
      "entidade": "leitura_medidor",
      "id": "0193a3d0-9abc-7def-8012-3456789abcde",
      "operacao": "upsert",
      "versao": 1,
      "dados": { "id": "0193a3d0-9abc-7def-8012-3456789abcde", "valor": "412873", "status_validacao": "valida" }
    }
  ],
  "proximo_cursor": "eyJzIjo4ODEyNDAyLCJlIjoiYzNmIn0",
  "tem_mais": false,
  "servidor_em": "2026-11-12T14:01:22Z"
}
```

- O cursor é opaco. Internamente codifica a posição num registro de alterações por empresa, com número de sequência monotônico gravado na mesma transação da alteração (proposta; detalhe no documento 5), mais um resumo do escopo.
- `operacao`: `upsert` (objeto completo no formato da API) ou `remocao` com `motivo` (`excluido`, `saiu_do_escopo`). Quando algo entra no escopo (OS atribuída ao técnico), o pull traz o objeto e seus filhos como `upsert`.
- O app aplica as alterações em ordem. Se houver comandos locais pendentes sobre o mesmo registro, o app reaplica esses comandos sobre a versão do servidor para exibição, até o push confirmar (proposta).
- **Retenção do cursor: 30 dias** (proposta), bem acima dos 7 dias de meta offline. Cursor mais antigo: `410 cursor_expirado`, e o app refaz o instantâneo **depois** de enviar os comandos pendentes.
- Espera longa opcional: `aguardar=25` mantém a requisição aberta até 25 segundos se não houver alterações (proposta). O servidor também envia push silencioso de "há alterações" ao dispositivo.

### 12.4 Push de comandos

`POST /v1/sincronizacao/comandos`

```
Authorization: Bearer eyJhbGciOiJFZERTQSIs...
X-Tenaz-Dispositivo-Id: 0192ff00-dddd-7eee-8fff-000000000005
X-Tenaz-App-Versao: 1.4.2
Content-Type: application/json
```

```json
{
  "enviado_em_dispositivo": "2026-11-12T14:00:58Z",
  "comandos": [
    {
      "id_comando": "0193a400-0000-7aaa-8bbb-000000000001",
      "tipo": "solicitacao.criar",
      "criado_em_dispositivo": "2026-11-10T07:12:00Z",
      "alvo": { "entidade": "solicitacao", "id": "0193a3ff-1111-7aaa-8bbb-0000000000a1" },
      "dados": { "ativo_id": "0192f7c4-9999-7aaa-8bbb-000000000112", "descricao": "Ruído no diferencial em subida." }
    },
    {
      "id_comando": "0193a400-0000-7aaa-8bbb-000000000002",
      "tipo": "ordem_servico.iniciar",
      "criado_em_dispositivo": "2026-11-10T08:30:00Z",
      "alvo": { "entidade": "ordem_servico", "id": "0193a2f0-5b6c-7d8e-9f01-23456789abcd", "versao_base": 3 },
      "dados": { "data_inicio": "2026-11-10T08:30:00Z", "ativo_parado": true }
    },
    {
      "id_comando": "0193a400-0000-7aaa-8bbb-000000000003",
      "tipo": "apontamento.criar",
      "criado_em_dispositivo": "2026-11-10T11:00:00Z",
      "alvo": { "entidade": "apontamento", "id": "0193a3ff-2222-7aaa-8bbb-0000000000a2" },
      "dados": { "ordem_servico_id": "0193a2f0-5b6c-7d8e-9f01-23456789abcd", "data_inicio": "2026-11-10T08:30:00Z", "data_fim": "2026-11-10T11:00:00Z", "tipo_hora": "normal" }
    },
    {
      "id_comando": "0193a400-0000-7aaa-8bbb-000000000004",
      "tipo": "leitura_medidor.criar",
      "criado_em_dispositivo": "2026-11-10T08:29:00Z",
      "alvo": { "entidade": "leitura_medidor", "id": "0193a3ff-3333-7aaa-8bbb-0000000000a3" },
      "dados": { "medidor_id": "0192f7d1-2222-7ccc-8ddd-000000000031", "valor": "18342.6", "data_leitura": "2026-11-10T08:29:00Z" }
    },
    {
      "id_comando": "0193a400-0000-7aaa-8bbb-000000000005",
      "tipo": "ordem_servico.atualizar",
      "criado_em_dispositivo": "2026-11-10T11:05:00Z",
      "alvo": { "entidade": "ordem_servico", "id": "0193a2f0-5b6c-7d8e-9f01-23456789abcd", "versao_base": 3 },
      "dados": { "descricao": "Vazamento no retentor; carcaça trincada também." }
    },
    {
      "id_comando": "0193a400-0000-7aaa-8bbb-000000000006",
      "tipo": "ordem_servico.concluir",
      "criado_em_dispositivo": "2026-11-10T11:10:00Z",
      "alvo": { "entidade": "ordem_servico", "id": "0193a2f0-5b6c-7d8e-9f01-23456789abcd", "versao_base": 3 },
      "dados": { "data_conclusao": "2026-11-10T11:10:00Z", "servico_realizado": "Troca do retentor e da carcaça." }
    }
  ]
}
```

- Até **100 comandos ou 1 MB** por lote (proposta). O app envia em ordem de criação.
- Cada comando é processado na sua própria transação, na ordem do lote. A falha de um comando não desfaz os anteriores. Comando que depende de outro rejeitado no mesmo lote (mesmo `alvo.id` ou referência a um id criado por ele) recebe `rejeitado` com `codigo: "dependencia_rejeitada"`.
- O `id_comando` fica guardado por **90 dias** (proposta). Reenvio devolve o mesmo resultado, com `status: "ja_aplicado"`.
- Erros que afetam o lote inteiro (`401`, `429`, `5xx`) não consomem comandos: o app repete o mesmo lote depois, com espera exponencial.
- Anexos vão antes, por `POST /v1/anexos` com o id gerado no dispositivo. Comando que referencia anexo ainda não enviado é aplicado e devolve `anexos_pendentes`; o vínculo se completa quando o upload é confirmado.

**Tipos de comando (MVP):** `solicitacao.criar`, `solicitacao.comentar`, `ordem_servico.criar`, `ordem_servico.atualizar`, `ordem_servico.reconhecer`, `ordem_servico.assumir`, `ordem_servico.iniciar`, `servico_os.atualizar`, `ordem_servico.pausar`, `ordem_servico.retomar`, `ordem_servico.concluir`, `tarefa_os.atualizar`, `apontamento.criar`, `apontamento.atualizar`, `peca_os.consumir`, `comentario_os.criar`, `assinatura_os.criar`, `leitura_medidor.criar`, `checklist_execucao.criar`, `bloqueio_energia.aplicar`, `bloqueio_energia.retirar`. **Fase 2:** `abastecimento.criar`, `pneu.aferir`, `sinistro.criar`, `peca_serializada.instalar`, `peca_serializada.remover`, `vistoria_locacao.criar`. Cada tipo usa o mesmo validador do endpoint REST equivalente e as mesmas regras de alcance (o `prestador` só envia comandos sobre OS no seu alcance). `encerrar`, `recusar-liberacao` e `cancelar` exigem conexão (proposta): são decisões de gestão sobre o estado consolidado. Em OS vinculada, comando sobre campo do outro lado volta `rejeitado` com `os_vinculada_somente_leitura`; a sincronização com a conta do cliente só começa depois que o comando é aplicado no servidor (22.6).

**Campos extras offline.** `ordem_servico.criar`, `ordem_servico.atualizar` e `ordem_servico.concluir` aceitam `campos_extras` e `causa_externa`. O app valida os campos obrigatórios com as definições do cache antes de enfileirar o comando (RF-OS-021). No servidor, a obrigatoriedade é avaliada pelas definições vigentes no `criado_em_dispositivo` do comando (proposta): campo que passou a ser obrigatório depois que o técnico concluiu offline não rejeita a conclusão: o comando volta `aplicado` com o aviso `campos_extras_pendentes`, e a OS fica sinalizada para revisão (RF-OS-021). Campo que já era obrigatório quando o comando foi criado e chegou vazio volta `rejeitado` com `campo_extra_obrigatorio`.

**Resposta `200 OK`**

```json
{
  "resultados": [
    { "id_comando": "0193a400-0000-7aaa-8bbb-000000000001", "status": "aplicado", "entidade": "solicitacao", "id": "0193a3ff-1111-7aaa-8bbb-0000000000a1", "versao": 1, "dados_servidor": { "numero": "SS-2026-000733" } },
    { "id_comando": "0193a400-0000-7aaa-8bbb-000000000002", "status": "aplicado", "entidade": "ordem_servico", "id": "0193a2f0-5b6c-7d8e-9f01-23456789abcd", "versao": 4 },
    { "id_comando": "0193a400-0000-7aaa-8bbb-000000000003", "status": "ja_aplicado", "entidade": "apontamento", "id": "0193a3ff-2222-7aaa-8bbb-0000000000a2", "versao": 1 },
    { "id_comando": "0193a400-0000-7aaa-8bbb-000000000004", "status": "aplicado", "entidade": "leitura_medidor", "id": "0193a3ff-3333-7aaa-8bbb-0000000000a3", "versao": 1, "avisos": [ { "codigo": "leitura_suspeita", "mensagem": "Salto acima do teto; aguardando validação do gestor." } ] },
    { "id_comando": "0193a400-0000-7aaa-8bbb-000000000005", "status": "aplicado_com_mesclagem", "entidade": "ordem_servico", "id": "0193a2f0-5b6c-7d8e-9f01-23456789abcd", "versao": 5, "campos_descartados": [ { "campo": "descricao", "valor_servidor": "Vazamento no retentor. Ver foto da carcaça.", "alterado_por_id": "0192f800-aaaa-7bbb-8ccc-000000000009", "alterado_em": "2026-11-10T12:40:00Z" } ] },
    {
      "id_comando": "0193a400-0000-7aaa-8bbb-000000000006",
      "status": "conflito",
      "entidade": "ordem_servico",
      "id": "0193a2f0-5b6c-7d8e-9f01-23456789abcd",
      "erro": {
        "type": "https://api.tenaz.com.br/erros/transicao_invalida",
        "title": "Ação não permitida no estado atual da OS",
        "status": 409,
        "codigo": "transicao_invalida",
        "detail": "A OS-2026-004812 foi cancelada por Ana Lima em 10/11/2026 às 10:02. A conclusão feita offline não foi aplicada; as horas apontadas foram mantidas.",
        "estado_atual": "cancelada",
        "transicoes_permitidas": []
      },
      "estado_servidor": { "status": "cancelada", "versao": 5 }
    }
  ],
  "servidor_em": "2026-11-12T14:01:01Z",
  "deslocamento_relogio_s": -3
}
```

| `status` | Significado | O que o app faz |
|---|---|---|
| `aplicado` | Gravado | Remove da fila; atualiza o registro local |
| `ja_aplicado` | Repetição de comando já processado | Remove da fila |
| `aplicado_com_mesclagem` | Gravado em parte; campos alterados no servidor depois venceram | Remove da fila; mostra os `campos_descartados` ao usuário |
| `conflito` | O estado do servidor impede o comando | Remove da fila; mostra o conflito e o estado atual; o usuário decide (por exemplo, abrir nova solicitação) |
| `rejeitado` | Erro de validação que não muda com nova tentativa | Remove da fila; mostra o erro para correção e reenvio como comando novo |

### 12.5 Regras de conflito

| Tipo de dado | Estratégia |
|---|---|
| **Só acréscimo** (apontamento, leitura, comentário, foto, execução de checklist, assinatura) | Não conflita. Só passa pela validação de domínio (por exemplo, leitura regressiva vinda do app vira `suspeita` em vez de ser recusada, para não perder o registro de campo, proposta) |
| **Edição de campos** (descrição, observação, tarefa marcada, cada chave de `campos_extras`) | Mesclagem por campo contra `versao_base` (em `campos_extras`, cada chave conta como um campo, proposta): campo não alterado no servidor desde `versao_base` é aplicado; campo alterado nos dois lados segue "último a escrever vence" por campo, comparando o horário do dispositivo (corrigido pelo deslocamento de relógio) com o horário da alteração no servidor. O lado perdedor volta em `campos_descartados` |
| **Transição de estado** | Validada contra o estado atual, não contra `versao_base`. Se o estado atual permitir, aplica. Se o estado já é o pretendido (dois técnicos iniciaram a mesma OS), aplica sem nova transição e inclui o usuário em `executores_ids` |
| **OS cancelada ou encerrada no servidor** | Transições e edições: `conflito`. Apontamentos e leituras: aplicados, porque registram fatos, e a gestão é notificada (proposta) |
| **OS reatribuída a outro técnico durante o offline** | Aplica, com aviso `os_reatribuida`, e notifica a gestão |
| **Dois técnicos assumem a mesma OS offline** | O primeiro `ordem_servico.assumir` aplicado vence; o segundo recebe `conflito` com `transicao_invalida` e o responsável atual |
| **OS vinculada encerrada, cancelada ou retirada pelo cliente durante o offline** (conta SRV) | Mesma regra da OS cancelada ou encerrada: transições e edições em `conflito`; apontamentos, componentes e fotos são aplicados na conta da subcontratada e só sincronizam se o vínculo ainda aceitar (senão ficam locais, com aviso `vinculo_fechado`) |
| **OS retirada do fornecedor durante o offline** (`prestador`) | Transições e edições: `rejeitado` com `recurso_nao_encontrado` (a OS saiu do alcance). Apontamentos, leituras e fotos: guardados numa fila de revisão da gestão, sem aplicar à OS, e a gestão é notificada (proposta) |
| **Consumo de peça** | Aplica mesmo que o saldo fique negativo, porque a peça já foi usada; gera alerta `estoque_negativo` para o almoxarife (proposta). Pela API online, saldo insuficiente continua devolvendo `422 estoque_insuficiente` |
| **Ativo inativado** | Leitura e execução de checklist: `rejeitado` com `ativo_inativo`; apontamento em OS existente: aplicado |
| **Versão de checklist substituída** | Regra da seção 10 |

**Relógio do dispositivo.** O servidor calcula o deslocamento entre `enviado_em_dispositivo` e o horário de recebimento. Até 2 minutos, os horários são usados como vieram. Acima de 10 minutos, o servidor corrige os horários do lote pelo deslocamento e marca os registros com `horario_corrigido=true` (proposta). Os dois valores, original e corrigido, ficam na auditoria.

---

## 13. Webhooks

Disponíveis no plano Profissional, desde o MVP (sumário 2.1 e 4).

### 13.1 Criar assinatura

`POST /v1/webhooks`

```json
{
  "url": "https://integracao.cliente.com.br/tenaz/webhooks",
  "descricao": "ERP — custos de OS e documentos",
  "eventos": ["ordem_servico.encerrada", "documento_vencimento.vencendo", "documento_vencimento.vencido", "leitura_medidor.registrada"],
  "filtros": { "unidade_ids": ["0192f790-0000-7aaa-8bbb-000000000001"] },
  "versao_evento": 1,
  "ativo": true
}
```

Resposta `201` com o `segredo` (`whsec_...`, 32 bytes aleatórios em base64url), exibido **só nesta resposta**. `eventos: ["*"]` assina todos os eventos, inclusive os criados no futuro.

Regras da URL: só `https`, porta 443 ou 8443; o host não pode resolver para IP privado, de loopback ou de link local (proteção contra SSRF), verificado a cada entrega; redirecionamentos não são seguidos. A Tenaz publica os IPs de saída dos webhooks (proposta).

### 13.2 Catálogo de eventos

Regra de nome (sumário 3.4): `<entidade>.<fato_no_particípio>`, sempre com o **nome completo** da entidade, sem abreviações: `ordem_servico.encerrada`, `documento_vencimento.vencendo`, `leitura_medidor.registrada`.

| Evento | Quando | Plano | Fase |
|---|---|---|---|
| `ativo.criado` | Ativo cadastrado | P | MVP |
| `ativo.atualizado` | Dados do ativo alterados | P | MVP |
| `ativo.status_alterado` | Mudança de `status_operacional` | P | MVP |
| `ativo.colocado_em_operacao` | Ativo passou de `em_implantacao` para `em_operacao` (com `data_entrada_operacao` e `projeto_id`, quando houver); um evento por ativo, inclusive no lote e na conclusão de projeto (RF-ATV-013) | P | MVP |
| `leitura_medidor.registrada` | Leitura de medidor gravada (inclui `origem` e `status_validacao`) | P | MVP |
| `solicitacao.criada` | Solicitação aberta | P | MVP |
| `solicitacao.aprovada` | Solicitação aprovada (com `ordem_servico_id`) | P | MVP |
| `solicitacao.recusada` | Solicitação recusada (inclusive por desistência) ou marcada duplicada | P | MVP |
| `ordem_servico.criada` | OS criada, com `origem` (`manual`, `solicitacao`, `plano`, `checklist`, `telematica`, `rascunho_ia`, `vistoria_locacao`, `prestador`, `cliente_conectado`, `projeto`) | P | MVP |
| `ordem_servico.atualizada` | Campos da OS alterados sem transição, inclusive `assumir`, atribuição ou retirada de fornecedor, troca de `responsavel_custo` e dados recebidos da conta conectada | P | MVP |
| `ordem_servico.status_alterado` | Qualquer transição, inclusive `desprogramar` e `recusar-liberacao` (com `status_anterior` e motivo) | P | MVP |
| `ordem_servico.programada` | OS programada | P | MVP |
| `ordem_servico.iniciada` | OS em execução | P | MVP |
| `ordem_servico.pausada` | OS pausada (com motivo) | P | MVP |
| `ordem_servico.concluida` | Trabalho técnico concluído | P | MVP |
| `ordem_servico.encerrada` | OS encerrada, com custos consolidados e registro NR-12 quando houver | P | MVP |
| `ordem_servico.cancelada` | OS cancelada (com motivo) | P | MVP |
| `plano_manutencao.vencendo` | Plano a vencer dentro da antecedência configurada | P | MVP |
| `checklist_execucao.registrada` | Checklist executado | P | MVP |
| `checklist_execucao.item_reprovado` | Execução com item reprovado | P | MVP |
| `item_estoque.abaixo_minimo` | Saldo disponível abaixo do mínimo | P | MVP |
| `movimentacao_estoque.registrada` | Movimentação de qualquer `tipo` (20.1), com `nfe_numero` e `nfe_chave` quando houver; no MVP, `entrada`, `saida_os`, `devolucao_os` e `ajuste` | P | MVP |
| `documento_vencimento.vencendo` | Documento a vencer em D-60, D-30 e D-7 (um evento por marco) | P | MVP |
| `documento_vencimento.vencido` | Documento venceu | P | MVP |
| `ativo.despacho_bloqueado` | A verificação "apto para despacho" do veículo passou a `false` (com os motivos, seção 16) | P | MVP |
| `ativo.despacho_liberado` | O veículo voltou a ficar apto para despacho | P | MVP |
| `conexao_empresa.ativada` | Conexão entre contas aceita (ou reativada depois de suspensa); emitido nas duas contas | P | MVP |
| `conexao_empresa.encerrada` | Conexão recusada ou encerrada; a sincronização parou; emitido nas duas contas | P | MVP |
| `vinculo_os.sincronizado` | Uma alteração da OS vinculada foi gravada na outra conta (criação, transição, apontamentos, componentes, fotos); emitido na conta que recebeu os dados | P | MVP |
| `vinculo_os.erro` | A sincronização de um vínculo falhou e esgotou as novas tentativas (`estado_sincronizacao = erro`); emitido na conta de origem da alteração | P | MVP |
| `regra_automacao.disparada` | **(proposta)** Uma regra configurável disparou e aplicou o efeito; `dados` traz o `disparo_regra` (regra, condição, fato, ativo, efeito e os ids criados na própria conta: solicitação ou bloqueio). Um evento por disparo; a rotina diária não repete disparo do mesmo fato (23.1) | P | MVP |
| `importacao.concluida` | Importação terminou (com resumo) | P | MVP |
| `exportacao.concluida` | Exportação completa disponível | P | MVP |
| `webhook.teste` | Disparado por `POST /v1/webhooks/{id}/testar` | P | MVP |
| `pedido_compra.aprovado` | Pedido de compra aprovado | P | F2 |
| `nota_fiscal_entrada.confirmada` | Entrada de NF-e confirmada | P | F2 |
| `capacitacao.vencendo` | Capacitação a vencer | P | F2 |
| `capacitacao.vencida` | Capacitação vencida | P | F2 |
| `abastecimento.registrado` | Abastecimento gravado | P | F2 |
| `abastecimento.anomalia_detectada` | Abastecimento com anomalia | P | F2 |
| `pneu.movimentado` | Montagem, desmontagem, recapagem ou sucata | P | F2 |
| `multa.registrada` | Multa cadastrada | P | F2 |
| `multa.prazo_indicacao_proximo` | Faltam 7 dias para o fim do prazo de 30 dias de indicação | P | F2 |
| `sinistro.registrado` | Sinistro cadastrado | P | F2 |
| `telematica.dtc_recebido` | DTC recebido e casado com um ativo | P | F2 |
| `ordem_reparo.enviada` | Ordem de reparo enviada ao fornecedor, com NF-e de remessa | P | F2 |
| `ordem_reparo.retornada` | Peça voltou do conserto, recondicionada, com custo e garantia | P | F2 |
| `contrato_locacao.ativado` | Contrato de locação passou a `vigente` | P | F2 |
| `contrato_locacao.encerrado` | Contrato de locação encerrado | P | F2 |
| `vistoria_locacao.registrada` | Vistoria de entrega ou devolução registrada (com avarias e OS geradas) | P | F2 |
| `fornecedor.sla_violado` | OS ou substituição de ativo fora do SLA de atendimento, de solução ou de substituição do contrato (prestação ou locação tomada) | P | F2 |
| `projeto.iniciado` | Projeto passou a `em_andamento` (RF-PRJ-002) | P | F2 |
| `projeto.concluido` | Projeto concluído, com os ativos colocados em operação (RF-PRJ-007) | P | F2 |
| `projeto.cancelado` | Projeto cancelado, com motivo e o destino de cada ativo `em_implantacao` | P | F2 |
| `contrato_cliente.ativado` | Contrato com cliente passou a `vigente` (RF-CTR-001) | P | F2 |
| `contrato_cliente.encerrado` | Contrato com cliente encerrado | P | F2 |
| `contrato_cliente.sla_violado` | OS de item coberto fora do prazo de atendimento ou de solução do contrato com cliente, com o prazo e o atraso | P | F2 |
| `medicao_contrato.fechada` | Medição do contrato fechada, com o total, o desconto e o link do relatório (RF-CTR-004) | P | F2 |

Só usuários internos da própria conta assinam webhooks (administrador ou token com `webhooks:gerenciar`); usuários `prestador` e `locatario` não recebem eventos. A conta de uma subcontratada (SRV) assina os eventos **da própria conta**, se o plano incluir webhooks (22.10); ela nunca recebe eventos da conta do cliente, e vice-versa. Os eventos de conexão e de vínculo trazem só dados da conta que os recebe: o número da OS do outro lado (`numero_os_remota`) e a razão social da outra empresa, nunca ids de outra conta.

Exemplo de `vinculo_os.sincronizado` recebido pela conta do cliente:

```json
{
  "id": "evt_0193d5a0b1c27f3a8b4c5d6e7f809a22",
  "tipo": "vinculo_os.sincronizado",
  "versao_evento": 1,
  "ocorrido_em": "2026-11-13T15:10:04Z",
  "empresa_id": "0192f700-0000-7aaa-8bbb-000000000001",
  "ambiente": "producao",
  "autor": { "tipo": "sincronizacao_conexao", "id": "0193d000-0000-7aaa-8bbb-000000000c01" },
  "dados": {
    "vinculo_os": {
      "id": "0193d100-0000-7aaa-8bbb-000000000d01",
      "conexao_empresa_id": "0193d000-0000-7aaa-8bbb-000000000c01",
      "ordem_servico_id": "0193d200-0000-7aaa-8bbb-000000000e01",
      "numero_os": "OS-2026-004950",
      "numero_os_remota": "OS-2026-000231",
      "lado_origem": "cliente",
      "estado_sincronizacao": "sincronizado",
      "grupos_recebidos": ["os_execucao", "apontamentos_horas", "componentes", "fotos_anexos"],
      "sincronizado_em": "2026-11-13T15:10:04Z"
    }
  }
}
```

### 13.3 Envelope

```
POST /tenaz/webhooks HTTP/1.1
Host: integracao.cliente.com.br
Content-Type: application/json; charset=utf-8
User-Agent: Tenaz-Webhooks/1.0
Tenaz-Webhook-Id: evt_0193a5c0d1e27f3a8b4c5d6e7f809a1b
Tenaz-Webhook-Timestamp: 1794489600
Tenaz-Webhook-Signature: v1=3b9f1c0e7d2a4f58b6c1e9a07d3f2b4c8e5a1d6f9b0c2e4a7d8f1b3c5e7a9d0f
Tenaz-Webhook-Tentativa: 1
```

```json
{
  "id": "evt_0193a5c0d1e27f3a8b4c5d6e7f809a1b",
  "tipo": "ordem_servico.encerrada",
  "versao_evento": 1,
  "ocorrido_em": "2026-11-13T11:02:00Z",
  "empresa_id": "0192f700-0000-7aaa-8bbb-000000000001",
  "ambiente": "producao",
  "autor": { "tipo": "usuario", "id": "0192f800-aaaa-7bbb-8ccc-000000000009" },
  "dados": {
    "ordem_servico": {
      "id": "0193a2f0-5b6c-7d8e-9f01-23456789abcd",
      "numero": "OS-2026-004812",
      "versao": 10,
      "tipo": "corretiva_planejada",
      "status": "encerrada",
      "ativo_id": "0192f7c4-1a2b-7c3d-8e4f-5a6b7c8d9e01",
      "unidade_id": "0192f790-0000-7aaa-8bbb-000000000001",
      "executor_tipo": "interno",
      "fornecedor_id": null,
      "responsavel_custo": "empresa",
      "contrato_locacao_id": null,
      "data_encerramento": "2026-11-13T11:02:00Z",
      "custo_mao_obra": "412.50",
      "custo_pecas": "86.90",
      "custo_terceiros": "0.00",
      "custo_total": "499.40"
    }
  },
  "dados_anteriores": { "status": "concluida" }
}
```

- `dados` traz o **retrato completo** do recurso no momento do evento, no mesmo formato da API em `versao_evento`. Para ordenar eventos do mesmo recurso, use `dados.<recurso>.versao`.
- `dados_anteriores` traz só os campos que mudaram, nos eventos de alteração.
- `autor.tipo`: `usuario`, `token_api`, `link_execucao` (prestador ocasional, com o `id` do link), `sincronizacao_conexao` (dado gravado pelo processo de sincronização entre contas, com o `id` da conexão; o usuário da outra conta aparece só pelo nome, em `autor.nome_remoto`), `sistema` (plano, regra de telemática, rotina de vencimentos), `agente_ia` (Fase 3).
- Corpo máximo de 256 KB. Recurso maior vai sem os filhos, com `dados_truncados: true`; o cliente busca o restante pela API.
- **Entrega pelo menos uma vez, sem garantia de ordem.** O receptor deve deduplicar por `id` (igual a `Tenaz-Webhook-Id`).

### 13.4 Assinatura HMAC-SHA256 com timestamp

1. A Tenaz monta a string `"{Tenaz-Webhook-Timestamp}.{corpo bruto}"`, com o timestamp em segundos Unix e o corpo exatamente como enviado.
2. Calcula `HMAC-SHA256(segredo, string)`, com o segredo completo (incluindo o prefixo `whsec_`) como chave em UTF-8, e codifica em hexadecimal minúsculo.
3. Envia em `Tenaz-Webhook-Signature: v1=<hex>`. Durante a rotação de segredo (24 horas), o cabeçalho traz duas assinaturas separadas por vírgula: `v1=<nova>,v1=<antiga>`.

O receptor deve:

- recalcular a assinatura sobre o **corpo bruto**, antes de qualquer parse de JSON;
- comparar em tempo constante com cada `v1=` do cabeçalho;
- recusar se `|agora − timestamp| > 300 segundos` (proteção contra repetição);
- deduplicar por `Tenaz-Webhook-Id`;
- responder `2xx` em até 10 segundos e processar depois, de forma assíncrona.

Verificação em Go:

```go
func VerificarWebhook(segredo string, cabecalhos http.Header, corpo []byte, agora time.Time) error {
	ts := cabecalhos.Get("Tenaz-Webhook-Timestamp")
	seg, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return errors.New("timestamp inválido")
	}
	if d := agora.Unix() - seg; d > 300 || d < -300 {
		return errors.New("timestamp fora da tolerância")
	}
	mac := hmac.New(sha256.New, []byte(segredo))
	mac.Write([]byte(ts + "."))
	mac.Write(corpo)
	esperada := hex.EncodeToString(mac.Sum(nil))
	for _, parte := range strings.Split(cabecalhos.Get("Tenaz-Webhook-Signature"), ",") {
		v := strings.TrimPrefix(strings.TrimSpace(parte), "v1=")
		if hmac.Equal([]byte(v), []byte(esperada)) {
			return nil
		}
	}
	return errors.New("assinatura inválida")
}
```

Verificação em TypeScript (Node.js):

```ts
import { createHmac, timingSafeEqual } from "node:crypto";

export function verificarWebhook(segredo: string, cabecalhos: Record<string, string>, corpoBruto: Buffer): boolean {
  const ts = cabecalhos["tenaz-webhook-timestamp"];
  if (!ts || Math.abs(Date.now() / 1000 - Number(ts)) > 300) return false;
  const esperada = createHmac("sha256", segredo).update(`${ts}.`).update(corpoBruto).digest("hex");
  return (cabecalhos["tenaz-webhook-signature"] ?? "")
    .split(",")
    .map((p) => p.trim().replace(/^v1=/, ""))
    .some((v) => v.length === esperada.length && timingSafeEqual(Buffer.from(v), Buffer.from(esperada)));
}
```

### 13.5 Entrega e novas tentativas

- Sucesso: qualquer `2xx` em até **10 segundos**. Qualquer outro resultado (3xx, 4xx, 5xx, timeout, erro de TLS ou DNS) é falha.
- Exceção: `410 Gone` desativa a assinatura na hora e notifica o administrador (proposta).
- Novas tentativas com **espera exponencial e variação de ±20%, por até 24 horas** desde a primeira tentativa (meta do estudo de mercado):

| Tentativa | Espera após a anterior | Tempo acumulado aproximado |
|---|---|---|
| 1 | — | 0 |
| 2 | 1 min | 1 min |
| 3 | 5 min | 6 min |
| 4 | 15 min | 21 min |
| 5 | 30 min | 51 min |
| 6 | 1 h | 1 h 51 min |
| 7 | 2 h | 3 h 51 min |
| 8 | 4 h | 7 h 51 min |
| 9 | 8 h | 15 h 51 min |
| 10 | 8 h | 23 h 51 min |

- Cada tentativa gera novo `Tenaz-Webhook-Timestamp` e nova assinatura; `Tenaz-Webhook-Id` e corpo não mudam.
- Respostas `429` com `Retry-After` do receptor são respeitadas, dentro da janela de 24 horas.
- Entregas para o mesmo endpoint são limitadas a 10 simultâneas (proposta), para não sobrecarregar o receptor.

### 13.6 DLQ, suspensão e painel de entregas

- Depois da última tentativa, a entrega vai para a **DLQ** (`status = dlq`), onde fica **30 dias** (proposta). O administrador recebe e-mail com o resumo diário.
- Se todas as entregas de uma assinatura caírem na DLQ por **3 dias seguidos**, a assinatura passa a `suspensa` (proposta). Eventos novos continuam guardados com `status = retido` por 30 dias e saem quando a assinatura for reativada (`POST /v1/webhooks/{id}/reativar`).
- Painel: `GET /v1/webhooks/{id}` traz os indicadores de entrega (`taxa_sucesso_24h`, `latencia_p95_ms`, `entregas_em_dlq`, `ultima_falha_em`). `GET /v1/webhooks/{id}/entregas?status=dlq,falhando&inicio=...` lista as entregas.

```json
{
  "dados": [
    {
      "id": "0193a6aa-0000-7aaa-8bbb-000000000901",
      "evento_id": "evt_0193a5c0d1e27f3a8b4c5d6e7f809a1b",
      "tipo": "ordem_servico.encerrada",
      "status": "dlq",
      "tentativas": 10,
      "primeira_tentativa_em": "2026-11-13T11:02:01Z",
      "ultima_tentativa_em": "2026-11-14T10:53:12Z",
      "proxima_tentativa_em": null,
      "ultima_resposta": { "status_http": 503, "duracao_ms": 10000, "erro": "timeout", "corpo_inicio": "" }
    }
  ],
  "paginacao": { "proximo_cursor": null, "tem_mais": false, "limite": 50 }
}
```

Estados de entrega: `pendente`, `entregue`, `falhando` (em novas tentativas), `dlq`, `retido`, `descartado` (assinatura excluída).

### 13.7 Reentrega manual

- `POST /v1/entregas-webhook/{id}/reentregar` → `202`. Reenvia uma entrega em qualquer estado, com o **mesmo** `Tenaz-Webhook-Id` (o receptor deduplica) e assinatura nova. Conta como nova série de tentativas de 24 horas.
- `POST /v1/webhooks/{id}/reentregar-dlq` com `{ "inicio": "...", "fim": "...", "tipos": ["ordem_servico.encerrada"] }` → `202` com `Location` de processamento. Reentrega em lote, na ordem de `ocorrido_em`, respeitando o limite de simultaneidade.

---

## 14. Integrações de entrada

### 14.1 Telemática

Fase 2, plano Profissional (conectores de telemática na tabela de recursos do sumário 2.1). Sensores IoT (leituras com origem `iot`) são da Fase 3, Enterprise, e usarão o mesmo modelo canônico.

**Arquitetura.** Três caminhos chegam ao mesmo modelo canônico de evento:

| Caminho | Como funciona | Provedores |
|---|---|---|
| Conector da Tenaz (coleta) | Um worker consulta a API do provedor periodicamente (proposta: a cada 15 minutos para leituras) com as credenciais da `conexao_telematica` | rFMS Volvo e Scania (padrão ACEA rFMS), Geotab, Cobli; depois Ticket Log, Sascar, Omnilink (APIs não verificadas) |
| Receptor de webhook do provedor | `POST /v1/telematica/receptor/{conexao_id}`, verificado com a assinatura do provedor | Provedores que enviam alertas de falha por webhook |
| Ingestão genérica | `POST /v1/telematica/eventos` com token de API e escopo `telematica:ingestao` | Rastreador próprio do cliente, gateway, integrador |

**Modelo canônico.** `POST /v1/telematica/eventos` aceita até 500 eventos por requisição:

```json
{
  "fonte": "rastreador_proprio",
  "eventos": [
    {
      "id_externo": "trk-88213-20261112T1350",
      "tipo": "dtc",
      "ocorrido_em": "2026-11-12T13:50:12Z",
      "veiculo": { "vin": "9BSR6X200N3912345", "placa": "RTA2E34", "id_externo": "unid-4471" },
      "dtc": {
        "protocolo": "j1939",
        "spn": 110,
        "fmi": 0,
        "codigo": null,
        "fonte_ecu": 0,
        "descricao": "Temperatura do líquido de arrefecimento acima da faixa normal",
        "ocorrencias": 3,
        "estado": "ativo",
        "lampada": "vermelha",
        "severidade": "critica",
        "acao_recomendada": "parar_imediatamente"
      }
    },
    {
      "id_externo": "trk-88213-20261112T1400-odo",
      "tipo": "leitura",
      "ocorrido_em": "2026-11-12T14:00:00Z",
      "veiculo": { "placa": "RTA2E34" },
      "leituras": [
        { "tipo_medidor": "km", "valor": "412990.4" },
        { "tipo_medidor": "horas", "valor": "18350.2" }
      ]
    },
    {
      "id_externo": "trk-88213-20261112T1405-fuel",
      "tipo": "combustivel",
      "ocorrido_em": "2026-11-12T14:05:00Z",
      "veiculo": { "placa": "RTA2E34" },
      "combustivel": { "evento": "abastecimento_detectado", "nivel_percentual": 94.0, "litros_estimados": "318.0", "consumo_total_litros": "152340.5" }
    },
    {
      "id_externo": "trk-88213-20261112T1407-tpms",
      "tipo": "pneu",
      "ocorrido_em": "2026-11-12T14:07:00Z",
      "veiculo": { "placa": "RTA2E34" },
      "pneu": { "posicao": "E2-ESQ-EXT", "numero_fogo": null, "pressao_psi": "78", "temperatura_c": "61", "alerta": "pressao_baixa" }
    },
    {
      "id_externo": "trk-88213-20261112T1410-pos",
      "tipo": "posicao",
      "ocorrido_em": "2026-11-12T14:10:00Z",
      "veiculo": { "placa": "RTA2E34" },
      "posicao": { "latitude": -23.5329, "longitude": -46.7920, "velocidade_kmh": 0, "ignicao": false }
    }
  ]
}
```

| Tipo | Campos principais | O que a Tenaz faz |
|---|---|---|
| `leitura` | `tipo_medidor`, `valor` | Grava `leitura_medidor` com origem `telematica`, pelas regras da seção 9; no máximo uma leitura gravada por medidor a cada 60 minutos, mas o valor atual sempre se atualiza (proposta) |
| `dtc` | `protocolo` (`j1939`, `obd2`, `proprietario`); J1939: `spn` (0 a 524.287) e `fmi` (0 a 31); OBD-II: `codigo` (ex. `P0217`); `estado` (`ativo`, `inativo`); `lampada`; `severidade` (`critica`, `alta`, `media`, `baixa`, `informativa`); `acao_recomendada` (`parar_imediatamente`, `na_proxima_parada`, `seguro_continuar`) | Aplica as regras de telemática; sem severidade do provedor, usa tabela própria por SPN/FMI, com padrão `media` (proposta). `severidade` é a classificação do DTC vinda do provedor, não a prioridade da OS: a regra decide a prioridade (`baixa`, `media`, `alta`, `urgente`) |
| `combustivel` | `evento` (`nivel`, `abastecimento_detectado`, `queda_brusca`), nível, litros, consumo total | Alimenta consumo e conciliação com cartão; queda brusca pode virar alerta |
| `pneu` | `posicao` (código de eixo e lado), pressão, temperatura, `alerta` | Atualiza o pneu montado na posição (Fase 2) e aplica regras |
| `posicao` | Latitude, longitude, velocidade, ignição | Guarda só a última posição do ativo; a Tenaz não é rastreador e não armazena trilha (proposta) |

**Casamento com o ativo**, nesta ordem:

1. Vínculo manual já feito para `conexao + veiculo.id_externo`.
2. **VIN** (17 caracteres, normalizado em maiúsculas).
3. **Placa** normalizada (maiúsculas, sem hífen nem espaço; aceita o formato antigo `AAA9999` e o Mercosul `AAA9A99`).

VIN e placa que apontam para ativos diferentes geram `nao_casado` com motivo `divergencia_vin_placa`. A Tenaz **nunca cria ativo automaticamente**. Eventos não casados ficam 7 dias (proposta) em `GET /v1/conexoes-telematica/{id}/veiculos-externos` e são reprocessados quando o vínculo é feito.

**Resposta `202 Accepted`**

```json
{
  "recebidos": 5,
  "resultados": [
    { "id_externo": "trk-88213-20261112T1350", "status": "aceito", "ativo_id": "0192f7c4-9999-7aaa-8bbb-000000000112" },
    { "id_externo": "trk-88213-20261112T1400-odo", "status": "aceito", "ativo_id": "0192f7c4-9999-7aaa-8bbb-000000000112" },
    { "id_externo": "trk-88213-20261112T1405-fuel", "status": "duplicado" },
    { "id_externo": "trk-88213-20261112T1407-tpms", "status": "aceito", "ativo_id": "0192f7c4-9999-7aaa-8bbb-000000000112" },
    { "id_externo": "trk-88213-20261112T1410-pos", "status": "aceito", "ativo_id": "0192f7c4-9999-7aaa-8bbb-000000000112" }
  ]
}
```

Status por evento: `aceito`, `duplicado` (mesmo `id_externo` na mesma fonte), `nao_casado`, `rejeitado` (com `erros`). O processamento das regras é assíncrono.

**Regras "DTC → solicitação ou OS"** (`/v1/regras-telematica`):

```json
{
  "nome": "DTC crítico de arrefecimento abre OS emergencial",
  "ativo": true,
  "condicao": {
    "tipo_evento": "dtc",
    "severidade_minima": "alta",
    "spn": [110, 111],
    "classe_ativo_ids": ["0192f7b0-0000-7aaa-8bbb-000000000031"]
  },
  "acao": {
    "tipo": "criar_os",
    "tipo_os": "corretiva_emergencial",
    "prioridade": "urgente",
    "bloquear_despacho": true,
    "notificar_perfis": ["gestor"]
  },
  "janela_deduplicacao_horas": 24
}
```

`acao.tipo`: `criar_solicitacao`, `criar_os`, `apenas_notificar`. A deduplicação usa `ativo + spn + fmi` (ou `codigo`) na janela: um DTC repetido acrescenta ocorrência e comentário na solicitação ou OS aberta, sem criar outra.

**Mapeamento dos provedores (alto nível, a verificar na implementação):**

| Provedor | Origem dos dados | Leituras | DTC | Observação |
|---|---|---|---|---|
| rFMS (Volvo, Scania) | API rFMS da montadora, sem hardware adicional | Odômetro total, horas de motor, combustível total | Conforme a versão do rFMS oferecida | Coleta periódica |
| Geotab | API do MyGeotab | Odômetro, horímetro | Falhas com severidade e ação recomendada (enriquecimento Noregon) | Coleta periódica |
| Cobli | API pública (avisos e registros de manutenção) | A confirmar | Avisos | Endpoints a mapear |
| Ticket Log, Sascar, Omnilink | A verificar | — | — | APIs não verificadas no estudo |

### 14.2 Cartão combustível

Fase 2, plano Profissional. Como as APIs das administradoras de cartão não foram verificadas, a primeira versão é **importação de arquivo** (`POST /v1/importacoes` com `tipo = abastecimentos_cartao` e `layout` do fornecedor: `ticket_log`, `generico` e outros a mapear). Conectores por API entram depois como `conexoes` do mesmo módulo (proposta).

Campos canônicos de `abastecimento` importado:

| Campo | Exemplo | Observação |
|---|---|---|
| `id_transacao_externa` | `"TL-20261112-0048812"` | Deduplicação |
| `data_abastecimento` | `"2026-11-12T09:41:00Z"` | |
| `placa` | `"RTA2E34"` | Casamento pela placa normalizada |
| `produto` | `"diesel_s10"` | `diesel_s10`, `diesel_s500`, `gasolina`, `etanol`, `gnv`, `arla_32` |
| `litros` | `"312.40"` | |
| `valor_unitario`, `valor_total` | `"6.0490"`, `"1889.71"` | Dinheiro como string |
| `hodometro` / `horimetro` | `"412990"` | Gera leitura com origem `cartao_combustivel` |
| `posto_cnpj`, `posto_nome`, `posto_municipio`, `posto_uf` | | |
| `condutor_matricula` | | Identificação do condutor conforme o layout; dados pessoais seguem o documento 9 |

Anomalias marcadas no abastecimento (proposta) e notificadas por `abastecimento.anomalia_detectada`: litros acima da capacidade do tanque; consumo (km/l ou l/h) fora da faixa da classe de ativo (±30%); dois abastecimentos do mesmo veículo em menos de 2 horas; hodômetro regressivo; produto incompatível (Arla em veículo sem SCR, gasolina em veículo a diesel).

### 14.3 Importação por CSV

MVP, todos os planos (web); por token de API, Profissional.

| `tipo` | Chave natural para atualizar | Fase |
|---|---|---|
| `localizacoes` | Caminho (`Planta 1/Linha 2`) | MVP |
| `classes_ativo` | Nome | MVP |
| `ativos` | `tag` | MVP |
| `dados_veiculo` | `tag` ou `placa` | MVP |
| `medidores` | `tag` + `tipo` | MVP |
| `leituras_medidor` | — (sempre insere) | MVP |
| `itens_estoque` | `codigo` | MVP |
| `saldos_iniciais` | `codigo` + almoxarifado | MVP |
| `fornecedores` | `cnpj` (com a coluna `papeis`) | MVP |
| `servicos_fornecedor` | `cnpj` do fornecedor + descrição | MVP |
| `planos_manutencao` | Nome + `tag` ou classe | MVP |
| `documentos_vencimento` | Tipo + dono + número | MVP |
| `pessoas` | `cpf` ou matrícula | MVP |
| `usuarios` | `email` | MVP |
| `abastecimentos_cartao` | `id_transacao_externa` | F2 |
| `pneus` | `numero_fogo` | F2 |
| `pecas_serializadas` | Código do item + `numero_serie` | F2 |
| `locatarios` | `cnpj` ou `cpf` | F2 |
| `projetos` | Nome do projeto + código do modelo; um projeto por linha, com TAG do ativo ou caminho da localização (RF-PRJ-002) (proposta) | F2 |

O tipo `ativos` aceita a coluna `fase_ciclo_vida` (padrão `em_operacao`) só na inserção; na atualização a coluna é ignorada, porque a fase só muda por `colocar-em-operacao` (23.3).

Formato: UTF-8 com ou sem BOM; separador `;` (padrão do Excel em português) ou `,`; separador decimal declarado (`,` ou `.`); datas `DD/MM/AAAA`, `DD/MM/AAAA HH:MM` (no fuso da unidade) ou ISO 8601; primeira linha com os nomes de coluna iguais aos campos da API; referências por chave natural (TAG, código, CNPJ), nunca por UUID. Limite: 50.000 linhas e 20 MB por arquivo (proposta). O modelo de cada tipo está em `GET /v1/importacoes/modelos/{tipo}`.

Fluxo em duas etapas:

1. `POST /v1/importacoes` → `202`. Valida tudo sem gravar.

```json
{
  "tipo": "ativos",
  "anexo_id": "0193a7a0-0000-7aaa-8bbb-000000000e01",
  "opcoes": { "modo": "inserir_ou_atualizar", "separador": ";", "separador_decimal": ",", "formato_data": "DD/MM/AAAA", "tudo_ou_nada": false }
}
```

2. Quando `status = validada`, `GET /v1/importacoes/{id}` mostra o resumo:

```json
{
  "id": "0193a7a1-0000-7aaa-8bbb-000000000e02",
  "tipo": "ativos",
  "status": "validada",
  "resumo": { "linhas_total": 1240, "validas": 1228, "com_erro": 12, "a_inserir": 1100, "a_atualizar": 128 },
  "erros_amostra": [
    { "linha": 17, "coluna": "classe_ativo", "codigo": "nao_encontrado", "mensagem": "Classe \"Empilhadeira eletrica\" não existe. Cadastre a classe ou corrija o nome." },
    { "linha": 233, "coluna": "pai_tag", "codigo": "hierarquia_profundidade_excedida", "mensagem": "O ativo ficaria no 7º nível da árvore." }
  ]
}
```

3. `POST /v1/importacoes/{id}/confirmar` grava as linhas válidas (ou nada, se `tudo_ou_nada` e houver erro), em lotes transacionais de 500 linhas. Estados: `validando`, `validada`, `processando`, `concluida`, `concluida_com_erros`, `falhou`, `cancelada`. A validação expira em 24 horas.

Cada registro gravado gera evento de auditoria com a referência da importação.

### 14.4 NF-e (XML)

Fase 2, plano Profissional. `POST /v1/notas-fiscais-entrada` com `{ "anexo_id": "...", "almoxarifado_id": "..." }` (XML enviado antes por `/v1/anexos`) ou `multipart/form-data` com o arquivo.

Validações:

- XML de NF-e modelo 55 autorizada (`nfeProc` com protocolo de autorização), conforme o esquema vigente; assinatura digital do XML verificada (proposta). Consulta da situação na SEFAZ não entra na primeira versão (a verificar).
- Chave de acesso com 44 dígitos e dígito verificador válido. Chave já importada: `409 nfe_ja_importada`.
- CNPJ destinatário igual ao de uma unidade da empresa: senão, `422 nfe_destinatario_divergente`.
- Grupos do XML que a Tenaz não conhece (por exemplo, os campos de IBS e CBS da fase de teste de 2026 da LC 214/2025) são ignorados sem erro; o XML original fica guardado.

Conciliação (`POST /v1/notas-fiscais-entrada/{id}/conciliar`):

- **Fornecedor** pelo CNPJ emitente; se não existir, é criado como rascunho com o papel `fornecedor_pecas` (proposta).
- **Pedido de compra** pelo `xPed`/`nItemPed` do item, quando informado.
- **Item de estoque** pelo de-para `fornecedor + cProd` → `item_estoque_id`, memorizado na primeira conciliação.
- **Destino de cada item:** entrada no almoxarifado ou aplicação direta em OS ou ativo (peça aplicada na hora, serviço de terceiro). Frete, seguro, outras despesas e descontos são rateados pelo valor dos itens.

`POST /v1/notas-fiscais-entrada/{id}/confirmar-entrada` gera as movimentações de estoque e os custos nas OS, e emite `nota_fiscal_entrada.confirmada`.

**NF-e de remessa e retorno.** Remessas e retornos a subcontratada e envios e retornos de conserto **não** passam por esta importação: a NF-e é emitida pelo ERP do cliente, e a Tenaz só registra `nfe_numero` e `nfe_chave` na movimentação ou na ordem de reparo (20.2 e 20.3). O tratamento fiscal (CFOP) depende da pendência P8.

NFS-e: a primeira versão aceita o XML ou o PDF como anexo, com lançamento manual do valor na OS ou no orçamento do fornecedor. Importação estruturada depende da adoção do padrão nacional de NFS-e pelos municípios dos clientes (a verificar).

---

## 15. Integrações de saída: conectores de ERP

Fase 2, plano **Enterprise** (sumário, seção 2.1). Esta seção é de alto nível; o desenho detalhado fica no documento 5.

**Princípios:**

- Os conectores rodam em workers da Tenaz e consomem os **mesmos eventos internos** que alimentam os webhooks. O cliente não precisa hospedar nada.
- Cada conector usa **somente a API oficial do ERP**. Nunca grava direto no banco do ERP.
- Todo envio é idempotente: a Tenaz guarda o id externo devolvido pelo ERP e não reenvia o que já foi aceito.
- De-para configurável (`/v1/conectores-erp/{id}/mapeamentos`): centro de custo, conta contábil, filial, item, fornecedor, unidade de medida.
- Cada execução é registrada com entrada, saída e erro, e pode ser reprocessada (`.../execucoes/{execucao_id}/reprocessar`). Falhas seguem a mesma política de novas tentativas dos webhooks (24 horas) e depois aparecem no painel de execuções.
- Credenciais do ERP ficam cifradas (KMS) e nunca voltam pela API.

**Fluxos padrão:**

| Fluxo | Direção | Conteúdo | Dono do dado |
|---|---|---|---|
| Cadastros mestres | ERP → Tenaz | Fornecedores, itens de estoque, centros de custo, filiais; opcionalmente bens do imobilizado como ativos | ERP |
| Requisição e pedido de compra | Tenaz → ERP | Requisição aprovada ou pedido de compra vira documento de compra no ERP; status volta para a Tenaz | Tenaz cria; ERP executa |
| Entrada de nota | ERP → Tenaz | Nota lançada no ERP dá entrada no estoque da Tenaz (alternativa à importação de XML) | ERP |
| Custos de manutenção | Tenaz → ERP | Custos de OS encerradas por centro de custo, ativo e conta, por OS ou consolidados no mês | Tenaz |
| Saldos de estoque | Um sentido, configurável | Saldos do almoxarifado | Tenaz **ou** ERP, nunca os dois |
| Remessa, retorno e conserto | Tenaz → ERP, depois ERP → Tenaz | A Tenaz pede ao ERP a emissão da NF-e de remessa ou de retorno; o ERP devolve número e chave, gravados em `nfe_numero` e `nfe_chave` | ERP emite; Tenaz registra |
| Cobrança de locação | Tenaz → ERP | Avarias com custo do locatário e horas ou km excedentes, por contrato | Tenaz apura; ERP fatura |

**Por ERP (a verificar na implementação):**

| ERP | Canal | Observação do estudo de mercado |
|---|---|---|
| Sankhya | API Gateway com AppKey | O ERP exige que toda escrita passe pelo Gateway e não reconhece registros gravados direto no banco |
| TOTVS Protheus | API REST do Protheus; Mensagem Única onde necessário | O módulo SIGAMNT integra por Mensagem Única e webservice TBC, com documentação antiga |
| Omie | API REST pública do Omie | Também é canal de venda (marketplace) |

Configuração de exemplo (`POST /v1/conectores-erp`):

```json
{
  "erp": "sankhya",
  "nome": "Sankhya — matriz",
  "credenciais": { "url_gateway": "https://api.sankhya.com.br", "app_key": "<valor>", "token": "<valor>" },
  "fluxos": {
    "cadastros_mestres": { "ativo": true, "frequencia_minutos": 60 },
    "pedido_compra": { "ativo": true },
    "custos_manutencao": { "ativo": true, "modo": "por_os" },
    "saldos_estoque": { "ativo": false }
  }
}
```

---

## 16. Bloqueio de despacho de veículo

Regra de negócio do estudo de mercado: **veículo com documento crítico vencido gera bloqueio de despacho, exposto por API para o TMS.** A Tenaz informa; a decisão de despachar continua do TMS e do operador.

"Apto para despacho" é uma **verificação calculada** na consulta, não uma entidade (sumário 3.6). As rotas mantêm o nome `/v1/bloqueios-despacho/consulta` por compatibilidade com o texto já distribuído a integradores (ver 26.1, A-3). O bloqueio manual de um veículo é um `bloqueio_ativo` (`/v1/bloqueios-ativo`, 6.13), que também deixa o ativo `parado`.

### 16.1 Motivos de bloqueio

| Código | Origem | Fase |
|---|---|---|
| `documento_vencido` | Documento com vencimento de tipo marcado `critico_despacho` (CRLV-e, cronotacógrafo, CIV, CIPP, seguro obrigatório, conforme configuração) vencido | MVP |
| `ativo_parado` | Veículo com `status_operacional` `parado` ou `em_manutencao` | MVP |
| `checklist_item_critico_reprovado` | Checklist com item `critico_despacho` reprovado, sem OS encerrada depois | MVP |
| `bloqueio_ativo` | Bloqueio de conformidade vigente (`bloqueio_ativo`), de qualquer `origem`: `manual`, `conformidade` ou `regra` (6.13) | MVP |
| `implemento_bloqueado` | Implemento acoplado (carreta) com bloqueio próprio | MVP |
| `dtc_critico_ativo` | DTC ativo com regra que manda bloquear | F2 |
| `pneu_abaixo_sulco_minimo` | Pneu montado abaixo do sulco mínimo configurado | F2 |

Alertas que **não** bloqueiam (exibidos junto): `documento_vencendo`, `preventiva_vencida`, `preventiva_vencendo`.

### 16.2 Consulta

`GET /v1/bloqueios-despacho/consulta?placa=RTA2E34` (ou `?vin=...`, ou `?ativo_id=...`), escopo `despacho:consulta`.

**Resposta `200 OK`** (`Cache-Control: private, max-age=60`)

```json
{
  "veiculo": {
    "ativo_id": "0192f7c4-9999-7aaa-8bbb-000000000112",
    "tag": "FRT-0112",
    "placa": "RTA2E34",
    "vin": "9BSR6X200N3912345",
    "status_operacional": "operando"
  },
  "liberado": false,
  "consultado_em": "2026-11-12T14:15:00Z",
  "valido_ate": "2026-11-12T14:16:00Z",
  "motivos": [
    {
      "codigo": "documento_vencido",
      "descricao": "Cronotacógrafo: verificação Inmetro vencida em 08/11/2026.",
      "ativo_id": "0192f7c4-9999-7aaa-8bbb-000000000112",
      "origem": { "entidade": "documento_vencimento", "id": "0193a0b0-0000-7aaa-8bbb-000000000c11" },
      "desde": "2026-11-09T03:00:00Z"
    }
  ],
  "alertas": [
    {
      "codigo": "documento_vencendo",
      "descricao": "CRLV-e 2026 vence em 24/11/2026.",
      "ativo_id": "0192f7c4-9999-7aaa-8bbb-000000000112",
      "origem": { "entidade": "documento_vencimento", "id": "0193a0b0-0000-7aaa-8bbb-000000000c12" }
    }
  ],
  "composicao": [
    { "ativo_id": "0192f7c4-9999-7aaa-8bbb-000000000231", "tag": "CAR-0031", "placa": "QWE4R56", "tipo": "implemento", "liberado": true }
  ]
}
```

- `liberado` é `false` se houver qualquer motivo no veículo ou nos implementos acoplados. `liberado` é a verificação "apto para despacho" do sumário 3.3, separada do status operacional: um veículo `operando` pode estar com `liberado: false` por documento crítico vencido.
- Placa ou VIN não encontrados: `404 recurso_nao_encontrado`. Veículo inativo: `liberado: false` com motivo `ativo_parado`.
- Consulta em lote: `POST /v1/bloqueios-despacho/consultas` com `{ "veiculos": [ { "placa": "RTA2E34" }, { "vin": "9BWZZZ377VT004251" } ] }` (até 100), resposta com um item por veículo na mesma forma, e `erro` no item quando não encontrado.
- Para evitar consulta a cada despacho, o TMS pode manter um cache alimentado pelos webhooks `ativo.despacho_bloqueado` e `ativo.despacho_liberado` e usar a consulta como confirmação. A política do TMS quando a Tenaz estiver indisponível (liberar com alerta ou bloquear) é decisão do cliente; a documentação recomenda registrar a decisão.
- Toda consulta fica registrada para auditoria (quem consultou, quando, resultado), com retenção de 90 dias (proposta).

---

## 17. Exportação completa dos dados da empresa

Disponível em **todos os planos**, desde o MVP: portabilidade reduz o medo de dependência do fornecedor. Meta do estudo de mercado: exportação completa em **até 24 horas**.

### 17.1 Pedido

`POST /v1/exportacoes-empresa` — perfil `administrador` com MFA validado nos últimos 15 minutos (proposta; senão `401 mfa_obrigatorio`), ou token com escopo `exportacao:empresa`.

```json
{
  "formato": "json",
  "incluir_anexos": true,
  "incluir_auditoria": true,
  "entidades": null,
  "notificar_emails": ["ti@cliente.com.br"]
}
```

- `formato`: `json` (NDJSON, um arquivo por entidade) ou `csv` (UTF-8, separador `;`).
- `entidades: null` exporta tudo; uma lista restringe.
- Resposta `202 Accepted`, `Location: /v1/exportacoes-empresa/{id}`. Uma exportação em andamento por empresa e até 4 por dia (seção 5).

### 17.2 Situação e download

`GET /v1/exportacoes-empresa/0193a8a0-0000-7aaa-8bbb-000000000f01`

```json
{
  "id": "0193a8a0-0000-7aaa-8bbb-000000000f01",
  "status": "concluida",
  "formato": "json",
  "solicitada_por_id": "0192f800-aaaa-7bbb-8ccc-000000000001",
  "solicitada_em": "2026-11-12T14:00:00Z",
  "concluida_em": "2026-11-12T15:12:40Z",
  "dados_ate": "2026-11-12T14:00:00Z",
  "tamanho_bytes": 7340032000,
  "arquivos": [
    { "nome": "tenaz-exportacao-2026-11-12T1400Z.part1.zip", "tamanho_bytes": 5368709120, "sha256": "a3f1...9c", "url_download": "https://exportacoes.tenaz.com.br/...", "url_expira_em": "2026-11-13T15:12:40Z" },
    { "nome": "tenaz-exportacao-2026-11-12T1400Z.part2.zip", "tamanho_bytes": 1971322880, "sha256": "7be0...41", "url_download": "https://exportacoes.tenaz.com.br/...", "url_expira_em": "2026-11-13T15:12:40Z" }
  ],
  "arquivos_apagados_em": "2026-11-19T15:12:40Z"
}
```

Estados: `na_fila`, `em_andamento`, `concluida`, `falhou`, `cancelada`, `expirada`. URLs pré-assinadas valem 24 horas e podem ser renovadas com novo `GET` enquanto os arquivos existirem (7 dias, proposta).

### 17.3 Conteúdo do pacote

```
tenaz-exportacao-2026-11-12T1400Z/
  manifest.json              versão do esquema, versão da API, data de corte, contagem e sha256 por arquivo
  esquema/                   JSON Schema de cada entidade, gerado do OpenAPI
  dados/
    ativo.ndjson
    ordem_servico.ndjson
    leitura_medidor.ndjson
    ...                      um arquivo por entidade, campos iguais aos da API
  anexos/
    <anexo_id>/<nome_arquivo>
  auditoria/
    evento_auditoria.ndjson
```

- O pacote é um retrato consistente no instante `dados_ate` (leitura num snapshot transacional ou réplica, proposta).
- Pacotes acima de 5 GB são divididos em partes (proposta).
- Credenciais, segredos de webhook, hashes de senha, tokens e hashes de links de uso único **não** são exportados. Não há dado de saúde a exportar: a Tenaz não o guarda (sumário 3.6).
- Cada exportação gera evento de auditoria e o webhook `exportacao.concluida`.
- Contas com SRV ou com conexões: o pacote traz `cliente_servico`, `especialidade_tecnico`, `pessoa_especialidade`, `servico_os`, os equipamentos `de_cliente` (inclusive as cópias sincronizadas, com o que já foi recebido), os `vinculo_os` da conta e a visão da conta sobre as suas `conexao_empresa`. Nunca traz dado da outra conta além do que já foi gravado nesta pela sincronização.
- **Agendamento** (Fase 2): `PUT /v1/exportacoes-empresa/agendamento` com `{ "frequencia": "mensal", "dia": 1, "formato": "json", "incluir_anexos": false }`.
- Após o fim do contrato, a empresa mantém acesso de leitura e exportação por 90 dias (proposta; política no documento 9).

---

## 18. Versionamento, depreciação e compatibilidade

### 18.1 O que é compatível

Entram em `/v1` sem nova versão, registradas no changelog:

- endpoint, recurso, ação, filtro ou expansão novos;
- campo novo opcional na requisição;
- campo novo na resposta;
- valor novo em enum marcado como extensível (`x-enum-extensivel`);
- código de erro novo;
- tipo de evento de webhook novo;
- cabeçalho novo opcional;
- validação mais flexível.

São **incompatíveis** e exigem nova versão principal: remover ou renomear campo, endpoint ou valor de enum; mudar tipo ou formato de campo; tornar obrigatório um campo opcional; mudar o significado de campo, código de erro ou estado; mudar valor padrão; endurecer validação de forma que requisições válidas passem a falhar.

Exceção: correções de segurança e adequações legais podem exigir mudança incompatível com aviso curto. Nesse caso, a Tenaz avisa os administradores das empresas afetadas com a maior antecedência possível e registra o motivo no changelog.

### 18.2 Política de depreciação (proposta)

| Item | Regra |
|---|---|
| Nova versão principal (`/v2`) | Só quando inevitável. `/v1` continua suportada por **pelo menos 18 meses** após o lançamento estável da `/v2`; contratos Enterprise podem estender |
| Campo ou endpoint depreciado dentro de `/v1` | Continua funcionando até o fim da `/v1`; aparece como `deprecated: true` no OpenAPI |
| Aviso nas respostas | Cabeçalhos `Deprecation` (RFC 9745), `Sunset` (RFC 8594) e `Link: <https://docs.tenaz.com.br/changelog#...>; rel="deprecation"` em toda resposta de endpoint depreciado |
| Aviso ativo | E-mail aos administradores das empresas cujos tokens usaram o item depreciado nos últimos 30 dias, no anúncio, 6 meses, 3 meses e 1 mês antes do `Sunset` |
| Visibilidade | `GET /v1/tokens-api/{id}` mostra `usos_depreciados` (endpoint, última chamada) |
| Após o `Sunset` | `410 Gone` com código `versao_descontinuada`, por mais 90 dias, antes de remover a rota |

Exemplo:

```
Deprecation: @1798761600
Sunset: Wed, 30 Jun 2027 23:59:59 GMT
Link: <https://docs.tenaz.com.br/changelog#2027-01-campo-custo>; rel="deprecation"
```

### 18.3 Regras para clientes

- Ignorar campos desconhecidos na resposta.
- Tratar valor desconhecido em enum extensível como "outro", sem erro. Os tipos TypeScript gerados representam enums extensíveis como união de literais mais `string` (proposta).
- Decidir por `codigo` de erro, e tratar código desconhecido pelo `status` HTTP.
- Não depender da ordem dos campos nem da ordem de itens sem `ordenar`.
- Seguir o `Location` de recursos criados e de processamentos assíncronos.

### 18.4 Recursos em beta

Endpoints com `x-estabilidade: beta` exigem o cabeçalho `Tenaz-Beta: <nome-do-recurso>`. Podem mudar com 30 dias de aviso e ficam fora do SLA. Sem o cabeçalho, devolvem `404`.

### 18.5 Versão do app de campo

- O app envia `X-Tenaz-App-Versao`. `GET /v1/configuracao-app` informa `versao_minima` e `versao_recomendada`.
- Abaixo da mínima, a API devolve `400 versao_app_nao_suportada` em todos os endpoints, **exceto** `POST /v1/sincronizacao/comandos`, que continua aceitando comandos por 30 dias depois que a versão deixa de ser suportada (proposta). Assim nenhum registro feito offline se perde por atualização pendente.
- O app nunca é forçado a atualizar no meio de uma OS: a tela de atualização aparece depois que a fila de comandos esvazia.

### 18.6 Versão dos eventos de webhook

Cada tipo de evento tem `versao_evento` inteira. Mudança incompatível no `dados` cria nova versão do evento; a assinatura escolhe a versão (`versao_evento`), e a anterior continua sendo emitida por 12 meses (proposta).

### 18.7 Changelog

Changelog público na documentação, com data, tipo (adição, alteração, depreciação, correção) e módulo, gerado a partir do diff do OpenAPI e revisado à mão. Toda mudança de contrato entra também no `CHANGELOG.md` do repositório, no mesmo pull request (RNF-040; o CI rejeita mudança sem registro). Mudanças relevantes também vão para os administradores por e-mail. Cada entrada do changelog traz a versão do produto (18.8) em que a mudança entrou.

### 18.8 Versão do produto (SemVer)

A Tenaz tem **duas numerações**, com papéis diferentes:

| Numeração | Formato | Onde aparece | Quando muda |
|---|---|---|---|
| Versão do produto | SemVer 2.0.0: `MAJOR.MINOR.PATCH` (ex.: `1.4.2`) | `GET /v1/versao`, cabeçalho `Tenaz-Versao`, rodapé do web, tela "Sobre" do app, changelog | A cada publicação em produção |
| Versão da API | Prefixo da rota (`/v1`) | URL | Só em quebra de compatibilidade da API pública (18.1) |

Regras do SemVer na Tenaz:

- **MAJOR** sobe quando há quebra de compatibilidade da API pública, nos termos da seção 18.1 (remover ou renomear campo, endpoint ou enum, mudar tipo ou significado, endurecer validação).
- **MINOR** sobe com funcionalidade nova compatível: endpoint, campo, filtro, evento de webhook, valor de enum extensível.
- **PATCH** sobe com correção que não muda o contrato.
- Versões de pré-lançamento usam o sufixo do SemVer (`1.5.0-rc.1`) e só existem no sandbox.

**Relação com o prefixo `/v1`:** os dois números são independentes. O prefixo só muda em quebra de compatibilidade; por isso, toda quebra que cria um prefixo novo (`/v2`) vem numa versão MAJOR nova do produto, mas o número do prefixo não acompanha o MAJOR (o produto `3.0.0` pode ser o que lança a `/v2`). Enquanto a `/v1` estiver suportada (18.2), versões novas do produto continuam a servi-la sem quebra: mudanças MINOR e PATCH nunca alteram o contrato de uma rota já publicada.

**`GET /v1/versao`** — público, sem autenticação, em todos os planos, desde o MVP (RNF-039). É exceção aceita à regra de rota no plural (sumário 3.6). Limite de 60 req/min por IP (seção 5). Resposta `200 OK`, com `Cache-Control: public, max-age=60`:

```json
{
  "versao": "1.4.2",
  "commit": "a1b2c3d",
  "data_build": "2026-11-12T14:03:00Z",
  "ambiente": "producao"
}
```

| Campo | Descrição |
|---|---|
| `versao` | Versão do produto em SemVer 2.0.0 |
| `commit` | Hash curto (7 caracteres) do commit publicado |
| `data_build` | Data e hora UTC do build |
| `ambiente` | `producao` ou `sandbox` |

A resposta **não expõe** dependências, versões de bibliotecas, sistema operacional, região, nomes de servidor nem qualquer outro dado de infraestrutura. O mesmo vale para o cabeçalho. No contrato gerado, o esquema da resposta se chama `Versao` (até a versão 0.4 deste documento, o código o gerava como `Dados`). **Divergência aberta (A-38):** o servidor devolve só `versao`, `commit` e `data_build`, como a seção 9.5 do documento 5; o campo `ambiente` desta seção não sai.

**Cabeçalho `Tenaz-Versao`.** Toda resposta da API, de sucesso ou de erro, inclusive `GET /v1/versao`, traz `Tenaz-Versao: 1.4.2` com a versão do produto que atendeu a requisição. O suporte usa esse valor junto com o `id_rastreio` para localizar o problema. O web e o app mostram a versão na interface (rodapé e tela "Sobre"), lida de `GET /v1/versao`; o app mostra também a própria versão (`X-Tenaz-App-Versao`, 18.5), que segue SemVer à parte.

---

## 19. Terceiros: prestadores, OS atribuída e link de uso único

Módulo TER (sumário 2 e 3.1). No MVP e em todos os planos: fornecedor com papéis, catálogo de serviços, OS executada por prestador, perfil `prestador` e link de uso único. Na Fase 2, plano Profissional: contratos de prestação com SLA, avaliação de desempenho e oficina externa com orçamento, aprovação, execução, nota e garantia.

### 19.1 Fornecedor e papéis

`fornecedor` é o cadastro único de empresas externas. Uma empresa pode ter vários papéis ao mesmo tempo.

`POST /v1/fornecedores`

```json
{
  "razao_social": "Hidrotec Manutenção Industrial Ltda.",
  "nome_fantasia": "Hidrotec",
  "cnpj": "12345678000190",
  "papeis": ["subcontratada", "assistencia_tecnica"],
  "contatos": [
    { "nome": "Rafael Moura", "email": "rafael@hidrotec.com.br", "telefone": "+5511987654321", "principal": true }
  ],
  "endereco": { "logradouro": "Rua das Indústrias, 450", "municipio": "Guarulhos", "uf": "SP", "cep": "07220000" }
}
```

| `papel` | Uso na API |
|---|---|
| `fornecedor_pecas` | Pedido de compra, NF-e de entrada |
| `subcontratada` | Executa OS com equipe própria; pode ter almoxarifado `em_terceiro` |
| `assistencia_tecnica` | Executa OS; recebe peças em ordem de reparo |
| `oficina` | Executa OS fora da empresa (antiga "oficina externa"), com orçamento |
| `locador` | Lado do contrato de locação `tomado` |
| `posto_combustivel` | Posto do abastecimento (`abastecimentos.fornecedor_id`); a importação do cartão casa o posto pelo `posto_cnpj` (14.2) |
| `seguradora` | Seguradora do sinistro (`sinistros.seguradora_id`) e da apólice registrada como documento com vencimento |

Um fornecedor com papel `subcontratada`, `assistencia_tecnica` ou `oficina` pode ter conta própria na Tenaz (SRV) e ser conectado à conta do cliente (`POST /v1/conexoes-empresa`, 22.5). `GET /v1/fornecedores/{id}` traz então `conexao_empresa: { "id": "...", "estado": "ativa" }`.

O papel é conferido em cada operação: OS atribuída exige `subcontratada`, `assistencia_tecnica` ou `oficina` (ou `locador` em ativo tomado); almoxarifado `em_terceiro`, `subcontratada`; ordem de reparo, `assistencia_tecnica`, `oficina` ou `subcontratada`; contrato `tomado`, `locador`. Papel ausente devolve `422 fornecedor_papel_incompativel`.

### 19.2 Catálogo de serviços

`POST /v1/servicos-fornecedor`

```json
{
  "fornecedor_id": "0193b100-0000-7aaa-8bbb-000000000101",
  "descricao": "Manutenção corretiva de compressor de parafuso",
  "classe_ativo_ids": ["0192f7b0-0000-7aaa-8bbb-000000000021"],
  "marcas": ["Atlas Copco", "Schulz"],
  "area_atendimento": { "municipios_ibge": ["3518800", "3550308"], "raio_km": null },
  "prazo_atendimento_horas": 4,
  "prazo_solucao_horas": 24,
  "preco_referencia": "220.00",
  "unidade_preco": "hora"
}
```

`unidade_preco`: `hora`, `visita` ou `servico`. A busca "quem atende este ativo" é `GET /v1/servicos-fornecedor?classe_ativo_id=...&marca=Atlas%20Copco&municipio_ibge=3518800`, ordenada por prazo de atendimento.

### 19.3 Atribuir OS a prestador

`POST /v1/ordens-servico/{id}/atribuir-fornecedor` — perfil Gestão, `If-Match`, `Idempotency-Key`. Aceita OS em `aberta`, `planejada`, `programada` ou `pausada` (proposta).

```
POST /v1/ordens-servico/0193a2f0-5b6c-7d8e-9f01-23456789abcd/atribuir-fornecedor HTTP/1.1
If-Match: "2"
Idempotency-Key: 9a7e1c3b-2d4f-4e6a-8b0c-1d2e3f4a5b6c
Content-Type: application/json
```

```json
{
  "fornecedor_id": "0193b100-0000-7aaa-8bbb-000000000101",
  "responsavel_id": "0193b1a0-0000-7aaa-8bbb-000000000201",
  "data_programada": "2026-11-13T12:00:00Z",
  "data_limite": "2026-11-14T12:00:00Z",
  "responsavel_custo": "empresa",
  "observacao_fornecedor": "Levar kit de vedação do eixo. Acesso pela portaria 2; apresentar crachá e o certificado de NR-12 na entrada.",
  "notificar": true
}
```

| Campo | Obrigatório | Regra |
|---|---|---|
| `fornecedor_id` | Sim | Fornecedor ativo com papel compatível (19.1) |
| `responsavel_id` | Não | Usuário `prestador` do mesmo fornecedor. Sem ele, a OS fica visível a todos os usuários `prestador` do fornecedor, e o primeiro que assumir vira responsável |
| `data_programada` | Não | Com `responsavel_id`, leva a OS a `programada` (mesmas regras de `programar`) |
| `data_limite` | Não | Padrão: agora + `prazo_solucao_horas` do serviço ou do contrato de prestação, quando houver |
| `responsavel_custo` | Não | Padrão conforme 8.6; diferente do padrão exige `justificativa_responsavel_custo` |
| `observacao_fornecedor` | Não | Texto visível ao prestador |
| `notificar` | Não | Padrão `true`: e-mail e push aos usuários `prestador` do fornecedor |

**Resposta `200 OK`** (`ETag: "3"`), trecho:

```json
{
  "id": "0193a2f0-5b6c-7d8e-9f01-23456789abcd",
  "numero": "OS-2026-004812",
  "versao": 3,
  "status": "programada",
  "prioridade": "alta",
  "executor_tipo": "fornecedor",
  "fornecedor_id": "0193b100-0000-7aaa-8bbb-000000000101",
  "responsavel_id": "0193b1a0-0000-7aaa-8bbb-000000000201",
  "responsavel_custo": "empresa",
  "contrato_locacao_id": null,
  "data_programada": "2026-11-13T12:00:00Z",
  "data_limite": "2026-11-14T12:00:00Z",
  "data_reconhecimento": null,
  "links_execucao_ativos": 0,
  "transicoes_permitidas": ["desprogramar", "reconhecer", "iniciar", "cancelar"]
}
```

Efeitos: a OS entra no alcance do fornecedor (4.9) e no escopo offline dos seus usuários `prestador`; o histórico da OS registra a atribuição; webhook `ordem_servico.atualizada`. O usuário `prestador` que executa OS conta como técnico cobrado (sumário 2.1).

**Fornecedor com conexão ativa (SRV).** Se o fornecedor tem `conexao_empresa` `ativa` e o ativo da OS está compartilhado na conexão, a atribuição cria o `vinculo_os` e a OS correspondente na conta da subcontratada (22.8): `responsavel_id` e `data_programada` não são aceitos (`422 validacao_falhou`; quem distribui ao técnico é a subcontratada), links de uso único não são gerados para essa OS (`409 os_vinculada_somente_leitura`) e a resposta traz `vinculo.estado_sincronizacao = "pendente"`. Com conexão `suspensa`, a atribuição devolve `409 conexao_inativa`, salvo com `sem_sincronizacao: true`, que trata o fornecedor como se não tivesse conta própria. Ativo fora de `ativos_compartilhados` da conexão devolve `422 validacao_falhou` (`erros[].codigo = "ativo_nao_compartilhado"`), salvo com `compartilhar_ativo: true`, que inclui o ativo na conexão na mesma transação (exige perfil Admin).

`POST /v1/ordens-servico/{id}/retirar-fornecedor` com `{ "motivo": "..." }` volta a OS para `executor_tipo = interno`, limpa `fornecedor_id` e o responsável externo, revoga links abertos e tira a OS do alcance do fornecedor. Apontamentos já feitos pelo prestador ficam na OS. Em OS vinculada, a retirada cancela a cópia na conta da subcontratada (motivo `retirada_pelo_cliente`) e fecha o vínculo; o que já foi sincronizado fica nas duas contas.

### 19.4 Link seguro de uso único

Para o prestador ocasional, sem usuário na Tenaz. O link permite registrar **uma** execução de **uma** OS atribuída a um fornecedor, sem login, e não conta como técnico cobrado.

**Gerar o link** — `POST /v1/links-execucao` (perfil Gestão):

```json
{
  "ordem_servico_id": "0193a2f0-5b6c-7d8e-9f01-23456789abcd",
  "nome_executor": "José Aparecido Lima",
  "documento_executor": "12345678909",
  "validade_horas": 72,
  "exigir_assinatura": true,
  "exigir_fotos": true,
  "enviar_para": { "email": "jose.lima@hidrotec.com.br", "whatsapp": null }
}
```

| Campo | Obrigatório | Regra |
|---|---|---|
| `ordem_servico_id` | Sim | OS com `executor_tipo = fornecedor` (senão `422 os_sem_fornecedor`), em estado não terminal |
| `nome_executor` | Sim | Nome que vai para a OS e para o registro NR-12, quando houver |
| `documento_executor` | Não | CPF do executor; guardado como dado pessoal (documento 9) |
| `validade_horas` | Não | Padrão do parâmetro da empresa (proposta: 72 h); máximo 168 h |
| `exigir_assinatura`, `exigir_fotos` | Não | Padrão `true` |
| `enviar_para` | Não | Sem destino, a URL só volta na resposta, para a empresa enviar por conta própria |

**Resposta `201 Created`**:

```json
{
  "id": "0193b2c0-0000-7aaa-8bbb-000000000301",
  "ordem_servico_id": "0193a2f0-5b6c-7d8e-9f01-23456789abcd",
  "fornecedor_id": "0193b100-0000-7aaa-8bbb-000000000101",
  "nome_executor": "José Aparecido Lima",
  "situacao": "ativo",
  "url": "https://app.tenaz.com.br/l/#lx_5Qm2Zp8Kd3Vt7Ww1Nc9Hb4Ry6Ls0Ga2Ef",
  "expira_em": "2026-11-15T14:00:00Z",
  "utilizado_em": null,
  "criado_em": "2026-11-12T14:00:00Z"
}
```

**Onde fica o token (sumário 3.6).** O token vai no **fragmento** da URL (`/l/#<token>`), nunca no caminho nem na query string. O navegador não envia o fragmento ao servidor: ele não aparece em log de acesso, em proxy, em CDN nem no cabeçalho `Referer`. A página `/l/` é estática; ela lê o fragmento, troca o token por uma sessão curta e apaga o fragmento da barra de endereço (`history.replaceState`). Nenhuma rota da API recebe o token no caminho.

- O token (`lx_` + 32 caracteres base62, 190 bits de entropia, proposta) aparece **só nesta resposta**. A Tenaz guarda o hash SHA-256; `GET /v1/links-execucao/{id}` nunca devolve `url`.
- `situacao`: `ativo`, `utilizado`, `expirado`, `revogado`. Um link por vez por OS (proposta): gerar outro revoga o anterior.
- O link é revogado automaticamente quando a OS é retirada do fornecedor, cancelada ou encerrada.

**Ativar o link** — `POST /v1/publico/links-execucao/ativar`. Sem cookie e sem CSRF; o token é a única credencial e vai no corpo:

```
POST /v1/publico/links-execucao/ativar HTTP/1.1
Host: app.tenaz.com.br
Content-Type: application/json
```

```json
{
  "token": "lx_5Qm2Zp8Kd3Vt7Ww1Nc9Hb4Ry6Ls0Ga2Ef"
}
```

**Resposta `200 OK`** (`Cache-Control: no-store`):

```json
{
  "sessao_link": "lxs_9Tq4Wm1Zc7Kp2Vn8Hd3Rb6Ys0Lf5Ge1Ja",
  "expira_em": "2026-11-13T12:50:00Z",
  "link": { "situacao": "ativo", "expira_em": "2026-11-15T14:00:00Z" }
}
```

| Regra | Detalhe |
|---|---|
| Sessão curta | `sessao_link` (`lxs_` + 32 caracteres base62, proposta) vale **30 minutos**, renovados a cada uso, até no máximo 4 horas e nunca além do `expira_em` do link. A Tenaz guarda só o hash |
| Uso da sessão | As demais chamadas enviam `Authorization: Bearer lxs_...`. A página guarda a sessão em memória e em `sessionStorage` (proposta), nunca em `localStorage` nem em cookie; por usar `Authorization`, não há CSRF |
| Escopo | `link`: a sessão só alcança a OS do link e só as três rotas `/v1/publico/links-execucao/atual/*` |
| Reativação | Ativar de novo, com o link ainda `ativo`, emite sessão nova e revoga a anterior (um aparelho por vez, proposta). Ativar não consome o link |
| Reconhecimento | A primeira ativação grava `data_reconhecimento` da OS, se vazia |
| Erros | Token inexistente ou malformado: `404 recurso_nao_encontrado`, sem distinguir de link de outra empresa. Link vencido: `410 link_execucao_expirado`; usado: `410 link_execucao_utilizado`; revogado: `410 link_execucao_revogado`. Sessão vencida ou revogada nas demais rotas: `401 credencial_invalida` (a página pede para abrir o link de novo) |
| Limites | 10 ativações por minuto por IP; 10 tokens inválidos seguidos bloqueiam o IP por 15 minutos (seção 5) |

**Consultar a OS do link** — `GET /v1/publico/links-execucao/atual`, com `Authorization: Bearer lxs_9Tq4Wm1Zc7Kp2Vn8Hd3Rb6Ys0Lf5Ge1Ja`. Resposta `200` com o **escopo mínimo** necessário para executar:

```json
{
  "empresa": { "nome_fantasia": "Metalúrgica Exemplo" },
  "ordem_servico": {
    "numero": "OS-2026-004812",
    "tipo": "corretiva_planejada",
    "status": "programada",
    "prioridade": "alta",
    "titulo": "Substituir retentor do eixo da unidade compressora",
    "descricao": "Vazamento de óleo identificado na SS-2026-000731.",
    "data_programada": "2026-11-13T12:00:00Z",
    "data_limite": "2026-11-14T12:00:00Z",
    "observacao_fornecedor": "Levar kit de vedação do eixo. Acesso pela portaria 2; apresentar crachá e o certificado de NR-12 na entrada.",
    "exige_registro_nr12": true,
    "tarefas": [
      { "id": "0193a2f1-0000-7aaa-8bbb-000000000001", "descricao": "Aplicar bloqueio e etiquetagem da fonte elétrica", "obrigatoria": true },
      { "id": "0193a2f1-0000-7aaa-8bbb-000000000002", "descricao": "Substituir retentor do eixo principal", "obrigatoria": true },
      { "id": "0193a2f1-0000-7aaa-8bbb-000000000003", "descricao": "Testar proteções e dispositivos de parada de emergência", "obrigatoria": true }
    ],
    "medidores": [ { "id": "0192f7d1-2222-7ccc-8ddd-000000000031", "tipo": "horas", "unidade": "h" } ],
    "campos_extras": [
      { "chave": "rompeu_lacre", "rotulo": "Rompeu lacre?", "tipo": "sim_nao", "obrigatorio_em": "conclusao", "ajuda": "Marque sim se algum lacre do equipamento foi retirado.", "valor": null }
    ]
  },
  "ativo": { "tag": "CMP-0347", "nome": "Compressor de parafuso 75 cv", "localizacao": "Planta 1 / Casa de compressores" },
  "executor": { "nome": "José Aparecido Lima" },
  "exigir_assinatura": true,
  "exigir_fotos": true,
  "expira_em": "2026-11-15T14:00:00Z"
}
```

O que **não** sai pelo link: ids internos da OS e do ativo, custos, histórico do ativo, outras OS, dados de usuários e técnicos, apontamentos de terceiros, documentos técnicos, unidade e CNPJ da empresa. Cabeçalhos da resposta: `Cache-Control: no-store`, `Referrer-Policy: no-referrer`, `X-Robots-Tag: noindex`.

**Enviar fotos e assinatura** — `POST /v1/publico/links-execucao/atual/anexos`, com a sessão do link e o mesmo corpo de `POST /v1/anexos` (2.13):

```json
{
  "id": "0193b2c1-0000-7aaa-8bbb-0000000000f3",
  "nome_arquivo": "assinatura-jose.png",
  "tipo_conteudo": "image/png",
  "tamanho_bytes": 18422,
  "sha256": "5f2a9c0e1b7d4a3f8e6c2b1a0d9f8e7c6b5a4d3c2b1a0f9e8d7c6b5a4f3e2d1c"
}
```

Resposta `201` com `id` e `url_upload` (válida por 15 minutos). A URL de upload só aceita o arquivo declarado, e o anexo fica preso ao link: só pode ser referenciado no registro de execução desse link. Limite de 20 anexos por link (proposta).

**Registrar a execução** — `POST /v1/publico/links-execucao/atual/registrar-execucao`, com a sessão do link e `Idempotency-Key` obrigatório:

```
POST /v1/publico/links-execucao/atual/registrar-execucao HTTP/1.1
Host: app.tenaz.com.br
Authorization: Bearer lxs_9Tq4Wm1Zc7Kp2Vn8Hd3Rb6Ys0Lf5Ge1Ja
Idempotency-Key: 4c1e7a2b-9d3f-4e5a-8b6c-0d1e2f3a4b5c
Content-Type: application/json
```

```json
{
  "data_inicio": "2026-11-13T12:20:00Z",
  "data_conclusao": "2026-11-13T15:05:00Z",
  "servico_realizado": "Trocado o retentor do eixo principal e limpa a carcaça. Teste de 30 minutos em carga sem vazamento.",
  "tarefas_concluidas_ids": [
    "0193a2f1-0000-7aaa-8bbb-000000000001",
    "0193a2f1-0000-7aaa-8bbb-000000000002",
    "0193a2f1-0000-7aaa-8bbb-000000000003"
  ],
  "leituras": [ { "medidor_id": "0192f7d1-2222-7ccc-8ddd-000000000031", "valor": "18343.4" } ],
  "pecas_utilizadas": [ { "descricao": "Retentor 45x62x8 NBR", "quantidade": "1", "fornecida_por": "fornecedor" } ],
  "registro_nr12": {
    "condicoes_seguranca": "Proteções reinstaladas; intertravamento da porta lateral e parada de emergência testados.",
    "indicacao_conclusiva": "apto",
    "restricoes": null
  },
  "anexos_ids": ["0193b2c1-0000-7aaa-8bbb-0000000000f1", "0193b2c1-0000-7aaa-8bbb-0000000000f2"],
  "assinatura_anexo_id": "0193b2c1-0000-7aaa-8bbb-0000000000f3",
  "campos_extras": { "rompeu_lacre": false },
  "causa_externa": null,
  "aceite_termo": true
}
```

**Resposta `200 OK`**:

```json
{
  "ordem_servico": { "numero": "OS-2026-004812", "status": "concluida" },
  "recebido_em": "2026-11-13T15:06:12Z",
  "comprovante": "Execução registrada em 13/11/2026 às 12:06 (horário de Brasília). A empresa vai conferir e liberar a OS.",
  "link": { "situacao": "utilizado" }
}
```

Regras do registro:

- Uma única chamada bem-sucedida consome o link (`situacao = utilizado`, `utilizado_em` gravado) e revoga a sessão. Repetir com a mesma `Idempotency-Key` devolve a mesma resposta (a idempotência é guardada por link, não por sessão); com outra chave, `410 link_execucao_utilizado`. Link vencido ou revogado depois da ativação: `410 link_execucao_expirado` ou `410 link_execucao_revogado`.
- A execução aplica, numa transação, as ações `iniciar` (se a OS não estava em execução) e `concluir`, com `data_inicio` e `data_conclusao` informadas, mais as leituras (origem `os`), as tarefas e os anexos. A OS vai para `concluida` e aguarda a liberação da empresa, que pode `recusar-liberacao` e gerar novo link.
- `data_reconhecimento` é gravada na primeira ativação do link, se ainda vazia.
- O executor entra na OS como texto (`nome_executor`), sem usuário; o autor nos eventos e na auditoria é `{ "tipo": "link_execucao", "id": "<id do link>" }`, com escopo de acesso `link`, IP e agente do navegador.
- Peças: só descrição e quantidade. O link nunca movimenta estoque; peça do estoque da empresa é baixada depois por um usuário interno. `fornecida_por`: `empresa` ou `fornecedor`.
- Assinatura eletrônica avançada (sumário 2.1): o registro guarda hash SHA-256 do conteúdo, carimbo de tempo, IP e aceite do termo. A validade jurídica para registro NR-12 depende da pendência P1.
- Validações e erros de campo seguem as mesmas regras do `concluir` (8.2); erro de validação **não** consome o link.
- **Campos extras pelo link** (RF-OS-021): a consulta mostra só as definições com `preenchido_pelo_prestador = true` que se aplicam à OS, pela `chave` (sem o id da definição), e o registro aceita só essas chaves; outra chave devolve `422 validacao_falhou`. Campo obrigatório na conclusão vazio devolve `422 campo_extra_obrigatorio` sem consumir o link. Campo obrigatório sem `preenchido_pelo_prestador` não impede o registro pelo link: a OS fica `concluida` com o aviso `campos_extras_pendentes`, e a empresa preenche antes de `encerrar` (8.2). O link também aceita `causa_externa` (RF-CTR-003), com as evidências em `anexos_ids`.

### 19.5 Orçamentos, contratos de prestação e avaliação (Fase 2)

- **Orçamento de fornecedor** (`/v1/orcamentos-fornecedor`) substitui o antigo orçamento de oficina e o serviço externo. Liga-se a uma OS (`ordem_servico_id`, OS com `executor_tipo = fornecedor`; senão `422 os_sem_fornecedor`) ou a uma ordem de reparo (`ordem_reparo_id`). Estados: `pendente`, `aprovado`, `recusado`. Aprovar grava o valor como custo previsto de terceiros; a nota do fornecedor, no encerramento, vira `custo_terceiros`. Fluxo da oficina: orçamento → aprovação → execução → nota → garantia (`garantia_dias` no orçamento aprovado).
- **Contrato de prestação** (`/v1/contratos-prestacao`): vigência, escopo, ativos ou classes cobertos, `sla_atendimento_horas`, `sla_solucao_horas`, preços. O servidor identifica o contrato vigente que cobre o ativo quando a OS é atribuída ao fornecedor e mede o SLA de atendimento (`data_reconhecimento` − atribuição) e de solução (`data_conclusao` − atribuição). Estouro emite `fornecedor.sla_violado`.
- **Avaliação** (`/v1/avaliacoes-fornecedor`): SLA cumprido e retrabalho (OS do mesmo ativo e modo de falha reaberta ou com liberação recusada dentro da garantia) são calculados; a nota e o comentário são informados pela gestão.

### 19.6 Isolamento do prestador

O alcance do perfil `prestador` está na seção 4.9. Resumo das garantias da API:

- toda consulta com credencial de `prestador` recebe, no repositório, a condição `fornecedor_id = usuario.fornecedor_id`, além do RLS por empresa;
- um mesmo fornecedor que atende duas empresas clientes tem usuários distintos em cada uma; uma credencial nunca vale para duas empresas (4.7);
- registro fora do alcance devolve `404`, nunca `403`;
- testes automatizados de vazamento entre fornecedores e entre empresas rodam em todo build (risco R4);
- a subcontratada com conta própria (SRV) não usa o perfil `prestador`: o isolamento entre a conta dela e a do cliente é o isolamento entre empresas (4.7), mais a lista fechada de campos da conexão (22.6, risco R5).

---

## 20. Estoque estendido: estoque em terceiros, peças serializadas e conserto

Fase 2, plano Profissional (sumário 2.1 e 4), salvo o estoque básico do MVP.

### 20.1 Almoxarifados e tipos de movimentação

`almoxarifado.tipo`: `proprio`, `em_terceiro` (estoque da empresa em poder de uma subcontratada, com `fornecedor_id` obrigatório) e `quarentena` (peças com defeito aguardando decisão). O **estoque geral** é a soma de todos os almoxarifados (`GET /v1/saldos-estoque/consolidado`).

| `tipo` | Origem → destino | Gerada por | Plano | NF-e |
|---|---|---|---|---|
| `entrada` | — → `proprio` | `POST /v1/movimentacoes-estoque` ou NF-e de entrada | E | Opcional |
| `saida_os` | qualquer → OS | `POST /v1/ordens-servico/{id}/pecas` | E | — |
| `devolucao_os` | OS → qualquer | `.../pecas/{peca_id}/estornar` | E | — |
| `ajuste` | qualquer, com motivo | `POST /v1/movimentacoes-estoque` | E | — |
| `transferencia` | `proprio` → `proprio` ou `quarentena` | `POST /v1/movimentacoes-estoque` | P | — |
| `remessa_terceiro` | `proprio` → `em_terceiro` | `POST /v1/movimentacoes-estoque` | P | `nfe_numero`, `nfe_chave` |
| `retorno_terceiro` | `em_terceiro` → `proprio` ou `quarentena` | `POST /v1/movimentacoes-estoque` | P | `nfe_numero`, `nfe_chave` |
| `envio_conserto` | `quarentena` → fornecedor | `POST /v1/ordens-reparo/{id}/enviar` | P | `nfe_numero`, `nfe_chave` |
| `retorno_conserto` | fornecedor → `proprio` | `POST /v1/ordens-reparo/{id}/registrar-retorno` | P | `nfe_numero`, `nfe_chave` |
| `baixa_sucata` | `quarentena` → — | `POST /v1/movimentacoes-estoque` ou `.../pecas-serializadas/{id}/baixar` | P | — |

Origem ou destino incompatível com o tipo devolve `422 almoxarifado_tipo_incompativel`. Em remessas e retornos, `nfe_numero` e `nfe_chave` são opcionais no registro (a nota pode chegar depois, por `PATCH` da movimentação, só nesses dois campos) e, quando informados, a chave é validada (`422 nfe_chave_invalida`). A Tenaz registra a NF-e emitida pelo ERP; não a emite. CFOPs e regras fiscais: pendência P8.

### 20.2 Remessa de peças a subcontratada

`POST /v1/movimentacoes-estoque` — perfil Estoque, `Idempotency-Key`.

```json
{
  "id": "0193b300-0000-7aaa-8bbb-000000000401",
  "tipo": "remessa_terceiro",
  "almoxarifado_origem_id": "0192fa00-0000-7aaa-8bbb-000000000001",
  "almoxarifado_destino_id": "0193b110-0000-7aaa-8bbb-000000000111",
  "data_movimentacao": "2026-11-12T16:00:00Z",
  "itens": [
    { "item_estoque_id": "0192fa10-3333-7ddd-8eee-000000000101", "quantidade": "20" },
    { "item_estoque_id": "0192fa10-3333-7ddd-8eee-000000000144", "quantidade": "2", "pecas_serializadas_ids": ["0193b200-0000-7aaa-8bbb-000000000501", "0193b200-0000-7aaa-8bbb-000000000502"] }
  ],
  "nfe_numero": "18342",
  "nfe_chave": "35261112345678000190550010000183421000183428",
  "observacao": "Estoque de consignação para o contrato de manutenção da linha 2."
}
```

| Campo | Regra |
|---|---|
| `almoxarifado_destino_id` | Almoxarifado `em_terceiro`; o `fornecedor_id` dele precisa ter o papel `subcontratada` |
| `itens[].quantidade` | Até o saldo disponível na origem (`422 estoque_insuficiente`) |
| `itens[].pecas_serializadas_ids` | Obrigatório para item com `controla_serie = true`, uma peça por unidade (`422 numero_serie_obrigatorio`); cada peça precisa estar `em_estoque` na origem (`409 peca_situacao_invalida`) |

**Resposta `201 Created`**:

```json
{
  "id": "0193b300-0000-7aaa-8bbb-000000000401",
  "versao": 1,
  "tipo": "remessa_terceiro",
  "status": "em_transito",
  "almoxarifado_origem_id": "0192fa00-0000-7aaa-8bbb-000000000001",
  "almoxarifado_destino_id": "0193b110-0000-7aaa-8bbb-000000000111",
  "fornecedor_id": "0193b100-0000-7aaa-8bbb-000000000101",
  "data_movimentacao": "2026-11-12T16:00:00Z",
  "itens": [
    { "item_estoque_id": "0192fa10-3333-7ddd-8eee-000000000101", "quantidade": "20", "custo_unitario": "18.4500", "custo_total": "369.00" },
    { "item_estoque_id": "0192fa10-3333-7ddd-8eee-000000000144", "quantidade": "2", "custo_unitario": "1240.0000", "custo_total": "2480.00", "pecas_serializadas_ids": ["0193b200-0000-7aaa-8bbb-000000000501", "0193b200-0000-7aaa-8bbb-000000000502"] }
  ],
  "nfe_numero": "18342",
  "nfe_chave": "35261112345678000190550010000183421000183428",
  "recebido_em": null,
  "recebido_por_id": null,
  "criado_em": "2026-11-12T16:00:04Z",
  "criado_por_id": "0192f800-aaaa-7bbb-8ccc-000000000031"
}
```

- O saldo sai da origem na hora e entra no destino como `em_transito` (as peças serializadas ficam com `situacao = em_transito`). **Saldo e peça `em_transito` não podem ser consumidos, reservados nem movimentados** até a confirmação de recebimento (sumário 3.6): `422 estoque_em_transito`, com `quantidade_em_transito`. O saldo disponível em `GET /v1/saldos-estoque` não inclui o que está em trânsito, que sai à parte em `quantidade_em_transito`. Quando a subcontratada confirma o recebimento (`POST /v1/movimentacoes-estoque/{id}/confirmar-recebimento`, perfil `prestador` do fornecedor ou Estoque), o saldo fica disponível no almoxarifado `em_terceiro` e a movimentação passa a `recebida`. Sem confirmação em 5 dias (parâmetro, proposta), a gestão é notificada.
- Divergência no recebimento é registrada na confirmação (`itens[].quantidade_recebida`) e gera `ajuste` com motivo automático.
- A subcontratada consome do próprio almoxarifado `em_terceiro` nas OS atribuídas a ela (`saida_os`). O **retorno** de sobras ou de peças retiradas é `tipo = retorno_terceiro`, do `em_terceiro` para `proprio` (ou para `quarentena`, se defeituosas), com NF-e de retorno.
- Webhook: `movimentacao_estoque.registrada`.

### 20.3 Envio de peça para conserto: ordem de reparo

Número `OR-AAAA-NNNNNN`, por empresa, reinicia a cada ano.

**Destino da peça com defeito.** Decidir o que acontece com a peça na `quarentena` (abrir ordem de reparo, declarar irreparável com destino `sucata` ou `devolvida_sem_reparo`, baixar como sucata) exige a permissão `est.decidir_destino` (sumário 3.6; padrão em `gestor`, `planejador` e `almoxarife`). Sem ela: `403 permissao_negada`. As demais ações da ordem de reparo (enviar, registrar orçamento e retorno) seguem os perfis do catálogo (6.8).

| Ação | Rota | De | Para | Campos obrigatórios | Efeitos |
|---|---|---|---|---|---|
| Abrir | `POST /v1/ordens-reparo` | — | `aberta` | `fornecedor_id`, `itens` (peças da `quarentena`) | Reserva as peças na quarentena |
| `enviar` | `POST .../enviar` | `aberta` | `enviada` | `data_envio`; `nfe_numero` e `nfe_chave` de remessa (podem vir depois, proposta) | Gera `envio_conserto`; peças → `em_conserto`; webhook `ordem_reparo.enviada` |
| `registrar-orcamento` | `POST .../registrar-orcamento` | `enviada` | `em_orcamento` | `valor`, `prazo_dias` | Cria `orcamento_fornecedor` ligado à ordem |
| `aprovar` | `POST .../aprovar` | `em_orcamento` | `aprovada` | — | Aprova o orçamento; `prazo_previsto` = aprovação + `prazo_dias` |
| `registrar-retorno` | `POST .../registrar-retorno` | `aprovada` | `retornada` | `data_retorno`, `almoxarifado_destino_id`, `custo_reparo`, `garantia_dias`; NF-e de retorno | Gera `retorno_conserto`; peça volta `recondicionada` com custo e garantia; webhook `ordem_reparo.retornada` |
| `declarar-irreparavel` | `POST .../declarar-irreparavel` | `enviada`, `em_orcamento`, `aprovada` | `irreparavel` | `motivo`, `destino` (`sucata` ou `devolvida_sem_reparo`) | `sucata`: `baixa_sucata`; `devolvida_sem_reparo`: `retorno_conserto` para a `quarentena` sem custo de reparo |
| `cancelar` | `POST .../cancelar` | `aberta` | `cancelada` | `motivo` | Libera as peças na quarentena |

Orçamento recusado (`POST /v1/orcamentos-fornecedor/{id}/recusar`) leva a `declarar-irreparavel` ou ao cancelamento da ordem com a peça de volta e abertura de nova ordem com outro fornecedor (sumário 3.3). Item com `reparavel = false` devolve `422 item_nao_reparavel`.

**Abrir** — `POST /v1/ordens-reparo`:

```json
{
  "fornecedor_id": "0193b100-0000-7aaa-8bbb-000000000102",
  "itens": [
    {
      "item_estoque_id": "0192fa10-3333-7ddd-8eee-000000000144",
      "peca_serializada_id": "0193b200-0000-7aaa-8bbb-000000000509",
      "defeito_relatado": "Bomba hidráulica com vazamento interno; pressão cai para 90 bar em carga.",
      "ordem_servico_origem_id": "0193a2f0-5b6c-7d8e-9f01-23456789abcd"
    }
  ],
  "prazo_previsto": "2026-11-26",
  "observacao": "Enviar com o laudo de remoção."
}
```

**Enviar** — `POST /v1/ordens-reparo/0193b400-.../enviar`, `If-Match: "1"`:

```json
{
  "data_envio": "2026-11-13T09:30:00Z",
  "nfe_numero": "18351",
  "nfe_chave": "35261112345678000190550010000183511000183515",
  "transportadora": "Coleta pelo fornecedor"
}
```

**Registrar o retorno** — `POST /v1/ordens-reparo/0193b400-0000-7aaa-8bbb-000000000601/registrar-retorno`, `If-Match: "4"`, `Idempotency-Key`:

```json
{
  "data_retorno": "2026-11-24T14:10:00Z",
  "almoxarifado_destino_id": "0192fa00-0000-7aaa-8bbb-000000000001",
  "nfe_numero": "7712",
  "nfe_chave": "35261198765432000110550010000077121000077129",
  "itens": [
    {
      "peca_serializada_id": "0193b200-0000-7aaa-8bbb-000000000509",
      "custo_reparo": "1380.00",
      "garantia_dias": 90,
      "servico_executado": "Troca do kit de vedação e das placas de desgaste; bancada a 210 bar sem vazamento.",
      "anexos_ids": ["0193b401-0000-7aaa-8bbb-0000000000a1"]
    }
  ]
}
```

**Resposta `200 OK`** (`ETag: "5"`):

```json
{
  "id": "0193b400-0000-7aaa-8bbb-000000000601",
  "numero": "OR-2026-000045",
  "versao": 5,
  "status": "retornada",
  "fornecedor_id": "0193b100-0000-7aaa-8bbb-000000000102",
  "orcamento_fornecedor_id": "0193b410-0000-7aaa-8bbb-000000000611",
  "data_envio": "2026-11-13T09:30:00Z",
  "prazo_previsto": "2026-11-26",
  "data_retorno": "2026-11-24T14:10:00Z",
  "nfe_remessa": { "numero": "18351", "chave": "35261112345678000190550010000183511000183515" },
  "nfe_retorno": { "numero": "7712", "chave": "35261198765432000110550010000077121000077129" },
  "custo_total": "1380.00",
  "itens": [
    {
      "peca_serializada_id": "0193b200-0000-7aaa-8bbb-000000000509",
      "numero_serie": "BH-A10V-772193",
      "situacao": "em_estoque",
      "recondicionada": true,
      "almoxarifado_id": "0192fa00-0000-7aaa-8bbb-000000000001",
      "custo_reparo": "1380.00",
      "garantia_dias": 90,
      "garantia_ate": "2027-02-22"
    }
  ],
  "movimentacoes_ids": ["0193b402-0000-7aaa-8bbb-000000000621", "0193b402-0000-7aaa-8bbb-000000000622"],
  "transicoes_permitidas": []
}
```

O custo do reparo entra no custo médio do item (proposta; regra de custeio no documento 4). Se a peça recondicionada falhar dentro da garantia, a OS que a remove pode abrir nova ordem de reparo com `responsavel_custo = garantia`.

### 20.4 Peças serializadas

Item com `controla_serie = true` é controlado peça a peça (`peca_serializada`): `numero_serie` único por item (`409 numero_serie_duplicado`), `situacao` (`em_estoque`, `instalada`, `em_transito`, `em_conserto`, `baixada`), `almoxarifado_id` ou `ativo_id`, `recondicionada`, garantia. Toda movimentação, consumo em OS, instalação e remoção desse item exige as peças (`422 numero_serie_obrigatorio`).

- `POST /v1/pecas-serializadas/{id}/instalar` com `{ "ativo_id": "...", "ordem_servico_id": "...", "posicao": "Bomba principal" }`: `em_estoque` → `instalada`, junto com a `saida_os`.
- `POST /v1/pecas-serializadas/{id}/remover` com `{ "ordem_servico_id": "...", "almoxarifado_destino_id": "<quarentena>", "defeituosa": true, "defeito_relatado": "..." }`: `instalada` → `em_estoque` na quarentena. Peça defeituosa reparável fica disponível para ordem de reparo.
- `GET /v1/pecas-serializadas/{id}/historico` lista instalações (ativo, OS, datas, horas ou km do medidor na instalação e na remoção, quando houver) e reparos (ordens de reparo, custo, garantia), base do tempo médio entre falhas da peça.
- Operação fora da situação exigida devolve `409 peca_situacao_invalida` com `situacao_atual`.

---

## 21. Locação: contratos, vistorias e responsabilidade de custo

Módulo LOC, Fase 2, plano Profissional (sumário 2.1 e 4).

### 21.1 Estados e regras

`contrato_locacao`: `rascunho` → `vigente` → `encerrado`, ou `cancelado`. Número `CL-AAAA-NNNNNN`, por empresa, reinicia a cada ano.

| Ação | De | Para | Exige |
|---|---|---|---|
| `ativar` | `rascunho` | `vigente` | Ao menos um item (`422 contrato_locacao_sem_itens`); matriz de responsabilidade preenchida; nenhum ativo em outro contrato vigente no período (`409 ativo_ja_locado`) |
| `encerrar` | `vigente` | `encerrado` | Todos os itens com vistoria de devolução (`422 vistoria_devolucao_pendente`, com `itens_pendentes`); `motivo_encerramento` quando antes de `data_fim` (`422 motivo_obrigatorio`) |
| `cancelar` | Só `rascunho` | `cancelado` | `motivo_cancelamento`. Contrato `vigente` não se cancela: termina por `encerrar` (`409 transicao_invalida`) |

Regras do sumário 3.6: **sublocação não é permitida** (ativo `de_terceiro` em contrato `cedido` devolve `422 sublocacao_nao_permitida`); ativo `de_terceiro` devolvido ao locador no fim do contrato `tomado` passa a `inativo`.

**Sentido:**

| `sentido` | Contraparte | Efeito nos ativos |
|---|---|---|
| `cedido` | `locatario_id` | Ativo próprio continua `operando` durante a locação; a localização exibida passa a ser o local de uso do contrato. As OS do ativo recebem `contrato_locacao_id` e o `responsavel_custo` padrão da matriz. O perfil `locatario` passa a ver o ativo (4.9) |
| `tomado` | `fornecedor_id` com papel `locador` (`422 fornecedor_papel_incompativel` sem ele) | O ativo é cadastrado (ou já existe) com `propriedade = de_terceiro`. Manutenção a cargo do locador pela matriz vira OS com `executor_tipo = fornecedor` (o locador) e `responsavel_custo = locador`. O SLA de substituição do locador é medido em `.../itens/{item_id}/substituir`. Na vistoria de devolução ao locador, o ativo passa a `inativo` |

### 21.2 Criar contrato de locação

**Cedido** — `POST /v1/contratos-locacao`:

```json
{
  "sentido": "cedido",
  "locatario_id": "0193c100-0000-7aaa-8bbb-000000000701",
  "unidade_id": "0192f790-0000-7aaa-8bbb-000000000001",
  "data_inicio": "2026-12-01",
  "data_fim": "2027-05-31",
  "local_uso": {
    "descricao": "Obra Residencial Parque das Águas",
    "endereco": "Av. Brasil, 2100, Campinas, SP",
    "latitude": -22.9071,
    "longitude": -47.0632
  },
  "franquia_horas": "200",
  "franquia_km": null,
  "periodo_franquia": "mes",
  "valor_excedente_hora": "95.00",
  "valor_excedente_km": null,
  "itens": [
    { "ativo_id": "0192f7c4-5555-7aaa-8bbb-000000000340", "data_entrega_prevista": "2026-12-01", "data_devolucao_prevista": "2027-05-31" }
  ],
  "responsabilidades": [
    { "categoria": "preventiva", "responsavel_custo": "empresa" },
    { "categoria": "corretiva_desgaste", "responsavel_custo": "empresa" },
    { "categoria": "avaria_mau_uso", "responsavel_custo": "locatario" },
    { "categoria": "pneus", "responsavel_custo": "locatario" },
    { "categoria": "deslocamento", "responsavel_custo": "locatario" }
  ],
  "observacao": "Retroescavadeira com operador do locatário."
}
```

**Resposta `201 Created`** (`Location: /v1/contratos-locacao/0193c200-...`, `ETag: "1"`):

```json
{
  "id": "0193c200-0000-7aaa-8bbb-000000000801",
  "numero": "CL-2026-000012",
  "versao": 1,
  "sentido": "cedido",
  "status": "rascunho",
  "locatario_id": "0193c100-0000-7aaa-8bbb-000000000701",
  "fornecedor_id": null,
  "unidade_id": "0192f790-0000-7aaa-8bbb-000000000001",
  "data_inicio": "2026-12-01",
  "data_fim": "2027-05-31",
  "local_uso": { "descricao": "Obra Residencial Parque das Águas", "endereco": "Av. Brasil, 2100, Campinas, SP", "latitude": -22.9071, "longitude": -47.0632 },
  "franquia_horas": "200",
  "franquia_km": null,
  "periodo_franquia": "mes",
  "valor_excedente_hora": "95.00",
  "valor_excedente_km": null,
  "sla_substituicao_horas": null,
  "itens": [
    {
      "id": "0193c201-0000-7aaa-8bbb-000000000811",
      "ativo_id": "0192f7c4-5555-7aaa-8bbb-000000000340",
      "data_entrega_prevista": "2026-12-01",
      "data_devolucao_prevista": "2027-05-31",
      "data_entrega": null,
      "data_devolucao": null,
      "vistoria_entrega_id": null,
      "vistoria_devolucao_id": null
    }
  ],
  "responsabilidades": [
    { "categoria": "preventiva", "responsavel_custo": "empresa" },
    { "categoria": "corretiva_desgaste", "responsavel_custo": "empresa" },
    { "categoria": "avaria_mau_uso", "responsavel_custo": "locatario" },
    { "categoria": "pneus", "responsavel_custo": "locatario" },
    { "categoria": "deslocamento", "responsavel_custo": "locatario" }
  ],
  "transicoes_permitidas": ["ativar", "cancelar"],
  "criado_em": "2026-11-12T17:00:00Z",
  "atualizado_em": "2026-11-12T17:00:00Z",
  "criado_por_id": "0192f800-aaaa-7bbb-8ccc-000000000009"
}
```

**Tomado** — mesmo endpoint, com o locador no lugar do locatário:

```json
{
  "sentido": "tomado",
  "fornecedor_id": "0193b100-0000-7aaa-8bbb-000000000150",
  "unidade_id": "0192f790-0000-7aaa-8bbb-000000000001",
  "data_inicio": "2026-11-20",
  "data_fim": "2027-02-20",
  "numero_contrato_contraparte": "LOC-88412/2026",
  "local_uso": { "descricao": "Centro de distribuição Jundiaí, doca 4" },
  "franquia_horas": "180",
  "franquia_km": null,
  "periodo_franquia": "mes",
  "valor_excedente_hora": "60.00",
  "valor_excedente_km": null,
  "sla_substituicao_horas": 24,
  "itens": [
    {
      "ativo": {
        "tipo": "equipamento",
        "classe_ativo_id": "0192f7b0-0000-7aaa-8bbb-000000000044",
        "tag": "EMP-L-0007",
        "nome": "Empilhadeira elétrica 2,5 t (locada)",
        "marca": "Exemplo",
        "modelo": "E25",
        "numero_serie": "E25-554120",
        "criticidade": "B"
      },
      "data_entrega_prevista": "2026-11-20",
      "data_devolucao_prevista": "2027-02-20"
    }
  ],
  "responsabilidades": [
    { "categoria": "preventiva", "responsavel_custo": "locador" },
    { "categoria": "corretiva_desgaste", "responsavel_custo": "locador" },
    { "categoria": "avaria_mau_uso", "responsavel_custo": "empresa" },
    { "categoria": "pneus", "responsavel_custo": "locador" },
    { "categoria": "deslocamento", "responsavel_custo": "locador" }
  ]
}
```

Resposta `201` na mesma forma, com `sentido: "tomado"`, `fornecedor_id` preenchido, `locatario_id: null` e o ativo criado (`itens[].ativo_id`) com `propriedade: "de_terceiro"`. No contrato tomado, `itens[]` aceita `ativo_id` de ativo já cadastrado ou o objeto `ativo` para cadastrá-lo na mesma transação.

| Campo | Regra |
|---|---|
| `local_uso` | Nome canônico (sumário 3.6). Objeto com `descricao` (obrigatória), `endereco`, `latitude`, `longitude`. No `cedido`, é a localização exibida do ativo durante a locação |
| `franquia_horas`, `franquia_km` | Nomes canônicos (sumário 3.6). Decimal em string, por `periodo_franquia`; `null` quando não há franquia daquele tipo. Ao menos um dos dois quando houver cobrança de excedente |
| `periodo_franquia` | `mes` ou `contrato` (proposta) |
| `valor_excedente_hora`, `valor_excedente_km` | Dinheiro por hora ou km acima da franquia (proposta) |

### 21.3 Matriz de responsabilidade

`responsabilidade_locacao` define, por categoria, quem paga a manutenção. As categorias são as do sumário 3.6: `preventiva`, `corretiva_desgaste`, `avaria_mau_uso`, `pneus`, `deslocamento`. A matriz tem exatamente uma linha por categoria; faltar ou repetir categoria devolve `422 validacao_falhou`. `responsavel_custo` aceita `empresa`, `locatario` (só em `cedido`) e `locador` (só em `tomado`).

`PUT /v1/contratos-locacao/{id}/responsabilidades` substitui a matriz inteira (corpo: `{ "responsabilidades": [ ... ] }`, mesma forma de 21.2). Em contrato `vigente`, a troca vale para OS criadas depois, exige `X-Tenaz-Motivo` e gera auditoria.

Como a OS recebe o padrão: o tipo de OS define a categoria (`preventiva`, `preditiva` e `inspecao_legal` → `preventiva`; `corretiva_planejada` e `corretiva_emergencial` → `corretiva_desgaste`, salvo indicação de mau uso; OS gerada por avaria na vistoria → `avaria_mau_uso`). Pneus e deslocamento são lançados como itens de custo da OS com a categoria própria (proposta). O padrão pode ser trocado com justificativa (8.6).

### 21.4 Vistoria de devolução com avaria gerando OS

`POST /v1/vistorias-locacao` — perfil Gestão ou Execução, `Idempotency-Key`. A vistoria de entrega usa o mesmo corpo com `tipo = entrega` e grava `data_entrega` no item; a de devolução exige a de entrega do mesmo item (`422 vistoria_entrega_pendente`).

```json
{
  "id": "0193c300-0000-7aaa-8bbb-000000000901",
  "contrato_locacao_id": "0193c200-0000-7aaa-8bbb-000000000801",
  "item_contrato_locacao_id": "0193c201-0000-7aaa-8bbb-000000000811",
  "tipo": "devolucao",
  "data_vistoria": "2027-05-31T15:30:00Z",
  "checklist_modelo_id": "0192fd00-0000-7aaa-8bbb-000000000640",
  "checklist_modelo_versao": 2,
  "respostas": [
    { "item_id": "0192fd41-0000-7aaa-8bbb-000000000001", "valor_booleano": true },
    { "item_id": "0192fd41-0000-7aaa-8bbb-000000000002", "valor_booleano": false, "observacao": "Braço da caçamba com trinca na solda do pino.", "anexos_ids": ["0193c2ff-0000-7aaa-8bbb-0000000000b1"] }
  ],
  "leituras": [ { "medidor_id": "0192f7d1-5555-7ccc-8ddd-000000000340", "valor": "4412.5" } ],
  "avarias": [
    {
      "descricao": "Trinca na solda do pino do braço da caçamba",
      "gravidade": "alta",
      "ativo_parado": true,
      "anexos_ids": ["0193c2ff-0000-7aaa-8bbb-0000000000b1", "0193c2ff-0000-7aaa-8bbb-0000000000b2"],
      "gerar_os": true,
      "responsavel_custo": "locatario"
    },
    {
      "descricao": "Risco na pintura da cabine",
      "gravidade": "baixa",
      "ativo_parado": false,
      "anexos_ids": ["0193c2ff-0000-7aaa-8bbb-0000000000b3"],
      "gerar_os": false
    }
  ],
  "assinaturas": [
    { "parte": "empresa", "nome": "Marcos Vieira", "assinatura_anexo_id": "0193c2ff-0000-7aaa-8bbb-0000000000c1" },
    { "parte": "locatario", "nome": "Paula Ribeiro", "documento": "98765432100", "assinatura_anexo_id": "0193c2ff-0000-7aaa-8bbb-0000000000c2" }
  ]
}
```

| Campo | Regra |
|---|---|
| `tipo` | `entrega` ou `devolucao` |
| `respostas` | Mesmas regras de checklist (seção 10); itens reprovados não geram solicitação na vistoria, viram avaria quando o vistoriador marcar |
| `leituras` | Origem `manual`, regras da seção 9. Na devolução, são a base das horas ou km excedentes à franquia |
| `avarias[].gerar_os` | `true` cria uma OS `corretiva_planejada` (ou `corretiva_emergencial` se `gravidade = alta` e `ativo_parado`, proposta) com `origem = vistoria_locacao`, `contrato_locacao_id`, `responsavel_custo = locatario` (padrão da categoria `avaria_mau_uso`) e as fotos da avaria |
| `assinaturas` | Uma por parte (`empresa`, `locatario` no cedido; `empresa`, `locador` no tomado). Parte que não assinou na hora assina depois em `POST /v1/vistorias-locacao/{id}/assinaturas`; o locatário pode assinar pelo portal |

**Resposta `201 Created`**:

```json
{
  "id": "0193c300-0000-7aaa-8bbb-000000000901",
  "versao": 1,
  "contrato_locacao_id": "0193c200-0000-7aaa-8bbb-000000000801",
  "numero_contrato": "CL-2026-000012",
  "item_contrato_locacao_id": "0193c201-0000-7aaa-8bbb-000000000811",
  "ativo_id": "0192f7c4-5555-7aaa-8bbb-000000000340",
  "tipo": "devolucao",
  "data_vistoria": "2027-05-31T15:30:00Z",
  "resultado": "com_avarias",
  "leituras_registradas": [
    { "id": "0193c301-0000-7aaa-8bbb-000000000911", "medidor_id": "0192f7d1-5555-7ccc-8ddd-000000000340", "valor": "4412.5", "status_validacao": "valida" }
  ],
  "uso_no_periodo": {
    "horas": { "inicial": "3180.0", "final": "4412.5", "utilizado": "1232.5", "franquia_horas": "1200", "excedente": "32.5" },
    "km": null
  },
  "avarias": [
    {
      "id": "0193c302-0000-7aaa-8bbb-000000000921",
      "descricao": "Trinca na solda do pino do braço da caçamba",
      "gravidade": "alta",
      "ordem_servico": { "id": "0193c310-0000-7aaa-8bbb-000000000931", "numero": "OS-2027-000418", "status": "aberta", "tipo": "corretiva_emergencial", "prioridade": "urgente", "responsavel_custo": "locatario" }
    },
    {
      "id": "0193c302-0000-7aaa-8bbb-000000000922",
      "descricao": "Risco na pintura da cabine",
      "gravidade": "baixa",
      "ordem_servico": null
    }
  ],
  "assinaturas_pendentes": [],
  "item": { "data_devolucao": "2027-05-31T15:30:00Z" },
  "status_operacional_ativo": "parado",
  "criado_em": "2027-05-31T15:41:10Z"
}
```

Efeitos: grava `data_devolucao` no item; a localização exibida do ativo volta à de cadastro; a avaria com `ativo_parado` coloca o ativo em `parado` (a OS gerada é corretiva emergencial); a OS com `responsavel_custo = locatario` e o excedente da franquia entram no relatório de cobrança (`GET /v1/contratos-locacao/{id}/cobranca`). Webhooks: `vistoria_locacao.registrada`, `leitura_medidor.registrada`, `ordem_servico.criada`. Quando todos os itens tiverem vistoria de devolução, o contrato pode ser encerrado (`POST .../encerrar`, webhook `contrato_locacao.encerrado`).

### 21.5 Relatório de cobrança

`GET /v1/contratos-locacao/{id}/cobranca?inicio=2026-12-01&fim=2027-05-31` devolve, por item: uso no período contra `franquia_horas` e `franquia_km` (excedente e valor), OS com `responsavel_custo = locatario` com o custo consolidado de cada uma (só OS `encerrada`; as abertas aparecem em `pendentes`), e o total. A Tenaz apura; a fatura é emitida pelo ERP (fluxo "Cobrança de locação", seção 15). O locatário vê o mesmo relatório dos próprios contratos.

---

## 22. SRV: gestão de serviços do prestador e conexão entre contas

Módulo SRV (sumário 2, 2.2 e 3.5), no **MVP**, vendido no plano **Serviços**. A subcontratada ou assistência técnica tem **conta própria**: uma `empresa` com o módulo SRV, que cadastra técnicos, clientes atendidos e equipamentos de clientes, abre e distribui OS, registra serviço, componentes e horas, e cujos técnicos usam o app de campo. Quando o cliente também usa a Tenaz, as duas contas se conectam e a OS trafega entre elas. Pacote Go: `servicos`. A decisão de arquitetura da conexão é o **ADR-009 Conexão entre contas** (documento 5); o risco é o **R5** do sumário.

### 22.1 Entidades e rotas

O SRV reaproveita as entidades do restante da Tenaz, com o significado do lado de quem presta o serviço:

| Necessidade | Entidade (nome canônico) | Rota |
|---|---|---|
| Técnico | `pessoa` + `usuario` com perfil `tecnico`; especialidades em `especialidade_tecnico` (catálogo da conta) e ligação `pessoa_especialidade` | `/v1/pessoas`, `/v1/usuarios`, `/v1/especialidades-tecnico` |
| Cliente atendido | `cliente_servico`: `tipo_pessoa` (`juridica`, `fisica`), `documento` (CNPJ ou CPF), `nome`, `contatos`, `enderecos_atendimento`, `conexao_empresa_id` (opcional) | `/v1/clientes-servico` |
| Equipamento de cliente | `ativo` com `propriedade = de_cliente` e `cliente_servico_id`; a cópia sincronizada guarda `conexao_empresa_id` e `ativo_remoto_id` | `/v1/ativos` |
| Serviço a realizar | Campo `descricao_servico` da OS + itens `servico_os` (referência a `servico_fornecedor` da própria conta, quantidade, preço de referência) | `/v1/ordens-servico/{id}/servicos` |
| Componentes | Itens de material da OS com `tipo_material` (`consumido`, `substituido`, `fornecido_pelo_cliente`), `numero_serie_retirado`, `numero_serie_instalado` | `/v1/ordens-servico/{id}/pecas` |
| Horas | `data_inicio` e `data_conclusao` da OS; `apontamento` por técnico | Ações da OS e `/v1/ordens-servico/{id}/apontamentos` |
| Relatório | PDF da OS com serviço, componentes, horas, fotos e assinatura do cliente | `GET /v1/ordens-servico/{id}/relatorio-pdf` |
| Conexão | `conexao_empresa`, no esquema da plataforma, fora do RLS por empresa: `empresa_cliente_id`, `empresa_prestadora_id`, `fornecedor_id` (na conta do cliente), `cliente_servico_id` (na conta da subcontratada), `estado`, `campos_compartilhados`, `data_convite`, `data_aceite`, `data_encerramento` | `/v1/conexoes-empresa` |
| Vínculo de OS | `vinculo_os`: `conexao_empresa_id`, `os_cliente_id`, `os_prestador_id`, `lado_origem` (`cliente`, `prestador`), `estado_sincronizacao` (`pendente`, `sincronizado`, `erro`) | `/v1/vinculos-os` e o objeto `vinculo` da OS |

**Ids de outra conta nunca saem na API.** `fornecedor_id` (conta do cliente) só aparece para o cliente; `cliente_servico_id` (conta da subcontratada), só para a subcontratada; `ativo_remoto_id` não sai para ninguém. No vínculo, cada conta vê o id da **própria** OS em `ordem_servico_id` e só o **número** da OS do outro lado em `numero_os_remota`. Os nomes `os_cliente_id` e `os_prestador_id` são colunas do banco (documento 2), não campos da API. `conexao_empresa.id` é o mesmo nas duas contas, porque o registro é da plataforma.

### 22.2 Plano e alcance da conta da subcontratada

- **Plano Serviços** (sumário 2.2): SRV completo (técnicos, clientes, equipamentos de clientes, OS de serviço, catálogo de serviços, componentes, app offline, assinatura do cliente, relatório por cliente, conexão com clientes Tenaz), estoque básico e indicadores essenciais. No catálogo, as rotas de plano `S` e as rotas `E` de ADM, ATV, MED, SOL, OS, CHK, estoque básico, IND essencial e MOB.
- Fora do plano Serviços: API v1 e webhooks, multi-unidade, estoque e compras completos, LOC, FRO, NR-13 e o restante do Profissional. A conta pode contratar também um dos planos da seção 2.1 do sumário, e então soma os recursos. Rota fora do plano: `403 plano_nao_inclui_recurso`.
- **Técnico cobrado** na conta da subcontratada: usuário ativo com perfil `tecnico`; `gestor` e `planejador` não pagam. No cliente, o técnico da subcontratada com conta própria **não** conta como técnico cobrado (ele nem tem usuário lá). `GET /v1/empresa/plano` mostra a contagem.
- Uma mesma empresa pode ser cliente numa conexão e subcontratada em outra; o papel vem de cada `conexao_empresa`.
- Na conta da subcontratada, o próprio catálogo de serviços é `servico_fornecedor` com `fornecedor_id = null` (serviço da própria conta, proposta; ver 26.1, A-20).

### 22.3 Técnicos e especialidades

**Criar especialidade** — `POST /v1/especialidades-tecnico` (Gestão):

```json
{
  "nome": "Compressores de parafuso",
  "descricao": "Manutenção corretiva e preventiva de compressores de parafuso até 150 cv.",
  "classe_ativo_ids": ["0193e050-0000-7aaa-8bbb-000000000051"]
}
```

Resposta `201` com `id: "0193e600-0000-7aaa-8bbb-000000000601"`, `ativo: true` e `quantidade_tecnicos: 0`.

**Criar técnico com especialidades** — `POST /v1/pessoas` (Gestão), com o usuário criado na mesma transação:

```
POST /v1/pessoas HTTP/1.1
Authorization: Bearer eyJhbGciOiJFZERTQSIs...
Idempotency-Key: 7b2e4c1a-3d5f-4a6b-9c8d-0e1f2a3b4c5d
Content-Type: application/json
```

```json
{
  "nome": "Bruno Tavares",
  "cpf": "52998224725",
  "telefone": "+5511976543210",
  "unidade_id": "0193e010-0000-7aaa-8bbb-000000000011",
  "custo_hora": "68.00",
  "especialidade_tecnico_ids": [
    "0193e600-0000-7aaa-8bbb-000000000601",
    "0193e600-0000-7aaa-8bbb-000000000602"
  ],
  "usuario": {
    "email": "bruno.tavares@hidrotec.com.br",
    "perfis": ["tecnico"],
    "unidade_ids": ["0193e010-0000-7aaa-8bbb-000000000011"],
    "enviar_convite": true
  }
}
```

**Resposta `201 Created`** (`Location: /v1/pessoas/0193e500-0000-7aaa-8bbb-000000000501`, `ETag: "1"`):

```json
{
  "id": "0193e500-0000-7aaa-8bbb-000000000501",
  "versao": 1,
  "nome": "Bruno Tavares",
  "cpf": "52998224725",
  "telefone": "+5511976543210",
  "unidade_id": "0193e010-0000-7aaa-8bbb-000000000011",
  "fornecedor_id": null,
  "custo_hora": "68.00",
  "especialidades": [
    { "id": "0193e600-0000-7aaa-8bbb-000000000601", "nome": "Compressores de parafuso" },
    { "id": "0193e600-0000-7aaa-8bbb-000000000602", "nome": "Hidráulica industrial" }
  ],
  "usuario": {
    "id": "0193e510-0000-7aaa-8bbb-000000000511",
    "email": "bruno.tavares@hidrotec.com.br",
    "perfis": ["tecnico"],
    "situacao": "convidado",
    "escopo_acesso": "interno"
  },
  "conta_como_tecnico_cobrado": true,
  "ativo": true,
  "criado_em": "2026-11-12T18:00:02Z",
  "atualizado_em": "2026-11-12T18:00:02Z",
  "criado_por_id": "0193e510-0000-7aaa-8bbb-000000000500"
}
```

- `especialidade_tecnico_ids` grava as linhas de `pessoa_especialidade`. `PATCH /v1/pessoas/{id}` com a lista substitui as especialidades; especialidade inativa ou de outra conta devolve `422 validacao_falhou` (`erros[].codigo = "inativo"` ou `"nao_encontrado"`).
- `usuario` é opcional; sem ele, a pessoa existe sem login e pode ser criada depois por `POST /v1/usuarios` com `pessoa_id`. O valor da cobrança muda quando o usuário fica ativo.
- `GET /v1/pessoas?especialidade_tecnico_id=...` lista os técnicos aptos; a tela de distribuição usa esse filtro. Atribuir OS a técnico sem a especialidade pedida na OS (`especialidade_tecnico_ids` da OS, opcional) é permitido, com aviso `tecnico_sem_especialidade` na resposta (proposta).
- Na Fase 2, as capacitações com validade (`/v1/capacitacoes`) continuam valendo também na conta SRV.

### 22.4 Clientes atendidos e equipamentos de clientes

**Criar cliente atendido** — `POST /v1/clientes-servico` (Gestão, plano `S`):

```json
{
  "tipo_pessoa": "juridica",
  "documento": "45678912000155",
  "nome": "Laticínios Serra Azul Ltda.",
  "nome_fantasia": "Serra Azul",
  "contatos": [
    { "nome": "Marta Oliveira", "funcao": "Gerente de manutenção", "email": "marta@serraazul.com.br", "telefone": "+5519998887766", "principal": true, "recebe_relatorio": true }
  ],
  "enderecos_atendimento": [
    {
      "descricao": "Fábrica de Amparo",
      "logradouro": "Rodovia SP-95, km 41",
      "municipio": "Amparo",
      "uf": "SP",
      "cep": "13900000",
      "latitude": -22.7011,
      "longitude": -46.7642,
      "instrucoes_acesso": "Portaria 1; EPI obrigatório na área de produção."
    }
  ],
  "observacao": "Contrato de manutenção mensal dos compressores."
}
```

**Resposta `201 Created`** (`Location: /v1/clientes-servico/0193e100-0000-7aaa-8bbb-000000000102`, `ETag: "1"`):

```json
{
  "id": "0193e100-0000-7aaa-8bbb-000000000102",
  "versao": 1,
  "tipo_pessoa": "juridica",
  "documento": "45678912000155",
  "nome": "Laticínios Serra Azul Ltda.",
  "nome_fantasia": "Serra Azul",
  "contatos": [
    { "id": "0193e101-0000-7aaa-8bbb-000000000111", "nome": "Marta Oliveira", "funcao": "Gerente de manutenção", "email": "marta@serraazul.com.br", "telefone": "+5519998887766", "principal": true, "recebe_relatorio": true }
  ],
  "enderecos_atendimento": [
    { "id": "0193e102-0000-7aaa-8bbb-000000000121", "descricao": "Fábrica de Amparo", "logradouro": "Rodovia SP-95, km 41", "municipio": "Amparo", "uf": "SP", "cep": "13900000", "latitude": -22.7011, "longitude": -46.7642, "instrucoes_acesso": "Portaria 1; EPI obrigatório na área de produção." }
  ],
  "conexao_empresa_id": null,
  "conectado": false,
  "quantidade_equipamentos": 0,
  "ativo": true,
  "criado_em": "2026-11-12T18:05:00Z",
  "atualizado_em": "2026-11-12T18:05:00Z",
  "criado_por_id": "0193e510-0000-7aaa-8bbb-000000000500"
}
```

| Campo | Regra |
|---|---|
| `tipo_pessoa` | `juridica` exige CNPJ (14 dígitos) em `documento`; `fisica`, CPF (11 dígitos). Dígito verificador conferido (`422 validacao_falhou`, `formato_invalido`) |
| `documento` | Único por conta entre clientes ativos (`409 conflito`) |
| `enderecos_atendimento` | Ao menos um. O endereço vira a localização exibida dos equipamentos do cliente e o destino no app |
| `conexao_empresa_id` | Só leitura; preenchido pelo aceite de uma conexão (22.5). Com conexão `ativa`, `nome` e `documento` ficam só leitura |
| Dados pessoais | Cliente `fisica` tem CPF, nome e contatos tratados como dados pessoais (documento 9) |

**Criar equipamento de cliente** — `POST /v1/ativos` (Gestão):

```json
{
  "tipo": "equipamento",
  "propriedade": "de_cliente",
  "cliente_servico_id": "0193e100-0000-7aaa-8bbb-000000000102",
  "endereco_atendimento_id": "0193e102-0000-7aaa-8bbb-000000000121",
  "classe_ativo_id": "0193e050-0000-7aaa-8bbb-000000000051",
  "tag": "LSA-CMP-01",
  "nome": "Compressor de parafuso 40 cv — sala de máquinas",
  "marca": "Schulz",
  "modelo": "SRP 3040",
  "numero_serie": "SRP3040-118822",
  "criticidade": "A",
  "medidores": [ { "tipo": "horas", "unidade": "h", "leitura_inicial": "21450.0" } ]
}
```

**Resposta `201 Created`** (trecho):

```json
{
  "id": "0193e210-0000-7aaa-8bbb-000000000211",
  "versao": 1,
  "tipo": "equipamento",
  "propriedade": "de_cliente",
  "cliente_servico_id": "0193e100-0000-7aaa-8bbb-000000000102",
  "endereco_atendimento_id": "0193e102-0000-7aaa-8bbb-000000000121",
  "tag": "LSA-CMP-01",
  "nome": "Compressor de parafuso 40 cv — sala de máquinas",
  "status_operacional": "operando",
  "codigo_qr": "tnz-a7Kq2",
  "conexao_empresa_id": null,
  "sincronizado": false,
  "medidores": [ { "id": "0193e211-0000-7aaa-8bbb-000000000221", "tipo": "horas", "unidade": "h", "ultima_leitura": "21450.0" } ]
}
```

- Equipamento `de_cliente` tem hierarquia, QR code, medidores, documentos técnicos, checklists e histórico como qualquer ativo. Não entra em KPIs de disponibilidade da subcontratada nem na contagem de veículos cobrados (proposta).
- TAG: na conta da subcontratada, a TAG de equipamento `de_cliente` é única **por cliente** (`cliente_servico_id` + `tag`), porque clientes diferentes podem usar a mesma TAG (proposta; confirmar no documento 2).
- `DELETE /v1/clientes-servico/{id}` com equipamentos ativos ou OS em aberto devolve `409 conflito`; operação nova com cliente inativo devolve `422 cliente_servico_inativo`.

**Equipamento sincronizado.** A cópia criada pela conexão (22.5) chega com `propriedade = de_cliente`, `cliente_servico_id` da conexão, `conexao_empresa_id` preenchido e `sincronizado = true`. Os campos do grupo `ativo_identificacao` (e de `ativo_medidores`, se compartilhado) vêm do cliente e ficam só leitura:

```json
{
  "type": "https://api.tenaz.com.br/erros/ativo_sincronizado_somente_leitura",
  "title": "Este campo vem da conta do cliente",
  "status": 409,
  "detail": "O equipamento CMP-0347 é sincronizado com a conta de Metalúrgica Exemplo. Peça ao cliente para alterar o número de série.",
  "codigo": "ativo_sincronizado_somente_leitura",
  "campos": ["numero_serie"],
  "id_rastreio": "01JD0A1B2C3D4E5F6G7H8J9K0L"
}
```

A subcontratada pode acrescentar à cópia o que é só dela (observações internas, documentos técnicos próprios, checklists próprios); isso não volta para o cliente. Com a conexão encerrada, a cópia vira equipamento `de_cliente` comum (`sincronizado = false`) e passa a ser editável; o histórico fica.

### 22.5 Conexão entre contas: convite, aceite e ciclo de vida

**Estados** (sumário 3.5):

| Ação | Rota | Quem | De | Para | Efeitos |
|---|---|---|---|---|---|
| `convidar` | `POST /v1/conexoes-empresa` | Cliente (Admin) | — | `convidada` | Envia e-mail com o código de convite; webhook nenhum |
| `aceitar` | `POST /v1/conexoes-empresa/{id}/aceitar` | Subcontratada (Admin) | `convidada` | `ativa` | Grava `data_aceite`; liga ou cria o `cliente_servico`; cria as cópias dos ativos compartilhados; webhook `conexao_empresa.ativada` nas duas contas |
| `recusar` | `POST .../recusar` | Subcontratada (Admin) | `convidada` | `encerrada` | Grava `data_encerramento` com motivo; webhook `conexao_empresa.encerrada` no cliente |
| `suspender` | `POST .../suspender` | Qualquer lado (Admin) | `ativa` | `suspensa` | A sincronização para; alterações ficam `pendente` |
| `reativar` | `POST .../reativar` | O lado que suspendeu (Admin) | `suspensa` | `ativa` | As pendências são sincronizadas em ordem; webhook `conexao_empresa.ativada` |
| `encerrar` | `POST .../encerrar` | Qualquer lado (Admin) | `convidada`, `ativa`, `suspensa` | `encerrada` | Final. Para a sincronização sem apagar histórico (22.9); webhook `conexao_empresa.encerrada` |

Ação fora do estado: `409 transicao_invalida`. Todas exigem `If-Match`; motivo obrigatório em `recusar`, `suspender` e `encerrar` (`422 motivo_obrigatorio`). Convite não aceito vence em 14 dias (proposta) e passa a `encerrada` com motivo `convite_expirado`.

**`campos_compartilhados` — lista fechada de grupos (risco R5).** O cliente escolhe os grupos no convite; a sincronização nunca copia campo fora deles. Grupos obrigatórios não podem ser retirados. A lista é fechada: grupo desconhecido devolve `422 validacao_falhou`, e grupo novo só entra por versão nova da API.

| Grupo | Sentido | Conteúdo | Obrigatório | Padrão |
|---|---|---|:-:|:-:|
| `ativo_identificacao` | cliente → subcontratada | TAG, nome, tipo, classe (nome), marca, modelo, número de série, criticidade, localização em texto | Sim | Sim |
| `ativo_medidores` | dois sentidos | Medidores e última leitura (cliente → subcontratada); leituras registradas na OS (subcontratada → cliente, origem `os`) | Não | Sim |
| `ativo_documentos_tecnicos` | cliente → subcontratada | Documentos técnicos marcados como compartilháveis | Não | Não |
| `os_dados_basicos` | cliente → subcontratada | Número, tipo, prioridade, título, descrição, `data_limite`, tarefas, `observacao_fornecedor` | Sim | Sim |
| `os_execucao` | subcontratada → cliente | Transições de execução (`programar`, `iniciar`, `pausar`, `retomar`, `concluir`), `data_programada`, `data_inicio`, `data_conclusao`, `motivo_pausa`, `servico_realizado`, `descricao_servico`, serviços realizados (descrição e quantidade), tarefas concluídas | Sim | Sim |
| `liberacao` | cliente → subcontratada | `encerrar`, `recusar-liberacao` com `motivo_recusa_liberacao`, `cancelar` com motivo, retirada do fornecedor | Sim | Sim |
| `apontamentos_horas` | subcontratada → cliente | Início, fim e `tipo_hora` de cada apontamento | Não | Sim |
| `nomes_tecnicos` | subcontratada → cliente | Nome dos técnicos do responsável e dos apontamentos (sem e-mail, telefone ou CPF) | Não | Sim |
| `componentes` | subcontratada → cliente | Itens de material com `tipo_material`, descrição, quantidade, `numero_serie_retirado`, `numero_serie_instalado` | Não | Sim |
| `fotos_anexos` | dois sentidos | Fotos e anexos da OS marcados com `compartilhar` | Não | Sim |
| `comentarios` | dois sentidos | Comentários da OS marcados com `compartilhar` | Não | Sim |
| `codigos_falha` | subcontratada → cliente | Modo, causa e atividade (pelo texto e pelo código ISO 14224 do catálogo padrão) | Não | Não |
| `assinatura_cliente` | subcontratada → cliente | Assinatura do responsável do cliente coletada em campo | Não | Sim |
| `registro_nr12` | subcontratada → cliente | Rascunho dos campos do item 12.11.2 preenchidos na conclusão | Não | Sim |
| `valores` | subcontratada → cliente | Preço dos serviços e dos componentes; no cliente, vira custo de terceiros previsto | Não | Não |

**Nunca compartilhados no lançamento** (sumário 3.6, RF-OS-021, RF-ATV-013, RF-PRJ-004): os campos extras da OS (`campos_extras`), a causa externa da parada (`causa_externa`, `decisao_exclusao`), a etapa de projeto da OS (`etapa_projeto_id`) e a fase de ciclo de vida do ativo (`fase_ciclo_vida`, `data_entrada_operacao`); também ficam fora o projeto e as etapas, as regras configuráveis e seus disparos, os contratos com clientes e as medições. Eles só passam a trafegar por um grupo novo, decidido no sumário. A OS `implantacao` trafega como qualquer OS, com o `tipo` no grupo `os_dados_basicos`.

Os ativos compartilhados ficam em `ativos_compartilhados`: lista de `ativo_ids`, `classe_ativo_ids`, `localizacao_ids` ou o `contrato_prestacao_id` (Fase 2) cujos ativos cobertos entram (sumário 3.5).

**Convidar** — `POST /v1/conexoes-empresa`, na conta do cliente (perfil `administrador`, `If-Match` não se aplica, `Idempotency-Key`):

```json
{
  "fornecedor_id": "0193b100-0000-7aaa-8bbb-000000000101",
  "email_convite": "rafael@hidrotec.com.br",
  "campos_compartilhados": [
    "ativo_identificacao", "ativo_medidores", "os_dados_basicos", "os_execucao", "liberacao",
    "apontamentos_horas", "nomes_tecnicos", "componentes", "fotos_anexos", "comentarios",
    "assinatura_cliente", "registro_nr12"
  ],
  "ativos_compartilhados": {
    "ativo_ids": [],
    "classe_ativo_ids": ["0192f7b0-0000-7aaa-8bbb-000000000021"],
    "localizacao_ids": [],
    "contrato_prestacao_id": null
  },
  "mensagem": "Vamos passar as OS de compressores para a conta de vocês na Tenaz."
}
```

| Campo | Regra |
|---|---|
| `fornecedor_id` | Fornecedor ativo com papel `subcontratada`, `assistencia_tecnica` ou `oficina` (`422 fornecedor_papel_incompativel`) e CNPJ preenchido |
| `email_convite` | Recebe o código de uso único. Se já houver conta na Tenaz com a mesma raiz de CNPJ, o convite também aparece na lista dela |
| Duplicidade | Conexão não encerrada para o mesmo fornecedor: `409 conexao_duplicada` |

**Resposta `201 Created`** (conta do cliente):

```json
{
  "id": "0193d000-0000-7aaa-8bbb-000000000c01",
  "versao": 1,
  "papel_desta_conta": "cliente",
  "estado": "convidada",
  "fornecedor_id": "0193b100-0000-7aaa-8bbb-000000000101",
  "empresa_remota": { "nome": "Hidrotec Manutenção Industrial Ltda.", "cnpj": "12345678000190", "conta_existente": true },
  "email_convite": "rafael@hidrotec.com.br",
  "campos_compartilhados": ["ativo_identificacao", "ativo_medidores", "os_dados_basicos", "os_execucao", "liberacao", "apontamentos_horas", "nomes_tecnicos", "componentes", "fotos_anexos", "comentarios", "assinatura_cliente", "registro_nr12"],
  "ativos_compartilhados": { "ativo_ids": [], "classe_ativo_ids": ["0192f7b0-0000-7aaa-8bbb-000000000021"], "localizacao_ids": [], "contrato_prestacao_id": null, "quantidade_ativos": 6 },
  "data_convite": "2026-11-12T19:00:00Z",
  "convite_expira_em": "2026-11-26T19:00:00Z",
  "data_aceite": null,
  "data_encerramento": null,
  "transicoes_permitidas": ["encerrar"]
}
```

O código do convite (`cvx_` + 24 caracteres base62, proposta) vai **só no e-mail**, nunca na resposta, e é guardado como hash. O link do e-mail leva o código no fragmento (`https://app.tenaz.com.br/conexoes/0193d000-0000-7aaa-8bbb-000000000c01/aceitar#cvx_...`), como o link de uso único. Quem ainda não tem conta cria a conta do plano Serviços por esse link e aceita em seguida.

**Aceitar** — `POST /v1/conexoes-empresa/0193d000-0000-7aaa-8bbb-000000000c01/aceitar`, na conta da subcontratada (perfil `administrador`), `If-Match: "1"`:

```json
{
  "codigo_convite": "cvx_Hq7Lm2Pz9Wd4Ks1Vb8Nr3Tc6",
  "cliente_servico_id": null,
  "criar_cliente_servico": {
    "endereco_atendimento": { "descricao": "Planta 1", "logradouro": "Av. Industrial, 1200", "municipio": "Guarulhos", "uf": "SP", "cep": "07220010" }
  }
}
```

**Resposta `200 OK`** (conta da subcontratada):

```json
{
  "id": "0193d000-0000-7aaa-8bbb-000000000c01",
  "versao": 2,
  "papel_desta_conta": "prestador",
  "estado": "ativa",
  "cliente_servico_id": "0193e100-0000-7aaa-8bbb-000000000101",
  "empresa_remota": { "nome": "Metalúrgica Exemplo S.A.", "cnpj": "11222333000181" },
  "campos_compartilhados": ["ativo_identificacao", "ativo_medidores", "os_dados_basicos", "os_execucao", "liberacao", "apontamentos_horas", "nomes_tecnicos", "componentes", "fotos_anexos", "comentarios", "assinatura_cliente", "registro_nr12"],
  "ativos_compartilhados": { "quantidade_ativos": 6, "estado_sincronizacao": "pendente" },
  "data_convite": "2026-11-12T19:00:00Z",
  "data_aceite": "2026-11-12T19:20:41Z",
  "transicoes_permitidas": ["suspender", "encerrar"]
}
```

- O aceite confere o código (hash, tempo constante) e a raiz do CNPJ da conta contra a do fornecedor convidado. Falha: `422 convite_conexao_invalido` com `motivo`; 5 códigos errados bloqueiam o convite (seção 5).
- `cliente_servico_id` liga a conexão a um cliente já cadastrado (com o mesmo documento); sem ele, `criar_cliente_servico` cria o cliente com os dados da conta do cliente. A subcontratada vê só `campos_compartilhados` e `ativos_compartilhados.quantidade_ativos`; o filtro de ativos do cliente (classes, localizações) é dado do cliente e não sai para ela.
- Ampliar `campos_compartilhados` depois do aceite (`PATCH` do cliente) deixa a conexão `ativa`, com `alteracao_pendente`, até um novo `aceitar` com `{ "confirmar_alteracao": true }` (sem código); reduzir vale na hora e para de enviar o grupo retirado.

### 22.6 Como a sincronização funciona (ADR-009)

A API pública não tem rota de sincronização entre contas: nenhum cliente da API lê ou grava na outra conta. O processo segue o ADR-009 (documento 5) e o risco R5:

| Princípio | Como aparece na API |
|---|---|
| Cópia por conta | Cada conta tem a sua OS, o seu equipamento e o seu registro de vínculo, consultados só pelas rotas da própria conta |
| Processo com papel próprio | Um worker com papel de banco próprio grava em cada conta numa **transação separada**, com `app.empresa_id` daquela conta e o RLS ligado. O autor nos eventos e na auditoria é `{ "tipo": "sincronizacao_conexao", "id": "<conexao_empresa_id>" }` |
| Lista fechada de campos | Só os grupos de `campos_compartilhados` saem da conta de origem; o filtro é aplicado na saída (outbox) e conferido de novo na entrada |
| Outbox nos dois lados | A alteração é gravada na conta de origem com um registro de saída na mesma transação; o worker entrega e marca `sincronizado`. Uma alteração nunca se perde se a outra conta estiver indisponível |
| Idempotência | Cada mensagem entre contas tem `mensagem_id` único e `sequencia` crescente por vínculo e sentido; o destino descarta `mensagem_id` repetido e aplica na ordem da sequência; repetir não duplica nada (2_ARQUITETURA_BANCO_DADOS.md, seção 2.10) |
| Ordem | Por vínculo, as alterações são aplicadas na ordem de `versao` da OS de origem |
| Auditoria | Evento de auditoria nas duas contas: o que saiu, o que entrou, com os grupos |
| Sem leitura cruzada | Nenhuma consulta junta dados de duas empresas; o worker lê só a outbox de uma conta e grava na outra |

**Estado do vínculo.** `pendente` (alteração aguardando entrega), `sincronizado` (última alteração gravada nos dois lados), `erro` (falhou depois das novas tentativas, com espera exponencial por até 24 horas, como os webhooks). Em `erro`, o vínculo emite `vinculo_os.erro`, mostra `ultimo_erro` e pode ser reprocessado (`POST /v1/vinculos-os/{id}/reprocessar`). Meta de atraso, com as duas contas disponíveis: até 60 segundos (proposta).

**Conflitos.** Cada campo tem um dono (22.8), então os dois lados nunca escrevem o mesmo campo. Transições do cliente que dependem da etapa de execução (`cancelar`, `retirar-fornecedor`) são validadas contra o estado da cópia na subcontratada no momento da entrega; se a cópia já estiver em `em_execucao` ou `concluida`, a entrega falha com `ultimo_erro.codigo = "transicao_invalida"`, a OS do cliente volta ao estado sincronizado e o gestor do cliente é notificado (proposta).

### 22.7 OS de serviço na conta da subcontratada

Exemplos para cliente **sem** conta na Tenaz (Laticínios Serra Azul, 22.4). O fluxo é o mesmo da seção 8, com os campos de serviço e de componentes.

**Abrir OS de serviço com técnico e serviço a realizar** — `POST /v1/ordens-servico` (Execução):

```json
{
  "tipo": "corretiva_planejada",
  "ativo_id": "0193e210-0000-7aaa-8bbb-000000000211",
  "titulo": "Compressor desarmando por alta temperatura",
  "descricao": "Cliente relata desarme após 40 minutos de operação.",
  "descricao_servico": "Diagnosticar alta temperatura, trocar elemento separador e óleo, limpar radiador.",
  "prioridade": "alta",
  "data_limite": "2026-11-14T21:00:00Z",
  "responsavel_id": "0193e510-0000-7aaa-8bbb-000000000511",
  "especialidade_tecnico_ids": ["0193e600-0000-7aaa-8bbb-000000000601"],
  "servicos": [
    { "servico_fornecedor_id": "0193e700-0000-7aaa-8bbb-000000000701", "quantidade": "3" },
    { "servico_fornecedor_id": "0193e700-0000-7aaa-8bbb-000000000702", "quantidade": "1" }
  ]
}
```

**Resposta `201 Created`** (`ETag: "1"`), trecho:

```json
{
  "id": "0193e310-0000-7aaa-8bbb-000000000311",
  "numero": "OS-2026-000198",
  "versao": 1,
  "status": "aberta",
  "origem": "manual",
  "tipo": "corretiva_planejada",
  "prioridade": "alta",
  "ativo_id": "0193e210-0000-7aaa-8bbb-000000000211",
  "cliente_servico_id": "0193e100-0000-7aaa-8bbb-000000000102",
  "descricao_servico": "Diagnosticar alta temperatura, trocar elemento separador e óleo, limpar radiador.",
  "servicos": [
    { "id": "0193e311-0000-7aaa-8bbb-000000000321", "servico_fornecedor_id": "0193e700-0000-7aaa-8bbb-000000000701", "descricao": "Mão de obra técnica — compressores", "quantidade": "3", "unidade_preco": "hora", "preco_referencia": "220.00", "valor_total": "660.00", "realizado": false },
    { "id": "0193e311-0000-7aaa-8bbb-000000000322", "servico_fornecedor_id": "0193e700-0000-7aaa-8bbb-000000000702", "descricao": "Visita técnica — região de Campinas", "quantidade": "1", "unidade_preco": "visita", "preco_referencia": "180.00", "valor_total": "180.00", "realizado": false }
  ],
  "responsavel_id": "0193e510-0000-7aaa-8bbb-000000000511",
  "executores_ids": ["0193e510-0000-7aaa-8bbb-000000000511"],
  "data_reconhecimento": null,
  "vinculo": null,
  "avisos": [],
  "transicoes_permitidas": ["planejar", "programar", "reconhecer", "cancelar"]
}
```

- `cliente_servico_id` da OS é só leitura, herdado do equipamento. OS de equipamento `de_cliente` exige que o cliente esteja ativo (`422 cliente_servico_inativo`).
- `servicos[].servico_fornecedor_id` precisa ser do catálogo da própria conta; `preco_referencia` padrão é o do catálogo e pode ser alterado por item.
- Sem `responsavel_id`, a OS fica na fila da equipe e qualquer técnico com acesso pode `assumir` (8.2).
- Pelo app, o técnico recebe a OS no próximo pull (seção 12) e a vê offline.

**Apontar início e fim** — o técnico inicia, aponta e conclui. `POST /v1/ordens-servico/0193e310-.../iniciar`, `If-Match: "2"`:

```json
{
  "data_inicio": "2026-11-13T13:10:00Z",
  "ativo_parado": true,
  "leituras": [ { "medidor_id": "0193e211-0000-7aaa-8bbb-000000000221", "valor": "21512.4" } ]
}
```

Resposta `200`, `ETag: "3"`, com `status: "em_execucao"`, `data_inicio: "2026-11-13T13:10:00Z"` e um apontamento aberto automaticamente (proposta). O técnico também pode apontar o período manualmente — `POST /v1/ordens-servico/0193e310-.../apontamentos`:

```json
{
  "id": "0193e320-0000-7aaa-8bbb-000000000331",
  "usuario_id": "0193e510-0000-7aaa-8bbb-000000000511",
  "data_inicio": "2026-11-13T13:10:00Z",
  "data_fim": "2026-11-13T15:55:00Z",
  "tipo_hora": "normal",
  "observacao": "Diagnóstico, troca do separador e do óleo, limpeza do radiador."
}
```

**Resposta `201 Created`**:

```json
{
  "id": "0193e320-0000-7aaa-8bbb-000000000331",
  "ordem_servico_id": "0193e310-0000-7aaa-8bbb-000000000311",
  "usuario_id": "0193e510-0000-7aaa-8bbb-000000000511",
  "pessoa_nome": "Bruno Tavares",
  "data_inicio": "2026-11-13T13:10:00Z",
  "data_fim": "2026-11-13T15:55:00Z",
  "duracao_min": 165,
  "tipo_hora": "normal",
  "custo_hora": "68.00",
  "custo": "187.00",
  "criado_em": "2026-11-13T15:56:10Z"
}
```

**Registrar componentes substituídos** — `POST /v1/ordens-servico/0193e310-.../pecas` (Execução), um item por chamada ou em lote com `itens`:

```json
{
  "itens": [
    {
      "tipo_material": "substituido",
      "item_estoque_id": "0193e800-0000-7aaa-8bbb-000000000801",
      "descricao": "Elemento separador ar/óleo SRP 3040",
      "quantidade": "1",
      "numero_serie_retirado": "SEP-2023-00412",
      "numero_serie_instalado": "SEP-2026-01877",
      "motivo_substituicao": "Saturado; diferencial de pressão 1,2 bar"
    },
    {
      "tipo_material": "consumido",
      "item_estoque_id": "0193e800-0000-7aaa-8bbb-000000000805",
      "descricao": "Óleo sintético para compressor (L)",
      "quantidade": "12"
    },
    {
      "tipo_material": "fornecido_pelo_cliente",
      "item_estoque_id": null,
      "descricao": "Filtro de ar (estoque do cliente)",
      "quantidade": "1",
      "numero_serie_instalado": null
    }
  ]
}
```

**Resposta `201 Created`**:

```json
{
  "dados": [
    { "id": "0193e330-0000-7aaa-8bbb-000000000341", "tipo_material": "substituido", "item_estoque_id": "0193e800-0000-7aaa-8bbb-000000000801", "descricao": "Elemento separador ar/óleo SRP 3040", "quantidade": "1", "numero_serie_retirado": "SEP-2023-00412", "numero_serie_instalado": "SEP-2026-01877", "custo_unitario": "412.0000", "custo_total": "412.00", "movimentacao_estoque_id": "0193e331-0000-7aaa-8bbb-000000000351" },
    { "id": "0193e330-0000-7aaa-8bbb-000000000342", "tipo_material": "consumido", "item_estoque_id": "0193e800-0000-7aaa-8bbb-000000000805", "descricao": "Óleo sintético para compressor (L)", "quantidade": "12", "custo_unitario": "38.5000", "custo_total": "462.00", "movimentacao_estoque_id": "0193e331-0000-7aaa-8bbb-000000000352" },
    { "id": "0193e330-0000-7aaa-8bbb-000000000343", "tipo_material": "fornecido_pelo_cliente", "item_estoque_id": null, "descricao": "Filtro de ar (estoque do cliente)", "quantidade": "1", "custo_unitario": null, "custo_total": "0.00", "movimentacao_estoque_id": null }
  ]
}
```

| `tipo_material` | Regra |
|---|---|
| `consumido` | Material gasto (óleo, vedação). Com `item_estoque_id`, gera `saida_os` do estoque da conta |
| `substituido` | Peça retirada e peça instalada. `numero_serie_retirado` e `numero_serie_instalado` quando houver; item com `controla_serie = true` exige `pecas_serializadas_ids` (Fase 2). Com `item_estoque_id`, gera `saida_os` da peça instalada |
| `fornecido_pelo_cliente` | Peça do cliente instalada pela subcontratada. Não movimenta estoque nem tem custo; entra no relatório |

Sem `item_estoque_id`, o item é só descritivo (`descricao` obrigatória). Saldo insuficiente: `422 estoque_insuficiente`; saldo em trânsito: `422 estoque_em_transito`.

**Concluir com a assinatura do cliente** — `POST /v1/ordens-servico/{id}/assinaturas` com `{ "papel": "cliente_atendido", "nome": "Marta Oliveira", "documento": null, "assinatura_anexo_id": "0193e340-0000-7aaa-8bbb-0000000003a1" }` e depois `POST .../concluir` com `data_conclusao`, `servico_realizado` e as fotos (8.3). Sem vínculo com conta de cliente, a própria subcontratada encerra a OS com `os.encerrar` (8.2).

**Gerar o relatório em PDF** — `GET /v1/ordens-servico/0193e310-0000-7aaa-8bbb-000000000311/relatorio-pdf?incluir_fotos=true&incluir_valores=false` (Leitura):

```
HTTP/1.1 200 OK
Content-Type: application/pdf
Content-Disposition: attachment; filename="OS-2026-000198.pdf"
Cache-Control: private, no-store
ETag: "7"
Tenaz-Hash-Sha256: 3e9b1f0c7a2d4e58b6c1e9a07d3f2b4c8e5a1d6f9b0c2e4a7d8f1b3c5e7a9d0f
```

| Parâmetro | Padrão | Efeito |
|---|---|---|
| `incluir_fotos` | `true` | Fotos da OS, até 24 por relatório, redimensionadas (proposta) |
| `incluir_valores` | `false` | Preço dos serviços e dos componentes |
| `incluir_apontamentos` | `true` | Horas por técnico; com `false`, só início, fim e total |

Conteúdo: cabeçalho da conta que presta o serviço (nome, CNPJ, logotipo); cliente atendido e endereço de atendimento; equipamento (TAG, nome, marca, modelo, número de série, leituras); número, tipo, datas de abertura, início, conclusão e encerramento; `descricao_servico`, serviços realizados e `servico_realizado`; componentes por `tipo_material`, com números de série retirado e instalado; apontamentos (técnico, início, fim, duração); fotos; assinaturas (técnico e cliente), com carimbo de tempo e o hash do conteúdo; QR code que leva à verificação do hash (proposta). O mesmo hash vai no cabeçalho `Tenaz-Hash-Sha256`.

- OS em qualquer estado gera o relatório; antes de `concluida`, o PDF traz a marca "Prévia". OS com mais de 24 fotos ou em lote (`POST /v1/relatorios/servicos_por_cliente/gerar`, com `cliente_servico_id` e período) devolve `202` com processamento (2.14).
- Na conta do cliente, o relatório de uma OS vinculada usa só os dados que a conta dele recebeu.

### 22.8 Ciclo da OS vinculada: do cliente para a subcontratada e de volta

**Quem altera o quê.** Numa OS com vínculo, cada campo e cada ação tem um dono. O outro lado lê (se o grupo for compartilhado) e recebe erro ao tentar alterar.

| Campo ou ação | Dono | Na conta do cliente | Na conta da subcontratada |
|---|---|---|---|
| `tipo`, `prioridade`, `titulo`, `descricao`, `data_limite`, `observacao_fornecedor` | Cliente | Altera até `concluida` | Lê. `PATCH` → `409 os_vinculada_somente_leitura` |
| `ativo_id` | Cliente | Fixo depois do vínculo (`422 campo_somente_leitura`) | Lê (cópia sincronizada) |
| Tarefas: criar, remover, editar texto | Cliente | Antes de `em_execucao` | `409 os_vinculada_somente_leitura` |
| Tarefas: marcar concluída | Subcontratada | Lê | Marca |
| `responsavel_custo`, `justificativa_responsavel_custo`, custos, `contrato_locacao_id` | Cliente | Altera (8.6) | Não recebe |
| `responsavel_id`, `executores_ids`, `programar`, `desprogramar`, `assumir`, `reconhecer` | Subcontratada | Lê os nomes, se `nomes_tecnicos`. Ação → `409 os_vinculada_somente_leitura` | Executa |
| `iniciar`, `pausar`, `retomar`, `concluir`; `data_inicio`, `data_conclusao`, `motivo_pausa`, `servico_realizado` | Subcontratada | Ação → `409 os_vinculada_somente_leitura` | Executa |
| `descricao_servico`, `servicos` (`servico_os`) | Subcontratada | Lê | Altera |
| Apontamentos | Subcontratada | Lê, se `apontamentos_horas`. `PATCH`/`DELETE` → `409 os_vinculada_somente_leitura` | Cria e corrige até o encerramento |
| Componentes (`pecas` com `tipo_material`) | Subcontratada | Lê, se `componentes`. Pode lançar consumo do **próprio** estoque, que fica só na conta dele | Registra |
| Leituras de medidor na OS | Subcontratada | Recebe, se `ativo_medidores` (origem `os`) | Registra |
| Comentários e anexos | Quem criou | Cria; marca `compartilhar` | Cria; marca `compartilhar`. Sem o grupo → `422 campo_nao_compartilhado` |
| Códigos de falha | Subcontratada | Lê, se `codigos_falha`; senão preenche os seus | Preenche |
| Assinatura do cliente em campo | Subcontratada coleta | Recebe, se `assinatura_cliente` | Registra |
| `registro_nr12` | Subcontratada preenche o rascunho; cliente completa no encerramento | Completa e assina em `encerrar` (8.4) | Preenche na conclusão |
| `encerrar`, `recusar-liberacao` | Cliente (`os.encerrar`) | Executa | `409 os_vinculada_somente_leitura` |
| `cancelar`, `retirar-fornecedor` | Cliente | Executa (estados da 8.2) | `409 os_vinculada_somente_leitura`; a subcontratada pede pelo comentário |

O objeto `vinculo` na OS (as duas contas):

| Campo | Descrição |
|---|---|
| `id` | Vínculo como visto por esta conta (não serve para consultar a outra) |
| `conexao_empresa_id` | Conexão (id da plataforma, igual nas duas contas) |
| `lado` | `cliente` ou `prestador`: o papel desta conta |
| `lado_origem` | `cliente` (a OS nasceu no cliente) ou `prestador` (a subcontratada abriu) |
| `empresa_remota` | Nome e CNPJ da outra conta |
| `numero_os_remota` | Número da OS na outra conta |
| `estado_sincronizacao` | `pendente`, `sincronizado`, `erro` |
| `sincronizado_em`, `ultimo_erro` | Última entrega bem-sucedida e último erro (`codigo`, `mensagem`, `em`) |
| `grupos_compartilhados` | `campos_compartilhados` da conexão no momento |
| `campos_somente_leitura` | Campos da OS que esta conta não pode alterar |

**Passo 1 — o cliente atribui a OS ao fornecedor conectado.** Conta do cliente: `POST /v1/ordens-servico/0193d200-0000-7aaa-8bbb-000000000e01/atribuir-fornecedor`, `If-Match: "2"`:

```json
{
  "fornecedor_id": "0193b100-0000-7aaa-8bbb-000000000101",
  "data_limite": "2026-11-14T21:00:00Z",
  "observacao_fornecedor": "Compressor 2 desarmando; acesso pela portaria 2."
}
```

Resposta `200` (`ETag: "3"`), trecho:

```json
{
  "id": "0193d200-0000-7aaa-8bbb-000000000e01",
  "numero": "OS-2026-004950",
  "versao": 3,
  "status": "aberta",
  "executor_tipo": "fornecedor",
  "fornecedor_id": "0193b100-0000-7aaa-8bbb-000000000101",
  "responsavel_id": null,
  "vinculo": {
    "id": "0193d100-0000-7aaa-8bbb-000000000d01",
    "conexao_empresa_id": "0193d000-0000-7aaa-8bbb-000000000c01",
    "lado": "cliente",
    "lado_origem": "cliente",
    "empresa_remota": { "nome": "Hidrotec Manutenção Industrial Ltda.", "cnpj": "12345678000190" },
    "numero_os_remota": null,
    "estado_sincronizacao": "pendente",
    "sincronizado_em": null,
    "ultimo_erro": null
  },
  "transicoes_permitidas": ["cancelar"]
}
```

`transicoes_permitidas` do cliente passa a ter só as ações que são dele.

**Passo 2 — a OS nasce na subcontratada.** O processo de sincronização cria a cópia na conta da Hidrotec. Conta da subcontratada: `GET /v1/ordens-servico/0193e300-0000-7aaa-8bbb-000000000301`, trecho:

```json
{
  "id": "0193e300-0000-7aaa-8bbb-000000000301",
  "numero": "OS-2026-000231",
  "versao": 1,
  "status": "aberta",
  "origem": "cliente_conectado",
  "tipo": "corretiva_planejada",
  "prioridade": "alta",
  "titulo": "Compressor 2 desarmando por temperatura",
  "ativo_id": "0193e200-0000-7aaa-8bbb-000000000201",
  "cliente_servico_id": "0193e100-0000-7aaa-8bbb-000000000101",
  "data_limite": "2026-11-14T21:00:00Z",
  "observacao_fornecedor": "Compressor 2 desarmando; acesso pela portaria 2.",
  "executor_tipo": "interno",
  "responsavel_id": null,
  "vinculo": {
    "id": "0193e400-0000-7aaa-8bbb-000000000401",
    "conexao_empresa_id": "0193d000-0000-7aaa-8bbb-000000000c01",
    "lado": "prestador",
    "lado_origem": "cliente",
    "empresa_remota": { "nome": "Metalúrgica Exemplo S.A.", "cnpj": "11222333000181" },
    "numero_os_remota": "OS-2026-004950",
    "estado_sincronizacao": "sincronizado",
    "sincronizado_em": "2026-11-13T12:00:09Z",
    "campos_somente_leitura": ["tipo", "prioridade", "titulo", "descricao", "data_limite", "observacao_fornecedor", "ativo_id", "tarefas.descricao"]
  },
  "transicoes_permitidas": ["planejar", "programar", "reconhecer", "assumir"]
}
```

Na conta da subcontratada, a OS é executada por técnico próprio (`executor_tipo = interno`) e recebe `origem = cliente_conectado` (sumário 3.6); o vínculo traz também `vinculo.lado_origem = cliente` (ver 26.1, A-21, resolvido). Uma OS `implantacao` do cliente chega com o mesmo `tipo`, sem `etapa_projeto_id`, e a cópia não passa pela validação de fase do ativo, que não é compartilhada (ver 26.1, A-31). Webhooks na conta da subcontratada: `ordem_servico.criada` e `vinculo_os.sincronizado`.

**Passo 3 — o técnico assume e executa** (conta da subcontratada). `POST .../0193e300-.../assumir` (`If-Match: "1"`), depois `iniciar`, `pecas` com `tipo_material = substituido`, `apontamentos` e `concluir`, iguais aos exemplos da 22.7. Cada passo gera uma entrega para o cliente com os grupos de `os_execucao`, `apontamentos_horas`, `componentes` e `fotos_anexos`.

**Passo 4 — o cliente recebe a conclusão.** Conta do cliente: `GET /v1/ordens-servico/0193d200-0000-7aaa-8bbb-000000000e01?expandir=apontamentos,pecas`, trecho:

```json
{
  "id": "0193d200-0000-7aaa-8bbb-000000000e01",
  "numero": "OS-2026-004950",
  "versao": 9,
  "status": "concluida",
  "data_reconhecimento": "2026-11-13T12:04:30Z",
  "data_inicio": "2026-11-13T13:10:00Z",
  "data_conclusao": "2026-11-13T15:58:00Z",
  "responsavel_nome_remoto": "Bruno Tavares",
  "servico_realizado": "Trocado o elemento separador e o óleo; radiador limpo. Teste de 40 minutos em carga sem desarme.",
  "apontamentos": [
    { "id": "0193d210-0000-7aaa-8bbb-000000000e11", "pessoa_nome": "Bruno Tavares", "data_inicio": "2026-11-13T13:10:00Z", "data_fim": "2026-11-13T15:55:00Z", "tipo_hora": "normal", "origem_remota": true }
  ],
  "pecas": [
    { "id": "0193d220-0000-7aaa-8bbb-000000000e21", "tipo_material": "substituido", "descricao": "Elemento separador ar/óleo", "quantidade": "1", "numero_serie_retirado": "SEP-2023-00412", "numero_serie_instalado": "SEP-2026-01877", "origem_remota": true }
  ],
  "custo_mao_obra": "0.00",
  "custo_terceiros": "0.00",
  "vinculo": {
    "id": "0193d100-0000-7aaa-8bbb-000000000d01",
    "lado": "cliente",
    "lado_origem": "cliente",
    "numero_os_remota": "OS-2026-000231",
    "estado_sincronizacao": "sincronizado",
    "sincronizado_em": "2026-11-13T15:58:07Z"
  },
  "transicoes_permitidas": ["encerrar", "recusar-liberacao"]
}
```

`data_reconhecimento`, `data_inicio` e `data_conclusao` vêm da subcontratada e alimentam o MTTA, o MTTR e o SLA do fornecedor no cliente. As horas do terceiro não entram em `custo_mao_obra`; com o grupo `valores`, o preço vira `custo_terceiros` previsto. Webhooks no cliente: `vinculo_os.sincronizado`, `ordem_servico.concluida`.

**Passo 5 — o cliente libera.** Conta do cliente: `POST /v1/ordens-servico/0193d200-.../encerrar`, `If-Match: "9"`, por usuário com `os.encerrar` (8.4 quando NR-12). Resposta `200` com `status: "encerrada"`. A sincronização leva o encerramento à cópia, que fica `encerrada` e imutável na conta da subcontratada (webhook `vinculo_os.sincronizado` lá). Se o cliente recusar (`recusar-liberacao` com `motivo_recusa_liberacao`), as duas cópias voltam a `em_execucao` com o motivo, e a subcontratada retoma.

**OS aberta pela subcontratada** (`lado_origem = prestador`). Na conta da subcontratada, `POST /v1/ordens-servico` em equipamento sincronizado:

```json
{
  "tipo": "corretiva_planejada",
  "ativo_id": "0193e200-0000-7aaa-8bbb-000000000201",
  "titulo": "Vazamento no resfriador de óleo encontrado na visita",
  "descricao": "Visto durante a OS-2026-000231. Recomenda troca do resfriador.",
  "prioridade": "media",
  "enviar_ao_cliente": true
}
```

A OS nasce `aberta` na subcontratada, com `vinculo.lado_origem = prestador` e `estado_sincronizacao = pendente`. Na conta do cliente, a sincronização cria a OS `aberta` com `origem = prestador`, `executor_tipo = fornecedor`, `fornecedor_id` da conexão e o `responsavel_custo` padrão da conta do cliente; ela segue a triagem do cliente (sumário 3.5). A partir daí, os donos são os da tabela acima: o cliente pode alterar prioridade, planejar, cancelar ou liberar; a subcontratada executa. `enviar_ao_cliente` é `true` por padrão em equipamento sincronizado; com `false`, a OS fica só na conta da subcontratada. Com a conexão `suspensa` ou `encerrada`, `enviar_ao_cliente: true` devolve `409 conexao_inativa`.

**Erros do ciclo.**

O cliente tenta corrigir um apontamento vindo da subcontratada — `PATCH /v1/ordens-servico/0193d200-.../apontamentos/0193d210-...`:

```json
{
  "type": "https://api.tenaz.com.br/erros/os_vinculada_somente_leitura",
  "title": "Este dado pertence à outra conta da OS vinculada",
  "status": 409,
  "detail": "O apontamento foi registrado pela Hidrotec Manutenção Industrial Ltda. na OS-2026-000231. Peça a correção pelo comentário da OS.",
  "codigo": "os_vinculada_somente_leitura",
  "lado_dono": "prestador",
  "campos": ["apontamentos"],
  "id_rastreio": "01JD0B2C3D4E5F6G7H8J9K0L1M"
}
```

A subcontratada tenta encerrar a cópia — `POST /v1/ordens-servico/0193e300-.../encerrar`:

```json
{
  "type": "https://api.tenaz.com.br/erros/os_vinculada_somente_leitura",
  "title": "Este dado pertence à outra conta da OS vinculada",
  "status": 409,
  "detail": "A liberação da OS-2026-000231 é feita pelo cliente Metalúrgica Exemplo S.A. na OS-2026-004950.",
  "codigo": "os_vinculada_somente_leitura",
  "lado_dono": "cliente",
  "campos": ["encerrar"],
  "id_rastreio": "01JD0B3D4E5F6G7H8J9K0L1M2N"
}
```

A subcontratada marca um comentário para o cliente numa conexão sem o grupo `comentarios` — `POST /v1/ordens-servico/0193e300-.../comentarios` com `{ "texto": "Resfriador com corrosão.", "compartilhar": true }`:

```json
{
  "type": "https://api.tenaz.com.br/erros/campo_nao_compartilhado",
  "title": "Este dado não é compartilhado pela conexão",
  "status": 422,
  "detail": "A conexão com Metalúrgica Exemplo S.A. não compartilha comentários. Grave o comentário sem compartilhar ou peça ao cliente para incluir o grupo.",
  "codigo": "campo_nao_compartilhado",
  "grupos": ["comentarios"],
  "campos": ["compartilhar"],
  "id_rastreio": "01JD0B4E5F6G7H8J9K0L1M2N3P"
}
```

O mesmo erro sai quando o cliente atribui OS com `exigir_assinatura_cliente: true` e a conexão não tem `assinatura_cliente`.

O cliente atribui OS com a conexão suspensa:

```json
{
  "type": "https://api.tenaz.com.br/erros/conexao_inativa",
  "title": "A conexão com a conta do fornecedor não está ativa",
  "status": 409,
  "detail": "A conexão com Hidrotec Manutenção Industrial Ltda. está suspensa desde 20/11/2026. Reative a conexão ou atribua sem sincronização.",
  "codigo": "conexao_inativa",
  "estado_conexao": "suspensa",
  "id_rastreio": "01JD0B5F6G7H8J9K0L1M2N3P4Q"
}
```

**Vínculo com erro.** `GET /v1/vinculos-os?estado_sincronizacao=erro` lista os vínculos parados; `GET /v1/vinculos-os/{id}` traz o histórico de entregas e `ultimo_erro` (por exemplo, `{ "codigo": "transicao_invalida", "mensagem": "Cancelamento recusado: a OS-2026-000231 já está em execução.", "em": "2026-11-13T13:12:00Z" }`); `POST /v1/vinculos-os/{id}/reprocessar` (`202`) reenfileira depois da correção.

### 22.9 Suspensão e encerramento da conexão

| Situação | OS vinculadas abertas | Equipamentos sincronizados | Histórico |
|---|---|---|---|
| `suspensa` | Continuam nas duas contas; alterações ficam `pendente` e saem na reativação, em ordem | Ficam só leitura na subcontratada | Intacto |
| `encerrada` | O vínculo fecha: cada lado fica com a sua cópia no último estado sincronizado. A OS do cliente volta a `executor_tipo = fornecedor` sem sincronização (o cliente pode retirar o fornecedor, cancelar ou usar usuário `prestador` ou link); a cópia da subcontratada vira OS comum, que ela pode concluir e encerrar | Viram `de_cliente` comuns (`sincronizado = false`), editáveis | Intacto nas duas contas; nada é apagado (sumário 3.5) |

Encerrar a conexão não cancela OS em nenhum lado. As duas contas recebem `conexao_empresa.encerrada` com `motivo` e a lista de vínculos que ficaram abertos (`vinculos_abertos`).

### 22.10 API pública para integradores da conta da subcontratada

A conta da subcontratada é uma empresa como as outras: o integrador dela (ERP, BI, sistema de faturamento) usa token de API **da conta dela**, e vê só o que está gravado nela.

| A API expõe ao integrador da subcontratada | Rotas e escopos |
|---|---|
| Clientes atendidos, especialidades, serviços a realizar e relatório da OS | `/v1/clientes-servico`, `/v1/especialidades-tecnico`, `/v1/ordens-servico/{id}/servicos`, `relatorio-pdf`; `servicos:leitura` / `servicos:escrita` |
| Técnicos (pessoas e usuários) | `/v1/pessoas`, `/v1/usuarios`; `pessoas:*`, `empresa:leitura` |
| Equipamentos de clientes, inclusive as cópias sincronizadas (só os grupos compartilhados) | `/v1/ativos?propriedade=de_cliente`; `ativos:*` (campos do cliente só leitura, 22.4) |
| OS próprias e vinculadas, apontamentos, componentes, comentários, fotos e assinaturas | `/v1/ordens-servico`; `ordens_servico:*`, com as mesmas regras de dono da 22.8 |
| Conexões, ativos compartilhados e vínculos, só leitura | `/v1/conexoes-empresa`, `/v1/vinculos-os`; `conexoes:leitura` |
| Webhooks da própria conta, inclusive `conexao_empresa.*` e `vinculo_os.*` | `/v1/webhooks`; `webhooks:gerenciar` |
| Exportação completa da conta (seção 17) | `exportacao:empresa` |

| A API **não** expõe | Por quê |
|---|---|
| Dados da conta do cliente além dos grupos compartilhados (custos, `responsavel_custo`, outros fornecedores, outros ativos, usuários) | Lista fechada de campos (22.5) e isolamento por empresa (4.7) |
| Ids da outra conta (`fornecedor_id` do cliente, ids de OS, de ativo e de usuários de lá) | Só números e nomes (22.1) |
| Ações de conexão (`convidar`, `aceitar`, `recusar`, `suspender`, `reativar`, `encerrar`, alterar grupos) | Ações administrativas reservadas ao usuário `administrador` (4.4) |
| Escrita em campo do cliente numa OS vinculada ou em equipamento sincronizado | `409 os_vinculada_somente_leitura` e `409 ativo_sincronizado_somente_leitura` |
| Eventos da conta do cliente | Cada conta assina só os próprios webhooks (13.2) |

**Plano.** Tokens de API e webhooks não estão no plano Serviços (sumário 2.2). A conta da subcontratada que contratar também o Profissional ganha tokens e webhooks com os limites dele (seção 5). Proposta a validar: incluir no plano Serviços um token de **leitura** (`servicos:leitura`, `ordens_servico:leitura`, `ativos:leitura`, `conexoes:leitura`) para o faturamento da subcontratada (26.1, A-22).

Do lado do cliente, o integrador vê a OS vinculada como OS atribuída a fornecedor, com o objeto `vinculo`, e recebe os dados da subcontratada nos campos de execução, apontamentos e componentes (marcados `origem_remota: true`), sem ids da conta dela.

### 22.11 Resumo de erros e eventos do SRV

| Código | HTTP | Quando |
|---|---|---|
| `conexao_inativa` | 409 | Operação que exige conexão `ativa` com conexão `convidada`, `suspensa` ou `encerrada` |
| `conexao_duplicada` | 409 | Segunda conexão não encerrada para o mesmo fornecedor |
| `convite_conexao_invalido` | 422 | Código errado ou vencido, ou raiz de CNPJ divergente |
| `campo_nao_compartilhado` | 422 | Dado marcado para ir à outra conta, mas o grupo não está em `campos_compartilhados` |
| `os_vinculada_somente_leitura` | 409 | Campo ou ação que pertence ao outro lado da OS vinculada |
| `ativo_sincronizado_somente_leitura` | 409 | Campo de equipamento sincronizado que vem do cliente |
| `cliente_servico_inativo` | 422 | Operação nova com cliente atendido inativo |

| Evento | Conta que recebe | Quando |
|---|---|---|
| `conexao_empresa.ativada` | As duas | Aceite ou reativação |
| `conexao_empresa.encerrada` | As duas | Recusa, encerramento ou convite vencido |
| `vinculo_os.sincronizado` | A que recebeu os dados | Alteração da OS vinculada gravada |
| `vinculo_os.erro` | A de origem da alteração | Sincronização esgotou as tentativas |

---

## 23. Configuração por segmento: regras, campos extras, fase de ciclo de vida e pacotes

A Tenaz atende segmentos diferentes **configurando blocos que o sistema já conhece** (sumário 3.7). A API reflete isso com catálogos fechados: condição, efeito, tipo de campo e tipo de item de pacote são enums **não extensíveis por configuração**; valor novo só entra por requisito no documento 1, com teste, e chega ao contrato como mudança compatível (18.1). Os estados canônicos, as fórmulas dos indicadores, as regras do app offline, o isolamento e a lista de grupos da conexão entre contas **não são configuráveis** por nenhuma rota (risco R7). Pacote Go: `configuracao` (regras e pacotes); os campos extras ficam no módulo OS e a fase de ciclo de vida, em ATV.

| Bloco | Rotas | Plano | Fase | Requisito |
|---|---|---|---|---|
| Regras configuráveis | `/v1/regras-automacao`, `/v1/disparos-regra` | E | MVP | RF-ADM-016 |
| Campos extras da OS | `/v1/campos-extras-os`, `/v1/ordens-servico/{id}/campos-extras` | E | MVP | RF-OS-021 |
| Fase de ciclo de vida | `/v1/ativos/{id}/colocar-em-operacao`, `/v1/ativos/colocar-em-operacao` | E | MVP | RF-ATV-013 |
| Pacotes de segmento | `/v1/pacotes-segmento`, `/v1/aplicacoes-pacote` | E | F2 | RF-ADM-017 |

### 23.1 Regras configuráveis (`regra_automacao`)

Uma regra tem **uma condição e um efeito**, escolhidos dos catálogos fechados do sumário 3.7, com parâmetros e um filtro opcional por classe de ativo, unidade e criticidade. Toda regra nasce `inativa` (proposta), para que a empresa simule antes de ativar; as que vêm de pacote também (RF-ADM-017). Criar, alterar, ativar, inativar e excluir exigem `adm.gerenciar_regras` (grupo Regras, 4.8).

**Catálogo de condições** (`condicao`):

| Valor | Parâmetros (`parametros_condicao`) | Quando é avaliada | Fase |
|---|---|---|---|
| `documento_vencido` | `tipo_documento_id` | No vencimento do documento do ativo e na rotina diária | MVP |
| `documento_a_vencer` | `tipo_documento_id`, `dias_antecedencia` | Na rotina diária | MVP |
| `os_concluida_com_campo` | `campo_extra_os_id` (campo do tipo `sim_nao`), `valor` (padrão `true`) | Na ação `concluir` da OS, inclusive pelo link e pelo app | MVP |
| `item_checklist_reprovado` | `checklist_modelo_id`, `item_id` (proposta) | No envio da execução do checklist | MVP |
| `leitura_fora_faixa` | `tipo_medidor` (opcional); a faixa é a do próprio medidor | No registro da leitura | MVP |
| `capacitacao_vencida` | `tipo_capacitacao` (proposta) | No vencimento da capacitação e na rotina diária | F2 (P) |
| `etapa_projeto_atrasada` | `dias_atraso` (padrão 0), `modelo_projeto_id` (opcional) | Na rotina diária | F2 (P, PRJ) |

**Catálogo de efeitos** (`efeito`):

| Valor | Parâmetros (`parametros_efeito`) | O que a API grava | Fase |
|---|---|---|---|
| `notificar` | `destinatarios` (perfis, `usuario_ids`, `responsavel_ativo`), `canais` (`app`, `email`, `push`) | Notificações (`/v1/notificacoes`) | MVP |
| `abrir_solicitacao` | `prioridade_sugerida`, `titulo` (com marcadores do fato, proposta) | `solicitacao` com `origem = regra` e `disparo_regra_id`, que segue a triagem normal (7.2) | MVP |
| `bloquear_ativo` | `motivo` | `bloqueio_ativo` com `origem = regra`: o ativo fica `parado` até a liberação manual com motivo (6.13) | MVP |
| `bloquear_ativo_ate_documento` | `motivo`, `tipo_documento_id` que levanta o bloqueio | O mesmo bloqueio, levantado sozinho quando um documento desse tipo é registrado no ativo com emissão posterior ao fato que disparou a regra (sumário 3.6) | MVP |
| `impedir_atribuicao` | — | Recusa a atribuição de OS a pessoa na condição (`422 capacitacao_vencida`) | F2 (P) |

**Combinações** (proposta, a confirmar no documento 4): `bloquear_ativo_ate_documento` exige `tipo_documento_id`; `abrir_solicitacao` não é aceito com `documento_a_vencer`, para não abrir solicitação repetida a cada dia; `impedir_atribuicao` só com `capacitacao_vencida`; `capacitacao_vencida` só com `notificar` ou `impedir_atribuicao`; `etapa_projeto_atrasada` só com `notificar`. Fora disso: `422 regra_combinacao_invalida`. Condição ou efeito da Fase 2 sem o módulo ou o plano: `403 plano_nao_inclui_recurso`.

**Criar** — `POST /v1/regras-automacao` (grupo Regras), com `Idempotency-Key`:

```json
{
  "nome": "Verificação metrológica vencida bloqueia o radar",
  "descricao": "O radar não volta a operar sem nova verificação do Inmetro.",
  "condicao": "documento_vencido",
  "parametros_condicao": { "tipo_documento_id": "0194a100-0000-7aaa-8bbb-000000000a01" },
  "filtro": { "classe_ativo_ids": ["0194a000-0000-7aaa-8bbb-000000000c01"], "unidade_ids": [], "criticidades": [] },
  "efeito": "bloquear_ativo_ate_documento",
  "parametros_efeito": {
    "tipo_documento_id": "0194a100-0000-7aaa-8bbb-000000000a01",
    "motivo": "Verificação metrológica vencida"
  }
}
```

**Resposta `201 Created`** (`Location: /v1/regras-automacao/0194a200-...`, `ETag: "1"`):

```json
{
  "id": "0194a200-0000-7aaa-8bbb-000000000r01",
  "versao": 1,
  "nome": "Verificação metrológica vencida bloqueia o radar",
  "situacao": "inativa",
  "condicao": "documento_vencido",
  "efeito": "bloquear_ativo_ate_documento",
  "parametros_condicao": { "tipo_documento_id": "0194a100-0000-7aaa-8bbb-000000000a01" },
  "parametros_efeito": { "tipo_documento_id": "0194a100-0000-7aaa-8bbb-000000000a01", "motivo": "Verificação metrológica vencida" },
  "filtro": { "classe_ativo_ids": ["0194a000-0000-7aaa-8bbb-000000000c01"], "unidade_ids": [], "criticidades": [] },
  "origem_pacote": null,
  "total_disparos": 0,
  "bloqueios_vigentes": 0,
  "transicoes_permitidas": ["ativar", "simular"],
  "criado_em": "2026-11-12T16:00:00Z",
  "atualizado_em": "2026-11-12T16:00:00Z",
  "criado_por_id": "0192f800-aaaa-7bbb-8ccc-000000000001"
}
```

`condicao` e `efeito` não mudam depois de criados (`422 campo_somente_leitura`): para trocar, cria-se outra regra. `origem_pacote` (proposta) traz o código do pacote quando a regra veio de uma aplicação (23.4).

**Simular** — `POST /v1/regras-automacao/{id}/simular` (Gestão), sem corpo ou com `{ "limite": 200 }`. Avalia a regra agora, ativa ou inativa, **sem gravar nada**: não cria disparo, solicitação, bloqueio nem notificação, e não exige `If-Match`. Para condições de evento (`os_concluida_com_campo`, `item_checklist_reprovado`, `leitura_fora_faixa`), considera os fatos dos últimos 30 dias (proposta). Resposta `200`:

```json
{
  "avaliado_em": "2026-11-12T16:05:00Z",
  "total_afetados": 2,
  "limite": 200,
  "ativos": [
    {
      "ativo_id": "0194a300-0000-7aaa-8bbb-000000000001",
      "tag": "RDR-0012",
      "fato": { "entidade": "documento_vencimento", "id": "0194a400-0000-7aaa-8bbb-000000000001", "descricao": "Verificação metrológica vencida em 02/11/2026." },
      "efeito_previsto": "bloquear_ativo_ate_documento",
      "ja_disparada": false
    }
  ]
}
```

**Ativar e inativar.** `ativar` respeita o limite de regras ativas por empresa (proposta: 50; `422 limite_regras_ativas_excedido`). `inativar` não levanta os bloqueios já criados pela regra (RF-ADM-016): a resposta traz `bloqueios_vigentes` com o `id` e a TAG de cada um, para a gestão decidir em `/v1/bloqueios-ativo/{id}/liberar`. Alteração por `PATCH` vale para os fatos seguintes; não refaz disparos passados.

**Avaliação e disparos.** A regra é avaliada no fato que a dispara e numa rotina diária idempotente: o par regra e fato dispara no máximo uma vez. Cada disparo gera um `disparo_regra`, visível em `GET /v1/regras-automacao/{id}/disparos` e `GET /v1/disparos-regra` (por ativo, por regra, por período), com `regra_automacao_id`, `condicao`, `efeito`, `ativo_id`, `fato` (`entidade` e `id`), `ocorrido_em`, `resultado` (`aplicado` ou `sem_efeito`, quando o efeito já estava em vigor, proposta) e os registros criados (`solicitacao_id`, `bloqueio_ativo_id`, `notificacoes`). Nos eventos e na auditoria, o autor é `{ "tipo": "regra_automacao", "id": "<id da regra>" }`. A exportação da lista de regras e de disparos usa `POST /v1/exportacoes-lista`.

**Limites do que a regra faz** (sumário 3.7). Uma regra nunca muda o estado de uma OS, nunca altera fórmula de indicador e nunca grava dado em outra conta. Os comportamentos padrão (alertas de documentos, item reprovado de checklist gerando solicitação) continuam sem regra; a regra acrescenta efeitos, não desliga os padrões.

Fluxo do exemplo dos radares (sumário 3.7), com as duas regras ativas:

```
Documento vence ──→ regra "verificação vencida" ──→ bloqueio_ativo (origem regra) ──→ ativo parado
OS concluída com rompeu_lacre = true ──→ regra "rompeu lacre" ──→ bloqueio_ativo (origem regra) ──→ ativo parado
Nova verificação registrada no ativo (emissão posterior ao fato) ──→ bloqueio levantado ──→ status recalculado
```

O pacote `fiscalizacao_eletronica` entrega essas regras `inativa` até a pendência **P9** (RF-ADM-017).

### 23.2 Campos extras da OS (`campo_extra_os`)

**Definição** — `POST /v1/campos-extras-os` (grupo Configuração, permissão `adm.gerenciar_configuracao`):

```json
{
  "chave": "rompeu_lacre",
  "rotulo": "Rompeu lacre?",
  "tipo": "sim_nao",
  "tipos_os": ["corretiva_planejada", "corretiva_emergencial", "preventiva"],
  "classe_ativo_ids": ["0194a000-0000-7aaa-8bbb-000000000c01"],
  "obrigatorio_em": "conclusao",
  "preenchido_pelo_prestador": true,
  "ajuda": "Marque sim se algum lacre do equipamento foi retirado."
}
```

| Campo | Regra |
|---|---|
| `chave` | Obrigatória, única na empresa, minúsculas sem acento (`^[a-z][a-z0-9_]{1,49}$`). Só leitura depois de haver valor gravado |
| `rotulo` | Obrigatório, até 80 caracteres. A OS mostra o rótulo vigente na data do registro |
| `tipo` | `texto`, `numero`, `data`, `lista`, `sim_nao`. Só leitura depois de haver valor gravado (cria-se outro campo) |
| `opcoes` | Obrigatória em `lista`: `[{ "valor": "...", "rotulo": "..." }]`; valor removido continua válido nas OS antigas |
| `unidade` | Opcional em `numero` (`mm`, `kWh`, texto livre) |
| `tipos_os`, `classe_ativo_ids` | Onde o campo aparece; lista vazia = todos |
| `obrigatorio_em` | `nenhum` (padrão), `abertura` ou `conclusao` (proposta de nome para o "momento da obrigatoriedade" do RF-OS-021) |
| `preenchido_pelo_prestador` | Marcação "preenchido pelo prestador" do RF-OS-021 (proposta; o nome do campo também é proposta). `false` por padrão. Com `true`, o campo aparece e é preenchido nos escopos `prestador` e `link`; sem ela, esses escopos não veem o campo |
| `ajuda`, `ordem`, `situacao` | Texto de ajuda; ordem de exibição (`PUT /v1/campos-extras-os/ordem`); `ativo` ou `inativo` |

Limite: 30 campos ativos por combinação de tipo de OS e classe de ativo (proposta; `422 limite_campos_extras_excedido`). Inativar (`DELETE`) não apaga valores gravados. Criar, editar, ordenar e inativar definições exigem `adm.gerenciar_configuracao` (sumário 3.6).

**Valores na OS.** `ordem_servico.campos_extras` é um objeto pela `chave`: `{ "rompeu_lacre": true, "lacres_retirados": "A1234, A1235" }`. Formato por tipo: `texto` (string, até 2.000 caracteres), `numero` (string decimal, 2.5), `data` (`YYYY-MM-DD`), `lista` (o `valor` de uma opção), `sim_nao` (booleano). São gravados por `POST` e `PATCH /v1/ordens-servico`, pelo corpo de `concluir`, por `PATCH /v1/ordens-servico/{id}/campos-extras`, pelo app (12.4) e pelo link (19.4). `PATCH` segue o JSON Merge Patch: chave com `null` apaga o valor.

`GET /v1/ordens-servico/{id}/campos-extras` devolve as definições que se aplicam à OS, com o valor atual:

```json
{
  "dados": [
    { "chave": "rompeu_lacre", "rotulo": "Rompeu lacre?", "tipo": "sim_nao", "obrigatorio_em": "conclusao", "preenchido_pelo_prestador": true, "valor": true },
    { "chave": "lacres_instalados", "rotulo": "Lacres instalados", "tipo": "texto", "obrigatorio_em": "conclusao", "preenchido_pelo_prestador": true, "valor": null }
  ]
}
```

**Obrigatoriedade.** Campo `abertura` vazio recusa a criação da OS; campo `conclusao` vazio recusa `concluir`; os dois com `422 campo_extra_obrigatorio`:

```json
{
  "type": "https://api.tenaz.com.br/erros/campo_extra_obrigatorio",
  "title": "Campo extra obrigatório não preenchido",
  "status": 422,
  "detail": "Preencha \"Lacres instalados\" antes de concluir a OS-2026-005120.",
  "codigo": "campo_extra_obrigatorio",
  "campos_faltantes": ["campos_extras.lacres_instalados"],
  "id_rastreio": "01JD2K4M6P8R0T2V4X6Z8B0D2F"
}
```

- A obrigatoriedade vale pelas definições vigentes na data da ação (no app, no `criado_em_dispositivo` do comando, 12.4). Campo novo obrigatório não afeta OS já concluídas (RF-OS-021).
- Campo obrigatório sem `preenchido_pelo_prestador` não impede a conclusão feita pelo `prestador`, pelo link ou pela subcontratada conectada (8.2): a OS fica `concluida` com o aviso `campos_extras_pendentes`, e `encerrar` passa a exigir o campo (proposta; ver 26.1, A-30).
- Os campos extras não entram em fórmula de indicador, não trafegam na conexão entre contas no lançamento (22.5) e aparecem no relatório da OS em PDF e na impressão.
- Filtro e exportação: `GET /v1/ordens-servico?campos_extras.rompeu_lacre=true` (igualdade e lista; faixa em `numero` e `data`, proposta) e `POST /v1/exportacoes-lista`.

### 23.3 Fase de ciclo de vida e entrada em operação

O ativo nasce `em_operacao` (padrão) ou `em_implantacao` (`POST /v1/ativos`, importação ou projeto). Enquanto `em_implantacao`, tem status `inativo`, fica fora dos KPIs e da cobrança de veículo, e seus planos ficam suspensos (6.2). Recebe OS `implantacao` (8.1). A passagem para `em_operacao` é única, pela ação `colocar_em_operacao`.

**Colocar em operação** — `POST /v1/ativos/{id}/colocar-em-operacao` (grupo Entrada em operação), com `If-Match` e `Idempotency-Key`:

```json
{
  "data_entrada_operacao": "2026-11-20T09:00:00Z",
  "termo_entrega_anexo_id": "0194a500-0000-7aaa-8bbb-000000000f01",
  "observacao": "Aferição aprovada; radar liberado pela prefeitura."
}
```

| Validação | Erro |
|---|---|
| O ativo está `em_implantacao` | `422 ativo_nao_em_implantacao` |
| Não há OS `implantacao` do ativo fora de `encerrada` ou `cancelada` (proposta) | `422 os_implantacao_pendente`, com `ordens_servico_pendentes` |
| Campos obrigatórios da classe preenchidos (proposta) | `422 campos_obrigatorios_transicao`, com `campos_faltantes` |
| Documentos críticos válidos, quando a classe os exige (ex.: verificação metrológica) (proposta) | `422 documentos_exigidos_pendentes`, com `tipos_documento_faltantes` |
| `data_entrada_operacao` (padrão: agora) não futura nem anterior à criação do ativo | `422 validacao_falhou` |
| Permissão `atv.colocar_em_operacao` | `403 permissao_negada` |

Resposta `200` com o ativo: `fase_ciclo_vida = em_operacao`, `status_operacional` recalculado (`operando`, ou `parado` se houver bloqueio vigente), `data_entrada_operacao`, `colocado_em_operacao_por_id` e `planos_ativados` (ids dos planos ativados a partir da data). Eventos: `ativo.colocado_em_operacao` e `ativo.status_alterado`. Se o ativo é item de contrato com cliente incluído "a partir da entrada em operação", o item passa a contar na medição dessa data (25.3).

**Em lote** — `POST /v1/ativos/colocar-em-operacao`, com `Idempotency-Key`, até 100 ativos. Cada ativo leva a sua `versao` (no lugar do `If-Match`) e é processado na própria transação, como os comandos do app (12.4); a falha de um não desfaz os outros:

```json
{
  "data_entrada_operacao": "2026-11-20T09:00:00Z",
  "ativos": [
    { "ativo_id": "0194a300-0000-7aaa-8bbb-000000000001", "versao": 6 },
    { "ativo_id": "0194a300-0000-7aaa-8bbb-000000000002", "versao": 3 }
  ]
}
```

```json
{
  "resultados": [
    { "ativo_id": "0194a300-0000-7aaa-8bbb-000000000001", "tag": "RDR-0012", "status": "aplicado", "versao": 7 },
    {
      "ativo_id": "0194a300-0000-7aaa-8bbb-000000000002",
      "tag": "RDR-0013",
      "status": "rejeitado",
      "erro": { "codigo": "os_implantacao_pendente", "status": 422, "detail": "A OS-2026-005131 (montagem) ainda está em execução.", "ordens_servico_pendentes": ["OS-2026-005131"] }
    }
  ]
}
```

A fase **não volta** para `em_implantacao`: a reinstalação num outro ponto é transferência com OS `implantacao` ou projeto (RF-ATV-013). `alterar-status-operacional` não reativa ativo `em_implantacao` (`422 ativo_em_implantacao`).

### 23.4 Pacotes de segmento (Fase 2)

O catálogo de pacotes (`pacote_segmento`) é **global**, mantido pela Tenaz, fora do RLS por empresa e só leitura para as contas (RF-ADM-017). Pacotes iniciais (sumário 3.7, proposta): `fiscalizacao_eletronica`, `transportadora`, `industria`, `comercio`, `locadora`, `prestador_servicos`. O pacote é identificado pelo `codigo` na rota (2.1) e tem `versao` inteira.

`GET /v1/pacotes-segmento/fiscalizacao_eletronica` devolve o conteúdo agrupado por `tipo_item`: `classe_ativo`, `campo_extra_os`, `procedimento`, `plano_manutencao`, `checklist_modelo`, `tipo_documento`, `regra_automacao`, `modelo_projeto`, `codigo_falha`, mais `modulos_recomendados` (só informativo).

**Pré-visualizar** — `POST /v1/pacotes-segmento/{codigo}/pre-visualizar` (grupo Configuração), corpo `{ "versao": 2 }` (padrão: a mais recente). Não grava nada:

```json
{
  "pacote_codigo": "fiscalizacao_eletronica",
  "versao": 2,
  "resumo": { "a_criar": 41, "a_ignorar": 6 },
  "itens": [
    { "tipo_item": "classe_ativo", "chave_natural": "RDR", "nome": "Radar fixo", "acao": "criar" },
    { "tipo_item": "tipo_documento", "chave_natural": "verificacao_metrologica", "nome": "Verificação metrológica", "acao": "ignorar", "motivo": "ja_existe" },
    { "tipo_item": "modelo_projeto", "chave_natural": "Implantação de ponto de fiscalização", "acao": "ignorar", "motivo": "modulo_nao_contratado" }
  ]
}
```

**Aplicar** — `POST /v1/pacotes-segmento/{codigo}/aplicar` (grupo Configuração, permissão `adm.gerenciar_configuracao`, com `Idempotency-Key`), corpo `{ "versao": 2 }`. Resposta `201 Created` com `Location: /v1/aplicacoes-pacote/{id}` e a `aplicacao_pacote` (pacote, versão, data, quem aplicou, `itens_criados` e `itens_ignorados` com o motivo).

- A aplicação **só cria o que falta**, identificado pela chave natural do item (sigla da classe, código do tipo de documento, chave do campo extra, nome da regra, do procedimento, do plano, do checklist e do modelo de projeto). Nunca altera nem apaga o que a conta já tem ou editou.
- Tudo o que o pacote cria vira configuração comum da conta, editável como se a empresa tivesse criado. Regras nascem `inativa`.
- Itens de módulo ou plano não contratado não são criados: aparecem como ignorados (`modulo_nao_contratado`, `plano_nao_inclui`).
- A aplicação é **atômica** (RF-ADM-017): roda numa transação e cria tudo o que falta ou nada. Falha em qualquer item desfaz a aplicação inteira e devolve o erro do item, sem `aplicacao_pacote` gravada. Aplicar de novo a mesma versão não cria nada.
- Escolher pacotes no roteiro de primeiros passos (RF-ADM-001) usa as mesmas rotas, a partir da Fase 2; no MVP, o roteiro segue sem pacotes.
- Aplicar pacote é ação administrativa: não está disponível para tokens de API (4.4).

---

## 24. PRJ: projetos de implantação

Módulo da Fase 2, plano Profissional (RF-PRJ-001 a RF-PRJ-007; sumário 3.1 a 3.3). O projeto organiza OS `implantacao` em etapas com dependências; quem executa é a OS (seção 8), atribuída a técnico interno ou a fornecedor (19.3), inclusive com conta conectada (22.8). O projeto e as etapas **não** trafegam na conexão entre contas no lançamento: só a OS (RF-PRJ-004).

### 24.1 Estados e ações

**Projeto** (sumário 3.3): `rascunho` → `em_andamento` → `concluido`; `em_andamento` ⇄ `suspenso`, com motivo; `cancelado` a partir de `rascunho`, `em_andamento` ou `suspenso`, com motivo. `concluido` e `cancelado` não aceitam edição (`409 registro_imutavel`).

| Ação | Rota | De | Para | Campos obrigatórios | Efeitos |
|---|---|---|---|---|---|
| Iniciar | `POST /v1/projetos/{id}/iniciar` | `rascunho` | `em_andamento` | Pelo menos um ativo ou uma localização; `responsavel_id` (`422 campos_obrigatorios_transicao`) | Torna `disponivel` as etapas sem predecessoras; gera as OS das etapas marcadas para geração automática (proposta); `projeto.iniciado` |
| Suspender | `.../suspender` | `em_andamento` | `suspenso` | `motivo` | Pausa os prazos das etapas, não as OS |
| Retomar | `.../retomar` | `suspenso` | `em_andamento` | — | Os prazos voltam a correr |
| Concluir | `.../concluir` | `em_andamento` | `concluido` | Todas as etapas `concluida` ou `dispensada` (`422 projeto_etapas_pendentes`) | Coloca os ativos em operação (24.6); `projeto.concluido` |
| Cancelar | `.../cancelar` | `rascunho`, `em_andamento`, `suspenso` | `cancelado` | `motivo`; `destino_ativos` para cada ativo `em_implantacao`; nenhuma OS em `em_execucao` (`422 projeto_com_os_em_execucao`) | `projeto.cancelado` |

**Etapa de projeto** (sumário 3.3): `pendente` → `disponivel` (todas as predecessoras `concluida` ou `dispensada`) → `em_andamento` → `concluida`; `dispensada` a partir de `pendente` ou `disponivel`. A passagem a `disponivel` é automática, nunca por rota; o nome evita confusão com a liberação da OS (sumário 3.6). Ações de etapa exigem o projeto `em_andamento` (senão `409 transicao_invalida`, com o estado do projeto).

| Ação | Rota | De | Para | Regras |
|---|---|---|---|---|
| Gerar OS | `POST /v1/projetos/{id}/etapas/{etapa_id}/gerar-os` | `disponivel`, `em_andamento` | Mesmo estado | Etapa `pendente`: `409 etapa_predecessora_pendente` (24.5) |
| Iniciar | `.../iniciar` | `disponivel` | `em_andamento` | Só em etapa sem OS; com OS, a etapa inicia sozinha quando a primeira OS dela inicia |
| Concluir | `.../concluir` | `em_andamento`; `disponivel` em etapa sem OS (proposta) | `concluida` | OS da etapa `encerrada` ou `cancelada`, ao menos uma `encerrada` se a etapa gera OS (`422 etapa_os_pendente`); documentos exigidos (`422 documentos_exigidos_pendentes`); campos exigidos (`422 campo_extra_obrigatorio`). Torna `disponivel` as sucessoras cujas predecessoras terminaram e notifica seus responsáveis |
| Dispensar | `.../dispensar` | `pendente`, `disponivel` | `dispensada` | `motivo`; permissão `prj.dispensar_etapa`; a etapa de entrega não é dispensável (`422 etapa_entrega_nao_dispensavel`, proposta). Torna `disponivel` as sucessoras |

Etapa com data prevista vencida e não concluída vem com `atrasada = true` e `dias_atraso` (condição calculada, não estado) e pode disparar regra com a condição `etapa_projeto_atrasada` (23.1).

### 24.2 Modelo de projeto

`POST /v1/modelos-projeto` (grupo Projetos, permissão `prj.gerenciar`). As etapas do modelo (`modelo_etapa_projeto`) vêm no corpo; as predecessoras se referem a outras etapas pela `chave` local:

```json
{
  "nome": "Implantação de ponto de fiscalização",
  "descricao": "Autorização, energia, fundação e poste, laços, montagem, aferição e entrega.",
  "classe_ativo_ids": ["0194a000-0000-7aaa-8bbb-000000000c01"],
  "etapas": [
    { "chave": "autorizacao", "nome": "Autorização do órgão de trânsito", "ordem": 1, "predecessoras": [], "duracao_prevista_dias": 20, "peso": 5, "executor_padrao": "interno", "gera_os": false, "tipos_documento_exigidos_ids": ["0194a100-0000-7aaa-8bbb-000000000a07"], "campos_exigidos": ["numero_autorizacao"], "valor_referencia": "0.00", "etapa_entrega": false },
    { "chave": "energia", "nome": "Ligação de energia", "ordem": 2, "predecessoras": ["autorizacao"], "duracao_prevista_dias": 30, "peso": 15, "executor_padrao": "fornecedor", "gera_os": true, "procedimento_id": "0194a600-0000-7aaa-8bbb-000000000p01", "campos_exigidos": ["protocolo_concessionaria"], "valor_referencia": "1800.00", "etapa_entrega": false },
    { "chave": "entrega", "nome": "Entrega e termo", "ordem": 7, "predecessoras": ["afericao"], "duracao_prevista_dias": 2, "peso": 5, "executor_padrao": "interno", "gera_os": false, "valor_referencia": "0.00", "etapa_entrega": true }
  ]
}
```

- Validações (RF-PRJ-001): pelo menos uma etapa; predecessoras do mesmo modelo e sem ciclo (`422 dependencia_ciclica`); `peso` maior ou igual a zero, com soma maior que zero; no máximo 40 etapas (proposta; `422 validacao_falhou` com `tamanho_maximo`); no máximo uma `etapa_entrega`.
- `campos_exigidos` são chaves de `campo_extra_os` que precisam estar preenchidas nas OS da etapa; em etapa sem OS, ficam na própria etapa (24.4).
- O modelo tem `numero_versao` (proposta), que sobe quando as etapas mudam. O projeto guarda `modelo_projeto_versao` e não muda quando o modelo muda.

### 24.3 Criar projeto

`POST /v1/projetos` (grupo Projetos), com `Idempotency-Key`:

```json
{
  "modelo_projeto_id": "0194a700-0000-7aaa-8bbb-000000000m01",
  "nome": "Ponto Av. Brasil x Rua 7 — sentido centro",
  "unidade_id": "0194a800-0000-7aaa-8bbb-000000000u01",
  "localizacao_id": "0194a810-0000-7aaa-8bbb-000000000l01",
  "ativo_ids": ["0194a300-0000-7aaa-8bbb-000000000001"],
  "cliente_servico_id": "0194a900-0000-7aaa-8bbb-000000000s01",
  "contrato_cliente_id": "0194aa00-0000-7aaa-8bbb-000000000t01",
  "responsavel_id": "0192f800-aaaa-7bbb-8ccc-000000000017",
  "data_inicio_prevista": "2026-12-01",
  "orcamento_etapas": [ { "chave": "energia", "valor": "2100.00" } ]
}
```

Resposta `201` com `numero` (`PJ-2026-000031`), `status = rascunho`, `modelo_projeto_versao`, as etapas copiadas com datas previstas calculadas pela duração e pelas dependências, `avanco_fisico_percentual = 0`, `custo_previsto` e `custo_realizado`. Ativo já em outro projeto não concluído nem cancelado: `409 ativo_em_outro_projeto` (proposta). Em `rascunho`, ativos e localização mudam por `PATCH` (`ativo_ids`). Vários projetos de uma vez: importação `projetos` (14.3).

### 24.4 Etapas

`GET /v1/projetos/{id}/etapas` devolve, por etapa: `id`, `chave`, `nome`, `status`, `predecessoras`, `executor_tipo` e `fornecedor_id`, `responsavel_id`, datas previstas e reais, `valor`, `atrasada`, `dias_atraso`, `ordens_servico` (id, número e status), `motivo_dispensa` e `exigencias`:

```json
{
  "exigencias": {
    "documentos": [ { "tipo_documento_id": "0194a100-0000-7aaa-8bbb-000000000a07", "nome": "Autorização do órgão", "situacao": "atendido", "documento_vencimento_id": "0194ab00-0000-7aaa-8bbb-000000000d01" } ],
    "campos": [ { "chave": "numero_autorizacao", "situacao": "pendente" } ]
  }
}
```

- Documento exigido é registrado no ativo por `/v1/documentos-vencimento` e vale para a etapa quando emitido depois do início dela (proposta).
- Em etapa sem OS, os campos exigidos são gravados por `PATCH /v1/projetos/{id}/etapas/{etapa_id}` no objeto `campos_extras` da etapa (proposta).
- `PATCH` da etapa altera executor, responsável, datas previstas e `valor`; não altera `status` (`422 campo_somente_leitura`).

### 24.5 OS de implantação da etapa

`POST /v1/projetos/{id}/etapas/{etapa_id}/gerar-os` (grupo Projetos), com `Idempotency-Key`:

```json
{ "por_ativo": true, "data_limite": "2026-12-20T21:00:00Z", "responsavel_id": null }
```

- Cria uma OS, ou uma por ativo do projeto com `por_ativo = true` (proposta), com `tipo = implantacao`, `origem = projeto`, `etapa_projeto_id`, o ativo ou a localização do projeto e as tarefas do procedimento da etapa. Resposta `201` com `ordens_servico` (id e número de cada uma).
- Com executor `fornecedor`, a OS nasce atribuída ao fornecedor da etapa, como em `atribuir-fornecedor` (19.3); com fornecedor de conexão `ativa`, a OS vai para a conta dele (22.8), sem o projeto.
- Geração automática quando a etapa fica `disponivel`, quando o modelo marca `gera_os`, é proposta do documento 1 (RF-PRJ-004) e fica ligada por parâmetro do projeto.
- As validações são as da OS (8.1), inclusive os campos extras obrigatórios na abertura. A primeira OS da etapa a entrar em `em_execucao` leva a etapa a `em_andamento`. O custo da OS entra no custo realizado do projeto, como custo de implantação.

### 24.6 Concluir o projeto e colocar os ativos em operação

`POST /v1/projetos/{id}/concluir` (permissões `prj.gerenciar` e `atv.colocar_em_operacao`), com `If-Match`. Corpo opcional: `{ "data_entrada_operacao": "2026-12-22T10:00:00Z", "termo_entrega_anexo_id": "..." }`.

- Exige todas as etapas `concluida` ou `dispensada` (`422 projeto_etapas_pendentes`, com `etapas_pendentes`).
- Executa `colocar_em_operacao` (23.3) em todos os ativos do projeto ainda `em_implantacao`, numa transação: se algum não puder, nada é concluído e a resposta é `422 ativos_nao_aptos_operacao`, com o impedimento de cada um (proposta):

```json
{
  "type": "https://api.tenaz.com.br/erros/ativos_nao_aptos_operacao",
  "title": "Há ativos que não podem entrar em operação",
  "status": 422,
  "detail": "1 ativo do PJ-2026-000031 não pode entrar em operação.",
  "codigo": "ativos_nao_aptos_operacao",
  "ativos": [ { "ativo_id": "0194a300-0000-7aaa-8bbb-000000000001", "tag": "RDR-0012", "codigo": "documentos_exigidos_pendentes", "tipos_documento_faltantes": ["Verificação metrológica"] } ],
  "id_rastreio": "01JD2M7Q9S1U3W5Y7A9C1E3G5J"
}
```

- Com sucesso: um `ativo.colocado_em_operacao` por ativo e `projeto.concluido`. Com `contrato_cliente_id`, os ativos entram como itens cobertos a partir da `data_entrada_operacao` (RF-PRJ-007, proposta; 25.3).
- O termo de entrega sai em `GET /v1/projetos/{id}/termo-entrega` (PDF) e recebe as assinaturas em `POST /v1/projetos/{id}/assinaturas` (`papel`: `responsavel` ou `representante_cliente`; `assinatura_anexo_id`), com a mesma assinatura eletrônica avançada da OS (P1).

**Cancelar** — `POST /v1/projetos/{id}/cancelar` com `motivo` e `destino_ativos: [ { "ativo_id": "...", "destino": "manter_em_implantacao" } ]` ou `"inativar"` com `motivo_inativacao`. O cancelamento não apaga ativos e não cancela as OS ainda não iniciadas: a resposta lista `ordens_servico_abertas` para decisão do planejador (proposta).

### 24.7 Medição de etapas para pagamento do prestador

Relatório persistido por fornecedor e período (RF-PRJ-006), base do pagamento da empreitada. O nome da entidade, `medicao_etapa`, é **proposta** (o sumário não a nomeia; ver 26.1, A-34). Situações: `em_apuracao` → `fechada` (proposta, como a medição do contrato).

- `POST /v1/medicoes-etapa` com `{ "fornecedor_id": "...", "periodo_inicio": "2026-12-01", "periodo_fim": "2026-12-31" }` cria a medição com as etapas `concluida` no período, executadas pelo fornecedor e ainda não medidas. Cada linha traz projeto, etapa, ativo, data de conclusão, `valor` (padrão: o valor da etapa, que vem do modelo ou do contrato de prestação, RF-TER-006), `valor_ajustado` e `motivo_ajuste`.
- Etapa `dispensada` não entra; uma etapa entra numa única medição (`409 etapa_ja_medida` quando a requisição lista explicitamente uma etapa já medida, proposta).
- `POST /v1/medicoes-etapa/{id}/ajustar` com `{ "etapa_projeto_id": "...", "valor": "1500.00", "motivo": "Poste reaproveitado; fundação não executada." }`, exige `prj.ajustar_medicao`.
- `fechar` congela a medição (`409 registro_imutavel` depois); `DELETE` só em `em_apuracao`. Relatório em PDF ou XLSX em `/relatorio`.
- A Tenaz não paga nem fatura: o pagamento sai do ERP.

---

## 25. CTR: contratos com clientes

Módulo da Fase 2, plano Profissional (RF-CTR-001 a RF-CTR-005). Controla o contrato de serviço com um `cliente_servico` que exige disponibilidade dos equipamentos (por exemplo, o contrato de fiscalização eletrônica com uma prefeitura) e apura a **disponibilidade contratual** mês a mês, com memória de cálculo. **Não muda** MTTR, MTBF nem a disponibilidade técnica dos indicadores (seção 11). O cliente final não tem acesso à API nem ao painel no lançamento (P10): a medição sai em PDF e XLSX.

### 25.1 Estados e ações

**Contrato com cliente** (sumário 3.3): `rascunho` → `vigente` → `encerrado`, ou `cancelado` só a partir de `rascunho`, como o contrato de locação. `encerrado` e `cancelado` não aceitam edição (`409 registro_imutavel`).

| Ação | Rota | De | Para | Campos obrigatórios e validações | Efeitos |
|---|---|---|---|---|---|
| Ativar | `POST /v1/contratos-cliente/{id}/ativar` | `rascunho` | `vigente` | Pelo menos um item (`422 contrato_cliente_sem_itens`), meta, calendário e faixas válidas (`422 faixas_desconto_invalidas`) | `contrato_cliente.ativado`; os prazos passam a valer nas OS dos itens |
| Encerrar | `.../encerrar` | `vigente` | `encerrado` | `motivo_encerramento` quando antes de `data_fim` | `contrato_cliente.encerrado`; a medição do último período continua em apuração |
| Cancelar | `.../cancelar` | `rascunho` | `cancelado` | `motivo` | — |

**Medição do contrato** (sumário 3.3): `em_apuracao` → `fechada`; `reabrir` volta a `em_apuracao`, com motivo e a permissão `ctr.reabrir_medicao`.

### 25.2 Criar contrato

`POST /v1/contratos-cliente` (grupo Contratos, permissão `ctr.gerenciar`), com `Idempotency-Key`:

```json
{
  "cliente_servico_id": "0194a900-0000-7aaa-8bbb-000000000s01",
  "numero_contrato_cliente": "Pregão eletrônico 041/2026",
  "data_inicio": "2027-01-01",
  "data_fim": "2028-12-31",
  "renovacao_automatica": false,
  "unidade_ids": ["0194a800-0000-7aaa-8bbb-000000000u01"],
  "meta_disponibilidade": 98.00,
  "calendario_apuracao": { "tipo": "continuo_24x7" },
  "prazos_prioridade": [
    { "prioridade": "urgente", "atendimento_horas": 4, "solucao_horas": 24 },
    { "prioridade": "alta", "atendimento_horas": 8, "solucao_horas": 48 }
  ],
  "faixas_desconto": [
    { "de_percentual": 95.00, "ate_percentual": 98.00, "desconto_percentual": 5.00 },
    { "de_percentual": 90.00, "ate_percentual": 95.00, "desconto_percentual": 10.00 },
    { "de_percentual": 0.00, "ate_percentual": 90.00, "desconto_percentual": 20.00 }
  ],
  "causas_externas_aceitas": ["vandalismo", "acidente_terceiro", "falta_energia"],
  "evidencia_causa_externa_obrigatoria": true,
  "valor_mensal_referencia": "48000.00"
}
```

- `calendario_apuracao.tipo`: `continuo_24x7` ou `calendario_contrato` (proposta de nome para o "calendário do contrato" do RF-CTR-001), este com `dias_semana`, `hora_inicio` e `hora_fim`.
- Percentuais são número decimal de 0 a 100 com até duas casas (2.5). As faixas são intervalos `[de, ate)`, contíguas da meta para baixo; disponibilidade igual ou acima da meta não tem desconto.
- `causas_externas_aceitas` é subconjunto da lista canônica de `causa_externa` (sumário 3.6).
- Resposta `201` com `numero` (`CT-2026-000004`) e `status = rascunho`.
- Em `vigente`, `PATCH` só altera `data_fim`, `renovacao_automatica`, contatos e anexos; mudança de meta, faixas ou prazos é aditivo com histórico, proposta do documento 1 ainda sem rota (ver 26.1, A-35).

### 25.3 Itens cobertos

`POST /v1/contratos-cliente/{id}/itens` com `{ "ativo_id": "...", "data_inclusao": "2027-01-01", "meta_disponibilidade": null, "valor_mensal": "1600.00" }`.

- O ativo pode ser `proprio`, `de_terceiro` ou `de_cliente`, e não pode estar `inativo` na inclusão, salvo `em_implantacao`, que só passa a contar a partir da `data_entrada_operacao` (RF-CTR-002).
- Um ativo fica em no máximo um contrato com cliente vigente no mesmo período (`409 ativo_em_outro_contrato_cliente`).
- Em `vigente`, a inclusão é aditivo; a exclusão grava `data_exclusao` por `PATCH`. O item conta na medição entre a inclusão (ou a entrada em operação, o que vier depois) e a exclusão.

### 25.4 Causa externa da parada

`ordem_servico.causa_externa` (8.1) é marcada na abertura ou na execução, pela web, pelo app, pelo `prestador` e pelo link de uso único (19.4), com as evidências em `anexos_ids` da OS. Campos só leitura acrescentados (proposta): `causa_externa_marcada_por_id` e `causa_externa_marcada_em`. A causa **não muda** MTTR, MTBF nem a disponibilidade técnica; só a medição do contrato desconta a parada, e só com `decisao_exclusao = aprovada` (25.5). A OS mostra `decisao_exclusao` só leitura, com `pendente` enquanto não houver decisão (onde o valor é gravado fica para o documento 2). Na conexão entre contas, a causa externa nunca é compartilhada no lançamento (sumário 3.6): a subcontratada conectada não a marca na OS do cliente (ver 26.1, A-32).

### 25.5 Medição do contrato

`POST /v1/medicoes-contrato` com `{ "contrato_cliente_id": "...", "periodo_inicio": "2027-01-01", "periodo_fim": "2027-01-31" }` abre a medição (`409 medicao_periodo_duplicado` se já houver). A rotina mensal abre sozinha a do mês anterior (proposta). Enquanto `em_apuracao`, os números são recalculados a cada consulta.

`GET /v1/medicoes-contrato/{id}`, trecho:

```json
{
  "id": "0194ac00-0000-7aaa-8bbb-000000000x01",
  "contrato_cliente_id": "0194aa00-0000-7aaa-8bbb-000000000t01",
  "numero_contrato": "CT-2026-000004",
  "periodo_inicio": "2027-01-01",
  "periodo_fim": "2027-01-31",
  "situacao": "em_apuracao",
  "itens": [
    {
      "ativo_id": "0194a300-0000-7aaa-8bbb-000000000001",
      "tag": "RDR-0012",
      "horas_requeridas": "744.00",
      "horas_paradas": "30.50",
      "horas_excluidas": "18.00",
      "disponibilidade_contratual": 98.32,
      "meta": 98.00,
      "faixa": null,
      "desconto_percentual": 0.00,
      "desconto_valor": "0.00",
      "paradas": [
        { "ordem_servico_id": "0194ad00-0000-7aaa-8bbb-000000000o01", "numero_os": "OS-2027-000112", "horas": "18.00", "causa_externa": "vandalismo", "decisao_exclusao": "aprovada", "decidido_por_id": "0192f800-aaaa-7bbb-8ccc-000000000001" },
        { "ordem_servico_id": "0194ad00-0000-7aaa-8bbb-000000000o02", "numero_os": "OS-2027-000140", "horas": "12.50", "causa_externa": "falta_energia", "decisao_exclusao": "pendente", "decidido_por_id": null }
      ]
    }
  ],
  "total": { "valor_referencia": "48000.00", "desconto_valor": "0.00", "valor_com_desconto": "48000.00" },
  "memoria_calculo": "Disponibilidade contratual = (horas requeridas − horas paradas não excluídas) ÷ horas requeridas, no calendário de apuração."
}
```

- Disponibilidade contratual do item = (horas requeridas − horas paradas não excluídas) ÷ horas requeridas, dentro do calendário de apuração; o tratamento de paradas que cruzam o período fica no documento 4 (RF-CTR-004).
- Cada parada com causa externa tem `decisao_exclusao` (sumário 3.6): `pendente` (inicial), `aprovada` ou `recusada`; só a `aprovada` é descontada. Parada sem causa externa vem com `decisao_exclusao = null`. `POST /v1/medicoes-contrato/{id}/decidir-exclusao` com `{ "ordem_servico_id": "...", "decisao_exclusao": "aprovada", "motivo": "Boletim de ocorrência anexado." }` (`aprovada` ou `recusada`), com a permissão `ctr.fechar_medicao`. Causa fora das aceitas pelo contrato: `422 causa_externa_nao_aceita`; sem evidência, quando o contrato exige: `422 validacao_falhou` (proposta).
- `fechar` exige a permissão `ctr.fechar_medicao` e nenhuma parada com `decisao_exclusao = pendente` (`422 medicao_exclusoes_pendentes`) e gera `medicao_contrato.fechada`. Medição `fechada` não é editada (`409 registro_imutavel`); mudança em OS depois do fechamento não altera a medição.
- `reabrir` exige `motivo` e `ctr.reabrir_medicao`; a resposta traz `divergencias` (o que mudou nas OS desde o fechamento, proposta).
- `GET /v1/medicoes-contrato/{id}/relatorio` gera o PDF ou o XLSX com a memória de cálculo. A Tenaz não fatura: o valor com desconto vai para o ERP.

**Prazos de atendimento e de solução.** Nas OS dos itens cobertos, o contrato vigente mede o prazo de atendimento (`data_reconhecimento` − `data_abertura`) e de solução (`data_conclusao` − `data_abertura`) pela prioridade, no calendário de apuração (proposta). O estouro gera `contrato_cliente.sla_violado` e notifica o `gestor` (RF-CTR-001).

### 25.6 Painel do contrato

`GET /v1/contratos-cliente/{id}/painel` (Gestão) usa o mesmo cálculo da medição, com o período até o momento (RF-CTR-005): disponibilidade parcial por item e do contrato, itens abaixo da meta ordenados pelo impacto no desconto, desconto projetado, OS abertas dos itens com os prazos e as paradas com causa externa aguardando decisão. O painel é interno até a decisão de **P10**.

---

## 26. Pendências e decisões em aberto

### 26.1 Inconsistências e lacunas encontradas no sumário executivo

Os 15 pontos da versão 0.1 (I-1 a I-15) foram resolvidos pelo sumário 0.3 e já estavam aplicados na versão 0.2 deste documento; a única parte que seguia aberta, a lista de entidades sem nome canônico de I-14, continua em A-2. A tabela abaixo mostra a situação de cada ponto registrado na versão 0.2, dos que surgiram no alinhamento com o sumário 0.4 e dos que surgiram no alinhamento com o sumário 0.5 (A-29 a A-36).

**Resolvidos pelo sumário 0.5 (seção 3.6) e aplicados na versão 0.4 deste documento:**

| # | Ponto | Como o sumário 0.5 resolveu, e o que mudou aqui |
|---|---|---|
| A-21 | `origem` da cópia da OS na conta da subcontratada | Valor `cliente_conectado` (sumário 3.6). A cópia passa a nascer com `origem = cliente_conectado` (8.1, 22.8, 13.2); a proposta anterior, de repetir a origem do cliente, foi abandonada |
| A-29 | Origem `regra` da OS sem efeito que a produza | Não existe origem `regra` na OS: nenhum efeito de regra abre OS; a regra abre solicitação com origem `regra` (7.2). A OS ganha só a origem `projeto` (8.1) |
| A-32 | `causa_externa` e conexão entre contas | Nunca compartilhada no lançamento, junto com campos extras, etapa de projeto e fase de ciclo de vida (22.5). Consequência: a subcontratada com conta própria não marca a causa externa na OS do cliente; quem marca é o cliente, o `prestador` ou o link (25.4). Se o produto quiser a marcação pela subcontratada conectada, é preciso um grupo novo no sumário |
| — | Estado da etapa de projeto | `disponivel`, e não "liberada" (24.1) |
| — | Decisão sobre a parada com causa externa | Campo `decisao_exclusao`: `pendente`, `aprovada`, `recusada` (25.5) |
| — | Permissões de configuração, PRJ e CTR | Nove permissões, com os perfis padrão do sumário (4.8): `adm.gerenciar_regras`, `adm.gerenciar_configuracao`, `atv.colocar_em_operacao`, `prj.gerenciar`, `prj.dispensar_etapa`, `prj.ajustar_medicao`, `ctr.gerenciar`, `ctr.fechar_medicao`, `ctr.reabrir_medicao` |
| — | Aplicação de pacote | Atômica: tudo o que falta ou nada (23.4) |
| — | Campo extra visível ao prestador | Marcação "preenchido pelo prestador" (proposta), que define o que os escopos `prestador` e `link` veem e preenchem; aqui, o campo `preenchido_pelo_prestador` (proposta de nome, 23.2) |

**Resolvidos pelo sumário 0.4 (seção 3.6) e já aplicados aqui:**

| # | Ponto | Como o sumário 0.4 resolveu, e o que mudou aqui |
|---|---|---|
| A-1 | `bloqueio_despacho` × `bloqueio_ativo` | `bloqueio_ativo` é a entidade que deixa o ativo `parado`; "apto para despacho" é verificação calculada; `bloqueio_energia` é o registro LOTO. O cadastro manual passou de `/v1/bloqueios-despacho` para `/v1/bloqueios-ativo` (6.13); o motivo `bloqueio_manual` virou `bloqueio_ativo` (16.1); os eventos `bloqueio_despacho.criado` e `.liberado` viraram `ativo.despacho_bloqueado` e `ativo.despacho_liberado` (13.2). Resta só o nome da rota de consulta (A-3) |
| A-4 | Nome da permissão de liberação | `os.encerrar`, no formato `<módulo>.<ação>`; substitui `ordem_servico.liberar` em todo o documento (4.8, 8.2). Nova permissão `est.decidir_destino` para o destino de peça com defeito (20.3) |
| A-5 | Ação "assumir" | `POST /v1/ordens-servico/{id}/assumir`: técnico se atribui OS `aberta` ou `planejada` da sua equipe; grava `data_reconhecimento` (8.2) |
| A-6 | Campos de LOC e OS não fixados | Fixados: `local_uso`, `franquia_horas`, `franquia_km`, categorias da matriz (`preventiva`, `corretiva_desgaste`, `avaria_mau_uso`, `pneus`, `deslocamento`), `justificativa_responsavel_custo`, `motivo_recusa_liberacao`, `origem` da OS com 8 valores (8.1; 10 valores desde a versão 0.4, com `cliente_conectado` e `projeto`). O objeto `franquia` dos exemplos virou os campos planos (21.2). Continuam propostas: ver A-10 |
| A-9 | Locatário no app de campo | Só portal web no lançamento do LOC (4.9, 12.2) |
| A-11 | Link de uso único com `{token}` no caminho | Token no fragmento (`/l/#<token>`), ativação por `POST /v1/publico/links-execucao/ativar` com o token no corpo e sessão curta; demais rotas em `/atual/*` (19.4). Nenhuma rota recebe token no caminho (2.1) |
| A-12 | Escopo `fornecedor` como valor de escopo | Valores `interno`, `prestador`, `locatario`, `link`, expostos em `escopo_acesso` (4.9) |
| A-13 | Desistência de solicitação e `duplicada` | `nova` → `recusada` direto, com motivo `desistencia`; `duplicada` é final, corrigida por `reabrir-como-nova` (7.1); origem `portal_locatario` (7.2) |
| A-14 | Locação: cancelamento, sublocação, ativo devolvido | Só `rascunho` cancela; `vigente` termina por `encerrar` (com motivo se antecipado); sublocação recusada com `422 sublocacao_nao_permitida`; ativo `de_terceiro` devolvido vai a `inativo` (21.1) |
| A-15 | Papéis de fornecedor | Acrescentados `posto_combustivel` e `seguradora` (19.1) |
| A-16 | Peça em trânsito | Não consumível antes da confirmação de recebimento: `422 estoque_em_transito` (20.2) |

**Abertos:**

| # | Ponto | Tratamento neste documento |
|---|---|---|
| A-2 | Entidades da API ainda sem nome no sumário: sessão, sessão do link de uso único, dispositivo, comentário, assinatura de OS, documento técnico, item de material da OS, saldo de estoque, composição de veículo, conexão de telemática, conector de ERP, importação, exportação, processamento, painel, notificação, inventário de estoque, consulta de despacho | Nomes desta especificação são propostas; fixar no documento 2 e acrescentar ao vocabulário |
| A-3 | Exceções à regra de rota: `/v1/webhooks` (`assinatura_webhook`), `/v1/sincronizacao/comandos` (`comando_sync`), `/v1/ia/rascunhos-*` (`rascunho_ia`), `/v1/ativos/{id}/enquadramento-nr13` (`equipamento_nr13`), `/v1/publico/links-execucao/ativar` e `/atual/*` (`link_execucao`, rota pública), `/v1/bloqueios-despacho/consulta` (verificação calculada, sem entidade). `GET /v1/versao` já é exceção aceita (sumário 3.6) | Mantidas por clareza para integradores. Registrar como exceções aceitas no sumário ou renomear antes da publicação do contrato (para a consulta de despacho, sugestão: `/v1/consultas-despacho`) |
| A-7 | Cancelamento de ordem de reparo depois do envio | Proposta: `cancelar` só de `aberta`; depois do envio, `declarar-irreparavel` com `destino = devolvida_sem_reparo` |
| A-8 | Data de reconhecimento da solicitação | O sumário 3.3 diz que a triagem grava o instante de reconhecimento, sem nomear o campo. Proposto `data_reconhecimento`, como na OS |
| A-10 | Campos de LOC, TER e EST ainda não fixados pelo sumário: `periodo_franquia`, `valor_excedente_hora`, `valor_excedente_km`, `sla_substituicao_horas`, `situacao` do link, `status` da movimentação (`em_transito`, `recebida`) | Propostas; confirmar no documento 2 |
| A-17 | Armazenamento do `vinculo_os`: uma linha em cada conta ou registro na plataforma (escolha do documento 2) | A API é neutra: `vinculo.id` identifica o vínculo **na conta que consulta**; nenhum id de OS da outra conta sai (22.1) |
| A-18 | Convite de conexão modelado como criação (`POST /v1/conexoes-empresa`, ação `convidar`), não como sub-rota de ação | Segue a regra de criação da seção 2.1. Confirmar com o documento 4 |
| A-19 | Só o cliente convida (sumário 3.5) | A subcontratada não convida o cliente. Se o produto quiser o caminho inverso, acrescentar uma ação na Fase 2 |
| A-20 | Catálogo próprio da subcontratada | Proposto como `servico_fornecedor` com `fornecedor_id = null` (22.2). Confirmar no documento 2 |
| A-22 | API e webhooks no plano Serviços | O sumário 2.2 não os inclui. Proposta: token só de leitura no plano Serviços (22.10). Depende de P4 |
| A-23 | TAG de equipamento `de_cliente` na conta da subcontratada | Proposta: única por cliente (`cliente_servico_id` + `tag`), não por conta (22.4). Confirmar no documento 2 |
| A-24 | Subcontratada que não aceita executar uma OS vinculada | O sumário 3.6 já decidiu a ação `devolver` (OS recebida e não iniciada volta ao cliente como `aberta`, sem executor, com motivo). A rota ainda não está neste documento, que nesta versão tratou só da configuração por segmento. Proposta para a próxima versão: `POST /v1/ordens-servico/{id}/devolver`, na conta da subcontratada, sincronizada pelo grupo `liberacao` ou por um grupo próprio. Divergência registrada |
| A-25 | Conflito de transição entre contas (cliente cancela enquanto a subcontratada já iniciou) | Proposta da 22.6: a entrega falha, a OS do cliente volta ao estado sincronizado e a gestão é avisada. Validar no documento 4 |
| A-26 | Envio do relatório PDF ao cliente atendido | O sumário 3.5 diz que a subcontratada "envia ao cliente o relatório"; a API só gera o PDF (22.7). Avaliar `enviar_para` com e-mail do contato marcado `recebe_relatorio` |
| A-27 | Equipamento `de_cliente` do tipo `veiculo` na conta da subcontratada | Proposta: não conta como veículo cobrado nem entra nos KPIs de frota da subcontratada (22.4). Confirmar no sumário 2.2 |
| A-28 | `reabrir-como-nova` | Rota proposta para a regra "reabrir como nova" do sumário 3.6 (7.1). Confirmar no documento 4 |
| A-30 | Campo extra obrigatório que o executor não preenche | Prestador e link só veem os campos com `preenchido_pelo_prestador`, e os campos extras não trafegam na conexão. Proposta (8.2, 23.2): esses executores concluem sem o campo, a OS fica com o aviso `campos_extras_pendentes` e `encerrar` passa a exigir o campo. Confirmar no documento 4 |
| A-31 | Fase de ciclo de vida fora da conexão | A cópia `de_cliente` na subcontratada fica `em_operacao`, e uma OS `implantacao` do cliente chega lá num ativo que, para a conta dela, já opera. Proposta (22.8): a cópia aceita o `tipo` sem a validação de fase. Confirmar no documento 4 |
| A-33 | Status do ativo na conexão | O RF-ADM-016 diz que o efeito da regra (ativo `parado`) chega à outra conta "pelo campo de status já compartilhado", mas o grupo `ativo_identificacao` (22.5) não lista `status_operacional`. Decidir no sumário se o status entra no grupo (mudança compatível) ou corrigir o RF-ADM-016 |
| A-34 | Nomes desta versão sem nome canônico no sumário | Propostas: entidades `medicao_etapa` (24.7) e `lancamento_custo_projeto` (6.18); campos `preenchido_pelo_prestador`, `obrigatorio_em`, `opcoes`, `numero_versao`, `modelo_projeto_versao`, `origem_pacote`, `disparo_regra_id`, `levanta_com_tipo_documento_id`, `levantado_automaticamente`, `colocado_em_operacao_por_id`, `termo_entrega_anexo_id`, `causa_externa_marcada_por_id`, `causa_externa_marcada_em`, `numero_contrato_cliente`, `renovacao_automatica`, `evidencia_causa_externa_obrigatoria`; valores `calendario_contrato` (calendário de apuração) e `em_apuracao`/`fechada` na medição de etapas; aviso `campos_extras_pendentes`; escopos `configuracao:*`, `projetos:*`, `contratos_cliente:*`; evento `regra_automacao.disparada`. Fixar no documento 2 e, os de vocabulário, no sumário |
| A-35 | Operações do documento 1 ainda sem rota | Aditivo de contrato com cliente com histórico (RF-CTR-001, proposta) e marcação de medição como enviada ao cliente ou ao ERP (RF-PRJ-006, RF-CTR-004, propostas). Hoje o `PATCH` de contrato `vigente` só altera data de fim, renovação, contatos e anexos (25.2) |
| A-36 | Alcance de `prj.gerenciar` e `ctr.gerenciar` | O sumário 3.6 lista modelos e projetos (criar, iniciar, suspender, retomar, concluir, cancelar) e contratos e itens. Proposta (4.8): `prj.gerenciar` cobre também etapas, `gerar-os`, lançamentos de custo e a geração e o fechamento da medição de etapas; `ctr.gerenciar` cobre também abrir a medição do contrato. Confirmar no sumário |
| A-37 | Lacunas da implementação de E01 na camada HTTP | (a) Idempotência guardada na memória de cada instância (2.9): correta só com uma instância; antes de escalar a API, o depósito precisa ser compartilhado (tabela no banco, que é do documento 2, ou Valkey); a obrigatoriedade da chave para o app e para tokens de API também fica para depois. (b) Rate limit da seção 5 ainda não implementado, inclusive o de login e de `/oauth/token`; hoje só vale o bloqueio progressivo da identidade. (c) `X-Forwarded-For` não é lido: o IP gravado é o da conexão; atrás do balanceador, é preciso a lista de proxies confiáveis. (d) `If-None-Match` e `304` não implementados (2.11). (e) Nas ações de estado e no `DELETE`, a versão do `If-Match` é conferida por leitura antes da ação, não na mesma transação. (f) `PATCH /v1/usuarios/{id}` com `perfis` e dados roda em duas transações (perfis primeiro, depois dados), e a troca de perfis não muda a versão do usuário. (g) Erro de `/oauth/token` com `401`, e não o `400` da RFC 6749 para `invalid_grant`, para manter o código `credencial_invalida` com um status só |
| A-38 | Campos de `GET /v1/versao` | A seção 18.8 lista `ambiente`; a seção 9.5 do documento 5 e o servidor devolvem só `versao`, `commit` e `data_build`. Decidir qual vale e alinhar o outro |
| A-39 | Escopo `link` fora do alcance | A seção 4.9 manda responder `404` também ao escopo `link` fora das rotas dele, mas hoje o código responde `403 permissao_negada`, porque `Ator.Externo()` não trata `link` como externo. Sem efeito em E01: a sessão do link de uso único só nasce em E06/E09. Na etapa do link, `Ator.Externo()` deve tratar `link` como externo, com `404` fora do alcance. A regra da seção 4.9 não muda |

**Divergências do documento 1 (seção 8.3) tratadas na versão 0.5:** ponto 39, grupo **Conexão** (`srv.gerenciar_conexao`, padrão `administrador` e `gestor`) nas rotas de conexão da 6.10, no lugar de "Admin", mantendo o lado de cada ação; ponto 42, `prestador` nomeado nas rotas de modelos e execuções de checklist (6.7), com o alcance das OS do seu fornecedor (4.9); ponto 43, `plh` nomeado em `PUT /v1/ativos/{id}/enquadramento-nr13` (6.13); ponto 44, `POST /v1/ordens-reparo/{id}/registrar-orcamento` no grupo das rotas de orçamento de fornecedor (Gestão, Compras, `prestador`; permissão `ter.registrar_orcamento`). A seção 4.8 passa a remeter ao catálogo de RF-ADM-005 e lista `srv.gerenciar_conexao` entre as permissões finas.

### 26.2 Pendências do sumário que afetam a API

| Pendência | Efeito na API |
|---|---|
| P1 — ICP-Brasil no item 1.6 da NR-1 | Forma de `assinatura_anexo_id` no encerramento NR-12, no registro pelo link de uso único (19.4) e nos relatórios NR-13. A assinatura eletrônica avançada já está em todos os planos no MVP; ICP-Brasil pode exigir endpoint de assinatura com certificado (A1/A3 ou nuvem) |
| P2 — Tabela de validade de CIV e CIPP | Hoje `documentos-vencimento` recebe a validade digitada; se a tabela for confirmada, entra um campo `validade_sugerida` calculado |
| P3 — API Senatran/Serpro | Sem API, `POST /v1/multas` e `indicar-condutor` são registros manuais; com API, entra um conector de importação |
| P5 — Nome e domínio | URLs base, prefixos de token (`tnz_`, `lx_`), nomes de cabeçalho (`Tenaz-*`, `X-Tenaz-*`) |
| P6 — Prazos NR-13 com SPIE | Valores devolvidos por `GET /v1/prazos-nr13` |
| P7 — Local da cópia de recuperação de desastre | Sem efeito no contrato da API. Pode afetar a documentação pública (região, RTO de região publicado) |
| P8 — CFOPs de remessa, retorno, conserto e locação | Pode acrescentar `cfop` e validações por tipo em `movimentacoes-estoque` e `ordens-reparo` (20.1). Hoje a API só registra `nfe_numero` e `nfe_chave` |
| P4 — Validação de preço e de prioridades | Plano Serviços (preço por técnico, desconto de canal) e inclusão de API de leitura nesse plano (22.10, A-22) |
| R4 — Prestador vendo dados de outra empresa ou de outro prestador | Filtro por fornecedor além do RLS e suíte de testes de vazamento em todo build (4.9 e 19.6) |
| R5 — Sincronização entre contas copia dado além do permitido ou para a conta errada | Lista fechada de grupos em `campos_compartilhados` (22.5), filtro na saída e na entrada, transações separadas por conta com papel próprio, nenhum id de outra conta na API (22.1), auditoria nos dois lados (22.6). A suíte de vazamento ganha um cenário com três contas (cliente, subcontratada conectada e terceira conta) sobre todos os endpoints e eventos de webhook |
| R6 — SRV no MVP aumenta o escopo | A API reaproveita OS, apontamento, ativos, pessoas e sincronização do app; o SRV soma 23 rotas próprias e 6 rotas novas em OS (6.20) |
| P9 — Regras metrológicas de medidores de velocidade | Conteúdo do pacote `fiscalizacao_eletronica` (23.4): tipo de documento "verificação metrológica" com validade digitada e as duas regras de bloqueio entregues `inativa` até a decisão. Sem efeito no contrato da API |
| P10 — Acesso do cliente final de um contrato | Hoje nenhum escopo externo acessa `contratos-cliente`, `medicoes-contrato` nem o painel (4.9); a medição sai em PDF e XLSX (25.5). Se o acesso for aprovado, entra um escopo de acesso novo (decidido no sumário 3.6), com filtro por `cliente_servico` e cenário próprio na suíte de vazamento |
| R7 — Configuração que vira sistema "faça você mesmo" | Catálogos fechados de condição, efeito, tipo de campo e tipo de item de pacote, sem enum extensível por configuração; estados, fórmulas e grupos da conexão sem rota de configuração (seção 23) |

### 26.3 Decisões a validar

Todas as marcações "proposta" deste documento, em especial: limites de rate limit por plano e das rotas públicas (seção 5); retenção do cursor de sincronização (30 dias) e do `id_comando` (90 dias); regras de conflito do offline (12.5); calendário de novas tentativas, DLQ de 30 dias e suspensão de webhooks (13.5 e 13.6); prazo de suporte da `/v1` após a `/v2` (18 meses); dinheiro e quantidades como string decimal (2.5); MFA recente para exportação (17.1); validade padrão e entropia do link de uso único e duração da sessão curta (30 minutos, até 4 horas) (19.4); confirmação de recebimento de remessa (20.2); alcance exato de `prestador` e `locatario` (4.9). No SRV: validade do convite de conexão (14 dias) e formato do código (22.5); grupos padrão de `campos_compartilhados` (22.5); meta de 60 segundos de atraso da sincronização entre contas e regra de conflito de transição (22.6); limite de 24 fotos no relatório PDF (22.7); aviso `tecnico_sem_especialidade` em vez de erro (22.3). Na configuração por segmento (versão 0.4): regra nascendo `inativa`, combinações permitidas de condição e efeito, janela de 30 dias da simulação e limite de 50 regras ativas (23.1); limite de 30 campos extras por combinação, nomes `obrigatorio_em` e `preenchido_pelo_prestador`, obrigatoriedade avaliada na data do comando offline e conclusão com `campos_extras_pendentes` (23.2, 12.4, 8.2); exigências de `colocar-em-operacao` sobre OS de implantação, campos e documentos da classe, e o lote de até 100 ativos (23.3); geração de uma OS por ativo, conclusão de projeto tudo ou nada, entidade `medicao_etapa` com `em_apuracao` e `fechada`, lançamentos de custo do projeto (24); prazos medidos no calendário de apuração, abertura automática da medição e `divergencias` na reabertura (25.5).

---

## Histórico de revisões

| Versão | Data | Autor | Alteração |
|---|---|---|---|
| 0.1 | 12/11/2026 | — | Primeira versão: convenções, autenticação, rate limit, catálogo por módulo, sincronização offline, webhooks, integrações, bloqueio de despacho, exportação e política de versões |
| 0.2 | 12/11/2026 | — | Alinhamento ao sumário executivo 0.3: prioridade `urgente`; transições da OS (`desprogramar`, `recusar-liberacao`, cancelamento só de `aberta`, `planejada`, `programada`, `pausada`); regras do status do ativo; `pessoa` e `/v1/pessoas`; perfis `comprador`, `plh`, `prestador`, `locatario` e alcance dos externos (4.9); campos novos do ativo e da OS; origem de leitura `iot`; eventos com nome completo da entidade; planos conforme a tabela do sumário 2.1 (catálogo padrão de falhas e registro de bloqueio no Essencial, MCP no Enterprise, Essencial com 1 unidade). Módulos novos TER e LOC e EST estendido (seções 19 a 21), com eventos e códigos de erro novos. Versão do produto em SemVer, `GET /v1/versao` e cabeçalho `Tenaz-Versao` (18.8). Pendências renumeradas para a seção 22, com as inconsistências da versão 0.1 marcadas como resolvidas |
| 0.3 | 12/11/2026 | — | Alinhamento ao sumário executivo 0.4. Permissões `<módulo>.<ação>`: `os.encerrar` (substitui `ordem_servico.liberar`) e `est.decidir_destino`. Escopo de acesso `interno`, `prestador`, `locatario`, `link`. Link de uso único com token no fragmento, `POST /v1/publico/links-execucao/ativar` e sessão curta (`/atual/*`); nenhum token no caminho. Ação `assumir`; `origem` da OS com 8 valores; solicitação `nova` → `recusada`, `duplicada` final e `reabrir-como-nova`, origem `portal_locatario`. Locação: `local_uso`, `franquia_horas`, `franquia_km`, cancelamento só de `rascunho`, sublocação proibida. Papéis `posto_combustivel` e `seguradora`; peça `em_transito` não consumível; `/v1/bloqueios-ativo` e eventos `ativo.despacho_*`. Módulo novo **SRV** (seções 6.10 e 22): especialidades, clientes atendidos, equipamentos `de_cliente`, serviços a realizar, componentes por `tipo_material`, relatório PDF, conexão entre contas, vínculo de OS com dono por campo, API para integradores da subcontratada. Eventos `conexao_empresa.*` e `vinculo_os.*`; 9 códigos de erro novos. Catálogo de 451 para 482 endpoints. Pendências na seção 23 |
| 0.4 | 12/11/2026 | — | Alinhamento ao sumário executivo 0.5: configuração por segmento (sumário 3.6 e 3.7). Regras configuráveis (`/v1/regras-automacao`, `/v1/disparos-regra`, simulação, ativar e inativar) e campos extras da OS (`/v1/campos-extras-os`, `/v1/ordens-servico/{id}/campos-extras`, `campos_extras` com marcação "preenchido pelo prestador") no MVP; fase de ciclo de vida do ativo (`fase_ciclo_vida`, `data_entrada_operacao`, `colocar-em-operacao` unitário e em lote) e tipo de OS `implantacao`; pacotes de segmento com aplicação atômica (`/v1/pacotes-segmento`, `/v1/aplicacoes-pacote`, Fase 2). Módulos novos **PRJ** (6.18 e 24: modelos, projetos com ações, etapas `pendente`, `disponivel`, `em_andamento`, `concluida`, `dispensada`, OS de implantação, termo de entrega, medição de etapas) e **CTR** (6.19 e 25: contratos com clientes, itens, painel, medições com `decisao_exclusao`). OS com origem `projeto` (e `cliente_conectado`, resolvendo A-21), `causa_externa`, `etapa_projeto_id`; sem origem `regra` na OS; solicitação com origem `regra`; `bloqueio_ativo.origem`; `/v1/clientes-servico` também no Profissional com CTR. Nove permissões novas (`adm.gerenciar_regras`, `adm.gerenciar_configuracao`, `atv.colocar_em_operacao`, `prj.gerenciar`, `prj.dispensar_etapa`, `prj.ajustar_medicao`, `ctr.gerenciar`, `ctr.fechar_medicao`, `ctr.reabrir_medicao`), escopos de token `configuracao:*`, `projetos:*`, `contratos_cliente:*` (proposta), 9 eventos de webhook novos (`regra_automacao.disparada` como proposta) e 24 códigos de erro novos. Catálogo de 482 para 564 endpoints. Seções novas 23 a 25; pendências renumeradas para a seção 26, com A-29 a A-36 |
| 0.5 | 12/11/2026 | — | Contrato de E01 (conta e acesso), gerado do código e alinhado ao sumário 0.6. Códigos de erro do código Go alinhados à seção 3.2 (snake_case: `recurso_nao_encontrado`, `permissao_negada`, `validacao_falhou`, `versao_desatualizada`, `limite_requisicoes_excedido`, `erro_interno`...) e códigos novos `senha_fraca`, `mfa_codigo_invalido`, `link_email_invalido`, `ultimo_administrador`, `pre_condicao_obrigatoria`; problema com `erros[].codigo`, `id_rastreio` e membros `desafio_id`, `desafio_expira_em`, `tentativas_restantes`, `versao_atual`, `error`. 19 rotas novas de ADM (convites, MFA do próprio usuário, troca de senha, `POST /oauth/authorize`, `POST /v1/empresas`, turnos, feriados, reativações, reenvio de convite, permissões avulsas) e `DELETE /v1/unidades/{id}` no MVP; catálogo de 564 para 583. Sessão web (cookie de desenvolvimento, CSRF, Origin, `desafio_id` como segredo, formato da sessão), fluxo do authorize, `empresa_escolhida_id`, idempotência e `If-Match` da E01, extensões `x-escopos-acesso` e `x-permissoes`. Pontos 39, 42, 43 e 44 da seção 8.3 do documento 1: grupo Conexão, `prestador` em checklists, `plh` no NR-13, grupo da rota de orçamento da ordem de reparo; seção 4.8 remete a RF-ADM-005. Esquema de `GET /v1/versao` chamado `Versao`. Pendências A-37 e A-38 |
| 0.5 | 13/11/2026 | — | Alinhamento ao código e ao contrato regerado de E01 (mesma versão): `mfa_obrigatorio` em `GET` e `PATCH /v1/empresa`; `mfa_pendente` na sessão (login, `mfa/verificar`, `GET /v1/auth/sessao`, troca de empresa) e como membro do `401 nao_autenticado` (3.1, 3.2); sessão web restrita à ativação do MFA com as cinco rotas que ela atende, recusa no app e em `/oauth/token`, `409 transicao_invalida` com `estado_atual = "exigido_pela_empresa"` em `mfa/desativar` e passo TOTP de uso único (RN-ADM-027, 4.2); `pessoa.situacao` só `ativa` e `inativa`; regras de RN-ADM-020 e RN-ADM-021 em `POST /v1/usuarios` e `DELETE /v1/pessoas/{id}` (6.1) |
| 0.6 | 13/11/2026 | — | Alinhamento ao sumário 0.7 e ao documento 1 versão 0.6: a Tenaz não guarda dado de saúde. Sai `pessoa.saude.ler` da tabela de permissões finas (4.8, agora doze) e a exceção ao formato `<módulo>.<ação>`; texto ajustado em 4.4 (escopos de token e ações sensíveis), 6.1 (detalhe da pessoa, permissão avulsa de exemplo, "Pessoa e usuário"), 6.13 (agenda de documentos), 6.15 e 13.6 ("saúde" da integração passa a "indicadores de entrega" e "situação de funcionamento"), 17.3 (pacote da exportação), 19.3 e 19.4 (exemplos de `observacao_fornecedor`) e 22.10 |
| 0.6 | 13/11/2026 | — | Seção 4.3: regras de verificação do access token do app (`exp` obrigatório, tolerância de 30 s, `iat` e `nbf` no futuro recusados, `iss` `tenaz`, `aud` `tenaz-api`, cabeçalho `alg: EdDSA`, `typ: JWT`, `kid` conhecido; `nbf` respeitado quando presente), depois da correção do defeito em que o servidor ignorava `nbf`. Pendência A-39 (escopo `link` recebe `403` em vez de `404` até a etapa do link) |
