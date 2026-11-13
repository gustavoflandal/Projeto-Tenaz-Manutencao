// Login (RF-MOB-001): abre o navegador do sistema na tela de autorização da
// web. O app nunca vê a senha. A escolha de empresa, quando há mais de uma, é
// feita nessa tela da web (seção 4.3 de docs/3_ESPECIFICACAO_APIS.md).

import * as WebBrowser from 'expo-web-browser';
import { useEffect, useState } from 'react';
import { ScrollView, StyleSheet, Text, View } from 'react-native';

import { chaveDoErro } from '../src/acesso/erros';
import { gerente } from '../src/acesso/nativo';
import { URI_DE_RETORNO } from '../src/acesso/retorno';
import { Botao } from '../src/componentes/Botao';
import { Mensagem } from '../src/componentes/Mensagem';
import { SeletorDeIdioma } from '../src/componentes/SeletorDeIdioma';
import { SeloVersao } from '../src/componentes/SeloVersao';
import { cor } from '../src/componentes/tema';
import type { ChaveDeTexto } from '../src/i18n';
import { useTraducao } from '../src/i18n/contexto';

export default function Login() {
  const { t } = useTraducao();
  const [ocupado, setOcupado] = useState(false);
  const [erro, setErro] = useState<ChaveDeTexto | null>(null);
  const [lembrete, setLembrete] = useState<string | null>(null);

  useEffect(() => {
    void gerente.ultimaConta().then((c) => {
      if (c) setLembrete(`${c.usuarioNome}, ${c.empresaNome}`);
    });
    const e = gerente.estado();
    if (e.tipo === 'saiu' && e.aviso === 'sessao_encerrada') setErro('mobile.login.sessaoTerminou');
  }, []);

  async function entrar() {
    setOcupado(true);
    setErro(null);
    try {
      const tentativa = await gerente.prepararLogin();
      const r = await WebBrowser.openAuthSessionAsync(tentativa.url, URI_DE_RETORNO);
      if (r.type === 'success') await gerente.concluirLogin(tentativa, r.url);
      // cancel e dismiss: a pessoa fechou o navegador, nada a dizer.
    } catch (e) {
      setErro(chaveDoErro(e));
    } finally {
      setOcupado(false);
    }
  }

  return (
    <View style={estilos.tela}>
      <ScrollView contentContainerStyle={estilos.conteudo}>
        <Text style={estilos.marca} accessibilityRole="header">
          Tenaz
        </Text>
        <Text style={estilos.slogan}>{t('mobile.login.slogan')}</Text>
        {lembrete ? <Text style={estilos.texto}>{t('mobile.login.ultimoAcesso', { conta: lembrete })}</Text> : null}
        {erro ? <Mensagem variante="erro" texto={t(erro)} /> : null}
        <Botao
          titulo={t('mobile.login.entrar')}
          aoPressionar={() => void entrar()}
          ocupado={ocupado}
          rotuloAcessivel={t('mobile.login.entrarRotulo')}
        />
        <SeletorDeIdioma />
      </ScrollView>
      <SeloVersao />
    </View>
  );
}

const estilos = StyleSheet.create({
  tela: { flex: 1, backgroundColor: cor.fundo },
  // Folga inferior para o conteúdo nunca ficar sob o selo (48 dp + 8 dp + respiro).
  conteudo: { flexGrow: 1, justifyContent: 'center', padding: 24, gap: 16, paddingBottom: 72 },
  marca: { fontSize: 36, fontWeight: '800', color: cor.texto },
  slogan: { fontSize: 16, color: cor.textoSecundario },
  texto: { fontSize: 14, color: cor.textoSecundario },
});
