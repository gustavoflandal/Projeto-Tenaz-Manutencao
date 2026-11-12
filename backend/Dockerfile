# Imagem do backend da Tenaz.
#
# Um binário, vários modos: a mesma imagem roda api, worker, sincronizador,
# ingestor e migrate (seção 4.1 de docs/5_ARQUITETURA_SOFTWARE.md).
#
# A imagem é construída UMA vez e promovida entre ambientes pelo digest, nunca
# reconstruída para produção (seção 9.4).
#
# Construção a partir da raiz do repositório:
#   docker build -f backend/Dockerfile \
#     --build-arg VERSAO=1.0.0 --build-arg COMMIT=a1b2c3d \
#     --build-arg DATA_BUILD=2026-11-12T14:05:00Z -t tenaz:1.0.0 backend

FROM golang:1.27-alpine AS construcao

# A versão NUNCA é digitada em arquivo de código: entra pelo build, derivada
# da tag git por ferramentas/versao.sh (RNF-039).
ARG VERSAO=0.0.0-dev
ARG COMMIT=desconhecido
ARG DATA_BUILD=
ARG AMBIENTE=producao

WORKDIR /src

# Camada de dependências separada: muda só quando go.mod ou go.sum mudam.
COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .

ENV CGO_ENABLED=0 GOOS=linux
RUN go build \
      -trimpath \
      -ldflags "-s -w \
        -X tenaz/internal/plataforma/versao.Versao=${VERSAO} \
        -X tenaz/internal/plataforma/versao.Commit=${COMMIT} \
        -X tenaz/internal/plataforma/versao.DataBuild=${DATA_BUILD} \
        -X tenaz/internal/plataforma/versao.Ambiente=${AMBIENTE}" \
      -o /tenaz ./cmd/tenaz

# Falha a construção se a versão não tiver sido injetada: imagem sem versão
# rastreável não vai para nenhum ambiente.
RUN /tenaz versao | grep -q "^versao=${VERSAO}$"

# ---------------------------------------------------------------------------

FROM gcr.io/distroless/static-debian12:nonroot

# Imagem mínima, sem shell e sem gerenciador de pacotes: menos superfície de
# ataque e nada para um invasor usar se chegar a executar código.
COPY --from=construcao /tenaz /tenaz

# Usuário sem privilégio, já definido na imagem base.
USER nonroot:nonroot

EXPOSE 8080

ENV TENAZ_AMBIENTE=producao \
    TENAZ_ENDERECO_HTTP=:8080

ENTRYPOINT ["/tenaz"]
CMD ["api"]
