# 013 — Conta e acesso (etapa E01)

**Data:** 2026-11-12 · **Tipo:** implementação
**Pedido por:** dono do produto · **Versão do produto:** —
**Situação:** em andamento. Este registro é atualizado a cada mudança, e não só no fechamento: a documentação é idempotente ao código ([CLAUDE.md](../../CLAUDE.md), seção 3.2).

## Pedido

"Leia os documentos da pasta e a memória da sessão e continue com o projeto." Com a E00 fechada, a etapa seguinte do plano é **E01 — Conta e acesso**, conforme [docs/prompts/E01_conta-e-acesso.md](../prompts/E01_conta-e-acesso.md). No início da sessão, o dono do produto pediu também que a infraestrutura fosse levantada no Docker.

Durante a etapa, o dono do produto determinou que **a documentação seja obrigatoriamente idempotente ao código**: todo fix, adendo, inclusão, exclusão ou alteração documentado, sem descompasso. A regra entrou em [CLAUDE.md](../../CLAUDE.md), seção 3.2.

## Decisões

### Do dono do produto

Registradas uma a uma na seção "Ajustes durante a etapa" do [prompt de E01](../prompts/E01_conta-e-acesso.md) e no [sumário 0.6](../0_SUMARIO_EXECUTIVO.md), seções 3.2 e 3.6.

| Assunto | Decisão |
|---|---|
| Tabelas de autenticação, ausentes do documento 2 | Nomes `sessao`, `token_renovacao`, `codigo_autorizacao`, `desafio_mfa`, `token_email` e os campos `identidade.falhas_login`, `bloqueada_ate`, `mfa_ativado_em` |
| Calendário da unidade (RF-ADM-002) | Tabelas próprias `turno_unidade` e `feriado_unidade`; os turnos são informados na criação, sem padrão imposto pelo sistema |
| Login por telefone e código | Adiado: pendência **P11** (provedor de SMS) |
| Catálogo de permissões | O catálogo de RF-ADM-005 é canônico por uma linha no sumário 3.6 |
| Permissão avulsa | Entidade `usuario_permissao` |
| Permissões do `administrador` | Todas, menos `pessoa.saude.ler`, `os.encerrar`, `est.decidir_destino` e `ter.manter_catalogo_proprio` |
| Situação do usuário | `convidado`, `ativo`, `inativo` |
| Responsável pela conta | Nasce com `administrador` e `gestor` |
| MFA obrigatório | Entra em E01, configurável por empresa num único interruptor (`empresa.mfa_obrigatorio`) |
| CNPJ da empresa | O `administrador` altera; a conexão entre contas trata a troca na etapa dela (proposta RN-ADM-031 recusada) |
| Pessoa | Vínculo `empregado`, `terceiro`, `autonomo`; situação só `ativa` ou `inativa`; pessoa inativa nunca é excluída e fica ligada ao histórico |
| Entregas consolidadas | Nova pasta [docs/entregas/](../entregas/README.md), com uma entrega por etapa aceita e em `main`. Registradas a concepção (C0) e a fundação técnica (E00); a E01 entra quando for aceita |
| Documentação | Obrigatoriamente idempotente ao código: toda mudança documentada no mesmo passo ([CLAUDE.md](../../CLAUDE.md), seção 3.2) |
| Trilha de auditoria configurável (ampliação de RF-ADM-011, dentro de E01) | Tela, módulo, tipos de operação (login, logout, inclusão, alteração, exclusão, deleção, visualização, impressão, extração), data e hora, usuário e IP em todo evento; login, logout, cadastro, alteração, exclusão, deleção e visualização de dado sensível sempre gravados; visualização, impressão e extração ligadas ou desligadas por tela pelo `administrador` (`captura_auditoria`); hash por evento e selo encadeado por empresa e período fora do banco (`selo_auditoria`); telas de consulta, verificação de integridade e configuração |

### Técnicas

- **Infraestrutura no Docker:** PostgreSQL 17, Valkey 8 e, novo, o Mailpit (SMTP de teste) no compose local de [infra/docker/](../../infra/docker/).
- **Partições de `evento_auditoria` sem `pg_partman`**, por função própria, resolvendo para a auditoria a pendência da etapa 011.
- **ADR-011** ([5_ARQUITETURA_SOFTWARE.md](../5_ARQUITETURA_SOFTWARE.md), seção 15): papel `tenaz_autenticacao` com pool exclusivo do módulo de autenticação, e dono próprio `tenaz_credencial` para as tabelas de credencial. Motivo: a revisão de isolamento demonstrou que, sem RLS em `plataforma`, o papel da aplicação lia e forjava credenciais de outra empresa (ISO-E01-02) e que o papel de manutenção lia todas as senhas (ISO-E01-08).
- **`app.escopo_id()` não é criada** (pendência E01-3 do documento 2): as políticas usam funções tipadas por fornecedor, locatário e link.
- **Fim de linha LF** fixado por `.gitattributes` e **número de tamanho fixo** no contrato da API: o Go desta máquina é de 32 bits e o contrato gerado divergia do CI.

## Mudanças até agora

| Documento ou arquivo | O que mudou | IDs afetados |
|---|---|---|
| [0_SUMARIO_EXECUTIVO.md](../0_SUMARIO_EXECUTIVO.md) | Versão 0.6: nomes da autenticação, do calendário e de `usuario_permissao`; decisões da seção 3.6; pendência P11 | P11 |
| [1_ESPECIFICACAO_REQUISITOS.md](../1_ESPECIFICACAO_REQUISITOS.md) | Versão 0.5: catálogo de 117 permissões e matriz dos 12 perfis; RF de ADM alinhados às decisões | RF-ADM-001, 002, 004, 005, 006, 011, 015, RF-MOB-001, RF-SRV-001, RF-SRV-003 |
| [2_ARQUITETURA_BANCO_DADOS.md](../2_ARQUITETURA_BANCO_DADOS.md) | Versão 0.5: DDL de autenticação, calendário, `usuario_permissao`, partições, papéis `tenaz_autenticacao` e `tenaz_credencial`, pendências E01-1 a E01-14 e achados ISO-E01-01 a 08 | — |
| [4_FLUXOS_PROCESSO.md](../4_FLUXOS_PROCESSO.md) | Versão 0.5: regras RN-ADM-017 a RN-ADM-031 (031 revogada), lacunas L-40 a L-43 | RN-ADM-017 a RN-ADM-031 |
| [5_ARQUITETURA_SOFTWARE.md](../5_ARQUITETURA_SOFTWARE.md) | Versão 0.6: ADR-011, seções 6.6 e 6.9, verificações da suíte. Emenda do ADR-011 na mesma etapa: dono `tenaz_credencial` das tabelas de credencial (ISO-E01-08), verificação "Dono das credenciais" na 9.3, `mfa_ultimo_passo` na tabela de privilégios da 6.9 (estava faltando). Conferência do documento contra as migrações 00004 e 00005 corrigiu quatro afirmações que o banco não sustentava: privilégios de esquema do `tenaz_credencial`, membros do papel em desenvolvimento e em produção, objetos de que ele é dono (também índices e tipos de linha) e o `DELETE` do dono usado pela purga | ADR-011 |
| [3_ESPECIFICACAO_APIS.md](../3_ESPECIFICACAO_APIS.md) | Versão 0.5: 19 rotas novas de E01, códigos de erro novos (`senha_fraca`, `mfa_codigo_invalido`, `link_email_invalido`, `ultimo_administrador`, `pre_condicao_obrigatoria`), membros do problema RFC 9457, cookie e CSRF, fluxo do `authorize` do app, grupo de permissão Conexão com `srv.gerenciar_conexao`, esquema `Versao`, pendências A-37 e A-38 | RF-ADM-001 a 006, 011, 015, RF-MOB-001 |
| [9_SEGURANCA_CONFIABILIDADE.md](../9_SEGURANCA_CONFIABILIDADE.md) | Versão 0.4: sessão web, renovação de 30 dias, permissões pela matriz canônica, exceção anti-escalada do administrador | RNF-017 |
| [execucao/1_MODELO_EXECUCAO.md](../execucao/1_MODELO_EXECUCAO.md), [execucao/2_EQUIPE_AGENTES.md](../execucao/2_EQUIPE_AGENTES.md) | Contagem de regras de negócio: 376 | — |
| [CLAUDE.md](../../CLAUDE.md) | Regra da documentação idempotente ao código (seção 3.2); `docs/entregas/` na tabela de fontes e na estrutura do repositório | — |
| [docs/entregas/](../entregas/README.md) | Pasta nova: regras, modelo e índice; entregas [C0](../entregas/C0_concepcao.md) e [E00](../entregas/E00_fundacao-tecnica.md) | — |
| [docs/README.md](../README.md), [execucao/1_MODELO_EXECUCAO.md](../execucao/1_MODELO_EXECUCAO.md) (versão 0.2), [execucao/3_SKILLS.md](../execucao/3_SKILLS.md) (versão 0.2), `.claude/skills/tenaz-fechar-etapa/SKILL.md` | A entrega consolidada passa a fazer parte do aceite e do fechamento de etapa | — |
| [CHANGELOG.md](../../CHANGELOG.md) | Entradas de E01 em Adicionado, Segurança e Corrigido | RNF-040 |
| `backend/db/migracoes/00003_auditoria.sql`, `00004_adm.sql`, `00005_autenticacao.sql` | Auditoria, ADM e autenticação | RF-ADM-001 a 006, 011, 015 |
| `backend/internal/modulos/adm/`, `backend/internal/plataforma/{autenticacao,autorizacao,auditoria,paginacao,resiliencia,email}/` | Casos de uso de ADM e autenticação, catálogo de permissões em Go | RF-ADM-001 a 006, 011, 015, RF-MOB-001 |
| `backend/testes/vazamento/` | Suíte de isolamento cobrindo as tabelas de E01, com provas de detector e o mecanismo de achados abertos | RNF-013, RNF-037 |
| `backend/internal/plataforma/api/`, `backend/internal/modulos/adm/rotas/`, `backend/internal/servidor/`, `backend/testes/integracao/`, `backend/api/openapi.yaml`, `pacotes/api-tipos/` | Rotas HTTP de ADM e de autenticação; códigos de erro do Go renomeados para os do documento 3; cabeçalhos `Tenaz-Versao` e `X-Request-Id` em toda resposta; parâmetro desconhecido → 400; idempotência em memória; testes por rota e de fluxo completo | RF-ADM-001 a 006, 011, 015, RF-MOB-001 |
| `web/`, `pacotes/{api-tipos,tokens-ui,i18n}/` | Fundação da web, selo de versão e "Sobre" | RNF-039, RNF-025 |
| `infra/docker/compose.yaml`, `.gitattributes`, `backend/internal/plataforma/erros/erros.go` | Mailpit, fim de linha LF, `status` em 64 bits | — |

## Verificação até agora

| Verificação | Resultado |
|---|---|
| Migrações 00003 a 00005 | Sobem do zero, descem até a 00002 e sobem de novo |
| Suíte de isolamento | Verde com 240 testes antes do ISO-E01-08; em ajuste pelo `tenaz-isolamento` depois da correção do dono das credenciais |
| Revisão da suíte pelo `tenaz-qa-automacao` | **Reprovada** na primeira passada: 8 de 15 mutações passaram sem detecção. Correções em andamento; nova revisão com as mesmas mutações antes do fechamento |
| Revisão do domínio pelo `tenaz-dominio-manutencao` | Sem bloqueio; quatro correções e regras RN-ADM escritas no documento 4 |
| Testes do backend | 242 verdes na entrega do `tenaz-backend-go`; vermelho no meio da rodada de rotas e de migração, a corrigir |
| `python ferramentas/verificar_docs.py` | 0 erros |

## Em andamento

- `tenaz-qa-automacao`: nova revisão da suíte de isolamento, com as mesmas mutações.
- `tenaz-backend-go`: rodada de correções (MFA obrigatório, reuso de TOTP, falha de MFA no bloqueio, responsável com `gestor`, turnos na criação sem padrão, sobreposição de turnos, permissão avulsa de externo restrita ao perfil, pessoa sem `desligada`, testes de autenticação pedidos pelo QA).
- `tenaz-produto`: ampliação de RF-ADM-011 com a auditoria configurável; depois banco, API, web e documentos 2, 3, 4, 7 e 9.
- `tenaz-arquiteto`: documento 5 alinhado ao código (exemplo de erro, camadas HTTP, A-38).
- Depois: parte 3 da suíte de isolamento (por rota), telas da web e do app, evidência visual, revisão final do QA e aceite.

## Pendências

- **Decisões de 2026-11-12 ainda não aplicadas nos documentos e no código** (registradas no prompt de E01, "Ajustes durante a etapa"): (1) a Tenaz não guarda dado de saúde — tirar ASO, exame toxicológico, `pessoa.saude.ler` e as regras de saúde do sumário, dos documentos 1, 2, 3, 4 e 9, do catálogo de permissões, das migrações e do Go; (2) nenhuma ferramenta ou recurso que não seja 100% gratuito ou open source — auditar a stack e os serviços, substituir o MinIO e rever a nuvem paga dos documentos; (3) detalhes da trilha de auditoria configurável, a aplicar nos documentos 2, 3 e 9 e no código.
- **Revisão da suíte de isolamento (2ª passada, `tenaz-qa-automacao`):** reprovada só pelo item B1 — função `SECURITY DEFINER` fora de `app` e `plataforma` não é vista pela suíte (vazamento demonstrado com as mutações m25 e m26). Dono: `tenaz-isolamento`.
- **Estado do código neste commit:** a rodada de correções do `tenaz-backend-go` estava em andamento quando o dono do produto pediu o commit; o código pode estar no meio dessa rodada.

- **P11:** provedor de SMS para o login por telefone ([sumário](../0_SUMARIO_EXECUTIVO.md), seção 6).
- **E01-11:** CPF e data de nascimento aparecem na trilha de auditoria; resolver com o `tenaz-conformidade-br` antes das tabelas de saúde ([documento 2](../2_ARQUITETURA_BANCO_DADOS.md), seção 17.3).
- **E01-13:** troca de CNPJ de uma conta com conexão, na etapa da conexão entre contas (E07).
- **ISO-E01-09** (risco residual, apontado na emenda do ADR-011): o esquema `plataforma` tem dono `tenaz_dono`; quem roubar a credencial do `tenaz_manutencao`, membro do dono, não lê as credenciais, mas pode apagar e recriar as tabelas de credencial para capturar o que for gravado depois, desligar o trigger de revogação ou desligar o RLS das tabelas de `app`. Opções em aberto no documento 5, seção 16: tirar o `tenaz_manutencao` de `tenaz_dono` (funções estreitas e `MAINTAIN`) ou deixá-lo só no pipeline. Prazo: antes da primeira implantação em homologação.
- **Migrador superusuário em desenvolvimento e CI** ([documento 2](../2_ARQUITETURA_BANCO_DADOS.md), seções 2.2 e 14.4): localmente e no CI as migrações rodam como superusuário (`POSTGRES_USER`), enquanto em homologação e produção o `tenaz_migrador` é só membro dos donos. Um erro de privilégio que só aparece sem superusuário escaparia até a primeira implantação. Pendência para o `tenaz-plataforma-sre`: o CI aplicar as migrações com um migrador não superusuário, igual ao de produção.
- **ISO-E01-10:** o `tenaz_sinc` não tem USAGE em `app`; falha fechada; corrigir na migração de E04.
- **A-37:** idempotência guardada em memória, correta só com uma instância; tabela compartilhada de chaves pendente ([documento 3](../3_ESPECIFICACAO_APIS.md), seção 26).
- **E01-14:** `conexao_empresa` (E04) precisa do mesmo tratamento de dono próprio das credenciais.
- **L-40:** se o calendário da unidade entra no tempo requerido dos indicadores ([documento 4](../4_FLUXOS_PROCESSO.md), seção 35.3).
- Go de 32 bits nesta máquina: recomendada a troca pelo de 64 bits, igual ao CI.
