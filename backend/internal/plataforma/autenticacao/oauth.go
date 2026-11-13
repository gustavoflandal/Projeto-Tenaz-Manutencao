package autenticacao

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"tenaz/internal/plataforma/autenticacao/internal/repo/gerado"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/erros"
	"tenaz/internal/plataforma/obs"
)

// OAuth 2.1 do app (documento 3, seção 4.3). O app é cliente público: sem
// segredo, com PKCE S256 obrigatório.
//
// A tela de autorização é servida pela web: o usuário entra com a sessão web
// normal (senha e MFA), escolhe a empresa, e a web chama
// EmitirCodigoAutorizacao com o Ator da sessão. Assim o app nunca vê a senha
// e o login do app usa exatamente as mesmas regras do login da web.

// PrefixoRefreshToken identifica o refresh token, para scanners de segredo.
const PrefixoRefreshToken = "tnz_rt_"

var reDesafioPKCE = regexp.MustCompile(`^[A-Za-z0-9_-]{43}$`)

// EntradaAutorizacao são os parâmetros de GET /oauth/authorize, já com a
// empresa escolhida na tela.
type EntradaAutorizacao struct {
	ClienteID   string // client_id
	URIRetorno  string // redirect_uri
	DesafioPKCE string // code_challenge
	MetodoPKCE  string // code_challenge_method; só S256
	EmpresaID   uuid.UUID
}

// ValidarCliente confere client_id e redirect_uri (byte a byte) antes de a
// tela de autorização abrir. Divergência nunca redireciona: a URI não é
// confiável (OAuth 2.1, seção 4.1.2.1).
func (s *Servico) ValidarCliente(clienteID, uriRetorno string) error {
	if clienteID != s.cfg.ClienteOAuthApp {
		return erros.EntradaInvalida.ComCampo("client_id", erros.CampoValorInvalido, "Cliente OAuth desconhecido.")
	}
	if !slices.Contains(s.cfg.URIsRetornoApp, uriRetorno) {
		return erros.EntradaInvalida.ComCampo("redirect_uri", erros.CampoValorInvalido, "URI de retorno não registrada.")
	}
	return nil
}

// EmitirCodigoAutorizacao emite o código de uso único (60 s) para a empresa
// escolhida. O Ator é o da sessão web da tela de autorização; a empresa pode
// ser qualquer uma da identidade, conferida contra acesso_empresa. Perfil sem
// mob.usar_app (o locatario, por exemplo) não recebe código.
func (s *Servico) EmitirCodigoAutorizacao(ctx context.Context, sessaoWeb autorizacao.Ator, e EntradaAutorizacao) (string, error) {
	if err := s.ValidarCliente(e.ClienteID, e.URIRetorno); err != nil {
		return "", err
	}
	if e.MetodoPKCE != "S256" {
		return "", erros.EntradaInvalida.ComCampo("code_challenge_method", erros.CampoValorInvalido, "Só S256 é aceito.")
	}
	if !reDesafioPKCE.MatchString(e.DesafioPKCE) {
		return "", erros.EntradaInvalida.ComCampo("code_challenge", erros.CampoFormatoInvalido,
			"code_challenge deve ser o sha256 do verificador em base64url, com 43 caracteres.")
	}

	var usuarioID uuid.UUID
	err := s.repo.Tx(ctx, func(q *gerado.Queries) error {
		var err error
		usuarioID, err = q.AuthBuscarAcesso(ctx, gerado.AuthBuscarAcessoParams{IdentidadeID: sessaoWeb.IdentidadeID, EmpresaID: e.EmpresaID})
		return naoEncontrado(err, erros.NaoEncontrado)
	})
	if err != nil {
		return "", err
	}
	a, err := s.CarregarAtor(ctx, e.EmpresaID, usuarioID)
	if err != nil {
		return "", err
	}
	if !a.Pode(autorizacao.MobUsarApp) {
		return "", erros.SemPermissao.Com("o perfil não usa o app (mob.usar_app)")
	}
	if err := s.marcarMFAPendente(ctx, &a); err != nil {
		return "", err
	}
	if a.SomenteAtivarMFA {
		return "", ErrMFAPendente()
	}

	codigo, hash, err := novoSegredo()
	if err != nil {
		return "", err
	}
	agora := s.agora()
	err = s.repo.Tx(ctx, func(q *gerado.Queries) error {
		return q.AuthCriarCodigoAutorizacao(ctx, gerado.AuthCriarCodigoAutorizacaoParams{
			ID: uuid.Must(uuid.NewV7()), CodigoHash: hash, IdentidadeID: sessaoWeb.IdentidadeID,
			EmpresaID: e.EmpresaID, UsuarioID: usuarioID, ClienteOauth: e.ClienteID, UriRetorno: e.URIRetorno,
			DesafioPkce: e.DesafioPKCE, CriadoEm: agora, ExpiraEm: agora.Add(s.cfg.ValidadeCodigoOAuth),
		})
	})
	if err != nil {
		return "", fmt.Errorf("criar código de autorização: %w", err)
	}
	return codigo, nil
}

// EntradaTrocaCodigo é POST /oauth/token com grant_type=authorization_code.
type EntradaTrocaCodigo struct {
	Codigo          string
	ClienteID       string
	URIRetorno      string
	VerificadorPKCE string
	DispositivoID   *uuid.UUID
}

// TrocarCodigo troca o código pelos tokens. Qualquer divergência (cliente,
// URI de retorno byte a byte, PKCE, validade) responde credencial_invalida, e
// o código fica gasto. Código reapresentado depois de usado revoga a sessão
// que ele gerou (OAuth 2.1, seção 4.1.3).
func (s *Servico) TrocarCodigo(ctx context.Context, e EntradaTrocaCodigo) (TokensApp, error) {
	if !segredoBemFormado(e.Codigo) {
		return TokensApp{}, erros.CredencialInvalida
	}
	hash := hashSegredo(e.Codigo)
	var cod gerado.AuthUsarCodigoAutorizacaoRow
	reuso := false
	err := s.repo.Tx(ctx, func(q *gerado.Queries) error {
		var err error
		cod, err = q.AuthUsarCodigoAutorizacao(ctx, hash)
		if !errors.Is(err, pgx.ErrNoRows) {
			return err
		}
		antigo, err := q.AuthBuscarCodigoAutorizacao(ctx, hash)
		if err != nil {
			return naoEncontrado(err, erros.CredencialInvalida)
		}
		if antigo.UsadoEm != nil && antigo.SessaoID != nil {
			reuso = true
			if _, err := q.AuthRevogarSessao(ctx, *antigo.SessaoID); err != nil {
				return err
			}
			return q.AuthRevogarFamilias(ctx, []uuid.UUID{*antigo.SessaoID})
		}
		return nil
	})
	if err != nil {
		return TokensApp{}, err
	}
	if cod.ID == uuid.Nil {
		if reuso {
			obs.Logger(ctx).Warn("código de autorização reapresentado; sessão revogada", "evento", "reuso_codigo_oauth")
		}
		return TokensApp{}, erros.CredencialInvalida
	}
	if cod.ClienteOauth != e.ClienteID || cod.UriRetorno != e.URIRetorno || !conferePKCE(e.VerificadorPKCE, cod.DesafioPkce) {
		return TokensApp{}, erros.CredencialInvalida
	}
	a, err := s.CarregarAtor(ctx, cod.EmpresaID, cod.UsuarioID)
	if err != nil {
		return TokensApp{}, err
	}
	if a.IdentidadeID != cod.IdentidadeID || !a.Pode(autorizacao.MobUsarApp) {
		return TokensApp{}, erros.CredencialInvalida
	}
	if err := s.marcarMFAPendente(ctx, &a); err != nil {
		return TokensApp{}, err
	}
	if a.SomenteAtivarMFA {
		return TokensApp{}, ErrMFAPendente()
	}

	var t TokensApp
	err = s.repo.NaTrilha(ctx, cod.EmpresaID, func(q *gerado.Queries) error {
		var err error
		t, err = s.criarSessaoApp(ctx, q, cod.IdentidadeID, cod.EmpresaID, cod.UsuarioID, e.DispositivoID)
		if err != nil {
			return err
		}
		if err := q.AuthLigarSessaoAoCodigo(ctx, gerado.AuthLigarSessaoAoCodigoParams{ID: cod.ID, SessaoID: &t.SessaoID}); err != nil {
			return err
		}
		return gravarEvento(ctx, q, cod.EmpresaID, cod.UsuarioID, eventoLogin(cod.IdentidadeID, CanalApp, t.SessaoID, ""))
	})
	if err != nil {
		return TokensApp{}, err
	}
	return t, nil
}

// criarSessaoApp cria a sessão do aparelho, o primeiro refresh token da
// família e o access token, na transação da autenticação.
func (s *Servico) criarSessaoApp(ctx context.Context, q *gerado.Queries, identidadeID, empresaID, usuarioID uuid.UUID, dispositivoID *uuid.UUID) (TokensApp, error) {
	agora := s.agora()
	sessaoID := uuid.Must(uuid.NewV7())
	ip, ua := tecnico(ctx)
	err := q.AuthCriarSessao(ctx, gerado.AuthCriarSessaoParams{
		ID: sessaoID, Canal: CanalApp, IdentidadeID: identidadeID, EmpresaID: empresaID, UsuarioID: usuarioID,
		DispositivoID: dispositivoID, Ip: ip, UserAgent: ua, CriadoEm: agora, ExpiraEm: agora.Add(s.cfg.ValidadeApp),
	})
	if err != nil {
		return TokensApp{}, fmt.Errorf("criar sessão do app: %w", err)
	}
	refresh, err := s.novoRefresh(ctx, q, sessaoID, nil, agora)
	if err != nil {
		return TokensApp{}, err
	}
	return s.montarTokens(usuarioID, empresaID, sessaoID, refresh, agora)
}

func (s *Servico) novoRefresh(ctx context.Context, q *gerado.Queries, familia uuid.UUID, anterior *uuid.UUID, agora time.Time) (string, error) {
	texto, hash, err := novoSegredo()
	if err != nil {
		return "", err
	}
	err = q.AuthCriarTokenRenovacao(ctx, gerado.AuthCriarTokenRenovacaoParams{
		ID: uuid.Must(uuid.NewV7()), FamiliaID: familia, AnteriorID: anterior, TokenHash: hash,
		CriadoEm: agora, ExpiraEm: agora.Add(s.cfg.ValidadeApp),
	})
	if err != nil {
		return "", fmt.Errorf("criar refresh token: %w", err)
	}
	return PrefixoRefreshToken + texto, nil
}

func (s *Servico) montarTokens(usuarioID, empresaID, sessaoID uuid.UUID, refresh string, agora time.Time) (TokensApp, error) {
	access, exp, err := s.d.Assinador.Emitir(usuarioID, empresaID, sessaoID, agora)
	if err != nil {
		return TokensApp{}, fmt.Errorf("assinar access token: %w", err)
	}
	return TokensApp{
		AccessToken: access, ExpiraEm: exp, ExpiraEmSegs: int64(exp.Sub(agora) / time.Second), RefreshToken: refresh,
		SessaoID: sessaoID, EmpresaID: empresaID, UsuarioID: usuarioID,
	}, nil
}

func segredoDoRefresh(token string) (string, bool) {
	texto, ok := strings.CutPrefix(strings.TrimSpace(token), PrefixoRefreshToken)
	return texto, ok && segredoBemFormado(texto)
}

// Renovar troca o refresh token por um novo e devolve um access token novo
// (rotação, documento 3, seção 4.3). Refresh token já usado e reapresentado é
// sinal de roubo: a família inteira e a sessão do aparelho são revogadas, e o
// usuário entra de novo. O app precisa serializar as renovações: duas
// simultâneas com o mesmo token contam como reuso (pendência E01-7).
func (s *Servico) Renovar(ctx context.Context, refreshToken, clienteID string) (TokensApp, error) {
	texto, ok := segredoDoRefresh(refreshToken)
	if !ok || clienteID != s.cfg.ClienteOAuthApp {
		return TokensApp{}, erros.CredencialInvalida
	}
	hash := hashSegredo(texto)
	agora := s.agora()

	var (
		sess    gerado.AuthBuscarSessaoRow
		refresh string
		reuso   bool
	)
	err := s.repo.Tx(ctx, func(q *gerado.Queries) error {
		usado, err := q.AuthUsarTokenRenovacao(ctx, hash)
		if errors.Is(err, pgx.ErrNoRows) {
			antigo, err := q.AuthBuscarTokenRenovacao(ctx, hash)
			if err != nil {
				return naoEncontrado(err, erros.CredencialInvalida)
			}
			if antigo.UsadoEm != nil && antigo.RevogadoEm == nil {
				// Reuso: revoga tudo e confirma a revogação antes de recusar.
				reuso = true
				if err := q.AuthRevogarFamilias(ctx, []uuid.UUID{antigo.FamiliaID}); err != nil {
					return err
				}
				_, err := q.AuthRevogarSessao(ctx, antigo.FamiliaID)
				return err
			}
			return nil
		}
		if err != nil {
			return err
		}
		sess, err = q.AuthBuscarSessao(ctx, usado.FamiliaID)
		if err != nil {
			return err
		}
		if sess.Canal != CanalApp || sess.RevogadaEm != nil || !sess.ExpiraEm.After(agora) {
			return erros.CredencialInvalida
		}
		refresh, err = s.novoRefresh(ctx, q, sess.ID, &usado.ID, agora)
		if err != nil {
			return err
		}
		return q.AuthEstenderSessao(ctx, gerado.AuthEstenderSessaoParams{ID: sess.ID, ExpiraEm: agora.Add(s.cfg.ValidadeApp)})
	})
	if err != nil {
		return TokensApp{}, err
	}
	if refresh == "" {
		if reuso {
			obs.Logger(ctx).Warn("refresh token reapresentado; família revogada", "evento", "reuso_refresh_token")
		}
		return TokensApp{}, erros.CredencialInvalida
	}

	// O usuário pode ter sido inativado ou perdido mob.usar_app desde a
	// última renovação: confere agora e, se for o caso, revoga.
	a, err := s.CarregarAtor(ctx, sess.EmpresaID, sess.UsuarioID)
	if err == nil && (a.IdentidadeID != sess.IdentidadeID || !a.Pode(autorizacao.MobUsarApp)) {
		err = erros.CredencialInvalida
	}
	if err == nil {
		err = s.marcarMFAPendente(ctx, &a)
	}
	if err == nil && a.SomenteAtivarMFA {
		// O refresh token já girou e o novo não será entregue: a sessão do
		// aparelho é revogada, e o app entra de novo depois de ativar o MFA.
		_ = s.revogarSessao(ctx, sess.ID)
		return TokensApp{}, ErrMFAPendente()
	}
	if err != nil {
		if errors.Is(err, erros.CredencialInvalida) {
			_ = s.revogarSessao(ctx, sess.ID)
		}
		return TokensApp{}, err
	}
	return s.montarTokens(sess.UsuarioID, sess.EmpresaID, sess.ID, refresh, agora)
}

// Revogar é POST /oauth/revoke (RFC 7009): revoga a família e a sessão do
// aparelho. Token desconhecido não é erro, para não servir de oráculo.
func (s *Servico) Revogar(ctx context.Context, refreshToken string) error {
	texto, ok := segredoDoRefresh(refreshToken)
	if !ok {
		return nil
	}
	return s.repo.Tx(ctx, func(q *gerado.Queries) error {
		t, err := q.AuthBuscarTokenRenovacao(ctx, hashSegredo(texto))
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		if err != nil {
			return err
		}
		if err := q.AuthRevogarFamilias(ctx, []uuid.UUID{t.FamiliaID}); err != nil {
			return err
		}
		_, err = q.AuthRevogarSessao(ctx, t.FamiliaID)
		return err
	})
}
