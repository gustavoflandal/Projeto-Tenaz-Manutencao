-- Trilha de auditoria e procedimento de preparação de tabela.
--
-- Fonte: docs/2_ARQUITETURA_BANCO_DADOS.md, seções 2.10.3, 4.5 e 9.2 a 9.5.
--
-- Vem antes de qualquer tabela de negócio porque preparar_tabela liga o
-- trigger de auditoria: uma tabela criada antes desta migração nasceria sem
-- trilha, e a trilha não se reconstrói depois.
--
-- Particionamento sem pg_partman: a extensão não acompanha a imagem oficial do
-- PostgreSQL (pendência da etapa 011). As partições mensais são criadas por
-- plataforma.criar_particoes_mensais(), com uma partição padrão como rede de
-- segurança para o caso de o job mensal falhar.

-- +goose Up

-- +goose StatementBegin
-- Todo objeto criado a partir daqui pertence a tenaz_dono (seção 4.1).
-- SET LOCAL nos limites: goose roda cada arquivo numa transação, e o limite não
-- deve sobreviver na conexão que volta ao pool do migrador.
SET ROLE tenaz_dono;
SET LOCAL lock_timeout = '3s';
SET LOCAL statement_timeout = '60s';
-- +goose StatementEnd

-- +goose StatementBegin
-- Conexão em processamento pelo sincronizador entre contas (seção 2.10.3).
-- Só a leitura da variável; plataforma.conferir_conexao() entra com a
-- conexão entre contas, em E04. Fora do sincronizador devolve nulo, e é com
-- esse nulo que tg_auditar grava conexao_empresa_id.
CREATE OR REPLACE FUNCTION app.conexao_atual() RETURNS uuid
  LANGUAGE sql STABLE PARALLEL SAFE
AS $$ SELECT nullif(current_setting('app.conexao_empresa_id', true), '')::uuid $$;
-- +goose StatementEnd

-- +goose StatementBegin
-- Tabelas somente inserção: evento_auditoria agora; transicao_status_os,
-- assinatura_eletronica e outras depois. Vale até para o dono e para o
-- migrador: desligar o trigger é DDL, que fica no histórico (seção 9.4).
CREATE OR REPLACE FUNCTION app.tg_somente_insercao() RETURNS trigger
  LANGUAGE plpgsql
AS $$
BEGIN
  RAISE EXCEPTION 'a tabela % é somente inserção', TG_TABLE_NAME USING ERRCODE = '42501';
END
$$;
-- +goose StatementEnd

-- +goose StatementBegin
-- Evento de auditoria (seção 9.2). Particionada por mês em ocorrido_em.
-- ocorrido_em é o instante do servidor (início da transação); o instante do
-- aparelho, quando a ação foi feita offline, vai em ocorrido_no_dispositivo_em
-- (RF-ADM-011).
CREATE TABLE app.evento_auditoria (
  empresa_id                 uuid        NOT NULL,
  id                         uuid        NOT NULL DEFAULT app.uuid_v7(),
  ocorrido_em                timestamptz NOT NULL DEFAULT now(),
  ator_tipo                  text        NOT NULL,
  usuario_id                 uuid,
  token_api_id               uuid,
  acao                       text        NOT NULL,
  entidade                   text        NOT NULL,     -- nome da tabela ou do recurso
  entidade_id                uuid,
  entidade_chave             text,                     -- chave não uuid (ex.: parametro_empresa.chave)
  antes                      jsonb,
  depois                     jsonb,
  campos_alterados           text[],
  origem                     text,
  ip                         inet,
  dispositivo_id             uuid,
  user_agent                 text,
  request_id                 text,                     -- na sincronização entre contas: o mensagem_id
  conexao_empresa_id         uuid,                     -- evento gravado pela sincronização entre contas
  ocorrido_no_dispositivo_em timestamptz,              -- instante do aparelho (ação offline)
  escopo                     text,                     -- escopo de acesso da credencial (RF-ADM-011); nulo sem contexto
  transacao                  xid8        NOT NULL DEFAULT pg_current_xact_id(),
  CONSTRAINT pk_evento_auditoria PRIMARY KEY (empresa_id, ocorrido_em, id),
  CONSTRAINT ck_evento_auditoria_ator CHECK (ator_tipo IN ('usuario', 'sistema', 'integracao', 'suporte_tenaz',
    'link_execucao', 'sincronizacao_conexao')),
  CONSTRAINT ck_evento_auditoria_acao CHECK (acao IN ('inclusao', 'alteracao', 'exclusao', 'login', 'login_falho',
    'logout', 'exportacao', 'visualizacao_sensivel', 'acesso_suporte', 'revisao', 'sincronizacao')),
  -- evento da sincronização sempre diz de qual conexão veio
  CONSTRAINT ck_evento_auditoria_conexao CHECK (ator_tipo <> 'sincronizacao_conexao' OR conexao_empresa_id IS NOT NULL),
  CONSTRAINT ck_evento_auditoria_origem CHECK (origem IN ('web', 'app', 'api', 'job', 'importacao')),
  CONSTRAINT ck_evento_auditoria_escopo CHECK (escopo IN ('interno', 'prestador', 'locatario', 'link'))
) PARTITION BY RANGE (ocorrido_em);
-- +goose StatementEnd

-- +goose StatementBegin
-- Histórico de um registro na tela do próprio registro (RF-ADM-011).
CREATE INDEX ix_evento_auditoria_entidade ON app.evento_auditoria (empresa_id, entidade, entidade_id, ocorrido_em DESC);
-- Filtro por usuário na tela de auditoria. O filtro por período usa a chave primária.
CREATE INDEX ix_evento_auditoria_usuario  ON app.evento_auditoria (empresa_id, usuario_id, ocorrido_em DESC);
-- +goose StatementEnd

-- +goose StatementBegin
-- Isolamento por comando, e não a policy isolamento_empresa FOR ALL: não existe
-- policy de UPDATE nem de DELETE. Mesmo com privilégio concedido por engano,
-- o RLS bloquearia a alteração (seção 9.4, camada 1).
ALTER TABLE app.evento_auditoria ENABLE ROW LEVEL SECURITY;
ALTER TABLE app.evento_auditoria FORCE  ROW LEVEL SECURITY;

CREATE POLICY auditoria_leitura ON app.evento_auditoria
  FOR SELECT TO tenaz_app, tenaz_job, tenaz_leitura
  USING (empresa_id = app.empresa_atual());

CREATE POLICY auditoria_insercao ON app.evento_auditoria
  FOR INSERT TO tenaz_app, tenaz_job, tenaz_sinc
  WITH CHECK (empresa_id = app.empresa_atual());

-- Escopo externo: o trigger grava o evento na transação do prestador, do
-- locatário ou do link, mas nenhum deles lê a trilha (seção 2.9.4).
CREATE POLICY escopo_externo_leitura ON app.evento_auditoria AS RESTRICTIVE
  FOR SELECT TO tenaz_app, tenaz_job, tenaz_leitura
  USING (app.escopo_interno());

-- Escopo externo só grava evento em nome de si mesmo: o usuario_id do evento é
-- o da transação (nulo no link). Sem isto, um prestador poderia inserir na
-- trilha um evento atribuído a outra pessoa (achado ISO-E01, sugestão da
-- revisão de isolamento). O trigger já grava app.usuario_atual(); a policy
-- trava o INSERT direto.
CREATE POLICY escopo_externo_insercao ON app.evento_auditoria AS RESTRICTIVE
  FOR INSERT TO tenaz_app, tenaz_job
  WITH CHECK (app.escopo_interno()
              OR (usuario_id IS NOT DISTINCT FROM app.usuario_atual()
                  AND escopo IS NOT DISTINCT FROM app.escopo_atual()));
-- O escopo gravado também é o da transação: um prestador não registra evento
-- como se fosse interno.

-- Evento de autenticação (entidade 'identidade': login, falha, logout, troca
-- de senha, MFA) é só do tenaz_autenticacao (ADR-011; documento 5, seção 6.9).
-- Os papéis de negócio e o sincronizador não o gravam nem por INSERT direto
-- (ISO-E01-07). O tg_auditar nunca produz essa entidade: identidade está em
-- plataforma, sem trigger de auditoria, e não há tabela app.identidade.
CREATE POLICY auditoria_sem_autenticacao ON app.evento_auditoria AS RESTRICTIVE
  FOR INSERT TO tenaz_app, tenaz_job, tenaz_sinc
  WITH CHECK (entidade <> 'identidade');

-- O sincronizador grava, mas não lê a trilha, e só eventos da conexão em
-- processamento (seção 2.10.4).
CREATE POLICY conexao_sinc ON app.evento_auditoria AS RESTRICTIVE
  FOR INSERT TO tenaz_sinc
  WITH CHECK (conexao_empresa_id = app.conexao_atual() AND ator_tipo = 'sincronizacao_conexao');

-- Os privilégios padrão do esquema dariam UPDATE e DELETE; eles saem aqui.
REVOKE ALL ON app.evento_auditoria FROM PUBLIC, tenaz_app, tenaz_job, tenaz_leitura;
GRANT SELECT, INSERT ON app.evento_auditoria TO tenaz_app, tenaz_job;
GRANT SELECT ON app.evento_auditoria TO tenaz_leitura;
-- O sincronizador recebe INSERT por coluna, como em toda tabela que ele toca
-- (seção 2.10.4): exatamente as colunas que tg_auditar lista no INSERT. Faltando
-- uma, toda escrita auditada do sincronizador falharia com 42501, mesmo com a
-- coluna nula (ISO-E01-05). As demais (id, ocorrido_em, transacao) ficam com o
-- valor padrão.
GRANT INSERT (empresa_id, ator_tipo, usuario_id, token_api_id, acao, entidade, entidade_id, entidade_chave,
              antes, depois, campos_alterados, origem, ip, dispositivo_id, user_agent, request_id,
              conexao_empresa_id, ocorrido_no_dispositivo_em, escopo)
  ON app.evento_auditoria TO tenaz_sinc;

-- Trigger de linha em tabela particionada é clonado em cada partição, inclusive
-- nas criadas depois.
CREATE TRIGGER trg_somente_insercao BEFORE UPDATE OR DELETE ON app.evento_auditoria
  FOR EACH ROW EXECUTE FUNCTION app.tg_somente_insercao();
-- TRUNCATE ignora RLS e não dispara trigger de linha: sem trava própria, o
-- dono (e o migrador) apagaria a trilha num comando. A trava fica em cada
-- partição, criada por criar_particoes_mensais, e não no pai: trigger de
-- comando não é clonado nas partições, e o TRUNCATE do pai recursa nelas e
-- dispara o trigger de cada uma — conferido no banco.
-- +goose StatementEnd

-- +goose StatementBegin
-- Criação das partições mensais (seção 9.5), no lugar do pg_partman.
--
-- SECURITY DEFINER, de dono tenaz_dono: a partição precisa ter o mesmo dono do
-- pai, senão escapa do RLS forçado e da suíte de propriedade. Quem chama não
-- importa — o job mensal roda como tenaz_manutencao, a migração como o dono.
--
-- Defesas:
--   - lista fechada de tabelas: a função não aceita tabela arbitrária, para não
--     virar um "crie qualquer coisa como dono" disponível a quem tem EXECUTE;
--   - search_path fixo e todo nome qualificado; pg_temp por último, para que
--     uma tabela temporária de mesmo nome não seja usada no lugar da real;
--   - EXECUTE revogado de PUBLIC e concedido só a tenaz_manutencao;
--   - CREATE TABLE + ATTACH PARTITION, e não CREATE TABLE ... PARTITION OF:
--     o ATTACH trava o pai em SHARE UPDATE EXCLUSIVE, e não em ACCESS
--     EXCLUSIVE, então a gravação de auditoria não para durante o job;
--   - toda partição nasce sem privilégio para os papéis da aplicação e com RLS
--     habilitado e forçado, sem policy: acesso direto à partição devolve zero
--     linhas mesmo que alguém conceda privilégio por engano. Pelo pai, valem as
--     policies do pai (o PostgreSQL não aplica as da partição nesse caminho).
--
--   - TimeZone fixo em UTC na função: os limites das partições são escritos
--     como texto, e não podem depender do fuso da sessão de quem chama;
--   - p_desde limitado a 120 meses antes do mês corrente: uma chamada não cria
--     centenas de tabelas por engano;
--   - nas tabelas somente inserção (hoje, a auditoria), cada partição recebe a
--     trava de TRUNCATE (trg_somente_insercao_truncate), que barra também o
--     TRUNCATE do pai, porque ele recursa nas partições.
--
-- Idempotente: partição que já existe é pulada. Devolve quantas criou.
-- Se a partição padrão já tiver linhas do mês a criar, o ATTACH falha com
-- 23514 e nada é perdido: o procedimento manual está na seção 9.5.
CREATE OR REPLACE FUNCTION plataforma.criar_particoes_mensais(
  p_tabela          regclass,
  p_meses_a_frente  integer DEFAULT 3,
  p_desde           date    DEFAULT NULL   -- primeiro mês; nulo = mês corrente (UTC)
) RETURNS integer
  LANGUAGE plpgsql VOLATILE SECURITY DEFINER
  SET search_path = pg_catalog, pg_temp
  SET lock_timeout = '3s'
  SET TimeZone = 'UTC'
AS $$
DECLARE
  v_esquema    text;
  v_nome       text;
  v_estrategia "char";
  v_mes        date;
  v_ultimo     date;
  v_particao   text;
  v_criadas    integer := 0;
  v_truncate   boolean;
BEGIN
  SELECT n.nspname, c.relname, pt.partstrat
    INTO v_esquema, v_nome, v_estrategia
  FROM pg_catalog.pg_class c
  JOIN pg_catalog.pg_namespace n ON n.oid = c.relnamespace
  LEFT JOIN pg_catalog.pg_partitioned_table pt ON pt.partrelid = c.oid
  WHERE c.oid = p_tabela;

  -- Lista fechada. Tabela particionada nova entra aqui por CREATE OR REPLACE,
  -- na migração que a cria, com revisão de isolamento.
  IF v_esquema IS NULL OR (v_esquema || '.' || v_nome) <> ALL (ARRAY['app.evento_auditoria']) THEN
    RAISE EXCEPTION 'tabela % fora da lista de particionamento mensal', p_tabela USING ERRCODE = '42501';
  END IF;
  IF v_estrategia IS DISTINCT FROM 'r' THEN
    RAISE EXCEPTION 'tabela % não é particionada por intervalo', p_tabela USING ERRCODE = '22023';
  END IF;
  IF p_meses_a_frente IS NULL OR p_meses_a_frente NOT BETWEEN 0 AND 24 THEN
    RAISE EXCEPTION 'p_meses_a_frente fora de 0 a 24: %', p_meses_a_frente USING ERRCODE = '22023';
  END IF;
  IF p_desde < (date_trunc('month', now() AT TIME ZONE 'UTC') - interval '120 months')::date THEN
    RAISE EXCEPTION 'p_desde anterior a 120 meses: %', p_desde USING ERRCODE = '22023';
  END IF;

  -- Tabelas somente inserção: cada partição ganha a trava de TRUNCATE.
  v_truncate := (v_esquema || '.' || v_nome) = ANY (ARRAY['app.evento_auditoria']);

  -- Partição padrão: rede de segurança. Recebe o que cair fora das mensais
  -- (job que não rodou, relógio adiantado). Deve ficar vazia; o monitoramento
  -- olha pg_stat_user_tables.n_live_tup, sem ler dado (seção 9.5).
  v_particao := v_nome || '_padrao';
  IF pg_catalog.to_regclass(format('%I.%I', v_esquema, v_particao)) IS NULL THEN
    EXECUTE format('CREATE TABLE %I.%I PARTITION OF %I.%I DEFAULT',
                   v_esquema, v_particao, v_esquema, v_nome);
    EXECUTE format('REVOKE ALL ON TABLE %I.%I FROM PUBLIC, tenaz_app, tenaz_job, tenaz_leitura, tenaz_sinc',
                   v_esquema, v_particao);
    EXECUTE format('ALTER TABLE %I.%I ENABLE ROW LEVEL SECURITY', v_esquema, v_particao);
    EXECUTE format('ALTER TABLE %I.%I FORCE ROW LEVEL SECURITY', v_esquema, v_particao);
    IF v_truncate THEN
      EXECUTE format('CREATE TRIGGER trg_somente_insercao_truncate BEFORE TRUNCATE ON %I.%I '
                     'FOR EACH STATEMENT EXECUTE FUNCTION app.tg_somente_insercao()', v_esquema, v_particao);
    END IF;
    v_criadas := v_criadas + 1;
  END IF;

  -- Meses em UTC, como todo instante do banco.
  v_ultimo := (date_trunc('month', now() AT TIME ZONE 'UTC') + make_interval(months => p_meses_a_frente))::date;
  v_mes    := date_trunc('month', coalesce(p_desde::timestamp, now() AT TIME ZONE 'UTC'))::date;

  WHILE v_mes <= v_ultimo LOOP
    v_particao := v_nome || '_' || to_char(v_mes, 'YYYY_MM');
    IF pg_catalog.to_regclass(format('%I.%I', v_esquema, v_particao)) IS NULL THEN
      EXECUTE format('CREATE TABLE %I.%I (LIKE %I.%I INCLUDING DEFAULTS INCLUDING CONSTRAINTS)',
                     v_esquema, v_particao, v_esquema, v_nome);
      EXECUTE format('REVOKE ALL ON TABLE %I.%I FROM PUBLIC, tenaz_app, tenaz_job, tenaz_leitura, tenaz_sinc',
                     v_esquema, v_particao);
      EXECUTE format('ALTER TABLE %I.%I ENABLE ROW LEVEL SECURITY', v_esquema, v_particao);
      EXECUTE format('ALTER TABLE %I.%I FORCE ROW LEVEL SECURITY', v_esquema, v_particao);
      IF v_truncate THEN
        EXECUTE format('CREATE TRIGGER trg_somente_insercao_truncate BEFORE TRUNCATE ON %I.%I '
                       'FOR EACH STATEMENT EXECUTE FUNCTION app.tg_somente_insercao()', v_esquema, v_particao);
      END IF;
      -- O ATTACH cria os índices e clona os triggers do pai.
      EXECUTE format('ALTER TABLE %I.%I ATTACH PARTITION %I.%I FOR VALUES FROM (%L) TO (%L)',
                     v_esquema, v_nome, v_esquema, v_particao,
                     (v_mes::timestamp AT TIME ZONE 'UTC'),
                     ((v_mes + interval '1 month')::timestamp AT TIME ZONE 'UTC'));
      v_criadas := v_criadas + 1;
    END IF;
    v_mes := (v_mes + interval '1 month')::date;
  END LOOP;

  RETURN v_criadas;
END
$$;
-- +goose StatementEnd

-- +goose StatementBegin
REVOKE ALL ON FUNCTION plataforma.criar_particoes_mensais(regclass, integer, date) FROM PUBLIC;
GRANT EXECUTE ON FUNCTION plataforma.criar_particoes_mensais(regclass, integer, date) TO tenaz_manutencao;
-- +goose StatementEnd

-- +goose StatementBegin
-- Chamada inicial: partição padrão, mês corrente e três meses à frente. A
-- partir daqui, o job mensal mantém sempre três meses prontos.
SELECT plataforma.criar_particoes_mensais('app.evento_auditoria', 3);
-- +goose StatementEnd

-- +goose StatementBegin
-- Trigger genérico de auditoria (seção 9.3). Roda com o papel de quem escreveu,
-- e não como dono: assim a policy auditoria_insercao confere que o evento é da
-- empresa da transação. O evento grava na mesma transação do fato; se o fato
-- for desfeito, o evento também é.
--
-- Contexto opcional lido aqui, além do de 00002: app.token_api_id,
-- app.dispositivo_id, app.user_agent e app.ocorrido_no_dispositivo_em
-- (seção 2.4). Ausentes, viram nulo.
CREATE OR REPLACE FUNCTION app.tg_auditar() RETURNS trigger
  LANGUAGE plpgsql
AS $$
DECLARE
  v_coluna_chave text   := TG_ARGV[0];
  v_ignorar      text[] := ARRAY['versao', 'alterado_xid', 'atualizado_em', 'atualizado_por'];
  -- nunca vão para a trilha, nem antes nem depois
  v_omitir       text[] := ARRAY['segredo_hash', 'token_hash', 'codigo_hash', 'csrf_hash', 'convite_hash',
                                 'senha_hash', 'mfa_segredo_cifrado'];
  v_antes        jsonb;
  v_depois       jsonb;
  v_linha        jsonb;
  v_campos       text[];
  v_chave        text;
  v_acao         text;
  v_uuid         constant text := '^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$';
BEGIN
  IF TG_OP = 'INSERT' THEN
    v_acao   := 'inclusao';
    v_depois := to_jsonb(NEW) - v_omitir - v_ignorar;
    v_linha  := to_jsonb(NEW);
  ELSIF TG_OP = 'DELETE' THEN
    v_acao  := 'exclusao';
    v_antes := to_jsonb(OLD) - v_omitir - v_ignorar;
    v_linha := to_jsonb(OLD);
  ELSE
    v_acao   := 'alteracao';
    v_antes  := to_jsonb(OLD) - v_omitir;
    v_depois := to_jsonb(NEW) - v_omitir;
    v_linha  := to_jsonb(NEW);
    SELECT array_agg(d.key ORDER BY d.key) INTO v_campos
    FROM jsonb_each(v_depois) d
    WHERE d.value IS DISTINCT FROM v_antes -> d.key
      AND NOT d.key = ANY (v_ignorar);
    IF v_campos IS NULL THEN
      RETURN NULL;                      -- nada relevante mudou
    END IF;
    SELECT jsonb_object_agg(k, v_antes -> k), jsonb_object_agg(k, v_depois -> k)
      INTO v_antes, v_depois
    FROM unnest(v_campos) AS k;
  END IF;

  v_chave := v_linha ->> v_coluna_chave;

  INSERT INTO app.evento_auditoria
    (empresa_id, ator_tipo, usuario_id, token_api_id, acao, entidade, entidade_id, entidade_chave,
     antes, depois, campos_alterados, origem, ip, dispositivo_id, user_agent, request_id,
     conexao_empresa_id, ocorrido_no_dispositivo_em, escopo)
  VALUES
    (coalesce(v_linha ->> 'empresa_id', v_linha ->> 'id')::uuid,   -- a tabela empresa não tem empresa_id
     coalesce(nullif(current_setting('app.ator_tipo', true), ''), 'sistema'),
     app.usuario_atual(),
     nullif(current_setting('app.token_api_id', true), '')::uuid,
     v_acao,
     TG_TABLE_NAME,
     CASE WHEN v_chave ~ v_uuid THEN v_chave::uuid END,
     CASE WHEN v_chave !~ v_uuid THEN v_chave END,
     v_antes, v_depois, v_campos,
     nullif(current_setting('app.origem', true), ''),
     nullif(current_setting('app.ip', true), '')::inet,
     nullif(current_setting('app.dispositivo_id', true), '')::uuid,
     nullif(current_setting('app.user_agent', true), ''),
     nullif(current_setting('app.request_id', true), ''),
     app.conexao_atual(),                      -- nulo fora do sincronizador
     nullif(current_setting('app.ocorrido_no_dispositivo_em', true), '')::timestamptz,
     app.escopo_atual());                     -- RF-ADM-011; nulo sem contexto
  RETURN NULL;
END
$$;
-- +goose StatementEnd

-- +goose StatementBegin
-- Preparação de tabela de negócio (seção 4.5). Toda migração que cria tabela
-- com empresa_id termina chamando este procedimento.
--
-- Recusa tabela cujo dono não seja tenaz_dono: o FORCE ROW LEVEL SECURITY não
-- vale para o dono, então uma tabela criada sem o SET ROLE da migração
-- contornaria o RLS. É a mesma regra que a suíte de vazamento cobra, aplicada
-- no instante em que o erro acontece, e não no CI.
CREATE OR REPLACE PROCEDURE app.preparar_tabela(
  p_tabela  regclass,
  p_auditar boolean DEFAULT true,
  p_chave   text    DEFAULT 'id',       -- coluna que identifica a entidade na auditoria
  p_escopo  text    DEFAULT 'interno'   -- 'interno': só usuários internos; 'proprio': a migração cria a policy (seção 2.9.4)
)
  LANGUAGE plpgsql
AS $$
BEGIN
  IF p_escopo NOT IN ('interno', 'proprio') THEN
    RAISE EXCEPTION 'p_escopo inválido: %', p_escopo USING ERRCODE = '22023';
  END IF;
  IF (SELECT c.relowner FROM pg_catalog.pg_class c WHERE c.oid = p_tabela) <> 'tenaz_dono'::regrole THEN
    RAISE EXCEPTION 'a tabela % não pertence a tenaz_dono; a migração esqueceu o SET ROLE tenaz_dono', p_tabela
      USING ERRCODE = '42501';
  END IF;

  EXECUTE format('ALTER TABLE %s ENABLE ROW LEVEL SECURITY', p_tabela);
  EXECUTE format('ALTER TABLE %s FORCE ROW LEVEL SECURITY', p_tabela);
  EXECUTE format('DROP POLICY IF EXISTS isolamento_empresa ON %s', p_tabela);
  EXECUTE format(
    'CREATE POLICY isolamento_empresa ON %s FOR ALL TO tenaz_app, tenaz_job, tenaz_leitura '
    'USING (empresa_id = app.empresa_atual()) WITH CHECK (empresa_id = app.empresa_atual())',
    p_tabela);
  IF p_escopo = 'interno' THEN
    EXECUTE format('DROP POLICY IF EXISTS escopo_externo ON %s', p_tabela);
    EXECUTE format(
      'CREATE POLICY escopo_externo ON %s AS RESTRICTIVE FOR ALL TO tenaz_app, tenaz_job, tenaz_leitura '
      'USING (app.escopo_interno()) WITH CHECK (app.escopo_interno())', p_tabela);
  END IF;
  EXECUTE format(
    'CREATE TRIGGER trg_padrao BEFORE INSERT OR UPDATE ON %s '
    'FOR EACH ROW EXECUTE FUNCTION app.tg_padrao()', p_tabela);
  IF p_auditar THEN
    EXECUTE format(
      'CREATE TRIGGER trg_auditoria AFTER INSERT OR UPDATE OR DELETE ON %s '
      'FOR EACH ROW EXECUTE FUNCTION app.tg_auditar(%L)', p_tabela, p_chave);
  END IF;
END
$$;
-- +goose StatementEnd

-- +goose StatementBegin
-- O procedimento é DDL de migração. Os privilégios padrão do esquema dariam
-- EXECUTE à aplicação; ele falharia de todo modo (só o dono altera a tabela),
-- mas o mínimo necessário é nenhum.
REVOKE ALL ON PROCEDURE app.preparar_tabela(regclass, boolean, text, text)
  FROM PUBLIC, tenaz_app, tenaz_job, tenaz_leitura;
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
SET ROLE tenaz_dono;
DROP PROCEDURE IF EXISTS app.preparar_tabela(regclass, boolean, text, text);
DROP FUNCTION IF EXISTS app.tg_auditar();
-- As partições caem junto com o pai.
DROP TABLE IF EXISTS app.evento_auditoria;
DROP FUNCTION IF EXISTS plataforma.criar_particoes_mensais(regclass, integer, date);
DROP FUNCTION IF EXISTS app.tg_somente_insercao();
DROP FUNCTION IF EXISTS app.conexao_atual();
-- +goose StatementEnd
