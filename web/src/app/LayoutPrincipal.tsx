// Casca da tela principal (seção 5.1 do documento 7): menu lateral, barra
// superior, área de conteúdo e o selo de versão no canto inferior direito.

import { useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { Outlet } from 'react-router';

import { MenuTema } from '../componentes/MenuTema';
import type { Sessao } from '../lib/tipos';
import { useSessao } from '../modulos/adm/sessao';
import { SeloVersao } from '../componentes/versao/SeloVersao';
import { useAcimaDe } from '../lib/breakpoints';

import { BarraLateral } from './BarraLateral';
import { MenuUsuario } from './MenuUsuario';

const CHAVE_MENU = 'tenaz.menu';

type PreferenciaMenu = 'expandido' | 'recolhido';

function lerPreferenciaMenu(): PreferenciaMenu | null {
  try {
    const salvo = globalThis.localStorage.getItem(CHAVE_MENU);
    return salvo === 'expandido' || salvo === 'recolhido' ? salvo : null;
  } catch {
    return null;
  }
}

/**
 * Menu expandido ou recolhido (seções 3.10 e 5.2 do documento 7): sem escolha
 * do usuário, expandido a partir de 1280 px e recolhido abaixo. Abaixo de
 * 768 px fica sempre recolhido, até a gaveta existir.
 */
function useMenu(): { readonly expandido: boolean; readonly alternavel: boolean; readonly alternar: () => void } {
  const acimaDeMd = useAcimaDe('md');
  const acimaDeXl = useAcimaDe('xl');
  const [preferencia, setPreferencia] = useState(lerPreferenciaMenu);
  const expandido = acimaDeMd && (preferencia === null ? acimaDeXl : preferencia === 'expandido');

  const alternar = useCallback(() => {
    const nova: PreferenciaMenu = expandido ? 'recolhido' : 'expandido';
    setPreferencia(nova);
    try {
      globalThis.localStorage.setItem(CHAVE_MENU, nova);
    } catch {
      // Sem armazenamento, a escolha vale até recarregar.
    }
  }, [expandido]);

  return { expandido, alternavel: acimaDeMd, alternar };
}

export function LayoutPrincipal() {
  const { data: sessao } = useSessao();
  // A guarda de rotas só monta a casca com sessão; sem ela, não há o que mostrar.
  return sessao ? <Casca sessao={sessao} /> : null;
}

function Casca({ sessao }: { readonly sessao: Sessao }) {
  const { t } = useTranslation();
  const { expandido, alternavel, alternar } = useMenu();

  return (
    <div className="flex min-h-screen bg-fundo text-texto">
      <a
        href="#conteudo"
        className="sr-only rounded-md bg-superficie px-4 py-2 texto-rotulo text-texto shadow-3 focus:not-sr-only focus:fixed focus:top-2 focus:left-2 focus:z-(--tz-z-dica)"
      >
        {t('layout.pularConteudo')}
      </a>

      <BarraLateral sessao={sessao} expandido={expandido} aoAlternar={alternavel ? alternar : undefined} />

      <div className="flex min-w-0 flex-1 flex-col">
        <header className="sticky top-0 z-(--tz-z-nav) flex h-14 shrink-0 items-center justify-between gap-2 border-b border-borda bg-superficie px-6">
          <p className="min-w-0 truncate texto-rotulo">{sessao.empresa.nome}</p>
          <div className="flex items-center gap-2">
            <MenuTema />
            <MenuUsuario sessao={sessao} />
          </div>
        </header>

        {/* Respiro inferior de 32 px para o selo nunca cobrir conteúdo (6.19). */}
        <main id="conteudo" tabIndex={-1} className="flex-1 px-6 pt-6 pb-8">
          <Outlet />
        </main>
      </div>

      <footer>
        <SeloVersao />
      </footer>
    </div>
  );
}
