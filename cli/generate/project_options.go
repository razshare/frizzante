package generate

import "embed"

type ProjectOptions struct {
	Name        string
	Efs         embed.FS
	Interactive bool
}
