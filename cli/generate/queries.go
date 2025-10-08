package generate

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/singleselect"
	"github.com/razshare/frizzante/tui/spinner"
)

func Queries(options QueriesOptions) (err error) {
	if options.SqlcYaml != "" && !files.IsFile(options.SqlcYaml) {
		messages.Infof("%s not found", options.SqlcYaml)
	}

	if options.SqlcYaml == "" {
		choices := make([]search.Choice, 0)

		if files.IsFile(filepath.Join("lib", "database", "sqlc.yaml")) {
			choices = append(choices, search.Choice{Id: "lib/database/sqlc.yaml", Description: "lib/database/sqlc.yaml"})
		}

		if files.IsFile(filepath.Join("lib", "database", "sqlite", "sqlc.yaml")) {
			choices = append(choices, search.Choice{Id: "lib/database/sqlite/sqlc.yaml", Description: "lib/database/sqlite/sqlc.yaml"})
		}

		choices = append(choices, search.Choice{Id: "other", Description: "other"})

		if options.Auto {
			options.SqlcYaml = choices[0].Id
		} else {
			options.SqlcYaml, err = singleselect.Sendf(choices, "where is your sqlc.yaml file located?")
		}
	}

	to := filepath.Dir(options.SqlcYaml)

	if _, err = exec.LookPath(options.Sqlc); err != nil && !files.IsFile(options.Sqlc) {
		if err = Sqlc(SqlcOptions{
			Sqlc:     options.Sqlc,
			Platform: options.Platform,
			Auto:     options.Auto,
		}); err != nil {
			return
		}
	}

	yaml := filepath.Join(to, "sqlc.yaml")

	if !files.IsFile(yaml) {
		return fmt.Errorf("%s not found", yaml)
	}

	var sqlc string
	if files.IsFile(options.Sqlc) {
		if sqlc, err = filepath.Rel(to, options.Sqlc); err != nil {
			return err
		}
	} else if sqlc, err = exec.LookPath(options.Sqlc); err != nil {
		sqlc = options.Sqlc
	}

	spin := spinner.New("generating queries")

	go spinner.Start(spin)
	if !messages.Command(to, os.Environ(), sqlc, "generate") {
		spinner.Stop(spin)
		err = errors.New("could not generate queries")
		return
	}
	spinner.Stop(spin)

	if err = FixImports(FixImportsOptions{Directory: to}); err != nil {
		return
	}

	messages.Success(
		"queries generated at database.Queries.*\n",
		to+"/queries.go",
	)

	return
}
