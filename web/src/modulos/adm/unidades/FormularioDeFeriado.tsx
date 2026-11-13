// Incluir feriado na unidade (RN-ADM-019): um por data, valendo o dia inteiro
// no fuso da unidade. A data é digitada como dd/mm/aaaa.

import { useState } from 'react';
import { useTranslation } from 'react-i18next';

import { Alerta } from '../../../componentes/Alerta';
import { Botao } from '../../../componentes/Botao';
import { CampoTexto } from '../../../componentes/Campo';
import { Formulario } from '../../../componentes/Formulario';
import { tratarErroDeEnvio } from '../../../lib/erros';
import { dataParaISO, mascararData } from '../../../lib/formatacao';
import { textoObrigatorio, useFormulario, z } from '../../../lib/formularios';
import { useIncluirFeriado } from '../unidades';

const esquema = z.object({
  data: z.string().trim().refine((valor) => dataParaISO(valor) !== null, 'validacao.dataInvalida'),
  descricao: textoObrigatorio('validacao.descricaoObrigatoria').max(100, 'validacao.textoLongo'),
});

export function FormularioDeFeriado({ unidadeId, aoConcluir }: { readonly unidadeId: string; readonly aoConcluir: () => void }) {
  const { t } = useTranslation();
  const incluir = useIncluirFeriado();
  const [mensagem, setMensagem] = useState<string | null>(null);
  const formulario = useFormulario(esquema, { defaultValues: { data: '', descricao: '' } });
  const { errors } = formulario.formState;

  const enviar = formulario.handleSubmit((v) => {
    setMensagem(null);
    const data = dataParaISO(v.data);
    if (data === null) return;
    incluir.mutate(
      { unidadeId, corpo: { data, descricao: v.descricao } },
      {
        onSuccess: aoConcluir,
        onError: (erro) => tratarErroDeEnvio(erro, { setError: formulario.setError, definirMensagem: setMensagem }),
      },
    );
  });

  return (
    <Formulario erros={errors} onSubmit={(e) => void enviar(e)} className="flex flex-col gap-4">
      {mensagem && <Alerta tipo="erro">{mensagem}</Alerta>}
      <CampoTexto
        rotulo={t('feriados.data')}
        apoio={t('feriados.dataApoio')}
        inputMode="numeric"
        placeholder={t('feriados.dataExemplo')}
        dado
        autoComplete="off"
        erro={errors.data}
        {...formulario.register('data', { onChange: (e) => (e.target.value = mascararData(e.target.value)) })}
      />
      <CampoTexto rotulo={t('feriados.descricao')} autoComplete="off" erro={errors.descricao} {...formulario.register('descricao')} />
      <div className="flex justify-end gap-2">
        <Botao onClick={aoConcluir}>{t('comum.voltar')}</Botao>
        <Botao type="submit" variante="primario" carregando={incluir.isPending}>
          {incluir.isPending ? t('comum.salvando') : t('feriados.incluir')}
        </Botao>
      </div>
    </Formulario>
  );
}
