package autenticacao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"tenaz/internal/plataforma/autenticacao/internal/repo/gerado"
	"tenaz/internal/plataforma/erros"
	"tenaz/internal/plataforma/obs"
)

// EntradaLogin é o corpo de POST /v1/auth/login.
type EntradaLogin struct {
	Email string
	Senha string
	// EmpresaID escolhe a empresa da sessão entre as da identidade. Vazio: a
	// da última sessão, ou a primeira. A escolha é conferida contra
	// acesso_empresa; uma empresa fora dele responde recurso_nao_encontrado.
	EmpresaID *uuid.UUID
}

// Motivos gravados em login_falho (campo depois.motivo).
const (
	motivoSenhaIncorreta = "senha_incorreta"
	motivoContaBloqueada = "conta_bloqueada"
	motivoMFAIncorreto   = "mfa_incorreto"
)

// Login confere e-mail e senha.
//
// Toda falha responde credencial_invalida, no mesmo tempo: e-mail que não
// existe, senha errada, conta bloqueada (temporária ou administrativamente),
// identidade só de SSO e identidade sem empresa habilitada. A diferença fica
// no log de segurança e na trilha, nunca na resposta (documento 9, seção 5.3).
//
// Senha certa com MFA ativo devolve o desafio; sem MFA, a sessão web.
func (s *Servico) Login(ctx context.Context, e EntradaLogin) (ResultadoLogin, error) {
	log := obs.Logger(ctx)
	emailNorm := normalizarEmail(e.Email)
	if emailNorm == "" || e.Senha == "" || len(e.Senha) > 4*SenhaMaximo {
		s.d.Hasher.GastarTempo(ctx, e.Senha)
		return ResultadoLogin{}, erros.CredencialInvalida
	}

	var id gerado.AuthBuscarIdentidadePorEmailRow
	err := s.repo.Tx(ctx, func(q *gerado.Queries) error {
		var err error
		id, err = q.AuthBuscarIdentidadePorEmail(ctx, emailNorm)
		return err
	})
	if errors.Is(err, pgx.ErrNoRows) {
		// Sem empresa: vai só para o log de segurança, nunca para a trilha
		// (seção 9.1 do documento 2). O e-mail não vai para o log.
		s.d.Hasher.GastarTempo(ctx, e.Senha)
		log.Warn("login com e-mail inexistente", "evento", "login_falho")
		return ResultadoLogin{}, erros.CredencialInvalida
	}
	if err != nil {
		return ResultadoLogin{}, fmt.Errorf("buscar identidade: %w", err)
	}

	agora := s.agora()
	if id.BloqueadaEm != nil || (id.BloqueadaAte != nil && id.BloqueadaAte.After(agora)) || id.SenhaHash == nil {
		s.d.Hasher.GastarTempo(ctx, e.Senha)
		if id.SenhaHash != nil {
			_ = s.eventoNasEmpresas(ctx, nil, eventoLoginFalho(id.ID, motivoContaBloqueada))
		}
		log.Warn("login recusado: conta bloqueada ou sem senha", "evento", "login_falho", "identidade_id", id.ID)
		return ResultadoLogin{}, erros.CredencialInvalida
	}

	ok, rehash, err := s.d.Hasher.Verificar(ctx, e.Senha, *id.SenhaHash)
	if err != nil {
		return ResultadoLogin{}, fmt.Errorf("conferir senha: %w", err)
	}
	if !ok {
		if err := s.registrarFalha(ctx, id.ID, agora); err != nil {
			return ResultadoLogin{}, err
		}
		_ = s.eventoNasEmpresas(ctx, nil, eventoLoginFalho(id.ID, motivoSenhaIncorreta))
		log.Warn("login com senha incorreta", "evento", "login_falho", "identidade_id", id.ID)
		return ResultadoLogin{}, erros.CredencialInvalida
	}

	if rehash {
		// Parâmetros do Argon2id mudaram: regrava com os atuais agora, que é o
		// único momento em que a senha em claro está disponível.
		novo, err := s.d.Hasher.Gerar(ctx, e.Senha)
		if err == nil {
			err = s.repo.Tx(ctx, func(q *gerado.Queries) error {
				return q.AuthDefinirSenha(ctx, gerado.AuthDefinirSenhaParams{ID: id.ID, SenhaHash: &novo})
			})
		}
		if err != nil {
			log.Error("regravar hash de senha", "identidade_id", id.ID, "erro", err)
		}
	}

	if id.MfaAtivadoEm != nil {
		return s.criarDesafioMFA(ctx, id.ID, agora)
	}
	return s.concluirLogin(ctx, id.ID, e.EmpresaID, false)
}

func eventoLoginFalho(identidadeID uuid.UUID, motivo string) eventoAuth {
	return eventoAuth{Acao: acaoLoginFalho, IdentidadeID: identidadeID, Depois: map[string]string{"motivo": motivo}}
}

// registrarFalha incrementa as falhas e, a cada 5 seguidas, bloqueia por um
// tempo que dobra a cada bloqueio: 15 min, 30 min, 1 h... até 24 h. O login
// certo zera a contagem (bloqueio progressivo, documento 9, seção 5.3).
func (s *Servico) registrarFalha(ctx context.Context, identidadeID uuid.UUID, agora time.Time) error {
	return s.repo.Tx(ctx, func(q *gerado.Queries) error {
		falhas, err := q.AuthIncrementarFalhasLogin(ctx, identidadeID)
		if err != nil {
			return fmt.Errorf("registrar falha de login: %w", err)
		}
		if bloqueio := s.duracaoBloqueio(int(falhas)); bloqueio > 0 {
			return q.AuthBloquearIdentidadeAte(ctx, gerado.AuthBloquearIdentidadeAteParams{
				ID: identidadeID, BloqueadaAte: ptr(agora.Add(bloqueio)),
			})
		}
		return nil
	})
}

func (s *Servico) duracaoBloqueio(falhas int) time.Duration {
	n := s.cfg.FalhasParaBloquear
	if n < 1 || falhas < n || falhas%n != 0 {
		return 0
	}
	d := s.cfg.BloqueioInicial
	for i := 1; i < falhas/n && d < s.cfg.BloqueioMaximo; i++ {
		d *= 2
	}
	return min(d, s.cfg.BloqueioMaximo)
}

func (s *Servico) criarDesafioMFA(ctx context.Context, identidadeID uuid.UUID, agora time.Time) (ResultadoLogin, error) {
	texto, hash, err := novoSegredo()
	if err != nil {
		return ResultadoLogin{}, err
	}
	expira := agora.Add(s.cfg.ValidadeDesafioMFA)
	err = s.repo.Tx(ctx, func(q *gerado.Queries) error {
		return q.AuthCriarDesafioMFA(ctx, gerado.AuthCriarDesafioMFAParams{
			ID: uuid.Must(uuid.NewV7()), TokenHash: hash, IdentidadeID: identidadeID, Canal: CanalWeb,
			Ip: ipDe(ctx), CriadoEm: agora, ExpiraEm: expira,
		})
	})
	if err != nil {
		return ResultadoLogin{}, fmt.Errorf("criar desafio de MFA: %w", err)
	}
	return ResultadoLogin{DesafioMFA: &DesafioMFA{DesafioID: texto, ExpiraEm: expira}}, nil
}

// EntradaVerificarMFA é o corpo de POST /v1/auth/mfa/verificar.
type EntradaVerificarMFA struct {
	DesafioID string
	Codigo    string
	EmpresaID *uuid.UUID
}

// VerificarMFA conclui o login com o código TOTP. Cada código errado consome
// uma das 5 tentativas do desafio e soma uma falha ao bloqueio progressivo da
// identidade (RN-ADM-027); esgotadas as tentativas, ou vencido o desafio, o
// login recomeça pela senha. Código de um passo TOTP já usado vale como
// errado. Código errado: mfa_codigo_invalido, com `tentativas_restantes` em
// Dados. Desafio inválido ou identidade bloqueada: credencial_invalida, sem
// somar falha.
func (s *Servico) VerificarMFA(ctx context.Context, e EntradaVerificarMFA) (ResultadoLogin, error) {
	if !segredoBemFormado(e.DesafioID) {
		return ResultadoLogin{}, erros.CredencialInvalida
	}
	agora := s.agora()
	var (
		identidadeID uuid.UUID
		restantes    int
		certo        bool
	)
	err := s.repo.Tx(ctx, func(q *gerado.Queries) error {
		d, err := q.AuthTravarDesafioMFA(ctx, hashSegredo(e.DesafioID))
		if err != nil {
			return naoEncontrado(err, erros.CredencialInvalida)
		}
		if d.ConcluidoEm != nil || !d.ExpiraEm.After(agora) || d.Tentativas >= 5 {
			return erros.CredencialInvalida
		}
		id, err := q.AuthTravarIdentidade(ctx, d.IdentidadeID)
		if err != nil {
			return err
		}
		if id.BloqueadaEm != nil || (id.BloqueadaAte != nil && id.BloqueadaAte.After(agora)) || id.MfaAtivadoEm == nil {
			return erros.CredencialInvalida
		}
		segredo, err := s.d.Cifrador.Decifrar(id.MfaSegredoCifrado, id.ID[:])
		if err != nil {
			return fmt.Errorf("decifrar segredo TOTP: %w", err)
		}
		identidadeID = id.ID
		certo, err = consumirTOTP(ctx, q, id.ID, segredo, e.Codigo, agora)
		if err != nil {
			return err
		}
		if certo {
			return q.AuthConcluirDesafioMFA(ctx, d.ID)
		}
		restantes = 5 - int(d.Tentativas) - 1
		// O incremento precisa ser confirmado mesmo com o código errado: por
		// isso a função devolve nil e o erro sai depois do COMMIT.
		return q.AuthIncrementarTentativaMFA(ctx, d.ID)
	})
	if err != nil {
		return ResultadoLogin{}, err
	}
	if !certo {
		if err := s.registrarFalha(ctx, identidadeID, agora); err != nil {
			return ResultadoLogin{}, err
		}
		_ = s.eventoNasEmpresas(ctx, nil, eventoLoginFalho(identidadeID, motivoMFAIncorreto))
		return ResultadoLogin{}, erros.MFACodigoInvalido.ComDado("tentativas_restantes", max(restantes, 0))
	}
	return s.concluirLogin(ctx, identidadeID, e.EmpresaID, true)
}

// escolherEmpresa decide a empresa da sessão.
func (s *Servico) escolherEmpresa(ctx context.Context, identidadeID uuid.UUID, pedida *uuid.UUID) (uuid.UUID, uuid.UUID, []gerado.AuthListarAcessosRow, error) {
	var (
		lista  []gerado.AuthListarAcessosRow
		ultima uuid.UUID
	)
	err := s.repo.Tx(ctx, func(q *gerado.Queries) error {
		var err error
		lista, err = q.AuthListarAcessos(ctx, identidadeID)
		if err != nil {
			return err
		}
		ultima, err = q.AuthUltimaEmpresa(ctx, identidadeID)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		return err
	})
	if err != nil {
		return uuid.Nil, uuid.Nil, nil, err
	}
	if len(lista) == 0 {
		// Nenhum usuário ativo em empresa nenhuma: para quem tenta entrar, é
		// igual a senha errada.
		return uuid.Nil, uuid.Nil, nil, erros.CredencialInvalida
	}
	alvo := ultima
	if pedida != nil {
		alvo = *pedida
	}
	for _, a := range lista {
		if a.EmpresaID == alvo {
			return a.EmpresaID, a.UsuarioID, lista, nil
		}
	}
	if pedida != nil {
		return uuid.Nil, uuid.Nil, nil, erros.NaoEncontrado.Com("empresa fora do acesso da identidade")
	}
	return lista[0].EmpresaID, lista[0].UsuarioID, lista, nil
}

// concluirLogin cria a sessão web na empresa escolhida (seção 6.9 do
// documento 5): (1) acessos no pool da autenticação; (2) ator e nomes das
// empresas no pool da aplicação; (3) na trilha da empresa da sessão, numa
// transação só, a sessão, o zerar das falhas e o evento de login; (4) o
// evento nas demais empresas, sem derrubar o login se falhar.
//
// mfaAtivo diz se a identidade tem segundo fator. Se a empresa escolhida o
// exige e ela não tem, a sessão nasce restrita à ativação (RN-ADM-027).
func (s *Servico) concluirLogin(ctx context.Context, identidadeID uuid.UUID, pedida *uuid.UUID, mfaAtivo bool) (ResultadoLogin, error) {
	empresaID, usuarioID, lista, err := s.escolherEmpresa(ctx, identidadeID, pedida)
	if err != nil {
		return ResultadoLogin{}, err
	}
	a, err := s.CarregarAtor(ctx, empresaID, usuarioID)
	if err != nil {
		return ResultadoLogin{}, err
	}
	if a.IdentidadeID != identidadeID {
		return ResultadoLogin{}, erros.CredencialInvalida
	}
	empresas := make([]EmpresaAcesso, 0, len(lista))
	for _, ac := range lista {
		nome, err := s.nomeEmpresa(ctx, ac.EmpresaID)
		if err != nil {
			return ResultadoLogin{}, err
		}
		acesso := EmpresaAcesso{EmpresaID: ac.EmpresaID, Nome: nome, UsuarioID: ac.UsuarioID}
		if ac.EmpresaID == empresaID {
			empresas = append([]EmpresaAcesso{acesso}, empresas...)
		} else {
			empresas = append(empresas, acesso)
		}
	}

	var sessao SessaoWeb
	err = s.repo.NaTrilha(ctx, empresaID, func(q *gerado.Queries) error {
		var err error
		sessao, err = s.criarSessaoWeb(ctx, q, identidadeID, empresaID, usuarioID)
		if err != nil {
			return err
		}
		if err := q.AuthRegistrarLoginCerto(ctx, gerado.AuthRegistrarLoginCertoParams{ID: identidadeID}); err != nil {
			return err
		}
		return gravarEvento(ctx, q, empresaID, usuarioID, eventoLogin(identidadeID, CanalWeb, sessao.SessaoID, ""))
	})
	if err != nil {
		return ResultadoLogin{}, err
	}
	sessao.MFAPendente = a.MFAObrigatorio && !mfaAtivo
	_ = s.eventoNasEmpresas(ctx, &empresaID, eventoAuth{
		Acao: acaoLogin, IdentidadeID: identidadeID, Depois: map[string]string{"canal": CanalWeb},
	})
	return ResultadoLogin{Sessao: &sessao, Empresas: empresas}, nil
}

func eventoLogin(identidadeID uuid.UUID, canal string, sessaoID uuid.UUID, motivo string) eventoAuth {
	depois := map[string]string{"canal": canal, "sessao_id": sessaoID.String()}
	if motivo != "" {
		depois["motivo"] = motivo
	}
	return eventoAuth{Acao: acaoLogin, IdentidadeID: identidadeID, Depois: depois}
}

func (s *Servico) revogarSessao(ctx context.Context, sessaoID uuid.UUID) error {
	return s.repo.Tx(ctx, func(q *gerado.Queries) error {
		if _, err := q.AuthRevogarSessao(ctx, sessaoID); err != nil {
			return err
		}
		return q.AuthRevogarFamilias(ctx, []uuid.UUID{sessaoID})
	})
}

// criarSessaoWeb grava a sessão na transação da autenticação.
func (s *Servico) criarSessaoWeb(ctx context.Context, q *gerado.Queries, identidadeID, empresaID, usuarioID uuid.UUID) (SessaoWeb, error) {
	cookie, hashCookie, err := novoSegredo()
	if err != nil {
		return SessaoWeb{}, err
	}
	csrf := tokenCSRF(cookie)
	agora := s.agora()
	sw := SessaoWeb{
		Cookie: cookie, CSRF: csrf, SessaoID: uuid.Must(uuid.NewV7()),
		EmpresaID: empresaID, UsuarioID: usuarioID, ExpiraEm: agora.Add(s.cfg.MaximoWeb),
	}
	ip, ua := tecnico(ctx)
	err = q.AuthCriarSessao(ctx, gerado.AuthCriarSessaoParams{
		ID: sw.SessaoID, Canal: CanalWeb, TokenHash: hashCookie, CsrfHash: hashSegredo(csrf),
		IdentidadeID: identidadeID, EmpresaID: empresaID, UsuarioID: usuarioID,
		Ip: ip, UserAgent: ua, CriadoEm: agora, ExpiraEm: sw.ExpiraEm,
	})
	if err != nil {
		return SessaoWeb{}, fmt.Errorf("criar sessão web: %w", err)
	}
	return sw, nil
}
