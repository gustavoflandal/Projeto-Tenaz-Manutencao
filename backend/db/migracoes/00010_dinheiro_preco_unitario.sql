-- Representação do dinheiro, primeira aplicação (ADR-013, Parte F, documento 5,
-- versão 0.9; documento 2, seção 1.4, versão 0.9).
--
-- usuario.custo_hora é preço unitário (por hora), não quantia: vai a
-- numeric(18,6), sem limite à escala da moeda, e passa a ser multiplicado pela
-- quantidade de horas com um único arredondamento (no apontamento, em E04).
-- Os valores gravados não mudam (numeric(12,2) cabe em numeric(18,6)).
-- usuario.moeda perde o DEFAULT 'BRL': a moeda é sempre explícita (o backend
-- grava desde a 00009), para que um INSERT esquecido falhe em vez de gravar BRL
-- numa empresa argentina ou paraguaia.
--
-- Os privilégios (GRANT por coluna a tenaz_app, tenaz_job, tenaz_leitura), os
-- gatilhos (trg_usuario_colunas_externo, trg_auditoria) e as políticas não
-- dependem do tipo e ficam como estão: ALTER COLUMN TYPE preserva ACL.
-- app.casas_moeda(char(3)) só nasce com a primeira coluna de quantia.
--
-- Moeda do usuário = moeda da empresa, no banco (achado do tenaz-isolamento): o
-- custo do apontamento multiplica custo_hora por horas na moeda da empresa, e
-- uma moeda errada daria custo errado em ordem de grandeza. UNIQUE (id, moeda)
-- em empresa e FK composta (empresa_id, moeda) em usuario, sem ON UPDATE: a
-- moeda da empresa é imutável (trg_empresa_pais_imutavel). É o padrão de toda
-- coluna de quantia com moeda (documento 2, seção 1.4). Os usuários existentes
-- receberam a moeda da empresa na 00009 (DEFAULT 'BRL' em empresas BRL), e o
-- ensaio com 3.210 usuários não achou violação.
--
-- Compatível com a versão anterior do código: o Go que lê custo_hora como
-- numeric continua lendo; o que grava a moeda explicitamente continua gravando.
-- Tabela pequena (uma linha por usuário), logo a reescrita da troca de escala é
-- curta; lock_timeout protege contra espera longa.

-- +goose Up

-- +goose StatementBegin
SET ROLE tenaz_dono;
SET LOCAL lock_timeout = '3s';
SET LOCAL statement_timeout = '60s';

-- A conferência roda como o migrador (sem RLS): com o RLS forçado, a própria
-- validação da FK não enxerga as linhas de empresa nenhuma e passaria em falso.
RESET ROLE;
DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM app.usuario u JOIN app.empresa e ON e.id = u.empresa_id WHERE u.moeda <> e.moeda) THEN
    RAISE EXCEPTION 'há usuário em moeda diferente da da sua empresa; corrija antes da 00010' USING ERRCODE = '23503';
  END IF;
END
$$;
SET ROLE tenaz_dono;

ALTER TABLE app.usuario
  ALTER COLUMN custo_hora TYPE numeric(18,6),
  ADD CONSTRAINT ck_usuario_custo_hora CHECK (custo_hora >= 0),
  ALTER COLUMN moeda DROP DEFAULT;

ALTER TABLE app.empresa ADD CONSTRAINT uq_empresa_id_moeda UNIQUE (id, moeda);
ALTER TABLE app.usuario ADD CONSTRAINT fk_usuario_moeda_empresa
  FOREIGN KEY (empresa_id, moeda) REFERENCES app.empresa (id, moeda);
RESET ROLE;
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
SET ROLE tenaz_dono;
SET LOCAL lock_timeout = '3s';
SET LOCAL statement_timeout = '60s';

-- A volta arredonda a 2 casas (o que a 00004 aceitava); falha se algum valor
-- passar de 9.999.999.999,99, o que o Up nunca produziu com dado da 00004.
ALTER TABLE app.usuario DROP CONSTRAINT fk_usuario_moeda_empresa;
ALTER TABLE app.empresa DROP CONSTRAINT uq_empresa_id_moeda;
ALTER TABLE app.usuario
  ALTER COLUMN moeda SET DEFAULT 'BRL',
  DROP CONSTRAINT ck_usuario_custo_hora,
  ALTER COLUMN custo_hora TYPE numeric(12,2);
RESET ROLE;
-- +goose StatementEnd
