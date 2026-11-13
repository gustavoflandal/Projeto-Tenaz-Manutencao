// Aviso de conflito de versão (412, documento 3, seção 2.11): o registro foi
// alterado por outra pessoa e a tela mostra a versão atual. O que a pessoa tinha
// digitado e não salvou aparece em lista, campo e valor, para digitar de novo.

import { useTranslation } from 'react-i18next';

import type { AlteracaoNaoSalva } from '../lib/erros';

import { Alerta } from './Alerta';

export function AlertaDeConflito({ alteracoes, className }: { readonly alteracoes: readonly AlteracaoNaoSalva[]; readonly className?: string }) {
  const { t } = useTranslation();
  return (
    <Alerta tipo="atencao" rolar {...(className ? { className } : {})}>
      <p>{t('erro.conflitoDeVersao')}</p>
      {alteracoes.length > 0 && (
        <>
          <p className="mt-2">{t('erro.alteracoesNaoSalvas')}</p>
          <ul className="mt-1 list-disc pl-5">
            {alteracoes.map((a) => (
              <li key={a.campo}>{t('erro.alteracaoItem', { campo: a.campo, valor: a.valor })}</li>
            ))}
          </ul>
        </>
      )}
    </Alerta>
  );
}
