package adm

import (
	"bytes"
	"context"
	"encoding/csv"
	"net/netip"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"

	"tenaz/internal/modulos/adm/internal/repo"
	"tenaz/internal/plataforma/auditoria"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/db/consultas"
	"tenaz/internal/plataforma/erros"
	"tenaz/internal/plataforma/paginacao"
	"tenaz/internal/plataforma/textos"
)

// Trilha de auditoria: RF-ADM-011. Só leitura; nenhum perfil altera evento.

// acoesAuditoria são os nove tipos de operação (sumário 3.6).
var acoesAuditoria = []string{
	auditoria.AcaoLogin, auditoria.AcaoLogout, auditoria.AcaoInclusao, auditoria.AcaoAlteracao, auditoria.AcaoExclusao,
	auditoria.AcaoDelecao, auditoria.AcaoVisualizacao, auditoria.AcaoImpressao, auditoria.AcaoExtracao,
}

// MaxEventosExportacao é o teto de linhas de uma exportação CSV da trilha.
// Acima dele a pessoa estreita o período; a exportação completa e verificável
// é a de RF-ADM-013.
const MaxEventosExportacao = 50000

// filtroValidado confere o filtro e devolve o IP normalizado.
func filtroValidado(f FiltroAuditoria) (FiltroAuditoria, error) {
	if f.Acao != nil && !slices.Contains(acoesAuditoria, *f.Acao) {
		return f, erros.EntradaInvalida.ComCampoChave("acao", erros.CampoValorInvalido, "campo.acao_desconhecida", nil)
	}
	if f.Modulo != nil {
		m := strings.ToLower(*f.Modulo)
		if !db.ModuloValido(m) {
			return f, erros.EntradaInvalida.ComCampoChave("modulo", erros.CampoValorInvalido, "campo.modulo_desconhecido", nil)
		}
		f.Modulo = &m
	}
	if f.Tela != nil && !auditoria.TelaValida(*f.Tela) {
		return f, erros.EntradaInvalida.ComCampoChave("tela", erros.CampoValorInvalido, "campo.tela_desconhecida", nil)
	}
	if f.IP != nil {
		a, err := netip.ParseAddr(strings.TrimSpace(*f.IP))
		if err != nil {
			return f, erros.EntradaInvalida.ComCampoChave("ip", erros.CampoFormatoInvalido, "campo.ip_formato", nil)
		}
		ip := a.Unmap().String()
		f.IP = &ip
		if f.De == nil {
			return f, erros.EntradaInvalida.ComCampoChave("ocorrido_em[gte]", erros.CampoObrigatorio, "campo.ip_exige_inicio", nil)
		}
	}
	if f.De != nil && f.Ate != nil && !f.De.Before(*f.Ate) {
		return f, erros.EntradaInvalida.ComCampoChave("ate", erros.CampoValorInvalido, "campo.periodo_fim_depois", nil)
	}
	return f, nil
}

func instanteTexto(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.UTC().Format(time.RFC3339Nano)
}

func textoSimples(p *string) string { return textoOuVazio(p, func(s string) string { return s }) }

func paramsListar(a autorizacao.Ator, f FiltroAuditoria, cEm *time.Time, cID *uuid.UUID, limite int32) consultas.AdmListarEventosAuditoriaParams {
	return consultas.AdmListarEventosAuditoriaParams{
		EmpresaID: a.EmpresaID, Entidade: f.Entidade, EntidadeID: f.EntidadeID, UsuarioID: f.UsuarioID,
		Acao: f.Acao, Modulo: f.Modulo, Tela: f.Tela, Ip: f.IP, De: f.De, Ate: f.Ate,
		CursorOcorridoEm: cEm, CursorID: cID, Limite: limite,
	}
}

func resumoDe(r consultas.AdmListarEventosAuditoriaRow) ResumoEvento {
	return ResumoEvento{
		ID: r.ID, OcorridoEm: r.OcorridoEm, AtorTipo: r.AtorTipo, UsuarioID: r.UsuarioID, UsuarioNome: r.UsuarioNome,
		TokenAPIID: r.TokenApiID, Acao: r.Acao, Resultado: r.Resultado, Entidade: r.Entidade, EntidadeID: r.EntidadeID,
		EntidadeChave: r.EntidadeChave, CamposAlterados: r.CamposAlterados, Modulo: r.Modulo, Tela: r.Tela,
		Sensivel: r.Sensivel, Origem: r.Origem, IP: vazioParaNil(r.Ip), DispositivoID: r.DispositivoID,
		RequestID: r.RequestID, ConexaoEmpresaID: r.ConexaoEmpresaID, OcorridoNoDispositivoEm: r.OcorridoNoDispositivoEm,
	}
}

func (s *servico) ListarEventosAuditoria(ctx context.Context, a autorizacao.Ator, f FiltroAuditoria) (paginacao.Pagina[ResumoEvento], error) {
	var pag paginacao.Pagina[ResumoEvento]
	if err := apenasInterno(a, autorizacao.AdmLerAuditoria); err != nil {
		return pag, err
	}
	f, err := filtroValidado(f)
	if err != nil {
		return pag, err
	}
	consulta := paginacao.Assinatura("eventos-auditoria",
		textoSimples(f.Entidade), textoOuVazio(f.EntidadeID, uuid.UUID.String), textoOuVazio(f.UsuarioID, uuid.UUID.String),
		textoSimples(f.Acao), textoSimples(f.Modulo), textoSimples(f.Tela), textoSimples(f.IP),
		instanteTexto(f.De), instanteTexto(f.Ate))
	agora := s.agora()
	pos, err := paginacao.Decodificar(f.Cursor, consulta, agora)
	if err != nil {
		return pag, err
	}
	limite := paginacao.Limite(f.Limite)
	cEm, cID := cursorParams(pos)
	err = s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		linhas, err := q.AdmListarEventosAuditoria(ctx, paramsListar(a, f, cEm, cID, limite+1))
		if err != nil {
			return err
		}
		res := make([]ResumoEvento, 0, len(linhas))
		for _, r := range linhas {
			res = append(res, resumoDe(r))
		}
		pag = paginacao.Montar(res, limite, consulta, agora, func(e ResumoEvento) paginacao.Posicao {
			return paginacao.Posicao{Instante: e.OcorridoEm, ID: e.ID}
		})
		return nil
	})
	return pag, err
}

func (s *servico) ObterEventoAuditoria(ctx context.Context, a autorizacao.Ator, id uuid.UUID) (Evento, error) {
	if err := apenasInterno(a, autorizacao.AdmLerAuditoria); err != nil {
		return Evento{}, err
	}
	var ev Evento
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		r, err := q.AdmBuscarEventoAuditoria(ctx, consultas.AdmBuscarEventoAuditoriaParams{EmpresaID: a.EmpresaID, ID: id})
		if err != nil {
			return naoEncontrado(err)
		}
		ev = Evento{
			ResumoEvento: ResumoEvento{
				ID: r.ID, OcorridoEm: r.OcorridoEm, AtorTipo: r.AtorTipo, UsuarioID: r.UsuarioID, UsuarioNome: r.UsuarioNome,
				TokenAPIID: r.TokenApiID, Acao: r.Acao, Resultado: r.Resultado, Entidade: r.Entidade, EntidadeID: r.EntidadeID,
				EntidadeChave: r.EntidadeChave, CamposAlterados: r.CamposAlterados, Modulo: r.Modulo, Tela: r.Tela,
				Sensivel: r.Sensivel, Origem: r.Origem, IP: vazioParaNil(r.Ip), DispositivoID: r.DispositivoID,
				RequestID: r.RequestID, ConexaoEmpresaID: r.ConexaoEmpresaID, OcorridoNoDispositivoEm: r.OcorridoNoDispositivoEm,
			},
			Antes: r.Antes, Depois: r.Depois, Detalhe: r.Detalhe, Hash: r.Hash, UserAgent: r.UserAgent,
		}
		return s.visualizou(ctx, a, "evento_auditoria", id, false)
	})
	return ev, err
}

// visualizou registra a abertura de um registro (RF-ADM-011): sempre, se
// mostra dado sensível; senão, só com a captura de visualização da tela de
// origem ligada.
func (s *servico) visualizou(ctx context.Context, a autorizacao.Ator, entidade string, id uuid.UUID, sensivel bool) error {
	return s.visualizouComDetalhe(ctx, a, entidade, id, sensivel, nil)
}

func (s *servico) visualizouComDetalhe(ctx context.Context, a autorizacao.Ator, entidade string, id uuid.UUID, sensivel bool, detalhe any) error {
	_, err := s.d.Auditoria.Visualizou(ctx, auditoria.Operacao{
		EmpresaID: a.EmpresaID, UsuarioID: ptr(a.UsuarioID), Entidade: entidade, EntidadeID: ptr(id),
		Sensivel: sensivel, Detalhe: detalhe,
	})
	return err
}

// colunasCSV é o cabeçalho da exportação da trilha; a ordem é parte do contrato.
var colunasCSV = []string{"id", "ocorrido_em", "usuario_id", "usuario_nome", "ator_tipo", "ip", "acao", "resultado",
	"modulo", "tela", "entidade", "entidade_id", "entidade_chave", "campos_alterados", "sensivel", "origem",
	"request_id", "hash"}

// ExportarEventosAuditoria gera o CSV com os mesmos filtros da consulta, do
// mais recente ao mais antigo. A exportação é uma extracao: segue a captura da
// tela de origem, inclusive a da própria trilha (RF-ADM-011, critério 4).
// Exige o período inteiro e no máximo MaxEventosExportacao linhas.
func (s *servico) ExportarEventosAuditoria(ctx context.Context, a autorizacao.Ator, f FiltroAuditoria) (ExportacaoAuditoria, error) {
	var res ExportacaoAuditoria
	if err := apenasInterno(a, autorizacao.AdmLerAuditoria); err != nil {
		return res, err
	}
	f, err := filtroValidado(f)
	if err != nil {
		return res, err
	}
	if f.De == nil {
		return res, erros.EntradaInvalida.ComCampoChave("ocorrido_em[gte]", erros.CampoObrigatorio, "campo.exportacao_inicio", nil)
	}
	if f.Ate == nil {
		return res, erros.EntradaInvalida.ComCampoChave("ocorrido_em[lt]", erros.CampoObrigatorio, "campo.exportacao_fim", nil)
	}
	err = s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		linhas, err := q.AdmListarEventosAuditoria(ctx, paramsListar(a, f, nil, nil, MaxEventosExportacao+1))
		if err != nil {
			return err
		}
		if len(linhas) > MaxEventosExportacao {
			return erros.EntradaInvalida.ComCampoChave("ocorrido_em[gte]", erros.CampoValorInvalido, "campo.exportacao_grande", nil)
		}
		var buf bytes.Buffer
		w := csv.NewWriter(&buf)
		_ = w.Write(colunasCSV)
		for _, r := range linhas {
			_ = w.Write([]string{
				r.ID.String(), r.OcorridoEm.UTC().Format("2006-01-02T15:04:05.000000Z"), textoUUID(r.UsuarioID),
				celula(textoPtr(r.UsuarioNome)), r.AtorTipo, r.Ip, r.Acao, r.Resultado, textoPtr(r.Modulo), textoPtr(r.Tela),
				celula(r.Entidade), textoUUID(r.EntidadeID), celula(textoPtr(r.EntidadeChave)),
				celula(strings.Join(r.CamposAlterados, ";")), boolTexto(r.Sensivel), textoPtr(r.Origem),
				textoPtr(r.RequestID), hexa(r.Hash),
			})
		}
		w.Flush()
		if err := w.Error(); err != nil {
			return err
		}
		res = ExportacaoAuditoria{CSV: buf.Bytes(), Quantidade: len(linhas)}
		_, err = s.d.Auditoria.Extraiu(ctx, auditoria.Operacao{
			EmpresaID: a.EmpresaID, UsuarioID: ptr(a.UsuarioID), Entidade: "evento_auditoria",
			Detalhe: map[string]any{"formato": "csv", "quantidade": len(linhas), "filtros": filtrosTexto(f)},
		})
		return err
	})
	return res, err
}

// filtrosTexto devolve os filtros informados, para o detalhe da extração.
func filtrosTexto(f FiltroAuditoria) map[string]string {
	m := map[string]string{}
	for k, v := range map[string]string{
		"entidade": textoSimples(f.Entidade), "entidade_id": textoOuVazio(f.EntidadeID, uuid.UUID.String),
		"usuario_id": textoOuVazio(f.UsuarioID, uuid.UUID.String), "acao": textoSimples(f.Acao),
		"modulo": textoSimples(f.Modulo), "tela": textoSimples(f.Tela), "ip": textoSimples(f.IP),
		"de": instanteTexto(f.De), "ate": instanteTexto(f.Ate),
	} {
		if v != "" {
			m[k] = v
		}
	}
	return m
}

// celula protege contra injeção de fórmula em planilha: texto livre que
// começa por = + - @ ou controle ganha um apóstrofo na frente.
func celula(s string) string {
	if s != "" && strings.ContainsRune("=+-@\t\r", rune(s[0])) {
		return "'" + s
	}
	return s
}

func textoPtr(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

func textoUUID(id *uuid.UUID) string {
	if id == nil {
		return ""
	}
	return id.String()
}

func boolTexto(b bool) string {
	if b {
		return "true"
	}
	return "false"
}

const digitosHex = "0123456789abcdef"

func hexa(b []byte) string {
	out := make([]byte, 0, len(b)*2)
	for _, c := range b {
		out = append(out, digitosHex[c>>4], digitosHex[c&0x0f])
	}
	return string(out)
}

// VerificarIntegridadeAuditoria devolve o estado de cada dia do período
// (inclusive, dias UTC): integro, violado (com o motivo e o primeiro evento
// alterado, quando há como apontar) ou nao_selado.
func (s *servico) VerificarIntegridadeAuditoria(ctx context.Context, a autorizacao.Ator, de, ate time.Time) ([]DiaVerificado, error) {
	if err := apenasInterno(a, autorizacao.AdmLerAuditoria); err != nil {
		return nil, err
	}
	if s.d.Verificador == nil {
		return nil, erros.ServicoIndisponivel.ComChave("detalhe.verificacao_nao_configurada", nil)
	}
	dias, err := s.d.Verificador.Verificar(db.ComModulo(ctx, moduloADM), a.EmpresaID, de, ate)
	if err != nil {
		return nil, err
	}
	res := make([]DiaVerificado, 0, len(dias))
	for _, d := range dias {
		res = append(res, DiaVerificado{Dia: d.Dia, Estado: d.Estado, Motivo: d.Motivo, EventoID: d.EventoID})
	}
	return res, nil
}

// capturaDe monta a captura com o nome da tela no idioma da requisição, pelo
// catálogo de textos do servidor (chave tela.<código>, ADR-013). O código da
// tela não muda; sem texto no catálogo, vale o nome escrito no código.
func capturaDe(ctx context.Context, c auditoria.Captura) CapturaAuditoria {
	nome := c.Tela.Nome
	if chave := "tela." + c.Tela.Codigo; textos.Padrao().Existe(textos.PtBR, chave) {
		nome = textos.T(textos.DoContexto(ctx), chave, nil)
	}
	return CapturaAuditoria{
		Tela: c.Tela.Codigo, Modulo: c.Tela.Modulo, Nome: nome, Visualizacao: c.Chaves.Visualizacao,
		Impressao: c.Chaves.Impressao, Extracao: c.Chaves.Extracao, Configurada: c.Configurada,
	}
}

// ListarCapturasAuditoria devolve todas as telas do catálogo com a captura da
// empresa (o padrão onde ela não mudou nada). adm.gerenciar_parametros.
func (s *servico) ListarCapturasAuditoria(ctx context.Context, a autorizacao.Ator) ([]CapturaAuditoria, error) {
	if err := apenasInterno(a, autorizacao.AdmGerenciarParametros); err != nil {
		return nil, err
	}
	var res []CapturaAuditoria
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		l, err := s.d.Auditoria.ListarCapturas(ctx, a.EmpresaID)
		if err != nil {
			return err
		}
		for _, c := range l {
			res = append(res, capturaDe(ctx, c))
		}
		return nil
	})
	return res, err
}

// AlterarCapturaAuditoria liga ou desliga, para uma tela, a captura de
// visualização, impressão e extração. Grava sempre as três chaves. A troca é
// uma alteração sempre gravada (trg_auditoria) e só vale para eventos futuros.
func (s *servico) AlterarCapturaAuditoria(ctx context.Context, a autorizacao.Ator, tela string, c EntradaCaptura) (CapturaAuditoria, error) {
	if err := apenasInterno(a, autorizacao.AdmGerenciarParametros); err != nil {
		return CapturaAuditoria{}, err
	}
	var res CapturaAuditoria
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		r, err := s.d.Auditoria.AlterarCaptura(ctx, a.EmpresaID, ptr(a.UsuarioID), tela, auditoria.Mudanca{
			Visualizacao: &c.Visualizacao, Impressao: &c.Impressao, Extracao: &c.Extracao,
		})
		res = capturaDe(ctx, r)
		return err
	})
	if err == nil {
		s.d.Auditoria.InvalidarCapturas(a.EmpresaID)
	}
	return res, err
}
