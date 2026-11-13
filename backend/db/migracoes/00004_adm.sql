-- Administração: empresa, unidade e calendário, identidade, pessoa, usuário,
-- perfil e as travas de escopo externo.
--
-- Fonte: docs/2_ARQUITETURA_BANCO_DADOS.md, seções 2.3, 2.9.3, 2.9.4 e 5.1.
-- Requisitos: RF-ADM-001, RF-ADM-002, RF-ADM-004, RF-ADM-005, RF-ADM-015.
--
-- Os 12 perfis padrão NÃO são semeados aqui: o serviço de criação de empresa
-- os cria junto com a empresa, a partir do catálogo de permissões versionado
-- no código (seção 14.2). Esta migração garante só os CHECKs que nenhum
-- perfil pode violar.
--
-- Ficam sem FK, por enquanto, pessoa.fornecedor_id, usuario.fornecedor_id e
-- usuario.locatario_id: fornecedor nasce em E06 e locatario na fase 2. As FKs
-- entram como NOT VALID + VALIDATE na migração de cada tabela (seção 17).

-- +goose Up

-- +goose StatementBegin
-- Correção da 00001, que não se edita depois de aplicada: ela fez
-- `REVOKE ALL ON SCHEMA public FROM PUBLIC`, o que tirou também o USAGE, e as
-- extensões (citext, pg_trgm, btree_gin, btree_gist) moram em public. Sem
-- USAGE, tenaz_dono não acha o tipo citext ao criar a tabela — e, o que é
-- pior, tenaz_app resolveria `email = $1` como text = text, pelo cast
-- implícito, e o login passaria a diferenciar maiúsculas em silêncio.
--
-- A intenção original (CVE-2018-1058) é tirar o CREATE de PUBLIC, e isso
-- continua: aqui volta só o USAGE, e só para os papéis do sistema, não para
-- PUBLIC. Roda como quem aplica a migração, dono do esquema public; por isso
-- vem antes do SET ROLE. O RESET é necessário porque o SET ROLE da migração
-- anterior sobrevive na conexão do goose.
RESET ROLE;
GRANT USAGE ON SCHEMA public TO tenaz_dono, tenaz_app, tenaz_job, tenaz_leitura, tenaz_sinc;
-- +goose StatementEnd

-- +goose StatementBegin
-- Dono das credenciais (achado ISO-E01-08). As tabelas de autenticação de
-- plataforma não têm RLS (são lidas antes de a empresa ser conhecida), então
-- quem é dono delas, ou membro com herança do dono, lê tudo, inclusive
-- senha_hash e segredo de MFA. Se o dono fosse tenaz_dono, o tenaz_manutencao
-- (LOGIN, membro de tenaz_dono com herança, para partições e VACUUM) leria
-- todas as credenciais.
--
-- Por isso as sete tabelas de credencial e as funções SECURITY DEFINER que as
-- tocam pertencem a tenaz_credencial: NOLOGIN, NOBYPASSRLS, sem membros além
-- de quem aplica migração (o pipeline, como já acontece com tenaz_dono). Nem
-- tenaz_dono nem tenaz_manutencao são membros dele, em nenhum sentido.
--
-- RLS forçado com policy só para tenaz_autenticacao (caminho a) foi descartado:
-- as funções SECURITY DEFINER do dono seriam barradas, e dar policy ao dono
-- reabriria o acesso ao tenaz_manutencao, que herda do dono. Tirar a herança do
-- tenaz_manutencao (caminho c) não resolve: o dono continuaria lendo, e
-- qualquer membro com SET ROLE também.
DO $$
BEGIN
  IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'tenaz_credencial') THEN
    CREATE ROLE tenaz_credencial NOLOGIN NOBYPASSRLS;
  END IF;
  -- Quem aplica a migração precisa assumir o papel para criar e alterar as tabelas.
  IF NOT pg_has_role(current_user, 'tenaz_credencial', 'MEMBER') THEN
    EXECUTE format('GRANT tenaz_credencial TO %I', current_user);
  END IF;
END
$$;
-- Cria tabelas em plataforma; usa citext (public) e app.uuid_v7(); a FK de
-- acesso_empresa resolve app.usuario.
GRANT USAGE, CREATE ON SCHEMA plataforma TO tenaz_credencial;
GRANT USAGE ON SCHEMA public, app TO tenaz_credencial;
-- +goose StatementEnd

-- +goose StatementBegin
SET ROLE tenaz_dono;
SET LOCAL lock_timeout = '3s';
SET LOCAL statement_timeout = '60s';
-- +goose StatementEnd

-- ------------------------------------------------------------------ empresa

-- +goose StatementBegin
-- Empresa cliente: unidade de isolamento (tenant). Não tem empresa_id, e sim
-- id; está na lista de tabelas globais (seção 2.11) e tem policies próprias.
CREATE TABLE app.empresa (
  id                       uuid        NOT NULL DEFAULT app.uuid_v7(),
  razao_social             text        NOT NULL,
  nome_fantasia            text,
  cnpj                     char(14),
  plano                    text        DEFAULT 'essencial',    -- nulo: conta só do plano Serviços
  srv_habilitado           boolean     NOT NULL DEFAULT false, -- módulo SRV, plano Serviços (sumário 2.2)
  mfa_obrigatorio          boolean     NOT NULL DEFAULT false, -- MFA exigido de todos os usuários, ligado pelo administrador (RF-ADM-006; sumário 3.6)
  situacao                 text        NOT NULL DEFAULT 'teste',
  fuso_horario             text        NOT NULL DEFAULT 'America/Sao_Paulo',
  retencao_auditoria_anos  smallint    NOT NULL DEFAULT 5,
  celula                   text        NOT NULL DEFAULT 'pool-01', -- cluster onde os dados residem (seção 15)
  teste_ate                date,
  versao bigint NOT NULL DEFAULT 1, alterado_xid xid8 NOT NULL DEFAULT pg_current_xact_id(),
  criado_em timestamptz NOT NULL DEFAULT now(), criado_por uuid,
  atualizado_em timestamptz NOT NULL DEFAULT now(), atualizado_por uuid, excluido_em timestamptz,
  CONSTRAINT pk_empresa PRIMARY KEY (id),
  CONSTRAINT uq_empresa_cnpj UNIQUE (cnpj),
  CONSTRAINT ck_empresa_plano CHECK (plano IN ('essencial', 'profissional', 'enterprise')),
  -- a subcontratada pode ter só o plano Serviços ou somá-lo a um dos planos acima
  CONSTRAINT ck_empresa_contratacao CHECK (plano IS NOT NULL OR srv_habilitado),
  CONSTRAINT ck_empresa_situacao CHECK (situacao IN ('teste', 'ativa', 'suspensa', 'cancelada')),
  CONSTRAINT ck_empresa_retencao CHECK (retencao_auditoria_anos BETWEEN 5 AND 30),
  CONSTRAINT ck_empresa_cnpj CHECK (cnpj ~ '^[0-9]{14}$')
);
-- +goose StatementEnd

-- +goose StatementBegin
-- Colunas padrão da empresa. app.tg_padrao() não serve: ela compara
-- NEW.empresa_id, que esta tabela não tem, e o PL/pgSQL falharia na execução.
-- A versão própria torna o id imutável, pelo mesmo motivo que tg_padrao torna
-- empresa_id imutável: uma empresa nunca vira outra.
CREATE OR REPLACE FUNCTION app.tg_padrao_empresa() RETURNS trigger
  LANGUAGE plpgsql
AS $$
BEGIN
  IF TG_OP = 'INSERT' THEN
    NEW.versao         := 1;
    NEW.criado_em      := now();
    NEW.criado_por     := coalesce(NEW.criado_por, app.usuario_atual());
    NEW.atualizado_em  := NEW.criado_em;
    NEW.atualizado_por := NEW.criado_por;
  ELSE
    IF NEW.id IS DISTINCT FROM OLD.id THEN
      RAISE EXCEPTION 'o id da empresa é imutável' USING ERRCODE = '42501';
    END IF;
    NEW.versao         := OLD.versao + 1;
    NEW.criado_em      := OLD.criado_em;
    NEW.criado_por     := OLD.criado_por;
    NEW.atualizado_em  := now();
    NEW.atualizado_por := app.usuario_atual();
  END IF;
  NEW.alterado_xid := pg_current_xact_id();
  RETURN NEW;
END
$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_padrao BEFORE INSERT OR UPDATE ON app.empresa
  FOR EACH ROW EXECUTE FUNCTION app.tg_padrao_empresa();
CREATE TRIGGER trg_auditoria AFTER INSERT OR UPDATE OR DELETE ON app.empresa
  FOR EACH ROW EXECUTE FUNCTION app.tg_auditar('id');

ALTER TABLE app.empresa ENABLE ROW LEVEL SECURITY;
ALTER TABLE app.empresa FORCE  ROW LEVEL SECURITY;

-- A própria empresa, por id (seção 2.3). A criação da conta abre a transação
-- já com o id novo em app.empresa_id, gerado no Go.
CREATE POLICY empresa_propria ON app.empresa
  FOR ALL TO tenaz_app, tenaz_leitura
  USING (id = app.empresa_atual()) WITH CHECK (id = app.empresa_atual());

-- Jobs listam as empresas para percorrê-las uma a uma; não criam nem apagam.
-- A listagem só vale FORA da transação de uma empresa: dentro dela, o job vê
-- só a empresa da transação, como qualquer papel (ISO-E01-04). USING (true)
-- deixaria o job, no contexto de A, ler B.
CREATE POLICY empresa_listagem_job ON app.empresa
  FOR SELECT TO tenaz_job USING (app.empresa_atual() IS NULL OR id = app.empresa_atual());
CREATE POLICY empresa_propria_job ON app.empresa
  FOR UPDATE TO tenaz_job
  USING (id = app.empresa_atual()) WITH CHECK (id = app.empresa_atual());

-- Escopo externo: prestador, locatário e link leem o cadastro da empresa
-- contratante (nome, fuso), mas não o alteram nem o apagam. Sem esta policy,
-- empresa_propria (FOR ALL) deixaria um prestador editar a razão social.
--
-- Fica fora tenaz_job, de propósito: o job lista as empresas ANTES de abrir a
-- transação de cada uma, portanto sem escopo, e a restritiva o deixaria cego.
-- O job nunca roda com escopo externo, e a escrita dele já é só UPDATE da
-- própria empresa (empresa_propria_job).
CREATE POLICY escopo_externo ON app.empresa AS RESTRICTIVE
  FOR ALL TO tenaz_app, tenaz_leitura
  USING (app.escopo_atual() IN ('interno', 'prestador', 'locatario', 'link'))
  WITH CHECK (app.escopo_interno());
CREATE POLICY escopo_externo_exclusao ON app.empresa AS RESTRICTIVE
  FOR DELETE TO tenaz_app, tenaz_leitura
  USING (app.escopo_interno());

-- Os privilégios padrão dariam INSERT e DELETE ao job, que não cria nem apaga
-- empresa; o restante fica como o padrão.
REVOKE INSERT, DELETE ON app.empresa FROM tenaz_job;
-- +goose StatementEnd

-- ------------------------------------------------------------------ unidade

-- +goose StatementBegin
-- Unidade: planta, filial, base ou loja.
CREATE TABLE app.unidade (
  empresa_id     uuid    NOT NULL,
  id             uuid    NOT NULL DEFAULT app.uuid_v7(),
  codigo         text    NOT NULL,
  nome           text    NOT NULL,
  cnpj           char(14),
  endereco       jsonb,
  fuso_horario   text,                       -- nulo herda o da empresa
  latitude       numeric(9,6),
  longitude      numeric(9,6),
  situacao       text    NOT NULL DEFAULT 'ativa',
  versao bigint NOT NULL DEFAULT 1, alterado_xid xid8 NOT NULL DEFAULT pg_current_xact_id(),
  criado_em timestamptz NOT NULL DEFAULT now(), criado_por uuid,
  atualizado_em timestamptz NOT NULL DEFAULT now(), atualizado_por uuid, excluido_em timestamptz,
  CONSTRAINT pk_unidade PRIMARY KEY (empresa_id, id),
  CONSTRAINT fk_unidade_empresa FOREIGN KEY (empresa_id) REFERENCES app.empresa (id),
  CONSTRAINT ck_unidade_situacao CHECK (situacao IN ('ativa', 'inativa')),
  CONSTRAINT ck_unidade_cnpj CHECK (cnpj ~ '^[0-9]{14}$')
);
-- RF-ADM-002: código único na empresa.
CREATE UNIQUE INDEX uq_unidade_codigo ON app.unidade (empresa_id, codigo) WHERE excluido_em IS NULL;
-- Descida do app (seção 11.3).
CREATE INDEX ix_unidade_sync ON app.unidade (empresa_id, alterado_xid);
CALL app.preparar_tabela('app.unidade', p_escopo => 'proprio');

-- Leitura para os quatro escopos (seção 2.9.4); gravação e exclusão só internas.
CREATE POLICY escopo_externo ON app.unidade AS RESTRICTIVE
  FOR ALL TO tenaz_app, tenaz_job, tenaz_leitura
  USING (app.escopo_atual() IN ('interno', 'prestador', 'locatario', 'link'))
  WITH CHECK (app.escopo_interno());
-- O WITH CHECK não alcança o DELETE: sem esta policy, um escopo que só lê
-- poderia apagar a linha que lê.
CREATE POLICY escopo_externo_exclusao ON app.unidade AS RESTRICTIVE
  FOR DELETE TO tenaz_app, tenaz_job, tenaz_leitura
  USING (app.escopo_interno());
-- +goose StatementEnd

-- +goose StatementBegin
-- Calendário de trabalho da unidade (RF-ADM-002): turnos por dia da semana.
-- "Pelo menos um turno" é validação do serviço, não do banco (seção 5.1).
CREATE TABLE app.turno_unidade (
  empresa_id   uuid     NOT NULL,
  id           uuid     NOT NULL DEFAULT app.uuid_v7(),
  unidade_id   uuid     NOT NULL,
  dia_semana   smallint NOT NULL,          -- ISO 8601, como extract(isodow): 1 = segunda ... 7 = domingo
  hora_inicio  time     NOT NULL,
  hora_fim     time     NOT NULL,          -- menor que o início: termina no dia seguinte; 24:00 = meia-noite
  nome         text,                       -- ex.: "1º turno"
  versao bigint NOT NULL DEFAULT 1, alterado_xid xid8 NOT NULL DEFAULT pg_current_xact_id(),
  criado_em timestamptz NOT NULL DEFAULT now(), criado_por uuid,
  atualizado_em timestamptz NOT NULL DEFAULT now(), atualizado_por uuid, excluido_em timestamptz,
  CONSTRAINT pk_turno_unidade PRIMARY KEY (empresa_id, id),
  CONSTRAINT fk_turno_unidade_unidade FOREIGN KEY (empresa_id, unidade_id) REFERENCES app.unidade (empresa_id, id),
  CONSTRAINT ck_turno_unidade_dia CHECK (dia_semana BETWEEN 1 AND 7),
  CONSTRAINT ck_turno_unidade_horario CHECK (hora_inicio < time '24:00' AND hora_fim <> hora_inicio)
);
-- Leitura do calendário por unidade (tela e tempo requerido dos KPIs) e
-- trava contra turno repetido no mesmo dia.
CREATE UNIQUE INDEX uq_turno_unidade_inicio ON app.turno_unidade (empresa_id, unidade_id, dia_semana, hora_inicio)
  WHERE excluido_em IS NULL;
CALL app.preparar_tabela('app.turno_unidade');
-- +goose StatementEnd

-- +goose StatementBegin
-- Feriados da unidade (RF-ADM-002). Feriado nacional é gravado em cada
-- unidade pelo serviço, para que a unidade possa trabalhar nele.
CREATE TABLE app.feriado_unidade (
  empresa_id    uuid  NOT NULL,
  id            uuid  NOT NULL DEFAULT app.uuid_v7(),
  unidade_id    uuid  NOT NULL,
  data_feriado  date  NOT NULL,
  descricao     text  NOT NULL,
  versao bigint NOT NULL DEFAULT 1, alterado_xid xid8 NOT NULL DEFAULT pg_current_xact_id(),
  criado_em timestamptz NOT NULL DEFAULT now(), criado_por uuid,
  atualizado_em timestamptz NOT NULL DEFAULT now(), atualizado_por uuid, excluido_em timestamptz,
  CONSTRAINT pk_feriado_unidade PRIMARY KEY (empresa_id, id),
  CONSTRAINT fk_feriado_unidade_unidade FOREIGN KEY (empresa_id, unidade_id) REFERENCES app.unidade (empresa_id, id),
  CONSTRAINT ck_feriado_unidade_descricao CHECK (btrim(descricao) <> '')
);
-- Um feriado por dia na unidade; serve também à leitura por período.
CREATE UNIQUE INDEX uq_feriado_unidade_data ON app.feriado_unidade (empresa_id, unidade_id, data_feriado)
  WHERE excluido_em IS NULL;
CALL app.preparar_tabela('app.feriado_unidade');
-- +goose StatementEnd

-- --------------------------------------------------------------- identidade

-- +goose StatementBegin
-- Identidade de login, global. Uma pessoa com o mesmo e-mail em duas empresas
-- tem uma identidade e dois usuários. Sem dado de negócio e sem RLS: é lida
-- antes de a empresa ser conhecida. Lista de tabelas globais (seção 2.11).
-- Dono tenaz_credencial (ISO-E01-08): nem tenaz_dono nem tenaz_manutencao a leem.
SET ROLE tenaz_credencial;
CREATE TABLE plataforma.identidade (
  id                    uuid        NOT NULL DEFAULT app.uuid_v7(),
  email                 citext      NOT NULL,
  senha_hash            text,                 -- argon2id; nulo quando só usa SSO
  mfa_segredo_cifrado   bytea,                -- cifrado com chave do KMS
  mfa_ativado_em        timestamptz,          -- TOTP confirmado; nulo: sem segundo fator
  mfa_ultimo_passo      bigint,               -- último passo TOTP aceito (RFC 6238, 5.2): o mesmo código
                                              -- não vale duas vezes dentro da janela
  email_verificado_em   timestamptz,
  falhas_login          smallint    NOT NULL DEFAULT 0,  -- falhas seguidas; zera no login certo
  bloqueada_ate         timestamptz,          -- bloqueio temporário após 5 falhas (RF-ADM-006)
  bloqueada_em          timestamptz,          -- bloqueio administrativo, até desbloqueio explícito
  ultimo_login_em       timestamptz,
  criado_em             timestamptz NOT NULL DEFAULT now(),
  atualizado_em         timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT pk_identidade PRIMARY KEY (id),
  CONSTRAINT uq_identidade_email UNIQUE (email),
  CONSTRAINT ck_identidade_falhas CHECK (falhas_login >= 0),
  -- MFA ativado exige o segredo guardado
  CONSTRAINT ck_identidade_mfa CHECK (mfa_ativado_em IS NULL OR mfa_segredo_cifrado IS NOT NULL),
  CONSTRAINT ck_identidade_mfa_passo CHECK (mfa_ultimo_passo >= 0)
);
-- REFERENCES para a FK de app.usuario.identidade_id: o dono de app.usuario
-- pode criar a FK, sem ler a tabela.
REVOKE ALL ON plataforma.identidade FROM PUBLIC;
GRANT REFERENCES (id) ON plataforma.identidade TO tenaz_dono;
SET ROLE tenaz_dono;
-- +goose StatementEnd

-- ------------------------------------------------------------------- pessoa

-- +goose StatementBegin
-- Pessoa: colaborador ou terceiro que executa, conduz ou responde por algo
-- (RF-ADM-015). Pode ou não ter usuário.
CREATE TABLE app.pessoa (
  empresa_id            uuid    NOT NULL,
  id                    uuid    NOT NULL DEFAULT app.uuid_v7(),
  nome                  text    NOT NULL,
  cpf                   char(11),               -- dado pessoal (LGPD); visível só a perfis autorizados
  data_nascimento       date,                   -- dado pessoal; usado em CNH e ASO
  matricula             text,
  cargo                 text,
  telefone              text,
  email                 citext,
  vinculo               text    NOT NULL DEFAULT 'empregado',
  fornecedor_id         uuid,                   -- técnico de fornecedor (vinculo = 'terceiro'); FK em E06
  unidade_id            uuid,
  registro_profissional text,                   -- CREA, CFT etc. (PLH, responsável NR-12)
  situacao              text    NOT NULL DEFAULT 'ativa',
  versao bigint NOT NULL DEFAULT 1, alterado_xid xid8 NOT NULL DEFAULT pg_current_xact_id(),
  criado_em timestamptz NOT NULL DEFAULT now(), criado_por uuid,
  atualizado_em timestamptz NOT NULL DEFAULT now(), atualizado_por uuid, excluido_em timestamptz,
  CONSTRAINT pk_pessoa PRIMARY KEY (empresa_id, id),
  CONSTRAINT fk_pessoa_empresa FOREIGN KEY (empresa_id) REFERENCES app.empresa (id),
  CONSTRAINT fk_pessoa_unidade FOREIGN KEY (empresa_id, unidade_id) REFERENCES app.unidade (empresa_id, id),
  CONSTRAINT ck_pessoa_cpf CHECK (cpf ~ '^[0-9]{11}$'),
  CONSTRAINT ck_pessoa_vinculo CHECK (vinculo IN ('empregado', 'terceiro', 'autonomo')),
  CONSTRAINT ck_pessoa_fornecedor CHECK (fornecedor_id IS NULL OR vinculo = 'terceiro'),
  -- sumário 3.6: só ativa ou inativa; pessoa inativa nunca é excluída e
  -- continua ligada ao histórico dos serviços
  CONSTRAINT ck_pessoa_situacao CHECK (situacao IN ('ativa', 'inativa'))
);
-- RF-ADM-015: CPF único na empresa, quando informado.
CREATE UNIQUE INDEX uq_pessoa_cpf ON app.pessoa (empresa_id, cpf) WHERE cpf IS NOT NULL AND excluido_em IS NULL;
-- Busca por trecho do nome na tela de pessoas.
CREATE INDEX ix_pessoa_nome ON app.pessoa USING gin (empresa_id, nome gin_trgm_ops);
-- Escopo do prestador: pessoas do seu fornecedor.
CREATE INDEX ix_pessoa_fornecedor ON app.pessoa (empresa_id, fornecedor_id) WHERE fornecedor_id IS NOT NULL;
CREATE INDEX ix_pessoa_sync ON app.pessoa (empresa_id, alterado_xid);
CALL app.preparar_tabela('app.pessoa', p_escopo => 'proprio');
-- +goose StatementEnd

-- ------------------------------------------------------------------ usuario

-- +goose StatementBegin
-- Usuário: login de acesso (RF-ADM-004). Pode estar ligado a uma pessoa;
-- usuário de prestador fica ligado ao fornecedor e o de locatário, ao locatário.
CREATE TABLE app.usuario (
  empresa_id         uuid    NOT NULL,
  id                 uuid    NOT NULL DEFAULT app.uuid_v7(),
  identidade_id      uuid,                   -- nulo até o convite ser aceito; depois, imutável
  pessoa_id          uuid,                   -- opcional
  fornecedor_id      uuid,                   -- perfil prestador; FK em E06
  locatario_id       uuid,                   -- perfil locatario (fase 2); FK com a tabela locatario
  nome               text    NOT NULL,
  email              citext,
  telefone           text,
  custo_hora         numeric(12,2),          -- custo de mão de obra apontada
  unidade_padrao_id  uuid,
  situacao           text    NOT NULL DEFAULT 'ativo',
  versao bigint NOT NULL DEFAULT 1, alterado_xid xid8 NOT NULL DEFAULT pg_current_xact_id(),
  criado_em timestamptz NOT NULL DEFAULT now(), criado_por uuid,
  atualizado_em timestamptz NOT NULL DEFAULT now(), atualizado_por uuid, excluido_em timestamptz,
  CONSTRAINT pk_usuario PRIMARY KEY (empresa_id, id),
  -- alvo das FKs da plataforma: a sessão, o código OAuth e o acesso à empresa
  -- só existem para a identidade que é de fato a do usuário naquela empresa
  CONSTRAINT uq_usuario_id_identidade UNIQUE (empresa_id, id, identidade_id),
  CONSTRAINT fk_usuario_empresa FOREIGN KEY (empresa_id) REFERENCES app.empresa (id),
  CONSTRAINT fk_usuario_identidade FOREIGN KEY (identidade_id) REFERENCES plataforma.identidade (id),
  CONSTRAINT fk_usuario_pessoa FOREIGN KEY (empresa_id, pessoa_id) REFERENCES app.pessoa (empresa_id, id),
  CONSTRAINT fk_usuario_unidade FOREIGN KEY (empresa_id, unidade_padrao_id) REFERENCES app.unidade (empresa_id, id),
  -- sumário 3.6: o bloqueio temporário por falhas de login é da identidade
  -- (identidade.bloqueada_ate), não um estado do usuário
  CONSTRAINT ck_usuario_situacao CHECK (situacao IN ('convidado', 'ativo', 'inativo')),
  CONSTRAINT ck_usuario_externo CHECK (num_nonnulls(fornecedor_id, locatario_id) <= 1)
);
-- Uma identidade tem no máximo um usuário por empresa.
CREATE UNIQUE INDEX uq_usuario_identidade ON app.usuario (empresa_id, identidade_id) WHERE identidade_id IS NOT NULL;
-- RF-ADM-004: e-mail único na empresa.
CREATE UNIQUE INDEX uq_usuario_email ON app.usuario (empresa_id, email) WHERE email IS NOT NULL AND excluido_em IS NULL;
-- RF-ADM-015: uma pessoa tem no máximo um usuário por empresa.
CREATE UNIQUE INDEX uq_usuario_pessoa ON app.usuario (empresa_id, pessoa_id) WHERE pessoa_id IS NOT NULL AND excluido_em IS NULL;
-- Busca por trecho do nome na tela de usuários.
CREATE INDEX ix_usuario_nome ON app.usuario USING gin (empresa_id, nome gin_trgm_ops);
-- Escopo externo: colegas do mesmo fornecedor ou locatário; desativação em lote
-- quando o fornecedor é inativado (documento 9, seção 5.3).
CREATE INDEX ix_usuario_fornecedor ON app.usuario (empresa_id, fornecedor_id) WHERE fornecedor_id IS NOT NULL;
CREATE INDEX ix_usuario_locatario ON app.usuario (empresa_id, locatario_id) WHERE locatario_id IS NOT NULL;
CREATE INDEX ix_usuario_sync ON app.usuario (empresa_id, alterado_xid);
CALL app.preparar_tabela('app.usuario', p_escopo => 'proprio');
-- As FKs compostas das tabelas de credencial apontam para estas colunas.
GRANT REFERENCES (empresa_id, id, identidade_id) ON app.usuario TO tenaz_credencial;
-- +goose StatementEnd

-- +goose StatementBegin
-- Escopo externo de usuario (seção 2.9.4): o usuário externo vê a si mesmo e
-- os colegas do mesmo fornecedor ou locatário, e só altera a própria linha,
-- sem trocar o vínculo — senão um prestador apagaria o próprio fornecedor_id e
-- passaria a ser interno. O link não vê usuário nenhum.
CREATE POLICY escopo_externo ON app.usuario AS RESTRICTIVE
  FOR ALL TO tenaz_app, tenaz_job, tenaz_leitura
  USING (
       app.escopo_interno()
    OR (app.escopo_atual() = 'prestador'
        AND (id = app.usuario_atual() OR fornecedor_id = app.fornecedor_atual()))
    OR (app.escopo_atual() = 'locatario'
        AND (id = app.usuario_atual() OR locatario_id = app.locatario_atual())))
  WITH CHECK (
       app.escopo_interno()
    OR (app.escopo_atual() = 'prestador'
        AND id = app.usuario_atual() AND fornecedor_id = app.fornecedor_atual())
    OR (app.escopo_atual() = 'locatario'
        AND id = app.usuario_atual() AND locatario_id = app.locatario_atual()));
CREATE POLICY escopo_externo_exclusao ON app.usuario AS RESTRICTIVE
  FOR DELETE TO tenaz_app, tenaz_job, tenaz_leitura
  USING (app.escopo_interno());
-- Cadastro de usuário é só interno. O WITH CHECK acima já só aceitaria a
-- própria linha, que existe (PK); esta policy fecha o INSERT por desenho, e
-- não por acaso de chave.
CREATE POLICY escopo_externo_insercao ON app.usuario AS RESTRICTIVE
  FOR INSERT TO tenaz_app, tenaz_job, tenaz_leitura
  WITH CHECK (app.escopo_interno());
-- +goose StatementEnd

-- +goose StatementBegin
-- O que o usuário externo altera na própria linha (ISO-E01-01). O RLS corta
-- linhas, não colunas: sem este trigger, o prestador ou o locatário mudaria
-- pessoa_id (e passaria a ler o CPF da pessoa interna pela policy de pessoa),
-- custo_hora, unidade_padrao_id, e-mail ou situacao.
--
-- Lista de permitidas, e não de proibidas: fora do escopo interno, só nome e
-- telefone mudam (PATCH /v1/usuarios/eu). Coluna nova nasce proibida para
-- externo. As colunas padrão mantidas por trg_padrao (que roda antes, pela
-- ordem alfabética dos triggers) ficam fora da comparação.
--
-- Escolhido em vez de "UPDATE só interno" porque o prestador precisa alterar o
-- próprio nome e telefone, e porque a regra vale também para o link e para
-- qualquer escopo externo que vier.
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
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_usuario_colunas_externo BEFORE UPDATE ON app.usuario
  FOR EACH ROW EXECUTE FUNCTION app.tg_usuario_colunas_externo();
-- +goose StatementEnd

-- +goose StatementBegin
-- Escopo externo de pessoa (seção 2.9.4): o prestador vê as pessoas do seu
-- fornecedor (inclusive técnicos sem login); prestador e locatário veem as
-- pessoas ligadas aos usuários que já enxergam — ele mesmo e os colegas. A
-- subconsulta passa pela policy de usuario, que não consulta pessoa: não há
-- ciclo. Nenhum escopo externo grava pessoa.
CREATE POLICY escopo_externo ON app.pessoa AS RESTRICTIVE
  FOR ALL TO tenaz_app, tenaz_job, tenaz_leitura
  USING (
       app.escopo_interno()
    OR (app.escopo_atual() = 'prestador' AND fornecedor_id = app.fornecedor_atual())
    OR (app.escopo_atual() IN ('prestador', 'locatario')
        AND id IN (SELECT u.pessoa_id FROM app.usuario u WHERE u.pessoa_id IS NOT NULL)))
  WITH CHECK (app.escopo_interno());
CREATE POLICY escopo_externo_exclusao ON app.pessoa AS RESTRICTIVE
  FOR DELETE TO tenaz_app, tenaz_job, tenaz_leitura
  USING (app.escopo_interno());
-- +goose StatementEnd

-- ----------------------------------------------------------- acesso_empresa

-- +goose StatementBegin
-- Em quais empresas a identidade entra. Mantida SÓ pelo trigger em
-- app.usuario: a aplicação lê, nunca grava. Expõe o vínculo, nunca dado de
-- negócio. Lista de tabelas globais (seção 2.11).
--
-- A FK composta para (empresa_id, usuario_id, identidade_id) garante que o
-- acesso aponta para um usuário daquela empresa cuja identidade é a mesma.
-- Dono tenaz_credencial (ISO-E01-08).
SET ROLE tenaz_credencial;
CREATE TABLE plataforma.acesso_empresa (
  identidade_id uuid        NOT NULL,
  empresa_id    uuid        NOT NULL,
  usuario_id    uuid        NOT NULL,
  habilitado    boolean     NOT NULL DEFAULT true,   -- usuario ativo e não excluído
  atualizado_em timestamptz NOT NULL DEFAULT now(),
  CONSTRAINT pk_acesso_empresa PRIMARY KEY (identidade_id, empresa_id),
  CONSTRAINT fk_acesso_empresa_identidade FOREIGN KEY (identidade_id)
    REFERENCES plataforma.identidade (id) ON DELETE CASCADE,
  CONSTRAINT fk_acesso_empresa_usuario FOREIGN KEY (empresa_id, usuario_id, identidade_id)
    REFERENCES app.usuario (empresa_id, id, identidade_id) ON DELETE CASCADE
);
REVOKE ALL ON plataforma.acesso_empresa FROM PUBLIC;
-- +goose StatementEnd

-- +goose StatementBegin
-- Manutenção de acesso_empresa. SECURITY DEFINER para que tenaz_app não
-- precise de privilégio de escrita na tabela: se tivesse, um defeito de código
-- poderia gravar um acesso a qualquer empresa. Com o trigger, o acesso só nasce
-- de uma linha real de app.usuario, que já passou pelo WITH CHECK da policy de
-- empresa da transação.
--
-- search_path fixo, pg_temp por último e todo nome qualificado. Está na lista
-- aprovada de funções SECURITY DEFINER (seção 2.8, consulta 6). Dono
-- tenaz_credencial (ISO-E01-08): roda com o privilégio do dono das credenciais,
-- e não do tenaz_dono, que não as alcança.
CREATE OR REPLACE FUNCTION plataforma.tg_acesso_empresa() RETURNS trigger
  LANGUAGE plpgsql SECURITY DEFINER
  SET search_path = pg_catalog, pg_temp
AS $$
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
$$;
-- +goose StatementEnd

-- +goose StatementBegin
REVOKE ALL ON FUNCTION plataforma.tg_acesso_empresa() FROM PUBLIC;
-- Criar o trigger exige EXECUTE na função para o dono de app.usuario; função
-- de trigger não se chama diretamente.
GRANT EXECUTE ON FUNCTION plataforma.tg_acesso_empresa() TO tenaz_dono;
SET ROLE tenaz_dono;
CREATE TRIGGER trg_acesso_empresa AFTER INSERT OR UPDATE OF identidade_id, situacao, excluido_em ON app.usuario
  FOR EACH ROW WHEN (NEW.identidade_id IS NOT NULL)
  EXECUTE FUNCTION plataforma.tg_acesso_empresa();
-- +goose StatementEnd

-- +goose StatementBegin
-- Privilégios na plataforma (seção 5.1). O esquema não tem privilégio padrão.
-- USAGE para o tenaz_app (funções seguras da conexão e catálogos globais) e
-- para o tenaz_job (purga). Nenhum dos dois tem privilégio em identidade nem
-- em acesso_empresa: as credenciais são só do papel tenaz_autenticacao, que
-- nasce na 00005 (ADR-011, documento 5; achado ISO-E01-02). A escrita em
-- acesso_empresa continua só do trigger acima.
GRANT USAGE ON SCHEMA plataforma TO tenaz_app, tenaz_job;
-- +goose StatementEnd

-- ------------------------------------------------------------------- perfil

-- +goose StatementBegin
-- Permissões que usuário externo (prestador, locatário) nunca recebe, nem por
-- perfil nem avulsa (RF-ADM-005, regras dos perfis padrão): os.encerrar; as
-- nove da configuração, de PRJ e de CTR (sumário 3.6); srv.gerenciar_conexao;
-- pessoa.saude.ler; e as do grupo Admin do catálogo.
--
-- O catálogo completo vive no Go. Aqui fica só esta lista curta, porque é a
-- parte que nenhum defeito de código pode violar. IMMUTABLE para poder entrar
-- em CHECK; mudar a lista é CREATE OR REPLACE seguido de recriar os CHECKs que
-- a usam com NOT VALID + VALIDATE, porque um CHECK não reconfere sozinho as
-- linhas antigas.
CREATE OR REPLACE FUNCTION app.permissoes_vedadas_externo() RETURNS text[]
  LANGUAGE sql IMMUTABLE PARALLEL SAFE
AS $$
  SELECT ARRAY[
    'os.encerrar',
    'adm.gerenciar_regras', 'adm.gerenciar_configuracao', 'atv.colocar_em_operacao',
    'prj.gerenciar', 'prj.dispensar_etapa', 'prj.ajustar_medicao',
    'ctr.gerenciar', 'ctr.fechar_medicao', 'ctr.reabrir_medicao',
    'srv.gerenciar_conexao',
    'pessoa.saude.ler',
    -- grupo Admin
    'adm.gerenciar_empresa', 'adm.gerenciar_plano', 'adm.gerenciar_usuarios', 'adm.gerenciar_perfis',
    'adm.gerenciar_parametros', 'adm.ler_auditoria', 'adm.gerenciar_sso', 'adm.exportar',
    'con.gerenciar_tipos_documento', 'int.gerenciar_tokens', 'int.gerenciar_webhooks',
    'int.gerenciar_conexoes_telematica', 'int.gerenciar_erp', 'ia.gerenciar'
  ]::text[]
$$;
-- +goose StatementEnd

-- +goose StatementBegin
-- Perfil: conjunto de permissões (RF-ADM-005). Os 12 perfis padrão são criados
-- pelo serviço junto com a empresa.
CREATE TABLE app.perfil (
  empresa_id   uuid    NOT NULL,
  id           uuid    NOT NULL DEFAULT app.uuid_v7(),
  codigo       text    NOT NULL,
  nome         text    NOT NULL,
  padrao       boolean NOT NULL DEFAULT false,
  -- perfil padrão de que este deriva (no padrão, o próprio código). Nulo:
  -- personalizado sem origem, sempre interno. É o que diz se o perfil é
  -- externo, inclusive o personalizado derivado de prestador ou locatario
  -- (Fase 2), que o código sozinho não revela.
  perfil_origem text,
  permissoes   text[]  NOT NULL DEFAULT '{}',   -- <módulo>.<ação>, ex.: {'os.executar','os.encerrar','est.decidir_destino'}
  versao bigint NOT NULL DEFAULT 1, alterado_xid xid8 NOT NULL DEFAULT pg_current_xact_id(),
  criado_em timestamptz NOT NULL DEFAULT now(), criado_por uuid,
  atualizado_em timestamptz NOT NULL DEFAULT now(), atualizado_por uuid, excluido_em timestamptz,
  CONSTRAINT pk_perfil PRIMARY KEY (empresa_id, id),
  CONSTRAINT uq_perfil_codigo UNIQUE (empresa_id, codigo),
  CONSTRAINT fk_perfil_empresa FOREIGN KEY (empresa_id) REFERENCES app.empresa (id),
  CONSTRAINT ck_perfil_padrao CHECK (NOT padrao OR codigo IN
    ('administrador', 'gestor', 'planejador', 'tecnico', 'almoxarife', 'comprador', 'plh', 'motorista',
     'solicitante', 'visualizador', 'prestador', 'locatario')),
  -- dados de saúde ficam fora de todo perfil padrão (sumário 3.1)
  CONSTRAINT ck_perfil_saude CHECK (NOT padrao OR NOT ('pessoa.saude.ler' = ANY (permissoes))),
  CONSTRAINT ck_perfil_origem CHECK (perfil_origem IN
    ('administrador', 'gestor', 'planejador', 'tecnico', 'almoxarife', 'comprador', 'plh', 'motorista',
     'solicitante', 'visualizador', 'prestador', 'locatario')),
  CONSTRAINT ck_perfil_padrao_origem CHECK (NOT padrao OR perfil_origem = codigo),
  -- a liberação de OS nunca é dada a usuário externo (sumário 3.6), nem por
  -- perfil personalizado derivado de prestador ou locatario
  CONSTRAINT ck_perfil_externo_encerrar CHECK (coalesce(perfil_origem, codigo) NOT IN ('prestador', 'locatario')
                                               OR NOT ('os.encerrar' = ANY (permissoes))),
  -- nem nenhuma das demais vedadas a externo
  CONSTRAINT ck_perfil_externo_vedadas CHECK (coalesce(perfil_origem, codigo) NOT IN ('prestador', 'locatario')
                                              OR NOT (permissoes && app.permissoes_vedadas_externo()))
);
CALL app.preparar_tabela('app.perfil');
-- +goose StatementEnd

-- +goose StatementBegin
-- O tipo do perfil não muda depois de criado (ISO-E01-03). O tipo é
-- coalesce(perfil_origem, codigo) e diz se o perfil é externo; se pudesse
-- mudar, um perfil já atribuído a um prestador viraria interno e receberia
-- os.encerrar, e os CHECKs e o trigger de atribuição, que olham o tipo no
-- momento da gravação, não veriam a atribuição antiga. perfil_origem e padrao
-- ficam imutáveis também, por serem parte do tipo. Mudar de tipo é criar outro
-- perfil.
CREATE OR REPLACE FUNCTION app.tg_perfil_tipo() RETURNS trigger
  LANGUAGE plpgsql
AS $$
BEGIN
  IF NEW.perfil_origem IS DISTINCT FROM OLD.perfil_origem
     OR NEW.padrao IS DISTINCT FROM OLD.padrao
     OR coalesce(NEW.perfil_origem, NEW.codigo) IS DISTINCT FROM coalesce(OLD.perfil_origem, OLD.codigo) THEN
    RAISE EXCEPTION 'o tipo do perfil (origem, padrão) é imutável depois de criado' USING ERRCODE = '42501';
  END IF;
  RETURN NEW;
END
$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_perfil_tipo BEFORE UPDATE OF codigo, perfil_origem, padrao ON app.perfil
  FOR EACH ROW EXECUTE FUNCTION app.tg_perfil_tipo();
-- +goose StatementEnd

-- +goose StatementBegin
-- Usuário pode ter vários perfis; unidade_id nulo = todas as unidades.
CREATE TABLE app.usuario_perfil (
  empresa_id  uuid NOT NULL,
  id          uuid NOT NULL DEFAULT app.uuid_v7(),
  usuario_id  uuid NOT NULL,
  perfil_id   uuid NOT NULL,
  unidade_id  uuid,
  versao bigint NOT NULL DEFAULT 1, alterado_xid xid8 NOT NULL DEFAULT pg_current_xact_id(),
  criado_em timestamptz NOT NULL DEFAULT now(), criado_por uuid,
  atualizado_em timestamptz NOT NULL DEFAULT now(), atualizado_por uuid, excluido_em timestamptz,
  CONSTRAINT pk_usuario_perfil PRIMARY KEY (empresa_id, id),
  -- também é o índice de "perfis do usuário", lido a cada requisição
  CONSTRAINT fk_usuario_perfil_usuario FOREIGN KEY (empresa_id, usuario_id) REFERENCES app.usuario (empresa_id, id),
  CONSTRAINT fk_usuario_perfil_perfil FOREIGN KEY (empresa_id, perfil_id) REFERENCES app.perfil (empresa_id, id),
  CONSTRAINT fk_usuario_perfil_unidade FOREIGN KEY (empresa_id, unidade_id) REFERENCES app.unidade (empresa_id, id)
);
-- Uma atribuição ativa por usuário, perfil e unidade (nula = todas). Parcial,
-- como as demais unicidades: retirar o perfil (excluido_em) e atribuí-lo de
-- novo não colide com a linha retirada. É também o índice de "perfis do
-- usuário", lido a cada requisição.
CREATE UNIQUE INDEX uq_usuario_perfil ON app.usuario_perfil (empresa_id, usuario_id, perfil_id, unidade_id)
  NULLS NOT DISTINCT WHERE excluido_em IS NULL;
-- Quem tem um perfil (contagem de técnicos cobrados, último administrador).
CREATE INDEX ix_usuario_perfil_perfil ON app.usuario_perfil (empresa_id, perfil_id) WHERE excluido_em IS NULL;
CALL app.preparar_tabela('app.usuario_perfil');
-- +goose StatementEnd

-- +goose StatementBegin
-- Usuário externo só tem o perfil do seu tipo; usuário interno nunca recebe
-- perfil externo (seção 2.9.4, RF-ADM-004).
--
-- Perfil ou usuário não encontrado: deixa passar, e a FK composta recusa com
-- 23503. A leitura é feita sob o RLS de quem grava, então um id de outra
-- empresa simplesmente não aparece aqui.
CREATE OR REPLACE FUNCTION app.tg_usuario_perfil_externo() RETURNS trigger
  LANGUAGE plpgsql
AS $$
DECLARE
  v_codigo     text;
  v_fornecedor uuid;
  v_locatario  uuid;
BEGIN
  IF NEW.excluido_em IS NOT NULL THEN
    RETURN NEW;                       -- retirar um perfil nunca cria incompatibilidade
  END IF;
  -- o tipo do perfil é o da origem: personalizado derivado de prestador é de prestador
  SELECT coalesce(p.perfil_origem, p.codigo) INTO v_codigo FROM app.perfil p
  WHERE p.empresa_id = NEW.empresa_id AND p.id = NEW.perfil_id;
  IF NOT FOUND THEN
    RETURN NEW;
  END IF;
  -- FOR SHARE: trava a linha do usuário até o fim da transação. Sem a trava,
  -- em READ COMMITTED, esta atribuição e uma troca de vínculo simultânea
  -- (tg_usuario_vinculo) se veriam mutuamente no estado antigo e as duas
  -- passariam (ISO-E01-06). Com ela, a troca de vínculo espera e, ao rodar,
  -- já enxerga a atribuição.
  SELECT u.fornecedor_id, u.locatario_id INTO v_fornecedor, v_locatario FROM app.usuario u
  WHERE u.empresa_id = NEW.empresa_id AND u.id = NEW.usuario_id
  FOR SHARE;
  IF NOT FOUND THEN
    RETURN NEW;
  END IF;
  IF (v_codigo = 'prestador') IS DISTINCT FROM (v_fornecedor IS NOT NULL)
     OR (v_codigo = 'locatario') IS DISTINCT FROM (v_locatario IS NOT NULL) THEN
    RAISE EXCEPTION 'perfil % incompatível com o vínculo externo do usuário', v_codigo USING ERRCODE = '23514';
  END IF;
  RETURN NEW;
END
$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_usuario_perfil_externo BEFORE INSERT OR UPDATE ON app.usuario_perfil
  FOR EACH ROW EXECUTE FUNCTION app.tg_usuario_perfil_externo();
-- +goose StatementEnd

-- +goose StatementBegin
-- Permissão avulsa concedida a um usuário (RF-ADM-005; sumário 3.2 e 3.6), como
-- pessoa.saude.ler, que nenhum perfil padrão tem. Auditada pelo preparar_tabela:
-- cada concessão e retirada é evento (RF-ADM-011). Retirar é exclusão lógica.
-- O formato é conferido aqui; a existência do código no catálogo, no Go.
CREATE TABLE app.usuario_permissao (
  empresa_id  uuid NOT NULL,
  id          uuid NOT NULL DEFAULT app.uuid_v7(),
  usuario_id  uuid NOT NULL,
  permissao   text NOT NULL,                -- <módulo>.<ação>; exceção aceita: pessoa.saude.ler
  versao bigint NOT NULL DEFAULT 1, alterado_xid xid8 NOT NULL DEFAULT pg_current_xact_id(),
  criado_em timestamptz NOT NULL DEFAULT now(), criado_por uuid,
  atualizado_em timestamptz NOT NULL DEFAULT now(), atualizado_por uuid, excluido_em timestamptz,
  CONSTRAINT pk_usuario_permissao PRIMARY KEY (empresa_id, id),
  CONSTRAINT fk_usuario_permissao_usuario FOREIGN KEY (empresa_id, usuario_id) REFERENCES app.usuario (empresa_id, id),
  -- a ação começa por letra e aceita dígito depois dela (con.gerenciar_nr13);
  -- os 117 códigos do catálogo de RF-ADM-005 passam neste padrão
  CONSTRAINT ck_usuario_permissao_formato CHECK (permissao ~ '^[a-z]{2,3}\.[a-z][a-z0-9_]*[a-z0-9]$'
                                                 OR permissao = 'pessoa.saude.ler')
);
-- Uma concessão ativa por usuário e permissão; também o índice da leitura das
-- permissões do usuário, feita a cada requisição.
CREATE UNIQUE INDEX uq_usuario_permissao ON app.usuario_permissao (empresa_id, usuario_id, permissao)
  WHERE excluido_em IS NULL;
CALL app.preparar_tabela('app.usuario_permissao');
-- +goose StatementEnd

-- +goose StatementBegin
-- Usuário externo nunca recebe permissão vedada, nem avulsa. Usuário não
-- encontrado (inclusive de outra empresa): deixa passar, e a FK recusa (23503).
CREATE OR REPLACE FUNCTION app.tg_usuario_permissao_externo() RETURNS trigger
  LANGUAGE plpgsql
AS $$
DECLARE
  v_fornecedor uuid;
  v_locatario  uuid;
BEGIN
  IF NEW.excluido_em IS NOT NULL OR NOT (NEW.permissao = ANY (app.permissoes_vedadas_externo())) THEN
    RETURN NEW;
  END IF;
  -- FOR SHARE pelo mesmo motivo de tg_usuario_perfil_externo (ISO-E01-06)
  SELECT u.fornecedor_id, u.locatario_id INTO v_fornecedor, v_locatario FROM app.usuario u
  WHERE u.empresa_id = NEW.empresa_id AND u.id = NEW.usuario_id
  FOR SHARE;
  IF NOT FOUND THEN
    RETURN NEW;
  END IF;
  IF v_fornecedor IS NOT NULL OR v_locatario IS NOT NULL THEN
    RAISE EXCEPTION 'permissão % é vedada a usuário externo', NEW.permissao USING ERRCODE = '23514';
  END IF;
  RETURN NEW;
END
$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_usuario_permissao_externo BEFORE INSERT OR UPDATE ON app.usuario_permissao
  FOR EACH ROW EXECUTE FUNCTION app.tg_usuario_permissao_externo();
-- +goose StatementEnd

-- +goose StatementBegin
-- Trigger simétrico em app.usuario: quem já tem perfil não troca de vínculo
-- externo (seria um técnico interno virando prestador com os perfis antigos),
-- e a identidade, uma vez definida, não muda — trocar identidade_id seria
-- entregar a conta a outro login. Para trocar, desliga-se o usuário e cria-se
-- outro.
--
-- A consulta a usuario_perfil roda sob o RLS de quem grava. Um usuário externo
-- não enxerga usuario_perfil (só interno), mas também não chega aqui com
-- vínculo trocado: o WITH CHECK da policy escopo_externo de usuario o recusa
-- antes.
CREATE OR REPLACE FUNCTION app.tg_usuario_vinculo() RETURNS trigger
  LANGUAGE plpgsql
AS $$
BEGIN
  IF OLD.identidade_id IS NOT NULL AND NEW.identidade_id IS DISTINCT FROM OLD.identidade_id THEN
    RAISE EXCEPTION 'a identidade do usuário é imutável depois de definida' USING ERRCODE = '42501';
  END IF;
  IF (NEW.fornecedor_id IS DISTINCT FROM OLD.fornecedor_id OR NEW.locatario_id IS DISTINCT FROM OLD.locatario_id)
     AND EXISTS (SELECT 1 FROM app.usuario_perfil up
                 WHERE up.empresa_id = NEW.empresa_id AND up.usuario_id = NEW.id AND up.excluido_em IS NULL) THEN
    RAISE EXCEPTION 'o usuário tem perfis; o vínculo externo não pode mudar' USING ERRCODE = '23514';
  END IF;
  -- interno com permissão avulsa vedada (ex.: pessoa.saude.ler) não vira externo
  IF (NEW.fornecedor_id IS NOT NULL OR NEW.locatario_id IS NOT NULL)
     AND (NEW.fornecedor_id IS DISTINCT FROM OLD.fornecedor_id OR NEW.locatario_id IS DISTINCT FROM OLD.locatario_id)
     AND EXISTS (SELECT 1 FROM app.usuario_permissao pu
                 WHERE pu.empresa_id = NEW.empresa_id AND pu.usuario_id = NEW.id AND pu.excluido_em IS NULL
                   AND pu.permissao = ANY (app.permissoes_vedadas_externo())) THEN
    RAISE EXCEPTION 'o usuário tem permissão avulsa vedada a usuário externo' USING ERRCODE = '23514';
  END IF;
  RETURN NEW;
END
$$;
-- +goose StatementEnd

-- +goose StatementBegin
CREATE TRIGGER trg_usuario_vinculo BEFORE UPDATE OF identidade_id, fornecedor_id, locatario_id ON app.usuario
  FOR EACH ROW EXECUTE FUNCTION app.tg_usuario_vinculo();
-- +goose StatementEnd

-- ---------------------------------------------------- conferência do escopo

-- +goose StatementBegin
-- Versão completa de app.conferir_escopo() (seção 2.9.3), agora que app.usuario
-- existe. Além do que a 00002 conferia, lê a linha do usuário (sob as policies
-- de quem chama) e recusa a transação se o escopo da sessão não bater com o
-- cadastro: prestador marcado como interno por engano, usuário que perdeu o
-- vínculo com o fornecedor durante a sessão, usuário bloqueado ou desligado.
--
-- Custo: uma busca pela chave primária por transação, no mesmo lote do
-- set_config (seção 2.4).
CREATE OR REPLACE FUNCTION app.conferir_escopo() RETURNS void
  LANGUAGE plpgsql STABLE
AS $$
DECLARE
  v_fornecedor uuid;
  v_locatario  uuid;
  v_situacao   text;
BEGIN
  IF app.escopo_atual() IS NULL
     OR app.escopo_atual() NOT IN ('interno', 'prestador', 'locatario', 'link') THEN
    RAISE EXCEPTION 'escopo ausente ou inválido' USING ERRCODE = '42501';
  END IF;

  IF app.empresa_atual() IS NULL THEN
    RAISE EXCEPTION 'empresa ausente no contexto' USING ERRCODE = '42501';
  END IF;

  IF app.usuario_atual() IS NULL THEN
    -- jobs e serviço de autenticação (interno) e link de uso único (link)
    IF app.escopo_atual() = 'link' AND app.link_hash_atual() IS NULL THEN
      RAISE EXCEPTION 'link sem segredo' USING ERRCODE = '42501';
    END IF;
    IF app.escopo_atual() IN ('prestador', 'locatario') THEN
      RAISE EXCEPTION 'escopo externo sem usuário' USING ERRCODE = '42501';
    END IF;
    RETURN;
  END IF;

  SELECT u.fornecedor_id, u.locatario_id, u.situacao
    INTO v_fornecedor, v_locatario, v_situacao
  FROM app.usuario u
  WHERE u.empresa_id = app.empresa_atual() AND u.id = app.usuario_atual() AND u.excluido_em IS NULL;

  IF NOT FOUND OR v_situacao <> 'ativo'
     OR v_fornecedor IS DISTINCT FROM app.fornecedor_atual()
     OR v_locatario  IS DISTINCT FROM app.locatario_atual()
     -- o CASE vai entre parênteses: dentro de um IF do PL/pgSQL, o primeiro
     -- THEN sem parêntese encerra a condição
     OR app.escopo_atual() <> (CASE WHEN v_fornecedor IS NOT NULL THEN 'prestador'
                                    WHEN v_locatario  IS NOT NULL THEN 'locatario'
                                    ELSE 'interno' END) THEN
    RAISE EXCEPTION 'escopo da sessão diverge do cadastro do usuário' USING ERRCODE = '42501';
  END IF;
END
$$;
-- +goose StatementEnd

-- +goose Down

-- +goose StatementBegin
SET ROLE tenaz_dono;

-- conferir_escopo volta à versão da fundação (00002), que não lê app.usuario.
CREATE OR REPLACE FUNCTION app.conferir_escopo() RETURNS void
  LANGUAGE plpgsql STABLE
AS $$
BEGIN
  IF app.escopo_atual() IS NULL
     OR app.escopo_atual() NOT IN ('interno', 'prestador', 'locatario', 'link') THEN
    RAISE EXCEPTION 'escopo ausente ou inválido' USING ERRCODE = '42501';
  END IF;

  IF app.empresa_atual() IS NULL THEN
    RAISE EXCEPTION 'empresa ausente no contexto' USING ERRCODE = '42501';
  END IF;

  IF app.usuario_atual() IS NULL THEN
    IF app.escopo_atual() = 'link' AND app.link_hash_atual() IS NULL THEN
      RAISE EXCEPTION 'link sem segredo' USING ERRCODE = '42501';
    END IF;
    IF app.escopo_atual() IN ('prestador', 'locatario') THEN
      RAISE EXCEPTION 'escopo externo sem usuário' USING ERRCODE = '42501';
    END IF;
  END IF;
END
$$;

DROP TABLE IF EXISTS app.usuario_permissao;
DROP TABLE IF EXISTS app.usuario_perfil;
DROP TABLE IF EXISTS app.perfil;
-- credenciais: dono tenaz_credencial
DROP TRIGGER IF EXISTS trg_acesso_empresa ON app.usuario;
SET ROLE tenaz_credencial;
DROP TABLE IF EXISTS plataforma.acesso_empresa;
DROP FUNCTION IF EXISTS plataforma.tg_acesso_empresa();
SET ROLE tenaz_dono;
-- a policy de pessoa consulta app.usuario e prende a tabela
DROP POLICY IF EXISTS escopo_externo ON app.pessoa;
DROP TABLE IF EXISTS app.usuario;
DROP TABLE IF EXISTS app.pessoa;
SET ROLE tenaz_credencial;
DROP TABLE IF EXISTS plataforma.identidade;
SET ROLE tenaz_dono;
DROP TABLE IF EXISTS app.feriado_unidade;
DROP TABLE IF EXISTS app.turno_unidade;
DROP TABLE IF EXISTS app.unidade;
DROP TABLE IF EXISTS app.empresa;
DROP FUNCTION IF EXISTS app.tg_usuario_vinculo();
DROP FUNCTION IF EXISTS app.tg_usuario_colunas_externo();
DROP FUNCTION IF EXISTS app.tg_perfil_tipo();
DROP FUNCTION IF EXISTS app.tg_usuario_permissao_externo();
DROP FUNCTION IF EXISTS app.permissoes_vedadas_externo();
DROP FUNCTION IF EXISTS app.tg_usuario_perfil_externo();
DROP FUNCTION IF EXISTS app.tg_padrao_empresa();
REVOKE USAGE ON SCHEMA plataforma FROM tenaz_app, tenaz_job;
RESET ROLE;
REVOKE USAGE ON SCHEMA public FROM tenaz_dono, tenaz_app, tenaz_job, tenaz_leitura, tenaz_sinc;
-- o papel tenaz_credencial não é removido (objeto do cluster); perde os acessos
REVOKE ALL ON SCHEMA plataforma FROM tenaz_credencial;
REVOKE USAGE ON SCHEMA public, app FROM tenaz_credencial;
-- +goose StatementEnd
