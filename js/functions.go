package js

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
)

// SetFunction sets a function.
func SetFunction(runtime *goja.Runtime, name string, call Function) error {
	return runtime.Set(name, call)
}

// SetFunctions sets a map of functions.
func SetFunctions(runtime *goja.Runtime, calls map[string]Function) error {
	for name, call := range calls {
		if err := runtime.Set(name, call); err != nil {
			return err
		}
	}

	return nil
}

// Bundle bundles JavaScript source code into a specific format given a root directory containing node_modules.
func Bundle(root string, format api.Format, source string) (string, error) {
	result := api.Build(api.BuildOptions{
		Bundle: true,
		Format: format,
		Write:  false,
		Stdin: &api.StdinOptions{
			Contents:   source,
			ResolveDir: root,
		},
	})

	for _, err := range result.Errors {
		return "", fmt.Errorf("%s in %s:%d:%d", err.Text, err.Location.File, err.Location.Line, err.Location.Column)
	}

	return string(result.OutputFiles[0].Contents), nil
}
