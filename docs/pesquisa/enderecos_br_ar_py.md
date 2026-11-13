# Endereço multipaís: código postal e subdivisão (BR, AR, PY) (consulta de 13/11/2026)

Nota de pesquisa para o endereço da unidade (documento 3, seção 2.18: `subdivisao` e `codigo_postal`). Não é parecer. Rótulos de certeza: **[oficial]** fonte oficial lida; **[secundária]** fonte não oficial coerente entre si; **[dúvida]** precisa de confirmação. Limite desta consulta: sites de `wikipedia.org` (en) e `codigopostal.paraguay.gov.py` não abriram na ferramenta; a ISO Online Browsing Platform não foi consultada diretamente.

## 1. Regra geral proposto para o perfil de país

`codigo_postal` normalizado: maiúsculas, sem espaço nem hífen. `subdivisao`: código ISO 3166-2 completo, lista fechada por país.

| País | Regex de `codigo_postal` (normalizado) | Certeza |
|---|---|---|
| BR | `^\d{8}$` (CEP, 8 dígitos; "12345-678" vira "12345678") | Prática consolidada dos Correios [secundária]; confirmar no site dos Correios se quiser citar |
| AR | `^([A-HJ-NP-Z]\d{4}[A-Z]{3}|\d{4})$` | CPA: [secundária]; 4 dígitos: [secundária], ver 3 |
| PY | `^\d{6}$` | Decreto 8718 e Correo Paraguayo: [oficial por citação], ver 2 |

## 2. Paraguai

- **Forma:** seis dígitos numéricos, sem letras. Dois primeiros = departamento, dois seguintes = distrito, dois últimos = bairro/localidade; em Asunción os dois primeiros são "00" (ex. `001526`, `001518`). Fonte: Correo Paraguayo, por reportagens e guias que citam o texto do **Decreto n.º 8718** (governo Cartes), que define o código como conjunto de caracteres numéricos de seis dígitos; vigência 60 dias após a publicação. Data exata do decreto e texto do artigo **não conferidos** [dúvida]. Sistema novo desde cerca de junho de 2018 [secundária].
- **Quem define:** Dirección Nacional de Correos del Paraguay (DINACOPA, "Correo Paraguayo"). Consulta oficial indicada pelo próprio Correo: https://codigopostal.paraguay.gov.py/ (não abriu na ferramenta; consultado só por referência).
- **Uso corrente em endereço de empresa:** no endereço aparece antes da cidade (`110935 Luque`), com o departamento fora de Asunción (geopostcodes, secundária, consulta 13/11/2026). Uso corrente real entre empresas: [dúvida]. Recomendo campo **opcional** no PY, formato validado quando preenchido.
- Fontes: https://www.geopostcodes.com/country/paraguay/address-format/ ; https://codigo-postal.org/en-us/paraguay/ ; https://correoparaguayo.gov.py/sitio/wp-content/uploads/2025/08/Compendioeneroadiciembrede2024.pdf (resultado de busca, não lido); https://www.lanacion.com.py/pais/2024/06/21/paraguay-conectado-al-mundo-a-traves-de-dinacopa-entre-192-paises/ (só confirma o papel da DINACOPA; não cita o código).

Subdivisões ISO 3166-2:PY (fonte: es.wikipedia ISO 3166-2:PY, consulta 13/11/2026 [secundária, espelha a ISO]; 17 departamentos + capital = 18 códigos):

| Código | Nome | Código | Nome |
|---|---|---|---|
| PY-ASU | Asunción (distrito capital) | PY-09 | Paraguarí |
| PY-1 | Concepción | PY-10 | Alto Paraná |
| PY-2 | San Pedro | PY-11 | Central |
| PY-3 | Cordillera | PY-12 | Ñeembucú |
| PY-4 | Guairá | PY-13 | Amambay |
| PY-5 | Caaguazú | PY-14 | Canindeyú |
| PY-6 | Caazapá | PY-15 | Presidente Hayes |
| PY-7 | Itapúa | PY-16 | Alto Paraguay |
| PY-8 | Misiones | PY-19 | Boquerón |

**Atenção ao formato do código:** a tabela da es.wikipedia mostrou `PY-01`..`PY-09` com zero à esquerda, enquanto outra fonte (resultado de busca em espelhos) e a conhecida forma da ISO usam `PY-1`..`PY-9` sem zero. Não consegui resolver na ISO OBP: **[dúvida]**. Minha leitura: a ISO usa `PY-1` a `PY-9` (um dígito); peça ao dono confirmar na ISO OBP (https://www.iso.org/obp/ui/#iso:code:3166:PY) antes de fixar a lista. Proposta: aceitar na validação apenas a forma escolhida, sem tolerar as duas. Obsoletos, não aceitar: PY-17 (Chaco) e PY-18 (Nueva Asunción), extintos em 1992.

Regex de `subdivisao` PY (forma sem zero): `^PY-(ASU|[1-9]|1[0-6]|19)$`.

## 3. Argentina

- **CPA (Código Postal Argentino):** 8 caracteres, `LNNNNLLL`: 1 letra da província, 4 dígitos (o código postal antigo), 3 letras (face de quadra). Exemplo `C1425DKF`. Criado pelo Correo Argentino (1994 segundo uma fonte; adotado nacionalmente em 1998/1999, Resolução 1368/1998 da Secretaría de Comunicaciones segundo uma fonte secundária; **norma não conferida** [dúvida]). Letras possíveis da província: as do ISO 3166-2:AR (sem I, O, Ñ). Consulta oficial: correoargentino.com.ar/formularios/cpa (por referência, não aberta).
- **Código antigo de 4 dígitos:** segue em uso como parte do CPA e é aceito na prática pelos operadores para localidades pequenas [secundária]. Não encontrei fonte oficial dizendo que seja inválido. Proposta: **aceitar os dois**, normalizando só o formato (maiúsculas, sem espaço/hífen).
- Validação cruzada útil (regra da Tenaz, não da lei): primeira letra do CPA = letra da `subdivisao` (`AR-C` com `C1425DKF`). A letra do CPA coincide com a ISO segundo as fontes [secundária]; para CPA de 4 dígitos não há cruzamento.
- Fontes: https://en.wikipedia.org/wiki/Postal_codes_in_Argentina ; https://infonegocios.info/plus/el-cpa-codigo-postal-argentino-esta-mas-vivo-que-nunca-infonegocios-x5009lfg-y-tu-empresa (ago/2012, não confirma validade atual do código de 4 dígitos).

Subdivisões ISO 3166-2:AR (fonte: https://statoids.com/uar.html, consulta 13/11/2026 [secundária]; 23 províncias + Ciudad Autónoma de Buenos Aires = 24):

| Código | Nome | Código | Nome |
|---|---|---|---|
| AR-A | Salta | AR-Q | Neuquén |
| AR-B | Buenos Aires | AR-R | Río Negro |
| AR-C | Ciudad Autónoma de Buenos Aires | AR-S | Santa Fe |
| AR-D | San Luis | AR-T | Tucumán |
| AR-E | Entre Ríos | AR-U | Chubut |
| AR-F | La Rioja | AR-V | Tierra del Fuego |
| AR-G | Santiago del Estero | AR-W | Corrientes |
| AR-H | Chaco | AR-X | Córdoba |
| AR-J | San Juan | AR-Y | Jujuy |
| AR-K | Catamarca | AR-Z | Santa Cruz |
| AR-L | La Pampa | AR-M | Mendoza |
| AR-N | Misiones | AR-P | Formosa |

Regex de `subdivisao` AR: `^AR-[A-HJ-NP-Z]$` (24 letras exatas: A-H, J-N, P-Z). Statoids registra o nome da capital como "Ciudad de Buenos Aires"; a ISO usa "Ciudad Autónoma de Buenos Aires" (Ciudad Autónoma, conforme Constitución art. 129): confirmar o rótulo.

## 4. Brasil

ISO 3166-2:BR (fonte: Wikipédia en, consulta 13/11/2026 [secundária]): 26 estados + Distrito Federal = 27 códigos, e **coincidem com as UFs**: AC, AL, AP, AM, BA, CE, DF, ES, GO, MA, MT, MS, MG, PA, PB, PR, PE, PI, RJ, RN, RS, RO, RR, SC, SP, SE, TO. Nomes: Acre, Alagoas, Amapá, Amazonas, Bahia, Ceará, Distrito Federal, Espírito Santo, Goiás, Maranhão, Mato Grosso, Mato Grosso do Sul, Minas Gerais, Pará, Paraíba, Paraná, Pernambuco, Piauí, Rio de Janeiro, Rio Grande do Norte, Rio Grande do Sul, Rondônia, Roraima, Santa Catarina, São Paulo, Sergipe, Tocantins.

Regex de `subdivisao` BR: `^BR-(AC|AL|AP|AM|BA|CE|DF|ES|GO|MA|MT|MS|MG|PA|PB|PR|PE|PI|RJ|RN|RS|RO|RR|SC|SP|SE|TO)$`.

## 5. Casos para `identificacao.json`

- PY: `001526` ok; `110935` ok; `00152` falha; `0015A6` falha; `001-526` normaliza para `001526`.
- AR: `C1425DKF` ok; `c1425 dkf` normaliza e passa; `1425` ok; `I1425DKF` falha (letra I); `C142DKF` falha; `C1425DK` falha.
- BR: `01310100` ok; `01310-100` normaliza e passa; `1310100` falha.
- `subdivisao`: `AR-C` ok; `AR-I` falha; `PY-ASU` ok; `PY-17` falha; `BR-SP` ok; `BR-XX` falha; `SP` sozinho falha.

## 6. Pendências para o apoio externo
1. Advogado/contador do PY: confirmar o texto do Decreto 8718 e se o código postal é exigido em algum documento fiscal (SIFEN) ou só uso postal.
2. ~~Confirmar na ISO OBP o formato `PY-1` versus `PY-01` e o nome oficial de AR-C.~~ **Resolvido em 13/11/2026** pelo coordenador, na base `iso_3166-2.json` do projeto iso-codes do Debian (https://salsa.debian.org/iso-codes-team/iso-codes, ramo `main`), que espelha a ISO 3166-2: o Paraguai tem 18 códigos, **sem zero à esquerda** (`PY-1` a `PY-16`, `PY-19`, `PY-ASU`); a Argentina, 24 (`AR-A` a `AR-Z`, sem I nem O), com `AR-C` = "Ciudad Autónoma de Buenos Aires"; o Brasil, 27, iguais às UFs. Aceita-se só a forma sem zero. Fonte secundária confiável (não é a ISO OBP, que não abriu na ferramenta); se o apoio externo mostrar divergência, vale a ISO.
3. Confirmar a norma que institui o CPA e a validade do código de 4 dígitos (Correo Argentino).
