# E01 — Conta e acesso

**Data do pedido:** 2026-11-12 · **Sprint do roadmap:** S1 · **Marco:** —
**Pedido por:** dono do produto

## Objetivo

Uma empresa passa a existir no sistema, com unidades, pessoas, usuários, perfis e permissões, trilha de auditoria e a versão visível na tela. A conta de subcontratada nasce aqui também, para que a trilha do módulo SRV possa correr em paralelo desde o início.

## Entrega

| Requisito | Observação de escopo nesta etapa |
|---|---|
| RF-ADM-001 | Sem o roteiro de primeiros passos, que fica em E12 |
| RF-ADM-002 | completo |
| RF-ADM-004 | completo |
| RF-ADM-005 | Com a permissão `os.encerrar`, a permissão `est.decidir_destino` e os escopos `interno`, `prestador`, `locatario`, `link` |
| RF-ADM-006 | completo |
| RF-ADM-011 | completo: trilha de auditoria com valores antes e depois |
| RF-ADM-015 | completo: entidade de pessoa, com e sem usuário |
| RF-MOB-001 | completo, com o selo de versão e a tela "Sobre" no app |
| RNF-039 | Selo de versão e "Sobre" na tela de login e na tela principal da web |
| RF-SRV-001 | Parcial: empresa com o módulo SRV habilitado, perfis da conta e escolha de empresa no login. Cobrança em E10, primeiros passos em E12 |

## Contexto obrigatório

| Documento | Seção | Por quê |
|---|---|---|
| [0_SUMARIO_EXECUTIVO.md](../0_SUMARIO_EXECUTIVO.md) | 3.1, 3.4, 3.6 | Nomes canônicos de empresa, unidade, pessoa, usuário, perfil; formato das permissões; decisões de nome |
| [1_ESPECIFICACAO_REQUISITOS.md](../1_ESPECIFICACAO_REQUISITOS.md) | 4 (ADM), 2 | Os requisitos e os atores |
| [2_ARQUITETURA_BANCO_DADOS.md](../2_ARQUITETURA_BANCO_DADOS.md) | 2, 5, 9 | Isolamento, DDL do MVP, trilha de auditoria |
| [3_ESPECIFICACAO_APIS.md](../3_ESPECIFICACAO_APIS.md) | 4 | Autenticação e autorização |
| [7_PADROES_DESIGN.md](../7_PADROES_DESIGN.md) | 3, 5, 8 | Tokens, layout da web, telas principais |
| [9_SEGURANCA_CONFIABILIDADE.md](../9_SEGURANCA_CONFIABILIDADE.md) | 5, 6 | Segurança e trilha de auditoria |

## Agentes previstos

Conduz `tenaz-produto` no escopo e `tenaz-backend-go` na implementação. Apoiam `tenaz-banco-dados`, `tenaz-isolamento`, `tenaz-api-contrato`, `tenaz-design-sistema`, `tenaz-frontend-web`, `tenaz-mobile-offline`, `tenaz-qa-automacao`, `tenaz-documentacao`.

## Portão de saída

| Verificação | Como se comprova |
|---|---|
| Permissões e escopos conforme o sumário | Lista de permissões no formato `<módulo>.<ação>`; `os.encerrar` presente em gestor e planejador e ausente em prestador e locatário |
| Pessoa separada de usuário | Pessoa cadastrada sem usuário e pessoa com usuário, nas duas direções |
| Auditoria completa | Cada alteração relevante com autor, instante e valores antes e depois |
| Isolamento | Suíte de vazamento cobrindo todas as tabelas criadas nesta etapa, inclusive as de pessoa e de perfil |
| Versão visível | Selo no canto inferior direito do login e da tela principal, na web e no app; "Sobre" com versão, versão da API, data do build, commit curto e ambiente |
| Acessibilidade | Verificação automática e manual de WCAG 2.1 AA nas telas desta etapa |

## Evidência esperada

`docs/screenshots/E01_conta-e-acesso/`: login, tela principal com o selo de versão, "Sobre", cadastro de empresa e unidade, cadastro de pessoa, cadastro de usuário com perfis, trilha de auditoria, seleção de empresa no login, nos dois temas.

## Fora do escopo

SSO e SAML (E10 e fase 2). Multi-unidade com usuário restrito a uma unidade, que é RF-ADM-003, da fase 2. Cobrança. Cadastro de ativos (E02).

## Decisões pendentes do dono do produto

Nenhuma trava esta etapa. A tabela de planos depende de **P4**, mas a cobrança só entra em E10.

## Ajustes durante a etapa

**2026-11-12 — Decisões do dono do produto no início da etapa.**

| Assunto | Decisão |
|---|---|
| Tabelas de autenticação, ausentes do documento 2 | Aprovados os nomes `sessao` (campo `canal`), `token_renovacao`, `codigo_autorizacao`, `desafio_mfa` e `token_email` (campo `finalidade`), no esquema da plataforma, e os campos `identidade.falhas_login`, `identidade.bloqueada_ate` e `identidade.mfa_ativado_em`. Entraram no sumário 0.6, seção 3.2 |
| Calendário de trabalho da unidade (RF-ADM-002), ausente do DDL de `unidade` | Tabelas próprias: `turno_unidade` e `feriado_unidade` |
| Login por telefone e código (RF-ADM-006, proposta) | Adiado até a escolha do provedor de SMS: pendência **P11** do sumário. E01 entrega e-mail e senha com MFA por TOTP |

| Catálogo de permissões (a matriz não existia em nenhum documento) | O catálogo e a matriz escritos em RF-ADM-005 (documento 1, versão 0.5) passam a ser canônicos por uma linha no sumário 3.6, sem repetir a lista |
| Permissão avulsa, como `pessoa.saude.ler` (sem lugar no banco) | Nova entidade `usuario_permissao` |
| Permissões do `administrador` | Todas, menos `pessoa.saude.ler`, `os.encerrar`, `est.decidir_destino` e `ter.manter_catalogo_proprio` |
| `usuario.situacao` (divergência entre os documentos 1 e 2) | `convidado`, `ativo`, `inativo` |
| MFA obrigatório por empresa (RF-ADM-006), não implementado na primeira passada | Implementar em E01, com o campo `empresa.mfa_obrigatorio`. Retificação do dono do produto: "o MFA obrigatório no login também deverá ser requisito configurável" — confirmado como um único interruptor por empresa, ligado ou desligado pelo `administrador`; sem exigência por perfil nem por usuário |
| Vínculo e situação da pessoa (lacuna L-42, valores do DDL sem nome no sumário) | Vínculo `empregado`, `terceiro` ou `autonomo`. Situação só `ativa` ou `inativa`, sem `desligada`: "quando inativo permanece no banco de dados [...] se houver vínculo com algum tipo de serviço que foi feito. Ex: transporte de equipamentos, pode estar vinculado a um motorista temporário" |
| Trilha de auditoria configurável (pedido do dono do produto: "um sistema de captura de logs e auditoria bem robusto que possibilitasse a configuração dos tipos de captura nas telas do sistema (login, logout, visualização, cadastro, alteração, deleção, exclusão, impressão e extração) sempre com data/hora, nome do usuário, IP da máquina, tipo da operação, nome da tela ou módulo, e controle da integridade por hash") | **Entra em E01**, ampliando RF-ADM-011. Login, logout, cadastro, alteração, exclusão, deleção e visualização de dado sensível são sempre gravados; o administrador liga ou desliga por tela visualização, impressão e extração. Integridade por hash de cada evento e selo encadeado por empresa e período, guardado fora do banco. Telas de consulta (com verificação de integridade do período) e de configuração entram na etapa |
| Detalhes da trilha (perguntas do `tenaz-produto` sobre RF-ADM-011) | Falha de login, acesso do suporte, revisão de leitura e sincronização viram atributos do evento, sem tipo próprio. Impressão e extração de dado sensível e a exportação completa da empresa são sempre gravadas. Aceitas: capturas padrão (impressão e extração ligadas, visualização desligada), visualização = registro aberto em detalhe, selo diário por dia UTC, armazenamento imutável com bloqueio de objeto, estado "ainda não selado", configuração com `adm.gerenciar_parametros` |
| Ferramentas e recursos | **Requisito do dono do produto: não usar nenhuma ferramenta ou recurso que não seja 100% gratuito ou open source.** O MinIO sai (deixou de ser open source na prática); buscar substituto compatível com S3 e com bloqueio de objeto. Toda a stack e os serviços citados nos documentos (inclusive nuvem paga, como AWS, KMS e GuardDuty) precisam ser auditados contra essa regra |
| Dados de saúde | **A Tenaz não guarda dado de saúde.** ASO e exame toxicológico saem do escopo, com a permissão `pessoa.saude.ler` e as regras de dado sensível de saúde. Aplicar no sumário, nos documentos 1, 2, 3, 4, 9, no catálogo de permissões e no código |
| Alcance da regra de ferramentas (auditoria do `tenaz-arquiteto`, ADR-012) | "Ambiente de desenvolvimento é diferente de ambiente de produção. A questão de gratuidade se aplica exclusivamente para ferramentas utilizadas diretamente pelo sistema. Documente isso como regra." O Docker Desktop continua no desenvolvimento (o dono do produto está dentro da gratuidade) |
| Hospedagem de produção (H1) | Nuvem gerenciada (AWS) como exceção explícita à regra; ADR-006 volta a valer |
| Selo da auditoria em desenvolvimento (H10) | Por enquanto, tudo numa única máquina, no Docker |
| Reforço da vedação de dado de saúde | Aviso fixo no cadastro da pessoa, no dossiê do motorista e na tela de anexo, e recusa de tipo de documento com nome ligado a saúde |
| Detalhes do reforço contra dado de saúde (perguntas do `tenaz-produto`) | Texto do aviso aceito: "A Tenaz não guarda dado de saúde. Não registre nem anexe ASO, exames, atestados ou qualquer informação médica." Saem `clinica` e `exame` da lista de termos recusados. O nome do arquivo anexado não é verificado, só o aviso. O aviso no envio de anexo entra em E02; o da pessoa, em E01. Tipos de documento antigos com nome vedado ficam só listados para revisão, sem inativação automática |
| Serviços pagos por natureza (H5 a H9: lojas, e-mail, SMS, IA, pagamento) | Decididos na etapa em que cada um entra |
| Troca de CNPJ da empresa (proposta RN-ADM-031: só pelo suporte) | Recusada: o `administrador` altera o CNPJ. A conexão entre contas trata a troca na sua etapa |
| Perfis do responsável pela conta (lacuna L-43) | `administrador` e `gestor` |
| Calendário padrão da unidade (lacuna L-40) | "A configuração de turnos deve ser item configurável": o responsável informa os turnos na criação; o sistema não impõe padrão, só oferece modelos para preencher |

Decisões técnicas tomadas pelo orquestrador, sem mudança de escopo:

- Validade do refresh token do app: **30 dias** de inatividade, como em RNF-017 (documento 1, dono do requisito). O documento 9 dizia 14 dias e é corrigido.
- Envio de e-mail em desenvolvimento por um servidor SMTP de teste no compose local (Mailpit), sem provedor real.
- Partições mensais de `evento_auditoria` criadas por função própria, sem `pg_partman`, que não acompanha a imagem oficial do PostgreSQL (pendência da etapa 011).
