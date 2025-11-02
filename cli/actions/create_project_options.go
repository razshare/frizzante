package actions

import "embed"

type CreateProjectOptions struct {
	Name string
	Go   string
	Air  string
	Bun  string
	Efs  embed.FS
}
