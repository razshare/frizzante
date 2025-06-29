package main

import "embed"

//go:embed .github
//go:embed makefile
//go:embed templates/project/app/dist
var testEfs embed.FS
