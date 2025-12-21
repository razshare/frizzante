package apps

import (
	"path/filepath"

	"github.com/razshare/frizzante/cli/extensions"
	flag "github.com/spf13/pflag"
)

func New() *App {
	go_ := flag.StringP("go", "", "go"+extensions.Find(), "sets the go binary location")
	air := flag.StringP("air", "", filepath.Join(".gen", "air", "air"+extensions.Find()), "sets the air binary location")
	bun := flag.StringP("bun", "", filepath.Join(".gen", "bun", "bun"+extensions.Find()), "sets the bun binary location")
	sqlc := flag.StringP("sqlc", "", filepath.Join(".gen", "sqlc", "sqlc"+extensions.Find()), "sets the sqlc binary location")
	sqlcYaml := flag.StringP("sqlc-yaml", "", "", "sets the sqlc configuration file location")
	tags := flag.StringP("tags", "", "", "sets build tags")
	databaseConnectionString := flag.StringP("database-connection-string", "", "", "sets the database connection string or file")
	databaseType := flag.StringP("database-type", "", "sqlc", "sets the type of database (currently only sqlite is supported)")

	return &App{
		Go:                       go_,
		Air:                      air,
		Bun:                      bun,
		Sqlc:                     sqlc,
		SqlcYaml:                 sqlcYaml,
		Tags:                     tags,
		DatabaseConnectionString: databaseConnectionString,
		DatabaseType:             databaseType,
	}
}
