package runtime

import (
	"main/lib/core/js"

	"github.com/dop251/goja"
)

// WithFunctions sets a map of functions.
func WithFunctions(run *goja.Runtime, calls map[string]js.Function) error {
	for name, call := range calls {
		if err := run.Set(name, call); err != nil {
			return err
		}
	}

	return nil
}
