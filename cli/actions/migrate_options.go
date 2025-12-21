package actions

import (
	"database/sql"

	"github.com/razshare/frizzante/platforms"
)

type MigrateOptions struct {
	Database *sql.DB
	Offset   string
	Target   string
	Sqlc     string
	SqlcYaml string
	Platform platforms.Platform
}
