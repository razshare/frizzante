package actions

import "github.com/razshare/frizzante/platforms"

type BuildOptions struct {
	Go       string
	Bun      string
	Tags     []string
	Platform platforms.Platform
}
