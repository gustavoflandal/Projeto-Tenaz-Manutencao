// Verificação de integridade da trilha (RF-ADM-011): para cada dia UTC do
// período, o servidor recalcula os hashes, o selo e a cadeia e compara com a
// cópia imutável. Estados: íntegro, violado (com motivo e o primeiro evento
// alterado) e ainda não selado. Se o armazenamento do selo não responde (503),
// nada se afirma: sem a cópia não há integridade comprovada.

import { useState, type FormEvent } from 'react';
import { useTranslation } from 'react-i18next';
import { useSearchParams } from 'react-router';

import { Alerta } from '../../../componentes/Alerta';
import { Botao } from '../../../componentes/Botao';
import { CabecalhoPagina } from '../../../componentes/CabecalhoPagina';
import { CampoTexto, type ErroDeCampo } from '../../../componentes/Campo';
import { ChipVerificacao } from '../../../componentes/ChipStatus';
import { Modal } from '../../../componentes/Dialogos';
import { Carregando, EstadoErro, EstadoVazio } from '../../../componentes/Estados';
import { Tabela, Td, Th, Tr } from '../../../componentes/Tabela';
import { dataParaISO, mascararData } from '../../../lib/formatacao';
import { traduzirChave, traduzirOuOriginal } from '../../../lib/i18n';
import { useVerificacao } from '../auditoria';
import { useFormatos } from '../formatos';

import { DetalheDoEvento } from './DetalheDoEvento';

const DIAS_MAXIMOS = 400;
const DIA_EM_MS = 86_400_000;

function diasEntre(deISO: string, ateISO: string): number {
  return Math.round((Date.parse(`${ateISO}T00:00:00Z`) - Date.parse(`${deISO}T00:00:00Z`)) / DIA_EM_MS) + 1;
}

/** Período padrão do formulário: os últimos 30 dias UTC, até hoje. */
function periodoPadrao(): { de: string; ate: string } {
  const hoje = new Date();
  const ate = hoje.toISOString().slice(0, 10);
  const de = new Date(hoje.getTime() - 29 * DIA_EM_MS).toISOString().slice(0, 10);
  return { de, ate };
}

const ROTULO_DO_RESUMO = {
  integro: 'integridade.integros',
  violado: 'integridade.violados',
  nao_selado: 'integridade.naoSelados',
} as const;

export function TelaIntegridade() {
  const { t } = useTranslation();
  const { data: formatarData } = useFormatos();
  const [parametros, definirParametros] = useSearchParams();
  const padrao = periodoPadrao();
  const de = parametros.get('de');
  const ate = parametros.get('ate');
  const consulta = useVerificacao(de, ate);
  const [campoDe, setCampoDe] = useState(formatarData(de ?? padrao.de));
  const [campoAte, setCampoAte] = useState(formatarData(ate ?? padrao.ate));
  const [erros, setErros] = useState<{ de?: string; ate?: string }>({});
  const [evento, setEvento] = useState<string | null>(null);

  const verificar = (e: FormEvent) => {
    e.preventDefault();
    const deISO = dataParaISO(campoDe);
    const ateISO = dataParaISO(campoAte);
    const novos: { de?: string; ate?: string } = {};
    if (deISO === null) novos.de = traduzirChave('validacao.dataInvalida');
    if (ateISO === null) novos.ate = traduzirChave('validacao.dataInvalida');
    if (deISO && ateISO) {
      if (deISO > ateISO) novos.ate = t('auditoria.periodoInvertido');
      else if (diasEntre(deISO, ateISO) > DIAS_MAXIMOS) novos.ate = t('integridade.periodoMaximo');
    }
    setErros(novos);
    if (Object.keys(novos).length > 0 || !deISO || !ateISO) return;
    definirParametros({ de: deISO, ate: ateISO }, { replace: true });
  };

  const erro = (mensagem: string | undefined): ErroDeCampo | undefined =>
    mensagem ? { type: 'servidor', message: mensagem } : undefined;
  const dias = consulta.data?.dados ?? [];
  const contar = (estado: string) => dias.filter((d) => d.estado === estado).length;

  return (
    <>
      <CabecalhoPagina
        titulo={t('integridade.titulo')}
        trilha={[{ rotulo: t('auditoria.titulo'), para: '/administracao/auditoria' }, { rotulo: t('integridade.titulo') }]}
      />
      <p className="mb-4 max-w-180 texto-corpo text-texto-secundario">{t('integridade.texto')}</p>

      <form noValidate role="search" aria-label={t('integridade.titulo')} onSubmit={verificar} className="mb-6 flex flex-wrap items-start gap-3">
        <div className="w-40">
          <CampoTexto
            rotulo={t('integridade.de')}
            name="de"
            dado
            inputMode="numeric"
            placeholder={t('feriados.dataExemplo')}
            value={campoDe}
            onChange={(ev) => setCampoDe(mascararData(ev.target.value))}
            erro={erro(erros.de)}
          />
        </div>
        <div className="w-40">
          <CampoTexto
            rotulo={t('integridade.ate')}
            name="ate"
            dado
            inputMode="numeric"
            placeholder={t('feriados.dataExemplo')}
            value={campoAte}
            onChange={(ev) => setCampoAte(mascararData(ev.target.value))}
            erro={erro(erros.ate)}
          />
        </div>
        <div className="pt-7">
          <Botao type="submit" variante="primario" carregando={consulta.isFetching}>
            {consulta.isFetching ? t('integridade.verificando') : t('integridade.verificar')}
          </Botao>
        </div>
      </form>

      {consulta.isFetching && <Carregando rotulo={t('integridade.carregando')} linhas={3} />}
      {consulta.isError && !consulta.isFetching && (
        <EstadoErro erro={consulta.error} aoTentarDeNovo={() => void consulta.refetch()} />
      )}
      {consulta.isSuccess && !consulta.isFetching && (
        <>
          {contar('violado') > 0 && (
            <Alerta tipo="erro" className="mb-4">
              {t('integridade.violadoAviso')}
            </Alerta>
          )}
          <dl aria-label={t('integridade.resumo')} className="mb-4 flex flex-wrap gap-6">
            {(['integro', 'violado', 'nao_selado'] as const).map((estado) => (
              <div key={estado}>
                <dt className="texto-sobrelinha text-texto-secundario">{t(ROTULO_DO_RESUMO[estado])}</dt>
                <dd className="texto-kpi">{contar(estado)}</dd>
              </div>
            ))}
          </dl>
          {dias.length === 0 ? (
            <EstadoVazio titulo={t('integridade.vazio')} />
          ) : (
            <Tabela legenda={t('integridade.titulo')}>
              <thead>
                <tr>
                  <Th>{t('integridade.dia')}</Th>
                  <Th>{t('integridade.estado')}</Th>
                  <Th>{t('integridade.motivo')}</Th>
                  <Th>{t('integridade.evento')}</Th>
                </tr>
              </thead>
              <tbody>
                {dias.map((d) => (
                  <Tr key={d.dia}>
                    <Td className="texto-dado">{formatarData(d.dia)}</Td>
                    <Td>
                      <ChipVerificacao estado={d.estado} />
                    </Td>
                    <Td>{d.motivo ? traduzirOuOriginal('integridade.motivos', d.motivo) : t('comum.naoInformado')}</Td>
                    <Td>
                      {d.evento_id ? (
                        <Botao
                          variante="terciario"
                          aria-label={t('integridade.verEventoDe', { dia: formatarData(d.dia) })}
                          onClick={() => setEvento(d.evento_id ?? null)}
                        >
                          {t('integridade.verEvento')}
                        </Botao>
                      ) : (
                        t('comum.naoInformado')
                      )}
                    </Td>
                  </Tr>
                ))}
              </tbody>
            </Tabela>
          )}
        </>
      )}

      <Modal aberto={evento !== null} aoMudar={(a) => !a && setEvento(null)} titulo={t('auditoria.detalheTitulo')}>
        {evento !== null && <DetalheDoEvento id={evento} />}
      </Modal>
    </>
  );
}
