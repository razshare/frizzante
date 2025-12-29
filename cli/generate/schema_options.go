package generate

import "database/sql"

type SchemaOptions struct {
	Sqlc     string
	SqlcYaml string
	Database *sql.DB
}
