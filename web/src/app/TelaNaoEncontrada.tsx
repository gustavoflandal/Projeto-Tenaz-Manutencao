import { Link } from 'react-router';
import { useTranslation } from 'react-i18next';

export function TelaNaoEncontrada() {
  const { t } = useTranslation();
  return (
    <>
      <h1 className="texto-h2">{t('naoEncontrada.titulo')}</h1>
      <p className="mt-2 texto-corpo text-texto-secundario">{t('naoEncontrada.texto')}</p>
      <Link to="/" className="mt-4 inline-flex texto-botao text-texto underline">
        {t('naoEncontrada.acao')}
      </Link>
    </>
  );
}
