// Casca das telas sem sessão (entrar, esqueci a senha, convite, cadastro...):
// fundo `fundo`, logo, cartão `superficie` e o selo de versão no canto inferior
// direito, fora do cartão (seção 5.6 do documento 7, RNF-039).

import { useEffect, useRef, type ReactNode } from 'react';

import { Logo } from '../../../componentes/Logo';
import { LinksDeIdioma } from '../../../componentes/SeletorDeIdioma';
import { SeloVersao } from '../../../componentes/versao/SeloVersao';
import { useTema } from '../../../lib/tema';

export interface CartaoPublicoProps {
  readonly titulo: string;
  readonly children: ReactNode;
  /** Cartão largo, para o cadastro da conta. */
  readonly largo?: boolean;
}

export function CartaoPublico({ titulo, children, largo }: CartaoPublicoProps) {
  const { resolvido } = useTema();
  const titulo1 = useRef<HTMLHeadingElement>(null);
  // Cada passo da tela (senha, código, empresa) é um cartão novo: o foco vai para o título dele.
  useEffect(() => titulo1.current?.focus(), []);
  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-fundo px-4 pt-16 pb-12 text-texto">
      <main id="conteudo" className={`flex w-full flex-col items-center gap-8 ${largo ? 'max-w-180' : 'max-w-100'}`}>
        <Logo forma="horizontal" variante={resolvido === 'escuro' ? 'negativo' : 'positivo'} className="h-10 w-auto" />
        <div className="w-full rounded-lg border border-borda bg-superficie p-5 shadow-1">
          <h1 ref={titulo1} tabIndex={-1} className="texto-h2">{titulo}</h1>
          <div className="mt-6 flex flex-col gap-4">{children}</div>
        </div>
        <LinksDeIdioma />
      </main>

      <footer>
        <SeloVersao />
      </footer>
    </div>
  );
}
