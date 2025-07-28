package main

import (
	"embed"
	frizzanteCli "github.com/razshare/frizzante/cli"
)

//go:embed version
//go:embed app/frizzante
//go:embed sqlc.yaml
var cliEfs embed.FS
var cli = frizzanteCli.Cli{Efs: cliEfs}

func main() { cli.OnStart() }
