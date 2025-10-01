package csr

import "embed"

type Config struct {
	App     string
	Efs     embed.FS
	UseDisk bool
}
