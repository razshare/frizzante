package codegen

import (
	"embed"
)

var Functions = map[string]func(efs embed.FS, base string) error{
	"air":      Air,
	"bun":      Bun,
	"session":  Session,
	"database": Database,
	"queries":  Queries,
	"core":     Core,
	"forms":    Forms,
	"links":    Links,
}
