// Tela de login (seção 5.6 do documento 7, RF-ADM-006).
//
// Só o layout: o contrato de autenticação ainda não existe em
// backend/api/openapi.yaml. Os campos aceitam digitação para conferência
// visual e de teclado, mas nada é enviado, e o botão fica desabilitado com o
// motivo ao lado (regra 4.3.7).

import { Eye, EyeOff } from 'lucide-react';
import { useId, useState, type FormEvent } from 'react';
import { useTranslation } from 'react-i18next';

import { Botao } from '../../componentes/Botao';
import { Logo } from '../../componentes/Logo';
import { SeloVersao } from '../../componentes/versao/SeloVersao';
import { useTema } from '../../lib/tema';

const CAMPO =
  'h-10 w-full rounded-md border border-borda-controle bg-superficie px-3 texto-corpo text-texto hover:border-borda-controle-hover focus-visible:border-foco';

export function TelaEntrar() {
  const { t } = useTranslation();
  const { resolvido } = useTema();
  const [senhaVisivel, setSenhaVisivel] = useState(false);
  const id = useId();
  const idEmail = `${id}-email`;
  const idSenha = `${id}-senha`;
  const idMotivo = `${id}-motivo`;

  const naoEnviar = (evento: FormEvent) => evento.preventDefault();

  return (
    <div className="flex min-h-screen flex-col items-center justify-center bg-fundo px-4 pt-16 pb-12 text-texto">
      <main className="flex w-full max-w-100 flex-col items-center gap-8">
        <Logo forma="horizontal" variante={resolvido === 'escuro' ? 'negativo' : 'positivo'} className="h-10 w-auto" />

        <div className="w-full rounded-lg border border-borda bg-superficie p-5 shadow-1">
          <h1 className="texto-h2">{t('entrar.titulo')}</h1>

          <form noValidate onSubmit={naoEnviar} className="mt-6 flex flex-col gap-4">
            <div className="flex flex-col gap-3">
              <label htmlFor={idEmail} className="texto-rotulo">
                {t('entrar.email')}
              </label>
              <input id={idEmail} name="email" type="email" autoComplete="username" className={CAMPO} />
            </div>

            <div className="flex flex-col gap-3">
              <label htmlFor={idSenha} className="texto-rotulo">
                {t('entrar.senha')}
              </label>
              <div className="relative">
                <input
                  id={idSenha}
                  name="senha"
                  type={senhaVisivel ? 'text' : 'password'}
                  autoComplete="current-password"
                  className={`${CAMPO} pr-12`}
                />
                <Botao
                  variante="icone"
                  aria-label={t('entrar.mostrarSenha')}
                  aria-pressed={senhaVisivel}
                  aria-controls={idSenha}
                  onClick={() => setSenhaVisivel((atual) => !atual)}
                  className="absolute top-0 right-0"
                >
                  {senhaVisivel ? <EyeOff size={20} aria-hidden="true" /> : <Eye size={20} aria-hidden="true" />}
                </Botao>
              </div>
            </div>

            <Botao type="submit" variante="primario" disabled aria-describedby={idMotivo} className="mt-2 w-full">
              {t('entrar.botao')}
            </Botao>
            <p id={idMotivo} className="texto-legenda text-texto-secundario">
              {t('entrar.indisponivel')}
            </p>
          </form>
        </div>
      </main>

      <footer>
        <SeloVersao />
      </footer>
    </div>
  );
}
