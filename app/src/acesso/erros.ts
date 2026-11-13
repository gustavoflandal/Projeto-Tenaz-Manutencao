import type { ChaveDeTexto } from '../i18n';

// Erros do acesso. O app trata o `codigo` estável do problem+json (seção 4 de
// docs/3_ESPECIFICACAO_APIS.md), nunca o texto.

/** Resposta de erro da API, já lida. */
export class ErroApi extends Error {
  constructor(
    readonly status: number,
    readonly codigo: string,
    /** Membro `error` da RFC 6749 em /oauth/token (invalid_grant). */
    readonly error: string | null,
    /** 401 nao_autenticado com mfa_pendente: a sessão só serve para ativar o MFA. */
    readonly mfaPendente: boolean,
    detalhe: string,
    readonly idRastreio: string | null = null,
  ) {
    super(detalhe);
    this.name = 'ErroApi';
  }
}

/** Sem conexão ou servidor inalcançável. Nunca derruba a sessão: o campo vive sem sinal. */
export class ErroDeRede extends Error {
  constructor(causa?: unknown) {
    super('Sem conexão com a Tenaz.');
    this.name = 'ErroDeRede';
    this.cause = causa;
  }
}

/** Refresh token recusado (vencido, revogado ou reusado): é preciso entrar de novo. A fila local continua intacta. */
export class ErroSessaoEncerrada extends Error {
  constructor() {
    super('Sua sessão terminou. Entre de novo; o que você registrou neste aparelho continua guardado.');
    this.name = 'ErroSessaoEncerrada';
  }
}

/**
 * A credencial do aparelho passou a ser de outra empresa ou usuário enquanto o
 * pedido esperava. O pedido NÃO é reenviado: enviaria dado de uma conta com a
 * credencial de outra. O que ficou na fila da conta original segue guardado.
 */
export class ErroContaMudou extends Error {
  constructor() {
    super('A conta ativa mudou durante o envio. Nada foi enviado.');
    this.name = 'ErroContaMudou';
  }
}

/** O servidor devolveu error=access_denied: perfil sem mob.usar_app (o locatário). */
export class ErroAcessoNegado extends Error {
  constructor() {
    super('Este perfil não usa o app de campo. Use a versão web da Tenaz.');
    this.name = 'ErroAcessoNegado';
  }
}

/** Retorno do navegador que não é desta tentativa de login (state diferente) ou sem código. */
export class ErroRetornoInvalido extends Error {
  constructor(motivo: string) {
    super(motivo);
    this.name = 'ErroRetornoInvalido';
  }
}

/**
 * Chave do texto, no catálogo @tenaz/i18n, que explica o erro ao usuário. O
 * texto em si sai no idioma escolhido, na tela: o erro nunca carrega frase
 * pronta nem repete o `detail` do servidor (que vem em um idioma só).
 */
export function chaveDoErro(e: unknown): ChaveDeTexto {
  if (e instanceof ErroDeRede) return 'mobile.erros.rede';
  if (e instanceof ErroAcessoNegado) return 'mobile.erros.acessoNegado';
  if (e instanceof ErroSessaoEncerrada) return 'mobile.erros.sessaoEncerrada';
  if (e instanceof ErroRetornoInvalido) return 'mobile.erros.retornoInvalido';
  if (e instanceof ErroContaMudou) return 'mobile.erros.contaMudou';
  if (e instanceof ErroApi) {
    if (e.mfaPendente) return 'mobile.erros.mfaPendente';
    if (e.status === 404) return 'mobile.erros.semAcesso';
    if (e.status === 429) return 'mobile.erros.muitasTentativas';
    if (e.status >= 500) return 'mobile.erros.indisponivel';
  }
  return 'mobile.erros.generico';
}
