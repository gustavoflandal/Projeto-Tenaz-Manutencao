// Telas de autenticação (RF-ADM-006): login com segundo fator e escolha de
// empresa, sessão restrita, guarda de rotas, recuperação de senha, convite e
// autorização do app. A API é simulada com as respostas do contrato.

import { screen, waitFor, within } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';

import { problema, simularAPI } from '../../../testes/api-simulada';
import { sessao, usuario as usuarioDeTeste } from '../../../testes/dados';
import { renderizarRota, rotaDaSessao } from '../../../testes/renderizar';

const DESAFIO = 'd'.repeat(43);
const TOKEN = 't'.repeat(43);

async function preencherLogin(usuario: ReturnType<typeof userEvent.setup>, senha = 'senha-correta-123') {
  await usuario.type(await screen.findByLabelText('E-mail'), 'ana@polar.com.br');
  await usuario.type(screen.getByLabelText('Senha'), senha);
  await usuario.click(screen.getByRole('button', { name: 'Entrar' }));
}

describe('login', () => {
  it('entra, guarda a sessão no cache e vai para o painel; o POST leva Idempotency-Key', async () => {
    const api = simularAPI({ 'POST /v1/auth/login': { corpo: sessao() }, 'GET /v1/usuarios/eu': { corpo: usuarioDeTeste({ id: 'u-1', versao: 1 }) } });
    const usuario = userEvent.setup();
    renderizarRota('/entrar');

    await preencherLogin(usuario);

    expect(await screen.findByRole('heading', { level: 1, name: 'Painel' })).toBeInTheDocument();
    const [chamada] = api.chamadasA('POST /v1/auth/login');
    expect(chamada?.corpo).toEqual({ email: 'ana@polar.com.br', senha: 'senha-correta-123' });
    expect(chamada?.cabecalhos.get('Idempotency-Key')).toBeTruthy();
    // Dependeu só do login: não houve nova leitura da sessão.
    expect(api.chamadasA('GET /v1/auth/sessao')).toHaveLength(0);
    expect(api.semRota).toEqual([]);
  });

  it('mostra a mesma mensagem para qualquer falha de credencial, sem sair da tela', async () => {
    simularAPI({ 'POST /v1/auth/login': problema(401, 'credencial_invalida') });
    const usuario = userEvent.setup();
    renderizarRota('/entrar');

    await preencherLogin(usuario, 'errada');

    expect(await screen.findByRole('alert')).toHaveTextContent('E-mail ou senha incorretos');
    expect(screen.getByRole('heading', { level: 1, name: 'Entrar' })).toBeInTheDocument();
  });

  it('valida e-mail e senha antes de enviar, com erro no campo e resumo com link', async () => {
    const api = simularAPI({});
    const usuario = userEvent.setup();
    renderizarRota('/entrar');

    await usuario.click(await screen.findByRole('button', { name: 'Entrar' }));

    const campo = await screen.findByLabelText('E-mail');
    await waitFor(() => expect(campo).toHaveAttribute('aria-invalid', 'true'));
    expect(campo).toHaveAccessibleDescription('Informe o e-mail.');
    expect(screen.getByRole('link', { name: 'Informe o e-mail.' })).toHaveAttribute('href', `#${campo.id}`);
    expect(api.chamadas).toHaveLength(0);
  });

  it('pede o código do segundo fator e termina o login em mfa/verificar', async () => {
    const api = simularAPI({
      'POST /v1/auth/login': problema(401, 'mfa_obrigatorio', { desafio_id: DESAFIO, desafio_expira_em: '2026-11-13T12:05:00Z' }),
      'POST /v1/auth/mfa/verificar': { corpo: sessao() },
    });
    const usuario = userEvent.setup();
    renderizarRota('/entrar');

    await preencherLogin(usuario);
    expect(await screen.findByRole('heading', { level: 1, name: 'Confirme o segundo fator' })).toBeInTheDocument();
    await usuario.type(screen.getByLabelText('Código'), '123456');
    await usuario.click(screen.getByRole('button', { name: 'Verificar código' }));

    expect(await screen.findByRole('heading', { level: 1, name: 'Painel' })).toBeInTheDocument();
    expect(api.chamadasA('POST /v1/auth/mfa/verificar')[0]?.corpo).toEqual({ desafio_id: DESAFIO, codigo: '123456' });
  });

  it('informa as tentativas restantes quando o código está errado', async () => {
    simularAPI({
      'POST /v1/auth/login': problema(401, 'mfa_obrigatorio', { desafio_id: DESAFIO }),
      'POST /v1/auth/mfa/verificar': problema(401, 'mfa_codigo_invalido', { tentativas_restantes: 3 }),
    });
    const usuario = userEvent.setup();
    renderizarRota('/entrar');

    await preencherLogin(usuario);
    await usuario.type(await screen.findByLabelText('Código'), '000000');
    await usuario.click(screen.getByRole('button', { name: 'Verificar código' }));

    expect(await screen.findByRole('alert')).toHaveTextContent('Restam 3 tentativas');
    expect(screen.getByLabelText('Código')).toHaveValue('');
  });

  it('volta para a senha quando o desafio acaba (sem tentativas)', async () => {
    simularAPI({
      'POST /v1/auth/login': problema(401, 'mfa_obrigatorio', { desafio_id: DESAFIO }),
      'POST /v1/auth/mfa/verificar': problema(401, 'mfa_codigo_invalido', { tentativas_restantes: 0 }),
    });
    const usuario = userEvent.setup();
    renderizarRota('/entrar');

    await preencherLogin(usuario);
    await usuario.type(await screen.findByLabelText('Código'), '000000');
    await usuario.click(screen.getByRole('button', { name: 'Verificar código' }));

    expect(await screen.findByRole('heading', { level: 1, name: 'Entrar' })).toBeInTheDocument();
    expect(screen.getByRole('status')).toHaveTextContent('O código venceu ou as tentativas acabaram');
  });

  it('com mais de uma empresa, pede a escolha e troca para a escolhida', async () => {
    const api = simularAPI({
      'POST /v1/auth/login': {
        corpo: sessao({
          empresas: [
            { id: 'e-1', nome: 'Refrigeração Polar', atual: true },
            { id: 'e-2', nome: 'Norte Transportes', atual: false },
          ],
        }),
      },
      'POST /v1/auth/trocar-empresa': {
        corpo: { ...sessao({ empresa: { nome: 'Norte Transportes', plano: 'profissional', srv_habilitado: false, pais: 'BR', moeda: 'BRL', idioma_padrao: 'pt-BR' } }), tokens: null },
      },
    });
    const usuario = userEvent.setup();
    renderizarRota('/entrar');

    await preencherLogin(usuario);
    expect(await screen.findByRole('heading', { level: 1, name: 'Escolha a empresa' })).toBeInTheDocument();
    expect(screen.getByRole('radio', { name: /Refrigeração Polar/ })).toBeChecked();
    await usuario.click(screen.getByRole('radio', { name: 'Norte Transportes' }));
    await usuario.click(screen.getByRole('button', { name: 'Continuar' }));

    expect(await screen.findByText('Norte Transportes')).toBeInTheDocument();
    expect(api.chamadasA('POST /v1/auth/trocar-empresa')[0]?.corpo).toEqual({ empresa_escolhida_id: 'e-2' });
  });

  it('com a sessão restrita (mfa_pendente), leva à ativação do segundo fator', async () => {
    simularAPI({ 'POST /v1/auth/login': { corpo: sessao({ mfa_pendente: true }) } });
    const usuario = userEvent.setup();
    renderizarRota('/entrar');

    await preencherLogin(usuario);

    expect(await screen.findByRole('heading', { level: 1, name: 'Ative o segundo fator' })).toBeInTheDocument();
  });
});

describe('guarda de rotas', () => {
  it('sem sessão, manda para o login guardando o endereço pedido', async () => {
    simularAPI(rotaDaSessao(null));
    const { roteador } = renderizarRota('/administracao/pessoas?busca=ana');

    expect(await screen.findByRole('heading', { level: 1, name: 'Entrar' })).toBeInTheDocument();
    expect(roteador.state.location.pathname).toBe('/entrar');
    expect(new URLSearchParams(roteador.state.location.search).get('retorno')).toBe('/administracao/pessoas?busca=ana');
  });

  it('com a sessão restrita, só a tela de ativação abre', async () => {
    simularAPI(rotaDaSessao(sessao({ mfa_pendente: true })));
    const { roteador } = renderizarRota('/administracao/pessoas');

    expect(await screen.findByRole('heading', { level: 1, name: 'Ative o segundo fator' })).toBeInTheDocument();
    expect(roteador.state.location.pathname).toBe('/ativar-mfa');
  });

  it('um 401 numa tela leva ao login com o aviso de sessão encerrada', async () => {
    simularAPI({ ...rotaDaSessao(sessao()), 'GET /v1/unidades': problema(401, 'nao_autenticado', { instance: '/v1/unidades' }) });
    renderizarRota('/administracao/unidades');

    expect(await screen.findByRole('heading', { level: 1, name: 'Entrar' })).toBeInTheDocument();
    expect(screen.getByRole('status')).toHaveTextContent('Sua sessão terminou por segurança');
  });

  it('um 403 mostra o estado de falta de permissão, sem tentar de novo', async () => {
    simularAPI({ ...rotaDaSessao(sessao()), 'GET /v1/unidades': problema(403, 'permissao_negada') });
    renderizarRota('/administracao/unidades');

    expect(await screen.findByRole('heading', { name: 'Sem acesso a esta área' })).toBeInTheDocument();
    expect(screen.queryByRole('button', { name: 'Tentar de novo' })).not.toBeInTheDocument();
  });

  it('sem a permissão da rota, nem chama a API', async () => {
    const api = simularAPI(rotaDaSessao(sessao({ permissoes: ['adm.ler'] })));
    renderizarRota('/administracao/usuarios');

    expect(await screen.findByRole('heading', { name: 'Sem acesso a esta área' })).toBeInTheDocument();
    expect(api.chamadasA('GET /v1/usuarios')).toHaveLength(0);
  });
});

describe('menu do usuário', () => {
  it('sai: POST com X-CSRF-Token, esvazia a sessão e volta ao login', async () => {
    const api = simularAPI({ ...rotaDaSessao(sessao()), 'POST /v1/auth/logout': { status: 204 } });
    const usuario = userEvent.setup();
    renderizarRota('/');

    await usuario.click(await screen.findByRole('button', { name: 'Menu de Ana Souza' }));
    await usuario.click(await screen.findByRole('menuitem', { name: 'Sair' }));

    expect(await screen.findByRole('heading', { level: 1, name: 'Entrar' })).toBeInTheDocument();
    await waitFor(() => expect(api.chamadasA('POST /v1/auth/logout')).toHaveLength(1));
    expect(api.chamadasA('POST /v1/auth/logout')[0]?.cabecalhos.get('X-CSRF-Token')).toBe('csrf-de-teste');
  });

  it('troca de empresa pelo menu', async () => {
    const duas = sessao({
      empresas: [
        { id: 'e-1', nome: 'Refrigeração Polar', atual: true },
        { id: 'e-2', nome: 'Norte Transportes', atual: false },
      ],
    });
    const api = simularAPI({
      ...rotaDaSessao(duas),
      'POST /v1/auth/trocar-empresa': {
        corpo: { ...sessao({ empresa: { nome: 'Norte Transportes', plano: null, srv_habilitado: true, pais: 'BR', moeda: 'BRL', idioma_padrao: 'pt-BR' } }), tokens: null },
      },
    });
    const usuario = userEvent.setup();
    renderizarRota('/');

    await usuario.click(await screen.findByRole('button', { name: 'Menu de Ana Souza' }));
    await usuario.click(await screen.findByRole('menuitem', { name: 'Norte Transportes' }));

    await waitFor(() => expect(api.chamadasA('POST /v1/auth/trocar-empresa')).toHaveLength(1));
    expect(api.chamadasA('POST /v1/auth/trocar-empresa')[0]?.corpo).toEqual({ empresa_escolhida_id: 'e-2' });
  });
});

describe('recuperação de senha', () => {
  it('responde igual, exista a conta ou não', async () => {
    const api = simularAPI({ 'POST /v1/auth/recuperar-senha': { status: 202 } });
    const usuario = userEvent.setup();
    renderizarRota('/esqueci-a-senha');

    await usuario.type(await screen.findByLabelText('E-mail'), 'quem@polar.com.br');
    await usuario.click(screen.getByRole('button', { name: 'Enviar link' }));

    expect(await screen.findByRole('status')).toHaveTextContent('Se quem@polar.com.br tem conta');
    expect(api.chamadasA('POST /v1/auth/recuperar-senha')[0]?.corpo).toEqual({ email: 'quem@polar.com.br' });
  });

  it('redefine a senha com o token do fragmento, enviado no corpo e tirado da barra', async () => {
    const api = simularAPI({ 'POST /v1/auth/redefinir-senha': { status: 204 } });
    const usuario = userEvent.setup();
    const { roteador } = renderizarRota(`/redefinir-senha#token=${TOKEN}`);

    await usuario.type(await screen.findByLabelText('Senha nova'), 'uma-senha-bem-longa');
    await usuario.type(screen.getByLabelText('Repita a senha nova'), 'uma-senha-bem-longa');
    await usuario.click(screen.getByRole('button', { name: 'Salvar senha nova' }));

    expect(await screen.findByRole('status')).toHaveTextContent('Senha alterada');
    const [chamada] = api.chamadasA('POST /v1/auth/redefinir-senha');
    expect(chamada?.corpo).toEqual({ token: TOKEN, nova_senha: 'uma-senha-bem-longa' });
    expect(chamada?.caminho).not.toContain(TOKEN);
    expect(roteador.state.location.hash).toBe('');
  });

  it('recusa senhas diferentes e curtas antes de enviar', async () => {
    const api = simularAPI({});
    const usuario = userEvent.setup();
    renderizarRota(`/redefinir-senha#token=${TOKEN}`);

    await usuario.type(await screen.findByLabelText('Senha nova'), 'curta');
    await usuario.type(screen.getByLabelText('Repita a senha nova'), 'outra');
    await usuario.click(screen.getByRole('button', { name: 'Salvar senha nova' }));

    expect(await screen.findAllByText('Use ao menos 10 caracteres.')).not.toHaveLength(0);
    expect(screen.getByLabelText('Repita a senha nova')).toHaveAccessibleDescription(/senhas não são iguais/);
    expect(api.chamadas).toHaveLength(0);
  });

  it('mostra o erro de senha fraca no campo', async () => {
    simularAPI({
      'POST /v1/auth/redefinir-senha': problema(422, 'senha_fraca', {
        erros: [{ campo: 'senha', codigo: 'valor_invalido', mensagem: 'Esta senha aparece em listas de senhas vazadas. Escolha outra.' }],
      }),
    });
    const usuario = userEvent.setup();
    renderizarRota(`/redefinir-senha#token=${TOKEN}`);

    await usuario.type(await screen.findByLabelText('Senha nova'), 'password1234');
    await usuario.type(screen.getByLabelText('Repita a senha nova'), 'password1234');
    await usuario.click(screen.getByRole('button', { name: 'Salvar senha nova' }));

    expect(await screen.findByLabelText('Senha nova')).toHaveAccessibleDescription(/listas de senhas vazadas/);
  });

  it('link vencido (410) oferece pedir outro', async () => {
    simularAPI({ 'POST /v1/auth/redefinir-senha': problema(410, 'link_email_invalido') });
    const usuario = userEvent.setup();
    renderizarRota(`/redefinir-senha#token=${TOKEN}`);

    await usuario.type(await screen.findByLabelText('Senha nova'), 'uma-senha-bem-longa');
    await usuario.type(screen.getByLabelText('Repita a senha nova'), 'uma-senha-bem-longa');
    await usuario.click(screen.getByRole('button', { name: 'Salvar senha nova' }));

    expect(await screen.findByRole('alert')).toHaveTextContent('Este link venceu');
    expect(screen.getByRole('link', { name: 'Pedir novo link' })).toBeInTheDocument();
  });

  it('sem token no fragmento, trata como link inválido sem chamar a API', async () => {
    const api = simularAPI({});
    renderizarRota('/redefinir-senha');
    expect(await screen.findByRole('alert')).toHaveTextContent('Este link venceu');
    expect(api.chamadas).toHaveLength(0);
  });
});

describe('convite', () => {
  it('consulta o convite, pede a senha nova e aceita', async () => {
    const api = simularAPI({
      'POST /v1/auth/convites/consultar': { corpo: { email: 'novo@polar.com.br', pede_senha: true } },
      'POST /v1/auth/convites/aceitar': { corpo: { identidade_nova: true } },
    });
    const usuario = userEvent.setup();
    renderizarRota(`/convite#token=${TOKEN}`);

    expect(await screen.findByText(/novo@polar.com.br/)).toBeInTheDocument();
    await usuario.type(screen.getByLabelText('Crie uma senha'), 'uma-senha-bem-longa');
    await usuario.type(screen.getByLabelText('Repita a senha nova'), 'uma-senha-bem-longa');
    await usuario.click(screen.getByRole('button', { name: 'Aceitar convite' }));

    expect(await screen.findByRole('status')).toHaveTextContent('Convite aceito');
    expect(api.chamadasA('POST /v1/auth/convites/consultar')[0]?.corpo).toEqual({ token: TOKEN });
    expect(api.chamadasA('POST /v1/auth/convites/aceitar')[0]?.corpo).toEqual({ token: TOKEN, senha: 'uma-senha-bem-longa' });
  });

  it('e-mail que já tem login aceita sem pedir senha', async () => {
    const api = simularAPI({
      'POST /v1/auth/convites/consultar': { corpo: { email: 'ana@polar.com.br', pede_senha: false } },
      'POST /v1/auth/convites/aceitar': { corpo: { identidade_nova: false } },
    });
    const usuario = userEvent.setup();
    renderizarRota(`/convite#token=${TOKEN}`);

    await usuario.click(await screen.findByRole('button', { name: 'Aceitar convite' }));

    expect(await screen.findByRole('status')).toHaveTextContent('Convite aceito');
    expect(api.chamadasA('POST /v1/auth/convites/aceitar')[0]?.corpo).toEqual({ token: TOKEN });
  });

  it('convite vencido, usado ou substituído (410) mostra o aviso', async () => {
    simularAPI({ 'POST /v1/auth/convites/consultar': problema(410, 'link_email_invalido') });
    renderizarRota(`/convite#token=${TOKEN}`);
    expect(await screen.findByRole('alert')).toHaveTextContent('Este convite venceu ou já foi usado');
  });
});

describe('ativação do segundo fator pela sessão restrita', () => {
  it('mostra o QR code e o segredo, e confirma com o código', async () => {
    const api = simularAPI({
      ...rotaDaSessao(sessao({ mfa_pendente: true })),
      'POST /v1/auth/mfa/iniciar': { corpo: { segredo_base32: 'JBSWY3DPEHPK3PXP', uri_otpauth: 'otpauth://totp/Tenaz:ana?secret=JBSWY3DPEHPK3PXP' } },
      'POST /v1/auth/mfa/confirmar': { status: 204 },
    });
    const usuario = userEvent.setup();
    renderizarRota('/ativar-mfa');

    await usuario.click(await screen.findByRole('button', { name: 'Ativar segundo fator' }));
    expect(await screen.findByRole('img', { name: 'QR code para o aplicativo autenticador' })).toBeInTheDocument();
    expect(screen.getByText('JBSWY3DPEHPK3PXP')).toBeInTheDocument();
    await usuario.type(screen.getByLabelText('Código de 6 dígitos'), '654321');
    await usuario.click(screen.getByRole('button', { name: 'Confirmar e ativar' }));

    await waitFor(() => expect(api.chamadasA('POST /v1/auth/mfa/confirmar')).toHaveLength(1));
    expect(api.chamadasA('POST /v1/auth/mfa/confirmar')[0]?.corpo).toEqual({ codigo: '654321' });
  });

  it('código errado aparece no campo', async () => {
    simularAPI({
      ...rotaDaSessao(sessao({ mfa_pendente: true })),
      'POST /v1/auth/mfa/iniciar': { corpo: { segredo_base32: 'JBSWY3DPEHPK3PXP', uri_otpauth: 'otpauth://totp/Tenaz:ana?secret=JBSWY3DPEHPK3PXP' } },
      'POST /v1/auth/mfa/confirmar': problema(401, 'mfa_codigo_invalido', { instance: '/v1/auth/mfa/confirmar' }),
    });
    const usuario = userEvent.setup();
    renderizarRota('/ativar-mfa');

    await usuario.click(await screen.findByRole('button', { name: 'Ativar segundo fator' }));
    await usuario.type(await screen.findByLabelText('Código de 6 dígitos'), '111111');
    await usuario.click(screen.getByRole('button', { name: 'Confirmar e ativar' }));

    expect(await screen.findByLabelText('Código de 6 dígitos')).toHaveAccessibleDescription(/Código incorreto/);
    // Um 401 de credencial digitada não encerra a sessão.
    expect(screen.getByRole('heading', { level: 1, name: 'Ative o segundo fator' })).toBeInTheDocument();
  });
});

describe('autorização do app (/autorizar-app)', () => {
  const PEDIDO =
    '/autorizar-app?response_type=code&client_id=tenaz-app&redirect_uri=br.com.tenaz.app%3A%2Foauth&code_challenge=abc&code_challenge_method=S256&state=xyz';

  it('com uma empresa, autoriza e navega até a URI de retorno', async () => {
    const api = simularAPI({
      ...rotaDaSessao(sessao()),
      'POST /oauth/authorize': { corpo: { uri_redirecionamento: 'br.com.tenaz.app:/oauth?code=c0d&state=xyz' } },
    });
    const assign = vi.fn();
    vi.spyOn(globalThis, 'location', 'get').mockReturnValue({ ...globalThis.location, origin: 'http://localhost:3000', assign } as unknown as Location);
    const usuario = userEvent.setup();
    renderizarRota(PEDIDO);

    await usuario.click(await screen.findByRole('button', { name: 'Autorizar o app' }));

    await waitFor(() => expect(assign).toHaveBeenCalledWith('br.com.tenaz.app:/oauth?code=c0d&state=xyz'));
    expect(api.chamadasA('POST /oauth/authorize')[0]?.corpo).toEqual({
      response_type: 'code',
      client_id: 'tenaz-app',
      redirect_uri: 'br.com.tenaz.app:/oauth',
      code_challenge: 'abc',
      code_challenge_method: 'S256',
      empresa_escolhida_id: 'e-1',
      state: 'xyz',
    });
  });

  it('descarta uma resposta que aponta para outro destino', async () => {
    simularAPI({
      ...rotaDaSessao(sessao()),
      'POST /oauth/authorize': { corpo: { uri_redirecionamento: 'https://outro.site/?code=c0d' } },
    });
    const assign = vi.fn();
    vi.spyOn(globalThis, 'location', 'get').mockReturnValue({ ...globalThis.location, origin: 'http://localhost:3000', assign } as unknown as Location);
    const usuario = userEvent.setup();
    renderizarRota(PEDIDO);

    await usuario.click(await screen.findByRole('button', { name: 'Autorizar o app' }));

    expect(await screen.findByRole('alert')).toHaveTextContent('destino inesperado');
    expect(assign).not.toHaveBeenCalled();
  });

  it('com mais de uma empresa, a pessoa escolhe qual autorizar', async () => {
    const api = simularAPI({
      ...rotaDaSessao(
        sessao({
          empresas: [
            { id: 'e-1', nome: 'Refrigeração Polar', atual: true },
            { id: 'e-2', nome: 'Norte Transportes', atual: false },
          ],
        }),
      ),
      'POST /oauth/authorize': { problema: undefined, corpo: { uri_redirecionamento: 'br.com.tenaz.app:/oauth?error=access_denied&state=xyz' } },
    });
    const assign = vi.fn();
    vi.spyOn(globalThis, 'location', 'get').mockReturnValue({ ...globalThis.location, origin: 'http://localhost:3000', assign } as unknown as Location);
    const usuario = userEvent.setup();
    renderizarRota(PEDIDO);

    await usuario.click(await screen.findByRole('radio', { name: 'Norte Transportes' }));
    await usuario.click(screen.getByRole('button', { name: 'Autorizar o app' }));

    await waitFor(() => expect(assign).toHaveBeenCalled());
    expect(api.chamadasA('POST /oauth/authorize')[0]?.corpo).toMatchObject({ empresa_escolhida_id: 'e-2' });
  });

  it('empresa fora do acesso (404) mostra o erro', async () => {
    simularAPI({ ...rotaDaSessao(sessao()), 'POST /oauth/authorize': problema(404, 'recurso_nao_encontrado') });
    const usuario = userEvent.setup();
    renderizarRota(PEDIDO);

    await usuario.click(await screen.findByRole('button', { name: 'Autorizar o app' }));

    expect(await screen.findByRole('alert')).toHaveTextContent('Não encontramos este registro');
  });

  it('pedido incompleto não chama a API', async () => {
    const api = simularAPI(rotaDaSessao(sessao()));
    renderizarRota('/autorizar-app?client_id=tenaz-app');
    expect(await screen.findByRole('alert')).toHaveTextContent('não veio do app');
    expect(within(document.body).queryByRole('button', { name: 'Autorizar o app' })).not.toBeInTheDocument();
    expect(api.chamadasA('POST /oauth/authorize')).toHaveLength(0);
  });

  it('sem sessão, passa pelo login e volta ao pedido do app', async () => {
    simularAPI({ ...rotaDaSessao(null), 'POST /v1/auth/login': { corpo: sessao() } });
    const usuario = userEvent.setup();
    const { roteador } = renderizarRota(PEDIDO);

    await preencherLogin(usuario);

    expect(await screen.findByRole('button', { name: 'Autorizar o app' })).toBeInTheDocument();
    expect(roteador.state.location.pathname).toBe('/autorizar-app');
    expect(roteador.state.location.search).toContain('code_challenge=abc');
  });
});
