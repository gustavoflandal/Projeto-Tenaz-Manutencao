// Package resiliencia dá a toda chamada externa o que a seção 8.2 do
// documento 5 exige: tempo limite, repetição com recuo exponencial e
// disjuntor (circuit breaker).
//
// Sem estado global: cada cliente externo tem o seu Disjuntor, criado no
// construtor do cliente.
package resiliencia

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

// ErrCircuitoAberto: o disjuntor está aberto e a chamada nem foi tentada.
var ErrCircuitoAberto = errors.New("circuito aberto: dependência indisponível")

// Permanente marca um erro que não adianta repetir (endereço recusado,
// credencial errada). Repetir só gasta tempo e cota do provedor.
func Permanente(err error) error {
	if err == nil {
		return nil
	}
	return &erroPermanente{err}
}

type erroPermanente struct{ error }

func (e *erroPermanente) Unwrap() error { return e.error }

// Politica define a repetição e o disjuntor.
type Politica struct {
	TempoLimite     time.Duration // por tentativa
	Tentativas      int           // total, contando a primeira
	RecuoInicial    time.Duration
	RecuoMaximo     time.Duration
	FalhasParaAbrir int           // falhas seguidas que abrem o circuito
	TempoAberto     time.Duration // quanto tempo fica aberto antes de deixar uma tentativa passar
}

// Disjuntor guarda o estado de uma dependência.
type Disjuntor struct {
	p       Politica
	mu      sync.Mutex
	falhas  int
	abertoE time.Time // aberto até
	agora   func() time.Time
	dormir  func(context.Context, time.Duration) error
}

// NovoDisjuntor monta o disjuntor com a política.
func NovoDisjuntor(p Politica) *Disjuntor {
	if p.Tentativas < 1 {
		p.Tentativas = 1
	}
	if p.FalhasParaAbrir < 1 {
		p.FalhasParaAbrir = 5
	}
	return &Disjuntor{p: p, agora: time.Now, dormir: dormir}
}

func dormir(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (d *Disjuntor) permitir() bool {
	d.mu.Lock()
	defer d.mu.Unlock()
	return !d.agora().Before(d.abertoE)
}

func (d *Disjuntor) registrar(err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err == nil {
		d.falhas = 0
		return
	}
	d.falhas++
	if d.falhas >= d.p.FalhasParaAbrir {
		d.abertoE = d.agora().Add(d.p.TempoAberto)
		d.falhas = 0
	}
}

// Executar chama fn com tempo limite por tentativa, repete os erros
// transitórios com recuo exponencial e alimenta o disjuntor.
func (d *Disjuntor) Executar(ctx context.Context, fn func(context.Context) error) error {
	recuo := d.p.RecuoInicial
	var ultimo error
	for tentativa := 1; tentativa <= d.p.Tentativas; tentativa++ {
		if !d.permitir() {
			return ErrCircuitoAberto
		}
		err := d.tentar(ctx, fn)
		var perm *erroPermanente
		if errors.As(err, &perm) {
			// Erro de quem chamou (endereço inválido), não da dependência:
			// não conta para abrir o circuito.
			return perm.error
		}
		d.registrar(err)
		if err == nil {
			return nil
		}
		ultimo = err
		if tentativa == d.p.Tentativas || ctx.Err() != nil {
			break
		}
		// Recuo com variação aleatória de até 50%, para que clientes que
		// falharam juntos não voltem juntos. Não é uso criptográfico.
		espera := recuo/2 + time.Duration(rand.Int64N(int64(recuo/2)+1)) //nolint:gosec // jitter de recuo
		if err := d.dormir(ctx, espera); err != nil {
			return fmt.Errorf("%w (depois de: %w)", err, ultimo)
		}
		recuo = min(recuo*2, d.p.RecuoMaximo)
	}
	return ultimo
}

func (d *Disjuntor) tentar(ctx context.Context, fn func(context.Context) error) error {
	if d.p.TempoLimite <= 0 {
		return fn(ctx)
	}
	ctxT, cancelar := context.WithTimeout(ctx, d.p.TempoLimite)
	defer cancelar()
	return fn(ctxT)
}
