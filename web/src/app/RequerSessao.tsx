// Guarda das rotas protegidas. Quem não tem sessão vai para o login e volta
// depois ao endereço pedido; quem tem a sessão restrita à ativação do MFA
// (`mfa_pendente`, RN-ADM-027) só alcança a tela de ativação.
//
// Esconder tela não é segurança: o servidor responde 401, 403 ou 404 em cada
// rota, e a web trata essas respostas. A guarda só evita telas que não podem abrir.

import { useQueryClient } from '@tanstack/react-query';
import { useTranslation } from 'react-i18next';
import { Navigate, Outlet, useLocation } from 'react-router';

import { Carregando, EstadoErro } from '../componentes/Estados';
import { chaves } from '../lib/consultas';
import { useAplicarIdiomaDoUsuario } from '../modulos/adm/idioma';
import { useSessao } from '../modulos/adm/sessao';

export interface RequerSessaoProps {
  /** Só para a tela de ativação do MFA: aceita a sessão restrita e recusa a plena. */
  readonly somenteRestrita?: boolean;
}

export function RequerSessao({ somenteRestrita = false }: RequerSessaoProps) {
  const { t } = useTranslation();
  const local = useLocation();
  const cliente = useQueryClient();
  const sessao = useSessao();

  if (sessao.isPending) {
    return (
      <div className="min-h-screen bg-fundo p-6">
        <Carregando rotulo={t('sessao.carregando')} linhas={3} />
      </div>
    );
  }

  if (sessao.isError) {
    return (
      <div className="min-h-screen bg-fundo text-texto">
        <EstadoErro erro={sessao.error} aoTentarDeNovo={() => void sessao.refetch()} />
      </div>
    );
  }

  if (sessao.data === null) {
    const retorno = encodeURIComponent(`${local.pathname}${local.search}`);
    const encerrada = cliente.getQueryData(chaves.sessaoEncerrada) === true;
    const destino = `/entrar?retorno=${retorno}${encerrada ? '&motivo=sessao' : ''}`;
    return <Navigate to={destino} replace />;
  }

  if (sessao.data.mfa_pendente && !somenteRestrita) return <Navigate to="/ativar-mfa" replace />;
  if (!sessao.data.mfa_pendente && somenteRestrita) return <Navigate to="/" replace />;

  return somenteRestrita ? <Outlet /> : <RotaComIdioma />;
}

/** Rota liberada: o idioma do usuário vale desde a primeira tela (a sessão restrita não lê o usuário). */
function RotaComIdioma() {
  useAplicarIdiomaDoUsuario();
  return <Outlet />;
}
