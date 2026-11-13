// Aviso fixo sobre dado de saúde (RF-ADM-015 e sumário 3.6): aparece em toda
// carga das telas de criar, editar e detalhar pessoa, para qualquer perfil.
// Não tem botão de fechar e nenhum parâmetro o desliga. Só existe em pt-BR:
// em espanhol o componente não renderiza nada (decisão do dono do produto;
// fica como customização, se o cliente pedir).

import { useTranslation } from 'react-i18next';

import { Alerta } from './Alerta';

export function AvisoDadoDeSaude() {
  const { t, i18n } = useTranslation();
  if (i18n.language !== 'pt-BR') return null;
  return (
    <Alerta tipo="informacao" className="max-w-180">
      {t('pessoa.avisoSaude')}
    </Alerta>
  );
}
