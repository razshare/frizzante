package actions

import (
	"embed"
)

type DevOptions struct {
	Go          string
	Air         string
	Bun         string
	Efs         embed.FS
	Tags        string
	Interactive bool
}
