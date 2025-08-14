package ssr

import "embed"

type Config struct {
	Efs   embed.FS
	App   string
	Disk  bool
	Limit int
}
