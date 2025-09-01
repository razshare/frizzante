package main

import (
	"embed"
	"errors"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/internal/cli"
	"github.com/razshare/frizzante/internal/cli/app"
	"github.com/razshare/frizzante/tui/messages"
	flag "github.com/spf13/pflag"
)

//go:embed logo.txt
//go:embed version
//go:embed internal/project.zip
var efs embed.FS
var frz = app.New()

func main() {
	flag.Parse()
	frz.Efs = efs
	if err := cli.Start(frz); err != nil {
		if !errors.Is(err, tea.ErrInterrupted) {
			messages.Fatal(err)
		}
		os.Exit(0)
	}
}
