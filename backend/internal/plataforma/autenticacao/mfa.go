package autenticacao

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"tenaz/internal/plataforma/autenticacao/internal/repo/gerado"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/erros"
)

// ConfiguracaoMFA é o que a tela de ativação mostra: o segredo para digitar
// e a URI otpauth:// para o QR code. Só existe nesta resposta; o banco guarda
// o segredo cifrado.
type ConfiguracaoMFA struct {
	SegredoBase32 string
	URIOTPAuth    string
}

// IniciarMFA gera o segredo TOTP e o grava pendente. Só passa a valer depois
// de ConfirmarMFA com um código certo. Chamar de novo antes de confirmar
// troca o segredo pendente.
func (s *Servico) IniciarMFA(ctx context.Context, a autorizacao.Ator) (ConfiguracaoMFA, error) {
	bruto, base32, err := novoSegredoTOTP()
	if err != nil {
		return ConfiguracaoMFA{}, err
	}
	cifrado, err := s.d.Cifrador.Cifrar(bruto, a.IdentidadeID[:])
	if err != nil {
		return ConfiguracaoMFA{}, fmt.Errorf("cifrar segredo TOTP: %w", err)
	}
	var emailID string
	err = s.repo.Tx(ctx, func(q *gerado.Queries) error {
		id, err := q.AuthTravarIdentidade(ctx, a.IdentidadeID)
		if err != nil {
			return naoEncontrado(err, erros.CredencialInvalida)
		}
		if id.MfaAtivadoEm != nil {
			return erros.TransicaoInvalida.ComChave("detalhe.mfa_ja_ativo", nil).ComDado("estado_atual", "ativo")
		}
		emailID = id.Email
		return q.AuthGravarSegredoMFAPendente(ctx, gerado.AuthGravarSegredoMFAPendenteParams{
			ID: id.ID, MfaSegredoCifrado: cifrado,
		})
	})
	if err != nil {
		return ConfiguracaoMFA{}, err
	}
	return ConfiguracaoMFA{SegredoBase32: base32, URIOTPAuth: uriOTPAuth(s.cfg.EmissorTOTP, emailID, base32)}, nil
}

// ConfirmarMFA ativa o segundo fator com o primeiro código certo.
func (s *Servico) ConfirmarMFA(ctx context.Context, a autorizacao.Ator, codigo string) error {
	agora := s.agora()
	err := s.repo.Tx(ctx, func(q *gerado.Queries) error {
		id, err := q.AuthTravarIdentidade(ctx, a.IdentidadeID)
		if err != nil {
			return naoEncontrado(err, erros.CredencialInvalida)
		}
		if id.MfaAtivadoEm != nil {
			return erros.TransicaoInvalida.ComChave("detalhe.mfa_ja_ativo", nil).ComDado("estado_atual", "ativo")
		}
		if id.MfaSegredoCifrado == nil {
			return erros.TransicaoInvalida.ComChave("detalhe.mfa_iniciar_antes", nil).ComDado("estado_atual", "inativo")
		}
		segredo, err := s.d.Cifrador.Decifrar(id.MfaSegredoCifrado, id.ID[:])
		if err != nil {
			return fmt.Errorf("decifrar segredo TOTP: %w", err)
		}
		// O passo do código de confirmação fica gravado: ele não serve de novo
		// no primeiro login (mfa_ultimo_passo, zerado a cada nova ativação).
		certo, err := consumirTOTP(ctx, q, id.ID, segredo, codigo, agora)
		if err != nil {
			return err
		}
		if !certo {
			return erros.MFACodigoInvalido
		}
		_, err = q.AuthAtivarMFA(ctx, id.ID)
		return err
	})
	if err != nil {
		return err
	}
	return s.eventoNasEmpresas(ctx, nil, eventoIdentidade(a.IdentidadeID, "mfa_ativado_em"))
}

// DesativarMFA desliga o segundo fator, exigindo um código válido. Recusado
// com transicao_invalida (estado_atual "exigido_pela_empresa") enquanto a
// identidade tiver acesso habilitado a alguma empresa que exige MFA
// (RN-ADM-027, proposta). Código errado soma falha ao bloqueio progressivo,
// como no login: uma sessão roubada não testa códigos sem limite.
func (s *Servico) DesativarMFA(ctx context.Context, a autorizacao.Ator, codigo string) error {
	exigida, err := s.algumaEmpresaExigeMFA(ctx, a.IdentidadeID)
	if err != nil {
		return err
	}
	if exigida {
		return erros.TransicaoInvalida.ComChave("detalhe.mfa_empresa_exige", nil).
			ComDado("estado_atual", "exigido_pela_empresa")
	}
	agora := s.agora()
	errado := false
	err = s.repo.Tx(ctx, func(q *gerado.Queries) error {
		id, err := q.AuthTravarIdentidade(ctx, a.IdentidadeID)
		if err != nil {
			return naoEncontrado(err, erros.CredencialInvalida)
		}
		if id.MfaAtivadoEm == nil {
			return erros.TransicaoInvalida.ComChave("detalhe.mfa_inativo", nil).ComDado("estado_atual", "inativo")
		}
		segredo, err := s.d.Cifrador.Decifrar(id.MfaSegredoCifrado, id.ID[:])
		if err != nil {
			return fmt.Errorf("decifrar segredo TOTP: %w", err)
		}
		certo, err := consumirTOTP(ctx, q, id.ID, segredo, codigo, agora)
		if err != nil {
			return err
		}
		if !certo {
			errado = true
			return nil
		}
		return q.AuthDesativarMFA(ctx, id.ID)
	})
	if err != nil {
		return err
	}
	if errado {
		if err := s.registrarFalha(ctx, a.IdentidadeID, agora); err != nil {
			return err
		}
		return erros.MFACodigoInvalido
	}
	return s.eventoNasEmpresas(ctx, nil, eventoIdentidade(a.IdentidadeID, "mfa_ativado_em"))
}

// algumaEmpresaExigeMFA lê empresa.mfa_obrigatorio de cada empresa em que a
// identidade tem acesso habilitado, no pool da aplicação.
func (s *Servico) algumaEmpresaExigeMFA(ctx context.Context, identidadeID uuid.UUID) (bool, error) {
	lista, err := s.acessos(ctx, identidadeID)
	if err != nil {
		return false, err
	}
	for _, ac := range lista {
		var exige bool
		err := s.d.Banco.NaEmpresa(ctx, ac.EmpresaID, db.Interno(nil), func(ctx context.Context) error {
			var err error
			exige, err = autorizacao.EmpresaExigeMFA(ctx, ac.EmpresaID)
			return err
		})
		if err != nil {
			return false, err
		}
		if exige {
			return true, nil
		}
	}
	return false, nil
}

// eventoIdentidade é a alteração de credencial na trilha: só o nome do campo,
// sem antes nem depois (seção 9.1 do documento 2).
func eventoIdentidade(identidadeID uuid.UUID, campo string) eventoAuth {
	return eventoAuth{Acao: acaoAlteracao, IdentidadeID: identidadeID, Campos: []string{campo}}
}
