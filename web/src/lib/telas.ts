// Código da tela atual, enviado em X-Tenaz-Tela em toda chamada à API
// (RF-ADM-011): é o que a trilha grava como "tela de origem" e o que as chaves
// de captura (visualização, impressão, extração) usam.
//
// O catálogo é fechado e vive no backend (backend/internal/plataforma/
// auditoria/catalogo.go, formato `<modulo>.<tela>`); o teste telas.test.ts
// confere esta lista contra aquele arquivo. Rota sem tela no catálogo (como
// "Minha conta") não envia o cabeçalho.

/** Telas do catálogo fechado, na ordem do backend. */
export const TELAS_DO_CATALOGO = [
  'adm.empresa',
  'adm.unidades',
  'adm.unidade_detalhe',
  'adm.pessoas',
  'adm.pessoa_detalhe',
  'adm.usuarios',
  'adm.usuario_detalhe',
  'adm.perfis',
  'adm.auditoria',
  'adm.parametros',
  'ind.painel',
  'atv.ativos',
  'atv.ativo_detalhe',
  'med.medidores',
  'sol.solicitacoes',
  'sol.solicitacao_detalhe',
  'os.ordens',
  'os.ordem_detalhe',
  'pre.planos',
  'chk.modelos',
  'est.itens',
  'est.movimentacoes',
  'ter.fornecedores',
  'srv.conexoes',
  'prj.projetos',
  'ctr.contratos',
  'loc.locatarios',
  'loc.contratos',
  'fro.veiculos',
  'con.conformidade',
  'int.integracoes',
] as const;

export type CodigoDeTela = (typeof TELAS_DO_CATALOGO)[number];

/** Códigos de módulo das telas do catálogo, sem repetição, na ordem do catálogo. */
export const MODULOS_DO_CATALOGO: readonly string[] = [...new Set(TELAS_DO_CATALOGO.map((t) => t.split('.')[0] ?? ''))];

/** Tela do catálogo para um caminho da web, ou null quando a rota não tem tela no catálogo. */
export function telaDoCaminho(caminho: string): CodigoDeTela | null {
  const p = caminho.replace(/\/+$/, '') || '/';
  if (p === '/') return 'ind.painel';
  if (p === '/administracao/empresa') return 'adm.empresa';
  if (p === '/administracao/unidades' || p === '/administracao/unidades/nova') return 'adm.unidades';
  if (p.startsWith('/administracao/unidades/')) return 'adm.unidade_detalhe';
  if (p === '/administracao/pessoas' || p === '/administracao/pessoas/nova') return 'adm.pessoas';
  if (p.startsWith('/administracao/pessoas/')) return 'adm.pessoa_detalhe';
  if (p === '/administracao/usuarios' || p === '/administracao/usuarios/novo') return 'adm.usuarios';
  if (p.startsWith('/administracao/usuarios/')) return 'adm.usuario_detalhe';
  if (p === '/administracao/perfis') return 'adm.perfis';
  if (p === '/administracao/auditoria' || p === '/administracao/auditoria/integridade') return 'adm.auditoria';
  if (p === '/administracao/auditoria/capturas') return 'adm.parametros';
  return null;
}
