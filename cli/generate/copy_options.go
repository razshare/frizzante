package generate

import "embed"

type CopyOptions struct {
	Ignore []string
	From   string
	To     string
	Efs    embed.FS
}
