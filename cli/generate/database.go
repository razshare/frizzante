package generate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func Database(options DatabaseOptions) (err error) {
	databaseType := strings.ToLower(options.Type)
	toDirectoryName := filepath.Join("lib", databaseType, "databases")

	fromDirectoryName := fmt.Sprintf("internal/additions/lib/%s/databases", databaseType)
	if err = Copy(CopyOptions{
		From: fromDirectoryName,
		To:   toDirectoryName,
		Efs:  options.Efs,
	}); err != nil {
		return
	}

	if databaseType == "sqlite" {
		spin := spinners.New("adding github.com/mattn/go-sqlite3")

		go spinners.Start(spin)
		if !messages.Command(messages.CommandOptions{
			Environment: os.Environ(),
			Program:     options.Go,
			Args:        []string{"get", "github.com/mattn/go-sqlite3"},
		}) {
			spinners.Stop(spin)
			err = errors.New("could not add github.com/mattn/go-sqlite3")
			return
		}

		spinners.Stop(spin)

		messages.Success("sqlite database is ready")

		if err = FixImports(FixImportsOptions{Directory: toDirectoryName}); err != nil {
			return
		}
	}

	return
}
