-- Fecha o evento retroativo (achado do tenaz-isolamento sobre a 00007).
--
-- O tenaz_app e o tenaz_job tinham INSERT em todas as colunas de
-- app.evento_auditoria, inclusive ocorrido_em, id, transacao e hash. Com
-- ocorrido_em à escolha, um evento podia ser gravado num dia ainda não selado
-- com data retroativa, ou num dia já selado (e a verificação o acusaria, mas só
-- depois). O trigger de BEFORE INSERT não resolve: o migrador, que grava as
-- fixtures e a conversão da 00007, também passa por ele. O privilégio por
-- coluna, sim: é o molde que o tenaz_sinc e o tenaz_autenticacao já seguem.
--
-- Ficam de fora do INSERT dos dois papéis: id, ocorrido_em e transacao (valores
-- padrão do servidor) e hash (do trigger). tg_auditar e o serviço de auditoria
-- não listam nenhuma delas. Fonte: documento 2, seção 9.2.

-- +goose Up

-- +goose StatementBegin
SET ROLE tenaz_dono;
SET LOCAL lock_timeout = '3s';
REVOKE INSERT ON app.evento_auditoria FROM tenaz_app, tenaz_job;
GRANT INSERT (empresa_id, ator_tipo, usuario_id, usuario_nome, token_api_id, acao, resultado, sensivel,
              entidade, entidade_id, entidade_chave, antes, depois, detalhe, campos_alterados, modulo, tela,
              origem, ip, dispositivo_id, user_agent, request_id, conexao_empresa_id,
              ocorrido_no_dispositivo_em, escopo)
  ON app.evento_auditoria TO tenaz_app, tenaz_job;
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
SET ROLE tenaz_dono;
REVOKE INSERT (empresa_id, ator_tipo, usuario_id, usuario_nome, token_api_id, acao, resultado, sensivel,
               entidade, entidade_id, entidade_chave, antes, depois, detalhe, campos_alterados, modulo, tela,
               origem, ip, dispositivo_id, user_agent, request_id, conexao_empresa_id,
               ocorrido_no_dispositivo_em, escopo)
  ON app.evento_auditoria FROM tenaz_app, tenaz_job;
GRANT INSERT ON app.evento_auditoria TO tenaz_app, tenaz_job;
-- +goose StatementEnd
