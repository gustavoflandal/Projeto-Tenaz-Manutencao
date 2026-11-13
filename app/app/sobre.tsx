// Sobre a Tenaz: versão do release em execução, binário da loja com o número
// de build, identificador da atualização sem loja, versão da API
// (GET /v1/versao, público), data do build, commit e ambiente. "Copiar
// informações" leva o mesmo texto para o suporte.

import * as Clipboard from 'expo-clipboard';
import { useEffect, useState } from 'react';
import { ScrollView, StyleSheet, Text, View } from 'react-native';

import { informacoesDoApp, urlDaApi } from '../src/acesso/nativo';
import { linhasDoSobre, textoParaCopiar, type VersaoDaApi } from '../src/acesso/versao';
import { Botao } from '../src/componentes/Botao';
import { cor } from '../src/componentes/tema';
import { useTraducao } from '../src/i18n/contexto';

const info = informacoesDoApp();

export default function Sobre() {
  const { t } = useTraducao();
  const [api, setApi] = useState<VersaoDaApi>({ estado: 'consultando' });
  const [copia, setCopia] = useState<'copiado' | 'falhou' | null>(null);

  useEffect(() => {
    let vivo = true;
    fetch(`${urlDaApi()}/v1/versao`)
      .then((r) => (r.ok ? (r.json() as Promise<{ versao: string }>) : Promise.reject(new Error(String(r.status)))))
      .then((v) => vivo && setApi({ estado: 'ok', versao: v.versao }))
      // Sem rede é o estado normal do campo: não é erro.
      .catch(() => vivo && setApi({ estado: 'indisponivel' }));
    return () => {
      vivo = false;
    };
  }, []);

  const linhas = linhasDoSobre(t, info, api);

  async function copiar() {
    try {
      await Clipboard.setStringAsync(textoParaCopiar(t, linhas));
      setCopia('copiado');
    } catch {
      setCopia('falhou');
    }
  }

  return (
    <ScrollView style={{ backgroundColor: cor.fundo }} contentContainerStyle={estilos.tela}>
      <Text style={estilos.titulo} accessibilityRole="header">
        {t('sobre.titulo')}
      </Text>
      {linhas.map((l) => (
        <View key={l.rotulo} style={estilos.linha} accessible accessibilityLabel={t('sobre.linhaCopia', { rotulo: l.rotulo, valor: l.valor })}>
          <Text style={estilos.rotulo}>{l.rotulo}</Text>
          <Text style={estilos.valor}>{l.valor}</Text>
        </View>
      ))}
      <Botao titulo={t('sobre.copiar')} variante="secundario" aoPressionar={() => void copiar()} />
      {copia ? (
        <Text style={estilos.retorno} accessibilityLiveRegion="polite">
          {copia === 'copiado' ? t('sobre.copiado') : t('sobre.copiaFalhou')}
        </Text>
      ) : null}
    </ScrollView>
  );
}

const estilos = StyleSheet.create({
  tela: { padding: 16, gap: 12 },
  titulo: { fontSize: 22, fontWeight: '700', color: cor.texto },
  linha: { backgroundColor: cor.superficie, borderRadius: 8, padding: 12, borderWidth: 1, borderColor: cor.borda, minHeight: 48 },
  rotulo: { fontSize: 14, color: cor.textoTerciario, textTransform: 'uppercase', letterSpacing: 0.6 },
  valor: { fontSize: 16, color: cor.texto, fontFamily: 'monospace' },
  retorno: { fontSize: 14, color: cor.textoSecundario },
});
