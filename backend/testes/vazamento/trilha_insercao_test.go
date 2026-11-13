package vazamento

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
)

// Restritivas de INSERT em app.evento_auditoria (seções 9.2 e 2.10.4 do
// documento 2; ADR-011). A permissiva auditoria_insercao só confere a
// empresa; quem limita O QUE cada papel grava são estas restritivas:
//
//	auditoria_sem_autenticacao  tenaz_app, tenaz_job, tenaz_sinc não gravam
//	                            evento da identidade (ISO-E01-07)
//	escopo_externo_insercao     escopo externo só grava evento em nome de si
//	conexao_sinc                o sincronizador só grava evento da conexão em
//	                            processamento, como sincronizacao_conexao
//
// Cada uma tem o seu molde (o que passa e o que é recusado PELA POLICY, com
// "row-level security" na mensagem) e a sua prova: sem a policy, o verificador
// acusa.

type tentativaTrilha struct {
	papel  string
	c      contexto
	evento tentativa
	aceita bool
}

func eventoTrilha(empresa, ator, usuario, acao, entidade, conexao string) tentativa {
	return eventoTrilhaNoEscopo(empresa, ator, usuario, acao, entidade, conexao, "")
}

// eventoTrilhaNoEscopo grava também a coluna escopo do evento, que a
// escopo_externo_insercao confere contra o escopo da transação.
func eventoTrilhaNoEscopo(empresa, ator, usuario, acao, entidade, conexao, escopo string) tentativa {
	return tentativa{
		fmt.Sprintf("empresa=%.8s ator=%s usuario=%.8s acao=%s entidade=%s conexao=%.8s escopo=%s",
			empresa, ator, usuario, acao, entidade, conexao, escopo),
		`INSERT INTO app.evento_auditoria (empresa_id, ator_tipo, usuario_id, acao, entidade, conexao_empresa_id, escopo)
		 VALUES ($1, $2, nullif($3, '')::uuid, $4, $5, nullif($6, '')::uuid, nullif($7, ''))`,
		[]any{empresa, ator, usuario, acao, entidade, conexao, escopo},
	}
}

// problemasDaTrilha executa as tentativas e devolve as divergências do molde.
func problemasDaTrilha(t *testing.T, ctx context.Context, tx pgx.Tx, tentativas []tentativaTrilha) []string {
	t.Helper()
	var problemas []string
	for _, x := range tentativas {
		err := tentarComo(ctx, tx, x.papel, x.c, x.evento.sql, x.evento.args...)
		switch {
		case x.aceita && err != nil:
			problemas = append(problemas, fmt.Sprintf("%s: evento legítimo recusado: %s: %v", x.papel, x.evento.rotulo, err))
		case !x.aceita && !recusaDeRLS(err):
			problemas = append(problemas, fmt.Sprintf("%s: não recusado pela policy: %s (resultado: %v)",
				x.papel, x.evento.rotulo, err))
		}
	}
	return problemas
}

// ------------------------------------------- auditoria_sem_autenticacao

func moldeSemAutenticacao(f *fixturesE01) []tentativaTrilha {
	a := f.A
	conexao := novoID()
	sinc := contexto{Empresa: a.ID, Conexao: conexao}
	return []tentativaTrilha{
		{papelApp, f.internoA(), eventoTrilha(a.ID, "usuario", a.UsuarioInterno, "alteracao", "unidade", ""), true},
		{papelApp, f.internoA(), eventoTrilha(a.ID, "usuario", a.UsuarioInterno, "login", "identidade", ""), false},
		{papelApp, f.internoA(), eventoTrilha(a.ID, "sistema", "", "alteracao", "identidade", ""), false},
		{papelJob, f.internoA(), eventoTrilha(a.ID, "sistema", "", "logout", "identidade", ""), false},
		{papelSinc, sinc, eventoTrilha(a.ID, "sincronizacao_conexao", "", "login", "identidade", conexao), false},
	}
}

func TestTrilhaSoAutenticacaoGravaEventoDeIdentidade(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		usagemDoSincronizadorEmApp(t, ctx, tx)
		for _, p := range problemasDaTrilha(t, ctx, tx, moldeSemAutenticacao(f)) {
			t.Errorf("%s (auditoria_sem_autenticacao, ISO-E01-07)", p)
		}
	})
}

func TestDetectorAcusaTrilhaSemAuditoriaSemAutenticacao(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		usagemDoSincronizadorEmApp(t, ctx, tx)
		estragar(t, ctx, tx, "DROP POLICY auditoria_sem_autenticacao ON app.evento_auditoria")
		p := strings.Join(problemasDaTrilha(t, ctx, tx, moldeSemAutenticacao(f)), " | ")
		if !strings.Contains(p, "tenaz_app: não recusado pela policy: empresa=") || !strings.Contains(p, "acao=login entidade=identidade") {
			t.Errorf("sem a policy, o evento de login da identidade pelo tenaz_app não foi acusado: %s", p)
		}
	})
}

// --------------------------------------------- escopo_externo_insercao

func moldeEscopoInsercao(f *fixturesE01) []tentativaTrilha {
	a := f.A
	return []tentativaTrilha{
		// o externo grava em nome de si, no próprio escopo
		{papelApp, f.prestadorF1(), eventoTrilhaNoEscopo(a.ID, "usuario", f.UsuarioF1, "alteracao", "usuario", "", "prestador"), true},
		{papelApp, f.locatarioL1(), eventoTrilhaNoEscopo(a.ID, "usuario", f.UsuarioL1, "alteracao", "usuario", "", "locatario"), true},
		{papelApp, f.linkA(), eventoTrilhaNoEscopo(a.ID, "link_execucao", "", "alteracao", "ordem_servico", "", "link"), true},
		// o interno grava evento de qualquer usuário da empresa (ex.: o serviço)
		{papelApp, f.internoA(), eventoTrilhaNoEscopo(a.ID, "sistema", f.UsuarioF1, "alteracao", "usuario", "", "interno"), true},
		// o externo não atribui o evento a outra pessoa, nem o deixa sem autor
		{papelApp, f.prestadorF1(), eventoTrilhaNoEscopo(a.ID, "usuario", a.UsuarioInterno, "exclusao", "usuario", "", "prestador"), false},
		{papelApp, f.prestadorF1(), eventoTrilhaNoEscopo(a.ID, "usuario", f.UsuarioF1Colega, "alteracao", "usuario", "", "prestador"), false},
		{papelApp, f.prestadorF1(), eventoTrilhaNoEscopo(a.ID, "sistema", "", "alteracao", "usuario", "", "prestador"), false},
		{papelApp, f.locatarioL1(), eventoTrilhaNoEscopo(a.ID, "usuario", f.UsuarioL2, "alteracao", "usuario", "", "locatario"), false},
		{papelApp, f.linkA(), eventoTrilhaNoEscopo(a.ID, "usuario", a.UsuarioInterno, "alteracao", "ordem_servico", "", "link"), false},
		// nem finge outro escopo: o prestador gravando o evento como se fosse interno
		{papelApp, f.prestadorF1(), eventoTrilhaNoEscopo(a.ID, "usuario", f.UsuarioF1, "alteracao", "usuario", "", "interno"), false},
	}
}

func TestTrilhaEscopoExternoSoGravaEmNomeDeSi(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		for _, p := range problemasDaTrilha(t, ctx, tx, moldeEscopoInsercao(f)) {
			t.Errorf("%s (escopo_externo_insercao)", p)
		}
	})
}

func TestDetectorAcusaTrilhaSemEscopoExternoInsercao(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		estragar(t, ctx, tx, "DROP POLICY escopo_externo_insercao ON app.evento_auditoria")
		p := strings.Join(problemasDaTrilha(t, ctx, tx, moldeEscopoInsercao(f)), " | ")
		if !strings.Contains(p, fmt.Sprintf("usuario=%.8s acao=exclusao", f.A.UsuarioInterno)) {
			t.Errorf("sem a policy, o prestador gravando evento em nome do interno não foi acusado: %s", p)
		}
	})
}

// --------------------------------------------------------- conexao_sinc

// usagemDoSincronizadorEmApp concede, só na transação desfeita, o USAGE em
// app que o tenaz_sinc ainda não tem (ver a revisão de isolamento: pendência
// para E04). Sem ele, toda tentativa do sincronizador seria recusada por falta
// de USAGE no esquema — e o teste da policy passaria pela razão errada.
func usagemDoSincronizadorEmApp(t *testing.T, ctx context.Context, tx pgx.Tx) {
	t.Helper()
	estragar(t, ctx, tx, "GRANT USAGE ON SCHEMA app TO tenaz_sinc")
}

func moldeConexaoSinc(f *fixturesE01) []tentativaTrilha {
	a, b := f.A, f.B
	conexao, outra := novoID(), novoID()
	c := contexto{Empresa: a.ID, Conexao: conexao}
	return []tentativaTrilha{
		{papelSinc, c, eventoTrilha(a.ID, "sincronizacao_conexao", "", "alteracao", "ordem_servico", conexao), true},
		{"tenaz_sinc_conexao", c, eventoTrilha(a.ID, "sincronizacao_conexao", "", "alteracao", "ordem_servico", conexao), true},
		// evento de outra conexão
		{papelSinc, c, eventoTrilha(a.ID, "sincronizacao_conexao", "", "alteracao", "ordem_servico", outra), false},
		// evento sem conexão, ou com outro ator
		{papelSinc, c, eventoTrilha(a.ID, "sincronizacao_conexao", "", "alteracao", "ordem_servico", ""), false},
		{papelSinc, c, eventoTrilha(a.ID, "usuario", a.UsuarioInterno, "alteracao", "ordem_servico", conexao), false},
		// sem a conexão no contexto
		{papelSinc, contexto{Empresa: a.ID}, eventoTrilha(a.ID, "sincronizacao_conexao", "", "alteracao", "ordem_servico", conexao), false},
		// evento de B na transação de A
		{papelSinc, c, eventoTrilha(b.ID, "sincronizacao_conexao", "", "alteracao", "ordem_servico", conexao), false},
	}
}

func TestTrilhaSincronizadorSoGravaAConexaoEmProcessamento(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		usagemDoSincronizadorEmApp(t, ctx, tx)
		for _, p := range problemasDaTrilha(t, ctx, tx, moldeConexaoSinc(f)) {
			t.Errorf("%s (conexao_sinc, R5)", p)
		}
		// O sincronizador grava, mas não lê a trilha.
		err := tentarComo(ctx, tx, papelSinc, contexto{Empresa: f.A.ID}, `SELECT count(*) FROM app.evento_auditoria`)
		exigirErro(t, "tenaz_sinc lê a trilha", err, erroPrivilegio)
	})
}

func TestDetectorAcusaTrilhaSemConexaoSinc(t *testing.T) {
	ctx := context.Background()
	f := fixtures(t)
	emTransacaoDesfeita(t, ctx, func(tx pgx.Tx) {
		usagemDoSincronizadorEmApp(t, ctx, tx)
		estragar(t, ctx, tx, "DROP POLICY conexao_sinc ON app.evento_auditoria")
		p := strings.Join(problemasDaTrilha(t, ctx, tx, moldeConexaoSinc(f)), " | ")
		if !strings.Contains(p, "tenaz_sinc: não recusado pela policy: empresa=") || !strings.Contains(p, "ator=usuario") {
			t.Errorf("sem a policy, o sincronizador gravando evento de outro ator não foi acusado: %s", p)
		}
	})
}
