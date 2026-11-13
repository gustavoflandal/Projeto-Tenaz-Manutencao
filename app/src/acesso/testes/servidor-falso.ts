// Servidor falso da Tenaz para os testes do acesso: implementa, em memória, o
// contrato de /oauth/token, /oauth/revoke, /v1/auth/sessao e
// /v1/auth/trocar-empresa (seção 4.3 de docs/3_ESPECIFICACAO_APIS.md), inclusive
// o que mais importa para o app: refresh rotativo com detecção de reuso, que
// revoga a família inteira.

import { createHash, randomBytes } from 'node:crypto';

import type { Cripto } from '../pkce';

export const criptoDeTeste: Cripto = {
  bytesAleatorios: (n) => new Uint8Array(randomBytes(n)),
  sha256: async (entrada) => new Uint8Array(createHash('sha256').update(entrada).digest()),
};

export interface EmpresaFalsa {
  id: string;
  nome: string;
}

export interface UsuarioFalso {
  id: string;
  nome: string;
  empresas: EmpresaFalsa[];
  /** Perfil sem mob.usar_app, como o locatário. */
  locatario?: boolean;
  /** Empresas que exigem MFA e onde este usuário não o ativou. */
  mfaPendenteEm?: string[];
}

interface Familia {
  id: number;
  usuario: UsuarioFalso;
  empresaId: string;
  revogada: boolean;
}

export interface RegistroDePedido {
  caminho: string;
  /** Empresa da credencial que o servidor viu; null sem Bearer válido. */
  empresaDoToken: string | null;
}

export class ServidorFalso {
  semRede = false;
  /** Tempo, em ms de relógio falso, que o access token vale. */
  validadeDoAccess = 900;
  relogio = 1_000_000_000_000;
  refreshChamadas = 0;
  pedidos: RegistroDePedido[] = [];
  /** Quando preenchido, a renovação espera esta promessa: permite empilhar chamadas concorrentes. */
  segurarRefresh: Promise<void> | null = null;
  /** Quando preenchido, trocar-empresa espera esta promessa. */
  segurarTroca: Promise<void> | null = null;

  private seq = 0;
  private codigos = new Map<string, { desafio: string; usuario: UsuarioFalso; empresaId: string }>();
  private access = new Map<string, { familia: Familia; expiraEm: number }>();
  private refresh = new Map<string, { familia: Familia; usado: boolean }>();
  private familias: Familia[] = [];
  readonly usuarios = new Map<string, UsuarioFalso>();

  constructor(...usuarios: UsuarioFalso[]) {
    for (const u of usuarios) this.usuarios.set(u.id, u);
  }

  /** Simula o que a web faz em POST /oauth/authorize: devolve a URI de retorno. */
  autorizar(usuarioId: string, empresaId: string, urlDeAutorizacao: string): string {
    const q = new URL(urlDeAutorizacao).searchParams;
    const usuario = this.usuarios.get(usuarioId);
    if (!usuario) throw new Error('usuário desconhecido');
    const state = q.get('state') ?? '';
    if (usuario.locatario) return `br.com.tenaz.app:/oauth?error=access_denied&state=${encodeURIComponent(state)}`;
    const codigo = `cod_${++this.seq}`;
    this.codigos.set(codigo, { desafio: q.get('code_challenge') ?? '', usuario, empresaId });
    return `br.com.tenaz.app:/oauth?code=${codigo}&state=${encodeURIComponent(state)}`;
  }

  /** Vence o access token no servidor (o relógio falso anda). */
  avancar(ms: number): void {
    this.relogio += ms;
  }

  readonly fetch = async (entrada: string | URL | Request, init?: RequestInit): Promise<Response> => {
    if (this.semRede) throw new TypeError('Network request failed');
    const url = new URL(String(entrada));
    const caminho = url.pathname;
    const cab = new Headers(init?.headers);
    const bearer = cab.get('authorization')?.replace(/^Bearer /, '') ?? null;
    const familiaDoToken = bearer ? this.familiaDoAccess(bearer) : null;
    this.pedidos.push({ caminho, empresaDoToken: familiaDoToken?.empresaId ?? null });
    const corpo = typeof init?.body === 'string' ? init.body : '';

    if (caminho === '/oauth/token') return this.token(new URLSearchParams(corpo));
    if (caminho === '/oauth/revoke') {
      const rt = this.refresh.get(new URLSearchParams(corpo).get('token') ?? '');
      if (rt) rt.familia.revogada = true;
      return new Response(null, { status: 200 });
    }
    if (!familiaDoToken) return problema(401, 'credencial_invalida', 'Token inválido ou vencido.');

    const usuario = familiaDoToken.usuario;
    const pendente = usuario.mfaPendenteEm?.includes(familiaDoToken.empresaId) ?? false;

    if (caminho === '/v1/auth/sessao') return json(200, this.sessao(familiaDoToken));
    if (caminho === '/v1/auth/trocar-empresa') {
      if (this.segurarTroca) await this.segurarTroca;
      const alvo = (JSON.parse(corpo) as { empresa_escolhida_id: string }).empresa_escolhida_id;
      if (!usuario.empresas.some((e) => e.id === alvo)) return problema(404, 'nao_encontrado', 'Empresa fora do acesso.');
      familiaDoToken.revogada = true;
      const nova = this.novaFamilia(usuario, alvo);
      const tokens = this.emitir(nova);
      return json(200, { ...this.sessao(nova), tokens });
    }
    if (pendente) return json(401, { type: 'x', title: 'Não autenticado', status: 401, codigo: 'nao_autenticado', mfa_pendente: true }, 'application/problem+json');
    return json(200, { ok: true, empresa: familiaDoToken.empresaId });
  };

  private async token(f: URLSearchParams): Promise<Response> {
    if (f.get('grant_type') === 'authorization_code') {
      const c = this.codigos.get(f.get('code') ?? '');
      this.codigos.delete(f.get('code') ?? '');
      const desafioCalculado = createHash('sha256').update(f.get('code_verifier') ?? '').digest('base64url');
      if (!c || f.get('redirect_uri') !== 'br.com.tenaz.app:/oauth' || f.get('client_id') !== 'tenaz-app' || desafioCalculado !== c.desafio) {
        return invalidGrant();
      }
      return json(200, this.emitir(this.novaFamilia(c.usuario, c.empresaId)));
    }
    this.refreshChamadas += 1;
    if (this.segurarRefresh) await this.segurarRefresh;
    const rt = this.refresh.get(f.get('refresh_token') ?? '');
    if (!rt) return invalidGrant();
    if (rt.familia.revogada) return invalidGrant();
    if (rt.usado) {
      // Reuso: derruba a família inteira (seção 4.3).
      rt.familia.revogada = true;
      return invalidGrant();
    }
    rt.usado = true;
    return json(200, this.emitir(rt.familia));
  }

  private novaFamilia(usuario: UsuarioFalso, empresaId: string): Familia {
    const f = { id: ++this.seq, usuario, empresaId, revogada: false };
    this.familias.push(f);
    return f;
  }

  private emitir(f: Familia) {
    const at = `at_${++this.seq}`;
    const rt = `tnz_rt_${++this.seq}`;
    this.access.set(at, { familia: f, expiraEm: this.relogio + this.validadeDoAccess * 1000 });
    this.refresh.set(rt, { familia: f, usado: false });
    return { access_token: at, token_type: 'Bearer', expires_in: this.validadeDoAccess, refresh_token: rt, scope: 'app_campo' };
  }

  private familiaDoAccess(token: string): Familia | null {
    const a = this.access.get(token);
    if (!a || a.familia.revogada || a.expiraEm <= this.relogio) return null;
    return a.familia;
  }

  private sessao(f: Familia) {
    const atual = f.usuario.empresas.find((e) => e.id === f.empresaId);
    return {
      canal: 'app',
      csrf_token: null,
      expira_em: null,
      escopo_acesso: 'interno',
      usuario: { id: f.usuario.id, nome: f.usuario.nome, email: 'x@exemplo.com.br', pessoa_id: f.usuario.id, mfa_ativo: false },
      empresa: { nome: atual?.nome ?? '', plano: 'profissional', srv_habilitado: false },
      perfis: ['tecnico'],
      permissoes: [],
      empresas: [...f.usuario.empresas]
        .sort((a, b) => Number(b.id === f.empresaId) - Number(a.id === f.empresaId))
        .map((e) => ({ id: e.id, nome: e.nome, atual: e.id === f.empresaId })),
      mfa_pendente: f.usuario.mfaPendenteEm?.includes(f.empresaId) ?? false,
    };
  }
}

function json(status: number, corpo: unknown, tipo = 'application/json'): Response {
  return new Response(JSON.stringify(corpo), { status, headers: { 'Content-Type': tipo } });
}

function problema(status: number, codigo: string, detail: string): Response {
  return json(status, { type: 'x', title: codigo, status, codigo, detail }, 'application/problem+json');
}

function invalidGrant(): Response {
  return json(
    401,
    { type: 'x', title: 'Credencial inválida', status: 401, codigo: 'credencial_invalida', error: 'invalid_grant', detail: 'Código ou refresh token inválido.' },
    'application/problem+json',
  );
}
