// Cofre sobre expo-secure-store (Keychain e Keystore).
//
// As chaves do SecureStore só aceitam letras, números, ".", "-" e "_": por isso
// as chaves do app usam ponto (tenaz.credencial), nunca dois-pontos.

import * as SecureStore from 'expo-secure-store';

import type { Cofre } from './cofre';

export class CofreExpo implements Cofre {
  ler(chave: string): Promise<string | null> {
    return SecureStore.getItemAsync(chave);
  }

  gravar(chave: string, valor: string): Promise<void> {
    // WHEN_UNLOCKED_THIS_DEVICE_ONLY: não vai para backup nem para outro
    // aparelho, e só abre com o aparelho desbloqueado. O sync em segundo plano
    // do campo roda com o aparelho em uso, então a restrição não atrapalha.
    return SecureStore.setItemAsync(chave, valor, {
      keychainAccessible: SecureStore.WHEN_UNLOCKED_THIS_DEVICE_ONLY,
    });
  }

  apagar(chave: string): Promise<void> {
    return SecureStore.deleteItemAsync(chave);
  }
}
