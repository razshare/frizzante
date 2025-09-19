package ssr

import "embed"

type Config struct {
	App   string
	Efs   embed.FS
	Limit int
	Disk  bool
}
