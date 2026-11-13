-- Eventos de auditoria gravados pelo serviço (seção 9.1 do documento 2), e
-- não pelo trg_auditoria: login, falha de login, logout, troca de senha,
-- ativação de MFA, visualização de dado sensível.
--
-- O contexto técnico (origem, IP, dispositivo, user agent, request_id,
-- token de API, instante do aparelho) vem das mesmas variáveis da transação
-- que app.tg_auditar() lê, fixadas por NaEmpresa: um evento gravado pelo
-- serviço e um gravado pelo trigger na mesma requisição ficam iguais nesses
-- campos. A empresa é a da transação; a policy auditoria_insercao recusa
-- qualquer outra.

-- name: InserirEventoAuditoria :exec
INSERT INTO app.evento_auditoria
  (empresa_id, ator_tipo, usuario_id, token_api_id, acao, entidade, entidade_id, entidade_chave,
   antes, depois, campos_alterados, origem, ip, dispositivo_id, user_agent, request_id,
   ocorrido_no_dispositivo_em)
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
   nullif(current_setting('app.ocorrido_no_dispositivo_em', true), '')::timestamptz);
