// Configuração do Tenaz Mobile.
//
// A versão NUNCA é digitada aqui: vem do CI, derivado da tag git por
// ferramentas/versao.sh (RNF-039 e seção 9.5 de
// docs/5_ARQUITETURA_SOFTWARE.md).
//
// Três números diferentes convivem, e confundi-los é o erro comum:
//
//   version          versão do RELEASE que gerou o binário (1.4.2)
//   buildNumber      inteiro crescente exigido pelas lojas
//   runtimeVersion   muda só quando a parte NATIVA muda; é o que decide se uma
//                    atualização sem loja pode chegar a um aparelho
//
// Como o app opera até 7 dias offline, um aparelho pode estar alguns releases
// atrás. O selo mostra a versão que está DE FATO rodando, que pode ser a de
// uma atualização sem loja mais nova que o binário instalado.

import type { ExpoConfig } from 'expo/config';

const versaoRelease = process.env.TENAZ_VERSAO ?? '0.0.0-dev';
const commit = process.env.TENAZ_COMMIT ?? 'desconhecido';
const dataBuild = process.env.TENAZ_DATA_BUILD ?? '';
const ambiente = process.env.TENAZ_AMBIENTE ?? 'desenvolvimento';
// Base da API. 10.0.2.2 é o computador do desenvolvedor visto pelo emulador Android.
const apiUrl = process.env.TENAZ_API_URL ?? 'http://10.0.2.2:8080';

// As lojas exigem inteiro crescente. No CI é o número da execução; localmente,
// 1 basta, porque build local não vai para loja.
const numeroDoBuild = Number(process.env.TENAZ_NUMERO_BUILD ?? '1');

// A loja recusa versão com sufixo de pré-lançamento no campo `version`.
const versaoParaLoja = versaoRelease.split('-')[0] ?? '0.0.0';

const config: ExpoConfig = {
  name: 'Tenaz',
  slug: 'tenaz',
  // O esquema br.com.tenaz.app recebe o retorno do login (redirect_uri
  // br.com.tenaz.app:/oauth, seção 4.3 de docs/3_ESPECIFICACAO_APIS.md).
  scheme: ['br.com.tenaz.app', 'tenaz'],
  version: versaoParaLoja,
  orientation: 'portrait',
  userInterfaceStyle: 'automatic',
  // A nova arquitetura do React Native é a única desde o SDK 55: não há mais chave `newArchEnabled`.

  // Política de impressão digital: a versão de runtime muda só quando a parte
  // nativa muda, o que é exatamente a condição para a atualização sem loja
  // poder ser aplicada.
  runtimeVersion: { policy: 'fingerprint' },

  icon: './assets/icone.png',
  ios: {
    bundleIdentifier: 'br.com.tenaz.app',
    buildNumber: String(numeroDoBuild),
    supportsTablet: true,
  },

  android: {
    package: 'br.com.tenaz.app',
    versionCode: numeroDoBuild,
    adaptiveIcon: {
      foregroundImage: './assets/icone-adaptativo.png',
      backgroundColor: '#14202E',
    },
    permissions: ['CAMERA', 'ACCESS_FINE_LOCATION', 'ACCESS_COARSE_LOCATION'],
  },

  plugins: [
    'expo-router',
    // A tela de abertura é plugin desde o SDK 55 (a chave `splash` saiu da configuração).
    ['expo-splash-screen', { image: './assets/splash.png', resizeMode: 'contain', backgroundColor: '#14202E' }],
    'expo-secure-store',
    'expo-sqlite',
    'expo-web-browser',
    'expo-local-authentication',
  ],

  // Lido em tempo de execução pelo selo de versão e pela tela "Sobre".
  // `versaoRelease` guarda a versão completa, com o sufixo de pré-lançamento
  // que o campo `version` não aceita.
  extra: {
    versaoRelease,
    commit,
    dataBuild,
    ambiente,
    apiUrl,
  },
};

export default config;
