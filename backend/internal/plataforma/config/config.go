// Package config lê a configuração do processo a partir do ambiente.
//
// Nenhum segredo tem valor padrão: ausência de segredo é erro na partida, não
// um valor conveniente que vaza para produção.
package config

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Config é a configuração de um processo da Tenaz.
type Config struct {
	Ambiente        string
	EnderecoHTTP    string
	BancoURL        string
	BancoMaxConn    int32
	TempoLimiteHTTP time.Duration
	NivelLog        string
	LogJSON         bool

	// BancoPapel: papel assumido por SET ROLE no pool da aplicação, para o
	// desenvolvimento local, em que TENAZ_BANCO_URL é a do migrador
	// (superusuário, que ignora o RLS). Em homologação e produção fica vazio:
	// a URL já é a de tenaz_app.
	BancoPapel string

	// Pool próprio das credenciais de autenticação (ADR-011; seção 6.9 do
	// documento 5), aberto só no modo api.
	BancoURLAutenticacao     string
	BancoAutenticacaoMaxConn int32

	Autenticacao Autenticacao

	// BancoURLJob: conexão do papel tenaz_job, dos subcomandos de job
	// (TENAZ_BANCO_URL_JOB). Em desenvolvimento, vazia, usa TENAZ_BANCO_URL
	// com SET ROLE tenaz_job.
	BancoURLJob string

	// Armazenamento é o S3 do selo da trilha de auditoria (ADR-012).
	Armazenamento Armazenamento

	// CookieSeguro: cookie de sessão __Host-tenaz_sessao com Secure
	// (TENAZ_COOKIE_SEGURO, padrão verdadeiro). Falso só em desenvolvimento
	// sobre http, em que o navegador recusa o prefixo __Host- sem Secure; aí
	// o cookie se chama tenaz_sessao (seção 4.2 do documento 3).
	CookieSeguro bool

	// OrigensWeb: origens aceitas em Origin nas requisições do navegador que
	// alteram dados (TENAZ_ORIGENS_WEB, separadas por vírgula). Vazio: a
	// origem de TENAZ_URL_WEB.
	OrigensWeb []string
}

// Autenticacao são os segredos e endereços do serviço de autenticação. Nenhum
// tem valor padrão: o modo api com banco recusa subir sem eles.
type Autenticacao struct {
	URLWeb              string // base dos links de e-mail (TENAZ_URL_WEB)
	ChaveMFA            string // TENAZ_CHAVE_MFA: 32 bytes em base64 (AES-256-GCM do segredo TOTP)
	ChaveJWT            string // TENAZ_CHAVE_JWT: semente Ed25519 de 32 bytes em base64
	ChaveJWTID          string // TENAZ_CHAVE_JWT_ID: kid
	ChavesJWTAnteriores string // TENAZ_CHAVES_JWT_ANTERIORES: kid=chave_publica_base64,... (rotação)
	SMTPEndereco        string // TENAZ_SMTP_ENDERECO, host:porta
	SMTPUsuario         string
	SMTPSenha           string
	EmailRemetente      string // TENAZ_EMAIL_REMETENTE
}

// Armazenamento é a configuração do armazenamento de objetos compatível com
// S3 (ADR-012, decisão 5): endpoint, região e estilo de caminho por
// configuração, para que trocar de servidor seja trocar variáveis. Em
// desenvolvimento é o SeaweedFS do compose, com a identidade tenaz_auditoria;
// em produção, o S3 da AWS. Sem valor padrão de segredo.
type Armazenamento struct {
	Endpoint         string // TENAZ_S3_ENDPOINT; vazio: o endpoint da AWS para a região
	Regiao           string // TENAZ_S3_REGIAO
	CaminhoExplicito bool   // TENAZ_S3_CAMINHO_EXPLICITO (path-style): o SeaweedFS exige
	BucketAuditoria  string // TENAZ_S3_BUCKET_AUDITORIA
	ChaveAcesso      string // TENAZ_S3_CHAVE_ACESSO
	ChaveSecreta     string // TENAZ_S3_CHAVE_SECRETA
}

// Configurado informa se há armazenamento configurado.
func (a Armazenamento) Configurado() bool {
	return a.Regiao != "" && a.BucketAuditoria != "" && a.ChaveAcesso != "" && a.ChaveSecreta != ""
}

// Ambientes aceitos. Produção tem regras próprias (ver Validar).
const (
	AmbienteDesenvolvimento = "desenvolvimento"
	AmbienteHomologacao     = "homologacao"
	AmbienteProducao        = "producao"
)

// DeAmbiente monta a configuração a partir das variáveis de ambiente.
func DeAmbiente() (Config, error) {
	c := Config{
		Ambiente:        texto("TENAZ_AMBIENTE", AmbienteDesenvolvimento),
		EnderecoHTTP:    texto("TENAZ_ENDERECO_HTTP", ":8080"),
		BancoURL:        os.Getenv("TENAZ_BANCO_URL"),
		NivelLog:        texto("TENAZ_NIVEL_LOG", "info"),
		TempoLimiteHTTP: 30 * time.Second,
	}

	maxConn, err := inteiro("TENAZ_BANCO_MAX_CONN", 10)
	if err != nil {
		return Config{}, err
	}
	c.BancoMaxConn = int32(maxConn)

	segundos, err := inteiro("TENAZ_TEMPO_LIMITE_HTTP_S", 30)
	if err != nil {
		return Config{}, err
	}
	c.TempoLimiteHTTP = time.Duration(segundos) * time.Second

	c.LogJSON = c.Ambiente != AmbienteDesenvolvimento

	c.BancoPapel = os.Getenv("TENAZ_BANCO_PAPEL")
	c.BancoURLAutenticacao = os.Getenv("TENAZ_BANCO_URL_AUTENTICACAO")
	maxAuth, err := inteiro("TENAZ_BANCO_AUTENTICACAO_MAX_CONN", 5)
	if err != nil {
		return Config{}, err
	}
	c.BancoAutenticacaoMaxConn = int32(maxAuth)

	c.Autenticacao = Autenticacao{
		URLWeb:              os.Getenv("TENAZ_URL_WEB"),
		ChaveMFA:            os.Getenv("TENAZ_CHAVE_MFA"),
		ChaveJWT:            os.Getenv("TENAZ_CHAVE_JWT"),
		ChaveJWTID:          os.Getenv("TENAZ_CHAVE_JWT_ID"),
		ChavesJWTAnteriores: os.Getenv("TENAZ_CHAVES_JWT_ANTERIORES"),
		SMTPEndereco:        os.Getenv("TENAZ_SMTP_ENDERECO"),
		SMTPUsuario:         os.Getenv("TENAZ_SMTP_USUARIO"),
		SMTPSenha:           os.Getenv("TENAZ_SMTP_SENHA"),
		EmailRemetente:      os.Getenv("TENAZ_EMAIL_REMETENTE"),
	}

	c.BancoURLJob = os.Getenv("TENAZ_BANCO_URL_JOB")
	c.Armazenamento = Armazenamento{
		Endpoint:        os.Getenv("TENAZ_S3_ENDPOINT"),
		Regiao:          os.Getenv("TENAZ_S3_REGIAO"),
		BucketAuditoria: os.Getenv("TENAZ_S3_BUCKET_AUDITORIA"),
		ChaveAcesso:     os.Getenv("TENAZ_S3_CHAVE_ACESSO"),
		ChaveSecreta:    os.Getenv("TENAZ_S3_CHAVE_SECRETA"),
	}
	switch v := strings.ToLower(strings.TrimSpace(os.Getenv("TENAZ_S3_CAMINHO_EXPLICITO"))); v {
	case "", "0", "false", "nao", "não":
	case "1", "true", "sim":
		c.Armazenamento.CaminhoExplicito = true
	default:
		return Config{}, fmt.Errorf("TENAZ_S3_CAMINHO_EXPLICITO deve ser true ou false, recebeu %q", v)
	}

	c.CookieSeguro = true
	switch v := strings.ToLower(strings.TrimSpace(os.Getenv("TENAZ_COOKIE_SEGURO"))); v {
	case "", "1", "true", "sim":
	case "0", "false", "nao", "não":
		c.CookieSeguro = false
	default:
		return Config{}, fmt.Errorf("TENAZ_COOKIE_SEGURO deve ser true ou false, recebeu %q", v)
	}
	for _, o := range strings.Split(os.Getenv("TENAZ_ORIGENS_WEB"), ",") {
		if o = strings.TrimRight(strings.TrimSpace(o), "/"); o != "" {
			c.OrigensWeb = append(c.OrigensWeb, o)
		}
	}

	return c, c.Validar()
}

// Validar recusa configuração incoerente antes de o processo subir.
func (c Config) Validar() error {
	switch c.Ambiente {
	case AmbienteDesenvolvimento, AmbienteHomologacao, AmbienteProducao:
	default:
		return fmt.Errorf("TENAZ_AMBIENTE inválido: %q (use desenvolvimento, homologacao ou producao)", c.Ambiente)
	}
	if c.BancoMaxConn < 1 {
		return fmt.Errorf("TENAZ_BANCO_MAX_CONN deve ser pelo menos 1, recebeu %d", c.BancoMaxConn)
	}
	if c.Ambiente == AmbienteProducao && strings.Contains(c.BancoURL, "sslmode=disable") {
		return fmt.Errorf("conexão sem TLS não é aceita em produção (RNF-032)")
	}
	if c.BancoAutenticacaoMaxConn < 1 {
		return fmt.Errorf("TENAZ_BANCO_AUTENTICACAO_MAX_CONN deve ser pelo menos 1, recebeu %d", c.BancoAutenticacaoMaxConn)
	}
	if !c.CookieSeguro && c.Ambiente != AmbienteDesenvolvimento {
		return fmt.Errorf("TENAZ_COOKIE_SEGURO=false só é aceito em desenvolvimento; fora dele o cookie é sempre __Host- com Secure")
	}
	if c.BancoPapel != "" && c.Ambiente != AmbienteDesenvolvimento {
		return fmt.Errorf("TENAZ_BANCO_PAPEL só existe em desenvolvimento; fora dele, a URL já é a do papel")
	}
	return nil
}

// BancoAutenticacao devolve a URL e o papel do pool da autenticação no modo
// api (ADR-011). Em homologação e produção a URL é obrigatória, diferente de
// TENAZ_BANCO_URL e, em produção, com TLS. Em desenvolvimento, vazia, usa
// TENAZ_BANCO_URL com SET ROLE tenaz_autenticacao.
func (c Config) BancoAutenticacao() (url, papel string, err error) {
	if c.Ambiente == AmbienteDesenvolvimento {
		if c.BancoURLAutenticacao == "" {
			return c.BancoURL, "tenaz_autenticacao", nil
		}
		return c.BancoURLAutenticacao, "", nil
	}
	switch {
	case c.BancoURLAutenticacao == "":
		return "", "", fmt.Errorf("TENAZ_BANCO_URL_AUTENTICACAO não definida (obrigatória no modo api, ADR-011)")
	case c.BancoURLAutenticacao == c.BancoURL:
		return "", "", fmt.Errorf("TENAZ_BANCO_URL_AUTENTICACAO não pode ser igual a TENAZ_BANCO_URL (ADR-011)")
	case c.Ambiente == AmbienteProducao && strings.Contains(c.BancoURLAutenticacao, "sslmode=disable"):
		return "", "", fmt.Errorf("conexão sem TLS não é aceita em produção (RNF-032)")
	}
	return c.BancoURLAutenticacao, "", nil
}

// BancoJob devolve a URL e o papel do pool dos jobs (tenaz_job). Em
// homologação e produção a URL é obrigatória e diferente das demais; em
// desenvolvimento, vazia, usa TENAZ_BANCO_URL com SET ROLE tenaz_job.
func (c Config) BancoJob() (url, papel string, err error) {
	if c.Ambiente == AmbienteDesenvolvimento {
		if c.BancoURLJob == "" {
			return c.BancoURL, "tenaz_job", nil
		}
		return c.BancoURLJob, "", nil
	}
	switch {
	case c.BancoURLJob == "":
		return "", "", fmt.Errorf("TENAZ_BANCO_URL_JOB não definida (obrigatória nos jobs)")
	case c.BancoURLJob == c.BancoURL || c.BancoURLJob == c.BancoURLAutenticacao:
		return "", "", fmt.Errorf("TENAZ_BANCO_URL_JOB não pode ser igual às URLs da aplicação nem da autenticação")
	case c.Ambiente == AmbienteProducao && strings.Contains(c.BancoURLJob, "sslmode=disable"):
		return "", "", fmt.Errorf("conexão sem TLS não é aceita em produção (RNF-032)")
	}
	return c.BancoURLJob, "", nil
}

// ExigeArmazenamento confere as variáveis do S3 do selo: sem elas o job do
// selo não roda e a verificação responde servico_indisponivel.
func (c Config) ExigeArmazenamento() error {
	a := c.Armazenamento
	var faltam []string
	for nome, v := range map[string]string{
		"TENAZ_S3_REGIAO": a.Regiao, "TENAZ_S3_BUCKET_AUDITORIA": a.BucketAuditoria,
		"TENAZ_S3_CHAVE_ACESSO": a.ChaveAcesso, "TENAZ_S3_CHAVE_SECRETA": a.ChaveSecreta,
	} {
		if v == "" {
			faltam = append(faltam, nome)
		}
	}
	if len(faltam) > 0 {
		sort.Strings(faltam)
		return fmt.Errorf("variáveis do armazenamento não definidas: %s", strings.Join(faltam, ", "))
	}
	return nil
}

// ExigeAutenticacao confere os segredos do serviço de autenticação.
func (c Config) ExigeAutenticacao() error {
	a := c.Autenticacao
	faltam := []string{}
	for nome, v := range map[string]string{
		"TENAZ_URL_WEB": a.URLWeb, "TENAZ_CHAVE_MFA": a.ChaveMFA, "TENAZ_CHAVE_JWT": a.ChaveJWT,
		"TENAZ_CHAVE_JWT_ID": a.ChaveJWTID, "TENAZ_SMTP_ENDERECO": a.SMTPEndereco, "TENAZ_EMAIL_REMETENTE": a.EmailRemetente,
	} {
		if v == "" {
			faltam = append(faltam, nome)
		}
	}
	if len(faltam) > 0 {
		sort.Strings(faltam)
		return fmt.Errorf("variáveis da autenticação não definidas: %s", strings.Join(faltam, ", "))
	}
	return nil
}

// ExigeBanco devolve erro quando o subcomando precisa de banco e não há URL.
func (c Config) ExigeBanco() error {
	if c.BancoURL == "" {
		return fmt.Errorf("TENAZ_BANCO_URL não definida")
	}
	return nil
}

func texto(chave, padrao string) string {
	if v := os.Getenv(chave); v != "" {
		return v
	}
	return padrao
}

func inteiro(chave string, padrao int) (int, error) {
	v := os.Getenv(chave)
	if v == "" {
		return padrao, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("%s deve ser um número inteiro, recebeu %q", chave, v)
	}
	return n, nil
}
