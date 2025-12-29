package actions

import (
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/razshare/frizzante/tui/messages"
)

func PackageWatch(options PackageWatchOptions) (err error) {
	if err = Touch(TouchOptions{}); err != nil {
		return
	}

	var group sync.WaitGroup
	group.Go(func() {
		messages.Command(messages.CommandOptions{
			DirectoryName: "app",
			Environment:   append(os.Environ(), "DEV=1"),
			Program:       options.Bun,
			Args:          []string{"x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=false", "--watch"},
		})
	})
	group.Go(func() {
		messages.Command(messages.CommandOptions{
			DirectoryName: "app",
			Environment:   os.Environ(),
			Program:       options.Bun,
			Args:          []string{"x", "vite", "build", "--logLevel=info", "--outDir=dist", "--emptyOutDir=false", "--watch", "--ssr=app.server.ts"},
		})
	})
	group.Go(func() {
		time.Sleep(time.Second)
		messages.Command(messages.CommandOptions{
			DirectoryName: "app",
			Environment:   append(os.Environ(), "DEV=1"),
			Program:       filepath.Join("app", "node_modules", ".bin", "esbuild"),
			Args:          []string{"--bundle", "--watch", "--outfile=dist/app.server.cjs", "--format=cjs", "--allow-overwrite", "dist/app.server.js"},
		})
	})
	group.Wait()
	return err
}
