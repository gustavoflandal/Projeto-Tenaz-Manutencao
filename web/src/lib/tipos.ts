// Nomes curtos para os tipos do contrato. Nenhum tipo de API é escrito à mão:
// tudo vem de @tenaz/api-tipos, gerado de backend/api/openapi.yaml.

import type { components } from '@tenaz/api-tipos';

type Esquemas = components['schemas'];

export type Sessao = Esquemas['Sessao'];
export type Empresa = Esquemas['Empresa'];
export type EmpresaAcesso = Esquemas['EmpresaAcesso'];
export type ResumoUnidade = Esquemas['ResumoUnidade'];
export type Unidade = Esquemas['Unidade'];
export type Turno = Esquemas['Turno'];
export type Feriado = Esquemas['Feriado'];
export type Pessoa = Esquemas['Pessoa'];
export type ResumoPessoa = Esquemas['ResumoPessoa'];
export type Usuario = Esquemas['Usuario'];
export type ResumoUsuario = Esquemas['ResumoUsuario'];
export type Perfil = Esquemas['Perfil'];
export type Permissao = Esquemas['Permissao'];
export type ResumoEvento = Esquemas['ResumoEvento'];
export type Evento = Esquemas['Evento'];
export type DiaVerificado = Esquemas['DiaVerificado'];
export type CapturaAuditoria = Esquemas['CapturaAuditoria'];
export type Paginacao = Esquemas['Paginacao'];

export type NovaEmpresa = Esquemas['NovaEmpresa'];
export type EntradaUnidade = Esquemas['EntradaUnidade'];
export type EntradaPessoa = Esquemas['EntradaPessoa'];
export type EntradaUsuario = Esquemas['EntradaUsuario'];

export type SituacaoPessoa = Pessoa['situacao'];
export type SituacaoUsuario = Usuario['situacao'];
export type Vinculo = NonNullable<Pessoa['vinculo']>;

/** Página de uma lista paginada por cursor (seção 2.6 do documento 3). */
export interface PaginaDeLista<T> {
  readonly dados: readonly T[] | null;
  readonly paginacao: Paginacao;
}
