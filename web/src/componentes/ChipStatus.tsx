// Chip de status (componente 6.5 do documento 7): ícone de 14 px e a palavra
// do estado, nunca só cor. Não é botão.

import { CircleCheck, CircleSlash, Clock, Mail, ShieldAlert, type LucideIcon } from 'lucide-react';
import { useTranslation } from 'react-i18next';

type Cor = 'operando' | 'preventiva' | 'inativo' | 'parado';

const ESTILO: Record<Cor, string> = {
  operando: 'bg-status-operando-fundo text-status-operando-texto',
  preventiva: 'bg-status-preventiva-fundo text-status-preventiva-texto',
  inativo: 'bg-status-inativo-fundo text-status-inativo-texto',
  parado: 'bg-status-parado-fundo text-status-parado-texto',
};

const ICONE_COR: Record<Cor, string> = {
  operando: 'text-status-operando',
  preventiva: 'text-status-preventiva',
  inativo: 'text-status-inativo',
  parado: 'text-status-parado',
};

interface Definicao {
  readonly cor: Cor;
  readonly icone: LucideIcon;
}

/** Situações do cadastro (empresa, pessoa, unidade e usuário). */
const SITUACOES: Record<'ativa' | 'ativo' | 'inativa' | 'inativo' | 'convidado', Definicao> = {
  ativa: { cor: 'operando', icone: CircleCheck },
  ativo: { cor: 'operando', icone: CircleCheck },
  convidado: { cor: 'preventiva', icone: Mail },
  inativa: { cor: 'inativo', icone: CircleSlash },
  inativo: { cor: 'inativo', icone: CircleSlash },
};

export type SituacaoDoChip = keyof typeof SITUACOES;

export function ChipSituacao({ situacao }: { readonly situacao: SituacaoDoChip }) {
  const { t } = useTranslation();
  const { cor, icone: Icone } = SITUACOES[situacao];
  return (
    <span className={`inline-flex items-center gap-1.5 rounded-sm py-0.75 pr-2.5 pl-2 texto-legenda font-semibold ${ESTILO[cor]}`}>
      <Icone size={14} aria-hidden="true" className={ICONE_COR[cor]} />
      {t(`situacao.${situacao}`)}
    </span>
  );
}

/** Estado de um dia na verificação de integridade da trilha: sempre texto e ícone, nunca só cor. */
const VERIFICACAO: Record<'integro' | 'violado' | 'nao_selado', Definicao> = {
  integro: { cor: 'operando', icone: CircleCheck },
  violado: { cor: 'parado', icone: ShieldAlert },
  nao_selado: { cor: 'inativo', icone: Clock },
};

export function ChipVerificacao({ estado }: { readonly estado: keyof typeof VERIFICACAO }) {
  const { t } = useTranslation();
  const { cor, icone: Icone } = VERIFICACAO[estado];
  return (
    <span className={`inline-flex items-center gap-1.5 rounded-sm py-0.75 pr-2.5 pl-2 texto-legenda font-semibold ${ESTILO[cor]}`}>
      <Icone size={14} aria-hidden="true" className={ICONE_COR[cor]} />
      {t(`integridade.estados.${estado}`)}
    </span>
  );
}
