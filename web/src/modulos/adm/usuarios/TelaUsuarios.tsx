// Usuários (RF-ADM-004): lista com filtros na URL e paginação por cursor, e a
// tela de convite.

import { Plus } from 'lucide-react';
import { useTranslation } from 'react-i18next';
import { Link, useNavigate, useSearchParams } from 'react-router';

import { Botao } from '../../../componentes/Botao';
import { BuscaPorTexto } from '../../../componentes/BuscaPorTexto';
import { CabecalhoPagina } from '../../../componentes/CabecalhoPagina';
import { CampoSelecao } from '../../../componentes/Campo';
import { ChipSituacao } from '../../../componentes/ChipStatus';
import { Carregando, EstadoErro, EstadoVazio } from '../../../componentes/Estados';
import { RodapeDeLista, Tabela, Td, Th, Tr } from '../../../componentes/Tabela';
import { listar } from '../../../lib/formatacao';
import { useTemPermissao } from '../sessao';
import { usePerfis, useUsuarios, type FiltroDeUsuarios } from '../usuarios';

import { FormularioDeUsuario } from './FormularioDeUsuario';

function filtroDaUrl(parametros: URLSearchParams): FiltroDeUsuarios {
  const busca = parametros.get('busca') ?? '';
  const situacao = parametros.get('situacao');
  const perfil = parametros.get('perfil') ?? '';
  return {
    ...(busca ? { busca } : {}),
    ...(situacao === 'convidado' || situacao === 'ativo' || situacao === 'inativo' ? { situacao } : {}),
    ...(perfil ? { perfil } : {}),
  };
}

export function TelaUsuarios() {
  const { t } = useTranslation();
  const navegar = useNavigate();
  const [parametros, definirParametros] = useSearchParams();
  const filtro = filtroDaUrl(parametros);
  const lista = useUsuarios(filtro);
  const perfis = usePerfis();
  const podeConvidar = useTemPermissao('adm.gerenciar_usuarios');
  const temFiltro = Object.keys(filtro).length > 0;

  const mudar = (chave: string, valor: string) => {
    const novos = new URLSearchParams(parametros);
    if (valor) novos.set(chave, valor);
    else novos.delete(chave);
    definirParametros(novos, { replace: true });
  };

  const acaoNovo = podeConvidar ? (
    <Botao variante="primario" onClick={() => void navegar('/administracao/usuarios/novo')}>
      <Plus size={20} aria-hidden="true" />
      {t('usuario.convidar')}
    </Botao>
  ) : undefined;

  const nomeDoPerfil = (codigo: string) => perfis.data?.dados?.find((p) => p.codigo === codigo)?.nome ?? codigo;

  return (
    <>
      <CabecalhoPagina titulo={t('usuario.titulo')} acoes={acaoNovo} />
      <div className="mb-4 flex flex-wrap items-end gap-3" role="group" aria-label={t('lista.filtros')}>
        <BuscaPorTexto rotulo={t('usuario.buscar')} valor={parametros.get('busca') ?? ''} aoBuscar={(valor) => mudar('busca', valor)} />
        <div className="w-44">
          <CampoSelecao rotulo={t('comum.situacao')} name="situacao" value={parametros.get('situacao') ?? ''} onChange={(e) => mudar('situacao', e.target.value)}>
            <option value="">{t('lista.todas')}</option>
            <option value="convidado">{t('situacao.convidado')}</option>
            <option value="ativo">{t('situacao.ativo')}</option>
            <option value="inativo">{t('situacao.inativo')}</option>
          </CampoSelecao>
        </div>
        <div className="w-52">
          <CampoSelecao rotulo={t('usuario.perfil')} name="perfil" value={parametros.get('perfil') ?? ''} onChange={(e) => mudar('perfil', e.target.value)}>
            <option value="">{t('lista.todos')}</option>
            {(perfis.data?.dados ?? []).map((p) => (
              <option key={p.codigo} value={p.codigo}>
                {p.nome}
              </option>
            ))}
          </CampoSelecao>
        </div>
        {temFiltro && (
          <Botao variante="terciario" onClick={() => definirParametros({}, { replace: true })}>
            {t('lista.limparFiltros')}
          </Botao>
        )}
      </div>

      {lista.carregando && <Carregando rotulo={t('usuario.carregando')} />}
      {lista.erro !== null && lista.erro !== undefined && <EstadoErro erro={lista.erro} aoTentarDeNovo={lista.recarregar} />}
      {!lista.carregando && !lista.erro && lista.itens.length === 0 && (
        <EstadoVazio
          titulo={temFiltro ? t('lista.semResultadoTitulo') : t('usuario.vazioTitulo')}
          texto={temFiltro ? t('lista.semResultado') : t('usuario.vazio')}
          acao={temFiltro ? undefined : acaoNovo}
        />
      )}
      {lista.itens.length > 0 && (
        <>
          <Tabela legenda={t('usuario.titulo')}>
            <thead>
              <tr>
                <Th>{t('usuario.nome')}</Th>
                <Th>{t('usuario.email')}</Th>
                <Th>{t('usuario.perfis')}</Th>
                <Th>{t('usuario.pessoa')}</Th>
                <Th>{t('comum.situacao')}</Th>
              </tr>
            </thead>
            <tbody>
              {lista.itens.map((u) => (
                <Tr key={u.id}>
                  <Td>
                    <Link to={`/administracao/usuarios/${u.id}`} className="underline">
                      {u.nome}
                    </Link>
                  </Td>
                  <Td>{u.email ?? t('comum.naoInformado')}</Td>
                  <Td>{listar((u.perfis ?? []).map(nomeDoPerfil)) || t('comum.naoInformado')}</Td>
                  <Td>{u.pessoa_id ? t('comum.sim') : t('comum.nao')}</Td>
                  <Td>
                    <ChipSituacao situacao={u.situacao} />
                  </Td>
                </Tr>
              ))}
            </tbody>
          </Tabela>
          <RodapeDeLista quantidade={lista.itens.length} temMais={lista.temMais} carregandoMais={lista.carregandoMais} aoCarregarMais={lista.carregarMais} />
        </>
      )}
    </>
  );
}

export function TelaUsuarioNovo() {
  const { t } = useTranslation();
  return (
    <>
      <CabecalhoPagina
        titulo={t('usuario.convidar')}
        trilha={[{ rotulo: t('usuario.titulo'), para: '/administracao/usuarios' }, { rotulo: t('usuario.convidar') }]}
      />
      <FormularioDeUsuario />
    </>
  );
}
