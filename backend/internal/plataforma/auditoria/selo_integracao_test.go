package auditoria_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"

	"tenaz/internal/plataforma/armazenamento"
	"tenaz/internal/plataforma/auditoria"
	"tenaz/internal/plataforma/db/dbteste"
	"tenaz/internal/plataforma/erros"
)

func hoje() time.Time {
	t := time.Now().UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func dia(deslocamento int) time.Time { return hoje().AddDate(0, 0, deslocamento) }

func chaveDoSelo(empresa uuid.UUID, d time.Time) string {
	return "selos/empresas/" + empresa.String() + "/" + d.Format("2006-01-02") + ".selo"
}

// cenario é uma empresa com eventos em dias passados, ainda sem selo:
// D-5 com 3 eventos, D-4 vazio, D-3 com 2, D-2 com 1.
type cenario struct {
	empresa uuid.UUID
	eventos map[int][]uuid.UUID // deslocamento do dia -> ids
}

func TestSeloEVerificacao(t *testing.T) {
	ctx := context.Background()
	amb, err := dbteste.Preparar(ctx)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(amb.Encerrar)

	exec := func(t *testing.T, sql string, args ...any) {
		t.Helper()
		if _, err := amb.Migrador.Exec(ctx, sql, args...); err != nil {
			t.Fatalf("%s: %v", sql, err)
		}
	}
	// novoCenario cria a empresa e semeia os eventos como migrador, com
	// ocorrido_em no passado (o hash vem do trigger).
	novoCenario := func(t *testing.T, anos int) cenario {
		t.Helper()
		c := cenario{empresa: uuid.Must(uuid.NewV7()), eventos: map[int][]uuid.UUID{}}
		exec(t, `INSERT INTO app.empresa (id, razao_social, retencao_auditoria_anos) VALUES ($1, 'Empresa do selo', $2)`, c.empresa, anos)
		for desloc, n := range map[int]int{-5: 3, -4: 0, -3: 2, -2: 1} {
			for i := 0; i < n; i++ {
				id := uuid.Must(uuid.NewV7())
				exec(t, `INSERT INTO app.evento_auditoria (empresa_id, id, ocorrido_em, ator_tipo, acao, entidade, entidade_id, modulo, request_id)
VALUES ($1, $2, $3, 'sistema', 'inclusao', 'unidade', $2, 'adm', $4)`,
					c.empresa, id, dia(desloc).Add(time.Duration(10+i)*time.Hour), fmt.Sprintf("semente-%d-%d", desloc, i))
				c.eventos[desloc] = append(c.eventos[desloc], id)
			}
		}
		return c
	}
	selador := func(t *testing.T, arm armazenamento.Armazenamento) *auditoria.Selador {
		t.Helper()
		s, err := auditoria.NovoSelador(auditoria.DependenciasSelador{Banco: amb.Job, Armazenamento: arm})
		if err != nil {
			t.Fatal(err)
		}
		return s
	}
	verificador := func(t *testing.T, arm armazenamento.Armazenamento) *auditoria.Verificador {
		t.Helper()
		v, err := auditoria.NovoVerificador(amb.App, arm)
		if err != nil {
			t.Fatal(err)
		}
		return v
	}
	// adulterar roda o comando como migrador com o trigger de somente inserção
	// desligado e religado na mesma transação: ninguém mais enxerga a janela.
	adulterar := func(t *testing.T, tabela string, comandos ...string) {
		t.Helper()
		tx, err := amb.Migrador.Begin(ctx)
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = tx.Rollback(ctx) }()
		if _, err := tx.Exec(ctx, fmt.Sprintf(`ALTER TABLE %s DISABLE TRIGGER trg_somente_insercao`, tabela)); err != nil {
			t.Fatal(err)
		}
		for _, c := range comandos {
			if _, err := tx.Exec(ctx, c); err != nil {
				t.Fatalf("%s: %v", c, err)
			}
		}
		if _, err := tx.Exec(ctx, fmt.Sprintf(`ALTER TABLE %s ENABLE TRIGGER trg_somente_insercao`, tabela)); err != nil {
			t.Fatal(err)
		}
		if err := tx.Commit(ctx); err != nil {
			t.Fatal(err)
		}
	}
	periodo := func(t *testing.T, v *auditoria.Verificador, empresa uuid.UUID, de, ate int) map[int]auditoria.DiaVerificado {
		t.Helper()
		dias, err := v.Verificar(ctx, empresa, dia(de), dia(ate))
		if err != nil {
			t.Fatal(err)
		}
		if len(dias) != ate-de+1 {
			t.Fatalf("esperados %d dias, vieram %d", ate-de+1, len(dias))
		}
		m := map[int]auditoria.DiaVerificado{}
		for i, d := range dias {
			m[de+i] = d
		}
		return m
	}
	esperar := func(t *testing.T, d auditoria.DiaVerificado, estado, motivo string) {
		t.Helper()
		m := ""
		if d.Motivo != nil {
			m = *d.Motivo
		}
		if d.Estado != estado || m != motivo {
			t.Errorf("dia %s: %s/%s, esperado %s/%s", d.Dia.Format("2006-01-02"), d.Estado, m, estado, motivo)
		}
	}

	t.Run("o job sela, copia ao armazenamento com a retenção da empresa e é idempotente", func(t *testing.T) {
		c := novoCenario(t, 7)
		arm := dbteste.NovaMemoria()
		sel := selador(t, arm)
		antes := time.Now().UTC()
		selados, enviados, err := sel.Empresa(ctx, c.empresa)
		depois := time.Now().UTC()
		if err != nil {
			t.Fatal(err)
		}
		// D-5 a D-2 sempre; D-1 (vazio) só depois de 01:00 UTC
		if selados < 4 || selados > 5 || enviados != selados {
			t.Fatalf("selados %d, enviados %d", selados, enviados)
		}
		for desloc, ids := range map[int]int{-5: 3, -4: 0, -3: 2, -2: 1} {
			var quantidade int64
			var chave, hashDia, hashEnc string
			var anterior *string
			var calculado time.Time
			if err := amb.Migrador.QueryRow(ctx, `SELECT quantidade, chave_objeto, encode(hash_dia, 'hex'), encode(hash_encadeado, 'hex'),
  encode(hash_anterior, 'hex'), calculado_em FROM app.selo_auditoria WHERE empresa_id = $1 AND dia = $2`, c.empresa, dia(desloc)).
				Scan(&quantidade, &chave, &hashDia, &hashEnc, &anterior, &calculado); err != nil {
				t.Fatalf("selo do dia D%d: %v", desloc, err)
			}
			if int(quantidade) != ids {
				t.Errorf("D%d: quantidade %d, esperada %d (o dia vazio também é selado)", desloc, quantidade, ids)
			}
			if chave != chaveDoSelo(c.empresa, dia(desloc)) {
				t.Errorf("D%d: chave %s", desloc, chave)
			}
			corpo, err := arm.Ler(ctx, chave)
			if err != nil {
				t.Fatalf("D%d: objeto ausente: %v", desloc, err)
			}
			ant := "null"
			if anterior != nil {
				ant = `"` + *anterior + `"`
			}
			// chaves em ordem de bytes, sem espaço, hash em hexadecimal, formato 1
			esperado := fmt.Sprintf(`{"calculado_em":"%s","dia":"%s","empresa_id":"%s","formato":1,"hash_anterior":%s,"hash_dia":"%s","hash_encadeado":"%s","quantidade":%d}`,
				calculado.UTC().Format("2006-01-02T15:04:05.000000Z"), dia(desloc).Format("2006-01-02"), c.empresa, ant, hashDia, hashEnc, quantidade)
			if string(corpo) != esperado {
				t.Errorf("D%d: conteúdo do selo\n  veio     %s\n  esperado %s", desloc, corpo, esperado)
			}
			ate := arm.Retencao(chave)
			if ate.Before(antes.AddDate(7, 0, 0)) || ate.After(depois.AddDate(7, 0, 0)) {
				t.Errorf("D%d: retenção até %v, esperada agora + 7 anos (retencao_auditoria_anos)", desloc, ate)
			}
		}
		// o primeiro selo não tem anterior; os outros encadeiam
		var primeiroTemAnterior bool
		if err := amb.Migrador.QueryRow(ctx, `SELECT hash_anterior IS NOT NULL FROM app.selo_auditoria WHERE empresa_id = $1 AND dia = $2`,
			c.empresa, dia(-5)).Scan(&primeiroTemAnterior); err != nil || primeiroTemAnterior {
			t.Errorf("o primeiro selo não tem hash_anterior: %v %v", primeiroTemAnterior, err)
		}

		// segunda rodada: nada a fazer
		s2, e2, err := sel.Empresa(ctx, c.empresa)
		if err != nil || s2 != 0 || e2 != 0 {
			t.Errorf("segunda rodada: selados %d, enviados %d, %v", s2, e2, err)
		}
		// copia que sumiu dos últimos 7 dias é reenviada
		arm.Apagar(chaveDoSelo(c.empresa, dia(-3)))
		s3, e3, err := sel.Empresa(ctx, c.empresa)
		if err != nil || s3 != 0 || e3 != 1 {
			t.Errorf("reenvio: selados %d, enviados %d, %v", s3, e3, err)
		}
		if ok, _ := arm.Existe(ctx, chaveDoSelo(c.empresa, dia(-3))); !ok {
			t.Error("o selo de D-3 deveria ter voltado")
		}
	})

	t.Run("falha do armazenamento não deixa dia sem cópia: a rodada seguinte reenvia", func(t *testing.T) {
		c := novoCenario(t, 5)
		arm := dbteste.NovaMemoria()
		arm.FalhaGravar = true
		sel := selador(t, arm)
		selados, enviados, err := sel.Empresa(ctx, c.empresa)
		if err == nil || selados < 4 || enviados != 0 {
			t.Fatalf("a falha de gravação tem de voltar: selados %d, enviados %d, %v", selados, enviados, err)
		}
		arm.FalhaGravar = false
		s2, e2, err := sel.Empresa(ctx, c.empresa)
		if err != nil || s2 != 0 || e2 != selados {
			t.Errorf("reenvio: selados %d, enviados %d (esperados 0 e %d), %v", s2, e2, selados, err)
		}
	})

	t.Run("Executar percorre as empresas e uma falha não interrompe as outras", func(t *testing.T) {
		c1, c2 := novoCenario(t, 5), novoCenario(t, 5)
		arm := dbteste.NovaMemoria()
		res, err := selador(t, arm).Executar(ctx)
		if err != nil {
			t.Fatal(err)
		}
		if res.Empresas < 2 || res.Selados < 8 {
			t.Errorf("resultado: %+v", res)
		}
		for _, e := range []uuid.UUID{c1.empresa, c2.empresa} {
			if ok, _ := arm.Existe(ctx, chaveDoSelo(e, dia(-5))); !ok {
				t.Errorf("empresa %s sem cópia do selo de D-5", e)
			}
		}
	})

	t.Run("verificação: íntegro, não selado e isolamento entre empresas", func(t *testing.T) {
		a, b := novoCenario(t, 5), novoCenario(t, 5)
		arm := dbteste.NovaMemoria()
		sel := selador(t, arm)
		for _, c := range []cenario{a, b} {
			if _, _, err := sel.Empresa(ctx, c.empresa); err != nil {
				t.Fatal(err)
			}
		}
		v := verificador(t, arm)
		r := periodo(t, v, a.empresa, -7, -2)
		esperar(t, r[-7], auditoria.EstadoNaoSelado, "anterior_ao_primeiro_selo")
		esperar(t, r[-6], auditoria.EstadoNaoSelado, "anterior_ao_primeiro_selo")
		for _, d := range []int{-5, -4, -3, -2} {
			esperar(t, r[d], auditoria.EstadoIntegro, "")
		}
		// o dia de hoje ainda não tem selo
		h := periodo(t, v, a.empresa, 0, 0)
		esperar(t, h[0], auditoria.EstadoNaoSelado, "ainda_nao_selado")

		// Adulterar a empresa A não aparece na verificação de B, e os selos de
		// B são os dela: o isolamento A/B da verificação.
		alvo := a.eventos[-3][0]
		adulterar(t, "app.evento_auditoria", fmt.Sprintf(`UPDATE app.evento_auditoria SET entidade = 'adulterada' WHERE id = '%s'`, alvo))
		ra := periodo(t, v, a.empresa, -5, -2)
		esperar(t, ra[-3], auditoria.EstadoViolado, "evento_alterado")
		if ra[-3].EventoID == nil || *ra[-3].EventoID != alvo {
			t.Errorf("a violação tem de apontar o evento adulterado %s: %v", alvo, ra[-3].EventoID)
		}
		esperar(t, ra[-5], auditoria.EstadoIntegro, "")
		esperar(t, ra[-2], auditoria.EstadoIntegro, "")
		rb := periodo(t, v, b.empresa, -5, -2)
		for _, d := range []int{-5, -4, -3, -2} {
			esperar(t, rb[d], auditoria.EstadoIntegro, "")
		}
		// empresa sem nenhum selo: nada de outra empresa aparece
		vazia := uuid.Must(uuid.NewV7())
		exec(t, `INSERT INTO app.empresa (id, razao_social) VALUES ($1, 'Sem eventos')`, vazia)
		rv := periodo(t, v, vazia, -5, -2)
		for _, d := range []int{-5, -4, -3, -2} {
			esperar(t, rv[d], auditoria.EstadoNaoSelado, "ainda_nao_selado")
		}
	})

	t.Run("violação: evento apagado, evento incluído depois, selo adulterado, cadeia quebrada", func(t *testing.T) {
		c := novoCenario(t, 5)
		arm := dbteste.NovaMemoria()
		if _, _, err := selador(t, arm).Empresa(ctx, c.empresa); err != nil {
			t.Fatal(err)
		}
		v := verificador(t, arm)

		// evento incluído num dia vazio já selado
		exec(t, `INSERT INTO app.evento_auditoria (empresa_id, ocorrido_em, ator_tipo, acao, entidade, modulo)
VALUES ($1, $2, 'sistema', 'inclusao', 'unidade', 'adm')`, c.empresa, dia(-4).Add(5*time.Hour))
		// evento incluído antes do primeiro selo
		exec(t, `INSERT INTO app.evento_auditoria (empresa_id, ocorrido_em, ator_tipo, acao, entidade, modulo)
VALUES ($1, $2, 'sistema', 'inclusao', 'unidade', 'adm')`, c.empresa, dia(-6).Add(5*time.Hour))
		// evento apagado
		adulterar(t, "app.evento_auditoria", fmt.Sprintf(`DELETE FROM app.evento_auditoria WHERE id = '%s'`, c.eventos[-5][0]))
		// selo com o hash do encadeamento trocado (hash_dia correto)
		adulterar(t, "app.selo_auditoria", fmt.Sprintf(`UPDATE app.selo_auditoria SET hash_encadeado = sha256('forjado'::bytea)
WHERE empresa_id = '%s' AND dia = '%s'`, c.empresa, dia(-3).Format("2006-01-02")))
		// selo do dia seguinte apontando para outro anterior não é caso aqui: o de D-2 fica com o anterior verdadeiro,
		// que já não encaixa no D-3 adulterado
		r := periodo(t, v, c.empresa, -6, -2)
		esperar(t, r[-6], auditoria.EstadoViolado, "evento_anterior_ao_primeiro_selo")
		esperar(t, r[-5], auditoria.EstadoViolado, "eventos_divergentes")
		esperar(t, r[-4], auditoria.EstadoViolado, "eventos_divergentes")
		esperar(t, r[-3], auditoria.EstadoViolado, "selo_adulterado")
		esperar(t, r[-2], auditoria.EstadoViolado, "cadeia_quebrada")

		// selo apagado no meio da cadeia
		d := novoCenario(t, 5)
		if _, _, err := selador(t, arm).Empresa(ctx, d.empresa); err != nil {
			t.Fatal(err)
		}
		adulterar(t, "app.selo_auditoria", fmt.Sprintf(`DELETE FROM app.selo_auditoria WHERE empresa_id = '%s' AND dia = '%s'`,
			d.empresa, dia(-4).Format("2006-01-02")))
		rd := periodo(t, v, d.empresa, -5, -2)
		esperar(t, rd[-4], auditoria.EstadoViolado, "selo_ausente")
	})

	t.Run("quem refaz eventos e selos coerentes no banco é pego pela cópia imutável", func(t *testing.T) {
		c := novoCenario(t, 5)
		arm := dbteste.NovaMemoria()
		if _, _, err := selador(t, arm).Empresa(ctx, c.empresa); err != nil {
			t.Fatal(err)
		}
		v := verificador(t, arm)
		alvo := c.eventos[-2][0]
		// Refaz o último dia de ponta a ponta: conteúdo, hash do evento, hash do dia e encadeado.
		adulterar(t, "app.evento_auditoria", fmt.Sprintf(`UPDATE app.evento_auditoria SET entidade = 'forjada' WHERE id = '%s'`, alvo),
			fmt.Sprintf(`UPDATE app.evento_auditoria e SET hash = app.hash_evento_auditoria(e) WHERE e.id = '%s'`, alvo))
		adulterar(t, "app.selo_auditoria", fmt.Sprintf(`UPDATE app.selo_auditoria s SET
  hash_dia = app.hash_dia_auditoria(s.empresa_id, s.dia, s.quantidade, (SELECT string_agg(e.hash, ''::bytea ORDER BY e.ocorrido_em, e.id)
    FROM app.evento_auditoria e WHERE e.empresa_id = s.empresa_id AND e.ocorrido_em >= s.dia AND e.ocorrido_em < s.dia + 1))
WHERE s.empresa_id = '%s' AND s.dia = '%s'`, c.empresa, dia(-2).Format("2006-01-02")),
			fmt.Sprintf(`UPDATE app.selo_auditoria SET hash_encadeado = sha256(coalesce(hash_anterior, ''::bytea) || hash_dia)
WHERE empresa_id = '%s' AND dia = '%s'`, c.empresa, dia(-2).Format("2006-01-02")))
		// o banco, sozinho, já não vê nada de errado
		tx, err := amb.Migrador.Begin(ctx)
		if err != nil {
			t.Fatal(err)
		}
		var estado string
		if err := tx.QueryRow(ctx, `SELECT set_config('app.empresa_id', $1, true)`, c.empresa.String()).Scan(&estado); err != nil {
			t.Fatal(err)
		}
		if err := tx.QueryRow(ctx, `SELECT estado FROM app.verificar_integridade_auditoria($1::date, $1::date)`, dia(-2)).Scan(&estado); err != nil {
			t.Fatal(err)
		}
		_ = tx.Rollback(ctx)
		if estado != auditoria.EstadoIntegro {
			t.Fatalf("a falsificação coerente deveria enganar o banco (estado %s)", estado)
		}
		r := periodo(t, v, c.empresa, -5, -2)
		esperar(t, r[-2], auditoria.EstadoViolado, auditoria.MotivoCopiaDivergente)
		esperar(t, r[-3], auditoria.EstadoIntegro, "")

		// cópia ausente
		arm.Apagar(chaveDoSelo(c.empresa, dia(-3)))
		r = periodo(t, v, c.empresa, -5, -2)
		esperar(t, r[-3], auditoria.EstadoViolado, auditoria.MotivoCopiaAusente)
		esperar(t, r[-5], auditoria.EstadoIntegro, "")
	})

	t.Run("período inválido e armazenamento ausente", func(t *testing.T) {
		v := verificador(t, dbteste.NovaMemoria())
		_, err := v.Verificar(ctx, uuid.Must(uuid.NewV7()), dia(-1), dia(-3))
		if !errors.Is(err, auditoria.ErroPeriodo) {
			t.Errorf("fim antes do início: %v", err)
		}
		_, err = v.Verificar(ctx, uuid.Must(uuid.NewV7()), dia(-500), dia(-1))
		if !errors.Is(err, auditoria.ErroPeriodo) {
			t.Errorf("mais de 400 dias: %v", err)
		}
		_, err = verificador(t, nil).Verificar(ctx, uuid.Must(uuid.NewV7()), dia(-3), dia(-1))
		if !errors.Is(err, erros.ServicoIndisponivel) {
			t.Errorf("sem armazenamento: %v", err)
		}
	})

	t.Run("o hash do banco é o do formato 1 recalculado no Go", func(t *testing.T) {
		c := novoCenario(t, 5)
		// valores que exigem escape, nulos e jsonb aninhado
		exec(t, `INSERT INTO app.evento_auditoria (empresa_id, ocorrido_em, ator_tipo, usuario_nome, acao, resultado, sensivel,
  entidade, entidade_chave, antes, depois, detalhe, campos_alterados, modulo, tela, origem, ip, user_agent, request_id, escopo)
VALUES ($1, now(), 'usuario', 'Zé "da Silva" \ '||chr(10)||chr(9)||chr(1)||' 日本 <&>', 'visualizacao', 'sucesso', true,
  'pessoa', 'chave "x"', NULL, '{"b": 1, "aa": [true, null, "x"], "c": {"z": 1.50, "y": "ç"}}'::jsonb,
  '{"filtros": {"de": "2026-01-01"}, "quantidade": 10}'::jsonb, ARRAY['nome', 'cpf'], 'adm', 'adm.pessoas', 'web',
  '2001:db8::1', 'Mozilla/5.0 "x" '||chr(8232), 'req-1', 'interno')`, c.empresa)
		n, div, err := dbteste.ConferirHashesDaTrilha(ctx, amb.Migrador, c.empresa)
		if err != nil {
			t.Fatal(err)
		}
		if n < 7 {
			t.Fatalf("só %d eventos conferidos", n)
		}
		for _, d := range div {
			t.Error(d)
		}
	})
}
