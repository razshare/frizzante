package actions

import (
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"time"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
)

func PackageWatch(options PackageWatchOptions) (err error) {
	if err = Touch(TouchOptions{}); err != nil {
		return
	}

	var bun string
	if files.IsFile(options.Bun) {
		if bun, err = filepath.Rel("app", options.Bun); err != nil {
			return
		}
	} else if bun, err = exec.LookPath(options.Bun); err != nil {
		bun = options.Bun
	}

	var group sync.WaitGroup
	group.Go(func() {
		messages.Command(
			"app",
			append(os.Environ(), "DEV=1"),
			bun, "x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=false", "--watch",
		)
	})
	group.Go(func() {
		messages.Command(
			"app",
			os.Environ(),
			bun, "x", "vite", "build", "--logLevel=info", "--outDir=dist", "--emptyOutDir=false", "--watch", "--ssr=app.server.ts",
		)
	})
	group.Go(func() {
		time.Sleep(time.Second)
		messages.Command(
			"app",
			append(os.Environ(), "DEV=1"),
			filepath.Join("node_modules", ".bin", "esbuild"),
			"--bundle", "--watch", "--outfile=dist/app.server.cjs", "--format=cjs", "--allow-overwrite", "dist/app.server.js",
		)
	})
	group.Wait()
	return err
}
