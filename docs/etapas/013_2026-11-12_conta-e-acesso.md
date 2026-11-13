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
| Ferramentas e recursos | Só software 100% gratuito ou open source no que o sistema usa diretamente; o ambiente de desenvolvimento fica fora da regra (Docker Desktop continua); AWS gerenciada em produção como exceção explícita (ADR-006 volta a valer); MinIO sai e o SeaweedFS entra no desenvolvimento (ADR-012); serviços pagos por natureza decididos na etapa em que entram |
| Dados de saúde | A Tenaz não guarda dado de saúde; o produto reforça com aviso fixo nas telas de pessoa, dossiê do motorista e anexo, e com recusa de tipo de documento com nome ligado a saúde |
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
| `infra/docker/compose.yaml`, `infra/docker/seaweedfs/security.toml`, `infra/docker/seaweedfs/s3.json`, `ferramentas/tarefas.sh`, [9_SEGURANCA_CONFIABILIDADE.md](../9_SEGURANCA_CONFIABILIDADE.md) (versão 0.5, complemento) | MinIO substituído pelo SeaweedFS 4.53 por digest, no `up` comum, com JWT de volume e filer, só a porta S3 publicada, identidades S3 por uso e telemetria desligada; imagens do compose por versão exata e digest (PostgreSQL 17.11, Valkey 8.1.11, Mailpit, aws-cli 2.40.2); `infra` espera todos saudáveis e cria o bucket `tenaz-auditoria` com bloqueio de objeto COMPLIANCE e versionamento; tarefas novas `armazenamento-preparar` e `armazenamento-provar` (19 de 19 em 2026-11-13) | RF-ADM-011, ADR-012 |

## Verificação até agora

| Verificação | Resultado |
|---|---|
| Migrações 00003 a 00005 | Sobem do zero, descem até a 00002 e sobem de novo |
| Suíte de isolamento | Verde com 240 testes antes do ISO-E01-08; em ajuste pelo `tenaz-isolamento` depois da correção do dono das credenciais |
| Revisão da suíte pelo `tenaz-qa-automacao` | 1ª passada **reprovada** (8 de 15 mutações sem detecção); 2ª **reprovada** só pelo B1 (função `SECURITY DEFINER` fora de `app` e `plataforma`); 3ª, em 2026-11-13, **aprovada para o banco**: das 20 mutações, só passou a m37 (função anexada a extensão de propósito, P-QA-1), e a regra das consultas sqlc aceita truques como `OR true` e filtro em comentário (P-QA-2). As duas e os resíduos (P-QA-3) estão sendo fechados pelo `tenaz-isolamento`. Parte por rota ainda não revisada |
| Revisão dos testes do backend pelo `tenaz-qa-automacao` | **Aprovada** em 2026-11-13: das 18 mutações no código (CSRF, `nbf`, zerar falhas, reuso de TOTP, trava do último administrador, 404 em outra empresa, `OR true` em consulta de credencial, `NaEmpresa`), só uma passou verde (reuso de refresh sem revogar a família, coberto pela revogação da sessão). Dois "deve corrigir" em correção pelo `tenaz-backend-go`: o analisador das consultas de credencial aceitava junção, aspas e `LATERAL`; dois testes comparavam o relógio do banco com o do host |
| Revisão do domínio pelo `tenaz-dominio-manutencao` | Sem bloqueio; quatro correções e regras RN-ADM escritas no documento 4 |
| Testes do backend | 242 verdes na entrega do `tenaz-backend-go`; vermelho no meio da rodada de rotas e de migração, a corrigir |
| `python ferramentas/verificar_docs.py` | 0 erros |

## Em andamento

Atualizado em 2026-11-13.

**Concluído desde o commit 310dee9** (pasta de trabalho, com `conferir` verde e 363 testes de Go verdes):

- Rodada de correções do `tenaz-backend-go`: MFA obrigatório por empresa, com sessão restrita a ativar o MFA (`mfa_pendente`, RN-ADM-027); código TOTP de uso único (`identidade.mfa_ultimo_passo`); código de MFA errado soma falha no bloqueio progressivo; responsável nasce `administrador` e `gestor` (RN-ADM-017); turnos da unidade inicial informados, sem padrão, e recusa de turnos sobrepostos (RN-ADM-019); permissão avulsa de usuário externo só dentro do perfil de origem (RN-ADM-026); pessoa sem `desligada`, pessoa com usuário ativo ou convidado não é inativada, usuário não é ligado a pessoa inativa, pessoa de prestador herda o fornecedor e o vínculo `terceiro` (RN-ADM-020, RN-ADM-021); catálogo de permissões com `adm.gerenciar_parametros` cobrindo a configuração das capturas e o grupo Conexão em `srv.gerenciar_conexao` (RF-ADM-005).
- Documento 1 com RF-ADM-011 ampliado (auditoria configurável, nove critérios de aceite) e documento 5 alinhado à suíte e ao código (seções 2.5, 5.2, 7, 8.2, 9.3, 9.5, 16).

**Concluído em 2026-11-13:** testes do backend pedidos pelo `tenaz-qa-automacao` (CSRF negativo, convite e link expirados, fronteira do bloqueio, mesma resposta para senha errada, bloqueio e e-mail inexistente, sessão web e refresh vencidos, TOTP reapresentado, último administrador em concorrência, 28 casos de outra empresa sem oráculo, situação do usuário, auditoria, `NaEmpresa`, checagem de consultas endurecida); um defeito real achado e corrigido (access token com `nbf` no futuro era aceito); saúde retirada do Go (catálogo com 116 permissões, `NaEmpresa` com 14 variáveis, contrato sem saúde); subida dos contêineres de teste serializada (falha intermitente no Windows).

**Verificação de 2026-11-13, com tudo isso na pasta de trabalho:** `ferramentas/tarefas.sh conferir` verde (documentação com 0 erros e 0 avisos, contrato e tipos em dia); `go test ./...` com **465 testes verdes** em 26 pacotes; suíte de isolamento verde (297 casos). Pronto para o commit que corrige o CI do commit 310dee9, quando o dono do produto pedir.

**Trabalhando agora:**

- `tenaz-api-contrato`: documento 3 com `mfa_obrigatorio`, `mfa_pendente` e a situação da pessoa, como já estão no contrato gerado.
- `tenaz-backend-go`: testes pedidos pelo `tenaz-qa-automacao` (CSRF negativo, convite e link expirados, sessão e refresh vencidos, concorrência do último administrador, auditoria, `NaEmpresa`, checagem de consultas endurecida).
- `tenaz-isolamento`: item B1 da 2ª revisão do QA (função `SECURITY DEFINER` fora de `app` e `plataforma`).

- **Dados de saúde fora do escopo** (decisão de 2026-11-13): feito no sumário 0.7 e no documento 1 versão 0.6 (catálogo com 116 permissões, sem exceção ao formato; nenhum RF ou RNF inteiro sem objeto) e no CHANGELOG (grupo Removido). Feito também: documento 2 versão 0.6 e migração nova `00006_remove_saude.sql` (as migrações 00001 a 00005 já estão em `main` e não se editam; Down e Up testados), documento 3 versão 0.6, documento 4 versão 0.6 (RN-ADM-008 removida), documento 6 versão 0.5, documento 9 versão 0.5 (seção 7.6 removida; o cliente controla ASO, NR-7, fora da Tenaz), suíte de isolamento sem uso da permissão, e a ressalva no sumário de que o "dado sensível" da Tenaz não é dado pessoal sensível da LGPD. Em andamento: documento 5 e o Go (autorização, `NaEmpresa`, ADM, contrato), que precisam entrar no mesmo commit da 00006. O teste que compara o catálogo do Go com o documento 1 falha até o Go acompanhar.
- Auditoria da stack (`tenaz-arquiteto`): feita; ADR-012 no documento 5 versão 0.7; SeaweedFS recomendado depois de ensaio prático de bloqueio de objeto (15 de 15 verificações; o ensaio mostrou que, sem `security.toml`, a API interna apagava bucket em retenção sem credencial, o que tornou o endurecimento obrigatório). Feito em 2026-11-13: emenda do ADR-012 com o alcance da regra e a exceção da AWS (ADR-006 volta a valer); SeaweedFS 4.53 no compose de desenvolvimento, endurecido, com bucket `tenaz-auditoria` em bloqueio de objeto e 19 de 19 provas (`armazenamento-provar`); todas as imagens por versão exata e digest; reforço da vedação de dado de saúde nos requisitos (aviso fixo e recusa de tipo de documento por nome). Pendentes: tarefa `licencas` e job de CI do armazenamento (`tenaz-plataforma-sre`).

**Depois:** decisões do dono do produto sobre os serviços pagos por natureza (hospedagem, lojas, SMS, e-mail), auditoria configurável (banco, API, telas), parte 3 da suíte de isolamento (por rota), telas da web e do app, evidência visual, revisão final do QA e aceite.

## Pendências

- **Decisões de 2026-11-12 ainda não aplicadas nos documentos e no código** (registradas no prompt de E01, "Ajustes durante a etapa"): (1) a Tenaz não guarda dado de saúde — tirar ASO, exame toxicológico, `pessoa.saude.ler` e as regras de saúde do sumário, dos documentos 1, 2, 3, 4 e 9, do catálogo de permissões, das migrações e do Go; (2) nenhuma ferramenta ou recurso que não seja 100% gratuito ou open source — auditar a stack e os serviços, substituir o MinIO e rever a nuvem paga dos documentos; (3) detalhes da trilha de auditoria configurável, a aplicar nos documentos 2, 3 e 9 e no código.
- **Revisão da suíte de isolamento (2ª passada, `tenaz-qa-automacao`):** reprovada só pelo item B1 — função `SECURITY DEFINER` fora de `app` e `plataforma` não é vista pela suíte (vazamento demonstrado com as mutações m25 e m26). Dono: `tenaz-isolamento`.
- **CI do commit 310dee9 vermelho:** o commit, pedido pelo dono do produto, pegou a rodada do `tenaz-backend-go` pela metade, e o contrato versionado ficou atrás do código. Corrigido na pasta de trabalho em 2026-11-13; entra no próximo commit junto com o documento 3.

- **P11:** provedor de SMS para o login por telefone ([sumário](../0_SUMARIO_EXECUTIVO.md), seção 6).
- **Escopo `link`** recebe `permissao_negada` em vez de 404 fora do alcance (`Ator.Externo()` não trata `link` como externo); corrigir na etapa do link de uso único (documento 3, seção 26).
- **Retirada de perfil ou de permissão avulsa** entra na trilha como `alteracao` de `excluido_em`, e RF-ADM-011 pede `exclusao` com a linha inteira em `antes`: parte da implementação da auditoria configurável.
- **E01-11:** CPF e data de nascimento aparecem na trilha de auditoria; resolver com o `tenaz-conformidade-br` antes das tabelas de saúde ([documento 2](../2_ARQUITETURA_BANCO_DADOS.md), seção 17.3).
- **E01-13:** troca de CNPJ de uma conta com conexão, na etapa da conexão entre contas (E07).
- **ISO-E01-09** (risco residual, apontado na emenda do ADR-011): o esquema `plataforma` tem dono `tenaz_dono`; quem roubar a credencial do `tenaz_manutencao`, membro do dono, não lê as credenciais, mas pode apagar e recriar as tabelas de credencial para capturar o que for gravado depois, desligar o trigger de revogação ou desligar o RLS das tabelas de `app`. Opções em aberto no documento 5, seção 16: tirar o `tenaz_manutencao` de `tenaz_dono` (funções estreitas e `MAINTAIN`) ou deixá-lo só no pipeline. Prazo: antes da primeira implantação em homologação.
- **Migrador superusuário em desenvolvimento e CI** ([documento 2](../2_ARQUITETURA_BANCO_DADOS.md), seções 2.2 e 14.4): localmente e no CI as migrações rodam como superusuário (`POSTGRES_USER`), enquanto em homologação e produção o `tenaz_migrador` é só membro dos donos. Um erro de privilégio que só aparece sem superusuário escaparia até a primeira implantação. Pendência para o `tenaz-plataforma-sre`: o CI aplicar as migrações com um migrador não superusuário, igual ao de produção.
- **ISO-E01-10:** o `tenaz_sinc` não tem USAGE em `app`; falha fechada; corrigir na migração de E04.
- **A-37:** idempotência guardada em memória, correta só com uma instância; tabela compartilhada de chaves pendente ([documento 3](../3_ESPECIFICACAO_APIS.md), seção 26).
- **E01-14:** `conexao_empresa` (E04) precisa do mesmo tratamento de dono próprio das credenciais.
- **Como o produto reforça que não recebe dado de saúde** (decisão do dono do produto; opções no [documento 9](../9_SEGURANCA_CONFIABILIDADE.md), seção 12): só termos de uso, aviso fixo no cadastro da pessoa e no dossiê do motorista, recusa de tipo de documento com nome de saúde, ou aviso na tela de anexo.
- **E01-11:** CPF e data de nascimento na trilha de auditoria; proposta de registrar só que o campo mudou, antes do primeiro cliente em produção (E12).
- **L-40:** se o calendário da unidade entra no tempo requerido dos indicadores ([documento 4](../4_FLUXOS_PROCESSO.md), seção 35.3).
- Go de 32 bits nesta máquina: recomendada a troca pelo de 64 bits, igual ao CI.
