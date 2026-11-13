// Minha conta (/conta): dados do próprio usuário, troca de senha e segundo
// fator. Serve a todos os escopos (interno, prestador e locatário): usa só as
// rotas do próprio usuário, que não pedem permissão fina (documento 3, 4.2).

import { useState } from 'react';
import { useTranslation } from 'react-i18next';

import { Alerta } from '../../../componentes/Alerta';
import { AlertaDeConflito } from '../../../componentes/AlertaDeConflito';
import { Botao } from '../../../componentes/Botao';
import { CabecalhoPagina } from '../../../componentes/CabecalhoPagina';
import { CampoSenha, CampoTexto } from '../../../componentes/Campo';
import { SeletorDeIdioma } from '../../../componentes/SeletorDeIdioma';
import { Carregando, EstadoErro } from '../../../componentes/Estados';
import { Formulario } from '../../../componentes/Formulario';
import { alteracoesNaoSalvas, aplicarErrosDoServidor, mensagemDoErro, tratarErroDeEnvio, type AlteracaoNaoSalva } from '../../../lib/erros';
import { senhaNova, useFormulario, vazioParaNulo, z } from '../../../lib/formularios';
import { useAlterarUsuarioEu, useUsuarioEu } from '../usuarios';
import { useSessao, useTrocarSenha } from '../sessao';

import { SecaoMFA } from './SecaoMFA';

const esquemaDados = z.object({
  nome: z.string().trim().min(1, 'validacao.nomeObrigatorio').max(200, 'validacao.textoLongo'),
  telefone: z.string().trim().max(30, 'validacao.textoLongo'),
});

const esquemaSenha = z
  .object({ senha_atual: z.string().min(1, 'validacao.senhaAtualObrigatoria'), nova_senha: senhaNova, confirmacao: z.string() })
  .refine((v) => v.nova_senha === v.confirmacao, { path: ['confirmacao'], message: 'validacao.senhasDiferentes' });

function SecaoDados() {
  const { t } = useTranslation();
  const consulta = useUsuarioEu();
  const alterar = useAlterarUsuarioEu();
  const [mensagem, setMensagem] = useState<string | null>(null);
  const [salvo, setSalvo] = useState(false);
  const [conflito, setConflito] = useState<readonly AlteracaoNaoSalva[] | null>(null);
  const usuario = consulta.data;
  const formulario = useFormulario(esquemaDados, {
    values: usuario ? { nome: usuario.nome, telefone: usuario.telefone ?? '' } : undefined,
  });

  if (consulta.isPending) return <Carregando rotulo={t('conta.carregando')} linhas={2} />;
  if (consulta.isError) return <EstadoErro erro={consulta.error} aoTentarDeNovo={() => void consulta.refetch()} />;

  const enviar = formulario.handleSubmit((v) => {
    setMensagem(null);
    setSalvo(false);
    setConflito(null);
    alterar.mutate(
      { versao: consulta.data.versao, corpo: { nome: v.nome, telefone: vazioParaNulo(v.telefone) } },
      {
        onSuccess: () => setSalvo(true),
        onError: (erro) => {
          tratarErroDeEnvio(erro, {
            setError: formulario.setError,
            definirMensagem: setMensagem,
            aoConflitoDeVersao: (alteracoes) => {
              setConflito(alteracoes);
              void consulta.refetch();
            },
            alteracoes: () =>
              alteracoesNaoSalvas(v, formulario.formState.dirtyFields, { nome: t('conta.nome'), telefone: t('conta.telefone') }),
          });
        },
      },
    );
  });

  return (
    <Formulario erros={formulario.formState.errors} onSubmit={(e) => void enviar(e)} className="flex max-w-120 flex-col gap-4">
      {mensagem && <Alerta tipo="erro" rolar>{mensagem}</Alerta>}
      {conflito && <AlertaDeConflito alteracoes={conflito} />}
      {salvo && <Alerta tipo="sucesso" rolar>{t('conta.salvo')}</Alerta>}
      <CampoTexto rotulo={t('conta.nome')} autoComplete="name" erro={formulario.formState.errors.nome} {...formulario.register('nome')} />
      <CampoTexto rotulo={t('conta.email')} value={consulta.data.email ?? ''} readOnly apoio={t('conta.emailApoio')} />
      <CampoTexto rotulo={t('conta.telefone')} type="tel" autoComplete="tel" opcional erro={formulario.formState.errors.telefone} {...formulario.register('telefone')} />
      <div>
        <Botao type="submit" variante="primario" carregando={alterar.isPending}>
          {alterar.isPending ? t('comum.salvando') : t('comum.salvar')}
        </Botao>
      </div>
    </Formulario>
  );
}

function SecaoSenha() {
  const { t } = useTranslation();
  const trocar = useTrocarSenha();
  const [mensagem, setMensagem] = useState<string | null>(null);
  const [trocada, setTrocada] = useState(false);
  const formulario = useFormulario(esquemaSenha, { defaultValues: { senha_atual: '', nova_senha: '', confirmacao: '' } });

  const enviar = formulario.handleSubmit((v) => {
    setMensagem(null);
    setTrocada(false);
    trocar.mutate(
      { senha_atual: v.senha_atual, nova_senha: v.nova_senha },
      {
        onSuccess: () => {
          setTrocada(true);
          formulario.reset();
        },
        onError: (erro) => {
          if (!aplicarErrosDoServidor(erro, formulario.setError, { mapa: { senha: 'nova_senha' } })) {
            setMensagem(mensagemDoErro(erro));
          }
        },
      },
    );
  });

  return (
    <Formulario erros={formulario.formState.errors} onSubmit={(e) => void enviar(e)} className="flex max-w-120 flex-col gap-4">
      {mensagem && <Alerta tipo="erro">{mensagem}</Alerta>}
      {trocada && <Alerta tipo="sucesso">{t('conta.senhaTrocada')}</Alerta>}
      <CampoSenha rotulo={t('conta.senhaAtual')} autoComplete="current-password" erro={formulario.formState.errors.senha_atual} {...formulario.register('senha_atual')} />
      <CampoSenha rotulo={t('redefinir.novaSenha')} apoio={t('redefinir.apoioSenha')} autoComplete="new-password" erro={formulario.formState.errors.nova_senha} {...formulario.register('nova_senha')} />
      <CampoSenha rotulo={t('redefinir.confirmacao')} autoComplete="new-password" erro={formulario.formState.errors.confirmacao} {...formulario.register('confirmacao')} />
      <p className="texto-legenda text-texto-secundario">{t('conta.senhaEncerraOutras')}</p>
      <div>
        <Botao type="submit" variante="primario" carregando={trocar.isPending}>
          {trocar.isPending ? t('comum.salvando') : t('conta.trocarSenha')}
        </Botao>
      </div>
    </Formulario>
  );
}

export function TelaConta() {
  const { t } = useTranslation();
  const { data: sessao } = useSessao();
  if (!sessao) return null;

  return (
    <>
      <CabecalhoPagina titulo={t('conta.titulo')} />
      <div className="flex flex-col gap-10">
        <section aria-labelledby="conta-dados" className="flex flex-col gap-4">
          <h2 id="conta-dados" className="texto-h3">
            {t('conta.dados')}
          </h2>
          <SecaoDados />
        </section>
        <section aria-labelledby="conta-idioma" className="flex flex-col gap-4">
          <h2 id="conta-idioma" className="texto-h3">
            {t('idioma.titulo')}
          </h2>
          <SeletorDeIdioma />
        </section>
        <section aria-labelledby="conta-senha" className="flex flex-col gap-4">
          <h2 id="conta-senha" className="texto-h3">
            {t('conta.senha')}
          </h2>
          <SecaoSenha />
        </section>
        <section aria-labelledby="conta-mfa" className="flex flex-col gap-4">
          <h2 id="conta-mfa" className="texto-h3">
            {t('conta.mfa')}
          </h2>
          <SecaoMFA mfaAtivo={sessao.usuario.mfa_ativo} />
        </section>
      </div>
    </>
  );
}
