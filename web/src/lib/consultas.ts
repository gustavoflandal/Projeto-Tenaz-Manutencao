// Cliente do TanStack Query: a única camada que guarda dado do servidor na
// web. Nenhum estado global mantém cópia de dado do servidor (seção 8.4 de
// docs/5_ARQUITETURA_SOFTWARE.md).

import { MutationCache, QueryCache, QueryClient } from '@tanstack/react-query';

import { registrarFonteDoCSRF } from './api';
import { ehErroDaAPI } from './erros';

/** Chaves de consulta, num só lugar, para invalidação previsível. */
export const chaves = {
  versaoDaAPI: ['plataforma', 'versao'] as const,
  sessao: ['sessao'] as const,
  /** Marca que a sessão acabou por 401, e não por saída pedida pela pessoa. */
  sessaoEncerrada: ['sessao-encerrada'] as const,
  empresa: ['empresa'] as const,
  unidades: ['unidades'] as const,
  unidade: (id: string) => ['unidades', id] as const,
  feriados: (unidadeId: string) => ['unidades', unidadeId, 'feriados'] as const,
  pessoas: (filtros?: object) => (filtros ? (['pessoas', 'lista', filtros] as const) : (['pessoas'] as const)),
  pessoa: (id: string) => ['pessoas', 'detalhe', id] as const,
  usuarios: (filtros?: object) => (filtros ? (['usuarios', 'lista', filtros] as const) : (['usuarios'] as const)),
  usuario: (id: string) => ['usuarios', 'detalhe', id] as const,
  perfis: ['perfis'] as const,
  permissoes: ['permissoes'] as const,
  eventos: (filtros?: object) => (filtros ? (['eventos', 'lista', filtros] as const) : (['eventos'] as const)),
  evento: (id: string) => ['eventos', 'detalhe', id] as const,
  verificacao: (de: string, ate: string) => ['verificacao-auditoria', de, ate] as const,
  capturas: ['capturas-auditoria'] as const,
  convite: ['convite'] as const,
  configuracaoMFA: ['mfa', 'configuracao'] as const,
};

/**
 * Um 401 numa chamada que exige sessão significa sessão perdida: o cache da
 * sessão é esvaziado e a guarda de rotas manda a pessoa entrar de novo. Fora
 * disso ficam o login, o segundo fator e a troca de senha, onde 401 é erro de
 * credencial digitada. Com `mfa_pendente`, a sessão continua e só fica restrita.
 */
function tratarSessaoPerdida(cliente: QueryClient, erro: unknown): void {
  if (!ehErroDaAPI(erro) || erro.status !== 401) return;
  if (erro.problema?.mfa_pendente === true) {
    void cliente.invalidateQueries({ queryKey: chaves.sessao });
    return;
  }
  const instancia = erro.problema?.instance ?? '';
  const ehRotaDeCredencial = /^\/v1\/auth\/(login|mfa\/(verificar|confirmar|desativar)|trocar-senha)/.test(instancia);
  if (ehRotaDeCredencial) return;
  if (erro.codigo !== 'nao_autenticado' && erro.codigo !== 'credencial_invalida') return;
  cliente.setQueryData(chaves.sessaoEncerrada, true);
  cliente.setQueryData(chaves.sessao, null);
}

export function criarClienteDeConsultas(): QueryClient {
  const cliente: QueryClient = new QueryClient({
    queryCache: new QueryCache({ onError: (erro) => tratarSessaoPerdida(cliente, erro) }),
    mutationCache: new MutationCache({ onError: (erro) => tratarSessaoPerdida(cliente, erro) }),
    defaultOptions: {
      queries: {
        // Erro de cliente (4xx) não melhora tentando de novo; falha de rede e 5xx, uma vez.
        retry: (tentativas, erro) => tentativas < 1 && !(ehErroDaAPI(erro) && (erro.status ?? 500) < 500),
        refetchOnWindowFocus: false,
      },
    },
  });
  return cliente;
}

/**
 * Liga o token CSRF das chamadas à sessão guardada neste cliente de consultas,
 * sem cópia em outro lugar. Chamado por quem monta a aplicação (num efeito, e
 * não na criação do cliente: o StrictMode cria o cliente duas vezes).
 */
export function ligarCSRFAoCliente(cliente: QueryClient): void {
  registrarFonteDoCSRF(() => {
    const sessao = cliente.getQueryData<{ csrf_token?: string | null } | null>(chaves.sessao);
    return sessao?.csrf_token;
  });
}
