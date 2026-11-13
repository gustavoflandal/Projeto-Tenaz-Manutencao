// Segundo fator por TOTP (RF-ADM-006): ativar com QR code e primeiro código, e
// desativar com um código. Serve à tela da conta e à sessão restrita de quem a
// empresa obriga a ativar (documento 3, seção 4.2).

import { useState } from 'react';
import { useTranslation } from 'react-i18next';

import { Alerta } from '../../../componentes/Alerta';
import { Botao } from '../../../componentes/Botao';
import { CampoTexto } from '../../../componentes/Campo';
import { CodigoQR } from '../../../componentes/CodigoQR';
import { Formulario } from '../../../componentes/Formulario';
import { ehErroDaAPI, mensagemDoErro } from '../../../lib/erros';
import { useFormulario, z } from '../../../lib/formularios';
import { useConfirmarMFA, useDesativarMFA, useIniciarMFA } from '../sessao';

const esquemaCodigo = z.object({
  codigo: z.string().trim().regex(/^[0-9]{6}$/, 'validacao.codigoMfa'),
});

function FormularioDeCodigo({
  rotuloBotao,
  rotuloEnviando,
  carregando,
  perigo,
  aoEnviar,
  erro,
}: {
  readonly rotuloBotao: string;
  readonly rotuloEnviando: string;
  readonly carregando: boolean;
  readonly perigo?: boolean;
  readonly aoEnviar: (codigo: string, definirErro: (mensagem: string) => void) => void;
  readonly erro: string | null;
}) {
  const { t } = useTranslation();
  const formulario = useFormulario(esquemaCodigo, { defaultValues: { codigo: '' } });
  return (
    <Formulario
      erros={formulario.formState.errors}
      onSubmit={(e) => void formulario.handleSubmit((v) => aoEnviar(v.codigo, (m) => formulario.setError('codigo', { type: 'servidor', message: m })))(e)}
      className="flex max-w-100 flex-col gap-4"
    >
      {erro && <Alerta tipo="erro">{erro}</Alerta>}
      <CampoTexto
        rotulo={t('mfa.codigo')}
        apoio={t('mfa.codigoApoio')}
        inputMode="numeric"
        autoComplete="one-time-code"
        maxLength={6}
        dado
        erro={formulario.formState.errors.codigo}
        {...formulario.register('codigo')}
      />
      <div>
        <Botao type="submit" variante={perigo ? 'perigo' : 'primario'} carregando={carregando}>
          {carregando ? rotuloEnviando : rotuloBotao}
        </Botao>
      </div>
    </Formulario>
  );
}

export interface SecaoMFAProps {
  readonly mfaAtivo: boolean;
  /** A empresa exige o segundo fator: não há como desativar (409 exigido_pela_empresa). */
  readonly exigidoPelaEmpresa?: boolean;
}

export function SecaoMFA({ mfaAtivo, exigidoPelaEmpresa }: SecaoMFAProps) {
  const { t } = useTranslation();
  const iniciar = useIniciarMFA();
  const confirmar = useConfirmarMFA();
  const desativar = useDesativarMFA();
  const [erro, setErro] = useState<string | null>(null);
  const [desativando, setDesativando] = useState(false);

  if (mfaAtivo) {
    return (
      <div className="flex flex-col gap-4">
        <Alerta tipo="sucesso">{t('mfa.ativo')}</Alerta>
        {exigidoPelaEmpresa ? (
          <p className="texto-corpo text-texto-secundario">{t('mfa.exigidoPelaEmpresa')}</p>
        ) : desativando ? (
          <>
            <p className="texto-corpo text-texto-secundario">{t('mfa.desativarTexto')}</p>
            <FormularioDeCodigo
              rotuloBotao={t('mfa.desativar')}
              rotuloEnviando={t('mfa.desativando')}
              carregando={desativar.isPending}
              perigo
              erro={erro}
              aoEnviar={(codigo, definirErro) => {
                setErro(null);
                desativar.mutate(codigo, {
                  onSuccess: () => {
                    setDesativando(false);
                    iniciar.reset();
                  },
                  onError: (e) => {
                    const mensagem = mensagemDoErro(e);
                    if (ehErroDaAPI(e) && e.codigo === 'mfa_codigo_invalido') definirErro(mensagem);
                    else setErro(mensagem);
                  },
                });
              }}
            />
          </>
        ) : (
          <div>
            <Botao onClick={() => setDesativando(true)}>{t('mfa.desativar')}</Botao>
          </div>
        )}
      </div>
    );
  }

  if (!iniciar.data) {
    return (
      <div className="flex flex-col gap-4">
        <p className="texto-corpo text-texto-secundario">{t('mfa.inativoTexto')}</p>
        {iniciar.isError && <Alerta tipo="erro">{mensagemDoErro(iniciar.error)}</Alerta>}
        <div>
          <Botao variante="primario" carregando={iniciar.isPending} onClick={() => iniciar.mutate()}>
            {iniciar.isPending ? t('mfa.preparando') : t('mfa.ativar')}
          </Botao>
        </div>
      </div>
    );
  }

  return (
    <div className="flex flex-col gap-4">
      <ol className="flex list-decimal flex-col gap-2 pl-5 texto-corpo">
        <li>{t('mfa.passo1')}</li>
        <li>{t('mfa.passo2')}</li>
      </ol>
      <div className="flex flex-wrap items-start gap-6">
        <CodigoQR conteudo={iniciar.data.uri_otpauth} rotulo={t('mfa.qrRotulo')} />
        <div className="flex flex-col gap-1">
          <p className="texto-legenda text-texto-secundario">{t('mfa.segredoRotulo')}</p>
          <p className="texto-dado-forte break-all select-all">{iniciar.data.segredo_base32}</p>
        </div>
      </div>
      <FormularioDeCodigo
        rotuloBotao={t('mfa.confirmar')}
        rotuloEnviando={t('mfa.confirmando')}
        carregando={confirmar.isPending}
        erro={erro}
        aoEnviar={(codigo, definirErro) => {
          setErro(null);
          confirmar.mutate(codigo, {
            onError: (e) => {
              const mensagem = mensagemDoErro(e);
              if (ehErroDaAPI(e) && e.codigo === 'mfa_codigo_invalido') definirErro(mensagem);
              else setErro(mensagem);
            },
          });
        }}
      />
    </div>
  );
}
