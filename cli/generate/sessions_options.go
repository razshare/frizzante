package generate

import "embed"

type SessionsOptions struct {
	Type string
	Efs  embed.FS
	Auto bool
}
