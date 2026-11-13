// Aceite de convite (/convite#token=...), RF-ADM-004 e RF-ADM-006. A página
// consulta o convite com o token do fragmento: se o e-mail ainda não tem
// login, pede a senha nova; se já tem (outra empresa), entra com a que tem.

import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router';

import { Alerta } from '../../../componentes/Alerta';
import { Botao } from '../../../componentes/Botao';
import { CampoSenha } from '../../../componentes/Campo';
import { Carregando, EstadoErro } from '../../../componentes/Estados';
import { Formulario } from '../../../componentes/Formulario';
import { aplicarErrosDoServidor, ehErroDaAPI, mensagemDoErro } from '../../../lib/erros';
import { senhaNova, useFormulario, z } from '../../../lib/formularios';
import { useTokenDoFragmento } from '../../../lib/fragmento';
import { useAceitarConvite, useConsultarConvite } from '../sessao';

import { CartaoPublico } from './CartaoPublico';

const esquemaComSenha = z
  .object({ senha: senhaNova, confirmacao: z.string() })
  .refine((v) => v.senha === v.confirmacao, { path: ['confirmacao'], message: 'validacao.senhasDiferentes' });

function LinkInvalido() {
  const { t } = useTranslation();
  return (
    <CartaoPublico titulo={t('convite.titulo')}>
      <Alerta tipo="erro">{t('convite.linkInvalido')}</Alerta>
      <Link to="/entrar" className="texto-corpo-p underline">
        {t('convite.irParaEntrar')}
      </Link>
    </CartaoPublico>
  );
}

export function TelaConvite() {
  const { t } = useTranslation();
  const token = useTokenDoFragmento();
  const consulta = useConsultarConvite(token);
  const aceitar = useAceitarConvite();
  const [aceito, setAceito] = useState(false);
  const [mensagem, setMensagem] = useState<string | null>(null);
  const formulario = useFormulario(esquemaComSenha, { defaultValues: { senha: '', confirmacao: '' } });

  if (token === null) return <LinkInvalido />;

  if (consulta.isPending) {
    return (
      <CartaoPublico titulo={t('convite.titulo')}>
        <Carregando rotulo={t('convite.carregando')} linhas={2} />
      </CartaoPublico>
    );
  }

  if (consulta.isError) {
    if (ehErroDaAPI(consulta.error) && consulta.error.codigo === 'link_email_invalido') return <LinkInvalido />;
    return (
      <CartaoPublico titulo={t('convite.titulo')}>
        <EstadoErro erro={consulta.error} aoTentarDeNovo={() => void consulta.refetch()} />
      </CartaoPublico>
    );
  }

  const convite = consulta.data;

  if (aceito) {
    return (
      <CartaoPublico titulo={t('convite.titulo')}>
        <Alerta tipo="sucesso">{t('convite.aceito', { email: convite.email })}</Alerta>
        <Link to="/entrar" className="texto-corpo-p underline">
          {t('convite.irParaEntrar')}
        </Link>
      </CartaoPublico>
    );
  }

  const enviar = (senha: string | undefined) => {
    setMensagem(null);
    aceitar.mutate(
      senha === undefined ? { token } : { token, senha },
      {
        onSuccess: () => setAceito(true),
        onError: (erro) => {
          if (ehErroDaAPI(erro) && erro.codigo === 'link_email_invalido') {
            setMensagem(t('convite.linkInvalido'));
            return;
          }
          if (aplicarErrosDoServidor(erro, formulario.setError)) return;
          setMensagem(mensagemDoErro(erro));
        },
      },
    );
  };

  return (
    <CartaoPublico titulo={t('convite.titulo')}>
      <p className="texto-corpo text-texto-secundario">
        {convite.pede_senha ? t('convite.textoNovaSenha', { email: convite.email }) : t('convite.textoSenhaExistente', { email: convite.email })}
      </p>
      {mensagem && <Alerta tipo="erro">{mensagem}</Alerta>}
      {convite.pede_senha ? (
        <Formulario
          erros={formulario.formState.errors}
          onSubmit={(e) => void formulario.handleSubmit((v) => enviar(v.senha))(e)}
          className="flex flex-col gap-4"
        >
          <CampoSenha
            rotulo={t('convite.senha')}
            apoio={t('redefinir.apoioSenha')}
            autoComplete="new-password"
            erro={formulario.formState.errors.senha}
            {...formulario.register('senha')}
          />
          <CampoSenha
            rotulo={t('redefinir.confirmacao')}
            autoComplete="new-password"
            erro={formulario.formState.errors.confirmacao}
            {...formulario.register('confirmacao')}
          />
          <Botao type="submit" variante="primario" carregando={aceitar.isPending} className="w-full">
            {aceitar.isPending ? t('convite.aceitando') : t('convite.botao')}
          </Botao>
        </Formulario>
      ) : (
        <Botao variante="primario" carregando={aceitar.isPending} onClick={() => enviar(undefined)} className="w-full">
          {aceitar.isPending ? t('convite.aceitando') : t('convite.botao')}
        </Botao>
      )}
    </CartaoPublico>
  );
}
