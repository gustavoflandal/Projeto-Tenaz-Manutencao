package auditoria_test

import (
	"context"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"

	"tenaz/internal/plataforma/armazenamento"
	"tenaz/internal/plataforma/auditoria"
	"tenaz/internal/plataforma/db/dbteste"
)

// O selo no armazenamento imutável (RF-ADM-011, critério 5; ADR-012): o
// objeto grava com retenção em modo de conformidade e o servidor recusa
// apagá-lo, encurtar a retenção ou trocar o modo. Sobrescrever a chave cria
// uma versão nova, e a protegida continua intacta; a verificação lê a versão
// corrente e por isso pega a troca (copia_divergente). Roda contra o SeaweedFS
// do compose (bash ferramentas/tarefas.sh infra); sem ele, se pula.
func TestSeloNoArmazenamentoImutavel(t *testing.T) {
	ctx := context.Background()
	s3c, cfg, err := dbteste.ArmazenamentoS3(ctx)
	if err != nil {
		t.Skipf("sem armazenamento S3: %v", err)
	}
	cli := s3.New(s3.Options{
		Region: cfg.Regiao, BaseEndpoint: aws.String(cfg.Endpoint), UsePathStyle: cfg.CaminhoExplicito,
		Credentials:                credentials.NewStaticCredentialsProvider(cfg.ChaveAcesso, cfg.ChaveSecreta, ""),
		RequestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
		ResponseChecksumValidation: aws.ResponseChecksumValidationWhenRequired,
	})
	chave := "selos/empresas/" + uuid.NewString() + "/2026-01-01.selo"
	original := []byte(`{"formato":1,"quantidade":3}`)
	// Retenção curta de propósito: o objeto de conformidade não se apaga antes
	// do prazo, nem pelo teste.
	ate := time.Now().UTC().Add(5 * time.Minute)

	if ok, err := s3c.Existe(ctx, chave); err != nil || ok {
		t.Fatalf("antes de gravar: existe %v, %v", ok, err)
	}
	if _, err := s3c.Ler(ctx, chave); !errors.Is(err, armazenamento.ErrNaoEncontrado) {
		t.Fatalf("ler objeto ausente: %v", err)
	}
	if err := s3c.GravarImutavel(ctx, chave, original, ate); err != nil {
		t.Fatalf("gravar: %v", err)
	}
	if ok, err := s3c.Existe(ctx, chave); err != nil || !ok {
		t.Fatalf("depois de gravar: existe %v, %v", ok, err)
	}
	corpo, err := s3c.Ler(ctx, chave)
	if err != nil || string(corpo) != string(original) {
		t.Fatalf("ler: %q %v", corpo, err)
	}

	// retenção gravada em modo de conformidade
	versoes, err := cli.ListObjectVersions(ctx, &s3.ListObjectVersionsInput{Bucket: &cfg.Bucket, Prefix: &chave})
	if err != nil || len(versoes.Versions) != 1 {
		t.Fatalf("versões: %v %v", versoes, err)
	}
	versao := versoes.Versions[0].VersionId
	ret, err := cli.GetObjectRetention(ctx, &s3.GetObjectRetentionInput{Bucket: &cfg.Bucket, Key: &chave, VersionId: versao})
	if err != nil || ret.Retention == nil || ret.Retention.Mode != types.ObjectLockRetentionModeCompliance {
		t.Fatalf("retenção: %+v %v", ret, err)
	}
	if d := ret.Retention.RetainUntilDate.Sub(ate); d < -2*time.Second || d > 2*time.Second {
		t.Errorf("retain-until %v, esperado %v", ret.Retention.RetainUntilDate, ate)
	}

	negado := func(nome string, err error) {
		t.Helper()
		if err == nil {
			t.Errorf("%s: o servidor deveria recusar", nome)
			return
		}
		if !strings.Contains(err.Error(), "AccessDenied") && !strings.Contains(err.Error(), "Forbidden") {
			t.Errorf("%s: recusa esperada como AccessDenied, veio %v", nome, err)
		}
	}
	_, err = cli.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: &cfg.Bucket, Key: &chave, VersionId: versao})
	negado("apagar a versão protegida", err)
	_, err = cli.DeleteObject(ctx, &s3.DeleteObjectInput{Bucket: &cfg.Bucket, Key: &chave, VersionId: versao, BypassGovernanceRetention: aws.Bool(true)})
	negado("apagar com bypass de governança", err)
	_, err = cli.PutObjectRetention(ctx, &s3.PutObjectRetentionInput{Bucket: &cfg.Bucket, Key: &chave, VersionId: versao,
		Retention: &types.ObjectLockRetention{Mode: types.ObjectLockRetentionModeCompliance, RetainUntilDate: aws.Time(time.Now().UTC().Add(time.Minute))}})
	negado("encurtar a retenção", err)
	_, err = cli.PutObjectRetention(ctx, &s3.PutObjectRetentionInput{Bucket: &cfg.Bucket, Key: &chave, VersionId: versao,
		Retention: &types.ObjectLockRetention{Mode: types.ObjectLockRetentionModeGovernance, RetainUntilDate: aws.Time(ate.Add(time.Hour))}})
	negado("trocar para governança", err)

	// Sobrescrever a chave cria uma versão nova; a protegida fica intacta e a
	// verificação, que lê a corrente, enxerga a troca.
	forjado := []byte(`{"formato":1,"quantidade":0}`)
	if err := s3c.GravarImutavel(ctx, chave, forjado, ate); err != nil {
		t.Fatalf("sobrescrever: %v", err)
	}
	corrente, err := s3c.Ler(ctx, chave)
	if err != nil || string(corrente) != string(forjado) {
		t.Fatalf("a versão corrente é a nova: %q %v", corrente, err)
	}
	antiga, err := cli.GetObject(ctx, &s3.GetObjectInput{Bucket: &cfg.Bucket, Key: &chave, VersionId: versao})
	if err != nil {
		t.Fatalf("a versão original sumiu: %v", err)
	}
	defer func() { _ = antiga.Body.Close() }()
	buf := make([]byte, 128)
	n, _ := antiga.Body.Read(buf)
	if string(buf[:n]) != string(original) {
		t.Errorf("a versão protegida mudou: %q", buf[:n])
	}
}

// Ponta a ponta com o servidor S3 de verdade: o job sela e copia, a
// verificação lê a cópia e dá íntegro, e a troca da cópia por outro conteúdo
// dá violado/copia_divergente.
func TestSeloPontaAPontaNoS3(t *testing.T) {
	ctx := context.Background()
	s3c, _, err := dbteste.ArmazenamentoS3(ctx)
	if err != nil {
		t.Skipf("sem armazenamento S3: %v", err)
	}
	amb, err := dbteste.Preparar(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(amb.Encerrar)

	empresa := uuid.Must(uuid.NewV7())
	if _, err := amb.Migrador.Exec(ctx, `INSERT INTO app.empresa (id, razao_social, retencao_auditoria_anos) VALUES ($1, 'Selo no S3', 5)`, empresa); err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 2; i++ {
		if _, err := amb.Migrador.Exec(ctx, `INSERT INTO app.evento_auditoria (empresa_id, ocorrido_em, ator_tipo, acao, entidade, modulo)
VALUES ($1, $2, 'sistema', 'inclusao', 'unidade', 'adm')`, empresa, dia(-3).Add(time.Duration(i+1)*time.Hour)); err != nil {
			t.Fatal(err)
		}
	}
	sel, err := auditoria.NovoSelador(auditoria.DependenciasSelador{Banco: amb.Job, Armazenamento: s3c})
	if err != nil {
		t.Fatal(err)
	}
	selados, enviados, err := sel.Empresa(ctx, empresa)
	if err != nil || selados < 2 || enviados != selados {
		t.Fatalf("job: selados %d, enviados %d, %v", selados, enviados, err)
	}
	v, err := auditoria.NovoVerificador(amb.App, s3c)
	if err != nil {
		t.Fatal(err)
	}
	dias, err := v.Verificar(ctx, empresa, dia(-3), dia(-2))
	if err != nil || len(dias) != 2 || dias[0].Estado != auditoria.EstadoIntegro || dias[1].Estado != auditoria.EstadoIntegro {
		t.Fatalf("verificação com a cópia: %+v %v", dias, err)
	}
	// o objeto é o que o banco calcula
	chave := chaveDoSelo(empresa, dia(-3))
	if err := s3c.GravarImutavel(ctx, chave, []byte(`{"formato":1,"forjado":true}`), time.Now().UTC().Add(5*time.Minute)); err != nil {
		t.Fatal(err)
	}
	dias, err = v.Verificar(ctx, empresa, dia(-3), dia(-3))
	if err != nil || dias[0].Estado != auditoria.EstadoViolado || dias[0].Motivo == nil || *dias[0].Motivo != auditoria.MotivoCopiaDivergente {
		t.Fatalf("cópia trocada: %+v %v", dias, err)
	}
}
