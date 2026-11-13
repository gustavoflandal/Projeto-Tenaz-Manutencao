// Troca de empresa sem sair (POST /v1/auth/trocar-empresa). Precisa de
// conexão: o servidor revoga a credencial da empresa atual e emite outra. O
// banco local de cada empresa fica separado e intacto.

import { useRouter } from 'expo-router';
import { useState } from 'react';
import { FlatList, Pressable, StyleSheet, Text, View } from 'react-native';

import { useAcesso } from '../src/acesso/contexto';
import { chaveDoErro, ErroDeRede } from '../src/acesso/erros';
import { gerente } from '../src/acesso/nativo';
import { Mensagem } from '../src/componentes/Mensagem';
import { cor } from '../src/componentes/tema';
import type { ChaveDeTexto } from '../src/i18n';
import { useTraducao } from '../src/i18n/contexto';

export default function Empresas() {
  const { estado } = useAcesso();
  const { t } = useTraducao();
  const roteador = useRouter();
  const [ocupado, setOcupado] = useState<string | null>(null);
  const [erro, setErro] = useState<ChaveDeTexto | null>(null);

  if (estado.tipo !== 'entrou') return null;
  const conta = estado.conta;

  async function trocar(id: string) {
    setOcupado(id);
    setErro(null);
    try {
      await gerente.trocarEmpresa(id);
      roteador.replace('/');
    } catch (e) {
      setErro(e instanceof ErroDeRede ? 'mobile.empresas.precisaConexao' : chaveDoErro(e));
    } finally {
      setOcupado(null);
    }
  }

  return (
    <View style={estilos.tela}>
      {erro ? <Mensagem variante="erro" texto={t(erro)} /> : null}
      <FlatList
        data={conta.empresas}
        keyExtractor={(e) => e.id}
        contentContainerStyle={{ gap: 8 }}
        ListEmptyComponent={<Text style={estilos.detalhe}>{t('mobile.empresas.vazio')}</Text>}
        renderItem={({ item }) => {
          const atual = item.id === conta.empresaId;
          return (
            <Pressable
              accessibilityRole="button"
              accessibilityLabel={
                atual
                  ? t('mobile.empresas.atualRotulo', { nome: item.nome })
                  : t('mobile.empresas.entrarRotulo', { nome: item.nome })
              }
              accessibilityState={{ selected: atual, disabled: atual || ocupado !== null, busy: ocupado === item.id }}
              disabled={atual || ocupado !== null}
              onPress={() => void trocar(item.id)}
              style={({ pressed }) => [
                estilos.item,
                atual && estilos.itemAtual,
                pressed && { backgroundColor: cor.superficieSelecionada },
              ]}
            >
              <Text style={estilos.nome}>{item.nome}</Text>
              <Text style={estilos.detalhe}>
                {atual ? t('mobile.empresas.atual') : ocupado === item.id ? t('mobile.empresas.trocando') : t('mobile.empresas.entrar')}
              </Text>
            </Pressable>
          );
        }}
      />
    </View>
  );
}

const estilos = StyleSheet.create({
  tela: { flex: 1, padding: 16, gap: 12, backgroundColor: cor.fundo },
  // Item de lista: 64 dp no mínimo, o item inteiro é o alvo (docs/7, 9.2).
  item: {
    minHeight: 64,
    padding: 14,
    borderRadius: 8,
    backgroundColor: cor.superficie,
    borderWidth: 1,
    borderColor: cor.borda,
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    gap: 12,
  },
  itemAtual: { borderColor: cor.texto, borderWidth: 2 },
  nome: { fontSize: 16, fontWeight: '600', color: cor.texto, flex: 1 },
  detalhe: { fontSize: 14, color: cor.textoSecundario },
});
