package databases

import (
	"database/sql"
	"embed"
)

type MigrateOptions struct {
	Efs      embed.FS
	Database *sql.DB
	Offset   string
	Target   string
}
