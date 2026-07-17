package generations

import "embed"

type DatabasesOptions struct {
	Strict bool
	Efs    embed.FS
}
