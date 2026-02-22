package actions

import "embed"

type BuildOptions struct {
	Go     string
	Bun    string
	Tags   string
	Strict bool
	Efs    embed.FS
}
