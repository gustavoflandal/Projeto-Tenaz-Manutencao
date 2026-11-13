// Ligação com os módulos nativos do Expo. Fica fora da lógica para os testes
// rodarem no Node sem Expo (o mesmo motivo de sync/deposito-expo.ts).

import * as Application from 'expo-application';
import Constants from 'expo-constants';
import * as Crypto from 'expo-crypto';
import * as LocalAuthentication from 'expo-local-authentication';

import { type AutenticadorLocal, BloqueioLocal } from './bloqueio';
import { CofreExpo } from './cofre-expo';
import type { Cripto } from './pkce';
import { GerenteDeSessoes } from './sessoes';
import { type InformacoesDoApp, lerInformacoes } from './versao';

const criptoExpo: Cripto = {
  bytesAleatorios: (n) => Crypto.getRandomBytes(n),
  sha256: async (entrada) => {
    const bytes = await Crypto.digest(Crypto.CryptoDigestAlgorithm.SHA256, new TextEncoder().encode(entrada));
    return new Uint8Array(bytes);
  },
};

const autenticadorExpo: AutenticadorLocal = {
  disponivel: async () =>
    (await LocalAuthentication.hasHardwareAsync()) && (await LocalAuthentication.isEnrolledAsync()),
  autenticar: async (motivo) => {
    // Permite o PIN do aparelho como alternativa à biometria.
    const r = await LocalAuthentication.authenticateAsync({
      promptMessage: motivo,
      cancelLabel: 'Cancelar',
      disableDeviceFallback: false,
    });
    return r.success;
  },
};

/** Base da API, de app.config.ts (TENAZ_API_URL). */
export function urlDaApi(): string {
  const url = Constants.expoConfig?.extra?.apiUrl;
  return typeof url === 'string' && url !== '' ? url : 'http://10.0.2.2:8080';
}

export function informacoesDoApp(): InformacoesDoApp {
  // `Updates` (atualização sem loja) entra com expo-updates, quando houver servidor de atualização.
  return lerInformacoes(
    Constants.expoConfig?.extra,
    { versao: Application.nativeApplicationVersion, build: Application.nativeBuildVersion },
    null,
  );
}

const cofre = new CofreExpo();

export const gerente = new GerenteDeSessoes({ cofre, cripto: criptoExpo, apiUrl: urlDaApi() });
export const bloqueio = new BloqueioLocal(cofre, autenticadorExpo);
