package generate

import (
	"embed"

	"github.com/razshare/frizzante/platforms"
)

type DatabaseOptions struct {
	Generate string
	Go       string
	Type     string
	Sqlc     string
	Efs      embed.FS
	Platform platforms.Platform
	Auto     bool
}
