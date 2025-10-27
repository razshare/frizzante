package generate

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/select_one"
	"github.com/razshare/frizzante/tui/spinners"
	"gopkg.in/yaml.v3"
)

func Schema(options SchemaOptions) (err error) {
	yamlFileName := options.SqlcYaml

	if yamlFileName != "" && !files.IsFile(yamlFileName) {
		messages.Infof("%s not found", yamlFileName)
	}

	if yamlFileName == "" {
		var names []string
		if names, err = files.FindWithSuffix("lib", "sqlc.yaml"); err != nil {
			return
		}

		choices := make([]search.Choice, len(names))
		for index, name := range names {
			choices[index] = search.Choice{Id: name}
		}

		choices = append(choices, search.Choice{Id: "other", Description: "use a different file"})

		if options.Auto {
			yamlFileName = choices[0].Id
		} else {
			yamlFileName, err = select_one.Sendf(choices, "where is your sqlc.yaml file located?")
		}
	}

	baseDirectory := filepath.Dir(yamlFileName)

	if _, err = exec.LookPath(options.Sqlc); err != nil && !files.IsFile(options.Sqlc) {
		if err = Sqlc(SqlcOptions{
			Sqlc:     options.Sqlc,
			Platform: options.Platform,
			Auto:     options.Auto,
		}); err != nil {
			return
		}
	}

	if !files.IsFile(yamlFileName) {
		return fmt.Errorf("%s not found", yamlFileName)
	}

	var data []byte
	if data, err = os.ReadFile(yamlFileName); err != nil {
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
