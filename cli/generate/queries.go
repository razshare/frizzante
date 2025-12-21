package generate

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func Queries(options QueriesOptions) (err error) {

	if !files.IsFile(options.SqlcYaml) {
		err = fmt.Errorf("%s not found", options.SqlcYaml)
		return
	}

	baseDirectory := filepath.Dir(options.SqlcYaml)

	if _, err = exec.LookPath(options.Sqlc); err != nil && !files.IsFile(options.Sqlc) {
		if err = Sqlc(SqlcOptions{
			Sqlc:     options.Sqlc,
			Platform: options.Platform,
		}); err != nil {
			return
		}
	}

	var sqlc string
	if files.IsFile(options.Sqlc) {
		if sqlc, err = filepath.Rel(baseDirectory, options.Sqlc); err != nil {
			return err
		}
	} else if sqlc, err = exec.LookPath(options.Sqlc); err != nil {
		sqlc = options.Sqlc
	}

	spin := spinners.New("generating queries")

	go spinners.Start(spin)
	if !messages.Command(messages.CommandOptions{
		DirectoryName: baseDirectory,
		Environment:   os.Environ(),
		Program:       sqlc,
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

	messages.Success(filepath.Join(baseDirectory, "queries.go"))

	return
}
