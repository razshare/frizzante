package actions

import (
	"database/sql"

	"github.com/razshare/frizzante/platforms"
)

type MigrateOptions struct {
	Offset   string
	Target   string
	Sqlc     string
	SqlcYaml string
	Database *sql.DB
	Platform platforms.Platform
	Auto     bool
}
