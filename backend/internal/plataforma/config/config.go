// Package config lê a configuração do processo a partir do ambiente.
//
// Nenhum segredo tem valor padrão: ausência de segredo é erro na partida, não
// um valor conveniente que vaza para produção.
package config

import (
	"fmt"
	"os"
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
