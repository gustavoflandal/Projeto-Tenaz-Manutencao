-- Endereço multipaís (ADR-013; documento 3, seção 2.18, versão 0.9; documento 2,
-- seção 5.2, versão 0.9).
--
-- app.unidade.endereco é jsonb (00004) e é a única coluna de endereço do esquema
-- (nenhuma outra tabela existente tem jsonb de endereço, nenhum CHECK nem função
-- olha as chaves). As chaves {uf, cep} viram {subdivisao, codigo_postal}:
--   subdivisao    = 'BR-' || upper(uf)           (código ISO 3166-2 completo)
--   codigo_postal = cep só com dígitos           (omitida se não sobrar dígito)
-- As demais chaves (municipio, logradouro, numero, complemento, bairro) ficam.
-- Não há chave 'pais': o país do endereço é o da empresa.
--
-- Conferência como migrador (RLS forçado esconde linhas do tenaz_dono): antes da
-- 00009 só havia Brasil, então uf/cep só podem existir em empresa BR. Se existir
-- endereço com uf/cep em empresa AR ou PY, a migração aborta (23514): prefixar com
-- 'BR-' ali seria dado errado. No ensaio, 29 endereços, todos de empresa BR.
--
-- Compatível com a versão anterior do código no que o banco exige (jsonb sem
-- CHECK aceita as duas formas); o Go novo lê e grava só a forma nova, e a
-- migração roda no mesmo release (a API nunca foi lançada: CHANGELOG em
-- [Não lançado]). Idempotente: só toca linhas que ainda têm uf ou cep.
-- trg_padrao fica ligado (versao sobe e o app offline recebe a forma nova);
-- trg_auditoria fica desligado só durante a cópia, como na 00009.

-- +goose Up

-- +goose StatementBegin
SET ROLE tenaz_dono;
SET LOCAL lock_timeout = '3s';
RESET ROLE;
DO $$
BEGIN
  IF EXISTS (
    SELECT 1 FROM app.unidade u JOIN app.empresa e ON e.id = u.empresa_id
    WHERE e.pais <> 'BR' AND (u.endereco ? 'uf' OR u.endereco ? 'cep')
  ) THEN
    RAISE EXCEPTION 'há endereço de unidade de empresa não brasileira com uf ou cep; corrija antes da 00011' USING ERRCODE = '23514';
  END IF;
END
$$;
SET ROLE tenaz_dono;

ALTER TABLE app.unidade DISABLE TRIGGER trg_auditoria;
RESET ROLE;
SET LOCAL statement_timeout = '15min';

UPDATE app.unidade
SET endereco = (endereco - 'uf' - 'cep') || jsonb_strip_nulls(jsonb_build_object(
      'subdivisao',    CASE WHEN nullif(btrim(endereco ->> 'uf'), '') IS NOT NULL
                            THEN 'BR-' || upper(btrim(endereco ->> 'uf')) END,
      'codigo_postal', nullif(regexp_replace(coalesce(endereco ->> 'cep', ''), '[^0-9]', '', 'g'), '')))
WHERE jsonb_typeof(endereco) = 'object' AND (endereco ? 'uf' OR endereco ? 'cep');

SET LOCAL statement_timeout = '60s';
SET ROLE tenaz_dono;
ALTER TABLE app.unidade ENABLE TRIGGER trg_auditoria;
RESET ROLE;
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
-- Inverso. Empresa BR: uf = sufixo de 'BR-xx'. Empresa AR ou PY (a forma antiga
-- nunca as teve): o sufixo e o código postal também voltam, sem perda, em uf e cep
-- (AR-C vira uf 'C'; PY-11 vira uf '11'); o código anterior não os validava na
-- leitura. Subdivisão sem hífen fica inteira em uf.
SET ROLE tenaz_dono;
SET LOCAL lock_timeout = '3s';
ALTER TABLE app.unidade DISABLE TRIGGER trg_auditoria;
RESET ROLE;
SET LOCAL statement_timeout = '15min';

UPDATE app.unidade
SET endereco = (endereco - 'subdivisao' - 'codigo_postal') || jsonb_strip_nulls(jsonb_build_object(
      'uf',  regexp_replace(endereco ->> 'subdivisao', '^[A-Z]{2}-', ''),
      'cep', endereco ->> 'codigo_postal'))
WHERE jsonb_typeof(endereco) = 'object' AND (endereco ? 'subdivisao' OR endereco ? 'codigo_postal');

SET LOCAL statement_timeout = '60s';
SET ROLE tenaz_dono;
ALTER TABLE app.unidade ENABLE TRIGGER trg_auditoria;
RESET ROLE;
-- +goose StatementEnd
