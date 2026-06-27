package databases

import "embed"

//go:embed source.sqlite
//go:embed migrations
var Efs embed.FS
