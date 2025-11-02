package actions

import (
	"embed"

	"github.com/razshare/frizzante/platforms"
)

type GenerateOptions struct {
	Generation string
	Go         string
	Air        string
	Bun        string
	Sqlc       string
	SqlcYaml   string
	Database   string
	Tags       []string
	Efs        embed.FS
	Platform   platforms.Platform
	Auto       bool
}
