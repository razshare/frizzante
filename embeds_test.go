package main

import "embed"

//go:embed .github
//go:embed makefile
//go:embed app/dist
var testEfs embed.FS
