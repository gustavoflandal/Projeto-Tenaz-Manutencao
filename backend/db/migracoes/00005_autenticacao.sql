-- Autenticação: sessão, refresh token, código OAuth, desafio de MFA e token
-- de e-mail.
--
-- Fonte: docs/2_ARQUITETURA_BANCO_DADOS.md, seção 5.1; docs/3_ESPECIFICACAO_APIS.md,
-- seções 4.2, 4.3 e 4.7. Requisitos: RF-ADM-004, RF-ADM-006, RNF-017.
-- Nomes aprovados pelo dono do produto na etapa E01 (sumário 0.6, seção 3.2).
--
-- Esquema plataforma, sem RLS: estas tabelas são consultadas antes de a
-- empresa ser conhecida (o cookie, o refresh token e o código chegam sem
-- empresa). Sem RLS, o privilégio é a defesa (ADR-011, documento 5, seções 6.9
-- e 15; achado ISO-E01-02):
--
--   1. Só o papel tenaz_autenticacao, usado só pelo pacote de autenticação
--      num pool próprio, alcança as sete tabelas de credencial (identidade e
--      acesso_empresa, da 00004, e as cinco desta). tenaz_app, tenaz_job,
--      tenaz_leitura e tenaz_sinc recebem "permissão negada". O tenaz_job só
--      executa a função de purga.
--   2. Nenhum segredo é guardado: só o sha256 dele.
--   3. A empresa de uma credencial é fixada na criação: o tenaz_autenticacao
--      não tem UPDATE em empresa_id, usuario_id, identidade_id, canal nem nos
--      hashes. Trocar de empresa é criar outra sessão (documento 3, seção 4.7).
--   4. FK composta para app.usuario (empresa_id, id, identidade_id): a sessão
--      só existe para a identidade que é de fato a do usuário naquela empresa.

-- +goose Up

-- +goose StatementBegin
-- Papel das credenciais (ADR-011). Nasce aqui, e não na 00001, que já está em
-- main e não se edita. Roda como quem aplica a migração, que cria papéis; o
-- RESET é necessário porque o SET ROLE da migração anterior sobrevive na
-- conexão do goose. Guarda de existência como na 00001: em provedor
-- gerenciado, os papéis vêm da infraestrutura como código.
--
-- NOLOGIN NOBYPASSRLS, sem objetos, sem relação de membro com tenaz_dono,
-- tenaz_app ou tenaz_job. A conexão usa tenaz_autenticacao_conexao, que herda
-- o papel, no mesmo par de tenaz_sinc e tenaz_sinc_conexao.
RESET ROLE;
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'tenaz_autenticacao') THEN
    CREATE ROLE tenaz_autenticacao NOLOGIN NOBYPASSRLS;
  END IF;
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'tenaz_autenticacao_conexao') THEN
    CREATE ROLE tenaz_autenticacao_conexao LOGIN NOBYPASSRLS IN ROLE tenaz_autenticacao;
  END IF;
END
$$;
-- O tipo citext de identidade.email e token_email.email mora em public (00004).
GRANT USAGE ON SCHEMA public TO tenaz_autenticacao;
-- +goose StatementEnd

-- +goose StatementBegin
-- As tabelas de credencial pertencem a tenaz_credencial, e não a tenaz_dono
-- (achado ISO-E01-08, migração 00004): sem RLS, o dono e quem herda dele
-- (tenaz_manutencao) leriam todas as credenciais. Os objetos de app tocados
-- aqui (policies e GRANTs de evento_auditoria) continuam com SET ROLE tenaz_dono.
SET ROLE tenaz_credencial;
SET LOCAL lock_timeout = '3s';
SET LOCAL statement_timeout = '60s';
-- +goose StatementEnd

-- +goose StatementBegin
-- Sessão web ou do app (documento 3, seções 4.2 e 4.3).
--
-- Web: o cookie __Host-tenaz_sessao leva um identificador aleatório de 256
-- bits; aqui fica só o sha256 (token_hash). O token CSRF também é guardado
-- como hash e comparado em tempo constante. Expira por inatividade (12 h, ou
-- menos por parâmetro da empresa, conferido pelo serviço contra ultimo_uso_em)
-- e no máximo em 7 dias (expira_em, com teto no CHECK).
--
-- App: a sessão é a do aparelho; o id dela é o claim sid do access token, e
-- ela é a família dos refresh tokens (token_renovacao.familia_id). Não tem
-- cookie nem CSRF. expira_em acompanha o refresh token vigente: 30 dias de
-- inatividade, empurrado a cada renovação (RNF-017).
CREATE TABLE plataforma.sessao (
  id              uuid        NOT NULL DEFAULT app.uuid_v7(),
  canal           text        NOT NULL,
  token_hash      bytea,                    -- web: sha256 do identificador do cookie
  csrf_hash       bytea,                    -- web: sha256 do token CSRF
  identidade_id   uuid        NOT NULL,
  empresa_id      uuid        NOT NULL,
  usuario_id      uuid        NOT NULL,
  dispositivo_id  uuid,                     -- app: aparelho; FK quando app.dispositivo existir
  ip              inet,
  user_agent      text,
  criado_em       timestamptz NOT NULL DEFAULT now(),
  ultimo_uso_em   timestamptz NOT NULL DEFAULT now(),
  expira_em       timestamptz NOT NULL,
  revogada_em     timestamptz,
  CONSTRAINT pk_sessao PRIMARY KEY (id),
  CONSTRAINT fk_sessao_identidade FOREIGN KEY (identidade_id)
    REFERENCES plataforma.identidade (id) ON DELETE CASCADE,
  CONSTRAINT fk_sessao_usuario FOREIGN KEY (empresa_id, usuario_id, identidade_id)
    REFERENCES app.usuario (empresa_id, id, identidade_id) ON DELETE CASCADE,
  CONSTRAINT ck_sessao_canal CHECK (canal IN ('web', 'app')),
  -- cookie e CSRF existem só na web
  CONSTRAINT ck_sessao_web CHECK ((canal = 'web') = (token_hash IS NOT NULL AND csrf_hash IS NOT NULL)),
  CONSTRAINT ck_sessao_hash CHECK (octet_length(token_hash) = 32 AND octet_length(csrf_hash) = 32),
  -- teto de 7 dias da sessão web, mesmo que o serviço erre a conta
  CONSTRAINT ck_sessao_validade CHECK (expira_em > criado_em
                                       AND (canal <> 'web' OR expira_em <= criado_em + interval '7 days'))
);
-- Cada requisição web acha a sessão pelo hash do cookie.
CREATE UNIQUE INDEX uq_sessao_token ON plataforma.sessao (token_hash) WHERE token_hash IS NOT NULL;
-- "Encerrar sessões" do usuário (RF-ADM-004, RF-ADM-006) e inativação.
CREATE INDEX ix_sessao_usuario ON plataforma.sessao (empresa_id, usuario_id) WHERE revogada_em IS NULL;
-- Troca de senha e bloqueio revogam todas as sessões da identidade, em todas as empresas.
CREATE INDEX ix_sessao_identidade ON plataforma.sessao (identidade_id) WHERE revogada_em IS NULL;
-- Purga das vencidas pelo job.
CREATE INDEX ix_sessao_expira ON plataforma.sessao (expira_em);
-- +goose StatementEnd

-- +goose StatementBegin
-- Refresh token rotativo do app (documento 3, seção 4.3; RNF-017).
--
-- Cada renovação marca usado_em no token apresentado e cria o sucessor, com o
-- mesmo familia_id e anterior_id apontando para ele. Reapresentar um token já
-- usado é sinal de roubo: o serviço revoga a família inteira — todos os tokens
-- com o mesmo familia_id e a sessão de id familia_id — e o usuário entra de
-- novo. A família é a sessão do app: familia_id referencia plataforma.sessao.
--
-- A rotação é atômica:
--   UPDATE plataforma.token_renovacao SET usado_em = now()
--   WHERE token_hash = $1 AND usado_em IS NULL AND revogado_em IS NULL AND expira_em > now()
--   RETURNING familia_id;
-- Zero linhas com o hash existente = reuso, expiração ou revogação.
CREATE TABLE plataforma.token_renovacao (
  id           uuid        NOT NULL DEFAULT app.uuid_v7(),
  familia_id   uuid        NOT NULL,       -- a sessão do app
  anterior_id  uuid,                       -- o token que este substituiu; nulo no primeiro
  token_hash   bytea       NOT NULL,       -- sha256 do refresh token (tnz_rt_...)
  criado_em    timestamptz NOT NULL DEFAULT now(),
  expira_em    timestamptz NOT NULL,       -- criado_em + 30 dias de inatividade
  usado_em     timestamptz,                -- trocado pelo sucessor
  revogado_em  timestamptz,
  CONSTRAINT pk_token_renovacao PRIMARY KEY (id),
  CONSTRAINT uq_token_renovacao_hash UNIQUE (token_hash),
  CONSTRAINT fk_token_renovacao_familia FOREIGN KEY (familia_id)
    REFERENCES plataforma.sessao (id) ON DELETE CASCADE,
  -- SET NULL, e não CASCADE: apagar um token antigo nunca pode levar junto o
  -- sucessor vivo. A purga normal é pela sessão, que leva a família inteira.
  CONSTRAINT fk_token_renovacao_anterior FOREIGN KEY (anterior_id)
    REFERENCES plataforma.token_renovacao (id) ON DELETE SET NULL,
  CONSTRAINT ck_token_renovacao_hash CHECK (octet_length(token_hash) = 32),
  CONSTRAINT ck_token_renovacao_validade CHECK (expira_em > criado_em
                                                AND expira_em <= criado_em + interval '30 days')
);
-- Revogação da família inteira, na detecção de reuso.
CREATE INDEX ix_token_renovacao_familia ON plataforma.token_renovacao (familia_id);
-- FK para o anterior: apagar um token acha o sucessor sem varredura.
CREATE INDEX ix_token_renovacao_anterior ON plataforma.token_renovacao (anterior_id) WHERE anterior_id IS NOT NULL;
-- +goose StatementEnd

-- +goose StatementBegin
-- Código de autorização OAuth 2.1 do app, com PKCE S256 (documento 3, seção
-- 4.3). Uso único e vida curta. A empresa é escolhida na tela de autorização,
-- antes de o código existir, e fica fixa nele.
--
-- desafio_pkce é o code_challenge: base64url sem preenchimento de um sha256,
-- 43 caracteres. Só S256 é aceito; o método não é guardado porque não varia.
-- Código reapresentado depois de usado revoga a sessão que ele gerou
-- (sessao_id), como recomenda o OAuth 2.1.
CREATE TABLE plataforma.codigo_autorizacao (
  id             uuid        NOT NULL DEFAULT app.uuid_v7(),
  codigo_hash    bytea       NOT NULL,     -- sha256 do código
  identidade_id  uuid        NOT NULL,
  empresa_id     uuid        NOT NULL,
  usuario_id     uuid        NOT NULL,
  cliente_oauth  text        NOT NULL,     -- client_id
  uri_retorno    text        NOT NULL,     -- redirect_uri, conferida byte a byte na troca
  desafio_pkce   text        NOT NULL,     -- code_challenge (S256)
  criado_em      timestamptz NOT NULL DEFAULT now(),
  expira_em      timestamptz NOT NULL,
  usado_em       timestamptz,
  sessao_id      uuid,                     -- sessão criada na troca
  CONSTRAINT pk_codigo_autorizacao PRIMARY KEY (id),
  CONSTRAINT uq_codigo_autorizacao_hash UNIQUE (codigo_hash),
  CONSTRAINT fk_codigo_autorizacao_identidade FOREIGN KEY (identidade_id)
    REFERENCES plataforma.identidade (id) ON DELETE CASCADE,
  CONSTRAINT fk_codigo_autorizacao_usuario FOREIGN KEY (empresa_id, usuario_id, identidade_id)
    REFERENCES app.usuario (empresa_id, id, identidade_id) ON DELETE CASCADE,
  CONSTRAINT fk_codigo_autorizacao_sessao FOREIGN KEY (sessao_id)
    REFERENCES plataforma.sessao (id) ON DELETE SET NULL,
  CONSTRAINT ck_codigo_autorizacao_hash CHECK (octet_length(codigo_hash) = 32),
  CONSTRAINT ck_codigo_autorizacao_cliente CHECK (cliente_oauth IN ('tenaz-app')),
  CONSTRAINT ck_codigo_autorizacao_pkce CHECK (desafio_pkce ~ '^[A-Za-z0-9_-]{43}$'),
  -- teto de 10 minutos (OAuth 2.1); o serviço usa 60 segundos
  CONSTRAINT ck_codigo_autorizacao_validade CHECK (expira_em > criado_em
                                                   AND expira_em <= criado_em + interval '10 minutes'),
  CONSTRAINT ck_codigo_autorizacao_uso CHECK (sessao_id IS NULL OR usado_em IS NOT NULL)
);
CREATE INDEX ix_codigo_autorizacao_sessao ON plataforma.codigo_autorizacao (sessao_id) WHERE sessao_id IS NOT NULL;
CREATE INDEX ix_codigo_autorizacao_expira ON plataforma.codigo_autorizacao (expira_em);
-- +goose StatementEnd

-- +goose StatementBegin
-- Desafio de MFA: o passo entre a senha certa e o TOTP (documento 3, seção
-- 4.2). Ainda não há empresa: a escolha vem depois do segundo fator.
--
-- O desafio_id da API é um segredo de 256 bits; aqui fica o sha256. Cada
-- código errado incrementa tentativas no mesmo UPDATE que confere o limite, e
-- o CHECK é o teto: a sexta tentativa não grava.
CREATE TABLE plataforma.desafio_mfa (
  id             uuid        NOT NULL DEFAULT app.uuid_v7(),
  token_hash     bytea       NOT NULL,
  identidade_id  uuid        NOT NULL,
  canal          text        NOT NULL,
  tentativas     smallint    NOT NULL DEFAULT 0,
  ip             inet,
  criado_em      timestamptz NOT NULL DEFAULT now(),
  expira_em      timestamptz NOT NULL,
  concluido_em   timestamptz,
  CONSTRAINT pk_desafio_mfa PRIMARY KEY (id),
  CONSTRAINT uq_desafio_mfa_hash UNIQUE (token_hash),
  CONSTRAINT fk_desafio_mfa_identidade FOREIGN KEY (identidade_id)
    REFERENCES plataforma.identidade (id) ON DELETE CASCADE,
  CONSTRAINT ck_desafio_mfa_hash CHECK (octet_length(token_hash) = 32),
  CONSTRAINT ck_desafio_mfa_canal CHECK (canal IN ('web', 'app')),
  CONSTRAINT ck_desafio_mfa_tentativas CHECK (tentativas BETWEEN 0 AND 5),
  -- teto de 10 minutos; o serviço usa 5
  CONSTRAINT ck_desafio_mfa_validade CHECK (expira_em > criado_em
                                            AND expira_em <= criado_em + interval '10 minutes')
);
CREATE INDEX ix_desafio_mfa_identidade ON plataforma.desafio_mfa (identidade_id);
CREATE INDEX ix_desafio_mfa_expira ON plataforma.desafio_mfa (expira_em);
-- +goose StatementEnd

-- +goose StatementBegin
-- Link enviado por e-mail (Mailpit em desenvolvimento): convite de usuário,
-- redefinição de senha e verificação de e-mail. Só o hash, uso único,
-- expiração com teto por finalidade.
--
-- O convite é de uma empresa e de um usuário (que nasce `convidado` e ainda
-- não tem identidade); redefinição e verificação são da identidade, sem
-- empresa.
CREATE TABLE plataforma.token_email (
  id             uuid        NOT NULL DEFAULT app.uuid_v7(),
  token_hash     bytea       NOT NULL,
  finalidade     text        NOT NULL,
  email          citext      NOT NULL,     -- destino no momento do envio
  identidade_id  uuid,
  empresa_id     uuid,
  usuario_id     uuid,
  criado_em      timestamptz NOT NULL DEFAULT now(),
  expira_em      timestamptz NOT NULL,
  usado_em       timestamptz,
  revogado_em    timestamptz,              -- novo envio revoga o anterior
  CONSTRAINT pk_token_email PRIMARY KEY (id),
  CONSTRAINT uq_token_email_hash UNIQUE (token_hash),
  CONSTRAINT fk_token_email_identidade FOREIGN KEY (identidade_id)
    REFERENCES plataforma.identidade (id) ON DELETE CASCADE,
  CONSTRAINT fk_token_email_usuario FOREIGN KEY (empresa_id, usuario_id)
    REFERENCES app.usuario (empresa_id, id) ON DELETE CASCADE,
  CONSTRAINT ck_token_email_hash CHECK (octet_length(token_hash) = 32),
  CONSTRAINT ck_token_email_finalidade CHECK (finalidade IN ('convite', 'redefinicao_senha', 'verificacao_email')),
  -- convite: empresa e usuário; demais: identidade, sem empresa
  CONSTRAINT ck_token_email_alvo CHECK (
    CASE finalidade
      WHEN 'convite' THEN empresa_id IS NOT NULL AND usuario_id IS NOT NULL
      ELSE identidade_id IS NOT NULL AND empresa_id IS NULL AND usuario_id IS NULL
    END),
  -- tetos (proposta): convite 7 dias, verificação 24 horas, redefinição 1 hora
  CONSTRAINT ck_token_email_validade CHECK (expira_em > criado_em AND expira_em <= criado_em +
    CASE finalidade
      WHEN 'convite'           THEN interval '7 days'
      WHEN 'verificacao_email' THEN interval '24 hours'
      ELSE                          interval '1 hour'
    END)
);
-- Novo convite revoga o anterior do mesmo usuário; "reenviar convite".
CREATE INDEX ix_token_email_usuario ON plataforma.token_email (empresa_id, usuario_id) WHERE usuario_id IS NOT NULL;
-- Nova redefinição revoga a anterior da mesma identidade.
CREATE INDEX ix_token_email_identidade ON plataforma.token_email (identidade_id) WHERE identidade_id IS NOT NULL;
CREATE INDEX ix_token_email_expira ON plataforma.token_email (expira_em);
-- +goose StatementEnd

-- +goose StatementBegin
-- Privilégios das credenciais (ADR-011; documento 5, seção 6.9). Primeiro tira
-- tudo de todos, depois dá ao tenaz_autenticacao exatamente o necessário:
-- SELECT e INSERT, UPDATE só nas colunas de ciclo de vida (uso, consumo,
-- revogação), nenhum DELETE nem TRUNCATE.
-- esquemas e funções de apoio: concedidos pelo dono deles
SET ROLE tenaz_dono;
GRANT USAGE ON SCHEMA app, plataforma TO tenaz_autenticacao;
GRANT EXECUTE ON FUNCTION app.uuid_v7(), app.empresa_atual() TO tenaz_autenticacao;
-- as sete tabelas: concedidas pelo dono delas
SET ROLE tenaz_credencial;

REVOKE ALL ON plataforma.identidade, plataforma.acesso_empresa, plataforma.sessao, plataforma.token_renovacao,
              plataforma.codigo_autorizacao, plataforma.desafio_mfa, plataforma.token_email
  FROM PUBLIC, tenaz_app, tenaz_job, tenaz_leitura, tenaz_sinc;

GRANT SELECT, INSERT ON plataforma.identidade TO tenaz_autenticacao;
GRANT UPDATE (email, senha_hash, mfa_segredo_cifrado, mfa_ativado_em, mfa_ultimo_passo, email_verificado_em, falhas_login,
              bloqueada_ate, bloqueada_em, ultimo_login_em, atualizado_em)
  ON plataforma.identidade TO tenaz_autenticacao;               -- id e criado_em nunca mudam
GRANT SELECT ON plataforma.acesso_empresa TO tenaz_autenticacao; -- a escrita é só do trigger
GRANT SELECT, INSERT ON plataforma.sessao, plataforma.token_renovacao, plataforma.codigo_autorizacao,
                        plataforma.desafio_mfa, plataforma.token_email TO tenaz_autenticacao;
GRANT UPDATE (ultimo_uso_em, expira_em, revogada_em) ON plataforma.sessao TO tenaz_autenticacao;
GRANT UPDATE (usado_em, revogado_em) ON plataforma.token_renovacao TO tenaz_autenticacao;
GRANT UPDATE (usado_em, sessao_id) ON plataforma.codigo_autorizacao TO tenaz_autenticacao;
GRANT UPDATE (tentativas, concluido_em) ON plataforma.desafio_mfa TO tenaz_autenticacao;
GRANT UPDATE (usado_em, revogado_em) ON plataforma.token_email TO tenaz_autenticacao;
-- +goose StatementEnd

-- +goose StatementBegin
SET ROLE tenaz_dono;
-- Auditoria de autenticação (documento 2, seção 9.1): login, falha, logout,
-- troca de senha e MFA, gravados pelo próprio tenaz_autenticacao em cada
-- empresa da identidade, na mesma transação da credencial. No molde do
-- tenaz_sinc: INSERT por coluna e nenhuma leitura da trilha.
GRANT INSERT (empresa_id, ator_tipo, usuario_id, acao, entidade, entidade_id, depois, campos_alterados,
              origem, ip, dispositivo_id, user_agent, request_id)
  ON app.evento_auditoria TO tenaz_autenticacao;

-- Permissiva: a empresa do evento é a da transação, como para os demais.
ALTER POLICY auditoria_insercao ON app.evento_auditoria
  TO tenaz_app, tenaz_job, tenaz_sinc, tenaz_autenticacao;

-- Restritiva: só eventos da identidade, só as quatro ações, e o usuário do
-- evento tem de ser, em acesso_empresa, o usuário daquela identidade naquela
-- empresa. O papel não grava nenhum outro tipo de evento nem atribui o evento
-- a outra pessoa. A consulta é pela chave primária de acesso_empresa.
--
-- Interação com as outras restritivas de evento_auditoria: escopo_externo_
-- insercao vale só para tenaz_app e tenaz_job, e conexao_sinc só para
-- tenaz_sinc; nenhuma alcança o tenaz_autenticacao, que grava sem escopo
-- definido. Ele não é escopo: é contido por esta policy e pelas colunas.
CREATE POLICY autenticacao_identidade ON app.evento_auditoria AS RESTRICTIVE
  FOR INSERT TO tenaz_autenticacao
  WITH CHECK (ator_tipo = 'usuario' AND entidade = 'identidade'
              AND acao IN ('login', 'login_falho', 'logout', 'alteracao')
              AND EXISTS (SELECT 1 FROM plataforma.acesso_empresa a
                          WHERE a.identidade_id = evento_auditoria.entidade_id
                            AND a.empresa_id = evento_auditoria.empresa_id
                            AND a.usuario_id = evento_auditoria.usuario_id));
-- +goose StatementEnd

-- +goose StatementBegin
SET ROLE tenaz_credencial;
-- Manutenção de acesso_empresa, agora com revogação (ADR-011, decisão 6).
-- Mesmo upsert e mesmo search_path da 00004. Quando o acesso fica desabilitado
-- (usuário que deixa de estar ativo, ou excluído), revoga na mesma transação as
-- sessões abertas dele naquela empresa e os refresh tokens das famílias
-- revogadas: a inativação corta as credenciais de forma atômica, sem que o ADM
-- toque tabela de credencial. app.conferir_escopo() continua recusando o
-- usuário inativo na requisição seguinte, como segunda barreira.
CREATE OR REPLACE FUNCTION plataforma.tg_acesso_empresa() RETURNS trigger
  LANGUAGE plpgsql SECURITY DEFINER
  SET search_path = pg_catalog, pg_temp
AS $$
DECLARE
  v_habilitado boolean := NEW.situacao = 'ativo' AND NEW.excluido_em IS NULL;
BEGIN
  INSERT INTO plataforma.acesso_empresa AS a (identidade_id, empresa_id, usuario_id, habilitado, atualizado_em)
  VALUES (NEW.identidade_id, NEW.empresa_id, NEW.id, v_habilitado, pg_catalog.now())
  ON CONFLICT (identidade_id, empresa_id) DO UPDATE
    SET usuario_id    = EXCLUDED.usuario_id,
        habilitado    = EXCLUDED.habilitado,
        atualizado_em = EXCLUDED.atualizado_em;

  IF NOT v_habilitado THEN
    WITH revogadas AS (
      UPDATE plataforma.sessao s
      SET revogada_em = pg_catalog.now()
      WHERE s.empresa_id = NEW.empresa_id AND s.usuario_id = NEW.id AND s.revogada_em IS NULL
      RETURNING s.id
    )
    UPDATE plataforma.token_renovacao t
    SET revogado_em = pg_catalog.now()
    WHERE t.familia_id IN (SELECT r.id FROM revogadas r) AND t.revogado_em IS NULL;
  END IF;
  RETURN NULL;
END
$$;
-- +goose StatementEnd

-- +goose StatementBegin
-- Purga das credenciais vencidas há mais de 30 dias (proposta: o prazo da
-- investigação de incidente). SECURITY DEFINER e sem parâmetro: não recebe
-- nenhum dado da requisição, por isso é o caso certo para uma função estreita.
-- O tenaz_job a executa por um job diário e não tem nenhum privilégio nas
-- tabelas. A ordem respeita as FKs; a sessão leva por cascata os tokens que
-- ainda restarem da família. Devolve o total apagado.
CREATE OR REPLACE FUNCTION plataforma.purgar_credenciais_vencidas() RETURNS bigint
  LANGUAGE plpgsql VOLATILE SECURITY DEFINER
  SET search_path = pg_catalog, pg_temp
AS $$
DECLARE
  v_limite timestamptz := pg_catalog.now() - interval '30 days';
  v_n      bigint;
  v_total  bigint := 0;
BEGIN
  DELETE FROM plataforma.codigo_autorizacao WHERE expira_em < v_limite;
  GET DIAGNOSTICS v_n = ROW_COUNT; v_total := v_total + v_n;
  DELETE FROM plataforma.desafio_mfa WHERE expira_em < v_limite;
  GET DIAGNOSTICS v_n = ROW_COUNT; v_total := v_total + v_n;
  DELETE FROM plataforma.token_email WHERE expira_em < v_limite;
  GET DIAGNOSTICS v_n = ROW_COUNT; v_total := v_total + v_n;
  DELETE FROM plataforma.token_renovacao WHERE expira_em < v_limite;
  GET DIAGNOSTICS v_n = ROW_COUNT; v_total := v_total + v_n;
  DELETE FROM plataforma.sessao WHERE expira_em < v_limite;
  GET DIAGNOSTICS v_n = ROW_COUNT; v_total := v_total + v_n;
  RETURN v_total;
END
$$;
-- +goose StatementEnd

-- +goose StatementBegin
REVOKE ALL ON FUNCTION plataforma.purgar_credenciais_vencidas() FROM PUBLIC;
GRANT EXECUTE ON FUNCTION plataforma.purgar_credenciais_vencidas() TO tenaz_job;
-- +goose StatementEnd

-- +goose StatementBegin
-- Volta ao dono comum no fim: o goose registra a versão na tabela de controle,
-- que pertence a tenaz_dono, com o papel corrente da conexão.
SET ROLE tenaz_dono;
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
SET ROLE tenaz_credencial;
DROP FUNCTION IF EXISTS plataforma.purgar_credenciais_vencidas();

-- tg_acesso_empresa volta à versão da 00004, sem a revogação (as tabelas
-- revogadas deixam de existir logo abaixo).
CREATE OR REPLACE FUNCTION plataforma.tg_acesso_empresa() RETURNS trigger
  LANGUAGE plpgsql SECURITY DEFINER
  SET search_path = pg_catalog, pg_temp
AS $f$
BEGIN
  INSERT INTO plataforma.acesso_empresa AS a (identidade_id, empresa_id, usuario_id, habilitado, atualizado_em)
  VALUES (NEW.identidade_id, NEW.empresa_id, NEW.id,
          NEW.situacao = 'ativo' AND NEW.excluido_em IS NULL, pg_catalog.now())
  ON CONFLICT (identidade_id, empresa_id) DO UPDATE
    SET usuario_id    = EXCLUDED.usuario_id,
        habilitado    = EXCLUDED.habilitado,
        atualizado_em = EXCLUDED.atualizado_em;
  RETURN NULL;
END
$f$;

REVOKE ALL ON plataforma.identidade, plataforma.acesso_empresa FROM tenaz_autenticacao;

SET ROLE tenaz_dono;
DROP POLICY IF EXISTS autenticacao_identidade ON app.evento_auditoria;
ALTER POLICY auditoria_insercao ON app.evento_auditoria TO tenaz_app, tenaz_job, tenaz_sinc;
REVOKE ALL ON app.evento_auditoria FROM tenaz_autenticacao;
REVOKE EXECUTE ON FUNCTION app.uuid_v7(), app.empresa_atual() FROM tenaz_autenticacao;
REVOKE USAGE ON SCHEMA app, plataforma FROM tenaz_autenticacao;

SET ROLE tenaz_credencial;

DROP TABLE IF EXISTS plataforma.token_email;
DROP TABLE IF EXISTS plataforma.desafio_mfa;
DROP TABLE IF EXISTS plataforma.codigo_autorizacao;
DROP TABLE IF EXISTS plataforma.token_renovacao;
DROP TABLE IF EXISTS plataforma.sessao;

-- O papel não é removido: é objeto do cluster e pode ter vindo da
-- infraestrutura, como os da 00001.
RESET ROLE;
REVOKE USAGE ON SCHEMA public FROM tenaz_autenticacao;
-- +goose StatementEnd
