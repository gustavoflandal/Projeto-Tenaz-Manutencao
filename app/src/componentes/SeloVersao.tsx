// Selo de versão (RNF-039): canto inferior direito, dentro da área segura,
// 8 dp da borda, alvo de toque de 48 x 48 dp, texto à direita e embaixo.
// Só aparece em telas sem botão fixo inferior (Hoje e login); nas demais, a
// versão fica em Sobre. Não usa Sinal, cor de status nem ícone.

import { useRouter } from 'expo-router';
import { Pressable, StyleSheet, Text } from 'react-native';
import { useSafeAreaInsets } from 'react-native-safe-area-context';

import { informacoesDoApp } from '../acesso/nativo';
import { rotuloAcessivelDoSelo, textoDoSelo } from '../acesso/versao';
import { useTraducao } from '../i18n/contexto';
import { ALVO_DE_TOQUE, cor } from './tema';

const info = informacoesDoApp();

export function SeloVersao() {
  const roteador = useRouter();
  const margem = useSafeAreaInsets();
  const { t } = useTraducao();
  return (
    <Pressable
      accessibilityRole="button"
      accessibilityLabel={rotuloAcessivelDoSelo(t, info)}
      onPress={() => roteador.push('/sobre')}
      style={[estilos.area, { right: 8 + margem.right, bottom: 8 + margem.bottom }]}
    >
      <Text style={estilos.texto}>{textoDoSelo(t, info)}</Text>
    </Pressable>
  );
}

const estilos = StyleSheet.create({
  area: {
    position: 'absolute',
    minWidth: ALVO_DE_TOQUE,
    height: ALVO_DE_TOQUE,
    justifyContent: 'flex-end',
    alignItems: 'flex-end',
  },
  texto: {
    fontSize: 12,
    lineHeight: 16,
    fontFamily: 'monospace',
    color: cor.textoTerciario,
    backgroundColor: cor.fundo,
    paddingHorizontal: 6,
    paddingVertical: 2,
    borderRadius: 4,
  },
});
