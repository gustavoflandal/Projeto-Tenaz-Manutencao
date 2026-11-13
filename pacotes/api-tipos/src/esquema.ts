// Gerado por pacotes/api-tipos/gerar.mjs a partir de backend/api/openapi.yaml.
// Não edite à mão: mude o contrato no Go, rode `bash ferramentas/tarefas.sh contrato`
// e depois `bash ferramentas/tarefas.sh api-tipos`.

export interface paths {
    "/oauth/authorize": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Início do OAuth 2.1 com PKCE do app
         * @description Confere client_id e redirect_uri e redireciona (302) para a tela de autorização da web, /autorizar-app, com os mesmos parâmetros. Lá a pessoa entra com a sessão web normal (senha e MFA), escolhe a empresa e a web chama POST /oauth/authorize. Cliente ou URI de retorno inválidos respondem 422 e nunca redirecionam.
         */
        get: operations["iniciarAutorizacaoApp"];
        put?: never;
        /**
         * Autorizar o app na empresa escolhida
         * @description Chamada pela tela /autorizar-app da web, com a sessão web (cookie, X-CSRF-Token e Origin). Emite o código de uso único (60 s) e devolve a redirect_uri com code e state. Perfil sem mob.usar_app (o locatario) recebe a redirect_uri com error=access_denied. Empresa fora do acesso: 404.
         */
        post: operations["autorizarApp"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/oauth/revoke": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Revogar o refresh token (sair no app)
         * @description RFC 7009: revoga a família de refresh tokens e a sessão do aparelho. O refresh token no corpo é a credencial. Token desconhecido também responde 200, para não servir de oráculo.
         */
        post: operations["revogarTokenApp"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/oauth/token": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Trocar o código ou renovar o refresh token
         * @description grant_type=authorization_code (code, redirect_uri, client_id, code_verifier) ou grant_type=refresh_token (refresh_token, client_id). Cada renovação devolve refresh token novo; um refresh token já usado e reapresentado revoga a família e a sessão do aparelho. Falha: 401 credencial_invalida com error=invalid_grant (RFC 6749). O app serializa as renovações.
         */
        post: operations["emitirTokensApp"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/convites/aceitar": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Aceitar o convite
         * @description Liga o login ao usuário convidado e o ativa; o link é gasto. Depois, a pessoa entra por POST /v1/auth/login. 410 link_email_invalido; 422 senha_fraca; 409 conflito se o e-mail já tem usuário na empresa.
         */
        post: operations["aceitarConvite"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/convites/consultar": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Consultar o convite do link
         * @description Confere o link sem gastá-lo, para a tela de aceite saber se pede senha. O token vai no corpo, nunca na URL. 410 link_email_invalido.
         */
        post: operations["consultarConvite"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/login": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Entrar com e-mail e senha
         * @description Cria a sessão web e devolve o cookie de sessão e o csrf_token. Com MFA ativo, responde 401 mfa_obrigatorio com desafio_id, e o login termina em POST /v1/auth/mfa/verificar. Toda falha de credencial responde 401 credencial_invalida, sem dizer se o e-mail existe (RF-ADM-006).
         */
        post: operations["entrar"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/logout": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Sair
         * @description Encerra a sessão no servidor (web ou app; no app, com a família de refresh tokens) e apaga o cookie.
         */
        post: operations["sair"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/mfa/confirmar": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Confirmar a ativação do segundo fator
         * @description Ativa o segundo fator com o primeiro código certo. 401 mfa_codigo_invalido.
         */
        post: operations["confirmarMFA"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/mfa/desativar": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Desativar o segundo fator
         * @description Exige um código válido. Código errado soma uma falha ao bloqueio progressivo, como no login, e código de um passo TOTP já usado vale como errado: 401 mfa_codigo_invalido. 409 transicao_invalida com estado_atual "inativo" se o segundo fator não está ativo, e com estado_atual "exigido_pela_empresa" enquanto a identidade tiver acesso a alguma empresa que exige MFA (RN-ADM-027).
         */
        post: operations["desativarMFA"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/mfa/iniciar": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Iniciar a ativação do segundo fator
         * @description Gera o segredo TOTP pendente e devolve o segredo e a URI do QR code, uma única vez. Só vale depois de POST /v1/auth/mfa/confirmar. 409 transicao_invalida se o segundo fator já está ativo.
         */
        post: operations["iniciarMFA"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/mfa/verificar": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Concluir o login com o segundo fator
         * @description Confere o código TOTP do desafio. Código errado responde 401 mfa_codigo_invalido com tentativas_restantes; esgotadas as 5 tentativas, ou vencido o desafio (5 minutos), o login recomeça pela senha.
         */
        post: operations["verificarMFA"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/recuperar-senha": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Pedir link de redefinição de senha
         * @description Envia o link (válido por 1 hora) se o e-mail tiver login. Responde 202 igual, exista ou não o e-mail, para não revelar quais e-mails existem.
         */
        post: operations["recuperarSenha"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/redefinir-senha": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Redefinir a senha com o link do e-mail
         * @description O token vem do fragmento do link (#token=...) e vai no corpo, nunca na URL. Redefinir encerra todas as sessões da pessoa, em todas as empresas. 410 link_email_invalido; 422 senha_fraca.
         */
        post: operations["redefinirSenha"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/sessao": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Sessão atual
         * @description Usuário, empresa, escopo, perfis, permissões, empresas do seletor e, na web, o csrf_token.
         */
        get: operations["obterSessao"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/trocar-empresa": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Trocar de empresa
         * @description Emite credencial NOVA na empresa escolhida e revoga a atual: uma credencial nunca vale para duas empresas (seção 4.7). Na web, cookie e csrf_token novos; no app, tokens novos. Empresa fora do acesso: 404.
         */
        post: operations["trocarEmpresa"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/auth/trocar-senha": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Trocar a própria senha
         * @description Confere a senha atual (errada conta como falha de login) e encerra as outras sessões da pessoa; a atual continua. 422 senha_fraca.
         */
        post: operations["trocarSenha"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/empresa": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Dados da empresa
         * @description A empresa da credencial. Traz ETag.
         */
        get: operations["obterEmpresa"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        /**
         * Alterar dados cadastrais
         * @description JSON Merge Patch com If-Match. Plano e módulo SRV mudam pela cobrança (RF-ADM-009), não aqui: 422 campo_somente_leitura. mfa_obrigatorio liga o segundo fator obrigatório: quem não o tem passa a entrar com sessão restrita à ativação (RN-ADM-027).
         */
        patch: operations["alterarEmpresa"];
        trace?: never;
    };
    "/v1/empresas": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Criar a conta da empresa (cadastro pelo site)
         * @description Cria a empresa com a unidade inicial, os 12 perfis padrão com a matriz de RF-ADM-005 e o responsável como administrador convidado, e envia o convite, que confirma o e-mail antes do primeiro acesso (RF-ADM-001; RF-SRV-001 com srv_habilitado). 409 conflito em CNPJ já cadastrado.
         */
        post: operations["criarEmpresa"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/eventos-auditoria": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Trilha de auditoria
         * @description Filtros por entidade, registro, usuário, ação e período; paginada por cursor, do mais recente (RF-ADM-011).
         */
        get: operations["listarEventosAuditoria"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/eventos-auditoria/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Evento de auditoria
         * @description Com os valores antes e depois. O evento é imutável.
         */
        get: operations["obterEventoAuditoria"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/perfis": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Listar perfis
         * @description Os 12 perfis padrão e os personalizados, com as permissões. Lista curta, inteira numa página.
         */
        get: operations["listarPerfis"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/permissoes": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Catálogo de permissões
         * @description O catálogo inteiro de RF-ADM-005, com fase e marcação de sensível. Lista curta, inteira numa página.
         */
        get: operations["listarPermissoes"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/pessoas": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Listar pessoas
         * @description Paginada por cursor, ordem -criado_em. O CPF sai mascarado; o inteiro só no detalhe.
         */
        get: operations["listarPessoas"];
        put?: never;
        /**
         * Cadastrar pessoa
         * @description Com ou sem usuário (RF-ADM-015). 409 conflito em CPF repetido na empresa.
         */
        post: operations["criarPessoa"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/pessoas/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Detalhar pessoa
         * @description Com o CPF inteiro; a leitura é registrada na trilha. O titular lê a própria pessoa sem adm.ler_pessoas. Traz ETag.
         */
        get: operations["obterPessoa"];
        put?: never;
        post?: never;
        /**
         * Inativar pessoa
         * @description A pessoa nunca é apagada: é inativada (seção 2.12). Com If-Match. 409 transicao_invalida se a pessoa já está inativa ou se está ligada a usuário ativo ou convidado: o usuário é inativado antes (RN-ADM-021).
         */
        delete: operations["inativarPessoa"];
        options?: never;
        head?: never;
        /**
         * Alterar pessoa
         * @description JSON Merge Patch com If-Match.
         */
        patch: operations["alterarPessoa"];
        trace?: never;
    };
    "/v1/pessoas/{id}/reativar": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Reativar pessoa
         * @description inativa → ativa, com If-Match. 409 transicao_invalida se já está ativa.
         */
        post: operations["reativarPessoa"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/unidades": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Listar unidades
         * @description Lista curta, inteira numa página (tem_mais falso), com as unidades ativas e inativas.
         */
        get: operations["listarUnidades"];
        put?: never;
        /**
         * Criar unidade
         * @description Com pelo menos um turno (RF-ADM-002). O Essencial tem 1 unidade: a segunda responde 403 plano_nao_inclui_recurso com plano_necessario. 409 conflito em código repetido.
         */
        post: operations["criarUnidade"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/unidades/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Detalhar unidade
         * @description Com os turnos. Traz ETag.
         */
        get: operations["obterUnidade"];
        put?: never;
        post?: never;
        /**
         * Inativar unidade
         * @description ativa → inativa, com If-Match; a unidade nunca é apagada (seção 2.12). A última unidade ativa não é inativada (422 validacao_falhou). 409 transicao_invalida se já está inativa.
         */
        delete: operations["inativarUnidade"];
        options?: never;
        head?: never;
        /**
         * Alterar unidade
         * @description Nome, código, endereço, fuso, CNPJ e coordenadas. JSON Merge Patch com If-Match.
         */
        patch: operations["alterarUnidade"];
        trace?: never;
    };
    "/v1/unidades/{id}/feriados": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Listar feriados da unidade
         * @description Lista curta, inteira numa página, com filtro de período opcional.
         */
        get: operations["listarFeriados"];
        put?: never;
        /**
         * Incluir feriado
         * @description 409 conflito em data repetida na unidade.
         */
        post: operations["incluirFeriado"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/unidades/{id}/feriados/{feriado_id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        post?: never;
        /**
         * Remover feriado
         * @description Com If-Match (ETag do feriado, que vem na lista).
         */
        delete: operations["removerFeriado"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/unidades/{id}/reativar": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Reativar unidade
         * @description inativa → ativa, com If-Match; respeita o limite de unidades do plano. 409 transicao_invalida se já está ativa.
         */
        post: operations["reativarUnidade"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/unidades/{id}/turnos": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Incluir turno
         * @description Turno por dia da semana no calendário da unidade (RF-ADM-002).
         */
        post: operations["incluirTurno"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/unidades/{id}/turnos/{turno_id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        post?: never;
        /**
         * Remover turno
         * @description Com If-Match. O último turno da unidade não é removido (422 validacao_falhou).
         */
        delete: operations["removerTurno"];
        options?: never;
        head?: never;
        /**
         * Alterar turno
         * @description JSON Merge Patch com If-Match (ETag do turno).
         */
        patch: operations["alterarTurno"];
        trace?: never;
    };
    "/v1/usuarios": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Listar usuários
         * @description Paginada por cursor, ordem -criado_em.
         */
        get: operations["listarUsuarios"];
        put?: never;
        /**
         * Criar usuário e enviar convite
         * @description O usuário nasce convidado e o convite (7 dias) vai por e-mail (RF-ADM-004). Regras anti-escalada: só se atribui perfil cujas permissões se tem (403 permissao_negada). prestador exige fornecedor_id; locatario, locatario_id. 409 conflito em e-mail já usado na empresa.
         */
        post: operations["criarUsuario"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/usuarios/{id}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Detalhar usuário
         * @description Com perfis e permissões avulsas. Traz ETag.
         */
        get: operations["obterUsuario"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        /**
         * Alterar perfis e dados do usuário
         * @description JSON Merge Patch com If-Match. perfis é a lista completa e substitui a atual. Ninguém altera os próprios perfis; só administrador altera administrador; 409 ultimo_administrador ao tirar o perfil do último administrador ativo. E-mail e vínculo externo não mudam (422 campo_somente_leitura no e-mail).
         */
        patch: operations["alterarUsuario"];
        trace?: never;
    };
    "/v1/usuarios/{id}/desativar": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Desativar usuário
         * @description situacao → inativo, com If-Match; as sessões e os refresh tokens do usuário caem na mesma transação. 409 ultimo_administrador; 409 transicao_invalida se já está inativo.
         */
        post: operations["desativarUsuario"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/usuarios/{id}/encerrar-sessoes": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Encerrar as sessões do usuário
         * @description Revoga as sessões web e do app e os refresh tokens do usuário nesta empresa; o logout vai para a trilha.
         */
        post: operations["encerrarSessoesUsuario"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/usuarios/{id}/permissoes": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Conceder permissão avulsa
         * @description Permissão do catálogo de RF-ADM-005 concedida a este usuário (usuario_permissao), como os.encerrar. Ninguém concede a si mesmo; só concede quem tem a permissão, salvo as duas que o administrador concede sem ter (os.encerrar e est.decidir_destino); usuário externo nunca recebe permissão vedada. A concessão vai para a trilha.
         */
        post: operations["concederPermissaoUsuario"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/usuarios/{id}/permissoes/{permissao}": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        post?: never;
        /**
         * Retirar permissão avulsa
         * @description A retirada vai para a trilha. Permissão que o usuário não tem: 404.
         */
        delete: operations["retirarPermissaoUsuario"];
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/usuarios/{id}/reativar": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Reativar usuário
         * @description inativo → ativo, com If-Match; quem nunca aceitou o convite volta a convidado e recebe convite novo. 409 transicao_invalida se não está inativo.
         */
        post: operations["reativarUsuario"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/usuarios/{id}/reenviar-convite": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        get?: never;
        put?: never;
        /**
         * Reenviar o convite
         * @description Revoga o link anterior e envia outro (7 dias). 409 transicao_invalida se o usuário não está convidado; 503 servico_indisponivel se o e-mail não saiu.
         */
        post: operations["reenviarConvite"];
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
    "/v1/usuarios/eu": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Próprio usuário
         * @description Usuário da credencial, com perfis e permissões avulsas. Vale para os perfis externos. Traz ETag.
         */
        get: operations["obterUsuarioEu"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        /**
         * Alterar os próprios dados
         * @description Só nome e telefone. JSON Merge Patch com If-Match.
         */
        patch: operations["alterarUsuarioEu"];
        trace?: never;
    };
    "/v1/versao": {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        /**
         * Versão publicada
         * @description Devolve a versão, o commit curto e a data do build do release em produção (RNF-039).
         */
        get: operations["obterVersao"];
        put?: never;
        post?: never;
        delete?: never;
        options?: never;
        head?: never;
        patch?: never;
        trace?: never;
    };
}
export type webhooks = Record<string, never>;
export interface components {
    schemas: {
        Autorizacao: {
            /** @description redirect_uri com code e state (ou error=access_denied), para a web navegar até ela e o sistema entregar ao app */
            uri_redirecionamento: string;
        };
        CampoErro: {
            /**
             * @description Campo em notação de ponto
             * @example responsavel.email
             */
            campo: string;
            /**
             * @description Vocabulário curto: obrigatorio, valor_invalido, tamanho_maximo, tamanho_minimo, formato_invalido, fora_da_faixa, nao_encontrado, duplicado, inativo
             * @example obrigatorio
             */
            codigo: string;
            /**
             * @description O que está errado e o que fazer
             * @example Informe o e-mail.
             */
            mensagem: string;
            /** @description Valores aceitos, quando o campo é uma enumeração */
            valores_aceitos?: string[] | null;
        };
        ConfiguracaoMFA: {
            /** @description Segredo para digitar no aplicativo autenticador */
            segredo_base32: string;
            /** @description URI otpauth:// para o QR code */
            uri_otpauth: string;
        };
        ContaCriada: {
            /** @description Falso: a conta existe, mas o e-mail do convite não saiu; peça o reenvio ao suporte */
            convite_enviado: boolean;
        };
        Convite: {
            /** @description E-mail convidado */
            email: string;
            /** @description Verdadeiro: o e-mail ainda não tem login e a tela pede a senha nova; falso: entra com a senha que já tem */
            pede_senha: boolean;
        };
        ConviteAceito: {
            /** @description Falso: o e-mail já entrava em outra empresa e usa a senha que tem */
            identidade_nova: boolean;
        };
        Empresa: {
            /** Format: date-time */
            atualizado_em: string;
            /** @description 14 dígitos, sem máscara */
            cnpj: string | null;
            /** Format: date-time */
            criado_em: string;
            /** @example America/Sao_Paulo */
            fuso_horario: string;
            /** Format: uuid */
            id: string;
            /** @description Segundo fator exigido de todos os usuários (RN-ADM-027): quem não o tem entra com sessão restrita à ativação */
            mfa_obrigatorio: boolean;
            nome_fantasia: string | null;
            /** @description essencial, profissional ou enterprise; nulo na conta só do plano Serviços */
            plano: string | null;
            razao_social: string;
            /** Format: int32 */
            retencao_auditoria_anos: number;
            /** @enum {string} */
            situacao: "teste" | "ativa" | "suspensa" | "cancelada";
            srv_habilitado: boolean;
            /** Format: date */
            teste_ate: string | null;
            /** Format: int64 */
            versao: number;
        };
        EmpresaAcesso: {
            /** @description Empresa da credencial atual */
            atual: boolean;
            /**
             * Format: uuid
             * @description Identificador para empresa_escolhida_id
             */
            id: string;
            nome: string;
        };
        EmpresaSessao: {
            nome: string;
            /**
             * @description Nulo na conta só do plano Serviços
             * @enum {string|null}
             */
            plano: "essencial" | "profissional" | "enterprise" | null;
            /** @description Módulo SRV habilitado (conta de subcontratada) */
            srv_habilitado: boolean;
        };
        Endereco: {
            bairro?: string;
            cep?: string;
            complemento?: string;
            logradouro?: string;
            municipio?: string;
            numero?: string;
            uf?: string;
        };
        EntradaAutorizarPOSTBody: {
            client_id: string;
            code_challenge: string;
            /** @enum {string} */
            code_challenge_method: "S256";
            /**
             * Format: uuid
             * @description Empresa da sessão do app, entre as de GET /v1/auth/sessao
             */
            empresa_escolhida_id: string;
            redirect_uri: string;
            /** @enum {string} */
            response_type: "code";
            state?: string;
        };
        EntradaCodigoMFABody: {
            codigo: string;
        };
        EntradaConviteAceitarBody: {
            /** @description Obrigatória quando pede_senha; ignorada quando o e-mail já tem login */
            senha?: string;
            token: string;
        };
        EntradaConviteConsultarBody: {
            /** @description Token do link do convite, lido do fragmento #token= */
            token: string;
        };
        EntradaFeriado: {
            /** Format: date */
            data: string;
            descricao: string;
        };
        EntradaLoginBody: {
            /** @description E-mail do login */
            email: string;
            /**
             * Format: uuid
             * @description Empresa da sessão, entre as da pessoa; ausente: a da última sessão
             */
            empresa_escolhida_id?: string;
            senha: string;
        };
        EntradaMFABody: {
            /** @description Código de 6 dígitos do aplicativo autenticador */
            codigo: string;
            /** @description desafio_id da resposta 401 mfa_obrigatorio */
            desafio_id: string;
            /** Format: uuid */
            empresa_escolhida_id?: string;
        };
        EntradaPermissao: {
            /**
             * @description Código do catálogo de RF-ADM-005
             * @example os.encerrar
             */
            permissao: string;
        };
        EntradaPessoa: {
            cargo?: string;
            /** @description Com ou sem máscara */
            cpf?: string;
            /** Format: date */
            data_nascimento?: string;
            email?: string;
            /** Format: uuid */
            fornecedor_id?: string;
            matricula?: string;
            nome: string;
            registro_profissional?: string;
            telefone?: string;
            /** Format: uuid */
            unidade_id?: string;
            /**
             * @description Ausente: empregado
             * @enum {string}
             */
            vinculo?: "empregado" | "terceiro" | "autonomo";
        };
        EntradaRecuperarSenhaBody: {
            email: string;
        };
        EntradaRedefinirSenhaBody: {
            nova_senha: string;
            /** @description Token do link do e-mail, lido do fragmento #token= */
            token: string;
        };
        EntradaRevogarBody: {
            client_id?: string;
            token: string;
            /** @enum {string} */
            token_type_hint?: "refresh_token";
        };
        EntradaTokenBody: {
            client_id: string;
            code?: string;
            code_verifier?: string;
            /** @enum {string} */
            grant_type: "authorization_code" | "refresh_token";
            redirect_uri?: string;
            refresh_token?: string;
        };
        EntradaTrocarEmpresaBody: {
            /**
             * Format: uuid
             * @description Empresa de destino, entre as de GET /v1/auth/sessao
             */
            empresa_escolhida_id: string;
        };
        EntradaTrocarSenhaBody: {
            nova_senha: string;
            senha_atual: string;
        };
        EntradaTurno: {
            /**
             * Format: int32
             * @description 1 = segunda ... 7 = domingo
             */
            dia_semana: number;
            /** @example 17:00 */
            hora_fim: string;
            /** @example 08:00 */
            hora_inicio: string;
            nome?: string;
        };
        EntradaUnidade: {
            cnpj?: string;
            codigo: string;
            endereco?: components["schemas"]["Endereco"];
            fuso_horario?: string;
            /** Format: double */
            latitude?: number;
            /** Format: double */
            longitude?: number;
            nome: string;
            /** @description Só na criação; depois, pelas rotas de turno. Obrigatório em POST /v1/unidades; no cadastro da conta, ausente vale segunda a sexta das 08:00 às 17:00 */
            turnos?: components["schemas"]["EntradaTurno"][] | null;
        };
        EntradaUsuario: {
            /** @example 125.50 */
            custo_hora?: string;
            email: string;
            /**
             * Format: uuid
             * @description Obrigatório no perfil prestador
             */
            fornecedor_id?: string;
            /**
             * Format: uuid
             * @description Obrigatório no perfil locatario
             */
            locatario_id?: string;
            nome: string;
            /** @description Cria a pessoa junto; não use junto com pessoa_id */
            nova_pessoa?: components["schemas"]["EntradaPessoa"];
            /** @description Códigos dos perfis: administrador, gestor, tecnico... */
            perfis: string[] | null;
            /**
             * Format: uuid
             * @description Liga a uma pessoa existente; não use junto com nova_pessoa
             */
            pessoa_id?: string;
            telefone?: string;
            /** Format: uuid */
            unidade_padrao_id?: string;
        };
        Evento: {
            acao: string;
            /** @description Valores antes da alteração; nulo na criação */
            antes: unknown;
            /** @enum {string} */
            ator_tipo: "usuario" | "sistema" | "integracao" | "suporte_tenaz" | "link_execucao" | "sincronizacao_conexao";
            campos_alterados: string[] | null;
            /** Format: uuid */
            conexao_empresa_id: string;
            /** @description Valores depois da alteração; nulo na exclusão */
            depois: unknown;
            /** Format: uuid */
            dispositivo_id: string;
            entidade: string;
            entidade_chave: string | null;
            /** Format: uuid */
            entidade_id: string;
            /** Format: uuid */
            id: string;
            ip: string | null;
            /** Format: date-time */
            ocorrido_em: string;
            /** Format: date-time */
            ocorrido_no_dispositivo_em: string | null;
            /** @enum {string|null} */
            origem: "web" | "app" | "api" | "job" | "importacao" | null;
            request_id: string | null;
            /** Format: uuid */
            token_api_id: string;
            user_agent: string | null;
            /** Format: uuid */
            usuario_id: string;
        };
        Feriado: {
            /** Format: date */
            data: string;
            descricao: string;
            /** Format: uuid */
            id: string;
            /** Format: uuid */
            unidade_id: string;
            /** Format: int64 */
            versao: number;
        };
        ListaFeriado: {
            dados: components["schemas"]["Feriado"][] | null;
            paginacao: components["schemas"]["Paginacao"];
        };
        ListaPerfil: {
            dados: components["schemas"]["Perfil"][] | null;
            paginacao: components["schemas"]["Paginacao"];
        };
        ListaPermissao: {
            dados: components["schemas"]["Permissao"][] | null;
            paginacao: components["schemas"]["Paginacao"];
        };
        ListaResumoEvento: {
            dados: components["schemas"]["ResumoEvento"][] | null;
            paginacao: components["schemas"]["Paginacao"];
        };
        ListaResumoPessoa: {
            dados: components["schemas"]["ResumoPessoa"][] | null;
            paginacao: components["schemas"]["Paginacao"];
        };
        ListaResumoUnidade: {
            dados: components["schemas"]["ResumoUnidade"][] | null;
            paginacao: components["schemas"]["Paginacao"];
        };
        ListaResumoUsuario: {
            dados: components["schemas"]["ResumoUsuario"][] | null;
            paginacao: components["schemas"]["Paginacao"];
        };
        NovaEmpresa: {
            /** @description Com ou sem máscara; dígito verificador conferido */
            cnpj: string;
            /** @description IANA; ausente: America/Sao_Paulo */
            fuso_horario?: string;
            nome_fantasia?: string;
            /**
             * @description Ausente só com srv_habilitado (conta do plano Serviços)
             * @enum {string}
             */
            plano?: "essencial" | "profissional" | "enterprise";
            razao_social: string;
            responsavel: components["schemas"]["Responsavel"];
            /** @description Conta de subcontratada, com o módulo SRV (RF-SRV-001) */
            srv_habilitado?: boolean;
            /** @description Unidade inicial; sem turnos, segunda a sexta das 08:00 às 17:00 */
            unidade: components["schemas"]["EntradaUnidade"];
        };
        Paginacao: {
            /**
             * Format: int32
             * @description Itens por página
             * @example 50
             */
            limite: number;
            /** @description Valor para o parâmetro cursor da próxima página; nulo quando não há mais */
            proximo_cursor: string | null;
            /** @description Há mais itens depois desta página */
            tem_mais: boolean;
        };
        PatchEmpresa: {
            /** @description Com ou sem máscara */
            cnpj?: string;
            fuso_horario?: string;
            /** @description Exige segundo fator de todos os usuários da empresa */
            mfa_obrigatorio?: boolean;
            nome_fantasia?: string | null;
            razao_social?: string;
        };
        PatchEndereco: {
            bairro?: string | null;
            cep?: string | null;
            complemento?: string | null;
            logradouro?: string | null;
            municipio?: string | null;
            numero?: string | null;
            uf?: string | null;
        } | null;
        PatchEu: {
            nome?: string;
            telefone?: string | null;
        };
        PatchPessoa: {
            cargo?: string | null;
            cpf?: string | null;
            /** Format: date */
            data_nascimento?: string | null;
            email?: string | null;
            /** Format: uuid */
            fornecedor_id?: string | null;
            matricula?: string | null;
            nome?: string;
            registro_profissional?: string | null;
            telefone?: string | null;
            /** Format: uuid */
            unidade_id?: string | null;
            /** @enum {string} */
            vinculo?: "empregado" | "terceiro" | "autonomo";
        };
        PatchTurno: {
            /** Format: int32 */
            dia_semana?: number;
            hora_fim?: string;
            hora_inicio?: string;
            nome?: string | null;
        };
        PatchUnidade: {
            cnpj?: string | null;
            codigo?: string;
            /** @description Mesclado campo a campo; null apaga o endereço */
            endereco?: components["schemas"]["PatchEndereco"];
            fuso_horario?: string | null;
            /** Format: double */
            latitude?: number | null;
            /** Format: double */
            longitude?: number | null;
            nome?: string;
        };
        PatchUsuario: {
            custo_hora?: string | null;
            nome?: string;
            /** @description Lista completa dos perfis; substitui a atual */
            perfis?: string[] | null;
            /** Format: uuid */
            pessoa_id?: string | null;
            telefone?: string | null;
            /** Format: uuid */
            unidade_padrao_id?: string | null;
        };
        Perfil: {
            codigo: string;
            /** Format: uuid */
            id: string;
            nome: string;
            /** @description Perfil padrão do sumário, que não se altera */
            padrao: boolean;
            perfil_origem: string | null;
            permissoes: string[] | null;
            /** Format: int64 */
            versao: number;
        };
        PerfilAtribuido: {
            codigo: string;
            nome: string;
            /** Format: uuid */
            perfil_id: string;
            /** @description Perfil padrão de origem */
            perfil_origem: string;
            /** Format: uuid */
            unidade_id: string;
        };
        Permissao: {
            /** @example os.encerrar */
            codigo: string;
            descricao: string;
            /** @enum {string} */
            fase: "MVP" | "F2" | "F3";
            /** @description Grupo da seção 4.8 do documento 3 */
            grupo: string;
            sensivel: boolean;
        };
        Pessoa: {
            /** Format: date-time */
            atualizado_em: string;
            cargo: string | null;
            /** @description 11 dígitos, sem máscara */
            cpf: string | null;
            /** Format: date-time */
            criado_em: string;
            /** Format: date */
            data_nascimento: string | null;
            email: string | null;
            /** Format: uuid */
            fornecedor_id: string;
            /** Format: uuid */
            id: string;
            matricula: string | null;
            nome: string;
            registro_profissional: string | null;
            /** @enum {string} */
            situacao: "ativa" | "inativa";
            telefone: string | null;
            /** Format: uuid */
            unidade_id: string;
            /**
             * Format: uuid
             * @description Usuário ligado à pessoa, se houver
             */
            usuario_id: string;
            /** Format: int64 */
            versao: number;
            /** @enum {string} */
            vinculo: "empregado" | "terceiro" | "autonomo";
        };
        Problema: {
            /** @description Campos exigidos que faltaram */
            campos_faltantes?: string[] | null;
            /**
             * @description Código estável do erro, em snake_case; é o que o cliente trata
             * @example validacao_falhou
             */
            codigo: string;
            /**
             * Format: date-time
             * @description mfa_obrigatorio: fim da validade do desafio (5 minutos)
             */
            desafio_expira_em?: string;
            /** @description mfa_obrigatorio: desafio do segundo fator, segredo de 256 bits e uso único; vai no corpo de POST /v1/auth/mfa/verificar */
            desafio_id?: string;
            /** @description Explicação desta ocorrência e o que fazer */
            detail?: string;
            /** @description Erros por campo, com o campo em notação de ponto */
            erros?: components["schemas"]["CampoErro"][] | null;
            /** @description transicao_invalida: estado atual do registro */
            estado_atual?: string;
            /** @description Igual ao X-Request-Id; o suporte usa para localizar o log */
            id_rastreio?: string;
            /**
             * @description Caminho da requisição
             * @example /v1/pessoas
             */
            instance?: string;
            /** @description plano_nao_inclui_recurso: plano mínimo do recurso */
            plano_necessario?: string;
            /** @description versao_desatualizada: o recurso como está gravado, para o cliente mostrar a diferença */
            recurso_atual?: unknown;
            /**
             * Format: int64
             * @description Código de status HTTP
             * @example 422
             */
            status: number;
            /**
             * Format: int32
             * @description mfa_codigo_invalido: tentativas que restam no desafio; em zero, o login recomeça pela senha
             */
            tentativas_restantes?: number;
            /**
             * @description Resumo legível, fixo por código
             * @example Os dados enviados não passaram na validação
             */
            title: string;
            /** @description transicao_invalida: ações aceitas no estado atual */
            transicoes_permitidas?: string[] | null;
            /**
             * @description URI que documenta o código do erro
             * @example https://api.tenaz.com.br/erros/validacao_falhou
             */
            type: string;
            /**
             * Format: int64
             * @description versao_desatualizada: versão gravada do registro
             */
            versao_atual?: number;
        } & {
            [key: string]: unknown;
        };
        Responsavel: {
            email: string;
            nome: string;
            telefone?: string;
        };
        ResumoEvento: {
            acao: string;
            /** @enum {string} */
            ator_tipo: "usuario" | "sistema" | "integracao" | "suporte_tenaz" | "link_execucao" | "sincronizacao_conexao";
            campos_alterados: string[] | null;
            /** Format: uuid */
            conexao_empresa_id: string;
            /** Format: uuid */
            dispositivo_id: string;
            entidade: string;
            entidade_chave: string | null;
            /** Format: uuid */
            entidade_id: string;
            /** Format: uuid */
            id: string;
            ip: string | null;
            /** Format: date-time */
            ocorrido_em: string;
            /** Format: date-time */
            ocorrido_no_dispositivo_em: string | null;
            /** @enum {string|null} */
            origem: "web" | "app" | "api" | "job" | "importacao" | null;
            request_id: string | null;
            /** Format: uuid */
            token_api_id: string;
            /** Format: uuid */
            usuario_id: string;
        };
        ResumoPessoa: {
            cargo: string | null;
            /** @example ***.***.789-01 */
            cpf_mascarado: string | null;
            /** Format: date-time */
            criado_em: string;
            /** Format: uuid */
            fornecedor_id: string;
            /** Format: uuid */
            id: string;
            matricula: string | null;
            nome: string;
            /** @enum {string} */
            situacao: "ativa" | "inativa";
            /** Format: uuid */
            unidade_id: string;
            /** Format: uuid */
            usuario_id: string;
            /** Format: int64 */
            versao: number;
            /** @enum {string} */
            vinculo: "empregado" | "terceiro" | "autonomo";
        };
        ResumoUnidade: {
            /** Format: date-time */
            atualizado_em: string;
            cnpj: string | null;
            codigo: string;
            /** Format: date-time */
            criado_em: string;
            endereco: components["schemas"]["Endereco"];
            /** @description Nulo: herda o da empresa */
            fuso_horario: string | null;
            /** Format: uuid */
            id: string;
            /** Format: double */
            latitude: number | null;
            /** Format: double */
            longitude: number | null;
            nome: string;
            /** @enum {string} */
            situacao: "ativa" | "inativa";
            /** Format: int64 */
            versao: number;
        };
        ResumoUsuario: {
            /** Format: date-time */
            criado_em: string;
            email: string | null;
            /** Format: uuid */
            fornecedor_id: string;
            /** Format: uuid */
            id: string;
            /** Format: uuid */
            locatario_id: string;
            nome: string;
            /** @description Perfis de origem */
            perfis: string[] | null;
            /** Format: uuid */
            pessoa_id: string;
            /** @enum {string} */
            situacao: "convidado" | "ativo" | "inativo";
            telefone: string | null;
            /** Format: int64 */
            versao: number;
        };
        Sessao: {
            /**
             * @description Canal da credencial
             * @enum {string}
             */
            canal: "web" | "app";
            /** @description Token CSRF da sessão web, para o cabeçalho X-CSRF-Token em POST, PATCH, PUT e DELETE; nulo no app */
            csrf_token: string | null;
            empresa: components["schemas"]["EmpresaSessao"];
            /** @description Empresas em que a pessoa entra, para o seletor; a atual vem primeiro */
            empresas: components["schemas"]["EmpresaAcesso"][] | null;
            /**
             * @description Escopo de acesso da credencial (seção 4.9)
             * @enum {string}
             */
            escopo_acesso: "interno" | "prestador" | "locatario" | "link";
            /**
             * Format: date-time
             * @description Teto absoluto da sessão web (7 dias); a inatividade de 12 horas vence antes. Nulo fora do login e da troca de empresa
             */
            expira_em: string | null;
            /** @description A empresa exige segundo fator e a pessoa ainda não o ativou: a sessão só serve para ler a sessão, sair, trocar de empresa e ativar o MFA; as demais rotas respondem 401 nao_autenticado com mfa_pendente (RN-ADM-027) */
            mfa_pendente: boolean;
            /** @description Códigos de origem dos perfis ativos */
            perfis: string[] | null;
            /** @description Permissões efetivas, em ordem alfabética (perfis e permissões avulsas) */
            permissoes: string[] | null;
            usuario: components["schemas"]["UsuarioSessao"];
        };
        SessoesEncerradas: {
            /** Format: int32 */
            sessoes_encerradas: number;
        };
        TokensApp: {
            /** @description JWT EdDSA, 15 minutos */
            access_token: string;
            /**
             * Format: int64
             * @description Segundos até o access token vencer
             * @example 900
             */
            expires_in: number;
            /** @description Refresh token rotativo (tnz_rt_...), 30 dias de inatividade; cada uso devolve outro */
            refresh_token: string;
            /** @enum {string} */
            scope: "app_campo";
            /** @enum {string} */
            token_type: "Bearer";
        };
        TrocaEmpresa: {
            /**
             * @description Canal da credencial
             * @enum {string}
             */
            canal: "web" | "app";
            /** @description Token CSRF da sessão web, para o cabeçalho X-CSRF-Token em POST, PATCH, PUT e DELETE; nulo no app */
            csrf_token: string | null;
            empresa: components["schemas"]["EmpresaSessao"];
            /** @description Empresas em que a pessoa entra, para o seletor; a atual vem primeiro */
            empresas: components["schemas"]["EmpresaAcesso"][] | null;
            /**
             * @description Escopo de acesso da credencial (seção 4.9)
             * @enum {string}
             */
            escopo_acesso: "interno" | "prestador" | "locatario" | "link";
            /**
             * Format: date-time
             * @description Teto absoluto da sessão web (7 dias); a inatividade de 12 horas vence antes. Nulo fora do login e da troca de empresa
             */
            expira_em: string | null;
            /** @description A empresa exige segundo fator e a pessoa ainda não o ativou: a sessão só serve para ler a sessão, sair, trocar de empresa e ativar o MFA; as demais rotas respondem 401 nao_autenticado com mfa_pendente (RN-ADM-027) */
            mfa_pendente: boolean;
            /** @description Códigos de origem dos perfis ativos */
            perfis: string[] | null;
            /** @description Permissões efetivas, em ordem alfabética (perfis e permissões avulsas) */
            permissoes: string[] | null;
            /** @description Tokens novos do app; nulo na web, que recebe o cookie novo */
            tokens: components["schemas"]["TokensApp"];
            usuario: components["schemas"]["UsuarioSessao"];
        };
        Turno: {
            /**
             * Format: int32
             * @description 1 = segunda ... 7 = domingo
             */
            dia_semana: number;
            /**
             * @description Menor que o início termina no dia seguinte; 24:00 é meia-noite
             * @example 17:00
             */
            hora_fim: string;
            /** @example 08:00 */
            hora_inicio: string;
            /** Format: uuid */
            id: string;
            nome: string | null;
            /** Format: uuid */
            unidade_id: string;
            /** Format: int64 */
            versao: number;
        };
        Unidade: {
            /** Format: date-time */
            atualizado_em: string;
            cnpj: string | null;
            codigo: string;
            /** Format: date-time */
            criado_em: string;
            endereco: components["schemas"]["Endereco"];
            /** @description Nulo: herda o da empresa */
            fuso_horario: string | null;
            /** Format: uuid */
            id: string;
            /** Format: double */
            latitude: number | null;
            /** Format: double */
            longitude: number | null;
            nome: string;
            /** @enum {string} */
            situacao: "ativa" | "inativa";
            turnos: components["schemas"]["Turno"][] | null;
            /** Format: int64 */
            versao: number;
        };
        Usuario: {
            /** Format: date-time */
            atualizado_em: string;
            /** Format: date-time */
            criado_em: string;
            /**
             * @description Custo por hora, texto decimal
             * @example 125.50
             */
            custo_hora: string | null;
            email: string | null;
            /** @enum {string} */
            escopo_acesso: "interno" | "prestador" | "locatario";
            /** Format: uuid */
            fornecedor_id: string;
            /** Format: uuid */
            id: string;
            /** Format: uuid */
            locatario_id: string;
            nome: string;
            perfis: components["schemas"]["PerfilAtribuido"][] | null;
            /** @description Permissões concedidas a este usuário além dos perfis, como os.encerrar */
            permissoes_avulsas: string[] | null;
            /** Format: uuid */
            pessoa_id: string;
            /** @enum {string} */
            situacao: "convidado" | "ativo" | "inativo";
            telefone: string | null;
            /** Format: uuid */
            unidade_padrao_id: string;
            /** Format: int64 */
            versao: number;
        };
        UsuarioCriado: {
            /** Format: date-time */
            atualizado_em: string;
            /** @description Falso: o usuário existe, mas o e-mail não saiu; use reenviar-convite */
            convite_enviado: boolean;
            /** Format: date-time */
            criado_em: string;
            /**
             * @description Custo por hora, texto decimal
             * @example 125.50
             */
            custo_hora: string | null;
            email: string | null;
            /** @enum {string} */
            escopo_acesso: "interno" | "prestador" | "locatario";
            /** Format: uuid */
            fornecedor_id: string;
            /** Format: uuid */
            id: string;
            /** Format: uuid */
            locatario_id: string;
            nome: string;
            perfis: components["schemas"]["PerfilAtribuido"][] | null;
            /** @description Permissões concedidas a este usuário além dos perfis, como os.encerrar */
            permissoes_avulsas: string[] | null;
            /** Format: uuid */
            pessoa_id: string;
            /** @enum {string} */
            situacao: "convidado" | "ativo" | "inativo";
            telefone: string | null;
            /** Format: uuid */
            unidade_padrao_id: string;
            /** Format: int64 */
            versao: number;
        };
        UsuarioSessao: {
            /** Format: email */
            email: string;
            /** Format: uuid */
            id: string;
            /** @description Segundo fator por TOTP ativo */
            mfa_ativo: boolean;
            nome: string;
            /** Format: uuid */
            pessoa_id: string;
        };
        Versao: {
            /**
             * @description Commit curto que gerou o build
             * @example a1b2c3d
             */
            commit: string;
            /**
             * @description Instante do build em UTC (RFC 3339)
             * @example 2026-11-12T14:05:00Z
             */
            data_build: string;
            /**
             * @description Versão do release em execução (SemVer)
             * @example 1.4.2
             */
            versao: string;
        };
    };
    responses: never;
    parameters: never;
    requestBodies: never;
    headers: never;
    pathItems: never;
}
export type $defs = Record<string, never>;
export interface operations {
    iniciarAutorizacaoApp: {
        parameters: {
            query: {
                client_id: string;
                code_challenge: string;
                code_challenge_method: "S256";
                redirect_uri: string;
                response_type: "code";
                state?: string;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description Found */
            302: {
                headers: {
                    "Cache-Control"?: string;
                    Location?: string;
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    autorizarApp: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["EntradaAutorizarPOSTBody"];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    "Cache-Control"?: string;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Autorizacao"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    revogarTokenApp: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/x-www-form-urlencoded": components["schemas"]["EntradaRevogarBody"];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    emitirTokensApp: {
        parameters: {
            query?: never;
            header?: {
                /** @description Dispositivo do app, gravado na sessão */
                "X-Tenaz-Dispositivo-Id"?: string;
            };
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/x-www-form-urlencoded": components["schemas"]["EntradaTokenBody"];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    "Cache-Control"?: string;
                    Pragma?: string;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TokensApp"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    aceitarConvite: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["EntradaConviteAceitarBody"];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ConviteAceito"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    consultarConvite: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["EntradaConviteConsultarBody"];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Convite"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    entrar: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["EntradaLoginBody"];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    "Cache-Control"?: string;
                    "Set-Cookie"?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Sessao"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Gone */
            410: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    sair: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description No Content */
            204: {
                headers: {
                    "Set-Cookie"?: string | null;
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    confirmarMFA: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["EntradaCodigoMFABody"];
            };
        };
        responses: {
            /** @description No Content */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    desativarMFA: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["EntradaCodigoMFABody"];
            };
        };
        responses: {
            /** @description No Content */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Failed */
            412: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Required */
            428: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    iniciarMFA: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    "Cache-Control"?: string;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ConfiguracaoMFA"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    verificarMFA: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["EntradaMFABody"];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    "Cache-Control"?: string;
                    "Set-Cookie"?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Sessao"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Gone */
            410: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    recuperarSenha: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["EntradaRecuperarSenhaBody"];
            };
        };
        responses: {
            /** @description Accepted */
            202: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Gone */
            410: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    redefinirSenha: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["EntradaRedefinirSenhaBody"];
            };
        };
        responses: {
            /** @description No Content */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Gone */
            410: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    obterSessao: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    "Cache-Control"?: string;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Sessao"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    trocarEmpresa: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["EntradaTrocarEmpresaBody"];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    "Cache-Control"?: string;
                    "Set-Cookie"?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["TrocaEmpresa"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    trocarSenha: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["EntradaTrocarSenhaBody"];
            };
        };
        responses: {
            /** @description No Content */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    obterEmpresa: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Empresa"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    alterarEmpresa: {
        parameters: {
            query?: never;
            header?: {
                /** @description ETag lido do recurso; obrigatório */
                "If-Match"?: string;
            };
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["PatchEmpresa"];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Empresa"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Failed */
            412: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Required */
            428: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    criarEmpresa: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["NovaEmpresa"];
            };
        };
        responses: {
            /** @description Created */
            201: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ContaCriada"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    listarEventosAuditoria: {
        parameters: {
            query?: {
                acao?: string;
                cursor?: string;
                entidade?: string;
                entidade_id?: string;
                limite?: number;
                /** @description A partir deste instante, inclusive */
                "ocorrido_em[gte]"?: string;
                /** @description Antes deste instante */
                "ocorrido_em[lt]"?: string;
                usuario_id?: string;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ListaResumoEvento"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    obterEventoAuditoria: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description Identificador (UUID) */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Evento"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    listarPerfis: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ListaPerfil"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    listarPermissoes: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ListaPermissao"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    listarPessoas: {
        parameters: {
            query?: {
                /** @description Trecho do nome */
                busca?: string;
                /** @description true: só pessoas com usuário; false: só sem usuário */
                com_usuario?: "true" | "false";
                /** @description proximo_cursor da página anterior */
                cursor?: string;
                fornecedor_id?: string;
                limite?: number;
                situacao?: "ativa" | "inativa";
                unidade_id?: string;
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ListaResumoPessoa"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    criarPessoa: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["EntradaPessoa"];
            };
        };
        responses: {
            /** @description Created */
            201: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Pessoa"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    obterPessoa: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description Identificador (UUID) */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Pessoa"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    inativarPessoa: {
        parameters: {
            query?: never;
            header?: {
                /** @description ETag lido do recurso; obrigatório */
                "If-Match"?: string;
            };
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description No Content */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Failed */
            412: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Required */
            428: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    alterarPessoa: {
        parameters: {
            query?: never;
            header?: {
                /** @description ETag lido do recurso; obrigatório */
                "If-Match"?: string;
            };
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["PatchPessoa"];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Pessoa"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Failed */
            412: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Required */
            428: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    reativarPessoa: {
        parameters: {
            query?: never;
            header?: {
                /** @description ETag lido do recurso; obrigatório */
                "If-Match"?: string;
            };
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Pessoa"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Failed */
            412: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Required */
            428: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    listarUnidades: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ListaResumoUnidade"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    criarUnidade: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["EntradaUnidade"];
            };
        };
        responses: {
            /** @description Created */
            201: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Unidade"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    obterUnidade: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description Identificador (UUID) */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Unidade"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    inativarUnidade: {
        parameters: {
            query?: never;
            header?: {
                /** @description ETag lido do recurso; obrigatório */
                "If-Match"?: string;
            };
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description No Content */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Failed */
            412: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Required */
            428: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    alterarUnidade: {
        parameters: {
            query?: never;
            header?: {
                /** @description ETag lido do recurso; obrigatório */
                "If-Match"?: string;
            };
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["PatchUnidade"];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Unidade"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Failed */
            412: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Required */
            428: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    listarFeriados: {
        parameters: {
            query?: {
                /** @description A partir desta data, inclusive */
                "data[gte]"?: string;
                /** @description Até esta data, inclusive */
                "data[lte]"?: string;
            };
            header?: never;
            path: {
                /** @description Unidade */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ListaFeriado"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    incluirFeriado: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["EntradaFeriado"];
            };
        };
        responses: {
            /** @description Created */
            201: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Feriado"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    removerFeriado: {
        parameters: {
            query?: never;
            header?: {
                /** @description ETag do feriado; obrigatório */
                "If-Match"?: string;
            };
            path: {
                feriado_id: string;
                /** @description Unidade */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description No Content */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Failed */
            412: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Required */
            428: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    reativarUnidade: {
        parameters: {
            query?: never;
            header?: {
                /** @description ETag lido do recurso; obrigatório */
                "If-Match"?: string;
            };
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Unidade"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Failed */
            412: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Required */
            428: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    incluirTurno: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["EntradaTurno"];
            };
        };
        responses: {
            /** @description Created */
            201: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Turno"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    removerTurno: {
        parameters: {
            query?: never;
            header?: {
                /** @description ETag do turno; obrigatório */
                "If-Match"?: string;
            };
            path: {
                /** @description Unidade */
                id: string;
                turno_id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description No Content */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Failed */
            412: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Required */
            428: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    alterarTurno: {
        parameters: {
            query?: never;
            header?: {
                /** @description ETag do turno; obrigatório */
                "If-Match"?: string;
            };
            path: {
                /** @description Unidade */
                id: string;
                turno_id: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["PatchTurno"];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Turno"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Failed */
            412: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Required */
            428: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    listarUsuarios: {
        parameters: {
            query?: {
                /** @description Trecho do nome ou do e-mail */
                busca?: string;
                cursor?: string;
                limite?: number;
                /** @description Código do perfil de origem */
                perfil?: string;
                situacao?: "convidado" | "ativo" | "inativo";
            };
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["ListaResumoUsuario"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    criarUsuario: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["EntradaUsuario"];
            };
        };
        responses: {
            /** @description Created */
            201: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["UsuarioCriado"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    obterUsuario: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description Identificador (UUID) */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Usuario"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    alterarUsuario: {
        parameters: {
            query?: never;
            header?: {
                /** @description ETag lido do recurso; obrigatório */
                "If-Match"?: string;
            };
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["PatchUsuario"];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Usuario"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Failed */
            412: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Required */
            428: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    desativarUsuario: {
        parameters: {
            query?: never;
            header?: {
                /** @description ETag lido do recurso; obrigatório */
                "If-Match"?: string;
            };
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Usuario"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Failed */
            412: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Required */
            428: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    encerrarSessoesUsuario: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description Identificador (UUID) */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["SessoesEncerradas"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    concederPermissaoUsuario: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["EntradaPermissao"];
            };
        };
        responses: {
            /** @description Created */
            201: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Usuario"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    retirarPermissaoUsuario: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                id: string;
                permissao: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description No Content */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Failed */
            412: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Required */
            428: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    reativarUsuario: {
        parameters: {
            query?: never;
            header?: {
                /** @description ETag lido do recurso; obrigatório */
                "If-Match"?: string;
            };
            path: {
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["UsuarioCriado"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Failed */
            412: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Required */
            428: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    reenviarConvite: {
        parameters: {
            query?: never;
            header?: never;
            path: {
                /** @description Identificador (UUID) */
                id: string;
            };
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description No Content */
            204: {
                headers: {
                    [name: string]: unknown;
                };
                content?: never;
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    obterUsuarioEu: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Usuario"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    alterarUsuarioEu: {
        parameters: {
            query?: never;
            header?: {
                /** @description ETag lido do recurso; obrigatório */
                "If-Match"?: string;
            };
            path?: never;
            cookie?: never;
        };
        requestBody: {
            content: {
                "application/json": components["schemas"]["PatchEu"];
            };
        };
        responses: {
            /** @description OK */
            200: {
                headers: {
                    ETag?: string | null;
                    Location?: string | null;
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Usuario"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unauthorized */
            401: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Forbidden */
            403: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Not Found */
            404: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Conflict */
            409: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Failed */
            412: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Precondition Required */
            428: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
    obterVersao: {
        parameters: {
            query?: never;
            header?: never;
            path?: never;
            cookie?: never;
        };
        requestBody?: never;
        responses: {
            /** @description OK */
            200: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/json": components["schemas"]["Versao"];
                };
            };
            /** @description Bad Request */
            400: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Unprocessable Entity */
            422: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
            /** @description Internal Server Error */
            500: {
                headers: {
                    [name: string]: unknown;
                };
                content: {
                    "application/problem+json": components["schemas"]["Problema"];
                };
            };
        };
    };
}
