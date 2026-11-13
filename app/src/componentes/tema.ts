// Cores do tema claro, lidas de pacotes/tokens-ui/tokens.json: a fonte única
// dos tokens (docs/7_PADROES_DESIGN.md, seção 3). Nada de cor transcrita à mão.
// Texto sobre Sinal é sempre Aço. Quando o app ganhar o tema escuro, a troca é
// escolher `tokens.temas.escuro` aqui.

import tokens from '@tenaz/tokens-ui/tokens.json';

const claro = tokens.temas.claro;

export const cor = {
  fundo: claro.fundo,
  superficie: claro.superficie,
  superficieSelecionada: claro['superficie-selecionada'],
  texto: claro.texto,
  textoSecundario: claro['texto-secundario'],
  textoTerciario: claro['texto-terciario'],
  textoDesabilitado: claro['texto-desabilitado'],
  borda: claro.borda,
  bordaControle: claro['borda-controle'],
  foco: claro.foco,
  acaoPrimaria: claro['acao-primaria'],
  acaoPrimariaPressionada: claro['acao-primaria-pressionada'],
  acaoPrimariaTexto: claro['acao-primaria-texto'],
  perigo: claro.perigo,
  perigoPressionado: claro['perigo-hover'],
  perigoTexto: claro['perigo-texto'],
  atencao: claro['status-atencao'],
  atencaoFundo: claro['status-atencao-fundo'],
  atencaoTexto: claro['status-atencao-texto'],
  navFundo: claro['nav-fundo'],
  navTexto: claro['nav-texto'],
} as const;

/** Alvo de toque mínimo no app: 48 dp. A ação primária tem 56 dp (9.2). */
export const ALVO_DE_TOQUE = 48;
export const ALVO_PRIMARIO = 56;
