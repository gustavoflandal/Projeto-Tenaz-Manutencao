package autenticacao

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"tenaz/internal/plataforma/autenticacao/internal/repo/gerado"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/email"
	"tenaz/internal/plataforma/erros"
	"tenaz/internal/plataforma/obs"
	"tenaz/internal/plataforma/textos"
)

// Finalidades de token_email.
const (
	finalidadeConvite     = "convite"
	finalidadeRedefinicao = "redefinicao_senha"
)

// O link leva o segredo no fragmento (#token=...): o fragmento não vai ao
// servidor na requisição da página, não aparece em log de acesso nem no
// cabeçalho Referer. A página lê o fragmento e o envia no corpo do POST.
func (s *Servico) link(caminho, segredo string) string {
	return s.cfg.URLWeb + caminho + "#token=" + segredo
}

// ---------------------------------------------------------------- senha

// RecuperarSenha envia o link de redefinição (1 hora) se o e-mail existir.
// Responde igual, e no mesmo tempo, exista ou não: o envio vai para segundo
// plano. Um novo pedido revoga o link anterior.
func (s *Servico) RecuperarSenha(ctx context.Context, emailInformado string) error {
	emailNorm := normalizarEmail(emailInformado)
	if emailNorm == "" {
		return nil
	}
	agora := s.agora()
	var segredo, destino string
	err := s.repo.Tx(ctx, func(q *gerado.Queries) error {
		id, err := q.AuthBuscarIdentidadePorEmail(ctx, emailNorm)
		if errors.Is(err, pgx.ErrNoRows) {
			return nil
		}
		if err != nil {
			return err
		}
		if id.BloqueadaEm != nil || id.SenhaHash == nil {
			return nil // bloqueio administrativo e conta só de SSO não redefinem senha
		}
		if err := q.AuthRevogarTokensEmailIdentidade(ctx, gerado.AuthRevogarTokensEmailIdentidadeParams{
			IdentidadeID: &id.ID, Finalidade: finalidadeRedefinicao,
		}); err != nil {
			return err
		}
		texto, hash, err := novoSegredo()
		if err != nil {
			return err
		}
		segredo, destino = texto, id.Email
		return q.AuthCriarTokenEmail(ctx, gerado.AuthCriarTokenEmailParams{
			ID: uuid.Must(uuid.NewV7()), TokenHash: hash, Finalidade: finalidadeRedefinicao, Email: id.Email,
			IdentidadeID: &id.ID, CriadoEm: agora, ExpiraEm: agora.Add(s.cfg.ValidadeRedefinicao),
		})
	})
	if err != nil {
		// Falha de infraestrutura também não muda a resposta: fica no log.
		obs.Logger(ctx).Error("preparar redefinição de senha", "erro", err)
		return nil
	}
	if segredo == "" {
		return nil
	}
	// E-mail disparado pela própria pessoa: o idioma negociado na requisição,
	// porque a identidade é global e não guarda idioma (ADR-013, Parte B, item 7).
	idioma := textos.DoContexto(ctx)
	msg := email.Mensagem{
		Para:    destino,
		Assunto: textos.T(idioma, "email.redefinicao_senha.assunto", nil),
		Texto:   textos.T(idioma, "email.redefinicao_senha.corpo", textos.Parametros{"link": s.link("/redefinir-senha", segredo)}),
	}
	s.enviarEmSegundoPlano(ctx, msg)
	return nil
}

// enviarEmSegundoPlano envia sem prender a resposta. O contexto do envio não
// é cancelado com a requisição, mas tem prazo próprio. O corpo, que leva o
// link, nunca vai para o log.
func (s *Servico) enviarEmSegundoPlano(ctx context.Context, m email.Mensagem) {
	log := obs.Logger(ctx)
	base := context.WithoutCancel(ctx)
	s.envios.Add(1)
	go func() {
		defer s.envios.Done()
		ctxEnvio, cancelar := context.WithTimeout(base, time.Minute)
		defer cancelar()
		if err := s.d.Email.Enviar(ctxEnvio, m); err != nil {
			log.Error("enviar e-mail", "assunto", m.Assunto, "erro", err)
		}
	}()
}

// RedefinirSenha define a senha nova com o link do e-mail. A política é
// conferida antes de gastar o link, para que uma senha fraca não obrigue a
// pedir outro. Redefinir revoga todas as sessões da identidade, em todas as
// empresas, e confirma o e-mail (quem abriu o link o controla).
func (s *Servico) RedefinirSenha(ctx context.Context, token, novaSenha string) error {
	if !segredoBemFormado(token) {
		return erros.LinkEmailInvalido
	}
	hash := hashSegredo(token)
	var alvo gerado.AuthBuscarTokenEmailValidoRow
	err := s.repo.Tx(ctx, func(q *gerado.Queries) error {
		var err error
		alvo, err = q.AuthBuscarTokenEmailValido(ctx, gerado.AuthBuscarTokenEmailValidoParams{TokenHash: hash, Finalidade: finalidadeRedefinicao})
		return naoEncontrado(err, erros.LinkEmailInvalido)
	})
	if err != nil {
		return err
	}
	if err := s.d.Politica.Validar(novaSenha, alvo.Email); err != nil {
		return err
	}
	senhaHash, err := s.d.Hasher.Gerar(ctx, novaSenha)
	if err != nil {
		return err
	}
	var identidadeID uuid.UUID
	err = s.repo.Tx(ctx, func(q *gerado.Queries) error {
		usado, err := q.AuthUsarTokenEmail(ctx, gerado.AuthUsarTokenEmailParams{TokenHash: hash, Finalidade: finalidadeRedefinicao})
		if err != nil {
			return naoEncontrado(err, erros.LinkEmailInvalido)
		}
		if usado.IdentidadeID == nil {
			return erros.LinkEmailInvalido
		}
		identidadeID = *usado.IdentidadeID
		if err := q.AuthDefinirSenha(ctx, gerado.AuthDefinirSenhaParams{ID: identidadeID, SenhaHash: &senhaHash}); err != nil {
			return err
		}
		if err := q.AuthMarcarEmailVerificado(ctx, identidadeID); err != nil {
			return err
		}
		return revogarSessoesIdentidade(ctx, q, identidadeID, nil)
	})
	if err != nil {
		return err
	}
	return s.eventoNasEmpresas(ctx, nil, eventoIdentidade(identidadeID, "senha_hash"))
}

// TrocarSenha troca a senha de quem está logado, conferindo a atual. A senha
// atual errada conta como falha de login (e leva ao bloqueio). As outras
// sessões da identidade são revogadas; a atual continua.
func (s *Servico) TrocarSenha(ctx context.Context, a autorizacao.Ator, atual, nova string) error {
	var id gerado.AuthBuscarIdentidadeRow
	err := s.repo.Tx(ctx, func(q *gerado.Queries) error {
		var err error
		id, err = q.AuthBuscarIdentidade(ctx, a.IdentidadeID)
		return naoEncontrado(err, erros.CredencialInvalida)
	})
	if err != nil {
		return err
	}
	if id.SenhaHash == nil {
		return erros.TransicaoInvalida.ComChave("detalhe.identidade_sso", nil).ComDado("estado_atual", "sso")
	}
	ok, _, err := s.d.Hasher.Verificar(ctx, atual, *id.SenhaHash)
	if err != nil {
		return err
	}
	if !ok {
		if err := s.registrarFalha(ctx, id.ID, s.agora()); err != nil {
			return err
		}
		return erros.CredencialInvalida.ComCampoChave("senha_atual", erros.CampoValorInvalido, "campo.senha_atual", nil)
	}
	if err := s.d.Politica.Validar(nova, id.Email); err != nil {
		return err
	}
	senhaHash, err := s.d.Hasher.Gerar(ctx, nova)
	if err != nil {
		return err
	}
	err = s.repo.Tx(ctx, func(q *gerado.Queries) error {
		if err := q.AuthDefinirSenha(ctx, gerado.AuthDefinirSenhaParams{ID: id.ID, SenhaHash: &senhaHash}); err != nil {
			return err
		}
		return revogarSessoesIdentidade(ctx, q, id.ID, &a.SessaoID)
	})
	if err != nil {
		return err
	}
	return s.eventoNasEmpresas(ctx, nil, eventoIdentidade(id.ID, "senha_hash"))
}

// --------------------------------------------------------------- convite

// EmitirConvite grava o link de convite (7 dias) do usuário `convidado`,
// revogando o anterior ("reenviar convite"), e envia o e-mail. Passo (2) do
// convite na seção 6.9 do documento 5: chame depois do COMMIT que criou o
// usuário; a chave estrangeira composta de token_email garante que o usuário
// é da empresa informada, que vem do ator autenticado (ou da conta recém
// criada).
//
// Falha no envio: servico_indisponivel. O usuário continua convidado, com o
// link gravado; o administrador reenvia.
func (s *Servico) EmitirConvite(ctx context.Context, empresaID, usuarioID uuid.UUID, emailDestino, nomeEmpresa, nomeConvidado, idiomaDestinatario string) error {
	agora := s.agora()
	segredo, hash, err := novoSegredo()
	if err != nil {
		return err
	}
	destino := normalizarEmail(emailDestino)
	err = s.repo.Tx(ctx, func(q *gerado.Queries) error {
		if err := q.AuthRevogarTokensEmailUsuario(ctx, gerado.AuthRevogarTokensEmailUsuarioParams{
			EmpresaID: &empresaID, UsuarioID: &usuarioID, Finalidade: finalidadeConvite,
		}); err != nil {
			return err
		}
		return q.AuthCriarTokenEmail(ctx, gerado.AuthCriarTokenEmailParams{
			ID: uuid.Must(uuid.NewV7()), TokenHash: hash, Finalidade: finalidadeConvite, Email: destino,
			EmpresaID: &empresaID, UsuarioID: &usuarioID, CriadoEm: agora, ExpiraEm: agora.Add(s.cfg.ValidadeConvite),
		})
	})
	if err != nil {
		return fmt.Errorf("emitir convite: %w", err)
	}
	// E-mail disparado por outra pessoa: o idioma do destinatário, que é o
	// usuario.idioma (no convite, o padrão da empresa).
	idioma := textos.PtBR
	if textos.Suportado(idiomaDestinatario) {
		idioma = textos.Idioma(idiomaDestinatario)
	}
	err = s.d.Email.Enviar(ctx, email.Mensagem{
		Para:    destino,
		Assunto: textos.T(idioma, "email.convite.assunto", textos.Parametros{"empresa": nomeEmpresa}),
		Texto: textos.T(idioma, "email.convite.corpo", textos.Parametros{
			"nome": nomeConvidado, "empresa": nomeEmpresa, "dias": int(s.cfg.ValidadeConvite.Hours() / 24),
			"link": s.link("/convite", segredo),
		}),
	})
	if err != nil {
		obs.Logger(ctx).Error("enviar convite", "erro", err)
		return erros.ServicoIndisponivel.ComChave("detalhe.convite_email_nao_saiu", nil).Por(err)
	}
	return nil
}

// ConviteValido é o convite localizado pelo link.
type ConviteValido struct {
	EmpresaID        uuid.UUID
	UsuarioID        uuid.UUID
	Email            string
	IdentidadeExiste bool // o e-mail já tem login: a senha não é pedida
}

// LocalizarConvite confere o link sem gastá-lo, para a tela de aceite saber
// se pede senha (identidade nova) ou não (o e-mail já entra em outra empresa).
func (s *Servico) LocalizarConvite(ctx context.Context, token string) (ConviteValido, error) {
	if !segredoBemFormado(token) {
		return ConviteValido{}, erros.LinkEmailInvalido
	}
	var c ConviteValido
	err := s.repo.Tx(ctx, func(q *gerado.Queries) error {
		t, err := q.AuthBuscarTokenEmailValido(ctx, gerado.AuthBuscarTokenEmailValidoParams{
			TokenHash: hashSegredo(token), Finalidade: finalidadeConvite,
		})
		if err != nil {
			return naoEncontrado(err, erros.LinkEmailInvalido)
		}
		if t.EmpresaID == nil || t.UsuarioID == nil {
			return erros.LinkEmailInvalido
		}
		c = ConviteValido{EmpresaID: *t.EmpresaID, UsuarioID: *t.UsuarioID, Email: t.Email}
		_, err = q.AuthBuscarIdentidadePorEmail(ctx, t.Email)
		switch {
		case err == nil:
			c.IdentidadeExiste = true
		case !errors.Is(err, pgx.ErrNoRows):
			return err
		}
		return nil
	})
	return c, err
}

// AceitarConvite conclui o convite. A identidade do e-mail é usada se já
// existir (a senha informada é ignorada: a pessoa entra com a que já tem) ou
// criada com a senha, conferida pela política. ligar é o caso de uso do
// módulo ADM que, na transação da empresa, liga a identidade ao usuário e o
// ativa; o link só é gasto depois que ligar confirma.
func (s *Servico) AceitarConvite(ctx context.Context, token, senha string,
	ligar func(ctx context.Context, c ConviteValido, identidadeID uuid.UUID) error) (ConviteValido, error) {
	c, err := s.LocalizarConvite(ctx, token)
	if err != nil {
		return ConviteValido{}, err
	}
	var senhaHash string
	if !c.IdentidadeExiste {
		if err := s.d.Politica.Validar(senha, c.Email); err != nil {
			return ConviteValido{}, err
		}
		if senhaHash, err = s.d.Hasher.Gerar(ctx, senha); err != nil {
			return ConviteValido{}, err
		}
	}
	var identidadeID uuid.UUID
	err = s.repo.Tx(ctx, func(q *gerado.Queries) error {
		if senhaHash != "" {
			novo, err := q.AuthCriarIdentidade(ctx, gerado.AuthCriarIdentidadeParams{
				ID: uuid.Must(uuid.NewV7()), Email: c.Email, SenhaHash: &senhaHash,
			})
			if err == nil {
				identidadeID = novo
				return nil
			}
			if !errors.Is(err, pgx.ErrNoRows) {
				return err
			}
			// Criada por outro convite no meio-tempo: usa a existente.
		}
		id, err := q.AuthBuscarIdentidadePorEmail(ctx, c.Email)
		if err != nil {
			return err
		}
		identidadeID = id.ID
		// O link chegou ao e-mail: ele está confirmado.
		return q.AuthMarcarEmailVerificado(ctx, id.ID)
	})
	if err != nil {
		return ConviteValido{}, err
	}
	if err := ligar(ctx, c, identidadeID); err != nil {
		return ConviteValido{}, err
	}
	// (3) consome o link e grava o evento, na trilha da empresa do convite.
	// O acesso_empresa já existe: o trigger o criou no passo (2).
	err = s.repo.NaTrilha(ctx, c.EmpresaID, func(q *gerado.Queries) error {
		_, err := q.AuthUsarTokenEmail(ctx, gerado.AuthUsarTokenEmailParams{TokenHash: hashSegredo(token), Finalidade: finalidadeConvite})
		if errors.Is(err, pgx.ErrNoRows) {
			return nil // gasto por um aceite simultâneo; o usuário já está ligado
		}
		if err != nil {
			return err
		}
		return gravarEvento(ctx, q, c.EmpresaID, c.UsuarioID, eventoAuth{
			Acao: acaoAlteracao, IdentidadeID: identidadeID, Campos: []string{"email_verificado_em"},
			Depois: map[string]string{"motivo": "convite_aceito"},
		})
	})
	return c, err
}
