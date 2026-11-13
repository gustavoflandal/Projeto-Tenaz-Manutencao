// Comando tenaz: um binário, vários modos.
//
// Subcomandos (seção 4.1 de docs/5_ARQUITETURA_SOFTWARE.md):
//
//	api            servidor HTTP
//	worker         processador da fila de jobs
//	sincronizador  sincronização entre contas conectadas (SRV)
//	ingestor       ingestão de telemática e IoT
//	migrate        aplica as migrações do banco
//	openapi        escreve o contrato gerado na saída padrão
//	versao         imprime versão, commit e data do build
package main

import (
	"context"
	"crypto/ed25519"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"tenaz/internal/modulos/adm"
	"tenaz/internal/plataforma/api"
	"tenaz/internal/plataforma/autenticacao"
	"tenaz/internal/plataforma/config"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/email"
	"tenaz/internal/plataforma/migracao"
	"tenaz/internal/plataforma/obs"
	"tenaz/internal/plataforma/versao"
	"tenaz/internal/servidor"
)

func main() {
	if err := executar(); err != nil {
		fmt.Fprintf(os.Stderr, "tenaz: %v\n", err)
		os.Exit(1)
	}
}

func executar() error {
	if len(os.Args) < 2 {
		uso()
		return errors.New("subcomando não informado")
	}

	sub := os.Args[1]

	// "versao" não carrega configuração: tem de funcionar em qualquer
	// ambiente, inclusive numa imagem sem variável definida, porque o CI o
	// usa para conferir a versão embutida contra a tag.
	if sub == "versao" {
		return imprimirVersao()
	}

	cfg, err := config.DeAmbiente()
	if err != nil {
		return err
	}
	versao.Ambiente = cfg.Ambiente

	// Produção não aceita versão de pré-lançamento (seção 9.5, regra 5).
	if cfg.Ambiente == config.AmbienteProducao && versao.PreLancamento() {
		return fmt.Errorf("versão de pré-lançamento (%s) não é aceita em produção", versao.Versao)
	}

	registro := obs.Novo(os.Stdout, sub, cfg.NivelLog, cfg.LogJSON)

	ctx, cancelar := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancelar()
	ctx = obs.ComLogger(ctx, registro)

	switch sub {
	case "api":
		return servirAPI(ctx, cfg)
	case "openapi":
		return escreverOpenAPI()
	case "worker", "sincronizador", "ingestor":
		return fmt.Errorf("subcomando %q ainda não implementado; chega nas etapas seguintes", sub)
	case "migrate":
		if err := cfg.ExigeBanco(); err != nil {
			return err
		}
		if err := migracao.Aplicar(ctx, cfg.BancoURL); err != nil {
			return err
		}
		v, err := migracao.Versao(ctx, cfg.BancoURL)
		if err != nil {
			return err
		}
		registro.Info("migrações aplicadas", "versao_esquema", v)
		return nil
	case "migrate-status":
		if err := cfg.ExigeBanco(); err != nil {
			return err
		}
		return migracao.Situacao(ctx, cfg.BancoURL)
	default:
		uso()
		return fmt.Errorf("subcomando desconhecido: %q", sub)
	}
}

func uso() {
	fmt.Fprint(os.Stderr, `uso: tenaz <subcomando>

  api            servidor HTTP
  worker         processador da fila de jobs
  sincronizador  sincronização entre contas conectadas
  ingestor       ingestão de telemática e IoT
  migrate        aplica as migrações do banco
  migrate-status mostra o estado das migrações
  openapi        escreve o contrato OpenAPI na saída padrão
  versao         imprime versão, commit e data do build
`)
}

func imprimirVersao() error {
	d := versao.Interno()
	fmt.Printf("versao=%s\ncommit=%s\ndata_build=%s\nambiente=%s\n",
		d.Versao, d.Commit, d.DataBuild, d.Ambiente)

	if ok, motivo := versao.Coerente(); !ok {
		fmt.Fprintf(os.Stderr, "AVISO: versão incoerente: %s\n", motivo)
	}
	return nil
}

func escreverOpenAPI() error {
	yaml, err := api.EscreverOpenAPI(servidor.Contrato().API())
	if err != nil {
		return fmt.Errorf("gerar OpenAPI: %w", err)
	}
	_, err = os.Stdout.Write(yaml)
	return err
}

func servirAPI(ctx context.Context, cfg config.Config) error {
	registro := obs.Logger(ctx)

	// O banco é opcional nesta etapa: GET /v1/versao tem de responder mesmo
	// sem banco, porque é o que a verificação sintética chama depois do deploy
	// para confirmar a versão publicada.
	var (
		banco *db.DB
		deps  servidor.Dependencias
	)
	if cfg.BancoURL != "" {
		var err error
		banco, err = db.AbrirComOpcoes(ctx, db.Opcoes{URL: cfg.BancoURL, MaxConn: cfg.BancoMaxConn, Papel: cfg.BancoPapel})
		if err != nil {
			return err
		}
		defer banco.Fechar()
		registro.Info("banco conectado")

		// Pool próprio das credenciais (ADR-011): só no modo api, nunca no
		// worker.
		autent, fechar, err := montarAutenticacao(ctx, cfg, banco)
		if err != nil {
			return err
		}
		defer fechar()
		registro.Info("autenticação pronta")

		modADM, err := adm.Novo(adm.Dependencias{Banco: banco, Autenticacao: autent})
		if err != nil {
			return err
		}
		deps = servidor.Dependencias{Autenticacao: autent, ADM: modADM}
	} else {
		registro.Warn("TENAZ_BANCO_URL não definida: subindo sem banco; só GET /v1/versao responde")
	}

	origens := cfg.OrigensWeb
	if len(origens) == 0 {
		if o := api.OrigemDe(cfg.Autenticacao.URLWeb); o != "" {
			origens = []string{o}
		}
	}
	montado := servidor.Montar(api.Config{
		CookieSeguro: cfg.CookieSeguro, OrigensWeb: origens, URLWeb: strings.TrimRight(cfg.Autenticacao.URLWeb, "/"),
	}, deps)

	srv := &http.Server{
		Addr:              cfg.EnderecoHTTP,
		Handler:           montado.Handler(),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       cfg.TempoLimiteHTTP,
		WriteTimeout:      cfg.TempoLimiteHTTP,
		IdleTimeout:       120 * time.Second,
	}

	erroServidor := make(chan error, 1)
	go func() {
		// A versão já vem do logger do processo; repeti-la aqui duplicava o
		// campo na linha de log.
		registro.Info("servidor HTTP no ar",
			"endereco", cfg.EnderecoHTTP,
			"ambiente", cfg.Ambiente,
		)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			erroServidor <- err
		}
	}()

	select {
	case err := <-erroServidor:
		return fmt.Errorf("servidor HTTP: %w", err)
	case <-ctx.Done():
		registro.Info("encerrando")
	}

	// Desligamento com prazo, para não cortar requisição em andamento.
	prazo, cancelar := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancelar()
	if err := srv.Shutdown(prazo); err != nil {
		return fmt.Errorf("encerrar servidor: %w", err)
	}
	return nil
}

// montarAutenticacao abre o pool da autenticação e monta o serviço com as
// chaves do ambiente. Em desenvolvimento o e-mail vai para o Mailpit.
func montarAutenticacao(ctx context.Context, cfg config.Config, banco *db.DB) (*autenticacao.Servico, func(), error) {
	if err := cfg.ExigeAutenticacao(); err != nil {
		return nil, nil, err
	}
	url, papel, err := cfg.BancoAutenticacao()
	if err != nil {
		return nil, nil, err
	}
	pool, err := db.AbrirPool(ctx, db.Opcoes{URL: url, MaxConn: cfg.BancoAutenticacaoMaxConn, Papel: papel})
	if err != nil {
		return nil, nil, fmt.Errorf("pool da autenticação: %w", err)
	}
	falhar := func(err error) (*autenticacao.Servico, func(), error) {
		pool.Close()
		return nil, nil, err
	}

	a := cfg.Autenticacao
	chaveMFA, err := base64.StdEncoding.DecodeString(a.ChaveMFA)
	if err != nil {
		return falhar(fmt.Errorf("TENAZ_CHAVE_MFA: %w", err))
	}
	cifrador, err := autenticacao.NovoCifrador(chaveMFA)
	if err != nil {
		return falhar(err)
	}
	semente, err := base64.StdEncoding.DecodeString(a.ChaveJWT)
	if err != nil {
		return falhar(fmt.Errorf("TENAZ_CHAVE_JWT: %w", err))
	}
	anteriores := map[string]ed25519.PublicKey{}
	for _, par := range strings.Split(a.ChavesJWTAnteriores, ",") {
		if strings.TrimSpace(par) == "" {
			continue
		}
		kid, b64, ok := strings.Cut(strings.TrimSpace(par), "=")
		pub, errDec := base64.StdEncoding.DecodeString(b64)
		if !ok || errDec != nil {
			return falhar(fmt.Errorf("TENAZ_CHAVES_JWT_ANTERIORES: par inválido %q", kid))
		}
		anteriores[kid] = ed25519.PublicKey(pub)
	}
	assinador, err := autenticacao.NovoAssinador(a.ChaveJWTID, semente, anteriores)
	if err != nil {
		return falhar(err)
	}
	hasher, err := autenticacao.NovoHasher(autenticacao.ParametrosPadrao(), 4)
	if err != nil {
		return falhar(err)
	}
	smtp, err := email.NovoSMTP(email.ConfigSMTP{
		Endereco: a.SMTPEndereco, Usuario: a.SMTPUsuario, Senha: a.SMTPSenha,
		ExigirTLS: cfg.Ambiente != config.AmbienteDesenvolvimento,
		Remetente: a.EmailRemetente, NomeRemetente: "Tenaz",
	})
	if err != nil {
		return falhar(err)
	}
	s, err := autenticacao.Novo(autenticacao.ConfigPadrao(a.URLWeb), autenticacao.Dependencias{
		Banco: banco, PoolAutenticacao: pool, Hasher: hasher, Politica: autenticacao.NovaPolitica(),
		Cifrador: cifrador, Assinador: assinador, Email: smtp,
	})
	if err != nil {
		return falhar(err)
	}
	return s, func() { s.Aguardar(); pool.Close() }, nil
}
