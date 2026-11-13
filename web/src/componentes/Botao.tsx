// Botão (componente 6.1 do documento 7), tamanho M (40 px).
//
// Primário: Sinal com texto Aço, uma vez por tela ou modal. Desabilitado em
// Aço 100 com texto Concreto; quem desabilita explica o motivo ao lado.

import type { ComponentProps } from 'react';

export type VarianteBotao = 'primario' | 'secundario' | 'terciario' | 'icone';

const BASE =
  'inline-flex shrink-0 items-center justify-center gap-2 rounded-md texto-botao transition-colors duration-(--tz-duracao-rapida) ease-padrao disabled:cursor-not-allowed';

const VARIANTES: Record<VarianteBotao, string> = {
  primario:
    'h-10 px-4 bg-acao-primaria text-acao-primaria-texto enabled:hover:bg-acao-primaria-hover enabled:active:bg-acao-primaria-pressionada disabled:bg-borda disabled:text-texto-desabilitado',
  secundario:
    'h-10 px-4 border border-borda-controle bg-transparent text-texto enabled:hover:border-borda-controle-hover enabled:hover:bg-superficie-hover enabled:active:bg-superficie-selecionada disabled:border-borda disabled:text-texto-desabilitado',
  terciario:
    'h-10 px-2 bg-transparent text-texto enabled:hover:underline disabled:text-texto-desabilitado',
  icone:
    'size-10 bg-transparent text-texto-secundario enabled:hover:bg-superficie-hover enabled:hover:text-texto disabled:text-texto-desabilitado',
};

export interface BotaoProps extends ComponentProps<'button'> {
  readonly variante?: VarianteBotao;
}

export function Botao({ variante = 'secundario', type = 'button', className, ...resto }: BotaoProps) {
  const classes = [BASE, VARIANTES[variante], className].filter(Boolean).join(' ');
  return <button type={type} className={classes} {...resto} />;
}
