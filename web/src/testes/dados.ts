// Dados de teste no formato do contrato (tipos de @tenaz/api-tipos).

import type {
  CapturaAuditoria,
  DiaVerificado,
  Empresa,
  Evento,
  Perfil,
  Permissao,
  Pessoa,
  ResumoEvento,
  ResumoPessoa,
  ResumoUnidade,
  ResumoUsuario,
  Sessao,
  Unidade,
  Usuario,
} from '../lib/tipos';

export const PERMISSOES_DE_ADMINISTRADOR = [
  'adm.exportar',
  'adm.gerenciar_empresa',
  'adm.gerenciar_parametros',
  'adm.gerenciar_perfis',
  'adm.gerenciar_pessoas',
  'adm.gerenciar_usuarios',
  'adm.ler',
  'adm.ler_auditoria',
  'adm.ler_pessoas',
  'adm.ler_usuarios',
];

export function sessao(extras: Partial<Sessao> = {}): Sessao {
  return {
    canal: 'web',
    csrf_token: 'csrf-de-teste',
    expira_em: '2026-11-20T12:00:00Z',
    escopo_acesso: 'interno',
    usuario: { id: 'u-1', nome: 'Ana Souza', email: 'ana@polar.com.br', pessoa_id: 'p-1', mfa_ativo: false, idioma: 'pt-BR' },
    empresa: { nome: 'Refrigeração Polar', plano: 'essencial', srv_habilitado: false, pais: 'BR', moeda: 'BRL', idioma_padrao: 'pt-BR' },
    perfis: ['administrador', 'gestor'],
    permissoes: PERMISSOES_DE_ADMINISTRADOR,
    empresas: [{ id: 'e-1', nome: 'Refrigeração Polar', atual: true }],
    mfa_pendente: false,
    ...extras,
  };
}

export function empresa(extras: Partial<Empresa> = {}): Empresa {
  return {
    id: 'e-1',
    razao_social: 'Refrigeração Polar Ltda',
    nome_fantasia: 'Polar',
    pais: 'BR',
    moeda: 'BRL',
    idioma_padrao: 'pt-BR',
    identificacao_tipo: 'cnpj',
    identificacao: '11222333000181',
    plano: 'essencial',
    srv_habilitado: false,
    mfa_obrigatorio: false,
    situacao: 'ativa',
    fuso_horario: 'America/Sao_Paulo',
    retencao_auditoria_anos: 5,
    teste_ate: null,
    versao: 3,
    criado_em: '2026-11-01T10:00:00Z',
    atualizado_em: '2026-11-02T10:00:00Z',
    ...extras,
  };
}

export function unidade(extras: Partial<Unidade> = {}): Unidade {
  return {
    id: 'un-1',
    codigo: 'MTZ',
    nome: 'Matriz Contagem',
    identificacao_tipo: null,
    identificacao: null,
    endereco: { municipio: 'Contagem', subdivisao: 'BR-MG' },
    fuso_horario: null,
    latitude: null,
    longitude: null,
    situacao: 'ativa',
    versao: 2,
    criado_em: '2026-11-01T10:00:00Z',
    atualizado_em: '2026-11-01T10:00:00Z',
    turnos: [
      { id: 't-1', unidade_id: 'un-1', dia_semana: 1, hora_inicio: '08:00', hora_fim: '17:00', nome: 'Comercial', versao: 1 },
      { id: 't-2', unidade_id: 'un-1', dia_semana: 2, hora_inicio: '08:00', hora_fim: '17:00', nome: null, versao: 4 },
    ],
    ...extras,
  };
}

export function resumoUnidade(u: Unidade = unidade()): ResumoUnidade {
  const { turnos, ...resumo } = u;
  void turnos;
  return resumo;
}

export function pessoa(extras: Partial<Pessoa> = {}): Pessoa {
  return {
    id: 'pe-1',
    nome: 'João Lima',
    identificacao_tipo: 'cpf',
    identificacao: '52998224725',
    data_nascimento: '1990-03-15',
    matricula: 'M-100',
    cargo: 'Técnico',
    telefone: null,
    email: null,
    vinculo: 'empregado',
    fornecedor_id: '',
    unidade_id: 'un-1',
    registro_profissional: null,
    situacao: 'ativa',
    usuario_id: '',
    versao: 5,
    criado_em: '2026-11-01T10:00:00Z',
    atualizado_em: '2026-11-01T10:00:00Z',
    ...extras,
  };
}

export function resumoPessoa(p: Pessoa = pessoa()): ResumoPessoa {
  return {
    id: p.id,
    nome: p.nome,
    identificacao_tipo: p.identificacao_tipo,
    identificacao_mascarada: '***.982.247-**',
    matricula: p.matricula,
    cargo: p.cargo,
    vinculo: p.vinculo,
    fornecedor_id: p.fornecedor_id,
    unidade_id: p.unidade_id,
    situacao: p.situacao,
    usuario_id: p.usuario_id,
    versao: p.versao,
    criado_em: p.criado_em,
  };
}

export function usuario(extras: Partial<Usuario> = {}): Usuario {
  return {
    id: 'u-2',
    pessoa_id: '',
    fornecedor_id: '',
    locatario_id: '',
    escopo_acesso: 'interno',
    nome: 'Carlos Lima',
    email: 'carlos@polar.com.br',
    telefone: null,
    idioma: 'pt-BR',
    moeda: 'BRL',
    pais: 'BR',
    custo_hora: null,
    unidade_padrao_id: '',
    situacao: 'ativo',
    perfis: [{ perfil_id: 'pf-1', codigo: 'tecnico', nome: 'Técnico', perfil_origem: 'tecnico', unidade_id: 'un-1' }],
    permissoes_avulsas: [],
    versao: 7,
    criado_em: '2026-11-01T10:00:00Z',
    atualizado_em: '2026-11-01T10:00:00Z',
    ...extras,
  };
}

export function resumoUsuario(u: Usuario = usuario()): ResumoUsuario {
  return {
    id: u.id,
    nome: u.nome,
    email: u.email,
    telefone: u.telefone,
    pessoa_id: u.pessoa_id,
    fornecedor_id: u.fornecedor_id,
    locatario_id: u.locatario_id,
    situacao: u.situacao,
    perfis: (u.perfis ?? []).map((p) => p.codigo),
    versao: u.versao,
    criado_em: u.criado_em,
  };
}

export const PERFIS: Perfil[] = [
  { id: 'pf-0', codigo: 'administrador', nome: 'Administrador', padrao: true, perfil_origem: null, permissoes: ['adm.ler', 'os.ler'], versao: 1 },
  { id: 'pf-1', codigo: 'tecnico', nome: 'Técnico', padrao: true, perfil_origem: null, permissoes: ['os.ler', 'os.executar'], versao: 1 },
  { id: 'pf-2', codigo: 'prestador', nome: 'Prestador', padrao: true, perfil_origem: null, permissoes: ['os.ler'], versao: 1 },
];

export const PERMISSOES: Permissao[] = [
  { codigo: 'os.encerrar', descricao: 'Liberar OS', grupo: 'OS', fase: 'MVP', sensivel: true },
  { codigo: 'os.ler', descricao: 'Ler OS', grupo: 'OS', fase: 'MVP', sensivel: false },
  { codigo: 'adm.ler', descricao: 'Ler empresa', grupo: 'ADM', fase: 'MVP', sensivel: false },
];

export function resumoEvento(extras: Partial<ResumoEvento> = {}): ResumoEvento {
  return {
    id: 'ev-1',
    ocorrido_em: '2026-11-12T14:05:00Z',
    ator_tipo: 'usuario',
    usuario_id: '11111111-2222-3333-4444-555555555555',
    usuario_nome: 'Ana Souza',
    token_api_id: '',
    acao: 'alteracao',
    resultado: 'sucesso',
    entidade: 'pessoa',
    entidade_id: 'pe-1',
    entidade_chave: 'João Lima',
    campos_alterados: ['cargo'],
    modulo: 'adm',
    tela: 'adm.pessoa_detalhe',
    sensivel: false,
    origem: 'web',
    ip: '203.0.113.7',
    dispositivo_id: '',
    request_id: 'req-1',
    conexao_empresa_id: '',
    ocorrido_no_dispositivo_em: null,
    ...extras,
  };
}

export function evento(extras: Partial<Evento> = {}): Evento {
  return {
    ...resumoEvento(),
    antes: { cargo: 'Auxiliar' },
    depois: { cargo: 'Técnico' },
    detalhe: null,
    hash: 'ab'.repeat(32),
    user_agent: 'teste',
    ...extras,
  };
}

export function pagina<T>(dados: T[], proximoCursor: string | null = null) {
  return { dados, paginacao: { proximo_cursor: proximoCursor, tem_mais: proximoCursor !== null, limite: 50 } };
}

export function diaVerificado(extras: Partial<DiaVerificado> = {}): DiaVerificado {
  return { dia: '2026-11-12', estado: 'integro', motivo: null, evento_id: '', ...extras };
}

export function captura(extras: Partial<CapturaAuditoria> = {}): CapturaAuditoria {
  return {
    tela: 'adm.pessoas',
    modulo: 'adm',
    nome: 'Pessoas',
    visualizacao: false,
    impressao: true,
    extracao: true,
    configurada: false,
    ...extras,
  };
}
