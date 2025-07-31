package main

import (
	"embed"
	"github.com/razshare/frizzante/cli"
)

//go:embed version
//go:embed app/frizzante
//go:embed sqlc.yaml
//go:embed queries.sql
//go:embed schema.sql
var cliEfs embed.FS
var frizzante = &cli.Cli{Efs: cliEfs}

func main() { frizzante.OnStart() }
