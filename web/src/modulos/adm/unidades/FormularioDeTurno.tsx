// Incluir ou alterar um turno da unidade (RN-ADM-019). A sobreposição com os
// outros turnos é conferida pelo servidor, que aponta o turno em conflito.

import { validarTurno } from '@tenaz/dominio';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';

import { Alerta } from '../../../componentes/Alerta';
import { AlertaDeConflito } from '../../../componentes/AlertaDeConflito';
import { Botao } from '../../../componentes/Botao';
import { CampoSelecao, CampoTexto } from '../../../componentes/Campo';
import { Formulario } from '../../../componentes/Formulario';
import { alteracoesNaoSalvas, tratarErroDeEnvio, type AlteracaoNaoSalva } from '../../../lib/erros';
import { useFormulario, vazioParaNulo, z } from '../../../lib/formularios';
import type { Turno } from '../../../lib/tipos';
import { useAlterarTurno, useIncluirTurno } from '../unidades';

const NOMES_DOS_DIAS = ['segunda', 'terca', 'quarta', 'quinta', 'sexta', 'sabado', 'domingo'] as const;

const esquema = z
  .object({
    dia_semana: z.string().min(1, 'validacao.obrigatorio'),
    nome: z.string().trim().max(60, 'validacao.nomeDoTurnoLongo'),
    hora_inicio: z.string().trim(),
    hora_fim: z.string().trim(),
  })
  .superRefine((v, contexto) => {
    const erro = validarTurno(v);
    if (erro === 'hora_inicio_invalida') contexto.addIssue({ code: 'custom', path: ['hora_inicio'], message: 'validacao.horaInicio' });
    if (erro === 'hora_fim_invalida') contexto.addIssue({ code: 'custom', path: ['hora_fim'], message: 'validacao.horaFim' });
    if (erro === 'fim_igual_ao_inicio') contexto.addIssue({ code: 'custom', path: ['hora_fim'], message: 'validacao.fimIgualAoInicio' });
  });

export interface FormularioDeTurnoProps {
  readonly unidadeId: string;
  /** Ausente: inclui. Presente: altera, com a versão do próprio turno. */
  readonly turno?: Turno;
  readonly aoConcluir: () => void;
}

export function FormularioDeTurno({ unidadeId, turno, aoConcluir }: FormularioDeTurnoProps) {
  const { t } = useTranslation();
  const incluir = useIncluirTurno();
  const alterar = useAlterarTurno();
  const [mensagem, setMensagem] = useState<string | null>(null);
  const [conflito, setConflito] = useState<readonly AlteracaoNaoSalva[] | null>(null);
  const formulario = useFormulario(esquema, {
    defaultValues: {
      dia_semana: String(turno?.dia_semana ?? 1),
      nome: turno?.nome ?? '',
      hora_inicio: turno?.hora_inicio ?? '',
      hora_fim: turno?.hora_fim ?? '',
    },
  });
  const { errors } = formulario.formState;
  const pendente = incluir.isPending || alterar.isPending;

  const enviar = formulario.handleSubmit((v) => {
    setMensagem(null);
    setConflito(null);
    const tratar = (erro: unknown) =>
      tratarErroDeEnvio(erro, {
        setError: formulario.setError,
        definirMensagem: setMensagem,
        aoConflitoDeVersao: setConflito,
        alteracoes: () =>
          alteracoesNaoSalvas(v, formulario.formState.dirtyFields, {
            dia_semana: t('turnos.dia'),
            nome: t('turnos.nome'),
            hora_inicio: t('turnos.inicio'),
            hora_fim: t('turnos.fim'),
          }),
      });
    if (turno) {
      alterar.mutate(
        {
          unidadeId,
          turnoId: turno.id,
          versao: turno.versao,
          corpo: { dia_semana: Number(v.dia_semana), nome: vazioParaNulo(v.nome), hora_inicio: v.hora_inicio, hora_fim: v.hora_fim },
        },
        { onSuccess: aoConcluir, onError: tratar },
      );
      return;
    }
    incluir.mutate(
      {
        unidadeId,
        corpo: {
          dia_semana: Number(v.dia_semana),
          hora_inicio: v.hora_inicio,
          hora_fim: v.hora_fim,
          ...(v.nome ? { nome: v.nome } : {}),
        },
      },
      { onSuccess: aoConcluir, onError: tratar },
    );
  });

  return (
    <Formulario erros={errors} onSubmit={(e) => void enviar(e)} className="flex flex-col gap-4">
      {mensagem && <Alerta tipo="erro">{mensagem}</Alerta>}
      {conflito && <AlertaDeConflito alteracoes={conflito} />}
      <CampoSelecao rotulo={t('turnos.dia')} erro={errors.dia_semana} {...formulario.register('dia_semana')}>
        {NOMES_DOS_DIAS.map((nome, i) => (
          <option key={nome} value={i + 1}>
            {t(`turnos.dias.${nome}`)}
          </option>
        ))}
      </CampoSelecao>
      <CampoTexto rotulo={t('turnos.nome')} opcional autoComplete="off" erro={errors.nome} {...formulario.register('nome')} />
      <div className="grid grid-cols-2 gap-4">
        <CampoTexto rotulo={t('turnos.inicio')} placeholder={t('turnos.exemploInicio')} inputMode="numeric" dado autoComplete="off" erro={errors.hora_inicio} {...formulario.register('hora_inicio')} />
        <CampoTexto rotulo={t('turnos.fim')} placeholder={t('turnos.exemploFim')} inputMode="numeric" dado autoComplete="off" erro={errors.hora_fim} {...formulario.register('hora_fim')} />
      </div>
      <p className="texto-legenda text-texto-secundario">{t('turnos.fimNoDiaSeguinte')}</p>
      <div className="flex justify-end gap-2">
        <Botao onClick={aoConcluir}>{t('comum.voltar')}</Botao>
        <Botao type="submit" variante="primario" carregando={pendente}>
          {pendente ? t('comum.salvando') : t('comum.salvar')}
        </Botao>
      </div>
    </Formulario>
  );
}
