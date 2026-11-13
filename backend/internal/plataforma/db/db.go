// Package db dá acesso ao PostgreSQL com o contexto de empresa e de escopo
// sempre definido na transação.
//
// Esta é a terceira das quatro camadas de isolamento da seção 1.2 de
// docs/2_ARQUITETURA_BANCO_DADOS.md. As outras três são o RLS forçado no
// banco, o filtro explícito por empresa_id em toda consulta e as chaves
// estrangeiras compostas com empresa_id.
//
// Regra inviolável: nenhum código toca dado de cliente fora de NaEmpresa.
package db

import (
	"context"
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"tenaz/internal/plataforma/db/consultas"
	"tenaz/internal/plataforma/erros"
)

// Tipos de escopo. São só estes quatro, pela seção 3.6 do sumário executivo.
const (
	EscopoInterno   = "interno"
	EscopoPrestador = "prestador"
	EscopoLocatario = "locatario"
	EscopoLink      = "link"
)

// Escopo é a camada de isolamento adicional aplicada dentro da empresa, para
// prestador, locatário e link de uso único (risco R4 do sumário).
//
// O escopo é resolvido no servidor a partir do banco, a cada requisição.
// Nunca vem do token, do corpo da requisição nem da URL.
type Escopo struct {
	Tipo string

	FornecedorID *uuid.UUID // escopo prestador
	LocatarioID  *uuid.UUID // escopo locatario
	LinkHashHex  string     // escopo link: sha256 do segredo, em hexadecimal

	UsuarioID    *uuid.UUID
	PodeLerSaude bool // permissão pessoa.saude.ler, das permissões efetivas
}

// Interno é o escopo de usuários da própria empresa, de jobs e de chaves de API.
func Interno(usuarioID *uuid.UUID) Escopo {
	return Escopo{Tipo: EscopoInterno, UsuarioID: usuarioID}
}

var hexSha256 = regexp.MustCompile(`^[0-9a-f]{64}$`)

// Valido recusa escopo incoerente antes de chegar ao banco. O banco também
// recusa, por app.conferir_escopo(); esta é a checagem barata que vem antes.
func (e Escopo) Valido() error {
	switch e.Tipo {
	case EscopoInterno:
		if e.FornecedorID != nil || e.LocatarioID != nil || e.LinkHashHex != "" {
			return errors.New("escopo interno não pode ter fornecedor, locatário nem link")
		}
	case EscopoPrestador:
		if e.FornecedorID == nil || e.LocatarioID != nil || e.LinkHashHex != "" {
			return errors.New("escopo prestador exige fornecedor_id, e só ele")
		}
		if e.UsuarioID == nil {
			return errors.New("escopo prestador exige usuário")
		}
	case EscopoLocatario:
		if e.LocatarioID == nil || e.FornecedorID != nil || e.LinkHashHex != "" {
			return errors.New("escopo locatario exige locatario_id, e só ele")
		}
		if e.UsuarioID == nil {
			return errors.New("escopo locatario exige usuário")
		}
	case EscopoLink:
		if !hexSha256.MatchString(e.LinkHashHex) {
			return errors.New("escopo link exige o sha256 do segredo, em hexadecimal")
		}
		if e.UsuarioID != nil || e.FornecedorID != nil || e.LocatarioID != nil {
			return errors.New("escopo link não tem usuário, fornecedor nem locatário")
		}
	default:
		return fmt.Errorf("escopo inválido: %q", e.Tipo)
	}
	// Dado de saúde nunca é lido fora do escopo interno: pessoa.saude.ler é
	// vedada a usuário externo (RF-ADM-005) e o link não tem usuário.
	if e.PodeLerSaude && e.Tipo != EscopoInterno {
		return errors.New("pessoa.saude.ler só existe no escopo interno")
	}
	return nil
}

// Opcoes configuram o pool.
type Opcoes struct {
	URL     string
	MaxConn int32

	// Papel, quando definido, é assumido por SET ROLE em toda conexão nova.
	// Serve ao desenvolvimento local e aos testes, em que a URL é a do
	// migrador (superusuário, que ignora o RLS): com Papel = "tenaz_app" a
	// aplicação roda com as mesmas policies de produção. Em produção a URL já
	// é a de tenaz_app e Papel fica vazio.
	Papel string
}

// DB é o pool de conexões.
type DB struct {
	pool *pgxpool.Pool
}

// Abrir monta o pool.
func Abrir(ctx context.Context, url string, maxConn int32) (*DB, error) {
	return AbrirComOpcoes(ctx, Opcoes{URL: url, MaxConn: maxConn})
}

var papelValido = regexp.MustCompile(`^tenaz_(app|job|leitura|autenticacao)$`)

// AbrirComOpcoes monta o pool da aplicação, usado só por NaEmpresa.
func AbrirComOpcoes(ctx context.Context, o Opcoes) (*DB, error) {
	pool, err := AbrirPool(ctx, o)
	if err != nil {
		return nil, err
	}
	return &DB{pool: pool}, nil
}

// AbrirPool monta um pool cru. Existe para o pool próprio da autenticação
// (TENAZ_BANCO_URL_AUTENTICACAO), que não passa por NaEmpresa: as tabelas de
// credencial ficam fora do alcance do pool da aplicação (ISO-E01-02).
func AbrirPool(ctx context.Context, o Opcoes) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(o.URL)
	if err != nil {
		return nil, fmt.Errorf("interpretar TENAZ_BANCO_URL: %w", err)
	}
	if o.MaxConn > 0 {
		cfg.MaxConns = o.MaxConn
	}
	if o.Papel != "" {
		// Lista fechada: o papel vai para um comando SQL sem parâmetro (SET
		// ROLE não aceita $1), então só nomes conhecidos passam.
		if !papelValido.MatchString(o.Papel) {
			return nil, fmt.Errorf("papel de banco não aceito: %q", o.Papel)
		}
		comando := "SET ROLE " + o.Papel
		cfg.AfterConnect = func(ctx context.Context, c *pgx.Conn) error {
			_, err := c.Exec(ctx, comando)
			return err
		}
	}

	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("abrir pool: %w", err)
	}
	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("alcançar o banco: %w", err)
	}
	return pool, nil
}

// Fechar devolve as conexões.
func (db *DB) Fechar() {
	if db != nil && db.pool != nil {
		db.pool.Close()
	}
}

// Pool expõe o pool para usos que não tocam dado de empresa: migração,
// verificação de saúde e as suítes de vazamento.
func (db *DB) Pool() *pgxpool.Pool { return db.pool }

type chaveTx struct{}

// ErrTransacaoAninhada: NaEmpresa chamado dentro de outra
// transação. Uma transação por vez é o que garante que nenhum código tenha
// duas empresas abertas ao mesmo tempo (seção 2.4 do documento 2). Quem já
// está dentro de uma transação usa a do contexto, por Tx ou Q.
var ErrTransacaoAninhada = erros.SemContextoEmpresa.Com("transação aninhada: use a transação do contexto")

// NaEmpresa abre uma transação, fixa o contexto de empresa, de escopo e da
// requisição e coloca a transação no contexto. Todo caso de uso que toca dado
// de cliente passa por aqui.
//
// O contexto é definido com set_config(..., is_local => true), local à
// transação. Nunca SET de sessão: no pool, um SET de sessão vazaria para a
// requisição seguinte, de outra empresa.
//
// Os dados da requisição (ator, origem, IP, dispositivo...) vêm de
// ComRequisicao; ausentes, a trilha grava nulo e o ator vira "usuario" quando
// o escopo tem usuário e "sistema" quando não tem.
func (db *DB) NaEmpresa(ctx context.Context, empresaID uuid.UUID, escopo Escopo, fn func(context.Context) error) error {
	if empresaID == uuid.Nil {
		return erros.SemContextoEmpresa.Com("empresa_id vazio em NaEmpresa")
	}
	if err := escopo.Valido(); err != nil {
		return erros.SemContextoEmpresa.Por(err)
	}
	if EmTransacao(ctx) {
		return ErrTransacaoAninhada
	}
	req, err := RequisicaoDe(ctx).normalizada(escopo)
	if err != nil {
		return erros.SemContextoEmpresa.Por(err)
	}

	return pgx.BeginTxFunc(ctx, db.pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		if err := fixarContexto(ctx, tx, empresaID, escopo, req); err != nil {
			return err
		}
		return fn(context.WithValue(ctx, chaveTx{}, tx))
	})
}

// contextoSQL fixa todas as variáveis numa única instrução: uma ida ao banco,
// junto com a conferência do escopo (seção 2.4 do documento 2).
const contextoSQL = `SELECT
  set_config('app.empresa_id',                 $1,  true),
  set_config('app.escopo',                     $2,  true),
  set_config('app.usuario_id',                 $3,  true),
  set_config('app.fornecedor_id',              $4,  true),
  set_config('app.locatario_id',               $5,  true),
  set_config('app.link_hash',                  $6,  true),
  set_config('app.saude',                      $7,  true),
  set_config('app.ator_tipo',                  $8,  true),
  set_config('app.origem',                     $9,  true),
  set_config('app.ip',                         $10, true),
  set_config('app.request_id',                 $11, true),
  set_config('app.user_agent',                 $12, true),
  set_config('app.token_api_id',               $13, true),
  set_config('app.dispositivo_id',             $14, true),
  set_config('app.ocorrido_no_dispositivo_em', $15, true)`

func fixarContexto(ctx context.Context, tx pgx.Tx, empresaID uuid.UUID, escopo Escopo, req Requisicao) error {
	lote := &pgx.Batch{}
	lote.Queue(contextoSQL,
		empresaID.String(), escopo.Tipo, textoUUID(escopo.UsuarioID), textoUUID(escopo.FornecedorID),
		textoUUID(escopo.LocatarioID), escopo.LinkHashHex, simNao(escopo.PodeLerSaude),
		req.AtorTipo, req.Origem, req.IP, req.RequestID, req.UserAgent,
		textoUUID(req.TokenAPIID), textoUUID(req.DispositivoID), textoInstante(req.OcorridoNoDispositivoEm))
	// O banco confere o escopo contra o cadastro e recusa a transação se
	// divergirem. Falha fechada.
	lote.Queue(`SELECT app.conferir_escopo()`)

	if err := tx.SendBatch(ctx, lote).Close(); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "42501" && escopo.UsuarioID != nil {
			// Usuário inativado, excluído ou com vínculo externo trocado
			// durante a sessão: a credencial deixou de valer (seção 2.9.3).
			return erros.CredencialInvalida.Com("o acesso do usuário mudou; entre de novo").Por(err)
		}
		return fmt.Errorf("fixar contexto da transação: %w", err)
	}
	return nil
}

// EmTransacao informa se o contexto já carrega uma transação.
func EmTransacao(ctx context.Context) bool {
	_, ok := ctx.Value(chaveTx{}).(pgx.Tx)
	return ok
}

// Tx recupera a transação do contexto. Devolve erro quando o chamador está
// fora de NaEmpresa, o que é falha de programação, não condição de execução.
func Tx(ctx context.Context) (pgx.Tx, error) {
	tx, ok := ctx.Value(chaveTx{}).(pgx.Tx)
	if !ok {
		return nil, erros.SemContextoEmpresa.Com("acesso a dado de cliente fora de NaEmpresa")
	}
	return tx, nil
}

// Q devolve as consultas geradas pelo sqlc ligadas à transação do contexto.
// É a única forma de obter consultas: não existe Queries sobre o pool.
func Q(ctx context.Context) (*consultas.Queries, error) {
	tx, err := Tx(ctx)
	if err != nil {
		return nil, err
	}
	return consultas.New(tx), nil
}

func textoUUID(id *uuid.UUID) string {
	if id == nil {
		return ""
	}
	return id.String()
}

func simNao(v bool) string {
	if v {
		return "sim"
	}
	return "nao"
}
