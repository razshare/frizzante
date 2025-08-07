package js

import (
	"fmt"
	"github.com/dop251/goja"
	"github.com/evanw/esbuild/pkg/api"
)

// New creates a new JavaScript runtime.
func New() *JavaScript {
	return &JavaScript{
		Runtime: goja.New(),
	}
}

// SetFunction sets a function.
func SetFunction(js *JavaScript, n string, f Function) error {
	return js.Set(n, f)
}

// SetFunctions sets a map of functions.
func SetFunctions(js *JavaScript, fs map[string]Function) error {
	for n, cb := range fs {
		err := js.Set(n, cb)
		if err != nil {
			return err
		}
	}

	return nil
}

// Bundle bundles source code into a specific format.
func Bundle(root string, format api.Format, source string) (bundle string, err error) {
	result := api.Build(api.BuildOptions{
		Bundle: true,
		Format: format,
		Write:  false,
		Stdin: &api.StdinOptions{
			Contents:   source,
			ResolveDir: root,
		},
	})

	for _, buildError := range result.Errors {
		return "", fmt.Errorf("%s in %s:%d:%d", buildError.Text, buildError.Location.File, buildError.Location.Line, buildError.Location.Column)
	}

	return string(result.OutputFiles[0].Contents), nil
}
