package generations

import "embed"

type CopyOptions struct {
	Strict bool
	Ignore []string
	From   string
	To     string
	Efs    embed.FS
}
