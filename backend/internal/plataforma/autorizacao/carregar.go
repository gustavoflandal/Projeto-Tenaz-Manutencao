package autorizacao

import (
	"context"
	"errors"
	"fmt"
	"slices"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/db/consultas"
	"tenaz/internal/plataforma/erros"
)

// Carregar lê o usuário, os perfis, as permissões avulsas e a empresa na
// transação do contexto e monta o Ator, com o escopo do cadastro (seção 2.9.2
// do documento 2). Deve rodar em NaEmpresa com escopo interno e sem usuário —
// o único ponto em que dado de usuário externo é lido assim.
//
// Usuário inexistente, excluído, convidado ou inativo: credencial_invalida.
// SessaoID, Canal e IdentidadeID da credencial ficam com quem chama conferir.
func Carregar(ctx context.Context, empresaID, usuarioID uuid.UUID) (Ator, error) {
	q, err := db.Q(ctx)
	if err != nil {
		return Ator{}, err
	}
	u, err := q.AutzBuscarUsuario(ctx, consultas.AutzBuscarUsuarioParams{EmpresaID: empresaID, ID: usuarioID})
	if errors.Is(err, pgx.ErrNoRows) {
		return Ator{}, erros.CredencialInvalida.Tecnico("usuário não encontrado")
	}
	if err != nil {
		return Ator{}, fmt.Errorf("ler usuário: %w", err)
	}
	if u.Situacao != "ativo" || u.IdentidadeID == nil {
		return Ator{}, erros.CredencialInvalida.Tecnico("usuário não está ativo")
	}
	emp, err := q.AutzBuscarEmpresa(ctx, empresaID)
	if errors.Is(err, pgx.ErrNoRows) {
		return Ator{}, erros.CredencialInvalida.Tecnico("empresa não encontrada")
	}
	if err != nil {
		return Ator{}, fmt.Errorf("ler empresa: %w", err)
	}
	perfis, err := q.AutzPerfisDoUsuario(ctx, consultas.AutzPerfisDoUsuarioParams{EmpresaID: empresaID, UsuarioID: usuarioID})
	if err != nil {
		return Ator{}, fmt.Errorf("ler perfis do usuário: %w", err)
	}
	avulsas, err := q.AutzPermissoesAvulsas(ctx, consultas.AutzPermissoesAvulsasParams{EmpresaID: empresaID, UsuarioID: usuarioID})
	if err != nil {
		return Ator{}, fmt.Errorf("ler permissões avulsas: %w", err)
	}

	uid := usuarioID
	escopo := db.Escopo{Tipo: db.EscopoInterno, UsuarioID: &uid}
	switch {
	case u.FornecedorID != nil:
		escopo = db.Escopo{Tipo: db.EscopoPrestador, FornecedorID: u.FornecedorID, UsuarioID: &uid}
	case u.LocatarioID != nil:
		escopo = db.Escopo{Tipo: db.EscopoLocatario, LocatarioID: u.LocatarioID, UsuarioID: &uid}
	}
	externo := escopo.Tipo != db.EscopoInterno

	conj := Conjunto{}
	adicionar := func(p string) {
		// Defesa a mais: o banco já impede, por CHECK e gatilho, que usuário
		// externo tenha permissão vedada; se um dia falhar, ela não vale aqui.
		if externo && VedadaExterno(p) {
			return
		}
		conj[p] = struct{}{}
	}
	var codigos []string
	for _, p := range perfis {
		if !slices.Contains(codigos, p.CodigoOrigem) {
			codigos = append(codigos, p.CodigoOrigem)
		}
		for _, perm := range p.Permissoes {
			adicionar(perm)
		}
	}
	for _, perm := range avulsas {
		adicionar(perm)
	}
	slices.Sort(codigos)

	return Ator{
		EmpresaID:      empresaID,
		UsuarioID:      usuarioID,
		IdentidadeID:   *u.IdentidadeID,
		Escopo:         escopo,
		PessoaID:       u.PessoaID,
		Perfis:         codigos,
		Permissoes:     conj,
		Plano:          emp.Plano,
		SrvHabilitado:  emp.SrvHabilitado,
		MFAObrigatorio: emp.MfaObrigatorio,
		Idioma:         u.Idioma,
		IdiomaEmpresa:  emp.IdiomaPadrao,
		Pais:           emp.Pais,
		Moeda:          emp.Moeda,
	}, nil
}

// EmpresaExigeMFA lê empresa.mfa_obrigatorio na transação do contexto.
func EmpresaExigeMFA(ctx context.Context, empresaID uuid.UUID) (bool, error) {
	q, err := db.Q(ctx)
	if err != nil {
		return false, err
	}
	emp, err := q.AutzBuscarEmpresa(ctx, empresaID)
	if err != nil {
		return false, fmt.Errorf("ler empresa: %w", err)
	}
	return emp.MfaObrigatorio, nil
}

// NomeEmpresa lê o nome de exibição da empresa da transação do contexto.
func NomeEmpresa(ctx context.Context, empresaID uuid.UUID) (string, error) {
	q, err := db.Q(ctx)
	if err != nil {
		return "", err
	}
	emp, err := q.AutzBuscarEmpresa(ctx, empresaID)
	if err != nil {
		return "", fmt.Errorf("ler empresa: %w", err)
	}
	if emp.NomeFantasia != nil && *emp.NomeFantasia != "" {
		return *emp.NomeFantasia, nil
	}
	return emp.RazaoSocial, nil
}

// NomeUsuario lê o nome do usuário na transação do contexto.
func NomeUsuario(ctx context.Context, empresaID, usuarioID uuid.UUID) (string, error) {
	q, err := db.Q(ctx)
	if err != nil {
		return "", err
	}
	u, err := q.AutzBuscarUsuario(ctx, consultas.AutzBuscarUsuarioParams{EmpresaID: empresaID, ID: usuarioID})
	if err != nil {
		return "", fmt.Errorf("ler usuário: %w", err)
	}
	return u.Nome, nil
}
