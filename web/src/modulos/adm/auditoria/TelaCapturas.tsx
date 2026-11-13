// Capturas da auditoria por tela (RF-ADM-011): o que a trilha grava além do
// que é sempre gravado. Por tela, três chaves (visualização, impressão e
// extração); cada troca vale na hora, só para eventos futuros, e é uma
// alteração sempre gravada na trilha. Exige `adm.gerenciar_parametros`.

import { useState } from 'react';
import { useTranslation } from 'react-i18next';

import { Alerta } from '../../../componentes/Alerta';
import { CabecalhoPagina } from '../../../componentes/CabecalhoPagina';
import { InterruptorDeTabela } from '../../../componentes/Campo';
import { Carregando, EstadoErro, EstadoVazio } from '../../../componentes/Estados';
import { Tabela, Td, Th, Tr } from '../../../componentes/Tabela';
import { mensagemDoErro } from '../../../lib/erros';
import type { CapturaAuditoria } from '../../../lib/tipos';
import { useCapturas, useDefinirCaptura } from '../auditoria';
import { CAPTURAS_CONFIGURAVEIS, type CapturaConfiguravel } from '../constantes';

function agruparPorModulo(telas: readonly CapturaAuditoria[]): Array<readonly [string, CapturaAuditoria[]]> {
  const modulos = [...new Set(telas.map((t) => t.modulo))];
  return modulos.map((modulo) => [modulo, telas.filter((t) => t.modulo === modulo)] as const);
}

export function TelaCapturas() {
  const { t } = useTranslation();
  const consulta = useCapturas();
  const definir = useDefinirCaptura();
  const [aviso, setAviso] = useState<string | null>(null);

  const alternar = (tela: CapturaAuditoria, captura: CapturaConfiguravel, ligada: boolean) => {
    setAviso(null);
    definir.mutate(
      {
        tela: tela.tela,
        chaves: {
          visualizacao: captura === 'visualizacao' ? ligada : tela.visualizacao,
          impressao: captura === 'impressao' ? ligada : tela.impressao,
          extracao: captura === 'extracao' ? ligada : tela.extracao,
        },
      },
      {
        onSuccess: () =>
          setAviso(
            t(ligada ? 'capturas.ligada' : 'capturas.desligada', {
              captura: t(`capturas.${captura}`),
              tela: tela.nome,
            }),
          ),
      },
    );
  };

  return (
    <>
      <CabecalhoPagina
        titulo={t('capturas.titulo')}
        trilha={[{ rotulo: t('auditoria.titulo'), para: '/administracao/auditoria' }, { rotulo: t('capturas.titulo') }]}
      />
      <p className="mb-4 max-w-180 texto-corpo text-texto-secundario">{t('capturas.texto')}</p>
      <div className="mb-6 max-w-180 rounded-md border border-borda bg-superficie p-3">
        <p className="texto-rotulo">{t('capturas.sempreTitulo')}</p>
        <p className="texto-corpo-p text-texto-secundario">{t('capturas.sempre')}</p>
        <p className="mt-1 texto-legenda text-texto-secundario">{t('capturas.padraoTexto')}</p>
      </div>

      {definir.isError && (
        <Alerta tipo="erro" className="mb-4">
          {mensagemDoErro(definir.error)}
        </Alerta>
      )}
      <div role="status" className="mb-4 min-h-6 texto-corpo-p">
        {aviso}
      </div>

      {consulta.isPending && <Carregando rotulo={t('capturas.carregando')} />}
      {consulta.isError && <EstadoErro erro={consulta.error} aoTentarDeNovo={() => void consulta.refetch()} />}
      {consulta.isSuccess && (consulta.data.dados ?? []).length === 0 && (
        <EstadoVazio titulo={t('capturas.vazioTitulo')} texto={t('capturas.vazio')} />
      )}
      {consulta.isSuccess && (
        <div className="flex flex-col gap-8">
          {agruparPorModulo(consulta.data.dados ?? []).map(([modulo, telas]) => (
            <section key={modulo} aria-labelledby={`capturas-${modulo}`}>
              <h2 id={`capturas-${modulo}`} className="mb-2 texto-h3">
                {t('capturas.moduloTitulo', { modulo })}
              </h2>
              <Tabela legenda={t('capturas.moduloTitulo', { modulo })}>
                <thead>
                  <tr>
                    <Th>{t('capturas.tela')}</Th>
                    {CAPTURAS_CONFIGURAVEIS.map((captura) => (
                      <Th key={captura}>{t(`capturas.${captura}`)}</Th>
                    ))}
                  </tr>
                </thead>
                <tbody>
                  {telas.map((tela) => (
                    <Tr key={tela.tela}>
                      <Td>
                        {tela.nome}
                        <span className="ml-2 texto-dado text-texto-secundario">{tela.tela}</span>
                        {!tela.configurada && <span className="ml-2 texto-legenda text-texto-secundario">{t('capturas.padrao')}</span>}
                      </Td>
                      {CAPTURAS_CONFIGURAVEIS.map((captura) => (
                        <Td key={captura}>
                          <InterruptorDeTabela
                            rotulo={t('capturas.ajuste', { captura: t(`capturas.${captura}`), tela: tela.nome })}
                            marcado={tela[captura]}
                            desabilitado={definir.isPending && definir.variables?.tela === tela.tela}
                            aoMudar={(ligada) => alternar(tela, captura, ligada)}
                          />
                        </Td>
                      ))}
                    </Tr>
                  ))}
                </tbody>
              </Tabela>
            </section>
          ))}
        </div>
      )}
    </>
  );
}
