// Ponto de entrada da web.
//
// Nesta etapa (E00) a web existe só para provar a esteira de versão: a
// interface de gestão começa em E01, com React, os tokens da marca e o selo de
// versão na tela de login e na tela principal (RF-ADM-001, RNF-039).

import { seloDeVersao, versaoDaWeb } from './versao';

const raiz = document.querySelector('#raiz');
if (raiz) {
  raiz.textContent = seloDeVersao(versaoDaWeb);
}
