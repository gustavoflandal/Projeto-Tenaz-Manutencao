// Valores antes e depois de um evento da trilha (RF-ADM-011): uma linha por
// campo, com o que mudou marcado por texto e não só por cor. Criação não tem
// "antes"; exclusão não tem "depois".

import { useTranslation } from 'react-i18next';

import { Tabela, Td, Th, Tr } from '../../../componentes/Tabela';

function comoObjeto(valor: unknown): Record<string, unknown> {
  return typeof valor === 'object' && valor !== null && !Array.isArray(valor) ? (valor as Record<string, unknown>) : {};
}

function mostrar(valor: unknown, vazio: string): string {
  if (valor === undefined || valor === null) return vazio;
  return typeof valor === 'string' ? valor : JSON.stringify(valor);
}

export interface DiferencaDoEventoProps {
  readonly antes: unknown;
  readonly depois: unknown;
  readonly camposAlterados: readonly string[];
}

export function DiferencaDoEvento({ antes, depois, camposAlterados }: DiferencaDoEventoProps) {
  const { t } = useTranslation();
  const a = comoObjeto(antes);
  const d = comoObjeto(depois);
  const campos = [...new Set([...Object.keys(a), ...Object.keys(d)])].sort();

  if (campos.length === 0) return <p className="texto-corpo-p text-texto-secundario">{t('auditoria.semValores')}</p>;

  return (
    <Tabela legenda={t('auditoria.valores')}>
      <thead>
        <tr>
          <Th>{t('auditoria.campo')}</Th>
          <Th>{t('auditoria.antes')}</Th>
          <Th>{t('auditoria.depois')}</Th>
        </tr>
      </thead>
      <tbody>
        {campos.map((campo) => {
          const mudou = camposAlterados.includes(campo) || mostrar(a[campo], '') !== mostrar(d[campo], '');
          return (
            <Tr key={campo}>
              <Td className="texto-dado">
                {campo}
                {mudou && camposAlterados.length > 0 && <span className="ml-2 texto-legenda text-texto-secundario">{t('auditoria.alterado')}</span>}
              </Td>
              <Td className="texto-dado break-all">{mostrar(a[campo], t('auditoria.semValor'))}</Td>
              <Td className="texto-dado break-all">{mostrar(d[campo], t('auditoria.semValor'))}</Td>
            </Tr>
          );
        })}
      </tbody>
    </Tabela>
  );
}
