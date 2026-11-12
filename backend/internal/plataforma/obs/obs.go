// Package obs concentra observabilidade: log estruturado com redação.
//
// Campos obrigatórios de toda linha emitida dentro de requisição ou job, pela
// seção 14.3 de docs/5_ARQUITETURA_SOFTWARE.md: ts, nivel, msg, servico,
// versao, trace_id, span_id, empresa_id, usuario_id, request_id, rota, status,
// duracao_ms.
package obs

import (
	"context"
	"io"
	"log/slog"
	"strings"
	"time"

	"tenaz/internal/plataforma/versao"
)

// Campos canônicos. Usar estas constantes, nunca a string solta.
const (
	CampoServico   = "servico"
	CampoVersao    = "versao"
	CampoTraceID   = "trace_id"
	CampoSpanID    = "span_id"
	CampoEmpresaID = "empresa_id"
	CampoUsuarioID = "usuario_id"
	CampoRequestID = "request_id"
	CampoRota      = "rota"
	CampoStatus    = "status"
	CampoDuracaoMs = "duracao_ms"
)

// Redigido substitui o valor de um campo sensível.
const Redigido = "[REDIGIDO]"

// sensiveis são os nomes de campo cujo valor nunca vai para o log, pela seção
// 14.3: senha, token, chave de API, segredo de conector, conteúdo de documento
// sensível, CPF completo, resultado de exame.
var sensiveis = []string{
	"senha", "password", "token", "secret", "segredo", "chave_api", "api_key",
	"authorization", "cookie", "cpf", "exame", "toxicologico", "aso",
	"refresh_token", "access_token", "credencial", "private_key", "assinatura_hmac",
}

// Sensivel informa se um nome de campo deve ter o valor redigido.
//
// A comparação é por trecho, sem diferença de maiúsculas e com hífen, ponto e
// espaço normalizados para sublinhado. Sem essa normalização, o cabeçalho
// "X-Api-Key" escapava da lista, que registra "api_key".
func Sensivel(nome string) bool {
	n := strings.ToLower(nome)
	n = strings.NewReplacer("-", "_", ".", "_", " ", "_").Replace(n)
	for _, s := range sensiveis {
		if strings.Contains(n, s) {
			return true
		}
	}
	return false
}

// redator é um ReplaceAttr do slog que mascara campos sensíveis em qualquer
// profundidade e renomeia os campos padrão para os nomes canônicos.
func redator(_ []string, a slog.Attr) slog.Attr {
	switch a.Key {
	case slog.TimeKey:
		// UTC no log, como no banco e na API (princípio 6 da seção 1.2 de
		// docs/2_ARQUITETURA_BANCO_DADOS.md). Horário local em log de
		// servidor torna impossível correlacionar incidente entre regiões.
		a.Key = "ts"
		if t, ok := a.Value.Any().(time.Time); ok {
			a.Value = slog.TimeValue(t.UTC())
		}
		return a
	case slog.LevelKey:
		a.Key = "nivel"
		return a
	case slog.MessageKey:
		a.Key = "msg"
		return a
	}
	if Sensivel(a.Key) {
		return slog.String(a.Key, Redigido)
	}
	return a
}

// Novo monta o logger do processo.
func Novo(destino io.Writer, servico, nivel string, json bool) *slog.Logger {
	opcoes := &slog.HandlerOptions{Level: aNivel(nivel), ReplaceAttr: redator}

	var h slog.Handler
	if json {
		h = slog.NewJSONHandler(destino, opcoes)
	} else {
		h = slog.NewTextHandler(destino, opcoes)
	}

	return slog.New(h).With(
		CampoServico, servico,
		CampoVersao, versao.Versao,
	)
}

func aNivel(n string) slog.Level {
	switch strings.ToLower(n) {
	case "debug":
		return slog.LevelDebug
	case "warn", "aviso":
		return slog.LevelWarn
	case "error", "erro":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

type chaveLogger struct{}

// ComLogger guarda o logger no contexto.
func ComLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, chaveLogger{}, l)
}

// Logger recupera o logger do contexto, ou o padrão.
func Logger(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(chaveLogger{}).(*slog.Logger); ok {
		return l
	}
	return slog.Default()
}
