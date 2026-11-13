// Busca por texto de uma lista: aplica ao enviar (Enter ou botão), e não a cada
// tecla, para não pedir ao servidor uma página por letra digitada.

import { Search } from 'lucide-react';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';

import { Botao } from './Botao';
import { CampoTexto } from './Campo';

export interface BuscaPorTextoProps {
  readonly rotulo: string;
  readonly valor: string;
  readonly aoBuscar: (valor: string) => void;
}

export function BuscaPorTexto({ rotulo, valor, aoBuscar }: BuscaPorTextoProps) {
  const { t } = useTranslation();
  // Mantém o que a pessoa digita; volta ao valor da URL quando ele muda por fora.
  const [digitado, setDigitado] = useState(valor);
  const [anterior, setAnterior] = useState(valor);
  if (valor !== anterior) {
    setAnterior(valor);
    setDigitado(valor);
  }

  const buscar = () => aoBuscar(digitado.trim());

  // Sem <form>: a busca também aparece dentro de formulários, e form dentro de form não vale.
  return (
    <div role="search" className="flex items-end gap-2">
      <div className="w-64">
        <CampoTexto
          rotulo={rotulo}
          type="search"
          value={digitado}
          onChange={(e) => setDigitado(e.target.value)}
          onKeyDown={(e) => {
            if (e.key === 'Enter') {
              e.preventDefault();
              buscar();
            }
          }}
          name="busca"
        />
      </div>
      <Botao onClick={buscar} aria-label={t('lista.buscar')}>
        <Search size={20} aria-hidden="true" />
        {t('lista.buscar')}
      </Botao>
    </div>
  );
}
