// Painel (IND). Nesta etapa, só o cabeçalho e o estado vazio: os indicadores
// chegam com os módulos de ativos e de ordens de serviço. Para quem entra como
// prestador ou locatário, avisa que o portal próprio chega em etapa futura.

import { useTranslation } from 'react-i18next';

import { Alerta } from '../../componentes/Alerta';
import { useSessao } from '../adm/sessao';

export function TelaPainel() {
  const { t } = useTranslation();
  const { data: sessao } = useSessao();
  const externo = sessao != null && sessao.escopo_acesso !== 'interno';
  return (
    <>
      <h1 className="texto-h2">{t('painel.titulo')}</h1>
      {externo && (
        <Alerta tipo="informacao" className="mt-4 max-w-180">
          {t('painel.portalEmBreve')}
        </Alerta>
      )}
      <p className="mt-2 max-w-180 texto-corpo text-texto-secundario">{t('painel.vazio')}</p>
    </>
  );
}
