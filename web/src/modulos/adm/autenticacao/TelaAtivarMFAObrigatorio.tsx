// Sessão restrita à ativação do segundo fator (RN-ADM-027, documento 3, seção
// 4.2). A empresa exige MFA e a pessoa ainda não o ativou: a sessão só serve
// para ler a sessão, sair, trocar de empresa e ativar o MFA. Confirmado o
// código, a mesma sessão volta a valer, sem novo login.

import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router';

import { Alerta } from '../../../componentes/Alerta';
import { Botao } from '../../../componentes/Botao';
import { useSair, useSessao } from '../sessao';

import { CartaoPublico } from './CartaoPublico';
import { EscolhaDeEmpresa } from './EscolhaDeEmpresa';
import { SecaoMFA } from './SecaoMFA';

export function TelaAtivarMFAObrigatorio() {
  const { t } = useTranslation();
  const navegar = useNavigate();
  const { data: sessao } = useSessao();
  const sair = useSair();
  if (!sessao) return null;

  return (
    <CartaoPublico titulo={t('mfaObrigatorio.titulo')}>
      <Alerta tipo="atencao">{t('mfaObrigatorio.texto', { empresa: sessao.empresa.nome })}</Alerta>
      <SecaoMFA mfaAtivo={sessao.usuario.mfa_ativo} exigidoPelaEmpresa />
      {(sessao.empresas?.length ?? 0) > 1 && (
        <details className="border-t border-borda pt-4">
          <summary className="cursor-pointer texto-rotulo">{t('mfaObrigatorio.outraEmpresa')}</summary>
          <div className="mt-3">
            <EscolhaDeEmpresa
              sessao={sessao}
              aoConcluir={() => void navegar('/', { replace: true })}
              rotuloConfirmar={t('mfaObrigatorio.trocar')}
            />
          </div>
        </details>
      )}
      <Botao
        variante="terciario"
        onClick={() => {
          void navegar('/entrar', { replace: true });
          sair.mutate();
        }}
      >
        {t('usuarioMenu.sair')}
      </Botao>
    </CartaoPublico>
  );
}
