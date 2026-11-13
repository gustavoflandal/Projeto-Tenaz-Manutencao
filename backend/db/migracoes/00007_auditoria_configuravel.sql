-- Trilha de auditoria configurável (RF-ADM-011, ampliado na etapa E01).
--
-- Fonte: docs/2_ARQUITETURA_BANCO_DADOS.md, versão 0.7, seções 9.1 a 9.5 e 5.1;
-- pendências E01-18 e E01-19. Decisões do dono do produto: docs/prompts/
-- E01_conta-e-acesso.md, "Ajustes durante a etapa"; sumário 0.7, seções 3.2 e 3.6.
--
-- O que esta migração faz:
--   1. evento_auditoria ganha usuario_nome, tela, modulo, resultado, sensivel,
--      detalhe e hash; a lista de ações passa a ter os nove tipos de operação
--      (login, logout, inclusao, alteracao, exclusao, delecao, visualizacao,
--      impressao, extracao). Falha de login, acesso do suporte, revisão e
--      sincronização deixam de ser tipo e viram atributo do evento.
--   2. Os eventos já gravados são convertidos para os valores novos (a trilha
--      ainda não existe em homologação nem em produção; ver a nota do passo 4).
--   3. O hash de cada evento é calculado por trigger, no banco, sobre uma
--      serialização canônica que se recalcula fora dele (seção 9.2.1).
--   4. tg_auditar passa a gravar exclusao (exclusão lógica, com a linha inteira
--      em antes) e delecao (remoção física) — pendência E01-18.
--   5. captura_auditoria (captura por tela de visualização, impressão e
--      extração) e selo_auditoria (selo diário encadeado, somente inserção),
--      com o cálculo do selo e a verificação de integridade no próprio banco.
--   6. plataforma.acesso_empresa guarda o nome do usuário, para que o evento de
--      login, gravado por um papel que não lê app.usuario, leve o nome verdadeiro.
--
-- Compatibilidade com a versão anterior do código: NÃO há. O CHECK de acao
-- recusa os valores antigos, como a 00006 recusou pessoa.saude.ler. Nada desta
-- trilha existe em homologação ou produção; o código Go que grava os valores
-- novos entra no MESMO commit (documento 2, seção 14.3).
--
-- 00001 a 00006 estão em main e não se editam.

-- +goose Up

-- +goose StatementBegin
SET ROLE tenaz_dono;
SET LOCAL lock_timeout = '3s';
SET LOCAL statement_timeout = '60s';
-- +goose StatementEnd

-- ============================================================ 1. nome do usuário em acesso_empresa

-- +goose StatementBegin
-- O papel das credenciais grava o evento de login em cada empresa da identidade
-- (seção 9.1) e não lê app.usuario. O nome do usuário naquela empresa vem de
-- acesso_empresa, mantida pelo trigger de app.usuario na mesma transação do
-- UPDATE: o nome do evento é o do instante, e a policy autenticacao_identidade
-- confere que o nome gravado é o desta linha.
SET ROLE tenaz_credencial;
ALTER TABLE plataforma.acesso_empresa ADD COLUMN usuario_nome text;

CREATE OR REPLACE FUNCTION plataforma.tg_acesso_empresa() RETURNS trigger
  LANGUAGE plpgsql SECURITY DEFINER
  SET search_path = pg_catalog, pg_temp
AS $$
DECLARE
  v_habilitado boolean := NEW.situacao = 'ativo' AND NEW.excluido_em IS NULL;
BEGIN
  INSERT INTO plataforma.acesso_empresa AS a
    (identidade_id, empresa_id, usuario_id, usuario_nome, habilitado, atualizado_em)
  VALUES (NEW.identidade_id, NEW.empresa_id, NEW.id, NEW.nome, v_habilitado, pg_catalog.now())
  ON CONFLICT (identidade_id, empresa_id) DO UPDATE
    SET usuario_id    = EXCLUDED.usuario_id,
        usuario_nome  = EXCLUDED.usuario_nome,
        habilitado    = EXCLUDED.habilitado,
        atualizado_em = EXCLUDED.atualizado_em;

  IF NOT v_habilitado THEN
    WITH revogadas AS (
      UPDATE plataforma.sessao s
      SET revogada_em = pg_catalog.now()
      WHERE s.empresa_id = NEW.empresa_id AND s.usuario_id = NEW.id AND s.revogada_em IS NULL
      RETURNING s.id
    )
    UPDATE plataforma.token_renovacao t
    SET revogado_em = pg_catalog.now()
    WHERE t.familia_id IN (SELECT r.id FROM revogadas r) AND t.revogado_em IS NULL;
  END IF;
  RETURN NULL;
END
$$;
-- +goose StatementEnd

-- +goose StatementBegin
-- O trigger passa a disparar também quando o nome muda.
SET ROLE tenaz_dono;
DROP TRIGGER trg_acesso_empresa ON app.usuario;
CREATE TRIGGER trg_acesso_empresa AFTER INSERT OR UPDATE OF identidade_id, nome, situacao, excluido_em ON app.usuario
  FOR EACH ROW WHEN (NEW.identidade_id IS NOT NULL)
  EXECUTE FUNCTION plataforma.tg_acesso_empresa();
-- +goose StatementEnd

-- +goose StatementBegin
-- Preenchimento do que já existe, como quem aplica a migração: o dono das
-- credenciais não lê app.usuario, e o dono de app.usuario não lê as credenciais.
RESET ROLE;
UPDATE plataforma.acesso_empresa a
SET usuario_nome = u.nome
FROM app.usuario u
WHERE u.empresa_id = a.empresa_id AND u.id = a.usuario_id;

SET ROLE tenaz_credencial;
ALTER TABLE plataforma.acesso_empresa ALTER COLUMN usuario_nome SET NOT NULL;
SET ROLE tenaz_dono;
-- +goose StatementEnd

-- ============================================================ 2. colunas novas do evento

-- +goose StatementBegin
-- Tudo aceita nulo ou tem valor padrão (seção 14.3). resultado e sensivel têm
-- padrão constante: a coluna nasce sem reescrever a tabela.
ALTER TABLE app.evento_auditoria
  ADD COLUMN usuario_nome text,                                 -- nome do usuário no instante do evento
  ADD COLUMN tela         text,                                 -- '<modulo>.<tela>', catálogo fechado no Go
  ADD COLUMN modulo       text,                                 -- código do sumário (seção 2), minúsculas
  ADD COLUMN resultado    text    NOT NULL DEFAULT 'sucesso',   -- 'falha' só em login
  ADD COLUMN sensivel     boolean NOT NULL DEFAULT false,       -- o registro mostrado tem dado sensível (CPF completo etc.)
  ADD COLUMN detalhe      jsonb,                                -- filtros, formato e quantidade (impressão, extração); nunca o conteúdo
  ADD COLUMN hash         bytea;                                -- sha256 do conteúdo do evento (seção 9.2.1)

ALTER TABLE app.evento_auditoria DROP CONSTRAINT ck_evento_auditoria_acao;
-- +goose StatementEnd

-- ============================================================ 3. serialização canônica e hash

-- +goose StatementBegin
-- Hash do evento (seção 9.2.1): sha256 (UTF-8) de um objeto JSON que reúne TODOS
-- os campos do evento, menos o próprio hash.
--   - O objeto externo é escrito sem espaço, com as 29 chaves em ordem de bytes e
--     'formato' = 1 (versiona a regra: se ela mudar, o evento antigo continua
--     verificável pela versão dele).
--   - Cadeia: JSON escapado (aspas, barra invertida, \b \f \n \r \t, \u00xx
--     minúsculo nos demais controles, resto em UTF-8). Uuid em minúsculas;
--     instante em UTC com microssegundos e sufixo Z; ip como o texto do inet;
--     transacao como texto; boolean true/false; ausente = null.
--   - antes, depois, detalhe e campos_alterados entram como o texto do jsonb do
--     PostgreSQL (jsonb_out: chaves por tamanho e depois por bytes, ': ' e ', '
--     como separadores). Uma forma canônica própria custava 300 us por evento
--     (medido); esta custa cerca de 40 us, e o formato do jsonb é o do disco, estável.
-- Tudo é uma só concatenação, sem laço nem recursão.
CREATE OR REPLACE FUNCTION app.json_valor(p anyelement) RETURNS text
  LANGUAGE sql IMMUTABLE PARALLEL SAFE
AS $$ SELECT coalesce(pg_catalog.to_jsonb(p)::text, 'null') $$;

CREATE OR REPLACE FUNCTION app.hash_evento_auditoria(e app.evento_auditoria) RETURNS bytea
  LANGUAGE sql STABLE PARALLEL SAFE
AS $$
SELECT pg_catalog.sha256(pg_catalog.convert_to(pg_catalog.concat(
  '{"acao":',                       app.json_valor(e.acao),
  ',"antes":',                      app.json_valor(e.antes),
  ',"ator_tipo":',                  app.json_valor(e.ator_tipo),
  ',"campos_alterados":',           app.json_valor(e.campos_alterados),
  ',"conexao_empresa_id":',         app.json_valor(e.conexao_empresa_id),
  ',"depois":',                     app.json_valor(e.depois),
  ',"detalhe":',                    app.json_valor(e.detalhe),
  ',"dispositivo_id":',             app.json_valor(e.dispositivo_id),
  ',"empresa_id":',                 app.json_valor(e.empresa_id),
  ',"entidade":',                   app.json_valor(e.entidade),
  ',"entidade_chave":',             app.json_valor(e.entidade_chave),
  ',"entidade_id":',                app.json_valor(e.entidade_id),
  ',"escopo":',                     app.json_valor(e.escopo),
  ',"formato":1',
  ',"id":',                         app.json_valor(e.id),
  ',"ip":',                         app.json_valor(e.ip::text),
  ',"modulo":',                     app.json_valor(e.modulo),
  ',"ocorrido_em":',                app.json_valor(pg_catalog.to_char(e.ocorrido_em AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS.US"Z"')),
  ',"ocorrido_no_dispositivo_em":', app.json_valor(pg_catalog.to_char(e.ocorrido_no_dispositivo_em AT TIME ZONE 'UTC', 'YYYY-MM-DD"T"HH24:MI:SS.US"Z"')),
  ',"origem":',                     app.json_valor(e.origem),
  ',"request_id":',                 app.json_valor(e.request_id),
  ',"resultado":',                  app.json_valor(e.resultado),
  ',"sensivel":',                   CASE WHEN e.sensivel THEN 'true' ELSE 'false' END,
  ',"tela":',                       app.json_valor(e.tela),
  ',"token_api_id":',               app.json_valor(e.token_api_id),
  ',"transacao":',                  app.json_valor(e.transacao::text),
  ',"user_agent":',                 app.json_valor(e.user_agent),
  ',"usuario_id":',                 app.json_valor(e.usuario_id),
  ',"usuario_nome":',               app.json_valor(e.usuario_nome),
  '}'), 'UTF8'))
$$;
-- +goose StatementEnd

-- +goose StatementBegin
-- BEFORE INSERT: completa o que vem do contexto da transação e calcula o hash.
-- Sobrescreve sempre o hash que o chamador mandar. Roda com o papel de quem
-- grava (tenaz_app, tenaz_job, tenaz_sinc, tenaz_autenticacao): nenhum deles
-- precisa de INSERT na coluna hash, porque o privilégio de coluna só vale para
-- as colunas da lista do INSERT, e não para as que o trigger altera.
--
--   usuario_nome  só do contexto (app.usuario_nome) e só se o usuário do evento
--                 é o da transação; o papel das credenciais manda o dele, e a
--                 policy autenticacao_identidade o confere com acesso_empresa;
--   modulo, tela  do contexto (app.modulo, app.tela) quando o chamador não manda;
--                 evento da identidade é sempre do módulo adm (RF-ADM-006).
CREATE OR REPLACE FUNCTION app.tg_evento_auditoria_preparar() RETURNS trigger
  LANGUAGE plpgsql
AS $$
BEGIN
  IF NEW.usuario_id IS NOT NULL AND NEW.usuario_nome IS NULL
     AND NEW.usuario_id IS NOT DISTINCT FROM app.usuario_atual() THEN
    NEW.usuario_nome := nullif(current_setting('app.usuario_nome', true), '');
  END IF;
  IF NEW.modulo IS NULL THEN
    IF NEW.entidade = 'identidade' THEN
      NEW.modulo := 'adm';
    ELSE
      NEW.modulo := lower(nullif(current_setting('app.modulo', true), ''));
    END IF;
  END IF;
  IF NEW.tela IS NULL THEN
    NEW.tela := nullif(current_setting('app.tela', true), '');
  END IF;
  NEW.hash := app.hash_evento_auditoria(NEW);
  RETURN NEW;
END
$$;
-- +goose StatementEnd

-- ============================================================ 4. conversão dos eventos existentes

-- +goose StatementBegin
-- A trilha é somente inserção, até para o dono (seção 9.4). Esta é a única vez em
-- que ela é reescrita: o trigger é desligado por DDL, que fica neste arquivo e no
-- histórico de migrações, e religado na mesma transação. Nenhum selo existe
-- ainda (selo_auditoria nasce abaixo), então nenhuma cadeia se quebra. Em banco
-- com trilha grande, esta passagem seria feita por partição e em lotes
-- (seção 14.3); aqui a trilha ainda não tem dado de homologação nem de produção.
--
-- Mapa dos valores antigos (decisão do dono do produto):
--   login_falho            -> login, resultado = falha
--   acesso_suporte         -> login (o ator_tipo = suporte_tenaz já o distingue)
--   exportacao             -> extracao
--   visualizacao_sensivel  -> visualizacao, sensivel = true
--   revisao, sincronizacao -> alteracao
--
-- Roda como quem aplica a migração (BYPASSRLS): o dono não passa pelo RLS
-- forçado e a trilha não tem policy de UPDATE.
ALTER TABLE app.evento_auditoria DISABLE TRIGGER trg_somente_insercao;
RESET ROLE;
SET LOCAL statement_timeout = '15min';

UPDATE app.evento_auditoria e
SET resultado    = CASE WHEN e.acao = 'login_falho' THEN 'falha' ELSE 'sucesso' END,
    sensivel     = (e.acao = 'visualizacao_sensivel'),
    acao         = CASE e.acao
                     WHEN 'login_falho'           THEN 'login'
                     WHEN 'acesso_suporte'        THEN 'login'
                     WHEN 'exportacao'            THEN 'extracao'
                     WHEN 'visualizacao_sensivel' THEN 'visualizacao'
                     WHEN 'revisao'               THEN 'alteracao'
                     WHEN 'sincronizacao'         THEN 'alteracao'
                     ELSE e.acao
                   END,
    -- todas as tabelas e todos os eventos gravados até aqui são de Administração
    modulo       = 'adm',
    usuario_nome = (SELECT u.nome FROM app.usuario u WHERE u.empresa_id = e.empresa_id AND u.id = e.usuario_id);

-- o hash depende de todos os campos: só agora, com o resto já convertido
UPDATE app.evento_auditoria e SET hash = app.hash_evento_auditoria(e);

SET LOCAL statement_timeout = '60s';
SET ROLE tenaz_dono;
ALTER TABLE app.evento_auditoria ENABLE TRIGGER trg_somente_insercao;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE app.evento_auditoria
  ADD CONSTRAINT ck_evento_auditoria_acao CHECK (acao IN ('login', 'logout', 'inclusao', 'alteracao', 'exclusao',
    'delecao', 'visualizacao', 'impressao', 'extracao')),
  -- falha é atributo do login; operação que falha desfaz a transação e não deixa evento
  ADD CONSTRAINT ck_evento_auditoria_resultado CHECK (resultado IN ('sucesso', 'falha')
    AND (resultado = 'sucesso' OR acao = 'login')),
  ADD CONSTRAINT ck_evento_auditoria_sensivel CHECK (NOT sensivel OR acao IN ('visualizacao', 'impressao', 'extracao')),
  ADD CONSTRAINT ck_evento_auditoria_modulo CHECK (modulo IN ('adm', 'atv', 'med', 'sol', 'os', 'pre', 'chk', 'est',
    'ter', 'srv', 'prj', 'ctr', 'loc', 'fro', 'con', 'ind', 'int', 'ia', 'mob')),
  -- a tela é '<modulo>.<tela>' e começa pelo módulo do evento
  ADD CONSTRAINT ck_evento_auditoria_tela CHECK (tela IS NULL OR
    (modulo IS NOT NULL AND tela ~ '^[a-z]+\.[a-z0-9_]+(\.[a-z0-9_]+)*$' AND pg_catalog.split_part(tela, '.', 1) = modulo)),
  ADD CONSTRAINT ck_evento_auditoria_hash CHECK (pg_catalog.octet_length(hash) = 32);
ALTER TABLE app.evento_auditoria ALTER COLUMN hash SET NOT NULL;

-- Trigger de linha em tabela particionada é clonado em cada partição, inclusive
-- nas criadas depois. Vem DEPOIS do preenchimento acima.
CREATE TRIGGER trg_evento_preparar BEFORE INSERT ON app.evento_auditoria
  FOR EACH ROW EXECUTE FUNCTION app.tg_evento_auditoria_preparar();
-- +goose StatementEnd

-- +goose StatementBegin
-- Consulta da tela de auditoria por módulo, tela e tipo (RF-ADM-011, filtros).
-- Medidas e justificativa: documento 2, seção 9.2.
CREATE INDEX ix_evento_auditoria_modulo ON app.evento_auditoria (empresa_id, modulo, tela, ocorrido_em DESC);
CREATE INDEX ix_evento_auditoria_acao   ON app.evento_auditoria (empresa_id, acao, ocorrido_em DESC);
-- +goose StatementEnd

-- +goose StatementBegin
-- Autenticação (ADR-011): os mesmos três tipos que o papel já gravava, agora sem
-- login_falho (é login com resultado falha), e o nome do usuário no evento é o
-- de acesso_empresa naquela empresa.
ALTER POLICY autenticacao_identidade ON app.evento_auditoria
  WITH CHECK (ator_tipo = 'usuario' AND entidade = 'identidade'
              AND acao IN ('login', 'logout', 'alteracao')
              AND EXISTS (SELECT 1 FROM plataforma.acesso_empresa a
                          WHERE a.identidade_id = evento_auditoria.entidade_id
                            AND a.empresa_id = evento_auditoria.empresa_id
                            AND a.usuario_id = evento_auditoria.usuario_id
                            AND a.usuario_nome = evento_auditoria.usuario_nome));

-- Além das 13 colunas da 00005: nome, tela e resultado. O módulo é sempre adm
-- (o trigger o fixa) e o hash é do trigger.
GRANT INSERT (usuario_nome, tela, resultado) ON app.evento_auditoria TO tenaz_autenticacao;
-- +goose StatementEnd

-- ============================================================ 5. tg_auditar: exclusão e deleção

-- +goose StatementBegin
-- E01-18. Exclusão lógica (excluido_em de nulo para preenchido) é 'exclusao' e
-- leva a linha inteira em antes: a retirada de um perfil ou de uma permissão
-- avulsa passa a dizer qual saiu. Remoção física é 'delecao'. Reativação
-- (excluido_em preenchido para nulo) e demais mudanças seguem como 'alteracao'.
-- tela, modulo e usuario_nome vêm do trigger de BEFORE INSERT do evento, e não
-- daqui: o INSERT abaixo é o mesmo da 00003 e o GRANT do tenaz_sinc não muda.
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
    v_acao  := 'delecao';
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
    -- exclusão lógica: excluido_em era nulo (ou a tabela não o tem, e então nunca chega aqui) e agora tem valor
    IF coalesce(jsonb_typeof(to_jsonb(OLD) -> 'excluido_em'), 'null') = 'null'
       AND jsonb_typeof(to_jsonb(NEW) -> 'excluido_em') = 'string' THEN
      v_acao := 'exclusao';
      -- a linha inteira antes; depois leva só o que mudou, entre eles o excluido_em
      SELECT jsonb_object_agg(k, v_depois -> k) INTO v_depois FROM unnest(v_campos) AS k;
      v_antes := to_jsonb(OLD) - v_omitir - v_ignorar;
    ELSE
      SELECT jsonb_object_agg(k, v_antes -> k), jsonb_object_agg(k, v_depois -> k)
        INTO v_antes, v_depois
      FROM unnest(v_campos) AS k;
    END IF;
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

-- ============================================================ 6. captura_auditoria

-- +goose StatementBegin
-- Captura por tela (RF-ADM-011): uma linha por empresa e tela, com uma chave para
-- cada tipo configurável. Sem linha, vale o padrão da criação da empresa:
-- impressão e extração ligadas, visualização desligada. A linha só nasce quando
-- alguém muda uma chave; o Go grava as TRÊS chaves de uma vez (o UPDATE do upsert
-- não pode devolver a outra chave ao padrão). Ligar ou desligar é alteração
-- sempre gravada: o trg_auditoria registra entidade_chave = tela e o campo
-- alterado. Login, logout, inclusão, alteração, exclusão, deleção e a
-- visualização de dado sensível não aparecem aqui: nunca se desligam.
--
-- Escopo interno: o prestador, o locatário e o link não leem nem escrevem a
-- configuração; o Go a carrega num contexto interno por empresa.
CREATE TABLE app.captura_auditoria (
  empresa_id    uuid    NOT NULL,
  id            uuid    NOT NULL DEFAULT app.uuid_v7(),
  tela          text    NOT NULL,                  -- '<modulo>.<tela>'; catálogo fechado no Go
  visualizacao  boolean NOT NULL DEFAULT false,
  impressao     boolean NOT NULL DEFAULT true,
  extracao      boolean NOT NULL DEFAULT true,
  versao bigint NOT NULL DEFAULT 1, alterado_xid xid8 NOT NULL DEFAULT pg_current_xact_id(),
  criado_em timestamptz NOT NULL DEFAULT now(), criado_por uuid,
  atualizado_em timestamptz NOT NULL DEFAULT now(), atualizado_por uuid, excluido_em timestamptz,
  CONSTRAINT pk_captura_auditoria PRIMARY KEY (empresa_id, id),
  CONSTRAINT uq_captura_auditoria_tela UNIQUE (empresa_id, tela),
  CONSTRAINT ck_captura_auditoria_tela CHECK (tela ~ '^[a-z]+\.[a-z0-9_]+(\.[a-z0-9_]+)*$'
    AND pg_catalog.split_part(tela, '.', 1) IN ('adm', 'atv', 'med', 'sol', 'os', 'pre', 'chk', 'est', 'ter', 'srv',
                                                'prj', 'ctr', 'loc', 'fro', 'con', 'ind', 'int', 'ia', 'mob'))
);
CALL app.preparar_tabela('app.captura_auditoria', true, 'tela');
-- +goose StatementEnd

-- ============================================================ 7. selo_auditoria

-- +goose StatementBegin
-- hash_dia (seção 9.4): sha256 de
--   "tenaz-selo-v1|<empresa_id>|<AAAA-MM-DD>|<quantidade>|" (UTF-8)
--   seguido dos hashes dos eventos do dia (32 bytes cada), em ordem (ocorrido_em, id).
-- O cabeçalho amarra o selo à empresa, ao dia e à contagem. Dia sem evento tem
-- selo com quantidade 0: sem isso, um evento inserido num dia vazio já passado
-- não seria notado.
CREATE OR REPLACE FUNCTION app.hash_dia_auditoria(p_empresa uuid, p_dia date, p_quantidade bigint, p_hashes bytea)
  RETURNS bytea
  LANGUAGE sql IMMUTABLE PARALLEL SAFE
AS $$
SELECT pg_catalog.sha256(
  pg_catalog.convert_to('tenaz-selo-v1|' || p_empresa::text || '|' || pg_catalog.to_char(p_dia, 'YYYY-MM-DD')
                        || '|' || p_quantidade::text || '|', 'UTF8') || p_hashes)
$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TABLE app.selo_auditoria (
  empresa_id      uuid        NOT NULL,
  dia             date        NOT NULL,                 -- dia UTC
  quantidade      bigint      NOT NULL,                 -- eventos do dia
  hash_dia        bytea       NOT NULL,                 -- app.hash_dia_auditoria()
  hash_anterior   bytea,                                -- hash_encadeado do dia anterior; nulo só no primeiro selo
  hash_encadeado  bytea       NOT NULL,                 -- sha256(coalesce(hash_anterior, '') || hash_dia)
  calculado_em    timestamptz NOT NULL DEFAULT now(),
  chave_objeto    text        NOT NULL,                 -- selos/empresas/<empresa_id>/<AAAA-MM-DD>.selo
  CONSTRAINT pk_selo_auditoria PRIMARY KEY (empresa_id, dia),
  CONSTRAINT ck_selo_auditoria_hash CHECK (pg_catalog.octet_length(hash_dia) = 32
    AND pg_catalog.octet_length(hash_encadeado) = 32
    AND (hash_anterior IS NULL OR pg_catalog.octet_length(hash_anterior) = 32)),
  CONSTRAINT ck_selo_auditoria_quantidade CHECK (quantidade >= 0)
);
-- +goose StatementEnd

-- +goose StatementBegin
-- BEFORE INSERT: o chamador diz só (empresa_id, dia); o resto o banco calcula e
-- sobrescreve, para que nem o job nem um defeito dele forjem um selo. Roda com o
-- papel de quem grava (tenaz_job), que lê a trilha da empresa pelo RLS.
--
--   - só se sela o dia que já passou, com uma hora de folga depois da meia-noite
--     UTC: ocorrido_em é o início da transação (seção 9.4);
--   - o selo é sempre o do dia seguinte ao último da empresa; o primeiro é o do
--     dia do primeiro evento: a cadeia não tem buraco;
--   - hash_anterior é o hash_encadeado do selo anterior, lido aqui.
CREATE OR REPLACE FUNCTION app.tg_selo_auditoria() RETURNS trigger
  LANGUAGE plpgsql
AS $$
DECLARE
  v_ultimo   app.selo_auditoria%ROWTYPE;
  v_limite   date := ((pg_catalog.now() AT TIME ZONE 'UTC') - interval '1 hour')::date;
  v_primeiro date;
  v_inicio   timestamptz := NEW.dia::timestamp AT TIME ZONE 'UTC';
  v_fim      timestamptz := (NEW.dia + 1)::timestamp AT TIME ZONE 'UTC';
  v_n        bigint;
  v_hashes   bytea;
BEGIN
  IF NEW.dia >= v_limite THEN
    RAISE EXCEPTION 'o dia % ainda não pode ser selado', NEW.dia USING ERRCODE = '22023';
  END IF;

  SELECT * INTO v_ultimo FROM app.selo_auditoria s
  WHERE s.empresa_id = NEW.empresa_id ORDER BY s.dia DESC LIMIT 1;
  IF FOUND THEN
    IF NEW.dia <> v_ultimo.dia + 1 THEN
      RAISE EXCEPTION 'o próximo selo é o do dia %, e não o de %', v_ultimo.dia + 1, NEW.dia USING ERRCODE = '22023';
    END IF;
    NEW.hash_anterior := v_ultimo.hash_encadeado;
  ELSE
    SELECT (min(e.ocorrido_em) AT TIME ZONE 'UTC')::date INTO v_primeiro
    FROM app.evento_auditoria e WHERE e.empresa_id = NEW.empresa_id;
    IF v_primeiro IS DISTINCT FROM NEW.dia THEN
      RAISE EXCEPTION 'o primeiro selo é o do dia do primeiro evento (%), e não o de %', v_primeiro, NEW.dia
        USING ERRCODE = '22023';
    END IF;
    NEW.hash_anterior := NULL;
  END IF;

  SELECT count(*), coalesce(string_agg(e.hash, ''::bytea ORDER BY e.ocorrido_em, e.id), ''::bytea)
    INTO v_n, v_hashes
  FROM app.evento_auditoria e
  WHERE e.empresa_id = NEW.empresa_id AND e.ocorrido_em >= v_inicio AND e.ocorrido_em < v_fim;

  NEW.quantidade     := v_n;
  NEW.hash_dia       := app.hash_dia_auditoria(NEW.empresa_id, NEW.dia, v_n, v_hashes);
  NEW.hash_encadeado := sha256(coalesce(NEW.hash_anterior, ''::bytea) || NEW.hash_dia);
  NEW.calculado_em   := now();
  NEW.chave_objeto   := 'selos/empresas/' || NEW.empresa_id::text || '/' || to_char(NEW.dia, 'YYYY-MM-DD') || '.selo';
  RETURN NEW;
END
$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_selo_calcular BEFORE INSERT ON app.selo_auditoria
  FOR EACH ROW EXECUTE FUNCTION app.tg_selo_auditoria();
-- Imutável como a trilha (seção 9.4): UPDATE e DELETE recusados até para o dono;
-- TRUNCATE ignora RLS e trigger de linha, então tem a trava própria.
CREATE TRIGGER trg_somente_insercao BEFORE UPDATE OR DELETE ON app.selo_auditoria
  FOR EACH ROW EXECUTE FUNCTION app.tg_somente_insercao();
CREATE TRIGGER trg_somente_insercao_truncate BEFORE TRUNCATE ON app.selo_auditoria
  FOR EACH STATEMENT EXECUTE FUNCTION app.tg_somente_insercao();

ALTER TABLE app.selo_auditoria ENABLE ROW LEVEL SECURITY;
ALTER TABLE app.selo_auditoria FORCE  ROW LEVEL SECURITY;

-- Policies por comando, como em evento_auditoria: nenhuma de UPDATE nem de DELETE.
CREATE POLICY selo_leitura ON app.selo_auditoria
  FOR SELECT TO tenaz_app, tenaz_job, tenaz_leitura
  USING (empresa_id = app.empresa_atual());
-- só o job do selo grava; a tela só lê
CREATE POLICY selo_insercao ON app.selo_auditoria
  FOR INSERT TO tenaz_job
  WITH CHECK (empresa_id = app.empresa_atual());
-- prestador, locatário e link não leem selo (seção 2.9.4)
CREATE POLICY escopo_externo ON app.selo_auditoria AS RESTRICTIVE
  FOR ALL TO tenaz_app, tenaz_job, tenaz_leitura
  USING (app.escopo_interno()) WITH CHECK (app.escopo_interno());

REVOKE ALL ON app.selo_auditoria FROM PUBLIC, tenaz_app, tenaz_job, tenaz_leitura;
GRANT SELECT ON app.selo_auditoria TO tenaz_app, tenaz_leitura;
GRANT SELECT, INSERT ON app.selo_auditoria TO tenaz_job;
-- +goose StatementEnd

-- +goose StatementBegin
-- Job diário do selo, por empresa: sela, em ordem, todos os dias pendentes e
-- devolve os dias selados, para o Go copiar cada selo ao armazenamento imutável.
-- Idempotente: sem dia pendente, devolve vazio. Só o tenaz_job executa.
CREATE OR REPLACE FUNCTION app.selar_dias_pendentes() RETURNS SETOF date
  LANGUAGE plpgsql VOLATILE
AS $$
DECLARE
  v_empresa uuid := app.empresa_atual();
  v_limite  date := ((pg_catalog.now() AT TIME ZONE 'UTC') - interval '1 hour')::date;
  v_dia     date;
BEGIN
  IF v_empresa IS NULL THEN
    RAISE EXCEPTION 'selar_dias_pendentes exige o contexto de uma empresa' USING ERRCODE = '22023';
  END IF;
  SELECT max(s.dia) + 1 INTO v_dia FROM app.selo_auditoria s WHERE s.empresa_id = v_empresa;
  IF v_dia IS NULL THEN
    SELECT (min(e.ocorrido_em) AT TIME ZONE 'UTC')::date INTO v_dia
    FROM app.evento_auditoria e WHERE e.empresa_id = v_empresa;
  END IF;
  WHILE v_dia IS NOT NULL AND v_dia < v_limite LOOP
    INSERT INTO app.selo_auditoria (empresa_id, dia) VALUES (v_empresa, v_dia);
    RETURN NEXT v_dia;
    v_dia := v_dia + 1;
  END LOOP;
  RETURN;
END
$$;
REVOKE ALL ON FUNCTION app.selar_dias_pendentes() FROM PUBLIC, tenaz_app, tenaz_leitura;
GRANT EXECUTE ON FUNCTION app.selar_dias_pendentes() TO tenaz_job;
-- +goose StatementEnd

-- +goose StatementBegin
-- Verificação de integridade (RF-ADM-011), na empresa da transação. Por dia:
--   integro      hash de cada evento confere, a contagem e o hash_dia conferem, o
--                selo confere e a cadeia confere;
--   violado      alguma conferência falhou; motivo diz qual e evento_id aponta o
--                primeiro evento alterado, quando há como apontar;
--   nao_selado   o dia ainda não tem selo (o estado nunca é gravado: é calculado aqui).
-- Motivos: evento_alterado, eventos_divergentes (evento apagado ou incluído fora
-- de ordem), selo_adulterado, cadeia_quebrada, selo_ausente,
-- evento_anterior_ao_primeiro_selo. Não olha a cópia no armazenamento imutável:
-- essa comparação é do Go. Até 400 dias por chamada.
CREATE OR REPLACE FUNCTION app.verificar_integridade_auditoria(p_de date, p_ate date)
  RETURNS TABLE (dia date, estado text, motivo text, evento_id uuid)
  LANGUAGE plpgsql STABLE
AS $$
DECLARE
  v_empresa  uuid := app.empresa_atual();
  v_dia      date;
  v_s        app.selo_auditoria%ROWTYPE;
  v_prev     app.selo_auditoria%ROWTYPE;
  v_inicio   timestamptz;
  v_fim      timestamptz;
  v_n        bigint;
  v_hashes   bytea;
  v_evento   uuid;
  v_primeiro date;
BEGIN
  IF v_empresa IS NULL THEN
    RAISE EXCEPTION 'verificar_integridade_auditoria exige o contexto de uma empresa' USING ERRCODE = '22023';
  END IF;
  IF p_de IS NULL OR p_ate IS NULL OR p_ate < p_de OR p_ate - p_de > 400 THEN
    RAISE EXCEPTION 'período inválido: de % até % (no máximo 400 dias)', p_de, p_ate USING ERRCODE = '22023';
  END IF;

  FOR v_dia IN SELECT g::date FROM generate_series(p_de::timestamp, p_ate::timestamp, interval '1 day') g LOOP
    dia := v_dia; estado := NULL; motivo := NULL; evento_id := NULL;
    v_inicio := v_dia::timestamp AT TIME ZONE 'UTC';
    v_fim    := (v_dia + 1)::timestamp AT TIME ZONE 'UTC';

    SELECT * INTO v_s FROM app.selo_auditoria s WHERE s.empresa_id = v_empresa AND s.dia = v_dia;
    IF NOT FOUND THEN
      SELECT min(s.dia) INTO v_primeiro FROM app.selo_auditoria s WHERE s.empresa_id = v_empresa;
      IF v_primeiro IS NULL OR v_dia > v_primeiro THEN
        -- sem selo nenhum, ou depois do primeiro: se há selo posterior, o dia foi apagado
        IF v_primeiro IS NOT NULL AND EXISTS (SELECT 1 FROM app.selo_auditoria s
                                              WHERE s.empresa_id = v_empresa AND s.dia > v_dia) THEN
          estado := 'violado'; motivo := 'selo_ausente';
        ELSE
          estado := 'nao_selado'; motivo := 'ainda_nao_selado';
        END IF;
      ELSE
        -- antes do primeiro selo: só pode haver evento se alguém o incluiu depois
        IF EXISTS (SELECT 1 FROM app.evento_auditoria e
                   WHERE e.empresa_id = v_empresa AND e.ocorrido_em >= v_inicio AND e.ocorrido_em < v_fim) THEN
          estado := 'violado'; motivo := 'evento_anterior_ao_primeiro_selo';
        ELSE
          estado := 'nao_selado'; motivo := 'anterior_ao_primeiro_selo';
        END IF;
      END IF;
      RETURN NEXT;
      CONTINUE;
    END IF;

    -- 1. o hash gravado em cada evento é o do conteúdo atual dele
    SELECT e.id INTO v_evento FROM app.evento_auditoria e
    WHERE e.empresa_id = v_empresa AND e.ocorrido_em >= v_inicio AND e.ocorrido_em < v_fim
      AND app.hash_evento_auditoria(e) IS DISTINCT FROM e.hash
    ORDER BY e.ocorrido_em, e.id LIMIT 1;
    IF FOUND THEN
      estado := 'violado'; motivo := 'evento_alterado'; evento_id := v_evento;
      RETURN NEXT; CONTINUE;
    END IF;

    -- 2. os eventos do dia são os do selo: nenhum apagado, nenhum incluído depois
    SELECT count(*), coalesce(string_agg(e.hash, ''::bytea ORDER BY e.ocorrido_em, e.id), ''::bytea)
      INTO v_n, v_hashes
    FROM app.evento_auditoria e
    WHERE e.empresa_id = v_empresa AND e.ocorrido_em >= v_inicio AND e.ocorrido_em < v_fim;
    IF v_n <> v_s.quantidade OR app.hash_dia_auditoria(v_empresa, v_dia, v_n, v_hashes) <> v_s.hash_dia THEN
      estado := 'violado'; motivo := 'eventos_divergentes';
      RETURN NEXT; CONTINUE;
    END IF;

    -- 3. o próprio selo
    IF sha256(coalesce(v_s.hash_anterior, ''::bytea) || v_s.hash_dia) <> v_s.hash_encadeado THEN
      estado := 'violado'; motivo := 'selo_adulterado';
      RETURN NEXT; CONTINUE;
    END IF;

    -- 4. a cadeia: encaixa no selo do dia anterior (ou é o primeiro)
    SELECT * INTO v_prev FROM app.selo_auditoria s WHERE s.empresa_id = v_empresa AND s.dia = v_dia - 1;
    IF (FOUND AND v_s.hash_anterior IS DISTINCT FROM v_prev.hash_encadeado)
       OR (NOT FOUND AND v_s.hash_anterior IS NOT NULL) THEN
      estado := 'violado'; motivo := 'cadeia_quebrada';
      RETURN NEXT; CONTINUE;
    END IF;

    estado := 'integro';
    RETURN NEXT;
  END LOOP;
  RETURN;
END
$$;
-- +goose StatementEnd

-- +goose StatementBegin
-- Volta ao dono comum no fim: o goose registra a versão na tabela de controle,
-- que pertence a tenaz_dono, com o papel corrente da conexão.
SET ROLE tenaz_dono;
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
SET ROLE tenaz_dono;
SET LOCAL lock_timeout = '3s';
SET LOCAL statement_timeout = '60s';

DROP FUNCTION IF EXISTS app.verificar_integridade_auditoria(date, date);
DROP FUNCTION IF EXISTS app.selar_dias_pendentes();
DROP TABLE IF EXISTS app.selo_auditoria;
DROP FUNCTION IF EXISTS app.tg_selo_auditoria();
DROP FUNCTION IF EXISTS app.hash_dia_auditoria(uuid, date, bigint, bytea);
DROP TABLE IF EXISTS app.captura_auditoria;

-- tg_auditar volta à versão da 00003: 'exclusao' para a remoção física e
-- alteração de excluido_em como 'alteracao'.
CREATE OR REPLACE FUNCTION app.tg_auditar() RETURNS trigger
  LANGUAGE plpgsql
AS $$
DECLARE
  v_coluna_chave text   := TG_ARGV[0];
  v_ignorar      text[] := ARRAY['versao', 'alterado_xid', 'atualizado_em', 'atualizado_por'];
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
      RETURN NULL;
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
    (coalesce(v_linha ->> 'empresa_id', v_linha ->> 'id')::uuid,
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
     app.conexao_atual(),
     nullif(current_setting('app.ocorrido_no_dispositivo_em', true), '')::timestamptz,
     app.escopo_atual());
  RETURN NULL;
END
$$;

-- autenticação como na 00005
ALTER POLICY autenticacao_identidade ON app.evento_auditoria
  WITH CHECK (ator_tipo = 'usuario' AND entidade = 'identidade'
              AND acao IN ('login', 'login_falho', 'logout', 'alteracao')
              AND EXISTS (SELECT 1 FROM plataforma.acesso_empresa a
                          WHERE a.identidade_id = evento_auditoria.entidade_id
                            AND a.empresa_id = evento_auditoria.empresa_id
                            AND a.usuario_id = evento_auditoria.usuario_id));
REVOKE INSERT (usuario_nome, tela, resultado) ON app.evento_auditoria FROM tenaz_autenticacao;

DROP TRIGGER IF EXISTS trg_evento_preparar ON app.evento_auditoria;
DROP INDEX IF EXISTS app.ix_evento_auditoria_acao;
DROP INDEX IF EXISTS app.ix_evento_auditoria_modulo;

ALTER TABLE app.evento_auditoria
  DROP CONSTRAINT ck_evento_auditoria_acao,
  DROP CONSTRAINT ck_evento_auditoria_resultado,
  DROP CONSTRAINT ck_evento_auditoria_sensivel,
  DROP CONSTRAINT ck_evento_auditoria_modulo,
  DROP CONSTRAINT ck_evento_auditoria_tela,
  DROP CONSTRAINT ck_evento_auditoria_hash;

-- Eventos de valor novo voltam ao mais próximo dos antigos; o que o valor novo
-- carregava a mais (tela, módulo, hash) some com as colunas.
ALTER TABLE app.evento_auditoria DISABLE TRIGGER trg_somente_insercao;
RESET ROLE;
SET LOCAL statement_timeout = '15min';
UPDATE app.evento_auditoria e
SET acao = CASE
             WHEN e.acao = 'login' AND e.resultado = 'falha'    THEN 'login_falho'
             WHEN e.acao = 'visualizacao'                       THEN 'visualizacao_sensivel'
             WHEN e.acao IN ('extracao', 'impressao')           THEN 'exportacao'
             WHEN e.acao = 'delecao'                            THEN 'exclusao'
             ELSE e.acao
           END;
SET LOCAL statement_timeout = '60s';
SET ROLE tenaz_dono;
ALTER TABLE app.evento_auditoria ENABLE TRIGGER trg_somente_insercao;

ALTER TABLE app.evento_auditoria
  DROP COLUMN hash, DROP COLUMN detalhe, DROP COLUMN sensivel, DROP COLUMN resultado,
  DROP COLUMN modulo, DROP COLUMN tela, DROP COLUMN usuario_nome;
ALTER TABLE app.evento_auditoria
  ADD CONSTRAINT ck_evento_auditoria_acao CHECK (acao IN ('inclusao', 'alteracao', 'exclusao', 'login', 'login_falho',
    'logout', 'exportacao', 'visualizacao_sensivel', 'acesso_suporte', 'revisao', 'sincronizacao'));

DROP FUNCTION IF EXISTS app.tg_evento_auditoria_preparar();
DROP FUNCTION IF EXISTS app.hash_evento_auditoria(app.evento_auditoria);
DROP FUNCTION IF EXISTS app.json_valor(anyelement);

-- acesso_empresa e o trigger voltam à 00005
DROP TRIGGER trg_acesso_empresa ON app.usuario;
CREATE TRIGGER trg_acesso_empresa AFTER INSERT OR UPDATE OF identidade_id, situacao, excluido_em ON app.usuario
  FOR EACH ROW WHEN (NEW.identidade_id IS NOT NULL)
  EXECUTE FUNCTION plataforma.tg_acesso_empresa();

SET ROLE tenaz_credencial;
CREATE OR REPLACE FUNCTION plataforma.tg_acesso_empresa() RETURNS trigger
  LANGUAGE plpgsql SECURITY DEFINER
  SET search_path = pg_catalog, pg_temp
AS $f$
DECLARE
  v_habilitado boolean := NEW.situacao = 'ativo' AND NEW.excluido_em IS NULL;
BEGIN
  INSERT INTO plataforma.acesso_empresa AS a (identidade_id, empresa_id, usuario_id, habilitado, atualizado_em)
  VALUES (NEW.identidade_id, NEW.empresa_id, NEW.id, v_habilitado, pg_catalog.now())
  ON CONFLICT (identidade_id, empresa_id) DO UPDATE
    SET usuario_id    = EXCLUDED.usuario_id,
        habilitado    = EXCLUDED.habilitado,
        atualizado_em = EXCLUDED.atualizado_em;

  IF NOT v_habilitado THEN
    WITH revogadas AS (
      UPDATE plataforma.sessao s
      SET revogada_em = pg_catalog.now()
      WHERE s.empresa_id = NEW.empresa_id AND s.usuario_id = NEW.id AND s.revogada_em IS NULL
      RETURNING s.id
    )
    UPDATE plataforma.token_renovacao t
    SET revogado_em = pg_catalog.now()
    WHERE t.familia_id IN (SELECT r.id FROM revogadas r) AND t.revogado_em IS NULL;
  END IF;
  RETURN NULL;
END
$f$;
ALTER TABLE plataforma.acesso_empresa DROP COLUMN usuario_nome;
SET ROLE tenaz_dono;
-- +goose StatementEnd
