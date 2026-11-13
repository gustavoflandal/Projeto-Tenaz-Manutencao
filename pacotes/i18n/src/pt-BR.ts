// Textos da interface em português do Brasil, idioma de origem da Tenaz.
//
// Nenhum texto de interface fica fora deste arquivo (seção 8.4 de
// docs/5_ARQUITETURA_SOFTWARE.md). Tom de voz e termos: seção 10 de
// docs/7_PADROES_DESIGN.md. Caixa de frase em títulos, botões e menus; sem
// exclamação; erro diz o que aconteceu e o que fazer.

export const ptBR = {
  app: {
    nome: 'Tenaz',
    titulo: 'Tenaz · Gestão de Manutenção e Frota',
  },

  ambiente: {
    producao: 'Produção',
    homologacao: 'Homologação',
    desenvolvimento: 'Desenvolvimento',
  },

  // Selo de versão (componente 6.19, RNF-039)
  selo: {
    texto: 'v{{versao}}',
    textoComAmbiente: 'v{{versao}} · {{ambiente}}',
    nomeAcessivel: 'Versão {{versao}}. Abrir informações sobre o sistema',
    nomeAcessivelComAmbiente: 'Versão {{versao}}, ambiente {{ambiente}}. Abrir informações sobre o sistema',
  },

  // Diálogo "Sobre" (componente 6.19)
  sobre: {
    titulo: 'Sobre a Tenaz',
    descricao: 'Versão em uso e dados do build, para informar ao suporte.',
    web: 'Web',
    api: 'API',
    dataBuild: 'Data do build',
    commit: 'Commit',
    ambiente: 'Ambiente',
    versao: 'v{{versao}}',
    apiConsultando: 'Consultando a API',
    apiSemResposta: 'A API não respondeu. Verifique a conexão e abra de novo.',
    apiUltimaVista: 'Sem resposta · última vista v{{versao}}',
    dataBuildAusente: 'Build local, sem data',
    commitAusente: 'Não informado',
    versaoNova: 'Há uma versão mais nova. Recarregue a página para atualizar.',
    recarregar: 'Recarregar',
    copiar: 'Copiar informações',
    copiado: 'Informações da versão copiadas.',
    copiaFalhou: 'Não foi possível copiar. Selecione as informações e copie à mão.',
    fechar: 'Fechar',
    linhaCopia: '{{rotulo}}: {{valor}}',
  },

  tema: {
    rotulo: 'Tema',
    botao: 'Tema: {{tema}}. Trocar tema',
    sistema: 'Igual ao sistema',
    claro: 'Claro',
    escuro: 'Escuro',
    sol: 'Sol (alto contraste)',
  },

  // Tela de login (seção 5.6)
  entrar: {
    titulo: 'Entrar',
    email: 'E-mail',
    senha: 'Senha',
    mostrarSenha: 'Mostrar senha',
    ocultarSenha: 'Ocultar senha',
    botao: 'Entrar',
    indisponivel: 'A entrada com e-mail e senha ainda não está disponível nesta versão.',
  },

  // Casca da tela principal (seção 5)
  layout: {
    pularConteudo: 'Pular para o conteúdo',
    menuPrincipal: 'Menu principal',
    recolherMenu: 'Recolher menu',
    expandirMenu: 'Expandir menu',
  },

  menu: {
    painel: 'Painel',
  },

  painel: {
    titulo: 'Painel',
    vazio: 'Os indicadores da empresa aparecem aqui depois do cadastro dos primeiros ativos e ordens de serviço.',
  },

  naoEncontrada: {
    titulo: 'Página não encontrada',
    texto: 'Este endereço não existe ou mudou. Volte ao painel pelo menu.',
    acao: 'Ir para o painel',
  },
} as const;

/** Forma dos textos: os outros idiomas, quando houver, seguem esta. */
export type Textos = typeof ptBR;
