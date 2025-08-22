package action

import (
	"github.com/razshare/frizzante/tui/messages"
	"os"
	"os/exec"
	"path/filepath"
)

func Pkg(opts PkgOptions) error {
	err := Touch(TouchOptions{App: opts.App})
	if err != nil {
		return err
	}

	bun, err := filepath.Rel(opts.App, opts.Bun)
	if err != nil {
		return err
	}

	ssr := exec.Command(bun, "x", "vite", "build", "--logLevel=info", "--outDir=dist", "--emptyOutDir=true", "--ssr=frizzante/core/scripts/server.ts")
	ssr.Dir = opts.App
	ssr.Env = append(os.Environ())
	ssr.Stderr = os.Stderr
	ssr.Stdout = os.Stdout
	ssr.Stdin = os.Stdin
	err = ssr.Run()
	if err != nil {
		return err
	}

	csr := exec.Command(bun, "x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=true")
	csr.Dir = opts.App
	csr.Env = append(os.Environ())
	csr.Stderr = os.Stderr
	csr.Stdout = os.Stdout
	csr.Stdin = os.Stdin
	err = csr.Run()
	if err != nil {
		return err
	}

	esb := exec.Command(filepath.Join("node_modules", ".bin", "esbuild"), "--bundle", "--outfile=dist/server.js", "--format=cjs", "--allow-overwrite", "dist/server.js")
	esb.Dir = opts.App
	esb.Env = append(os.Environ())
	esb.Stderr = os.Stderr
	esb.Stdout = os.Stdout
	esb.Stdin = os.Stdin
	err = esb.Run()
	if err != nil {
		return err
	}

	messages.Success("project app package generated in ", filepath.Join(opts.App, "dist"))

	return nil
}
