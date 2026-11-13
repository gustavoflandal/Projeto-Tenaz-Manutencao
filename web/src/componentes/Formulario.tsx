// Estrutura comum dos formulários (seção 6.2 do documento 7): resumo de erros
// no topo com um link para cada campo, e o identificador que liga cada campo
// ao seu link. A validação roda ao sair do campo e ao salvar, nunca a cada
// tecla; ao salvar com erro, o foco vai para o primeiro campo com erro.

import { CircleX } from 'lucide-react';
import { useId, type ComponentProps } from 'react';
import type { FieldErrors } from 'react-hook-form';
import { useTranslation } from 'react-i18next';

import { ProvedorDoFormulario, idDoCampo, textoDoErro, type ErroDeCampo } from './Campo';

export interface ItemDeErro {
  readonly nome: string;
  readonly erro: ErroDeCampo;
}

function ehErroDeCampo(valor: unknown): valor is ErroDeCampo & { readonly message: string } {
  return typeof valor === 'object' && valor !== null && 'message' in valor && typeof valor.message === 'string';
}

/** Achata os erros do formulário (aninhados) numa lista de nome e mensagem. */
export function listarErros(erros: FieldErrors | object, prefixo = ''): ItemDeErro[] {
  const itens: ItemDeErro[] = [];
  for (const [chave, valor] of Object.entries(erros)) {
    if (chave === 'ref' || typeof valor !== 'object' || valor === null) continue;
    const nome = prefixo ? `${prefixo}.${chave}` : chave;
    if (ehErroDeCampo(valor)) itens.push({ nome, erro: valor });
    itens.push(...listarErros(valor, nome));
  }
  return itens;
}

export interface FormularioProps extends Omit<ComponentProps<'form'>, 'noValidate'> {
  readonly erros?: FieldErrors | object;
  /** Rótulo legível por campo, para o resumo; sem ele o link usa a mensagem. */
  readonly rotulos?: Readonly<Record<string, string>>;
  /** Identificador estável do formulário; sem ele, um gerado. */
  readonly idDoFormulario?: string;
}

export function Formulario({ erros = {}, rotulos, idDoFormulario, children, ...resto }: FormularioProps) {
  const gerado = useId();
  const id = idDoFormulario ?? gerado;
  const itens = listarErros(erros);
  return (
    <ProvedorDoFormulario value={id}>
      <form noValidate {...resto}>
        {itens.length > 0 && <ResumoDeErros idDoFormulario={id} itens={itens} rotulos={rotulos} />}
        {children}
      </form>
    </ProvedorDoFormulario>
  );
}

function ResumoDeErros({
  idDoFormulario,
  itens,
  rotulos,
}: {
  readonly idDoFormulario: string;
  readonly itens: readonly ItemDeErro[];
  readonly rotulos: Readonly<Record<string, string>> | undefined;
}) {
  const { t } = useTranslation();
  return (
    <div role="alert" className="mb-4 rounded-md border-l-4 border-status-parado bg-status-parado-fundo p-3 text-texto">
      <p className="flex items-center gap-2 texto-h4">
        <CircleX size={20} aria-hidden="true" className="text-status-parado" />
        {t('formulario.resumoTitulo', { count: itens.length })}
      </p>
      <ul className="mt-2 flex flex-col gap-1">
        {itens.map(({ nome, erro }) => {
          const mensagem = textoDoErro(erro) ?? '';
          const rotulo = rotulos?.[nome];
          return (
            <li key={nome} className="texto-corpo-p">
              <a href={`#${idDoCampo(idDoFormulario, nome)}`} className="underline">
                {rotulo ? t('formulario.itemDeErro', { campo: rotulo, mensagem }) : mensagem}
              </a>
            </li>
          );
        })}
      </ul>
    </div>
  );
}
