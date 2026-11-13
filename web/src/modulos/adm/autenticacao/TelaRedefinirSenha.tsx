// Redefinição de senha pelo link do e-mail (RF-ADM-006). O token vem no
// fragmento da URL (`#token=`), é lido, tirado da barra de endereço e enviado
// no corpo da chamada. Redefinir encerra todas as sessões da pessoa.

import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router';

import { Alerta } from '../../../componentes/Alerta';
import { Botao } from '../../../componentes/Botao';
import { CampoSenha } from '../../../componentes/Campo';
import { Formulario } from '../../../componentes/Formulario';
import { aplicarErrosDoServidor, ehErroDaAPI, mensagemDoErro } from '../../../lib/erros';
import { senhaNova, useFormulario, z } from '../../../lib/formularios';
import { useTokenDoFragmento } from '../../../lib/fragmento';
import { useRedefinirSenha } from '../sessao';

import { CartaoPublico } from './CartaoPublico';

const esquema = z
  .object({ nova_senha: senhaNova, confirmacao: z.string() })
  .refine((v) => v.nova_senha === v.confirmacao, { path: ['confirmacao'], message: 'validacao.senhasDiferentes' });

export function TelaRedefinirSenha() {
  const { t } = useTranslation();
  const token = useTokenDoFragmento();
  const redefinir = useRedefinirSenha();
  const [concluido, setConcluido] = useState(false);
  const [linkInvalido, setLinkInvalido] = useState(token === null);
  const [mensagem, setMensagem] = useState<string | null>(null);
  const formulario = useFormulario(esquema, { defaultValues: { nova_senha: '', confirmacao: '' } });

  const enviar = formulario.handleSubmit((valores) => {
    if (token === null) return;
    setMensagem(null);
    redefinir.mutate(
      { token, nova_senha: valores.nova_senha },
      {
        onSuccess: () => setConcluido(true),
        onError: (erro) => {
          if (ehErroDaAPI(erro) && erro.codigo === 'link_email_invalido') {
            setLinkInvalido(true);
            return;
          }
          if (aplicarErrosDoServidor(erro, formulario.setError, { mapa: { senha: 'nova_senha' } })) return;
          setMensagem(mensagemDoErro(erro));
        },
      },
    );
  });

  if (concluido) {
    return (
      <CartaoPublico titulo={t('redefinir.titulo')}>
        <Alerta tipo="sucesso">{t('redefinir.concluido')}</Alerta>
        <Link to="/entrar" className="texto-corpo-p underline">
          {t('redefinir.entrar')}
        </Link>
      </CartaoPublico>
    );
  }

  if (linkInvalido) {
    return (
      <CartaoPublico titulo={t('redefinir.titulo')}>
        <Alerta tipo="erro">{t('redefinir.linkInvalido')}</Alerta>
        <Link to="/esqueci-a-senha" className="texto-corpo-p underline">
          {t('redefinir.pedirOutro')}
        </Link>
      </CartaoPublico>
    );
  }

  return (
    <CartaoPublico titulo={t('redefinir.titulo')}>
      {mensagem && <Alerta tipo="erro">{mensagem}</Alerta>}
      <Formulario erros={formulario.formState.errors} onSubmit={(e) => void enviar(e)} className="flex flex-col gap-4">
        <CampoSenha
          rotulo={t('redefinir.novaSenha')}
          apoio={t('redefinir.apoioSenha')}
          autoComplete="new-password"
          erro={formulario.formState.errors.nova_senha}
          {...formulario.register('nova_senha')}
        />
        <CampoSenha
          rotulo={t('redefinir.confirmacao')}
          autoComplete="new-password"
          erro={formulario.formState.errors.confirmacao}
          {...formulario.register('confirmacao')}
        />
        <Botao type="submit" variante="primario" carregando={redefinir.isPending} className="w-full">
          {redefinir.isPending ? t('redefinir.salvando') : t('redefinir.botao')}
        </Botao>
      </Formulario>
    </CartaoPublico>
  );
}
