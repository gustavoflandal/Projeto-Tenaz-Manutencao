// Trilha de auditoria (RF-ADM-011): consulta com filtros por usuário, tela,
// módulo, tipo de operação, período, IP, entidade e registro, paginada por
// cursor; exportação em CSV do período; detalhe do evento com antes, depois e hash.

import { Download, ShieldAlert, X } from 'lucide-react';
import { useState, type FormEvent } from 'react';
import { useTranslation } from 'react-i18next';
import { useSearchParams } from 'react-router';

import { Alerta } from '../../../componentes/Alerta';
import { Botao } from '../../../componentes/Botao';
import { CabecalhoPagina } from '../../../componentes/CabecalhoPagina';
import { CampoSelecao, CampoTexto, type ErroDeCampo } from '../../../componentes/Campo';
import { Modal } from '../../../componentes/Dialogos';
import { Carregando, EstadoErro, EstadoVazio } from '../../../componentes/Estados';
import { RodapeDeLista, Tabela, Td, Th, Tr } from '../../../componentes/Tabela';
import { baixarArquivo } from '../../../lib/download';
import { mensagemDoErro } from '../../../lib/erros';
import { dataParaISO, listar, mascararData } from '../../../lib/formatacao';
import { traduzirChave, traduzirOuOriginal } from '../../../lib/i18n';
import { MODULOS_DO_CATALOGO, TELAS_DO_CATALOGO } from '../../../lib/telas';
import type { ResumoEvento } from '../../../lib/tipos';
import { exportarEventos, useEventos, type FiltroDeEventos } from '../auditoria';
import { ACOES_DE_AUDITORIA } from '../constantes';
import { useDataHora, useFormatos } from '../formatos';

import { DetalheDoEvento } from './DetalheDoEvento';

/** Início do dia (local) como instante RFC 3339. */
function inicioDoDia(isoDaData: string, deslocamentoEmDias = 0): string {
  const [ano, mes, dia] = isoDaData.split('-').map(Number);
  return new Date(ano ?? 0, (mes ?? 1) - 1, (dia ?? 1) + deslocamentoEmDias).toISOString();
}

/** Filtro da API a partir da URL. O fim do período é inclusivo na tela e exclusivo na API. */
export function filtroDaUrl(parametros: URLSearchParams): FiltroDeEventos {
  const texto = (chave: string) => parametros.get(chave) ?? '';
  const de = texto('de');
  const ate = texto('ate');
  return {
    ...(texto('entidade') ? { entidade: texto('entidade') } : {}),
    ...(texto('entidade_id') ? { entidade_id: texto('entidade_id') } : {}),
    ...(texto('usuario_id') ? { usuario_id: texto('usuario_id') } : {}),
    ...(texto('acao') ? { acao: texto('acao') } : {}),
    ...(texto('modulo') ? { modulo: texto('modulo') } : {}),
    ...(texto('tela') ? { tela: texto('tela') } : {}),
    ...(texto('ip') ? { ip: texto('ip') } : {}),
    ...(de ? { de: inicioDoDia(de) } : {}),
    ...(ate ? { ate: inicioDoDia(ate, 1) } : {}),
  };
}

interface ErrosDoFiltro {
  de?: string;
  ate?: string;
  ip?: string;
  periodo?: string;
}

function FormularioDeFiltros({
  parametros,
  aoAplicar,
  aoLimpar,
  temFiltro,
}: {
  readonly parametros: URLSearchParams;
  readonly aoAplicar: (novos: URLSearchParams) => void;
  readonly aoLimpar: () => void;
  readonly temFiltro: boolean;
}) {
  const { t } = useTranslation();
  const { data: dataDaEmpresa } = useFormatos();
  const dataDaUrl = (iso: string | null): string => (iso ? dataDaEmpresa(iso) : '');
  const [de, setDe] = useState(dataDaUrl(parametros.get('de')));
  const [ate, setAte] = useState(dataDaUrl(parametros.get('ate')));
  const [erros, setErros] = useState<ErrosDoFiltro>({});

  const aplicar = (evento: FormEvent<HTMLFormElement>) => {
    evento.preventDefault();
    const dados = new FormData(evento.currentTarget);
    const valor = (nome: string) => String(dados.get(nome) ?? '').trim();
    const deISO = de.trim() === '' ? '' : dataParaISO(de);
    const ateISO = ate.trim() === '' ? '' : dataParaISO(ate);
    const novos: ErrosDoFiltro = {};
    if (deISO === null) novos.de = traduzirChave('validacao.dataInvalida');
    if (ateISO === null) novos.ate = traduzirChave('validacao.dataInvalida');
    if (deISO && ateISO && deISO > ateISO) novos.periodo = t('auditoria.periodoInvertido');
    if (valor('ip') && !deISO) novos.ip = t('auditoria.ipExigeInicio');
    setErros(novos);
    if (Object.keys(novos).length > 0) return;

    const url = new URLSearchParams();
    // Registro e usuário só entram por link (histórico do registro): ficam como estão.
    for (const fixo of ['entidade_id', 'usuario_id']) {
      const atual = parametros.get(fixo);
      if (atual) url.set(fixo, atual);
    }
    for (const chave of ['entidade', 'modulo', 'tela', 'acao', 'ip']) {
      if (valor(chave)) url.set(chave, valor(chave));
    }
    if (deISO) url.set('de', deISO);
    if (ateISO) url.set('ate', ateISO);
    aoAplicar(url);
  };

  const erro = (mensagem: string | undefined): ErroDeCampo | undefined =>
    mensagem ? { type: 'servidor', message: mensagem } : undefined;

  return (
    <form noValidate role="search" aria-label={t('lista.filtros')} onSubmit={aplicar} className="mb-4 flex flex-col gap-3">
      <div className="flex flex-wrap items-start gap-3">
        <div className="w-40">
          <CampoTexto
            rotulo={t('auditoria.entidade')}
            name="entidade"
            placeholder={t('auditoria.entidadeExemplo')}
            defaultValue={parametros.get('entidade') ?? ''}
          />
        </div>
        <div className="w-44">
          <CampoSelecao rotulo={t('auditoria.modulo')} name="modulo" defaultValue={parametros.get('modulo') ?? ''}>
            <option value="">{t('auditoria.todosModulos')}</option>
            {MODULOS_DO_CATALOGO.map((modulo) => (
              <option key={modulo} value={modulo}>
                {modulo}
              </option>
            ))}
          </CampoSelecao>
        </div>
        <div className="w-56">
          <CampoSelecao rotulo={t('auditoria.tela')} name="tela" defaultValue={parametros.get('tela') ?? ''}>
            <option value="">{t('auditoria.todasTelas')}</option>
            {TELAS_DO_CATALOGO.map((tela) => (
              <option key={tela} value={tela}>
                {tela}
              </option>
            ))}
          </CampoSelecao>
        </div>
        <div className="w-52">
          <CampoSelecao rotulo={t('auditoria.acao')} name="acao" defaultValue={parametros.get('acao') ?? ''}>
            <option value="">{t('lista.todas')}</option>
            {ACOES_DE_AUDITORIA.map((acao) => (
              <option key={acao} value={acao}>
                {t(`auditoria.acoes.${acao}`)}
              </option>
            ))}
          </CampoSelecao>
        </div>
      </div>
      <div className="flex flex-wrap items-start gap-3">
        <div className="w-36">
          <CampoTexto
            rotulo={t('auditoria.de')}
            name="de"
            dado
            inputMode="numeric"
            placeholder={t('feriados.dataExemplo')}
            value={de}
            onChange={(e) => setDe(mascararData(e.target.value))}
            erro={erro(erros.de)}
          />
        </div>
        <div className="w-36">
          <CampoTexto
            rotulo={t('auditoria.ate')}
            name="ate"
            dado
            inputMode="numeric"
            placeholder={t('feriados.dataExemplo')}
            value={ate}
            onChange={(e) => setAte(mascararData(e.target.value))}
            erro={erro(erros.ate ?? erros.periodo)}
          />
        </div>
        <div className="w-48">
          <CampoTexto
            rotulo={t('auditoria.ip')}
            name="ip"
            dado
            autoComplete="off"
            apoio={t('auditoria.ipApoio')}
            defaultValue={parametros.get('ip') ?? ''}
            erro={erro(erros.ip)}
          />
        </div>
        <div className="flex items-end gap-2 pt-7">
          <Botao type="submit" variante="primario">
            {t('lista.buscar')}
          </Botao>
          {temFiltro && (
            <Botao variante="terciario" onClick={aoLimpar}>
              {t('lista.limparFiltros')}
            </Botao>
          )}
        </div>
      </div>
    </form>
  );
}

export function TelaAuditoria() {
  const { t } = useTranslation();
  const dataHora = useDataHora();
  const [parametros, definirParametros] = useSearchParams();
  const filtro = filtroDaUrl(parametros);
  const lista = useEventos(filtro);
  const [aberto, setAberto] = useState<string | null>(null);
  const [exportando, setExportando] = useState(false);
  const [aviso, setAviso] = useState<{ tipo: 'sucesso' | 'erro'; texto: string } | null>(null);
  const temFiltro = [...parametros.keys()].length > 0;

  const removerDaUrl = (chave: string) => {
    const novos = new URLSearchParams(parametros);
    novos.delete(chave);
    definirParametros(novos, { replace: true });
  };

  const exportar = async () => {
    setAviso(null);
    if (filtro.de === undefined || filtro.ate === undefined) {
      setAviso({ tipo: 'erro', texto: t('auditoria.exportarExigePeriodo') });
      return;
    }
    setExportando(true);
    try {
      baixarArquivo(await exportarEventos(filtro), `trilha-auditoria-${parametros.get('de')}-${parametros.get('ate')}.csv`);
      setAviso({ tipo: 'sucesso', texto: t('auditoria.exportado') });
    } catch (erro) {
      setAviso({ tipo: 'erro', texto: mensagemDoErro(erro) });
    } finally {
      setExportando(false);
    }
  };

  const registro = parametros.get('entidade_id');
  const usuarioId = parametros.get('usuario_id');
  return (
    <>
      <CabecalhoPagina
        titulo={t('auditoria.titulo')}
        acoes={
          <Botao carregando={exportando} onClick={() => void exportar()}>
            <Download size={20} aria-hidden="true" />
            {exportando ? t('auditoria.exportando') : t('auditoria.exportar')}
          </Botao>
        }
      />
      <p className="mb-4 max-w-180 texto-corpo text-texto-secundario">{t('auditoria.texto')}</p>

      <FormularioDeFiltros
        key={parametros.toString()}
        parametros={parametros}
        aoAplicar={(novos) => definirParametros(novos, { replace: true })}
        aoLimpar={() => definirParametros({}, { replace: true })}
        temFiltro={temFiltro}
      />

      {aviso && (
        <Alerta tipo={aviso.tipo} className="mb-4">
          {aviso.texto}
        </Alerta>
      )}
      {(registro || usuarioId) && (
        <ul className="mb-4 flex flex-wrap gap-2" aria-label={t('auditoria.filtrosAtivos')}>
          {registro && (
            <li className="inline-flex items-center gap-1 rounded-sm bg-superficie-selecionada py-1 pr-1 pl-3 texto-legenda">
              {t('auditoria.filtroRegistro', { id: registro })}
              <Botao variante="icone" className="size-8" aria-label={t('auditoria.removerFiltroRegistro')} onClick={() => removerDaUrl('entidade_id')}>
                <X size={16} aria-hidden="true" />
              </Botao>
            </li>
          )}
          {usuarioId && (
            <li className="inline-flex items-center gap-1 rounded-sm bg-superficie-selecionada py-1 pr-1 pl-3 texto-legenda">
              {t('auditoria.filtroUsuario', { id: usuarioId })}
              <Botao variante="icone" className="size-8" aria-label={t('auditoria.removerFiltroUsuario')} onClick={() => removerDaUrl('usuario_id')}>
                <X size={16} aria-hidden="true" />
              </Botao>
            </li>
          )}
        </ul>
      )}

      {lista.carregando && <Carregando rotulo={t('auditoria.carregando')} />}
      {lista.erro !== null && lista.erro !== undefined && <EstadoErro erro={lista.erro} aoTentarDeNovo={lista.recarregar} />}
      {!lista.carregando && !lista.erro && lista.itens.length === 0 && (
        <EstadoVazio
          titulo={temFiltro ? t('lista.semResultadoTitulo') : t('auditoria.vazioTitulo')}
          texto={temFiltro ? t('lista.semResultado') : t('auditoria.vazio')}
        />
      )}
      {lista.itens.length > 0 && (
        <>
          <Tabela legenda={t('auditoria.titulo')}>
            <thead>
              <tr>
                <Th>{t('auditoria.quando')}</Th>
                <Th>{t('auditoria.usuario')}</Th>
                <Th>{t('auditoria.acao')}</Th>
                <Th>{t('auditoria.tela')}</Th>
                <Th>{t('auditoria.entidade')}</Th>
                <Th>{t('auditoria.ip')}</Th>
                <Th>
                  <span className="sr-only">{t('comum.acoes')}</span>
                </Th>
              </tr>
            </thead>
            <tbody>
              {lista.itens.map((e: ResumoEvento) => {
                const acao = traduzirOuOriginal('auditoria.acoes', e.acao);
                const quando = dataHora(e.ocorrido_em) ?? e.ocorrido_em;
                return (
                  <Tr key={e.id}>
                    <Td className="texto-dado whitespace-nowrap">{quando}</Td>
                    <Td>{e.usuario_nome || (e.usuario_id ? e.usuario_id.slice(0, 8) : t(`auditoria.atores.${e.ator_tipo}`))}</Td>
                    <Td>
                      {acao}
                      {e.resultado === 'falha' && <span className="ml-2 texto-legenda font-semibold">{t('auditoria.resultados.falha')}</span>}
                      {e.sensivel && (
                        <span className="ml-2 inline-flex align-middle">
                          <ShieldAlert size={14} aria-hidden="true" />
                          <span className="sr-only">{t('auditoria.sensivel')}</span>
                        </span>
                      )}
                    </Td>
                    <Td className="texto-dado">{e.tela ?? e.modulo ?? t('comum.naoInformado')}</Td>
                    <Td>{listar([e.entidade, e.entidade_chave])}</Td>
                    <Td className="texto-dado">{e.ip ?? t('comum.naoInformado')}</Td>
                    <Td>
                      <Botao variante="terciario" onClick={() => setAberto(e.id)} aria-label={t('auditoria.verDetalhesDe', { acao, quando })}>
                        {t('auditoria.verDetalhes')}
                      </Botao>
                    </Td>
                  </Tr>
                );
              })}
            </tbody>
          </Tabela>
          <RodapeDeLista quantidade={lista.itens.length} temMais={lista.temMais} carregandoMais={lista.carregandoMais} aoCarregarMais={lista.carregarMais} />
        </>
      )}

      <Modal aberto={aberto !== null} aoMudar={(a) => !a && setAberto(null)} titulo={t('auditoria.detalheTitulo')}>
        {aberto !== null && <DetalheDoEvento id={aberto} />}
      </Modal>
    </>
  );
}
