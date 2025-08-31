package action

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/files"
	"github.com/razshare/frizzante/internal/tui/messages"
)

func Package(options PacakgeOptions) (err error) {
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

	ssr := exec.Command(bun, "x", "vite", "build", "--logLevel=info", "--outDir=dist", "--emptyOutDir=true", "--ssr=app.server.ts")
	ssr.Dir = options.App
	ssr.Env = append(os.Environ())
	ssr.Stderr = os.Stderr
	ssr.Stdout = os.Stdout
	ssr.Stdin = os.Stdin
	if err = ssr.Run(); err != nil {
		return
	}

	csr := exec.Command(bun, "x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=true")
	csr.Dir = options.App
	csr.Env = append(os.Environ())
	csr.Stderr = os.Stderr
	csr.Stdout = os.Stdout
	csr.Stdin = os.Stdin
	if err = csr.Run(); err != nil {
		return err
	}

	esb := exec.Command(filepath.Join("node_modules", ".bin", "esbuild"), "--bundle", "--outfile=dist/app.server.js", "--format=cjs", "--allow-overwrite", "dist/app.server.js")
	esb.Dir = options.App
	esb.Env = append(os.Environ())
	esb.Stderr = os.Stderr
	esb.Stdout = os.Stdout
	esb.Stdin = os.Stdin
	if err = esb.Run(); err != nil {
		return
	}

	messages.Success("project app package generated in ", filepath.Join(options.App, "dist"))

	if files.IsDirectory(filepath.Join("lib", "core", "svelte", "ssr")) {
		if err = os.RemoveAll(filepath.Join("lib", "core", "svelte", "ssr", "app")); err != nil {
			return
		}

		if err = files.CopyDirectory(filepath.Join(options.App, "dist"), filepath.Join("lib", "core", "svelte", "ssr", "app", "dist")); err != nil {
			return
		}
	}

	return
}
