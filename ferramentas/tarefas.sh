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

infra() {
  titulo "Subindo PostgreSQL e Valkey"
  "${COMPOSE[@]}" up -d
  printf 'aguardando o banco'
  for _ in $(seq 1 60); do
    if "${COMPOSE[@]}" exec -T postgres pg_isready -U tenaz_migrador -d tenaz >/dev/null 2>&1; then
      echo; verde "infraestrutura no ar"; return 0
    fi
    printf '.'; sleep 1
  done
  echo; vermelho "o banco não respondeu no prazo"; return 1
}

infra_parar()  { "${COMPOSE[@]}" stop; }
infra_limpar() { "${COMPOSE[@]}" down -v; }
infra_estado() { "${COMPOSE[@]}" ps; }

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
  infra                sobe PostgreSQL e Valkey locais
  infra-parar          para os contêineres, preservando os dados
  infra-limpar         para e APAGA os dados
  infra-estado         mostra o estado dos contêineres
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
