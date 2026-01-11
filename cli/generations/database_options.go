package generations

import "embed"

type DatabasesOptions struct {
	Strict bool
	Go     string
	Type   string
	Sqlc   string
	Efs    embed.FS
}
