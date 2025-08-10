package js

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
)

// SetFunction sets a function.
func SetFunction(rt *goja.Runtime, n string, f Function) error {
	return rt.Set(n, f)
}

// SetFunctions sets a map of functions.
func SetFunctions(rt *goja.Runtime, fs map[string]Function) error {
	for n, cb := range fs {
		err := rt.Set(n, cb)
		if err != nil {
			return err
		}
	}

	return nil
}

// Bundle given a directory containing a node_modules subdirectory,
// bundles source code into a specific format.
func Bundle(d string, f api.Format, s string) (bundle string, err error) {
	result := api.Build(api.BuildOptions{
		Bundle: true,
		Format: f,
		Write:  false,
		Stdin: &api.StdinOptions{
			Contents:   s,
			ResolveDir: d,
		},
	})

	for _, buildError := range result.Errors {
		return "", fmt.Errorf("%s in %s:%d:%d", buildError.Text, buildError.Location.File, buildError.Location.Line, buildError.Location.Column)
	}

	return string(result.OutputFiles[0].Contents), nil
}
