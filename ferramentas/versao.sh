#!/usr/bin/env bash
# Deriva a versão do release a partir da tag git e a imprime como pares
# chave=valor, para o build injetar nos artefatos.
#
# Regra de RNF-039: a versão NUNCA é digitada à mão em arquivo de código.
# Vem da tag `vX.Y.Z`. Fora de uma tag, a versão ganha sufixo de
# pré-lançamento com a contagem de commits e o commit curto, e produção
# recusa versão de pré-lançamento.
#
# Uso:
#   ferramentas/versao.sh              # pares chave=valor
#   ferramentas/versao.sh --versao     # só a versão
#   eval "$(ferramentas/versao.sh --export)"
set -euo pipefail

if ! git rev-parse --git-dir >/dev/null 2>&1; then
  echo "erro: fora de um repositório git" >&2
  exit 1
fi

COMMIT_LONGO="$(git rev-parse HEAD)"
COMMIT="$(git rev-parse --short=7 HEAD)"
DATA_BUILD="$(date -u +%Y-%m-%dT%H:%M:%SZ)"

# Árvore suja não gera versão limpa: o artefato não corresponderia ao commit.
# Usa `git status --porcelain`, que conta arquivo não rastreado, para bater com
# o que o próprio Go grava em vcs.modified no binário. Com `git diff-index`,
# um arquivo novo não rastreado passava limpo aqui e sujo no Go, e a
# verificação de coerência acusava divergência sem motivo aparente.
SUJO=""
if [[ -n "$(git status --porcelain 2>/dev/null)" ]]; then
  SUJO="+sujo"
fi

TAG_EXATA="$(git describe --tags --exact-match 2>/dev/null || true)"

if [[ -n "$TAG_EXATA" ]]; then
  if [[ ! "$TAG_EXATA" =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-rc\.[0-9]+)?$ ]]; then
    echo "erro: tag $TAG_EXATA fora do formato vMAJOR.MINOR.PATCH ou vMAJOR.MINOR.PATCH-rc.N" >&2
    exit 1
  fi
  VERSAO="${TAG_EXATA#v}${SUJO}"
else
  ULTIMA="$(git describe --tags --abbrev=0 2>/dev/null || true)"
  if [[ -n "$ULTIMA" ]]; then
    DISTANCIA="$(git rev-list --count "${ULTIMA}..HEAD")"
    BASE="${ULTIMA#v}"
    # Próxima versão de desenvolvimento: incrementa o MINOR da última tag.
    IFS='.' read -r MA MI _PA <<<"${BASE%%-*}"
    VERSAO="${MA}.$((MI + 1)).0-dev.${DISTANCIA}+${COMMIT}${SUJO}"
  else
    DISTANCIA="$(git rev-list --count HEAD)"
    VERSAO="0.1.0-dev.${DISTANCIA}+${COMMIT}${SUJO}"
  fi
fi

AMBIENTE="${TENAZ_AMBIENTE:-desenvolvimento}"

case "${1:-}" in
  --versao)
    echo "$VERSAO"
    ;;
  --commit)
    echo "$COMMIT"
    ;;
  --ldflags)
    P="tenaz/internal/plataforma/versao"
    echo "-X ${P}.Versao=${VERSAO} -X ${P}.Commit=${COMMIT} -X ${P}.DataBuild=${DATA_BUILD} -X ${P}.Ambiente=${AMBIENTE}"
    ;;
  --export)
    echo "export TENAZ_VERSAO='${VERSAO}'"
    echo "export TENAZ_COMMIT='${COMMIT}'"
    echo "export TENAZ_COMMIT_LONGO='${COMMIT_LONGO}'"
    echo "export TENAZ_DATA_BUILD='${DATA_BUILD}'"
    ;;
  *)
    echo "versao=${VERSAO}"
    echo "commit=${COMMIT}"
    echo "commit_longo=${COMMIT_LONGO}"
    echo "data_build=${DATA_BUILD}"
    echo "ambiente=${AMBIENTE}"
    ;;
esac
