// Tela de login (seção 5.6 do documento 7, RF-ADM-006).
//
// Três passos na mesma página: e-mail e senha; código do segundo fator (TOTP)
// quando o servidor responde 401 mfa_obrigatorio com o desafio; e a escolha da
// empresa quando a pessoa tem acesso a mais de uma (documento 3, seção 4.2).
// Toda falha de credencial recebe a mesma mensagem, sem dizer se o e-mail existe.

import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Link, useNavigate, useSearchParams } from 'react-router';

import { Alerta } from '../../../componentes/Alerta';
import { Botao } from '../../../componentes/Botao';
import { CampoSenha, CampoTexto } from '../../../componentes/Campo';
import { Formulario } from '../../../componentes/Formulario';
import type { Problema } from '../../../lib/api';
import { aplicarErrosDoServidor, ehErroDaAPI, mensagemDoErro } from '../../../lib/erros';
import { emailObrigatorio, useFormulario, z } from '../../../lib/formularios';
import { retornoSeguro } from '../../../lib/retorno';
import type { Sessao } from '../../../lib/tipos';

import { CartaoPublico } from './CartaoPublico';
import { EscolhaDeEmpresa } from './EscolhaDeEmpresa';
import { useEntrar, useVerificarMFA } from '../sessao';

const esquemaCredenciais = z.object({
  email: emailObrigatorio,
  senha: z.string().min(1, 'validacao.senhaObrigatoria'),
});

const esquemaCodigo = z.object({
  codigo: z.string().trim().regex(/^[0-9]{6}$/, 'validacao.codigoMfa'),
});

type Desafio = Pick<Problema, 'desafio_id'> & { readonly desafio_id: string };

export function TelaEntrar() {
  const { t } = useTranslation();
  const navegar = useNavigate();
  const [parametros] = useSearchParams();
  const retorno = retornoSeguro(parametros.get('retorno'));
  const sessaoEncerrada = parametros.get('motivo') === 'sessao';

  const [desafio, setDesafio] = useState<Desafio | null>(null);
  const [sessaoAberta, setSessaoAberta] = useState<Sessao | null>(null);
  const [aviso, setAviso] = useState<string | null>(null);
  const [mensagem, setMensagem] = useState<string | null>(null);
  const [emailDigitado, setEmailDigitado] = useState('');

  const entrar = useEntrar();
  const verificar = useVerificarMFA();

  const credenciais = useFormulario(esquemaCredenciais, { defaultValues: { email: '', senha: '' } });
  const codigo = useFormulario(esquemaCodigo, { defaultValues: { codigo: '' } });

  const concluir = (sessao: Sessao) => {
    if (sessao.mfa_pendente) {
      void navegar('/ativar-mfa', { replace: true });
    } else if ((sessao.empresas?.length ?? 0) > 1) {
      setSessaoAberta(sessao);
    } else {
      void navegar(retorno, { replace: true });
    }
  };

  const enviarCredenciais = credenciais.handleSubmit((valores) => {
    setMensagem(null);
    setAviso(null);
    setEmailDigitado(valores.email);
    entrar.mutate(
      { email: valores.email, senha: valores.senha },
      {
        onSuccess: concluir,
        onError: (erro) => {
          if (ehErroDaAPI(erro) && erro.codigo === 'mfa_obrigatorio' && erro.problema?.desafio_id) {
            setDesafio({ desafio_id: erro.problema.desafio_id });
            return;
          }
          if (aplicarErrosDoServidor(erro, credenciais.setError)) return;
          setMensagem(mensagemDoErro(erro));
        },
      },
    );
  });

  const enviarCodigo = codigo.handleSubmit((valores) => {
    if (!desafio) return;
    setMensagem(null);
    verificar.mutate(
      { desafio_id: desafio.desafio_id, codigo: valores.codigo },
      {
        onSuccess: concluir,
        onError: (erro) => {
          const restantes = ehErroDaAPI(erro) ? erro.problema?.tentativas_restantes : undefined;
          if (ehErroDaAPI(erro) && erro.codigo === 'mfa_codigo_invalido' && restantes !== 0) {
            codigo.setValue('codigo', '');
            setMensagem(mensagemDoErro(erro));
            return;
          }
          // Desafio vencido, gasto ou sem tentativas: o login recomeça pela senha.
          setDesafio(null);
          credenciais.resetField('senha');
          codigo.reset();
          setAviso(t('entrar.desafioEncerrado'));
        },
      },
    );
  });

  if (sessaoAberta) {
    return (
      <CartaoPublico titulo={t('entrar.escolhaEmpresaTitulo')}>
        <EscolhaDeEmpresa sessao={sessaoAberta} aoConcluir={() => void navegar(retorno, { replace: true })} />
      </CartaoPublico>
    );
  }

  if (desafio) {
    return (
      <CartaoPublico titulo={t('entrar.mfaTitulo')}>
        <p className="texto-corpo text-texto-secundario">{t('entrar.mfaTexto', { email: emailDigitado })}</p>
        {mensagem && <Alerta tipo="erro">{mensagem}</Alerta>}
        <Formulario erros={codigo.formState.errors} onSubmit={(e) => void enviarCodigo(e)} className="flex flex-col gap-4">
          <CampoTexto
            rotulo={t('entrar.codigo')}
            apoio={t('entrar.codigoApoio')}
            inputMode="numeric"
            autoComplete="one-time-code"
            maxLength={6}
            dado
            erro={codigo.formState.errors.codigo}
            {...codigo.register('codigo')}
          />
          <Botao type="submit" variante="primario" carregando={verificar.isPending} className="w-full">
            {verificar.isPending ? t('entrar.verificando') : t('entrar.verificar')}
          </Botao>
          <Botao
            variante="terciario"
            onClick={() => {
              setDesafio(null);
              setMensagem(null);
              codigo.reset();
            }}
          >
            {t('entrar.voltarParaSenha')}
          </Botao>
        </Formulario>
      </CartaoPublico>
    );
  }

  return (
    <CartaoPublico titulo={t('entrar.titulo')}>
      {sessaoEncerrada && <Alerta tipo="informacao">{t('entrar.sessaoEncerrada')}</Alerta>}
      {aviso && <Alerta tipo="atencao">{aviso}</Alerta>}
      {mensagem && <Alerta tipo="erro">{mensagem}</Alerta>}
      <Formulario erros={credenciais.formState.errors} onSubmit={(e) => void enviarCredenciais(e)} className="flex flex-col gap-4">
        <CampoTexto
          rotulo={t('entrar.email')}
          type="email"
          autoComplete="username"
          erro={credenciais.formState.errors.email}
          {...credenciais.register('email')}
        />
        <CampoSenha
          rotulo={t('entrar.senha')}
          autoComplete="current-password"
          erro={credenciais.formState.errors.senha}
          {...credenciais.register('senha')}
        />
        <Botao type="submit" variante="primario" carregando={entrar.isPending} className="mt-2 w-full">
          {entrar.isPending ? t('entrar.entrando') : t('entrar.botao')}
        </Botao>
        <Link to="/esqueci-a-senha" className="texto-corpo-p underline">
          {t('entrar.esqueci')}
        </Link>
      </Formulario>
      <p className="border-t border-borda pt-4 texto-corpo-p text-texto-secundario">
        {t('entrar.semConta')}{' '}
        <Link to="/criar-conta" className="underline">
          {t('entrar.criarConta')}
        </Link>
      </p>
    </CartaoPublico>
  );
}
