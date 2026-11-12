-- Bootstrap do cluster: extensões, papéis, esquemas e permissões padrão.
--
-- Fonte: docs/2_ARQUITETURA_BANCO_DADOS.md, seções 2.2 e 4.1.
--
-- Esta é a única migração que roda como o papel migrador, porque cria os
-- próprios papéis. Da 00002 em diante, toda migração começa com
-- `SET ROLE tenaz_dono`, para que todo objeto tenha o mesmo dono. Isso não é
-- estética: `ALTER DEFAULT PRIVILEGES FOR ROLE tenaz_dono` só alcança objetos
-- criados POR tenaz_dono, e `FORCE ROW LEVEL SECURITY` não vale para o dono da
-- tabela — uma tabela de dono superusuário contornaria o RLS.
--
-- Os papéis têm guarda de existência: em provedor gerenciado eles vêm da
-- infraestrutura como código, e a guarda deixa esta parte sem efeito.

-- +goose Up

-- +goose StatementBegin
-- Extensões que acompanham o PostgreSQL padrão.
CREATE EXTENSION IF NOT EXISTS citext;      -- e-mail sem diferença de maiúsculas
CREATE EXTENSION IF NOT EXISTS btree_gist;  -- EXCLUDE com uuid + intervalo
CREATE EXTENSION IF NOT EXISTS btree_gin;   -- GIN composto (empresa_id, texto)
CREATE EXTENSION IF NOT EXISTS pg_trgm;     -- busca por trecho em TAG, nome, descrição
-- pg_partman não entra aqui: não acompanha a imagem oficial do PostgreSQL e só
-- é necessária na migração que criar a primeira tabela particionada
-- (leitura_medidor). Ver a pendência no registro da etapa.
-- +goose StatementEnd

-- +goose StatementBegin
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'tenaz_dono') THEN
    CREATE ROLE tenaz_dono NOLOGIN;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'tenaz_app') THEN
    CREATE ROLE tenaz_app LOGIN NOBYPASSRLS NOINHERIT;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'tenaz_job') THEN
    CREATE ROLE tenaz_job LOGIN NOBYPASSRLS NOINHERIT;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'tenaz_leitura') THEN
    CREATE ROLE tenaz_leitura LOGIN NOBYPASSRLS NOINHERIT;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'tenaz_manutencao') THEN
    CREATE ROLE tenaz_manutencao LOGIN NOBYPASSRLS IN ROLE tenaz_dono;
  END IF;
  -- Papel do sincronizador entre contas conectadas (SRV). Fica fora dos
  -- privilégios padrão de propósito: cada tabela que ele toca recebe GRANT por
  -- coluna na própria migração, conferido pelo teste de catálogo (risco R5).
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'tenaz_sinc') THEN
    CREATE ROLE tenaz_sinc NOLOGIN NOBYPASSRLS;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'tenaz_sinc_conexao') THEN
    CREATE ROLE tenaz_sinc_conexao LOGIN NOBYPASSRLS IN ROLE tenaz_sinc;
  END IF;

  -- Quem aplica a migração precisa poder assumir o papel do dono.
  IF NOT pg_has_role(current_user, 'tenaz_dono', 'MEMBER') THEN
    EXECUTE format('GRANT tenaz_dono TO %I', current_user);
  END IF;
END
$$;
-- +goose StatementEnd

-- +goose StatementBegin
-- app: dados de empresa, sempre com RLS.
-- plataforma: dados globais da Tenaz, sem dado de negócio.
CREATE SCHEMA IF NOT EXISTS app AUTHORIZATION tenaz_dono;

-- O esquema plataforma já existe neste ponto: o executor das migrações o cria
-- antes, porque a tabela de controle do goose mora nele. Aqui ele só passa
-- para o dono correto.
CREATE SCHEMA IF NOT EXISTS plataforma AUTHORIZATION tenaz_dono;
ALTER SCHEMA plataforma OWNER TO tenaz_dono;

REVOKE ALL ON SCHEMA public FROM PUBLIC;
GRANT USAGE ON SCHEMA app TO tenaz_app, tenaz_job, tenaz_leitura;

-- A tabela de controle do goose foi criada pelo executor antes desta migração
-- e pertence a ele. Ela passa para tenaz_dono aqui, uma única vez, para que as
-- migrações seguintes possam registrar a própria versão depois do
-- `SET ROLE tenaz_dono`. Sem isto, cada migração precisaria lembrar de um
-- `RESET ROLE` no fim — disciplina que uma hora falha.
ALTER TABLE IF EXISTS plataforma.migracao OWNER TO tenaz_dono;
-- +goose StatementEnd

-- +goose StatementBegin
-- Privilégios padrão: valem para objetos criados por tenaz_dono daqui em
-- diante, que é a razão de toda migração assumir esse papel.
ALTER DEFAULT PRIVILEGES FOR ROLE tenaz_dono IN SCHEMA app
  GRANT SELECT, INSERT, UPDATE, DELETE ON TABLES TO tenaz_app, tenaz_job;
ALTER DEFAULT PRIVILEGES FOR ROLE tenaz_dono IN SCHEMA app
  GRANT SELECT ON TABLES TO tenaz_leitura;
ALTER DEFAULT PRIVILEGES FOR ROLE tenaz_dono IN SCHEMA app
  GRANT EXECUTE ON FUNCTIONS TO tenaz_app, tenaz_job, tenaz_leitura;
-- TRUNCATE nunca é concedido: TRUNCATE ignora RLS.
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
ALTER DEFAULT PRIVILEGES FOR ROLE tenaz_dono IN SCHEMA app
  REVOKE SELECT, INSERT, UPDATE, DELETE ON TABLES FROM tenaz_app, tenaz_job;
ALTER DEFAULT PRIVILEGES FOR ROLE tenaz_dono IN SCHEMA app
  REVOKE SELECT ON TABLES FROM tenaz_leitura;
ALTER DEFAULT PRIVILEGES FOR ROLE tenaz_dono IN SCHEMA app
  REVOKE EXECUTE ON FUNCTIONS FROM tenaz_app, tenaz_job, tenaz_leitura;
DROP SCHEMA IF EXISTS app;
-- O esquema plataforma não é removido: guarda a tabela de controle das
-- próprias migrações. Papéis também não: são objetos do cluster e podem ter
-- vindo da infraestrutura, fora desta migração.
-- +goose StatementEnd
