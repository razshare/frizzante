package js

import (
	"fmt"
	"github.com/evanw/esbuild/pkg/api"
)

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
