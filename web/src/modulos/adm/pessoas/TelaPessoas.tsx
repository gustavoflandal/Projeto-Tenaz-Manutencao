// Pessoas (RF-ADM-015): lista com a identificação mascarada, criação e detalhe
// com a identificação completa. Filtros na URL, para compartilhar e voltar; a lista é paginada por
// cursor no servidor, nunca filtrada no navegador.

import { Plus } from 'lucide-react';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Link, useNavigate, useParams, useSearchParams } from 'react-router';

import { Alerta } from '../../../componentes/Alerta';
import { AlertaDeConflito } from '../../../componentes/AlertaDeConflito';
import { AvisoDadoDeSaude } from '../../../componentes/AvisoDadoDeSaude';
import { Botao } from '../../../componentes/Botao';
import { CabecalhoPagina } from '../../../componentes/CabecalhoPagina';
import { BuscaPorTexto } from '../../../componentes/BuscaPorTexto';
import { CampoSelecao } from '../../../componentes/Campo';
import { ChipSituacao } from '../../../componentes/ChipStatus';
import { Confirmacao } from '../../../componentes/Dialogos';
import { Carregando, EstadoErro, EstadoVazio } from '../../../componentes/Estados';
import { RodapeDeLista, Tabela, Td, Th, Tr } from '../../../componentes/Tabela';
import { ehNaoEncontrado, mensagemDoErro, type AlteracaoNaoSalva } from '../../../lib/erros';
import { traduzirOuOriginal } from '../../../lib/i18n';
import { useTemPermissao } from '../sessao';
import { useInativarPessoa, usePessoa, usePessoas, useReativarPessoa, type FiltroDePessoas } from '../pessoas';

import { FormularioDePessoa } from './FormularioDePessoa';

/** Tipo e identificação mascarada (`DNI *****678`); nulo quando a pessoa não tem identificação. */
function documentoMascarado(tipo: string | null, mascarada: string | null): string | null {
  if (!mascarada) return null;
  return tipo ? `${traduzirOuOriginal('identificacao.tipos', tipo.toLowerCase())} ${mascarada}` : mascarada;
}

function filtroDaUrl(parametros: URLSearchParams): FiltroDePessoas {
  const busca = parametros.get('busca') ?? '';
  const situacao = parametros.get('situacao');
  const comUsuario = parametros.get('com_usuario');
  return {
    ...(busca ? { busca } : {}),
    ...(situacao === 'ativa' || situacao === 'inativa' ? { situacao } : {}),
    ...(comUsuario === 'true' || comUsuario === 'false' ? { com_usuario: comUsuario } : {}),
  };
}

export function TelaPessoas() {
  const { t } = useTranslation();
  const navegar = useNavigate();
  const [parametros, definirParametros] = useSearchParams();
  const filtro = filtroDaUrl(parametros);
  const lista = usePessoas(filtro);
  const podeCriar = useTemPermissao('adm.gerenciar_pessoas');
  const temFiltro = Object.keys(filtro).length > 0;

  const mudar = (chave: string, valor: string) => {
    const novos = new URLSearchParams(parametros);
    if (valor) novos.set(chave, valor);
    else novos.delete(chave);
    definirParametros(novos, { replace: true });
  };

  const acaoNova = podeCriar ? (
    <Botao variante="primario" onClick={() => void navegar('/administracao/pessoas/nova')}>
      <Plus size={20} aria-hidden="true" />
      {t('pessoa.nova')}
    </Botao>
  ) : undefined;

  return (
    <>
      <CabecalhoPagina titulo={t('pessoa.titulo')} acoes={acaoNova} />
      <div className="mb-4 flex flex-wrap items-end gap-3" role="search" aria-label={t('lista.filtros')}>
        <BuscaPorTexto rotulo={t('pessoa.buscar')} valor={parametros.get('busca') ?? ''} aoBuscar={(valor) => mudar('busca', valor)} />
        <div className="w-44">
          <CampoSelecao rotulo={t('comum.situacao')} value={parametros.get('situacao') ?? ''} onChange={(e) => mudar('situacao', e.target.value)} name="situacao">
            <option value="">{t('lista.todas')}</option>
            <option value="ativa">{t('situacao.ativa')}</option>
            <option value="inativa">{t('situacao.inativa')}</option>
          </CampoSelecao>
        </div>
        <div className="w-52">
          <CampoSelecao rotulo={t('pessoa.usuarioFiltro')} value={parametros.get('com_usuario') ?? ''} onChange={(e) => mudar('com_usuario', e.target.value)} name="com_usuario">
            <option value="">{t('lista.todas')}</option>
            <option value="true">{t('pessoa.comUsuario')}</option>
            <option value="false">{t('pessoa.semUsuario')}</option>
          </CampoSelecao>
        </div>
        {temFiltro && (
          <Botao variante="terciario" onClick={() => definirParametros({}, { replace: true })}>
            {t('lista.limparFiltros')}
          </Botao>
        )}
      </div>

      {lista.carregando && <Carregando rotulo={t('pessoa.carregando')} />}
      {lista.erro !== null && lista.erro !== undefined && <EstadoErro erro={lista.erro} aoTentarDeNovo={lista.recarregar} />}
      {!lista.carregando && !lista.erro && lista.itens.length === 0 && (
        <EstadoVazio
          titulo={temFiltro ? t('lista.semResultadoTitulo') : t('pessoa.vazioTitulo')}
          texto={temFiltro ? t('lista.semResultado') : t('pessoa.vazio')}
          acao={temFiltro ? undefined : acaoNova}
        />
      )}
      {lista.itens.length > 0 && (
        <>
          <Tabela legenda={t('pessoa.titulo')}>
            <thead>
              <tr>
                <Th>{t('pessoa.nome')}</Th>
                <Th>{t('identificacao.documento')}</Th>
                <Th>{t('pessoa.matricula')}</Th>
                <Th>{t('pessoa.cargo')}</Th>
                <Th>{t('pessoa.vinculo')}</Th>
                <Th>{t('pessoa.usuario')}</Th>
                <Th>{t('comum.situacao')}</Th>
              </tr>
            </thead>
            <tbody>
              {lista.itens.map((p) => (
                <Tr key={p.id}>
                  <Td>
                    <Link to={`/administracao/pessoas/${p.id}`} className="underline">
                      {p.nome}
                    </Link>
                  </Td>
                  <Td className="texto-dado">{documentoMascarado(p.identificacao_tipo, p.identificacao_mascarada) ?? t('comum.naoInformado')}</Td>
                  <Td className="texto-dado">{p.matricula ?? t('comum.naoInformado')}</Td>
                  <Td>{p.cargo ?? t('comum.naoInformado')}</Td>
                  <Td>{t(`vinculos.${p.vinculo}`)}</Td>
                  <Td>{p.usuario_id ? t('comum.sim') : t('comum.nao')}</Td>
                  <Td>
                    <ChipSituacao situacao={p.situacao} />
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

export function TelaPessoaNova() {
  const { t } = useTranslation();
  return (
    <>
      <CabecalhoPagina
        titulo={t('pessoa.nova')}
        trilha={[{ rotulo: t('pessoa.titulo'), para: '/administracao/pessoas' }, { rotulo: t('pessoa.nova') }]}
      />
      <FormularioDePessoa />
    </>
  );
}

export function TelaPessoaDetalhe() {
  const { t } = useTranslation();
  const { id = '' } = useParams();
  const consulta = usePessoa(id);
  const inativar = useInativarPessoa();
  const reativar = useReativarPessoa();
  const podeEditar = useTemPermissao('adm.gerenciar_pessoas');
  const podeVerAuditoria = useTemPermissao('adm.ler_auditoria');
  const podeVerUsuarios = useTemPermissao('adm.ler_usuarios');
  const [confirmando, setConfirmando] = useState(false);
  const [conflito, setConflito] = useState<readonly AlteracaoNaoSalva[] | null>(null);

  if (consulta.isPending) {
    return (
      <>
        <CabecalhoPagina titulo={t('pessoa.titulo')} />
        <div className="mb-4">
          <AvisoDadoDeSaude />
        </div>
        <Carregando rotulo={t('pessoa.carregando')} linhas={4} />
      </>
    );
  }
  if (consulta.isError) {
    if (ehNaoEncontrado(consulta.error)) {
      return <EstadoVazio titulo={t('pessoa.naoEncontradaTitulo')} texto={t('pessoa.naoEncontrada')} />;
    }
    return <EstadoErro erro={consulta.error} aoTentarDeNovo={() => void consulta.refetch()} />;
  }

  const pessoa = consulta.data;
  const ativa = pessoa.situacao === 'ativa';
  const acoes = (
    <>
      {podeVerAuditoria && (
        <Link
          to={`/administracao/auditoria?entidade=pessoa&entidade_id=${pessoa.id}`}
          className="inline-flex h-10 items-center px-2 texto-botao underline"
        >
          {t('auditoria.historicoDoRegistro')}
        </Link>
      )}
      {podeEditar &&
        (ativa ? (
          <Botao onClick={() => setConfirmando(true)}>{t('pessoa.inativar')}</Botao>
        ) : (
          <Botao
            variante="primario"
            carregando={reativar.isPending}
            onClick={() => reativar.mutate({ id: pessoa.id, versao: pessoa.versao })}
          >
            {t('pessoa.reativar')}
          </Botao>
        ))}
    </>
  );

  return (
    <>
      <CabecalhoPagina
        titulo={pessoa.nome}
        trilha={[{ rotulo: t('pessoa.titulo'), para: '/administracao/pessoas' }, { rotulo: pessoa.nome }]}
        metadados={
          <>
            <ChipSituacao situacao={pessoa.situacao} />
            {pessoa.usuario_id && podeVerUsuarios ? (
              <Link to={`/administracao/usuarios/${pessoa.usuario_id}`} className="underline">
                {t('pessoa.verUsuario')}
              </Link>
            ) : (
              <span>{pessoa.usuario_id ? t('pessoa.temUsuario') : t('pessoa.semUsuarioTexto')}</span>
            )}
          </>
        }
        acoes={acoes}
      />
      {reativar.isError && <Alerta tipo="erro" className="mb-4">{mensagemDoErro(reativar.error)}</Alerta>}
      {conflito && <AlertaDeConflito alteracoes={conflito} className="mb-4" />}
      {!ativa && <Alerta tipo="informacao" className="mb-4">{t('pessoa.inativaTexto')}</Alerta>}
      <FormularioDePessoa
        pessoa={pessoa}
        podeEditar={podeEditar && ativa}
        aoConflitoDeVersao={(alteracoes) => {
          setConflito(alteracoes);
          void consulta.refetch();
        }}
      />

      <Confirmacao
        aberto={confirmando}
        aoMudar={(aberto) => {
          setConfirmando(aberto);
          if (!aberto) inativar.reset();
        }}
        titulo={t('pessoa.inativarTitulo', { nome: pessoa.nome })}
        texto={t('pessoa.inativarTexto')}
        rotuloConfirmar={t('pessoa.inativarBotao')}
        carregando={inativar.isPending}
        erro={inativar.isError ? <Alerta tipo="erro">{mensagemDoErro(inativar.error)}</Alerta> : undefined}
        aoConfirmar={() => inativar.mutate({ id: pessoa.id, versao: pessoa.versao }, { onSuccess: () => setConfirmando(false) })}
      />
    </>
  );
}
