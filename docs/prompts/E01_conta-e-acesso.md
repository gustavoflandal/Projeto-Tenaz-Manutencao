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

Nenhum até agora.
