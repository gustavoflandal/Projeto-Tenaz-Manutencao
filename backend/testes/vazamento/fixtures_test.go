package vazamento

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"fmt"
	"sync"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Fixtures da etapa E01: duas empresas, A e B, com pelo menos uma linha em
// cada tabela de app e de plataforma criada até aqui; e, dentro de A, dois
// fornecedores (F1, F2) e dois locatários (L1, L2), cada um com usuário e
// pessoa, para o escopo externo (seção 2.9.5 do documento 2).
//
// As fixtures são gravadas pelo PRÓPRIO tenaz_app, no contexto de cada
// empresa, como faz o serviço de criação de conta (seção 2.3): é a prova de
// que o caminho normal passa pelas policies, e não só o caminho de ataque é
// recusado. Os identificadores são aleatórios a cada execução, para a suíte
// poder rodar mais de uma vez no mesmo banco (TENAZ_TESTE_BANCO_URL).
//
// fornecedor_id e locatario_id são uuids soltos: as tabelas fornecedor (E06) e
// locatario (fase 2) ainda não existem, e as FKs entram com elas (E01-1).
//
// A empresa B tem um usuário prestador com o MESMO fornecedor_id de F1. Na
// vida real os ids diferem; aqui ele simula o prestador que atende duas
// contratantes (seção 9.3 do documento 5) e prova que o escopo nunca alarga o
// corte por empresa.

type empresaFixture struct {
	ID               string
	Unidade          string
	Turno            string
	Feriado          string
	PessoaInterna    string // pessoa com CPF e sem login: o alvo do vazamento por pessoa_id
	PessoaDoInterno  string // pessoa com CPF ligada ao usuário interno: o CPF de quem tem login
	Identidade       string
	UsuarioInterno   string // ativo, com perfil gestor e permissão avulsa
	UsuarioInativo   string
	PerfilGestor     string
	PerfilPrestador  string
	PerfilLocatario  string
	UsuarioPerfil    string
	UsuarioPermissao string
	SessaoWeb        string
	SessaoApp        string
	CodigoAutoriz    string
	DesafioMFA       string
	TokenEmail       string
	TokenRenovacao   string
}

type fixturesE01 struct {
	A, B empresaFixture

	// Escopo externo, todos na empresa A.
	F1, F2, L1, L2 string

	UsuarioF1, UsuarioF1Colega, UsuarioF2 string
	UsuarioL1, UsuarioL1Colega, UsuarioL2 string
	IdentidadeF1                          string
	PessoaUsuarioF1, PessoaF1SemLogin     string
	PessoaUsuarioF2, PessoaF2SemLogin     string
	PessoaUsuarioL1, PessoaUsuarioL2      string
	PerfilPrestadorAtribuido              string // usuario_perfil de UsuarioF1

	// Na empresa B, com fornecedor_id igual a F1.
	UsuarioBMesmoFornecedor string

	LinkHash string // hash de um link sem linha em link_execucao (E06)
}

var (
	fixOnce sync.Once
	fixDado *fixturesE01
	fixErro error
)

// fixtures povoa o banco uma vez por execução da suíte.
func fixtures(t *testing.T) *fixturesE01 {
	t.Helper()
	fixOnce.Do(func() {
		fixDado, fixErro = povoar(context.Background())
	})
	if fixErro != nil {
		t.Fatalf("povoar fixtures de E01: %v", fixErro)
	}
	return fixDado
}

func novoID() string { return uuid.Must(uuid.NewV7()).String() }

func hashAleatorio() []byte {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	h := sha256.Sum256(b)
	return h[:]
}

func cpfAleatorio() string {
	b := make([]byte, 11)
	if _, err := rand.Read(b); err != nil {
		panic(err)
	}
	for i := range b {
		b[i] = '0' + b[i]%10
	}
	return string(b)
}

func povoar(ctx context.Context) (*fixturesE01, error) {
	f := &fixturesE01{F1: novoID(), F2: novoID(), L1: novoID(), L2: novoID()}
	f.LinkHash = fmt.Sprintf("%x", hashAleatorio())

	for _, e := range []*empresaFixture{&f.A, &f.B} {
		if err := povoarEmpresa(ctx, e); err != nil {
			return nil, err
		}
	}
	if err := povoarExternos(ctx, f); err != nil {
		return nil, err
	}
	return f, nil
}

// povoarEmpresa repete a criação de conta (RF-ADM-001): transação com o id
// novo em app.empresa_id, escopo interno, sem usuário.
func povoarEmpresa(ctx context.Context, e *empresaFixture) error {
	e.ID = novoID()
	c := contexto{Empresa: e.ID, Escopo: "interno"}

	return confirmarComoPapel(ctx, papelApp, c, func(tx pgx.Tx) error {
		passos := []struct {
			rotulo string
			sql    string
			args   []any
			saida  *string
		}{
			{"empresa", `INSERT INTO app.empresa (id, razao_social) VALUES ($1, 'Empresa de teste')`,
				[]any{e.ID}, nil},
			{"unidade", `INSERT INTO app.unidade (empresa_id, codigo, nome) VALUES ($1, 'U1', 'Unidade') RETURNING id`,
				[]any{e.ID}, &e.Unidade},
			{"turno", `INSERT INTO app.turno_unidade (empresa_id, unidade_id, dia_semana, hora_inicio, hora_fim)
			           VALUES ($1, $2, 1, '08:00', '17:00') RETURNING id`, nil, &e.Turno},
			{"feriado", `INSERT INTO app.feriado_unidade (empresa_id, unidade_id, data_feriado, descricao)
			             VALUES ($1, $2, '2026-12-25', 'Natal') RETURNING id`, nil, &e.Feriado},
			{"pessoa interna", `INSERT INTO app.pessoa (empresa_id, nome, cpf, data_nascimento)
			                    VALUES ($1, 'Pessoa interna', $2, '1980-01-01') RETURNING id`,
				[]any{e.ID, cpfAleatorio()}, &e.PessoaInterna},
			{"identidade", `INSERT INTO plataforma.identidade (email, senha_hash) VALUES ($1, 'argon2id-falso') RETURNING id`,
				[]any{"interno-" + e.ID + "@teste.tenaz"}, &e.Identidade},
			{"pessoa do interno", `INSERT INTO app.pessoa (empresa_id, nome, cpf, data_nascimento)
			                       VALUES ($1, 'Pessoa do usuário interno', $2, '1975-05-05') RETURNING id`,
				[]any{e.ID, cpfAleatorio()}, &e.PessoaDoInterno},
			{"usuário interno", `INSERT INTO app.usuario (empresa_id, nome, identidade_id, pessoa_id)
			                     VALUES ($1, 'Interno', $2, $3) RETURNING id`,
				nil, &e.UsuarioInterno},
			{"usuário inativo", `INSERT INTO app.usuario (empresa_id, nome, situacao) VALUES ($1, 'Inativo', 'inativo') RETURNING id`,
				[]any{e.ID}, &e.UsuarioInativo},
			{"perfil gestor", `INSERT INTO app.perfil (empresa_id, codigo, nome, padrao, perfil_origem, permissoes)
			                   VALUES ($1, 'gestor', 'Gestor', true, 'gestor', '{os.encerrar,os.executar}') RETURNING id`,
				[]any{e.ID}, &e.PerfilGestor},
			{"perfil prestador", `INSERT INTO app.perfil (empresa_id, codigo, nome, padrao, perfil_origem, permissoes)
			                      VALUES ($1, 'prestador', 'Prestador', true, 'prestador', '{os.executar}') RETURNING id`,
				[]any{e.ID}, &e.PerfilPrestador},
			{"perfil locatário", `INSERT INTO app.perfil (empresa_id, codigo, nome, padrao, perfil_origem, permissoes)
			                      VALUES ($1, 'locatario', 'Locatário', true, 'locatario', '{sol.abrir}') RETURNING id`,
				[]any{e.ID}, &e.PerfilLocatario},
			{"usuario_perfil", `INSERT INTO app.usuario_perfil (empresa_id, usuario_id, perfil_id) VALUES ($1, $2, $3) RETURNING id`,
				nil, &e.UsuarioPerfil},
			{"usuario_permissao", `INSERT INTO app.usuario_permissao (empresa_id, usuario_id, permissao)
			                       VALUES ($1, $2, 'os.encerrar') RETURNING id`, nil, &e.UsuarioPermissao},
		}
		for _, p := range passos {
			// Argumentos que dependem de passos anteriores são montados aqui.
			args := p.args
			switch p.rotulo {
			case "turno", "feriado":
				args = []any{e.ID, e.Unidade}
			case "usuário interno":
				args = []any{e.ID, e.Identidade, e.PessoaDoInterno}
			case "usuario_perfil":
				args = []any{e.ID, e.UsuarioInterno, e.PerfilGestor}
			case "usuario_permissao":
				args = []any{e.ID, e.UsuarioInterno}
			}
			var err error
			if p.rotulo == "identidade" {
				// Credencial só pelo papel da autenticação (ADR-011).
				err = comoAutenticacao(ctx, tx, func() error {
					return tx.QueryRow(ctx, p.sql, args...).Scan(p.saida)
				})
			} else if p.saida != nil {
				err = tx.QueryRow(ctx, p.sql, args...).Scan(p.saida)
			} else {
				_, err = tx.Exec(ctx, p.sql, args...)
			}
			if err != nil {
				return fmt.Errorf("empresa %s, %s: %w", e.ID, p.rotulo, err)
			}
		}
		return comoAutenticacao(ctx, tx, func() error { return povoarCredenciais(ctx, tx, e) })
	})
}

// comoAutenticacao troca, dentro da transação da fixture, para o papel das
// credenciais (ADR-011) e volta ao tenaz_app. É o caminho do serviço: só o
// tenaz_autenticacao grava nas tabelas de autenticação.
func comoAutenticacao(ctx context.Context, tx pgx.Tx, fn func() error) error {
	if _, err := tx.Exec(ctx, "SET LOCAL ROLE tenaz_autenticacao"); err != nil {
		return fmt.Errorf("assumir tenaz_autenticacao: %w", err)
	}
	if err := fn(); err != nil {
		return err
	}
	_, err := tx.Exec(ctx, "SET LOCAL ROLE "+pgx.Identifier{papelApp}.Sanitize())
	return err
}

// povoarCredenciais grava uma linha em cada tabela de autenticação, pelo
// mesmo caminho do serviço: INSERT do tenaz_autenticacao (ADR-011).
func povoarCredenciais(ctx context.Context, tx pgx.Tx, e *empresaFixture) error {
	passos := []struct {
		rotulo string
		sql    string
		args   []any
		saida  *string
	}{
		{"sessão web", `INSERT INTO plataforma.sessao (canal, token_hash, csrf_hash, identidade_id, empresa_id, usuario_id, expira_em)
		                VALUES ('web', $1, $2, $3, $4, $5, now() + interval '12 hours') RETURNING id`,
			[]any{hashAleatorio(), hashAleatorio(), e.Identidade, e.ID, e.UsuarioInterno}, &e.SessaoWeb},
		{"sessão app", `INSERT INTO plataforma.sessao (canal, identidade_id, empresa_id, usuario_id, expira_em)
		                VALUES ('app', $1, $2, $3, now() + interval '30 days') RETURNING id`,
			[]any{e.Identidade, e.ID, e.UsuarioInterno}, &e.SessaoApp},
		{"código OAuth", `INSERT INTO plataforma.codigo_autorizacao (codigo_hash, identidade_id, empresa_id, usuario_id,
		                    cliente_oauth, uri_retorno, desafio_pkce, expira_em)
		                  VALUES ($1, $2, $3, $4, 'tenaz-app', 'tenaz://retorno', $5, now() + interval '60 seconds') RETURNING id`,
			[]any{hashAleatorio(), e.Identidade, e.ID, e.UsuarioInterno, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQ"}, &e.CodigoAutoriz},
		{"desafio MFA", `INSERT INTO plataforma.desafio_mfa (token_hash, identidade_id, canal, expira_em)
		                 VALUES ($1, $2, 'web', now() + interval '5 minutes') RETURNING id`,
			[]any{hashAleatorio(), e.Identidade}, &e.DesafioMFA},
		{"convite", `INSERT INTO plataforma.token_email (token_hash, finalidade, email, empresa_id, usuario_id, expira_em)
		             VALUES ($1, 'convite', 'convite@teste.tenaz', $2, $3, now() + interval '7 days') RETURNING id`,
			[]any{hashAleatorio(), e.ID, e.UsuarioInativo}, &e.TokenEmail},
		{"refresh token", `INSERT INTO plataforma.token_renovacao (familia_id, token_hash, expira_em)
		                   VALUES ($1, $2, now() + interval '30 days') RETURNING id`, nil, &e.TokenRenovacao},
	}
	for _, p := range passos {
		args := p.args
		if p.rotulo == "refresh token" {
			args = []any{e.SessaoApp, hashAleatorio()}
		}
		if err := tx.QueryRow(ctx, p.sql, args...).Scan(p.saida); err != nil {
			return fmt.Errorf("empresa %s, %s: %w", e.ID, p.rotulo, err)
		}
	}
	return nil
}

// povoarExternos cria, em A, os usuários e pessoas de F1, F2, L1 e L2; e, em
// B, o usuário com o mesmo fornecedor_id de F1.
func povoarExternos(ctx context.Context, f *fixturesE01) error {
	a := f.A.ID
	err := confirmarComoPapel(ctx, papelApp, contexto{Empresa: a, Escopo: "interno"}, func(tx pgx.Tx) error {
		inserir := func(rotulo, sql string, saida *string, args ...any) error {
			if err := tx.QueryRow(ctx, sql, args...).Scan(saida); err != nil {
				return fmt.Errorf("%s: %w", rotulo, err)
			}
			return nil
		}
		const pessoaTerceiro = `INSERT INTO app.pessoa (empresa_id, nome, cpf, vinculo, fornecedor_id)
		                        VALUES ($1, $2, $3, 'terceiro', $4) RETURNING id`
		const pessoaSemFornecedor = `INSERT INTO app.pessoa (empresa_id, nome, cpf, vinculo)
		                             VALUES ($1, $2, $3, 'autonomo') RETURNING id`
		const usuarioPrestador = `INSERT INTO app.usuario (empresa_id, nome, fornecedor_id, pessoa_id, identidade_id)
		                          VALUES ($1, $2, $3, $4, $5) RETURNING id`
		const usuarioLocatario = `INSERT INTO app.usuario (empresa_id, nome, locatario_id, pessoa_id)
		                          VALUES ($1, $2, $3, $4) RETURNING id`

		if err := comoAutenticacao(ctx, tx, func() error {
			return inserir("identidade F1", `INSERT INTO plataforma.identidade (email) VALUES ($1) RETURNING id`,
				&f.IdentidadeF1, "f1-"+a+"@teste.tenaz")
		}); err != nil {
			return err
		}
		for _, p := range []struct {
			rotulo string
			sql    string
			saida  *string
			args   []any
		}{
			{"pessoa do usuário F1", pessoaTerceiro, &f.PessoaUsuarioF1, []any{a, "Técnico F1", cpfAleatorio(), f.F1}},
			{"pessoa F1 sem login", pessoaTerceiro, &f.PessoaF1SemLogin, []any{a, "Ajudante F1", cpfAleatorio(), f.F1}},
			{"pessoa do usuário F2", pessoaTerceiro, &f.PessoaUsuarioF2, []any{a, "Técnico F2", cpfAleatorio(), f.F2}},
			{"pessoa F2 sem login", pessoaTerceiro, &f.PessoaF2SemLogin, []any{a, "Ajudante F2", cpfAleatorio(), f.F2}},
			{"pessoa do usuário L1", pessoaSemFornecedor, &f.PessoaUsuarioL1, []any{a, "Locatário L1", cpfAleatorio()}},
			{"pessoa do usuário L2", pessoaSemFornecedor, &f.PessoaUsuarioL2, []any{a, "Locatário L2", cpfAleatorio()}},
		} {
			if err := inserir(p.rotulo, p.sql, p.saida, p.args...); err != nil {
				return err
			}
		}
		for _, p := range []struct {
			rotulo string
			sql    string
			saida  *string
			args   []any
		}{
			{"usuário F1", usuarioPrestador, &f.UsuarioF1, []any{a, "Prestador F1", f.F1, f.PessoaUsuarioF1, f.IdentidadeF1}},
			{"colega F1", usuarioPrestador, &f.UsuarioF1Colega, []any{a, "Colega F1", f.F1, nil, nil}},
			{"usuário F2", usuarioPrestador, &f.UsuarioF2, []any{a, "Prestador F2", f.F2, f.PessoaUsuarioF2, nil}},
			{"usuário L1", usuarioLocatario, &f.UsuarioL1, []any{a, "Locatário L1", f.L1, f.PessoaUsuarioL1}},
			{"colega L1", usuarioLocatario, &f.UsuarioL1Colega, []any{a, "Colega L1", f.L1, nil}},
			{"usuário L2", usuarioLocatario, &f.UsuarioL2, []any{a, "Locatário L2", f.L2, f.PessoaUsuarioL2}},
		} {
			if err := inserir(p.rotulo, p.sql, p.saida, p.args...); err != nil {
				return err
			}
		}
		return inserir("perfil prestador de F1",
			`INSERT INTO app.usuario_perfil (empresa_id, usuario_id, perfil_id) VALUES ($1, $2, $3) RETURNING id`,
			&f.PerfilPrestadorAtribuido, a, f.UsuarioF1, f.A.PerfilPrestador)
	})
	if err != nil {
		return fmt.Errorf("externos de A: %w", err)
	}

	return confirmarComoPapel(ctx, papelApp, contexto{Empresa: f.B.ID, Escopo: "interno"}, func(tx pgx.Tx) error {
		return tx.QueryRow(ctx,
			`INSERT INTO app.usuario (empresa_id, nome, fornecedor_id) VALUES ($1, 'Prestador de B', $2) RETURNING id`,
			f.B.ID, f.F1).Scan(&f.UsuarioBMesmoFornecedor)
	})
}

// Contextos prontos.
func (f *fixturesE01) internoA() contexto { return contexto{Empresa: f.A.ID, Escopo: "interno"} }

func (f *fixturesE01) prestadorF1() contexto {
	return contexto{Empresa: f.A.ID, Escopo: "prestador", Usuario: f.UsuarioF1, Fornecedor: f.F1}
}

func (f *fixturesE01) locatarioL1() contexto {
	return contexto{Empresa: f.A.ID, Escopo: "locatario", Usuario: f.UsuarioL1, Locatario: f.L1}
}

func (f *fixturesE01) linkA() contexto {
	return contexto{Empresa: f.A.ID, Escopo: "link", LinkHash: f.LinkHash}
}
