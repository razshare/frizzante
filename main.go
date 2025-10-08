package main

import (
	"embed"
	"errors"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/cli"
	app_ "github.com/razshare/frizzante/cli/app"
	"github.com/razshare/frizzante/tui/messages"
	flag "github.com/spf13/pflag"
)

//go:embed logo.txt
//go:embed version
//go:embed internal/additions/**
//go:embed internal/project/**
//go:embed internal/project/lib/core/views/render/.gitignore
//go:embed internal/project/app/.gitignore
//go:embed internal/project/app/.npmrc
//go:embed internal/project/app/.prettierrc
//go:embed internal/project/app/.prettierignore
var efs embed.FS
var frz = app_.New()

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
