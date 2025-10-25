package actions

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func Package(options PackageOptions) (err error) {
	spin := spinners.New(fmt.Sprintf("packaging javascript application in app/dist"))
	go spinners.Start(spin)
	defer spinners.Stop(spin)

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

	if !messages.Command(
		"app",
		os.Environ(),
		bun, "x", "vite", "build", "--logLevel=info", "--outDir=dist", "--emptyOutDir=true", "--ssr=app.server.ts") {
		return
	}

	if !messages.Command(
		"app",
		os.Environ(),
		bun, "x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=true") {
		return
	}

	if !messages.Command(
		"app",
		os.Environ(),
		filepath.Join("node_modules", ".bin", "esbuild"),
		"--bundle", "--outfile=dist/app.server.cjs", "--format=cjs", "--allow-overwrite", "dist/app.server.js") {
		return
	}

	if err = os.RemoveAll(filepath.Join("app", "dist", "assets")); err != nil {
		return
	}

	if err = os.RemoveAll(filepath.Join("app", "dist", "app.server.js")); err != nil {
		return
	}

	messages.Success("app/dist generated")

	if !options.Prod && files.IsDirectory(filepath.Join("lib", "core", "views", "render")) {
		if files.IsDirectory(filepath.Join("app", "dist")) {
			if err = files.CopyDirectory(
				filepath.Join("app", "dist"),
				filepath.Join("lib", "core", "views", "render", "app", "dist"),
			); err != nil {
				return
			}

			messages.Success("app/dist copied to lib/core/views/render")
		}
	}

	return
}
