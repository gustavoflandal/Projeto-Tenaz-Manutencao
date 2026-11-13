import { Redirect, Stack, useSegments } from 'expo-router';
import { ActivityIndicator, StyleSheet, Text, View } from 'react-native';
import { SafeAreaProvider } from 'react-native-safe-area-context';

import { ProvedorDeAcesso, useAcesso } from '../src/acesso/contexto';
import { Botao } from '../src/componentes/Botao';
import { cor } from '../src/componentes/tema';
import { ProvedorDeIdioma, useTraducao } from '../src/i18n/contexto';

function Guarda() {
  const { estado, bloqueado, desbloquear } = useAcesso();
  const { t } = useTraducao();
  const segmentos = useSegments();
  const naTelaDeLogin = segmentos[0] === 'login';

  if (estado.tipo === 'iniciando') {
    return (
      <View style={estilos.centro} accessibilityLabel={t('mobile.carregando')}>
        <ActivityIndicator color={cor.texto} />
      </View>
    );
  }
  if (estado.tipo === 'saiu' && !naTelaDeLogin) return <Redirect href="/login" />;
  if (estado.tipo === 'entrou' && naTelaDeLogin) return <Redirect href="/" />;

  if (estado.tipo === 'entrou' && bloqueado) {
    return (
      <View style={estilos.centro}>
        <Text style={estilos.titulo} accessibilityRole="header">
          {t('mobile.bloqueio.titulo')}
        </Text>
        <Text style={estilos.texto}>{t('mobile.bloqueio.texto')}</Text>
        <Botao
          titulo={t('mobile.bloqueio.desbloquear')}
          aoPressionar={() => void desbloquear(t('mobile.bloqueio.motivoDesbloquear'))}
        />
      </View>
    );
  }

  return (
    <Stack
      screenOptions={{
        headerStyle: { backgroundColor: cor.navFundo },
        headerTintColor: cor.navTexto,
        contentStyle: { backgroundColor: cor.fundo },
      }}
    >
      <Stack.Screen name="index" options={{ title: t('mobile.hoje.titulo') }} />
      <Stack.Screen name="login" options={{ headerShown: false }} />
      <Stack.Screen name="empresas" options={{ title: t('mobile.empresas.titulo') }} />
      <Stack.Screen name="sobre" options={{ title: t('sobre.titulo') }} />
    </Stack>
  );
}

export default function Raiz() {
  return (
    <SafeAreaProvider>
      <ProvedorDeIdioma>
        <ProvedorDeAcesso>
          <Guarda />
        </ProvedorDeAcesso>
      </ProvedorDeIdioma>
    </SafeAreaProvider>
  );
}

const estilos = StyleSheet.create({
  centro: { flex: 1, alignItems: 'center', justifyContent: 'center', padding: 24, gap: 16, backgroundColor: cor.fundo },
  titulo: { fontSize: 22, fontWeight: '700', color: cor.texto },
  texto: { fontSize: 16, color: cor.textoSecundario, textAlign: 'center' },
});
