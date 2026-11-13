// Escolha da empresa (documento 3, seção 4.2): a pessoa com acesso a mais de
// uma empresa escolhe em qual trabalhar. A atual vem primeiro e já marcada;
// escolher outra chama a troca de empresa, que devolve a sessão da escolhida.

import { useState } from 'react';
import { useTranslation } from 'react-i18next';

import { Alerta } from '../../../componentes/Alerta';
import { Botao } from '../../../componentes/Botao';
import { mensagemDoErro } from '../../../lib/erros';
import type { Sessao } from '../../../lib/tipos';
import { useTrocarEmpresa } from '../sessao';

export interface EscolhaDeEmpresaProps {
  readonly sessao: Sessao;
  /** Depois de escolher (e de trocar, se for outra). */
  readonly aoConcluir: () => void;
  /** Usado em /autorizar-app: devolve a escolha sem trocar a sessão da web. */
  readonly aoEscolher?: (empresaId: string) => void;
  readonly rotuloConfirmar?: string;
  readonly carregando?: boolean;
}

export function EscolhaDeEmpresa({ sessao, aoConcluir, aoEscolher, rotuloConfirmar, carregando }: EscolhaDeEmpresaProps) {
  const { t } = useTranslation();
  const empresas = sessao.empresas ?? [];
  const atual = empresas.find((e) => e.atual) ?? empresas[0];
  const [escolhida, setEscolhida] = useState(atual?.id ?? '');
  const trocar = useTrocarEmpresa();

  const confirmar = () => {
    if (aoEscolher) {
      aoEscolher(escolhida);
      return;
    }
    if (escolhida === atual?.id) {
      aoConcluir();
      return;
    }
    trocar.mutate(escolhida, { onSuccess: aoConcluir });
  };

  return (
    <form
      noValidate
      onSubmit={(evento) => {
        evento.preventDefault();
        confirmar();
      }}
      className="flex flex-col gap-4"
    >
      {trocar.isError && <Alerta tipo="erro">{mensagemDoErro(trocar.error)}</Alerta>}
      <fieldset className="flex flex-col gap-2">
        <legend className="mb-2 texto-corpo text-texto-secundario">{t('entrar.escolhaEmpresaTexto')}</legend>
        {empresas.map((empresa) => (
          <label
            key={empresa.id}
            className="flex min-h-10 cursor-pointer items-center gap-3 rounded-md border-2 border-borda-controle px-3 py-2 texto-corpo hover:border-borda-controle-hover has-checked:border-texto"
          >
            <input
              type="radio"
              name="empresa"
              value={empresa.id}
              checked={escolhida === empresa.id}
              onChange={() => setEscolhida(empresa.id)}
              className="size-5 accent-texto"
            />
            <span className="flex-1">{empresa.nome}</span>
            {empresa.atual && <span className="texto-legenda text-texto-secundario">{t('entrar.empresaAtual')}</span>}
          </label>
        ))}
      </fieldset>
      <Botao type="submit" variante="primario" carregando={trocar.isPending || carregando} className="w-full">
        {rotuloConfirmar ?? t('entrar.continuar')}
      </Botao>
    </form>
  );
}
