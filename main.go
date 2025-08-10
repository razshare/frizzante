package main

import (
	"embed"
	"github.com/razshare/frizzante/cli"
)

//go:embed clilogo.txt
//go:embed database.sqlite
//go:embed version
//go:embed app/frizzante
//go:embed sqlc.yaml
//go:embed queries.sql
//go:embed schema.sql
var efs embed.FS

func main() { cli.OnStart(efs) }
