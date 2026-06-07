package main

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/cli/apps"
	"github.com/razshare/frizzante/cli/extensions"
	"github.com/razshare/frizzante/cli/menus"
	"github.com/razshare/frizzante/tui/messages"
	flag "github.com/spf13/pflag"
)

//go:embed logo.txt
//go:embed version
//go:embed internal/project
//go:embed internal/project/.air.toml
//go:embed internal/project/.zed/debug.json
//go:embed internal/project/.vscode/launch.json
//go:embed internal/project/.vscode/tasks.json
//go:embed internal/project/app/.gitignore
//go:embed internal/project/app/.npmrc
//go:embed internal/project/app/.prettierrc
//go:embed internal/project/app/.prettierignore
var efs embed.FS

func main() {
	app := apps.App{
		Efs:          efs,
		Go:           flag.StringP("go", "", fmt.Sprintf("go%s", extensions.Find()), "sets the go binary file location"),
		Air:          flag.StringP("air", "", filepath.Join(".gen", "air", fmt.Sprintf("air%s", extensions.Find())), "sets the air binary file location"),
		Bun:          flag.StringP("bun", "", filepath.Join(".gen", "bun", fmt.Sprintf("bun%s", extensions.Find())), "sets the bun binary file location"),
		Sqlc:         flag.StringP("sqlc", "", filepath.Join(".gen", "sqlc", fmt.Sprintf("sqlc%s", extensions.Find())), "sets the sqlc binary file location"),
		SqlcYaml:     flag.StringP("sqlc-yaml", "", "", "sets the sqlc configuration file location"),
		Tags:         flag.StringP("tags", "", "", "sets build tags"),
		Database:     flag.StringP("database", "", "", "sets the database connection string; used for running migrations and snapshots"),
		DatabaseType: flag.StringP("database-type", "", "", "sets the type of database to use; currently only sqlite is supported"),
		Context:      flag.StringP("context", "", "js", "sets the context; used with <frizzante add>; currently only \"js\" context is supported"),
		Output:       flag.StringP("output", "o", filepath.Join(".gen", "bin", "app"+extensions.Find()), "sets the output name of the build process"),
		Strict:       flag.BoolP("strict", "s", false, "enables strict mode; program will stop if any required arguments or flags are missing; useful in ci/cd pipelines"),
		Incremental:  flag.BoolP("incremental", "", false, "enables incremental mode for svelte-check"),
	}
	flag.Parse()
	if _, err := menus.Activate(&menus.Main, app, flag.Args()); err != nil {
		if !errors.Is(err, tea.ErrInterrupted) {
			messages.Fatal(err)
		}
		os.Exit(0)
	}
}
