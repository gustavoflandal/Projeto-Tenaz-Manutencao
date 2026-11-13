// Tabela de dados (componente 6.4 do documento 7), versão simples: cabeçalho
// com escopo de coluna, legenda para leitor de tela, linhas com hover em
// `superficie-hover`. Ordenação e colunas configuráveis entram com os módulos
// que precisam delas; as listas desta etapa são ordenadas pelo servidor.

import { useEffect, useRef, type ComponentProps, type ReactNode } from 'react';
import { useTranslation } from 'react-i18next';

import { Botao } from './Botao';

export interface TabelaProps {
  /** Nome da tabela para leitor de tela. */
  readonly legenda: string;
  readonly children: ReactNode;
  readonly ocupado?: boolean;
}

export function Tabela({ legenda, children, ocupado }: TabelaProps) {
  return (
    <div className="overflow-x-auto rounded-lg border border-borda bg-superficie">
      <table aria-busy={ocupado || undefined} className="w-full border-collapse text-left">
        <caption className="sr-only">{legenda}</caption>
        {children}
      </table>
    </div>
  );
}

export function Th({ className, ...resto }: ComponentProps<'th'>) {
  return (
    <th
      scope="col"
      className={['h-10 border-b border-borda px-3 texto-rotulo text-texto-secundario', className].filter(Boolean).join(' ')}
      {...resto}
    />
  );
}

export function Td({ className, ...resto }: ComponentProps<'td'>) {
  return <td className={['h-10 border-b border-borda px-3 py-2 texto-corpo-p', className].filter(Boolean).join(' ')} {...resto} />;
}

/** Linha da tabela. Só a linha inteira clicável pede `interativa` (hover em `superficie-hover`). */
export function Tr({ className, interativa, ...resto }: ComponentProps<'tr'> & { readonly interativa?: boolean }) {
  return <tr className={[interativa && 'hover:bg-superficie-hover', className].filter(Boolean).join(' ')} {...resto} />;
}

export interface RodapeDeListaProps {
  readonly quantidade: number;
  readonly temMais: boolean;
  readonly carregandoMais: boolean;
  readonly aoCarregarMais: () => void;
}

/** Paginação por cursor: mostra quantos itens estão na tela e pede a próxima página. */
export function RodapeDeLista({ quantidade, temMais, carregandoMais, aoCarregarMais }: RodapeDeListaProps) {
  const { t } = useTranslation();
  const contagem = useRef<HTMLParagraphElement>(null);
  const tinhaMais = useRef(temMais);
  // O botão "Carregar mais" some ao fim da lista: o foco vai para a contagem, e não se perde.
  useEffect(() => {
    if (tinhaMais.current && !temMais) contagem.current?.focus();
    tinhaMais.current = temMais;
  }, [temMais]);
  return (
    <div className="mt-3 flex items-center justify-between gap-3">
      <p ref={contagem} tabIndex={-1} role="status" className="texto-dado text-texto-secundario">
        {temMais ? t('lista.carregados', { count: quantidade }) : t('lista.fimDaLista', { count: quantidade })}
      </p>
      {temMais && (
        <Botao carregando={carregandoMais} onClick={aoCarregarMais}>
          {carregandoMais ? t('lista.carregandoMais') : t('lista.carregarMais')}
        </Botao>
      )}
    </div>
  );
}
