package actions

import (
	"os"
	"path/filepath"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/razshare/frizzante/cli/extensions"
	"github.com/razshare/frizzante/cli/generations"
	"github.com/razshare/frizzante/internal/project/lib/core/esbuild"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
	"github.com/razshare/frizzante/tui/spinners"
)

func Build(options BuildOptions) (err error) {
	if !files.IsFile(".air.toml") {
		if err = generations.AirConfig(generations.AirConfigOptions{
			Efs: options.Efs,
		}); err != nil {
			return
		}
	}
	if !messages.Command(messages.CommandOptions{
		Environment:   os.Environ(),
		DirectoryName: "app",
		Program:       options.Bun,
		Args:          []string{"x", "vite", "build", "--logLevel=info", "--outDir=dist/client", "--emptyOutDir=false"},
	}) {
		messages.Error("could not build client bundles")
		return
	}
	if !messages.Command(messages.CommandOptions{
		Environment:   os.Environ(),
		DirectoryName: "app",
		Program:       options.Bun,
		Args:          []string{"x", "vite", "build", "--logLevel=info", "--outDir=dist/server", "--emptyOutDir=true", "--ssr=app.server.ts"},
	}) {
		messages.Error("could not build server bundle")
		return
	}
	// Vite build will convert app.server.ts into a new file "app.server.js",
	// however this new file could contain "import" statements and "require" calls.
	//
	// We need app.server.js be self-contained, meaning we
	// cannot allow it to have any "import" statements or "require" calls
	// because Goja does not implement a "require()" function.
	//
	// To solve this issue we need to run the app.server.js
	// through Esbuild, which will give us a self-container js script.
	var sourceData []byte
	if sourceData, err = os.ReadFile(filepath.Join("app", "dist", "server", "app.server.js")); err != nil {
		return
	}
	var sourceStringBundled string
	if sourceStringBundled, err = esbuild.Bundle("app", api.FormatCommonJS, string(sourceData)); err != nil {
		return
	}
	if err = os.RemoveAll(filepath.Join("app", "dist", "server", "app.server.js")); err != nil {
		return
	}
	if err = os.WriteFile(filepath.Join("app", "dist", "server", "app.server.js"), []byte(sourceStringBundled), os.ModePerm); err != nil {
		return
	}
	extension := extensions.Find()
	spin := spinners.New("building binary")
	go spinners.Start(spin)
	if !messages.Command(messages.CommandOptions{
		Environment: os.Environ(),
		Program:     options.Go,
		Args:        []string{"build", "-o=" + filepath.Join(".gen", "bin", "app"+extension), "."},
	}) {
		messages.Error("could not build go source code")
		return
	}
	spinners.Stop(spin)
	messages.Success("project built into ", filepath.Join(".gen", "bin", "app"+extension))
	return
}
