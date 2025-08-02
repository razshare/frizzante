package javascript

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
func (javascript *JavaScript) SetFunction(name string, callback Function) error {
	return javascript.Set(name, callback)
}

// SetFunctions sets a map of functions.
func (javascript *JavaScript) SetFunctions(functions map[string]Function) error {
	for name, callback := range functions {
		err := javascript.Set(name, callback)
		if err != nil {
			return err
		}
	}

	return nil
}

// Bundle bundles source code into a specific format.
func Bundle(root string, format api.Format, sourceCode string) (bundle string, err error) {
	result := api.Build(api.BuildOptions{
		Bundle: true,
		Format: format,
		Write:  false,
		Stdin: &api.StdinOptions{
			Contents:   sourceCode,
			ResolveDir: root,
		},
	})

	for _, buildError := range result.Errors {
		return "", fmt.Errorf("%s in %s:%d:%d", buildError.Text, buildError.Location.File, buildError.Location.Line, buildError.Location.Column)
	}

	return string(result.OutputFiles[0].Contents), nil
}
