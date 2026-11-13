import { render, screen } from '@testing-library/react';
import { afterEach, describe, expect, it } from 'vitest';

import { i18n } from '../lib/i18n';

import { AvisoDadoDeSaude } from './AvisoDadoDeSaude';

describe('aviso de dado de saúde', () => {
  afterEach(async () => {
    await i18n.changeLanguage('pt-BR');
  });

  it('aparece em pt-BR', async () => {
    await i18n.changeLanguage('pt-BR');
    render(<AvisoDadoDeSaude />);
    expect(screen.getByText(/não guarda dado de saúde/i)).toBeInTheDocument();
  });

  it('não renderiza nada em es-419', async () => {
    await i18n.changeLanguage('es-419');
    const { container } = render(<AvisoDadoDeSaude />);
    expect(container).toBeEmptyDOMElement();
  });
});
