// Menu lateral (seção 5.2 do documento 7): fundo Aço, 240 px expandido e
// 64 px recolhido. Itens aparecem conforme perfil e plano; enquanto a
// autorização não existe (E01), só o Painel.

import {
  Building2,
  ChartLine,
  ChevronsLeft,
  ChevronsRight,
  Contact,
  History,
  MapPin,
  ShieldCheck,
  ShieldEllipsis,
  SlidersHorizontal,
  Users,
  type LucideIcon,
} from 'lucide-react';
import type { Textos } from '@tenaz/i18n';
import { Tooltip } from 'radix-ui';
import { Link, NavLink } from 'react-router';
import { useTranslation } from 'react-i18next';

import { Logo } from '../componentes/Logo';
import type { Sessao } from '../lib/tipos';
import { temPermissao } from '../modulos/adm/sessao';

type ChaveDoMenu = `menu.${Exclude<keyof Textos['menu'], 'emBreve' | 'emBreveNome' | 'administracao'>}`;

interface ItemDoMenu {
  readonly rota: string;
  readonly rotulo: ChaveDoMenu;
  readonly icone: LucideIcon;
  /** Permissão que mostra o item. Item sem permissão não aparece (seção 5.2). */
  readonly permissao?: string;
}

interface GrupoDoMenu {
  readonly titulo?: 'menu.administracao';
  readonly itens: readonly ItemDoMenu[];
}

const GRUPOS: readonly GrupoDoMenu[] = [
  { itens: [{ rota: '/', rotulo: 'menu.painel', icone: ChartLine }] },
  {
    titulo: 'menu.administracao',
    itens: [
      { rota: '/administracao/empresa', rotulo: 'menu.empresa', icone: Building2, permissao: 'adm.ler' },
      { rota: '/administracao/unidades', rotulo: 'menu.unidades', icone: MapPin, permissao: 'adm.ler' },
      { rota: '/administracao/pessoas', rotulo: 'menu.pessoas', icone: Contact, permissao: 'adm.ler_pessoas' },
      { rota: '/administracao/usuarios', rotulo: 'menu.usuarios', icone: Users, permissao: 'adm.ler_usuarios' },
      { rota: '/administracao/perfis', rotulo: 'menu.perfis', icone: ShieldCheck, permissao: 'adm.ler_usuarios' },
      { rota: '/administracao/auditoria', rotulo: 'menu.auditoria', icone: History, permissao: 'adm.ler_auditoria' },
      {
        rota: '/administracao/auditoria/integridade',
        rotulo: 'menu.integridade',
        icone: ShieldEllipsis,
        permissao: 'adm.ler_auditoria',
      },
      {
        rota: '/administracao/auditoria/capturas',
        rotulo: 'menu.capturas',
        icone: SlidersHorizontal,
        permissao: 'adm.gerenciar_parametros',
      },
    ],
  },
];

// Sobre o menu (superfície Aço), o anel de foco é Névoa (11.1).
const FOCO_SOBRE_ACO = 'focus-visible:outline-nav-texto';

export interface BarraLateralProps {
  readonly sessao: Sessao;
  readonly expandido: boolean;
  /** Ausente quando a largura da tela não permite expandir o menu. */
  readonly aoAlternar?: () => void;
}

/**
 * Item do menu. Recolhido, mostra só o ícone, com o nome em tooltip (abre no
 * hover após 400 ms e no foco, componente 6.18) e em `aria-label`.
 */
function ItemAtivo({ item, expandido }: { readonly item: ItemDoMenu; readonly expandido: boolean }) {
  const { t } = useTranslation();
  const Icone = item.icone;
  const ligacao = (
    <NavLink
      to={item.rota}
      end
      aria-label={expandido ? undefined : t(item.rotulo)}
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
      {expandido && <span>{t(item.rotulo)}</span>}
    </NavLink>
  );
  if (expandido) return ligacao;
  return (
    <Tooltip.Root>
      <Tooltip.Trigger asChild>{ligacao}</Tooltip.Trigger>
      <Tooltip.Portal>
        <Tooltip.Content
          side="right"
          sideOffset={8}
          className="z-(--tz-z-dica) rounded-sm bg-nav-fundo px-2 py-1 texto-legenda text-nav-texto shadow-3"
        >
          {t(item.rotulo)}
        </Tooltip.Content>
      </Tooltip.Portal>
    </Tooltip.Root>
  );
}

export function BarraLateral({ sessao, expandido, aoAlternar }: BarraLateralProps) {
  const { t } = useTranslation();
  const IconeAlternar = expandido ? ChevronsLeft : ChevronsRight;
  const grupos = GRUPOS.map((grupo) => ({
    ...grupo,
    itens: grupo.itens.filter((item) => item.permissao === undefined || temPermissao(sessao, item.permissao)),
  })).filter((grupo) => grupo.itens.length > 0);

  return (
    <Tooltip.Provider delayDuration={400} skipDelayDuration={0}>
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

      <div className="flex flex-1 flex-col gap-4 overflow-y-auto px-2 py-2">
        {grupos.map((grupo, indice) => (
          <div key={grupo.titulo ?? indice} role="group" aria-label={grupo.titulo ? t(grupo.titulo) : undefined}>
            {grupo.titulo && expandido && (
              <p aria-hidden="true" className="mb-1 px-3 texto-sobrelinha text-nav-texto-secundario">
                {t(grupo.titulo)}
              </p>
            )}
            <ul className="flex flex-col gap-1">
              {grupo.itens.map((item) => (
                <li key={item.rota}>
                  <ItemAtivo item={item} expandido={expandido} />
                </li>
              ))}
            </ul>
          </div>
        ))}
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
    </Tooltip.Provider>
  );
}
