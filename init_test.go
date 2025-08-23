package main

import "embed"

//go:embed .github
//go:embed makefile
//go:embed app/dist
//go:embed template/project.zip
//go:embed template/lib
var Tefs embed.FS
