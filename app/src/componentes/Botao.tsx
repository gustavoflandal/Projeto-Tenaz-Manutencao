import { ActivityIndicator, Pressable, StyleSheet, Text } from 'react-native';

import { ALVO_DE_TOQUE, ALVO_PRIMARIO, cor } from './tema';

type Variante = 'primario' | 'secundario' | 'perigo';

interface Props {
  titulo: string;
  aoPressionar: () => void;
  /** primario: Sinal com texto Aço (56 dp). secundario: contorno. perigo: ação que sai ou descarta. */
  variante?: Variante;
  ocupado?: boolean;
  desabilitado?: boolean;
  /** Nome para o leitor de tela, quando o título sozinho não basta. */
  rotuloAcessivel?: string;
}

const ESTILOS: Record<Variante, { fundo: string; pressionado: string; texto: string; borda?: string }> = {
  primario: { fundo: cor.acaoPrimaria, pressionado: cor.acaoPrimariaPressionada, texto: cor.acaoPrimariaTexto },
  secundario: { fundo: cor.superficie, pressionado: cor.superficieSelecionada, texto: cor.texto, borda: cor.bordaControle },
  perigo: { fundo: cor.perigo, pressionado: cor.perigoPressionado, texto: cor.perigoTexto },
};

export function Botao({ titulo, aoPressionar, variante = 'primario', ocupado, desabilitado, rotuloAcessivel }: Props) {
  const inativo = !!(desabilitado || ocupado);
  const e = ESTILOS[variante];
  return (
    <Pressable
      accessibilityRole="button"
      accessibilityLabel={rotuloAcessivel ?? titulo}
      accessibilityState={{ disabled: inativo, busy: !!ocupado }}
      disabled={inativo}
      onPress={aoPressionar}
      style={({ pressed }) => [
        estilos.base,
        { minHeight: variante === 'primario' ? ALVO_PRIMARIO : ALVO_DE_TOQUE },
        // Desabilitado usa fundo e texto de tokens, não opacidade: opacidade derruba o
        // contraste do primário a 1,9:1; assim o texto segue legível e o botão claramente inativo.
        inativo
          ? { backgroundColor: cor.superficieSelecionada, borderWidth: 0 }
          : {
              backgroundColor: pressed ? e.pressionado : e.fundo,
              borderWidth: e.borda ? 1.5 : 0,
              borderColor: e.borda,
            },
      ]}
    >
      {ocupado ? (
        <ActivityIndicator color={cor.texto} />
      ) : (
        <Text style={[estilos.texto, { color: inativo ? cor.textoSecundario : e.texto }]}>{titulo}</Text>
      )}
    </Pressable>
  );
}

const estilos = StyleSheet.create({
  base: {
    paddingHorizontal: 20,
    borderRadius: 8,
    alignItems: 'center',
    justifyContent: 'center',
  },
  texto: { fontSize: 16, fontWeight: '600' },
});
