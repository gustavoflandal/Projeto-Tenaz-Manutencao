// Package esquema embute os arquivos SQL do banco no binário.
//
// O nome do pacote difere do nome da pasta de propósito: a pasta `db` é a
// definida na estrutura do repositório (seção 7 de
// docs/5_ARQUITETURA_SOFTWARE.md) e guarda SQL, não Go; e `db` já é o nome do
// pacote de acesso ao banco em internal/plataforma/db. Este arquivo existe
// apenas porque go:embed não alcança diretório acima do seu.
package esquema

import "embed"

// Migracoes são os arquivos do goose, aplicados em ordem de número.
//
//go:embed migracoes/*.sql
var Migracoes embed.FS

// DiretorioMigracoes é o caminho dentro de Migracoes.
const DiretorioMigracoes = "migracoes"

// Sementes são os dados de desenvolvimento e de teste: empresas fictícias
// usadas pelas suítes de vazamento e pelo ambiente local.
//
//go:embed all:sementes
var Sementes embed.FS
