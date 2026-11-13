package adm

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"

	"tenaz/internal/modulos/adm/internal/dominio"
	"tenaz/internal/modulos/adm/internal/repo"
	"tenaz/internal/plataforma/autenticacao"
	"tenaz/internal/plataforma/autorizacao"
	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/db/consultas"
	"tenaz/internal/plataforma/erros"
	"tenaz/internal/plataforma/obs"
	"tenaz/internal/plataforma/pais"
)

// CriarConta: RF-ADM-001 e RF-SRV-001 (parcial). Tudo numa transação aberta
// já com o id novo em app.empresa_id (policy empresa_propria); o convite sai
// depois do COMMIT, pelo pool da autenticação.
func (s *servico) CriarConta(ctx context.Context, e EntradaCriarConta) (ContaCriada, error) {
	perfil, ee := dominio.PerfilDoPais("pais", e.Pais)
	if ee != nil {
		return ContaCriada{}, ee
	}
	ident, ee := dominio.IdentificacaoDaEmpresa("", perfil, e.IdentificacaoTipo, e.Identificacao)
	if ee != nil {
		return ContaCriada{}, ee
	}
	// Moeda e idioma padrão vêm do perfil do país e são gravados explicitamente;
	// o idioma padrão pode ser trocado já na criação (é da empresa e o
	// administrador o altera depois). O país e a moeda nunca mudam pela API.
	idiomaPadrao := perfil.IdiomaPadrao
	if e.IdiomaPadrao != nil && strings.TrimSpace(*e.IdiomaPadrao) != "" {
		if idiomaPadrao, ee = dominio.Idioma("idioma_padrao", *e.IdiomaPadrao); ee != nil {
			return ContaCriada{}, ee
		}
	}
	// O responsável fica com o idioma da tela de cadastro, se for suportado.
	idiomaResp := idiomaPadrao
	if pais.IdiomaValido(e.IdiomaResponsavel) {
		idiomaResp = e.IdiomaResponsavel
	}
	razao, ee := dominio.Texto("razao_social", e.RazaoSocial, 200)
	if ee != nil {
		return ContaCriada{}, ee
	}
	fantasia, ee := dominio.TextoOpcional("nome_fantasia", e.NomeFantasia, 200)
	if ee != nil {
		return ContaCriada{}, ee
	}
	if e.FusoHorario == "" {
		e.FusoHorario = perfil.FusoPadrao
	}
	fuso, ee := dominio.FusoDoPais("fuso_horario", perfil, e.FusoHorario)
	if ee != nil {
		return ContaCriada{}, ee
	}
	plano := e.Plano
	if plano == nil && !e.SrvHabilitado {
		plano = ptr(dominio.PlanoEssencial)
	}
	if ee := dominio.PlanoValido(plano, e.SrvHabilitado); ee != nil {
		return ContaCriada{}, ee
	}
	nomeResp, ee := dominio.Texto("responsavel.nome", e.Responsavel.Nome, 200)
	if ee != nil {
		return ContaCriada{}, ee
	}
	emailResp, ee := dominio.Email("responsavel.email", e.Responsavel.Email)
	if ee != nil {
		return ContaCriada{}, ee
	}
	telResp, ee := dominio.TextoOpcional("responsavel.telefone", e.Responsavel.Telefone, 30)
	if ee != nil {
		return ContaCriada{}, ee
	}
	// Sem turno informado, a conta é recusada: o sistema não impõe calendário
	// (RN-ADM-017, RN-ADM-019).
	un, err := validarUnidade("unidade.", e.Unidade, true)
	if err != nil {
		return ContaCriada{}, err
	}
	if ee := dominio.ConferirTipoDaUnidade("unidade.", perfil, un.ident); ee != nil {
		return ContaCriada{}, ee
	}
	if ee := conferirFusoDaUnidade("unidade.", perfil, un.fuso); ee != nil {
		return ContaCriada{}, ee
	}
	if ee := conferirEnderecoDaUnidade("unidade.", perfil, &un, e.Unidade.Endereco); ee != nil {
		return ContaCriada{}, ee
	}

	res := ContaCriada{EmpresaID: novoID(), UnidadeID: novoID(), UsuarioID: novoID()}
	err = s.d.Banco.NaEmpresa(db.ComModulo(ctx, moduloADM), res.EmpresaID, db.Interno(nil), func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		err = q.AdmCriarEmpresa(ctx, consultas.AdmCriarEmpresaParams{
			ID: res.EmpresaID, RazaoSocial: razao, NomeFantasia: fantasia, Pais: perfil.Codigo, Moeda: perfil.Moeda.Codigo(),
			IdiomaPadrao: idiomaPadrao, IdentificacaoTipo: &ident.Tipo, Identificacao: &ident.Valor, Plano: plano,
			SrvHabilitado: e.SrvHabilitado, FusoHorario: fuso,
		})
		if err != nil {
			return repo.Traduzir(err)
		}
		if err := criarUnidadeTx(ctx, q, res.EmpresaID, res.UnidadeID, un); err != nil {
			return err
		}
		perfis, err := criarPerfisPadraoTx(ctx, q, res.EmpresaID)
		if err != nil {
			return err
		}
		err = q.AdmCriarUsuario(ctx, consultas.AdmCriarUsuarioParams{
			EmpresaID: res.EmpresaID, ID: res.UsuarioID, Nome: nomeResp, Email: &emailResp, Telefone: telResp,
			Moeda: perfil.Moeda.Codigo(), Idioma: idiomaResp,
		})
		if err != nil {
			return repo.Traduzir(err)
		}
		// O responsável nasce administrador e gestor: a empresa de uma pessoa
		// só já libera OS (os.encerrar), e ninguém altera os próprios perfis
		// (RN-ADM-017, RN-ADM-025).
		for _, p := range []string{autorizacao.PerfilAdministrador, autorizacao.PerfilGestor} {
			if err := q.AdmColocarPerfilUsuario(ctx, consultas.AdmColocarPerfilUsuarioParams{
				EmpresaID: res.EmpresaID, ID: novoID(), UsuarioID: res.UsuarioID, PerfilID: perfis[p],
			}); err != nil {
				return repo.Traduzir(err)
			}
		}
		return nil
	})
	if err != nil {
		return ContaCriada{}, err
	}
	nome := razao
	if fantasia != nil {
		nome = *fantasia
	}
	res.ConviteEnviado = s.convidar(ctx, res.EmpresaID, res.UsuarioID, emailResp, nome, nomeResp, idiomaResp) == nil
	return res, nil
}

// criarPerfisPadraoTx cria os 12 perfis padrão com a matriz de RF-ADM-005,
// fonte única no pacote autorizacao. Devolve o id de cada um pelo código.
func criarPerfisPadraoTx(ctx context.Context, q *consultas.Queries, empresaID uuid.UUID) (map[string]uuid.UUID, error) {
	ids := map[string]uuid.UUID{}
	for _, p := range autorizacao.PerfisPadrao() {
		id := novoID()
		err := q.AdmCriarPerfil(ctx, consultas.AdmCriarPerfilParams{
			EmpresaID: empresaID, ID: id, Codigo: p.Codigo, Nome: p.Nome, Padrao: true,
			PerfilOrigem: ptr(p.Codigo), Permissoes: autorizacao.PermissoesDoPerfil(p.Codigo),
		})
		if err != nil {
			return nil, fmt.Errorf("criar perfil %s: %w", p.Codigo, repo.Traduzir(err))
		}
		ids[p.Codigo] = id
	}
	return ids, nil
}

// convidar grava e envia o convite, depois do COMMIT que criou o usuário.
// Falha fica no log e volta como erro, mas não desfaz o usuário: o
// administrador reenvia.
func (s *servico) convidar(ctx context.Context, empresaID, usuarioID uuid.UUID, email, nomeEmpresa, nomeConvidado, idioma string) error {
	err := s.d.Autenticacao.EmitirConvite(ctx, empresaID, usuarioID, email, nomeEmpresa, nomeConvidado, idioma)
	if err != nil {
		obs.Logger(ctx).Error("convite não enviado", "empresa_id", empresaID, "usuario_id", usuarioID, "erro", err)
	}
	return err
}

// AceitarConvite liga a identidade ao usuário convidado, na transação da
// empresa do convite. O ator é "sistema": ainda não há sessão.
func (s *servico) AceitarConvite(ctx context.Context, token, senha string, idioma *string) (ConviteAceito, error) {
	// O idioma em que a pessoa aceitou passa a ser o dela (ADR-013).
	if idioma != nil {
		v, ee := dominio.Idioma("idioma", *idioma)
		if ee != nil {
			return ConviteAceito{}, ee
		}
		idioma = &v
	}
	c, err := s.d.Autenticacao.AceitarConvite(ctx, token, senha,
		func(ctx context.Context, c autenticacao.ConviteValido, identidadeID uuid.UUID) error {
			return s.d.Banco.NaEmpresa(db.ComModulo(ctx, moduloADM), c.EmpresaID, db.Interno(nil), func(ctx context.Context) error {
				q, err := repo.Q(ctx)
				if err != nil {
					return err
				}
				n, err := q.AdmAceitarConviteUsuario(ctx, consultas.AdmAceitarConviteUsuarioParams{
					EmpresaID: c.EmpresaID, ID: c.UsuarioID, IdentidadeID: &identidadeID, Idioma: idioma,
				})
				if err != nil {
					return repo.Traduzir(err)
				}
				if n == 0 {
					// Usuário inativado, excluído ou já ativo desde o envio.
					return erros.LinkEmailInvalido
				}
				return nil
			})
		})
	if err != nil {
		return ConviteAceito{}, err
	}
	return ConviteAceito{EmpresaID: c.EmpresaID, UsuarioID: c.UsuarioID, IdentidadeNova: !c.IdentidadeExiste}, nil
}

// ------------------------------------------------------------------ empresa

func (s *servico) ObterEmpresa(ctx context.Context, a autorizacao.Ator) (Empresa, error) {
	if err := apenasInterno(a, autorizacao.AdmLer); err != nil {
		return Empresa{}, err
	}
	var emp Empresa
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		var err error
		emp, err = lerEmpresa(ctx, a.EmpresaID)
		return err
	})
	return emp, err
}

func lerEmpresa(ctx context.Context, id uuid.UUID) (Empresa, error) {
	q, err := repo.Q(ctx)
	if err != nil {
		return Empresa{}, err
	}
	r, err := q.AdmBuscarEmpresa(ctx, id)
	if err != nil {
		return Empresa{}, naoEncontrado(err)
	}
	return Empresa{
		ID: r.ID, RazaoSocial: r.RazaoSocial, NomeFantasia: r.NomeFantasia, Pais: r.Pais, Moeda: r.Moeda,
		IdiomaPadrao: r.IdiomaPadrao, IdentificacaoTipo: r.IdentificacaoTipo, Identificacao: r.Identificacao, Plano: r.Plano,
		SrvHabilitado: r.SrvHabilitado, MFAObrigatorio: r.MfaObrigatorio, Situacao: r.Situacao, FusoHorario: r.FusoHorario,
		RetencaoAuditoriaAnos: int32(r.RetencaoAuditoriaAnos), TesteAte: vazioParaNil(r.TesteAte),
		Versao: r.Versao, CriadoEm: r.CriadoEm, AtualizadoEm: r.AtualizadoEm,
	}, nil
}

func (s *servico) EditarEmpresa(ctx context.Context, a autorizacao.Ator, e EntradaEmpresa) (Empresa, error) {
	if err := apenasInterno(a, autorizacao.AdmGerenciarEmpresa); err != nil {
		return Empresa{}, err
	}
	razao, ee := dominio.Texto("razao_social", e.RazaoSocial, 200)
	if ee != nil {
		return Empresa{}, ee
	}
	fantasia, ee := dominio.TextoOpcional("nome_fantasia", e.NomeFantasia, 200)
	if ee != nil {
		return Empresa{}, ee
	}
	fuso, ee := dominio.FusoHorario("fuso_horario", e.FusoHorario)
	if ee != nil {
		return Empresa{}, ee
	}
	idioma, ee := dominio.Idioma("idioma_padrao", e.IdiomaPadrao)
	if ee != nil {
		return Empresa{}, ee
	}
	var emp Empresa
	err := s.naEmpresa(ctx, a, func(ctx context.Context) error {
		q, err := repo.Q(ctx)
		if err != nil {
			return err
		}
		atual, err := q.AdmTravarEmpresa(ctx, a.EmpresaID)
		if err != nil {
			return naoEncontrado(err)
		}
		if err := conferirVersao(e.VersaoEsperada, atual.Versao); err != nil {
			return err
		}
		// A identificação é do tipo do país da empresa, que não muda.
		perfil, ok := pais.Perfil(atual.Pais)
		if !ok {
			return fmt.Errorf("empresa com país %q fora do perfil de país", atual.Pais)
		}
		ident, ee := dominio.IdentificacaoDaEmpresa("", perfil, e.IdentificacaoTipo, e.Identificacao)
		if ee != nil {
			return ee
		}
		// O fuso é um dos do país da empresa.
		if _, ee := dominio.FusoDoPais("fuso_horario", perfil, fuso); ee != nil {
			return ee
		}
		err = q.AdmAtualizarEmpresa(ctx, consultas.AdmAtualizarEmpresaParams{
			ID: a.EmpresaID, RazaoSocial: razao, NomeFantasia: fantasia, IdentificacaoTipo: &ident.Tipo,
			Identificacao: &ident.Valor, IdiomaPadrao: idioma, FusoHorario: fuso, MfaObrigatorio: e.MFAObrigatorio,
		})
		if err != nil {
			return repo.Traduzir(err)
		}
		emp, err = lerEmpresa(ctx, a.EmpresaID)
		return err
	})
	return emp, err
}

// -------------------------------------------------------- unidade (comum)

type unidadeValida struct {
	codigo, nome string
	ident        *dominio.Identificacao
	fuso         *string
	endereco     []byte
	lat, lon     *float64
	turnos       []dominio.Turno
	nomesTurnos  []*string
}

func validarUnidade(prefixo string, e EntradaUnidade, exigirTurnos bool) (unidadeValida, error) {
	var u unidadeValida
	var ee *erros.Erro
	if u.codigo, ee = dominio.Texto(prefixo+"codigo", e.Codigo, 20); ee != nil {
		return u, ee
	}
	if u.nome, ee = dominio.Texto(prefixo+"nome", e.Nome, 200); ee != nil {
		return u, ee
	}
	if u.ident, ee = dominio.IdentificacaoDaUnidade(prefixo, e.IdentificacaoTipo, e.Identificacao); ee != nil {
		return u, ee
	}
	if e.FusoHorario != nil && *e.FusoHorario != "" {
		f, ee := dominio.FusoHorario(prefixo+"fuso_horario", *e.FusoHorario)
		if ee != nil {
			return u, ee
		}
		u.fuso = &f
	}
	if e.Latitude != nil && (*e.Latitude < -90 || *e.Latitude > 90) {
		return u, erros.EntradaInvalida.ComCampoChave(prefixo+"latitude", erros.CampoForaDaFaixa, "campo.latitude", nil)
	}
	if e.Longitude != nil && (*e.Longitude < -180 || *e.Longitude > 180) {
		return u, erros.EntradaInvalida.ComCampoChave(prefixo+"longitude", erros.CampoForaDaFaixa, "campo.longitude", nil)
	}
	u.lat, u.lon = e.Latitude, e.Longitude
	if e.Endereco != nil {
		b, err := json.Marshal(e.Endereco)
		if err != nil {
			return u, err
		}
		u.endereco = b
	}
	if exigirTurnos && len(e.Turnos) == 0 {
		return u, dominio.ErrUltimoTurno()
	}
	for i, t := range e.Turnos {
		dt := dominio.Turno{DiaSemana: t.DiaSemana, HoraInicio: t.HoraInicio, HoraFim: t.HoraFim}
		if ee := dominio.ValidarTurno(fmt.Sprintf("%sturnos[%d]", prefixo, i), dt); ee != nil {
			return u, ee
		}
		nome, ee := dominio.TextoOpcional(fmt.Sprintf("%sturnos[%d].nome", prefixo, i), t.Nome, 60)
		if ee != nil {
			return u, ee
		}
		u.turnos = append(u.turnos, dt)
		u.nomesTurnos = append(u.nomesTurnos, nome)
	}
	if i, j := dominio.PrimeiraSobreposicao(u.turnos); i >= 0 {
		return u, dominio.ErrSobreposicao(fmt.Sprintf("%sturnos[%d]", prefixo, j), u.turnos[i])
	}
	return u, nil
}

// conferirFusoDaUnidade confere o fuso da unidade, se informado, contra a lista
// do país da empresa (a unidade não tem país próprio, Q1).
func conferirFusoDaUnidade(prefixo string, perfil pais.PerfilPais, fuso *string) *erros.Erro {
	if fuso == nil {
		return nil
	}
	_, ee := dominio.FusoDoPais(prefixo+"fuso_horario", perfil, *fuso)
	return ee
}

// conferirEnderecoDaUnidade confere subdivisão e código postal do endereço
// contra o perfil do país da empresa e regrava o endereço com o código postal
// normalizado (documento 3, 2.18). Sem endereço, não há o que conferir.
func conferirEnderecoDaUnidade(prefixo string, perfil pais.PerfilPais, u *unidadeValida, e *Endereco) *erros.Erro {
	if e == nil {
		return nil
	}
	n, ee := dominio.SubdivisaoECodigoPostal(prefixo+"endereco.", perfil, e.Subdivisao, e.CodigoPostal)
	if ee != nil {
		return ee
	}
	norm := *e
	norm.CodigoPostal = n
	b, err := json.Marshal(norm)
	if err != nil {
		return erros.Interno
	}
	u.endereco = b
	return nil
}

func (u unidadeValida) identTipo() *string {
	if u.ident == nil {
		return nil
	}
	return &u.ident.Tipo
}

func (u unidadeValida) identValor() *string {
	if u.ident == nil {
		return nil
	}
	return &u.ident.Valor
}

func criarUnidadeTx(ctx context.Context, q *consultas.Queries, empresaID, unidadeID uuid.UUID, u unidadeValida) error {
	err := q.AdmCriarUnidade(ctx, consultas.AdmCriarUnidadeParams{
		EmpresaID: empresaID, ID: unidadeID, Codigo: u.codigo, Nome: u.nome, IdentificacaoTipo: u.identTipo(), Identificacao: u.identValor(),
		Endereco: u.endereco, FusoHorario: u.fuso, Latitude: u.lat, Longitude: u.lon,
	})
	if err != nil {
		return repo.Traduzir(err)
	}
	for i, t := range u.turnos {
		err := q.AdmCriarTurno(ctx, consultas.AdmCriarTurnoParams{
			EmpresaID: empresaID, ID: novoID(), UnidadeID: unidadeID, DiaSemana: int16(t.DiaSemana),
			HoraInicio: t.HoraInicio, HoraFim: t.HoraFim, Nome: u.nomesTurnos[i],
		})
		if err != nil {
			return repo.Traduzir(err)
		}
	}
	return nil
}
