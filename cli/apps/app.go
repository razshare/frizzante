package apps

import "embed"

type App struct {
	Go                       *string
	Air                      *string
	Bun                      *string
	Sqlc                     *string
	SqlcYaml                 *string
	Tags                     *string
	DatabaseConnectionString *string
	DatabaseType             *string
	Efs                      embed.FS
}
