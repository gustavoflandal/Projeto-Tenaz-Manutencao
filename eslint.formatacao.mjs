// Trava de formatação para pessoas (ADR-013, Parte G, item 4, documento 5):
// número, dinheiro, lista, tempo relativo e data só pelas funções de
// `pacotes/dominio/src/formatos.ts`. Este arquivo não tem dependência: a web, o
// app e cada pacote o importam na sua configuração do ESLint.
//
// Barrados: Intl.NumberFormat, Intl.RelativeTimeFormat, Intl.ListFormat,
// Intl.DateTimeFormat e toLocaleString, toLocaleDateString, toLocaleTimeString.
// Permitidos: Intl.Collator (ordenação) e Intl.PluralRules (plural).
// A única exceção é o Intl.DateTimeFormat de formatos.ts, motor de fuso, onde
// valem as regras de `regrasDoMotorDeFuso`.

const INTL_DE_TEXTO = 'NumberFormat|RelativeTimeFormat|ListFormat|DurationFormat|DateTimeFormat';
const TO_LOCALE = 'toLocaleString|toLocaleDateString|toLocaleTimeString';

const MENSAGEM_INTL =
  'Formatação para pessoas só por pacotes/dominio/src/formatos.ts (ADR-013, Parte G): sem Intl gerando texto.';
const MENSAGEM_TO_LOCALE = 'toLocaleString, toLocaleDateString e toLocaleTimeString são proibidos (ADR-013, Parte G): use formatos.ts.';

const seletorDeIntl = (nomes) =>
  `MemberExpression[object.name="Intl"][property.name=/^(${nomes})$/], MemberExpression[object.object.name="globalThis"][object.property.name="Intl"][property.name=/^(${nomes})$/]`;

const seletorDeToLocale = `CallExpression[callee.property.name=/^(${TO_LOCALE})$/]`;

/** Regras do código de interface: nada de Intl gerando texto nem toLocale*. */
export const regrasDeFormatacao = {
  'no-restricted-syntax': [
    'error',
    { selector: seletorDeIntl(INTL_DE_TEXTO), message: MENSAGEM_INTL },
    { selector: seletorDeToLocale, message: MENSAGEM_TO_LOCALE },
  ],
};

/** Regras de formatos.ts: o Intl.DateTimeFormat vale só como motor de fuso, sem `format` nem `formatRange`. */
export const regrasDoMotorDeFuso = {
  'no-restricted-syntax': [
    'error',
    { selector: seletorDeIntl('NumberFormat|RelativeTimeFormat|ListFormat|DurationFormat'), message: MENSAGEM_INTL },
    { selector: seletorDeToLocale, message: MENSAGEM_TO_LOCALE },
    {
      selector: 'CallExpression[callee.property.name=/^(format|formatRange|formatRangeToParts)$/]',
      message: 'O motor de fuso usa só formatToParts: o texto é montado por formatos.ts, nunca pelo Intl.',
    },
  ],
};
