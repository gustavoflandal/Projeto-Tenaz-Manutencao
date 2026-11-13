// Estados de tela (seções 6.14 e 6.15 do documento 7): vazio, erro de carga,
// sem permissão e carregando.

import { CloudOff, Lock, Search, type LucideIcon } from 'lucide-react';
import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';

import { ehSemPermissao, mensagemDoErro } from '../lib/erros';

import { Botao } from './Botao';

export interface EstadoVazioProps {
  readonly icone?: LucideIcon;
  readonly titulo: string;
  readonly texto?: string;
  readonly acao?: ReactNode;
}

export function EstadoVazio({ icone: Icone = Search, titulo, texto, acao }: EstadoVazioProps) {
  return (
    <div className="flex flex-col items-center gap-3 px-6 py-12 text-center">
      <Icone size={48} aria-hidden="true" className="text-texto-terciario" />
      <h2 className="texto-h3">{titulo}</h2>
      {texto && <p className="max-w-120 texto-corpo text-texto-secundario">{texto}</p>}
      {acao}
    </div>
  );
}

export interface EstadoErroProps {
  readonly erro: unknown;
  readonly aoTentarDeNovo?: () => void;
}

/** Erro de carga: a frase do erro e "Tentar de novo". */
export function EstadoErro({ erro, aoTentarDeNovo }: EstadoErroProps) {
  const { t } = useTranslation();
  if (ehSemPermissao(erro)) return <EstadoSemPermissao />;
  return (
    <div role="alert" className="flex flex-col items-center gap-3 px-6 py-12 text-center">
      <CloudOff size={48} aria-hidden="true" className="text-texto-terciario" />
      <p className="max-w-120 texto-corpo">{mensagemDoErro(erro)}</p>
      {aoTentarDeNovo && <Botao onClick={aoTentarDeNovo}>{t('comum.tentarDeNovo')}</Botao>}
    </div>
  );
}

export function EstadoSemPermissao({ texto }: { readonly texto?: string }) {
  const { t } = useTranslation();
  return <EstadoVazio icone={Lock} titulo={t('comum.semPermissaoTitulo')} texto={texto ?? t('comum.semPermissao')} />;
}

/** Esqueleto de carregamento: contêiner `aria-busy` com texto oculto (6.15). */
export function Carregando({ rotulo, linhas = 5 }: { readonly rotulo: string; readonly linhas?: number }) {
  return (
    <div aria-busy="true" className="flex flex-col gap-3 py-2">
      <span className="sr-only" role="status">
        {rotulo}
      </span>
      {Array.from({ length: linhas }, (_, i) => (
        <div key={i} aria-hidden="true" className="h-10 rounded-md bg-superficie-esqueleto" />
      ))}
    </div>
  );
}
