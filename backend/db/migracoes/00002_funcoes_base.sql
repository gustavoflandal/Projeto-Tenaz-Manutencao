-- Funções de contexto, UUID v7 e o trigger das colunas padrão.
--
-- Fonte: docs/2_ARQUITETURA_BANCO_DADOS.md, seções 2.3, 2.9.2, 4.2 a 4.4.
--
-- Nenhuma tabela de negócio nasce aqui. Esta migração cria o chão em que elas
-- pisam: sem ela, uma tabela criada depois não tem como falhar fechada.

-- +goose Up

-- +goose StatementBegin
-- Todo objeto criado a partir daqui pertence a tenaz_dono (seção 4.1).
SET ROLE tenaz_dono;
-- +goose StatementEnd

-- Funções de contexto. Todas STABLE e em SQL puro, para o planejador expandir
-- dentro da consulta e usar o índice que começa por empresa_id.
--
-- O nullif(..., '') é essencial: depois de um COMMIT, uma variável definida com
-- is_local => true volta a string vazia, não a nulo. Sem ele, o cast falharia
-- com erro em vez de devolver zero linhas — e o erro em si já contaria ao
-- chamador que havia algo ali.
--
-- Nenhuma delas recebe coluna da linha como argumento, então nenhuma precisa
-- ser LEAKPROOF (marcar como LEAKPROOF exige superusuário e é risco por si).

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION app.empresa_atual() RETURNS uuid
  LANGUAGE sql STABLE PARALLEL SAFE
AS $$ SELECT nullif(current_setting('app.empresa_id', true), '')::uuid $$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION app.usuario_atual() RETURNS uuid
  LANGUAGE sql STABLE PARALLEL SAFE
AS $$ SELECT nullif(current_setting('app.usuario_id', true), '')::uuid $$;
-- +goose StatementEnd

-- +goose StatementBegin
-- 'interno', 'prestador', 'locatario' ou 'link'. Sem valor: NULL, e toda policy
-- de escopo falha fechada, inclusive para usuário interno.
CREATE OR REPLACE FUNCTION app.escopo_atual() RETURNS text
  LANGUAGE sql STABLE PARALLEL SAFE
AS $$ SELECT nullif(current_setting('app.escopo', true), '') $$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION app.escopo_interno() RETURNS boolean
  LANGUAGE sql STABLE PARALLEL SAFE
AS $$ SELECT app.escopo_atual() = 'interno' $$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION app.fornecedor_atual() RETURNS uuid
  LANGUAGE sql STABLE PARALLEL SAFE
AS $$ SELECT nullif(current_setting('app.fornecedor_id', true), '')::uuid $$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE OR REPLACE FUNCTION app.locatario_atual() RETURNS uuid
  LANGUAGE sql STABLE PARALLEL SAFE
AS $$ SELECT nullif(current_setting('app.locatario_id', true), '')::uuid $$;
-- +goose StatementEnd

-- +goose StatementBegin
-- sha256 do segredo do link de uso único, em hexadecimal
CREATE OR REPLACE FUNCTION app.link_hash_atual() RETURNS bytea
  LANGUAGE sql STABLE PARALLEL SAFE
AS $$ SELECT decode(nullif(current_setting('app.link_hash', true), ''), 'hex') $$;
-- +goose StatementEnd

-- +goose StatementBegin
-- permissão pessoa.saude.ler, que não está em nenhum perfil padrão
CREATE OR REPLACE FUNCTION app.pode_ler_saude() RETURNS boolean
  LANGUAGE sql STABLE PARALLEL SAFE
AS $$ SELECT coalesce(current_setting('app.saude', true), '') = 'sim' $$;
-- +goose StatementEnd

-- +goose StatementBegin
-- Conferência do escopo contra o cadastro, chamada pelo helper NaEmpresa logo
-- depois de fixar o contexto da transação.
--
-- Esta é a versão da fundação: valida o valor do escopo, exige empresa e exige
-- o segredo quando o escopo é de link. A conferência contra app.usuario
-- (prestador marcado como interno por engano, usuário que perdeu o vínculo com
-- o fornecedor durante a sessão) entra na migração de ADM, quando app.usuario
-- existir, por CREATE OR REPLACE — RF-ADM-004 e RF-ADM-005.
CREATE OR REPLACE FUNCTION app.conferir_escopo() RETURNS void
  LANGUAGE plpgsql STABLE
AS $$
BEGIN
  IF app.escopo_atual() IS NULL
     OR app.escopo_atual() NOT IN ('interno', 'prestador', 'locatario', 'link') THEN
    RAISE EXCEPTION 'escopo ausente ou inválido' USING ERRCODE = '42501';
  END IF;

  IF app.empresa_atual() IS NULL THEN
    RAISE EXCEPTION 'empresa ausente no contexto' USING ERRCODE = '42501';
  END IF;

  IF app.usuario_atual() IS NULL THEN
    IF app.escopo_atual() = 'link' AND app.link_hash_atual() IS NULL THEN
      RAISE EXCEPTION 'link sem segredo' USING ERRCODE = '42501';
    END IF;
    IF app.escopo_atual() IN ('prestador', 'locatario') THEN
      RAISE EXCEPTION 'escopo externo sem usuário' USING ERRCODE = '42501';
    END IF;
  END IF;
END
$$;
-- +goose StatementEnd

-- +goose StatementBegin
-- UUID v7 como rede de segurança no banco. O backend e o app geram o seu com a
-- biblioteca; o valor padrão existe para o caso de uma inserção direta.
-- No PostgreSQL 18 há uuidv7() nativo; nas versões 16 e 17, troca os 48 bits
-- iniciais de um UUID v4 pelo instante em milissegundos e ajusta os bits de
-- versão de 0100 para 0111.
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM pg_proc
    WHERE proname = 'uuidv7' AND pronamespace = 'pg_catalog'::regnamespace
  ) THEN
    EXECUTE $f$
      CREATE OR REPLACE FUNCTION app.uuid_v7() RETURNS uuid
        LANGUAGE sql VOLATILE PARALLEL SAFE
      AS 'SELECT uuidv7()'
    $f$;
  ELSE
    EXECUTE $f$
      CREATE OR REPLACE FUNCTION app.uuid_v7() RETURNS uuid
        LANGUAGE sql VOLATILE PARALLEL SAFE
      AS $corpo$
        SELECT encode(
                 set_bit(set_bit(
                   overlay(uuid_send(gen_random_uuid())
                           PLACING substring(int8send((extract(epoch FROM clock_timestamp()) * 1000)::bigint) FROM 3)
                           FROM 1 FOR 6),
                 52, 1), 53, 1),
               'hex')::uuid
      $corpo$
    $f$;
  END IF;
END
$$;
-- +goose StatementEnd

-- +goose StatementBegin
-- Trigger das colunas padrão (seção 4.4). Mantém versao, alterado_xid, criado e
-- atualizado, e torna empresa_id imutável: uma linha nunca muda de empresa, nem
-- por erro de código.
CREATE OR REPLACE FUNCTION app.tg_padrao() RETURNS trigger
  LANGUAGE plpgsql
AS $$
BEGIN
  IF TG_OP = 'INSERT' THEN
    NEW.versao         := 1;
    NEW.criado_em      := now();
    NEW.criado_por     := coalesce(NEW.criado_por, app.usuario_atual());
    NEW.atualizado_em  := NEW.criado_em;
    NEW.atualizado_por := NEW.criado_por;
  ELSE
    IF NEW.empresa_id IS DISTINCT FROM OLD.empresa_id THEN
      RAISE EXCEPTION 'empresa_id é imutável' USING ERRCODE = '42501';
    END IF;
    NEW.versao         := OLD.versao + 1;
    NEW.criado_em      := OLD.criado_em;
    NEW.criado_por     := OLD.criado_por;
    NEW.atualizado_em  := now();
    NEW.atualizado_por := app.usuario_atual();
  END IF;
  NEW.alterado_xid := pg_current_xact_id();
  RETURN NEW;
END
$$;
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
SET ROLE tenaz_dono;
DROP FUNCTION IF EXISTS app.tg_padrao();
DROP FUNCTION IF EXISTS app.uuid_v7();
DROP FUNCTION IF EXISTS app.conferir_escopo();
DROP FUNCTION IF EXISTS app.pode_ler_saude();
DROP FUNCTION IF EXISTS app.link_hash_atual();
DROP FUNCTION IF EXISTS app.locatario_atual();
DROP FUNCTION IF EXISTS app.fornecedor_atual();
DROP FUNCTION IF EXISTS app.escopo_interno();
DROP FUNCTION IF EXISTS app.escopo_atual();
DROP FUNCTION IF EXISTS app.usuario_atual();
DROP FUNCTION IF EXISTS app.empresa_atual();
-- +goose StatementEnd
