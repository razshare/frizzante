package action

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
)

func PackageWatch(options PackageWatchOptions) (err error) {
	if err = Touch(TouchOptions{App: options.App}); err != nil {
		return
	}

	var bun string
	if files.IsFile(options.Bun) {
		if bun, err = filepath.Rel(options.App, options.Bun); err != nil {
			return
		}
	} else if bun, err = exec.LookPath(options.Bun); err != nil {
		bun = options.Bun
	}

	ssr := exec.Command(bun, "x", "vite", "build", "--logLevel=info", "--outDir=dist", "--emptyOutDir=false", "--watch", "--ssr=app.server.ts")
	ssr.Dir = options.App
	ssr.Env = append(os.Environ(), "DEV=1")
	//ssr.Stderr = os.Stderr
	//ssr.Stdout = os.Stdout
	//ssr.Stdin = os.Stdin
	if err = ssr.Start(); err != nil {
		if ssr.Err != nil {
			messages.Error(ssr.Err.Error())
		}
		return
	}

	//messages.Success("vite server watcher launched")

	csr := exec.Command(bun, "x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=false", "--watch")
	csr.Dir = options.App
	csr.Env = append(os.Environ())
	//csr.Stderr = os.Stderr
	//csr.Stdout = os.Stdout
	//csr.Stdin = os.Stdin
	if err = csr.Start(); err != nil {
		if csr.Err != nil {
			messages.Error(csr.Err.Error())
		}
		return
	}

	//messages.Success("vite client watcher launched")

	if err = csr.Wait(); err != nil {
		return
	}

	return ssr.Wait()
}
