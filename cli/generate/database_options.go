package generate

import (
	"embed"

	"github.com/razshare/frizzante/platforms"
)

type DatabaseOptions struct {
	Go       string
	Type     string
	Sqlc     string
	Efs      embed.FS
	Platform platforms.Platform
}
