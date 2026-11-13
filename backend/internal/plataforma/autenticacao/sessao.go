package autenticacao

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"tenaz/internal/plataforma/autenticacao/internal/repo/gerado"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/erros"
)

// tecnico devolve IP e user agent da requisição já saneados para o banco: um
// IP malformado quebraria o cast para inet e derrubaria o login.
func tecnico(ctx context.Context) (ip, userAgent string) {
	r := db.RequisicaoDe(ctx)
	if a, err := netip.ParseAddr(r.IP); err == nil {
		ip = a.Unmap().String()
	}
	userAgent = r.UserAgent
	if len(userAgent) > 512 {
		userAgent = strings.ToValidUTF8(userAgent[:512], "")
	}
	return ip, userAgent
}

func ipDe(ctx context.Context) string {
	ip, _ := tecnico(ctx)
	return ip
}

// AutenticarWeb resolve a sessão web do cookie e devolve o Ator, com
// permissões e escopo lidos do banco nesta requisição.
//
// exigirCSRF: verdadeiro em POST, PATCH, PUT e DELETE; csrf é o cabeçalho
// X-CSRF-Token. A conferência de Origin é da camada HTTP.
func (s *Servico) AutenticarWeb(ctx context.Context, cookie, csrf string, exigirCSRF bool) (autorizacao.Ator, error) {
	sess, err := s.sessaoWebValida(ctx, cookie)
	if err != nil {
		return autorizacao.Ator{}, err
	}
	if exigirCSRF && (csrf == "" || !iguais(hashSegredo(csrf), sess.CsrfHash)) {
		return autorizacao.Ator{}, erros.CSRFInvalido
	}
	a, err := s.CarregarAtor(ctx, sess.EmpresaID, sess.UsuarioID)
	if err != nil {
		return autorizacao.Ator{}, err
	}
	if a.IdentidadeID != sess.IdentidadeID {
		return autorizacao.Ator{}, erros.CredencialInvalida
	}
	if err := s.marcarMFAPendente(ctx, &a); err != nil {
		return autorizacao.Ator{}, err
	}
	a.SessaoID, a.Canal = sess.ID, CanalWeb
	return a, nil
}

func (s *Servico) sessaoWebValida(ctx context.Context, cookie string) (gerado.AuthBuscarSessaoPorTokenRow, error) {
	if !segredoBemFormado(cookie) {
		return gerado.AuthBuscarSessaoPorTokenRow{}, erros.CredencialInvalida
	}
	agora := s.agora()
	var sess gerado.AuthBuscarSessaoPorTokenRow
	err := s.repo.Tx(ctx, func(q *gerado.Queries) error {
		var err error
		sess, err = q.AuthBuscarSessaoPorToken(ctx, hashSegredo(cookie))
		if err != nil {
			return naoEncontrado(err, erros.CredencialInvalida)
		}
		if sess.Canal != CanalWeb || sess.RevogadaEm != nil || !sess.ExpiraEm.After(agora) ||
			agora.Sub(sess.UltimoUsoEm) > s.cfg.InatividadeWeb {
			return erros.CredencialInvalida
		}
		// Inatividade conta a partir do último uso; regravar a cada
		// requisição seria uma escrita por chamada, então só depois de um
		// intervalo.
		if agora.Sub(sess.UltimoUsoEm) > s.cfg.IntervaloToqueWeb {
			return q.AuthTocarSessao(ctx, sess.ID)
		}
		return nil
	})
	return sess, err
}

// InfoSessao é a resposta de GET /v1/auth/sessao.
type InfoSessao struct {
	Ator     autorizacao.Ator
	CSRF     string
	Nome     string
	Email    string
	MFAAtivo bool
	// MFAPendente: sessão restrita à ativação do MFA (RN-ADM-027).
	MFAPendente bool
	Empresas    []EmpresaAcesso
}

// ObterSessaoWeb devolve a sessão, com o token CSRF (derivado do cookie) e as
// empresas da identidade para o seletor.
func (s *Servico) ObterSessaoWeb(ctx context.Context, cookie string) (InfoSessao, error) {
	a, err := s.AutenticarWeb(ctx, cookie, "", false)
	if err != nil {
		return InfoSessao{}, err
	}
	info, err := s.ObterSessao(ctx, a)
	info.CSRF = tokenCSRF(cookie)
	return info, err
}

// ObterSessao devolve os dados da sessão de um ator já autenticado (web ou
// app), sem o token CSRF, que só existe na sessão web.
func (s *Servico) ObterSessao(ctx context.Context, a autorizacao.Ator) (InfoSessao, error) {
	info := InfoSessao{Ator: a, MFAPendente: a.SomenteAtivarMFA}
	err := s.d.Banco.NaEmpresa(ctx, a.EmpresaID, a.Escopo, func(ctx context.Context) error {
		var err error
		info.Nome, err = autorizacao.NomeUsuario(ctx, a.EmpresaID, a.UsuarioID)
		return err
	})
	if err != nil {
		return InfoSessao{}, err
	}
	err = s.repo.Tx(ctx, func(q *gerado.Queries) error {
		id, err := q.AuthBuscarIdentidade(ctx, a.IdentidadeID)
		if err != nil {
			return err
		}
		info.Email, info.MFAAtivo = id.Email, id.MfaAtivadoEm != nil
		return nil
	})
	if err != nil {
		return InfoSessao{}, err
	}
	info.Empresas, err = s.ListarEmpresas(ctx, a.IdentidadeID)
	return info, err
}

// ListarEmpresas devolve as empresas habilitadas da identidade, com o nome
// lido na transação de cada uma.
func (s *Servico) ListarEmpresas(ctx context.Context, identidadeID uuid.UUID) ([]EmpresaAcesso, error) {
	lista, err := s.acessos(ctx, identidadeID)
	if err != nil {
		return nil, err
	}
	out := make([]EmpresaAcesso, 0, len(lista))
	for _, a := range lista {
		nome, err := s.nomeEmpresa(ctx, a.EmpresaID)
		if err != nil {
			return nil, err
		}
		out = append(out, EmpresaAcesso{EmpresaID: a.EmpresaID, Nome: nome, UsuarioID: a.UsuarioID})
	}
	return out, nil
}

// AutenticarApp confere o access token (Authorization: Bearer) e a sessão do
// aparelho, e devolve o Ator. Quem perdeu mob.usar_app perde o app na hora.
func (s *Servico) AutenticarApp(ctx context.Context, bearer string) (autorizacao.Ator, error) {
	c, err := s.d.Assinador.Verificar(strings.TrimSpace(bearer), s.agora())
	if err != nil {
		return autorizacao.Ator{}, erros.CredencialInvalida
	}
	usuarioID, err1 := uuid.Parse(c.Sujeito)
	empresaID, err2 := uuid.Parse(c.Empresa)
	sessaoID, err3 := uuid.Parse(c.Sessao)
	if err1 != nil || err2 != nil || err3 != nil {
		return autorizacao.Ator{}, erros.CredencialInvalida
	}
	agora := s.agora()
	var sess gerado.AuthBuscarSessaoRow
	err = s.repo.Tx(ctx, func(q *gerado.Queries) error {
		var err error
		sess, err = q.AuthBuscarSessao(ctx, sessaoID)
		if err != nil {
			return naoEncontrado(err, erros.CredencialInvalida)
		}
		if sess.Canal != CanalApp || sess.RevogadaEm != nil || !sess.ExpiraEm.After(agora) {
			return erros.CredencialInvalida
		}
		return nil
	})
	if err != nil {
		return autorizacao.Ator{}, err
	}
	// A empresa vem do token e é conferida contra a sessão (documento 3,
	// seção 4.7): um token nunca vale para outra empresa.
	if sess.EmpresaID != empresaID || sess.UsuarioID != usuarioID {
		return autorizacao.Ator{}, erros.CredencialInvalida
	}
	a, err := s.CarregarAtor(ctx, empresaID, usuarioID)
	if err != nil {
		return autorizacao.Ator{}, err
	}
	if a.IdentidadeID != sess.IdentidadeID {
		return autorizacao.Ator{}, erros.CredencialInvalida
	}
	if !a.Pode(autorizacao.MobUsarApp) {
		return autorizacao.Ator{}, erros.SemPermissao.ComChave("detalhe.perfil_sem_app", nil)
	}
	// Empresa que passou a exigir MFA: o app para até a ativação pela web.
	if err := s.marcarMFAPendente(ctx, &a); err != nil {
		return autorizacao.Ator{}, err
	}
	if a.SomenteAtivarMFA {
		return autorizacao.Ator{}, ErrMFAPendente()
	}
	a.SessaoID, a.Canal = sess.ID, CanalApp
	return a, nil
}

// Logout encerra a sessão do ator (web ou app; no app, com a família de
// refresh tokens) e grava o logout na trilha da empresa da sessão, na mesma
// transação. As outras empresas não recebem o logout: ele é de uma sessão,
// que é de uma empresa só.
func (s *Servico) Logout(ctx context.Context, a autorizacao.Ator) error {
	return s.repo.NaTrilha(ctx, a.EmpresaID, func(q *gerado.Queries) error {
		if _, err := q.AuthRevogarSessao(ctx, a.SessaoID); err != nil {
			return fmt.Errorf("revogar sessão: %w", err)
		}
		if err := q.AuthRevogarFamilias(ctx, []uuid.UUID{a.SessaoID}); err != nil {
			return err
		}
		return gravarEvento(ctx, q, a.EmpresaID, a.UsuarioID, eventoAuth{
			Acao: acaoLogout, IdentidadeID: a.IdentidadeID,
			Depois: map[string]string{"canal": a.Canal, "sessao_id": a.SessaoID.String()},
		})
	})
}

// ResultadoTroca traz a credencial nova da troca de empresa: sessão web ou
// tokens do app, conforme o canal da credencial de origem.
type ResultadoTroca struct {
	Web *SessaoWeb
	App *TokensApp
}

// TrocarEmpresa emite uma credencial NOVA na empresa pedida e revoga a atual
// (documento 3, seção 4.7): uma credencial nunca vale para duas empresas.
// Empresa fora do acesso da identidade: recurso_nao_encontrado. Sequência da seção 6.9
// do documento 5: acesso (autenticação), ator (aplicação), revogação, sessão
// nova e evento na trilha da empresa nova (autenticação, uma transação).
func (s *Servico) TrocarEmpresa(ctx context.Context, a autorizacao.Ator, empresaID uuid.UUID) (ResultadoTroca, error) {
	var usuarioID uuid.UUID
	err := s.repo.Tx(ctx, func(q *gerado.Queries) error {
		var err error
		usuarioID, err = q.AuthBuscarAcesso(ctx, gerado.AuthBuscarAcessoParams{IdentidadeID: a.IdentidadeID, EmpresaID: empresaID})
		return naoEncontrado(err, erros.NaoEncontrado)
	})
	if err != nil {
		return ResultadoTroca{}, err
	}
	novo, err := s.CarregarAtor(ctx, empresaID, usuarioID)
	if err != nil {
		return ResultadoTroca{}, err
	}
	if novo.IdentidadeID != a.IdentidadeID {
		return ResultadoTroca{}, erros.CredencialInvalida
	}
	if a.Canal == CanalApp && !novo.Pode(autorizacao.MobUsarApp) {
		return ResultadoTroca{}, erros.SemPermissao.ComChave("detalhe.perfil_sem_app_empresa", nil)
	}
	// Empresa de destino que exige MFA: na web, a sessão nova nasce restrita à
	// ativação; no app, a troca é recusada até a ativação (RN-ADM-027).
	if err := s.marcarMFAPendente(ctx, &novo); err != nil {
		return ResultadoTroca{}, err
	}
	if a.Canal == CanalApp && novo.SomenteAtivarMFA {
		return ResultadoTroca{}, ErrMFAPendente()
	}
	canal := a.Canal
	if canal == "" {
		canal = CanalWeb
	}

	var res ResultadoTroca
	err = s.repo.NaTrilha(ctx, empresaID, func(q *gerado.Queries) error {
		if _, err := q.AuthRevogarSessao(ctx, a.SessaoID); err != nil {
			return err
		}
		if err := q.AuthRevogarFamilias(ctx, []uuid.UUID{a.SessaoID}); err != nil {
			return err
		}
		var sessaoNova uuid.UUID
		if canal == CanalApp {
			t, err := s.criarSessaoApp(ctx, q, a.IdentidadeID, empresaID, usuarioID, nil)
			if err != nil {
				return err
			}
			res.App, sessaoNova = &t, t.SessaoID
		} else {
			w, err := s.criarSessaoWeb(ctx, q, a.IdentidadeID, empresaID, usuarioID)
			if err != nil {
				return err
			}
			w.MFAPendente = novo.SomenteAtivarMFA
			res.Web, sessaoNova = &w, w.SessaoID
		}
		return gravarEvento(ctx, q, empresaID, usuarioID, eventoLogin(a.IdentidadeID, canal, sessaoNova, "troca_de_empresa"))
	})
	return res, err
}

// EncerrarSessoesUsuario revoga as sessões e os refresh tokens do usuário na
// empresa do administrador (RF-ADM-004; POST /v1/usuarios/{id}/encerrar-
// sessoes) e grava o logout na trilha dessa empresa, na mesma transação. O
// evento é do usuário alvo (a política da trilha exige o usuário da
// identidade); quem pediu vai em depois.encerrado_por. Quem chama já
// conferiu a permissão e que o usuário é da empresa. Usuário que nunca
// aceitou o convite não tem sessão: devolve 0.
//
// A inativação não chama isto: o trigger de acesso_empresa revoga as
// credenciais na mesma transação (seção 6.9 do documento 5).
func (s *Servico) EncerrarSessoesUsuario(ctx context.Context, admin autorizacao.Ator, usuarioID uuid.UUID) (int, error) {
	var n int
	err := s.repo.NaTrilha(ctx, admin.EmpresaID, func(q *gerado.Queries) error {
		identidadeID, err := q.AuthIdentidadeDoUsuario(ctx, gerado.AuthIdentidadeDoUsuarioParams{
			EmpresaID: admin.EmpresaID, UsuarioID: usuarioID,
		})
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		if err != nil {
			return err
		}
		ids, err := q.AuthRevogarSessoesUsuario(ctx, gerado.AuthRevogarSessoesUsuarioParams{EmpresaID: admin.EmpresaID, UsuarioID: usuarioID})
		if err != nil {
			return fmt.Errorf("revogar sessões do usuário: %w", err)
		}
		n = len(ids)
		if n > 0 {
			if err := q.AuthRevogarFamilias(ctx, ids); err != nil {
				return err
			}
		}
		return gravarEvento(ctx, q, admin.EmpresaID, usuarioID, eventoAuth{
			Acao: acaoLogout, IdentidadeID: identidadeID,
			Depois: map[string]any{"motivo": "encerradas_pelo_administrador", "encerrado_por": admin.UsuarioID, "sessoes": n},
		})
	})
	return n, err
}

// revogarSessoesIdentidade revoga as sessões da identidade em todas as
// empresas, menos a indicada (troca e redefinição de senha).
func revogarSessoesIdentidade(ctx context.Context, q *gerado.Queries, identidadeID uuid.UUID, exceto *uuid.UUID) error {
	ids, err := q.AuthRevogarSessoesIdentidade(ctx, gerado.AuthRevogarSessoesIdentidadeParams{
		IdentidadeID: identidadeID, ExcetoSessaoID: exceto,
	})
	if err != nil {
		return fmt.Errorf("revogar sessões da identidade: %w", err)
	}
	if len(ids) == 0 {
		return nil
	}
	return q.AuthRevogarFamilias(ctx, ids)
}
