// Unidades (RF-ADM-002): lista, criação e detalhe com turnos e feriados.

import { Plus } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { Link, useNavigate } from 'react-router';

import { Botao } from '../../../componentes/Botao';
import { CabecalhoPagina } from '../../../componentes/CabecalhoPagina';
import { ChipSituacao } from '../../../componentes/ChipStatus';
import { Carregando, EstadoErro, EstadoVazio } from '../../../componentes/Estados';
import { Tabela, Td, Th, Tr } from '../../../componentes/Tabela';
import { listar } from '../../../lib/formatacao';
import { useTemPermissao } from '../sessao';
import { useUnidades } from '../unidades';

import { nomeDaSubdivisao } from './CampoDeSubdivisao';
import { FormularioDeUnidade } from './FormularioDeUnidade';

export function TelaUnidades() {
  const { t } = useTranslation();
  const navegar = useNavigate();
  const consulta = useUnidades();
  const podeCriar = useTemPermissao('adm.gerenciar_empresa');

  const acaoNova = podeCriar ? (
    <Botao variante="primario" onClick={() => void navegar('/administracao/unidades/nova')}>
      <Plus size={20} aria-hidden="true" />
      {t('unidade.nova')}
    </Botao>
  ) : undefined;

  if (consulta.isPending) {
    return (
      <>
        <CabecalhoPagina titulo={t('unidade.titulo')} />
        <Carregando rotulo={t('unidade.carregando')} />
      </>
    );
  }
  if (consulta.isError) return <EstadoErro erro={consulta.error} aoTentarDeNovo={() => void consulta.refetch()} />;

  const unidades = consulta.data.dados ?? [];
  return (
    <>
      <CabecalhoPagina titulo={t('unidade.titulo')} acoes={acaoNova} />
      {unidades.length === 0 ? (
        <EstadoVazio titulo={t('unidade.vazioTitulo')} texto={t('unidade.vazio')} acao={acaoNova} />
      ) : (
        <Tabela legenda={t('unidade.titulo')}>
          <thead>
            <tr>
              <Th>{t('unidade.codigo')}</Th>
              <Th>{t('unidade.nome')}</Th>
              <Th>{t('unidade.localizacao')}</Th>
              <Th>{t('comum.situacao')}</Th>
            </tr>
          </thead>
          <tbody>
            {unidades.map((u) => (
              <Tr key={u.id}>
                <Td className="texto-dado-forte">{u.codigo}</Td>
                <Td>
                  <Link to={`/administracao/unidades/${u.id}`} className="underline">
                    {u.nome}
                  </Link>
                </Td>
                <Td>{listar([u.endereco.municipio, u.endereco.subdivisao ? nomeDaSubdivisao(u.endereco.subdivisao) : undefined]) || t('comum.naoInformado')}</Td>
                <Td>
                  <ChipSituacao situacao={u.situacao} />
                </Td>
              </Tr>
            ))}
          </tbody>
        </Tabela>
      )}
    </>
  );
}

export function TelaUnidadeNova() {
  const { t } = useTranslation();
  return (
    <>
      <CabecalhoPagina
        titulo={t('unidade.nova')}
        trilha={[{ rotulo: t('unidade.titulo'), para: '/administracao/unidades' }, { rotulo: t('unidade.nova') }]}
      />
      <FormularioDeUnidade />
    </>
  );
}
