package main

import (
	"embed"
	"errors"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/cli/apps"
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
var app = apps.New()

func main() {
	flag.Parse()
	messages.Prefix = configs.Styles.Menu.PaddingRight(1).Render("│")
	app.Efs = efs
	query := strings.Join(flag.Args(), " ")
	if err := cli.StartApp(cli.StartAppOptions{
		App:   *app,
		Query: query,
	}); err != nil {
		if !errors.Is(err, tea.ErrInterrupted) {
			messages.Fatal(err)
		}
		os.Exit(0)
	}
}
