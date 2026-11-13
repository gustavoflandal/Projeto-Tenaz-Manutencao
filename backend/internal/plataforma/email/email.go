// Package email envia os e-mails transacionais da Tenaz: convite, redefinição
// de senha e verificação de e-mail (token_email, RF-ADM-004, RF-ADM-006).
//
// Em desenvolvimento o destino é o Mailpit do compose local (SMTP em
// localhost:1025, caixa em http://localhost:8025), que não entrega nada a
// ninguém. Em produção, o provedor real por SMTP com STARTTLS obrigatório.
//
// O corpo leva o link com o segredo de uso único. Por isso nenhum erro deste
// pacote inclui o corpo, e quem chama nunca registra a mensagem em log.
package email

import (
	"bytes"
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"mime"
	"mime/quotedprintable"
	"net"
	"net/mail"
	"net/smtp"
	"net/textproto"
	"strings"
	"time"

	"tenaz/internal/plataforma/resiliencia"
)

// Mensagem é um e-mail de texto simples.
type Mensagem struct {
	Para    string
	Assunto string
	Texto   string
}

// Enviador é a porta de envio. Os testes usam uma implementação em memória.
type Enviador interface {
	Enviar(ctx context.Context, m Mensagem) error
}

// ConfigSMTP configura o envio.
type ConfigSMTP struct {
	Endereco      string // host:porta
	Usuario       string // vazio: sem autenticação (Mailpit)
	Senha         string
	ExigirTLS     bool // STARTTLS obrigatório; sempre verdadeiro fora do desenvolvimento
	Remetente     string
	NomeRemetente string
}

// SMTP envia por SMTP com tempo limite, repetição e disjuntor.
type SMTP struct {
	cfg       ConfigSMTP
	remetente *mail.Address
	disjuntor *resiliencia.Disjuntor
}

// NovoSMTP valida a configuração e monta o enviador.
func NovoSMTP(cfg ConfigSMTP) (*SMTP, error) {
	if _, _, err := net.SplitHostPort(cfg.Endereco); err != nil {
		return nil, fmt.Errorf("endereço SMTP inválido %q: %w", cfg.Endereco, err)
	}
	rem, err := mail.ParseAddress(cfg.Remetente)
	if err != nil {
		return nil, fmt.Errorf("remetente inválido: %w", err)
	}
	rem.Name = cfg.NomeRemetente
	if cfg.Usuario != "" && !cfg.ExigirTLS {
		return nil, errors.New("autenticação SMTP sem TLS não é aceita: a senha iria em claro")
	}
	return &SMTP{
		cfg:       cfg,
		remetente: rem,
		disjuntor: resiliencia.NovoDisjuntor(resiliencia.Politica{
			TempoLimite:     10 * time.Second,
			Tentativas:      3,
			RecuoInicial:    500 * time.Millisecond,
			RecuoMaximo:     4 * time.Second,
			FalhasParaAbrir: 5,
			TempoAberto:     30 * time.Second,
		}),
	}, nil
}

// Enviar monta e envia a mensagem.
func (s *SMTP) Enviar(ctx context.Context, m Mensagem) error {
	para, err := mail.ParseAddress(m.Para)
	if err != nil {
		return resiliencia.Permanente(fmt.Errorf("destinatário inválido: %w", err))
	}
	corpo, err := s.montar(para, m)
	if err != nil {
		return resiliencia.Permanente(err)
	}
	err = s.disjuntor.Executar(ctx, func(ctx context.Context) error {
		return s.entregar(ctx, para.Address, corpo)
	})
	if err != nil {
		return fmt.Errorf("enviar e-mail: %w", err)
	}
	return nil
}

func (s *SMTP) montar(para *mail.Address, m Mensagem) ([]byte, error) {
	// Quebra de linha no assunto seria injeção de cabeçalho.
	if strings.ContainsAny(m.Assunto, "\r\n") {
		return nil, errors.New("assunto com quebra de linha")
	}
	id := make([]byte, 16)
	if _, err := rand.Read(id); err != nil {
		return nil, err
	}
	dominio := s.remetente.Address[strings.LastIndex(s.remetente.Address, "@")+1:]

	var b bytes.Buffer
	cab := textproto.MIMEHeader{}
	cab.Set("From", s.remetente.String())
	cab.Set("To", para.String())
	cab.Set("Subject", mime.QEncoding.Encode("utf-8", m.Assunto))
	cab.Set("Date", time.Now().UTC().Format(time.RFC1123Z))
	cab.Set("Message-ID", "<"+hex.EncodeToString(id)+"@"+dominio+">")
	cab.Set("MIME-Version", "1.0")
	cab.Set("Content-Type", "text/plain; charset=UTF-8")
	cab.Set("Content-Transfer-Encoding", "quoted-printable")
	for _, k := range []string{"From", "To", "Subject", "Date", "Message-ID", "MIME-Version", "Content-Type",
		"Content-Transfer-Encoding"} {
		fmt.Fprintf(&b, "%s: %s\r\n", k, cab.Get(k))
	}
	b.WriteString("\r\n")
	qp := quotedprintable.NewWriter(&b)
	if _, err := qp.Write([]byte(strings.ReplaceAll(m.Texto, "\n", "\r\n"))); err != nil {
		return nil, err
	}
	if err := qp.Close(); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (s *SMTP) entregar(ctx context.Context, para string, corpo []byte) error {
	host, _, _ := net.SplitHostPort(s.cfg.Endereco)
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", s.cfg.Endereco)
	if err != nil {
		return err
	}
	if prazo, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(prazo)
	}
	c, err := smtp.NewClient(conn, host)
	if err != nil {
		_ = conn.Close()
		return err
	}
	defer c.Close()

	if ok, _ := c.Extension("STARTTLS"); ok {
		if err := c.StartTLS(&tls.Config{ServerName: host, MinVersion: tls.VersionTLS12}); err != nil {
			return err
		}
	} else if s.cfg.ExigirTLS {
		return resiliencia.Permanente(errors.New("o servidor SMTP não oferece STARTTLS"))
	}
	if s.cfg.Usuario != "" {
		if err := c.Auth(smtp.PlainAuth("", s.cfg.Usuario, s.cfg.Senha, host)); err != nil {
			return resiliencia.Permanente(fmt.Errorf("autenticar no SMTP: %w", err))
		}
	}
	if err := c.Mail(s.remetente.Address); err != nil {
		return classificar(err)
	}
	if err := c.Rcpt(para); err != nil {
		return classificar(err)
	}
	w, err := c.Data()
	if err != nil {
		return classificar(err)
	}
	if _, err := w.Write(corpo); err != nil {
		return err
	}
	if err := w.Close(); err != nil {
		return classificar(err)
	}
	return c.Quit()
}

// classificar: resposta 5xx do servidor é definitiva; 4xx e erro de rede,
// transitórios.
func classificar(err error) error {
	var tp *textproto.Error
	if errors.As(err, &tp) && tp.Code >= 500 {
		return resiliencia.Permanente(err)
	}
	return err
}
