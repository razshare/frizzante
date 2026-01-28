package actions

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/razshare/frizzante/tui/messages"
)

func PackageWatch(options PackageWatchOptions) (err error) {
	var group sync.WaitGroup
	group.Go(func() {
		if !messages.Command(messages.CommandOptions{
			DirectoryName: "app",
			Environment:   append(os.Environ(), "DEV=1"),
			Program:       options.Bun,
			Args:          []string{"x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=true", "--watch"},
		}) {
			messages.Error("vite failed to generate dist/client/*")
		}
	})
	group.Go(func() {
		if !messages.Command(messages.CommandOptions{
			DirectoryName: "app",
			Environment:   append(os.Environ(), "DEV=1"),
			Program:       options.Bun,
			Args:          []string{"x", "vite", "build", "--logLevel=info", "--outDir=dist/server", "--emptyOutDir=true", "--watch", "--ssr=app.server.ts"},
		}) {
			messages.Error("vite failed to generate dist/server/app.server.js")
		}
	})
	group.Go(func() {
		if !messages.Command(messages.CommandOptions{
			DirectoryName: "app",
			Environment:   append(os.Environ(), "DEV=1"),
			Program:       filepath.Join("app", "node_modules", ".bin", "esbuild"),
			Args:          []string{"--bundle", "--watch", "--outfile=dist/server/app.server.cjs", "--format=cjs", "--allow-overwrite", "dist/server/app.server.js"},
		}) {
			messages.Error("esbuild failed to convert dist/server/app.server.js into dist/server/app.server.cjs")
		}
	})
	group.Wait()
	return
}
