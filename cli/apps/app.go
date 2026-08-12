package apps

import "embed"

type App struct {
	Efs          embed.FS
	Go           *string
	Air          *string
	Bun          *string
	Sqlc         *string
	SqlcYaml     *string
	Tags         *string
	Database     *string
	DatabaseType *string
	Context      *string
	Output       *string
	AskModel     *string
	AskHost      *string
	AskProtocol  *string
	Strict       *bool
	Incremental  *bool
}
