package actions

import (
	"errors"
	"os"
	"path/filepath"

	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
)

func Package(options PackageOptions) (err error) {
	if err = Touch(TouchOptions{}); err != nil {
		return
	}
	if err = os.RemoveAll(filepath.Join("app", "dist")); err != nil {
		return
	}
	if !messages.Command(messages.CommandOptions{
		Environment:   os.Environ(),
		DirectoryName: "app",
		Program:       options.Bun,
		Args:          []string{"x", "vite", "build", "--logLevel=info", "--outDir=dist", "--emptyOutDir=false", "--ssr=app.server.ts"},
	}) {
		err = errors.New("could not build server bundle")
		return
	}
	if !messages.Command(messages.CommandOptions{
		Environment:   os.Environ(),
		DirectoryName: "app",
		Program:       filepath.Join("app", "node_modules", ".bin", "esbuild"),
		Args:          []string{"--bundle", "--outfile=dist/app.server.cjs", "--format=cjs", "--allow-overwrite", "dist/app.server.js"},
	}) {
		err = errors.New("could not normalize server bundle")
		return
	}
	if !messages.Command(messages.CommandOptions{
		Environment:   os.Environ(),
		DirectoryName: "app",
		Program:       options.Bun,
		Args:          []string{"x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=false"},
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
	if !options.Production && files.IsDirectory(filepath.Join("lib", "core", "ssr")) {
		if files.IsDirectory(filepath.Join("app", "dist")) {
			if err = files.CopyDirectory(
				filepath.Join("app", "dist"),
				filepath.Join("lib", "core", "ssr", "app", "dist"),
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
