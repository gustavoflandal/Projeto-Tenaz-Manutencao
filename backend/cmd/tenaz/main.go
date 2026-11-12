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
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"tenaz/internal/plataforma/api"
	"tenaz/internal/plataforma/config"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/migracao"
	"tenaz/internal/plataforma/obs"
	"tenaz/internal/plataforma/versao"
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
	_, a := api.Novo()
	yaml, err := api.EscreverOpenAPI(a)
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
	var banco *db.DB
	if cfg.BancoURL != "" {
		var err error
		banco, err = db.Abrir(ctx, cfg.BancoURL, cfg.BancoMaxConn)
		if err != nil {
			return err
		}
		defer banco.Fechar()
		registro.Info("banco conectado")
	} else {
		registro.Warn("TENAZ_BANCO_URL não definida: subindo sem banco")
	}

	mux, _ := api.Novo()

	servidor := &http.Server{
		Addr:              cfg.EnderecoHTTP,
		Handler:           mux,
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
		if err := servidor.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
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
	if err := servidor.Shutdown(prazo); err != nil {
		return fmt.Errorf("encerrar servidor: %w", err)
	}
	return nil
}
