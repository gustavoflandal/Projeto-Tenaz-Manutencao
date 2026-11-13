// Cabeçalho da página (seção 5.4 do documento 7): trilha de navegação, título
// em h1 visual `h2`, metadados e ações à direita (uma primária, no máximo).

import { ChevronRight } from 'lucide-react';
import type { ReactNode } from 'react';
import { useTranslation } from 'react-i18next';
import { Link } from 'react-router';

export interface MigalhaDePao {
  readonly rotulo: string;
  readonly para?: string;
}

export interface CabecalhoPaginaProps {
  readonly titulo: string;
  readonly trilha?: readonly MigalhaDePao[];
  readonly metadados?: ReactNode;
  readonly acoes?: ReactNode;
}

export function CabecalhoPagina({ titulo, trilha, metadados, acoes }: CabecalhoPaginaProps) {
  const { t } = useTranslation();
  return (
    <header className="mb-6 flex flex-col gap-2">
      {trilha && trilha.length > 0 && (
        <nav aria-label={t('layout.trilha')}>
          <ol className="flex flex-wrap items-center gap-1 texto-legenda text-texto-secundario">
            {trilha.map((item, i) => (
              <li key={item.rotulo} className="flex items-center gap-1">
                {i > 0 && <ChevronRight size={14} aria-hidden="true" />}
                {item.para ? (
                  <Link to={item.para} className="underline">
                    {item.rotulo}
                  </Link>
                ) : (
                  <span aria-current="page">{item.rotulo}</span>
                )}
              </li>
            ))}
          </ol>
        </nav>
      )}
      <div className="flex flex-wrap items-center justify-between gap-3">
        <h1 className="texto-h2">{titulo}</h1>
        {acoes && <div className="flex flex-wrap items-center justify-end gap-2">{acoes}</div>}
      </div>
      {metadados && <div className="flex flex-wrap items-center gap-3 texto-corpo-p text-texto-secundario">{metadados}</div>}
    </header>
  );
}
