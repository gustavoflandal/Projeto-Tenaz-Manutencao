// Perfis (RF-ADM-005), somente leitura: os doze perfis padrão com as permissões
// de cada um. Perfil padrão não se altera. As permissões sensíveis aparecem
// destacadas quando o catálogo está ao alcance de quem olha.

import { ShieldAlert } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { CabecalhoPagina } from '../../../componentes/CabecalhoPagina';
import { Carregando, EstadoErro, EstadoVazio } from '../../../componentes/Estados';
import { useTemPermissao } from '../sessao';
import { usePerfis, usePermissoes } from '../usuarios';

export function TelaPerfis() {
  const { t } = useTranslation();
  const perfis = usePerfis();
  const podeVerCatalogo = useTemPermissao('adm.gerenciar_perfis');
  const catalogo = usePermissoes(podeVerCatalogo);
  const sensiveis = new Set((catalogo.data?.dados ?? []).filter((p) => p.sensivel).map((p) => p.codigo));

  return (
    <>
      <CabecalhoPagina titulo={t('perfil.titulo')} />
      <p className="mb-4 max-w-180 texto-corpo text-texto-secundario">{t('perfil.texto')}</p>
      {perfis.isPending && <Carregando rotulo={t('perfil.carregando')} />}
      {perfis.isError && <EstadoErro erro={perfis.error} aoTentarDeNovo={() => void perfis.refetch()} />}
      {perfis.isSuccess && (perfis.data.dados ?? []).length === 0 && (
        <EstadoVazio titulo={t('perfil.vazioTitulo')} texto={t('perfil.vazio')} />
      )}
      {perfis.isSuccess && (
        <ul className="flex max-w-180 flex-col gap-3">
          {(perfis.data.dados ?? []).map((perfil) => (
            <li key={perfil.id} className="rounded-lg border border-borda bg-superficie">
              <details>
                <summary className="flex min-h-12 cursor-pointer flex-wrap items-center justify-between gap-2 px-4 py-2">
                  <span className="texto-h4">{perfil.nome}</span>
                  <span className="texto-legenda text-texto-secundario">
                    {t('perfil.quantidade', { count: perfil.permissoes?.length ?? 0 })}
                    {perfil.padrao ? ` · ${t('perfil.padrao')}` : ''}
                  </span>
                </summary>
                <div className="border-t border-borda px-4 py-3">
                  <p className="mb-2 texto-legenda text-texto-secundario">{t('perfil.codigo', { codigo: perfil.codigo })}</p>
                  <ul className="flex flex-wrap gap-2">
                    {(perfil.permissoes ?? []).map((codigo) => (
                      <li
                        key={codigo}
                        className="inline-flex items-center gap-1 rounded-sm bg-superficie-selecionada px-2 py-1 texto-dado"
                      >
                        {sensiveis.has(codigo) && (
                          <>
                            <ShieldAlert size={14} aria-hidden="true" />
                            <span className="sr-only">{t('perfil.sensivel')}</span>
                          </>
                        )}
                        {codigo}
                      </li>
                    ))}
                  </ul>
                </div>
              </details>
            </li>
          ))}
        </ul>
      )}
    </>
  );
}
