// Package servidor compõe a API: a camada HTTP da plataforma e as rotas de
// cada módulo, num ponto só. O processo (tenaz api), a geração do contrato
// (tenaz openapi) e os testes de rota montam a API por aqui, para que o
// contrato gerado seja exatamente o que o servidor atende.
package servidor

import (
	"tenaz/internal/modulos/adm"
	admrotas "tenaz/internal/modulos/adm/rotas"
	"tenaz/internal/plataforma/api"
	"tenaz/internal/plataforma/autenticacao"
)

// Dependencias são os serviços que as rotas usam. Na geração do contrato,
// ficam nil: as rotas existem e responderiam 503.
type Dependencias struct {
	Autenticacao *autenticacao.Servico
	ADM          adm.Servico
}

// Montar monta a API com todas as rotas.
func Montar(cfg api.Config, d Dependencias) *api.Servidor {
	s := api.Novo(cfg, d.Autenticacao)
	admrotas.Registrar(s, d.ADM)
	return s
}

// Contrato monta a API só para gerar o contrato OpenAPI.
func Contrato() *api.Servidor {
	return Montar(api.ConfigContrato(), Dependencias{})
}
