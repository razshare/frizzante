package runtime

import (
	"main/lib/core/js"

	"github.com/dop251/goja"
)

// WithFunction sets a function.
func WithFunction(run *goja.Runtime, name string, call js.Function) error {
	return run.Set(name, call)
}
