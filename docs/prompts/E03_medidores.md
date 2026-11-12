# E03 — Medidores e composição

**Data do pedido:** 2026-11-12 · **Sprint do roadmap:** S3 · **Marco:** M1
**Pedido por:** dono do produto

## Objetivo

Todo ativo passa a ter medidores confiáveis: km, horas, ciclos, com validação de regressão, teto de salto e troca de hodômetro ou horímetro. É a base da preventiva por medidor e dos indicadores de custo por km, e o erro aqui se propaga para os dois.

## Entrega

| Requisito | Observação de escopo nesta etapa |
|---|---|
| RF-MED-001 a RF-MED-008 | completos |
| RF-ATV-004 | Composição de veículos |
| RF-ATV-009, RF-ATV-011 | completos |
| RF-ATV-010 | Parcial. Completo em E09 |
| RF-MOB-003 | Parcial: leitura de medidor offline |
| RF-SRV-005 | Equipamentos de clientes, com propriedade `de_cliente` |
| RF-SRV-003 | Parcial: técnicos com especialidades e capacitações. Documentos dos técnicos em E09 |

## Contexto obrigatório

| Documento | Seção | Por quê |
|---|---|---|
| [0_SUMARIO_EXECUTIVO.md](../0_SUMARIO_EXECUTIVO.md) | 3.1, 3.5 | Medidor, leitura e origens; equipamento de cliente na conta da subcontratada |
| [2_ARQUITETURA_BANCO_DADOS.md](../2_ARQUITETURA_BANCO_DADOS.md) | 7 | Medidores, leituras e segmento de medidor |
| [4_FLUXOS_PROCESSO.md](../4_FLUXOS_PROCESSO.md) | Fluxo 04 | Leitura de medidor, com validações e exceções |
| [3_ESPECIFICACAO_APIS.md](../3_ESPECIFICACAO_APIS.md) | 9 | Registrar leitura |

## Agentes previstos

Conduz `tenaz-dominio-manutencao`. Apoiam `tenaz-banco-dados`, `tenaz-backend-go`, `tenaz-api-contrato`, `tenaz-frontend-web`, `tenaz-mobile-offline`, `tenaz-qa-automacao`, `tenaz-conexao-contas` (equipamentos de cliente).

## Portão de saída

| Verificação | Como se comprova |
|---|---|
| **Marco M1** | Empresa de teste com 2.000 ativos importados em menos de 2 minutos; conta de subcontratada com clientes, equipamentos `de_cliente` e técnicos com especialidades |
| Validação de leitura | Regressão, teto de salto e troca de medidor com deslocamento, cada caso com teste |
| Casos compartilhados | A mesma bateria de casos de validação passando no Go e no TypeScript, para que o app valide offline igual ao servidor |
| Origem registrada | Toda leitura com origem gravada, conforme os valores do sumário |
| Leitura offline | Leitura registrada em modo avião, sincronizada depois, sem duplicar |

## Evidência esperada

`docs/screenshots/E03_medidores/`: cadastro de medidor, registro de leitura na web, registro no app, histórico com a troca de medidor marcada, recusa por regressão, composição de veículo, ficha de equipamento de cliente.

## Fora do escopo

Leitura por telemática, cartão combustível e IoT, que são das fases 2 e 3. Planos preventivos (E06). Indicadores (E09).

## Decisões pendentes do dono do produto

Nenhuma.

## Ajustes durante a etapa

Nenhum até agora.
