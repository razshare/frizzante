package actions

import "embed"

type CreateProjectOptions struct {
	Strict bool
	Name   string
	Efs    embed.FS
}
