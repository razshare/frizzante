package generations

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/v2/internal/project/lib/core/files"
	"github.com/razshare/frizzante/v2/tui/confirm"
	"github.com/razshare/frizzante/v2/tui/inputs"
	"github.com/razshare/frizzante/v2/tui/messages"
	"github.com/razshare/frizzante/v2/tui/search"
	"github.com/razshare/frizzante/v2/tui/select_one"
	"github.com/razshare/frizzante/v2/tui/spinners"
)

func Schema(options SchemaOptions) (err error) {
	sqlcYaml := options.SqlcYaml
	if sqlcYaml == "" {
		if options.Strict {
			err = errors.New("no sqlc.yaml file provided")
			return
		}
		var names []string
		if names, err = files.FindWithSuffix(".", "sqlc.yaml"); err != nil {
			return
		}
		choices := make([]search.Choice, len(names))
		for index, name := range names {
			choices[index] = search.Choice{Id: name}
		}
		choices = append(choices, search.Choice{Id: "other", Description: "other"})
		sqlcYaml, err = select_one.Sendf(choices, "where is your sqlc.yaml file located?")
		if sqlcYaml == "other" {
			if sqlcYaml, err = inputs.Send("where is your sqlc.yaml file located?"); err != nil {
				return
			}
		}
	}
	if !files.IsFile(options.Sqlc) {
		if options.Strict {
			err = fmt.Errorf("%s is missing", options.Sqlc)
			return
		}
		var yesInstall bool
		if yesInstall, err = confirm.Sendf(true, "%s is missing. Install?", options.Sqlc); err != nil {
			return
		}
		if !yesInstall {
			err = errors.New("cannot continue generating schema because sqlc is missing")
		}
		if err = Sqlc(SqlcOptions{Sqlc: options.Sqlc}); err != nil {
			return
		}
	}
	baseDirectory := filepath.Dir(sqlcYaml)
	spin := spinners.New("generating schema")
	go spinners.Start(spin)
	if !messages.Command(messages.CommandOptions{
		DirectoryName: baseDirectory,
		Environment:   os.Environ(),
		Program:       options.Sqlc,
		Args:          []string{"generate"},
	}) {
		spinners.Stop(spin)
		err = errors.New("could not generate schema")
		return
	}
	spinners.Stop(spin)
	if err = FixImports(FixImportsOptions{Directory: baseDirectory}); err != nil {
		return
	}
	messages.Success("schema generated")
	messages.Tip(
		"## usage example\n",
		"databases.Queries.FindTodosBySessionId(http.Request.Context(), \"some-session-id-123-...\")",
	)
	return
}
