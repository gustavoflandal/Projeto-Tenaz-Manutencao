// Usuários, perfis e permissões (RF-ADM-004, RF-ADM-005).

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { components } from '@tenaz/api-tipos';

import { api, executar, obter, seVersao } from '../../lib/api';
import { chaves } from '../../lib/consultas';
import { LIMITE_DA_PAGINA, useListaPorCursor } from '../../lib/lista';

type Esquemas = components['schemas'];

export interface FiltroDeUsuarios {
  readonly busca?: string;
  readonly situacao?: 'convidado' | 'ativo' | 'inativo';
  readonly perfil?: string;
}

export function useUsuarios(filtro: FiltroDeUsuarios) {
  return useListaPorCursor(chaves.usuarios(filtro), (cursor) =>
    obter(
      api.GET('/v1/usuarios', {
        params: { query: { ...filtro, limite: LIMITE_DA_PAGINA, ...(cursor === undefined ? {} : { cursor }) } },
      }),
    ),
  );
}

export function useUsuario(id: string) {
  return useQuery({
    queryKey: chaves.usuario(id),
    queryFn: () => obter(api.GET('/v1/usuarios/{id}', { params: { path: { id } } })),
  });
}

function useInvalidarUsuarios() {
  const cliente = useQueryClient();
  return () => void cliente.invalidateQueries({ queryKey: chaves.usuarios() });
}

export function useCriarUsuario() {
  const invalidar = useInvalidarUsuarios();
  const cliente = useQueryClient();
  return useMutation({
    mutationFn: (corpo: Esquemas['EntradaUsuario']) => obter(api.POST('/v1/usuarios', { body: corpo })),
    onSuccess: () => {
      invalidar();
      // A pessoa criada junto entra na lista de pessoas.
      void cliente.invalidateQueries({ queryKey: chaves.pessoas() });
    },
  });
}

export function useAlterarUsuario() {
  const invalidar = useInvalidarUsuarios();
  return useMutation({
    mutationFn: ({ id, versao, corpo }: { id: string; versao: number; corpo: Esquemas['PatchUsuario'] }) =>
      obter(
        api.PATCH('/v1/usuarios/{id}', { params: { path: { id }, header: { 'If-Match': seVersao(versao) } }, body: corpo }),
      ),
    onSuccess: invalidar,
  });
}

export function useDesativarUsuario() {
  const invalidar = useInvalidarUsuarios();
  return useMutation({
    mutationFn: ({ id, versao }: { id: string; versao: number }) =>
      obter(
        api.POST('/v1/usuarios/{id}/desativar', { params: { path: { id }, header: { 'If-Match': seVersao(versao) } } }),
      ),
    onSuccess: invalidar,
  });
}

export function useReativarUsuario() {
  const invalidar = useInvalidarUsuarios();
  return useMutation({
    mutationFn: ({ id, versao }: { id: string; versao: number }) =>
      obter(
        api.POST('/v1/usuarios/{id}/reativar', { params: { path: { id }, header: { 'If-Match': seVersao(versao) } } }),
      ),
    onSuccess: invalidar,
  });
}

export function useReenviarConvite() {
  return useMutation({
    mutationFn: (id: string) => executar(api.POST('/v1/usuarios/{id}/reenviar-convite', { params: { path: { id } } })),
  });
}

export function useEncerrarSessoes() {
  return useMutation({
    mutationFn: (id: string) => obter(api.POST('/v1/usuarios/{id}/encerrar-sessoes', { params: { path: { id } } })),
  });
}

export function useConcederPermissao() {
  const invalidar = useInvalidarUsuarios();
  return useMutation({
    mutationFn: ({ id, permissao }: { id: string; permissao: string }) =>
      obter(api.POST('/v1/usuarios/{id}/permissoes', { params: { path: { id } }, body: { permissao } })),
    onSuccess: invalidar,
  });
}

export function useRetirarPermissao() {
  const invalidar = useInvalidarUsuarios();
  return useMutation({
    mutationFn: ({ id, permissao }: { id: string; permissao: string }) =>
      executar(api.DELETE('/v1/usuarios/{id}/permissoes/{permissao}', { params: { path: { id, permissao } } })),
    onSuccess: invalidar,
  });
}

// ---------------------------------------------------------- perfis e permissões

export function usePerfis() {
  return useQuery({ queryKey: chaves.perfis, queryFn: () => obter(api.GET('/v1/perfis')) });
}

/** Catálogo de permissões; exige `adm.gerenciar_perfis`. */
export function usePermissoes(habilitado: boolean) {
  return useQuery({
    queryKey: chaves.permissoes,
    queryFn: () => obter(api.GET('/v1/permissoes')),
    enabled: habilitado,
    staleTime: 10 * 60_000,
  });
}

// ------------------------------------------------------------ o próprio usuário

export const chaveUsuarioEu = ['usuarios', 'eu'] as const;

export function lerUsuarioEu() {
  return obter(api.GET('/v1/usuarios/eu'));
}

export function useUsuarioEu() {
  return useQuery({ queryKey: chaveUsuarioEu, queryFn: lerUsuarioEu });
}

export function useAlterarUsuarioEu() {
  const cliente = useQueryClient();
  return useMutation({
    mutationFn: ({ versao, corpo }: { versao: number; corpo: Esquemas['PatchEu'] }) =>
      obter(api.PATCH('/v1/usuarios/eu', { params: { header: { 'If-Match': seVersao(versao) } }, body: corpo })),
    onSuccess: (usuario) => {
      cliente.setQueryData(chaveUsuarioEu, usuario);
      void cliente.invalidateQueries({ queryKey: chaves.sessao });
    },
  });
}
