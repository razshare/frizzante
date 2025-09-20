package action

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
)

func Package(options PackageOptions) (err error) {
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
	//ssr.Stderr = os.Stderr
	//ssr.Stdout = os.Stdout
	//ssr.Stdin = os.Stdin
	if err = ssr.Run(); err != nil {
		if ssr.Err != nil {
			messages.Error(ssr.Err.Error())
		}
		return
	}

	csr := exec.Command(bun, "x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=true")
	csr.Dir = options.App
	csr.Env = append(os.Environ())
	//csr.Stderr = os.Stderr
	//csr.Stdout = os.Stdout
	//csr.Stdin = os.Stdin
	if err = csr.Run(); err != nil {
		if csr.Err != nil {
			messages.Error(csr.Err.Error())
		}
		return
	}

	esb := exec.Command(filepath.Join("node_modules", ".bin", "esbuild"), "--bundle", "--outfile=dist/app.server.js", "--format=cjs", "--allow-overwrite", "dist/app.server.js")
	esb.Dir = options.App
	esb.Env = append(os.Environ())
	//esb.Stderr = os.Stderr
	//esb.Stdout = os.Stdout
	//esb.Stdin = os.Stdin
	if err = esb.Run(); err != nil {
		if esb.Err != nil {
			messages.Error(esb.Err.Error())
		}
		return
	}

	messages.Successf("%s generated", filepath.Join(options.App, "dist"))

	if !options.Prod && files.IsDirectory(filepath.Join("lib", "core", "view", "ssr")) {
		if files.IsDirectory(filepath.Join(options.App, "dist")) {
			if err = files.CopyDirectory(
				filepath.Join(options.App, "dist"),
				filepath.Join("lib", "core", "view", "ssr", "app", "dist"),
			); err != nil {
				return
			}

			messages.Successf("%s copied to  %s", filepath.Join(options.App, "dist"), filepath.Join("lib", "core", "view", "ssr"))
		}
	}

	return
}
