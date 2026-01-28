package actions

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/razshare/frizzante/tui/messages"
)

func PackageWatch(options PackageWatchOptions) (err error) {
	var group sync.WaitGroup
	group.Go(func() {
		messages.Command(messages.CommandOptions{
			DirectoryName: "app",
			Environment:   os.Environ(),
			Program:       options.Bun,
			Args:          []string{"x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=true", "--watch"},
		})
	})
	group.Go(func() {
		messages.Command(messages.CommandOptions{
			DirectoryName: "app",
			Environment:   append(os.Environ(), "DEV=1"),
			Program:       options.Bun,
			Args:          []string{"x", "vite", "build", "--logLevel=info", "--outDir=dist/server", "--emptyOutDir=true", "--watch", "--ssr=app.server.ts"},
		})
	})
	group.Go(func() {
		time.Sleep(time.Second)
		messages.Command(messages.CommandOptions{
			DirectoryName: "app",
			Environment:   append(os.Environ(), "DEV=1"),
			Program:       filepath.Join("app", "node_modules", ".bin", "esbuild"),
			Args:          []string{"--bundle", "--watch", "--outfile=dist/server/app.server.cjs", "--format=cjs", "--allow-overwrite", "dist/server/app.server.js"},
		})
	})
	group.Wait()
	return
}
