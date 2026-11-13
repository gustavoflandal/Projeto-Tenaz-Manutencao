// Modais e confirmações (componente 6.13 do documento 7).
//
// Confirmação só para ação irreversível ou de grande alcance. O título é a
// pergunta com o objeto, o texto diz a consequência, o botão repete o verbo e
// o que sai sem agir é "Voltar", nunca "Cancelar". O foco inicial vai para o
// botão menos destrutivo.

import { X } from 'lucide-react';
import { AlertDialog, Dialog } from 'radix-ui';
import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';

import { Botao } from './Botao';

const VEU = 'fixed inset-0 z-(--tz-z-veu) bg-veu';
const CONTEUDO =
  'fixed inset-x-4 top-1/2 z-(--tz-z-modal) mx-auto -translate-y-1/2 rounded-lg border border-borda bg-superficie-elevada p-6 text-texto shadow-4';

export interface ConfirmacaoProps {
  readonly aberto: boolean;
  readonly aoMudar: (aberto: boolean) => void;
  readonly titulo: string;
  readonly texto: string;
  readonly rotuloConfirmar: string;
  /** Ação destrutiva: botão de perigo. Ação importante, mas não destrutiva: primário. */
  readonly perigo?: boolean;
  readonly carregando?: boolean;
  readonly aoConfirmar: () => void;
  /** Mensagem de erro da tentativa anterior, mostrada dentro do diálogo. */
  readonly erro?: ReactNode;
}

export function Confirmacao({
  aberto,
  aoMudar,
  titulo,
  texto,
  rotuloConfirmar,
  perigo,
  carregando,
  aoConfirmar,
  erro,
}: ConfirmacaoProps) {
  const { t } = useTranslation();
  return (
    <AlertDialog.Root open={aberto} onOpenChange={aoMudar}>
      <AlertDialog.Portal>
        <AlertDialog.Overlay className={VEU} />
        <AlertDialog.Content className={`${CONTEUDO} max-w-100`}>
          <AlertDialog.Title className="texto-h3">{titulo}</AlertDialog.Title>
          <AlertDialog.Description className="mt-3 texto-corpo text-texto-secundario">{texto}</AlertDialog.Description>
          {erro && <div className="mt-4">{erro}</div>}
          <div className="mt-6 flex justify-end gap-2">
            <AlertDialog.Cancel asChild>
              <Botao variante="secundario" disabled={carregando}>
                {t('comum.voltar')}
              </Botao>
            </AlertDialog.Cancel>
            {/* Sem AlertDialog.Action: ele fecharia o diálogo antes do fim da chamada. */}
            <Botao variante={perigo ? 'perigo' : 'primario'} carregando={carregando} onClick={aoConfirmar}>
              {rotuloConfirmar}
            </Botao>
          </div>
        </AlertDialog.Content>
      </AlertDialog.Portal>
    </AlertDialog.Root>
  );
}

export interface ModalProps {
  readonly aberto: boolean;
  readonly aoMudar: (aberto: boolean) => void;
  readonly titulo: string;
  readonly descricao?: string;
  readonly children: ReactNode;
}

/** Modal de formulário curto (560 px). O foco fica preso dentro e Esc fecha. */
export function Modal({ aberto, aoMudar, titulo, descricao, children }: ModalProps) {
  const { t } = useTranslation();
  return (
    <Dialog.Root open={aberto} onOpenChange={aoMudar}>
      <Dialog.Portal>
        <Dialog.Overlay className={VEU} />
        <Dialog.Content className={`${CONTEUDO} max-h-screen max-w-140 overflow-y-auto`}>
          <div className="flex items-start justify-between gap-4">
            <Dialog.Title className="texto-h3">{titulo}</Dialog.Title>
            <Dialog.Close asChild>
              <Botao variante="icone" aria-label={t('comum.fechar')} className="-mt-2 -mr-2">
                <X size={20} aria-hidden="true" />
              </Botao>
            </Dialog.Close>
          </div>
          <Dialog.Description className={descricao ? 'mt-2 texto-corpo text-texto-secundario' : 'sr-only'}>
            {descricao ?? titulo}
          </Dialog.Description>
          <div className="mt-4">{children}</div>
        </Dialog.Content>
      </Dialog.Portal>
    </Dialog.Root>
  );
}
