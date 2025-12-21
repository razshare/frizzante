package main

import (
	"embed"
	"errors"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/cli/apps"
	"github.com/razshare/frizzante/cli/extensions"
	"github.com/razshare/frizzante/tui/configs"
	"github.com/razshare/frizzante/tui/messages"
	flag "github.com/spf13/pflag"
)

//go:embed logo.txt
//go:embed version
//go:embed internal/additions/**
//go:embed internal/project/**
//go:embed internal/project/lib/core/views/render/.gitignore
//go:embed internal/project/.vscode/launch.json
//go:embed internal/project/app/.gitignore
//go:embed internal/project/app/.npmrc
//go:embed internal/project/app/.prettierrc
//go:embed internal/project/app/.prettierignore
var efs embed.FS

func main() {
	flag.Parse()
	messages.Prefix = configs.Styles.Menu.PaddingRight(1).Render("│")
	query := strings.Join(flag.Args(), " ")
	app := apps.App{
		Efs: efs,
		Modifiers: apps.Modifiers{
			Go:                       flag.StringP("go", "", "go"+extensions.Find(), "sets the go binary location"),
			Air:                      flag.StringP("air", "", filepath.Join(".gen", "air", "air"+extensions.Find()), "sets the air binary location"),
			Bun:                      flag.StringP("bun", "", filepath.Join(".gen", "bun", "bun"+extensions.Find()), "sets the bun binary location"),
			Sqlc:                     flag.StringP("sqlc", "", filepath.Join(".gen", "sqlc", "sqlc"+extensions.Find()), "sets the sqlc binary location"),
			SqlcYaml:                 flag.StringP("sqlc-yaml", "", "", "sets the sqlc configuration file location"),
			Tags:                     flag.StringP("tags", "", "", "sets build tags"),
			DatabaseConnectionString: flag.StringP("database-connection-string", "", "", "sets the database connection string or file"),
			DatabaseType:             flag.StringP("database-type", "", "sqlc", "sets the type of database (currently only sqlite is supported)"),
		},
	}
	if err := cli.StartApp(cli.StartAppOptions{App: app, Query: query}); err != nil {
		if !errors.Is(err, tea.ErrInterrupted) {
			messages.Fatal(err)
		}
		os.Exit(0)
	}
}
