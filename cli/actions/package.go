package actions

import (
	"errors"
	"os"
	"path/filepath"
	"sync"

	"github.com/razshare/frizzante/tui/messages"
)

func Package(options PackageOptions) (err error) {
	var group sync.WaitGroup
	group.Go(func() {
		if !messages.Command(messages.CommandOptions{
			Environment:   os.Environ(),
			DirectoryName: "app",
			Program:       options.Bun,
			Args:          []string{"x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=false"},
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
			Args:          []string{"x", "vite", "build", "--logLevel=info", "--outDir=dist/server", "--emptyOutDir=true", "--ssr=app.server.ts"},
		}) {
			if err == nil {
				err = errors.New("could not build server bundle")
			}
			return
		}
		if err = os.Rename(filepath.Join("app", "dist", "server", "app.server.js"), filepath.Join("app", "dist", "server", "app.server.1.js")); err != nil {
			return
		}
		if !messages.Command(messages.CommandOptions{
			Environment:   os.Environ(),
			DirectoryName: "app",
			Program:       filepath.Join("app", "node_modules", ".bin", "esbuild"),
			Args:          []string{"--bundle", "--outfile=dist/server/app.server.js", "--format=cjs", "--allow-overwrite", "dist/server/app.server.1.js"},
		}) {
			if err == nil {
				err = errors.New("could not normalize server bundle")
			}
			return
		}
		err = os.RemoveAll(filepath.Join("app", "dist", "server", "app.server.js"))
	})
	group.Wait()
	messages.Success("app/dist generated")
	return
}
