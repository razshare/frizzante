package generators

import (
	"embed"
	"github.com/razshare/frizzante/codegen/cmd"
	"github.com/razshare/frizzante/codegen/copy"
	"github.com/razshare/frizzante/codegen/copy_modded"
	"github.com/razshare/frizzante/codegen/download"
)

var Functions = map[string]func(efs embed.FS){
	"air":      download.Air,
	"bun":      download.Bun,
	"sqlc":     download.Sqlc,
	"session":  copy_modded.Session,
	"database": copy_modded.Database,
	"queries":  cmd.Queries,
	"core":     copy.Core,
	"forms":    copy.Forms,
	"links":    copy.Links,
}
