// Idioma da interface: o do aparelho na primeira vez; depois, o que a pessoa
// escolher, guardado localmente até a API ter `usuario.idioma` (pendência 13.40
// de docs/7_PADROES_DESIGN.md). Vale antes do login, por isso não depende de sessão.

import { getLocales } from 'expo-localization';
import { createContext, type ReactNode, useCallback, useContext, useEffect, useMemo, useState } from 'react';

import type { Cofre } from '../acesso/cofre';
import { CofreExpo } from '../acesso/cofre-expo';
import { escolherIdioma, IDIOMAS, type Idioma, tradutor, type Traduzir } from './index';

const CHAVE = 'tenaz.idioma';
const cofre: Cofre = new CofreExpo();

interface Valor {
  idioma: Idioma;
  idiomas: readonly Idioma[];
  t: Traduzir;
  trocarIdioma: (novo: Idioma) => Promise<void>;
}

const Contexto = createContext<Valor | null>(null);

function tagsDoAparelho(): string[] {
  return getLocales().map((l) => l.languageTag);
}

export function ProvedorDeIdioma({ children }: { children: ReactNode }) {
  const [idioma, setIdioma] = useState<Idioma>(() => escolherIdioma(null, tagsDoAparelho()));

  useEffect(() => {
    void cofre.ler(CHAVE).then((guardado) => setIdioma(escolherIdioma(guardado, tagsDoAparelho())));
  }, []);

  const trocarIdioma = useCallback(async (novo: Idioma) => {
    setIdioma(novo);
    await cofre.gravar(CHAVE, novo);
  }, []);

  const valor = useMemo<Valor>(
    () => ({ idioma, idiomas: IDIOMAS, t: tradutor(idioma), trocarIdioma }),
    [idioma, trocarIdioma],
  );
  return <Contexto.Provider value={valor}>{children}</Contexto.Provider>;
}

export function useTraducao(): Valor {
  const v = useContext(Contexto);
  if (!v) throw new Error('useTraducao fora do ProvedorDeIdioma');
  return v;
}
