# Skills do Projeto — Tenaz

**Versão:** 0.3 · **Data:** 13/11/2026 · **Status:** Em elaboração
**Produto:** Tenaz · Gestão de Manutenção e Frota

Uma skill é o procedimento obrigatório de uma tarefa recorrente deste projeto. Ela existe para que o mesmo trabalho seja feito do mesmo jeito, com os mesmos portões, independentemente de quem o faça e de quando. A definição executável de cada skill está em [../../.claude/skills/](../../.claude/skills/).

**Diferença entre agente e skill.** O agente é *quem* faz, com uma especialidade e um conjunto de documentos sob sua guarda. A skill é *como* se faz, passo a passo, com o que não pode faltar. Um agente usa várias skills; uma skill é usada por vários agentes.

---

## 1. Catálogo

| Skill | Quando dispara | Agentes que mais usam |
|---|---|---|
| [`tenaz-contexto`](../../.claude/skills/tenaz-contexto/SKILL.md) | Começo de qualquer sessão ou tarefa no repositório | Todos |
| [`tenaz-implementar-requisito`](../../.claude/skills/tenaz-implementar-requisito/SKILL.md) | Há um `RF-` ou `RNF-` para entregar | `tenaz-produto`, `tenaz-backend-go` |
| [`tenaz-tabela-nova`](../../.claude/skills/tenaz-tabela-nova/SKILL.md) | Nasce ou muda uma tabela, coluna ou política | `tenaz-banco-dados`, `tenaz-isolamento` |
| [`tenaz-endpoint-novo`](../../.claude/skills/tenaz-endpoint-novo/SKILL.md) | Nasce ou muda uma rota, um erro ou um evento de webhook | `tenaz-api-contrato` |
| [`tenaz-operacao-offline`](../../.claude/skills/tenaz-operacao-offline/SKILL.md) | Uma operação precisa funcionar sem conexão | `tenaz-mobile-offline` |
| [`tenaz-campo-compartilhado`](../../.claude/skills/tenaz-campo-compartilhado/SKILL.md) | Um campo pode trafegar entre contas conectadas | `tenaz-conexao-contas`, `tenaz-isolamento` |
| [`tenaz-tela-nova`](../../.claude/skills/tenaz-tela-nova/SKILL.md) | Nasce ou muda uma tela na web ou no app | `tenaz-design-sistema`, `tenaz-frontend-web`, `tenaz-mobile-offline` |
| [`tenaz-evidencia-visual`](../../.claude/skills/tenaz-evidencia-visual/SKILL.md) | Uma tela entregue precisa de prova | Agentes de interface |
| [`tenaz-indicador`](../../.claude/skills/tenaz-indicador/SKILL.md) | Há um KPI para calcular ou exibir | `tenaz-dominio-manutencao`, `tenaz-integracoes-dados` |
| [`tenaz-prompt-etapa`](../../.claude/skills/tenaz-prompt-etapa/SKILL.md) | Uma etapa vai começar | Orquestrador, com o dono do produto |
| [`tenaz-fechar-etapa`](../../.claude/skills/tenaz-fechar-etapa/SKILL.md) | Uma etapa vai terminar | `tenaz-documentacao` |
| [`tenaz-release`](../../.claude/skills/tenaz-release/SKILL.md) | Vai sair uma versão em produção | `tenaz-plataforma-sre` |

---

## 2. Para que serve cada uma

### 2.1 `tenaz-contexto`

Orienta qualquer sessão nova: a ordem das fontes da verdade, o documento dono de cada assunto, o vocabulário canônico que não se improvisa, o que está decidido, o que é pendência do dono do produto e qual agente chamar para quê. É a primeira skill de toda tarefa, inclusive de tarefas que parecem triviais, porque o custo de usar um nome fora do vocabulário é propagar o erro por documentos, banco, API, web e app ao mesmo tempo.

### 2.2 `tenaz-implementar-requisito`

Leva um requisito do texto até o aceite, passando pelos 9 itens do critério de pronto do [roadmap](../6_ROADMAP.md), seção 3.1. Define a ordem: ler o requisito e as regras de negócio que o governam, conferir o modelo de dados, conferir o contrato, implementar, testar no nível certo, provar o isolamento, medir o desempenho, capturar a evidência visual, documentar. É a skill que impede o atalho mais comum: implementar o caminho feliz e esquecer validação, auditoria, offline e acessibilidade.

### 2.3 `tenaz-tabela-nova`

Toda tabela de negócio nasce com `empresa_id`, RLS habilitado **e** forçado, política de empresa com `USING` e `WITH CHECK`, política `escopo_externo` restritiva e, quando o dado puder trafegar entre contas, classificação na lista fechada da conexão. A skill inclui a entrada na suíte de vazamento e a migração no padrão expandir e contrair. Uma tabela que entra fora deste procedimento quebra o build, por desenho.

### 2.4 `tenaz-endpoint-novo`

Mantém o contrato único: tipos Go anotados, OpenAPI regenerado e versionado, comparação com a branch principal para barrar mudança incompatível, erro no formato RFC 9457 com código estável, paginação por cursor, idempotência por cabeçalho, concorrência por ETag, escopo de acesso declarado e tipos TypeScript regenerados. Web e app consomem o mesmo contrato da API pública; não existe API interna paralela.

### 2.5 `tenaz-operacao-offline`

Transforma uma operação em operação de campo: comando na fila local, identificador gerado no aparelho, idempotência no servidor, numeração provisória quando houver número legível, resolução de conflito definida antes de escrever código e teste automatizado de desconexão, ação, reconexão e conflito. O app trata a conexão como exceção, não como estado normal.

### 2.6 `tenaz-campo-compartilhado`

Todo campo novo de OS ou de ativo é classificado como compartilhável ou nunca compartilhado antes de existir. A skill cobre a lista fechada de campos, o grupo opcional, a transação separada por conta, a auditoria nos dois lados e o acréscimo do caso na suíte de três contas. É o procedimento que sustenta o risco R5 do [sumário](../0_SUMARIO_EXECUTIVO.md).

### 2.7 `tenaz-tela-nova`

Tela que usa só os tokens da marca, componentes do [documento 7](../7_PADROES_DESIGN.md), contraste conferido, estados de carregamento, vazio e erro escritos, microcopy no tom de voz definido, navegação por teclado, selo de versão não encoberto e ajuda contextual mínima. Termina com a evidência visual capturada.

### 2.8 `tenaz-evidencia-visual`

Define como a prova de uma tela é produzida e guardada em [../screenshots/](../screenshots/): o que capturar, em que largura, em que tema, com que massa de dados, com que nome. A evidência serve ao aceite, à documentação, ao manual de operação e à comparação entre versões.

### 2.9 `tenaz-indicador`

Implementa um indicador exatamente pela fórmula do [documento 4](../4_FLUXOS_PROCESSO.md), com a base de tempo correta, o tratamento de OS sem ativo, o efeito das paradas e o texto de "como calculamos" exigido por RF-IND-002. Indicador sem explicação visível ao usuário não é indicador entregue.

### 2.10 `tenaz-prompt-etapa`

Escreve o prompt da etapa em [../prompts/](../prompts/) antes de qualquer trabalho: objetivo, requisitos pelos IDs, documentos a ler, agentes previstos, portão de saída, evidência esperada e o que está explicitamente fora do escopo. Um prompt bom é o que permite a outra pessoa, meses depois, entender o que se pediu e por quê.

### 2.11 `tenaz-fechar-etapa`

Fecha a etapa conforme [CLAUDE.md](../../CLAUDE.md), seção 3: documentos afetados com versão subida e linha no histórico, `CHANGELOG.md` na seção `[Não lançado]`, registro em [../etapas/](../etapas/) no formato fixo, entrega consolidada em [../entregas/](../entregas/README.md) quando a etapa é aceita, sumário quando o escopo mudar, verificador sem erro e PDFs regerados.

### 2.12 `tenaz-release`

Fecha uma versão em produção conforme RNF-039 e RNF-040 e a seção 3.3 do [roadmap](../6_ROADMAP.md): tag anotada, seção do CHANGELOG, versão igual em web, app, API e workers, injetada pelo build, nunca digitada, verificações de CI e confirmação depois do deploy.

---

## 3. Regras das skills

1. **Skill é procedimento, não conhecimento.** O que uma skill contém é a sequência de passos e os portões. O conteúdo técnico mora nos documentos de `docs/`, e a skill aponta para eles.
2. **Uma skill não duplica documento.** Quando a skill e o documento divergem, vale o documento, e a skill é corrigida no mesmo commit.
3. **Skill nova nasce de repetição observada.** Tarefa feita três vezes de formas diferentes é candidata a skill.
4. **Toda skill termina em verificação.** Se o último passo não é um comando que se roda ou uma evidência que se olha, a skill está incompleta.
5. **Skill alterada atualiza este documento** no mesmo commit, e a mudança entra no registro da etapa.

---

**Histórico de versões**

- v0.1 (12/11/2026): primeira versão. Doze skills, com gatilho, finalidade e as cinco regras de manutenção.
- v0.2 (12/11/2026): `tenaz-fechar-etapa` passa a escrever a entrega consolidada em `docs/entregas/` quando a etapa é aceita (pedido do dono do produto na etapa E01).
- v0.3 (13/11/2026): `tenaz-fechar-etapa` passa a fazer o commit no branch da etapa e a abrir o pull request para `main` (modelo de execução, seção 4.1).
