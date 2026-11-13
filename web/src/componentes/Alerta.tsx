// Alerta em linha (componente 6.12 do documento 7): ícone, texto em `texto`,
// barra lateral de 4 px na cor do tipo. `role="alert"` para erro e
// `role="status"` para informação, sucesso e atenção.

import { CircleCheck, CircleX, Info, TriangleAlert, type LucideIcon } from 'lucide-react';
import { useEffect, useRef, type ReactNode } from 'react';

export type TipoDeAlerta = 'informacao' | 'sucesso' | 'atencao' | 'erro';

const ESTILO: Record<TipoDeAlerta, { readonly icone: LucideIcon; readonly classe: string; readonly cor: string }> = {
  informacao: { icone: Info, classe: 'bg-superficie border-texto', cor: 'text-texto' },
  sucesso: { icone: CircleCheck, classe: 'bg-status-operando-fundo border-status-operando', cor: 'text-status-operando' },
  atencao: { icone: TriangleAlert, classe: 'bg-status-atencao-fundo border-status-atencao', cor: 'text-status-atencao' },
  erro: { icone: CircleX, classe: 'bg-status-parado-fundo border-status-parado', cor: 'text-status-parado' },
};

export interface AlertaProps {
  readonly tipo: TipoDeAlerta;
  readonly titulo?: string;
  readonly children?: ReactNode;
  /** Até duas ações terciárias, ao pé do alerta. */
  readonly acoes?: ReactNode;
  readonly className?: string;
  /** Rola até o alerta ao aparecer: para sucesso ou erro no topo de formulário comprido. */
  readonly rolar?: boolean;
}

export function Alerta({ tipo, titulo, children, acoes, className, rolar }: AlertaProps) {
  const { icone: Icone, classe, cor } = ESTILO[tipo];
  const caixa = useRef<HTMLDivElement>(null);
  useEffect(() => {
    if (rolar) caixa.current?.scrollIntoView({ block: 'nearest' });
  }, [rolar]);
  return (
    <div
      ref={caixa}
      role={tipo === 'erro' ? 'alert' : 'status'}
      className={['flex items-start gap-3 rounded-md border-l-4 p-3 text-texto', classe, className].filter(Boolean).join(' ')}
    >
      <Icone size={20} aria-hidden="true" className={`mt-0.5 shrink-0 ${cor}`} />
      <div className="flex min-w-0 flex-col items-start gap-1">
        {titulo && <p className="texto-h4">{titulo}</p>}
        {children && <div className="texto-corpo-p">{children}</div>}
        {acoes && <div className="mt-1 flex flex-wrap gap-2">{acoes}</div>}
      </div>
    </div>
  );
}
