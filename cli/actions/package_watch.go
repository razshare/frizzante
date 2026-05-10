package actions

import (
	"errors"
	"os"
	"sync"

	"github.com/razshare/frizzante/tui/messages"
)

func PackageWatch(options PackageWatchOptions) (err error) {
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
			Environment:   os.Environ(),
			DirectoryName: "app",
			Program:       options.Bun,
			Args:          []string{"x", "vite", "build", "--logLevel=info", "--outDir=dist/server", "--emptyOutDir=true", "--ssr=app.server.ts", "--watch"},
		}) {
			messages.Error("could not build server bundle")
		}
	})
	group.Wait()
	return
}
