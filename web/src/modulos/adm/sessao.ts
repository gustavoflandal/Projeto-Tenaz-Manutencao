// Sessão e autenticação (RF-ADM-006, seção 4.2 do documento 3).
//
// A sessão é dado do servidor: vive só no cache do TanStack Query (cookie
// httpOnly no navegador; o csrf_token sai dela). Não há estado global paralelo.

import { useMutation, useQuery, useQueryClient, type QueryClient } from '@tanstack/react-query';
import type { components } from '@tenaz/api-tipos';

import { api, ErroDaAPI, executar, obter, type Problema } from '../../lib/api';
import { chaves } from '../../lib/consultas';
import type { Sessao } from '../../lib/tipos';

type Esquemas = components['schemas'];

// ------------------------------------------------------------------- leitura

async function lerSessao(): Promise<Sessao | null> {
  const { data, response, error } = await api.GET('/v1/auth/sessao');
  if (response.status === 401) return null;
  if (!response.ok || data === undefined) {
    throw new ErroDaAPI(`GET /v1/auth/sessao devolveu ${response.status}`, response.status, error as Problema | undefined);
  }
  return data;
}

/** Sessão atual; `null` quando não há sessão (não é erro). */
export function useSessao() {
  return useQuery({
    queryKey: chaves.sessao,
    queryFn: lerSessao,
    staleTime: 5 * 60_000,
    retry: false,
  });
}

/** A sessão atual tem a permissão? Esconder botão não é segurança: o servidor decide. */
export function useTemPermissao(codigo: string): boolean {
  const { data } = useSessao();
  return data?.permissoes?.includes(codigo) ?? false;
}

export function temPermissao(sessao: Sessao | null | undefined, codigo: string): boolean {
  return sessao?.permissoes?.includes(codigo) ?? false;
}

/** Guarda a sessão nova e descarta todo dado de outra sessão ou empresa. */
export function aplicarSessao(cliente: QueryClient, sessao: Sessao): void {
  cliente.removeQueries({ predicate: (consulta) => consulta.queryKey[0] !== chaves.sessao[0] });
  cliente.setQueryData(chaves.sessao, sessao);
}

// -------------------------------------------------------------------- entrar

export function useEntrar() {
  const cliente = useQueryClient();
  return useMutation({
    mutationFn: (entrada: Esquemas['EntradaLoginBody']) => obter(api.POST('/v1/auth/login', { body: entrada })),
    onSuccess: (sessao) => aplicarSessao(cliente, sessao),
  });
}

export function useVerificarMFA() {
  const cliente = useQueryClient();
  return useMutation({
    mutationFn: (entrada: Esquemas['EntradaMFABody']) => obter(api.POST('/v1/auth/mfa/verificar', { body: entrada })),
    onSuccess: (sessao) => aplicarSessao(cliente, sessao),
  });
}

export function useSair() {
  const cliente = useQueryClient();
  return useMutation({
    mutationFn: () => executar(api.POST('/v1/auth/logout')),
    onSettled: () => {
      // Mesmo que a chamada falhe, a tela deixa de mostrar dado da sessão.
      cliente.removeQueries({ predicate: (consulta) => consulta.queryKey[0] !== chaves.sessao[0] });
      cliente.setQueryData(chaves.sessao, null);
    },
  });
}

export function useTrocarEmpresa() {
  const cliente = useQueryClient();
  return useMutation({
    mutationFn: (empresaId: string) =>
      obter(api.POST('/v1/auth/trocar-empresa', { body: { empresa_escolhida_id: empresaId } })),
    onSuccess: (sessao) => aplicarSessao(cliente, sessao),
  });
}

// ------------------------------------------------------ senha e convite

export function useRecuperarSenha() {
  return useMutation({
    mutationFn: (email: string) => executar(api.POST('/v1/auth/recuperar-senha', { body: { email } })),
  });
}

export function useRedefinirSenha() {
  return useMutation({
    mutationFn: (entrada: Esquemas['EntradaRedefinirSenhaBody']) =>
      executar(api.POST('/v1/auth/redefinir-senha', { body: entrada })),
  });
}

export function useConsultarConvite(token: string | null) {
  return useQuery({
    queryKey: chaves.convite,
    queryFn: () => obter(api.POST('/v1/auth/convites/consultar', { body: { token: token ?? '' } })),
    enabled: token !== null,
    retry: false,
    // O token é segredo de uso único: nada dele fica no cache além da tela.
    gcTime: 0,
    staleTime: Infinity,
  });
}

export function useAceitarConvite() {
  return useMutation({
    mutationFn: (entrada: Esquemas['EntradaConviteAceitarBody']) =>
      obter(api.POST('/v1/auth/convites/aceitar', { body: entrada })),
  });
}

// ------------------------------------------------------ conta do usuário

export function useTrocarSenha() {
  return useMutation({
    mutationFn: (entrada: Esquemas['EntradaTrocarSenhaBody']) =>
      executar(api.POST('/v1/auth/trocar-senha', { body: entrada })),
  });
}

export function useIniciarMFA() {
  return useMutation({ mutationFn: () => obter(api.POST('/v1/auth/mfa/iniciar')) });
}

export function useConfirmarMFA() {
  const cliente = useQueryClient();
  return useMutation({
    mutationFn: (codigo: string) => executar(api.POST('/v1/auth/mfa/confirmar', { body: { codigo } })),
    onSuccess: () => cliente.invalidateQueries({ queryKey: chaves.sessao }),
  });
}

export function useDesativarMFA() {
  const cliente = useQueryClient();
  return useMutation({
    mutationFn: (codigo: string) => executar(api.POST('/v1/auth/mfa/desativar', { body: { codigo } })),
    onSuccess: () => cliente.invalidateQueries({ queryKey: chaves.sessao }),
  });
}

// ------------------------------------------------------------ cadastro e app

export function useCriarEmpresa() {
  return useMutation({
    mutationFn: (entrada: Esquemas['NovaEmpresa']) => obter(api.POST('/v1/empresas', { body: entrada })),
  });
}

export function useAutorizarApp() {
  return useMutation({
    mutationFn: (entrada: Esquemas['EntradaAutorizarPOSTBody']) => obter(api.POST('/oauth/authorize', { body: entrada })),
  });
}
