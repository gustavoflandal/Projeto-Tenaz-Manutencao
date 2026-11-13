package adm

import (
	"context"
	"errors"
	"math/big"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"tenaz/internal/plataforma/autenticacao"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/erros"
)

// Os casos de uso ficam neste pacote, e não em adm/internal/app como na
// seção 5.2 do documento 5: os DTOs da interface pública estão aqui, e um
// pacote app que os usasse e fosse construído daqui formaria ciclo de
// importação. Domínio e repositório continuam em adm/internal.

// ErrUltimoAdministrador: a operação deixaria a empresa sem administrador
// ativo (RF-ADM-004). Código ultimo_administrador da seção 3.2 do documento 3.
var ErrUltimoAdministrador = erros.Novo("ultimo_administrador", "A empresa ficaria sem administrador", http.StatusConflict)

// Autenticacao é o que o ADM usa do serviço de autenticação da plataforma
// (sequências da seção 6.9 do documento 5). Toda chamada acontece FORA de
// NaEmpresa, depois do COMMIT: as credenciais ficam num pool e num papel de
// banco próprios (ADR-011). A inativação de usuário não chama nada daqui: o
// trigger de acesso_empresa revoga as sessões na mesma transação.
type Autenticacao interface {
	EmitirConvite(ctx context.Context, empresaID, usuarioID uuid.UUID, email, nomeEmpresa, nomeConvidado string) error
	EncerrarSessoesUsuario(ctx context.Context, admin autorizacao.Ator, usuarioID uuid.UUID) (int, error)
	AceitarConvite(ctx context.Context, token, senha string,
		ligar func(ctx context.Context, c autenticacao.ConviteValido, identidadeID uuid.UUID) error) (autenticacao.ConviteValido, error)
}

// Dependencias do módulo.
type Dependencias struct {
	Banco        *db.DB // pool da aplicação
	Autenticacao Autenticacao
	Agora        func() time.Time
}

type servico struct {
	d Dependencias
}

// Novo monta o serviço do módulo.
func Novo(d Dependencias) (Servico, error) {
	if d.Banco == nil || d.Autenticacao == nil {
		return nil, errors.New("adm: dependência ausente")
	}
	if d.Agora == nil {
		d.Agora = time.Now
	}
	return &servico{d: d}, nil
}

func (s *servico) agora() time.Time { return s.d.Agora().UTC() }

// naEmpresa abre a transação com a empresa e o escopo do ator.
func (s *servico) naEmpresa(ctx context.Context, a autorizacao.Ator, fn func(context.Context) error) error {
	return s.d.Banco.NaEmpresa(ctx, a.EmpresaID, a.Escopo, fn)
}

// apenasInterno: as rotas de ADM não existem para usuário externo, além de
// /v1/usuarios/eu (documento 3, seção 4.9): 404, não 403.
func apenasInterno(a autorizacao.Ator, permissoes ...string) error {
	if a.Externo() {
		return erros.NaoEncontrado
	}
	return a.Exigir(permissoes...)
}

func naoEncontrado(err error) error {
	if errors.Is(err, pgx.ErrNoRows) {
		return erros.NaoEncontrado
	}
	return err
}

func conferirVersao(esperada *int64, atual int64) error {
	if esperada != nil && *esperada != atual {
		return erros.VersaoDesatualizada.Com("o registro está na versão %d", atual).ComDado("versao_atual", atual)
	}
	return nil
}

func novoID() uuid.UUID { return uuid.Must(uuid.NewV7()) }

func ptr[T any](v T) *T { return &v }

func vazioParaNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// centavos converte numeric(12,2) em centavos, sem ponto flutuante.
func centavos(n pgtype.Numeric) *int64 {
	if !n.Valid || n.Int == nil || n.NaN {
		return nil
	}
	v := new(big.Int).Set(n.Int)
	exp := n.Exp + 2
	dez := big.NewInt(10)
	for ; exp > 0; exp-- {
		v.Mul(v, dez)
	}
	for ; exp < 0; exp++ {
		v.Quo(v, dez)
	}
	if !v.IsInt64() {
		return nil
	}
	return ptr(v.Int64())
}
