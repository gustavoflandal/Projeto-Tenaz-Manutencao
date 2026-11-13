-- País, identificação, moeda e idioma (ADR-013, aceito em 13/11/2026; sumário
-- 0.8, seção 3.6, linha "País, identificação, moeda e idioma").
--
-- Fonte: docs/2_ARQUITETURA_BANCO_DADOS.md, versão 0.8, seções 1.4, 5.1 e 8.2.
-- O ADR citava 00008, número que já é o da trava da trilha.
--
-- O que faz:
--   - empresa: pais (BR, AR, PY), moeda (uma por empresa, dada pelo país),
--     idioma_padrao, e o par identificacao_tipo + identificacao no lugar de cnpj;
--   - unidade: o par no lugar de cnpj (a unidade não tem país próprio: Q1);
--   - pessoa: o par no lugar de cpf, aceitando os tipos de qualquer país;
--   - usuario: idioma, e moeda ao lado de custo_hora;
--   - pais e moeda da empresa viram imutáveis (operação assistida do suporte,
--     que é migração, e não a API).
--
-- Validação em duas camadas: o banco confere a FORMA de cada tipo
-- (app.identificacao_valida); o dígito verificador é do Go e do TypeScript.
-- CNPJ aceita letras nas 12 primeiras posições (CNPJ alfanumérico). CUIT e RUC só
-- pela forma, até a P13.
--
-- Expandir e contrair numa migração só, e cnpj e cpf removidos aqui: exceção à
-- regra do documento 2 (seção 14.3), justificada no ADR-013 porque nenhuma
-- versão antiga do código roda em ambiente nenhum; o Go novo entra no mesmo
-- commit. 00001 a 00008 não se editam.

-- +goose Up

-- +goose StatementBegin
SET ROLE tenaz_dono;
SET LOCAL lock_timeout = '3s';
SET LOCAL statement_timeout = '60s';
-- +goose StatementEnd

-- +goose StatementBegin
-- Forma de cada tipo, normalizada: maiúsculas, só letras e dígitos, sem ponto,
-- barra nem hífen (a máscara é só de exibição). Tipo novo entra por CREATE OR
-- REPLACE desta função, com requisito e teste.
CREATE OR REPLACE FUNCTION app.identificacao_valida(p_tipo text, p_valor text) RETURNS boolean
  LANGUAGE sql IMMUTABLE PARALLEL SAFE
AS $$
SELECT CASE p_tipo
  WHEN 'cnpj' THEN p_valor ~ '^[0-9A-Z]{12}[0-9]{2}$'   -- alfanumérico nas 12 primeiras; os 2 dígitos verificadores são numéricos
  WHEN 'cpf'  THEN p_valor ~ '^[0-9]{11}$'
  WHEN 'cuit' THEN p_valor ~ '^[0-9]{11}$'
  WHEN 'cuil' THEN p_valor ~ '^[0-9]{11}$'
  WHEN 'dni'  THEN p_valor ~ '^[0-9]{7,8}$'
  WHEN 'ruc'  THEN p_valor ~ '^[0-9A-Z]{6,10}$'          -- só a forma; RUC antigo pode ter letra (P13)
  WHEN 'ci'   THEN p_valor ~ '^[0-9]{4,8}$'
  ELSE false
END
$$;
-- +goose StatementEnd

-- +goose StatementBegin
-- Os defaults (BR, BRL, pt-BR) deixam a coluna nascer sem reescrever a tabela e
-- dão valor às empresas existentes (todas brasileiras). A criação de conta no Go
-- grava os três explicitamente, a partir do país escolhido.
ALTER TABLE app.empresa
  ADD COLUMN pais               char(2)  NOT NULL DEFAULT 'BR',
  ADD COLUMN moeda              char(3)  NOT NULL DEFAULT 'BRL',
  ADD COLUMN idioma_padrao      text     NOT NULL DEFAULT 'pt-BR',
  ADD COLUMN identificacao_tipo text,
  ADD COLUMN identificacao      text;
ALTER TABLE app.unidade
  ADD COLUMN identificacao_tipo text,
  ADD COLUMN identificacao      text;
ALTER TABLE app.pessoa
  ADD COLUMN identificacao_tipo text,
  ADD COLUMN identificacao      text;
ALTER TABLE app.usuario
  ADD COLUMN idioma text     NOT NULL DEFAULT 'pt-BR',
  ADD COLUMN moeda  char(3)  NOT NULL DEFAULT 'BRL';     -- moeda de custo_hora
-- +goose StatementEnd

-- +goose StatementBegin
-- Cópia dos valores. Como quem aplica a migração (BYPASSRLS: o dono não passa
-- pelo RLS forçado). trg_padrao continua ligado: versao e alterado_xid sobem, e o
-- app offline recebe as linhas na forma nova. trg_auditoria fica desligado só
-- durante a cópia: ela não é ação de usuário, e geraria um evento por linha.
ALTER TABLE app.empresa DISABLE TRIGGER trg_auditoria;
ALTER TABLE app.unidade DISABLE TRIGGER trg_auditoria;
ALTER TABLE app.pessoa  DISABLE TRIGGER trg_auditoria;
RESET ROLE;
SET LOCAL statement_timeout = '15min';

UPDATE app.empresa SET identificacao_tipo = 'cnpj', identificacao = upper(cnpj) WHERE cnpj IS NOT NULL;
UPDATE app.unidade SET identificacao_tipo = 'cnpj', identificacao = upper(cnpj) WHERE cnpj IS NOT NULL;
UPDATE app.pessoa  SET identificacao_tipo = 'cpf',  identificacao = cpf         WHERE cpf  IS NOT NULL;

SET LOCAL statement_timeout = '60s';
SET ROLE tenaz_dono;
ALTER TABLE app.empresa ENABLE TRIGGER trg_auditoria;
ALTER TABLE app.unidade ENABLE TRIGGER trg_auditoria;
ALTER TABLE app.pessoa  ENABLE TRIGGER trg_auditoria;
-- +goose StatementEnd

-- +goose StatementBegin
-- Remoção do antigo, na mesma migração (ver o cabeçalho).
ALTER TABLE app.empresa DROP CONSTRAINT uq_empresa_cnpj, DROP CONSTRAINT ck_empresa_cnpj, DROP COLUMN cnpj;
ALTER TABLE app.unidade DROP CONSTRAINT ck_unidade_cnpj, DROP COLUMN cnpj;
DROP INDEX app.uq_pessoa_cpf;
ALTER TABLE app.pessoa  DROP CONSTRAINT ck_pessoa_cpf, DROP COLUMN cpf;
-- +goose StatementEnd

-- +goose StatementBegin
ALTER TABLE app.empresa
  ADD CONSTRAINT ck_empresa_pais CHECK (pais IN ('BR', 'AR', 'PY')),
  -- uma moeda por empresa, dada pelo país (PYG sem casas decimais)
  ADD CONSTRAINT ck_empresa_moeda CHECK ((pais, moeda) IN (('BR', 'BRL'), ('AR', 'ARS'), ('PY', 'PYG'))),
  ADD CONSTRAINT ck_empresa_idioma CHECK (idioma_padrao IN ('pt-BR', 'es-419')),
  -- os dois preenchidos ou os dois nulos; a empresa só aceita o tipo de pessoa jurídica do próprio país
  ADD CONSTRAINT ck_empresa_identificacao CHECK ((identificacao_tipo IS NULL) = (identificacao IS NULL)
    AND (identificacao IS NULL OR (app.identificacao_valida(identificacao_tipo, identificacao)
         AND (pais, identificacao_tipo) IN (('BR', 'cnpj'), ('AR', 'cuit'), ('PY', 'ruc'))))),
  -- única em toda a plataforma, como o uq_empresa_cnpj que substitui
  ADD CONSTRAINT uq_empresa_identificacao UNIQUE (identificacao_tipo, identificacao);

-- A unidade herda o país da empresa: aceita só os tipos de pessoa jurídica; a
-- coerência com o país da empresa é do serviço (o banco não cruza tabelas num CHECK).
ALTER TABLE app.unidade
  ADD CONSTRAINT ck_unidade_identificacao CHECK ((identificacao_tipo IS NULL) = (identificacao IS NULL)
    AND (identificacao IS NULL OR (identificacao_tipo IN ('cnpj', 'cuit', 'ruc')
         AND app.identificacao_valida(identificacao_tipo, identificacao))));

-- A pessoa aceita os tipos de pessoa física de QUALQUER país suportado: uma
-- transportadora brasileira pode ter motorista paraguaio.
ALTER TABLE app.pessoa
  ADD CONSTRAINT ck_pessoa_identificacao CHECK ((identificacao_tipo IS NULL) = (identificacao IS NULL)
    AND (identificacao IS NULL OR (identificacao_tipo IN ('cpf', 'cuil', 'dni', 'ci')
         AND app.identificacao_valida(identificacao_tipo, identificacao))));
-- RF-ADM-015: identificação única na empresa, quando informada.
CREATE UNIQUE INDEX uq_pessoa_identificacao ON app.pessoa (empresa_id, identificacao_tipo, identificacao)
  WHERE identificacao IS NOT NULL AND excluido_em IS NULL;

ALTER TABLE app.usuario
  ADD CONSTRAINT ck_usuario_idioma CHECK (idioma IN ('pt-BR', 'es-419')),
  ADD CONSTRAINT ck_usuario_moeda CHECK (moeda IN ('BRL', 'ARS', 'PYG'));
-- +goose StatementEnd

-- +goose StatementBegin
-- País e moeda não mudam: trocar o país muda identificação, moeda, conteúdo
-- legal e o aditivo de transferência de dados (P14). Operação assistida do
-- suporte, por migração, nunca pela API. O prefixo do número da OS (OS- no Brasil,
-- OT- na Argentina e no Paraguai) decorre do país e depende desta trava.
CREATE OR REPLACE FUNCTION app.tg_empresa_pais_imutavel() RETURNS trigger
  LANGUAGE plpgsql
AS $$
BEGIN
  IF NEW.pais IS DISTINCT FROM OLD.pais OR NEW.moeda IS DISTINCT FROM OLD.moeda THEN
    RAISE EXCEPTION 'o país e a moeda da empresa não mudam' USING ERRCODE = '42501';
  END IF;
  RETURN NEW;
END
$$;
CREATE TRIGGER trg_empresa_pais_imutavel BEFORE UPDATE ON app.empresa
  FOR EACH ROW EXECUTE FUNCTION app.tg_empresa_pais_imutavel();
-- +goose StatementEnd

-- +goose StatementBegin
-- O idioma é do próprio usuário (perfil), inclusive o do prestador e do
-- locatário: entra na lista de colunas que o escopo externo altera.
CREATE OR REPLACE FUNCTION app.tg_usuario_colunas_externo() RETURNS trigger
  LANGUAGE plpgsql
AS $$
DECLARE
  v_livres text[] := ARRAY['nome', 'telefone', 'idioma',
                           'versao', 'alterado_xid', 'atualizado_em', 'atualizado_por'];
BEGIN
  IF app.escopo_atual() IN ('prestador', 'locatario', 'link')
     AND (to_jsonb(NEW) - v_livres) IS DISTINCT FROM (to_jsonb(OLD) - v_livres) THEN
    RAISE EXCEPTION 'fora do escopo interno, o usuário só altera nome, telefone e idioma' USING ERRCODE = '42501';
  END IF;
  RETURN NEW;
END
$$;
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
SET ROLE tenaz_dono;
SET LOCAL lock_timeout = '3s';
SET LOCAL statement_timeout = '60s';

-- A versão anterior só representa CNPJ e CPF brasileiros, com dígitos. A
-- conferência roda como o migrador: sob RLS forçado, tenaz_dono não enxerga linha
-- alguma e a recusa passaria em falso (achado do teste de migração; seção 14.2).
RESET ROLE;
DO $$
BEGIN
  IF EXISTS (SELECT 1 FROM app.empresa WHERE pais <> 'BR')
     OR EXISTS (SELECT 1 FROM app.empresa WHERE identificacao_tipo IS DISTINCT FROM 'cnpj' AND identificacao IS NOT NULL)
     OR EXISTS (SELECT 1 FROM app.unidade WHERE identificacao_tipo IS DISTINCT FROM 'cnpj' AND identificacao IS NOT NULL)
     OR EXISTS (SELECT 1 FROM app.pessoa WHERE identificacao_tipo IS DISTINCT FROM 'cpf' AND identificacao IS NOT NULL) THEN
    RAISE EXCEPTION 'Down recusado: há empresa de outro país ou identificação que a versão anterior não representa';
  END IF;
END
$$;
SET ROLE tenaz_dono;

DROP TRIGGER trg_empresa_pais_imutavel ON app.empresa;
DROP FUNCTION app.tg_empresa_pais_imutavel();

CREATE OR REPLACE FUNCTION app.tg_usuario_colunas_externo() RETURNS trigger
  LANGUAGE plpgsql
AS $$
DECLARE
  v_livres text[] := ARRAY['nome', 'telefone',
                           'versao', 'alterado_xid', 'atualizado_em', 'atualizado_por'];
BEGIN
  IF app.escopo_atual() IN ('prestador', 'locatario', 'link')
     AND (to_jsonb(NEW) - v_livres) IS DISTINCT FROM (to_jsonb(OLD) - v_livres) THEN
    RAISE EXCEPTION 'fora do escopo interno, o usuário só altera nome e telefone' USING ERRCODE = '42501';
  END IF;
  RETURN NEW;
END
$$;

ALTER TABLE app.empresa ADD COLUMN cnpj char(14);
ALTER TABLE app.unidade ADD COLUMN cnpj char(14);
ALTER TABLE app.pessoa  ADD COLUMN cpf  char(11);

ALTER TABLE app.empresa DISABLE TRIGGER trg_auditoria;
ALTER TABLE app.unidade DISABLE TRIGGER trg_auditoria;
ALTER TABLE app.pessoa  DISABLE TRIGGER trg_auditoria;
RESET ROLE;
UPDATE app.empresa SET cnpj = identificacao WHERE identificacao_tipo = 'cnpj';
UPDATE app.unidade SET cnpj = identificacao WHERE identificacao_tipo = 'cnpj';
UPDATE app.pessoa  SET cpf  = identificacao WHERE identificacao_tipo = 'cpf';
SET ROLE tenaz_dono;
ALTER TABLE app.empresa ENABLE TRIGGER trg_auditoria;
ALTER TABLE app.unidade ENABLE TRIGGER trg_auditoria;
ALTER TABLE app.pessoa  ENABLE TRIGGER trg_auditoria;

ALTER TABLE app.usuario DROP CONSTRAINT ck_usuario_idioma, DROP CONSTRAINT ck_usuario_moeda,
  DROP COLUMN idioma, DROP COLUMN moeda;

DROP INDEX app.uq_pessoa_identificacao;
ALTER TABLE app.pessoa DROP CONSTRAINT ck_pessoa_identificacao, DROP COLUMN identificacao_tipo, DROP COLUMN identificacao;
ALTER TABLE app.unidade DROP CONSTRAINT ck_unidade_identificacao, DROP COLUMN identificacao_tipo, DROP COLUMN identificacao;
ALTER TABLE app.empresa
  DROP CONSTRAINT uq_empresa_identificacao, DROP CONSTRAINT ck_empresa_identificacao,
  DROP CONSTRAINT ck_empresa_idioma, DROP CONSTRAINT ck_empresa_moeda, DROP CONSTRAINT ck_empresa_pais,
  DROP COLUMN identificacao_tipo, DROP COLUMN identificacao, DROP COLUMN idioma_padrao,
  DROP COLUMN moeda, DROP COLUMN pais;

ALTER TABLE app.empresa ADD CONSTRAINT uq_empresa_cnpj UNIQUE (cnpj),
  ADD CONSTRAINT ck_empresa_cnpj CHECK (cnpj ~ '^[0-9]{14}$');
ALTER TABLE app.unidade ADD CONSTRAINT ck_unidade_cnpj CHECK (cnpj ~ '^[0-9]{14}$');
ALTER TABLE app.pessoa ADD CONSTRAINT ck_pessoa_cpf CHECK (cpf ~ '^[0-9]{11}$');
CREATE UNIQUE INDEX uq_pessoa_cpf ON app.pessoa (empresa_id, cpf) WHERE cpf IS NOT NULL AND excluido_em IS NULL;

DROP FUNCTION app.identificacao_valida(text, text);
-- +goose StatementEnd
