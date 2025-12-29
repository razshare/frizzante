package generate

import "embed"

type LinksOptions struct {
	Strict bool
	Efs    embed.FS
}
