package generations

import "embed"

type CoreOptions struct {
	Strict bool
	Efs    embed.FS
}
