package main

import (
	"embed"
	frizzanteCli "github.com/razshare/frizzante/cli"
)

//go:embed version
//go:embed app/frizzante
var efs embed.FS
var cli = frizzanteCli.Cli{Efs: efs}

func main() { cli.Start() }
