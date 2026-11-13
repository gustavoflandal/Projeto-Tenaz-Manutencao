// Diálogo "Sobre" (componente 6.19 do documento 7, RNF-039).

import { Copy, Info, X } from 'lucide-react';
import { Dialog } from 'radix-ui';
import { useTranslation } from 'react-i18next';

import { Botao } from '../Botao';

import { useCopiar, useInformacoesVersao } from './useInformacoesVersao';

/** Conteúdo do diálogo; montado só quando aberto, e então consulta a API. */
export function DialogoSobre() {
  const { t } = useTranslation();
  const { linhas, estadoDaAPI, versaoNovaDisponivel, textoParaCopiar } = useInformacoesVersao();
  const { estado: estadoCopia, copiar } = useCopiar();

  let mensagemCopia = '';
  if (estadoCopia === 'copiado') mensagemCopia = t('sobre.copiado');
  else if (estadoCopia === 'falhou') mensagemCopia = t('sobre.copiaFalhou');

  return (
    <Dialog.Portal>
      <Dialog.Overlay className="fixed inset-0 z-(--tz-z-veu) bg-veu" />
      <Dialog.Content className="fixed inset-x-4 top-1/2 z-(--tz-z-modal) mx-auto max-w-100 -translate-y-1/2 rounded-lg border border-borda bg-superficie-elevada p-6 text-texto shadow-4">
        <div className="flex items-start justify-between gap-4">
          <Dialog.Title className="texto-h3">{t('sobre.titulo')}</Dialog.Title>
          <Dialog.Close asChild>
            <Botao variante="icone" aria-label={t('sobre.fechar')} className="-mt-2 -mr-2">
              <X size={20} aria-hidden="true" />
            </Botao>
          </Dialog.Close>
        </div>
        <Dialog.Description className="sr-only">{t('sobre.descricao')}</Dialog.Description>

        {versaoNovaDisponivel && (
          <div
            role="status"
            className="mt-4 flex items-start gap-3 rounded-md border-l-4 border-texto bg-superficie p-3"
          >
            <Info size={20} aria-hidden="true" className="mt-0.5 shrink-0" />
            <div className="flex flex-col items-start gap-2">
              <p className="texto-corpo-p">{t('sobre.versaoNova')}</p>
              <Botao variante="secundario" onClick={() => globalThis.location.reload()}>
                {t('sobre.recarregar')}
              </Botao>
            </div>
          </div>
        )}

        <dl className="mt-6 flex flex-col gap-3" aria-busy={estadoDaAPI === 'consultando'}>
          {linhas.map((linha) => (
            <div key={linha.chave} className="flex items-baseline gap-4">
              <dt className="w-32 shrink-0 texto-sobrelinha text-texto-secundario">{linha.rotulo}</dt>
              <dd className="min-w-0 texto-dado break-words" data-linha={linha.chave}>
                {linha.valor}
              </dd>
            </div>
          ))}
        </dl>

        <p role="status" className="mt-4 min-h-5 texto-legenda text-texto-secundario">
          {mensagemCopia}
        </p>

        <div className="mt-4 flex items-center justify-between gap-2">
          <Botao variante="secundario" onClick={() => void copiar(textoParaCopiar)}>
            <Copy size={20} aria-hidden="true" />
            {t('sobre.copiar')}
          </Botao>
          <Dialog.Close asChild>
            <Botao variante="secundario">{t('sobre.fechar')}</Botao>
          </Dialog.Close>
        </div>
      </Dialog.Content>
    </Dialog.Portal>
  );
}
