package main

import (
	"embed"
	"errors"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/cli/app"
	"github.com/razshare/frizzante/tui/messages"
	flag "github.com/spf13/pflag"
)

//go:embed logo.txt
//go:embed version
//go:embed internal/additions/**
//go:embed internal/project/**
//go:embed internal/project/lib/core/view/ssr/.gitignore
//go:embed internal/project/app/.gitignore
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
