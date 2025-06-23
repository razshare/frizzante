package frizzante

import "embed"

//go:embed app/dist
//go:embed makefile
var efs embed.FS
