// Botão (componente 6.1 do documento 7), tamanho M (40 px).
//
// Primário: Sinal com texto Aço, uma vez por tela ou modal. Desabilitado em
// Aço 100 com texto Concreto; quem desabilita explica o motivo ao lado.
// Carregando: mantém a largura, mostra o indicador e ignora o clique; quem usa
// troca o texto pelo gerúndio ("Salvando").

import { Loader2 } from 'lucide-react';
import type { ComponentProps, MouseEvent } from 'react';

export type VarianteBotao = 'primario' | 'secundario' | 'terciario' | 'perigo' | 'icone';

const BASE =
  'inline-flex shrink-0 items-center justify-center gap-2 rounded-md texto-botao transition-colors duration-(--tz-duracao-rapida) ease-padrao disabled:cursor-not-allowed aria-busy:cursor-progress';

const VARIANTES: Record<VarianteBotao, string> = {
  primario:
    'h-10 px-4 bg-acao-primaria text-acao-primaria-texto enabled:hover:bg-acao-primaria-hover enabled:active:bg-acao-primaria-pressionada disabled:bg-borda disabled:text-texto-desabilitado',
  secundario:
    'h-10 px-4 border border-borda-controle bg-transparent text-texto enabled:hover:border-borda-controle-hover enabled:hover:bg-superficie-hover enabled:active:bg-superficie-selecionada disabled:border-borda disabled:text-texto-desabilitado',
  terciario:
    'h-10 px-2 bg-transparent text-texto enabled:hover:underline disabled:text-texto-desabilitado',
  perigo:
    'h-10 px-4 bg-perigo text-perigo-texto enabled:hover:bg-perigo-hover disabled:bg-borda disabled:text-texto-desabilitado',
  icone:
    'size-10 bg-transparent text-texto-secundario enabled:hover:bg-superficie-hover enabled:hover:text-texto disabled:text-texto-desabilitado',
};

export interface BotaoProps extends ComponentProps<'button'> {
  readonly variante?: VarianteBotao;
  /** Envio em andamento: `aria-busy`, indicador giratório e clique ignorado. */
  readonly carregando?: boolean;
}

export function Botao({
  variante = 'secundario',
  type = 'button',
  className,
  carregando = false,
  children,
  onClick,
  ...resto
}: BotaoProps) {
  const classes = [BASE, VARIANTES[variante], className].filter(Boolean).join(' ');
  const aoClicar = (evento: MouseEvent<HTMLButtonElement>) => {
    if (carregando) {
      evento.preventDefault();
      return;
    }
    onClick?.(evento);
  };
  return (
    <button type={type} className={classes} aria-busy={carregando || undefined} onClick={aoClicar} {...resto}>
      {carregando && <Loader2 size={16} aria-hidden="true" className="motion-safe:animate-spin" />}
      {children}
    </button>
  );
}
