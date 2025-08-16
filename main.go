package main

import (
	"embed"
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/cli/on"
	"github.com/razshare/frizzante/tui/messages"
	"os"
)

//go:embed clilogo.txt
//go:embed database.sqlite
//go:embed version
//go:embed template/lib
//go:embed template/app/frizzante
//go:embed sqlc.yaml
//go:embed queries.sql
//go:embed schema.sql
var efs embed.FS
var c = cli.New()

func main() {
	c.Efs = efs
	err := on.Start(c)
	if err != nil {
		if errors.Is(err, tea.ErrInterrupted) {
			os.Exit(0)
		}
		messages.Fatal(err)
	}
}
