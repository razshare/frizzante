package actions

import (
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/cli/generations"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
)

func Configure(options ConfigureOptions) (err error) {
	if !files.IsFile(options.Bun) {
		if err = generations.Bun(generations.BunOptions{Bun: options.Bun}); err != nil {
			return
		}
	}
	if !files.IsFile(options.Air) {
		if err = generations.Air(generations.AirOptions{Air: options.Air}); err != nil {
			return
		}
	}
	if err = Install(InstallOptions{Go: options.Go, Bun: options.Bun}); err != nil {
		return
	}
	if err = Build(BuildOptions{Go: options.Go, Bun: options.Bun, Tags: "ci"}); err != nil {
		return
	}
	if !messages.Command(messages.CommandOptions{
		Environment: os.Environ(),
		Program:     filepath.Join(".gen", "bin", "app"),
	}) {
		return
	}
	messages.Success("project configured")
	return
}
