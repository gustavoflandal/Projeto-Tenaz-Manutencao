// Tratamento de erro da API: do código estável (seção 3.2 do documento 3) à
// mensagem em português, e dos erros por campo (`erros`) ao formulário.
//
// Formato das mensagens: o que aconteceu e o que fazer (seção 10.3 do
// documento 7). Os textos ficam em @tenaz/i18n; o servidor só escolhe o código.

import type { FieldValues, Path, UseFormSetError } from 'react-hook-form';

import { ErroDaAPI, type Problema } from './api';
import { listar } from './formatacao';
import { existeTexto, traduzirChave } from './i18n';

export { ErroDaAPI };

export function ehErroDaAPI(erro: unknown): erro is ErroDaAPI {
  return erro instanceof ErroDaAPI;
}

/** Falha de rede: o navegador não chegou ao servidor. */
export function ehFalhaDeRede(erro: unknown): boolean {
  return erro instanceof TypeError;
}

/** Tipo de erro de campo vindo do servidor: a mensagem já está em português. */
export const TIPO_ERRO_DO_SERVIDOR = 'servidor';

/**
 * Mensagem para um erro que não é de campo. Usa o texto do código; sem texto
 * para o código, o status; e por último o `detail` do servidor.
 */
export function mensagemDoErro(erro: unknown): string {
  if (ehFalhaDeRede(erro)) return traduzirChave('erro.semConexao');
  if (!ehErroDaAPI(erro)) return traduzirChave('erro.inesperado');

  const { codigo, status, problema } = erro;
  // Erro de validação que não coube em nenhum campo: o servidor diz o que está errado.
  if (status === 422 && problema?.erros?.length) return problema.erros.map((e) => e.mensagem).join(' ');
  if (codigo === 'mfa_codigo_invalido' && problema?.tentativas_restantes !== undefined) {
    return traduzirChave('erro.codigos.mfa_codigo_invalido_restantes', { count: problema.tentativas_restantes });
  }
  if (codigo === 'transicao_invalida' && problema?.estado_atual === 'exigido_pela_empresa') {
    return traduzirChave('erro.codigos.mfa_exigido_pela_empresa');
  }
  if (codigo && existeTexto(`erro.codigos.${codigo}`)) return traduzirChave(`erro.codigos.${codigo}`);
  if (status !== undefined && status >= 500) {
    return traduzirChave('erro.servidor', { rastreio: problema?.id_rastreio ?? '-' });
  }
  if (status !== undefined && existeTexto(`erro.status.${status}`)) return traduzirChave(`erro.status.${status}`);
  return problema?.detail ?? traduzirChave('erro.inesperado');
}

export function ehConflitoDeVersao(erro: unknown): boolean {
  return ehErroDaAPI(erro) && erro.status === 412;
}

/** Erro 404: o registro não existe, é de outra empresa ou está fora do alcance. */
export function ehNaoEncontrado(erro: unknown): boolean {
  return ehErroDaAPI(erro) && erro.status === 404;
}

export function ehSemPermissao(erro: unknown): boolean {
  return ehErroDaAPI(erro) && erro.status === 403 && erro.codigo !== 'csrf_invalido';
}

/** Notação do servidor (`turnos[0].hora_fim`) para a do formulário (`turnos.0.hora_fim`). */
export function caminhoDoCampo(campo: string): string {
  return campo.replaceAll(/\[(\d+)\]/g, '.$1');
}

export interface OpcoesDeErrosDoServidor {
  /** Caminho do objeto no formulário, quando o corpo da API está dentro dele (`unidade`, no cadastro da conta). */
  readonly prefixo?: string;
  /** Nome do campo no servidor para o nome no formulário (`senha` do servidor é `nova_senha` na redefinição). */
  readonly mapa?: Readonly<Record<string, string>>;
  /** Tradução geral do nome do campo, para formulários de nomes planos (`unidade.codigo` vira `unidade_codigo`). */
  readonly transformar?: (campo: string) => string;
}

/**
 * Leva os erros por campo do servidor (422 validacao_falhou, 409 duplicado)
 * para o formulário. Devolve `true` se algum campo recebeu erro.
 */
export function aplicarErrosDoServidor<T extends FieldValues>(
  erro: unknown,
  setError: UseFormSetError<T>,
  { prefixo = '', mapa = {}, transformar }: OpcoesDeErrosDoServidor = {},
): boolean {
  if (!ehErroDaAPI(erro) || !erro.problema?.erros?.length) return false;
  for (const e of erro.problema.erros) {
    const nome = transformar ? transformar(e.campo) : (mapa[e.campo] ?? e.campo);
    const caminho = caminhoDoCampo(prefixo ? `${prefixo}.${nome}` : nome);
    setError(caminho as Path<T>, { type: TIPO_ERRO_DO_SERVIDOR, message: e.mensagem });
  }
  return true;
}

/** Problema com a lista de erros por campo (para os testes e o resumo). */
export function errosPorCampo(erro: unknown): NonNullable<Problema['erros']> {
  return ehErroDaAPI(erro) ? (erro.problema?.erros ?? []) : [];
}

/** Campo que a pessoa alterou e ainda não salvou, com o valor digitado. */
export interface AlteracaoNaoSalva {
  readonly campo: string;
  readonly valor: string;
}

/**
 * Alterações digitadas e não salvas de um formulário: só os campos sujos, com o
 * rótulo da tela e o valor. Serve ao aviso de conflito de versão (412), que
 * recarrega o registro e descarta o que foi digitado.
 */
export function alteracoesNaoSalvas(
  valores: object,
  sujos: object,
  rotulos: Readonly<Record<string, string>>,
): AlteracaoNaoSalva[] {
  const resultado: AlteracaoNaoSalva[] = [];
  for (const [campo, valor] of Object.entries(valores)) {
    const sujo = (sujos as Record<string, unknown>)[campo];
    const mudou = sujo === true || (Array.isArray(sujo) && sujo.some(Boolean));
    if (!mudou) continue;
    let texto: string;
    if (Array.isArray(valor)) texto = listar(valor.map(String));
    else if (typeof valor === 'boolean') texto = traduzirChave(valor ? 'comum.sim' : 'comum.nao');
    else texto = String(valor ?? '');
    resultado.push({ campo: rotulos[campo] ?? campo, valor: texto === '' ? traduzirChave('comum.vazio') : texto });
  }
  return resultado;
}

export interface TratamentoDeEnvio<T extends FieldValues> {
  readonly setError: UseFormSetError<T>;
  readonly definirMensagem: (mensagem: string | null) => void;
  /** 412 versao_desatualizada: quem salva recarrega o recurso e avisa, listando o que foi digitado. */
  readonly aoConflitoDeVersao?: (alteracoes: readonly AlteracaoNaoSalva[]) => void;
  /** Alterações digitadas, para o aviso do conflito; calculada só se houver conflito. */
  readonly alteracoes?: () => readonly AlteracaoNaoSalva[];
  readonly opcoes?: OpcoesDeErrosDoServidor;
}

/**
 * Caminho comum de uma tentativa de salvar que falhou: conflito de versão,
 * erro por campo (mostrado no campo) ou mensagem geral no topo do formulário.
 */
export function tratarErroDeEnvio<T extends FieldValues>(erro: unknown, tratamento: TratamentoDeEnvio<T>): void {
  if (ehConflitoDeVersao(erro) && tratamento.aoConflitoDeVersao) {
    tratamento.aoConflitoDeVersao(tratamento.alteracoes?.() ?? []);
    return;
  }
  if (aplicarErrosDoServidor(erro, tratamento.setError, tratamento.opcoes)) {
    tratamento.definirMensagem(null);
    return;
  }
  tratamento.definirMensagem(mensagemDoErro(erro));
}
