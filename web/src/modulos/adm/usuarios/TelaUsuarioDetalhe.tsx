// Detalhe do usuário (RF-ADM-004 e RF-ADM-005): dados e perfis, permissões
// avulsas e as ações de ciclo de vida (desativar, reativar, reenviar convite,
// encerrar sessões). O servidor decide cada uma; aqui o botão só aparece para
// quem tem a permissão.

import type { Moeda } from '@tenaz/dominio';
import { X } from 'lucide-react';
import { useMemo, useState } from 'react';
import { Controller } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { Link, useLocation, useParams } from 'react-router';

import { Alerta } from '../../../componentes/Alerta';
import { AlertaDeConflito } from '../../../componentes/AlertaDeConflito';
import { Botao } from '../../../componentes/Botao';
import { CabecalhoPagina } from '../../../componentes/CabecalhoPagina';
import { CampoSelecao, CampoTexto } from '../../../componentes/Campo';
import { ChipSituacao } from '../../../componentes/ChipStatus';
import { Confirmacao } from '../../../componentes/Dialogos';
import { Carregando, EstadoErro, EstadoVazio } from '../../../componentes/Estados';
import { Formulario } from '../../../componentes/Formulario';
import { alteracoesNaoSalvas, ehNaoEncontrado, mensagemDoErro, tratarErroDeEnvio, type AlteracaoNaoSalva } from '../../../lib/erros';
import { custoHoraOpcional, precoUnitarioParaOContrato, useFormulario, vazioParaNulo, z } from '../../../lib/formularios';
import type { Usuario } from '../../../lib/tipos';
import { useFormatos } from '../formatos';
import { useMoedaDeCusto, useRegiao } from '../regiao';
import { useSessao, useTemPermissao } from '../sessao';
import { useUnidades } from '../unidades';
import {
  useAlterarUsuario,
  useConcederPermissao,
  useDesativarUsuario,
  useEncerrarSessoes,
  usePerfis,
  usePermissoes,
  useReativarUsuario,
  useReenviarConvite,
  useRetirarPermissao,
  useUsuario,
} from '../usuarios';

import { SeletorDePerfis } from './SeletorDePerfis';

const esquemaDoFormulario = (moeda: Moeda | null) =>
  z.object({
    nome: z.string().trim().min(1, 'validacao.nomeObrigatorio').max(200, 'validacao.textoLongo'),
    telefone: z.string().trim().max(30, 'validacao.textoLongo'),
    custo_hora: custoHoraOpcional(moeda),
    unidade_padrao_id: z.string(),
    perfis: z.array(z.string()).min(1, 'validacao.perfilObrigatorio'),
  });

function codigosDosPerfis(usuario: Usuario): string[] {
  return (usuario.perfis ?? []).map((p) => p.codigo);
}

function DadosDoUsuario({
  usuario,
  podeEditar,
  aoConflito,
}: {
  readonly usuario: Usuario;
  readonly podeEditar: boolean;
  readonly aoConflito: (alteracoes: readonly AlteracaoNaoSalva[]) => void;
}) {
  const { t } = useTranslation();
  const custo = useMoedaDeCusto();
  const { moeda } = useRegiao();
  const { precoUnitario } = useFormatos();
  const esquema = useMemo(() => esquemaDoFormulario(moeda), [moeda]);
  const apoioDoCusto =
    usuario.custo_hora && moeda
      ? t('usuario.custoHoraAtual', { valor: precoUnitario(usuario.custo_hora) })
      : custo.semCentavos
        ? t('usuario.custoHoraApoioSemCentavos')
        : t('usuario.custoHoraApoio');
  const alterar = useAlterarUsuario();
  const perfis = usePerfis();
  const unidades = useUnidades();
  const [mensagem, setMensagem] = useState<string | null>(null);
  const [salvo, setSalvo] = useState(false);
  const formulario = useFormulario(esquema, {
    values: {
      nome: usuario.nome,
      telefone: usuario.telefone ?? '',
      custo_hora: usuario.custo_hora ?? '',
      unidade_padrao_id: usuario.unidade_padrao_id,
      perfis: codigosDosPerfis(usuario),
    },
  });
  const { errors, dirtyFields } = formulario.formState;

  const enviar = formulario.handleSubmit((v) => {
    setMensagem(null);
    setSalvo(false);
    const atuais = codigosDosPerfis(usuario);
    const perfisMudaram = v.perfis.length !== atuais.length || v.perfis.some((p) => !atuais.includes(p));
    alterar.mutate(
      {
        id: usuario.id,
        versao: usuario.versao,
        corpo: {
          ...(dirtyFields.nome ? { nome: v.nome } : {}),
          ...(dirtyFields.telefone ? { telefone: vazioParaNulo(v.telefone) } : {}),
          ...(dirtyFields.custo_hora ? { custo_hora: vazioParaNulo(v.custo_hora) && precoUnitarioParaOContrato(v.custo_hora, moeda) } : {}),
          ...(dirtyFields.unidade_padrao_id ? { unidade_padrao_id: v.unidade_padrao_id || null } : {}),
          ...(perfisMudaram ? { perfis: v.perfis } : {}),
        },
      },
      {
        onSuccess: () => setSalvo(true),
        onError: (erro) =>
          tratarErroDeEnvio(erro, {
            setError: formulario.setError,
            definirMensagem: setMensagem,
            aoConflitoDeVersao: aoConflito,
            alteracoes: () =>
              alteracoesNaoSalvas(v, dirtyFields, {
                nome: t('usuario.nome'),
                telefone: t('usuario.telefone'),
                custo_hora: t('usuario.custoHora', { moeda: custo.simbolo }),
                unidade_padrao_id: t('usuario.unidadePadrao'),
                perfis: t('usuario.perfis'),
              }),
          }),
      },
    );
  });

  return (
    <Formulario erros={errors} onSubmit={(e) => void enviar(e)} className="flex max-w-180 flex-col gap-6">
      {mensagem && <Alerta tipo="erro" rolar>{mensagem}</Alerta>}
      {salvo && <Alerta tipo="sucesso" rolar>{t('usuario.salvo')}</Alerta>}
      <div className="grid gap-4 md:grid-cols-2">
        <CampoTexto rotulo={t('usuario.nome')} readOnly={!podeEditar} erro={errors.nome} {...formulario.register('nome')} />
        <CampoTexto rotulo={t('usuario.email')} value={usuario.email ?? ''} readOnly apoio={t('usuario.emailApoio')} />
        <CampoTexto rotulo={t('usuario.telefone')} opcional type="tel" readOnly={!podeEditar} erro={errors.telefone} {...formulario.register('telefone')} />
        <CampoTexto rotulo={t('usuario.custoHora', { moeda: custo.simbolo })} apoio={apoioDoCusto} opcional dado inputMode="decimal" readOnly={!podeEditar} erro={errors.custo_hora} {...formulario.register('custo_hora')} />
        <CampoSelecao rotulo={t('usuario.unidadePadrao')} opcional disabled={!podeEditar} erro={errors.unidade_padrao_id} {...formulario.register('unidade_padrao_id')}>
          <option value="">{t('pessoa.semUnidade')}</option>
          {(unidades.data?.dados ?? []).map((u) => (
            <option key={u.id} value={u.id}>
              {u.nome}
            </option>
          ))}
        </CampoSelecao>
      </div>
      <Controller
        control={formulario.control}
        name="perfis"
        render={({ field, fieldState }) => (
          <SeletorDePerfis
            perfis={perfis.data?.dados ?? []}
            valor={field.value}
            aoMudar={field.onChange}
            erro={fieldState.error}
            somenteLeitura={!podeEditar || usuario.escopo_acesso !== 'interno'}
          />
        )}
      />
      {podeEditar && (
        <div className="flex justify-end">
          <Botao type="submit" variante="primario" carregando={alterar.isPending}>
            {alterar.isPending ? t('comum.salvando') : t('comum.salvar')}
          </Botao>
        </div>
      )}
    </Formulario>
  );
}

function PermissoesAvulsas({ usuario, podeEditar }: { readonly usuario: Usuario; readonly podeEditar: boolean }) {
  const { t } = useTranslation();
  const { data: sessao } = useSessao();
  const podeVerCatalogo = useTemPermissao('adm.gerenciar_perfis');
  const catalogo = usePermissoes(podeVerCatalogo && podeEditar);
  const conceder = useConcederPermissao();
  const retirar = useRetirarPermissao();
  const [escolhida, setEscolhida] = useState('');
  const avulsas = usuario.permissoes_avulsas ?? [];
  const ehEuMesmo = sessao?.usuario.id === usuario.id;
  const disponiveis = (catalogo.data?.dados ?? []).filter((p) => !avulsas.includes(p.codigo));
  const grupos = [...new Set(disponiveis.map((p) => p.grupo))].map(
    (grupo) => [grupo, disponiveis.filter((p) => p.grupo === grupo)] as const,
  );

  return (
    <div className="flex max-w-180 flex-col gap-4">
      <p className="texto-corpo text-texto-secundario">{t('usuario.avulsasTexto')}</p>
      {(conceder.isError || retirar.isError) && (
        <Alerta tipo="erro">{mensagemDoErro(conceder.error ?? retirar.error)}</Alerta>
      )}
      {avulsas.length === 0 ? (
        <p className="texto-corpo-p text-texto-secundario">{t('usuario.avulsasVazio')}</p>
      ) : (
        <ul className="flex flex-wrap gap-2">
          {avulsas.map((codigo) => (
            <li key={codigo} className="inline-flex items-center gap-1 rounded-sm bg-superficie-selecionada py-1 pr-1 pl-3 texto-dado">
              {codigo}
              {podeEditar && !ehEuMesmo && (
                <Botao
                  variante="icone"
                  className="size-8"
                  aria-label={t('usuario.retirarPermissao', { permissao: codigo })}
                  carregando={retirar.isPending && retirar.variables?.permissao === codigo}
                  onClick={() => retirar.mutate({ id: usuario.id, permissao: codigo })}
                >
                  <X size={16} aria-hidden="true" />
                </Botao>
              )}
            </li>
          ))}
        </ul>
      )}
      {podeEditar && ehEuMesmo && <p className="texto-legenda text-texto-secundario">{t('usuario.avulsasProprias')}</p>}
      {podeEditar && !ehEuMesmo && podeVerCatalogo && (
        <div className="flex items-end gap-3">
          <div className="w-96 max-w-full">
            <CampoSelecao rotulo={t('usuario.conceder')} name="permissao" value={escolhida} onChange={(e) => setEscolhida(e.target.value)} disabled={catalogo.isPending}>
              <option value="">{t('usuario.escolhaPermissao')}</option>
              {grupos.map(([grupo, itens]) => (
                <optgroup key={grupo} label={grupo}>
                  {itens.map((p) => (
                    <option key={p.codigo} value={p.codigo}>
                      {p.codigo}
                    </option>
                  ))}
                </optgroup>
              ))}
            </CampoSelecao>
          </div>
          <Botao
            disabled={escolhida === ''}
            carregando={conceder.isPending}
            onClick={() => conceder.mutate({ id: usuario.id, permissao: escolhida }, { onSuccess: () => setEscolhida('') })}
          >
            {t('usuario.concederBotao')}
          </Botao>
        </div>
      )}
      {podeEditar && !ehEuMesmo && !podeVerCatalogo && <p className="texto-legenda text-texto-secundario">{t('usuario.semCatalogo')}</p>}
    </div>
  );
}

export function TelaUsuarioDetalhe() {
  const { t } = useTranslation();
  const { id = '' } = useParams();
  const local = useLocation();
  const consulta = useUsuario(id);
  const desativar = useDesativarUsuario();
  const reativar = useReativarUsuario();
  const reenviar = useReenviarConvite();
  const encerrar = useEncerrarSessoes();
  const podeEditar = useTemPermissao('adm.gerenciar_usuarios');
  const podeVerAuditoria = useTemPermissao('adm.ler_auditoria');
  const podeVerPessoas = useTemPermissao('adm.ler_pessoas');
  const [confirmando, setConfirmando] = useState<'desativar' | 'encerrar' | null>(null);
  const [conflito, setConflito] = useState<readonly AlteracaoNaoSalva[] | null>(null);
  const conviteDoCadastro = (local.state as { conviteEnviado?: boolean } | null)?.conviteEnviado;
  const [aviso, setAviso] = useState<{ tipo: 'sucesso' | 'atencao'; texto: string } | null>(
    conviteDoCadastro === undefined
      ? null
      : conviteDoCadastro
        ? { tipo: 'sucesso', texto: t('usuario.conviteEnviado') }
        : { tipo: 'atencao', texto: t('usuario.conviteNaoEnviado') },
  );

  if (consulta.isPending) {
    return (
      <>
        <CabecalhoPagina titulo={t('usuario.titulo')} />
        <Carregando rotulo={t('usuario.carregando')} linhas={4} />
      </>
    );
  }
  if (consulta.isError) {
    if (ehNaoEncontrado(consulta.error)) {
      return <EstadoVazio titulo={t('usuario.naoEncontradoTitulo')} texto={t('usuario.naoEncontrado')} />;
    }
    return <EstadoErro erro={consulta.error} aoTentarDeNovo={() => void consulta.refetch()} />;
  }

  const usuario = consulta.data;
  const interno = usuario.escopo_acesso === 'interno';
  const acoes = (
    <>
      {podeVerAuditoria && (
        <Link
          to={`/administracao/auditoria?entidade=usuario&entidade_id=${usuario.id}`}
          className="inline-flex h-10 items-center px-2 texto-botao underline"
        >
          {t('auditoria.historicoDoRegistro')}
        </Link>
      )}
      {podeEditar && usuario.situacao === 'convidado' && (
        <Botao
          carregando={reenviar.isPending}
          onClick={() => reenviar.mutate(usuario.id, { onSuccess: () => setAviso({ tipo: 'sucesso', texto: t('usuario.conviteReenviado', { email: usuario.email ?? '' }) }) })}
        >
          {t('usuario.reenviarConvite')}
        </Botao>
      )}
      {podeEditar && usuario.situacao !== 'inativo' && <Botao onClick={() => setConfirmando('encerrar')}>{t('usuario.encerrarSessoes')}</Botao>}
      {podeEditar && usuario.situacao !== 'inativo' && <Botao onClick={() => setConfirmando('desativar')}>{t('usuario.desativar')}</Botao>}
      {podeEditar && usuario.situacao === 'inativo' && (
        <Botao
          variante="primario"
          carregando={reativar.isPending}
          onClick={() =>
            reativar.mutate(
              { id: usuario.id, versao: usuario.versao },
              { onSuccess: (r) => setAviso({ tipo: 'sucesso', texto: r.convite_enviado ? t('usuario.reativadoComConvite') : t('usuario.reativado') }) },
            )
          }
        >
          {t('usuario.reativar')}
        </Botao>
      )}
    </>
  );

  return (
    <>
      <CabecalhoPagina
        titulo={usuario.nome}
        trilha={[{ rotulo: t('usuario.titulo'), para: '/administracao/usuarios' }, { rotulo: usuario.nome }]}
        metadados={
          <>
            <ChipSituacao situacao={usuario.situacao} />
            <span>{usuario.email}</span>
            {usuario.pessoa_id && podeVerPessoas && (
              <Link to={`/administracao/pessoas/${usuario.pessoa_id}`} className="underline">
                {t('usuario.verPessoa')}
              </Link>
            )}
          </>
        }
        acoes={acoes}
      />
      <div className="mb-4 flex flex-col gap-3">
        {aviso && <Alerta tipo={aviso.tipo}>{aviso.texto}</Alerta>}
        {(reativar.isError || reenviar.isError) && <Alerta tipo="erro">{mensagemDoErro(reativar.error ?? reenviar.error)}</Alerta>}
        {conflito && <AlertaDeConflito alteracoes={conflito} />}
        {encerrar.isSuccess && <Alerta tipo="sucesso">{t('usuario.sessoesEncerradas', { count: encerrar.data.sessoes_encerradas })}</Alerta>}
        {!interno && <Alerta tipo="informacao">{t('usuario.externo')}</Alerta>}
      </div>

      <div className="flex flex-col gap-10">
        <section aria-labelledby="usuario-dados" className="flex flex-col gap-4">
          <h2 id="usuario-dados" className="texto-h3">
            {t('usuario.dados')}
          </h2>
          <DadosDoUsuario
            usuario={usuario}
            podeEditar={podeEditar && usuario.situacao !== 'inativo'}
            aoConflito={(alteracoes) => {
              setConflito(alteracoes);
              void consulta.refetch();
            }}
          />
        </section>
        <section aria-labelledby="usuario-avulsas" className="flex flex-col gap-4">
          <h2 id="usuario-avulsas" className="texto-h3">
            {t('usuario.avulsas')}
          </h2>
          <PermissoesAvulsas usuario={usuario} podeEditar={podeEditar && usuario.situacao !== 'inativo'} />
        </section>
      </div>

      <Confirmacao
        aberto={confirmando === 'desativar'}
        aoMudar={(aberto) => {
          if (!aberto) {
            setConfirmando(null);
            desativar.reset();
          }
        }}
        titulo={t('usuario.desativarTitulo', { nome: usuario.nome })}
        texto={t('usuario.desativarTexto')}
        rotuloConfirmar={t('usuario.desativarBotao')}
        perigo
        carregando={desativar.isPending}
        erro={desativar.isError ? <Alerta tipo="erro">{mensagemDoErro(desativar.error)}</Alerta> : undefined}
        aoConfirmar={() => desativar.mutate({ id: usuario.id, versao: usuario.versao }, { onSuccess: () => setConfirmando(null) })}
      />
      <Confirmacao
        aberto={confirmando === 'encerrar'}
        aoMudar={(aberto) => {
          if (!aberto) {
            setConfirmando(null);
            encerrar.reset();
          }
        }}
        titulo={t('usuario.encerrarTitulo', { nome: usuario.nome })}
        texto={t('usuario.encerrarTexto')}
        rotuloConfirmar={t('usuario.encerrarBotao')}
        carregando={encerrar.isPending}
        erro={encerrar.isError ? <Alerta tipo="erro">{mensagemDoErro(encerrar.error)}</Alerta> : undefined}
        aoConfirmar={() => encerrar.mutate(usuario.id, { onSuccess: () => setConfirmando(null) })}
      />
    </>
  );
}
