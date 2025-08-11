package codegen

import (
	"embed"
	"github.com/razshare/frizzante/gen/copy"
	"github.com/razshare/frizzante/gen/download"
)

var Generators = map[string]func(efs embed.FS){
	"air":   download.Air,
	"bun":   download.Bun,
	"sqlc":  download.Sqlc,
	"core":  copy.Core,
	"forms": copy.Forms,
	"links": copy.Links,
}
