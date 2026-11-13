// Painel (IND). Nesta etapa, só o cabeçalho e o estado vazio: os indicadores
// chegam com os módulos de ativos e de ordens de serviço.

import { useTranslation } from 'react-i18next';

export function TelaPainel() {
  const { t } = useTranslation();
  return (
    <>
      <h1 className="texto-h2">{t('painel.titulo')}</h1>
      <p className="mt-2 max-w-180 texto-corpo text-texto-secundario">{t('painel.vazio')}</p>
    </>
  );
}
