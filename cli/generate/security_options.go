package generate

import "embed"

type SecurityOptions struct {
	Strict bool
	Efs    embed.FS
}
