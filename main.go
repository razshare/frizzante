package main

import (
	"embed"
	"errors"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/razshare/frizzante/cli"
	"github.com/razshare/frizzante/cli/app"
	"github.com/razshare/frizzante/tui/messages"
	flag "github.com/spf13/pflag"
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
var c = app.New()

func main() {
	flag.Parse()
	c.Efs = efs
	if err := cli.Start(c); err != nil {
		if !errors.Is(err, tea.ErrInterrupted) {
			messages.Fatal(err)
		}
		os.Exit(0)
	}
}
