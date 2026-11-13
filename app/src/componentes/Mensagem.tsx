// Alerta em linha (docs/7_PADROES_DESIGN.md, 9.11 e 10.3): texto na cor `texto`,
// com ícone e barra lateral na cor do status. Texto em `perigo` sobre `fundo`
// mede 4,32:1 e reprova AA, por isso a cor de status fica só no ícone e na barra.

import type { ReactNode } from 'react';
import { StyleSheet, Text, View } from 'react-native';

import { cor } from './tema';

interface Props {
  variante: 'erro' | 'atencao';
  titulo?: string;
  texto: string;
  children?: ReactNode;
}

export function Mensagem({ variante, titulo, texto, children }: Props) {
  const status = variante === 'erro' ? cor.perigo : cor.atencao;
  return (
    <View
      accessibilityRole="alert"
      accessibilityLiveRegion="polite"
      style={[
        estilos.caixa,
        { backgroundColor: variante === 'erro' ? cor.superficie : cor.atencaoFundo, borderLeftColor: status },
      ]}
    >
      {/* Ícone decorativo: o leitor de tela lê o texto, não o "!". */}
      <View
        accessibilityElementsHidden
        importantForAccessibility="no-hide-descendants"
        style={[estilos.icone, { backgroundColor: status }]}
      >
        <Text style={[estilos.glifo, { color: variante === 'erro' ? cor.perigoTexto : cor.atencaoTexto }]}>!</Text>
      </View>
      <View style={estilos.corpo}>
        {titulo ? <Text style={estilos.titulo}>{titulo}</Text> : null}
        <Text style={estilos.texto}>{texto}</Text>
        {children}
      </View>
    </View>
  );
}

const estilos = StyleSheet.create({
  caixa: {
    flexDirection: 'row',
    gap: 12,
    padding: 12,
    borderRadius: 8,
    borderWidth: 1,
    borderColor: cor.borda,
    borderLeftWidth: 6,
  },
  icone: { width: 24, height: 24, borderRadius: 12, alignItems: 'center', justifyContent: 'center' },
  glifo: { fontSize: 15, fontWeight: '800' },
  corpo: { flex: 1, gap: 8 },
  titulo: { fontSize: 17, fontWeight: '700', color: cor.texto },
  texto: { fontSize: 16, color: cor.texto },
});
