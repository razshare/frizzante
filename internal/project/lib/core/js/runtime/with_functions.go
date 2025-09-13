package runtime

import (
	"github.com/dop251/goja"
	"github.com/razshare/frizzante/internal/project/lib/core/js"
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
