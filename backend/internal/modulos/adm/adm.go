// Package adm é a API pública do módulo ADM — Administração (seção 5.1 do
// documento 5): conta da empresa, unidades com turnos e feriados, pessoas,
// usuários, perfis, permissões avulsas e a consulta da trilha de auditoria
// (RF-ADM-001, 002, 004, 005, 011, 015; RF-SRV-001 parcial).
//
// Os outros módulos e a camada HTTP usam só este pacote; o que fica em
// adm/internal o compilador não deixa importar.
//
// Convenções desta interface:
//   - todo método recebe o Ator resolvido pela autenticação nesta requisição
//     (exceto CriarConta e AceitarConvite, que são públicos);
//   - Editar* recebe o estado completo dos campos editáveis: a camada HTTP
//     faz a fusão do PATCH com o registro lido; VersaoEsperada vem do If-Match;
//   - erros são *erros.Erro com código estável; o mapa está em cada método.
package adm

import (
	"context"
	"time"

	"github.com/google/uuid"

	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/paginacao"
)

// Servico é a porta do módulo ADM.
//
// Erros comuns a todos os métodos com Ator: permissao_negada (403) quando falta
// a permissão; recurso_nao_encontrado (404) para registro inexistente, de outra
// empresa ou fora do alcance — e para usuário externo em qualquer rota de ADM
// além de ObterEu e EditarEu; validacao_falhou (422) com Campos;
// versao_desatualizada (412) quando VersaoEsperada não confere.
type Servico interface {
	// CriarConta cria a empresa com a unidade inicial (e seus turnos), os 12
	// perfis padrão com a matriz de RF-ADM-005 e o usuário administrador
	// `convidado`, e envia o convite (que confirma o e-mail do responsável
	// antes do primeiro acesso). Público. conflito em cnpj já cadastrado.
	CriarConta(ctx context.Context, e EntradaCriarConta) (ContaCriada, error)

	// AceitarConvite liga o login ao usuário convidado e o ativa. Senha só é
	// exigida quando o e-mail ainda não tem login (ConviteValido.
	// IdentidadeExiste, de autenticacao.LocalizarConvite). Público.
	// link_email_invalido, senha_fraca, conflito (o e-mail já tem usuário na
	// empresa).
	AceitarConvite(ctx context.Context, token, senha string) (ConviteAceito, error)

	// Empresa: adm.ler para ler, adm.gerenciar_empresa para alterar.
	ObterEmpresa(ctx context.Context, a autorizacao.Ator) (Empresa, error)
	EditarEmpresa(ctx context.Context, a autorizacao.Ator, e EntradaEmpresa) (Empresa, error)

	// Unidades: adm.ler para ler; adm.gerenciar_empresa para o resto.
	// CriarUnidade: plano_nao_inclui_recurso (403, plano_necessario) na
	// segunda unidade do Essencial. InativarUnidade recusa a última unidade
	// ativa. Turnos: RemoverTurno recusa o último (validacao_falhou).
	ListarUnidades(ctx context.Context, a autorizacao.Ator) ([]Unidade, error)
	ObterUnidade(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (Unidade, error)
	CriarUnidade(ctx context.Context, a autorizacao.Ator, e EntradaUnidade) (Unidade, error)
	EditarUnidade(ctx context.Context, a autorizacao.Ator, id uuid.UUID, e EntradaUnidade) (Unidade, error)
	InativarUnidade(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (Unidade, error)
	ReativarUnidade(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (Unidade, error)
	IncluirTurno(ctx context.Context, a autorizacao.Ator, unidadeID uuid.UUID, e EntradaTurno) (Turno, error)
	AlterarTurno(ctx context.Context, a autorizacao.Ator, unidadeID, turnoID uuid.UUID, e EntradaTurno) (Turno, error)
	RemoverTurno(ctx context.Context, a autorizacao.Ator, unidadeID, turnoID uuid.UUID) error
	ListarFeriados(ctx context.Context, a autorizacao.Ator, unidadeID uuid.UUID, de, ate *string) ([]Feriado, error)
	IncluirFeriado(ctx context.Context, a autorizacao.Ator, unidadeID uuid.UUID, e EntradaFeriado) (Feriado, error)
	RemoverFeriado(ctx context.Context, a autorizacao.Ator, unidadeID, feriadoID uuid.UUID) error

	// Pessoas: adm.ler_pessoas para ler (o titular lê a própria);
	// adm.gerenciar_pessoas para o resto. A lista traz o CPF mascarado; o
	// detalhe traz o CPF inteiro e grava visualizacao_sensivel na trilha.
	// conflito em cpf repetido. Pessoa nunca é excluída: é inativada.
	ListarPessoas(ctx context.Context, a autorizacao.Ator, f FiltroPessoas) (paginacao.Pagina[ResumoPessoa], error)
	ObterPessoa(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (Pessoa, error)
	CriarPessoa(ctx context.Context, a autorizacao.Ator, e EntradaPessoa) (Pessoa, error)
	EditarPessoa(ctx context.Context, a autorizacao.Ator, id uuid.UUID, e EntradaPessoa) (Pessoa, error)
	InativarPessoa(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (Pessoa, error)
	ReativarPessoa(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (Pessoa, error)

	// Usuários: adm.ler_usuarios para ler; adm.gerenciar_usuarios para o
	// resto. Regras anti-escalada (documento 9, seção 5.4): ninguém altera os
	// próprios perfis nem permissões; só administrador altera administrador;
	// só se atribui perfil cujas permissões se tem (com as isenções do
	// administrador). ultimo_administrador (409) ao inativar ou tirar o
	// perfil do último administrador ativo. transicao_invalida (409,
	// estado_atual) em inativar inativo, reativar ativo, reenviar convite de
	// quem não é convidado.
	ListarUsuarios(ctx context.Context, a autorizacao.Ator, f FiltroUsuarios) (paginacao.Pagina[ResumoUsuario], error)
	ObterUsuario(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (Usuario, error)
	CriarUsuario(ctx context.Context, a autorizacao.Ator, e EntradaUsuario) (UsuarioCriado, error)
	EditarUsuario(ctx context.Context, a autorizacao.Ator, id uuid.UUID, e EntradaEditarUsuario) (Usuario, error)
	DefinirPerfis(ctx context.Context, a autorizacao.Ator, id uuid.UUID, perfis []string) (Usuario, error)
	InativarUsuario(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (Usuario, error)
	ReativarUsuario(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (UsuarioCriado, error)
	ReenviarConvite(ctx context.Context, a autorizacao.Ator, id uuid.UUID) error
	EncerrarSessoes(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (int32, error)
	ConcederPermissao(ctx context.Context, a autorizacao.Ator, id uuid.UUID, permissao string) (Usuario, error)
	RetirarPermissao(ctx context.Context, a autorizacao.Ator, id uuid.UUID, permissao string) (Usuario, error)

	// Próprio usuário (GET e PATCH /v1/usuarios/eu): qualquer autenticado,
	// inclusive externo. Só nome e telefone mudam (ISO-E01-01).
	ObterEu(ctx context.Context, a autorizacao.Ator) (Usuario, error)
	EditarEu(ctx context.Context, a autorizacao.Ator, e EntradaEu) (Usuario, error)

	// Perfis (adm.ler_usuarios) e catálogo de permissões (adm.gerenciar_perfis).
	ListarPerfis(ctx context.Context, a autorizacao.Ator) ([]Perfil, error)
	ListarPermissoes(ctx context.Context, a autorizacao.Ator) ([]autorizacao.Permissao, error)

	// Trilha de auditoria (adm.ler_auditoria), com filtros e cursor.
	ListarEventosAuditoria(ctx context.Context, a autorizacao.Ator, f FiltroAuditoria) (paginacao.Pagina[ResumoEvento], error)
	ObterEventoAuditoria(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (Evento, error)
}

// ------------------------------------------------------------------ conta

// EntradaCriarConta é o cadastro pelo site (RF-ADM-001, RF-SRV-001).
type EntradaCriarConta struct {
	CNPJ          string // com ou sem máscara
	RazaoSocial   string
	NomeFantasia  *string
	FusoHorario   string  // IANA; vazio: America/Sao_Paulo
	Plano         *string // essencial, profissional, enterprise; nulo só com SrvHabilitado
	SrvHabilitado bool    // módulo SRV, plano Serviços
	Responsavel   EntradaResponsavel
	Unidade       EntradaUnidade // unidade inicial; sem turnos, segunda a sexta 08:00 às 17:00
}

// EntradaResponsavel é o responsável, que nasce administrador convidado.
type EntradaResponsavel struct {
	Nome     string
	Email    string
	Telefone *string
}

// ContaCriada é o resultado de CriarConta.
type ContaCriada struct {
	EmpresaID      uuid.UUID
	UnidadeID      uuid.UUID
	UsuarioID      uuid.UUID
	ConviteEnviado bool // falso: a conta existe, o e-mail não saiu (fica no log)
}

// ConviteAceito é o resultado de AceitarConvite.
type ConviteAceito struct {
	EmpresaID      uuid.UUID
	UsuarioID      uuid.UUID
	IdentidadeNova bool // falso: o e-mail já entrava em outra empresa e usa a senha que tem
}

// ----------------------------------------------------------------- empresa

// Empresa é a empresa da credencial.
type Empresa struct {
	ID                    uuid.UUID
	RazaoSocial           string
	NomeFantasia          *string
	CNPJ                  *string
	Plano                 *string
	SrvHabilitado         bool
	MFAObrigatorio        bool   // segundo fator exigido de todos os usuários (RN-ADM-027)
	Situacao              string // teste, ativa, suspensa, cancelada
	FusoHorario           string
	RetencaoAuditoriaAnos int32
	TesteAte              *string // AAAA-MM-DD
	Versao                int64
	CriadoEm              time.Time
	AtualizadoEm          time.Time
}

// EntradaEmpresa são os dados cadastrais editáveis. Plano e SRV mudam pela
// cobrança (RF-ADM-009, E10), não aqui.
type EntradaEmpresa struct {
	RazaoSocial    string
	NomeFantasia   *string
	CNPJ           string
	FusoHorario    string
	MFAObrigatorio bool // RN-ADM-027; a troca vai para a trilha (RN-ADM-030)
	VersaoEsperada *int64
}

// ----------------------------------------------------------------- unidade

// Endereco da unidade (jsonb).
type Endereco struct {
	Logradouro  string `json:"logradouro,omitempty"`
	Numero      string `json:"numero,omitempty"`
	Complemento string `json:"complemento,omitempty"`
	Bairro      string `json:"bairro,omitempty"`
	Municipio   string `json:"municipio,omitempty"`
	UF          string `json:"uf,omitempty"`
	CEP         string `json:"cep,omitempty"`
}

// EntradaUnidade cria ou edita a unidade. Turnos só na criação; depois, pelas
// operações de turno.
type EntradaUnidade struct {
	Codigo         string
	Nome           string
	CNPJ           *string
	Endereco       *Endereco
	FusoHorario    *string // nulo herda o da empresa
	Latitude       *float64
	Longitude      *float64
	Turnos         []EntradaTurno
	VersaoEsperada *int64
}

// Unidade com o calendário.
type Unidade struct {
	ID           uuid.UUID
	Codigo       string
	Nome         string
	CNPJ         *string
	Endereco     *Endereco
	FusoHorario  *string
	Latitude     *float64
	Longitude    *float64
	Situacao     string // ativa, inativa
	Versao       int64
	CriadoEm     time.Time
	AtualizadoEm time.Time
	Turnos       []Turno // no detalhe
}

// EntradaTurno é um turno por dia da semana.
type EntradaTurno struct {
	DiaSemana  int32  // 1 = segunda ... 7 = domingo
	HoraInicio string // HH:MM
	HoraFim    string // HH:MM; menor que o início termina no dia seguinte; 24:00 = meia-noite
	Nome       *string
}

// Turno da unidade.
type Turno struct {
	ID         uuid.UUID
	UnidadeID  uuid.UUID
	DiaSemana  int32
	HoraInicio string
	HoraFim    string
	Nome       *string
	Versao     int64
}

// EntradaFeriado é um feriado da unidade.
type EntradaFeriado struct {
	Data      string // AAAA-MM-DD
	Descricao string
}

// Feriado da unidade.
type Feriado struct {
	ID        uuid.UUID
	UnidadeID uuid.UUID
	Data      string
	Descricao string
	Versao    int64
}

// ------------------------------------------------------------------ pessoa

// Vínculos e situações da pessoa (ck_pessoa_vinculo, ck_pessoa_situacao).
const (
	VinculoEmpregado = "empregado"
	VinculoTerceiro  = "terceiro"
	VinculoAutonomo  = "autonomo"
)

// EntradaPessoa cria ou edita a pessoa (RF-ADM-015).
type EntradaPessoa struct {
	Nome                 string
	CPF                  *string // com ou sem máscara
	DataNascimento       *string // AAAA-MM-DD
	Matricula            *string
	Cargo                *string
	Telefone             *string
	Email                *string
	Vinculo              string // empregado (padrão), terceiro, autonomo
	FornecedorID         *uuid.UUID
	UnidadeID            *uuid.UUID
	RegistroProfissional *string
	VersaoEsperada       *int64
}

// Pessoa no detalhe, com o CPF inteiro.
type Pessoa struct {
	ID                   uuid.UUID
	Nome                 string
	CPF                  *string
	DataNascimento       *string
	Matricula            *string
	Cargo                *string
	Telefone             *string
	Email                *string
	Vinculo              string
	FornecedorID         *uuid.UUID
	UnidadeID            *uuid.UUID
	RegistroProfissional *string
	Situacao             string // ativa, inativa (RN-ADM-020)
	UsuarioID            *uuid.UUID
	Versao               int64
	CriadoEm             time.Time
	AtualizadoEm         time.Time
}

// ResumoPessoa é a linha da lista, com o CPF mascarado.
type ResumoPessoa struct {
	ID           uuid.UUID
	Nome         string
	CPFMascarado *string
	Matricula    *string
	Cargo        *string
	Vinculo      string
	FornecedorID *uuid.UUID
	UnidadeID    *uuid.UUID
	Situacao     string
	UsuarioID    *uuid.UUID
	Versao       int64
	CriadoEm     time.Time
}

// FiltroPessoas são os filtros de GET /v1/pessoas.
type FiltroPessoas struct {
	Busca        *string // trecho do nome
	Situacao     *string
	FornecedorID *uuid.UUID
	UnidadeID    *uuid.UUID
	ComUsuario   *bool
	Cursor       string
	Limite       int32
}

// ----------------------------------------------------------------- usuário

// EntradaUsuario cria o usuário e envia o convite (RF-ADM-004).
type EntradaUsuario struct {
	Nome              string
	Email             string
	Telefone          *string
	Perfis            []string       // códigos: administrador, gestor, tecnico...
	FornecedorID      *uuid.UUID     // obrigatório no perfil prestador
	LocatarioID       *uuid.UUID     // obrigatório no perfil locatario
	PessoaID          *uuid.UUID     // liga a uma pessoa existente...
	NovaPessoa        *EntradaPessoa // ...ou cria a pessoa junto (nunca os dois)
	CustoHoraCentavos *int64
	UnidadePadraoID   *uuid.UUID
}

// EntradaEditarUsuario são os dados editáveis por quem administra. E-mail e
// vínculo externo não mudam: para trocar, inativa-se e cria-se outro.
type EntradaEditarUsuario struct {
	Nome              string
	Telefone          *string
	CustoHoraCentavos *int64
	UnidadePadraoID   *uuid.UUID
	PessoaID          *uuid.UUID
	VersaoEsperada    *int64
}

// EntradaEu é o que o próprio usuário altera.
type EntradaEu struct {
	Nome     string
	Telefone *string
}

// PerfilAtribuido é um perfil do usuário.
type PerfilAtribuido struct {
	PerfilID     uuid.UUID
	Codigo       string
	Nome         string
	CodigoOrigem string // tipo do perfil
	UnidadeID    *uuid.UUID
}

// Usuario no detalhe.
type Usuario struct {
	ID                uuid.UUID
	IdentidadeLigada  bool // o convite foi aceito
	PessoaID          *uuid.UUID
	FornecedorID      *uuid.UUID
	LocatarioID       *uuid.UUID
	EscopoAcesso      string // interno, prestador, locatario
	Nome              string
	Email             *string
	Telefone          *string
	CustoHoraCentavos *int64
	UnidadePadraoID   *uuid.UUID
	Situacao          string // convidado, ativo, inativo
	Perfis            []PerfilAtribuido
	PermissoesAvulsas []string
	Versao            int64
	CriadoEm          time.Time
	AtualizadoEm      time.Time
}

// UsuarioCriado traz o usuário e se o convite saiu.
type UsuarioCriado struct {
	Usuario        Usuario
	ConviteEnviado bool
}

// ResumoUsuario é a linha da lista.
type ResumoUsuario struct {
	ID           uuid.UUID
	Nome         string
	Email        *string
	Telefone     *string
	PessoaID     *uuid.UUID
	FornecedorID *uuid.UUID
	LocatarioID  *uuid.UUID
	Situacao     string
	Perfis       []string // códigos de origem
	Versao       int64
	CriadoEm     time.Time
}

// FiltroUsuarios são os filtros de GET /v1/usuarios.
type FiltroUsuarios struct {
	Busca    *string // trecho do nome ou do e-mail
	Situacao *string
	Perfil   *string // código de origem
	Cursor   string
	Limite   int32
}

// Perfil da empresa.
type Perfil struct {
	ID           uuid.UUID
	Codigo       string
	Nome         string
	Padrao       bool
	PerfilOrigem *string
	Permissoes   []string
	Versao       int64
}

// --------------------------------------------------------------- auditoria

// FiltroAuditoria são os filtros de GET /v1/eventos-auditoria.
type FiltroAuditoria struct {
	Entidade   *string
	EntidadeID *uuid.UUID
	UsuarioID  *uuid.UUID
	Acao       *string
	De         *time.Time // inclusive
	Ate        *time.Time // exclusive
	Cursor     string
	Limite     int32
}

// ResumoEvento é a linha da trilha, sem antes e depois.
type ResumoEvento struct {
	ID                      uuid.UUID
	OcorridoEm              time.Time
	AtorTipo                string
	UsuarioID               *uuid.UUID
	TokenAPIID              *uuid.UUID
	Acao                    string
	Entidade                string
	EntidadeID              *uuid.UUID
	EntidadeChave           *string
	CamposAlterados         []string
	Origem                  *string
	IP                      *string
	DispositivoID           *uuid.UUID
	RequestID               *string
	ConexaoEmpresaID        *uuid.UUID
	OcorridoNoDispositivoEm *time.Time
}

// Evento é o evento com os valores antes e depois (JSON cru).
type Evento struct {
	ResumoEvento
	Antes     []byte
	Depois    []byte
	UserAgent *string
}
