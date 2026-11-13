-- A Tenaz não guarda dado de saúde.
--
-- Decisão do dono do produto, 13/11/2026 (sumário 0.7, seção 3.6, linha
-- "Dados de saúde"; documento 1, versão 0.6): ASO, exame toxicológico e
-- qualquer outro dado de saúde da pessoa saem do escopo. A permissão
-- pessoa.saude.ler deixa de existir, e com ela a única exceção ao formato
-- <módulo>.<ação>. Fonte: docs/2_ARQUITETURA_BANCO_DADOS.md, versão 0.6.
--
-- Migração nova, e não edição: 00001 a 00005 já estão em main (commit
-- 310dee9), e migração aplicada não se edita (seção 14.2).
--
-- Entra no MESMO commit que a mudança do Go: o helper NaEmpresa deixa de
-- definir app.saude e de carregar PodeLerSaude, o catálogo perde
-- pessoa.saude.ler, e os testes que concedem essa permissão (inclusive as
-- fixtures da suíte de vazamento) passam a usar outra. Antes disso, o CHECK de
-- formato recriado aqui recusa o código, e eles falham.

-- +goose Up

-- +goose StatementBegin
-- Limpeza de dados, como quem aplica a migração: o dono não passa pelo RLS
-- forçado das tabelas de app, e o migrador tem BYPASSRLS. O RESET é necessário
-- porque o SET ROLE da migração anterior sobrevive na conexão do goose.
--
-- Nenhuma conta usa a permissão ainda (nada em homologação nem em produção);
-- isto existe para que o CHECK abaixo valide em qualquer banco. A concessão
-- avulsa é apagada, e não só excluída logicamente, porque o CHECK vale também
-- para a linha excluída; o trg_auditoria registra a exclusão na trilha.
RESET ROLE;
DELETE FROM app.usuario_permissao WHERE permissao = 'pessoa.saude.ler';
UPDATE app.perfil SET permissoes = array_remove(permissoes, 'pessoa.saude.ler')
WHERE 'pessoa.saude.ler' = ANY (permissoes);
-- +goose StatementEnd

-- +goose StatementBegin
SET ROLE tenaz_dono;
SET LOCAL lock_timeout = '3s';
SET LOCAL statement_timeout = '60s';
-- +goose StatementEnd

-- +goose StatementBegin
-- Sem a permissão, o CHECK que a mantinha fora dos perfis padrão perde o objeto.
ALTER TABLE app.perfil DROP CONSTRAINT ck_perfil_saude;

-- Formato sem exceção: toda permissão é <módulo>.<ação>. NOT VALID + VALIDATE,
-- como toda troca de CHECK (seção 14.3).
ALTER TABLE app.usuario_permissao DROP CONSTRAINT ck_usuario_permissao_formato;
ALTER TABLE app.usuario_permissao ADD CONSTRAINT ck_usuario_permissao_formato
  CHECK (permissao ~ '^[a-z]{2,3}\.[a-z][a-z0-9_]*[a-z0-9]$') NOT VALID;
ALTER TABLE app.usuario_permissao VALIDATE CONSTRAINT ck_usuario_permissao_formato;
COMMENT ON COLUMN app.usuario_permissao.permissao IS '<módulo>.<ação>, do catálogo de RF-ADM-005, sem exceção';
-- +goose StatementEnd

-- +goose StatementBegin
-- Lista de vedadas a usuário externo sem pessoa.saude.ler. Mudar a lista é
-- CREATE OR REPLACE mais a recriação do CHECK que a usa (seção 5.1): o CHECK
-- não reconfere sozinho as linhas antigas.
CREATE OR REPLACE FUNCTION app.permissoes_vedadas_externo() RETURNS text[]
  LANGUAGE sql IMMUTABLE PARALLEL SAFE
AS $$
  SELECT ARRAY[
    'os.encerrar',
    'adm.gerenciar_regras', 'adm.gerenciar_configuracao', 'atv.colocar_em_operacao',
    'prj.gerenciar', 'prj.dispensar_etapa', 'prj.ajustar_medicao',
    'ctr.gerenciar', 'ctr.fechar_medicao', 'ctr.reabrir_medicao',
    'srv.gerenciar_conexao',
    -- grupo Admin
    'adm.gerenciar_empresa', 'adm.gerenciar_plano', 'adm.gerenciar_usuarios', 'adm.gerenciar_perfis',
    'adm.gerenciar_parametros', 'adm.ler_auditoria', 'adm.gerenciar_sso', 'adm.exportar',
    'con.gerenciar_tipos_documento', 'int.gerenciar_tokens', 'int.gerenciar_webhooks',
    'int.gerenciar_conexoes_telematica', 'int.gerenciar_erp', 'ia.gerenciar'
  ]::text[]
$$;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE app.perfil DROP CONSTRAINT ck_perfil_externo_vedadas;
ALTER TABLE app.perfil ADD CONSTRAINT ck_perfil_externo_vedadas
  CHECK (coalesce(perfil_origem, codigo) NOT IN ('prestador', 'locatario')
         OR NOT (permissoes && app.permissoes_vedadas_externo())) NOT VALID;
ALTER TABLE app.perfil VALIDATE CONSTRAINT ck_perfil_externo_vedadas;
-- +goose StatementEnd

-- +goose StatementBegin
-- A função de contexto da permissão. Nenhuma policy a usa: a policy de saúde
-- de documento_vencimento (seção 5.10) nunca chegou a ser criada. A variável
-- app.saude, se o Go ainda a definir, fica sem efeito.
DROP FUNCTION app.pode_ler_saude();
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
-- Devolve a estrutura da 00004 e da 00002. As concessões apagadas no Up não
-- voltam: estão só na trilha de auditoria.
SET ROLE tenaz_dono;

CREATE OR REPLACE FUNCTION app.pode_ler_saude() RETURNS boolean
  LANGUAGE sql STABLE PARALLEL SAFE
AS $$ SELECT coalesce(current_setting('app.saude', true), '') = 'sim' $$;

CREATE OR REPLACE FUNCTION app.permissoes_vedadas_externo() RETURNS text[]
  LANGUAGE sql IMMUTABLE PARALLEL SAFE
AS $$
  SELECT ARRAY[
    'os.encerrar',
    'adm.gerenciar_regras', 'adm.gerenciar_configuracao', 'atv.colocar_em_operacao',
    'prj.gerenciar', 'prj.dispensar_etapa', 'prj.ajustar_medicao',
    'ctr.gerenciar', 'ctr.fechar_medicao', 'ctr.reabrir_medicao',
    'srv.gerenciar_conexao',
    'pessoa.saude.ler',
    -- grupo Admin
    'adm.gerenciar_empresa', 'adm.gerenciar_plano', 'adm.gerenciar_usuarios', 'adm.gerenciar_perfis',
    'adm.gerenciar_parametros', 'adm.ler_auditoria', 'adm.gerenciar_sso', 'adm.exportar',
    'con.gerenciar_tipos_documento', 'int.gerenciar_tokens', 'int.gerenciar_webhooks',
    'int.gerenciar_conexoes_telematica', 'int.gerenciar_erp', 'ia.gerenciar'
  ]::text[]
$$;

ALTER TABLE app.perfil DROP CONSTRAINT ck_perfil_externo_vedadas;
ALTER TABLE app.perfil ADD CONSTRAINT ck_perfil_externo_vedadas
  CHECK (coalesce(perfil_origem, codigo) NOT IN ('prestador', 'locatario')
         OR NOT (permissoes && app.permissoes_vedadas_externo()));

ALTER TABLE app.usuario_permissao DROP CONSTRAINT ck_usuario_permissao_formato;
ALTER TABLE app.usuario_permissao ADD CONSTRAINT ck_usuario_permissao_formato
  CHECK (permissao ~ '^[a-z]{2,3}\.[a-z][a-z0-9_]*[a-z0-9]$' OR permissao = 'pessoa.saude.ler');
COMMENT ON COLUMN app.usuario_permissao.permissao IS NULL;

ALTER TABLE app.perfil ADD CONSTRAINT ck_perfil_saude
  CHECK (NOT padrao OR NOT ('pessoa.saude.ler' = ANY (permissoes)));
-- +goose StatementEnd
