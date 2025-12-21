package actions

import (
	"errors"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
)

func Package(options PackageOptions) (err error) {
	if err = Touch(TouchOptions{}); err != nil {
		return
	}

	var bun string
	if files.IsFile(options.Bun) {
		if bun, err = filepath.Rel("app", options.Bun); err != nil {
			return
		}
	} else if bun, err = exec.LookPath(options.Bun); err != nil {
		err = nil // we dont' care, we fallback to options.Bun
		bun = options.Bun
	}

	if !messages.Command(messages.CommandOptions{
		Environment:   os.Environ(),
		DirectoryName: "app",
		Program:       bun,
		Args:          []string{"x", "vite", "build", "--logLevel=info", "--outDir=dist", "--emptyOutDir=true", "--ssr=app.server.ts"},
	}) {
		err = errors.New("could not build server bundle")
		return
	}

	if !messages.Command(messages.CommandOptions{
		Environment:   os.Environ(),
		DirectoryName: "app",
		Program:       filepath.Join("node_modules", ".bin", "esbuild"),
		Args:          []string{"--bundle", "--outfile=dist/app.server.cjs", "--format=cjs", "--allow-overwrite", "dist/app.server.js"},
	}) {
		err = errors.New("could not normalize server bundle")
		return
	}

	if !messages.Command(messages.CommandOptions{
		Environment:   os.Environ(),
		DirectoryName: "app",
		Program:       bun,
		Args:          []string{"x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=true"},
	}) {
		err = errors.New("could not build client bundles")
		return
	}

	if err = os.RemoveAll(filepath.Join("app", "dist", "assets")); err != nil {
		return
	}

	if err = os.RemoveAll(filepath.Join("app", "dist", "app.server.js")); err != nil {
		return
	}

	messages.Success("app/dist generated")

	if !options.Production && files.IsDirectory(filepath.Join("lib", "core", "views", "render")) {
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

	if !options.Production && files.IsDirectory(filepath.Join("lib", "core", "send")) {
		if files.IsDirectory(filepath.Join("app", "dist")) {
			if err = files.CopyDirectory(
				filepath.Join("app", "dist"),
				filepath.Join("lib", "core", "send", "app", "dist"),
			); err != nil {
				return
			}

			messages.Success("app/dist copied to lib/core/send")
		}
	}

	return
}
