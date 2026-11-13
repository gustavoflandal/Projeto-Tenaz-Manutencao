// Package autenticacao é o serviço de autenticação da Tenaz (RF-ADM-004,
// RF-ADM-006, RF-MOB-001, RNF-017; documento 3, seção 4; documento 9, seção
// 5.3): senha com Argon2id, bloqueio temporário, MFA por TOTP, sessão web por
// cookie com CSRF, OAuth 2.1 com PKCE para o app, refresh token rotativo com
// detecção de reuso, escolha e troca de empresa, links de e-mail (convite e
// redefinição de senha) e os eventos de auditoria de login.
//
// Dois pools, de propósito (achado ISO-E01-02):
//
//   - o pool da autenticação, injetado, é o único que toca as tabelas de
//     credencial do esquema plataforma, pelo repositório em internal/repo;
//   - o pool da aplicação, por NaEmpresa, lê o usuário, os perfis e a empresa
//     (autorizacao.Carregar) e grava os eventos na trilha de cada empresa.
//
// Nenhuma operação mistura os dois numa transação: o repositório recusa ser
// chamado de dentro de NaEmpresa. Onde a ordem importa, a escrita de
// credencial vem primeiro e a trilha depois; se a trilha falhar, a credencial
// é revogada (não existe login sem evento).
//
// É plataforma, não módulo (seção 5.1 do documento 5): os módulos chamam este
// serviço; ele não importa módulo nenhum.
package autenticacao

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"tenaz/internal/plataforma/autenticacao/internal/repo"
	"tenaz/internal/plataforma/autenticacao/internal/repo/gerado"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/email"
	"tenaz/internal/plataforma/erros"
	"tenaz/internal/plataforma/obs"
)

// Canais de sessão (sessao.canal).
const (
	CanalWeb = "web"
	CanalApp = "app"
)

// Config são os prazos e os valores fixos da autenticação. Todos têm teto no
// banco (CHECKs de 00005_autenticacao.sql); estes são os valores usados.
type Config struct {
	URLWeb string // base dos links enviados por e-mail, sem barra no fim

	InatividadeWeb time.Duration // 12 h (RNF-017)
	MaximoWeb      time.Duration // 7 dias (documento 3, seção 4.2)
	ValidadeApp    time.Duration // 30 dias de inatividade (RNF-017)

	ValidadeCodigoOAuth time.Duration // 60 s
	ValidadeDesafioMFA  time.Duration // 5 min
	ValidadeConvite     time.Duration // 7 dias
	ValidadeRedefinicao time.Duration // 1 h

	FalhasParaBloquear int           // 5 (RF-ADM-006)
	BloqueioInicial    time.Duration // 15 min, dobrando a cada novo bloqueio
	BloqueioMaximo     time.Duration // 24 h

	ClienteOAuthApp   string   // client_id do app
	URIsRetornoApp    []string // redirect_uri aceitas, conferidas byte a byte
	EmissorTOTP       string
	IntervaloToqueWeb time.Duration // não regrava ultimo_uso_em a cada requisição
}

// ConfigPadrao devolve os valores da documentação.
func ConfigPadrao(urlWeb string) Config {
	return Config{
		URLWeb:              strings.TrimRight(urlWeb, "/"),
		InatividadeWeb:      12 * time.Hour,
		MaximoWeb:           7 * 24 * time.Hour,
		ValidadeApp:         30 * 24 * time.Hour,
		ValidadeCodigoOAuth: 60 * time.Second,
		ValidadeDesafioMFA:  5 * time.Minute,
		ValidadeConvite:     7 * 24 * time.Hour,
		ValidadeRedefinicao: time.Hour,
		FalhasParaBloquear:  5,
		BloqueioInicial:     15 * time.Minute,
		BloqueioMaximo:      24 * time.Hour,
		ClienteOAuthApp:     "tenaz-app",
		URIsRetornoApp:      []string{"br.com.tenaz.app:/oauth"},
		EmissorTOTP:         "Tenaz",
		IntervaloToqueWeb:   time.Minute,
	}
}

// Dependencias do serviço, todas obrigatórias.
type Dependencias struct {
	// Banco é o pool da aplicação: só NaEmpresa.
	Banco *db.DB
	// PoolAutenticacao é o pool das tabelas de credencial
	// (TENAZ_BANCO_URL_AUTENTICACAO). Ponto de troca do papel de banco.
	PoolAutenticacao *pgxpool.Pool

	Hasher    *Hasher
	Politica  *Politica
	Cifrador  *Cifrador
	Assinador *Assinador
	Email     email.Enviador
	Agora     func() time.Time // relógio; nil usa time.Now
}

// Servico é o serviço de autenticação.
type Servico struct {
	cfg  Config
	d    Dependencias
	repo *repo.Repo

	// envios em segundo plano (recuperação de senha), para o desligamento
	// esperar por eles
	envios sync.WaitGroup
}

// Novo monta o serviço.
func Novo(cfg Config, d Dependencias) (*Servico, error) {
	if d.Banco == nil || d.PoolAutenticacao == nil || d.Hasher == nil || d.Politica == nil || d.Cifrador == nil ||
		d.Assinador == nil || d.Email == nil {
		return nil, errors.New("autenticacao: dependência ausente")
	}
	if cfg.URLWeb == "" || cfg.ClienteOAuthApp == "" || len(cfg.URIsRetornoApp) == 0 {
		return nil, errors.New("autenticacao: configuração incompleta")
	}
	if cfg.ValidadeApp > 30*24*time.Hour || cfg.MaximoWeb > 7*24*time.Hour ||
		cfg.ValidadeConvite > 7*24*time.Hour || cfg.ValidadeRedefinicao > time.Hour {
		return nil, errors.New("autenticacao: validade acima do teto do banco")
	}
	if d.Agora == nil {
		d.Agora = time.Now
	}
	return &Servico{cfg: cfg, d: d, repo: repo.Novo(d.PoolAutenticacao)}, nil
}

// Aguardar espera os envios de e-mail em segundo plano. O processo chama no
// desligamento; os testes, antes de conferir a caixa.
func (s *Servico) Aguardar() { s.envios.Wait() }

func (s *Servico) agora() time.Time { return s.d.Agora().UTC() }

// ------------------------------------------------------------------- tipos

// EmpresaAcesso é uma empresa em que a identidade entra.
type EmpresaAcesso struct {
	EmpresaID uuid.UUID
	Nome      string // nome fantasia, ou a razão social
	UsuarioID uuid.UUID
}

// DesafioMFA é o passo entre a senha certa e o TOTP. A camada HTTP responde
// 401 mfa_obrigatorio com o DesafioID.
type DesafioMFA struct {
	DesafioID string
	ExpiraEm  time.Time
}

// SessaoWeb é a sessão criada. Cookie vai em __Host-tenaz_sessao (HttpOnly,
// Secure, SameSite=Lax, Path=/); CSRF, no corpo, para o cabeçalho X-CSRF-Token.
type SessaoWeb struct {
	Cookie    string
	CSRF      string
	SessaoID  uuid.UUID
	EmpresaID uuid.UUID
	UsuarioID uuid.UUID
	ExpiraEm  time.Time // teto absoluto; a inatividade de 12 h vence antes

	// MFAPendente: a empresa exige MFA e a identidade não o tem. A sessão só
	// serve para ativá-lo (RN-ADM-027).
	MFAPendente bool
}

// ResultadoLogin traz o desafio de MFA ou a sessão, nunca os dois.
type ResultadoLogin struct {
	DesafioMFA *DesafioMFA
	Sessao     *SessaoWeb
	Empresas   []EmpresaAcesso // com a sessão: as empresas para o seletor
}

// TokensApp é a resposta de POST /oauth/token.
type TokensApp struct {
	AccessToken  string
	ExpiraEm     time.Time
	ExpiraEmSegs int64 // expires_in
	RefreshToken string
	SessaoID     uuid.UUID
	EmpresaID    uuid.UUID
	UsuarioID    uuid.UUID
}

// ------------------------------------------------------------- auxiliares

// naoEncontrado traduz "nenhuma linha" para o erro dado.
func naoEncontrado(err error, alvo *erros.Erro) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return alvo
	}
	return err
}

func normalizarEmail(e string) string {
	return strings.ToLower(strings.TrimSpace(e))
}

func ptr[T any](v T) *T { return &v }

// ErrMFAPendente: a empresa exige segundo fator e a identidade ainda não o
// ativou (RN-ADM-027). A sessão web existe, mas só serve para ativá-lo; o app
// e as demais rotas respondem como não autenticado.
func ErrMFAPendente() *erros.Erro {
	return erros.NaoAutenticado.ComChave("detalhe.mfa_empresa_ativar", nil).
		ComDado("mfa_pendente", true)
}

// mfaAtivo lê, por chave exata, se a identidade tem MFA ativo.
func (s *Servico) mfaAtivo(ctx context.Context, identidadeID uuid.UUID) (bool, error) {
	var ativo bool
	err := s.repo.Tx(ctx, func(q *gerado.Queries) error {
		id, err := q.AuthBuscarIdentidade(ctx, identidadeID)
		if err != nil {
			return naoEncontrado(err, erros.CredencialInvalida)
		}
		ativo = id.MfaAtivadoEm != nil
		return nil
	})
	return ativo, err
}

// marcarMFAPendente liga a restrição da sessão quando a empresa exige MFA e
// a identidade não o tem.
func (s *Servico) marcarMFAPendente(ctx context.Context, a *autorizacao.Ator) error {
	if !a.MFAObrigatorio {
		return nil
	}
	ativo, err := s.mfaAtivo(ctx, a.IdentidadeID)
	if err != nil {
		return err
	}
	a.SomenteAtivarMFA = !ativo
	return nil
}

// consumirTOTP confere o código e grava o passo, numa escrita atômica: código
// certo de um passo já usado vale como errado.
func consumirTOTP(ctx context.Context, q *gerado.Queries, identidadeID uuid.UUID, segredo []byte, codigo string, agora time.Time) (bool, error) {
	passo, ok := passoTOTP(segredo, codigo, agora)
	if !ok {
		return false, nil
	}
	n, err := q.AuthUsarPassoTOTP(ctx, gerado.AuthUsarPassoTOTPParams{ID: identidadeID, Passo: ptr(int64(passo))})
	if err != nil {
		return false, err
	}
	return n == 1, nil
}

// CarregarAtor lê o Ator no pool da aplicação, em NaEmpresa com escopo
// interno e sem usuário (seção 2.9.2 do documento 2).
func (s *Servico) CarregarAtor(ctx context.Context, empresaID, usuarioID uuid.UUID) (autorizacao.Ator, error) {
	var a autorizacao.Ator
	err := s.d.Banco.NaEmpresa(ctx, empresaID, db.Interno(nil), func(ctx context.Context) error {
		var err error
		a, err = autorizacao.Carregar(ctx, empresaID, usuarioID)
		return err
	})
	return a, err
}

// acessos lista as empresas habilitadas da identidade (busca por chave
// exata, identidade_id).
func (s *Servico) acessos(ctx context.Context, identidadeID uuid.UUID) ([]gerado.AuthListarAcessosRow, error) {
	var lista []gerado.AuthListarAcessosRow
	err := s.repo.Tx(ctx, func(q *gerado.Queries) error {
		var err error
		lista, err = q.AuthListarAcessos(ctx, identidadeID)
		return err
	})
	return lista, err
}

// eventoAuth é um evento de autenticação: ator "usuario", entidade
// "identidade" (seção 6.9 do documento 5).
type eventoAuth struct {
	Acao         string // login, logout, alteracao
	Falha        bool   // resultado = falha: só em login
	IdentidadeID uuid.UUID
	Depois       any
	Campos       []string
}

// Ações aceitas pela política restritiva do tenaz_autenticacao na trilha.
const (
	acaoLogin     = "login"
	acaoLogout    = "logout"
	acaoAlteracao = "alteracao"
)

// Resultados do evento (evento_auditoria.resultado): a falha de login é um
// login com resultado falha, e o motivo vai em depois.
const (
	resultadoSucesso = "sucesso"
	resultadoFalha   = "falha"
)

// gravarEvento grava o evento na transação NaTrilha de q. usuarioID é o
// usuário da identidade naquela empresa (a política confere em
// acesso_empresa).
func gravarEvento(ctx context.Context, q *gerado.Queries, empresaID, usuarioID uuid.UUID, ev eventoAuth) error {
	var depois []byte
	if ev.Depois != nil {
		b, err := json.Marshal(ev.Depois)
		if err != nil {
			return err
		}
		depois = b
	}
	resultado := resultadoSucesso
	if ev.Falha {
		resultado = resultadoFalha
	}
	nome, err := q.AuthNomeDoAcesso(ctx, gerado.AuthNomeDoAcessoParams{IdentidadeID: ev.IdentidadeID, EmpresaID: empresaID})
	if err != nil {
		return fmt.Errorf("ler o nome do usuário para o evento %s: %w", ev.Acao, err)
	}
	req := db.RequisicaoDe(ctx).Saneada()
	var tela *string
	if req.Tela != "" {
		tela = &req.Tela
	}
	err = q.AuthInserirEvento(ctx, gerado.AuthInserirEventoParams{
		EmpresaID: empresaID, UsuarioID: &usuarioID, UsuarioNome: &nome, Acao: ev.Acao, Resultado: resultado,
		IdentidadeID: &ev.IdentidadeID, Depois: depois, CamposAlterados: ev.Campos,
		DispositivoID: req.DispositivoID, Tela: tela,
	})
	if err != nil {
		return fmt.Errorf("gravar evento %s na trilha: %w", ev.Acao, err)
	}
	return nil
}

// eventoNasEmpresas grava o evento na trilha de cada empresa em que a
// identidade tem acesso habilitado (seção 9.1 do documento 2), uma transação
// por empresa. O evento de uma empresa nunca cita outra: dizer à empresa B
// que a pessoa também entra na empresa A seria vazamento entre empresas.
//
// Falha numa empresa não interrompe as outras; o erro volta no fim.
func (s *Servico) eventoNasEmpresas(ctx context.Context, exceto *uuid.UUID, ev eventoAuth) error {
	lista, err := s.acessos(ctx, ev.IdentidadeID)
	if err != nil {
		return err
	}
	var falhas []error
	for _, a := range lista {
		if exceto != nil && a.EmpresaID == *exceto {
			continue
		}
		err := s.repo.NaTrilha(ctx, a.EmpresaID, func(q *gerado.Queries) error {
			return gravarEvento(ctx, q, a.EmpresaID, a.UsuarioID, ev)
		})
		if err != nil {
			falhas = append(falhas, err)
		}
	}
	if err := errors.Join(falhas...); err != nil {
		obs.Logger(ctx).Error("falha ao gravar evento de autenticação na trilha", "acao", ev.Acao, "erro", err)
		return err
	}
	return nil
}

// nomeEmpresa lê o nome de exibição da empresa, no pool da aplicação.
func (s *Servico) nomeEmpresa(ctx context.Context, empresaID uuid.UUID) (string, error) {
	var nome string
	err := s.d.Banco.NaEmpresa(ctx, empresaID, db.Interno(nil), func(ctx context.Context) error {
		var err error
		nome, err = autorizacao.NomeEmpresa(ctx, empresaID)
		return err
	})
	return nome, err
}
