package generate

import "embed"

type AirConfigOptions struct {
	Tags []string
	Efs  embed.FS
}
