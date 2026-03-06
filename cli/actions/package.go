package actions

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/razshare/frizzante/internal/project/lib/core/esbuild"
	"github.com/razshare/frizzante/internal/project/lib/core/files"
	"github.com/razshare/frizzante/tui/messages"
)

func Package(options PackageOptions) (err error) {
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
	var data []byte
	if data, err = os.ReadFile(filepath.Join("app", "dist", "server", "app.server.js")); err != nil {
		return
	}
	// currently svelte imports `node:crypto`, which will break our runtime,
	// so we need to strip it off from the bundle.
	// see issues #17762 and #17771:
	// https://github.com/sveltejs/svelte/issues/17762
	// https://github.com/sveltejs/svelte/issues/17771
	source := strings.Replace(string(data), `await obfuscated_import(`, "await import(", 1)
	emptyTsFileName := fmt.Sprintf(".%s%s", string(os.PathSeparator), filepath.Join(".gen", "empty.ts"))
	if !files.IsFile(emptyTsFileName) {
		if err = os.WriteFile(emptyTsFileName, []byte("export default {}"), os.ModePerm); err != nil {
			return
		}
	}
	var sourceStringBundled string
	if sourceStringBundled, err = esbuild.Bundle("app", api.FormatCommonJS, source, map[string]string{
		"node:crypto": strings.ReplaceAll(emptyTsFileName, string(os.PathSeparator), "/"),
	}); err != nil {
		return
	}
	if err = os.RemoveAll(emptyTsFileName); err != nil {
		return
	}
	if err = os.RemoveAll(filepath.Join("app", "dist", "server", "app.server.js")); err != nil {
		return
	}
	if err = os.WriteFile(filepath.Join("app", "dist", "server", "app.server.js"), []byte(sourceStringBundled), os.ModePerm); err != nil {
		return
	}
	messages.Success("project app packaged into ", filepath.Join("app", "dist"))
	return
}
