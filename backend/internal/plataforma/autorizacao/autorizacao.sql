-- Leitura das permissões efetivas e do escopo do usuário, a cada requisição
-- (seção 2.9.2 do documento 2). Roda em NaEmpresa com escopo interno e sem
-- usuário, pelo pool da aplicação, e filtra empresa_id também de forma
-- explícita.

-- name: AutzBuscarUsuario :one
SELECT id, identidade_id, pessoa_id, fornecedor_id, locatario_id, nome, email, situacao, idioma
FROM app.usuario
WHERE empresa_id = sqlc.arg(empresa_id) AND id = sqlc.arg(id) AND excluido_em IS NULL;

-- Perfis ativos do usuário, com o código de origem (o tipo do perfil) e as
-- permissões. unidade_id (RF-ADM-003, Fase 2) ainda não restringe.
-- name: AutzPerfisDoUsuario :many
SELECT coalesce(p.perfil_origem, p.codigo)::text AS codigo_origem, p.permissoes
FROM app.usuario_perfil up
JOIN app.perfil p ON p.empresa_id = up.empresa_id AND p.id = up.perfil_id
WHERE up.empresa_id = sqlc.arg(empresa_id) AND up.usuario_id = sqlc.arg(usuario_id)
  AND up.excluido_em IS NULL AND p.excluido_em IS NULL;

-- name: AutzPermissoesAvulsas :many
SELECT permissao
FROM app.usuario_permissao
WHERE empresa_id = sqlc.arg(empresa_id) AND usuario_id = sqlc.arg(usuario_id) AND excluido_em IS NULL;

-- name: AutzBuscarEmpresa :one
SELECT id, razao_social, nome_fantasia, plano, srv_habilitado, mfa_obrigatorio, situacao, idioma_padrao, pais, moeda
FROM app.empresa
WHERE id = sqlc.arg(id) AND excluido_em IS NULL;
