-- Consultas do serviço de autenticação (RF-ADM-004, RF-ADM-006, RNF-017).
--
-- As tabelas do esquema plataforma não têm RLS: são lidas antes de a empresa
-- ser conhecida (seção 5.1 do documento 2). Aqui nenhuma consulta devolve
-- segredo: só hashes, comparados no Go em tempo constante.
--
-- Este arquivo só toca as tabelas de credencial do esquema plataforma, e roda
-- no pool próprio da autenticação (achado ISO-E01-02): nunca pelo pool da
-- aplicação nem por NaEmpresa. Busca de identidade é sempre por chave exata
-- (id, e-mail ou hash), nunca por listagem. As leituras de usuário, perfis e
-- empresa estão em internal/plataforma/autorizacao/autorizacao.sql.

-- ---------------------------------------------------------------- identidade

-- name: AuthBuscarIdentidadePorEmail :one
SELECT id, email, senha_hash, mfa_segredo_cifrado, mfa_ativado_em, email_verificado_em,
       falhas_login, bloqueada_ate, bloqueada_em
FROM plataforma.identidade
WHERE email = sqlc.arg(email);

-- name: AuthBuscarIdentidade :one
SELECT id, email, senha_hash, mfa_segredo_cifrado, mfa_ativado_em, email_verificado_em,
       falhas_login, bloqueada_ate, bloqueada_em
FROM plataforma.identidade
WHERE id = sqlc.arg(id);

-- name: AuthTravarIdentidade :one
SELECT id, email, senha_hash, mfa_segredo_cifrado, mfa_ativado_em, email_verificado_em,
       falhas_login, bloqueada_ate, bloqueada_em
FROM plataforma.identidade
WHERE id = sqlc.arg(id)
FOR UPDATE;

-- name: AuthCriarIdentidade :one
INSERT INTO plataforma.identidade (id, email, senha_hash, email_verificado_em)
VALUES (sqlc.arg(id), sqlc.arg(email), sqlc.arg(senha_hash), now())
ON CONFLICT (email) DO NOTHING
RETURNING id;

-- Falha de login: incremento atômico, para que duas tentativas simultâneas
-- não contem como uma. O bloqueio é decidido no Go a partir do total.
-- name: AuthIncrementarFalhasLogin :one
UPDATE plataforma.identidade
SET falhas_login = least(falhas_login + 1, 32000), atualizado_em = now()
WHERE id = sqlc.arg(id)
RETURNING falhas_login;

-- name: AuthBloquearIdentidadeAte :exec
UPDATE plataforma.identidade
SET bloqueada_ate = sqlc.arg(bloqueada_ate), atualizado_em = now()
WHERE id = sqlc.arg(id);

-- Login certo: zera as falhas e, quando os parâmetros do Argon2id mudaram,
-- grava o hash recalculado.
-- name: AuthRegistrarLoginCerto :exec
UPDATE plataforma.identidade
SET falhas_login = 0, bloqueada_ate = NULL, ultimo_login_em = now(), atualizado_em = now(),
    senha_hash = coalesce(sqlc.narg(senha_hash_novo), senha_hash)
WHERE id = sqlc.arg(id);

-- name: AuthDefinirSenha :exec
UPDATE plataforma.identidade
SET senha_hash = sqlc.arg(senha_hash), falhas_login = 0, bloqueada_ate = NULL, atualizado_em = now()
WHERE id = sqlc.arg(id);

-- name: AuthMarcarEmailVerificado :exec
UPDATE plataforma.identidade
SET email_verificado_em = coalesce(email_verificado_em, now()), atualizado_em = now()
WHERE id = sqlc.arg(id);

-- Ativação em dois passos: o segredo é gravado pendente (mfa_ativado_em nulo)
-- e só vale depois de confirmado com um código.
-- name: AuthGravarSegredoMFAPendente :exec
UPDATE plataforma.identidade
SET mfa_segredo_cifrado = sqlc.arg(mfa_segredo_cifrado), mfa_ativado_em = NULL, mfa_ultimo_passo = NULL,
    atualizado_em = now()
WHERE id = sqlc.arg(id) AND mfa_ativado_em IS NULL;

-- name: AuthAtivarMFA :execrows
UPDATE plataforma.identidade
SET mfa_ativado_em = now(), atualizado_em = now()
WHERE id = sqlc.arg(id) AND mfa_segredo_cifrado IS NOT NULL AND mfa_ativado_em IS NULL;

-- Consumo do passo TOTP (RFC 6238, seção 5.2): o mesmo código não vale duas
-- vezes. Atômico: zero linhas = passo já usado (ou anterior ao último aceito).
-- name: AuthUsarPassoTOTP :execrows
UPDATE plataforma.identidade
SET mfa_ultimo_passo = sqlc.arg(passo), atualizado_em = now()
WHERE id = sqlc.arg(id) AND (mfa_ultimo_passo IS NULL OR mfa_ultimo_passo < sqlc.arg(passo));

-- name: AuthDesativarMFA :exec
UPDATE plataforma.identidade
SET mfa_segredo_cifrado = NULL, mfa_ativado_em = NULL, mfa_ultimo_passo = NULL, atualizado_em = now()
WHERE id = sqlc.arg(id);

-- ------------------------------------------------------------ acesso_empresa

-- name: AuthListarAcessos :many
SELECT empresa_id, usuario_id, usuario_nome
FROM plataforma.acesso_empresa
WHERE identidade_id = sqlc.arg(identidade_id) AND habilitado
ORDER BY empresa_id;

-- Identidade de um usuário numa empresa, para o evento de encerrar sessões a
-- pedido do administrador. Filtro (empresa_id, usuario_id).
-- name: AuthIdentidadeDoUsuario :one
SELECT identidade_id
FROM plataforma.acesso_empresa
WHERE empresa_id = sqlc.arg(empresa_id) AND usuario_id = sqlc.arg(usuario_id);

-- Nome do usuário naquela empresa, para o evento da trilha (evento_auditoria.
-- usuario_nome). Filtro (identidade_id, empresa_id), a chave primária.
-- name: AuthNomeDoAcesso :one
SELECT usuario_nome
FROM plataforma.acesso_empresa
WHERE identidade_id = sqlc.arg(identidade_id) AND empresa_id = sqlc.arg(empresa_id);

-- name: AuthBuscarAcesso :one
SELECT usuario_id
FROM plataforma.acesso_empresa
WHERE identidade_id = sqlc.arg(identidade_id) AND empresa_id = sqlc.arg(empresa_id) AND habilitado;

-- Empresa da sessão mais recente da identidade: a escolha padrão no login.
-- name: AuthUltimaEmpresa :one
SELECT s.empresa_id
FROM plataforma.sessao s
JOIN plataforma.acesso_empresa a
  ON a.identidade_id = s.identidade_id AND a.empresa_id = s.empresa_id AND a.habilitado
WHERE s.identidade_id = sqlc.arg(identidade_id)
ORDER BY s.criado_em DESC
LIMIT 1;

-- -------------------------------------------------------------------- sessao

-- Os prazos têm teto em CHECK contra criado_em. criado_em vai explícito, do
-- mesmo relógio que calculou expira_em: com o padrão now(), um relógio do
-- servidor de aplicação poucos milissegundos à frente do banco faria o teto
-- exato (7 dias, 30 dias, 1 hora) falhar.

-- name: AuthCriarSessao :exec
INSERT INTO plataforma.sessao
  (id, canal, token_hash, csrf_hash, identidade_id, empresa_id, usuario_id, dispositivo_id, ip, user_agent,
   criado_em, ultimo_uso_em, expira_em)
VALUES
  (sqlc.arg(id), sqlc.arg(canal), sqlc.narg(token_hash), sqlc.narg(csrf_hash), sqlc.arg(identidade_id),
   sqlc.arg(empresa_id), sqlc.arg(usuario_id), sqlc.narg(dispositivo_id),
   nullif(sqlc.arg(ip)::text, '')::inet, nullif(sqlc.arg(user_agent)::text, ''), sqlc.arg(criado_em),
   sqlc.arg(criado_em), sqlc.arg(expira_em));

-- name: AuthBuscarSessaoPorToken :one
SELECT id, canal, csrf_hash, identidade_id, empresa_id, usuario_id, criado_em, ultimo_uso_em, expira_em,
       revogada_em
FROM plataforma.sessao
WHERE token_hash = sqlc.arg(token_hash);

-- name: AuthBuscarSessao :one
SELECT id, canal, csrf_hash, identidade_id, empresa_id, usuario_id, criado_em, ultimo_uso_em, expira_em,
       revogada_em
FROM plataforma.sessao
WHERE id = sqlc.arg(id);

-- name: AuthTocarSessao :exec
UPDATE plataforma.sessao
SET ultimo_uso_em = now()
WHERE id = sqlc.arg(id) AND revogada_em IS NULL;

-- name: AuthEstenderSessao :exec
UPDATE plataforma.sessao
SET ultimo_uso_em = now(), expira_em = sqlc.arg(expira_em)
WHERE id = sqlc.arg(id) AND revogada_em IS NULL;

-- name: AuthRevogarSessao :execrows
UPDATE plataforma.sessao
SET revogada_em = now()
WHERE id = sqlc.arg(id) AND revogada_em IS NULL;

-- name: AuthRevogarSessoesUsuario :many
UPDATE plataforma.sessao
SET revogada_em = now()
WHERE empresa_id = sqlc.arg(empresa_id) AND usuario_id = sqlc.arg(usuario_id) AND revogada_em IS NULL
RETURNING id;

-- Troca e redefinição de senha: todas as sessões da identidade, em todas as
-- empresas, menos a que fez a troca (quando houver).
-- name: AuthRevogarSessoesIdentidade :many
UPDATE plataforma.sessao
SET revogada_em = now()
WHERE identidade_id = sqlc.arg(identidade_id) AND revogada_em IS NULL
  AND id IS DISTINCT FROM sqlc.narg(exceto_sessao_id)::uuid
RETURNING id;

-- ----------------------------------------------------------- token_renovacao

-- name: AuthCriarTokenRenovacao :exec
INSERT INTO plataforma.token_renovacao (id, familia_id, anterior_id, token_hash, criado_em, expira_em)
VALUES (sqlc.arg(id), sqlc.arg(familia_id), sqlc.narg(anterior_id), sqlc.arg(token_hash), sqlc.arg(criado_em),
        sqlc.arg(expira_em));

-- Rotação atômica (seção 5.1 do documento 2): zero linhas com o hash
-- existente é reuso, expiração ou revogação.
-- name: AuthUsarTokenRenovacao :one
UPDATE plataforma.token_renovacao
SET usado_em = now()
WHERE token_hash = sqlc.arg(token_hash) AND usado_em IS NULL AND revogado_em IS NULL AND expira_em > now()
RETURNING id, familia_id;

-- name: AuthBuscarTokenRenovacao :one
SELECT id, familia_id, usado_em, revogado_em, expira_em
FROM plataforma.token_renovacao
WHERE token_hash = sqlc.arg(token_hash);

-- name: AuthRevogarFamilias :exec
UPDATE plataforma.token_renovacao
SET revogado_em = now()
WHERE familia_id = ANY (sqlc.arg(familias)::uuid[]) AND revogado_em IS NULL;

-- -------------------------------------------------------- codigo_autorizacao

-- name: AuthCriarCodigoAutorizacao :exec
INSERT INTO plataforma.codigo_autorizacao
  (id, codigo_hash, identidade_id, empresa_id, usuario_id, cliente_oauth, uri_retorno, desafio_pkce, criado_em,
   expira_em)
VALUES
  (sqlc.arg(id), sqlc.arg(codigo_hash), sqlc.arg(identidade_id), sqlc.arg(empresa_id), sqlc.arg(usuario_id),
   sqlc.arg(cliente_oauth), sqlc.arg(uri_retorno), sqlc.arg(desafio_pkce), sqlc.arg(criado_em),
   sqlc.arg(expira_em));

-- name: AuthUsarCodigoAutorizacao :one
UPDATE plataforma.codigo_autorizacao
SET usado_em = now()
WHERE codigo_hash = sqlc.arg(codigo_hash) AND usado_em IS NULL AND expira_em > now()
RETURNING id, identidade_id, empresa_id, usuario_id, cliente_oauth, uri_retorno, desafio_pkce;

-- name: AuthBuscarCodigoAutorizacao :one
SELECT id, usado_em, sessao_id
FROM plataforma.codigo_autorizacao
WHERE codigo_hash = sqlc.arg(codigo_hash);

-- name: AuthLigarSessaoAoCodigo :exec
UPDATE plataforma.codigo_autorizacao
SET sessao_id = sqlc.arg(sessao_id)
WHERE id = sqlc.arg(id);

-- --------------------------------------------------------------- desafio_mfa

-- name: AuthCriarDesafioMFA :exec
INSERT INTO plataforma.desafio_mfa (id, token_hash, identidade_id, canal, ip, criado_em, expira_em)
VALUES (sqlc.arg(id), sqlc.arg(token_hash), sqlc.arg(identidade_id), sqlc.arg(canal),
        nullif(sqlc.arg(ip)::text, '')::inet, sqlc.arg(criado_em), sqlc.arg(expira_em));

-- name: AuthTravarDesafioMFA :one
SELECT id, identidade_id, canal, tentativas, expira_em, concluido_em
FROM plataforma.desafio_mfa
WHERE token_hash = sqlc.arg(token_hash)
FOR UPDATE;

-- name: AuthIncrementarTentativaMFA :exec
UPDATE plataforma.desafio_mfa
SET tentativas = tentativas + 1
WHERE id = sqlc.arg(id) AND tentativas < 5;

-- name: AuthConcluirDesafioMFA :exec
UPDATE plataforma.desafio_mfa
SET concluido_em = now()
WHERE id = sqlc.arg(id) AND concluido_em IS NULL;

-- --------------------------------------------------------------- token_email

-- name: AuthCriarTokenEmail :exec
INSERT INTO plataforma.token_email
  (id, token_hash, finalidade, email, identidade_id, empresa_id, usuario_id, criado_em, expira_em)
VALUES (sqlc.arg(id), sqlc.arg(token_hash), sqlc.arg(finalidade), sqlc.arg(email), sqlc.narg(identidade_id),
        sqlc.narg(empresa_id), sqlc.narg(usuario_id), sqlc.arg(criado_em), sqlc.arg(expira_em));

-- Novo envio revoga o anterior do mesmo alvo e da mesma finalidade.
-- name: AuthRevogarTokensEmailUsuario :exec
UPDATE plataforma.token_email
SET revogado_em = now()
WHERE empresa_id = sqlc.arg(empresa_id) AND usuario_id = sqlc.arg(usuario_id) AND finalidade = sqlc.arg(finalidade)
  AND usado_em IS NULL AND revogado_em IS NULL;

-- name: AuthRevogarTokensEmailIdentidade :exec
UPDATE plataforma.token_email
SET revogado_em = now()
WHERE identidade_id = sqlc.arg(identidade_id) AND finalidade = sqlc.arg(finalidade)
  AND usado_em IS NULL AND revogado_em IS NULL;

-- Leitura sem consumo: descobre a empresa do convite antes de abrir a
-- transação dela. O consumo acontece dentro dessa transação.
-- name: AuthBuscarTokenEmailValido :one
SELECT id, finalidade, email, identidade_id, empresa_id, usuario_id
FROM plataforma.token_email
WHERE token_hash = sqlc.arg(token_hash) AND finalidade = sqlc.arg(finalidade)
  AND usado_em IS NULL AND revogado_em IS NULL AND expira_em > now();

-- name: AuthUsarTokenEmail :one
UPDATE plataforma.token_email
SET usado_em = now()
WHERE token_hash = sqlc.arg(token_hash) AND finalidade = sqlc.arg(finalidade)
  AND usado_em IS NULL AND revogado_em IS NULL AND expira_em > now()
RETURNING id, email, identidade_id, empresa_id, usuario_id;

-- ------------------------------------------------- trilha (app.evento_auditoria)

-- Evento de autenticação na trilha da empresa (seção 6.9 do documento 5),
-- gravado na transação NaTrilha, que fixou app.empresa_id e os dados da
-- requisição. Só estas 16 colunas têm INSERT (as 13 da 00005, mais nome, resultado e tela da 00007) para tenaz_autenticacao, e não há
-- SELECT: sem RETURNING. ator_tipo e entidade são fixos; a política restritiva
-- do papel confere a ação e que usuario_id é o da identidade naquela empresa.
-- O nome do usuário no evento é o de acesso_empresa naquela empresa
-- (AuthNomeDoAcesso): a policy autenticacao_identidade confere que o nome
-- gravado é o desta linha.
-- name: AuthInserirEvento :exec
INSERT INTO app.evento_auditoria
  (empresa_id, ator_tipo, usuario_id, usuario_nome, acao, resultado, entidade, entidade_id, depois,
   campos_alterados, origem, ip, dispositivo_id, user_agent, request_id, tela)
VALUES
  (sqlc.arg(empresa_id), 'usuario', sqlc.arg(usuario_id), sqlc.arg(usuario_nome), sqlc.arg(acao),
   sqlc.arg(resultado), 'identidade', sqlc.arg(identidade_id),
   sqlc.narg(depois), sqlc.narg(campos_alterados)::text[], nullif(current_setting('app.origem', true), ''),
   nullif(current_setting('app.ip', true), '')::inet, sqlc.narg(dispositivo_id),
   nullif(current_setting('app.user_agent', true), ''), nullif(current_setting('app.request_id', true), ''),
   sqlc.narg(tela));
