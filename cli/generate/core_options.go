package generate

import "embed"

type CoreOptions struct {
	Efs  embed.FS
	Auto bool
}
