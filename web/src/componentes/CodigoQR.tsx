// QR code da ativação do segundo fator. Desenhado como SVG a partir da matriz
// de módulos; a cor vem de `currentColor`, sem valor de cor escrito à mão. O
// fundo é sempre claro, para o leitor do aplicativo autenticador funcionar
// também no tema escuro.

import { encode } from 'uqr';
import { useMemo } from 'react';

const MARGEM = 4; // zona silenciosa de 4 módulos, como manda a norma do QR

export function CodigoQR({ conteudo, rotulo }: { readonly conteudo: string; readonly rotulo: string }) {
  const caminho = useMemo(() => {
    const matriz = encode(conteudo, { ecc: 'M', border: 0 }).data;
    const partes: string[] = [];
    matriz.forEach((linha, y) => {
      linha.forEach((escuro, x) => {
        if (escuro) partes.push(`M${x + MARGEM} ${y + MARGEM}h1v1h-1z`);
      });
    });
    return { d: partes.join(''), tamanho: matriz.length + MARGEM * 2 };
  }, [conteudo]);

  return (
    <svg
      role="img"
      aria-label={rotulo}
      viewBox={`0 0 ${caminho.tamanho} ${caminho.tamanho}`}
      className="size-48 rounded-md bg-nav-texto text-nav-fundo"
      shapeRendering="crispEdges"
    >
      <path d={caminho.d} fill="currentColor" />
    </svg>
  );
}
