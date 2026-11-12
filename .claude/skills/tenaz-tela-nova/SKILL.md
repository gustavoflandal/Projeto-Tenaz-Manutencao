---
name: tenaz-tela-nova
description: Use sempre que nascer ou mudar uma tela na web ou no app da Tenaz. Garante tokens da marca, componentes do documento 7, contraste WCAG 2.1 AA, os quatro estados, microcopy no tom de voz, navegação por teclado, selo de versão visível, ajuda contextual e a captura de evidência.
---

# Tela nova na Tenaz

A tela é onde o sistema encontra uma pessoa de luva, no pátio, com uma mão livre, ou um planejador com 40 OS abertas. Os dois precisam entender a tela sem treinamento formal — é a proposta de valor 6 do [sumário](../../../docs/0_SUMARIO_EXECUTIVO.md).

## 1. Antes de desenhar

1. Procure a tela em [docs/7_PADROES_DESIGN.md](../../../docs/7_PADROES_DESIGN.md), seção 8 (telas principais) e seção 9 (app de campo). Muitas já estão especificadas.
2. Procure o componente no catálogo da seção 6. Componente novo é exceção que se justifica e se registra.
3. Confira o requisito: que dado a tela mostra, que ação permite, que permissão exige.
4. Confira o contrato da API: a tela só mostra campo que existe no contrato gerado. Falta de campo é assunto de `tenaz-api-contrato`, não motivo para inventar.

## 2. Tokens e marca

| Item | Regra |
|---|---|
| Cores | Só tokens. Aço `#14202E` e Sinal `#FF6A13` como base |
| Texto sobre Sinal | **Sempre Aço.** Nunca branco |
| Contorno de foco | Aço |
| Chip de status no tema claro | Texto em Aço |
| Texto terciário | Aço 400 |
| Tipografia | Archivo Expanded, Inter, JetBrains Mono |
| Espaçamento, raio, sombra | Só tokens |
| Textos | Só pelo pacote de internacionalização. Nenhum texto escrito direto no componente |

## 3. Status, prioridade e criticidade

Três coisas distintas, com cores distintas, que a interface mistura com frequência:

- **Status operacional:** `operando`, `em_manutencao`, `parado`, `inativo`. "Em manutenção" usa a cor de Preventiva.
- **Atenção:** **não é status.** É condição calculada sobre `operando` (plano vencendo, documento vencendo, leitura fora da faixa) e usa a cor de Atenção.
- **Apto para despacho:** verificação separada. Um veículo pode estar `operando` e impedido de despacho. São dois indicadores na tela.
- **Prioridade** (`baixa`, `media`, `alta`, `urgente`) e **criticidade** (A, B, C) são escalas diferentes e não compartilham paleta.

## 4. Os quatro estados

Toda tela tem os quatro desenhados antes de ser implementada:

| Estado | O que precisa |
|---|---|
| Padrão | Com massa de dados plausível, não com uma linha |
| Carregando | Sem salto de layout quando o conteúdo chega |
| Vazio | Texto que diz **o que a pessoa pode fazer**, não só "nenhum registro" |
| Erro | O que aconteceu e o que fazer, em português, sem culpar o usuário e sem jargão |

Mais os estados que o requisito exigir: validação, offline, bloqueado, somente leitura, sem permissão.

## 5. App de campo

- Alvo de toque grande; uso com uma mão; legível sob sol.
- Nada que dependa de passar o cursor.
- Ação de campo funciona offline, com [`tenaz-operacao-offline`](../tenaz-operacao-offline/SKILL.md).
- Número provisório mostrado como provisório.
- Indicação clara de que há coisa na fila de sincronização.

## 6. Acessibilidade

Meta WCAG 2.1 AA (RNF-025):

1. Rótulo em todo campo. Placeholder não é rótulo.
2. Navegação completa por teclado na web, com foco visível.
3. Contraste conferido com número, pela seção 4 do documento 7, inclusive texto sobre cor de status.
4. Ordem de leitura coerente; região e cabeçalho marcados.
5. Nenhuma informação transmitida **só** por cor. Status tem texto ou forma, além da cor.
6. Mensagem de erro associada ao campo.

## 7. Selo de versão e ajuda

- O selo de versão fica no **canto inferior direito** da tela principal e da tela de login, na web e no app (RNF-039). Nenhum componente, rodapé ou aviso o encobre.
- Ajuda contextual mínima na tela é o critério 8 dos nove critérios de pronto. Uma frase que explique o que a tela faz já cumpre; nenhuma não cumpre.

## 8. Fechar

1. Capture a evidência com [`tenaz-evidencia-visual`](../tenaz-evidencia-visual/SKILL.md): estado padrão, mais um por estado exigido, nos dois temas.
2. Rode a verificação automática de acessibilidade e faça a manual de teclado e leitor de tela.
3. Peça a revisão de `tenaz-design-sistema`. **Quem implementa não aprova o próprio desenho.**
4. [`tenaz-fechar-etapa`](../tenaz-fechar-etapa/SKILL.md).
