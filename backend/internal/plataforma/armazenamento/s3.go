package armazenamento

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"

	"tenaz/internal/plataforma/resiliencia"
)

// ConfigS3 é a configuração do cliente S3. Nada tem valor padrão que
// funcione fora do desenvolvimento: endpoint, região, bucket e credencial
// vêm do ambiente (internal/plataforma/config).
type ConfigS3 struct {
	Endpoint         string // vazio: o endpoint padrão da AWS para a região
	Regiao           string
	CaminhoExplicito bool // path-style: obrigatório no SeaweedFS; a AWS aceita os dois
	Bucket           string
	ChaveAcesso      string
	ChaveSecreta     string
}

// Validar recusa configuração incompleta.
func (c ConfigS3) Validar() error {
	var falta []string
	for nome, v := range map[string]string{
		"regiao": c.Regiao, "bucket": c.Bucket, "chave de acesso": c.ChaveAcesso, "chave secreta": c.ChaveSecreta,
	} {
		if v == "" {
			falta = append(falta, nome)
		}
	}
	if len(falta) > 0 {
		return fmt.Errorf("armazenamento S3 sem: %v", falta)
	}
	return nil
}

// S3 implementa Armazenamento sobre o aws-sdk-go-v2.
type S3 struct {
	cli       *s3.Client
	bucket    string
	disjuntor *resiliencia.Disjuntor
}

// Política de toda chamada ao armazenamento (seção 8.2 do documento 5): 10 s
// por tentativa, três tentativas com recuo, circuito aberto por 30 s depois de
// cinco falhas seguidas.
func politicaPadrao() resiliencia.Politica {
	return resiliencia.Politica{
		TempoLimite: 10 * time.Second, Tentativas: 3,
		RecuoInicial: 200 * time.Millisecond, RecuoMaximo: 2 * time.Second,
		FalhasParaAbrir: 5, TempoAberto: 30 * time.Second,
	}
}

// NovoS3 monta o cliente. O disjuntor é próprio deste cliente.
func NovoS3(c ConfigS3) (*S3, error) {
	if err := c.Validar(); err != nil {
		return nil, err
	}
	cli := s3.New(s3.Options{
		Region:       c.Regiao,
		Credentials:  credentials.NewStaticCredentialsProvider(c.ChaveAcesso, c.ChaveSecreta, ""),
		BaseEndpoint: ptrOuNil(c.Endpoint),
		UsePathStyle: c.CaminhoExplicito,
		// Só soma e confere checksum quando o servidor exige: servidores
		// compatíveis com S3 não aceitam todos os trailers do SDK novo.
		RequestChecksumCalculation: aws.RequestChecksumCalculationWhenRequired,
		ResponseChecksumValidation: aws.ResponseChecksumValidationWhenRequired,
		// A repetição é da política da Tenaz (resiliencia), não do SDK.
		RetryMaxAttempts: 1,
	})
	return &S3{cli: cli, bucket: c.Bucket, disjuntor: resiliencia.NovoDisjuntor(politicaPadrao())}, nil
}

func ptrOuNil(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// GravarImutavel implementa Armazenamento.
func (s *S3) GravarImutavel(ctx context.Context, chave string, corpo []byte, retencaoAte time.Time) error {
	return s.disjuntor.Executar(ctx, func(ctx context.Context) error {
		_, err := s.cli.PutObject(ctx, &s3.PutObjectInput{
			Bucket:                    &s.bucket,
			Key:                       &chave,
			Body:                      bytes.NewReader(corpo),
			ContentType:               aws.String("application/json"),
			ObjectLockMode:            types.ObjectLockModeCompliance,
			ObjectLockRetainUntilDate: aws.Time(retencaoAte.UTC()),
		})
		return classificar(err, "gravar "+chave)
	})
}

// Existe implementa Armazenamento.
func (s *S3) Existe(ctx context.Context, chave string) (bool, error) {
	existe := false
	err := s.disjuntor.Executar(ctx, func(ctx context.Context) error {
		_, err := s.cli.HeadObject(ctx, &s3.HeadObjectInput{Bucket: &s.bucket, Key: &chave})
		if ehNaoEncontrado(err) {
			existe = false
			return nil
		}
		if err != nil {
			return classificar(err, "conferir "+chave)
		}
		existe = true
		return nil
	})
	return existe, err
}

// Ler implementa Armazenamento.
func (s *S3) Ler(ctx context.Context, chave string) ([]byte, error) {
	var corpo []byte
	err := s.disjuntor.Executar(ctx, func(ctx context.Context) error {
		r, err := s.cli.GetObject(ctx, &s3.GetObjectInput{Bucket: &s.bucket, Key: &chave})
		if ehNaoEncontrado(err) {
			return resiliencia.Permanente(ErrNaoEncontrado)
		}
		if err != nil {
			return classificar(err, "ler "+chave)
		}
		defer func() { _ = r.Body.Close() }()
		corpo, err = io.ReadAll(io.LimitReader(r.Body, 1<<20))
		if err != nil {
			return fmt.Errorf("ler corpo de %s: %w", chave, err)
		}
		return nil
	})
	if errors.Is(err, ErrNaoEncontrado) {
		return nil, ErrNaoEncontrado
	}
	return corpo, err
}

func ehNaoEncontrado(err error) bool {
	if err == nil {
		return false
	}
	var nf *types.NotFound
	var nsk *types.NoSuchKey
	if errors.As(err, &nf) || errors.As(err, &nsk) {
		return true
	}
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch ae.ErrorCode() {
		case "NotFound", "NoSuchKey":
			return true
		}
	}
	return false
}

// classificar embrulha o erro e marca como permanente o que repetir não
// resolve: credencial errada, acesso negado, pedido malformado.
func classificar(err error, contexto string) error {
	if err == nil {
		return nil
	}
	err = fmt.Errorf("armazenamento: %s: %w", contexto, err)
	var ae smithy.APIError
	if errors.As(err, &ae) {
		switch ae.ErrorCode() {
		case "AccessDenied", "InvalidAccessKeyId", "SignatureDoesNotMatch", "NoSuchBucket", "InvalidArgument",
			"InvalidRequest", "MalformedXML", "ObjectLockConfigurationNotFoundError":
			return resiliencia.Permanente(err)
		}
	}
	return err
}
