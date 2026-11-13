// Seletor da subdivisão do endereço (documento 3, seção 2.18): a lista é a do
// país da empresa, com o nome oficial de cada uma; o código ISO 3166-2 é o valor.

import { SUBDIVISOES, type Pais } from '@tenaz/dominio';
import { useTranslation } from 'react-i18next';

import { CampoSelecao, type CampoSelecaoProps } from '../../../componentes/Campo';
import { traduzirOuOriginal } from '../../../lib/i18n';

type Props = Omit<CampoSelecaoProps, 'rotulo' | 'children' | 'opcional'> & {
  /** País da empresa; sem ele (ainda carregando), a lista do Brasil. */
  readonly pais: Pais | null;
  /** Valor já gravado: continua escolhível mesmo se não for do país atual. */
  readonly atual?: string | undefined;
};

export function CampoDeSubdivisao({ pais, atual, ...resto }: Props) {
  const { t } = useTranslation();
  const lista = SUBDIVISOES[pais ?? 'BR'];
  const fora = atual && !lista.includes(atual) ? atual : null;
  return (
    <CampoSelecao rotulo={t(`subdivisao.rotulo.${pais ?? 'BR'}`)} opcional {...resto}>
      <option value="">{t('subdivisao.escolha')}</option>
      {fora && <option value={fora}>{fora}</option>}
      {lista.map((codigo) => (
        <option key={codigo} value={codigo}>
          {traduzirOuOriginal('subdivisao.nomes', codigo)}
        </option>
      ))}
    </CampoSelecao>
  );
}

/** Nome oficial da subdivisão para listas e detalhes; sem catálogo, o próprio código. */
export function nomeDaSubdivisao(codigo: string | undefined | null): string {
  return codigo ? traduzirOuOriginal('subdivisao.nomes', codigo) : '';
}
