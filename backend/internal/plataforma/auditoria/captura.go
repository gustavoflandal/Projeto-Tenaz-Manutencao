package auditoria

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"

	"tenaz/internal/plataforma/db"
	"tenaz/internal/plataforma/db/consultas"
	"tenaz/internal/plataforma/erros"
)

// ErrTelaDesconhecida: a tela não está no catálogo fechado. Código
// validacao_falhou, campo tela.
var ErrTelaDesconhecida = erros.EntradaInvalida.ComCampoChave("tela", erros.CampoValorInvalido, "campo.tela_desconhecida", nil)

// Chaves são as três chaves de uma tela (RF-ADM-011): só estes três tipos são
// configuráveis. Login, logout, inclusão, alteração, exclusão, deleção e a
// visualização de dado sensível nunca se desligam.
type Chaves struct {
	Visualizacao bool
	Impressao    bool
	Extracao     bool
}

// ChavesPadrao vale para a tela que ainda não tem linha em captura_auditoria:
// impressão e extração ligadas, visualização desligada (documento 2, seção
// 9.2.2). Os mesmos valores são o DEFAULT das colunas.
var ChavesPadrao = Chaves{Visualizacao: false, Impressao: true, Extracao: true}

// Captura é a configuração de uma tela, com os dados do catálogo.
type Captura struct {
	Tela   Tela
	Chaves Chaves
	// Configurada: a empresa já mudou alguma chave desta tela; senão vale o
	// padrão.
	Configurada bool
}

// Mudanca é a alteração pedida: chave ausente não muda.
type Mudanca struct {
	Visualizacao *bool
	Impressao    *bool
	Extracao     *bool
}

// Tipo de operação configurável.
type Tipo string

// Os três tipos configuráveis.
const (
	TipoVisualizacao Tipo = AcaoVisualizacao
	TipoImpressao    Tipo = AcaoImpressao
	TipoExtracao     Tipo = AcaoExtracao
)

func (c Chaves) liga(t Tipo) bool {
	switch t {
	case TipoVisualizacao:
		return c.Visualizacao
	case TipoImpressao:
		return c.Impressao
	case TipoExtracao:
		return c.Extracao
	}
	return false
}

// Dependencias do serviço.
type Dependencias struct {
	Banco *db.DB
	Agora func() time.Time
	// ValidadeCache é por quanto tempo uma instância guarda a configuração de
	// uma empresa sem reler (padrão 30 s). A alteração na própria instância
	// invalida na hora; outra instância enxerga em até esta validade.
	ValidadeCache time.Duration
}

type entradaCache struct {
	mapa     map[string]Chaves
	expiraEm time.Time
}

// Servico é a captura por tela e a gravação de visualização, impressão e
// extração. Sem estado global: o cache é do serviço.
type Servico struct {
	d     Dependencias
	mu    sync.Mutex
	cache map[uuid.UUID]entradaCache
}

// Novo monta o serviço.
func Novo(d Dependencias) (*Servico, error) {
	if d.Banco == nil {
		return nil, errors.New("auditoria: banco ausente")
	}
	if d.Agora == nil {
		d.Agora = time.Now
	}
	if d.ValidadeCache <= 0 {
		d.ValidadeCache = 30 * time.Second
	}
	return &Servico{d: d, cache: map[uuid.UUID]entradaCache{}}, nil
}

// ListarCapturas devolve todas as telas do catálogo com a configuração da
// empresa (o padrão onde não há linha). Na transação do contexto, escopo
// interno (o externo não lê a tabela).
func (s *Servico) ListarCapturas(ctx context.Context, empresaID uuid.UUID) ([]Captura, error) {
	mapa, err := lerCapturas(ctx, empresaID)
	if err != nil {
		return nil, err
	}
	res := make([]Captura, 0, len(telas))
	for _, t := range telas {
		c, ok := mapa[t.Codigo]
		if !ok {
			c = ChavesPadrao
		}
		res = append(res, Captura{Tela: t, Chaves: c, Configurada: ok})
	}
	return res, nil
}

// AlterarCaptura aplica a mudança e grava SEMPRE as três chaves (o upsert não
// pode devolver a outra chave ao padrão). A tela fora do catálogo é recusada.
// Na transação do contexto; o chamador chama InvalidarCapturas depois do
// COMMIT. A troca da chave aparece na trilha como inclusao ou alteracao de
// captura_auditoria, pelo trg_auditoria.
func (s *Servico) AlterarCaptura(ctx context.Context, empresaID uuid.UUID, usuarioID *uuid.UUID, tela string, m Mudanca) (Captura, error) {
	t, ok := BuscarTela(tela)
	if !ok {
		return Captura{}, ErrTelaDesconhecida
	}
	q, err := db.Q(ctx)
	if err != nil {
		return Captura{}, err
	}
	mapa, err := lerCapturas(ctx, empresaID)
	if err != nil {
		return Captura{}, err
	}
	c, configurada := mapa[tela]
	if !configurada {
		c = ChavesPadrao
	}
	if m.Visualizacao != nil {
		c.Visualizacao = *m.Visualizacao
	}
	if m.Impressao != nil {
		c.Impressao = *m.Impressao
	}
	if m.Extracao != nil {
		c.Extracao = *m.Extracao
	}
	if err := q.AudGravarCaptura(ctx, consultas.AudGravarCapturaParams{
		EmpresaID: empresaID, Tela: tela, Visualizacao: c.Visualizacao, Impressao: c.Impressao, Extracao: c.Extracao,
		UsuarioID: usuarioID,
	}); err != nil {
		return Captura{}, fmt.Errorf("gravar captura da tela %s: %w", tela, err)
	}
	return Captura{Tela: t, Chaves: c, Configurada: true}, nil
}

func lerCapturas(ctx context.Context, empresaID uuid.UUID) (map[string]Chaves, error) {
	q, err := db.Q(ctx)
	if err != nil {
		return nil, err
	}
	linhas, err := q.AudListarCapturas(ctx, empresaID)
	if err != nil {
		return nil, fmt.Errorf("ler capturas de auditoria: %w", err)
	}
	mapa := make(map[string]Chaves, len(linhas))
	for _, l := range linhas {
		mapa[l.Tela] = Chaves{Visualizacao: l.Visualizacao, Impressao: l.Impressao, Extracao: l.Extracao}
	}
	return mapa, nil
}

// InvalidarCapturas descarta o cache da empresa. Chamar depois do COMMIT de
// AlterarCaptura.
func (s *Servico) InvalidarCapturas(empresaID uuid.UUID) {
	s.mu.Lock()
	delete(s.cache, empresaID)
	s.mu.Unlock()
}

// ChavesDaTela devolve as chaves vigentes da tela para a empresa. É lido num
// contexto interno próprio, FORA da transação do chamador (db.SemTransacao), e
// guardado em cache: a requisição de prestador, locatário ou link não lê
// captura_auditoria, e esta leitura serve também a ela. Tela vazia ou fora do
// catálogo: o padrão.
func (s *Servico) ChavesDaTela(ctx context.Context, empresaID uuid.UUID, tela string) (Chaves, error) {
	if tela == "" || !TelaValida(tela) {
		return ChavesPadrao, nil
	}
	mapa, err := s.capturasDaEmpresa(ctx, empresaID)
	if err != nil {
		return Chaves{}, err
	}
	if c, ok := mapa[tela]; ok {
		return c, nil
	}
	return ChavesPadrao, nil
}

func (s *Servico) capturasDaEmpresa(ctx context.Context, empresaID uuid.UUID) (map[string]Chaves, error) {
	agora := s.d.Agora()
	s.mu.Lock()
	e, ok := s.cache[empresaID]
	s.mu.Unlock()
	if ok && agora.Before(e.expiraEm) {
		return e.mapa, nil
	}
	var mapa map[string]Chaves
	err := s.d.Banco.NaEmpresa(db.SemTransacao(ctx), empresaID, db.Interno(nil), func(ctx context.Context) error {
		var err error
		mapa, err = lerCapturas(ctx, empresaID)
		return err
	})
	if err != nil {
		return nil, fmt.Errorf("carregar capturas da empresa: %w", err)
	}
	s.mu.Lock()
	s.cache[empresaID] = entradaCache{mapa: mapa, expiraEm: agora.Add(s.d.ValidadeCache)}
	s.mu.Unlock()
	return mapa, nil
}
