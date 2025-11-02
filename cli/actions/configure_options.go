package actions

import (
	"embed"

	"github.com/razshare/frizzante/platforms"
)

type ConfigureOptions struct {
	Go       string
	Air      string
	Bun      string
	Efs      embed.FS
	Platform platforms.Platform
	Auto     bool
}
