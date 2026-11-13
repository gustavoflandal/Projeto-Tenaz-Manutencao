// Recuperação de senha (RF-ADM-006): pede o e-mail e responde sempre igual,
// exista a conta ou não (202), para não revelar quem tem conta.

import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router';

import { Alerta } from '../../../componentes/Alerta';
import { Botao } from '../../../componentes/Botao';
import { CampoTexto } from '../../../componentes/Campo';
import { Formulario } from '../../../componentes/Formulario';
import { aplicarErrosDoServidor, mensagemDoErro } from '../../../lib/erros';
import { emailObrigatorio, useFormulario, z } from '../../../lib/formularios';
import { useRecuperarSenha } from '../sessao';

import { CartaoPublico } from './CartaoPublico';

const esquema = z.object({ email: emailObrigatorio });

export function TelaEsqueciSenha() {
  const { t } = useTranslation();
  const recuperar = useRecuperarSenha();
  const [enviadoPara, setEnviadoPara] = useState<string | null>(null);
  const [mensagem, setMensagem] = useState<string | null>(null);
  const formulario = useFormulario(esquema, { defaultValues: { email: '' } });

  const enviar = formulario.handleSubmit((valores) => {
    setMensagem(null);
    recuperar.mutate(valores.email, {
      onSuccess: () => setEnviadoPara(valores.email),
      onError: (erro) => {
        if (!aplicarErrosDoServidor(erro, formulario.setError)) setMensagem(mensagemDoErro(erro));
      },
    });
  });

  if (enviadoPara) {
    return (
      <CartaoPublico titulo={t('esqueci.titulo')}>
        <Alerta tipo="sucesso">{t('esqueci.enviado', { email: enviadoPara })}</Alerta>
        <Link to="/entrar" className="texto-corpo-p underline">
          {t('esqueci.voltar')}
        </Link>
      </CartaoPublico>
    );
  }

  return (
    <CartaoPublico titulo={t('esqueci.titulo')}>
      <p className="texto-corpo text-texto-secundario">{t('esqueci.texto')}</p>
      {mensagem && <Alerta tipo="erro">{mensagem}</Alerta>}
      <Formulario erros={formulario.formState.errors} onSubmit={(e) => void enviar(e)} className="flex flex-col gap-4">
        <CampoTexto
          rotulo={t('entrar.email')}
          type="email"
          autoComplete="username"
          erro={formulario.formState.errors.email}
          {...formulario.register('email')}
        />
        <Botao type="submit" variante="primario" carregando={recuperar.isPending} className="w-full">
          {recuperar.isPending ? t('esqueci.enviando') : t('esqueci.botao')}
        </Botao>
        <Link to="/entrar" className="texto-corpo-p underline">
          {t('esqueci.voltar')}
        </Link>
      </Formulario>
    </CartaoPublico>
  );
}
