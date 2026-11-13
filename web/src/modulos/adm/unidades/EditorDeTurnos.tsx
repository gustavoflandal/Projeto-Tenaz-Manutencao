// Editor de turnos por dia da semana (RF-ADM-001 e RF-ADM-002, RN-ADM-019).
//
// O sistema não impõe calendário: o formulário nasce vazio e os botões de
// modelo ("Comercial" e "Contínuo 24x7") só preenchem os horários, que seguem
// editáveis. Sobreposição, horário inválido e "sem turno" aparecem no campo.

import { DIAS_DA_SEMANA, turnosDoModelo, type ModeloDeTurnos } from '@tenaz/dominio';
import { Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { useFieldArray, type Control, type FieldErrors, type FieldValues, type Path, type UseFormRegister } from 'react-hook-form';
import { useTranslation } from 'react-i18next';

import { Botao } from '../../../componentes/Botao';
import { CampoTexto, textoDoErro, type ErroDeCampo } from '../../../componentes/Campo';
import { Confirmacao } from '../../../componentes/Dialogos';

import { turnosDoFormulario, type TurnoDoFormulario } from './turnos';

interface ComTurnos extends FieldValues {
  turnos: TurnoDoFormulario[];
}

export interface EditorDeTurnosProps<T extends ComTurnos> {
  readonly control: Control<T>;
  readonly register: UseFormRegister<T>;
  readonly errors: FieldErrors<T>;
}

const NOMES_DOS_DIAS = ['segunda', 'terca', 'quarta', 'quinta', 'sexta', 'sabado', 'domingo'] as const;

function erroDeLinha(erros: unknown, indice: number): ErroDeCampo | undefined {
  const lista = erros as Array<ErroDeCampo | undefined> | undefined;
  return lista?.[indice];
}

export function EditorDeTurnos<T extends ComTurnos>({ control, register, errors }: EditorDeTurnosProps<T>) {
  const { t } = useTranslation();
  const { fields, append, remove, replace } = useFieldArray({ control, name: 'turnos' as never });

  const erros = (errors as FieldErrors<ComTurnos>).turnos;
  const erroGeral: ErroDeCampo | undefined =
    (erros as { root?: ErroDeCampo } | undefined)?.root ?? (erros as ErroDeCampo | undefined);
  const mensagemGeral = erroGeral?.message ? textoDoErro(erroGeral) : undefined;

  const [modeloPendente, setModeloPendente] = useState<ModeloDeTurnos | null>(null);
  const preencher = (modelo: ModeloDeTurnos) => replace(turnosDoFormulario(turnosDoModelo(modelo)) as never);
  // O modelo só preenche; se já há turnos digitados, pergunta antes de trocá-los.
  const pedirModelo = (modelo: ModeloDeTurnos) => (fields.length > 0 ? setModeloPendente(modelo) : preencher(modelo));

  const nomeDoCampo = (indice: number, campo: string) => `turnos.${indice}.${campo}` as Path<T>;

  return (
    <fieldset className="flex flex-col gap-4">
      <legend className="texto-h4">{t('turnos.titulo')}</legend>
      <p className="texto-corpo text-texto-secundario">{t('turnos.texto')}</p>

      <div className="flex flex-wrap items-center gap-2">
        <span className="texto-rotulo">{t('turnos.modelos')}</span>
        <Botao onClick={() => pedirModelo('comercial')}>{t('turnos.modeloComercial')}</Botao>
        <Botao onClick={() => pedirModelo('continuo')}>{t('turnos.modeloContinuo')}</Botao>
      </div>
      <p className="-mt-2 texto-legenda text-texto-secundario">{t('turnos.modelosApoio')}</p>

      {mensagemGeral && (
        <p role="alert" className="texto-corpo text-texto">
          {mensagemGeral}
        </p>
      )}

      <div className="flex flex-col gap-3">
        {DIAS_DA_SEMANA.map((dia) => {
          const linhas = fields
            .map((campo, indice) => ({ campo, indice }))
            .filter(({ campo }) => (campo as unknown as TurnoDoFormulario).dia_semana === dia);
          const nomeDoDia = t(`turnos.dias.${NOMES_DOS_DIAS[dia - 1] ?? 'segunda'}`);
          return (
            <fieldset key={dia} className="rounded-md border border-borda p-3">
              <legend className="px-1 texto-rotulo">{nomeDoDia}</legend>
              {linhas.length === 0 && <p className="texto-corpo-p text-texto-secundario">{t('turnos.semTurno')}</p>}
              <ul className="flex flex-col gap-3">
                {linhas.map(({ campo, indice }) => {
                  const erro = erroDeLinha(erros, indice);
                  const mensagem = erro?.message ? textoDoErro(erro) : undefined;
                  return (
                    <li key={campo.id} className="flex flex-col gap-2">
                      <div className="flex flex-wrap items-start gap-3">
                        <input type="hidden" {...register(nomeDoCampo(indice, 'dia_semana'), { valueAsNumber: true })} />
                        <div className="min-w-40 flex-2">
                        <CampoTexto
                          rotulo={t('turnos.nome')}
                          opcional
                          autoComplete="off"
                          erro={(erro as Record<string, ErroDeCampo | undefined> | undefined)?.nome}
                          {...register(nomeDoCampo(indice, 'nome'))}
                        />
                        </div>
                        <div className="min-w-28 flex-1">
                        <CampoTexto
                          rotulo={t('turnos.inicio')}
                          placeholder={t('turnos.exemploInicio')}
                          inputMode="numeric"
                          autoComplete="off"
                          dado
                          erro={(erro as Record<string, ErroDeCampo | undefined> | undefined)?.hora_inicio}
                          {...register(nomeDoCampo(indice, 'hora_inicio'))}
                        />
                        </div>
                        <div className="min-w-28 flex-1">
                        <CampoTexto
                          rotulo={t('turnos.fim')}
                          placeholder={t('turnos.exemploFim')}
                          inputMode="numeric"
                          autoComplete="off"
                          dado
                          erro={(erro as Record<string, ErroDeCampo | undefined> | undefined)?.hora_fim}
                          {...register(nomeDoCampo(indice, 'hora_fim'))}
                        />
                        </div>
                        <Botao
                          variante="icone"
                          aria-label={t('turnos.remover', { dia: nomeDoDia })}
                          onClick={() => remove(indice)}
                          className="mt-7"
                        >
                          <Trash2 size={20} aria-hidden="true" />
                        </Botao>
                      </div>
                      {mensagem && (
                        <p role="alert" className="texto-legenda text-texto">
                          {mensagem}
                        </p>
                      )}
                    </li>
                  );
                })}
              </ul>
              <div className="mt-3">
                <Botao
                  onClick={() => append({ dia_semana: dia, nome: '', hora_inicio: '', hora_fim: '' } as never)}
                  aria-label={t('turnos.adicionarEm', { dia: nomeDoDia })}
                >
                  <Plus size={16} aria-hidden="true" />
                  {t('turnos.adicionar')}
                </Botao>
              </div>
            </fieldset>
          );
        })}
      </div>
      <p className="texto-legenda text-texto-secundario">{t('turnos.fimNoDiaSeguinte')}</p>
      <Confirmacao
        aberto={modeloPendente !== null}
        aoMudar={(aberto) => !aberto && setModeloPendente(null)}
        titulo={t('turnos.substituirTitulo')}
        texto={t('turnos.substituirTexto', { modelo: modeloPendente ? t(modeloPendente === 'comercial' ? 'turnos.modeloComercial' : 'turnos.modeloContinuo') : '' })}
        rotuloConfirmar={t('turnos.substituirBotao')}
        aoConfirmar={() => {
          if (modeloPendente) preencher(modeloPendente);
          setModeloPendente(null);
        }}
      />
    </fieldset>
  );
}
