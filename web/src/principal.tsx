// Ponto de entrada da aplicação principal da web.

import './app/estilos.css';

import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';

import { Aplicacao } from './app/Aplicacao';
import { i18n } from './lib/i18n';
import { iniciarTema } from './lib/tema';

// O tema é aplicado antes da primeira renderização, para não piscar.
iniciarTema();
document.title = i18n.t('app.titulo');

const raiz = document.querySelector('#raiz');
if (raiz) {
  createRoot(raiz).render(
    <StrictMode>
      <Aplicacao />
    </StrictMode>,
  );
}
