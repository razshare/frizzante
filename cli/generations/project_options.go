package generations

import "embed"

type ProjectOptions struct {
	Interactive bool
	Name        string
	Efs         embed.FS
}
