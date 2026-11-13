-- Consultas do módulo ADM (RF-ADM-001, 002, 004, 005, 011, 015).
--
-- Regra que o teste de vazamento cobra: toda consulta sobre tabela de cliente
-- cita empresa_id no WHERE, mesmo com RLS ligado (duas defesas, de propósito).
--
-- Tipos sem equivalente direto no Go saem convertidos aqui, para o domínio
-- não lidar com pgtype: hora como 'HH24:MI', data como 'YYYY-MM-DD' (vazio
-- quando nula) e coordenadas em float8. Dinheiro entra em centavos (bigint) e
-- sai como numeric exato, convertido para centavos no Go.

-- ------------------------------------------------------------------- empresa

-- name: AdmCriarEmpresa :exec
INSERT INTO app.empresa (id, razao_social, nome_fantasia, cnpj, plano, srv_habilitado, fuso_horario)
VALUES (sqlc.arg(id), sqlc.arg(razao_social), sqlc.narg(nome_fantasia), sqlc.arg(cnpj), sqlc.narg(plano),
        sqlc.arg(srv_habilitado), sqlc.arg(fuso_horario));

-- name: AdmBuscarEmpresa :one
SELECT id, razao_social, nome_fantasia, cnpj, plano, srv_habilitado, mfa_obrigatorio, situacao, fuso_horario,
       retencao_auditoria_anos, coalesce(to_char(teste_ate, 'YYYY-MM-DD'), '')::text AS teste_ate, versao, criado_em, atualizado_em
FROM app.empresa
WHERE id = sqlc.arg(id) AND excluido_em IS NULL;

-- name: AdmTravarEmpresa :one
SELECT id, plano, srv_habilitado, versao
FROM app.empresa
WHERE id = sqlc.arg(id) AND excluido_em IS NULL
FOR UPDATE;

-- name: AdmAtualizarEmpresa :exec
UPDATE app.empresa
SET razao_social = sqlc.arg(razao_social), nome_fantasia = sqlc.narg(nome_fantasia), cnpj = sqlc.arg(cnpj),
    fuso_horario = sqlc.arg(fuso_horario), mfa_obrigatorio = sqlc.arg(mfa_obrigatorio)
WHERE id = sqlc.arg(id) AND excluido_em IS NULL;

-- ------------------------------------------------------------------- unidade

-- name: AdmCriarUnidade :exec
INSERT INTO app.unidade (empresa_id, id, codigo, nome, cnpj, endereco, fuso_horario, latitude, longitude)
VALUES (sqlc.arg(empresa_id), sqlc.arg(id), sqlc.arg(codigo), sqlc.arg(nome), sqlc.narg(cnpj), sqlc.narg(endereco),
        sqlc.narg(fuso_horario), sqlc.narg(latitude)::float8, sqlc.narg(longitude)::float8);

-- name: AdmBuscarUnidade :one
SELECT id, codigo, nome, cnpj, endereco, fuso_horario, latitude, longitude, situacao, versao, criado_em, atualizado_em
FROM app.unidade
WHERE empresa_id = sqlc.arg(empresa_id) AND id = sqlc.arg(id) AND excluido_em IS NULL;

-- name: AdmTravarUnidade :one
SELECT id, situacao, versao
FROM app.unidade
WHERE empresa_id = sqlc.arg(empresa_id) AND id = sqlc.arg(id) AND excluido_em IS NULL
FOR UPDATE;

-- name: AdmListarUnidades :many
SELECT id, codigo, nome, cnpj, endereco, fuso_horario, latitude, longitude, situacao, versao, criado_em, atualizado_em
FROM app.unidade
WHERE empresa_id = sqlc.arg(empresa_id) AND excluido_em IS NULL
ORDER BY nome, id;

-- A regra "Essencial tem uma unidade" conta as unidades não excluídas,
-- inclusive as inativas: inativar não libera vaga.
-- name: AdmContarUnidades :one
SELECT count(*)::bigint
FROM app.unidade
WHERE empresa_id = sqlc.arg(empresa_id) AND excluido_em IS NULL;

-- name: AdmContarUnidadesAtivas :one
SELECT count(*)::bigint
FROM app.unidade
WHERE empresa_id = sqlc.arg(empresa_id) AND excluido_em IS NULL AND situacao = 'ativa';

-- name: AdmAtualizarUnidade :exec
UPDATE app.unidade
SET codigo = sqlc.arg(codigo), nome = sqlc.arg(nome), cnpj = sqlc.narg(cnpj), endereco = sqlc.narg(endereco),
    fuso_horario = sqlc.narg(fuso_horario), latitude = sqlc.narg(latitude)::float8,
    longitude = sqlc.narg(longitude)::float8
WHERE empresa_id = sqlc.arg(empresa_id) AND id = sqlc.arg(id) AND excluido_em IS NULL;

-- name: AdmDefinirSituacaoUnidade :exec
UPDATE app.unidade
SET situacao = sqlc.arg(situacao)
WHERE empresa_id = sqlc.arg(empresa_id) AND id = sqlc.arg(id) AND excluido_em IS NULL;

-- --------------------------------------------------------- turno e feriado

-- name: AdmCriarTurno :exec
INSERT INTO app.turno_unidade (empresa_id, id, unidade_id, dia_semana, hora_inicio, hora_fim, nome)
VALUES (sqlc.arg(empresa_id), sqlc.arg(id), sqlc.arg(unidade_id), sqlc.arg(dia_semana),
        sqlc.arg(hora_inicio)::text::time, sqlc.arg(hora_fim)::text::time, sqlc.narg(nome));

-- name: AdmListarTurnos :many
SELECT id, unidade_id, dia_semana, to_char(hora_inicio, 'HH24:MI')::text AS hora_inicio,
       to_char(hora_fim, 'HH24:MI')::text AS hora_fim, nome, versao
FROM app.turno_unidade
WHERE empresa_id = sqlc.arg(empresa_id) AND unidade_id = sqlc.arg(unidade_id) AND excluido_em IS NULL
ORDER BY dia_semana, hora_inicio;

-- name: AdmBuscarTurno :one
SELECT id, unidade_id, dia_semana, to_char(hora_inicio, 'HH24:MI')::text AS hora_inicio,
       to_char(hora_fim, 'HH24:MI')::text AS hora_fim, nome, versao
FROM app.turno_unidade
WHERE empresa_id = sqlc.arg(empresa_id) AND id = sqlc.arg(id) AND excluido_em IS NULL;

-- name: AdmAtualizarTurno :exec
UPDATE app.turno_unidade
SET dia_semana = sqlc.arg(dia_semana), hora_inicio = sqlc.arg(hora_inicio)::text::time,
    hora_fim = sqlc.arg(hora_fim)::text::time, nome = sqlc.narg(nome)
WHERE empresa_id = sqlc.arg(empresa_id) AND id = sqlc.arg(id) AND excluido_em IS NULL;

-- name: AdmRemoverTurno :exec
UPDATE app.turno_unidade
SET excluido_em = now()
WHERE empresa_id = sqlc.arg(empresa_id) AND id = sqlc.arg(id) AND excluido_em IS NULL;

-- name: AdmContarTurnos :one
SELECT count(*)::bigint
FROM app.turno_unidade
WHERE empresa_id = sqlc.arg(empresa_id) AND unidade_id = sqlc.arg(unidade_id) AND excluido_em IS NULL;

-- name: AdmCriarFeriado :exec
INSERT INTO app.feriado_unidade (empresa_id, id, unidade_id, data_feriado, descricao)
VALUES (sqlc.arg(empresa_id), sqlc.arg(id), sqlc.arg(unidade_id), sqlc.arg(data_feriado)::text::date,
        sqlc.arg(descricao));

-- name: AdmListarFeriados :many
SELECT id, unidade_id, to_char(data_feriado, 'YYYY-MM-DD')::text AS data_feriado, descricao, versao
FROM app.feriado_unidade
WHERE empresa_id = sqlc.arg(empresa_id) AND unidade_id = sqlc.arg(unidade_id) AND excluido_em IS NULL
  AND (sqlc.narg(de)::text IS NULL OR data_feriado >= sqlc.narg(de)::text::date)
  AND (sqlc.narg(ate)::text IS NULL OR data_feriado <= sqlc.narg(ate)::text::date)
ORDER BY data_feriado;

-- name: AdmRemoverFeriado :execrows
UPDATE app.feriado_unidade
SET excluido_em = now()
WHERE empresa_id = sqlc.arg(empresa_id) AND unidade_id = sqlc.arg(unidade_id) AND id = sqlc.arg(id)
  AND excluido_em IS NULL;

-- -------------------------------------------------------------------- perfil

-- name: AdmCriarPerfil :exec
INSERT INTO app.perfil (empresa_id, id, codigo, nome, padrao, perfil_origem, permissoes)
VALUES (sqlc.arg(empresa_id), sqlc.arg(id), sqlc.arg(codigo), sqlc.arg(nome), sqlc.arg(padrao),
        sqlc.narg(perfil_origem), sqlc.arg(permissoes)::text[]);

-- name: AdmListarPerfis :many
SELECT id, codigo, nome, padrao, perfil_origem, permissoes, versao
FROM app.perfil
WHERE empresa_id = sqlc.arg(empresa_id) AND excluido_em IS NULL
ORDER BY padrao DESC, nome, id;

-- name: AdmBuscarPerfisPorIDs :many
SELECT id, codigo, nome, padrao, perfil_origem, permissoes, versao
FROM app.perfil
WHERE empresa_id = sqlc.arg(empresa_id) AND id = ANY (sqlc.arg(ids)::uuid[]) AND excluido_em IS NULL;

-- name: AdmBuscarPerfisPorCodigos :many
SELECT id, codigo, nome, padrao, perfil_origem, permissoes, versao
FROM app.perfil
WHERE empresa_id = sqlc.arg(empresa_id) AND codigo = ANY (sqlc.arg(codigos)::text[]) AND excluido_em IS NULL;

-- Serializa as operações que podem deixar a empresa sem administrador:
-- duas inativações simultâneas esperam uma pela outra aqui.
-- name: AdmTravarPerfilAdministrador :one
SELECT id
FROM app.perfil
WHERE empresa_id = sqlc.arg(empresa_id) AND codigo = 'administrador' AND padrao
FOR UPDATE;

-- name: AdmContarAdministradoresAtivos :one
SELECT count(DISTINCT u.id)::bigint
FROM app.usuario u
JOIN app.usuario_perfil up ON up.empresa_id = u.empresa_id AND up.usuario_id = u.id AND up.excluido_em IS NULL
JOIN app.perfil p ON p.empresa_id = up.empresa_id AND p.id = up.perfil_id AND p.excluido_em IS NULL
WHERE u.empresa_id = sqlc.arg(empresa_id) AND u.situacao = 'ativo' AND u.excluido_em IS NULL
  AND coalesce(p.perfil_origem, p.codigo) = 'administrador'
  AND u.id IS DISTINCT FROM sqlc.narg(exceto_usuario_id)::uuid;

-- ------------------------------------------------------------ usuario_perfil

-- name: AdmPerfisDoUsuario :many
SELECT up.id, up.perfil_id, up.unidade_id, p.codigo, p.nome, p.padrao,
       coalesce(p.perfil_origem, p.codigo)::text AS codigo_origem, p.permissoes
FROM app.usuario_perfil up
JOIN app.perfil p ON p.empresa_id = up.empresa_id AND p.id = up.perfil_id
WHERE up.empresa_id = sqlc.arg(empresa_id) AND up.usuario_id = sqlc.arg(usuario_id)
  AND up.excluido_em IS NULL AND p.excluido_em IS NULL
ORDER BY p.nome;

-- A unicidade de usuario_perfil é parcial (uq_usuario_perfil, só as linhas
-- com excluido_em nulo, NULLS NOT DISTINCT): recolocar um perfil retirado
-- cria outra linha, e o histórico da retirada fica.
-- name: AdmColocarPerfilUsuario :exec
INSERT INTO app.usuario_perfil (empresa_id, id, usuario_id, perfil_id, unidade_id)
VALUES (sqlc.arg(empresa_id), sqlc.arg(id), sqlc.arg(usuario_id), sqlc.arg(perfil_id), sqlc.narg(unidade_id))
ON CONFLICT (empresa_id, usuario_id, perfil_id, unidade_id) WHERE excluido_em IS NULL DO NOTHING;

-- name: AdmRetirarPerfilUsuario :exec
UPDATE app.usuario_perfil
SET excluido_em = now()
WHERE empresa_id = sqlc.arg(empresa_id) AND id = sqlc.arg(id) AND excluido_em IS NULL;

-- --------------------------------------------------------- usuario_permissao

-- name: AdmConcederPermissao :exec
INSERT INTO app.usuario_permissao (empresa_id, id, usuario_id, permissao)
VALUES (sqlc.arg(empresa_id), sqlc.arg(id), sqlc.arg(usuario_id), sqlc.arg(permissao));

-- name: AdmRetirarPermissao :execrows
UPDATE app.usuario_permissao
SET excluido_em = now()
WHERE empresa_id = sqlc.arg(empresa_id) AND usuario_id = sqlc.arg(usuario_id) AND permissao = sqlc.arg(permissao)
  AND excluido_em IS NULL;

-- name: AdmPermissoesAvulsas :many
SELECT id, permissao, criado_em, criado_por
FROM app.usuario_permissao
WHERE empresa_id = sqlc.arg(empresa_id) AND usuario_id = sqlc.arg(usuario_id) AND excluido_em IS NULL
ORDER BY permissao;

-- -------------------------------------------------------------------- pessoa

-- name: AdmCriarPessoa :exec
INSERT INTO app.pessoa (empresa_id, id, nome, cpf, data_nascimento, matricula, cargo, telefone, email, vinculo,
                        fornecedor_id, unidade_id, registro_profissional)
VALUES (sqlc.arg(empresa_id), sqlc.arg(id), sqlc.arg(nome), sqlc.narg(cpf), sqlc.narg(data_nascimento)::text::date,
        sqlc.narg(matricula), sqlc.narg(cargo), sqlc.narg(telefone), sqlc.narg(email), sqlc.arg(vinculo),
        sqlc.narg(fornecedor_id), sqlc.narg(unidade_id), sqlc.narg(registro_profissional));

-- name: AdmBuscarPessoa :one
SELECT p.id, p.nome, p.cpf, coalesce(to_char(p.data_nascimento, 'YYYY-MM-DD'), '')::text AS data_nascimento,
       p.matricula, p.cargo, p.telefone, p.email, p.vinculo, p.fornecedor_id, p.unidade_id,
       p.registro_profissional, p.situacao, p.versao, p.criado_em, p.atualizado_em,
       u.id AS usuario_id
FROM app.pessoa p
LEFT JOIN app.usuario u ON u.empresa_id = p.empresa_id AND u.pessoa_id = p.id AND u.excluido_em IS NULL
WHERE p.empresa_id = sqlc.arg(empresa_id) AND p.id = sqlc.arg(id) AND p.excluido_em IS NULL;

-- name: AdmTravarPessoa :one
SELECT id, fornecedor_id, situacao, versao
FROM app.pessoa
WHERE empresa_id = sqlc.arg(empresa_id) AND id = sqlc.arg(id) AND excluido_em IS NULL
FOR UPDATE;

-- name: AdmAtualizarPessoa :exec
UPDATE app.pessoa
SET nome = sqlc.arg(nome), cpf = sqlc.narg(cpf), data_nascimento = sqlc.narg(data_nascimento)::text::date,
    matricula = sqlc.narg(matricula), cargo = sqlc.narg(cargo), telefone = sqlc.narg(telefone),
    email = sqlc.narg(email), vinculo = sqlc.arg(vinculo), fornecedor_id = sqlc.narg(fornecedor_id),
    unidade_id = sqlc.narg(unidade_id), registro_profissional = sqlc.narg(registro_profissional)
WHERE empresa_id = sqlc.arg(empresa_id) AND id = sqlc.arg(id) AND excluido_em IS NULL;

-- name: AdmDefinirSituacaoPessoa :exec
UPDATE app.pessoa
SET situacao = sqlc.arg(situacao)
WHERE empresa_id = sqlc.arg(empresa_id) AND id = sqlc.arg(id) AND excluido_em IS NULL;

-- Lista por cursor (criado_em, id), do mais novo para o mais antigo. A busca
-- usa o índice GIN de trigramas; o padrão chega já escapado do Go.
-- name: AdmListarPessoas :many
SELECT p.id, p.nome, p.cpf, p.matricula, p.cargo, p.vinculo, p.fornecedor_id, p.unidade_id, p.situacao,
       p.versao, p.criado_em,
       u.id AS usuario_id
FROM app.pessoa p
LEFT JOIN app.usuario u ON u.empresa_id = p.empresa_id AND u.pessoa_id = p.id AND u.excluido_em IS NULL
WHERE p.empresa_id = sqlc.arg(empresa_id) AND p.excluido_em IS NULL
  AND (sqlc.narg(busca)::text IS NULL OR p.nome ILIKE sqlc.narg(busca)::text)
  AND (sqlc.narg(situacao)::text IS NULL OR p.situacao = sqlc.narg(situacao)::text)
  AND (sqlc.narg(fornecedor_id)::uuid IS NULL OR p.fornecedor_id = sqlc.narg(fornecedor_id)::uuid)
  AND (sqlc.narg(unidade_id)::uuid IS NULL OR p.unidade_id = sqlc.narg(unidade_id)::uuid)
  AND (sqlc.narg(com_usuario)::boolean IS NULL OR (u.id IS NOT NULL) = sqlc.narg(com_usuario)::boolean)
  AND (sqlc.narg(cursor_criado_em)::timestamptz IS NULL
       OR (p.criado_em, p.id) < (sqlc.narg(cursor_criado_em)::timestamptz, sqlc.narg(cursor_id)::uuid))
ORDER BY p.criado_em DESC, p.id DESC
LIMIT sqlc.arg(limite);

-- ------------------------------------------------------------------- usuario

-- name: AdmCriarUsuario :exec
INSERT INTO app.usuario (empresa_id, id, pessoa_id, fornecedor_id, locatario_id, nome, email, telefone,
                         custo_hora, unidade_padrao_id, situacao)
VALUES (sqlc.arg(empresa_id), sqlc.arg(id), sqlc.narg(pessoa_id), sqlc.narg(fornecedor_id), sqlc.narg(locatario_id),
        sqlc.arg(nome), sqlc.arg(email), sqlc.narg(telefone),
        (sqlc.narg(custo_hora_centavos)::bigint / 100.0)::numeric(12,2), sqlc.narg(unidade_padrao_id),
        'convidado');

-- name: AdmBuscarUsuario :one
SELECT id, identidade_id, pessoa_id, fornecedor_id, locatario_id, nome, email, telefone,
       custo_hora, unidade_padrao_id, situacao, versao, criado_em,
       atualizado_em
FROM app.usuario
WHERE empresa_id = sqlc.arg(empresa_id) AND id = sqlc.arg(id) AND excluido_em IS NULL;

-- name: AdmTravarUsuario :one
SELECT id, identidade_id, pessoa_id, fornecedor_id, locatario_id, email, situacao, versao
FROM app.usuario
WHERE empresa_id = sqlc.arg(empresa_id) AND id = sqlc.arg(id) AND excluido_em IS NULL
FOR UPDATE;

-- name: AdmAtualizarUsuario :exec
UPDATE app.usuario
SET nome = sqlc.arg(nome), telefone = sqlc.narg(telefone),
    custo_hora = (sqlc.narg(custo_hora_centavos)::bigint / 100.0)::numeric(12,2),
    unidade_padrao_id = sqlc.narg(unidade_padrao_id), pessoa_id = sqlc.narg(pessoa_id)
WHERE empresa_id = sqlc.arg(empresa_id) AND id = sqlc.arg(id) AND excluido_em IS NULL;

-- O próprio usuário (PATCH /v1/usuarios/eu) muda só nome e telefone, inclusive
-- o externo, que não pode tocar nas demais colunas da própria linha
-- (ISO-E01-01).
-- name: AdmAtualizarEu :exec
UPDATE app.usuario
SET nome = sqlc.arg(nome), telefone = sqlc.narg(telefone)
WHERE empresa_id = sqlc.arg(empresa_id) AND id = sqlc.arg(id) AND excluido_em IS NULL;

-- name: AdmDefinirSituacaoUsuario :exec
UPDATE app.usuario
SET situacao = sqlc.arg(situacao)
WHERE empresa_id = sqlc.arg(empresa_id) AND id = sqlc.arg(id) AND excluido_em IS NULL;

-- O convite aceito liga a identidade e ativa o usuário, numa única escrita.
-- name: AdmAceitarConviteUsuario :execrows
UPDATE app.usuario
SET identidade_id = sqlc.arg(identidade_id), situacao = 'ativo'
WHERE empresa_id = sqlc.arg(empresa_id) AND id = sqlc.arg(id) AND situacao = 'convidado'
  AND identidade_id IS NULL AND excluido_em IS NULL;

-- name: AdmBuscarUsuarioDaPessoa :one
SELECT id, fornecedor_id, situacao
FROM app.usuario
WHERE empresa_id = sqlc.arg(empresa_id) AND pessoa_id = sqlc.arg(pessoa_id) AND excluido_em IS NULL;

-- name: AdmListarUsuarios :many
SELECT u.id, u.pessoa_id, u.fornecedor_id, u.locatario_id, u.nome, u.email, u.telefone, u.situacao, u.versao,
       u.criado_em,
       coalesce((SELECT array_agg(coalesce(p.perfil_origem, p.codigo) ORDER BY p.codigo)
                 FROM app.usuario_perfil up
                 JOIN app.perfil p ON p.empresa_id = up.empresa_id AND p.id = up.perfil_id AND p.excluido_em IS NULL
                 WHERE up.empresa_id = u.empresa_id AND up.usuario_id = u.id AND up.excluido_em IS NULL),
                '{}')::text[] AS perfis
FROM app.usuario u
WHERE u.empresa_id = sqlc.arg(empresa_id) AND u.excluido_em IS NULL
  AND (sqlc.narg(busca)::text IS NULL OR u.nome ILIKE sqlc.narg(busca)::text OR u.email ILIKE sqlc.narg(busca)::text)
  AND (sqlc.narg(situacao)::text IS NULL OR u.situacao = sqlc.narg(situacao)::text)
  AND (sqlc.narg(perfil)::text IS NULL OR EXISTS (
        SELECT 1 FROM app.usuario_perfil up
        JOIN app.perfil p ON p.empresa_id = up.empresa_id AND p.id = up.perfil_id AND p.excluido_em IS NULL
        WHERE up.empresa_id = u.empresa_id AND up.usuario_id = u.id AND up.excluido_em IS NULL
          AND coalesce(p.perfil_origem, p.codigo) = sqlc.narg(perfil)::text))
  AND (sqlc.narg(cursor_criado_em)::timestamptz IS NULL
       OR (u.criado_em, u.id) < (sqlc.narg(cursor_criado_em)::timestamptz, sqlc.narg(cursor_id)::uuid))
ORDER BY u.criado_em DESC, u.id DESC
LIMIT sqlc.arg(limite);

-- --------------------------------------------------------------- auditoria

-- Trilha com filtros (RF-ADM-011). Ordem (ocorrido_em, id) decrescente, a da
-- chave primária; o intervalo de datas permite ao planejador descartar as
-- partições fora do período.
-- name: AdmListarEventosAuditoria :many
SELECT id, ocorrido_em, ator_tipo, usuario_id, token_api_id, acao, entidade, entidade_id, entidade_chave,
       campos_alterados, origem, coalesce(host(ip), '')::text AS ip, dispositivo_id, request_id, conexao_empresa_id,
       ocorrido_no_dispositivo_em
FROM app.evento_auditoria
WHERE empresa_id = sqlc.arg(empresa_id)
  AND (sqlc.narg(entidade)::text IS NULL OR entidade = sqlc.narg(entidade)::text)
  AND (sqlc.narg(entidade_id)::uuid IS NULL OR entidade_id = sqlc.narg(entidade_id)::uuid)
  AND (sqlc.narg(usuario_id)::uuid IS NULL OR usuario_id = sqlc.narg(usuario_id)::uuid)
  AND (sqlc.narg(acao)::text IS NULL OR acao = sqlc.narg(acao)::text)
  AND (sqlc.narg(de)::timestamptz IS NULL OR ocorrido_em >= sqlc.narg(de)::timestamptz)
  AND (sqlc.narg(ate)::timestamptz IS NULL OR ocorrido_em < sqlc.narg(ate)::timestamptz)
  AND (sqlc.narg(cursor_ocorrido_em)::timestamptz IS NULL
       OR (ocorrido_em, id) < (sqlc.narg(cursor_ocorrido_em)::timestamptz, sqlc.narg(cursor_id)::uuid))
ORDER BY ocorrido_em DESC, id DESC
LIMIT sqlc.arg(limite);

-- name: AdmBuscarEventoAuditoria :one
SELECT id, ocorrido_em, ator_tipo, usuario_id, token_api_id, acao, entidade, entidade_id, entidade_chave,
       antes, depois, campos_alterados, origem, coalesce(host(ip), '')::text AS ip, dispositivo_id, user_agent, request_id,
       conexao_empresa_id, ocorrido_no_dispositivo_em
FROM app.evento_auditoria
WHERE empresa_id = sqlc.arg(empresa_id) AND id = sqlc.arg(id);
