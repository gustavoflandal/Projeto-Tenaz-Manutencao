package db

import (
	"context"
	"fmt"
	"net/netip"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
)

// Atores da trilha de auditoria (ck_evento_auditoria_ator, seção 9.2 do
// documento 2).
const (
	AtorUsuario              = "usuario"
	AtorSistema              = "sistema"
	AtorIntegracao           = "integracao"
	AtorSuporteTenaz         = "suporte_tenaz"
	AtorLinkExecucao         = "link_execucao"
	AtorSincronizacaoConexao = "sincronizacao_conexao"
)

// Origens da trilha (ck_evento_auditoria_origem).
const (
	OrigemWeb        = "web"
	OrigemApp        = "app"
	OrigemAPI        = "api"
	OrigemJob        = "job"
	OrigemImportacao = "importacao"
)

// Requisicao é o contexto opcional da auditoria (seção 2.4 do documento 2),
// gravado por app.tg_auditar() em todo evento da transação. O middleware HTTP
// preenche uma vez por requisição, com ComRequisicao; NaEmpresa lê daqui.
//
// Nada aqui decide empresa nem escopo: isso vem só da credencial, pelo Escopo.
type Requisicao struct {
	AtorTipo                string // vazio: "usuario" com usuário no escopo, "sistema" sem
	Origem                  string // web, app, api, job, importacao; vazio grava nulo
	IP                      string
	UserAgent               string
	RequestID               string
	TokenAPIID              *uuid.UUID
	DispositivoID           *uuid.UUID
	OcorridoNoDispositivoEm *time.Time // instante do aparelho, em comando feito offline

	// Modulo é o código do módulo que atende a chamada (sumário, seção 2, em
	// minúsculas): vai para evento_auditoria.modulo. Quem monta o serviço do
	// módulo o fixa; vazio grava nulo.
	Modulo string
	// Tela é a tela de origem, '<modulo>.<tela>', já validada contra o
	// catálogo fechado (X-Tenaz-Tela). Só vale quando começa pelo módulo da
	// chamada, como exige ck_evento_auditoria_tela; senão é descartada. Vazia
	// em API, job e sincronização.
	Tela string
}

// TelaDoModulo devolve a tela só se ela é do módulo da chamada.
func (r Requisicao) TelaDoModulo() string { return telaCoerente(strings.ToLower(r.Modulo), r.Tela) }

// Módulos aceitos em evento_auditoria.modulo (ck_evento_auditoria_modulo).
var modulosValidos = map[string]struct{}{
	"adm": {}, "atv": {}, "med": {}, "sol": {}, "os": {}, "pre": {}, "chk": {}, "est": {}, "ter": {}, "srv": {},
	"prj": {}, "ctr": {}, "loc": {}, "fro": {}, "con": {}, "ind": {}, "int": {}, "ia": {}, "mob": {},
}

// ModuloValido informa se o código é um dos 19 módulos do sumário.
func ModuloValido(m string) bool {
	_, ok := modulosValidos[m]
	return ok
}

// ComModulo devolve o contexto com o módulo da chamada fixado, preservando o
// resto da requisição.
func ComModulo(ctx context.Context, modulo string) context.Context {
	r := RequisicaoDe(ctx)
	if r.Modulo == modulo {
		return ctx
	}
	r.Modulo = modulo
	return ComRequisicao(ctx, r)
}

type chaveRequisicao struct{}

// ComRequisicao guarda os dados da requisição no contexto.
func ComRequisicao(ctx context.Context, r Requisicao) context.Context {
	return context.WithValue(ctx, chaveRequisicao{}, r)
}

// RequisicaoDe devolve os dados da requisição; vazios, se não houver.
func RequisicaoDe(ctx context.Context) Requisicao {
	r, _ := ctx.Value(chaveRequisicao{}).(Requisicao)
	return r
}

const (
	maxUserAgent = 512
	maxRequestID = 128
)

// normalizada confere os valores antes de irem para o banco: um valor fora
// do CHECK da trilha faria a transação inteira falhar no primeiro evento, e
// um IP malformado quebraria o cast para inet. Campo técnico inválido vira
// vazio; ator ou origem inválidos são erro de programação.
func (r Requisicao) normalizada(escopo Escopo) (Requisicao, error) {
	if r.AtorTipo == "" {
		switch {
		case escopo.Tipo == EscopoLink:
			r.AtorTipo = AtorLinkExecucao
		case escopo.UsuarioID != nil:
			r.AtorTipo = AtorUsuario
		default:
			r.AtorTipo = AtorSistema
		}
	}
	switch r.AtorTipo {
	case AtorUsuario, AtorSistema, AtorIntegracao, AtorSuporteTenaz, AtorLinkExecucao:
	case AtorSincronizacaoConexao:
		// Só o pool do sincronizador grava com este ator, e ele chega com a
		// conexão (E04). Aqui é sempre erro.
		return r, fmt.Errorf("ator %q exige o pool do sincronizador", r.AtorTipo)
	default:
		return r, fmt.Errorf("ator inválido: %q", r.AtorTipo)
	}
	switch r.Origem {
	case "", OrigemWeb, OrigemApp, OrigemAPI, OrigemJob, OrigemImportacao:
	default:
		return r, fmt.Errorf("origem inválida: %q", r.Origem)
	}
	if r.IP != "" {
		if a, err := netip.ParseAddr(r.IP); err != nil {
			r.IP = ""
		} else {
			r.IP = a.Unmap().String()
		}
	}
	r.UserAgent = cortar(r.UserAgent, maxUserAgent)
	r.RequestID = cortar(r.RequestID, maxRequestID)
	r.Modulo = strings.ToLower(r.Modulo)
	if r.Modulo != "" && !ModuloValido(r.Modulo) {
		return r, fmt.Errorf("módulo inválido: %q", r.Modulo)
	}
	r.Tela = telaCoerente(r.Modulo, r.Tela)
	return r, nil
}

// telaCoerente descarta a tela que não é do módulo da chamada: o CHECK da
// trilha exige que ela comece pelo módulo, e uma tela de outro módulo (a
// listagem de usuários chamada pela tela do painel, por exemplo) faria a
// transação inteira falhar no primeiro evento.
func telaCoerente(modulo, tela string) string {
	if modulo == "" || !strings.HasPrefix(tela, modulo+".") {
		return ""
	}
	return tela
}

// Saneada devolve os campos técnicos prontos para o banco, sem decidir ator:
// IP inválido vira vazio, user agent e request_id são cortados e origem fora
// da lista vira vazio. É o que o repositório da autenticação usa na
// transação da trilha, que não tem escopo nem usuário (seção 6.9 do
// documento 5).
func (r Requisicao) Saneada() Requisicao {
	switch r.Origem {
	case OrigemWeb, OrigemApp, OrigemAPI, OrigemJob, OrigemImportacao:
	default:
		r.Origem = ""
	}
	if r.IP != "" {
		if a, err := netip.ParseAddr(r.IP); err != nil {
			r.IP = ""
		} else {
			r.IP = a.Unmap().String()
		}
	}
	r.UserAgent = cortar(r.UserAgent, maxUserAgent)
	r.RequestID = cortar(r.RequestID, maxRequestID)
	// O evento da identidade é sempre do módulo adm (RF-ADM-006).
	r.Modulo = "adm"
	r.Tela = telaCoerente(r.Modulo, r.Tela)
	return r
}

// cortar limita o tamanho em bytes sem partir um caractere UTF-8 ao meio.
func cortar(s string, max int) string {
	if len(s) <= max {
		return s
	}
	s = s[:max]
	for !utf8.ValidString(s) {
		s = s[:len(s)-1]
	}
	return s
}

func textoInstante(t *time.Time) string {
	if t == nil || t.IsZero() {
		return ""
	}
	return t.UTC().Format(time.RFC3339Nano)
}
