package runtime

import (
	"github.com/dop251/goja"
	"github.com/razshare/frizzante/js"
)

// WithFunction sets a function.
func WithFunction(run *goja.Runtime, name string, call js.Function) error {
	return run.Set(name, call)
}
