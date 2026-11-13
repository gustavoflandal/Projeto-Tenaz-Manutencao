package rotas

import (
	"encoding/hex"
	"encoding/json"
	"time"

	"github.com/google/uuid"

	"tenaz/internal/modulos/adm"
	"tenaz/internal/plataforma/api"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/dinheiro"
)

// Tipos do contrato do ADM. Campos em snake_case, iguais às colunas do banco
// (seção 2.2 do documento 3); campo sem valor sai null, não é omitido;
// dinheiro é texto decimal; data-hora em UTC com Z.

// ----------------------------------------------------------------- empresa

// Empresa é a empresa da credencial (GET /v1/empresa). É a única resposta
// com o identificador da empresa (seção 2.2).
type Empresa struct {
	ID                    uuid.UUID `json:"id" format:"uuid"`
	RazaoSocial           string    `json:"razao_social"`
	NomeFantasia          *string   `json:"nome_fantasia"`
	Pais                  string    `json:"pais" enum:"BR,AR,PY" doc:"Define a moeda, o idioma padrão sugerido e o tipo de identificação; só leitura depois da criação da conta"`
	Moeda                 string    `json:"moeda" enum:"BRL,ARS,PYG" doc:"Moeda da empresa, dada pelo país; só leitura"`
	IdiomaPadrao          string    `json:"idioma_padrao" enum:"pt-BR,es-419" doc:"Idioma que o usuário convidado recebe; cada usuário muda o seu em PATCH /v1/usuarios/eu"`
	IdentificacaoTipo     *string   `json:"identificacao_tipo" doc:"cnpj no BR, cuit na AR, ruc no PY"`
	Identificacao         *string   `json:"identificacao" doc:"Normalizada: maiúsculas, só letras e dígitos, sem máscara"`
	Plano                 *string   `json:"plano" doc:"essencial, profissional ou enterprise; nulo na conta só do plano Serviços"`
	SrvHabilitado         bool      `json:"srv_habilitado"`
	MFAObrigatorio        bool      `json:"mfa_obrigatorio" doc:"Segundo fator exigido de todos os usuários (RN-ADM-027): quem não o tem entra com sessão restrita à ativação"`
	Situacao              string    `json:"situacao" enum:"teste,ativa,suspensa,cancelada"`
	FusoHorario           string    `json:"fuso_horario" example:"America/Sao_Paulo"`
	RetencaoAuditoriaAnos int32     `json:"retencao_auditoria_anos"`
	TesteAte              *string   `json:"teste_ate" format:"date"`
	Versao                int64     `json:"versao"`
	CriadoEm              time.Time `json:"criado_em"`
	AtualizadoEm          time.Time `json:"atualizado_em"`
}

func paraEmpresa(e adm.Empresa) Empresa {
	return Empresa{
		ID: e.ID, RazaoSocial: e.RazaoSocial, NomeFantasia: e.NomeFantasia, Pais: e.Pais, Moeda: e.Moeda,
		IdiomaPadrao: e.IdiomaPadrao, IdentificacaoTipo: e.IdentificacaoTipo, Identificacao: e.Identificacao, Plano: e.Plano,
		SrvHabilitado: e.SrvHabilitado, MFAObrigatorio: e.MFAObrigatorio, Situacao: e.Situacao, FusoHorario: e.FusoHorario,
		RetencaoAuditoriaAnos: e.RetencaoAuditoriaAnos, TesteAte: e.TesteAte, Versao: e.Versao,
		CriadoEm: api.Instante(e.CriadoEm), AtualizadoEm: api.Instante(e.AtualizadoEm),
	}
}

// PatchEmpresa são os campos que PATCH /v1/empresa altera.
type PatchEmpresa struct {
	RazaoSocial  *string `json:"razao_social,omitempty" minLength:"1" maxLength:"200"`
	NomeFantasia *string `json:"nome_fantasia,omitempty" maxLength:"200" nullable:"true"`
	// País e moeda não mudam pela API: enviá-los é 422 campo_somente_leitura.
	IdentificacaoTipo *string `json:"identificacao_tipo,omitempty" enum:"cnpj,cuit,ruc" doc:"Só o tipo do país da empresa"`
	Identificacao     *string `json:"identificacao,omitempty" maxLength:"24" doc:"Com ou sem máscara; a forma e, no CNPJ, o dígito verificador são conferidos"`
	IdiomaPadrao      *string `json:"idioma_padrao,omitempty" enum:"pt-BR,es-419"`
	FusoHorario       *string `json:"fuso_horario,omitempty" maxLength:"64"`
	// MFAObrigatorio liga ou desliga o segundo fator obrigatório (RN-ADM-027).
	MFAObrigatorio *bool `json:"mfa_obrigatorio,omitempty" doc:"Exige segundo fator de todos os usuários da empresa"`
}

type camposEmpresa struct {
	RazaoSocial       string  `json:"razao_social"`
	NomeFantasia      *string `json:"nome_fantasia"`
	IdentificacaoTipo string  `json:"identificacao_tipo"`
	Identificacao     string  `json:"identificacao"`
	IdiomaPadrao      string  `json:"idioma_padrao"`
	FusoHorario       string  `json:"fuso_horario"`
	MFAObrigatorio    bool    `json:"mfa_obrigatorio"`
}

// Responsavel é quem cria a conta: nasce administrador convidado.
type Responsavel struct {
	Nome     string  `json:"nome" minLength:"1" maxLength:"200"`
	Email    string  `json:"email" maxLength:"254"`
	Telefone *string `json:"telefone,omitempty" maxLength:"30"`
}

// NovaEmpresa é o corpo de POST /v1/empresas (cadastro pelo site).
type NovaEmpresa struct {
	Pais              string         `json:"pais" enum:"BR,AR,PY" doc:"Define a moeda (BRL, ARS, PYG), o idioma padrão sugerido e o tipo de identificação; não muda depois"`
	IdentificacaoTipo string         `json:"identificacao_tipo" enum:"cnpj,cuit,ruc" doc:"Só o tipo de pessoa jurídica do país: cnpj no BR, cuit na AR, ruc no PY"`
	Identificacao     string         `json:"identificacao" maxLength:"24" doc:"Com ou sem máscara; a forma e, no CNPJ, o dígito verificador são conferidos"`
	IdiomaPadrao      *string        `json:"idioma_padrao,omitempty" enum:"pt-BR,es-419" doc:"Ausente: o do país (pt-BR no BR, es-419 na AR e no PY)"`
	RazaoSocial       string         `json:"razao_social" minLength:"1" maxLength:"200"`
	NomeFantasia      *string        `json:"nome_fantasia,omitempty" maxLength:"200"`
	FusoHorario       string         `json:"fuso_horario,omitempty" maxLength:"64" doc:"IANA; ausente: o do país (America/Sao_Paulo, America/Argentina/Buenos_Aires, America/Asuncion)"`
	Plano             *string        `json:"plano,omitempty" enum:"essencial,profissional,enterprise" doc:"Ausente só com srv_habilitado (conta do plano Serviços)"`
	SrvHabilitado     bool           `json:"srv_habilitado,omitempty" doc:"Conta de subcontratada, com o módulo SRV (RF-SRV-001)"`
	Responsavel       Responsavel    `json:"responsavel"`
	Unidade           EntradaUnidade `json:"unidade" doc:"Unidade inicial; sem turnos, segunda a sexta das 08:00 às 17:00"`
}

// ContaCriada é a resposta de POST /v1/empresas. Nenhum identificador sai
// para quem ainda não tem credencial.
type ContaCriada struct {
	ConviteEnviado bool `json:"convite_enviado" doc:"Falso: a conta existe, mas o e-mail do convite não saiu; peça o reenvio ao suporte"`
}

// ----------------------------------------------------------------- unidade

// Endereco da unidade.
type Endereco struct {
	Logradouro   string `json:"logradouro,omitempty" maxLength:"200"`
	Numero       string `json:"numero,omitempty" maxLength:"20"`
	Complemento  string `json:"complemento,omitempty" maxLength:"100"`
	Bairro       string `json:"bairro,omitempty" maxLength:"100"`
	Municipio    string `json:"municipio,omitempty" maxLength:"100"`
	Subdivisao   string `json:"subdivisao,omitempty" maxLength:"6" doc:"Código ISO 3166-2 do país da empresa (BR-SP, AR-C, PY-ASU)"`
	CodigoPostal string `json:"codigo_postal,omitempty" maxLength:"16" doc:"CEP (BR), CPA ou código de 4 dígitos (AR), 6 dígitos (PY); recebido com ou sem máscara, devolvido normalizado"`
}

// PatchEndereco é o endereço no PATCH, que aceita null para apagar.
type PatchEndereco struct {
	_            struct{} `nullable:"true"`
	Logradouro   string   `json:"logradouro,omitempty" maxLength:"200" nullable:"true"`
	Numero       string   `json:"numero,omitempty" maxLength:"20" nullable:"true"`
	Complemento  string   `json:"complemento,omitempty" maxLength:"100" nullable:"true"`
	Bairro       string   `json:"bairro,omitempty" maxLength:"100" nullable:"true"`
	Municipio    string   `json:"municipio,omitempty" maxLength:"100" nullable:"true"`
	Subdivisao   string   `json:"subdivisao,omitempty" maxLength:"6" nullable:"true"`
	CodigoPostal string   `json:"codigo_postal,omitempty" maxLength:"16" nullable:"true"`
}

func paraEndereco(e *adm.Endereco) *Endereco {
	if e == nil {
		return nil
	}
	v := Endereco(*e)
	return &v
}

func deEndereco(e *Endereco) *adm.Endereco {
	if e == nil {
		return nil
	}
	v := adm.Endereco(*e)
	return &v
}

// Turno é um turno da unidade por dia da semana.
type Turno struct {
	ID         uuid.UUID `json:"id" format:"uuid"`
	UnidadeID  uuid.UUID `json:"unidade_id" format:"uuid"`
	DiaSemana  int32     `json:"dia_semana" minimum:"1" maximum:"7" doc:"1 = segunda ... 7 = domingo"`
	HoraInicio string    `json:"hora_inicio" example:"08:00"`
	HoraFim    string    `json:"hora_fim" example:"17:00" doc:"Menor que o início termina no dia seguinte; 24:00 é meia-noite"`
	Nome       *string   `json:"nome"`
	Versao     int64     `json:"versao"`
}

func paraTurno(t adm.Turno) Turno {
	return Turno{ID: t.ID, UnidadeID: t.UnidadeID, DiaSemana: t.DiaSemana, HoraInicio: t.HoraInicio, HoraFim: t.HoraFim, Nome: t.Nome, Versao: t.Versao}
}

// EntradaTurno cria um turno.
type EntradaTurno struct {
	DiaSemana  int32   `json:"dia_semana" minimum:"1" maximum:"7" doc:"1 = segunda ... 7 = domingo"`
	HoraInicio string  `json:"hora_inicio" maxLength:"5" example:"08:00"`
	HoraFim    string  `json:"hora_fim" maxLength:"5" example:"17:00"`
	Nome       *string `json:"nome,omitempty" maxLength:"60"`
}

func deTurno(e EntradaTurno) adm.EntradaTurno {
	return adm.EntradaTurno{DiaSemana: e.DiaSemana, HoraInicio: e.HoraInicio, HoraFim: e.HoraFim, Nome: e.Nome}
}

// PatchTurno altera o turno.
type PatchTurno struct {
	DiaSemana  *int32  `json:"dia_semana,omitempty" minimum:"1" maximum:"7"`
	HoraInicio *string `json:"hora_inicio,omitempty" maxLength:"5"`
	HoraFim    *string `json:"hora_fim,omitempty" maxLength:"5"`
	Nome       *string `json:"nome,omitempty" maxLength:"60" nullable:"true"`
}

type camposTurno struct {
	DiaSemana  int32   `json:"dia_semana"`
	HoraInicio string  `json:"hora_inicio"`
	HoraFim    string  `json:"hora_fim"`
	Nome       *string `json:"nome"`
}

// Feriado da unidade.
type Feriado struct {
	ID        uuid.UUID `json:"id" format:"uuid"`
	UnidadeID uuid.UUID `json:"unidade_id" format:"uuid"`
	Data      string    `json:"data" format:"date"`
	Descricao string    `json:"descricao"`
	Versao    int64     `json:"versao"`
}

func paraFeriado(f adm.Feriado) Feriado {
	return Feriado{ID: f.ID, UnidadeID: f.UnidadeID, Data: f.Data, Descricao: f.Descricao, Versao: f.Versao}
}

// EntradaFeriado cria um feriado.
type EntradaFeriado struct {
	Data      string `json:"data" format:"date"`
	Descricao string `json:"descricao" minLength:"1" maxLength:"100"`
}

// Unidade no detalhe, com os turnos.
type Unidade struct {
	ResumoUnidade
	Turnos []Turno `json:"turnos"`
}

// ResumoUnidade é a linha da lista.
type ResumoUnidade struct {
	ID                uuid.UUID `json:"id" format:"uuid"`
	Codigo            string    `json:"codigo"`
	Nome              string    `json:"nome"`
	IdentificacaoTipo *string   `json:"identificacao_tipo" doc:"Só o tipo do país da empresa"`
	Identificacao     *string   `json:"identificacao"`
	Endereco          *Endereco `json:"endereco"`
	FusoHorario       *string   `json:"fuso_horario" doc:"Nulo: herda o da empresa"`
	Latitude          *float64  `json:"latitude"`
	Longitude         *float64  `json:"longitude"`
	Situacao          string    `json:"situacao" enum:"ativa,inativa"`
	Versao            int64     `json:"versao"`
	CriadoEm          time.Time `json:"criado_em"`
	AtualizadoEm      time.Time `json:"atualizado_em"`
}

func paraResumoUnidade(u adm.Unidade) ResumoUnidade {
	return ResumoUnidade{
		ID: u.ID, Codigo: u.Codigo, Nome: u.Nome, IdentificacaoTipo: u.IdentificacaoTipo, Identificacao: u.Identificacao,
		Endereco: paraEndereco(u.Endereco), FusoHorario: u.FusoHorario,
		Latitude: u.Latitude, Longitude: u.Longitude, Situacao: u.Situacao, Versao: u.Versao,
		CriadoEm: api.Instante(u.CriadoEm), AtualizadoEm: api.Instante(u.AtualizadoEm),
	}
}

func paraUnidade(u adm.Unidade) Unidade {
	out := Unidade{ResumoUnidade: paraResumoUnidade(u), Turnos: []Turno{}}
	for _, t := range u.Turnos {
		out.Turnos = append(out.Turnos, paraTurno(t))
	}
	return out
}

// EntradaUnidade cria a unidade.
type EntradaUnidade struct {
	Codigo            string         `json:"codigo" minLength:"1" maxLength:"20"`
	Nome              string         `json:"nome" minLength:"1" maxLength:"200"`
	IdentificacaoTipo *string        `json:"identificacao_tipo,omitempty" enum:"cnpj,cuit,ruc" doc:"Só o tipo do país da empresa; vai junto com identificacao"`
	Identificacao     *string        `json:"identificacao,omitempty" maxLength:"24"`
	Endereco          *Endereco      `json:"endereco,omitempty"`
	FusoHorario       *string        `json:"fuso_horario,omitempty" maxLength:"64"`
	Latitude          *float64       `json:"latitude,omitempty" minimum:"-90" maximum:"90"`
	Longitude         *float64       `json:"longitude,omitempty" minimum:"-180" maximum:"180"`
	Turnos            []EntradaTurno `json:"turnos,omitempty" maxItems:"70" doc:"Só na criação; depois, pelas rotas de turno. Obrigatório em POST /v1/unidades; no cadastro da conta, ausente vale segunda a sexta das 08:00 às 17:00"`
}

func deUnidade(e EntradaUnidade) adm.EntradaUnidade {
	u := adm.EntradaUnidade{
		Codigo: e.Codigo, Nome: e.Nome, IdentificacaoTipo: e.IdentificacaoTipo, Identificacao: e.Identificacao,
		Endereco: deEndereco(e.Endereco), FusoHorario: e.FusoHorario,
		Latitude: e.Latitude, Longitude: e.Longitude,
	}
	for _, t := range e.Turnos {
		u.Turnos = append(u.Turnos, deTurno(t))
	}
	return u
}

// PatchUnidade altera a unidade. Turnos e feriados têm rotas próprias.
type PatchUnidade struct {
	Codigo            *string        `json:"codigo,omitempty" minLength:"1" maxLength:"20"`
	Nome              *string        `json:"nome,omitempty" minLength:"1" maxLength:"200"`
	IdentificacaoTipo *string        `json:"identificacao_tipo,omitempty" nullable:"true" doc:"cnpj, cuit ou ruc: o do país da empresa; vai junto com identificacao"`
	Identificacao     *string        `json:"identificacao,omitempty" maxLength:"24" nullable:"true"`
	Endereco          *PatchEndereco `json:"endereco,omitempty" doc:"Mesclado campo a campo; null apaga o endereço"`
	FusoHorario       *string        `json:"fuso_horario,omitempty" maxLength:"64" nullable:"true"`
	Latitude          *float64       `json:"latitude,omitempty" minimum:"-90" maximum:"90" nullable:"true"`
	Longitude         *float64       `json:"longitude,omitempty" minimum:"-180" maximum:"180" nullable:"true"`
}

type camposUnidade struct {
	Codigo            string    `json:"codigo"`
	Nome              string    `json:"nome"`
	IdentificacaoTipo *string   `json:"identificacao_tipo"`
	Identificacao     *string   `json:"identificacao"`
	Endereco          *Endereco `json:"endereco"`
	FusoHorario       *string   `json:"fuso_horario"`
	Latitude          *float64  `json:"latitude"`
	Longitude         *float64  `json:"longitude"`
}

// ------------------------------------------------------------------ pessoa

// Pessoa no detalhe, com a identificação inteira, que é dado sensível (a leitura
// é registrada na trilha com sensivel).
type Pessoa struct {
	ID                   uuid.UUID  `json:"id" format:"uuid"`
	Nome                 string     `json:"nome"`
	IdentificacaoTipo    *string    `json:"identificacao_tipo" doc:"cpf, cuil, dni ou ci"`
	Identificacao        *string    `json:"identificacao" doc:"Normalizada, sem máscara; dado sensível"`
	DataNascimento       *string    `json:"data_nascimento" format:"date"`
	Matricula            *string    `json:"matricula"`
	Cargo                *string    `json:"cargo"`
	Telefone             *string    `json:"telefone"`
	Email                *string    `json:"email"`
	Vinculo              string     `json:"vinculo" enum:"empregado,terceiro,autonomo"`
	FornecedorID         *uuid.UUID `json:"fornecedor_id" format:"uuid"`
	UnidadeID            *uuid.UUID `json:"unidade_id" format:"uuid"`
	RegistroProfissional *string    `json:"registro_profissional"`
	Situacao             string     `json:"situacao" enum:"ativa,inativa"`
	UsuarioID            *uuid.UUID `json:"usuario_id" format:"uuid" doc:"Usuário ligado à pessoa, se houver"`
	Versao               int64      `json:"versao"`
	CriadoEm             time.Time  `json:"criado_em"`
	AtualizadoEm         time.Time  `json:"atualizado_em"`
}

func paraPessoa(p adm.Pessoa) Pessoa {
	return Pessoa{
		ID: p.ID, Nome: p.Nome, IdentificacaoTipo: p.IdentificacaoTipo, Identificacao: p.Identificacao,
		DataNascimento: p.DataNascimento, Matricula: p.Matricula, Cargo: p.Cargo,
		Telefone: p.Telefone, Email: p.Email, Vinculo: p.Vinculo, FornecedorID: p.FornecedorID, UnidadeID: p.UnidadeID,
		RegistroProfissional: p.RegistroProfissional, Situacao: p.Situacao, UsuarioID: p.UsuarioID, Versao: p.Versao,
		CriadoEm: api.Instante(p.CriadoEm), AtualizadoEm: api.Instante(p.AtualizadoEm),
	}
}

// ResumoPessoa é a linha da lista, com a identificação mascarada.
type ResumoPessoa struct {
	ID                     uuid.UUID  `json:"id" format:"uuid"`
	Nome                   string     `json:"nome"`
	IdentificacaoTipo      *string    `json:"identificacao_tipo"`
	IdentificacaoMascarada *string    `json:"identificacao_mascarada" example:"***.982.247-**" doc:"CPF mostra o miolo; os demais tipos, os 3 últimos caracteres"`
	Matricula              *string    `json:"matricula"`
	Cargo                  *string    `json:"cargo"`
	Vinculo                string     `json:"vinculo" enum:"empregado,terceiro,autonomo"`
	FornecedorID           *uuid.UUID `json:"fornecedor_id" format:"uuid"`
	UnidadeID              *uuid.UUID `json:"unidade_id" format:"uuid"`
	Situacao               string     `json:"situacao" enum:"ativa,inativa"`
	UsuarioID              *uuid.UUID `json:"usuario_id" format:"uuid"`
	Versao                 int64      `json:"versao"`
	CriadoEm               time.Time  `json:"criado_em"`
}

func paraResumoPessoa(p adm.ResumoPessoa) ResumoPessoa {
	return ResumoPessoa{
		ID: p.ID, Nome: p.Nome, IdentificacaoTipo: p.IdentificacaoTipo, IdentificacaoMascarada: p.IdentificacaoMascarada, Matricula: p.Matricula, Cargo: p.Cargo, Vinculo: p.Vinculo,
		FornecedorID: p.FornecedorID, UnidadeID: p.UnidadeID, Situacao: p.Situacao, UsuarioID: p.UsuarioID, Versao: p.Versao,
		CriadoEm: api.Instante(p.CriadoEm),
	}
}

// EntradaPessoa cria a pessoa.
type EntradaPessoa struct {
	Nome                 string     `json:"nome" minLength:"1" maxLength:"200"`
	IdentificacaoTipo    *string    `json:"identificacao_tipo,omitempty" enum:"cpf,cuil,dni,ci" doc:"Qualquer país suportado; vai junto com identificacao"`
	Identificacao        *string    `json:"identificacao,omitempty" maxLength:"24" doc:"Com ou sem máscara; a forma e, no CPF, o dígito verificador são conferidos"`
	DataNascimento       *string    `json:"data_nascimento,omitempty" format:"date"`
	Matricula            *string    `json:"matricula,omitempty" maxLength:"40"`
	Cargo                *string    `json:"cargo,omitempty" maxLength:"100"`
	Telefone             *string    `json:"telefone,omitempty" maxLength:"30"`
	Email                *string    `json:"email,omitempty" maxLength:"254"`
	Vinculo              string     `json:"vinculo,omitempty" enum:"empregado,terceiro,autonomo" doc:"Ausente: empregado"`
	FornecedorID         *uuid.UUID `json:"fornecedor_id,omitempty" format:"uuid"`
	UnidadeID            *uuid.UUID `json:"unidade_id,omitempty" format:"uuid"`
	RegistroProfissional *string    `json:"registro_profissional,omitempty" maxLength:"60"`
}

func dePessoa(e EntradaPessoa) adm.EntradaPessoa {
	return adm.EntradaPessoa{
		Nome: e.Nome, IdentificacaoTipo: e.IdentificacaoTipo, Identificacao: e.Identificacao, DataNascimento: e.DataNascimento, Matricula: e.Matricula, Cargo: e.Cargo, Telefone: e.Telefone,
		Email: e.Email, Vinculo: e.Vinculo, FornecedorID: e.FornecedorID, UnidadeID: e.UnidadeID,
		RegistroProfissional: e.RegistroProfissional,
	}
}

// PatchPessoa altera a pessoa.
type PatchPessoa struct {
	Nome                 *string    `json:"nome,omitempty" minLength:"1" maxLength:"200"`
	IdentificacaoTipo    *string    `json:"identificacao_tipo,omitempty" nullable:"true" doc:"cpf, cuil, dni ou ci; vai junto com identificacao"`
	Identificacao        *string    `json:"identificacao,omitempty" maxLength:"24" nullable:"true"`
	DataNascimento       *string    `json:"data_nascimento,omitempty" format:"date" nullable:"true"`
	Matricula            *string    `json:"matricula,omitempty" maxLength:"40" nullable:"true"`
	Cargo                *string    `json:"cargo,omitempty" maxLength:"100" nullable:"true"`
	Telefone             *string    `json:"telefone,omitempty" maxLength:"30" nullable:"true"`
	Email                *string    `json:"email,omitempty" maxLength:"254" nullable:"true"`
	Vinculo              *string    `json:"vinculo,omitempty" enum:"empregado,terceiro,autonomo"`
	FornecedorID         *uuid.UUID `json:"fornecedor_id,omitempty" format:"uuid" nullable:"true"`
	UnidadeID            *uuid.UUID `json:"unidade_id,omitempty" format:"uuid" nullable:"true"`
	RegistroProfissional *string    `json:"registro_profissional,omitempty" maxLength:"60" nullable:"true"`
}

type camposPessoa struct {
	Nome                 string     `json:"nome"`
	IdentificacaoTipo    *string    `json:"identificacao_tipo"`
	Identificacao        *string    `json:"identificacao"`
	DataNascimento       *string    `json:"data_nascimento"`
	Matricula            *string    `json:"matricula"`
	Cargo                *string    `json:"cargo"`
	Telefone             *string    `json:"telefone"`
	Email                *string    `json:"email"`
	Vinculo              string     `json:"vinculo"`
	FornecedorID         *uuid.UUID `json:"fornecedor_id"`
	UnidadeID            *uuid.UUID `json:"unidade_id"`
	RegistroProfissional *string    `json:"registro_profissional"`
}

// ----------------------------------------------------------------- usuário

// PerfilAtribuido é um perfil do usuário.
type PerfilAtribuido struct {
	PerfilID     uuid.UUID  `json:"perfil_id" format:"uuid"`
	Codigo       string     `json:"codigo"`
	Nome         string     `json:"nome"`
	PerfilOrigem string     `json:"perfil_origem" doc:"Perfil padrão de origem"`
	UnidadeID    *uuid.UUID `json:"unidade_id" format:"uuid"`
}

// Usuario no detalhe.
type Usuario struct {
	ID                uuid.UUID         `json:"id" format:"uuid"`
	PessoaID          *uuid.UUID        `json:"pessoa_id" format:"uuid"`
	FornecedorID      *uuid.UUID        `json:"fornecedor_id" format:"uuid"`
	LocatarioID       *uuid.UUID        `json:"locatario_id" format:"uuid"`
	EscopoAcesso      string            `json:"escopo_acesso" enum:"interno,prestador,locatario"`
	Nome              string            `json:"nome"`
	Email             *string           `json:"email"`
	Telefone          *string           `json:"telefone"`
	CustoHora         *string           `json:"custo_hora" pattern:"^(0|[1-9][0-9]{0,11})(\\.[0-9]{1,6})?$" doc:"Custo por hora: preço unitário em texto decimal, com ponto, no mínimo as casas da moeda e até 6; sem sinal" example:"125.50"`
	Pais              string            `json:"pais" enum:"BR,AR,PY" doc:"País da empresa; só leitura. A interface lê este campo e nunca o deriva da moeda"`
	Moeda             string            `json:"moeda" enum:"BRL,ARS,PYG" doc:"Moeda de custo_hora: a da empresa; só leitura"`
	Idioma            string            `json:"idioma" enum:"pt-BR,es-419" doc:"Idioma do usuário; só ele o altera, em PATCH /v1/usuarios/eu"`
	UnidadePadraoID   *uuid.UUID        `json:"unidade_padrao_id" format:"uuid"`
	Situacao          string            `json:"situacao" enum:"convidado,ativo,inativo"`
	Perfis            []PerfilAtribuido `json:"perfis"`
	PermissoesAvulsas []string          `json:"permissoes_avulsas" doc:"Permissões concedidas a este usuário além dos perfis, como os.encerrar"`
	Versao            int64             `json:"versao"`
	CriadoEm          time.Time         `json:"criado_em"`
	AtualizadoEm      time.Time         `json:"atualizado_em"`
}

// precoTexto é o texto canônico do preço unitário (documento 5, Parte F, item 4):
// no mínimo as casas da moeda e no máximo seis, sem zeros à direita além delas.
func precoTexto(p *dinheiro.PrecoUnitario) *string {
	if p == nil {
		return nil
	}
	s := p.String()
	return &s
}

func paraUsuario(u adm.Usuario) Usuario {
	out := Usuario{
		ID: u.ID, PessoaID: u.PessoaID, FornecedorID: u.FornecedorID, LocatarioID: u.LocatarioID, EscopoAcesso: u.EscopoAcesso,
		Nome: u.Nome, Email: u.Email, Telefone: u.Telefone, CustoHora: precoTexto(u.CustoHora),
		Pais: u.Pais, Moeda: u.Moeda, Idioma: u.Idioma, UnidadePadraoID: u.UnidadePadraoID, Situacao: u.Situacao, Perfis: []PerfilAtribuido{},
		PermissoesAvulsas: append([]string{}, u.PermissoesAvulsas...), Versao: u.Versao,
		CriadoEm: api.Instante(u.CriadoEm), AtualizadoEm: api.Instante(u.AtualizadoEm),
	}
	for _, p := range u.Perfis {
		out.Perfis = append(out.Perfis, PerfilAtribuido{PerfilID: p.PerfilID, Codigo: p.Codigo, Nome: p.Nome, PerfilOrigem: p.CodigoOrigem, UnidadeID: p.UnidadeID})
	}
	return out
}

// UsuarioCriado é o usuário com a informação do convite.
type UsuarioCriado struct {
	Usuario
	ConviteEnviado bool `json:"convite_enviado" doc:"Falso: o usuário existe, mas o e-mail não saiu; use reenviar-convite"`
}

// ResumoUsuario é a linha da lista.
type ResumoUsuario struct {
	ID           uuid.UUID  `json:"id" format:"uuid"`
	Nome         string     `json:"nome"`
	Email        *string    `json:"email"`
	Telefone     *string    `json:"telefone"`
	PessoaID     *uuid.UUID `json:"pessoa_id" format:"uuid"`
	FornecedorID *uuid.UUID `json:"fornecedor_id" format:"uuid"`
	LocatarioID  *uuid.UUID `json:"locatario_id" format:"uuid"`
	Situacao     string     `json:"situacao" enum:"convidado,ativo,inativo"`
	Perfis       []string   `json:"perfis" doc:"Perfis de origem"`
	Versao       int64      `json:"versao"`
	CriadoEm     time.Time  `json:"criado_em"`
}

func paraResumoUsuario(u adm.ResumoUsuario) ResumoUsuario {
	return ResumoUsuario{
		ID: u.ID, Nome: u.Nome, Email: u.Email, Telefone: u.Telefone, PessoaID: u.PessoaID, FornecedorID: u.FornecedorID,
		LocatarioID: u.LocatarioID, Situacao: u.Situacao, Perfis: append([]string{}, u.Perfis...), Versao: u.Versao,
		CriadoEm: api.Instante(u.CriadoEm),
	}
}

// EntradaUsuario cria o usuário e envia o convite.
type EntradaUsuario struct {
	Nome            string         `json:"nome" minLength:"1" maxLength:"200"`
	Email           string         `json:"email" maxLength:"254"`
	Telefone        *string        `json:"telefone,omitempty" maxLength:"30"`
	Perfis          []string       `json:"perfis" minItems:"1" maxItems:"12" doc:"Códigos dos perfis: administrador, gestor, tecnico..."`
	FornecedorID    *uuid.UUID     `json:"fornecedor_id,omitempty" format:"uuid" doc:"Obrigatório no perfil prestador"`
	LocatarioID     *uuid.UUID     `json:"locatario_id,omitempty" format:"uuid" doc:"Obrigatório no perfil locatario"`
	PessoaID        *uuid.UUID     `json:"pessoa_id,omitempty" format:"uuid" doc:"Liga a uma pessoa existente; não use junto com nova_pessoa"`
	NovaPessoa      *EntradaPessoa `json:"nova_pessoa,omitempty" doc:"Cria a pessoa junto; não use junto com pessoa_id"`
	CustoHora       *string        `json:"custo_hora,omitempty" pattern:"^(0|[1-9][0-9]{0,11})(\\.[0-9]{1,6})?$" doc:"Preço unitário em texto decimal, até 6 casas, sem sinal; recusado, sem arredondar, o que a moeda da empresa não representa" example:"125.50"`
	UnidadePadraoID *uuid.UUID     `json:"unidade_padrao_id,omitempty" format:"uuid"`
}

// PatchUsuario altera o usuário. E-mail e vínculo externo não mudam.
type PatchUsuario struct {
	Nome            *string    `json:"nome,omitempty" minLength:"1" maxLength:"200"`
	Telefone        *string    `json:"telefone,omitempty" maxLength:"30" nullable:"true"`
	CustoHora       *string    `json:"custo_hora,omitempty" pattern:"^(0|[1-9][0-9]{0,11})(\\.[0-9]{1,6})?$" nullable:"true"`
	UnidadePadraoID *uuid.UUID `json:"unidade_padrao_id,omitempty" format:"uuid" nullable:"true"`
	PessoaID        *uuid.UUID `json:"pessoa_id,omitempty" format:"uuid" nullable:"true"`
	Perfis          []string   `json:"perfis,omitempty" minItems:"1" maxItems:"12" doc:"Lista completa dos perfis; substitui a atual"`
}

type camposUsuario struct {
	Nome            string     `json:"nome"`
	Telefone        *string    `json:"telefone"`
	CustoHora       *string    `json:"custo_hora"`
	UnidadePadraoID *uuid.UUID `json:"unidade_padrao_id"`
	PessoaID        *uuid.UUID `json:"pessoa_id"`
}

// PatchEu altera os dados do próprio usuário.
type PatchEu struct {
	Nome     *string `json:"nome,omitempty" minLength:"1" maxLength:"200"`
	Telefone *string `json:"telefone,omitempty" maxLength:"30" nullable:"true"`
	Idioma   *string `json:"idioma,omitempty" enum:"pt-BR,es-419" doc:"Só o próprio usuário altera o idioma"`
}

type camposEu struct {
	Nome     string  `json:"nome"`
	Telefone *string `json:"telefone"`
	Idioma   string  `json:"idioma"`
}

// EntradaPermissao concede uma permissão avulsa.
type EntradaPermissao struct {
	Permissao string `json:"permissao" maxLength:"60" example:"os.encerrar" doc:"Código do catálogo de RF-ADM-005"`
}

// SessoesEncerradas é a resposta de encerrar-sessoes.
type SessoesEncerradas struct {
	SessoesEncerradas int32 `json:"sessoes_encerradas"`
}

// ------------------------------------------------------ perfis e permissões

// Perfil da empresa.
type Perfil struct {
	ID           uuid.UUID `json:"id" format:"uuid"`
	Codigo       string    `json:"codigo"`
	Nome         string    `json:"nome"`
	Padrao       bool      `json:"padrao" doc:"Perfil padrão do sumário, que não se altera"`
	PerfilOrigem *string   `json:"perfil_origem"`
	Permissoes   []string  `json:"permissoes"`
	Versao       int64     `json:"versao"`
}

func paraPerfil(p adm.Perfil) Perfil {
	return Perfil{ID: p.ID, Codigo: p.Codigo, Nome: p.Nome, Padrao: p.Padrao, PerfilOrigem: p.PerfilOrigem,
		Permissoes: append([]string{}, p.Permissoes...), Versao: p.Versao}
}

// Permissao é uma linha do catálogo de RF-ADM-005.
type Permissao struct {
	Codigo    string `json:"codigo" example:"os.encerrar"`
	Descricao string `json:"descricao"`
	Grupo     string `json:"grupo" doc:"Grupo da seção 4.8 do documento 3"`
	Fase      string `json:"fase" enum:"MVP,F2,F3"`
	Sensivel  bool   `json:"sensivel"`
}

func paraPermissao(p autorizacao.Permissao) Permissao {
	return Permissao{Codigo: p.Codigo, Descricao: p.Descricao, Grupo: p.Grupo, Fase: string(p.Fase), Sensivel: p.Sensivel}
}

// --------------------------------------------------------------- auditoria

// ResumoEvento é a linha da trilha de auditoria.
type ResumoEvento struct {
	ID                      uuid.UUID  `json:"id" format:"uuid"`
	OcorridoEm              time.Time  `json:"ocorrido_em"`
	AtorTipo                string     `json:"ator_tipo" enum:"usuario,sistema,integracao,suporte_tenaz,link_execucao,sincronizacao_conexao"`
	UsuarioID               *uuid.UUID `json:"usuario_id" format:"uuid"`
	UsuarioNome             *string    `json:"usuario_nome" doc:"Nome do usuário no instante do evento; a troca de nome não reescreve o histórico"`
	TokenAPIID              *uuid.UUID `json:"token_api_id" format:"uuid"`
	Acao                    string     `json:"acao" enum:"login,logout,inclusao,alteracao,exclusao,delecao,visualizacao,impressao,extracao" doc:"Tipo de operação (lista fechada, sumário 3.6)"`
	Resultado               string     `json:"resultado" enum:"sucesso,falha" doc:"falha só em login; o motivo está em depois do detalhe do evento"`
	Entidade                string     `json:"entidade"`
	EntidadeID              *uuid.UUID `json:"entidade_id" format:"uuid"`
	EntidadeChave           *string    `json:"entidade_chave"`
	CamposAlterados         []string   `json:"campos_alterados"`
	Modulo                  *string    `json:"modulo" doc:"Código do módulo (sumário, seção 2), em minúsculas"`
	Tela                    *string    `json:"tela" doc:"Tela de origem, <modulo>.<tela> do catálogo fechado; nula em API, job e sincronização"`
	Sensivel                bool       `json:"sensivel" doc:"O registro mostrado, impresso ou extraído tem dado sensível"`
	Origem                  *string    `json:"origem" enum:"web,app,api,job,importacao"`
	IP                      *string    `json:"ip"`
	DispositivoID           *uuid.UUID `json:"dispositivo_id" format:"uuid"`
	RequestID               *string    `json:"request_id"`
	ConexaoEmpresaID        *uuid.UUID `json:"conexao_empresa_id" format:"uuid"`
	OcorridoNoDispositivoEm *time.Time `json:"ocorrido_no_dispositivo_em"`
}

func paraResumoEvento(e adm.ResumoEvento) ResumoEvento {
	campos := e.CamposAlterados
	if campos == nil {
		campos = []string{}
	}
	return ResumoEvento{
		ID: e.ID, OcorridoEm: api.Instante(e.OcorridoEm), AtorTipo: e.AtorTipo, UsuarioID: e.UsuarioID, UsuarioNome: e.UsuarioNome,
		TokenAPIID: e.TokenAPIID, Acao: e.Acao, Resultado: e.Resultado, Entidade: e.Entidade, EntidadeID: e.EntidadeID,
		EntidadeChave: e.EntidadeChave, CamposAlterados: campos, Modulo: e.Modulo, Tela: e.Tela, Sensivel: e.Sensivel,
		Origem: e.Origem, IP: e.IP, DispositivoID: e.DispositivoID, RequestID: e.RequestID, ConexaoEmpresaID: e.ConexaoEmpresaID,
		OcorridoNoDispositivoEm: api.InstanteOpcional(e.OcorridoNoDispositivoEm),
	}
}

// Evento é o evento com os valores antes e depois.
type Evento struct {
	ResumoEvento
	Antes     json.RawMessage `json:"antes" doc:"Valores antes da alteração; nulo na criação"`
	Depois    json.RawMessage `json:"depois" doc:"Valores depois da alteração; nulo na exclusão. Na falha de login, traz o motivo"`
	Detalhe   json.RawMessage `json:"detalhe" doc:"Filtros, formato e quantidade de uma impressão ou extração; nunca o conteúdo mostrado"`
	Hash      string          `json:"hash" pattern:"^[0-9a-f]{64}$" doc:"sha256, em hexadecimal, do conteúdo do evento (formato 1, documento 2, seção 9.2.1)"`
	UserAgent *string         `json:"user_agent"`
}

func paraEvento(e adm.Evento) Evento {
	return Evento{ResumoEvento: paraResumoEvento(e.ResumoEvento), Antes: jsonOuNulo(e.Antes), Depois: jsonOuNulo(e.Depois),
		Detalhe: jsonOuNulo(e.Detalhe), Hash: hex.EncodeToString(e.Hash), UserAgent: e.UserAgent}
}

// DiaVerificado é o resultado da verificação de integridade de um dia UTC.
type DiaVerificado struct {
	Dia      string     `json:"dia" format:"date" doc:"Dia UTC"`
	Estado   string     `json:"estado" enum:"integro,violado,nao_selado"`
	Motivo   *string    `json:"motivo" doc:"Em violado: evento_alterado, eventos_divergentes, selo_adulterado, cadeia_quebrada, selo_ausente, evento_anterior_ao_primeiro_selo, copia_ausente ou copia_divergente. Em nao_selado: ainda_nao_selado ou anterior_ao_primeiro_selo"`
	EventoID *uuid.UUID `json:"evento_id" format:"uuid" doc:"Primeiro evento alterado, quando há como apontar"`
}

func paraDiaVerificado(d adm.DiaVerificado) DiaVerificado {
	return DiaVerificado{Dia: d.Dia.UTC().Format("2006-01-02"), Estado: d.Estado, Motivo: d.Motivo, EventoID: d.EventoID}
}

// CapturaAuditoria é a configuração de uma tela do catálogo.
type CapturaAuditoria struct {
	Tela         string `json:"tela" example:"adm.pessoas" doc:"<modulo>.<tela> do catálogo fechado"`
	Modulo       string `json:"modulo" example:"adm"`
	Nome         string `json:"nome" example:"Pessoas"`
	Visualizacao bool   `json:"visualizacao" doc:"Grava um evento por registro aberto nesta tela"`
	Impressao    bool   `json:"impressao" doc:"Grava a impressão (PDF ou tela) feita nesta tela"`
	Extracao     bool   `json:"extracao" doc:"Grava a extração (lista, relatório, exportação) feita nesta tela"`
	Configurada  bool   `json:"configurada" doc:"Falso: a empresa ainda não mudou nenhuma chave desta tela e vale o padrão (impressão e extração ligadas, visualização desligada)"`
}

func paraCapturaAuditoria(c adm.CapturaAuditoria) CapturaAuditoria {
	return CapturaAuditoria{Tela: c.Tela, Modulo: c.Modulo, Nome: c.Nome, Visualizacao: c.Visualizacao,
		Impressao: c.Impressao, Extracao: c.Extracao, Configurada: c.Configurada}
}

// EntradaCapturaAuditoria são as três chaves de uma tela: sempre as três.
type EntradaCapturaAuditoria struct {
	Visualizacao bool `json:"visualizacao"`
	Impressao    bool `json:"impressao"`
	Extracao     bool `json:"extracao"`
}

func jsonOuNulo(b []byte) json.RawMessage {
	if len(b) == 0 {
		return json.RawMessage("null")
	}
	return json.RawMessage(b)
}
