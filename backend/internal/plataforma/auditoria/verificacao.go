package auditoria

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"tenaz/internal/plataforma/armazenamento"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/db/consultas"
	"tenaz/internal/plataforma/erros"
)

// Estados da verificação de integridade de um dia (documento 2, seção 9.4).
// O estado é calculado a cada verificação e nunca gravado.
const (
	EstadoIntegro         = "integro"
	EstadoViolado         = "violado"
	EstadoNaoSelado       = "nao_selado"
	MotivoCopiaAusente    = "copia_ausente"
	MotivoCopiaDivergente = "copia_divergente"
)

// DiaVerificado é o resultado de um dia.
type DiaVerificado struct {
	Dia      time.Time // dia UTC
	Estado   string
	Motivo   *string
	EventoID *uuid.UUID // primeiro evento alterado, quando há como apontar
}

// ErroPeriodo: período fora do aceito pela verificação (de..ate, até 400 dias).
var ErroPeriodo = erros.EntradaInvalida.ComCampoChave("ate", erros.CampoValorInvalido, "campo.periodo_verificacao", nil)

const maxDiasVerificacao = 400

// Verificador junta a recontagem do banco com a conferência da cópia no
// armazenamento imutável (RF-ADM-011).
type Verificador struct {
	banco *db.DB
	arm   armazenamento.Armazenamento
}

// NovoVerificador monta o verificador. Sem armazenamento, a verificação
// responde servico_indisponivel: sem a cópia não há como afirmar que o período
// está íntegro.
func NovoVerificador(banco *db.DB, arm armazenamento.Armazenamento) (*Verificador, error) {
	if banco == nil {
		return nil, errors.New("verificador: banco ausente")
	}
	return &Verificador{banco: banco, arm: arm}, nil
}

// Verificar confere os dias de de a ate (inclusive, dias UTC) da empresa. Em
// duas etapas: o banco recalcula tudo (app.verificar_integridade_auditoria,
// numa transação de leitura da empresa) e, fora da transação, cada dia íntegro
// tem o selo comparado com o objeto do armazenamento: objeto ausente dá
// violado/copia_ausente, e objeto diferente, violado/copia_divergente. Só a
// comparação com a cópia pega quem, com acesso total ao banco, refez eventos e
// selos coerentes entre si.
func (v *Verificador) Verificar(ctx context.Context, empresaID uuid.UUID, de, ate time.Time) ([]DiaVerificado, error) {
	de, ate = diaUTC(de), diaUTC(ate)
	if ate.Before(de) || ate.Sub(de) > maxDiasVerificacao*24*time.Hour {
		return nil, ErroPeriodo
	}
	if v.arm == nil {
		return nil, erros.ServicoIndisponivel.ComChave("detalhe.selo_nao_configurado", nil)
	}

	var (
		dias  []DiaVerificado
		selos = map[string]Selo{}
	)
	err := v.banco.NaEmpresa(ctx, empresaID, db.Interno(nil), func(ctx context.Context) error {
		tx, err := db.Tx(ctx)
		if err != nil {
			return err
		}
		rows, err := tx.Query(ctx,
			`SELECT dia, estado, motivo, evento_id FROM app.verificar_integridade_auditoria($1::date, $2::date)`, de, ate)
		if err != nil {
			return fmt.Errorf("verificar integridade da trilha: %w", err)
		}
		defer rows.Close()
		for rows.Next() {
			var (
				d      pgtype.Date
				estado string
				motivo *string
				evento *uuid.UUID
			)
			if err := rows.Scan(&d, &estado, &motivo, &evento); err != nil {
				return err
			}
			dias = append(dias, DiaVerificado{Dia: d.Time.UTC(), Estado: estado, Motivo: motivo, EventoID: evento})
		}
		if err := rows.Err(); err != nil {
			return err
		}
		rows.Close()
		q, err := db.Q(ctx)
		if err != nil {
			return err
		}
		linhas, err := q.AudListarSelos(ctx, consultasListarSelos(empresaID, de, ate))
		if err != nil {
			return fmt.Errorf("listar selos: %w", err)
		}
		for _, l := range linhas {
			s := Selo{EmpresaID: empresaID, Dia: l.Dia.Time, Quantidade: l.Quantidade, HashDia: l.HashDia,
				HashAnterior: l.HashAnterior, HashEncadeado: l.HashEncadeado, CalculadoEm: l.CalculadoEm,
				ChaveObjeto: l.ChaveObjeto}
			selos[s.Dia.UTC().Format("2006-01-02")] = s
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if err := v.conferirCopias(ctx, dias, selos); err != nil {
		return nil, err
	}
	return dias, nil
}

// conferirCopias compara, em paralelo limitado, o selo do banco com o objeto
// do armazenamento, nos dias que o banco deu como íntegros.
func (v *Verificador) conferirCopias(ctx context.Context, dias []DiaVerificado, selos map[string]Selo) error {
	var (
		wg       sync.WaitGroup
		mu       sync.Mutex
		primeiro error
		vaga     = make(chan struct{}, 8)
	)
	for i := range dias {
		d := &dias[i]
		if d.Estado != EstadoIntegro {
			continue
		}
		s, ok := selos[d.Dia.Format("2006-01-02")]
		if !ok {
			continue // o banco já teria dito selo_ausente
		}
		wg.Add(1)
		vaga <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-vaga }()
			motivo, err := v.compararCopia(ctx, s)
			mu.Lock()
			defer mu.Unlock()
			if err != nil && primeiro == nil {
				primeiro = err
			}
			if motivo != "" {
				d.Estado, d.Motivo = EstadoViolado, &motivo
			}
		}()
	}
	wg.Wait()
	if primeiro != nil {
		return erros.ServicoIndisponivel.ComChave("detalhe.selo_copia", nil).Por(primeiro)
	}
	return nil
}

func (v *Verificador) compararCopia(ctx context.Context, s Selo) (string, error) {
	corpo, err := v.arm.Ler(ctx, s.ChaveObjeto)
	if errors.Is(err, armazenamento.ErrNaoEncontrado) {
		return MotivoCopiaAusente, nil
	}
	if err != nil {
		return "", err
	}
	if !bytes.Equal(corpo, s.Conteudo()) {
		return MotivoCopiaDivergente, nil
	}
	return "", nil
}

func diaUTC(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func consultasListarSelos(empresaID uuid.UUID, de, ate time.Time) consultas.AudListarSelosParams {
	return consultas.AudListarSelosParams{EmpresaID: empresaID, De: pgtype.Date{Time: diaUTC(de), Valid: true},
		Ate: pgtype.Date{Time: diaUTC(ate), Valid: true}}
}
