package codegen

import "github.com/razshare/frizzante/cli"

var Functions = map[string]func(c *cli.Cli, base string) error{
	"air":      Air,
	"bun":      Bun,
	"session":  Session,
	"database": Database,
	"queries":  Queries,
	"core":     Core,
	"forms":    Forms,
	"links":    Links,
}
