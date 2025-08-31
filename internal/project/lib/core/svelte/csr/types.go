package csr

import "embed"

type Config struct {
	Efs  embed.FS
	App  string
	Disk bool
}
