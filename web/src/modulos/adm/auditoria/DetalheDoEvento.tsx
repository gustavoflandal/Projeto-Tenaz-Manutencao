// Detalhe de um evento da trilha (RF-ADM-011): quem, quando, de onde, em qual
// tela, o hash do conteúdo e os valores antes e depois. Abrir o detalhe é uma
// visualização, gravada pelo servidor quando a captura da tela de origem está ligada.

import { useTranslation } from 'react-i18next';

import { Carregando, EstadoErro } from '../../../componentes/Estados';
import { listar } from '../../../lib/formatacao';
import { traduzirOuOriginal } from '../../../lib/i18n';
import { useEvento } from '../auditoria';
import { useDataHora } from '../formatos';

import { DiferencaDoEvento } from './DiferencaDoEvento';

export function DetalheDoEvento({ id }: { readonly id: string }) {
  const { t } = useTranslation();
  const consulta = useEvento(id);
  const dataHora = useDataHora();
  if (consulta.isPending) return <Carregando rotulo={t('auditoria.carregandoEvento')} linhas={3} />;
  if (consulta.isError) return <EstadoErro erro={consulta.error} aoTentarDeNovo={() => void consulta.refetch()} />;
  const e = consulta.data;
  const detalhe = e.detalhe === null || e.detalhe === undefined ? null : JSON.stringify(e.detalhe, null, 2);
  const linhas: ReadonlyArray<readonly [string, string, boolean]> = [
    [t('auditoria.quando'), dataHora(e.ocorrido_em) ?? e.ocorrido_em, true],
    [t('auditoria.acao'), traduzirOuOriginal('auditoria.acoes', e.acao), false],
    [t('auditoria.resultado'), t(`auditoria.resultados.${e.resultado}`), false],
    [t('auditoria.usuario'), e.usuario_nome || e.usuario_id || t(`auditoria.atores.${e.ator_tipo}`), false],
    [t('auditoria.tela'), e.tela ?? e.modulo ?? t('comum.naoInformado'), true],
    [t('auditoria.entidade'), listar([e.entidade, e.entidade_chave]), false],
    [t('auditoria.sensivel'), e.sensivel ? t('comum.sim') : t('comum.nao'), false],
    [t('auditoria.ip'), e.ip ?? t('comum.naoInformado'), true],
    [t('auditoria.origem'), e.origem ? t(`auditoria.origens.${e.origem}`) : t('comum.naoInformado'), false],
    [t('auditoria.requisicao'), e.request_id ?? t('comum.naoInformado'), true],
    [t('auditoria.hash'), e.hash, true],
  ];
  return (
    <div className="flex flex-col gap-4">
      <dl className="flex flex-col gap-2">
        {linhas.map(([rotulo, valor, dado]) => (
          <div key={rotulo} className="flex items-baseline gap-4">
            <dt className="w-28 shrink-0 texto-sobrelinha text-texto-secundario">{rotulo}</dt>
            <dd className={`min-w-0 break-all ${dado ? 'texto-dado' : 'texto-corpo-p'}`}>{valor}</dd>
          </div>
        ))}
      </dl>
      {detalhe && (
        <section aria-labelledby={`detalhe-${e.id}`}>
          <h3 id={`detalhe-${e.id}`} className="mb-1 texto-sobrelinha text-texto-secundario">
            {t('auditoria.detalhe')}
          </h3>
          <pre className="overflow-x-auto rounded-md bg-superficie-selecionada p-3 texto-dado">{detalhe}</pre>
        </section>
      )}
      <DiferencaDoEvento antes={e.antes} depois={e.depois} camposAlterados={e.campos_alterados ?? []} />
    </div>
  );
}
