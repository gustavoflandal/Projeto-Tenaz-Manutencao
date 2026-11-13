#!/usr/bin/env bash
# Tarefas de desenvolvimento e de CI da Tenaz.
#
# O CI chama exatamente estes subcomandos, para que "passa na minha máquina" e
# "passa no CI" signifiquem a mesma coisa. Os scripts do package.json na raiz
# apontam para cá, conforme a decisão de orquestração da seção 7 de
# docs/5_ARQUITETURA_SOFTWARE.md.
#
# Roda em Git Bash no Windows e em bash no Linux. Não usa make, que não existe
# por padrão no Windows.
#
#   ferramentas/tarefas.sh ajuda
set -euo pipefail

RAIZ="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$RAIZ"

COMPOSE=(docker compose -f infra/docker/compose.yaml)
BANCO_DEV="postgres://tenaz_migrador:tenaz_dev@127.0.0.1:5433/tenaz?sslmode=disable"

vermelho() { printf '\033[31m%s\033[0m\n' "$*" >&2; }
verde()    { printf '\033[32m%s\033[0m\n' "$*"; }
titulo()   { printf '\n\033[1m== %s\033[0m\n' "$*"; }

# ----------------------------------------------------------------------- infra

# Sobe tudo e espera cada serviço ficar saudável pelo healthcheck do compose;
# depois prepara o bucket da trilha de auditoria.
infra() {
  titulo "Subindo PostgreSQL, Valkey, Mailpit e SeaweedFS (S3)"
  if ! "${COMPOSE[@]}" up -d --wait --wait-timeout 180; then
    "${COMPOSE[@]}" ps
    vermelho "algum serviço não ficou saudável no prazo"
    return 1
  fi
  "${COMPOSE[@]}" ps
  armazenamento_preparar
  verde "infraestrutura no ar"
}

infra_parar()  { "${COMPOSE[@]}" stop; }
infra_limpar() { "${COMPOSE[@]}" down -v; }
infra_estado() { "${COMPOSE[@]}" ps; }

# ------------------------------------------------------------- armazenamento
#
# SeaweedFS do compose, falado só pela API S3 (ADR-012, decisões 4 e 5, em
# docs/5_ARQUITETURA_SOFTWARE.md). O cliente é o contêiner `s3-cli` do compose
# (aws-cli com versão e digest fixos): nada é instalado no host. As credenciais
# abaixo são as de desenvolvimento de infra/docker/seaweedfs/s3.json.

S3_BUCKET_AUDITORIA="tenaz-auditoria"
# Retenção padrão do bucket, em modo de conformidade. Em produção é a da trilha
# (5 anos no mínimo, RNF-019); aqui é 1 dia, para que o objeto de teste não
# fique preso por anos. `infra-limpar` apaga o volume e, com ele, tudo.
S3_RETENCAO_DIAS_DEV=1
S3_AUDITORIA_CHAVE="tenaz_auditoria"
S3_AUDITORIA_SEGREDO="tenaz_auditoria_dev_secreto"

# s3 <args do aws>: roda o aws-cli na rede do compose, como `tenaz_operacao`.
# MSYS_NO_PATHCONV impede o Git Bash de reescrever caminhos como /tmp/x.
s3() {
  MSYS_NO_PATHCONV=1 "${COMPOSE[@]}" --profile ferramentas run --rm --no-deps -T s3-cli "$@"
}

# s3_roteiro <args>: roda no contêiner do cliente o roteiro bash lido da
# entrada padrão, com o aws e o curl da imagem.
s3_roteiro() {
  MSYS_NO_PATHCONV=1 "${COMPOSE[@]}" --profile ferramentas run --rm --no-deps -T \
    -e S3_AUDITORIA_CHAVE="$S3_AUDITORIA_CHAVE" -e S3_AUDITORIA_SEGREDO="$S3_AUDITORIA_SEGREDO" \
    --entrypoint bash s3-cli -s -- "$@"
}

# Cria o bucket da trilha, se não existir, com bloqueio de objeto (que liga o
# versionamento) e retenção padrão COMPLIANCE; e confere os dois, exista ou não.
# Idempotente: rodar de novo não muda nada.
armazenamento_preparar() {
  titulo "Bucket $S3_BUCKET_AUDITORIA (bloqueio de objeto, versionamento)"
  s3_roteiro "$S3_BUCKET_AUDITORIA" "$S3_RETENCAO_DIAS_DEV" <<'ROTEIRO'
set -euo pipefail
B="$1"; DIAS="$2"
if aws s3api head-bucket --bucket "$B" >/dev/null 2>&1; then
  echo "bucket $B já existe"
else
  aws s3api create-bucket --bucket "$B" --object-lock-enabled-for-bucket >/dev/null
  aws s3api put-object-lock-configuration --bucket "$B" --object-lock-configuration \
    "{\"ObjectLockEnabled\":\"Enabled\",\"Rule\":{\"DefaultRetention\":{\"Mode\":\"COMPLIANCE\",\"Days\":$DIAS}}}"
  echo "bucket $B criado"
fi
versionamento="$(aws s3api get-bucket-versioning --bucket "$B" --query Status --output text)"
bloqueio="$(aws s3api get-object-lock-configuration --bucket "$B" \
  --query 'ObjectLockConfiguration.[ObjectLockEnabled,Rule.DefaultRetention.Mode,Rule.DefaultRetention.Days]' --output text)"
echo "versionamento: $versionamento"
echo "bloqueio de objeto: $bloqueio"
[[ "$versionamento" == "Enabled" ]] || { echo "versionamento desligado em $B" >&2; exit 1; }
[[ "$bloqueio" == Enabled*COMPLIANCE* ]] || { echo "bloqueio de objeto fora do modo COMPLIANCE em $B" >&2; exit 1; }
ROTEIRO
  verde "bucket pronto"
}

# Prova o endurecimento do ADR-012 contra o contêiner do compose: o objeto em
# retenção COMPLIANCE não é apagado, encurtado, rebaixado nem sobrescrito pela
# API S3, nem pela identidade Admin; a identidade da trilha não mexe no
# bloqueio nem apaga o bucket; o anônimo não lê; a URL pré-assinada vale e a
# adulterada não; e a API do filer não é alcançável pela rede e, sem JWT,
# recebe 401. É a semente do teste de integração do armazenamento (decisão 5)
# e qualifica toda versão nova da imagem. Cada execução deixa dois objetos
# pequenos em ensaio/, retidos por 1 hora e 1 dia; somem com `infra-limpar`.
armazenamento_provar() {
  titulo "Prova do armazenamento (ADR-012): bloqueio de objeto e filer"
  local chave saida versao codigo
  chave="ensaio/$(date -u +%Y%m%dT%H%M%SZ)-$RANDOM.selo"

  if ! saida="$(s3_roteiro "$S3_BUCKET_AUDITORIA" "$chave" <<'ROTEIRO'
set -uo pipefail
B="$1"; K="$2"; FALHAS=0
ok()    { printf '  ok     %s\n' "$*"; }
falha() { printf '  FALHA  %s\n' "$*"; FALHAS=$((FALHAS + 1)); }
aud()   { AWS_ACCESS_KEY_ID="$S3_AUDITORIA_CHAVE" AWS_SECRET_ACCESS_KEY="$S3_AUDITORIA_SEGREDO" aws "$@"; }
# negado <descrição> <comando...>: o comando tem de falhar com AccessDenied.
negado() {
  local d="$1"; shift; local s rc
  s="$("$@" 2>&1)"; rc=$?
  if [[ $rc -ne 0 && "$s" == *AccessDenied* ]]; then ok "$d: AccessDenied"; else falha "$d: veio rc=$rc $s"; fi
}
ATE="$(date -u -d '+1 hour' +%Y-%m-%dT%H:%M:%SZ)"
CURTO="$(date -u -d '+5 minutes' +%Y-%m-%dT%H:%M:%SZ)"

[[ "$(aws s3api get-bucket-versioning --bucket "$B" --query Status --output text)" == Enabled ]] \
  && ok "1. versionamento ligado" || falha "1. versionamento"
[[ "$(aws s3api get-object-lock-configuration --bucket "$B" --query ObjectLockConfiguration.Rule.DefaultRetention.Mode --output text)" == COMPLIANCE ]] \
  && ok "2. retenção padrão COMPLIANCE" || falha "2. retenção padrão"

printf 'selo-original' > /tmp/original
V="$(aud s3api put-object --bucket "$B" --key "$K" --body /tmp/original \
  --object-lock-mode COMPLIANCE --object-lock-retain-until-date "$ATE" --query VersionId --output text 2>&1)"
[[ -n "$V" && "$V" != *rror* ]] && ok "3. identidade da trilha grava o selo em COMPLIANCE até $ATE (versão $V)" || falha "3. gravar: $V"
[[ "$(aud s3api get-object-retention --bucket "$B" --key "$K" --version-id "$V" --query Retention.Mode --output text)" == COMPLIANCE ]] \
  && ok "4. retenção da versão lida de volta: COMPLIANCE" || falha "4. retenção da versão"

negado "5. Admin apaga a versão protegida" aws s3api delete-object --bucket "$B" --key "$K" --version-id "$V"
negado "6. Admin apaga com bypass de governança" aws s3api delete-object --bucket "$B" --key "$K" --version-id "$V" --bypass-governance-retention
negado "7. Admin encurta a retenção" aws s3api put-object-retention --bucket "$B" --key "$K" --version-id "$V" \
  --retention "{\"Mode\":\"COMPLIANCE\",\"RetainUntilDate\":\"$CURTO\"}"
negado "8. Admin rebaixa para GOVERNANCE" aws s3api put-object-retention --bucket "$B" --key "$K" --version-id "$V" \
  --retention "{\"Mode\":\"GOVERNANCE\",\"RetainUntilDate\":\"$ATE\"}" --bypass-governance-retention
negado "9. identidade da trilha muda o bloqueio do bucket" aud s3api put-object-lock-configuration --bucket "$B" \
  --object-lock-configuration '{"ObjectLockEnabled":"Enabled","Rule":{"DefaultRetention":{"Mode":"GOVERNANCE","Days":1}}}'
negado "10. identidade da trilha apaga o bucket" aud s3api delete-bucket --bucket "$B"

printf 'selo-FORJADO' > /tmp/forjado
V2="$(aud s3api put-object --bucket "$B" --key "$K" --body /tmp/forjado --query VersionId --output text 2>&1)"
[[ -n "$V2" && "$V2" != "$V" && "$V2" != *rror* ]] && ok "11. sobrescrever cria versão nova ($V2), a protegida fica" || falha "11. sobrescrever: $V2"
[[ "$(aws s3api delete-object --bucket "$B" --key "$K" --query DeleteMarker --output text 2>&1)" == True ]] \
  && ok "12. apagar sem versão só cria marcador de exclusão" || falha "12. marcador de exclusão"
if s="$(aws s3api delete-bucket --bucket "$B" 2>&1)"; then falha "13. Admin apagou o bucket"
else ok "13. Admin apaga o bucket: recusado ($(grep -o '([A-Za-z]*)' <<<"$s" | head -1 | tr -d '()'))"; fi

codigo="$(curl -s -o /dev/null -w '%{http_code}' "$AWS_ENDPOINT_URL/$B/$K")"
[[ "$codigo" == 403 ]] && ok "14. leitura anônima: HTTP 403" || falha "14. leitura anônima: HTTP $codigo"

printf 'anexo' > /tmp/anexo
aud s3api put-object --bucket "$B" --key "$K.anexo" --body /tmp/anexo >/dev/null
URL="$(aud s3 presign "s3://$B/$K.anexo" --expires-in 300)"
codigo="$(curl -s -o /dev/null -w '%{http_code}' "$URL")"
[[ "$codigo" == 200 ]] && ok "15. URL pré-assinada de 5 minutos: HTTP 200" || falha "15. pré-assinada: HTTP $codigo"
codigo="$(curl -s -o /dev/null -w '%{http_code}' "${URL/.anexo/}")"
[[ "$codigo" == 403 ]] && ok "16. URL pré-assinada adulterada: HTTP 403" || falha "16. adulterada: HTTP $codigo"

codigo="$(curl -s -m 5 -o /dev/null -w '%{http_code}' -X DELETE "http://seaweedfs:8888/buckets/$B?recursive=true&ignoreRecursiveError=true")"; rc=$?
[[ $rc -ne 0 && "$codigo" == 000 ]] && ok "17. API do filer (8888) inalcançável pela rede do compose (curl rc=$rc)" \
  || falha "17. filer respondeu pela rede: HTTP $codigo"

echo "VERSAO=$V"
exit $((FALHAS > 0))
ROTEIRO
  )"; then
    echo "$saida"; vermelho "prova do armazenamento falhou"; return 1
  fi
  echo "$saida" | grep -v '^VERSAO='
  versao="$(echo "$saida" | sed -n 's/^VERSAO=//p' | tr -d '\r')"

  # De dentro do próprio contêiner, onde o filer escuta: sem JWT, 401.
  codigo="$(MSYS_NO_PATHCONV=1 "${COMPOSE[@]}" exec -T seaweedfs curl -s -o /dev/null -w '%{http_code}' \
    -X DELETE "http://127.0.0.1:8888/buckets/$S3_BUCKET_AUDITORIA?recursive=true&ignoreRecursiveError=true")"
  if [[ "$codigo" != 401 ]]; then
    vermelho "  FALHA  18. DELETE recursivo no filer sem JWT: HTTP $codigo (esperado 401)"; return 1
  fi
  echo "  ok     18. DELETE recursivo no filer sem JWT, de dentro do contêiner: HTTP 401"

  if ! saida="$(s3_roteiro "$S3_BUCKET_AUDITORIA" "$chave" "$versao" <<'ROTEIRO'
set -euo pipefail
aws s3api get-object --bucket "$1" --key "$2" --version-id "$3" /tmp/lido >/dev/null
cat /tmp/lido
ROTEIRO
  )"; then
    echo "$saida"; vermelho "  FALHA  19. a versão protegida sumiu"; return 1
  fi
  if [[ "$(echo "$saida" | tr -d '\r')" != "selo-original" ]]; then
    vermelho "  FALHA  19. conteúdo da versão protegida mudou: $saida"; return 1
  fi
  echo "  ok     19. depois de tudo, a versão protegida continua lá com o conteúdo original"
  verde "armazenamento endurecido: 19 de 19"
}

migrar() {
  titulo "Aplicando migrações"
  (cd backend && TENAZ_BANCO_URL="$BANCO_DEV" go run ./cmd/tenaz migrate)
}

migrar_situacao() {
  (cd backend && TENAZ_BANCO_URL="$BANCO_DEV" go run ./cmd/tenaz migrate-status)
}

# --------------------------------------------------------------------- backend

compilar() {
  titulo "Compilando com a versão injetada pela tag"
  mkdir -p bin
  (cd backend && go build -ldflags "$(bash "$RAIZ/ferramentas/versao.sh" --ldflags)" -o ../bin/tenaz ./cmd/tenaz)
  verde "bin/tenaz"
}

# Ambiente da API local. As chaves abaixo são SÓ de desenvolvimento, estão no
# repositório de propósito e não valem em nenhum outro ambiente; em
# homologação e produção vêm do cofre. O pool da aplicação assume tenaz_app e o
# da autenticação, sem TENAZ_BANCO_URL_AUTENTICACAO, assume tenaz_autenticacao
# (ADR-011), os dois sobre a URL do migrador. Sobre http, o cookie de sessão
# sai como tenaz_sessao, sem Secure (TENAZ_COOKIE_SEGURO=false, recusado fora
# do desenvolvimento; seção 4.2 do documento 3).
api() {
  (cd backend && \
    TENAZ_BANCO_URL="$BANCO_DEV" \
    TENAZ_BANCO_PAPEL="tenaz_app" \
    TENAZ_BANCO_AUTENTICACAO_MAX_CONN="5" \
    TENAZ_URL_WEB="http://localhost:5173" \
    TENAZ_COOKIE_SEGURO="false" \
    TENAZ_CHAVE_MFA="HwN8RJDGAV6P2aR295hByxI8whzgK2LmY7w4HqabfXI=" \
    TENAZ_CHAVE_JWT="bYC9XDoeAyjFKg4SQh2/DJnEACmHjXLJSc5nz5Ebx+c=" \
    TENAZ_CHAVE_JWT_ID="dev-1" \
    TENAZ_SMTP_ENDERECO="127.0.0.1:1025" \
    TENAZ_EMAIL_REMETENTE="nao-responda@tenaz.local" \
    go run ./cmd/tenaz api)
}

formatar() {
  (cd backend && gofmt -w ./cmd ./internal ./testes ./db)
  verde "formatado"
}

conferir_formato() {
  titulo "Formato do Go"
  local saida
  saida="$(cd backend && gofmt -l ./cmd ./internal ./testes ./db)"
  if [[ -n "$saida" ]]; then
    vermelho "arquivos fora do formato:"
    echo "$saida" >&2
    vermelho "rode: ferramentas/tarefas.sh formatar"
    return 1
  fi
  verde "formato ok"
}

vet() {
  titulo "Análise estática"
  (cd backend && go vet ./...)
  verde "análise ok"
}

testar() {
  titulo "Testes (contêiner descartável)"
  (cd backend && go test ./... -count=1)
}

testar_rapido() {
  titulo "Testes reaproveitando o banco local"
  (cd backend && TENAZ_TESTE_BANCO_URL="$BANCO_DEV" go test ./... -count=1)
}

# Falha aqui BLOQUEIA merge e deploy, e é tratada como incidente de segurança
# de severidade alta, mesmo acontecendo só no CI.
vazamento() {
  titulo "Suíte de isolamento (RNF-013, RNF-037, RNF-041)"
  (cd backend && go test ./testes/vazamento/ -count=1 -v)
}

# -------------------------------------------------------------------- contrato

contrato() {
  titulo "Regerando o contrato OpenAPI"
  (cd backend && go run ./cmd/tenaz openapi > api/openapi.yaml)
  verde "backend/api/openapi.yaml regerado"
  api_tipos
}

# Tipos TypeScript gerados do contrato, em pacotes/api-tipos. Quem muda o
# contrato roda `contrato`, que já chama esta tarefa.
api_tipos() {
  titulo "Regerando os tipos TypeScript do contrato"
  "${PNPM[@]}" --filter @tenaz/api-tipos gerar
}

conferir_contrato() {
  titulo "Contrato OpenAPI em dia"
  local gerado
  gerado="$(mktemp)"
  (cd backend && go run ./cmd/tenaz openapi) > "$gerado"
  if ! diff -u backend/api/openapi.yaml "$gerado"; then
    rm -f "$gerado"
    vermelho "o contrato versionado difere do gerado pelo código."
    vermelho "alguém mudou a API sem regenerar. Rode: ferramentas/tarefas.sh contrato"
    return 1
  fi
  rm -f "$gerado"
  verde "contrato em dia"
}

# -------------------------------------------------------------- javascript

# pnpm pelo npx: `corepack enable` precisa de privilégio de administrador no
# Windows. Quem tiver pnpm global pode trocar PNPM por "pnpm".
PNPM=(npx --yes pnpm@10)

instalar_js() {
  titulo "Instalando dependências JavaScript"
  "${PNPM[@]}" install
}

testar_js() {
  titulo "Testes de JavaScript (web e app)"
  "${PNPM[@]}" -r test
}

tipos_js() {
  titulo "Tipos de JavaScript"
  "${PNPM[@]}" -r typecheck
}

lint_js() {
  titulo "Análise estática de JavaScript"
  "${PNPM[@]}" -r lint
}

tokens_ui() {
  titulo "Regerando o CSS dos tokens de design"
  "${PNPM[@]}" --filter @tenaz/tokens-ui gerar
}

# Arquivos gerados versionados: tipos do contrato e CSS dos tokens.
conferir_gerados_js() {
  titulo "Arquivos gerados de JavaScript em dia"
  "${PNPM[@]}" --filter @tenaz/api-tipos conferir
  "${PNPM[@]}" --filter @tenaz/tokens-ui conferir
}

compilar_web() {
  titulo "Compilando a web com a versão injetada"
  eval "$(bash "$RAIZ/ferramentas/versao.sh" --export)"
  (cd web && TENAZ_VERSAO="$TENAZ_VERSAO" TENAZ_COMMIT="$TENAZ_COMMIT"      TENAZ_DATA_BUILD="$TENAZ_DATA_BUILD" "${PNPM[@]}" exec vite build)
  verde "web/dist"
}

# ------------------------------------------------------------------ documentos

conferir_docs() {
  titulo "Integridade da documentação"
  python ferramentas/verificar_docs.py
}

versao() { bash "$RAIZ/ferramentas/versao.sh"; }

# ------------------------------------------------------------------- conjuntos

conferir() {
  conferir_formato
  vet
  conferir_contrato
  conferir_gerados_js
  tipos_js
  lint_js
  conferir_docs
  verde $'\nTodas as conferências passaram'
}

tudo() {
  conferir
  testar
  testar_js
  verde $'\nConferências e testes completos'
}

ajuda() {
  cat <<'TXT'
uso: ferramentas/tarefas.sh <tarefa>

Infraestrutura
  infra                sobe PostgreSQL, Valkey, Mailpit e SeaweedFS (S3), espera
                       todos saudáveis e prepara o bucket da trilha de auditoria
  infra-parar          para os contêineres, preservando os dados
  infra-limpar         para e APAGA os dados
  infra-estado         mostra o estado dos contêineres
  armazenamento-preparar  cria o bucket tenaz-auditoria, se faltar, com bloqueio
                       de objeto COMPLIANCE e versionamento, e confere os dois
  armazenamento-provar    prova o endurecimento do ADR-012: retenção COMPLIANCE
                       resiste à API S3 e o filer sem JWT recebe 401
  migrar               aplica as migrações no banco local
  migrar-situacao      mostra o estado das migrações

Backend
  compilar             compila com a versão injetada pela tag
  api                  sobe a API local
  formatar             formata o código Go
  conferir-formato     falha se houver código fora do formato
  vet                  análise estática
  testar               todos os testes, com contêiner descartável
  testar-rapido        todos os testes, reaproveitando o banco local
  vazamento            suíte de isolamento, com saída detalhada

Contrato
  contrato             regera backend/api/openapi.yaml e os tipos TypeScript
  conferir-contrato    falha se o contrato versionado divergir do gerado
  api-tipos            regera só os tipos TypeScript (pacotes/api-tipos)

JavaScript
  instalar-js          instala as dependências da web e do app
  testar-js            testes da web e do app
  tipos-js             conferência de tipos da web e do app
  lint-js              análise estática (ESLint) da web
  tokens-ui            regera o CSS dos tokens a partir de tokens.json
  conferir-gerados-js  falha se tipos da API ou CSS dos tokens divergirem
  compilar-web         compila a web com a versão injetada

Documentação
  conferir-docs        verifica links, IDs, versões, etapas e PDFs
  versao               mostra a versão derivada da tag git

Conjuntos
  conferir             formato, análise, contrato, gerados, tipos, lint e documentação
  tudo                 conferir + todos os testes (o que o CI roda)
TXT
}

tarefa="${1:-ajuda}"
shift || true
case "$tarefa" in
  infra)             infra "$@" ;;
  infra-parar)       infra_parar "$@" ;;
  infra-limpar)      infra_limpar "$@" ;;
  infra-estado)      infra_estado "$@" ;;
  armazenamento-preparar) armazenamento_preparar "$@" ;;
  armazenamento-provar)   armazenamento_provar "$@" ;;
  migrar)            migrar "$@" ;;
  migrar-situacao)   migrar_situacao "$@" ;;
  compilar)          compilar "$@" ;;
  api)               api "$@" ;;
  formatar)          formatar "$@" ;;
  conferir-formato)  conferir_formato "$@" ;;
  vet)               vet "$@" ;;
  testar)            testar "$@" ;;
  testar-rapido)     testar_rapido "$@" ;;
  vazamento)         vazamento "$@" ;;
  contrato)          contrato "$@" ;;
  conferir-contrato) conferir_contrato "$@" ;;
  api-tipos)         api_tipos "$@" ;;
  instalar-js)       instalar_js "$@" ;;
  testar-js)         testar_js "$@" ;;
  tipos-js)          tipos_js "$@" ;;
  lint-js)           lint_js "$@" ;;
  tokens-ui)         tokens_ui "$@" ;;
  conferir-gerados-js) conferir_gerados_js "$@" ;;
  compilar-web)      compilar_web "$@" ;;
  conferir-docs)     conferir_docs "$@" ;;
  versao)            versao "$@" ;;
  conferir)          conferir "$@" ;;
  tudo)              tudo "$@" ;;
  ajuda|-h|--help)   ajuda ;;
  *)                 vermelho "tarefa desconhecida: $tarefa"; echo; ajuda; exit 1 ;;
esac
