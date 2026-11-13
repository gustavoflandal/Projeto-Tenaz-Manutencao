package resiliencia

import (
	"context"
	"errors"
	"testing"
	"time"
)

func semEspera(d *Disjuntor) *Disjuntor {
	d.dormir = func(context.Context, time.Duration) error { return nil }
	return d
}

func TestRepeteTransitorioENaoRepetePermanente(t *testing.T) {
	d := semEspera(NovoDisjuntor(Politica{Tentativas: 3, RecuoInicial: time.Millisecond, RecuoMaximo: time.Millisecond,
		FalhasParaAbrir: 10, TempoAberto: time.Minute}))
	n := 0
	err := d.Executar(context.Background(), func(context.Context) error {
		n++
		if n < 3 {
			return errors.New("transitório")
		}
		return nil
	})
	if err != nil || n != 3 {
		t.Fatalf("três tentativas até dar certo: n=%d %v", n, err)
	}
	n = 0
	perm := errors.New("endereço recusado")
	err = d.Executar(context.Background(), func(context.Context) error { n++; return Permanente(perm) })
	if !errors.Is(err, perm) || n != 1 {
		t.Fatalf("permanente não repete: n=%d %v", n, err)
	}
}

func TestDisjuntorAbreEFecha(t *testing.T) {
	agora := time.Now()
	d := semEspera(NovoDisjuntor(Politica{Tentativas: 1, FalhasParaAbrir: 2, TempoAberto: time.Minute}))
	d.agora = func() time.Time { return agora }
	falha := func(context.Context) error { return errors.New("fora") }
	_ = d.Executar(context.Background(), falha)
	_ = d.Executar(context.Background(), falha)
	if err := d.Executar(context.Background(), falha); !errors.Is(err, ErrCircuitoAberto) {
		t.Fatalf("aberto depois de 2 falhas: %v", err)
	}
	agora = agora.Add(2 * time.Minute)
	if err := d.Executar(context.Background(), func(context.Context) error { return nil }); err != nil {
		t.Fatalf("passado o tempo, deixa tentar: %v", err)
	}
}

func TestTempoLimitePorTentativa(t *testing.T) {
	d := semEspera(NovoDisjuntor(Politica{TempoLimite: 10 * time.Millisecond, Tentativas: 1, FalhasParaAbrir: 5}))
	err := d.Executar(context.Background(), func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("tempo limite: %v", err)
	}
}
