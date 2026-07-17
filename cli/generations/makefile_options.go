package generations

import "embed"

type MakefileOptions struct {
	Strict bool
	Efs    embed.FS
}
