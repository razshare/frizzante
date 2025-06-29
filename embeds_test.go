package main

import "embed"

//go:embed .github
//go:embed makefile
var tefs embed.FS
