// Glossário bilíngue pt-BR <-> es-419 do vocabulário canônico da Tenaz.
//
// Fonte única das traduções aprovadas dos termos da interface. O documento 7
// (seção 10.5) cita este arquivo. Termo novo entra aqui antes de entrar nos
// catálogos de texto (regra 3 do agente tenaz-idiomas). O nome técnico
// (entidade, campo, estado, permissão, rota, evento) NÃO se traduz: a coluna
// `id` só liga o termo ao sumário executivo, seção 3.
//
// Fontes: F1 = glossário e blog de CMMS em espanhol (Tractian, eMaint, Cegid,
// Rockwell Automation es-ar, Capterra Argentina), que usam "orden de trabajo",
// "activo", "mantenimiento preventivo", "inventario"; F2 = normas ISO 55000
// (activo) e ISO 14224 (falla); F3 = uso corrente do mercado latino-americano de
// manutenção industrial e de frota; F4 = convenção do projeto para fechar uma
// ambiguidade, sem fonte externa. `revisar: true` marca o que um falante nativo
// da Argentina ou do Paraguai deve confirmar antes da primeira tela em produção.

export type FonteTermo = 'F1' | 'F2' | 'F3' | 'F4';

export interface TermoGlossario {
  /** Nome canônico do sumário (seção 3) ou grupo do termo. */
  id: string;
  pt: string;
  es: string;
  /** Alternativas descartadas, em espanhol; o motivo vai em `nota`. */
  descartados: string[];
  fonte: FonteTermo[];
  nota?: string;
  revisar?: boolean;
  /**
   * Teste de uso: expressões (sem distinção de caixa, aplicadas aos textos em
   * espanhol) que não podem aparecer em nenhum texto do catálogo es-419.
   */
  esProibidos?: string[];
  /**
   * Teste de consistência: quando um texto pt-BR casa `ptDetecta`, o texto
   * es-419 da mesma chave tem de casar `esExige`.
   */
  ptDetecta?: string;
  esExige?: string;
}

const t = (x: TermoGlossario): TermoGlossario => x;

export const glossario: readonly TermoGlossario[] = [
  // ---- Pontos sensíveis ------------------------------------------------
  t({
    id: 'ordem_servico',
    pt: 'ordem de serviço (OS)',
    es: 'orden de trabajo (OT)',
    descartados: ['orden de servicio', 'OS', 'ticket', 'pedido de trabajo'],
    fonte: ['F1', 'F3'],
    nota: '"Orden de trabajo" é o termo dos CMMS em espanhol e do uso industrial; "orden de servicio" soa a atendimento a cliente ou a venda e colide com o módulo Servicios (SRV). A sigla é OT. O prefixo do número (OS-AAAA-NNNNNN) é dado, não texto, e fica como está; confirmar com o dono do produto se vira OT-.',
    esProibidos: ['\\b[oó]rden(es)? de servicio\\b', '\\bOS\\b(?![.])'],
    ptDetecta: 'ordem de serviço|ordens de serviço',
    esExige: '[oó]rden(es)? de trabajo',
  }),
  t({
    id: 'solicitacao',
    pt: 'solicitação',
    es: 'solicitud',
    descartados: ['ticket', 'pedido', 'requerimiento', 'reclamo'],
    fonte: ['F1', 'F3'],
    nota: 'Forma completa: "solicitud de mantenimiento". A sigla SS do número coincide com "solicitud de servicio".',
    esProibidos: ['\\bticket(s)?\\b'],
  }),
  t({
    id: 'ativo',
    pt: 'ativo',
    es: 'activo',
    descartados: ['bien', 'equipo (como termo geral)', 'máquina'],
    fonte: ['F1', 'F2'],
    nota: 'ISO 55000 e os CMMS em espanhol usam "activo". "Equipo" fica só para o tipo de ativo `equipamento`.',
    ptDetecta: '\\bativos?\\b',
    esExige: '\\bactivos?\\b',
  }),
  t({
    id: 'encerrada',
    pt: 'encerrada / encerrar',
    es: 'cerrada / cerrar',
    descartados: ['finalizada', 'terminada', 'clausurada'],
    fonte: ['F1', 'F3'],
    nota: '`encerrada` (liberada, custos consolidados) é diferente de `concluida`. "Cerrar una OT" é o uso do mercado e fica reservado ao encerramento; "concluida" nunca é "cerrada".',
  }),
  t({
    id: 'concluida',
    pt: 'concluída / concluir',
    es: 'concluida / concluir',
    descartados: ['terminada', 'finalizada', 'completada'],
    fonte: ['F4'],
    nota: 'Mantém a distinção `concluida` x `encerrada` do documento 7, seção 10.2.',
  }),
  t({
    id: 'parado',
    pt: 'parado',
    es: 'detenido',
    descartados: ['parado', 'fuera de servicio', 'en paro', 'inoperativo'],
    fonte: ['F3'],
    nota: '"Parado" quer dizer "de pé" ou "desempregado" em parte da América Latina; "en paro" é regionalismo (México, Colômbia); "fuera de servicio" diz outra coisa (baixado). "Detenido" é entendido em toda a região ("equipo detenido").',
    revisar: true,
    esProibidos: ['\\bparad[oa]s?\\b'],
  }),
  t({
    id: 'em_implantacao',
    pt: 'em implantação / implantação',
    es: 'en implementación / implementación',
    descartados: ['en implantación', 'en instalación', 'en obra', 'en puesta en marcha'],
    fonte: ['F4'],
    nota: '"Implantación" é pouco usada fora da Espanha. "Instalación" é descartada porque o documento 7 já a evita. A ação `colocar_em_operacao` é "Poner en operación".',
    revisar: true,
    esProibidos: ['\\bimplantaci[oó]n\\b'],
  }),
  t({
    id: 'triagem',
    pt: 'triagem',
    es: 'evaluación',
    descartados: ['triaje', 'clasificación', 'revisión'],
    fonte: ['F4'],
    nota: '"Triaje" tem conotação médica. "Revisión" fica para a ação de auditoria `revisao`. Estado `em_triagem` = "En evaluación".',
    esProibidos: ['\\btriaje\\b'],
  }),
  t({
    id: 'locacao',
    pt: 'locação / locatário / locador',
    es: 'alquiler / arrendatario / arrendador',
    descartados: ['locación', 'locatario', 'locador', 'inquilino', 'renta'],
    fonte: ['F3'],
    nota: '"Alquiler" é neutro para equipamento e frota. "Locatario/locador" é o termo legal argentino, entendido em toda parte, mas "arrendatario/arrendador" é o par neutro do continente. A conformidade do Mercosul (P12) confirma antes do contrato de locação.',
    revisar: true,
    esProibidos: ['\\binquilinos?\\b'],
  }),
  t({
    id: 'cedido_tomado',
    pt: 'cedido / tomado',
    es: 'cedido en alquiler / tomado en alquiler',
    descartados: ['alquilado a', 'alquilado de'],
    fonte: ['F4'],
  }),
  t({
    id: 'almoxarifado',
    pt: 'almoxarifado',
    es: 'almacén',
    descartados: ['depósito', 'bodega', 'pañol'],
    fonte: ['F3'],
    nota: '"Depósito" e "pañol" são argentinos; "bodega" é chileno e colombiano. "Almacén" é neutro. O perfil `almoxarife` é "almacenista".',
  }),
  t({
    id: 'estoque',
    pt: 'estoque',
    es: 'inventario',
    descartados: ['stock', 'existencias'],
    fonte: ['F1'],
    nota: 'CMMS em espanhol usam "control de inventario". "Stock" é comum na Argentina, mas não neutro.',
    esProibidos: ['\\bstock\\b'],
  }),
  t({
    id: 'peca',
    pt: 'peça',
    es: 'repuesto',
    descartados: ['pieza', 'refacción', 'recambio'],
    fonte: ['F3'],
    nota: '"Refacción" é mexicano; "recambio" é peninsular.',
    esProibidos: ['\\brefacci[oó]n(es)?\\b'],
  }),
  t({
    id: 'falha',
    pt: 'falha',
    es: 'falla',
    descartados: ['fallo', 'avería', 'defecto'],
    fonte: ['F2', 'F3'],
    nota: 'A ISO 14224 em espanhol da América Latina traduz "failure" como "falla". "Fallo" e "avería" são peninsulares.',
    esProibidos: ['\\bfallos?\\b'],
  }),
  t({
    id: 'lista_verificacao',
    pt: 'checklist',
    es: 'lista de verificación',
    descartados: ['checklist', 'lista de chequeo', 'lista de control'],
    fonte: ['F3'],
    nota: '"Lista de chequeo" é colombiano; "lista de control" é peninsular.',
    esProibidos: ['\\bchecklists?\\b'],
  }),
  t({
    id: 'trilha_auditoria',
    pt: 'trilha de auditoria',
    es: 'pista de auditoría',
    descartados: ['rastro de auditoría', 'registro de auditoría', 'historial de auditoría'],
    fonte: ['F3'],
    nota: '"Pista de auditoría" é o termo de contabilidade e TI da região para audit trail.',
    ptDetecta: 'trilha de auditoria',
    esExige: 'pista de auditor[ií]a',
  }),
  t({
    id: 'entrar',
    pt: 'entrar / sair (sessão)',
    es: 'iniciar sesión / cerrar sesión',
    descartados: ['ingresar', 'entrar', 'salir', 'loguearse'],
    fonte: ['F3'],
    esProibidos: ['\\bingres(ar|a|e|o)\\b'],
    nota: '"Ingresar" é correto mas ambíguo com "ingresar datos"; "iniciar sesión" é o padrão dos sistemas. A ação de auditoria `login` é "Inicio de sesión".',
  }),
  t({
    id: 'email',
    pt: 'e-mail',
    es: 'correo electrónico',
    descartados: ['e-mail', 'email', 'mail'],
    fonte: ['F3'],
    esProibidos: ['\\be-?mail\\b'],
  }),
  t({
    id: 'senha',
    pt: 'senha',
    es: 'contraseña',
    descartados: ['clave', 'password', 'pass'],
    fonte: ['F3'],
    nota: '"Clave" é comum no Paraguai e na Argentina, mas "contraseña" é neutro.',
    esProibidos: ['\\bpassword\\b'],
  }),
  t({
    id: 'segundo_fator',
    pt: 'segundo fator',
    es: 'segundo factor',
    descartados: ['doble factor', 'verificación en dos pasos', '2FA'],
    fonte: ['F3'],
    nota: '"Aplicativo autenticador" = "app de autenticación".',
    ptDetecta: 'segundo fator',
    esExige: 'segundo factor',
  }),
  t({
    id: 'sem_conexao',
    pt: 'sem conexão',
    es: 'sin conexión',
    descartados: ['offline', 'fuera de línea', 'desconectado'],
    fonte: ['F3'],
    esProibidos: ['\\boffline\\b', 'fuera de l[ií]nea'],
  }),
  t({
    id: 'conta',
    pt: 'conta (cada empresa em que o usuário trabalha)',
    es: 'cuenta',
    descartados: ['tenant', 'inquilino', 'organización'],
    fonte: ['F3'],
    esProibidos: ['\\btenants?\\b'],
  }),
  t({
    id: 'link',
    pt: 'link',
    es: 'enlace',
    descartados: ['link', 'vínculo'],
    fonte: ['F3'],
    nota: '"Vínculo" é relação entre registros (pessoa e fornecedor), não endereço.',
    esProibidos: ['\\blinks?\\b'],
  }),
  t({
    id: 'tom_tratamento',
    pt: 'você / imperativo',
    es: 'usted (imperativo de usted) ou forma impessoal',
    descartados: ['tú', 'vos', 'imperativo de tú ("ingresa")', 'imperativo de vos ("ingresá")'],
    fonte: ['F4'],
    nota: 'Sem voseo e sem tuteo. Imperativo de "usted" ("Inicie", "Complete", "Elija"), como o pt-BR usa o imperativo sem pronome.',
    esProibidos: [
      '\\b(vos|tenés|podés|querés|sos|ingresá|hacé|elegí|confirmá|revisá|seleccioná|escribí|intentá|completá|usá|volvé)\\b',
      '\\b(tú|tus?|puedes|tienes|eres|quieres)\\b',
    ],
  }),

  t({
    id: 'vincular',
    pt: 'vincular / vinculada (pessoa e usuário)',
    es: 'vincular / vinculada',
    descartados: ['asociar', 'asociada', 'ligar'],
    fonte: ['F4'],
    nota: 'O verbo de relação entre registros (pessoa e usuário, perfil e fornecedor) é "vincular", como na interface. "Vínculo" (substantivo) é também o campo `vinculo` da pessoa (empregado, terceiro, autônomo). Decisão da revisão dos textos do servidor (13/11/2026).',
  }),
  t({
    id: 'tela',
    pt: 'tela',
    es: 'pantalla',
    descartados: ['ventana', 'vista', 'página'],
    fonte: ['F3'],
    nota: 'Uso corrente em software na região; mesmo termo da interface e dos textos de erro do servidor.',
  }),
  t({
    id: 'chave_idempotencia',
    pt: 'chave de idempotência',
    es: 'identificador de idempotencia',
    descartados: ['clave de idempotencia', 'llave de idempotencia'],
    fonte: ['F4'],
    revisar: true,
    nota: '"Clave" é senha em parte da região (ver `senha`); por isso o cabeçalho Idempotency-Key aparece em espanhol como "identificador". O nome do cabeçalho não se traduz.',
  }),
  t({
    id: 'desafio_mfa',
    pt: 'código (e não "desafio") do segundo fator',
    es: 'código (no "desafío") del segundo factor',
    descartados: ['desafío'],
    fonte: ['F4'],
    nota: '"Desafio" (challenge) é jargão da API e não aparece em texto ao usuário, em nenhum dos dois idiomas; usa-se "código". O nome técnico `desafio` continua no código.',
    esProibidos: ['\\bdesaf[ií]os?\\b'],
  }),

  // ---- Entidades -------------------------------------------------------
  t({ id: 'empresa', pt: 'empresa', es: 'empresa', descartados: ['compañía', 'organización'], fonte: ['F3'] }),
  t({ id: 'unidade', pt: 'unidade', es: 'unidad', descartados: ['sucursal', 'sede'], fonte: ['F3'], nota: 'Planta, filial, base ou loja do documento 0 viram "unidad" no menu; texto livre pode citar "planta" ou "sucursal".' }),
  t({ id: 'localizacao', pt: 'localização', es: 'ubicación', descartados: ['localización', 'sitio'], fonte: ['F3'] }),
  t({ id: 'classe_ativo', pt: 'classe de ativo', es: 'clase de activo', descartados: ['categoría de activo'], fonte: ['F3'], nota: '"Tipo de activo" é o campo `tipo` (equipamento, veículo, implemento, componente).' }),
  t({ id: 'criticidade', pt: 'criticidade', es: 'criticidad', descartados: ['importancia', 'gravedad'], fonte: ['F2', 'F3'] }),
  t({ id: 'medidor', pt: 'medidor', es: 'medidor', descartados: ['contador'], fonte: ['F3'] }),
  t({ id: 'leitura', pt: 'leitura', es: 'lectura', descartados: ['medición'], fonte: ['F3'], nota: '"Medición" fica para a medição do contrato com cliente (documento 7, 10.2).' }),
  t({ id: 'plano_manutencao', pt: 'plano de manutenção', es: 'plan de mantenimiento', descartados: ['programa de mantenimiento', 'rutina'], fonte: ['F1'] }),
  t({ id: 'gatilho_plano', pt: 'gatilho (do plano)', es: 'disparador (del plan)', descartados: ['gatillo', 'trigger', 'activador'], fonte: ['F4'], revisar: true, nota: '"Gatillo" é literal e soa a arma em parte da região.' }),
  t({ id: 'procedimento', pt: 'procedimento', es: 'procedimiento', descartados: ['instructivo'], fonte: ['F3'] }),
  t({ id: 'tarefa', pt: 'tarefa', es: 'tarea', descartados: ['actividad'], fonte: ['F1'], nota: '"Actividad" é do código de falha (modo, causa, atividade).' }),
  t({ id: 'apontamento', pt: 'apontamento (de horas)', es: 'registro de horas', descartados: ['apuntamiento', 'imputación de horas'], fonte: ['F3'] }),
  t({ id: 'pessoa', pt: 'pessoa', es: 'persona', descartados: ['empleado'], fonte: ['F3'] }),
  t({ id: 'usuario', pt: 'usuário', es: 'usuario', descartados: ['operador'], fonte: ['F3'] }),
  t({ id: 'perfil', pt: 'perfil', es: 'perfil', descartados: ['rol'], fonte: ['F3'], ptDetecta: '\\bperfil\\b|\\bperfis\\b', esExige: '\\bperfil(es)?\\b' }),
  t({ id: 'permissao', pt: 'permissão', es: 'permiso', descartados: ['autorización', 'privilegio'], fonte: ['F3'] }),
  t({ id: 'fornecedor', pt: 'fornecedor', es: 'proveedor', descartados: ['tercero', 'socio'], fonte: ['F3'] }),
  t({ id: 'subcontratada', pt: 'subcontratada', es: 'subcontratista', descartados: ['tercerizada', 'contratista'], fonte: ['F3'] }),
  t({ id: 'assistencia_tecnica', pt: 'assistência técnica', es: 'servicio técnico', descartados: ['asistencia técnica'], fonte: ['F3'], nota: '"Servicio técnico autorizado" é a forma comercial.' }),
  t({ id: 'oficina', pt: 'oficina', es: 'taller', descartados: ['oficina'], fonte: ['F3'], nota: '"Oficina" em espanhol é escritório.' }),
  t({ id: 'prestador', pt: 'prestador', es: 'prestador', descartados: ['tercerizado', 'proveedor de servicios'], fonte: ['F4'] }),
  t({ id: 'ordem_reparo', pt: 'ordem de reparo', es: 'orden de reparación', descartados: ['orden de arreglo'], fonte: ['F3'] }),
  t({ id: 'remessa_retorno', pt: 'remessa a terceiro / retorno de terceiro', es: 'envío a tercero / retorno de tercero', descartados: ['remesa', 'préstamo'], fonte: ['F4'], nota: '"Remesa" em espanhol é envio de dinheiro.' }),
  t({ id: 'pedido_compra', pt: 'pedido de compra', es: 'orden de compra', descartados: ['pedido de compra'], fonte: ['F3'] }),
  t({ id: 'requisicao_compra', pt: 'requisição de compra', es: 'requisición de compra', descartados: ['pedido interno', 'solicitud de compra'], fonte: ['F3'] }),
  t({ id: 'vistoria', pt: 'vistoria de entrega / de devolução', es: 'inspección de entrega / de devolución', descartados: ['check-in', 'check-out', 'peritaje'], fonte: ['F4'] }),
  t({ id: 'projeto', pt: 'projeto / etapa', es: 'proyecto / etapa', descartados: ['obra', 'fase'], fonte: ['F3'] }),
  t({ id: 'contrato', pt: 'contrato com cliente / de prestação / de locação', es: 'contrato con cliente / de prestación / de alquiler', descartados: ['contrato de servicio'], fonte: ['F3'] }),
  t({ id: 'cliente_atendido', pt: 'cliente atendido', es: 'cliente atendido', descartados: ['cliente final'], fonte: ['F4'] }),
  t({ id: 'documento_vencimento', pt: 'documento com vencimento', es: 'documento con vencimiento', descartados: ['documento vencible'], fonte: ['F3'] }),
  t({ id: 'capacitacao', pt: 'capacitação', es: 'capacitación', descartados: ['formación', 'certificación'], fonte: ['F3'] }),
  t({ id: 'regra', pt: 'regra / condição / efeito', es: 'regla / condición / efecto', descartados: ['automatización', 'flujo de trabajo'], fonte: ['F4'] }),
  t({ id: 'campo_extra', pt: 'campo extra', es: 'campo adicional', descartados: ['campo personalizado', 'campo custom'], fonte: ['F4'] }),
  t({ id: 'pacote_segmento', pt: 'pacote de segmento', es: 'paquete de segmento', descartados: ['plantilla', 'kit'], fonte: ['F4'] }),

  // ---- Estados da OS ---------------------------------------------------
  t({ id: 'os_aberta', pt: 'Aberta', es: 'Abierta', descartados: ['Nueva'], fonte: ['F1'] }),
  t({ id: 'os_planejada', pt: 'Planejada', es: 'Planificada', descartados: ['Planeada'], fonte: ['F1'], nota: '"Planeada" é mexicano.' }),
  t({ id: 'os_programada', pt: 'Programada', es: 'Programada', descartados: ['Agendada'], fonte: ['F1'] }),
  t({ id: 'os_em_execucao', pt: 'Em execução', es: 'En ejecución', descartados: ['En curso', 'En proceso'], fonte: ['F1'] }),
  t({ id: 'os_pausada', pt: 'Pausada', es: 'Pausada', descartados: ['Suspendida', 'En pausa'], fonte: ['F3'] }),
  t({ id: 'os_concluida', pt: 'Concluída', es: 'Concluida', descartados: ['Completada', 'Finalizada'], fonte: ['F4'] }),
  t({ id: 'os_encerrada', pt: 'Encerrada', es: 'Cerrada', descartados: ['Finalizada', 'Terminada'], fonte: ['F1'] }),
  t({ id: 'os_cancelada', pt: 'Cancelada', es: 'Cancelada', descartados: ['Anulada'], fonte: ['F3'] }),

  // ---- Solicitação -----------------------------------------------------
  t({ id: 'sol_nova', pt: 'Nova', es: 'Nueva', descartados: [], fonte: ['F3'] }),
  t({ id: 'sol_em_triagem', pt: 'Em triagem', es: 'En evaluación', descartados: ['En triaje'], fonte: ['F4'] }),
  t({ id: 'sol_aprovada', pt: 'Aprovada', es: 'Aprobada', descartados: ['Autorizada'], fonte: ['F3'] }),
  t({ id: 'sol_recusada', pt: 'Recusada', es: 'Rechazada', descartados: ['Denegada'], fonte: ['F3'] }),
  t({ id: 'sol_duplicada', pt: 'Duplicada', es: 'Duplicada', descartados: ['Repetida'], fonte: ['F3'] }),

  // ---- Status do ativo -------------------------------------------------
  t({ id: 'ativo_operando', pt: 'Operando', es: 'En operación', descartados: ['Funcionando'], fonte: ['F3'] }),
  t({ id: 'ativo_em_manutencao', pt: 'Em manutenção', es: 'En mantenimiento', descartados: ['En reparación'], fonte: ['F3'] }),
  t({ id: 'ativo_parado', pt: 'Parado', es: 'Detenido', descartados: ['Parado', 'En paro'], fonte: ['F3'], revisar: true }),
  t({ id: 'ativo_inativo', pt: 'Inativo', es: 'Inactivo', descartados: ['Dado de baja'], fonte: ['F3'], nota: '"Dado de baja" é baixa definitiva, outro conceito.' }),

  // ---- Prioridade e tipo de OS ----------------------------------------
  t({ id: 'prio_baixa', pt: 'Baixa', es: 'Baja', descartados: [], fonte: ['F3'] }),
  t({ id: 'prio_media', pt: 'Média', es: 'Media', descartados: ['Normal'], fonte: ['F3'] }),
  t({ id: 'prio_alta', pt: 'Alta', es: 'Alta', descartados: [], fonte: ['F3'] }),
  t({ id: 'prio_urgente', pt: 'Urgente', es: 'Urgente', descartados: ['Crítica', 'Emergencia'], fonte: ['F3'], nota: 'Como no pt-BR: "crítica" confunde com criticidade e "emergencia" é tipo de OS.' }),
  t({ id: 'tipo_preventiva', pt: 'Preventiva', es: 'Preventiva', descartados: ['Programada'], fonte: ['F1'] }),
  t({ id: 'tipo_preditiva', pt: 'Preditiva', es: 'Predictiva', descartados: ['Predecible'], fonte: ['F1'] }),
  t({ id: 'tipo_corretiva_planejada', pt: 'Corretiva planejada', es: 'Correctiva planificada', descartados: ['Correctiva programada'], fonte: ['F1'] }),
  t({ id: 'tipo_corretiva_emergencial', pt: 'Corretiva emergencial', es: 'Correctiva de emergencia', descartados: ['Correctiva urgente', 'Reactiva'], fonte: ['F1'] }),
  t({ id: 'tipo_inspecao_legal', pt: 'Inspeção legal', es: 'Inspección legal', descartados: ['Inspección reglamentaria'], fonte: ['F3'], revisar: true }),
  t({ id: 'tipo_melhoria', pt: 'Melhoria', es: 'Mejora', descartados: ['Mejoramiento'], fonte: ['F1'] }),
  t({ id: 'tipo_implantacao', pt: 'Implantação', es: 'Implementación', descartados: ['Implantación', 'Instalación'], fonte: ['F4'], revisar: true }),

  // ---- Perfis ----------------------------------------------------------
  t({ id: 'perfil_administrador', pt: 'Administrador', es: 'Administrador', descartados: [], fonte: ['F3'] }),
  t({ id: 'perfil_gestor', pt: 'Gestor', es: 'Gestor', descartados: ['Gerente'], fonte: ['F3'] }),
  t({ id: 'perfil_planejador', pt: 'Planejador', es: 'Planificador', descartados: ['Planeador'], fonte: ['F3'] }),
  t({ id: 'perfil_tecnico', pt: 'Técnico', es: 'Técnico', descartados: ['Mecánico'], fonte: ['F3'] }),
  t({ id: 'perfil_almoxarife', pt: 'Almoxarife', es: 'Almacenista', descartados: ['Pañolero', 'Bodeguero'], fonte: ['F3'] }),
  t({ id: 'perfil_comprador', pt: 'Comprador', es: 'Comprador', descartados: [], fonte: ['F3'] }),
  t({ id: 'perfil_plh', pt: 'PLH (pessoa legalmente habilitada)', es: 'Profesional habilitado', descartados: ['PLH'], fonte: ['F4'], revisar: true, nota: 'PLH é conceito da NR-13 brasileira; o equivalente na Argentina e no Paraguai depende de P12.' }),
  t({ id: 'perfil_motorista', pt: 'Motorista', es: 'Conductor', descartados: ['Chofer', 'Camionero'], fonte: ['F3'], nota: '"Chofer" é regional.' }),
  t({ id: 'perfil_solicitante', pt: 'Solicitante', es: 'Solicitante', descartados: [], fonte: ['F3'] }),
  t({ id: 'perfil_visualizador', pt: 'Visualizador', es: 'Visualizador', descartados: ['Observador', 'Lector'], fonte: ['F4'] }),
  t({ id: 'perfil_prestador', pt: 'Prestador', es: 'Prestador', descartados: [], fonte: ['F4'] }),
  t({ id: 'perfil_locatario', pt: 'Locatário', es: 'Arrendatario', descartados: ['Locatario', 'Inquilino'], fonte: ['F3'], revisar: true }),

  // ---- Módulos e nomes comerciais -------------------------------------
  t({ id: 'mod_adm', pt: 'Administração', es: 'Administración', descartados: [], fonte: ['F3'] }),
  t({ id: 'mod_atv', pt: 'Ativos', es: 'Activos', descartados: [], fonte: ['F1'] }),
  t({ id: 'mod_med', pt: 'Medidores', es: 'Medidores', descartados: [], fonte: ['F3'] }),
  t({ id: 'mod_sol', pt: 'Solicitações', es: 'Solicitudes', descartados: [], fonte: ['F3'] }),
  t({ id: 'mod_os', pt: 'Ordens de serviço', es: 'Órdenes de trabajo', descartados: ['Órdenes de servicio'], fonte: ['F1'] }),
  t({ id: 'mod_pre', pt: 'Planos de manutenção', es: 'Planes de mantenimiento', descartados: [], fonte: ['F1'] }),
  t({ id: 'mod_chk', pt: 'Checklists', es: 'Listas de verificación', descartados: [], fonte: ['F3'] }),
  t({ id: 'mod_est', pt: 'Estoque e compras', es: 'Inventario y compras', descartados: ['Stock y compras'], fonte: ['F1'] }),
  t({ id: 'mod_ter', pt: 'Terceiros', es: 'Terceros', descartados: [], fonte: ['F3'] }),
  t({ id: 'mod_srv', pt: 'Gestão de serviços do prestador', es: 'Gestión de servicios del prestador', descartados: [], fonte: ['F4'] }),
  t({ id: 'mod_prj', pt: 'Projetos e implantação', es: 'Proyectos e implementación', descartados: [], fonte: ['F4'] }),
  t({ id: 'mod_ctr', pt: 'Contratos com clientes', es: 'Contratos con clientes', descartados: [], fonte: ['F3'] }),
  t({ id: 'mod_loc', pt: 'Locação', es: 'Alquiler', descartados: ['Locación'], fonte: ['F3'], revisar: true }),
  t({ id: 'mod_fro', pt: 'Frota', es: 'Flota', descartados: ['Parque automotor'], fonte: ['F3'] }),
  t({ id: 'mod_con', pt: 'Conformidade', es: 'Cumplimiento normativo', descartados: ['Conformidad', 'Compliance'], fonte: ['F3'] }),
  t({ id: 'mod_ind', pt: 'Indicadores e relatórios', es: 'Indicadores e informes', descartados: ['KPIs y reportes'], fonte: ['F3'], nota: '"Informe" é o termo de relatório de negócio; "reporte" é anglicismo aceito, mas evitado.' }),
  t({ id: 'mod_int', pt: 'Integrações', es: 'Integraciones', descartados: [], fonte: ['F3'] }),
  t({ id: 'mod_ia', pt: 'Assistente', es: 'Asistente', descartados: [], fonte: ['F3'] }),
  t({ id: 'mod_mob', pt: 'App de campo', es: 'App de campo', descartados: ['Aplicación móvil'], fonte: ['F3'] }),
  t({ id: 'marca_slogan', pt: 'Cada ativo pronto para operar.', es: 'Cada activo listo para operar.', descartados: [], fonte: ['F4'] }),
  t({ id: 'marca_comerciais', pt: 'Tenaz Planta / Frota / Locação / Projetos / Serviços / Mobile', es: 'Tenaz Planta / Flota / Alquiler / Proyectos / Servicios / Mobile', descartados: [], fonte: ['F4'], nota: '"Tenaz" já é palavra espanhola (perseverante); Planta e Mobile ficam, os demais seguem o glossário.' }),
  t({ id: 'plano_comercial', pt: 'Essencial / Profissional / Enterprise', es: 'Esencial / Profesional / Enterprise', descartados: [], fonte: ['F3'] }),

  // ---- Trilha de auditoria --------------------------------------------
  t({ id: 'aud_inclusao', pt: 'Inclusão', es: 'Creación', descartados: ['Alta', 'Inclusión'], fonte: ['F3'], nota: '"Alta" (registro novo) é argentino e confunde com a prioridade Alta.' }),
  t({ id: 'aud_alteracao', pt: 'Alteração', es: 'Modificación', descartados: ['Alteración'], fonte: ['F3'], nota: '"Alteración" sugere dano ou adulteração.' }),
  t({ id: 'aud_exclusao', pt: 'Exclusão', es: 'Eliminación', descartados: ['Exclusión', 'Baja'], fonte: ['F3'] }),
  t({ id: 'aud_login', pt: 'Entrada / Falha de entrada / Saída', es: 'Inicio de sesión / Inicio de sesión fallido / Cierre de sesión', descartados: ['Ingreso', 'Egreso'], fonte: ['F3'] }),
  t({ id: 'aud_origem', pt: 'Origem', es: 'Origen', descartados: ['Procedencia'], fonte: ['F3'] }),
  t({ id: 'aud_sincronizacao', pt: 'Sincronização entre contas', es: 'Sincronización entre cuentas', descartados: [], fonte: ['F3'] }),

  // ---- País, moeda e identificação (ADR-013) ---------------------------
  t({ id: 'pais', pt: 'País', es: 'País', descartados: [], fonte: ['F4'] }),
  t({ id: 'pais_nomes', pt: 'Brasil / Argentina / Paraguai', es: 'Brasil / Argentina / Paraguay', descartados: ['Brazil'], fonte: ['F4'] }),
  t({ id: 'moeda', pt: 'Moeda', es: 'Moneda', descartados: ['Divisa'], fonte: ['F4'] }),
  t({ id: 'moeda_nomes', pt: 'Real brasileiro / Peso argentino / Guarani paraguaio', es: 'Real brasileño / Peso argentino / Guaraní paraguayo', descartados: ['Dólar'], fonte: ['F4'], revisar: true, nota: 'Nome oficial da moeda mais o código ISO entre parênteses. "Guaraní" leva acento em espanhol e fica no singular como nome da moeda; "Real brasileño" é a forma usada em espanhol (não "real brasilero").' }),
  t({ id: 'identificacao', pt: 'Identificação', es: 'Identificación', descartados: ['Documento fiscal'], fonte: ['F4'], nota: 'Campo genérico de empresa, unidade e pessoa; o tipo (CNPJ, CPF, CUIT, CUIL, DNI, RUC, CI) vem do país.' }),
  t({ id: 'identificacao_siglas', pt: 'CNPJ / CPF / CUIT / CUIL / DNI / RUC / CI', es: 'CNPJ / CPF / CUIT / CUIL / DNI / RUC / CI', descartados: [], fonte: ['F4'], nota: 'Siglas oficiais de cada país: não se traduzem.' }),
  t({ id: 'idioma_padrao', pt: 'Idioma padrão', es: 'Idioma predeterminado', descartados: ['Idioma por defecto'], fonte: ['F4'], revisar: true, nota: '"Predeterminado" é o termo dos sistemas operacionais em espanhol; "por defecto" também é entendido e fica descartado só por coerência com o texto de apoio.' }),
  t({ id: 'codigo_postal', pt: 'CEP (Brasil) / Código postal', es: 'CEP (Brasil) / Código postal', descartados: ['CP'], fonte: ['F4'], nota: 'No Brasil a sigla CEP continua; na Argentina e no Paraguai, "Código postal".' }),
  t({ id: 'subdivisao', pt: 'Estado (Brasil) / Província (Argentina) / Departamento (Paraguai)', es: 'Estado (Brasil) / Provincia (Argentina) / Departamento (Paraguay)', descartados: ['UF', 'Sigla da província ou do departamento'], fonte: ['F4'], nota: 'Seletor com o nome oficial (ISO 3166-2 completo no contrato); o rótulo muda pelo país da empresa. A revisão de idiomas apontou que o rótulo antigo soava estranho.' }),
  t({ id: 'localidade', pt: 'Município (Brasil) / Cidade ou localidade', es: 'Municipio (Brasil) / Ciudad o localidad', descartados: ['Localidad (sozinho)', 'Distrito'], fonte: ['F4'], revisar: true, nota: 'Na Argentina a unidade menor é a "localidad"; no Paraguai, "ciudad" ou "distrito". O rótulo com as duas palavras cobre os dois países; falante nativo confirma.' }),
];

/** Termos do glossário que viram regra de teste no catálogo es-419. */
export const regrasDoGlossario = glossario.filter(
  (g) => g.esProibidos !== undefined || (g.ptDetecta !== undefined && g.esExige !== undefined),
);
