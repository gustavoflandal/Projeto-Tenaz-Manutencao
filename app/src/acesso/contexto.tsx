// Estado de acesso para as telas: sessão atual e bloqueio local.

import { createContext, type ReactNode, useCallback, useContext, useEffect, useState } from 'react';
import { AppState } from 'react-native';

import { bloqueio, gerente } from './nativo';
import type { EstadoDeAcesso } from './sessoes';

interface ValorDoContexto {
  estado: EstadoDeAcesso;
  /** O bloqueio local está pedindo biometria ou PIN. Só esconde a interface; a fila e a sincronização seguem. */
  bloqueado: boolean;
  desbloquear: (motivo: string) => Promise<boolean>;
  bloqueioAtivo: boolean;
  alternarBloqueio: (motivoAtivar: string) => Promise<boolean>;
}

const Contexto = createContext<ValorDoContexto | null>(null);

export function ProvedorDeAcesso({ children }: { children: ReactNode }) {
  const [estado, setEstado] = useState<EstadoDeAcesso>(gerente.estado());
  const [bloqueado, setBloqueado] = useState(false);
  const [bloqueioAtivo, setBloqueioAtivo] = useState(false);

  useEffect(() => {
    const cancelar = gerente.assinar(setEstado);
    void (async () => {
      await bloqueio.iniciar();
      await gerente.carregar();
      setBloqueado(bloqueio.estaBloqueado());
      setBloqueioAtivo(bloqueio.ativo());
      // Sair sem rede deixa a revogação guardada; repete quando houver sinal.
      gerente.tentarRevogacoesPendentes().catch(() => undefined);
    })();
    const ouvinte = AppState.addEventListener('change', (s) => {
      if (s === 'active') {
        bloqueio.aoVoltar();
        setBloqueado(bloqueio.estaBloqueado());
      } else {
        bloqueio.aoIrParaSegundoPlano();
      }
    });
    return () => {
      cancelar();
      ouvinte.remove();
    };
  }, []);

  const desbloquear = useCallback(async (motivo: string) => {
    const ok = await bloqueio.desbloquear(motivo);
    setBloqueado(bloqueio.estaBloqueado());
    return ok;
  }, []);

  const alternarBloqueio = useCallback(async (motivoAtivar: string) => {
    if (bloqueio.ativo()) await bloqueio.desativar();
    else if (!(await bloqueio.ativar(motivoAtivar))) return false;
    setBloqueioAtivo(bloqueio.ativo());
    setBloqueado(false);
    return true;
  }, []);

  return (
    <Contexto.Provider value={{ estado, bloqueado, desbloquear, bloqueioAtivo, alternarBloqueio }}>
      {children}
    </Contexto.Provider>
  );
}

export function useAcesso(): ValorDoContexto {
  const v = useContext(Contexto);
  if (!v) throw new Error('useAcesso fora do ProvedorDeAcesso');
  return v;
}
