// Package armazenamento é a porta da plataforma para o armazenamento de
// objetos compatível com S3 (ADR-012, decisão 5 do documento 5): o backend
// fala só S3, com endpoint, região e estilo de caminho por configuração, de
// modo que trocar o servidor (SeaweedFS no desenvolvimento, S3 da AWS em
// produção) é trocar configuração. Nenhuma API própria de servidor, nenhum
// minio-go.
//
// Hoje a porta serve ao selo diário da trilha de auditoria (RF-ADM-011): a
// gravação é imutável, em modo de conformidade (Object Lock COMPLIANCE), e a
// leitura serve à verificação de integridade.
package armazenamento

import (
	"context"
	"errors"
	"time"
)

// ErrNaoEncontrado: o objeto não existe.
var ErrNaoEncontrado = errors.New("objeto não encontrado")

// Armazenamento é o que a plataforma usa do armazenamento de objetos.
type Armazenamento interface {
	// GravarImutavel grava o objeto com Object Lock em modo de conformidade,
	// retido até retencaoAte. Gravar de novo a mesma chave cria uma versão
	// nova (o bucket é versionado); a versão protegida continua intacta.
	GravarImutavel(ctx context.Context, chave string, corpo []byte, retencaoAte time.Time) error
	// Existe confere por HEAD se há objeto na chave.
	Existe(ctx context.Context, chave string) (bool, error)
	// Ler devolve a versão corrente do objeto; ErrNaoEncontrado se não há.
	Ler(ctx context.Context, chave string) ([]byte, error)
}
