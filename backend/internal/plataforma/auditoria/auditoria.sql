-- Eventos de auditoria gravados pelo serviço (seção 9.1 do documento 2), e
-- não pelo trg_auditoria: visualização, impressão, extração, alteração no
-- nível do fato. Login, logout e falha de login são do papel da autenticação
-- (ADR-011), com consulta própria.
--
-- O contexto técnico (origem, IP, dispositivo, user agent, request_id,
-- token de API, instante do aparelho) vem das mesmas variáveis da transação
-- que app.tg_auditar() lê, fixadas por NaEmpresa: um evento gravado pelo
-- serviço e um gravado pelo trigger na mesma requisição ficam iguais nesses
-- campos. O módulo, a tela e o nome do usuário também vêm de lá, pelo trigger
-- app.tg_evento_auditoria_preparar(), que ainda calcula o hash. A empresa é a
-- da transação; a policy auditoria_insercao recusa qualquer outra.

-- name: InserirEventoAuditoria :exec
INSERT INTO app.evento_auditoria
  (empresa_id, ator_tipo, usuario_id, token_api_id, acao, entidade, entidade_id, entidade_chave,
   antes, depois, campos_alterados, origem, ip, dispositivo_id, user_agent, request_id,
   ocorrido_no_dispositivo_em, sensivel, detalhe, escopo)
VALUES
  (sqlc.arg(empresa_id),
   sqlc.arg(ator_tipo),
   sqlc.narg(usuario_id),
   nullif(current_setting('app.token_api_id', true), '')::uuid,
   sqlc.arg(acao),
   sqlc.arg(entidade),
   sqlc.narg(entidade_id),
   sqlc.narg(entidade_chave),
   sqlc.narg(antes),
   sqlc.narg(depois),
   sqlc.narg(campos_alterados)::text[],
   nullif(current_setting('app.origem', true), ''),
   nullif(current_setting('app.ip', true), '')::inet,
   nullif(current_setting('app.dispositivo_id', true), '')::uuid,
   nullif(current_setting('app.user_agent', true), ''),
   nullif(current_setting('app.request_id', true), ''),
   nullif(current_setting('app.ocorrido_no_dispositivo_em', true), '')::timestamptz,
   sqlc.arg(sensivel)::boolean,
   sqlc.narg(detalhe),
   app.escopo_atual());

-- ------------------------------------------------------------------ captura

-- Sem linha vale o padrão: impressão e extração ligadas, visualização
-- desligada (documento 2, seção 9.2.2). Só o escopo interno lê.
-- name: AudListarCapturas :many
SELECT tela, visualizacao, impressao, extracao
FROM app.captura_auditoria
WHERE empresa_id = sqlc.arg(empresa_id) AND excluido_em IS NULL
ORDER BY tela;

-- As três chaves em cada gravação: o UPDATE do upsert não pode devolver a
-- outra chave ao padrão.
-- name: AudGravarCaptura :exec
INSERT INTO app.captura_auditoria (empresa_id, tela, visualizacao, impressao, extracao, criado_por, atualizado_por)
VALUES (sqlc.arg(empresa_id), sqlc.arg(tela), sqlc.arg(visualizacao), sqlc.arg(impressao), sqlc.arg(extracao),
        sqlc.narg(usuario_id), sqlc.narg(usuario_id))
ON CONFLICT (empresa_id, tela) DO UPDATE
SET visualizacao   = EXCLUDED.visualizacao,
    impressao      = EXCLUDED.impressao,
    extracao       = EXCLUDED.extracao,
    atualizado_por = EXCLUDED.atualizado_por,
    excluido_em    = NULL;

-- ---------------------------------------------------------------------- selo

-- Job do selo, como tenaz_job, na transação da empresa: sela os dias pendentes
-- e devolve os dias selados.
-- name: AudSelarDiasPendentes :many
SELECT app.selar_dias_pendentes()::date AS dia;

-- Selos de um período, para enviar ao armazenamento e para a verificação.
-- name: AudListarSelos :many
SELECT dia, quantidade, hash_dia, hash_anterior, hash_encadeado, calculado_em, chave_objeto
FROM app.selo_auditoria
WHERE empresa_id = sqlc.arg(empresa_id) AND dia >= sqlc.arg(de)::date AND dia <= sqlc.arg(ate)::date
ORDER BY dia;

-- app.verificar_integridade_auditoria(de, ate) devolve TABLE: o sqlc não
-- resolve as colunas de função que retorna tabela, e a chamada fica em
-- verificacao.go, com pgx e parâmetros.

-- Retenção da empresa, para o retain-until do selo.
-- name: AudRetencaoEmpresa :one
SELECT retencao_auditoria_anos::int AS anos
FROM app.empresa
WHERE id = sqlc.arg(empresa_id);
