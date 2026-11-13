// Detalhe da unidade: dados, turnos e feriados (RF-ADM-002, RN-ADM-019).
// Turno e feriado têm versão própria: usam o ETag deles, não o da unidade.

import { Pencil, Plus, Trash2 } from 'lucide-react';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useParams } from 'react-router';

import { Alerta } from '../../../componentes/Alerta';
import { AlertaDeConflito } from '../../../componentes/AlertaDeConflito';
import { Botao } from '../../../componentes/Botao';
import { CabecalhoPagina } from '../../../componentes/CabecalhoPagina';
import { ChipSituacao } from '../../../componentes/ChipStatus';
import { Confirmacao } from '../../../componentes/Dialogos';
import { Carregando, EstadoErro, EstadoVazio } from '../../../componentes/Estados';
import { Tabela, Td, Th, Tr } from '../../../componentes/Tabela';
import { ehNaoEncontrado, mensagemDoErro, type AlteracaoNaoSalva } from '../../../lib/erros';
import type { Feriado, Turno, Unidade } from '../../../lib/tipos';
import { useFormatos } from '../formatos';
import { useTemPermissao } from '../sessao';
import {
  useFeriados,
  useInativarUnidade,
  useReativarUnidade,
  useRemoverFeriado,
  useRemoverTurno,
  useUnidade,
} from '../unidades';

import { FormularioDeFeriado } from './FormularioDeFeriado';
import { FormularioDeTurno } from './FormularioDeTurno';
import { FormularioDeUnidade } from './FormularioDeUnidade';
import { Modal } from '../../../componentes/Dialogos';

const NOMES_DOS_DIAS = ['segunda', 'terca', 'quarta', 'quinta', 'sexta', 'sabado', 'domingo'] as const;

function ordenarTurnos(turnos: readonly Turno[]): Turno[] {
  return [...turnos].sort((a, b) => a.dia_semana - b.dia_semana || a.hora_inicio.localeCompare(b.hora_inicio));
}

function TurnosDaUnidade({ unidade, podeEditar }: { readonly unidade: Unidade; readonly podeEditar: boolean }) {
  const { t } = useTranslation();
  const remover = useRemoverTurno();
  const [editando, setEditando] = useState<Turno | 'novo' | null>(null);
  const [removendo, setRemovendo] = useState<Turno | null>(null);
  const turnos = ordenarTurnos(unidade.turnos ?? []);
  const nomeDoDia = (dia: number) => t(`turnos.dias.${NOMES_DOS_DIAS[dia - 1] ?? 'segunda'}`);

  return (
    <div className="flex flex-col gap-3">
      {podeEditar && (
        <div>
          <Botao onClick={() => setEditando('novo')}>
            <Plus size={16} aria-hidden="true" />
            {t('turnos.adicionar')}
          </Botao>
        </div>
      )}
      {turnos.length === 0 ? (
        <EstadoVazio titulo={t('turnos.vazioTitulo')} texto={t('turnos.vazio')} />
      ) : (
        <Tabela legenda={t('turnos.titulo')}>
          <thead>
            <tr>
              <Th>{t('turnos.dia')}</Th>
              <Th>{t('turnos.nome')}</Th>
              <Th>{t('turnos.horario')}</Th>
              {podeEditar && (
                <Th>
                  <span className="sr-only">{t('comum.acoes')}</span>
                </Th>
              )}
            </tr>
          </thead>
          <tbody>
            {turnos.map((turno) => (
              <Tr key={turno.id}>
                <Td>{nomeDoDia(turno.dia_semana)}</Td>
                <Td>{turno.nome ?? t('comum.naoInformado')}</Td>
                <Td className="texto-dado">{t('turnos.intervalo', { inicio: turno.hora_inicio, fim: turno.hora_fim })}</Td>
                {podeEditar && (
                  <Td className="w-28">
                    <div className="flex">
                      <Botao
                        variante="icone"
                        aria-label={t('turnos.editarDe', { dia: nomeDoDia(turno.dia_semana), horario: t('turnos.intervalo', { inicio: turno.hora_inicio, fim: turno.hora_fim }) })}
                        onClick={() => setEditando(turno)}
                      >
                        <Pencil size={20} aria-hidden="true" />
                      </Botao>
                      <Botao
                        variante="icone"
                        aria-label={t('turnos.removerDe', { dia: nomeDoDia(turno.dia_semana), horario: t('turnos.intervalo', { inicio: turno.hora_inicio, fim: turno.hora_fim }) })}
                        onClick={() => setRemovendo(turno)}
                      >
                        <Trash2 size={20} aria-hidden="true" />
                      </Botao>
                    </div>
                  </Td>
                )}
              </Tr>
            ))}
          </tbody>
        </Tabela>
      )}

      <Modal
        aberto={editando !== null}
        aoMudar={(aberto) => !aberto && setEditando(null)}
        titulo={editando === 'novo' ? t('turnos.novoTitulo') : t('turnos.editarTitulo')}
      >
        {editando !== null && (
          <FormularioDeTurno
            unidadeId={unidade.id}
            turno={editando === 'novo' ? undefined : editando}
            aoConcluir={() => setEditando(null)}
          />
        )}
      </Modal>

      <Confirmacao
        aberto={removendo !== null}
        aoMudar={(aberto) => {
          if (!aberto) {
            setRemovendo(null);
            remover.reset();
          }
        }}
        titulo={t('turnos.removerTitulo')}
        texto={t('turnos.removerTexto')}
        rotuloConfirmar={t('turnos.removerBotao')}
        perigo
        carregando={remover.isPending}
        erro={remover.isError ? <Alerta tipo="erro">{mensagemDoErro(remover.error)}</Alerta> : undefined}
        aoConfirmar={() =>
          removendo &&
          remover.mutate(
            { unidadeId: unidade.id, turnoId: removendo.id, versao: removendo.versao },
            { onSuccess: () => setRemovendo(null) },
          )
        }
      />
    </div>
  );
}

function FeriadosDaUnidade({ unidade, podeEditar }: { readonly unidade: Unidade; readonly podeEditar: boolean }) {
  const { t } = useTranslation();
  const { data: formatarData } = useFormatos();
  const consulta = useFeriados(unidade.id);
  const remover = useRemoverFeriado();
  const [incluindo, setIncluindo] = useState(false);
  const [removendo, setRemovendo] = useState<Feriado | null>(null);

  if (consulta.isPending) return <Carregando rotulo={t('feriados.carregando')} linhas={2} />;
  if (consulta.isError) return <EstadoErro erro={consulta.error} aoTentarDeNovo={() => void consulta.refetch()} />;
  const feriados = [...(consulta.data.dados ?? [])].sort((a, b) => a.data.localeCompare(b.data));

  return (
    <div className="flex flex-col gap-3">
      {podeEditar && (
        <div>
          <Botao onClick={() => setIncluindo(true)}>
            <Plus size={16} aria-hidden="true" />
            {t('feriados.adicionar')}
          </Botao>
        </div>
      )}
      {feriados.length === 0 ? (
        <EstadoVazio titulo={t('feriados.vazioTitulo')} texto={t('feriados.vazio')} />
      ) : (
        <Tabela legenda={t('feriados.titulo')}>
          <thead>
            <tr>
              <Th>{t('feriados.data')}</Th>
              <Th>{t('feriados.descricao')}</Th>
              {podeEditar && (
                <Th>
                  <span className="sr-only">{t('comum.acoes')}</span>
                </Th>
              )}
            </tr>
          </thead>
          <tbody>
            {feriados.map((feriado) => (
              <Tr key={feriado.id}>
                <Td className="texto-dado">{formatarData(feriado.data)}</Td>
                <Td>{feriado.descricao}</Td>
                {podeEditar && (
                  <Td className="w-16">
                    <Botao
                      variante="icone"
                      aria-label={t('feriados.removerDe', { descricao: feriado.descricao, data: formatarData(feriado.data) })}
                      onClick={() => setRemovendo(feriado)}
                    >
                      <Trash2 size={20} aria-hidden="true" />
                    </Botao>
                  </Td>
                )}
              </Tr>
            ))}
          </tbody>
        </Tabela>
      )}

      <Modal aberto={incluindo} aoMudar={setIncluindo} titulo={t('feriados.novoTitulo')} descricao={t('feriados.novoDescricao')}>
        <FormularioDeFeriado unidadeId={unidade.id} aoConcluir={() => setIncluindo(false)} />
      </Modal>

      <Confirmacao
        aberto={removendo !== null}
        aoMudar={(aberto) => {
          if (!aberto) {
            setRemovendo(null);
            remover.reset();
          }
        }}
        titulo={t('feriados.removerTitulo', { descricao: removendo?.descricao ?? '' })}
        texto={t('feriados.removerTexto')}
        rotuloConfirmar={t('feriados.removerBotao')}
        perigo
        carregando={remover.isPending}
        erro={remover.isError ? <Alerta tipo="erro">{mensagemDoErro(remover.error)}</Alerta> : undefined}
        aoConfirmar={() =>
          removendo &&
          remover.mutate(
            { unidadeId: unidade.id, feriadoId: removendo.id, versao: removendo.versao },
            { onSuccess: () => setRemovendo(null) },
          )
        }
      />
    </div>
  );
}

export function TelaUnidadeDetalhe() {
  const { t } = useTranslation();
  const { id = '' } = useParams();
  const consulta = useUnidade(id);
  const inativar = useInativarUnidade();
  const reativar = useReativarUnidade();
  const podeEditar = useTemPermissao('adm.gerenciar_empresa');
  const [confirmando, setConfirmando] = useState(false);
  const [conflito, setConflito] = useState<readonly AlteracaoNaoSalva[] | null>(null);

  if (consulta.isPending) {
    return (
      <>
        <CabecalhoPagina titulo={t('unidade.titulo')} />
        <Carregando rotulo={t('unidade.carregando')} linhas={4} />
      </>
    );
  }
  if (consulta.isError) {
    if (ehNaoEncontrado(consulta.error)) {
      return <EstadoVazio titulo={t('unidade.naoEncontradaTitulo')} texto={t('unidade.naoEncontrada')} />;
    }
    return <EstadoErro erro={consulta.error} aoTentarDeNovo={() => void consulta.refetch()} />;
  }

  const unidade = consulta.data;
  const ativa = unidade.situacao === 'ativa';
  const acoes = podeEditar ? (
    ativa ? (
      <Botao onClick={() => setConfirmando(true)}>{t('unidade.inativar')}</Botao>
    ) : (
      <Botao
        variante="primario"
        carregando={reativar.isPending}
        onClick={() => reativar.mutate({ id: unidade.id, versao: unidade.versao })}
      >
        {t('unidade.reativar')}
      </Botao>
    )
  ) : undefined;

  return (
    <>
      <CabecalhoPagina
        titulo={unidade.nome}
        trilha={[{ rotulo: t('unidade.titulo'), para: '/administracao/unidades' }, { rotulo: unidade.codigo }]}
        metadados={<ChipSituacao situacao={unidade.situacao} />}
        acoes={acoes}
      />
      {reativar.isError && <Alerta tipo="erro" className="mb-4">{mensagemDoErro(reativar.error)}</Alerta>}
      {conflito && <AlertaDeConflito alteracoes={conflito} className="mb-4" />}

      <div className="flex flex-col gap-10">
        <section aria-labelledby="unidade-dados" className="flex flex-col gap-4">
          <h2 id="unidade-dados" className="texto-h3">
            {t('unidade.dados')}
          </h2>
          <FormularioDeUnidade
            unidade={unidade}
            podeEditar={podeEditar && ativa}
            aoConflitoDeVersao={(alteracoes) => {
              setConflito(alteracoes);
              void consulta.refetch();
            }}
          />
        </section>
        <section aria-labelledby="unidade-turnos" className="flex flex-col gap-4">
          <h2 id="unidade-turnos" className="texto-h3">
            {t('turnos.titulo')}
          </h2>
          <TurnosDaUnidade unidade={unidade} podeEditar={podeEditar && ativa} />
        </section>
        <section aria-labelledby="unidade-feriados" className="flex flex-col gap-4">
          <h2 id="unidade-feriados" className="texto-h3">
            {t('feriados.titulo')}
          </h2>
          <FeriadosDaUnidade unidade={unidade} podeEditar={podeEditar && ativa} />
        </section>
      </div>

      <Confirmacao
        aberto={confirmando}
        aoMudar={(aberto) => {
          setConfirmando(aberto);
          if (!aberto) inativar.reset();
        }}
        titulo={t('unidade.inativarTitulo', { nome: unidade.nome })}
        texto={t('unidade.inativarTexto')}
        rotuloConfirmar={t('unidade.inativarBotao')}
        carregando={inativar.isPending}
        erro={inativar.isError ? <Alerta tipo="erro">{mensagemDoErro(inativar.error)}</Alerta> : undefined}
        aoConfirmar={() => inativar.mutate({ id: unidade.id, versao: unidade.versao }, { onSuccess: () => setConfirmando(false) })}
      />
    </>
  );
}
