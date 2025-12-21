package generate

import (
	"database/sql"

	"github.com/razshare/frizzante/platforms"
)

type SchemaOptions struct {
	Sqlc     string
	SqlcYaml string
	Database *sql.DB
	Platform platforms.Platform
}
