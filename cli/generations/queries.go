package generations

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/confirm"
	"github.com/razshare/frizzante/tui/inputs"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/search"
	"github.com/razshare/frizzante/tui/select_one"
	"github.com/razshare/frizzante/tui/spinners"
)

func Queries(options QueriesOptions) (err error) {
	sqlcYaml := options.SqlcYaml
	if sqlcYaml == "" {
		if options.Strict {
			err = errors.New("no sqlc.yaml file provided")
			return
		}
		var names []string
		if names, err = files.FindWithSuffix("lib", "sqlc.yaml"); err != nil {
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
			err = errors.New("cannot continue generating queries because sqlc is missing")
		}
		if err = Sqlc(SqlcOptions{Sqlc: options.Sqlc}); err != nil {
			return
		}
	}
	baseDirectory := filepath.Dir(sqlcYaml)
	spin := spinners.New("generating queries")
	go spinners.Start(spin)
	if !messages.Command(messages.CommandOptions{
		DirectoryName: baseDirectory,
		Environment:   os.Environ(),
		Program:       options.Sqlc,
		Args:          []string{"generate"},
	}) {
		spinners.Stop(spin)
		err = errors.New("could not generate queries")
		return
	}
	spinners.Stop(spin)
	if err = FixImports(FixImportsOptions{Directory: baseDirectory}); err != nil {
		return
	}
	messages.Success("queries generated")
	messages.Tip(
		"## usage example\n",
		"databases.Queries.FindTodosBySessionId(client.Request.Context(), \"some-session-id-123-...\")",
	)
	return
}
