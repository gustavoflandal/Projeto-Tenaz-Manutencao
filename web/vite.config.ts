import { writeFileSync } from 'node:fs';
import { resolve } from 'node:path';

import { defineConfig, type Plugin } from 'vite';

// A versão NUNCA é digitada em arquivo de código: vem do CI, derivada da tag
// git por ferramentas/versao.sh (RNF-039 e seção 9.5 de
// docs/5_ARQUITETURA_SOFTWARE.md).
const versao = process.env.TENAZ_VERSAO ?? '0.0.0-dev';
const commit = process.env.TENAZ_COMMIT ?? 'desconhecido';
const dataBuild = process.env.TENAZ_DATA_BUILD ?? '';
const ambiente = process.env.TENAZ_AMBIENTE ?? 'desenvolvimento';

/**
 * Publica `versao.json` junto dos arquivos estáticos.
 *
 * Serve a duas coisas: a verificação do CI, que confere se a versão publicada
 * é a da tag, e a própria aplicação, que compara o arquivo com a versão
 * embutida para avisar o usuário quando há versão mais nova — necessário
 * porque a web fica aberta por horas no mesmo pregão do navegador.
 */
function publicarVersao(): Plugin {
  return {
    name: 'tenaz-publicar-versao',
    apply: 'build',
    closeBundle() {
      const destino = resolve(__dirname, 'dist', 'versao.json');
      const corpo = JSON.stringify(
        { versao, commit, data_build: dataBuild, ambiente },
        null,
        2,
      );
      writeFileSync(destino, `${corpo}\n`, 'utf8');
    },
  };
}

export default defineConfig({
  plugins: [publicarVersao()],
  define: {
    __TENAZ_VERSAO__: JSON.stringify(versao),
    __TENAZ_COMMIT__: JSON.stringify(commit),
    __TENAZ_DATA_BUILD__: JSON.stringify(dataBuild),
    __TENAZ_AMBIENTE__: JSON.stringify(ambiente),
  },
  build: {
    outDir: 'dist',
    sourcemap: true,
  },
  server: {
    port: 5173,
    proxy: {
      // A web consome os mesmos endpoints da API pública. Não existe API
      // interna paralela (princípio 4 da arquitetura).
      '/v1': { target: 'http://127.0.0.1:8080', changeOrigin: true },
    },
  },
});
