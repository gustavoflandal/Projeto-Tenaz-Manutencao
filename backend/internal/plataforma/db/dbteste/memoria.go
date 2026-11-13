package dbteste

import (
	"context"
	"errors"
	"sync"
	"time"

	"tenaz/internal/plataforma/armazenamento"
)

// Memoria é um armazenamento de objetos em memória, para os testes do selo da
// trilha que não precisam do SeaweedFS: guarda o corpo e a retenção de cada
// chave, pode falhar na gravação e deixa o teste apagar uma cópia.
type Memoria struct {
	mu          sync.Mutex
	objetos     map[string][]byte
	retencao    map[string]time.Time
	FalhaGravar bool // quando verdadeiro, GravarImutavel devolve erro
}

var _ armazenamento.Armazenamento = (*Memoria)(nil)

// NovaMemoria devolve um armazenamento vazio.
func NovaMemoria() *Memoria {
	return &Memoria{objetos: map[string][]byte{}, retencao: map[string]time.Time{}}
}

// GravarImutavel implementa armazenamento.Armazenamento.
func (m *Memoria) GravarImutavel(_ context.Context, chave string, corpo []byte, ate time.Time) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.FalhaGravar {
		return errors.New("armazenamento fora do ar")
	}
	m.objetos[chave] = append([]byte(nil), corpo...)
	m.retencao[chave] = ate
	return nil
}

// Existe implementa armazenamento.Armazenamento.
func (m *Memoria) Existe(_ context.Context, chave string) (bool, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	_, ok := m.objetos[chave]
	return ok, nil
}

// Ler implementa armazenamento.Armazenamento.
func (m *Memoria) Ler(_ context.Context, chave string) ([]byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	b, ok := m.objetos[chave]
	if !ok {
		return nil, armazenamento.ErrNaoEncontrado
	}
	return append([]byte(nil), b...), nil
}

// Apagar tira a cópia, como quem a tivesse removido do armazenamento.
func (m *Memoria) Apagar(chave string) {
	m.mu.Lock()
	delete(m.objetos, chave)
	m.mu.Unlock()
}

// Retencao devolve até quando a chave ficou retida.
func (m *Memoria) Retencao(chave string) time.Time {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.retencao[chave]
}
