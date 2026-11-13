// Seletor de idioma: o nome de cada idioma no próprio idioma (docs/7, 10.7).
// Vale antes do login, por isso aparece na tela de login e em Hoje.

import { Pressable, StyleSheet, Text, View } from 'react-native';

import { useTraducao } from '../i18n/contexto';
import { ALVO_DE_TOQUE, cor } from './tema';

export function SeletorDeIdioma() {
  const { idioma, idiomas, trocarIdioma, t } = useTraducao();
  return (
    <View accessibilityRole="radiogroup" accessibilityLabel={t('mobile.idioma.rotulo')} style={estilos.grupo}>
      <Text style={estilos.rotulo}>{t('mobile.idioma.rotulo')}</Text>
      <View style={estilos.opcoes}>
        {idiomas.map((i) => {
          const nome = t(`idiomas.${i}`);
          const atual = i === idioma;
          return (
            <Pressable
              key={i}
              accessibilityRole="radio"
              accessibilityLabel={t('mobile.idioma.opcaoRotulo', { nome })}
              accessibilityState={{ selected: atual, checked: atual }}
              onPress={() => void trocarIdioma(i)}
              style={({ pressed }) => [
                estilos.opcao,
                atual && estilos.opcaoAtual,
                pressed && { backgroundColor: cor.superficieSelecionada },
              ]}
            >
              <Text style={[estilos.nome, atual && { fontWeight: '700' }]}>{nome}</Text>
            </Pressable>
          );
        })}
      </View>
    </View>
  );
}

const estilos = StyleSheet.create({
  grupo: { gap: 8 },
  rotulo: { fontSize: 14, color: cor.textoTerciario },
  opcoes: { flexDirection: 'row', flexWrap: 'wrap', gap: 8 },
  opcao: {
    minHeight: ALVO_DE_TOQUE,
    paddingHorizontal: 14,
    borderRadius: 8,
    borderWidth: 1.5,
    borderColor: cor.bordaControle,
    backgroundColor: cor.superficie,
    justifyContent: 'center',
  },
  opcaoAtual: { borderColor: cor.texto, borderWidth: 2 },
  nome: { fontSize: 15, color: cor.texto },
});
