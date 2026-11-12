---
name: tenaz-implementar-requisito
description: Use para levar um requisito funcional ou não funcional da Tenaz do texto até o aceite, cumprindo os nove critérios de pronto. Dispara sempre que a tarefa for "implementar RF-..." ou "entregar RNF-...", e também quando um requisito já implementado precisar de correção.
---

# Implementar um requisito da Tenaz

Procedimento completo, na ordem. O atalho mais comum deste projeto é implementar o caminho feliz e esquecer validação, auditoria, isolamento, offline, acessibilidade e documentação. Os nove critérios existem para fechar essa porta.

## 1. Entender

1. Leia o requisito **no documento dono**, [docs/1_ESPECIFICACAO_REQUISITOS.md](../../../docs/1_ESPECIFICACAO_REQUISITOS.md), pelo ID. Nunca por paráfrase de terceiro.
2. Leia o **fluxo** correspondente em [docs/4_FLUXOS_PROCESSO.md](../../../docs/4_FLUXOS_PROCESSO.md) e todas as regras de negócio do módulo que o governam.
3. Confirme a **fase**: o requisito é do MVP, da fase 2 ou da fase 3? Trabalhar requisito de fase futura é mudança de escopo.
4. Confira o **vocabulário canônico** no [sumário](../../../docs/0_SUMARIO_EXECUTIVO.md), seção 3: entidade, campos, estados, permissão, rota, evento.
5. Liste os **casos de borda** antes de escrever qualquer código. Se você não consegue listá-los, não entendeu o requisito: chame `tenaz-dominio-manutencao`.

## 2. Modelar

1. Veja se o modelo de dados já está especificado em [docs/2_ARQUITETURA_BANCO_DADOS.md](../../../docs/2_ARQUITETURA_BANCO_DADOS.md). Quase sempre está.
2. Tabela, coluna ou política nova: use [`tenaz-tabela-nova`](../tenaz-tabela-nova/SKILL.md). Sem exceção.
3. Campo novo de OS ou de ativo: classifique como compartilhável ou nunca compartilhado, com [`tenaz-campo-compartilhado`](../tenaz-campo-compartilhado/SKILL.md). Classificação **antes** de o campo existir.

## 3. Contratar

1. Veja se a rota já está em [docs/3_ESPECIFICACAO_APIS.md](../../../docs/3_ESPECIFICACAO_APIS.md).
2. Rota, erro ou evento novo: use [`tenaz-endpoint-novo`](../tenaz-endpoint-novo/SKILL.md).
3. Declare o escopo de acesso da rota: `interno`, `prestador`, `locatario`, `link`.

## 4. Implementar

1. Domínio e caso de uso primeiro, com erro de domínio tipado e código estável.
2. Transação com contexto de empresa; repositório filtrando por empresa **também** de forma explícita, mesmo com RLS ativo.
3. Auditoria onde o requisito pedir (RF-ADM-011): autor, instante, valores antes e depois.
4. Interface só depois do contrato gerado, com [`tenaz-tela-nova`](../tenaz-tela-nova/SKILL.md).
5. Operação de campo: [`tenaz-operacao-offline`](../tenaz-operacao-offline/SKILL.md).
6. Indicador: [`tenaz-indicador`](../tenaz-indicador/SKILL.md).

## 5. Verificar — os nove critérios de pronto

Da seção 3.1 do [roadmap](../../../docs/6_ROADMAP.md). Marque um por um, com a prova ao lado.

| # | Critério | Prova |
|---|---|---|
| 1 | Validações e regras com teste automatizado | Saída da suíte |
| 2 | Tabela nova com RLS forçado e no teste de vazamento; tabela de escopo externo nos testes de RNF-037 ou RNF-038; campo tocado pela sincronização no teste de três contas de RNF-041 | Saída das suítes de isolamento |
| 3 | Endpoints na especificação OpenAPI, e são os mesmos que web e app usam | Contrato gerado igual ao versionado |
| 4 | Alterações relevantes gerando evento de auditoria | Consulta à trilha depois da operação |
| 5 | Uso offline previsto: teste com desconexão, ação, reconexão e conflito | Saída do teste de sincronização |
| 6 | Telas novas com acessibilidade verificada e sem esconder o selo de versão | Relatório de acessibilidade e captura |
| 7 | p95 dentro de RNF-007 em homologação com massa de 2.000 ativos | Saída do teste de carga |
| 8 | Ajuda contextual mínima na tela | Captura |
| 9 | Aceite em homologação; requisitos de CON também pelo consultor de SST; requisitos com NF-e também pelo contador | Registro do aceite |

**Critério faltando significa requisito não entregue.** Não existe "pronto com pendência".

## 6. Documentar

Use [`tenaz-fechar-etapa`](../tenaz-fechar-etapa/SKILL.md): documento afetado com versão subida, `CHANGELOG.md`, registro da etapa, sumário quando couber, verificador sem erro, PDFs regerados.

## Erros que esta skill previne

- Implementar requisito de outra fase sem perceber.
- Criar tabela sem política de escopo externo, que só aparece quando o prestador entra, etapas depois.
- Criar campo e descobrir na suíte de três contas que ele vaza entre contas.
- Entregar tela sem estado vazio, sem erro e sem ajuda contextual.
- Fechar sem entrada no `CHANGELOG.md`, o que viola RNF-040.
- Declarar pronto sem rodar o teste de carga, e descobrir o p95 no endurecimento.
