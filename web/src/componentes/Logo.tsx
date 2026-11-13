// Logo da Tenaz (seção 2.3 do documento 7). Os arquivos vêm de marca/logo/svg,
// sem cópia: a marca é a fonte.

import { useTranslation } from 'react-i18next';

import horizontalNegativo from '../../../marca/logo/svg/tenaz-horizontal-negativo.svg';
import horizontalPositivo from '../../../marca/logo/svg/tenaz-horizontal-positivo.svg';
import simboloNegativo from '../../../marca/logo/svg/tenaz-simbolo-negativo.svg';
import simboloPositivo from '../../../marca/logo/svg/tenaz-simbolo-positivo.svg';

const ARQUIVOS = {
  horizontal: { positivo: horizontalPositivo, negativo: horizontalNegativo },
  simbolo: { positivo: simboloPositivo, negativo: simboloNegativo },
} as const;

export interface LogoProps {
  readonly forma: 'horizontal' | 'simbolo';
  /** Positivo sobre fundo claro; negativo sobre Aço e no tema escuro. */
  readonly variante: 'positivo' | 'negativo';
  readonly className?: string;
}

export function Logo({ forma, variante, className }: LogoProps) {
  const { t } = useTranslation();
  return <img src={ARQUIVOS[forma][variante]} alt={t('app.nome')} className={className} />;
}
