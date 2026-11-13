// Selo de versão (componente 6.19 do documento 7, RNF-039).
//
// Fixo no canto inferior direito, abaixo de menus, modais e toasts (camada
// `z-fixo`). Fundo opaco, para manter o contraste se algo rolar por baixo.
// Deve ser o último elemento focável da página: quem usa o selo o coloca no
// fim do layout.

import { Dialog } from 'radix-ui';
import { useState } from 'react';

import { nomeAcessivelDoSelo, seloDeVersao } from '../../lib/versao';

import { DialogoSobre } from './DialogoSobre';

export function SeloVersao() {
  const [aberto, setAberto] = useState(false);
  return (
    <Dialog.Root open={aberto} onOpenChange={setAberto}>
      <Dialog.Trigger asChild>
        <button
          type="button"
          aria-label={nomeAcessivelDoSelo()}
          className="fixed right-2 bottom-2 z-(--tz-z-fixo) inline-flex min-h-6 items-center rounded-xs bg-fundo px-1.5 py-0.5 texto-selo text-texto-terciario hover:text-texto"
        >
          {seloDeVersao()}
        </button>
      </Dialog.Trigger>
      {aberto && <DialogoSobre />}
    </Dialog.Root>
  );
}
