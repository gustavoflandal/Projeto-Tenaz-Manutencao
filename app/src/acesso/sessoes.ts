// Sessão do app: login OAuth 2.1 com PKCE, tokens no cofre, renovação
// serializada, troca de empresa e saída (seção 4.3 de
// docs/3_ESPECIFICACAO_APIS.md e RF-MOB-001).
//
// Regras que este arquivo garante e que os testes trancam:
//
//  1. UMA renovação por vez. O refresh token é rotativo: duas renovações com o
//     mesmo token contam como REUSO no servidor e derrubam a família inteira.
//     Quem chega durante uma renovação espera o resultado dela.
//  2. O refresh token novo é gravado no cofre ANTES de a renovação responder a
//     quem esperava.
//  3. Sem rede, a sessão não cai. Só `invalid_grant` encerra a sessão.
//  4. Encerrar sessão (ou sair) NUNCA toca no banco local nem na fila de
//     comandos: só nas credenciais. O banco é um arquivo por empresa e usuário.
//  5. Um pedido feito para a conta A nunca sai com a credencial da conta B,
//     nem depois de uma troca de empresa no meio do caminho.

import type { components } from '@tenaz/api-tipos';

import { arquivoDaConta } from '../sync/deposito-expo';
import type { Cofre } from './cofre';
import {
  ErroAcessoNegado,
  ErroApi,
  ErroContaMudou,
  ErroDeRede,
  ErroRetornoInvalido,
  ErroSessaoEncerrada,
} from './erros';
import { type Cripto, gerarPKCE, uuidV4 } from './pkce';
import { CLIENT_ID, interpretarRetorno, URI_DE_RETORNO, urlDeAutorizacao } from './retorno';

type SessaoDaApi = components['schemas']['Sessao'];
type TokensDaApi = components['schemas']['TokensApp'];
type TrocaDaApi = components['schemas']['TrocaEmpresa'];

const CHAVES = {
  credencial: 'tenaz.credencial',
  sessao: 'tenaz.sessao',
  ultimaConta: 'tenaz.ultima-conta',
  dispositivo: 'tenaz.dispositivo',
  revogacoes: 'tenaz.revogacoes-pendentes',
} as const;

/** Renova com folga: um access token que vence em menos de 60 s já conta como vencido. */
const FOLGA_DO_TOKEN_MS = 60_000;
const TEMPO_LIMITE_MS = 20_000;

export interface Credencial {
  accessToken: string;
  /** Instante (ms desde 1970) em que o access token vence. */
  expiraEm: number;
  refreshToken: string;
  empresaId: string;
  usuarioId: string;
}

/** A conta: uma empresa e um usuário. Define o banco local e a credencial. */
export interface ChaveDeConta {
  empresaId: string;
  usuarioId: string;
}

export interface EmpresaDoSeletor {
  id: string;
  nome: string;
}

/** O que o app guarda da sessão para abrir sem rede. Sem permissões: elas vivem no banco local da conta. */
export interface ContaAtiva extends ChaveDeConta {
  empresaNome: string;
  usuarioNome: string;
  empresas: EmpresaDoSeletor[];
  escopoAcesso: string;
  /** A empresa exige MFA e este usuário ainda não ativou: só serve para ativar na web, trocar de empresa ou sair. */
  mfaPendente: boolean;
  /** Arquivo do banco local desta conta (um por empresa e usuário). */
  arquivoBanco: string;
}

export type AvisoDeSaida = 'sessao_encerrada';

export type EstadoDeAcesso =
  | { tipo: 'iniciando' }
  | { tipo: 'saiu'; aviso: AvisoDeSaida | null }
  | { tipo: 'entrou'; conta: ContaAtiva };

export interface TentativaDeLogin {
  /** URL a abrir no navegador do sistema. */
  url: string;
  verificador: string;
  state: string;
}

export interface OpcoesDoGerente {
  cofre: Cofre;
  cripto: Cripto;
  /** Base da API, sem barra final. */
  apiUrl: string;
  http?: typeof fetch;
  agora?: () => number;
}

interface OpcoesDePedido {
  metodo?: 'GET' | 'POST';
  json?: unknown;
  formulario?: Record<string, string>;
  token?: string;
  cabecalhos?: Record<string, string>;
}

export class GerenteDeSessoes {
  private readonly cofre: Cofre;
  private readonly cripto: Cripto;
  private readonly apiUrl: string;
  private readonly http: typeof fetch;
  private readonly agora: () => number;

  private credencial: Credencial | null = null;
  private estadoAtual: EstadoDeAcesso = { tipo: 'iniciando' };
  private readonly ouvintes = new Set<(e: EstadoDeAcesso) => void>();

  /** Fila de operações que mudam a credencial (renovar, trocar, entrar, sair). */
  private trava: Promise<unknown> = Promise.resolve();
  /** A renovação em andamento, compartilhada por todos que chegam durante ela. */
  private renovacaoEmVoo: Promise<Credencial> | null = null;

  constructor(opcoes: OpcoesDoGerente) {
    this.cofre = opcoes.cofre;
    this.cripto = opcoes.cripto;
    this.apiUrl = opcoes.apiUrl.replace(/\/+$/, '');
    this.http = opcoes.http ?? fetch;
    this.agora = opcoes.agora ?? Date.now;
  }

  // ---- estado -------------------------------------------------------------

  estado(): EstadoDeAcesso {
    return this.estadoAtual;
  }

  assinar(ouvinte: (e: EstadoDeAcesso) => void): () => void {
    this.ouvintes.add(ouvinte);
    return () => this.ouvintes.delete(ouvinte);
  }

  private publicar(e: EstadoDeAcesso): void {
    this.estadoAtual = e;
    for (const o of this.ouvintes) o(e);
  }

  /** Nomes da última conta que entrou, para a tela de login lembrar de quem é a fila guardada. */
  async ultimaConta(): Promise<{ empresaNome: string; usuarioNome: string } | null> {
    const bruto = await this.cofre.ler(CHAVES.ultimaConta);
    return bruto ? (JSON.parse(bruto) as { empresaNome: string; usuarioNome: string }) : null;
  }

  // ---- abertura do app (funciona sem rede) --------------------------------

  /** Lê o cofre e decide o estado. Não faz nenhuma chamada de rede: o app abre offline. */
  async carregar(): Promise<EstadoDeAcesso> {
    const bruta = await this.cofre.ler(CHAVES.credencial);
    const sessao = await this.cofre.ler(CHAVES.sessao);
    if (!bruta || !sessao) {
      this.credencial = null;
      this.publicar({ tipo: 'saiu', aviso: null });
      return this.estadoAtual;
    }
    this.credencial = JSON.parse(bruta) as Credencial;
    this.publicar({ tipo: 'entrou', conta: JSON.parse(sessao) as ContaAtiva });
    return this.estadoAtual;
  }

  // ---- login --------------------------------------------------------------

  /** Gera PKCE e state e monta a URL de autorização. O verificador fica só na memória até a troca do código. */
  async prepararLogin(): Promise<TentativaDeLogin> {
    const p = await gerarPKCE(this.cripto);
    return { url: urlDeAutorizacao(this.apiUrl, p.desafio, p.state), verificador: p.verificador, state: p.state };
  }

  /**
   * Conclui o login com a URL que o navegador devolveu. Lança ErroAcessoNegado
   * para o locatário (access_denied) e ErroRetornoInvalido quando o state não
   * é desta tentativa.
   */
  concluirLogin(tentativa: TentativaDeLogin, urlDeRetorno: string): Promise<ContaAtiva> {
    const retorno = interpretarRetorno(urlDeRetorno, tentativa.state);
    if (retorno.tipo === 'negado') return Promise.reject(new ErroAcessoNegado());
    if (retorno.tipo === 'erro') return Promise.reject(new ErroRetornoInvalido(retorno.motivo));

    return this.exclusivo(async () => {
      const dispositivo = await this.idDoDispositivo();
      const tokens = await this.pedir<TokensDaApi>('/oauth/token', {
        metodo: 'POST',
        formulario: {
          grant_type: 'authorization_code',
          code: retorno.codigo,
          redirect_uri: URI_DE_RETORNO,
          client_id: CLIENT_ID,
          code_verifier: tentativa.verificador,
        },
        cabecalhos: { 'X-Tenaz-Dispositivo-Id': dispositivo },
      });

      let sessao: SessaoDaApi;
      try {
        sessao = await this.pedir<SessaoDaApi>('/v1/auth/sessao', { token: tokens.access_token });
      } catch (e) {
        // O código já foi gasto: sem saber a conta, a credencial não serve.
        await this.revogarOuGuardar(tokens.refresh_token);
        throw e;
      }
      return this.adotar(tokens, sessao);
    });
  }

  // ---- uso autenticado ----------------------------------------------------

  /**
   * Chamada autenticada em nome de UMA conta. O chamador (a sincronização) diz
   * de qual conta é o pedido; se a credencial do aparelho for de outra, o
   * pedido não sai (ErroContaMudou). Renova e repete uma vez em 401, menos em
   * `mfa_pendente`, que não é token vencido.
   */
  autenticado<T>(conta: ChaveDeConta, caminho: string, opcoes: Omit<OpcoesDePedido, 'token'> = {}): Promise<T> {
    return this.autenticadoCom<T>(conta, caminho, opcoes, (usado) => this.renovar(usado));
  }

  private async autenticadoCom<T>(
    conta: ChaveDeConta,
    caminho: string,
    opcoes: Omit<OpcoesDePedido, 'token'>,
    renovador: (tokenUsado: string) => Promise<Credencial>,
  ): Promise<T> {
    let cred = this.exigirCredencial(conta);
    if (cred.expiraEm - this.agora() < FOLGA_DO_TOKEN_MS) {
      cred = await renovador(cred.accessToken);
      this.exigirMesmaConta(cred, conta);
    }
    try {
      return await this.pedir<T>(caminho, { ...opcoes, token: cred.accessToken });
    } catch (e) {
      if (e instanceof ErroApi && e.status === 401 && !e.mfaPendente) {
        const nova = await renovador(cred.accessToken);
        this.exigirMesmaConta(nova, conta);
        return this.pedir<T>(caminho, { ...opcoes, token: nova.accessToken });
      }
      throw e;
    }
  }

  /**
   * Renovação serializada. Chamadas concorrentes compartilham UMA só ida ao
   * servidor; quem chega depois de a renovação terminar, com o token antigo na
   * mão, recebe a credencial já renovada sem nova ida.
   */
  renovar(tokenUsado?: string): Promise<Credencial> {
    if (this.renovacaoEmVoo) return this.renovacaoEmVoo;
    const p: Promise<Credencial> = this.exclusivo(async () => {
      const atual = this.credencial;
      if (!atual) throw new ErroSessaoEncerrada();
      if (tokenUsado !== undefined && atual.accessToken !== tokenUsado) return atual;
      return this.renovarAgora(atual);
    }).finally(() => {
      if (this.renovacaoEmVoo === p) this.renovacaoEmVoo = null;
    });
    this.renovacaoEmVoo = p;
    return p;
  }

  /** Só chamar com a trava na mão. */
  private async renovarAgora(atual: Credencial): Promise<Credencial> {
    let tokens: TokensDaApi;
    try {
      tokens = await this.pedir<TokensDaApi>('/oauth/token', {
        metodo: 'POST',
        formulario: { grant_type: 'refresh_token', refresh_token: atual.refreshToken, client_id: CLIENT_ID },
      });
    } catch (e) {
      if (e instanceof ErroApi && (e.error === 'invalid_grant' || e.status === 401)) {
        // Credencial morta (vencida, revogada ou reusada). Só as credenciais
        // saem: o banco local e a fila ficam como estão.
        await this.apagarCredencial('sessao_encerrada');
        throw new ErroSessaoEncerrada();
      }
      throw e; // rede, 5xx e 429: a sessão continua, tenta-se depois
    }
    const nova = this.montarCredencial(tokens, atual);
    await this.gravarCredencial(nova);
    return nova;
  }

  // ---- empresa ------------------------------------------------------------

  /**
   * Troca de empresa sem sair (POST /v1/auth/trocar-empresa). Precisa de
   * conexão: a credencial antiga é revogada pelo servidor, e a nova é de outra
   * empresa. Segura a trava do começo ao fim, para nenhuma renovação rodar com
   * o refresh token que está sendo trocado. Vale também com `mfa_pendente`.
   */
  trocarEmpresa(empresaId: string): Promise<ContaAtiva> {
    return this.exclusivo(async () => {
      const antes = this.credencial;
      if (!antes) throw new ErroSessaoEncerrada();
      const resposta = await this.autenticadoCom<TrocaDaApi>(
        antes,
        '/v1/auth/trocar-empresa',
        { metodo: 'POST', json: { empresa_escolhida_id: empresaId } },
        async (usado) => {
          const atual = this.credencial;
          if (!atual) throw new ErroSessaoEncerrada();
          return atual.accessToken !== usado ? atual : this.renovarAgora(atual);
        },
      );
      if (!resposta.tokens) throw new ErroRetornoInvalido('a troca de empresa não devolveu tokens.');
      return this.adotar(resposta.tokens, resposta);
    });
  }

  /** Relê a sessão no servidor (por exemplo, depois de a pessoa ativar o MFA na web). */
  async atualizarSessao(): Promise<ContaAtiva> {
    const cred = this.credencial;
    if (!cred) throw new ErroSessaoEncerrada();
    const sessao = await this.autenticado<SessaoDaApi>(cred, '/v1/auth/sessao');
    return this.exclusivo(async () => {
      this.exigirMesmaConta(this.credencial, cred);
      const conta = this.contaDaSessao(sessao);
      await this.cofre.gravar(CHAVES.sessao, JSON.stringify(conta));
      this.publicar({ tipo: 'entrou', conta });
      return conta;
    });
  }

  // ---- saída --------------------------------------------------------------

  /**
   * Sai: apaga as credenciais locais primeiro (sair funciona sem rede) e
   * depois revoga a família no servidor. Sem rede, a revogação fica guardada e
   * é repetida na próxima abertura. NÃO apaga o banco local: a fila de comandos
   * do técnico continua guardada para quando ele entrar de novo.
   */
  sair(): Promise<void> {
    return this.exclusivo(async () => {
      const cred = this.credencial;
      await this.apagarCredencial(null);
      if (cred) await this.revogarOuGuardar(cred.refreshToken);
    });
  }

  /** Repete as revogações que não saíram por falta de rede. Chamar depois de entrar e ao reabrir com sinal. */
  async tentarRevogacoesPendentes(): Promise<void> {
    const pendentes = await this.lerRevogacoes();
    const restantes: string[] = [];
    for (const token of pendentes) {
      try {
        await this.pedir('/oauth/revoke', { metodo: 'POST', formulario: { token, token_type_hint: 'refresh_token' } });
      } catch (e) {
        if (e instanceof ErroDeRede) restantes.push(token);
        // Qualquer resposta do servidor (200 inclusive para token desconhecido) encerra o assunto.
      }
    }
    if (restantes.length === 0) await this.cofre.apagar(CHAVES.revogacoes);
    else await this.cofre.gravar(CHAVES.revogacoes, JSON.stringify(restantes));
  }

  // ---- internos -----------------------------------------------------------

  private exclusivo<T>(fn: () => Promise<T>): Promise<T> {
    const proxima = this.trava.then(fn, fn);
    this.trava = proxima.catch(() => undefined);
    return proxima;
  }

  private exigirCredencial(conta: ChaveDeConta): Credencial {
    if (!this.credencial) throw new ErroSessaoEncerrada();
    this.exigirMesmaConta(this.credencial, conta);
    return this.credencial;
  }

  private exigirMesmaConta(cred: ChaveDeConta | null, conta: ChaveDeConta): void {
    if (!cred || cred.empresaId !== conta.empresaId || cred.usuarioId !== conta.usuarioId) {
      throw new ErroContaMudou();
    }
  }

  private empresaAtual(s: SessaoDaApi): EmpresaDoSeletor {
    const e = (s.empresas ?? []).find((x) => x.atual);
    if (!e) throw new ErroRetornoInvalido('a sessão não informa a empresa atual.');
    return { id: e.id, nome: e.nome };
  }

  private contaDaSessao(s: SessaoDaApi): ContaAtiva {
    const empresa = this.empresaAtual(s);
    return {
      empresaId: empresa.id,
      usuarioId: s.usuario.id,
      empresaNome: s.empresa.nome,
      usuarioNome: s.usuario.nome,
      empresas: (s.empresas ?? []).map((x) => ({ id: x.id, nome: x.nome })),
      escopoAcesso: s.escopo_acesso,
      mfaPendente: s.mfa_pendente,
      arquivoBanco: arquivoDaConta(empresa.id, s.usuario.id),
    };
  }

  private montarCredencial(t: TokensDaApi, conta: ChaveDeConta): Credencial {
    return {
      accessToken: t.access_token,
      expiraEm: this.agora() + t.expires_in * 1000,
      refreshToken: t.refresh_token,
      empresaId: conta.empresaId,
      usuarioId: conta.usuarioId,
    };
  }

  /** Grava credencial e conta no cofre e publica o estado. Só chamar com a trava na mão. */
  private async adotar(tokens: TokensDaApi, sessao: SessaoDaApi): Promise<ContaAtiva> {
    const conta = this.contaDaSessao(sessao);
    const cred = this.montarCredencial(tokens, conta);
    await this.gravarCredencial(cred);
    await this.cofre.gravar(CHAVES.sessao, JSON.stringify(conta));
    await this.cofre.gravar(
      CHAVES.ultimaConta,
      JSON.stringify({ empresaNome: conta.empresaNome, usuarioNome: conta.usuarioNome }),
    );
    this.publicar({ tipo: 'entrou', conta });
    return conta;
  }

  private async gravarCredencial(c: Credencial): Promise<void> {
    await this.cofre.gravar(CHAVES.credencial, JSON.stringify(c));
    this.credencial = c;
  }

  private async apagarCredencial(aviso: AvisoDeSaida | null): Promise<void> {
    this.credencial = null;
    await this.cofre.apagar(CHAVES.credencial);
    await this.cofre.apagar(CHAVES.sessao);
    this.publicar({ tipo: 'saiu', aviso });
  }

  private async lerRevogacoes(): Promise<string[]> {
    const bruto = await this.cofre.ler(CHAVES.revogacoes);
    return bruto ? (JSON.parse(bruto) as string[]) : [];
  }

  private async revogarOuGuardar(refreshToken: string): Promise<void> {
    try {
      await this.pedir('/oauth/revoke', {
        metodo: 'POST',
        formulario: { token: refreshToken, token_type_hint: 'refresh_token' },
      });
    } catch (e) {
      if (e instanceof ErroDeRede) {
        const lista = await this.lerRevogacoes();
        await this.cofre.gravar(CHAVES.revogacoes, JSON.stringify([...lista, refreshToken]));
      }
      // Resposta do servidor: o token já não vale, nada a guardar.
    }
  }

  private async idDoDispositivo(): Promise<string> {
    const existente = await this.cofre.ler(CHAVES.dispositivo);
    if (existente) return existente;
    const novo = uuidV4(this.cripto);
    await this.cofre.gravar(CHAVES.dispositivo, novo);
    return novo;
  }

  /** Uma ida ao servidor. Falha de rede vira ErroDeRede; resposta de erro vira ErroApi. */
  private async pedir<T>(caminho: string, op: OpcoesDePedido = {}): Promise<T> {
    const cabecalhos: Record<string, string> = { Accept: 'application/json, application/problem+json', ...op.cabecalhos };
    let corpo: string | undefined;
    if (op.formulario) {
      cabecalhos['Content-Type'] = 'application/x-www-form-urlencoded';
      corpo = Object.entries(op.formulario)
        .map(([k, v]) => `${encodeURIComponent(k)}=${encodeURIComponent(v)}`)
        .join('&');
    } else if (op.json !== undefined) {
      cabecalhos['Content-Type'] = 'application/json';
      corpo = JSON.stringify(op.json);
    }
    if (op.token) cabecalhos.Authorization = `Bearer ${op.token}`;

    const controle = new AbortController();
    const limite = setTimeout(() => controle.abort(), TEMPO_LIMITE_MS);
    let resposta: Response;
    try {
      resposta = await this.http(`${this.apiUrl}${caminho}`, {
        method: op.metodo ?? 'GET',
        headers: cabecalhos,
        body: corpo,
        signal: controle.signal as never,
      });
    } catch (e) {
      throw new ErroDeRede(e);
    } finally {
      clearTimeout(limite);
    }

    if (resposta.ok) {
      const texto = await resposta.text();
      return (texto ? JSON.parse(texto) : undefined) as T;
    }
    throw await lerErro(resposta);
  }
}

/** Lê o problem+json; corpo que não é JSON (proxy, página de erro) vira erro genérico do status. */
async function lerErro(r: Response): Promise<ErroApi> {
  let p: Record<string, unknown> = {};
  try {
    p = JSON.parse(await r.text()) as Record<string, unknown>;
  } catch {
    // mantém p vazio
  }
  const texto = (v: unknown): string | null => (typeof v === 'string' ? v : null);
  return new ErroApi(
    r.status,
    texto(p.codigo) ?? 'erro_desconhecido',
    texto(p.error),
    p.mfa_pendente === true,
    texto(p.detail) ?? texto(p.title) ?? `Erro ${r.status}`,
    texto(p.id_rastreio),
  );
}
