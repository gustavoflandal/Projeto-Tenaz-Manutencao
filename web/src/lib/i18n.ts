// Inicialização do i18next com os textos de @tenaz/i18n.
//
// Os recursos vêm embutidos no pacote; a inicialização é síncrona, então
// qualquer módulo que importe este arquivo já pode traduzir.

import { ESPACO_DE_NOMES, IDIOMA_PADRAO, recursos, type Textos } from '@tenaz/i18n';
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

declare module 'i18next' {
  interface CustomTypeOptions {
    defaultNS: typeof ESPACO_DE_NOMES;
    resources: { [ESPACO_DE_NOMES]: Textos };
  }
}

if (!i18n.isInitialized) {
  void i18n.use(initReactI18next).init({
    lng: IDIOMA_PADRAO,
    fallbackLng: IDIOMA_PADRAO,
    resources: recursos,
    ns: [ESPACO_DE_NOMES],
    defaultNS: ESPACO_DE_NOMES,
    initAsync: false,
    // React já escapa o conteúdo; escapar de novo mostraria entidades HTML.
    interpolation: { escapeValue: false },
  });
}

export { i18n };
