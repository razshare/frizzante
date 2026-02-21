package actions

import (
	"errors"
	"os"
	"sync"

	"github.com/razshare/frizzante/cli/generations"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
)

func Dev(options DevOptions) (err error) {
	if !files.IsFile(".air.toml") {
		if err = generations.AirConfig(generations.AirConfigOptions{
			Efs: options.Efs,
		}); err != nil {
			return
		}
	}
	var group sync.WaitGroup
	group.Go(func() {
		if !messages.Command(messages.CommandOptions{
			Environment:   os.Environ(),
			DirectoryName: "app",
			Program:       options.Bun,
			Args:          []string{"x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=false", "--watch"},
		}) {
			if err == nil {
				err = errors.New("could not build client bundles")
			}
			return
		}
	})
	group.Go(func() {
		if !messages.Command(messages.CommandOptions{
			Environment:   append(os.Environ(), "DEV=1"),
			DirectoryName: "app",
			Program:       options.Bun,
			Args:          []string{"x", "vite", "build", "--logLevel=info", "--outDir=dist/server", "--emptyOutDir=true", "--ssr=app.server.ts", "--watch"},
		}) {
			messages.Error("could not build server bundle")
		}
	})
	group.Go(func() {
		if !messages.Command(messages.CommandOptions{
			Environment: append(os.Environ(), "DEV=1"),
			Program:     options.Air,
		}) {
			messages.Error("air failed")
		}
	})
	group.Wait()
	return
}
