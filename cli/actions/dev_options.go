package actions

import (
	"embed"
)

type DevOptions struct {
	Efs    embed.FS
	Go     string
	Air    string
	Bun    string
	Strict bool
}
