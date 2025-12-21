package generate

import (
	"database/sql"
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
	"gopkg.in/yaml.v3"
)

func Schema(options SchemaOptions) (err error) {
	baseDirectory := filepath.Dir(options.SqlcYaml)

	if _, err = exec.LookPath(options.Sqlc); err != nil && !files.IsFile(options.Sqlc) {
		if err = Sqlc(SqlcOptions{
			Sqlc:     options.Sqlc,
			Platform: options.Platform,
		}); err != nil {
			return
		}
	}

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

	if data, err = os.ReadFile(filepath.Join(baseDirectory, schema)); err != nil {
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
