// Escolha do idioma da interface (documento 7, 10.7): o nome de cada idioma no
// próprio idioma, sem bandeira. Dois formatos: lista de opções em Minha conta e
// link de texto no pé do cartão das telas sem sessão. A troca passa por
// `useEscolherIdioma` (modulos/adm/idioma.ts): guarda no navegador e, havendo
// sessão, grava em `usuario.idioma` (PATCH /v1/usuarios/eu).

import { IDIOMAS, ehIdioma, type Idioma } from '@tenaz/i18n';
import { useTranslation } from 'react-i18next';

import { useEscolherIdioma } from '../modulos/adm/idioma';

import { Botao } from './Botao';

/** Opções de rádio, para Minha conta. */
export function SeletorDeIdioma() {
  const { t, i18n } = useTranslation();
  const escolher = useEscolherIdioma();
  return (
    <fieldset className="flex flex-col gap-2">
      <legend className="mb-1 texto-corpo text-texto-secundario">{t('idioma.texto')}</legend>
      {IDIOMAS.map((idioma) => (
        <label key={idioma} lang={idioma} className="flex min-h-8 cursor-pointer items-center gap-3 texto-corpo">
          <input
            type="radio"
            name="idioma"
            value={idioma}
            checked={i18n.language === idioma}
            onChange={(e) => {
              if (ehIdioma(e.target.value)) void escolher(e.target.value);
            }}
            className="size-5 accent-texto"
          />
          {t(`idiomas.${idioma}`)}
        </label>
      ))}
    </fieldset>
  );
}

/** Links de texto no pé do cartão das telas sem sessão. */
export function LinksDeIdioma() {
  const { t, i18n } = useTranslation();
  const escolher = useEscolherIdioma();
  return (
    <nav aria-label={t('idioma.rodape')} className="flex justify-center gap-4">
      {IDIOMAS.map((idioma: Idioma) => (
        <Botao
          key={idioma}
          variante="terciario"
          lang={idioma}
          aria-current={i18n.language === idioma ? 'true' : undefined}
          onClick={() => void escolher(idioma)}
          className={i18n.language === idioma ? 'underline' : ''}
        >
          {t(`idiomas.${idioma}`)}
        </Botao>
      ))}
    </nav>
  );
}
