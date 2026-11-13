// Tela principal (casca de "Hoje"): empresa e usuário da sessão, estado vazio da
// lista, troca de empresa, idioma, bloqueio local, Sobre e Sair. As OS entram
// com o RF-MOB-002 em diante.

import { useRouter } from 'expo-router';
import { useState } from 'react';
import { ScrollView, StyleSheet, Switch, Text, View } from 'react-native';

import { useAcesso } from '../src/acesso/contexto';
import { chaveDoErro } from '../src/acesso/erros';
import { gerente } from '../src/acesso/nativo';
import { Botao } from '../src/componentes/Botao';
import { Mensagem } from '../src/componentes/Mensagem';
import { SeletorDeIdioma } from '../src/componentes/SeletorDeIdioma';
import { SeloVersao } from '../src/componentes/SeloVersao';
import { ALVO_DE_TOQUE, cor } from '../src/componentes/tema';
import type { ChaveDeTexto } from '../src/i18n';
import { useTraducao } from '../src/i18n/contexto';

export default function Principal() {
  const { estado, bloqueioAtivo, alternarBloqueio } = useAcesso();
  const { t } = useTraducao();
  const roteador = useRouter();
  const [aviso, setAviso] = useState<ChaveDeTexto | null>(null);
  const [ocupado, setOcupado] = useState(false);

  if (estado.tipo !== 'entrou') return null;
  const conta = estado.conta;

  async function executar(acao: () => Promise<unknown>) {
    setOcupado(true);
    setAviso(null);
    try {
      await acao();
    } catch (e) {
      setAviso(chaveDoErro(e));
    } finally {
      setOcupado(false);
    }
  }

  return (
    <View style={estilos.tela}>
      <ScrollView contentContainerStyle={estilos.conteudo}>
        <Text style={estilos.titulo} accessibilityRole="header">
          {conta.empresaNome}
        </Text>
        <Text style={estilos.texto}>{conta.usuarioNome}</Text>

        {conta.mfaPendente ? (
          <Mensagem variante="atencao" titulo={t('mobile.mfa.titulo')} texto={t('mobile.mfa.texto')}>
            <Botao
              titulo={t('mobile.mfa.verificar')}
              variante="secundario"
              ocupado={ocupado}
              aoPressionar={() => void executar(() => gerente.atualizarSessao())}
            />
          </Mensagem>
        ) : (
          <View style={estilos.cartao}>
            <Text style={estilos.cartaoTitulo}>{t('mobile.hoje.vazioTitulo')}</Text>
            <Text style={estilos.texto}>{t('mobile.hoje.vazioTexto')}</Text>
          </View>
        )}

        {aviso ? <Mensagem variante="erro" texto={t(aviso)} /> : null}

        {conta.empresas.length > 1 ? (
          <Botao titulo={t('mobile.acoes.trocarEmpresa')} variante="secundario" aoPressionar={() => roteador.push('/empresas')} />
        ) : null}

        <SeletorDeIdioma />

        <View style={estilos.linha}>
          <Text style={[estilos.texto, { flex: 1 }]}>{t('mobile.bloqueio.rotulo')}</Text>
          <Switch
            accessibilityLabel={t('mobile.bloqueio.rotuloAcessivel')}
            value={bloqueioAtivo}
            trackColor={{ false: cor.bordaControle, true: cor.texto }}
            thumbColor={cor.superficie}
            onValueChange={() =>
              void executar(async () => {
                if (!(await alternarBloqueio(t('mobile.bloqueio.motivoAtivar')))) setAviso('mobile.bloqueio.indisponivel');
              })
            }
            style={{ minHeight: ALVO_DE_TOQUE }}
          />
        </View>

        <Botao titulo={t('mobile.acoes.sobre')} variante="secundario" aoPressionar={() => roteador.push('/sobre')} />
        <Botao
          titulo={t('mobile.acoes.sair')}
          variante="perigo"
          ocupado={ocupado}
          rotuloAcessivel={t('mobile.acoes.sairRotulo')}
          aoPressionar={() => void executar(() => gerente.sair())}
        />
        <Text style={estilos.apoio}>{t('mobile.acoes.sairAviso')}</Text>
      </ScrollView>
      <SeloVersao />
    </View>
  );
}

const estilos = StyleSheet.create({
  tela: { flex: 1, backgroundColor: cor.fundo },
  // Folga inferior para o conteúdo nunca ficar sob o selo (48 dp + 8 dp + 16 dp de respiro).
  conteudo: { padding: 16, gap: 16, paddingBottom: 72 },
  titulo: { fontSize: 24, fontWeight: '700', color: cor.texto },
  texto: { fontSize: 16, color: cor.textoSecundario },
  apoio: { fontSize: 14, color: cor.textoTerciario },
  cartao: { backgroundColor: cor.superficie, borderRadius: 8, padding: 16, gap: 4, borderWidth: 1, borderColor: cor.borda },
  cartaoTitulo: { fontSize: 17, fontWeight: '700', color: cor.texto },
  linha: { flexDirection: 'row', alignItems: 'center', gap: 12, minHeight: ALVO_DE_TOQUE },
});
