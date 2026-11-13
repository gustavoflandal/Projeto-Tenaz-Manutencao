// Gera os tipos TypeScript do contrato OpenAPI (backend/api/openapi.yaml).
//
// O contrato é gerado dos tipos Go (seção 2.5 de docs/5_ARQUITETURA_SOFTWARE.md);
// estes tipos são gerados do contrato. Web e app importam daqui, e um campo
// renomeado no Go quebra a compilação do TypeScript no mesmo pull request.
//
// Uso: node gerar.mjs           escreve src/esquema.ts
//      node gerar.mjs --conferir falha se o arquivo versionado divergir
//
// Pela raiz: bash ferramentas/tarefas.sh api-tipos

import { readFileSync, writeFileSync } from 'node:fs';
import { dirname, resolve } from 'node:path';
import { fileURLToPath, pathToFileURL } from 'node:url';

import openapiTS, { astToString } from 'openapi-typescript';

const AQUI = dirname(fileURLToPath(import.meta.url));

export const CONTRATO = resolve(AQUI, '..', '..', 'backend', 'api', 'openapi.yaml');
export const DESTINO = resolve(AQUI, 'src', 'esquema.ts');

const CABECALHO = `// Gerado por pacotes/api-tipos/gerar.mjs a partir de backend/api/openapi.yaml.
// Não edite à mão: mude o contrato no Go, rode \`bash ferramentas/tarefas.sh contrato\`
// e depois \`bash ferramentas/tarefas.sh api-tipos\`.

`;

/** Tipos gerados do contrato, como texto. */
export async function gerarTipos(contrato = CONTRATO) {
  const ast = await openapiTS(pathToFileURL(contrato), { alphabetize: true });
  return `${CABECALHO}${astToString(ast)}`;
}

const executadoDireto = process.argv[1] && resolve(process.argv[1]) === fileURLToPath(import.meta.url);
if (executadoDireto) {
  const gerado = await gerarTipos();
  if (process.argv.includes('--conferir')) {
    if (readFileSync(DESTINO, 'utf8').replace(/\r\n/g, '\n') !== gerado) {
      console.error('pacotes/api-tipos/src/esquema.ts difere do contrato. Rode: bash ferramentas/tarefas.sh api-tipos');
      process.exit(1);
    }
    console.log('tipos da API em dia');
  } else {
    writeFileSync(DESTINO, gerado, 'utf8');
    console.log('pacotes/api-tipos/src/esquema.ts gerado');
  }
}
