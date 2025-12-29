package actions

import (
	"github.com/razshare/frizzante/cli/generate"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
)

func Configure(options ConfigureOptions) (err error) {
	if !files.IsFile(options.Bun) {
		if err = generate.Bun(generate.BunOptions{Bun: options.Bun}); err != nil {
			return
		}
	}
	if !files.IsFile(options.Air) {
		if err = generate.Air(generate.AirOptions{Air: options.Air}); err != nil {
			return
		}
	}
	if err = Install(InstallOptions{Go: options.Go, Bun: options.Bun}); err != nil {
		return
	}
	if err = Package(PackageOptions{Bun: options.Bun}); err != nil {
		return
	}
	messages.Success("project configured")
	return
}
