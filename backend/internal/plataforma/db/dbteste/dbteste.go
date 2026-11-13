// Package dbteste prepara um PostgreSQL real para os testes de integração,
// no mesmo padrão da suíte de vazamento: contêiner descartável, ou o banco de
// TENAZ_TESTE_BANCO_URL quando definida (ciclo local mais rápido; o banco
// precisa estar com as migrações atuais).
//
// Os pools de teste assumem os papéis de produção por SET ROLE sobre a URL do
// migrador: o da aplicação como tenaz_app (RLS valendo) e o da autenticação
// como tenaz_autenticacao (ADR-011). Só código de teste importa este pacote.
package dbteste

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"

	"tenaz/internal/plataforma/armazenamento"
	"tenaz/internal/plataforma/auditoria"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/migracao"
)

// ImagemPostgres é a imagem do PostgreSQL de teste: a mesma do compose de
// desenvolvimento (infra/docker), fixada por versão e digest, para que o
// teste, o desenvolvimento e o CI rodem o mesmo binário. Muda junto com o
// compose. Exportada para a suíte de vazamento usar a mesma.
const ImagemPostgres = "postgres:17.11-alpine3.24@sha256:e5d2edd28b1f8d40a0a6c0c9cdc545b19b56a0d5d88f94bc11a97d970d4700d0"

// Ambiente são os três pools de um teste de integração.
type Ambiente struct {
	URL      string        // migrador
	App      *db.DB        // tenaz_app, só NaEmpresa
	Auth     *pgxpool.Pool // tenaz_autenticacao
	Job      *db.DB        // tenaz_job: jobs (selo da trilha), só NaEmpresa e a listagem de empresas
	Migrador *pgxpool.Pool // superusuário: só para preparar e inspecionar o cenário
	encerrar func()
}

// Encerrar fecha os pools e o contêiner.
func (a *Ambiente) Encerrar() {
	if a.App != nil {
		a.App.Fechar()
	}
	if a.Auth != nil {
		a.Auth.Close()
	}
	if a.Job != nil {
		a.Job.Fechar()
	}
	if a.Migrador != nil {
		a.Migrador.Close()
	}
	if a.encerrar != nil {
		a.encerrar()
	}
}

// Preparar sobe o banco, aplica as migrações e abre os pools.
func Preparar(ctx context.Context) (*Ambiente, error) {
	url, encerrar, err := subir(ctx)
	if err != nil {
		return nil, err
	}
	amb := &Ambiente{URL: url, encerrar: encerrar}
	if err := migracao.Aplicar(ctx, url); err != nil {
		amb.Encerrar()
		return nil, fmt.Errorf("aplicar migrações: %w", err)
	}
	if amb.Migrador, err = pgxpool.New(ctx, url); err != nil {
		amb.Encerrar()
		return nil, err
	}
	if amb.App, err = db.AbrirComOpcoes(ctx, db.Opcoes{URL: url, MaxConn: 8, Papel: "tenaz_app"}); err != nil {
		amb.Encerrar()
		return nil, err
	}
	if amb.Job, err = db.AbrirComOpcoes(ctx, db.Opcoes{URL: url, MaxConn: 4, Papel: "tenaz_job"}); err != nil {
		amb.Encerrar()
		return nil, fmt.Errorf("pool do job (o papel tenaz_job existe?): %w", err)
	}
	if amb.Auth, err = db.AbrirPool(ctx, db.Opcoes{URL: url, MaxConn: 4, Papel: "tenaz_autenticacao"}); err != nil {
		amb.Encerrar()
		return nil, fmt.Errorf("pool da autenticação (o papel tenaz_autenticacao existe?): %w", err)
	}
	return amb, nil
}

func subir(ctx context.Context) (string, func(), error) {
	if url := os.Getenv("TENAZ_TESTE_BANCO_URL"); url != "" {
		return url, func() {}, nil
	}
	// go test ./... roda os pacotes em processos paralelos, e vários deles
	// sobem um contêiner. No Docker Desktop do Windows, clientes criados ao
	// mesmo tempo fazem o testcontainers falhar ao detectar o Docker ("failed
	// to create Docker provider") por quase um minuto. A subida é serializada
	// entre os processos por um arquivo de trava; os testes continuam em
	// paralelo, cada pacote no seu banco.
	soltar, err := travarSubida(ctx)
	if err != nil {
		return "", nil, err
	}
	var cont *postgres.PostgresContainer
	for tentativa := 1; ; tentativa++ {
		cont, err = postgres.Run(ctx, ImagemPostgres,
			postgres.WithDatabase("tenaz"),
			postgres.WithUsername("tenaz_migrador"),
			postgres.WithPassword("tenaz_teste"),
			testcontainers.WithWaitStrategy(
				wait.ForLog("database system is ready to accept connections").
					WithOccurrence(2).
					WithStartupTimeout(90*time.Second),
			),
		)
		// Falha de detecção do Docker é transitória: repete com recuo;
		// qualquer outro erro sai na hora.
		if err == nil || tentativa == 4 || !strings.Contains(err.Error(), "failed to create Docker provider") {
			break
		}
		time.Sleep(time.Duration(tentativa) * 3 * time.Second)
	}
	soltar()
	if err != nil {
		return "", nil, fmt.Errorf("subir PostgreSQL de teste: %w", err)
	}
	encerrar := func() { _ = testcontainers.TerminateContainer(cont) }
	url, err := cont.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		encerrar()
		return "", nil, err
	}
	return url, encerrar, nil
}

// travarSubida espera a vez de subir o contêiner: cria o arquivo de trava
// com O_EXCL e o apaga ao soltar. Trava mais velha que 3 minutos é de um
// processo que morreu no meio e é descartada.
func travarSubida(ctx context.Context) (func(), error) {
	caminho := filepath.Join(os.TempDir(), "tenaz-dbteste-subida.trava")
	for {
		f, err := os.OpenFile(caminho, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o600)
		if err == nil {
			_ = f.Close()
			return func() { _ = os.Remove(caminho) }, nil
		}
		if !errors.Is(err, fs.ErrExist) {
			return nil, fmt.Errorf("trava da subida do PostgreSQL de teste: %w", err)
		}
		if info, err := os.Stat(caminho); err == nil && time.Since(info.ModTime()) > 3*time.Minute {
			_ = os.Remove(caminho)
			continue
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(250 * time.Millisecond):
		}
	}
}

// ArmazenamentoS3 devolve o cliente S3 do selo para os testes: o SeaweedFS do
// compose (ferramentas/tarefas.sh infra), com a identidade tenaz_auditoria e
// o bucket tenaz-auditoria, ou o que TENAZ_TESTE_S3_ENDPOINT,
// TENAZ_TESTE_S3_CHAVE_ACESSO e TENAZ_TESTE_S3_CHAVE_SECRETA indicarem.
// Devolve nil e o motivo quando o armazenamento não responde: o teste que o
// exige se pula, dizendo como subir a infraestrutura.
func ArmazenamentoS3(ctx context.Context) (*armazenamento.S3, armazenamento.ConfigS3, error) {
	cfg := armazenamento.ConfigS3{
		Endpoint: valorOu("TENAZ_TESTE_S3_ENDPOINT", "http://127.0.0.1:8333"), Regiao: "us-east-1",
		CaminhoExplicito: true, Bucket: valorOu("TENAZ_TESTE_S3_BUCKET", "tenaz-auditoria"),
		ChaveAcesso:  valorOu("TENAZ_TESTE_S3_CHAVE_ACESSO", "tenaz_auditoria"),
		ChaveSecreta: valorOu("TENAZ_TESTE_S3_CHAVE_SECRETA", "tenaz_auditoria_dev_secreto"),
	}
	s3, err := armazenamento.NovoS3(cfg)
	if err != nil {
		return nil, cfg, err
	}
	ctxT, cancelar := context.WithTimeout(ctx, 5*time.Second)
	defer cancelar()
	if _, err := s3.Existe(ctxT, "selos/_ping"); err != nil {
		return nil, cfg, fmt.Errorf("armazenamento S3 de teste fora do ar em %s (rode: bash ferramentas/tarefas.sh infra): %w", cfg.Endpoint, err)
	}
	return s3, cfg, nil
}

func valorOu(chave, padrao string) string {
	if v := os.Getenv(chave); v != "" {
		return v
	}
	return padrao
}

// ConferirHashesDaTrilha recalcula no Go, pelo formato 1 do documento 2
// (seção 9.2.1), o hash de todo evento da empresa e o compara com o que
// app.hash_evento_auditoria gravou. Devolve quantos conferiu e uma descrição
// de cada divergência. Lê como migrador: é o verificador externo, que recebe
// os campos como texto do banco.
func ConferirHashesDaTrilha(ctx context.Context, pool *pgxpool.Pool, empresaID uuid.UUID) (conferidos int, divergencias []string, err error) {
	linhas, err := pool.Query(ctx, `SELECT id::text, empresa_id::text, ocorrido_em, ator_tipo, usuario_id::text, usuario_nome,
  token_api_id::text, acao, resultado, sensivel, entidade, entidade_id::text, entidade_chave, antes::text, depois::text,
  detalhe::text, to_jsonb(campos_alterados)::text, modulo, tela, origem, ip::text, dispositivo_id::text, user_agent,
  request_id, conexao_empresa_id::text, ocorrido_no_dispositivo_em, escopo, transacao::text, hash
FROM app.evento_auditoria WHERE empresa_id = $1 ORDER BY ocorrido_em, id`, empresaID)
	if err != nil {
		return 0, nil, err
	}
	defer linhas.Close()
	for linhas.Next() {
		var (
			e    auditoria.ConteudoEvento
			hash []byte
		)
		if err := linhas.Scan(&e.ID, &e.EmpresaID, &e.OcorridoEm, &e.AtorTipo, &e.UsuarioID, &e.UsuarioNome,
			&e.TokenAPIID, &e.Acao, &e.Resultado, &e.Sensivel, &e.Entidade, &e.EntidadeID, &e.EntidadeChave, &e.Antes,
			&e.Depois, &e.Detalhe, &e.CamposAlterados, &e.Modulo, &e.Tela, &e.Origem, &e.IP, &e.DispositivoID,
			&e.UserAgent, &e.RequestID, &e.ConexaoEmpresaID, &e.OcorridoNoDispositivoEm, &e.Escopo, &e.Transacao,
			&hash); err != nil {
			return conferidos, divergencias, err
		}
		conferidos++
		if h := e.Hash(); !bytes.Equal(h[:], hash) {
			divergencias = append(divergencias, fmt.Sprintf("evento %s (%s/%s): banco %x, Go %x sobre %s",
				e.ID, e.Entidade, e.Acao, hash, h, e.Serializar()))
		}
	}
	return conferidos, divergencias, linhas.Err()
}
