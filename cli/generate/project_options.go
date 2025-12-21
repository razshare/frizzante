package generate

import "embed"

type ProjectOptions struct {
	Name        string
	Go          string
	Efs         embed.FS
	Interactive bool
}
