// Inicialização do i18next com os textos de @tenaz/i18n.
//
// Os recursos vêm embutidos no pacote; a inicialização é síncrona, então
// qualquer módulo que importe este arquivo já pode traduzir.

import { ESPACO_DE_NOMES, IDIOMA_PADRAO, IDIOMAS, recursos, type Idioma, type Textos } from '@tenaz/i18n';
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

import { guardarIdioma, idiomaInicial } from './idioma';

declare module 'i18next' {
  interface CustomTypeOptions {
    defaultNS: typeof ESPACO_DE_NOMES;
    resources: { [ESPACO_DE_NOMES]: Textos };
  }
}

if (!i18n.isInitialized) {
  void i18n.use(initReactI18next).init({
    lng: idiomaInicial(),
    fallbackLng: IDIOMA_PADRAO,
    supportedLngs: [...IDIOMAS],
    resources: recursos,
    ns: [ESPACO_DE_NOMES],
    defaultNS: ESPACO_DE_NOMES,
    initAsync: false,
    // React já escapa o conteúdo; escapar de novo mostraria entidades HTML.
    interpolation: { escapeValue: false },
  });
}

/** Idioma e título do documento acompanham a língua da interface (WCAG 3.1.1). */
function aplicarAoDocumento(idioma: string): void {
  if (typeof document === 'undefined') return;
  document.documentElement.lang = idioma;
  document.title = i18n.t('app.titulo');
}

i18n.on('languageChanged', aplicarAoDocumento);
aplicarAoDocumento(i18n.language);

/** Troca o idioma e guarda a escolha no navegador. */
export async function definirIdioma(idioma: Idioma): Promise<void> {
  guardarIdioma(idioma);
  await i18n.changeLanguage(idioma);
}

export { i18n };

/** Chave de texto conhecida pelo pacote de i18n? */
export function existeTexto(chave: string): boolean {
  return i18n.exists(chave);
}

/**
 * Traduz uma chave montada em tempo de execução (um código de erro do
 * servidor, uma mensagem de esquema de validação). Chave desconhecida volta
 * como veio: os testes e o tipo dos esquemas impedem que isso chegue ao usuário.
 */
export function traduzirChave(chave: string, opcoes?: Record<string, unknown>): string {
  // A chave é montada em tempo de execução; o tipo estrito das chaves não se aplica.
  return String((i18n.t as unknown as (chave: string, opcoes?: Record<string, unknown>) => unknown)(chave, opcoes));
}

/** Texto de um valor que o contrato não fecha em enumeração (a ação de um evento): sem texto conhecido, o próprio valor. */
export function traduzirOuOriginal(prefixo: string, valor: string): string {
  const chave = `${prefixo}.${valor}`;
  return i18n.exists(chave) ? traduzirChave(chave) : valor;
}
