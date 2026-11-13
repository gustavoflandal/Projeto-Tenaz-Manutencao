// Empresa (RF-ADM-001): dados cadastrais, CNPJ editável pelo `administrador` e
// o interruptor do MFA obrigatório (RF-ADM-006).

import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query';
import type { components } from '@tenaz/api-tipos';

import { api, obter, seVersao } from '../../lib/api';
import { chaves } from '../../lib/consultas';

export function useEmpresa() {
  return useQuery({ queryKey: chaves.empresa, queryFn: () => obter(api.GET('/v1/empresa')) });
}

export function useAlterarEmpresa() {
  const cliente = useQueryClient();
  return useMutation({
    mutationFn: ({ versao, corpo }: { versao: number; corpo: components['schemas']['PatchEmpresa'] }) =>
      obter(api.PATCH('/v1/empresa', { params: { header: { 'If-Match': seVersao(versao) } }, body: corpo })),
    onSuccess: (empresa) => {
      cliente.setQueryData(chaves.empresa, empresa);
      // Ligar o MFA obrigatório pode deixar a sessão atual restrita (RN-ADM-027).
      void cliente.invalidateQueries({ queryKey: chaves.sessao });
    },
  });
}
