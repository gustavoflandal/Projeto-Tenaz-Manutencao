// Menu lateral (seção 5.2 do documento 7): fundo Aço, 240 px expandido e
// 64 px recolhido. Itens aparecem conforme perfil e plano; enquanto a
// autorização não existe (E01), só o Painel.

import { ChartLine, ChevronsLeft, ChevronsRight, type LucideIcon } from 'lucide-react';
import { Link, NavLink } from 'react-router';
import { useTranslation } from 'react-i18next';

import { Logo } from '../componentes/Logo';

interface ItemDoMenu {
  readonly rota: string;
  readonly rotulo: 'menu.painel';
  readonly icone: LucideIcon;
}

const ITENS: readonly ItemDoMenu[] = [{ rota: '/', rotulo: 'menu.painel', icone: ChartLine }];

// Sobre o menu (superfície Aço), o anel de foco é Névoa (11.1).
const FOCO_SOBRE_ACO = 'focus-visible:outline-nav-texto';

export interface BarraLateralProps {
  readonly expandido: boolean;
  /** Ausente quando a largura da tela não permite expandir o menu. */
  readonly aoAlternar?: () => void;
}

export function BarraLateral({ expandido, aoAlternar }: BarraLateralProps) {
  const { t } = useTranslation();
  const IconeAlternar = expandido ? ChevronsLeft : ChevronsRight;

  return (
    <nav
      id="menu-principal"
      aria-label={t('layout.menuPrincipal')}
      className={`sticky top-0 flex h-screen shrink-0 flex-col bg-nav-fundo text-nav-texto ${expandido ? 'w-60' : 'w-16'}`}
    >
      <div className={`flex h-14 shrink-0 items-center ${expandido ? 'px-4' : 'justify-center'}`}>
        <Link to="/" className={`inline-flex rounded-xs ${FOCO_SOBRE_ACO}`}>
          {expandido ? (
            <Logo forma="horizontal" variante="negativo" className="h-6 w-auto" />
          ) : (
            <Logo forma="simbolo" variante="negativo" className="size-8" />
          )}
        </Link>
      </div>

      <div className="flex-1 overflow-y-auto px-2 py-2">
        <ul className="flex flex-col gap-1">
          {ITENS.map(({ rota, rotulo, icone: Icone }) => (
            <li key={rota}>
              <NavLink
                to={rota}
                end
                aria-label={expandido ? undefined : t(rotulo)}
                className={({ isActive }) =>
                  [
                    'flex h-10 items-center gap-3 rounded-md border-l-3 texto-corpo-p',
                    expandido ? 'px-3' : 'justify-center',
                    isActive
                      ? 'border-nav-texto bg-nav-item-ativo text-nav-texto'
                      : 'border-transparent text-nav-texto hover:bg-nav-item-ativo',
                    FOCO_SOBRE_ACO,
                  ].join(' ')
                }
              >
                <Icone size={20} aria-hidden="true" className="shrink-0 text-nav-texto-secundario" />
                {expandido && <span>{t(rotulo)}</span>}
              </NavLink>
            </li>
          ))}
        </ul>
      </div>

      {aoAlternar && (
        <div
          className={`flex shrink-0 border-t border-nav-item-ativo p-2 ${expandido ? 'justify-end' : 'justify-center'}`}
        >
          <button
            type="button"
            onClick={aoAlternar}
            aria-expanded={expandido}
            aria-controls="menu-principal"
            aria-label={expandido ? t('layout.recolherMenu') : t('layout.expandirMenu')}
            className={`inline-flex size-10 items-center justify-center rounded-md text-nav-texto-secundario hover:bg-nav-item-ativo hover:text-nav-texto ${FOCO_SOBRE_ACO}`}
          >
            <IconeAlternar size={20} aria-hidden="true" />
          </button>
        </div>
      )}
    </nav>
  );
}
