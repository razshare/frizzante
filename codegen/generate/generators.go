package generate

import (
	"embed"
	"github.com/razshare/frizzante/codegen/copy"
	"github.com/razshare/frizzante/codegen/download"
	"github.com/razshare/frizzante/codegen/parse"
)

var Available = map[string]func(efs embed.FS){
	"air":     download.Air,
	"bun":     download.Bun,
	"sqlc":    download.Sqlc,
	"session": parse.Session,
	"core":    copy.Core,
	"forms":   copy.Forms,
	"links":   copy.Links,
}
