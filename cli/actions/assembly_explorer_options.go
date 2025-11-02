package actions

import "github.com/razshare/frizzante/platforms"

type AssemblyExplorerOptions struct {
	Go       string
	Bun      string
	Tags     []string
	Platform platforms.Platform
	Auto     bool
}
