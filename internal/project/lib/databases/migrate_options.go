package databases

import (
	"database/sql"
	"embed"
)

type MigrateOptions struct {
	Efs          embed.FS
	Database     *sql.DB
	Direction    MigrationDirection
	FirstVersion string
	Stop         func(version string) bool
	Exists       func(version string) (exists bool, err error)
	Add          func(version string) (err error)
	Remove       func(version string) (err error)
}
