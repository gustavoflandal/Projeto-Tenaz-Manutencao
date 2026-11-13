package auditoria

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	"tenaz/internal/plataforma/armazenamento"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/obs"
)

// FormatoSelo versiona o conteúdo do objeto do selo (documento 2, seção 9.4).
const FormatoSelo = 1

// diasDeConferencia é a janela em que o job confere por HEAD que a cópia
// existe e reenvia a que faltar: falha de rede não deixa dia sem cópia.
const diasDeConferencia = 7

// Selo é uma linha de selo_auditoria.
type Selo struct {
	EmpresaID     uuid.UUID
	Dia           time.Time // dia UTC, à meia-noite
	Quantidade    int64
	HashDia       []byte
	HashAnterior  []byte // nulo só no primeiro selo
	HashEncadeado []byte
	CalculadoEm   time.Time
	ChaveObjeto   string
}

const formatoInstante = "2006-01-02T15:04:05.000000Z"

// Conteudo é o JSON do objeto do selo (documento 2, seção 9.4): sem espaços e
// com as chaves em ordem de bytes (calculado_em, dia, empresa_id, formato,
// hash_anterior, hash_dia, hash_encadeado, quantidade). Os hashes vão em
// hexadecimal minúsculo; hash_anterior é null no primeiro selo. Os valores são
// só uuid, data, hexadecimal e número, sem caractere que peça escape.
func (s Selo) Conteudo() []byte {
	var b strings.Builder
	b.WriteString(`{"calculado_em":"` + s.CalculadoEm.UTC().Format(formatoInstante) + `"`)
	b.WriteString(`,"dia":"` + s.Dia.UTC().Format("2006-01-02") + `"`)
	b.WriteString(`,"empresa_id":"` + s.EmpresaID.String() + `"`)
	b.WriteString(`,"formato":` + strconv.Itoa(FormatoSelo))
	if s.HashAnterior == nil {
		b.WriteString(`,"hash_anterior":null`)
	} else {
		b.WriteString(`,"hash_anterior":"` + hex.EncodeToString(s.HashAnterior) + `"`)
	}
	b.WriteString(`,"hash_dia":"` + hex.EncodeToString(s.HashDia) + `"`)
	b.WriteString(`,"hash_encadeado":"` + hex.EncodeToString(s.HashEncadeado) + `"`)
	b.WriteString(`,"quantidade":` + strconv.FormatInt(s.Quantidade, 10) + `}`)
	return []byte(b.String())
}

// DependenciasSelador do job do selo.
type DependenciasSelador struct {
	// Banco é o pool do papel tenaz_job: o job só lê a trilha da empresa,
	// grava selo_auditoria (pelo trg_selo_calcular) e executa
	// app.selar_dias_pendentes().
	Banco         *db.DB
	Armazenamento armazenamento.Armazenamento
	Agora         func() time.Time
}

// Selador é o job diário do selo da trilha (RF-ADM-011; documento 2, seção
// 9.4). Não há agendador de jobs no projeto (ADR-005, River, ainda proposta):
// o serviço é chamado pelo subcomando `tenaz selar-auditoria`, que o
// agendamento da infraestrutura dispara uma vez por dia, depois de 01:00 UTC.
type Selador struct {
	d DependenciasSelador
}

// NovoSelador monta o job.
func NovoSelador(d DependenciasSelador) (*Selador, error) {
	if d.Banco == nil || d.Armazenamento == nil {
		return nil, errors.New("selador: dependência ausente")
	}
	if d.Agora == nil {
		d.Agora = time.Now
	}
	return &Selador{d: d}, nil
}

// ResultadoSelo resume uma rodada.
type ResultadoSelo struct {
	Empresas int // empresas percorridas
	Selados  int // dias selados no banco nesta rodada
	Enviados int // objetos gravados no armazenamento imutável
}

// Executar percorre as empresas, uma transação por empresa, e para cada uma:
// sela os dias pendentes no banco (o banco calcula tudo), grava no
// armazenamento imutável o selo de cada dia novo e dos últimos 7 dias que
// estejam sem cópia. Falha numa empresa não interrompe as outras; os erros
// voltam juntos no fim. Idempotente.
func (j *Selador) Executar(ctx context.Context) (ResultadoSelo, error) {
	var res ResultadoSelo
	empresas, err := j.empresas(ctx)
	if err != nil {
		return res, err
	}
	var falhas []error
	for _, id := range empresas {
		res.Empresas++
		selados, enviados, err := j.Empresa(ctx, id)
		res.Selados += selados
		res.Enviados += enviados
		if err != nil {
			falhas = append(falhas, fmt.Errorf("empresa %s: %w", id, err))
		}
	}
	if err := errors.Join(falhas...); err != nil {
		obs.Logger(ctx).Error("falha ao selar a trilha de auditoria", "erro", err)
		return res, err
	}
	obs.Logger(ctx).Info("trilha de auditoria selada", "empresas", res.Empresas, "dias_selados", res.Selados, "objetos_enviados", res.Enviados)
	return res, nil
}

// empresas lista as empresas FORA de qualquer transação de empresa: é a única
// leitura sem contexto, e a policy empresa_listagem_job só a permite quando
// app.empresa_atual() é nulo (documento 2, 00004).
func (j *Selador) empresas(ctx context.Context) ([]uuid.UUID, error) {
	rows, err := j.d.Banco.Pool().Query(ctx, `SELECT id FROM app.empresa ORDER BY id`)
	if err != nil {
		return nil, fmt.Errorf("listar empresas para selar: %w", err)
	}
	defer rows.Close()
	var ids []uuid.UUID
	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, rows.Err()
}

// Empresa sela e copia os selos de uma empresa só: o corpo de Executar, que
// também serve para refazer uma empresa que falhou. Devolve os dias selados no
// banco e os objetos gravados no armazenamento imutável.
func (j *Selador) Empresa(ctx context.Context, id uuid.UUID) (selados, enviados int, err error) {
	ctx = db.ComRequisicao(ctx, db.Requisicao{Origem: db.OrigemJob, Modulo: "adm", RequestID: uuid.NewString()})
	agora := j.d.Agora().UTC()
	var (
		anos  int32
		selos []Selo
	)
	// 1. sela no banco e lê o que precisa ser conferido, numa transação só
	err = j.d.Banco.NaEmpresa(ctx, id, db.Interno(nil), func(ctx context.Context) error {
		q, err := db.Q(ctx)
		if err != nil {
			return err
		}
		if anos, err = q.AudRetencaoEmpresa(ctx, id); err != nil {
			return fmt.Errorf("ler retenção da empresa: %w", err)
		}
		novos, err := q.AudSelarDiasPendentes(ctx)
		if err != nil {
			return fmt.Errorf("selar dias pendentes: %w", err)
		}
		selados = len(novos)
		de := agora.Truncate(24*time.Hour).AddDate(0, 0, -(diasDeConferencia + 1))
		for _, d := range novos {
			if d.Valid && d.Time.Before(de) {
				de = d.Time
			}
		}
		linhas, err := q.AudListarSelos(ctx, consultasListarSelos(id, de, agora))
		if err != nil {
			return fmt.Errorf("listar selos: %w", err)
		}
		for _, l := range linhas {
			selos = append(selos, Selo{
				EmpresaID: id, Dia: l.Dia.Time, Quantidade: l.Quantidade, HashDia: l.HashDia,
				HashAnterior: l.HashAnterior, HashEncadeado: l.HashEncadeado, CalculadoEm: l.CalculadoEm,
				ChaveObjeto: l.ChaveObjeto,
			})
		}
		return nil
	})
	if err != nil {
		return 0, 0, err
	}

	// 2. fora da transação: o que falta no armazenamento imutável
	ate := agora.AddDate(int(anos), 0, 0)
	var falhas []error
	for _, s := range selos {
		existe, err := j.d.Armazenamento.Existe(ctx, s.ChaveObjeto)
		if err != nil {
			falhas = append(falhas, err)
			continue
		}
		if existe {
			continue
		}
		if err := j.d.Armazenamento.GravarImutavel(ctx, s.ChaveObjeto, s.Conteudo(), ate); err != nil {
			falhas = append(falhas, err)
			continue
		}
		enviados++
	}
	return selados, enviados, errors.Join(falhas...)
}
