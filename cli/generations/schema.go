package generations

import (
	"database/sql"
	"errors"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
	"gopkg.in/yaml.v3"
)

func Schema(options SchemaOptions) (err error) {
	directoryName := filepath.Dir(options.SqlcYaml)

	var data []byte
	if data, err = os.ReadFile(options.SqlcYaml); err != nil {
		return
	}

	type Configuration struct {
		Sql []struct {
			Schema string `yaml:"schema"`
		} `yaml:"sql"`
	}

	var config Configuration
	if err = yaml.Unmarshal(data, &config); err != nil {
		return
	}

	var schema string
	if len(config.Sql) == 0 {
		err = errors.New("sqlc configuration is missing a schema definition")
		return
	}

	schema = config.Sql[0].Schema

	spin := spinners.Newf("migrating database schema using %s", schema)
	go spinners.Start(spin)

	if data, err = os.ReadFile(filepath.Join(directoryName, schema)); err != nil {
		spinners.Stop(spin)
		return
	}

	var transaction *sql.Tx
	if transaction, err = options.Database.Begin(); err != nil {
		return
	}

	if query := string(data); query != "" {
		if _, err = transaction.Exec(query); err != nil {
			if rerr := transaction.Rollback(); rerr != nil {
				err = rerr
				spinners.Stop(spin)
				return
			}
			spinners.Stop(spin)
			return
		}
	}

	spinners.Stop(spin)

	if err = transaction.Commit(); err != nil {
		return
	}

	messages.Success("database schema updated successfully")

	return
}
