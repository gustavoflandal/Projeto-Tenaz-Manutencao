---
name: tenaz-tabela-nova
description: Use SEMPRE que nascer ou mudar uma tabela, coluna, enum, índice ou política no banco da Tenaz. Garante empresa_id, RLS forçado, política de escopo externo, classificação dos campos da conexão entre contas, entrada na suíte de vazamento e migração no padrão expandir e contrair. Nenhuma tabela entra fora deste procedimento.
---

# Tabela nova no banco da Tenaz

Uma tabela que entra sem este procedimento é um vazamento esperando acontecer. Os riscos R2, R4 e R5 do [sumário](../../../docs/0_SUMARIO_EXECUTIVO.md) são os únicos críticos do projeto, e todos passam por aqui.

## 1. Antes de criar

1. Confira se a tabela já está especificada em [docs/2_ARQUITETURA_BANCO_DADOS.md](../../../docs/2_ARQUITETURA_BANCO_DADOS.md), seções 5 e 13. Grande parte do modelo do MVP e das fases 2 e 3 já existe.
2. Confira o nome canônico da entidade no [sumário](../../../docs/0_SUMARIO_EXECUTIVO.md), seções 3.1 e 3.2. Nome que não está lá entra **primeiro** no sumário.
3. Decida: é tabela de **dados de cliente** ou **catálogo global**?

## 2. Colunas obrigatórias

| Coluna | Regra |
|---|---|
| Chave primária | UUIDv7. Ordenável por tempo, não enumerável, gerável no app offline |
| `empresa_id` | **Obrigatória** em toda tabela de dados de cliente. Sem exceção |
| Datas de criação e alteração | Em UTC |
| Autoria | Quem criou e quem alterou, quando o requisito pedir trilha |

Nomes em `snake_case`, português, sem acento. Dinheiro em inteiro de centavos mais código da moeda, **nunca** ponto flutuante. Marcos de tempo da OS com prefixo `data_`.

## 3. Isolamento — as três camadas

### 3.1 RLS por empresa (obrigatório)

1. RLS **habilitado** e **forçado** na tabela.
2. Política com `USING` **e** `WITH CHECK`. Só `USING` protege leitura e deixa a escrita aberta.
3. Sem o contexto de empresa definido, a tabela devolve **zero linha**.

### 3.2 Escopo externo (obrigatório)

Toda tabela com `empresa_id` recebe a política de escopo externo **restritiva**. Ela é adicional ao RLS, não substituta.

- A tabela é legível por `prestador`, `locatario` ou `link`? Então entra na **lista de acesso externo**, guardada no repositório, com o predicado que limita o que cada escopo vê.
- Não é legível por escopo externo? O predicado recusa, e a rota responde 404 para esses escopos.

A suíte verifica que a lista do repositório **bate exatamente** com as tabelas cujo predicado aceita escopo externo. Divergência quebra o build.

### 3.3 Conexão entre contas (quando aplicável)

A tabela é lida ou gravada pelo processo de sincronização entre contas?

1. Recebe a política de sincronização **restritiva**.
2. O papel do sincronizador recebe privilégio **exatamente** nas colunas da lista fechada, nunca na tabela inteira.
3. Cada coluna é classificada com [`tenaz-campo-compartilhado`](../tenaz-campo-compartilhado/SKILL.md).
4. O papel **não tem** privilégio de contornar RLS e **não é dono** de tabela.

## 4. Catálogo global

Tabela sem `empresa_id` só existe como catálogo global e **precisa entrar na lista explícita de tabelas globais**, revisada em revisão de código. Tabela fora da lista e sem `empresa_id` quebra o build.

Pergunta de controle: duas empresas diferentes podem ver a mesma linha sem problema? Se a resposta tem qualquer "depende", não é catálogo global.

## 5. Suíte de vazamento

Acrescente a tabela nas verificações da seção 9.3 de [docs/5_ARQUITETURA_SOFTWARE.md](../../../docs/5_ARQUITETURA_SOFTWARE.md):

- Leitura e escrita cruzada entre as duas empresas fictícias.
- Comportamento sem contexto de empresa.
- Leitura cruzada por rota, com token de uma empresa e identificador de outra: 404, nunca 403.
- Quando houver escopo externo: entre dois fornecedores, entre dois locatários, e com o mesmo CNPJ em duas empresas.
- Quando houver sincronização: na suíte de três contas, com valor sentinela nos campos internos que **não** podem trafegar.
- Arquivos: anexo de uma empresa sem URL assinada para outra.

Semente de dados povoada **em todas as tabelas**: tabela vazia passa em qualquer teste de vazamento.

## 6. Migração

Padrão **expandir e contrair**:

1. Compatível com a versão anterior do código. As duas versões convivem durante o deploy.
2. Coluna nova aceita nulo ou tem valor padrão; preenchimento é passo separado.
3. **Remoção de coluna só no release seguinte** ao que parou de usá-la.
4. Renomear é adicionar, copiar, parar de usar e remover depois.
5. Índice em tabela grande criado de forma concorrente; migração que trava tabela grande é dividida.
6. Todo índice novo acompanha a consulta que o justifica, e começa por `empresa_id` quando a consulta filtra por empresa.

## 7. Fechar

1. Atualize [docs/2_ARQUITETURA_BANCO_DADOS.md](../../../docs/2_ARQUITETURA_BANCO_DADOS.md): DDL, diagrama se houver, versão do documento e linha no histórico.
2. Peça a revisão de `tenaz-isolamento`. **Quem cria a política não aprova a política.**
3. Rode as suítes de integração e de vazamento e cole a saída.
4. Feche com [`tenaz-fechar-etapa`](../tenaz-fechar-etapa/SKILL.md).

## Verificação final

```bash
python ferramentas/verificar_docs.py
```

Mais a suíte de vazamento verde, com a verificação de cobertura de RLS e de política restritiva passando.
