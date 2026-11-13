// Formulários: react-hook-form com esquema zod. A mensagem de cada regra do
// esquema é uma chave de @tenaz/i18n; quem mostra o erro a traduz. A validação
// roda ao sair do campo e ao salvar, nunca a cada tecla (seção 6.2 do documento 7).

import { zodResolver } from '@hookform/resolvers/zod';
import {
  cpaConfereComSubdivisao,
  ErroDeDinheiro,
  identificacaoValida,
  lerPrecoUnitario,
  moedaDoCodigo,
  senhaTemTamanho,
  validarCodigoPostal,
  type Moeda,
  type Pais,
} from '@tenaz/dominio';
import { useForm, type FieldValues, type Resolver, type UseFormProps } from 'react-hook-form';
import { z } from 'zod';

export { z };

export function useFormulario<Esquema extends z.ZodType<FieldValues, FieldValues>>(
  esquema: Esquema,
  opcoes?: Omit<UseFormProps<z.input<Esquema>, unknown, z.output<Esquema>>, 'resolver'>,
) {
  return useForm<z.input<Esquema>, unknown, z.output<Esquema>>({
    resolver: zodResolver(esquema) as unknown as Resolver<z.input<Esquema>, unknown, z.output<Esquema>>,
    mode: 'onBlur',
    reValidateMode: 'onBlur',
    ...opcoes,
  });
}

// ------------------------------------------------------- regras reutilizáveis

/** Texto obrigatório, aparado. */
export const textoObrigatorio = (mensagem = 'validacao.obrigatorio') => z.string().trim().min(1, mensagem);

export const emailObrigatorio = z
  .string()
  .trim()
  .min(1, 'validacao.emailObrigatorio')
  .pipe(z.email('validacao.emailInvalido'));

/** E-mail que pode ficar vazio. */
export const emailOpcional = z
  .string()
  .trim()
  .refine((valor) => valor === '' || z.email().safeParse(valor).success, 'validacao.emailInvalido');

/**
 * Confere a identificação de um tipo (o do país) dentro de um `superRefine`.
 * A mensagem é a chave do tipo em `validacao.identificacaoInvalida`; o servidor
 * repete a conferência e tem a palavra final.
 */
export function conferirIdentificacao(
  ctx: z.RefinementCtx,
  tipo: string,
  valor: string,
  opcoes: { readonly obrigatoria: boolean; readonly caminho?: string },
): void {
  const caminho = opcoes.caminho ?? 'identificacao';
  if (valor.trim() === '') {
    if (opcoes.obrigatoria) ctx.addIssue({ code: 'custom', path: [caminho], message: 'validacao.identificacaoObrigatoria' });
    return;
  }
  if (!identificacaoValida(tipo, valor)) {
    ctx.addIssue({ code: 'custom', path: [caminho], message: `validacao.identificacaoInvalida.${tipo}` });
  }
}

/**
 * Preço unitário digitado (sem sinal, até seis casas) para o texto decimal do
 * contrato, com ponto. A vírgula decimal é aceita na digitação; a leitura é a do
 * `dinheiro.ts` (ADR-013, Parte F). Devolve `null` quando a entrada é recusada.
 * Sem a moeda (usuário ainda não carregou), vale a gramática do preço, que é a
 * mesma em qualquer moeda.
 */
export function precoUnitarioParaOContrato(texto: string, moeda: Moeda | null): string | null {
  const aparado = texto.trim();
  const canonico = aparado.includes('.') ? aparado : aparado.replace(',', '.');
  try {
    lerPrecoUnitario(canonico, moedaDoCodigo(moeda ?? 'BRL'));
  } catch (erro) {
    if (erro instanceof ErroDeDinheiro && erro.codigo === 'formato_invalido') return null;
    throw erro;
  }
  return canonico;
}

/** Custo por hora opcional: preço unitário sem sinal, na moeda da empresa. */
export const custoHoraOpcional = (moeda: Moeda | null) =>
  z.string().refine((texto) => texto.trim() === '' || precoUnitarioParaOContrato(texto, moeda) !== null, 'validacao.custoHora');

/**
 * Confere o código postal pelo país da empresa e, na Argentina, a letra do CPA
 * contra a província (documento 3, seção 2.18). Os dois campos são opcionais; o
 * servidor repete a conferência e tem a palavra final.
 */
export function conferirEndereco(
  ctx: z.RefinementCtx,
  pais: Pais | null,
  campos: { readonly subdivisao: string; readonly codigoPostal: string },
  caminhoDoCodigoPostal = 'codigo_postal',
): void {
  const texto = campos.codigoPostal.trim();
  if (pais === null || texto === '') return;
  const lido = validarCodigoPostal(pais, texto);
  if (!lido.valido) {
    ctx.addIssue({ code: 'custom', path: [caminhoDoCodigoPostal], message: `validacao.codigoPostalInvalido.${pais}` });
  } else if (pais === 'AR' && !cpaConfereComSubdivisao(campos.subdivisao, lido.normalizado)) {
    ctx.addIssue({ code: 'custom', path: [caminhoDoCodigoPostal], message: 'validacao.cpaNaoConfere' });
  }
}

export const senhaNova = z.string().refine(senhaTemTamanho, 'validacao.senhaCurta');

/** Texto opcional: vazio vira `undefined` no envio. */
export function vazioParaIndefinido(valor: string | undefined): string | undefined {
  const aparado = valor?.trim();
  return aparado ? aparado : undefined;
}

/** Texto opcional em PATCH: vazio vira `null`, que apaga o valor no servidor. */
export function vazioParaNulo(valor: string | undefined): string | null {
  return vazioParaIndefinido(valor) ?? null;
}
